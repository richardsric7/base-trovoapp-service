package users

import (
	"time"
	assets "trovo-wallet-api/internal/components/assets/models"
)

type UserJSON struct {
	ID                     string                    `json:"-"`
	Username               string                    `json:"username"`
	Email                  string                    `json:"email"`
	ImageThumbnailURL      string                    `json:"imageThumbnailURL"`
	FirstName              string                    `json:"firstName"`
	LastName               string                    `json:"lastName"`
	Mobile                 string                    `json:"mobile"`
	PublicKey              string                    `json:"publicKey"`
	PrimarySigner          string                    `json:"primarySigner"`
	Referrer               string                    `json:"referrer"`
	ReferralLink           string                    `json:"referralLink"`
	ReferralQrCode         string                    `json:"referralQrCode"`
	PushNotificationToken  string                    `json:"pushNotificationToken"`
	Corporate              int                       `json:"corporate"`
	MobileVerified         int                       `json:"mobileVerified"`
	MembershipType         int                       `json:"membershipType"`
	MembershipExpiry       time.Time                 `json:"membershipExpiry"`
	KYCVerified            int                       `json:"kycVerified"`
	AccountRecoveryEnabled int                       `json:"accountRecoveryEnabled"`
	UserWallets            []UserWalletJSON          `json:"userWallets"`
	Verified               int                       `json:"verified"`
	Suspended              int                       `json:"suspended"`
	HasSecurityQuestions   int                       `json:"hasSecurityQuestions"`
	CuratedSwapList        []assets.CuratedSwapAsset `json:"curatedSwapList"`
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
	WalletType              int                    `json:"walletType"`      //0=normal, 1= assetIssuing, 2= MarketMaking, 3 = bulkPayment
	WalletThreshold         int                    `json:"walletThreshold"` //0=no shared access, 1 = view-Only shared access, 2 = approver is present
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
	FullName        string    `json:"fullName"`
	Permission      string    `json:"permission"`
}

type AuthJSON struct {
	CreatedAt           time.Time `json:"createdAt"`
	UpdatedAt           time.Time `json:"updatedAt"`
	ID                  string    `json:"id"`
	WalletOwnerUsername string    `json:"walletOwnerUsername"`
	WalletPublicKey     string    `json:"walletPublicKey"`
	Alias               string    `json:"alias"`
	Initiator           string    `json:"initiator"`
	TransactionType     string    `json:"transactionType"`
	Description         string    `json:"description"`
	ApprovalsNeeded     int       `json:"approvalsNeeded"`
	ApprovalsGotten     int       `gorm:"not null;default:0" json:"approvalsGotten"`
	TransactionStatus   string    `json:"transactionStatus"`
	RejectedBy          string    `json:"rejectedBy"`
	ReasonForRejection  string    `json:"reasonForRejection"`
	ApprovedBy          string    `json:"approvedBy"`
	Transaction         string    `json:"transaction,omitempty"`
}

type PaginatedAuths struct {
	Pages        int        `json:"pages"`
	CurrentPage  int        `json:"currentPage"`
	TotalRecords int        `json:"totalRecords"`
	Limit        int        `json:"limit"`
	Records      []AuthJSON `json:"records"`
}
