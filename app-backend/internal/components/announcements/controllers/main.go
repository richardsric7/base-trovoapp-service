package announcements

import (
	announcementServices "trovo-wallet-api/internal/components/announcements/services"
	dbCon "trovo-wallet-api/internal/db"
	"trovo-wallet-api/internal/sharedconfig"

	"net/http"

	"github.com/gin-gonic/gin"
)

// Init initializes /v1/assets endpoint
func Init(router *gin.Engine, gc *sharedconfig.GlobalConfig) {

	router.GET("/v1/announcements", getAnnouncementsHandler(gc))
	router.GET("/v1/app-version", getAppVersionHandler(gc))

}

// getAnnouncementsHandler godoc
// @Summary Get active in-app announcements
// @Description Returns announcement banners/messages to show in the client apps, targeted by the caller's IP-derived location. No authentication required. Response is cached for up to 30 minutes.
// @Tags Announcements
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /v1/announcements [get]
func getAnnouncementsHandler(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {

		cacheKey := "[GET] /v1/announcements"

		ok, status, response := gc.RedisCache.CachedHttpResponse(cacheKey)

		if ok {
			// log.Printf("[%v], served from cache\n", cacheKey)
			c.JSON(status, response)
			return
		}

		dbCon.PrintDBStats("GET /v1/announcements", gc.DB)

		announcements, _ := announcementServices.HandleGetAnnouncement(c.ClientIP(), gc.DB)

		c.JSON(http.StatusOK, announcements)

		cacheDurationInSeconds := 30 * 60 //30 minutes
		gc.RedisCache.CacheHttpResponse(cacheKey, http.StatusOK, announcements, cacheDurationInSeconds)

	}
}

// getAppVersionHandler godoc
// @Summary Get the minimum/latest supported app version
// @Description Returns the current app-version requirements the mobile client uses to prompt for a forced or optional update. No authentication required. Response is cached for up to 30 minutes.
// @Tags Announcements
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /v1/app-version [get]
func getAppVersionHandler(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {

		cacheKey := "[GET] /v1/app-version"

		ok, status, response := gc.RedisCache.CachedHttpResponse(cacheKey)

		if ok {
			// log.Printf("[%v], served from cache\n", cacheKey)
			c.JSON(status, response)
			return
		}

		dbCon.PrintDBStats("GET /v1/app-version", gc.DB)

		appVersion := announcementServices.GetAppVersion(gc.DB)

		c.JSON(http.StatusOK, appVersion)

		cacheDurationInSeconds := 30 * 60 //30 minutes
		gc.RedisCache.CacheHttpResponse(cacheKey, http.StatusOK, appVersion, cacheDurationInSeconds)

	}
}
