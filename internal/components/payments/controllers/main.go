package payments

import (
	"bytes"
	"fmt"
	"os"
	"sync"
	"time"
	tPayErrors "trovo-wallet-api/internal/components/payments/errors"
	paymentModels "trovo-wallet-api/internal/components/payments/models"
	payments "trovo-wallet-api/internal/components/payments/services"
	usersDB "trovo-wallet-api/internal/components/users/db"
	userModels "trovo-wallet-api/internal/components/users/models"
	userServices "trovo-wallet-api/internal/components/users/services"
	conDB "trovo-wallet-api/internal/db"
	"trovo-wallet-api/internal/sharedconfig"

	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

// Init initializes the controller
func Init(router *gin.Engine, callBackRetryChan chan userModels.RetryCallbacks, gc *sharedconfig.GlobalConfig) {
	//start routine to resend failed payment callbacks
	// type retrySling struct {
	// 	Req   *sling.Sling
	// 	Count int
	// }

	go func(c chan userModels.RetryCallbacks) {
		maxCallbackCount := 10
		//loop
		var workerWaitGroup sync.WaitGroup
		numWorkers := 5

		notificationWorker := func(wm *sync.WaitGroup, wid int, callbackChannel chan userModels.RetryCallbacks) {
			log.Println("[paymentNotification] started worker -> ", wid)
			//start worker here
			defer wm.Done()
			for {
				callbackObj := <-callbackChannel
				if callbackObj.Count > maxCallbackCount {
					//skip retries
					continue
				}
				resp, err := http.Post(callbackObj.CallbackURL, "application/json", callbackObj.Req)

				if err != nil {
					//send back into channel to retry later
					log.Printf("Callback retry failed: [%+v], Worker ID: [%v]\n", callbackObj, wid)
					if callbackObj.Count <= maxCallbackCount {
						callbackObj.Count++
						time.Sleep(500 * time.Millisecond)
						callbackChannel <- callbackObj
					}
					continue
				}
				if resp.StatusCode >= 300 {
					//send back into channel to retry later
					log.Printf("Callback retry failed: [%+v], Worker ID: [%v]\n", callbackObj, wid)
					if callbackObj.Count <= maxCallbackCount {
						callbackObj.Count++
						time.Sleep(500 * time.Millisecond)
						callbackChannel <- callbackObj
					}
					continue
				}
				log.Printf("[paymentNotification]Worker[%v] ######@@@@######@@ successful to: [%v]\n", wid, callbackObj.CallbackURL)

			}
			//loop work
		}

		for i := 0; i < numWorkers; i++ {
			workerWaitGroup.Add(1)
			go notificationWorker(&workerWaitGroup, i+1, c)
		}
		log.Println("#####@started Routine to retry failed Payment callbacks....")

		workerWaitGroup.Wait()

	}(callBackRetryChan)

	router.POST("/v1/users/payment", middleware.AuthenticationMiddlewareUsingTimestamp(), func(c *gin.Context) {
		var err error

		//get user DB record
		accountSignerUser, getUserError := userModels.UserSigner(middleware.ExtractSigner(c)).GetOwner(gc.DB)

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
		primaryAccountAlias := accountSignerUser.Username
		if primaryAccountAlias == os.Getenv("LOG_TARGET_USER") || middleware.ExtractPublicKey(c) == os.Getenv("LOG_TARGET_USER_PK") {
			log.Printf("[CUSTOM LOG] %v error:%v\n", primaryAccountAlias, getUserError)
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
			log.Printf("[FAILED PAYMENT] UNMARSHAL FAIL from [%v], error: [%v]\n", primaryAccountAlias, err)

			log.Print(err)
			c.JSON(invalidJSON.HTTPCode(), invalidJSON.JSONError())
			return
		}
		if primaryAccountAlias == os.Getenv("LOG_TARGET_USER") || middleware.ExtractPublicKey(c) == os.Getenv("LOG_TARGET_USER_PK") {
			log.Printf("[CUSTOM LOG] paymentInfo %+v\n", paymentInfo)
		}

		conDB.PrintDBStats(fmt.Sprintf("POST /v1/users/%v/payments %v", c.Param("primaryAccountAlias"), primaryAccountAlias), gc.DB)

		//check if username is reserved. Reserved usernames should not send payments.
		//TODO: cache this
		_, checkReservedUserError := usersDB.UsernameIsReserved(primaryAccountAlias, gc.DB)
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
		userWallet, temp, getWalletError := usersDB.GetWallet(middleware.ExtractPublicKey(c), gc.DB)
		if primaryAccountAlias == os.Getenv("LOG_TARGET_USER") || middleware.ExtractPublicKey(c) == os.Getenv("LOG_TARGET_USER_PK") {
			log.Printf("[CUSTOM LOG] %v error:%v\n", primaryAccountAlias, getWalletError)
		}
		if getWalletError != nil {
			log.Printf("[FAILED PAYMENT] ERROR GETTING USER FROM DB from [%v], error: [%v]\n", primaryAccountAlias, getWalletError)

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
				Err:        "error-account-not-primary-account-alias",
				ErrMessage: "only primary/subwallets are allowed for payment requests",
				Code:       http.StatusForbidden,
			}

			c.JSON(errAccountIsTemp.HTTPCode(), errAccountIsTemp.JSONError())
			return

		}
		if userWallet.WalletType != 0 {
			c.JSON(http.StatusForbidden, gin.H{"error": "error-wallet-type-forbidden", "message": "Operation not allowed on any special type of wallets. Only standard wallets are allowed."})
			return
		}
		{
			//prevent wallets with approver from using this endpoint
			if userWallet.SharedAccessEnabled == 1 && userWallet.WalletCountApproverAccess(gc) > 0 {
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
		if len(paymentInfo.Destination) == 56 || len(paymentInfo.Destination) == 69 {
			destinationWallet, _, getDestinationWalletError = usersDB.GetWallet(paymentInfo.Destination, gc.DB)
			if getDestinationWalletError == nil {
				paymentInfo.Messages = append(paymentInfo.Messages, fmt.Sprintf("Notice: Address[%v] belongs to the wallet alias [%v]", paymentInfo.Destination, destinationWallet.Alias))
				paymentInfo.Destination = destinationWallet.Alias
			}
		}

		//check if receiver is reserved. Reserved usernames should not be sent payments.

		if len(paymentInfo.Destination) != 56 && len(paymentInfo.Destination) != 69 {
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

			destinationUser, getDestinationUserError = usersDB.GetUser(paymentInfo.Destination, gc.DB)
			if getDestinationUserError != nil {
				ex := &tPayErrors.ErrorPaymentDestinationDoesNotExist{}
				c.JSON(ex.HTTPCode(), ex.JSONError())
				return
			}
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
		}
		paymentInfoReturned, returnedDestination, paymentError := userServices.Pay(&accountSignerUser, &userWallet, &paymentInfo, gc)
		if primaryAccountAlias == os.Getenv("LOG_TARGET_USER") || middleware.ExtractPublicKey(c) == os.Getenv("LOG_TARGET_USER_PK") {
			log.Printf("[CUSTOM LOG] returned Payment Error: [%v]\n", paymentError)

			if paymentInfoReturned != nil {
				log.Printf("[CUSTOM LOG] returned Payment Info: [%+v]\n", *paymentInfoReturned)
			}

		}
		if paymentError != nil {
			log.Printf("[FAILED PAYMENT]  from [%v] to [%v], error: [%v]\n", primaryAccountAlias, paymentInfo.Destination, paymentError)
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

			payments.UpdateAndLogUserPaymentGeoInformation(&accountSignerUser, paymentInfoReturned, gc.DB)
			senderPaymentHistoryCacheKey := fmt.Sprintf("[GET] /v1/users/payments/%v", middleware.ExtractPublicKey(c))
			senderCacheKey := fmt.Sprintf("[GET] /v1/users/%v", primaryAccountAlias)

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
				receiverTempCacheKey = fmt.Sprintf("GetBalance_%s", *destinationWallet.TempPublicKey)
			}
			if len(userWallet.ID) == 56 {
				senderTempCacheKey = fmt.Sprintf("GetBalance_%s", *userWallet.TempPublicKey)

			}

			if returnedDestination != nil {
				destinationUsername := strings.TrimSpace(strings.ToLower(destinationUser.Username))

				receiverCacheKey := fmt.Sprintf("[GET] /v1/users/%v", destinationUsername)
				receiverPaymentHistoryCacheKey := fmt.Sprintf("[GET] /v1/users/payments/%v", middleware.ExtractPublicKey(c))
				gc.RedisCache.InvalidateCachedHttpResponse(receiverCacheKey, receiverPaymentHistoryCacheKey)
				gc.RedisCache.InvalidateCachedHttpResponse(receiverPaymentHistoryCacheKey, senderBalanceCacheKey, receiverBalanceCacheKey)
			}
			gc.RedisCache.InvalidateCachedHttpResponse(senderBalanceCacheKey, senderTempCacheKey, receiverBalanceCacheKey, receiverTempCacheKey, sNFT, rNTF)

			if paymentInfoReturned.TransactionID == "PENDING_AUTH" {
				c.JSON(http.StatusOK, paymentInfoReturned)

				{
					accessList := userWallet.GetPermissionList(gc.DB)
					// send push notifications
					assetCode := paymentInfo.AssetCode
					if assetCode == "" {
						assetCode = os.Getenv("NATIVE_ASSET_CODE")
					}
					dataPayload := make(map[string]string)
					dataPayload["route"] = "pendingAuth"
					for _, a := range accessList {

						if a.Permission == "APPROVER" {
							ph, e := usersDB.GetUser(a.TargetUsername, gc.DB)
							if e == nil {
								ph.SendPushMessage("Trovo: Payment request awaiting approval!", fmt.Sprintf("You have a payment transaction of %v %v to %v initiated by %v from the wallet with alias %v, which is now awaiting approval from you or any other approver.", paymentInfo.Amount, assetCode, paymentInfo.Destination, accountSignerUser.Username, userWallet.Alias), "", dataPayload, gc)

							}

						}
					}
					accountSignerUser.SendPushMessage("Trovo: Payment Request Submitted!", fmt.Sprintf("You have successfully submitted payment request of %v %v to %v on the wallet with alias %v. The listed approvers have been notifed to attend to the request.", paymentInfo.Amount, assetCode, paymentInfo.Destination, userWallet.Alias), "", dataPayload, gc)

				}

				return
			}

			c.JSON(http.StatusOK, paymentInfoReturned)

			{

				//start callback process here

				//TODO: make callback request if callback is available

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
					c := userModels.RetryCallbacks{Req: responseBody, CallbackURL: d, Count: 0}
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
				if getDestinationWalletError == nil {
					dataPayload := make(map[string]string)
					dataPayload["none"] = ""
					if destinationWallet.SharedAccessEnabled == 1 {
						if destinationWallet.HasViewOnlyAccess(gc) {
							u, e := destinationWallet.GetWalletOwner(gc.DB)
							if e == nil {
								if u.PushNotificationToken != nil {

									u.SendPushMessage("Trovo: Shared Wallet Credited!", fmt.Sprintf("You have received %v %v from %v to your shared wallet with alias %v", paymentInfo.Amount, assetCode, userWallet.Alias, paymentInfo.Destination), "", dataPayload, gc)

								}
							}

						}
						for _, v := range destinationWallet.Permissions {
							u, e := userModels.Username(v.TargetUsername).GetSimpleUser(gc.DB)
							if e != nil {
								continue
							}
							if u.PushNotificationToken != nil {

								u.SendPushMessage("Trovo: Shared Wallet Credited!", fmt.Sprintf("You have received %v %v from %v to your shared wallet with alias %v", paymentInfo.Amount, assetCode, userWallet.Alias, paymentInfo.Destination), "", dataPayload, gc)

							}
						}
					}
				}
				accountSignerUser.SendPushMessage("Trovo: Wallet Debited!", fmt.Sprintf("You have successfully sent %v %v from your wallet with alias %v to %v", paymentInfo.Amount, assetCode, userWallet.Alias, paymentInfo.Destination), "", dataPayload, gc)

			}

		}

	})

	router.POST("/v1/shared-access/payment", middleware.AuthenticationMiddlewareUsingTimestamp(), func(c *gin.Context) {
		var err error
		{
			//check if pending shared access modify exists
			if userServices.CheckPendingSharedAccessApproval(middleware.ExtractPublicKey(c), gc.DB) {
				c.JSON(http.StatusForbidden, gin.H{"error": "error-pending-shared-access-op", "message": "There is a pending shared access operation on this wallet and must be completed first before attempting to send payment from this wallet."})
				return
			}
		}
		//get user DB record
		accountSignerUser, getUserError := userModels.UserSigner(middleware.ExtractSigner(c)).GetOwner(gc.DB)

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
		hasInitiatorAccess := false
		// check if user has initiator access to wallet.
		for _, p := range accountSignerUser.WalletsSharedWithUser {
			if p.WalletPublicKey == middleware.ExtractPublicKey(c) && p.TargetUsername == accountSignerUser.Username && p.Permission == "INITIATOR" {
				hasInitiatorAccess = true
			}
		}
		if !hasInitiatorAccess && !userModels.UserWalletID(middleware.ExtractPublicKey(c)).PublicKeyHasViewOnlyAccess(gc) {
			c.JSON(http.StatusForbidden, gin.H{"error": "error-unauthorized-access", "message": "You do not have an initiator permission on this wallet."})
			return
		}
		primaryAccountAlias := accountSignerUser.Username
		if primaryAccountAlias == os.Getenv("LOG_TARGET_USER") || middleware.ExtractPublicKey(c) == os.Getenv("LOG_TARGET_USER_PK") {
			log.Printf("[CUSTOM LOG] %v error:%v\n", primaryAccountAlias, getUserError)
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

		var paymentInfo paymentModels.PaymentInfo
		// var err error

		data, _ := io.ReadAll(c.Request.Body)

		err = json.Unmarshal(data, &paymentInfo)

		var invalidJSON tErrors.ErrorInvalidJSON

		if err != nil {
			log.Printf("[FAILED PAYMENT] UNMARSHAL FAIL from [%v], error: [%v]\n", primaryAccountAlias, err)

			log.Print(err)
			c.JSON(invalidJSON.HTTPCode(), invalidJSON.JSONError())
			return
		}
		if primaryAccountAlias == os.Getenv("LOG_TARGET_USER") || middleware.ExtractPublicKey(c) == os.Getenv("LOG_TARGET_USER_PK") {
			log.Printf("[CUSTOM LOG] paymentInfo %+v\n", paymentInfo)
		}

		conDB.PrintDBStats(fmt.Sprintf("POST /v1/users/%v/payments %v", c.Param("primaryAccountAlias"), primaryAccountAlias), gc.DB)

		//check if username is reserved. Reserved usernames should not send payments.
		//TODO: cache this
		_, checkReservedUserError := usersDB.UsernameIsReserved(primaryAccountAlias, gc.DB)
		if checkReservedUserError != nil {

			var ex tErrors.GenericError
			var ok bool

			ex, ok = checkReservedUserError.(tErrors.GenericError)
			if ok {
				c.JSON(ex.HTTPCode(), ex.JSONError())
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": checkReservedUserError.Error(), "message": checkReservedUserError.Error()})
			}
			return
		}
		//get the wallet you are sending payment from
		userWallet, temp, getWalletError := usersDB.GetWallet(middleware.ExtractPublicKey(c), gc.DB)
		if primaryAccountAlias == os.Getenv("LOG_TARGET_USER") || middleware.ExtractPublicKey(c) == os.Getenv("LOG_TARGET_USER_PK") {
			log.Printf("[CUSTOM LOG] %v error:%v\n", primaryAccountAlias, getWalletError)
		}
		if getWalletError != nil {
			log.Printf("[FAILED PAYMENT] ERROR GETTING USER FROM DB from [%v], error: [%v]\n", primaryAccountAlias, getWalletError)

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
				Err:        "error-account-not-primary-account-alias",
				ErrMessage: "only primary/subwallets are allowed for payment requests",
				Code:       http.StatusForbidden,
			}

			c.JSON(errAccountIsTemp.HTTPCode(), errAccountIsTemp.JSONError())
			return

		}
		if userWallet.WalletType != 0 {
			c.JSON(http.StatusForbidden, gin.H{"error": "error-wallet-type-forbidden", "message": "Operation not allowed on any special type of wallets. Only standard wallets are allowed."})
			return
		}
		{
			//prevent wallets with approver from using this endpoint
			if userWallet.SharedAccessEnabled == 0 {
				errAccountIsTemp := &tErrors.CustomError{
					Param:      "ID",
					Err:        "error-wallet-without-shared-access-not-allowed",
					ErrMessage: "This wallet does not have shared access enabled to use this endpoint.",
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
		if len(paymentInfo.Destination) == 56 || len(paymentInfo.Destination) == 69 {
			destinationWallet, _, getDestinationWalletError = usersDB.GetWallet(paymentInfo.Destination, gc.DB)
			if getDestinationWalletError == nil {
				paymentInfo.Messages = append(paymentInfo.Messages, fmt.Sprintf("Notice: Bantu Address[%v] belongs to the wallet alias [%v]", paymentInfo.Destination, destinationWallet.Alias))
				paymentInfo.Destination = destinationWallet.Alias
			}
		}

		//check if receiver is reserved. Reserved usernames should not be sent payments.

		if len(paymentInfo.Destination) != 56 && len(paymentInfo.Destination) != 69 {
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

			destinationUser, getDestinationUserError = usersDB.GetUser(paymentInfo.Destination, gc.DB)
			if getDestinationUserError != nil {
				ex := &tPayErrors.ErrorPaymentDestinationDoesNotExist{}
				c.JSON(ex.HTTPCode(), ex.JSONError())
				return
			}
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
		}
		paymentInfoReturned, returnedDestination, paymentError := userServices.Pay(&accountSignerUser, &userWallet, &paymentInfo, gc)
		if primaryAccountAlias == os.Getenv("LOG_TARGET_USER") || middleware.ExtractPublicKey(c) == os.Getenv("LOG_TARGET_USER_PK") {
			log.Printf("[CUSTOM LOG] returned Payment Error: [%v]\n", paymentError)

			if paymentInfoReturned != nil {
				log.Printf("[CUSTOM LOG] returned Payment Info: [%+v]\n", *paymentInfoReturned)
			}

		}
		if paymentError != nil {
			log.Printf("[FAILED PAYMENT]  from [%v] to [%v], error: [%v]\n", primaryAccountAlias, paymentInfo.Destination, paymentError)
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

			payments.UpdateAndLogUserPaymentGeoInformation(&accountSignerUser, paymentInfoReturned, gc.DB)
			senderPaymentHistoryCacheKey := fmt.Sprintf("[GET] /v1/users/payments/%v", middleware.ExtractPublicKey(c))
			senderCacheKey := fmt.Sprintf("[GET] /v1/users/%v", primaryAccountAlias)

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
				receiverTempCacheKey = fmt.Sprintf("GetBalance_%s", *destinationWallet.TempPublicKey)
			}
			if len(userWallet.ID) == 56 {
				senderTempCacheKey = fmt.Sprintf("GetBalance_%s", *userWallet.TempPublicKey)

			}

			if returnedDestination != nil {
				destinationUsername := strings.TrimSpace(strings.ToLower(destinationUser.Username))

				receiverCacheKey := fmt.Sprintf("[GET] /v1/users/%v", destinationUsername)
				receiverPaymentHistoryCacheKey := fmt.Sprintf("[GET] /v1/users/payments/%v", middleware.ExtractPublicKey(c))
				gc.RedisCache.InvalidateCachedHttpResponse(receiverCacheKey, receiverPaymentHistoryCacheKey)
				gc.RedisCache.InvalidateCachedHttpResponse(receiverPaymentHistoryCacheKey, senderBalanceCacheKey, receiverBalanceCacheKey)
			}
			gc.RedisCache.InvalidateCachedHttpResponse(senderBalanceCacheKey, senderTempCacheKey, receiverBalanceCacheKey, receiverTempCacheKey, sNFT, rNTF)

			if paymentInfoReturned.TransactionID == "PENDING_AUTH" {
				c.JSON(http.StatusOK, paymentInfoReturned)

				{
					accessList := userWallet.GetPermissionList(gc.DB)
					// send push notifications
					assetCode := paymentInfo.AssetCode
					if assetCode == "" {
						assetCode = os.Getenv("NATIVE_ASSET_CODE")
					}
					dataPayload := make(map[string]string)
					dataPayload["route"] = "pendingAuth"
					for _, a := range accessList {

						if a.Permission == "APPROVER" {
							ph, e := usersDB.GetUser(a.TargetUsername, gc.DB)
							if e == nil {
								ph.SendPushMessage("Trovo: Payment request awaiting approval!", fmt.Sprintf("You have a payment transaction of %v %v to %v initiated by %v from the wallet with alias %v, which is now awaiting approval from you or any other approver.", paymentInfo.Amount, assetCode, paymentInfo.Destination, accountSignerUser.Username, userWallet.Alias), "", dataPayload, gc)

							}

						}
					}
					accountSignerUser.SendPushMessage("Trovo: Payment Request Submitted!", fmt.Sprintf("You have successfully submitted payment request of %v %v to %v on the wallet with alias %v. The listed approvers have been notifed to attend to the request.", paymentInfo.Amount, assetCode, paymentInfo.Destination, userWallet.Alias), "", dataPayload, gc)

				}

				return
			}

			c.JSON(http.StatusOK, paymentInfoReturned)

			{

				//start callback process here

				//TODO: make callback request if callback is available

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
					c := userModels.RetryCallbacks{Req: responseBody, CallbackURL: d, Count: 0}
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
				if getDestinationWalletError == nil {
					destinationUser.SendPushMessage("Trovo: Wallet Credited!", fmt.Sprintf("You have received %v %v from %v to your wallet with alias %v", paymentInfo.Amount, assetCode, userWallet.Alias, paymentInfo.Destination), "", dataPayload, gc)
				}
				accountSignerUser.SendPushMessage("Trovo: Wallet Debited!", fmt.Sprintf("You have successfully sent %v %v from your wallet with alias %v to %v", paymentInfo.Amount, assetCode, userWallet.Alias, paymentInfo.Destination), "", dataPayload, gc)

			}

		}

	})

}
