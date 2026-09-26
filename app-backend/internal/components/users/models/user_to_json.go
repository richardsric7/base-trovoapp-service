package users

import (
	"fmt"
	"strings"
	"trovo-wallet-api/internal/sharedconfig"
)

func (u *User) ToJSON(gc *sharedconfig.GlobalConfig) (jsonObj UserJSON) {
	jsonObj.ID = u.ID
	jsonObj.Username = u.Username
	jsonObj.Email = u.Email
	jsonObj.FirstName = u.FirstName
	jsonObj.Address = u.Address
	jsonObj.PrimarySigner = u.PrimarySigner
	jsonObj.Corporate = u.Corporate
	jsonObj.MobileVerified = u.MobileVerified
	jsonObj.MembershipType = u.MembershipType
	jsonObj.KYCVerified = u.KYCVerified
	jsonObj.AccountRecoveryEnabled = u.AccountRecoveryEnabled
	jsonObj.Verified = u.Verified
	jsonObj.Suspended = u.Suspended
	jsonObj.IsMerchant = u.IsMerchant
	jsonObj.MerchantOnline = u.MerchantOnline
	jsonObj.HasSecurityQuestions = u.HasSecurityQuestions
	jsonObj.CuratedSwapList = u.GetCuratedSwapList(gc)
	jsonObj.PatronMembership = u.PatronMembership

	{
		if u.LastName != nil {
			jsonObj.LastName = *u.LastName
		}
		if u.ImageThumbnailURL != nil {
			jsonObj.ImageThumbnailURL = *u.ImageThumbnailURL
		}
		if u.Mobile != nil {
			jsonObj.Mobile = *u.Mobile
		}
		if u.MembershipExpiry != nil {
			jsonObj.MembershipExpiry = *u.MembershipExpiry
		}
		if u.Referrer != nil {
			jsonObj.Referrer = *u.Referrer
		}
		if u.ReferralLink != nil {
			jsonObj.ReferralLink = *u.ReferralLink
		}
		if u.ReferralQrCode != nil {
			jsonObj.ReferralQrCode = *u.ReferralQrCode
		}
		if u.PushNotificationToken != nil {
			jsonObj.PushNotificationToken = *u.PushNotificationToken
		}
		if u.CountryCode != nil {
			jsonObj.CountryCode = *u.CountryCode
		}

	}
	{
		//referral statistics
		l1, l2, l3 := u.GetDownlines(gc)

		jsonObj.DownlineStats.Level1 = uint64(len(l1))
		jsonObj.DownlineStats.Level2 = uint64(len(l2))
		jsonObj.DownlineStats.Level3 = uint64(len(l3))
	}
	{
		//referral statistics
		l1, l2, l3 := u.GetUplines(gc)

		jsonObj.Uplines.Level1 = l1
		jsonObj.Uplines.Level2 = l2
		jsonObj.Uplines.Level3 = l3
	}
	if u.UserWallets != nil {
		// log.Println("[UserToJSON] started user wallets json")
		if len(u.UserWallets) > 0 {
			for _, uw := range u.UserWallets {
				//get user permissions

				uwJson := uw.ToJSON(gc)
				// log.Printf("[UserToJSON] added user wallet [%+v]\n", uwJson)
				jsonObj.UserWallets = append(jsonObj.UserWallets, uwJson)
			}
		}

	} else {
		jsonObj.UserWallets = make([]UserWalletJSON, 0)
	}
	//curated swap list

	//user fiat payment methods
	if u.UserFiatPaymentMethods != nil {
		jsonObj.UserFiatPaymentMethods = u.UserFiatPaymentMethods

	} else {
		jsonObj.UserFiatPaymentMethods = make([]UserFiatPaymentMethod, 0)
	}
	//user fiat payment methods

	// log.Println("[UserToJSON] ended user wallets json and returning data")
	return jsonObj
}

func (uw *UserWallet) ToJSON(gc *sharedconfig.GlobalConfig) (jsonObj UserWalletJSON) {
	jsonObj.CreatedAt = uw.CreatedAt
	jsonObj.ID = uw.ID
	jsonObj.Alias = uw.Alias
	jsonObj.Signer = uw.Signer
	jsonObj.UserID = uw.UserID
	jsonObj.SharedAccessEnabled = uw.SharedAccessEnabled
	jsonObj.WalletType = uw.WalletType
	viewOnlyAccess := true
	hasApprover := false

	if uw.SharedAccessEnabled == 1 {
		//get shared access
		for _, permission := range uw.Permissions {
			if permission.Permission != "VIEW-ONLY" {
				viewOnlyAccess = false
			}
			if permission.Permission == "APPROVER" {
				hasApprover = true
			}
			jsonObj.Permissions = append(jsonObj.Permissions, permission.ToJSON(gc))
		}
		jsonObj.NumberOfApprovalsNeeded = uw.NumberOfApprovalsNeeded
		jsonObj.SharedAccessCreatedAt = uw.SharedAccessCreatedAt
		jsonObj.SharedAccessUpdatedAt = uw.SharedAccessUpdatedAt

	} else {
		jsonObj.Permissions = make([]WalletPermissionJSON, 0)
		viewOnlyAccess = false
		hasApprover = false
	}
	if viewOnlyAccess {
		jsonObj.WalletThreshold = 1
	} else if hasApprover {
		jsonObj.WalletThreshold = 2
	} else {
		jsonObj.WalletThreshold = 0
	}

	//nullable
	{
		if uw.TempAddress != nil {
			jsonObj.TempAddress = *uw.TempAddress
		}
		if uw.Tag != nil {
			jsonObj.Tag = *uw.Tag
		}
		if uw.Description != nil {
			jsonObj.Description = *uw.Description
		}
		if uw.LinkedWalletAddress != nil {
			jsonObj.LinkedWalletAddress = *uw.LinkedWalletAddress
		}

	}
	jsonObj.PrimaryWallet = uw.PrimaryWallet
	if uw.Tag == nil && !strings.Contains(uw.Alias, "_") {
		jsonObj.PrimaryWallet = 1
	}

	return
}

func (wa *WalletPermission) ToJSON(gc *sharedconfig.GlobalConfig) (jsonObj WalletPermissionJSON) {
	jsonObj.CreatedAt = wa.CreatedAt
	jsonObj.UpdatedAt = wa.UpdatedAt
	jsonObj.WalletAddress = wa.WalletAddress
	jsonObj.TargetUsername = wa.TargetUsername
	jsonObj.Permission = wa.Permission
	var name string
	u, e := Username(wa.TargetUsername).GetSimpleUser(gc.DB, gc)
	if e == nil {
		name = u.FirstName
		if u.LastName != nil {
			name = fmt.Sprintf("%s %s", name, *u.LastName)
		}
	}
	jsonObj.FullName = name
	return
}

func (a *PendingAuth) ToJSON(gc *sharedconfig.GlobalConfig) (jsonObj AuthJSON) {

	if a == nil {
		return
	}
	var walletOwnerUsername, walletAlias string
	// walletOwner,e:=UserWalletID(a.WalletAddress).GetWalletOwner(gc.DB)
	wallet, e := UserWalletID(a.WalletAddress).GetWallet(gc.DB, gc)
	if e == nil {
		if wallet.Tag != nil {
			//subwallet
			walletOwnerUsername = strings.Split(wallet.Alias, "_")[0]
		} else {
			//primary wallet
			walletOwnerUsername = wallet.Alias
		}
		walletAlias = wallet.Alias

	}
	jsonObj = AuthJSON{
		CreatedAt:           a.CreatedAt,
		UpdatedAt:           a.UpdatedAt,
		ID:                  a.ID,
		WalletOwnerUsername: walletOwnerUsername,
		WalletAddress:       a.WalletAddress,
		Alias:               walletAlias,
		Initiator:           a.Initiator,
		TransactionType:     a.TransactionType,
		Description:         a.Description,
		ApprovalsNeeded:     a.ApprovalsNeeded,
		ApprovalsGotten:     a.ApprovalsGotten,

		TransactionStatus: a.TransactionStatus,
		Transaction:       a.TransactionXdr,
	}
	if a.RejectedBy != nil {
		jsonObj.RejectedBy = *a.RejectedBy
	}
	if a.ApprovedBy != nil {
		jsonObj.ApprovedBy = *a.ApprovedBy
	}
	if a.ReasonForRejection != nil {
		jsonObj.ReasonForRejection = *a.ReasonForRejection
	}
	if a.TransactionID != nil {
		jsonObj.TransactionID = *a.TransactionID
	}

	// if a.PendingTransactionSignatures != nil {
	// 	if len(a.PendingTransactionSignatures) > 0 {
	// 		for i, sig := range a.PendingTransactionSignatures {
	// 			jsonObj.ApprovedBy = fmt.Sprintf("%s on %s", sig.Approver, carbon.Time2Carbon(sig.CreatedAt).ToDateString())
	// 			if i+1 < len(a.PendingTransactionSignatures) {
	// 				jsonObj.ApprovedBy = fmt.Sprintf("%s, ", jsonObj.ApprovedBy)
	// 			}
	// 		}
	// 	}
	// }

	return jsonObj

}
