// Package health serves the liveness and readiness checks a hosting
// platform (Docker, Kubernetes, a load balancer) polls.
package health

import (
	"context"
	"net/http"
	"time"

	"trovo-wallet-api/internal/sharedconfig"

	"github.com/gin-gonic/gin"
)

// checkTimeout bounds each dependency check, so a hung dependency makes
// /ready answer 503 quickly instead of hanging the probe.
const checkTimeout = 3 * time.Second

// Init registers GET /health and GET /ready. Neither needs authentication
// and neither reveals more than up/down per dependency.
func Init(router *gin.Engine, gc *sharedconfig.GlobalConfig) {
	router.GET("/health", getHealth)
	router.GET("/ready", getReady(gc))
}

// getHealth godoc
// @Summary Liveness check
// @Description 200 while the process is running. Use it as the liveness check.
// @Tags Health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func getHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// getReady godoc
// @Summary Readiness check
// @Description Checks the main database, the tracking database, Redis (when caching is on) and the Base RPC node. 200 when all are up, 503 otherwise.
// @Tags Health
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 503 {object} map[string]interface{}
// @Router /ready [get]
func getReady(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		checks := Check(c.Request.Context(), gc)
		status, code := "ready", http.StatusOK
		for _, v := range checks {
			if v != "up" {
				status, code = "not ready", http.StatusServiceUnavailable
				break
			}
		}
		c.JSON(code, gin.H{"status": status, "checks": checks})
	}
}

// Check reports "up" or "down" for each dependency app-backend needs to
// serve requests. Redis is only checked when caching is on.
func Check(ctx context.Context, gc *sharedconfig.GlobalConfig) map[string]string {
	checks := map[string]string{
		"database":          pingDB(ctx, gc, false),
		"tracking_database": pingDB(ctx, gc, true),
		"blockchain":        "down",
	}
	if gc.RedisCache != nil && gc.RedisCache.Enabled && gc.RedisCache.Client != nil {
		cctx, cancel := context.WithTimeout(ctx, checkTimeout)
		checks["redis"] = upDown(gc.RedisCache.Client.Ping(cctx).Err())
		cancel()
	}
	if gc.BantuExpansionClient != nil {
		cctx, cancel := context.WithTimeout(ctx, checkTimeout)
		_, err := gc.BantuExpansionClient.BlockNumber(cctx)
		cancel()
		checks["blockchain"] = upDown(err)
	}
	return checks
}

func pingDB(ctx context.Context, gc *sharedconfig.GlobalConfig, tracking bool) string {
	db := gc.DB
	if tracking {
		db = gc.RoachDB
	}
	if db == nil {
		return "down"
	}
	sqlDB, err := db.DB()
	if err != nil {
		return "down"
	}
	cctx, cancel := context.WithTimeout(ctx, checkTimeout)
	defer cancel()
	return upDown(sqlDB.PingContext(cctx))
}

func upDown(err error) string {
	if err != nil {
		return "down"
	}
	return "up"
}
