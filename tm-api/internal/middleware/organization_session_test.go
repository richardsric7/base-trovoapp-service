package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"admin-panel-dashboard/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestOrganizationSessionVersionRevokesIssuedToken(t *testing.T) {
	t.Setenv("JWT_SECRET", "organization-session-test-secret")
	gin.SetMode(gin.TestMode)
	database, err := gorm.Open(sqlite.Open("file:"+uuid.NewString()+"?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	if err := database.AutoMigrate(&models.Organization{}, &models.OrganizationMember{}); err != nil {
		t.Fatalf("migrate database: %v", err)
	}
	organization := models.Organization{
		ID: uuid.NewString(), Name: "Session Org", Email: uuid.NewString() + "@example.com",
		Type: models.OrganizationTypeAssetManager, Status: models.OrganizationStatusActive, CreatedBy: "test",
	}
	stakeholderID := uint64(20)
	stakeholderType := "asset_manager"
	organization.StakeholderID = &stakeholderID
	organization.StakeholderType = &stakeholderType
	member := models.OrganizationMember{
		ID: uuid.NewString(), OrganizationID: organization.ID, Email: uuid.NewString() + "@example.com",
		Role: string(models.OrganizationSuperAdmin), Status: models.OrganizationStatusActive, CreatedBy: "test",
	}
	if err := database.Create(&organization).Error; err != nil {
		t.Fatalf("seed organization: %v", err)
	}
	if err := database.Create(&member).Error; err != nil {
		t.Fatalf("seed member: %v", err)
	}
	token, err := GenerateOrganizationTokenWithDB(database, &member)
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}

	requestStatus := func(value string) int {
		router := gin.New()
		router.GET("/protected", OrganizationAuthMiddleware(database), func(c *gin.Context) { c.Status(http.StatusNoContent) })
		request := httptest.NewRequest(http.MethodGet, "/protected", nil)
		request.Header.Set("Authorization", "Bearer "+value)
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)
		return response.Code
	}
	stakeholderRequestStatus := func(value string) int {
		router := gin.New()
		router.GET("/protected", StakeholderAuthMiddleware(database), func(c *gin.Context) { c.Status(http.StatusNoContent) })
		request := httptest.NewRequest(http.MethodGet, "/protected", nil)
		request.Header.Set("Authorization", "Bearer "+value)
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)
		return response.Code
	}
	if status := requestStatus(token); status != http.StatusNoContent {
		t.Fatalf("fresh token status = %d, want 204", status)
	}
	if status := stakeholderRequestStatus(token); status != http.StatusNoContent {
		t.Fatalf("fresh stakeholder token status = %d, want 204", status)
	}
	if err := database.Model(&models.OrganizationMember{}).Where("id = ?", member.ID).
		UpdateColumn("session_version", gorm.Expr("session_version + 1")).Error; err != nil {
		t.Fatalf("revoke session: %v", err)
	}
	if status := requestStatus(token); status != http.StatusUnauthorized {
		t.Fatalf("revoked token status = %d, want 401", status)
	}
	if status := stakeholderRequestStatus(token); status != http.StatusUnauthorized {
		t.Fatalf("revoked stakeholder token status = %d, want 401", status)
	}
	if err := database.First(&member, "id = ?", member.ID).Error; err != nil {
		t.Fatalf("reload member: %v", err)
	}
	newToken, err := GenerateOrganizationTokenWithDB(database, &member)
	if err != nil {
		t.Fatalf("generate replacement token: %v", err)
	}
	if status := requestStatus(newToken); status != http.StatusNoContent {
		t.Fatalf("replacement token status = %d, want 204", status)
	}
	if status := stakeholderRequestStatus(newToken); status != http.StatusNoContent {
		t.Fatalf("replacement stakeholder token status = %d, want 204", status)
	}
}

func TestActiveOrganizationAuthDoesNotRequirePortalStakeholderRole(t *testing.T) {
	t.Setenv("JWT_SECRET", "organization-compliance-test-secret")
	gin.SetMode(gin.TestMode)
	database, err := gorm.Open(sqlite.Open("file:"+uuid.NewString()+"?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	if err := database.AutoMigrate(&models.Organization{}, &models.OrganizationMember{}); err != nil {
		t.Fatalf("migrate database: %v", err)
	}
	stakeholderID := uint64(44)
	stakeholderType := models.StakeholderTypeLegal
	organization := models.Organization{
		ID: uuid.NewString(), Name: "Legal Org", Email: uuid.NewString() + "@example.com",
		Type: stakeholderType, Status: models.OrganizationStatusActive, CreatedBy: "test",
		StakeholderID: &stakeholderID, StakeholderType: &stakeholderType,
	}
	member := models.OrganizationMember{
		ID: uuid.NewString(), OrganizationID: organization.ID, Email: uuid.NewString() + "@example.com",
		Role: string(models.OrganizationSuperAdmin), Status: models.OrganizationStatusActive, CreatedBy: "test",
	}
	if err := database.Create(&organization).Error; err != nil {
		t.Fatal(err)
	}
	if err := database.Create(&member).Error; err != nil {
		t.Fatal(err)
	}
	token, err := GenerateOrganizationTokenWithDB(database, &member)
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}

	request := func(middlewareChain ...gin.HandlerFunc) int {
		router := gin.New()
		router.GET("/protected", append(middlewareChain, func(c *gin.Context) { c.Status(http.StatusNoContent) })...)
		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", token)
		response := httptest.NewRecorder()
		router.ServeHTTP(response, req)
		return response.Code
	}
	if status := request(OrganizationAuthMiddleware(database), RequireActiveOrganization(database)); status != http.StatusNoContent {
		t.Fatalf("generic organization compliance auth status=%d, want 204", status)
	}
	if status := request(StakeholderAuthMiddleware(database)); status != http.StatusForbidden {
		t.Fatalf("portal stakeholder auth status=%d, want 403 for unsupported role", status)
	}
	if err := database.Model(&organization).Update("status", models.OrganizationStatusSuspended).Error; err != nil {
		t.Fatal(err)
	}
	if status := request(OrganizationAuthMiddleware(database), RequireActiveOrganization(database)); status != http.StatusForbidden {
		t.Fatalf("suspended organization status=%d, want 403", status)
	}
}
