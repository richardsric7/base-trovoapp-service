package swagger

import (
	"admin-panel-dashboard/internal/observe"
	serverModels "admin-panel-dashboard/internal/server/models"
	"net/http"
	"os"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/ecnepsnai/discord"
	"github.com/gin-gonic/gin"
)

// Init initializes
func Init(router *gin.Engine, s *serverModels.Server) {

	// Serve Swagger UI files with cache-busting headers
	router.GET("/swagger/*any", func(c *gin.Context) {
		c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")
		ginSwagger.WrapHandler(swaggerFiles.Handler)(c)
	})

	// Your existing routes
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	apiV1 := router.Group("/api/v1")

	apiV1.GET("/swagger/*any", func(c *gin.Context) {
		c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")
		ginSwagger.WrapHandler(swaggerFiles.Handler)(c)
	})
}

func LogDiscordError(msg string) {
	// Also count it, so the failure shows on the error-rate dashboard
	// instead of only in Discord. Additive: Discord is unchanged.
	observe.RecordHandledFailure(msg)

	discord.WebhookURL = "https://discord.com/api/webhooks/865931042795290636/jObHzZWdnbhX1jomOSZQX8Ip5AXLArh87PI4-ZQ8u6ssnRbZuVdY_iPxz5qoWkHUlZwS"
	if len(os.Getenv("500_ERROR_WEBHOOK")) > 50 {
		discord.WebhookURL = os.Getenv("500_ERROR_WEBHOOK")
	}
	err := discord.Say(msg)
	if err != nil {
		return
	}
}
