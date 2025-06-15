package callbacks

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	userModels "trovo-wallet-api/internal/components/users/models"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
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
		// log headers
		log.Printf("<><><><><><><><><>%+v\n<><><><><><><><><><><>\n", c.Request.Header)
		secret := gc.GetKycConfig("doja").SecretKey
		if secret == "" {
			c.JSON(http.StatusInternalServerError, "error")
			return
		}

		// Read body
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			log.Println("[KYC WEBHOOK ERROR] Unable to read request body")
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
		// var event map[string]interface{}
		var event userModels.DojaKYCResponse
		if hmac.Equal([]byte(expectedMAC), []byte(signature)) {

			if err := json.Unmarshal([]byte(body), &event); err != nil {
				log.Println("[KYC WEBHOOK ERROR] Invalid JSON")

				c.JSON(http.StatusBadRequest, "Invalid JSON")
				return
			}

			// Do something with event
			log.Println("[KYC WEBHOOK] Valid webhook received:", event)
		} else {
			log.Printf("[KYC WEBHOOK ERROR] Invalid signature. x-dojah-signature: [%v], Expected Mac: [%v]\n", signature, expectedMAC)

			// c.JSON(http.StatusUnauthorized, "Invalid signature")
			// return
		}
		if err := json.Unmarshal([]byte(body), &event); err != nil {
			log.Println("[KYC WEBHOOK ERROR] Invalid JSON")

			c.JSON(http.StatusBadRequest, "Invalid JSON")
			return
		}
		// get the KYC widget info
		widget := userModels.DojaWidgetID(event.WidgetID).GetByID(gc.DB)
		//check if it is pending.
		if len(widget.ID) == 0 {
			log.Println("[KYC WEBHOOK ERROR] Invalid WidgetID", event.WidgetID)

			c.JSON(http.StatusBadRequest, "Invalid WidgetID")
			return
		}

		// check if the user exists
		user, userErr := userModels.Username(event.Metadata.UserID).GetSimpleUser(gc.DB, gc)

		if userErr != nil {
			gc.LogDiscordFailedRequest(fmt.Sprintf("[KYC WEBHOOK ERROR] User ID [%v] in metadata is invalid\n", event.Metadata.UserID))
			log.Printf("[KYC WEBHOOK ERROR] User ID [%v] in metadata is invalid\n", event.Metadata.UserID)

			c.JSON(http.StatusOK, "Invalid Metadata:UserID")
			return
		}

		// get user KYC progress
		progress := userModels.Username(event.Metadata.UserID).GetUserDojaKYCProgress(gc.DB)
		dbTx := gc.DB.Begin()
		defer dbTx.Rollback()
		//check if it is pending.
		if strings.EqualFold(event.VerificationStatus, "Pending") {
			//mark the user progress as pending and send push notification
			if widget.Level == 1 {
				progress.KYCLevel1Submitted = 1
			} else if widget.Level == 2 {
				progress.KYCLevel2Submitted = 1
			} else if widget.Level == 3 {
				progress.KYCLevel3Submitted = 1
			} else if widget.Level == 4 {
				progress.KYCLevel3Submitted = 1
			}

			//send PN
			dataPayload := make(map[string]string)
			dataPayload["route"] = ""
			title := fmt.Sprintf("KYC Level %v now pending confirmation", widget.Level)
			msg := fmt.Sprintf("KYC Level %v has been submitted and is now pending confirmation. Please wait for it to finish before continuing to other levels.", widget.Level)
			user.SendPushMessage(title, msg, "", dataPayload, gc)
		}

		//check if it is completed.
		if strings.EqualFold(event.VerificationStatus, "Completed") {
			//mark the user progress as Completed and send push notification
			if widget.Level == 1 {
				progress.KYCLevel1Submitted = 1
				progress.KYCLevel1Completed = 1
			} else if widget.Level == 2 {
				progress.KYCLevel2Submitted = 1
				progress.KYCLevel2Completed = 1
			} else if widget.Level == 3 {
				progress.KYCLevel3Submitted = 1
				progress.KYCLevel3Completed = 1
			} else if widget.Level == 4 {
				progress.KYCLevel3Submitted = 1
				progress.KYCLevel4Completed = 1
			}
			user.KYCVerified = widget.Level

			e := dbTx.Omit(clause.Associations).Save(&user).Error
			if e != nil {
				errMsg := fmt.Sprintf("[KYC WEBHOOK ERROR] Unable to save progressfor [%v] due to [%v]", event.Metadata.UserID, e)
				gc.LogDiscordFailedRequest(errMsg)
				log.Println(errMsg)

				c.JSON(http.StatusInternalServerError, "error")
				return
			}
			//send PN
			dataPayload := make(map[string]string)
			dataPayload["route"] = ""
			title := fmt.Sprintf("KYC Level %v now completed.", widget.Level)
			nextLevel := widget.Level + 1
			if nextLevel > 4 {
				nextLevel = 0
			}
			msg := fmt.Sprintf("KYC Level %v has been completed.%v", widget.Level, func() string {
				if nextLevel == 0 {
					return ""
				}
				return fmt.Sprintf(" Please proceed to next level (%v) when ready.", nextLevel)
			}())
			user.SendPushMessage(title, msg, "", dataPayload, gc)
		}

		//check if it is failed.
		if strings.EqualFold(event.VerificationStatus, "Failed") {
			//reset user progress and send push notification
			if widget.Level == 1 {
				progress.KYCLevel1Submitted = 0
				progress.KYCLevel1Completed = 0
			} else if widget.Level == 2 {
				progress.KYCLevel2Submitted = 0
				progress.KYCLevel2Completed = 0
			} else if widget.Level == 3 {
				progress.KYCLevel3Submitted = 0
				progress.KYCLevel3Completed = 0
			} else if widget.Level == 4 {
				progress.KYCLevel3Submitted = 0
				progress.KYCLevel4Completed = 0
			}

			//send PN
			dataPayload := make(map[string]string)
			dataPayload["route"] = ""
			title := fmt.Sprintf("KYC Level %v failed.", widget.Level)

			msg := fmt.Sprintf("KYC Level %v failed verification. Please resubmit your correct information to try again.", widget.Level)
			user.SendPushMessage(title, msg, "", dataPayload, gc)
		}

		e := dbTx.Save(&progress).Error
		if e != nil {
			errMsg := fmt.Sprintf("[KYC WEBHOOK ERROR] Unable to save progress for [%v] due to [%v]", event.Metadata.UserID, e)
			gc.LogDiscordFailedRequest(errMsg)
			log.Println(errMsg)

			c.JSON(http.StatusInternalServerError, "error")
			return
		}

		dbTx.Commit()

		c.JSON(http.StatusOK, "success")
	})

}
