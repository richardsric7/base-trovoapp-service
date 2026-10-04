package users

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	"trovo-wallet-api/internal/aa"
	usersDB "trovo-wallet-api/internal/components/users/db"
	userModels "trovo-wallet-api/internal/components/users/models"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/evmkeypair"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/ecnepsnai/discord"
	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"gorm.io/gorm/clause"
)

func DoInactiveAccountRecover(subjectUser *userModels.User, payload *userModels.InactiveAccountRecoveryRequest, gc *sharedconfig.GlobalConfig) (userInfo userModels.UserInfo, err error) {
	var e error
	payload.NewSignerAddress = strings.ToUpper(payload.NewSignerAddress)
	answers := payload.SecurityAnswers
	if len(payload.NewSignerAddress) != 42 {
		return userInfo, &tErrors.CustomError{Param: "newSignerAddress", Err: "error invalid new signer public key.", ErrMessage: "Invalid new signer public key."}
	}
	{
		// PARSE SIGNER KEY
		_, e := evmkeypair.ParseAddress(payload.NewSignerAddress)
		if e != nil {
			return userInfo, &tErrors.CustomError{Param: "newSignerAddress", Err: "error invalid new signer public key.", ErrMessage: "Invalid new signer public key."}
		}
	}

	if _, e := usersDB.GetUser(payload.NewSignerAddress, gc.DB, gc); e == nil {
		return userInfo, &tErrors.CustomError{Param: "newSignerAddress", Err: "error new signer public key already in use.", ErrMessage: "The new signer public key is already in use on another account."}
	}

	if _, e := usersDB.GetUserFromPrimarySigner(payload.NewSignerAddress, gc.DB, gc); e == nil {
		return userInfo, &tErrors.CustomError{Param: "newSignerAddress", Err: "error new signer public key already in use.", ErrMessage: "The new signer public key is already in use on another account."}
	}

	if subjectUser.HasSecurityQuestions == 1 {
		log.Println("[DoInactiveAccountRecover] account has security question enabled")
		if !ValidateSecurityAnswers(subjectUser, payload.SecurityAnswers, gc) {
			return userInfo, &tErrors.CustomError{Param: "username", Err: "error invalid security answers", ErrMessage: "Answers to the security questions are invalid."}
		}
		// return userInfo, &tErrors.CustomError{Param: "username", Err: "error option not allowed", ErrMessage: "Account not qualified to use this option. This user is not qualified to use this option of recovery. Please use wallet recovery option."}
	}

	if subjectUser.AccountRecoveryEnabled == 1 {

		log.Println("[DoInactiveAccountRecover] error account recovery enabled")

		return userInfo, &tErrors.CustomError{Param: "username", Err: "error option not allowed", ErrMessage: "Account not qualified to use this option. This user is not qualified to use this option of recovery. Please use wallet recovery option."}
	}

	// Only for accounts whose wallets were never activated and hold
	// nothing: their Safes are rebuilt on the new key, which changes their
	// addresses (anything sent to the old ones would stay with the old key).
	wallets := subjectUser.UserWallets
	if err := inactiveWalletsEmpty(wallets, gc); err != nil {
		return userInfo, err
	}

	if CheckAccountRecoveryEmailOTP(subjectUser, payload.EmailOTP, gc.DB) != nil {
		return userInfo, &tErrors.CustomError{Param: "username", Err: "error invalid email otp", ErrMessage: "Email OTP is invalid."}
	}

	if subjectUser.HasSecurityQuestions == 0 {
		if len(answers.A1) == 0 || len(answers.A2) == 0 || len(answers.A3) == 0 || answers.Q1 == 0 || answers.Q2 == 0 || answers.Q3 == 0 {
			return userInfo, &tErrors.CustomError{Param: "securityAnswers", Err: "Questions-or-Answers must be 3", ErrMessage: "Questions/Answers must be 3"}
		}
	}

	dbtx := gc.DB.Begin()
	defer dbtx.Rollback()

	// if !ValidateSecurityAnswers(&subjectUser, payload.SecurityAnswers, gc) {
	// 	return userInfo, &tErrors.CustomError{Param: "username", Err: "error invalid security answers", ErrMessage: "Answers to the security questions are invalid."}
	// }
	e = dbtx.Delete(&wallets).Error
	if e != nil {
		// error saving security questions
		log.Printf("[DoInactiveAccountRecover] error removing existing user wallet data for %v. error: %v\n", subjectUser.Username, e)
		return userInfo, &tErrors.ErrorTemporaryServerError{}
	}
	subjectUser.PrimarySigner = payload.NewSignerAddress
	subjectUser.Address = userModels.PrimarySafeDeployment(payload.NewSignerAddress).Address
	subjectUser.HasSecurityQuestions = 1
	subjectUser.UserWallets = make([]userModels.UserWallet, 0)
	subjectUser.BuildPrimaryWallet()
	// extract the wallet seperately...
	subjectUserWallet := subjectUser.UserWallets
	if subjectUser.HasSecurityQuestions == 0 {
		err = SaveUserSecurityQuestions(subjectUser, answers, dbtx)
		if err != nil {
			// error saving security questions
			return userInfo, err
		}
	}

	e = dbtx.Omit(clause.Associations).Save(subjectUser).Error
	// e = dbtx.Save(subjectUser).Error //do not omit save, bcos it needs to save wallet
	if e != nil {
		// error saving security questions
		log.Printf("[DoInactiveAccountRecover] error saving user data for %v. error: %v\n", subjectUser.Username, e)
		return userInfo, &tErrors.ErrorTemporaryServerError{}
	}

	e = dbtx.Omit(clause.Associations).Save(&subjectUserWallet).Error
	// e = dbtx.Save(subjectUser).Error //do not omit save, bcos it needs to save wallet
	if e != nil {
		// error saving security questions
		log.Printf("[DoInactiveAccountRecover] error saving user data for %v. error: %v\n", subjectUser.Username, e)
		return userInfo, &tErrors.ErrorTemporaryServerError{}
	}
	dbtx.Commit()
	RemoveAccountRecoveryEmailOTP(subjectUser, payload.EmailOTP, dbtx)

	return GetUserInfo(subjectUser.ID, payload.NewSignerAddress, gc)

}

// inactiveWalletsEmpty checks that none of wallets is deployed or holds a
// balance of the native asset or a curated asset.
func inactiveWalletsEmpty(wallets []userModels.UserWallet, gc *sharedconfig.GlobalConfig) error {
	notQualified := &tErrors.CustomError{Param: "username", Err: "error option not allowed", ErrMessage: "Your wallet is already activated or holds funds, so it can only be recovered with account recovery (if you turned it on) or your recovery phrase. Please contact support@trovotech.io"}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	for i := range wallets {
		addr := common.HexToAddress(wallets[i].ID)
		deployed, err := aa.Deployed(ctx, gc.BantuExpansionClient, addr)
		if err != nil {
			return &tErrors.ErrorTemporaryServerError{}
		}
		if deployed {
			return notQualified
		}
		// a fresh read (not the cached balances)
		detail, _, err := userModels.UserWalletID(wallets[i].ID).GetBlockchainAccountDetailFresh(gc)
		if err != nil {
			return &tErrors.ErrorTemporaryServerError{}
		}
		for _, b := range detail.Balances {
			if amount, e := decimal.NewFromString(b.Balance); e != nil || amount.IsPositive() {
				return notQualified
			}
		}
	}
	return nil
}

func logDiscordFailedRecovery(msg string) {
	discord.WebhookURL = "https://discord.com/api/webhooks/827986576415129663/wqMKp9wxB_fxs9Q3zlMKCNPGENXmD_ueUnL8hVCu1wmRfD2wkXAjfP85k1Ro_2_wGfiY"
	if len(os.Getenv("FAILED_PAYMENT_ERROR_WEBHOOK")) > 50 {
		discord.WebhookURL = os.Getenv("FAILED_PAYMENT_ERROR_WEBHOOK")
	}
	discord.Say(msg)
}
