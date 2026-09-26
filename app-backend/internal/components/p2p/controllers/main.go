package p2p

import (
	"log"
	"time"
	p2pServices "trovo-wallet-api/internal/components/p2p/services"
	"trovo-wallet-api/internal/middleware"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/gin-gonic/gin"
)

// Init registers every P2P route under /v1/p2p/... (Plan Section 12's
// prefix correction) and starts the three background sweeps this module
// needs in the absence of a Deposit Router event indexer / job scheduler:
// order-expiry, escrow-deposit reconciliation, and marketplace ranking
// (Plan Sections 24, 33/42, 14).
func Init(router *gin.Engine, gc *sharedconfig.GlobalConfig) {
	auth := middleware.AuthenticationMiddlewareUsingTimestamp()

	// Offers / marketplace
	router.POST("/v1/p2p/offers", auth, postOffersHandler(gc))
	router.GET("/v1/p2p/offers", auth, getMarketplaceOffersHandler(gc))
	router.GET("/v1/p2p/offers/:offerID", auth, getOfferHandler(gc))
	router.GET("/v1/p2p/offers/:offerID/quote", auth, getOfferQuoteHandler(gc))
	router.PUT("/v1/p2p/offers/:offerID", auth, putOfferHandler(gc))
	router.POST("/v1/p2p/offers/:offerID/activate", auth, postActivateOfferHandler(gc))
	router.POST("/v1/p2p/offers/:offerID/pause", auth, postPauseOfferHandler(gc))
	router.POST("/v1/p2p/offers/:offerID/close", auth, postCloseOfferHandler(gc))
	router.GET("/v1/p2p/my-offers", auth, getMyOffersHandler(gc))

	// Orders
	router.POST("/v1/p2p/orders", auth, postOrdersHandler(gc))
	router.GET("/v1/p2p/orders", auth, getMyOrdersHandler(gc))
	router.GET("/v1/p2p/orders/:orderID", auth, getOrderHandler(gc))
	router.POST("/v1/p2p/orders/:orderID/accept", auth, postAcceptOrderHandler(gc))
	router.POST("/v1/p2p/orders/:orderID/reject", auth, postRejectOrderHandler(gc))
	router.POST("/v1/p2p/orders/:orderID/cancel", auth, postCancelOrderHandler(gc))
	router.POST("/v1/p2p/orders/:orderID/merchant-cancel", auth, postMerchantCancelOrderHandler(gc))

	// Escrow deposit
	router.POST("/v1/p2p/orders/:orderID/escrow-deposit", auth, postEscrowDepositHandler(gc))
	router.POST("/v1/p2p/orders/:orderID/escrow-deposit/regenerate-shortlink", auth, postRegenerateEscrowShortlinkHandler(gc))

	// Fiat payment stage
	router.POST("/v1/p2p/orders/:orderID/payment-sent", auth, postPaymentSentHandler(gc))
	router.POST("/v1/p2p/orders/:orderID/payment-confirmed", auth, postPaymentConfirmedHandler(gc))

	// Disputes
	router.POST("/v1/p2p/orders/:orderID/disputes", auth, postOpenDisputeHandler(gc))
	router.GET("/v1/p2p/orders/:orderID/dispute", auth, getOpenDisputeForOrderHandler(gc))
	router.POST("/v1/p2p/disputes/:disputeID/merchant-confirms-payment", auth, postMerchantConfirmsPaymentHandler(gc))
	router.POST("/v1/p2p/disputes/:disputeID/buyer-confirms-not-paid", auth, postBuyerConfirmsNotPaidHandler(gc))

	// Admin/arbiter dispute resolution - app-backend has no admin/staff
	// user model of its own (that lives in tm-api); this is a server-to-
	// server endpoint an admin-authorized tm-api call can reach, the same
	// API-key trust boundary as every other app-backend<->tm-api
	// integration point (internal/components/servicelinks).
	router.POST("/v1/p2p/disputes/:disputeID/admin-resolve", middleware.AuthenticationMiddlewareUsingAPIKey(gc), postAdminResolveDisputeHandler(gc))

	// Refunds
	router.GET("/v1/p2p/refunds", auth, getMyRefundsHandler(gc))
	router.POST("/v1/p2p/refunds/:refundID/claim", auth, postClaimRefundHandler(gc))

	// Performance / trust signals
	router.GET("/v1/p2p/merchants/:merchantID/performance", auth, getMerchantPerformanceHandler(gc))
	router.GET("/v1/p2p/orders/:orderID/customer-performance", auth, getOrderCustomerPerformanceHandler(gc))
	router.GET("/v1/p2p/my-performance", auth, getMyPerformanceHandler(gc))

	startBackgroundSweeps(gc)
}

// startBackgroundSweeps runs the order-expiry, escrow-deposit
// reconciliation, and marketplace ranking sweeps on a fixed interval. There
// is no job-scheduler primitive in app-backend beyond plain ticker
// goroutines (the same pattern payments/controllers/main.go uses for
// callback retries).
func startBackgroundSweeps(gc *sharedconfig.GlobalConfig) {
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			if n, err := p2pServices.ExpireStaleOrders(gc); err != nil {
				log.Printf("[p2p:ExpireStaleOrders] error: %v\n", err)
			} else if n > 0 {
				log.Printf("[p2p:ExpireStaleOrders] expired %d stale order(s)\n", n)
			}
		}
	}()
	go func() {
		ticker := time.NewTicker(2 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			p2pServices.RunEscrowReconciliationSweep(gc)
		}
	}()
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			if err := p2pServices.CalculateRanking(gc); err != nil {
				log.Printf("[p2p:CalculateRanking] error: %v\n", err)
			}
		}
	}()
}
