package db

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

type TokenizationReadClient interface {
	ListAssets(ctx context.Context, filters AssetFilters) (AssetPage, error)
	ListAllAssets(ctx context.Context, filters AssetFilters) ([]TokenizedAsset, error)
	GetAsset(ctx context.Context, idOrCode string) (*TokenizedAsset, error)
}

// TokenizationDetailReadClient is an optional richer read contract implemented
// by the production WalletDB client. Keeping it separate preserves lightweight
// fakes while allowing the portal to read Wallet-owned related records without
// duplicating or mutating them in AdminDB.
type TokenizationDetailReadClient interface {
	TokenizationReadClient
	GetAssetRelatedData(ctx context.Context, asset TokenizedAsset) (TokenizedAssetRelatedData, error)
}

type AssetFilters struct {
	Page                   int
	Limit                  int
	Status                 *int
	AssetType              string
	Sector                 string
	Search                 string
	DateFrom               *time.Time
	DateTo                 *time.Time
	AssetIDs               []string
	AssetCodes             []string
	ManagerStakeholderID   *uint64
	CustodianStakeholderID *uint64
	TrusteeStakeholderID   *uint64
	LegalAdviserID         *uint64
	FinancialAdviserID     *uint64
	IssuingHouseID         *uint64
	RatingAgencyID         *uint64
}

type AssetPage struct {
	Records []TokenizedAsset
	Total   int64
	Page    int
	Limit   int
}

type TokenizedAsset struct {
	ID                                       string    `gorm:"column:id" json:"id"`
	CreatedAt                                time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt                                time.Time `gorm:"column:updated_at" json:"updatedAt"`
	DateSubmitted                            time.Time `gorm:"column:date_submitted" json:"dateSubmitted"`
	DateOfApproval                           time.Time `gorm:"column:date_of_approval" json:"dateOfApproval"`
	MintingDate                              time.Time `gorm:"column:minting_date" json:"mintingDate"`
	AssetSector                              string    `gorm:"column:asset_sector" json:"assetSector"`
	AssetSubSector                           string    `gorm:"column:asset_sub_sector" json:"assetSubSector"`
	AssetType                                string    `gorm:"column:asset_type" json:"assetType"`
	AssetName                                string    `gorm:"column:asset_name" json:"assetName"`
	AssetCode                                string    `gorm:"column:asset_code" json:"assetCode"`
	AssetDescription                         string    `gorm:"column:asset_description" json:"assetDescription"`
	AssetLogo                                string    `gorm:"column:asset_logo" json:"assetLogo"`
	AssetWebsite                             string    `gorm:"column:asset_website" json:"assetWebsite"`
	OfferingType                             string    `gorm:"column:offering_type" json:"offeringType"`
	AssetCountryLocation                     string    `gorm:"column:asset_country_location" json:"assetCountryLocation"`
	AssetPhysicalAddress                     string    `gorm:"column:asset_physical_address" json:"assetPhysicalAddress"`
	AssetLongitude                           string    `gorm:"column:asset_longitude" json:"assetLongitude"`
	AssetLatitude                            string    `gorm:"column:asset_latitude" json:"assetLatitude"`
	OwnershipType                            string    `gorm:"column:ownership_type" json:"ownershipType"`
	OwnershipKind                            string    `gorm:"column:ownership_kind" json:"ownershipKind"`
	AssetOwnerName                           string    `gorm:"column:asset_owner_name" json:"assetOwnerName"`
	AssetOwnerAddress                        string    `gorm:"column:asset_owner_address" json:"assetOwnerAddress"`
	AssetQuoteCurrency                       string    `gorm:"column:asset_quote_currency" json:"assetQuoteCurrency"`
	AssetCurrentValue                        float64   `gorm:"column:asset_current_value" json:"assetCurrentValue"`
	ValueOfTokenizedAsset                    float64   `gorm:"column:value_of_tokenized_asset" json:"valueOfTokenizedAsset"`
	AssetOwnerRetainedOrContributedValue     float64   `gorm:"column:asset_owner_retained_or_contributed_value" json:"assetOwnerRetainedOrContributedValue"`
	AssetMiscCostOutsideValuation            float64   `gorm:"column:asset_msc_cost_outisde_of_valuation" json:"assetMscCostOutisdeOfValuation"`
	NumberOfTokenToBeIssued                  float64   `gorm:"column:number_of_token_to_be_issued" json:"numberOfTokenToBeIssued"`
	MaxNumberOfTokenAvailableForSale         float64   `gorm:"column:max_number_of_token_available_for_sale" json:"maxNumberOfTokenAvailableForSale"`
	NumberOfTokenToBeSold                    float64   `gorm:"column:number_of_token_to_be_sold" json:"numberOfTokenToBeSold"`
	TotalTokenHeldByManager                  float64   `gorm:"column:total_token_held_by_manager" json:"totalTokenHeldByManager"`
	PricePerToken                            float64   `gorm:"column:price_per_token" json:"pricePerToken"`
	SalesStart                               time.Time `gorm:"column:sales_start" json:"salesStart"`
	SalesEnd                                 time.Time `gorm:"column:sales_end" json:"salesEnd"`
	CapQuantity                              float64   `gorm:"column:cap_quantity" json:"capQuantity"`
	CapAmountInFiat                          float64   `gorm:"column:cap_amount_in_fiat" json:"capAmountInFiat"`
	CapDurationInDays                        int       `gorm:"column:cap_duration_in_days" json:"capDurationInDays"`
	ProceedCycle                             string    `gorm:"column:proceed_cycle" json:"proceedCycle"`
	ProceedPayoutCurrency                    string    `gorm:"column:proceed_payout_currency" json:"proceedPayoutCurrency"`
	ProceedPayoutType                        int       `gorm:"column:proceed_payout_type" json:"proceedPayoutType"`
	AssetTokenizationStatus                  int       `gorm:"column:asset_tokenization_status" json:"assetTokenizationStatus"`
	AssetManagerID                           uint64    `gorm:"column:asset_manager_id" json:"assetManagerId"`
	ApprovedAssetCustodianID                 uint64    `gorm:"column:approved_asset_custodian_id" json:"approvedAssetCustodianId"`
	TrusteeID                                uint64    `gorm:"column:trustee_id" json:"trusteeId"`
	LegalAdviserID                           uint64    `gorm:"column:legal_adviser_id" json:"legalAdviserId"`
	FinancialAdviserID                       uint64    `gorm:"column:financial_adviser_id" json:"financialAdviserId"`
	AssetIssuingHouseID                      uint64    `gorm:"column:asset_issuing_house_id" json:"assetIssuingHouseId"`
	RatingAgencyID                           uint64    `gorm:"column:rating_agency_id" json:"ratingAgencyId"`
	CustodianFeeValue                        float64   `gorm:"column:custodian_fee_value" json:"custodianFeeValue"`
	FeeInFiat                                float64   `gorm:"column:fee_in_fiat" json:"feeInFiat"`
	FeeInAsset                               float64   `gorm:"column:fee_in_asset" json:"feeInAsset"`
	FeeInAssetPercent                        float64   `gorm:"column:fee_in_asset_percent" json:"feeInAssetPercent"`
	TokenizationFeeID                        *uint64   `gorm:"column:tokenization_fee_id" json:"tokenizationFeeId"`
	SECTokenizationFeeValue                  float64   `gorm:"column:sec_tokenization_fee_value" json:"SECTokenizationFeeValue"`
	AssetManagerFeeValue                     float64   `gorm:"column:asset_manager_fee_value" json:"assetManagerFeeValue"`
	IssuingHouseFeeValue                     float64   `gorm:"column:issuing_house_fee_value" json:"issuingHouseFeeValue"`
	LegalAndProfessionalFeeValue             float64   `gorm:"column:legal_and_professional_fee_value" json:"legalAndProfessionalFeeValue"`
	RatingAgencyFeeValue                     float64   `gorm:"column:rating_agency_fee_value" json:"ratingAgencyFeeValue"`
	TrusteeFeeValue                          float64   `gorm:"column:trustee_fee_value" json:"trusteeFeeValue"`
	VATValue                                 float64   `gorm:"column:vat_value" json:"vatValue"`
	InitiatorUsername                        string    `gorm:"column:initiator_username" json:"initiatorUsername"`
	IssuingWalletAddress                     string    `gorm:"column:issuing_wallet_address" json:"issuingWalletAddress"`
	IssuingWalletAlias                       string    `gorm:"column:issuing_wallet_alias" json:"issuingWalletAlias"`
	MarketMakingWallet                       string    `gorm:"column:market_making_wallet" json:"marketMakingWallet"`
	InitialOwnerPreferredWalletAddress       string    `gorm:"column:initial_owner_preferred_wallet_address" json:"initialOwnerPreferredWalletAddress"`
	WalletToHoldAssetsNotForSale             string    `gorm:"column:wallet_to_hold_assets_not_for_sale" json:"walletToHoldAssetsNotForSale"`
	ExemptedCountries                        string    `gorm:"column:exempted_countries" json:"exemptedCountries"`
	ProtectionMethods                        string    `gorm:"column:protection_methods" json:"protectionMethods"`
	InsuranceCompanyName                     string    `gorm:"column:insurance_company_name" json:"insuranceCompanyName"`
	InsurancePolicyNumber                    string    `gorm:"column:insurance_policy_number" json:"insurancePolicyNumber"`
	InsurancePolicyHolder                    string    `gorm:"column:insurance_policy_holder" json:"insurancePolicyHolder"`
	PercentageValueOfInsurance               float64   `gorm:"column:percentage_value_of_insurance" json:"percentageValueOfInsurance"`
	IsFreeFromLiensAndEncumbrances           int       `gorm:"column:is_free_from_liens_and_encumbrances" json:"isFreeFromLiensAndEncumbrances"`
	AgreeTransferTitleToCustodian            int       `gorm:"column:agree_transfer_title_to_custodian" json:"agreeTransferTitleToCustodian"`
	ContractualProtectionRevGuarantees       int       `gorm:"column:contractual_protection_rev_guarantees" json:"contractualProtectionRevGuarantees"`
	ContractualProtectionPerfBond            int       `gorm:"column:contractual_protection_perf_bond" json:"contractualProtectionPerfBond"`
	ContractualProtectionSLA                 int       `gorm:"column:contractual_protection_sla" json:"contractualProtectionSLA"`
	RiskSharingMechanismPPPs                 int       `gorm:"column:risk_sharing_mechanism_pp_ps" json:"riskSharingMechanismPPPs"`
	RiskSharingMechanismHedgeInstruments     int       `gorm:"column:risk_sharing_mechanism_hedge_instruments" json:"riskSharingMechanismHedgeInstruments"`
	RiskSharingMechanismCompletionGuarantees int       `gorm:"column:risk_sharing_mechanism_completion_guarantees" json:"riskSharingMechanismCompletionGuarantees"`
	IndependentMonitoringList                string    `gorm:"column:independent_monitoring_list" json:"independentMonitoringList"`
	ESGSafeguardsSustainabilityCerts         int       `gorm:"column:esg_safeguards_sus_certs" json:"eSGSafeguardsSusCerts"`
	ESGSafeguardsCommunityEngagementPlan     int       `gorm:"column:esg_safeguards_comm_eng_plans" json:"eSGSafeguardsCommEngPlans"`
	SecurityMeasuresSurveillanceSystems      int       `gorm:"column:security_measures_surveilance_systems" json:"securityMeasuresSurveilanceSystems"`
	SecurityMeasuresOnSitePersonnel          int       `gorm:"column:security_measures_on_site_security_personnel" json:"securityMeasuresOnSiteSecurityPersonnel"`
	SecurityMeasuresPerimeterSecurity        int       `gorm:"column:security_measures_perimeter_security" json:"securityMeasuresPerimeterSecurity"`
	SecurityMeasuresCriticalInfrastructure   int       `gorm:"column:security_measures_critical_infra_protections" json:"securityMeasuresCriticalInfraProtections"`
	LegalAdvisor                             string    `gorm:"column:legal_advisor" json:"legalAdvisor"`
	FinancialAdvisor                         string    `gorm:"column:financial_advisor" json:"financialAdvisor"`
	OtherAssetProtection                     string    `gorm:"column:other_asset_protection" json:"otherAssetProtection"`
	BankID                                   *uint64   `gorm:"column:bank_id" json:"bankId"`
	AccountNumber                            string    `gorm:"column:account_number" json:"accountNumber"`
	BeneficiaryName                          string    `gorm:"column:beneficiary_name" json:"beneficiaryName"`
	TokenizationApplicationFee               float64   `gorm:"column:tokenization_application_fee" json:"tokenizationApplicationFee"`
	TokenizationApplicationFeeAsset          string    `gorm:"column:tokenization_application_fee_asset" json:"tokenizationApplicationFeeAsset"`
	VettingStatus                            int       `gorm:"column:vetting_status" json:"vettingStatus"`
	DueDiligenceFail                         int       `gorm:"column:due_diligence_fail" json:"dueDiligenceFail"`
	TokenHolderCount                         int64     `gorm:"->;column:token_holder_count" json:"tokenHolderCount"`
}

type AssetTokenizationDocument struct {
	ID               uint64    `gorm:"column:id;primaryKey;autoIncrement:false" json:"id"`
	CreatedAt        time.Time `gorm:"column:created_at" json:"created_at"`
	TokenizedAssetID string    `gorm:"column:tokenized_asset_id" json:"tokenized_asset_id"`
	DocumentType     string    `gorm:"column:document_type" json:"document_type"`
	DocumentTitle    string    `gorm:"column:document_title" json:"document_title"`
	DocumentURL      string    `gorm:"column:document_url" json:"document_url"`
	ShowToPublic     int       `gorm:"column:show_to_public" json:"show_to_public"`
}

func (AssetTokenizationDocument) TableName() string { return "asset_tokenization_documents" }

type TokenizationFeeProofOfPayment struct {
	ID                             uint64    `gorm:"column:id;primaryKey;autoIncrement:false" json:"id"`
	CreatedAt                      time.Time `gorm:"column:created_at" json:"created_at"`
	TokenizationFeePaymentMethodID string    `gorm:"column:tokenization_fee_payment_method_id" json:"payment_method_id"`
	TokenizedAssetID               string    `gorm:"column:tokenized_asset_id" json:"tokenized_asset_id"`
	TransactionReference           string    `gorm:"column:transaction_reference" json:"transaction_reference"`
	DocumentURL                    string    `gorm:"column:document_url" json:"document_url"`
}

func (TokenizationFeeProofOfPayment) TableName() string { return "tokenization_fee_proof_of_payments" }

type WalletBank struct {
	ID          uint64 `gorm:"column:id;primaryKey;autoIncrement:false" json:"id"`
	BankName    string `gorm:"column:bank_name" json:"bank_name"`
	CountryCode string `gorm:"column:country_code" json:"country_code"`
}

func (WalletBank) TableName() string { return "banks" }

type WalletTokenizationFee struct {
	ID                 uint64  `gorm:"column:id;primaryKey;autoIncrement:false" json:"id"`
	FeeFiatPercentage  float64 `gorm:"column:fee_fiat_percentage" json:"fee_fiat_percentage"`
	FeeFiatCap         float64 `gorm:"column:fee_fiat_cap" json:"fee_fiat_cap"`
	FeeAssetPercentage float64 `gorm:"column:fee_asset_percentage" json:"fee_asset_percentage"`
	FeeDescription     string  `gorm:"column:fee_description" json:"fee_description"`
}

func (WalletTokenizationFee) TableName() string { return "tokenization_fees" }

type WalletTokenizationStatus struct {
	ID          uint64 `gorm:"column:id;primaryKey;autoIncrement:false" json:"id"`
	Description string `gorm:"column:description" json:"description"`
}

func (WalletTokenizationStatus) TableName() string { return "tokenization_statuses" }

// WalletAssetType is the wallet's id->name lookup for the asset_type column stored
// on a tokenized asset. Used to resolve the numeric asset_type into a human label
// so the stakeholder (org) detail endpoint can return it directly, instead of the
// org portal having to call the admin-only tokenization proxy for the type list.
type WalletAssetType struct {
	ID        uint64 `gorm:"column:id;primaryKey;autoIncrement:false" json:"id"`
	AssetType string `gorm:"column:asset_type" json:"assetType"`
}

func (WalletAssetType) TableName() string { return "tokenized_asset_types" }

type TokenizedAssetRelatedData struct {
	Documents     []AssetTokenizationDocument
	PaymentProofs []TokenizationFeeProofOfPayment
	Bank          *WalletBank
	PreferredFee  *WalletTokenizationFee
	StatusCatalog []WalletTokenizationStatus
	AssetTypeName string
}

func (a TokenizedAsset) ConfiguredFeeValue() float64 {
	return a.FeeInFiat + a.SECTokenizationFeeValue + a.CustodianFeeValue + a.AssetManagerFeeValue +
		a.IssuingHouseFeeValue + a.LegalAndProfessionalFeeValue + a.RatingAgencyFeeValue + a.TrusteeFeeValue + a.VATValue
}

func (c *GormTokenizationReadClient) ListAllAssets(ctx context.Context, filters AssetFilters) ([]TokenizedAsset, error) {
	query := c.db.WithContext(ctx).Model(&TokenizedAsset{})
	query = applyAssetFilters(query, filters)
	var records []TokenizedAsset
	if err := withTokenHolderCount(query).Order("created_at desc").Find(&records).Error; err != nil {
		return nil, err
	}
	return records, nil
}

func (TokenizedAsset) TableName() string {
	return "tokenized_assets"
}

type GormTokenizationReadClient struct {
	db *gorm.DB
}

func NewGormTokenizationReadClient(db *gorm.DB) *GormTokenizationReadClient {
	return &GormTokenizationReadClient{db: db}
}

func (c *GormTokenizationReadClient) ListAssets(ctx context.Context, filters AssetFilters) (AssetPage, error) {
	page, limit := normalizePageLimit(filters.Page, filters.Limit)
	query := c.db.WithContext(ctx).Model(&TokenizedAsset{})
	query = applyAssetFilters(query, filters)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return AssetPage{}, err
	}

	var records []TokenizedAsset
	if err := withTokenHolderCount(query).Order("created_at desc").Offset((page - 1) * limit).Limit(limit).Find(&records).Error; err != nil {
		return AssetPage{}, err
	}
	return AssetPage{Records: records, Total: total, Page: page, Limit: limit}, nil
}

func (c *GormTokenizationReadClient) GetAsset(ctx context.Context, idOrCode string) (*TokenizedAsset, error) {
	idOrCode = strings.TrimSpace(idOrCode)
	if idOrCode == "" {
		return nil, gorm.ErrRecordNotFound
	}
	var asset TokenizedAsset
	err := c.db.WithContext(ctx).
		Scopes(withTokenHolderCount).
		Where("id = ? OR asset_code = ?", idOrCode, idOrCode).
		First(&asset).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}
	return &asset, nil
}

func (c *GormTokenizationReadClient) GetAssetRelatedData(ctx context.Context, asset TokenizedAsset) (TokenizedAssetRelatedData, error) {
	result := TokenizedAssetRelatedData{
		Documents:     []AssetTokenizationDocument{},
		PaymentProofs: []TokenizationFeeProofOfPayment{},
		StatusCatalog: []WalletTokenizationStatus{},
	}
	if err := c.db.WithContext(ctx).Where("tokenized_asset_id = ?", asset.ID).Order("created_at asc").Find(&result.Documents).Error; err != nil {
		return result, err
	}
	if err := c.db.WithContext(ctx).Where("tokenized_asset_id = ?", asset.ID).Order("created_at asc").Find(&result.PaymentProofs).Error; err != nil {
		return result, err
	}
	if asset.BankID != nil && *asset.BankID > 0 {
		var bank WalletBank
		if err := c.db.WithContext(ctx).First(&bank, "id = ?", *asset.BankID).Error; err == nil {
			result.Bank = &bank
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return result, err
		}
	}
	if asset.TokenizationFeeID != nil && *asset.TokenizationFeeID > 0 {
		var fee WalletTokenizationFee
		if err := c.db.WithContext(ctx).First(&fee, "id = ?", *asset.TokenizationFeeID).Error; err == nil {
			result.PreferredFee = &fee
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return result, err
		}
	}
	if err := c.db.WithContext(ctx).Order("id asc").Find(&result.StatusCatalog).Error; err != nil {
		return result, err
	}
	// Resolve the human asset-type label from the wallet lookup so the org detail
	// endpoint returns it directly (asset.AssetType is a numeric id). Best-effort:
	// a miss leaves it empty and the response keeps the raw value.
	if strings.TrimSpace(asset.AssetType) != "" {
		var assetType WalletAssetType
		if err := c.db.WithContext(ctx).First(&assetType, "id = ?", strings.TrimSpace(asset.AssetType)).Error; err == nil {
			result.AssetTypeName = assetType.AssetType
		}
	}
	return result, nil
}

func withTokenHolderCount(query *gorm.DB) *gorm.DB {
	return query.Select(`tokenized_assets.*,
		(SELECT COUNT(DISTINCT wallet_address)
		 FROM tokenized_asset_subscriptions
		 WHERE tokenized_asset_id = tokenized_assets.id) AS token_holder_count`)
}

func applyAssetFilters(query *gorm.DB, filters AssetFilters) *gorm.DB {
	if filters.Status != nil {
		query = query.Where("asset_tokenization_status = ?", *filters.Status)
	}
	if strings.TrimSpace(filters.AssetType) != "" {
		query = query.Where("asset_type = ?", strings.TrimSpace(filters.AssetType))
	}
	if strings.TrimSpace(filters.Sector) != "" {
		query = query.Where("asset_sector = ?", strings.TrimSpace(filters.Sector))
	}
	if strings.TrimSpace(filters.Search) != "" {
		search := "%" + strings.TrimSpace(filters.Search) + "%"
		query = query.Where("(asset_name ILIKE ? OR asset_code ILIKE ?)", search, search)
	}
	if filters.DateFrom != nil {
		query = query.Where("created_at >= ?", *filters.DateFrom)
	}
	if filters.DateTo != nil {
		query = query.Where("created_at <= ?", *filters.DateTo)
	}

	visibilityClauses := make([]string, 0, 4)
	visibilityArgs := make([]interface{}, 0, 4)
	if filters.ManagerStakeholderID != nil {
		visibilityClauses = append(visibilityClauses, "asset_manager_id = ?")
		visibilityArgs = append(visibilityArgs, *filters.ManagerStakeholderID)
	}
	if filters.CustodianStakeholderID != nil {
		visibilityClauses = append(visibilityClauses, "approved_asset_custodian_id = ?")
		visibilityArgs = append(visibilityArgs, *filters.CustodianStakeholderID)
	}
	if filters.TrusteeStakeholderID != nil {
		visibilityClauses = append(visibilityClauses, "trustee_id = ?")
		visibilityArgs = append(visibilityArgs, *filters.TrusteeStakeholderID)
	}
	if filters.LegalAdviserID != nil {
		visibilityClauses = append(visibilityClauses, "legal_adviser_id = ?")
		visibilityArgs = append(visibilityArgs, *filters.LegalAdviserID)
	}
	if filters.FinancialAdviserID != nil {
		visibilityClauses = append(visibilityClauses, "financial_adviser_id = ?")
		visibilityArgs = append(visibilityArgs, *filters.FinancialAdviserID)
	}
	if filters.IssuingHouseID != nil {
		visibilityClauses = append(visibilityClauses, "asset_issuing_house_id = ?")
		visibilityArgs = append(visibilityArgs, *filters.IssuingHouseID)
	}
	if filters.RatingAgencyID != nil {
		visibilityClauses = append(visibilityClauses, "rating_agency_id = ?")
		visibilityArgs = append(visibilityArgs, *filters.RatingAgencyID)
	}
	if len(filters.AssetIDs) > 0 {
		visibilityClauses = append(visibilityClauses, "id IN ?")
		visibilityArgs = append(visibilityArgs, filters.AssetIDs)
	}
	if len(filters.AssetCodes) > 0 {
		visibilityClauses = append(visibilityClauses, "asset_code IN ?")
		visibilityArgs = append(visibilityArgs, filters.AssetCodes)
	}
	if len(visibilityClauses) > 0 {
		query = query.Where(strings.Join(visibilityClauses, " OR "), visibilityArgs...)
	}

	return query
}

func normalizePageLimit(page, limit int) (int, int) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	return page, limit
}

func ParseOptionalStatus(value string) (*int, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, nil
	}
	status, err := strconv.Atoi(value)
	if err != nil {
		return nil, err
	}
	return &status, nil
}
