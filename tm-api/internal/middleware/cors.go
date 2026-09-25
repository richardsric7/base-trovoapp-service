package middleware

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// allowedOrigins reads CORS_ALLOWED_ORIGINS (comma-separated). Empty/unset means no explicit
// allow-list is configured — see CORSMiddleware for what that means for origin matching.
func allowedOrigins() []string {
	raw := os.Getenv("CORS_ALLOWED_ORIGINS")

	var origins []string
	for _, origin := range strings.Split(raw, ",") {
		origin = strings.TrimSpace(origin)
		if len(origin) > 0 {
			origins = append(origins, origin)
		}
	}
	return origins
}

// CORSMiddleware echoes back the request's Origin header instead of literal "*" — "*" combined
// with Allow-Credentials is an invalid combination per the CORS spec that compliant browsers
// reject outright for any credentialed request. When CORS_ALLOWED_ORIGINS is set, only origins in
// that list are echoed back (opt-in restriction); when it's unset, every origin is echoed,
// preserving the previous "any origin" behavior while still never emitting the invalid "*"
// alongside credentials. Deliberately not defaulting the allow-list to
// TROVO_MANAGER_FRONTEND_URL alone: that single value doesn't cover every real deployed frontend
// origin (the organisation portal in particular), and defaulting to it silently CORS-blocked
// every other origin, including org member logins.
func CORSMiddleware() gin.HandlerFunc {
	origins := allowedOrigins()

	return func(c *gin.Context) {
		reqOrigin := c.Request.Header.Get("Origin")
		if len(reqOrigin) > 0 {
			if len(origins) == 0 {
				c.Writer.Header().Set("Access-Control-Allow-Origin", reqOrigin)
			} else {
				for _, allowed := range origins {
					if reqOrigin == allowed {
						c.Writer.Header().Set("Access-Control-Allow-Origin", reqOrigin)
						break
					}
				}
			}
		}
		c.Writer.Header().Set("Vary", "Origin")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE,PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
