package middleware

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

const corsAllowHeaders = "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-TW-PUBLIC-KEY, X-TW-SIGNATURE, X-TW-SIGNER, X-TW-TIMESTAMP, X-TW-DEVICE-ID, X-TW-APP-VERSION"

// CORSMiddleware controls which websites' pages may call app-backend from a
// browser. (Phones and servers are not subject to CORS.)
//
// CORS_ALLOWED_ORIGINS (comma-separated, e.g. https://wallet.trovo.app)
// limits it to those sites: their Origin is echoed back, any other gets no
// Access-Control-Allow-Origin and the browser blocks the response. Unset,
// any site may call it ("*"), as before. No cookies are used (requests are
// signed), so credentials are never allowed - "*" with credentials is
// invalid anyway.
func CORSMiddleware() gin.HandlerFunc {
	var allowed []string
	for _, o := range strings.Split(os.Getenv("CORS_ALLOWED_ORIGINS"), ",") {
		if o = strings.TrimRight(strings.TrimSpace(o), "/"); o != "" {
			allowed = append(allowed, o)
		}
	}
	return func(c *gin.Context) {
		if len(allowed) == 0 {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		} else {
			c.Writer.Header().Add("Vary", "Origin")
			origin := c.Request.Header.Get("Origin")
			for _, a := range allowed {
				if origin == a {
					c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
					break
				}
			}
		}
		c.Writer.Header().Set("Access-Control-Allow-Headers", corsAllowHeaders)
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
