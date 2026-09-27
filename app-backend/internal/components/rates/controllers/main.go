package rates

import (
	"net/http"
	ratesService "trovo-wallet-api/internal/components/rates/services"
	dbCon "trovo-wallet-api/internal/db"
	"trovo-wallet-api/internal/middleware"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/gin-gonic/gin"
)

// Init initializes /v2/assets endpoint
func Init(router *gin.Engine, gc *sharedconfig.GlobalConfig) {

	router.GET("/v1/rates", middleware.AuthenticationMiddlewareUsingTimestamp(), getRatesHandler(gc))

}

// getRatesHandler godoc
// @Summary Get current exchange rates
// @Description Returns the current exchange rates the wallet uses to price assets against fiat/reference currencies. Response is cached for up to 30 minutes.
// @Tags Rates
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Security SignatureAuth
// @Router /v1/rates [get]
func getRatesHandler(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {

		cacheKey := "[GET] /v1/rates"

		ok, status, response := gc.RedisCache.CachedHttpResponse(cacheKey)

		if ok {
			// log.Printf("[%v], served from cache\n", cacheKey)
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

	}
}
