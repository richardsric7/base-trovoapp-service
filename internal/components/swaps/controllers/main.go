package payments

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	swaperrors "trovo-wallet-api/internal/components/swaps/errors"
	swapModels "trovo-wallet-api/internal/components/swaps/models"
	swapServices "trovo-wallet-api/internal/components/swaps/services"
	usersdb "trovo-wallet-api/internal/components/users/db"
	userModels "trovo-wallet-api/internal/components/users/models"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/middleware"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/gin-gonic/gin"
)

//Init initializes the controller
func Init(router *gin.Engine, gc *sharedconfig.GlobalConfig) {
	router.POST("/v1/users/swap", middleware.AuthenticationMiddlewareUsingTimestamp(), func(c *gin.Context) {

		// suppliedUsername := strings.TrimSpace(strings.ToLower(c.Param("primaryWalletAlias")))
		// uDec, e := base64.URLEncoding.DecodeString(c.Param("primaryWalletAlias"))
		// if e == nil {
		// 	//check if the decoded contains any non-english character
		// 	invalidChars := 0

		// 	acceptedChars := "abcdefghijklmnopqrstuvwxyz_1234567890/"
		// 	for _, c := range uDec {

		// 		if !strings.Contains(acceptedChars, strings.TrimSpace(strings.ToLower(string(c)))) {
		// 			invalidChars++
		// 		}

		// 	}
		// 	if invalidChars == 0 {
		// 		suppliedUsername = string(uDec)
		// 	}

		// }

		// if suppliedUsername == "null" {
		// 	log.Printf("user cannot be %v\n", suppliedUsername)
		// 	c.JSON(http.StatusBadRequest, gin.H{"error": "user cannot be null"})
		// 	return
		// }
		//check if username is reserved. Reserved usernames should not send payments.

		signerOwner, getUserError := usersdb.GetUser(middleware.ExtractSigner(c), gc.DB)

		if getUserError != nil {

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
		// _, checkReservedUserError := usersdb.UsernameIsReserved(signerOwner.Username, gc.DB)
		// if checkReservedUserError != nil {

		// 	var ex tErrors.GenericError
		// 	var ok bool

		// 	ex, ok = checkReservedUserError.(tErrors.GenericError)
		// 	if ok {
		// 		c.JSON(ex.HTTPCode(), ex.JSONError())
		// 	} else {
		// 		c.JSON(http.StatusBadRequest, gin.H{"error": checkReservedUserError.Error()})
		// 	}
		// 	return
		// }

		wallet, temp, getWalletError := usersdb.GetWallet(middleware.ExtractPublicKey(c), gc.DB)

		if getWalletError != nil {

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
		if temp {
			c.JSON(http.StatusBadRequest, (&tErrors.CustomError{Param: "publicKey", Err: "error-temporary-account-forbidden", ErrMessage: "temporary accounts are forbidden from making payment or swap requests", Code: http.StatusForbidden}).JSONError())
			return
		}

		primarySigner := signerOwner.PrimarySigner

		walletOwner, errWalletOwner := userModels.UserWalletID(wallet.ID).GetWalletOwner(gc.DB)

		if errWalletOwner != nil {

			var ex tErrors.GenericError
			var ok bool

			ex, ok = errWalletOwner.(tErrors.GenericError)
			if ok {
				c.JSON(ex.HTTPCode(), ex.JSONError())
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": errWalletOwner.Error()})
			}
			return
		}

		if primarySigner != middleware.ExtractSigner(c) {
			c.JSON(http.StatusBadRequest, (&swaperrors.ErrorInvalidPaymentSender{}).JSONError())
			return
		}
		//TODO: check if the wallet belongs to the person making swap or if the person has permission to do swap.

		var swapInfo swapModels.SwapSendInfo
		var err error

		data, _ := ioutil.ReadAll(c.Request.Body)

		err = json.Unmarshal(data, &swapInfo)

		var invalidJSON tErrors.ErrorInvalidJSON

		if err != nil {
			log.Print(err)
			c.JSON(invalidJSON.HTTPCode(), invalidJSON.JSONError())
			return
		}

		swapError := swapServices.SwapSend(&walletOwner, &wallet, middleware.ExtractSigner(c), &swapInfo, gc)

		if swapError != nil {
			var ex tErrors.GenericError
			var ok bool

			ex, ok = swapError.(tErrors.GenericError)
			if ok {
				c.JSON(ex.HTTPCode(), ex.JSONError())
			} else {
				log.Print(swapError)
				c.JSON(http.StatusBadRequest, gin.H{"error": swapError.Error()})
			}

			return
		}

		if len(swapInfo.TransactionID) == 0 {
			c.JSON(http.StatusAccepted, swapInfo)

		} else {
			dataPayload := make(map[string]string)
			dataPayload["route"] = "basicTransactionHistory"
			var fromAsset, toAsset string
			if swapInfo.SourceAssetIssuer == "" || swapInfo.SourceAssetIssuer == "native" {
				fromAsset = "XBN"
			} else {
				fromAsset = swapInfo.SourceAssetCode
			}
			if swapInfo.DestinationAssetIssuer == "" || swapInfo.DestinationAssetIssuer == "native" {
				toAsset = "XBN"
			} else {
				toAsset = swapInfo.DestinationAssetCode
			}
			if swapInfo.TransactionID != "PENDING_AUTH" {

				signerOwner.SendPushMessage(fmt.Sprintf("Trovo: %v Swapped to %v on %v!", fromAsset, toAsset, wallet.Alias), fmt.Sprintf("You have successfully swapped %v%v to %v%v from your wallet with alias %v", swapInfo.SourceAmount, fromAsset, swapInfo.SwappedEstimate, toAsset, wallet.Alias), "", dataPayload, gc)

			} else {
				//send push notification to all parties that need to authorize
				signerOwner.SendPushMessage(fmt.Sprintf("Trovo: Request to swap %v to %v on %v submitted!", fromAsset, toAsset, wallet.Alias), fmt.Sprintf("You have submitted request to swap %v%v to %v from the wallet with alias %v", swapInfo.SourceAmount, fromAsset, toAsset, wallet.Alias), "", dataPayload, gc)

			}

			//log user current location
			signerOwner.PublicIP = c.ClientIP()
			swapServices.UpdateAndLogUserSwapGeoInformation(&signerOwner, &swapInfo, gc.DB)
			senderPaymentHistoryCacheKey := fmt.Sprintf("[GET] /v1/users/%v/payments", signerOwner.Username)
			senderCacheKey := fmt.Sprintf("[GET] /v1/users/%v", signerOwner.Username)

			gc.RedisCache.InvalidateCachedHttpResponse(senderCacheKey, senderPaymentHistoryCacheKey)
			c.JSON(http.StatusOK, swapInfo)

		}

	})

}
