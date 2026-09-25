package observability

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"
)

// ---- Response types (also the Swagger DTOs) ---------------------------------

// DependencyView is one backing service a backend depends on.
type DependencyView struct {
	Name      string `json:"name" example:"admin_db"`
	Status    string `json:"status" example:"up"`
	LatencyMs int64  `json:"latency_ms" example:"12"`
	Error     string `json:"error,omitempty"`
	Optional  bool   `json:"optional,omitempty"`
}

// StreamView reports a worker service's progress.
//
// Only background workers report this - the payment history engine streams
// operations from the blockchain, and "is the process alive" says nothing
// about whether it is still doing its job. A request-driven API has no
// equivalent, so this is absent for those.
type StreamView struct {
	Status              string `json:"status" example:"up"`
	OperationsProcessed uint64 `json:"operations_processed" example:"14238"`
	SecondsSinceLastOp  int64  `json:"seconds_since_last_operation" example:"12"`
	LastOperationAt     string `json:"last_operation_at,omitempty"`
	Note                string `json:"note,omitempty"`
}

// ServiceHealth is one backend's health card.
type ServiceHealth struct {
	Name          string           `json:"name" example:"admin-api"`
	Status        string           `json:"status" example:"up"`
	Version       string           `json:"version" example:"a1b2c3d"`
	UptimeSeconds int64            `json:"uptime_seconds" example:"86400"`
	Dependencies  []DependencyView `json:"dependencies"`
	// Stream is set only for worker services that report one.
	Stream *StreamView `json:"stream,omitempty"`
	Error  string      `json:"error,omitempty"`
}

// HostHealth is the machine the platform runs on.
type HostHealth struct {
	CPUPercent    *float64 `json:"cpu_percent"`
	MemoryPercent *float64 `json:"memory_percent"`
	DiskPercent   *float64 `json:"disk_percent"`
	Available     bool     `json:"available"`
}

// HealthResponse is the payload behind the Health screen.
type HealthResponse struct {
	Services   []ServiceHealth `json:"services"`
	Host       HostHealth      `json:"host"`
	CheckedAt  string          `json:"checked_at"`
	GrafanaURL string          `json:"grafana_url,omitempty"`
}

// Issue is one grouped error from the crash reporter.
type Issue struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Culprit   string `json:"culprit,omitempty"`
	Level     string `json:"level" example:"error"`
	Status    string `json:"status" example:"unresolved"`
	Count     int64  `json:"count" example:"14"`
	App       string `json:"app" example:"admin-web"`
	FirstSeen string `json:"first_seen"`
	LastSeen  string `json:"last_seen"`
	Permalink string `json:"permalink,omitempty"`
}

// IssuesResponse is the payload behind the Issues screen.
type IssuesResponse struct {
	Issues       []Issue `json:"issues"`
	Available    bool    `json:"available"`
	Reason       string  `json:"reason,omitempty"`
	GlitchTipURL string  `json:"glitchtip_url,omitempty"`
}

// SeriesPoint is one point on a trend line.
type SeriesPoint struct {
	Time  string  `json:"time"`
	Value float64 `json:"value"`
}

// Series is one named line on a chart.
type Series struct {
	Name   string        `json:"name"`
	Points []SeriesPoint `json:"points"`
}

// TrendsResponse is the payload behind the Trends screen.
type TrendsResponse struct {
	RequestRate []Series `json:"request_rate"`
	ErrorRate   []Series `json:"error_rate"`
	// HandledFailures counts failures the services caught and handled rather
	// than returning as a 5xx - most of them in background workers that serve
	// no HTTP request at all, so ErrorRate above cannot see them however
	// badly a money path is breaking.
	HandledFailures []Series `json:"handled_failures"`
	LatencyP95      []Series `json:"latency_p95"`
	Window          string   `json:"window" example:"6h"`
	Available       bool     `json:"available"`
	Reason          string   `json:"reason,omitempty"`
}

// LogLine is one structured log entry.
type LogLine struct {
	Time       string `json:"time"`
	Service    string `json:"service"`
	Level      string `json:"level,omitempty"`
	Method     string `json:"method,omitempty"`
	Path       string `json:"path,omitempty"`
	Status     int    `json:"status,omitempty"`
	DurationMs int64  `json:"duration_ms,omitempty"`
	Message    string `json:"message"`
}

// TraceResponse is the payload behind the Trace screen.
type TraceResponse struct {
	RequestID string    `json:"request_id"`
	Lines     []LogLine `json:"lines"`
	Available bool      `json:"available"`
	Reason    string    `json:"reason,omitempty"`
}

// RecentErrorsResponse lists recent failing requests, so an admin has
// something to click rather than needing a request id to hand.
type RecentErrorsResponse struct {
	Lines     []LogLine `json:"lines"`
	Available bool      `json:"available"`
	Reason    string    `json:"reason,omitempty"`
}

// ---- Cache ------------------------------------------------------------------

// cacheTTL is short: these screens are near-real-time, but without any cache a
// room full of admins with the page open would multiply into a steady query
// load on Prometheus and Loki.
const cacheTTL = 15 * time.Second

type cacheEntry struct {
	value     any
	expiresAt time.Time
}

var (
	cacheMu sync.Mutex
	cache   = map[string]cacheEntry{}
)

// cached returns a memoised value, recomputing it when stale. Errors are not
// cached - a transient upstream failure should not persist for the full TTL.
func cached[T any](key string, build func() (T, error)) (T, error) {
	cacheMu.Lock()
	if e, ok := cache[key]; ok && time.Now().Before(e.expiresAt) {
		if v, ok := e.value.(T); ok {
			cacheMu.Unlock()
			return v, nil
		}
	}
	cacheMu.Unlock()

	v, err := build()
	if err != nil {
		return v, err
	}

	cacheMu.Lock()
	cache[key] = cacheEntry{value: v, expiresAt: time.Now().Add(cacheTTL)}
	cacheMu.Unlock()
	return v, nil
}

// ---- Health -----------------------------------------------------------------

// readyPayload mirrors what every backend's /ready endpoint returns.
type readyPayload struct {
	Status        string `json:"status"`
	Service       string `json:"service"`
	Version       string `json:"version"`
	UptimeSeconds int64  `json:"uptime_seconds"`
	Dependencies  []struct {
		Name      string `json:"name"`
		Status    string `json:"status"`
		LatencyMs int64  `json:"latency_ms"`
		Error     string `json:"error"`
		Optional  bool   `json:"optional"`
	} `json:"dependencies"`
	// Present only on worker services; absent on request-driven APIs.
	Stream *struct {
		Status              string `json:"status"`
		OperationsProcessed uint64 `json:"operations_processed"`
		SecondsSinceLastOp  int64  `json:"seconds_since_last_operation"`
		LastOperationAt     string `json:"last_operation_at"`
		Note                string `json:"note"`
	} `json:"stream"`
}

// Health polls every configured service's /ready endpoint in parallel and adds
// host metrics from Prometheus.
func (c Config) Health(ctx context.Context) HealthResponse {
	out := HealthResponse{
		Services:   make([]ServiceHealth, len(c.Services)),
		CheckedAt:  time.Now().UTC().Format(time.RFC3339),
		GrafanaURL: c.GrafanaURL,
	}

	var wg sync.WaitGroup
	for i, target := range c.Services {
		wg.Add(1)
		go func(i int, target ServiceTarget) {
			defer wg.Done()
			out.Services[i] = c.serviceHealth(ctx, target)
		}(i, target)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		out.Host = c.hostHealth(ctx)
	}()
	wg.Wait()

	return out
}

// serviceHealth fetches and translates one service's /ready response. A service
// that cannot be reached is reported as down with the reason, rather than
// omitted - an absent card would look like the service does not exist.
func (c Config) serviceHealth(ctx context.Context, target ServiceTarget) ServiceHealth {
	view := ServiceHealth{Name: target.Name, Status: "down", Dependencies: []DependencyView{}}

	var payload readyPayload
	if err := getJSON(ctx, target.URL+"/ready", nil, &payload); err != nil {
		view.Error = "unreachable: " + err.Error()
		return view
	}

	view.Status = payload.Status
	view.Version = payload.Version
	view.UptimeSeconds = payload.UptimeSeconds
	for _, d := range payload.Dependencies {
		view.Dependencies = append(view.Dependencies, DependencyView{
			Name:      d.Name,
			Status:    d.Status,
			LatencyMs: d.LatencyMs,
			Error:     d.Error,
			Optional:  d.Optional,
		})
	}
	if payload.Stream != nil {
		view.Stream = &StreamView{
			Status:              payload.Stream.Status,
			OperationsProcessed: payload.Stream.OperationsProcessed,
			SecondsSinceLastOp:  payload.Stream.SecondsSinceLastOp,
			LastOperationAt:     payload.Stream.LastOperationAt,
			Note:                payload.Stream.Note,
		}
	}
	return view
}

// hostHealth reads CPU, memory and disk from node-exporter via Prometheus.
func (c Config) hostHealth(ctx context.Context) HostHealth {
	if !c.HasPrometheus() {
		return HostHealth{}
	}
	host := HostHealth{Available: true}
	host.CPUPercent = c.scalar(ctx, `100 * (1 - avg(rate(node_cpu_seconds_total{mode="idle"}[5m])))`)
	host.MemoryPercent = c.scalar(ctx, `100 * (1 - node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes)`)
	host.DiskPercent = c.scalar(ctx, `100 * (1 - node_filesystem_avail_bytes{mountpoint="/"} / node_filesystem_size_bytes{mountpoint="/"})`)
	return host
}

// scalar runs an instant query expected to return a single number, returning
// nil when there is no usable sample so the UI can show "-" rather than zero.
func (c Config) scalar(ctx context.Context, promQL string) *float64 {
	res, err := c.queryInstant(ctx, promQL)
	if err != nil || len(res.Data.Result) == 0 {
		return nil
	}
	v, ok := sampleValue(res.Data.Result[0].Value)
	if !ok {
		return nil
	}
	return &v
}

// ---- Issues -----------------------------------------------------------------

// Issues lists grouped errors from the crash reporter across every app.
func (c Config) Issues(ctx context.Context, limit int, project, query string) IssuesResponse {
	if !c.HasGlitchTip() {
		return IssuesResponse{
			Issues:    []Issue{},
			Available: false,
			Reason:    "crash reporting is not configured for this environment",
		}
	}

	raw, err := c.fetchIssues(ctx, limit, project, query)
	if err != nil {
		return IssuesResponse{
			Issues:    []Issue{},
			Available: false,
			Reason:    "could not reach the crash reporter: " + err.Error(),
		}
	}

	issues := make([]Issue, 0, len(raw))
	for _, r := range raw {
		app := r.Project.Slug
		if app == "" {
			app = r.Project.Name
		}
		issues = append(issues, Issue{
			ID:        r.ID,
			Title:     r.Title,
			Culprit:   r.Culprit,
			Level:     r.Level,
			Status:    r.Status,
			Count:     issueCount(r.Count),
			App:       app,
			FirstSeen: r.FirstSeen,
			LastSeen:  r.LastSeen,
			Permalink: r.Permalink,
		})
	}
	return IssuesResponse{Issues: issues, Available: true, GlitchTipURL: c.PublicGlitchTipURL()}
}

// ---- Trends -----------------------------------------------------------------

// windows maps the UI's window options to a duration and a sensible step. The
// step keeps each series to roughly 60-120 points: enough to see shape,
// few enough to stay responsive.
var windows = map[string]struct {
	Duration time.Duration
	Step     time.Duration
}{
	"1h":  {time.Hour, time.Minute},
	"6h":  {6 * time.Hour, 5 * time.Minute},
	"24h": {24 * time.Hour, 15 * time.Minute},
}

// Trends returns request rate, error rate and p95 latency per service.
func (c Config) Trends(ctx context.Context, window string) TrendsResponse {
	w, ok := windows[window]
	if !ok {
		window, w = "6h", windows["6h"]
	}
	if !c.HasPrometheus() {
		return TrendsResponse{
			Window:      window,
			Available:   false,
			Reason:      "metrics are not configured for this environment",
			RequestRate: []Series{}, ErrorRate: []Series{}, LatencyP95: []Series{},
			HandledFailures: []Series{},
		}
	}

	out := TrendsResponse{Window: window, Available: true}
	var wg sync.WaitGroup
	wg.Add(4)

	go func() {
		defer wg.Done()
		out.RequestRate = c.series(ctx, `sum by (service) (rate(http_requests_total[5m]))`, w.Duration, w.Step)
	}()
	go func() {
		defer wg.Done()
		// Percentage of requests failing, rather than a raw count - a rate of
		// 3/s means nothing without knowing the total.
		out.ErrorRate = c.series(ctx,
			`100 * sum by (service) (rate(http_requests_total{status=~"5.."}[5m])) / clamp_min(sum by (service) (rate(http_requests_total[5m])), 0.001)`,
			w.Duration, w.Step)
	}()
	go func() {
		defer wg.Done()
		// A raw per-minute rate rather than a percentage: these failures have
		// no request to be a percentage of. A background worker failing every
		// poll produces a steady non-zero line here while ErrorRate above
		// stays at a flat, accurate zero.
		out.HandledFailures = c.series(ctx,
			`60 * sum by (service) (rate(handled_failures_total[5m]))`,
			w.Duration, w.Step)
	}()
	go func() {
		defer wg.Done()
		out.LatencyP95 = c.series(ctx,
			`histogram_quantile(0.95, sum by (service, le) (rate(http_request_duration_seconds_bucket[5m])))`,
			w.Duration, w.Step)
	}()
	wg.Wait()

	return out
}

// series runs a range query and converts it into named lines for the chart.
func (c Config) series(ctx context.Context, promQL string, window, step time.Duration) []Series {
	res, err := c.queryRange(ctx, promQL, window, step)
	if err != nil {
		return []Series{}
	}

	out := make([]Series, 0, len(res.Data.Result))
	for _, r := range res.Data.Result {
		name := r.Metric["service"]
		if name == "" {
			name = "unknown"
		}
		points := make([]SeriesPoint, 0, len(r.Values))
		for _, pair := range r.Values {
			ts, okT := sampleTime(pair)
			v, okV := sampleValue(pair)
			if !okT || !okV {
				continue
			}
			points = append(points, SeriesPoint{Time: ts.Format(time.RFC3339), Value: v})
		}
		out = append(out, Series{Name: name, Points: points})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

// ---- Trace ------------------------------------------------------------------

// requestIDPattern is deliberately strict. The request id is the only
// user-supplied value that reaches Loki, and LogQL treats quotes and braces as
// syntax - so anything outside this set is rejected rather than escaped.
var requestIDPattern = regexp.MustCompile(`^[A-Za-z0-9._-]{1,64}$`)

// ValidRequestID reports whether an id is safe to embed in a LogQL query.
func ValidRequestID(id string) bool { return requestIDPattern.MatchString(id) }

// Trace returns every log line carrying the given request id, across all
// services, in time order - the whole point of the request-id plumbing.
func (c Config) Trace(ctx context.Context, requestID string) TraceResponse {
	out := TraceResponse{RequestID: requestID, Lines: []LogLine{}}
	if !c.HasLoki() {
		out.Reason = "log search is not configured for this environment"
		return out
	}

	// The id is validated by the handler before reaching here; quoting it into
	// a line filter is safe for the character set ValidRequestID allows.
	logQL := fmt.Sprintf(`{container=~".+"} |= %q`, requestID)
	res, err := c.queryLoki(ctx, logQL, 24*time.Hour, 500)
	if err != nil {
		out.Reason = "could not reach log search: " + err.Error()
		return out
	}

	out.Lines = flattenLoki(res)
	out.Available = true
	return out
}

// RecentErrors lists the most recent failing requests, so the trace screen has
// something to offer an admin who does not already have a request id.
func (c Config) RecentErrors(ctx context.Context, limit int) RecentErrorsResponse {
	out := RecentErrorsResponse{Lines: []LogLine{}}
	if !c.HasLoki() {
		out.Reason = "log search is not configured for this environment"
		return out
	}

	logQL := `{container=~".+"} | json | status >= 400`
	res, err := c.queryLoki(ctx, logQL, time.Hour, limit)
	if err != nil {
		out.Reason = "could not reach log search: " + err.Error()
		return out
	}

	lines := flattenLoki(res)
	// Newest first: an admin looking for "what just broke" wants the latest.
	sort.Slice(lines, func(i, j int) bool { return lines[i].Time > lines[j].Time })
	if len(lines) > limit {
		lines = lines[:limit]
	}
	out.Lines = lines
	out.Available = true
	return out
}

// structuredLine is the JSON the observe package writes for each request.
type structuredLine struct {
	Level      string `json:"level"`
	Service    string `json:"service"`
	RequestID  string `json:"request_id"`
	Method     string `json:"method"`
	Path       string `json:"path"`
	Status     int    `json:"status"`
	DurationMs int64  `json:"duration_ms"`
	Msg        string `json:"msg"`
}

// flattenLoki turns Loki's stream/value structure into a flat, time-ordered
// list, parsing the structured fields where a line is our own JSON and falling
// back to the raw text where it is not (Traefik, Postgres and friends).
func flattenLoki(res lokiResult) []LogLine {
	var lines []LogLine
	for _, stream := range res.Data.Result {
		service := stream.Stream["container"]
		for _, pair := range stream.Values {
			if len(pair) < 2 {
				continue
			}
			ts := lokiTimestamp(pair[0])
			raw := pair[1]

			line := LogLine{Time: ts, Service: service, Message: raw}
			var parsed structuredLine
			if err := json.Unmarshal([]byte(raw), &parsed); err == nil && parsed.Service != "" {
				line.Level = parsed.Level
				line.Method = parsed.Method
				line.Path = parsed.Path
				line.Status = parsed.Status
				line.DurationMs = parsed.DurationMs
				if parsed.Service != "" {
					line.Service = parsed.Service
				}
				if parsed.Msg != "" {
					line.Message = parsed.Msg
				}
			}
			lines = append(lines, line)
		}
	}
	sort.Slice(lines, func(i, j int) bool { return lines[i].Time < lines[j].Time })
	if lines == nil {
		lines = []LogLine{}
	}
	return lines
}

// lokiTimestamp converts Loki's nanosecond string timestamp to RFC3339.
func lokiTimestamp(ns string) string {
	n, err := parseInt64(ns)
	if err != nil {
		return ""
	}
	return time.Unix(0, n).UTC().Format(time.RFC3339Nano)
}

func parseInt64(s string) (int64, error) {
	var n int64
	_, err := fmt.Sscanf(strings.TrimSpace(s), "%d", &n)
	return n, err
}
