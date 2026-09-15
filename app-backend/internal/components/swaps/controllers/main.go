package payments

import (
	"trovo-wallet-api/internal/middleware"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/gin-gonic/gin"
)

// Init initializes the controller

func Init(router *gin.Engine, gc *sharedconfig.GlobalConfig) {

	router.POST("/v1/users/swap", middleware.AuthenticationMiddlewareUsingTimestamp(), postUsersSwapHandler(gc))
	router.POST("/v1/shared-access/swap", middleware.AuthenticationMiddlewareUsingTimestamp(), postSharedAccessSwapHandler(gc))
}
