package users

import "strings"

func (u *User) ToJSON() (jsonObj UserJSON) {
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

		for _, uw := range u.UserWallets {
			uwJson := uw.ToJSON()
			// log.Printf("[UserToJSON] added user wallet [%+v]\n", uwJson)
			jsonObj.UserWallets = append(jsonObj.UserWallets, uwJson)

		}

	}
	// log.Println("[UserToJSON] ended user wallets json and returning data")
	return jsonObj
}

func (uw *UserWallet) ToJSON() (jsonObj UserWalletJSON) {
	jsonObj.CreatedAt = uw.CreatedAt
	jsonObj.ID = uw.ID
	jsonObj.Alias = uw.Alias
	jsonObj.Signer = uw.Signer
	jsonObj.UserID = uw.UserID
	jsonObj.ManagedAccessEnabled = uw.ManagedAccessEnabled
	if uw.ManagedAccessEnabled == 1 {
		majson := uw.UserWalletManagedAccess.ToJSON()
		jsonObj.UserWalletManagedAccess = &majson
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

func (uwma *UserWalletManagedAccess) ToJSON() (jsonObj UserWalletManagedAccessJSON) {
	jsonObj.CreatedAt = uwma.CreatedAt
	jsonObj.UpdatedAt = uwma.UpdatedAt
	jsonObj.ID = uwma.ID
	jsonObj.UserWalletID = uwma.UserWalletID
	for _, uwmaAL := range uwma.AccessList {
		jsonObj.AccessList = append(jsonObj.AccessList, uwmaAL.ToJSON())
	}
	return
}

func (wa *WalletAccess) ToJSON() (jsonObj WalletAccessJSON) {
	jsonObj.CreatedAt = wa.CreatedAt
	jsonObj.UpdatedAt = wa.UpdatedAt
	jsonObj.Username = wa.Username
	jsonObj.AccessLevel = wa.AccessLevel
	jsonObj.UserWalletManagedAccessID = wa.UserWalletManagedAccessID
	return
}
