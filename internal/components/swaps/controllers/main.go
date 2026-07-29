package payments

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	swaperrors "trovo-wallet-api/internal/components/swaps/errors"
	swapModels "trovo-wallet-api/internal/components/swaps/models"
	swapServices "trovo-wallet-api/internal/components/swaps/services"
	usersdb "trovo-wallet-api/internal/components/users/db"
	userModels "trovo-wallet-api/internal/components/users/models"
	userServices "trovo-wallet-api/internal/components/users/services"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/middleware"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/gin-gonic/gin"
)

// Init initializes the controller
func Init(router *gin.Engine, gc *sharedconfig.GlobalConfig) {
	router.POST("/v1/users/swap", middleware.AuthenticationMiddlewareUsingTimestamp(), func(c *gin.Context) {

		signerOwner, getUserError := usersdb.GetUser(middleware.ExtractSigner(c), gc.DB, gc)

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
		{
			//check if pending shared access op exists
			if userServices.CheckPendingSharedAccessApproval(middleware.ExtractPublicKey(c), gc.DB) {
				c.JSON(http.StatusForbidden, gin.H{"error": "error-pending-shared-access-op", "message": "There is a pending shared access operation on this wallet and must be completed first before attempting to send payment from this wallet."})
				return
			}
		}

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

		if signerOwner.PrimarySigner != wallet.Signer {
			c.JSON(http.StatusForbidden, gin.H{"error": "error-unauthorized-access", "message": "You do not have permission on this wallet."})
			return
		}

		if wallet.WalletType == 2 || wallet.WalletType == 3 {
			c.JSON(http.StatusForbidden, gin.H{"error": "error-wallet-type-not-allowed", "message": "Market Making & Bulk Payment wallets are not allowed for this operation."})
			return
		}

		if temp {
			c.JSON(http.StatusBadRequest, (&tErrors.CustomError{Param: "publicKey", Err: "error-temporary-account-forbidden", ErrMessage: "temporary accounts are forbidden from making payment or swap requests", Code: http.StatusForbidden}).JSONError())
			return
		}

		primarySigner := signerOwner.PrimarySigner

		walletOwner, errWalletOwner := userModels.UserWalletID(wallet.ID).GetWalletOwner(gc.DB, gc)

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

		var swapInfo swapModels.SwapSendInfo
		var err error

		data, _ := io.ReadAll(c.Request.Body)

		err = json.Unmarshal(data, &swapInfo)

		var invalidJSON tErrors.ErrorInvalidJSON

		if err != nil {
			log.Print(err)
			c.JSON(invalidJSON.HTTPCode(), invalidJSON.JSONError())
			return
		}

		swapError := swapServices.SwapSend(&signerOwner, &walletOwner, &wallet, &swapInfo, gc)

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
				fromAsset = os.Getenv("NATIVE_ASSET_CODE")
			} else {
				fromAsset = swapInfo.SourceAssetCode
			}
			if swapInfo.DestinationAssetIssuer == "" || swapInfo.DestinationAssetIssuer == "native" {
				toAsset = os.Getenv("NATIVE_ASSET_CODE")
			} else {
				toAsset = swapInfo.DestinationAssetCode
			}
			if swapInfo.TransactionID != "PENDING_AUTH" {

				signerOwner.SendPushMessage(fmt.Sprintf("Trovo: %v Swapped to %v on %v!", fromAsset, toAsset, wallet.Alias), fmt.Sprintf("You have successfully swapped %v%v to %v%v from your wallet with alias %v", swapInfo.SourceAmount, fromAsset, swapInfo.SwappedEstimate, toAsset, wallet.Alias), "", dataPayload, gc)

			}

			//log user current location
			signerOwner.PublicIP = c.ClientIP()
			if len(c.GetHeader("Cf-Connecting-Ip")) > 4 {
				signerOwner.PublicIP = c.GetHeader("Cf-Connecting-Ip")
			}
			swapServices.UpdateAndLogUserSwapGeoInformation(&signerOwner, &swapInfo, gc.DB)
			senderPaymentHistoryCacheKey := fmt.Sprintf("[GET] /v1/users/%v/payments", signerOwner.Username)
			senderCacheKey := fmt.Sprintf("[GET] /v1/users/%v", signerOwner.Username)

			gc.RedisCache.InvalidateCachedHttpResponse(senderCacheKey, senderPaymentHistoryCacheKey)
			c.JSON(http.StatusOK, swapInfo)

		}

	})

	router.POST("/v1/shared-access/swap", middleware.AuthenticationMiddlewareUsingTimestamp(), func(c *gin.Context) {

		{
			//check if pending shared access op exists
			if userServices.CheckPendingSharedAccessApproval(middleware.ExtractPublicKey(c), gc.DB) {
				c.JSON(http.StatusForbidden, gin.H{"error": "error-pending-shared-access-op", "message": "There is a pending shared access operation on this wallet and must be completed first before attempting to send payment from this wallet."})
				return
			}
		}

		signerOwner, getUserError := usersdb.GetUserFromPrimarySigner(middleware.ExtractSigner(c), gc.DB, gc)

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
		if signerOwner.Suspended == 1 {
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
		if wallet.WalletType == 2 || wallet.WalletType == 3 {
			c.JSON(http.StatusForbidden, gin.H{"error": "error-wallet-type-not-allowed", "message": "Market Making & Bulk Payment wallets are not allowed for this operation."})
			return
		}

		if temp {
			c.JSON(http.StatusBadRequest, (&tErrors.CustomError{Param: "publicKey", Err: "error-temporary-account-forbidden", ErrMessage: "temporary accounts are forbidden from making payment or swap requests", Code: http.StatusForbidden}).JSONError())
			return
		}

		walletOwner, errWalletOwner := userModels.UserWalletID(wallet.ID).GetWalletOwner(gc.DB, gc)

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

		if wallet.WalletType != 0 {
			c.JSON(http.StatusForbidden, gin.H{"error": "error-wallet-type-forbidden", "message": "Operation not allowed on any special type of wallets. Only standard wallets are allowed."})
			return
		}
		hasInitiatorAccess := false
		// check if user has initiator access to wallet.
		isViewOnly := wallet.HasViewOnlyAccess(gc)

		if isViewOnly && wallet.Signer != signerOwner.PrimarySigner {
			c.JSON(http.StatusForbidden, gin.H{"error": "error-unauthorized-access", "message": "You do not have permission on this wallet."})
			return
		}

		if wallet.SharedAccessEnabled == 0 {
			c.JSON(http.StatusForbidden, gin.H{"error": "error-unauthorized-access", "message": "Invalid shared access request on wallet without shared access."})
			return
		}
		// check if user has initiator access to wallet.
		if wallet.SharedAccessEnabled == 1 && !isViewOnly {

			for _, p := range signerOwner.WalletsSharedWithUser {
				if p.WalletPublicKey == middleware.ExtractPublicKey(c) && p.TargetUsername == signerOwner.Username && p.Permission == "INITIATOR" {
					hasInitiatorAccess = true
				}
			}
			if !hasInitiatorAccess && !wallet.HasViewOnlyAccess(gc) {
				c.JSON(http.StatusForbidden, gin.H{"error": "error-unauthorized-access", "message": "You do not have an initiator permission on this wallet."})
				return
			}
		}

		var swapInfo swapModels.SwapSendInfo
		var err error

		data, _ := io.ReadAll(c.Request.Body)

		err = json.Unmarshal(data, &swapInfo)

		var invalidJSON tErrors.ErrorInvalidJSON

		if err != nil {
			log.Print(err)
			c.JSON(invalidJSON.HTTPCode(), invalidJSON.JSONError())
			return
		}

		swapError := swapServices.SwapSend(&signerOwner, &walletOwner, &wallet, &swapInfo, gc)

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

		if swapInfo.Commit == 0 {
			c.JSON(http.StatusAccepted, swapInfo)

		} else {
			dataPayload := make(map[string]string)
			dataPayload["route"] = "pendingApproval"

			if swapInfo.TransactionID == "PENDING_AUTH" {
				c.JSON(http.StatusOK, swapInfo)

				{
					accessList := wallet.GetPermissionList(gc.DB)
					// send push notifications
					for _, a := range accessList {

						if a.Permission == "APPROVER" {
							ph, e := usersdb.GetUser(a.TargetUsername, gc.DB, gc)
							if e == nil {
								ph.SendPushMessage(fmt.Sprintf("Trovo: SWAP %v %v awaiting approval!", swapInfo.SourceAmount, swapInfo.Memo), fmt.Sprintf("%v initiated swap request from %v now waiting for an approval. Request: %v", signerOwner.Username, wallet.Alias, swapInfo.ReturnedDescription), "", dataPayload, gc)

							}

						}
					}
					signerOwner.SendPushMessage("Trovo: swap request Submitted!", fmt.Sprintf("You have successfully submitted swap request %v.\nThe listed approvers have been notifed to attend to the request.", swapInfo.ReturnedDescription), "", dataPayload, gc)

				}

				return
			}

			//log user current location
			signerOwner.PublicIP = c.ClientIP()
			if len(c.GetHeader("Cf-Connecting-Ip")) > 4 {
				signerOwner.PublicIP = c.GetHeader("Cf-Connecting-Ip")
			}
			swapServices.UpdateAndLogUserSwapGeoInformation(&signerOwner, &swapInfo, gc.DB)
			senderPaymentHistoryCacheKey := fmt.Sprintf("[GET] /v1/users/%v/payments", signerOwner.Username)
			senderCacheKey := fmt.Sprintf("[GET] /v1/users/%v", signerOwner.Username)

			gc.RedisCache.InvalidateCachedHttpResponse(senderCacheKey, senderPaymentHistoryCacheKey)
			c.JSON(http.StatusOK, swapInfo)

		}

	})

}
