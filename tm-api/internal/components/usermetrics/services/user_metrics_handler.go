package usermetrics

import (
	usermetricsDB "admin-panel-dashboard/internal/components/usermetrics/db"
	userServices "admin-panel-dashboard/internal/components/users/services"
	p2pErrors "admin-panel-dashboard/internal/errors"
	"admin-panel-dashboard/internal/middleware"
	"admin-panel-dashboard/internal/models"
	serverResponse "admin-panel-dashboard/internal/server/response"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// @Summary Get user metrics
// @Description Retrieves metrics related to user activity.
// @ID GetUserMetrics
// @Tags Metrics
// @Security JwtTokenAuth
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Success 200 {object} response.Data
// @Failure 400,401,500 {object} object
// @Router /users/metrics [get]
func GetUserMetrics(walletDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		ad, _ := middleware.ExtractTokenMetadata(c.Request)
		userID := ad.UserID

		userInfo, err := userServices.GetUser(userID, walletDB)
		if err != nil {
			log.Println("[METRICS] error for user:", userInfo.Username, "error: ", err)

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
		log.Println("checking user data: ", userInfo, userID)
		userMetrics, err := usermetricsDB.GetUserMetrics(walletDB)
		if err != nil {
			log.Println("[METRICS] error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
		serverResponse.JSON(c, http.StatusOK, "userMetrics fetched successfully", userMetrics, nil)
	}

}

// @Summary Get wallet distribution data
// @Description Retrieves data about the distribution of user wallets.
// @ID GetWalletDistribution
// @Tags Wallets
// @Security JwtTokenAuth
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Success 200 {object} response.Data
// @Failure 400,401,500 {object} object
// @Router /wallet/distribution [get]
func GetWalletDistribution(walletDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		ad, _ := middleware.ExtractTokenMetadata(c.Request)
		userID := ad.UserID

		userInfo, err := userServices.GetUser(userID, walletDB)
		if err != nil {
			log.Println("[METRICS] error for user:", userInfo.Username, "error: ", err)

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
		log.Println("checking user data: ", userInfo, userID)
		walletDistribution, err := usermetricsDB.GetWalletVersionDistributionData(walletDB)
		if err != nil {
			log.Println("[METRICS] error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
		serverResponse.JSON(c, http.StatusOK, "walletDistribution data fetched successfully", walletDistribution, nil)
	}

}

// @Summary Get top referrer count
// @Description Retrieves the count of top referrers.
// @ID GetReferrerCount
// @Tags Referrers
// @Security JwtTokenAuth
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Success 200 {object} response.Data
// @Failure 400,401,500 {object} object
// @Router /referrer/count [get]
func GetReferrerCount(walletDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		ad, _ := middleware.ExtractTokenMetadata(c.Request)
		userID := ad.UserID

		userInfo, err := userServices.GetUser(userID, walletDB)
		if err != nil {
			log.Println("[METRICS] error for user:", userInfo.Username, "error: ", err)

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
		log.Println("checking user data: ", userInfo, userID)
		topReferrer, err := usermetricsDB.GetTopReferrers(walletDB)
		if err != nil {
			log.Println("[METRICS] error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
		serverResponse.JSON(c, http.StatusOK, "topReferrer data fetched successfully", topReferrer, nil)
	}

}

// @Summary Get recent registrations
// @Description Retrieves information about recent user registrations with pagination and optional filtering by username, email, or phone number.
// @ID GetRecentRegistrations
// @Tags Registrations
// @Security JwtTokenAuth
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Param page query int false "Page number for pagination" default(1)
// @Param pageSize query int false "Number of items per page" default(10)
// @Param username query string false "Filter by username"
// @Param email query string false "Filter by user email"
// @Param phone query string false "Filter by user phone number"
// @Param first_name query string false "Filter by user first_name"
// @Param last_name query string false "Filter by user last name"
// @Param search query string false "General search across multiple fields"
// @Success 200 {object} response.Data
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /recent/registrations [get]
func GetRecentRegistrations(walletDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		ad, _ := middleware.ExtractTokenMetadata(c.Request)
		userID := ad.UserID

		userInfo, err := userServices.GetUser(userID, walletDB)
		if err != nil {
			log.Println("[METRICS] error for user:", userInfo.Username, "error: ", err)

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
		log.Println("checking user data: ", userInfo, userID)
		page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
		if err != nil || page < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page number"})
			return
		}

		pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
		if err != nil || pageSize <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page size"})
			return
		}
		usernameFilter := c.Query("username")
		emailFilter := c.Query("email")
		phoneFilter := c.Query("phone")
		firstName := c.Query("first_name")
		lastName := c.Query("last_name")
		search := c.Query("search")

		// Retrieve data with pagination and filtering
		req := models.UserRequestDTO{
			Page:      page,
			PageSize:  pageSize,
			Username:  usernameFilter,
			Email:     emailFilter,
			Phone:     phoneFilter,
			FirstName: firstName,
			LastName:  lastName,
			Search:    search,
		}

		recentRegistrations, total, err := usermetricsDB.GetRecentRegistrations(req, walletDB)
		if err != nil {
			log.Println("[METRICS] error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
		// Send the response
		// Send the paginated response
		user := models.UserResponse{
			Data:     recentRegistrations,
			Total:    total,
			Page:     page,
			PageSize: pageSize,
		}
		serverResponse.JSON(c, http.StatusOK, "recent registration list fetched successfully", user, nil)
	}

}

// @Summary Get country statistics
// @Description Retrieves statistics related to user distribution by country.
// @ID GetCountryStatistics
// @Tags Statistics
// @Security JwtTokenAuth
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Success 200 {object} response.Data
// @Failure 400,401,500 {object} object
// @Router /country/statistics [get]
func GetCountryStatistics(walletDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		ad, _ := middleware.ExtractTokenMetadata(c.Request)
		userID := ad.UserID

		userInfo, err := userServices.GetUser(userID, walletDB)
		if err != nil {
			log.Println("[METRICS] error for user:", userInfo.Username, "error: ", err)

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
		log.Println("checking user data: ", userInfo, userID)
		countryStatistics, err := usermetricsDB.CountUsersByCountry(walletDB)
		if err != nil {
			log.Println("[METRICS] error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
		// Send the response
		serverResponse.JSON(c, http.StatusOK, "countryStatistics fetched successfully", countryStatistics, nil)
	}

}

// @Summary Get P2P statistics
// @Description Retrieves statistics related to P2P metrics.
// @ID GetP2PStatistics
// @Tags P2P
// @Security JwtTokenAuth
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Success 200 {object} response.Data
// @Failure 400,401,500 {object} object
// // @Router /p2p/statistics [get]
func GetP2PStatistics(walletDB, p2pDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		ad, _ := middleware.ExtractTokenMetadata(c.Request)
		userID := ad.UserID

		userInfo, err := userServices.GetUser(userID, walletDB)
		if err != nil {
			log.Println("[METRICS] error for user:", userInfo.Username, "error: ", err)

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
		log.Println("checking user data: ", userInfo, userID)
		p2pMetrics, err := usermetricsDB.GetP2PMetrics(p2pDB)
		if err != nil {
			log.Println("[METRICS] error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		serverResponse.JSON(c, http.StatusOK, "p2pMetrics fetched successfully", p2pMetrics, nil)
	}

}

// GetP2PStatisticsNew retrieves statistics based on filters or custom date range.
// @Summary Get P2P statistics
// @Description Retrieves statistics related to P2P metrics based on the selected filter or custom range.
// @ID GetP2PStatisticsNew
// @Tags P2P
// @Security JwtTokenAuth
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Param start_date query string false "The start date for the custom range (format: YYYY-MM-DD)."
// @Param end_date query string false "The end date for the custom range (format: YYYY-MM-DD)."
// @Param filter query string false "The time filter to apply (e.g., today, yesterday, last_week, last_month, last_3_months, last_year, custom)." Enums(today, yesterday, last_week, last_month, last_3_months, last_year, custom)
// @Success 200 {object} Response
// @Failure 400,401,500 {object} object
// @Router /p2p/statistics [get]
func GetP2PStatisticsNew(walletDB, p2pDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract filter parameter and dates
		filter := c.DefaultQuery("filter", "last_week") // Default to "last_week"
		startDateParam := c.Query("start_date")
		endDateParam := c.Query("end_date")

		// Parse start_date and end_date with UTC timezone awareness
		var startDate, endDate *time.Time
		if startDateParam != "" {
			parsedStartDate, err := time.Parse("2006-01-02", startDateParam)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_date format. Use YYYY-MM-DD."})
				return
			}
			utcStartDate := parsedStartDate.UTC()
			startDate = &utcStartDate
		}
		if endDateParam != "" {
			parsedEndDate, err := time.Parse("2006-01-02", endDateParam)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_date format. Use YYYY-MM-DD."})
				return
			}
			utcEndDate := parsedEndDate.UTC()
			endDate = &utcEndDate
		}

		// Automatically set filter to "custom" if start_date and end_date are provided
		if startDate != nil && endDate != nil {
			filter = "custom"
		}

		// Extract user metadata from JWT
		ad, _ := middleware.ExtractTokenMetadata(c.Request)
		userID := ad.UserID

		// Fetch user info from database
		_, err := userServices.GetUser(userID, walletDB)
		if err != nil {
			log.Println("[METRICS] Error for user:", userID, "error:", err)
			if ex, ok := err.(p2pErrors.GenericError); ok {
				c.JSON(ex.HTTPCode(), ex.JSONError())
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}
			return
		}

		// Fetch P2P metrics using the filter
		p2pMetrics, err := usermetricsDB.GetP2PMetricsNew(p2pDB, filter, startDate, endDate)
		if err != nil {
			log.Println("[METRICS] Error fetching P2P metrics:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		// Send success response
		serverResponse.JSON(c, http.StatusOK, "P2P Metrics fetched successfully", p2pMetrics, nil)
	}
}

// @Summary Get list of trades on appeal
// @Description Retrieves the list of trades that are currently on appeal with pagination, filtering, and search.
// @ID GetAppealList
// @Tags Appeals
// @Security JwtTokenAuth
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Param page query int false "Page number for pagination" default(1)
// @Param pageSize query int false "Number of items per page" default(10)
// @Param appeal_id query string false "Filter by appeal ID"
// @Param customer_name query string false "Filter by customer name"
// @Param description query string false "Filter by description"
// @Param updated_on query string false "Filter by updated date (YYYY-MM-DD)"
// @Param status query string false "Filter by status e.g suspended, "", and "" "
// @Param search query string false "General search across multiple fields"
// @Success 200 {object} response.Data
// @Failure 400 {object} models.ErrorResponse "Invalid request parameters"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /appeal/list [get]
func GetAppealList(walletDB, p2pDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		ad, _ := middleware.ExtractTokenMetadata(c.Request)
		userID := ad.UserID

		userInfo, err := userServices.GetUser(userID, walletDB)
		if err != nil {
			log.Println("[METRICS] error for user:", userInfo.Username, "error: ", err)

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
		log.Println("checking user data: ", userInfo, userID)
		page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
		if err != nil || page < 1 {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid page number"})
			return
		}

		pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
		if err != nil || pageSize <= 0 {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid page size"})
			return
		}
		req := models.AppealListRequest{
			Page:         page,
			PageSize:     pageSize,
			AppealID:     c.Query("appeal_id"),
			CustomerName: c.Query("customer_name"),
			Description:  c.Query("description"),
			UpdatedOn:    c.Query("updated_on"),
			Status:       c.Query("status"),
			Search:       c.Query("search"),
		}
		// Retrieve data with pagination, filtering, and search
		appealList, total, err := usermetricsDB.GetTradesOnAppealList(req, p2pDB)
		if err != nil {
			log.Println("[METRICS] error:", err)
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal server error"})
			return
		}

		appealListData := models.AppealResponse{
			Data:     appealList,
			Total:    total,
			Page:     req.Page,
			PageSize: req.PageSize,
		}
		serverResponse.JSON(c, http.StatusOK, "Appeal List Data fetched successfully", appealListData, nil)
	}

}

// GetTradeStatistics @Summary Get Trade Statistics
// @Description Retrieve trade statistics for a user
// @Tags P2P
// @Accept  json
// @Produce  json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Success 200 {object} usermetrics.TradeStatistics
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /trades/statistics [get]
func GetTradeStatistics(walletDB, p2pDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		ad, _ := middleware.ExtractTokenMetadata(c.Request)
		userID := ad.UserID

		userInfo, err := userServices.GetUser(userID, walletDB)
		if err != nil {
			log.Println("[METRICS] error for user:", userInfo.Username, "error: ", err)

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
		log.Println("checking user data: ", userInfo, userID)

		// Call service layer
		stats, err := usermetricsDB.GetTradeStatistics(p2pDB)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, stats)
	}
}

// @Summary Get list of trades
// @Description Retrieves the list of trades for users.
// @ID GetTradeList
// @Tags P2P
// @Security JwtTokenAuth
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Param id query string false "Filter by ID"
// @Param order_type query string false "Filter by order type"
// @Param created_at query string false "Filter by creation date"
// @Param updated_at query string false "Filter by update date"
// @Param expires_at query string false "Filter by expiration date"
// @Param accepted_at query string false "Filter by acceptance date"
// @Param cancel_after query string false "Filter by cancel after date"
// @Param offer_id query string false "Filter by offer ID"
// @Param offer_type query string false "Filter by offer type"
// @Param offer_maker query string false "Filter by offer maker"
// @Param offer_maker_phone query string false "Filter by offer maker phone"
// @Param offer_maker_country_code query string false "Filter by offer maker country code"
// @Param offer_max_time_per_transaction query uint false "Filter by max time per transaction"
// @Param offer_taker query string false "Filter by offer taker"
// @Param offer_taker_phone query string false "Filter by offer taker phone"
// @Param offer_payment_method_id query string false "Filter by offer payment method ID"
// @Param offer_payment_channel_id query string false "Filter by offer payment channel ID"
// @Param offer_payment_method_name query string false "Filter by offer payment method name"
// @Param offer_payment_method_destination_account query string false "Filter by offer payment method destination account"
// @Param offer_payment_method_memo query string false "Filter by offer payment method memo"
// @Param offer_payment_method_bank_name query string false "Filter by offer payment method bank name"
// @Param offer_payment_method_account_opening_branch query string false "Filter by offer payment method account opening branch"
// @Param offer_currency_payment_method_id query string false "Filter by offer currency payment method ID"
// @Param offer_currency_payment_channel_id query string false "Filter by offer currency payment channel ID"
// @Param offer_currency_payment_method_name query string false "Filter by offer currency payment method name"
// @Param offer_currency_payment_method_destination_account query string false "Filter by offer currency payment method destination account"
// @Param offer_currency_payment_method_memo query string false "Filter by offer currency payment method memo"
// @Param offer_currency_payment_method_bank_name query string false "Filter by offer currency payment method bank name"
// @Param offer_currency_payment_method_account_opening_branch query string false "Filter by offer currency payment method account opening branch"
// @Param offer_currency_payment_method_country_code query string false "Filter by offer currency payment method country code"
// @Param offer_currency_payment_method_currency_id query string false "Filter by offer currency payment method currency ID"
// @Param offer_currency_id query string false "Filter by offer currency ID"
// @Param offer_asset_amount query float64 false "Filter by offer asset amount"
// @Param offer_asset_id query string false "Filter by offer asset ID"
// @Param offer_asset_price query float64 false "Filter by offer asset price"
// @Param offer_min_trade_amount query float64 false "Filter by offer min trade amount"
// @Param offer_max_trade_amount query float64 false "Filter by offer max trade amount"
// @Param offer_remark query string false "Filter by offer remark"
// @Param taker_payment_method_id query string false "Filter by taker payment method ID"
// @Param taker_payment_channel_id query string false "Filter by taker payment channel ID"
// @Param taker_payment_method_name query string false "Filter by taker payment method name"
// @Param taker_payment_method_destination_account query string false "Filter by taker payment method destination account"
// @Param taker_payment_method_memo query string false "Filter by taker payment method memo"
// @Param taker_payment_method_bank_name query string false "Filter by taker payment method bank name"
// @Param taker_payment_method_account_opening_branch query string false "Filter by taker payment method account opening branch"
// @Param taker_payment_method_country_code query string false "Filter by taker payment method country code"
// @Param taker_payment_method_currency_id query string false "Filter by taker payment method currency ID"
// @Param order_escrow_address query string false "Filter by order escrow address"
// @Param order_payment_memo query string false "Filter by order payment memo"
// @Param order_amount query float64 false "Filter by order amount"
// @Param order_maker_fee query float64 false "Filter by order maker fee"
// @Param order_taker_fee query float64 false "Filter by order taker fee"
// @Param order_escrow_transaction_id query string false "Filter by order escrow transaction ID"
// @Param order_asset_release_transaction_id query string false "Filter by order asset release transaction ID"
// @Param order_status_id query uint false "Filter by order status ID"
// @Param order_status query string false "Filter by order status"
// @Param dynamic_link query string false "Filter by dynamic link"
// @Param qr_code query string false "Filter by QR code"
// @Param fiat_deposit_transaction_id query string false "Filter by fiat deposit transaction ID"
// @Param page query int false "Page number for pagination" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Success 200 {object} response.Data
// @Failure 400 {object} models.ErrorResponse "Invalid request"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /orders/trade [get]
func GetTradeList(walletDB, p2pDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		ad, _ := middleware.ExtractTokenMetadata(c.Request)
		userID := ad.UserID

		userInfo, err := userServices.GetUser(userID, walletDB)
		if err != nil {
			log.Println("[METRICS] error for user:", userInfo.Username, "error: ", err)

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
		log.Println("checking user data: ", userInfo, userID)
		var req models.TradeListRequest
		if err := c.ShouldBindQuery(&req); err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid request parameters"})
			return
		}

		tradeList, total, err := usermetricsDB.GetAllTradesList(req, p2pDB)
		if err != nil {
			log.Println("[METRICS] error:", err)
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Internal server error"})
			return
		}

		serverResponse.JSON(c, http.StatusOK, "Orders List fetched successfully", gin.H{
			"data":     tradeList,
			"total":    total,
			"page":     req.Page,
			"pageSize": req.PageSize,
		}, nil)
	}
	//todo: how to get total completed trades from db, not just total trade
	//todo: also trading hours
	// how to compute total transactions since there are multiple currencies
	//trding volume graph
}

// @Summary Get list of trades where user is taker or maker
// @Description Retrieves the list of trades for users.
// @ID GetTradeListByUserName
// @Tags P2P
// @Security JwtTokenAuth
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Param username query string false "Filters tradelist by username and returns all trades that matches that username."
// @Param id query string false "Filter by ID"
// @Param order_type query string false "Filter by order type"
// @Param created_at query string false "Filter by creation date"
// @Param updated_at query string false "Filter by update date"
// @Param expires_at query string false "Filter by expiration date"
// @Param accepted_at query string false "Filter by acceptance date"
// @Param cancel_after query string false "Filter by cancel after date"
// @Param offer_id query string false "Filter by offer ID"
// @Param offer_type query string false "Filter by offer type"
// @Param offer_maker query string false "Filter by offer maker"
// @Param offer_maker_phone query string false "Filter by offer maker phone"
// @Param offer_maker_country_code query string false "Filter by offer maker country code"
// @Param offer_max_time_per_transaction query uint false "Filter by max time per transaction"
// @Param offer_taker query string false "Filter by offer taker"
// @Param offer_taker_phone query string false "Filter by offer taker phone"
// @Param offer_payment_method_id query string false "Filter by offer payment method ID"
// @Param offer_payment_channel_id query string false "Filter by offer payment channel ID"
// @Param offer_payment_method_name query string false "Filter by offer payment method name"
// @Param offer_payment_method_destination_account query string false "Filter by offer payment method destination account"
// @Param offer_payment_method_memo query string false "Filter by offer payment method memo"
// @Param offer_payment_method_bank_name query string false "Filter by offer payment method bank name"
// @Param offer_payment_method_account_opening_branch query string false "Filter by offer payment method account opening branch"
// @Param offer_currency_payment_method_id query string false "Filter by offer currency payment method ID"
// @Param offer_currency_payment_channel_id query string false "Filter by offer currency payment channel ID"
// @Param offer_currency_payment_method_name query string false "Filter by offer currency payment method name"
// @Param offer_currency_payment_method_destination_account query string false "Filter by offer currency payment method destination account"
// @Param offer_currency_payment_method_memo query string false "Filter by offer currency payment method memo"
// @Param offer_currency_payment_method_bank_name query string false "Filter by offer currency payment method bank name"
// @Param offer_currency_payment_method_account_opening_branch query string false "Filter by offer currency payment method account opening branch"
// @Param offer_currency_payment_method_country_code query string false "Filter by offer currency payment method country code"
// @Param offer_currency_payment_method_currency_id query string false "Filter by offer currency payment method currency ID"
// @Param offer_currency_id query string false "Filter by offer currency ID"
// @Param offer_asset_amount query float64 false "Filter by offer asset amount"
// @Param offer_asset_id query string false "Filter by offer asset ID"
// @Param offer_asset_price query float64 false "Filter by offer asset price"
// @Param offer_min_trade_amount query float64 false "Filter by offer min trade amount"
// @Param offer_max_trade_amount query float64 false "Filter by offer max trade amount"
// @Param offer_remark query string false "Filter by offer remark"
// @Param taker_payment_method_id query string false "Filter by taker payment method ID"
// @Param taker_payment_channel_id query string false "Filter by taker payment channel ID"
// @Param taker_payment_method_name query string false "Filter by taker payment method name"
// @Param taker_payment_method_destination_account query string false "Filter by taker payment method destination account"
// @Param taker_payment_method_memo query string false "Filter by taker payment method memo"
// @Param taker_payment_method_bank_name query string false "Filter by taker payment method bank name"
// @Param taker_payment_method_account_opening_branch query string false "Filter by taker payment method account opening branch"
// @Param taker_payment_method_country_code query string false "Filter by taker payment method country code"
// @Param taker_payment_method_currency_id query string false "Filter by taker payment method currency ID"
// @Param order_escrow_address query string false "Filter by order escrow address"
// @Param order_payment_memo query string false "Filter by order payment memo"
// @Param order_amount query float64 false "Filter by order amount"
// @Param order_maker_fee query float64 false "Filter by order maker fee"
// @Param order_taker_fee query float64 false "Filter by order taker fee"
// @Param order_escrow_transaction_id query string false "Filter by order escrow transaction ID"
// @Param order_asset_release_transaction_id query string false "Filter by order asset release transaction ID"
// @Param order_status_id query uint false "Filter by order status ID"
// @Param order_status query string false "Filter by order status"
// @Param dynamic_link query string false "Filter by dynamic link"
// @Param qr_code query string false "Filter by QR code"
// @Param fiat_deposit_transaction_id query string false "Filter by fiat deposit transaction ID"
// @Param page query int false "Page number for pagination" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Success 200 {object} response.Data
// @Failure 400 {object} models.ErrorResponse "Invalid request"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /orders/trade/list [get]
func GetTradeListByUserName(walletDB, p2pDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		ad, _ := middleware.ExtractTokenMetadata(c.Request)
		userID := ad.UserID

		userInfo, err := userServices.GetUser(userID, walletDB)
		if err != nil {
			log.Println("[METRICS] error for user:", userInfo.Username, "error: ", err)

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
		log.Println("checking user data: ", userInfo, userID)
		var req models.TradeListRequest
		if err := c.ShouldBindQuery(&req); err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid request parameters"})
			return
		}

		tradeList, total, err := usermetricsDB.GetAllTradesList(req, p2pDB)
		if err != nil {
			log.Println("[METRICS] error:", err)
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Internal server error"})
			return
		}

		serverResponse.JSON(c, http.StatusOK, "Orders List fetched successfully", gin.H{
			"data":     tradeList,
			"total":    total,
			"page":     req.Page,
			"pageSize": req.PageSize,
		}, nil)
	}
	//todo: how to get total completed trades from db, not just total trade
	//todo: also trading hours
	// how to compute total transactions since there are multiple currencies
	//trding volume graph
}

// @Summary Get list of users
// @Description Retrieves information about list of users with pagination and optional filtering by username, email, or phone number.
// @ID GetUserList
// @Tags Registrations
// @Security JwtTokenAuth
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Param page query int false "Page number for pagination" default(1)
// @Param pageSize query int false "Number of items per page" default(10)
// @Param username query string false "Filter by username"
// @Param email query string false "Filter by user email"
// @Param phone query string false "Filter by user phone number"
// @Param first_name query string false "Filter by user first_name"
// @Param last_name query string false "Filter by user last name"
// @Param city query string false "Filter by user city"
// @Param public_key query string false "Filter by public_key"
// @Param search query string false "General search across multiple fields"
// @Success 200 {object} response.Data
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /users [get]
func GetUserList(walletDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		ad, _ := middleware.ExtractTokenMetadata(c.Request)
		userID := ad.UserID

		userInfo, err := userServices.GetUser(userID, walletDB)
		if err != nil {
			log.Println("[METRICS] error for user:", userInfo.Username, "error: ", err)

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
		log.Println("checking user data: ", userInfo, userID)
		page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
		if err != nil || page < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page number"})
			return
		}

		pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
		if err != nil || pageSize <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page size"})
			return
		}
		usernameFilter := c.Query("username")
		emailFilter := c.Query("email")
		phoneFilter := c.Query("phone")
		firstName := c.Query("first_name")
		lastName := c.Query("last_name")
		search := c.Query("search")
		city := c.Query("city")
		publicKey := c.Query("public_key")

		// Retrieve data with pagination and filtering
		req := models.UserRequestDTO{
			Page:      page,
			PageSize:  pageSize,
			Username:  usernameFilter,
			Email:     emailFilter,
			Phone:     phoneFilter,
			FirstName: firstName,
			LastName:  lastName,
			City:      city,
			Search:    search,
			PublicKey: publicKey,
		}
		tradeList, total, err := usermetricsDB.GetUserList(req, walletDB)
		if err != nil {
			log.Println("[METRICS] error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		user := models.UserResponseData{
			Data:     tradeList,
			Total:    total,
			Page:     page,
			PageSize: pageSize,
		}
		serverResponse.JSON(c, http.StatusOK, "userList fetched successfully", user, nil)
	}
}

func GetUserProfileInfo(walletDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

	}
	// db, err := db.TrovoWalletDb()
	//if err != nil {
	//	return
	//}
	//ad, _ := middleware.ExtractTokenMetadata(c.Request)
	//userID := ad.UserID
	//
	//userInfo, err := userServices.GetUser(userID, db)
	//if err != nil {
	//	log.Println("[METRICS] error for user:", userInfo.Username, "error: ", err)
	//
	//	var ex p2pErrors.GenericError
	//	var ok bool
	//
	//	ex, ok = err.(p2pErrors.GenericError)
	//	var statusCode int = 0
	//	var response interface{}
	//
	//	if ok {
	//		statusCode = ex.HTTPCode()
	//		response = ex.JSONError()
	//	} else {
	//		statusCode = http.StatusBadRequest
	//		response = gin.H{"error": err.Error()}
	//	}
	//
	//	c.JSON(statusCode, response)
	//	return
	//}

	// userInf, err := usermetrics.GetUserProfileInfo()
	//if err != nil {
	//	log.Println("[METRICS] error:", err)
	//	c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	//	return
	//}
	//
	//// Send the response
	//c.JSON(http.StatusOK, userInf)
}

//// @Summary Get user wallets
//// @Description Retrieves the wallets information of the authenticated user.
//// @ID GetUserWallets
//// @Tags Users
//// @Security JwtTokenAuth
//// @Produce json
//// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Success 200 {object} response.Data
//// @Failure 400,401,500 {object} object
//// @Router /user/wallets [get]
// func GetUserWallets(c *gin.Context) {
//	walletDb, err := db.TrovoWalletDb()
//	if err != nil {
//		return
//	}
//	//ad, _ := middleware.ExtractTokenMetadata(c.Request)
//	//userID := ad.UserID
//	//
//	//userInfo, err := userServices.GetUser(userID, walletDb)
//	//if err != nil {
//	//	log.Println("[METRICS] error for user:", userInfo.Username, "error: ", err)
//	//
//	//	var ex p2pErrors.GenericError
//	//	var ok bool
//	//
//	//	ex, ok = err.(p2pErrors.GenericError)
//	//	var statusCode int = 0
//	//	var response interface{}
//	//
//	//	if ok {
//	//		statusCode = ex.HTTPCode()
//	//		response = ex.JSONError()
//	//	} else {
//	//		statusCode = http.StatusBadRequest
//	//		response = gin.H{"error": err.Error()}
//	//	}
//	//
//	//	c.JSON(statusCode, response)
//	//	return
//	//}
//
//	wallets, err := usermetrics.GetUserList()
//	if err != nil {
//		log.Println("[METRICS] error:", err)
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
//		return
//	}
//
//	// Send the response
//	c.JSON(http.StatusOK, wallets)
//}

// @Summary Get User Profile Info
// @Description Retrieves user profile info
// @ID GetUserProfile
// @Tags Registrations
// @Security JwtTokenAuth
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Param id query string false "Get by user ID"
// @Param email query string false "Get by user email"
// @Param username query string false "Get by username"
// @Success 200 {object} response.Data
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /users/profile [get]
func GetUserProfile(walletDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		ad, _ := middleware.ExtractTokenMetadata(c.Request)
		userID := ad.UserID

		userInfo, err := userServices.GetUser(userID, walletDB)
		if err != nil {
			log.Println("[METRICS] error for user:", userInfo.Username, "error: ", err)

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

		usernameFilter := c.Query("username")
		emailFilter := c.Query("email")
		id := c.Query("id")
		// Retrieve data with pagination and filtering
		req := models.UserProfile{
			ID:       id,
			Username: usernameFilter,
			Email:    emailFilter,
		}
		userInf, err := usermetricsDB.GetUserProfileInfo(req, walletDB)
		if err != nil {
			log.Println("[METRICS] error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
		req.UserID = userInf.ID
		wallets, err := usermetricsDB.GetUserWallets(req, walletDB)
		if err != nil {
			log.Println("[METRICS] error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		responseData := models.UserInfoResponse{
			UserInfo: userInf,
			Wallets:  wallets,
		}

		serverResponse.JSON(c, http.StatusOK, "user Profile fetched successfully", responseData, nil)
	}
}
