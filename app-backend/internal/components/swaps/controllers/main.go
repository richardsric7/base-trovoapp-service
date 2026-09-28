package payments

import (
	"time"
	"trovo-wallet-api/internal/middleware"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/gin-gonic/gin"
)

// Init initializes the controller

func Init(router *gin.Engine, gc *sharedconfig.GlobalConfig) {

	router.POST("/v1/users/swap", middleware.AuthenticationMiddlewareUsingTimestamp(), middleware.RateLimitMiddleware(gc, "swap", 10, time.Minute), postUsersSwapHandler(gc))
	router.POST("/v1/shared-access/swap", middleware.AuthenticationMiddlewareUsingTimestamp(), middleware.RateLimitMiddleware(gc, "swap", 10, time.Minute), postSharedAccessSwapHandler(gc))
}
