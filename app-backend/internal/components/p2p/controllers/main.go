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
// prefix correction) and starts the two background sweeps this module
// needs in the absence of a Deposit Router event indexer / job scheduler:
// order-expiry and escrow-deposit reconciliation (Plan Sections 24, 33/42).
func Init(router *gin.Engine, gc *sharedconfig.GlobalConfig) {
	auth := middleware.AuthenticationMiddlewareUsingTimestamp()

	// Offers / marketplace
	router.POST("/v1/p2p/offers", auth, postOffersHandler(gc))
	router.GET("/v1/p2p/offers", auth, getMarketplaceOffersHandler(gc))
	router.GET("/v1/p2p/offers/:offerID", auth, getOfferHandler(gc))
	router.POST("/v1/p2p/offers/:offerID/activate", auth, postActivateOfferHandler(gc))
	router.POST("/v1/p2p/offers/:offerID/pause", auth, postPauseOfferHandler(gc))
	router.GET("/v1/p2p/my-offers", auth, getMyOffersHandler(gc))

	// Orders
	router.POST("/v1/p2p/orders", auth, postOrdersHandler(gc))
	router.GET("/v1/p2p/orders", auth, getMyOrdersHandler(gc))
	router.GET("/v1/p2p/orders/:orderID", auth, getOrderHandler(gc))
	router.POST("/v1/p2p/orders/:orderID/accept", auth, postAcceptOrderHandler(gc))
	router.POST("/v1/p2p/orders/:orderID/reject", auth, postRejectOrderHandler(gc))
	router.POST("/v1/p2p/orders/:orderID/cancel", auth, postCancelOrderHandler(gc))

	// Escrow deposit
	router.POST("/v1/p2p/orders/:orderID/escrow-deposit", auth, postEscrowDepositHandler(gc))

	// Fiat payment stage
	router.POST("/v1/p2p/orders/:orderID/payment-sent", auth, postPaymentSentHandler(gc))
	router.POST("/v1/p2p/orders/:orderID/payment-confirmed", auth, postPaymentConfirmedHandler(gc))

	// Disputes
	router.POST("/v1/p2p/orders/:orderID/disputes", auth, postOpenDisputeHandler(gc))
	router.POST("/v1/p2p/disputes/:disputeID/merchant-confirms-payment", auth, postMerchantConfirmsPaymentHandler(gc))
	router.POST("/v1/p2p/disputes/:disputeID/buyer-confirms-not-paid", auth, postBuyerConfirmsNotPaidHandler(gc))

	startBackgroundSweeps(gc)
}

// startBackgroundSweeps runs the order-expiry and escrow-deposit
// reconciliation sweeps on a fixed interval. There is no job-scheduler
// primitive in app-backend beyond plain ticker goroutines (the same pattern
// payments/controllers/main.go uses for callback retries).
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
}
