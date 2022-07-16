package users

import (
	"time"

	"github.com/gofrs/uuid"
)

type User struct {
	CreatedAt             time.Time    `json:"createdAt"`
	UpdatedAt             time.Time    `json:"updatedAt"`
	LastUpdatedMobileOn   time.Time    `json:"lastUpdatedMobileOn"`
	ID                    string       `json:"id"`
	Username              string       `gorm:"size:16; index:idx_user_unique_username, unique" json:"username"`
	Email                 string       `gorm:"size:45; index:idx_user_unique_email, unique" json:"email"`
	ImageThumbnailURL     *string      `json:"imageThumbnailURL"`
	FirstName             string       `gorm:"size:50" json:"firstName"`
	LastName              *string      `gorm:"size:50" json:"lastName"`
	Mobile                *string      `gorm:"size:16; index:idx_user_unique_phone, unique" json:"mobile"`
	PublicKey             string       `gorm:"size:56; index:idx_user_unique_public_key, unique" json:"publicKey"`
	PrimarySigner         string       `gorm:"size:56; index:idx_user_unique_primary_signer, unique" json:"primarySigner"`
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
	CreatedAt               time.Time               `json:"createdAt"`
	UpdatedAt               time.Time               `json:"updatedAt"`
	ID                      string                  `gorm:"size:56" json:"publicKey"`
	TempPublicKey           *string                 `gorm:"size:56;index:idx_user_wallet_temp_key;null"`
	Tag                     *string                 `gorm:"null;size:16" json:"tag"`
	Description             *string                 `gorm:"null;size:100" json:"description"`
	Alias                   string                  `gorm:"size:27; index:idx_unique_alias, unique" json:"alias"` //primaryUsername_tag for sub wallets
	Signer                  string                  `gorm:"size:56; index:idx_user_wallet_signer" json:"signer"`  //if ID is same as signer, then it is a primary wallet
	UserID                  string                  `gorm:"type:integer;not null; default:0;index:idx_user_wallets_user_id" json:"userId"`
	ManagedAccessEnabled    uint                    `gorm:"type:integer;not null; default:0" json:"managedAccessEnabled"`
	UserWalletManagedAccess UserWalletManagedAccess `json:"userWalletManagedAccess"`
	Tracked                 int                     `gorm:"type:integer;not null;default:0" json:"-"`
}

type UserWalletManagedAccess struct {
	CreatedAt           time.Time      `json:"createdAt"`
	UpdatedAt           time.Time      `json:"updatedAt"`
	ID                  string         `gorm:"" json:"accessId"`
	UserWalletID        string         `gorm:"size:56; index:idx_manage_access_user_wallet_id" json:"publicKey"`
	NumberOfAuthorizers uint           `gorm:"type:integer; default:1" json:"numberOfAuthorizers"`
	AccessList          []WalletAccess `json:"accessList"`
}
type UserWalletManagedAccessInfo struct {
	UserWalletManagedAccessID string             `json:"userWalletManagedAccessId"`
	NumberOfAuthorizers       uint               `json:"numberOfAuthorizers"`
	PublicKey                 string             `json:"publicKey"`
	AccessList                []WalletAccessInfo `json:"accessList"`
	Transaction               string             `json:"transaction"`
	TransactionSignature      string             `json:"transactionSignature"`
	TransactionID             string             `json:"transactionId"`
	NetworkPassPhrase         string             `json:"networkPassPhrase"`
	Messages                  []string           `json:"messages"`
	SignatureRequired         uint               `json:"signatureRequired"`
}
type WalletAccess struct {
	CreatedAt                 time.Time `json:"createdAt"`
	UpdatedAt                 time.Time `json:"updatedAt"`
	ID                        string
	Username                  string `gorm:"size:16;not null; index:access_level_permission,unique" json:"username"`
	AccessLevel               string `gorm:"size:10;not null; index:access_level_permission,unique" json:"accessLevel"`
	UserWalletManagedAccessID string `gorm:"not null;index:idx_wallet_access_wallet_access_id" json:"userWalletManagedAccessId"`
}
type WalletAccessInfo struct {
	ID                        string
	Username                  string `json:"username"`
	AccessLevel               string `json:"accessLevel"`
	UserWalletManagedAccessID string `json:"userWalletManagedAccessId"`
	Name                      string `json:"name"`
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

//UserWalletManagedAccessID is type for wallet access id
type UserWalletManagedAccessID string

//UserWalletID is type for wallet/sub-wallet Public Key
type UserWalletID string

type TrackedWallet struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	PublicKey string    `gorm:"index:idx_tracked_wallet_public_key,unique"`
	Alias     string    `gorm:"index:idx_tracked_wallet_alias"`
	Name      string    `gorm:"index:idx_tracked_wallet_name"`
}

type TrackedPublicKey struct {
	PublicKey string `gorm:"primaryKey"`
}

type SubWalletInfo struct {
	PublicKey               string   `json:"publicKey"`
	WalletTag               string   `json:"walletTag"`
	WalletDescription       string   `json:"walletDescription"`
	Alias                   string   `json:"alias"`
	Transaction             string   `json:"transaction"`
	PrimarySignature        string   `json:"primarySignature"`
	SubWalletSignature      string   `json:"subWalletSignature"`
	TransactionID           string   `json:"transactionId"`
	NetworkPassPhrase       string   `json:"networkPassPhrase"`
	ChannelAccount          string   `json:"channelAccount"`
	ChannelAccountSignature string   `json:"channelAccountSignature"`
	SubWalletMustSign       int      `json:"subWalletMustSign"`
	Messages                []string `json:"messages"`
}
