package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	serverModels "admin-panel-dashboard/internal/server/models"

	"github.com/gin-gonic/gin"
)

// TestRoutesRegisteredAndUnauthenticated boots the real router with the health
// component mounted and exercises every path a probe might use. It asserts the
// endpoints answer without credentials - a probe that needs a token cannot
// report on a service whose auth path is broken.
func TestRoutesRegisteredAndUnauthenticated(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// No databases configured: every dependency reports "skipped", which must
	// not make the service unready.
	Init(router, &serverModels.Server{})

	for _, path := range []string{"/health", "/ready", "/api/v1/health", "/api/v1/ready"} {
		t.Run(path, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, path, nil)
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			if rec.Code != http.StatusOK {
				t.Fatalf("%s = %d, want 200 (body: %s)", path, rec.Code, rec.Body.String())
			}

			var body map[string]any
			if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
				t.Fatalf("%s returned invalid JSON: %v", path, err)
			}
			if body["status"] != "up" {
				t.Errorf("%s status = %v, want up", path, body["status"])
			}
			if body["service"] != "admin-api" {
				t.Errorf("%s service = %v, want admin-api", path, body["service"])
			}
			if _, ok := body["version"]; !ok {
				t.Errorf("%s missing version", path)
			}
		})
	}
}

// TestReadinessListsDependencies confirms /ready reports each dependency by
// name, which is what makes a failure diagnosable at a glance.
func TestReadinessListsDependencies(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	Init(router, &serverModels.Server{})

	req := httptest.NewRequest(http.MethodGet, "/ready", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	var body struct {
		Dependencies []struct {
			Name   string `json:"name"`
			Status string `json:"status"`
		} `json:"dependencies"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	want := map[string]bool{"admin_db": false, "wallet_db": false, "p2p_db": false, "redis": false, "wallet_api": false}
	for _, dep := range body.Dependencies {
		if _, ok := want[dep.Name]; !ok {
			t.Errorf("unexpected dependency %q", dep.Name)
			continue
		}
		want[dep.Name] = true
	}
	for name, found := range want {
		if !found {
			t.Errorf("dependency %q missing from readiness report", name)
		}
	}
}
