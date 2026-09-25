package usermetrics

import (
	"admin-panel-dashboard/internal/components/accesslog"
	userMetricServices "admin-panel-dashboard/internal/components/usermetrics/services"
	"admin-panel-dashboard/internal/middleware"
	"admin-panel-dashboard/internal/models"
	serverModels "admin-panel-dashboard/internal/server/models"

	"github.com/gin-gonic/gin"
)

// Init initializes
func Init(router *gin.Engine, s *serverModels.Server) {

	// login callback url

	apiV1 := router.Group("/api/v1")

	apiV1.GET("/users/metrics", middleware.JwtTokenAuthMiddleware(s.AdminDB), userMetricServices.GetUserMetrics(s.TrovoWalletDB))
	apiV1.GET("/wallet/distribution", middleware.JwtTokenAuthMiddleware(s.AdminDB), userMetricServices.GetWalletDistribution(s.TrovoWalletDB))
	apiV1.GET("/referrer/count", middleware.JwtTokenAuthMiddleware(s.AdminDB), userMetricServices.GetReferrerCount(s.TrovoWalletDB))
	apiV1.GET("/recent/registrations", middleware.JwtTokenAuthMiddleware(s.AdminDB), userMetricServices.GetRecentRegistrations(s.TrovoWalletDB))
	apiV1.GET("/country/statistics", middleware.JwtTokenAuthMiddleware(s.AdminDB), userMetricServices.GetCountryStatistics(s.TrovoWalletDB))
	apiV1.GET("/p2p/statistics", middleware.JwtTokenAuthMiddleware(s.AdminDB), userMetricServices.GetP2PStatisticsNew(s.TrovoWalletDB, s.P2P))
	apiV1.GET("/appeal/list", middleware.JwtTokenAuthMiddleware(s.AdminDB), userMetricServices.GetAppealList(s.TrovoWalletDB, s.P2P))
	apiV1.GET("/trades/statistics", middleware.JwtTokenAuthMiddleware(s.AdminDB), userMetricServices.GetTradeStatistics(s.TrovoWalletDB, s.P2P))
	apiV1.GET("/orders/trade", middleware.JwtTokenAuthMiddleware(s.AdminDB), userMetricServices.GetTradeList(s.TrovoWalletDB, s.P2P))
	apiV1.GET("/orders/trade/list", middleware.JwtTokenAuthMiddleware(s.AdminDB), userMetricServices.GetTradeListByUserName(s.TrovoWalletDB, s.P2P))
	apiV1.GET("/users", middleware.JwtTokenAuthMiddleware(s.AdminDB), userMetricServices.GetUserList(s.TrovoWalletDB))
	apiV1.GET("/users/profile", middleware.JwtTokenAuthMiddleware(s.AdminDB), userMetricServices.GetUserProfile(s.TrovoWalletDB))
	apiV1.GET("/wallet-balances/:walletAddress", middleware.JwtTokenAuthMiddleware(s.AdminDB), userMetricServices.GetWalletBalances(s.TrovoWalletDB))
	apiV1.GET("/fiat/payments", middleware.JwtTokenAuthMiddleware(s.AdminDB), userMetricServices.GetFiatPayments(s.TrovoWalletDB))

	// Fee Management endpoints
	apiV1.GET("/fee/collections", middleware.JwtTokenAuthMiddleware(s.AdminDB), userMetricServices.GetFeeCollectionsHandler(s.TrovoWalletDB))
	apiV1.GET("/fee/configs", middleware.JwtTokenAuthMiddleware(s.AdminDB), userMetricServices.GetServiceLinkServiceFeesHandler(s.TrovoWalletDB))
	apiV1.GET("/fee/configs/:id", middleware.JwtTokenAuthMiddleware(s.AdminDB), userMetricServices.GetServiceLinkServiceFeeByIDHandler(s.TrovoWalletDB))
	apiV1.POST("/fee/configs", middleware.JwtTokenAuthMiddleware(s.AdminDB), accesslog.Audit(s.AdminDB, models.EventFeeConfigChange, models.AccessCategoryConfig, accesslog.BodyField("service_link_id")), userMetricServices.SaveServiceLinkServiceFeeHandler(s.TrovoWalletDB))
	apiV1.DELETE("/fee/configs/:id", middleware.JwtTokenAuthMiddleware(s.AdminDB), accesslog.Audit(s.AdminDB, models.EventFeeConfigDelete, models.AccessCategoryConfig, accesslog.Param("id")), userMetricServices.DeleteServiceLinkServiceFeeHandler(s.TrovoWalletDB))

	// Service Link endpoints
	apiV1.GET("/service-links", middleware.JwtTokenAuthMiddleware(s.AdminDB), userMetricServices.GetServiceLinksHandler(s.TrovoWalletDB))

	// Tokenization endpoints
	// apiV1.GET("/public/tokenization", middleware.JwtTokenAuthMiddleware(s.AdminDB), userMetricServices.GetPublicTokenization(s.TrovoWalletDB))
	//apiV1.GET("/tokenization", middleware.JwtTokenAuthMiddleware(s.AdminDB), userMetricServices.GetTokenization(s.TrovoWalletDB))
	apiV1.GET("/tokenization", middleware.JwtTokenAuthMiddleware(s.AdminDB), userMetricServices.GetTrovoManagerTokenizationWithRaw(s.TrovoWalletDB))
	apiV1.GET("/tokenization/detail/:tokenizedAssetID", middleware.JwtTokenAuthMiddleware(s.AdminDB), userMetricServices.GetTrovoManagerTokenizationDetailWithRaw(s.TrovoWalletDB))
	apiV1.GET("/tokenization/list", middleware.JwtTokenAuthMiddleware(s.AdminDB), userMetricServices.GetTrovoManagerTokenizationListWithRaw(s.TrovoWalletDB))
	apiV1.GET("/tokenization/statistics", middleware.JwtTokenAuthMiddleware(s.AdminDB), userMetricServices.GetTokenizedAssetStatisticsHandler(s.TrovoWalletDB))
	// apiV1.P("/tokenization", middleware.JwtTokenAuthMiddleware(s.AdminDB), userMetricServices.SubmitTokenization(s.TrovoWalletDB))
	apiV1.PUT("/tokenization/update/:tokenizedAssetID", middleware.JwtTokenAuthMiddleware(s.AdminDB), userMetricServices.UpdateTrovoManagerTokenizationWithRaw(s.TrovoWalletDB))
	apiV1.PUT("/tokenization/vet/:tokenizedAssetID", middleware.JwtTokenAuthMiddleware(s.AdminDB), accesslog.Audit(s.AdminDB, models.EventTokenizationVet, models.AccessCategoryAsset, accesslog.Param("tokenizedAssetID")), userMetricServices.VetTrovoManagerTokenizationWithRaw(s.TrovoWalletDB))
	apiV1.PUT("/tokenization/salesdate/:tokenizedAssetID", middleware.JwtTokenAuthMiddleware(s.AdminDB), userMetricServices.UpdateTrovoManagerTokenizationSalesDatesWithRaw(s.TrovoWalletDB))
	apiV1.POST("/tokenization/fee/:tokenizedAssetID", middleware.JwtTokenAuthMiddleware(s.AdminDB), userMetricServices.AcknowledgeTokenizationFeePaymentWithRaw(s.TrovoWalletDB))
	apiV1.POST("/tokenization/failed/:tokenizedAssetID", middleware.JwtTokenAuthMiddleware(s.AdminDB), accesslog.Audit(s.AdminDB, models.EventTokenizationFail, models.AccessCategoryAsset, accesslog.Param("tokenizedAssetID")), userMetricServices.FailDueDiligenceWithRaw(s.TrovoWalletDB))
	apiV1.POST("/tokenization/mint/:tokenizedAssetID", middleware.JwtTokenAuthMiddleware(s.AdminDB), accesslog.Audit(s.AdminDB, models.EventTokenizationMint, models.AccessCategoryAsset, accesslog.Param("tokenizedAssetID")), userMetricServices.ConfirmAndMintTokenizationWithRaw(s.TrovoWalletDB))
	// apiV1.PUT("/tokenization/document", middleware.JwtTokenAuthMiddleware(s.AdminDB), userMetricServices.UploadTokenizationDocument(s.TrovoWalletDB))
	apiV1.PUT("/tokenization/document", middleware.JwtTokenAuthMiddleware(s.AdminDB), userMetricServices.UploadTrovoManagerTokenizationDocumentWithRaw(s.TrovoWalletDB))
	apiV1.PUT("/tokenization/logo/:tokenizedAssetID", middleware.JwtTokenAuthMiddleware(s.AdminDB), userMetricServices.UploadTrovoManagerTokenizationLogo(s.TrovoWalletDB))
	apiV1.DELETE("/tokenization/document/:documentID", middleware.JwtTokenAuthMiddleware(s.AdminDB), userMetricServices.DeleteTrovoManagerTokenizationDocumentWithRaw(s.TrovoWalletDB))

	// P2P USERS MANAGEMENT
	apiV1.GET("/p2p/users", middleware.JwtTokenAuthMiddleware(s.AdminDB), userMetricServices.GetP2PUserList(s.TrovoWalletDB, s.P2P))
	apiV1.GET("p2p/transaction/history", middleware.JwtTokenAuthMiddleware(s.AdminDB), userMetricServices.GetTransactionHistory(s.P2P))

	// partners
	apiV1.POST("/partners/save", middleware.JwtTokenAuthMiddleware(s.AdminDB), userMetricServices.SavePartnerHandler(s.TrovoWalletDB))
	apiV1.DELETE("/partners/delete", middleware.JwtTokenAuthMiddleware(s.AdminDB), userMetricServices.DeletePartnerHandler(s.TrovoWalletDB))
	apiV1.GET("/partners/list", middleware.JwtTokenAuthMiddleware(s.AdminDB), userMetricServices.ListPartnersHandler(s.TrovoWalletDB))
	apiV1.GET("/partners/one", middleware.JwtTokenAuthMiddleware(s.AdminDB), userMetricServices.GetPartnerHandler(s.TrovoWalletDB))

	// country management endpoints
	apiV1.POST("/country/save", middleware.JwtTokenAuthMiddleware(s.AdminDB), userMetricServices.SaveCountryHandler(s.TrovoWalletDB))
	apiV1.GET("/country/list", middleware.JwtTokenAuthMiddleware(s.AdminDB), userMetricServices.GetCountriesHandler(s.TrovoWalletDB))
	apiV1.DELETE("/country/delete/:id", middleware.JwtTokenAuthMiddleware(s.AdminDB), userMetricServices.DeleteCountryHandler(s.TrovoWalletDB))
	apiV1.GET("/country/get/:id", middleware.JwtTokenAuthMiddleware(s.AdminDB), userMetricServices.GetCountryByIDHandler(s.TrovoWalletDB))

	// Faucet Config endpoints
	apiV1.POST("/faucet/configs", middleware.JwtTokenAuthMiddleware(s.AdminDB), userMetricServices.SaveFaucetConfigHandler(s.TrovoWalletDB))
	apiV1.GET("/faucet/configs", middleware.JwtTokenAuthMiddleware(s.AdminDB), userMetricServices.GetAllFaucetConfigsHandler(s.TrovoWalletDB))
	apiV1.GET("/faucet/configs/:id", middleware.JwtTokenAuthMiddleware(s.AdminDB), userMetricServices.GetFaucetConfigByIDHandler(s.TrovoWalletDB))
	apiV1.DELETE("/faucet/configs/:id", middleware.JwtTokenAuthMiddleware(s.AdminDB), userMetricServices.DeleteFaucetConfigHandler(s.TrovoWalletDB))

	// KYC Config endpoints
	apiV1.POST("/kyc/configs", middleware.JwtTokenAuthMiddleware(s.AdminDB), accesslog.Audit(s.AdminDB, models.EventKycConfigChange, models.AccessCategoryConfig, accesslog.BodyField("service_provider")), userMetricServices.SaveKycConfigHandler(s.TrovoWalletDB))
	apiV1.GET("/kyc/configs", middleware.JwtTokenAuthMiddleware(s.AdminDB), userMetricServices.GetAllKycConfigsHandler(s.TrovoWalletDB))
	apiV1.GET("/kyc/configs/:id", middleware.JwtTokenAuthMiddleware(s.AdminDB), userMetricServices.GetKycConfigByIDHandler(s.TrovoWalletDB))
	apiV1.DELETE("/kyc/configs/:id", middleware.JwtTokenAuthMiddleware(s.AdminDB), userMetricServices.DeleteKycConfigHandler(s.TrovoWalletDB))

	// Doja Widget endpoints
	apiV1.POST("/doja/widgets", middleware.JwtTokenAuthMiddleware(s.AdminDB), userMetricServices.SaveDojaWidgetHandler(s.TrovoWalletDB))
	apiV1.GET("/doja/widgets", middleware.JwtTokenAuthMiddleware(s.AdminDB), userMetricServices.GetAllDojaWidgetsHandler(s.TrovoWalletDB))
	apiV1.GET("/doja/widgets/:id", middleware.JwtTokenAuthMiddleware(s.AdminDB), userMetricServices.GetDojaWidgetByIDHandler(s.TrovoWalletDB))
	apiV1.DELETE("/doja/widgets/:id", middleware.JwtTokenAuthMiddleware(s.AdminDB), userMetricServices.DeleteDojaWidgetHandler(s.TrovoWalletDB))

	// KYC Level endpoints
	apiV1.POST("/kyc/levels", middleware.JwtTokenAuthMiddleware(s.AdminDB), accesslog.Audit(s.AdminDB, models.EventKycLevelChange, models.AccessCategoryConfig, accesslog.BodyField("user_category")), userMetricServices.SaveKycLevelHandler(s.TrovoWalletDB))
	apiV1.GET("/kyc/levels", middleware.JwtTokenAuthMiddleware(s.AdminDB), userMetricServices.GetAllKycLevelsHandler(s.TrovoWalletDB))
	apiV1.GET("/kyc/levels/:id", middleware.JwtTokenAuthMiddleware(s.AdminDB), userMetricServices.GetKycLevelByIDHandler(s.TrovoWalletDB))
	apiV1.DELETE("/kyc/levels/:id", middleware.JwtTokenAuthMiddleware(s.AdminDB), userMetricServices.DeleteKycLevelHandler(s.TrovoWalletDB))

	// Tokenization Minting endpoints
	apiV1.POST("/tokenization/minting/users", middleware.JwtTokenAuthMiddleware(s.AdminDB), accesslog.Audit(s.AdminDB, models.EventMintingUserGrant, models.AccessCategoryAsset, accesslog.BodyField("username")), userMetricServices.CreateMintingUser(s.TrovoWalletDB))
	apiV1.DELETE("/tokenization/minting/users", middleware.JwtTokenAuthMiddleware(s.AdminDB), accesslog.Audit(s.AdminDB, models.EventMintingUserRevoke, models.AccessCategoryAsset, nil), userMetricServices.DeleteMintingUser(s.TrovoWalletDB))
	apiV1.GET("/tokenization/minting/users", middleware.JwtTokenAuthMiddleware(s.AdminDB), userMetricServices.GetAllMintingUsers(s.TrovoWalletDB))
	apiV1.GET("/tokenization/minting/users/search", middleware.JwtTokenAuthMiddleware(s.AdminDB), userMetricServices.SearchMintingUsers(s.TrovoWalletDB))
}

// func Init(router *gin.Engine, db *gorm.DB, gc *models.GlobalConfig, m *sync.Mutex) {
//	router.GET("/v1/users/metrics", middleware.JwtTokenAuthMiddleware(), GetUserMetrics)
//	router.GET("/v1/wallet/distribution", middleware.JwtTokenAuthMiddleware(), GetWalletDistribution)
//	router.GET("/v1/referrer/count", middleware.JwtTokenAuthMiddleware(), GetReferrerCount)
//	router.GET("/v1/recent/registrations", middleware.JwtTokenAuthMiddleware(), GetRecentRegistrations)
//	router.GET("/v1/country/statistics", middleware.JwtTokenAuthMiddleware(), GetCountryStatistics)
//	router.GET("/v1/p2p/statistics", middleware.JwtTokenAuthMiddleware(), GetP2PStatistics)
//	router.GET("/v1/appeal/list", middleware.JwtTokenAuthMiddleware(), GetAppealList)
//	router.GET("/v1/orders/trade", middleware.JwtTokenAuthMiddleware(), GetTradeList)
//	router.GET("/v1/users/list", middleware.JwtTokenAuthMiddleware(), GetUserList)
//	router.GET("/v1/profile/info", middleware.JwtTokenAuthMiddleware(), GetUserProfile)
//	router.GET("/v1/user/wallets", middleware.JwtTokenAuthMiddleware(), GetUserWallets)
//}

//// @Summary Get user metrics
//// @Description Retrieves metrics related to user activity.
//// @ID GetUserMetrics
//// @Tags Metrics
//// @Security JwtTokenAuth
//// @Produce json
//// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
//// @Success 200 {object} models.WalletUserMetricsResponse
//// @Failure 400,401,500 {object} object
//// @Router /v1/users/metrics [get]
// func GetUserMetrics(c *gin.Context) {
//	walletDb, err := db.TrovoWalletDb()
//	if err != nil {
//		return
//	}
//	ad, _ := middleware.ExtractTokenMetadata(c.Request)
//	userID := ad.UserID
//
//	userInfo, err := userServices.GetUser(userID, walletDb)
//	if err != nil {
//		log.Println("[METRICS] error for user:", userInfo.Username, "error: ", err)
//
//		var ex p2pErrors.GenericError
//		var ok bool
//
//		ex, ok = err.(p2pErrors.GenericError)
//		var statusCode int = 0
//		var response interface{}
//
//		if ok {
//			statusCode = ex.HTTPCode()
//			response = ex.JSONError()
//		} else {
//			statusCode = http.StatusBadRequest
//			response = gin.H{"error": err.Error()}
//		}
//
//		c.JSON(statusCode, response)
//		return
//	}
//
//	userMetrics, err := usermetrics.GetUserMetrics()
//	if err != nil {
//		log.Println("[METRICS] error:", err)
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
//		return
//	}
//
//	c.JSON(http.StatusOK, userMetrics)
//}
//
//// @Summary Get wallet distribution data
//// @Description Retrieves data about the distribution of user wallets.
//// @ID GetWalletDistribution
//// @Tags Wallets
//// @Security JwtTokenAuth
//// @Produce json
//// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
//// @Success 200 {object} models.WalletVersionDistribution
//// @Failure 400,401,500 {object} object
//// @Router /v1/wallet/distribution [get]
//func GetWalletDistribution(c *gin.Context) {
//	walletDb, err := db.TrovoWalletDb()
//	if err != nil {
//		return
//	}
//	ad, _ := middleware.ExtractTokenMetadata(c.Request)
//	userID := ad.UserID
//
//	userInfo, err := userServices.GetUser(userID, walletDb)
//	if err != nil {
//		log.Println("[METRICS] error for user:", userInfo.Username, "error: ", err)
//
//		var ex p2pErrors.GenericError
//		var ok bool
//
//		ex, ok = err.(p2pErrors.GenericError)
//		var statusCode int = 0
//		var response interface{}
//
//		if ok {
//			statusCode = ex.HTTPCode()
//			response = ex.JSONError()
//		} else {
//			statusCode = http.StatusBadRequest
//			response = gin.H{"error": err.Error()}
//		}
//
//		c.JSON(statusCode, response)
//		return
//	}
//
//	walletDistribution, err := usermetrics.GetWalletVersionDistributionData()
//	if err != nil {
//		log.Println("[METRICS] error:", err)
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
//		return
//	}
//	c.JSON(http.StatusOK, walletDistribution)
//}
//
//// @Summary Get top referrer count
//// @Description Retrieves the count of top referrers.
//// @ID GetReferrerCount
//// @Tags Referrers
//// @Security JwtTokenAuth
//// @Produce json
//// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
//// @Success 200 {object} models.ReferrerCount
//// @Failure 400,401,500 {object} object
//// @Router /v1/referrer/count [get]
//func GetReferrerCount(c *gin.Context) {
//	walletDb, err := db.TrovoWalletDb()
//	if err != nil {
//		return
//	}
//	ad, _ := middleware.ExtractTokenMetadata(c.Request)
//	userID := ad.UserID
//
//	userInfo, err := userServices.GetUser(userID, walletDb)
//	if err != nil {
//		log.Println("[METRICS] error for user:", userInfo.Username, "error: ", err)
//
//		var ex p2pErrors.GenericError
//		var ok bool
//
//		ex, ok = err.(p2pErrors.GenericError)
//		var statusCode int = 0
//		var response interface{}
//
//		if ok {
//			statusCode = ex.HTTPCode()
//			response = ex.JSONError()
//		} else {
//			statusCode = http.StatusBadRequest
//			response = gin.H{"error": err.Error()}
//		}
//
//		c.JSON(statusCode, response)
//		return
//	}
//
//	topReferrer, err := usermetrics.GetTopReferrers()
//	if err != nil {
//		log.Println("[METRICS] error:", err)
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
//		return
//	}
//	c.JSON(http.StatusOK, topReferrer)
//}
//
//// @Summary Get recent registrations
//// @Description Retrieves information about recent user registrations.
//// @ID GetRecentRegistrations
//// @Tags Registrations
//// @Security JwtTokenAuth
//// @Produce json
//// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
//// @Success 200 {object} models.UserDTO
//// @Failure 400,401,500 {object} object
//// @Router /v1/recent/registrations [get]
//func GetRecentRegistrations(c *gin.Context) {
//	walletDb, err := db.TrovoWalletDb()
//	if err != nil {
//		return
//	}
//	ad, _ := middleware.ExtractTokenMetadata(c.Request)
//	userID := ad.UserID
//
//	userInfo, err := userServices.GetUser(userID, walletDb)
//	if err != nil {
//		log.Println("[METRICS] error for user:", userInfo.Username, "error: ", err)
//
//		var ex p2pErrors.GenericError
//		var ok bool
//
//		ex, ok = err.(p2pErrors.GenericError)
//		var statusCode int = 0
//		var response interface{}
//
//		if ok {
//			statusCode = ex.HTTPCode()
//			response = ex.JSONError()
//		} else {
//			statusCode = http.StatusBadRequest
//			response = gin.H{"error": err.Error()}
//		}
//
//		c.JSON(statusCode, response)
//		return
//	}
//
//	recentRegistrations, err := usermetrics.GetRecentRegistrations(10)
//	if err != nil {
//		log.Println("[METRICS] error:", err)
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
//		return
//	}
//	// Send the response
//	c.JSON(http.StatusOK, recentRegistrations)
//}
//
//// @Summary Get country statistics
//// @Description Retrieves statistics related to user distribution by country.
//// @ID GetCountryStatistics
//// @Tags Statistics
//// @Security JwtTokenAuth
//// @Produce json
//// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
//// @Success 200 {object} models.CountryUserCount
//// @Failure 400,401,500 {object} object
//// @Router /v1/country/statistics [get]
//func GetCountryStatistics(c *gin.Context) {
//	walletDb, err := db.TrovoWalletDb()
//	if err != nil {
//		return
//	}
//	ad, _ := middleware.ExtractTokenMetadata(c.Request)
//	userID := ad.UserID
//
//	userInfo, err := userServices.GetUser(userID, walletDb)
//	if err != nil {
//		log.Println("[METRICS] error for user:", userInfo.Username, "error: ", err)
//
//		var ex p2pErrors.GenericError
//		var ok bool
//
//		ex, ok = err.(p2pErrors.GenericError)
//		var statusCode int = 0
//		var response interface{}
//
//		if ok {
//			statusCode = ex.HTTPCode()
//			response = ex.JSONError()
//		} else {
//			statusCode = http.StatusBadRequest
//			response = gin.H{"error": err.Error()}
//		}
//
//		c.JSON(statusCode, response)
//		return
//	}
//
//	countryStatistics, err := usermetrics.CountUsersByCountry()
//	if err != nil {
//		log.Println("[METRICS] error:", err)
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
//		return
//	}
//	// Send the response
//	c.JSON(http.StatusOK, countryStatistics)
//}
//
//// @Summary Get P2P statistics
//// @Description Retrieves statistics related to P2P metrics.
//// @ID GetP2PStatistics
//// @Tags P2P
//// @Security JwtTokenAuth
//// @Produce json
//// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
//// @Success 200 {object} models.P2PUserMetricsResponse
//// @Failure 400,401,500 {object} object
//// @Router /v1/p2p/statistics [get]
//func GetP2PStatistics(c *gin.Context) {
//	walletDb, err := db.TrovoWalletDb()
//	if err != nil {
//		return
//	}
//	ad, _ := middleware.ExtractTokenMetadata(c.Request)
//	userID := ad.UserID
//
//	userInfo, err := userServices.GetUser(userID, walletDb)
//	if err != nil {
//		log.Println("[METRICS] error for user:", userInfo.Username, "error: ", err)
//
//		var ex p2pErrors.GenericError
//		var ok bool
//
//		ex, ok = err.(p2pErrors.GenericError)
//		var statusCode int = 0
//		var response interface{}
//
//		if ok {
//			statusCode = ex.HTTPCode()
//			response = ex.JSONError()
//		} else {
//			statusCode = http.StatusBadRequest
//			response = gin.H{"error": err.Error()}
//		}
//
//		c.JSON(statusCode, response)
//		return
//	}
//
//	p2pMetrics, err := usermetrics.GetP2PMetrics()
//	if err != nil {
//		log.Println("[METRICS] error:", err)
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
//		return
//	}
//
//	c.JSON(http.StatusOK, p2pMetrics)
//}
//
//// @Summary Get list of trades on appeal
//// @Description Retrieves the list of trades that are currently on appeal.
//// @ID GetAppealList
//// @Tags Appeals
//// @Security JwtTokenAuth
//// @Produce json
//// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
//// @Success 200 {object} models.TradesOnAppeal
//// @Failure 400,401,500 {object} object
//// @Router /v1/appeal/list [get]
//func GetAppealList(c *gin.Context) {
//	walletDb, err := db.TrovoWalletDb()
//	if err != nil {
//		return
//	}
//	ad, _ := middleware.ExtractTokenMetadata(c.Request)
//	userID := ad.UserID
//
//	userInfo, err := userServices.GetUser(userID, walletDb)
//	if err != nil {
//		log.Println("[METRICS] error for user:", userInfo.Username, "error: ", err)
//
//		var ex p2pErrors.GenericError
//		var ok bool
//
//		ex, ok = err.(p2pErrors.GenericError)
//		var statusCode int = 0
//		var response interface{}
//
//		if ok {
//			statusCode = ex.HTTPCode()
//			response = ex.JSONError()
//		} else {
//			statusCode = http.StatusBadRequest
//			response = gin.H{"error": err.Error()}
//		}
//
//		c.JSON(statusCode, response)
//		return
//	}
//
//	appealList, err := usermetrics.GetTradesOnAppealList()
//	if err != nil {
//		log.Println("[METRICS] error:", err)
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
//		return
//	}
//
//	// Send the response
//	c.JSON(http.StatusOK, appealList)
//}
//
//// @Summary Get list of trades
//// @Description Retrieves the list of trades for the authenticated user.
//// @ID GetTradeList
//// @Tags Trades
//// @Security JwtTokenAuth
//// @Produce json
//// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
//// @Success 200 {object} models.P2POrder
//// @Failure 400,401,500 {object} object
//// @Router /v1/orders/trade [get]
//func GetTradeList(c *gin.Context) {
//	walletDb, err := db.TrovoWalletDb()
//	if err != nil {
//		return
//	}
//	ad, _ := middleware.ExtractTokenMetadata(c.Request)
//	userID := ad.UserID
//
//	userInfo, err := userServices.GetUser(userID, walletDb)
//	if err != nil {
//		log.Println("[METRICS] error for user:", userInfo.Username, "error: ", err)
//
//		var ex p2pErrors.GenericError
//		var ok bool
//
//		ex, ok = err.(p2pErrors.GenericError)
//		var statusCode int = 0
//		var response interface{}
//
//		if ok {
//			statusCode = ex.HTTPCode()
//			response = ex.JSONError()
//		} else {
//			statusCode = http.StatusBadRequest
//			response = gin.H{"error": err.Error()}
//		}
//
//		c.JSON(statusCode, response)
//		return
//	}
//
//	tradeList, err := usermetrics.GetAllTradesList()
//	if err != nil {
//		log.Println("[METRICS] error:", err)
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
//		return
//	}
//
//	// Send the response
//	c.JSON(http.StatusOK, tradeList)
//}
//
//// @Summary Get list of users
//// @Description Retrieves the list of users for the authenticated user.
//// @ID GetUserList
//// @Tags Users
//// @Security JwtTokenAuth
//// @Produce json
//// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
//// @Success 200 {object} models.UserDto
//// @Failure 400,401,500 {object} object
//// @Router /v1/users/list [get]
//func GetUserList(c *gin.Context) {
//	walletDb, err := db.TrovoWalletDb()
//	if err != nil {
//		return
//	}
//	ad, _ := middleware.ExtractTokenMetadata(c.Request)
//	userID := ad.UserID
//
//	userInfo, err := userServices.GetUser(userID, walletDb)
//	if err != nil {
//		log.Println("[METRICS] error for user:", userInfo.Username, "error: ", err)
//
//		var ex p2pErrors.GenericError
//		var ok bool
//
//		ex, ok = err.(p2pErrors.GenericError)
//		var statusCode int = 0
//		var response interface{}
//
//		if ok {
//			statusCode = ex.HTTPCode()
//			response = ex.JSONError()
//		} else {
//			statusCode = http.StatusBadRequest
//			response = gin.H{"error": err.Error()}
//		}
//
//		c.JSON(statusCode, response)
//		return
//	}
//
//	tradeList, err := usermetrics.GetUserList()
//	if err != nil {
//		log.Println("[METRICS] error:", err)
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
//		return
//	}
//
//	// Send the response
//	c.JSON(http.StatusOK, tradeList)
//}
//
//// @Summary Get user profile information
//// @Description Retrieves the profile information of the authenticated user.
//// @ID GetUserProfileInfo
//// @Tags Users
//// @Security JwtTokenAuth
//// @Produce json
//// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
//// @Success 200 {object} models.UserDto
//// @Failure 400,401,500 {object} object
//// @Router /v1/profile/info [get]
//func GetUserProfileInfo(c *gin.Context) {
//	db, err := db.TrovoWalletDb()
//	if err != nil {
//		return
//	}
//	ad, _ := middleware.ExtractTokenMetadata(c.Request)
//	userID := ad.UserID
//
//	userInfo, err := userServices.GetUser(userID, db)
//	if err != nil {
//		log.Println("[METRICS] error for user:", userInfo.Username, "error: ", err)
//
//		var ex p2pErrors.GenericError
//		var ok bool
//
//		ex, ok = err.(p2pErrors.GenericError)
//		var statusCode int = 0
//		var response interface{}
//
//		if ok {
//			statusCode = ex.HTTPCode()
//			response = ex.JSONError()
//		} else {
//			statusCode = http.StatusBadRequest
//			response = gin.H{"error": err.Error()}
//		}
//
//		c.JSON(statusCode, response)
//		return
//	}
//
//	userInf, err := usermetrics.GetUserList()
//	if err != nil {
//		log.Println("[METRICS] error:", err)
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
//		return
//	}
//
//	// Send the response
//	c.JSON(http.StatusOK, userInf)
//}
//
//// @Summary Get user wallets
//// @Description Retrieves the wallets information of the authenticated user.
//// @ID GetUserWallets
//// @Tags Users
//// @Security JwtTokenAuth
//// @Produce json
//// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
//// @Success 200 {object} models.UserDto
//// @Failure 400,401,500 {object} object
//// @Router /v1/user/wallets [get]
//func GetUserWallets(c *gin.Context) {
//	walletDb, err := db.TrovoWalletDb()
//	if err != nil {
//		return
//	}
//	ad, _ := middleware.ExtractTokenMetadata(c.Request)
//	userID := ad.UserID
//
//	userInfo, err := userServices.GetUser(userID, walletDb)
//	if err != nil {
//		log.Println("[METRICS] error for user:", userInfo.Username, "error: ", err)
//
//		var ex p2pErrors.GenericError
//		var ok bool
//
//		ex, ok = err.(p2pErrors.GenericError)
//		var statusCode int = 0
//		var response interface{}
//
//		if ok {
//			statusCode = ex.HTTPCode()
//			response = ex.JSONError()
//		} else {
//			statusCode = http.StatusBadRequest
//			response = gin.H{"error": err.Error()}
//		}
//
//		c.JSON(statusCode, response)
//		return
//	}
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
//
//// @Summary Get user profile information
//// @Description Retrieves the profile information of the authenticated user.
//// @ID GetUserProfile
//// @Tags Users
//// @Security JwtTokenAuth
//// @Produce json
//// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
//// @Success 200 {object} models.UserDto
//// @Failure 400,401,500 {object} object
//// @Router /v1/profile/info [get]
//func GetUserProfile(c *gin.Context) {
//	walletDb, err := db.TrovoWalletDb()
//	if err != nil {
//		return
//	}
//	ad, _ := middleware.ExtractTokenMetadata(c.Request)
//	userID := ad.UserID
//
//	userInfo, err := userServices.GetUser(userID, walletDb)
//	if err != nil {
//		log.Println("[METRICS] error for user:", userInfo.Username, "error: ", err)
//
//		var ex p2pErrors.GenericError
//		var ok bool
//
//		ex, ok = err.(p2pErrors.GenericError)
//		var statusCode int = 0
//		var response interface{}
//
//		if ok {
//			statusCode = ex.HTTPCode()
//			response = ex.JSONError()
//		} else {
//			statusCode = http.StatusBadRequest
//			response = gin.H{"error": err.Error()}
//		}
//
//		c.JSON(statusCode, response)
//		return
//	}
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
