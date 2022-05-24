package users

import (
	"trovo-wallet-api/internal/cache"
	merchantServices "trovo-wallet-api/internal/components/merchants/services"
	usersDB "trovo-wallet-api/internal/components/users/db"
	usermodels "trovo-wallet-api/internal/components/users/models"
	users "trovo-wallet-api/internal/components/users/services"
	conDB "trovo-wallet-api/internal/db"
	bantupayErrors "trovo-wallet-api/internal/errors"

	"fmt"
	"os"

	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"trovo-wallet-api/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/stellar/go/clients/horizonclient"
	"gorm.io/gorm"
)

// Init initializes /v2/users endpoint
func Init(router *gin.Engine, db *gorm.DB, redisCache *cache.RedisCache, dynamicLinkServiceUrlChan chan string) {
	//websocket stream
	router.GET("/v2/users/:targetUser/ws", func(c *gin.Context) {

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
		conDB.PrintDBStats(fmt.Sprintf("/v2/users/%v/ws", identifier), db)
		log.Printf("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<Websocket connection detected for %v>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>\n", identifier)

		users.UserWebSocketAPI(c, db, redisCache)

	})

	router.GET("/v2/users/:targetUser/payments", middleware.AuthenticationMiddleware(), func(c *gin.Context) {

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
		limitStr := c.DefaultQuery("limit", "25")
		limit, _ := strconv.ParseUint(limitStr, 10, 64)
		orderStr := c.DefaultQuery("order", "desc")
		order := horizonclient.Order(orderStr)
		cursor := c.Query("cursor")
		if len(cursor) > 5 {
			order = horizonclient.OrderAsc
		}
		forTransactionHash := c.Query("forTransactionHash")
		includeH := c.DefaultQuery("includeHash", "false")
		tempStr := c.DefaultQuery("temp", "false")
		var includeHash, temp bool
		if includeH == "true" {
			includeHash = true
		}
		if tempStr == "true" {
			temp = true
		}

		cacheKey := fmt.Sprintf("[GET] /v2/users/%v/payments", identifier)
		cacheKeyParameters := fmt.Sprintf("limit=%v&order=%v&cursor=%v&forTransactionHash=%v&includeHash=%v&temp=%v", limit, orderStr, cursor, forTransactionHash, includeHash, temp)

		ok, status, response := redisCache.CachedHttpResponseWithParameters(cacheKey, cacheKeyParameters)

		if ok {
			log.Printf("[%v]/[%v], served from cache\n", cacheKey, cacheKeyParameters)
			c.JSON(status, response)
			return
		}

		conDB.PrintDBStats(cacheKey, db)

		cacheDurationInSeconds := 2 * 60 //2 minutes

		paymentHistory, err := users.GetUserPaymentHistory(identifier, db, middleware.ExtractPublicKey(c), uint(limit), order, cursor, forTransactionHash, includeHash, temp)

		if err != nil {
			var ex bantupayErrors.GenericError
			var ok bool

			var _statusCode int = 0
			var _response interface{}

			ex, ok = err.(bantupayErrors.GenericError)
			if ok {
				_statusCode = ex.HTTPCode()
				_response = ex.JSONError()

			} else {
				_statusCode = http.StatusBadRequest
				_response = gin.H{"error": err.Error()}

			}

			c.JSON(_statusCode, _response)
			redisCache.CacheHttpResponseWithParameters(cacheKey, cacheKeyParameters, _statusCode, _response, cacheDurationInSeconds)
			return
		}
		c.JSON(http.StatusOK, paymentHistory)
		redisCache.CacheHttpResponseWithParameters(cacheKey, cacheKeyParameters, http.StatusOK, paymentHistory, cacheDurationInSeconds)

	})

	router.GET("/v2/users/:targetUser", middleware.AuthenticationMiddleware(), func(c *gin.Context) {
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
		cacheKey := fmt.Sprintf("[GET] /v2/users/%v", identifier)

		conDB.PrintDBStats(fmt.Sprintf("/v2/users/%v", identifier), db)

		var budsInfo usermodels.BudsInfo

		//check if type is import
		queryType := strings.ToLower(c.Query("type"))
		if queryType == "import" {
			redisCache.InvalidateCachedHttpResponse(cacheKey)

			log.Println("BUDS wallet import request received from:", identifier, "for:", middleware.ExtractPublicKey(c), "........")
			//perform import specific tasks

			budsInfo, err = users.GetUser(identifier, db, middleware.ExtractPublicKey(c), dynamicLinkServiceUrlChan, redisCache)
			if err != nil {
				log.Println("[IMPORT BUDS] error for user:", identifier, "error: ", err)

				var ex bantupayErrors.GenericError
				var ok bool

				ex, ok = err.(bantupayErrors.GenericError)
				if ok {
					c.JSON(ex.HTTPCode(), ex.JSONError())
				} else {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				}
				return
			}
			if len(budsInfo.ReferralLink) == 0 {
				rld, errLink := merchantServices.GenerateReferralLink(budsInfo.Username, dynamicLinkServiceUrlChan, redisCache)
				if errLink == nil {

					user, errUFetch := usersDB.GetUserInfo(budsInfo.Username, db)
					if errUFetch == nil {
						user.ReferralLink = &rld.DynamicLink
						user.ReferralQrCode = &rld.QRCode
						errUFetch = db.Save(&user).Error
						if errUFetch != nil {
							log.Println("[IMPORT BUDS] error updating referral link for user:", identifier, "error: ", errUFetch)
						}
						if errUFetch == nil {
							budsInfo, err = users.GetUser(identifier, db, middleware.ExtractPublicKey(c), dynamicLinkServiceUrlChan, redisCache)
							if err != nil {
								cacheDurationInSeconds := 1 * 60 //1 minutes
								redisCache.CacheHttpResponse(cacheKey, http.StatusOK, budsInfo, cacheDurationInSeconds)
							}
						}
					}
				}

			}

			c.JSON(http.StatusOK, budsInfo)
			return

		}

		if queryType != "import" {
			//let's do some caching here too.

			ok, status, response := redisCache.CachedHttpResponse(cacheKey)

			if ok {
				log.Printf("[%v], served from cache\n", cacheKey)
				c.JSON(status, response)
				return
			}

		}

		cacheDurationInSeconds := 1 * 60 //1 minutes

		budsInfo, err = users.GetUser(identifier, db, middleware.ExtractPublicKey(c), dynamicLinkServiceUrlChan, redisCache)

		if err != nil {
			log.Println("[GET BUDS] error for user:", identifier, "error: ", err)

			var ex bantupayErrors.GenericError
			var ok bool

			ex, ok = err.(bantupayErrors.GenericError)
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
			redisCache.CacheHttpResponse(cacheKey, statusCode, response, cacheDurationInSeconds)
			return
		}

		c.JSON(http.StatusOK, budsInfo)
		redisCache.CacheHttpResponse(cacheKey, http.StatusOK, budsInfo, cacheDurationInSeconds)

	})

	router.POST("/v2/users", middleware.AuthenticationMiddleware(), func(c *gin.Context) {
		var err error
		// db, err := conDB.OpenDb()
		// if err != nil {
		// 	log.Println("--------------------DB error in POST USERS ENDPOINT:", err)
		// 	return
		// }

		var userRegistrationInfo usermodels.UserRegistrationInfo
		// var err error

		data, _ := ioutil.ReadAll(c.Request.Body)

		err = json.Unmarshal(data, &userRegistrationInfo)

		userRegistrationInfo.PublicKey = middleware.ExtractPublicKey(c)
		userRegistrationInfo.PublicIP = c.ClientIP()

		var invalidJSON bantupayErrors.ErrorInvalidJSON

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

		userRegistrationInfo.Username = strings.ToLower(userRegistrationInfo.Username)

		conDB.PrintDBStats(fmt.Sprintf("POST /v2/users %v", userRegistrationInfo.Username), db)

		var emailSent bool

		_, emailSent, err = users.RegisterUser(userRegistrationInfo, db, pool, dynamicLinkServiceUrlChan, redisCache)

		if err != nil {
			var ex bantupayErrors.GenericError
			var ok bool

			ex, ok = err.(bantupayErrors.GenericError)
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
		} else {
			//registration completed
			//pay influencer reward if applicable
			if os.Getenv("INFLUENCER_REWARD_ENABLED") == "1" {
				users.PayInfluencerReward(userRegistrationInfo, db, c)
			}
			//return response
			c.JSON(http.StatusOK, gin.H{"publicKey": userRegistrationInfo.PublicKey})
		}
	})

	router.PUT("/v2/users/:targetUser", middleware.AuthenticationMiddleware(), func(c *gin.Context) {
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

		conDB.PrintDBStats(fmt.Sprintf("PUT /v2/users %v", identifier), db)
		if identifier == "null" {
			log.Printf("user cannot be %v\n", identifier)
			c.JSON(http.StatusBadRequest, gin.H{"error": "user cannot be null"})
			return
		}
		var userUpdateInfo usermodels.UserUpdateInfo

		data, _ := ioutil.ReadAll(c.Request.Body)

		err = json.Unmarshal(data, &userUpdateInfo)

		var invalidJSON bantupayErrors.ErrorInvalidJSON

		if err != nil {
			c.JSON(http.StatusBadRequest, invalidJSON.JSONError())
			return
		}

		_, err = users.UpdateUser(identifier, middleware.ExtractPublicKey(c), userUpdateInfo, c.ClientIP(), db)

		if err != nil {
			var ex bantupayErrors.GenericError
			var ok bool

			ex, ok = err.(bantupayErrors.GenericError)
			if ok {
				c.JSON(http.StatusBadRequest, ex.JSONError())
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}
			return
		}

		cacheKey := fmt.Sprintf("[GET] /v2/users/%v", identifier)
		senderPaymentHistoryCacheKey := fmt.Sprintf("[GET] /v2/users/%v/payments", identifier)
		senderCacheKey := fmt.Sprintf("[GET] /v2/users/%v", identifier)

		redisCache.InvalidateCachedHttpResponse(senderCacheKey, senderPaymentHistoryCacheKey)
		redisCache.InvalidateCachedHttpResponse(cacheKey)

		//At this point, there was no error.
		//But either the email was sent or not.

		// c.JSON(http.StatusOK, gin.H{"username": c.Param("identifier")})
		budsInfo, err := users.GetUser(identifier, db, middleware.ExtractPublicKey(c), dynamicLinkServiceUrlChan, redisCache)

		if err != nil {
			var ex bantupayErrors.GenericError
			var ok bool

			ex, ok = err.(bantupayErrors.GenericError)
			if ok {
				c.JSON(ex.HTTPCode(), ex.JSONError())
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}
			return
		}

		c.JSON(http.StatusOK, budsInfo)

	})

	router.PUT("/v2/users/:targetUser/actions/claim-asset", middleware.AuthenticationMiddleware(), func(c *gin.Context) {
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

		conDB.PrintDBStats(fmt.Sprintf("PUT /v2/users/:identifier/actions/claim-asset %v", identifier), db)
		if identifier == "null" {
			log.Printf("user cannot be %v\n", identifier)
			c.JSON(http.StatusBadRequest, gin.H{"error": "user cannot be null"})
			return
		}
		var pendingAssetToClaim usermodels.PendingAssetToClaim
		// var err error

		data, _ := ioutil.ReadAll(c.Request.Body)

		err = json.Unmarshal(data, &pendingAssetToClaim)

		var invalidJSON bantupayErrors.ErrorInvalidJSON

		if err != nil {
			c.JSON(http.StatusBadRequest, invalidJSON.JSONError())
			return
		}

		returnedPendingAssetToClaim, complete, err := users.ClaimPendingAsset(identifier, middleware.ExtractPublicKey(c), &pendingAssetToClaim, db)

		if err != nil {
			var ex bantupayErrors.GenericError
			var ok bool

			ex, ok = err.(bantupayErrors.GenericError)
			if ok {
				c.JSON(ex.HTTPCode(), ex.JSONError())
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}
			return
		}

		cacheKey := fmt.Sprintf("[GET] /v2/users/%v", identifier)
		paymentHistoryCacheKey := fmt.Sprintf("[GET] /v2/users/%v/payments", identifier)
		senderPaymentHistoryCacheKey := fmt.Sprintf("[GET] /v2/users/%v/payments", identifier)
		senderCacheKey := fmt.Sprintf("[GET] /v2/users/%v", identifier)

		redisCache.InvalidateCachedHttpResponse(senderCacheKey, senderPaymentHistoryCacheKey)

		redisCache.InvalidateCachedHttpResponse(cacheKey, paymentHistoryCacheKey)

		if complete {
			c.JSON(http.StatusOK, returnedPendingAssetToClaim)
		} else {
			c.JSON(http.StatusAccepted, returnedPendingAssetToClaim)
		}

	})

	router.PUT("/v2/users/:targetUser/actions/verify-mobile/:verificationCode", middleware.AuthenticationMiddleware(), func(c *gin.Context) {
		var err error

		identifier := strings.TrimSpace(strings.ToLower(c.Param("targetUser")))
		verificationCode := strings.TrimSpace(c.Param("verificationCode"))
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
		conDB.PrintDBStats(fmt.Sprintf("PUT /v2/users/%v/actions/verify-mobile/%v", identifier, verificationCode), db)

		// var err error

		userInfo, err := usersDB.GetUserInfo(identifier, db)
		if err != nil {
			var ex bantupayErrors.GenericError
			var ok bool

			ex, ok = err.(bantupayErrors.GenericError)
			if ok {
				c.JSON(ex.HTTPCode(), ex.JSONError())
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}
			return
		}

		err = users.CheckPhoneVerificationCode(&userInfo, verificationCode, db, c)

		if err != nil {
			var ex bantupayErrors.GenericError
			var ok bool

			ex, ok = err.(bantupayErrors.GenericError)
			if ok {
				c.JSON(ex.HTTPCode(), ex.JSONError())
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}
			return
		}

		budsInfo, err := users.GetUser(identifier, db, middleware.ExtractPublicKey(c), dynamicLinkServiceUrlChan, redisCache)
		if err != nil {
			var ex bantupayErrors.GenericError
			var ok bool

			ex, ok = err.(bantupayErrors.GenericError)
			if ok {
				c.JSON(ex.HTTPCode(), ex.JSONError())
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}
			return
		}
		cacheKey := fmt.Sprintf("[GET] /v2/users/%v", identifier)
		paymentHistoryCacheKey := fmt.Sprintf("[GET] /v2/users/%v/payments", identifier)
		senderPaymentHistoryCacheKey := fmt.Sprintf("[GET] /v2/users/%v/payments", identifier)
		senderCacheKey := fmt.Sprintf("[GET] /v2/users/%v", identifier)

		redisCache.InvalidateCachedHttpResponse(senderCacheKey, senderPaymentHistoryCacheKey)

		redisCache.InvalidateCachedHttpResponse(cacheKey, paymentHistoryCacheKey)

		c.JSON(http.StatusOK, budsInfo)

	})

	//update mobile number
	router.PUT("/v2/users/:targetUser/actions/update-mobile/:mobile", middleware.AuthenticationMiddleware(), func(c *gin.Context) {
		var err error

		identifier := strings.TrimSpace(strings.ToLower(c.Param("targetUser")))
		mobile := strings.TrimSpace(c.Param("mobile"))
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
		conDB.PrintDBStats(fmt.Sprintf("PUT /v2/users/%v/actions/update-mobile/%v", identifier, mobile), db)

		// var err error

		userInfo, err := usersDB.GetUserInfo(identifier, db)
		if err != nil {
			var ex bantupayErrors.GenericError
			var ok bool

			ex, ok = err.(bantupayErrors.GenericError)
			if ok {
				c.JSON(ex.HTTPCode(), ex.JSONError())
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}
			return
		}

		err = users.UpdatePhoneNumber(&userInfo, mobile, db)

		if err != nil {
			var ex bantupayErrors.GenericError
			var ok bool

			ex, ok = err.(bantupayErrors.GenericError)
			if ok {
				c.JSON(ex.HTTPCode(), ex.JSONError())
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}
			return
		}

		budsInfo, err := users.GetUser(identifier, db, middleware.ExtractPublicKey(c), dynamicLinkServiceUrlChan, redisCache)
		if err != nil {
			var ex bantupayErrors.GenericError
			var ok bool

			ex, ok = err.(bantupayErrors.GenericError)
			if ok {
				c.JSON(ex.HTTPCode(), ex.JSONError())
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}
			return
		}
		cacheKey := fmt.Sprintf("[GET] /v2/users/%v", identifier)
		paymentHistoryCacheKey := fmt.Sprintf("[GET] /v2/users/%v/payments", identifier)
		senderPaymentHistoryCacheKey := fmt.Sprintf("[GET] /v2/users/%v/payments", identifier)
		senderCacheKey := fmt.Sprintf("[GET] /v2/users/%v", identifier)

		redisCache.InvalidateCachedHttpResponse(senderCacheKey, senderPaymentHistoryCacheKey)

		redisCache.InvalidateCachedHttpResponse(cacheKey, paymentHistoryCacheKey)

		c.JSON(http.StatusOK, budsInfo)

	})

	router.PUT("/v2/users/:targetUser/actions/request-mobile-otp", middleware.AuthenticationMiddleware(), func(c *gin.Context) {
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
		conDB.PrintDBStats(fmt.Sprintf("PUT /v2/users/%v/actions/request-mobile-otp", identifier), db)

		// var err error

		userInfo, err := usersDB.GetUserInfo(identifier, db)
		if err != nil {
			var ex bantupayErrors.GenericError
			var ok bool

			ex, ok = err.(bantupayErrors.GenericError)
			if ok {
				c.JSON(ex.HTTPCode(), ex.JSONError())
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}
			return
		}

		err = users.SendPhoneVerificationCode(&userInfo, db)

		if err != nil {
			var ex bantupayErrors.GenericError
			var ok bool

			ex, ok = err.(bantupayErrors.GenericError)
			if ok {
				c.JSON(ex.HTTPCode(), ex.JSONError())
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}
			return
		}

		budsInfo, err := users.GetUser(identifier, db, middleware.ExtractPublicKey(c), dynamicLinkServiceUrlChan, redisCache)
		if err != nil {
			var ex bantupayErrors.GenericError
			var ok bool

			ex, ok = err.(bantupayErrors.GenericError)
			if ok {
				c.JSON(ex.HTTPCode(), ex.JSONError())
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}
			return
		}
		cacheKey := fmt.Sprintf("[GET] /v2/users/%v", identifier)
		paymentHistoryCacheKey := fmt.Sprintf("[GET] /v2/users/%v/payments", identifier)
		senderPaymentHistoryCacheKey := fmt.Sprintf("[GET] /v2/users/%v/payments", identifier)
		senderCacheKey := fmt.Sprintf("[GET] /v2/users/%v", identifier)

		redisCache.InvalidateCachedHttpResponse(senderCacheKey, senderPaymentHistoryCacheKey)

		redisCache.InvalidateCachedHttpResponse(cacheKey, paymentHistoryCacheKey)

		c.JSON(http.StatusOK, budsInfo)

	})

	router.GET("/v2/users/:targetUser/generate/payment", middleware.AuthenticationMiddleware(), func(c *gin.Context) {
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

		cacheKey := fmt.Sprintf("[GET] /v2/users/%v/generate/payment", identifier)
		cacheKeyParameters := fmt.Sprintf("%v", c.Request.URL.RawQuery)

		{
			//search cache

			// cacheKeyParameters := fmt.Sprintf("limit=%v&order=%v&cursor=%v&forTransactionHash=%v&includeHash=%v&temp=%v", limit, orderStr, cursor, forTransactionHash, includeHash, temp)

			ok, status, response := redisCache.CachedHttpResponseWithParameters(cacheKey, cacheKeyParameters)

			if ok {
				log.Printf("[%v]/[%v], served from cache\n", cacheKey, cacheKeyParameters)
				c.JSON(status, response)
				return
			}
		}

		conDB.PrintDBStats(fmt.Sprintf("GET /v2/users/%v/generate/payment?paymentDestination=%v&assetCode=%v&assetIssuer=%v&amount=%v&memo=%v", identifier, paymentDestination, assetCode, assetIssuer, amount, memo), db)

		_, err = usersDB.GetUserInfo(identifier, db)

		if err != nil {
			log.Println("[GET UserInfo] error for user:", identifier, "error: ", err)

			var ex bantupayErrors.GenericError
			var ok bool

			ex, ok = err.(bantupayErrors.GenericError)
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
			cacheDurationInSeconds := 525600 * 3 * 60 //3yrs

			redisCache.CacheHttpResponseWithParameters(cacheKey, cacheKeyParameters, http.StatusOK, data, cacheDurationInSeconds)
		}
		c.JSON(http.StatusOK, data)

	})

}
