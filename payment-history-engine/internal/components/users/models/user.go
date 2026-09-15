package users

import "time"

type User struct {
	CreatedAt              time.Time    `json:"createdAt"`
	UpdatedAt              time.Time    `json:"updatedAt"`
	LastUpdatedMobileOn    time.Time    `json:"lastUpdatedMobileOn"`
	LastRecoveredAccountOn time.Time    `json:"lastRecoveredAccountOn"`
	ID                     string       `json:"id"`
	Username               string       `gorm:"size:16; index:idx_user_unique_username, unique" json:"username"`
	Email                  string       `gorm:"size:45; index:idx_user_unique_email, unique" json:"email"`
	FirstName              string       `gorm:"size:50" json:"firstName"`
	LastName               *string      `gorm:"size:50" json:"lastName"`
	Mobile                 *string      `gorm:"size:16; index:idx_user_unique_phone, unique" json:"mobile"`
	PublicKey              string       `gorm:"size:56; index:idx_user_unique_public_key, unique" json:"publicKey"`
	PrimarySigner          string       `gorm:"size:56; index:idx_user_unique_primary_signer, unique" json:"primarySigner"`
	PushNotificationToken  *string      `json:"pushNotificationToken"`
	Corporate              int          `gorm:"type:integer;not null; default:0" json:"corporate"`
	UserWallets            []UserWallet `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"userWallets"`
}
type UserWallet struct {
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
	ID            string    `gorm:"size:56" json:"publicKey"`
	TempPublicKey *string   `gorm:"size:56;index:idx_user_wallet_temp_key;null"`
	Tag           *string   `gorm:"null;size:16" json:"tag"`
	Description   *string   `gorm:"null;size:100" json:"description"`
	Alias         string    `gorm:"size:27; index:idx_unique_alias, unique" json:"alias"` //primaryUsername_tag for sub wallets
	Signer        string    `gorm:"size:56; index:idx_user_wallet_signer" json:"signer"`  //if ID is same as signer, then it is a primary wallet
	UserID        string    `gorm:"type:integer;not null; default:0;index:idx_user_wallets_user_id" json:"userId"`
	Tracked       int       `gorm:"type:integer;not null;default:0" json:"-"`
	PrimaryWallet int       `gorm:"type:integer;not null;default:0" json:"primaryWallet"`
}

type WalletPermission struct {
	CreatedAt       time.Time `json:"-"`
	UpdatedAt       time.Time `json:"-"`
	ID              string
	WalletPublicKey string `gorm:"size:90;not null; index:access_level_permission,unique;index:idx_public_key_shared" json:"walletPublicKey"`
	TargetUsername  string `gorm:"size:16;not null; index:access_level_permission,unique;" json:"targetUsername"`
	Permission      string `gorm:"size:10;not null; index:access_level_permission,unique" json:"permission"`
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

// UserWalletManagedAccessID is type for wallet access id
type UserWalletManagedAccessID string

// UserWalletID is type for wallet/sub-wallet Public Key
type UserWalletID string
