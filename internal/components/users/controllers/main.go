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

	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
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
	router.GET("/v1/stream/ws/:targetUser", func(c *gin.Context) {

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

	router.GET("/v1/users/payments/:targetPublicKeyForHistory", middleware.AuthenticationMiddlewareUsingTimestamp(), func(c *gin.Context) {
		// var err error

		targetPublicKeyForHistory := strings.TrimSpace(strings.ToUpper(c.Param("targetPublicKeyForHistory")))
		cacheKey := fmt.Sprintf("[GET] /v1/users/payments/%v", targetPublicKeyForHistory)
		cacheKeyParameters := c.Request.URL.RequestURI()
		{
			// check cache
			ok, status, response := gc.RedisCache.CachedHttpResponseWithParameters(cacheKey, cacheKeyParameters)

			if ok {
				log.Printf("[%v]/[%v], served from cache\n", cacheKey, cacheKeyParameters)
				c.JSON(status, response)
				return
			}

		}
		// cacheDurationInSeconds := 1 * 60 //1 minutes
		cacheDurationInSeconds := 60 //1 minutes
		conDB.PrintDBStats(fmt.Sprintf("/v1/users/payments/%v", targetPublicKeyForHistory), gc.DB)

		signerUser, err := usersDB.GetUser(middleware.ExtractSigner(c), gc.DB)

		if err != nil {
			log.Println("[GET USER] error for signer:", middleware.ExtractSigner(c), "error: ", err)

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
			gc.RedisCache.CacheHttpResponseWithParameters(cacheKey, cacheKeyParameters, statusCode, response, cacheDurationInSeconds)
			return
		}

		{
			// gc.RedisCache.InvalidateCachedHttpResponse(cacheKey)

			//check if the owner is the one accessing it or if the one accessing it has access to access it.

			if signerUser.PrimarySigner != middleware.ExtractSigner(c) && !userServices.HasAccessToPublicKey(signerUser.PublicKey, targetPublicKeyForHistory, gc) {
				te := &tErrors.ErrorInvalidAuthorization{}

				log.Println("[GET HISTORY] Invalid signer for user:", signerUser.Username, "error: ", err)
				c.JSON(te.HTTPCode(), te.JSONError())
				return
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
		pnt := c.Query("pnt")

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
					//update the push notification token, if it is different
					if len(pnt) > 10 {
						usersDB.UpdatePushNotificationToken(v.ID, &pnt, gc.DB)
						userInfo.UserData.PushNotificationToken = pnt

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

		data, _ := io.ReadAll(c.Request.Body)

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

		data, _ := io.ReadAll(c.Request.Body)

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
			pns.SendFirebaseMessage(*user.PushNotificationToken, "New Sub-wallet Added!", fmt.Sprintf("You have successfully added a new sub wallet tagged [%v].", returnedSubwalletInfo.WalletTag), "", dataPayload, gc.PushNotificationClient, gc.PNSContext)
		}

		//At this point, there was no error.

		c.JSON(http.StatusOK, returnedSubwalletInfo)
	})

	router.PUT("/v1/users/upload-picture", middleware.AuthenticationMiddlewareUsingTimestamp(), func(c *gin.Context) {
		var err error
		f, err := c.FormFile("profilePicture")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if f.Size > 700000 {
			//greater than 700kb
			c.JSON(http.StatusBadRequest, gin.H{"error": "picture cannot be more than 700kb in file size"})
			return
		}
		blobFile, err := f.Open()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "error attempting to validate the picture uploaded"})

			return
		}
		fnameSplit := strings.Split(f.Filename, ".")
		fileExtension := fnameSplit[len(fnameSplit)-1]

		{
			//check for unsupported extension
			if !strings.EqualFold(fileExtension, "jpg") && !strings.EqualFold(fileExtension, "jpeg") && !strings.EqualFold(fileExtension, "png") && !strings.EqualFold(fileExtension, "gif") {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Unsurported picture format. Only jpg, jpeg, png and gif are supported"})

				return
			}
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

		conDB.PrintDBStats(fmt.Sprintf("PUT /v1/users/upload-picture %v", user.Username), gc.DB)

		url, err := userServices.UploadProfilePicture(&user, blobFile, fmt.Sprintf("%s.%s", user.Username, fileExtension), gc)

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

		if user.PushNotificationToken != nil && len(url) > 0 {
			dataPayload := make(map[string]string)
			dataPayload["route"] = ""
			pns.SendFirebaseMessage(*user.PushNotificationToken, "Profile picture updated!", fmt.Sprintf("You have successfully updated profile picture on your account [%v].", user.Username), url, dataPayload, gc.PushNotificationClient, gc.PNSContext)
		}

		userCacheKey := fmt.Sprintf("[GET] /v1/users/%v", user.Username)

		gc.RedisCache.InvalidateCachedHttpResponse(userCacheKey)

		//At this point, there was no error.

		c.JSON(http.StatusOK, url)
	})

	router.POST("/v1/users/trust-asset", middleware.AuthenticationMiddlewareUsingTimestamp(), func(c *gin.Context) {
		var err error

		var trustLineInfo userModels.Trustline
		// var err error

		data, _ := io.ReadAll(c.Request.Body)

		err = json.Unmarshal(data, &trustLineInfo)

		var invalidJSON tErrors.ErrorInvalidJSON

		if err != nil {
			c.JSON(http.StatusBadRequest, invalidJSON.JSONError())
			return
		}

		signerUser, err := usersDB.GetUser(middleware.ExtractSigner(c), gc.DB)

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
		walletOwner, err := usersDB.GetUser(middleware.ExtractPublicKey(c), gc.DB)

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
		wallet, _, err := usersDB.GetWallet(middleware.ExtractPublicKey(c), gc.DB)

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

		conDB.PrintDBStats(fmt.Sprintf("POST /v1/users/trust-asset %v", wallet.Alias), gc.DB)

		returnedTrustLineInfo, err := userServices.TrustAsset(&signerUser, &wallet, &trustLineInfo, gc)

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

		if walletOwner.PushNotificationToken != nil && len(returnedTrustLineInfo.TransactionID) > 0 && returnedTrustLineInfo.TransactionID != "PENDING_AUTH" {
			dataPayload := make(map[string]string)
			dataPayload["route"] = ""
			pns.SendFirebaseMessage(*walletOwner.PushNotificationToken, fmt.Sprintf("Asset %v Accepted on %v!", trustLineInfo.AssetCode, wallet.Alias), fmt.Sprintf("You have successfully added the asset [%v] to the list of your trusted assets that you can receive on the wallet with alias [%v].", returnedTrustLineInfo.AssetCode, wallet.Alias), "", dataPayload, gc.PushNotificationClient, gc.PNSContext)
		}

		//At this point, there was no error.

		c.JSON(http.StatusOK, returnedTrustLineInfo)
	})

	router.POST("/v1/users/remove-asset", middleware.AuthenticationMiddlewareUsingTimestamp(), func(c *gin.Context) {
		var err error

		var trustLineInfo userModels.Trustline
		// var err error

		data, _ := io.ReadAll(c.Request.Body)

		err = json.Unmarshal(data, &trustLineInfo)

		var invalidJSON tErrors.ErrorInvalidJSON

		if err != nil {
			c.JSON(http.StatusBadRequest, invalidJSON.JSONError())
			return
		}

		signerUser, err := usersDB.GetUser(middleware.ExtractSigner(c), gc.DB)

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
		walletOwner, err := usersDB.GetUser(middleware.ExtractPublicKey(c), gc.DB)

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
		wallet, _, err := usersDB.GetWallet(middleware.ExtractPublicKey(c), gc.DB)

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

		conDB.PrintDBStats(fmt.Sprintf("POST /v1/users/remove-asset %v", wallet.Alias), gc.DB)

		returnedTrustLineInfo, err := userServices.RemoveAssetTrust(&signerUser, &wallet, &trustLineInfo, gc)

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

		if walletOwner.PushNotificationToken != nil && len(returnedTrustLineInfo.TransactionID) > 0 && returnedTrustLineInfo.TransactionID != "PENDING_AUTH" {
			dataPayload := make(map[string]string)
			dataPayload["route"] = ""
			pns.SendFirebaseMessage(*walletOwner.PushNotificationToken, fmt.Sprintf("Asset %v Removed from %v!", trustLineInfo.AssetCode, wallet.Alias), fmt.Sprintf("You have successfully removed the asset [%v] from the list of your trusted assets that you can receive on the wallet with alias [%v].", returnedTrustLineInfo.AssetCode, wallet.Alias), "", dataPayload, gc.PushNotificationClient, gc.PNSContext)
		}

		//At this point, there was no error.

		c.JSON(http.StatusOK, returnedTrustLineInfo)
	})

	router.PUT("/v1/users/actions/claim-asset", middleware.AuthenticationMiddlewareUsingTimestamp(), func(c *gin.Context) {
		var err error

		_, err = usersDB.GetUser(middleware.ExtractSigner(c), gc.DB)

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
		walletOwner, err := usersDB.GetUser(middleware.ExtractPublicKey(c), gc.DB)

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
		wallet, _, err := usersDB.GetWallet(middleware.ExtractPublicKey(c), gc.DB)

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

		conDB.PrintDBStats(fmt.Sprintf("PUT /v1/users/actions/claim-asset %v", middleware.ExtractPublicKey(c)), gc.DB)

		var pendingAssetToClaim userModels.PendingAssetToClaim
		// var err error

		data, _ := io.ReadAll(c.Request.Body)

		err = json.Unmarshal(data, &pendingAssetToClaim)

		var invalidJSON tErrors.ErrorInvalidJSON

		if err != nil {
			c.JSON(http.StatusBadRequest, invalidJSON.JSONError())
			return
		}

		_, complete, err := userServices.ClaimPendingAsset(&walletOwner, &wallet, &pendingAssetToClaim, gc.DB)

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

		var ownerBalanceCacheKey, tempCacheKey, sNFT string

		ownerBalanceCacheKey = fmt.Sprintf("GetBalance_%s", middleware.ExtractPublicKey(c))
		sNFT = fmt.Sprintf("GetNFTs_%s", middleware.ExtractPublicKey(c))

		tempCacheKey = fmt.Sprintf("GetBalance_%s", *wallet.TempPublicKey)

		userCacheKey := fmt.Sprintf("[GET] /v1/users/%v", walletOwner.Username)
		paymentPaymentHistoryCacheKey := fmt.Sprintf("[GET] /v1/users/payments/%v", middleware.ExtractPublicKey(c))

		gc.RedisCache.InvalidateCachedHttpResponse(ownerBalanceCacheKey, tempCacheKey, userCacheKey, paymentPaymentHistoryCacheKey, sNFT)
		log.Printf("[CLAIM ASSET] Transaction Signature: [%v]\n", pendingAssetToClaim.TransactionSignature)
		if complete {
			c.JSON(http.StatusOK, pendingAssetToClaim)
		} else {
			c.JSON(http.StatusAccepted, pendingAssetToClaim)
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
