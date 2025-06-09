package callbacks

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
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
			log.Printf("[1L FAILED CALBACK] UNMARSHAL FAIL, error: [%v]\nPayload:[%v]\n", err, string(data))

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
			c.JSON(http.StatusOK, "success")
		} else {
			log.Printf("[SAVE CALLBACK]other callback info: %+v\n", callbackObj)
			c.JSON(http.StatusOK, "success")
		}

	})

	router.POST("/v1/callbacks/doja/webhook", func(c *gin.Context) {

		secret := gc.GetKycConfig("doja").SecretKey
		if secret == "" {
			c.JSON(http.StatusInternalServerError, "error")
			return
		}

		// Read body
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, "Unable to read request body")
			return
		}
		// save record
		gc.SaveKycWebhookData("doja", string(body))
		// Compute HMAC
		mac := hmac.New(sha256.New, []byte(secret))
		mac.Write(body)
		expectedMAC := hex.EncodeToString(mac.Sum(nil))

		// Get signature from headers
		signature := c.GetHeader("x-dojah-signature")

		if hmac.Equal([]byte(expectedMAC), []byte(signature)) {
			var event map[string]interface{}
			if err := json.Unmarshal(body, &event); err != nil {
				c.JSON(http.StatusBadRequest, "Invalid JSON")
				return
			}

			// Do something with event
			log.Println("Valid webhook received:", event)
		} else {
			c.JSON(http.StatusUnauthorized, "Invalid signature")
			return
		}

		c.JSON(http.StatusOK, "success")
	})

}
