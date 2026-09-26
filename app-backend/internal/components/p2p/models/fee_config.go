package p2p

import "time"

// Fee configuration statuses
const (
	FeeConfigStatusActive   = "ACTIVE"
	FeeConfigStatusInactive = "INACTIVE"
)

// TradeFeeConfiguration is a versioned, per-country P2P fee schedule
// (Plan Section 52). VAT is applied to the fee amount, not the principal.
type TradeFeeConfiguration struct {
	ID                          string     `json:"id" gorm:"primaryKey;size:36"`
	CountryCode                 string     `json:"countryCode" gorm:"size:3;not null;index:idx_p2p_fee_config_country_code"`
	BuyerPlatformFeePercent     string     `json:"buyerPlatformFeePercent" gorm:"size:20;not null;default:'0'"`
	BuyerRegulatoryFeePercent   string     `json:"buyerRegulatoryFeePercent" gorm:"size:20;not null;default:'0'"`
	SellerPlatformFeePercent    string     `json:"sellerPlatformFeePercent" gorm:"size:20;not null;default:'0'"`
	SellerRegulatoryFeePercent  string     `json:"sellerRegulatoryFeePercent" gorm:"size:20;not null;default:'0'"`
	VatPercent                  string     `json:"vatPercent" gorm:"size:20;not null;default:'0'"`
	EffectiveFrom               time.Time  `json:"effectiveFrom"`
	EffectiveTo                 *time.Time `json:"effectiveTo"`
	Version                     int        `json:"version" gorm:"not null;default:1"`
	Status                      string     `json:"status" gorm:"size:10;not null;default:'ACTIVE';index:idx_p2p_fee_config_status"`
	CreatedAt                   time.Time  `json:"createdAt"`
	UpdatedAt                   time.Time  `json:"updatedAt"`
}

// FeeCollectionWalletConfiguration is a versioned, per-country fee wallet
// mapping (Plan Section 53 - kept as a DB-backed table per your decision,
// a deliberate upgrade over the rest of app-backend's single-global-env-var-
// per-product pattern).
type FeeCollectionWalletConfiguration struct {
	ID                          string     `json:"id" gorm:"primaryKey;size:36"`
	CountryCode                 string     `json:"countryCode" gorm:"size:3;not null;index:idx_p2p_fee_wallet_country_code"`
	PlatformFeeWalletAddress    string     `json:"platformFeeWalletAddress" gorm:"size:56;not null"`
	VatWalletAddress             string     `json:"vatWalletAddress" gorm:"size:56;not null"`
	RegulatoryFeeWalletAddress  string     `json:"regulatoryFeeWalletAddress" gorm:"size:56;not null"`
	EffectiveFrom               time.Time  `json:"effectiveFrom"`
	EffectiveTo                 *time.Time `json:"effectiveTo"`
	Version                     int        `json:"version" gorm:"not null;default:1"`
	Status                      string     `json:"status" gorm:"size:10;not null;default:'ACTIVE';index:idx_p2p_fee_wallet_status"`
	CreatedAt                   time.Time  `json:"createdAt"`
	UpdatedAt                   time.Time  `json:"updatedAt"`
}
