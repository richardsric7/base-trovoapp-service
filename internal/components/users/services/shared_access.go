package users

import (
	"encoding/json"
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

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/txnbuild"
	"gorm.io/gorm/clause"
)

func WalletCountViewOnlyAccess(wallet *userModels.UserWallet, gc *sharedconfig.GlobalConfig) (accessCount uint) {
	if wallet.SharedAccessEnabled == 0 {
		return 0
	}
	for _, access := range wallet.Permissions {
		if access.Permission == "VIEW-ONLY" {
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
	for _, access := range wallet.Permissions {
		if access.Permission != "VIEW-ONLY" {
			return false
		}
	}

	return
}

func WalletCountApproverAccess(wallet *userModels.UserWallet, gc *sharedconfig.GlobalConfig) (accessCount uint) {
	if wallet.SharedAccessEnabled == 0 {
		return 0
	}
	for _, access := range wallet.Permissions {
		if access.Permission == "APPROVER" {
			accessCount++
		}
	}

	return
}

func WalletCountInitiatorAccess(wallet *userModels.UserWallet, gc *sharedconfig.GlobalConfig) (accessCount uint) {
	if wallet.SharedAccessEnabled == 0 {
		return 0
	}
	for _, access := range wallet.Permissions {
		if access.Permission == "INITIATOR" {
			accessCount++
		}
	}

	return
}

func PublicKeyCountViewOnlyAccess(publicKey string, gc *sharedconfig.GlobalConfig) (accessCount uint) {
	if publicKey == "" {
		return 0
	}
	accessList := userModels.UserWalletID(publicKey).GetPermissionList(gc.DB)
	if len(accessList) == 0 {
		return 0
	}

	for _, access := range accessList {
		if access.Permission == "VIEW-ONLY" {
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
	accessList := userModels.UserWalletID(publicKey).GetPermissionList(gc.DB)
	if len(accessList) == 0 {
		return false
	}

	for _, access := range accessList {
		if access.Permission != "VIEW-ONLY" {
			return false
		}
	}

	return
}

func PublicKeyHasViewOnlyAccessWACL(publicKey string, accessList []userModels.WalletPermissionInfo, gc *sharedconfig.GlobalConfig) (viewOnly bool) {
	viewOnly = true
	if publicKey == "" {
		return false
	}

	for _, access := range accessList {
		if access.Permission != "VIEW-ONLY" {
			return false
		}
	}

	return
}

func PublicKeyCountApproverAccessWACL(publicKey string, accessList []userModels.WalletPermissionInfo, gc *sharedconfig.GlobalConfig) (accessCount uint) {
	if publicKey == "" {
		return 0
	}

	for _, access := range accessList {
		if access.Permission == "APPROVER" {
			accessCount++
		}
	}

	return
}

func PublicKeyCountInitiatorAccessWACL(publicKey string, accessList []userModels.WalletPermissionInfo, gc *sharedconfig.GlobalConfig) (accessCount uint) {
	if publicKey == "" {
		return 0
	}

	for _, access := range accessList {
		if access.Permission == "INITIATOR" {
			accessCount++
		}
	}

	return
}

func PublicKeyCountApproverAccess(publicKey string, gc *sharedconfig.GlobalConfig) (accessCount uint) {
	if publicKey == "" {
		return 0
	}
	acl := userModels.UserWalletID(publicKey).GetPermissionList(gc.DB)
	if len(acl) == 0 {
		return 0
	}

	for _, access := range acl {
		if access.Permission == "APPROVER" {
			accessCount++
		}
	}

	return
}

func PublicKeyCountInitiatorAccess(publicKey string, gc *sharedconfig.GlobalConfig) (accessCount uint) {
	if publicKey == "" {
		return 0
	}
	acl := userModels.UserWalletID(publicKey).GetPermissionList(gc.DB)
	if len(acl) == 0 {
		return 0
	}

	for _, access := range acl {
		if access.Permission == "INITIATOR" {
			accessCount++
		}
	}

	return
}

func CreateSharedWalletAccess(signerUser *userModels.User, walletOwner *userModels.User, wallet *userModels.UserWallet, accessInfo *userModels.UserWalletSharedAccessInfo, gc *sharedconfig.GlobalConfig) (returnedWallet userModels.UserWallet, err error) {

	// var  userModels.UserWalletSharedAccess
	accessInfo.Messages = make([]string, 0)
	if len(accessInfo.Permissions) == 0 {
		return returnedWallet, &tErrors.CustomError{
			Param:      "permissions",
			Err:        "error-permission-list-is-empty",
			ErrMessage: "Permission list is empty",
			Code:       http.StatusBadRequest,
		}
	}

	if w, v := wallet.IsValidLinkedWallet(gc); v {
		return returnedWallet, &tErrors.CustomError{
			Param:      "publicKey",
			Err:        "error-linked-wallets-not-allowed",
			ErrMessage: fmt.Sprintf("Linked Wallets are not allowed to be shared directly. Plase share the access on %v and it will mirror to this wallet.", w.Alias),
			Code:       http.StatusBadRequest,
		}
	}
	var linkedWallet userModels.UserWallet
	var hasLinkedWallet bool
	if wallet.WalletType == 1 && wallet.LinkedWalletPublicKey != nil {
		// set the linked wallet if it is a tokenization wallet
		hasLinkedWallet = true
		// accessInfo.LinkedWalletSignatureRequired = 1
		// accessInfo.LinkedWalletPublicKey = *wallet.LinkedWalletPublicKey
		linkedWallet, err = userModels.UserWalletID(*wallet.LinkedWalletPublicKey).GetWallet(gc.DB, gc)
		if err != nil {
			return returnedWallet, &tErrors.CustomError{
				Param:      "linkedWalletPublicKey",
				Err:        "error-getting-linked-wallet",
				ErrMessage: "Linked Wallet could not be validated.",
				Code:       http.StatusBadRequest,
			}
		}

	}
	if len(accessInfo.Permissions) == 1 && accessInfo.Permissions[0].TargetUsername == walletOwner.Username {
		return returnedWallet, &tErrors.CustomError{
			Param:      "permissions",
			Err:        "error-permission-unacceptable-access",
			ErrMessage: "Shared access cannot be enabled with view-only permission granted to yourself.",
			Code:       http.StatusForbidden,
		}
	}

	if wallet.WalletType == 2 || wallet.WalletType == 3 && accessInfo.NumberOfApprovalsNeeded > 0 {
		return returnedWallet, &tErrors.CustomError{
			Param:      "publicKey",
			Err:        "error-wallet-type-not-allowed",
			ErrMessage: "Market Making & Bulk Payment wallets allows only view access to be enabled.",
			Code:       http.StatusForbidden,
		}
	}

	var accessListInfo, viewOnly []userModels.WalletPermissionInfo
	var accessList, linkedWalletAccessList []userModels.WalletPermission
	var numberOfSubmittedApprovers int
	var numberOfSubmittedInitiators int
	var selfApprover int
	var approverUsers []*userModels.User

	if wallet.SharedAccessEnabled == 1 {
		return returnedWallet, &tErrors.CustomError{
			Param:      "id",
			Err:        "error-shared-access-already-active-on-wallet",
			ErrMessage: "Shared access already activated on wallet. Use option to update the shared access.",
			Code:       http.StatusForbidden,
		}
	}

	if wallet.PrimaryWallet == 1 {
		if !PublicKeyHasViewOnlyAccessWACL(accessInfo.WalletPublicKey, accessInfo.Permissions, gc) {
			return returnedWallet, &tErrors.ErrorOnlyViewAccessAllowedInPrimaryWallet{}
		}
	}
	checkAccess := make(map[string]userModels.WalletPermissionInfo, 0)
	// uniquePermission := make([]userModels.WalletPermissionInfo, 0)

	for _, v := range accessInfo.Permissions {

		if v.TargetUsername == walletOwner.Username && v.Permission == "VIEW-ONLY" {
			//skip owners being added as view only access. Owners have view access by default.
			continue
		}
		v.TargetUsername = strings.ToLower(v.TargetUsername)
		if _, ok := checkAccess[v.TargetUsername+v.Permission]; ok {
			continue
		}
		if v.Permission == "VIEW-ONLY" {
			viewOnly = append(viewOnly, v)
		}

		// checkAccess[v.TargetUsername+v.Permission] = v
		// uniquePermission = append(uniquePermission, v)

		permissionID := uuid.NewString()

		//check if username is valid
		if strings.Contains(v.TargetUsername, "_") {

			//it is a subwallet and cannot be given access
			return returnedWallet, &tErrors.CustomError{
				Param:      "username",
				Err:        "error-subwallet-not-allowed",
				ErrMessage: fmt.Sprintf("Access can only be granted to trovo wallet account, not a subwallet [%v]", v.TargetUsername),
				Code:       http.StatusForbidden,
			}
		}
		u, e := userModels.Username(v.TargetUsername).GetSimpleUser(gc.DB, gc)
		if e != nil {
			return returnedWallet, &tErrors.CustomError{
				Param:      "username",
				Err:        "error-trovo-wallet-account-invalid",
				ErrMessage: fmt.Sprintf("Trovo wallet account [%v] could not be validated at this time.", v.TargetUsername),
				Code:       http.StatusForbidden,
			}
		}
		if u.Username != v.TargetUsername {

			return returnedWallet, &tErrors.CustomError{
				Param:      "username",
				Err:        "error-trovo-wallet-account-invalid",
				ErrMessage: fmt.Sprintf("Trovo wallet account [%v] is not Trovo wallet account username.", v.TargetUsername),
				Code:       http.StatusForbidden,
			}

		}
		name := fmt.Sprintf("%v", u.FirstName)
		if u.LastName != nil {
			name = fmt.Sprintf("%v %v", name, *u.LastName)
		}
		pi := userModels.WalletPermissionInfo{
			TargetUsername:        v.TargetUsername,
			Name:                  name,
			Permission:            v.Permission,
			WalletPublicKey:       wallet.ID,
			WalletAlias:           wallet.Alias,
			PushNotificationToken: u.PushNotificationToken,
		}
		//infor of shared access users
		accessListInfo = append(accessListInfo, pi)
		accessList = append(accessList, userModels.WalletPermission{
			ID:              permissionID,
			TargetUsername:  v.TargetUsername,
			Permission:      v.Permission,
			WalletPublicKey: wallet.ID,
		})
		if v.Permission == "APPROVER" {
			numberOfSubmittedApprovers++
			if v.TargetUsername == walletOwner.Username {
				selfApprover = 1
			}
			approverUsers = append(approverUsers, &u)
		}
		if v.Permission == "INITIATOR" {
			numberOfSubmittedInitiators++

		}
		checkAccess[v.TargetUsername+v.Permission] = pi

		{
			//if it is a tokenized wallet, then build linkedwallet access list
			if hasLinkedWallet {
				linkedWalletAccessList = append(linkedWalletAccessList, userModels.WalletPermission{
					ID:              uuid.NewString(),
					TargetUsername:  v.TargetUsername,
					Permission:      v.Permission,
					WalletPublicKey: linkedWallet.ID,
				})
			}
		}
	}
	displayMessage := false
	if len(viewOnly) > 0 {
		message := "Unnecessary VIEW-ONLY access for these accounts where removed:"
		for _, v := range viewOnly {
			for _, a := range accessInfo.Permissions {
				if v.TargetUsername == a.TargetUsername && (a.Permission == "APPROVER" || a.Permission == "INITIATOR") {
					// remove the view only since the approver and initiator has view access already
					displayMessage = true
					delete(checkAccess, v.TargetUsername+"VIEW-ONLY")
					message = fmt.Sprintf("%v,%v", message, v.TargetUsername)
				}
			}
		}
		// rebuild the list
		accessListInfo = make([]userModels.WalletPermissionInfo, 0)
		for _, ca := range checkAccess {
			accessListInfo = append(accessListInfo, ca)
		}
		if displayMessage {
			accessInfo.Messages = append(accessInfo.Messages, message)
		}

	}

	if numberOfSubmittedApprovers <= accessInfo.NumberOfApprovalsNeeded && accessInfo.NumberOfApprovalsNeeded > 1 {
		//number of authorizers does not reach the minimum threshold needed. cannot proceed so as to prevent account lockout
		return returnedWallet, &tErrors.CustomError{
			Param:      "numberOfApprovalsNeeded",
			Err:        "error-approvers-not-enough",
			ErrMessage: fmt.Sprintf("Please ensure that your list of approvers configured are greater than the minimum number required to approve a transaction [%v]", accessInfo.NumberOfApprovalsNeeded),
			Code:       http.StatusForbidden,
		}
	}
	if numberOfSubmittedApprovers > 0 && accessInfo.NumberOfApprovalsNeeded == 0 {
		//number of APPROVERS does not reach the minimum threshold needed. cannot proceed so as to prevent account lockout
		return returnedWallet, &tErrors.CustomError{
			Param:      "numberOfApprovalsNeeded",
			Err:        "error-approvers-not-enough",
			ErrMessage: "You must specify the number of approvers needed to approve transactions on this wallet",
			Code:       http.StatusForbidden,
		}
	}
	if numberOfSubmittedApprovers > 0 && numberOfSubmittedInitiators == 0 {
		return returnedWallet, &tErrors.CustomError{
			Param:      "numberOfApprovalsNeeded",
			Err:        "error-initiator-missing",
			ErrMessage: "You must specify at least one initiator when an approver is specified.",
			Code:       http.StatusForbidden,
		}
	}
	if numberOfSubmittedApprovers == 0 && numberOfSubmittedInitiators > 0 {
		return returnedWallet, &tErrors.CustomError{
			Param:      "numberOfApprovers",
			Err:        "error-initiator-missing",
			ErrMessage: fmt.Sprintf("You must specify at least [%v] approvers when an approver is specified.", accessInfo.NumberOfApprovalsNeeded+1),
			Code:       http.StatusForbidden,
		}
	}

	accessInfo.Permissions = accessListInfo
	// if approver exists, then owner must sign transaction to add them as signers
	if (numberOfSubmittedApprovers - selfApprover) > 0 {
		accessInfo.SignatureRequired = 1

	}
	//prepare database execution
	dbTX := gc.DB.Begin()
	defer dbTX.Rollback()
	wallet.SharedAccessEnabled = 1
	wallet.NumberOfApprovalsNeeded = accessInfo.NumberOfApprovalsNeeded

	errDB := dbTX.Create(&accessList).Error
	if err != nil {
		log.Printf("[CreateSharedWalletAccess] error saving access list:%v\n AccessList:%+v\n", errDB, accessList)
		return returnedWallet, &tErrors.ErrorTemporaryServerError{}
	}

	errDB = dbTX.Omit(clause.Associations).Save(wallet).Error
	if err != nil {
		log.Printf("[CreateSharedWalletAccess] error saving shared access status of the wallet:%v\n sharedAccess:%+v\n", errDB, wallet)
		return returnedWallet, &tErrors.ErrorTemporaryServerError{}
	}

	//if has linked wallet, clone the wallet property
	if hasLinkedWallet {
		linkedWallet.SharedAccessEnabled = wallet.SharedAccessEnabled
		linkedWallet.NumberOfApprovalsNeeded = wallet.NumberOfApprovalsNeeded
		errDB := dbTX.Create(&linkedWalletAccessList).Error
		if err != nil {
			log.Printf("[CreateSharedWalletAccess] error saving linked wallet access list:%v\n Linked wallet AccessList:%+v\n", errDB, linkedWalletAccessList)
			return returnedWallet, &tErrors.ErrorTemporaryServerError{}
		}

		errDB = dbTX.Omit(clause.Associations).Save(linkedWallet).Error
		if err != nil {
			log.Printf("[CreateSharedWalletAccess] error saving shared access status of the linked wallet:%v\n linkedsharedAccess:%+v\n", errDB, linkedWallet)
			return returnedWallet, &tErrors.ErrorTemporaryServerError{}
		}
	}

	xdrBase64, messages, walletMustSign, errGenXdr := generateCreateSharedAccessXdr(wallet, walletOwner, approverUsers, accessInfo.Permissions, accessInfo.NumberOfApprovalsNeeded, gc)
	if errGenXdr != nil {
		return returnedWallet, errGenXdr
	}
	accessInfo.Messages = append(accessInfo.Messages, messages...)
	if walletMustSign {
		accessInfo.SignatureRequired = 1

	}
	accessInfo.NetworkPassPhrase = network.GetBlockchainNetworkPassPhrase()

	accessInfo.Transaction = xdrBase64

	if len(accessInfo.TransactionSignature) == 0 {
		return returnedWallet, nil
	}
	// extract signature and submit transaction
	//submit to blockchain
	signatures := make(map[string]string, 0)
	signatures[wallet.Signer] = accessInfo.TransactionSignature
	// if len(accessInfo.LinkedWalletTransactionSignature) > 10 {
	// 	signatures[linkedWallet.Signer] = accessInfo.LinkedWalletTransactionSignature

	// }
	// txnHash, err := network.SubmitXdrWithSignature(gc.BantuExpansionClient, wallet.Signer, xdrBase64, accessInfo.TransactionSignature)
	txnHash, err := network.SubmitXdrWithSignatures(gc.BantuExpansionClient, xdrBase64, signatures, gc.DB)
	if err != nil {
		log.Printf("Error submitting shared access txn [%+v] transaction: %s\n", accessInfo, err.Error())
		// logDiscordFailedRecovery(fmt.Sprintf("Error submitting shared access txn [%+v] transaction: %s", accessInfo, err.Error()))
		return returnedWallet, &tErrors.ErrorTemporaryServerError{}
	}
	accessInfo.TransactionID = txnHash

	dbTX.Commit()
	// invalidate cache
	{

		for _, v := range accessList {
			u, e := userModels.Username(v.TargetUsername).GetSimpleUser(gc.DB, gc)
			if e == nil {
				u.InvalidateUserCache(gc)
			}
		}
		signerUser.InvalidateUserCache(gc)
	}
	return returnedWallet, nil
}

func ModifySharedWalletAccess(signerUser *userModels.User, walletOwner *userModels.User, wallet *userModels.UserWallet, accessInfo *userModels.ModifySharedAccessInfo, gc *sharedconfig.GlobalConfig) (revokedList, modifiedList, addedList []userModels.WalletPermission, linkedRevokedList, linkedModifiedList, linkedAddedList []userModels.WalletPermission, err error) {

	//prepare database execution
	dbTX := gc.DB.Begin()
	defer dbTX.Rollback()
	if wallet.WalletType == 2 || wallet.WalletType == 3 {
		log.Printf("[ModifySharedWalletAccess] error wallet type not allowed for %+v\n", accessInfo)
		err = &tErrors.CustomError{
			Param:      "publicKey",
			Err:        "error-wallet-type-not-allowed",
			ErrMessage: "Market Making & Bulk Payment wallets are not allowed for this operation.",
			Code:       http.StatusForbidden,
		}
		return
	}

	if w, v := wallet.IsValidLinkedWallet(gc); v {
		err = &tErrors.CustomError{
			Param:      "publicKey",
			Err:        "error-linked-wallets-not-allowed",
			ErrMessage: fmt.Sprintf("Linked Wallets are not allowed to be shared directly. Plase share the access on %v and it will mirror to this wallet.", w.Alias),
			Code:       http.StatusBadRequest,
		}
		return
	}

	var linkedWallet userModels.UserWallet
	var hasLinkedWallet bool
	if wallet.WalletType == 1 && wallet.LinkedWalletPublicKey != nil {
		// set the linked wallet if it is a tokenization wallet
		hasLinkedWallet = true

		linkedWallet, err = userModels.UserWalletID(*wallet.LinkedWalletPublicKey).GetWallet(gc.DB, gc)
		if err != nil {
			err = &tErrors.CustomError{
				Param:      "linkedWalletPublicKey",
				Err:        "error-getting-linked-wallet",
				ErrMessage: "Linked Wallet could not be validated.",
				Code:       http.StatusBadRequest,
			}
			return
		}

	}
	if len(accessInfo.ModifiedPermissions) == 0 && len(accessInfo.AddedPermissions) == 0 && len(accessInfo.RevokedPermissions) == 0 {
		log.Printf("[ModifySharedWalletAccess] error no operations for %+v\n", accessInfo)
		err = &tErrors.CustomError{
			Param:      "permissions",
			Err:        "error-no-operations-to-perform",
			ErrMessage: "No operation can be performed since no permissions to add or revoke or modify",
			Code:       http.StatusBadRequest,
		}
		return
	}
	viewOnly := make(map[string]string, 0)
	// oldApproverPermissionMap := make(map[string]userModels.WalletPermission, 0)
	walletID := userModels.UserWalletID(accessInfo.WalletPublicKey)
	var linkedWalletID userModels.UserWalletID
	if hasLinkedWallet {
		linkedWalletID = userModels.UserWalletID(linkedWallet.ID)
	}

	fw, err := walletID.GetWallet(gc.DB, gc)
	if err != nil {
		log.Printf("[ModifySharedWalletAccess] error could not get wallet object for %v %+v\n", accessInfo.WalletPublicKey, accessInfo)
		return
	}
	wallet = &fw
	var oldNumberOfApprovers int
	for _, perm := range wallet.Permissions {

		if perm.Permission == "APPROVER" {
			// oldApproverPermissionMap[perm.TargetUsername] = perm
			oldNumberOfApprovers++
		}
	}
	if oldNumberOfApprovers > 0 {
		accessInfo.MultiParty = 1
	}
	ops := make([]txnbuild.Operation, 0)
	accessInfo.Messages = make([]string, 0)
	var revokedListInfo, modifiedListInfo, addedListInfo []userModels.WalletPermissionInfo

	var numberOfSubmittedApprovers int
	var numberOfSubmittedInitiators int

	if wallet.SharedAccessEnabled == 0 {
		log.Printf("[ModifySharedWalletAccess] error shared access not enabled for %+v\n", accessInfo)
		err = &tErrors.CustomError{
			Param:      "id",
			Err:        "error-shared-access-not-active-on-wallet",
			ErrMessage: "Shared access not activated on wallet. Use option to enable shared access.",
			Code:       http.StatusForbidden,
		}
		return
	}

	if wallet.PrimaryWallet == 1 {
		if !PublicKeyHasViewOnlyAccessWACL(accessInfo.WalletPublicKey, accessInfo.AddedPermissions, gc) || !PublicKeyHasViewOnlyAccessWACL(accessInfo.WalletPublicKey, accessInfo.ModifiedPermissions, gc) {
			log.Printf("[ModifySharedWalletAccess] error only view -only access allowed for primary wallet for %+v\n", accessInfo)
			err = &tErrors.ErrorOnlyViewAccessAllowedInPrimaryWallet{}
			return
		}
	}
	checkAccess := make(map[string]userModels.WalletPermissionInfo, 0)

	//revoked permissions
	//get the permissions to be revoked
	for _, v := range accessInfo.RevokedPermissions {
		v.TargetUsername = strings.ToLower(v.TargetUsername)
		//check if username is valid
		if strings.Contains(v.TargetUsername, "_") {
			log.Printf("[ModifySharedWalletAccess] error subwallet not allowed for %+v\n %+v\n", accessInfo, v)
			//it is a subwallet and cannot be given access
			err = &tErrors.CustomError{
				Param:      "username",
				Err:        "error-subwallet-not-allowed",
				ErrMessage: fmt.Sprintf("Access can only be granted/revoked to/from trovo wallet account, not a subwallet [%v]", v.TargetUsername),
				Code:       http.StatusForbidden,
			}
			return
		}
		userPermission, e := walletID.GetUserPermissionOnWallet(v.TargetUsername, v.Permission, dbTX)
		if e != nil {
			log.Printf("[ModifySharedWalletAccess] error unable to locate existing permission for %+v\n %+v\n", accessInfo, v)
			err = &tErrors.CustomError{
				Param:      "username",
				Err:        "error-trovo-wallet-account-invalid",
				ErrMessage: fmt.Sprintf("Permission for account [%v] could not be located on this wallet at this time.", v.TargetUsername),
				Code:       http.StatusForbidden,
			}
			return
		}
		var linkedUserPermission userModels.WalletPermission
		var eLinked error

		u, e := userModels.Username(v.TargetUsername).GetSimpleUser(gc.DB, gc)
		if e != nil {
			log.Printf("[ModifySharedWalletAccess] error getting user object for %+v\n %+v\nerror: %v", accessInfo, v, e)
			err = &tErrors.CustomError{
				Param:      "username",
				Err:        "error-trovo-wallet-account-invalid",
				ErrMessage: fmt.Sprintf("Trovo wallet account [%v] could not be validated at this time.", v.TargetUsername),
				Code:       http.StatusForbidden,
			}
			return
		}
		name := fmt.Sprintf("%v", u.FirstName)
		if u.LastName != nil {
			name = fmt.Sprintf("%v %v", name, *u.LastName)
		}
		revokedListInfo = append(revokedListInfo, userModels.WalletPermissionInfo{
			TargetUsername:        v.TargetUsername,
			Name:                  name,
			Permission:            v.Permission,
			WalletPublicKey:       wallet.ID,
			WalletAlias:           wallet.Alias,
			PushNotificationToken: u.PushNotificationToken,
		})
		revokedList = append(revokedList, userModels.WalletPermission{
			ID:              userPermission.ID,
			TargetUsername:  v.TargetUsername,
			Permission:      v.Permission,
			WalletPublicKey: wallet.ID,
		})
		if hasLinkedWallet {
			linkedUserPermission, eLinked = linkedWalletID.GetUserPermissionOnWallet(v.TargetUsername, v.Permission, dbTX)
			if eLinked != nil {
				log.Printf("[ModifySharedWalletAccess] error unable to locate existing permission on linked wallet for %+v\n %+v\n", accessInfo, v)
				err = &tErrors.CustomError{
					Param:      "username",
					Err:        "error-trovo-wallet-account-invalid",
					ErrMessage: fmt.Sprintf("Permission for account [%v] could not be located on the linked wallet at this time.", v.TargetUsername),
					Code:       http.StatusForbidden,
				}
				return
			}

			linkedRevokedList = append(linkedRevokedList, userModels.WalletPermission{
				ID:              linkedUserPermission.ID,
				TargetUsername:  v.TargetUsername,
				Permission:      v.Permission,
				WalletPublicKey: linkedWallet.ID,
			})
		}

		if v.Permission == "APPROVER" {
			//get ops to add.
			op, e := generateRemoveSharedAccessOps(wallet, walletOwner, &u, gc)
			if e != nil {
				log.Printf("[ModifySharedWalletAccess] error generating blockchain operation for %+v: error: %v\n", v, e)

				// err = &tErrors.CustomError{
				// 	Param:      "username",
				// 	Err:        "error-trovo-wallet-account-invalid",
				// 	ErrMessage: fmt.Sprintf("Trovo wallet account [%v] could not be validated on the blockchain at this time.", v.TargetUsername),
				// 	Code:       http.StatusBadRequest,
				// }

				return

			} else {
				ops = append(ops, op)
			}

			if hasLinkedWallet {
				//get ops to add.
				op, e := generateRemoveSharedAccessOps(&linkedWallet, walletOwner, &u, gc)
				if e != nil {
					log.Printf("[ModifySharedWalletAccess] error generating blockchain operation for %+v ON linked wallet: error: %v\n", v, e)
					return

				} else {
					ops = append(ops, op)
				}

			}

		}

	}

	//modified permissions
	for _, v := range accessInfo.ModifiedPermissions {

		v.TargetUsername = strings.ToLower(v.TargetUsername)
		if strings.Contains(v.TargetUsername, "_") {
			log.Printf("[ModifySharedWalletAccess] error subwallet not allowed for %+v\n %+v\n", accessInfo, v)

			//it is a subWallet and cannot be given access
			err = &tErrors.CustomError{
				Param:      "username",
				Err:        "error-subwallet-not-allowed",
				ErrMessage: fmt.Sprintf("Access can only be granted to trovo wallet account, not a subwallet [%v]", v.TargetUsername),
				Code:       http.StatusForbidden,
			}
			return
		}

		var linkedEPermission userModels.WalletPermission
		var eLinked error
		ePermission, e := walletID.GetUserPermissionOnWallet(v.TargetUsername, v.Permission, dbTX)
		if e != nil {
			log.Printf("[ModifySharedWalletAccess] error unable to locate existing permission for %+v\n %+v\nerror:%v\n", accessInfo, v, e)
			err = &tErrors.CustomError{
				Param:      "username",
				Err:        "error-trovo-wallet-account-invalid",
				ErrMessage: fmt.Sprintf("Permission for account [%v] could not be located on this wallet at this time.", v.TargetUsername),
				Code:       http.StatusBadRequest,
			}
			return
		}

		if ePermission.Permission == v.Permission {
			// no changes to be made
			continue
		}
		if v.TargetUsername == walletOwner.Username && v.Permission == "VIEW-ONLY" {
			//skip adding wallet owner as VIEW-ONLY
			accessInfo.Messages = append(accessInfo.Messages, fmt.Sprintf("%v already has view access as owner, so the assigned View access has been skipped.", v.TargetUsername))
			continue
		}

		if _, ok := checkAccess[v.TargetUsername+v.Permission]; ok {
			continue
		}

		//check if username is valid

		u, e := userModels.Username(v.TargetUsername).GetSimpleUser(gc.DB, gc)
		if e != nil {
			log.Printf("[ModifySharedWalletAccess] error getting user object for %+v\n %+v\nerror: %v", accessInfo, v, e)

			err = &tErrors.CustomError{
				Param:      "username",
				Err:        "error-trovo-wallet-account-invalid",
				ErrMessage: fmt.Sprintf("Trovo wallet account [%v] could not be validated at this time.", v.TargetUsername),
				Code:       http.StatusBadRequest,
			}
			return
		}

		name := fmt.Sprintf("%v", u.FirstName)
		if u.LastName != nil {
			name = fmt.Sprintf("%v %v", name, *u.LastName)
		}
		modifiedListInfo = append(modifiedListInfo, userModels.WalletPermissionInfo{
			TargetUsername:        v.TargetUsername,
			Name:                  name,
			Permission:            v.Permission,
			WalletPublicKey:       wallet.ID,
			WalletAlias:           wallet.Alias,
			PushNotificationToken: u.PushNotificationToken,
		})
		modifiedList = append(modifiedList, userModels.WalletPermission{
			CreatedAt:       ePermission.CreatedAt,
			UpdatedAt:       ePermission.UpdatedAt,
			ID:              ePermission.ID,
			WalletPublicKey: wallet.ID,
			TargetUsername:  ePermission.TargetUsername,
			Permission:      v.Permission, //modify the permission
		})
		if hasLinkedWallet {
			linkedEPermission, eLinked = linkedWalletID.GetUserPermissionOnWallet(v.TargetUsername, v.Permission, dbTX)
			if eLinked != nil {
				log.Printf("[ModifySharedWalletAccess] error unable to locate existing permission on linked wallet for %+v\n %+v\nerror:%v\n", accessInfo, v, eLinked)
				err = &tErrors.CustomError{
					Param:      "username",
					Err:        "error-trovo-wallet-account-invalid",
					ErrMessage: fmt.Sprintf("Permission for account [%v] could not be located on the linked wallet at this time.", v.TargetUsername),
					Code:       http.StatusBadRequest,
				}
				return
			}

			linkedModifiedList = append(linkedModifiedList, userModels.WalletPermission{
				CreatedAt:       linkedEPermission.CreatedAt,
				UpdatedAt:       linkedEPermission.UpdatedAt,
				ID:              linkedEPermission.ID,
				WalletPublicKey: linkedWallet.ID,
				TargetUsername:  linkedEPermission.TargetUsername,
				Permission:      v.Permission, //modify the permission
			})

		}

		// v.permission is the new peremission
		if v.Permission == "APPROVER" {
			// attempt to add the public key as signer
			o, m, e := generateAddSharedAccessOps(wallet, walletOwner, &u, gc)
			if e != nil {
				log.Printf("[ModifySharedWalletAccess] error generating blockchain operation to add access for %+v: error: %v\n", v, e)

				err = &tErrors.CustomError{
					Param:      "username",
					Err:        "error-trovo-wallet-account-invalid",
					ErrMessage: fmt.Sprintf("Trovo wallet account [%v] could not be validated on the blockchain at this time. Unable to add this access for this user.", v.TargetUsername),
					Code:       http.StatusBadRequest,
				}

				return
			}

			ops = append(ops, o...)
			accessInfo.Messages = append(accessInfo.Messages, m...)
			if hasLinkedWallet {
				// attempt to add the public key as signer
				o, m, e := generateAddSharedAccessOps(&linkedWallet, walletOwner, &u, gc)
				if e != nil {
					log.Printf("[ModifySharedWalletAccess] error generating blockchain operation to add access on linked wallet for %+v: error: %v\n", v, e)

					err = &tErrors.CustomError{
						Param:      "username",
						Err:        "error-trovo-wallet-account-invalid",
						ErrMessage: fmt.Sprintf("Trovo wallet account [%v] could not be validated on the blockchain at this time. Unable to add this access on linked wallet for this user.", v.TargetUsername),
						Code:       http.StatusBadRequest,
					}

					return
				}

				ops = append(ops, o...)
				accessInfo.Messages = append(accessInfo.Messages, m...)
			}
		}
		// if is a downgrade of access from approver
		if ePermission.Permission == "APPROVER" {
			// attempt to remove the public key as signer
			op, e := generateRemoveSharedAccessOps(wallet, walletOwner, &u, gc)
			if e != nil {
				log.Printf("[ModifySharedWalletAccess] error generating blockchain operation for %+v: error: %v\n", v, e)

				err = &tErrors.CustomError{
					Param:      "username",
					Err:        "error-trovo-wallet-account-invalid",
					ErrMessage: fmt.Sprintf("Trovo wallet account [%v] could not be validated on the blockchain at this time. Unable to remove the approver access from blockchain at this time.", v.TargetUsername),
					Code:       http.StatusBadRequest,
				}

				return

			}

			ops = append(ops, op)

			if hasLinkedWallet {
				// attempt to remove the public key as signer
				op, e := generateRemoveSharedAccessOps(&linkedWallet, walletOwner, &u, gc)
				if e != nil {
					log.Printf("[ModifySharedWalletAccess] error generating blockchain operation on linked wallet for %+v: error: %v\n", v, e)

					err = &tErrors.CustomError{
						Param:      "username",
						Err:        "error-trovo-wallet-account-invalid",
						ErrMessage: fmt.Sprintf("Trovo wallet account [%v] could not be validated on the blockchain at this time. Unable to remove the approver access on linked wallet from blockchain at this time.", v.TargetUsername),
						Code:       http.StatusBadRequest,
					}

					return

				}

				ops = append(ops, op)

			}

		}
		checkAccess[v.TargetUsername+v.Permission] = userModels.WalletPermissionInfo{
			TargetUsername:        v.TargetUsername,
			Name:                  name,
			Permission:            v.Permission,
			WalletPublicKey:       wallet.ID,
			WalletAlias:           wallet.Alias,
			PushNotificationToken: u.PushNotificationToken,
		}

	}
	//added permissions
	for _, v := range accessInfo.AddedPermissions {

		v.TargetUsername = strings.ToLower(v.TargetUsername)
		if strings.Contains(v.TargetUsername, "_") {
			log.Printf("[ModifySharedWalletAccess] error generating blockchain operation for %+v\n", v)

			//it is a subWallet and cannot be given access
			err = &tErrors.CustomError{
				Param:      "username",
				Err:        "error-subwallet-not-allowed",
				ErrMessage: fmt.Sprintf("Access can only be granted to trovo wallet account, not a subwallet [%v]", v.TargetUsername),
				Code:       http.StatusForbidden,
			}
			return
		}
		if _, ok := checkAccess[v.TargetUsername+v.Permission]; ok {
			continue
		}
		_, e := walletID.GetUserPermissionOnWallet(v.TargetUsername, v.Permission, dbTX)
		if e == nil {
			// permission exists, so cannot be added
			continue
		}
		if v.TargetUsername == walletOwner.Username && v.Permission == "VIEW-ONLY" {
			//skip owners being added as view only access. Owners have view access by default.
			accessInfo.Messages = append(accessInfo.Messages, fmt.Sprintf("%v already has view access as owner, so the assigned View access has been skipped.", v.TargetUsername))
			continue
		}

		permissionID := uuid.NewString()

		u, e := userModels.Username(v.TargetUsername).GetSimpleUser(gc.DB, gc)
		if e != nil {
			log.Printf("[ModifySharedWalletAccess] error validating user object for %+v\n", v)

			err = &tErrors.CustomError{
				Param:      "username",
				Err:        "error-trovo-wallet-account-invalid",
				ErrMessage: fmt.Sprintf("Trovo wallet account [%v] could not be validated at this time.", v.TargetUsername),
				Code:       http.StatusForbidden,
			}
			return
		}

		name := fmt.Sprintf("%v", u.FirstName)
		if u.LastName != nil {
			name = fmt.Sprintf("%v %v", name, *u.LastName)
		}
		// infor of shared access users
		addedListInfo = append(addedListInfo, userModels.WalletPermissionInfo{
			TargetUsername:        v.TargetUsername,
			Name:                  name,
			Permission:            v.Permission,
			WalletPublicKey:       wallet.ID,
			WalletAlias:           wallet.Alias,
			PushNotificationToken: u.PushNotificationToken,
		})

		addedList = append(addedList, userModels.WalletPermission{
			ID:              permissionID,
			TargetUsername:  v.TargetUsername,
			Permission:      v.Permission,
			WalletPublicKey: wallet.ID,
		})

		if hasLinkedWallet {
			addedList = append(addedList, userModels.WalletPermission{
				ID:              uuid.NewString(),
				TargetUsername:  v.TargetUsername,
				Permission:      v.Permission,
				WalletPublicKey: linkedWallet.ID,
			})
		}

		if v.Permission == "APPROVER" {
			// attempt to add the public key as signer
			o, m, e := generateAddSharedAccessOps(wallet, walletOwner, &u, gc)
			if e != nil {
				log.Printf("[ModifySharedWalletAccess] error generating blockchain operation to add access for %+v: error: %v\n", v, e)

				err = &tErrors.CustomError{
					Param:      "username",
					Err:        "error-trovo-wallet-account-not-validated",
					ErrMessage: fmt.Sprintf("Trovo wallet account [%v] could not be validated on the blockchain at this time. Unable to add this access at this time.", v.TargetUsername),
					Code:       http.StatusBadRequest,
				}

				return
			}

			ops = append(ops, o...)
			accessInfo.Messages = append(accessInfo.Messages, m...)

			{
				oRecoveredAccount := generateRemoveRecoveredAccountAccessOps(wallet, u.Username, gc)
				if len(oRecoveredAccount) > 0 {
					ops = append(ops, oRecoveredAccount...)
				}
			}

			if hasLinkedWallet {
				// attempt to add the public key as signer
				o, m, e := generateAddSharedAccessOps(&linkedWallet, walletOwner, &u, gc)
				if e != nil {
					log.Printf("[ModifySharedWalletAccess] error generating blockchain operation to add access on linked wallet for %+v: error: %v\n", v, e)

					err = &tErrors.CustomError{
						Param:      "username",
						Err:        "error-trovo-wallet-account-not-validated",
						ErrMessage: fmt.Sprintf("Trovo wallet account [%v] could not be validated on the blockchain at this time. Unable to add this access on linked wallet at this time.", v.TargetUsername),
						Code:       http.StatusBadRequest,
					}

					return
				}

				ops = append(ops, o...)
				accessInfo.Messages = append(accessInfo.Messages, m...)

				{
					oRecoveredAccount := generateRemoveRecoveredAccountAccessOps(&linkedWallet, u.Username, gc)
					if len(oRecoveredAccount) > 0 {
						ops = append(ops, oRecoveredAccount...)
					}
				}
			}

		}
		checkAccess[v.TargetUsername+v.Permission] = userModels.WalletPermissionInfo{
			TargetUsername:        v.TargetUsername,
			Name:                  name,
			Permission:            v.Permission,
			WalletPublicKey:       wallet.ID,
			WalletAlias:           wallet.Alias,
			PushNotificationToken: u.PushNotificationToken,
		}

	}
	{
		//delete, save and create new records to be sure of what the real state now is.
		if len(revokedList) > 0 {
			e := dbTX.Delete(&revokedList).Error
			if e != nil {
				log.Println("[ModifySharedWalletAccess] error deleting revoked list: ", e)
				err = &tErrors.ErrorTemporaryServerError{}
				return
			}
			if hasLinkedWallet {
				e := dbTX.Delete(&linkedRevokedList).Error
				if e != nil {
					log.Println("[ModifySharedWalletAccess] error deleting linked revoked list: ", e)
					err = &tErrors.ErrorTemporaryServerError{}
					return
				}
			}
		}
		if len(modifiedList) > 0 {
			e := dbTX.Save(&modifiedList).Error
			if e != nil {
				log.Println("[ModifySharedWalletAccess] error saving modified list: ", e)
				err = &tErrors.ErrorTemporaryServerError{}
				return
			}
			if hasLinkedWallet {
				e := dbTX.Save(&linkedModifiedList).Error
				if e != nil {
					log.Println("[ModifySharedWalletAccess] error saving linked modified list: ", e)
					err = &tErrors.ErrorTemporaryServerError{}
					return
				}
			}
		}
		if len(addedList) > 0 {
			e := dbTX.Create(&addedList).Error
			if e != nil {
				log.Println("[ModifySharedWalletAccess] error creating added permissions: ", e)
				err = &tErrors.ErrorTemporaryServerError{}
				return
			}

			if hasLinkedWallet {
				e := dbTX.Create(&linkedAddedList).Error
				if e != nil {
					log.Println("[ModifySharedWalletAccess] error creating linked added permissions: ", e)
					err = &tErrors.ErrorTemporaryServerError{}
					return
				}
			}
		}

	}
	//invalidate all existing cache relating to this wallet
	var updatedLinkedWallet userModels.UserWallet
	var eLinked error
	wallet.InvalidateUserCache(gc)

	//it was successfully saved. now refresh the list to know the standing.
	updatedWallet, e := walletID.GetWallet(dbTX, gc)
	if e != nil {
		log.Println("[ModifySharedWalletAccess] error fetching updated wallet: ", e)
		err = &tErrors.ErrorTemporaryServerError{}
		return
	}
	if len(updatedWallet.Permissions) == 0 {
		err = &tErrors.CustomError{
			Param:      "permissions",
			Err:        "error-invalid-operation",
			ErrMessage: "Removing all access permissions is same as disabling shared access on the wallet. Please use the option to disable shared access on this wallet.",
			Code:       http.StatusForbidden,
		}
		return
	}
	if hasLinkedWallet {
		linkedWallet.InvalidateUserCache(gc)
		//it was successfully saved. now refresh the list to know the standing.
		updatedLinkedWallet, eLinked = walletID.GetWallet(dbTX, gc)
		if eLinked != nil {
			log.Println("[ModifySharedWalletAccess] error fetching updated wallet: ", eLinked)
			err = &tErrors.ErrorTemporaryServerError{}
			return
		}
	}

	{
		for _, p := range updatedWallet.Permissions {
			if p.Permission == "VIEW-ONLY" {
				viewOnly[p.TargetUsername] = p.Permission
			}
			if p.Permission == "INITIATOR" {
				numberOfSubmittedInitiators++
			}
			if p.Permission == "APPROVER" {
				numberOfSubmittedApprovers++
			}
		}
	}
	{
		for _, p := range updatedWallet.Permissions {
			if p.Permission == "VIEW-ONLY" {
				continue
			}
			_, ok := viewOnly[p.TargetUsername]
			if ok {
				err = &tErrors.CustomError{
					Param:      "numberOfApprovalsNeeded",
					Err:        "error-invalid-permission",
					ErrMessage: fmt.Sprintf("%v cannot have VIEW access after being granted an %v access on the same wallet.", p.TargetUsername, p.Permission),
					Code:       http.StatusForbidden,
				}
				return
			}

		}
	}

	if numberOfSubmittedApprovers <= accessInfo.NumberOfApprovalsNeeded && accessInfo.NumberOfApprovalsNeeded > 1 {
		//number of authorizers does not reach the minimum threshold needed. cannot proceed so as to prevent account lockout
		err = &tErrors.CustomError{
			Param:      "numberOfApprovalsNeeded",
			Err:        "error-approvers-not-enough",
			ErrMessage: fmt.Sprintf("Please ensure that your list of approvers configured are greater than the minimum number required to approve a transaction [%v]", accessInfo.NumberOfApprovalsNeeded),
			Code:       http.StatusForbidden,
		}
		return
	}
	if numberOfSubmittedApprovers > 0 && accessInfo.NumberOfApprovalsNeeded == 0 {
		//number of APPROVERS does not reach the minimum threshold needed. cannot proceed so as to prevent account lockout
		err = &tErrors.CustomError{
			Param:      "numberOfApprovalsNeeded",
			Err:        "error-approvers-not-enough",
			ErrMessage: "You must specify the number of approvers needed to approve transactions on this wallet",
			Code:       http.StatusForbidden,
		}
		return
	}

	if numberOfSubmittedApprovers > 0 && numberOfSubmittedInitiators == 0 {
		err = &tErrors.CustomError{
			Param:      "numberOfApprovalsNeeded",
			Err:        "error-initiator-missing",
			ErrMessage: "You must specify at least one initiator when an approver is specified.",
			Code:       http.StatusForbidden,
		}
		return
	}

	if numberOfSubmittedApprovers == 0 && numberOfSubmittedInitiators > 0 {
		err = &tErrors.CustomError{
			Param:      "numberOfApprovers",
			Err:        "error-initiator-missing",
			ErrMessage: "You cannot have an initiator access when you have not specified approvers.",
			Code:       http.StatusForbidden,
		}
		return
	}

	updatedWallet.SharedAccessEnabled = 1
	updatedWallet.NumberOfApprovalsNeeded = accessInfo.NumberOfApprovalsNeeded

	errDB := dbTX.Omit(clause.Associations).Save(&updatedWallet).Error
	if errDB != nil {
		log.Printf("[ModifySharedWalletAccess] error saving shared access status of the wallet:%v\n sharedAccess:%+v\n", errDB, updatedWallet)
		err = &tErrors.ErrorTemporaryServerError{}
		return
	}
	if hasLinkedWallet {
		updatedLinkedWallet.SharedAccessEnabled = updatedWallet.SharedAccessEnabled
		updatedLinkedWallet.NumberOfApprovalsNeeded = updatedWallet.NumberOfApprovalsNeeded

		errDB := dbTX.Omit(clause.Associations).Save(&updatedLinkedWallet).Error
		if errDB != nil {
			log.Printf("[ModifySharedWalletAccess] error saving shared access status of the linked wallet:%v\n sharedAccess:%+v\n", errDB, updatedLinkedWallet)
			err = &tErrors.ErrorTemporaryServerError{}
			return
		}
	}
	xdrBase64, transactionSource, messages, errGenXdr := generateModifySharedAccessXdr(wallet, walletOwner, numberOfSubmittedApprovers, accessInfo.NumberOfApprovalsNeeded, oldNumberOfApprovers, ops, gc)
	if errGenXdr != nil {
		err = errGenXdr
		return
	}
	accessInfo.TransactionSource = transactionSource
	accessInfo.Messages = append(accessInfo.Messages, messages...)

	accessInfo.NetworkPassPhrase = network.GetBlockchainNetworkPassPhrase()

	accessInfo.Transaction = xdrBase64

	// accessInfo.RevokedPermissions = revokedListInfo
	// accessInfo.ModifiedPermissions = modifiedListInfo
	// accessInfo.AddedPermissions = addedListInfo

	if oldNumberOfApprovers == 0 {
		accessInfo.SignatureRequired = 1
	}
	accessInfo.RevokedPermissions = revokedListInfo
	accessInfo.ModifiedPermissions = modifiedListInfo
	accessInfo.AddedPermissions = addedListInfo
	if len(accessInfo.TransactionSignature) == 0 && accessInfo.Commit == 0 {
		err = nil
		return
	}
	if len(accessInfo.TransactionSignature) > 0 && accessInfo.Commit == 0 {
		// extract signature and submit transaction
		//submit to blockchain
		var txnHash string
		signatures := make(map[string]string, 0)
		signatures[wallet.Signer] = accessInfo.TransactionSignature

		txnHash, err = network.SubmitXdrWithSignatures(gc.BantuExpansionClient, xdrBase64, signatures, gc.DB)
		if err != nil {
			log.Printf("Error submitting shared access txn [%+v] transaction: %s\n", accessInfo, err.Error())
			// logDiscordFailedRecovery(fmt.Sprintf("Error submitting shared access txn [%+v] transaction: %s", accessInfo, err.Error()))
			err = &tErrors.ErrorTemporaryServerError{}
			return
		}
		accessInfo.TransactionID = txnHash

		dbTX.Commit()
		err = nil
		wallet.InvalidateUserCache(gc)
		return

	}

	if len(accessInfo.TransactionSignature) == 0 && oldNumberOfApprovers > 0 && accessInfo.Commit == 1 {
		//save as pending transaction request
		description := fmt.Sprintf("Modify shared access on wallet %v.", wallet.Alias)
		// description := fmt.Sprintf("%v\n%v", wallet.Alias)
		id := uuid.New().String()

		if len(revokedList) > 0 {
			revokedUsers := ""
			for i, u := range revokedList {
				revokedUsers = fmt.Sprintf("%s(%v)", u.TargetUsername, u.Permission)
				if i < len(revokedList)-1 {
					revokedUsers = fmt.Sprintf("%s, ", revokedUsers)
				}
			}
			description = fmt.Sprintf("%s\n Permissions to revoke:%v.", description, revokedUsers)
		}

		if len(modifiedList) > 0 {
			modifiedUsers := ""
			for i, u := range modifiedList {
				modifiedUsers = fmt.Sprintf("%s(%v)", u.TargetUsername, u.Permission)
				if i < len(revokedList)-1 {
					modifiedUsers = fmt.Sprintf("%s, ", modifiedUsers)
				}
			}
			description = fmt.Sprintf("%s\n Modifying Permissions:%v.", description, modifiedUsers)
		}

		if len(addedList) > 0 {
			addedUsers := ""
			for i, u := range addedList {
				addedUsers = fmt.Sprintf("%s(%v)", u.TargetUsername, u.Permission)
				if i < len(addedList)-1 {
					addedUsers = fmt.Sprintf("%s, ", addedUsers)
				}
			}
			description = fmt.Sprintf("%s\n Adding New Permissions:%v.", description, addedUsers)
		}
		transactionByte, _ := json.Marshal(*accessInfo)
		transactionStr := string(transactionByte)
		pendingAuth := userModels.PendingAuth{
			ID:                       id,
			Initiator:                signerUser.Username,
			InitiatorSignerPublicKey: signerUser.PrimarySigner,
			WalletPublicKey:          wallet.ID,
			TransactionType:          "MODIFY SHARED ACCESS",
			Description:              description,
			TransactionSource:        transactionSource,
			ApprovalsNeeded:          accessInfo.NumberOfApprovalsNeeded,
			TransactionXdr:           xdrBase64,
			TransactionInfoStr:       &transactionStr,
		}
		// rollback all the other changes since the changes can only apply when approvals are completed.
		dbTX.Rollback()
		//save and commit this to database
		e := gc.DB.Create(&pendingAuth).Error
		if e != nil {
			log.Printf("Error saving modify shared access txn [%+v] transaction on pending auth table: %s\n", pendingAuth, e.Error())
			err = &tErrors.ErrorTemporaryServerError{}
			return
		}
		//set the transaction id
		accessInfo.TransactionID = "PENDING_AUTH"

		return

	}
	err = &tErrors.CustomError{
		Param:      "transaction",
		Err:        "error-known-state",
		ErrMessage: "This request does not meet any known conditions for execution.",
		Code:       http.StatusNotFound,
	}
	return
}

func RemoveSharedWalletAccess(signerUser *userModels.User, wallet *userModels.UserWallet, accessInfo *userModels.DisableSharedAccessInfo, gc *sharedconfig.GlobalConfig) (err error) {
	// var managedAccess userModels.UserWalletSharedAccess
	var hasLinkedWallet bool
	var linkedWallet userModels.UserWallet
	if wallet.WalletType == 1 && wallet.LinkedWalletPublicKey != nil {
		hasLinkedWallet = true
	}

	if w, v := wallet.IsValidLinkedWallet(gc); v {
		return &tErrors.CustomError{
			Param:      "publicKey",
			Err:        "error-linked-wallets-not-allowed",
			ErrMessage: fmt.Sprintf("Linked Wallets are not allowed to be shared directly. Plase share the access on %v and it will mirror to this wallet.", w.Alias),
			Code:       http.StatusBadRequest,
		}

	}
	var linkedAccessList []userModels.WalletPermission
	var errLinked error
	accessList := wallet.Permissions
	if hasLinkedWallet {
		linkedWallet, errLinked = userModels.UserWalletID(*wallet.LinkedWalletPublicKey).GetWallet(gc.DB, gc)
		if errLinked != nil {
			return &tErrors.CustomError{
				Param:      "username",
				Err:        "error-confirming-linked-wallet",
				ErrMessage: "Unable to confirm linked wallet at this time. Please try again after some minutes.",
				Code:       http.StatusForbidden,
			}
		}
		linkedAccessList = linkedWallet.Permissions
	}
	var numberOfApprovers int
	// var numberOfSubmittedInitiators int
	// var selfApprover int
	var approverUsers []*userModels.User

	if wallet.SharedAccessEnabled == 0 {
		return &tErrors.CustomError{
			Param:      "id",
			Err:        "error-shared-access-not-active-on-wallet",
			ErrMessage: "Shared access not activated on wallet.",
			Code:       http.StatusForbidden,
		}
	}
	if wallet.WalletType == 2 || wallet.WalletType == 3 {
		return &tErrors.CustomError{
			Param:      "publicKey",
			Err:        "error-wallet-type-not-allowed",
			ErrMessage: "Market Making & Bulk Payment wallets are not allowed for this operation.",
			Code:       http.StatusForbidden,
		}
	}

	approvalsNeeded := wallet.NumberOfApprovalsNeeded
	var userPermissions string
	walletOwner, e := wallet.GetWalletOwner(gc.DB, gc)
	if e != nil {
		return &tErrors.CustomError{
			Param:      "username",
			Err:        "error-confirming-wallet-owner",
			ErrMessage: "Unable to confirm wallet owner at this time. Please try again after some minutes.",
			Code:       http.StatusForbidden,
		}
	}

	for i, v := range accessList {
		u, e := userModels.Username(v.TargetUsername).GetSimpleUser(gc.DB, gc)
		if e != nil {
			return &tErrors.CustomError{
				Param:      "username",
				Err:        "error-trovo-wallet-account-invalid",
				ErrMessage: fmt.Sprintf("Trovo wallet account [%v] could not be validated at this time.", v.TargetUsername),
				Code:       http.StatusForbidden,
			}
		}
		if v.Permission == "APPROVER" {
			numberOfApprovers++
			// if v.TargetUsername == walletOwner.Username {
			// 	selfApprover = 1
			// }
			approverUsers = append(approverUsers, &u)
		}
		{
			accessInfo.Permissions = append(accessInfo.Permissions, v.ToWalletPermissionInfo(&u, wallet, gc))
		}
		userPermissions = fmt.Sprintf("%s(%v)", v.TargetUsername, v.Permission)
		if i < len(accessList)-1 {
			userPermissions = fmt.Sprintf("%s, ", userPermissions)
		}
	}

	if numberOfApprovers > 0 {
		accessInfo.MultiParty = 1

	}

	xdrBase64, transactionSource, messages, walletMustSign, _, errGenXdr := generateRemoveSharedAccessXdr(wallet, &walletOwner, approverUsers, numberOfApprovers, gc)
	if errGenXdr != nil {
		return errGenXdr
	}
	accessInfo.Messages = messages
	if walletMustSign {
		accessInfo.SignatureRequired = 1

	}
	accessInfo.TransactionSource = transactionSource
	accessInfo.NetworkPassPhrase = network.GetBlockchainNetworkPassPhrase()

	accessInfo.Transaction = xdrBase64

	if len(accessInfo.TransactionSignature) == 0 {
		if accessInfo.Commit == 0 {
			return nil
		}
	}

	if numberOfApprovers > 0 && accessInfo.Commit == 1 {
		//SET transaction id to pending auth
		accessInfo.MultiParty = 1
		if accessInfo.Commit == 1 {
			accessInfo.TransactionID = "PENDING_AUTH"

			id := uuid.New().String()

			description := fmt.Sprintf("Disabling shared access on wallet %v.\n This will remove the permissions:\n%v", wallet.Alias, userPermissions)
			transactionByte, _ := json.Marshal(*accessInfo)
			transactionStr := string(transactionByte)
			pendingAuth := userModels.PendingAuth{
				ID:                       id,
				Initiator:                signerUser.Username,
				InitiatorSignerPublicKey: signerUser.PrimarySigner,
				WalletPublicKey:          wallet.ID,
				TransactionType:          "DISABLE SHARED ACCESS",
				Description:              description,
				TransactionSource:        accessInfo.TransactionSource,
				ApprovalsNeeded:          approvalsNeeded,
				TransactionXdr:           xdrBase64,
				TransactionInfoStr:       &transactionStr,
			}
			// rollback all the other changes since the changes can only apply when approvals are completed.
			// dbTX.Rollback()
			//save and commit this to database
			e := gc.DB.Create(&pendingAuth).Error
			if e != nil {
				log.Printf("Error saving disable shared access txn [%+v] transaction on pending auth table: %s\n", pendingAuth, e.Error())
				return &tErrors.ErrorTemporaryServerError{}
			}

		}

		return nil
	}
	if len(accessInfo.TransactionSignature) == 0 && wallet.HasViewOnlyAccess(gc) {

		return nil

	}
	if !wallet.HasViewOnlyAccess(gc) {
		log.Println("[RemoveSharedWalletAccess] wallet is not view only and did not meet condition for multiparty")
		return &tErrors.ErrorTemporaryServerError{}

	}
	dbTX := gc.DB.Begin()
	defer dbTX.Rollback()
	wallet.SharedAccessEnabled = 0
	wallet.NumberOfApprovalsNeeded = 0
	if hasLinkedWallet {
		linkedWallet.SharedAccessEnabled = 0
		linkedWallet.NumberOfApprovalsNeeded = 0
	}
	e = dbTX.Omit(clause.Associations).Save(wallet).Error
	if e != nil {
		log.Println("[RemoveSharedWalletAccess] error saving wallet", e)
		return &tErrors.ErrorTemporaryServerError{}
	}
	e = dbTX.Delete(&accessList).Error
	if e != nil {
		log.Println("[RemoveSharedWalletAccess] error deleting access list", e)
		return &tErrors.ErrorTemporaryServerError{}
	}
	if hasLinkedWallet {
		e = dbTX.Omit(clause.Associations).Save(&linkedWallet).Error
		if e != nil {
			log.Println("[RemoveSharedWalletAccess] error saving linked wallet", e)
			return &tErrors.ErrorTemporaryServerError{}
		}
		e = dbTX.Delete(&linkedAccessList).Error
		if e != nil {
			log.Println("[RemoveSharedWalletAccess] error deleting linked access list", e)
			return &tErrors.ErrorTemporaryServerError{}
		}
	}

	// extract signature and submit transaction
	//getting here means it does not contain approvers
	//submit to blockchain
	signatures := make(map[string]string, 0)
	signatures[wallet.Signer] = accessInfo.TransactionSignature
	txnHash, err := network.SubmitXdrWithSignatures(gc.BantuExpansionClient, xdrBase64, signatures, gc.DB)
	if err != nil {
		log.Printf("[RemoveSharedWalletAccess] Error submitting disable shared access txn [%+v] transaction: %s\n", accessInfo, err.Error())
		// logDiscordFailedRecovery(fmt.Sprintf("Error submitting shared access txn [%+v] transaction: %s", accessInfo, err.Error()))
		return &tErrors.ErrorTemporaryServerError{}
	}
	accessInfo.TransactionID = txnHash

	dbTX.Commit()
	wallet.InvalidateUserCache(gc)
	return nil

}

func generateCreateSharedAccessXdr(wallet *userModels.UserWallet, walletOwner *userModels.User, approvers []*userModels.User, accessInfo []userModels.WalletPermissionInfo, authThreshold int, gc *sharedconfig.GlobalConfig) (xdrbase64 string, messages []string, walletMustSign bool, err error) {
	client := gc.BantuExpansionClient
	ops := make([]txnbuild.Operation, 0)
	messages = make([]string, 0)
	var hasLinkedWallet bool
	var errLinkedWallet error
	var linkedWallet userModels.UserWallet
	if wallet.WalletType == 1 && wallet.LinkedWalletPublicKey != nil {
		// set the linked wallet if it is a tokenization wallet
		hasLinkedWallet = true

		linkedWallet, errLinkedWallet = userModels.UserWalletID(*wallet.LinkedWalletPublicKey).GetWallet(gc.DB, gc)
		if errLinkedWallet != nil {
			err = &tErrors.CustomError{
				Param:      "linkedWalletPublicKey",
				Err:        "error-getting-linked-wallet",
				ErrMessage: "Linked Wallet could not be validated.",
				Code:       http.StatusBadRequest,
			}
			return "", messages, walletMustSign, err
		}
	}
	// totalNativeBalanceNeeded := decimal.Zero
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
	if !walletAccountExists || (walletAccountNativeBalance).LessThan(minBalance) {
		log.Printf("[generateCreateSharedAccessXdr] by [%v] shared WalletAccount underfunded \n", wallet.Alias)

		err = &tErrors.CustomError{
			Param:      "publicKey",
			Err:        "error-wallet-underfunded",
			ErrMessage: fmt.Sprintf("Wallet %v does not have enough %v balance to perform this operation", wallet.Alias, os.Getenv("NATIVE_ASSET_CODE")),
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
				{
					if hasLinkedWallet {
						ops = append(ops, &txnbuild.SetOptions{
							Signer: &txnbuild.Signer{
								Address: recoveryAddress,
								Weight:  0,
							},
							SourceAccount: linkedWallet.ID,
						})
					}
				}

				//add message about disabling recovery on that wallet
				messages = append(messages, fmt.Sprintf("Account Recovery on this wallet %v has to be disabled so as to enable shared access.", wallet.Alias))
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
			{
				if hasLinkedWallet {
					ops = append(ops, &txnbuild.SetOptions{
						Signer: &txnbuild.Signer{
							Address: user3p.PrimarySigner,
							Weight:  1,
						},
						SourceAccount: linkedWallet.ID,
					})
				}
			}

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
				{
					if hasLinkedWallet {
						ops = append(ops, &txnbuild.SetOptions{
							Signer: &txnbuild.Signer{
								Address: user3p.PrimarySigner,
								Weight:  1,
							},
							SourceAccount: linkedWallet.ID,
						})
					}
				}

				walletMustSign = true
			}

		}
		if totalUsersToFund > 0 {
			totalNativeBalanceNeeded := activationAmount.Mul(decimal.NewFromInt(totalUsersToFund))
			if walletAccountNativeBalance.LessThan(totalNativeBalanceNeeded) {
				//not enough balance to perform this.
				log.Printf("[generateCreateSharedAccessXdr] by [%v] Shared Access WalletAccount underfunded \n", wallet.Alias)

				err = &tErrors.CustomError{
					Param:      "walletPublicKey",
					Err:        "error-wallet-underfunded",
					ErrMessage: fmt.Sprintf("Wallet %v needs more than %v %v balance to perform this operation", wallet.Alias, totalNativeBalanceNeeded.String(), os.Getenv("NATIVE_ASSET_CODE")),
					Code:       404,
				}
				return "", messages, walletMustSign, err
			}
		}
	}

	//activating shared access is free. No fee.

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
			{
				if hasLinkedWallet {
					ops = append(ops, &txnbuild.SetOptions{
						LowThreshold:    txnbuild.NewThreshold(txnbuild.Threshold(authThreshold)),
						MediumThreshold: txnbuild.NewThreshold(txnbuild.Threshold(authThreshold)),
						HighThreshold:   txnbuild.NewThreshold(txnbuild.Threshold(authThreshold)),
						SourceAccount:   linkedWallet.ID,
					})
				}
			}
		}
	}

	if len(ops) == 0 {
		// no operations to sign. create a dummy ops, will be ignored on next try.
		ops = append(ops, &txnbuild.SetOptions{
			LowThreshold:    txnbuild.NewThreshold(txnbuild.Threshold(0)),
			MediumThreshold: txnbuild.NewThreshold(txnbuild.Threshold(0)),
			HighThreshold:   txnbuild.NewThreshold(txnbuild.Threshold(0)),
			SourceAccount:   wallet.ID,
		})
		walletMustSign = true
		{
			if hasLinkedWallet {
				ops = append(ops, &txnbuild.SetOptions{
					LowThreshold:    txnbuild.NewThreshold(txnbuild.Threshold(0)),
					MediumThreshold: txnbuild.NewThreshold(txnbuild.Threshold(0)),
					HighThreshold:   txnbuild.NewThreshold(txnbuild.Threshold(0)),
					SourceAccount:   linkedWallet.ID,
				})
			}
		}

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

func generateModifySharedAccessXdr(wallet *userModels.UserWallet, walletOwner *userModels.User, numberOfSubmittedApprovers, numberOfApprovalsNeeded, oldNumberOfApprovers int, ops []txnbuild.Operation, gc *sharedconfig.GlobalConfig) (xdrbase64, transactionSource string, messages []string, err error) {
	var hasLinkedWallet bool
	if wallet.WalletType == 1 && wallet.LinkedWalletPublicKey != nil {
		hasLinkedWallet = true
	}

	client := gc.BantuExpansionClient
	messages = make([]string, 0)
	// totalNativeBalanceNeeded := decimal.Zero
	var activationAmount = decimal.NewFromFloat(6)
	var minBalance = decimal.NewFromFloat(3.0)
	if len(os.Getenv("SUB_WALLET_ACTIVATION_AMOUNT")) > 0 {
		activationAmount = decimal.RequireFromString(os.Getenv("SUB_WALLET_ACTIVATION_AMOUNT"))
	}
	if len(os.Getenv("WALLET_MINIMUM_BALANCE")) > 0 {
		minBalance = decimal.RequireFromString(os.Getenv("WALLET_MINIMUM_BALANCE"))
	}

	if numberOfSubmittedApprovers == 0 && numberOfApprovalsNeeded > 0 {
		err = &tErrors.CustomError{
			Param:      "numberOfApprovers",
			Err:        "error-no-approver-specified",
			ErrMessage: "No approvers specifieds",
			Code:       404,
		}
		return "", "", messages, err
	}
	chanAccount := <-gc.ChannelAccounts
	defer func(c *keypair.Full) {
		gc.ChannelAccounts <- c
	}(chanAccount)
	// paymentInfo.Messages = messages
	_, _, _, _, chanSourceAccount, _ := network.BlockchainAccountProperties(client, chanAccount.Address(), txnbuild.NativeAsset{})

	//check if primary account has native enough native balance
	var nativeAsset txnbuild.Asset = txnbuild.NativeAsset{}
	walletAccountExists, _, walletAccountNativeBalance, _, walletSourceAccount, errWalletAct := network.BlockchainAccountProperties(client, wallet.ID, nativeAsset)
	if errWalletAct != nil {
		log.Printf("[generateModifySharedAccessXdr] by [%v] for shared Account Properties error:[%v] \n", wallet.Alias, errWalletAct)

		return "", "", messages, errWalletAct
	}
	if !walletAccountExists || (walletAccountNativeBalance.Sub(activationAmount)).LessThan(minBalance.Mul(decimal.NewFromInt(int64(numberOfSubmittedApprovers)))) {
		log.Printf("[generateModifySharedAccessXdr] by [%v] shared WalletAccount underfunded. Needs at least %v %v\n", wallet.Alias, (minBalance.Mul(decimal.NewFromInt(int64(numberOfSubmittedApprovers)))).Truncate(7).String(), os.Getenv("NATIVE_ASSET_CODE"))

		err = &tErrors.CustomError{
			Param:      "publicKey",
			Err:        "error-wallet-underfunded",
			ErrMessage: fmt.Sprintf("Wallet %v needs minimum of %v %v balance to perform this operation.", wallet.Alias, (minBalance.Mul(decimal.NewFromInt(int64(numberOfSubmittedApprovers)))).Truncate(7).String(), os.Getenv("NATIVE_ASSET_CODE")),
			Code:       http.StatusBadRequest,
		}
		return "", "", messages, err
	}
	{
		//check if account recovery is enabled, then disable it on the wallet.
		if walletOwner.AccountRecoveryEnabled == 1 && numberOfSubmittedApprovers > 0 {
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
				if hasLinkedWallet {
					//recovery a signer to the wallet. remove it
					ops = append(ops, &txnbuild.SetOptions{
						Signer: &txnbuild.Signer{
							Address: recoveryAddress,
							Weight:  0,
						},
						SourceAccount: *wallet.LinkedWalletPublicKey,
					})

					//add message about disabling recovery on that wallet
					messages = append(messages, "Account Recovery on the linked wallet has to be disabled so as to enable shared access.")

				}
			}
		}
	}

	{
		//adjust account threshold
		if numberOfApprovalsNeeded > 0 || len(ops) == 0 {
			// len(ops) == 0 prevents empty ops error
			ops = append(ops, &txnbuild.SetOptions{
				LowThreshold:    txnbuild.NewThreshold(txnbuild.Threshold(numberOfApprovalsNeeded)),
				MediumThreshold: txnbuild.NewThreshold(txnbuild.Threshold(numberOfApprovalsNeeded)),
				HighThreshold:   txnbuild.NewThreshold(txnbuild.Threshold(numberOfApprovalsNeeded)),
				SourceAccount:   wallet.ID,
			})

			if hasLinkedWallet {
				ops = append(ops, &txnbuild.SetOptions{
					LowThreshold:    txnbuild.NewThreshold(txnbuild.Threshold(numberOfApprovalsNeeded)),
					MediumThreshold: txnbuild.NewThreshold(txnbuild.Threshold(numberOfApprovalsNeeded)),
					HighThreshold:   txnbuild.NewThreshold(txnbuild.Threshold(numberOfApprovalsNeeded)),
					SourceAccount:   *wallet.LinkedWalletPublicKey,
				})
			}
		}
	}

	// Construct the transaction that holds the operations to execute on the network
	var tx *txnbuild.Transaction
	if oldNumberOfApprovers > 0 {
		//multiparty
		transactionSource = chanSourceAccount.AccountID
		tx, err = txnbuild.NewTransaction(
			txnbuild.TransactionParams{
				SourceAccount:        chanSourceAccount,
				IncrementSequenceNum: true,
				Operations:           ops,
				BaseFee:              2000,
				Preconditions: txnbuild.Preconditions{
					TimeBounds: txnbuild.NewInfiniteTimeout(),
				},
				Memo: txnbuild.MemoText("Modify shared access"),
			},
		)
	} else {
		//single signer
		tx, err = txnbuild.NewTransaction(
			txnbuild.TransactionParams{
				SourceAccount:        walletSourceAccount,
				IncrementSequenceNum: true,
				Operations:           ops,
				BaseFee:              2000,
				Preconditions: txnbuild.Preconditions{
					TimeBounds: txnbuild.NewInfiniteTimeout(),
				},
				Memo: txnbuild.MemoText("Modify shared access"),
			},
		)
	}

	if err != nil {
		log.Println("[generateModifySharedAccessXdr] error constructing transaction ", err)
		return "", "", messages, err
	}
	if oldNumberOfApprovers > 0 {
		tx, err = tx.Sign(gc.BantuNetworkPassphrase, chanAccount)
		if err != nil {
			log.Println("[generateModifySharedAccessXdr] error signing transaction with chan account", err)
			return "", "", messages, err
		}
	}

	var xdrBase64 string

	xdrBase64, err = tx.Base64()
	if err != nil {
		log.Println("[generateModifySharedAccessXdr] error getting txn base64", err)
		return "", "", messages, err
	}

	return xdrBase64, transactionSource, messages, nil

}
func generateAddSharedAccessOps(wallet *userModels.UserWallet, walletOwner *userModels.User, approver *userModels.User, gc *sharedconfig.GlobalConfig) (ops []txnbuild.Operation, messages []string, err error) {
	client := gc.BantuExpansionClient
	var activationAmount = decimal.NewFromFloat(6)
	// var minBalance = decimal.NewFromFloat(3.0)
	if len(os.Getenv("SUB_WALLET_ACTIVATION_AMOUNT")) > 0 {
		activationAmount = decimal.RequireFromString(os.Getenv("SUB_WALLET_ACTIVATION_AMOUNT"))
	}
	ops = make([]txnbuild.Operation, 0)
	messages = make([]string, 0)

	//check if primary account has native enough native balance
	var nativeAsset txnbuild.Asset = txnbuild.NativeAsset{}
	_, _, _, _, walletSourceAccount, errWalletAct := network.BlockchainAccountProperties(client, wallet.ID, nativeAsset)
	if errWalletAct != nil {
		log.Printf("[generateCreateSharedAccessXdr] by [%v] for shared Account Properties error:[%v] \n", wallet.Alias, errWalletAct)

		return ops, messages, errWalletAct
	}

	//check access list to know if you would activate the user wallets before proceeding.

	//ensure u r using the account signer, since the account may have been recovered, or may be recovered in the future, changing the signer, but retaining the primary key
	approverAccountExists, _, _, _, _, _ := network.BlockchainAccountProperties(client, approver.PrimarySigner, nativeAsset)
	if !approverAccountExists {
		//if subwallet is not activated
		//build transaction that will activate the primary signer from the assigning wallet
		ops = append(ops, &txnbuild.CreateAccount{
			Destination:   approver.PrimarySigner,
			Amount:        activationAmount.String(),
			SourceAccount: wallet.ID,
		})

		messages = append(messages, fmt.Sprintf("%v %v will be deducted from wallet %v and be used to activate the approver account %v.", activationAmount.String(), os.Getenv("NATIVE_ASSET_CODE"), wallet.Alias, approver.Username))

		//after creation, it now exists with enough balance to add signer wallet as signer
		//ENSURE IT IS NOT the
		ops = append(ops, &txnbuild.SetOptions{
			Signer: &txnbuild.Signer{
				Address: approver.PrimarySigner,
				Weight:  1,
			},
			SourceAccount: wallet.ID,
		})

	}

	if approverAccountExists {
		//account exists, check if it already it a signer in the wallet

		//after topping up, it now has enough balance to add primary wallet as signer if it is not already a signer
		if !wallet.SignerIsValidWA(approver.PrimarySigner, walletSourceAccount) {

			if walletOwner.PrimarySigner != wallet.Signer {

				ops = append(ops, &txnbuild.SetOptions{
					Signer: &txnbuild.Signer{
						Address: approver.PrimarySigner,
						Weight:  1,
					},
					SourceAccount: wallet.ID,
				})
			} else {
				//primary signer and master signer
				ops = append(ops, &txnbuild.SetOptions{
					MasterWeight:  txnbuild.NewThreshold(1),
					SourceAccount: wallet.ID,
				})
			}

		}

	}

	return ops, messages, nil

}

func generateRemoveSharedAccessXdr(wallet *userModels.UserWallet, walletOwner *userModels.User, approvers []*userModels.User, numberOfApprovers int, gc *sharedconfig.GlobalConfig) (xdrbase64, transactionSource string, messages []string, walletMustSign, multipartySign bool, err error) {

	var hasLinkedWallet bool
	// var linkedWallet userModels.UserWallet
	if wallet.WalletType == 1 && wallet.LinkedWalletPublicKey != nil {
		hasLinkedWallet = true
	}
	// if hasLinkedWallet {
	// 	linkedWallet, _ = userModels.UserWalletID(*wallet.LinkedWalletPublicKey).GetWallet(gc.DB, gc)

	// }
	client := gc.BantuExpansionClient
	ops := make([]txnbuild.Operation, 0)
	messages = make([]string, 0)
	// totalNativeBalanceNeeded := decimal.Zero
	var activationAmount = decimal.NewFromFloat(6)
	var minBalance = decimal.NewFromFloat(3.0)
	if len(os.Getenv("SUB_WALLET_ACTIVATION_AMOUNT")) > 0 {
		activationAmount = decimal.RequireFromString(os.Getenv("SUB_WALLET_ACTIVATION_AMOUNT"))
	}
	if len(os.Getenv("WALLET_MINIMUM_BALANCE")) > 0 {
		minBalance = decimal.RequireFromString(os.Getenv("WALLET_MINIMUM_BALANCE"))
	}
	chanAccount := <-gc.ChannelAccounts
	defer func(c *keypair.Full) {
		gc.ChannelAccounts <- c
	}(chanAccount)
	// paymentInfo.Messages = messages
	_, _, _, _, chanSourceAccount, _ := network.BlockchainAccountProperties(client, chanAccount.Address(), txnbuild.NativeAsset{})

	//check if primary account has native enough native balance
	var nativeAsset txnbuild.Asset = txnbuild.NativeAsset{}
	walletAccountExists, _, walletAccountNativeBalance, _, walletSourceAccount, errWalletAct := network.BlockchainAccountProperties(client, wallet.ID, nativeAsset)
	if errWalletAct != nil {
		log.Printf("[generateRemoveSharedAccessXdr] by [%v] for shared Account Properties error:[%v] \n", wallet.Alias, errWalletAct)

		return "", "", messages, walletMustSign, multipartySign, errWalletAct
	}
	if !walletAccountExists || (walletAccountNativeBalance.Sub(activationAmount)).LessThan(minBalance) {
		log.Printf("[generateRemoveSharedAccessXdr] by [%v] shared WalletAccount underfunded \n", wallet.Alias)

		err = &tErrors.CustomError{
			Param:      "publicKey",
			Err:        "error-wallet-underfunded",
			ErrMessage: fmt.Sprintf("Wallet %v does not have enough XBN balance to perform this operation", wallet.Alias),
			Code:       404,
		}
		return "", "", messages, walletMustSign, multipartySign, err
	}
	{
		//check if it is multiparty signature that is required.
		if wallet.WalletCountApproverAccess(gc) > 0 {
			multipartySign = true
		}
		//check if account recovery is enabled, then re-enable it on the wallet.
		if walletOwner.AccountRecoveryEnabled == 1 && multipartySign {
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
				messages = append(messages, "Account Recovery on this wallet will be enabled.")
				walletMustSign = true

				if hasLinkedWallet {
					//recovery a signer to the wallet. remove it
					ops = append(ops, &txnbuild.SetOptions{
						Signer: &txnbuild.Signer{
							Address: recoveryAddress,
							Weight:  1,
						},
						SourceAccount: *wallet.LinkedWalletPublicKey,
					})

					//add message about disabling recovery on that wallet
					messages = append(messages, "Account Recovery on the linked wallet will be enabled.")
				}
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

			//remove signer if already a signer
			if wallet.SignerIsValidWA(user3p.PrimarySigner, walletSourceAccount) && walletOwner.PrimarySigner != wallet.Signer {

				ops = append(ops, &txnbuild.SetOptions{
					Signer: &txnbuild.Signer{
						Address: user3p.PrimarySigner,
						Weight:  0,
					},
					SourceAccount: wallet.ID,
				})

				walletMustSign = true

				if hasLinkedWallet {
					ops = append(ops, &txnbuild.SetOptions{
						Signer: &txnbuild.Signer{
							Address: user3p.PrimarySigner,
							Weight:  0,
						},
						SourceAccount: *wallet.LinkedWalletPublicKey,
					})
				}
			}

		}
		if totalUsersToFund > 0 {
			totalNativeBalanceNeeded := activationAmount.Mul(decimal.NewFromInt(totalUsersToFund))
			if walletAccountNativeBalance.LessThan(totalNativeBalanceNeeded) {
				//not enough balance to perform this.
				log.Printf("[generateRemoveSharedAccessXdr] by [%v] Shared Access WalletAccount underfunded \n", wallet.Alias)

				err = &tErrors.CustomError{
					Param:      "walletPublicKey",
					Err:        "error-wallet-underfunded",
					ErrMessage: fmt.Sprintf("Wallet %v needs more than %v %v balance to perform this operation", wallet.Alias, totalNativeBalanceNeeded.String(), os.Getenv("NATIVE_ASSET_CODE")),
					Code:       404,
				}
				return "", "", messages, walletMustSign, multipartySign, err
			}
		}
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
		if hasLinkedWallet {
			ops = append(ops, &txnbuild.SetOptions{
				LowThreshold:    txnbuild.NewThreshold(txnbuild.Threshold(0)),
				MediumThreshold: txnbuild.NewThreshold(txnbuild.Threshold(0)),
				HighThreshold:   txnbuild.NewThreshold(txnbuild.Threshold(0)),
				SourceAccount:   wallet.ID,
			})
		}
	}

	if len(ops) == 0 {
		// no operations to sign
		return "no-ops", "", messages, walletMustSign, multipartySign, nil
	}
	var tx *txnbuild.Transaction
	if numberOfApprovers > 0 {
		//multiparty
		transactionSource = chanSourceAccount.AccountID
		tx, err = txnbuild.NewTransaction(
			txnbuild.TransactionParams{
				SourceAccount:        chanSourceAccount,
				IncrementSequenceNum: true,
				Operations:           ops,
				BaseFee:              txnbuild.MinBaseFee,
				Preconditions: txnbuild.Preconditions{
					TimeBounds: txnbuild.NewInfiniteTimeout(),
				},
				Memo: txnbuild.MemoText("Disable shared access"),
			},
		)
	} else {
		//single signer
		tx, err = txnbuild.NewTransaction(
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
	}

	if err != nil {
		log.Println("[generateRemoveSharedAccessXdr] error constructing transaction ", err)
		return "", "", messages, walletMustSign, multipartySign, err
	}

	if numberOfApprovers > 0 {
		tx, err = tx.Sign(gc.BantuNetworkPassphrase, chanAccount)
		if err != nil {
			log.Println("[generateModifySharedAccessXdr] error signing transaction with chan account", err)
			return "", "", messages, walletMustSign, multipartySign, err
		}
	}
	var xdrBase64 string

	xdrBase64, err = tx.Base64()
	if err != nil {
		log.Println("[generateRemoveSharedAccessXdr] error getting txn base64", err)
		return "", "", messages, walletMustSign, multipartySign, err
	}

	return xdrBase64, transactionSource, messages, walletMustSign, multipartySign, nil

}

func generateRemoveSharedAccessOps(wallet *userModels.UserWallet, walletOwner *userModels.User, approver *userModels.User, gc *sharedconfig.GlobalConfig) (op txnbuild.Operation, err error) {
	client := gc.BantuExpansionClient

	//check if primary account has native enough native balance
	var nativeAsset txnbuild.Asset = txnbuild.NativeAsset{}
	_, _, _, _, walletSourceAccount, errWalletAct := network.BlockchainAccountProperties(client, wallet.ID, nativeAsset)
	if errWalletAct != nil {
		log.Printf("[generateRemoveSharedAccessXdr] by [%v] for shared Account Properties error:[%v] \n", wallet.Alias, errWalletAct)

		return op, errWalletAct
	}

	//ensure u r using the account signer, since the account may have been recovered, or may be recovered in the future, changing the signer, but retaining the primary key
	approverAccountExists, _, _, _, _, _ := network.BlockchainAccountProperties(client, approver.PrimarySigner, nativeAsset)

	if approverAccountExists {
		//account exists, check if it already it a signer in the wallet

		//remove signer if already a signer
		if wallet.SignerIsValidWA(approver.PrimarySigner, walletSourceAccount) {
			if walletOwner.PrimarySigner != wallet.Signer {

				op = &txnbuild.SetOptions{
					Signer: &txnbuild.Signer{
						Address: approver.PrimarySigner,
						Weight:  0,
					},
					SourceAccount: wallet.ID,
				}

				return op, nil
			} else {
				//primary signer and master signer
				op = &txnbuild.SetOptions{
					MasterWeight:  txnbuild.NewThreshold(0),
					SourceAccount: wallet.ID,
				}

				return op, nil
			}
		}

	}
	return op, &tErrors.ErrorTemporaryServerError{}
}

func generateRemoveRecoveredAccountAccessOps(wallet *userModels.UserWallet, approverUsernameAdded string, gc *sharedconfig.GlobalConfig) (ops []txnbuild.Operation) {
	client := gc.BantuExpansionClient
	listOfRecovery := make([]userModels.UserAccountRecoveryLog, 0)
	gc.DB.Where("username = ?", approverUsernameAdded).Find(&listOfRecovery)
	ops = make([]txnbuild.Operation, 0)
	if len(listOfRecovery) == 0 || listOfRecovery == nil {
		//no account recovery done so far
		return
	}
	//check if primary account has native enough native balance
	var nativeAsset txnbuild.Asset = txnbuild.NativeAsset{}
	_, _, _, _, walletSourceAccount, errWalletAct := network.BlockchainAccountProperties(client, wallet.ID, nativeAsset)
	if errWalletAct != nil {
		log.Printf("[generateRemoveRecoveredAccountAccessOps] by [%v] for shared Account Properties error:[%v] \n", wallet.Alias, errWalletAct)

	}

	for _, aRec := range listOfRecovery {
		//ensure u r using the account signer, since the account may have been recovered, or may be recovered in the future, changing the signer, but retaining the primary key
		approverAccountExists, _, _, _, _, _ := network.BlockchainAccountProperties(client, aRec.OldSignerPublicKey, nativeAsset)

		if approverAccountExists {
			//account exists, check if it already it a signer in the wallet

			//remove signer if already a signer
			if wallet.SignerIsValidWA(aRec.OldSignerPublicKey, walletSourceAccount) {
				if aRec.MasterWallet == 1 {
					//primary signer and master signer
					ops = append(ops, &txnbuild.SetOptions{
						MasterWeight:  txnbuild.NewThreshold(0),
						SourceAccount: wallet.ID,
					})

				} else {

					ops = append(ops, &txnbuild.SetOptions{
						Signer: &txnbuild.Signer{
							Address: aRec.OldSignerPublicKey,
							Weight:  0,
						},
						SourceAccount: wallet.ID,
					})

				}
			}

		}
	}

	return ops
}

func HasAccessToPublicKey(signerPublicKey, targetPublicKey string, gc *sharedconfig.GlobalConfig) (hasAccess bool) {
	signerUser, err := usersDB.GetUserFromPrimarySigner(signerPublicKey, gc.DB, gc)

	if err != nil {
		return false
	}
	wallet, temp, err := usersDB.GetWallet(targetPublicKey, gc.DB)
	if err != nil {
		return false
	}
	if temp {
		return false
	}
	if wallet.SharedAccessEnabled == 0 {
		return false
	}
	return wallet.SignerHasAccess(&signerUser, gc)

}

func HasInitiatorPermissionToPublicKey(ownerSignerPublicKey, targetPublicKey string, gc *sharedconfig.GlobalConfig) (hasAccess bool) {
	user, err := usersDB.GetUserFromPrimarySigner(ownerSignerPublicKey, gc.DB, gc)

	if err != nil {
		return false
	}
	walletPermissions := user.WalletsSharedWithUser
	if len(walletPermissions) == 0 {
		return false
	}
	for _, walletAccess := range walletPermissions {
		if walletAccess.WalletPublicKey == targetPublicKey && walletAccess.Permission == "INITIATOR" {
			return true
		}
	}

	return false
}

func SignerHasInitiatorPermissionToPublicKey(signerOwner userModels.User, targetPublicKey string, gc *sharedconfig.GlobalConfig) (hasAccess bool) {

	walletPermissions := signerOwner.WalletsSharedWithUser
	if len(walletPermissions) == 0 {
		return false
	}
	for _, walletAccess := range walletPermissions {
		if walletAccess.WalletPublicKey == targetPublicKey && walletAccess.Permission == "INITIATOR" {
			return true
		}
	}

	return false
}
