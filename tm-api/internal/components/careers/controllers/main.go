package careers

import (
	careerServices "admin-panel-dashboard/internal/components/careers/services"
	"admin-panel-dashboard/internal/middleware"
	serverModels "admin-panel-dashboard/internal/server/models"

	"github.com/gin-gonic/gin"
)

// Init initializes the careers module routes
func Init(router *gin.Engine, s *serverModels.Server) {
	apiV1 := router.Group("/api/v1")

	// Public endpoints (no authentication required - only published roles)
	apiV1.GET("/careers/roles", careerServices.GetPublicCareerRolesHandler(s.AdminDB))
	apiV1.GET("/careers/roles/:id", careerServices.GetPublicCareerRoleByIDHandler(s.AdminDB))

	// Private/Admin endpoints (authentication required)
	adminRoutes := apiV1.Group("/admin/careers")
	adminRoutes.Use(middleware.JwtTokenAuthMiddleware(s.AdminDB))
	{
		adminRoutes.GET("/roles", careerServices.GetAdminCareerRolesHandler(s.AdminDB))
		adminRoutes.GET("/roles/:id", careerServices.GetAdminCareerRoleByIDHandler(s.AdminDB))
		adminRoutes.POST("/roles", careerServices.CreateCareerRoleHandler(s.AdminDB, s.TrovoWalletDB))
		adminRoutes.PUT("/roles/:id", careerServices.UpdateCareerRoleHandler(s.AdminDB, s.TrovoWalletDB))
		adminRoutes.DELETE("/roles/:id", careerServices.DeleteCareerRoleHandler(s.AdminDB, s.TrovoWalletDB))
		adminRoutes.PUT("/roles/:id/publish", careerServices.PublishCareerRoleHandler(s.AdminDB, s.TrovoWalletDB))
		adminRoutes.PUT("/roles/:id/status", careerServices.UpdateCareerRoleStatusHandler(s.AdminDB, s.TrovoWalletDB))
	}
}
