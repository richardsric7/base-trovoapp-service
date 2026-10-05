package controllers

import (
	"admin-panel-dashboard/internal/components/accesslog"
	"admin-panel-dashboard/internal/components/publicmarkets"
	"admin-panel-dashboard/internal/middleware"
	"admin-panel-dashboard/internal/models"
	serverModels "admin-panel-dashboard/internal/server/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Service is the Public Markets service on app-backend's database.
func Service(s *serverModels.Server) *publicmarkets.Service {
	return &publicmarkets.Service{DB: s.TrovoWalletDB, Chain: publicmarkets.NewContractReader()}
}

// can checks the calling admin's role for a permission.
func can(adminDB *gorm.DB) func(*gin.Context, string) bool {
	return func(c *gin.Context, permission string) bool {
		id, ok := c.Get("trovo_admin_id")
		if !ok {
			return false
		}
		var admin models.AdminUser
		if adminDB.First(&admin, "id = ?", id).Error != nil {
			return false
		}
		return models.HasPermission(adminDB, &admin, permission)
	}
}

// Init registers the Public Markets admin endpoints (Trovo admins; every
// change is audited).
func Init(router *gin.Engine, s *serverModels.Server) {
	h := publicmarkets.Handlers{S: Service(s), Can: can(s.AdminDB)}
	auth := middleware.JwtTokenAuthMiddleware(s.AdminDB)
	audit := func(event string, target func(*gin.Context) string) gin.HandlerFunc {
		return accesslog.Audit(s.AdminDB, event, models.AccessCategoryPublicMarkets, target)
	}
	id := accesslog.Param("id")

	g := router.Group("/api/v1/public-markets", auth)
	g.GET("/me", h.Me)
	g.GET("/overview", h.Overview)
	g.GET("/health", h.Health)
	g.GET("/jobs", h.Jobs)
	g.GET("/partner-events", h.PartnerEvents)

	g.GET("/assets", h.Assets)
	g.POST("/assets", audit(models.EventPMAssetChange, accesslog.BodyField("assetCode")), h.CreateAsset)
	g.GET("/assets/:id", h.Asset)
	g.PUT("/assets/:id", audit(models.EventPMAssetChange, id), h.UpdateAsset)
	g.POST("/assets/:id/contract", audit(models.EventPMContract, id), h.RegisterContract)
	g.POST("/assets/:id/go-live", audit(models.EventPMAssetLive, id), h.GoLive)
	g.POST("/assets/:id/halt", audit(models.EventPMAssetHalt, id), h.Halt)
	g.POST("/assets/:id/resume", audit(models.EventPMAssetResume, id), h.Resume)
	g.POST("/assets/:id/price", audit(models.EventPMPrice, id), h.SetPrice)
	g.POST("/assets/:id/position", audit(models.EventPMPosition, id), h.RecordPosition)
	g.GET("/assets/:id/prices", h.PriceHistory)
	g.GET("/prices", h.Prices)

	g.GET("/orders", h.Orders)
	g.GET("/orders/:id", h.Order)
	g.GET("/batches", h.Batches)
	g.POST("/batches/:id/approve", audit(models.EventPMBatchApprove, id), h.ApproveBatch)
	g.POST("/batches/:id/reject", audit(models.EventPMBatchReject, id), h.RejectBatch)
	g.POST("/batches/:id/roll", audit(models.EventPMBatchRoll, id), h.RollBatch)
	g.GET("/instructions", h.Instructions)
	g.POST("/instructions/:id/retry", audit(models.EventPMInstruction, id), h.RetryInstruction)
	g.POST("/instructions/:id/handled", audit(models.EventPMInstruction, id), h.MarkHandled)
	g.POST("/instructions/:id/outcome", audit(models.EventPMOutcome, id), h.RecordOutcome)

	g.GET("/reconciliation", h.Reconciliation)
	g.GET("/reconciliation/runs", h.ReconciliationRuns)
	g.POST("/reconciliation/run", audit(models.EventPMReconcile, accesslog.BodyField("target")), h.RunReconciliation)
	g.POST("/position-feed", audit(models.EventPMReconcile, accesslog.BodyField("target")), h.PositionFeed)

	g.GET("/corporate-actions", h.CorporateActions)
	g.POST("/corporate-actions", audit(models.EventPMActionDeclare, accesslog.BodyField("assetCode")), h.Declare)
	g.GET("/corporate-actions/:id", h.CorporateAction)
	g.POST("/corporate-actions/:id/approve", audit(models.EventPMActionApprove, id), h.ApproveAction)
	g.POST("/corporate-actions/:id/cancel", audit(models.EventPMActionCancel, id), h.CancelAction)
	g.GET("/confirmations", h.Confirmations)

	g.GET("/exchanges", h.Exchanges)
	g.POST("/exchanges", audit(models.EventPMExchangeChange, accesslog.BodyField("serviceLinkId")), h.Onboard)
	g.GET("/exchanges/:id", h.Exchange)
	g.PUT("/exchanges/:id", audit(models.EventPMExchangeChange, id), h.UpdateExchange)
	g.POST("/exchanges/:id/rotate-secret", audit(models.EventPMExchangeSecret, id), h.RotateSecret)
	g.POST("/exchanges/:id/suspend", audit(models.EventPMExchangeStatus, id), h.SuspendExchange)
	g.POST("/exchanges/:id/activate", audit(models.EventPMExchangeStatus, id), h.ActivateExchange)
	g.POST("/exchanges/:id/withdrawals", audit(models.EventPMExchangeWithdraw, id), h.Withdraw)
	g.POST("/exchanges/:id/replay-dead-letters", audit(models.EventPMWebhookReplay, id), h.ReplayDeadLetters)
	g.POST("/webhooks/:id/replay", audit(models.EventPMWebhookReplay, id), h.ReplayWebhook)
	g.GET("/wallets", h.Wallets)

	g.GET("/settings", h.Settings)
	g.PUT("/settings", audit(models.EventPMSettings, nil), h.UpdateSettings)
	g.GET("/dealing-members", h.DealingMembers)
	g.POST("/dealing-members", audit(models.EventPMPartnerChange, accesslog.BodyField("code")), h.SaveDealingMember)
	g.PUT("/dealing-members/:id", audit(models.EventPMPartnerChange, id), h.SaveDealingMember)
	g.GET("/custodians", h.Custodians)
	g.PUT("/custodians/:id", audit(models.EventPMPartnerChange, id), h.ConfigureCustodian)
}
