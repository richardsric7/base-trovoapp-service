package controllers

import (
	"admin-panel-dashboard/internal/components/accesslog"
	stakeholderDB "admin-panel-dashboard/internal/components/stakeholder/db"
	"admin-panel-dashboard/internal/components/stakeholder/handlers"
	"admin-panel-dashboard/internal/components/stakeholder/models"
	"admin-panel-dashboard/internal/components/stakeholder/services"
	"admin-panel-dashboard/internal/middleware"
	rootModels "admin-panel-dashboard/internal/models"
	serverModels "admin-panel-dashboard/internal/server/models"

	"github.com/gin-gonic/gin"
)

func Init(router *gin.Engine, s *serverModels.Server) {
	apiV1 := router.Group("/api/v1")
	apiV1.GET("/stakeholder/shared/health", handlers.NewHandler(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil).Health)

	tokenizationClient := stakeholderDB.NewGormTokenizationReadClient(s.TrovoWalletDB)
	auditService := services.NewAuditService(s.AdminDB)
	notificationService := services.NewNotificationService(s.AdminDB)
	authorizationService := services.NewAuthorizationService(s.AdminDB, s.GC.ServiceLink)
	assetService := services.NewAssetService(s.AdminDB, tokenizationClient, auditService)
	profileService := services.NewProfileService(s.AdminDB, auditService)
	documentStorage := services.NewWalletDocumentStorageClient(s.GC.ServiceLink)
	distributionPayouts := services.NewWalletDistributionPayoutClient(s.GC.ServiceLink)
	documentService := services.NewDocumentService(s.AdminDB, assetService, auditService, documentStorage)
	fundReleaseService := services.NewFundReleaseService(s.AdminDB, assetService, auditService, notificationService, authorizationService, documentService)
	dueDiligenceService := services.NewDueDiligenceService(s.AdminDB, assetService, documentService, auditService, notificationService)
	financialService := services.NewFinancialService(s.AdminDB, assetService, auditService, notificationService, authorizationService, documentService, distributionPayouts)
	custodianService := services.NewCustodianService(s.AdminDB, assetService, documentService, auditService)
	operationsService := services.NewOperationsService(s.AdminDB, assetService, documentService, auditService, notificationService)
	dashboardService := services.NewDashboardService(s.AdminDB, assetService)
	structuringService := services.NewStructuringService(s.AdminDB, assetService, documentService, auditService)
	fundManagementService := services.NewFundManagementService(s.AdminDB, s.TrovoWalletDB, assetService)
	complianceRequirementService := services.NewComplianceRequirementService(s.AdminDB, auditService, documentService)

	h := handlers.NewHandler(
		s.AdminDB,
		assetService,
		profileService,
		authorizationService,
		notificationService,
		fundReleaseService,
		dueDiligenceService,
		financialService,
		custodianService,
		documentService,
		operationsService,
		dashboardService,
		structuringService,
		fundManagementService,
		complianceRequirementService,
	)

	admin := apiV1.Group("/stakeholder/admin")
	admin.POST("/asset-assignments", middleware.JwtTokenAuthMiddleware(s.AdminDB),
		accesslog.Audit(s.AdminDB, rootModels.EventAssetAssignment, rootModels.AccessCategoryAsset, accesslog.BodyField("asset_id")),
		h.CreateAssetAssignment)
	admin.POST("/custodian-compliance", middleware.AuthenticateSuperAdmin(s.AdminDB),
		accesslog.Audit(s.AdminDB, rootModels.EventComplianceCreate, rootModels.AccessCategoryAsset, accesslog.BodyField("asset_id")),
		h.CreateCompliance)
	admin.GET("/custodian-compliance", middleware.AuthenticateSuperAdmin(s.AdminDB), h.ListAdminCompliance)
	admin.GET("/compliance-templates", middleware.AuthenticateSuperAdmin(s.AdminDB), h.ListComplianceTemplates)
	admin.POST("/compliance-templates", middleware.AuthenticateSuperAdmin(s.AdminDB), h.CreateComplianceTemplate)
	admin.PUT("/compliance-templates/:id", middleware.AuthenticateSuperAdmin(s.AdminDB), h.UpdateComplianceTemplate)
	admin.GET("/compliance-requirements", middleware.AuthenticateSuperAdmin(s.AdminDB), h.ListComplianceRequirements)
	admin.POST("/compliance-requirements", middleware.AuthenticateSuperAdmin(s.AdminDB), h.CreateComplianceRequirement)
	admin.GET("/compliance-requirements/:id", middleware.AuthenticateSuperAdmin(s.AdminDB), h.GetComplianceRequirement)
	admin.PUT("/compliance-requirements/:id/review/start", middleware.AuthenticateSuperAdmin(s.AdminDB), h.StartComplianceRequirementReview)
	admin.PUT("/compliance-requirements/:id/review", middleware.AuthenticateSuperAdmin(s.AdminDB), h.ReviewComplianceRequirement)
	admin.GET("/compliance-requirements/:id/documents/:doc_id/download", middleware.AuthenticateSuperAdmin(s.AdminDB), h.DownloadComplianceRequirementDocumentForAdmin)

	// Compliance applies to every authenticated organization, including types
	// that do not have a stakeholder-portal dashboard role.
	org := apiV1.Group("/stakeholder/org", middleware.OrganizationAuthMiddleware(s.AdminDB), middleware.RequireActiveOrganization(s.AdminDB))
	org.GET("/compliance", h.ListOrgCompliance)
	org.POST("/compliance/:item_id/submit", h.SubmitOrgCompliance)
	org.PUT("/compliance/:item_id", h.UpdateOrgCompliance)
	org.POST("/compliance/documents/upload", h.UploadOrganizationComplianceDocument)
	org.GET("/compliance/documents/:doc_id/download", h.DownloadOrganizationComplianceDocument)

	stakeholder := apiV1.Group("/stakeholder", middleware.StakeholderAuthMiddleware(s.AdminDB))
	shared := stakeholder.Group("/shared")
	shared.GET("/profile", h.GetProfile)
	shared.PUT("/profile", h.UpdateProfile)
	shared.PUT("/profile/notifications", h.UpdateNotificationPreferences)
	shared.GET("/assets", h.ListAssets)
	shared.GET("/assets/:asset_id", h.GetAsset)
	shared.GET("/assets/:asset_id/activities", h.ListAssetActivities)
	shared.POST("/authorizations", h.CreateAuthorization)
	shared.POST("/authorizations/:challenge_id/verify", h.VerifyAuthorization)
	shared.POST("/documents", h.CreateDocument)
	shared.POST("/documents/upload", h.UploadDocument)
	shared.GET("/document-categories", h.DocumentCategories)
	shared.GET("/documents", h.ListDocuments)
	shared.GET("/documents/:doc_id", h.GetDocument)
	shared.GET("/documents/:doc_id/download", h.DownloadDocument)
	shared.GET("/notifications", h.ListNotifications)
	shared.PUT("/notifications/:id/read", h.MarkNotificationRead)
	shared.GET("/audit-trail", h.ListAuditTrail)

	trustee := stakeholder.Group("/trustee", middleware.RequireStakeholderRole(models.DashboardRoleTrustee))
	trustee.GET("/dashboard", h.Dashboard)
	trustee.GET("/fund-management/summary", h.GetTrusteeFundManagementSummary)
	trustee.GET("/due-diligence/:asset_id", h.GetDueDiligence)
	trustee.PUT("/due-diligence/:asset_id/categories", h.UpdateDueDiligenceCategory)
	trustee.PUT("/due-diligence/:asset_id/items/:item_id", h.UpdateDueDiligenceItem)
	trustee.POST("/due-diligence/:asset_id/approve", h.ApproveDueDiligence)
	trustee.POST("/due-diligence/:asset_id/reject", h.RejectDueDiligence)
	trustee.GET("/fund-releases", h.ListFundReleases)
	trustee.GET("/fund-releases/:request_id", h.GetFundRelease)
	trustee.POST("/fund-releases/:request_id/approve", h.ApproveFundRelease)
	trustee.POST("/fund-releases/:request_id/reject", h.RejectFundRelease)
	trustee.GET("/distributions", h.ListDistributions)
	trustee.GET("/distributions/:dist_id", h.GetDistribution)
	trustee.GET("/distributions/:dist_id/payouts/download", h.DownloadDistributionPayouts)
	trustee.POST("/distributions/:dist_id/authorize", h.AuthorizeDistribution)
	trustee.POST("/distributions/:dist_id/reject", h.RejectDistribution)
	trustee.GET("/distributions/history", h.ListDistributionHistory)

	custodian := stakeholder.Group("/custodian", middleware.RequireStakeholderRole(models.DashboardRoleAssetCustodian))
	custodian.GET("/dashboard", h.Dashboard)
	custodian.POST("/accounts", h.CreateAccount)
	custodian.GET("/accounts", h.ListAccounts)
	custodian.GET("/accounts/:account_id", h.GetAccount)
	custodian.PUT("/accounts/:account_id", h.UpdateAccount)
	custodian.GET("/fund-releases", h.ListFundReleases)
	custodian.GET("/fund-releases/:request_id", h.GetFundRelease)
	custodian.POST("/fund-releases/:request_id/execute", h.ExecuteFundRelease)
	custodian.PUT("/fund-releases/:request_id/status", h.UpdateFundReleaseStatus)
	custodian.GET("/compliance", h.ListCompliance)
	custodian.PUT("/compliance/:item_id", h.UpdateCompliance)

	manager := stakeholder.Group("/asset-manager", middleware.RequireStakeholderRole(models.DashboardRoleAssetManager))
	manager.GET("/dashboard", h.Dashboard)
	manager.GET("/assets", h.ListAssets)
	manager.PUT("/assets/:asset_id", h.UpdateAssetOperation)
	manager.GET("/revenue", h.ListRevenue)
	manager.POST("/revenue", h.CreateRevenue)
	manager.POST("/revenue/submit-distribution", h.SubmitDistribution)
	manager.GET("/valuations/:asset_id", h.ListValuations)
	manager.POST("/valuations", h.CreateValuation)
	manager.POST("/valuations/:id/request-independent", h.RequestIndependentValuation)
	manager.POST("/fund-releases", h.CreateFundRelease)
	manager.GET("/fund-releases", h.ListFundReleases)
	manager.GET("/fund-releases/summary", h.GetAssetManagerFundSummary)
	manager.GET("/fund-releases/:request_id", h.GetFundRelease)
	manager.POST("/reports", h.CreateReport)
	manager.GET("/reports", h.ListReports)
	manager.GET("/report-types", h.ListReportTypes)
	manager.POST("/reports/:report_id/submit", h.SubmitReport)

	// Tier 1 — Legal Adviser & Financial Adviser: view assigned assets (via
	// /shared/*) + confirm their A5 workstream complete.
	legalAdviser := stakeholder.Group("/legal-adviser", middleware.RequireStakeholderRole(models.DashboardRoleLegalAdviser))
	legalAdviser.GET("/dashboard", h.Dashboard)
	legalAdviser.POST("/structuring/:asset_id/complete", h.ConfirmStructuring)

	financialAdviser := stakeholder.Group("/financial-adviser", middleware.RequireStakeholderRole(models.DashboardRoleFinancialAdviser))
	financialAdviser.GET("/dashboard", h.Dashboard)
	financialAdviser.POST("/structuring/:asset_id/complete", h.ConfirmStructuring)

	// Tier 2 — Issuing House & Rating Agency: view-only (assigned assets via
	// /shared/*, dashboard). No portal actions per PRD §2.3.
	issuingHouse := stakeholder.Group("/issuing-house", middleware.RequireStakeholderRole(models.DashboardRoleIssuingHouse))
	issuingHouse.GET("/dashboard", h.Dashboard)

	ratingAgency := stakeholder.Group("/rating-agency", middleware.RequireStakeholderRole(models.DashboardRoleRatingAgency))
	ratingAgency.GET("/dashboard", h.Dashboard)
}
