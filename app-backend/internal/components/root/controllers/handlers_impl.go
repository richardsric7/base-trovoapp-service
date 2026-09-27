package root

import (
	"net/http"
	"os"
	root "trovo-wallet-api/internal/components/root/services"
	"trovo-wallet-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

// getHandler godoc
// @Summary Health/info check for the API
// @Description Returns basic service info (used as a lightweight liveness check; also echoes the caller's wallet address back if the X-TW-PUBLIC-KEY header is present).
// @Tags Root
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Security SignatureAuth
// @Router / [get]
func getHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		rootInfo := root.GetRootInfo()
		rootInfo.Address = middleware.ExtractAddress(c)
		c.JSON(200, rootInfo)
	}
}

// getDotwellKnownAppleAppSiteAssociationHandler godoc
// @Summary Apple App Site Association file (universal links)
// @Description Serves the well-known apple-app-site-association document iOS uses to verify this domain is allowed to open links in the Trovo mobile app (universal links / deep linking).
// @Tags Root
// @Produce json
// @Success 200 {object} map[string]interface{}
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
// @Summary Android Asset Links file (app links)
// @Description Serves the well-known assetlinks.json document Android uses to verify this domain is allowed to open links in the Trovo mobile app (Android App Links / deep linking).
// @Tags Root
// @Produce json
// @Success 200 {array} map[string]interface{}
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
