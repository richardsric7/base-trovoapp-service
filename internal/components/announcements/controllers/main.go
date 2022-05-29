package announcements

import (
	"log"
	"trovo-wallet-api/internal/cache"
	announcementServices "trovo-wallet-api/internal/components/announcements/services"
	dbCon "trovo-wallet-api/internal/db"

	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Init initializes /v1/assets endpoint
func Init(router *gin.Engine, db *gorm.DB, redisCache *cache.RedisCache) {

	router.GET("/v1/announcements", func(c *gin.Context) {

		cacheKey := "[GET] /v1/announcements"

		ok, status, response := redisCache.CachedHttpResponse(cacheKey)

		if ok {
			log.Printf("[%v], served from cache\n", cacheKey)
			c.JSON(status, response)
			return
		}

		dbCon.PrintDBStats("GET /v1/announcements", db)

		announcements, _ := announcementServices.HandleGetAnnouncement(c.ClientIP(), db)

		c.JSON(http.StatusOK, announcements)

		cacheDurationInSeconds := 30 * 60 //30 minutes
		redisCache.CacheHttpResponse(cacheKey, http.StatusOK, announcements, cacheDurationInSeconds)

	})
	router.GET("/v1/app-version", func(c *gin.Context) {

		cacheKey := "[GET] /v1/app-version"

		ok, status, response := redisCache.CachedHttpResponse(cacheKey)

		if ok {
			log.Printf("[%v], served from cache\n", cacheKey)
			c.JSON(status, response)
			return
		}

		dbCon.PrintDBStats("GET /v1/app-version", db)

		appVersion := announcementServices.GetAppVersion(db)

		c.JSON(http.StatusOK, appVersion)

		cacheDurationInSeconds := 30 * 60 //30 minutes
		redisCache.CacheHttpResponse(cacheKey, http.StatusOK, appVersion, cacheDurationInSeconds)

	})

}
