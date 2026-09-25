package models

import (
	"time"
)

// TokenizationStakeholderType represents predefined stakeholder types
type TokenizationStakeholderType struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Type        string    `gorm:"unique;not null" json:"type"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TokenizationStakeholder represents a stakeholder and their fees
type TokenizationStakeholder struct {
	ID                uint                        `gorm:"primaryKey" json:"id"`
	StakeholderTypeID uint                        `gorm:"not null" json:"stakeholder_type_id"`
	StakeholderType   TokenizationStakeholderType `gorm:"foreignKey:StakeholderTypeID" json:"stakeholder_type"`
	StakeholderName   string                      `gorm:"not null" json:"stakeholder_name"`
	StakeholderEmail  string                      `gorm:"not null" json:"stakeholder_email"`
	Description       string                      `json:"description"`
	FeePercentage     float64                     `json:"fee_percentage"`
	FixedFee          float64                     `json:"fixed_fee"`
	CreatedAt         time.Time                   `json:"created_at"`
	UpdatedAt         time.Time                   `json:"updated_at"`
}

// StakeholderRequestDTO for filtering and pagination
type StakeholderRequestDTO struct {
	Page            int    `form:"page" json:"page"`
	PageSize        int    `form:"page_size" json:"page_size"`
	Search          string `form:"search" json:"search"`
	StakeholderName string `form:"stakeholder_name" json:"stakeholder_name"`
	StakeholderType uint   `form:"stakeholder_type" json:"stakeholder_type"`
	Email           string `form:"email" json:"email"`
}

// StakeholderResponse for paginated responses
type StakeholderResponse struct {
	Data     []TokenizationStakeholder `json:"data"`
	Total    int64                     `json:"total"`
	Page     int                       `json:"page"`
	PageSize int                       `json:"page_size"`
}

// StakeholderCreateRequest for creating a new stakeholder
type StakeholderCreateRequest struct {
	StakeholderTypeID uint    `json:"stakeholder_type_id" binding:"required"`
	StakeholderName   string  `json:"stakeholder_name" binding:"required"`
	StakeholderEmail  string  `json:"stakeholder_email" binding:"required,email"`
	Description       string  `json:"description"`
	FeePercentage     float64 `json:"fee_percentage"`
	FixedFee          float64 `json:"fixed_fee"`
}

// StakeholderUpdateRequest for updating a stakeholder
type StakeholderUpdateRequest struct {
	StakeholderTypeID uint    `json:"stakeholder_type_id"`
	StakeholderName   string  `json:"stakeholder_name"`
	StakeholderEmail  string  `json:"stakeholder_email" binding:"omitempty,email"`
	Description       string  `json:"description"`
	FeePercentage     float64 `json:"fee_percentage"`
	FixedFee          float64 `json:"fixed_fee"`
}

// ApprovedAssetCustodian represents an approved custodian
type ApprovedAssetCustodian struct {
	ID                    uint    `json:"id" gorm:"primaryKey"`
	AssetCustodianName    string  `json:"asset_custodian_name" example:"First Bank Custody"`
	AssetCustodianAddress string  `json:"asset_custodian_address" example:"Plot 5, Custody Street, VI"`
	AssetCustodianCountry string  `json:"asset_custodian_country" example:"NG"`
	RequirementDocument   string  `json:"requirement_document" example:"https://example.com/document.pdf"`
	FeePercent            float64 `json:"fee_percent" example:"0.25"`
	FeeFixed              float64 `json:"fee_fixed" example:"1000"`
}

// AssetIssuingHouse represents an issuing house
type AssetIssuingHouse struct {
	ID                       uint    `json:"id" gorm:"primaryKey"`
	AssetIssuingHouseName    string  `json:"asset_issuing_house_name" example:"Vetiva Capital"`
	AssetIssuingHouseAddress string  `json:"asset_issuing_house_address" example:"7th Floor, Marina Tower"`
	AssetIssuingHouseCountry string  `json:"asset_issuing_house_country" example:"NG"`
	FeePercent               float64 `json:"fee_percent" example:"0.25"`
	FeeFixed                 float64 `json:"fee_fixed" example:"1000"`
}

// SavePartnerRequest represents a flexible request format
type SavePartnerRequest struct {
	// Type of partner e.g. asset_manager, asset_issuing_house, approved_asset_custodian
	Type string `json:"type" binding:"required" example:"asset_manager"`

	// Action to perform: create or update
	Action string `json:"action" binding:"required" example:"create"`

	// Payload holds the partner details
	Payload map[string]interface{} `json:"payload" binding:"required"`
}

type AssetManager struct {
	ID                  int64   `json:"id" db:"id"`
	AssetManagerName    string  `json:"asset_manager_name" db:"asset_manager_name"`
	AssetManagerAddress string  `json:"asset_manager_address" db:"asset_manager_address"`
	AssetManagerCountry string  `json:"asset_manager_country" db:"asset_manager_country"`
	FeePercent          float64 `json:"fee_percent" db:"fee_percent"`
	FeeFixed            float64 `json:"fee_fixed" db:"fee_fixed"`
}

// AssetManagerRequest defines the expected data payload for Asset Manager operations
// Action and Type are provided via query parameters.
// @Description Data payload for Asset Manager. ID is required in the body for update actions.
// @Example {"asset_manager_name": "FundQuest Capital", "asset_manager_address": "123 Ikoyi Crescent", "asset_manager_country": "Nigeria", "fee_percent": 0.25, "fee_fixed": 1000}
// @Example (Update) {"id": 1, "asset_manager_name": "FundQuest Capital Updated", "asset_manager_address": "123 Ikoyi Crescent", "asset_manager_country": "Nigeria", "fee_percent": 0.30, "fee_fixed": 1100}
type AssetManagerRequest struct {
	// ID is required in the JSON body *only* for update actions.
	ID                  int64   `json:"id,omitempty"`
	AssetManagerName    string  `json:"asset_manager_name" binding:"required"`
	AssetManagerAddress string  `json:"asset_manager_address" binding:"required"`
	AssetManagerCountry string  `json:"asset_manager_country" binding:"required"`
	FeePercent          float64 `json:"fee_percent"`
	FeeFixed            float64 `json:"fee_fixed"`
}

// AssetIssuingHouseRequest defines the expected data payload for Issuing House operations
// Action and Type are provided via query parameters.
// @Description Data payload for Asset Issuing House. ID is required in the body for update actions.
// @Example {"asset_issuing_house_name": "Vetiva Capital", "asset_issuing_house_address": "7th Floor, Marina Tower", "asset_issuing_house_country": "NG", "fee_percent": 0.25, "fee_fixed": 1000}
// @Example (Update) {"id": 1, "asset_issuing_house_name": "Vetiva Capital Updated", "asset_issuing_house_address": "7th Floor, Marina Tower", "asset_issuing_house_country": "NG", "fee_percent": 0.30, "fee_fixed": 1100}
type AssetIssuingHouseRequest struct {
	// ID is required in the JSON body *only* for update actions.
	ID                       uint    `json:"id,omitempty"`
	AssetIssuingHouseName    string  `json:"asset_issuing_house_name" binding:"required"`
	AssetIssuingHouseAddress string  `json:"asset_issuing_house_address" binding:"required"`
	AssetIssuingHouseCountry string  `json:"asset_issuing_house_country" binding:"required"`
	FeePercent               float64 `json:"fee_percent"`
	FeeFixed                 float64 `json:"fee_fixed"`
}

// ApprovedAssetCustodianRequest defines the expected data payload for Approved Custodian operations
// Action and Type are provided via query parameters.
// @Description Data payload for Approved Asset Custodian. ID is required in the body for update actions.
// @Example {"asset_custodian_name": "First Bank Custody", "asset_custodian_address": "Plot 5, Custody Street, VI", "asset_custodian_country": "NG", "requirement_document": "https://example.com/doc.pdf", "fee_percent": "0.25", "fee_fixed": "1000"}
// @Example (Update) {"id": 1, "asset_custodian_name": "First Bank Custody Updated", "asset_custodian_address": "Plot 5, Custody Street, VI", "asset_custodian_country": "NG", "requirement_document": "https://example.com/doc_updated.pdf", "fee_percent": "0.30", "fee_fixed": "1100"}
type ApprovedAssetCustodianRequest struct {
	// ID is required in the JSON body *only* for update actions.
	ID                    uint   `json:"id,omitempty"`
	AssetCustodianName    string `json:"asset_custodian_name" binding:"required"`
	AssetCustodianAddress string `json:"asset_custodian_address" binding:"required"`
	AssetCustodianCountry string `json:"asset_custodian_country" binding:"required"`
	RequirementDocument   string `json:"requirement_document" binding:"required"`
	FeePercent            string `json:"fee_percent"` // Matches model's string type
	FeeFixed              string `json:"fee_fixed"`   // Matches model's string type
}

// StandardSuccessResponse represents a successful operation
// @Description standard success response
//
//	@Example {
//	  "message": "Success"
//	}
type StandardSuccessResponse struct {
	Message string `json:"message"`
}

type Country struct {
	ID                              int64    `json:"id"`
	CountryCode                     string   `json:"country_code"`
	SecTokenizationFeePercent       *float64 `json:"sec_tokenization_fee_percent"`
	SecTradeFeePercent              *float64 `json:"sec_trade_fee_percent"`
	RegionName                      *string  `json:"region_name"`
	SecTokenizationFee              *float64 `json:"sec_tokenization_fee"`
	SecTradeFee                     *float64 `json:"sec_trade_fee"`
	MinTokenizationFee              *float64 `json:"min_tokenization_fee"`
	TokenizationApplicationFee      *float64 `json:"tokenization_application_fee"`
	TokenizationApplicationFeeAsset *string  `json:"tokenization_application_fee_asset"`
	CountryName                     *string  `json:"country_name"`
	QuoteCurrencyCode               *string  `json:"quote_currency_code"`
	LegalAndProfessionalFeePercent  *float64 `json:"legal_and_professional_fee_percent"`
	RatingAgencyFeePercent          *float64 `json:"rating_agency_fee_percent"`
	VatPercent                      *float64 `json:"vat_percent"`
	RegulatorName                   *string  `json:"regulator_name"`
	SecTokenizationFeeFixed         *float64 `json:"sec_tokenization_fee_fixed"`
	SecTradeFeeFixed                *float64 `json:"sec_trade_fee_fixed"`
	LegalAndProfessionalFeeFixed    *float64 `json:"legal_and_professional_fee_fixed"`
	RatingAgencyFeeFixed            *float64 `json:"rating_agency_fee_fixed"`
	FiatLabel                       *string  `json:"fiat_label"`
	FiatGlyph                       *string  `json:"fiat_glyph"`
}

type CountryRequest struct {
	Action                                   string  `json:"action" binding:"required"` // "create" or "update"
	CountryCode                              string  `json:"country_code" binding:"required"`
	RegionName                               string  `json:"region_name"`
	CountryName                              string  `json:"country_name"`
	QuoteCurrencyCode                        string  `json:"quote_currency_code"`
	RegulatorName                            string  `json:"regulator_name"`
	FiatLabel                                string  `json:"fiat_label"`
	FiatGlyph                                string  `json:"fiat_glyph"`
	SecTokenizationFeePercent                float64 `json:"sec_tokenization_fee_percent"`
	SecTradeFeePercent                       float64 `json:"sec_trade_fee_percent"`
	SecTokenizationFee                       float64 `json:"sec_tokenization_fee"`
	SecTradeFee                              float64 `json:"sec_trade_fee"`
	MinTokenizationFee                       float64 `json:"min_tokenization_fee"`
	TokenizationApplicationFee               float64 `json:"tokenization_application_fee"`
	TokenizationApplicationFeeAsset          string  `json:"tokenization_application_fee_asset"`
	LegalAndProfessionalFeePercent           float64 `json:"legal_and_professional_fee_percent"`
	RatingAgencyFeePercent                   float64 `json:"rating_agency_fee_percent"`
	VATPercent                               float64 `json:"vat_percent"`
	SecTokenizationFeeFixed                  float64 `json:"sec_tokenization_fee_fixed"`
	SecTradeFeeFixed                         float64 `json:"sec_trade_fee_fixed"`
	LegalAndProfessionalFeeFixed             float64 `json:"legal_and_professional_fee_fixed"`
	RatingAgencyFeeFixed                     float64 `json:"rating_agency_fee_fixed"`
	MinTrovBalanceForTokenizationApplication float64 `json:"min_trov_balance_for_tokenization_application"`
	FiatActivationAmount                     float64 `json:"fiat_activation_amount"`
	TrovTokenActivationPercent               float64 `json:"trov_token_activation_percent"`
}

type CountryResponse struct {
	// No ID: the wallet-owned `countries` table is keyed on country_code.
	CountryCode                              string  `json:"country_code"`
	RegionName                               string  `json:"region_name"`
	CountryName                              string  `json:"country_name"`
	QuoteCurrencyCode                        string  `json:"quote_currency_code"`
	RegulatorName                            string  `json:"regulator_name"`
	FiatLabel                                string  `json:"fiat_label"`
	FiatGlyph                                string  `json:"fiat_glyph"`
	SecTokenizationFeePercent                float64 `json:"sec_tokenization_fee_percent"`
	SecTradeFeePercent                       float64 `json:"sec_trade_fee_percent"`
	SecTokenizationFee                       float64 `json:"sec_tokenization_fee"`
	SecTradeFee                              float64 `json:"sec_trade_fee"`
	MinTokenizationFee                       float64 `json:"min_tokenization_fee"`
	TokenizationApplicationFee               float64 `json:"tokenization_application_fee"`
	TokenizationApplicationFeeAsset          string  `json:"tokenization_application_fee_asset"`
	LegalAndProfessionalFeePercent           float64 `json:"legal_and_professional_fee_percent"`
	RatingAgencyFeePercent                   float64 `json:"rating_agency_fee_percent"`
	VATPercent                               float64 `json:"vat_percent"`
	SecTokenizationFeeFixed                  float64 `json:"sec_tokenization_fee_fixed"`
	SecTradeFeeFixed                         float64 `json:"sec_trade_fee_fixed"`
	LegalAndProfessionalFeeFixed             float64 `json:"legal_and_professional_fee_fixed"`
	RatingAgencyFeeFixed                     float64 `json:"rating_agency_fee_fixed"`
	MinTrovBalanceForTokenizationApplication float64 `json:"min_trov_balance_for_tokenization_application"`
	FiatActivationAmount                     float64 `json:"fiat_activation_amount"`
	TrovTokenActivationPercent               float64 `json:"trov_token_activation_percent"`
}

type LegalAndProfessionalsPartners struct {
	ID             uint    `json:"id,omitempty" gorm:"primaryKey"`
	PartnerName    string  `json:"partner_name" gorm:"not null"`
	PartnerAddress string  `json:"partner_address" gorm:"not null"`
	PartnerCountry string  `json:"partner_country" gorm:"not null"`
	FeePercent     float64 `json:"fee_percent" gorm:"not null"`
	FeeFixed       float64 `json:"fee_fixed" gorm:"not null"`
}

// set table name to legal_and_profesional_partners
func (LegalAndProfessionalsPartners) TableName() string {
	return "legal_and_profesional_partners"
}

type LegalAndProfessionalsRequest struct {
	ID             uint   `json:"id,omitempty"`
	PartnerName    string `json:"partner_name"`
	PartnerAddress string `json:"partner_address"`
	PartnerCountry string `json:"partner_country"`
	FeePercent     string `json:"fee_percent"`
	FeeFixed       string `json:"fee_fixed"`
}

type RatingAgency struct {
	ID            uint    `json:"id,omitempty" gorm:"primaryKey"`
	AgencyName    string  `json:"agency_name" gorm:"not null"`
	AgencyAddress string  `json:"agency_address" gorm:"not null"`
	AgencyCountry string  `json:"agency_country" gorm:"not null"`
	FeePercent    float64 `json:"fee_percent" gorm:"not null"`
	FeeFixed      float64 `json:"fee_fixed" gorm:"not null"`
}

// TableName set table name to rating_agencies
func (RatingAgency) TableName() string {
	return "rating_agencies"
}

type RatingAgencyRequest struct {
	ID            uint   `json:"id,omitempty"`
	AgencyName    string `json:"agency_name"`
	AgencyAddress string `json:"agency_address"`
	AgencyCountry string `json:"agency_country"`
	FeePercent    string `json:"fee_percent"`
	FeeFixed      string `json:"fee_fixed"`
}

type Trustees struct {
	ID             uint    `json:"id,omitempty" gorm:"primaryKey"`
	TrusteeName    string  `json:"trustee_name" gorm:"not null"`
	TrusteeAddress string  `json:"trustee_address" gorm:"not null"`
	TrusteeCountry string  `json:"trustee_country" gorm:"not null"`
	FeePercent     float64 `json:"fee_percent" gorm:"not null"`
	FeeFixed       float64 `json:"fee_fixed" gorm:"not null"`
}

// TableName sets table name for Trustees
func (Trustees) TableName() string {
	return "trustees"
}

type TrusteesRequest struct {
	ID             uint   `json:"id,omitempty"`
	TrusteeName    string `json:"trustee_name"`
	TrusteeAddress string `json:"trustee_address"`
	TrusteeCountry string `json:"trustee_country"`
	FeePercent     string `json:"fee_percent"`
	FeeFixed       string `json:"fee_fixed"`
}

// LegalAdviser is a distinct A5 stakeholder role that prepares legal documentation.
// Separate from LegalAndProfessionalsPartners.
type LegalAdviser struct {
	ID             uint    `json:"id,omitempty" gorm:"primaryKey"`
	AdviserName    string  `json:"adviser_name" gorm:"not null"`
	AdviserAddress string  `json:"adviser_address" gorm:"not null"`
	AdviserCountry string  `json:"adviser_country" gorm:"not null"`
	FeePercent     float64 `json:"fee_percent" gorm:"not null"`
	FeeFixed       float64 `json:"fee_fixed" gorm:"not null"`
}

// TableName sets table name for LegalAdviser
func (LegalAdviser) TableName() string {
	return "legal_advisers"
}

type LegalAdviserRequest struct {
	ID             uint   `json:"id,omitempty"`
	AdviserName    string `json:"adviser_name"`
	AdviserAddress string `json:"adviser_address"`
	AdviserCountry string `json:"adviser_country"`
	FeePercent     string `json:"fee_percent"`
	FeeFixed       string `json:"fee_fixed"`
}

// FinancialAdviser is a distinct A5 stakeholder role that structures the financial
// terms of the tokenization. Separate from LegalAndProfessionalsPartners.
type FinancialAdviser struct {
	ID             uint    `json:"id,omitempty" gorm:"primaryKey"`
	AdviserName    string  `json:"adviser_name" gorm:"not null"`
	AdviserAddress string  `json:"adviser_address" gorm:"not null"`
	AdviserCountry string  `json:"adviser_country" gorm:"not null"`
	FeePercent     float64 `json:"fee_percent" gorm:"not null"`
	FeeFixed       float64 `json:"fee_fixed" gorm:"not null"`
}

// TableName sets table name for FinancialAdviser
func (FinancialAdviser) TableName() string {
	return "financial_advisers"
}

type FinancialAdviserRequest struct {
	ID             uint   `json:"id,omitempty"`
	AdviserName    string `json:"adviser_name"`
	AdviserAddress string `json:"adviser_address"`
	AdviserCountry string `json:"adviser_country"`
	FeePercent     string `json:"fee_percent"`
	FeeFixed       string `json:"fee_fixed"`
}
