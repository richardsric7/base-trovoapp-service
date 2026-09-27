package p2p

import (
	"encoding/json"
	"io"
	"net/http"
	p2pServices "trovo-wallet-api/internal/components/p2p/services"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/gin-gonic/gin"
)

type paymentMethodRequest struct {
	PaymentChannel string `json:"paymentChannel"`
	Provider       string `json:"provider"`
	Account        string `json:"account"`
}

// postPaymentMethodsHandler godoc
// @Summary Add a new merchant payment method
// @Description Saves a new fiat settlement channel (e.g. bank transfer, mobile money) that the caller can later select from when creating SELL offers. Also backs the offer-creation form's inline "add new payment method" quick-create.
// @Tags P2P
// @Accept json
// @Produce json
// @Param body body paymentMethodRequest true "paymentChannel, provider, account"
// @Success 201 {object} map[string]interface{} "Created payment method"
// @Failure 400 {object} map[string]interface{}
// @Security SignatureAuth
// @Router /v1/p2p/payment-methods [post]
func postPaymentMethodsHandler(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := currentUser(c, gc)
		if err != nil {
			writeError(c, err)
			return
		}
		data, _ := io.ReadAll(c.Request.Body)
		var req paymentMethodRequest
		if err := json.Unmarshal(data, &req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "error-invalid-json"})
			return
		}
		pm, err := p2pServices.CreatePaymentMethod(gc, user.ID, user.Username, p2pServices.CreatePaymentMethodInput{
			PaymentChannel: req.PaymentChannel,
			Provider:       req.Provider,
			Account:        req.Account,
		})
		if err != nil {
			writeError(c, err)
			return
		}
		c.JSON(http.StatusCreated, pm)
	}
}

// getMyPaymentMethodsHandler godoc
// @Summary List my saved P2P payment methods
// @Tags P2P
// @Produce json
// @Success 200 {object} map[string]interface{} "data: list of payment methods"
// @Failure 500 {object} map[string]interface{}
// @Security SignatureAuth
// @Router /v1/p2p/payment-methods [get]
func getMyPaymentMethodsHandler(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := currentUser(c, gc)
		if err != nil {
			writeError(c, err)
			return
		}
		methods, err := p2pServices.ListMyPaymentMethods(gc.DB, user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error-temporary-server-error"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": methods})
	}
}

// putPaymentMethodHandler godoc
// @Summary Update a saved P2P payment method
// @Tags P2P
// @Accept json
// @Produce json
// @Param paymentMethodID path string true "Payment method ID"
// @Param body body paymentMethodRequest true "paymentChannel, provider, account"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Security SignatureAuth
// @Router /v1/p2p/payment-methods/{paymentMethodID} [put]
func putPaymentMethodHandler(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := currentUser(c, gc)
		if err != nil {
			writeError(c, err)
			return
		}
		data, _ := io.ReadAll(c.Request.Body)
		var req paymentMethodRequest
		if err := json.Unmarshal(data, &req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "error-invalid-json"})
			return
		}
		pm, err := p2pServices.UpdatePaymentMethod(gc, c.Param("paymentMethodID"), user.ID, p2pServices.UpdatePaymentMethodInput{
			PaymentChannel: req.PaymentChannel,
			Provider:       req.Provider,
			Account:        req.Account,
		})
		if err != nil {
			writeError(c, err)
			return
		}
		c.JSON(http.StatusOK, pm)
	}
}

type setPaymentMethodActiveRequest struct {
	Active bool `json:"active"`
}

// putPaymentMethodActiveHandler godoc
// @Summary Activate or deactivate a P2P payment method
// @Description Toggles a payment method active/inactive. Payment methods are never deleted, only deactivated. Fails if the payment method is still referenced by a live offer.
// putPaymentMethodActiveHandler toggles active/inactive - never a delete
// (Plan: "you cannot delete payment method"). Fails with 409 if the
// payment method is still in use by a live offer.
// @Tags P2P
// @Accept json
// @Produce json
// @Param paymentMethodID path string true "Payment method ID"
// @Param body body setPaymentMethodActiveRequest true "active: true or false"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{} "Payment method still in use by a live offer"
// @Security SignatureAuth
// @Router /v1/p2p/payment-methods/{paymentMethodID}/active [put]
func putPaymentMethodActiveHandler(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := currentUser(c, gc)
		if err != nil {
			writeError(c, err)
			return
		}
		data, _ := io.ReadAll(c.Request.Body)
		var req setPaymentMethodActiveRequest
		if err := json.Unmarshal(data, &req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "error-invalid-json"})
			return
		}
		pm, err := p2pServices.SetPaymentMethodActive(gc, c.Param("paymentMethodID"), user.ID, req.Active)
		if err != nil {
			writeError(c, err)
			return
		}
		c.JSON(http.StatusOK, pm)
	}
}
