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
	Corporate              int              `json:"corporate"`
	MobileVerified         int              `json:"mobileVerified"`
	MembershipType         int              `json:"membershipType"`
	MembershipExpiry       time.Time        `json:"membershipExpiry"`
	KYCVerified            int              `json:"kycVerified"`
	AccountRecoveryEnabled int              `json:"accountRecoveryEnabled"`
	UserWallets            []UserWalletJSON `json:"userWallets"`
	Verified               int              `json:"verified"`
	Suspended              int              `json:"suspended"`
	HasSecurityQuestions   int              `json:"hasSecurityQuestions"`
}

type UserWalletJSON struct {
	CreatedAt               time.Time              `json:"createdAt"`
	ID                      string                 `json:"publicKey"`
	TempPublicKey           string                 `json:"-"`
	Tag                     string                 `json:"tag"`
	Description             string                 `json:"description"`
	Alias                   string                 `json:"alias"`  //primaryUsername_tag for sub wallets
	Signer                  string                 `json:"signer"` //if ID is same as signer, then it is a primary wallet
	UserID                  string                 `json:"userId"`
	SharedAccessEnabled     int                    `json:"sharedAccessEnabled"`
	PrimaryWallet           int                    `json:"primaryWallet"`
	NumberOfApprovalsNeeded int                    `json:"numberOfApprovalsNeeded"`
	Permissions             []WalletPermissionJSON `json:"permissions"`
	SharedAccessCreatedAt   time.Time              `json:"sharedAccessCreatedAt"`
	SharedAccessUpdatedAt   time.Time              `json:"sharedAccessUpdatedAt"`
}

type WalletPermissionJSON struct {
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
	WalletPublicKey string    `json:"walletPublicKey"`
	TargetUsername  string    `json:"targetUsername"`
	Permission      string    `json:"permission"`
	// UserWalletSharedAccessID string    `json:"userWalletSharedAccessId"`
}
