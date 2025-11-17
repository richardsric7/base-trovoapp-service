package servicelinks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"
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
	"gorm.io/gorm/clause"

	"github.com/gin-gonic/gin"
)

// Init initializes /v1/services endpoint
func Init(router *gin.Engine, gc *sharedconfig.GlobalConfig) {

	//retryCallbacks stores failed callbacks
	type retryCallbacks struct {
		Req         *bytes.Buffer
		CallbackURL string
		Count       int
	}

	callBackRetryChan := make(chan retryCallbacks, 20000)
	go func(c chan retryCallbacks) {
		log.Println("#####@started Routine to retry failed auth/events/login callbacks....")
		//loop
		for {
			callbackObj := <-c
			if callbackObj.Count > 100 {
				//skip 100 retries
				continue
			}
			_, err := http.Post(callbackObj.CallbackURL, "application/json", callbackObj.Req)

			if err != nil {
				//send back into channel to retry later
				log.Printf("Callback retry failed: [%+v]\n", callbackObj)
				if callbackObj.Count <= 99 {
					callbackObj.Count++
					c <- callbackObj
				}

			}
			//wait 1 second
			time.Sleep(1 * time.Second)
		}

	}(callBackRetryChan)

	{
		//auto expire login sessions that where that are not within valid time.

		go func() {
			log.Println("@@@@@Started routine to Auto remove <SERVICELINK> LoginSessions")
			period := time.Duration(3)
			if os.Getenv("SERVICE_LINK_LOGIN_REQUEST_VALIDITY") != "" {
				m, e := decimal.NewFromString(os.Getenv("SERVICE_LINK_LOGIN_REQUEST_VALIDITY"))
				if e == nil {
					if m.IsPositive() {
						period = time.Duration(m.IntPart())
					}
				}
			}
			for {
				err := gc.DB.Where("created_at < ?", time.Now().Add(-1*period*time.Minute)).Delete(servicelinkModels.ServiceLinkLoginSession{}).Error
				if err != nil {
					log.Printf("[Expire Login Sessions Routine]unable to delete expired login sessions due to error [%v]\n", err)
				}
				time.Sleep(30 * time.Second)
			}
		}()
	}

	{
		//auto expire authorizations that are not within valid time.

		go func() {
			log.Println("@@@@@Started routine to Auto remove <service> Authorizations")

			for {
				err := gc.DB.Where("expires_at < ?", time.Now()).Delete(servicelinkModels.ServiceLinkAuthorization{}).Error
				if err != nil {
					log.Printf("[Expire Authorizations Routine]unable to delete expired authorizations requests due to error [%v]\n", err)
				}
				time.Sleep(30 * time.Second)
			}
		}()
	}

	{
		//auto expire events that are not within valid time.

		go func() {
			log.Println("@@@@@Started routine to Auto remove <service> events")

			for {
				err := gc.DB.Where("expires_at < ?", time.Now()).Delete(servicelinkModels.ServiceLinkEvent{}).Error
				if err != nil {
					log.Printf("[Expire Events Routine]unable to delete expired events requests due to error [%v]\n", err)
				}
				time.Sleep(30 * time.Second)
			}
		}()
	}

	//service login request
	router.POST("/v1/servicelinks/login/request/:targetUser", middleware.AuthenticationMiddlewareUsingAPIKey(gc), func(c *gin.Context) {
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

	})

	//user login approval url; uses signature algorithm bcos it is only called by trovoApp.
	router.POST("/v1/users/servicelinks/login/approval/:targetUser", middleware.AuthenticationMiddlewareUsingTimestamp(), func(c *gin.Context) {

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

	})

	//service login verify url
	router.GET("/v1/servicelinks/login/verify/:ownerUsername/:targetUser/:loginID", middleware.AuthenticationMiddlewareUsingAPIKey(gc), func(c *gin.Context) {

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
	})

	//service login refresh token url
	router.POST("/v1/servicelinks/token/refresh", func(c *gin.Context) {

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

	})

	//service token verify
	router.POST("/v1/servicelinks/token/verify", func(c *gin.Context) {
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

	})

	//service token verify
	router.DELETE("/v1/servicelinks/token", func(c *gin.Context) {
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

	})

	//service authorization request
	router.POST("/v1/servicelinks/authorize/request/:targetUser", middleware.AuthenticationMiddlewareUsingAPIKey(gc), func(c *gin.Context) {
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
	})

	//service authorization tokenizedAsset
	router.POST("/v1/servicelinks/authorize/tokenized-asset", middleware.AuthenticationMiddlewareUsingAPIKey(gc), func(c *gin.Context) {

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
	})

	//service event link request
	router.POST("/v1/servicelinks/events/request", middleware.AuthenticationMiddlewareUsingAPIKey(gc), func(c *gin.Context) {

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
	})

	//user authorization approval url
	router.POST("/v1/users/servicelinks/authorize/approval/:targetUser", middleware.AuthenticationMiddlewareUsingTimestamp(), func(c *gin.Context) {

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

	})

	//user events approval url
	router.POST("/v1/users/servicelinks/events/approval/:targetUser", middleware.AuthenticationMiddlewareUsingTimestamp(), func(c *gin.Context) {

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

	})

	//service authorization verify url
	router.GET("/v1/servicelinks/authorize/verify/:targetUser", middleware.AuthenticationMiddlewareUsingTimestamp(), func(c *gin.Context) {
		// var err error

		ownerUsername := strings.TrimSpace(strings.ToLower(c.Query("ownerUsername")))
		trovoUser := strings.TrimSpace(strings.ToLower(c.Param("targetUser")))
		authID := strings.TrimSpace(strings.ToLower(c.Query("authId")))

		conDB.PrintDBStats(fmt.Sprintf("GET /v1/servicelinks/authorize/verify/%v %v/%v", trovoUser, ownerUsername, authID), gc.DB)
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

			//could not save login session
			response := gin.H{"error": "error-unauthorized-request", "data": "unauthorizedRequest", "message": "authorization awaiting approval"}
			statusCode := http.StatusUnauthorized
			c.JSON(statusCode, response)
			return

		}
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	//service payment request
	router.GET("/v1/servicelinks/payment/request/:targetUser", middleware.AuthenticationMiddlewareUsingAPIKey(gc), func(c *gin.Context) {

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

	})

	//service payment request
	router.GET("/v1/servicelinks/tokenized-asset/:assetCode", middleware.AuthenticationMiddlewareUsingAPIKey(gc), func(c *gin.Context) {

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

	})

	//SERVICELINK USER INFO request
	router.GET("/v1/servicelinks/:ownerUsername/:targetUser/userinfo", middleware.AuthenticationMiddlewareUsingAPIKey(gc), func(c *gin.Context) {

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

	})

	//SERVICE push notification request
	router.POST("/v1/servicelinks/:ownerUsername/:targetUser/push", middleware.AuthenticationMiddlewareUsingAPIKey(gc), func(c *gin.Context) {

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
	})

	//register user from service link
	router.POST("/v1/trovo-api/users/onboard", middleware.AuthenticationMiddlewareUsingAPIKey(gc), func(c *gin.Context) {

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
	})

	//update user kyc from service link
	router.POST("/v1/trovo-api/users/update-kyc", middleware.AuthenticationMiddlewareUsingAPIKey(gc), func(c *gin.Context) {

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
	})

	//mint token from service link
	router.POST("/v1/trovo-api/tokens/mint", middleware.AuthenticationMiddlewareUsingAPIKey(gc), func(c *gin.Context) {

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

		c.JSON(http.StatusOK, data)
	})

	//create new subwallet from service link
	router.POST("/v1/trovo-api/users/subwallet", middleware.AuthenticationMiddlewareUsingAPIKey(gc), func(c *gin.Context) {

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

		var subWalletInfo userModels.SubWalletInfo
		// var err error

		data, _ := io.ReadAll(c.Request.Body)

		err = json.Unmarshal(data, &subWalletInfo)

		var invalidJSON tErrors.ErrorInvalidJSON

		if err != nil {
			c.JSON(http.StatusBadRequest, invalidJSON.JSONError())
			return
		}

		conDB.PrintDBStats(fmt.Sprintf("POST /v1/trovo-api/users/subwallet %v", signerUser.Username), gc.DB)

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

		c.JSON(http.StatusOK, returnedSubwalletInfo)
	})

	//Get tokenization parameters from service link
	router.GET("/v1/trovo-api/assets/parameters", middleware.AuthenticationMiddlewareUsingAPIKey(gc), func(c *gin.Context) {

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
		countries := userServices.GetCountryConfigs(gc.DB)

		c.JSON(http.StatusOK, gin.H{"assetSectors": sectorList, "assetSubSectors": subsectorList, "assetTypes": assetTypes, "assetCustodians": custdians, "assetManagers": managers,
			"assetIssuingHouses": issuers, "tokenizationFees": fees, "tokenizationCurrencies": currencies, "assetProtectionOptions": apo, "assetProceedCycle": apc,
			"publicListingAllowedCountries": ac, "tokenizationDocumentTypes": docTypes,
			"tokenizationStatuses": statuses, "feePaymentMethods": fpms, "countryConfigs": countries})

	})

	//Get tokenization bank list from service link
	router.GET("/v1/trovo-api/assets/bank-list/:countryCode", middleware.AuthenticationMiddlewareUsingAPIKey(gc), func(c *gin.Context) {

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

	})

	//Get tokenization list for Admin from service link
	router.GET("/v1/trovo-api/assets/admin/list", middleware.AuthenticationMiddlewareUsingAPIKey(gc), func(c *gin.Context) {

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

	})

	//Get tokenization list for marketvfrom service link
	router.GET("/v1/trovo-api/assets/market/list", middleware.AuthenticationMiddlewareUsingAPIKey(gc), func(c *gin.Context) {

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

	})

	//Apply for tokenization from service link
	router.POST("/v1/trovo-api/assets/apply", middleware.AuthenticationMiddlewareUsingAPIKey(gc), func(c *gin.Context) {

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
	})

	//Upload logo for tokenization from service link
	router.PUT("/v1/trovo-api/assets/logo", middleware.AuthenticationMiddlewareUsingAPIKey(gc), func(c *gin.Context) {

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
	})

	//Upload documents for tokenization from service link
	router.PUT("/v1/trovo-api/assets/documents", middleware.AuthenticationMiddlewareUsingAPIKey(gc), func(c *gin.Context) {

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
	})

	//Upload fee payment documents for tokenization from service link
	router.PUT("/v1/trovo-api/assets/fees/document", middleware.AuthenticationMiddlewareUsingAPIKey(gc), func(c *gin.Context) {

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

	})

	//confirm fee payment for tokenization from service link
	router.POST("/v1/trovo-api/assets/fees/confirm/:tokenizationID", middleware.AuthenticationMiddlewareUsingAPIKey(gc), func(c *gin.Context) {

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

	})

	//DELETE tokenization from service link
	router.DELETE("/v1/trovo-api/assets/:tokenizationID", middleware.AuthenticationMiddlewareUsingAPIKey(gc), func(c *gin.Context) {

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
	})

	//DELETE tokenization document from service link
	router.DELETE("/v1/trovo-api/assets/documents/:documentID", middleware.AuthenticationMiddlewareUsingAPIKey(gc), func(c *gin.Context) {

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
	})

	//DELETE tokenization fee document from service link
	router.DELETE("/v1/trovo-api/assets/fees/documents/:documentID", middleware.AuthenticationMiddlewareUsingAPIKey(gc), func(c *gin.Context) {

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
	})

	//Confirm Application for tokenization from service link
	router.POST("/v1/trovo-api/assets/confirm-application/:tokenizationID", middleware.AuthenticationMiddlewareUsingAPIKey(gc), func(c *gin.Context) {

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
	})

	//Purchase tokenized asset from service link
	router.POST("/v1/trovo-api/assets/buy/:tokenizationID", middleware.AuthenticationMiddlewareUsingAPIKey(gc), func(c *gin.Context) {

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

		tokenizedAssetID := c.Param("tokenizedAssetID")
		user, err := usersDB.GetUserFromPrimarySigner(middleware.ExtractSigner(c), gc.DB, gc)

		if err != nil {
			log.Println("[GET USERINFO] error for user:", middleware.ExtractSigner(c), "error: ", err)

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
		subscriberWallet, temp, getWalletError := usersDB.GetWallet(middleware.ExtractPublicKey(c), gc.DB)

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

		tokenizedAsset, _, err := userServices.GetTokenizedAssetByID(tokenizedAssetID, gc.DB)
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

		var tInput userModels.TokenizedAssetSubscriptionInput

		data, _ := io.ReadAll(c.Request.Body)
		// log.Println(string(data))
		err = json.Unmarshal(data, &tInput)

		var invalidJSON tErrors.ErrorInvalidJSON

		if err != nil {
			c.JSON(http.StatusBadRequest, invalidJSON.JSONError())
			return
		}

		sub, err := userServices.SubscribeToTokenizedAsset(&user, &subscriberWallet, &tokenizedAsset, &tInput, gc)
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

		c.JSON(http.StatusOK, tInput)

		if user.PushNotificationToken != nil && len(tInput.TransactionID) > 0 && tInput.TransactionID != "PENDING_AUTH" {
			dataPayload := make(map[string]string)
			dataPayload["route"] = "assetSubscription"
			user.SendPushMessage(fmt.Sprintf("You have successfully subscribed to %v", *tokenizedAsset.AssetCode), fmt.Sprintf("You have successfully purchased %v %v worth of %v on the wallet with alias [%v].", sub.Amount, *tokenizedAsset.AssetQuoteCurrency, *tokenizedAsset.AssetCode, subscriberWallet.Alias), "", dataPayload, gc)
		}
		user.InvalidateUserCache(gc)
	})

}
