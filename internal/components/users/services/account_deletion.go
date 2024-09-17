package users

import (
	"fmt"
	"log"
	"os"
	"time"
	userBc "trovo-wallet-api/internal/components/users/blockchain"
	userModels "trovo-wallet-api/internal/components/users/models"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/network"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/stellar/go/protocols/horizon"
	"github.com/stellar/go/txnbuild"
)

func AccountDeletion(user *userModels.User, payload *userModels.UserAccountDeletionPayload, gc *sharedconfig.GlobalConfig) (err error) {
	var e error
	client := gc.BantuExpansionClient
	multiAccessWallets := make([]userModels.UserWallet, 0)
	ops := make([]txnbuild.Operation, 0)
	messages := make([]string, 0)
	payload.Messages = make([]string, 0)
	if user.HasSecurityQuestions == 0 {
		return &tErrors.CustomError{Param: "username", Err: "error security answers not set", ErrMessage: "security answers has not been set for this account."}
	}
	if user.HasSharedAccessInAnyWallet(gc) {
		return &tErrors.CustomError{Param: "username", Err: "error shared access enabled in wallet", ErrMessage: "Account cannot be deleted while shared access is still active on any wallet. Please ensure all shared access is removed from this account before proceeding."}
	}
	if user.AccountRecoveryEnabled == 1 {
		return &tErrors.CustomError{Param: "username", Err: "error account recovery enabled", ErrMessage: "Account cannot be deleted while account recovery is still active. Please disable account recovery on this account before proceeding."}
	}

	wallet, _ := userModels.UserWalletID(user.PublicKey).GetWallet(gc.DB, gc)
	homeDomain := "trovotech.io"
	ops = append(ops, &txnbuild.SetOptions{
		HomeDomain:    &homeDomain,
		SourceAccount: wallet.ID,
	})

	dbtx := gc.DB.Begin()
	defer dbtx.Rollback()
	user.AccountRecoveryEnabled = 0
	exp := time.Now().AddDate(1, 0, 0)
	user.AccountRecoveryExpiresOn = &exp
	dbErr := dbtx.Save(user).Error
	if dbErr != nil {
		log.Printf("[AccountDeletion] Error saving account recovery state: %v\n", dbErr)
		return &tErrors.ErrorTemporaryServerError{}
	}

	var userAccount horizon.Account
	if userAccount, e = userBc.GetBlockchainAccountDetail(user.PublicKey); e != nil {
		if e.Error() == "error-blockchain-account-not-activated" {
			return &tErrors.CustomError{Param: "username", Err: "error primary account not yet activated", ErrMessage: fmt.Sprintf("Primary account is not yet activated. Please send upto 50 %v to the primary wallet to continue.", os.Getenv("NATIVE_ASSET_CODE"))}
		}
		log.Printf("[AccountDeletion] Error on blockchain validating primary account %v\n", e)
		return &tErrors.ErrorTemporaryServerError{}

	}

	var xdrBase64 string
	payload.Messages = messages
	payload.NetworkPassPhrase = network.GetBlockchainNetworkPassPhrase()
	// oldTrx := payload.Transaction

	// generate new transaction

	memo := "Account Deletion Request"
	log.Println("[AccountDeletion] Memo:", memo)

	if len(ops) == 0 {
		return &tErrors.CustomError{Param: "username", Err: "error no operations to perform", ErrMessage: "Could not find any operations to perform for this action."}

	}
	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &userAccount,
			IncrementSequenceNum: true,
			Operations:           ops,
			BaseFee:              2000,
			Preconditions:        txnbuild.Preconditions{TimeBounds: txnbuild.NewInfiniteTimeout()},
			Memo:                 txnbuild.MemoText(memo),
		},
	)

	if err != nil {
		log.Println("[AccountDeletion]error constructing transaction ", err)
		return &tErrors.ErrorTemporaryServerError{}
	}

	xdrBase64, err = tx.Base64()
	if err != nil {
		log.Println("[AccountDeletion] error getting txn base64", err)
		return &tErrors.ErrorTemporaryServerError{}
	}
	if len(payload.Transaction) == 0 || len(payload.TransactionSignature) == 0 {
		payload.Transaction = xdrBase64

		return nil
	}

	if len(payload.TransactionSignature) == 0 {
		return &tErrors.CustomError{Param: "TransactionSignature", Err: "error transaction signature is required", ErrMessage: "Transaction signature is required."}
	}
	//submit to blockchain. sending the transaction that was signed, bcos the new one may differ based on quantity of fee asset

	txnHash, err := network.SubmitXdrWithSignature(client, user.PrimarySigner, payload.Transaction, payload.TransactionSignature)
	if err != nil {
		logDiscordFailedRecovery(fmt.Sprintf("Error submitting account recovery enable [%+v] transaction: %s", payload, err.Error()))
		return
	}
	payload.TransactionID = txnHash
	dbtx.Commit()
	if err == nil && len(multiAccessWallets) > 0 {
		log.Printf("Skipped wallets %+v\n", multiAccessWallets)
	}
	user.InvalidateUserCache(gc)
	owner, _ := userModels.Username(user.Username).GetFullUser(gc.DB, gc)
	user = &owner

	return nil

}
