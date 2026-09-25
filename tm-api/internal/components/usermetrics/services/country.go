package usermetrics

import (
	usermetricsDB "admin-panel-dashboard/internal/components/usermetrics/db"
	"admin-panel-dashboard/internal/middleware"
	"admin-panel-dashboard/internal/models"
	serverResponse "admin-panel-dashboard/internal/server/response"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SaveCountryHandler handles create/update country
// @Summary Save country
// @Description Create or update a country. Action must be one of: `create`, `update`. ID should only be provided for `update`.
// @Tags Country
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Param data body models.CountryRequest true "Country payload" example({"action": "create", "country_code": "PL", "region_name": "North America", "country_name": "United States", "quote_currency_code": "USD", "regulator_name": "SEC", "fiat_label": "USD", "fiat_glyph": "$", "sec_tokenization_fee_percent": 0.1, "sec_trade_fee_percent": 0.05, "tokenization_application_fee_asset": "USDT", "legal_and_professional_fee_percent": 0.2, "rating_agency_fee_percent": 0.02, "vat_percent": 0.0, "sec_tokenization_fee_fixed": 100.0, "sec_trade_fee_fixed": 50.0, "legal_and_professional_fee_fixed": 200.0, "rating_agency_fee_fixed": 20.0})
// @Success 200 {object} models.StandardSuccessResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /country/save [post]
func SaveCountryHandler(walletDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		_, err := middleware.ExtractTokenMetadata(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized access"})
			return
		}

		var req models.CountryRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request format"})
			return
		}

		if req.Action != "create" && req.Action != "update" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid action: must be 'create' or 'update'"})
			return
		}

		if req.Action == "update" && req.CountryCode == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "country_code is required for update"})
			return
		}

		if err := usermetricsDB.SaveCountry(walletDB, req); err != nil {
			log.Println("[SAVE_COUNTRY] error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		message := "Country updated successfully"
		if req.Action == "create" {
			message = "Country created successfully"
		}

		serverResponse.JSON(c, http.StatusOK, message, nil, nil)
	}
}

// GetCountriesHandler fetches all countries
// @Summary Get countries
// @Description Get all countries
// @Tags Country
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Success 200 {array} models.CountryResponse
// @Failure 401 {object} map[string]string
// @Router /country/list [get]
func GetCountriesHandler(walletDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		_, err := middleware.ExtractTokenMetadata(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized access"})
			return
		}

		countries, err := usermetricsDB.ListCountries(walletDB)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, countries)
	}
}

// DeleteCountryHandler deletes a country
// @Summary Delete country
// @Description Delete a country by ID
// @Tags Country
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Param id path int true "Country ID"
// @Success 200 {object} models.StandardSuccessResponse
// @Failure 400 {object} map[string]string
// @Router /country/delete/{id} [delete]
func DeleteCountryHandler(walletDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		_, err := middleware.ExtractTokenMetadata(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized access"})
			return
		}

		id := c.Param("id")
		if err := usermetricsDB.DeleteCountry(walletDB, id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		serverResponse.JSON(c, http.StatusOK, "Country deleted successfully", nil, nil)
	}
}

// GetCountryByIDHandler returns a single country by ID
// @Summary Get country by ID
// @Description Fetch a single country using its ID
// @Tags Country
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Param id path int true "Country ID"
// @Success 200 {object} models.CountryResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /country/get/{id} [get]
func GetCountryByIDHandler(walletDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		_, err := middleware.ExtractTokenMetadata(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized access"})
			return
		}

		id := c.Param("id")
		country, err := usermetricsDB.GetCountryByID(walletDB, id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "country not found"})
			return
		}

		c.JSON(http.StatusOK, country)
	}
}
