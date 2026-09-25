package observe

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	healthservices "admin-panel-dashboard/internal/components/health/services"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
)

func newRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(RequestID(), Logger(), Recovery(), Metrics())
	return r
}

func TestRequestIDGeneratedAndReturned(t *testing.T) {
	r := newRouter()
	var seen string
	r.GET("/x", func(c *gin.Context) {
		seen = RequestIDFrom(c)
		c.Status(http.StatusOK)
	})

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/x", nil))

	header := rec.Header().Get(RequestIDHeader)
	if header == "" {
		t.Fatal("response must carry the request ID header")
	}
	if seen != header {
		t.Errorf("handler saw %q but response header was %q - they must match", seen, header)
	}
}

func TestRequestIDHonoursUpstreamValue(t *testing.T) {
	// An ID assigned upstream must survive, or a request cannot be traced
	// across service hops.
	r := newRouter()
	r.GET("/x", func(c *gin.Context) { c.Status(http.StatusOK) })

	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	req.Header.Set(RequestIDHeader, "upstream-abc-123")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if got := rec.Header().Get(RequestIDHeader); got != "upstream-abc-123" {
		t.Errorf("request ID = %q, want the upstream value preserved", got)
	}
}

func TestRequestIDRejectsMalformedUpstreamValue(t *testing.T) {
	// A header with newlines would forge extra lines in the logs; an
	// over-long one would bloat every line. Both must be replaced.
	for name, bad := range map[string]string{
		"newline":  "abc\ndef",
		"too long": strings.Repeat("x", 200),
		"spaces":   "abc def",
	} {
		t.Run(name, func(t *testing.T) {
			r := newRouter()
			r.GET("/x", func(c *gin.Context) { c.Status(http.StatusOK) })

			req := httptest.NewRequest(http.MethodGet, "/x", nil)
			req.Header.Set(RequestIDHeader, bad)
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)

			if got := rec.Header().Get(RequestIDHeader); got == bad {
				t.Errorf("malformed upstream ID %q was accepted", bad)
			}
		})
	}
}

func TestRecoveryConvertsPanicTo500WithRequestID(t *testing.T) {
	r := newRouter()
	r.GET("/boom", func(c *gin.Context) { panic("kaboom") })

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/boom", nil))

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500", rec.Code)
	}

	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("panic response is not JSON: %v", err)
	}
	// The ID is what lets a user's bug report be found in the logs.
	if body["request_id"] == "" {
		t.Error("panic response must include the request ID")
	}
	if body["error"] == "" {
		t.Error("panic response must include an error message")
	}
}

func TestRecoveryDoesNotLeakPanicDetail(t *testing.T) {
	// The panic value can contain internal state; it belongs in the logs and
	// the crash report, not in the HTTP response.
	r := newRouter()
	r.GET("/boom", func(c *gin.Context) { panic("connection string postgres://user:hunter2@db") })

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/boom", nil))

	if strings.Contains(rec.Body.String(), "hunter2") {
		t.Fatalf("panic detail leaked into the response body: %s", rec.Body.String())
	}
}

func TestServerKeepsServingAfterPanic(t *testing.T) {
	r := newRouter()
	r.GET("/boom", func(c *gin.Context) { panic("kaboom") })
	r.GET("/ok", func(c *gin.Context) { c.String(http.StatusOK, "fine") })

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/boom", nil))

	rec2 := httptest.NewRecorder()
	r.ServeHTTP(rec2, httptest.NewRequest(http.MethodGet, "/ok", nil))
	if rec2.Code != http.StatusOK {
		t.Fatalf("a panic on one route broke a later request: got %d", rec2.Code)
	}
}

func TestMetricsEndpointExposesRequestCounters(t *testing.T) {
	r := newRouter()
	r.GET("/metrics", MetricsHandler())
	r.GET("/counted", func(c *gin.Context) { c.Status(http.StatusOK) })

	r.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/counted", nil))

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/metrics", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("/metrics status = %d, want 200", rec.Code)
	}
	body := rec.Body.String()
	for _, want := range []string{"http_requests_total", "http_request_duration_seconds", `route="/counted"`} {
		if !strings.Contains(body, want) {
			t.Errorf("/metrics output missing %q", want)
		}
	}
}

func TestMetricsGroupsUnmatchedRoutes(t *testing.T) {
	// A 404 scan across many random URLs must not create a time series per
	// URL, which would blow up Prometheus cardinality.
	r := newRouter()
	r.GET("/metrics", MetricsHandler())
	for _, p := range []string{"/nope-1", "/nope-2", "/nope-3"} {
		r.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, p, nil))
	}

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/metrics", nil))
	body := rec.Body.String()

	if !strings.Contains(body, `route="unmatched"`) {
		t.Error(`unmatched requests should be grouped under route="unmatched"`)
	}
	for _, p := range []string{"/nope-1", "/nope-2", "/nope-3"} {
		if strings.Contains(body, `route="`+p+`"`) {
			t.Errorf("unmatched path %q became its own metric label", p)
		}
	}
}

func TestScrubEventRemovesCredentials(t *testing.T) {
	event := &sentry.Event{
		Request: &sentry.Request{
			Headers: map[string]string{
				"Authorization": "Bearer secret-token",
				"Cookie":        "session=abc",
				"User-Agent":    "test-agent",
			},
			QueryString: "token=leaked",
			Cookies:     "session=abc",
			Data:        `{"password":"hunter2"}`,
		},
	}

	got := scrubEvent(event, nil)

	for _, h := range []string{"Authorization", "Cookie"} {
		if _, present := got.Request.Headers[h]; present {
			t.Errorf("%s header must be stripped from crash reports", h)
		}
	}
	if got.Request.Headers["User-Agent"] != "test-agent" {
		t.Error("harmless headers should be kept - they aid diagnosis")
	}
	if got.Request.QueryString != "" || got.Request.Cookies != "" || got.Request.Data != "" {
		t.Error("query string, cookies and body must all be cleared")
	}
}

func TestInitWithoutDSNDisablesReportingButDoesNotFail(t *testing.T) {
	t.Setenv("SENTRY_DSN", "")
	initialised = false
	flush := Init("test")
	defer flush()

	if Enabled() {
		t.Error("crash reporting should be disabled when no DSN is configured")
	}
	// And a panic must still be recovered and logged locally.
	r := newRouter()
	r.GET("/boom", func(c *gin.Context) { panic("kaboom") })
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/boom", nil))
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500 even with reporting disabled", rec.Code)
	}
}

func TestLoggerEmitsStructuredJSONWithRequestID(t *testing.T) {
	// The log line is the product here, so assert on its parsed fields.
	r := newRouter()
	r.GET("/logged", func(c *gin.Context) { c.Status(http.StatusTeapot) })

	out := captureStdout(t, func() {
		r.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/logged", nil))
	})

	var line map[string]any
	if err := json.Unmarshal([]byte(strings.TrimSpace(out)), &line); err != nil {
		t.Fatalf("log line is not valid JSON (%v): %s", err, out)
	}
	if line["request_id"] == "" || line["request_id"] == nil {
		t.Error("log line must carry the request ID")
	}
	if line["path"] != "/logged" {
		t.Errorf("path = %v, want /logged", line["path"])
	}
	if line["status"] != float64(http.StatusTeapot) {
		t.Errorf("status = %v, want 418", line["status"])
	}
	if line["level"] != "warn" {
		t.Errorf("level = %v, want warn for a 4xx", line["level"])
	}
	if line["service"] != "admin-api" {
		t.Errorf("service = %v, want admin-api", line["service"])
	}
}

func TestLoggerSkipsProbeEndpoints(t *testing.T) {
	// Probes run every few seconds; logging them would swamp retention.
	r := newRouter()
	r.GET("/health", func(c *gin.Context) { c.Status(http.StatusOK) })

	out := captureStdout(t, func() {
		r.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/health", nil))
	})

	if strings.TrimSpace(out) != "" {
		t.Errorf("/health should not be logged, got: %s", out)
	}
}

func TestLoggerMarksServerErrors(t *testing.T) {
	r := newRouter()
	r.GET("/fail", func(c *gin.Context) { c.Status(http.StatusInternalServerError) })

	out := captureStdout(t, func() {
		r.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/fail", nil))
	})

	var line map[string]any
	if err := json.Unmarshal([]byte(strings.TrimSpace(out)), &line); err != nil {
		t.Fatalf("not JSON: %s", out)
	}
	if line["level"] != "error" {
		t.Errorf("level = %v, want error for a 5xx", line["level"])
	}
}

func TestInitFallsBackToBuildVersion(t *testing.T) {
	// Regression: Init was passed os.Getenv("APP_VERSION"), which is empty on
	// deployed builds because the version is stamped in at link time. Crash
	// reports then showed a blank release, so an error could not be tied to
	// the build that produced it.
	t.Setenv("APP_VERSION", "stamped-build-123")
	t.Setenv("SENTRY_DSN", "")
	initialised = false

	flush := Init("")
	defer flush()

	if got := healthservices.Version(); got != "stamped-build-123" {
		t.Fatalf("version resolver returned %q, want the stamped build", got)
	}
}

// A handler that returns 500 without panicking must still be reported: those
// are the ones Recovery() cannot see, and they were reaching GlitchTip only
// from the browser.
func TestLoggerReportsNonPanicServerErrors(t *testing.T) {
	r := newRouter()
	r.GET("/broken", func(c *gin.Context) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "boom"})
	})

	out := captureStdout(t, func() {
		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest("GET", "/broken", nil))
		if w.Code != http.StatusInternalServerError {
			t.Fatalf("status = %d, want 500", w.Code)
		}
	})

	// Reporting is disabled in tests (no DSN), so the observable contract is
	// that the request is logged at error level and nothing panics on the way
	// through reportServerError.
	if !strings.Contains(out, `"level":"error"`) {
		t.Fatalf("expected an error-level log line, got: %s", out)
	}
}

// 4xx stay out of crash reporting on purpose - they are expected outcomes.
func TestLoggerDoesNotReportClientErrors(t *testing.T) {
	r := newRouter()
	r.GET("/bad", func(c *gin.Context) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "nope"})
	})

	out := captureStdout(t, func() {
		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest("GET", "/bad", nil))
	})

	if !strings.Contains(out, `"level":"warn"`) {
		t.Fatalf("expected a warn-level log line for 4xx, got: %s", out)
	}
}
