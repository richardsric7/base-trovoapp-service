package users

import (
	paymentServices "trovo-wallet-api/internal/components/payments/services"
	usersDB "trovo-wallet-api/internal/components/users/db"
	userModels "trovo-wallet-api/internal/components/users/models"
	userServices "trovo-wallet-api/internal/components/users/services"
	conDB "trovo-wallet-api/internal/db"
	dl "trovo-wallet-api/internal/dynamiclinks"
	tErrors "trovo-wallet-api/internal/errors"
	pns "trovo-wallet-api/internal/pns"

	"fmt"

	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"trovo-wallet-api/internal/middleware"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/gin-gonic/gin"
)

// Init initializes /v1/users endpoint
func Init(router *gin.Engine, gc *sharedconfig.GlobalConfig) {
	//websocket stream
	router.GET("/v1/users/:targetUser/ws", func(c *gin.Context) {

		identifier := strings.TrimSpace(strings.ToLower(c.Param("targetUser")))
		uDec, e := base64.URLEncoding.DecodeString(c.Param("targetUser"))
		if e == nil {
			//check if the decoded contains any non-english character
			invalidChars := 0

			acceptedChars := "abcdefghijklmnopqrstuvwxyz_1234567890/"
			for _, c := range uDec {

				if !strings.Contains(acceptedChars, strings.TrimSpace(strings.ToLower(string(c)))) {
					invalidChars++
				}

			}
			if invalidChars == 0 {
				identifier = string(uDec)
			}

		}
		if identifier == "null" {
			log.Printf("user cannot be %v\n", identifier)
			c.JSON(http.StatusBadRequest, gin.H{"error": "user cannot be null"})
			return
		}
		conDB.PrintDBStats(fmt.Sprintf("/v1/users/%v/ws", identifier), gc.DB)
		log.Printf("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<Websocket connection detected for %v>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>\n", identifier)

		userServices.UserWebSocketAPI(c, gc)

	})

	router.GET("/v1/users/:targetUser/payments/:targetPublicKeyForHistory", middleware.AuthenticationMiddlewareUsingTimestamp(), func(c *gin.Context) {
		// var err error

		identifier := strings.TrimSpace(strings.ToLower(c.Param("targetUser")))
		targetPublicKeyForHistory := strings.TrimSpace(strings.ToUpper(c.Param("targetPublicKeyForHistory")))
		uDec, e := base64.URLEncoding.DecodeString(c.Param("targetUser"))
		if e == nil {
			//check if the decoded contains any non-english character
			invalidChars := 0

			acceptedChars := "abcdefghijklmnopqrstuvwxyz_1234567890/"
			for _, c := range uDec {

				if !strings.Contains(acceptedChars, strings.TrimSpace(strings.ToLower(string(c)))) {
					invalidChars++
				}

			}
			if invalidChars == 0 {
				identifier = string(uDec)
			}

		}
		if identifier == "null" {
			log.Printf("[GET PAYMENTS HISTORY] user cannot be %v\n", identifier)
			c.JSON(http.StatusBadRequest, gin.H{"error": "user cannot be null"})
			return
		}
		cacheKey := fmt.Sprintf("[GET] /v1/users/%v/payments/%v", identifier, targetPublicKeyForHistory)

		conDB.PrintDBStats(fmt.Sprintf(" /v1/users/%v/payments/%v", identifier, targetPublicKeyForHistory), gc.DB)

		// cacheDurationInSeconds := 1 * 60 //1 minutes
		cacheDurationInSeconds := 10 //1 minutes

		userInfo, err := userServices.GetUserInfo(identifier, gc, c)

		if err != nil {
			log.Println("[GET USERINFO] error for user:", identifier, "error: ", err)

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
			gc.RedisCache.CacheHttpResponse(cacheKey, statusCode, response, cacheDurationInSeconds)
			return
		}

		{
			gc.RedisCache.InvalidateCachedHttpResponse(cacheKey)

			//check if the owner is the one accessing it or if the one accessing it has access to access it.
			for _, v := range userInfo.UserData.UserWallets {
				if v.PrimaryWallet == 1 {
					if v.Signer != middleware.ExtractSigner(c) && !userServices.HasAccessToPublicKey(userInfo.UserData.PublicKey, targetPublicKeyForHistory, gc) {
						te := &tErrors.ErrorInvalidAuthorization{}

						log.Println("[GET HISTORY] Invalid signer for user:", identifier, "error: ", err)
						c.JSON(te.HTTPCode(), te.JSONError())
						return
					}
				}
			}

		}
		//Get Payment history
		historyRecords := paymentServices.GetPaymentHistory(targetPublicKeyForHistory, gc, c)

		c.JSON(http.StatusOK, historyRecords)
		gc.RedisCache.CacheHttpResponse(cacheKey, http.StatusOK, historyRecords, cacheDurationInSeconds)

	})

	router.GET("/v1/users/:targetUser", middleware.AuthenticationMiddlewareUsingTimestamp(), func(c *gin.Context) {
		// var err error

		identifier := strings.TrimSpace(strings.ToLower(c.Param("targetUser")))
		uDec, e := base64.URLEncoding.DecodeString(c.Param("targetUser"))
		if e == nil {
			//check if the decoded contains any non-english character
			invalidChars := 0

			acceptedChars := "abcdefghijklmnopqrstuvwxyz_1234567890/"
			for _, c := range uDec {

				if !strings.Contains(acceptedChars, strings.TrimSpace(strings.ToLower(string(c)))) {
					invalidChars++
				}

			}
			if invalidChars == 0 {
				identifier = string(uDec)
			}

		}
		if identifier == "null" {
			log.Printf("user cannot be %v\n", identifier)
			c.JSON(http.StatusBadRequest, gin.H{"error": "user cannot be null"})
			return
		}
		cacheKey := fmt.Sprintf("[GET] /v1/users/%v", identifier)

		conDB.PrintDBStats(fmt.Sprintf("/v1/users/%v", identifier), gc.DB)

		//check if type is import
		queryType := strings.ToLower(c.Query("type"))

		cacheDurationInSeconds := 1 * 60 //1 minutes

		userInfo, err := userServices.GetUserInfo(identifier, gc, c)

		if err != nil {
			log.Println("[GET USERINFO] error for user:", identifier, "error: ", err)

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
			gc.RedisCache.CacheHttpResponse(cacheKey, statusCode, response, cacheDurationInSeconds)
			return
		}

		if queryType == "import" {
			gc.RedisCache.InvalidateCachedHttpResponse(cacheKey)

			log.Println("Wallet import request received from:", identifier, "for:", middleware.ExtractPublicKey(c), "........")
			//perform import specific tasks

			//check if username key matches with import credential
			if userInfo.AssetBalances == nil && userInfo.DefaultAssets == nil && userInfo.UserData.UserWallets == nil {
				te := &tErrors.CustomError{Param: "username",
					Err:        "error-wallet-does-not-belong-to-username",
					ErrMessage: fmt.Sprintf("The wallet you are importing does not belong to %s", identifier),
					Code:       http.StatusBadRequest}

				log.Println("[Wallet import] Invalid import-credential for user:", identifier, "error: ", err)
				c.JSON(te.HTTPCode(), te.JSONError())
				return
			}
			//check if the owner is the one importing it
			for _, v := range userInfo.UserData.UserWallets {
				if v.PrimaryWallet == 1 {
					if v.Signer != middleware.ExtractSigner(c) {
						te := &tErrors.ErrorInvalidAuthorization{}

						log.Println("[Wallet import] Invalid signer for user:", identifier, "error: ", err)
						c.JSON(te.HTTPCode(), te.JSONError())
						return
					}
				}
			}

		}

		c.JSON(http.StatusOK, userInfo)
		gc.RedisCache.CacheHttpResponse(cacheKey, http.StatusOK, userInfo, cacheDurationInSeconds)

	})

	router.POST("/v1/users", middleware.AuthenticationMiddlewareUsingTimestamp(), func(c *gin.Context) {
		var err error

		var userRegistrationInfo userModels.UserRegistrationInfo
		// var err error

		data, _ := ioutil.ReadAll(c.Request.Body)

		err = json.Unmarshal(data, &userRegistrationInfo)

		userRegistrationInfo.PublicKey = middleware.ExtractSigner(c)
		userRegistrationInfo.PublicIP = c.ClientIP()

		var invalidJSON tErrors.ErrorInvalidJSON

		if err != nil {
			c.JSON(http.StatusBadRequest, invalidJSON.JSONError())
			return
		}

		if strings.Contains(userRegistrationInfo.Username, "%") {
			u, e := url.QueryUnescape(userRegistrationInfo.Username)
			if e == nil {
				userRegistrationInfo.Username = u
			}
		}

		//replace _ and /
		userRegistrationInfo.Username = strings.ReplaceAll(strings.ReplaceAll(strings.ToLower(userRegistrationInfo.Username), "_", ""), "/", "")

		conDB.PrintDBStats(fmt.Sprintf("POST /v1/users %v", userRegistrationInfo.Username), gc.DB)

		var emailSent bool

		_, emailSent, err = userServices.RegisterUser(userRegistrationInfo, gc)

		if err != nil {
			var ex tErrors.GenericError
			var ok bool

			ex, ok = err.(tErrors.GenericError)
			if ok {
				c.JSON(http.StatusBadRequest, ex.JSONError())
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}
			return
		}

		//At this point, there was no error.
		//But either the email was sent or not.
		if emailSent {
			c.JSON(http.StatusAccepted, gin.H{"message": "Verification code sent to your email"})
			//send push notificationMessage
			if len(userRegistrationInfo.PushNotificationToken) > 50 {
				dataPayload := make(map[string]string)
				dataPayload["route"] = ""
				pns.SendFirebaseMessage(userRegistrationInfo.PushNotificationToken, "Verification code sent to your email", "Please check your email to get the verification code. It is only valid today.", "", dataPayload, gc.PushNotificationClient, gc.PNSContext)
			}
		} else {
			//registration completed
			dataPayload := make(map[string]string)
			dataPayload["route"] = ""
			//return response
			c.JSON(http.StatusOK, gin.H{"message": userRegistrationInfo.PublicKey})
			pns.SendFirebaseMessage(userRegistrationInfo.PushNotificationToken, "Registration completed!", fmt.Sprintf("Congratulations! Your trovo wallet account has successfully been created. To receive payment, you can share your primary account username  %s (also known as your alias) to your friends or you can use your public key for payments outside of Trovo Ecosystem. Please take the very important step to backup your wallet or use the available option to enable Account Recovery (Terms and Conditions apply). Thank you!", userRegistrationInfo.Username), "", dataPayload, gc.PushNotificationClient, gc.PNSContext)

		}
	})

	router.POST("/v1/users/subwallet", middleware.AuthenticationMiddlewareUsingTimestamp(), func(c *gin.Context) {
		var err error

		var subWalletInfo userModels.SubWalletInfo
		// var err error

		data, _ := ioutil.ReadAll(c.Request.Body)

		err = json.Unmarshal(data, &subWalletInfo)

		var invalidJSON tErrors.ErrorInvalidJSON

		if err != nil {
			c.JSON(http.StatusBadRequest, invalidJSON.JSONError())
			return
		}

		user, err := usersDB.GetUser(middleware.ExtractSigner(c), gc.DB)

		if err != nil {
			var ex tErrors.GenericError
			var ok bool

			ex, ok = err.(tErrors.GenericError)
			if ok {
				c.JSON(http.StatusBadRequest, ex.JSONError())
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}
			return
		}

		conDB.PrintDBStats(fmt.Sprintf("POST /v1/users/subwallet %v", user.Username), gc.DB)

		returnedSubwalletInfo, err := userServices.CreateNewSubWallet(&user, &subWalletInfo, gc)

		if err != nil {
			var ex tErrors.GenericError
			var ok bool

			ex, ok = err.(tErrors.GenericError)
			if ok {
				c.JSON(http.StatusBadRequest, ex.JSONError())
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}
			return
		}

		if user.PushNotificationToken != nil && len(returnedSubwalletInfo.TransactionID) > 0 && returnedSubwalletInfo.TransactionID != "PENDING_AUTH" {
			dataPayload := make(map[string]string)
			dataPayload["route"] = ""
			pns.SendFirebaseMessage(*user.PushNotificationToken, "New Sub-wallet Added!", "You have successfully added a new sub wallet.", "", dataPayload, gc.PushNotificationClient, gc.PNSContext)
		}

		//At this point, there was no error.

		c.JSON(http.StatusOK, returnedSubwalletInfo)
	})

	router.PUT("/v1/users/:targetUser/actions/claim-asset", middleware.AuthenticationMiddlewareUsingTimestamp(), func(c *gin.Context) {
		var err error

		identifier := strings.TrimSpace(strings.ToLower(c.Param("targetUser")))
		uDec, e := base64.URLEncoding.DecodeString(c.Param("targetUser"))
		if e == nil {
			//check if the decoded contains any non-english character
			invalidChars := 0

			acceptedChars := "abcdefghijklmnopqrstuvwxyz_1234567890/"
			for _, c := range uDec {

				if !strings.Contains(acceptedChars, strings.TrimSpace(strings.ToLower(string(c)))) {
					invalidChars++
				}

			}
			if invalidChars == 0 {
				identifier = string(uDec)
			}

		}

		conDB.PrintDBStats(fmt.Sprintf("PUT /v1/users/:identifier/actions/claim-asset %v", identifier), gc.DB)
		if identifier == "null" {
			log.Printf("user cannot be %v\n", identifier)
			c.JSON(http.StatusBadRequest, gin.H{"error": "user cannot be null"})
			return
		}
		var pendingAssetToClaim userModels.PendingAssetToClaim
		// var err error

		data, _ := ioutil.ReadAll(c.Request.Body)

		err = json.Unmarshal(data, &pendingAssetToClaim)

		var invalidJSON tErrors.ErrorInvalidJSON

		if err != nil {
			c.JSON(http.StatusBadRequest, invalidJSON.JSONError())
			return
		}

		returnedPendingAssetToClaim, complete, err := userServices.ClaimPendingAsset(identifier, middleware.ExtractPublicKey(c), &pendingAssetToClaim, gc.DB)

		if err != nil {
			var ex tErrors.GenericError
			var ok bool

			ex, ok = err.(tErrors.GenericError)
			if ok {
				c.JSON(ex.HTTPCode(), ex.JSONError())
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}
			return
		}

		cacheKey := fmt.Sprintf("[GET] /v1/users/%v", identifier)
		paymentHistoryCacheKey := fmt.Sprintf("[GET] /v1/users/%v/payments", identifier)
		senderCacheKey := fmt.Sprintf("[GET] /v1/users/%v", identifier)
		senderPaymentHistoryCacheKey := fmt.Sprintf("[GET] /v1/users/%v/payments", identifier)

		gc.RedisCache.InvalidateCachedHttpResponse(senderCacheKey, senderPaymentHistoryCacheKey)

		gc.RedisCache.InvalidateCachedHttpResponse(cacheKey, paymentHistoryCacheKey)

		if complete {
			c.JSON(http.StatusOK, returnedPendingAssetToClaim)
		} else {
			c.JSON(http.StatusAccepted, returnedPendingAssetToClaim)
		}

	})

	router.GET("/v1/users/:targetUser/generate/payment", middleware.AuthenticationMiddlewareUsingTimestamp(), func(c *gin.Context) {
		var err error

		identifier := strings.TrimSpace(strings.ToLower(c.Param("targetUser")))
		uDec, e := base64.URLEncoding.DecodeString(c.Param("targetUser"))
		if e == nil {
			//check if the decoded contains any non-english character
			invalidChars := 0

			acceptedChars := "abcdefghijklmnopqrstuvwxyz_1234567890/"
			for _, c := range uDec {

				if !strings.Contains(acceptedChars, strings.TrimSpace(strings.ToLower(string(c)))) {
					invalidChars++
				}

			}
			if invalidChars == 0 {
				identifier = string(uDec)
			}

		}

		if identifier == "null" {
			log.Printf("user cannot be %v\n", identifier)
			c.JSON(http.StatusBadRequest, gin.H{"error": "user cannot be null"})
			return
		}

		paymentDestination := strings.TrimSpace(strings.ToLower(c.Query("paymentDestination")))
		if len(paymentDestination) == 56 {
			paymentDestination = strings.ToUpper(paymentDestination)
		}
		assetCode := strings.TrimSpace(strings.ToUpper(c.Query("assetCode")))
		assetIssuer := strings.TrimSpace(strings.ToUpper(c.Query("assetIssuer")))
		amount := strings.TrimSpace(c.Query("amount"))
		memo := strings.TrimSpace(c.Query("memo"))

		cacheKey := fmt.Sprintf("[GET] /v1/users/%v/generate/payment", identifier)
		cacheKeyParameters := fmt.Sprintf("%v", c.Request.URL.RawQuery)

		{
			//search cache

			// cacheKeyParameters := fmt.Sprintf("limit=%v&order=%v&cursor=%v&forTransactionHash=%v&includeHash=%v&temp=%v", limit, orderStr, cursor, forTransactionHash, includeHash, temp)

			ok, status, response := gc.RedisCache.CachedHttpResponseWithParameters(cacheKey, cacheKeyParameters)

			if ok {
				log.Printf("[%v]/[%v], served from cache\n", cacheKey, cacheKeyParameters)
				c.JSON(status, response)
				return
			}
		}

		conDB.PrintDBStats(fmt.Sprintf("GET /v1/users/%v/generate/payment?paymentDestination=%v&assetCode=%v&assetIssuer=%v&amount=%v&memo=%v", identifier, paymentDestination, assetCode, assetIssuer, amount, memo), gc.DB)

		_, err = usersDB.GetUser(identifier, gc.DB)

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
			cacheDurationInSeconds := 525600 * 3 * 60 //3yrs

			gc.RedisCache.CacheHttpResponseWithParameters(cacheKey, cacheKeyParameters, http.StatusOK, data, cacheDurationInSeconds)
		}
		c.JSON(http.StatusOK, data)

	})

}
