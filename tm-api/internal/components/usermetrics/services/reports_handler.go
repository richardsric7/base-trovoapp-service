package usermetrics

import (
	usermetricsDB "admin-panel-dashboard/internal/components/usermetrics/db"
	"admin-panel-dashboard/internal/middleware"
	"admin-panel-dashboard/internal/models"
	serverResponse "admin-panel-dashboard/internal/server/response"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// reportsPermission is the seeded AdminPermission ("ACCESS_REPORTS", see
// internal/db/main.go's seedPermissions) every report handler below is
// gated on - granted to SUPER_ADMIN and EDIT_LEVEL_ADMIN, not
// VIEW_ONLY_ADMIN. This is real RBAC (models.HasPermission resolves the
// caller's role -> RolePermission -> AdminPermission), not just "any
// logged-in admin" the way the older P2P statistics endpoints check.
const reportsPermission = "ACCESS_REPORTS"

// checkReportsAccess resolves the calling admin from their JWT and
// verifies they hold reportsPermission, writing the appropriate error
// response itself on failure (401 if the admin can't be resolved at all,
// 403 if they're a real admin without report access).
func checkReportsAccess(c *gin.Context, adminDB *gorm.DB) bool {
	ad, err := middleware.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return false
	}
	var admin models.AdminUser
	if err := adminDB.First(&admin, "wallet_user_id = ?", ad.UserID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return false
	}
	if !models.HasPermission(adminDB, &admin, reportsPermission) {
		c.JSON(http.StatusForbidden, gin.H{"error": "you do not have permission to access reports"})
		return false
	}
	return true
}

// @Summary Get P2P trading volume report
// @Description Daily trading volume/order-count trend over a date range.
// @Tags P2P Reports
// @Security JwtTokenAuth
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Param range query string false "7d, 30d, 90d, or 1y" default(30d)
// @Success 200 {object} response.Data
// @Failure 401,403,500 {object} object
// @Router /p2p/reports/volume [get]
func GetP2PVolumeReportHandler(adminDB, p2pDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !checkReportsAccess(c, adminDB) {
			return
		}
		from, to := usermetricsDB.ResolveReportRange(c.Query("range"))
		report, err := usermetricsDB.GetP2PVolumeReport(p2pDB, from, to)
		if err != nil {
			log.Println("[REPORTS] volume error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
		serverResponse.JSON(c, http.StatusOK, "Volume report fetched successfully", report, nil)
	}
}

// @Summary Get P2P distribution report
// @Description Trade breakdown by asset, currency, country, and order status.
// @Tags P2P Reports
// @Security JwtTokenAuth
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Param range query string false "7d, 30d, 90d, or 1y" default(30d)
// @Success 200 {object} response.Data
// @Failure 401,403,500 {object} object
// @Router /p2p/reports/distribution [get]
func GetP2PDistributionReportHandler(adminDB, p2pDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !checkReportsAccess(c, adminDB) {
			return
		}
		from, to := usermetricsDB.ResolveReportRange(c.Query("range"))
		report, err := usermetricsDB.GetP2PDistributionReport(p2pDB, from, to)
		if err != nil {
			log.Println("[REPORTS] distribution error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
		serverResponse.JSON(c, http.StatusOK, "Distribution report fetched successfully", report, nil)
	}
}

// @Summary Get P2P dispute report
// @Description Dispute open/resolve trend, resolution rate, and breakdowns by subject/resolution.
// @Tags P2P Reports
// @Security JwtTokenAuth
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Param range query string false "7d, 30d, 90d, or 1y" default(30d)
// @Success 200 {object} response.Data
// @Failure 401,403,500 {object} object
// @Router /p2p/reports/disputes [get]
func GetP2PDisputeReportHandler(adminDB, p2pDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !checkReportsAccess(c, adminDB) {
			return
		}
		from, to := usermetricsDB.ResolveReportRange(c.Query("range"))
		report, err := usermetricsDB.GetP2PDisputeReport(p2pDB, from, to)
		if err != nil {
			log.Println("[REPORTS] disputes error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
		serverResponse.JSON(c, http.StatusOK, "Dispute report fetched successfully", report, nil)
	}
}

// @Summary Get P2P revenue report
// @Description Platform/regulatory fee and VAT revenue collected on completed orders, trended over time.
// @Tags P2P Reports
// @Security JwtTokenAuth
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Param range query string false "7d, 30d, 90d, or 1y" default(30d)
// @Success 200 {object} response.Data
// @Failure 401,403,500 {object} object
// @Router /p2p/reports/revenue [get]
func GetP2PRevenueReportHandler(adminDB, p2pDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !checkReportsAccess(c, adminDB) {
			return
		}
		from, to := usermetricsDB.ResolveReportRange(c.Query("range"))
		report, err := usermetricsDB.GetP2PRevenueReport(p2pDB, from, to)
		if err != nil {
			log.Println("[REPORTS] revenue error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
		serverResponse.JSON(c, http.StatusOK, "Revenue report fetched successfully", report, nil)
	}
}

// @Summary Get P2P growth report
// @Description New offers and first-time merchants trended over a date range.
// @Tags P2P Reports
// @Security JwtTokenAuth
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Param range query string false "7d, 30d, 90d, or 1y" default(30d)
// @Success 200 {object} response.Data
// @Failure 401,403,500 {object} object
// @Router /p2p/reports/growth [get]
func GetP2PGrowthReportHandler(adminDB, p2pDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !checkReportsAccess(c, adminDB) {
			return
		}
		from, to := usermetricsDB.ResolveReportRange(c.Query("range"))
		report, err := usermetricsDB.GetP2PGrowthReport(p2pDB, from, to)
		if err != nil {
			log.Println("[REPORTS] growth error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
		serverResponse.JSON(c, http.StatusOK, "Growth report fetched successfully", report, nil)
	}
}

// @Summary Get P2P hourly activity report
// @Description Order count by hour-of-day (UTC), summed across the date range - for staffing/support-hours planning.
// @Tags P2P Reports
// @Security JwtTokenAuth
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Param range query string false "7d, 30d, 90d, or 1y" default(30d)
// @Success 200 {object} response.Data
// @Failure 401,403,500 {object} object
// @Router /p2p/reports/hourly-activity [get]
func GetP2PHourlyActivityReportHandler(adminDB, p2pDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !checkReportsAccess(c, adminDB) {
			return
		}
		from, to := usermetricsDB.ResolveReportRange(c.Query("range"))
		report, err := usermetricsDB.GetP2PHourlyActivityReport(p2pDB, from, to)
		if err != nil {
			log.Println("[REPORTS] hourly activity error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
		serverResponse.JSON(c, http.StatusOK, "Hourly activity report fetched successfully", report, nil)
	}
}

// @Summary Get P2P merchant leaderboard report
// @Description Paginated merchant performance leaderboard, sortable by trades/volume/completion rate/disputes.
// @Tags P2P Reports
// @Security JwtTokenAuth
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Items per page" default(10)
// @Param sortBy query string false "completedTrades, completedVolume, completionRate, or disputesOpened" default(completedTrades)
// @Success 200 {object} response.Data
// @Failure 401,403,500 {object} object
// @Router /p2p/reports/merchants [get]
func GetP2PMerchantLeaderboardHandler(adminDB, p2pDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !checkReportsAccess(c, adminDB) {
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
		leaderboard, total, err := usermetricsDB.GetP2PMerchantLeaderboard(p2pDB, usermetricsDB.MerchantLeaderboardRequest{
			Page:     page,
			PageSize: pageSize,
			SortBy:   c.Query("sortBy"),
		})
		if err != nil {
			log.Println("[REPORTS] merchant leaderboard error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
		serverResponse.JSON(c, http.StatusOK, "Merchant leaderboard fetched successfully", gin.H{
			"data":     leaderboard,
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
		}, nil)
	}
}
