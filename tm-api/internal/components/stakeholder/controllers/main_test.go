package controllers

import (
	"testing"

	coreModels "admin-panel-dashboard/internal/models"
	serverModels "admin-panel-dashboard/internal/server/models"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestInitRegistersCompleteComplianceWorkflow(t *testing.T) {
	gin.SetMode(gin.TestMode)
	database, err := gorm.Open(sqlite.Open("file:compliance-routes?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	router := gin.New()
	Init(router, &serverModels.Server{AdminDB: database, GC: &coreModels.GlobalConfig{}})

	routes := make(map[string]struct{})
	for _, route := range router.Routes() {
		routes[route.Method+" "+route.Path] = struct{}{}
	}
	expected := []string{
		"POST /api/v1/stakeholder/admin/custodian-compliance",
		"GET /api/v1/stakeholder/admin/custodian-compliance",
		"GET /api/v1/stakeholder/admin/compliance-templates",
		"POST /api/v1/stakeholder/admin/compliance-templates",
		"PUT /api/v1/stakeholder/admin/compliance-templates/:id",
		"GET /api/v1/stakeholder/admin/compliance-requirements",
		"POST /api/v1/stakeholder/admin/compliance-requirements",
		"GET /api/v1/stakeholder/admin/compliance-requirements/:id",
		"PUT /api/v1/stakeholder/admin/compliance-requirements/:id/review/start",
		"PUT /api/v1/stakeholder/admin/compliance-requirements/:id/review",
		"GET /api/v1/stakeholder/admin/compliance-requirements/:id/documents/:doc_id/download",
		"GET /api/v1/stakeholder/org/compliance",
		"PUT /api/v1/stakeholder/org/compliance/:item_id",
		"POST /api/v1/stakeholder/org/compliance/:item_id/submit",
		"POST /api/v1/stakeholder/org/compliance/documents/upload",
		"GET /api/v1/stakeholder/org/compliance/documents/:doc_id/download",
	}
	for _, route := range expected {
		if _, ok := routes[route]; !ok {
			t.Errorf("route %q is not registered", route)
		}
	}
}
