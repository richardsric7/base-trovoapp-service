package users

import "time"

type FaceMatchVerificationCred struct {
	Message string `json:"message"`
	Data    struct {
		SelfieVideo UploadData `json:"selfie-video"`
		Passport    UploadData `json:"passport"`
	} `json:"data"`
}

type AWSFields struct {
	Bucket            string `json:"bucket"`
	Policy            string `json:"Policy"`
	XAmzDate          string `json:"X-Amz-Date"`
	XAmzAlgorithm     string `json:"X-Amz-Algorithm"`
	XAmzSignature     string `json:"X-Amz-Signature"`
	XAmzSecurityToken string `json:"X-Amz-Security-Token"`
	XAmzCredential    string `json:"X-Amz-Credential"`
	Key               string `json:"key"`
}

type DriverLicenseVerificationCred struct {
	Message string `json:"message"`
	Data    struct {
		DrivingLicenseFront UploadData `json:"driving_license-front"`
		DrivingLicenseBack  UploadData `json:"driving_license-back"`
	} `json:"data"`
}

type VerificationResponse struct {
	Message string `json:"message"`
	Data    struct {
		VerificationID string `json:"verificationId"`
	} `json:"data"`
}

type OneLiquidityErrorResponse struct {
	Message string `json:"message"`
}

type VerificationRequest struct {
	AccountID string `json:"accountId"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type PassportVerificationCred struct {
	Message string `json:"message"`
	Data    struct {
		PassportUploadData UploadData `json:"passport"`
	} `json:"data"`
}

type UploadData struct {
	Fields AWSFields `json:"fields"`
	URL    string    `json:"url"`
}

type NationalIDVerificationCred struct {
	Message string `json:"message"`
	Data    struct {
		NationalIDBack  UploadData `json:"national_id-back"`
		NationalIDFront UploadData `json:"national_id-front"`
	} `json:"data"`
}

type ProofOfResidenceCred struct {
	Message string `json:"message"`
	Data    struct {
		VerificationID string `json:"verificationId"`
		DocType        string `json:"docType"`
		PresignedURL   struct {
			ProofOfResidency UploadData `json:"proof_of_residency"`
		} `json:"presignedUrl"`
	} `json:"data"`
}

type FacematchPassport struct {
	CreatedAt         time.Time `json:"-"`
	UpdatedAt         time.Time `json:"-"`
	ID                string    `gorm:"size:100;not null"`
	UserID            string    `gorm:"size:100;not null"`
	VerificationID    string    `gorm:"size:100;not null"`
	SelfieVideoFormat string    `gorm:"size:100;not null"`
	PictureFormat     string    `gorm:"size:100;not null"`
	SelfieKey         string    `gorm:"size:100;not null"`
	DocumentKey       string    `gorm:"size:100;not null"`
	KycStatus         int       `gorm:"type:integer;default:0;not null"`
	KycData           *string   `gorm:"null"`
}

type FacematchNationalID struct {
	CreatedAt         time.Time `json:"-"`
	UpdatedAt         time.Time `json:"-"`
	ID                string    `gorm:"size:100;not null"`
	UserID            string    `gorm:"size:100;not null"`
	VerificationID    string    `gorm:"size:100;not null"`
	SelfieVideoFormat string    `gorm:"size:100;not null"`
	PictureFormat     string    `gorm:"size:100;not null"`
	SelfieKey         string    `gorm:"size:100;not null"`
	DocumentKey       string    `gorm:"size:100;not null"`
	KycStatus         int       `gorm:"type:integer;default:0;not null"`
	KycData           *string   `gorm:"null"`
}

type FacematchDrivingLicense struct {
	CreatedAt         time.Time `json:"-"`
	UpdatedAt         time.Time `json:"-"`
	ID                string    `gorm:"size:100;not null"`
	UserID            string    `gorm:"size:100;not null"`
	VerificationID    string    `gorm:"size:100;not null"`
	SelfieVideoFormat string    `gorm:"size:100;not null"`
	PictureFormat     string    `gorm:"size:100;not null"`
	SelfieKey         string    `gorm:"size:100;not null"`
	DocumentKey       string    `gorm:"size:100;not null"`
	KycStatus         int       `gorm:"type:integer;default:0;not null"`
	KycData           *string   `gorm:"null"`
}

type GovermentIDProofOfResidency struct {
	CreatedAt      time.Time `json:"-"`
	UpdatedAt      time.Time `json:"-"`
	ID             string    `gorm:"size:100;not null"`
	UserID         string    `gorm:"size:100;not null"`
	VerificationID string    `gorm:"size:100;not null"`
	PictureFormat  string    `gorm:"size:100;not null"`
	DocumentKey    string    `gorm:"size:100;not null"`
	KycStatus      int       `gorm:"type:integer;default:0;not null"`
	KycData        *string   `gorm:"null"`
}

type StartDocumentRequest struct {
	Inputs         []DocumentInputs `json:"inputs"`
	VerificationID string           `json:"verificationId"`
}

type DocumentInputs struct {
	Country  string `json:"country"`
	DocType  string `json:"docType"`
	MimeType string `json:"mimeType"`
	Key      string `json:"key"`
	Side     string `json:"side,omitempty"`
}

type OKResponse struct {
	Message string `json:"message"`
}

type CryptoWalletDepositAddress struct {
	ID                   string    `json:"id"`
	CreatedAt            time.Time `gorm:"default:now()" json:"createdAt"`
	UserID               string    `gorm:"not null;" json:"-"`
	TrovoWalletPublicKey string    `json:"TrovoWalletPublicKey"`
	Currency             string    `gorm:"not null;size:12" json:"currency"`
	DepositAddress       string    `gorm:"not null;size:100" json:"depositAddress"`
	Network              string    `gorm:"not null;size:100" json:"network"`
}

type CryptoSubwalletResponse struct {
	Message string          `json:"message"`
	Data    CryptoSubWallet `json:"data"`
}
type CryptoSubwalletsResponse struct {
	Message string            `json:"message"`
	Data    []CryptoSubWallet `json:"data"`
}
type CryptoSubWallet struct {
	WalletID     string          `json:"walletId"`
	UID          string          `json:"uid"`
	Addresses    []CryptoAddress `json:"addresses"`
	Currency     string          `json:"currency"`
	IntegratorPk string          `json:"integratorPk,omitempty"`
	CreatedAt    string          `json:"createdAt,omitempty"`
	UpdatedAt    string          `json:"updatedAt,omitempty"`
}
type CryptoAddress struct {
	Address string `json:"address"`
	Network string `json:"network"`
}

type CryptoDepositResponse struct {
	Message string              `json:"message"`
	Data    DepositResponseItem `json:"data"`
}
type DepositResponseItem struct {
	DepositID   string `json:"depositId"`
	TxID        string `json:"txId"`
	Amount      string `json:"amount"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
	Currency    string `json:"currency"`
	Decimal     int    `json:"decimal"`
	Fees        string `json:"fees"`
	FromAddress string `json:"fromAddress"`
	IsCompleted bool   `json:"isCompleted"`
	IsValid     bool   `json:"isValid"`
	IsVerified  bool   `json:"isVerified"`
	ToAddress   string `json:"toAddress"`
}

type CryptoWithdrawalNetworksResponse struct {
	Message string              `json:"message"`
	Data    []WithdrawalNetwork `json:"data"`
}

type CryptoWithdrawal struct {
	TrovoWalletPublicKey string `gorm:"size:100" json:"TrovoWalletPublicKey"`
	WithdrawalID         string `gorm:"primaryKey" json:"withdrawalId"`
	Amount               int    `json:"amount"`
	CreatedAt            string `json:"createdAt"`
	UpdatedAt            string `json:"updatedAt"`
	Currency             string `gorm:"size:100" json:"currency"`
	Network              string `gorm:"size:100" json:"network"`
	ToAddress            string `gorm:"size:100" json:"toAddress"`
	Status               string `json:"status"`
}

type SubWalletInput struct {
	Currency string `json:"currency"`
	UID      string `json:"uid"`
}
