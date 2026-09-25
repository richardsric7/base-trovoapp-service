package observability

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// fakeReady stands in for a backend's /ready endpoint.
func fakeReady(t *testing.T, payload string, status int) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/ready" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		w.Write([]byte(payload))
	}))
}

func TestHealthReportsDependenciesPerService(t *testing.T) {
	svc := fakeReady(t, `{"status":"up","service":"admin-api","version":"abc123","uptime_seconds":3600,
		"dependencies":[{"name":"admin_db","status":"up","latency_ms":12},
		                {"name":"redis","status":"down","latency_ms":0,"error":"refused","optional":true}]}`, 200)
	defer svc.Close()

	cfg := Config{Services: []ServiceTarget{{Name: "admin-api", URL: svc.URL}}}
	out := cfg.Health(context.Background())

	if len(out.Services) != 1 {
		t.Fatalf("got %d services, want 1", len(out.Services))
	}
	s := out.Services[0]
	if s.Status != "up" || s.Version != "abc123" || s.UptimeSeconds != 3600 {
		t.Errorf("service summary wrong: %+v", s)
	}
	if len(s.Dependencies) != 2 {
		t.Fatalf("got %d dependencies, want 2", len(s.Dependencies))
	}
	if s.Dependencies[1].Error != "refused" || !s.Dependencies[1].Optional {
		t.Errorf("a failing optional dependency must keep its reason: %+v", s.Dependencies[1])
	}
}

func TestHealthReportsUnreachableServiceRatherThanOmittingIt(t *testing.T) {
	// A service that cannot be reached must still appear, as down. Omitting it
	// would look like the service does not exist.
	srv := fakeReady(t, `{}`, 200)
	url := srv.URL
	srv.Close()

	cfg := Config{Services: []ServiceTarget{{Name: "wallet-api", URL: url}}}
	out := cfg.Health(context.Background())

	if len(out.Services) != 1 {
		t.Fatalf("got %d services, want 1", len(out.Services))
	}
	if out.Services[0].Status != "down" {
		t.Errorf("status = %q, want down", out.Services[0].Status)
	}
	if out.Services[0].Error == "" {
		t.Error("an unreachable service must say why")
	}
}

func TestHealthQueriesServicesConcurrently(t *testing.T) {
	// Three slow services must not take three times as long: an admin opening
	// a health page during an incident should not wait on serial polling.
	slow := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.Write([]byte(`{"status":"up","service":"x"}`))
	}))
	defer slow.Close()

	cfg := Config{Services: []ServiceTarget{
		{Name: "a", URL: slow.URL}, {Name: "b", URL: slow.URL}, {Name: "c", URL: slow.URL},
	}}

	start := time.Now()
	cfg.Health(context.Background())
	elapsed := time.Since(start)

	if elapsed > 500*time.Millisecond {
		t.Errorf("took %v for 3 services at 200ms each - looks serial", elapsed)
	}
}

func TestIssuesUnavailableWhenNotConfigured(t *testing.T) {
	// A deployment without crash reporting must say so, not return an empty
	// list that reads as "no errors".
	out := Config{}.Issues(context.Background(), 10, "", "")
	if out.Available {
		t.Error("issues should be unavailable with no crash reporter configured")
	}
	if out.Reason == "" {
		t.Error("unavailable screens must explain why")
	}
	if out.Issues == nil {
		t.Error("issues must be an empty array, not null, so the UI can map over it")
	}
}

func TestIssuesNormalisesUpstreamPayload(t *testing.T) {
	gt := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.Header.Get("Authorization"), "Bearer ") {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.Write([]byte(`[{"id":"7","title":"TypeError: x is undefined","culprit":"NavBar",
			"level":"error","status":"unresolved","count":"14",
			"firstSeen":"2026-09-01T10:00:00Z","lastSeen":"2026-09-09T12:00:00Z",
			"project":{"name":"Admin Web","slug":"admin-web"}}]`))
	}))
	defer gt.Close()

	cfg := Config{GlitchTipURL: gt.URL, GlitchTipToken: "tok", GlitchTipOrg: "trovo"}
	out := cfg.Issues(context.Background(), 10, "", "")

	if !out.Available || len(out.Issues) != 1 {
		t.Fatalf("expected one issue, got %+v", out)
	}
	got := out.Issues[0]
	// count arrives as a string from this version; it must still be a number.
	if got.Count != 14 {
		t.Errorf("count = %d, want 14 (string counts must be parsed)", got.Count)
	}
	if got.App != "admin-web" {
		t.Errorf("app = %q, want admin-web", got.App)
	}
}

func TestIssuesSurvivesUpstreamFailure(t *testing.T) {
	gt := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer gt.Close()

	cfg := Config{GlitchTipURL: gt.URL, GlitchTipToken: "tok", GlitchTipOrg: "trovo"}
	out := cfg.Issues(context.Background(), 10, "", "")

	// The crash reporter being down must not 500 the admin API - that is
	// precisely when someone is looking at this screen.
	if out.Available {
		t.Error("should report unavailable when the upstream fails")
	}
	if out.Reason == "" {
		t.Error("must explain the failure")
	}
}

func TestValidRequestIDRejectsInjection(t *testing.T) {
	// The request id is the only user-supplied value that reaches LogQL.
	valid := []string{"f3bb087a-d3a0-41f0-871d-c9e0ce2f5142", "abc123", "a.b-c_d"}
	for _, id := range valid {
		if !ValidRequestID(id) {
			t.Errorf("%q should be accepted", id)
		}
	}

	invalid := []string{
		`" or true or "`,        // quote breaks out of the line filter
		`abc"}|="x`,             // brace and pipe are LogQL syntax
		"abc def",               // whitespace
		"abc\ndef",              // newline
		"",                      // empty
		strings.Repeat("a", 65), // over length
	}
	for _, id := range invalid {
		if ValidRequestID(id) {
			t.Errorf("%q must be rejected - it could alter the log query", id)
		}
	}
}

func TestTraceParsesStructuredLines(t *testing.T) {
	loki := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Confirm the request id actually reached the query.
		if !strings.Contains(r.URL.RawQuery, "req-123") {
			t.Errorf("request id missing from LogQL: %s", r.URL.RawQuery)
		}
		w.Write([]byte(`{"data":{"result":[{"stream":{"container":"admin-api-dev"},
			"values":[["1757000000000000000","{\"level\":\"error\",\"service\":\"admin-api\",\"request_id\":\"req-123\",\"method\":\"GET\",\"path\":\"/x\",\"status\":500,\"duration_ms\":42,\"msg\":\"request\"}"]]}]}}`))
	}))
	defer loki.Close()

	cfg := Config{LokiURL: loki.URL}
	out := cfg.Trace(context.Background(), "req-123")

	if !out.Available || len(out.Lines) != 1 {
		t.Fatalf("expected one line, got %+v", out)
	}
	line := out.Lines[0]
	if line.Status != 500 || line.Method != "GET" || line.Path != "/x" || line.DurationMs != 42 {
		t.Errorf("structured fields not parsed: %+v", line)
	}
	if line.Service != "admin-api" {
		t.Errorf("service = %q, want admin-api from the log body", line.Service)
	}
}

func TestTraceKeepsUnstructuredLines(t *testing.T) {
	// Traefik and Postgres do not write our JSON; those lines must survive as
	// raw text rather than being dropped.
	loki := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"data":{"result":[{"stream":{"container":"traefik"},
			"values":[["1757000000000000000","plain text log line"]]}]}}`))
	}))
	defer loki.Close()

	out := Config{LokiURL: loki.URL}.Trace(context.Background(), "req-123")
	if len(out.Lines) != 1 || out.Lines[0].Message != "plain text log line" {
		t.Fatalf("unstructured line lost: %+v", out.Lines)
	}
	if out.Lines[0].Service != "traefik" {
		t.Errorf("service should fall back to the container label, got %q", out.Lines[0].Service)
	}
}

func TestTrendsReturnsSeriesPerService(t *testing.T) {
	prom := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"data":{"result":[
			{"metric":{"service":"wallet-api"},"values":[[1757000000,"1.5"],[1757000060,"2.5"]]},
			{"metric":{"service":"admin-api"},"values":[[1757000000,"0.5"]]}]}}`))
	}))
	defer prom.Close()

	out := Config{PrometheusURL: prom.URL}.Trends(context.Background(), "1h")

	if !out.Available {
		t.Fatal("trends should be available with Prometheus configured")
	}
	if len(out.RequestRate) != 2 {
		t.Fatalf("got %d series, want 2", len(out.RequestRate))
	}
	// Sorted so the chart legend is stable between refreshes.
	if out.RequestRate[0].Name != "admin-api" {
		t.Errorf("series should be sorted by name, got %q first", out.RequestRate[0].Name)
	}
	if len(out.RequestRate[1].Points) != 2 || out.RequestRate[1].Points[1].Value != 2.5 {
		t.Errorf("points not parsed: %+v", out.RequestRate[1].Points)
	}
}

func TestTrendsRejectsUnknownWindow(t *testing.T) {
	// An unrecognised window must fall back, not produce an unbounded query.
	out := Config{}.Trends(context.Background(), "99y")
	if out.Window != "6h" {
		t.Errorf("window = %q, want the 6h default", out.Window)
	}
}

func TestScalarReturnsNilRatherThanZeroWhenAbsent(t *testing.T) {
	// Showing 0% disk usage when Prometheus has no data would be a lie.
	prom := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"data":{"result":[]}}`))
	}))
	defer prom.Close()

	host := Config{PrometheusURL: prom.URL}.hostHealth(context.Background())
	if host.DiskPercent != nil {
		t.Errorf("absent metric should be nil, got %v", *host.DiskPercent)
	}
}

func TestLoadConfigParsesServiceList(t *testing.T) {
	t.Setenv("OBSERVABILITY_SERVICES", "admin-api=http://admin:8080, p2p-api=http://p2p-dev:8080/ ,broken")
	cfg := LoadConfig()

	if len(cfg.Services) != 2 {
		t.Fatalf("got %d services, want 2 (malformed entries skipped)", len(cfg.Services))
	}
	if cfg.Services[1].Name != "p2p-api" || cfg.Services[1].URL != "http://p2p-dev:8080" {
		t.Errorf("trailing slash and spaces should be trimmed: %+v", cfg.Services[1])
	}
}

func TestCapabilitiesReflectConfiguration(t *testing.T) {
	cfg := Config{PrometheusURL: "http://p", Services: []ServiceTarget{{Name: "a", URL: "http://a"}}}
	if !cfg.HasPrometheus() || cfg.HasLoki() || cfg.HasGlitchTip() {
		t.Errorf("capability detection wrong: prom=%v loki=%v gt=%v",
			cfg.HasPrometheus(), cfg.HasLoki(), cfg.HasGlitchTip())
	}

	// A token without an org slug cannot address any endpoint, so it is not
	// "configured".
	partial := Config{GlitchTipURL: "http://g", GlitchTipToken: "t"}
	if partial.HasGlitchTip() {
		t.Error("crash reporting needs url, token and org - partial config must not count")
	}
}

func TestResponsesAreJSONSerialisable(t *testing.T) {
	// Empty slices must serialise as [] not null, or the UI cannot map over them.
	for name, v := range map[string]any{
		"issues":        Config{}.Issues(context.Background(), 10, "", ""),
		"trends":        Config{}.Trends(context.Background(), "1h"),
		"trace":         Config{}.Trace(context.Background(), "abc"),
		"recent-errors": Config{}.RecentErrors(context.Background(), 10),
	} {
		b, err := json.Marshal(v)
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		if strings.Contains(string(b), ":null") {
			t.Errorf("%s serialises a null collection: %s", name, b)
		}
	}
}

func TestPublicGlitchTipURLIsWhatReachesTheBrowser(t *testing.T) {
	// The address this service calls and the address a user's browser can
	// reach are different: one is a container name on the internal network.
	// Handing the internal one to the UI produces links that cannot resolve.
	cfg := Config{
		GlitchTipURL:       "http://glitchtip-web:8000",
		GlitchTipPublicURL: "https://errors.dev.trovo.app",
		GlitchTipToken:     "t",
		GlitchTipOrg:       "o",
	}
	if got := cfg.PublicGlitchTipURL(); got != "https://errors.dev.trovo.app" {
		t.Errorf("link URL = %q, want the public hostname", got)
	}

	// With no public URL configured, fall back to the internal one rather than
	// omitting the link - an obviously broken link is easier to notice and fix
	// than a silently missing feature.
	fallback := Config{GlitchTipURL: "http://glitchtip-web:8000"}
	if got := fallback.PublicGlitchTipURL(); got != "http://glitchtip-web:8000" {
		t.Errorf("fallback = %q, want the internal URL", got)
	}
}

func TestIssuesLinkUsesPublicURL(t *testing.T) {
	gt := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`[]`))
	}))
	defer gt.Close()

	cfg := Config{
		GlitchTipURL:       gt.URL,
		GlitchTipPublicURL: "https://errors.dev.trovo.app",
		GlitchTipToken:     "t",
		GlitchTipOrg:       "o",
	}
	out := cfg.Issues(context.Background(), 10, "", "")
	if out.GlitchTipURL != "https://errors.dev.trovo.app" {
		t.Errorf("issues link = %q, want the public hostname", out.GlitchTipURL)
	}
}

func TestWorkerStreamStatusIsPassedThrough(t *testing.T) {
	// The payment history engine reports whether its blockchain stream is still
	// advancing. That is the only meaningful health signal for a worker - "the
	// process is alive" says nothing about whether payment history is still
	// being written - so it must survive the trip through this gateway.
	svc := fakeReady(t, `{"status":"up","service":"payment-history-engine","version":"321fa58",
		"uptime_seconds":3600,
		"stream":{"status":"stalled","operations_processed":14238,
		          "seconds_since_last_operation":1820,
		          "last_operation_at":"2026-09-10T12:00:00Z",
		          "note":"no operations processed recently - the stream may have stopped"},
		"dependencies":[{"name":"wallet_db","status":"up","latency_ms":10}]}`, 200)
	defer svc.Close()

	cfg := Config{Services: []ServiceTarget{{Name: "payment-history", URL: svc.URL}}}
	out := cfg.Health(context.Background())

	if len(out.Services) != 1 {
		t.Fatalf("got %d services, want 1", len(out.Services))
	}
	stream := out.Services[0].Stream
	if stream == nil {
		t.Fatal("stream status was dropped - the one signal that matters for a worker")
	}
	if stream.Status != "stalled" {
		t.Errorf("status = %q, want stalled", stream.Status)
	}
	if stream.OperationsProcessed != 14238 {
		t.Errorf("operations = %d, want 14238", stream.OperationsProcessed)
	}
	if stream.SecondsSinceLastOp != 1820 {
		t.Errorf("seconds since last op = %d, want 1820", stream.SecondsSinceLastOp)
	}
	if stream.Note == "" {
		t.Error("a stalled stream must carry its explanation through")
	}
}

func TestRequestDrivenServicesHaveNoStream(t *testing.T) {
	// An API has no stream to report, and the field must be absent rather than
	// present-and-empty, so the interface can omit the section entirely.
	svc := fakeReady(t, `{"status":"up","service":"admin-api","version":"abc",
		"dependencies":[{"name":"admin_db","status":"up","latency_ms":12}]}`, 200)
	defer svc.Close()

	out := Config{Services: []ServiceTarget{{Name: "admin-api", URL: svc.URL}}}.Health(context.Background())
	if out.Services[0].Stream != nil {
		t.Errorf("an API should report no stream, got %+v", out.Services[0].Stream)
	}

	// And it must not appear in the JSON at all.
	b, _ := json.Marshal(out.Services[0])
	if strings.Contains(string(b), "stream") {
		t.Errorf("stream key should be omitted for APIs: %s", b)
	}
}
