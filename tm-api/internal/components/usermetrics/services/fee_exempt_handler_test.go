package usermetrics

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"admin-panel-dashboard/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestFeeExemptUsersAdmin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, err := gorm.Open(sqlite.Open("file:"+url.QueryEscape(t.Name())+"?mode=memory&cache=shared"), &gorm.Config{Logger: logger.Discard})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&models.FeeExemptUser{}); err != nil {
		t.Fatal(err)
	}
	db.Exec(`CREATE TABLE users (username TEXT PRIMARY KEY)`)
	db.Exec(`INSERT INTO users (username) VALUES ('Treasury'), ('alice')`)

	router := func(authType string) *gin.Engine {
		r := gin.New()
		r.Use(func(c *gin.Context) {
			c.Set("auth_type", authType)
			c.Set("trovo_admin_email", "ops@trovo.test")
		})
		r.GET("/fee/exempt-users", GetFeeExemptUsersHandler(db))
		r.POST("/fee/exempt-users", AddFeeExemptUserHandler(db))
		r.DELETE("/fee/exempt-users/:username", RemoveFeeExemptUserHandler(db))
		return r
	}
	do := func(r *gin.Engine, method, path string, body interface{}) *httptest.ResponseRecorder {
		var buf bytes.Buffer
		if body != nil {
			json.NewEncoder(&buf).Encode(body)
		}
		req := httptest.NewRequest(method, path, &buf)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		return w
	}
	admin, member := router("trovo_admin"), router("organization_member")

	if w := do(member, http.MethodPost, "/fee/exempt-users", map[string]string{"username": "treasury", "reason": "x"}); w.Code != http.StatusForbidden {
		t.Fatalf("an organization member added an exemption: %d", w.Code)
	}
	if w := do(admin, http.MethodPost, "/fee/exempt-users", map[string]string{"username": "nobody", "reason": "x"}); w.Code != http.StatusBadRequest {
		t.Fatalf("an unknown username: %d %s", w.Code, w.Body)
	}
	if w := do(admin, http.MethodPost, "/fee/exempt-users", map[string]string{"username": "treasury", "reason": "platform market maker"}); w.Code != http.StatusOK {
		t.Fatalf("add: %d %s", w.Code, w.Body)
	}
	var rows []models.FeeExemptUser
	db.Find(&rows)
	if len(rows) != 1 || rows[0].Username != "Treasury" || rows[0].AddedBy != "ops@trovo.test" || rows[0].Reason != "platform market maker" {
		t.Fatalf("stored: %+v", rows)
	}
	w := do(admin, http.MethodGet, "/fee/exempt-users", nil)
	var listed struct {
		Data []models.FeeExemptUser `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &listed)
	if w.Code != http.StatusOK || len(listed.Data) != 1 {
		t.Fatalf("list: %d %s", w.Code, w.Body)
	}
	if w := do(admin, http.MethodDelete, "/fee/exempt-users/TREASURY", nil); w.Code != http.StatusOK {
		t.Fatalf("remove: %d", w.Code)
	}
	var n int64
	db.Model(&models.FeeExemptUser{}).Count(&n)
	if n != 0 {
		t.Fatalf("%d exemptions left after removal", n)
	}
}
