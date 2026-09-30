package users

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"
	"trovo-wallet-api/internal/basetxn"
	usersDB "trovo-wallet-api/internal/components/users/db"
	userModels "trovo-wallet-api/internal/components/users/models"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/network"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/google/uuid"
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

func AddressCountViewOnlyAccess(publicKey string, gc *sharedconfig.GlobalConfig) (accessCount uint) {
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

func AddressHasViewOnlyAccess(publicKey string, gc *sharedconfig.GlobalConfig) (viewOnly bool) {
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

func AddressHasViewOnlyAccessWACL(publicKey string, accessList []userModels.WalletPermissionInfo, gc *sharedconfig.GlobalConfig) (viewOnly bool) {
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

func AddressCountApproverAccessWACL(publicKey string, accessList []userModels.WalletPermissionInfo, gc *sharedconfig.GlobalConfig) (accessCount uint) {
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

func AddressCountInitiatorAccessWACL(publicKey string, accessList []userModels.WalletPermissionInfo, gc *sharedconfig.GlobalConfig) (accessCount uint) {
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

func AddressCountApproverAccess(publicKey string, gc *sharedconfig.GlobalConfig) (accessCount uint) {
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

func AddressCountInitiatorAccess(publicKey string, gc *sharedconfig.GlobalConfig) (accessCount uint) {
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
	if sErr := signerUser.EnsureNotSuspended(); sErr != nil {
		return returnedWallet, sErr
	}
	if sErr := walletOwner.EnsureNotSuspended(); sErr != nil {
		return returnedWallet, sErr
	}

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
	if wallet.WalletType == 1 && wallet.LinkedWalletAddress != nil {
		// set the linked wallet if it is a tokenization wallet
		hasLinkedWallet = true
		// accessInfo.LinkedWalletSignatureRequired = 1
		// accessInfo.LinkedWalletAddress = *wallet.LinkedWalletAddress
		linkedWallet, err = userModels.UserWalletID(*wallet.LinkedWalletAddress).GetWallet(gc.DB, gc)
		if err != nil {
			return returnedWallet, &tErrors.CustomError{
				Param:      "linkedWalletAddress",
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
		if !AddressHasViewOnlyAccessWACL(accessInfo.WalletAddress, accessInfo.Permissions, gc) {
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
				ErrMessage: fmt.Sprintf("Access can only be granted to TrovoApp account, not a subwallet [%v]", v.TargetUsername),
				Code:       http.StatusForbidden,
			}
		}
		u, e := userModels.Username(v.TargetUsername).GetSimpleUser(gc.DB, gc)
		if e != nil {
			return returnedWallet, &tErrors.CustomError{
				Param:      "username",
				Err:        "error-trovo-wallet-account-invalid",
				ErrMessage: fmt.Sprintf("TrovoApp account [%v] could not be validated at this time.", v.TargetUsername),
				Code:       http.StatusForbidden,
			}
		}
		if u.Username != v.TargetUsername {

			return returnedWallet, &tErrors.CustomError{
				Param:      "username",
				Err:        "error-trovo-wallet-account-invalid",
				ErrMessage: fmt.Sprintf("TrovoApp account [%v] is not TrovoApp account username.", v.TargetUsername),
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
			WalletAddress:         wallet.ID,
			WalletAlias:           wallet.Alias,
			PushNotificationToken: u.PushNotificationToken,
		}
		//infor of shared access users
		accessListInfo = append(accessListInfo, pi)
		accessList = append(accessList, userModels.WalletPermission{
			ID:             permissionID,
			TargetUsername: v.TargetUsername,
			Permission:     v.Permission,
			WalletAddress:  wallet.ID,
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
					ID:             uuid.NewString(),
					TargetUsername: v.TargetUsername,
					Permission:     v.Permission,
					WalletAddress:  linkedWallet.ID,
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

	errDB := dbTX.Omit(clause.Associations).Create(&accessList).Error
	if errDB != nil {
		log.Printf("[CreateSharedWalletAccess] error saving access list:%v\n AccessList:%+v\n", errDB, accessList)
		return returnedWallet, &tErrors.ErrorTemporaryServerError{}
	}

	errDB = dbTX.Omit(clause.Associations).Save(wallet).Error
	if errDB != nil {
		log.Printf("[CreateSharedWalletAccess] error saving shared access status of the wallet:%v\n sharedAccess:%+v\n", errDB, wallet)
		return returnedWallet, &tErrors.ErrorTemporaryServerError{}
	}

	//if has linked wallet, clone the wallet property
	if hasLinkedWallet {
		linkedWallet.SharedAccessEnabled = wallet.SharedAccessEnabled
		linkedWallet.NumberOfApprovalsNeeded = wallet.NumberOfApprovalsNeeded
		errDB := dbTX.Omit(clause.Associations).Create(&linkedWalletAccessList).Error
		if errDB != nil {
			log.Printf("[CreateSharedWalletAccess] error saving linked wallet access list:%v\n Linked wallet AccessList:%+v\n", errDB, linkedWalletAccessList)
			return returnedWallet, &tErrors.ErrorTemporaryServerError{}
		}

		errDB = dbTX.Omit(clause.Associations).Save(linkedWallet).Error
		if errDB != nil {
			log.Printf("[CreateSharedWalletAccess] error saving shared access status of the linked wallet:%v\n linkedsharedAccess:%+v\n", errDB, linkedWallet)
			return returnedWallet, &tErrors.ErrorTemporaryServerError{}
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	accessInfo.NetworkPassPhrase = network.GetBlockchainNetworkPassPhrase()
	accessInfo.SignatureRequired = 1
	target, threshold := sharedAccessTarget(walletOwner, approverUsers, accessInfo.NumberOfApprovalsNeeded)
	var linkedPtr *userModels.UserWallet
	if hasLinkedWallet {
		linkedPtr = &linkedWallet
	}
	calls, err := sharedAccessCalls(ctx, wallet, linkedPtr, target, threshold, gc)
	if err != nil {
		return returnedWallet, err
	}
	switch {
	case len(calls) == 0:
		// only view-only/initiator permissions: nothing changes on-chain, so
		// the owner signs a statement of exactly this change
		statement := sharedAccessAuthorizationPayload(wallet.ID, accessInfo.NumberOfApprovalsNeeded, accessInfo.Permissions)
		accessInfo.Transaction = statement
		if len(accessInfo.TransactionSignature) == 0 {
			return returnedWallet, nil
		}
		if err := verifyStatementSignature(wallet.Signer, statement, accessInfo.TransactionSignature); err != nil {
			return returnedWallet, err
		}
	case len(accessInfo.TransactionSignature) == 0:
		// approvers become co-owners of the wallet's Safe: the owner signs
		// the operation adding them and setting the threshold
		op, err := prepareSharedAccessOperation(ctx, signerUser, walletOwner, wallet, calls, target, threshold, 0, gc)
		if err != nil {
			return returnedWallet, err
		}
		accessInfo.Transaction = op.Transaction
		accessInfo.Messages = append(accessInfo.Messages, fmt.Sprintf("The %v approver(s) become co-signers of this wallet; %v signature(s) will be needed for each transaction.", len(target)-1, threshold))
		if hasLinkedWallet {
			accessInfo.Messages = append(accessInfo.Messages, fmt.Sprintf("The linked wallet %v gets the same co-signers.", linkedWallet.Alias))
		}
		accessInfo.Messages = append(accessInfo.Messages, operationMessages(op.Prepared)...)
		return returnedWallet, nil
	default:
		rec, p, err := loadSharedAccessOperation(accessInfo.Transaction, wallet, target, threshold, gc)
		if err != nil {
			return returnedWallet, err
		}
		hash, err := SignSingleOwnerOperation(ctx, rec, p, walletOwner.PrimarySigner, accessInfo.TransactionSignature, gc)
		if err != nil {
			return returnedWallet, err
		}
		accessInfo.TransactionID = hash
	}

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
	if sErr := signerUser.EnsureNotSuspended(); sErr != nil {
		return nil, nil, nil, nil, nil, nil, sErr
	}
	if sErr := walletOwner.EnsureNotSuspended(); sErr != nil {
		return nil, nil, nil, nil, nil, nil, sErr
	}

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
	if wallet.WalletType == 1 && wallet.LinkedWalletAddress != nil {
		// set the linked wallet if it is a tokenization wallet
		hasLinkedWallet = true

		linkedWallet, err = userModels.UserWalletID(*wallet.LinkedWalletAddress).GetWallet(gc.DB, gc)
		if err != nil {
			err = &tErrors.CustomError{
				Param:      "linkedWalletAddress",
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
	walletID := userModels.UserWalletID(accessInfo.WalletAddress)
	var linkedWalletID userModels.UserWalletID
	if hasLinkedWallet {
		linkedWalletID = userModels.UserWalletID(linkedWallet.ID)
	}

	fw, err := walletID.GetWallet(gc.DB, gc)
	if err != nil {
		log.Printf("[ModifySharedWalletAccess] error could not get wallet object for %v %+v\n", accessInfo.WalletAddress, accessInfo)
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
	ops := make([]basetxn.Operation, 0)
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
		if !AddressHasViewOnlyAccessWACL(accessInfo.WalletAddress, accessInfo.AddedPermissions, gc) || !AddressHasViewOnlyAccessWACL(accessInfo.WalletAddress, accessInfo.ModifiedPermissions, gc) {
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
				ErrMessage: fmt.Sprintf("Access can only be granted/revoked to/from TrovoApp account, not a subwallet [%v]", v.TargetUsername),
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
				ErrMessage: fmt.Sprintf("TrovoApp account [%v] could not be validated at this time.", v.TargetUsername),
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
			WalletAddress:         wallet.ID,
			WalletAlias:           wallet.Alias,
			PushNotificationToken: u.PushNotificationToken,
		})
		revokedList = append(revokedList, userModels.WalletPermission{
			ID:             userPermission.ID,
			TargetUsername: v.TargetUsername,
			Permission:     v.Permission,
			WalletAddress:  wallet.ID,
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
				ID:             linkedUserPermission.ID,
				TargetUsername: v.TargetUsername,
				Permission:     v.Permission,
				WalletAddress:  linkedWallet.ID,
			})
		}

		if v.Permission == "APPROVER" {
			//get ops to add.
			op, ignore, e := generateRemoveSharedAccessOps(wallet, walletOwner, &u, gc)
			if e != nil {
				log.Printf("[ModifySharedWalletAccess] error generating blockchain operation for %+v: error: %v\n", v, e)

				// err = &tErrors.CustomError{
				// 	Param:      "username",
				// 	Err:        "error-trovo-wallet-account-invalid",
				// 	ErrMessage: fmt.Sprintf("TrovoApp account [%v] could not be validated on the blockchain at this time.", v.TargetUsername),
				// 	Code:       http.StatusBadRequest,
				// }

				return

			} else {
				if !ignore {
					ops = append(ops, op)
				}
			}

			if hasLinkedWallet {
				//get ops to add.
				op, ignore, e := generateRemoveSharedAccessOps(&linkedWallet, walletOwner, &u, gc)
				if e != nil {
					log.Printf("[ModifySharedWalletAccess] error generating blockchain operation for %+v ON linked wallet: error: %v\n", v, e)
					return

				} else {
					if !ignore {
						ops = append(ops, op)
					}
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
				ErrMessage: fmt.Sprintf("Access can only be granted to TrovoApp account, not a subwallet [%v]", v.TargetUsername),
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
				ErrMessage: fmt.Sprintf("TrovoApp account [%v] could not be validated at this time.", v.TargetUsername),
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
			WalletAddress:         wallet.ID,
			WalletAlias:           wallet.Alias,
			PushNotificationToken: u.PushNotificationToken,
		})
		modifiedList = append(modifiedList, userModels.WalletPermission{
			CreatedAt:      ePermission.CreatedAt,
			UpdatedAt:      ePermission.UpdatedAt,
			ID:             ePermission.ID,
			WalletAddress:  wallet.ID,
			TargetUsername: ePermission.TargetUsername,
			Permission:     v.Permission, //modify the permission
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
				CreatedAt:      linkedEPermission.CreatedAt,
				UpdatedAt:      linkedEPermission.UpdatedAt,
				ID:             linkedEPermission.ID,
				WalletAddress:  linkedWallet.ID,
				TargetUsername: linkedEPermission.TargetUsername,
				Permission:     v.Permission, //modify the permission
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
					ErrMessage: fmt.Sprintf("TrovoApp account [%v] could not be validated on the blockchain at this time. Unable to add this access for this user.", v.TargetUsername),
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
						ErrMessage: fmt.Sprintf("TrovoApp account [%v] could not be validated on the blockchain at this time. Unable to add this access on linked wallet for this user.", v.TargetUsername),
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
			op, ignore, e := generateRemoveSharedAccessOps(wallet, walletOwner, &u, gc)
			if e != nil {
				log.Printf("[ModifySharedWalletAccess] error generating blockchain operation for %+v: error: %v\n", v, e)

				err = &tErrors.CustomError{
					Param:      "username",
					Err:        "error-trovo-wallet-account-invalid",
					ErrMessage: fmt.Sprintf("TrovoApp account [%v] could not be validated on the blockchain at this time. Unable to remove the approver access from blockchain at this time.", v.TargetUsername),
					Code:       http.StatusBadRequest,
				}

				return

			}
			if !ignore {
				ops = append(ops, op)

			}

			if hasLinkedWallet {
				// attempt to remove the public key as signer
				op, ignore, e := generateRemoveSharedAccessOps(&linkedWallet, walletOwner, &u, gc)
				if e != nil {
					log.Printf("[ModifySharedWalletAccess] error generating blockchain operation on linked wallet for %+v: error: %v\n", v, e)

					err = &tErrors.CustomError{
						Param:      "username",
						Err:        "error-trovo-wallet-account-invalid",
						ErrMessage: fmt.Sprintf("TrovoApp account [%v] could not be validated on the blockchain at this time. Unable to remove the approver access on linked wallet from blockchain at this time.", v.TargetUsername),
						Code:       http.StatusBadRequest,
					}

					return

				}
				if !ignore {
					ops = append(ops, op)
				}
			}

		}
		checkAccess[v.TargetUsername+v.Permission] = userModels.WalletPermissionInfo{
			TargetUsername:        v.TargetUsername,
			Name:                  name,
			Permission:            v.Permission,
			WalletAddress:         wallet.ID,
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
				ErrMessage: fmt.Sprintf("Access can only be granted to TrovoApp account, not a subwallet [%v]", v.TargetUsername),
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
				ErrMessage: fmt.Sprintf("TrovoApp account [%v] could not be validated at this time.", v.TargetUsername),
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
			WalletAddress:         wallet.ID,
			WalletAlias:           wallet.Alias,
			PushNotificationToken: u.PushNotificationToken,
		})

		addedList = append(addedList, userModels.WalletPermission{
			ID:             permissionID,
			TargetUsername: v.TargetUsername,
			Permission:     v.Permission,
			WalletAddress:  wallet.ID,
		})

		if hasLinkedWallet {
			linkedAddedList = append(linkedAddedList, userModels.WalletPermission{
				ID:             uuid.NewString(),
				TargetUsername: v.TargetUsername,
				Permission:     v.Permission,
				WalletAddress:  linkedWallet.ID,
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
					ErrMessage: fmt.Sprintf("TrovoApp account [%v] could not be validated on the blockchain at this time. Unable to add this access at this time.", v.TargetUsername),
					Code:       http.StatusBadRequest,
				}

				return
			}

			ops = append(ops, o...)
			accessInfo.Messages = append(accessInfo.Messages, m...)

			{
				oRecoveredAccount, ignore := generateRemoveRecoveredAccountAccessOps(wallet, u.Username, gc)
				if len(oRecoveredAccount) > 0 && !ignore {
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
						ErrMessage: fmt.Sprintf("TrovoApp account [%v] could not be validated on the blockchain at this time. Unable to add this access on linked wallet at this time.", v.TargetUsername),
						Code:       http.StatusBadRequest,
					}

					return
				}

				ops = append(ops, o...)
				accessInfo.Messages = append(accessInfo.Messages, m...)

				{
					oRecoveredAccount, ignore := generateRemoveRecoveredAccountAccessOps(&linkedWallet, u.Username, gc)
					if len(oRecoveredAccount) > 0 && !ignore {
						ops = append(ops, oRecoveredAccount...)
					}
				}
			}

		}
		checkAccess[v.TargetUsername+v.Permission] = userModels.WalletPermissionInfo{
			TargetUsername:        v.TargetUsername,
			Name:                  name,
			Permission:            v.Permission,
			WalletAddress:         wallet.ID,
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
			if hasLinkedWallet && len(linkedRevokedList) > 0 {
				e := dbTX.Delete(&linkedRevokedList).Error
				if e != nil {
					log.Println("[ModifySharedWalletAccess] error deleting linked revoked list: ", e)
					err = &tErrors.ErrorTemporaryServerError{}
					return
				}
			}
		}
		if len(modifiedList) > 0 {
			e := dbTX.Omit(clause.Associations).Save(&modifiedList).Error
			if e != nil {
				log.Println("[ModifySharedWalletAccess] error saving modified list: ", e)
				err = &tErrors.ErrorTemporaryServerError{}
				return
			}
			if hasLinkedWallet && len(linkedModifiedList) > 0 {
				e := dbTX.Omit(clause.Associations).Save(&linkedModifiedList).Error
				if e != nil {
					log.Println("[ModifySharedWalletAccess] error saving linked modified list: ", e)
					err = &tErrors.ErrorTemporaryServerError{}
					return
				}
			}
		}
		if len(addedList) > 0 {
			e := dbTX.Omit(clause.Associations).Create(&addedList).Error
			if e != nil {
				log.Println("[ModifySharedWalletAccess] error creating added permissions: ", e)
				err = &tErrors.ErrorTemporaryServerError{}
				return
			}

			if hasLinkedWallet && len(linkedAddedList) > 0 {
				e := dbTX.Omit(clause.Associations).Create(&linkedAddedList).Error
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
	accessInfo.NetworkPassPhrase = network.GetBlockchainNetworkPassPhrase()
	accessInfo.RevokedPermissions = revokedListInfo
	accessInfo.ModifiedPermissions = modifiedListInfo
	accessInfo.AddedPermissions = addedListInfo
	if accessInfo.DryRun {
		// an approval completing only needs the permission lists
		return
	}
	_ = ops // the Safe owner changes are planned from the final permissions below

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	var finalApprovers []*userModels.User
	var finalPermissions []userModels.WalletPermissionInfo
	for _, p := range updatedWallet.Permissions {
		finalPermissions = append(finalPermissions, userModels.WalletPermissionInfo{TargetUsername: p.TargetUsername, Permission: p.Permission})
		if p.Permission != "APPROVER" {
			continue
		}
		u, e := userModels.Username(p.TargetUsername).GetSimpleUser(gc.DB, gc)
		if e != nil {
			err = &tErrors.ErrorTemporaryServerError{}
			return
		}
		finalApprovers = append(finalApprovers, &u)
	}
	target, threshold := sharedAccessTarget(walletOwner, finalApprovers, accessInfo.NumberOfApprovalsNeeded)
	var linkedPtr *userModels.UserWallet
	if hasLinkedWallet {
		linkedPtr = &linkedWallet
	}
	calls, e := sharedAccessCalls(ctx, wallet, linkedPtr, target, threshold, gc)
	if e != nil {
		err = e
		return
	}
	statement := sharedAccessAuthorizationPayload(wallet.ID, accessInfo.NumberOfApprovalsNeeded, finalPermissions)
	coSigners := func(msgs []string) []string {
		if len(calls) == 0 {
			return msgs
		}
		msgs = append(msgs, fmt.Sprintf("After this change the wallet has %v co-signer(s); %v signature(s) are needed for each transaction.", len(target)-1, threshold))
		if hasLinkedWallet {
			msgs = append(msgs, fmt.Sprintf("The linked wallet %v gets the same co-signers.", linkedWallet.Alias))
		}
		return msgs
	}

	if oldNumberOfApprovers == 0 {
		// the owner alone controls the wallet today, and signs the change
		accessInfo.SignatureRequired = 1
		switch {
		case len(calls) == 0:
			accessInfo.Transaction = statement
			if len(accessInfo.TransactionSignature) == 0 {
				return
			}
			if e := verifyStatementSignature(wallet.Signer, statement, accessInfo.TransactionSignature); e != nil {
				err = e
				return
			}
		case len(accessInfo.TransactionSignature) == 0:
			op, e := prepareSharedAccessOperation(ctx, signerUser, walletOwner, wallet, calls, target, threshold, 0, gc)
			if e != nil {
				err = e
				return
			}
			accessInfo.Transaction = op.Transaction
			accessInfo.Messages = append(coSigners(accessInfo.Messages), operationMessages(op.Prepared)...)
			return
		default:
			rec, p, e := loadSharedAccessOperation(accessInfo.Transaction, wallet, target, threshold, gc)
			if e != nil {
				err = e
				return
			}
			hash, e := SignSingleOwnerOperation(ctx, rec, p, walletOwner.PrimarySigner, accessInfo.TransactionSignature, gc)
			if e != nil {
				err = e
				return
			}
			accessInfo.TransactionID = hash
		}
		if e := dbTX.Commit().Error; e != nil {
			log.Printf("[ModifySharedWalletAccess] saving %v: %v", wallet.ID, e)
			err = &tErrors.ErrorTemporaryServerError{}
			return
		}
		wallet.InvalidateUserCache(gc)
		return
	}

	// approvers must approve the change: the initiator previews it
	// (commit=0), then submits it for approval (commit=1)
	approvalsNeeded := wallet.NumberOfApprovalsNeeded
	if accessInfo.Commit == 0 {
		if len(calls) == 0 {
			accessInfo.Transaction = statement
			return
		}
		op, e := prepareSharedAccessOperation(ctx, signerUser, walletOwner, wallet, calls, target, threshold, sharedWalletOperationValidity(), gc)
		if e != nil {
			err = e
			return
		}
		accessInfo.Transaction = op.Transaction
		accessInfo.Messages = append(coSigners(accessInfo.Messages), operationMessages(op.Prepared)...)
		return
	}
	if len(calls) == 0 {
		if accessInfo.Transaction != statement {
			err = &tErrors.CustomError{Param: "transaction", Err: "transaction mismatch", ErrMessage: "The transaction does not match this change. Please start again.", Code: http.StatusBadRequest}
			return
		}
	} else {
		_, p, e := loadSharedAccessOperation(accessInfo.Transaction, wallet, target, threshold, gc)
		if e != nil {
			err = e
			return
		}
		// the Safe's current threshold decides how many approvals it takes
		approvalsNeeded = int(p.Threshold)
	}
	{
		description := fmt.Sprintf("Modify shared access on wallet %v.", wallet.Alias)

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
			ID:                     uuid.New().String(),
			Initiator:              signerUser.Username,
			InitiatorSignerAddress: signerUser.PrimarySigner,
			WalletAddress:          wallet.ID,
			TransactionType:        "MODIFY SHARED ACCESS",
			Description:            description,
			ApprovalsNeeded:        approvalsNeeded,
			TransactionXdr:         accessInfo.Transaction,
			TransactionInfoStr:     &transactionStr,
		}
		// the changes apply only once the approvals are complete
		dbTX.Rollback()
		if e := gc.DB.Omit(clause.Associations).Create(&pendingAuth).Error; e != nil {
			log.Printf("[ModifySharedWalletAccess] saving approval request for %v: %v", wallet.ID, e)
			err = &tErrors.ErrorTemporaryServerError{}
			return
		}
		accessInfo.TransactionID = "PENDING_AUTH"
	}
	return
}

func RemoveSharedWalletAccess(signerUser *userModels.User, wallet *userModels.UserWallet, accessInfo *userModels.DisableSharedAccessInfo, gc *sharedconfig.GlobalConfig) (err error) {
	// var managedAccess userModels.UserWalletSharedAccess
	var hasLinkedWallet bool
	var linkedWallet userModels.UserWallet
	if wallet.WalletType == 1 && wallet.LinkedWalletAddress != nil {
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
		linkedWallet, errLinked = userModels.UserWalletID(*wallet.LinkedWalletAddress).GetWallet(gc.DB, gc)
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
				ErrMessage: fmt.Sprintf("TrovoApp account [%v] could not be validated at this time.", v.TargetUsername),
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

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	accessInfo.NetworkPassPhrase = network.GetBlockchainNetworkPassPhrase()
	// back to the owner alone (threshold 1) on the wallet and its linked wallet
	target, threshold := sharedAccessTarget(&walletOwner, nil, 0)
	var linkedPtr *userModels.UserWallet
	if hasLinkedWallet {
		linkedPtr = &linkedWallet
	}
	calls, err := sharedAccessCalls(ctx, wallet, linkedPtr, target, threshold, gc)
	if err != nil {
		return err
	}
	statement := sharedAccessAuthorizationPayload(wallet.ID, 0, nil)
	_ = approverUsers

	if numberOfApprovers > 0 {
		// the approvers approve removing themselves
		accessInfo.MultiParty = 1
		if accessInfo.Commit == 0 {
			if len(calls) == 0 {
				accessInfo.Transaction = statement
				return nil
			}
			op, err := prepareSharedAccessOperation(ctx, signerUser, &walletOwner, wallet, calls, target, threshold, sharedWalletOperationValidity(), gc)
			if err != nil {
				return err
			}
			accessInfo.Transaction = op.Transaction
			accessInfo.Messages = append(accessInfo.Messages, "Once approved, the approvers stop being co-signers of this wallet.")
			accessInfo.Messages = append(accessInfo.Messages, operationMessages(op.Prepared)...)
			return nil
		}
		if len(calls) == 0 {
			if accessInfo.Transaction != statement {
				return &tErrors.CustomError{Param: "transaction", Err: "transaction mismatch", ErrMessage: "The transaction does not match this change. Please start again.", Code: http.StatusBadRequest}
			}
		} else {
			_, p, err := loadSharedAccessOperation(accessInfo.Transaction, wallet, target, threshold, gc)
			if err != nil {
				return err
			}
			// the Safe's current threshold decides how many approvals it takes
			approvalsNeeded = int(p.Threshold)
		}
		accessInfo.TransactionID = "PENDING_AUTH"
		description := fmt.Sprintf("Disabling shared access on wallet %v.\n This will remove the permissions:\n%v", wallet.Alias, userPermissions)
		transactionByte, _ := json.Marshal(*accessInfo)
		transactionStr := string(transactionByte)
		pendingAuth := userModels.PendingAuth{
			ID:                     uuid.New().String(),
			Initiator:              signerUser.Username,
			InitiatorSignerAddress: signerUser.PrimarySigner,
			WalletAddress:          wallet.ID,
			TransactionType:        "DISABLE SHARED ACCESS",
			Description:            description,
			ApprovalsNeeded:        approvalsNeeded,
			TransactionXdr:         accessInfo.Transaction,
			TransactionInfoStr:     &transactionStr,
		}
		if e := gc.DB.Omit(clause.Associations).Create(&pendingAuth).Error; e != nil {
			log.Printf("[RemoveSharedWalletAccess] saving approval request for %v: %v", wallet.ID, e)
			return &tErrors.ErrorTemporaryServerError{}
		}
		return nil
	}

	// no approvers: the owner alone disables it
	accessInfo.SignatureRequired = 1
	switch {
	case len(calls) == 0:
		accessInfo.Transaction = statement
		if len(accessInfo.TransactionSignature) == 0 {
			return nil
		}
		if err := verifyStatementSignature(wallet.Signer, statement, accessInfo.TransactionSignature); err != nil {
			return err
		}
	case len(accessInfo.TransactionSignature) == 0:
		op, err := prepareSharedAccessOperation(ctx, signerUser, &walletOwner, wallet, calls, target, threshold, 0, gc)
		if err != nil {
			return err
		}
		accessInfo.Transaction = op.Transaction
		accessInfo.Messages = append(accessInfo.Messages, operationMessages(op.Prepared)...)
		return nil
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

	if len(calls) > 0 {
		rec, p, err := loadSharedAccessOperation(accessInfo.Transaction, wallet, target, threshold, gc)
		if err != nil {
			return err
		}
		hash, err := SignSingleOwnerOperation(ctx, rec, p, walletOwner.PrimarySigner, accessInfo.TransactionSignature, gc)
		if err != nil {
			return err
		}
		accessInfo.TransactionID = hash
	}

	dbTX.Commit()
	wallet.InvalidateUserCache(gc)
	return nil

}

func generateAddSharedAccessOps(wallet *userModels.UserWallet, walletOwner *userModels.User, approver *userModels.User, gc *sharedconfig.GlobalConfig) (ops []basetxn.Operation, messages []string, err error) {
	// Safe owners are planned from the final permission set instead
	// (sharedAccessCalls); nothing to build per approver.
	return nil, nil, nil
}

func generateRemoveSharedAccessOps(wallet *userModels.UserWallet, walletOwner *userModels.User, approver *userModels.User, gc *sharedconfig.GlobalConfig) (op basetxn.Operation, ignore bool, err error) {
	// see generateAddSharedAccessOps
	return nil, true, nil
}

func generateRemoveRecoveredAccountAccessOps(wallet *userModels.UserWallet, approverUsernameAdded string, gc *sharedconfig.GlobalConfig) (ops []basetxn.Operation, ignore bool) {
	// see generateAddSharedAccessOps; a recovered approver's old key is
	// replaced when the owner set is planned from current signers
	return nil, true
}

func HasAccessToAddress(signerAddress, targetAddress string, gc *sharedconfig.GlobalConfig) (hasAccess bool) {
	signerUser, err := usersDB.GetUserFromPrimarySigner(signerAddress, gc.DB, gc)

	if err != nil {
		return false
	}
	wallet, err := usersDB.GetWallet(targetAddress, gc.DB)
	if err != nil {
		return false
	}
	if wallet.SharedAccessEnabled == 0 {
		return false
	}
	return wallet.SignerHasAccess(&signerUser, gc)

}

func HasInitiatorPermissionToAddress(ownerSignerAddress, targetAddress string, gc *sharedconfig.GlobalConfig) (hasAccess bool) {
	user, err := usersDB.GetUserFromPrimarySigner(ownerSignerAddress, gc.DB, gc)

	if err != nil {
		return false
	}
	walletPermissions := user.WalletsSharedWithUser
	if len(walletPermissions) == 0 {
		return false
	}
	for _, walletAccess := range walletPermissions {
		if walletAccess.WalletAddress == targetAddress && walletAccess.Permission == "INITIATOR" {
			return true
		}
	}

	return false
}

func SignerHasInitiatorPermissionToAddress(signerOwner userModels.User, targetAddress string, gc *sharedconfig.GlobalConfig) (hasAccess bool) {

	walletPermissions := signerOwner.WalletsSharedWithUser
	if len(walletPermissions) == 0 {
		return false
	}
	for _, walletAccess := range walletPermissions {
		if walletAccess.WalletAddress == targetAddress && walletAccess.Permission == "INITIATOR" {
			return true
		}
	}

	return false
}

// sharedAccessAuthorizationPayload is the base64 statement a wallet signer
// signs to authorize an app-layer-only shared access change (the apps'
// signBase64Txn decodes and signs exactly these bytes). It names the
// wallet, the approval threshold and every permission granted, in a
// deterministic order, so a signature authorizes only this change.
func sharedAccessAuthorizationPayload(walletID string, approvalsNeeded int, permissions []userModels.WalletPermissionInfo) string {
	perms := make([]string, 0, len(permissions))
	for _, p := range permissions {
		perms = append(perms, strings.ToLower(p.TargetUsername)+"="+strings.ToUpper(p.Permission))
	}
	sort.Strings(perms)
	statement := fmt.Sprintf("Shared access for wallet %s: %d approval(s) needed; permissions: %s", walletID, approvalsNeeded, strings.Join(perms, ", "))
	return base64.StdEncoding.EncodeToString([]byte(statement))
}
