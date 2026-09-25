package usermetrics

import (
	userServices "admin-panel-dashboard/internal/components/users/services"
	p2pErrors "admin-panel-dashboard/internal/errors"
	"admin-panel-dashboard/internal/middleware"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetTokenizedAssetStatisticsHandler retrieves aggregated statistics for tokenized assets
// @Summary Get tokenized asset statistics
// @Description Retrieves aggregated statistics including counts and total values for all, approved, pending, and rejected tokenized assets
// @Tags Tokenization
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tokenization/statistics [get]
func GetTokenizedAssetStatisticsHandler(walletDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		ad, err := middleware.ExtractTokenMetadata(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		userID := ad.UserID
		userInfo, err := userServices.GetUser(userID, walletDB)
		if err != nil {
			log.Println("[TOKENIZATION_STATISTICS] error for user:", userInfo.Username, "error:", err)
			var ex p2pErrors.GenericError
			var ok bool
			ex, ok = err.(p2pErrors.GenericError)
			var statusCode int
			var response interface{}
			if ok {
				statusCode = ex.HTTPCode()
				response = ex.JSONError()
			} else {
				statusCode = http.StatusBadRequest
				response = gin.H{"error": err.Error()}
			}
			c.JSON(statusCode, response)
			return
		}

		// Call trovo-manager tokenization stat endpoint
		result, err := makeRequestWithRaw(http.MethodGet, "/tokenization/stat", c.GetHeader("Authorization"), nil)
		if err != nil {
			log.Printf("[TOKENIZATION_STATISTICS] error fetching statistics: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Data(http.StatusOK, "application/json", result)
	}
}
