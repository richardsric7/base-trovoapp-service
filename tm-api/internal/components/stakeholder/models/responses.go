package models

import "time"

type PaginationMeta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

type PaginatedResponse struct {
	Records interface{}      `json:"records"`
	Meta    PaginationMeta   `json:"meta"`
	Summary *AssetStatistics `json:"summary,omitempty"`
}

type AssetStatistics struct {
	Total      int `json:"total"`
	Active     int `json:"active"`
	Liquidated int `json:"liquidated"`
}

type StakeholderProfileResponse struct {
	Member       StakeholderProfileMember       `json:"member"`
	Organization StakeholderProfileOrganization `json:"organization"`
	Stakeholder  StakeholderProfileLinkage      `json:"stakeholder"`
	Wallet       StakeholderProfileWallet       `json:"wallet"`
}

type StakeholderProfileMember struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	Status    string `json:"status"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
}

type StakeholderProfileOrganization struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Type   string `json:"type"`
	Status string `json:"status"`
}

type StakeholderProfileLinkage struct {
	ID            *uint64 `json:"stakeholder_id,omitempty"`
	Type          string  `json:"stakeholder_type"`
	DashboardRole string  `json:"dashboard_role"`
}

type StakeholderProfileWallet struct {
	TrovoWalletUsername *string `json:"trovo_wallet_username,omitempty"`
	IsWalletLinked      bool    `json:"is_wallet_linked"`
}

type AssetActionFlags struct {
	CanApproveDueDiligence   bool `json:"canApproveDueDiligence"`
	CanApproveFundRelease    bool `json:"canApproveFundRelease"`
	CanAuthorizeDistribution bool `json:"canAuthorizeDistribution"`
	CanExecuteFundRelease    bool `json:"canExecuteFundRelease"`
	CanUpdateAccount         bool `json:"canUpdateAccount"`
	CanRequestFundRelease    bool `json:"canRequestFundRelease"`
	CanRecordRevenue         bool `json:"canRecordRevenue"`
	CanSubmitValuation       bool `json:"canSubmitValuation"`
}

type TokenizedAssetResponse struct {
	ID                    string                      `json:"id"`
	AssetCode             string                      `json:"assetCode"`
	AssetName             string                      `json:"assetName"`
	AssetType             string                      `json:"assetType"`
	AssetSector           string                      `json:"assetSector"`
	AssetSubSector        string                      `json:"assetSubSector"`
	Status                int                         `json:"assetTokenizationStatus"`
	AssetQuoteCurrency    string                      `json:"assetQuoteCurrency"`
	AssetCurrentValue     float64                     `json:"assetCurrentValue"`
	ValueOfTokenizedAsset float64                     `json:"valueOfTokenizedAsset"`
	NumberOfTokensIssued  float64                     `json:"numberOfTokenToBeIssued"`
	AssetManagerID        uint64                      `json:"assetManagerId"`
	ApprovedCustodianID   uint64                      `json:"approvedAssetCustodianId"`
	TrusteeID             uint64                      `json:"trusteeId"`
	TokenHolderCount      int64                       `json:"tokenHolderCount"`
	ComplianceStatus      string                      `json:"complianceStatus"`
	CustodyStatus         string                      `json:"custodyStatus"`
	CustodianFeeValue     float64                     `json:"-"`
	ConfiguredFeeValue    float64                     `json:"-"`
	CreatedAt             time.Time                   `json:"createdAt"`
	UpdatedAt             time.Time                   `json:"updatedAt"`
	Actions               AssetActionFlags            `json:"actions"`
	Assignment            *StakeholderAssetAssignment `json:"assignment,omitempty"`
}

type AssignedStakeholderResponse struct {
	Role              string  `json:"role"`
	StakeholderID     uint64  `json:"stakeholder_id"`
	StakeholderType   string  `json:"stakeholder_type"`
	OrganizationID    *string `json:"organization_id,omitempty"`
	OrganizationName  string  `json:"organization_name,omitempty"`
	OrganizationEmail string  `json:"organization_email,omitempty"`
}

type AssetWalletDetails struct {
	IssuingWalletPublicKey             string `json:"issuing_wallet_public_key,omitempty"`
	IssuingWalletAlias                 string `json:"issuing_wallet_alias,omitempty"`
	MarketMakingWallet                 string `json:"market_making_wallet,omitempty"`
	InitialOwnerPreferredWalletAddress string `json:"initial_owner_preferred_wallet_address,omitempty"`
	WalletToHoldAssetsNotForSale       string `json:"wallet_to_hold_assets_not_for_sale,omitempty"`
}

type TokenizedAssetDetailResponse struct {
	TokenizedAssetResponse
	TokenizerUsername     string                        `json:"tokenizer_username,omitempty"`
	AssignedStakeholders  []AssignedStakeholderResponse `json:"assigned_stakeholders"`
	Wallets               AssetWalletDetails            `json:"wallets"`
	ExemptedCountries     []string                      `json:"exempted_countries"`
	Documents             []StakeholderDocument         `json:"documents"`
	AssetProfile          AssetProfileDetails           `json:"asset_profile"`
	Ownership             AssetOwnershipDetails         `json:"ownership"`
	RiskAndCompliance     AssetRiskAndComplianceDetails `json:"risk_and_compliance"`
	Protection            AssetProtectionDetails        `json:"asset_protection"`
	Offering              AssetOfferingDetails          `json:"offering"`
	Fees                  AssetTokenizationFeeDetails   `json:"tokenization_fees"`
	ReceivingAccount      AssetReceivingAccountDetails  `json:"receiving_account"`
	ApplicationFee        AssetApplicationFeeDetails    `json:"application_fee"`
	TokenizationDocuments []WalletAssetDocumentResponse `json:"tokenization_documents"`
	Timeline              AssetTokenizationTimeline     `json:"tokenization_timeline"`
	OperationalUpdate     *StakeholderAssetOperation    `json:"operational_update,omitempty"`
}

type AssetRiskAndComplianceDetails struct {
	RiskAssessmentScore *float64   `json:"risk_assessment_score"`
	MaximumRiskScore    float64    `json:"maximum_risk_score"`
	RiskLevel           *string    `json:"risk_level"`
	ComplianceStatus    string     `json:"compliance_status"`
	LastAuditDate       *time.Time `json:"last_audit_date"`
}

type AssetProfileDetails struct {
	Description     string `json:"description,omitempty"`
	LogoURL         string `json:"logo_url,omitempty"`
	Website         string `json:"website,omitempty"`
	OfferingType    string `json:"offering_type,omitempty"`
	Country         string `json:"country,omitempty"`
	PhysicalAddress string `json:"physical_address,omitempty"`
	Longitude       string `json:"longitude,omitempty"`
	Latitude        string `json:"latitude,omitempty"`
}

type AssetOwnershipDetails struct {
	Type                       string  `json:"type,omitempty"`
	Kind                       string  `json:"kind,omitempty"`
	OwnerName                  string  `json:"owner_name,omitempty"`
	OwnerAddress               string  `json:"owner_address,omitempty"`
	RetainedOrContributedValue float64 `json:"retained_or_contributed_value"`
	CostOutsideValuation       float64 `json:"cost_outside_valuation"`
}

type AssetProtectionDetails struct {
	Methods                          []string `json:"methods"`
	InsuranceCompanyName             string   `json:"insurance_company_name,omitempty"`
	InsurancePolicyNumber            string   `json:"insurance_policy_number,omitempty"`
	InsurancePolicyHolder            string   `json:"insurance_policy_holder,omitempty"`
	InsuranceCoveragePercentage      float64  `json:"insurance_coverage_percentage"`
	FreeFromLiensAndEncumbrances     bool     `json:"free_from_liens_and_encumbrances"`
	TransferTitleToCustodian         bool     `json:"transfer_title_to_custodian"`
	RevenueGuarantees                bool     `json:"revenue_guarantees"`
	PerformanceBond                  bool     `json:"performance_bond"`
	ServiceLevelAgreement            bool     `json:"service_level_agreement"`
	CompletionGuarantees             bool     `json:"completion_guarantees"`
	PublicPrivatePartnerships        bool     `json:"public_private_partnerships"`
	HedgeInstruments                 bool     `json:"hedge_instruments"`
	IndependentMonitoringList        string   `json:"independent_monitoring_list,omitempty"`
	ESGSustainabilityCertifications  bool     `json:"esg_sustainability_certifications"`
	ESGCommunityEngagementPlan       bool     `json:"esg_community_engagement_plan"`
	SurveillanceSystems              bool     `json:"surveillance_systems"`
	OnSiteSecurityPersonnel          bool     `json:"on_site_security_personnel"`
	PerimeterSecurity                bool     `json:"perimeter_security"`
	CriticalInfrastructureProtection bool     `json:"critical_infrastructure_protection"`
	LegalAdviser                     string   `json:"legal_adviser,omitempty"`
	FinancialAdviser                 string   `json:"financial_adviser,omitempty"`
	Other                            string   `json:"other,omitempty"`
}

type AssetOfferingDetails struct {
	TokensIssued         float64   `json:"tokens_issued"`
	TokensForSale        float64   `json:"tokens_for_sale"`
	TokensNotForSale     float64   `json:"tokens_not_for_sale"`
	MaximumTokensForSale float64   `json:"maximum_tokens_for_sale"`
	PricePerToken        float64   `json:"price_per_token"`
	AmountToBeRaised     float64   `json:"amount_to_be_raised"`
	SalesStart           time.Time `json:"sales_start,omitempty"`
	SalesEnd             time.Time `json:"sales_end,omitempty"`
	CapQuantity          float64   `json:"cap_quantity"`
	CapAmount            float64   `json:"cap_amount"`
	CapDurationDays      int       `json:"cap_duration_days"`
	PayoutCycle          string    `json:"payout_cycle,omitempty"`
	PayoutCurrency       string    `json:"payout_currency,omitempty"`
	PayoutType           int       `json:"payout_type"`
}

type PreferredTokenizationFee struct {
	ID              uint64  `json:"id"`
	Description     string  `json:"description,omitempty"`
	FiatPercentage  float64 `json:"fiat_percentage"`
	FiatCap         float64 `json:"fiat_cap"`
	AssetPercentage float64 `json:"asset_percentage"`
}

type AssetTokenizationFeeDetails struct {
	FiatAmount      float64                   `json:"fiat_amount"`
	AssetAmount     float64                   `json:"asset_amount"`
	AssetPercentage float64                   `json:"asset_percentage"`
	Preferred       *PreferredTokenizationFee `json:"preferred,omitempty"`
}

type AssetReceivingAccountDetails struct {
	BankID        *uint64 `json:"bank_id,omitempty"`
	BankName      string  `json:"bank_name,omitempty"`
	BankCountry   string  `json:"bank_country,omitempty"`
	AccountName   string  `json:"account_name,omitempty"`
	AccountNumber string  `json:"account_number,omitempty"`
}

type AssetApplicationFeePaymentResponse struct {
	ID                   uint64    `json:"id"`
	PaymentMethodID      string    `json:"payment_method_id,omitempty"`
	TransactionReference string    `json:"transaction_reference,omitempty"`
	DocumentURL          string    `json:"document_url,omitempty"`
	CreatedAt            time.Time `json:"created_at,omitempty"`
}

type AssetApplicationFeeDetails struct {
	Amount   float64                              `json:"amount"`
	Asset    string                               `json:"asset,omitempty"`
	Payments []AssetApplicationFeePaymentResponse `json:"payments"`
}

type WalletAssetDocumentResponse struct {
	ID           uint64    `json:"id"`
	DocumentType string    `json:"document_type,omitempty"`
	Title        string    `json:"title,omitempty"`
	URL          string    `json:"url,omitempty"`
	Public       bool      `json:"public"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
}

type AssetTokenizationTimelineStage struct {
	StatusID    uint64 `json:"status_id"`
	Description string `json:"description,omitempty"`
	Reached     bool   `json:"reached"`
	Current     bool   `json:"current"`
}

type AssetTokenizationTimeline struct {
	CurrentStatusID uint64                           `json:"current_status_id"`
	DateSubmitted   *time.Time                       `json:"date_submitted,omitempty"`
	DateApproved    *time.Time                       `json:"date_approved,omitempty"`
	MintingDate     *time.Time                       `json:"minting_date,omitempty"`
	Stages          []AssetTokenizationTimelineStage `json:"stages"`
}

type AuthorizationChallengeResponse struct {
	ID          string    `json:"id"`
	AuthID      string    `json:"auth_id"`
	DynamicLink string    `json:"dynamic_link,omitempty"`
	QRCode      string    `json:"qr_code,omitempty"`
	Status      string    `json:"status"`
	ExpiresAt   time.Time `json:"expires_at"`
}

type DueDiligenceCategoryResponse struct {
	Category           string             `json:"category"`
	Status             string             `json:"status"`
	Notes              string             `json:"notes,omitempty"`
	VerifiedByMemberID *string            `json:"verified_by_member_id,omitempty"`
	VerifiedAt         *time.Time         `json:"verified_at,omitempty"`
	Items              []DueDiligenceItem `json:"items"`
}

type DistributionAmountResponse struct {
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
}

type DistributionBreakdownResponse struct {
	TotalTokens                string                     `json:"total_tokens"`
	TotalTokenHolders          int64                      `json:"total_token_holders"`
	NetIncome                  DistributionAmountResponse `json:"net_income"`
	DistributionAmountPerToken DistributionAmountResponse `json:"distribution_amount_per_token"`
	Calculation                string                     `json:"calculation"`
}

type DistributionPayoutRecordResponse struct {
	ID                    string    `json:"id"`
	BeneficiaryPublicKey  string    `json:"beneficiary_public_key"`
	ConfirmedTokenBalance string    `json:"confirmed_token_balance"`
	Amount                string    `json:"amount"`
	Currency              string    `json:"currency"`
	CannotReceiveAsset    bool      `json:"cannot_receive_asset"`
	Paid                  bool      `json:"paid"`
	CreatedAt             time.Time `json:"created_at"`
}

type DistributionPayoutResponse struct {
	DistributionID       string                             `json:"distribution_id"`
	TokenizedAssetID     string                             `json:"tokenized_asset_id"`
	Batch                string                             `json:"batch,omitempty"`
	Amount               string                             `json:"amount"`
	Currency             string                             `json:"currency"`
	AmountPerToken       string                             `json:"amount_per_token"`
	Status               string                             `json:"status"`
	PaymentScheduleReady bool                               `json:"payment_schedule_ready"`
	PayoutCompleted      bool                               `json:"payout_completed"`
	Payouts              []DistributionPayoutRecordResponse `json:"payouts"`
	CreatedAt            time.Time                          `json:"created_at,omitempty"`
	UpdatedAt            time.Time                          `json:"updated_at,omitempty"`
}

type DistributionDetailResponse struct {
	Distribution Distribution                  `json:"distribution"`
	Requester    FundReleaseRequesterResponse  `json:"requester"`
	Breakdown    DistributionBreakdownResponse `json:"breakdown"`
	Payout       DistributionPayoutResponse    `json:"payout"`
}

type DistributionListItemResponse struct {
	Distribution
	Requester FundReleaseRequesterResponse `json:"requester"`
}

type ReportTypeResponse struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

type DashboardResponse struct {
	Role             string                     `json:"role"`
	Summary          map[string]interface{}     `json:"summary"`
	RecentAssets     []TokenizedAssetResponse   `json:"recent_assets"`
	PortfolioSummary *CustodianPortfolioSummary `json:"portfolio_summary"`
	RecentActivity   []DashboardActivity        `json:"recent_activity"`
}

type CurrencyTotal struct {
	Currency string `json:"currency"`
	Amount   string `json:"amount"`
}

type CustodianPortfolioCategory struct {
	Category   string          `json:"category"`
	AssetCount int             `json:"asset_count"`
	Values     []CurrencyTotal `json:"values"`
}

type CustodianPortfolioSummary struct {
	AssetCount int                          `json:"asset_count"`
	Values     []CurrencyTotal              `json:"values"`
	Categories []CustodianPortfolioCategory `json:"categories"`
}

type DashboardActivity struct {
	ID                string    `json:"id"`
	Source            string    `json:"source"`
	Type              string    `json:"type"`
	Title             string    `json:"title"`
	Message           string    `json:"message,omitempty"`
	Status            string    `json:"status,omitempty"`
	ActorOrganization string    `json:"actor_organization,omitempty"`
	ActorName         string    `json:"actor_name,omitempty"`
	EntityType        string    `json:"entity_type,omitempty"`
	EntityID          string    `json:"entity_id,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
}

type FundReleaseRequesterResponse struct {
	MemberID         string `json:"member_id"`
	Name             string `json:"name,omitempty"`
	FirstName        string `json:"first_name,omitempty"`
	LastName         string `json:"last_name,omitempty"`
	Email            string `json:"email,omitempty"`
	OrganizationID   string `json:"organization_id"`
	OrganizationName string `json:"organization_name,omitempty"`
}

type FundReleaseListItemResponse struct {
	FundReleaseRequest
	Requester FundReleaseRequesterResponse  `json:"requester"`
	Reviewer  *FundReleaseRequesterResponse `json:"reviewer,omitempty"`
}

type FundReleaseDocumentResponse struct {
	ID               string `json:"id"`
	Title            string `json:"title"`
	Category         string `json:"category"`
	OriginalFilename string `json:"original_filename,omitempty"`
	MimeType         string `json:"mime_type,omitempty"`
	SizeBytes        int64  `json:"size_bytes,omitempty"`
	ExternalURL      string `json:"external_url,omitempty"`
	DownloadPath     string `json:"download_path,omitempty"`
}

type FundReleaseReceivingAccountResponse struct {
	Bank          string `json:"bank,omitempty"`
	AccountName   string `json:"account_name,omitempty"`
	AccountNumber string `json:"account_number,omitempty"`
}

type FundReleaseDetailResponse struct {
	FundReleaseRequest
	Requester           FundReleaseRequesterResponse        `json:"requester"`
	SupportingDocuments []FundReleaseDocumentResponse       `json:"supporting_documents"`
	ReceivingAccount    FundReleaseReceivingAccountResponse `json:"receiving_account"`
}

type FundReleaseAggregate struct {
	Count   int             `json:"count"`
	Amounts []CurrencyTotal `json:"amounts"`
}

type AssetManagerFundSummary struct {
	Requested FundReleaseAggregate `json:"requested"`
	Approved  FundReleaseAggregate `json:"approved"`
	Released  FundReleaseAggregate `json:"released"`
	Pending   FundReleaseAggregate `json:"pending"`
	Rejected  FundReleaseAggregate `json:"rejected"`
	Remaining []CurrencyTotal      `json:"remaining_balance"`
}

type TrusteeFundManagementSummary struct {
	AssetCount                  int             `json:"asset_count"`
	TotalAmountProcessed        []CurrencyTotal `json:"total_amount_processed"`
	TotalAmountFromPrimarySales []CurrencyTotal `json:"total_amount_from_primary_sales"`
	TotalIncomeFromAssets       []CurrencyTotal `json:"total_income_from_assets"`
	TotalAmountPaidToIssuers    []CurrencyTotal `json:"total_amount_paid_to_issuers"`
	TotalPaidToInvestors        []CurrencyTotal `json:"total_paid_to_investors"`
	MilestonePaymentBalance     []CurrencyTotal `json:"milestone_payment_balance"`
	TotalFeesGenerated          []CurrencyTotal `json:"total_fees_generated"`
}
