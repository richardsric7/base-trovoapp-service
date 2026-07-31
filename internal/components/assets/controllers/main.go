package assets

import (
	"trovo-wallet-api/internal/middleware"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/gin-gonic/gin"
)

// Init initializes /v2/assets endpoint

func Init(router *gin.Engine, gc *sharedconfig.GlobalConfig) {

	router.GET("/v1/curated-assets", middleware.AuthenticationMiddlewareUsingTimestamp(), getCuratedAssetsHandler(gc))
}
