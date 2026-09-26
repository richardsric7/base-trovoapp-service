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

// postPaymentMethodsHandler saves a new merchant payment method - also the
// backing call for the offer-creation form's inline "Add new payment
// method" quick-create (the first item in its payment-method dropdown).
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

// putPaymentMethodActiveHandler toggles active/inactive - never a delete
// (Plan: "you cannot delete payment method"). Fails with 409 if the
// payment method is still in use by a live offer.
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
