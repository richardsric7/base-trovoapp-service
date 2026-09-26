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

// checkAdminAuth validates the calling admin's JWT-derived userID against
// the wallet users table, writing the appropriate error response itself
// on failure. Every P2P admin handler below needs this same check (it
// used to be repeated verbatim in each one).
func checkAdminAuth(c *gin.Context, walletDB *gorm.DB) bool {
	ad, _ := middleware.ExtractTokenMetadata(c.Request)
	if _, err := userServices.GetUser(ad.UserID, walletDB); err != nil {
		log.Println("[METRICS] error for user:", ad.UserID, "error: ", err)
		if ex, ok := err.(p2pErrors.GenericError); ok {
			c.JSON(ex.HTTPCode(), ex.JSONError())
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return false
	}
	return true
}

// @Summary Get P2P statistics
// @Description Retrieves aggregate P2P marketplace statistics (offers, orders, disputes) from the P2P module.
// @ID GetP2PStatistics
// @Tags P2P
// @Security JwtTokenAuth
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Success 200 {object} response.Data
// @Failure 400,401,500 {object} object
// @Router /p2p/statistics [get]
func GetP2PStatistics(walletDB, p2pDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !checkAdminAuth(c, walletDB) {
			return
		}
		stats, err := usermetricsDB.GetP2PStatistics(p2pDB)
		if err != nil {
			log.Println("[METRICS] error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
		serverResponse.JSON(c, http.StatusOK, "P2P statistics fetched successfully", stats, nil)
	}
}

// GetTradeStatistics @Summary Get Trade Statistics
// @Description Retrieve P2P trade statistics (completion rate, top traders, recent trades) from the P2P module.
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
		if !checkAdminAuth(c, walletDB) {
			return
		}
		stats, err := usermetricsDB.GetTradeStatistics(p2pDB)
		if err != nil {
			log.Println("[METRICS] error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		serverResponse.JSON(c, http.StatusOK, "Trade statistics fetched successfully", stats, nil)
	}
}

// @Summary Get list of P2P trades (orders)
// @Description Retrieves a paginated, filterable list of P2P orders from the P2P module.
// @ID GetTradeList
// @Tags P2P
// @Security JwtTokenAuth
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Param page query int false "Page number for pagination" default(1)
// @Param pageSize query int false "Number of items per page" default(10)
// @Param offerType query string false "Filter by offer type (BUY or SELL)"
// @Param status query string false "Filter by order status"
// @Param username query string false "Filter by merchant or customer username"
// @Param createdAt query string false "Filter by creation date (YYYY-MM-DD)"
// @Success 200 {object} response.Data
// @Failure 400 {object} models.ErrorResponse "Invalid request"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /orders/trade [get]
func GetTradeList(walletDB, p2pDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !checkAdminAuth(c, walletDB) {
			return
		}
		page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
		if err != nil || page < 1 {
			page = 1
		}
		pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
		if err != nil || pageSize <= 0 {
			pageSize = 10
		}
		req := usermetricsDB.TradeListRequest{
			Page:      page,
			PageSize:  pageSize,
			OfferType: c.Query("offerType"),
			Status:    c.Query("status"),
			Username:  c.Query("username"),
			CreatedAt: c.Query("createdAt"),
		}
		tradeList, total, err := usermetricsDB.GetTradeList(p2pDB, req)
		if err != nil {
			log.Println("[METRICS] error:", err)
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Internal server error"})
			return
		}
		serverResponse.JSON(c, http.StatusOK, "Orders list fetched successfully", gin.H{
			"data":     tradeList,
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
		}, nil)
	}
}

// @Summary Get list of P2P users
// @Description Retrieves wallet users with their P2P trading performance (merchant/customer) from the P2P module.
// @ID GetP2PUsers
// @Tags P2P
// @Security JwtTokenAuth
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Param page query int false "Page number for pagination" default(1)
// @Param pageSize query int false "Number of items per page" default(10)
// @Param username query string false "Filter by username"
// @Param email query string false "Filter by email"
// @Param phone query string false "Filter by phone"
// @Param search query string false "General search across username/email/phone"
// @Success 200 {object} response.Data
// @Failure 400,401,500 {object} object
// @Router /p2p/users [get]
func GetP2PUsers(walletDB, p2pDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !checkAdminAuth(c, walletDB) {
			return
		}
		page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
		if err != nil || page < 1 {
			page = 1
		}
		pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
		if err != nil || pageSize <= 0 {
			pageSize = 10
		}
		req := usermetricsDB.P2PUserListRequest{
			Page:     page,
			PageSize: pageSize,
			Username: c.Query("username"),
			Email:    c.Query("email"),
			Phone:    c.Query("phone"),
			Search:   c.Query("search"),
		}
		users, total, err := usermetricsDB.GetP2PUsers(walletDB, p2pDB, req)
		if err != nil {
			log.Println("[METRICS] error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
		serverResponse.JSON(c, http.StatusOK, "P2P users fetched successfully", gin.H{
			"data":     users,
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
		}, nil)
	}
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
// @Param address query string false "Filter by address"
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
		address := c.Query("address")

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
			Address:   address,
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
