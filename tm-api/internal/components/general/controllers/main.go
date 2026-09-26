package general

import (
	"admin-panel-dashboard/internal/components/accesslog"
	adminServices "admin-panel-dashboard/internal/components/general/services"
	"admin-panel-dashboard/internal/components/observability"
	"admin-panel-dashboard/internal/middleware"
	"admin-panel-dashboard/internal/models"
	"admin-panel-dashboard/internal/observe"
	serverModels "admin-panel-dashboard/internal/server/models"
	"log"
	"os"

	"github.com/ecnepsnai/discord"
	"github.com/gin-gonic/gin"
)

// Init initializes
func Init(router *gin.Engine, s *serverModels.Server) {

	// login callback url

	apiV1 := router.Group("/api/v1")

	// bodyTarget extracts a target identifier from a JSON body field for the audit
	// log, without consuming the body the handler needs (it restores c.Request.Body).
	m := models.AccessCategoryAccount

	apiV1.POST("/admin/invite", middleware.AuthenticateSuperAdmin(s.AdminDB),
		accesslog.Audit(s.AdminDB, models.EventAdminInvite, m, accesslog.BodyField("email_or_username")),
		adminServices.AddAdmin(s))
	// apiV1.POST("/admin/invite/accept", AcceptRejectInvitation())
	apiV1.PATCH("/admin/suspend", middleware.AuthenticateSuperAdmin(s.AdminDB),
		accesslog.Audit(s.AdminDB, models.EventAdminSuspend, m, accesslog.BodyField("email")),
		adminServices.SuspendAdmin(s))
	// Suspend/lift are distinct endpoints (not a toggle) - each requires
	// its own mandatory, logged reason. While suspended, app-backend
	// blocks the user's wallet(s) from either side of any transaction.
	apiV1.PATCH("/admin/users/suspend", middleware.AuthenticateSuperAdmin(s.AdminDB),
		accesslog.Audit(s.AdminDB, models.EventUserSuspend, m, accesslog.BodyField("email")),
		adminServices.SuspendUser(s))
	apiV1.PATCH("/admin/users/lift-suspension", middleware.AuthenticateSuperAdmin(s.AdminDB),
		accesslog.Audit(s.AdminDB, models.EventUserLiftSuspend, m, accesslog.BodyField("email")),
		adminServices.LiftUserSuspension(s))
	apiV1.GET("/users/suspension-history/:email", middleware.AuthenticateSuperAdmin(s.AdminDB), adminServices.GetUserSuspensionHistoryByEmail(s.AdminDB))
	apiV1.PATCH("/admin/unsuspend", middleware.AuthenticateSuperAdmin(s.AdminDB),
		accesslog.Audit(s.AdminDB, models.EventAdminUnsuspend, m, accesslog.BodyField("email")),
		adminServices.UnsuspendAdmin(s))
	apiV1.PATCH("/admin/modify/status", middleware.AuthenticateSuperAdmin(s.AdminDB),
		accesslog.Audit(s.AdminDB, models.EventAdminStatusChange, m, accesslog.BodyField("email")),
		adminServices.ChangeAdminStatus(s))
	apiV1.GET("/suspension/history", middleware.AuthenticateSuperAdmin(s.AdminDB), adminServices.GetAllSuspensionHistory(s.AdminDB))
	apiV1.GET("/suspension/history/:email", middleware.AuthenticateSuperAdmin(s.AdminDB), adminServices.GetSuspensionHistoryByEmail(s.AdminDB))
	apiV1.GET("/admin/list", middleware.AuthenticateSuperAdmin(s.AdminDB), adminServices.ListAdminsHandler(s))
	apiV1.POST("admin/remove", middleware.AuthenticateSuperAdmin(s.AdminDB),
		accesslog.Audit(s.AdminDB, models.EventAdminRemove, m, accesslog.BodyField("email")),
		adminServices.RemoveAdmin(s.AdminDB))
	apiV1.GET("/payment/history", middleware.JwtTokenAuthMiddleware(s.AdminDB), adminServices.GetPaymentHistory(s.TrovoWalletDB))
	// apiV1.POST("/configurations/roles", middleware.JwtTokenAuthMiddleware(s.AdminDB), controllers.CreateRole)

	// all configurations
	apiV1.GET("/configurations/all", middleware.JwtTokenAuthMiddleware(s.AdminDB), adminServices.GetAllConfigurations(s))
	apiV1.POST("/configurations/bulk", middleware.JwtTokenAuthMiddleware(s.AdminDB),
		accesslog.Audit(s.AdminDB, models.EventConfigBulkUpdate, models.AccessCategoryConfig, nil),
		adminServices.BulkCreateOrUpdateConfigurations(s))
	apiV1.GET("/configurations", middleware.JwtTokenAuthMiddleware(s.AdminDB), adminServices.ReadConfigurations(s)) // Read configurations
	apiV1.DELETE("/configurations", middleware.JwtTokenAuthMiddleware(s.AdminDB),
		accesslog.Audit(s.AdminDB, models.EventConfigDelete, models.AccessCategoryConfig, func(c *gin.Context) string { return c.Query("subject") }),
		adminServices.DeleteConfiguration(s)) // Delete configuration

	// System health module: the Trovo Manager observability screens. The admin
	// API queries Prometheus, Loki and the crash reporter server-side, so the
	// browser never holds credentials for the monitoring stack.
	observability.Init(apiV1, s)

	// Admin access / audit trail (security log). Super-admin only.
	apiV1.GET("/audit-trail", middleware.AuthenticateSuperAdmin(s.AdminDB), accesslog.ListHandler(s))
	apiV1.GET("/audit-trail/:id", middleware.AuthenticateSuperAdmin(s.AdminDB), accesslog.DetailHandler(s))

}
func LogDiscordError(msg string) {
	// Also count it, so the failure shows on the error-rate dashboard
	// instead of only in Discord. Additive: Discord is unchanged.
	observe.RecordHandledFailure(msg)

	discord.WebhookURL = "https://discord.com/api/webhooks/865931042795290636/jObHzZWdnbhX1jomOSZQX8Ip5AXLArh87PI4-ZQ8u6ssnRbZuVdY_iPxz5qoWkHUlZwS"
	if len(os.Getenv("500_ERROR_WEBHOOK")) > 50 {
		discord.WebhookURL = os.Getenv("500_ERROR_WEBHOOK")
	}
	err := discord.Say(msg)
	if err != nil {
		log.Println("Failed to send discord error message", err)
		return
	}
}
