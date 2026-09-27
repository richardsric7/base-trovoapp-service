package p2p

import (
	"encoding/json"
	"io"
	"net/http"
	p2pServices "trovo-wallet-api/internal/components/p2p/services"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/gin-gonic/gin"
)

type openDisputeRequest struct {
	Subject     string   `json:"subject"`
	Description string   `json:"description"`
	Evidence    []string `json:"evidence"`
}

// postOpenDisputeHandler godoc
// @Summary Open a dispute on a P2P order
// @Description Lets either party to an order raise a dispute (e.g. payment sent but not confirmed) with a subject, description, and supporting evidence (e.g. screenshot URLs), pausing the order's normal flow for arbitration.
// @Tags P2P
// @Accept json
// @Produce json
// @Param orderID path string true "Order ID"
// @Param body body openDisputeRequest true "subject, description, evidence (list of URLs)"
// @Success 201 {object} map[string]interface{} "Created dispute"
// @Failure 400 {object} map[string]interface{}
// @Security SignatureAuth
// @Router /v1/p2p/orders/{orderID}/disputes [post]
func postOpenDisputeHandler(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := currentUser(c, gc)
		if err != nil {
			writeError(c, err)
			return
		}
		data, _ := io.ReadAll(c.Request.Body)
		var req openDisputeRequest
		if err := json.Unmarshal(data, &req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "error-invalid-json"})
			return
		}
		dispute, err := p2pServices.OpenDispute(gc, c.Param("orderID"), user.ID, req.Subject, req.Description, req.Evidence)
		if err != nil {
			writeError(c, err)
			return
		}
		c.JSON(http.StatusCreated, dispute)
	}
}

// getOpenDisputeForOrderHandler godoc
// @Summary Get the open dispute for a P2P order
// @Tags P2P
// @Produce json
// @Param orderID path string true "Order ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{} "No open dispute for this order"
// @Security SignatureAuth
// @Router /v1/p2p/orders/{orderID}/dispute [get]
func getOpenDisputeForOrderHandler(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		dispute, err := p2pServices.GetOpenDisputeForOrder(gc, c.Param("orderID"))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "error-no-open-dispute"})
			return
		}
		c.JSON(http.StatusOK, dispute)
	}
}

// postMerchantConfirmsPaymentHandler godoc
// @Summary Merchant self-resolves a dispute by confirming payment was received
// @Description A same-side resolution shortcut: the merchant confirms the buyer's fiat payment actually did arrive, resolving the dispute in the buyer's favor without escalating to admin arbitration.
// @Tags P2P
// @Produce json
// @Param disputeID path string true "Dispute ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{} "Not the merchant, or invalid dispute state"
// @Security SignatureAuth
// @Router /v1/p2p/disputes/{disputeID}/merchant-confirms-payment [post]
func postMerchantConfirmsPaymentHandler(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := currentUser(c, gc)
		if err != nil {
			writeError(c, err)
			return
		}
		order, err := p2pServices.MerchantConfirmsPayment(gc, c.Param("disputeID"), user.ID)
		if err != nil {
			writeError(c, err)
			return
		}
		c.JSON(http.StatusOK, order)
	}
}

// postBuyerConfirmsNotPaidHandler godoc
// @Summary Buyer self-resolves a dispute by confirming they did not pay
// @Description A same-side resolution shortcut: the buyer confirms they never actually sent the fiat payment, resolving the dispute in the merchant's favor without escalating to admin arbitration.
// @Tags P2P
// @Produce json
// @Param disputeID path string true "Dispute ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{} "Not the customer, or invalid dispute state"
// @Security SignatureAuth
// @Router /v1/p2p/disputes/{disputeID}/buyer-confirms-not-paid [post]
func postBuyerConfirmsNotPaidHandler(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := currentUser(c, gc)
		if err != nil {
			writeError(c, err)
			return
		}
		order, err := p2pServices.BuyerConfirmsPaymentNotMade(gc, c.Param("disputeID"), user.ID)
		if err != nil {
			writeError(c, err)
			return
		}
		c.JSON(http.StatusOK, order)
	}
}
