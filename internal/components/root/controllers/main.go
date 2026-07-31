package root

import (
	"trovo-wallet-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

func Init(router *gin.Engine) {

	//Returns organisation running this bantupay api instance

	//Returns organisation running this bantupay api instance

	// Serve Android App Links file

	router.GET("/", middleware.AuthenticationMiddlewareUsingTimestamp(), getHandler())
	router.GET("/.well-known/stellar.toml", middleware.AuthenticationMiddlewareUsingTimestamp(), getDotwellKnownStellarDottomlHandler())
	router.GET("/.well-known/apple-app-site-association", getDotwellKnownAppleAppSiteAssociationHandler())
	router.GET("/.well-known/assetlinks.json", getDotwellKnownAssetlinksDotjsonHandler())
}
