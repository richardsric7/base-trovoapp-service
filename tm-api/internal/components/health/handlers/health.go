// Package handlers exposes the liveness and readiness probes over HTTP.
package handlers

import (
	"net/http"

	"admin-panel-dashboard/internal/components/health/services"

	"github.com/gin-gonic/gin"
)

// Handler serves the health endpoints.
type Handler struct {
	service *services.Service
}

// NewHandler builds the health handler.
func NewHandler(s *services.Service) *Handler {
	return &Handler{service: s}
}

// Liveness godoc
//
//	@Summary		Liveness probe
//	@Description	Reports that the API process is running. Performs no dependency
//	@Description	checks, so it stays 200 even while a database is unreachable -
//	@Description	use /ready to decide whether the service can serve requests.
//	@Description	Intended for uptime monitoring. No authentication required.
//	@Tags			Health
//	@Produce		json
//	@Success		200	{object}	services.LivenessResponse
//	@Router			/health [get]
func (h *Handler) Liveness(c *gin.Context) {
	c.JSON(http.StatusOK, h.service.Liveness())
}

// Readiness godoc
//
//	@Summary		Readiness probe
//	@Description	Checks every backing dependency (admin, wallet and P2P
//	@Description	databases, Redis, the Trovo Wallet API) and reports each one
//	@Description	individually. Returns 503 when a required dependency is down,
//	@Description	so an unhealthy container is never given traffic. Optional
//	@Description	dependencies (Redis, wallet API) being down report the service
//	@Description	as "degraded" but still return 200. No authentication required.
//	@Tags			Health
//	@Produce		json
//	@Success		200	{object}	services.ReadinessResponse	"Ready (status up or degraded)"
//	@Failure		503	{object}	services.ReadinessResponse	"A required dependency is down"
//	@Router			/ready [get]
func (h *Handler) Readiness(c *gin.Context) {
	report, status := h.service.Readiness(c.Request.Context())
	c.JSON(status, report)
}
