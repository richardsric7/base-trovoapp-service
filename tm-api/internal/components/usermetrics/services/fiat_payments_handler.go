package usermetrics

import (
	usermetricsDB "admin-panel-dashboard/internal/components/usermetrics/db"
	userServices "admin-panel-dashboard/internal/components/users/services"
	"admin-panel-dashboard/internal/middleware"
	"admin-panel-dashboard/internal/models"
	serverResponse "admin-panel-dashboard/internal/server/response"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// @Summary List fiat payment records
// @Description Retrieves paginated fiat payment records aggregated from fiat_payment and fiat_payment_invoice tables.
// @ID GetFiatPayments
// @Tags FiatPayments
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Param record_type query string false "Filter by record type" Enums(payment,invoice)
// @Param service_provider query string false "Filter by service provider"
// @Param username query string false "Filter by username"
// @Param payment_type query string false "Filter by payment type"
// @Param status query string false "Filter by invoice status"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} models.PaginatedFiatPaymentResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /fiat/payments [get]
func GetFiatPayments(walletDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract user info from token
		ad, _ := middleware.ExtractTokenMetadata(c.Request)
		userID := ad.UserID

		userInfo, err := userServices.GetUser(userID, walletDB)
		if err != nil {
			log.Println("[FIAT_PAYMENTS] error for user:", userInfo.Username, "error:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
		if err != nil || page < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page number"})
			return
		}

		pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "10"))
		if err != nil || pageSize < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page size"})
			return
		}

		recordType := c.Query("record_type")
		if recordType != "" && recordType != models.FiatRecordTypePayment && recordType != models.FiatRecordTypeInvoice {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid record_type. Must be either 'payment' or 'invoice'"})
			return
		}

		filter := models.FiatPaymentListFilter{
			RecordType:      recordType,
			ServiceProvider: c.Query("service_provider"),
			Username:        c.Query("username"),
			PaymentType:     c.Query("payment_type"),
			Status:          c.Query("status"),
			Page:            page,
			PageSize:        pageSize,
		}

		data, err := usermetricsDB.GetPaginatedFiatPaymentRecords(walletDB, filter)
		if err != nil {
			log.Println("[FIAT_PAYMENTS] error fetching records:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch fiat payment records"})
			return
		}

		serverResponse.JSON(c, http.StatusOK, "fiat payment records fetched successfully", data, nil)
	}
}
