package market

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"trovo-wallet-api/internal/middleware"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/gin-gonic/gin"
)

// Init registers the market's read routes. They need no signature (market
// data is public) and are rate-limited per caller.
func Init(router *gin.Engine, gc *sharedconfig.GlobalConfig) {
	limit := middleware.RateLimitMiddleware(gc, "market", 120, time.Minute)
	router.GET("/v1/market/pairs", limit, getPairs(gc))
	router.GET("/v1/market/orderbook", limit, getOrderBook(gc))
	router.GET("/v1/market/trades", limit, getTrades(gc))
	router.GET("/v1/market/candles", limit, getCandles(gc))
}

func queryInt(c *gin.Context, name string, def, max int) int {
	n, err := strconv.Atoi(c.Query(name))
	if err != nil || n <= 0 {
		return def
	}
	if n > max {
		return max
	}
	return n
}

func respondError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrNotConfigured):
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "error-market-not-configured", "message": err.Error()})
	case errors.Is(err, ErrBadPair):
		c.JSON(http.StatusBadRequest, gin.H{"error": "error-invalid-pair", "message": err.Error()})
	default:
		log.Printf("[market] %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error-temporary", "message": "Market data is not available right now. Please try again."})
	}
}

// now is the clock the handlers use (tests set it).
var now = time.Now

func timeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 20*time.Second)
}

// getPairs godoc
// @Summary Market pairs
// @Description Every tradable token against each quote currency (NAIRA_ASSET, DOLLAR_ASSET) with its last price, 24-hour change, high, low, volume and best bid and ask. Busiest first.
// @Tags market
// @Produce json
// @Success 200 {object} map[string]interface{} "{pairs: []Pair}"
// @Failure 503 {object} map[string]interface{}
// @Router /v1/market/pairs [get]
func getPairs(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := timeout()
		defer cancel()
		pairs, err := Pairs(ctx, gc, now())
		if err != nil {
			respondError(c, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"pairs": pairs})
	}
}

// getOrderBook godoc
// @Summary Order book of a pair
// @Description Asks (selling base, amount in base, cheapest first) and bids (buying base, amount in counter, best first), prices in counter per base.
// @Tags market
// @Produce json
// @Param base query string true "Base token contract address"
// @Param counter query string true "Counter (quote) token contract address"
// @Param limit query int false "Levels per side (default 20, at most 100)"
// @Success 200 {object} Book
// @Failure 400 {object} map[string]interface{}
// @Router /v1/market/orderbook [get]
func getOrderBook(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := timeout()
		defer cancel()
		book, err := OrderBook(ctx, gc, c.Query("base"), c.Query("counter"), queryInt(c, "limit", 20, 100))
		if err != nil {
			respondError(c, err)
			return
		}
		c.JSON(http.StatusOK, book)
	}
}

// getTrades godoc
// @Summary Recent trades of a pair
// @Description The last trades (up to 30 days back), newest first.
// @Tags market
// @Produce json
// @Param base query string true "Base token contract address"
// @Param counter query string true "Counter (quote) token contract address"
// @Param limit query int false "Trades (default 50, at most 200)"
// @Success 200 {object} map[string]interface{} "{trades: []Trade}"
// @Failure 400 {object} map[string]interface{}
// @Router /v1/market/trades [get]
func getTrades(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := timeout()
		defer cancel()
		trades, err := RecentTrades(ctx, gc, c.Query("base"), c.Query("counter"), queryInt(c, "limit", 50, 200), now())
		if err != nil {
			respondError(c, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"trades": trades})
	}
}

// getCandles godoc
// @Summary Price candles of a pair
// @Description Open, high, low, close and volume per period, oldest first; periods without trades are left out.
// @Tags market
// @Produce json
// @Param base query string true "Base token contract address"
// @Param counter query string true "Counter (quote) token contract address"
// @Param resolution query string false "15m, 1h, 4h, 1d or 1w (default 1h)"
// @Param limit query int false "Periods back (default 100, at most 500)"
// @Success 200 {object} map[string]interface{} "{resolution, candles: []Candle}"
// @Failure 400 {object} map[string]interface{}
// @Router /v1/market/candles [get]
func getCandles(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		res := c.DefaultQuery("resolution", "1h")
		d, ok := Resolutions[res]
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "error-invalid-resolution", "message": "resolution must be 15m, 1h, 4h, 1d or 1w"})
			return
		}
		ctx, cancel := timeout()
		defer cancel()
		candles, err := Candles(ctx, gc, c.Query("base"), c.Query("counter"), d, queryInt(c, "limit", 100, 500), now())
		if err != nil {
			respondError(c, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"resolution": res, "candles": candles})
	}
}
