// Package publicmarkets registers the Public Markets routes: the apps'
// /v1/public-markets, the exchange partners' /v1/trovo-api/public-markets
// and the Custodian / Dealing Member callbacks.
package publicmarkets

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	pmModels "trovo-wallet-api/internal/components/publicmarkets/models"
	"trovo-wallet-api/internal/components/publicmarkets/partners"
	pmsvc "trovo-wallet-api/internal/components/publicmarkets/services"
	servicelinkServices "trovo-wallet-api/internal/components/servicelinks/services"
	usersDB "trovo-wallet-api/internal/components/users/db"
	userModels "trovo-wallet-api/internal/components/users/models"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/middleware"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"gorm.io/gorm/clause"
)

// Engine is the running engine (set by Init).
var Engine *pmsvc.Engine

// Init builds the engine, registers the routes and starts the engine's
// background loops.
func Init(router *gin.Engine, gc *sharedconfig.GlobalConfig) {
	Engine = pmsvc.NewEngine(gc, pmsvc.NewBaseChain())
	Engine.Notify = func(username, title, body string, data map[string]string) {
		user, err := usersDB.GetUser(username, gc.DB, gc)
		if err != nil || user.Username == "" {
			return
		}
		user.SendPushMessage(title, body, "", data, gc)
		gc.PublishUserStreamEvent(username, "publicMarketsEvent", gin.H{"title": title, "body": body, "data": data})
	}
	Register(router, gc, Engine)
	Engine.Start()
}

// Register adds the Public Markets routes, served by e.
func Register(router *gin.Engine, gc *sharedconfig.GlobalConfig, e *pmsvc.Engine) {
	h := handlers{e: e, gc: gc}

	// public listing
	public := middleware.RateLimitMiddleware(gc, "public-markets-read", 120, time.Minute)
	router.GET("/v1/public-markets", public, h.listAssets)
	router.GET("/v1/public-markets/assets/:assetCode", public, h.getAsset)
	router.GET("/v1/public-markets/assets/:assetCode/prices", public, h.getPrices)

	// the app (signed requests)
	auth := middleware.AuthenticationMiddlewareUsingTimestamp()
	router.GET("/v1/public-markets/assets/:assetCode/quote", auth, h.quote)
	router.POST("/v1/public-markets/assets/:assetCode/buy", auth, middleware.RateLimitMiddleware(gc, "public-markets-trade", 30, time.Minute), h.buy)
	router.POST("/v1/public-markets/assets/:assetCode/sell", auth, middleware.RateLimitMiddleware(gc, "public-markets-trade", 30, time.Minute), h.sell)
	router.GET("/v1/public-markets/portfolio", auth, h.portfolio)
	router.GET("/v1/public-markets/orders", auth, h.myOrders)
	router.GET("/v1/public-markets/orders/:orderID", auth, h.myOrder)
	router.GET("/v1/public-markets/dividends", auth, h.myDividends)

	// exchange partners, on the service-link API (/v1/trovo-api/...)
	ex := router.Group("/v1/trovo-api/public-markets", hmacAPIKey(), middleware.AuthenticationMiddlewareUsingAPIKey(gc), h.exchangeAuth)
	read := middleware.RateLimitMiddleware(gc, "trovo-api-public-markets-read", 120, time.Minute)
	write := middleware.RateLimitMiddleware(gc, "trovo-api-public-markets-write", 60, time.Minute)
	ex.GET("/assets", read, h.exAssets)
	ex.GET("/assets/:assetCode", read, h.exAsset)
	ex.POST("/wallets", write, h.idempotent, h.exProvision)
	ex.GET("/wallets/:walletID", read, h.exWallet)
	ex.POST("/creation", write, h.idempotent, h.exCreation)
	ex.POST("/redemption", write, h.idempotent, h.exRedemption)
	ex.GET("/orders/:orderID", read, h.exOrder)
	ex.POST("/confirmations", write, h.idempotent, h.exConfirm)
	ex.GET("/account", read, h.exAccount)
	ex.POST("/deposits", write, h.idempotent, h.exDeposit)

	// Custodian and Dealing Member callbacks (HMAC-signed)
	router.POST("/v1/custodian-partners/callbacks/settlement", h.partnerCallback("CUSTODIAN", pmsvc.EventSettlement))
	router.POST("/v1/custodian-partners/callbacks/corporate-action", h.partnerCallback("CUSTODIAN", pmsvc.EventCorporateAction))
	router.POST("/v1/custodian-partners/callbacks/positions", h.partnerCallback("CUSTODIAN", pmsvc.EventPositionFeed))
	router.POST("/v1/dealing-member-partners/callbacks/execution", h.partnerCallback("DEALING_MEMBER", pmsvc.EventExecution))
}

type handlers struct {
	e  *pmsvc.Engine
	gc *sharedconfig.GlobalConfig
}

func fail(c *gin.Context, err error) {
	if e, ok := pmsvc.AsError(err); ok {
		c.JSON(e.Status, e)
		return
	}
	if ge, ok := err.(tErrors.GenericError); ok {
		c.JSON(ge.HTTPCode(), ge.JSONError())
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{"error": "error-temporary", "message": "Something went wrong. Please try again."})
}

func ctxTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 60*time.Second)
}

// ------------------------------------------------------------- public

// listAssets godoc
// @Summary List Public Markets assets
// @Description Tokenized NGX equities and FMDQ bonds with their reference prices and market session. Filter by market (NGX, FMDQ), type (EQUITY, BOND) and a search.
// @Tags Public Markets
// @Produce json
// @Param market query string false "NGX or FMDQ"
// @Param type query string false "EQUITY or BOND"
// @Param search query string false "Code or name"
// @Success 200 {object} map[string]interface{}
// @Router /v1/public-markets [get]
func (h handlers) listAssets(c *gin.Context) {
	ctx, cancel := ctxTimeout()
	defer cancel()
	c.JSON(http.StatusOK, gin.H{"assets": h.e.ListAssets(ctx, c.Query("market"), c.Query("type"), c.Query("search"))})
}

// getAsset godoc
// @Summary A Public Markets asset
// @Description Price, key statistics, custody chain, market session and corporate actions.
// @Tags Public Markets
// @Produce json
// @Param assetCode path string true "Asset code (MTNN-T) or ticker (MTNN)"
// @Success 200 {object} pmsvc.AssetView
// @Router /v1/public-markets/assets/{assetCode} [get]
func (h handlers) getAsset(c *gin.Context) {
	a, err := h.e.AssetByCode(c.Param("assetCode"))
	if err != nil {
		fail(c, err)
		return
	}
	ctx, cancel := ctxTimeout()
	defer cancel()
	c.JSON(http.StatusOK, h.e.ViewAsset(ctx, a, true))
}

// getPrices godoc
// @Summary Price chart of an asset
// @Tags Public Markets
// @Produce json
// @Param assetCode path string true "Asset code"
// @Param range query string false "1D, 1W, 1M, 3M, 1Y or All" default(1D)
// @Success 200 {object} map[string]interface{}
// @Router /v1/public-markets/assets/{assetCode}/prices [get]
func (h handlers) getPrices(c *gin.Context) {
	a, err := h.e.AssetByCode(c.Param("assetCode"))
	if err != nil {
		fail(c, err)
		return
	}
	rng := c.DefaultQuery("range", "1D")
	c.JSON(http.StatusOK, gin.H{"assetCode": a.AssetCode, "range": rng, "points": h.e.PriceHistory(a, rng), "previousClose": a.PreviousClose})
}

// ------------------------------------------------------------- app

func (h handlers) user(c *gin.Context) (*userModels.User, bool) {
	u, err := userModels.GetSlimUserBySigner(middleware.ExtractSigner(c), h.gc.DB, h.gc)
	if err != nil || u.Username == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "error-unknown-user", "message": "Unknown user."})
		return nil, false
	}
	return &u, true
}

// quote godoc
// @Summary Quote a buy or sell
// @Description Price, Trovo fee, estimated quantity or proceeds, and whether it fills instantly (fast / netted) or at the next session (slow).
// @Tags Public Markets
// @Produce json
// @Param assetCode path string true "Asset code"
// @Param side query string true "buy or sell"
// @Param amount query string false "buy: amount to spend (funding currency, fee included)"
// @Param quantity query string false "sell: tokens"
// @Success 200 {object} pmsvc.Quote
// @Security SignatureAuth
// @Router /v1/public-markets/assets/{assetCode}/quote [get]
func (h handlers) quote(c *gin.Context) {
	a, err := h.e.AssetByCode(c.Param("assetCode"))
	if err != nil {
		fail(c, err)
		return
	}
	ctx, cancel := ctxTimeout()
	defer cancel()
	var q pmsvc.Quote
	if strings.EqualFold(c.Query("side"), "sell") {
		q, err = h.e.QuoteRedemption(ctx, a, decimalOf(c.Query("quantity")))
	} else {
		q, err = h.e.QuoteCreation(ctx, a, decimalOf(c.Query("amount")))
	}
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(http.StatusOK, q)
}

func decimalOf(s string) decimal.Decimal {
	v, _ := decimal.NewFromString(strings.ReplaceAll(strings.TrimSpace(s), ",", ""))
	return v
}

// buy godoc
// @Summary Buy a Public Markets asset
// @Description Two calls, like every wallet operation. Without transactionSignature: returns the quote and the operation (transaction) to sign - the wallet pays the amount to the Public Markets treasury. With transaction and transactionSignature: submits it and places the order, which fills once the payment is mined (instantly from Custodian inventory, or at the next session).
// @Tags Public Markets
// @Accept json
// @Produce json
// @Param assetCode path string true "Asset code"
// @Param body body pmsvc.TradeRequest true "walletAddress and amount; then transaction and transactionSignature"
// @Success 200 {object} pmsvc.TradeResponse
// @Security SignatureAuth
// @Router /v1/public-markets/assets/{assetCode}/buy [post]
func (h handlers) buy(c *gin.Context) { h.trade(c, true) }

// sell godoc
// @Summary Sell a Public Markets asset
// @Description Two calls. The wallet sends the tokens to the asset's issuing Safe; the sale is netted against the day's demand and paid at once, or sold by the Dealing Member and paid after the Custodian settles.
// @Tags Public Markets
// @Accept json
// @Produce json
// @Param assetCode path string true "Asset code"
// @Param body body pmsvc.TradeRequest true "walletAddress and quantity; then transaction and transactionSignature"
// @Success 200 {object} pmsvc.TradeResponse
// @Security SignatureAuth
// @Router /v1/public-markets/assets/{assetCode}/sell [post]
func (h handlers) sell(c *gin.Context) { h.trade(c, false) }

func (h handlers) trade(c *gin.Context, buy bool) {
	u, ok := h.user(c)
	if !ok {
		return
	}
	a, err := h.e.AssetByCode(c.Param("assetCode"))
	if err != nil {
		fail(c, err)
		return
	}
	var r pmsvc.TradeRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error-invalid-json", "message": "Invalid request body."})
		return
	}
	ctx, cancel := ctxTimeout()
	defer cancel()
	var res *pmsvc.TradeResponse
	if buy {
		res, err = h.e.AppBuy(ctx, u, a, r)
	} else {
		res, err = h.e.AppSell(ctx, u, a, r)
	}
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(http.StatusOK, res)
}

// portfolio godoc
// @Summary My Stocks
// @Description Holdings across the user's wallets with average cost, returns, income, open orders and recent activity.
// @Tags Public Markets
// @Produce json
// @Success 200 {object} pmsvc.Portfolio
// @Security SignatureAuth
// @Router /v1/public-markets/portfolio [get]
func (h handlers) portfolio(c *gin.Context) {
	u, ok := h.user(c)
	if !ok {
		return
	}
	ctx, cancel := ctxTimeout()
	defer cancel()
	c.JSON(http.StatusOK, h.e.PortfolioOf(ctx, u))
}

// myOrders godoc
// @Summary My Public Markets orders
// @Tags Public Markets
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Security SignatureAuth
// @Router /v1/public-markets/orders [get]
func (h handlers) myOrders(c *gin.Context) {
	u, ok := h.user(c)
	if !ok {
		return
	}
	var orders []pmModels.Order
	h.gc.DB.Where("username = ?", u.Username).Order("created_at DESC").Limit(100).Find(&orders)
	c.JSON(http.StatusOK, gin.H{"orders": orders})
}

// myOrder godoc
// @Summary One of my orders, with its timeline
// @Tags Public Markets
// @Produce json
// @Param orderID path string true "Order ID (PM-CR-..., PM-RD-...)"
// @Success 200 {object} pmsvc.OrderView
// @Security SignatureAuth
// @Router /v1/public-markets/orders/{orderID} [get]
func (h handlers) myOrder(c *gin.Context) {
	u, ok := h.user(c)
	if !ok {
		return
	}
	var o pmModels.Order
	if h.gc.DB.First(&o, "id = ? AND username = ?", c.Param("orderID"), u.Username).Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "error-unknown-order", "message": "Unknown order."})
		return
	}
	c.JSON(http.StatusOK, h.e.OrderViewOf(&o))
}

// myDividends godoc
// @Summary My dividends and coupons
// @Description Each with its breakdown: units on the record date, gross, withholding tax and the net paid.
// @Tags Public Markets
// @Produce json
// @Param assetCode query string false "Only this asset"
// @Success 200 {object} map[string]interface{}
// @Security SignatureAuth
// @Router /v1/public-markets/dividends [get]
func (h handlers) myDividends(c *gin.Context) {
	u, ok := h.user(c)
	if !ok {
		return
	}
	c.JSON(http.StatusOK, gin.H{"dividends": h.e.DividendsOf(u, c.Query("assetCode"))})
}

// ------------------------------------------------------------- exchanges

// hmacAPIKey accepts the integration specification's
// "Authorization: HMAC {apiKey}:{signature}" as well as the service-link
// API key header.
func hmacAPIKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		if auth := c.GetHeader("Authorization"); strings.HasPrefix(auth, "HMAC ") {
			kv := strings.SplitN(strings.TrimPrefix(auth, "HMAC "), ":", 2)
			if len(kv) == 2 {
				c.Request.Header.Set("X-TW-SERVICE-LINK-API-KEY", strings.TrimSpace(kv[0]))
				c.Set("pmSignature", strings.TrimSpace(kv[1]))
			}
		}
		c.Next()
	}
}

// exchangeAuth admits an active, onboarded exchange whose request is
// signed with its signing secret.
func (h handlers) exchangeAuth(c *gin.Context) {
	sl, err := servicelinkServices.GetServiceLinkByAPIKey(middleware.ExtractServiceLinkApiKey(c), h.gc.DB)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "error-invalid-api-key", "message": "Invalid API key."})
		return
	}
	if sl.Inactive != 0 || sl.Suspended != 0 || sl.Verified == 0 {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "error-partner-suspended", "message": "This service link is not active."})
		return
	}
	p, err := h.e.ExchangePartnerOf(sl.ID)
	if err != nil {
		e, _ := pmsvc.AsError(err)
		c.AbortWithStatusJSON(e.Status, e)
		return
	}
	body, _ := io.ReadAll(io.LimitReader(c.Request.Body, 1<<20))
	c.Request.Body = io.NopCloser(bytes.NewReader(body))
	sig := c.GetString("pmSignature")
	if sig == "" {
		sig = c.GetHeader("X-Trovotech-Signature")
	}
	if err := h.e.VerifyExchangeSignature(p, c.GetHeader("X-Trovotech-Timestamp"), body, sig); err != nil {
		e, _ := pmsvc.AsError(err)
		c.AbortWithStatusJSON(e.Status, e)
		return
	}
	c.Set("pmServiceLinkID", sl.ID)
	c.Set("pmBody", body)
	c.Next()
}

// recorder captures a response so it can be replayed for a repeated
// Idempotency-Key.
type recorder struct {
	gin.ResponseWriter
	buf bytes.Buffer
}

func (r *recorder) Write(b []byte) (int, error) {
	r.buf.Write(b)
	return r.ResponseWriter.Write(b)
}

// idempotent makes a write replayable: the first request with a key runs,
// a repeat returns the original response, the same key with a different
// body is refused.
func (h handlers) idempotent(c *gin.Context) {
	key := strings.TrimSpace(c.GetHeader("Idempotency-Key"))
	if key == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "error-missing-idempotency-key", "message": "The Idempotency-Key header is required.", "field": "Idempotency-Key"})
		return
	}
	sl := c.GetString("pmServiceLinkID")
	body, _ := c.Get("pmBody")
	sum := sha256.Sum256(append([]byte(c.FullPath()+"|"), body.([]byte)...))
	hash := hex.EncodeToString(sum[:])
	rec := pmModels.IdempotencyRecord{ServiceLinkID: sl, Key: key, Path: c.FullPath(), RequestHash: hash, CreatedAt: time.Now().UTC()}
	res := h.gc.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&rec)
	if res.Error != nil {
		fail(c, res.Error)
		c.Abort()
		return
	}
	if res.RowsAffected == 0 {
		var prev pmModels.IdempotencyRecord
		h.gc.DB.First(&prev, "service_link_id = ? AND key = ?", sl, key)
		switch {
		case prev.RequestHash != hash:
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": "error-idempotency-key-reused", "message": "This Idempotency-Key was used for a different request."})
		case !prev.Done:
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": "error-request-in-progress", "message": "A request with this Idempotency-Key is still being processed."})
		default:
			c.Header("Idempotent-Replayed", "true")
			c.Data(prev.Status, "application/json; charset=utf-8", []byte(prev.Body))
			c.Abort()
		}
		return
	}
	w := &recorder{ResponseWriter: c.Writer}
	c.Writer = w
	c.Next()
	status := c.Writer.Status()
	if status >= 500 {
		h.gc.DB.Delete(&pmModels.IdempotencyRecord{}, rec.ID) // a server error may be retried
		return
	}
	h.gc.DB.Model(&pmModels.IdempotencyRecord{}).Where("id = ?", rec.ID).Updates(map[string]interface{}{"status": status, "body": w.buf.String(), "done": true})
}

func (h handlers) bind(c *gin.Context, v interface{}) bool {
	body, _ := c.Get("pmBody")
	if err := json.Unmarshal(body.([]byte), v); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error-invalid-json", "message": "Invalid request body."})
		return false
	}
	return true
}

// exAssets godoc
// @Summary Public Markets assets (exchange partners)
// @Description Assets an exchange can offer its customers, with reference prices and whether they accept orders.
// @Tags Public Markets (exchange partners)
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Security ApiKeyAuth
// @Router /v1/trovo-api/public-markets/assets [get]
func (h handlers) exAssets(c *gin.Context) { h.listAssets(c) }

// exAsset godoc
// @Summary A Public Markets asset with its reference price (exchange partners)
// @Tags Public Markets (exchange partners)
// @Produce json
// @Param assetCode path string true "Asset code"
// @Success 200 {object} pmsvc.AssetView
// @Security ApiKeyAuth
// @Router /v1/trovo-api/public-markets/assets/{assetCode} [get]
func (h handlers) exAsset(c *gin.Context) { h.getAsset(c) }

// exProvision godoc
// @Summary Provision a customer wallet (§6.3.1)
// @Description Opens an individually addressed wallet for one of the exchange's customers, from the minimum identity and tax data Trovotech needs (full KYC stays with the exchange). A repeated externalUserRef returns the same wallet. A missing field is refused with a field-level error.
// @Tags Public Markets (exchange partners)
// @Accept json
// @Produce json
// @Param Idempotency-Key header string true "The exchange's idempotency key"
// @Param body body pmsvc.ProvisionRequest true "externalUserRef, legalName, taxIdentifier, residencyCountry, nationality, ndpaConsent"
// @Success 201 {object} pmsvc.WalletView
// @Failure 400 {object} map[string]interface{} "field-level error"
// @Security ApiKeyAuth
// @Router /v1/trovo-api/public-markets/wallets [post]
func (h handlers) exProvision(c *gin.Context) {
	var r pmsvc.ProvisionRequest
	if !h.bind(c, &r) {
		return
	}
	ctx, cancel := ctxTimeout()
	defer cancel()
	w, created, err := h.e.ProvisionWallet(ctx, c.GetString("pmServiceLinkID"), r)
	if err != nil {
		fail(c, err)
		return
	}
	status := http.StatusOK
	if created {
		status = http.StatusCreated
	}
	c.JSON(status, w)
}

// exWallet godoc
// @Summary A customer wallet's holdings (exchange partners)
// @Tags Public Markets (exchange partners)
// @Produce json
// @Param walletID path string true "walletId"
// @Success 200 {object} map[string]interface{}
// @Security ApiKeyAuth
// @Router /v1/trovo-api/public-markets/wallets/{walletID} [get]
func (h handlers) exWallet(c *gin.Context) {
	w, err := h.e.PartnerWallet(c.GetString("pmServiceLinkID"), c.Param("walletID"))
	if err != nil {
		fail(c, err)
		return
	}
	var rows []pmModels.LedgerEntry
	h.gc.DB.Where("LOWER(wallet_address) = ? AND balance <> '0'", strings.ToLower(w.WalletAddress)).Find(&rows)
	holdings := []gin.H{}
	for _, r := range rows {
		var a pmModels.Asset
		if h.gc.DB.First(&a, "id = ?", r.AssetID).Error != nil {
			continue
		}
		qty := decimalOf(r.Balance).Shift(-int32(a.TokenDecimals))
		holdings = append(holdings, gin.H{"assetCode": a.AssetCode, "quantity": qty.String(), "price": a.LastPrice, "value": qty.Mul(decimalOf(a.LastPrice)).StringFixed(2)})
	}
	c.JSON(http.StatusOK, gin.H{"walletId": w.ID, "publicKey": w.WalletAddress, "status": w.Status, "externalUserRef": w.ExternalUserRef, "holdings": holdings})
}

// exCreation godoc
// @Summary Place a creation order (§6.3.2)
// @Description Buys for a customer's wallet, paid from the exchange's prefunded balance. Answers 202 pending; a creation.settled webhook follows once the Custodian confirms settlement. Runs on the same engine as Trovo App orders.
// @Tags Public Markets (exchange partners)
// @Accept json
// @Produce json
// @Param Idempotency-Key header string true "The exchange's idempotency key"
// @Param body body pmsvc.ExchangeOrderRequest true "walletId, assetCode, amount, amountCurrency (NGN), externalOrderRef"
// @Success 202 {object} pmsvc.ExchangeOrderView
// @Security ApiKeyAuth
// @Router /v1/trovo-api/public-markets/creation [post]
func (h handlers) exCreation(c *gin.Context) { h.exOrderPlace(c, true) }

// exRedemption godoc
// @Summary Place a redemption order (§6.3.3)
// @Description Sells from a customer's wallet: quantity (tokens) or amount (NGN). Answers 202 pending; redemption.settled follows, with the proceeds credited to the exchange's balance.
// @Tags Public Markets (exchange partners)
// @Accept json
// @Produce json
// @Param Idempotency-Key header string true "The exchange's idempotency key"
// @Param body body pmsvc.ExchangeOrderRequest true "walletId, assetCode, quantity or amount, externalOrderRef"
// @Success 202 {object} pmsvc.ExchangeOrderView
// @Security ApiKeyAuth
// @Router /v1/trovo-api/public-markets/redemption [post]
func (h handlers) exRedemption(c *gin.Context) { h.exOrderPlace(c, false) }

func (h handlers) exOrderPlace(c *gin.Context, creation bool) {
	var r pmsvc.ExchangeOrderRequest
	if !h.bind(c, &r) {
		return
	}
	ctx, cancel := ctxTimeout()
	defer cancel()
	var o *pmModels.Order
	var err error
	if creation {
		o, err = h.e.PlaceExchangeCreation(ctx, c.GetString("pmServiceLinkID"), r)
	} else {
		o, err = h.e.PlaceExchangeRedemption(ctx, c.GetString("pmServiceLinkID"), r)
	}
	if err != nil {
		fail(c, err)
		return
	}
	h.e.Kick()
	c.JSON(http.StatusAccepted, pmsvc.ExchangeView(o))
}

// exOrder godoc
// @Summary An order's status (exchange partners)
// @Description pending, settled, rejected or failed (internal states are not exposed).
// @Tags Public Markets (exchange partners)
// @Produce json
// @Param orderID path string true "trovotechOrderId"
// @Success 200 {object} pmsvc.ExchangeOrderView
// @Security ApiKeyAuth
// @Router /v1/trovo-api/public-markets/orders/{orderID} [get]
func (h handlers) exOrder(c *gin.Context) {
	var o pmModels.Order
	if h.gc.DB.First(&o, "id = ? AND service_link_id = ?", c.Param("orderID"), c.GetString("pmServiceLinkID")).Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "error-unknown-order", "message": "Unknown order."})
		return
	}
	c.JSON(http.StatusOK, pmsvc.ExchangeView(&o))
}

// exConfirm godoc
// @Summary Confirm receipt of a webhook (§6.4.3)
// @Description Required for every dividend.paid: the exchange confirms it passed the (already withheld) amount on. Unconfirmed events past the SLA are escalated.
// @Tags Public Markets (exchange partners)
// @Accept json
// @Produce json
// @Param Idempotency-Key header string true "The exchange's idempotency key"
// @Param body body pmsvc.ConfirmationRequest true "eventId, or event and walletId; confirmedAt"
// @Success 200 {object} map[string]interface{}
// @Security ApiKeyAuth
// @Router /v1/trovo-api/public-markets/confirmations [post]
func (h handlers) exConfirm(c *gin.Context) {
	var r pmsvc.ConfirmationRequest
	if !h.bind(c, &r) {
		return
	}
	w, err := h.e.Confirm(c.GetString("pmServiceLinkID"), r)
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"eventId": w.ID, "event": w.Event, "walletId": w.PartnerWalletID, "confirmedAt": w.ConfirmedAt})
}

// exAccount godoc
// @Summary The exchange's prefunded balance and its movements
// @Tags Public Markets (exchange partners)
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Security ApiKeyAuth
// @Router /v1/trovo-api/public-markets/account [get]
func (h handlers) exAccount(c *gin.Context) {
	sl := c.GetString("pmServiceLinkID")
	var p pmModels.ExchangePartner
	h.gc.DB.First(&p, "service_link_id = ?", sl)
	var entries []pmModels.ExchangeLedgerEntry
	h.gc.DB.Where("service_link_id = ?", sl).Order("id DESC").Limit(100).Find(&entries)
	treasury := ""
	if h.e.Chain != nil {
		treasury = h.e.Chain.TreasurySafe().Hex()
	}
	c.JSON(http.StatusOK, gin.H{"balance": p.Balance, "currency": pmsvc.LoadSettings(h.gc.DB).FundingAssetCode, "fundingAddress": p.FundingAddress,
		"depositAddress": treasury, "environment": p.Environment, "entries": entries})
}

// exDeposit godoc
// @Summary Credit a deposit to the prefunded balance
// @Description After sending the funding stablecoin from the registered funding wallet to the Public Markets treasury, submit the transaction hash to have it credited.
// @Tags Public Markets (exchange partners)
// @Accept json
// @Produce json
// @Param Idempotency-Key header string true "The exchange's idempotency key"
// @Param body body map[string]string true "txHash"
// @Success 200 {object} map[string]interface{}
// @Security ApiKeyAuth
// @Router /v1/trovo-api/public-markets/deposits [post]
func (h handlers) exDeposit(c *gin.Context) {
	var r struct {
		TxHash string `json:"txHash"`
	}
	if !h.bind(c, &r) {
		return
	}
	ctx, cancel := ctxTimeout()
	defer cancel()
	amount, err := h.e.VerifyDeposit(ctx, c.GetString("pmServiceLinkID"), r.TxHash, "exchange")
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"txHash": r.TxHash, "credited": amount.String()})
}

// ------------------------------------------------------------- partner callbacks

// partnerCallback receives a Custodian or Dealing Member webhook: the
// partner names itself in X-Partner-Code and signs the body
// (X-Trovotech-Timestamp + X-Partner-Signature, HMAC-SHA256 with the
// secret of its credentials). Every webhook is idempotent on the
// partner's reference.
func (h handlers) partnerCallback(partnerType, kind string) gin.HandlerFunc {
	return func(c *gin.Context) {
		code := strings.TrimSpace(c.GetHeader("X-Partner-Code"))
		var credRef string
		switch partnerType {
		case "CUSTODIAN":
			var p pmModels.Custodian
			if h.gc.DB.First(&p, "code = ? AND active = ?", code, true).Error != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "error-unknown-partner"})
				return
			}
			credRef = p.CredentialsRef
		default:
			var p pmModels.DealingMember
			if h.gc.DB.First(&p, "code = ? AND active = ?", code, true).Error != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "error-unknown-partner"})
				return
			}
			credRef = p.CredentialsRef
		}
		body, _ := io.ReadAll(io.LimitReader(c.Request.Body, 1<<20))
		creds, err := partners.ResolveCredentials(credRef)
		sig := c.GetHeader("X-Partner-Signature")
		if sig == "" {
			sig = c.GetHeader("X-Custodian-Signature")
		}
		if err != nil || creds.Secret == "" || !partners.VerifySignature(creds.Secret, c.GetHeader("X-Trovotech-Timestamp"), body, sig) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "error-invalid-signature"})
			return
		}
		ref := referenceOf(kind, body, c.GetHeader("Idempotency-Key"))
		dup, err := h.e.ReceivePartnerEvent(partnerType+":"+code, kind, ref, body)
		if err != nil {
			var pe *pmsvc.Error
			if errors.As(err, &pe) {
				c.JSON(pe.Status, pe)
				return
			}
			// stored; it is retried by the engine
			c.JSON(http.StatusAccepted, gin.H{"status": "accepted"})
			return
		}
		if dup {
			c.JSON(http.StatusOK, gin.H{"status": "duplicate"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "processed"})
	}
}

// referenceOf is the partner's own reference of a webhook (what makes a
// redelivery a duplicate).
func referenceOf(kind string, body []byte, idemKey string) string {
	var m map[string]interface{}
	_ = json.Unmarshal(body, &m)
	str := func(k string) string {
		if v, ok := m[k].(string); ok {
			return v
		}
		return ""
	}
	switch kind {
	case pmsvc.EventSettlement:
		return str("instructionId") + ":" + str("status") + ":" + str("custodianReference")
	case pmsvc.EventExecution:
		return str("orderId") + ":" + str("status")
	case pmsvc.EventCorporateAction:
		return str("assetCode") + ":" + str("eventType") + ":" + str("recordDate")
	}
	if idemKey != "" {
		return idemKey
	}
	sum := sha256.Sum256(body)
	return hex.EncodeToString(sum[:])
}
