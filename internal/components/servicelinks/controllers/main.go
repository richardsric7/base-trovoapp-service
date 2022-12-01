package servicelinks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
	servicelinkModels "trovo-wallet-api/internal/components/servicelinks/models"
	servicelinkServices "trovo-wallet-api/internal/components/servicelinks/services"
	userServices "trovo-wallet-api/internal/components/users/services"

	conDB "trovo-wallet-api/internal/db"
	dl "trovo-wallet-api/internal/dynamiclinks"
	tErrors "trovo-wallet-api/internal/errors"
	pns "trovo-wallet-api/internal/pns"
	"trovo-wallet-api/internal/sharedconfig"

	"log"
	"net/http"
	"strings"
	"trovo-wallet-api/internal/middleware"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

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

		err = gc.DB.Create(&newLoginSession).Error
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
		// 	pns.SendFirebaseMessage(*userInfo.PushNotificationToken, fmt.Sprintf("Login for [%v] requested!", userInfo.Username), fmt.Sprintf("Your Trovo Wallet username [%v] has been used to request a login session on [%v] service using [%v]. Click to continue.", userInfo.Username, mInfo.LongName, serviceLinkRequestInput.DeviceInfo), "", dataPayload, gc.PushNotificationClient, gc.PNSContext)
		// }

		c.JSON(http.StatusOK, data)

	})

	//user login approval url; uses signature algorithm bcos it is only called by trovo wallet.
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
			response := gin.H{"error": "error-invalid-user-access", "data": "Authentication", "message": "this login request does not belong to your Trovo wallet"}
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

		err = gc.DB.Save(&loginSession).Error
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
			pns.SendFirebaseMessage(*userInfo.PushNotificationToken, fmt.Sprintf("Login for %v authorized!", userInfo.Username), fmt.Sprintf("Your Trovo Wallet username %v has been authorized to login on %v service.", userInfo.Username, mInfo.LongName), "", dataPayload, gc.PushNotificationClient, gc.PNSContext)
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
		cacheKey = fmt.Sprintf("[GET] /v1/servicelinks/%v/%v/login/%v", mInfo.OwnerUsername, userInfo.Username, loginID)
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
				log.Printf("[%v], served from cache\n", cacheKey)
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
			cacheDurationInSeconds := 60 //1 minutes

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
		mapToken := map[string]string{}
		if err := c.ShouldBindJSON(&mapToken); err != nil {
			log.Printf("[HandleRefreshToken] error could not bind json body: %v\n", err)
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "unable to retrieve refresh token"})
			return
		}
		refreshToken := mapToken["refreshToken"]

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
		if token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token is invalid"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "success", "token": token})

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
			log.Println("Login Authorization Request Input JSON Error:", err)
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
		err = gc.DB.Create(&authorizationData).Error
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
		// 	pns.SendFirebaseMessage(*userInfo.PushNotificationToken, fmt.Sprintf("Authorization for trovo account %v requested!", userInfo.Username), fmt.Sprintf("Your Trovo Wallet username [%v] has been used to request an authorization session on [%v] service using [%v]. Click to continue.", userInfo.Username, mInfo.LongName, serviceLinkRequestInput.DeviceInfo), "", dataPayload, gc.PushNotificationClient, gc.PNSContext)
		// }
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
				response = gin.H{"error": err.Error()}
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
		err = gc.DB.Create(&authorizationData).Error
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
			response := gin.H{"error": "error-invalid-user-access", "data": "Authentication", "message": "2FA/Authorization request does not belong to your Trovo Wallet"}
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

		err = gc.DB.Save(&authData).Error
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
			pns.SendFirebaseMessage(*userInfo.PushNotificationToken, fmt.Sprintf("2FA Action for trovo account %v authorized!", userInfo.Username), fmt.Sprintf("Your Trovo Wallet username %v has been used to authorize a 2FA action on %v service.", userInfo.Username, mInfo.LongName), "", dataPayload, gc.PushNotificationClient, gc.PNSContext)
		}
		//return report to user and not keep them waiting.
		c.JSON(http.StatusOK, gin.H{"message": "success"})

		//TODO: make callback request if callback is availble
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
			pns.SendFirebaseMessage(*userInfo.PushNotificationToken, fmt.Sprintf("Event Registration/Participation for %v authorized!", userInfo.Username), fmt.Sprintf("Your Trovo Wallet username %v has been used to authorize an event registration/participation action on %v service.", userInfo.Username, mInfo.LongName), "", dataPayload, gc.PushNotificationClient, gc.PNSContext)
		}

		//return report to user and not keep them waiting.
		c.JSON(http.StatusOK, gin.H{"message": "success"})

		//TODO: make callback request if callback is available
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
				log.Printf("[%v]/[%v], served from cache\n", cacheKey, cacheKeyParameters)
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

}
