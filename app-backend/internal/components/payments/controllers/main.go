package payments

import (
	"time"
	userModels "trovo-wallet-api/internal/components/users/models"
	"trovo-wallet-api/internal/sharedconfig"

	"trovo-wallet-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

// Init initializes the controller

func Init(router *gin.Engine, callBackRetryChan chan userModels.RetryCallbacks, gc *sharedconfig.GlobalConfig) {

	// payment notifications are delivered (and retried) through
	// gc.SendCallback; nothing is sent on callBackRetryChan any more

	router.POST("/v1/users/payment", middleware.AuthenticationMiddlewareUsingTimestamp(), middleware.RateLimitMiddleware(gc, "payment", 20, time.Minute), postUsersPaymentHandler(callBackRetryChan, gc))
	router.POST("/v1/shared-access/payment", middleware.AuthenticationMiddlewareUsingTimestamp(), middleware.RateLimitMiddleware(gc, "payment", 20, time.Minute), postSharedAccessPaymentHandler(callBackRetryChan, gc))
}
