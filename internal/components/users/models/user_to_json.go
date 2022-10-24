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
	jsonObj.PublicKey = u.PublicKey
	jsonObj.PrimarySigner = u.PrimarySigner
	jsonObj.Corporate = u.Corporate
	jsonObj.MobileVerified = u.MobileVerified
	jsonObj.MembershipType = u.MembershipType
	jsonObj.KYCVerified = u.KYCVerified
	jsonObj.AccountRecoveryEnabled = u.AccountRecoveryEnabled
	jsonObj.Verified = u.Verified
	jsonObj.Suspended = u.Suspended
	jsonObj.HasSecurityQuestions = u.HasSecurityQuestions
	// log.Println("[UserToJSON] set basic params")
	//nullable
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
	jsonObj.AssetIssuerWallet = uw.AssetIssuerWallet
	viewOnlyAccess := 0

	if uw.SharedAccessEnabled == 1 {
		//get shared access
		for _, permission := range uw.Permissions {
			if permission.Permission != "VIEW-ONLY" {
				viewOnlyAccess = 1
			}
			jsonObj.Permissions = append(jsonObj.Permissions, permission.ToJSON())
		}
		jsonObj.NumberOfApprovalsNeeded = uw.NumberOfApprovalsNeeded
		jsonObj.SharedAccessCreatedAt = uw.SharedAccessCreatedAt
		jsonObj.SharedAccessUpdatedAt = uw.SharedAccessUpdatedAt

	} else {
		jsonObj.Permissions = make([]WalletPermissionJSON, 0)
	}
	jsonObj.HasViewOnlyAccess = viewOnlyAccess

	//nullable
	{
		if uw.TempPublicKey != nil {
			jsonObj.TempPublicKey = *uw.TempPublicKey
		}
		if uw.Tag != nil {
			jsonObj.Tag = *uw.Tag
		}
		if uw.Description != nil {
			jsonObj.Description = *uw.Description
		}

	}
	jsonObj.PrimaryWallet = uw.PrimaryWallet
	if uw.Tag == nil && !strings.Contains(uw.Alias, "_") {
		jsonObj.PrimaryWallet = 1
	}

	return
}

func (wa *WalletPermission) ToJSON() (jsonObj WalletPermissionJSON) {
	jsonObj.CreatedAt = wa.CreatedAt
	jsonObj.UpdatedAt = wa.UpdatedAt
	jsonObj.WalletPublicKey = wa.WalletPublicKey
	jsonObj.TargetUsername = wa.TargetUsername
	jsonObj.Permission = wa.Permission
	return
}

func (a *PendingAuth) ToJSON(gc *sharedconfig.GlobalConfig) (jsonObj AuthJSON) {

	if a == nil {
		return
	}
	var walletOwnerUsername, walletAlias string
	// walletOwner,e:=UserWalletID(a.WalletPublicKey).GetWalletOwner(gc.DB)
	wallet, e := UserWalletID(a.WalletPublicKey).GetWallet(gc.DB)
	if e == nil {
		if wallet.Tag != nil {
			walletOwnerUsername = *wallet.Tag
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
		Alias:               walletAlias,
		Initiator:           a.Initiator,
		TransactionType:     a.TransactionType,
		Description:         a.Description,
		ApprovalsNeeded:     a.ApprovalsNeeded,
		ApprovalsGotten:     a.ApprovalsGotten,
		TransactionStatus:   a.TransactionStatus,
		Transaction:         a.TransactionXdr,
	}
	if a.RejectedBy != nil {
		jsonObj.RejectedBy = *a.RejectedBy
	}
	if a.ReasonForRejection != nil {
		jsonObj.ReasonForRejection = *a.ReasonForRejection
	}

	if a.PendingTransactionSignatures != nil {
		if len(a.PendingTransactionSignatures) > 0 {
			for i, sig := range a.PendingTransactionSignatures {
				jsonObj.ApprovedBy = fmt.Sprintf("%s|%s", sig.Approver, sig.CreatedAt.Format("2006-01-02"))
				if i < len(a.PendingTransactionSignatures)+1 {
					jsonObj.ApprovedBy = fmt.Sprintf("%s,\n", jsonObj.ApprovedBy)
				}
			}
		}
	}

	return jsonObj

}
