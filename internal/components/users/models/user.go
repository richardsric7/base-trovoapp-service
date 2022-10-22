package users

import (
	"time"

	"github.com/gofrs/uuid"
)

type User struct {
	CreatedAt                time.Time          `json:"createdAt"`
	UpdatedAt                time.Time          `json:"updatedAt"`
	LastUpdatedMobileOn      time.Time          `json:"lastUpdatedMobileOn"`
	LastRecoveredAccountOn   time.Time          `json:"lastRecoveredAccountOn"`
	ID                       string             `json:"id"`
	Username                 string             `gorm:"size:16; index:idx_user_unique_username, unique" json:"username"`
	Email                    string             `gorm:"size:45; index:idx_user_unique_email, unique" json:"email"`
	ImageThumbnailURL        *string            `json:"imageThumbnailURL"`
	FirstName                string             `gorm:"size:50" json:"firstName"`
	LastName                 *string            `gorm:"size:50" json:"lastName"`
	Mobile                   *string            `gorm:"size:16; index:idx_user_unique_phone, unique" json:"mobile"`
	PublicKey                string             `gorm:"size:56; index:idx_user_unique_public_key, unique" json:"publicKey"`
	PrimarySigner            string             `gorm:"size:56; index:idx_user_unique_primary_signer, unique" json:"primarySigner"`
	Referrer                 *string            `gorm:"size:16; index:idx_user_referrer" json:"referrer"`
	ReferralLink             *string            `json:"referralLink"`
	ReferralQrCode           *string            `json:"referralQrCode"`
	PushNotificationToken    *string            `json:"pushNotificationToken"`
	Corporate                int                `gorm:"type:integer;not null; default:0" json:"corporate"`
	MobileVerified           int                `gorm:"type:integer;not null; default:0" json:"mobileVerified"`
	MembershipType           int                `gorm:"type:integer;not null; default:0" json:"membershipType"`
	MembershipExpiry         *time.Time         `json:"membershipExpiry"`
	KYCVerified              int                `gorm:"type:integer;not null; default:0" json:"kycVerified"`
	AccountRecoveryEnabled   int                `gorm:"type:integer;not null; default:0" json:"accountRecoveryEnabled"`
	AccountRecoveryExpiresOn *time.Time         `gorm:"null" json:"accountRecoveryExpiresOn"`
	UserWallets              []UserWallet       `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"userWallets"`
	PublicIP                 string             `gorm:"size:45" json:"publicIP"`
	CountryCode              *string            `gorm:"size:2;null"`
	Latitude                 *float64           `gorm:"null"`
	Longitude                *float64           `gorm:"null"`
	City                     *string            `gorm:"null;size:100"`
	Region                   *string            `gorm:"null;size:100"`
	RegionName               *string            `gorm:"null;size:100"`
	TimeZone                 *string            `gorm:"null;size:100"`
	ISP                      *string            `gorm:"null;size:150"`
	HasSecurityQuestions     int                `gorm:"type:integer;not null;default:0" json:"hasSecurityQuestions"`
	Verified                 int                `gorm:"type:integer;not null;default:0" json:"verified"`
	Suspended                int                `gorm:"type:integer;not null;default:0" json:"suspended"`
	SuspensionReason         *string            `gorm:"null" json:"suspensionReason"`
	WalletsSharedWithUser    []WalletPermission `gorm:"foreignKey:TargetUsername;references:Username;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type UserWallet struct {
	CreatedAt               time.Time          `json:"createdAt"`
	UpdatedAt               time.Time          `json:"updatedAt"`
	ID                      string             `gorm:"size:56" json:"publicKey"`
	TempPublicKey           *string            `gorm:"size:56;index:idx_user_wallet_temp_key;null"`
	Tag                     *string            `gorm:"null;size:16" json:"tag"`
	Description             *string            `gorm:"null;size:100" json:"description"`
	Alias                   string             `gorm:"size:27; index:idx_unique_alias, unique" json:"alias"` //primaryUsername_tag for sub wallets
	Signer                  string             `gorm:"size:56; index:idx_user_wallet_signer" json:"signer"`  //if ID is same as signer, then it is a primary wallet
	UserID                  string             `gorm:"type:integer;not null; default:0;index:idx_user_wallets_user_id" json:"userId"`
	SharedAccessEnabled     int                `gorm:"type:integer;not null; default:0" json:"sharedAccessEnabled"`
	Tracked                 int                `gorm:"type:integer;not null;default:0" json:"-"`
	PrimaryWallet           int                `gorm:"type:integer;not null;default:0" json:"primaryWallet"`
	NumberOfApprovalsNeeded int                `gorm:"type:integer; default:0" json:"numberOfApprovalsNeeded"`
	AssetIssuerWallet       int                `gorm:"type:integer; default:0" json:"assetIssuerWallet"`
	Permissions             []WalletPermission `gorm:"foreignKey:WalletPublicKey;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"permissions"`
	SharedAccessCreatedAt   time.Time          `json:"sharedAccessCreatedAt"`
	SharedAccessUpdatedAt   time.Time          `json:"sharedAccessUpdatedAt"`
}

type WalletPermission struct {
	CreatedAt       time.Time `json:"-"`
	UpdatedAt       time.Time `json:"-"`
	ID              string
	WalletPublicKey string `gorm:"size:90;not null; index:access_level_permission,unique;index:idx_public_key_shared" json:"walletPublicKey"`
	TargetUsername  string `gorm:"size:16;not null; index:access_level_permission,unique;" json:"targetUsername"`
	Permission      string `gorm:"size:10;not null; index:access_level_permission,unique" json:"permission"`
}
type UserWalletSharedAccessInfo struct {
	WalletPublicKey         string                 `json:"walletPublicKey"`
	NumberOfApprovalsNeeded int                    `json:"numberOfApprovalsNeeded"`
	Permissions             []WalletPermissionInfo `json:"permissions"`
	Transaction             string                 `json:"transaction"`
	TransactionSignature    string                 `json:"transactionSignature"`
	TransactionID           string                 `json:"transactionId"`
	NetworkPassPhrase       string                 `json:"networkPassPhrase"`
	Messages                []string               `json:"messages"`
	SignatureRequired       int                    `json:"signatureRequired"`
	Approvers               []User                 `json:"-"`
	Initiators              []User                 `json:"-"`
	Viewers                 []User                 `json:"-"`
}
type DisableSharedAccessInfo struct {
	WalletPublicKey      string                 `json:"walletPublicKey"`
	Transaction          string                 `json:"transaction"`
	TransactionSignature string                 `json:"transactionSignature"`
	TransactionID        string                 `json:"transactionId"`
	NetworkPassPhrase    string                 `json:"networkPassPhrase"`
	Messages             []string               `json:"messages"`
	SignatureRequired    int                    `json:"signatureRequired"`
	MultiParty           int                    `json:"multiParty"`
	Permissions          []WalletPermissionInfo `json:"-"`
	Commit               int                    `json:"commit"`
}
type WalletPermissionInfo struct {
	ID                    string  `json:"Id"`
	WalletPublicKey       string  `json:"-"`
	WalletAlias           string  `json:"-"`
	TargetUsername        string  `json:"targetUsername"`
	Name                  string  `json:"name"`
	Permission            string  `json:"permission"`
	PushNotificationToken *string `json:"-"`
}

type Permissions struct {
	ID         uint64
	Permission string `gorm:"size:text" json:"permission"`
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

// UserWalletsharedAccessID is type for wallet access id
// type UserWalletSharedAccessID string

// UserWalletID is type for wallet/sub-wallet Public Key
type UserWalletID string

// Username is a type for username of a user
type Username string

type WalletAlias string

// UserSigner is type for signer Public Key
type UserSigner string

type ApprovalID string

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
	AssetIssuerWallet       int      `json:"assetIssuerWallet"`
	Messages                []string `json:"messages"`
}

type SecurityQuestion struct {
	ID       uint64
	Question string
}

type UserSecurityAnswer struct {
	ID        uint64    `json:"id"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	Username  string    `gorm:"size:20;not null; index:unique_user_security_answer,unique" json:"-"`
	Q1        uint64    `gorm:"not null;" json:"q1"`
	A1        string    `gorm:"size:150;not null;" json:"a1"`
	Q2        uint64    `gorm:"not null;" json:"q2"`
	A2        string    `gorm:"size:150;not null;" json:"a2"`
	Q3        uint64    `gorm:"not null;" json:"q3"`
	A3        string    `gorm:"size:150;not null;" json:"a3"`
}

type UserAccountRecoveryPayload struct {
	Transaction          string             `json:"transaction"`
	TransactionSignature string             `json:"transactionSignature"`
	TransactionID        string             `json:"transactionId"`
	NetworkPassPhrase    string             `json:"networkPassPhrase"`
	Messages             []string           `json:"messages"`
	SecurityAnswers      UserSecurityAnswer `json:"securityAnswers"`
}

type AccountRecoveryRequest struct {
	NewSignerPublicKey                string             `json:"newSignerPublicKey"`
	DisableOldSignerFromPrimaryWallet uint64             `json:"disableOldSignerFromPrimaryWallet"`
	Commit                            uint64             `json:"commit"`
	Messages                          []string           `json:"messages"`
	SecurityAnswers                   UserSecurityAnswer `json:"securityAnswers"`
	EmailOTP                          string             `json:"emailOtp"`
	Username                          string             `json:"username"`
	TransactionID                     string             `json:"transactionId"`
}
type InactiveAccountRecoveryRequest struct {
	NewSignerPublicKey string             `json:"newSignerPublicKey"`
	SecurityAnswers    UserSecurityAnswer `json:"securityAnswers"`
	EmailOTP           string             `json:"emailOtp"`
	Username           string             `json:"username"`
}

type PendingAuth struct {
	CreatedAt                    time.Time                     `json:"createdAt"`
	UpdatedAt                    time.Time                     `gorm:"default:now()" json:"updatedAt"`
	ID                           string                        `gorm:"size:56" json:"id"`
	Initiator                    string                        `gorm:"size:20;not null;index:idx_pending_auth_initiator" json:"initiator"`
	InitiatorSignerPublicKey     string                        `gorm:"size:56;not null;index:idx_pending_auth_signer_public_key" json:"initiatorSignerPublicKey"`
	WalletPublicKey              string                        `gorm:"size:56;not null;index:idx_pending_auth_wallet_public_key" json:"walletPublicKey"`
	TransactionType              string                        `gorm:"size:28;not null;index:idx_pending_auth_transaction_type" json:"transactionType"`
	Description                  string                        `gorm:"not null;" json:"description"`
	ApprovalsNeeded              int                           `gorm:"not null;" json:"approvalsNeeded"`
	ApprovalsGotten              int                           `gorm:"not null;default:0" json:"approvalsGotten"`
	TransactionStatus            string                        `gorm:"size:20;not null;default:'PENDING';index:idx_pending_auth_transaction_status" json:"transactionStatus"`
	RejectedBy                   *string                       `gorm:"size:20;null;index:idx_pending_auth_rejected_by" json:"rejectedBy"`
	ReasonForRejection           *string                       `gorm:"size:200;null;" json:"reasonForRejection"`
	TransactionXdr               string                        `gorm:"not null;" json:"transactionXdr"`
	TransactionInfoStr           *string                       `gorm:"null;" json:"transactionInfoStr"`
	TransactionID                *string                       `gorm:"size:70;null;index:idx_pending_auth_transaction_id" json:"transactionID"`
	PendingTransactionSignatures []PendingTransactionSignature `json:"pendingTransactionSignatures"`
}
type PendingTransactionSignature struct {
	CreatedAt                time.Time `json:"createdAt"`
	ID                       string    `gorm:"size:56"`
	PendingAuthID            string    `gorm:"size:56;not null;index:idx_pending_trxsig_pending_auth,unique"`
	Approver                 string    `gorm:"size:20;not null;index:idx_pending_trxsig_pending_auth,unique;index:idx_pending_trxsig_approver" json:"approver"`
	ApproverSignerPublicKey  string    `gorm:"size:56;not null;index:idx_pending_trxsig_approver_signer" json:"approverSignerPublicKey"`
	TransactionWithSignature string    `gorm:"not null" json:"transactionWithSignature"`
}

type ApprovalPayload struct {
	Transaction          string `json:"transaction"`
	TransactionSignature string `json:"transactionSignature"`
	RemainingApprovals   int    `json:"remainingApproval"`
	NetworkPassPhrase    string `json:"networkPassPhrase"`
}
