package users

import (
	"admin-panel-dashboard/internal/components/accesslog"
	users "admin-panel-dashboard/internal/components/users/services"
	conDB "admin-panel-dashboard/internal/db"
	p2pErrors "admin-panel-dashboard/internal/errors"
	"admin-panel-dashboard/internal/models"
	"admin-panel-dashboard/internal/observe"
	serverModels "admin-panel-dashboard/internal/server/models"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"admin-panel-dashboard/internal/middleware"
	"log"
	"net/http"
	"strings"

	"github.com/ecnepsnai/discord"
	"github.com/gin-gonic/gin"
)

// Init initializes /v2/users endpoint
func Init(router *gin.Engine, s *serverModels.Server) {
	// get user detail
	router.GET("/v1/users/detail/:targetUser", middleware.JwtTokenAuthMiddleware(s.AdminDB), func(c *gin.Context) {
		var err error

		userIdentifier := strings.ReplaceAll(strings.ToLower(c.Param("targetUser")), " ", "")

		cacheKey := fmt.Sprintf("[GET] /v1/detail/users/%v", userIdentifier)
		{
			// search cache

			// cacheKeyParameters := fmt.Sprintf("limit=%v&order=%v&cursor=%v&forTransactionHash=%v&includeHash=%v&temp=%v", limit, orderStr, cursor, forTransactionHash, includeHash, temp)

			ok, status, response := s.GC.Cache.CachedHttpResponse(cacheKey)

			if ok {
				log.Printf("[%v], served from cache\n", cacheKey)
				c.JSON(status, response)
				return
			}
		}

		conDB.PrintDBStats(fmt.Sprintf("/v1/detail/users/%v", userIdentifier), s.P2P)

		var userInfo models.UserInfo

		cacheDurationInSeconds := 2 * 60 // 2 minutes

		userInfo, err = users.GetUser(userIdentifier, s.P2P)

		if err != nil {
			log.Println("[GET BUDS] error for user:", userIdentifier, "error: ", err)

			var ex p2pErrors.GenericError
			var ok bool

			ex, ok = err.(p2pErrors.GenericError)
			statusCode := 0
			var response interface{}

			if ok {
				statusCode = ex.HTTPCode()
				response = ex.JSONError()
			} else {
				statusCode = http.StatusBadRequest
				response = gin.H{"error": err.Error()}
			}

			c.JSON(statusCode, response)
			// gc.Cache.CacheHttpResponse(cacheKey, statusCode, response, cacheDurationInSeconds)
			return
		}
		// check wrong access
		ad, _ := middleware.ExtractTokenMetadata(c.Request)
		if ad.UserID != userInfo.ID {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized access"})
			return
		}
		s.GC.Cache.CacheHttpResponse(cacheKey, http.StatusOK, userInfo, cacheDurationInSeconds)

		c.JSON(http.StatusOK, userInfo)

	})

	// update User Contact Phone
	router.PUT("/v1/users/update", middleware.JwtTokenAuthMiddleware(s.AdminDB), func(c *gin.Context) {
		var err error
		ad, _ := middleware.ExtractTokenMetadata(c.Request)
		userID := ad.UserID
		userInfo, err := users.GetUser(userID, s.P2P)

		if err != nil {
			var ex p2pErrors.GenericError
			var ok bool

			ex, ok = err.(p2pErrors.GenericError)
			if ok {
				c.JSON(http.StatusBadRequest, ex.JSONError())
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}
			return
		}

		identifier := userInfo.Username

		conDB.PrintDBStats(fmt.Sprintf("PUT /v1/users/update %v", identifier), s.P2P)

		var userUpdateInfo models.UserUpdateInfo
		// var err error

		data, _ := io.ReadAll(c.Request.Body)

		err = json.Unmarshal(data, &userUpdateInfo)

		var invalidJSON p2pErrors.ErrorInvalidJSON

		if err != nil {
			c.JSON(http.StatusBadRequest, invalidJSON.JSONError())
			return
		}

		_, err = users.UpdateUser(identifier, middleware.ExtractPublicKey(c), userUpdateInfo, c.ClientIP(), s.P2P)

		if err != nil {
			var ex p2pErrors.GenericError
			var ok bool

			ex, ok = err.(p2pErrors.GenericError)
			if ok {
				c.JSON(http.StatusBadRequest, ex.JSONError())
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}
			return
		}

		// At this point, there was no error.
		//But either the email was sent or not.

		// c.JSON(http.StatusOK, gin.H{"username": c.Param("identifier")})
		userInfo, err = users.GetUser(identifier, s.P2P)

		if err != nil {
			var ex p2pErrors.GenericError
			var ok bool

			ex, ok = err.(p2pErrors.GenericError)
			if ok {
				c.JSON(ex.HTTPCode(), ex.JSONError())
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}
			return
		}
		{

			log.Println("[NotifyOrderListeners].....")

			s.GC.SendPNAndSocketUserChanged(userInfo)

		}
		cacheKey := fmt.Sprintf("[GET] /v1/detail/users/%v", identifier)
		cacheKeyAllOffers := fmt.Sprintf("[GET] /v1/offers %v", "ALL")

		s.GC.Cache.InvalidateCachedHttpResponse(cacheKey, cacheKeyAllOffers)

		c.JSON(http.StatusOK, userInfo)

	})

	// update KYC level
	router.PUT("/v1/admin/users/update/:targetUser", middleware.JwtTokenAuthMiddleware(s.AdminDB),
		accesslog.Audit(s.AdminDB, models.EventUserKycChange, models.AccessCategoryAccount, accesslog.Param("targetUser")),
		func(c *gin.Context) {
			var err error
			ad, _ := middleware.ExtractTokenMetadata(c.Request)
			userID := ad.UserID
			userInfo, err := users.GetUser(userID, s.P2P)
			if err != nil {
				log.Println("[] error for user:", userInfo.Username, "error: ", err)

				var ex p2pErrors.GenericError
				var ok bool

				ex, ok = err.(p2pErrors.GenericError)
				statusCode := 0
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
			// Check admin access
			if userInfo.AdminLevel == 0 {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access. Only admins can acess this endpoint"})
				return
			}

			identifier := strings.ReplaceAll(strings.ToLower(c.Param("targetUser")), " ", "")

			conDB.PrintDBStats(fmt.Sprintf("PUT /v1/admin/users/update %v", identifier), s.P2P)

			var userUpdateInfo models.KYCUpdateInfo
			// var err error

			data, _ := io.ReadAll(c.Request.Body)

			err = json.Unmarshal(data, &userUpdateInfo)

			var invalidJSON p2pErrors.ErrorInvalidJSON

			if err != nil {
				c.JSON(http.StatusBadRequest, invalidJSON.JSONError())
				return
			}
			userInfo, err = users.GetUser(identifier, s.P2P)

			if err != nil {
				var ex p2pErrors.GenericError
				var ok bool

				ex, ok = err.(p2pErrors.GenericError)
				if ok {
					c.JSON(http.StatusBadRequest, ex.JSONError())
				} else {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				}
				return
			}

			userKYCBefore := userInfo.KYCLevel
			_, err = users.UpdateUserKYC(identifier, middleware.ExtractPublicKey(c), userUpdateInfo, c.ClientIP(), s.P2P)

			if err != nil {
				var ex p2pErrors.GenericError
				var ok bool

				ex, ok = err.(p2pErrors.GenericError)
				if ok {
					c.JSON(http.StatusBadRequest, ex.JSONError())
				} else {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				}
				return
			}

			// At this point, there was no error.
			//But either the email was sent or not.

			// c.JSON(http.StatusOK, gin.H{"username": c.Param("identifier")})
			userInfo, err = users.GetUser(identifier, s.P2P)

			if err != nil {
				var ex p2pErrors.GenericError
				var ok bool

				ex, ok = err.(p2pErrors.GenericError)
				if ok {
					c.JSON(ex.HTTPCode(), ex.JSONError())
				} else {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				}
				return
			}
			{

				s.GC.SendPNAndSocketUserKYCChanged(userInfo, int(userKYCBefore))

				log.Println("[NotifyOrderListeners].....")

			}
			cacheKey := fmt.Sprintf("[GET] /v1/detail/users/%v", identifier)
			cacheKeyAllOffers := fmt.Sprintf("[GET] /v1/offers %v", "ALL")

			s.GC.Cache.InvalidateCachedHttpResponse(cacheKey, cacheKeyAllOffers)

			c.JSON(http.StatusOK, userInfo)

		})

	// Toggle Offline
	router.PUT("/v1/users/toggle", middleware.JwtTokenAuthMiddleware(s.AdminDB),
		accesslog.Audit(s.AdminDB, models.EventUserToggle, models.AccessCategoryAccount, nil),
		func(c *gin.Context) {
			var err error
			ad, _ := middleware.ExtractTokenMetadata(c.Request)
			userID := ad.UserID
			userInfo, err := users.GetUser(userID, s.P2P)

			if err != nil {
				log.Println("[TOGGLE USER STATE] error toggling user:", userInfo.Username, "error: ", err)

				var ex p2pErrors.GenericError
				var ok bool

				ex, ok = err.(p2pErrors.GenericError)
				statusCode := 0
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
			identifier := userInfo.Username
			if userInfo.KYCLevel == 0 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Offline/Online toggle is for merchants only."})
				return
			}

			conDB.PrintDBStats(fmt.Sprintf("PUT /v1/users/toggle %v current state: [%v]", identifier, userInfo.Offline), s.P2P)

			updatedState, err := userInfo.ToggleOffline(s.P2P)
			log.Println("[TOGGLE USER STATE] updated state toggling user:", userInfo.Username, "state: ", updatedState)

			if err != nil {
				var ex p2pErrors.GenericError
				var ok bool

				ex, ok = err.(p2pErrors.GenericError)
				if ok {
					c.JSON(http.StatusBadRequest, ex.JSONError())
				} else {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				}
				return
			}
			cacheKey := fmt.Sprintf("[GET] /v1/detail/users/%v", identifier)
			cacheKeyAllOffers := fmt.Sprintf("[GET] /v1/offers %v", "ALL")

			s.GC.Cache.InvalidateCachedHttpResponse(cacheKey, cacheKeyAllOffers)
			s.GC.Cache.InvalidateCachedHttpResponse(cacheKey, cacheKeyAllOffers)

			userInfo, _ = users.GetUser(userID, s.P2P)

			{

				log.Println("[NotifyOrderListeners].....")
				pl := struct {
					User models.UserInfo `json:"user"`
				}{
					User: userInfo,
				}
				s.GC.BroadcastToUserConnectionStreams(userInfo.Username, "userChanged", pl)
				if userInfo.Offline == 1 {
					pl := struct {
						State string `json:"state"`
					}{
						State: "offline",
					}
					s.GC.BroadcastToUserConnectionStreams(userInfo.Username, "userOffline", pl)

				} else {
					pl := struct {
						State string `json:"state"`
					}{
						State: "online",
					}
					s.GC.BroadcastToUserConnectionStreams(userInfo.Username, "userOnline", pl)

				}

			}

			c.JSON(http.StatusOK, userInfo)

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
		log.Println("Error sending discord message", err)
		return
	}
}
