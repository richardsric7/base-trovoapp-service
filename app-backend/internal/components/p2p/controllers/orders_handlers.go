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
