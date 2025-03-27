package assets

import (
	assets "trovo-wallet-api/internal/components/assets/services"
	"trovo-wallet-api/internal/middleware"
	"trovo-wallet-api/internal/sharedconfig"

	"net/http"

	"github.com/gin-gonic/gin"
)

// Init initializes /v2/assets endpoint
func Init(router *gin.Engine, gc *sharedconfig.GlobalConfig) {

	router.GET("/v1/curated-assets", middleware.AuthenticationMiddlewareUsingTimestamp(), func(c *gin.Context) {
		cacheKey := "[GET] /v1/curated-assets"

		ok, status, response := gc.RedisCache.CachedHttpResponse(cacheKey)

		if ok {
			// log.Printf("[%v], served from cache\n", cacheKey)
			c.JSON(status, response)
			return
		}

		// dbCon.PrintDBStats("[GET] /v1/curated-assets", gc.DB)

		curatedAssets := assets.GetCuratedAssets(gc)

		c.JSON(http.StatusOK, curatedAssets)

		cacheDurationInSeconds := 20
		gc.RedisCache.CacheHttpResponse(cacheKey, http.StatusOK, curatedAssets, cacheDurationInSeconds)

	})
}
