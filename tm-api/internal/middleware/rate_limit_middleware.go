package middleware

import (
	"admin-panel-dashboard/internal/models"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimitMiddleware returns a fixed-window rate limiter backed by Redis
// INCR+EXPIRE, so the limit holds across every instance behind the load
// balancer instead of being multiplied per-replica the way an in-memory
// counter would be. key namespaces this limiter's counters from any other
// caller's (routes sharing a key share a budget) - the global default
// registered in main.go uses "global".
//
// Callers are identified by the X-TW-SIGNER header when present, else the
// X-TW-SERVICE-LINK-API-KEY header, falling back to client IP. Most tm-api
// routes are behind JWT auth rather than these headers, so IP is the
// common case here - still useful as blanket abuse insurance for an admin
// panel with a much smaller abuse surface than the wallet API.
//
// defaultLimit is used unless overridden by RATE_LIMIT_<KEY>_PER_MINUTE or
// the global RATE_LIMIT_REQUESTS_PER_MINUTE env var (see
// rateLimitPerWindow). Set RATE_LIMIT_ENABLED=0 to disable rate limiting
// entirely. When Redis is disabled or unreachable, this middleware no-ops
// rather than failing requests closed.
func RateLimitMiddleware(gc *models.GlobalConfig, key string, defaultLimit int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		if os.Getenv("RATE_LIMIT_ENABLED") == "0" {
			c.Next()
			return
		}
		if gc.Cache == nil || !gc.Cache.Enabled || gc.Cache.Client == nil {
			c.Next()
			return
		}

		identity := ExtractSigner(c)
		if identity == "" {
			identity = ExtractServiceLinkApiKey(c)
		}
		if identity == "" {
			identity = c.ClientIP()
		}

		redisKey := "ratelimit:" + key + ":" + identity
		limit := rateLimitPerWindow(key, defaultLimit)

		count, err := gc.Cache.Client.Incr(gc.Cache.Context, redisKey).Result()
		if err != nil {
			log.Printf("[RateLimitMiddleware] redis INCR failed for [%s], allowing request through: %v\n", redisKey, err)
			c.Next()
			return
		}
		if count == 1 {
			gc.Cache.Client.Expire(gc.Cache.Context, redisKey, window)
		}

		if count > int64(limit) {
			c.Header("Retry-After", strconv.Itoa(int(window.Seconds())))
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "rate-limit-exceeded",
				"message": "Too many requests, please try again later.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// rateLimitPerWindow resolves the effective limit for key: an explicit
// RATE_LIMIT_<KEY>_PER_MINUTE override (key upper-cased, hyphens to
// underscores) takes precedence, then the global
// RATE_LIMIT_REQUESTS_PER_MINUTE, then defaultLimit.
func rateLimitPerWindow(key string, defaultLimit int) int {
	envKey := "RATE_LIMIT_" + strings.ToUpper(strings.ReplaceAll(key, "-", "_")) + "_PER_MINUTE"
	if n, ok := positiveEnvInt(envKey); ok {
		return n
	}
	if n, ok := positiveEnvInt("RATE_LIMIT_REQUESTS_PER_MINUTE"); ok {
		return n
	}
	return defaultLimit
}

func positiveEnvInt(envKey string) (int, bool) {
	raw := os.Getenv(envKey)
	if raw == "" {
		return 0, false
	}
	n, err := strconv.Atoi(raw)
	if err != nil || n <= 0 {
		return 0, false
	}
	return n, true
}
