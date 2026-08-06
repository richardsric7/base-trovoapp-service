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

	router.GET("/v1/announcements", func(c *gin.Context) {

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

	})
	router.GET("/v1/app-version", func(c *gin.Context) {

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

	})

}
