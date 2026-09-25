package usermetrics

import (
	usermetricsDB "admin-panel-dashboard/internal/components/usermetrics/db"
	"admin-panel-dashboard/internal/server/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// @Summary Get all service links
// @Description Retrieves a list of all service links.
// @ID GetAllServiceLinks
// @Tags ServiceLinks
// @Security JwtTokenAuth
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Success 200 {object} response.Data{data=[]usermetricsDB.ServiceLink}
// @Failure 500 {object} map[string]string
// @Router /service-links [get]
func GetServiceLinksHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		links, err := usermetricsDB.GetAllServiceLinks(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		response.JSON(c, http.StatusOK, "Service links fetched successfully", links, nil)
	}
}
