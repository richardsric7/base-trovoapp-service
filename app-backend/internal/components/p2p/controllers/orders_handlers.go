package p2p

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	p2pServices "trovo-wallet-api/internal/components/p2p/services"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/gin-gonic/gin"
)

type createOrderRequest struct {
	OfferID              string `json:"offerId"`
	SpecifiedAssetAmount string `json:"specifiedAssetAmount"`
}

// postOrdersHandler godoc
// @Summary Create a new P2P order against an offer
// @Description Creates an order (a trade) against an existing offer for a specified asset amount. Starts the order lifecycle (pending merchant acceptance, then escrow deposit, then fiat payment, then release).
// @Tags P2P
// @Accept json
// @Produce json
// @Param body body createOrderRequest true "offerId and specifiedAssetAmount"
// @Success 201 {object} map[string]interface{} "Created order"
// @Failure 400 {object} map[string]interface{} "Invalid JSON or validation error"
// @Security SignatureAuth
// @Router /v1/p2p/orders [post]
func postOrdersHandler(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := currentUser(c, gc)
		if err != nil {
			writeError(c, err)
			return
		}
		data, _ := io.ReadAll(c.Request.Body)
		var req createOrderRequest
		if err := json.Unmarshal(data, &req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "error-invalid-json"})
			return
		}
		order, err := p2pServices.CreateOrder(gc, user.ID, user.Username, p2pServices.CreateOrderInput{
			OfferID:              req.OfferID,
			SpecifiedAssetAmount: req.SpecifiedAssetAmount,
		})
		if err != nil {
			writeError(c, err)
			return
		}
		c.JSON(http.StatusCreated, order)
	}
}

// getMyOrdersHandler godoc
// @Summary List my P2P orders
// @Description Lists orders the caller is party to, either as customer or merchant, with optional filters and pagination.
// @Tags P2P
// @Produce json
// @Param role query string false "Filter by role: customer or merchant"
// @Param status query string false "Filter by order status"
// @Param page query int false "Page number (default 1)"
// @Param pageSize query int false "Results per page (default 20)"
// @Success 200 {object} map[string]interface{} "data (orders), total, page, pageSize"
// @Failure 500 {object} map[string]interface{}
// @Security SignatureAuth
// @Router /v1/p2p/orders [get]
func getMyOrdersHandler(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := currentUser(c, gc)
		if err != nil {
			writeError(c, err)
			return
		}
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
		orders, total, err := p2pServices.ListOrders(gc.DB, p2pServices.ListOrdersFilter{
			UserID:   user.ID,
			Role:     c.Query("role"),
			Status:   c.Query("status"),
			Page:     page,
			PageSize: pageSize,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error-temporary-server-error"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": orders, "total": total, "page": page, "pageSize": pageSize})
	}
}

// getOrderHandler godoc
// @Summary Get a single P2P order by ID
// @Description Returns order details. Only the order's customer or merchant may view it. Also best-effort reconciles any escrow deposit made via the shared deposit link since the last check.
// @Tags P2P
// @Produce json
// @Param orderID path string true "Order ID"
// @Success 200 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{} "Caller is not a party to this order"
// @Failure 404 {object} map[string]interface{} "Order not found"
// @Security SignatureAuth
// @Router /v1/p2p/orders/{orderID} [get]
func getOrderHandler(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := currentUser(c, gc)
		if err != nil {
			writeError(c, err)
			return
		}
		order, err := p2pServices.GetOrderByID(gc.DB, c.Param("orderID"))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "error-order-not-found"})
			return
		}
		if order.CustomerUserID != user.ID && order.MerchantUserID != user.ID {
			c.JSON(http.StatusForbidden, gin.H{"error": "error-forbidden"})
			return
		}
		// Best-effort reconciliation on read, so a third-party deposit made
		// via the shared link (Plan Section 33) is reflected without the
		// caller having to wait for the next background sweep.
		_ = p2pServices.ReconcileEscrowDepositsFromPaymentHistory(gc, &order)
		c.JSON(http.StatusOK, order)
	}
}

// postAcceptOrderHandler godoc
// @Summary Merchant accepts a pending P2P order
// @Tags P2P
// @Produce json
// @Param orderID path string true "Order ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{} "Not the merchant, or invalid order state"
// @Security SignatureAuth
// @Router /v1/p2p/orders/{orderID}/accept [post]
func postAcceptOrderHandler(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := currentUser(c, gc)
		if err != nil {
			writeError(c, err)
			return
		}
		order, err := p2pServices.AcceptOrder(gc, c.Param("orderID"), user.ID)
		if err != nil {
			writeError(c, err)
			return
		}
		c.JSON(http.StatusOK, order)
	}
}

// postRejectOrderHandler godoc
// @Summary Merchant rejects a pending P2P order
// @Tags P2P
// @Produce json
// @Param orderID path string true "Order ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{} "Not the merchant, or invalid order state"
// @Security SignatureAuth
// @Router /v1/p2p/orders/{orderID}/reject [post]
func postRejectOrderHandler(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := currentUser(c, gc)
		if err != nil {
			writeError(c, err)
			return
		}
		order, err := p2pServices.RejectOrder(gc, c.Param("orderID"), user.ID)
		if err != nil {
			writeError(c, err)
			return
		}
		c.JSON(http.StatusOK, order)
	}
}

// postCancelOrderHandler godoc
// @Summary Customer cancels their own P2P order
// @Tags P2P
// @Produce json
// @Param orderID path string true "Order ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{} "Not the customer, or invalid order state"
// @Security SignatureAuth
// @Router /v1/p2p/orders/{orderID}/cancel [post]
func postCancelOrderHandler(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := currentUser(c, gc)
		if err != nil {
			writeError(c, err)
			return
		}
		order, err := p2pServices.CancelOrder(gc, c.Param("orderID"), user.ID)
		if err != nil {
			writeError(c, err)
			return
		}
		c.JSON(http.StatusOK, order)
	}
}

// postMerchantCancelOrderHandler godoc
// @Summary Merchant cancels a P2P order
// @Tags P2P
// @Produce json
// @Param orderID path string true "Order ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{} "Not the merchant, or invalid order state"
// @Security SignatureAuth
// @Router /v1/p2p/orders/{orderID}/merchant-cancel [post]
func postMerchantCancelOrderHandler(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := currentUser(c, gc)
		if err != nil {
			writeError(c, err)
			return
		}
		order, err := p2pServices.MerchantCancelOrder(gc, c.Param("orderID"), user.ID)
		if err != nil {
			writeError(c, err)
			return
		}
		c.JSON(http.StatusOK, order)
	}
}

// postPaymentSentHandler godoc
// @Summary Customer marks fiat payment as sent
// @Description Called by the buying customer after they have sent the off-platform fiat payment to the merchant's payment method, to notify the merchant to check and confirm receipt.
// @Tags P2P
// @Produce json
// @Param orderID path string true "Order ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{} "Not the customer, or invalid order state"
// @Security SignatureAuth
// @Router /v1/p2p/orders/{orderID}/payment-sent [post]
func postPaymentSentHandler(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := currentUser(c, gc)
		if err != nil {
			writeError(c, err)
			return
		}
		order, err := p2pServices.MarkFiatPaymentSent(gc, c.Param("orderID"), user.ID)
		if err != nil {
			writeError(c, err)
			return
		}
		c.JSON(http.StatusOK, order)
	}
}

// postPaymentConfirmedHandler godoc
// @Summary Merchant confirms fiat payment received
// @Description Called by the merchant after they have verified the buyer's fiat payment landed. Triggers release of the escrowed asset to the customer.
// @Tags P2P
// @Produce json
// @Param orderID path string true "Order ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{} "Not the merchant, or invalid order state"
// @Security SignatureAuth
// @Router /v1/p2p/orders/{orderID}/payment-confirmed [post]
func postPaymentConfirmedHandler(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := currentUser(c, gc)
		if err != nil {
			writeError(c, err)
			return
		}
		order, err := p2pServices.ConfirmFiatPaymentReceived(gc, c.Param("orderID"), user.ID)
		if err != nil {
			writeError(c, err)
			return
		}
		c.JSON(http.StatusOK, order)
	}
}
