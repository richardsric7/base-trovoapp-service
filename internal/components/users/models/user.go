package users

import "time"

type User struct {
	CreatedAt             time.Time    `json:"createdAt"`
	UpdatedAt             time.Time    `json:"updatedAt"`
	ID                    string       `json:"id"`
	Username              string       `gorm:"size:16; index:idx_unique_username, unique" json:"username"`
	Email                 string       `gorm:"size:45; index:idx_unique_email, unique" json:"email"`
	ImageThumbnail        string       `json:"imageThumbnail"`
	FirstName             string       `gorm:"size:50" json:"firstName"`
	LastName              string       `gorm:"size:50" json:"lastName"`
	Mobile                *string      `gorm:"size:16" json:"mobile"`
	PublicKey             string       `gorm:"size:56" json:"publicKey"`
	Referrer              *string      `gorm:"size:16" json:"referrer"`
	ReferralLink          *string      `json:"referralLink"`
	ReferralQrCode        *string      `json:"referralQrCode"`
	PushNotificationToken *string      `json:"pushNotificationToken"`
	Corporate             uint         `gorm:"type:tinyint; default:0" json:"corporate"`
	MobileVerified        uint         `gorm:"type:tinyint; default:0" json:"mobileVerified"`
	MemershipType         uint         `gorm:"type:tinyint; default:0" json:"membershipType"`
	MembershipExpiry      *time.Time   `json:"membershipExpiry"`
	KYCVerified           uint         `gorm:"type:tinyint; default:0" json:"kycVerified"`
	WalletRecoveryEnabled uint         `gorm:"type:tinyint; default:0" json:"walletRecoveryEnabled"`
	UserWallets           []UserWallet `json:"userWallets"`
}

type UserWallet struct {
	CreatedAt               time.Time               `json:"createdAt"`
	UpdatedAt               time.Time               `json:"updatedAt"`
	ID                      string                  `gorm:"size:56" json:"publicKey"`
	Tag                     string                  `gorm:"size:10" json:"tag"`
	Description             string                  `gorm:"size:100" json:"description"`
	Link                    *string                 `json:"link"`
	QrCode                  *string                 `json:"qrCode"`
	Alias                   string                  `gorm:"size:27; index:idx_unique_alias, unique" json:"alias"` //primaryUsername_tag for sub wallets
	Signer                  string                  `gorm:"size:56; index:idx_sub_wallet_signer" json:"signer"`   //if ID is same as signer, then it is a primary wallet
	UserID                  string                  `json:"userId"`
	ManagedAccessEnabled    uint                    `gorm:"type:tinyint; default:0" json:"managedAccessEnabled"`
	UserWalletManagedAccess UserWalletManagedAccess `json:"userWalletManagedAccess"`
}

type UserWalletManagedAccess struct {
	CreatedAt           time.Time      `json:"createdAt"`
	UpdatedAt           time.Time      `json:"updatedAt"`
	ID                  string         `gorm:"" json:"accessId"`
	UserWalletID        string         `gorm:"size:56" json:"publicKey"`
	NumberOfAuthorizers uint           `gorm:"type:tinyint; default:1" json:"numberOfAuthorizers"`
	AccessList          []WalletAccess `json:"accessList"`
}
type WalletAccess struct {
	Username                  string `gorm:"size:16; primaryKey" json:"username"`
	AccessLevel               string `gorm:"size:10" json:"accessLevel"`
	UserWalletManagedAccessID string `gorm:"size:56" json:"userWalletManagedAccessId"`
}

type AccessLevel struct {
	ID          uint64
	AccessLevel string `gorm:"size:text" json:"accessList"`
}
