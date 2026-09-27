package p2p

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	p2pServices "trovo-wallet-api/internal/components/p2p/services"
	userModels "trovo-wallet-api/internal/components/users/models"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/middleware"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/gin-gonic/gin"
)

// currentUser resolves the authenticated caller from the signed request,
// matching postUsersPaymentHandler's own resolution pattern exactly.
func currentUser(c *gin.Context, gc *sharedconfig.GlobalConfig) (userModels.User, error) {
	return userModels.UserSigner(middleware.ExtractSigner(c)).GetOwner(gc.DB, gc)
}

func writeError(c *gin.Context, err error) {
	if ex, ok := err.(tErrors.GenericError); ok {
		c.JSON(ex.HTTPCode(), ex.JSONError())
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
}

type createOfferRequest struct {
	OfferType             string `json:"offerType"`
	Asset                 string `json:"asset"`
	PaymentMethodID       string `json:"paymentMethodId"`
	MerchantPayoutAddress string `json:"merchantPayoutAddress"`
	Country               string `json:"country"`
	CountryCode           string `json:"countryCode"`
	Currency              string `json:"currency"`
	PriceType             string `json:"priceType"`
	Price                 string `json:"price"`
	PriceMargin           string `json:"priceMargin"`
	MinOrderAmount        string `json:"minOrderAmount"`
	MaxOrderAmount        string `json:"maxOrderAmount"`
	AvailableLiquidity    string `json:"availableLiquidity"`
	Remark                string `json:"remark"`
}

// postOffersHandler godoc
// @Summary Create a new P2P buy/sell offer
// @Description Creates a merchant offer to buy or sell an asset for fiat on the P2P marketplace. The caller becomes the merchant on the offer.
// @Tags P2P
// @Accept json
// @Produce json
// @Param body body createOfferRequest true "Offer details: offerType (BUY/SELL), asset, paymentMethodId (an existing saved payment method), merchantPayoutAddress, country, countryCode, currency, priceType, price, priceMargin, minOrderAmount, maxOrderAmount, availableLiquidity, remark"
// @Success 201 {object} map[string]interface{} "Created offer"
// @Failure 400 {object} map[string]interface{} "Invalid JSON or validation error"
// @Security SignatureAuth
// @Router /v1/p2p/offers [post]
func postOffersHandler(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := currentUser(c, gc)
		if err != nil {
			writeError(c, err)
			return
		}
		data, _ := io.ReadAll(c.Request.Body)
		var req createOfferRequest
		if err := json.Unmarshal(data, &req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "error-invalid-json"})
			return
		}
		offer, err := p2pServices.CreateOffer(gc, user.Username, user.ID, p2pServices.CreateOfferInput{
			OfferType:             req.OfferType,
			Asset:                 req.Asset,
			PaymentMethodID:       req.PaymentMethodID,
			MerchantPayoutAddress: req.MerchantPayoutAddress,
			Country:               req.Country,
			CountryCode:           req.CountryCode,
			Currency:              req.Currency,
			PriceType:             req.PriceType,
			Price:                 req.Price,
			PriceMargin:           req.PriceMargin,
			MinOrderAmount:        req.MinOrderAmount,
			MaxOrderAmount:        req.MaxOrderAmount,
			AvailableLiquidity:    req.AvailableLiquidity,
			Remark:                req.Remark,
		})
		if err != nil {
			writeError(c, err)
			return
		}
		c.JSON(http.StatusCreated, offer)
	}
}

type updateOfferRequest struct {
	PaymentMethodID       string `json:"paymentMethodId"`
	MerchantPayoutAddress string `json:"merchantPayoutAddress"`
	PriceType             string `json:"priceType"`
	Price                 string `json:"price"`
	PriceMargin           string `json:"priceMargin"`
	MinOrderAmount        string `json:"minOrderAmount"`
	MaxOrderAmount        string `json:"maxOrderAmount"`
	AvailableLiquidity    string `json:"availableLiquidity"`
	Remark                string `json:"remark"`
}

// putOfferHandler godoc
// @Summary Update an existing P2P offer
// @Description Updates the editable fields of an offer the caller owns as merchant (price, limits, payout address, payment method, remark). Offer type and asset cannot be changed after creation.
// @Tags P2P
// @Accept json
// @Produce json
// @Param offerID path string true "Offer ID"
// @Param body body updateOfferRequest true "Fields to update"
// @Success 200 {object} map[string]interface{} "Updated offer"
// @Failure 400 {object} map[string]interface{} "Invalid JSON or validation error"
// @Security SignatureAuth
// @Router /v1/p2p/offers/{offerID} [put]
func putOfferHandler(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := currentUser(c, gc)
		if err != nil {
			writeError(c, err)
			return
		}
		data, _ := io.ReadAll(c.Request.Body)
		var req updateOfferRequest
		if err := json.Unmarshal(data, &req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "error-invalid-json"})
			return
		}
		offer, err := p2pServices.UpdateOffer(gc, c.Param("offerID"), user.ID, p2pServices.UpdateOfferInput{
			PaymentMethodID:       req.PaymentMethodID,
			MerchantPayoutAddress: req.MerchantPayoutAddress,
			PriceType:             req.PriceType,
			Price:                 req.Price,
			PriceMargin:           req.PriceMargin,
			MinOrderAmount:        req.MinOrderAmount,
			MaxOrderAmount:        req.MaxOrderAmount,
			AvailableLiquidity:    req.AvailableLiquidity,
			Remark:                req.Remark,
		})
		if err != nil {
			writeError(c, err)
			return
		}
		c.JSON(http.StatusOK, offer)
	}
}

// postCloseOfferHandler godoc
// @Summary Permanently close a P2P offer
// @Description Closes an offer the caller owns as merchant. A closed offer stops appearing on the marketplace and cannot be reactivated (unlike pausing).
// @Tags P2P
// @Produce json
// @Param offerID path string true "Offer ID"
// @Success 200 {object} map[string]interface{} "Closed offer"
// @Failure 400 {object} map[string]interface{} "Not the offer's owner, or invalid state"
// @Security SignatureAuth
// @Router /v1/p2p/offers/{offerID}/close [post]
func postCloseOfferHandler(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := currentUser(c, gc)
		if err != nil {
			writeError(c, err)
			return
		}
		offer, err := p2pServices.CloseOffer(gc, c.Param("offerID"), user.ID)
		if err != nil {
			writeError(c, err)
			return
		}
		c.JSON(http.StatusOK, offer)
	}
}

// getMarketplaceOffersHandler godoc
// @Summary Browse P2P marketplace offers
// @Description Lists active offers on the P2P marketplace, with optional filters and pagination. Used to populate the marketplace buy/sell listing screen.
// @Tags P2P
// @Produce json
// @Param offerType query string false "Filter by BUY or SELL"
// @Param asset query string false "Filter by asset code"
// @Param assetClassId query int false "Filter by asset class ID"
// @Param countryCode query string false "Filter by country code"
// @Param currency query string false "Filter by fiat currency"
// @Param page query int false "Page number (default 1)"
// @Param pageSize query int false "Results per page (default 20)"
// @Success 200 {object} map[string]interface{} "data (offers), total, page, pageSize"
// @Failure 500 {object} map[string]interface{}
// @Security SignatureAuth
// @Router /v1/p2p/offers [get]
func getMarketplaceOffersHandler(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
		assetClassID, _ := strconv.ParseUint(c.Query("assetClassId"), 10, 64)
		offers, total, err := p2pServices.ListMarketplaceOffers(gc.DB, p2pServices.MarketplaceFilter{
			OfferType:    c.Query("offerType"),
			Asset:        c.Query("asset"),
			AssetClassID: assetClassID,
			CountryCode:  c.Query("countryCode"),
			Currency:     c.Query("currency"),
			Page:         page,
			PageSize:     pageSize,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error-temporary-server-error"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": offers, "total": total, "page": page, "pageSize": pageSize})
	}
}

// getMarketplaceFacetsHandler godoc
// @Summary List available P2P marketplace filter options
// @Description Returns the distinct asset and currency values currently offered on the marketplace, for populating filter dropdowns.
// getMarketplaceFacetsHandler lists the distinct asset/currency values
// worth offering as marketplace filter options right now (Plan: filter by
// currency and asset, in addition to the new asset category filter).
// @Tags P2P
// @Produce json
// @Success 200 {object} map[string]interface{} "Facet lists"
// @Failure 500 {object} map[string]interface{}
// @Security SignatureAuth
// @Router /v1/p2p/marketplace/facets [get]
func getMarketplaceFacetsHandler(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		facets, err := p2pServices.ListMarketplaceFacets(gc.DB)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error-temporary-server-error"})
			return
		}
		c.JSON(http.StatusOK, facets)
	}
}

// getAssetClassesHandler godoc
// @Summary List P2P asset categories
// @Description Returns every asset category the marketplace filter can narrow by (e.g. stablecoin, token). Small, mostly-static reference list.
// @Tags P2P
// @Produce json
// @Success 200 {object} map[string]interface{} "data: list of asset classes"
// @Failure 500 {object} map[string]interface{}
// @Security SignatureAuth
// @Router /v1/p2p/asset-classes [get]
func getAssetClassesHandler(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		classes, err := p2pServices.ListAssetClasses(gc.DB)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error-temporary-server-error"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": classes})
	}
}

// getOfferHandler godoc
// @Summary Get a single P2P offer by ID
// @Tags P2P
// @Produce json
// @Param offerID path string true "Offer ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{} "Offer not found"
// @Security SignatureAuth
// @Router /v1/p2p/offers/{offerID} [get]
func getOfferHandler(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		offer, err := p2pServices.GetOfferByID(gc.DB, c.Param("offerID"))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "error-offer-not-found"})
			return
		}
		c.JSON(http.StatusOK, offer)
	}
}

// getOfferQuoteHandler godoc
// @Summary Get a fee quote for an offer before creating an order
// @Description Returns an itemized breakdown (fees, VAT, charges, net amounts for both buyer and seller) for a hypothetical order against this offer at the given amount, without creating or persisting anything. Used by the order-creation screen to show the customer real numbers before they commit.
// @Tags P2P
// @Produce json
// @Param offerID path string true "Offer ID"
// @Param amount query string true "Specified asset amount to quote"
// @Success 200 {object} map[string]interface{} "Fee breakdown"
// @Failure 400 {object} map[string]interface{}
// @Security SignatureAuth
// @Router /v1/p2p/offers/{offerID}/quote [get]
func getOfferQuoteHandler(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		breakdown, offer, err := p2pServices.QuoteOrderFees(gc, c.Param("offerID"), c.Query("amount"))
		if err != nil {
			writeError(c, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"offerId":                 offer.ID,
			"specifiedAssetAmount":    c.Query("amount"),
			"paymentAmount":           breakdown.PaymentAmount.String(),
			"buyerTotalFees":          breakdown.BuyerTotalFees.String(),
			"buyerTotalVat":           breakdown.BuyerTotalVat.String(),
			"buyerTotalCharges":       breakdown.BuyerTotalCharges.String(),
			"buyerNetAssetAmount":     breakdown.BuyerNetAssetAmount.String(),
			"sellerTotalFees":         breakdown.SellerTotalFees.String(),
			"sellerTotalVat":          breakdown.SellerTotalVat.String(),
			"sellerTotalCharges":      breakdown.SellerTotalCharges.String(),
			"sellerEscrowAssetAmount": breakdown.SellerEscrowAssetAmount.String(),
		})
	}
}

// postActivateOfferHandler godoc
// @Summary Activate a paused P2P offer
// @Description Makes a paused offer visible and orderable on the marketplace again.
// @Tags P2P
// @Produce json
// @Param offerID path string true "Offer ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{} "Not the offer's owner, or invalid state"
// @Security SignatureAuth
// @Router /v1/p2p/offers/{offerID}/activate [post]
func postActivateOfferHandler(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := currentUser(c, gc)
		if err != nil {
			writeError(c, err)
			return
		}
		offer, err := p2pServices.ActivateOffer(gc, c.Param("offerID"), user.ID)
		if err != nil {
			writeError(c, err)
			return
		}
		c.JSON(http.StatusOK, offer)
	}
}

// postPauseOfferHandler godoc
// @Summary Pause a P2P offer
// @Description Temporarily hides an offer from the marketplace without closing it permanently. Can be reactivated later.
// @Tags P2P
// @Produce json
// @Param offerID path string true "Offer ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{} "Not the offer's owner, or invalid state"
// @Security SignatureAuth
// @Router /v1/p2p/offers/{offerID}/pause [post]
func postPauseOfferHandler(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := currentUser(c, gc)
		if err != nil {
			writeError(c, err)
			return
		}
		offer, err := p2pServices.PauseOffer(gc, c.Param("offerID"), user.ID)
		if err != nil {
			writeError(c, err)
			return
		}
		c.JSON(http.StatusOK, offer)
	}
}

// getMyOffersHandler godoc
// @Summary List my own P2P offers
// @Description Lists every offer the caller owns as merchant, in any state (active, paused, closed).
// @Tags P2P
// @Produce json
// @Success 200 {object} map[string]interface{} "data: list of offers"
// @Failure 500 {object} map[string]interface{}
// @Security SignatureAuth
// @Router /v1/p2p/my-offers [get]
func getMyOffersHandler(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := currentUser(c, gc)
		if err != nil {
			writeError(c, err)
			return
		}
		offers, err := p2pServices.ListMerchantOffers(gc.DB, user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error-temporary-server-error"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": offers})
	}
}
