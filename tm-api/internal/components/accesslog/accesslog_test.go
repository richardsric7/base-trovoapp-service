package accesslog

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"admin-panel-dashboard/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func testDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file:"+uuid.NewString()+"?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&models.AdminAccessLog{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return db
}

// recordSync is used directly so the test is deterministic (Record is async). It
// also passes an empty IP so no live geo lookup is attempted.
func TestRecordSync_WritesRow(t *testing.T) {
	db := testDB(t)
	id := uint(7)
	recordSync(db, Event{
		Event:    models.EventAdminSuspend,
		Status:   models.AccessStatusSuccessful,
		Category: models.AccessCategoryAccount,
		Actor:    Actor{AdminID: &id, Username: "jane", Email: "jane@x.com", Role: "SUPER_ADMIN"},
		Target:   "bob@x.com",
	})

	var got models.AdminAccessLog
	if err := db.First(&got).Error; err != nil {
		t.Fatalf("expected a row: %v", err)
	}
	if got.Event != models.EventAdminSuspend || got.Status != models.AccessStatusSuccessful {
		t.Errorf("event/status mismatch: %+v", got)
	}
	if got.Category != models.AccessCategoryAccount {
		t.Errorf("category = %q, want %q", got.Category, models.AccessCategoryAccount)
	}
	if got.Action != "Admin suspend" {
		t.Errorf("action = %q, want %q", got.Action, "Admin suspend")
	}
	if got.Target != "bob@x.com" || got.Username != "jane" {
		t.Errorf("target/username mismatch: %+v", got)
	}
	if got.ActorAdminID == nil || *got.ActorAdminID != 7 {
		t.Errorf("actor id mismatch: %+v", got.ActorAdminID)
	}
	if got.OccurredAt.IsZero() {
		t.Error("occurred_at not set")
	}
}

func TestHumanise(t *testing.T) {
	cases := map[string]string{
		models.EventLoginSuccess:      "Login",
		models.EventLoginFailure:      "Login failed",
		models.EventLogout:            "Logout",
		models.EventAdminStatusChange: "Admin status change",
		models.EventSecretWrite:       "Secret write",
	}
	for in, want := range cases {
		if got := humanise(in); got != want {
			t.Errorf("humanise(%q) = %q, want %q", in, got, want)
		}
	}
}

// TestAudit_FailedStatus verifies the middleware records status=failed when the
// handler responds 4xx, and status=successful on 2xx.
func TestAudit_StatusFromResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)

	for _, tc := range []struct {
		name       string
		handlerErr bool
		want       string
	}{
		{"success", false, models.AccessStatusSuccessful},
		{"failure", true, models.AccessStatusFailed},
	} {
		t.Run(tc.name, func(t *testing.T) {
			db := testDB(t)
			r := gin.New()
			r.POST("/x", Audit(db, models.EventAdminInvite, models.AccessCategoryAccount, nil), func(c *gin.Context) {
				if tc.handlerErr {
					c.JSON(http.StatusBadRequest, gin.H{"error": "nope"})
					return
				}
				c.JSON(http.StatusOK, gin.H{"ok": true})
			})

			req := httptest.NewRequest(http.MethodPost, "/x", nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			// Record runs async; the row is written after the handler. Poll briefly.
			row := waitForRow(t, db)
			if row.Status != tc.want {
				t.Errorf("status = %q, want %q", row.Status, tc.want)
			}
		})
	}
}

// TestAudit_BodyFieldTargetPreservesBody verifies the audit middleware can read a
// body field for the target WITHOUT consuming the body the handler needs.
func TestAudit_BodyFieldTargetPreservesBody(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testDB(t)

	var handlerSawEmail string
	r := gin.New()
	r.POST("/x", Audit(db, models.EventAdminSuspend, models.AccessCategoryAccount, BodyField("email")), func(c *gin.Context) {
		var body struct {
			Email string `json:"email"`
		}
		_ = c.ShouldBindJSON(&body)
		handlerSawEmail = body.Email
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	payload, _ := json.Marshal(map[string]string{"email": "victim@x.com"})
	req := httptest.NewRequest(http.MethodPost, "/x", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if handlerSawEmail != "victim@x.com" {
		t.Errorf("handler did not see the body: got %q (audit middleware consumed it?)", handlerSawEmail)
	}
	row := waitForRow(t, db)
	if row.Target != "victim@x.com" {
		t.Errorf("target = %q, want %q", row.Target, "victim@x.com")
	}
}

func TestCanonicalStatus(t *testing.T) {
	cases := map[string]string{
		"Successful": models.AccessStatusSuccessful,
		"failed":     models.AccessStatusFailed,
		"PENDING":    models.AccessStatusPending,
		"garbage":    "",
		"":           "",
	}
	for in, want := range cases {
		if got := canonicalStatus(in); got != want {
			t.Errorf("canonicalStatus(%q) = %q, want %q", in, got, want)
		}
	}
}

// waitForRow polls for the async-written row a bounded number of times.
func waitForRow(t *testing.T, db *gorm.DB) models.AdminAccessLog {
	t.Helper()
	for i := 0; i < 200; i++ {
		var row models.AdminAccessLog
		if err := db.First(&row).Error; err == nil {
			return row
		}
		time.Sleep(5 * time.Millisecond)
	}
	t.Fatal("no access-log row was written within the polling window")
	return models.AdminAccessLog{}
}
