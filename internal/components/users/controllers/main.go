package users

import (
	"os"
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
	"github.com/stellar/go/keypair"
)

// Init initializes /v1/users endpoint
func Init(router *gin.Engine, callBackRetryChan chan userModels.RetryCallbacks, gc *sharedconfig.GlobalConfig) {
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

		_, err := keypair.ParseAddress(targetPublicKeyForHistory)
		if err != nil {

			statusCode := http.StatusBadRequest
			response := gin.H{"error": "error invalid address", "message": "Only valid addresses are allowed"}

			c.JSON(statusCode, response)
			return
		}
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
		cacheDurationInSeconds := 20 //in seconds
		conDB.PrintDBStats(fmt.Sprintf("/v1/users/payments/%v", targetPublicKeyForHistory), gc.DB)

		signerUser, err := usersDB.GetUserFromPrimarySigner(middleware.ExtractSigner(c), gc.DB, gc)

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
				response = gin.H{"error": err.Error(), "message": err.Error()}
			}

			c.JSON(statusCode, response)
			gc.RedisCache.CacheHttpResponseWithParameters(cacheKey, cacheKeyParameters, statusCode, response, cacheDurationInSeconds)
			return
		}
		wallet, temp, err := usersDB.GetWallet(targetPublicKeyForHistory, gc.DB)

		if err != nil {
			log.Println("[GET TARGET USER] error for PUBLIC KEY:", targetPublicKeyForHistory, "error: ", err)

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
			// gc.RedisCache.CacheHttpResponseWithParameters(cacheKey, cacheKeyParameters, statusCode, response, cacheDurationInSeconds)
			return
		}
		if temp {

			statusCode := http.StatusBadRequest
			response := gin.H{"error": "error only main wallets allowed", "message": "Only main wallets are allowed. The address you provided is not a main wallet."}

			c.JSON(statusCode, response)
			return
		}

		targetOwnerUser, err := usersDB.GetUser(targetPublicKeyForHistory, gc.DB, gc)

		if err != nil {
			log.Println("[GET TARGET USER] error for PUBLIC KEY:", targetPublicKeyForHistory, "error: ", err)

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
			// gc.RedisCache.CacheHttpResponseWithParameters(cacheKey, cacheKeyParameters, statusCode, response, cacheDurationInSeconds)
			return
		}

		{
			// gc.RedisCache.InvalidateCachedHttpResponse(cacheKey)

			//check if the owner is the one accessing it or if the one accessing it has access to access it.

			if (signerUser.Username != targetOwnerUser.Username) && !wallet.SignerHasAccess(&signerUser, gc) {
				te := &tErrors.ErrorInvalidAuthorization{}

				log.Println("[GET HISTORY] Invalid access for user:", signerUser.Username, "error: ", te.Error())
				c.JSON(te.HTTPCode(), te.JSONError())
				return
			}

		}
		//Get Payment history
		historyRecords := paymentServices.GetPaymentHistory(targetPublicKeyForHistory, gc, c)

		c.JSON(http.StatusOK, historyRecords)
		// gc.RedisCache.CacheHttpResponse(cacheKey, http.StatusOK, historyRecords, cacheDurationInSeconds)
		gc.RedisCache.CacheHttpResponseWithParameters(cacheKey, cacheKeyParameters, http.StatusOK, historyRecords, cacheDurationInSeconds)

	})

	router.GET("/v1/users/:targetUser", middleware.AuthenticationMiddlewareUsingTimestamp(), func(c *gin.Context) {
		// var err error//true-client-ip
		if os.Getenv("LOG_IP_ADDRESS") == "1" {
			log.Printf("<<<<<<<<<<<<<>>>>>>>>>>>IP address: %v\nCLIENTIP: %v\nTrue CLient IP: %v", c.GetHeader(strings.ToUpper("x-forwarded-for")), c.ClientIP(), c.GetHeader(strings.ToUpper("true-client-ip")))
			// log.Printf("<<<<<<<<<<<<<<>>>>>>>>>>>%+v\n", c)
		}
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
		if queryType == "import" {
			u, e := usersDB.GetUser(identifier, gc.DB, gc)
			if e == nil {
				u.InvalidateUserCache(gc)
			}
		}
		userInfo, err := userServices.GetUserInfo(identifier, middleware.ExtractSigner(c), gc)

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
				response = gin.H{"error": err.Error(), "message": err.Error()}
			}
			if queryType == "import" {
				c.JSON(http.StatusNotFound, response)
				return
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
						usersDB.UpdatePushNotificationToken(v.ID, &pnt, gc.DB, gc)
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

		user, err := usersDB.GetUserFromPrimarySigner(middleware.ExtractSigner(c), gc.DB, gc)

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
		// io.ReadAll(blobFile)
		fnameSplit := strings.Split(f.Filename, ".")
		fileExtension := fnameSplit[len(fnameSplit)-1]

		{
			//check for unsupported extension
			if !strings.EqualFold(fileExtension, "jpg") && !strings.EqualFold(fileExtension, "jpeg") && !strings.EqualFold(fileExtension, "png") && !strings.EqualFold(fileExtension, "gif") {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Unsurported picture format. Only jpg, jpeg, png and gif are supported"})

				return
			}
		}

		user, err := usersDB.GetUserFromPrimarySigner(middleware.ExtractSigner(c), gc.DB, gc)

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

	router.POST("/v1/users/asset/opt-in", middleware.AuthenticationMiddlewareUsingTimestamp(), func(c *gin.Context) {
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

		signerUser, err := usersDB.GetUserFromPrimarySigner(middleware.ExtractSigner(c), gc.DB, gc)

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
		walletOwner, err := usersDB.GetUser(middleware.ExtractPublicKey(c), gc.DB, gc)

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
		if wallet.WalletType != 0 {
			c.JSON(http.StatusForbidden, gin.H{"error": "error-wallet-type-forbidden", "message": "Operation not allowed on any special type of wallets. Only standard wallets are allowed."})
			return
		}
		if signerUser.PrimarySigner != wallet.Signer {
			c.JSON(http.StatusForbidden, gin.H{"error": "error-unauthorized-access", "message": "You do not have permission on this wallet."})
			return
		}
		conDB.PrintDBStats(fmt.Sprintf("POST /v1/users/opt-in %v", wallet.Alias), gc.DB)

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
			pns.SendFirebaseMessage(*walletOwner.PushNotificationToken, fmt.Sprintf("Asset %v opted in on %v!", trustLineInfo.AssetCode, wallet.Alias), fmt.Sprintf("You have successfully added the asset [%v] to the list of your trusted assets that you can receive on the wallet with alias [%v].", returnedTrustLineInfo.AssetCode, wallet.Alias), "", dataPayload, gc.PushNotificationClient, gc.PNSContext)
		}
		walletOwner.InvalidateUserCache(gc)
		//At this point, there was no error.

		c.JSON(http.StatusOK, returnedTrustLineInfo)
	})

	router.POST("/v1/shared-access/users/asset/opt-in", middleware.AuthenticationMiddlewareUsingTimestamp(), func(c *gin.Context) {
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

		signerUser, err := usersDB.GetUserFromPrimarySigner(middleware.ExtractSigner(c), gc.DB, gc)

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
		{
			//check if pending shared access modify exists
			if userServices.CheckPendingSharedAccessApproval(middleware.ExtractPublicKey(c), gc.DB) {
				c.JSON(http.StatusForbidden, gin.H{"error": "error-pending-shared-access-op", "message": "There is a pending shared access operation on this wallet and must be completed first before attempting to send payment from this wallet."})
				return
			}
		}
		if wallet.WalletType != 0 {
			c.JSON(http.StatusForbidden, gin.H{"error": "error-wallet-type-forbidden", "message": "Operation not allowed on any special type of wallets. Only standard wallets are allowed."})
			return
		}

		hasInitiatorAccess := false
		isViewOnly := wallet.HasViewOnlyAccess(gc)

		if isViewOnly && wallet.Signer != signerUser.PrimarySigner {
			c.JSON(http.StatusForbidden, gin.H{"error": "error-unauthorized-access", "message": "You do not have permission on this wallet."})
			return
		}
		// check if user has initiator access to wallet.
		if wallet.SharedAccessEnabled == 1 && !isViewOnly {
			// check if user has initiator access to wallet.
			for _, p := range signerUser.WalletsSharedWithUser {
				if p.WalletPublicKey == middleware.ExtractPublicKey(c) && p.TargetUsername == signerUser.Username && p.Permission == "INITIATOR" {
					hasInitiatorAccess = true
				}
			}
			if !hasInitiatorAccess && !wallet.HasViewOnlyAccess(gc) {
				c.JSON(http.StatusForbidden, gin.H{"error": "error-unauthorized-access", "message": "You do not have an initiator permission on this wallet."})
				return
			}
		}

		conDB.PrintDBStats(fmt.Sprintf("POST /v1/shared-access/users/asset/opt-in %v", wallet.Alias), gc.DB)

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

		if returnedTrustLineInfo.Commit == 0 {
			c.JSON(http.StatusAccepted, returnedTrustLineInfo)

			return
		}

		//At this point, there was no error.

		c.JSON(http.StatusOK, returnedTrustLineInfo)
		if returnedTrustLineInfo.TransactionID == "PENDING_AUTH" {
			//start push notificationMessage
			notificationList := make(map[string]string)
			permissionList := wallet.Permissions
			for _, v := range permissionList {
				u, e := userModels.Username(v.TargetUsername).GetSimpleUser(gc.DB, gc)
				if e != nil {
					continue
				}
				if u.PushNotificationToken == nil {
					continue
				}

				if _, ok := notificationList[*u.PushNotificationToken]; ok {
					continue
				}

				dataPayload := make(map[string]string)
				dataPayload["route"] = "pendingApproval"

				u.SendPushMessage(fmt.Sprintf("%v opt-in request from %v!", trustLineInfo.AssetCode, wallet.Alias), fmt.Sprintf("Request: %v", returnedTrustLineInfo.ReturnedDescription), "", dataPayload, gc)
				notificationList[*u.PushNotificationToken] = v.TargetUsername

			}
		}
	})

	router.DELETE("/v1/users/asset/opt-out", middleware.AuthenticationMiddlewareUsingTimestamp(), func(c *gin.Context) {
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

		signerUser, err := usersDB.GetUserFromPrimarySigner(middleware.ExtractSigner(c), gc.DB, gc)

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
		walletOwner, err := usersDB.GetUser(middleware.ExtractPublicKey(c), gc.DB, gc)

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

		if walletOwner.Username != signerUser.Username {
			c.JSON(http.StatusForbidden, gin.H{"error": "error-access-forbidden", "message": "Access forbidden. Wallet does not belong to you."})
			return
		}
		if wallet.WalletType != 0 {
			c.JSON(http.StatusForbidden, gin.H{"error": "error-wallet-type-forbidden", "message": "Operation not allowed on any special type of wallets. Only standard wallets are allowed."})
			return
		}
		conDB.PrintDBStats(fmt.Sprintf("DELETE /v1/users/asset/opt-out %v", wallet.Alias), gc.DB)

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
			pns.SendFirebaseMessage(*walletOwner.PushNotificationToken, fmt.Sprintf("Opt-out asset %v on %v!", trustLineInfo.AssetCode, wallet.Alias), fmt.Sprintf("You have successfully opted-out of the asset [%v] from the list of your trusted assets that you can receive on the wallet with alias [%v].", returnedTrustLineInfo.AssetCode, wallet.Alias), "", dataPayload, gc.PushNotificationClient, gc.PNSContext)
		}
		walletOwner.InvalidateUserCache(gc)
		//At this point, there was no error.

		c.JSON(http.StatusOK, returnedTrustLineInfo)
	})

	router.POST("/v1/shared-access/users/asset/opt-out", middleware.AuthenticationMiddlewareUsingTimestamp(), func(c *gin.Context) {
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

		signerUser, err := usersDB.GetUserFromPrimarySigner(middleware.ExtractSigner(c), gc.DB, gc)

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
		{
			//check if pending shared access modify exists
			if userServices.CheckPendingSharedAccessApproval(middleware.ExtractPublicKey(c), gc.DB) {
				c.JSON(http.StatusForbidden, gin.H{"error": "error-pending-shared-access-op", "message": "There is a pending shared access operation on this wallet and must be completed first before attempting to send payment from this wallet."})
				return
			}
		}
		if wallet.WalletType != 0 {
			c.JSON(http.StatusForbidden, gin.H{"error": "error-wallet-type-forbidden", "message": "Operation not allowed on any special type of wallets. Only standard wallets are allowed."})
			return
		}

		hasInitiatorAccess := false
		isViewOnly := wallet.HasViewOnlyAccess(gc)

		if isViewOnly && wallet.Signer != signerUser.PrimarySigner {
			c.JSON(http.StatusForbidden, gin.H{"error": "error-unauthorized-access", "message": "You do not have permission on this wallet."})
			return
		}
		// check if user has initiator access to wallet.
		if wallet.SharedAccessEnabled == 1 && !isViewOnly {
			// check if user has initiator access to wallet.
			for _, p := range signerUser.WalletsSharedWithUser {
				if p.WalletPublicKey == middleware.ExtractPublicKey(c) && p.TargetUsername == signerUser.Username && p.Permission == "INITIATOR" {
					hasInitiatorAccess = true
				}
			}
			if !hasInitiatorAccess && !userModels.UserWalletID(middleware.ExtractPublicKey(c)).PublicKeyHasViewOnlyAccess(gc) {
				c.JSON(http.StatusForbidden, gin.H{"error": "error-unauthorized-access", "message": "You do not have an initiator permission on this wallet."})
				return
			}
		}

		conDB.PrintDBStats(fmt.Sprintf("POST /v1/shared-access/users/asset/opt-out %v", wallet.Alias), gc.DB)

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
		if returnedTrustLineInfo.Commit == 0 {
			c.JSON(http.StatusAccepted, returnedTrustLineInfo)
			return
		}
		//At this point, there was no error.

		c.JSON(http.StatusOK, returnedTrustLineInfo)
		if returnedTrustLineInfo.TransactionID == "PENDING_AUTH" {
			//start push notificationMessage
			notificationList := make(map[string]string)
			permissionList := wallet.Permissions
			for _, v := range permissionList {
				u, e := userModels.Username(v.TargetUsername).GetSimpleUser(gc.DB, gc)
				if e != nil {
					continue
				}
				if u.PushNotificationToken == nil {
					continue
				}

				if _, ok := notificationList[*u.PushNotificationToken]; ok {
					continue
				}

				dataPayload := make(map[string]string)
				dataPayload["route"] = "pendingApproval"

				u.SendPushMessage(fmt.Sprintf("%v opt-out request from %v!", trustLineInfo.AssetCode, wallet.Alias), fmt.Sprintf("Request: %v", returnedTrustLineInfo.ReturnedDescription), "", dataPayload, gc)
				notificationList[*u.PushNotificationToken] = v.TargetUsername
			}
		}

	})

	router.PUT("/v1/users/actions/claim-asset", middleware.AuthenticationMiddlewareUsingTimestamp(), func(c *gin.Context) {
		var err error

		signerUser, err := usersDB.GetUserFromPrimarySigner(middleware.ExtractSigner(c), gc.DB, gc)

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
		walletOwner, err := usersDB.GetUser(middleware.ExtractPublicKey(c), gc.DB, gc)

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
		if walletOwner.Username != signerUser.Username {
			c.JSON(http.StatusForbidden, gin.H{"error": "error-access-forbidden", "message": "Access forbidden. Wallet does not belong to you."})
			return
		}
		if wallet.WalletType != 0 {
			c.JSON(http.StatusForbidden, gin.H{"error": "error-wallet-type-forbidden", "message": "Operation not allowed on any special type of wallets. Only standard wallets are allowed."})
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

		_, complete, err := userServices.ClaimPendingAsset(&signerUser, &wallet, &pendingAssetToClaim, gc)

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
			dataPayload := make(map[string]string)
			dataPayload["none"] = ""
			if len(pendingAssetToClaim.TransactionID) > 0 {
				walletOwner.SendPushMessage(fmt.Sprintf("%v pending balance on wallet %v has been claimed!", pendingAssetToClaim.AssetCode, wallet.Alias), fmt.Sprintf("%v pending balance claimed!", pendingAssetToClaim.AssetCode), "", dataPayload, gc)
			}
			walletOwner.InvalidateUserCache(gc)
		} else {
			c.JSON(http.StatusAccepted, pendingAssetToClaim)
		}

	})

	router.DELETE("/v1/users/actions/reject-asset", middleware.AuthenticationMiddlewareUsingTimestamp(), func(c *gin.Context) {
		var err error

		signerUser, err := usersDB.GetUserFromPrimarySigner(middleware.ExtractSigner(c), gc.DB, gc)

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
		walletOwner, err := usersDB.GetUser(middleware.ExtractPublicKey(c), gc.DB, gc)

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
		if walletOwner.Username != signerUser.Username {
			c.JSON(http.StatusForbidden, gin.H{"error": "error-access-forbidden", "message": "Access forbidden. Wallet does not belong to you."})
			return
		}
		if wallet.WalletType != 0 {
			c.JSON(http.StatusForbidden, gin.H{"error": "error-wallet-type-forbidden", "message": "Operation not allowed on any special type of wallets. Only standard wallets are allowed."})
			return
		}
		conDB.PrintDBStats(fmt.Sprintf("DELETE /v1/users/actions/reject-asset %v", middleware.ExtractPublicKey(c)), gc.DB)

		var pendingAssetToClaim userModels.PendingAssetToClaim
		// var err error

		data, _ := io.ReadAll(c.Request.Body)

		err = json.Unmarshal(data, &pendingAssetToClaim)

		var invalidJSON tErrors.ErrorInvalidJSON

		if err != nil {
			c.JSON(http.StatusBadRequest, invalidJSON.JSONError())
			return
		}

		_, complete, err := userServices.RejectPendingAsset(&signerUser, &wallet, &pendingAssetToClaim, gc)

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
		log.Printf("[REJECT ASSET] Transaction Signature: [%v]\n", pendingAssetToClaim.TransactionSignature)
		if complete {
			c.JSON(http.StatusOK, pendingAssetToClaim)
			dataPayload := make(map[string]string)
			dataPayload["none"] = ""
			if len(pendingAssetToClaim.TransactionID) > 0 {
				walletOwner.SendPushMessage(fmt.Sprintf("%v pending balance rejected on wallet %v!", pendingAssetToClaim.AssetCode, wallet.Alias), fmt.Sprintf("%v pending balance rejected", pendingAssetToClaim.AssetCode), "", dataPayload, gc)
			}

			walletOwner.InvalidateUserCache(gc)

		} else {
			c.JSON(http.StatusAccepted, pendingAssetToClaim)
		}

	})

	router.PUT("/v1/shared-access/users/actions/claim-asset", middleware.AuthenticationMiddlewareUsingTimestamp(), func(c *gin.Context) {
		var err error

		signerUser, err := usersDB.GetUserFromPrimarySigner(middleware.ExtractSigner(c), gc.DB, gc)

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
		{
			//check if pending shared access modify exists
			if userServices.CheckPendingSharedAccessApproval(middleware.ExtractPublicKey(c), gc.DB) {
				c.JSON(http.StatusForbidden, gin.H{"error": "error-pending-shared-access-op", "message": "There is a pending shared access operation on this wallet and must be completed first before attempting to send payment from this wallet."})
				return
			}
		}
		if wallet.WalletType != 0 {
			c.JSON(http.StatusForbidden, gin.H{"error": "error-wallet-type-forbidden", "message": "Operation not allowed on any special type of wallets. Only standard wallets are allowed."})
			return
		}

		hasInitiatorAccess := false
		isViewOnly := wallet.HasViewOnlyAccess(gc)

		if isViewOnly && wallet.Signer != signerUser.PrimarySigner {
			c.JSON(http.StatusForbidden, gin.H{"error": "error-unauthorized-access", "message": "You do not have permission on this wallet."})
			return
		}
		// check if user has initiator access to wallet.
		if wallet.SharedAccessEnabled == 1 && !isViewOnly {
			// check if user has initiator access to wallet.

			for _, p := range signerUser.WalletsSharedWithUser {
				if p.WalletPublicKey == middleware.ExtractPublicKey(c) && p.TargetUsername == signerUser.Username && p.Permission == "INITIATOR" {
					hasInitiatorAccess = true
				}
			}
			if !hasInitiatorAccess && !userModels.UserWalletID(middleware.ExtractPublicKey(c)).PublicKeyHasViewOnlyAccess(gc) {
				c.JSON(http.StatusForbidden, gin.H{"error": "error-unauthorized-access", "message": "You do not have an initiator permission on this wallet."})
				return
			}
		}

		conDB.PrintDBStats(fmt.Sprintf("PUT /v1/shared-access/users/actions/claim-asset %v", middleware.ExtractPublicKey(c)), gc.DB)

		var pendingAssetToClaim userModels.PendingAssetToClaim
		// var err error

		data, _ := io.ReadAll(c.Request.Body)

		err = json.Unmarshal(data, &pendingAssetToClaim)

		var invalidJSON tErrors.ErrorInvalidJSON

		if err != nil {
			c.JSON(http.StatusBadRequest, invalidJSON.JSONError())
			return
		}

		_, complete, err := userServices.ClaimPendingAsset(&signerUser, &wallet, &pendingAssetToClaim, gc)

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

		userCacheKey := fmt.Sprintf("[GET] /v1/users/%v", wallet.Alias)
		paymentPaymentHistoryCacheKey := fmt.Sprintf("[GET] /v1/users/payments/%v", middleware.ExtractPublicKey(c))

		gc.RedisCache.InvalidateCachedHttpResponse(ownerBalanceCacheKey, tempCacheKey, userCacheKey, paymentPaymentHistoryCacheKey, sNFT)
		log.Printf("[CLAIM ASSET] Transaction Signature: [%v]\n", pendingAssetToClaim.TransactionSignature)
		if complete {

			c.JSON(http.StatusOK, pendingAssetToClaim)
			{
				//start push notificationMessage
				notificationList := make(map[string]string)
				permissionList := wallet.Permissions
				for _, v := range permissionList {
					u, e := userModels.Username(v.TargetUsername).GetSimpleUser(gc.DB, gc)
					if e != nil {
						continue
					}

					dataPayload := make(map[string]string)
					dataPayload["route"] = "pendingApproval"
					if pendingAssetToClaim.TransactionID == "PENDING_AUTH" && u.PushNotificationToken != nil {

						if _, ok := notificationList[*u.PushNotificationToken]; ok {
							continue
						}
						u.SendPushMessage(fmt.Sprintf("%v submitted %v request on wallet %v!", signerUser.Username, "ACCEPT PENDING ASSET", wallet.Alias), fmt.Sprintf("Request: %v", pendingAssetToClaim.ReturnedDescription), "", dataPayload, gc)
						notificationList[*u.PushNotificationToken] = v.TargetUsername

					}

				}
			}
		} else {
			c.JSON(http.StatusAccepted, pendingAssetToClaim)
		}

	})

	router.DELETE("/v1/shared-access/users/actions/reject-asset", middleware.AuthenticationMiddlewareUsingTimestamp(), func(c *gin.Context) {
		var err error

		signerUser, err := usersDB.GetUserFromPrimarySigner(middleware.ExtractSigner(c), gc.DB, gc)

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
		{
			//check if pending shared access modify exists
			if userServices.CheckPendingSharedAccessApproval(middleware.ExtractPublicKey(c), gc.DB) {
				c.JSON(http.StatusForbidden, gin.H{"error": "error-pending-shared-access-op", "message": "There is a pending shared access operation on this wallet and must be completed first before attempting to send payment from this wallet."})
				return
			}
		}
		if wallet.WalletType != 0 {
			c.JSON(http.StatusForbidden, gin.H{"error": "error-wallet-type-forbidden", "message": "Operation not allowed on any special type of wallets. Only standard wallets are allowed."})
			return
		}

		hasInitiatorAccess := false
		isViewOnly := wallet.HasViewOnlyAccess(gc)

		if isViewOnly && wallet.Signer != signerUser.PrimarySigner {
			c.JSON(http.StatusForbidden, gin.H{"error": "error-unauthorized-access", "message": "You do not have an initiator permission on this wallet."})
			return
		}
		// check if user has initiator access to wallet.
		if wallet.SharedAccessEnabled == 1 && !isViewOnly {
			// check if user has initiator access to wallet.
			for _, p := range signerUser.WalletsSharedWithUser {
				if p.WalletPublicKey == middleware.ExtractPublicKey(c) && p.TargetUsername == signerUser.Username && p.Permission == "INITIATOR" {
					hasInitiatorAccess = true
				}
			}
			if !hasInitiatorAccess && !wallet.HasViewOnlyAccess(gc) {
				c.JSON(http.StatusForbidden, gin.H{"error": "error-unauthorized-access", "message": "You do not have an initiator permission on this wallet."})
				return
			}
		}

		conDB.PrintDBStats(fmt.Sprintf("REJECT /v1/shared-access/users/actions/reject-asset %v", middleware.ExtractPublicKey(c)), gc.DB)

		var pendingAssetToClaim userModels.PendingAssetToClaim
		// var err error

		data, _ := io.ReadAll(c.Request.Body)

		err = json.Unmarshal(data, &pendingAssetToClaim)

		var invalidJSON tErrors.ErrorInvalidJSON

		if err != nil {
			c.JSON(http.StatusBadRequest, invalidJSON.JSONError())
			return
		}

		_, complete, err := userServices.RejectPendingAsset(&signerUser, &wallet, &pendingAssetToClaim, gc)

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

		userCacheKey := fmt.Sprintf("[GET] /v1/users/%v", wallet.Alias)
		paymentPaymentHistoryCacheKey := fmt.Sprintf("[GET] /v1/users/payments/%v", middleware.ExtractPublicKey(c))

		gc.RedisCache.InvalidateCachedHttpResponse(ownerBalanceCacheKey, tempCacheKey, userCacheKey, paymentPaymentHistoryCacheKey, sNFT)
		log.Printf("[REJECT ASSET] Transaction Signature: [%v]\n", pendingAssetToClaim.TransactionSignature)
		if complete {
			c.JSON(http.StatusOK, pendingAssetToClaim)
			{
				//start push notificationMessage

				permissionList := wallet.Permissions
				for _, v := range permissionList {
					u, e := userModels.Username(v.TargetUsername).GetSimpleUser(gc.DB, gc)
					if e != nil {
						continue
					}

					dataPayload := make(map[string]string)
					dataPayload["none"] = ""
					if pendingAssetToClaim.TransactionID == "PENDING_AUTH" {
						u.SendPushMessage(fmt.Sprintf("%v submitted %v request on wallet %v!", signerUser.Username, "REJECT PENDING ASSET", wallet.Alias), fmt.Sprintf("Request: %v", pendingAssetToClaim.ReturnedDescription), "", dataPayload, gc)

					}

				}
			}
		} else {
			c.JSON(http.StatusAccepted, pendingAssetToClaim)
		}

	})

	router.GET("/v1/users/payment/generate/:targetUser", middleware.AuthenticationMiddlewareUsingTimestamp(), func(c *gin.Context) {
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

		cacheKey := fmt.Sprintf("[GET] /v1/users/payment/generate/%v", identifier)
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

		_, err = usersDB.GetUser(identifier, gc.DB, gc)

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

	router.POST("/v1/security-questions", middleware.AuthenticationMiddlewareUsingTimestamp(), func(c *gin.Context) {
		var err error

		var answers userModels.UserSecurityAnswer
		// var err error
		data, _ := io.ReadAll(c.Request.Body)

		err = json.Unmarshal(data, &answers)

		var invalidJSON tErrors.ErrorInvalidJSON

		if err != nil {
			c.JSON(http.StatusBadRequest, invalidJSON.JSONError())
			return
		}
		if len(answers.A1) == 0 || len(answers.A2) == 0 || len(answers.A3) == 0 || answers.Q1 == 0 || answers.Q2 == 0 || answers.Q3 == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Questions/Answers must be 3"})
			return
		}

		user, err := usersDB.GetUserFromPrimarySigner(middleware.ExtractSigner(c), gc.DB, gc)

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

		conDB.PrintDBStats(fmt.Sprintf("POST /v1/users/security-questions %v", user.Username), gc.DB)

		err = userServices.SaveUserSecurityQuestions(&user, answers, gc.DB)

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

		if user.PushNotificationToken != nil {
			dataPayload := make(map[string]string)
			dataPayload["none"] = ""
			pns.SendFirebaseMessage(*user.PushNotificationToken, "Secret Questions/Answers saved!", fmt.Sprintf("You have successfully saved secret questions in your account [%v].", user.Username), "", dataPayload, gc.PushNotificationClient, gc.PNSContext)
		}
		user.InvalidateUserCache(gc)
		//At this point, there was no error.

		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	router.GET("/v1/security-questions/:targetUser", middleware.AuthenticationMiddlewareUsingTimestamp(), func(c *gin.Context) {
		// var err error
		targetUser := strings.TrimSpace(strings.ToLower(c.Param("targetUser")))
		user, err := usersDB.GetUser(targetUser, gc.DB, gc)

		if err != nil {
			log.Println("[GET QUESTIONS] error for user:", targetUser, "error: ", err)

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

		conDB.PrintDBStats(fmt.Sprintf("security-questions %v", user.Username), gc.DB)

		securityQuestions := userServices.GetSecurityQuestions(&user, gc)

		userSecurityAnswers, _ := userServices.GetUserSecurityAnswers(&user, gc)
		userSecurityAnswers.A1 = ""
		userSecurityAnswers.A2 = ""
		userSecurityAnswers.A3 = ""

		c.JSON(http.StatusOK, gin.H{"securityQuestions": securityQuestions, "userSecurityAnswers": userSecurityAnswers})

	})

	router.POST("/v1/verify-answers/:targetUser", middleware.AuthenticationMiddlewareUsingTimestamp(), func(c *gin.Context) {
		var err error
		targetUser := strings.TrimSpace(strings.ToLower(c.Param("targetUser")))
		var answers userModels.UserSecurityAnswer
		// var err error

		data, _ := io.ReadAll(c.Request.Body)

		err = json.Unmarshal(data, &answers)

		var invalidJSON tErrors.ErrorInvalidJSON

		if err != nil {
			c.JSON(http.StatusBadRequest, invalidJSON.JSONError())
			return
		}

		if answers.Q1 == 0 || answers.Q2 == 0 || answers.Q3 == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Questions must be 3"})
			return
		}
		if len(answers.A1) == 0 || len(answers.A2) == 0 || len(answers.A3) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Answers must be 3"})
			return
		}

		user, err := usersDB.GetUser(targetUser, gc.DB, gc)

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

		conDB.PrintDBStats(fmt.Sprintf("POST /v1/users/verify-answers/%v", user.Username), gc.DB)

		if !userServices.ValidateSecurityAnswers(&user, answers, gc) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "The answers provided are invalid."})
			return
		}

		//At this point, there was no error.

		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	router.POST("/v1/account/recovery/request-email-otp/:targetUser", middleware.AuthenticationMiddlewareUsingTimestamp(), func(c *gin.Context) {
		var err error

		targetUser := strings.TrimSpace(c.Param("targetUser"))

		user, err := usersDB.GetUser(targetUser, gc.DB, gc)

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

		conDB.PrintDBStats(fmt.Sprintf("POST /v1/account/recovery/request-email-otp/%v", user.Username), gc.DB)

		if err = userServices.SendAccountRecoveryEmailOTP(&user, gc.DB); err != nil {
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

		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	router.POST("/v1/account/recovery/verify-email-otp/:targetUser/:otp", middleware.AuthenticationMiddlewareUsingTimestamp(), func(c *gin.Context) {
		var err error
		otp := strings.TrimSpace(c.Param("otp"))
		targetUser := strings.TrimSpace(c.Param("targetUser"))

		user, err := usersDB.GetUser(targetUser, gc.DB, gc)

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

		conDB.PrintDBStats(fmt.Sprintf("POST /v1/verify-email-otp/%v/%v", user.Username, otp), gc.DB)

		if err = userServices.CheckAccountRecoveryEmailOTP(&user, otp, gc.DB); err != nil {
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

		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	router.POST("/v1/verify-email-otp/:targetUser/:otp", middleware.AuthenticationMiddlewareUsingTimestamp(), func(c *gin.Context) {
		var err error
		otp := strings.TrimSpace(c.Param("otp"))
		targetUser := strings.TrimSpace(c.Param("targetUser"))

		user, err := usersDB.GetUser(targetUser, gc.DB, gc)

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

		conDB.PrintDBStats(fmt.Sprintf("POST /v1/verify-email-otp/%v/%v", user.Username, otp), gc.DB)

		if err = userServices.CheckAccountRecoveryEmailOTP(&user, otp, gc.DB); err != nil {
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

		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	router.POST("/v1/users/account/recovery", middleware.AuthenticationMiddlewareUsingTimestamp(), func(c *gin.Context) {
		var err error

		var payload userModels.UserAccountRecoveryPayload
		// var err error

		data, _ := io.ReadAll(c.Request.Body)

		err = json.Unmarshal(data, &payload)

		var invalidJSON tErrors.ErrorInvalidJSON

		if err != nil {
			c.JSON(http.StatusBadRequest, invalidJSON.JSONError())
			return
		}

		user, err := usersDB.GetUserFromPrimarySigner(middleware.ExtractSigner(c), gc.DB, gc)

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

		conDB.PrintDBStats(fmt.Sprintf("POST /v1/users/account/recovery  %v", user.Username), gc.DB)
		err = userServices.EnableAccountRecovery(&user, &payload, gc)
		if err != nil {
			log.Printf("Failed to enable account recovery: %v\n", err)
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
		if len(payload.TransactionID) > 0 {
			if user.PushNotificationToken != nil {
				dataPayload := make(map[string]string)
				dataPayload["none"] = ""
				pns.SendFirebaseMessage(*user.PushNotificationToken, "Account Recovery Enabled!", fmt.Sprintf("Congratulations! You have successfully enabled account recovery service on your account [%v]. Your account will be recovered by Trovotech should you lose your secret key. The service will expire on %v. We will notify you when it is time to renew the service to keep your account recovery active.", user.Username, user.AccountRecoveryExpiresOn.Format("01-02-2006 15:04:05")), "", dataPayload, gc.PushNotificationClient, gc.PNSContext)
				c.JSON(http.StatusOK, payload)
			}
		} else {
			c.JSON(http.StatusAccepted, payload)
		}

	})

	router.DELETE("/v1/users/account/recovery", middleware.AuthenticationMiddlewareUsingTimestamp(), func(c *gin.Context) {
		var err error

		var payload userModels.UserAccountRecoveryPayload
		// var err error

		data, _ := io.ReadAll(c.Request.Body)

		err = json.Unmarshal(data, &payload)

		var invalidJSON tErrors.ErrorInvalidJSON

		if err != nil {
			c.JSON(http.StatusBadRequest, invalidJSON.JSONError())
			return
		}

		user, err := usersDB.GetUserFromPrimarySigner(middleware.ExtractSigner(c), gc.DB, gc)

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

		conDB.PrintDBStats(fmt.Sprintf("DELETE /v1/users/account/recovery  %v", user.Username), gc.DB)
		err = userServices.DisableAccountRecovery(&user, &payload, gc)
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
		if len(payload.TransactionID) > 0 {
			if user.PushNotificationToken != nil {
				dataPayload := make(map[string]string)
				dataPayload["none"] = ""
				pns.SendFirebaseMessage(*user.PushNotificationToken, "Account Recovery Disabled!", fmt.Sprintf("Congratulations! You have successfully disabled account recovery feature on your account [%v].", user.Username), "", dataPayload, gc.PushNotificationClient, gc.PNSContext)
			}
			c.JSON(http.StatusOK, payload)
		} else {
			c.JSON(http.StatusAccepted, payload)
		}

	})

	router.POST("/v1/users/account/recover", middleware.AuthenticationMiddlewareUsingTimestamp(), func(c *gin.Context) {
		var err error

		var payload userModels.AccountRecoveryRequest
		// var err error

		data, _ := io.ReadAll(c.Request.Body)

		err = json.Unmarshal(data, &payload)

		var invalidJSON tErrors.ErrorInvalidJSON

		if err != nil {
			c.JSON(http.StatusBadRequest, invalidJSON.JSONError())
			return
		}

		user, err := usersDB.GetUser(payload.Username, gc.DB, gc)

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

		conDB.PrintDBStats(fmt.Sprintf("POST /v1/users/account/recover  %v", user.Username), gc.DB)
		_, err = userServices.DoAccountRecovery(&user, &payload, gc)
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
		if payload.Commit == 1 {
			c.JSON(http.StatusOK, payload)

			if user.PushNotificationToken != nil {
				dataPayload := make(map[string]string)
				dataPayload["none"] = ""
				pns.SendFirebaseMessage(*user.PushNotificationToken, "Account Recovery Successful!", fmt.Sprintf("Congratulations! You have successfully recovered your account [%v]. Please import the new secret key using the same username specified.", user.Username), "", dataPayload, gc.PushNotificationClient, gc.PNSContext)
			}
		} else {
			c.JSON(http.StatusAccepted, payload)
		}

	})

	router.POST("/v1/users/inactive-account/recover", middleware.AuthenticationMiddlewareUsingTimestamp(), func(c *gin.Context) {
		var err error

		var payload userModels.InactiveAccountRecoveryRequest
		// var err error

		data, _ := io.ReadAll(c.Request.Body)

		err = json.Unmarshal(data, &payload)

		var invalidJSON tErrors.ErrorInvalidJSON

		if err != nil {
			c.JSON(http.StatusBadRequest, invalidJSON.JSONError())
			return
		}

		subjectUser, err := usersDB.GetUser(payload.Username, gc.DB, gc)

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

		conDB.PrintDBStats(fmt.Sprintf("POST /v1/users/inactive-account/recover  %v", subjectUser.Username), gc.DB)
		userInfo, err := userServices.DoInactiveAccountRecover(&subjectUser, &payload, gc)
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

		c.JSON(http.StatusOK, userInfo)

		if subjectUser.PushNotificationToken != nil {
			dataPayload := make(map[string]string)
			dataPayload["none"] = ""
			pns.SendFirebaseMessage(*subjectUser.PushNotificationToken, "Account updated successfully!", fmt.Sprintf("Congratulations! You have successfully updated your account [%v] with the new secret key and security questions. Please import the new secret key using the same username specified.", subjectUser.Username), "", dataPayload, gc.PushNotificationClient, gc.PNSContext)
		}

	})

	router.POST("/v1/shared-access/users/account", middleware.AuthenticationMiddlewareUsingTimestamp(), func(c *gin.Context) {
		var err error

		signerUser, err := usersDB.GetUserFromPrimarySigner(middleware.ExtractSigner(c), gc.DB, gc)

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
		walletOwner, err := usersDB.GetUser(middleware.ExtractPublicKey(c), gc.DB, gc)

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
		if signerUser.PrimarySigner != wallet.Signer {
			c.JSON(http.StatusForbidden, gin.H{"error": "error-unauthorized-access", "message": "You do not have permission on this wallet."})
			return
		}
		conDB.PrintDBStats(fmt.Sprintf("POST /v1/shared-access/users/account %v", middleware.ExtractPublicKey(c)), gc.DB)

		var sharedAccessInfo userModels.UserWalletSharedAccessInfo
		// var err error

		data, _ := io.ReadAll(c.Request.Body)
		log.Println(string(data))
		err = json.Unmarshal(data, &sharedAccessInfo)

		var invalidJSON tErrors.ErrorInvalidJSON

		if err != nil {
			c.JSON(http.StatusBadRequest, invalidJSON.JSONError())
			return
		}
		sharedAccessInfo.WalletPublicKey = middleware.ExtractPublicKey(c)
		log.Printf("[DEBUG] sharedAccess %+v\n", sharedAccessInfo)
		_, err = userServices.CreateSharedWalletAccess(&signerUser, &walletOwner, &wallet, &sharedAccessInfo, gc)

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
		log.Printf("[CREATE SHARED ACCESS] Transaction Signature: [%v]\n", sharedAccessInfo.TransactionSignature)
		if len(sharedAccessInfo.TransactionID) > 0 {
			for _, v := range sharedAccessInfo.Permissions {

				if v.PushNotificationToken != nil {
					dataPayload := make(map[string]string)
					dataPayload["none"] = ""
					pns.SendFirebaseMessage(*v.PushNotificationToken, fmt.Sprintf("%v permission granted on wallet %v!", v.Permission, v.WalletAlias), fmt.Sprintf("You have been granted %v permission on the wallet [%v]. Please navigate to the section for third-party wallet access whenever you wish to perform tasks relating to this wallet.", v.Permission, v.WalletAlias), "", dataPayload, gc.PushNotificationClient, gc.PNSContext)
				}
			}
			c.JSON(http.StatusOK, sharedAccessInfo)

		} else {
			c.JSON(http.StatusAccepted, sharedAccessInfo)
		}

	})

	//modify shared access
	router.PUT("/v1/shared-access/users/account", middleware.AuthenticationMiddlewareUsingTimestamp(), func(c *gin.Context) {
		var err error

		signerUser, err := usersDB.GetUserFromPrimarySigner(middleware.ExtractSigner(c), gc.DB, gc)

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
		walletOwner, err := usersDB.GetUser(middleware.ExtractPublicKey(c), gc.DB, gc)

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
		walletOwner.InvalidateUserCache(gc)
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

		hasInitiatorAccess := false
		isViewOnly := wallet.HasViewOnlyAccess(gc)
		// check if user has initiator access to wallet.

		if isViewOnly && wallet.Signer != signerUser.PrimarySigner {
			c.JSON(http.StatusForbidden, gin.H{"error": "error-unauthorized-access", "message": "You do not have an initiator permission on this wallet."})
			return
		}
		// check if user has initiator access to wallet.
		if wallet.SharedAccessEnabled == 1 && !isViewOnly {
			for _, p := range signerUser.WalletsSharedWithUser {
				if p.WalletPublicKey == middleware.ExtractPublicKey(c) && p.TargetUsername == signerUser.Username && p.Permission == "INITIATOR" {
					hasInitiatorAccess = true
				}
			}
			if !hasInitiatorAccess && !userModels.UserWalletID(middleware.ExtractPublicKey(c)).PublicKeyHasViewOnlyAccess(gc) {
				c.JSON(http.StatusForbidden, gin.H{"error": "error-unauthorized-access", "message": "You do not have an initiator permission on this wallet."})
				return
			}
		}

		conDB.PrintDBStats(fmt.Sprintf("PUT /v1/shared-access/users/account %v", middleware.ExtractPublicKey(c)), gc.DB)

		var sharedAccessInfo userModels.ModifySharedAccessInfo
		// var err error

		data, _ := io.ReadAll(c.Request.Body)
		log.Println(string(data))
		err = json.Unmarshal(data, &sharedAccessInfo)

		var invalidJSON tErrors.ErrorInvalidJSON

		if err != nil {
			c.JSON(http.StatusBadRequest, invalidJSON.JSONError())
			return
		}
		sharedAccessInfo.WalletPublicKey = middleware.ExtractPublicKey(c)
		log.Printf("[DEBUG] modify %+v\n", sharedAccessInfo)
		_, _, _, err = userServices.ModifySharedWalletAccess(&signerUser, &walletOwner, &wallet, &sharedAccessInfo, gc)

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

		var ownerBalanceCacheKey, tempCacheKey, sNFT string

		ownerBalanceCacheKey = fmt.Sprintf("GetBalance_%s", middleware.ExtractPublicKey(c))
		sNFT = fmt.Sprintf("GetNFTs_%s", middleware.ExtractPublicKey(c))

		tempCacheKey = fmt.Sprintf("GetBalance_%s", *wallet.TempPublicKey)

		userCacheKey := fmt.Sprintf("[GET] /v1/users/%v", walletOwner.Username)
		paymentPaymentHistoryCacheKey := fmt.Sprintf("[GET] /v1/users/payments/%v", middleware.ExtractPublicKey(c))

		gc.RedisCache.InvalidateCachedHttpResponse(ownerBalanceCacheKey, tempCacheKey, userCacheKey, paymentPaymentHistoryCacheKey, sNFT)
		log.Printf("[MODIFY SHARED ACCESS] Transaction Signature: [%v]\n", sharedAccessInfo.TransactionSignature)
		if len(sharedAccessInfo.TransactionID) > 0 {
			if sharedAccessInfo.TransactionID == "PENDING_AUTH" {
				wallet, _, _ := usersDB.GetWallet(middleware.ExtractPublicKey(c), gc.DB)
				//saved to pending auth table for disabling shared access
				notificationList := make(map[string]string)
				for _, v := range wallet.Permissions {
					if v.Permission == "APPROVER" {
						u, e := userModels.Username(v.TargetUsername).GetSimpleUser(gc.DB, gc)
						if e == nil && u.PushNotificationToken != nil {

							if _, ok := notificationList[*u.PushNotificationToken]; ok {
								continue
							}

							log.Println("notifying approver:", v.TargetUsername)
							dataPayload := make(map[string]string)
							dataPayload["route"] = "pendingApproval"
							u.SendPushMessage(fmt.Sprintf("Pending Approval: Modify shared access on wallet %v!", wallet.Alias), fmt.Sprintf("You have a pending approval to modify shared access on the wallet %v. Please tap to choose the appropriate action.", wallet.Alias), "", dataPayload, gc)
							notificationList[*u.PushNotificationToken] = v.TargetUsername
						}
					}
				}
				c.JSON(http.StatusOK, sharedAccessInfo)
				return
			}
			log.Println("notifying walletOwner:", walletOwner.Username)
			dataPayload := make(map[string]string)
			dataPayload["route"] = "pendingApproval"
			walletOwner.SendPushMessage(fmt.Sprintf("Pending Approval: Modify shared access on wallet %v!", wallet.Alias), fmt.Sprintf("You have a pending approval to modify shared access on the wallet %v. Please tap to choose the appropriate action.", wallet.Alias), "", dataPayload, gc)

			c.JSON(http.StatusOK, sharedAccessInfo)
			return

		} else {
			c.JSON(http.StatusAccepted, sharedAccessInfo)
		}

	})

	router.DELETE("/v1/shared-access/users/account", middleware.AuthenticationMiddlewareUsingTimestamp(), func(c *gin.Context) {
		var err error
		var sharedAccessInfo userModels.DisableSharedAccessInfo
		// var err error

		data, _ := io.ReadAll(c.Request.Body)
		log.Println(string(data))
		err = json.Unmarshal(data, &sharedAccessInfo)

		var invalidJSON tErrors.ErrorInvalidJSON

		if err != nil {
			c.JSON(http.StatusBadRequest, invalidJSON.JSONError())
			return
		}
		signerUser, err := usersDB.GetUserFromPrimarySigner(middleware.ExtractSigner(c), gc.DB, gc)

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
		walletOwner, err := usersDB.GetUser(middleware.ExtractPublicKey(c), gc.DB, gc)

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

		//check if owner is the initiator
		var isInitiator bool
		isViewOnly := wallet.HasViewOnlyAccess(gc)
		pl := wallet.Permissions
		if len(pl) == 0 {
			//reject request
			c.JSON(http.StatusBadRequest, gin.H{"error": "error-permissions-not-found", "message": "Could not determine permissions on this wallet at this time. Please try again later."})
			return
		}
		if wallet.SharedAccessEnabled == 0 {
			//reject request
			c.JSON(http.StatusBadRequest, gin.H{"error": "error-shared-access-not-enabled", "message": "Shared access is not currently enabled on this wallet."})
			return
		}
		for _, p := range pl {

			if p.Permission == "INITIATOR" || p.Permission == "APPROVER" {
				isViewOnly = false
			}
			if p.Permission == "INITIATOR" && p.TargetUsername == signerUser.Username {
				isInitiator = true

			}

		}

		if signerUser.Username == walletOwner.Username && !isViewOnly && !isInitiator {
			//reject request
			c.JSON(http.StatusBadRequest, gin.H{"error": "error-not-an-initiator", "message": "You do not have initiator permission on this wallet. Only an initiator can submit this transaction."})
			return
		}

		if signerUser.Username != walletOwner.Username && !isInitiator {
			//reject request
			c.JSON(http.StatusBadRequest, gin.H{"error": "error-not-an-initiator", "message": "You do not have initiator permission on this wallet. Only an initiator can submit this transaction."})
			return
		}

		conDB.PrintDBStats(fmt.Sprintf("DELETE /v1/shared-access/users/account %v", middleware.ExtractPublicKey(c)), gc.DB)

		sharedAccessInfo.WalletPublicKey = middleware.ExtractPublicKey(c)

		log.Printf("[DEBUG] sharedAccess %+v\n", sharedAccessInfo)
		err = userServices.RemoveSharedWalletAccess(&signerUser, &wallet, &sharedAccessInfo, gc)

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
		log.Printf("[REMOVED SHARED ACCESS] Transaction Signature: [%v]\n", sharedAccessInfo.TransactionSignature)
		if len(sharedAccessInfo.TransactionID) > 0 {
			if sharedAccessInfo.TransactionID == "PENDING_AUTH" {
				//saved to pending auth table for disabling shared access
				for _, v := range sharedAccessInfo.Permissions {

					if v.PushNotificationToken != nil && v.Permission == "APPROVER" {
						log.Println("notifying approver:", v.TargetUsername)
						dataPayload := make(map[string]string)
						dataPayload["link"] = "authPending"
						pns.SendFirebaseMessage(*v.PushNotificationToken, fmt.Sprintf("Pending Approval: Disable shared access on wallet %v!", v.WalletAlias), fmt.Sprintf("You have a pending approval to disable shared access on the wallet %v. Please tap to choose the appropriate action.", v.WalletAlias), "", dataPayload, gc.PushNotificationClient, gc.PNSContext)
					}
				}
				c.JSON(http.StatusOK, sharedAccessInfo)
				return
			} else {
				//transaction was completed successfully
				notificationList := make(map[string]string)
				for _, v := range sharedAccessInfo.Permissions {

					if v.PushNotificationToken != nil {

						if _, ok := notificationList[*v.PushNotificationToken]; ok {
							continue
						}

						dataPayload := make(map[string]string)
						dataPayload["none"] = ""
						pns.SendFirebaseMessage(*v.PushNotificationToken, fmt.Sprintf("Your %v permission on wallet %v has been removed!", v.Permission, v.WalletAlias), fmt.Sprintf("Your %v permission on the wallet [%v] has been removed as shared access has been disabled on the wallet. the wallet is no longer available on your list of shared access wallets.", v.Permission, v.WalletAlias), "", dataPayload, gc.PushNotificationClient, gc.PNSContext)
						notificationList[*v.PushNotificationToken] = v.TargetUsername
					}
				}
				c.JSON(http.StatusOK, sharedAccessInfo)
				return
			}
		} else {
			c.JSON(http.StatusAccepted, sharedAccessInfo)
		}

	})

	//get transaction list
	router.GET("/v1/shared-access/approvals", middleware.AuthenticationMiddlewareUsingTimestamp(), func(c *gin.Context) {

		conDB.PrintDBStats(fmt.Sprintf("[GET] /v1/shared-access/approvals %v", middleware.ExtractSigner(c)), gc.DB)

		signerUser, err := usersDB.GetUserFromPrimarySigner(middleware.ExtractSigner(c), gc.DB, gc)

		if err != nil {
			log.Println("[GET GetApprovalRequest] error for signer:", middleware.ExtractSigner(c), "error: ", err)

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
			// gc.RedisCache.CacheHttpResponseWithParameters(cacheKey, cacheKeyParameters, statusCode, response, cacheDurationInSeconds)
			return
		}
		permittedPublicKeys := make([]string, 0)

		for _, k := range signerUser.WalletsSharedWithUser {
			if k.Permission == "VIEW-ONLY" {
				continue
			}
			// view only is not permitted to see transactions
			permittedPublicKeys = append(permittedPublicKeys, k.WalletPublicKey)
		}

		//Get Payment history
		historyRecords := userServices.GetApprovalList(&signerUser, permittedPublicKeys, gc, c)

		c.JSON(http.StatusOK, historyRecords)

	})

	//get specific transaction
	router.GET("/v1/shared-access/approval/:ID", middleware.AuthenticationMiddlewareUsingTimestamp(), func(c *gin.Context) {
		// var err error

		approvalID := c.Param("ID")
		conDB.PrintDBStats(fmt.Sprintf("[GET] /v1/shared-access/approval/%v", approvalID), gc.DB)

		signerUser, err := usersDB.GetUserFromPrimarySigner(middleware.ExtractSigner(c), gc.DB, gc)

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
				response = gin.H{"error": err.Error(), "message": err.Error()}
			}

			c.JSON(statusCode, response)
			// gc.RedisCache.CacheHttpResponseWithParameters(cacheKey, cacheKeyParameters, statusCode, response, cacheDurationInSeconds)
			return
		}
		permittedPublicKeys := make([]string, 0)

		for _, k := range signerUser.WalletsSharedWithUser {
			if k.Permission == "VIEW-ONLY" {
				continue
			}
			// view only is not permitted to see transactions
			permittedPublicKeys = append(permittedPublicKeys, k.WalletPublicKey)
		}
		if len(permittedPublicKeys) == 0 {
			c.JSON(http.StatusForbidden, gin.H{"error": "error-access-forbidden", "message": "You do not have needed permissions to access section."})
			return
		}

		//Get Payment history
		approvalRequestJSON, err := userServices.GetApprovalRequestJSON(approvalID, gc)
		if err != nil {
			log.Println("[GET GetApprovalRequest] error for signer:", signerUser.Username, "error: ", err)

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
			// gc.RedisCache.CacheHttpResponseWithParameters(cacheKey, cacheKeyParameters, statusCode, response, cacheDurationInSeconds)
			return
		}
		c.JSON(http.StatusOK, approvalRequestJSON)
		// gc.RedisCache.CacheHttpResponse(cacheKey, http.StatusOK, historyRecords, cacheDurationInSeconds)
		// gc.RedisCache.CacheHttpResponseWithParameters(cacheKey, cacheKeyParameters, http.StatusOK, historyRecords, cacheDurationInSeconds)

	})

	//submit specific approval signature
	router.POST("/v1/shared-access/approval/:ID", middleware.AuthenticationMiddlewareUsingTimestamp(), func(c *gin.Context) {
		// var err error

		var payload userModels.ApprovalPayload
		// var err error

		data, _ := io.ReadAll(c.Request.Body)

		err := json.Unmarshal(data, &payload)

		var invalidJSON tErrors.ErrorInvalidJSON

		if err != nil {
			c.JSON(http.StatusBadRequest, invalidJSON.JSONError())
			return
		}
		approvalID := c.Param("ID")
		payload.DeviceID = c.GetHeader("X-TW-DEVICE-ID")
		conDB.PrintDBStats(fmt.Sprintf("[POST] /v1/shared-access/approval/%v", approvalID), gc.DB)

		signerUser, err := usersDB.GetUserFromPrimarySigner(middleware.ExtractSigner(c), gc.DB, gc)

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
				response = gin.H{"error": err.Error(), "message": err.Error()}
			}

			c.JSON(statusCode, response)
			// gc.RedisCache.CacheHttpResponseWithParameters(cacheKey, cacheKeyParameters, statusCode, response, cacheDurationInSeconds)
			return
		}
		permittedPublicKeys := make(map[string]string, 0)

		for _, k := range signerUser.WalletsSharedWithUser {
			if k.Permission != "APPROVER" {
				continue
			}
			// view only is not permitted to see transactions
			permittedPublicKeys[k.WalletPublicKey] = k.WalletPublicKey
			// permittedPublicKeys = append(permittedPublicKeys, k.WalletPublicKey)
		}
		if len(permittedPublicKeys) == 0 {
			c.JSON(http.StatusForbidden, gin.H{"error": "error-access-forbidden", "message": "You do not have needed permissions to access section."})
			return
		}

		approvalRequest, err := userServices.GetApprovalRequest(approvalID, gc)
		if err != nil {
			log.Println("[POST ApproveRequest] error for signer:", signerUser.Username, "error: ", err)

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
			// gc.RedisCache.CacheHttpResponseWithParameters(cacheKey, cacheKeyParameters, statusCode, response, cacheDurationInSeconds)
			return
		}

		if _, ok := permittedPublicKeys[approvalRequest.WalletPublicKey]; !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "error-access-forbidden", "message": "You do not have needed permission."})
			return
		}
		//check if user already approved before

		{
			for _, a := range approvalRequest.PendingTransactionSignatures {
				if a.Approver == signerUser.Username {
					c.JSON(http.StatusForbidden, gin.H{"error": "error-duplicate-approval", "message": "An approval from you already exists. You can only submit one approval."})
					return
				}
			}
		}

		err = userServices.ApproveTransaction(&signerUser, &approvalRequest, &payload, callBackRetryChan, gc)
		if err != nil {
			log.Println("[POST ApproveRequest] error for signer:", signerUser.Username, "error: ", err)

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
			// gc.RedisCache.CacheHttpResponseWithParameters(cacheKey, cacheKeyParameters, statusCode, response, cacheDurationInSeconds)
			return
		}
		if len(payload.TransactionSignature) == 0 {
			c.JSON(http.StatusAccepted, payload)
			return
		}

		c.JSON(http.StatusOK, payload)
		// gc.RedisCache.CacheHttpResponse(cacheKey, http.StatusOK, historyRecords, cacheDurationInSeconds)
		// gc.RedisCache.CacheHttpResponseWithParameters(cacheKey, cacheKeyParameters, http.StatusOK, historyRecords, cacheDurationInSeconds)
		{
			notificationList := make(map[string]string)
			//start push notificationMessage
			wallet, e := userModels.UserWalletID(approvalRequest.WalletPublicKey).GetWallet(gc.DB, gc)
			if e != nil {
				return
			}
			permissionList := wallet.Permissions
			for _, v := range permissionList {

				u, e := userModels.Username(v.TargetUsername).GetSimpleUser(gc.DB, gc)
				if e != nil {
					continue
				}

				if u.PushNotificationToken != nil && v.Permission != "VIEW-ONLY" {
					if _, ok := notificationList[*u.PushNotificationToken]; ok {
						continue
					}
					dataPayload := make(map[string]string)
					dataPayload["route"] = "pendingApproval"
					if approvalRequest.TransactionStatus != "COMPLETED" {
						pns.SendFirebaseMessage(*u.PushNotificationToken, fmt.Sprintf("%v Submitted an approval on wallet %v!", signerUser.Username, wallet.Alias), fmt.Sprintf("%v submitted an approval for request:\n%v\nApproval stage is now %v/%v", signerUser.Username, approvalRequest.Description, approvalRequest.ApprovalsGotten, approvalRequest.ApprovalsNeeded), "", dataPayload, gc.PushNotificationClient, gc.PNSContext)
						notificationList[*u.PushNotificationToken] = v.TargetUsername
					}
				}
			}
		}
	})

	//reject specific approval signature
	router.DELETE("/v1/shared-access/approval/:ID", middleware.AuthenticationMiddlewareUsingTimestamp(), func(c *gin.Context) {
		// var err error

		var payload userModels.RejectPayload
		// var err error

		data, _ := io.ReadAll(c.Request.Body)

		err := json.Unmarshal(data, &payload)

		var invalidJSON tErrors.ErrorInvalidJSON

		if err != nil {
			c.JSON(http.StatusBadRequest, invalidJSON.JSONError())
			return
		}
		approvalID := c.Param("ID")
		conDB.PrintDBStats(fmt.Sprintf("[DELETE] /v1/shared-access/approval/%v", approvalID), gc.DB)

		signerUser, err := usersDB.GetUserFromPrimarySigner(middleware.ExtractSigner(c), gc.DB, gc)

		if err != nil {
			log.Println("[GET RejectRequest] error for signer:", middleware.ExtractSigner(c), "error: ", err)

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
			// gc.RedisCache.CacheHttpResponseWithParameters(cacheKey, cacheKeyParameters, statusCode, response, cacheDurationInSeconds)
			return
		}
		permittedPublicKeys := make(map[string]string, 0)

		for _, k := range signerUser.WalletsSharedWithUser {
			if k.Permission != "APPROVER" {
				continue
			}
			// view only is not permitted to see transactions
			permittedPublicKeys[k.WalletPublicKey] = k.WalletPublicKey
			// permittedPublicKeys = append(permittedPublicKeys, k.WalletPublicKey)
		}
		if len(permittedPublicKeys) == 0 {
			c.JSON(http.StatusForbidden, gin.H{"error": "error-access-forbidden", "message": "You do not have needed permissions to access section."})
			return
		}

		approvalRequest, err := userServices.GetApprovalRequest(approvalID, gc)
		if err != nil {
			log.Println("[POST RejectRequest] error for signer:", signerUser.Username, "error: ", err)

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
			// gc.RedisCache.CacheHttpResponseWithParameters(cacheKey, cacheKeyParameters, statusCode, response, cacheDurationInSeconds)
			return
		}

		//check if user already approved before

		{
			for _, a := range approvalRequest.PendingTransactionSignatures {
				if a.Approver == signerUser.Username {
					c.JSON(http.StatusForbidden, gin.H{"error": "error-duplicate-approval", "message": "An approval from you already exists. You can only submit one approval."})
					return
				}
			}
		}

		if _, ok := permittedPublicKeys[approvalRequest.WalletPublicKey]; !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "error-access-forbidden", "message": "You do not have needed permission."})
			return
		}

		err = userServices.RejectTransaction(&signerUser, &approvalRequest, &payload, gc)
		if err != nil {
			log.Println("[POST RejectRequest] error for signer:", signerUser.Username, "error: ", err)

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
			// gc.RedisCache.CacheHttpResponseWithParameters(cacheKey, cacheKeyParameters, statusCode, response, cacheDurationInSeconds)
			return
		}

		c.JSON(http.StatusOK, payload)
		// gc.RedisCache.CacheHttpResponse(cacheKey, http.StatusOK, historyRecords, cacheDurationInSeconds)
		// gc.RedisCache.CacheHttpResponseWithParameters(cacheKey, cacheKeyParameters, http.StatusOK, historyRecords, cacheDurationInSeconds)
		{
			notificationList := make(map[string]string)

			//start push notificationMessage
			wallet, e := userModels.UserWalletID(approvalRequest.WalletPublicKey).GetWallet(gc.DB, gc)
			if e != nil {
				return
			}
			permissionList := wallet.Permissions
			for _, v := range permissionList {
				u, e := userModels.Username(v.TargetUsername).GetSimpleUser(gc.DB, gc)
				if e != nil {
					continue
				}

				dataPayload := make(map[string]string)
				dataPayload["route"] = "pendingApproval"
				if approvalRequest.TransactionStatus == "REJECTED" && u.PushNotificationToken != nil {

					// notificationList := make(map[string]string)

					if _, ok := notificationList[*u.PushNotificationToken]; ok {
						continue
					}
					// notificationList[*u.PushNotificationToken] = v.TargetUsername

					u.SendPushMessage(fmt.Sprintf("%v rejected %v request on wallet %v!", signerUser.Username, approvalRequest.TransactionType, wallet.Alias), fmt.Sprintf("Reason: %v\nRequest:%v", *approvalRequest.ReasonForRejection, approvalRequest.Description), "", dataPayload, gc)
					notificationList[*u.PushNotificationToken] = v.TargetUsername
				}

			}
		}
	})

	//get specific  wallet balance
	router.GET("/v1/shared-access/wallet-balances", middleware.AuthenticationMiddlewareUsingTimestamp(), func(c *gin.Context) {
		// var err error

		conDB.PrintDBStats(fmt.Sprintf("[GET] /v1/shared-access/wallet-balances %v", middleware.ExtractPublicKey(c)), gc.DB)

		signerUser, err := usersDB.GetUserFromPrimarySigner(middleware.ExtractSigner(c), gc.DB, gc)

		if err != nil {
			log.Println("[GET Wallet Balances] error for signer:", middleware.ExtractSigner(c), "error: ", err)

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
			// gc.RedisCache.CacheHttpResponseWithParameters(cacheKey, cacheKeyParameters, statusCode, response, cacheDurationInSeconds)
			return
		}
		permitted := false

		for _, k := range signerUser.WalletsSharedWithUser {
			if k.TargetUsername == signerUser.Username && k.WalletPublicKey == middleware.ExtractPublicKey(c) {
				permitted = true
			}
		}
		if !permitted {
			c.JSON(http.StatusForbidden, gin.H{"error": "error-access-forbidden", "message": "You do not have needed permissions to access this wallet."})
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

		assetBalances, err := wallet.GetWalletAssetBalances(gc)
		if err != nil {
			log.Println("[GET Wallet Balances] error for signer:", signerUser.Username, "error: ", err)

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
		nfts, err := wallet.GetNFTs(false, gc)
		if err != nil {
			log.Println("[GET Wallet Balances] error for signer:", signerUser.Username, "error: ", err)

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

		c.JSON(http.StatusOK, gin.H{"assetBalances": assetBalances, "nfts": nfts})

	})

	// market making
	{
		router.POST("/v1/users/trades", middleware.AuthenticationMiddlewareUsingTimestamp(), func(c *gin.Context) {
			var err error

			var makeOfferRequest userModels.MarketOfferRequest
			// var err error

			data, _ := io.ReadAll(c.Request.Body)

			err = json.Unmarshal(data, &makeOfferRequest)

			var invalidJSON tErrors.ErrorInvalidJSON

			if err != nil {
				c.JSON(http.StatusBadRequest, invalidJSON.JSONError())
				return
			}

			signerUser, err := usersDB.GetUserFromPrimarySigner(middleware.ExtractSigner(c), gc.DB, gc)

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
			walletOwner, err := usersDB.GetUser(middleware.ExtractPublicKey(c), gc.DB, gc)

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
			if wallet.WalletType != 0 {
				c.JSON(http.StatusForbidden, gin.H{"error": "error-wallet-type-forbidden", "message": "Operation not allowed on any special type of wallets. Only standard wallets are allowed."})
				return
			}
			if signerUser.PrimarySigner != wallet.Signer {
				c.JSON(http.StatusForbidden, gin.H{"error": "error-unauthorized-access", "message": "You do not have permission on this wallet."})
				return
			}
			conDB.PrintDBStats(fmt.Sprintf("POST /v1/users/trades %v", wallet.Alias), gc.DB)

			err = userServices.MakeOffer(&signerUser, &walletOwner, &wallet, &makeOfferRequest, gc)

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

			if walletOwner.PushNotificationToken != nil && len(makeOfferRequest.TransactionID) > 0 && makeOfferRequest.TransactionID != "PENDING_AUTH" {
				dataPayload := make(map[string]string)
				dataPayload["route"] = ""
				pns.SendFirebaseMessage(*walletOwner.PushNotificationToken, fmt.Sprintf("%v %v %v offer accepted on %v!", makeOfferRequest.Quantity, makeOfferRequest.AssetCode, makeOfferRequest.OfferType, wallet.Alias), fmt.Sprintf("You have successfully submitted a market offer to %v %v %v @ %v %v on the wallet with alias [%v].", makeOfferRequest.OfferType, makeOfferRequest.Quantity, makeOfferRequest.AssetCode, makeOfferRequest.PricePerUnit, makeOfferRequest.CurrencyCode, wallet.Alias), "", dataPayload, gc.PushNotificationClient, gc.PNSContext)
			}
			walletOwner.InvalidateUserCache(gc)
			//At this point, there was no error.

			c.JSON(http.StatusOK, makeOfferRequest)
		})

	}
	//CRYPTO
	{
		router.GET("/v1/crypto/withrawal-history/:currency/:targetPublicKeyForHistory", middleware.AuthenticationMiddlewareUsingTimestamp(), func(c *gin.Context) {
			// var err error
			currency := strings.ToUpper(c.Param("currency"))

			targetPublicKeyForHistory := strings.TrimSpace(strings.ToUpper(c.Param("targetPublicKeyForHistory")))

			_, err := keypair.ParseAddress(targetPublicKeyForHistory)
			if err != nil {

				statusCode := http.StatusBadRequest
				response := gin.H{"error": "error invalid address", "message": "Only valid addresses are allowed"}

				c.JSON(statusCode, response)
				return
			}
			cacheKey := fmt.Sprintf("[GET] /v1/crypto/withrawal-history/%v/%v", currency, targetPublicKeyForHistory)
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
			cacheDurationInSeconds := 20 //in seconds
			conDB.PrintDBStats(fmt.Sprintf("/v1/crypto/withrawal-history/%v/%v", currency, targetPublicKeyForHistory), gc.DB)

			signerUser, err := usersDB.GetUserFromPrimarySigner(middleware.ExtractSigner(c), gc.DB, gc)

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
					response = gin.H{"error": err.Error(), "message": err.Error()}
				}

				c.JSON(statusCode, response)
				gc.RedisCache.CacheHttpResponseWithParameters(cacheKey, cacheKeyParameters, statusCode, response, cacheDurationInSeconds)
				return
			}
			wallet, temp, err := usersDB.GetWallet(targetPublicKeyForHistory, gc.DB)

			if err != nil {
				log.Println("[GET TARGET USER] error for PUBLIC KEY:", targetPublicKeyForHistory, "error: ", err)

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
				// gc.RedisCache.CacheHttpResponseWithParameters(cacheKey, cacheKeyParameters, statusCode, response, cacheDurationInSeconds)
				return
			}
			if temp {

				statusCode := http.StatusBadRequest
				response := gin.H{"error": "error only main wallets allowed", "message": "Only main wallets are allowed. The address you provided is not a main wallet."}

				c.JSON(statusCode, response)
				return
			}

			targetOwnerUser, err := usersDB.GetUser(targetPublicKeyForHistory, gc.DB, gc)

			if err != nil {
				log.Println("[GET TARGET USER] error for PUBLIC KEY:", targetPublicKeyForHistory, "error: ", err)

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
				// gc.RedisCache.CacheHttpResponseWithParameters(cacheKey, cacheKeyParameters, statusCode, response, cacheDurationInSeconds)
				return
			}

			{
				// gc.RedisCache.InvalidateCachedHttpResponse(cacheKey)

				//check if the owner is the one accessing it or if the one accessing it has access to access it.

				if (signerUser.Username != targetOwnerUser.Username) && !wallet.SignerHasAccess(&signerUser, gc) {
					te := &tErrors.ErrorInvalidAuthorization{}

					log.Println("[GET HISTORY] Invalid access for user:", signerUser.Username, "error: ", te.Error())
					c.JSON(te.HTTPCode(), te.JSONError())
					return
				}

			}
			//Get Payment history
			historyRecords := paymentServices.GetCryptoWithdrawalHistory(targetPublicKeyForHistory, currency, gc, c)

			c.JSON(http.StatusOK, historyRecords)
			// gc.RedisCache.CacheHttpResponse(cacheKey, http.StatusOK, historyRecords, cacheDurationInSeconds)
			gc.RedisCache.CacheHttpResponseWithParameters(cacheKey, cacheKeyParameters, http.StatusOK, historyRecords, cacheDurationInSeconds)

		})

		router.GET("/v1/crypto/deposit-history/:currency/:targetPublicKeyForHistory", middleware.AuthenticationMiddlewareUsingTimestamp(), func(c *gin.Context) {
			// var err error
			currency := strings.ToUpper(c.Param("currency"))

			targetPublicKeyForHistory := strings.TrimSpace(strings.ToUpper(c.Param("targetPublicKeyForHistory")))

			_, err := keypair.ParseAddress(targetPublicKeyForHistory)
			if err != nil {

				statusCode := http.StatusBadRequest
				response := gin.H{"error": "error invalid address", "message": "Only valid addresses are allowed"}

				c.JSON(statusCode, response)
				return
			}
			cacheKey := fmt.Sprintf("[GET] /v1/crypto/deposit-history/%v/%v", currency, targetPublicKeyForHistory)
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
			cacheDurationInSeconds := 20 //in seconds
			conDB.PrintDBStats(fmt.Sprintf("/v1/crypto/deposit-history/%v/%v", currency, targetPublicKeyForHistory), gc.DB)

			signerUser, err := usersDB.GetUserFromPrimarySigner(middleware.ExtractSigner(c), gc.DB, gc)

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
					response = gin.H{"error": err.Error(), "message": err.Error()}
				}

				c.JSON(statusCode, response)
				gc.RedisCache.CacheHttpResponseWithParameters(cacheKey, cacheKeyParameters, statusCode, response, cacheDurationInSeconds)
				return
			}
			wallet, temp, err := usersDB.GetWallet(targetPublicKeyForHistory, gc.DB)

			if err != nil {
				log.Println("[GET TARGET USER] error for PUBLIC KEY:", targetPublicKeyForHistory, "error: ", err)

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
				// gc.RedisCache.CacheHttpResponseWithParameters(cacheKey, cacheKeyParameters, statusCode, response, cacheDurationInSeconds)
				return
			}
			if temp {

				statusCode := http.StatusBadRequest
				response := gin.H{"error": "error only main wallets allowed", "message": "Only main wallets are allowed. The address you provided is not a main wallet."}

				c.JSON(statusCode, response)
				return
			}

			targetOwnerUser, err := usersDB.GetUser(targetPublicKeyForHistory, gc.DB, gc)

			if err != nil {
				log.Println("[GET TARGET USER] error for PUBLIC KEY:", targetPublicKeyForHistory, "error: ", err)

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
				// gc.RedisCache.CacheHttpResponseWithParameters(cacheKey, cacheKeyParameters, statusCode, response, cacheDurationInSeconds)
				return
			}

			{
				// gc.RedisCache.InvalidateCachedHttpResponse(cacheKey)

				//check if the owner is the one accessing it or if the one accessing it has access to access it.

				if (signerUser.Username != targetOwnerUser.Username) && !wallet.SignerHasAccess(&signerUser, gc) {
					te := &tErrors.ErrorInvalidAuthorization{}

					log.Println("[GET HISTORY] Invalid access for user:", signerUser.Username, "error: ", te.Error())
					c.JSON(te.HTTPCode(), te.JSONError())
					return
				}

			}
			//Get Payment history
			historyRecords := paymentServices.GetCryptoDepositHistory(targetPublicKeyForHistory, currency, gc, c)

			c.JSON(http.StatusOK, historyRecords)
			// gc.RedisCache.CacheHttpResponse(cacheKey, http.StatusOK, historyRecords, cacheDurationInSeconds)
			gc.RedisCache.CacheHttpResponseWithParameters(cacheKey, cacheKeyParameters, http.StatusOK, historyRecords, cacheDurationInSeconds)

		})

		//get specific  wallet balance, middleware.AuthenticationMiddlewareUsingTimestamp()
		router.GET("/v1/crypto/withdrawal-networks/:currency", middleware.AuthenticationMiddlewareUsingTimestamp(), func(c *gin.Context) {
			// var err error
			currency := strings.ToUpper(c.Param("currency"))

			wdlNetworks, err := userServices.GetWithdrawalNetworks(currency, gc)
			if err != nil {
				log.Println("[GET Wdl networks] error: ", err)

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
				// gc.RedisCache.CacheHttpResponseWithParameters(cacheKey, cacheKeyParameters, statusCode, response, cacheDurationInSeconds)
				return
			}
			c.JSON(http.StatusOK, gin.H{"networks": wdlNetworks, "serviceFee": os.Getenv("CRYPTO_WITHDRAWAL_SERVICE_FEE")})

		})

		router.POST("/v1/crypto/withdrawals", middleware.AuthenticationMiddlewareUsingTimestamp(), func(c *gin.Context) {
			var err error
			accountSignerUser, getUserError := userModels.UserSigner(middleware.ExtractSigner(c)).GetOwner(gc.DB, gc)

			if getUserError != nil {
				log.Printf("[GENERATE DEPOSIT ADDRESS] ERROR GETTING USER FROM DB from [%v], error: [%v]\n", middleware.ExtractSigner(c), getUserError)

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

			walletOwner, err := usersDB.GetUser(middleware.ExtractPublicKey(c), gc.DB, gc)

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

			var wdlInput userModels.WithdrawalRequestInput

			data, _ := io.ReadAll(c.Request.Body)
			// log.Println(string(data))
			err = json.Unmarshal(data, &wdlInput)

			var invalidJSON tErrors.ErrorInvalidJSON

			if err != nil {
				c.JSON(http.StatusBadRequest, invalidJSON.JSONError())
				return
			}
			wdlInput.Currency = strings.ToUpper(wdlInput.Currency)
			wallet, temp, err := usersDB.GetWallet(middleware.ExtractPublicKey(c), gc.DB)
			if temp {
				c.JSON(http.StatusBadRequest, gin.H{"error": "error-wallet-forbidden", "message": "Wallet forbidden."})
				return
			}

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

			err = userServices.QueueWithdrawalRequest(&accountSignerUser, &wallet, &wdlInput, gc)
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

			if len(wdlInput.TransactionID) == 0 {
				c.JSON(http.StatusAccepted, wdlInput)
				return
			} else {
				c.JSON(http.StatusOK, wdlInput)
			}

			if walletOwner.PushNotificationToken != nil && len(wdlInput.TransactionID) > 0 && wdlInput.TransactionID != "PENDING_AUTH" {
				dataPayload := make(map[string]string)
				dataPayload["route"] = "cryptoHistory"
				walletOwner.SendPushMessage(fmt.Sprintf("%v %v withdrawal on %v has been submitted!", wdlInput.AmountSubmitted, wdlInput.Currency, wallet.Alias), fmt.Sprintf("You have successfully submitted a withdrawal request for %v %v on the wallet with alias [%v].", wdlInput.AmountSubmitted, wdlInput.Currency, wallet.Alias), "", dataPayload, gc)
			}
			walletOwner.InvalidateUserCache(gc)

		})

		router.POST("/v1/shared-access/crypto/withdrawals", middleware.AuthenticationMiddlewareUsingTimestamp(), func(c *gin.Context) {
			var err error
			accountSignerUser, getUserError := userModels.UserSigner(middleware.ExtractSigner(c)).GetOwner(gc.DB, gc)

			if getUserError != nil {
				log.Printf("[GENERATE DEPOSIT ADDRESS] ERROR GETTING USER FROM DB from [%v], error: [%v]\n", middleware.ExtractSigner(c), getUserError)

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

			walletOwner, err := usersDB.GetUser(middleware.ExtractPublicKey(c), gc.DB, gc)

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

			var wdlInput userModels.WithdrawalRequestInput

			data, _ := io.ReadAll(c.Request.Body)
			log.Println(string(data))
			err = json.Unmarshal(data, &wdlInput)

			var invalidJSON tErrors.ErrorInvalidJSON

			if err != nil {
				c.JSON(http.StatusBadRequest, invalidJSON.JSONError())
				return
			}
			wdlInput.Currency = strings.ToUpper(wdlInput.Currency)
			wallet, temp, err := usersDB.GetWallet(middleware.ExtractPublicKey(c), gc.DB)
			if temp {
				c.JSON(http.StatusBadRequest, gin.H{"error": "error-wallet-forbidden", "message": "Wallet forbidden."})
				return
			}

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

			//check if shared wallet, then check if user has access
			if wallet.SharedAccessEnabled == 1 && wallet.NumberOfApprovalsNeeded > 0 {
				//check if signer has access
				hasInitiatorAccess := false
				// check if user has initiator access to wallet.
				for _, p := range accountSignerUser.WalletsSharedWithUser {
					if p.WalletPublicKey == middleware.ExtractPublicKey(c) && p.TargetUsername == accountSignerUser.Username && p.Permission == "INITIATOR" {
						hasInitiatorAccess = true
					}
				}
				if !hasInitiatorAccess {
					c.JSON(http.StatusForbidden, gin.H{"error": "error-unauthorized-access", "message": "You do not have an initiator permission on this wallet."})
					return
				}
			}

			if wallet.HasViewOnlyAccess(gc) {
				if wallet.UserID != accountSignerUser.ID {
					c.JSON(http.StatusForbidden, gin.H{"error": "error-unauthorized-access", "message": "You do not have permission to access this wallet."})
					return
				}
			}

			err = userServices.QueueWithdrawalRequest(&accountSignerUser, &wallet, &wdlInput, gc)
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
			if wdlInput.TransactionID != "PENDING_AUTH" {
				c.JSON(http.StatusAccepted, wdlInput)
				return
			} else {
				c.JSON(http.StatusOK, wdlInput)
			}

			if walletOwner.PushNotificationToken != nil && len(wdlInput.TransactionID) > 0 && wdlInput.TransactionID != "PENDING_AUTH" {
				dataPayload := make(map[string]string)
				dataPayload["route"] = ""
				accountSignerUser.SendPushMessage(fmt.Sprintf("%v %v withdrawal request on %v has been submitted!", wdlInput.AmountSubmitted, wdlInput.Currency, wallet.Alias), fmt.Sprintf("You have successfully submitted a withdrawal request for %v %v on the wallet with alias [%v]. All approvers have been notified.", wdlInput.AmountSubmitted, wdlInput.Currency, wallet.Alias), "", dataPayload, gc)

			}
			{
				//start push notificationMessage

				permissionList := wallet.Permissions
				for _, v := range permissionList {
					u, e := userModels.Username(v.TargetUsername).GetSimpleUser(gc.DB, gc)
					if e != nil {
						continue
					}

					dataPayload := make(map[string]string)
					dataPayload["route"] = "pendingApproval"
					if wdlInput.TransactionID == "PENDING_AUTH" {
						u.SendPushMessage(fmt.Sprintf("%v %v withdrawal request submitted on %v!", wdlInput.AmountSubmitted, wdlInput.Currency, wallet.Alias), fmt.Sprintf("Request:\n %v", wdlInput.ReturnedDescription), "", dataPayload, gc)
					}

				}
			}
			walletOwner.InvalidateUserCache(gc)

		})

		router.POST("/v1/crypto/generate-addresses/:currency", middleware.AuthenticationMiddlewareUsingTimestamp(), func(c *gin.Context) {
			// var err error
			currency := strings.ToUpper(c.Param("currency"))
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
			accountSignerUser, getUserError := userModels.UserSigner(middleware.ExtractSigner(c)).GetOwner(gc.DB, gc)

			if getUserError != nil {
				log.Printf("[GENERATE DEPOSIT ADDRESS] ERROR GETTING USER FROM DB from [%v], error: [%v]\n", middleware.ExtractSigner(c), getUserError)

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
			//check if shared wallet, then check if user has access
			if wallet.SharedAccessEnabled == 1 && wallet.NumberOfApprovalsNeeded > 0 {
				//check if signer has access
				hasInitiatorAccess := false
				// check if user has initiator access to wallet.
				for _, p := range accountSignerUser.WalletsSharedWithUser {
					if p.WalletPublicKey == middleware.ExtractPublicKey(c) && p.TargetUsername == accountSignerUser.Username && p.Permission == "INITIATOR" {
						hasInitiatorAccess = true
					}
				}
				if !hasInitiatorAccess {
					c.JSON(http.StatusForbidden, gin.H{"error": "error-unauthorized-access", "message": "You do not have an initiator permission on this wallet."})
					return
				}
			}

			if wallet.SharedAccessEnabled == 0 || wallet.HasViewOnlyAccess(gc) {
				if wallet.UserID != accountSignerUser.ID {
					c.JSON(http.StatusForbidden, gin.H{"error": "error-unauthorized-access", "message": "You do not have permission to access this wallet."})
					return
				}
			}

			depositAddresses, err := userServices.GenerateDepositAddresses(&wallet, currency, gc)
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
			c.JSON(http.StatusOK, depositAddresses)

		})

	}

}
