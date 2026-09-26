package p2p

import "time"

// Order lifecycle statuses (Plan Section 24)
const (
	OrderStatusAwaitingApproval           = "AWAITING_APPROVAL"
	OrderStatusAwaitingEscrowDeposit      = "AWAITING_ESCROW_DEPOSIT"
	OrderStatusAwaitingPayment            = "AWAITING_PAYMENT"
	OrderStatusAwaitingPaymentConfirmation = "AWAITING_PAYMENT_CONFIRMATION"
	OrderStatusCompleted                  = "COMPLETED"
	OrderStatusRejected                   = "REJECTED"
	OrderStatusCancelled                  = "CANCELLED"
	OrderStatusExpired                    = "EXPIRED"
)

// Escrow deposit states (Plan Section 43)
const (
	EscrowDepositStatusNotDeposited = "NOT_DEPOSITED"
	EscrowDepositStatusPartial      = "PARTIAL"
	EscrowDepositStatusFull         = "FULL"
	EscrowDepositStatusOverpaid     = "OVERPAID"
	EscrowDepositStatusInvalid      = "INVALID"
	EscrowDepositStatusRefundable   = "REFUNDABLE"
	EscrowDepositStatusConfirmed    = "CONFIRMED"
)

// Order is a single P2P trade instance created against an Offer.
type Order struct {
	ID                 string `json:"id" gorm:"primaryKey;size:36"`
	OfferID            string `json:"offerId" gorm:"size:36;not null;index:idx_p2p_order_offer_id"`
	CustomerUserID     string `json:"customerUserId" gorm:"size:100;not null;index:idx_p2p_order_customer_user_id"`
	CustomerUsername   string `json:"customerUsername" gorm:"size:70;not null"`
	MerchantUserID     string `json:"merchantUserId" gorm:"size:100;not null;index:idx_p2p_order_merchant_user_id"`
	MerchantUsername   string `json:"merchantUsername" gorm:"size:70;not null"`
	OfferType          string `json:"offerType" gorm:"size:10;not null"`
	Asset              string `json:"asset" gorm:"size:12;not null"`
	AssetContractAddress string `json:"assetContractAddress" gorm:"size:56;not null;default:''"`

	PaymentMethodSnapshot PaymentMethod `json:"paymentMethodSnapshot" gorm:"embedded;embeddedPrefix:payment_method_snapshot_"`
	Country               string        `json:"country" gorm:"size:60;not null;default:''"`
	CountryCode            string       `json:"countryCode" gorm:"size:3;not null;default:''"`
	Currency               string       `json:"currency" gorm:"size:10;not null"`
	Price                  string       `json:"price" gorm:"size:60;not null"`
	SpecifiedAssetAmount   string       `json:"specifiedAssetAmount" gorm:"size:60;not null"`
	PaymentAmount          string       `json:"paymentAmount" gorm:"size:60;not null"`

	BuyerPlatformFee     string `json:"buyerPlatformFee" gorm:"size:60;not null;default:'0'"`
	BuyerRegulatoryFee   string `json:"buyerRegulatoryFee" gorm:"size:60;not null;default:'0'"`
	BuyerPlatformFeeVat  string `json:"buyerPlatformFeeVat" gorm:"size:60;not null;default:'0'"`
	BuyerRegulatoryFeeVat string `json:"buyerRegulatoryFeeVat" gorm:"size:60;not null;default:'0'"`
	BuyerTotalFees       string `json:"buyerTotalFees" gorm:"size:60;not null;default:'0'"`
	BuyerTotalVat        string `json:"buyerTotalVat" gorm:"size:60;not null;default:'0'"`
	BuyerTotalCharges    string `json:"buyerTotalCharges" gorm:"size:60;not null;default:'0'"`
	BuyerNetAssetAmount  string `json:"buyerNetAssetAmount" gorm:"size:60;not null;default:'0'"`

	SellerPlatformFee       string `json:"sellerPlatformFee" gorm:"size:60;not null;default:'0'"`
	SellerRegulatoryFee     string `json:"sellerRegulatoryFee" gorm:"size:60;not null;default:'0'"`
	SellerPlatformFeeVat    string `json:"sellerPlatformFeeVat" gorm:"size:60;not null;default:'0'"`
	SellerRegulatoryFeeVat  string `json:"sellerRegulatoryFeeVat" gorm:"size:60;not null;default:'0'"`
	SellerTotalFees         string `json:"sellerTotalFees" gorm:"size:60;not null;default:'0'"`
	SellerTotalVat          string `json:"sellerTotalVat" gorm:"size:60;not null;default:'0'"`
	SellerTotalCharges      string `json:"sellerTotalCharges" gorm:"size:60;not null;default:'0'"`
	SellerEscrowAssetAmount string `json:"sellerEscrowAssetAmount" gorm:"size:60;not null;default:'0'"`

	CombinedPlatformFee   string `json:"combinedPlatformFee" gorm:"size:60;not null;default:'0'"`
	CombinedRegulatoryFee string `json:"combinedRegulatoryFee" gorm:"size:60;not null;default:'0'"`
	CombinedVat           string `json:"combinedVat" gorm:"size:60;not null;default:'0'"`

	PlatformFeeWalletAddress   string `json:"platformFeeWalletAddress" gorm:"size:56;not null;default:''"`
	RegulatoryFeeWalletAddress string `json:"regulatoryFeeWalletAddress" gorm:"size:56;not null;default:''"`
	VatWalletAddress           string `json:"vatWalletAddress" gorm:"size:56;not null;default:''"`
	FeeConfigurationVersion    int    `json:"feeConfigurationVersion" gorm:"not null;default:0"`
	FeeConfigurationSnapshot   string `json:"feeConfigurationSnapshot" gorm:"type:text"` // JSON snapshot of TradeFeeConfiguration at order-creation time

	EscrowDepositShortlink        string `json:"escrowDepositShortlink" gorm:"size:255;default:''"`
	EscrowDepositQRCode           string `json:"escrowDepositQrCode" gorm:"type:text"`
	EscrowDepositTransactionHash  string `json:"escrowDepositTransactionHash" gorm:"size:100;default:''"`
	AssetReleaseTransactionHash   string `json:"assetReleaseTransactionHash" gorm:"size:100;default:''"`

	EscrowDepositStatus    string `json:"escrowDepositStatus" gorm:"size:20;not null;default:'NOT_DEPOSITED'"`
	ExpectedEscrowAmount   string `json:"expectedEscrowAmount" gorm:"size:60;not null;default:'0'"`
	DepositedEscrowAmount  string `json:"depositedEscrowAmount" gorm:"size:60;not null;default:'0'"`
	RefundableAmount       string `json:"refundableAmount" gorm:"size:60;not null;default:'0'"`

	// AssetDepositor/AssetRecipient/FiatPayer/FiatRecipient are resolved from
	// OfferType per Plan Section 11 and stored explicitly so downstream code
	// (escrow deposit, settlement, refunds) never re-derives the role mapping.
	AssetDepositor string `json:"assetDepositor" gorm:"size:100;not null"`
	AssetRecipient string `json:"assetRecipient" gorm:"size:100;not null"`
	FiatPayer      string `json:"fiatPayer" gorm:"size:100;not null"`
	FiatRecipient  string `json:"fiatRecipient" gorm:"size:100;not null"`

	OrderStatus                     string     `json:"orderStatus" gorm:"size:32;not null;default:'AWAITING_APPROVAL';index:idx_p2p_order_status"`
	IsDisputed                      bool       `json:"isDisputed" gorm:"not null;default:false"`
	BuyerPayoutAddress               string     `json:"buyerPayoutAddress" gorm:"size:56;default:''"`
	MerchantPaymentDetailsSnapshot   string     `json:"merchantPaymentDetailsSnapshot" gorm:"type:text"`
	ExpiresAt                        *time.Time `json:"expiresAt"`
	CompletedAt                      *time.Time `json:"completedAt"`
	PaymentSentAt                    *time.Time `json:"paymentSentAt"`
	PaymentConfirmedAt               *time.Time `json:"paymentConfirmedAt"`
	AssetReleasedAt                  *time.Time `json:"assetReleasedAt"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
