package assets

import (
	assets "trovo-wallet-api/internal/components/assets/services"
	"trovo-wallet-api/internal/sharedconfig"

	"net/http"

	"github.com/gin-gonic/gin"
)

// getCuratedAssetsHandler godoc
// @Summary GET /v1/curated-assets
// @Tags assets
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /v1/curated-assets [get]
func getCuratedAssetsHandler(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
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

	}
}
