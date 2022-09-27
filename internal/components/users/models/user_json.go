package users

import "time"

type UserJSON struct {
	ID                     string           `json:"-"`
	Username               string           `json:"username"`
	Email                  string           `json:"email"`
	ImageThumbnailURL      string           `json:"imageThumbnailURL"`
	FirstName              string           `json:"firstName"`
	LastName               string           `json:"lastName"`
	Mobile                 string           `json:"mobile"`
	PublicKey              string           `json:"publicKey"`
	PrimarySigner          string           `json:"primarySigner"`
	Referrer               string           `json:"referrer"`
	ReferralLink           string           `json:"referralLink"`
	ReferralQrCode         string           `json:"referralQrCode"`
	PushNotificationToken  string           `json:"pushNotificationToken"`
	Corporate              uint             `json:"corporate"`
	MobileVerified         uint             `json:"mobileVerified"`
	MembershipType         uint             `json:"membershipType"`
	MembershipExpiry       time.Time        `json:"membershipExpiry"`
	KYCVerified            uint             `json:"kycVerified"`
	AccountRecoveryEnabled uint             `json:"accountRecoveryEnabled"`
	UserWallets            []UserWalletJSON `json:"userWallets"`
	Verified               uint             `json:"verified"`
	Suspended              uint             `json:"suspended"`
	HasSecurityQuestions   uint             `json:"hasSecurityQuestions"`
}

type UserWalletJSON struct {
	CreatedAt               time.Time                    `json:"createdAt"`
	ID                      string                       `json:"publicKey"`
	TempPublicKey           string                       `json:"-"`
	Tag                     string                       `json:"tag"`
	Description             string                       `json:"description"`
	Alias                   string                       `json:"alias"`  //primaryUsername_tag for sub wallets
	Signer                  string                       `json:"signer"` //if ID is same as signer, then it is a primary wallet
	UserID                  string                       `json:"userId"`
	ManagedAccessEnabled    uint                         `json:"managedAccessEnabled"`
	PrimaryWallet           uint                         `json:"primaryWallet"`
	UserWalletManagedAccess *UserWalletManagedAccessJSON `json:"userWalletManagedAccess,omitempty"`
}

type UserWalletManagedAccessJSON struct {
	CreatedAt           time.Time          `json:"createdAt"`
	UpdatedAt           time.Time          `json:"updatedAt"`
	ID                  string             `json:"accessId"`
	UserWalletID        string             `json:"publicKey"`
	NumberOfAuthorizers uint               `json:"numberOfAuthorizers"`
	AccessList          []WalletAccessJSON `json:"accessList"`
}

type WalletAccessJSON struct {
	CreatedAt                 time.Time `json:"createdAt"`
	UpdatedAt                 time.Time `json:"updatedAt"`
	Username                  string    `json:"username"`
	AccessLevel               string    `json:"accessLevel"`
	UserWalletManagedAccessID string    `json:"userWalletManagedAccessId"`
}
