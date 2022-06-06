package payments

import (
	"bytes"
	"fmt"
	"os"
	"sync"
	"time"
	paymentErrors "trovo-wallet-api/internal/components/payments/errors"
	paymentmodels "trovo-wallet-api/internal/components/payments/models"
	payments "trovo-wallet-api/internal/components/payments/services"
	usersDB "trovo-wallet-api/internal/components/users/db"
	conDB "trovo-wallet-api/internal/db"
	"trovo-wallet-api/internal/sharedconfig"

	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"trovo-wallet-api/internal/errors"
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

	router.POST("/v2/users/:targetUser/payments", middleware.AuthenticationMiddlewareUsingTimestamp(), func(c *gin.Context) {
		var err error

		suppliedUsername := strings.TrimSpace(strings.ToLower(c.Param("targetUser")))
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
				suppliedUsername = string(uDec)
			}

		}
		if suppliedUsername == "null" {
			log.Printf("user cannot be %v\n", suppliedUsername)
			c.JSON(http.StatusBadRequest, gin.H{"error": "user cannot be null"})
			return
		}
		conDB.PrintDBStats(fmt.Sprintf("POST /v2/users/%v/payments %v", c.Param("targetUser"), suppliedUsername), gc.DB)

		//check if username is reserved. Reserved usernames should not send payments.
		//TODO: cache this
		_, checkReservedUserError := usersDB.UsernameIsReserved(suppliedUsername, gc.DB)
		if checkReservedUserError != nil {

			var ex errors.GenericError
			var ok bool

			ex, ok = checkReservedUserError.(errors.GenericError)
			if ok {
				c.JSON(ex.HTTPCode(), ex.JSONError())
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": checkReservedUserError.Error()})
			}
			return
		}

		userInfo, getUserError := usersDB.GetUser(suppliedUsername, gc.DB)
		if suppliedUsername == os.Getenv("LOG_TARGET_USER") || middleware.ExtractPublicKey(c) == os.Getenv("LOG_TARGET_USER_PK") {
			log.Printf("[CUSTOM LOG] %v error:%v\n", suppliedUsername, getUserError)
		}
		if getUserError != nil {
			log.Printf("[FAILIED PAYMENT] ERROR GETTING USER INFO ()()()()()()()()() from [%v], error: [%v]\n", suppliedUsername, getUserError)

			var ex errors.GenericError
			var ok bool

			ex, ok = getUserError.(errors.GenericError)
			if ok {
				c.JSON(ex.HTTPCode(), ex.JSONError())
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": getUserError.Error()})
			}
			return
		}
		if userInfo.Suspended == 1 {
			getUserError = &errors.ErrorUsernameIsSuspended{}

			var ex errors.GenericError
			var ok bool

			ex, ok = getUserError.(errors.GenericError)
			if ok {
				c.JSON(ex.HTTPCode(), ex.JSONError())
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": getUserError.Error()})
			}
			return

		}

		suppliedPublicKey := userInfo.PublicKey

		if suppliedPublicKey != middleware.ExtractPublicKey(c) {
			log.Printf("[FAILIED PAYMENT] INVALID PAYMENT SENDER PUBLIC KEY ()()()()()()()()() from [%v], error: [%v]\n", suppliedUsername, err)
			c.JSON(http.StatusBadRequest, (&paymentErrors.ErrorInvalidPaymentSender{}).JSONError())
			return
		}

		var paymentInfo paymentmodels.PaymentInfo
		// var err error

		data, _ := ioutil.ReadAll(c.Request.Body)

		err = json.Unmarshal(data, &paymentInfo)

		var invalidJSON errors.ErrorInvalidJSON

		if err != nil {
			log.Printf("[FAILIED PAYMENT] UNMARSHAL FAIL from [%v], error: [%v]\n", suppliedUsername, err)

			log.Print(err)
			c.JSON(invalidJSON.HTTPCode(), invalidJSON.JSONError())
			return
		}
		if suppliedUsername == os.Getenv("LOG_TARGET_USER") || middleware.ExtractPublicKey(c) == os.Getenv("LOG_TARGET_USER_PK") {
			log.Printf("[CUSTOM LOG] paymentInfo %+v\n", paymentInfo)
		}
		//check if the public key exists in BUDS and then transform to username
		paymentInfo.Messages = make([]string, 0)
		if len(paymentInfo.Destination) == 56 {
			tUser, uerr := usersDB.GetUser(paymentInfo.Destination, gc.DB)
			if uerr == nil {
				paymentInfo.Messages = append(paymentInfo.Messages, fmt.Sprintf("Notice: TrovoWallet Address[%v] belongs to username [%v]", paymentInfo.Destination, tUser.Username))
				paymentInfo.Destination = tUser.Username
			}
		}

		//check if receiver is reserved. Reserved usernames should not be sent payments.
		if len(paymentInfo.Destination) != 56 {
			//skip public key payments
			_, checkReservedReceiverError := usersDB.UsernameIsReserved(paymentInfo.Destination, gc.DB)
			if checkReservedReceiverError != nil {

				var ex errors.GenericError
				var ok bool

				ex, ok = checkReservedReceiverError.(errors.GenericError)
				if ok {
					c.JSON(ex.HTTPCode(), ex.JSONError())
				} else {
					c.JSON(http.StatusBadRequest, gin.H{"error": checkReservedReceiverError.Error()})
				}
				return
			}
		}
		paymentInfoReturned, destinationUser, paymentError := payments.Pay(&userInfo, middleware.ExtractPublicKey(c), &paymentInfo, gc.DB)
		if suppliedUsername == os.Getenv("LOG_TARGET_USER") || middleware.ExtractPublicKey(c) == os.Getenv("LOG_TARGET_USER_PK") {
			log.Printf("[CUSTOM LOG] returned Payment Error: [%v]\n", paymentError)

			if paymentInfoReturned != nil {
				log.Printf("[CUSTOM LOG] returned Payment Info: [%+v]\n", *paymentInfoReturned)
			}

		}
		if paymentError != nil {
			log.Printf("[FAILIED PAYMENT]  from [%v] to [%v], error: [%v]\n", suppliedUsername, paymentInfo.Destination, paymentError)
			var ex errors.GenericError
			var ok bool

			ex, ok = paymentError.(errors.GenericError)
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
			userInfo.PublicIP = c.ClientIP()
			payments.UpdateAndLogUserPaymentGeoInfomation(&userInfo, paymentInfoReturned, gc.DB)
			senderPaymentHistoryCacheKey := fmt.Sprintf("[GET] /v2/users/%v/payments", suppliedUsername)
			senderCacheKey := fmt.Sprintf("[GET] /v2/users/%v", suppliedUsername)

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
						Sender:          suppliedUsername,
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
					destinationUser.SendPushMessage("Trovo: Wallet Credited!", fmt.Sprintf("You have received %v %v from %v", paymentInfo.Amount, assetCode, userInfo.Username), "", gc)
				}
				userInfo.SendPushMessage("Trovo: Wallet Debited!", fmt.Sprintf("You have successfully sent %v %v to %v", paymentInfo.Amount, assetCode, paymentInfo.Destination), "", gc)

			}

		}

	})

}
