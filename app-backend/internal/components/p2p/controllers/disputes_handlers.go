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
