package payments

import (
	"sync"
	"time"
	userModels "trovo-wallet-api/internal/components/users/models"
	"trovo-wallet-api/internal/sharedconfig"

	"log"
	"net/http"
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

	router.POST("/v1/users/payment", middleware.AuthenticationMiddlewareUsingTimestamp(), postUsersPaymentHandler(callBackRetryChan, gc))
	router.POST("/v1/shared-access/payment", middleware.AuthenticationMiddlewareUsingTimestamp(), postSharedAccessPaymentHandler(callBackRetryChan, gc))
}
