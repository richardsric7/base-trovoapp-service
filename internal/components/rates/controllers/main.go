package rates

import (
	"log"
	"net/http"
	ratesService "trovo-wallet-api/internal/components/rates/services"
	dbCon "trovo-wallet-api/internal/db"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/gin-gonic/gin"
)

// Init initializes /v2/assets endpoint
func Init(router *gin.Engine, gc *sharedconfig.GlobalConfig) {

	router.GET("/v1/rates", func(c *gin.Context) {

		cacheKey := "[GET] /v1/rates"

		ok, status, response := gc.RedisCache.CachedHttpResponse(cacheKey)

		if ok {
			log.Printf("[%v], served from cache\n", cacheKey)
			c.JSON(status, response)
			return
		}

		dbCon.PrintDBStats("GET /v1/rates", gc.DB)

		rate := ratesService.HandleGetRates(gc.DB)

		c.JSON(http.StatusOK, rate)
		cacheDurationInSeconds := 30 * 60 //2 hours
		if len(rate) > 0 {
			gc.RedisCache.CacheHttpResponse(cacheKey, http.StatusOK, rate, cacheDurationInSeconds)

		}

	})

}
