package users

import (
	"encoding/json"
	"net/http"

	usersDB "trovo-wallet-api/internal/components/users/db"
	userModels "trovo-wallet-api/internal/components/users/models"
	userServices "trovo-wallet-api/internal/components/users/services"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/middleware"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/gin-gonic/gin"
)

// stablerailAvailable answers 503 while Stablerail is not configured and
// enabled.
func stablerailAvailable(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !userServices.StablerailEnabled(gc) {
			c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{"error": "error-fiat-unavailable", "message": "Bank deposits and withdrawals are not available at this time."})
			return
		}
		c.Next()
	}
}

// getUsersStablerailProfileHandler godoc
// @Summary GET /v1/users/stablerail/profile
// @Description Whether bank deposits and withdrawals are available, whether the user's BVN is verified for them, and the smallest bank withdrawal.
// @Tags users
// @Produce json
// @Success 200 {object} userModels.StablerailProfile
// @Failure 400 {object} map[string]interface{}
// @Router /v1/users/stablerail/profile [get]
func getUsersStablerailProfileHandler(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := usersDB.GetUserFromPrimarySigner(middleware.ExtractSigner(c), gc.DB, gc)
		if err != nil {
			writeError(c, err)
			return
		}
		c.JSON(http.StatusOK, userServices.GetStablerailProfile(&user, gc))
	}
}

// postUsersStablerailWithdrawHandler godoc
// @Summary POST /v1/users/stablerail/withdraw
// @Description Withdraws cNGN from the wallet to a Nigerian bank account. Like a payment it takes two calls: the first (amount, accountNumber, bankCode) returns the transaction to sign; the second (the same body plus transaction and transactionSignature, or commit for shared wallets) submits it. The cNGN goes to the user's Stablerail wallet, which then pays the bank account; follow it with GET /v1/users/stablerail/withdrawals.
// @Tags users
// @Accept json
// @Produce json
// @Param body body userModels.BankWithdrawalRequest true "Withdrawal"
// @Success 200 {object} userModels.BankWithdrawalRequest
// @Failure 400 {object} map[string]interface{}
// @Failure 503 {object} map[string]interface{}
// @Router /v1/users/stablerail/withdraw [post]
func postUsersStablerailWithdrawHandler(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req userModels.BankWithdrawalRequest
		if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
			var invalidJSON tErrors.ErrorInvalidJSON
			c.JSON(invalidJSON.HTTPCode(), invalidJSON.JSONError())
			return
		}
		signerUser, walletOwner, wallet, ok := marketWallet(c, gc)
		if !ok {
			return
		}
		if err := userServices.WithdrawToBank(&signerUser, &walletOwner, &wallet, &req, gc); err != nil {
			writeError(c, err)
			return
		}
		c.JSON(http.StatusOK, req)
	}
}

// getUsersStablerailWithdrawalsHandler godoc
// @Summary GET /v1/users/stablerail/withdrawals
// @Description The user's bank withdrawals, newest first, with their status: DEPOSITING, DEPOSITED, REQUESTING, DEPOSIT_FAILED, REQUEST_FAILED, then Stablerail's (pending, processing, bank_verification, transfer_pending, transfer_confirmed, payout_pending, completed, failed, cancelled).
// @Tags users
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /v1/users/stablerail/withdrawals [get]
func getUsersStablerailWithdrawalsHandler(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := usersDB.GetUserFromPrimarySigner(middleware.ExtractSigner(c), gc.DB, gc)
		if err != nil {
			writeError(c, err)
			return
		}
		list, err := userServices.ListBankWithdrawals(user.Username, gc)
		if err != nil {
			writeError(c, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"withdrawals": list})
	}
}
