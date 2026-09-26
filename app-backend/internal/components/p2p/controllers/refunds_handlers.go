package p2p

import (
	"net/http"
	p2pServices "trovo-wallet-api/internal/components/p2p/services"
	"trovo-wallet-api/internal/middleware"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/gin-gonic/gin"
)

// getMyRefundsHandler lists every refund owed to the caller's wallet
// address (Plan Section 45 - Refund.Sender is a wallet address, the same
// identifier the escrow-deposit endpoint already authenticates by).
func getMyRefundsHandler(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		refunds, err := p2pServices.ListMyRefunds(gc.DB, middleware.ExtractAddress(c))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error-temporary-server-error"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": refunds})
	}
}

// postClaimRefundHandler implements Plan Section 45's claim sequence.
func postClaimRefundHandler(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		refund, err := p2pServices.ClaimRefund(gc, c.Param("refundID"), middleware.ExtractAddress(c))
		if err != nil {
			writeError(c, err)
			return
		}
		c.JSON(http.StatusOK, refund)
	}
}
