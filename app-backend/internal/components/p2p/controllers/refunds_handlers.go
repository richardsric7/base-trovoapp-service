package p2p

import (
	"net/http"
	p2pServices "trovo-wallet-api/internal/components/p2p/services"
	"trovo-wallet-api/internal/middleware"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/gin-gonic/gin"
)

// getMyRefundsHandler godoc
// @Summary List refunds owed to my wallet
// @Description Lists every P2P escrow refund owed to the caller's wallet address, e.g. from a cancelled or disputed order where the deposited asset needs to be returned.
// @Tags P2P
// @Produce json
// @Success 200 {object} map[string]interface{} "data: list of refunds"
// @Failure 500 {object} map[string]interface{}
// @Security SignatureAuth
// @Router /v1/p2p/refunds [get]
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

// postClaimRefundHandler godoc
// @Summary Claim an owed P2P escrow refund
// @Description Triggers the on-chain transfer of a pending refund back to the caller's wallet.
// @Tags P2P
// @Produce json
// @Param refundID path string true "Refund ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Security SignatureAuth
// @Router /v1/p2p/refunds/{refundID}/claim [post]
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
