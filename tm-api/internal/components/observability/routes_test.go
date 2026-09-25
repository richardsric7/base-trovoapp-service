package observability

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// mountForTest registers the handlers without the auth guard, so the routing
// and payload shape can be exercised directly.
func mountForTest(cfg Config) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	g := r.Group("/api/v1/observability")
	g.GET("/capabilities", CapabilitiesHandler(cfg))
	g.GET("/health", HealthHandler(cfg))
	g.GET("/issues", IssuesHandler(cfg))
	g.GET("/trends", TrendsHandler(cfg))
	g.GET("/trace/:request_id", TraceHandler(cfg))
	g.GET("/recent-errors", RecentErrorsHandler(cfg))
	return r
}

func TestEveryEndpointAnswersWithNothingConfigured(t *testing.T) {
	// A deployment with no monitoring stack must still serve the screens, so
	// the UI can render "not configured" rather than erroring.
	r := mountForTest(Config{})

	for _, path := range []string{
		"/api/v1/observability/capabilities",
		"/api/v1/observability/health",
		"/api/v1/observability/issues",
		"/api/v1/observability/trends",
		"/api/v1/observability/recent-errors",
		"/api/v1/observability/trace/abc123",
	} {
		t.Run(path, func(t *testing.T) {
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))

			if rec.Code != http.StatusOK {
				t.Fatalf("%s = %d, want 200 (body: %s)", path, rec.Code, rec.Body.String())
			}
			var body struct {
				Data json.RawMessage `json:"data"`
			}
			if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
				t.Fatalf("%s: invalid JSON: %v", path, err)
			}
			if len(body.Data) == 0 {
				t.Errorf("%s: response missing a data envelope", path)
			}
		})
	}
}

func TestTraceRejectsMalformedRequestID(t *testing.T) {
	r := mountForTest(Config{LokiURL: "http://unused"})

	rec := httptest.NewRecorder()
	// A quote would break out of the LogQL line filter if it were interpolated.
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, `/api/v1/observability/trace/abc%22or%22x`, nil))

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400 for a malformed request id", rec.Code)
	}
}

func TestIssuesLimitIsBounded(t *testing.T) {
	// A caller must not be able to ask the crash reporter for an unbounded page.
	var seen string
	gt := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		seen = req.URL.RawQuery
		w.Write([]byte(`[]`))
	}))
	defer gt.Close()

	r := mountForTest(Config{GlitchTipURL: gt.URL, GlitchTipToken: "t", GlitchTipOrg: "o"})
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/v1/observability/issues?limit=99999", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
	if !contains(seen, "limit=100") {
		t.Errorf("limit not clamped, upstream query was %q", seen)
	}
}

func contains(haystack, needle string) bool {
	return len(haystack) >= len(needle) && (haystack == needle || len(needle) == 0 ||
		func() bool {
			for i := 0; i+len(needle) <= len(haystack); i++ {
				if haystack[i:i+len(needle)] == needle {
					return true
				}
			}
			return false
		}())
}
