package users

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	usersDB "trovo-wallet-api/internal/components/users/db"
	userModels "trovo-wallet-api/internal/components/users/models"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/network"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/gofrs/uuid"
	"github.com/shopspring/decimal"
	"github.com/stellar/go/txnbuild"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func WalletCountViewOnlyAccess(wallet *userModels.UserWallet, gc *sharedconfig.GlobalConfig) (accessCount uint) {
	if wallet.ManagedAccessEnabled == 0 {
		return 0
	}
	for _, access := range wallet.UserWalletManagedAccess.AccessList {
		if access.AccessLevel == "VIEW-ONLY" {
			accessCount++
		}
	}

	return
}
func WalletHasViewOnlyAccess(wallet *userModels.UserWallet, gc *sharedconfig.GlobalConfig) (viewOnly bool) {
	viewOnly = true
	if wallet.ManagedAccessEnabled == 0 {
		return false
	}
	for _, access := range wallet.UserWalletManagedAccess.AccessList {
		if access.AccessLevel != "VIEW-ONLY" {
			return false
		}
	}

	return
}
func WalletCountAuthorizerAccess(wallet *userModels.UserWallet, gc *sharedconfig.GlobalConfig) (accessCount uint) {
	if wallet.ManagedAccessEnabled == 0 {
		return 0
	}
	for _, access := range wallet.UserWalletManagedAccess.AccessList {
		if access.AccessLevel == "AUTHORIZER" {
			accessCount++
		}
	}

	return
}
func WalletCountInitiatorAccess(wallet *userModels.UserWallet, gc *sharedconfig.GlobalConfig) (accessCount uint) {
	if wallet.ManagedAccessEnabled == 0 {
		return 0
	}
	for _, access := range wallet.UserWalletManagedAccess.AccessList {
		if access.AccessLevel == "INITIATOR" {
			accessCount++
		}
	}

	return
}

func PublicKeyCountViewOnlyAccess(publicKey string, gc *sharedconfig.GlobalConfig) (accessCount uint) {
	if publicKey == "" {
		return 0
	}
	wallet, err := userModels.UserWalletID(publicKey).GetWallet(gc.DB)
	if err != nil {
		return 0
	}

	if wallet.ManagedAccessEnabled == 0 {
		return 0
	}
	for _, access := range wallet.UserWalletManagedAccess.AccessList {
		if access.AccessLevel == "VIEW-ONLY" {
			accessCount++
		}
	}

	return
}

func PublicKeyHasViewOnlyAccess(publicKey string, gc *sharedconfig.GlobalConfig) (viewOnly bool) {
	viewOnly = true
	if publicKey == "" {
		return false
	}
	wallet, err := userModels.UserWalletID(publicKey).GetWallet(gc.DB)
	if err != nil {
		return false
	}

	if wallet.ManagedAccessEnabled == 0 {
		return false
	}
	for _, access := range wallet.UserWalletManagedAccess.AccessList {
		if access.AccessLevel != "VIEW-ONLY" {
			return false
		}
	}

	return
}

func PublicKeyCountAuthorizerAccess(publicKey string, gc *sharedconfig.GlobalConfig) (accessCount uint) {
	if publicKey == "" {
		return 0
	}
	wallet, err := userModels.UserWalletID(publicKey).GetWallet(gc.DB)
	if err != nil {
		return 0
	}

	if wallet.ManagedAccessEnabled == 0 {
		return 0
	}
	for _, access := range wallet.UserWalletManagedAccess.AccessList {
		if access.AccessLevel == "AUTHORIZER" {
			accessCount++
		}
	}

	return
}

func PublicKeyCountInitiatorAccess(publicKey string, gc *sharedconfig.GlobalConfig) (accessCount uint) {
	if publicKey == "" {
		return 0
	}
	wallet, err := userModels.UserWalletID(publicKey).GetWallet(gc.DB)
	if err != nil {
		return 0
	}

	if wallet.ManagedAccessEnabled == 0 {
		return 0
	}
	for _, access := range wallet.UserWalletManagedAccess.AccessList {
		if access.AccessLevel == "INITIATOR" {
			accessCount++
		}
	}

	return
}

func CreateMultiWalletAccess(signerPublicKey string, accessInfo *userModels.UserWalletManagedAccessInfo, gc *sharedconfig.GlobalConfig) (managedAccess userModels.UserWalletManagedAccess, err error) {
	// var managedAccess userModels.UserWalletManagedAccess

	if len(accessInfo.AccessList) == 0 {
		return managedAccess, &tErrors.CustomError{
			Param:      "accessList",
			Err:        "error-access-list-is-empty",
			ErrMessage: "Access list is empty",
			Code:       http.StatusBadRequest,
		}
	}
	var accessList []userModels.WalletAccess
	var numOfAuthorizers, selfAuthorizer uint
	var authorizerUsers []*userModels.User
	e := gc.DB.Preload(clause.Associations).Where("user_wallet_id = ?", accessInfo.PublicKey).First(&managedAccess).Error
	uuid, _ := uuid.NewV4()
	accessID := uuid.String()
	if err != nil {
		if !errors.Is(e, gorm.ErrRecordNotFound) {
			return managedAccess, &tErrors.ErrorTemporaryServerError{}
		}
		//record not found. needs to create it
		// var accessList []userModels.WalletAccess
		// var numOfAuthorizers uint
		walletOwner, e := userModels.UserWalletID(accessInfo.PublicKey).GetWalletOwner(gc.DB)
		if e != nil {
			return managedAccess, &tErrors.CustomError{
				Param:      "username",
				Err:        "error-confirming-wallet-owner",
				ErrMessage: "Unable to confirm wallet owner at this time. Please try again after some minutes.",
				Code:       http.StatusForbidden,
			}
		}
		wallet, e := userModels.UserWalletID(accessInfo.PublicKey).GetWallet(gc.DB)
		if e != nil {
			return managedAccess, &tErrors.CustomError{
				Param:      "username",
				Err:        "error-confirming-wallet",
				ErrMessage: "Unable to confirm wallet at this time. Please try again after some minutes.",
				Code:       http.StatusForbidden,
			}
		}
		for _, v := range accessInfo.AccessList {
			//check if username is valid
			v.Username = strings.ToLower(v.Username)
			if strings.Contains(v.Username, "_") {

				//it is a subwallet and cannot be given access
				return managedAccess, &tErrors.CustomError{
					Param:      "username",
					Err:        "error-subwallet-not-allowed",
					ErrMessage: fmt.Sprintf("Access can only be granted to trovo wallet account, not a subwallet [%v]", v.Username),
					Code:       http.StatusForbidden,
				}
			}
			u, e := usersDB.GetUser(v.Username, gc.DB)
			if e != nil {
				return managedAccess, &tErrors.CustomError{
					Param:      "username",
					Err:        "error-trovo-wallet-account-invalid",
					ErrMessage: fmt.Sprintf("Trovo wallet account [%v] could not be validated at this time.", v.Username),
					Code:       http.StatusForbidden,
				}
			}
			if u.Username != v.Username {
				if e != nil {
					return managedAccess, &tErrors.CustomError{
						Param:      "username",
						Err:        "error-trovo-wallet-account-invalid",
						ErrMessage: fmt.Sprintf("Trovo wallet account [%v] is not Trovo wallet account username.", v.Username),
						Code:       http.StatusForbidden,
					}
				}
			}

			accessList = append(accessList, userModels.WalletAccess{
				UserWalletManagedAccessID: accessID,
				Username:                  v.Username,
				AccessLevel:               v.AccessLevel,
			})
			if v.AccessLevel == "AUTHORIZER" {
				numOfAuthorizers++
				if v.Username == walletOwner.Username {
					selfAuthorizer = 1
				}
			}

			authorizerUsers = append(authorizerUsers, &u)
		}
		if numOfAuthorizers < accessInfo.NumberOfAuthorizers {
			//number of authorizers does not reach the minimum threshold needed. cannot proceed so as to prevent account lockout
			return managedAccess, &tErrors.CustomError{
				Param:      "numberOfAuthorizers",
				Err:        "error-authorizers-not-enough",
				ErrMessage: fmt.Sprintf("Authorizers (%v) assigned to wallet are less than the number (%v) required to authorize transaction. Try specifying the number of authorizers are reduce the number of authorizers required to %v", numOfAuthorizers, accessInfo.NumberOfAuthorizers, numOfAuthorizers),
				Code:       http.StatusForbidden,
			}
		}

		managedAccess = userModels.UserWalletManagedAccess{
			ID:                  accessID,
			UserWalletID:        accessInfo.PublicKey,
			NumberOfAuthorizers: accessInfo.NumberOfAuthorizers,
			AccessList:          accessList,
		}
		// if authorizers exists, then owner must sign transaction to add them as signers
		if (numOfAuthorizers - selfAuthorizer) > 0 {
			accessInfo.SignatureRequired = 1

		}
		xdrBase64, messages, walletMustSign, errGenXdr := generateCreateMultiWalletAccessXdr(&wallet, authorizerUsers, int(accessInfo.NumberOfAuthorizers), gc)
		if errGenXdr != nil {
			return managedAccess, errGenXdr
		}
		accessInfo.Messages = messages
		if walletMustSign {
			accessInfo.SignatureRequired = 1

		}
		if xdrBase64 == "no-ops" {
			accessInfo.Transaction = xdrBase64
		}
		if len(accessInfo.TransactionSignature) == 0 {
			return managedAccess, nil
		}

	}
	//existing access was retrieved

	return managedAccess, &tErrors.CustomError{
		Param:      "id",
		Err:        "error-access-management-already-active-on-wallet",
		ErrMessage: "Access management already activated on wallet. Use option to update the access or update the access list.",
		Code:       http.StatusForbidden,
	}
}

func UpdateMultiWalletAccess(signerPublicKey string, accessInfo *userModels.UserWalletManagedAccessInfo, gc *sharedconfig.GlobalConfig) (managedAccess userModels.UserWalletManagedAccess, err error) {
	// var managedAccess userModels.UserWalletManagedAccess

	if len(accessInfo.AccessList) == 0 {
		return managedAccess, &tErrors.CustomError{
			Param:      "accessList",
			Err:        "error-access-list-is-empty",
			ErrMessage: "Access list is empty",
			Code:       http.StatusBadRequest,
		}
	}
	var accessList []userModels.WalletAccess
	var numOfAuthorizers uint
	e := gc.DB.Preload(clause.Associations).Where("user_wallet_id = ?", accessInfo.PublicKey).First(&managedAccess).Error

	if err != nil {
		if !errors.Is(e, gorm.ErrRecordNotFound) {
			return managedAccess, &tErrors.ErrorTemporaryServerError{}
		}
		//record not found. needs to create it
		// var accessList []userModels.WalletAccess
		// var numOfAuthorizers uint
		for _, v := range accessInfo.AccessList {
			//check if username is valid
			v.Username = strings.ToLower(v.Username)
			if strings.Contains(v.Username, "_") {

				//it is a subwallet and cannot be given access
				return managedAccess, &tErrors.CustomError{
					Param:      "username",
					Err:        "error-subwallet-not-allowed",
					ErrMessage: fmt.Sprintf("Access can only be granted to trovo wallet account, not a subwallet [%v]", v.Username),
					Code:       http.StatusForbidden,
				}
			}
			u, e := usersDB.GetUser(v.Username, gc.DB)
			if e != nil {
				return managedAccess, &tErrors.CustomError{
					Param:      "username",
					Err:        "error-trovo-wallet-account-invalid",
					ErrMessage: fmt.Sprintf("Trovo wallet account [%v] could not be validated at this time.", v.Username),
					Code:       http.StatusForbidden,
				}
			}
			if u.Username != v.Username {
				if e != nil {
					return managedAccess, &tErrors.CustomError{
						Param:      "username",
						Err:        "error-trovo-wallet-account-invalid",
						ErrMessage: fmt.Sprintf("Trovo wallet account [%v] is not Trovo wallet account username.", v.Username),
						Code:       http.StatusForbidden,
					}
				}
			}

		}
		if numOfAuthorizers < accessInfo.NumberOfAuthorizers {
			//number of authorizers does not reach the minimum threshold needed. cannot proceed so as to prevent account lockout
			return managedAccess, &tErrors.CustomError{
				Param:      "numberOfAuthorizers",
				Err:        "error-authorizers-not-enough",
				ErrMessage: fmt.Sprintf("Authorizers (%v) assigned to wallet are less than the number (%v) required to authorize transaction. Try specifying the number of authorizers are reduce the number of authorizers required to %v", numOfAuthorizers, accessInfo.NumberOfAuthorizers, numOfAuthorizers),
				Code:       http.StatusForbidden,
			}
		}
		uuid, _ := uuid.NewV4()
		accessID := uuid.String()
		managedAccess = userModels.UserWalletManagedAccess{
			ID:                  accessID,
			UserWalletID:        accessInfo.PublicKey,
			NumberOfAuthorizers: accessInfo.NumberOfAuthorizers,
			AccessList:          accessList,
		}
	}
	//existing access was retrieved
	if managedAccess.ID != accessInfo.UserWalletManagedAccessID {
		return managedAccess, &tErrors.CustomError{
			Param:      "userWalletManagedAccessID",
			Err:        "error-user-wallet-managed-access-id-mismatch",
			ErrMessage: "Access ID mismatch",
			Code:       http.StatusBadRequest,
		}
	}
	for _, v := range accessInfo.AccessList {
		//check if username is valid
		v.Username = strings.ToLower(v.Username)
		if strings.Contains(v.Username, "_") {

			//it is a subwallet and cannot be given access
			return managedAccess, &tErrors.CustomError{
				Param:      "username",
				Err:        "error-subwallet-not-allowed",
				ErrMessage: fmt.Sprintf("Access can only be granted to trovo wallet account, not a subwallet [%v]", v.Username),
				Code:       http.StatusForbidden,
			}
		}
		u, e := usersDB.GetUser(v.Username, gc.DB)
		if e != nil {
			return managedAccess, &tErrors.CustomError{
				Param:      "username",
				Err:        "error-trovo-wallet-account-invalid",
				ErrMessage: fmt.Sprintf("Trovo wallet account [%v] could not be validated at this time.", v.Username),
				Code:       http.StatusForbidden,
			}
		}
		if u.Username != v.Username {
			if e != nil {
				return managedAccess, &tErrors.CustomError{
					Param:      "username",
					Err:        "error-trovo-wallet-account-invalid",
					ErrMessage: fmt.Sprintf("Trovo wallet account [%v] is not Trovo wallet account username.", v.Username),
					Code:       http.StatusForbidden,
				}
			}
		}

		if v.AccessLevel == "AUTHORIZER" {
			numOfAuthorizers++
		}
	}
	if numOfAuthorizers < accessInfo.NumberOfAuthorizers {
		//number of authorizers does not reach the minimum threshold needed. cannot proceed so as to prevent account lockout
		return managedAccess, &tErrors.CustomError{
			Param:      "numberOfAuthorizers",
			Err:        "error-authorizers-not-enough",
			ErrMessage: fmt.Sprintf("Authorizers (%v) assigned to wallet are less than the number (%v) specified as required to authorize transaction. Try specifying the number of authorizers to be at least equal the number of required %v", numOfAuthorizers, accessInfo.NumberOfAuthorizers, numOfAuthorizers),
			Code:       http.StatusForbidden,
		}
	}
	return
}

func generateCreateMultiWalletAccessXdr(wallet *userModels.UserWallet, authorizers []*userModels.User, authThreshold int, gc *sharedconfig.GlobalConfig) (xdrbase64 string, messages []string, walletMustSign bool, err error) {
	client := gc.BantuExpansionClient
	ops := make([]txnbuild.Operation, 0)
	messages = make([]string, 0)
	totalNativeBalanceNeeded := decimal.Zero
	var activationAmount = decimal.NewFromFloat(6)
	var minBalance = decimal.NewFromFloat(3.0)
	if len(os.Getenv("SUB_WALLET_ACTIVATION_AMOUNT")) > 0 {
		activationAmount = decimal.RequireFromString(os.Getenv("SUB_WALLET_ACTIVATION_AMOUNT"))
	}
	if len(os.Getenv("WALLET_MINIMUM_BALANCE")) > 0 {
		minBalance = decimal.RequireFromString(os.Getenv("WALLET_MINIMUM_BALANCE"))
	}

	if len(authorizers) == 0 {
		err = &tErrors.CustomError{
			Param:      "numberOfAuthorizers",
			Err:        "error-no-authorizers-specified",
			ErrMessage: fmt.Sprintf("%v account does not have enough XBN balance to perform this operation", wallet.Alias),
			Code:       404,
		}
		return "", messages, walletMustSign, err
	}
	//check if it is first call to create sub-wallet

	//populate the subwallet Info and generate the transaction

	//check if primary account has native enough native balance
	var nativeAsset txnbuild.Asset = txnbuild.NativeAsset{}
	walletAccountExists, _, walletAccountNativeBalance, _, walletSourceAccount, errWalletAct := network.BlockchainAccountProperties(client, wallet.ID, nativeAsset)
	if errWalletAct != nil {
		log.Printf("[generateCreateMultiWalletAccessXdr] by [%v] for MultiAccess Account Properties error:[%v] \n", wallet.Alias, errWalletAct)

		return "", messages, walletMustSign, errWalletAct
	}
	if !walletAccountExists || (walletAccountNativeBalance.Sub(activationAmount)).LessThan(minBalance) {
		log.Printf("[generateCreateMultiWalletAccessXdr] by [%v] MultiAccess WalletAccount underfunded \n", wallet.Alias)

		err = &tErrors.CustomError{
			Param:      "publicKey",
			Err:        "error-wallet-underfunded",
			ErrMessage: fmt.Sprintf("Wallet %v does not have enough XBN balance to perform this operation", wallet.Alias),
			Code:       404,
		}
		return "", messages, walletMustSign, err
	}
	//check access list to know if you would activate the user wallets before proceeding.
	for _, user3p := range authorizers {
		var totalUsersToFund int64
		//ensure u r using the account signer, since the account may have been recovered, or may be recovered in the future, changing the signer, but retaining the primary key
		authorizerAccountExists, _, _, _, _, _ := network.BlockchainAccountProperties(client, user3p.PrimarySigner, nativeAsset)
		if !authorizerAccountExists {
			//if subwallet is not activated
			//build transaction that will activate the primary signer from the assigning wallet

			ops = append(ops, &txnbuild.CreateAccount{
				Destination:   user3p.PrimarySigner,
				Amount:        activationAmount.String(),
				SourceAccount: wallet.ID,
			})

			//after creation, it now exists with enough balance to add signer wallet as signer
			ops = append(ops, &txnbuild.SetOptions{
				Signer: &txnbuild.Signer{
					Address: user3p.PrimarySigner,
					Weight:  1,
				},
				SourceAccount: wallet.ID,
			})
			totalUsersToFund++

			//since an operation now exists, wallet must sign
			walletMustSign = true
			// messages = append(messages, fmt.Sprintf("Important: %v XBN will be deducted from wallet %v to used to activate the sub-wallet.", activationAmount.String()))

		}

		if authorizerAccountExists {
			//account exists, check if it already it a signer in the wallet

			//after topping up, it now has enough balance to add primary wallet as signer if it is not already a signer
			if !wallet.SignerIsValidWA(user3p.PrimarySigner, walletSourceAccount) {

				ops = append(ops, &txnbuild.SetOptions{
					Signer: &txnbuild.Signer{
						Address: user3p.PrimarySigner,
						Weight:  1,
					},
					SourceAccount: wallet.ID,
				})

				walletMustSign = true
			}

		}
		if totalUsersToFund > 0 {
			totalNativeBalanceNeeded = activationAmount.Mul(decimal.NewFromInt(totalUsersToFund))
			if walletAccountNativeBalance.LessThan(totalNativeBalanceNeeded) {
				//not enough balance to perform this.
				log.Printf("[generateCreateMultiWalletAccessXdr] by [%v] MultiAccess WalletAccount underfunded \n", wallet.Alias)

				err = &tErrors.CustomError{
					Param:      "publicKey",
					Err:        "error-wallet-underfunded",
					ErrMessage: fmt.Sprintf("Wallet %v needs more than %v XBN balance to perform this operation", wallet.Alias, totalNativeBalanceNeeded.String()),
					Code:       404,
				}
				return "", messages, walletMustSign, err
			}
		}
	}

	//TODO: if account exists and subwallet has enough balance, we add the operation to pay TROVO fee from primary Wallet

	// Construct the transaction that holds the operations to execute on the network
	if len(ops) == 0 {
		// no operations to sign
		return "no-ops", messages, walletMustSign, nil
	}

	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        walletSourceAccount,
			IncrementSequenceNum: true,
			Operations:           ops,
			BaseFee:              txnbuild.MinBaseFee,
			Preconditions: txnbuild.Preconditions{
				TimeBounds: txnbuild.NewInfiniteTimeout(),
			},
			Memo: txnbuild.MemoText("Create Signers"),
		},
	)
	if err != nil {
		log.Println("[generateCreateMultiWalletAccessXdr] error constructing transaction ", err)
		return "", messages, walletMustSign, err
	}

	var xdrBase64 string

	xdrBase64, err = tx.Base64()
	if err != nil {
		log.Println("[generateCreateMultiWalletAccessXdr] error getting txn base64", err)
		return "", messages, walletMustSign, err
	}
	messages = append(messages, fmt.Sprintf("Important: %v XBN will be deducted from wallet %v to be used to activate/fund the authorizer(s).", totalNativeBalanceNeeded.String(), wallet.Alias))

	return xdrBase64, messages, walletMustSign, nil

}
