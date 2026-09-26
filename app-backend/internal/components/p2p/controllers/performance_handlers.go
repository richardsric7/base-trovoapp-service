package p2p

import (
	"net/http"
	p2pServices "trovo-wallet-api/internal/components/p2p/services"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/gin-gonic/gin"
)

// getMerchantPerformanceHandler backs the marketplace trust signals on an
// offer card (Plan Section 8/26) - any authenticated caller can look up any
// merchant's performance, the same visibility the marketplace listing
// itself already gives every merchant's offers.
func getMerchantPerformanceHandler(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, err := currentUser(c, gc); err != nil {
			writeError(c, err)
			return
		}
		perf, err := p2pServices.GetMerchantPerformance(gc.DB, c.Param("merchantID"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error-temporary-server-error"})
			return
		}
		c.JSON(http.StatusOK, perf)
	}
}

// getMyPerformanceHandler returns the caller's own CustomerPerformance -
// shown to a merchant reviewing an order before/after acceptance (Plan
// Section 26).
func getMyPerformanceHandler(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := currentUser(c, gc)
		if err != nil {
			writeError(c, err)
			return
		}
		perf, err := p2pServices.GetCustomerPerformance(gc.DB, user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error-temporary-server-error"})
			return
		}
		c.JSON(http.StatusOK, perf)
	}
}
