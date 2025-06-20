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
	"os"
	"strings"
	payments "trovo-wallet-api/internal/components/payments/models"
	userModels "trovo-wallet-api/internal/components/users/models"
	userServices "trovo-wallet-api/internal/components/users/services"
	"trovo-wallet-api/internal/middleware"

	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/stellar/go/keypair"
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
		log.Printf("[KYC WEBHOOK ERROR] <><><><><><><><><>%+v\n<><><><><><><><><><><>\n", c.Request.Header)
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
		dojahIP := c.ClientIP()
		if len(c.GetHeader("Cf-Connecting-Ip")) > 4 {
			dojahIP = c.GetHeader("Cf-Connecting-Ip")
		}
		if dojahIP == "20.112.64.208" {

			// if err := json.Unmarshal([]byte(body), &event); err != nil {
			// 	log.Println("[KYC WEBHOOK ERROR] Invalid JSON")

			// 	c.JSON(http.StatusBadRequest, "Invalid JSON")
			// 	return
			// }

			// Do something with event
			log.Println("[KYC WEBHOOK] ✅ ✅ ✅ ✅ ✅ ✅ ✅ ✅ ✅ ✅ ✅ ✅ ✅ ✅ Valid webhook received:", event)
		} else {
			// log.Printf("[KYC WEBHOOK ERROR] Invalid signature. x-dojah-signature: [%v], Expected Mac: [%v]\n", signature, expectedMAC)
			log.Printf("[KYC WEBHOOK ERROR] Invalid IP. x-dojah-signature: [%v], Expected Mac: [%v], IP: [%v]\n", signature, expectedMAC, dojahIP)

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
		takeAction := false
		//check if it is pending.
		if strings.EqualFold(event.VerificationStatus, "Pending") {
			//mark the user progress as pending and send push notification
			if widget.Level == 1 {

				if progress.KYCLevel1Completed == 0 {
					takeAction = true
					progress.KYCLevel1Submitted = 1
				}

			} else if widget.Level == 2 {
				if progress.KYCLevel2Completed == 0 {
					takeAction = true
					progress.KYCLevel2Submitted = 1
				}
			} else if widget.Level == 3 {
				if progress.KYCLevel3Completed == 0 {
					takeAction = true
					progress.KYCLevel3Submitted = 1
				}
			} else if widget.Level == 4 {
				if progress.KYCLevel1Completed == 0 {
					takeAction = true
					progress.KYCLevel3Submitted = 1
				}
			}

			//send PN
			dataPayload := make(map[string]string)
			dataPayload["route"] = ""
			title := fmt.Sprintf("KYC Level %v now pending confirmation", widget.Level)
			msg := fmt.Sprintf("KYC Level %v has been submitted and is now pending confirmation. Please wait for it to finish before continuing to other levels.", widget.Level)
			if takeAction {

				user.SendPushMessage(title, msg, "", dataPayload, gc)
			}
		}

		//check if it is completed.
		if strings.EqualFold(event.VerificationStatus, "Completed") {
			//mark the user progress as Completed and send push notification
			if widget.Level == 1 {
				if progress.KYCLevel1Completed == 0 {
					takeAction = true
					progress.KYCLevel1Submitted = 1
					progress.KYCLevel1Completed = 1
				}

			} else if widget.Level == 2 {
				if progress.KYCLevel2Completed == 0 {
					takeAction = true
					progress.KYCLevel2Submitted = 1
					progress.KYCLevel2Completed = 1
				}

			} else if widget.Level == 3 {
				if progress.KYCLevel3Completed == 0 {
					takeAction = true
					progress.KYCLevel3Submitted = 1
					progress.KYCLevel3Completed = 1
				}

			} else if widget.Level == 4 {
				if progress.KYCLevel4Completed == 0 {
					takeAction = true
					progress.KYCLevel3Submitted = 1
					progress.KYCLevel4Completed = 1
				}

			}
			if takeAction {
				user.KYCVerified = widget.Level

				e := dbTx.Omit(clause.Associations).Save(&user).Error
				if e != nil {
					errMsg := fmt.Sprintf("[KYC WEBHOOK ERROR] Unable to save progress for [%v] due to [%v]", event.Metadata.UserID, e)
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
		}

		//check if it is failed.
		if strings.EqualFold(event.VerificationStatus, "Failed") {
			//reset user progress and send push notification
			if widget.Level == 1 {
				if progress.KYCLevel1Completed == 0 {
					takeAction = true
					progress.KYCLevel1Submitted = 0
					progress.KYCLevel1Completed = 0
				}

			} else if widget.Level == 2 {
				if progress.KYCLevel2Completed == 0 {
					takeAction = true
					progress.KYCLevel2Submitted = 0
					progress.KYCLevel2Completed = 0
				}
			} else if widget.Level == 3 {
				if progress.KYCLevel3Completed == 0 {
					takeAction = true
					progress.KYCLevel3Submitted = 0
					progress.KYCLevel3Completed = 0
				}
			} else if widget.Level == 4 {
				if progress.KYCLevel4Completed == 0 {
					takeAction = true
					progress.KYCLevel3Submitted = 0
					progress.KYCLevel4Completed = 0
				}
			}
			if takeAction {
				//send PN
				dataPayload := make(map[string]string)
				dataPayload["route"] = ""
				title := fmt.Sprintf("KYC Level %v failed.", widget.Level)

				msg := fmt.Sprintf("KYC Level %v failed verification. Please resubmit your correct information to try again.", widget.Level)
				user.SendPushMessage(title, msg, "", dataPayload, gc)
			}

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

	router.POST("/v1/callbacks/flutterwave/webhook", func(c *gin.Context) {
		// log headers
		log.Printf("[FLUTTERWAVE WEBHOOK ERROR] <><><><><><><><><>%+v\n<><><><><><><><><><><>\n", c.Request.Header)
		pcc, err := userServices.GetPaymentConfigByServiceProvider("flutterwave", gc)
		if err != nil {
			c.JSON(http.StatusInternalServerError, "error")
			return
		}

		// Read body
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			log.Println("[FLUTTERWAVE WEBHOOK ERROR] Unable to read request body")
			c.JSON(http.StatusBadRequest, "Unable to read request body")
			return
		}
		// save record
		userServices.SavePaymentWebhookData("flutterwave", string(body), gc)
		// Compute HMAC

		// Get signature from headers
		hash := c.GetHeader("verif-hash")
		// var event map[string]interface{}
		var event userModels.FlutterwaveWebhook

		if hash == pcc.VerificationHash {

			// Do something with event
			log.Println("[FLUTTERWAVE WEBHOOK] ✅ ✅ ✅ ✅ ✅ ✅ ✅ ✅ ✅ ✅ ✅ ✅ ✅ ✅ Valid webhook received:", pcc)
		} else {
			// log.Printf("[KYC WEBHOOK ERROR] Invalid signature. x-dojah-signature: [%v], Expected Mac: [%v]\n", signature, expectedMAC)
			log.Printf("[FLUTTERWAVE WEBHOOK ERROR] hash [%v]\n", hash)

			c.JSON(http.StatusUnauthorized, "Invalid signature")
			return
		}
		if err := json.Unmarshal([]byte(body), &event); err != nil {
			log.Println("[FLUTTERWAVE WEBHOOK ERROR] Invalid JSON")

			c.JSON(http.StatusBadRequest, "Invalid JSON")
			return
		}

		// check if the user exists
		user, userErr := userModels.Username(event.MetaData.UserID).GetSimpleUser(gc.DB, gc)

		if userErr != nil {
			gc.LogDiscordFailedRequest(fmt.Sprintf("[FLUTTERWAVE WEBHOOK ERROR] User ID [%v] in metadata is invalid\n", event.MetaData.UserID))
			log.Printf("[FLUTTERWAVE WEBHOOK ERROR] User ID [%v] in metadata is invalid\n", event.MetaData.UserID)

			c.JSON(http.StatusInternalServerError, "Invalid Metadata:UserID")
			return
		}
		//fetch data from flutterwave using the reference id.
		{
			//fetch from flutterwave
		}
		if strings.EqualFold(event.MetaData.Product, "activation") && strings.EqualFold(event.Data.Status, "successful") {

			//process value

			//get percentage for Gas
			cc := userModels.CountryCode(*user.CountryCode).GetConfig(gc)
			if len(cc.CountryCode) == 0 {
				//use default
				cc = userModels.CountryCode("NG").GetConfig(gc)

			}
			trovAmount := decimal.NewFromFloat(cc.TrovTokenActivationPercent / 100).Mul(decimal.NewFromInt(int64(event.Data.Amount)))
			gasAmount := decimal.NewFromInt(int64(event.Data.Amount)).Sub(trovAmount)

			trovAsset := gc.GetCuratedAssetByCode("TROV")
			cngnAsset := gc.GetCuratedAssetByCode("CNGN")

			// orderBookTrov, err := gc.GetOrderBook(trovAsset.AssetCode, trovAsset.AssetIssuer, cngnAsset.AssetCode, cngnAsset.AssetIssuer)

			// if err != nil {
			// 	gc.LogDiscordFailedRequest(fmt.Sprintf("[FLUTTERWAVE WEBHOOK ERROR] unable to fetch orderbook request for [%v]/[%v]. Err: %v\n", trovAsset.AssetCode, cngnAsset.AssetCode, err))
			// 	log.Printf("[FLUTTERWAVE WEBHOOK ERROR] unable to fetch orderbook request for [%v]/[%v]. Err: %v\n", trovAsset.AssetCode, cngnAsset.AssetCode, err)

			// 	c.JSON(http.StatusInternalServerError, "Invalid Metadata:UserID")
			// 	return
			// }
			// if len(orderBookTrov.Asks) == 0 {
			// 	gc.LogDiscordFailedRequest(fmt.Sprintf("[FLUTTERWAVE WEBHOOK ERROR] no market price for [%v]/[%v]. Err: %v\n", trovAsset.AssetCode, cngnAsset.AssetCode, err))
			// 	log.Printf("[FLUTTERWAVE WEBHOOK ERROR] no market price for [%v]/[%v]. Err: %v\n", trovAsset.AssetCode, cngnAsset.AssetCode, err)

			// 	c.JSON(http.StatusInternalServerError, "Invalid Metadata:UserID")
			// 	return
			// }
			//get the price of trov
			// priceOfTrov := orderBookTrov.Asks[0].Price

			// orderBookGas, err := gc.GetOrderBook(os.Getenv("NATIVE_ASSET_CODE"), "", cngnAsset.AssetCode, cngnAsset.AssetIssuer)

			// if err != nil {
			// 	gc.LogDiscordFailedRequest(fmt.Sprintf("[FLUTTERWAVE WEBHOOK ERROR] unable to fetch orderbook request for [%v]/[%v]. Err: %v\n", os.Getenv("NATIVE_ASSET_CODE"), cngnAsset.AssetCode, err))
			// 	log.Printf("[FLUTTERWAVE WEBHOOK ERROR] unable to fetch orderbook request for [%v]/[%v]. Err: %v\n", os.Getenv("NATIVE_ASSET_CODE"), cngnAsset.AssetCode, err)

			// 	c.JSON(http.StatusInternalServerError, "Invalid Metadata:UserID")
			// 	return
			// }
			// if len(orderBookGas.Asks) == 0 {
			// 	gc.LogDiscordFailedRequest(fmt.Sprintf("[FLUTTERWAVE WEBHOOK ERROR] no market price for [%v]/[%v]. Err: %v\n", os.Getenv("NATIVE_ASSET_CODE"), cngnAsset.AssetCode, err))
			// 	log.Printf("[FLUTTERWAVE WEBHOOK ERROR] no market price for [%v]/[%v]. Err: %v\n", os.Getenv("NATIVE_ASSET_CODE"), cngnAsset.AssetCode, err)

			// 	c.JSON(http.StatusInternalServerError, "Invalid Metadata:UserID")
			// 	return
			// }
			// priceOfGas := orderBookGas.Asks[0].Price
			// trovToDispense := trovAmount.Div(decimal.RequireFromString(priceOfTrov)).Truncate(7)
			// gasToDispense := gasAmount.Div(decimal.RequireFromString(priceOfGas)).Truncate(7)
			trovToDispense := userServices.GetSwapEstimate(cngnAsset.AssetCode, cngnAsset.AssetIssuer, trovAmount.String(), trovAsset.AssetCode, trovAsset.AssetIssuer, gc)
			gasToDispense := userServices.GetSwapEstimate(cngnAsset.AssetCode, cngnAsset.AssetIssuer, gasAmount.String(), os.Getenv("NATIVE_ASSET_CODE"), "", gc)
			// get faucet foir activation
			faucet, err := userServices.GetFaucetConfigByUserCase("ACTIVATION", gc)
			if err != nil {
				gc.LogDiscordFailedRequest(fmt.Sprintf("[FLUTTERWAVE WEBHOOK ERROR] unable to fetch FAUCET request for [%v]. Err: %v\n", "ACTIVATION", err))
				log.Printf("[FLUTTERWAVE WEBHOOK ERROR] unable to fetch FAUCET request for [%v]. Err: %v\n", "ACTIVATION", err)

				c.JSON(http.StatusInternalServerError, "Invalid Metadata:UserID")
				return
			}

			faucetKP := keypair.MustParseFull(faucet.SecretKey)
			sourceWallet, err := userModels.UserWalletID(faucetKP.Address()).GetWallet(gc.DB, gc)
			if err != nil {
				gc.LogDiscordFailedRequest(fmt.Sprintf("[FLUTTERWAVE WEBHOOK ERROR] unable to fetch FAUCET wallet for [%v]. Err: %v\n", "ACTIVATION", err))
				log.Printf("[FLUTTERWAVE WEBHOOK ERROR] unable to fetch FAUCET wallet for [%v]. Err: %v\n", "ACTIVATION", err)

				c.JSON(http.StatusInternalServerError, "Invalid Metadata:UserID")
				return
			}
			signerUser, err := userModels.UserWalletID(faucetKP.Address()).GetWalletOwner(gc.DB, gc)
			if err != nil {
				gc.LogDiscordFailedRequest(fmt.Sprintf("[FLUTTERWAVE WEBHOOK ERROR] unable to fetch FAUCET wallet owner for [%v]. Err: %v\n", "ACTIVATION", err))
				log.Printf("[FLUTTERWAVE WEBHOOK ERROR] unable to fetch FAUCET wallet owner for [%v]. Err: %v\n", "ACTIVATION", err)

				c.JSON(http.StatusInternalServerError, "Invalid Metadata:UserID")
				return
			}

			////////////////////////START GAS
			payGas := payments.PaymentInfo{
				Destination: user.Username,
				Memo:        "ACTIVATION",
				AssetIssuer: "",
				AssetCode:   os.Getenv("NATIVE_ASSET_CODE"),
				Amount:      gasToDispense,
			}

			_, _, err = userServices.Pay(&signerUser, &sourceWallet, &payGas, gc)
			if err != nil {
				gc.LogDiscordFailedRequest(fmt.Sprintf("[FLUTTERWAVE WEBHOOK ERROR] error processing %v gas payment for [%v]. Err: %v\n", gasToDispense, "ACTIVATION", err))
				log.Printf("[FLUTTERWAVE WEBHOOK ERROR] error processing %v gas payment for [%v]. Err: %v\n", gasToDispense, "ACTIVATION", err)

				c.JSON(http.StatusInternalServerError, "Invalid Metadata:UserID")
				return
			}

			//sign payment
			if len(payGas.Transaction) > 0 {
				payGas.Commit = 0
				signedBase64, err := middleware.SignBase64Txn(faucetKP.Seed(), payGas.Transaction, payGas.NetworkPassPhrase)
				if err != nil {
					gc.LogDiscordFailedRequest(fmt.Sprintf("[FLUTTERWAVE WEBHOOK ERROR] error signing gas payment for [%v]. Err: %v\n", "ACTIVATION", err))
					log.Printf("[FLUTTERWAVE WEBHOOK ERROR] error signing gas payment for [%v]. Err: %v\n", "ACTIVATION", err)

					c.JSON(http.StatusInternalServerError, "Invalid Metadata:UserID")
					return
				}

				payGas.TransactionSignature = signedBase64
			}

			_, _, err = userServices.Pay(&signerUser, &sourceWallet, &payGas, gc)
			if err != nil {
				gc.LogDiscordFailedRequest(fmt.Sprintf("[FLUTTERWAVE WEBHOOK ERROR] error processing 2nd leg %v gas payment for [%v]. Err: %v\n", gasToDispense, "ACTIVATION", err))
				log.Printf("[FLUTTERWAVE WEBHOOK ERROR] error processing 2nd leg %v gas payment for [%v]. Err: %v\n", gasToDispense, "ACTIVATION", err)

				c.JSON(http.StatusInternalServerError, "Invalid Metadata:UserID")
				return
			}
			////////////////////END GAS

			////////////////////////START TROV
			payTrov := payments.PaymentInfo{
				Destination: user.Username,
				Memo:        "ACTIVATION",
				AssetIssuer: trovAsset.AssetIssuer,
				AssetCode:   trovAsset.AssetCode,
				Amount:      trovToDispense,
			}

			_, _, err = userServices.Pay(&signerUser, &sourceWallet, &payTrov, gc)
			if err != nil {
				gc.LogDiscordFailedRequest(fmt.Sprintf("[FLUTTERWAVE WEBHOOK ERROR] error processing %v Trov payment for [%v]. Err: %v\n", trovToDispense, "ACTIVATION", err))
				log.Printf("[FLUTTERWAVE WEBHOOK ERROR] error processing %v trov payment for [%v]. Err: %v\n", trovToDispense, "ACTIVATION", err)

				c.JSON(http.StatusInternalServerError, "Invalid Metadata:UserID")
				return
			}

			//sign payment
			if len(payTrov.Transaction) > 0 {
				payTrov.Commit = 0
				signedBase64, err := middleware.SignBase64Txn(faucetKP.Seed(), payTrov.Transaction, payTrov.NetworkPassPhrase)
				if err != nil {
					gc.LogDiscordFailedRequest(fmt.Sprintf("[FLUTTERWAVE WEBHOOK ERROR] error signing trov payment for [%v]. Err: %v\n", "ACTIVATION", err))
					log.Printf("[FLUTTERWAVE WEBHOOK ERROR] error signing trov payment for [%v]. Err: %v\n", "ACTIVATION", err)

					c.JSON(http.StatusInternalServerError, "Invalid Metadata:UserID")
					return
				}

				payTrov.TransactionSignature = signedBase64
			}

			_, _, err = userServices.Pay(&signerUser, &sourceWallet, &payTrov, gc)
			if err != nil {
				gc.LogDiscordFailedRequest(fmt.Sprintf("[FLUTTERWAVE WEBHOOK ERROR] error processing 2nd leg %v trov payment for [%v]. Err: %v\n", trovToDispense, "ACTIVATION", err))
				log.Printf("[FLUTTERWAVE WEBHOOK ERROR] error processing 2nd leg %v trov payment for [%v]. Err: %v\n", trovToDispense, "ACTIVATION", err)

				c.JSON(http.StatusInternalServerError, "Invalid Metadata:UserID")
				return
			}
			////////////////////END TROV

			//save payment data
			err = userServices.SaveUserPaymentData(user.Username, "flutterwave", "ACTIVATION", event.Data.TxRef, float64(event.Data.Amount), gc)
			if err != nil {
				gc.LogDiscordFailedRequest(fmt.Sprintf("[FLUTTERWAVE WEBHOOK ERROR] error saving payment data [%+v]. Err: %v\n", event, err))
				log.Printf("[FLUTTERWAVE WEBHOOK ERROR] error saving payment data [%+v]. Err: %v\n", event, err)
			}
			//send PN
			dataPayload := make(map[string]string)
			dataPayload["route"] = ""
			title := fmt.Sprintf("Account activation payment of %v%v now completed.", event.Data.Currency, event.Data.Amount)

			msg := fmt.Sprintf("Payment of %v%v for account activation has been confirmed. %v of Gas and %v%v has been dispensed to your wallet %v. Please check your pending asset to accept the TROV utility token.", event.Data.Currency, event.Data.Amount, gasToDispense, trovToDispense, "TROV", user.Username)
			user.SendPushMessage(title, msg, "", dataPayload, gc)
		}
		user.InvalidateUserCache(gc)
		c.JSON(http.StatusOK, "success")
	})

}
