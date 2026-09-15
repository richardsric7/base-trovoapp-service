// Package health exposes liveness and readiness probes for the payment history
// engine.
//
// This service is a worker, not an API: it polls Base blocks/logs and writes
// payment history. That shape hides failure well - a stalled poll loop
// returns no errors to anyone, and nobody notices until a user asks why
// their transaction history stopped updating. So readiness here reports not
// just "can I reach my dependencies" but "is block processing actually advancing".
package health

import (
	"context"
	"net/http"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"trovo-wallet-payment-history-engine/internal/cache"

	"github.com/ethereum/go-ethereum/ethclient"
	"gorm.io/gorm"
)

// Status values reported for the service and each dependency.
const (
	StatusUp       = "up"
	StatusDown     = "down"
	StatusDegraded = "degraded"
	StatusSkipped  = "skipped"
	StatusStalled  = "stalled"
	StatusIdle     = "idle"
)

// checkTimeout bounds each dependency check so a hanging database cannot hang
// the probe itself.
const checkTimeout = 3 * time.Second

// stallThreshold is how long the stream may go without processing an operation
// before readiness reports it stalled.
//
// Quiet periods are normal on this network, so this is deliberately generous:
// the aim is to catch a stream that has died, not to alarm at every lull. A
// service that has genuinely processed nothing for fifteen minutes during
// business hours is worth looking at.
const stallThreshold = 15 * time.Minute

var (
	startedAt = time.Now()

	// buildVersion is set at link time; see the Dockerfile.
	buildVersion string

	// lastOperationNanos is the last time the stream processed an operation,
	// as a Unix nanosecond count. Held as an atomic so the worker goroutines
	// can update it without contending on a lock in their hot path.
	//
	// Uses the atomic.LoadInt64/StoreInt64 functions rather than the atomic.Int64
	// type, which needs Go 1.19 - this module targets 1.18.
	lastOperationNanos int64

	// operationsProcessed counts operations handled since start, which turns
	// "is it working" into a number an operator can watch move.
	operationsProcessed uint64

	// lastCursor is the most recent Stellar paging token written. Guarded by a
	// mutex because it is a string, read rarely and written often enough that
	// tearing would matter.
	cursorMu   sync.RWMutex
	lastCursor string

	// hasWork records whether the engine currently has anything to stream. The
	// engine watches specific wallets; with none registered there is nothing to
	// process, and zero operations is then the correct and healthy state rather
	// than evidence of a stall.
	//
	// Without this distinction a quiet environment reports itself degraded
	// forever, which is the kind of false alarm that teaches people to ignore
	// the indicator.
	hasWorkFlag int32
)

// SetHasWork is called by the stream monitors to say whether there is anything
// to process - that is, whether any wallets are being tracked.
func SetHasWork(has bool) {
	var v int32
	if has {
		v = 1
	}
	atomic.StoreInt32(&hasWorkFlag, v)
}

func hasWork() bool { return atomic.LoadInt32(&hasWorkFlag) == 1 }

// RecordOperation is called by the stream workers each time an operation is
// processed. It is the signal that the engine is alive in the way that matters:
// the process being up means little if the stream behind it has stopped.
func RecordOperation(cursor string) {
	atomic.StoreInt64(&lastOperationNanos, time.Now().UnixNano())
	atomic.AddUint64(&operationsProcessed, 1)
	if cursor != "" {
		cursorMu.Lock()
		lastCursor = cursor
		cursorMu.Unlock()
	}
}

// Version resolves the build stamp, the APP_VERSION env var, or "unknown".
func Version() string {
	if buildVersion != "" {
		return buildVersion
	}
	if v := os.Getenv("APP_VERSION"); v != "" {
		return v
	}
	return "unknown"
}

// DependencyStatus is one row of the readiness report.
type DependencyStatus struct {
	Name       string `json:"name" example:"wallet_db"`
	Status     string `json:"status" example:"up"`
	LatencyMs  int64  `json:"latency_ms" example:"4"`
	LatencyUs  int64  `json:"latency_us" example:"4210"`
	Error      string `json:"error,omitempty"`
	Optional   bool   `json:"optional,omitempty"`
	SkipReason string `json:"skip_reason,omitempty"`
}

// StreamStatus reports whether the engine is still making progress.
type StreamStatus struct {
	Status              string `json:"status" example:"up"`
	OperationsProcessed uint64 `json:"operations_processed" example:"14238"`
	LastOperationAt     string `json:"last_operation_at,omitempty"`
	SecondsSinceLastOp  int64  `json:"seconds_since_last_operation"`
	LastCursor          string `json:"last_cursor,omitempty"`
	Note                string `json:"note,omitempty"`
}

// LivenessResponse is the body of GET /health.
type LivenessResponse struct {
	Status        string `json:"status" example:"up"`
	Service       string `json:"service" example:"payment-history-engine"`
	Version       string `json:"version" example:"a1b2c3d"`
	UptimeSeconds int64  `json:"uptime_seconds" example:"3600"`
	Time          string `json:"time" example:"2026-09-10T12:00:00Z"`
}

// ReadinessResponse is the body of GET /ready.
type ReadinessResponse struct {
	Status        string             `json:"status" example:"up"`
	Service       string             `json:"service" example:"payment-history-engine"`
	Version       string             `json:"version" example:"a1b2c3d"`
	UptimeSeconds int64              `json:"uptime_seconds" example:"3600"`
	Time          string             `json:"time" example:"2026-09-10T12:00:00Z"`
	Stream        StreamStatus       `json:"stream"`
	Dependencies  []DependencyStatus `json:"dependencies"`
}

// Service holds the handles the probes check against.
type Service struct {
	db      *gorm.DB
	roachDB *gorm.DB
	cache   *cache.RedisCache
	rpcURL  string
	client  *http.Client
}

// New builds the probe service. Any handle may be nil; the matching check then
// reports "skipped" rather than failing. blockchainClient may be nil (its RPC
// URL is read directly from env, matching network.GetBlockchainClient's own
// resolution, since *ethclient.Client does not expose the URL it dialed).
func New(db, roachDB *gorm.DB, redisCache *cache.RedisCache, blockchainClient *ethclient.Client) *Service {
	s := &Service{
		db:      db,
		roachDB: roachDB,
		cache:   redisCache,
		client:  &http.Client{Timeout: checkTimeout},
	}
	if blockchainClient != nil {
		s.rpcURL = os.Getenv("BASE_RPC_URL")
		if s.rpcURL == "" {
			s.rpcURL = os.Getenv("EXPANSION_URL")
		}
	}
	return s
}

// Liveness reports that the process is running, touching nothing external.
func (s *Service) Liveness() LivenessResponse {
	now := time.Now()
	return LivenessResponse{
		Status:        StatusUp,
		Service:       "payment-history-engine",
		Version:       Version(),
		UptimeSeconds: int64(now.Sub(startedAt).Seconds()),
		Time:          now.UTC().Format(time.RFC3339),
	}
}

// stream reports on the engine's actual work.
func stream() StreamStatus {
	out := StreamStatus{
		OperationsProcessed: atomic.LoadUint64(&operationsProcessed),
		Status:              StatusUp,
	}

	cursorMu.RLock()
	out.LastCursor = lastCursor
	cursorMu.RUnlock()

	// Nothing to do is not the same as failing to do it. An engine with no
	// tracked wallets has no work, so no operations is correct.
	if !hasWork() {
		out.Status = StatusIdle
		out.Note = "no wallets are being tracked - nothing to process"
		return out
	}

	nanos := atomic.LoadInt64(&lastOperationNanos)
	if nanos == 0 {
		// Nothing processed yet. During startup this is expected; after a long
		// uptime with work available it means the stream never began, which is
		// worth saying rather than reporting a misleading "up".
		out.Status = StatusSkipped
		if time.Since(startedAt) > stallThreshold {
			out.Status = StatusStalled
			out.Note = "wallets are being tracked but nothing has been processed since start"
		} else {
			out.Note = "starting up - no operations processed yet"
		}
		return out
	}

	last := time.Unix(0, nanos)
	out.LastOperationAt = last.UTC().Format(time.RFC3339)
	out.SecondsSinceLastOp = int64(time.Since(last).Seconds())

	if time.Since(last) > stallThreshold {
		out.Status = StatusStalled
		out.Note = "wallets are being tracked but nothing has been processed recently - the stream may have stopped"
	}
	return out
}

// Readiness checks dependencies in parallel and reports stream progress.
//
// Returns 503 only when a required dependency is down. A stalled stream is
// reported as degraded rather than unready: restarting the container would not
// necessarily restart a stream that stopped because the network went quiet, and
// marking it unhealthy would produce a restart loop rather than a fix. The
// stalled status is there to be seen and alerted on, not acted on automatically.
func (s *Service) Readiness(ctx context.Context) (ReadinessResponse, int) {
	ctx, cancel := context.WithTimeout(ctx, checkTimeout)
	defer cancel()

	checks := []func(context.Context) DependencyStatus{
		func(c context.Context) DependencyStatus { return s.checkDB(c, "wallet_db", s.db, false) },
		func(c context.Context) DependencyStatus { return s.checkDB(c, "roach_db", s.roachDB, true) },
		s.checkRedis,
		s.checkBlockchainRPC,
	}

	results := make([]DependencyStatus, len(checks))
	var wg sync.WaitGroup
	for i, check := range checks {
		wg.Add(1)
		go func(i int, check func(context.Context) DependencyStatus) {
			defer wg.Done()
			results[i] = check(ctx)
		}(i, check)
	}
	wg.Wait()

	now := time.Now()
	report := ReadinessResponse{
		Status:        StatusUp,
		Service:       "payment-history-engine",
		Version:       Version(),
		UptimeSeconds: int64(now.Sub(startedAt).Seconds()),
		Time:          now.UTC().Format(time.RFC3339),
		Stream:        stream(),
		Dependencies:  results,
	}

	status := http.StatusOK
	for _, dep := range results {
		if dep.Status != StatusDown {
			continue
		}
		if dep.Optional {
			if report.Status == StatusUp {
				report.Status = StatusDegraded
			}
			continue
		}
		report.Status = StatusDown
		status = http.StatusServiceUnavailable
	}

	if report.Status == StatusUp && report.Stream.Status == StatusStalled {
		report.Status = StatusDegraded
	}
	return report, status
}

func (s *Service) checkDB(ctx context.Context, name string, db *gorm.DB, optional bool) (dep DependencyStatus) {
	// Named return so the deferred latency write lands in the returned value.
	dep = DependencyStatus{Name: name, Optional: optional}
	start := time.Now()
	defer func() {
		elapsed := time.Since(start)
		dep.LatencyMs = elapsed.Milliseconds()
		dep.LatencyUs = elapsed.Microseconds()
	}()

	if db == nil {
		dep.Status = StatusSkipped
		dep.SkipReason = "not configured"
		return dep
	}
	sqlDB, err := db.DB()
	if err != nil {
		dep.Status = StatusDown
		dep.Error = err.Error()
		return dep
	}
	if err := sqlDB.PingContext(ctx); err != nil {
		dep.Status = StatusDown
		dep.Error = err.Error()
		return dep
	}
	dep.Status = StatusUp
	return dep
}

func (s *Service) checkRedis(ctx context.Context) (dep DependencyStatus) {
	dep = DependencyStatus{Name: "redis", Optional: true}
	start := time.Now()
	defer func() {
		elapsed := time.Since(start)
		dep.LatencyMs = elapsed.Milliseconds()
		dep.LatencyUs = elapsed.Microseconds()
	}()

	if os.Getenv("ENABLE_CACHING") != "1" {
		dep.Status = StatusSkipped
		dep.SkipReason = "caching disabled"
		return dep
	}
	if s.cache == nil || s.cache.Client == nil {
		dep.Status = StatusSkipped
		dep.SkipReason = "not configured"
		return dep
	}
	if err := s.cache.Client.Ping(ctx).Err(); err != nil {
		dep.Status = StatusDown
		dep.Error = err.Error()
		return dep
	}
	dep.Status = StatusUp
	return dep
}

// checkBlockchainRPC confirms the Base RPC node is reachable. It is the
// source of every block/log this engine processes, so it is reported
// prominently - but marked optional, because a restart cannot fix an
// unreachable node and the engine resumes from its saved cursor once the
// node returns.
func (s *Service) checkBlockchainRPC(ctx context.Context) (dep DependencyStatus) {
	dep = DependencyStatus{Name: "base_rpc", Optional: true}
	start := time.Now()
	defer func() {
		elapsed := time.Since(start)
		dep.LatencyMs = elapsed.Milliseconds()
		dep.LatencyUs = elapsed.Microseconds()
	}()

	if s.rpcURL == "" {
		dep.Status = StatusSkipped
		dep.SkipReason = "not configured"
		return dep
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.rpcURL, nil)
	if err != nil {
		dep.Status = StatusDown
		dep.Error = err.Error()
		return dep
	}
	resp, err := s.client.Do(req)
	if err != nil {
		dep.Status = StatusDown
		dep.Error = err.Error()
		return dep
	}
	defer resp.Body.Close()

	// Any HTTP response proves the node is listening; the status code does not
	// matter for reachability.
	dep.Status = StatusUp
	return dep
}
