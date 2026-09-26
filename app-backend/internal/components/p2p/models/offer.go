package p2p

import "time"

// Offer types
const (
	OfferTypeBuy  = "BUY"
	OfferTypeSell = "SELL"
)

// Offer statuses
const (
	OfferStatusDraft          = "DRAFT"
	OfferStatusActive         = "ACTIVE"
	OfferStatusPaused         = "PAUSED"
	OfferStatusOutOfLiquidity = "OUT_OF_LIQUIDITY"
	OfferStatusSuspended      = "SUSPENDED"
	OfferStatusExpired        = "EXPIRED"
	OfferStatusClosed         = "CLOSED"
)

// Offer availability
const (
	OfferAvailabilityOnline  = "ONLINE"
	OfferAvailabilityOffline = "OFFLINE"
)

// PaymentMethod is the merchant's fiat settlement channel for an offer/order.
type PaymentMethod struct {
	PaymentChannel string `json:"paymentChannel" gorm:"size:30;not null;default:''"` // e.g. BANK_TRANSFER, MOBILE_MONEY, E_WALLET
	Provider       string `json:"provider" gorm:"size:60;not null;default:''"`       // e.g. GTBANK, MTN, PAYPAL
	Account        string `json:"account" gorm:"size:100;not null;default:''"`
}

// Offer is a merchant's standing buy/sell listing on the P2P marketplace.
type Offer struct {
	ID                  string        `json:"id" gorm:"primaryKey;size:36"`
	MerchantUsername    string        `json:"merchantUsername" gorm:"size:70;not null;index:idx_p2p_offer_merchant_username"`
	MerchantUserID      string        `json:"merchantUserId" gorm:"size:100;not null;index:idx_p2p_offer_merchant_user_id"`
	OfferType           string        `json:"offerType" gorm:"size:10;not null;index:idx_p2p_offer_type"`
	Asset               string        `json:"asset" gorm:"size:12;not null;index:idx_p2p_offer_asset"`
	ContractAddress     string        `json:"contractAddress" gorm:"size:56;not null;default:''"`
	PaymentMethod       PaymentMethod `json:"paymentMethod" gorm:"embedded;embeddedPrefix:payment_method_"`
	Country             string        `json:"country" gorm:"size:60;not null;default:''"`
	CountryCode         string        `json:"countryCode" gorm:"size:3;not null;index:idx_p2p_offer_country_code"`
	Currency            string        `json:"currency" gorm:"size:10;not null;index:idx_p2p_offer_currency"`
	PriceType           string        `json:"priceType" gorm:"size:20;not null;default:'FIXED'"` // FIXED or FLOATING
	Price               string        `json:"price" gorm:"size:60;not null"`
	PriceMargin         string        `json:"priceMargin" gorm:"size:60;not null;default:'0'"`
	MinOrderAmount      string        `json:"minOrderAmount" gorm:"size:60;not null"`
	MaxOrderAmount      string        `json:"maxOrderAmount" gorm:"size:60;not null"`
	AvailableLiquidity  string        `json:"availableLiquidity" gorm:"size:60;not null;default:'0'"`
	ReservedLiquidity   string        `json:"reservedLiquidity" gorm:"size:60;not null;default:'0'"`
	Remark              string        `json:"remark" gorm:"size:500;default:''"`
	AvailabilityStatus  string        `json:"availabilityStatus" gorm:"size:10;not null;default:'OFFLINE'"`
	Status              string        `json:"status" gorm:"size:20;not null;default:'DRAFT';index:idx_p2p_offer_status"`
	Version             int           `json:"version" gorm:"not null;default:1"`
	CreatedAt           time.Time     `json:"createdAt"`
	UpdatedAt           time.Time     `json:"updatedAt"`
}
