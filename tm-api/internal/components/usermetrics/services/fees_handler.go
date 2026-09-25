package usermetrics

import (
	usermetricsDB "admin-panel-dashboard/internal/components/usermetrics/db"
	"admin-panel-dashboard/internal/server/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// @Summary Get all fee collections
// @Description Retrieves a list of all fee collections.
// @ID GetFeeCollections
// @Tags Fees
// @Security JwtTokenAuth
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Success 200 {object} response.Data{data=[]usermetricsDB.FeeCollection}
// @Failure 500 {object} map[string]string
// @Router /fee/collections [get]
func GetFeeCollectionsHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		fees, err := usermetricsDB.GetAllFeeCollections(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		response.JSON(c, http.StatusOK, "Fee collections fetched successfully", fees, nil)
	}
}

// @Summary Get all service link service fees
// @Description Retrieves a list of all service link service fees.
// @ID GetServiceLinkServiceFees
// @Tags Fees
// @Security JwtTokenAuth
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Success 200 {object} response.Data{data=[]usermetricsDB.ServiceLinkServiceFee}
// @Failure 500 {object} map[string]string
// @Router /fee/configs [get]
func GetServiceLinkServiceFeesHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		fees, err := usermetricsDB.GetAllServiceLinkServiceFees(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		response.JSON(c, http.StatusOK, "Service fees fetched successfully", fees, nil)
	}
}

// @Summary Get service link service fee by ID
// @Description Retrieves a service link service fee by its Service Link ID.
// @ID GetServiceLinkServiceFeeByID
// @Tags Fees
// @Security JwtTokenAuth
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Param id path string true "Service Link ID"
// @Success 200 {object} response.Data{data=usermetricsDB.ServiceLinkServiceFee}
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /fee/configs/{id} [get]
func GetServiceLinkServiceFeeByIDHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		fee, err := usermetricsDB.GetServiceLinkServiceFeeByID(db, id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Service fee not found"})
			return
		}
		response.JSON(c, http.StatusOK, "Service fee fetched successfully", fee, nil)
	}
}

// SaveServiceLinkServiceFeeRequest represents the request payload for saving a service fee
type SaveServiceLinkServiceFeeRequest struct {
	ServiceLinkID string  `json:"service_link_id" binding:"required"`
	PaymentFee    float64 `json:"payment_fee"`
	SwapFee       float64 `json:"swap_fee"`
	SubwalletFee  float64 `json:"subwallet_fee"`
}

// @Summary Save service link service fee
// @Description Creates or updates a service link service fee.
// @ID SaveServiceLinkServiceFee
// @Tags Fees
// @Security JwtTokenAuth
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Param fee body SaveServiceLinkServiceFeeRequest true "Service Fee Data"
// @Success 200 {object} response.Data{data=usermetricsDB.ServiceLinkServiceFee}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /fee/configs [post]
func SaveServiceLinkServiceFeeHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req SaveServiceLinkServiceFeeRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		fee := usermetricsDB.ServiceLinkServiceFee{
			ServiceLinkID: req.ServiceLinkID,
			PaymentFee:    req.PaymentFee,
			SwapFee:       req.SwapFee,
			SubwalletFee:  req.SubwalletFee,
		}

		if err := usermetricsDB.SaveServiceLinkServiceFee(db, &fee); err != nil {
			if err.Error() == "service link not found" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "the provided service_link_id is invalid"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		response.JSON(c, http.StatusOK, "Service fee saved successfully", fee, nil)
	}
}

// @Summary Delete service link service fee
// @Description Deletes a service link service fee by its ID.
// @ID DeleteServiceLinkServiceFee
// @Tags Fees
// @Security JwtTokenAuth
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Param id path string true "Service Link ID"
// @Success 200 {object} response.Data
// @Failure 500 {object} map[string]string
// @Router /fee/configs/{id} [delete]
func DeleteServiceLinkServiceFeeHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if err := usermetricsDB.DeleteServiceLinkServiceFee(db, id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		response.JSON(c, http.StatusOK, "Service fee deleted successfully", nil, nil)
	}
}
