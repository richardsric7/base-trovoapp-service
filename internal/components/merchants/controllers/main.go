package merchants

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"time"
	"trovo-wallet-api/internal/cache"
	merchantModels "trovo-wallet-api/internal/components/merchants/models"
	merchantServices "trovo-wallet-api/internal/components/merchants/services"
	users "trovo-wallet-api/internal/components/users/services"
	conDB "trovo-wallet-api/internal/db"

	"fmt"
	tErrors "trovo-wallet-api/internal/errors"

	"log"
	"net/http"
	"strings"
	"trovo-wallet-api/internal/middleware"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Init initializes /v2/merchants endpoint
func Init(router *gin.Engine, db *gorm.DB, redisCache *cache.RedisCache, dynamicLinkServiceUrlChan chan string) {

	//retryCallbacks stores failed callbacks
	type retryCallbacks struct {
		Req         *bytes.Buffer
		CallbackURL string
		Count       int
	}

	callBackRetryChan := make(chan retryCallbacks, 20000)
	go func(c chan retryCallbacks) {
		log.Println("#####@started Routine to retry failed Payment callbacks....")
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
			log.Println("@@@@@Started routine to Auto remove <merchant> LoginSessions")
			period := time.Duration(3)
			if os.Getenv("MERCHANT_LOGIN_REQUEST_VALIDITY") != "" {
				m, e := decimal.NewFromString(os.Getenv("MERCHANT_LOGIN_REQUEST_VALIDITY"))
				if e == nil {
					if m.IsPositive() {
						period = time.Duration(m.IntPart())
					}
				}
			}
			for {
				err := db.Where("created_at < ?", time.Now().Add(-1*period*time.Minute)).Delete(merchantModels.MerchantLoginSession{}).Error
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
			log.Println("@@@@@Started routine to Auto remove <merchant> Authorizations")
			// period := time.Duration(3)
			// if os.Getenv("MERCHANT_AUTHORIZATION_REQUEST_VALIDITY") != "" {
			// 	m, e := decimal.NewFromString(os.Getenv("MERCHANT_AUTHORIZATION_REQUEST_VALIDITY"))
			// 	if e == nil {
			// 		if m.IsPositive() {
			// 			period = time.Duration(m.IntPart())
			// 		}
			// 	}
			// }
			for {
				err := db.Where("expires_at < ?", time.Now()).Delete(merchantModels.MerchantAuthorization{}).Error
				if err != nil {
					log.Printf("[Expire Authorizations Routine]unable to delete expired authorizations requests due to error [%v]\n", err)
				}
				time.Sleep(30 * time.Second)
			}
		}()
	}

	//merchant login request
	router.POST("/v2/merchants/:merchantID/:targetUser/login", middleware.AuthenticationMiddleware(), func(c *gin.Context) {

		identifier := strings.TrimSpace(strings.ToLower(c.Param("merchantID")))
		bantupayUser := strings.TrimSpace(strings.ToLower(c.Param("targetUser")))
		if bantupayUser == "null" {
			log.Printf("user cannot be %v\n", bantupayUser)
			c.JSON(http.StatusBadRequest, gin.H{"error": "user cannot be null"})
			return
		}
		conDB.PrintDBStats(fmt.Sprintf("POST /v2/merchants/%v/%v/login?", identifier, bantupayUser), db)
		mInfo, err := merchantServices.GetMerchant(identifier, db)

		if err != nil {
			log.Println("[GET MERCHANT] error for merchant:", identifier, "error: ", err)

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
			response := gin.H{"error": "error-invalid-merchant-access", "data": "Authentication", "message": "Authentication failed"}
			c.JSON(statusCode, response)
			return
		}
		// log.Printf("Merchant Infor: %+v\n", mInfo)
		if mInfo.LoginPermission == 0 {
			//wrong access
			statusCode := http.StatusUnauthorized
			response := gin.H{"error": "error-invalid-merchant-access", "data": "Permission", "message": "login permission not enabled for this merchant"}
			c.JSON(statusCode, response)
			return
		}

		userInfo, err := users.GetUserForMerchants(bantupayUser, mInfo, db, dynamicLinkServiceUrlChan, redisCache)

		if err != nil {
			log.Println("[GET UserInfo] error for user:", bantupayUser, "error: ", err)

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

		var merchantRequestInput merchantModels.MerchantRequestInput
		reqBody, _ := ioutil.ReadAll(c.Request.Body)

		err = json.Unmarshal(reqBody, &merchantRequestInput)

		var invalidJSON tErrors.ErrorInvalidJSON

		if err != nil {
			log.Println("Login Request Input JSON Error:", err)
			c.JSON(http.StatusBadRequest, invalidJSON.JSONError())
			return
		}

		//store login data for verification
		//does not exist. create new one.
		loginID := uuid.NewString()
		newLoginSession := merchantModels.MerchantLoginSession{
			MerchantUsername: identifier,
			WalletUsername:   userInfo.Username,
			ID:               loginID,
		}

		if len(merchantRequestInput.CallbackURL) > 0 {
			newLoginSession.CallbackURL = &merchantRequestInput.CallbackURL
		}

		err = db.Create(&newLoginSession).Error
		if err != nil {
			//could not create login session
			response := gin.H{"error": "error-temporary-server-error", "data": "temporaryServerError", "message": "Temporary Server Error. Contact support."}
			statusCode := http.StatusServiceUnavailable
			c.JSON(statusCode, response)
			return
		}

		//respond with deeplink and QRCode for login
		data, err := merchantServices.GenerateLoginData(mInfo.BantupayUsername, mInfo.ShortName, userInfo.Username, loginID, merchantRequestInput.DeviceInfo, dynamicLinkServiceUrlChan, redisCache)
		if err != nil {
			//could not create login session
			response := gin.H{"error": "error-temporary-server-error", "data": "temporaryServerError", "message": "Temporary Server Error. Contact support."}
			statusCode := http.StatusServiceUnavailable
			c.JSON(statusCode, response)
			return
		}
		c.JSON(http.StatusOK, data)

	})

	//user login approval url
	router.POST("/v2/users/merchants/:targetUser/login/:merchantID/:loginID", middleware.AuthenticationMiddleware(), func(c *gin.Context) {

		identifier := strings.TrimSpace(strings.ToLower(c.Param("targetUser")))
		merchant := strings.TrimSpace(strings.ToLower(c.Param("merchantID")))
		loginID := strings.TrimSpace(strings.ToLower(c.Param("loginID")))
		if identifier == "null" {
			log.Printf("user cannot be %v\n", identifier)
			c.JSON(http.StatusBadRequest, gin.H{"error": "user cannot be null"})
			return
		}
		conDB.PrintDBStats(fmt.Sprintf("POST /v2/users/%v/login/%v/%v", identifier, merchant, loginID), db)
		mInfo, err := merchantServices.GetMerchant(merchant, db)

		if err != nil {
			log.Println("[GET MERCHANT] error for merchant:", merchant, "error: ", err)

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
			response := gin.H{"error": "error-invalid-merchant-access", "data": "Permission", "message": "login permission not enabled for this merchant"}
			c.JSON(statusCode, response)
			return
		}

		userInfo, err := users.GetUserForMerchants(identifier, mInfo, db, dynamicLinkServiceUrlChan, redisCache)

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
		if len(strings.TrimSpace(userInfo.Mobile)) < 6 {
			log.Println("[GET UserInfo] error for user:", identifier, "error: user does not have a valid phone number")
			err = &tErrors.CustomError{
				Param:      "mobile",
				Err:        "no valid phone number",
				ErrMessage: "only users with valid phone number are allows to use this service. please update your mobile number.",
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
		if middleware.ExtractPublicKey(c) != userInfo.PublicKey {
			//wrong access
			statusCode := http.StatusUnauthorized
			response := gin.H{"error": "error-invalid-user-access", "data": "Authentication", "message": "this login request does not belong to your Bantupay wallet"}
			c.JSON(statusCode, response)
			return
		}
		//store login data for verification
		loginSession, err := merchantServices.GetLoginSession(merchant, userInfo.Username, loginID, db)
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
		if loginSession.Authorized == 1 {
			response := gin.H{"error": "error-login-session-does-not-exist", "data": userInfo.Username, "message": "invalid/expired login request"}
			statusCode := http.StatusNotFound
			c.JSON(statusCode, response)
			return
		}
		//exists and needs to be updated
		loginSession.Authorized = 1

		err = db.Save(&loginSession).Error
		if err != nil {
			//could not save login session
			response := gin.H{"error": "error-temporary-server-error", "data": "temporaryServerError", "message": "Temporary Server Error. Contact support."}
			statusCode := http.StatusServiceUnavailable
			c.JSON(statusCode, response)
			return
		}

		cacheKey := fmt.Sprintf("[GET] /v2/merchants/%v/%v/login/%v", mInfo.BantupayUsername, userInfo.Username, loginID)
		redisCache.InvalidateCachedHttpResponse(cacheKey)

		//return repsonse to user and  not keep them waiting.
		c.JSON(http.StatusOK, gin.H{"message": "success"})

		//make callback request if callback is availble

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
			body, err = ioutil.ReadAll(resp.Body)
			if err != nil {
				//send to retry channel
				log.Println("[LoginAuthCallback] callback failed:", err)
				c := retryCallbacks{Req: responseBody, CallbackURL: d, Count: 0}
				callBackRetryChan <- c
			} else {
				log.Printf("[LoginAuthCallback] Login authorization Callback successful to: [%v], Response:[%v]\n\n", d, string(body))

			}
		}
		cacheKey = fmt.Sprintf("[GET] /v2/merchants/%v/%v/login/%v", mInfo.BantupayUsername, userInfo.Username, loginID)
		redisCache.InvalidateCachedHttpResponse(cacheKey)

	})

	//merchant login verify url
	router.GET("/v2/merchants/:merchantID/:targetUser/login/:loginID", middleware.AuthenticationMiddleware(), func(c *gin.Context) {

		identifier := strings.TrimSpace(strings.ToLower(c.Param("merchantID")))
		bantupayUser := strings.TrimSpace(strings.ToLower(c.Param("targetUser")))
		loginID := strings.TrimSpace(strings.ToLower(c.Param("loginID")))
		if bantupayUser == "null" {
			log.Printf("user cannot be %v\n", bantupayUser)
			c.JSON(http.StatusBadRequest, gin.H{"error": "user cannot be null"})
			return
		}
		conDB.PrintDBStats(fmt.Sprintf("GET /v2/merchants/%v/%v/login/%v", identifier, bantupayUser, loginID), db)

		cacheKey := fmt.Sprintf("[GET] /v2/merchants/%v/%v/login/%v", identifier, bantupayUser, loginID)
		{
			//search cache

			ok, status, response := redisCache.CachedHttpResponse(cacheKey)

			if ok {
				log.Printf("[%v], served from cache\n", cacheKey)
				c.JSON(status, response)
				return
			}
		}

		mInfo, err := merchantServices.GetMerchant(identifier, db)

		if err != nil {
			log.Println("[GET MERCHANT USER LOGIN] error for merchant:", identifier, "error: ", err)

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
			response := gin.H{"error": "error-invalid-merchant-access", "data": "Authentication", "message": "Authentication failed"}
			c.JSON(statusCode, response)
			cacheDurationInSeconds := 60 //1 minutes

			redisCache.CacheHttpResponse(cacheKey, statusCode, response, cacheDurationInSeconds)
			return
		}
		// log.Printf("Merchant Info: %+v\n", mInfo)
		if mInfo.LoginPermission == 0 {
			//wrong access
			statusCode := http.StatusUnauthorized
			response := gin.H{"error": "error-invalid-merchant-access", "data": "Permission", "message": "login permission not enabled for this merchant"}
			c.JSON(statusCode, response)
			cacheDurationInSeconds := 60 //1 minutes

			redisCache.CacheHttpResponse(cacheKey, statusCode, response, cacheDurationInSeconds)
			return
		}

		userInfo, err := users.GetUserForMerchants(bantupayUser, mInfo, db, dynamicLinkServiceUrlChan, redisCache)

		if err != nil {
			log.Println("[GET UserInfo] error for user:", bantupayUser, "error: ", err)

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

			redisCache.CacheHttpResponse(cacheKey, statusCode, response, cacheDurationInSeconds)
			return
		}

		//store login data for verification
		loginSession, err := merchantServices.GetLoginSession(identifier, userInfo.Username, loginID, db)
		if err != nil {
			log.Printf("[error Verifyig Login] for user [%v], error [%v]]\n", bantupayUser, err)
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

			redisCache.CacheHttpResponse(cacheKey, statusCode, response, cacheDurationInSeconds)

			c.JSON(statusCode, response)
			return

		}
		if loginSession.Authorized == 0 {

			//could not save login session
			response := gin.H{"error": "error-unauthorized-login-session", "data": "unauthorizedLoginSession", "message": "login session awaiting authorization"}
			statusCode := http.StatusUnauthorized

			cacheDurationInSeconds := 60 //1 minutes

			redisCache.CacheHttpResponse(cacheKey, statusCode, response, cacheDurationInSeconds)

			c.JSON(statusCode, response)
			return

		}

		cacheDurationInSeconds := 2 * 60 //2 minutes

		redisCache.CacheHttpResponse(cacheKey, http.StatusOK, userInfo, cacheDurationInSeconds)

		c.JSON(http.StatusOK, userInfo)
	})

	//merchant authorization request
	router.POST("/v2/merchants/:merchantID/:targetUser/authorize", middleware.AuthenticationMiddleware(), func(c *gin.Context) {

		identifier := strings.TrimSpace(strings.ToLower(c.Param("merchantID")))
		bantupayUser := strings.TrimSpace(strings.ToLower(c.Param("targetUser")))

		conDB.PrintDBStats(fmt.Sprintf("POST /v2/merchants/%v/%v/authorize?", identifier, bantupayUser), db)
		if bantupayUser == "null" {
			log.Printf("user cannot be %v\n", bantupayUser)
			c.JSON(http.StatusBadRequest, gin.H{"error": "user cannot be null"})
			return
		}
		mInfo, err := merchantServices.GetMerchant(identifier, db)

		if err != nil {
			log.Println("[GET MERCHANT] error for merchant:", identifier, "error: ", err)

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
			response := gin.H{"error": "error-invalid-merchant-access", "data": "Authentication", "message": "Authentication failed"}
			c.JSON(statusCode, response)
			return
		}
		// log.Printf("Merchant Infor: %+v\n", mInfo)
		if mInfo.AuthorizationPermission == 0 {
			//wrong access
			statusCode := http.StatusUnauthorized
			response := gin.H{"error": "error-invalid-merchant-access", "data": "Permission", "message": "authorization permission not enabled for this merchant"}
			c.JSON(statusCode, response)
			return
		}

		userInfo, err := users.GetUserForMerchants(bantupayUser, mInfo, db, dynamicLinkServiceUrlChan, redisCache)

		if err != nil {
			log.Println("[GET UserInfo] error for user:", bantupayUser, "error: ", err)

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

		var merchantRequestInput merchantModels.MerchantRequestInput
		reqBody, _ := ioutil.ReadAll(c.Request.Body)

		err = json.Unmarshal(reqBody, &merchantRequestInput)

		var invalidJSON tErrors.ErrorInvalidJSON

		if err != nil {
			log.Println("LoAithorization Request Input JSON Error:", err)
			c.JSON(http.StatusBadRequest, invalidJSON.JSONError())
			return
		}

		authID := uuid.NewString()
		authorizationData := merchantModels.MerchantAuthorization{ID: authID, MerchantUsername: mInfo.BantupayUsername,
			WalletUsername: userInfo.Username}
		if len(merchantRequestInput.CallbackURL) > 0 {
			authorizationData.CallbackURL = &merchantRequestInput.CallbackURL
		}
		period := time.Duration(3)
		if os.Getenv("MERCHANT_AUTHORIZATION_REQUEST_VALIDITY") != "" {
			m, e := decimal.NewFromString(os.Getenv("MERCHANT_AUTHORIZATION_REQUEST_VALIDITY"))
			if e == nil {
				if m.IsPositive() {
					period = time.Duration(m.IntPart())
				}
			}
		}
		if merchantRequestInput.ValidityInMinutes > 0 {
			//check if validtity was submitted
			period = time.Duration(merchantRequestInput.ValidityInMinutes)

		}
		authorizationData.ExpiresAt = time.Now().Add(period * time.Minute)
		err = db.Create(&authorizationData).Error
		if err != nil {
			log.Printf("unable to create authorization: %v\n", err)
			//could not save login session
			response := gin.H{"error": "error-temporary-server-error", "data": "temporaryServerError", "message": "Temporary Server Error. Contact support."}
			statusCode := http.StatusServiceUnavailable
			c.JSON(statusCode, response)
			return
		}

		data, err := merchantServices.GenerateAuthorizationData(mInfo.BantupayUsername, mInfo.ShortName, merchantRequestInput.AuthDescription, userInfo.Username, merchantRequestInput.DeviceInfo, authID, dynamicLinkServiceUrlChan, redisCache)
		if err != nil {
			//could not create authorization session
			response := gin.H{"error": "error-temporary-server-error", "data": "temporaryServerError", "message": "Temporary Server Error. Contact support."}
			statusCode := http.StatusServiceUnavailable
			c.JSON(statusCode, response)
			return
		}
		c.JSON(http.StatusOK, data)
	})

	//merchant push notification request
	router.POST("/v2/merchants/:merchantID/:targetUser/push", middleware.AuthenticationMiddleware(), func(c *gin.Context) {

		identifier := strings.TrimSpace(strings.ToLower(c.Param("merchantID")))
		bantupayUser := strings.TrimSpace(strings.ToLower(c.Param("targetUser")))

		conDB.PrintDBStats(fmt.Sprintf("POST /v2/merchants/%v/%v/push?", identifier, bantupayUser), db)
		mInfo, err := merchantServices.GetMerchant(identifier, db)

		if err != nil {
			log.Println("[GET MERCHANT FOR PUSH] error for merchant:", identifier, "error: ", err)

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
			response := gin.H{"error": "error-invalid-merchant-access", "data": "Authentication", "message": "Authentication failed"}
			c.JSON(statusCode, response)
			return
		}
		// log.Printf("Merchant Infor: %+v\n", mInfo)
		if mInfo.PushNotificationPermission == 0 {
			//wrong access
			statusCode := http.StatusUnauthorized
			response := gin.H{"error": "error-invalid-merchant-access", "data": "Permission", "message": "Push notification permission not enabled for this merchant"}
			c.JSON(statusCode, response)
			return
		}

		userInfo, err := users.GetUserForMerchants(bantupayUser, mInfo, db, dynamicLinkServiceUrlChan, redisCache)

		if err != nil {
			log.Println("[GET UserInfo] error for user:", bantupayUser, "error: ", err)

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
		if len(userInfo.PushNotificationToken) < 50 {
			//no valid token set. user cannot receive push notification

			c.JSON(http.StatusOK, successResponseData)
			return
		}

		var merchantRequestInput merchantModels.MerchantPushNotificationInput
		reqBody, _ := ioutil.ReadAll(c.Request.Body)

		err = json.Unmarshal(reqBody, &merchantRequestInput)

		var invalidJSON tErrors.ErrorInvalidJSON

		if err != nil {
			log.Println("Push Notification Request Input JSON Error:", err)
			c.JSON(http.StatusBadRequest, invalidJSON.JSONError())
			return
		}
		//check if message is set
		if len(merchantRequestInput.Message) == 0 {

			c.JSON(http.StatusOK, successResponseData)
			return
		}

		if len(merchantRequestInput.Message) > 100 {
			message100Bytes := make([]byte, 0)

			//trim to 100 bytes
			for _, c := range []byte(merchantRequestInput.Message) {
				if (len(message100Bytes) + len(string(c))) <= 100 {
					message100Bytes = append(message100Bytes, c)
					if len(message100Bytes) == 100 {
						break
					}
				}
			}
			merchantRequestInput.Message = string(message100Bytes)
		}

		//Push Message
		merchantRequestInput.PushMessage(userInfo.PushNotificationToken)

		c.JSON(http.StatusOK, successResponseData)
	})

	//user authorization approval url
	router.POST("/v2/users/merchants/:targetUser/authorize/:merchantID/:authID", middleware.AuthenticationMiddleware(), func(c *gin.Context) {

		identifier := strings.TrimSpace(strings.ToLower(c.Param("targetUser")))
		merchant := strings.TrimSpace(strings.ToLower(c.Param("merchantID")))
		authID := strings.TrimSpace(c.Param("authID"))

		conDB.PrintDBStats(fmt.Sprintf("POST /v2/users/merchants/%v/authorize/%v/%v", identifier, merchant, authID), db)
		if identifier == "null" {
			log.Printf("user cannot be %v\n", identifier)
			c.JSON(http.StatusBadRequest, gin.H{"error": "user cannot be null"})
			return
		}
		mInfo, err := merchantServices.GetMerchant(merchant, db)

		if err != nil {
			log.Println("[GET MERCHANT] error for merchant:", merchant, "error: ", err)

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
			response := gin.H{"error": "error-invalid-merchant-access", "data": "Permission", "message": "authorization permission not enabled for this merchant"}
			log.Printf("%+v\n", response)
			c.JSON(statusCode, response)
			return
		}

		//check if merchant is for an event registration/reward merchant: [2 = registration, 1 = reward, 0 = none]
		if mInfo.RewardOnly == 2 {

			userInfo, err := users.GetUserForMerchants(middleware.ExtractPublicKey(c), mInfo, db, dynamicLinkServiceUrlChan, redisCache)

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

			//get authorization data for user
			authData, err := merchantServices.GetEventAuthorizationData(merchant, authID, db)
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
			//return report to user and not keep them waiting.
			c.JSON(http.StatusOK, gin.H{"message": "success"})

			//TODO: make callback request if callback is available
			if authData.CallbackURL != nil {

				//make callback request
				// callbackResponse := new(map[string]interface{})
				type payload struct {
					AuthID     string `json:"authId"`
					TargetUser string `json:"targetUser"`
					DeviceID   string `json:"deviceId"`
				}
				jsonPayload := payload{AuthID: authData.ID, TargetUser: userInfo.Username, DeviceID: c.GetHeader("X-BANTUPAY-DEVICE-ID")}
				/////
				d := *authData.CallbackURL
				body, err := json.Marshal(jsonPayload)
				if err != nil {
					log.Printf("[EVENT CALLBACK AuthCallback] could not unmarshal callback message due to [%v]\n", err)

				}
				log.Printf("[EVENT CALLBACK AuthCallback] JSON STRING: [%v]\n", string(body))

				responseBody := bytes.NewBuffer(body)
				//Leverage Go's HTTP Post function to make request
				resp, err := http.Post(d, "application/json", responseBody)
				//Handle Error
				if err != nil {
					log.Printf("[EVENT CALLBACK AuthCallback] could not send callback message due to [%v]\n", err)
					return
				}
				defer resp.Body.Close()
				//Read the response body
				body, err = ioutil.ReadAll(resp.Body)
				if err != nil {
					//send to retry channel
					log.Println("[EVENT CALLBACK AuthCallback] callback failed:", err)
					c := retryCallbacks{Req: responseBody, CallbackURL: d, Count: 0}
					callBackRetryChan <- c
				} else {
					log.Printf("[EVENT CALLBACK AuthCallback] authorization Callback successful to: [%v], Response:[%v]\n\n", d, string(body))

				}

			}

			return

		} else if mInfo.RewardOnly == 1 {

			userInfo, err := users.GetUserForMerchants(middleware.ExtractPublicKey(c), mInfo, db, dynamicLinkServiceUrlChan, redisCache)

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

			//get authorization data for user
			authData, err := merchantServices.GetRewardOnlyAuthorizationData(merchant, authID, db)
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
			//return report to user and not keep them waiting.
			c.JSON(http.StatusOK, gin.H{"message": "success"})

			//TODO: make callback request if callback is available
			if authData.CallbackURL != nil {

				//make callback request
				// callbackResponse := new(map[string]interface{})
				type payload struct {
					AuthID     string `json:"authId"`
					TargetUser string `json:"targetUser"`
					DeviceID   string `json:"deviceId"`
				}
				jsonPayload := payload{AuthID: authData.ID, TargetUser: userInfo.Username, DeviceID: c.GetHeader("X-BANTUPAY-DEVICE-ID")}

				/////
				d := *authData.CallbackURL
				body, err := json.Marshal(jsonPayload)
				if err != nil {
					log.Printf("[REWARD AuthCallback] could not unmarshal callback message due to [%v]\n", err)

				}
				log.Printf("[REWARD AuthCallback] JSON STRING: [%v]\n", string(body))

				responseBody := bytes.NewBuffer(body)
				//Leverage Go's HTTP Post function to make request
				resp, err := http.Post(d, "application/json", responseBody)
				//Handle Error
				if err != nil {
					log.Printf("[REWARD AuthCallback] could not send callback message due to [%v]\n", err)
					return
				}
				defer resp.Body.Close()
				//Read the response body
				body, err = ioutil.ReadAll(resp.Body)
				if err != nil {
					//send to retry channel
					log.Println("[REWARD AuthCallback] callback failed:", err)
					c := retryCallbacks{Req: responseBody, CallbackURL: d, Count: 0}
					callBackRetryChan <- c
				} else {
					log.Printf("[REWARD AuthCallback] authorization Callback successful to: [%v], Response:[%v]\n\n", d, string(body))

				}

			}

			return

		} else {

			userInfo, err := users.GetUserForMerchants(identifier, mInfo, db, dynamicLinkServiceUrlChan, redisCache)

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
			if (middleware.ExtractPublicKey(c) != userInfo.PublicKey) && mInfo.RewardOnly == 0 {
				//wrong access
				statusCode := http.StatusUnauthorized
				response := gin.H{"error": "error-invalid-user-access", "data": "Authentication", "message": "2FA/Authorization request does not belong to your Bantupay wallet"}
				c.JSON(statusCode, response)
				return
			}
			//get authorization data for user
			authData, err := merchantServices.GetUserAuthorizationData(merchant, userInfo.Username, authID, db)
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
			if authData.Authorized == 1 {
				response := gin.H{"error": "error-authorization-does-not-exist", "data": userInfo.Username, "message": "invalid/expired authorization request"}
				statusCode := http.StatusNotFound
				c.JSON(statusCode, response)
				return
			}
			//exists and needs to be updated
			authData.Authorized = 1

			err = db.Save(&authData).Error
			if err != nil {
				//could not save authData
				response := gin.H{"error": "error-temporary-server-error", "data": "temporaryServerError", "message": "Temporary Server Error. Contact support."}
				statusCode := http.StatusServiceUnavailable
				c.JSON(statusCode, response)
				return
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
				jsonPayload := payload{AuthID: authData.ID, TargetUser: authData.WalletUsername, DeviceID: c.GetHeader("X-BANTUPAY-DEVICE-ID")}

				/////
				d := *authData.CallbackURL
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
				body, err = ioutil.ReadAll(resp.Body)
				if err != nil {
					//send to retry channel
					log.Println("[LoginAuthCallback] callback failed:", err)
					c := retryCallbacks{Req: responseBody, CallbackURL: d, Count: 0}
					callBackRetryChan <- c
				} else {
					log.Printf("[LoginAuthCallback] Login authorization Callback successful to: [%v], Response:[%v]\n\n", d, string(body))

				}

			}
		}

	})

	//merchant authorization verify url
	router.GET("/v2/merchants/:merchantID/:targetUser/authorize/:authID", middleware.AuthenticationMiddleware(), func(c *gin.Context) {
		// var err error

		identifier := strings.TrimSpace(strings.ToLower(c.Param("merchantID")))
		bantupayUser := strings.TrimSpace(strings.ToLower(c.Param("targetUser")))
		authID := strings.TrimSpace(strings.ToLower(c.Param("authID")))

		conDB.PrintDBStats(fmt.Sprintf("GET /v2/merchants/%v/%v/authorize/%v", identifier, bantupayUser, authID), db)
		mInfo, err := merchantServices.GetMerchant(identifier, db)

		if err != nil {
			log.Println("[GET MERCHANT] error for merchant:", identifier, "error: ", err)

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
			response := gin.H{"error": "error-invalid-merchant-access", "data": "Authentication", "message": "Authentication failed"}
			c.JSON(statusCode, response)
			return
		}
		// log.Printf("Merchant Info: %+v\n", mInfo)
		if mInfo.AuthorizationPermission == 0 {
			//wrong access
			statusCode := http.StatusUnauthorized
			response := gin.H{"error": "error-invalid-merchant-access", "data": "Permission", "message": "authorization permission not enabled for this merchant"}
			c.JSON(statusCode, response)
			return
		}

		userInfo, err := users.GetUserForMerchants(bantupayUser, mInfo, db, dynamicLinkServiceUrlChan, redisCache)

		if err != nil {
			log.Println("[GET UserInfo] error for user:", bantupayUser, "error: ", err)

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

		authData, err := merchantServices.GetUserAuthorizationData(identifier, userInfo.Username, authID, db)
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

	//merchant payment request
	router.GET("/v2/merchants/:merchantID/:targetUser/payment", middleware.AuthenticationMiddleware(), func(c *gin.Context) {

		identifier := strings.TrimSpace(strings.ToLower(c.Param("merchantID")))
		bantupayUser := strings.TrimSpace(strings.ToLower(c.Param("targetUser")))
		paymentDestination := strings.TrimSpace(strings.ToLower(c.Query("paymentDestination")))
		if len(paymentDestination) == 56 {
			paymentDestination = strings.ToUpper(paymentDestination)
		}
		assetCode := strings.TrimSpace(strings.ToUpper(c.Query("assetCode")))
		assetIssuer := strings.TrimSpace(strings.ToUpper(c.Query("assetIssuer")))
		amount := strings.TrimSpace(c.Query("amount"))
		memo := strings.TrimSpace(c.Query("memo"))

		cacheKey := fmt.Sprintf("[GET] /v2/merchants/%v/%v/payment", identifier, bantupayUser)
		cacheKeyParameters := fmt.Sprintf("%v", c.Request.URL.RawQuery)

		{
			//search cache

			ok, status, response := redisCache.CachedHttpResponseWithParameters(cacheKey, cacheKeyParameters)

			if ok {
				log.Printf("[%v]/[%v], served from cache\n", cacheKey, cacheKeyParameters)
				c.JSON(status, response)
				return
			}
		}

		conDB.PrintDBStats(fmt.Sprintf("GET /v2/merchants/%v/%v/payment?paymentDestination=%v&assetCode=%v&assetIssuer=%v&amount=%v&memo=%v", identifier, bantupayUser, paymentDestination, assetCode, assetIssuer, amount, memo), db)
		mInfo, err := merchantServices.GetMerchant(identifier, db)

		if err != nil {
			log.Println("[GET MERCHANT PAYMENT DATA] error for merchant:", identifier, "error: ", err)

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
			response := gin.H{"error": "error-invalid-merchant-access", "data": "Authentication", "message": "Authentication failed"}
			c.JSON(statusCode, response)
			return
		}
		// log.Printf("Merchant Infor: %+v\n", mInfo)
		if mInfo.PaymentPermission == 0 {
			//wrong access
			statusCode := http.StatusUnauthorized
			response := gin.H{"error": "error-invalid-merchant-access", "data": "Permission", "message": "payment permission not enabled for this merchant"}
			c.JSON(statusCode, response)
			return
		}

		_, err = users.GetUserForMerchants(bantupayUser, mInfo, db, dynamicLinkServiceUrlChan, redisCache)

		if err != nil {
			log.Println("[GET UserInfo] error for user:", bantupayUser, "error: ", err)

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
		data, err := merchantServices.GeneratePaymentData(paymentDestination, assetCode, assetIssuer, amount, memo, dynamicLinkServiceUrlChan, redisCache)
		if err != nil {
			//could not create login session
			response := gin.H{"error": "error-temporary-server-error", "data": "temporaryServerError", "message": "Temporary Server Error. Contact support."}
			statusCode := http.StatusServiceUnavailable
			c.JSON(statusCode, response)
			return
		}
		{
			cacheDurationInSeconds := 20 * 60 //2 minutes

			redisCache.CacheHttpResponseWithParameters(cacheKey, cacheKeyParameters, http.StatusOK, data, cacheDurationInSeconds)
		}
		c.JSON(http.StatusOK, data)

	})

	//merchant payment request
	router.GET("/v2/merchants/:merchantID/:targetUser/userinfo", middleware.AuthenticationMiddleware(), func(c *gin.Context) {

		identifier := strings.TrimSpace(strings.ToLower(c.Param("merchantID")))
		bantupayUser := strings.TrimSpace(strings.ToLower(c.Param("targetUser")))
		{
			//do not cache this so that it brings the latest data
			// cacheKey := fmt.Sprintf("[GET] /v2/merchants/%v/%v/payment", identifier, bantupayUser)
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
		conDB.PrintDBStats(fmt.Sprintf("GET /v2/merchants/%v/%v/userinfo", identifier, bantupayUser), db)
		mInfo, err := merchantServices.GetMerchant(identifier, db)

		if err != nil {
			log.Println("[GET USER DATA] error for merchant:", identifier, "error: ", err)

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
			response := gin.H{"error": "error-invalid-merchant-access", "data": "Authentication", "message": "Authentication failed"}
			c.JSON(statusCode, response)
			return
		}
		// log.Printf("Merchant Info: %+v\n", mInfo)
		if mInfo.AllowUserInfo == 0 {
			//wrong access
			statusCode := http.StatusUnauthorized
			response := gin.H{"error": "error-invalid-merchant-access", "data": "Permission", "message": "userInfo permission not enabled for this merchant"}
			c.JSON(statusCode, response)
			return
		}

		data, err := users.GetUserForMerchants(bantupayUser, mInfo, db, dynamicLinkServiceUrlChan, redisCache)

		if err != nil {
			log.Println("[GET UserInfo] error for user:", bantupayUser, "error: ", err)

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

}
