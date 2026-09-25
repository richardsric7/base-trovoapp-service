package controllers

import (
	"admin-panel-dashboard/internal/components/accesslog"
	"admin-panel-dashboard/internal/components/organizations/handlers"
	"admin-panel-dashboard/internal/middleware"
	"admin-panel-dashboard/internal/models"
	serverModels "admin-panel-dashboard/internal/server/models"

	"github.com/gin-gonic/gin"
)

// Init initializes the organization routes
func Init(router *gin.Engine, s *serverModels.Server) {
	apiV1 := router.Group("/api/v1")

	apiV1.POST("/organizations/invite/root-user", middleware.JwtTokenAuthMiddleware(s.AdminDB),
		accesslog.Audit(s.AdminDB, models.EventOrgCreate, models.AccessCategoryOrg, accesslog.BodyField("organization_name")),
		handlers.TrovoAdminInvitesOrganizationRootUser(s))
	apiV1.PUT("/organizations/update/:organization_id", middleware.JwtTokenAuthMiddleware(s.AdminDB),
		accesslog.Audit(s.AdminDB, models.EventOrgUpdate, models.AccessCategoryOrg, accesslog.Param("organization_id")),
		handlers.TrovoAdminUpdatesOrganizationRootUser(s))
	apiV1.PUT("/organizations/deactivate/:organization_id", middleware.JwtTokenAuthMiddleware(s.AdminDB),
		accesslog.Audit(s.AdminDB, models.EventOrgDeactivate, models.AccessCategoryOrg, accesslog.Param("organization_id")),
		handlers.TrovoAdminDeactivatesOrganization(s))
	apiV1.POST("/organizations/invite/validate/:invite_id", handlers.RootUserValidateInvite(s))
	//apiV1.POST("/organizations/invite/validate-email/:invite_id", handlers.OrganizationUserValidatesEmail(s))
	apiV1.POST("/organizations/setup-password/:invite_id", handlers.SetupUserPassword(s))
	//apiV1.POST("/organizations/members", middleware.OrganizationAuthMiddleware(s.AdminDB), handlers.InviteOrganizationMember(s))
	apiV1.POST("/organizations/members", middleware.AllowOrgOrTrovoAdminNormalized(s.AdminDB),
		accesslog.Audit(s.AdminDB, models.EventOrgMemberInvite, models.AccessCategoryOrg, accesslog.BodyField("email")),
		handlers.InviteOrganizationMember(s))
	apiV1.POST("/organizations/member/login", handlers.OrganizationMemberLogin(s))
	apiV1.POST("/organizations/member/logout", middleware.OrganizationAuthMiddleware(s.AdminDB), handlers.OrganizationMemberLogout(s))

	// Wallet linking: Step 1 link (DB only, once), Step 2 request-auth, Step 3 request-login + verify-login (token)
	//apiV1.POST("/organizations/members/link-wallet", middleware.AllowOrgOrTrovoAdminNormalized(s.AdminDB), handlers.InitiateWalletLinkHandler(s))
	apiV1.POST("/organizations/members/link-wallet/request-auth", middleware.AllowOrgOrTrovoAdminNormalized(s.AdminDB), handlers.RequestWalletLinkAuthHandler(s)) //
	//apiV1.POST("/organizations/members/link-wallet/request-login", middleware.AllowOrgOrTrovoAdminNormalized(s.AdminDB), handlers.RequestWalletLinkLoginHandler(s))

	//apiV1.POST("/organizations/members/link-wallet/verify-login", middleware.AllowOrgOrTrovoAdminNormalized(s.AdminDB), handlers.VerifyWalletLinkLoginHandler(s))
	apiV1.POST("/organizations/members/link-wallet/verify", middleware.AllowOrgOrTrovoAdminNormalized(s.AdminDB),
		accesslog.Audit(s.AdminDB, models.EventWalletLink, models.AccessCategoryOrg, nil),
		handlers.VerifyWalletLinkHandler(s)) //

	apiV1.GET("/organizations", middleware.JwtTokenAuthMiddleware(s.AdminDB), handlers.ListOrganizationsPaginatedHandler(s))
	//apiV1.POST("/organizations/resend-password-setup", handlers.ResendPasswordSetupToken(s))
	apiV1.GET("/organizations/members", middleware.AllowOrgOrTrovoAdminNormalized(s.AdminDB), handlers.ListOrganizationMembers(s))
	apiV1.GET("/organizations/details", middleware.AllowOrgOrTrovoAdminNormalized(s.AdminDB), handlers.GetOrganizationDetails(s))

	// Member password reset routes
	router.POST("/organizations/members/reset-password", handlers.RequestMemberPasswordReset)
}
