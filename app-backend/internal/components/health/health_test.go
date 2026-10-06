package health

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"trovo-wallet-api/internal/sharedconfig"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func serve(t *testing.T, gc *sharedconfig.GlobalConfig, path string) (int, map[string]interface{}) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	r := gin.New()
	Init(r, gc)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, path, nil))
	var body map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &body)
	return w.Code, body
}

func TestHealthAlwaysOK(t *testing.T) {
	code, body := serve(t, &sharedconfig.GlobalConfig{}, "/health")
	if code != http.StatusOK || body["status"] != "ok" {
		t.Fatalf("got %d %v", code, body)
	}
}

func TestReadyReportsEachDependency(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	// databases up, no RPC client: not ready, and it says which part is down
	code, body := serve(t, &sharedconfig.GlobalConfig{DB: db, RoachDB: db}, "/ready")
	if code != http.StatusServiceUnavailable {
		t.Fatalf("want 503 without a blockchain client, got %d %v", code, body)
	}
	checks, _ := body["checks"].(map[string]interface{})
	if checks["database"] != "up" || checks["tracking_database"] != "up" || checks["blockchain"] != "down" {
		t.Fatalf("unexpected checks %v", checks)
	}
	if _, ok := checks["redis"]; ok {
		t.Fatalf("redis is not checked when caching is off: %v", checks)
	}
}
