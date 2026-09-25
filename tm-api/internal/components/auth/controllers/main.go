package auth

import (
	authServices "admin-panel-dashboard/internal/components/auth/services"
	"admin-panel-dashboard/internal/observe"
	serverModels "admin-panel-dashboard/internal/server/models"
	"log"
	"os"

	"github.com/ecnepsnai/discord"
	"github.com/gin-gonic/gin"
)

// Init initializes
func Init(router *gin.Engine, s *serverModels.Server) {

	// login callback url

	apiV1 := router.Group("/api/v1")

	apiV1.POST("/callbacks/login/:serviceName", func(c *gin.Context) {
		authServices.LoginCallback(c, s)
	})
	apiV1.POST("/callbacks/auth/:serviceName", func(c *gin.Context) {
		authServices.AuthCallback(c, s)
	})

	apiV1.POST("/login", func(c *gin.Context) {
		authServices.Login(c, s)
	})
	apiV1.GET("/users/verify/:targetUser/:loginID", func(c *gin.Context) {
		authServices.VerifyLoginID(c, s)
	})
	// Subscribe side of the login push-notification flow (Section 3, login-tab-throttling-fix-plan.md)
	// — BroadcastToLoginID already fires from LoginCallback on every real approval; this is what
	// lets a browser actually receive it, instead of only ever finding out via polling.
	apiV1.GET("/login/stream/:loginID", func(c *gin.Context) {
		authServices.LoginNotificationStream(c, s)
	})
	apiV1.POST("/login/token/refresh", func(c *gin.Context) {
		authServices.RefreshToken(c, s)
	})
	apiV1.POST("/logout", func(c *gin.Context) {
		authServices.Logout(c, s)
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
		log.Println("Failed to send discord error message", err)
		return
	}
}
