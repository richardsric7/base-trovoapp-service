package root

import (
	root "admin-panel-dashboard/internal/components/root/services"

	"github.com/gin-gonic/gin"
)

func Init(router *gin.Engine) {

	// Returns state of the api service
	router.GET("/", func(c *gin.Context) {

		rootInfo := root.GetRootInfo()
		rootInfo.Service = "TROVOTECH ADMIN DASHBOARD API Platform"
		rootInfo.Status = "ACTIVE"
		c.JSON(200, rootInfo)
	})
}
