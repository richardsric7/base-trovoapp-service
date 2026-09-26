package p2p

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	p2pModels "trovo-wallet-api/internal/components/p2p/models"
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
	OfferType          string                 `json:"offerType"`
	Asset              string                 `json:"asset"`
	PaymentMethod      p2pModels.PaymentMethod `json:"paymentMethod"`
	Country            string                 `json:"country"`
	CountryCode        string                 `json:"countryCode"`
	Currency           string                 `json:"currency"`
	PriceType          string                 `json:"priceType"`
	Price              string                 `json:"price"`
	PriceMargin        string                 `json:"priceMargin"`
	MinOrderAmount     string                 `json:"minOrderAmount"`
	MaxOrderAmount     string                 `json:"maxOrderAmount"`
	AvailableLiquidity string                 `json:"availableLiquidity"`
	Remark             string                 `json:"remark"`
}

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
			OfferType:          req.OfferType,
			Asset:              req.Asset,
			PaymentMethod:      req.PaymentMethod,
			Country:            req.Country,
			CountryCode:        req.CountryCode,
			Currency:           req.Currency,
			PriceType:          req.PriceType,
			Price:              req.Price,
			PriceMargin:        req.PriceMargin,
			MinOrderAmount:     req.MinOrderAmount,
			MaxOrderAmount:     req.MaxOrderAmount,
			AvailableLiquidity: req.AvailableLiquidity,
			Remark:             req.Remark,
		})
		if err != nil {
			writeError(c, err)
			return
		}
		c.JSON(http.StatusCreated, offer)
	}
}

func getMarketplaceOffersHandler(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
		offers, total, err := p2pServices.ListMarketplaceOffers(gc.DB, p2pServices.MarketplaceFilter{
			OfferType:   c.Query("offerType"),
			Asset:       c.Query("asset"),
			CountryCode: c.Query("countryCode"),
			Currency:    c.Query("currency"),
			Page:        page,
			PageSize:    pageSize,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error-temporary-server-error"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": offers, "total": total, "page": page, "pageSize": pageSize})
	}
}

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
