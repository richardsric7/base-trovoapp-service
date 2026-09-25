// Package services implements the liveness and readiness probes for the admin API.
//
// Two distinct questions are answered here, and the distinction matters:
//
//	Liveness  ("is the process alive?")  - no dependency calls, always cheap,
//	          always 200 while the process can serve traffic. Uptime monitoring
//	          polls this; a failure means the container should be restarted.
//	Readiness ("can it serve requests?") - checks every backing dependency and
//	          returns 503 naming whichever one is down. Docker's healthcheck
//	          uses this, so a container whose database connection is broken is
//	          never given traffic.
package services

import (
	"context"
	"net/http"
	"os"
	"sync"
	"time"

	"admin-panel-dashboard/internal/cache"
	"admin-panel-dashboard/internal/models"

	"gorm.io/gorm"
)

// Status values reported for the service overall and for each dependency.
const (
	StatusUp       = "up"
	StatusDown     = "down"
	StatusDegraded = "degraded"
	StatusSkipped  = "skipped"
)

// checkTimeout bounds every individual dependency check. Readiness must answer
// quickly even when a dependency is hanging rather than refusing connections -
// a probe that blocks is indistinguishable from a probe that fails, but takes
// far longer to tell you so.
const checkTimeout = 3 * time.Second

// startedAt is captured at package init so uptime is reported from process
// start rather than from the first request.
var startedAt = time.Now()

// buildVersion is overridable at link time:
//
//	go build -ldflags "-X admin-panel-dashboard/internal/components/health/services.buildVersion=$(git rev-parse --short HEAD)"
//
// It falls back to the APP_VERSION env var, then to "unknown", so an
// un-stamped build still reports something sensible.
var buildVersion string

// LivenessResponse is the body of GET /health.
type LivenessResponse struct {
	Status        string `json:"status" example:"up"`
	Service       string `json:"service" example:"admin-api"`
	Version       string `json:"version" example:"a1b2c3d"`
	UptimeSeconds int64  `json:"uptime_seconds" example:"3600"`
	Time          string `json:"time" example:"2026-09-08T12:00:00Z"`
}

// DependencyStatus is one row of the readiness report.
type DependencyStatus struct {
	Name       string `json:"name" example:"admin_db"`
	Status     string `json:"status" example:"up"`
	LatencyMs  int64  `json:"latency_ms" example:"4"`
	LatencyUs  int64  `json:"latency_us" example:"420"`
	Error      string `json:"error,omitempty"`
	Optional   bool   `json:"optional,omitempty"`
	SkipReason string `json:"skip_reason,omitempty"`
}

// ReadinessResponse is the body of GET /ready.
type ReadinessResponse struct {
	Status        string             `json:"status" example:"up"`
	Service       string             `json:"service" example:"admin-api"`
	Version       string             `json:"version" example:"a1b2c3d"`
	UptimeSeconds int64              `json:"uptime_seconds" example:"3600"`
	Time          string             `json:"time" example:"2026-09-08T12:00:00Z"`
	Dependencies  []DependencyStatus `json:"dependencies"`
}

// Service runs the probes against the handles the server already holds.
type Service struct {
	adminDB  *gorm.DB
	walletDB *gorm.DB
	p2pDB    *gorm.DB
	cache    *cache.RedisCache
	wallet   walletPinger
	client   *http.Client
}

// walletPinger is the small slice of ServiceLink the readiness check needs.
// Narrowing it to an interface keeps the check testable without a live wallet.
type walletPinger interface {
	BaseURL() string
}

// serviceLinkAdapter exposes the ServiceLink base URL through walletPinger.
type serviceLinkAdapter struct{ baseURL string }

func (a serviceLinkAdapter) BaseURL() string { return a.baseURL }

// New builds the probe service from the running server's dependencies. Any of
// them may be nil (Redis when caching is disabled, ServiceLink when unset) and
// the corresponding check is reported as skipped rather than failing.
func New(adminDB, walletDB, p2pDB *gorm.DB, gc *models.GlobalConfig) *Service {
	s := &Service{
		adminDB:  adminDB,
		walletDB: walletDB,
		p2pDB:    p2pDB,
		client:   &http.Client{Timeout: checkTimeout},
	}
	if gc != nil {
		s.cache = gc.Cache
		if gc.ServiceLink != nil && gc.ServiceLink.ApiBaseUrl != "" {
			s.wallet = serviceLinkAdapter{baseURL: gc.ServiceLink.ApiBaseUrl}
		}
	}
	return s
}

// Version resolves the build stamp, env override, or "unknown". Exported so
// crash reports carry the same build identifier that /health reports.
func Version() string {
	if buildVersion != "" {
		return buildVersion
	}
	if v := os.Getenv("APP_VERSION"); v != "" {
		return v
	}
	return "unknown"
}

// Liveness reports that the process is running. It deliberately touches
// nothing external: a database outage must not cause the orchestrator to
// restart an otherwise healthy process.
func (s *Service) Liveness() LivenessResponse {
	now := time.Now()
	return LivenessResponse{
		Status:        StatusUp,
		Service:       "admin-api",
		Version:       Version(),
		UptimeSeconds: int64(now.Sub(startedAt).Seconds()),
		Time:          now.UTC().Format(time.RFC3339),
	}
}

// Readiness checks every dependency in parallel and reports the aggregate.
// Returns the report and the HTTP status the handler should use: 200 when
// every required dependency is up, 503 otherwise.
func (s *Service) Readiness(ctx context.Context) (ReadinessResponse, int) {
	ctx, cancel := context.WithTimeout(ctx, checkTimeout)
	defer cancel()

	checks := []func(context.Context) DependencyStatus{
		func(c context.Context) DependencyStatus { return s.checkDB(c, "admin_db", s.adminDB, false) },
		func(c context.Context) DependencyStatus { return s.checkDB(c, "wallet_db", s.walletDB, false) },
		func(c context.Context) DependencyStatus { return s.checkDB(c, "p2p_db", s.p2pDB, false) },
		s.checkRedis,
		s.checkWallet,
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
		Service:       "admin-api",
		Version:       Version(),
		UptimeSeconds: int64(now.Sub(startedAt).Seconds()),
		Time:          now.UTC().Format(time.RFC3339),
		Dependencies:  results,
	}

	// A required dependency being down makes the service not ready. An
	// optional one being down is reported as degraded but still serves
	// traffic - losing the cache slows the API, it does not break it.
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
	return report, status
}

// checkDB verifies the pool can round-trip a ping to the database.
func (s *Service) checkDB(ctx context.Context, name string, db *gorm.DB, optional bool) (dep DependencyStatus) {
	// Named return: the deferred latency write below must land in the value
	// actually returned. With an unnamed return the result is copied out
	// before the defer runs, so every latency reported as 0.
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

// checkRedis pings the cache when caching is enabled. Redis is optional: the
// API runs without it, so its absence degrades rather than fails readiness.
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

// checkWallet confirms the Trovo Wallet API is reachable. Several admin
// endpoints proxy to it, so it is reported, but it is optional: the wallet
// being down must not take the admin dashboard out of rotation with it.
func (s *Service) checkWallet(ctx context.Context) (dep DependencyStatus) {
	dep = DependencyStatus{Name: "wallet_api", Optional: true}
	start := time.Now()
	defer func() {
		elapsed := time.Since(start)
		dep.LatencyMs = elapsed.Milliseconds()
		dep.LatencyUs = elapsed.Microseconds()
	}()

	if s.wallet == nil {
		dep.Status = StatusSkipped
		dep.SkipReason = "not configured"
		return dep
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.wallet.BaseURL(), nil)
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

	// Any HTTP response proves the service is listening and routable. The
	// status code is irrelevant here: the base URL may legitimately answer
	// 401 or 404, which still means "the wallet API is up".
	dep.Status = StatusUp
	return dep
}
