package middleware

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/gin-gonic/gin"
)

// RateLimitMiddleware returns a fixed-window rate limiter backed by Redis
// INCR+EXPIRE, so the limit holds across every instance behind the load
// balancer instead of being multiplied per-replica the way an in-memory
// counter would be. key namespaces this limiter's counters from any other
// endpoint's (e.g. "payment-history", "swap") - routes sharing a key share
// a budget.
//
// Callers are identified by the X-TW-SIGNER header when present (matches
// AuthenticationMiddlewareUsingTimestamp's signer-authenticated routes),
// else the X-TW-SERVICE-LINK-API-KEY header (matches
// AuthenticationMiddlewareUsingAPIKey's routes - keying by the caller's own
// key rather than IP, since a service link's calls may fan out through a
// shared egress IP), falling back to client IP when neither is present.
//
// defaultLimit is used unless overridden by RATE_LIMIT_<KEY>_PER_MINUTE or
// the global RATE_LIMIT_REQUESTS_PER_MINUTE env var (see
// rateLimitPerWindow). Set RATE_LIMIT_ENABLED=0 to disable rate limiting
// entirely. When Redis is disabled or unreachable, this middleware no-ops -
// same graceful-degradation posture as the rest of the caching layer -
// rather than failing requests closed.
func RateLimitMiddleware(gc *sharedconfig.GlobalConfig, key string, defaultLimit int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		if os.Getenv("RATE_LIMIT_ENABLED") == "0" {
			c.Next()
			return
		}
		if gc.RedisCache == nil || !gc.RedisCache.Enabled || gc.RedisCache.Client == nil {
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

		count, err := gc.RedisCache.Client.Incr(gc.RedisCache.Context, redisKey).Result()
		if err != nil {
			log.Printf("[RateLimitMiddleware] redis INCR failed for [%s], allowing request through: %v\n", redisKey, err)
			c.Next()
			return
		}
		if count == 1 {
			gc.RedisCache.Client.Expire(gc.RedisCache.Context, redisKey, window)
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
