package payments

import (
	"bytes"
	"fmt"
	"os"
	"sync"
	"time"
	paymentsDB "trovo-wallet-api/internal/components/payments/db"
	tPayErrors "trovo-wallet-api/internal/components/payments/errors"
	paymentModels "trovo-wallet-api/internal/components/payments/models"
	payments "trovo-wallet-api/internal/components/payments/services"
	conDB "trovo-wallet-api/internal/db"
	"trovo-wallet-api/internal/sharedconfig"

	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

//Init initializes the controller
func Init(router *gin.Engine, gc *sharedconfig.GlobalConfig) {
	//start routine to resend failed payment callbacks
	// type retrySling struct {
	// 	Req   *sling.Sling
	// 	Count int
	// }

	//retryCallbacks stores failed callbacks
	type retryCallbacks struct {
		Req         *bytes.Buffer
		CallbackURL string
		Count       int
	}

	callBackRetryChan := make(chan retryCallbacks, 200000)
	go func(c chan retryCallbacks) {
		maxCallbackCount := 10
		//loop
		var workerWaitGroup sync.WaitGroup
		numWorkers := 5

		notificationWorker := func(wm *sync.WaitGroup, wid int, callbackChannel chan retryCallbacks) {
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

	router.POST("/v1/users/:primaryAccountAlias/payments", middleware.AuthenticationMiddlewareUsingTimestamp(), func(c *gin.Context) {
		var err error

		primaryAccountAlias := strings.TrimSpace(strings.ToLower(c.Param("primaryAccountAlias")))
		uDec, e := base64.URLEncoding.DecodeString(c.Param("primaryAccountAlias"))
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
				primaryAccountAlias = string(uDec)
			}

		}
		if primaryAccountAlias == "null" {
			log.Printf("[POST PAYMENTS]user cannot be %v\n", primaryAccountAlias)
			c.JSON(http.StatusBadRequest, gin.H{"error": "user cannot be null"})
			return
		}
		conDB.PrintDBStats(fmt.Sprintf("POST /v1/users/%v/payments %v", c.Param("primaryAccountAlias"), primaryAccountAlias), gc.DB)

		//check if username is reserved. Reserved usernames should not send payments.
		//TODO: cache this
		_, checkReservedUserError := paymentsDB.UsernameIsReserved(primaryAccountAlias, gc.DB)
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
		userWallet, temp, getWalletError := paymentsDB.GetWallet(middleware.ExtractPublicKey(c), gc.DB)
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
				c.JSON(http.StatusBadRequest, gin.H{"error": getWalletError.Error()})
			}
			return
		}

		if temp || userWallet.Tag != nil {
			errAccountIsTemp := &tErrors.CustomError{
				Param:      "Username",
				Err:        "error-account-not-primary-account-alias",
				ErrMessage: "only sender primary accounts are allowed for payment requests",
				Code:       http.StatusForbidden,
			}

			c.JSON(errAccountIsTemp.HTTPCode(), errAccountIsTemp.JSONError())
			return

		}

		//get user DB record
		owner, getUserError := paymentsDB.GetUser(primaryAccountAlias, gc.DB)
		if primaryAccountAlias == os.Getenv("LOG_TARGET_USER") || middleware.ExtractPublicKey(c) == os.Getenv("LOG_TARGET_USER_PK") {
			log.Printf("[CUSTOM LOG] %v error:%v\n", primaryAccountAlias, getUserError)
		}
		if getUserError != nil {
			log.Printf("[FAILED PAYMENT] ERROR GETTING USER FROM DB from [%v], error: [%v]\n", primaryAccountAlias, getUserError)

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

		if owner.Suspended == 1 {
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

		primaryAccountSigner := owner.PublicKey

		if primaryAccountSigner != middleware.ExtractSigner(c) {
			log.Printf("[FAILED PAYMENT] INVALID PAYMENT SIGNER IN HEADER from [%v], error: [%v]\n", primaryAccountAlias, err)
			c.JSON(http.StatusBadRequest, (&tPayErrors.ErrorInvalidPaymentSender{}).JSONError())
			return
		}

		var paymentInfo paymentModels.PaymentInfo
		// var err error

		data, _ := ioutil.ReadAll(c.Request.Body)

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

		var destinationUser paymentsDB.User
		var destinationWallet paymentsDB.UserWallet
		var getDestinationUserError, getDestinationWalletError error
		//check if the public key exists in TROVO and then transform to username
		paymentInfo.Messages = make([]string, 0)
		if len(paymentInfo.Destination) == 56 {
			destinationWallet, _, getDestinationWalletError = paymentsDB.GetWallet(paymentInfo.Destination, gc.DB)
			if getDestinationWalletError == nil {
				paymentInfo.Messages = append(paymentInfo.Messages, fmt.Sprintf("Notice: Bantu Address[%v] belongs to the wallet alias [%v]", paymentInfo.Destination, destinationWallet.Alias))
				paymentInfo.Destination = destinationWallet.Alias
			}
		}

		//check if receiver is reserved. Reserved usernames should not be sent payments.

		if len(paymentInfo.Destination) != 56 {
			//skip public key payments
			_, checkReservedReceiverError := paymentsDB.UsernameIsReserved(paymentInfo.Destination, gc.DB)
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

			destinationUser, getDestinationUserError = paymentsDB.GetUser(paymentInfo.Destination, gc.DB)
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
		paymentInfoReturned, destinationUser, paymentError := payments.Pay(&owner, &userWallet, &paymentInfo, gc.DB)
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
			owner.PublicIP = c.ClientIP()
			
			payments.UpdateAndLogUserPaymentGeoInformation(&owner, paymentInfoReturned, gc.DB)
			senderPaymentHistoryCacheKey := fmt.Sprintf("[GET] /v2/users/%v/payments", primaryAccountAlias)
			senderCacheKey := fmt.Sprintf("[GET] /v2/users/%v", primaryAccountAlias)

			gc.RedisCache.InvalidateCachedHttpResponse(senderCacheKey, senderPaymentHistoryCacheKey)
			gc.RedisCache.InvalidateCachedHttpResponse(senderPaymentHistoryCacheKey)
			if destinationUser != nil {
				destinationUsername := strings.TrimSpace(strings.ToLower(destinationUser.Username))

				receiverCacheKey := fmt.Sprintf("[GET] /v2/users/%v", destinationUsername)
				receiverPaymentHistoryCacheKey := fmt.Sprintf("[GET] /v2/users/%v/payments", destinationUsername)
				gc.RedisCache.InvalidateCachedHttpResponse(receiverCacheKey, receiverPaymentHistoryCacheKey)
			}
			c.JSON(http.StatusOK, paymentInfoReturned)

			{

				//start callback process here

				//TODO: make callback request if callback is availble

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
						assetCode = "XBN"
					}
					jsonPayload := payload{
						Destination:     paymentInfoReturned.Destination,
						Sender:          primaryAccountAlias,
						Amount:          paymentInfoReturned.Amount,
						AssetCode:       assetCode,
						AssetIssuer:     paymentInfoReturned.AssetIssuer,
						TransactionID:   paymentInfoReturned.TransactionID,
						TransactionMemo: paymentInfoReturned.Memo,
						TransactionTime: time.Now(),
						DeviceID:        c.GetHeader("X-TrovoWallet-DEVICE-ID"),
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
					assetCode = "XBN"
				}
				if destinationUser != nil {
					destinationUser.SendPushMessage("Trovo: Wallet Credited!", fmt.Sprintf("You have received %v %v from %v", paymentInfo.Amount, assetCode, userDB.Username), "", gc)
				}
				userDB.SendPushMessage("Trovo: Wallet Debited!", fmt.Sprintf("You have successfully sent %v %v to %v", paymentInfo.Amount, assetCode, paymentInfo.Destination), "", gc)

			}

		}

	})

}
