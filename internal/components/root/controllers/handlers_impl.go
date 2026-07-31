package root

import (
	"net/http"
	"os"
	root "trovo-wallet-api/internal/components/root/services"
	"trovo-wallet-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

// getHandler godoc
// @Summary GET /
// @Tags root
// @Router / [get]
func getHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		rootInfo := root.GetRootInfo()
		rootInfo.PublicKey = middleware.ExtractPublicKey(c)
		c.JSON(200, rootInfo)
	}
}


// getDotwellKnownStellarDottomlHandler godoc
// @Summary GET /.well-known/stellar.toml
// @Tags root
// @Router /.well-known/stellar.toml [get]
func getDotwellKnownStellarDottomlHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		rootInfo := root.GetRootInfo()
		rootInfo.PublicKey = middleware.ExtractPublicKey(c)
		c.JSON(200, rootInfo)
	}
}


// getDotwellKnownAppleAppSiteAssociationHandler godoc
// @Summary GET /.well-known/apple-app-site-association
// @Tags root
// @Router /.well-known/apple-app-site-association [get]
func getDotwellKnownAppleAppSiteAssociationHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusOK, gin.H{
			"applinks": gin.H{
				"apps": []string{},
				"details": []gin.H{
					{
						"appID": os.Getenv("DYNAMIC_LINKS_IOS_BUNDLE_ID"),
						"paths": []string{"*"},
					},
				},
			},
		})
	}
}


// getDotwellKnownAssetlinksDotjsonHandler godoc
// @Summary GET /.well-known/assetlinks.json
// @Tags root
// @Router /.well-known/assetlinks.json [get]
func getDotwellKnownAssetlinksDotjsonHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusOK, []gin.H{
			{
				"relation": []string{"delegate_permission/common.handle_all_urls"},
				"target": gin.H{
					"namespace": "android_app",
					"package":   os.Getenv("DYNAMIC_LINKS_ANDROID_PACKAGE_NAME"),
					"sha256_cert_fingerprints": []string{
						"12:34:56:78:90:AB:CD:EF:12:34:56:78:90:AB:CD:EF:12:34:56:78:90:AB:CD:EF:12:34:56:78:90:AB:CD:EF",
					},
				},
			},
		})
	}
}
