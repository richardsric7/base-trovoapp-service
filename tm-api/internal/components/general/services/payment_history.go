package server

import (
	userServices "admin-panel-dashboard/internal/components/users/services"
	p2pErrors "admin-panel-dashboard/internal/errors"
	"admin-panel-dashboard/internal/middleware"
	"admin-panel-dashboard/internal/models"
	serverResponse "admin-panel-dashboard/internal/server/response"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// @Summary Get payment history
// @Description Retrieves the payment history with pagination, filtering, and search functionality.
// @ID GetPaymentHistory
// @Tags Payments
// @Security JwtTokenAuth
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Param page query int false "Page number for pagination" default(1)
// @Param pageSize query int false "Number of items per page" default(10)
// @Param transactionType query string false "Filter by transaction type (e.g., credit, debit)"
// @Param transactionDate query string false "Filter by transaction date (YYYY-MM-DD)"
// @Param from query string false "Filter by sender alias or name"
// @Param to query string false "Filter by recipient alias or name"
// @Param memo query string false "Filter by transaction memo"
// @Param contractAddress query string false "Filter by asset issuer"
// @Param assetCode query string false "Filter by asset code"
// @Param search query string false "Search across multiple fields (e.g., from, to, transaction ID)"
// @Success 200 {object} response.Data
// @Failure 400 {object} models.ErrorResponse "Invalid request parameters"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /payment/history [get]
func GetPaymentHistory(walletDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract user metadata and verify
		ad, _ := middleware.ExtractTokenMetadata(c.Request)
		userID := ad.UserID

		_, err := userServices.GetUser(userID, walletDB)
		if err != nil {
			log.Printf("[ERROR] user %s retrieval failed: %v", userID, err)
			respondWithError(c, err)
			return
		}

		// Validate pagination parameters
		page, pageSize, err := validatePagination(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
			return
		}

		// Construct request object for filtering
		req := models.PaymentHistoryRequest{
			Page:            page,
			PageSize:        pageSize,
			TransactionType: c.Query("transactionType"),
			TransactionDate: c.Query("transactionDate"),
			From:            c.Query("from"),
			To:              c.Query("to"),
			Memo:            c.Query("memo"),
			ContractAddress: c.Query("contractAddress"),
			AssetCode:       c.Query("assetCode"),
			Search:          c.Query("search"),
		}

		// Retrieve payment history
		paymentHistoryList, total, err := GetPaymentHistoryList(req, walletDB)
		if err != nil {
			log.Printf("[ERROR] payment history retrieval failed: %v", err)
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal server error"})
			return
		}

		// Respond with data
		paymentHistoryResponse := models.PaymentHistoryResponse{
			Data:     paymentHistoryList,
			Total:    total,
			Page:     req.Page,
			PageSize: req.PageSize,
		}
		serverResponse.JSON(c, http.StatusOK, "payment history fetched successfully", paymentHistoryResponse, nil)
	}
}

// validatePagination extracts and validates pagination parameters from the request context.
func validatePagination(c *gin.Context) (int, int, error) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		return 0, 0, fmt.Errorf("invalid page number")
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if err != nil || pageSize <= 0 {
		return 0, 0, fmt.Errorf("invalid page size")
	}

	return page, pageSize, nil
}

// respondWithError handles sending a structured error response.
func respondWithError(c *gin.Context, err error) {
	var ex p2pErrors.GenericError
	var ok bool

	if ex, ok = err.(p2pErrors.GenericError); ok {
		// If the error is a GenericError, use its HTTP status code and formatted message
		c.JSON(ex.HTTPCode(), ex.JSONError())
	} else {
		// Fallback to a generic 400 Bad Request if no specific error type is recognized
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func GetPaymentHistoryList(req models.PaymentHistoryRequest, db *gorm.DB) ([]models.PaymentHistoryJSON, int, error) {
	var paymentHistory []models.PaymentHistory
	var total int64

	query := db.Model(&models.PaymentHistory{})

	if req.TransactionType != "" {
		query = query.Where("transaction_type = ?", req.TransactionType)
	}
	if req.TransactionDate != "" {
		query = query.Where("DATE(transaction_date) = ?", req.TransactionDate)
	}
	if req.From != "" {
		query = query.Where(`"from" LIKE ?`, "%"+req.From+"%")
	}
	if req.To != "" {
		query = query.Where(`"to" LIKE ?`, "%"+req.To+"%")
	}
	if req.Memo != "" {
		query = query.Where("memo LIKE ?", "%"+req.Memo+"%")
	}
	if req.ContractAddress != "" {
		query = query.Where("contract_address LIKE ?", "%"+req.ContractAddress+"%")
	}
	if req.AssetCode != "" {
		query = query.Where("asset_code = ?", req.AssetCode)
	}
	if req.Search != "" {
		query = query.Where("from LIKE ? OR to LIKE ? OR transaction_id LIKE ?", "%"+req.Search+"%", "%"+req.Search+"%", "%"+req.Search+"%")
	}

	if req.Search != "" {
		// Using GORM's Or() method for cleaner multi-column search
		search := "%" + req.Search + "%"
		query = query.Where(
			db.Or(
				// db.Where("\"from\" ILIKE ?", search),
				//db.Where("\"to\" ILIKE ?", search),
				db.Where("transaction_id ILIKE ?", search),
				db.Where("memo ILIKE ?", search),
				db.Where("contract_address ILIKE ?", search),
				db.Where("asset_code ILIKE ?", search),
				db.Where("CAST(amount AS TEXT) ILIKE ?", search),
			),
		)
	}

	query.Count(&total)

	err := query.
		Offset((req.Page - 1) * req.PageSize).
		Limit(req.PageSize).
		Order("transaction_date DESC").
		Find(&paymentHistory).Error

	if err != nil {
		return nil, 0, err
	}

	// Transform data into JSON format
	paymentHistoryJSON := []models.PaymentHistoryJSON{}
	for _, ph := range paymentHistory {
		paymentHistoryJSON = append(paymentHistoryJSON, models.PaymentHistoryJSON{
			TransactionDate: ph.TransactionDate,
			TransactionType: ph.TransactionType,
			From:            coalesce(ph.From),
			FromAddress:     ph.FromAddress,
			To:              coalesce(ph.To),
			ToAddress:       ph.ToAddress,
			Memo:            coalesce(ph.Memo),
			ContractAddress: coalesce(ph.ContractAddress),
			AssetCode:       ph.AssetCode,
			Amount:          ph.Amount,
			TransactionID:   ph.TransactionID,
		})
	}

	return paymentHistoryJSON, int(total), nil
}

// coalesce handles nil pointers
func coalesce(input *string) string {
	if input == nil {
		return ""
	}
	return *input
}
