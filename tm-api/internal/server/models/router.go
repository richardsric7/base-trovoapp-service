package server

import (
	"github.com/gin-gonic/gin"
)

// func (s *Server) defineRoutes(router *gin.Engine, db *gorm.DB, gc *models.GlobalConfig, m *sync.Mutex) {
//	apiV1 := router.Group("/api/v1")
//	//apiV1.Use(middleware.Authenticate())
//	apiV1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
//
//	apiV1.GET("/users/metrics", middleware.JwtTokenAuthMiddleware(), GetUserMetrics)
//	apiV1.GET("/wallet/distribution", middleware.JwtTokenAuthMiddleware(), GetWalletDistribution)
//	apiV1.GET("/referrer/count", middleware.JwtTokenAuthMiddleware(), GetReferrerCount)
//	apiV1.GET("/recent/registrations", middleware.JwtTokenAuthMiddleware(), GetRecentRegistrations)
//	apiV1.GET("/country/statistics", middleware.JwtTokenAuthMiddleware(), GetCountryStatistics)
//	apiV1.GET("/p2p/statistics", middleware.JwtTokenAuthMiddleware(), GetP2PStatistics)
//	apiV1.GET("/appeal/list", middleware.JwtTokenAuthMiddleware(), GetAppealList)
//	apiV1.GET("/orders/trade", middleware.JwtTokenAuthMiddleware(), GetTradeList)
//	apiV1.GET("/users/list", middleware.JwtTokenAuthMiddleware(), GetUserList)
//	apiV1.GET("/profile/info", middleware.JwtTokenAuthMiddleware(), GetProfileInfo)
//	apiV1.GET("/user/wallets", middleware.JwtTokenAuthMiddleware(), GetUserWallets)
//}

func (s *Server) SetupRouter(r *gin.Engine) *gin.Engine {

	// apiV1.GET("/users/metrics", GetUserMetrics)
	// apiV1.GET("/wallet/distribution", GetWalletDistribution)
	// apiV1.GET("/referrer/count", GetReferrerCount)
	// apiV1.GET("/recent/registrations", GetRecentRegistrations)
	// apiV1.GET("/country/statistics", GetCountryStatistics)
	// apiV1.GET("/p2p/statistics", middleware.JwtTokenAuthMiddleware(), GetP2PStatistics)
	// apiV1.GET("/appeal/list", GetAppealList)
	// apiV1.GET("/orders/trade", GetTradeList)
	// apiV1.GET("/users/list", GetUserList)
	// apiV1.GET("/profile/info", GetProfileInfo)
	// apiV1.GET("/user/wallets", GetUserWallets)

	// Admin invite/permission routes

	return r
}

func (s *Server) SetupRouterParams(r *gin.Engine) *gin.Engine {

	// Request logging and panic recovery are installed in main.go by the
	// observe package (structured JSON logs with a request ID, and recovery
	// that reports crashes to GlitchTip). The gin.LoggerWithFormatter and
	// gin.Recovery that used to live here were removed rather than kept
	// alongside them: two loggers meant every request was logged twice, and
	// gin.Recovery would have swallowed panics before observe.Recovery could
	// report them.

	// CORS is handled by middleware.CORSMiddleware (main.go), which echoes the
	// allowed origin. The stock cors layer that used to sit here answered with
	// Access-Control-Allow-Origin: * (plus Allow-Credentials), which browsers
	// reject - so unmatched routes (404s) surfaced as misleading CORS errors.

	return r
}
