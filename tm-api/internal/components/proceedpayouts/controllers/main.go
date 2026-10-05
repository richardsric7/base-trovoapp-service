package controllers

import (
	"admin-panel-dashboard/internal/components/accesslog"
	"admin-panel-dashboard/internal/components/proceedpayouts"
	"admin-panel-dashboard/internal/middleware"
	"admin-panel-dashboard/internal/models"
	serverModels "admin-panel-dashboard/internal/server/models"

	"github.com/gin-gonic/gin"
)

// Service is the payout service on the server's databases.
func Service(s *serverModels.Server) *proceedpayouts.Service {
	return &proceedpayouts.Service{DB: s.TrovoWalletDB, AdminDB: s.AdminDB, Cache: s.GC.Cache}
}

// Init registers the proceeds payout admin endpoints (Trovo admins only;
// every change is audited).
func Init(router *gin.Engine, s *serverModels.Server) {
	h := proceedpayouts.Handlers{S: Service(s)}
	auth := middleware.JwtTokenAuthMiddleware(s.AdminDB)
	audit := func(event string, target func(*gin.Context) string) gin.HandlerFunc {
		return accesslog.Audit(s.AdminDB, event, models.AccessCategoryPayout, target)
	}
	id := accesslog.Param("id")
	item := accesslog.Param("itemId")

	g := router.Group("/api/v1/proceed-payouts")
	g.GET("", auth, h.List)
	g.GET("/engine", auth, h.Engine)
	g.POST("/engine/halt", auth, audit(models.EventPayoutEngineHalt, nil), h.Halt)
	g.POST("/engine/unhalt", auth, audit(models.EventPayoutEngineRun, nil), h.Unhalt)
	g.POST("/engine/sweep", auth, audit(models.EventPayoutSweep, accesslog.BodyField("token")), h.Sweep)
	g.GET("/fee-config", auth, h.FeeConfig)
	g.PUT("/fee-config", auth, audit(models.EventPayoutFeeConfig, accesslog.BodyField("feeWallet")), h.SetFeeConfig)
	g.GET("/reports/payouts", auth, h.PayoutsReport)
	g.GET("/reports/fees", auth, h.FeesReport)
	g.GET("/:id", auth, h.Detail)
	g.GET("/:id/items", auth, h.Items)
	g.POST("/:id/prepare", auth, audit(models.EventPayoutPrepare, id), h.Prepare)
	g.PUT("/:id/fee", auth, audit(models.EventPayoutFee, id), h.SetFee)
	g.POST("/:id/approve", auth, audit(models.EventPayoutApprove, id), h.Approve)
	g.POST("/:id/reject", auth, audit(models.EventPayoutReject, id), h.Reject)
	g.POST("/:id/confirm-funding", auth, audit(models.EventPayoutFund, id), h.ConfirmFunding)
	g.POST("/:id/pause", auth, audit(models.EventPayoutPause, id), h.Pause)
	g.POST("/:id/resume", auth, audit(models.EventPayoutResume, id), h.Resume)
	g.POST("/:id/cancel", auth, audit(models.EventPayoutCancel, id), h.Cancel)
	g.POST("/:id/retry-failed", auth, audit(models.EventPayoutRetry, id), h.RetryFailed)
	g.POST("/:id/items/:itemId/exclude", auth, audit(models.EventPayoutItemExclude, item), h.ExcludeItem)
	g.POST("/:id/items/:itemId/include", auth, audit(models.EventPayoutItemInclude, item), h.IncludeItem)
	g.POST("/:id/items/:itemId/mark-paid", auth, audit(models.EventPayoutItemPaid, item), h.MarkPaid)
}
