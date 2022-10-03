package users

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	bc "trovo-wallet-api/internal/blockchainalgofuncs"
	userBc "trovo-wallet-api/internal/components/users/blockchain"
	usersDB "trovo-wallet-api/internal/components/users/db"
	userModels "trovo-wallet-api/internal/components/users/models"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/network"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/gofrs/uuid"
	"github.com/shopspring/decimal"
	"github.com/stellar/go/txnbuild"
	"gorm.io/gorm"
)

func WalletCountViewOnlyAccess(wallet *userModels.UserWallet, gc *sharedconfig.GlobalConfig) (accessCount uint) {
	if wallet.SharedAccessEnabled == 0 {
		return 0
	}
	for _, access := range wallet.UserWalletSharedAccess.AccessList {
		if access.AccessLevel == "VIEW-ONLY" {
			accessCount++
		}
	}

	return
}
func WalletHasViewOnlyAccess(wallet *userModels.UserWallet, gc *sharedconfig.GlobalConfig) (viewOnly bool) {
	viewOnly = true
	if wallet.SharedAccessEnabled == 0 {
		return false
	}
	for _, access := range wallet.UserWalletSharedAccess.AccessList {
		if access.AccessLevel != "VIEW-ONLY" {
			return false
		}
	}

	return
}
func WalletCountAuthorizerAccess(wallet *userModels.UserWallet, gc *sharedconfig.GlobalConfig) (accessCount uint) {
	if wallet.SharedAccessEnabled == 0 {
		return 0
	}
	for _, access := range wallet.UserWalletSharedAccess.AccessList {
		if access.AccessLevel == "AUTHORIZER" {
			accessCount++
		}
	}

	return
}
func WalletCountInitiatorAccess(wallet *userModels.UserWallet, gc *sharedconfig.GlobalConfig) (accessCount uint) {
	if wallet.SharedAccessEnabled == 0 {
		return 0
	}
	for _, access := range wallet.UserWalletSharedAccess.AccessList {
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

	if wallet.SharedAccessEnabled == 0 {
		return 0
	}
	for _, access := range wallet.UserWalletSharedAccess.AccessList {
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
	accessList := userModels.UserWalletID(publicKey).GetAccessList(gc.DB)
	if len(accessList) == 0 {
		return false
	}

	for _, access := range accessList {
		if access.AccessLevel != "VIEW-ONLY" {
			return false
		}
	}

	return
}

func PublicKeyHasViewOnlyAccessWACL(publicKey string, accessList []userModels.WalletAccessInfo, gc *sharedconfig.GlobalConfig) (viewOnly bool) {
	viewOnly = true
	if publicKey == "" {
		return false
	}

	for _, access := range accessList {
		if access.AccessLevel != "VIEW-ONLY" {
			return false
		}
	}

	return
}

func PublicKeyCountApproverAccessWACL(publicKey string, accessList []userModels.WalletAccessInfo, gc *sharedconfig.GlobalConfig) (accessCount uint) {
	if publicKey == "" {
		return 0
	}

	for _, access := range accessList {
		if access.AccessLevel == "APPROVER" {
			accessCount++
		}
	}

	return
}
func PublicKeyCountApproverAccess(publicKey string, gc *sharedconfig.GlobalConfig) (accessCount uint) {
	if publicKey == "" {
		return 0
	}
	wallet, err := userModels.UserWalletID(publicKey).GetWallet(gc.DB)
	if err != nil {
		return 0
	}

	if wallet.SharedAccessEnabled == 0 {
		return 0
	}
	for _, access := range wallet.UserWalletSharedAccess.AccessList {
		if access.AccessLevel == "APPROVER" {
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

	if wallet.SharedAccessEnabled == 0 {
		return 0
	}
	for _, access := range wallet.UserWalletSharedAccess.AccessList {
		if access.AccessLevel == "INITIATOR" {
			accessCount++
		}
	}

	return
}

func CreateSharedWalletAccess(signerPublicKey string, accessInfo *userModels.UserWalletSharedAccessInfo, gc *sharedconfig.GlobalConfig) (managedAccess userModels.UserWalletSharedAccess, err error) {
	// var managedAccess userModels.UserWalletSharedAccess

	if len(accessInfo.AccessList) == 0 {
		return managedAccess, &tErrors.CustomError{
			Param:      "accessList",
			Err:        "error-access-list-is-empty",
			ErrMessage: "Access list is empty",
			Code:       http.StatusBadRequest,
		}
	}
	var accessListInfo []userModels.WalletAccessInfo
	var accessList []userModels.WalletAccess
	var numberOfSubmittedApprovers int
	var numberOfSubmittedInitiators int
	var selfApprover int
	var approverUsers []*userModels.User
	uuidSharedAccess, _ := uuid.NewV4()
	sharedAccessID := uuidSharedAccess.String()

	wallet, e := userModels.UserWalletID(accessInfo.WalletPublicKey).GetWallet(gc.DB)
	if e != nil {
		log.Println("invalid wallet address: ", accessInfo.WalletPublicKey)
		return managedAccess, &tErrors.ErrorInvalidWallet{
			PublicKey: accessInfo.WalletPublicKey,
		}
	}
	if wallet.SharedAccessEnabled == 1 {
		return managedAccess, &tErrors.CustomError{
			Param:      "id",
			Err:        "error-shared-access-already-active-on-wallet",
			ErrMessage: "Shared access already activated on wallet. Use option to update the shared access.",
			Code:       http.StatusForbidden,
		}
	}

	if wallet.PrimaryWallet == 1 {
		if !PublicKeyHasViewOnlyAccessWACL(accessInfo.WalletPublicKey, accessInfo.AccessList, gc) {
			return managedAccess, &tErrors.ErrorOnlyViewAccessAllowedInPrimaryWallet{}
		}
	}
	e = gc.DB.Where("user_wallet_id = ?", accessInfo.WalletPublicKey).First(&managedAccess).Error
	if e != nil {
		if !errors.Is(e, gorm.ErrRecordNotFound) {
			return managedAccess, &tErrors.ErrorTemporaryServerError{}
		}

		walletOwner, e := userModels.UserWalletID(accessInfo.WalletPublicKey).GetWalletOwner(gc.DB)
		if e != nil {
			return managedAccess, &tErrors.CustomError{
				Param:      "username",
				Err:        "error-confirming-wallet-owner",
				ErrMessage: "Unable to confirm wallet owner at this time. Please try again after some minutes.",
				Code:       http.StatusForbidden,
			}
		}

		for _, v := range accessInfo.AccessList {
			uuidAccess, _ := uuid.NewV4()
			walletAccessID := uuidAccess.String()
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
			name := fmt.Sprintf("%v", u.FirstName)
			if u.LastName == nil {
				name = fmt.Sprintf("%v %v", name, *u.LastName)
			}
			//infor of shared access users
			accessListInfo = append(accessListInfo, userModels.WalletAccessInfo{
				UserWalletSharedAccessID: sharedAccessID,
				Username:                 v.Username,
				Name:                     name,
				AccessLevel:              v.AccessLevel,
				WalletPublicKey:          wallet.ID,
			})
			accessList = append(accessList, userModels.WalletAccess{
				ID:                       walletAccessID,
				UserWalletSharedAccessID: sharedAccessID,
				TargetUsername:           v.Username,
				AccessLevel:              v.AccessLevel,
				WalletPublicKey:          wallet.ID,
			})
			if v.AccessLevel == "APPROVER" {
				numberOfSubmittedApprovers++
				if v.Username == walletOwner.Username {
					selfApprover = 1
				}
				approverUsers = append(approverUsers, &u)
			}
			if v.AccessLevel == "INITIATOR" {
				numberOfSubmittedInitiators++
				// if v.Username == walletOwner.Username {
				// 	selfApprover = 1
				// }
			}

		}
		if numberOfSubmittedApprovers <= accessInfo.NumberOfApprovers && accessInfo.NumberOfApprovers > 1 {
			//number of authorizers does not reach the minimum threshold needed. cannot proceed so as to prevent account lockout
			return managedAccess, &tErrors.CustomError{
				Param:      "numberOfApprovers",
				Err:        "error-approvers-not-enough",
				ErrMessage: fmt.Sprintf("Please ensure that your list of approvers are greater than the minimum number required to approve a transaction [%v]", accessInfo.NumberOfApprovers),
				Code:       http.StatusForbidden,
			}
		}
		if numberOfSubmittedApprovers > 0 && accessInfo.NumberOfApprovers == 0 {
			//number of APPROVERS does not reach the minimum threshold needed. cannot proceed so as to prevent account lockout
			return managedAccess, &tErrors.CustomError{
				Param:      "numberOfApprovers",
				Err:        "error-approvers-not-enough",
				ErrMessage: "You must specify the number of approvers needed to approve transactions on this wallet",
				Code:       http.StatusForbidden,
			}
		}
		if numberOfSubmittedApprovers > 0 && numberOfSubmittedInitiators == 0 {
			return managedAccess, &tErrors.CustomError{
				Param:      "numberOfApprovers",
				Err:        "error-initiator-missing",
				ErrMessage: "You must specify at least one initiator when an approver is specified.",
				Code:       http.StatusForbidden,
			}
		}
		if numberOfSubmittedApprovers == 0 && numberOfSubmittedInitiators > 0 {
			return managedAccess, &tErrors.CustomError{
				Param:      "numberOfApprovers",
				Err:        "error-initiator-missing",
				ErrMessage: fmt.Sprintf("You must specify at least [%v] approvers when an approver is specified.", accessInfo.NumberOfApprovers+1),
				Code:       http.StatusForbidden,
			}
		}

		managedAccess = userModels.UserWalletSharedAccess{
			ID:                sharedAccessID,
			UserWalletID:      accessInfo.WalletPublicKey,
			NumberOfApprovers: accessInfo.NumberOfApprovers,
			// AccessList:        accessList,
		}
		accessInfo.AccessList = accessListInfo
		// if approver exists, then owner must sign transaction to add them as signers
		if (numberOfSubmittedApprovers - selfApprover) > 0 {
			accessInfo.SignatureRequired = 1

		}
		//prepare database execution
		dbTX := gc.DB.Begin()
		defer dbTX.Rollback()
		wallet.SharedAccessEnabled = 1
		errDB := dbTX.Save(&wallet).Error
		if err != nil {
			log.Printf("[CreateSharedWalletAccess] error saving shared access status of the wallet:%v\n sharedAccess:%+v\n", errDB, managedAccess)
			return managedAccess, &tErrors.ErrorTemporaryServerError{}
		}
		errDB = dbTX.Create(&managedAccess).Error
		if err != nil {
			log.Printf("[CreateSharedWalletAccess] error saving shared access:%v\n sharedAccess:%+v\n", errDB, managedAccess)
			return managedAccess, &tErrors.ErrorTemporaryServerError{}
		}
		errDB = dbTX.Create(&accessList).Error
		if err != nil {
			log.Printf("[CreateSharedWalletAccess] error saving access list:%v\n AccessList:%+v\n", errDB, accessList)
			return managedAccess, &tErrors.ErrorTemporaryServerError{}
		}
		xdrBase64, messages, walletMustSign, errGenXdr := generateCreateSharedAccessXdr(&wallet, &walletOwner, approverUsers, accessInfo.AccessList, accessInfo.NumberOfApprovers, gc)
		if errGenXdr != nil {
			return managedAccess, errGenXdr
		}
		accessInfo.Messages = messages
		if walletMustSign {
			accessInfo.SignatureRequired = 1

		}
		accessInfo.NetworkPassPhrase = network.GetBlockchainNetworkPassPhrase()

		accessInfo.Transaction = xdrBase64

		if len(accessInfo.TransactionSignature) == 0 {
			return managedAccess, nil
		}
		// extract signature and submit transaction
		//submit to blockchain

		txnHash, err := network.SubmitXdrWithSignature(gc.BantuExpansionClient, wallet.Signer, xdrBase64, accessInfo.TransactionSignature)
		if err != nil {
			log.Printf("Error submitting shared access txn [%+v] transaction: %s\n", accessInfo, err.Error())
			// logDiscordFailedRecovery(fmt.Sprintf("Error submitting shared access txn [%+v] transaction: %s", accessInfo, err.Error()))
			return managedAccess, &tErrors.ErrorTemporaryServerError{}
		}
		accessInfo.TransactionID = txnHash

		dbTX.Commit()
		return managedAccess, nil

	}
	//existing access was retrieved
	log.Printf("[CreateSharedWalletAccess] existing access: %+v", managedAccess)
	return managedAccess, &tErrors.CustomError{
		Param:      "id",
		Err:        "error-shared-access-already-active-on-wallet",
		ErrMessage: "Shared access already activated on wallet. Use option to update the shared access.",
		Code:       http.StatusForbidden,
	}
}
func RemoveSharedWalletAccess(signerPublicKey string, accessInfo *userModels.DisableSharedAccessInfo, gc *sharedconfig.GlobalConfig) (err error) {
	var managedAccess userModels.UserWalletSharedAccess

	accessList := make([]userModels.WalletAccess, 0)
	var numberOfApprovers int
	// var numberOfSubmittedInitiators int
	var selfApprover int
	var approverUsers []*userModels.User
	dbTX := gc.DB.Begin()
	defer dbTX.Rollback()
	wallet, e := userModels.UserWalletID(accessInfo.WalletPublicKey).GetWallet(dbTX)
	if e != nil {
		log.Println("invalid wallet address: ", accessInfo.WalletPublicKey)
		return &tErrors.ErrorInvalidWallet{
			PublicKey: accessInfo.WalletPublicKey,
		}
	}
	if wallet.SharedAccessEnabled == 0 {
		return &tErrors.CustomError{
			Param:      "id",
			Err:        "error-shared-access-not-active-on-wallet",
			ErrMessage: "Shared access not activated on wallet.",
			Code:       http.StatusForbidden,
		}
	}
	wallet.SharedAccessEnabled = 0
	dbTX.Save(&wallet)
	e = dbTX.Where("user_wallet_id = ?", accessInfo.WalletPublicKey).First(&managedAccess).Error

	if e != nil {
		if errors.Is(e, gorm.ErrRecordNotFound) {
			return &tErrors.CustomError{
				Param:      "Id",
				Err:        "error-could-not-find-shared-access",
				ErrMessage: "Could not find shared access records for wallet.",
				Code:       http.StatusForbidden,
			}
		}
		return &tErrors.ErrorTemporaryServerError{}
	}

	e = dbTX.Where("wallet_public_key = ?", accessInfo.WalletPublicKey).Find(&accessList).Error

	if e != nil {
		if errors.Is(e, gorm.ErrRecordNotFound) {
			return &tErrors.CustomError{
				Param:      "Id",
				Err:        "error-could-not-find-access-list",
				ErrMessage: "Could not find access list for wallet.",
				Code:       http.StatusForbidden,
			}
		}
		return &tErrors.ErrorTemporaryServerError{}
	}

	walletOwner, e := userModels.UserWalletID(accessInfo.WalletPublicKey).GetWalletOwner(dbTX)
	if e != nil {
		return &tErrors.CustomError{
			Param:      "username",
			Err:        "error-confirming-wallet-owner",
			ErrMessage: "Unable to confirm wallet owner at this time. Please try again after some minutes.",
			Code:       http.StatusForbidden,
		}
	}

	for _, v := range accessList {
		u, e := usersDB.GetUser(v.TargetUsername, dbTX)
		if e != nil {
			return &tErrors.CustomError{
				Param:      "username",
				Err:        "error-trovo-wallet-account-invalid",
				ErrMessage: fmt.Sprintf("Trovo wallet account [%v] could not be validated at this time.", v.TargetUsername),
				Code:       http.StatusForbidden,
			}
		}
		if v.AccessLevel == "APPROVER" {
			numberOfApprovers++
			if v.TargetUsername == walletOwner.Username {
				selfApprover = 1
			}
			approverUsers = append(approverUsers, &u)
		}

	}

	e = dbTX.Delete(&accessList).Error
	if e != nil {
		log.Println("error deleting access list", e)
		return &tErrors.ErrorTemporaryServerError{}
	}

	// if approver exists, then owner must sign transaction to add them as signers
	if (numberOfApprovers - selfApprover) > 0 {
		accessInfo.SignatureRequired = 1

	}

	errDB := dbTX.Delete(&managedAccess).Error
	if err != nil {
		log.Printf("[CreateSharedWalletAccess] error saving shared access:%v\n sharedAccess:%+v\n", errDB, managedAccess)
		return &tErrors.ErrorTemporaryServerError{}
	}

	xdrBase64, messages, walletMustSign, errGenXdr := generateRemoveSharedAccessXdr(&wallet, &walletOwner, approverUsers, gc)
	if errGenXdr != nil {
		return errGenXdr
	}
	accessInfo.Messages = messages
	if walletMustSign {
		accessInfo.SignatureRequired = 1

	}
	accessInfo.NetworkPassPhrase = network.GetBlockchainNetworkPassPhrase()

	accessInfo.Transaction = xdrBase64

	if len(accessInfo.TransactionSignature) == 0 {
		return nil
	}
	// extract signature and submit transaction
	//submit to blockchain

	txnHash, err := network.SubmitXdrWithSignature(gc.BantuExpansionClient, wallet.Signer, xdrBase64, accessInfo.TransactionSignature)
	if err != nil {
		log.Printf("Error submitting shared access txn [%+v] transaction: %s\n", accessInfo, err.Error())
		// logDiscordFailedRecovery(fmt.Sprintf("Error submitting shared access txn [%+v] transaction: %s", accessInfo, err.Error()))
		return &tErrors.ErrorTemporaryServerError{}
	}
	accessInfo.TransactionID = txnHash

	dbTX.Commit()
	return nil

}

func generateCreateSharedAccessXdr(wallet *userModels.UserWallet, walletOwner *userModels.User, approvers []*userModels.User, accessInfo []userModels.WalletAccessInfo, authThreshold int, gc *sharedconfig.GlobalConfig) (xdrbase64 string, messages []string, walletMustSign bool, err error) {
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

	if len(approvers) == 0 && authThreshold > 0 {
		err = &tErrors.CustomError{
			Param:      "numberOfApprovers",
			Err:        "error-no-approver-specified",
			ErrMessage: "No approvers specifieds",
			Code:       404,
		}
		return "", messages, walletMustSign, err
	}

	//check if primary account has native enough native balance
	var nativeAsset txnbuild.Asset = txnbuild.NativeAsset{}
	walletAccountExists, _, walletAccountNativeBalance, _, walletSourceAccount, errWalletAct := network.BlockchainAccountProperties(client, wallet.ID, nativeAsset)
	if errWalletAct != nil {
		log.Printf("[generateCreateSharedAccessXdr] by [%v] for shared Account Properties error:[%v] \n", wallet.Alias, errWalletAct)

		return "", messages, walletMustSign, errWalletAct
	}
	if !walletAccountExists || (walletAccountNativeBalance.Sub(activationAmount)).LessThan(minBalance) {
		log.Printf("[generateCreateSharedAccessXdr] by [%v] shared WalletAccount underfunded \n", wallet.Alias)

		err = &tErrors.CustomError{
			Param:      "publicKey",
			Err:        "error-wallet-underfunded",
			ErrMessage: fmt.Sprintf("Wallet %v does not have enough XBN balance to perform this operation", wallet.Alias),
			Code:       404,
		}
		return "", messages, walletMustSign, err
	}
	{
		//check if account recovery is enabled, then disable it on the wallet.
		if walletOwner.AccountRecoveryEnabled == 1 && !PublicKeyHasViewOnlyAccessWACL(wallet.ID, accessInfo, gc) && PublicKeyCountApproverAccessWACL(wallet.ID, accessInfo, gc) > 1 {
			// get the recovery keypair
			recoveryAddress := bc.GetRecoveryAccountAddress(walletOwner.Username, walletOwner.PublicKey)

			if userBc.SignerIsValid(wallet.ID, recoveryAddress) {
				//recovery a signer to the wallet. remove it
				ops = append(ops, &txnbuild.SetOptions{
					Signer: &txnbuild.Signer{
						Address: recoveryAddress,
						Weight:  0,
					},
					SourceAccount: wallet.ID,
				})

				//add message about disabling recovery on that wallet
				messages = append(messages, "Account Recovery on this wallet has to be disabled so as to enable shared access.")
				walletMustSign = true
			}
		}
	}

	//check access list to know if you would activate the user wallets before proceeding.
	for _, user3p := range approvers {
		var totalUsersToFund int64
		//ensure u r using the account signer, since the account may have been recovered, or may be recovered in the future, changing the signer, but retaining the primary key
		approverAccountExists, _, _, _, _, _ := network.BlockchainAccountProperties(client, user3p.PrimarySigner, nativeAsset)
		if !approverAccountExists {
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

		if approverAccountExists {
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
				log.Printf("[generateCreateSharedAccessXdr] by [%v] MultiAccess WalletAccount underfunded \n", wallet.Alias)

				err = &tErrors.CustomError{
					Param:      "walletPublicKey",
					Err:        "error-wallet-underfunded",
					ErrMessage: fmt.Sprintf("Wallet %v needs more than %v XBN balance to perform this operation", wallet.Alias, totalNativeBalanceNeeded.String()),
					Code:       404,
				}
				return "", messages, walletMustSign, err
			}
		}
	}

	//TODO: if account exists and subwallet has enough balance, we add the operation to pay TROVO fee from primary Wallet
	{
		//process service fee
		if len(os.Getenv("SHARED_ACCESS_FEE_ASSET_ISSUER")) == 0 {
			ops = append(ops, &txnbuild.Payment{
				Destination:   os.Getenv("SHARED_ACCESS_FEE_ADDRESS"),
				Amount:        os.Getenv("SHARED_ACCESS_FEE_AMOUNT"),
				SourceAccount: wallet.ID,
				Asset:         txnbuild.NativeAsset{},
			})
			messages = append(messages, fmt.Sprintf("%v %v will be deducted from wallet %v as service fee.", os.Getenv("SHARED_ACCESS_FEE_AMOUNT"), os.Getenv("NATIVE_ASSET_CODE"), wallet.Alias))

		} else {
			ops = append(ops, &txnbuild.Payment{
				Destination:   os.Getenv("SHARED_ACCESS_FEE_ADDRESS"),
				Amount:        os.Getenv("SHARED_ACCESS_FEE_AMOUNT"),
				SourceAccount: wallet.ID,
				Asset:         txnbuild.CreditAsset{Code: os.Getenv("SHARED_ACCESS_FEE_ASSET_CODE"), Issuer: os.Getenv("SHARED_ACCESS_FEE_ASSET_ISSUER")},
			})
			messages = append(messages, fmt.Sprintf("%v %v will be deducted from wallet %v as service fee.", os.Getenv("SHARED_ACCESS_FEE_AMOUNT"), os.Getenv("SHARED_ACCESS_FEE_ASSET_CODE"), wallet.Alias))

		}

		walletMustSign = true
	}
	// Construct the transaction that holds the operations to execute on the network
	{
		//adjust account threshold
		if authThreshold > 0 {
			ops = append(ops, &txnbuild.SetOptions{
				LowThreshold:    txnbuild.NewThreshold(txnbuild.Threshold(authThreshold)),
				MediumThreshold: txnbuild.NewThreshold(txnbuild.Threshold(authThreshold)),
				HighThreshold:   txnbuild.NewThreshold(txnbuild.Threshold(authThreshold)),
				SourceAccount:   wallet.ID,
			})
			walletMustSign = true
		}
	}

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
			Memo: txnbuild.MemoText("Create shared access"),
		},
	)
	if err != nil {
		log.Println("[generateCreateSharedAccessXdr] error constructing transaction ", err)
		return "", messages, walletMustSign, err
	}

	var xdrBase64 string

	xdrBase64, err = tx.Base64()
	if err != nil {
		log.Println("[generateCreateSharedAccessXdr] error getting txn base64", err)
		return "", messages, walletMustSign, err
	}

	return xdrBase64, messages, walletMustSign, nil

}

func generateRemoveSharedAccessXdr(wallet *userModels.UserWallet, walletOwner *userModels.User, approvers []*userModels.User, gc *sharedconfig.GlobalConfig) (xdrbase64 string, messages []string, walletMustSign bool, err error) {
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

	//check if primary account has native enough native balance
	var nativeAsset txnbuild.Asset = txnbuild.NativeAsset{}
	walletAccountExists, _, walletAccountNativeBalance, _, walletSourceAccount, errWalletAct := network.BlockchainAccountProperties(client, wallet.ID, nativeAsset)
	if errWalletAct != nil {
		log.Printf("[generateRemoveSharedAccessXdr] by [%v] for shared Account Properties error:[%v] \n", wallet.Alias, errWalletAct)

		return "", messages, walletMustSign, errWalletAct
	}
	if !walletAccountExists || (walletAccountNativeBalance.Sub(activationAmount)).LessThan(minBalance) {
		log.Printf("[generateRemoveSharedAccessXdr] by [%v] shared WalletAccount underfunded \n", wallet.Alias)

		err = &tErrors.CustomError{
			Param:      "publicKey",
			Err:        "error-wallet-underfunded",
			ErrMessage: fmt.Sprintf("Wallet %v does not have enough XBN balance to perform this operation", wallet.Alias),
			Code:       404,
		}
		return "", messages, walletMustSign, err
	}
	{
		//check if account recovery is enabled, then re-enable it on the wallet.
		if walletOwner.AccountRecoveryEnabled == 1 && PublicKeyHasViewOnlyAccess(wallet.ID, gc) && PublicKeyCountApproverAccess(wallet.ID, gc) > 1 {
			// get the recovery keypair
			recoveryAddress := bc.GetRecoveryAccountAddress(walletOwner.Username, walletOwner.PublicKey)

			if !userBc.SignerIsValid(wallet.ID, recoveryAddress) {
				//recovery a signer to the wallet. remove it
				ops = append(ops, &txnbuild.SetOptions{
					Signer: &txnbuild.Signer{
						Address: recoveryAddress,
						Weight:  1,
					},
					SourceAccount: wallet.ID,
				})

				//add message about disabling recovery on that wallet
				messages = append(messages, "Account Recovery on this wallet has been enabled.")
				walletMustSign = true
			}
		}
	}

	//check access list to know if you would activate the user wallets before proceeding.
	for _, user3p := range approvers {
		var totalUsersToFund int64
		//ensure u r using the account signer, since the account may have been recovered, or may be recovered in the future, changing the signer, but retaining the primary key
		approverAccountExists, _, _, _, _, _ := network.BlockchainAccountProperties(client, user3p.PrimarySigner, nativeAsset)

		if approverAccountExists {
			//account exists, check if it already it a signer in the wallet

			//after topping up, it now has enough balance to add primary wallet as signer if it is not already a signer
			if !wallet.SignerIsValidWA(user3p.PrimarySigner, walletSourceAccount) {

				ops = append(ops, &txnbuild.SetOptions{
					Signer: &txnbuild.Signer{
						Address: user3p.PrimarySigner,
						Weight:  0,
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
				log.Printf("[generateRemoveSharedAccessXdr] by [%v] MultiAccess WalletAccount underfunded \n", wallet.Alias)

				err = &tErrors.CustomError{
					Param:      "walletPublicKey",
					Err:        "error-wallet-underfunded",
					ErrMessage: fmt.Sprintf("Wallet %v needs more than %v XBN balance to perform this operation", wallet.Alias, totalNativeBalanceNeeded.String()),
					Code:       404,
				}
				return "", messages, walletMustSign, err
			}
		}
	}

	//TODO: if account exists and subwallet has enough balance, we add the operation to pay TROVO fee from primary Wallet
	{
		//process service fee
		if len(os.Getenv("SHARED_ACCESS_FEE_ASSET_ISSUER")) == 0 {
			ops = append(ops, &txnbuild.Payment{
				Destination:   os.Getenv("SHARED_ACCESS_FEE_ADDRESS"),
				Amount:        os.Getenv("SHARED_ACCESS_FEE_AMOUNT"),
				SourceAccount: wallet.ID,
				Asset:         txnbuild.NativeAsset{},
			})
			messages = append(messages, fmt.Sprintf("%v %v will be deducted from wallet %v as service fee.", os.Getenv("SHARED_ACCESS_FEE_AMOUNT"), os.Getenv("NATIVE_ASSET_CODE"), wallet.Alias))

		} else {
			ops = append(ops, &txnbuild.Payment{
				Destination:   os.Getenv("SHARED_ACCESS_FEE_ADDRESS"),
				Amount:        os.Getenv("SHARED_ACCESS_FEE_AMOUNT"),
				SourceAccount: wallet.ID,
				Asset:         txnbuild.CreditAsset{Code: os.Getenv("SHARED_ACCESS_FEE_ASSET_CODE"), Issuer: os.Getenv("SHARED_ACCESS_FEE_ASSET_ISSUER")},
			})
			messages = append(messages, fmt.Sprintf("%v %v will be deducted from wallet %v as service fee.", os.Getenv("SHARED_ACCESS_FEE_AMOUNT"), os.Getenv("SHARED_ACCESS_FEE_ASSET_CODE"), wallet.Alias))

		}

		walletMustSign = true
	}
	// Construct the transaction that holds the operations to execute on the network
	{
		//adjust account threshold

		ops = append(ops, &txnbuild.SetOptions{
			LowThreshold:    txnbuild.NewThreshold(txnbuild.Threshold(0)),
			MediumThreshold: txnbuild.NewThreshold(txnbuild.Threshold(0)),
			HighThreshold:   txnbuild.NewThreshold(txnbuild.Threshold(0)),
			SourceAccount:   wallet.ID,
		})
		walletMustSign = true

	}

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
			Memo: txnbuild.MemoText("Disable shared access"),
		},
	)
	if err != nil {
		log.Println("[generateRemoveSharedAccessXdr] error constructing transaction ", err)
		return "", messages, walletMustSign, err
	}

	var xdrBase64 string

	xdrBase64, err = tx.Base64()
	if err != nil {
		log.Println("[generateRemoveSharedAccessXdr] error getting txn base64", err)
		return "", messages, walletMustSign, err
	}

	return xdrBase64, messages, walletMustSign, nil

}
