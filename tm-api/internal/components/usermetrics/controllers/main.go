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
	apiV1.GET("/p2p/statistics", middleware.JwtTokenAuthMiddleware(s.AdminDB), userMetricServices.GetP2PStatistics(s.TrovoWalletDB, s.P2P))
	apiV1.GET("/trades/statistics", middleware.JwtTokenAuthMiddleware(s.AdminDB), userMetricServices.GetTradeStatistics(s.TrovoWalletDB, s.P2P))
	apiV1.GET("/orders/trade", middleware.JwtTokenAuthMiddleware(s.AdminDB), userMetricServices.GetTradeList(s.TrovoWalletDB, s.P2P))

	// P2P Market Reports - gated on the ACCESS_REPORTS permission (see
	// reportsPermission in usermetrics/services/reports_handler.go), not
	// just "any logged-in admin" like the statistics endpoints above.
	apiV1.GET("/p2p/reports/volume", middleware.JwtTokenAuthMiddleware(s.AdminDB), userMetricServices.GetP2PVolumeReportHandler(s.AdminDB, s.P2P))
	apiV1.GET("/p2p/reports/distribution", middleware.JwtTokenAuthMiddleware(s.AdminDB), userMetricServices.GetP2PDistributionReportHandler(s.AdminDB, s.P2P))
	apiV1.GET("/p2p/reports/disputes", middleware.JwtTokenAuthMiddleware(s.AdminDB), userMetricServices.GetP2PDisputeReportHandler(s.AdminDB, s.P2P))
	apiV1.GET("/p2p/reports/revenue", middleware.JwtTokenAuthMiddleware(s.AdminDB), userMetricServices.GetP2PRevenueReportHandler(s.AdminDB, s.P2P))
	apiV1.GET("/p2p/reports/growth", middleware.JwtTokenAuthMiddleware(s.AdminDB), userMetricServices.GetP2PGrowthReportHandler(s.AdminDB, s.P2P))
	apiV1.GET("/p2p/reports/hourly-activity", middleware.JwtTokenAuthMiddleware(s.AdminDB), userMetricServices.GetP2PHourlyActivityReportHandler(s.AdminDB, s.P2P))

	// Curated assets management - the platform's asset catalog, including
	// the P2PEnabled flag that gates whether an asset can be used to
	// create a P2P offer or appears in marketplace search (enforced by
	// app-backend's P2P module against the same p2p_enabled column).
	apiV1.GET("/assets/curated", middleware.JwtTokenAuthMiddleware(s.AdminDB), userMetricServices.GetCuratedAssetsHandler(s.TrovoWalletDB))
	apiV1.GET("/assets/curated/:id", middleware.JwtTokenAuthMiddleware(s.AdminDB), userMetricServices.GetCuratedAssetByIDHandler(s.TrovoWalletDB))
	apiV1.POST("/assets/curated", middleware.JwtTokenAuthMiddleware(s.AdminDB), accesslog.Audit(s.AdminDB, models.EventAssetCurationChange, models.AccessCategoryAsset, accesslog.BodyField("assetCode")), userMetricServices.SaveCuratedAssetHandler(s.TrovoWalletDB))
	apiV1.PUT("/assets/curated/:id/p2p-enabled", middleware.JwtTokenAuthMiddleware(s.AdminDB), accesslog.Audit(s.AdminDB, models.EventAssetCurationChange, models.AccessCategoryAsset, accesslog.Param("id")), userMetricServices.SetCuratedAssetP2PEnabledHandler(s.TrovoWalletDB))
	// Deactivate/reactivate - the supported way to retire a curated asset;
	// there is intentionally no delete endpoint, since existing wallets,
	// offers and orders can still reference the asset by code.
	apiV1.PUT("/assets/curated/:id/inactive", middleware.JwtTokenAuthMiddleware(s.AdminDB), accesslog.Audit(s.AdminDB, models.EventAssetCurationChange, models.AccessCategoryAsset, accesslog.Param("id")), userMetricServices.SetCuratedAssetInactiveHandler(s.TrovoWalletDB))
	apiV1.GET("/asset-classes", middleware.JwtTokenAuthMiddleware(s.AdminDB), userMetricServices.GetAssetClassesHandler(s.TrovoWalletDB))
	apiV1.GET("/p2p/reports/merchants", middleware.JwtTokenAuthMiddleware(s.AdminDB), userMetricServices.GetP2PMerchantLeaderboardHandler(s.AdminDB, s.P2P))
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
	apiV1.GET("/p2p/users", middleware.JwtTokenAuthMiddleware(s.AdminDB), userMetricServices.GetP2PUsers(s.TrovoWalletDB, s.P2P))

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

