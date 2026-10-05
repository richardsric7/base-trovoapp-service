package users

import (
	"encoding/json"
	"io"
	"net/http"

	usersDB "trovo-wallet-api/internal/components/users/db"
	userModels "trovo-wallet-api/internal/components/users/models"
	userServices "trovo-wallet-api/internal/components/users/services"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/middleware"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/gin-gonic/gin"
)

func writeError(c *gin.Context, err error) {
	if ex, ok := err.(tErrors.GenericError); ok {
		c.JSON(ex.HTTPCode(), ex.JSONError())
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": err.Error()})
}

// marketWallet loads the signer, the wallet owner and the standard wallet a
// market-making request acts on, writing the error response when it fails.
func marketWallet(c *gin.Context, gc *sharedconfig.GlobalConfig) (signerUser, walletOwner userModels.User, wallet userModels.UserWallet, ok bool) {
	var err error
	if signerUser, err = usersDB.GetSlimUserFromPrimarySigner(middleware.ExtractSigner(c), gc.DB, gc); err != nil {
		writeError(c, err)
		return
	}
	if walletOwner, err = usersDB.GetUser(middleware.ExtractAddress(c), gc.DB, gc); err != nil {
		writeError(c, err)
		return
	}
	if wallet, err = usersDB.GetWallet(middleware.ExtractAddress(c), gc.DB); err != nil {
		writeError(c, err)
		return
	}
	if wallet.WalletType != 0 {
		c.JSON(http.StatusForbidden, gin.H{"error": "error-wallet-type-forbidden", "message": "Operation not allowed on any special type of wallets. Only standard wallets are allowed."})
		return
	}
	if signerUser.PrimarySigner != wallet.Signer {
		c.JSON(http.StatusForbidden, gin.H{"error": "error-unauthorized-access", "message": "You do not have permission on this wallet."})
		return
	}
	return signerUser, walletOwner, wallet, true
}

// getUsersTradesHandler godoc
// @Summary GET /v1/users/trades
// @Description The wallet's market-making offers, with what each still sells and whether it is open on the offer book.
// @Tags users
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /v1/users/trades [get]
func getUsersTradesHandler(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		_, _, wallet, ok := marketWallet(c, gc)
		if !ok {
			return
		}
		offers, err := userServices.ListOffers(&wallet, gc)
		if err != nil {
			writeError(c, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"offers": offers})
	}
}

// deleteUsersTradesIDHandler godoc
// @Summary DELETE /v1/users/trades/:id
// @Description Cancels a market-making offer, returning what is left of it to the wallet. Like placing one it takes two calls: the first returns the transaction to sign, the second (with transaction and transactionSignature, or commit for shared wallets) submits it.
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "Market offer id"
// @Param body body userModels.DeleteOfferRequest false "Signed transaction"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /v1/users/trades/{id} [delete]
func deleteUsersTradesIDHandler(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req userModels.DeleteOfferRequest
		if data, _ := io.ReadAll(c.Request.Body); len(data) > 0 {
			if err := json.Unmarshal(data, &req); err != nil {
				var invalidJSON tErrors.ErrorInvalidJSON
				c.JSON(invalidJSON.HTTPCode(), invalidJSON.JSONError())
				return
			}
		}
		req.ID = c.Param("id")
		signerUser, walletOwner, wallet, ok := marketWallet(c, gc)
		if !ok {
			return
		}
		if err := userServices.CancelOffer(&signerUser, &walletOwner, &wallet, &req, gc); err != nil {
			writeError(c, err)
			return
		}
		walletOwner.InvalidateUserCache(gc)
		c.JSON(http.StatusOK, req)
	}
}
