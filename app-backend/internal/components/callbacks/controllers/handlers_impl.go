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
	"trovo-wallet-api/internal/network"
	"trovo-wallet-api/internal/sharedconfig"

	evmkeypair "trovo-wallet-api/internal/evmkeypair"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"gorm.io/gorm/clause"
)

// postCallbacks1lHandler godoc
// @Summary POST /v1/callbacks/1l
// @Tags callbacks
// @Accept json
// @Produce json
// @Param body body userModels.CallbackDeposit true "Callback deposit payload"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /v1/callbacks/1l [post]
func postCallbacks1lHandler(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
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

	}
}

// postCallbacksDojaWebhookHandler godoc
// @Summary POST /v1/callbacks/doja/webhook
// @Tags callbacks
// @Accept json
// @Produce json
// @Param body body userModels.DojaKYCResponse true "Doja KYC webhook payload"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /v1/callbacks/doja/webhook [post]
func postCallbacksDojaWebhookHandler(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// log headers
		log.Printf("[KYC WEBHOOK ERROR] <><><><><><><><><>%+v\n<><><><><><><><><><><>\n", c.Request.Header)
		secret := gc.GetKycConfig("doja").SecretKey
		if secret == "" {
			c.JSON(http.StatusInternalServerError, "error")
			return
		}

		/*
			{
			  "metadata": {
			    "ipinfo": {
			      "status": "success",
			      "country": "Nigeria",
			      "city": "Lagos",
			      "district": "",
			      "zip": "",
			      "lat": 6.4474,
			      "lon": 3.3903,
			      "timezone": "Africa/Lagos",
			      "isp": "MTN NIGERIA Communication limited",
			      "org": "MTN Nigeria",
			      "as": "AS29465 MTN NIGERIA Communication limited",
			      "mobile": true,
			      "proxy": false,
			      "hosting": true,
			      "query": "102.88.115.133",
			      "region_name": "Lagos"
			    },
			    "device_info": "Mozilla/5.0 (Linux; Android 10; Infinix X692 Build/QP1A.190711.020; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/137.0.7151.115 Mobile Safari/537.36",
			    "user_id": "trovo"
			  },
			  "data": {
			    "index": {
			      "data": {

			      },
			      "message": "Successfully continued to the main checks.",
			      "status": true
			    },
			    "email": {
			      "data": {
			        "email": "info@trovotech.io"
			      },
			      "status": true,
			      "message": "info@trovotech.io validation Successful"
			    },
			    "user_data": {
			      "data": {
			        "first_name": "App",
			        "last_name": "Trovo",
			        "dob": "1977-08-25",
			        "email": "info@trovotech.io"
			      },
			      "message": "",
			      "status": true
			    },
			    "countries": {
			      "data": {
			        "country": "Nigeria"
			      },
			      "message": "Successfully continued to the next step.",
			      "status": true
			    },
			    "government_data": {
			      "data": {
			        "bvn": {
			          "entity": {
			            "customer": "6bb82c41-e15e-4308-b99d-e9640818eca9",
			            "app_id": null,
			            "bvn": "22324280081",
			            "first_name": "IFEANYI",
			            "last_name": "OKERE",
			            "middle_name": "",
			            "gender": "Male",
			            "date_of_birth": "01-Jun-1982",
			            "phone_number1": "08011111111",
			            "phone_number2": "",
			            "image_url": null,
			            "email": "",
			            "enrollment_bank": "",
			            "enrollment_branch": "",
			            "level_of_account": "",
			            "lga_of_origin": "",
			            "lga_of_residence": "",
			            "marital_status": "",
			            "name_on_card": "",
			            "nationality": "",
			            "nin": "",
			            "registration_date": "",
			            "residential_address": "",
			            "state_of_origin": "",
			            "state_of_residence": "",
			            "title": "",
			            "type": "BASIC",
			            "xc": null,
			            "sc": true,
			            "watch_listed": "",
			            "createdAt": "2024-11-23T18:38:32.000Z",
			            "updatedAt": "2024-11-23T18:38:32.000Z"
			          }
			        }
			      },
			      "message": "",
			      "status": true
			    },
			    "phone_number": {
			      "data": {
			        "phone": "2348034477900"
			      },
			      "message": "2348034477900 validation Successful",
			      "status": true
			    },
			    "address": {
			      "message": "Address Verification Failed",
			      "status": false,
			      "data": {
			        "location": {
			          "user_location": {
			            "latitude": "6.5191501",
			            "longitude": "3.2921159",
			            "name": "50 Ariyo Akinloye Street, Lagos, Nigeria"
			          },
			          "address_location": {
			            "latitude": "6.5211771",
			            "longitude": "3.2925486"
			          }
			        }
			      }
			    },
			    "additional_document": [
			      {
			        "document_type": "image",
			        "document_url": "https://images.dojah.io/image_686433e6da66af0047d5468887dd61b2-7a60-430b-a8f4-a0410fc8041b_1751397716.jpg"
			      }
			    ]
			  },
			  "id_type": "BVN",
			  "value": "22222222222",
			  "message": "Successfully completed the verification.",
			  "reference_id": "DJ-80F3FD18A1",
			  "widget_id": "67e6968cc5f45aec8ae35e0a",
			  "verification_mode": "OTP",
			  "verification_type": "BVN",
			  "verification_value": "22222222222",
			  "verification_url": "https://app.dojah.io/verifications/bio-data/30d31b0a-8ff3-4bf6-8fd2-455187e3e54e",
			  "selfie_url": "",
			  "status": true,
			  "aml": {
			    "status": false
			  },
			  "verification_status": "Completed"
			}
		*/
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
				//trigger CNGN onboarding
				/*			  "id_type": "BVN",
							  "value": "22222222222",
				*/
				if strings.EqualFold(event.IDType, "bvn") && len(event.Value) > 5 {
					// bvn is valid
					sriomsg, srierr := userServices.StablerailInitiateOnboardUser(&user, event.Value, gc)
					if srierr != nil {
						//perform operation trigger for stable rail for user
						gc.LogDiscordFailedRequest(fmt.Sprintf("FAILED to trigger stablerail onboarding for user %v with BVN %v", user.Username, event.Value))
						//TODO: save the detail for retry later,
						retryLater := userModels.StablerailOnboardUserRetry{
							TrovoUsername: user.Username,
							BVN:           event.Value,
						}
						e := dbTx.Omit(clause.Associations).Save(&retryLater).Error
						if e != nil {
							errMsg := fmt.Sprintf("[KYC WEBHOOK ERROR] Unable to save stablrail later retry task for [%v] due to [%v]\nRetry data: [%+v]", user.Username, e, retryLater)
							gc.LogDiscordFailedRequest(errMsg)
							log.Println(errMsg)

							c.JSON(http.StatusInternalServerError, "error")
							return
						}

					} else {
						//
						log.Printf("StablerailInitiate onboarding message for %v: %v\n", user.Username, sriomsg)
					}

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
	}
}

// postCallbacksFlutterwaveWebhookHandler godoc
// @Summary POST /v1/callbacks/flutterwave/webhook
// @Tags callbacks
// @Accept json
// @Produce json
// @Param body body userModels.FlutterwaveWebhook true "Flutterwave webhook payload"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /v1/callbacks/flutterwave/webhook [post]
func postCallbacksFlutterwaveWebhookHandler(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
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

			// orderBookTrov, err := gc.GetOrderBook(trovAsset.AssetCode, trovAsset.ContractAddress, cngnAsset.AssetCode, cngnAsset.ContractAddress)

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

			// orderBookGas, err := gc.GetOrderBook(os.Getenv("NATIVE_ASSET_CODE"), "", cngnAsset.AssetCode, cngnAsset.ContractAddress)

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
			trovToDispense := userServices.GetSwapEstimate(cngnAsset.AssetCode, cngnAsset.ContractAddress, trovAmount.String(), trovAsset.AssetCode, trovAsset.ContractAddress, gc)
			gasToDispense := userServices.GetSwapEstimate(cngnAsset.AssetCode, cngnAsset.ContractAddress, gasAmount.String(), os.Getenv("NATIVE_ASSET_CODE"), "", gc)
			// get faucet foir activation
			faucet, err := userServices.GetFaucetConfigByUserCase("ACTIVATION", gc)
			if err != nil {
				gc.LogDiscordFailedRequest(fmt.Sprintf("[FLUTTERWAVE WEBHOOK ERROR] unable to fetch FAUCET request for [%v]. Err: %v\n", "ACTIVATION", err))
				log.Printf("[FLUTTERWAVE WEBHOOK ERROR] unable to fetch FAUCET request for [%v]. Err: %v\n", "ACTIVATION", err)

				c.JSON(http.StatusInternalServerError, "Invalid Metadata:UserID")
				return
			}

			faucetKP := evmkeypair.MustParseFull(faucet.SecretKey)
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
				Destination:     user.Username,
				Memo:            "ACTIVATION",
				ContractAddress: "",
				AssetCode:       os.Getenv("NATIVE_ASSET_CODE"),
				Amount:          gasToDispense,
			}

			rpg, _, err := userServices.Pay(&signerUser, &sourceWallet, &payGas, gc)
			if err != nil {
				gc.LogDiscordFailedRequest(fmt.Sprintf("[FLUTTERWAVE WEBHOOK ERROR] error processing %v gas payment for [%v]. Err: %v\n", gasToDispense, "ACTIVATION", err))
				log.Printf("[FLUTTERWAVE WEBHOOK ERROR] error processing %v gas payment for [%v]. Err: %v\n", gasToDispense, "ACTIVATION", err)

				c.JSON(http.StatusInternalServerError, "Invalid Metadata:UserID")
				return
			}

			//sign payment
			if len(rpg.Transaction) > 0 {
				rpg.Commit = 0
				signedBase64, err := middleware.SignBase64Txn(faucetKP.Seed(), rpg.Transaction, gc.BantuNetworkPassphrase)
				if err != nil {
					gc.LogDiscordFailedRequest(fmt.Sprintf("[FLUTTERWAVE WEBHOOK ERROR] error signing gas payment for [%v]. Err: %v\n", "ACTIVATION", err))
					log.Printf("[FLUTTERWAVE WEBHOOK ERROR] error signing gas payment for [%v]. Err: %v\n", "ACTIVATION", err)

					c.JSON(http.StatusInternalServerError, "Invalid Metadata:UserID")
					return
				}

				rpg.TransactionSignature = signedBase64
			}

			_, _, err = userServices.Pay(&signerUser, &sourceWallet, rpg, gc)
			if err != nil {
				gc.LogDiscordFailedRequest(fmt.Sprintf("[FLUTTERWAVE WEBHOOK ERROR] error processing 2nd leg %v gas payment for [%v]. Err: %v\n", gasToDispense, "ACTIVATION", err))
				log.Printf("[FLUTTERWAVE WEBHOOK ERROR] error processing 2nd leg %v gas payment for [%v]. Err: %v\n", gasToDispense, "ACTIVATION", err)

				c.JSON(http.StatusInternalServerError, "Invalid Metadata:UserID")
				return
			}
			////////////////////END GAS

			////////////////////////START TROV
			payTrov := payments.PaymentInfo{
				Destination:     user.Username,
				Memo:            "ACTIVATION",
				ContractAddress: trovAsset.ContractAddress,
				AssetCode:       trovAsset.AssetCode,
				Amount:          trovToDispense,
			}

			rpt, _, err := userServices.Pay(&signerUser, &sourceWallet, &payTrov, gc)
			if err != nil {
				gc.LogDiscordFailedRequest(fmt.Sprintf("[FLUTTERWAVE WEBHOOK ERROR] error processing %v Trov payment for [%v]. Err: %v\n", trovToDispense, "ACTIVATION", err))
				log.Printf("[FLUTTERWAVE WEBHOOK ERROR] error processing %v trov payment for [%v]. Err: %v\n", trovToDispense, "ACTIVATION", err)

				c.JSON(http.StatusInternalServerError, "Invalid Metadata:UserID")
				return
			}

			//sign payment
			if len(rpt.Transaction) > 0 {
				rpt.Commit = 0
				signedBase64, err := middleware.SignBase64Txn(faucetKP.Seed(), rpt.Transaction, gc.BantuNetworkPassphrase)
				if err != nil {
					gc.LogDiscordFailedRequest(fmt.Sprintf("[FLUTTERWAVE WEBHOOK ERROR] error signing trov payment for [%v]. Err: %v\n", "ACTIVATION", err))
					log.Printf("[FLUTTERWAVE WEBHOOK ERROR] error signing trov payment for [%v]. Err: %v\n", "ACTIVATION", err)

					c.JSON(http.StatusInternalServerError, "Invalid Metadata:UserID")
					return
				}

				rpt.TransactionSignature = signedBase64
			}

			_, _, err = userServices.Pay(&signerUser, &sourceWallet, rpt, gc)
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
			//save payment invoices
			err = userServices.SaveUserPaymentInvoiceData(user.Username, "flutterwave", "ACTIVATION", event.Data.TxRef, "COMPLETED", &user.Username, &user.Address, nil, nil, nil, float64(event.Data.Amount), gc)
			if err != nil {
				gc.LogDiscordFailedRequest(fmt.Sprintf("[FLUTTERWAVE WEBHOOK ERROR] error saving payment invoice [%+v]. Err: %v\n", event, err))
				log.Printf("[FLUTTERWAVE WEBHOOK ERROR] error saving payment invoice [%+v]. Err: %v\n", event, err)
			}
			//send PN
			dataPayload := make(map[string]string)
			dataPayload["route"] = ""
			title := fmt.Sprintf("Account activation payment of %v%v now completed.", event.Data.Currency, event.Data.Amount)

			msg := fmt.Sprintf("Payment of %v%v for account activation has been confirmed. %v of Gas and %v%v has been dispensed to your wallet %v. Please check your pending asset to accept the TROV utility token.", event.Data.Currency, event.Data.Amount, gasToDispense, trovToDispense, "TROV", user.Username)
			user.SendPushMessage(title, msg, "", dataPayload, gc)
		}

		//Condition to process asset purchase
		if strings.EqualFold(event.MetaData.Product, "ASSET PURCHASE") && strings.EqualFold(event.Data.Status, "successful") {

			var invoice userModels.FiatPaymentInvoice
			invoiceErr := gc.DB.Where("id = ? AND status = ? AND payment_type = ?", event.Data.TxRef, "PENDING", "ASSET PURCHASE").First(&invoice).Error
			if invoiceErr != nil {
				// not found, or already processed by an earlier delivery of this webhook - nothing to do
				log.Printf("[FLUTTERWAVE WEBHOOK] no pending asset purchase invoice found for tx_ref [%v]: %v\n", event.Data.TxRef, invoiceErr)
				user.InvalidateUserCache(gc)
				c.JSON(http.StatusOK, "success")
				return
			}

			if invoice.Transaction == nil || invoice.TransactionSignature == nil {
				gc.LogDiscordFailedRequest(fmt.Sprintf("[FLUTTERWAVE WEBHOOK ERROR] asset purchase invoice [%v] has no signed transaction to submit\n", invoice.ID))
				log.Printf("[FLUTTERWAVE WEBHOOK ERROR] asset purchase invoice [%v] has no signed transaction to submit\n", invoice.ID)
				c.JSON(http.StatusOK, "success")
				return
			}

			var subscription userModels.TokenizedAssetSubscription
			subErr := gc.DB.Where("id = ?", invoice.ID).First(&subscription).Error
			if subErr != nil {
				gc.LogDiscordFailedRequest(fmt.Sprintf("[FLUTTERWAVE WEBHOOK ERROR] no tokenized asset subscription found for asset purchase invoice [%v]: %v\n", invoice.ID, subErr))
				log.Printf("[FLUTTERWAVE WEBHOOK ERROR] no tokenized asset subscription found for asset purchase invoice [%v]: %v\n", invoice.ID, subErr)
				c.JSON(http.StatusOK, "success")
				return
			}

			// resolve the actual invoice owner rather than assuming it matches the webhook's
			// MetaData.UserID-derived `user` - that field is caller-supplied and only used above for
			// the activation flow, so it isn't a reliable stand-in for who this invoice belongs to
			subscriber, subscriberErr := userModels.Username(subscription.SubscriberUsername).GetSimpleUser(gc.DB, gc)
			if subscriberErr != nil {
				gc.LogDiscordFailedRequest(fmt.Sprintf("[FLUTTERWAVE WEBHOOK ERROR] unable to load subscriber [%v] for asset purchase invoice [%v]: %v\n", subscription.SubscriberUsername, invoice.ID, subscriberErr))
				log.Printf("[FLUTTERWAVE WEBHOOK ERROR] unable to load subscriber [%v] for asset purchase invoice [%v]: %v\n", subscription.SubscriberUsername, invoice.ID, subscriberErr)
				c.JSON(http.StatusOK, "success")
				return
			}

			if len(subscription.TransactionID) == 0 {

				txnHash, submitErr := network.SubmitXdrWithSignature(gc.BantuExpansionClient, subscriber.PrimarySigner, *invoice.Transaction, *invoice.TransactionSignature)
				if submitErr != nil {
					gc.LogDiscordFailedRequest(fmt.Sprintf("[FLUTTERWAVE WEBHOOK ERROR] error submitting asset purchase invoice [%v] to blockchain: %v\n", invoice.ID, submitErr))
					log.Printf("[FLUTTERWAVE WEBHOOK ERROR] error submitting asset purchase invoice [%v] to blockchain: %v\n", invoice.ID, submitErr)
					// leave the invoice PENDING - it may still be resubmitted (a duplicated webhook
					// delivery, or a manual retry); recovery otherwise relies on the 2-day expiry sweep
					c.JSON(http.StatusOK, "success")
					return
				}

				subscription.TransactionID = txnHash
				gc.DB.Model(&userModels.TokenizedAssetSubscription{}).Where("id = ?", subscription.ID).Update("transaction_id", txnHash)

				saveErr := userServices.SaveUserPaymentData(subscriber.Username, invoice.ServiceProvider, invoice.PaymentType, txnHash, invoice.Amount, gc)
				if saveErr != nil {
					gc.LogDiscordFailedRequest(fmt.Sprintf("[FLUTTERWAVE WEBHOOK ERROR] error saving payment data for asset purchase invoice [%v]: %v\n", invoice.ID, saveErr))
					log.Printf("[FLUTTERWAVE WEBHOOK ERROR] error saving payment data for asset purchase invoice [%v]: %v\n", invoice.ID, saveErr)
				}

				gc.DB.Model(&userModels.FiatPaymentInvoice{}).Where("id = ?", invoice.ID).Update("status", "COMPLETED")

				if invoice.TransactionSource != nil {
					gc.ReleaseInUseChannelAccount(*invoice.TransactionSource)
				}

				subscriber.InvalidateUserCache(gc)

				walletAlias := ""
				if invoice.WalletAlias != nil {
					walletAlias = *invoice.WalletAlias
				}
				dataPayload := make(map[string]string)
				dataPayload["route"] = "assetSubscription"
				title := fmt.Sprintf("Asset purchase payment of %v%v now completed.", event.Data.Currency, event.Data.Amount)
				msg := fmt.Sprintf("Your payment of %v%v for %v has been confirmed and dispensed to your wallet %v.", event.Data.Currency, event.Data.Amount, subscription.AssetCode, walletAlias)
				subscriber.SendPushMessage(title, msg, "", dataPayload, gc)
			}
		}
		user.InvalidateUserCache(gc)
		c.JSON(http.StatusOK, "success")
	}
}
