// Package controllers registers the health endpoints.
package controllers

import (
	"admin-panel-dashboard/internal/components/health/handlers"
	"admin-panel-dashboard/internal/components/health/services"
	serverModels "admin-panel-dashboard/internal/server/models"

	"github.com/gin-gonic/gin"
)

// Init mounts the liveness and readiness probes.
//
// Both are registered at the server root AND under /api/v1 so that whichever
// convention a caller uses - the container healthcheck, uptime monitoring, or
// a load balancer - finds them. They are deliberately unauthenticated: a probe
// that needs a token cannot report on a service whose auth path is broken, and
// neither endpoint discloses anything beyond up/down and dependency names.
func Init(router *gin.Engine, s *serverModels.Server) {
	h := handlers.NewHandler(services.New(s.AdminDB, s.TrovoWalletDB, s.P2P, s.GC))

	router.GET("/health", h.Liveness)
	router.GET("/ready", h.Readiness)

	apiV1 := router.Group("/api/v1")
	apiV1.GET("/health", h.Liveness)
	apiV1.GET("/ready", h.Readiness)
}
