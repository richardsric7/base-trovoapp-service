package root

import (
	root "trovo-wallet-api/internal/components/root/services"
	"trovo-wallet-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

func Init(router *gin.Engine) {

	//Returns organisation running this bantupay api instance
	router.GET("/", middleware.AuthenticationMiddlewareUsingTimestamp(), func(c *gin.Context) {

		rootInfo := root.GetRootInfo()
		rootInfo.PublicKey = middleware.ExtractPublicKey(c)
		c.JSON(200, rootInfo)
	})
	//Returns organisation running this bantupay api instance
	router.GET("/.well-known/stellar.toml", middleware.AuthenticationMiddlewareUsingTimestamp(), func(c *gin.Context) {

		rootInfo := root.GetRootInfo()
		rootInfo.PublicKey = middleware.ExtractPublicKey(c)
		c.JSON(200, rootInfo)
	})
}
