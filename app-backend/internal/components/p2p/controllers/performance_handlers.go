package p2p

import (
	"net/http"
	p2pServices "trovo-wallet-api/internal/components/p2p/services"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/gin-gonic/gin"
)

// getMerchantPerformanceHandler godoc
// @Summary Get a P2P merchant's trust/performance stats
// @Description Backs the trust signals shown on a marketplace offer card (e.g. completion rate, average release time). Any authenticated caller can look up any merchant's performance, the same visibility the marketplace listing itself already gives.
// getMerchantPerformanceHandler backs the marketplace trust signals on an
// offer card (Plan Section 8/26) - any authenticated caller can look up any
// merchant's performance, the same visibility the marketplace listing
// itself already gives every merchant's offers.
// @Tags P2P
// @Produce json
// @Param merchantID path string true "Merchant's user ID"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security SignatureAuth
// @Router /v1/p2p/merchants/{merchantID}/performance [get]
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

// getOrderCustomerPerformanceHandler godoc
// @Summary Get a customer's trading performance for one of my orders
// @Description Lets a merchant look up the performance/trust history of the customer on a specific order they're party to, before deciding to accept it. Unlike merchant performance (public via the marketplace), a customer's history is private, so this only works if the caller is the merchant on orderID - it never accepts an arbitrary customer ID.
// getOrderCustomerPerformanceHandler lets a merchant look up the
// performance of the specific customer they have an order with, before
// deciding to accept it (Plan Section 8/26) - unlike merchant performance
// (already public via the marketplace listing), a customer's trading
// history is private, so this is scoped to a real order relationship
// rather than taking an arbitrary customer id: the caller must be the
// merchant on orderID, and the customer looked up is that order's own
// customer, never a caller-supplied id.
// @Tags P2P
// @Produce json
// @Param orderID path string true "Order ID"
// @Success 200 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{} "Not the merchant on this order"
// @Failure 404 {object} map[string]interface{} "Order not found"
// @Security SignatureAuth
// @Router /v1/p2p/orders/{orderID}/customer-performance [get]
func getOrderCustomerPerformanceHandler(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
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
		if order.MerchantUserID != user.ID {
			c.JSON(http.StatusForbidden, gin.H{"error": "error-forbidden", "message": "You are not the merchant on this order"})
			return
		}
		perf, err := p2pServices.GetCustomerPerformance(gc.DB, order.CustomerUserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error-temporary-server-error"})
			return
		}
		c.JSON(http.StatusOK, perf)
	}
}

// getMyPerformanceHandler godoc
// @Summary Get my own P2P trading performance
// @Description Returns the caller's own customer performance/trust stats - the same data a merchant sees when reviewing an order before or after acceptance.
// @Tags P2P
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security SignatureAuth
// @Router /v1/p2p/my-performance [get]
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
