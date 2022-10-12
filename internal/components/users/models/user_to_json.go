package users

import (
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
	if uw.SharedAccessEnabled == 1 {
		//get shared access
		for _, permision := range uw.Permissions {
			jsonObj.Permissions = append(jsonObj.Permissions, permision.ToJSON())
		}
		jsonObj.NumberOfApprovalsNeeded = uw.NumberOfApprovalsNeeded
		jsonObj.SharedAccessCreatedAt = uw.SharedAccessCreatedAt
		jsonObj.SharedAccessUpdatedAt = uw.SharedAccessUpdatedAt

	} else {
		jsonObj.Permissions = make([]WalletPermissionJSON, 0)
	}

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
