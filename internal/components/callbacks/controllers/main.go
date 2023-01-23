package callbacks

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	userModels "trovo-wallet-api/internal/components/users/models"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/gin-gonic/gin"
)

func Init(router *gin.Engine, callBackRetryChan chan userModels.RetryCallbacks, gc *sharedconfig.GlobalConfig) {

	router.POST("/v1/callbacks/1l", func(c *gin.Context) {
		var callbackObj userModels.CallbackDeposit

		var err error

		data, _ := io.ReadAll(c.Request.Body)

		err = json.Unmarshal(data, &callbackObj)

		var invalidJSON tErrors.ErrorInvalidJSON

		if err != nil {
			log.Printf("[1L FAILED CALBACK] UNMARSHAL FAIL, error: [%v]\nBody[%v]", err, string(data))

			// log.Print(err)
			c.JSON(invalidJSON.HTTPCode(), invalidJSON.JSONError())
			return
		}

		if strings.EqualFold(callbackObj.Event, "WALLET_DEPOSIT_COMPLETED") {

			err = callbackObj.SaveDepositCallback(gc)
			if err != nil {
				log.Printf("[DEPOSIT CALLBACK FAILED RETRIEVAL] ERROR GETTING DEPOSIT ADDRESS OBJECT FROM DB for [%v, %v, %v], error: [%v]\n", callbackObj.Data.Currency, callbackObj.Data.Network, callbackObj.Data.ToAddress, err)

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

		} else {
			log.Printf("[SAVE CALLBACK]Uncompleted deposit: %+v\n", callbackObj)
		}

	})

}
