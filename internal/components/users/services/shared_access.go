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
	if len(accessInfo.Permissions) == 1 && accessInfo.Permissions[0].TargetUsername == walletOwner.Username {
		return returnedWallet, &tErrors.CustomError{
			Param:      "permissions",
			Err:        "error-permission-unacceptable-access",
			ErrMessage: "Shared access cannot be enabled with view-only permission granted to yourself.",
			Code:       http.StatusForbidden,
		}
	}

	if wallet.WalletType == 2 || wallet.WalletType == 3 {
		return returnedWallet, &tErrors.CustomError{
			Param:      "publicKey",
			Err:        "error-wallet-type-not-allowed",
			ErrMessage: "Market Making & Bulk Payment wallets are not allowed for this operation.",
			Code:       http.StatusForbidden,
		}
	}

	var accessListInfo, viewOnly []userModels.WalletPermissionInfo
	var accessList []userModels.WalletPermission
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
			if e != nil {
				return returnedWallet, &tErrors.CustomError{
					Param:      "username",
					Err:        "error-trovo-wallet-account-invalid",
					ErrMessage: fmt.Sprintf("Trovo wallet account [%v] is not Trovo wallet account username.", v.TargetUsername),
					Code:       http.StatusForbidden,
				}
			}
		}
		name := fmt.Sprintf("%v", u.FirstName)
		if u.LastName == nil {
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
	}
	if len(viewOnly) > 0 {
		message := "Unnecessary VIEW-ONLY access for these accounts where removed:"
		for _, v := range viewOnly {
			for _, a := range accessInfo.Permissions {
				if v.TargetUsername == a.TargetUsername && (a.Permission == "APPROVER" || a.Permission == "INITIATOR") {
					// remove the view only since the approver and initiator has view access already

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
		accessInfo.Messages = append(accessInfo.Messages, message)
	}

	if numberOfSubmittedApprovers <= accessInfo.NumberOfApprovalsNeeded && accessInfo.NumberOfApprovalsNeeded > 1 {
		//number of authorizers does not reach the minimum threshold needed. cannot proceed so as to prevent account lockout
		return returnedWallet, &tErrors.CustomError{
			Param:      "numberOfApprovalsNeeded",
			Err:        "error-approvers-not-enough",
			ErrMessage: fmt.Sprintf("Please ensure that your list of approvers are greater than the minimum number required to approve a transaction [%v]", accessInfo.NumberOfApprovalsNeeded),
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

	errDB = dbTX.Save(wallet).Error
	if err != nil {
		log.Printf("[CreateSharedWalletAccess] error saving shared access status of the wallet:%v\n sharedAccess:%+v\n", errDB, wallet)
		return returnedWallet, &tErrors.ErrorTemporaryServerError{}
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

	txnHash, err := network.SubmitXdrWithSignature(gc.BantuExpansionClient, wallet.Signer, xdrBase64, accessInfo.TransactionSignature)
	if err != nil {
		log.Printf("Error submitting shared access txn [%+v] transaction: %s\n", accessInfo, err.Error())
		// logDiscordFailedRecovery(fmt.Sprintf("Error submitting shared access txn [%+v] transaction: %s", accessInfo, err.Error()))
		return returnedWallet, &tErrors.ErrorTemporaryServerError{}
	}
	accessInfo.TransactionID = txnHash

	dbTX.Commit()
	return returnedWallet, nil
}

func ModifySharedWalletAccess(signerUser *userModels.User, walletOwner *userModels.User, wallet *userModels.UserWallet, accessInfo *userModels.ModifySharedAccessInfo, gc *sharedconfig.GlobalConfig) (revokedList, modifiedList, addedList []userModels.WalletPermission, err error) {
	// var  userModels.UserWalletSharedAccess
	//prepare database execution
	dbTX := gc.DB.Begin()
	defer dbTX.Rollback()
	if wallet.WalletType == 2 || wallet.WalletType == 3 {
		err = &tErrors.CustomError{
			Param:      "publicKey",
			Err:        "error-wallet-type-not-allowed",
			ErrMessage: "Market Making & Bulk Payment wallets are not allowed for this operation.",
			Code:       http.StatusForbidden,
		}
		return
	}
	if len(accessInfo.ModifiedPermissions) == 0 && len(accessInfo.AddedPermissions) == 0 && len(accessInfo.RevokedPermissions) == 0 {
		err = &tErrors.CustomError{
			Param:      "permissions",
			Err:        "error-no-operations-to-perform",
			ErrMessage: "No operation can be performed since no permissions to add or revoke or modify",
			Code:       http.StatusBadRequest,
		}
		return
	}
	viewOnly := make(map[string]string, 0)
	var oldNumberOfApprovers int
	for _, perm := range wallet.Permissions {
		if perm.Permission == "APPROVER" {
			oldNumberOfApprovers++
		}
	}
	ops := make([]txnbuild.Operation, 0)
	accessInfo.Messages = make([]string, 0)
	var revokedListInfo, modifiedListInfo, addedListInfo []userModels.WalletPermissionInfo

	var numberOfSubmittedApprovers int
	var numberOfSubmittedInitiators int

	walletID := userModels.UserWalletID(accessInfo.WalletPublicKey)
	if wallet.SharedAccessEnabled == 0 {
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

			//it is a subwallet and cannot be given access
			err = &tErrors.CustomError{
				Param:      "username",
				Err:        "error-subwallet-not-allowed",
				ErrMessage: fmt.Sprintf("Access can only be granted/revoked to/from trovo wallet account, not a subwallet [%v]", v.TargetUsername),
				Code:       http.StatusForbidden,
			}
			return
		}
		userPermission, e := walletID.GetUserPermissionOnWallet(v.TargetUsername, dbTX)
		if e != nil {
			continue
		}

		u, e := userModels.Username(v.TargetUsername).GetSimpleUser(gc.DB, gc)
		if e != nil {
			err = &tErrors.CustomError{
				Param:      "username",
				Err:        "error-trovo-wallet-account-invalid",
				ErrMessage: fmt.Sprintf("Trovo wallet account [%v] could not be validated at this time.", v.TargetUsername),
				Code:       http.StatusForbidden,
			}
			return
		}
		name := fmt.Sprintf("%v", u.FirstName)
		if u.LastName == nil {
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
			CreatedAt:       userPermission.CreatedAt,
			UpdatedAt:       userPermission.UpdatedAt,
			ID:              userPermission.ID,
			TargetUsername:  v.TargetUsername,
			Permission:      v.Permission,
			WalletPublicKey: wallet.ID,
		})
		if v.Permission == "APPROVER" {
			//get ops to add.
			op, e := generateRemoveSharedAccessOps(wallet, walletOwner, &u, gc)
			if e == nil {
				ops = append(ops, op)
			}
		}

	}

	//modified permissions
	for _, v := range accessInfo.ModifiedPermissions {

		v.TargetUsername = strings.ToLower(v.TargetUsername)
		if strings.Contains(v.TargetUsername, "_") {

			//it is a subWallet and cannot be given access
			err = &tErrors.CustomError{
				Param:      "username",
				Err:        "error-subwallet-not-allowed",
				ErrMessage: fmt.Sprintf("Access can only be granted to trovo wallet account, not a subwallet [%v]", v.TargetUsername),
				Code:       http.StatusForbidden,
			}
			return
		}
		ePermission, e := walletID.GetUserPermissionOnWallet(v.TargetUsername, dbTX)
		if e != nil {
			continue
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
			err = &tErrors.CustomError{
				Param:      "username",
				Err:        "error-trovo-wallet-account-invalid",
				ErrMessage: fmt.Sprintf("Trovo wallet account [%v] could not be validated at this time.", v.TargetUsername),
				Code:       http.StatusForbidden,
			}
			return
		}

		name := fmt.Sprintf("%v", u.FirstName)
		if u.LastName == nil {
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

		// v.permission is the new peremission
		if v.Permission == "APPROVER" {
			// attempt to add the public key as signer
			o, m, e := generateAddSharedAccessOps(wallet, walletOwner, &u, gc)
			if e == nil {
				ops = append(ops, o...)
				accessInfo.Messages = append(accessInfo.Messages, m...)
			}

		}
		// if is a downgrade of access from approver
		if ePermission.Permission == "APPROVER" {
			// attempt to remove the public key as signer
			op, e := generateRemoveSharedAccessOps(wallet, walletOwner, &u, gc)
			if e == nil {
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

			//it is a subWallet and cannot be given access
			err = &tErrors.CustomError{
				Param:      "username",
				Err:        "error-subwallet-not-allowed",
				ErrMessage: fmt.Sprintf("Access can only be granted to trovo wallet account, not a subwallet [%v]", v.TargetUsername),
				Code:       http.StatusForbidden,
			}
			return
		}
		if v.TargetUsername == walletOwner.Username && v.Permission == "VIEW-ONLY" {
			//skip owners being added as view only access. Owners have view access by default.
			accessInfo.Messages = append(accessInfo.Messages, fmt.Sprintf("%v already has view access as owner, so the assigned View access has been skipped.", v.TargetUsername))
			continue
		}
		if _, ok := checkAccess[v.TargetUsername+v.Permission]; ok {
			continue
		}
		_, e := walletID.GetUserPermissionOnWallet(v.TargetUsername, dbTX)
		if e == nil {
			// permission exists, so cannot be added
			continue
		}

		permissionID := uuid.NewString()

		u, e := userModels.Username(v.TargetUsername).GetSimpleUser(gc.DB, gc)
		if e != nil {
			err = &tErrors.CustomError{
				Param:      "username",
				Err:        "error-trovo-wallet-account-invalid",
				ErrMessage: fmt.Sprintf("Trovo wallet account [%v] could not be validated at this time.", v.TargetUsername),
				Code:       http.StatusForbidden,
			}
			return
		}

		name := fmt.Sprintf("%v", u.FirstName)
		if u.LastName == nil {
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
		if v.Permission == "APPROVER" {
			// attempt to add the public key as signer
			o, m, e := generateAddSharedAccessOps(wallet, walletOwner, &u, gc)
			if e == nil {
				ops = append(ops, o...)
				accessInfo.Messages = append(accessInfo.Messages, m...)
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
		}
		if len(modifiedList) > 0 {
			e := dbTX.Save(&modifiedList).Error
			if e != nil {
				log.Println("[ModifySharedWalletAccess] error saving modified list: ", e)
				err = &tErrors.ErrorTemporaryServerError{}
				return
			}
		}
		if len(addedList) > 0 {
			e := dbTX.Create(&addedList).Error
			if e != nil {
				log.Println("[ModifySharedWalletAccess] error creating added permissions: ", e)
				err = &tErrors.ErrorTemporaryServerError{}
				return
			}
		}

	}
	//invalidate all existing cache relating to this wallet
	wallet.InvalidateUserCache(gc)
	//it was successfully saved. now refresh the list to know the standing.
	updatedWallet, e := walletID.GetWallet(dbTX)
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
			ErrMessage: fmt.Sprintf("Please ensure that your list of approvers are greater than the minimum number required to approve a transaction [%v]", accessInfo.NumberOfApprovalsNeeded),
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

	errDB := dbTX.Save(&updatedWallet).Error
	if errDB != nil {
		log.Printf("[ModifySharedWalletAccess] error saving shared access status of the wallet:%v\n sharedAccess:%+v\n", errDB, wallet)
		err = &tErrors.ErrorTemporaryServerError{}
		return
	}

	xdrBase64, messages, errGenXdr := generateModifySharedAccessXdr(wallet, walletOwner, numberOfSubmittedApprovers, accessInfo.NumberOfApprovalsNeeded, oldNumberOfApprovers, ops, gc)
	if errGenXdr != nil {
		err = errGenXdr
		return
	}

	accessInfo.Messages = append(accessInfo.Messages, messages...)

	accessInfo.NetworkPassPhrase = network.GetBlockchainNetworkPassPhrase()

	accessInfo.Transaction = xdrBase64

	// accessInfo.RevokedPermissions = revokedListInfo
	// accessInfo.ModifiedPermissions = modifiedListInfo
	// accessInfo.AddedPermissions = addedListInfo

	if oldNumberOfApprovers > 0 {
		accessInfo.MultiParty = 1
	}
	accessInfo.RevokedPermissions = revokedListInfo
	accessInfo.ModifiedPermissions = modifiedListInfo
	accessInfo.AddedPermissions = addedListInfo
	if len(accessInfo.TransactionSignature) == 0 && accessInfo.Commit == 0 {
		err = nil
		return
	}
	if len(accessInfo.TransactionSignature) > 0 && oldNumberOfApprovers == 0 && accessInfo.Commit == 1 {
		// extract signature and submit transaction
		//submit to blockchain
		var txnHash string
		txnHash, err = network.SubmitXdrWithSignature(gc.BantuExpansionClient, wallet.Signer, xdrBase64, accessInfo.TransactionSignature)
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
				revokedUsers = fmt.Sprintf("%s|%v", u.TargetUsername, u.Permission)
				if i < len(revokedList)-1 {
					revokedUsers = fmt.Sprintf("%s, ", revokedUsers)
				}
			}
			description = fmt.Sprintf("%s\nPermissions to revoke :%v", description, revokedUsers)
		}

		if len(modifiedList) > 0 {
			modifiedUsers := ""
			for i, u := range modifiedList {
				modifiedUsers = fmt.Sprintf("%s|%v", u.TargetUsername, u.Permission)
				if i < len(revokedList)-1 {
					modifiedUsers = fmt.Sprintf("%s, ", modifiedUsers)
				}
			}
			description = fmt.Sprintf("%s\nModifying Permissions :%v", description, modifiedUsers)
		}

		if len(addedList) > 0 {
			addedUsers := ""
			for i, u := range addedList {
				addedUsers = fmt.Sprintf("%s|%v", u.TargetUsername, u.Permission)
				if i < len(addedList)-1 {
					addedUsers = fmt.Sprintf("%s, ", addedUsers)
				}
			}
			description = fmt.Sprintf("%s.\nAdding New Permissions :%v.", description, addedUsers)
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

	accessList := wallet.Permissions
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
	userPermissions := make([]string, 0)
	walletOwner, e := wallet.GetWalletOwner(gc.DB, gc)
	if e != nil {
		return &tErrors.CustomError{
			Param:      "username",
			Err:        "error-confirming-wallet-owner",
			ErrMessage: "Unable to confirm wallet owner at this time. Please try again after some minutes.",
			Code:       http.StatusForbidden,
		}
	}

	for _, v := range accessList {
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
		userPermissions = append(userPermissions, fmt.Sprintf("%s - (%v)", v.TargetUsername, v.Permission))

	}

	// // if approver exists, then owner must sign transaction to add them as signers
	// if (numberOfApprovers - selfApprover) > 0 {
	// 	accessInfo.SignatureRequired = 1

	// }

	xdrBase64, messages, walletMustSign, _, errGenXdr := generateRemoveSharedAccessXdr(wallet, &walletOwner, approverUsers, numberOfApprovers, gc)
	if errGenXdr != nil {
		return errGenXdr
	}
	accessInfo.Messages = messages
	if walletMustSign {
		accessInfo.SignatureRequired = 1

	}
	accessInfo.NetworkPassPhrase = network.GetBlockchainNetworkPassPhrase()

	accessInfo.Transaction = xdrBase64
	if numberOfApprovers > 0 && len(accessInfo.TransactionSignature) == 0 {
		accessInfo.MultiParty = 1
		if accessInfo.Commit == 0 {
			return nil
		}

	}

	if len(accessInfo.TransactionSignature) == 0 {
		if accessInfo.Commit == 0 {
			return nil
		}
	}

	if numberOfApprovers > 0 && len(accessInfo.TransactionSignature) > 0 {
		//SET transaction id to pending auth
		accessInfo.MultiParty = 1
		if accessInfo.Commit == 1 {
			accessInfo.TransactionID = "PENDING_AUTH"

			//TODO: queue transaction and notify signers
			id := uuid.New().String()
			description := fmt.Sprintf("Disabling shared access on wallet %v.\nThis will remove the permissions:\n%v", wallet.Alias, userPermissions)
			transactionByte, _ := json.Marshal(*accessInfo)
			transactionStr := string(transactionByte)
			pendingAuth := userModels.PendingAuth{
				ID:                       id,
				Initiator:                signerUser.Username,
				InitiatorSignerPublicKey: signerUser.PrimarySigner,
				WalletPublicKey:          wallet.ID,
				TransactionType:          "DISABLE SHARED ACCESS",
				Description:              description,
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

	if accessInfo.Commit == 0 {
		return nil
	}
	dbTX := gc.DB.Begin()
	defer dbTX.Rollback()
	wallet.SharedAccessEnabled = 0
	wallet.NumberOfApprovalsNeeded = 0
	e = dbTX.Delete(&accessList).Error
	if e != nil {
		log.Println("[RemoveSharedWalletAccess] error deleting access list", e)
		return &tErrors.ErrorTemporaryServerError{}
	}
	// extract signature and submit transaction
	//getting here means it does not contain approvers
	//submit to blockchain

	txnHash, err := network.SubmitXdrWithSignature(gc.BantuExpansionClient, wallet.Signer, xdrBase64, accessInfo.TransactionSignature)
	if err != nil {
		log.Printf("[RemoveSharedWalletAccess] Error submitting disable shared access txn [%+v] transaction: %s\n", accessInfo, err.Error())
		// logDiscordFailedRecovery(fmt.Sprintf("Error submitting shared access txn [%+v] transaction: %s", accessInfo, err.Error()))
		return &tErrors.ErrorTemporaryServerError{}
	}
	accessInfo.TransactionID = txnHash

	dbTX.Commit()
	return nil

}

func generateCreateSharedAccessXdr(wallet *userModels.UserWallet, walletOwner *userModels.User, approvers []*userModels.User, accessInfo []userModels.WalletPermissionInfo, authThreshold int, gc *sharedconfig.GlobalConfig) (xdrbase64 string, messages []string, walletMustSign bool, err error) {
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
			totalNativeBalanceNeeded := activationAmount.Mul(decimal.NewFromInt(totalUsersToFund))
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

	//activating shared access is free. No fee, except for view only access.
	//TODO: if account exists and subwallet has enough balance, we add the operation to pay TROVO fee from primary Wallet
	serviceFee, e := decimal.NewFromString(os.Getenv("SHARED_ACCESS_FEE_AMOUNT"))
	if e != nil {
		serviceFee = decimal.Zero
	}
	if serviceFee.IsPositive() {
		if len(approvers) == 0 {
			//process service fee
			if len(os.Getenv("SHARED_ACCESS_FEE_ASSET_ISSUER")) == 56 {
				ops = append(ops, &txnbuild.Payment{
					Destination:   os.Getenv("SHARED_ACCESS_FEE_ADDRESS"),
					Amount:        os.Getenv("SHARED_ACCESS_FEE_AMOUNT"),
					SourceAccount: wallet.ID,
					Asset:         txnbuild.CreditAsset{Code: os.Getenv("SHARED_ACCESS_FEE_ASSET_CODE"), Issuer: os.Getenv("SHARED_ACCESS_FEE_ASSET_ISSUER")},
				})
				messages = append(messages, fmt.Sprintf("%v %v will be deducted from wallet %v as service fee for creating view only access.", os.Getenv("SHARED_ACCESS_FEE_AMOUNT"), os.Getenv("SHARED_ACCESS_FEE_ASSET_CODE"), wallet.Alias))

			} else {
				ops = append(ops, &txnbuild.Payment{
					Destination:   os.Getenv("SHARED_ACCESS_FEE_ADDRESS"),
					Amount:        os.Getenv("SHARED_ACCESS_FEE_AMOUNT"),
					SourceAccount: wallet.ID,
					Asset:         txnbuild.NativeAsset{},
				})
				messages = append(messages, fmt.Sprintf("%v %v will be deducted from wallet %v as service fee for creating view only access.", os.Getenv("SHARED_ACCESS_FEE_AMOUNT"), os.Getenv("NATIVE_ASSET_CODE"), wallet.Alias))

			}

		}
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
		// no operations to sign. create a dummy ops, will be ignored on next try.
		ops = append(ops, &txnbuild.SetOptions{
			LowThreshold:    txnbuild.NewThreshold(txnbuild.Threshold(0)),
			MediumThreshold: txnbuild.NewThreshold(txnbuild.Threshold(0)),
			HighThreshold:   txnbuild.NewThreshold(txnbuild.Threshold(0)),
			SourceAccount:   wallet.ID,
		})
		walletMustSign = true
		// return "no-ops", messages, walletMustSign, nil
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

func generateModifySharedAccessXdr(wallet *userModels.UserWallet, walletOwner *userModels.User, numberOfSubmittedApprovers, numberOfApprovalsNeeded, oldNumberOfApprovers int, ops []txnbuild.Operation, gc *sharedconfig.GlobalConfig) (xdrbase64 string, messages []string, err error) {
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
		return "", messages, err
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

		return "", messages, errWalletAct
	}
	if !walletAccountExists || (walletAccountNativeBalance.Sub(activationAmount)).LessThan(minBalance.Mul(decimal.NewFromInt(int64(numberOfSubmittedApprovers)))) {
		log.Printf("[generateModifySharedAccessXdr] by [%v] shared WalletAccount underfunded. Needs at least %v %v\n", wallet.Alias, (minBalance.Mul(decimal.NewFromInt(int64(numberOfSubmittedApprovers)))).Truncate(7).String(), os.Getenv("NATIVE_ASSET_CODE"))

		err = &tErrors.CustomError{
			Param:      "publicKey",
			Err:        "error-wallet-underfunded",
			ErrMessage: fmt.Sprintf("Wallet %v needs minimum of %v %v balance to perform this operation.", wallet.Alias, (minBalance.Mul(decimal.NewFromInt(int64(numberOfSubmittedApprovers)))).Truncate(7).String(), os.Getenv("NATIVE_ASSET_CODE")),
			Code:       http.StatusBadRequest,
		}
		return "", messages, err
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
			}
		}
	}

	//TODO: if account exists and subwallet has enough balance, we add the operation to pay TROVO fee from primary Wallet

	serviceFee, e := decimal.NewFromString(os.Getenv("SHARED_ACCESS_FEE_AMOUNT"))
	if e != nil {
		serviceFee = decimal.Zero
	}
	if serviceFee.IsPositive() {
		{
			//process service fee
			if len(os.Getenv("SHARED_ACCESS_FEE_ASSET_ISSUER")) == 56 {
				ops = append(ops, &txnbuild.Payment{
					Destination:   os.Getenv("SHARED_ACCESS_FEE_ADDRESS"),
					Amount:        os.Getenv("SHARED_ACCESS_FEE_AMOUNT"),
					SourceAccount: wallet.ID,
					Asset:         txnbuild.CreditAsset{Code: os.Getenv("SHARED_ACCESS_FEE_ASSET_CODE"), Issuer: os.Getenv("SHARED_ACCESS_FEE_ASSET_ISSUER")},
				})
				messages = append(messages, fmt.Sprintf("%v %v will be deducted from wallet %v as service fee.", os.Getenv("SHARED_ACCESS_FEE_AMOUNT"), os.Getenv("SHARED_ACCESS_FEE_ASSET_CODE"), wallet.Alias))

			} else {
				ops = append(ops, &txnbuild.Payment{
					Destination:   os.Getenv("SHARED_ACCESS_FEE_ADDRESS"),
					Amount:        os.Getenv("SHARED_ACCESS_FEE_AMOUNT"),
					SourceAccount: wallet.ID,
					Asset:         txnbuild.NativeAsset{},
				})
				messages = append(messages, fmt.Sprintf("%v %v will be deducted from wallet %v as service fee.", os.Getenv("SHARED_ACCESS_FEE_AMOUNT"), os.Getenv("NATIVE_ASSET_CODE"), wallet.Alias))

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
		}
	}
	// if len(ops) == 0 {
	// 	// no operations to sign. create a dummy ops, will be ignored on next try.
	// 	ops = append(ops, &txnbuild.SetOptions{
	// 		LowThreshold:    txnbuild.NewThreshold(txnbuild.Threshold(0)),
	// 		MediumThreshold: txnbuild.NewThreshold(txnbuild.Threshold(0)),
	// 		HighThreshold:   txnbuild.NewThreshold(txnbuild.Threshold(0)),
	// 		SourceAccount:   wallet.ID,
	// 	})
	// }

	// Construct the transaction that holds the operations to execute on the network
	var tx *txnbuild.Transaction
	if oldNumberOfApprovers > 0 {
		//multiparty
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
				BaseFee:              txnbuild.MinBaseFee,
				Preconditions: txnbuild.Preconditions{
					TimeBounds: txnbuild.NewInfiniteTimeout(),
				},
				Memo: txnbuild.MemoText("Modify shared access"),
			},
		)
	}

	if err != nil {
		log.Println("[generateModifySharedAccessXdr] error constructing transaction ", err)
		return "", messages, err
	}
	if oldNumberOfApprovers > 0 {
		tx, err = tx.Sign(gc.BantuNetworkPassphrase, chanAccount)
		if err != nil {
			log.Println("[generateModifySharedAccessXdr] error signing transaction with chan account", err)
			return "", messages, err
		}
	}

	var xdrBase64 string

	xdrBase64, err = tx.Base64()
	if err != nil {
		log.Println("[generateModifySharedAccessXdr] error getting txn base64", err)
		return "", messages, err
	}

	return xdrBase64, messages, nil

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

			ops = append(ops, &txnbuild.SetOptions{
				Signer: &txnbuild.Signer{
					Address: approver.PrimarySigner,
					Weight:  1,
				},
				SourceAccount: wallet.ID,
			})

		}

	}

	return ops, messages, nil

}

func generateRemoveSharedAccessXdr(wallet *userModels.UserWallet, walletOwner *userModels.User, approvers []*userModels.User, numberOfApprovers int, gc *sharedconfig.GlobalConfig) (xdrbase64 string, messages []string, walletMustSign, multipartySign bool, err error) {
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

		return "", messages, walletMustSign, multipartySign, errWalletAct
	}
	if !walletAccountExists || (walletAccountNativeBalance.Sub(activationAmount)).LessThan(minBalance) {
		log.Printf("[generateRemoveSharedAccessXdr] by [%v] shared WalletAccount underfunded \n", wallet.Alias)

		err = &tErrors.CustomError{
			Param:      "publicKey",
			Err:        "error-wallet-underfunded",
			ErrMessage: fmt.Sprintf("Wallet %v does not have enough XBN balance to perform this operation", wallet.Alias),
			Code:       404,
		}
		return "", messages, walletMustSign, multipartySign, err
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
			}
		}
		// if walletOwner.AccountRecoveryEnabled == 1 && PublicKeyHasViewOnlyAccess(wallet.ID, gc){
		// 	// get the recovery keypair
		// 	recoveryAddress := bc.GetRecoveryAccountAddress(walletOwner.Username, walletOwner.PublicKey)

		// 	if !userBc.SignerIsValid(wallet.ID, recoveryAddress) {
		// 		//recovery a signer to the wallet. remove it
		// 		ops = append(ops, &txnbuild.SetOptions{
		// 			Signer: &txnbuild.Signer{
		// 				Address: recoveryAddress,
		// 				Weight:  1,
		// 			},
		// 			SourceAccount: wallet.ID,
		// 		})

		// 		//add message about disabling recovery on that wallet
		// 		messages = append(messages, "Account Recovery on this wallet has been enabled.")
		// 		walletMustSign = true
		// 	}
		// }
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
					ErrMessage: fmt.Sprintf("Wallet %v needs more than %v XBN balance to perform this operation", wallet.Alias, totalNativeBalanceNeeded.String()),
					Code:       404,
				}
				return "", messages, walletMustSign, multipartySign, err
			}
		}
	}

	//TODO: if account exists and subwallet has enough balance, we add the operation to pay TROVO fee from primary Wallet
	fee := decimal.RequireFromString(os.Getenv("SHARED_ACCESS_FEE_AMOUNT"))
	if !fee.IsZero() {
		{
			//process service fee
			if len(os.Getenv("SHARED_ACCESS_FEE_ASSET_ISSUER")) != 56 {
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
		return "no-ops", messages, walletMustSign, multipartySign, nil
	}
	var tx *txnbuild.Transaction
	if numberOfApprovers > 0 {
		//multiparty
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
		return "", messages, walletMustSign, multipartySign, err
	}

	if numberOfApprovers > 0 {
		tx, err = tx.Sign(gc.BantuNetworkPassphrase, chanAccount)
		if err != nil {
			log.Println("[generateModifySharedAccessXdr] error signing transaction with chan account", err)
			return "", messages, walletMustSign, multipartySign, err
		}
	}
	var xdrBase64 string

	xdrBase64, err = tx.Base64()
	if err != nil {
		log.Println("[generateRemoveSharedAccessXdr] error getting txn base64", err)
		return "", messages, walletMustSign, multipartySign, err
	}

	return xdrBase64, messages, walletMustSign, multipartySign, nil

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
		if wallet.SignerIsValidWA(approver.PrimarySigner, walletSourceAccount) && walletOwner.PrimarySigner != wallet.Signer {

			op = &txnbuild.SetOptions{
				Signer: &txnbuild.Signer{
					Address: approver.PrimarySigner,
					Weight:  0,
				},
				SourceAccount: wallet.ID,
			}

			return op, nil
		}

	}
	return op, &tErrors.ErrorTemporaryServerError{}
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
