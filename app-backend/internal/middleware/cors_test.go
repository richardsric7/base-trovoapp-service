package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func corsOrigin(t *testing.T, allowed, origin string) string {
	t.Helper()
	t.Setenv("CORS_ALLOWED_ORIGINS", allowed)
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(CORSMiddleware())
	r.GET("/x", func(c *gin.Context) { c.Status(http.StatusOK) })
	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	req.Header.Set("Origin", origin)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Header().Get("Access-Control-Allow-Credentials") != "" {
		t.Fatal("credentials must never be allowed")
	}
	return w.Header().Get("Access-Control-Allow-Origin")
}

func TestCORS(t *testing.T) {
	if got := corsOrigin(t, "", "https://evil.example"); got != "*" {
		t.Errorf("unset: got %q, want *", got)
	}
	list := "https://wallet.trovo.example, https://staging.trovo.example/"
	if got := corsOrigin(t, list, "https://wallet.trovo.example"); got != "https://wallet.trovo.example" {
		t.Errorf("listed origin: got %q", got)
	}
	if got := corsOrigin(t, list, "https://staging.trovo.example"); got != "https://staging.trovo.example" {
		t.Errorf("trailing slash in setting: got %q", got)
	}
	if got := corsOrigin(t, list, "https://evil.example"); got != "" {
		t.Errorf("unlisted origin must get no allow header, got %q", got)
	}
}
