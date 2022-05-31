package users

import "time"

type User struct {
	CreatedAt             time.Time    `json:"createdAt"`
	UpdatedAt             time.Time    `json:"updatedAt"`
	LastUpdatedMobileOn   time.Time    `json:"lastUpdatedMobileOn"`
	ID                    string       `json:"id"`
	Username              string       `gorm:"size:16; index:idx_user_unique_username, unique" json:"username"`
	Email                 string       `gorm:"size:45; index:idx_user_unique_email, unique" json:"email"`
	ImageThumbnailURL     *string      `json:"imageThumbnailURL"`
	FirstName             string       `gorm:"size:50" json:"firstName"`
	LastName              string       `gorm:"size:50" json:"lastName"`
	Mobile                *string      `gorm:"size:16; index:idx_user_unique_phone, unique" json:"mobile"`
	PublicKey             string       `gorm:"size:56; index:idx_user_unique_public_key, unique" json:"publicKey"`
	Referrer              *string      `gorm:"size:16; index:idx_user_referrer" json:"referrer"`
	ReferralLink          *string      `json:"referralLink"`
	ReferralQrCode        *string      `json:"referralQrCode"`
	PushNotificationToken *string      `json:"pushNotificationToken"`
	Corporate             uint         `gorm:"type:integer;not null; default:0" json:"corporate"`
	MobileVerified        uint         `gorm:"type:integer;not null; default:0" json:"mobileVerified"`
	MembershipType        uint         `gorm:"type:integer;not null; default:0" json:"membershipType"`
	MembershipExpiry      *time.Time   `json:"membershipExpiry"`
	KYCVerified           uint         `gorm:"type:integer;not null; default:0" json:"kycVerified"`
	WalletRecoveryEnabled uint         `gorm:"type:integer;not null; default:0" json:"walletRecoveryEnabled"`
	UserWallets           []UserWallet `json:"userWallets"`
	PublicIP              string       `gorm:"size:45" json:"publicIP"`
	CountryCode           *string      `gorm:"size:2;null"`
	Latitude              *float64     `gorm:"null"`
	Longitude             *float64     `gorm:"null"`
	City                  *string      `gorm:"null;size:100"`
	Region                *string      `gorm:"null;size:100"`
	RegionName            *string      `gorm:"null;size:100"`
	TimeZone              *string      `gorm:"null;size:100"`
	ISP                   *string      `gorm:"null;size:150"`
	Verified              int          `gorm:"type:integer;not null;default:0" json:"verified"`
	Suspended             int          `gorm:"type:integer;not null;default:0" json:"suspended"`
	SuspensionReason      *string      `gorm:"null" json:"suspensionReason"`
}

type UserWallet struct {
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
	ID            string    `gorm:"size:56" json:"publicKey"`
	TempPublicKey *string   `gorm:"size:56;index:idx_user_wallet_temp_key;null"`
	Tag           string    `gorm:"size:10" json:"tag"`
	Description   string    `gorm:"size:100" json:"description"`
	// Link                    *string                 `json:"link"` //rather generated on the fly since it is tired to Asset
	// QrCode                  *string                 `json:"qrCode"` //generated on the fly since it is tied to asset
	Alias                   string                  `gorm:"size:27; index:idx_unique_alias, unique" json:"alias"` //primaryUsername_tag for sub wallets
	Signer                  string                  `gorm:"size:56; index:idx_user_wallet_signer" json:"signer"`  //if ID is same as signer, then it is a primary wallet
	UserID                  string                  `gorm:"type:integer;not null; default:0;index:idx_user_wallets_user_id" json:"userId"`
	ManagedAccessEnabled    uint                    `gorm:"type:integer;not null; default:0" json:"managedAccessEnabled"`
	UserWalletManagedAccess UserWalletManagedAccess `json:"userWalletManagedAccess"`
}

type UserWalletManagedAccess struct {
	CreatedAt           time.Time      `json:"createdAt"`
	UpdatedAt           time.Time      `json:"updatedAt"`
	ID                  string         `gorm:"" json:"accessId"`
	UserWalletID        string         `gorm:"size:56; index:idx_manage_access_user_wallet_id" json:"publicKey"`
	NumberOfAuthorizers uint           `gorm:"type:integer; default:1" json:"numberOfAuthorizers"`
	AccessList          []WalletAccess `json:"accessList"`
}
type WalletAccess struct {
	Username                  string `gorm:"size:16; primaryKey" json:"username"`
	AccessLevel               string `gorm:"size:10" json:"accessLevel"`
	UserWalletManagedAccessID string `gorm:"index:idx_wallet_access_wallet_access_id" json:"userWalletManagedAccessId"`
}

type AccessLevel struct {
	ID          uint64
	AccessLevel string `gorm:"size:text" json:"accessList"`
}

type UserRegistrationInfo struct {
	Username              string `json:"username"`
	Email                 string `json:"email"`
	FirstName             string `json:"firstName"`
	LastName              string `json:"lastName"`
	Mobile                string `json:"mobile"`
	MobileCountryCode     string `json:"mobileCountryCode"`
	PublicKey             string `json:"publicKey"`
	Referrer              string `json:"referrer"`
	PushNotificationToken string `json:"pushNotificationToken"`
	Corporate             uint   `json:"corporate"`
	VerificationCode      string `json:"verificationCode"`
	PublicIP              string `json:"-"`
}
