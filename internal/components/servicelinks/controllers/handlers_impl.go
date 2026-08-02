package servicelinks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"
	tPayErrors "trovo-wallet-api/internal/components/payments/errors"
	paymentModels "trovo-wallet-api/internal/components/payments/models"
	paymentServices "trovo-wallet-api/internal/components/payments/services"

	servicelinkModels "trovo-wallet-api/internal/components/servicelinks/models"
	servicelinkServices "trovo-wallet-api/internal/components/servicelinks/services"
	usersDB "trovo-wallet-api/internal/components/users/db"
	userModels "trovo-wallet-api/internal/components/users/models"
	userServices "trovo-wallet-api/internal/components/users/services"
	conDB "trovo-wallet-api/internal/db"
	dl "trovo-wallet-api/internal/dynamiclinks"
	tErrors "trovo-wallet-api/internal/errors"
	pns "trovo-wallet-api/internal/pns"
	"trovo-wallet-api/internal/sharedconfig"

	"log"
	"net/http"
	"net/url"
	"strings"
	"trovo-wallet-api/internal/middleware"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stellar/go/keypair"
	"gorm.io/gorm/clause"

	"github.com/gin-gonic/gin"
)

// postServicelinksLoginRequestTargetUserHandler godoc
// @Summary POST /v1/servicelinks/login/request/:targetUser
// @Tags servicelinks
// @Accept json
// @Produce json
// @Param targetUser path string true "Target user identifier"
// @Param loginDescription query string false "Login description"
// @Param body body servicelinkModels.ServiceLinkRequestInput true "Login request payload"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /v1/servicelinks/login/request/{targetUser} [post]
func postServicelinksLoginRequestTargetUserHandler(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		mInfo, err := servicelinkServices.GetServiceLinkByAPIKey(middleware.ExtractServiceLinkApiKey(c), gc.DB)

		if err != nil {
			log.Println("[GET SERVICE] error for SERVICE:", middleware.ExtractServiceLinkApiKey(c), "error: ", err)

			var ex tErrors.GenericError
			var ok bool

			ex, ok = err.(tErrors.GenericError)
			var statusCode int = 0
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
		trovoUser := strings.TrimSpace(strings.ToLower(c.Param("targetUser")))
		if trovoUser == "null" {
			log.Printf("user cannot be %v\n", trovoUser)
			c.JSON(http.StatusBadRequest, gin.H{"error": "user cannot be null"})
			return
		}
		var serviceLinkRequestInput servicelinkModels.ServiceLinkRequestInput
		reqBody, _ := io.ReadAll(c.Request.Body)

		err = json.Unmarshal(reqBody, &serviceLinkRequestInput)

		var invalidJSON tErrors.ErrorInvalidJSON

		if err != nil {
			log.Println("Login Request Input JSON Error:", err)
			c.JSON(http.StatusBadRequest, invalidJSON.JSONError())
			return
		}
		identifier := mInfo.OwnerUsername
		if identifier == "" {
			log.Println("service account cannot be empty")
			c.JSON(http.StatusBadRequest, gin.H{"error": "service account cannot be empty"})
			return
		}
		conDB.PrintDBStats(fmt.Sprintf("POST /v1/servicelinks/login/request/%v %v/%v", trovoUser, identifier, middleware.ExtractServiceLinkApiKey(c)), gc.DB)

		if mInfo.LoginPermission == 0 {
			//wrong access
			statusCode := http.StatusUnauthorized
			response := gin.H{"error": "error-invalid-service-access", "data": "Permission", "message": "login permission not enabled for this service"}
			c.JSON(statusCode, response)
			return
		}

		userInfo, err := servicelinkServices.GetUserForServiceLink(trovoUser, mInfo, gc.DB, gc)

		if err != nil {
			log.Println("[GET UserInfo] error for user:", trovoUser, "error: ", err)

			var ex tErrors.GenericError
			var ok bool

			ex, ok = err.(tErrors.GenericError)
			var statusCode int = 0
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
		loginDescription := c.Query("loginDescription")
		//store login data for verification
		//does not exist. create new one.
		loginID := uuid.NewString()
		newLoginSession := servicelinkModels.ServiceLinkLoginSession{
			ApiKey:         middleware.ExtractServiceLinkApiKey(c),
			OwnerUsername:  identifier,
			WalletUsername: userInfo.Username,
			ID:             loginID,
		}

		if len(serviceLinkRequestInput.CallbackURL) > 0 {
			newLoginSession.CallbackURL = &serviceLinkRequestInput.CallbackURL
		}

		err = gc.DB.Omit(clause.Associations).Create(&newLoginSession).Error
		if err != nil {
			//could not create login session
			response := gin.H{"error": "error-temporary-server-error", "data": "temporaryServerError", "message": "Temporary Server Error. Contact support."}
			statusCode := http.StatusServiceUnavailable
			c.JSON(statusCode, response)
			return
		}

		//respond with deep-link and QRCode for login
		data, err := dl.GenerateLoginData(mInfo.OwnerUsername, mInfo.ShortName, userInfo.Username, loginID, serviceLinkRequestInput.DeviceInfo, loginDescription, gc)
		if err != nil {
			//could not create login session
			response := gin.H{"error": "error-temporary-server-error", "data": "temporaryServerError", "message": "Temporary Server Error. Contact support."}
			statusCode := http.StatusServiceUnavailable
			c.JSON(statusCode, response)
			return
		}

		// if userInfo.PushNotificationToken != nil {
		// 	dataPayload := make(map[string]string)
		// 	dataPayload["link"] = data.DynamicLink
		// 	pns.SendFirebaseMessage(*userInfo.PushNotificationToken, fmt.Sprintf("Login for [%v] requested!", userInfo.Username), fmt.Sprintf("Your TrovoApp username [%v] has been used to request a login session on [%v] service using [%v]. Click to continue.", userInfo.Username, mInfo.LongName, serviceLinkRequestInput.DeviceInfo), "", dataPayload, gc.PushNotificationClient, gc.PNSContext)
		// }

		c.JSON(http.StatusOK, data)

	}
}

// postUsersServicelinksLoginApprovalTargetUserHandler godoc
// @Summary POST /v1/users/servicelinks/login/approval/:targetUser
// @Tags servicelinks
// @Produce json
// @Param targetUser path string true "Target user identifier"
// @Param ownerUsername query string true "Owner username"
// @Param loginId query string true "Login ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /v1/users/servicelinks/login/approval/{targetUser} [post]
func postUsersServicelinksLoginApprovalTargetUserHandler(callBackRetryChan chan retryCallbacks, gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {

		targetUser := strings.TrimSpace(strings.ToLower(c.Param("targetUser")))
		ownerUsername := strings.TrimSpace(strings.ToLower(c.Query("ownerUsername")))
		loginID := strings.TrimSpace(strings.ToLower(c.Query("loginId")))
		if targetUser == "null" || targetUser == "" {
			log.Printf("user cannot be %v\n", targetUser)
			c.JSON(http.StatusBadRequest, gin.H{"error": "user cannot be null"})
			return
		}
		if len(ownerUsername) == 0 {
			log.Printf("service owner cannot be empty%v\n", targetUser)
			c.JSON(http.StatusBadRequest, gin.H{"error": "service cannot be empty"})
			return
		}
		if len(loginID) == 0 {
			log.Printf("loginID cannot be empty%v\n", targetUser)
			c.JSON(http.StatusBadRequest, gin.H{"error": "loginId cannot be empty"})
			return
		}
		//store login data for verification
		loginSession, err := servicelinkServices.GetLoginSession(ownerUsername, targetUser, loginID, gc.DB)
		if err != nil {

			//other system error
			var ex tErrors.GenericError
			var ok bool

			ex, ok = err.(tErrors.GenericError)
			var statusCode int = 0
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
		conDB.PrintDBStats(fmt.Sprintf("POST /v1/users/servicelinks/login/approval/%v %v/%v", targetUser, ownerUsername, loginID), gc.DB)
		mInfo, err := servicelinkServices.GetServiceLinkByAPIKey(loginSession.ApiKey, gc.DB)

		if err != nil {
			log.Println("[GET SERVICE INFO] error for SERVICE:", ownerUsername, loginSession.ApiKey, "error: ", err)

			var ex tErrors.GenericError
			var ok bool

			ex, ok = err.(tErrors.GenericError)
			var statusCode int = 0
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

		if mInfo.LoginPermission == 0 {
			//wrong access
			statusCode := http.StatusUnauthorized
			response := gin.H{"error": "error-invalid-service-access", "data": "Permission", "message": "login permission not enabled for this service"}
			c.JSON(statusCode, response)
			return
		}

		//validate ownership
		if mInfo.OwnerUsername != ownerUsername {
			//wrong access
			statusCode := http.StatusUnauthorized
			response := gin.H{"error": "error-invalid-user-access", "data": "Authentication", "message": "this login request does not belong to service owner specified"}
			c.JSON(statusCode, response)
			return
		}

		userInfo, err := servicelinkServices.GetUserFromPrimarySigner(middleware.ExtractSigner(c), gc.DB, gc)

		if err != nil {
			log.Println("[GET UserInfo] error for signer:", middleware.ExtractSigner(c), "error: ", err)

			var ex tErrors.GenericError
			var ok bool

			ex, ok = err.(tErrors.GenericError)
			var statusCode int = 0
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
		//validate ownership
		if targetUser != userInfo.Username {
			//wrong access
			statusCode := http.StatusUnauthorized
			response := gin.H{"error": "error-invalid-user-access", "data": "Authentication", "message": "this login request does not belong to your TrovoApp"}
			c.JSON(statusCode, response)
			return
		}

		if userInfo.Mobile == nil {
			log.Println("[GET UserInfo] error for user:", targetUser, "error: user does not have a valid phone number")
			err = &tErrors.CustomError{
				Param:      "mobile",
				Err:        "error no valid phone number",
				ErrMessage: "Only users with valid phone number are allows to use this service. Please update your mobile number.",
				Code:       http.StatusBadRequest,
			}
			var ex tErrors.GenericError
			var ok bool

			ex, ok = err.(tErrors.GenericError)
			var statusCode int = 0
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

		if loginSession.Authorized == 1 {
			response := gin.H{"error": "error-login-session-does-not-exist", "data": userInfo.Username, "message": "invalid/expired login request"}
			statusCode := http.StatusNotFound
			c.JSON(statusCode, response)
			return
		}
		//exists and needs to be updated
		loginSession.Authorized = 1

		err = gc.DB.Omit(clause.Associations).Save(&loginSession).Error
		if err != nil {
			//could not save login session
			response := gin.H{"error": "error-temporary-server-error", "data": "temporaryServerError", "message": "Temporary Server Error. Contact support."}
			statusCode := http.StatusServiceUnavailable
			c.JSON(statusCode, response)
			return
		}

		cacheKey := fmt.Sprintf("[GET] /v1/users/servicelinks/login/approval/%v %v", userInfo.Username, loginID)
		gc.RedisCache.InvalidateCachedHttpResponse(cacheKey)

		if userInfo.PushNotificationToken != nil {
			dataPayload := make(map[string]string)
			dataPayload["route"] = ""
			pns.SendFirebaseMessage(*userInfo.PushNotificationToken, fmt.Sprintf("Login for %v authorized!", userInfo.Username), fmt.Sprintf("Your TrovoApp username %v has been authorized to login on %v service.", userInfo.Username, mInfo.LongName), "", dataPayload, gc.PushNotificationClient, gc.PNSContext)
		}

		//return response to user and  not keep them waiting.
		c.JSON(http.StatusOK, gin.H{"message": "success"})

		//make callback request if callback is available

		if loginSession.CallbackURL != nil {
			//make callback request
			// callbackResponse := new(map[string]interface{})
			type payload struct {
				LoginID    string `json:"loginId"`
				TargetUser string `json:"targetUser"`
			}
			jsonPayload := payload{LoginID: loginSession.ID, TargetUser: loginSession.WalletUsername}

			/////
			d := *loginSession.CallbackURL
			body, err := json.Marshal(jsonPayload)
			if err != nil {
				log.Printf("[LoginAuthCallback] could not unmarshal callback message due to [%v]\n", err)

			}
			log.Printf("[LoginAuthCallback] JSON STRING: [%v]\n", string(body))

			responseBody := bytes.NewBuffer(body)
			//Leverage Go's HTTP Post function to make request
			resp, err := http.Post(d, "application/json", responseBody)
			//Handle Error
			if err != nil {
				log.Printf("[LoginAuthCallback] could not send callback message due to [%v]\n", err)
				return
			}
			defer resp.Body.Close()
			//Read the response body
			body, err = io.ReadAll(resp.Body)
			if err != nil {
				//send to retry channel
				log.Println("[LoginAuthCallback] callback failed:", err)
				c := retryCallbacks{Req: responseBody, CallbackURL: d, Count: 0}
				callBackRetryChan <- c
			} else {
				log.Printf("[LoginAuthCallback] Login authorization Callback successful to: [%v], Response:[%v]\n\n", d, string(body))

			}
		}
		cacheKey = fmt.Sprintf("GET /v1/servicelinks/login/verify/%v/%v/%v", mInfo.OwnerUsername, userInfo.Username, loginID)
		gc.RedisCache.InvalidateCachedHttpResponse(cacheKey)

	}
}

// getServicelinksLoginVerifyOwnerUsernameTargetUserLoginIDHandler godoc
// @Summary GET /v1/servicelinks/login/verify/:ownerUsername/:targetUser/:loginID
// @Tags servicelinks
// @Produce json
// @Param ownerUsername path string true "Owner username"
// @Param targetUser path string true "Target user"
// @Param loginID path string true "Login ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /v1/servicelinks/login/verify/{ownerUsername}/{targetUser}/{loginID} [get]
func getServicelinksLoginVerifyOwnerUsernameTargetUserLoginIDHandler(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {

		ownerUsername := strings.TrimSpace(strings.ToLower(c.Param("ownerUsername")))
		trovoUser := strings.TrimSpace(strings.ToLower(c.Param("targetUser")))
		loginID := strings.TrimSpace(strings.ToLower(c.Param("loginID")))
		if trovoUser == "null" {
			log.Printf("user cannot be %v\n", trovoUser)
			c.JSON(http.StatusBadRequest, gin.H{"error": "user cannot be null"})
			return
		}
		conDB.PrintDBStats(fmt.Sprintf("GET /v1/servicelinks/login/verify/%v/%v/%v", ownerUsername, trovoUser, loginID), gc.DB)

		cacheKey := fmt.Sprintf("GET /v1/servicelinks/login/verify/%v/%v/%v", ownerUsername, trovoUser, loginID)
		{
			//search cache

			ok, status, response := gc.RedisCache.CachedHttpResponse(cacheKey)

			if ok {
				// log.Printf("[%v], served from cache\n", cacheKey)
				c.JSON(status, response)
				return
			}
		}

		mInfo, err := servicelinkServices.GetServiceLinkByAPIKey(middleware.ExtractServiceLinkApiKey(c), gc.DB)

		if err != nil {
			log.Println("[GET SERVICE FOR USER LOGIN] error for servicelink:", middleware.ExtractServiceLinkApiKey(c), "error: ", err)

			var ex tErrors.GenericError
			var ok bool

			ex, ok = err.(tErrors.GenericError)
			var statusCode int = 0
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
		if mInfo.OwnerUsername != ownerUsername {
			//wrong access
			statusCode := http.StatusUnauthorized
			response := gin.H{"error": "error-invalid-service-access", "data": "Authentication", "message": "Authentication failed"}
			c.JSON(statusCode, response)
			cacheDurationInSeconds := 60 //1 minutes

			gc.RedisCache.CacheHttpResponse(cacheKey, statusCode, response, cacheDurationInSeconds)
			return
		}
		// log.Printf("service Info: %+v\n", mInfo)
		if mInfo.LoginPermission == 0 {
			//wrong access
			statusCode := http.StatusUnauthorized
			response := gin.H{"error": "error-invalid-service-access", "data": "Permission", "message": "login permission not enabled for this service"}
			c.JSON(statusCode, response)
			cacheDurationInSeconds := 60 //1 minutes

			gc.RedisCache.CacheHttpResponse(cacheKey, statusCode, response, cacheDurationInSeconds)
			return
		}

		userInfo, err := servicelinkServices.GetUserForServiceLink(trovoUser, mInfo, gc.DB, gc)

		if err != nil {
			log.Println("[GET UserInfo] error for user:", trovoUser, "error: ", err)

			var ex tErrors.GenericError
			var ok bool

			ex, ok = err.(tErrors.GenericError)
			var statusCode int = 0
			var response interface{}

			if ok {
				statusCode = ex.HTTPCode()
				response = ex.JSONError()
			} else {
				statusCode = http.StatusBadRequest
				response = gin.H{"error": err.Error()}
			}

			c.JSON(statusCode, response)
			cacheDurationInSeconds := 5 //1 minutes

			gc.RedisCache.CacheHttpResponse(cacheKey, statusCode, response, cacheDurationInSeconds)
			return
		}

		//get login data for verification
		loginSession, err := servicelinkServices.GetLoginSession(ownerUsername, userInfo.Username, loginID, gc.DB)
		if err != nil {
			log.Printf("[error Verifying Login] for user [%v], error [%v]]\n", trovoUser, err)
			var ex tErrors.GenericError
			var ok bool

			ex, ok = err.(tErrors.GenericError)
			var statusCode int = 0
			var response interface{}

			if ok {
				statusCode = ex.HTTPCode()
				response = ex.JSONError()
			} else {
				statusCode = http.StatusBadRequest
				response = gin.H{"error": err.Error()}
			}

			cacheDurationInSeconds := 60 //1 minute

			gc.RedisCache.CacheHttpResponse(cacheKey, statusCode, response, cacheDurationInSeconds)

			c.JSON(statusCode, response)
			return

		}
		if loginSession.Authorized == 0 {

			//could not save login session
			response := gin.H{"error": "error-unauthorized-login-session", "data": "unauthorizedLoginSession", "message": "login session awaiting authorization"}
			statusCode := http.StatusUnauthorized

			cacheDurationInSeconds := 60 //1 minutes

			gc.RedisCache.CacheHttpResponse(cacheKey, statusCode, response, cacheDurationInSeconds)

			c.JSON(statusCode, response)
			return

		}

		data, err := userServices.LogUserIn(userInfo, gc)
		if err != nil {
			log.Printf("[error Verifying Login] for user [%v], error [%v]]\n", trovoUser, err)
			var ex tErrors.GenericError
			var ok bool

			ex, ok = err.(tErrors.GenericError)
			var statusCode int = 0
			var response interface{}

			if ok {
				statusCode = ex.HTTPCode()
				response = ex.JSONError()
			} else {
				statusCode = http.StatusBadRequest
				response = gin.H{"error": err.Error()}
			}

			cacheDurationInSeconds := 60 //1 minute

			gc.RedisCache.CacheHttpResponse(cacheKey, statusCode, response, cacheDurationInSeconds)

			c.JSON(statusCode, response)
			return

		}
		cacheDurationInSeconds := 2 * 60 //2 minutes

		gc.RedisCache.CacheHttpResponse(cacheKey, http.StatusOK, data, cacheDurationInSeconds)

		c.JSON(http.StatusOK, data)
	}
}

// postServicelinksTokenRefreshHandler godoc
// @Summary POST /v1/servicelinks/token/refresh
// @Tags servicelinks
// @Accept json
// @Produce json
// @Param body body object true "Refresh token payload with refreshToken field"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /v1/servicelinks/token/refresh [post]
func postServicelinksTokenRefreshHandler(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {

		type MapToken struct {
			RefreshToken string `json:"refreshToken"`
		}
		var mapToken MapToken
		reqBody, _ := io.ReadAll(c.Request.Body)

		err := json.Unmarshal(reqBody, &mapToken)

		// var invalidJSON tErrors.ErrorInvalidJSON

		if err != nil {
			log.Printf("[HandleRefreshToken] Login Request Input JSON Error:%v\nBody:%v\n", err, string(reqBody))
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "error-unable-to-retrieve-refresh-token", "message": "Unable to retrieve refresh token."})
			return
		}
		if len(mapToken.RefreshToken) == 0 {
			log.Println("[HandleRefreshToken] refresh token is empty...")
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "error-empty-refresh-token", "message": "Refresh token is empty"})
			return
		}
		refreshToken := mapToken.RefreshToken

		//verify the token
		refreshReponse, err := userServices.RefreshToken(refreshToken, gc)
		if err != nil { //if any goes wrong
			log.Printf("[Refresh Handle] error loging in: %v\n", err)
			var ex tErrors.GenericError
			var ok bool

			ex, ok = err.(tErrors.GenericError)
			var statusCode int = 0
			var response interface{}

			if ok {
				statusCode = ex.HTTPCode()
				response = ex.JSONError()
			} else {
				statusCode = http.StatusUnauthorized
				response = gin.H{"error": "unauthorized to perform this action"}
			}

			c.JSON(statusCode, response)
			return
		}

		c.JSON(http.StatusCreated, refreshReponse)

	}
}

// postServicelinksTokenVerifyHandler godoc
// @Summary POST /v1/servicelinks/token/verify
// @Tags servicelinks
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /v1/servicelinks/token/verify [post]
func postServicelinksTokenVerifyHandler(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := middleware.VerifyToken(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token is invalid"})
			return
		}
		resp := struct {
			Message string     `json:"message"`
			Token   *jwt.Token `json:"token"`
		}{
			Message: "success",
			Token:   token,
		}

		c.JSON(http.StatusOK, resp)

	}
}

// deleteServicelinksTokenHandler godoc
// @Summary DELETE /v1/servicelinks/token
// @Tags servicelinks
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /v1/servicelinks/token [delete]
func deleteServicelinksTokenHandler(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		au, err := middleware.ExtractTokenMetadata(c.Request)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"message": "successfully logged out"})
			return
		}

		deleted, delErr := middleware.DeleteAuth(au.AccessUUID, gc.RedisCache)
		if delErr != nil || deleted == 0 { //if any goes wrong
			c.JSON(http.StatusOK, gin.H{"message": "successfully logged out"})
			return

		}

	}
}

// postServicelinksAuthorizeRequestTargetUserHandler godoc
// @Summary POST /v1/servicelinks/authorize/request/:targetUser
// @Tags servicelinks
// @Accept json
// @Produce json
// @Param targetUser path string true "Target user identifier"
// @Param body body servicelinkModels.ServiceLinkRequestInput true "Authorize request payload"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /v1/servicelinks/authorize/request/{targetUser} [post]
func postServicelinksAuthorizeRequestTargetUserHandler(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		trovoUser := strings.TrimSpace(strings.ToLower(c.Param("targetUser")))

		if trovoUser == "null" {
			log.Printf("user cannot be %v\n", trovoUser)
			c.JSON(http.StatusBadRequest, gin.H{"error": "user cannot be null"})
			return
		}

		mInfo, err := servicelinkServices.GetServiceLinkByAPIKey(middleware.ExtractServiceLinkApiKey(c), gc.DB)

		if err != nil {
			log.Println("[GET service] error for service:", middleware.ExtractServiceLinkApiKey(c), "error: ", err)

			var ex tErrors.GenericError
			var ok bool

			ex, ok = err.(tErrors.GenericError)
			var statusCode int = 0
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

		ownerUsername := mInfo.OwnerUsername

		conDB.PrintDBStats(fmt.Sprintf("POST /v1/servicelinks/authorize/request/%v %v", trovoUser, ownerUsername), gc.DB)

		if mInfo.AuthorizationPermission == 0 {
			//wrong access
			statusCode := http.StatusUnauthorized
			response := gin.H{"error": "error-invalid-service-access", "data": "Permission", "message": "authorization permission not enabled for this service"}
			c.JSON(statusCode, response)
			return
		}

		userInfo, err := servicelinkServices.GetUserForServiceLink(trovoUser, mInfo, gc.DB, gc)

		if err != nil {
			log.Println("[GET UserInfo] error for user:", trovoUser, "error: ", err)

			var ex tErrors.GenericError
			var ok bool

			ex, ok = err.(tErrors.GenericError)
			var statusCode int = 0
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

		var serviceLinkRequestInput servicelinkModels.ServiceLinkRequestInput
		reqBody, _ := io.ReadAll(c.Request.Body)

		err = json.Unmarshal(reqBody, &serviceLinkRequestInput)

		var invalidJSON tErrors.ErrorInvalidJSON

		if err != nil {
			log.Println("Authorization Request Input JSON Error:", err)
			c.JSON(http.StatusBadRequest, invalidJSON.JSONError())
			return
		}

		authID := uuid.NewString()
		authorizationData := servicelinkModels.ServiceLinkAuthorization{

			ID:             authID,
			ApiKey:         middleware.ExtractServiceLinkApiKey(c),
			OwnerUsername:  mInfo.OwnerUsername,
			WalletUsername: userInfo.Username,
		}
		if len(serviceLinkRequestInput.CallbackURL) > 0 {
			authorizationData.CallbackURL = &serviceLinkRequestInput.CallbackURL
		}
		period := time.Duration(3)
		if os.Getenv("SERVICE_LINK_AUTHORIZATION_REQUEST_VALIDITY") != "" {
			m, e := decimal.NewFromString(os.Getenv("SERVICE_LINK_AUTHORIZATION_REQUEST_VALIDITY"))
			if e == nil {
				if m.IsPositive() {
					period = time.Duration(m.IntPart())
				}
			}
		}
		if serviceLinkRequestInput.ValidityInMinutes > 0 {
			//check if validtity was submitted
			period = time.Duration(serviceLinkRequestInput.ValidityInMinutes)

		}
		authorizationData.ExpiresAt = time.Now().Add(period * time.Minute)
		err = gc.DB.Omit(clause.Associations).Create(&authorizationData).Error
		if err != nil {
			log.Printf("unable to create authorization: %v\n", err)
			//could not save login session
			response := gin.H{"error": "error-temporary-server-error", "data": "temporaryServerError", "message": "Temporary Server Error. Contact support."}
			statusCode := http.StatusServiceUnavailable
			c.JSON(statusCode, response)
			return
		}

		data, err := dl.GenerateAuthorizationData(mInfo.OwnerUsername, mInfo.ShortName, serviceLinkRequestInput.AuthDescription, userInfo.Username, serviceLinkRequestInput.DeviceInfo, authID, gc)
		if err != nil {
			//could not create authorization session
			response := gin.H{"error": "error-temporary-server-error", "data": "temporaryServerError", "message": "Temporary Server Error. Contact support."}
			statusCode := http.StatusServiceUnavailable
			c.JSON(statusCode, response)
			return
		}
		// if userInfo.PushNotificationToken != nil {
		// 	dataPayload := make(map[string]string)
		// 	dataPayload["route"] = data.DynamicLink
		// 	pns.SendFirebaseMessage(*userInfo.PushNotificationToken, fmt.Sprintf("Authorization for trovo account %v requested!", userInfo.Username), fmt.Sprintf("Your TrovoApp username [%v] has been used to request an authorization session on [%v] service using [%v]. Click to continue.", userInfo.Username, mInfo.LongName, serviceLinkRequestInput.DeviceInfo), "", dataPayload, gc.PushNotificationClient, gc.PNSContext)
		// }
		c.JSON(http.StatusOK, data)
	}
}

// postServicelinksAuthorizeTokenizedAssetHandler godoc
// @Summary POST /v1/servicelinks/authorize/tokenized-asset
// @Tags servicelinks
// @Accept json
// @Produce json
// @Param body body servicelinkModels.ServiceLinkTokenizedAssetAuthRequestInput true "Tokenized asset auth request payload"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /v1/servicelinks/authorize/tokenized-asset [post]
func postServicelinksAuthorizeTokenizedAssetHandler(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {

		mInfo, err := servicelinkServices.GetServiceLinkByAPIKey(middleware.ExtractServiceLinkApiKey(c), gc.DB)

		if err != nil {
			log.Println("[GET service] error for service:", middleware.ExtractServiceLinkApiKey(c), "error: ", err)

			var ex tErrors.GenericError
			var ok bool

			ex, ok = err.(tErrors.GenericError)
			var statusCode int = 0
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

		ownerUsername := mInfo.OwnerUsername

		conDB.PrintDBStats(fmt.Sprintf("POST /v1/servicelinks/authorize/tokenized-asset %v", ownerUsername), gc.DB)

		if mInfo.TokenizedAssetAuthorizationPermission == 0 {
			//wrong access
			statusCode := http.StatusUnauthorized
			response := gin.H{"error": "error-invalid-service-access", "data": "Permission", "message": "Permission to authorize tokenized asset not enabled for this service"}
			c.JSON(statusCode, response)
			return
		}

		var serviceLinkRequestInput servicelinkModels.ServiceLinkTokenizedAssetAuthRequestInput
		reqBody, _ := io.ReadAll(c.Request.Body)

		err = json.Unmarshal(reqBody, &serviceLinkRequestInput)

		var invalidJSON tErrors.ErrorInvalidJSON

		if err != nil {
			log.Println("ServiceLinkTokenizedAssetAuthRequestInput Request Input JSON Error:", err)
			c.JSON(http.StatusBadRequest, invalidJSON.JSONError())
			return
		}
		if serviceLinkRequestInput.AssetCode == "" {
			log.Println("ServiceLinkTokenizedAssetAuthRequestInput Request Input Asset Code Error")
			c.JSON(http.StatusBadRequest, gin.H{"error": "error-no-asset-code", "message": "Asset Code is empty"})
			return
		}

		if serviceLinkRequestInput.UnsignedTransaction == "" {
			log.Println("ServiceLinkTokenizedAssetAuthRequestInput Request Input Transaction Error")
			c.JSON(http.StatusBadRequest, gin.H{"error": "error-no-unsignedTransaction", "message": "UnsignedTransaction is empty"})
			return
		}
		// get tokenized asset by the asset code
		t := gc.GetTokenizedAssetByCode(serviceLinkRequestInput.AssetCode)

		if t.ID == "" {
			log.Println("ServiceLinkTokenizedAssetAuthRequestInput Request Input Asset Code Error")
			c.JSON(http.StatusBadRequest, gin.H{"error": "error-no-asset-code", "message": "Asset Code is not for tokenized asset"})
			return
		}

		data, err := servicelinkServices.GetTransactionSignature(&serviceLinkRequestInput, gc)
		if err != nil {
			//could not create authorization session
			response := gin.H{"error": "error-temporary-server-error", "data": "temporaryServerError", "message": "Temporary Server Error. Contact support."}
			statusCode := http.StatusServiceUnavailable
			c.JSON(statusCode, response)
			return
		}

		c.JSON(http.StatusOK, data)
	}
}

// postServicelinksEventsRequestHandler godoc
// @Summary POST /v1/servicelinks/events/request
// @Tags servicelinks
// @Accept json
// @Produce json
// @Param body body servicelinkModels.ServiceLinkEventRequestInput true "Event request payload"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /v1/servicelinks/events/request [post]
func postServicelinksEventsRequestHandler(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {

		mInfo, err := servicelinkServices.GetServiceLinkByAPIKey(middleware.ExtractServiceLinkApiKey(c), gc.DB)

		if err != nil {
			log.Println("[GET service] error for service:", middleware.ExtractServiceLinkApiKey(c), "error: ", err)

			var ex tErrors.GenericError
			var ok bool

			ex, ok = err.(tErrors.GenericError)
			var statusCode int = 0
			var response interface{}

			if ok {
				statusCode = ex.HTTPCode()
				response = ex.JSONError()
			} else {
				statusCode = http.StatusBadRequest
				response = gin.H{"error": err.Error(), "message": err.Error()}
			}

			c.JSON(statusCode, response)
			return
		}

		ownerUsername := mInfo.OwnerUsername

		conDB.PrintDBStats(fmt.Sprintf("POST /v1/servicelinks/events/request %v", ownerUsername), gc.DB)

		if mInfo.EventPermission == 0 {
			//wrong access
			statusCode := http.StatusUnauthorized
			response := gin.H{"error": "error-invalid-service-access", "data": "Permission", "message": "Event link generation permission not enabled for this service"}
			c.JSON(statusCode, response)
			return
		}

		var serviceLinkRequestInput servicelinkModels.ServiceLinkEventRequestInput
		reqBody, _ := io.ReadAll(c.Request.Body)

		err = json.Unmarshal(reqBody, &serviceLinkRequestInput)

		var invalidJSON tErrors.ErrorInvalidJSON

		if err != nil {
			log.Println("Event Request Input JSON Error:", err)
			c.JSON(http.StatusBadRequest, invalidJSON.JSONError())
			return
		}

		eventID := uuid.NewString()
		authorizationData := servicelinkModels.ServiceLinkEvent{
			ID:            eventID,
			ApiKey:        middleware.ExtractServiceLinkApiKey(c),
			OwnerUsername: mInfo.OwnerUsername,
		}
		if len(serviceLinkRequestInput.CallbackURL) > 0 {
			authorizationData.CallbackURL = &serviceLinkRequestInput.CallbackURL
		}
		period := time.Duration(3)
		if os.Getenv("SERVICE_LINK_AUTHORIZATION_REQUEST_VALIDITY") != "" {
			m, e := decimal.NewFromString(os.Getenv("SERVICE_LINK_AUTHORIZATION_REQUEST_VALIDITY"))
			if e == nil {
				if m.IsPositive() {
					period = time.Duration(m.IntPart())
				}
			}
		}
		if serviceLinkRequestInput.ValidityInMinutes > 0 {
			//check if validtity was submitted
			period = time.Duration(serviceLinkRequestInput.ValidityInMinutes)

		}
		authorizationData.ExpiresAt = time.Now().Add(period * time.Minute)
		err = gc.DB.Omit(clause.Associations).Create(&authorizationData).Error
		if err != nil {
			log.Printf("unable to create event: %v\n", err)
			//could not save login session
			response := gin.H{"error": "error-temporary-server-error", "data": "temporaryServerError", "message": "Temporary Server Error. Contact support."}
			statusCode := http.StatusServiceUnavailable
			c.JSON(statusCode, response)
			return
		}

		data, err := dl.GenerateEventData(mInfo.OwnerUsername, mInfo.ShortName, serviceLinkRequestInput.EventDescription, serviceLinkRequestInput.DeviceInfo, eventID, gc)
		if err != nil {
			//could not create authorization session
			response := gin.H{"error": "error-temporary-server-error", "data": "temporaryServerError", "message": "Temporary Server Error. Contact support."}
			statusCode := http.StatusServiceUnavailable
			c.JSON(statusCode, response)
			return
		}

		c.JSON(http.StatusOK, data)
	}
}

// postUsersServicelinksAuthorizeApprovalTargetUserHandler godoc
// @Summary POST /v1/users/servicelinks/authorize/approval/:targetUser
// @Tags servicelinks
// @Produce json
// @Param targetUser path string true "Target user identifier"
// @Param ownerUsername query string true "Owner username"
// @Param authId query string true "Authorization ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /v1/users/servicelinks/authorize/approval/{targetUser} [post]
func postUsersServicelinksAuthorizeApprovalTargetUserHandler(callBackRetryChan chan retryCallbacks, gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {

		identifier := strings.TrimSpace(strings.ToLower(c.Param("targetUser")))
		if identifier == "null" {
			log.Printf("user cannot be %v\n", identifier)
			c.JSON(http.StatusBadRequest, gin.H{"error": "user cannot be null"})
			return
		}
		ownerUsername := strings.TrimSpace(strings.ToLower(c.Query("ownerUsername")))
		authID := strings.TrimSpace(c.Query("authId"))
		//get authorization data for user
		authData, err := servicelinkServices.GetUserAuthorizationData(ownerUsername, identifier, authID, gc.DB)
		if err != nil {

			//other system error
			var ex tErrors.GenericError
			var ok bool

			ex, ok = err.(tErrors.GenericError)
			var statusCode int = 0
			var response interface{}

			if ok {
				statusCode = ex.HTTPCode()
				response = ex.JSONError()
			} else {
				statusCode = http.StatusBadRequest
				response = gin.H{"error": err.Error()}
			}
			log.Println("Event authorization error:", response)
			c.JSON(statusCode, response)
			return

		}
		conDB.PrintDBStats(fmt.Sprintf("POST /v1/users/servicelinks/authorize/%v %v/%v", identifier, ownerUsername, authID), gc.DB)

		mInfo, err := servicelinkServices.GetServiceLinkByAPIKey(authData.ApiKey, gc.DB)

		if err != nil {
			log.Println("[GET SERVICE] error for SERVICELINK:", ownerUsername, "error: ", err)

			var ex tErrors.GenericError
			var ok bool

			ex, ok = err.(tErrors.GenericError)
			var statusCode int = 0
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

		if mInfo.AuthorizationPermission == 0 {
			//wrong access

			statusCode := http.StatusUnauthorized
			response := gin.H{"error": "error-invalid-service-access", "data": "Permission", "message": "authorization permission not enabled for this service"}
			log.Printf("%+v\n", response)
			c.JSON(statusCode, response)
			return
		}

		userInfo, err := servicelinkServices.GetUserFromPrimarySigner(middleware.ExtractSigner(c), gc.DB, gc)

		if err != nil {
			log.Println("[GET UserInfo] error for user:", identifier, "error: ", err)

			var ex tErrors.GenericError
			var ok bool

			ex, ok = err.(tErrors.GenericError)
			var statusCode int = 0
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
		//check if service is for an event registration/reward service: [2 = registration, 1 = reward, 0 = none]

		//2FA must belong to the caller account

		if identifier != userInfo.Username {
			//wrong access
			statusCode := http.StatusUnauthorized
			response := gin.H{"error": "error-invalid-user-access", "data": "Authentication", "message": "2FA/Authorization request does not belong to your TrovoApp"}
			c.JSON(statusCode, response)
			return
		}
		//get authorization data for user

		if authData.Authorized == 1 {
			response := gin.H{"error": "error-authorization-does-not-exist", "data": userInfo.Username, "message": "invalid/expired authorization request"}
			statusCode := http.StatusNotFound
			c.JSON(statusCode, response)
			return
		}
		//exists and needs to be updated
		authData.Authorized = 1

		err = gc.DB.Omit(clause.Associations).Save(&authData).Error
		if err != nil {
			//could not save authData
			response := gin.H{"error": "error-temporary-server-error", "data": "temporaryServerError", "message": "Temporary Server Error. Contact support."}
			statusCode := http.StatusServiceUnavailable
			c.JSON(statusCode, response)
			return
		}

		if userInfo.PushNotificationToken != nil {
			dataPayload := make(map[string]string)
			dataPayload["route"] = ""
			pns.SendFirebaseMessage(*userInfo.PushNotificationToken, fmt.Sprintf("2FA Action for trovo account %v authorized!", userInfo.Username), fmt.Sprintf("Your TrovoApp username %v has been used to authorize a 2FA action on %v service.", userInfo.Username, mInfo.LongName), "", dataPayload, gc.PushNotificationClient, gc.PNSContext)
		}
		//return report to user and not keep them waiting.
		c.JSON(http.StatusOK, gin.H{"message": "success"})

		if authData.CallbackURL != nil {

			//make callback request
			// callbackResponse := new(map[string]interface{})
			type payload struct {
				AuthID     string `json:"authId"`
				TargetUser string `json:"targetUser"`
				DeviceID   string `json:"deviceId"`
			}
			jsonPayload := payload{AuthID: authData.ID, TargetUser: authData.WalletUsername, DeviceID: c.GetHeader("X-TW-DEVICE-ID")}

			/////
			d := *authData.CallbackURL
			body, err := json.Marshal(jsonPayload)
			if err != nil {
				log.Printf("[AuthCallback] could not unmarshal callback message due to [%v]\n", err)

			}
			log.Printf("[AuthCallback] JSON STRING: [%v]\n", string(body))

			responseBody := bytes.NewBuffer(body)
			//Leverage Go's HTTP Post function to make request
			resp, err := http.Post(d, "application/json", responseBody)
			//Handle Error
			if err != nil {
				log.Printf("[AuthCallback] could not send callback message due to [%v]\n", err)
				return
			}
			defer resp.Body.Close()
			//Read the response body
			body, err = io.ReadAll(resp.Body)
			if err != nil {
				//send to retry channel
				log.Println("[AuthCallback] callback failed:", err)
				c := retryCallbacks{Req: responseBody, CallbackURL: d, Count: 0}
				callBackRetryChan <- c
			} else {
				log.Printf("[AuthCallback] 2FA authorization Callback successful to: [%v], Response:[%v]\n\n", d, string(body))

			}

		}

	}
}

// postUsersServicelinksEventsApprovalTargetUserHandler godoc
// @Summary POST /v1/users/servicelinks/events/approval/:targetUser
// @Tags servicelinks
// @Produce json
// @Param targetUser path string true "Target user identifier"
// @Param ownerUsername query string true "Owner username"
// @Param eventId query string true "Event ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /v1/users/servicelinks/events/approval/{targetUser} [post]
func postUsersServicelinksEventsApprovalTargetUserHandler(callBackRetryChan chan retryCallbacks, gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {

		identifier := strings.TrimSpace(strings.ToLower(c.Param("targetUser")))
		if identifier == "null" {
			log.Printf("user cannot be %v\n", identifier)
			c.JSON(http.StatusBadRequest, gin.H{"error": "user cannot be null"})
			return
		}
		ownerUsername := strings.TrimSpace(strings.ToLower(c.Query("ownerUsername")))
		eventID := strings.TrimSpace(c.Query("eventId"))
		//get authorization data for user
		eventData, err := servicelinkServices.GetEventAuthorizationData(ownerUsername, eventID, gc.DB)
		if err != nil {

			//other system error
			var ex tErrors.GenericError
			var ok bool

			ex, ok = err.(tErrors.GenericError)
			var statusCode int = 0
			var response interface{}

			if ok {
				statusCode = ex.HTTPCode()
				response = ex.JSONError()
			} else {
				statusCode = http.StatusBadRequest
				response = gin.H{"error": err.Error()}
			}
			log.Println("Event authorization error:", response)
			c.JSON(statusCode, response)
			return

		}
		conDB.PrintDBStats(fmt.Sprintf("POST /v1/users/servicelinks/events/approval/%v %v/%v", identifier, ownerUsername, eventID), gc.DB)

		mInfo, err := servicelinkServices.GetServiceLinkByAPIKey(eventData.ApiKey, gc.DB)

		if err != nil {
			log.Println("[GET SERVICE] error for SERVICELINK:", ownerUsername, "error: ", err)

			var ex tErrors.GenericError
			var ok bool

			ex, ok = err.(tErrors.GenericError)
			var statusCode int = 0
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

		if mInfo.EventPermission == 0 {
			//wrong access

			statusCode := http.StatusUnauthorized
			response := gin.H{"error": "error-invalid-service-access", "data": "Permission", "message": "Event permission not enabled for this service."}
			log.Printf("%+v\n", response)
			c.JSON(statusCode, response)
			return
		}

		userInfo, err := servicelinkServices.GetUserFromPrimarySigner(middleware.ExtractSigner(c), gc.DB, gc)

		if err != nil {
			log.Println("[GET UserInfo] error for user:", identifier, "error: ", err)

			var ex tErrors.GenericError
			var ok bool

			ex, ok = err.(tErrors.GenericError)
			var statusCode int = 0
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

		if userInfo.MobileVerified == 0 && os.Getenv("ENABLE_MOBILE_VERIFICATION") == "1" {

			statusCode := http.StatusBadRequest
			response := gin.H{"error": "error-invalid-user-access", "data": "Authentication", "message": "Your mobile phone number is not verified. Please make sure you have the latest version of the app and then go to settings/profile and click to verify your phone number before you can continue with this request."}

			c.JSON(statusCode, response)
			return
		}

		if userInfo.PushNotificationToken != nil {
			dataPayload := make(map[string]string)
			dataPayload["route"] = ""
			pns.SendFirebaseMessage(*userInfo.PushNotificationToken, fmt.Sprintf("Event Registration/Participation for %v authorized!", userInfo.Username), fmt.Sprintf("Your TrovoApp username %v has been used to authorize an event registration/participation action on %v service.", userInfo.Username, mInfo.LongName), "", dataPayload, gc.PushNotificationClient, gc.PNSContext)
		}

		//return report to user and not keep them waiting.
		c.JSON(http.StatusOK, gin.H{"message": "success"})

		if eventData.CallbackURL != nil {

			//make callback request
			// callbackResponse := new(map[string]interface{})
			type payload struct {
				EventID    string `json:"eventId"`
				TargetUser string `json:"targetUser"`
				DeviceID   string `json:"deviceId"`
			}
			jsonPayload := payload{EventID: eventData.ID, TargetUser: userInfo.Username, DeviceID: c.GetHeader("X-TW-DEVICE-ID")}
			/////
			d := *eventData.CallbackURL
			body, err := json.Marshal(jsonPayload)
			if err != nil {
				log.Printf("[EVENT CALLBACK] could not unmarshal callback message due to [%v]\n", err)

			}
			log.Printf("[EVENT CALLBACK] JSON STRING: [%v]\n", string(body))

			responseBody := bytes.NewBuffer(body)
			//Leverage Go's HTTP Post function to make request
			resp, err := http.Post(d, "application/json", responseBody)
			//Handle Error
			if err != nil {
				log.Printf("[EVENT CALLBACK] could not send callback message due to [%v]\n", err)
				return
			}
			defer resp.Body.Close()
			//Read the response body
			body, err = io.ReadAll(resp.Body)
			if err != nil {
				//send to retry channel
				log.Println("[EVENT CALLBACK AuthCallback] callback failed:", err)
				c := retryCallbacks{Req: responseBody, CallbackURL: d, Count: 0}
				callBackRetryChan <- c
			} else {
				log.Printf("[EVENT CALLBACK] Event Callback successful to: [%v], Response:[%v]\n\n", d, string(body))

			}

		}

	}
}

// getServicelinksAuthorizeVerifyOwnerUsernameTargetUserAuthIdHandler godoc
// @Summary GET /v1/servicelinks/authorize/verify/:ownerUsername/:targetUser/:authId
// @Tags servicelinks
// @Produce json
// @Param ownerUsername path string true "Owner username"
// @Param targetUser path string true "Target user"
// @Param authId path string true "Authorization ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /v1/servicelinks/authorize/verify/{ownerUsername}/{targetUser}/{authId} [get]
func getServicelinksAuthorizeVerifyOwnerUsernameTargetUserAuthIdHandler(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// var err error

		ownerUsername := strings.TrimSpace(strings.ToLower(c.Param("ownerUsername")))
		trovoUser := strings.TrimSpace(strings.ToLower(c.Param("targetUser")))
		authID := strings.TrimSpace(strings.ToLower(c.Param("authId")))

		conDB.PrintDBStats(fmt.Sprintf("GET /v1/servicelinks/authorize/verify/%v/%v/%v", ownerUsername, trovoUser, authID), gc.DB)
		mInfo, err := servicelinkServices.GetServiceLinkByAPIKey(middleware.ExtractServiceLinkApiKey(c), gc.DB)

		if err != nil {
			log.Println("[GET SERVICE INFO] error for SERVICE:", ownerUsername, "error: ", err)

			var ex tErrors.GenericError
			var ok bool

			ex, ok = err.(tErrors.GenericError)
			var statusCode int = 0
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

		if mInfo.AuthorizationPermission == 0 {
			//wrong access
			statusCode := http.StatusUnauthorized
			response := gin.H{"error": "error-invalid-service-access", "data": "Permission", "message": "authorization permission not enabled for this service"}
			c.JSON(statusCode, response)
			return
		}
		if mInfo.OwnerUsername != ownerUsername {
			//wrong access
			statusCode := http.StatusUnauthorized
			response := gin.H{"error": "error-invalid-service-access", "data": "Authentication", "message": "Authentication failed"}
			c.JSON(statusCode, response)

			return
		}
		userInfo, err := servicelinkServices.GetUserForServiceLink(trovoUser, mInfo, gc.DB, gc)

		if err != nil {
			log.Println("[GET UserInfo] error for user:", trovoUser, "error: ", err)

			var ex tErrors.GenericError
			var ok bool

			ex, ok = err.(tErrors.GenericError)
			var statusCode int = 0
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

		authData, err := servicelinkServices.GetUserAuthorizationData(ownerUsername, userInfo.Username, authID, gc.DB)
		if err != nil {

			var ex tErrors.GenericError
			var ok bool

			ex, ok = err.(tErrors.GenericError)
			var statusCode int = 0
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
		if authData.Authorized == 0 {

			//authorixation not approved
			response := gin.H{"error": "error-unauthorized-request", "data": "unauthorizedRequest", "message": "authorization awaiting approval"}
			statusCode := http.StatusUnauthorized
			c.JSON(statusCode, response)
			return

		}
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	}
}

// getServicelinksAppAuthorizeVerifyOwnerUsernameTargetUserAuthIdHandler godoc
// @Summary GET /v1/servicelinks/app/authorize/verify/:ownerUsername/:targetUser/:authId
// @Tags servicelinks
// @Produce json
// @Param ownerUsername path string true "Owner username"
// @Param targetUser path string true "Target user"
// @Param authId path string true "Authorization ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /v1/servicelinks/app/authorize/verify/{ownerUsername}/{targetUser}/{authId} [get]
func getServicelinksAppAuthorizeVerifyOwnerUsernameTargetUserAuthIdHandler(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// var err error

		ownerUsername := strings.TrimSpace(strings.ToLower(c.Param("ownerUsername")))
		trovoUser := strings.TrimSpace(strings.ToLower(c.Param("targetUser")))
		authID := strings.TrimSpace(strings.ToLower(c.Param("authId")))

		conDB.PrintDBStats(fmt.Sprintf("GET /v1/servicelinks/authorize/verify/%v/%v/%v", ownerUsername, trovoUser, authID), gc.DB)
		mInfo, err := servicelinkServices.GetServiceLinkByAPIKey(middleware.ExtractServiceLinkApiKey(c), gc.DB)

		if err != nil {
			log.Println("[GET SERVICE INFO] error for SERVICE:", ownerUsername, "error: ", err)

			var ex tErrors.GenericError
			var ok bool

			ex, ok = err.(tErrors.GenericError)
			var statusCode int = 0
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
		if mInfo.PublicKey != middleware.ExtractPublicKey(c) {
			//wrong access
			statusCode := http.StatusUnauthorized
			response := gin.H{"error": "error-invalid-service-access", "data": "Authentication", "message": "Authentication failed"}
			c.JSON(statusCode, response)
			return
		}
		// log.Printf("service Info: %+v\n", mInfo)
		if mInfo.AuthorizationPermission == 0 {
			//wrong access
			statusCode := http.StatusUnauthorized
			response := gin.H{"error": "error-invalid-service-access", "data": "Permission", "message": "authorization permission not enabled for this service"}
			c.JSON(statusCode, response)
			return
		}

		userInfo, err := servicelinkServices.GetUserForServiceLink(trovoUser, mInfo, gc.DB, gc)

		if err != nil {
			log.Println("[GET UserInfo] error for user:", trovoUser, "error: ", err)

			var ex tErrors.GenericError
			var ok bool

			ex, ok = err.(tErrors.GenericError)
			var statusCode int = 0
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

		authData, err := servicelinkServices.GetUserAuthorizationData(ownerUsername, userInfo.Username, authID, gc.DB)
		if err != nil {

			var ex tErrors.GenericError
			var ok bool

			ex, ok = err.(tErrors.GenericError)
			var statusCode int = 0
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
		if authData.Authorized == 0 {

			//authorixation not approved
			response := gin.H{"error": "error-unauthorized-request", "data": "unauthorizedRequest", "message": "authorization awaiting approval"}
			statusCode := http.StatusUnauthorized
			c.JSON(statusCode, response)
			return

		}
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	}
}

// getServicelinksPaymentRequestTargetUserHandler godoc
// @Summary GET /v1/servicelinks/payment/request/:targetUser
// @Tags servicelinks
// @Produce json
// @Param targetUser path string true "Target user identifier"
// @Param ownerUsername query string false "Owner username"
// @Param paymentDestination query string false "Payment destination"
// @Param assetCode query string false "Asset code"
// @Param assetIssuer query string false "Asset issuer"
// @Param amount query string false "Amount"
// @Param memo query string false "Memo"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /v1/servicelinks/payment/request/{targetUser} [get]
func getServicelinksPaymentRequestTargetUserHandler(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {

		ownerUsername := strings.TrimSpace(strings.ToLower(c.Query("ownerUsername")))
		trovoUser := strings.TrimSpace(strings.ToLower(c.Param("targetUser")))
		paymentDestination := strings.TrimSpace(strings.ToLower(c.Query("paymentDestination")))
		if len(paymentDestination) == 56 {
			paymentDestination = strings.ToUpper(paymentDestination)
		}
		assetCode := strings.TrimSpace(strings.ToUpper(c.Query("assetCode")))
		assetIssuer := strings.TrimSpace(strings.ToUpper(c.Query("assetIssuer")))
		amount := strings.TrimSpace(c.Query("amount"))
		memo := strings.TrimSpace(c.Query("memo"))

		cacheKey := fmt.Sprintf("[GET] /v1/servicelinks/payment/request/%v %v", trovoUser, ownerUsername)
		cacheKeyParameters := fmt.Sprintf("%v", c.Request.URL.RawQuery)

		{
			//search cache

			ok, status, response := gc.RedisCache.CachedHttpResponseWithParameters(cacheKey, cacheKeyParameters)

			if ok {
				// log.Printf("[%v]/[%v], served from cache\n", cacheKey, cacheKeyParameters)
				c.JSON(status, response)
				return
			}
		}

		conDB.PrintDBStats(fmt.Sprintf("GET /v1/servicelinks/payment/%v/%v/?paymentDestination=%v&assetCode=%v&assetIssuer=%v&amount=%v&memo=%v", trovoUser, ownerUsername, paymentDestination, assetCode, assetIssuer, amount, memo), gc.DB)
		mInfo, err := servicelinkServices.GetServiceLinkByAPIKey(middleware.ExtractServiceLinkApiKey(c), gc.DB)

		if err != nil {
			log.Println("[GET SERVICE PAYMENT DATA] error for SERVICE:", ownerUsername, "error: ", err)

			var ex tErrors.GenericError
			var ok bool

			ex, ok = err.(tErrors.GenericError)
			var statusCode int = 0
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
		// if mInfo.PublicKey != middleware.ExtractPublicKey(c) {
		// 	//wrong access
		// 	statusCode := http.StatusUnauthorized
		// 	response := gin.H{"error": "error-invalid-service-access", "data": "Authentication", "message": "Authentication failed"}
		// 	c.JSON(statusCode, response)
		// 	return
		// }
		// log.Printf("service Infor: %+v\n", mInfo)
		if mInfo.PaymentPermission == 0 {
			//wrong access
			statusCode := http.StatusUnauthorized
			response := gin.H{"error": "error-invalid-service-access", "data": "Permission", "message": "payment permission not enabled for this service"}
			c.JSON(statusCode, response)
			return
		}

		_, err = servicelinkServices.GetUserForServiceLink(trovoUser, mInfo, gc.DB, gc)

		if err != nil {
			log.Println("[GET UserInfo] error for user:", trovoUser, "error: ", err)

			var ex tErrors.GenericError
			var ok bool

			ex, ok = err.(tErrors.GenericError)
			var statusCode int = 0
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

		//generate payment data
		data, err := dl.GeneratePaymentData(paymentDestination, assetCode, assetIssuer, amount, memo, gc)
		if err != nil {
			//could not create login session
			response := gin.H{"error": "error-temporary-server-error", "data": "temporaryServerError", "message": "Temporary Server Error. Contact support."}
			statusCode := http.StatusServiceUnavailable
			c.JSON(statusCode, response)
			return
		}
		{
			cacheDurationInSeconds := 20 * 60 //2 minutes

			gc.RedisCache.CacheHttpResponseWithParameters(cacheKey, cacheKeyParameters, http.StatusOK, data, cacheDurationInSeconds)
		}
		c.JSON(http.StatusOK, data)

	}
}

// getTrovoApiPaymentRequestTargetUserHandler godoc
// @Summary GET /v1/trovo-api/payment/request/:targetUser
// @Tags servicelinks
// @Produce json
// @Param targetUser path string true "Target user identifier"
// @Param ownerUsername query string false "Owner username"
// @Param paymentDestination query string false "Payment destination"
// @Param assetCode query string false "Asset code"
// @Param assetIssuer query string false "Asset issuer"
// @Param amount query string false "Amount"
// @Param memo query string false "Memo"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /v1/trovo-api/payment/request/{targetUser} [get]
func getTrovoApiPaymentRequestTargetUserHandler(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {

		ownerUsername := strings.TrimSpace(strings.ToLower(c.Query("ownerUsername")))
		trovoUser := strings.TrimSpace(strings.ToLower(c.Param("targetUser")))
		paymentDestination := strings.TrimSpace(strings.ToLower(c.Query("paymentDestination")))
		if len(paymentDestination) == 56 {
			paymentDestination = strings.ToUpper(paymentDestination)
		}
		assetCode := strings.TrimSpace(strings.ToUpper(c.Query("assetCode")))
		assetIssuer := strings.TrimSpace(strings.ToUpper(c.Query("assetIssuer")))
		amount := strings.TrimSpace(c.Query("amount"))
		memo := strings.TrimSpace(c.Query("memo"))

		cacheKey := fmt.Sprintf("[GET] /v1/servicelinks/payment/request/%v %v", trovoUser, ownerUsername)
		cacheKeyParameters := fmt.Sprintf("%v", c.Request.URL.RawQuery)

		{
			//search cache

			ok, status, response := gc.RedisCache.CachedHttpResponseWithParameters(cacheKey, cacheKeyParameters)

			if ok {
				// log.Printf("[%v]/[%v], served from cache\n", cacheKey, cacheKeyParameters)
				c.JSON(status, response)
				return
			}
		}

		conDB.PrintDBStats(fmt.Sprintf("GET /v1/servicelinks/payment/%v/%v/?paymentDestination=%v&assetCode=%v&assetIssuer=%v&amount=%v&memo=%v", trovoUser, ownerUsername, paymentDestination, assetCode, assetIssuer, amount, memo), gc.DB)
		mInfo, err := servicelinkServices.GetServiceLinkByAPIKey(middleware.ExtractServiceLinkApiKey(c), gc.DB)

		if err != nil {
			log.Println("[GET SERVICE PAYMENT DATA] error for SERVICE:", ownerUsername, "error: ", err)

			var ex tErrors.GenericError
			var ok bool

			ex, ok = err.(tErrors.GenericError)
			var statusCode int = 0
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
		// if mInfo.PublicKey != middleware.ExtractPublicKey(c) {
		// 	//wrong access
		// 	statusCode := http.StatusUnauthorized
		// 	response := gin.H{"error": "error-invalid-service-access", "data": "Authentication", "message": "Authentication failed"}
		// 	c.JSON(statusCode, response)
		// 	return
		// }
		// log.Printf("service Infor: %+v\n", mInfo)
		if mInfo.PaymentPermission == 0 {
			//wrong access
			statusCode := http.StatusUnauthorized
			response := gin.H{"error": "error-invalid-service-access", "data": "Permission", "message": "payment permission not enabled for this service"}
			c.JSON(statusCode, response)
			return
		}

		_, err = servicelinkServices.GetUserForServiceLink(trovoUser, mInfo, gc.DB, gc)

		if err != nil {
			log.Println("[GET UserInfo] error for user:", trovoUser, "error: ", err)

			var ex tErrors.GenericError
			var ok bool

			ex, ok = err.(tErrors.GenericError)
			var statusCode int = 0
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

		//generate payment data
		data, err := dl.GeneratePaymentData(paymentDestination, assetCode, assetIssuer, amount, memo, gc)
		if err != nil {
			//could not create login session
			response := gin.H{"error": "error-temporary-server-error", "data": "temporaryServerError", "message": "Temporary Server Error. Contact support."}
			statusCode := http.StatusServiceUnavailable
			c.JSON(statusCode, response)
			return
		}
		{
			cacheDurationInSeconds := 20 * 60 //2 minutes

			gc.RedisCache.CacheHttpResponseWithParameters(cacheKey, cacheKeyParameters, http.StatusOK, data, cacheDurationInSeconds)
		}
		c.JSON(http.StatusOK, data)

	}
}

// getServicelinksTokenizedAssetAssetCodeHandler godoc
// @Summary GET /v1/servicelinks/tokenized-asset/:assetCode
// @Tags servicelinks
// @Produce json
// @Param assetCode path string true "Asset code"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /v1/servicelinks/tokenized-asset/{assetCode} [get]
func getServicelinksTokenizedAssetAssetCodeHandler(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {

		assetCode := strings.TrimSpace(strings.ToLower(c.Param("assetCode")))

		cacheKey := fmt.Sprintf("[GET] /v1/servicelinks/tokenized-asset/:assetCode/%v", assetCode)

		{
			//search cache

			ok, status, response := gc.RedisCache.CachedHttpResponse(cacheKey)

			if ok {
				// log.Printf("[%v], served from cache\n", cacheKey)
				c.JSON(status, response)
				return
			}
		}

		mInfo, err := servicelinkServices.GetServiceLinkByAPIKey(middleware.ExtractServiceLinkApiKey(c), gc.DB)

		if err != nil {
			log.Println("[GET SERVICE TOKENIZED ASSET DATA] error for SERVICE error: ", err)

			var ex tErrors.GenericError
			var ok bool

			ex, ok = err.(tErrors.GenericError)
			var statusCode int = 0
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

		if mInfo.TokenInfoPermission == 0 {
			//wrong access
			statusCode := http.StatusUnauthorized
			response := gin.H{"error": "error-invalid-service-access", "data": "Permission", "message": "Token information permission not enabled for this service"}
			c.JSON(statusCode, response)
			return
		}

		gcta := gc.GetTokenizedAssetByCode(assetCode)
		rgcta := &gcta
		data := rgcta.ToJSON()
		{
			cacheDurationInSeconds := 20 * 60 //2 minutes

			gc.RedisCache.CacheHttpResponse(cacheKey, http.StatusOK, data, cacheDurationInSeconds)
		}
		c.JSON(http.StatusOK, data)

	}
}

// getServicelinksOwnerUsernameTargetUserUserinfoHandler godoc
// @Summary GET /v1/servicelinks/:ownerUsername/:targetUser/userinfo
// @Tags servicelinks
// @Produce json
// @Param ownerUsername path string true "Owner username"
// @Param targetUser path string true "Target user"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /v1/servicelinks/{ownerUsername}/{targetUser}/userinfo [get]
func getServicelinksOwnerUsernameTargetUserUserinfoHandler(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {

		identifier := strings.TrimSpace(strings.ToLower(c.Param("ownerUsername")))
		trovoUser := strings.TrimSpace(strings.ToLower(c.Param("targetUser")))
		{
			//do not cache this so that it brings the latest data
			// cacheKey := fmt.Sprintf("[GET] /v1/servicelinks/%v/%v/payment", identifier, trovoUser)
			// cacheKeyParameters := fmt.Sprintf("%v", c.Request.URL.RawQuery)

			// {
			// 	//search cache

			// 	ok, status, response := redisCache.CachedHttpResponseWithParameters(cacheKey, cacheKeyParameters)

			// 	if ok {
			// 		log.Printf("[%v]/[%v], served from cache\n", cacheKey, cacheKeyParameters)
			// 		c.JSON(status, response)
			// 		return
			// 	}
			// }
		}
		conDB.PrintDBStats(fmt.Sprintf("GET /v1/servicelinks/%v/%v/userinfo", identifier, trovoUser), gc.DB)
		mInfo, err := servicelinkServices.GetServiceLinkByAPIKey(middleware.ExtractServiceLinkApiKey(c), gc.DB)

		if err != nil {
			log.Println("[GET USER DATA] error for SERVICE:", identifier, "error: ", err)

			var ex tErrors.GenericError
			var ok bool

			ex, ok = err.(tErrors.GenericError)
			var statusCode int = 0
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
		// if mInfo.PublicKey != middleware.ExtractPublicKey(c) {
		// 	//wrong access
		// 	statusCode := http.StatusUnauthorized
		// 	response := gin.H{"error": "error-invalid-service-access", "data": "Authentication", "message": "Authentication failed"}
		// 	c.JSON(statusCode, response)
		// 	return
		// }
		// log.Printf("service Info: %+v\n", mInfo)
		if mInfo.AllowUserInfo == 0 {
			//wrong access
			statusCode := http.StatusUnauthorized
			response := gin.H{"error": "error-invalid-service-access", "data": "Permission", "message": "userInfo permission not enabled for this service"}
			c.JSON(statusCode, response)
			return
		}

		data, err := servicelinkServices.GetUserForServiceLink(trovoUser, mInfo, gc.DB, gc)

		if err != nil {
			log.Println("[GET UserInfo] error for user:", trovoUser, "error: ", err)

			var ex tErrors.GenericError
			var ok bool

			ex, ok = err.(tErrors.GenericError)
			var statusCode int = 0
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

		c.JSON(http.StatusOK, data)

	}
}

// postServicelinksOwnerUsernameTargetUserPushHandler godoc
// @Summary POST /v1/servicelinks/:ownerUsername/:targetUser/push
// @Tags servicelinks
// @Accept json
// @Produce json
// @Param ownerUsername path string true "Owner username"
// @Param targetUser path string true "Target user"
// @Param body body servicelinkModels.ServiceLinkPushNotificationInput true "Push notification payload"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /v1/servicelinks/{ownerUsername}/{targetUser}/push [post]
func postServicelinksOwnerUsernameTargetUserPushHandler(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {

		ownerUsername := strings.TrimSpace(strings.ToLower(c.Param("ownerUsername")))
		trovoUser := strings.TrimSpace(strings.ToLower(c.Param("targetUser")))

		conDB.PrintDBStats(fmt.Sprintf("POST /v1/servicelinks/%v/%v/push?", ownerUsername, trovoUser), gc.DB)
		mInfo, err := servicelinkServices.GetServiceLinkByAPIKey(middleware.ExtractServiceLinkApiKey(c), gc.DB)

		if err != nil {
			log.Println("[GET service FOR PUSH] error for service:", ownerUsername, "error: ", err)

			var ex tErrors.GenericError
			var ok bool

			ex, ok = err.(tErrors.GenericError)
			var statusCode int = 0
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
		// if mInfo.PublicKey != middleware.ExtractPublicKey(c) {
		// 	//wrong access
		// 	statusCode := http.StatusUnauthorized
		// 	response := gin.H{"error": "error-invalid-service-access", "data": "Authentication", "message": "Authentication failed"}
		// 	c.JSON(statusCode, response)
		// 	return
		// }
		// log.Printf("service Infor: %+v\n", mInfo)
		if mInfo.PushNotificationPermission == 0 {
			//wrong access
			statusCode := http.StatusUnauthorized
			response := gin.H{"error": "error-invalid-service-access", "data": "Permission", "message": "Push notification permission not enabled for this service"}
			c.JSON(statusCode, response)
			return
		}

		userInfo, err := servicelinkServices.GetUserForServiceLink(trovoUser, mInfo, gc.DB, gc)

		if err != nil {
			log.Println("[GET UserInfo] error for user:", trovoUser, "error: ", err)

			var ex tErrors.GenericError
			var ok bool

			ex, ok = err.(tErrors.GenericError)
			var statusCode int = 0
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
		successResponseData := struct {
			Data string `json:"data"`
		}{
			Data: "OK",
		}
		//check if user has Push Notification Token set
		if userInfo.PushNotificationToken == nil {
			//no valid token set. user cannot receive push notification

			c.JSON(http.StatusOK, successResponseData)
			return
		}

		var serviceLinkRequestInput servicelinkModels.ServiceLinkPushNotificationInput
		reqBody, _ := io.ReadAll(c.Request.Body)

		err = json.Unmarshal(reqBody, &serviceLinkRequestInput)

		var invalidJSON tErrors.ErrorInvalidJSON

		if err != nil {
			log.Println("Push Notification Request Input JSON Error:", err)
			c.JSON(http.StatusBadRequest, invalidJSON.JSONError())
			return
		}
		//check if message is set
		if len(serviceLinkRequestInput.Message) == 0 {

			c.JSON(http.StatusOK, successResponseData)
			return
		}

		if len(serviceLinkRequestInput.Message) > 1000 {
			messageBytes := []byte(serviceLinkRequestInput.Message)
			msgB := messageBytes[0:999]
			serviceLinkRequestInput.Message = string(msgB)
		}

		//Push Message
		dataPayload := make(map[string]string)
		if serviceLinkRequestInput.Action == "login" {
			dataPayload["route"] = "login"
		} else if serviceLinkRequestInput.Action == "payment" {
			dataPayload["route"] = "payment"
		} else if serviceLinkRequestInput.Action == "auth" {
			dataPayload["route"] = "auth"
		} else if serviceLinkRequestInput.Action == "event" {
			dataPayload["route"] = "event"
		} else {
			dataPayload["route"] = "none"
		}

		userInfo.SendPushMessage(serviceLinkRequestInput.Title, serviceLinkRequestInput.Message, serviceLinkRequestInput.ImageURI, dataPayload, gc)

		c.JSON(http.StatusOK, successResponseData)
	}
}

// postTrovoApiUsersOnboardHandler godoc
// @Summary POST /v1/trovo-api/users/onboard
// @Tags servicelinks
// @Accept json
// @Produce json
// @Param body body userModels.UserRegistrationInfo true "User registration payload"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /v1/trovo-api/users/onboard [post]
func postTrovoApiUsersOnboardHandler(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {

		mInfo, err := servicelinkServices.GetServiceLinkByAPIKey(middleware.ExtractServiceLinkApiKey(c), gc.DB)

		if err != nil {
			log.Println("[GET service] error for service:", middleware.ExtractServiceLinkApiKey(c), "error: ", err)

			var ex tErrors.GenericError
			var ok bool

			ex, ok = err.(tErrors.GenericError)
			var statusCode int = 0
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

		// ownerUsername := mInfo.OwnerUsername

		if mInfo.CreateUsersPermission == 0 {
			//wrong access
			statusCode := http.StatusUnauthorized
			response := gin.H{"error": "error-invalid-service-access", "data": "Permission", "message": "Permission to create users not enabled for this service"}
			c.JSON(statusCode, response)
			return
		}

		var userRegistrationInfo userModels.UserRegistrationInfo
		// var err error

		data, _ := io.ReadAll(c.Request.Body)

		err = json.Unmarshal(data, &userRegistrationInfo)

		if len(userRegistrationInfo.PublicKey) != 56 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "error-invalid-public-key", "data": "publicKey", "message": "Invalid Public key."})
			return
		}
		if len(userRegistrationInfo.PrimarySigner) != 56 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "error-invalid-primary-signer-public-key", "data": "primarySigner", "message": "Invalid Primary Signer Public key."})
			return
		}

		userRegistrationInfo.CreatedByServiceLinkID = mInfo.ID
		userRegistrationInfo.PublicIP = c.ClientIP()
		if len(c.GetHeader("Cf-Connecting-Ip")) > 4 {
			userRegistrationInfo.PublicIP = c.GetHeader("Cf-Connecting-Ip")
		}

		var invalidJSON tErrors.ErrorInvalidJSON

		if err != nil {
			c.JSON(http.StatusBadRequest, invalidJSON.JSONError())
			return
		}

		//Generate username for the serviceLink user
		uid := uuid.NewString()
		endBatch := strings.Split(uid, "-")[len(strings.Split(uid, "-"))-1]
		// append it to the mInfo.Username
		ownerIdentifier := mInfo.OwnerUsername
		if len(ownerIdentifier) > 6 {
			ownerIdentifier = ownerIdentifier[0:5]
		}
		userRegistrationInfo.Username = ownerIdentifier + endBatch[len(endBatch)-6:] //append the last 6 characters of the endBatch
		//if referral is allowed, set the referrer
		if mInfo.AllowReferralForRegisteredUsers == 1 {
			userRegistrationInfo.Referrer = mInfo.OwnerUsername
		}

		//replace _ and /
		userRegistrationInfo.Username = strings.ReplaceAll(strings.ReplaceAll(strings.ToLower(userRegistrationInfo.Username), "_", ""), "/", "")

		_, _, err = userServices.RegisterUser(userRegistrationInfo, gc)

		if err != nil {
			var ex tErrors.GenericError
			var ok bool

			ex, ok = err.(tErrors.GenericError)
			if ok {
				c.JSON(http.StatusBadRequest, ex.JSONError())
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": err.Error()})
			}
			return
		}

		//At this point, there was no error.
		c.JSON(http.StatusOK, userRegistrationInfo)
	}
}

// postTrovoApiUsersUpdateKycHandler godoc
// @Summary POST /v1/trovo-api/users/update-kyc
// @Tags servicelinks
// @Accept json
// @Produce json
// @Param body body servicelinkModels.ServiceLinkUpdateKycInput true "KYC update payload"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /v1/trovo-api/users/update-kyc [post]
func postTrovoApiUsersUpdateKycHandler(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {

		mInfo, err := servicelinkServices.GetServiceLinkByAPIKey(middleware.ExtractServiceLinkApiKey(c), gc.DB)

		if err != nil {
			log.Println("[GET service] error for service:", middleware.ExtractServiceLinkApiKey(c), "error: ", err)

			var ex tErrors.GenericError
			var ok bool

			ex, ok = err.(tErrors.GenericError)
			var statusCode int = 0
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

		// ownerUsername := mInfo.OwnerUsername

		if mInfo.CreateUsersPermission == 0 {
			//wrong access
			statusCode := http.StatusUnauthorized
			response := gin.H{"error": "error-invalid-service-access", "data": "Permission", "message": "Permission to create users not enabled for this service"}
			c.JSON(statusCode, response)
			return
		}

		var kycData servicelinkModels.ServiceLinkUpdateKycInput
		// var err error

		data, _ := io.ReadAll(c.Request.Body)

		err = json.Unmarshal(data, &kycData)

		var invalidJSON tErrors.ErrorInvalidJSON

		if err != nil {
			c.JSON(http.StatusBadRequest, invalidJSON.JSONError())
			return
		}

		//validate the json data
		sliceOfByte := []byte(kycData.KycJsonData)

		if len(sliceOfByte) < 10 {
			statusCode := http.StatusBadRequest
			response := gin.H{"error": "error-invalid-kyc-data", "data": "KycJsonData", "message": "KYC data is invalid."}
			c.JSON(statusCode, response)
			return

		}

		if strings.Contains(kycData.TargetTrovoUsername, "%") {
			u, e := url.QueryUnescape(kycData.TargetTrovoUsername)
			if e == nil {
				kycData.TargetTrovoUsername = u
			}
		}
		if kycData.KycStatus != 1 && kycData.KycStatus != 2 && kycData.KycStatus != 3 && kycData.KycStatus != 4 {
			//wrong status
			statusCode := http.StatusBadRequest
			response := gin.H{"error": "error-invalid-kyc-status", "data": "KycStatus", "message": "KYC status is invalid."}
			c.JSON(statusCode, response)
			return
		}

		//get the user object
		targetUser, err := userModels.Username(kycData.KycStatus).GetSimpleUser(gc.DB, gc)
		if err != nil {
			log.Println("[GET servicelink target user] error for getting user:", kycData.TargetTrovoUsername, "error: ", err)

			var ex tErrors.GenericError
			var ok bool

			ex, ok = err.(tErrors.GenericError)
			var statusCode int = 0
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

		//check if the user belongs to the service link
		if targetUser.CreatedByServiceLinkID != nil {
			if mInfo.ID != *targetUser.CreatedByServiceLinkID {

				statusCode := http.StatusUnauthorized
				response := gin.H{"error": "error-invalid-user-access", "data": "KycStatus", "message": "User was not created by your service."}
				c.JSON(statusCode, response)
				return
			}
		}

		err = servicelinkServices.UpdateUserKYCStatus(&targetUser, kycData.KycStatus, kycData.KycJsonData, gc)

		if err != nil {
			var ex tErrors.GenericError
			var ok bool

			ex, ok = err.(tErrors.GenericError)
			if ok {
				c.JSON(http.StatusBadRequest, ex.JSONError())
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": err.Error()})
			}
			return
		}

		//At this point, there was no error.

		c.JSON(http.StatusOK, kycData)
	}
}

// postTrovoApiTokensMintHandler godoc
// @Summary POST /v1/trovo-api/tokens/mint
// @Tags servicelinks
// @Accept json
// @Produce json
// @Param body body userModels.MintingInfo true "Minting payload"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /v1/trovo-api/tokens/mint [post]
func postTrovoApiTokensMintHandler(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {

		mInfo, err := servicelinkServices.GetServiceLinkByAPIKey(middleware.ExtractServiceLinkApiKey(c), gc.DB)

		if err != nil {
			log.Println("[GET service] error for service:", middleware.ExtractServiceLinkApiKey(c), "error: ", err)

			var ex tErrors.GenericError
			var ok bool

			ex, ok = err.(tErrors.GenericError)
			var statusCode int = 0
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

		// ownerUsername := mInfo.OwnerUsername

		if mInfo.CreateUsersPermission == 0 {
			//wrong access
			statusCode := http.StatusUnauthorized
			response := gin.H{"error": "error-invalid-service-access", "data": "Permission", "message": "Permission to create users not enabled for this service"}
			c.JSON(statusCode, response)
			return
		}

		signerUser, err := usersDB.GetUser(mInfo.OwnerUsername, gc.DB, gc)

		if err != nil {
			var ex tErrors.GenericError
			var ok bool

			ex, ok = err.(tErrors.GenericError)
			if ok {
				c.JSON(http.StatusBadRequest, ex.JSONError())
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": err.Error()})
			}
			return
		}

		wallet, _, err := usersDB.GetWallet(middleware.ExtractPublicKey(c), gc.DB)

		if err != nil {
			var ex tErrors.GenericError
			var ok bool

			ex, ok = err.(tErrors.GenericError)
			if ok {
				c.JSON(http.StatusBadRequest, ex.JSONError())
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": err.Error()})
			}
			return
		}

		var mintingData userModels.MintingInfo
		// var err error

		data, _ := io.ReadAll(c.Request.Body)

		err = json.Unmarshal(data, &mintingData)

		var invalidJSON tErrors.ErrorInvalidJSON

		if err != nil {
			c.JSON(http.StatusBadRequest, invalidJSON.JSONError())
			return
		}

		if len(mintingData.Destination) == 0 {
			//wrong status
			statusCode := http.StatusBadRequest
			response := gin.H{"error": "error-invalid-minting-destination", "data": "destination", "message": "Invalid minting Destination."}
			c.JSON(statusCode, response)
			return
		}

		if len(mintingData.AssetCode) == 0 {
			//wrong status
			statusCode := http.StatusBadRequest
			response := gin.H{"error": "error-invalid-minting-asset-code", "data": "assetCode", "message": "Invalid Asset Code."}
			c.JSON(statusCode, response)
			return
		}

		if len(mintingData.AssetIssuer) == 0 {
			//wrong status
			statusCode := http.StatusBadRequest
			response := gin.H{"error": "error-invalid-minting-asset-issuer", "data": "assetIssuer", "message": "Invalid Asset Issuer."}
			c.JSON(statusCode, response)
			return
		}

		if len(mintingData.Amount) == 0 {
			//wrong status
			statusCode := http.StatusBadRequest
			response := gin.H{"error": "error-invalid-minting-amount", "data": "amount", "message": "Invalid Amount."}
			c.JSON(statusCode, response)
			return
		}

		_, _, err = userServices.MintAsset(&signerUser, &wallet, &mintingData, gc)

		if err != nil {
			var ex tErrors.GenericError
			var ok bool

			ex, ok = err.(tErrors.GenericError)
			if ok {
				c.JSON(http.StatusBadRequest, ex.JSONError())
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": err.Error()})
			}
			return
		}
		//At this point, there was no error.

		c.JSON(http.StatusOK, mintingData)
	}
}

// getTrovoApiUsersBalanceWalletPublicKeyHandler godoc
// @Summary GET /v1/trovo-api/users/balance/:walletPublicKey
// @Tags servicelinks
// @Produce json
// @Param walletPublicKey path string true "Wallet public key"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /v1/trovo-api/users/balance/{walletPublicKey} [get]
func getTrovoApiUsersBalanceWalletPublicKeyHandler(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {

		mInfo, err := servicelinkServices.GetServiceLinkByAPIKey(middleware.ExtractServiceLinkApiKey(c), gc.DB)

		if err != nil {
			log.Println("[GET service] error for service:", middleware.ExtractServiceLinkApiKey(c), "error: ", err)

			var ex tErrors.GenericError
			var ok bool

			ex, ok = err.(tErrors.GenericError)
			var statusCode int = 0
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

		// ownerUsername := mInfo.OwnerUsername

		if mInfo.CreateUsersPermission == 0 {
			//wrong access
			statusCode := http.StatusUnauthorized
			response := gin.H{"error": "error-invalid-service-access", "data": "Permission", "message": "Permission to create users not enabled for this service"}
			c.JSON(statusCode, response)
			return
		}

		signerUser, err := usersDB.GetUser(mInfo.OwnerUsername, gc.DB, gc)

		if err != nil {
			var ex tErrors.GenericError
			var ok bool

			ex, ok = err.(tErrors.GenericError)
			if ok {
				c.JSON(http.StatusBadRequest, ex.JSONError())
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": err.Error()})
			}
			return
		}
		walletPublicKey := c.Param("walletPublicKey")

		wallet, _, err := usersDB.GetWallet(walletPublicKey, gc.DB)

		if err != nil {
			var ex tErrors.GenericError
			var ok bool

			ex, ok = err.(tErrors.GenericError)
			if ok {
				c.JSON(http.StatusBadRequest, ex.JSONError())
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": err.Error()})
			}
			return
		}
		// check wallet owner to be sure of the oriviledge to access it.
		walletOwner, err := wallet.GetWalletOwner(gc.DB, gc)

		if err != nil {
			var ex tErrors.GenericError
			var ok bool

			ex, ok = err.(tErrors.GenericError)
			if ok {
				c.JSON(http.StatusBadRequest, ex.JSONError())
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": err.Error()})
			}
			return
		}
		if (walletOwner.Username != signerUser.Username) && walletOwner.CreatedByServiceLinkID == nil {
			//attempt to check ballance of wallet outside of your control.
			statusCode := http.StatusUnauthorized
			response := gin.H{"error": "error-invalid-service-access", "data": "Permission", "message": "Cannot access balance of wallets outside of your scope."}
			c.JSON(statusCode, response)
			return
		}
		if (walletOwner.Username != signerUser.Username) && walletOwner.CreatedByServiceLinkID != nil {
			// check if service link validates
			if *walletOwner.CreatedByServiceLinkID != mInfo.ID {
				//attempt to check ballance of wallet outside of your control.
				statusCode := http.StatusUnauthorized
				response := gin.H{"error": "error-invalid-service-access", "data": "Permission", "message": "Cannot access balance of wallets outside of your scope."}
				c.JSON(statusCode, response)
				return
			}

		}

		assetBalances, err := wallet.GetWalletAssetBalances(gc)
		if err != nil {
			log.Printf("[GET Wallet Balances] error for signer:%v, publicKey: %v, error: %v", signerUser.Username, walletPublicKey, err)

			var ex tErrors.GenericError
			var ok bool

			ex, ok = err.(tErrors.GenericError)
			var statusCode int = 0
			var response interface{}

			if ok {
				statusCode = ex.HTTPCode()
				response = ex.JSONError()
			} else {
				statusCode = http.StatusBadRequest
				response = gin.H{"error": err.Error(), "message": err.Error()}
			}

			c.JSON(statusCode, response)
			return
		}

		c.JSON(http.StatusOK, assetBalances)
	}
}

// getTrovoApiUsersPaymentHistoryWalletPublicKeyHandler godoc
// @Summary GET /v1/trovo-api/users/payment-history/:walletPublicKey
// @Tags servicelinks
// @Produce json
// @Param walletPublicKey path string true "Wallet public key"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /v1/trovo-api/users/payment-history/{walletPublicKey} [get]
func getTrovoApiUsersPaymentHistoryWalletPublicKeyHandler(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {

		mInfo, err := servicelinkServices.GetServiceLinkByAPIKey(middleware.ExtractServiceLinkApiKey(c), gc.DB)

		if err != nil {
			log.Println("[GET service] error for service:", middleware.ExtractServiceLinkApiKey(c), "error: ", err)

			var ex tErrors.GenericError
			var ok bool

			ex, ok = err.(tErrors.GenericError)
			var statusCode int = 0
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

		// ownerUsername := mInfo.OwnerUsername

		if mInfo.CreateUsersPermission == 0 {
			//wrong access
			statusCode := http.StatusUnauthorized
			response := gin.H{"error": "error-invalid-service-access", "data": "Permission", "message": "Permission to create users not enabled for this service"}
			c.JSON(statusCode, response)
			return
		}

		signerUser, err := usersDB.GetUser(mInfo.OwnerUsername, gc.DB, gc)

		if err != nil {
			var ex tErrors.GenericError
			var ok bool

			ex, ok = err.(tErrors.GenericError)
			if ok {
				c.JSON(http.StatusBadRequest, ex.JSONError())
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": err.Error()})
			}
			return
		}
		walletPublicKey := c.Param("walletPublicKey")

		wallet, temp, err := usersDB.GetWallet(walletPublicKey, gc.DB)

		if err != nil {
			var ex tErrors.GenericError
			var ok bool

			ex, ok = err.(tErrors.GenericError)
			if ok {
				c.JSON(http.StatusBadRequest, ex.JSONError())
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": err.Error()})
			}
			return
		}

		if temp {

			statusCode := http.StatusBadRequest
			response := gin.H{"error": "error only main wallets allowed", "message": "Only main wallets are allowed. The address you provided is not a main wallet."}

			c.JSON(statusCode, response)
			return
		}
		// check wallet owner to be sure of the oriviledge to access it.
		walletOwner, err := wallet.GetWalletOwner(gc.DB, gc)

		if err != nil {
			var ex tErrors.GenericError
			var ok bool

			ex, ok = err.(tErrors.GenericError)
			if ok {
				c.JSON(http.StatusBadRequest, ex.JSONError())
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": err.Error()})
			}
			return
		}
		if (walletOwner.Username != signerUser.Username) && walletOwner.CreatedByServiceLinkID == nil {
			//attempt to check ballance of wallet outside of your control.
			statusCode := http.StatusUnauthorized
			response := gin.H{"error": "error-invalid-service-access", "data": "Permission", "message": "Cannot access history of wallets outside of your scope."}
			c.JSON(statusCode, response)
			return
		}
		if (walletOwner.Username != signerUser.Username) && walletOwner.CreatedByServiceLinkID != nil {
			// check if service link validates
			if *walletOwner.CreatedByServiceLinkID != mInfo.ID {
				//attempt to check ballance of wallet outside of your control.
				statusCode := http.StatusUnauthorized
				response := gin.H{"error": "error-invalid-service-access", "data": "Permission", "message": "Cannot access history of wallets outside of your scope."}
				c.JSON(statusCode, response)
				return
			}

		}

		_, err = keypair.ParseAddress(walletPublicKey)
		if err != nil {

			statusCode := http.StatusBadRequest
			response := gin.H{"error": "error invalid address", "message": "Only valid addresses are allowed"}

			c.JSON(statusCode, response)
			return
		}
		cacheKey := fmt.Sprintf("[GET] /v1/users/payments/%v", walletPublicKey)
		cacheKeyParameters := c.Request.URL.RequestURI()
		{
			// check cache
			ok, status, response := gc.RedisCache.CachedHttpResponseWithParameters(cacheKey, cacheKeyParameters)

			if ok {
				// log.Printf("[%v]/[%v], served from cache\n", cacheKey, cacheKeyParameters)
				c.JSON(status, response)
				return
			}

		}
		// cacheDurationInSeconds := 1 * 60 //1 minutes
		cacheDurationInSeconds := 20 //in seconds

		//Get Payment history
		historyRecords := paymentServices.GetPaymentHistory(walletPublicKey, gc, c)

		c.JSON(http.StatusOK, historyRecords)
		// gc.RedisCache.CacheHttpResponse(cacheKey, http.StatusOK, historyRecords, cacheDurationInSeconds)
		gc.RedisCache.CacheHttpResponseWithParameters(cacheKey, cacheKeyParameters, http.StatusOK, historyRecords, cacheDurationInSeconds)

	}
}

// postTrovoApiUsersPaymentHandler godoc
// @Summary POST /v1/trovo-api/users/payment
// @Tags servicelinks
// @Accept json
// @Produce json
// @Param body body paymentModels.PaymentInfo true "Payment payload"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /v1/trovo-api/users/payment [post]
func postTrovoApiUsersPaymentHandler(callBackRetryChan chan retryCallbacks, gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {

		mInfo, err := servicelinkServices.GetServiceLinkByAPIKey(middleware.ExtractServiceLinkApiKey(c), gc.DB)

		if err != nil {
			log.Println("[GET service] error for service:", middleware.ExtractServiceLinkApiKey(c), "error: ", err)

			var ex tErrors.GenericError
			var ok bool

			ex, ok = err.(tErrors.GenericError)
			var statusCode int = 0
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

		// ownerUsername := mInfo.OwnerUsername

		if mInfo.CreateUsersPermission == 0 {
			//wrong access
			statusCode := http.StatusUnauthorized
			response := gin.H{"error": "error-invalid-service-access", "data": "Permission", "message": "Permission to create users not enabled for this service"}
			c.JSON(statusCode, response)
			return
		}

		//get user DB record
		accountSignerUser, getUserError := userModels.UserSigner(middleware.ExtractSigner(c)).GetOwner(gc.DB, gc)

		if getUserError != nil {
			log.Printf("[FAILED PAYMENT] ERROR GETTING USER FROM DB from [%v], error: [%v]\n", middleware.ExtractSigner(c), getUserError)

			var ex tErrors.GenericError
			var ok bool

			ex, ok = getUserError.(tErrors.GenericError)
			if ok {
				c.JSON(ex.HTTPCode(), ex.JSONError())
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": getUserError.Error()})
			}
			return
		}
		signerAccountAlias := accountSignerUser.Username
		if signerAccountAlias == os.Getenv("LOG_TARGET_USER") || middleware.ExtractPublicKey(c) == os.Getenv("LOG_TARGET_USER_PK") {
			log.Printf("[CUSTOM LOG] %v error:%v\n", signerAccountAlias, getUserError)
		}
		if accountSignerUser.Suspended == 1 {
			getUserError = &tErrors.CustomError{
				Param:      "Username",
				Err:        "error-account-suspended",
				ErrMessage: "Your account is currently suspended. Please contact support (support@trovotech.io) for more information.",
				Code:       http.StatusForbidden,
			}

			var ex tErrors.GenericError
			var ok bool

			ex, ok = getUserError.(tErrors.GenericError)
			if ok {
				c.JSON(ex.HTTPCode(), ex.JSONError())
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": getUserError.Error()})
			}
			return

		}

		{
			//check if pending shared access modify exists
			if userServices.CheckPendingSharedAccessApproval(middleware.ExtractPublicKey(c), gc.DB) {
				c.JSON(http.StatusForbidden, gin.H{"error": "error-pending-shared-access-op", "message": "There is a pending shared access operation on this wallet and must be completed first before attempting to send payment from this wallet."})
				return
			}
		}

		var paymentInfo paymentModels.PaymentInfo
		// var err error

		data, _ := io.ReadAll(c.Request.Body)

		err = json.Unmarshal(data, &paymentInfo)

		var invalidJSON tErrors.ErrorInvalidJSON

		if err != nil {
			log.Printf("[FAILED PAYMENT] UNMARSHAL FAIL from [%v], error: [%v]\n", signerAccountAlias, err)

			log.Print(err)
			c.JSON(invalidJSON.HTTPCode(), invalidJSON.JSONError())
			return
		}
		if signerAccountAlias == os.Getenv("LOG_TARGET_USER") || middleware.ExtractPublicKey(c) == os.Getenv("LOG_TARGET_USER_PK") {
			log.Printf("[CUSTOM LOG] paymentInfo %+v\n", paymentInfo)
		}

		conDB.PrintDBStats(fmt.Sprintf("POST /v1/users/payment %v", signerAccountAlias), gc.DB)

		//initiatlize message holder
		paymentInfo.Messages = make([]string, 0)
		//check if username is reserved. Reserved usernames should not send payments.

		_, checkReservedUserError := usersDB.UsernameIsReserved(signerAccountAlias, gc.DB)
		if checkReservedUserError != nil {

			var ex tErrors.GenericError
			var ok bool

			ex, ok = checkReservedUserError.(tErrors.GenericError)
			if ok {
				c.JSON(ex.HTTPCode(), ex.JSONError())
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": checkReservedUserError.Error()})
			}
			return
		}
		//get the wallet you are sending payment from
		sourceWallet, temp, getWalletError := usersDB.GetWallet(middleware.ExtractPublicKey(c), gc.DB)
		if signerAccountAlias == os.Getenv("LOG_TARGET_USER") || middleware.ExtractPublicKey(c) == os.Getenv("LOG_TARGET_USER_PK") {
			log.Printf("[CUSTOM LOG] %v error:%v\n", signerAccountAlias, getWalletError)
		}
		if getWalletError != nil {
			log.Printf("[FAILED PAYMENT] ERROR GETTING USER FROM DB from [%v], error: [%v]\n", signerAccountAlias, getWalletError)

			var ex tErrors.GenericError
			var ok bool

			ex, ok = getWalletError.(tErrors.GenericError)
			if ok {
				c.JSON(ex.HTTPCode(), ex.JSONError())
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": getWalletError.Error(), "message": getWalletError.Error()})
			}
			return
		}

		if sourceWallet.WalletType == 1 || sourceWallet.WalletType == 2 || sourceWallet.WalletType == 3 {
			c.JSON(http.StatusForbidden, gin.H{"error": "error-wallet-type-not-allowed", "message": "Token Minting/Market Making/Bulk Payment wallets are not allowed for this operation."})
			return
		}

		if temp {
			errAccountIsTemp := &tErrors.CustomError{
				Param:      "Username",
				Err:        "error-account-not-primary-account-alias",
				ErrMessage: "only primary/subwallets are allowed for payment requests",
				Code:       http.StatusForbidden,
			}

			c.JSON(errAccountIsTemp.HTTPCode(), errAccountIsTemp.JSONError())
			return

		}
		if sourceWallet.WalletType != 0 {
			c.JSON(http.StatusForbidden, gin.H{"error": "error-wallet-type-forbidden", "message": "Operation not allowed on any special type of wallets. Only standard wallets are allowed."})
			return
		}
		{
			//prevent wallets with approver from using this endpoint
			if sourceWallet.SharedAccessEnabled == 1 && sourceWallet.WalletCountApproverAccess(gc) > 0 {
				errAccountIsTemp := &tErrors.CustomError{
					Param:      "ID",
					Err:        "error-wallet-with-shared-access-not-allowed",
					ErrMessage: "This wallet has approver access enabled. Please let someone with an INITIATOR access submit the request.",
					Code:       http.StatusForbidden,
				}

				c.JSON(errAccountIsTemp.HTTPCode(), errAccountIsTemp.JSONError())
				return
			}
		}
		var destinationUser userModels.User
		var destinationWallet userModels.UserWallet
		var getDestinationUserError, getDestinationWalletError error
		//check if the public key exists in TROVO and then transform to username
		paymentInfo.Messages = make([]string, 0)
		publicKeyPayment := len(paymentInfo.Destination) == 56 || len(paymentInfo.Destination) == 69
		if publicKeyPayment {
			paymentInfo.Destination = strings.ToUpper(paymentInfo.Destination)
			destinationWallet, _, getDestinationWalletError = usersDB.GetWallet(paymentInfo.Destination, gc.DB)
			if getDestinationWalletError == nil {
				destinationUser, _ = destinationWallet.GetWalletOwner(gc.DB, gc)
				paymentInfo.Messages = append(paymentInfo.Messages, fmt.Sprintf("Notice: Address[%v] belongs to the wallet alias [%v] and has been used as destination", paymentInfo.Destination, destinationWallet.Alias))
				paymentInfo.Destination = destinationWallet.Alias
				publicKeyPayment = false
			}
		} else {
			//check if it is an email
			if strings.Contains(paymentInfo.Destination, "@") {
				//an email...replace the user info
				destinationUser, err = usersDB.GetUser(paymentInfo.Destination, gc.DB, gc)
				if err == nil {
					destinationWallet, _, getDestinationWalletError = usersDB.GetWallet(destinationUser.Username, gc.DB)
					if getDestinationWalletError == nil {
						paymentInfo.Messages = append(paymentInfo.Messages, fmt.Sprintf("Notice: Email [%v] belongs to the username [%v] and has been used as destination", paymentInfo.Destination, destinationUser.Username))
						paymentInfo.Destination = destinationUser.Username
						publicKeyPayment = false
					}
				}
			} else if strings.Contains(paymentInfo.Destination, "+") {
				//a phone...replace the user info
				destinationUser, err = usersDB.GetUser(paymentInfo.Destination, gc.DB, gc)
				if err == nil {
					destinationWallet, _, getDestinationWalletError = usersDB.GetWallet(destinationUser.Username, gc.DB)
					if getDestinationWalletError == nil {
						paymentInfo.Messages = append(paymentInfo.Messages, fmt.Sprintf("Notice: Phone [%v] belongs to the username [%v] and has been used as destination", paymentInfo.Destination, destinationUser.Username))
						paymentInfo.Destination = destinationUser.Username
						publicKeyPayment = false
					}
				}
			} else {
				//wallet alias...
				destinationWallet, _, getDestinationWalletError = usersDB.GetWallet(paymentInfo.Destination, gc.DB)
				if getDestinationWalletError == nil {
					//get destination user:
					destinationUser, _ = destinationWallet.GetWalletOwner(gc.DB, gc)
					publicKeyPayment = false
				}
			}

			paymentInfo.Destination = strings.ToLower(paymentInfo.Destination)
		}

		//check if receiver is reserved. Reserved usernames should not be sent payments.

		if !publicKeyPayment {
			//skip public key payments
			_, checkReservedReceiverError := usersDB.UsernameIsReserved(paymentInfo.Destination, gc.DB)
			if checkReservedReceiverError != nil {

				var ex tErrors.GenericError
				var ok bool

				ex, ok = checkReservedReceiverError.(tErrors.GenericError)
				if ok {
					c.JSON(ex.HTTPCode(), ex.JSONError())
				} else {
					c.JSON(http.StatusBadRequest, gin.H{"error": checkReservedReceiverError.Error()})
				}
				return
			}

			destinationUser, getDestinationUserError = usersDB.GetUser(paymentInfo.Destination, gc.DB, gc)
			if getDestinationUserError != nil {
				ex := &tPayErrors.ErrorPaymentDestinationDoesNotExist{}
				c.JSON(ex.HTTPCode(), ex.JSONError())
				return
			}

		}
		paymentInfoReturned, returnedDestination, paymentError := userServices.Pay(&accountSignerUser, &sourceWallet, &paymentInfo, gc)
		if signerAccountAlias == os.Getenv("LOG_TARGET_USER") || middleware.ExtractPublicKey(c) == os.Getenv("LOG_TARGET_USER_PK") {
			log.Printf("[CUSTOM LOG] returned Payment Error: [%v]\n", paymentError)

			if paymentInfoReturned != nil {
				log.Printf("[CUSTOM LOG] returned Payment Info: [%+v]\n", *paymentInfoReturned)
			}

		}
		if paymentError != nil {
			log.Printf("[FAILED PAYMENT]  from [%v] to [%v], error: [%v]\n", signerAccountAlias, paymentInfo.Destination, paymentError)
			var ex tErrors.GenericError
			var ok bool

			ex, ok = paymentError.(tErrors.GenericError)
			if ok {
				c.JSON(ex.HTTPCode(), ex.JSONError())
			} else {
				log.Print(paymentError)
				c.JSON(http.StatusBadRequest, gin.H{"error": paymentError.Error()})
			}

			return
		}

		if len(paymentInfo.TransactionID) == 0 {
			c.JSON(http.StatusAccepted, paymentInfoReturned)

		} else {

			//log user current location
			accountSignerUser.PublicIP = c.ClientIP()
			if len(c.GetHeader("Cf-Connecting-Ip")) > 4 {
				accountSignerUser.PublicIP = c.GetHeader("Cf-Connecting-Ip")
			}

			paymentServices.UpdateAndLogUserPaymentGeoInformation(&accountSignerUser, paymentInfoReturned, gc.DB)
			senderPaymentHistoryCacheKey := fmt.Sprintf("[GET] /v1/users/payments/%v", middleware.ExtractPublicKey(c))
			senderCacheKey := fmt.Sprintf("[GET] /v1/users/%v", signerAccountAlias)

			gc.RedisCache.InvalidateCachedHttpResponse(senderCacheKey, senderPaymentHistoryCacheKey)
			gc.RedisCache.InvalidateCachedHttpResponse(senderPaymentHistoryCacheKey)
			var senderBalanceCacheKey, senderTempCacheKey, receiverBalanceCacheKey, receiverTempCacheKey, rNTF, sNFT string

			senderBalanceCacheKey = fmt.Sprintf("GetBalance_%s", middleware.ExtractPublicKey(c))
			sNFT = fmt.Sprintf("GetNFTs_%s", middleware.ExtractPublicKey(c))
			if returnedDestination != nil {
				receiverBalanceCacheKey = fmt.Sprintf("GetBalance_%s", returnedDestination.PublicKey)
				rNTF = fmt.Sprintf("GetNFTs_%s", returnedDestination.PublicKey)

			}
			if len(destinationWallet.ID) == 56 {
				destinationWallet.InvalidateUserCache(gc)
				if destinationWallet.TempPublicKey != nil {

					receiverTempCacheKey = fmt.Sprintf("GetBalance_%s", *destinationWallet.TempPublicKey)
				}
			}
			if len(sourceWallet.ID) == 56 {
				if sourceWallet.TempPublicKey != nil {

					senderTempCacheKey = fmt.Sprintf("GetBalance_%s", *sourceWallet.TempPublicKey)
				}

			}

			if returnedDestination != nil {
				destinationUsername := strings.TrimSpace(strings.ToLower(destinationUser.Username))

				receiverCacheKey := fmt.Sprintf("[GET] /v1/users/%v", destinationUsername)
				receiverPaymentHistoryCacheKey := fmt.Sprintf("[GET] /v1/users/payments/%v", middleware.ExtractPublicKey(c))
				gc.RedisCache.InvalidateCachedHttpResponse(receiverCacheKey, receiverPaymentHistoryCacheKey)
				gc.RedisCache.InvalidateCachedHttpResponse(receiverPaymentHistoryCacheKey, senderBalanceCacheKey, receiverBalanceCacheKey)
				returnedDestination.InvalidateUserCache(gc)
				destinationWallet.InvalidateUserCache(gc)
			}
			gc.RedisCache.InvalidateCachedHttpResponse(senderBalanceCacheKey, senderTempCacheKey, receiverBalanceCacheKey, receiverTempCacheKey, sNFT, rNTF)
			accountSignerUser.InvalidateUserCache(gc)
			sourceWallet.InvalidateUserCache(gc)
			c.JSON(http.StatusOK, paymentInfoReturned)

			{

				//start callback process here

				if d, ok := paymentInfoReturned.CallbackURLS["orderPaymentCallbackUrl"]; ok && len(d) > 5 {
					log.Printf("[paymentNotification] found notification callbackUrl: [%v]\n\n", d)
					//make callback request
					// callbackResponse := new(map[string]interface{})
					type payload struct {
						Destination     string    `json:"destination"`
						Sender          string    `json:"sender"`
						Amount          string    `json:"amount"`
						AssetCode       string    `json:"assetCode"`
						AssetIssuer     string    `json:"assetIssuer"`
						TransactionID   string    `json:"transactionId"`
						TransactionMemo string    `json:"transactionMemo"`
						TransactionTime time.Time `json:"transactionTime"`
						DeviceID        string    `json:"deviceId"`
					}
					assetCode := paymentInfoReturned.AssetCode
					if paymentInfoReturned.AssetIssuer == "" {
						assetCode = os.Getenv("NATIVE_ASSET_CODE")
					}
					senderWallet, _, _ := usersDB.GetWallet(middleware.ExtractPublicKey(c), gc.DB)
					jsonPayload := payload{
						Destination:     paymentInfoReturned.Destination,
						Sender:          senderWallet.Alias,
						Amount:          paymentInfoReturned.Amount,
						AssetCode:       assetCode,
						AssetIssuer:     paymentInfoReturned.AssetIssuer,
						TransactionID:   paymentInfoReturned.TransactionID,
						TransactionMemo: paymentInfoReturned.Memo,
						TransactionTime: time.Now(),
						DeviceID:        c.GetHeader("X-TW-DEVICE-ID"),
					}
					/////
					body, err := json.Marshal(jsonPayload)
					if err != nil {
						log.Printf("[paymentNotification] could not unmarshal callback message due to [%v]\n", err)

					}
					log.Printf("[paymentNotification] JSON STRING: [%v]\n", string(body))

					responseBody := bytes.NewBuffer(body)
					//Leverage Go's HTTP Post function to make request
					c := retryCallbacks{Req: responseBody, CallbackURL: d, Count: 0}
					callBackRetryChan <- c
				}
			}

			{
				// send push notifications
				assetCode := paymentInfo.AssetCode
				if assetCode == "" {
					assetCode = os.Getenv("NATIVE_ASSET_CODE")
				}
				dataPayload := make(map[string]string)
				dataPayload["route"] = "basicTransactionHistory"
				if !publicKeyPayment {
					// dataPayload := make(map[string]string)
					// dataPayload["none"] = ""
					if destinationWallet.SharedAccessEnabled == 1 {
						if destinationWallet.HasViewOnlyAccess(gc) {
							u, e := destinationWallet.GetWalletOwner(gc.DB, gc)
							if e == nil {
								if u.PushNotificationToken != nil {

									u.SendPushMessage("Trovo: Shared Wallet Credited!", fmt.Sprintf("You have received %v %v from %v to your shared wallet with alias %v", paymentInfo.Amount, assetCode, sourceWallet.Alias, paymentInfo.Destination), "", dataPayload, gc)
									u.InvalidateUserCache(gc)

								}
							}

						}
						for _, v := range destinationWallet.Permissions {
							u, e := userModels.Username(v.TargetUsername).GetSimpleUser(gc.DB, gc)
							if e != nil {
								continue
							}
							if u.PushNotificationToken != nil {

								u.SendPushMessage("Trovo: Shared Wallet Credited!", fmt.Sprintf("You have received %v %v from %v to your shared wallet with alias %v", paymentInfo.Amount, assetCode, sourceWallet.Alias, paymentInfo.Destination), "", dataPayload, gc)
								u.InvalidateUserCache(gc)
							}
						}
					} else {
						// 				if !publicKeyPayment {
						// 	dw, _ := userModels.WalletAlias(paymentInfoReturned.Destination).GetWallet(gc.DB, gc)
						// 	destinationUser, _ = dw.GetWalletOwner(gc.DB, gc)
						// 	destinationUser.SendPushMessage("Trovo: Wallet Credited!", fmt.Sprintf("You have received %v %v from %v to your wallet with alias %v", paymentInfo.Amount, assetCode, sourceWallet.Alias, paymentInfo.Destination), "", dataPayload, gc)
						// 	destinationUser.InvalidateUserCache(gc)
						// }
						dw, _ := userModels.WalletAlias(paymentInfoReturned.Destination).GetWallet(gc.DB, gc)
						u, e := dw.GetWalletOwner(gc.DB, gc)
						if e == nil {
							if u.PushNotificationToken != nil {

								u.SendPushMessage("Trovo: Wallet Credited!", fmt.Sprintf("You have received %v %v from %v to your wallet with alias %v", paymentInfo.Amount, assetCode, sourceWallet.Alias, paymentInfo.Destination), "", dataPayload, gc)
								u.InvalidateUserCache(gc)

							}
						}
					}
				}
				accountSignerUser.SendPushMessage("Trovo: Wallet Debited!", fmt.Sprintf("You have successfully sent %v %v from your wallet with alias %v to %v", paymentInfo.Amount, assetCode, sourceWallet.Alias, paymentInfo.Destination), "", dataPayload, gc)

			}

		}

	}
}

// postTrovoApiUsersSubwalletHandler godoc
// @Summary POST /v1/trovo-api/users/subwallet
// @Tags servicelinks
// @Accept json
// @Produce json
// @Param body body userModels.ServiceLinkSubWalletInfo true "Sub-wallet payload"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /v1/trovo-api/users/subwallet [post]
func postTrovoApiUsersSubwalletHandler(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {

		mInfo, err := servicelinkServices.GetServiceLinkByAPIKey(middleware.ExtractServiceLinkApiKey(c), gc.DB)

		if err != nil {
			log.Println("[GET service] error for service:", middleware.ExtractServiceLinkApiKey(c), "error: ", err)

			var ex tErrors.GenericError
			var ok bool

			ex, ok = err.(tErrors.GenericError)
			var statusCode int = 0
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

		// ownerUsername := mInfo.OwnerUsername

		if mInfo.CreateUsersPermission == 0 {
			//wrong access
			statusCode := http.StatusUnauthorized
			response := gin.H{"error": "error-invalid-service-access", "data": "Permission", "message": "Permission to create users not enabled for this service"}
			c.JSON(statusCode, response)
			return
		}

		signerUser, err := usersDB.GetUser(mInfo.OwnerUsername, gc.DB, gc)

		if err != nil {
			var ex tErrors.GenericError
			var ok bool

			ex, ok = err.(tErrors.GenericError)
			if ok {
				c.JSON(http.StatusBadRequest, ex.JSONError())
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": err.Error()})
			}
			return
		}

		var input userModels.ServiceLinkSubWalletInfo
		// var err error

		data, _ := io.ReadAll(c.Request.Body)

		err = json.Unmarshal(data, &input)

		var invalidJSON tErrors.ErrorInvalidJSON

		if err != nil {
			c.JSON(http.StatusBadRequest, invalidJSON.JSONError())
			return
		}

		conDB.PrintDBStats(fmt.Sprintf("POST /v1/trovo-api/users/subwallet %v", signerUser.Username), gc.DB)
		subWalletInfo := input.ToSubwalletInfo(gc)
		returnedSubwalletInfo, err := userServices.CreateNewSubWallet(&signerUser, &subWalletInfo, gc)

		if err != nil {
			var ex tErrors.GenericError
			var ok bool

			ex, ok = err.(tErrors.GenericError)
			if ok {
				c.JSON(http.StatusBadRequest, ex.JSONError())
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": err.Error()})
			}
			return
		}

		if signerUser.PushNotificationToken != nil && len(returnedSubwalletInfo.TransactionID) > 0 && returnedSubwalletInfo.TransactionID != "PENDING_AUTH" {
			dataPayload := make(map[string]string)
			dataPayload["route"] = ""
			pns.SendFirebaseMessage(*signerUser.PushNotificationToken, "New Sub-wallet Added!", fmt.Sprintf("You have successfully added a new sub wallet tagged [%v].", returnedSubwalletInfo.WalletTag), "", dataPayload, gc.PushNotificationClient, gc.PNSContext)
		}

		//At this point, there was no error.
		input = returnedSubwalletInfo.ToServiceLinkSubwalletInfo(gc)
		c.JSON(http.StatusOK, input)
	}
}

// getTrovoApiAssetsParametersHandler godoc
// @Summary GET /v1/trovo-api/assets/parameters
// @Tags servicelinks
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /v1/trovo-api/assets/parameters [get]
func getTrovoApiAssetsParametersHandler(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {

		mInfo, err := servicelinkServices.GetServiceLinkByAPIKey(middleware.ExtractServiceLinkApiKey(c), gc.DB)

		if err != nil {
			log.Println("[GET service] error for service:", middleware.ExtractServiceLinkApiKey(c), "error: ", err)

			var ex tErrors.GenericError
			var ok bool

			ex, ok = err.(tErrors.GenericError)
			var statusCode int = 0
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

		// ownerUsername := mInfo.OwnerUsername

		if mInfo.CreateUsersPermission == 0 {
			//wrong access
			statusCode := http.StatusUnauthorized
			response := gin.H{"error": "error-invalid-service-access", "data": "Permission", "message": "Permission to create users not enabled for this service"}
			c.JSON(statusCode, response)
			return
		}

		sectorList := userServices.GetTokenizedAssetSectorList(gc.DB)
		subsectorList := userServices.GetTokenizedAssetSubSectorList(gc.DB)
		assetTypes := userServices.GetTokenizedAssetTypes(gc.DB)
		docTypes := userServices.GetAssetTokenizationDocumentTypes(gc.DB)
		custdians := userServices.GetApprovedAssetCustodians(gc.DB)
		managers := userServices.GetAssetManagers(gc.DB)
		issuers := userServices.GetAssetIssuingHouses(gc.DB)
		fees := userServices.GetTokenizationFees(gc.DB)
		statuses := userServices.GetTokenizationStatuses(gc.DB)
		// log.Printf("\n[TOKENIZATION FEES] %+v\n\n", fees)
		currencies := userServices.GetTokenizationCurrencies(gc.DB)
		apo := userServices.GetAssetProtectionOptions(gc.DB)
		apc := userServices.GetAssetProceedCycle(gc.DB)
		ac := userServices.GetTokenizationPublicAssetAllowedCountries(gc.DB)
		fpms := userServices.GetTokenizationFeePaymentMethods(gc.DB)
		countries := userServices.GetCountries(gc.DB)
		countryConfigs := userServices.GetCountryConfigs(gc.DB)

		c.JSON(http.StatusOK, gin.H{"assetSectors": sectorList, "assetSubSectors": subsectorList, "assetTypes": assetTypes, "assetCustodians": custdians, "assetManagers": managers,
			"assetIssuingHouses": issuers, "tokenizationFees": fees, "tokenizationCurrencies": currencies, "assetProtectionOptions": apo, "assetProceedCycle": apc,
			"publicListingAllowedCountries": ac, "tokenizationDocumentTypes": docTypes,
			"tokenizationStatuses": statuses, "feePaymentMethods": fpms, "countryConfigs": countryConfigs, "countries": countries})

	}
}

// getTrovoApiAssetsBankListCountryCodeHandler godoc
// @Summary GET /v1/trovo-api/assets/bank-list/:countryCode
// @Tags servicelinks
// @Produce json
// @Param countryCode path string true "Country code"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /v1/trovo-api/assets/bank-list/{countryCode} [get]
func getTrovoApiAssetsBankListCountryCodeHandler(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {

		mInfo, err := servicelinkServices.GetServiceLinkByAPIKey(middleware.ExtractServiceLinkApiKey(c), gc.DB)

		if err != nil {
			log.Println("[GET service] error for service:", middleware.ExtractServiceLinkApiKey(c), "error: ", err)

			var ex tErrors.GenericError
			var ok bool

			ex, ok = err.(tErrors.GenericError)
			var statusCode int = 0
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

		// ownerUsername := mInfo.OwnerUsername

		if mInfo.CreateUsersPermission == 0 {
			//wrong access
			statusCode := http.StatusUnauthorized
			response := gin.H{"error": "error-invalid-service-access", "data": "Permission", "message": "Permission to create users not enabled for this service"}
			c.JSON(statusCode, response)
			return
		}
		countryCode := c.Param("countryCode")
		// cacheKey := fmt.Sprintf("[GET] /v1/patron/%v", identifier)

		bankList := userServices.GetBanks(countryCode, gc.DB)

		c.JSON(http.StatusOK, gin.H{"bankList": bankList})

	}
}

// getTrovoApiAssetsAdminListHandler godoc
// @Summary GET /v1/trovo-api/assets/admin/list
// @Tags servicelinks
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /v1/trovo-api/assets/admin/list [get]
func getTrovoApiAssetsAdminListHandler(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {

		mInfo, err := servicelinkServices.GetServiceLinkByAPIKey(middleware.ExtractServiceLinkApiKey(c), gc.DB)

		if err != nil {
			log.Println("[GET service] error for service:", middleware.ExtractServiceLinkApiKey(c), "error: ", err)

			var ex tErrors.GenericError
			var ok bool

			ex, ok = err.(tErrors.GenericError)
			var statusCode int = 0
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

		// ownerUsername := mInfo.OwnerUsername

		if mInfo.CreateUsersPermission == 0 {
			//wrong access
			statusCode := http.StatusUnauthorized
			response := gin.H{"error": "error-invalid-service-access", "data": "Permission", "message": "Permission to create users not enabled for this service"}
			c.JSON(statusCode, response)
			return
		}

		user, getUserError := userModels.Username(mInfo.OwnerUsername).GetFullUser(gc.DB, gc)

		if getUserError != nil {

			var ex tErrors.GenericError
			var ok bool

			ex, ok = getUserError.(tErrors.GenericError)
			if ok {
				c.JSON(ex.HTTPCode(), ex.JSONError())
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": getUserError.Error(), "message": getUserError.Error()})
			}
			return
		}
		tokenizationList := userServices.GetTokenizationList(&user, true, gc, c)

		c.JSON(http.StatusOK, tokenizationList)

	}
}

// getTrovoApiAssetsMarketplaceListHandler godoc
// @Summary GET /v1/trovo-api/assets/marketplace/list
// @Tags servicelinks
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /v1/trovo-api/assets/marketplace/list [get]
func getTrovoApiAssetsMarketplaceListHandler(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {

		mInfo, err := servicelinkServices.GetServiceLinkByAPIKey(middleware.ExtractServiceLinkApiKey(c), gc.DB)

		if err != nil {
			log.Println("[GET service] error for service:", middleware.ExtractServiceLinkApiKey(c), "error: ", err)

			var ex tErrors.GenericError
			var ok bool

			ex, ok = err.(tErrors.GenericError)
			var statusCode int = 0
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

		// ownerUsername := mInfo.OwnerUsername

		if mInfo.CreateUsersPermission == 0 {
			//wrong access
			statusCode := http.StatusUnauthorized
			response := gin.H{"error": "error-invalid-service-access", "data": "Permission", "message": "Permission to create users not enabled for this service"}
			c.JSON(statusCode, response)
			return
		}

		user, getUserError := userModels.Username(mInfo.OwnerUsername).GetFullUser(gc.DB, gc)

		if getUserError != nil {

			var ex tErrors.GenericError
			var ok bool

			ex, ok = getUserError.(tErrors.GenericError)
			if ok {
				c.JSON(ex.HTTPCode(), ex.JSONError())
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": getUserError.Error(), "message": getUserError.Error()})
			}
			return
		}
		tokenizationList := userServices.GetTokenizationList(&user, false, gc, c)

		c.JSON(http.StatusOK, tokenizationList)

	}
}

// postTrovoApiAssetsApplyHandler godoc
// @Summary POST /v1/trovo-api/assets/apply
// @Tags servicelinks
// @Accept json
// @Produce json
// @Param body body userModels.TokenizedAssetJSONInput true "Tokenized asset application payload"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /v1/trovo-api/assets/apply [post]
func postTrovoApiAssetsApplyHandler(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {

		mInfo, err := servicelinkServices.GetServiceLinkByAPIKey(middleware.ExtractServiceLinkApiKey(c), gc.DB)

		if err != nil {
			log.Println("[GET service] error for service:", middleware.ExtractServiceLinkApiKey(c), "error: ", err)

			var ex tErrors.GenericError
			var ok bool

			ex, ok = err.(tErrors.GenericError)
			var statusCode int = 0
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

		// ownerUsername := mInfo.OwnerUsername

		if mInfo.CreateUsersPermission == 0 {
			//wrong access
			statusCode := http.StatusUnauthorized
			response := gin.H{"error": "error-invalid-service-access", "data": "Permission", "message": "Permission to create users not enabled for this service"}
			c.JSON(statusCode, response)
			return
		}

		initiator, getUserError := userModels.Username(mInfo.OwnerUsername).GetFullUser(gc.DB, gc)

		if getUserError != nil {

			var ex tErrors.GenericError
			var ok bool

			ex, ok = getUserError.(tErrors.GenericError)
			if ok {
				c.JSON(ex.HTTPCode(), ex.JSONError())
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": getUserError.Error(), "message": getUserError.Error()})
			}
			return
		}

		var tInput userModels.TokenizedAssetJSONInput

		data, _ := io.ReadAll(c.Request.Body)
		// log.Println(string(data))
		err = json.Unmarshal(data, &tInput)

		var invalidJSON tErrors.ErrorInvalidJSON

		if err != nil {
			c.JSON(http.StatusBadRequest, invalidJSON.JSONError())
			return
		}

		//perform request action
		ta, err := userServices.SubmitTokenizationAssetInfoByInitiator(&initiator, &tInput, gc)

		if err != nil {

			var ex tErrors.GenericError
			var ok bool

			ex, ok = err.(tErrors.GenericError)
			if ok {
				c.JSON(ex.HTTPCode(), ex.JSONError())
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": err.Error()})
			}
			return
		}

		c.JSON(http.StatusOK, ta.ToJSON(gc))
	}
}

// putTrovoApiAssetsLogoHandler godoc
// @Summary PUT /v1/trovo-api/assets/logo
// @Tags servicelinks
// @Accept multipart/form-data
// @Produce json
// @Param documentFile formData file true "Logo image file"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /v1/trovo-api/assets/logo [put]
func putTrovoApiAssetsLogoHandler(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {

		mInfo, err := servicelinkServices.GetServiceLinkByAPIKey(middleware.ExtractServiceLinkApiKey(c), gc.DB)

		if err != nil {
			log.Println("[GET service] error for service:", middleware.ExtractServiceLinkApiKey(c), "error: ", err)

			var ex tErrors.GenericError
			var ok bool

			ex, ok = err.(tErrors.GenericError)
			var statusCode int = 0
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

		// ownerUsername := mInfo.OwnerUsername

		if mInfo.CreateUsersPermission == 0 {
			//wrong access
			statusCode := http.StatusUnauthorized
			response := gin.H{"error": "error-invalid-service-access", "data": "Permission", "message": "Permission to create users not enabled for this service"}
			c.JSON(statusCode, response)
			return
		}

		initiator, getUserError := userModels.Username(mInfo.OwnerUsername).GetFullUser(gc.DB, gc)

		if getUserError != nil {

			var ex tErrors.GenericError
			var ok bool

			ex, ok = getUserError.(tErrors.GenericError)
			if ok {
				c.JSON(ex.HTTPCode(), ex.JSONError())
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": getUserError.Error(), "message": getUserError.Error()})
			}
			return
		}

		const MAX_UPLOAD_SIZE = 1024 * 1024 // 1MB
		r := c.Request
		// r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)
		if err := r.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "document cannot be more than 900kb in file size", "message": "document cannot be more than 900kb in file size"})
			return
		}

		f, fileHeader, err := r.FormFile("documentFile")

		if err != nil {
			log.Printf("Error Getting Uploaded file with param DocumentFile:%v\n", err)
			c.JSON(http.StatusForbidden, gin.H{"error": "error-no-ducument-file", "message": "There is no documentFile attached with request"})
			return
		}
		defer f.Close()
		blobFile, err := fileHeader.Open()

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "error attempting to validate the document uploaded", "message": "error attempting to validate the logo uploaded"})

			return
		}
		defer blobFile.Close()

		fnameSplit := strings.Split(fileHeader.Filename, ".")
		fileExtension := fnameSplit[len(fnameSplit)-1]

		{
			//check for unsupported extension
			if !strings.EqualFold(fileExtension, "jpg") && !strings.EqualFold(fileExtension, "jpeg") && !strings.EqualFold(fileExtension, "png") && !strings.EqualFold(fileExtension, "gif") {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Unsurported document format. Only jpg, jpeg, png, gif are supported", "message": "Unsurported document format. Only jpg, jpeg, png and gif are supported"})

				return
			}
		}

		t, _, _ := userModels.Username(initiator.Username).GetOpenTokenizedAssetByInitiatorUsername(gc.DB)

		if len(t.ID) < 5 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Tokenized Asset not valid", "message": "Tokenized Asset not valid"})
			return
		}
		if t.AssetTokenizationStatus > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Tokenization Request cannot be altered at this stage through this option. Please use the option within the tokenization detail."})
			return
		}

		url, err := userServices.UploadTokenizationAssetLogo(&initiator, &t, blobFile, fmt.Sprintf("%s-%s-%s.%s", initiator.Username, "logo", t.ID, fileExtension), gc)

		if err != nil {
			var ex tErrors.GenericError
			var ok bool

			ex, ok = err.(tErrors.GenericError)
			if ok {
				c.JSON(http.StatusBadRequest, ex.JSONError())
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": err.Error()})
			}
			return
		}

		if initiator.PushNotificationToken != nil && len(url) > 0 {
			dataPayload := make(map[string]string)
			dataPayload["route"] = ""
			pns.SendFirebaseMessage(*initiator.PushNotificationToken, "Asset Logo updated!", "You have successfully uploaded asset logo.", url, dataPayload, gc.PushNotificationClient, gc.PNSContext)
		}

		userCacheKey := fmt.Sprintf("[GET] /v1/users/%v", initiator.Username)

		gc.RedisCache.InvalidateCachedHttpResponse(userCacheKey)

		//At this point, there was no error.

		c.JSON(http.StatusOK, url)
	}
}

// putTrovoApiAssetsDocumentsHandler godoc
// @Summary PUT /v1/trovo-api/assets/documents
// @Tags servicelinks
// @Accept multipart/form-data
// @Produce json
// @Param documentFile formData file true "Document file"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /v1/trovo-api/assets/documents [put]
func putTrovoApiAssetsDocumentsHandler(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {

		mInfo, err := servicelinkServices.GetServiceLinkByAPIKey(middleware.ExtractServiceLinkApiKey(c), gc.DB)

		if err != nil {
			log.Println("[GET service] error for service:", middleware.ExtractServiceLinkApiKey(c), "error: ", err)

			var ex tErrors.GenericError
			var ok bool

			ex, ok = err.(tErrors.GenericError)
			var statusCode int = 0
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

		// ownerUsername := mInfo.OwnerUsername

		if mInfo.CreateUsersPermission == 0 {
			//wrong access
			statusCode := http.StatusUnauthorized
			response := gin.H{"error": "error-invalid-service-access", "data": "Permission", "message": "Permission to create users not enabled for this service"}
			c.JSON(statusCode, response)
			return
		}

		initiator, getUserError := userModels.Username(mInfo.OwnerUsername).GetFullUser(gc.DB, gc)

		if getUserError != nil {

			var ex tErrors.GenericError
			var ok bool

			ex, ok = getUserError.(tErrors.GenericError)
			if ok {
				c.JSON(ex.HTTPCode(), ex.JSONError())
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": getUserError.Error(), "message": getUserError.Error()})
			}
			return
		}

		const MAX_UPLOAD_SIZE = 1024 * 1024 // 1MB
		r := c.Request
		// r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)
		if err := r.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "document cannot be more than 900kb in file size", "message": "document cannot be more than 900kb in file size"})
			return
		}

		f, fileHeader, err := r.FormFile("documentFile")

		if err != nil {
			log.Printf("Error Getting Uploaded file with param DocumentFile:%v\n", err)
			c.JSON(http.StatusForbidden, gin.H{"error": "error-no-ducument-file", "message": "There is no documentFile attached with request"})
			return
		}
		defer f.Close()
		blobFile, err := fileHeader.Open()

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "error attempting to validate the document uploaded", "message": "error attempting to validate the document uploaded"})

			return
		}
		defer blobFile.Close()

		fnameSplit := strings.Split(fileHeader.Filename, ".")
		fileExtension := fnameSplit[len(fnameSplit)-1]

		{
			//check for unsupported extension
			if !strings.EqualFold(fileExtension, "jpg") && !strings.EqualFold(fileExtension, "jpeg") && !strings.EqualFold(fileExtension, "png") && !strings.EqualFold(fileExtension, "gif") && !strings.EqualFold(fileExtension, "pdf") {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Unsurported document format. Only jpg, jpeg, png, gif and pdf are supported", "message": "Unsurported document format. Only jpg, jpeg, png, gif and pdf are supported"})

				return
			}
		}

		var tokenizationInput userModels.AssetTokenizationInputDocument

		err = c.ShouldBind(&tokenizationInput)
		// data, _ := io.ReadAll(c.Request.Body)
		// // log.Println(string(data))
		// err = json.Unmarshal(data, &tokenizationInput)

		var invalidJSON tErrors.ErrorInvalidJSON

		if err != nil {
			log.Printf("Error Getting Uploaded file with param DocumentFile:%+v\n error: %v", r.Body, err)

			c.JSON(http.StatusBadRequest, invalidJSON.JSONError())
			return
		}
		if tokenizationInput.DocumentType == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "document type not specified", "message": "document type not specified"})
			return
		}
		if len(tokenizationInput.DocumentTitle) < 5 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Document title not valid. Must be at least 5 characters long", "message": "Document title not valid. Must be at least 5 characters long"})
			return
		}
		t, _, _ := userModels.Username(initiator.Username).GetOpenTokenizedAssetByInitiatorUsername(gc.DB)

		if len(t.ID) < 5 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Tokenized Asset not valid", "message": "Tokenized Asset not valid"})
			return
		}
		if t.AssetTokenizationStatus > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Tokenization Request cannot be altered at this stage through this option. Please use the option within the tokenization detail."})
			return
		}
		tokenizationInput.TokenizedAssetID = t.ID
		if len(tokenizationInput.TokenizedAssetID) < 5 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Tokenized Asset not valid", "message": "Tokenized Asset not valid. Ensure you ahve an open tokenizstion and try again."})
			return
		}

		url, err := userServices.UploadTokenizationDocument(&initiator, blobFile, fmt.Sprintf("%s-%s-%s.%s", initiator.Username, tokenizationInput.DocumentTitle, tokenizationInput.TokenizedAssetID, fileExtension), &tokenizationInput, gc)

		if err != nil {
			var ex tErrors.GenericError
			var ok bool

			ex, ok = err.(tErrors.GenericError)
			if ok {
				c.JSON(http.StatusBadRequest, ex.JSONError())
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": err.Error()})
			}
			return
		}

		if initiator.PushNotificationToken != nil && len(url) > 0 {
			dataPayload := make(map[string]string)
			dataPayload["route"] = ""
			pns.SendFirebaseMessage(*initiator.PushNotificationToken, tokenizationInput.DocumentTitle+" updated!", fmt.Sprintf("You have successfully uploaded %v[%v].", tokenizationInput.DocumentTitle, tokenizationInput.DocumentType), url, dataPayload, gc.PushNotificationClient, gc.PNSContext)
		}

		userCacheKey := fmt.Sprintf("[GET] /v1/users/%v", initiator.Username)

		gc.RedisCache.InvalidateCachedHttpResponse(userCacheKey)

		//At this point, there was no error.

		c.JSON(http.StatusOK, url)
	}
}

// putTrovoApiAssetsFeesDocumentHandler godoc
// @Summary PUT /v1/trovo-api/assets/fees/document
// @Tags servicelinks
// @Accept multipart/form-data
// @Produce json
// @Param documentFile formData file true "Proof of payment file"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /v1/trovo-api/assets/fees/document [put]
func putTrovoApiAssetsFeesDocumentHandler(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {

		mInfo, err := servicelinkServices.GetServiceLinkByAPIKey(middleware.ExtractServiceLinkApiKey(c), gc.DB)

		if err != nil {
			log.Println("[GET service] error for service:", middleware.ExtractServiceLinkApiKey(c), "error: ", err)

			var ex tErrors.GenericError
			var ok bool

			ex, ok = err.(tErrors.GenericError)
			var statusCode int = 0
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

		if mInfo.CreateUsersPermission == 0 {
			//wrong access
			statusCode := http.StatusUnauthorized
			response := gin.H{"error": "error-invalid-service-access", "data": "Permission", "message": "Permission to create users not enabled for this service"}
			c.JSON(statusCode, response)
			return
		}

		initiator, getUserError := userModels.Username(mInfo.OwnerUsername).GetFullUser(gc.DB, gc)

		if getUserError != nil {

			var ex tErrors.GenericError
			var ok bool

			ex, ok = getUserError.(tErrors.GenericError)
			if ok {
				c.JSON(ex.HTTPCode(), ex.JSONError())
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": getUserError.Error(), "message": getUserError.Error()})
			}
			return
		}

		const MAX_UPLOAD_SIZE = 1024 * 1024 // 1MB
		r := c.Request
		// r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)
		if err = r.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil {
			log.Printf("[TOKENIZE  FEE PAYMENT DOCUMENT] ERROR PARSING MULTIPART FORM, error: [%v]\n", err)

			c.JSON(http.StatusBadRequest, gin.H{"error": "error-invalid-size", "message": "Document cannot be more than 900kb in file size"})
			return
		}

		f, fileHeader, err := r.FormFile("documentFile")
		if err != nil {
			log.Printf("[TOKENIZE  FEE PAYMENT DOCUMENT] Error Getting Uploaded file with param DocumentFile:%v\n", err)
			c.JSON(http.StatusForbidden, gin.H{"error": "error-no-ducument-file", "message": "There is no documentFile attached with request"})
			return
		}
		defer f.Close()
		blobFile, err := fileHeader.Open()

		if err != nil {
			log.Printf("[TOKENIZE  FEE PAYMENT DOCUMENT] Error opening file with param DocumentFile:%v\n", err)

			c.JSON(http.StatusBadRequest, gin.H{"error": "error-unable-to-validate", "message": "error attempting to validate the document uploaded"})

			return
		}
		defer blobFile.Close()

		fnameSplit := strings.Split(fileHeader.Filename, ".")
		fileExtension := fnameSplit[len(fnameSplit)-1]

		{
			//check for unsupported extension
			if !strings.EqualFold(fileExtension, "jpg") && !strings.EqualFold(fileExtension, "jpeg") && !strings.EqualFold(fileExtension, "png") && !strings.EqualFold(fileExtension, "gif") && !strings.EqualFold(fileExtension, "pdf") {
				log.Printf("[TOKENIZE  FEE PAYMENT DOCUMENT] Error unsorported file format file:%v\n", fileExtension)

				c.JSON(http.StatusBadRequest, gin.H{"error": "error-unsurported-format", "message": "Unsurported document format. Only jpg, jpeg, png, gif and pdf are supported"})

				return
			}
		}
		var tokenizationInput userModels.TokenizationFeeProofOfPaymentInput

		err = c.ShouldBind(&tokenizationInput)
		// data, _ := io.ReadAll(c.Request.Body)
		// // log.Println(string(data))
		// err = json.Unmarshal(data, &tokenizationInput)

		var invalidJSON tErrors.ErrorInvalidJSON

		if err != nil {
			log.Printf("Error binding to fileds:%+v\n error: %v", r.Body, err)

			c.JSON(http.StatusBadRequest, invalidJSON.JSONError())
			return
		}
		if tokenizationInput.TokenizationFeePaymentMethodID == "" {
			log.Printf("[TOKENIZE FEE PAYMENT DOCUMENT] Error payment method not specified:%+v\n", tokenizationInput)

			c.JSON(http.StatusBadRequest, gin.H{"error": "payment-type-not-specified", "message": "payment type not specified"})
			return
		}
		//c.Param("tokenizedAssetID")
		t, _, err := userModels.Username(initiator.Username).GetFeeReadyTokenizedAssetApplicationByInitiatorUsername(gc.DB)
		if err != nil {
			log.Printf("[TOKENIZE  FEE PAYMENT DOCUMENT] ERROR GETTING TOKENIZED ASSERT FROM DB [%v], error: [%v]\n", initiator.Username, err)

			var ex tErrors.GenericError
			var ok bool

			ex, ok = err.(tErrors.GenericError)
			if ok {
				c.JSON(ex.HTTPCode(), ex.JSONError())
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": err.Error()})
			}
			return
		}

		if t.InitiatorUsername != initiator.Username {
			log.Printf("[TOKENIZE FEE PAYMENT DOCUMENT] Error invalid tokenized asset by owner:%+v\n", initiator.Username)

			c.JSON(http.StatusBadRequest, gin.H{"error": "error-tokenized-asset-not-valid", "message": "Tokenized Asset not valid for you."})
			return
		}
		if len(t.ID) < 5 {
			log.Printf("[TOKENIZE FEE PAYMENT DOCUMENT] Error invalid tokenized asset by owner:%+v\n", initiator.Username)

			c.JSON(http.StatusBadRequest, gin.H{"error": "error-tokenized-asset-not-valid", "message": "Tokenized Asset not valid"})
			return
		}
		if t.AssetTokenizationStatus < 1 {
			log.Printf("[TOKENIZE FEE PAYMENT DOCUMENT] Error status still open:%+v\n", tokenizationInput)

			c.JSON(http.StatusBadRequest, gin.H{"error": "error-invalid-status", "message": "Tokenization Request cannot be altered at this stage through this option. Please use the option within the tokenization detail."})
			return
		}

		if t.VettingStatus == 0 {
			log.Printf("[TOKENIZE FEE PAYMENT DOCUMENT] Error. still not vetted:%+v\n", tokenizationInput)

			c.JSON(http.StatusBadRequest, gin.H{"error": "error-invalid-status", "message": "tokenization Request is still being vetted by the team, therefore cannot be modified or updated at this time. Please excercise patience."})
			return
		}

		url, err := userServices.UploadTokenizationFeeProofOfPaymentDocument(&initiator, t.ID, blobFile, fmt.Sprintf("%s-%s-%s.%s", initiator.Username, uuid.NewString(), t.ID, fileExtension), &tokenizationInput, gc)

		if err != nil {
			log.Printf("[TOKENIZE FEE PAYMENT DOCUMENT] Error uploading document:%+v\n", err)

			var ex tErrors.GenericError
			var ok bool

			ex, ok = err.(tErrors.GenericError)
			if ok {
				c.JSON(http.StatusBadRequest, ex.JSONError())
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": err.Error()})
			}
			return
		}

		c.JSON(http.StatusOK, url)

		if initiator.PushNotificationToken != nil && len(url) > 0 {
			dataPayload := make(map[string]string)
			dataPayload["route"] = ""
			pns.SendFirebaseMessage(*initiator.PushNotificationToken, "Proof of payment updated!", "You have successfully uploaded proof of payment.", url, dataPayload, gc.PushNotificationClient, gc.PNSContext)
		}

		userCacheKey := fmt.Sprintf("[GET] /v1/users/%v", initiator.Username)

		gc.RedisCache.InvalidateCachedHttpResponse(userCacheKey)

	}
}

// postTrovoApiAssetsFeesConfirmTokenizationIDHandler godoc
// @Summary POST /v1/trovo-api/assets/fees/confirm/:tokenizationID
// @Tags servicelinks
// @Produce json
// @Param tokenizationID path string true "Tokenization ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /v1/trovo-api/assets/fees/confirm/{tokenizationID} [post]
func postTrovoApiAssetsFeesConfirmTokenizationIDHandler(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {

		mInfo, err := servicelinkServices.GetServiceLinkByAPIKey(middleware.ExtractServiceLinkApiKey(c), gc.DB)

		if err != nil {
			log.Println("[GET service] error for service:", middleware.ExtractServiceLinkApiKey(c), "error: ", err)

			var ex tErrors.GenericError
			var ok bool

			ex, ok = err.(tErrors.GenericError)
			var statusCode int = 0
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

		if mInfo.CreateUsersPermission == 0 {
			//wrong access
			statusCode := http.StatusUnauthorized
			response := gin.H{"error": "error-invalid-service-access", "data": "Permission", "message": "Permission to create users not enabled for this service"}
			c.JSON(statusCode, response)
			return
		}

		initiator, getUserError := userModels.Username(mInfo.OwnerUsername).GetFullUser(gc.DB, gc)

		if getUserError != nil {

			var ex tErrors.GenericError
			var ok bool

			ex, ok = getUserError.(tErrors.GenericError)
			if ok {
				c.JSON(ex.HTTPCode(), ex.JSONError())
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": getUserError.Error(), "message": getUserError.Error()})
			}
			return
		}

		//confirm request
		ta, err := userServices.ConfirmTokenizationFeePaymentByInitiator(&initiator, c.Param("tokenizationID"), gc)

		if err != nil {
			log.Printf("[ConfirmTokenizationFeePaymentByInitiator] tokenizedAsset: %+v\n Error: %v\n", err)

			var ex tErrors.GenericError
			var ok bool

			ex, ok = err.(tErrors.GenericError)
			if ok {
				c.JSON(ex.HTTPCode(), ex.JSONError())
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": err.Error()})
			}
			return
		}

		c.JSON(http.StatusOK, ta)

	}
}

// deleteTrovoApiAssetsTokenizationIDHandler godoc
// @Summary DELETE /v1/trovo-api/assets/:tokenizationID
// @Tags servicelinks
// @Produce json
// @Param tokenizationID path string true "Tokenization ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /v1/trovo-api/assets/{tokenizationID} [delete]
func deleteTrovoApiAssetsTokenizationIDHandler(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {

		mInfo, err := servicelinkServices.GetServiceLinkByAPIKey(middleware.ExtractServiceLinkApiKey(c), gc.DB)

		if err != nil {
			log.Println("[GET service] error for service:", middleware.ExtractServiceLinkApiKey(c), "error: ", err)

			var ex tErrors.GenericError
			var ok bool

			ex, ok = err.(tErrors.GenericError)
			var statusCode int = 0
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

		// ownerUsername := mInfo.OwnerUsername

		if mInfo.CreateUsersPermission == 0 {
			//wrong access
			statusCode := http.StatusUnauthorized
			response := gin.H{"error": "error-invalid-service-access", "data": "Permission", "message": "Permission to create users not enabled for this service"}
			c.JSON(statusCode, response)
			return
		}

		initiator, getUserError := userModels.Username(mInfo.OwnerUsername).GetFullUser(gc.DB, gc)

		if getUserError != nil {

			var ex tErrors.GenericError
			var ok bool

			ex, ok = getUserError.(tErrors.GenericError)
			if ok {
				c.JSON(ex.HTTPCode(), ex.JSONError())
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": getUserError.Error(), "message": getUserError.Error()})
			}
			return
		}

		ta, err := userServices.DeleteTokenization(&initiator, c.Param("tokenizationID"), gc)

		if err != nil {

			var ex tErrors.GenericError
			var ok bool

			ex, ok = err.(tErrors.GenericError)
			if ok {
				c.JSON(ex.HTTPCode(), ex.JSONError())
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": err.Error()})
			}
			return
		}

		c.JSON(http.StatusOK, ta)
	}
}

// deleteTrovoApiAssetsDocumentsDocumentIDHandler godoc
// @Summary DELETE /v1/trovo-api/assets/documents/:documentID
// @Tags servicelinks
// @Produce json
// @Param documentID path string true "Document ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /v1/trovo-api/assets/documents/{documentID} [delete]
func deleteTrovoApiAssetsDocumentsDocumentIDHandler(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {

		mInfo, err := servicelinkServices.GetServiceLinkByAPIKey(middleware.ExtractServiceLinkApiKey(c), gc.DB)

		if err != nil {
			log.Println("[GET service] error for service:", middleware.ExtractServiceLinkApiKey(c), "error: ", err)

			var ex tErrors.GenericError
			var ok bool

			ex, ok = err.(tErrors.GenericError)
			var statusCode int = 0
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

		// ownerUsername := mInfo.OwnerUsername

		if mInfo.CreateUsersPermission == 0 {
			//wrong access
			statusCode := http.StatusUnauthorized
			response := gin.H{"error": "error-invalid-service-access", "data": "Permission", "message": "Permission to create users not enabled for this service"}
			c.JSON(statusCode, response)
			return
		}

		initiator, getUserError := userModels.Username(mInfo.OwnerUsername).GetFullUser(gc.DB, gc)

		if getUserError != nil {

			var ex tErrors.GenericError
			var ok bool

			ex, ok = getUserError.(tErrors.GenericError)
			if ok {
				c.JSON(ex.HTTPCode(), ex.JSONError())
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": getUserError.Error(), "message": getUserError.Error()})
			}
			return
		}
		documentIDStr := c.Param("documentID")
		documentID, _ := strconv.ParseUint(documentIDStr, 10, 64)

		t, _, _ := userModels.Username(initiator.Username).GetOpenTokenizedAssetByInitiatorUsername(gc.DB)

		if t.InitiatorUsername != initiator.Username {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Tokenized Asset not valid", "message": "Tokenized Asset not valid for you."})
			return
		}

		if len(t.ID) < 5 {
			c.JSON(http.StatusBadRequest, gin.H{"error": initiator.Username + " does not have a valid tokenized asset."})
			return
		}
		if t.AssetTokenizationStatus > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Tokenization Request cannot be altered at this stage through this option. Please use the option withint the tokenization detail."})
			return
		}

		document, err := userServices.DeleteTokenizationDocument(&initiator, documentID, gc)

		if err != nil {
			var ex tErrors.GenericError
			var ok bool

			ex, ok = err.(tErrors.GenericError)
			if ok {
				c.JSON(http.StatusBadRequest, ex.JSONError())
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": err.Error()})
			}
			return
		}

		if initiator.PushNotificationToken != nil {
			dataPayload := make(map[string]string)
			dataPayload["route"] = ""
			pns.SendFirebaseMessage(*initiator.PushNotificationToken, document.DocumentTitle+" deleted!", fmt.Sprintf("You have successfully deleted %v[%v].", document.DocumentTitle, document.DocumentType), "", dataPayload, gc.PushNotificationClient, gc.PNSContext)
		}

		userCacheKey := fmt.Sprintf("[GET] /v1/users/%v", initiator.Username)

		gc.RedisCache.InvalidateCachedHttpResponse(userCacheKey)

		//At this point, there was no error.

		c.JSON(http.StatusOK, document)
	}
}

// deleteTrovoApiAssetsFeesDocumentsDocumentIDHandler godoc
// @Summary DELETE /v1/trovo-api/assets/fees/documents/:documentID
// @Tags servicelinks
// @Produce json
// @Param documentID path string true "Document ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /v1/trovo-api/assets/fees/documents/{documentID} [delete]
func deleteTrovoApiAssetsFeesDocumentsDocumentIDHandler(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {

		mInfo, err := servicelinkServices.GetServiceLinkByAPIKey(middleware.ExtractServiceLinkApiKey(c), gc.DB)

		if err != nil {
			log.Println("[GET service] error for service:", middleware.ExtractServiceLinkApiKey(c), "error: ", err)

			var ex tErrors.GenericError
			var ok bool

			ex, ok = err.(tErrors.GenericError)
			var statusCode int = 0
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

		// ownerUsername := mInfo.OwnerUsername

		if mInfo.CreateUsersPermission == 0 {
			//wrong access
			statusCode := http.StatusUnauthorized
			response := gin.H{"error": "error-invalid-service-access", "data": "Permission", "message": "Permission to create users not enabled for this service"}
			c.JSON(statusCode, response)
			return
		}

		initiator, getUserError := userModels.Username(mInfo.OwnerUsername).GetFullUser(gc.DB, gc)

		if getUserError != nil {

			var ex tErrors.GenericError
			var ok bool

			ex, ok = getUserError.(tErrors.GenericError)
			if ok {
				c.JSON(ex.HTTPCode(), ex.JSONError())
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": getUserError.Error(), "message": getUserError.Error()})
			}
			return
		}
		documentIDStr := c.Param("documentID")
		documentID, _ := strconv.ParseUint(documentIDStr, 10, 64)

		t, _, _ := userModels.Username(initiator.Username).GetOpenTokenizedAssetByInitiatorUsername(gc.DB)

		if t.InitiatorUsername != initiator.Username {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Tokenized Asset not valid", "message": "Tokenized Asset not valid for you."})
			return
		}

		if len(t.ID) < 5 {
			c.JSON(http.StatusBadRequest, gin.H{"error": initiator.Username + " does not have a valid tokenized asset."})
			return
		}
		if t.AssetTokenizationStatus != 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Tokenization Request cannot be altered at this stage through this option. Please use the option within the tokenization detail."})
			return
		}

		document, err := userServices.DeleteTokenizationFeePaymentDocument(&initiator, documentID, gc)

		if err != nil {
			var ex tErrors.GenericError
			var ok bool

			ex, ok = err.(tErrors.GenericError)
			if ok {
				c.JSON(http.StatusBadRequest, ex.JSONError())
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": err.Error()})
			}
			return
		}

		if initiator.PushNotificationToken != nil {
			dataPayload := make(map[string]string)
			dataPayload["route"] = ""
			pns.SendFirebaseMessage(*initiator.PushNotificationToken, "Proof of payment deleted!", fmt.Sprintf("You have successfully deleted document with ID [%v].", documentID), "", dataPayload, gc.PushNotificationClient, gc.PNSContext)
		}

		userCacheKey := fmt.Sprintf("[GET] /v1/users/%v", initiator.Username)

		gc.RedisCache.InvalidateCachedHttpResponse(userCacheKey)

		//At this point, there was no error.

		c.JSON(http.StatusOK, document)
	}
}

// postTrovoApiAssetsConfirmApplicationTokenizationIDHandler godoc
// @Summary POST /v1/trovo-api/assets/confirm-application/:tokenizationID
// @Tags servicelinks
// @Accept json
// @Produce json
// @Param tokenizationID path string true "Tokenization ID"
// @Param body body userModels.ConfirmTokenizedAssetJSONInput true "Confirm tokenization payload"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /v1/trovo-api/assets/confirm-application/{tokenizationID} [post]
func postTrovoApiAssetsConfirmApplicationTokenizationIDHandler(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {

		mInfo, err := servicelinkServices.GetServiceLinkByAPIKey(middleware.ExtractServiceLinkApiKey(c), gc.DB)

		if err != nil {
			log.Println("[GET service] error for service:", middleware.ExtractServiceLinkApiKey(c), "error: ", err)

			var ex tErrors.GenericError
			var ok bool

			ex, ok = err.(tErrors.GenericError)
			var statusCode int = 0
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

		if mInfo.CreateUsersPermission == 0 {
			//wrong access
			statusCode := http.StatusUnauthorized
			response := gin.H{"error": "error-invalid-service-access", "data": "Permission", "message": "Permission to create users not enabled for this service"}
			c.JSON(statusCode, response)
			return
		}

		initiator, getUserError := userModels.Username(mInfo.OwnerUsername).GetFullUser(gc.DB, gc)

		if getUserError != nil {

			var ex tErrors.GenericError
			var ok bool

			ex, ok = getUserError.(tErrors.GenericError)
			if ok {
				c.JSON(ex.HTTPCode(), ex.JSONError())
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": getUserError.Error(), "message": getUserError.Error()})
			}
			return
		}

		var confirmationInput userModels.ConfirmTokenizedAssetJSONInput

		data, _ := io.ReadAll(c.Request.Body)
		// log.Println(string(data))
		err = json.Unmarshal(data, &confirmationInput)

		var invalidJSON tErrors.ErrorInvalidJSON

		if err != nil {
			c.JSON(http.StatusBadRequest, invalidJSON.JSONError())
			return
		}
		//confirm request
		_, err = userServices.ConfirmTokenizationApplicationInfoByInitiator(&initiator, c.Param("tokenizationID"), &confirmationInput, gc)

		if err != nil {

			var ex tErrors.GenericError
			var ok bool

			ex, ok = err.(tErrors.GenericError)
			if ok {
				c.JSON(ex.HTTPCode(), ex.JSONError())
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": err.Error()})
			}
			return
		}

		c.JSON(http.StatusOK, confirmationInput)
	}
}

// postTrovoApiAssetsMarketplacePrimaryHandler godoc
// @Summary POST /v1/trovo-api/assets/marketplace/primary
// @Tags servicelinks
// @Accept json
// @Produce json
// @Param body body userModels.TokenizedAssetPrimarySalesPurchaseInputForServiceLink true "Primary sales purchase payload"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /v1/trovo-api/assets/marketplace/primary [post]
func postTrovoApiAssetsMarketplacePrimaryHandler(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {

		mInfo, err := servicelinkServices.GetServiceLinkByAPIKey(middleware.ExtractServiceLinkApiKey(c), gc.DB)

		if err != nil {
			log.Println("[GET service] error for service:", middleware.ExtractServiceLinkApiKey(c), "error: ", err)

			var ex tErrors.GenericError
			var ok bool

			ex, ok = err.(tErrors.GenericError)
			var statusCode int = 0
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

		if mInfo.CreateUsersPermission == 0 {
			//wrong access
			statusCode := http.StatusUnauthorized
			response := gin.H{"error": "error-invalid-service-access", "data": "Permission", "message": "Permission to create users not enabled for this service"}
			c.JSON(statusCode, response)
			return
		}

		var tInput userModels.TokenizedAssetPrimarySalesPurchaseInputForServiceLink

		data, _ := io.ReadAll(c.Request.Body)
		// log.Println(string(data))
		err = json.Unmarshal(data, &tInput)

		var invalidJSON tErrors.ErrorInvalidJSON

		if err != nil {
			c.JSON(http.StatusBadRequest, invalidJSON.JSONError())
			return
		}
		user, err := usersDB.GetUser(tInput.PurchaserUsername, gc.DB, gc)

		if err != nil {
			log.Println("[GET USERINFO] error for user:", tInput.PurchaserUsername, "error: ", err)

			var ex tErrors.GenericError
			var ok bool

			ex, ok = err.(tErrors.GenericError)
			var statusCode int = 0
			var response interface{}

			if ok {
				statusCode = ex.HTTPCode()
				response = ex.JSONError()
			} else {
				statusCode = http.StatusBadRequest
				response = gin.H{"error": err.Error(), "message": err.Error()}
			}

			c.JSON(statusCode, response)
			return
		}
		//get the wallet you are sending payment from
		destinationWallet, temp, getWalletError := usersDB.GetWallet(tInput.DestinationWalletPublicKey, gc.DB)

		if getWalletError != nil {

			var ex tErrors.GenericError
			var ok bool

			ex, ok = getWalletError.(tErrors.GenericError)
			if ok {
				c.JSON(ex.HTTPCode(), ex.JSONError())
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": getWalletError.Error(), "message": getWalletError.Error()})
			}
			return
		}

		if temp {
			errAccountIsTemp := &tErrors.CustomError{
				Param:      "Username",
				Err:        "error-account-not-temporary-wallet",
				ErrMessage: "Only normal/standard wallets are allowed for this request.",
				Code:       http.StatusForbidden,
			}

			c.JSON(errAccountIsTemp.HTTPCode(), errAccountIsTemp.JSONError())
			return

		}

		tokenizedAsset, _, err := userServices.GetTokenizedAssetByID(tInput.TokenizedAssetID, gc.DB)
		if err != nil {

			var ex tErrors.GenericError
			var ok bool

			ex, ok = err.(tErrors.GenericError)
			if ok {
				c.JSON(ex.HTTPCode(), ex.JSONError())
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": err.Error()})
			}
			return
		}
		subscriptionInput := tInput.ToSubscriptionInput(gc)
		sub, err := userServices.SubscribeToTokenizedAsset(&user, &destinationWallet, &tokenizedAsset, &subscriptionInput, gc)
		if err != nil {

			var ex tErrors.GenericError
			var ok bool

			ex, ok = err.(tErrors.GenericError)
			if ok {
				c.JSON(ex.HTTPCode(), ex.JSONError())
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": err.Error()})
			}
			return
		}
		tInput = subscriptionInput.ToServiceLinkInput(gc)
		c.JSON(http.StatusOK, tInput)

		if user.PushNotificationToken != nil && len(tInput.TransactionID) > 0 && tInput.TransactionID != "PENDING_AUTH" {
			dataPayload := make(map[string]string)
			dataPayload["route"] = "assetSubscription"
			user.SendPushMessage(fmt.Sprintf("You have successfully subscribed to %v", *tokenizedAsset.AssetCode), fmt.Sprintf("You have successfully purchased %v %v worth of %v on the wallet with alias [%v].", sub.Amount, *tokenizedAsset.AssetQuoteCurrency, *tokenizedAsset.AssetCode, destinationWallet.Alias), "", dataPayload, gc)
		}
		user.InvalidateUserCache(gc)
	}
}
