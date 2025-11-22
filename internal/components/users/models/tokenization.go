package users

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
	"trovo-wallet-api/internal/dynamiclinks"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/network"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/protocols/horizon"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gorm.io/gorm/clause"
)

/*
*
*****DocumentType and codes****
ProofOfAssetExistence
ProofOfAssetOwnership
ProofOfAssetStatusVerification
AssetCustodianAgreement
ProofOfAssetManager
AssetProtectionDocument
AssetValuationCertificate
AssetOwnerGovernmentID
ProofOfAssetCondtion
ThirdPartyTokenizationAgreement
ThirdPartyAssetOwnerBusinessRegistration
ThirdPartyAssetOwnerProofOfAddress
SEC Registration/Tokenization Approval
ComplianceWithLocalRegulation
ComplianceWithEnvironmentalStandard
EnvironmentalImpactAssessmentReport
ProofOfLegalFinancialCounsel
LegalFinancialAdvisorsContract
ProofOfExistingMortgagesOrLiens&Asset
ProofOfOutstandingLoansOnAsset
ProofOfLegalDisputeOrEncumbrancesOnAsset
ProofOfAssetAddress
TitleDeedsOrCertificates
OwnershipAgreements
EngineeringReportsForConstruction

*
*/

type JsonForm struct {
	ID         uint64
	JsonString string
}
type ExistingAssetValidationAssetDocument struct {
	ID                                uint64 `gorm:"" json:"-" form:"-"`
	ProofOfAssetExistence             int    `gorm:"default:1" json:"proofOfAssetExistence"`             //1
	ProofOfAssetAddress               int    `gorm:"default:1" json:"proofOfAssetAddress"`               //22
	ProofOfAssetOwnership             int    `gorm:"default:1" json:"proofOfAssetOwnership"`             //2
	TitleDeedsOrCertificates          int    `gorm:"default:1" json:"titleDeedsOrCertificates"`          //23
	OwnershipAgreements               int    `gorm:"default:1" json:"ownershipAgreements"`               //24
	ProofOfAssetStatusVerification    int    `gorm:"default:1" json:"proofOfAssetStatusVerification"`    //3
	ProofOfAssetCondtion              int    `gorm:"default:1" json:"proofOfAssetCondtion"`              //9
	AssetOwnerGovernmentID            int    `gorm:"default:1" json:"assetOwnerGovernmentId"`            //8
	EngineeringReportsForConstruction int    `gorm:"default:1" json:"engineeringReportsForConstruction"` //25

}

type TokenizedAsset struct {
	ID                                           string                          `json:"id"`
	CreatedAt                                    time.Time                       `json:"createdAt"`
	UpdatedAt                                    time.Time                       `json:"updatedAt"`
	InitiatorUsername                            string                          `gorm:"size:50;not null" json:"initiatorUsername"`
	AssetSector                                  *string                         `json:"assetSector"`
	AssetSubSector                               *string                         `json:"assetSubSector"`
	AssetType                                    *string                         `json:"assetType"`
	AssetName                                    *string                         `json:"assetName"`
	ApprovedAssetCustodianID                     uint64                          `gorm:"not null;default:0" json:"approvedAssetCustodianId"`
	ApprovedAssetCustodian                       ApprovedAssetCustodian          `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"approvedAssetCustodianInfo"`
	OfferingType                                 *string                         `gorm:"default:'PRIVATE'" json:"offeringType"` //PRIVATE, PUBLIC
	ClosedGroupID                                *string                         `gorm:"null" json:"closedGroupId"`
	ClosedGroup                                  ClosedGroup                     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"closedGroupInfo"`
	SecApproval                                  int                             `gorm:"default:0" json:"secApproval"`
	SecApprovalIdNumber                          *string                         `json:"secApprovalIdNumber"`
	IssuingWalletPublicKey                       *string                         `gorm:"size:60" json:"issuingWalletPublicKey"`
	IssuingWalletAlias                           *string                         `gorm:"size:60" json:"issuingWalletAlias"`
	MarketMakingWallet                           *string                         `json:"marketMakingWallet"`
	AssetDescription                             *string                         `json:"assetDescription"`
	AssetCountryLocation                         *string                         `gorm:"not null;size:2;default'NG'" json:"assetCountryLocation"`
	AssetPhysicalAddress                         *string                         `json:"assetPhysicalAddress"`
	AssetLongitude                               *string                         `json:"assetLongitude"`
	AssetLatitude                                *string                         `json:"assetLatitude"`
	OwnershipType                                *string                         `json:"ownershipType"` //DIRECT,THIRD-PARTY
	OwnershipKind                                *string                         `json:"ownershipKind"` //INDIVIDUAL,CORPORATE,
	InitialOwnerPreferredWalletAddress           *string                         `json:"initialOwnerPreferredWalletAddress"`
	AssetOwnerName                               *string                         `json:"assetOwnerName"`
	AssetOwnerAddress                            *string                         `json:"assetOwnerAddress"`
	AssetManagerID                               uint64                          `gorm:"default:0" json:"assetManagerId"`
	AssetManager                                 AssetManager                    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"assetManagerInfo"`
	AssetIssuingHouseID                          uint64                          `gorm:"not null;default:0" json:"assetIssuingHouseId"`
	AssetIssuingHouse                            AssetIssuingHouse               `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"assetIssuingHouseInfo"`
	LegalAndProfesionalPartnerID                 uint64                          `gorm:"not null;default:0" json:"legalAndProfesionalPartnerId"`
	LegalAndProfesionalPartner                   LegalAndProfesionalPartner      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"legalAndProfesionalPartnerInfo"`
	RatingAgencyID                               uint64                          `gorm:"not null;default:0" json:"ratingAgencyId"`
	RatingAgency                                 RatingAgency                    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"ratingAgencyInfo"`
	TrusteeID                                    uint64                          `gorm:"not null;default:0" json:"trusteeId"`
	Trustee                                      Trustee                         `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"trusteeInfo"`
	AssetQuoteCurrency                           *string                         `gorm:"default:'CNGN'" json:"assetQuoteCurrency"`
	AssetCurrentValue                            float64                         `gorm:"default:0" json:"assetCurrentValue"`
	AssetOwnerRetainedOrContributedValue         float64                         `gorm:"default:0" json:"assetOwnerRetainedOrContributedValue"`
	AssetMscCostOutisdeOfValuation               float64                         `gorm:"default:0" json:"assetMscCostOutisdeOfValuation"`
	ValueOfTokenizedAsset                        float64                         `gorm:"default:0" json:"valueOfTokenizedAsset"`
	ProtectionMethods                            *string                         `json:"protectionMethods"` //csv format
	InsuranceCompanyName                         *string                         `json:"insuranceCompanyName"`
	InsurancePolicyNumber                        *string                         `json:"insurancePolicyNumber"`
	InsurancePolicyHolder                        *string                         `json:"insurancePolicyHolder"`
	PercentageValueOfInsurance                   float64                         `gorm:"default:0" json:"percentageValueOfInsurance"`
	IsFreeFromLiensAndEncumbrances               int                             `gorm:"default:0" json:"IsFreeFromLiensAndEncumbrances"`
	AssetAlreadyExists                           int                             `gorm:"default:0" json:"assetAlreadyExists"`
	VettingStatus                                int                             `gorm:"default:0" json:"vettingStatus"`
	DueDiligenceFail                             int                             `gorm:"default:0" json:"dueDiligenceFail"` //0=False(success/in-progress), 1= true (failed).
	DueDiligenceFailureReason                    *string                         `json:"dueDiligenceFailureReason"`
	AssetTokenizationDocuments                   []AssetTokenizationDocument     `json:"AssetTokenizationDocuments"`
	ProofOfPaymentDocuments                      []TokenizationFeeProofOfPayment `json:"ProofOfPaymentDocuments"`
	AssetCode                                    *string                         `gorm:"size:12; index:idx_unique_tokenized_asset_code,unique" json:"assetCode"`
	AssetLogo                                    *string                         `json:"assetLogo"`
	AssetWebsite                                 *string                         `gorm:"default:'trovotech.io'" json:"assetWebsite"`
	InitialValueOfTokenizedAsset                 float64                         `gorm:"default:0" json:"initialValueOfTokenizedAsset"`
	NumberOfTokenToBeIssued                      float64                         `gorm:"default:0" json:"numberOfTokenToBeIssued"`
	MaxNumberOfTokenAvailableForSale             float64                         `gorm:"default:0" json:"maxNumberOfTokenAvailableForSale"`
	FeeInAsset                                   float64                         `gorm:"default:0" json:"feeInAsset"`
	FeeInAssetPercent                            float64                         `gorm:"default:0" json:"feeInAssetPercent"`
	FeeInFiat                                    float64                         `gorm:"default:0" json:"feeInFiat"` //tokenization fee in fiat
	NumberOfTokenToBeSold                        float64                         `gorm:"default:0" json:"numberOfTokenToBeSold"`
	TotalTokenHeldByManager                      float64                         `gorm:"default:0" json:"totalTokenHeldByManager"`
	WalletToHoldAssetsNotForSale                 *string                         `json:"walletToHoldAssetsNotForSale"`
	PricePerToken                                float64                         `gorm:"default:0" json:"pricePerToken"`
	SalesStart                                   time.Time                       `json:"salesStart"`
	SalesEnd                                     time.Time                       `json:"salesEnd"`
	CapOnPurchase                                int                             `gorm:"default:0" json:"capOnPurchase"`
	CapQuantity                                  float64                         `gorm:"default:0" json:"capQuantity"`
	CapAmountInFiat                              float64                         `gorm:"default:0" json:"capAmountInFiat"`
	CapDurationInDays                            int                             `gorm:"default:0" json:"capDurationInDays"`
	ProceedCycle                                 *string                         `gorm:"size:50" json:"proceedCycle"`
	TokenizationFeeID                            *uint64                         `gorm:"default:0" json:"tokenizationFeeId"`
	TokenizationFee                              TokenizationFee                 `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"tokenizationFee"`
	ExcludeSecFee                                int                             `gorm:"default:0" json:"excludeSecFee"`
	SECTokenizationFeePercent                    float64                         `gorm:"default:0" json:"SECTokenizationFeePercent"`
	SECTokenizationFeeFixed                      float64                         `gorm:"default:0" json:"SECTokenizationFeeFixed"`
	SECTokenizationFeeValue                      float64                         `gorm:"default:0" json:"SECTokenizationFeeValue"`
	CustodianFeePercent                          float64                         `gorm:"default:0" json:"custodianFeePercent"`
	CustodianFeeFixed                            float64                         `gorm:"default:0" json:"custodianFeeFixed"`
	CustodianFeeValue                            float64                         `gorm:"default:0" json:"custodianFeeValue"`
	AssetManagerFeeValue                         float64                         `gorm:"default:0" json:"assetManagerFeeValue"`
	AssetManagerFeePercent                       float64                         `gorm:"default:0" json:"assetManagerFeePercent"`
	AssetManagerFeeFixed                         float64                         `gorm:"default:0" json:"assetManagerFeeFixed"`
	IssuingHouseFeeValue                         float64                         `gorm:"default:0" json:"issuingHouseFeeValue"`
	IssuingHouseFeePercent                       float64                         `gorm:"default:0" json:"issuingHouseFeePercent"`
	IssuingHouseFeeFixed                         float64                         `gorm:"default:0" json:"issuingHouseFeeFixed"`
	LegalAndProfessionalFeePercent               float64                         `gorm:"default:0" json:"legalAndProfessionalFeePercent"`
	LegalAndProfessionalFeeFixed                 float64                         `gorm:"default:0" json:"legalAndProfessionalFeeFixed"`
	LegalAndProfessionalFeeValue                 float64                         `gorm:"default:0" json:"legalAndProfessionalFeeValue"`
	RatingAgencyFeePercent                       float64                         `gorm:"default:0" json:"ratingAgencyFeePercent"`
	RatingAgencyFeeFixed                         float64                         `gorm:"default:0" json:"ratingAgencyFeeFixed"`
	RatingAgencyFeeValue                         float64                         `gorm:"default:0" json:"ratingAgencyFeeValue"`
	TrusteeFeePercent                            float64                         `gorm:"default:0" json:"trusteeFeePercent"`
	TrusteeFeeFixed                              float64                         `gorm:"default:0" json:"trusteeFeeFixed"`
	TrusteeFeeValue                              float64                         `gorm:"default:0" json:"trusteeFeeValue"`
	VATPercent                                   float64                         `gorm:"default:0" json:"vatPercent"`
	VATValue                                     float64                         `gorm:"default:0" json:"vatValue"`
	VATInAsset                                   float64                         `gorm:"default:0" json:"vatInAsset"`
	ProceedPayoutCurrency                        *string                         `json:"proceedPayoutCurrency"`
	ProceedPayoutType                            int                             `gorm:"default:0" json:"proceedPayoutType"` // FIAT=1, CRYPTO=0
	ExemptedCountries                            *string                         `json:"exemptedCountries"`
	HasAdditionalKYCRequirements                 int                             `gorm:"default:0" json:"hasAdditionalKYCRequirements"`
	AdditionalKYCRequirements                    *string                         `json:"additionalKYCRequirements"`
	InvestorAccreditationRequired                int                             `gorm:"default:0" json:"investorAccreditationRequired"`
	AssetTokenizationStatus                      int                             `gorm:"default:0" json:"assetTokenizationStatus"`
	LastUpdatedBy                                *string                         `gorm:"null" json:"lastUpdatedBy"`
	TokenizationTransaction                      *string                         `gorm:"null" json:"-"`
	AgreeTransferTitleToCustodian                int                             `gorm:"default:0" json:"agreeTransferTitleToCustodian"`
	ContractualProtectionRevGuarantees           int                             `gorm:"default:0" json:"contractualProtectionRevGuarantees"`
	ContractualProtectionPerfBond                int                             `gorm:"default:0" json:"contractualProtectionPerfBond"`
	ContractualProtectionSLA                     int                             `gorm:"default:0" json:"contractualProtectionSLA"`
	RiskSharingMechanismPPPs                     int                             `gorm:"default:0" json:"riskSharingMechanismPPPs"`
	RiskSharingMechanismHedgeInstruments         int                             `gorm:"default:0" json:"riskSharingMechanismHedgeInstruments"`
	RiskSharingMechanismCompletionGuarantees     int                             `gorm:"default:0" json:"riskSharingMechanismCompletionGuarantees"`
	IndependentMonitoringList                    *string                         `gorm:"null" json:"independentMonitoringList"`
	ESGSafeguardsSusCerts                        int                             `gorm:"default:0" json:"eSGSafeguardsSusCerts"`
	ESGSafeguardsCommEngPlans                    int                             `gorm:"default:0" json:"eSGSafeguardsCommEngPlans"`
	SecurityMeasuresAccessControl                int                             `gorm:"default:0" json:"securityMeasuresAccessControl"`
	SecurityMeasuresSurveilanceSystems           int                             `gorm:"default:0" json:"securityMeasuresSurveilanceSystems"`
	SecurityMeasuresOnSiteSecurityPersonnel      int                             `gorm:"default:0" json:"securityMeasuresOnSiteSecurityPersonnel"`
	SecurityMeasuresPerimeterSecurity            int                             `gorm:"default:0" json:"securityMeasuresPerimeterSecurity"`
	SecurityMeasuresCriticalInfraProtections     int                             `gorm:"default:0" json:"securityMeasuresCriticalInfraProtections"`
	OtherAssetProtection                         *string                         `gorm:"null" json:"otherAssetProtection"`
	LegalAdvisor                                 *string                         `gorm:"null" json:"legalAdvisor"`
	FinancialAdvisor                             *string                         `gorm:"null" json:"financialAdvisor"`
	UndertakingNoLien                            int                             `gorm:"default:0" json:"undertakingNoLien"`
	UndertakingNotCollateral                     int                             `gorm:"default:0" json:"undertakingNotCollateral"`
	UndertakingNoClaims                          int                             `gorm:"default:0" json:"undertakingNoClaims"`
	UndertakingNoForeclosure                     int                             `gorm:"default:0" json:"undertakingNoForeclosure"`
	ComplianceNoViolation                        int                             `gorm:"default:0" json:"complianceNoViolation"`
	ComplianceAllPermits                         int                             `gorm:"default:0" json:"complianceAllPermits"`
	OutstandingFinancialRespNoDebts              int                             `gorm:"default:0" json:"outstandingFinancialRespNoDebts"`
	OutstandingFinancialRespNoHiddenLiabilities  int                             `gorm:"default:0" json:"outstandingFinancialRespNoHiddenLiabilities"`
	RiskManagementFullyInsured                   int                             `gorm:"default:0" json:"riskManagementFullyInsured"`
	RiskManagementDeclaredValue                  int                             `gorm:"default:0" json:"riskManagementDeclaredValue"`
	PhysicalConditionSound                       int                             `gorm:"default:0" json:"physicalConditionSound"`
	PhysicalConditionNolease                     int                             `gorm:"default:0" json:"physicalConditionNolease"`
	PhysicalConditionNoUndisclosedEasements      int                             `gorm:"default:0" json:"physicalConditionNoUndisclosedEasements"`
	MintingInitators                             *string                         `gorm:"null" json:"mintingInitators"` //CSV of approvers
	MintingApprovers                             *string                         `gorm:"null" json:"mintingApprovers"` //csv of initators
	BankID                                       *uint64                         `gorm:"null" json:"bankId"`
	Bank                                         Bank                            `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"bankInfo"`
	AccountNumber                                *string                         `gorm:"null" json:"accountNumber"`
	BeneficiaryName                              *string                         `gorm:"null" json:"beneficiaryName"`
	TokenizationApplicationFee                   float64                         `gorm:"default:0" json:"tokenizationApplicationFee"`
	TokenizationApplicationFeeAsset              string                          `gorm:"not null;default:'TROV:GAXMBPVA2GNG6A3NV6Q664VZASMROS5ZACKSMTPVCRIKPOJIV43A2CTJ'" json:"tokenizationApplicationFeeAsset"`
	ProjectStrategicObjectives                   *string                         `json:"projectStrategicObjectives"`
	ProjectDevelopmentTimeline                   *string                         `json:"projectDevelopmentTimeline"`
	ProjectKeyMilestoneAndDates                  *string                         `json:"projectKeyMilestoneAndDates"`
	ProjectScope                                 *string                         `json:"projectScope"`
	ProjectEconomicBenefits                      *string                         `json:"projectEconomicBenefits"`
	ProjectExpectedNoOfJobs                      int                             `json:"projectExpectedNoOfJobs"`
	ProjectIntendedSocialBenefits                *string                         `json:"projectIntendedSocialBenefits"`
	ProjectTechnicalPartners                     *string                         `json:"projectTechnicalPartners"`
	ProjectFinancialPartners                     *string                         `json:"projectFinancialPartners"`
	EstimatedProjectIRR                          float64                         `json:"estimatedProjectIRR"`
	EstimatedProjectROI                          float64                         `json:"estimatedProjectROI"`
	EstimatedProjectNPV                          float64                         `json:"estimatedProjectNPV"`
	EstimatedProjectPaybackPeriodsInMonths       int                             `json:"estimatedProjectPaybackPeriodsInMonths"`
	KeyAssumptionsList                           *string                         `json:"keyAssumptionsList"`
	ProjectIdentifiedLegalRisks                  *string                         `json:"projectIdentifiedLegalRisks"`
	ProjectIdentifiedRegulatoryRisks             *string                         `json:"projectIdentifiedRegulatoryRisks"`
	ProjectIdentifiedOperationalOrExecutionRisks *string                         `json:"projectIdentifiedOperationalOrExecutionRisks"`
	ProjectIdentifiedMarketRisks                 *string                         `json:"projectIdentifiedMarketRisks"`
	ProjectIdentifiedOtherRelevantRisks          *string                         `json:"projectIdentifiedOtherRelevantRisks"`
	DeepLink                                     *string                         `gorm:"null" json:"deepLink"`
	IsinOrSerialNumber                           *string                         `json:"isinOrSerialNumber"` /// begin the mutual fund fields
	InstrumentName                               *string                         `json:"instrumentName"`
	InstrumentType                               *string                         `json:"instrumentType"`
	TotalIssueSize                               float64                         `gorm:"default:0" json:"totalIssueSize"`
	IssueDate                                    time.Time                       `json:"issueDate"`
	MaturityDate                                 time.Time                       `json:"maturityDate"`
	FaceValuePerUnit                             float64                         `gorm:"default:0" json:"faceValuePerUnit"`
	CouponOrInterestRateType                     *string                         `json:"couponOrInterestRateType"`
	CouponOrInterestRate                         float64                         `json:"couponOrInterestRate"`
	ReferenceIndex                               *string                         `json:"referenceIndex"`
	ExpectedYield                                float64                         `gorm:"default:0" json:"expectedYield"`
	EarlyRedemptionOptionInvestor                *string                         `json:"earlyRedemptionOptionInvestor"`
	MinimumInvestmentAmount                      float64                         `json:"minimumInvestmentAmount"`
	TaxTreatmentTokenHolders                     *string                         `json:"taxTreatmentTokenHolders"`
	PaymentStructureToTokenHolders               *string                         `json:"paymentStructureToTokenHolders"`
	RedemptionMethod                             *string                         `json:"redemptionMethod"`
	PaymentStructure                             *string                         `json:"paymentStructure"`
	RepaymentMethod                              *string                         `json:"repaymentMethod"`
	PaymentCycle                                 *string                         `json:"paymentCycle"`
	WeightedAverageLife                          float64                         `json:"weightedAverageLife"`
	UnderlyingAssetPoolSize                      float64                         `json:"underlyingAssetPoolSize"`
	PoolComposition                              *string                         `json:"poolComposition"`
	CreditEnhancementMethod                      *string                         `json:"creditEnhancementMethod"`
	SummaryOfUseOfProceeds                       *string                         `json:"summaryOfUseOfProceeds"`
	CreditRatingIfAny                            *string                         `json:"creditRatingIfAny"`
	IssuerName                                   *string                         `json:"issuerName"`
	IssuerType                                   *string                         `json:"issuerType"`
	IssuerContactPerson                          *string                         `json:"issuerContactPerson"`
	ContactEmail                                 *string                         `json:"contactEmail"`
	ContactPhoneNumber                           *string                         `json:"contactPhoneNumber"`
	BriefCompanyOverview                         *string                         `json:"briefCompanyOverview"`
	MortgageOriginators                          *string                         `json:"mortgageOriginators"`
	Servicer                                     *string                         `json:"servicer"`
	InstrumentTrustee                            *string                         `json:"instrumentTrustee"`
	InstrumentCustodians                         *string                         `json:"instrumentCustodian"`
	InstrumentLegalAdvisor                       *string                         `json:"instrumentLegalAdvisor"`
	SpecialPurposeVehicle                        *string                         `json:"specialPurposeVehicle"`
	InstrumentAssetManagerOrAdministrator        *string                         `json:"instrumentAssetManagerOrAdministrator"`
	UnderwriterIfAny                             *string                         `json:"underwriterIfAny"`
	InstrumentCreditRatingAgency                 *string                         `json:"instrumentCreditRatingAgency"`
	AuditorOrVerifier                            *string                         `json:"auditorOrVerifier"`
	CreditRiskAssessment                         *string                         `json:"creditRiskAssessment"`
	CreditRating                                 *string                         `json:"creditRating"`
	PrepaymentRisk                               *string                         `json:"prepaymentRisk"`
	InterestRateRisk                             *string                         `json:"interestRateRisk"`
	StructuralComplexityRisk                     *string                         `json:"structuralComplexityRisk"`
	LegalOrRegulatoryRisk                        *string                         `json:"legalOrRegulatoryRisk"`
	OperationalRisk                              *string                         `json:"operationalRisk"`
	MarketRisk                                   *string                         `json:"marketRisk"`
	EsgRisk                                      *string                         `json:"esgRisk"`
	MitigationMeasures                           *string                         `json:"mitigationMeasures"`
	IssuingAuthority                             *string                         `json:"issuingAuthority"`
	RegulatoryApprovalId                         *string                         `json:"regulatoryApprovalId"`
	LicenseApprovalReferenceNumber               *string                         `json:"licenseApprovalReferenceNumber"`
	ListingStatus                                *string                         `json:"listingStatus"`
	CurrencyOfIssuance                           *string                         `json:"currencyOfIssuance"`
	CouponRateType                               *string                         `json:"couponRateType"`
	CouponRate                                   float64                         `gorm:"default:0" json:"couponRate"`
	SpreadOrMargin                               float64                         `gorm:"default:0" json:"spreadOrMargin"`
	ResetFrequency                               *string                         `json:"resetFrequency"`
	CouponPaymentFrequency                       *string                         `json:"couponPaymentFrequency"`
	RedemptionStructure                          *string                         `json:"redemptionStructure"`
	EarlyRedemptionOption                        *string                         `json:"earlyRedemptionOption"`
	EarlyRedemptionPenalty                       *string                         `json:"earlyRedemptionPenalty"`
	TaxTreatment                                 *string                         `json:"taxTreatment"`
	NavOrMarketValueUpdates                      *string                         `json:"navOrMarketValueUpdates"`
	ImpactMetrics                                *string                         `json:"impactMetrics"`
	LegalBacking                                 *string                         `json:"legalBacking"`
	DefaultHistory                               *string                         `json:"defaultHistory"`
	RiskFactorsSummary                           *string                         `json:"riskFactorsSummary"`
	PayingAgent                                  *string                         `json:"payingAgent"`
	Auditor                                      *string                         `json:"auditor"`
	RegistrarOrCSCSAgent                         *string                         `json:"registrarOrCSCSAgent"`
	FundStructure                                *string                         `json:"fundStructure"`
	AssetManagementCompanyName                   *string                         `json:"assetManagementCompanyName"`
	FundManagers                                 *string                         `json:"fundManagers"`
	RegulatoryLicenseNumber                      *string                         `json:"regulatoryLicenseNumber"`
	FundLaunchDate                               *time.Time                      `json:"fundLaunchDate"`
	TotalExpenseRatio                            float64                         `gorm:"default:0" json:"totalExpenseRatio"`
	ExitLoadRedemptionFee                        float64                         `json:"exitLoadRedemptionFee"`
	Tenure                                       int                             `gorm:"default:0" json:"tenure"`
	InitialNetAssetValue                         float64                         `json:"initialNetAssetValue"`
	NavUpdateFrequency                           *string                         `json:"navUpdateFrequency"`
	NavCalculationMethod                         *string                         `json:"navCalculationMethod"`
	RedemptionRules                              *string                         `json:"redemptionRules"`
	LockInPeriod                                 int                             `gorm:"default:0" json:"lockInPeriod"`
	EntryLoad                                    float64                         `gorm:"default:0" json:"entryLoad"`
	PerformanceFee                               float64                         `gorm:"default:0" json:"performanceFee"`
	DividendPolicy                               *string                         `json:"dividendPolicy"`
	LiquidityProfile                             *string                         `json:"liquidityProfile"`
	DistributionFrequency                        *string                         `json:"distributionFrequency"`
	DistributionMethod                           *string                         `json:"distributionMethod"`
	BenchmarkComparisonMethod                    *string                         `json:"benchmarkComparisonMethod"`
	FeeBreakdownSummary                          *string                         `json:"feeBreakdownSummary"`
	InvestmentObjective                          *string                         `json:"investmentObjective"`
	EquityStrategy                               *string                         `json:"equityStrategy"`
	MarketCapitalizationFocus                    *string                         `json:"marketCapitalizationFocus"`
	BenchmarkIndex                               *string                         `json:"benchmarkIndex"`
	SectorExposureLimits                         *string                         `json:"sectorExposureLimits"`
	TopHoldings                                  *string                         `json:"topHoldings"`
	GeographicExposure                           *string                         `json:"geographicExposure"`
	RiskProfile                                  *string                         `json:"riskProfile"`
	VolatilityEstimate                           float64                         `gorm:"default:0" json:"volatilityEstimate"`
	DividendYield                                float64                         `gorm:"default:0" json:"dividendYield"`
	TrusteeName                                  *string                         `json:"trusteeName"`
	FundAdministrator                            *string                         `json:"fundAdministrator"`
	InvestmentCommitteeMembers                   *string                         `json:"investmentCommitteeMembers"`
	IsinOrSecFundCode                            *string                         `json:"isinOrSecFundCode"`
	ExitLoadOrRedemptionFee                      float64                         `json:"exitLoadOrRedemptionFee"`
	FundRiskRating                               *string                         `json:"fundRiskRating"`
	AssetAllocation                              *string                         `json:"assetAllocation"`
	AverageMaturity                              float64                         `gorm:"default:0" json:"averageMaturity"`
	YieldToMaturity                              float64                         `gorm:"default:0" json:"yieldToMaturity"`
	CreditRatingProfile                          *string                         `json:"creditRatingProfile"`
	LockInPeriodPortfolio                        int                             `gorm:"default:0" json:"lockInPeriodPortfolio"`
	PerformanceFeeIfAny                          float64                         `json:"performanceFeeIfAny"`
	TopEquityHoldings                            *string                         `json:"topEquityHoldings"`
	TargetAllocation                             *string                         `json:"targetAllocation"`
	AllowedAllocationRange                       *string                         `json:"allowedAllocationRange"`
	AssetClassesIncluded                         *string                         `json:"assetClassesIncluded"`
	RebalancingFrequency                         *string                         `json:"rebalancingFrequency"`
	BenchmarkIndexComposite                      *string                         `json:"benchmarkIndexComposite"`
	TopEquityHoldingsList                        *string                         `json:"topEquityHoldingsList"`
	TopDebtHoldings                              *string                         `json:"topDebtHoldings"`
	CreditRatingDistribution                     *string                         `json:"creditRatingDistribution"`
	AverageMaturityHybrid                        float64                         `json:"averageMaturityHybrid"`
	YieldToMaturityHybrid                        float64                         `json:"yieldToMaturityHybrid"`
	TitleOfIssuance                              *string                         `json:"titleOfIssuance"`
	TypeOfCommercialPaper                        *string                         `json:"typeOfCommercialPaper"`
	PricingYield                                 float64                         `gorm:"default:0" json:"pricingYield"`
	UseOfProceeds                                *string                         `json:"useOfProceeds"`
	Ranking                                      *string                         `json:"ranking"`
	BackingSecurity                              *string                         `json:"backingSecurity"`
	IssuerRegistrationNumber                     *string                         `json:"issuerRegistrationNumber"`
	IncorporationDate                            *time.Time                      `json:"incorporationDate"`
	RcNumber                                     *string                         `json:"rcNumber"`
	TaxIdNumber                                  *string                         `json:"taxIdNumber"`
	OfficeAddress                                *string                         `json:"officeAddress"`
	Rating                                       *string                         `json:"rating"`
	PartiesInvolvedIssuer                        *string                         `json:"partiesInvolvedIssuer"`
	PartiesInvolvedArranger                      *string                         `json:"partiesInvolvedArranger"`
	PartiesInvolvedLegalAdviser                  *string                         `json:"partiesInvolvedLegalAdviser"`
	PartiesInvolvedAuditor                       *string                         `json:"partiesInvolvedAuditor"`
	PartiesInvolvedRatingAgency                  *string                         `json:"partiesInvolvedRatingAgency"`
	PartiesInvolvedCustodian                     *string                         `json:"partiesInvolvedCustodian"`
	PartiesInvolvedTrustee                       *string                         `json:"partiesInvolvedTrustee"`
	PartiesInvolvedAuditorVerifier               *string                         `json:"partiesInvolvedAuditorVerifier"`
	SecurityRiskLegalBacking                     *string                         `json:"securityRiskLegalBacking"`
	SecurityRiskCollateral                       *string                         `json:"securityRiskCollateral"`
	SecurityRiskDefaultHistory                   *string                         `json:"securityRiskDefaultHistory"`
	SecurityRiskCreditRating                     *string                         `json:"securityRiskCreditRating"`
	SecurityRiskRiskFactorsSummary               *string                         `json:"securityRiskRiskFactorsSummary"`
	SecurityRiskBusinessRisk                     *string                         `json:"securityRiskBusinessRisk"`
	SecurityRiskDefaultRisk                      *string                         `json:"securityRiskDefaultRisk"`
	SecurityRiskLiquidityRisk                    *string                         `json:"securityRiskLiquidityRisk"`
	SecurityRiskRegulatoryRisk                   *string                         `json:"securityRiskRegulatoryRisk"`
	SecurityRiskMarketRisk                       *string                         `json:"securityRiskMarketRisk"`
	SecurityRiskOperationalRisk                  *string                         `json:"securityRiskOperationalRisk"`
	SecurityRiskMitigationMeasures               *string                         `json:"securityRiskMitigationMeasures"`
	CommodityType                                *string                         `json:"commodityType"`
	CommodityDescription                         *string                         `json:"commodityDescription"`
	Quantity                                     int                             `gorm:"default:0" json:"quantity"`
	QualityGrade                                 *string                         `json:"qualityGrade"`
	IssuerContactInfo                            *string                         `json:"issuerContactInfo"`
	WarehouseName                                *string                         `json:"warehouseName"`
	WarehouseOperatorName                        *string                         `json:"warehouseOperatorName"`
	WarehouseLicenseNumber                       *string                         `json:"warehouseLicenseNumber"`
	WarehouseLocation                            *string                         `json:"warehouseLocation"`
	WrNumber                                     *string                         `json:"wrNumber"`
	WrIssueDate                                  *time.Time                      `json:"wrIssueDate"`
	WrExpiryDate                                 *time.Time                      `json:"wrExpiryDate"`
	WrSystemRegistration                         *string                         `json:"wrSystemRegistration"`
	WrRegistrationNumber                         *string                         `json:"wrRegistrationNumber"`
	WrVerifier                                   *string                         `json:"wrVerifier"`
	StorageCondition                             *string                         `json:"storageCondition"`
	WarehouseAccreditationBody                   *string                         `json:"warehouseAccreditationBody"`
	MinimumPurchaseAmount                        float64                         `json:"minimumPurchaseAmount"`
	AutoRollover                                 *string                         `json:"autoRollover"`
	CurrentBeneficialOwner                       *string                         `json:"currentBeneficialOwner"`
	WrCustodianName                              *string                         `json:"wrCustodianName"`
	OwnershipRightsRepresented                   *string                         `json:"ownershipRightsRepresented"`
	TrusteeOrThirdPartyOversight                 *string                         `json:"trusteeOrThirdPartyOversight"`
	LienOrEncumbrances                           *string                         `json:"lienOrEncumbrances"`
	AssetValuation                               float64                         `gorm:"default:0" json:"assetValuation"`
	ValuationDate                                *time.Time                      `json:"valuationDate"`
	ValuationMethodology                         *string                         `json:"valuationMethodology"`
	TokenizationObjective                        *string                         `json:"tokenizationObjective"`
	HoldingPeriod                                int                             `gorm:"default:0" json:"holdingPeriod"`
	RedemptionMechanism                          *string                         `json:"redemptionMechanism"`
	PartiesInvolvedUnderwriter                   *string                         `json:"partiesInvolvedUnderwriter"`
	PartiesInvolvedAssetManager                  *string                         `json:"partiesInvolvedAssetManager"`
	PartiesInvolvedLegalAdvisor                  *string                         `json:"partiesInvolvedLegalAdvisor"`
	PartiesInvolvedRegulator                     *string                         `json:"partiesInvolvedRegulator"`
	RisksMarketRisk                              *string                         `json:"risksMarketRisk"`
	RisksStorageRisk                             *string                         `json:"risksStorageRisk"`
	RisksTitleRisk                               *string                         `json:"risksTitleRisk"`
	RisksFraudRisk                               *string                         `json:"risksFraudRisk"`
	RisksInsuranceRisk                           *string                         `json:"risksInsuranceRisk"`
	RisksOperationalRisk                         *string                         `json:"risksOperationalRisk"`
	RisksRegulatoryRisk                          *string                         `json:"risksRegulatoryRisk"`
	RisksLiquidityRisk                           *string                         `json:"risksLiquidityRisk"`
	RisksForceMajeureRisk                        *string                         `json:"risksForceMajeureRisk"`
	RisksEarlyRedemptionRisk                     *string                         `json:"risksEarlyRedemptionRisk"`
	RisksMitigationMeasures                      *string                         `json:"risksMitigationMeasures"`
	RisksInsuranceCoverageSummary                *string                         `json:"risksInsuranceCoverageSummary"`
	RisksInsuranceProvider                       *string                         `json:"risksInsuranceProvider"`
	RisksCoverageValue                           float64                         `json:"risksCoverageValue"`
	QuanlityStandard                             *string                         `json:"quanlityStandard"`
	IssuerContactInformation                     *string                         `json:"issuerContactInformation"`
	VaultCustodianName                           *string                         `json:"vaultCustodianName"`
	VaultOperator                                *string                         `json:"vaultOperator"`
	VaultLicenseNumber                           *string                         `json:"vaultLicenseNumber"`
	VaultLocation                                *string                         `json:"vaultLocation"`
	Number                                       *string                         `json:"number"`
	IssuerDate                                   *time.Time                      `json:"issuerDate"`
	ExpiryDate                                   *time.Time                      `json:"expiryDate"`
	RegistryRecord                               *string                         `json:"registryRecord"`
	Verifier                                     *string                         `json:"verifier"`
	StorageConditions                            *string                         `json:"storageConditions"`
	VaultAccreditationBody                       *string                         `json:"vaultAccreditationBody"`
	OwnershipLegalHolder                         *string                         `json:"ownershipLegalHolder"`
	OwnershipCustodianName                       *string                         `json:"ownershipCustodianName"`
	OwnershipTrustee                             *string                         `json:"ownershipTrustee"`
	OwnershipLienOrEncumbrances                  *string                         `json:"ownershipLienOrEncumbrances"`
	ValuationAssetValuation                      *string                         `json:"valuationAssetValuation"`
	HoldingLockinPeriod                          int                             `gorm:"default:0" json:"holdingLockinPeriod"`
	InsuranceMarketRisk                          *string                         `json:"insuranceMarketRisk"`
	InsuranceStorageRisk                         *string                         `json:"insuranceStorageRisk"`
	InsuranceTitleRisk                           *string                         `json:"insuranceTitleRisk"`
	InsuranceFraudRisk                           *string                         `json:"insuranceFraudRisk"`
	InsuranceInsuranceRisk                       *string                         `json:"insuranceInsuranceRisk"`
	InsuranceOperationalRisk                     *string                         `json:"insuranceOperationalRisk"`
	InsuranceRegulatoryRisk                      *string                         `json:"insuranceRegulatoryRisk"`
	InsuranceLiquidityRisk                       *string                         `json:"insuranceLiquidityRisk"`
	InsuranceForceMajeureRisk                    *string                         `json:"insuranceForceMajeureRisk"`
	InsuranceEarlyRedemptionRisk                 *string                         `json:"insuranceEarlyRedemptionRisk"`
	InsuranceMitigationMeasures                  *string                         `json:"insuranceMitigationMeasures"`
	InsuranceInsuranceCoverageSummary            *string                         `json:"insuranceInsuranceCoverageSummary"`
	InsuranceInsuranceProvider                   *string                         `json:"insuranceInsuranceProvider"`
	InsuranceCoverageValue                       *string                         `json:"insuranceCoverageValue"`
	IssuerRegistrationNo                         *string                         `json:"issuerRegistrationNo"`
	SectorAndIndustry                            *string                         `json:"sectorAndIndustry"`
	LicenseOrPermitNumber                        *string                         `json:"licenseOrPermitNumber"`
	IssuerAdditionalInfo                         *string                         `json:"issuerAdditionalInfo"`
	IsinSerialNumber                             *string                         `json:"isinSerialNumber"`
	EsgOrImpactMetrics                           *string                         `json:"esgOrImpactMetrics"`
	InstrumentAdditionalInfo                     *string                         `json:"instrumentAdditionalInfo"`
	SecurityType                                 *string                         `json:"securityType"`
	CollateralDescription                        *string                         `json:"collateralDescription"`
	CovenantSummary                              *string                         `json:"covenantSummary"`
	CovenantTestingFrequency                     *string                         `json:"covenantTestingFrequency"`
	EventOfDefaultClauses                        *string                         `json:"eventOfDefaultClauses"`
	LegalEnforcementMechanism                    *string                         `json:"legalEnforcementMechanism"`
	Guarantee                                    *string                         `json:"guarantee"`
	RecoveryEstimate                             float64                         `gorm:"default:0" json:"recoveryEstimate"`
	RiskProfileAdditionalInfo                    *string                         `json:"riskProfileAdditionalInfo"`
	LicenseNumber                                *string                         `json:"licenseNumber"`
	ExitLoadFee                                  float64                         `gorm:"default:0" json:"exitLoadFee"`
	EntryLoadFee                                 float64                         `gorm:"default:0" json:"entryLoadFee"`
	PortfolioLockInPeriod                        int                             `gorm:"default:0" json:"portfolioLockInPeriod"`
	PortfolioPerformanceFee                      float64                         `gorm:"default:0" json:"portfolioPerformanceFee"`
	FundingStructure                             int                             `gorm:"default:0" json:"fundingStructure"` //0=equity, 1= debt, 2= hybrid
	EquityPercentage                             float64                         `gorm:"default:0" json:"equityPercentage"`
	DebtPercentage                               float64                         `gorm:"default:0" json:"debtPercentage"`
	DebtInstrumentType                           *string                         `json:"debtInstrumentType"`
	PrincipalPaymentMethod                       *string                         `json:"principalPaymentMethod"`
	DebtInstrumentRepaymentSource                *string                         `json:"debtInstrumentRepaymentSource"`
	DebtInstrumentGuaranteesOrEnhancements       *string                         `json:"debtInstrumentGuaranteesOrEnhancements"`
	DebtInstrumentDefaultAndRecoveryTerms        *string                         `json:"debtInstrumentDefaultAndRecoveryTerms"`
	DebtInstrumentRepaymentFrequency             *string                         `json:"debtInstrumentRepaymentFrequency"`
	InterestRepaymentFrequency                   *string                         `json:"interestRepaymentFrequency"`
	DcsrDetails                                  *string                         `json:"dcsrDetails"`
	SinkingFundStructure                         *string                         `json:"sinkingFundStructure"`
	CovenantMonitoringAgent                      *string                         `json:"covenantMonitoringAgent"`
	RightOfRecourse                              *string                         `json:"rightOfRecourse"`
	DebtInstrumentInterestRate                   float64                         `gorm:"default:0" json:"debtInstrumentInterestRate"`
	DcsrRatio                                    float64                         `gorm:"default:0" json:"dcsrRatio"`
	LtvRatio                                     float64                         `gorm:"default:0" json:"ltvRatio"`
	InterestCoverageRatio                        float64                         `gorm:"default:0" json:"interestCoverageRatio"`
	MaximumLeverageRatio                         float64                         `gorm:"default:0" json:"maximumLeverageRatio"`
	GracePeriod                                  int                             `gorm:"default:0" json:"gracePeriod"`
	TrusteeAppointed                             int                             `gorm:"default:0" json:"trusteeAppointed"`
	ReserveFundInPlace                           int                             `gorm:"default:0" json:"reserveFundInPlace"`
	SecurityOrCollateralOffered                  *string                         `json:"securityOrCollateralOffered"`
	FundInstrumentType                           *string                         `json:"fundInstrumentType"`
	InstrumentRatingAgency                       *string                         `json:"instrumentRatingAgency"`
	PortfolioTopHoldings                         *string                         `json:"portfolioTopHoldings"`
	CreditRatingAgency                           *string                         `json:"CreditRatingAgency"`
	TrusteeRegNumber                             *string                         `json:"trusteeRegNumber"`
	CreditEnhancerOrGuarantor                    *string                         `json:"creditEnhancerOrGuarantor"`
	BondStructuringAdvisor                       *string                         `json:"bondStructuringAdvisor"`
	EntitiesAdditionalInfo                       *string                         `json:"entitiesAdditionalInfo"`
	AuthorizedRepresentativeName                 *string                         `json:"authorizedRepresentativeName"`
	AuthorizedRepresentativeTitleOrPosition      *string                         `json:"authorizedRepresentativeTitleOrPosition"`
	AuthorizedRepresentativeEmail                *string                         `json:"authorizedRepresentativeEmail"`
	AcceptTokenizationTermsAndAgreement          int                             `gorm:"default:0" json:"acceptTokenizationTermsAndAgreement"`
	AttestInformationAccurateAndVerifiable       int                             `gorm:"default:0" json:"attestInformationAccurateAndVerifiable"`
	AcknowledgedSuitabilityCriteria              int                             `gorm:"default:0" json:"acknowledgedSuitabilityCriteria"`
	MinimumKycTier                               *string                         `json:"minimumKycTier"`
	InvestorCategory                             *string                         `json:"investorCategory"`
	WithholdingTaxDisclosure                     int                             `gorm:"default:0" json:"withholdingTaxDisclosure"`
}

type TokenizedAssetID string

type TokenizedAssetJSONInput struct {
	AssetSector                                  string    `json:"assetSector"`
	AssetSubSector                               string    `json:"assetSubSector"`
	AssetType                                    string    `json:"assetType"`
	AssetName                                    string    `json:"assetName"`
	AssetWebsite                                 string    `json:"assetWebsite"`
	ApprovedAssetCustodianID                     uint64    `gorm:"not null" json:"approvedAssetCustodianId"`
	OfferingType                                 string    `gorm:"default:'PRIVATE'" json:"offeringType"` //PRIVATE, PUBLIC
	ClosedGroupID                                string    `gorm:"null" json:"closedGroupId"`
	SecApproval                                  int       `gorm:"default:0" json:"secApproval"`
	SecApprovalIdNumber                          string    `json:"secApprovalIdNumber"`
	MarketMakingWallet                           string    `json:"marketMakingWallet"`
	AssetDescription                             string    `json:"assetDescription"`
	AssetCountryLocation                         string    `json:"assetCountryLocation"`
	AssetPhysicalAddress                         string    `json:"assetPhysicalAddress"`
	AssetLongitude                               string    `json:"assetLongitude"`
	AssetLatitude                                string    `json:"assetLatitude"`
	OwnershipType                                string    `json:"ownershipType"` //DIRECT,THIRD-PARTY
	OwnershipKind                                string    `json:"ownershipKind"` //INDIVIDUAL,CORPORATE,
	InitialOwnerPreferredWalletAddress           string    `json:"initialOwnerPreferredWalletAddress"`
	AssetOwnerName                               string    `json:"assetOwnerName"`
	AssetOwnerRetainedOrContributedValue         float64   `json:"assetOwnerRetainedOrContributedValue"`
	AssetOwnerAddress                            string    `json:"assetOwnerAddress"`
	AssetManagerID                               uint64    `json:"assetManagerId"`
	AssetIssuingHouseID                          uint64    `gorm:"not null" json:"assetIssuingHouseId"`
	LegalAndProfesionalPartnerID                 uint64    `gorm:"not null" json:"legalAndProfesionalPartnerId"`
	RatingAgencyID                               uint64    `gorm:"not null" json:"ratingAgencyId"`
	TrusteeID                                    uint64    `gorm:"not null" json:"trusteeId"`
	AssetQuoteCurrency                           string    `gorm:"default:'CNGN'" json:"assetQuoteCurrency"`
	AssetCurrentValue                            float64   `gorm:"default:0" json:"assetCurrentValue"`
	AssetMscCostOutisdeOfValuation               float64   `gorm:"default:0" json:"assetMscCostOutisdeOfValuation"`
	ProtectionMethods                            string    `json:"protectionMethods"` //csv format
	InsuranceCompanyName                         string    `json:"insuranceCompanyName"`
	InsurancePolicyNumber                        string    `json:"insurancePolicyNumber"`
	InsurancePolicyHolder                        string    `json:"insurancePolicyHolder"`
	PercentageValueOfInsurance                   float64   `gorm:"default:0" json:"percentageValueOfInsurance"`
	IsFreeFromLiensAndEncumbrances               int       `gorm:"default:0" json:"IsFreeFromLiensAndEncumbrances"`
	AssetAlreadyExists                           int       `gorm:"default:1" json:"assetAlreadyExists"`
	AssetCode                                    string    `json:"assetCode"`
	AssetLogo                                    string    `json:"assetLogo"`
	NumberOfTokenToBeIssued                      float64   `json:"numberOfTokenToBeIssued"`
	NumberOfTokenToBeSold                        float64   `json:"numberOfTokenToBeSold"`
	TotalTokenHeldByManager                      float64   `json:"totalTokenHeldByManager"`
	WalletToHoldAssetsNotForSale                 string    `json:"walletToHoldAssetsNotForSale"` //wallet that the original owner wants to use to receive their portion of tokenized asset that are not meant for sale.
	PricePerToken                                float64   `json:"pricePerToken"`
	SalesStart                                   time.Time `json:"salesStart"`
	SalesEnd                                     time.Time `json:"salesEnd"`
	CapOnPurchase                                int       `gorm:"default:0" json:"capOnPurchase"`
	CapQuantity                                  float64   `gorm:"default:0" json:"capQuantity"`
	CapAmountInFiat                              float64   `gorm:"default:0" json:"capAmountInFiat"`
	CapDurationInDays                            int       `gorm:"default:0" json:"capDurationInDays"`
	ProceedCycle                                 string    `gorm:"size:50" json:"proceedCycle"`
	TokenizationFeeID                            uint64    `json:"tokenizationFeeId"`
	ProceedPayoutCurrency                        string    `json:"proceedPayoutCurrency"`
	ProceedPayoutType                            int       `gorm:"default:0" json:"proceedPayoutType"` // FIAT=1, CRYPTO=0. FIAT requires fiat payment method.
	ExemptedCountries                            string    `json:"exemptedCountries"`
	HasAdditionalKYCRequirements                 int       `gorm:"default:0" json:"hasAdditionalKYCRequirements"`
	AdditionalKYCRequirements                    string    `json:"additionalKYCRequirements"`
	InvestorAccreditationRequired                int       `gorm:"default:0" json:"investorAccreditationRequired"`
	Messages                                     []string  `json:"messages"`
	AgreeTransferTitleToCustodian                int       `gorm:"default:0" json:"agreeTransferTitleToCustodian"`
	ContractualProtectionRevGuarantees           int       `gorm:"default:0" json:"contractualProtectionRevGuarantees"`
	ContractualProtectionPerfBond                int       `gorm:"default:0" json:"contractualProtectionPerfBond"`
	ContractualProtectionSLA                     int       `gorm:"default:0" json:"contractualProtectionSLA"`
	RiskSharingMechanismPPPs                     int       `gorm:"default:0" json:"riskSharingMechanismPPPs"`
	RiskSharingMechanismHedgeInstruments         int       `gorm:"default:0" json:"riskSharingMechanismHedgeInstruments"`
	RiskSharingMechanismCompletionGuarantees     int       `gorm:"default:0" json:"riskSharingMechanismCompletionGuarantees"`
	IndependentMonitoringList                    string    `gorm:"null" json:"independentMonitoringList"`
	ESGSafeguardsSusCerts                        int       `gorm:"default:0" json:"eSGSafeguardsSusCerts"`
	ESGSafeguardsCommEngPlans                    int       `gorm:"default:0" json:"eSGSafeguardsCommEngPlans"`
	SecurityMeasuresAccessControl                int       `gorm:"default:0" json:"securityMeasuresAccessControl"`
	SecurityMeasuresSurveilanceSystems           int       `gorm:"default:0" json:"securityMeasuresSurveilanceSystems"`
	SecurityMeasuresOnSiteSecurityPersonnel      int       `gorm:"default:0" json:"securityMeasuresOnSiteSecurityPersonnel"`
	SecurityMeasuresPerimeterSecurity            int       `gorm:"default:0" json:"securityMeasuresPerimeterSecurity"`
	SecurityMeasuresCriticalInfraProtections     int       `gorm:"default:0" json:"securityMeasuresCriticalInfraProtections"`
	OtherAssetProtection                         string    `gorm:"null" json:"otherAssetProtection"`
	LegalAdvisor                                 string    `gorm:"null" json:"legalAdvisor"`
	FinancialAdvisor                             string    `gorm:"null" json:"financialAdvisor"`
	UndertakingNoLien                            int       `gorm:"default:0" json:"undertakingNoLien"`
	UndertakingNotCollateral                     int       `gorm:"default:0" json:"undertakingNotCollateral"`
	UndertakingNoClaims                          int       `gorm:"default:0" json:"undertakingNoClaims"`
	UndertakingNoForeclosure                     int       `gorm:"default:0" json:"undertakingNoForeclosure"`
	ComplianceNoViolation                        int       `gorm:"default:0" json:"complianceNoViolation"`
	ComplianceAllPermits                         int       `gorm:"default:0" json:"complianceAllPermits"`
	OutstandingFinancialRespNoDebts              int       `gorm:"default:0" json:"outstandingFinancialRespNoDebts"`
	OutstandingFinancialRespNoHiddenLiabilities  int       `gorm:"default:0" json:"outstandingFinancialRespNoHiddenLiabilities"`
	RiskManagementFullyInsured                   int       `gorm:"default:0" json:"riskManagementFullyInsured"`
	RiskManagementDeclaredValue                  int       `gorm:"default:0" json:"riskManagementDeclaredValue"`
	PhysicalConditionSound                       int       `gorm:"default:0" json:"physicalConditionSound"`
	PhysicalConditionNolease                     int       `gorm:"default:0" json:"physicalConditionNolease"`
	PhysicalConditionNoUndisclosedEasements      int       `gorm:"default:0" json:"physicalConditionNoUndisclosedEasements"`
	MintingInitators                             string    `gorm:"null" json:"mintingInitators"` //CSV of approvers
	MintingApprovers                             string    `gorm:"null" json:"mintingApprovers"` //csv of initators
	BankID                                       uint64    `gorm:"null" json:"bankId"`
	AccountNumber                                string    `gorm:"null" json:"accountNumber"`
	BeneficiaryName                              string    `gorm:"null" json:"beneficiaryName"`
	ProjectStrategicObjectives                   string    `json:"projectStrategicObjectives"`
	ProjectDevelopmentTimeline                   string    `json:"projectDevelopmentTimeline"`
	ProjectKeyMilestoneAndDates                  string    `json:"projectKeyMilestoneAndDates"`
	ProjectScope                                 string    `json:"projectScope"`
	ProjectEconomicBenefits                      string    `json:"projectEconomicBenefits"`
	ProjectExpectedNoOfJobs                      int       `json:"projectExpectedNoOfJobs"`
	ProjectIntendedSocialBenefits                string    `json:"projectIntendedSocialBenefits"`
	ProjectTechnicalPartners                     string    `json:"projectTechnicalPartners"`
	ProjectFinancialPartners                     string    `json:"projectFinancialPartners"`
	EstimatedProjectIRR                          float64   `json:"estimatedProjectIRR"`
	EstimatedProjectROI                          float64   `json:"estimatedProjectROI"`
	EstimatedProjectNPV                          float64   `json:"estimatedProjectNPV"`
	EstimatedProjectPaybackPeriodsInMonths       int       `json:"estimatedProjectPaybackPeriodsInMonths"`
	KeyAssumptionsList                           string    `json:"keyAssumptionsList"`
	ProjectIdentifiedLegalRisks                  string    `json:"projectIdentifiedLegalRisks"`
	ProjectIdentifiedRegulatoryRisks             string    `json:"projectIdentifiedRegulatoryRisks"`
	ProjectIdentifiedOperationalOrExecutionRisks string    `json:"projectIdentifiedOperationalOrExecutionRisks"`
	ProjectIdentifiedMarketRisks                 string    `json:"projectIdentifiedMarketRisks"`
	ProjectIdentifiedOtherRelevantRisks          string    `json:"projectIdentifiedOtherRelevantRisks"`
	//////
	IsinOrSerialNumber                      string    `json:"isinOrSerialNumber"` /// begin the mutual fund fields
	InstrumentName                          string    `json:"instrumentName"`
	InstrumentType                          string    `json:"instrumentType"`
	TotalIssueSize                          float64   `json:"totalIssueSize"`
	IssueDate                               time.Time `json:"issueDate"`
	MaturityDate                            time.Time `json:"maturityDate"`
	FaceValuePerUnit                        float64   `json:"faceValuePerUnit"`
	CouponOrInterestRateType                string    `json:"couponOrInterestRateType"`
	CouponOrInterestRate                    float64   `json:"couponOrInterestRate"`
	ReferenceIndex                          string    `json:"referenceIndex"`
	ExpectedYield                           float64   `json:"expectedYield"`
	EarlyRedemptionOptionInvestor           string    `json:"earlyRedemptionOptionInvestor"`
	MinimumInvestmentAmount                 float64   `json:"minimumInvestmentAmount"`
	TaxTreatmentTokenHolders                string    `json:"taxTreatmentTokenHolders"`
	PaymentStructureToTokenHolders          string    `json:"paymentStructureToTokenHolders"`
	RedemptionMethod                        string    `json:"redemptionMethod"`
	PaymentStructure                        string    `json:"paymentStructure"`
	RepaymentMethod                         string    `json:"repaymentMethod"`
	PaymentCycle                            string    `json:"paymentCycle"`
	WeightedAverageLife                     float64   `json:"weightedAverageLife"`
	UnderlyingAssetPoolSize                 float64   `json:"underlyingAssetPoolSize"`
	PoolComposition                         string    `json:"poolComposition"`
	CreditEnhancementMethod                 string    `json:"creditEnhancementMethod"`
	SummaryOfUseOfProceeds                  string    `json:"summaryOfUseOfProceeds"`
	CreditRatingIfAny                       string    `json:"creditRatingIfAny"`
	IssuerName                              string    `json:"issuerName"`
	IssuerType                              string    `json:"issuerType"`
	IssuerContactPerson                     string    `json:"issuerContactPerson"`
	ContactEmail                            string    `json:"contactEmail"`
	ContactPhoneNumber                      string    `json:"contactPhoneNumber"`
	BriefCompanyOverview                    string    `json:"briefCompanyOverview"`
	MortgageOriginators                     string    `json:"mortgageOriginators"`
	Servicer                                string    `json:"servicer"`
	InstrumentTrustee                       string    `json:"instrumentTrustee"`
	InstrumentCustodians                    string    `json:"instrumentCustodian"`
	InstrumentLegalAdvisor                  string    `json:"instrumentLegalAdvisor"`
	SpecialPurposeVehicle                   string    `json:"specialPurposeVehicle"`
	InstrumentAssetManagerOrAdministrator   string    `json:"instrumentAssetManagerOrAdministrator"`
	UnderwriterIfAny                        string    `json:"underwriterIfAny"`
	InstrumentCreditRatingAgency            string    `json:"instrumentCreditRatingAgency"`
	AuditorOrVerifier                       string    `json:"auditorOrVerifier"`
	CreditRiskAssessment                    string    `json:"creditRiskAssessment"`
	CreditRating                            string    `json:"creditRating"`
	PrepaymentRisk                          string    `json:"prepaymentRisk"`
	InterestRateRisk                        string    `json:"interestRateRisk"`
	StructuralComplexityRisk                string    `json:"structuralComplexityRisk"`
	LegalOrRegulatoryRisk                   string    `json:"legalOrRegulatoryRisk"`
	OperationalRisk                         string    `json:"operationalRisk"`
	MarketRisk                              string    `json:"marketRisk"`
	EsgRisk                                 string    `json:"esgRisk"`
	MitigationMeasures                      string    `json:"mitigationMeasures"`
	IssuingAuthority                        string    `json:"issuingAuthority"`
	RegulatoryApprovalId                    string    `json:"regulatoryApprovalId"`
	LicenseApprovalReferenceNumber          string    `json:"licenseApprovalReferenceNumber"`
	ListingStatus                           string    `json:"listingStatus"`
	CurrencyOfIssuance                      string    `json:"currencyOfIssuance"`
	CouponRateType                          string    `json:"couponRateType"`
	CouponRate                              float64   `json:"couponRate"`
	SpreadOrMargin                          float64   `json:"spreadOrMargin"`
	ResetFrequency                          string    `json:"resetFrequency"`
	CouponPaymentFrequency                  string    `json:"couponPaymentFrequency"`
	RedemptionStructure                     string    `json:"redemptionStructure"`
	EarlyRedemptionOption                   string    `json:"earlyRedemptionOption"`
	EarlyRedemptionPenalty                  string    `json:"earlyRedemptionPenalty"`
	TaxTreatment                            string    `json:"taxTreatment"`
	NavOrMarketValueUpdates                 string    `json:"navOrMarketValueUpdates"`
	ImpactMetrics                           string    `json:"impactMetrics"`
	LegalBacking                            string    `json:"legalBacking"`
	DefaultHistory                          string    `json:"defaultHistory"`
	RiskFactorsSummary                      string    `json:"riskFactorsSummary"`
	PayingAgent                             string    `json:"payingAgent"`
	Auditor                                 string    `json:"auditor"`
	RegistrarOrCSCSAgent                    string    `json:"registrarOrCSCSAgent"`
	FundStructure                           string    `json:"fundStructure"`
	AssetManagementCompanyName              string    `json:"assetManagementCompanyName"`
	FundManagers                            string    `json:"fundManagers"`
	RegulatoryLicenseNumber                 string    `json:"regulatoryLicenseNumber"`
	FundLaunchDate                          time.Time `json:"fundLaunchDate"`
	TotalExpenseRatio                       float64   `json:"totalExpenseRatio"`
	ExitLoadRedemptionFee                   float64   `json:"exitLoadRedemptionFee"`
	Tenure                                  int       `json:"tenure"`
	InitialNetAssetValue                    float64   `json:"initialNetAssetValue"`
	NavUpdateFrequency                      string    `json:"navUpdateFrequency"`
	NavCalculationMethod                    string    `json:"navCalculationMethod"`
	RedemptionRules                         string    `json:"redemptionRules"`
	LockInPeriod                            int       `json:"lockInPeriod"`
	EntryLoad                               float64   `json:"entryLoad"`
	PerformanceFee                          float64   `json:"performanceFee"`
	DividendPolicy                          string    `json:"dividendPolicy"`
	LiquidityProfile                        string    `json:"liquidityProfile"`
	DistributionFrequency                   string    `json:"distributionFrequency"`
	DistributionMethod                      string    `json:"distributionMethod"`
	BenchmarkComparisonMethod               string    `json:"benchmarkComparisonMethod"`
	FeeBreakdownSummary                     string    `json:"feeBreakdownSummary"`
	InvestmentObjective                     string    `json:"investmentObjective"`
	EquityStrategy                          string    `json:"equityStrategy"`
	MarketCapitalizationFocus               string    `json:"marketCapitalizationFocus"`
	BenchmarkIndex                          string    `json:"benchmarkIndex"`
	SectorExposureLimits                    string    `json:"sectorExposureLimits"`
	TopHoldings                             string    `json:"topHoldings"`
	GeographicExposure                      string    `json:"geographicExposure"`
	RiskProfile                             string    `json:"riskProfile"`
	VolatilityEstimate                      float64   `json:"volatilityEstimate"`
	DividendYield                           float64   `json:"dividendYield"`
	TrusteeName                             string    `json:"trusteeName"`
	FundAdministrator                       string    `json:"fundAdministrator"`
	InvestmentCommitteeMembers              string    `json:"investmentCommitteeMembers"`
	IsinOrSecFundCode                       string    `json:"isinOrSecFundCode"`
	ExitLoadOrRedemptionFee                 float64   `json:"exitLoadOrRedemptionFee"`
	FundRiskRating                          string    `json:"fundRiskRating"`
	AssetAllocation                         string    `json:"assetAllocation"`
	AverageMaturity                         float64   `json:"averageMaturity"`
	YieldToMaturity                         float64   `json:"yieldToMaturity"`
	CreditRatingProfile                     string    `json:"creditRatingProfile"`
	LockInPeriodPortfolio                   int       `json:"lockInPeriodPortfolio"`
	PerformanceFeeIfAny                     float64   `json:"performanceFeeIfAny"`
	TopEquityHoldings                       string    `json:"topEquityHoldings"`
	TargetAllocation                        string    `json:"targetAllocation"`
	AllowedAllocationRange                  string    `json:"allowedAllocationRange"`
	AssetClassesIncluded                    string    `json:"assetClassesIncluded"`
	RebalancingFrequency                    string    `json:"rebalancingFrequency"`
	BenchmarkIndexComposite                 string    `json:"benchmarkIndexComposite"`
	TopEquityHoldingsList                   string    `json:"topEquityHoldingsList"`
	TopDebtHoldings                         string    `json:"topDebtHoldings"`
	CreditRatingDistribution                string    `json:"creditRatingDistribution"`
	AverageMaturityHybrid                   float64   `json:"averageMaturityHybrid"`
	YieldToMaturityHybrid                   float64   `json:"yieldToMaturityHybrid"`
	TitleOfIssuance                         string    `json:"titleOfIssuance"`
	TypeOfCommercialPaper                   string    `json:"typeOfCommercialPaper"`
	PricingYield                            float64   `json:"pricingYield"`
	UseOfProceeds                           string    `json:"useOfProceeds"`
	Ranking                                 string    `json:"ranking"`
	BackingSecurity                         string    `json:"backingSecurity"`
	IssuerRegistrationNumber                string    `json:"issuerRegistrationNumber"`
	IncorporationDate                       time.Time `json:"incorporationDate"`
	RcNumber                                string    `json:"rcNumber"`
	TaxIdNumber                             string    `json:"taxIdNumber"`
	OfficeAddress                           string    `json:"officeAddress"`
	Rating                                  string    `json:"rating"`
	PartiesInvolvedIssuer                   string    `json:"partiesInvolvedIssuer"`
	PartiesInvolvedArranger                 string    `json:"partiesInvolvedArranger"`
	PartiesInvolvedLegalAdviser             string    `json:"partiesInvolvedLegalAdviser"`
	PartiesInvolvedAuditor                  string    `json:"partiesInvolvedAuditor"`
	PartiesInvolvedRatingAgency             string    `json:"partiesInvolvedRatingAgency"`
	PartiesInvolvedCustodian                string    `json:"partiesInvolvedCustodian"`
	PartiesInvolvedTrustee                  string    `json:"partiesInvolvedTrustee"`
	PartiesInvolvedAuditorVerifier          string    `json:"partiesInvolvedAuditorVerifier"`
	SecurityRiskLegalBacking                string    `json:"securityRiskLegalBacking"`
	SecurityRiskCollateral                  string    `json:"securityRiskCollateral"`
	SecurityRiskDefaultHistory              string    `json:"securityRiskDefaultHistory"`
	SecurityRiskCreditRating                string    `json:"securityRiskCreditRating"`
	SecurityRiskRiskFactorsSummary          string    `json:"securityRiskRiskFactorsSummary"`
	SecurityRiskBusinessRisk                string    `json:"securityRiskBusinessRisk"`
	SecurityRiskDefaultRisk                 string    `json:"securityRiskDefaultRisk"`
	SecurityRiskLiquidityRisk               string    `json:"securityRiskLiquidityRisk"`
	SecurityRiskRegulatoryRisk              string    `json:"securityRiskRegulatoryRisk"`
	SecurityRiskMarketRisk                  string    `json:"securityRiskMarketRisk"`
	SecurityRiskOperationalRisk             string    `json:"securityRiskOperationalRisk"`
	SecurityRiskMitigationMeasures          string    `json:"securityRiskMitigationMeasures"`
	CommodityType                           string    `json:"commodityType"`
	CommodityDescription                    string    `json:"commodityDescription"`
	Quantity                                int       `json:"quantity"`
	QualityGrade                            string    `json:"qualityGrade"`
	IssuerContactInfo                       string    `json:"issuerContactInfo"`
	WarehouseName                           string    `json:"warehouseName"`
	WarehouseOperatorName                   string    `json:"warehouseOperatorName"`
	WarehouseLicenseNumber                  string    `json:"warehouseLicenseNumber"`
	WarehouseLocation                       string    `json:"warehouseLocation"`
	WrNumber                                string    `json:"wrNumber"`
	WrIssueDate                             time.Time `json:"wrIssueDate"`
	WrExpiryDate                            time.Time `json:"wrExpiryDate"`
	WrSystemRegistration                    string    `json:"wrSystemRegistration"`
	WrRegistrationNumber                    string    `json:"wrRegistrationNumber"`
	WrVerifier                              string    `json:"wrVerifier"`
	StorageCondition                        string    `json:"storageCondition"`
	WarehouseAccreditationBody              string    `json:"warehouseAccreditationBody"`
	MinimumPurchaseAmount                   float64   `json:"minimumPurchaseAmount"`
	AutoRollover                            string    `json:"autoRollover"`
	CurrentBeneficialOwner                  string    `json:"currentBeneficialOwner"`
	WrCustodianName                         string    `json:"wrCustodianName"`
	OwnershipRightsRepresented              string    `json:"ownershipRightsRepresented"`
	TrusteeOrThirdPartyOversight            string    `json:"trusteeOrThirdPartyOversight"`
	LienOrEncumbrances                      string    `json:"lienOrEncumbrances"`
	AssetValuation                          float64   `json:"assetValuation"`
	ValuationDate                           time.Time `json:"valuationDate"`
	ValuationMethodology                    string    `json:"valuationMethodology"`
	TokenizationObjective                   string    `json:"tokenizationObjective"`
	HoldingPeriod                           int       `json:"holdingPeriod"`
	RedemptionMechanism                     string    `json:"redemptionMechanism"`
	PartiesInvolvedUnderwriter              string    `json:"partiesInvolvedUnderwriter"`
	PartiesInvolvedAssetManager             string    `json:"partiesInvolvedAssetManager"`
	PartiesInvolvedLegalAdvisor             string    `json:"partiesInvolvedLegalAdvisor"`
	PartiesInvolvedRegulator                string    `json:"partiesInvolvedRegulator"`
	RisksMarketRisk                         string    `json:"risksMarketRisk"`
	RisksStorageRisk                        string    `json:"risksStorageRisk"`
	RisksTitleRisk                          string    `json:"risksTitleRisk"`
	RisksFraudRisk                          string    `json:"risksFraudRisk"`
	RisksInsuranceRisk                      string    `json:"risksInsuranceRisk"`
	RisksOperationalRisk                    string    `json:"risksOperationalRisk"`
	RisksRegulatoryRisk                     string    `json:"risksRegulatoryRisk"`
	RisksLiquidityRisk                      string    `json:"risksLiquidityRisk"`
	RisksForceMajeureRisk                   string    `json:"risksForceMajeureRisk"`
	RisksEarlyRedemptionRisk                string    `json:"risksEarlyRedemptionRisk"`
	RisksMitigationMeasures                 string    `json:"risksMitigationMeasures"`
	RisksInsuranceCoverageSummary           string    `json:"risksInsuranceCoverageSummary"`
	RisksInsuranceProvider                  string    `json:"risksInsuranceProvider"`
	RisksCoverageValue                      float64   `json:"risksCoverageValue"`
	QuanlityStandard                        string    `json:"quanlityStandard"`
	IssuerContactInformation                string    `json:"issuerContactInformation"`
	VaultCustodianName                      string    `json:"vaultCustodianName"`
	VaultOperator                           string    `json:"vaultOperator"`
	VaultLicenseNumber                      string    `json:"vaultLicenseNumber"`
	VaultLocation                           string    `json:"vaultLocation"`
	Number                                  string    `json:"number"`
	IssuerDate                              time.Time `json:"issuerDate"`
	ExpiryDate                              time.Time `json:"expiryDate"`
	RegistryRecord                          string    `json:"registryRecord"`
	Verifier                                string    `json:"verifier"`
	StorageConditions                       string    `json:"storageConditions"`
	VaultAccreditationBody                  string    `json:"vaultAccreditationBody"`
	OwnershipLegalHolder                    string    `json:"ownershipLegalHolder"`
	OwnershipCustodianName                  string    `json:"ownershipCustodianName"`
	OwnershipTrustee                        string    `json:"ownershipTrustee"`
	OwnershipLienOrEncumbrances             string    `json:"ownershipLienOrEncumbrances"`
	ValuationAssetValuation                 string    `json:"valuationAssetValuation"`
	HoldingLockinPeriod                     int       `json:"holdingLockinPeriod"`
	InsuranceMarketRisk                     string    `json:"insuranceMarketRisk"`
	InsuranceStorageRisk                    string    `json:"insuranceStorageRisk"`
	InsuranceTitleRisk                      string    `json:"insuranceTitleRisk"`
	InsuranceFraudRisk                      string    `json:"insuranceFraudRisk"`
	InsuranceInsuranceRisk                  string    `json:"insuranceInsuranceRisk"`
	InsuranceOperationalRisk                string    `json:"insuranceOperationalRisk"`
	InsuranceRegulatoryRisk                 string    `json:"insuranceRegulatoryRisk"`
	InsuranceLiquidityRisk                  string    `json:"insuranceLiquidityRisk"`
	InsuranceForceMajeureRisk               string    `json:"insuranceForceMajeureRisk"`
	InsuranceEarlyRedemptionRisk            string    `json:"insuranceEarlyRedemptionRisk"`
	InsuranceMitigationMeasures             string    `json:"insuranceMitigationMeasures"`
	InsuranceInsuranceCoverageSummary       string    `json:"insuranceInsuranceCoverageSummary"`
	InsuranceInsuranceProvider              string    `json:"insuranceInsuranceProvider"`
	InsuranceCoverageValue                  string    `json:"insuranceCoverageValue"`
	IssuerRegistrationNo                    string    `json:"issuerRegistrationNo"`
	SectorAndIndustry                       string    `json:"sectorAndIndustry"`
	LicenseOrPermitNumber                   string    `json:"licenseOrPermitNumber"`
	IssuerAdditionalInfo                    string    `json:"issuerAdditionalInfo"`
	IsinSerialNumber                        string    `json:"isinSerialNumber"`
	EsgOrImpactMetrics                      string    `json:"esgOrImpactMetrics"`
	InstrumentAdditionalInfo                string    `json:"instrumentAdditionalInfo"`
	SecurityType                            string    `json:"securityType"`
	CollateralDescription                   string    `json:"collateralDescription"`
	CovenantSummary                         string    `json:"covenantSummary"`
	CovenantTestingFrequency                string    `json:"covenantTestingFrequency"`
	EventOfDefaultClauses                   string    `json:"eventOfDefaultClauses"`
	LegalEnforcementMechanism               string    `json:"legalEnforcementMechanism"`
	Guarantee                               string    `json:"guarantee"`
	RecoveryEstimate                        float64   `json:"recoveryEstimate"`
	RiskProfileAdditionalInfo               string    `json:"riskProfileAdditionalInfo"`
	LicenseNumber                           string    `json:"licenseNumber"`
	ExitLoadFee                             float64   `json:"exitLoadFee"`
	EntryLoadFee                            float64   `json:"entryLoadFee"`
	PortfolioLockInPeriod                   int       `json:"portfolioLockInPeriod"`
	PortfolioPerformanceFee                 float64   `json:"portfolioPerformanceFee"`
	FundingStructure                        int       `gorm:"default:0" json:"fundingStructure"` //0=equity, 1= debt, 2= hybrid
	EquityPercentage                        float64   `gorm:"default:0" json:"equityPercentage"`
	DebtPercentage                          float64   `gorm:"default:0" json:"debtPercentage"`
	DebtInstrumentType                      string    `json:"debtInstrumentType"`
	PrincipalPaymentMethod                  string    `json:"principalPaymentMethod"`
	DebtInstrumentRepaymentSource           string    `json:"debtInstrumentRepaymentSource"`
	DebtInstrumentGuaranteesOrEnhancements  string    `json:"debtInstrumentGuaranteesOrEnhancements"`
	DebtInstrumentDefaultAndRecoveryTerms   string    `json:"debtInstrumentDefaultAndRecoveryTerms"`
	DebtInstrumentRepaymentFrequency        string    `json:"debtInstrumentRepaymentFrequency"`
	InterestRepaymentFrequency              string    `json:"interestRepaymentFrequency"`
	DcsrDetails                             string    `json:"dcsrDetails"`
	SinkingFundStructure                    string    `json:"sinkingFundStructure"`
	CovenantMonitoringAgent                 string    `json:"covenantMonitoringAgent"`
	RightOfRecourse                         string    `json:"rightOfRecourse"`
	DebtInstrumentInterestRate              float64   `gorm:"default:0" json:"debtInstrumentInterestRate"`
	DcsrRatio                               float64   `gorm:"default:0" json:"dcsrRatio"`
	LtvRatio                                float64   `gorm:"default:0" json:"ltvRatio"`
	InterestCoverageRatio                   float64   `gorm:"default:0" json:"interestCoverageRatio"`
	MaximumLeverageRatio                    float64   `gorm:"default:0" json:"maximumLeverageRatio"`
	GracePeriod                             int       `gorm:"default:0" json:"gracePeriod"`
	TrusteeAppointed                        int       `gorm:"default:0" json:"trusteeAppointed"`
	ReserveFundInPlace                      int       `gorm:"default:0" json:"reserveFundInPlace"`
	SecurityOrCollateralOffered             string    `json:"securityOrCollateralOffered"`
	FundInstrumentType                      string    `json:"fundInstrumentType"`
	InstrumentRatingAgency                  string    `json:"instrumentRatingAgency"`
	PortfolioTopHoldings                    string    `json:"portfolioTopHoldings"`
	CreditRatingAgency                      string    `json:"CreditRatingAgency"`
	TrusteeRegNumber                        string    `json:"trusteeRegNumber"`
	CreditEnhancerOrGuarantor               string    `json:"creditEnhancerOrGuarantor"`
	BondStructuringAdvisor                  string    `json:"bondStructuringAdvisor"`
	EntitiesAdditionalInfo                  string    `json:"entitiesAdditionalInfo"`
	AuthorizedRepresentativeName            string    `json:"authorizedRepresentativeName"`
	AuthorizedRepresentativeTitleOrPosition string    `json:"authorizedRepresentativeTitleOrPosition"`
	AuthorizedRepresentativeEmail           string    `json:"authorizedRepresentativeEmail"`
	AcceptTokenizationTermsAndAgreement     int       `gorm:"default:0" json:"acceptTokenizationTermsAndAgreement"`
	AttestInformationAccurateAndVerifiable  int       `gorm:"default:0" json:"attestInformationAccurateAndVerifiable"`
	AcknowledgedSuitabilityCriteria         int       `gorm:"default:0" json:"acknowledgedSuitabilityCriteria"`
	MinimumKycTier                          string    `json:"minimumKycTier"`
	InvestorCategory                        string    `json:"investorCategory"`
	WithholdingTaxDisclosure                int       `gorm:"default:0" json:"withholdingTaxDisclosure"`
}

type ConfirmTokenizedAssetJSONInput struct {
	Messages             []string `json:"messages"`
	Transaction          string   `json:"transaction"`
	TransactionSignature string   `json:"transactionSignature"`
	TransactionID        string   `json:"transactionId"`
	NetworkPassPhrase    string   `json:"networkPassPhrase"`
}

type VetTokenizedAssetJSONInput struct {
	ApprovedAssetCustodianID             uint64  `gorm:"not null" json:"approvedAssetCustodianId"`
	ApprovedAssetCustodianFeePercent     float64 `gorm:"not null" json:"ApprovedAssetCustodianFeePercent"`
	ApprovedAssetCustodianFeeFixed       float64 `gorm:"not null" json:"ApprovedAssetCustodianFeeFixed"`
	AssetManagerID                       uint64  `json:"assetManagerId"`
	AssetManagerFeePercent               float64 `json:"assetManagerFeePercent"`
	AssetManagerFeeFixed                 float64 `json:"assetManagerFeeFixed"`
	AssetIssuingHouseID                  uint64  `json:"assetIssuingHouseId"`
	AssetIssuingHouseFeePercent          float64 `json:"assetIssuingHouseFeePercent"`
	AssetIssuingHouseFeeFixed            float64 `json:"assetIssuingHouseFeeFixed"`
	LegalAndProfesionalPartnerID         uint64  `json:"legalAndProfesionalPartnerId"`
	LegalAndProfesionalPartnerFeePercent float64 `json:"legalAndProfesionalPartnerFeePercent"`
	LegalAndProfesionalPartnerFeeFixed   float64 `json:"legalAndProfesionalPartnerFeeFixed"`
	RatingAgencyID                       uint64  `json:"ratingAgencyId"`
	RatingAgencyFeePercent               float64 `json:"ratingAgencyFeePercent"`
	RatingAgencyFeeFixed                 float64 `json:"ratingAgencyFeeFixed"`
	TrusteeID                            uint64  `json:"trusteeId"`
	TrusteeFeePercent                    float64 `json:"trusteeFeePercent"`
	TrusteeFeeFixed                      float64 `json:"trusteeFeeFixed"`
	CountryCode                          string  `json:"CountryCode"`
	ProceedPayoutCurrency                string  `json:"proceedPayoutCurrency"`
	AssetQuoteCurrency                   string  `gorm:"default:'CNGN'" json:"assetQuoteCurrency"`
	ExcludeSecFee                        int     `gorm:"default:0" json:"excludeSecFee"`
}

type TokenizedAssetJSON struct {
	ID                                           string                          `json:"id"`
	CreatedAt                                    time.Time                       `json:"createdAt"`
	UpdatedAt                                    time.Time                       `json:"updatedAt"`
	InitiatorUsername                            string                          `gorm:"size:50;not null" json:"initiatorUsername"`
	AssetSector                                  string                          `json:"assetSector"`
	AssetSubSector                               string                          `json:"assetSubSector"`
	AssetType                                    string                          `json:"assetType"`
	AssetName                                    string                          `json:"assetName"`
	AssetWebsite                                 string                          `json:"assetWebsite"`
	ApprovedAssetCustodianID                     uint64                          `gorm:"not null" json:"approvedAssetCustodianId"`
	ApprovedAssetCustodian                       ApprovedAssetCustodian          `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"approvedAssetCustodianInfo"`
	AssetIssuingHouseID                          uint64                          `gorm:"not null" json:"assetIssuingHouseId"`
	AssetIssuingHouse                            AssetIssuingHouse               `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"assetIssuingHouseInfo"`
	LegalAndProfesionalPartnerID                 uint64                          `gorm:"not null" json:"legalAndProfesionalPartnerId"`
	LegalAndProfesionalPartner                   LegalAndProfesionalPartner      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"legalAndProfesionalPartnerInfo"`
	RatingAgencyID                               uint64                          `gorm:"not null" json:"ratingAgencyId"`
	RatingAgency                                 RatingAgency                    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"ratingAgencyInfo"`
	TrusteeID                                    uint64                          `gorm:"not null" json:"trusteeId"`
	Trustee                                      Trustee                         `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"trusteeInfo"`
	OfferingType                                 string                          `gorm:"default:'PRIVATE'" json:"offeringType"` //PRIVATE, PUBLIC
	ClosedGroupID                                string                          `gorm:"null" json:"closedGroupId"`
	ClosedGroup                                  ClosedGroup                     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"closedGroupInfo"`
	SecApproval                                  int                             `gorm:"default:0" json:"secApproval"`
	SecApprovalIdNumber                          string                          `json:"secApprovalIdNumber"`
	IssuingWalletPublicKey                       string                          `gorm:"size:60" json:"issuingWalletPublicKey"`
	IssuingWalletAlias                           string                          `gorm:"size:60" json:"issuingWalletAlias"`
	MarketMakingWallet                           string                          `json:"marketMakingWallet"`
	AssetDescription                             string                          `json:"assetDescription"`
	AssetCountryLocation                         string                          `json:"assetCountryLocation"`
	AssetPhysicalAddress                         string                          `json:"assetPhysicalAddress"`
	AssetLongitude                               string                          `json:"assetLongitude"`
	AssetLatitude                                string                          `json:"assetLatitude"`
	OwnershipType                                string                          `json:"ownershipType"` //DIRECT,THIRD-PARTY
	OwnershipKind                                string                          `json:"ownershipKind"` //INDIVIDUAL,CORPORATE,
	InitialOwnerPreferredWalletAddress           string                          `json:"initialOwnerPreferredWalletAddress"`
	AssetOwnerName                               string                          `json:"assetOwnerName"`
	AssetOwnerRetainedOrContributedValue         float64                         `json:"assetOwnerRetainedOrContributedValue"`
	AssetOwnerAddress                            string                          `json:"assetOwnerAddress"`
	AssetManagerID                               uint64                          `json:"assetManagerId"`
	AssetManager                                 AssetManager                    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"assetManagerInfo"`
	AssetQuoteCurrency                           string                          `gorm:"default:'NGN'" json:"assetQuoteCurrency"`
	AssetCurrentValue                            float64                         `gorm:"default:0" json:"assetCurrentValue"`
	AssetMscCostOutisdeOfValuation               float64                         `gorm:"default:0" json:"assetMscCostOutisdeOfValuation"`
	ValueOfTokenizedAsset                        float64                         `gorm:"default:0" json:"valueOfTokenizedAsset"`
	ProtectionMethods                            string                          `json:"protectionMethods"` //csv format
	InsuranceCompanyName                         string                          `json:"insuranceCompanyName"`
	InsurancePolicyNumber                        string                          `json:"insurancePolicyNumber"`
	InsurancePolicyHolder                        string                          `json:"insurancePolicyHolder"`
	PercentageValueOfInsurance                   float64                         `gorm:"default:0" json:"percentageValueOfInsurance"`
	IsFreeFromLiensAndEncumbrances               int                             `gorm:"default:0" json:"IsFreeFromLiensAndEncumbrances"`
	AssetAlreadyExists                           int                             `gorm:"default:0" json:"assetAlreadyExists"`
	VettingStatus                                int                             `gorm:"default:0" json:"vettingStatus"`
	DueDiligenceFail                             int                             `gorm:"default:0" json:"dueDiligenceFail"` //0=False(success/in-progress), 1= true (failed).
	DueDiligenceFailureReason                    string                          `json:"dueDiligenceFailureReason"`
	AssetTokenizationDocuments                   []AssetTokenizationDocument     `json:"AssetTokenizationDocuments"`
	ProofOfPaymentDocuments                      []TokenizationFeeProofOfPayment `json:"ProofOfPaymentDocuments"`
	AssetCode                                    string                          `json:"assetCode"`
	AssetLogo                                    string                          `json:"assetLogo"`
	NumberOfTokenToBeIssued                      float64                         `json:"numberOfTokenToBeIssued"`
	MaxNumberOfTokenAvailableForSale             float64                         `gorm:"default:0" json:"maxNumberOfTokenAvailableForSale"`
	FeeInAsset                                   float64                         `gorm:"default:0" json:"feeInAsset"`
	FeeInAssetPercent                            float64                         `gorm:"default:0" json:"feeInAssetPercent"`
	FeeInFiat                                    float64                         `gorm:"default:0" json:"feeInFiat"`
	NumberOfTokenToBeSold                        float64                         `gorm:"default:0" json:"numberOfTokenToBeSold"`
	TotalTokenHeldByManager                      float64                         `gorm:"default:0" json:"totalTokenHeldByManager"`
	WalletToHoldAssetsNotForSale                 string                          `json:"walletToHoldAssetsNotForSale"`
	PricePerToken                                float64                         `gorm:"default:0" json:"pricePerToken"`
	SalesStart                                   time.Time                       `json:"salesStart"`
	SalesEnd                                     time.Time                       `json:"salesEnd"`
	CapOnPurchase                                int                             `gorm:"default:0" json:"capOnPurchase"`
	CapQuantity                                  float64                         `gorm:"default:0" json:"capQuantity"`
	CapAmountInFiat                              float64                         `gorm:"default:0" json:"capAmountInFiat"`
	CapDurationInDays                            int                             `gorm:"default:0" json:"capDurationInDays"`
	ProceedCycle                                 string                          `gorm:"size:50" json:"proceedCycle"`
	TokenizationFeeID                            uint64                          `json:"tokenizationFeeId"`
	TokenizationFee                              TokenizationFee                 `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"tokenizationFee"`
	CountryConfig                                Country                         `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"countryConfig"`
	ExcludeSecFee                                int                             `gorm:"default:0" json:"excludeSecFee"`
	SECTokenizationFeePercent                    float64                         `gorm:"default:0" json:"SECTokenizationFeePercent"`
	SECTokenizationFeeFixed                      float64                         `gorm:"default:0" json:"SECTokenizationFeeFixed"`
	SECTokenizationFeeValue                      float64                         `gorm:"default:0" json:"SECTokenizationFeeValue"`
	CustodianFeePercent                          float64                         `gorm:"default:0" json:"custodianFeePercent"`
	CustodianFeeFixed                            float64                         `gorm:"default:0" json:"custodianFeeFixed"`
	CustodianFeeValue                            float64                         `gorm:"default:0" json:"custodianFeeValue"`
	AssetManagerFeeValue                         float64                         `gorm:"default:0" json:"assetManagerFeeValue"`
	AssetManagerFeePercent                       float64                         `gorm:"default:0" json:"assetManagerFeePercent"`
	AssetManagerFeeFixed                         float64                         `gorm:"default:0" json:"assetManagerFeeFixed"`
	IssuingHouseFeeValue                         float64                         `gorm:"default:0" json:"issuingHouseFeeValue"`
	IssuingHouseFeePercent                       float64                         `gorm:"default:0" json:"issuingHouseFeePercent"`
	IssuingHouseFeeFixed                         float64                         `gorm:"default:0" json:"issuingHouseFeeFixed"`
	LegalAndProfessionalFeePercent               float64                         `gorm:"default:0" json:"legalAndProfessionalFeePercent"`
	LegalAndProfessionalFeeFixed                 float64                         `gorm:"default:0" json:"legalAndProfessionalFeeFixed"`
	LegalAndProfessionalFeeValue                 float64                         `gorm:"default:0" json:"legalAndProfessionalFeeValue"`
	RatingAgencyFeePercent                       float64                         `gorm:"default:0" json:"ratingAgencyFeePercent"`
	RatingAgencyFeeFixed                         float64                         `gorm:"default:0" json:"ratingAgencyFeeFixed"`
	RatingAgencyFeeValue                         float64                         `gorm:"default:0" json:"ratingAgencyFeeValue"`
	TrusteeFeePercent                            float64                         `gorm:"default:0" json:"trusteeFeePercent"`
	TrusteeFeeFixed                              float64                         `gorm:"default:0" json:"trusteeFeeFixed"`
	TrusteeFeeValue                              float64                         `gorm:"default:0" json:"trusteeFeeValue"`
	VATPercent                                   float64                         `gorm:"default:0" json:"vatPercent"`
	VATValue                                     float64                         `gorm:"default:0" json:"vatValue"`
	VATInAsset                                   float64                         `gorm:"default:0" json:"vatInAsset"`
	ProceedPayoutCurrency                        string                          `json:"proceedPayoutCurrency"`
	ProceedPayoutType                            int                             `gorm:"default:0" json:"proceedPayoutType"` // FIAT=1, CRYPTO=0
	ExemptedCountries                            string                          `json:"exemptedCountries"`
	HasAdditionalKYCRequirements                 int                             `gorm:"default:0" json:"hasAdditionalKYCRequirements"`
	AdditionalKYCRequirements                    string                          `json:"additionalKYCRequirements"`
	InvestorAccreditationRequired                int                             `gorm:"default:0" json:"investorAccreditationRequired"`
	AssetTokenizationStatus                      int                             `gorm:"default:0" json:"assetTokenizationStatus"`
	LastUpdatedBy                                string                          `gorm:"null" json:"lastUpdatedBy"`
	AgreeTransferTitleToCustodian                int                             `gorm:"default:0" json:"agreeTransferTitleToCustodian"`
	ContractualProtectionRevGuarantees           int                             `gorm:"default:0" json:"contractualProtectionRevGuarantees"`
	ContractualProtectionPerfBond                int                             `gorm:"default:0" json:"contractualProtectionPerfBond"`
	ContractualProtectionSLA                     int                             `gorm:"default:0" json:"contractualProtectionSLA"`
	RiskSharingMechanismPPPs                     int                             `gorm:"default:0" json:"riskSharingMechanismPPPs"`
	RiskSharingMechanismHedgeInstruments         int                             `gorm:"default:0" json:"riskSharingMechanismHedgeInstruments"`
	RiskSharingMechanismCompletionGuarantees     int                             `gorm:"default:0" json:"riskSharingMechanismCompletionGuarantees"`
	IndependentMonitoringList                    string                          `gorm:"null" json:"independentMonitoringList"`
	ESGSafeguardsSusCerts                        int                             `gorm:"default:0" json:"eSGSafeguardsSusCerts"`
	ESGSafeguardsCommEngPlans                    int                             `gorm:"default:0" json:"eSGSafeguardsCommEngPlans"`
	SecurityMeasuresAccessControl                int                             `gorm:"default:0" json:"securityMeasuresAccessControl"`
	SecurityMeasuresSurveilanceSystems           int                             `gorm:"default:0" json:"securityMeasuresSurveilanceSystems"`
	SecurityMeasuresOnSiteSecurityPersonnel      int                             `gorm:"default:0" json:"securityMeasuresOnSiteSecurityPersonnel"`
	SecurityMeasuresPerimeterSecurity            int                             `gorm:"default:0" json:"securityMeasuresPerimeterSecurity"`
	SecurityMeasuresCriticalInfraProtections     int                             `gorm:"default:0" json:"securityMeasuresCriticalInfraProtections"`
	OtherAssetProtection                         string                          `gorm:"null" json:"otherAssetProtection"`
	LegalAdvisor                                 string                          `gorm:"null" json:"legalAdvisor"`
	FinancialAdvisor                             string                          `gorm:"null" json:"financialAdvisor"`
	UndertakingNoLien                            int                             `gorm:"default:0" json:"undertakingNoLien"`
	UndertakingNotCollateral                     int                             `gorm:"default:0" json:"undertakingNotCollateral"`
	UndertakingNoClaims                          int                             `gorm:"default:0" json:"undertakingNoClaims"`
	UndertakingNoForeclosure                     int                             `gorm:"default:0" json:"undertakingNoForeclosure"`
	ComplianceNoViolation                        int                             `gorm:"default:0" json:"complianceNoViolation"`
	ComplianceAllPermits                         int                             `gorm:"default:0" json:"complianceAllPermits"`
	OutstandingFinancialRespNoDebts              int                             `gorm:"default:0" json:"outstandingFinancialRespNoDebts"`
	OutstandingFinancialRespNoHiddenLiabilities  int                             `gorm:"default:0" json:"outstandingFinancialRespNoHiddenLiabilities"`
	RiskManagementFullyInsured                   int                             `gorm:"default:0" json:"riskManagementFullyInsured"`
	RiskManagementDeclaredValue                  int                             `gorm:"default:0" json:"riskManagementDeclaredValue"`
	PhysicalConditionSound                       int                             `gorm:"default:0" json:"physicalConditionSound"`
	PhysicalConditionNolease                     int                             `gorm:"default:0" json:"physicalConditionNolease"`
	PhysicalConditionNoUndisclosedEasements      int                             `gorm:"default:0" json:"physicalConditionNoUndisclosedEasements"`
	MintingInitators                             string                          `gorm:"null" json:"mintingInitators"` //CSV of approvers
	MintingApprovers                             string                          `gorm:"null" json:"mintingApprovers"` //csv of initators
	BankID                                       uint64                          `gorm:"null" json:"bankId"`
	Bank                                         Bank                            `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"bankInfo"`
	AccountNumber                                string                          `gorm:"null" json:"accountNumber"`
	BeneficiaryName                              string                          `gorm:"null" json:"beneficiaryName"`
	TokenizationApplicationFee                   float64                         `json:"tokenizationApplicationFee"`
	TokenizationApplicationFeeAsset              string                          `json:"tokenizationApplicationFeeAsset"`
	ProjectStrategicObjectives                   string                          `json:"projectStrategicObjectives"`
	ProjectDevelopmentTimeline                   string                          `json:"projectDevelopmentTimeline"`
	ProjectKeyMilestoneAndDates                  string                          `json:"projectKeyMilestoneAndDates"`
	ProjectScope                                 string                          `json:"projectScope"`
	ProjectEconomicBenefits                      string                          `json:"projectEconomicBenefits"`
	ProjectExpectedNoOfJobs                      int                             `json:"projectExpectedNoOfJobs"`
	ProjectIntendedSocialBenefits                string                          `json:"projectIntendedSocialBenefits"`
	ProjectTechnicalPartners                     string                          `json:"projectTechnicalPartners"`
	ProjectFinancialPartners                     string                          `json:"projectFinancialPartners"`
	EstimatedProjectIRR                          float64                         `json:"estimatedProjectIRR"`
	EstimatedProjectROI                          float64                         `json:"estimatedProjectROI"`
	EstimatedProjectNPV                          float64                         `json:"estimatedProjectNPV"`
	EstimatedProjectPaybackPeriodsInMonths       int                             `json:"estimatedProjectPaybackPeriodsInMonths"`
	KeyAssumptionsList                           string                          `json:"keyAssumptionsList"`
	ProjectIdentifiedLegalRisks                  string                          `json:"projectIdentifiedLegalRisks"`
	ProjectIdentifiedRegulatoryRisks             string                          `json:"projectIdentifiedRegulatoryRisks"`
	ProjectIdentifiedOperationalOrExecutionRisks string                          `json:"projectIdentifiedOperationalOrExecutionRisks"`
	ProjectIdentifiedMarketRisks                 string                          `json:"projectIdentifiedMarketRisks"`
	ProjectIdentifiedOtherRelevantRisks          string                          `json:"projectIdentifiedOtherRelevantRisks"`
	NumberOfExpressedInterests                   int64                           `json:"numberOfExpressedInterests"`
	NumberOfSubscribers                          int64                           `json:"numberOfSubscribers"`
	QuantityOfTokensSold                         float64                         `json:"quantityOfTokensSold"`
	QuantityOfTokensSoldInFiat                   float64                         `json:"quantityOfTokensSoldInFiat"`
	PurchaseCommitments                          float64                         `json:"purchaseCommitments"`
	DeepLink                                     string                          `json:"deepLink"`
	IsinOrSerialNumber                           string                          `json:"isinOrSerialNumber"` /// begin the mutual fund fields
	InstrumentName                               string                          `json:"instrumentName"`
	InstrumentType                               string                          `json:"instrumentType"`
	TotalIssueSize                               float64                         `json:"totalIssueSize"`
	IssueDate                                    time.Time                       `json:"issueDate"`
	MaturityDate                                 time.Time                       `json:"maturityDate"`
	FaceValuePerUnit                             float64                         `json:"faceValuePerUnit"`
	CouponOrInterestRateType                     string                          `json:"couponOrInterestRateType"`
	CouponOrInterestRate                         float64                         `json:"couponOrInterestRate"`
	ReferenceIndex                               string                          `json:"referenceIndex"`
	ExpectedYield                                float64                         `json:"expectedYield"`
	EarlyRedemptionOptionInvestor                string                          `json:"earlyRedemptionOptionInvestor"`
	MinimumInvestmentAmount                      float64                         `json:"minimumInvestmentAmount"`
	TaxTreatmentTokenHolders                     string                          `json:"taxTreatmentTokenHolders"`
	PaymentStructureToTokenHolders               string                          `json:"paymentStructureToTokenHolders"`
	RedemptionMethod                             string                          `json:"redemptionMethod"`
	PaymentStructure                             string                          `json:"paymentStructure"`
	RepaymentMethod                              string                          `json:"repaymentMethod"`
	PaymentCycle                                 string                          `json:"paymentCycle"`
	WeightedAverageLife                          float64                         `json:"weightedAverageLife"`
	UnderlyingAssetPoolSize                      float64                         `json:"underlyingAssetPoolSize"`
	PoolComposition                              string                          `json:"poolComposition"`
	CreditEnhancementMethod                      string                          `json:"creditEnhancementMethod"`
	SummaryOfUseOfProceeds                       string                          `json:"summaryOfUseOfProceeds"`
	CreditRatingIfAny                            string                          `json:"creditRatingIfAny"`
	IssuerName                                   string                          `json:"issuerName"`
	IssuerType                                   string                          `json:"issuerType"`
	IssuerContactPerson                          string                          `json:"issuerContactPerson"`
	ContactEmail                                 string                          `json:"contactEmail"`
	ContactPhoneNumber                           string                          `json:"contactPhoneNumber"`
	BriefCompanyOverview                         string                          `json:"briefCompanyOverview"`
	MortgageOriginators                          string                          `json:"mortgageOriginators"`
	Servicer                                     string                          `json:"servicer"`
	InstrumentTrustee                            string                          `json:"instrumentTrustee"`
	InstrumentCustodians                         string                          `json:"instrumentCustodian"`
	InstrumentLegalAdvisor                       string                          `json:"instrumentLegalAdvisor"`
	SpecialPurposeVehicle                        string                          `json:"specialPurposeVehicle"`
	InstrumentAssetManagerOrAdministrator        string                          `json:"instrumentAssetManagerOrAdministrator"`
	UnderwriterIfAny                             string                          `json:"underwriterIfAny"`
	InstrumentCreditRatingAgency                 string                          `json:"instrumentCreditRatingAgency"`
	AuditorOrVerifier                            string                          `json:"auditorOrVerifier"`
	CreditRiskAssessment                         string                          `json:"creditRiskAssessment"`
	CreditRating                                 string                          `json:"creditRating"`
	PrepaymentRisk                               string                          `json:"prepaymentRisk"`
	InterestRateRisk                             string                          `json:"interestRateRisk"`
	StructuralComplexityRisk                     string                          `json:"structuralComplexityRisk"`
	LegalOrRegulatoryRisk                        string                          `json:"legalOrRegulatoryRisk"`
	OperationalRisk                              string                          `json:"operationalRisk"`
	MarketRisk                                   string                          `json:"marketRisk"`
	EsgRisk                                      string                          `json:"esgRisk"`
	MitigationMeasures                           string                          `json:"mitigationMeasures"`
	IssuingAuthority                             string                          `json:"issuingAuthority"`
	RegulatoryApprovalId                         string                          `json:"regulatoryApprovalId"`
	LicenseApprovalReferenceNumber               string                          `json:"licenseApprovalReferenceNumber"`
	ListingStatus                                string                          `json:"listingStatus"`
	CurrencyOfIssuance                           string                          `json:"currencyOfIssuance"`
	CouponRateType                               string                          `json:"couponRateType"`
	CouponRate                                   float64                         `json:"couponRate"`
	SpreadOrMargin                               float64                         `json:"spreadOrMargin"`
	ResetFrequency                               string                          `json:"resetFrequency"`
	CouponPaymentFrequency                       string                          `json:"couponPaymentFrequency"`
	RedemptionStructure                          string                          `json:"redemptionStructure"`
	EarlyRedemptionOption                        string                          `json:"earlyRedemptionOption"`
	EarlyRedemptionPenalty                       string                          `json:"earlyRedemptionPenalty"`
	TaxTreatment                                 string                          `json:"taxTreatment"`
	NavOrMarketValueUpdates                      string                          `json:"navOrMarketValueUpdates"`
	ImpactMetrics                                string                          `json:"impactMetrics"`
	LegalBacking                                 string                          `json:"legalBacking"`
	DefaultHistory                               string                          `json:"defaultHistory"`
	RiskFactorsSummary                           string                          `json:"riskFactorsSummary"`
	PayingAgent                                  string                          `json:"payingAgent"`
	Auditor                                      string                          `json:"auditor"`
	RegistrarOrCSCSAgent                         string                          `json:"registrarOrCSCSAgent"`
	FundStructure                                string                          `json:"fundStructure"`
	AssetManagementCompanyName                   string                          `json:"assetManagementCompanyName"`
	FundManagers                                 string                          `json:"fundManagers"`
	RegulatoryLicenseNumber                      string                          `json:"regulatoryLicenseNumber"`
	FundLaunchDate                               time.Time                       `json:"fundLaunchDate"`
	TotalExpenseRatio                            float64                         `json:"totalExpenseRatio"`
	ExitLoadRedemptionFee                        float64                         `json:"exitLoadRedemptionFee"`
	Tenure                                       int                             `json:"tenure"`
	InitialNetAssetValue                         float64                         `json:"initialNetAssetValue"`
	NavUpdateFrequency                           string                          `json:"navUpdateFrequency"`
	NavCalculationMethod                         string                          `json:"navCalculationMethod"`
	RedemptionRules                              string                          `json:"redemptionRules"`
	LockInPeriod                                 int                             `json:"lockInPeriod"`
	EntryLoad                                    float64                         `json:"entryLoad"`
	PerformanceFee                               float64                         `json:"performanceFee"`
	DividendPolicy                               string                          `json:"dividendPolicy"`
	LiquidityProfile                             string                          `json:"liquidityProfile"`
	DistributionFrequency                        string                          `json:"distributionFrequency"`
	DistributionMethod                           string                          `json:"distributionMethod"`
	BenchmarkComparisonMethod                    string                          `json:"benchmarkComparisonMethod"`
	FeeBreakdownSummary                          string                          `json:"feeBreakdownSummary"`
	InvestmentObjective                          string                          `json:"investmentObjective"`
	EquityStrategy                               string                          `json:"equityStrategy"`
	MarketCapitalizationFocus                    string                          `json:"marketCapitalizationFocus"`
	BenchmarkIndex                               string                          `json:"benchmarkIndex"`
	SectorExposureLimits                         string                          `json:"sectorExposureLimits"`
	TopHoldings                                  string                          `json:"topHoldings"`
	GeographicExposure                           string                          `json:"geographicExposure"`
	RiskProfile                                  string                          `json:"riskProfile"`
	VolatilityEstimate                           float64                         `json:"volatilityEstimate"`
	DividendYield                                float64                         `json:"dividendYield"`
	TrusteeName                                  string                          `json:"trusteeName"`
	FundAdministrator                            string                          `json:"fundAdministrator"`
	InvestmentCommitteeMembers                   string                          `json:"investmentCommitteeMembers"`
	IsinOrSecFundCode                            string                          `json:"isinOrSecFundCode"`
	ExitLoadOrRedemptionFee                      float64                         `json:"exitLoadOrRedemptionFee"`
	FundRiskRating                               string                          `json:"fundRiskRating"`
	AssetAllocation                              string                          `json:"assetAllocation"`
	AverageMaturity                              float64                         `json:"averageMaturity"`
	YieldToMaturity                              float64                         `json:"yieldToMaturity"`
	CreditRatingProfile                          string                          `json:"creditRatingProfile"`
	LockInPeriodPortfolio                        int                             `json:"lockInPeriodPortfolio"`
	PerformanceFeeIfAny                          float64                         `json:"performanceFeeIfAny"`
	TopEquityHoldings                            string                          `json:"topEquityHoldings"`
	TargetAllocation                             string                          `json:"targetAllocation"`
	AllowedAllocationRange                       string                          `json:"allowedAllocationRange"`
	AssetClassesIncluded                         string                          `json:"assetClassesIncluded"`
	RebalancingFrequency                         string                          `json:"rebalancingFrequency"`
	BenchmarkIndexComposite                      string                          `json:"benchmarkIndexComposite"`
	TopEquityHoldingsList                        string                          `json:"topEquityHoldingsList"`
	TopDebtHoldings                              string                          `json:"topDebtHoldings"`
	CreditRatingDistribution                     string                          `json:"creditRatingDistribution"`
	AverageMaturityHybrid                        float64                         `json:"averageMaturityHybrid"`
	YieldToMaturityHybrid                        float64                         `json:"yieldToMaturityHybrid"`
	TitleOfIssuance                              string                          `json:"titleOfIssuance"`
	TypeOfCommercialPaper                        string                          `json:"typeOfCommercialPaper"`
	PricingYield                                 float64                         `json:"pricingYield"`
	UseOfProceeds                                string                          `json:"useOfProceeds"`
	Ranking                                      string                          `json:"ranking"`
	BackingSecurity                              string                          `json:"backingSecurity"`
	IssuerRegistrationNumber                     string                          `json:"issuerRegistrationNumber"`
	IncorporationDate                            time.Time                       `json:"incorporationDate"`
	RcNumber                                     string                          `json:"rcNumber"`
	TaxIdNumber                                  string                          `json:"taxIdNumber"`
	OfficeAddress                                string                          `json:"officeAddress"`
	Rating                                       string                          `json:"rating"`
	PartiesInvolvedIssuer                        string                          `json:"partiesInvolvedIssuer"`
	PartiesInvolvedArranger                      string                          `json:"partiesInvolvedArranger"`
	PartiesInvolvedLegalAdviser                  string                          `json:"partiesInvolvedLegalAdviser"`
	PartiesInvolvedAuditor                       string                          `json:"partiesInvolvedAuditor"`
	PartiesInvolvedRatingAgency                  string                          `json:"partiesInvolvedRatingAgency"`
	PartiesInvolvedCustodian                     string                          `json:"partiesInvolvedCustodian"`
	PartiesInvolvedTrustee                       string                          `json:"partiesInvolvedTrustee"`
	PartiesInvolvedAuditorVerifier               string                          `json:"partiesInvolvedAuditorVerifier"`
	SecurityRiskLegalBacking                     string                          `json:"securityRiskLegalBacking"`
	SecurityRiskCollateral                       string                          `json:"securityRiskCollateral"`
	SecurityRiskDefaultHistory                   string                          `json:"securityRiskDefaultHistory"`
	SecurityRiskCreditRating                     string                          `json:"securityRiskCreditRating"`
	SecurityRiskRiskFactorsSummary               string                          `json:"securityRiskRiskFactorsSummary"`
	SecurityRiskBusinessRisk                     string                          `json:"securityRiskBusinessRisk"`
	SecurityRiskDefaultRisk                      string                          `json:"securityRiskDefaultRisk"`
	SecurityRiskLiquidityRisk                    string                          `json:"securityRiskLiquidityRisk"`
	SecurityRiskRegulatoryRisk                   string                          `json:"securityRiskRegulatoryRisk"`
	SecurityRiskMarketRisk                       string                          `json:"securityRiskMarketRisk"`
	SecurityRiskOperationalRisk                  string                          `json:"securityRiskOperationalRisk"`
	SecurityRiskMitigationMeasures               string                          `json:"securityRiskMitigationMeasures"`
	CommodityType                                string                          `json:"commodityType"`
	CommodityDescription                         string                          `json:"commodityDescription"`
	Quantity                                     int                             `json:"quantity"`
	QualityGrade                                 string                          `json:"qualityGrade"`
	IssuerContactInfo                            string                          `json:"issuerContactInfo"`
	WarehouseName                                string                          `json:"warehouseName"`
	WarehouseOperatorName                        string                          `json:"warehouseOperatorName"`
	WarehouseLicenseNumber                       string                          `json:"warehouseLicenseNumber"`
	WarehouseLocation                            string                          `json:"warehouseLocation"`
	WrNumber                                     string                          `json:"wrNumber"`
	WrIssueDate                                  time.Time                       `json:"wrIssueDate"`
	WrExpiryDate                                 time.Time                       `json:"wrExpiryDate"`
	WrSystemRegistration                         string                          `json:"wrSystemRegistration"`
	WrRegistrationNumber                         string                          `json:"wrRegistrationNumber"`
	WrVerifier                                   string                          `json:"wrVerifier"`
	StorageCondition                             string                          `json:"storageCondition"`
	WarehouseAccreditationBody                   string                          `json:"warehouseAccreditationBody"`
	MinimumPurchaseAmount                        float64                         `json:"minimumPurchaseAmount"`
	AutoRollover                                 string                          `json:"autoRollover"`
	CurrentBeneficialOwner                       string                          `json:"currentBeneficialOwner"`
	WrCustodianName                              string                          `json:"wrCustodianName"`
	OwnershipRightsRepresented                   string                          `json:"ownershipRightsRepresented"`
	TrusteeOrThirdPartyOversight                 string                          `json:"trusteeOrThirdPartyOversight"`
	LienOrEncumbrances                           string                          `json:"lienOrEncumbrances"`
	AssetValuation                               float64                         `json:"assetValuation"`
	ValuationDate                                time.Time                       `json:"valuationDate"`
	ValuationMethodology                         string                          `json:"valuationMethodology"`
	TokenizationObjective                        string                          `json:"tokenizationObjective"`
	HoldingPeriod                                int                             `json:"holdingPeriod"`
	RedemptionMechanism                          string                          `json:"redemptionMechanism"`
	PartiesInvolvedUnderwriter                   string                          `json:"partiesInvolvedUnderwriter"`
	PartiesInvolvedAssetManager                  string                          `json:"partiesInvolvedAssetManager"`
	PartiesInvolvedLegalAdvisor                  string                          `json:"partiesInvolvedLegalAdvisor"`
	PartiesInvolvedRegulator                     string                          `json:"partiesInvolvedRegulator"`
	RisksMarketRisk                              string                          `json:"risksMarketRisk"`
	RisksStorageRisk                             string                          `json:"risksStorageRisk"`
	RisksTitleRisk                               string                          `json:"risksTitleRisk"`
	RisksFraudRisk                               string                          `json:"risksFraudRisk"`
	RisksInsuranceRisk                           string                          `json:"risksInsuranceRisk"`
	RisksOperationalRisk                         string                          `json:"risksOperationalRisk"`
	RisksRegulatoryRisk                          string                          `json:"risksRegulatoryRisk"`
	RisksLiquidityRisk                           string                          `json:"risksLiquidityRisk"`
	RisksForceMajeureRisk                        string                          `json:"risksForceMajeureRisk"`
	RisksEarlyRedemptionRisk                     string                          `json:"risksEarlyRedemptionRisk"`
	RisksMitigationMeasures                      string                          `json:"risksMitigationMeasures"`
	RisksInsuranceCoverageSummary                string                          `json:"risksInsuranceCoverageSummary"`
	RisksInsuranceProvider                       string                          `json:"risksInsuranceProvider"`
	RisksCoverageValue                           float64                         `json:"risksCoverageValue"`
	QuanlityStandard                             string                          `json:"quanlityStandard"`
	IssuerContactInformation                     string                          `json:"issuerContactInformation"`
	VaultCustodianName                           string                          `json:"vaultCustodianName"`
	VaultOperator                                string                          `json:"vaultOperator"`
	VaultLicenseNumber                           string                          `json:"vaultLicenseNumber"`
	VaultLocation                                string                          `json:"vaultLocation"`
	Number                                       string                          `json:"number"`
	IssuerDate                                   time.Time                       `json:"issuerDate"`
	ExpiryDate                                   time.Time                       `json:"expiryDate"`
	RegistryRecord                               string                          `json:"registryRecord"`
	Verifier                                     string                          `json:"verifier"`
	StorageConditions                            string                          `json:"storageConditions"`
	VaultAccreditationBody                       string                          `json:"vaultAccreditationBody"`
	OwnershipLegalHolder                         string                          `json:"ownershipLegalHolder"`
	OwnershipCustodianName                       string                          `json:"ownershipCustodianName"`
	OwnershipTrustee                             string                          `json:"ownershipTrustee"`
	OwnershipLienOrEncumbrances                  string                          `json:"ownershipLienOrEncumbrances"`
	ValuationAssetValuation                      string                          `json:"valuationAssetValuation"`
	HoldingLockinPeriod                          int                             `json:"holdingLockinPeriod"`
	InsuranceMarketRisk                          string                          `json:"insuranceMarketRisk"`
	InsuranceStorageRisk                         string                          `json:"insuranceStorageRisk"`
	InsuranceTitleRisk                           string                          `json:"insuranceTitleRisk"`
	InsuranceFraudRisk                           string                          `json:"insuranceFraudRisk"`
	InsuranceInsuranceRisk                       string                          `json:"insuranceInsuranceRisk"`
	InsuranceOperationalRisk                     string                          `json:"insuranceOperationalRisk"`
	InsuranceRegulatoryRisk                      string                          `json:"insuranceRegulatoryRisk"`
	InsuranceLiquidityRisk                       string                          `json:"insuranceLiquidityRisk"`
	InsuranceForceMajeureRisk                    string                          `json:"insuranceForceMajeureRisk"`
	InsuranceEarlyRedemptionRisk                 string                          `json:"insuranceEarlyRedemptionRisk"`
	InsuranceMitigationMeasures                  string                          `json:"insuranceMitigationMeasures"`
	InsuranceInsuranceCoverageSummary            string                          `json:"insuranceInsuranceCoverageSummary"`
	InsuranceInsuranceProvider                   string                          `json:"insuranceInsuranceProvider"`
	InsuranceCoverageValue                       string                          `json:"insuranceCoverageValue"`
	IssuerRegistrationNo                         string                          `json:"issuerRegistrationNo"`
	SectorAndIndustry                            string                          `json:"sectorAndIndustry"`
	LicenseOrPermitNumber                        string                          `json:"licenseOrPermitNumber"`
	IssuerAdditionalInfo                         string                          `json:"issuerAdditionalInfo"`
	IsinSerialNumber                             string                          `json:"isinSerialNumber"`
	EsgOrImpactMetrics                           string                          `json:"esgOrImpactMetrics"`
	InstrumentAdditionalInfo                     string                          `json:"instrumentAdditionalInfo"`
	SecurityType                                 string                          `json:"securityType"`
	CollateralDescription                        string                          `json:"collateralDescription"`
	CovenantSummary                              string                          `json:"covenantSummary"`
	CovenantTestingFrequency                     string                          `json:"covenantTestingFrequency"`
	EventOfDefaultClauses                        string                          `json:"eventOfDefaultClauses"`
	LegalEnforcementMechanism                    string                          `json:"legalEnforcementMechanism"`
	Guarantee                                    string                          `json:"guarantee"`
	RecoveryEstimate                             float64                         `json:"recoveryEstimate"`
	RiskProfileAdditionalInfo                    string                          `json:"riskProfileAdditionalInfo"`
	LicenseNumber                                string                          `json:"licenseNumber"`
	ExitLoadFee                                  float64                         `json:"exitLoadFee"`
	EntryLoadFee                                 float64                         `json:"entryLoadFee"`
	PortfolioLockInPeriod                        int                             `json:"portfolioLockInPeriod"`
	PortfolioPerformanceFee                      float64                         `json:"portfolioPerformanceFee"`
	FundingStructure                             int                             `gorm:"default:0" json:"fundingStructure"` //0=equity, 1= debt, 2= hybrid
	EquityPercentage                             float64                         `gorm:"default:0" json:"equityPercentage"`
	DebtPercentage                               float64                         `gorm:"default:0" json:"debtPercentage"`
	DebtInstrumentType                           string                          `json:"debtInstrumentType"`
	PrincipalPaymentMethod                       string                          `json:"principalPaymentMethod"`
	DebtInstrumentRepaymentSource                string                          `json:"debtInstrumentRepaymentSource"`
	DebtInstrumentGuaranteesOrEnhancements       string                          `json:"debtInstrumentGuaranteesOrEnhancements"`
	DebtInstrumentDefaultAndRecoveryTerms        string                          `json:"debtInstrumentDefaultAndRecoveryTerms"`
	DebtInstrumentRepaymentFrequency             string                          `json:"debtInstrumentRepaymentFrequency"`
	InterestRepaymentFrequency                   string                          `json:"interestRepaymentFrequency"`
	DcsrDetails                                  string                          `json:"dcsrDetails"`
	SinkingFundStructure                         string                          `json:"sinkingFundStructure"`
	CovenantMonitoringAgent                      string                          `json:"covenantMonitoringAgent"`
	RightOfRecourse                              string                          `json:"rightOfRecourse"`
	DebtInstrumentInterestRate                   float64                         `gorm:"default:0" json:"debtInstrumentInterestRate"`
	DcsrRatio                                    float64                         `gorm:"default:0" json:"dcsrRatio"`
	LtvRatio                                     float64                         `gorm:"default:0" json:"ltvRatio"`
	InterestCoverageRatio                        float64                         `gorm:"default:0" json:"interestCoverageRatio"`
	MaximumLeverageRatio                         float64                         `gorm:"default:0" json:"maximumLeverageRatio"`
	GracePeriod                                  int                             `gorm:"default:0" json:"gracePeriod"`
	TrusteeAppointed                             int                             `gorm:"default:0" json:"trusteeAppointed"`
	ReserveFundInPlace                           int                             `gorm:"default:0" json:"reserveFundInPlace"`
	SecurityOrCollateralOffered                  string                          `json:"securityOrCollateralOffered"`
	FundInstrumentType                           string                          `json:"fundInstrumentType"`
	InstrumentRatingAgency                       string                          `json:"instrumentRatingAgency"`
	PortfolioTopHoldings                         string                          `json:"portfolioTopHoldings"`
	CreditRatingAgency                           string                          `json:"CreditRatingAgency"`
	TrusteeRegNumber                             string                          `json:"trusteeRegNumber"`
	CreditEnhancerOrGuarantor                    string                          `json:"creditEnhancerOrGuarantor"`
	BondStructuringAdvisor                       string                          `json:"bondStructuringAdvisor"`
	EntitiesAdditionalInfo                       string                          `json:"entitiesAdditionalInfo"`
	AuthorizedRepresentativeName                 string                          `json:"authorizedRepresentativeName"`
	AuthorizedRepresentativeTitleOrPosition      string                          `json:"authorizedRepresentativeTitleOrPosition"`
	AuthorizedRepresentativeEmail                string                          `json:"authorizedRepresentativeEmail"`
	AcceptTokenizationTermsAndAgreement          int                             `gorm:"default:0" json:"acceptTokenizationTermsAndAgreement"`
	AttestInformationAccurateAndVerifiable       int                             `gorm:"default:0" json:"attestInformationAccurateAndVerifiable"`
	AcknowledgedSuitabilityCriteria              int                             `gorm:"default:0" json:"acknowledgedSuitabilityCriteria"`
	MinimumKycTier                               string                          `json:"minimumKycTier"`
	InvestorCategory                             string                          `json:"investorCategory"`
	WithholdingTaxDisclosure                     int                             `gorm:"default:0" json:"withholdingTaxDisclosure"`
}

type TokenizedAssetSector struct {
	ID                  string `gorm:"primaryKey;size:100" json:"sector"`
	RequirementDocument string `gorm:"" json:"requirementDocument"`
}

type TokenizedAssetSubSector struct {
	ID                     string `gorm:"primaryKey;size:100" json:"subSector"`
	TokenizedAssetSectorID string `gorm:"size:100" json:"assetSectorId"`
}

type TokenizationCurrency struct {
	AssetCode   string `gorm:"primaryKey;size:12" json:"assetCode"`
	AssetIssuer string `gorm:"size:68" json:"assetIssuer"`
	Label       string `gorm:"size:12" json:"label"`
}
type TokenizationPublicAssetAllowedCountryCode struct {
	ID string `gorm:"size:3" json:"id"`
}

type TokenizationStatus struct {
	ID          uint64 `gorm:"" json:"id"`
	Description string `json:"description"`
	// Inactive       int    `gorm:"default:0" json:"-"`
}

type TokenizationFee struct {
	ID                 uint64  `gorm:"" json:"id"`
	FeeFiatPercentage  float64 `json:"feeFiatPercentage"`
	FeeFiatCap         float64 `json:"feeFiatCap"` //Minimum fee
	FeeAssetPercentage float64 `json:"feeAssetPercentage"`
	// FeeAssetCap        float64 `json:"feeAssetCap"` //Minimum fee
	FeeDescription string `json:"feeDescription"`
	CountryCode    string `gorm:"size:2;default:'NG'" json:"countryCode"`
	Inactive       int    `gorm:"default:0" json:"-"`
}

type TokenizationMintingApprover struct {
	ID       uint64 `gorm:"" json:"id"`
	Approver string `gorm:"not null;size:16; index:idx__mintapprover_unique_user, unique" json:"approver"`
}
type TokenizationMintingInitiator struct {
	ID        uint64 `gorm:"" json:"id"`
	Initiator string `gorm:"not null;size:16; index:idx__mintinitiator_unique_user, unique" json:"initiator"`
}

type TokenizationFeePaymentMethod struct {
	ID               string  `gorm:"" json:"id"` //DIGITAL ASSET, FIAT, ETC
	FeeDescription   string  `json:"feeDescription"`
	ExtraDescription *string `json:"extraDescription"`
	Inactive         int     `gorm:"default:0" json:"-"`
}

type TokenizationFeeProofOfPayment struct {
	ID                             uint64 `json:"id"`
	CreatedAt                      time.Time
	TokenizationFeePaymentMethodID string  `gorm:"" json:"tokenizationFeePaymentMethodID"`
	TokenizedAssetID               string  `json:"tokenizedAssetId"`
	TransactionReference           *string `json:"transactionReference"`
	DocumentUrl                    string  `json:"documentUrl"`
}

type TokenizationFeeProofOfPaymentInput struct {
	TokenizationFeePaymentMethodID string `gorm:"" json:"tokenizationFeePaymentMethodID" form:"tokenizationFeePaymentMethodID"`
	TransactionReference           string `json:"transactionReference" form:"transactionReference"`
	// DocumentFile     *multipart.File `form:"documentFile"`
}

type TokenizedAssetType struct {
	ID                        uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	TokenizedAssetSubSectorID string `gorm:"size:100;index:idx_asset_type_unique,unique" json:"assetSubSectorId"`
	AssetType                 string `gorm:"size:100;index:idx_asset_type_unique,unique" json:"assetType"`
}

type ApprovedAssetCustodian struct {
	ID                    uint64  `gorm:"" json:"id"`
	AssetCustodianName    string  `gorm:"size:100" json:"assetCustodianName"`
	AssetCustodianAddress string  `json:"assetCustodianAddress"`
	AssetCustodianCountry string  `gorm:"size:3" json:"assetCustodianCountry"`
	RequirementDocument   string  `gorm:"" json:"requirementDocument"`
	FeePercent            float64 `gorm:"default:0" json:"FeePercent"`
	FeeFixed              float64 `gorm:"default:0" json:"FeeFixed"`
}
type AssetManager struct {
	ID                  uint64  `gorm:"" json:"id"`
	AssetManagerName    string  `gorm:"size:100" json:"assetManagerName"`
	AssetManagerAddress string  `json:"assetManagerAddress"`
	AssetManagerCountry string  `gorm:"size:3" json:"assetManagerCountry"`
	FeePercent          float64 `gorm:"default:0" json:"FeePercent"`
	FeeFixed            float64 `gorm:"default:0" json:"FeeFixed"`
}

type AssetIssuingHouse struct {
	ID                       uint64  `gorm:"" json:"id"`
	AssetIssuingHouseName    string  `gorm:"size:100" json:"assetIssuingHouseName"`
	AssetIssuingHouseAddress string  `json:"assetIssuingHouseAddress"`
	AssetIssuingHouseCountry string  `gorm:"size:3" json:"assetIssuingHouseCountry"`
	FeePercent               float64 `gorm:"default:0" json:"FeePercent"`
	FeeFixed                 float64 `gorm:"default:0" json:"FeeFixed"`
}

type LegalAndProfesionalPartner struct {
	ID             uint64  `gorm:"" json:"id"`
	PartnerName    string  `gorm:"size:100" json:"partnerName"`
	PartnerAddress string  `json:"partnerAddress"`
	PartnerCountry string  `gorm:"size:3" json:"partnerCountry"`
	FeePercent     float64 `gorm:"default:0" json:"FeePercent"`
	FeeFixed       float64 `gorm:"default:0" json:"FeeFixed"`
}

type RatingAgency struct {
	ID            uint64  `gorm:"" json:"id"`
	AgencyName    string  `gorm:"size:100" json:"agencyName"`
	AgencyAddress string  `json:"agencyAddress"`
	AgencyCountry string  `gorm:"size:3" json:"agencyCountry"`
	FeePercent    float64 `gorm:"default:0" json:"FeePercent"`
	FeeFixed      float64 `gorm:"default:0" json:"FeeFixed"`
}

type Trustee struct {
	ID             uint64  `gorm:"" json:"id"`
	TrusteeName    string  `gorm:"size:100" json:"trusteeName"`
	TrusteeAddress string  `json:"trusteeAddress"`
	TrusteeCountry string  `gorm:"size:3" json:"trusteeCountry"`
	FeePercent     float64 `gorm:"default:0" json:"FeePercent"`
	FeeFixed       float64 `gorm:"default:0" json:"FeeFixed"`
}

type AssetProtectionOption struct {
	ID string `gorm:"size:100;primaryKey" json:"id"`
}

type ProceedCycle struct {
	ID string `gorm:"size:50" json:"id"`
}

type AssetTokenizationDocument struct {
	ID               uint64 `json:"id"`
	CreatedAt        time.Time
	TokenizedAssetID string `json:"tokenizedAssetId"`
	DocumentType     string `json:"documentType"`
	DocumentTitle    string `json:"documentTitle"`
	DocumentUrl      string `json:"documentUrl"`
	ShowToPublic     int    `gorm:"default:0" json:"showToPublic"`
}
type AssetTokenizationInputDocument struct {
	ID               uint64    `json:"-" form:"-"`
	CreatedAt        time.Time `json:"-" form:"-"`
	TokenizedAssetID string    `json:"tokenizedAssetId" form:"tokenizedAssetId"`
	DocumentType     string    `json:"documentType" form:"documentType"`
	DocumentTitle    string    `json:"documentTitle" form:"documentTitle"`
	// DocumentFile     *multipart.File `form:"documentFile"`
	// DocumentFile     string `json:"-"`// this is not included in struct for input. already extracted by c.FormFile
}

type FailDueDiligence struct {
	Reason string `json:"reason" form:"reason"`
}

type AssetTokenizationDocumentType struct {
	ID                      uint64 `json:"-" form:"-"`
	DocumentType            string `json:"documentType" form:"documentType"`
	DocumentTypeDescription string `json:"documentTypeDescription" form:"documentTypeDescription"`
	DocumentCategory        string `json:"documentCategory" form:"documentCategory"`
	IsPublic                int    `gorm:"default:0" json:"isPublic" form:"isPublic"`
	// ColumnName              string `json:"-" form:"-"`
}

type IssuingWalletPublicKey string

type ExistingAssetValidationAssetInformation struct {
	ID                             uint64 `gorm:"" json:"-" form:"-"`
	AssetAlreadyExists0            int    `gorm:"default:1" json:"assetAlreadyExists0"`
	AssetAlreadyExists1            int    `gorm:"default:1" json:"assetAlreadyExists1"`
	AssetDescription               int    `gorm:"default:1" json:"assetDescription"`
	AssetPhysicalAddress           int    `gorm:"default:1" json:"assetPhysicalAddress"`
	AssetCountryLocation           int    `gorm:"default:1" json:"assetCountryLocation"`
	AssetLongitude                 int    `gorm:"default:1" json:"assetLongitude"`
	AssetLatitude                  int    `gorm:"default:1" json:"assetLatitude"`
	OwnershipTypeDirect            int    `gorm:"default:1" json:"ownershipTypeDirect"`
	OwnershipType3p                int    `gorm:"default:1" json:"ownershipType3p"` //3rd party
	OwnershipKindIndividual        int    `gorm:"default:1" json:"ownershipKindIndividual"`
	OwnershipKindCorporate         int    `gorm:"default:1" json:"ownershipKindCorporate"`
	AssetOwnerName                 int    `gorm:"default:1" json:"assetOwnerName"`
	AssetOwnerAddress              int    `gorm:"default:1" json:"assetOwnerAddress"`
	ApprovedAssetCustodianID       int    `gorm:"default:1" json:"approvedAssetCustodianId"`
	AssetManagerID                 int    `gorm:"default:1" json:"assetManagerId"`
	AssetCurrentValue              int    `gorm:"default:1" json:"AssetCurrentValue"` //current(existing asset)/proposed (non-Existing asset) value
	ValueOfTokenizedAsset          int    `gorm:"default:1" json:"valueOfTokenizedAsset"`
	AssetMscCostOutisdeOfValuation int    `gorm:"default:1" json:"assetMscCostOutisdeOfValuation"`
	ProtectionMethods              int    `gorm:"default:1" json:"protectionMethods"`
	InsuranceCompanyName           int    `gorm:"default:1" json:"insuranceCompanyName"`
	InsurancePolicyNumber          int    `gorm:"default:1" json:"insurancePolicyNumber"`
	InsurancePolicyHolder          int    `gorm:"default:1" json:"insurancePolicyHolder"`
	PercentageValueOfInsurance     int    `gorm:"default:1" json:"percentageValueOfInsurance"`
	IsFreeFromLiensAndEncumbrances int    `gorm:"default:1" json:"isFreeFromLiensAndEncumbrances"`
}

type ExistingAssetValidationAssetTokenInfo struct {
	ID uint64 `gorm:"" json:"-" form:"-"`
}

type NonExistingAssetValidationAssetInformation struct {
	ID                             uint64 `gorm:"" json:"-" form:"-"`
	AssetAlreadyExists0            int    `gorm:"default:1" json:"assetAlreadyExists0"`
	AssetAlreadyExists1            int    `gorm:"default:1" json:"assetAlreadyExists1"`
	AssetDescription               int    `gorm:"default:1" json:"assetDescription"`
	AssetPhysicalAddress           int    `gorm:"default:1" json:"assetPhysicalAddress"`
	AssetCountryLocation           int    `gorm:"default:1" json:"assetCountryLocation"`
	AssetLongitude                 int    `gorm:"default:1" json:"assetLongitude"`
	AssetLatitude                  int    `gorm:"default:1" json:"assetLatitude"`
	OwnershipTypeDirect            int    `gorm:"default:1" json:"ownershipTypeDirect"`
	OwnershipType3p                int    `gorm:"default:1" json:"ownershipType3p"` //3rd party
	OwnershipKindIndividual        int    `gorm:"default:1" json:"ownershipKindIndividual"`
	OwnershipKindCorporate         int    `gorm:"default:1" json:"ownershipKindCorporate"`
	AssetOwnerName                 int    `gorm:"default:1" json:"assetOwnerName"`
	AssetOwnerAddress              int    `gorm:"default:1" json:"assetOwnerAddress"`
	ApprovedAssetCustodianID       int    `gorm:"default:1" json:"approvedAssetCustodianId"`
	AssetManagerID                 int    `gorm:"default:1" json:"assetManagerId"`
	AssetCurrentValue              int    `gorm:"default:1" json:"AssetCurrentValue"` //current(existing asset)/proposed (non-Existing asset) value
	ValueOfTokenizedAsset          int    `gorm:"default:1" json:"valueOfTokenizedAsset"`
	AssetMscCostOutisdeOfValuation int    `gorm:"default:1" json:"assetMscCostOutisdeOfValuation"`
	ProtectionMethods              int    `gorm:"default:1" json:"protectionMethods"`
	InsuranceCompanyName           int    `gorm:"default:1" json:"insuranceCompanyName"`
	InsurancePolicyNumber          int    `gorm:"default:1" json:"insurancePolicyNumber"`
	InsurancePolicyHolder          int    `gorm:"default:1" json:"insurancePolicyHolder"`
	PercentageValueOfInsurance     int    `gorm:"default:1" json:"percentageValueOfInsurance"`
	IsFreeFromLiensAndEncumbrances int    `gorm:"default:1" json:"isFreeFromLiensAndEncumbrances"`
}
type TokenMinting struct {
	TokenizedAssetID     string   `json:"tokenizedAssetId"`
	Destination          string   `json:"destination" `
	Amount               string   `json:"amount" `
	AssetCode            string   `json:"assetCode"`
	AssetIssuer          string   `json:"assetIssuer"`
	Transaction          string   `json:"transaction"`
	TransactionSignature string   `json:"transactionSignature"`
	TransactionID        string   `json:"transactionId"`
	NetworkPassPhrase    string   `json:"networkPassPhrase"`
	TransactionSource    string   `json:"-"`
	ReturnedDescription  string   `json:"-"`
	Commit               int      `json:"commit"`
	Messages             []string `json:"messages"`
}

type TokenizedAssetSubscription struct {
	ID                 string         `gorm:"" json:"-" form:"-"`
	CreatedAt          time.Time      `json:"createdAt"`
	UpdatedAt          time.Time      `json:"updatedAt"`
	TokenizedAssetID   string         `gorm:"not null;size:100" json:"tokenizedAssetId"`
	TokenizedAsset     TokenizedAsset `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"tokenizedAssetInfo"`
	AssetCode          string         `gorm:"not null;size:12" json:"assetCode"`
	AssetIssuer        string         `gorm:"not null;size:100" json:"assetIssuer"`
	WalletAlias        string         `gorm:"not null;size:100" json:"walletAlias"`
	WalletPublicKey    string         `gorm:"not null;size:100" json:"walletPublicKey"`
	Amount             float64        `json:"amount"` //fiat Amount in tokenized asset quote currency
	Price              float64        `json:"price"`  // in tokenized asset price in quote currency
	SubscriberUsername string         `gorm:"not null;size:100" json:"subscriberUsername"`
	TransactionID      string         `json:"transactionId"`
}
type ExpressionOfInterest struct {
	ID                 uint64         `gorm:"" json:"-" form:"-"`
	CreatedAt          time.Time      `json:"createdAt"`
	UpdatedAt          time.Time      `json:"updatedAt"`
	TokenizedAssetID   string         `gorm:"not null;size:100" json:"tokenizedAssetId"`
	TokenizedAsset     TokenizedAsset `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"tokenizedAssetInfo"`
	AssetCode          string         `gorm:"not null;size:12" json:"assetCode"`
	AssetIssuer        string         `gorm:"not null;size:100" json:"assetIssuer"`
	Amount             float64        `json:"amount"` //fiat Amount in tokenized asset quote currency
	Price              float64        `json:"price"`
	SubscriberUsername string         `gorm:"not null;size:100" json:"subscriberUsername"`
}

type TokenizedAssetSubscriptionInput struct {
	TokenizedAssetID     string   `json:"tokenizedAssetId"`
	SubscriberUsername   string   `json:"subscriberUsername"`
	WalletPublicKey      string   `json:"walletPublicKey"`
	Amount               float64  `json:"amount"` //fiat Amount in tokenized asset quote currency
	SwappedEstimate      string   `json:"swappedEstimate"`
	Transaction          string   `json:"transaction"`
	TransactionSignature string   `json:"transactionSignature"`
	TransactionID        string   `json:"transactionId"`
	NetworkPassPhrase    string   `json:"networkPassPhrase"`
	Messages             []string `json:"messages"`
	Memo                 string   `json:"memo"`
	Multiparty           int      `json:"-"`
	TransactionSource    string   `json:"-"`
	SignatureRequired    int      `json:"signatureRequired"`
	Commit               int      `json:"commit"`
	ReturnedDescription  string   `json:"-"`
}
type TokenizedAssetPrimarySalesPurchaseInputForServiceLink struct {
	TokenizedAssetID           string   `json:"tokenizedAssetId"`
	PurchaserUsername          string   `json:"purchaserUsername"`
	DestinationWalletPublicKey string   `json:"destinationWalletPublicKey"`
	AmountInAssetCurrency      float64  `json:"amountInAssetCurrency"` //fiat Amount in tokenized asset quote currency
	SwappedEstimate            string   `json:"swappedEstimate"`
	Transaction                string   `json:"transaction"`
	TransactionSignature       string   `json:"transactionSignature"`
	TransactionID              string   `json:"transactionId"`
	NetworkPassPhrase          string   `json:"networkPassPhrase"`
	Messages                   []string `json:"messages"`
	Memo                       string   `json:"memo"`
	SignatureRequired          int      `json:"signatureRequired"`
	Commit                     int      `json:"commit"`
}

func (i *TokenizedAssetSubscriptionInput) ToServiceLinkInput(gc *sharedconfig.GlobalConfig) (si TokenizedAssetPrimarySalesPurchaseInputForServiceLink) {
	si.TokenizedAssetID = i.TokenizedAssetID
	si.PurchaserUsername = i.SubscriberUsername
	si.DestinationWalletPublicKey = i.WalletPublicKey
	si.AmountInAssetCurrency = i.Amount
	si.SwappedEstimate = i.SwappedEstimate
	si.Transaction = i.Transaction
	si.TransactionSignature = i.TransactionSignature
	si.TransactionID = i.TransactionID
	si.NetworkPassPhrase = i.NetworkPassPhrase
	si.Messages = i.Messages
	si.Memo = i.Memo
	si.SignatureRequired = i.SignatureRequired
	si.Commit = i.Commit
	return si
}

func (i *TokenizedAssetPrimarySalesPurchaseInputForServiceLink) ToSubscriptionInput(gc *sharedconfig.GlobalConfig) (si TokenizedAssetSubscriptionInput) {
	si.TokenizedAssetID = i.TokenizedAssetID
	si.SubscriberUsername = i.PurchaserUsername
	si.WalletPublicKey = i.DestinationWalletPublicKey
	si.Amount = i.AmountInAssetCurrency
	si.SwappedEstimate = i.SwappedEstimate
	si.Transaction = i.Transaction
	si.TransactionSignature = i.TransactionSignature
	si.TransactionID = i.TransactionID
	si.NetworkPassPhrase = i.NetworkPassPhrase
	si.Messages = i.Messages
	si.Memo = i.Memo
	si.SignatureRequired = i.SignatureRequired
	si.Commit = i.Commit
	return si
}

type ExpressionOfInterestInput struct {
	Amount float64 `json:"amount"` //fiat Amount in tokenized asset quote currency
}

type TokenizedAssetSalesDatesInput struct {
	SalesStart string `json:"salesStart"`
	SalesEnd   string `json:"salesEnd"`
}

type NonExistingAssetValidationAssetDocument struct {
	ID uint64 `gorm:"" json:"-" form:"-"`
}
type NonExistingAssetValidationAssetTokenInfo struct {
	ID uint64 `gorm:"" json:"-" form:"-"`
}

type ProceedPayout struct {
	ID                          uint64         `gorm:"" json:"-" form:"-"`
	CreatedAt                   time.Time      `json:"createdAt"`
	UpdatedAt                   time.Time      `json:"updatedAt"`
	TokenizedAssetID            string         `gorm:"not null;size:100" json:"tokenizedAssetId"`
	TokenizedAsset              TokenizedAsset `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"tokenizedAssetInfo"`
	Batch                       string         `gorm:"not null;size:100;index:,unique" json:"batch"` //asset code + payout cycle + month + year
	DepositedProceedAmount      float64        `json:"DepositedProceedAmount"`                       //fiat Amount in tokenized asset quote currency
	PlatformFee                 float64        `json:"PlatformFee"`                                  //fiat Amount in tokenized asset quote currency
	ProceedPayoutAmount         float64        `json:"proceedPayoutAmount"`                          //fiat Amount in tokenized asset quote currency
	AmountPerTokenizedAssetHeld float64        `json:"amountPerTokenizedAssetHeld"`                  //Amount of proceed for each tokenized asset in quote currency
	PaymentScheduleReady        int            `gorm:"default:0" json:"paymentScheduleReady"`        //tracks if payment schedule is ready
	PayoutCompleted             int            `gorm:"default:0" json:"PayoutCompleted"`             //tracks if payment is completed

}

type TokenizedAssetPayoutSchedule struct {
	ID                             string         `gorm:"" json:"-" form:"-"`
	CreatedAt                      time.Time      `json:"createdAt"`
	TokenizedAssetID               string         `gorm:"not null;size:100" json:"tokenizedAssetId"`
	TokenizedAsset                 TokenizedAsset `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"tokenizedAssetInfo"`
	Batch                          string         `gorm:"not null;size:100;index:," json:"batch"` //asset code + payout cycle + month + year
	PayoutAssetCode                string         `gorm:"not null;size:12" json:"payoutAssetCode"`
	PayoutAssetIssuer              string         `gorm:"not null;size:100" json:"payoutAssetIssuer"`
	BeneficiaryPublicKey           string         `gorm:"not null;size:100" json:"beneficiaryPublicKey"`
	ConfirmedTokenizedAssetBalance float64        `json:"confirmedTokenizedAssetBalance"` //asset balance at the time of preparing schedule
	AmountToReceive                float64        `json:"amountToReceive"`
	CannotReceiveAsset             int            `gorm:"default:0" json:"CannotReceiveAsset"` //checks if the beneficiary can receive the asset or not.
}

type TokenizedAssetPayoutEngineTask struct {
	ID                             uint64    `gorm:"" json:"-" form:"-"`
	TokenizedAssetPayoutScheduleID string    `gorm:"size:100;index:,unique" json:"tokenizedAssetPayoutScheduleID"`
	CreatedAt                      time.Time `json:"createdAt"`
	MemoFromBatch                  string    `gorm:"not null;size:28;index:," json:"batch"` //asset code + payout cycle + month + year
	PayoutAssetCode                string    `gorm:"not null;size:12" json:"payoutAssetCode"`
	PayoutAssetIssuer              string    `gorm:"not null;size:100" json:"payoutAssetIssuer"`
	BeneficiaryPublicKey           string    `gorm:"not null;size:100" json:"beneficiaryPublicKey"`
	AmountToReceive                string    `json:"amountToReceive"`
	Paid                           int       `gorm:"default:0" json:"paid"`
}

type PostTokenizationTrustlineCandidate struct {
	ID          uint64
	PublicKey   string
	Description string
}

func (p PostTokenizationTrustlineCandidate) GetUntrustedTokenizedAssets(gc *sharedconfig.GlobalConfig) (untrusted []string) {
	knownAssets := make(map[string]struct{})
	known := make([]string, 0)
	//get curated tokenized assets
	assets := gc.GetCuratedAssetByClassID(3, false)
	//get account balance
	account, _, err := UserWalletID(p.PublicKey).GetBlockchainAccountDetail(gc)
	if err != nil {
		return
	}

	for _, a := range assets {
		knownAssets[a.AssetCode+":"+a.AssetIssuer] = struct{}{}
	}

	for _, b := range account.Balances {
		key := b.Code + ":" + b.Issuer
		if _, exists := knownAssets[key]; exists {
			known = append(known, key)
		}
	}
	//range through known and remove the ones that exists
	if len(known) > 0 {
		//remove the ones that exists
		for _, v := range known {
			delete(knownAssets, v)
		}
	}

	for key := range knownAssets {
		untrusted = append(untrusted, key)
	}

	return untrusted
}

func (p *ProceedPayout) CreateBatch() error {
	if len(p.TokenizedAssetID) == 0 {
		return &tErrors.CustomError{Err: "error invalid tokenizedAssetId", ErrMessage: "tokenized asset identification is invalid"}
	}
	if p.TokenizedAsset.ProceedCycle == nil {
		return &tErrors.CustomError{Err: "error proceed-cycle-not-set", ErrMessage: "Proceed Cycle was not set for this project"}

	}
	if *p.TokenizedAsset.ProceedCycle != "None" {
		return &tErrors.CustomError{Err: "error proceed-cycle-not-set", ErrMessage: "Proceed Cycle was not set for this project"}

	}
	//asset code + payout cycle + month + year
	t := time.Now()
	p.Batch = fmt.Sprintf("%v%v%v%v", p.TokenizedAsset.AssetCode, *p.TokenizedAsset.ProceedCycle, t.Month().String(), t.Year())
	return nil
}

type TokenizedAssetCode string
type AssetTokenizationDocumentID uint64

func (did AssetTokenizationDocumentID) GetTokenization(gc *sharedconfig.GlobalConfig) (t TokenizedAsset) {
	//get document
	var d AssetTokenizationDocument
	e := gc.DB.Where("id = ?", uint64(did)).First(&d).Error
	if e != nil {
		log.Printf("[AssetTokenizationDocumentID::GetTokenization] Error getting tokenized asset for document ID %v, %v\n", uint64(did), e)
		return
	}

	e = gc.DB.Preload(clause.Associations).Order("asset_tokenization_status ASC").Where("id = ?", string(d.TokenizedAssetID)).First(&t).Error
	if e != nil {
		log.Printf("[AssetTokenizationDocumentID::GetTokenization] Error getting tokenized asset for %v, %v\n", d.TokenizedAssetID, e)
	}
	return
}

func (i IssuingWalletPublicKey) GetTokenizationByDocumentID(did uint64, gc *sharedconfig.GlobalConfig) (t TokenizedAsset) {
	return AssetTokenizationDocumentID(did).GetTokenization(gc)
}

// GetTokenization gets the tokenized asset by issuing wallet and returns the first one ordered by the asset tokenization status from 0.
func (i IssuingWalletPublicKey) GetTokenizationByAssetCode(assetCode string, gc *sharedconfig.GlobalConfig) (t TokenizedAsset) {
	e := gc.DB.Preload(clause.Associations).Order("asset_tokenization_status ASC").Where("issuing_wallet_public_key = ? AND asset_code = upper(?)", string(i), assetCode).First(&t).Error
	if e != nil {
		log.Printf("[IssuingWalletPublicKey::GetTokenization] Error getting tokenized asset for %v, %v\n", string(i), e)
	}
	return
}
func (i IssuingWalletPublicKey) GetTokenizationByID(id string, gc *sharedconfig.GlobalConfig) (t TokenizedAsset) {
	e := gc.DB.Preload(clause.Associations).Where("issuing_wallet_public_key = ? AND id = ?", string(i), id).First(&t).Error
	if e != nil {
		log.Printf("[IssuingWalletPublicKey::GetTokenization] Error getting tokenized asset for %v, %v\n", string(i), e)
	}
	return
}

func (i TokenizedAssetID) GetTokenization(gc *sharedconfig.GlobalConfig) (t TokenizedAsset) {
	e := gc.DB.Preload(clause.Associations).Where("id = ?", string(i)).First(&t).Error
	if e != nil {
		log.Printf("[TokenizedAssetID::GetTokenization] Error getting tokenized asset for %v, %v\n", string(i), e)
	}
	return
}

func (t TokenizedAssetCode) IsTokenizedAsset(gc *sharedconfig.GlobalConfig) bool {
	strt := strings.Split(string(t), ":")
	if len(strt) != 2 {
		return false
	}
	code, issuer := strt[0], strt[1]

	e := gc.DB.Where("Asset_Tokenization_Status > 4 AND Asset_Code = upper(?) AND Issuing_Wallet_Public_Key = upper(?)", code, issuer).First(&TokenizedAsset{}).Error

	return e == nil
}

func (i IssuingWalletPublicKey) GetTokenizationFeeByID(feeID uint64, gc *sharedconfig.GlobalConfig) (fee TokenizationFee) {
	gc.DB.Preload(clause.Associations).Where("id = ?", feeID).First(&fee)

	return
}

func (t *TokenizedAsset) GetTokenizationFeeByID(feeID uint64, gc *sharedconfig.GlobalConfig) (fee TokenizationFee) {
	gc.DB.Preload(clause.Associations).Where("id = ?", feeID).First(&fee)

	return
}
func (t *TokenizedAsset) GetExpressedInterestByUsername(subscriber string, gc *sharedconfig.GlobalConfig) (exp ExpressionOfInterest, err error) {
	if t == nil {
		log.Println("[TokenizedAsset::GetExpressedInterestByUsername] Error tokenized asset is nil")
		err = &tErrors.ErrorTemporaryServerError{}
		return
	}
	var issuerWalletPublicKey string
	if t.IssuingWalletPublicKey != nil {
		issuerWalletPublicKey = *t.IssuingWalletPublicKey
	}
	err = gc.DB.Preload(clause.Associations).Where("Tokenized_Asset_ID = ? AND Asset_Issuer = ? AND Subscriber_Username = ?", t.ID, issuerWalletPublicKey, subscriber).First(&exp).Error

	return
}

// SumExpressedInterest sums the amount that is commited in expression of interest
func (t *TokenizedAsset) SumExpressedInterest(gc *sharedconfig.GlobalConfig) (sum float64) {
	if t == nil {
		log.Println("[TokenizedAsset::SumExpressedInterest] Error tokenized asset is nil")
		return
	}
	var issuerWalletPublicKey string
	if t.IssuingWalletPublicKey != nil {
		issuerWalletPublicKey = *t.IssuingWalletPublicKey
	}
	gc.DB.Model(&ExpressionOfInterest{}).Where("Tokenized_Asset_ID = ? AND Asset_Issuer = ?", t.ID, issuerWalletPublicKey).Select("case when sum(Amount) is not null then sum(Amount) else 0 end").Row().Scan(&sum)

	return
}

func (t *TokenizedAsset) CountExpressedInterests(gc *sharedconfig.GlobalConfig) (count int64) {
	if t == nil {
		log.Println("[TokenizedAsset::CountExpressedInterests] Error tokenized asset is nil")
		return
	}
	var issuerWalletPublicKey string
	if t.IssuingWalletPublicKey != nil {
		issuerWalletPublicKey = *t.IssuingWalletPublicKey
	}
	gc.DB.Model(ExpressionOfInterest{}).Where("Tokenized_Asset_ID = ? AND Asset_Issuer = ?", t.ID, issuerWalletPublicKey).Count(&count)

	return
}

func (t *TokenizedAsset) CountNumberOfSubscribers(gc *sharedconfig.GlobalConfig) (count int64) {
	if t == nil {
		log.Println("[TokenizedAsset::CountNumberOfSubscribers] Error tokenized asset is nil")
		return
	}
	var issuerWalletPublicKey string
	if t.IssuingWalletPublicKey != nil {
		issuerWalletPublicKey = *t.IssuingWalletPublicKey
	}
	gc.DB.Model(TokenizedAssetSubscription{}).Where("Tokenized_Asset_ID = ? AND Asset_Issuer = ?", t.ID, issuerWalletPublicKey).Count(&count)

	return
}

// SumAmountSoldInFiat sums the total amount in fiat sold so far
func (t *TokenizedAsset) SumAmountSoldInFiat(gc *sharedconfig.GlobalConfig) (sum float64) {
	if t == nil {
		log.Println("[TokenizedAsset::SumAmountSoldInFiat] Error tokenized asset is nil")
		return
	}
	var issuerWalletPublicKey string
	if t.IssuingWalletPublicKey != nil {
		issuerWalletPublicKey = *t.IssuingWalletPublicKey
	}
	gc.DB.Model(TokenizedAssetSubscription{}).Where("Tokenized_Asset_ID = ? AND Asset_Issuer = ?", t.ID, issuerWalletPublicKey).Select("case when sum(Amount) is not null then sum(Amount) else 0 end").Row().Scan(&sum)

	return
}

// SumAmountBoughtByWalletOwner sums the amout that is purchased by the wallet owner
func (t *TokenizedAsset) SumAmountBoughtByWalletOwner(walletAlias string, gc *sharedconfig.GlobalConfig) (sum float64) {
	if t == nil {
		log.Println("[TokenizedAsset::SumAmountBoughtByWalletOwner] Error tokenized asset is nil")
		return
	}
	ownerUsername := strings.Split(walletAlias, "_")[0]
	var issuerWalletPublicKey string
	if t.IssuingWalletPublicKey != nil {
		issuerWalletPublicKey = *t.IssuingWalletPublicKey
	}
	gc.DB.Model(&TokenizedAssetSubscription{}).Where("Tokenized_Asset_ID = ? AND asset_issuer = ? AND Wallet_Alias LIKE ?", t.ID, issuerWalletPublicKey, strings.ToLower(ownerUsername)+"%").Select("case when sum(Amount) is not null then sum(Amount) else 0 end").Row().Scan(&sum)

	return
}

func (t *TokenizedAsset) GetTokenizedAssetSubscriptionByWalletPublicKey(subscriberWalletPublicKey string, gc *sharedconfig.GlobalConfig) (sub TokenizedAssetSubscription, err error) {
	if t == nil {
		log.Printf("[TokenizedAsset::GetTokenizedAssetSubscriptionByWalletPublicKey] Error tokenized asset is nil %v\n", subscriberWalletPublicKey)
		err = &tErrors.ErrorTemporaryServerError{}
		return
	}
	var issuerWalletPublicKey string
	if t.IssuingWalletPublicKey != nil {
		issuerWalletPublicKey = *t.IssuingWalletPublicKey
	}
	err = gc.DB.Preload(clause.Associations).Where("Tokenized_Asset_ID = ? AND Asset_Issuer = ? AND Wallet_Public_Key = ?", t.ID, issuerWalletPublicKey, subscriberWalletPublicKey).First(&sub).Error

	return
}

func (t *TokenizedAsset) UpdateTokenizationFeeByID(feeID uint64, gc *sharedconfig.GlobalConfig) (fee TokenizationFee) {
	if feeID == 0 {
		return
	}
	gc.DB.Preload(clause.Associations).Where("id = ?", feeID).First(&fee)

	if fee.ID > 0 {
		t.TokenizationFee = fee
		t.TokenizationFeeID = &feeID
	}

	return
}

func (t *TokenizedAsset) UpdateBank(gc *sharedconfig.GlobalConfig) (bank Bank) {
	if t.BankID == nil {
		return
	}
	if *t.BankID == 0 {
		return
	}

	gc.DB.Preload(clause.Associations).Where("id = ?", *t.BankID).First(&bank)

	if bank.ID > 0 {
		t.Bank = bank

	}
	return
}

func (t *TokenizedAsset) UpdateTokenizedAssetFromInput(ti *TokenizedAssetJSONInput, gc *sharedconfig.GlobalConfig) TokenizedAsset {
	titleCaser := cases.Title(language.English)
	replacer := strings.NewReplacer("\r", "", "\n", "", " ", "")
	if t.AssetTokenizationStatus < 4 {
		t.HasAdditionalKYCRequirements = ti.HasAdditionalKYCRequirements

		if len(ti.AdditionalKYCRequirements) > 0 && t.HasAdditionalKYCRequirements > 0 {

			t.AdditionalKYCRequirements = &ti.AdditionalKYCRequirements
		} else {
			t.AdditionalKYCRequirements = nil
			t.HasAdditionalKYCRequirements = 0
		}
		//clear any DD failure flags
		t.DueDiligenceFail = 0
		t.DueDiligenceFailureReason = nil

		if len(ti.AssetSector) > 0 {

			t.AssetSector = &ti.AssetSector
		}

		if len(ti.AssetSubSector) > 0 {

			t.AssetSubSector = &ti.AssetSubSector
		}

		if len(ti.AssetType) > 0 {

			t.AssetType = &ti.AssetType
		}

		if len(ti.AssetWebsite) > 0 {
			ti.AssetWebsite = strings.ToLower(ti.AssetWebsite)
			t.AssetWebsite = &ti.AssetWebsite
		} else {
			t.AssetWebsite = nil
		}

		if len(ti.MintingApprovers) > 0 {
			trimmed := replacer.Replace(ti.MintingApprovers)
			ti.MintingApprovers = trimmed
			t.MintingApprovers = &trimmed
		} else {
			//set default
			csvStr := t.GetMintingApproversInCSV(gc)
			t.MintingApprovers = &csvStr
		}

		if len(ti.MintingInitators) > 0 {
			trimmed := replacer.Replace(ti.MintingInitators)
			ti.MintingInitators = trimmed
			t.MintingInitators = &trimmed
		} else {
			//set default
			csvStr := t.GetMintingInitiatorsInCSV(gc)
			t.MintingInitators = &csvStr
		}

		if t.AssetTokenizationStatus < 3 {
			//after ccreting wallet, do not allow change of name and code.

			if len(ti.AssetName) > 0 {
				ti.AssetName = titleCaser.String(ti.AssetName)
				t.AssetName = &ti.AssetName
			}
			if len(ti.AssetCode) > 0 {
				ti.AssetCode = strings.ToUpper(ti.AssetCode)

				t.AssetCode = &ti.AssetCode
			}

			if t.AssetTokenizationStatus < 2 {
				//value and currency, asset manager and custodian and issuing house can change here.

				if len(ti.AssetCountryLocation) > 0 {

					t.AssetCountryLocation = &ti.AssetCountryLocation
				}

				if len(ti.AssetQuoteCurrency) > 0 {
					ti.AssetQuoteCurrency = strings.ToUpper(ti.AssetQuoteCurrency)
					t.AssetQuoteCurrency = &ti.AssetQuoteCurrency
				} else {
					t.AssetQuoteCurrency = nil
				}
				t.ApprovedAssetCustodianID = ti.ApprovedAssetCustodianID
				if t.ApprovedAssetCustodianID == 0 {
					//set default
					t.ApprovedAssetCustodianID = 1
				}
				t.AssetIssuingHouseID = ti.AssetIssuingHouseID
				if t.AssetIssuingHouseID == 0 {
					//set default
					t.AssetIssuingHouseID = 1
				}

				t.AssetOwnerRetainedOrContributedValue = ti.AssetOwnerRetainedOrContributedValue
				t.AssetManagerID = ti.AssetManagerID
				if t.AssetManagerID == 0 {
					//set default
					t.AssetManagerID = 1
				}
				t.LegalAndProfesionalPartnerID = ti.LegalAndProfesionalPartnerID
				if t.LegalAndProfesionalPartnerID == 0 {
					//set default
					t.LegalAndProfesionalPartnerID = 1
				}
				t.RatingAgencyID = ti.RatingAgencyID
				if t.RatingAgencyID == 0 {
					//set default
					t.RatingAgencyID = 1
				}

				t.AssetCurrentValue = ti.AssetCurrentValue
				t.AssetMscCostOutisdeOfValuation = ti.AssetMscCostOutisdeOfValuation

			}
			///end change of currency and issuing house and custodian

		}

		if len(ti.OfferingType) > 0 {

			t.OfferingType = &ti.OfferingType
		}

		if len(ti.ClosedGroupID) > 0 {

			t.ClosedGroupID = &ti.ClosedGroupID
		} else {
			t.ClosedGroupID = nil
		}

		if len(ti.SecApprovalIdNumber) > 0 && ti.SecApproval > 0 {

			t.SecApproval = ti.SecApproval

			t.SecApprovalIdNumber = &ti.SecApprovalIdNumber
		} else {
			t.SecApproval = 0

			t.SecApprovalIdNumber = nil
		}

		if len(ti.MarketMakingWallet) > 0 {

			t.MarketMakingWallet = &ti.MarketMakingWallet
		}

		if len(ti.AssetDescription) > 0 {
			ti.AssetDescription = titleCaser.String(ti.AssetDescription)
			t.AssetDescription = &ti.AssetDescription
		}

		if len(ti.AssetLongitude) > 0 {

			t.AssetLongitude = &ti.AssetLongitude
		}

		if len(ti.AssetLatitude) > 0 {

			t.AssetLatitude = &ti.AssetLatitude
		}

		if len(ti.OwnershipType) > 0 {

			t.OwnershipType = &ti.OwnershipType
		}

		if len(ti.OwnershipKind) > 0 {

			t.OwnershipKind = &ti.OwnershipKind
		}

		if len(ti.AssetOwnerName) > 0 {
			ti.AssetOwnerName = strings.ToUpper(ti.AssetOwnerName)

			t.AssetOwnerName = &ti.AssetOwnerName
		}
		if len(ti.ProceedPayoutCurrency) > 0 {
			ti.ProceedPayoutCurrency = strings.ToUpper(ti.ProceedPayoutCurrency)
			t.ProceedPayoutCurrency = &ti.ProceedPayoutCurrency
		} else {
			t.ProceedPayoutCurrency = nil
		}
	}
	///end status < 4 here

	if len(ti.AssetPhysicalAddress) > 0 {
		ti.AssetPhysicalAddress = titleCaser.String(ti.AssetPhysicalAddress)

		t.AssetPhysicalAddress = &ti.AssetPhysicalAddress
	}

	if len(ti.AssetOwnerAddress) > 0 {
		ti.AssetOwnerAddress = titleCaser.String(ti.AssetOwnerAddress)

		t.AssetOwnerAddress = &ti.AssetOwnerAddress
	}

	if len(ti.ProtectionMethods) > 0 {

		t.ProtectionMethods = &ti.ProtectionMethods
	} else {
		t.ProtectionMethods = nil
	}

	if len(ti.InsuranceCompanyName) > 0 {

		t.InsuranceCompanyName = &ti.InsuranceCompanyName
	} else {
		t.InsuranceCompanyName = nil

	}

	if len(ti.InsurancePolicyHolder) > 0 {

		t.InsurancePolicyHolder = &ti.InsurancePolicyHolder
	} else {
		t.InsurancePolicyHolder = nil
	}

	if len(ti.InitialOwnerPreferredWalletAddress) > 0 {

		t.InitialOwnerPreferredWalletAddress = &ti.InitialOwnerPreferredWalletAddress
	} else {
		t.InitialOwnerPreferredWalletAddress = nil
	}

	if len(ti.InsurancePolicyNumber) > 0 {

		t.InsurancePolicyNumber = &ti.InsurancePolicyNumber
	} else {
		t.InsurancePolicyNumber = nil
	}

	t.PercentageValueOfInsurance = ti.PercentageValueOfInsurance
	t.IsFreeFromLiensAndEncumbrances = ti.IsFreeFromLiensAndEncumbrances
	t.AssetAlreadyExists = ti.AssetAlreadyExists
	if t.AssetTokenizationStatus < 2 {
		t.NumberOfTokenToBeIssued = ti.NumberOfTokenToBeIssued
	}
	var feeCompo TokenizationFee
	var feeInAsset, VATAsset float64

	// /////
	var cConfig Country
	var custodyFee, assetMgtFee float64
	var secFee float64
	var vat, issuingHouseFeeValue, legalAndProfessionalFee, ratingAgencyFee, totalChargedFeesForVat, trusteeFee float64
	if t.AssetCountryLocation != nil {
		cConfig = CountryCode(*t.AssetCountryLocation).GetConfig(gc)
	}

	if t.AssetTokenizationStatus < 4 && t.AssetCurrentValue > 0 {

		if t.AssetTokenizationStatus < 2 {
			// t.NumberOfTokenToBeIssued = ti.NumberOfTokenToBeIssued
			//Do not change fees for assets that have been approved
			//only when fee has not been paid
			if ti.TokenizationFeeID > 0 {
				// fee has been selected
				t.TokenizationFeeID = &ti.TokenizationFeeID
				feeCompo = t.UpdateTokenizationFeeByID(ti.TokenizationFeeID, gc)

				// Calculate Fees
				feeInAsset = decimal.NewFromFloat(t.NumberOfTokenToBeIssued * (feeCompo.FeeAssetPercentage / 100)).Truncate(7).InexactFloat64()
				t.FeeInAsset = feeInAsset
				log.Printf("[UpdateTokenizedAssetFromInput] Calculated Fee In Asset:= %v\n", decimal.NewFromFloat(feeInAsset).String())
				t.FeeInAssetPercent = feeCompo.FeeAssetPercentage

				t.FeeInFiat = decimal.NewFromFloat(t.AssetCurrentValue * (feeCompo.FeeFiatPercentage / 100)).Truncate(7).InexactFloat64()
				if feeCompo.FeeFiatCap > t.FeeInFiat {
					t.FeeInFiat = feeCompo.FeeFiatCap
				}
				log.Printf("[UpdateTokenizedAssetFromInput] Calculated Fee In Fiat:= %v\n", decimal.NewFromFloat(t.FeeInFiat).String())

				VATAsset = decimal.NewFromFloat(feeInAsset * (cConfig.VATPercent / 100)).Truncate(7).InexactFloat64()
				t.VATInAsset = VATAsset
				log.Printf("[UpdateTokenizedAssetFromInput] Calculated VAT Asset:= %v\n", decimal.NewFromFloat(VATAsset).String())

			}

			// get tokenization fees.
			if t.ExcludeSecFee == 0 {
				t.SECTokenizationFeePercent = cConfig.SECTokenizationFeePercent
				t.SECTokenizationFeeFixed = cConfig.SECTokenizationFeeFixed
				secFee = (decimal.NewFromFloat(t.AssetCurrentValue * (cConfig.SECTokenizationFeePercent / 100)).Truncate(7)).Add(decimal.NewFromFloat(cConfig.SECTokenizationFeeFixed).Truncate(7)).InexactFloat64()
				t.SECTokenizationFeeValue = secFee
			}

			log.Printf("[UpdateTokenizedAssetFromInput] Calculated SEC Fee Value:= %v\n", decimal.NewFromFloat(secFee).String())

			{
				// t.ApprovedAssetCustodianID = ti.ApprovedAssetCustodianID
				// custodian := ApprovedCustodianID(t.ApprovedAssetCustodianID).GetApprovedCustodian(gc)
				// t.CustodianFeeFixed = custodian.FeeFixed
				// t.CustodianFeePercent = custodian.FeePercent
				custodyFee = (decimal.NewFromFloat(t.AssetCurrentValue * (t.CustodianFeePercent / 100)).Truncate(7)).Add(decimal.NewFromFloat(t.CustodianFeeFixed).Truncate(7)).InexactFloat64()
				t.CustodianFeeValue = custodyFee
				log.Printf("[UpdateTokenizedAssetFromInput] Custodian Fee Value:= %v\n", decimal.NewFromFloat(custodyFee).String())

			}

			{
				// t.AssetManagerID = ti.AssetManagerID
				// assetManager := AssetManagerID(t.AssetManagerID).GetAssetManager(gc)
				// t.AssetManagerFeePercent = assetManager.FeePercent
				// t.AssetManagerFeeFixed = assetManager.FeeFixed
				assetMgtFee = (decimal.NewFromFloat(t.AssetCurrentValue * (t.AssetManagerFeePercent / 100)).Truncate(7)).Add(decimal.NewFromFloat(t.AssetManagerFeeFixed).Truncate(7)).InexactFloat64()
				t.AssetManagerFeeValue = assetMgtFee
				log.Printf("[UpdateTokenizedAssetFromInput] Calculated Asset Managment fee Value:= %v\n", decimal.NewFromFloat(assetMgtFee).String())

			}

			{
				// t.AssetIssuingHouseID = ti.AssetIssuingHouseID
				// issuingHouse := IssuingHouseID(t.AssetIssuingHouseID).GetAssetIssuingHouse(gc)
				// t.IssuingHouseFeeFixed = issuingHouse.FeeFixed
				// t.IssuingHouseFeePercent = issuingHouse.FeePercent
				issuingHouseFeeValue = (decimal.NewFromFloat(t.AssetCurrentValue * (t.IssuingHouseFeePercent / 100)).Truncate(7)).Add(decimal.NewFromFloat(t.IssuingHouseFeeFixed).Truncate(7)).InexactFloat64()
				t.IssuingHouseFeeValue = issuingHouseFeeValue
				log.Printf("[UpdateTokenizedAssetFromInput] Calculated Issuing House Fee Value:= %v\n", decimal.NewFromFloat(issuingHouseFeeValue).String())

			}

			{
				// lpp := LegalAndProfesionalPartnerID(t.LegalAndProfesionalPartnerID).GetLegalAndProfesionalPartner(gc)
				// t.LegalAndProfessionalFeeFixed = lpp.FeeFixed
				// t.LegalAndProfessionalFeePercent = lpp.FeePercent
				legalAndProfessionalFee = (decimal.NewFromFloat(t.AssetCurrentValue * (t.LegalAndProfessionalFeePercent / 100)).Truncate(7)).Add(decimal.NewFromFloat(t.LegalAndProfessionalFeeFixed).Truncate(7)).InexactFloat64()
				t.LegalAndProfessionalFeeValue = legalAndProfessionalFee
				log.Printf("[UpdateTokenizedAssetFromInput] Calculated Legal&Professional Fee Value:= %v\n", decimal.NewFromFloat(legalAndProfessionalFee).String())

			}

			{
				// ra := RatingAgencyID(t.RatingAgencyID).GetRatingAgency(gc)
				// t.RatingAgencyFeeFixed = ra.FeeFixed
				// t.RatingAgencyFeePercent = ra.FeePercent
				ratingAgencyFee = (decimal.NewFromFloat(t.AssetCurrentValue * (t.RatingAgencyFeePercent / 100)).Truncate(7)).Add(decimal.NewFromFloat(t.RatingAgencyFeeFixed).Truncate(7)).InexactFloat64()
				t.RatingAgencyFeeValue = ratingAgencyFee
				log.Printf("[UpdateTokenizedAssetFromInput] Calculated Rating Agency Fee Value:= %v\n", decimal.NewFromFloat(ratingAgencyFee).String())

			}

			{
				// ra := TrusteeID(t.TrusteeID).GetTrustee(gc)
				// t.TrusteeFeeFixed = ra.FeeFixed
				// t.TrusteeFeePercent = ra.FeePercent
				trusteeFee = (decimal.NewFromFloat(t.AssetCurrentValue * (t.TrusteeFeePercent / 100)).Truncate(7)).Add(decimal.NewFromFloat(t.TrusteeFeeFixed).Truncate(7)).InexactFloat64()
				t.TrusteeFeeValue = trusteeFee
				log.Printf("[UpdateTokenizedAssetFromInput] Calculated Trustee Fee Value:= %v\n", decimal.NewFromFloat(trusteeFee).String())

			}

			totalChargedFeesForVat = secFee + custodyFee + assetMgtFee + t.FeeInFiat + issuingHouseFeeValue + legalAndProfessionalFee + ratingAgencyFee + trusteeFee
			log.Printf("[UpdateTokenizedAssetFromInput] Calculated Total charged Fees for VAT:= %v\n", decimal.NewFromFloat(totalChargedFeesForVat).String())

			vat = decimal.NewFromFloat(totalChargedFeesForVat * (cConfig.VATPercent / 100)).Truncate(7).InexactFloat64()

			t.VATPercent = cConfig.VATPercent
			t.VATValue = vat
			log.Printf("[UpdateTokenizedAssetFromInput] Calculated VAT Value:= %v\n", decimal.NewFromFloat(vat).String())

			initialValueOfTokenizedAsset := (t.AssetCurrentValue + t.AssetMscCostOutisdeOfValuation + totalChargedFeesForVat + vat) // - *this is not to be shown on the App*
			t.InitialValueOfTokenizedAsset = initialValueOfTokenizedAsset
			log.Printf("[UpdateTokenizedAssetFromInput] Calculated Initial Value Of Tokenized Asset:= %v\n", decimal.NewFromFloat(initialValueOfTokenizedAsset).String())

			if t.NumberOfTokenToBeIssued > 0 {
				// t.PricePerToken = decimal.NewFromFloat(t.ValueOfTokenizedAsset / t.NumberOfTokenToBeIssued).Truncate(7).InexactFloat64()
				pricePerToken := decimal.NewFromFloat((initialValueOfTokenizedAsset / (t.NumberOfTokenToBeIssued - feeInAsset - VATAsset)))
				t.PricePerToken = pricePerToken.Round(7).InexactFloat64()
				log.Printf("[UpdateTokenizedAssetFromInput] Calculated PricePerToken:= %v, Rounded 7DP Value:= %v\n", pricePerToken.String(), decimal.NewFromFloat(t.PricePerToken).String())
				// auto calculate, token to be held is less the fee. token not to be sold
				totalTokenHeldByManager := decimal.NewFromFloat(t.AssetOwnerRetainedOrContributedValue / t.PricePerToken)
				t.TotalTokenHeldByManager = totalTokenHeldByManager.Truncate(7).InexactFloat64()
				log.Printf("[UpdateTokenizedAssetFromInput] Calculated Total Token Held By Manager:= %v, Truncated 7DP Value:= %v\n", totalTokenHeldByManager.String(), decimal.NewFromFloat(t.TotalTokenHeldByManager).String())

				{
					//ensure correct the number of token to be sold.
					maxTokenToBeSold := decimal.NewFromFloat(t.NumberOfTokenToBeIssued - feeInAsset - t.TotalTokenHeldByManager - VATAsset)

					t.MaxNumberOfTokenAvailableForSale = maxTokenToBeSold.Truncate(7).InexactFloat64()
					t.NumberOfTokenToBeSold = t.MaxNumberOfTokenAvailableForSale
					log.Printf("[UpdateTokenizedAssetFromInput] Calculated Max Token to be sold:= %v, Truncated 7DP Value:= %v\n", maxTokenToBeSold.String(), decimal.NewFromFloat(t.NumberOfTokenToBeSold).String())

				}

				finalValueOfTokenizedAsset := decimal.NewFromFloat(t.NumberOfTokenToBeIssued * t.PricePerToken)
				t.ValueOfTokenizedAsset = finalValueOfTokenizedAsset.Truncate(7).InexactFloat64()
				log.Printf("[UpdateTokenizedAssetFromInput] Calculated Final Value of tokenized Asset:= %v, TRUNCATED 7DP value := %v\n", finalValueOfTokenizedAsset.String(), decimal.NewFromFloat(t.ValueOfTokenizedAsset).String())

			}

			if len(ti.WalletToHoldAssetsNotForSale) > 0 {

				t.WalletToHoldAssetsNotForSale = &ti.WalletToHoldAssetsNotForSale
			} else {
				t.WalletToHoldAssetsNotForSale = nil
			}

		}

	}
	if t.AssetTokenizationStatus < 5 {
		t.SalesStart = ti.SalesStart
	}

	t.SalesEnd = ti.SalesEnd
	t.CapOnPurchase = ti.CapOnPurchase
	t.CapQuantity = ti.CapQuantity
	t.CapDurationInDays = ti.CapDurationInDays
	t.CapAmountInFiat = decimal.NewFromFloat(ti.CapAmountInFiat).Truncate(7).InexactFloat64()
	if ti.CapOnPurchase == 1 && t.PricePerToken > 0 && ti.CapAmountInFiat > 0 {
		t.CapQuantity = decimal.NewFromFloat(ti.CapAmountInFiat / t.PricePerToken).Truncate(7).InexactFloat64()
	}

	if len(ti.ProceedCycle) > 0 {

		t.ProceedCycle = &ti.ProceedCycle
	} else {
		t.ProceedCycle = nil
	}

	if len(ti.ExemptedCountries) > 0 {

		t.ExemptedCountries = &ti.ExemptedCountries
	} else {
		t.ExemptedCountries = nil
	}

	t.HasAdditionalKYCRequirements = ti.HasAdditionalKYCRequirements

	if len(ti.AdditionalKYCRequirements) > 0 && t.HasAdditionalKYCRequirements > 0 {

		t.AdditionalKYCRequirements = &ti.AdditionalKYCRequirements
	} else {
		t.AdditionalKYCRequirements = nil
		t.HasAdditionalKYCRequirements = 0
	}

	t.InvestorAccreditationRequired = ti.InvestorAccreditationRequired
	t.AgreeTransferTitleToCustodian = ti.AgreeTransferTitleToCustodian
	t.ContractualProtectionRevGuarantees = ti.ContractualProtectionRevGuarantees
	t.ContractualProtectionPerfBond = ti.ContractualProtectionPerfBond
	t.ContractualProtectionSLA = ti.ContractualProtectionSLA
	t.RiskSharingMechanismPPPs = ti.RiskSharingMechanismPPPs
	t.RiskSharingMechanismHedgeInstruments = ti.RiskSharingMechanismHedgeInstruments
	t.RiskSharingMechanismCompletionGuarantees = ti.RiskSharingMechanismCompletionGuarantees
	if len(ti.IndependentMonitoringList) > 0 {

		t.IndependentMonitoringList = &ti.IndependentMonitoringList
	} else {
		t.IndependentMonitoringList = nil
	}
	t.ESGSafeguardsSusCerts = ti.ESGSafeguardsSusCerts
	t.ESGSafeguardsCommEngPlans = ti.ESGSafeguardsCommEngPlans
	t.SecurityMeasuresAccessControl = ti.SecurityMeasuresAccessControl
	t.SecurityMeasuresSurveilanceSystems = ti.SecurityMeasuresSurveilanceSystems
	t.SecurityMeasuresOnSiteSecurityPersonnel = ti.SecurityMeasuresOnSiteSecurityPersonnel
	t.SecurityMeasuresPerimeterSecurity = ti.SecurityMeasuresPerimeterSecurity
	t.SecurityMeasuresCriticalInfraProtections = ti.SecurityMeasuresCriticalInfraProtections
	if len(ti.OtherAssetProtection) > 0 {

		t.OtherAssetProtection = &ti.OtherAssetProtection
	} else {
		t.OtherAssetProtection = nil

	}
	if len(ti.LegalAdvisor) > 0 {
		ti.LegalAdvisor = strings.ToUpper(ti.LegalAdvisor)

		t.LegalAdvisor = &ti.LegalAdvisor
	} else {
		t.LegalAdvisor = nil
	}

	if len(ti.FinancialAdvisor) > 0 {
		ti.FinancialAdvisor = strings.ToUpper(ti.FinancialAdvisor)

		t.FinancialAdvisor = &ti.FinancialAdvisor
	} else {
		t.FinancialAdvisor = nil
	}

	if ti.BankID > 0 {

		t.BankID = &ti.BankID
		t.UpdateBank(gc)

	}

	if len(ti.BeneficiaryName) > 0 {
		ti.BeneficiaryName = strings.ToUpper(ti.BeneficiaryName)

		t.BeneficiaryName = &ti.BeneficiaryName

	}

	if len(ti.AccountNumber) > 0 {
		t.AccountNumber = &ti.AccountNumber

	}

	t.UndertakingNoLien = ti.UndertakingNoLien
	t.UndertakingNotCollateral = ti.UndertakingNotCollateral
	t.UndertakingNoClaims = ti.UndertakingNoClaims
	t.UndertakingNoForeclosure = ti.UndertakingNoForeclosure
	t.ComplianceNoViolation = ti.ComplianceNoViolation
	t.ComplianceAllPermits = ti.ComplianceAllPermits
	t.OutstandingFinancialRespNoDebts = ti.OutstandingFinancialRespNoDebts
	t.OutstandingFinancialRespNoHiddenLiabilities = ti.OutstandingFinancialRespNoHiddenLiabilities
	t.RiskManagementFullyInsured = ti.RiskManagementFullyInsured
	t.RiskManagementDeclaredValue = ti.RiskManagementDeclaredValue
	t.PhysicalConditionSound = ti.PhysicalConditionSound
	t.PhysicalConditionNolease = ti.PhysicalConditionNolease
	t.PhysicalConditionNoUndisclosedEasements = ti.PhysicalConditionNoUndisclosedEasements

	// asNE(ti.ProjectStrategicObjectives, t.ProjectStrategicObjectives)
	if ne(ti.ProjectStrategicObjectives) {
		t.ProjectStrategicObjectives = &ti.ProjectStrategicObjectives
	}

	// asNE(ti.ProjectDevelopmentTimeline, t.ProjectDevelopmentTimeline)
	if ne(ti.ProjectDevelopmentTimeline) {
		t.ProjectDevelopmentTimeline = &ti.ProjectDevelopmentTimeline
	}

	// asNE(ti.ProjectKeyMilestoneAndDates, t.ProjectKeyMilestoneAndDates)
	if ne(ti.ProjectKeyMilestoneAndDates) {
		t.ProjectKeyMilestoneAndDates = &ti.ProjectKeyMilestoneAndDates
	}

	// asNE(ti.ProjectScope, t.ProjectScope)
	if ne(ti.ProjectScope) {
		t.ProjectScope = &ti.ProjectScope
	}

	// asNE(ti.ProjectEconomicBenefits, t.ProjectEconomicBenefits)
	if ne(ti.ProjectEconomicBenefits) {
		t.ProjectEconomicBenefits = &ti.ProjectEconomicBenefits
	}

	t.ProjectExpectedNoOfJobs = ti.ProjectExpectedNoOfJobs

	// asNE(ti.ProjectIntendedSocialBenefits, t.ProjectIntendedSocialBenefits)
	if ne(ti.ProjectIntendedSocialBenefits) {
		t.ProjectIntendedSocialBenefits = &ti.ProjectIntendedSocialBenefits
	}

	// asNE(ti.ProjectTechnicalPartners, t.ProjectTechnicalPartners)
	if ne(ti.ProjectTechnicalPartners) {
		t.ProjectTechnicalPartners = &ti.ProjectTechnicalPartners
	}

	// asNE(ti.ProjectFinancialPartners, t.ProjectFinancialPartners)
	if ne(ti.ProjectFinancialPartners) {
		t.ProjectFinancialPartners = &ti.ProjectFinancialPartners
	}

	t.EstimatedProjectIRR = ti.EstimatedProjectIRR
	t.EstimatedProjectROI = ti.EstimatedProjectROI
	t.EstimatedProjectNPV = ti.EstimatedProjectNPV
	t.EstimatedProjectPaybackPeriodsInMonths = ti.EstimatedProjectPaybackPeriodsInMonths

	// asNE(ti.KeyAssumptionsList, t.KeyAssumptionsList)
	if ne(ti.KeyAssumptionsList) {
		t.KeyAssumptionsList = &ti.KeyAssumptionsList
	}
	// asNE(ti.ProjectIdentifiedLegalRisks, t.ProjectIdentifiedLegalRisks)
	if ne(ti.ProjectIdentifiedLegalRisks) {
		t.ProjectIdentifiedLegalRisks = &ti.ProjectIdentifiedLegalRisks
	}
	// asNE(ti.ProjectIdentifiedRegulatoryRisks, t.ProjectIdentifiedRegulatoryRisks)
	if ne(ti.ProjectIdentifiedRegulatoryRisks) {
		t.ProjectIdentifiedRegulatoryRisks = &ti.ProjectIdentifiedRegulatoryRisks
	}

	// asNE(ti.ProjectIdentifiedOperationalOrExecutionRisks, t.ProjectIdentifiedOperationalOrExecutionRisks)
	if ne(ti.ProjectIdentifiedOperationalOrExecutionRisks) {
		t.ProjectIdentifiedOperationalOrExecutionRisks = &ti.ProjectIdentifiedOperationalOrExecutionRisks
	}
	// asNE(ti.ProjectIdentifiedMarketRisks, t.ProjectIdentifiedMarketRisks)
	if ne(ti.ProjectIdentifiedMarketRisks) {
		t.ProjectIdentifiedMarketRisks = &ti.ProjectIdentifiedMarketRisks
	}
	// asNE(ti.ProjectIdentifiedOtherRelevantRisks, t.ProjectIdentifiedOtherRelevantRisks)
	if ne(ti.ProjectIdentifiedOtherRelevantRisks) {
		t.ProjectIdentifiedOtherRelevantRisks = &ti.ProjectIdentifiedOtherRelevantRisks
	}
	if t.DeepLink == nil && t.AssetCode != nil && t.IssuingWalletPublicKey != nil {
		//set deep link
		p, _ := dynamiclinks.GenerateTokenizedAssetDeeplink(*t.AssetCode, *t.IssuingWalletPublicKey, gc)
		if len(p.DynamicLink) > 0 {
			t.DeepLink = &p.DynamicLink
		}
	}

	t.TotalIssueSize = ti.TotalIssueSize
	t.IssueDate = ti.IssueDate
	t.MaturityDate = ti.MaturityDate
	t.FaceValuePerUnit = ti.FaceValuePerUnit
	t.CouponOrInterestRate = ti.CouponOrInterestRate
	t.ExpectedYield = ti.ExpectedYield
	t.MinimumInvestmentAmount = ti.MinimumInvestmentAmount
	t.WeightedAverageLife = ti.WeightedAverageLife
	t.UnderlyingAssetPoolSize = ti.UnderlyingAssetPoolSize
	t.CouponRate = ti.CouponRate
	t.SpreadOrMargin = ti.SpreadOrMargin
	t.TotalExpenseRatio = ti.TotalExpenseRatio
	t.ExitLoadRedemptionFee = ti.ExitLoadRedemptionFee
	t.Tenure = ti.Tenure
	t.InitialNetAssetValue = ti.InitialNetAssetValue
	t.EntryLoad = ti.EntryLoad
	t.PerformanceFee = ti.PerformanceFee
	t.VolatilityEstimate = ti.VolatilityEstimate
	t.DividendYield = ti.DividendYield
	t.ExitLoadOrRedemptionFee = ti.ExitLoadOrRedemptionFee
	t.AverageMaturity = ti.AverageMaturity
	t.YieldToMaturity = ti.YieldToMaturity
	t.PerformanceFeeIfAny = ti.PerformanceFeeIfAny
	t.AverageMaturityHybrid = ti.AverageMaturityHybrid
	t.YieldToMaturityHybrid = ti.YieldToMaturityHybrid
	t.PricingYield = ti.PricingYield
	/////
	t.Quantity = ti.Quantity
	t.MinimumPurchaseAmount = ti.MinimumPurchaseAmount
	t.AssetValuation = ti.AssetValuation
	t.HoldingPeriod = ti.HoldingPeriod
	t.RisksCoverageValue = ti.RisksCoverageValue
	t.RecoveryEstimate = ti.RecoveryEstimate
	t.ExitLoadFee = ti.ExitLoadFee
	t.EntryLoadFee = ti.EntryLoadFee
	t.PortfolioPerformanceFee = ti.PortfolioPerformanceFee
	//////////

	if len(ti.IsinOrSerialNumber) > 0 {
		t.IsinOrSerialNumber = &ti.IsinOrSerialNumber
	} else {
		t.IsinOrSerialNumber = nil
	}

	if len(ti.InstrumentName) > 0 {
		t.InstrumentName = &ti.InstrumentName
	} else {
		t.InstrumentName = nil
	}

	if len(ti.InstrumentType) > 0 {
		t.InstrumentType = &ti.InstrumentType
	} else {
		t.InstrumentType = nil
	}

	if len(ti.CouponOrInterestRateType) > 0 {
		t.CouponOrInterestRateType = &ti.CouponOrInterestRateType
	} else {
		t.CouponOrInterestRateType = nil
	}

	if len(ti.ReferenceIndex) > 0 {
		t.ReferenceIndex = &ti.ReferenceIndex
	} else {
		t.ReferenceIndex = nil
	}

	if len(ti.EarlyRedemptionOptionInvestor) > 0 {
		t.EarlyRedemptionOptionInvestor = &ti.EarlyRedemptionOptionInvestor
	} else {
		t.EarlyRedemptionOptionInvestor = nil
	}

	if len(ti.TaxTreatmentTokenHolders) > 0 {
		t.TaxTreatmentTokenHolders = &ti.TaxTreatmentTokenHolders
	} else {
		t.TaxTreatmentTokenHolders = nil
	}

	if len(ti.PaymentStructureToTokenHolders) > 0 {
		t.PaymentStructureToTokenHolders = &ti.PaymentStructureToTokenHolders
	} else {
		t.PaymentStructureToTokenHolders = nil
	}

	if len(ti.RedemptionMethod) > 0 {
		t.RedemptionMethod = &ti.RedemptionMethod
	} else {
		t.RedemptionMethod = nil
	}

	if len(ti.PaymentStructure) > 0 {
		t.PaymentStructure = &ti.PaymentStructure
	} else {
		t.PaymentStructure = nil
	}

	if len(ti.RepaymentMethod) > 0 {
		t.RepaymentMethod = &ti.RepaymentMethod
	} else {
		t.RepaymentMethod = nil
	}

	if len(ti.PaymentCycle) > 0 {
		t.PaymentCycle = &ti.PaymentCycle
	} else {
		t.PaymentCycle = nil
	}

	if len(ti.PoolComposition) > 0 {
		t.PoolComposition = &ti.PoolComposition
	} else {
		t.PoolComposition = nil
	}

	if len(ti.CreditEnhancementMethod) > 0 {
		t.CreditEnhancementMethod = &ti.CreditEnhancementMethod
	} else {
		t.CreditEnhancementMethod = nil
	}

	if len(ti.SummaryOfUseOfProceeds) > 0 {
		t.SummaryOfUseOfProceeds = &ti.SummaryOfUseOfProceeds
	} else {
		t.SummaryOfUseOfProceeds = nil
	}

	if len(ti.CreditRatingIfAny) > 0 {
		t.CreditRatingIfAny = &ti.CreditRatingIfAny
	} else {
		t.CreditRatingIfAny = nil
	}

	if len(ti.IssuerName) > 0 {
		t.IssuerName = &ti.IssuerName
	} else {
		t.IssuerName = nil
	}

	if len(ti.IssuerType) > 0 {
		t.IssuerType = &ti.IssuerType
	} else {
		t.IssuerType = nil
	}

	if len(ti.IssuerContactPerson) > 0 {
		t.IssuerContactPerson = &ti.IssuerContactPerson
	} else {
		t.IssuerContactPerson = nil
	}

	if len(ti.ContactEmail) > 0 {
		t.ContactEmail = &ti.ContactEmail
	} else {
		t.ContactEmail = nil
	}

	if len(ti.ContactPhoneNumber) > 0 {
		t.ContactPhoneNumber = &ti.ContactPhoneNumber
	} else {
		t.ContactPhoneNumber = nil
	}

	if len(ti.BriefCompanyOverview) > 0 {
		t.BriefCompanyOverview = &ti.BriefCompanyOverview
	} else {
		t.BriefCompanyOverview = nil
	}

	if len(ti.MortgageOriginators) > 0 {
		t.MortgageOriginators = &ti.MortgageOriginators
	} else {
		t.MortgageOriginators = nil
	}

	if len(ti.Servicer) > 0 {
		t.Servicer = &ti.Servicer
	} else {
		t.Servicer = nil
	}

	if len(ti.InstrumentTrustee) > 0 {
		t.InstrumentTrustee = &ti.InstrumentTrustee
	} else {
		t.InstrumentTrustee = nil
	}

	if len(ti.InstrumentCustodians) > 0 {
		t.InstrumentCustodians = &ti.InstrumentCustodians
	} else {
		t.InstrumentCustodians = nil
	}

	if len(ti.InstrumentLegalAdvisor) > 0 {
		t.InstrumentLegalAdvisor = &ti.InstrumentLegalAdvisor
	} else {
		t.InstrumentLegalAdvisor = nil
	}

	if len(ti.SpecialPurposeVehicle) > 0 {
		t.SpecialPurposeVehicle = &ti.SpecialPurposeVehicle
	} else {
		t.SpecialPurposeVehicle = nil
	}

	if len(ti.InstrumentAssetManagerOrAdministrator) > 0 {
		t.InstrumentAssetManagerOrAdministrator = &ti.InstrumentAssetManagerOrAdministrator
	} else {
		t.InstrumentAssetManagerOrAdministrator = nil
	}

	if len(ti.UnderwriterIfAny) > 0 {
		t.UnderwriterIfAny = &ti.UnderwriterIfAny
	} else {
		t.UnderwriterIfAny = nil
	}

	if len(ti.InstrumentCreditRatingAgency) > 0 {
		t.InstrumentCreditRatingAgency = &ti.InstrumentCreditRatingAgency
	} else {
		t.InstrumentCreditRatingAgency = nil
	}

	if len(ti.AuditorOrVerifier) > 0 {
		t.AuditorOrVerifier = &ti.AuditorOrVerifier
	} else {
		t.AuditorOrVerifier = nil
	}

	if len(ti.CreditRiskAssessment) > 0 {
		t.CreditRiskAssessment = &ti.CreditRiskAssessment
	} else {
		t.CreditRiskAssessment = nil
	}

	if len(ti.CreditRating) > 0 {
		t.CreditRating = &ti.CreditRating
	} else {
		t.CreditRating = nil
	}

	if len(ti.PrepaymentRisk) > 0 {
		t.PrepaymentRisk = &ti.PrepaymentRisk
	} else {
		t.PrepaymentRisk = nil
	}

	if len(ti.InterestRateRisk) > 0 {
		t.InterestRateRisk = &ti.InterestRateRisk
	} else {
		t.InterestRateRisk = nil
	}

	if len(ti.StructuralComplexityRisk) > 0 {
		t.StructuralComplexityRisk = &ti.StructuralComplexityRisk
	} else {
		t.StructuralComplexityRisk = nil
	}

	if len(ti.LegalOrRegulatoryRisk) > 0 {
		t.LegalOrRegulatoryRisk = &ti.LegalOrRegulatoryRisk
	} else {
		t.LegalOrRegulatoryRisk = nil
	}

	if len(ti.OperationalRisk) > 0 {
		t.OperationalRisk = &ti.OperationalRisk
	} else {
		t.OperationalRisk = nil
	}

	if len(ti.MarketRisk) > 0 {
		t.MarketRisk = &ti.MarketRisk
	} else {
		t.MarketRisk = nil
	}

	if len(ti.EsgRisk) > 0 {
		t.EsgRisk = &ti.EsgRisk
	} else {
		t.EsgRisk = nil
	}

	if len(ti.MitigationMeasures) > 0 {
		t.MitigationMeasures = &ti.MitigationMeasures
	} else {
		t.MitigationMeasures = nil
	}

	if len(ti.IssuingAuthority) > 0 {
		t.IssuingAuthority = &ti.IssuingAuthority
	} else {
		t.IssuingAuthority = nil
	}

	if len(ti.RegulatoryApprovalId) > 0 {
		t.RegulatoryApprovalId = &ti.RegulatoryApprovalId
	} else {
		t.RegulatoryApprovalId = nil
	}

	if len(ti.LicenseApprovalReferenceNumber) > 0 {
		t.LicenseApprovalReferenceNumber = &ti.LicenseApprovalReferenceNumber
	} else {
		t.LicenseApprovalReferenceNumber = nil
	}

	if len(ti.ListingStatus) > 0 {
		t.ListingStatus = &ti.ListingStatus
	} else {
		t.ListingStatus = nil
	}

	if len(ti.CurrencyOfIssuance) > 0 {
		t.CurrencyOfIssuance = &ti.CurrencyOfIssuance
	} else {
		t.CurrencyOfIssuance = nil
	}

	if len(ti.CouponRateType) > 0 {
		t.CouponRateType = &ti.CouponRateType
	} else {
		t.CouponRateType = nil
	}

	if len(ti.ResetFrequency) > 0 {
		t.ResetFrequency = &ti.ResetFrequency
	} else {
		t.ResetFrequency = nil
	}

	if len(ti.CouponPaymentFrequency) > 0 {
		t.CouponPaymentFrequency = &ti.CouponPaymentFrequency
	} else {
		t.CouponPaymentFrequency = nil
	}

	if len(ti.RedemptionStructure) > 0 {
		t.RedemptionStructure = &ti.RedemptionStructure
	} else {
		t.RedemptionStructure = nil
	}

	if len(ti.EarlyRedemptionOption) > 0 {
		t.EarlyRedemptionOption = &ti.EarlyRedemptionOption
	} else {
		t.EarlyRedemptionOption = nil
	}

	if len(ti.EarlyRedemptionPenalty) > 0 {
		t.EarlyRedemptionPenalty = &ti.EarlyRedemptionPenalty
	} else {
		t.EarlyRedemptionPenalty = nil
	}

	if len(ti.TaxTreatment) > 0 {
		t.TaxTreatment = &ti.TaxTreatment
	} else {
		t.TaxTreatment = nil
	}

	if len(ti.NavOrMarketValueUpdates) > 0 {
		t.NavOrMarketValueUpdates = &ti.NavOrMarketValueUpdates
	} else {
		t.NavOrMarketValueUpdates = nil
	}

	if len(ti.ImpactMetrics) > 0 {
		t.ImpactMetrics = &ti.ImpactMetrics
	} else {
		t.ImpactMetrics = nil
	}

	if len(ti.LegalBacking) > 0 {
		t.LegalBacking = &ti.LegalBacking
	} else {
		t.LegalBacking = nil
	}

	if len(ti.DefaultHistory) > 0 {
		t.DefaultHistory = &ti.DefaultHistory
	} else {
		t.DefaultHistory = nil
	}

	if len(ti.RiskFactorsSummary) > 0 {
		t.RiskFactorsSummary = &ti.RiskFactorsSummary
	} else {
		t.RiskFactorsSummary = nil
	}

	if len(ti.PayingAgent) > 0 {
		t.PayingAgent = &ti.PayingAgent
	} else {
		t.PayingAgent = nil
	}

	if len(ti.Auditor) > 0 {
		t.Auditor = &ti.Auditor
	} else {
		t.Auditor = nil
	}

	if len(ti.RegistrarOrCSCSAgent) > 0 {
		t.RegistrarOrCSCSAgent = &ti.RegistrarOrCSCSAgent
	} else {
		t.RegistrarOrCSCSAgent = nil
	}

	if len(ti.FundStructure) > 0 {
		t.FundStructure = &ti.FundStructure
	} else {
		t.FundStructure = nil
	}

	if len(ti.AssetManagementCompanyName) > 0 {
		t.AssetManagementCompanyName = &ti.AssetManagementCompanyName
	} else {
		t.AssetManagementCompanyName = nil
	}

	if len(ti.FundManagers) > 0 {
		t.FundManagers = &ti.FundManagers
	} else {
		t.FundManagers = nil
	}

	if len(ti.RegulatoryLicenseNumber) > 0 {
		t.RegulatoryLicenseNumber = &ti.RegulatoryLicenseNumber
	} else {
		t.RegulatoryLicenseNumber = nil
	}

	t.FundLaunchDate = &ti.FundLaunchDate

	if len(ti.NavUpdateFrequency) > 0 {
		t.NavUpdateFrequency = &ti.NavUpdateFrequency
	} else {
		t.NavUpdateFrequency = nil
	}

	if len(ti.NavCalculationMethod) > 0 {
		t.NavCalculationMethod = &ti.NavCalculationMethod
	} else {
		t.NavCalculationMethod = nil
	}

	if len(ti.RedemptionRules) > 0 {
		t.RedemptionRules = &ti.RedemptionRules
	} else {
		t.RedemptionRules = nil
	}

	t.LockInPeriod = ti.LockInPeriod

	if len(ti.DividendPolicy) > 0 {
		t.DividendPolicy = &ti.DividendPolicy
	} else {
		t.DividendPolicy = nil
	}

	if len(ti.LiquidityProfile) > 0 {
		t.LiquidityProfile = &ti.LiquidityProfile
	} else {
		t.LiquidityProfile = nil
	}

	if len(ti.DistributionFrequency) > 0 {
		t.DistributionFrequency = &ti.DistributionFrequency
	} else {
		t.DistributionFrequency = nil
	}

	if len(ti.DistributionMethod) > 0 {
		t.DistributionMethod = &ti.DistributionMethod
	} else {
		t.DistributionMethod = nil
	}

	if len(ti.BenchmarkComparisonMethod) > 0 {
		t.BenchmarkComparisonMethod = &ti.BenchmarkComparisonMethod
	} else {
		t.BenchmarkComparisonMethod = nil
	}

	if len(ti.FeeBreakdownSummary) > 0 {
		t.FeeBreakdownSummary = &ti.FeeBreakdownSummary
	} else {
		t.FeeBreakdownSummary = nil
	}

	if len(ti.InvestmentObjective) > 0 {
		t.InvestmentObjective = &ti.InvestmentObjective
	} else {
		t.InvestmentObjective = nil
	}

	if len(ti.EquityStrategy) > 0 {
		t.EquityStrategy = &ti.EquityStrategy
	} else {
		t.EquityStrategy = nil
	}

	if len(ti.MarketCapitalizationFocus) > 0 {
		t.MarketCapitalizationFocus = &ti.MarketCapitalizationFocus
	} else {
		t.MarketCapitalizationFocus = nil
	}

	if len(ti.BenchmarkIndex) > 0 {
		t.BenchmarkIndex = &ti.BenchmarkIndex
	} else {
		t.BenchmarkIndex = nil
	}

	if len(ti.SectorExposureLimits) > 0 {
		t.SectorExposureLimits = &ti.SectorExposureLimits
	} else {
		t.SectorExposureLimits = nil
	}

	if len(ti.TopHoldings) > 0 {
		t.TopHoldings = &ti.TopHoldings
	} else {
		t.TopHoldings = nil
	}

	if len(ti.GeographicExposure) > 0 {
		t.GeographicExposure = &ti.GeographicExposure
	} else {
		t.GeographicExposure = nil
	}

	if len(ti.RiskProfile) > 0 {
		t.RiskProfile = &ti.RiskProfile
	} else {
		t.RiskProfile = nil
	}

	if len(ti.TrusteeName) > 0 {
		t.TrusteeName = &ti.TrusteeName
	} else {
		t.TrusteeName = nil
	}

	if len(ti.FundAdministrator) > 0 {
		t.FundAdministrator = &ti.FundAdministrator
	} else {
		t.FundAdministrator = nil
	}

	if len(ti.InvestmentCommitteeMembers) > 0 {
		t.InvestmentCommitteeMembers = &ti.InvestmentCommitteeMembers
	} else {
		t.InvestmentCommitteeMembers = nil
	}

	if len(ti.IsinOrSecFundCode) > 0 {
		t.IsinOrSecFundCode = &ti.IsinOrSecFundCode
	} else {
		t.IsinOrSecFundCode = nil
	}

	if len(ti.FundRiskRating) > 0 {
		t.FundRiskRating = &ti.FundRiskRating
	} else {
		t.FundRiskRating = nil
	}

	if len(ti.AssetAllocation) > 0 {
		t.AssetAllocation = &ti.AssetAllocation
	} else {
		t.AssetAllocation = nil
	}

	if len(ti.CreditRatingProfile) > 0 {
		t.CreditRatingProfile = &ti.CreditRatingProfile
	} else {
		t.CreditRatingProfile = nil
	}

	t.LockInPeriodPortfolio = ti.LockInPeriodPortfolio

	if len(ti.TopEquityHoldings) > 0 {
		t.TopEquityHoldings = &ti.TopEquityHoldings
	} else {
		t.TopEquityHoldings = nil
	}

	if len(ti.TargetAllocation) > 0 {
		t.TargetAllocation = &ti.TargetAllocation
	} else {
		t.TargetAllocation = nil
	}

	if len(ti.AllowedAllocationRange) > 0 {
		t.AllowedAllocationRange = &ti.AllowedAllocationRange
	} else {
		t.AllowedAllocationRange = nil
	}

	if len(ti.AssetClassesIncluded) > 0 {
		t.AssetClassesIncluded = &ti.AssetClassesIncluded
	} else {
		t.AssetClassesIncluded = nil
	}

	if len(ti.RebalancingFrequency) > 0 {
		t.RebalancingFrequency = &ti.RebalancingFrequency
	} else {
		t.RebalancingFrequency = nil
	}

	if len(ti.BenchmarkIndexComposite) > 0 {
		t.BenchmarkIndexComposite = &ti.BenchmarkIndexComposite
	} else {
		t.BenchmarkIndexComposite = nil
	}

	if len(ti.TopEquityHoldingsList) > 0 {
		t.TopEquityHoldingsList = &ti.TopEquityHoldingsList
	} else {
		t.TopEquityHoldingsList = nil
	}

	if len(ti.TopDebtHoldings) > 0 {
		t.TopDebtHoldings = &ti.TopDebtHoldings
	} else {
		t.TopDebtHoldings = nil
	}

	if len(ti.CreditRatingDistribution) > 0 {
		t.CreditRatingDistribution = &ti.CreditRatingDistribution
	} else {
		t.CreditRatingDistribution = nil
	}

	if len(ti.TitleOfIssuance) > 0 {
		t.TitleOfIssuance = &ti.TitleOfIssuance
	} else {
		t.TitleOfIssuance = nil
	}

	if len(ti.TypeOfCommercialPaper) > 0 {
		t.TypeOfCommercialPaper = &ti.TypeOfCommercialPaper
	} else {
		t.TypeOfCommercialPaper = nil
	}

	if len(ti.UseOfProceeds) > 0 {
		t.UseOfProceeds = &ti.UseOfProceeds
	} else {
		t.UseOfProceeds = nil
	}

	if len(ti.Ranking) > 0 {
		t.Ranking = &ti.Ranking
	} else {
		t.Ranking = nil
	}

	if len(ti.BackingSecurity) > 0 {
		t.BackingSecurity = &ti.BackingSecurity
	} else {
		t.BackingSecurity = nil
	}

	if len(ti.IssuerRegistrationNumber) > 0 {
		t.IssuerRegistrationNumber = &ti.IssuerRegistrationNumber
	} else {
		t.IssuerRegistrationNumber = nil
	}

	t.IncorporationDate = &ti.IncorporationDate

	if len(ti.RcNumber) > 0 {
		t.RcNumber = &ti.RcNumber
	} else {
		t.RcNumber = nil
	}

	if len(ti.TaxIdNumber) > 0 {
		t.TaxIdNumber = &ti.TaxIdNumber
	} else {
		t.TaxIdNumber = nil
	}

	if len(ti.OfficeAddress) > 0 {
		t.OfficeAddress = &ti.OfficeAddress
	} else {
		t.OfficeAddress = nil
	}

	if len(ti.Rating) > 0 {
		t.Rating = &ti.Rating
	} else {
		t.Rating = nil
	}

	////////

	if len(ti.PartiesInvolvedIssuer) > 0 {
		t.PartiesInvolvedIssuer = &ti.PartiesInvolvedIssuer
	} else {
		t.PartiesInvolvedIssuer = nil
	}

	if len(ti.PartiesInvolvedArranger) > 0 {
		t.PartiesInvolvedArranger = &ti.PartiesInvolvedArranger
	} else {
		t.PartiesInvolvedArranger = nil
	}

	if len(ti.PartiesInvolvedLegalAdviser) > 0 {
		t.PartiesInvolvedLegalAdviser = &ti.PartiesInvolvedLegalAdviser
	} else {
		t.PartiesInvolvedLegalAdviser = nil
	}

	if len(ti.PartiesInvolvedAuditor) > 0 {
		t.PartiesInvolvedAuditor = &ti.PartiesInvolvedAuditor
	} else {
		t.PartiesInvolvedAuditor = nil
	}

	if len(ti.PartiesInvolvedRatingAgency) > 0 {
		t.PartiesInvolvedRatingAgency = &ti.PartiesInvolvedRatingAgency
	} else {
		t.PartiesInvolvedRatingAgency = nil
	}

	if len(ti.PartiesInvolvedCustodian) > 0 {
		t.PartiesInvolvedCustodian = &ti.PartiesInvolvedCustodian
	} else {
		t.PartiesInvolvedCustodian = nil
	}

	if len(ti.PartiesInvolvedTrustee) > 0 {
		t.PartiesInvolvedTrustee = &ti.PartiesInvolvedTrustee
	} else {
		t.PartiesInvolvedTrustee = nil
	}

	if len(ti.PartiesInvolvedAuditorVerifier) > 0 {
		t.PartiesInvolvedAuditorVerifier = &ti.PartiesInvolvedAuditorVerifier
	} else {
		t.PartiesInvolvedAuditorVerifier = nil
	}

	if len(ti.SecurityRiskLegalBacking) > 0 {
		t.SecurityRiskLegalBacking = &ti.SecurityRiskLegalBacking
	} else {
		t.SecurityRiskLegalBacking = nil
	}

	if len(ti.SecurityRiskCollateral) > 0 {
		t.SecurityRiskCollateral = &ti.SecurityRiskCollateral
	} else {
		t.SecurityRiskCollateral = nil
	}

	if len(ti.SecurityRiskDefaultHistory) > 0 {
		t.SecurityRiskDefaultHistory = &ti.SecurityRiskDefaultHistory
	} else {
		t.SecurityRiskDefaultHistory = nil
	}

	if len(ti.SecurityRiskCreditRating) > 0 {
		t.SecurityRiskCreditRating = &ti.SecurityRiskCreditRating
	} else {
		t.SecurityRiskCreditRating = nil
	}

	if len(ti.SecurityRiskRiskFactorsSummary) > 0 {
		t.SecurityRiskRiskFactorsSummary = &ti.SecurityRiskRiskFactorsSummary
	} else {
		t.SecurityRiskRiskFactorsSummary = nil
	}

	if len(ti.SecurityRiskBusinessRisk) > 0 {
		t.SecurityRiskBusinessRisk = &ti.SecurityRiskBusinessRisk
	} else {
		t.SecurityRiskBusinessRisk = nil
	}

	if len(ti.SecurityRiskDefaultRisk) > 0 {
		t.SecurityRiskDefaultRisk = &ti.SecurityRiskDefaultRisk
	} else {
		t.SecurityRiskDefaultRisk = nil
	}

	if len(ti.SecurityRiskLiquidityRisk) > 0 {
		t.SecurityRiskLiquidityRisk = &ti.SecurityRiskLiquidityRisk
	} else {
		t.SecurityRiskLiquidityRisk = nil
	}

	if len(ti.SecurityRiskRegulatoryRisk) > 0 {
		t.SecurityRiskRegulatoryRisk = &ti.SecurityRiskRegulatoryRisk
	} else {
		t.SecurityRiskRegulatoryRisk = nil
	}

	if len(ti.SecurityRiskMarketRisk) > 0 {
		t.SecurityRiskMarketRisk = &ti.SecurityRiskMarketRisk
	} else {
		t.SecurityRiskMarketRisk = nil
	}

	if len(ti.SecurityRiskOperationalRisk) > 0 {
		t.SecurityRiskOperationalRisk = &ti.SecurityRiskOperationalRisk
	} else {
		t.SecurityRiskOperationalRisk = nil
	}

	if len(ti.SecurityRiskMitigationMeasures) > 0 {
		t.SecurityRiskMitigationMeasures = &ti.SecurityRiskMitigationMeasures
	} else {
		t.SecurityRiskMitigationMeasures = nil
	}

	if len(ti.CommodityType) > 0 {
		t.CommodityType = &ti.CommodityType
	} else {
		t.CommodityType = nil
	}

	if len(ti.CommodityDescription) > 0 {
		t.CommodityDescription = &ti.CommodityDescription
	} else {
		t.CommodityDescription = nil
	}

	if len(ti.QualityGrade) > 0 {
		t.QualityGrade = &ti.QualityGrade
	} else {
		t.QualityGrade = nil
	}

	if len(ti.IssuerContactInfo) > 0 {
		t.IssuerContactInfo = &ti.IssuerContactInfo
	} else {
		t.IssuerContactInfo = nil
	}

	if len(ti.WarehouseName) > 0 {
		t.WarehouseName = &ti.WarehouseName
	} else {
		t.WarehouseName = nil
	}

	if len(ti.WarehouseOperatorName) > 0 {
		t.WarehouseOperatorName = &ti.WarehouseOperatorName
	} else {
		t.WarehouseOperatorName = nil
	}

	if len(ti.WarehouseLicenseNumber) > 0 {
		t.WarehouseLicenseNumber = &ti.WarehouseLicenseNumber
	} else {
		t.WarehouseLicenseNumber = nil
	}

	if len(ti.WarehouseLocation) > 0 {
		t.WarehouseLocation = &ti.WarehouseLocation
	} else {
		t.WarehouseLocation = nil
	}

	if len(ti.WrNumber) > 0 {
		t.WrNumber = &ti.WrNumber
	} else {
		t.WrNumber = nil
	}

	t.WrIssueDate = &ti.WrIssueDate

	t.WrExpiryDate = &ti.WrExpiryDate

	if len(ti.WrSystemRegistration) > 0 {
		t.WrSystemRegistration = &ti.WrSystemRegistration
	} else {
		t.WrSystemRegistration = nil
	}

	if len(ti.WrRegistrationNumber) > 0 {
		t.WrRegistrationNumber = &ti.WrRegistrationNumber
	} else {
		t.WrRegistrationNumber = nil
	}

	if len(ti.WrVerifier) > 0 {
		t.WrVerifier = &ti.WrVerifier
	} else {
		t.WrVerifier = nil
	}

	if len(ti.StorageCondition) > 0 {
		t.StorageCondition = &ti.StorageCondition
	} else {
		t.StorageCondition = nil
	}

	if len(ti.WarehouseAccreditationBody) > 0 {
		t.WarehouseAccreditationBody = &ti.WarehouseAccreditationBody
	} else {
		t.WarehouseAccreditationBody = nil
	}

	if len(ti.AutoRollover) > 0 {
		t.AutoRollover = &ti.AutoRollover
	} else {
		t.AutoRollover = nil
	}

	if len(ti.CurrentBeneficialOwner) > 0 {
		t.CurrentBeneficialOwner = &ti.CurrentBeneficialOwner
	} else {
		t.CurrentBeneficialOwner = nil
	}

	if len(ti.WrCustodianName) > 0 {
		t.WrCustodianName = &ti.WrCustodianName
	} else {
		t.WrCustodianName = nil
	}

	if len(ti.OwnershipRightsRepresented) > 0 {
		t.OwnershipRightsRepresented = &ti.OwnershipRightsRepresented
	} else {
		t.OwnershipRightsRepresented = nil
	}

	if len(ti.TrusteeOrThirdPartyOversight) > 0 {
		t.TrusteeOrThirdPartyOversight = &ti.TrusteeOrThirdPartyOversight
	} else {
		t.TrusteeOrThirdPartyOversight = nil
	}

	if len(ti.LienOrEncumbrances) > 0 {
		t.LienOrEncumbrances = &ti.LienOrEncumbrances
	} else {
		t.LienOrEncumbrances = nil
	}

	t.ValuationDate = &ti.ValuationDate

	if len(ti.ValuationMethodology) > 0 {
		t.ValuationMethodology = &ti.ValuationMethodology
	} else {
		t.ValuationMethodology = nil
	}

	if len(ti.TokenizationObjective) > 0 {
		t.TokenizationObjective = &ti.TokenizationObjective
	} else {
		t.TokenizationObjective = nil
	}

	if len(ti.RedemptionMechanism) > 0 {
		t.RedemptionMechanism = &ti.RedemptionMechanism
	} else {
		t.RedemptionMechanism = nil
	}

	if len(ti.PartiesInvolvedUnderwriter) > 0 {
		t.PartiesInvolvedUnderwriter = &ti.PartiesInvolvedUnderwriter
	} else {
		t.PartiesInvolvedUnderwriter = nil
	}

	if len(ti.PartiesInvolvedAssetManager) > 0 {
		t.PartiesInvolvedAssetManager = &ti.PartiesInvolvedAssetManager
	} else {
		t.PartiesInvolvedAssetManager = nil
	}

	if len(ti.PartiesInvolvedLegalAdvisor) > 0 {
		t.PartiesInvolvedLegalAdvisor = &ti.PartiesInvolvedLegalAdvisor
	} else {
		t.PartiesInvolvedLegalAdvisor = nil
	}

	if len(ti.PartiesInvolvedRegulator) > 0 {
		t.PartiesInvolvedRegulator = &ti.PartiesInvolvedRegulator
	} else {
		t.PartiesInvolvedRegulator = nil
	}

	if len(ti.RisksMarketRisk) > 0 {
		t.RisksMarketRisk = &ti.RisksMarketRisk
	} else {
		t.RisksMarketRisk = nil
	}

	if len(ti.RisksStorageRisk) > 0 {
		t.RisksStorageRisk = &ti.RisksStorageRisk
	} else {
		t.RisksStorageRisk = nil
	}

	if len(ti.RisksTitleRisk) > 0 {
		t.RisksTitleRisk = &ti.RisksTitleRisk
	} else {
		t.RisksTitleRisk = nil
	}

	if len(ti.RisksFraudRisk) > 0 {
		t.RisksFraudRisk = &ti.RisksFraudRisk
	} else {
		t.RisksFraudRisk = nil
	}

	if len(ti.RisksInsuranceRisk) > 0 {
		t.RisksInsuranceRisk = &ti.RisksInsuranceRisk
	} else {
		t.RisksInsuranceRisk = nil
	}

	if len(ti.RisksOperationalRisk) > 0 {
		t.RisksOperationalRisk = &ti.RisksOperationalRisk
	} else {
		t.RisksOperationalRisk = nil
	}

	if len(ti.RisksRegulatoryRisk) > 0 {
		t.RisksRegulatoryRisk = &ti.RisksRegulatoryRisk
	} else {
		t.RisksRegulatoryRisk = nil
	}

	if len(ti.RisksLiquidityRisk) > 0 {
		t.RisksLiquidityRisk = &ti.RisksLiquidityRisk
	} else {
		t.RisksLiquidityRisk = nil
	}

	if len(ti.RisksForceMajeureRisk) > 0 {
		t.RisksForceMajeureRisk = &ti.RisksForceMajeureRisk
	} else {
		t.RisksForceMajeureRisk = nil
	}

	if len(ti.RisksEarlyRedemptionRisk) > 0 {
		t.RisksEarlyRedemptionRisk = &ti.RisksEarlyRedemptionRisk
	} else {
		t.RisksEarlyRedemptionRisk = nil
	}

	if len(ti.RisksMitigationMeasures) > 0 {
		t.RisksMitigationMeasures = &ti.RisksMitigationMeasures
	} else {
		t.RisksMitigationMeasures = nil
	}

	if len(ti.RisksInsuranceCoverageSummary) > 0 {
		t.RisksInsuranceCoverageSummary = &ti.RisksInsuranceCoverageSummary
	} else {
		t.RisksInsuranceCoverageSummary = nil
	}

	if len(ti.RisksInsuranceProvider) > 0 {
		t.RisksInsuranceProvider = &ti.RisksInsuranceProvider
	} else {
		t.RisksInsuranceProvider = nil
	}

	if len(ti.QuanlityStandard) > 0 {
		t.QuanlityStandard = &ti.QuanlityStandard
	} else {
		t.QuanlityStandard = nil
	}

	if len(ti.IssuerContactInformation) > 0 {
		t.IssuerContactInformation = &ti.IssuerContactInformation
	} else {
		t.IssuerContactInformation = nil
	}

	if len(ti.VaultCustodianName) > 0 {
		t.VaultCustodianName = &ti.VaultCustodianName
	} else {
		t.VaultCustodianName = nil
	}

	if len(ti.VaultOperator) > 0 {
		t.VaultOperator = &ti.VaultOperator
	} else {
		t.VaultOperator = nil
	}

	if len(ti.VaultLicenseNumber) > 0 {
		t.VaultLicenseNumber = &ti.VaultLicenseNumber
	} else {
		t.VaultLicenseNumber = nil
	}

	if len(ti.VaultLocation) > 0 {
		t.VaultLocation = &ti.VaultLocation
	} else {
		t.VaultLocation = nil
	}

	if len(ti.Number) > 0 {
		t.Number = &ti.Number
	} else {
		t.Number = nil
	}

	t.IssuerDate = &ti.IssuerDate

	t.ExpiryDate = &ti.ExpiryDate

	if len(ti.RegistryRecord) > 0 {
		t.RegistryRecord = &ti.RegistryRecord
	} else {
		t.RegistryRecord = nil
	}

	if len(ti.Verifier) > 0 {
		t.Verifier = &ti.Verifier
	} else {
		t.Verifier = nil
	}

	if len(ti.StorageConditions) > 0 {
		t.StorageConditions = &ti.StorageConditions
	} else {
		t.StorageConditions = nil
	}

	if len(ti.VaultAccreditationBody) > 0 {
		t.VaultAccreditationBody = &ti.VaultAccreditationBody
	} else {
		t.VaultAccreditationBody = nil
	}

	if len(ti.OwnershipLegalHolder) > 0 {
		t.OwnershipLegalHolder = &ti.OwnershipLegalHolder
	} else {
		t.OwnershipLegalHolder = nil
	}

	if len(ti.OwnershipCustodianName) > 0 {
		t.OwnershipCustodianName = &ti.OwnershipCustodianName
	} else {
		t.OwnershipCustodianName = nil
	}

	if len(ti.OwnershipTrustee) > 0 {
		t.OwnershipTrustee = &ti.OwnershipTrustee
	} else {
		t.OwnershipTrustee = nil
	}

	if len(ti.OwnershipLienOrEncumbrances) > 0 {
		t.OwnershipLienOrEncumbrances = &ti.OwnershipLienOrEncumbrances
	} else {
		t.OwnershipLienOrEncumbrances = nil
	}

	if len(ti.ValuationAssetValuation) > 0 {
		t.ValuationAssetValuation = &ti.ValuationAssetValuation
	} else {
		t.ValuationAssetValuation = nil
	}

	t.HoldingLockinPeriod = ti.HoldingLockinPeriod

	if len(ti.InsuranceMarketRisk) > 0 {
		t.InsuranceMarketRisk = &ti.InsuranceMarketRisk
	} else {
		t.InsuranceMarketRisk = nil
	}

	if len(ti.InsuranceStorageRisk) > 0 {
		t.InsuranceStorageRisk = &ti.InsuranceStorageRisk
	} else {
		t.InsuranceStorageRisk = nil
	}

	if len(ti.InsuranceTitleRisk) > 0 {
		t.InsuranceTitleRisk = &ti.InsuranceTitleRisk
	} else {
		t.InsuranceTitleRisk = nil
	}

	if len(ti.InsuranceFraudRisk) > 0 {
		t.InsuranceFraudRisk = &ti.InsuranceFraudRisk
	} else {
		t.InsuranceFraudRisk = nil
	}

	if len(ti.InsuranceInsuranceRisk) > 0 {
		t.InsuranceInsuranceRisk = &ti.InsuranceInsuranceRisk
	} else {
		t.InsuranceInsuranceRisk = nil
	}

	if len(ti.InsuranceOperationalRisk) > 0 {
		t.InsuranceOperationalRisk = &ti.InsuranceOperationalRisk
	} else {
		t.InsuranceOperationalRisk = nil
	}

	if len(ti.InsuranceRegulatoryRisk) > 0 {
		t.InsuranceRegulatoryRisk = &ti.InsuranceRegulatoryRisk
	} else {
		t.InsuranceRegulatoryRisk = nil
	}

	if len(ti.InsuranceLiquidityRisk) > 0 {
		t.InsuranceLiquidityRisk = &ti.InsuranceLiquidityRisk
	} else {
		t.InsuranceLiquidityRisk = nil
	}

	if len(ti.InsuranceForceMajeureRisk) > 0 {
		t.InsuranceForceMajeureRisk = &ti.InsuranceForceMajeureRisk
	} else {
		t.InsuranceForceMajeureRisk = nil
	}

	if len(ti.InsuranceEarlyRedemptionRisk) > 0 {
		t.InsuranceEarlyRedemptionRisk = &ti.InsuranceEarlyRedemptionRisk
	} else {
		t.InsuranceEarlyRedemptionRisk = nil
	}

	if len(ti.InsuranceMitigationMeasures) > 0 {
		t.InsuranceMitigationMeasures = &ti.InsuranceMitigationMeasures
	} else {
		t.InsuranceMitigationMeasures = nil
	}

	if len(ti.InsuranceInsuranceCoverageSummary) > 0 {
		t.InsuranceInsuranceCoverageSummary = &ti.InsuranceInsuranceCoverageSummary
	} else {
		t.InsuranceInsuranceCoverageSummary = nil
	}

	if len(ti.InsuranceInsuranceProvider) > 0 {
		t.InsuranceInsuranceProvider = &ti.InsuranceInsuranceProvider
	} else {
		t.InsuranceInsuranceProvider = nil
	}

	if len(ti.InsuranceCoverageValue) > 0 {
		t.InsuranceCoverageValue = &ti.InsuranceCoverageValue
	} else {
		t.InsuranceCoverageValue = nil
	}

	if len(ti.IssuerRegistrationNo) > 0 {
		t.IssuerRegistrationNo = &ti.IssuerRegistrationNo
	} else {
		t.IssuerRegistrationNo = nil
	}

	if len(ti.SectorAndIndustry) > 0 {
		t.SectorAndIndustry = &ti.SectorAndIndustry
	} else {
		t.SectorAndIndustry = nil
	}

	if len(ti.LicenseOrPermitNumber) > 0 {
		t.LicenseOrPermitNumber = &ti.LicenseOrPermitNumber
	} else {
		t.LicenseOrPermitNumber = nil
	}

	if len(ti.IssuerAdditionalInfo) > 0 {
		t.IssuerAdditionalInfo = &ti.IssuerAdditionalInfo
	} else {
		t.IssuerAdditionalInfo = nil
	}

	if len(ti.IsinSerialNumber) > 0 {
		t.IsinSerialNumber = &ti.IsinSerialNumber
	} else {
		t.IsinSerialNumber = nil
	}

	if len(ti.EsgOrImpactMetrics) > 0 {
		t.EsgOrImpactMetrics = &ti.EsgOrImpactMetrics
	} else {
		t.EsgOrImpactMetrics = nil
	}

	if len(ti.InstrumentAdditionalInfo) > 0 {
		t.InstrumentAdditionalInfo = &ti.InstrumentAdditionalInfo
	} else {
		t.InstrumentAdditionalInfo = nil
	}

	if len(ti.SecurityType) > 0 {
		t.SecurityType = &ti.SecurityType
	} else {
		t.SecurityType = nil
	}

	if len(ti.CollateralDescription) > 0 {
		t.CollateralDescription = &ti.CollateralDescription
	} else {
		t.CollateralDescription = nil
	}

	if len(ti.CovenantSummary) > 0 {
		t.CovenantSummary = &ti.CovenantSummary
	} else {
		t.CovenantSummary = nil
	}

	if len(ti.CovenantTestingFrequency) > 0 {
		t.CovenantTestingFrequency = &ti.CovenantTestingFrequency
	} else {
		t.CovenantTestingFrequency = nil
	}

	if len(ti.EventOfDefaultClauses) > 0 {
		t.EventOfDefaultClauses = &ti.EventOfDefaultClauses
	} else {
		t.EventOfDefaultClauses = nil
	}

	if len(ti.LegalEnforcementMechanism) > 0 {
		t.LegalEnforcementMechanism = &ti.LegalEnforcementMechanism
	} else {
		t.LegalEnforcementMechanism = nil
	}

	if len(ti.Guarantee) > 0 {
		t.Guarantee = &ti.Guarantee
	} else {
		t.Guarantee = nil
	}

	if len(ti.RiskProfileAdditionalInfo) > 0 {
		t.RiskProfileAdditionalInfo = &ti.RiskProfileAdditionalInfo
	} else {
		t.RiskProfileAdditionalInfo = nil
	}

	if len(ti.LicenseNumber) > 0 {
		t.LicenseNumber = &ti.LicenseNumber
	} else {
		t.LicenseNumber = nil
	}

	t.PortfolioLockInPeriod = ti.PortfolioLockInPeriod

	t.FundingStructure = ti.FundingStructure
	t.EquityPercentage = ti.EquityPercentage
	t.DebtPercentage = ti.DebtPercentage
	if len(ti.DebtInstrumentType) > 0 {
		t.DebtInstrumentType = &ti.DebtInstrumentType
	} else {
		t.DebtInstrumentType = nil
	}

	if len(ti.PrincipalPaymentMethod) > 0 {
		t.PrincipalPaymentMethod = &ti.PrincipalPaymentMethod
	} else {
		t.PrincipalPaymentMethod = nil
	}

	if len(ti.DebtInstrumentRepaymentSource) > 0 {
		t.DebtInstrumentRepaymentSource = &ti.DebtInstrumentRepaymentSource
	} else {
		t.DebtInstrumentRepaymentSource = nil
	}

	if len(ti.DebtInstrumentGuaranteesOrEnhancements) > 0 {
		t.DebtInstrumentGuaranteesOrEnhancements = &ti.DebtInstrumentGuaranteesOrEnhancements
	} else {
		t.DebtInstrumentGuaranteesOrEnhancements = nil
	}

	if len(ti.DebtInstrumentDefaultAndRecoveryTerms) > 0 {
		t.DebtInstrumentDefaultAndRecoveryTerms = &ti.DebtInstrumentDefaultAndRecoveryTerms
	} else {
		t.DebtInstrumentDefaultAndRecoveryTerms = nil
	}

	if len(ti.DebtInstrumentRepaymentFrequency) > 0 {
		t.DebtInstrumentRepaymentFrequency = &ti.DebtInstrumentRepaymentFrequency
	} else {
		t.DebtInstrumentRepaymentFrequency = nil
	}

	if len(ti.InterestRepaymentFrequency) > 0 {
		t.InterestRepaymentFrequency = &ti.InterestRepaymentFrequency
	} else {
		t.InterestRepaymentFrequency = nil
	}

	if len(ti.DcsrDetails) > 0 {
		t.DcsrDetails = &ti.DcsrDetails
	} else {
		t.DcsrDetails = nil
	}

	if len(ti.SinkingFundStructure) > 0 {
		t.SinkingFundStructure = &ti.SinkingFundStructure
	} else {
		t.SinkingFundStructure = nil
	}

	if len(ti.CovenantMonitoringAgent) > 0 {
		t.CovenantMonitoringAgent = &ti.CovenantMonitoringAgent
	} else {
		t.CovenantMonitoringAgent = nil
	}

	if len(ti.RightOfRecourse) > 0 {
		t.RightOfRecourse = &ti.RightOfRecourse
	} else {
		t.RightOfRecourse = nil
	}
	t.DebtInstrumentInterestRate = ti.DebtInstrumentInterestRate
	t.DcsrRatio = ti.DcsrRatio
	t.LtvRatio = ti.LtvRatio
	t.InterestCoverageRatio = ti.InterestCoverageRatio
	t.MaximumLeverageRatio = ti.MaximumLeverageRatio
	t.GracePeriod = ti.GracePeriod
	t.TrusteeAppointed = ti.TrusteeAppointed
	t.ReserveFundInPlace = ti.ReserveFundInPlace

	if len(ti.SecurityOrCollateralOffered) > 0 {
		t.SecurityOrCollateralOffered = &ti.SecurityOrCollateralOffered
	} else {
		t.SecurityOrCollateralOffered = nil
	}
	if len(ti.FundInstrumentType) > 0 {
		t.FundInstrumentType = &ti.FundInstrumentType
	} else {
		t.FundInstrumentType = nil
	}
	if len(ti.InstrumentRatingAgency) > 0 {
		t.InstrumentRatingAgency = &ti.InstrumentRatingAgency
	} else {
		t.InstrumentRatingAgency = nil
	}
	if len(ti.PortfolioTopHoldings) > 0 {
		t.PortfolioTopHoldings = &ti.PortfolioTopHoldings
	} else {
		t.PortfolioTopHoldings = nil
	}

	if len(ti.CreditRatingAgency) > 0 {
		t.CreditRatingAgency = &ti.CreditRatingAgency
	} else {
		t.CreditRatingAgency = nil
	}

	if len(ti.TrusteeRegNumber) > 0 {
		t.TrusteeRegNumber = &ti.TrusteeRegNumber
	} else {
		t.TrusteeRegNumber = nil
	}

	if len(ti.CreditEnhancerOrGuarantor) > 0 {
		t.CreditEnhancerOrGuarantor = &ti.CreditEnhancerOrGuarantor
	} else {
		t.CreditEnhancerOrGuarantor = nil
	}

	if len(ti.BondStructuringAdvisor) > 0 {
		t.BondStructuringAdvisor = &ti.BondStructuringAdvisor
	} else {
		t.BondStructuringAdvisor = nil
	}

	if len(ti.EntitiesAdditionalInfo) > 0 {
		t.EntitiesAdditionalInfo = &ti.EntitiesAdditionalInfo
	} else {
		t.EntitiesAdditionalInfo = nil
	}

	if len(ti.AuthorizedRepresentativeName) > 0 {
		t.AuthorizedRepresentativeName = &ti.AuthorizedRepresentativeName
	} else {
		t.AuthorizedRepresentativeName = nil
	}

	if len(ti.AuthorizedRepresentativeTitleOrPosition) > 0 {
		t.AuthorizedRepresentativeTitleOrPosition = &ti.AuthorizedRepresentativeTitleOrPosition
	} else {
		t.AuthorizedRepresentativeTitleOrPosition = nil
	}

	if len(ti.AuthorizedRepresentativeEmail) > 0 {
		t.AuthorizedRepresentativeEmail = &ti.AuthorizedRepresentativeEmail
	} else {
		t.AuthorizedRepresentativeEmail = nil
	}

	t.AcceptTokenizationTermsAndAgreement = ti.AcceptTokenizationTermsAndAgreement
	t.AttestInformationAccurateAndVerifiable = ti.AttestInformationAccurateAndVerifiable
	t.AcknowledgedSuitabilityCriteria = ti.AcknowledgedSuitabilityCriteria

	if len(ti.MinimumKycTier) > 0 {
		t.MinimumKycTier = &ti.MinimumKycTier
	} else {
		t.MinimumKycTier = nil
	}

	if len(ti.InvestorCategory) > 0 {
		t.InvestorCategory = &ti.InvestorCategory
	} else {
		t.InvestorCategory = nil
	}

	t.WithholdingTaxDisclosure = ti.WithholdingTaxDisclosure

	//////////

	return *t

}

func (t *TokenizedAsset) UpdateCalculation(gc *sharedconfig.GlobalConfig) {
	//TODO: set the SEC fee, Custody fee, Asset manager fee and recover feeInFiat from total asset value
	if t.AssetTokenizationStatus > 2 || t.AssetCurrentValue == 0 {
		//do not modify when it is already passed stage for fees or if currenct value is not set
		return
	}
	var feeCompo TokenizationFee
	var feeInAsset float64

	if *t.TokenizationFeeID > 0 {
		// fee has been selected
		feeCompo = t.UpdateTokenizationFeeByID(*t.TokenizationFeeID, gc)

		// Calculate Fees

		feeInAsset = decimal.NewFromFloat(t.NumberOfTokenToBeIssued * (feeCompo.FeeAssetPercentage / 100)).Truncate(7).InexactFloat64()
		t.FeeInAsset = feeInAsset
		log.Printf("[UpdateCalculation] Calculated Fee In Asset:= %v\n", decimal.NewFromFloat(feeInAsset).String())

		t.FeeInAssetPercent = feeCompo.FeeAssetPercentage
		t.FeeInFiat = decimal.NewFromFloat(t.AssetCurrentValue * (feeCompo.FeeFiatPercentage / 100)).Truncate(7).InexactFloat64()
		if feeCompo.FeeFiatCap > t.FeeInFiat {
			t.FeeInFiat = feeCompo.FeeFiatCap
		}
		log.Printf("[UpdateCalculation] Calculated Fee In Fiat:= %v\n", decimal.NewFromFloat(t.FeeInFiat).String())

	}
	// get SEC tokenization fee.
	// /////

	/**
	Let Initial value of Tokenized Asset before considering fee in asset be
	t.InitialValueOfTokenizedAsset - *this is not to be shown on the App*

	and let Final Value of Tokenized Asset  considering the fee in asset be.
	t.FinalValueOfTokenizedAsset - *This is what is shown on the App*


	t.InitialValueOfTokenizedAsset = (t.AssetCurrentValue + t.AssetMscCostOutisdeOfValuation + totalChargedFeesForVat + vat) - *this is not to be shown on the App*

	t.PricePerToken =(t.InitialValueOfTokenizedAsset / (t.NumberOfTokenToBeIssued – feeInAsset-VATAsset))

	t.TotalTokenHeldByManager = (t.AssetOwnerRetainedOrContributedValue / t.PricePerToken) -*This is Token Not for Sale*

	maxTokenToBeSold := (t.NumberOfTokenToBeIssued - feeInAsset - t.TotalTokenHeldByManager)

	t.FinalValueOfTokenizedAsset = (t.NumberOfTokenToBeIssued * t.PricePerToken) – *This is what is shown on the App*
		**/
	var cConfig Country
	var custodyFee, assetMgtFee float64
	var secFee float64
	var vat, issuingHouseFeeValue, legalAndProfessionalFee, ratingAgencyFee, totalChargedFeesForVat, trusteeFee float64
	if t.AssetCountryLocation != nil {
		cConfig = CountryCode(*t.AssetCountryLocation).GetConfig(gc)

	}
	VATAsset := decimal.NewFromFloat(feeInAsset * (cConfig.VATPercent / 100)).Truncate(7).InexactFloat64()
	t.VATInAsset = VATAsset
	log.Printf("[UpdateCalculation] Calculated VAT Asset:= %v\n", decimal.NewFromFloat(VATAsset).String())

	if t.AssetCurrentValue > 0 {
		if t.ExcludeSecFee == 0 {
			t.SECTokenizationFeePercent = cConfig.SECTokenizationFeePercent
			t.SECTokenizationFeeFixed = cConfig.SECTokenizationFeeFixed
			secFee = (decimal.NewFromFloat(t.AssetCurrentValue * (cConfig.SECTokenizationFeePercent / 100)).Truncate(7)).Add(decimal.NewFromFloat(cConfig.SECTokenizationFeeFixed).Truncate(7)).InexactFloat64()
			t.SECTokenizationFeeValue = secFee

		}

		log.Printf("[UpdateCalculation] Calculated SEC Fee Value:= %v\n", decimal.NewFromFloat(secFee).String())

		{
			// custodian := ApprovedCustodianID(t.ApprovedAssetCustodianID).GetApprovedCustodian(gc)
			// t.CustodianFeeFixed = custodian.FeeFixed
			// t.CustodianFeePercent = custodian.FeePercent
			custodyFee = (decimal.NewFromFloat(t.AssetCurrentValue * (t.CustodianFeePercent / 100)).Truncate(7)).Add(decimal.NewFromFloat(t.CustodianFeeFixed).Truncate(7)).InexactFloat64()
			t.CustodianFeeValue = custodyFee
			log.Printf("[UpdateCalculation] Custodian Fee Value:= %v\n", decimal.NewFromFloat(custodyFee).String())

		}

		{
			// assetManager := AssetManagerID(t.AssetManagerID).GetAssetManager(gc)
			// t.AssetManagerFeePercent = assetManager.FeePercent
			// t.AssetManagerFeeFixed = assetManager.FeeFixed
			assetMgtFee = (decimal.NewFromFloat(t.AssetCurrentValue * (t.AssetManagerFeePercent / 100)).Truncate(7)).Add(decimal.NewFromFloat(t.AssetManagerFeeFixed).Truncate(7)).InexactFloat64()
			t.AssetManagerFeeValue = assetMgtFee
			log.Printf("[UpdateCalculation] Calculated Asset Managment fee Value:= %v\n", decimal.NewFromFloat(assetMgtFee).String())

		}

		{
			// issuingHouse := IssuingHouseID(t.AssetIssuingHouseID).GetAssetIssuingHouse(gc)
			// t.IssuingHouseFeeFixed = issuingHouse.FeeFixed
			// t.IssuingHouseFeePercent = issuingHouse.FeePercent
			issuingHouseFeeValue = (decimal.NewFromFloat(t.AssetCurrentValue * (t.IssuingHouseFeePercent / 100)).Truncate(7)).Add(decimal.NewFromFloat(t.IssuingHouseFeeFixed).Truncate(7)).InexactFloat64()
			t.IssuingHouseFeeValue = issuingHouseFeeValue
			log.Printf("[UpdateCalculation] Calculated Issuing House Fee Value:= %v\n", decimal.NewFromFloat(issuingHouseFeeValue).String())

		}

		{
			// lpp := LegalAndProfesionalPartnerID(t.LegalAndProfesionalPartnerID).GetLegalAndProfesionalPartner(gc)
			// t.LegalAndProfessionalFeeFixed = lpp.FeeFixed
			// t.LegalAndProfessionalFeePercent = lpp.FeePercent
			legalAndProfessionalFee = (decimal.NewFromFloat(t.AssetCurrentValue * (t.LegalAndProfessionalFeePercent / 100)).Truncate(7)).Add(decimal.NewFromFloat(t.LegalAndProfessionalFeeFixed).Truncate(7)).InexactFloat64()
			t.LegalAndProfessionalFeeValue = legalAndProfessionalFee
			log.Printf("[UpdateCalculation] Calculated Legal&Professional Fee Value:= %v\n", decimal.NewFromFloat(legalAndProfessionalFee).String())

		}

		{
			// ra := RatingAgencyID(t.RatingAgencyID).GetRatingAgency(gc)
			// t.RatingAgencyFeeFixed = ra.FeeFixed
			// t.RatingAgencyFeePercent = ra.FeePercent
			ratingAgencyFee = (decimal.NewFromFloat(t.AssetCurrentValue * (t.RatingAgencyFeePercent / 100)).Truncate(7)).Add(decimal.NewFromFloat(t.RatingAgencyFeeFixed).Truncate(7)).InexactFloat64()
			t.RatingAgencyFeeValue = ratingAgencyFee
			log.Printf("[UpdateCalculation] Calculated Rating Agency Fee Value:= %v\n", decimal.NewFromFloat(ratingAgencyFee).String())

		}

		{
			// ra := TrusteeID(t.TrusteeID).GetTrustee(gc)
			// t.TrusteeFeeFixed = ra.FeeFixed
			// t.TrusteeFeePercent = ra.FeePercent
			trusteeFee = (decimal.NewFromFloat(t.AssetCurrentValue * (t.TrusteeFeePercent / 100)).Truncate(7)).Add(decimal.NewFromFloat(t.TrusteeFeeFixed).Truncate(7)).InexactFloat64()
			t.TrusteeFeeValue = trusteeFee
			log.Printf("[UpdateCalculation] Calculated Trustee Fee Value:= %v\n", decimal.NewFromFloat(trusteeFee).String())

		}

		totalChargedFeesForVat = secFee + custodyFee + assetMgtFee + t.FeeInFiat + issuingHouseFeeValue + legalAndProfessionalFee + ratingAgencyFee + trusteeFee
		log.Printf("[UpdateCalculation] Calculated Total charged Fees for VAT:= %v\n", decimal.NewFromFloat(totalChargedFeesForVat).String())

		vat = decimal.NewFromFloat(totalChargedFeesForVat * (cConfig.VATPercent / 100)).Truncate(7).InexactFloat64()

		t.VATPercent = cConfig.VATPercent
		t.VATValue = vat
		log.Printf("[UpdateCalculation] Calculated VAT Value:= %v\n", decimal.NewFromFloat(vat).String())

	}

	initialValueOfTokenizedAsset := (t.AssetCurrentValue + t.AssetMscCostOutisdeOfValuation + totalChargedFeesForVat + vat) // - *this is not to be shown on the App*
	t.InitialValueOfTokenizedAsset = initialValueOfTokenizedAsset
	log.Printf("[UpdateCalculation] Calculated Initial Value Of Tokenized Asset:= %v\n", decimal.NewFromFloat(initialValueOfTokenizedAsset).String())
	if t.NumberOfTokenToBeIssued > 0 {
		// t.PricePerToken = decimal.NewFromFloat(t.ValueOfTokenizedAsset / t.NumberOfTokenToBeIssued).Truncate(7).InexactFloat64()
		pricePerToken := decimal.NewFromFloat((initialValueOfTokenizedAsset / (t.NumberOfTokenToBeIssued - feeInAsset - VATAsset)))
		t.PricePerToken = pricePerToken.Round(7).InexactFloat64()
		log.Printf("[UpdateCalculation] Calculated PricePerToken:= %v, Rounded 7DP Value:= %v\n", pricePerToken.String(), decimal.NewFromFloat(t.PricePerToken).String())
		// auto calculate, token to be held is less the fee. token not to be sold
		totalTokenHeldByManager := decimal.NewFromFloat(t.AssetOwnerRetainedOrContributedValue / t.PricePerToken)
		t.TotalTokenHeldByManager = totalTokenHeldByManager.Truncate(7).InexactFloat64()
		log.Printf("[UpdateCalculation] Calculated Total Token Held By Manager:= %v, Truncated 7DP Value:= %v\n", totalTokenHeldByManager.String(), decimal.NewFromFloat(t.TotalTokenHeldByManager).String())

		// feeInAsset = decimal.NewFromFloat(feeInAssetFiatEquivalent / t.PricePerToken).Truncate(7).InexactFloat64()
		t.FeeInAsset = feeInAsset

		{
			//ensure correct the number of token to be sold.
			// maxTokenToBeSold := decimal.NewFromFloat(t.NumberOfTokenToBeIssued - feeInAsset - t.TotalTokenHeldByManager)
			maxTokenToBeSold := decimal.NewFromFloat(t.NumberOfTokenToBeIssued - feeInAsset - t.TotalTokenHeldByManager - VATAsset)

			t.MaxNumberOfTokenAvailableForSale = maxTokenToBeSold.Truncate(7).InexactFloat64()
			t.NumberOfTokenToBeSold = t.MaxNumberOfTokenAvailableForSale
			log.Printf("[UpdateCalculation] Calculated Max Token to be sold:= %v, Truncated 7DP Value:= %v\n", maxTokenToBeSold.String(), decimal.NewFromFloat(t.NumberOfTokenToBeSold).String())

		}

		finalValueOfTokenizedAsset := decimal.NewFromFloat(t.NumberOfTokenToBeIssued * t.PricePerToken)
		t.ValueOfTokenizedAsset = finalValueOfTokenizedAsset.Truncate(7).InexactFloat64()
		log.Printf("[UpdateCalculation] Calculated Final Value of tokenized Asset:= %v, TRUNCATED 7DP value := %v\n", finalValueOfTokenizedAsset.String(), decimal.NewFromFloat(t.ValueOfTokenizedAsset).String())
	}

	if t.CapOnPurchase == 1 && t.PricePerToken > 0 && t.CapAmountInFiat > 0 {
		t.CapQuantity = decimal.NewFromFloat(t.CapAmountInFiat / t.PricePerToken).Truncate(7).InexactFloat64()
	}
}

func (ti *TokenizedAsset) GetMintingInitiators(gc *sharedconfig.GlobalConfig) (us []TokenizationMintingInitiator) {
	us = make([]TokenizationMintingInitiator, 0)

	gc.DB.Find(&us)
	return
}

func (ti *TokenizedAsset) GetMintingInitiatorsInCSV(gc *sharedconfig.GlobalConfig) (csvStr string) {
	us := make([]TokenizationMintingInitiator, 0)
	replacer := strings.NewReplacer("\r", "", "\n", "", " ", "")
	gc.DB.Find(&us)
	trimmed := replacer.Replace(TokenizationMintingInitiators(us).ToCSV())
	csvStr = trimmed
	return
}

func (ti *TokenizedAsset) GetMintingApprovers(gc *sharedconfig.GlobalConfig) (us []TokenizationMintingApprover) {
	us = make([]TokenizationMintingApprover, 0)

	gc.DB.Find(&us)

	return
}

func (ti *TokenizedAsset) GetMintingApproversInCSV(gc *sharedconfig.GlobalConfig) (csvStr string) {
	us := make([]TokenizationMintingApprover, 0)
	replacer := strings.NewReplacer("\r", "", "\n", "", " ", "")
	gc.DB.Find(&us)
	trimmed := replacer.Replace(TokenizationMintingApprovers(us).ToCSV())
	csvStr = trimmed
	return
}

func (ti *TokenizedAsset) ToJSON(gc *sharedconfig.GlobalConfig) (t TokenizedAssetJSON) {
	titleCaser := cases.Title(language.English)
	replacer := strings.NewReplacer("\r", "", "\n", "", " ", "")
	t.ID = ti.ID
	t.CreatedAt = ti.CreatedAt
	t.UpdatedAt = ti.UpdatedAt
	t.InitiatorUsername = ti.InitiatorUsername
	t.FeeInAsset = ti.FeeInAsset
	t.FeeInFiat = ti.FeeInFiat
	t.FeeInAssetPercent = ti.FeeInAssetPercent
	t.SECTokenizationFeePercent = ti.SECTokenizationFeePercent
	t.SECTokenizationFeeFixed = ti.SECTokenizationFeeFixed
	t.SECTokenizationFeeValue = ti.SECTokenizationFeeValue
	t.AssetManagerID = ti.AssetManagerID
	t.AssetManagerFeePercent = ti.AssetManagerFeePercent
	t.AssetManagerFeeFixed = ti.AssetManagerFeeFixed
	t.AssetManagerFeeValue = ti.AssetManagerFeeValue
	t.IssuingHouseFeeFixed = ti.IssuingHouseFeeFixed
	t.IssuingHouseFeePercent = ti.IssuingHouseFeePercent
	t.IssuingHouseFeeValue = ti.IssuingHouseFeeValue
	t.ApprovedAssetCustodianID = ti.ApprovedAssetCustodianID
	t.CustodianFeePercent = ti.CustodianFeePercent
	t.CustodianFeeFixed = ti.CustodianFeeFixed
	t.CustodianFeeValue = ti.CustodianFeeValue
	t.LegalAndProfesionalPartnerID = ti.LegalAndProfesionalPartnerID
	t.LegalAndProfessionalFeePercent = ti.LegalAndProfessionalFeePercent
	t.LegalAndProfessionalFeeFixed = ti.LegalAndProfessionalFeeFixed
	t.LegalAndProfessionalFeeValue = ti.LegalAndProfessionalFeeValue
	t.RatingAgencyID = ti.RatingAgencyID
	t.RatingAgencyFeePercent = ti.RatingAgencyFeePercent
	t.RatingAgencyFeeFixed = ti.RatingAgencyFeeFixed
	t.RatingAgencyFeeValue = ti.RatingAgencyFeeValue
	t.TrusteeID = ti.TrusteeID
	t.TrusteeFeePercent = ti.TrusteeFeePercent
	t.TrusteeFeeFixed = ti.TrusteeFeeFixed
	t.TrusteeFeeValue = ti.TrusteeFeeValue
	t.VATPercent = ti.VATPercent
	t.VATValue = ti.VATValue
	t.VATInAsset = ti.VATInAsset

	t.VettingStatus = ti.VettingStatus
	t.DueDiligenceFail = ti.DueDiligenceFail

	if ti.DueDiligenceFailureReason != nil {

		t.DueDiligenceFailureReason = *ti.DueDiligenceFailureReason

	}

	t.AssetOwnerRetainedOrContributedValue = ti.AssetOwnerRetainedOrContributedValue
	log.Printf("[TokenizedAsset:ToJSON]AssetOwnerRetainedOrContributedValue[%v]:= %v\n", t.ID, decimal.NewFromFloat(ti.AssetOwnerRetainedOrContributedValue).String())

	if ti.InitialOwnerPreferredWalletAddress != nil {

		t.InitialOwnerPreferredWalletAddress = *ti.InitialOwnerPreferredWalletAddress

	}

	if ti.AdditionalKYCRequirements != nil {

		t.AdditionalKYCRequirements = *ti.AdditionalKYCRequirements

	}
	if ti.AssetSector != nil {

		t.AssetSector = *ti.AssetSector
	}

	if ti.AssetSubSector != nil {

		t.AssetSubSector = *ti.AssetSubSector
	}

	if ti.MintingApprovers != nil {
		trimmed := replacer.Replace(*ti.MintingApprovers)
		t.MintingApprovers = trimmed
	} else {
		//get default minting approvers.
		t.MintingApprovers = ti.GetMintingApproversInCSV(gc)
	}

	if ti.MintingInitators != nil {

		trimmed := replacer.Replace(*ti.MintingInitators)
		t.MintingInitators = trimmed
	} else {
		//get default minting approvers.
		t.MintingInitators = ti.GetMintingInitiatorsInCSV(gc)
	}

	if ti.AssetType != nil {

		t.AssetType = *ti.AssetType
	}
	if ti.AssetName != nil {

		t.AssetName = titleCaser.String(*ti.AssetName)

	}
	if ti.AssetWebsite != nil {
		t.AssetWebsite = strings.ToLower(*ti.AssetWebsite)

	}

	if ti.AssetCode != nil {
		t.AssetCode = strings.ToUpper(*ti.AssetCode)

	}
	if ti.AssetLogo != nil {
		t.AssetLogo = *ti.AssetLogo

	}

	if ti.ApprovedAssetCustodianID > 0 {
		t.ApprovedAssetCustodianID = ti.ApprovedAssetCustodianID
		t.ApprovedAssetCustodian = ti.ApprovedAssetCustodian
	}

	if ti.AssetIssuingHouseID > 0 {
		t.AssetIssuingHouseID = ti.AssetIssuingHouseID
		t.AssetIssuingHouse = ti.AssetIssuingHouse
	}

	if ti.LegalAndProfesionalPartnerID > 0 {
		t.LegalAndProfesionalPartnerID = ti.LegalAndProfesionalPartnerID
		t.LegalAndProfesionalPartner = ti.LegalAndProfesionalPartner
	}

	if ti.RatingAgencyID > 0 {
		t.RatingAgencyID = ti.RatingAgencyID
		t.RatingAgency = ti.RatingAgency
	}
	if ti.TrusteeID > 0 {
		t.TrusteeID = ti.TrusteeID
		t.Trustee = ti.Trustee
	}
	if ti.OfferingType != nil {
		t.OfferingType = *ti.OfferingType
	}
	if ti.ClosedGroupID != nil {
		t.ClosedGroupID = *ti.ClosedGroupID
		t.ClosedGroup = ti.ClosedGroup
	}

	if ti.SecApproval > 0 {
		t.SecApproval = ti.SecApproval
		t.SecApprovalIdNumber = *ti.SecApprovalIdNumber
	}

	if ti.IssuingWalletPublicKey != nil {
		t.IssuingWalletPublicKey = *ti.IssuingWalletPublicKey
	}
	if ti.IssuingWalletAlias != nil {
		t.IssuingWalletAlias = *ti.IssuingWalletAlias
	}
	if ti.MarketMakingWallet != nil {
		t.MarketMakingWallet = *ti.MarketMakingWallet
	}
	if ti.AssetDescription != nil {
		t.AssetDescription = titleCaser.String(*ti.AssetDescription)
	}
	if ti.AssetCountryLocation != nil {
		t.AssetCountryLocation = *ti.AssetCountryLocation
		t.CountryConfig = CountryCode(t.AssetCountryLocation).GetConfig(gc)
	}

	if ti.AssetPhysicalAddress != nil {
		t.AssetPhysicalAddress = titleCaser.String(*ti.AssetPhysicalAddress)
	}
	if ti.AssetLongitude != nil {
		t.AssetLongitude = *ti.AssetLongitude
	}

	if ti.AssetLatitude != nil {
		t.AssetLatitude = *ti.AssetLatitude
	}
	if ti.OwnershipType != nil {
		t.OwnershipType = *ti.OwnershipType
	}
	if ti.OwnershipKind != nil {
		t.OwnershipKind = *ti.OwnershipKind
	}
	if ti.AssetOwnerName != nil {
		t.AssetOwnerName = strings.ToUpper(*ti.AssetOwnerName)
	}
	if ti.AssetOwnerAddress != nil {
		t.AssetOwnerAddress = titleCaser.String(*ti.AssetOwnerAddress)
	}
	if ti.AssetManagerID > 0 {
		t.AssetManagerID = ti.AssetManagerID
		t.AssetManager = ti.AssetManager
	}

	if ti.TokenizationApplicationFee > 0 {
		t.TokenizationApplicationFee = ti.TokenizationApplicationFee
		t.TokenizationApplicationFeeAsset = ti.TokenizationApplicationFeeAsset
	} 
	// else {

	// 	t.TokenizationApplicationFee = t.CountryConfig.TokenizationApplicationFee
	// 	t.TokenizationApplicationFeeAsset = t.CountryConfig.TokenizationApplicationFeeAsset

	// }

	if ti.AssetQuoteCurrency != nil {
		t.AssetQuoteCurrency = *ti.AssetQuoteCurrency
	}
	t.AssetCurrentValue = ti.AssetCurrentValue

	t.AssetMscCostOutisdeOfValuation = ti.AssetMscCostOutisdeOfValuation
	t.ValueOfTokenizedAsset = ti.ValueOfTokenizedAsset

	if ti.ProtectionMethods != nil {
		t.ProtectionMethods = *ti.ProtectionMethods
	}
	if ti.InsuranceCompanyName != nil {
		t.InsuranceCompanyName = *ti.InsuranceCompanyName
	}
	if ti.InsurancePolicyHolder != nil {
		t.InsurancePolicyHolder = *ti.InsurancePolicyHolder
	}
	if ti.InsurancePolicyNumber != nil {
		t.InsurancePolicyNumber = *ti.InsurancePolicyNumber
	}
	t.PercentageValueOfInsurance = ti.PercentageValueOfInsurance
	t.IsFreeFromLiensAndEncumbrances = ti.IsFreeFromLiensAndEncumbrances
	t.AssetAlreadyExists = ti.AssetAlreadyExists
	t.AssetTokenizationDocuments = ti.AssetTokenizationDocuments
	t.ProofOfPaymentDocuments = ti.ProofOfPaymentDocuments
	t.NumberOfTokenToBeIssued = ti.NumberOfTokenToBeIssued
	t.NumberOfTokenToBeSold = ti.NumberOfTokenToBeSold
	t.TotalTokenHeldByManager = ti.TotalTokenHeldByManager

	if ti.WalletToHoldAssetsNotForSale != nil {
		t.WalletToHoldAssetsNotForSale = *ti.WalletToHoldAssetsNotForSale
	}
	t.PricePerToken = ti.PricePerToken
	t.SalesStart = ti.SalesStart
	t.SalesEnd = ti.SalesEnd
	t.CapOnPurchase = ti.CapOnPurchase
	t.CapQuantity = ti.CapQuantity
	t.CapDurationInDays = ti.CapDurationInDays
	t.CapAmountInFiat = ti.CapAmountInFiat

	if ti.ProceedCycle != nil {
		t.ProceedCycle = *ti.ProceedCycle
	}
	if ti.ProceedPayoutCurrency != nil {
		t.ProceedPayoutCurrency = *ti.ProceedPayoutCurrency
	}
	if ti.TokenizationFeeID != nil {
		t.TokenizationFeeID = *ti.TokenizationFeeID
		t.TokenizationFee = ti.TokenizationFee
	}
	if ti.ExemptedCountries != nil {
		t.ExemptedCountries = *ti.ExemptedCountries
	}

	if ti.AdditionalKYCRequirements != nil {

		t.HasAdditionalKYCRequirements = ti.HasAdditionalKYCRequirements
		t.AdditionalKYCRequirements = *ti.AdditionalKYCRequirements
	}

	t.InvestorAccreditationRequired = ti.InvestorAccreditationRequired
	if ti.LastUpdatedBy != nil {
		t.LastUpdatedBy = *ti.LastUpdatedBy
	}
	t.AssetTokenizationStatus = ti.AssetTokenizationStatus
	t.AgreeTransferTitleToCustodian = ti.AgreeTransferTitleToCustodian
	t.ContractualProtectionRevGuarantees = ti.ContractualProtectionRevGuarantees
	t.ContractualProtectionPerfBond = ti.ContractualProtectionPerfBond
	t.ContractualProtectionSLA = ti.ContractualProtectionSLA
	t.RiskSharingMechanismPPPs = ti.RiskSharingMechanismPPPs
	t.RiskSharingMechanismHedgeInstruments = ti.RiskSharingMechanismHedgeInstruments
	t.RiskSharingMechanismCompletionGuarantees = ti.RiskSharingMechanismCompletionGuarantees
	if ti.IndependentMonitoringList != nil {

		t.IndependentMonitoringList = *ti.IndependentMonitoringList
	}
	t.ESGSafeguardsSusCerts = ti.ESGSafeguardsSusCerts
	t.ESGSafeguardsCommEngPlans = ti.ESGSafeguardsCommEngPlans
	t.SecurityMeasuresAccessControl = ti.SecurityMeasuresAccessControl
	t.SecurityMeasuresSurveilanceSystems = ti.SecurityMeasuresSurveilanceSystems
	t.SecurityMeasuresOnSiteSecurityPersonnel = ti.SecurityMeasuresOnSiteSecurityPersonnel
	t.SecurityMeasuresPerimeterSecurity = ti.SecurityMeasuresPerimeterSecurity
	t.SecurityMeasuresCriticalInfraProtections = ti.SecurityMeasuresCriticalInfraProtections
	if ti.OtherAssetProtection != nil {

		t.OtherAssetProtection = *ti.OtherAssetProtection
	}
	if ti.LegalAdvisor != nil {

		t.LegalAdvisor = strings.ToUpper(*ti.LegalAdvisor)
	}
	if ti.FinancialAdvisor != nil {

		t.FinancialAdvisor = strings.ToUpper(*ti.FinancialAdvisor)
	}

	if ti.BankID != nil {

		t.BankID = *ti.BankID
		t.Bank = ti.UpdateBank(gc)
	}

	if ti.BeneficiaryName != nil {
		t.BeneficiaryName = strings.ToUpper(*ti.BeneficiaryName)
	}
	if ti.AccountNumber != nil {
		t.AccountNumber = *ti.AccountNumber
	}

	t.UndertakingNoLien = ti.UndertakingNoLien
	t.UndertakingNotCollateral = ti.UndertakingNotCollateral
	t.UndertakingNoClaims = ti.UndertakingNoClaims
	t.UndertakingNoForeclosure = ti.UndertakingNoForeclosure
	t.ComplianceNoViolation = ti.ComplianceNoViolation
	t.ComplianceAllPermits = ti.ComplianceAllPermits
	t.OutstandingFinancialRespNoDebts = ti.OutstandingFinancialRespNoDebts
	t.OutstandingFinancialRespNoHiddenLiabilities = ti.OutstandingFinancialRespNoHiddenLiabilities
	t.RiskManagementFullyInsured = ti.RiskManagementFullyInsured
	t.RiskManagementDeclaredValue = ti.RiskManagementDeclaredValue
	t.PhysicalConditionSound = ti.PhysicalConditionSound
	t.PhysicalConditionNolease = ti.PhysicalConditionNolease
	t.PhysicalConditionNoUndisclosedEasements = ti.PhysicalConditionNoUndisclosedEasements

	if ti.ProjectStrategicObjectives != nil {
		t.ProjectStrategicObjectives = *ti.ProjectStrategicObjectives
	}

	// asNN(ti.ProjectStrategicObjectives, &t.ProjectStrategicObjectives)

	// asNN(ti.ProjectDevelopmentTimeline, &t.ProjectDevelopmentTimeline)

	if ti.ProjectDevelopmentTimeline != nil {
		t.ProjectDevelopmentTimeline = *ti.ProjectDevelopmentTimeline
	}

	// asNN(ti.ProjectKeyMilestoneAndDates, &t.ProjectKeyMilestoneAndDates)
	if ti.ProjectKeyMilestoneAndDates != nil {
		t.ProjectKeyMilestoneAndDates = *ti.ProjectKeyMilestoneAndDates
	}

	// asNN(ti.ProjectScope, &t.ProjectScope)

	if ti.ProjectScope != nil {
		t.ProjectScope = *ti.ProjectScope
	}
	// asNN(ti.ProjectEconomicBenefits, &t.ProjectEconomicBenefits)
	if ti.ProjectEconomicBenefits != nil {
		t.ProjectEconomicBenefits = *ti.ProjectEconomicBenefits
	}

	t.ProjectExpectedNoOfJobs = ti.ProjectExpectedNoOfJobs

	// asNN(ti.ProjectIntendedSocialBenefits, &t.ProjectIntendedSocialBenefits)
	if ti.ProjectIntendedSocialBenefits != nil {
		t.ProjectIntendedSocialBenefits = *ti.ProjectIntendedSocialBenefits
	}

	// asNN(ti.ProjectTechnicalPartners, &t.ProjectTechnicalPartners)
	if ti.ProjectTechnicalPartners != nil {
		t.ProjectTechnicalPartners = *ti.ProjectTechnicalPartners
	}
	// asNN(ti.ProjectFinancialPartners, &t.ProjectFinancialPartners)
	if ti.ProjectFinancialPartners != nil {
		t.ProjectFinancialPartners = *ti.ProjectFinancialPartners
	}

	t.EstimatedProjectIRR = ti.EstimatedProjectIRR
	t.EstimatedProjectROI = ti.EstimatedProjectROI
	t.EstimatedProjectNPV = ti.EstimatedProjectNPV
	t.EstimatedProjectPaybackPeriodsInMonths = ti.EstimatedProjectPaybackPeriodsInMonths

	// asNN(ti.KeyAssumptionsList, &t.KeyAssumptionsList)
	if ti.KeyAssumptionsList != nil {
		t.KeyAssumptionsList = *ti.KeyAssumptionsList
	}

	// asNN(ti.ProjectIdentifiedLegalRisks, &t.ProjectIdentifiedLegalRisks)

	if ti.ProjectIdentifiedLegalRisks != nil {
		t.ProjectIdentifiedLegalRisks = *ti.ProjectIdentifiedLegalRisks
	}

	// asNN(ti.ProjectIdentifiedRegulatoryRisks, &t.ProjectIdentifiedRegulatoryRisks)
	if ti.ProjectIdentifiedRegulatoryRisks != nil {
		t.ProjectIdentifiedRegulatoryRisks = *ti.ProjectIdentifiedRegulatoryRisks
	}

	// asNN(ti.ProjectIdentifiedOperationalOrExecutionRisks, &t.ProjectIdentifiedOperationalOrExecutionRisks)
	if ti.ProjectIdentifiedOperationalOrExecutionRisks != nil {
		t.ProjectIdentifiedOperationalOrExecutionRisks = *ti.ProjectIdentifiedOperationalOrExecutionRisks
	}

	// asNN(ti.ProjectIdentifiedMarketRisks, &t.ProjectIdentifiedMarketRisks)
	if ti.ProjectIdentifiedMarketRisks != nil {
		t.ProjectIdentifiedMarketRisks = *ti.ProjectIdentifiedMarketRisks
	}
	// asNN(ti.ProjectIdentifiedOtherRelevantRisks, &t.ProjectIdentifiedOtherRelevantRisks)
	if ti.ProjectIdentifiedOtherRelevantRisks != nil {
		t.ProjectIdentifiedOtherRelevantRisks = *ti.ProjectIdentifiedOtherRelevantRisks
	}

	t.NumberOfExpressedInterests = ti.CountExpressedInterests(gc)
	t.NumberOfSubscribers = ti.CountNumberOfSubscribers(gc)
	t.QuantityOfTokensSoldInFiat = ti.SumAmountSoldInFiat(gc)
	if t.PricePerToken > 0 {
		t.QuantityOfTokensSold = decimal.NewFromFloat(t.QuantityOfTokensSoldInFiat / t.PricePerToken).Truncate(7).InexactFloat64()
	}
	t.PurchaseCommitments = ti.SumExpressedInterest(gc)
	if ti.DeepLink == nil && ti.AssetCode != nil && ti.IssuingWalletPublicKey != nil {
		//set deep link
		p, _ := dynamiclinks.GenerateTokenizedAssetDeeplink(*ti.AssetCode, *ti.IssuingWalletPublicKey, gc)
		if len(p.DynamicLink) > 0 {
			ti.DeepLink = &p.DynamicLink
			//save the tokenized asset information
			gc.DB.Omit(clause.Associations).Save(ti)
		}
	}
	if ti.DeepLink != nil {
		t.DeepLink = *ti.DeepLink
	}

	t.TotalIssueSize = ti.TotalIssueSize
	t.FaceValuePerUnit = ti.FaceValuePerUnit
	t.CouponOrInterestRate = ti.CouponOrInterestRate
	t.ExpectedYield = ti.ExpectedYield
	t.MinimumInvestmentAmount = ti.MinimumInvestmentAmount
	t.WeightedAverageLife = ti.WeightedAverageLife
	t.UnderlyingAssetPoolSize = ti.UnderlyingAssetPoolSize
	t.CouponRate = ti.CouponRate
	t.SpreadOrMargin = ti.SpreadOrMargin
	t.TotalExpenseRatio = ti.TotalExpenseRatio
	t.ExitLoadRedemptionFee = ti.ExitLoadRedemptionFee
	t.Tenure = ti.Tenure
	t.InitialNetAssetValue = ti.InitialNetAssetValue
	t.EntryLoad = ti.EntryLoad
	t.PerformanceFee = ti.PerformanceFee
	t.VolatilityEstimate = ti.VolatilityEstimate
	t.DividendYield = ti.DividendYield
	t.ExitLoadOrRedemptionFee = ti.ExitLoadOrRedemptionFee
	t.AverageMaturity = ti.AverageMaturity
	t.YieldToMaturity = ti.YieldToMaturity
	t.PerformanceFeeIfAny = ti.PerformanceFeeIfAny
	t.AverageMaturityHybrid = ti.AverageMaturityHybrid
	t.YieldToMaturityHybrid = ti.YieldToMaturityHybrid
	t.PricingYield = ti.PricingYield
	t.Quantity = ti.Quantity
	t.MinimumPurchaseAmount = ti.MinimumPurchaseAmount
	t.AssetValuation = ti.AssetValuation
	t.HoldingPeriod = ti.HoldingPeriod
	t.RisksCoverageValue = ti.RisksCoverageValue
	t.RecoveryEstimate = ti.RecoveryEstimate
	t.ExitLoadFee = ti.ExitLoadFee
	t.EntryLoadFee = ti.EntryLoadFee
	t.PortfolioPerformanceFee = ti.PortfolioPerformanceFee

	//////
	if ti.IsinOrSerialNumber != nil {
		t.IsinOrSerialNumber = *ti.IsinOrSerialNumber
	}
	if ti.InstrumentName != nil {
		t.InstrumentName = *ti.InstrumentName
	}
	if ti.InstrumentType != nil {
		t.InstrumentType = *ti.InstrumentType
	}
	if ti.CouponOrInterestRateType != nil {
		t.CouponOrInterestRateType = *ti.CouponOrInterestRateType
	}
	if ti.ReferenceIndex != nil {
		t.ReferenceIndex = *ti.ReferenceIndex
	}
	if ti.EarlyRedemptionOptionInvestor != nil {
		t.EarlyRedemptionOptionInvestor = *ti.EarlyRedemptionOptionInvestor
	}
	if ti.TaxTreatmentTokenHolders != nil {
		t.TaxTreatmentTokenHolders = *ti.TaxTreatmentTokenHolders
	}
	if ti.PaymentStructureToTokenHolders != nil {
		t.PaymentStructureToTokenHolders = *ti.PaymentStructureToTokenHolders
	}
	if ti.RedemptionMethod != nil {
		t.RedemptionMethod = *ti.RedemptionMethod
	}
	if ti.PaymentStructure != nil {
		t.PaymentStructure = *ti.PaymentStructure
	}
	if ti.RepaymentMethod != nil {
		t.RepaymentMethod = *ti.RepaymentMethod
	}
	if ti.PaymentCycle != nil {
		t.PaymentCycle = *ti.PaymentCycle
	}
	if ti.PoolComposition != nil {
		t.PoolComposition = *ti.PoolComposition
	}
	if ti.CreditEnhancementMethod != nil {
		t.CreditEnhancementMethod = *ti.CreditEnhancementMethod
	}
	if ti.SummaryOfUseOfProceeds != nil {
		t.SummaryOfUseOfProceeds = *ti.SummaryOfUseOfProceeds
	}
	if ti.CreditRatingIfAny != nil {
		t.CreditRatingIfAny = *ti.CreditRatingIfAny
	}
	if ti.IssuerName != nil {
		t.IssuerName = *ti.IssuerName
	}
	if ti.IssuerType != nil {
		t.IssuerType = *ti.IssuerType
	}
	if ti.IssuerContactPerson != nil {
		t.IssuerContactPerson = *ti.IssuerContactPerson
	}
	if ti.ContactEmail != nil {
		t.ContactEmail = *ti.ContactEmail
	}
	if ti.ContactPhoneNumber != nil {
		t.ContactPhoneNumber = *ti.ContactPhoneNumber
	}
	if ti.BriefCompanyOverview != nil {
		t.BriefCompanyOverview = *ti.BriefCompanyOverview
	}
	if ti.MortgageOriginators != nil {
		t.MortgageOriginators = *ti.MortgageOriginators
	}
	if ti.Servicer != nil {
		t.Servicer = *ti.Servicer
	}
	if ti.InstrumentTrustee != nil {
		t.InstrumentTrustee = *ti.InstrumentTrustee
	}
	if ti.InstrumentCustodians != nil {
		t.InstrumentCustodians = *ti.InstrumentCustodians
	}
	if ti.InstrumentLegalAdvisor != nil {
		t.InstrumentLegalAdvisor = *ti.InstrumentLegalAdvisor
	}
	if ti.SpecialPurposeVehicle != nil {
		t.SpecialPurposeVehicle = *ti.SpecialPurposeVehicle
	}
	if ti.InstrumentAssetManagerOrAdministrator != nil {
		t.InstrumentAssetManagerOrAdministrator = *ti.InstrumentAssetManagerOrAdministrator
	}
	if ti.UnderwriterIfAny != nil {
		t.UnderwriterIfAny = *ti.UnderwriterIfAny
	}
	if ti.InstrumentCreditRatingAgency != nil {
		t.InstrumentCreditRatingAgency = *ti.InstrumentCreditRatingAgency
	}
	if ti.AuditorOrVerifier != nil {
		t.AuditorOrVerifier = *ti.AuditorOrVerifier
	}
	if ti.CreditRiskAssessment != nil {
		t.CreditRiskAssessment = *ti.CreditRiskAssessment
	}
	if ti.CreditRating != nil {
		t.CreditRating = *ti.CreditRating
	}
	if ti.PrepaymentRisk != nil {
		t.PrepaymentRisk = *ti.PrepaymentRisk
	}
	if ti.InterestRateRisk != nil {
		t.InterestRateRisk = *ti.InterestRateRisk
	}
	if ti.StructuralComplexityRisk != nil {
		t.StructuralComplexityRisk = *ti.StructuralComplexityRisk
	}
	if ti.LegalOrRegulatoryRisk != nil {
		t.LegalOrRegulatoryRisk = *ti.LegalOrRegulatoryRisk
	}
	if ti.OperationalRisk != nil {
		t.OperationalRisk = *ti.OperationalRisk
	}
	if ti.MarketRisk != nil {
		t.MarketRisk = *ti.MarketRisk
	}
	if ti.EsgRisk != nil {
		t.EsgRisk = *ti.EsgRisk
	}
	if ti.MitigationMeasures != nil {
		t.MitigationMeasures = *ti.MitigationMeasures
	}
	if ti.IssuingAuthority != nil {
		t.IssuingAuthority = *ti.IssuingAuthority
	}
	if ti.RegulatoryApprovalId != nil {
		t.RegulatoryApprovalId = *ti.RegulatoryApprovalId
	}
	if ti.LicenseApprovalReferenceNumber != nil {
		t.LicenseApprovalReferenceNumber = *ti.LicenseApprovalReferenceNumber
	}
	if ti.ListingStatus != nil {
		t.ListingStatus = *ti.ListingStatus
	}
	if ti.CurrencyOfIssuance != nil {
		t.CurrencyOfIssuance = *ti.CurrencyOfIssuance
	}
	if ti.CouponRateType != nil {
		t.CouponRateType = *ti.CouponRateType
	}
	if ti.ResetFrequency != nil {
		t.ResetFrequency = *ti.ResetFrequency
	}
	if ti.CouponPaymentFrequency != nil {
		t.CouponPaymentFrequency = *ti.CouponPaymentFrequency
	}
	if ti.RedemptionStructure != nil {
		t.RedemptionStructure = *ti.RedemptionStructure
	}
	if ti.EarlyRedemptionOption != nil {
		t.EarlyRedemptionOption = *ti.EarlyRedemptionOption
	}
	if ti.EarlyRedemptionPenalty != nil {
		t.EarlyRedemptionPenalty = *ti.EarlyRedemptionPenalty
	}
	if ti.TaxTreatment != nil {
		t.TaxTreatment = *ti.TaxTreatment
	}
	if ti.NavOrMarketValueUpdates != nil {
		t.NavOrMarketValueUpdates = *ti.NavOrMarketValueUpdates
	}
	if ti.ImpactMetrics != nil {
		t.ImpactMetrics = *ti.ImpactMetrics
	}
	if ti.LegalBacking != nil {
		t.LegalBacking = *ti.LegalBacking
	}
	if ti.DefaultHistory != nil {
		t.DefaultHistory = *ti.DefaultHistory
	}
	if ti.RiskFactorsSummary != nil {
		t.RiskFactorsSummary = *ti.RiskFactorsSummary
	}
	if ti.PayingAgent != nil {
		t.PayingAgent = *ti.PayingAgent
	}
	if ti.Auditor != nil {
		t.Auditor = *ti.Auditor
	}
	if ti.RegistrarOrCSCSAgent != nil {
		t.RegistrarOrCSCSAgent = *ti.RegistrarOrCSCSAgent
	}
	if ti.FundStructure != nil {
		t.FundStructure = *ti.FundStructure
	}
	if ti.AssetManagementCompanyName != nil {
		t.AssetManagementCompanyName = *ti.AssetManagementCompanyName
	}
	if ti.FundManagers != nil {
		t.FundManagers = *ti.FundManagers
	}
	if ti.RegulatoryLicenseNumber != nil {
		t.RegulatoryLicenseNumber = *ti.RegulatoryLicenseNumber
	}
	if ti.FundLaunchDate != nil {
		t.FundLaunchDate = *ti.FundLaunchDate
	}
	if ti.NavUpdateFrequency != nil {
		t.NavUpdateFrequency = *ti.NavUpdateFrequency
	}
	if ti.NavCalculationMethod != nil {
		t.NavCalculationMethod = *ti.NavCalculationMethod
	}
	if ti.RedemptionRules != nil {
		t.RedemptionRules = *ti.RedemptionRules
	}
	t.LockInPeriod = ti.LockInPeriod

	if ti.DividendPolicy != nil {
		t.DividendPolicy = *ti.DividendPolicy
	}
	if ti.LiquidityProfile != nil {
		t.LiquidityProfile = *ti.LiquidityProfile
	}
	if ti.DistributionFrequency != nil {
		t.DistributionFrequency = *ti.DistributionFrequency
	}
	if ti.DistributionMethod != nil {
		t.DistributionMethod = *ti.DistributionMethod
	}
	if ti.BenchmarkComparisonMethod != nil {
		t.BenchmarkComparisonMethod = *ti.BenchmarkComparisonMethod
	}
	if ti.FeeBreakdownSummary != nil {
		t.FeeBreakdownSummary = *ti.FeeBreakdownSummary
	}
	if ti.InvestmentObjective != nil {
		t.InvestmentObjective = *ti.InvestmentObjective
	}
	if ti.EquityStrategy != nil {
		t.EquityStrategy = *ti.EquityStrategy
	}
	if ti.MarketCapitalizationFocus != nil {
		t.MarketCapitalizationFocus = *ti.MarketCapitalizationFocus
	}
	if ti.BenchmarkIndex != nil {
		t.BenchmarkIndex = *ti.BenchmarkIndex
	}
	if ti.SectorExposureLimits != nil {
		t.SectorExposureLimits = *ti.SectorExposureLimits
	}
	if ti.TopHoldings != nil {
		t.TopHoldings = *ti.TopHoldings
	}
	if ti.GeographicExposure != nil {
		t.GeographicExposure = *ti.GeographicExposure
	}
	if ti.RiskProfile != nil {
		t.RiskProfile = *ti.RiskProfile
	}
	if ti.TrusteeName != nil {
		t.TrusteeName = *ti.TrusteeName
	}
	if ti.FundAdministrator != nil {
		t.FundAdministrator = *ti.FundAdministrator
	}
	if ti.InvestmentCommitteeMembers != nil {
		t.InvestmentCommitteeMembers = *ti.InvestmentCommitteeMembers
	}
	if ti.IsinOrSecFundCode != nil {
		t.IsinOrSecFundCode = *ti.IsinOrSecFundCode
	}
	if ti.FundRiskRating != nil {
		t.FundRiskRating = *ti.FundRiskRating
	}
	if ti.AssetAllocation != nil {
		t.AssetAllocation = *ti.AssetAllocation
	}
	if ti.CreditRatingProfile != nil {
		t.CreditRatingProfile = *ti.CreditRatingProfile
	}

	t.LockInPeriodPortfolio = ti.LockInPeriodPortfolio

	if ti.TopEquityHoldings != nil {
		t.TopEquityHoldings = *ti.TopEquityHoldings
	}
	if ti.TargetAllocation != nil {
		t.TargetAllocation = *ti.TargetAllocation
	}
	if ti.AllowedAllocationRange != nil {
		t.AllowedAllocationRange = *ti.AllowedAllocationRange
	}
	if ti.AssetClassesIncluded != nil {
		t.AssetClassesIncluded = *ti.AssetClassesIncluded
	}
	if ti.RebalancingFrequency != nil {
		t.RebalancingFrequency = *ti.RebalancingFrequency
	}
	if ti.BenchmarkIndexComposite != nil {
		t.BenchmarkIndexComposite = *ti.BenchmarkIndexComposite
	}
	if ti.TopEquityHoldingsList != nil {
		t.TopEquityHoldingsList = *ti.TopEquityHoldingsList
	}
	if ti.TopDebtHoldings != nil {
		t.TopDebtHoldings = *ti.TopDebtHoldings
	}
	if ti.CreditRatingDistribution != nil {
		t.CreditRatingDistribution = *ti.CreditRatingDistribution
	}
	if ti.TitleOfIssuance != nil {
		t.TitleOfIssuance = *ti.TitleOfIssuance
	}
	if ti.TypeOfCommercialPaper != nil {
		t.TypeOfCommercialPaper = *ti.TypeOfCommercialPaper
	}
	if ti.UseOfProceeds != nil {
		t.UseOfProceeds = *ti.UseOfProceeds
	}
	if ti.Ranking != nil {
		t.Ranking = *ti.Ranking
	}
	if ti.BackingSecurity != nil {
		t.BackingSecurity = *ti.BackingSecurity
	}
	if ti.IssuerRegistrationNumber != nil {
		t.IssuerRegistrationNumber = *ti.IssuerRegistrationNumber
	}
	if ti.IncorporationDate != nil {
		t.IncorporationDate = *ti.IncorporationDate
	}
	if ti.RcNumber != nil {
		t.RcNumber = *ti.RcNumber
	}
	if ti.TaxIdNumber != nil {
		t.TaxIdNumber = *ti.TaxIdNumber
	}
	if ti.OfficeAddress != nil {
		t.OfficeAddress = *ti.OfficeAddress
	}
	if ti.Rating != nil {
		t.Rating = *ti.Rating
	}

	/////

	if ti.PartiesInvolvedIssuer != nil {
		t.PartiesInvolvedIssuer = *ti.PartiesInvolvedIssuer
	}
	if ti.PartiesInvolvedArranger != nil {
		t.PartiesInvolvedArranger = *ti.PartiesInvolvedArranger
	}
	if ti.PartiesInvolvedLegalAdviser != nil {
		t.PartiesInvolvedLegalAdviser = *ti.PartiesInvolvedLegalAdviser
	}
	if ti.PartiesInvolvedAuditor != nil {
		t.PartiesInvolvedAuditor = *ti.PartiesInvolvedAuditor
	}
	if ti.PartiesInvolvedRatingAgency != nil {
		t.PartiesInvolvedRatingAgency = *ti.PartiesInvolvedRatingAgency
	}
	if ti.PartiesInvolvedCustodian != nil {
		t.PartiesInvolvedCustodian = *ti.PartiesInvolvedCustodian
	}
	if ti.PartiesInvolvedTrustee != nil {
		t.PartiesInvolvedTrustee = *ti.PartiesInvolvedTrustee
	}
	if ti.PartiesInvolvedAuditorVerifier != nil {
		t.PartiesInvolvedAuditorVerifier = *ti.PartiesInvolvedAuditorVerifier
	}
	if ti.SecurityRiskLegalBacking != nil {
		t.SecurityRiskLegalBacking = *ti.SecurityRiskLegalBacking
	}
	if ti.SecurityRiskCollateral != nil {
		t.SecurityRiskCollateral = *ti.SecurityRiskCollateral
	}
	if ti.SecurityRiskDefaultHistory != nil {
		t.SecurityRiskDefaultHistory = *ti.SecurityRiskDefaultHistory
	}
	if ti.SecurityRiskCreditRating != nil {
		t.SecurityRiskCreditRating = *ti.SecurityRiskCreditRating
	}
	if ti.SecurityRiskRiskFactorsSummary != nil {
		t.SecurityRiskRiskFactorsSummary = *ti.SecurityRiskRiskFactorsSummary
	}
	if ti.SecurityRiskBusinessRisk != nil {
		t.SecurityRiskBusinessRisk = *ti.SecurityRiskBusinessRisk
	}
	if ti.SecurityRiskDefaultRisk != nil {
		t.SecurityRiskDefaultRisk = *ti.SecurityRiskDefaultRisk
	}
	if ti.SecurityRiskLiquidityRisk != nil {
		t.SecurityRiskLiquidityRisk = *ti.SecurityRiskLiquidityRisk
	}
	if ti.SecurityRiskRegulatoryRisk != nil {
		t.SecurityRiskRegulatoryRisk = *ti.SecurityRiskRegulatoryRisk
	}
	if ti.SecurityRiskMarketRisk != nil {
		t.SecurityRiskMarketRisk = *ti.SecurityRiskMarketRisk
	}
	if ti.SecurityRiskOperationalRisk != nil {
		t.SecurityRiskOperationalRisk = *ti.SecurityRiskOperationalRisk
	}
	if ti.SecurityRiskMitigationMeasures != nil {
		t.SecurityRiskMitigationMeasures = *ti.SecurityRiskMitigationMeasures
	}
	if ti.CommodityType != nil {
		t.CommodityType = *ti.CommodityType
	}
	if ti.CommodityDescription != nil {
		t.CommodityDescription = *ti.CommodityDescription
	}
	if ti.QualityGrade != nil {
		t.QualityGrade = *ti.QualityGrade
	}
	if ti.IssuerContactInfo != nil {
		t.IssuerContactInfo = *ti.IssuerContactInfo
	}
	if ti.WarehouseName != nil {
		t.WarehouseName = *ti.WarehouseName
	}
	if ti.WarehouseOperatorName != nil {
		t.WarehouseOperatorName = *ti.WarehouseOperatorName
	}
	if ti.WarehouseLicenseNumber != nil {
		t.WarehouseLicenseNumber = *ti.WarehouseLicenseNumber
	}
	if ti.WarehouseLocation != nil {
		t.WarehouseLocation = *ti.WarehouseLocation
	}
	if ti.WrNumber != nil {
		t.WrNumber = *ti.WrNumber
	}
	if ti.WrIssueDate != nil {
		t.WrIssueDate = *ti.WrIssueDate
	}
	if ti.WrExpiryDate != nil {
		t.WrExpiryDate = *ti.WrExpiryDate
	}
	if ti.WrSystemRegistration != nil {
		t.WrSystemRegistration = *ti.WrSystemRegistration
	}
	if ti.WrRegistrationNumber != nil {
		t.WrRegistrationNumber = *ti.WrRegistrationNumber
	}
	if ti.WrVerifier != nil {
		t.WrVerifier = *ti.WrVerifier
	}
	if ti.StorageCondition != nil {
		t.StorageCondition = *ti.StorageCondition
	}
	if ti.WarehouseAccreditationBody != nil {
		t.WarehouseAccreditationBody = *ti.WarehouseAccreditationBody
	}
	if ti.AutoRollover != nil {
		t.AutoRollover = *ti.AutoRollover
	}
	if ti.CurrentBeneficialOwner != nil {
		t.CurrentBeneficialOwner = *ti.CurrentBeneficialOwner
	}
	if ti.WrCustodianName != nil {
		t.WrCustodianName = *ti.WrCustodianName
	}
	if ti.OwnershipRightsRepresented != nil {
		t.OwnershipRightsRepresented = *ti.OwnershipRightsRepresented
	}
	if ti.TrusteeOrThirdPartyOversight != nil {
		t.TrusteeOrThirdPartyOversight = *ti.TrusteeOrThirdPartyOversight
	}
	if ti.LienOrEncumbrances != nil {
		t.LienOrEncumbrances = *ti.LienOrEncumbrances
	}
	if ti.ValuationDate != nil {
		t.ValuationDate = *ti.ValuationDate
	}
	if ti.ValuationMethodology != nil {
		t.ValuationMethodology = *ti.ValuationMethodology
	}
	if ti.TokenizationObjective != nil {
		t.TokenizationObjective = *ti.TokenizationObjective
	}
	if ti.RedemptionMechanism != nil {
		t.RedemptionMechanism = *ti.RedemptionMechanism
	}
	if ti.PartiesInvolvedUnderwriter != nil {
		t.PartiesInvolvedUnderwriter = *ti.PartiesInvolvedUnderwriter
	}
	if ti.PartiesInvolvedAssetManager != nil {
		t.PartiesInvolvedAssetManager = *ti.PartiesInvolvedAssetManager
	}
	if ti.PartiesInvolvedLegalAdvisor != nil {
		t.PartiesInvolvedLegalAdvisor = *ti.PartiesInvolvedLegalAdvisor
	}
	if ti.PartiesInvolvedRegulator != nil {
		t.PartiesInvolvedRegulator = *ti.PartiesInvolvedRegulator
	}
	if ti.RisksMarketRisk != nil {
		t.RisksMarketRisk = *ti.RisksMarketRisk
	}
	if ti.RisksStorageRisk != nil {
		t.RisksStorageRisk = *ti.RisksStorageRisk
	}
	if ti.RisksTitleRisk != nil {
		t.RisksTitleRisk = *ti.RisksTitleRisk
	}
	if ti.RisksFraudRisk != nil {
		t.RisksFraudRisk = *ti.RisksFraudRisk
	}
	if ti.RisksInsuranceRisk != nil {
		t.RisksInsuranceRisk = *ti.RisksInsuranceRisk
	}
	if ti.RisksOperationalRisk != nil {
		t.RisksOperationalRisk = *ti.RisksOperationalRisk
	}
	if ti.RisksRegulatoryRisk != nil {
		t.RisksRegulatoryRisk = *ti.RisksRegulatoryRisk
	}
	if ti.RisksLiquidityRisk != nil {
		t.RisksLiquidityRisk = *ti.RisksLiquidityRisk
	}
	if ti.RisksForceMajeureRisk != nil {
		t.RisksForceMajeureRisk = *ti.RisksForceMajeureRisk
	}
	if ti.RisksEarlyRedemptionRisk != nil {
		t.RisksEarlyRedemptionRisk = *ti.RisksEarlyRedemptionRisk
	}
	if ti.RisksMitigationMeasures != nil {
		t.RisksMitigationMeasures = *ti.RisksMitigationMeasures
	}
	if ti.RisksInsuranceCoverageSummary != nil {
		t.RisksInsuranceCoverageSummary = *ti.RisksInsuranceCoverageSummary
	}
	if ti.RisksInsuranceProvider != nil {
		t.RisksInsuranceProvider = *ti.RisksInsuranceProvider
	}
	if ti.QuanlityStandard != nil {
		t.QuanlityStandard = *ti.QuanlityStandard
	}
	if ti.IssuerContactInformation != nil {
		t.IssuerContactInformation = *ti.IssuerContactInformation
	}
	if ti.VaultCustodianName != nil {
		t.VaultCustodianName = *ti.VaultCustodianName
	}
	if ti.VaultOperator != nil {
		t.VaultOperator = *ti.VaultOperator
	}
	if ti.VaultLicenseNumber != nil {
		t.VaultLicenseNumber = *ti.VaultLicenseNumber
	}
	if ti.VaultLocation != nil {
		t.VaultLocation = *ti.VaultLocation
	}
	if ti.Number != nil {
		t.Number = *ti.Number
	}
	if ti.IssuerDate != nil {
		t.IssuerDate = *ti.IssuerDate
	}
	if ti.ExpiryDate != nil {
		t.ExpiryDate = *ti.ExpiryDate
	}
	if ti.RegistryRecord != nil {
		t.RegistryRecord = *ti.RegistryRecord
	}
	if ti.Verifier != nil {
		t.Verifier = *ti.Verifier
	}
	if ti.StorageConditions != nil {
		t.StorageConditions = *ti.StorageConditions
	}
	if ti.VaultAccreditationBody != nil {
		t.VaultAccreditationBody = *ti.VaultAccreditationBody
	}
	if ti.OwnershipLegalHolder != nil {
		t.OwnershipLegalHolder = *ti.OwnershipLegalHolder
	}
	if ti.OwnershipCustodianName != nil {
		t.OwnershipCustodianName = *ti.OwnershipCustodianName
	}
	if ti.OwnershipTrustee != nil {
		t.OwnershipTrustee = *ti.OwnershipTrustee
	}
	if ti.OwnershipLienOrEncumbrances != nil {
		t.OwnershipLienOrEncumbrances = *ti.OwnershipLienOrEncumbrances
	}
	if ti.ValuationAssetValuation != nil {
		t.ValuationAssetValuation = *ti.ValuationAssetValuation
	}

	t.HoldingLockinPeriod = ti.HoldingLockinPeriod

	if ti.InsuranceMarketRisk != nil {
		t.InsuranceMarketRisk = *ti.InsuranceMarketRisk
	}
	if ti.InsuranceStorageRisk != nil {
		t.InsuranceStorageRisk = *ti.InsuranceStorageRisk
	}
	if ti.InsuranceTitleRisk != nil {
		t.InsuranceTitleRisk = *ti.InsuranceTitleRisk
	}
	if ti.InsuranceFraudRisk != nil {
		t.InsuranceFraudRisk = *ti.InsuranceFraudRisk
	}
	if ti.InsuranceInsuranceRisk != nil {
		t.InsuranceInsuranceRisk = *ti.InsuranceInsuranceRisk
	}
	if ti.InsuranceOperationalRisk != nil {
		t.InsuranceOperationalRisk = *ti.InsuranceOperationalRisk
	}
	if ti.InsuranceRegulatoryRisk != nil {
		t.InsuranceRegulatoryRisk = *ti.InsuranceRegulatoryRisk
	}
	if ti.InsuranceLiquidityRisk != nil {
		t.InsuranceLiquidityRisk = *ti.InsuranceLiquidityRisk
	}
	if ti.InsuranceForceMajeureRisk != nil {
		t.InsuranceForceMajeureRisk = *ti.InsuranceForceMajeureRisk
	}
	if ti.InsuranceEarlyRedemptionRisk != nil {
		t.InsuranceEarlyRedemptionRisk = *ti.InsuranceEarlyRedemptionRisk
	}
	if ti.InsuranceMitigationMeasures != nil {
		t.InsuranceMitigationMeasures = *ti.InsuranceMitigationMeasures
	}
	if ti.InsuranceInsuranceCoverageSummary != nil {
		t.InsuranceInsuranceCoverageSummary = *ti.InsuranceInsuranceCoverageSummary
	}
	if ti.InsuranceInsuranceProvider != nil {
		t.InsuranceInsuranceProvider = *ti.InsuranceInsuranceProvider
	}
	if ti.InsuranceCoverageValue != nil {
		t.InsuranceCoverageValue = *ti.InsuranceCoverageValue
	}
	if ti.IssuerRegistrationNo != nil {
		t.IssuerRegistrationNo = *ti.IssuerRegistrationNo
	}
	if ti.SectorAndIndustry != nil {
		t.SectorAndIndustry = *ti.SectorAndIndustry
	}
	if ti.LicenseOrPermitNumber != nil {
		t.LicenseOrPermitNumber = *ti.LicenseOrPermitNumber
	}
	if ti.IssuerAdditionalInfo != nil {
		t.IssuerAdditionalInfo = *ti.IssuerAdditionalInfo
	}
	if ti.IsinSerialNumber != nil {
		t.IsinSerialNumber = *ti.IsinSerialNumber
	}
	if ti.EsgOrImpactMetrics != nil {
		t.EsgOrImpactMetrics = *ti.EsgOrImpactMetrics
	}
	if ti.InstrumentAdditionalInfo != nil {
		t.InstrumentAdditionalInfo = *ti.InstrumentAdditionalInfo
	}
	if ti.SecurityType != nil {
		t.SecurityType = *ti.SecurityType
	}
	if ti.CollateralDescription != nil {
		t.CollateralDescription = *ti.CollateralDescription
	}
	if ti.CovenantSummary != nil {
		t.CovenantSummary = *ti.CovenantSummary
	}
	if ti.CovenantTestingFrequency != nil {
		t.CovenantTestingFrequency = *ti.CovenantTestingFrequency
	}
	if ti.EventOfDefaultClauses != nil {
		t.EventOfDefaultClauses = *ti.EventOfDefaultClauses
	}
	if ti.LegalEnforcementMechanism != nil {
		t.LegalEnforcementMechanism = *ti.LegalEnforcementMechanism
	}
	if ti.Guarantee != nil {
		t.Guarantee = *ti.Guarantee
	}
	if ti.RiskProfileAdditionalInfo != nil {
		t.RiskProfileAdditionalInfo = *ti.RiskProfileAdditionalInfo
	}
	if ti.LicenseNumber != nil {
		t.LicenseNumber = *ti.LicenseNumber
	}
	t.PortfolioLockInPeriod = ti.PortfolioLockInPeriod

	//////

	t.IssueDate = ti.IssueDate

	t.MaturityDate = ti.MaturityDate
	t.FundingStructure = ti.FundingStructure
	t.EquityPercentage = ti.EquityPercentage
	t.DebtPercentage = ti.DebtPercentage

	if ti.DebtInstrumentType != nil {
		t.DebtInstrumentType = *ti.DebtInstrumentType
	}
	if ti.PrincipalPaymentMethod != nil {
		t.PrincipalPaymentMethod = *ti.PrincipalPaymentMethod
	}
	if ti.DebtInstrumentRepaymentSource != nil {
		t.DebtInstrumentRepaymentSource = *ti.DebtInstrumentRepaymentSource
	}
	if ti.DebtInstrumentGuaranteesOrEnhancements != nil {
		t.DebtInstrumentGuaranteesOrEnhancements = *ti.DebtInstrumentGuaranteesOrEnhancements
	}
	if ti.DebtInstrumentDefaultAndRecoveryTerms != nil {
		t.DebtInstrumentDefaultAndRecoveryTerms = *ti.DebtInstrumentDefaultAndRecoveryTerms
	}
	if ti.DebtInstrumentRepaymentFrequency != nil {
		t.DebtInstrumentRepaymentFrequency = *ti.DebtInstrumentRepaymentFrequency
	}
	if ti.InterestRepaymentFrequency != nil {
		t.InterestRepaymentFrequency = *ti.InterestRepaymentFrequency
	}
	if ti.DcsrDetails != nil {
		t.DcsrDetails = *ti.DcsrDetails
	}
	if ti.SinkingFundStructure != nil {
		t.SinkingFundStructure = *ti.SinkingFundStructure
	}
	if ti.CovenantMonitoringAgent != nil {
		t.CovenantMonitoringAgent = *ti.CovenantMonitoringAgent
	}
	if ti.RightOfRecourse != nil {
		t.RightOfRecourse = *ti.RightOfRecourse
	}
	t.DebtInstrumentInterestRate = ti.DebtInstrumentInterestRate
	t.DcsrRatio = ti.DcsrRatio
	t.LtvRatio = ti.LtvRatio
	t.InterestCoverageRatio = ti.InterestCoverageRatio
	t.MaximumLeverageRatio = ti.MaximumLeverageRatio
	t.GracePeriod = ti.GracePeriod
	t.TrusteeAppointed = ti.TrusteeAppointed
	t.ReserveFundInPlace = ti.ReserveFundInPlace

	if ti.SecurityOrCollateralOffered != nil {
		t.SecurityOrCollateralOffered = *ti.SecurityOrCollateralOffered
	}
	if ti.FundInstrumentType != nil {
		t.FundInstrumentType = *ti.FundInstrumentType
	}
	if ti.InstrumentRatingAgency != nil {
		t.InstrumentRatingAgency = *ti.InstrumentRatingAgency
	}
	if ti.PortfolioTopHoldings != nil {
		t.PortfolioTopHoldings = *ti.PortfolioTopHoldings
	}

	if ti.CreditRatingAgency != nil {
		t.CreditRatingAgency = *ti.CreditRatingAgency
	}

	if ti.TrusteeRegNumber != nil {
		t.TrusteeRegNumber = *ti.TrusteeRegNumber
	}

	if ti.CreditEnhancerOrGuarantor != nil {
		t.CreditEnhancerOrGuarantor = *ti.CreditEnhancerOrGuarantor
	}

	if ti.BondStructuringAdvisor != nil {
		t.BondStructuringAdvisor = *ti.BondStructuringAdvisor
	}

	if ti.EntitiesAdditionalInfo != nil {
		t.EntitiesAdditionalInfo = *ti.EntitiesAdditionalInfo
	}

	if ti.AuthorizedRepresentativeName != nil {
		t.AuthorizedRepresentativeName = *ti.AuthorizedRepresentativeName
	}

	if ti.AuthorizedRepresentativeTitleOrPosition != nil {
		t.AuthorizedRepresentativeTitleOrPosition = *ti.AuthorizedRepresentativeTitleOrPosition
	}

	if ti.AuthorizedRepresentativeEmail != nil {
		t.AuthorizedRepresentativeEmail = *ti.AuthorizedRepresentativeEmail
	}

	t.AcceptTokenizationTermsAndAgreement = ti.AcceptTokenizationTermsAndAgreement
	t.AttestInformationAccurateAndVerifiable = ti.AttestInformationAccurateAndVerifiable
	t.AcknowledgedSuitabilityCriteria = ti.AcknowledgedSuitabilityCriteria

	if ti.MinimumKycTier != nil {
		t.MinimumKycTier = *ti.MinimumKycTier
	}

	if ti.InvestorCategory != nil {
		t.InvestorCategory = *ti.InvestorCategory
	}
	t.WithholdingTaxDisclosure = ti.WithholdingTaxDisclosure

	return t

}

type PaginatedTokenizedAssets struct {
	Pages        int                  `json:"pages"`
	CurrentPage  int                  `json:"currentPage"`
	TotalRecords int                  `json:"totalRecords"`
	Limit        int                  `json:"limit"`
	Records      []TokenizedAssetJSON `json:"records"`
}

type PaginatedTokenizedAssetSubscription struct {
	Pages        int                          `json:"pages"`
	CurrentPage  int                          `json:"currentPage"`
	TotalRecords int                          `json:"totalRecords"`
	Limit        int                          `json:"limit"`
	Records      []TokenizedAssetSubscription `json:"records"`
}
type PaginatedExpressionOfInterest struct {
	Pages        int                    `json:"pages"`
	CurrentPage  int                    `json:"currentPage"`
	TotalRecords int                    `json:"totalRecords"`
	Limit        int                    `json:"limit"`
	Records      []ExpressionOfInterest `json:"records"`
}

func (tas *TokenizedAssetSubscription) UpdateTokenizedAssetSubscriptionFromInput(subscriberUsername string, subscriberWallet *UserWallet, input *TokenizedAssetSubscriptionInput, ta *TokenizedAsset, gc *sharedconfig.GlobalConfig) (ts TokenizedAssetSubscription) {
	if tas == nil {
		return
	}

	if ta.AssetTokenizationStatus < 4 {
		return
	}
	tas.ID = uuid.NewString()
	tas.TokenizedAssetID = ta.ID
	tas.AssetCode = *ta.AssetCode
	tas.AssetIssuer = *ta.IssuingWalletPublicKey
	tas.WalletAlias = subscriberWallet.Alias
	tas.WalletPublicKey = subscriberWallet.ID
	tas.Amount = decimal.NewFromFloat(input.Amount).Truncate(7).InexactFloat64()
	tas.Price = ta.PricePerToken
	tas.SubscriberUsername = subscriberUsername
	return *tas
}

func (e *ExpressionOfInterest) UpdateExpressionOfInterestFromInput(subscriberUsername string, input *ExpressionOfInterestInput, ta *TokenizedAsset, gc *sharedconfig.GlobalConfig) (ei ExpressionOfInterest) {
	if e == nil {
		return
	}

	if ta.AssetTokenizationStatus < 4 {
		return
	}
	e.TokenizedAssetID = ta.ID
	e.AssetCode = *ta.AssetCode
	e.AssetIssuer = *ta.IssuingWalletPublicKey
	e.Amount = decimal.NewFromFloat(input.Amount).Truncate(7).InexactFloat64()
	e.Price = ta.PricePerToken
	e.SubscriberUsername = subscriberUsername
	return *e
}

// TokenizationMintingApprovers is a custom type wrapping a slice of TokenizationMintingApprover
type TokenizationMintingApprovers []TokenizationMintingApprover

// TokenizationMintingInitiators is a custom type wrapping a slice of TokenizationMintingInitiator
type TokenizationMintingInitiators []TokenizationMintingInitiator

// ToCSV converts the TokenizationMintingApprovers receiver to a CSV string of Approver values
func (tma TokenizationMintingApprovers) ToCSV() string {
	if len(tma) == 0 {
		return "" // Return empty string for empty slice
	}

	// Prepare a slice of approvers values
	a := make([]string, 0, len(tma)) // Pre-allocate capacity for efficiency
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	for _, v := range tma {
		a = append(a, v.Approver)
	}

	// Write the approvers values as a single CSV row
	err := writer.Write(a)
	if err != nil {
		return "" // Return empty string on error
	}

	// Flush the writer to ensure all data is written to the buffer
	writer.Flush()
	if err := writer.Error(); err != nil {
		return "" // Return empty string on flush error
	}

	// Return the CSV string
	return buf.String()
}

// ToCSV converts the TokenizationMintinginitiators receiver to a CSV string of initator values
func (tma TokenizationMintingInitiators) ToCSV() string {
	if len(tma) == 0 {
		return "" // Return empty string for empty slice
	}

	// Prepare a slice of initiators values
	a := make([]string, 0, len(tma)) // Pre-allocate capacity for efficiency
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	for _, v := range tma {
		a = append(a, v.Initiator)
	}

	// Write the initiator values as a single CSV row
	err := writer.Write(a)
	if err != nil {
		return "" // Return empty string on error
	}

	// Flush the writer to ensure all data is written to the buffer
	writer.Flush()
	if err := writer.Error(); err != nil {
		return "" // Return empty string on flush error
	}

	// Return the CSV string
	return buf.String()
}

type AssetManagerID uint64
type ApprovedCustodianID uint64
type IssuingHouseID uint64
type LegalAndProfesionalPartnerID uint64
type RatingAgencyID uint64
type TrusteeID uint64

func (a AssetManagerID) GetAssetManager(gc *sharedconfig.GlobalConfig) (assetManager AssetManager) {
	gc.DB.Where("id = ?", uint64(a)).First(&assetManager)

	return
}

func (a ApprovedCustodianID) GetApprovedCustodian(gc *sharedconfig.GlobalConfig) (custodian ApprovedAssetCustodian) {
	gc.DB.Where("id = ?", uint64(a)).First(&custodian)

	return
}

func (a IssuingHouseID) GetAssetIssuingHouse(gc *sharedconfig.GlobalConfig) (assetIssuingHouse AssetIssuingHouse) {
	gc.DB.Where("id = ?", uint64(a)).First(&assetIssuingHouse)

	return
}

func (a LegalAndProfesionalPartnerID) GetLegalAndProfesionalPartner(gc *sharedconfig.GlobalConfig) (lp LegalAndProfesionalPartner) {
	gc.DB.Where("id = ?", uint64(a)).First(&lp)

	return
}

func (a RatingAgencyID) GetRatingAgency(gc *sharedconfig.GlobalConfig) (ra RatingAgency) {
	gc.DB.Where("id = ?", uint64(a)).First(&ra)

	return
}

func (a TrusteeID) GetTrustee(gc *sharedconfig.GlobalConfig) (ts Trustee) {
	gc.DB.Where("id = ?", uint64(a)).First(&ts)

	return
}

// ne retruns true is string s is not empty
func ne(s string) bool {
	return len(s) > 0
}

// nn returns true is *string s is not nil
func NN(s *string) bool {
	return s != nil

}

// // asNN assigns the source to destination if source is not nil and empties destination if source is nil
// func asNN(source *string, destination *string) (r bool) {
// 	if destination != nil {
// 		r = false
// 	}
// 	if nn(source) {
// 		destination = source
// 	} else {
// 		empty := ""
// 		destination = &empty
// 	}
// 	return destination == nil
// }

// // asNE assigns the source to destination if source is not empty and nils destination if source is empty
// func asNE(source string, destination *string) (r bool) {
// 	if destination != nil {
// 		r = false
// 	}
// 	if ne(source) {
// 		destination = &source
// 	} else {
// 		destination = nil
// 	}
// 	return destination == nil
// }

// GetMarketOffer fetches the offer information using public key
func (t *TokenizedAsset) GetMarketOffers(gc *sharedconfig.GlobalConfig) (marketOffers horizon.OffersPage, err error) {
	if t.MarketMakingWallet == nil {
		return
	}
	seller := *t.MarketMakingWallet
	client := network.GetBlockchainClient()
	// var offerRequest horizonclient.OfferRequest

	//real account
	offerRequest := horizonclient.OfferRequest{Seller: seller}

	marketOffers, err = client.Offers(offerRequest)
	if err != nil {
		// log.Printf("[GetBlockchainAccountDetail]: %v, error: [%v]", accountRequest.AccountID, err)
		if strings.Contains(err.Error(), "timeout") || strings.Contains(err.Error(), "handshake") || strings.Contains(err.Error(), "no such host") || strings.Contains(err.Error(), "timeout") || strings.Contains(err.Error(), "dial") {
			log.Printf("[GetMarketOffer Network Failure]: %s\n", "Error Connecting to Expansion Service")
			return marketOffers, &tErrors.ErrorTemporaryServerError{}
		} else {
			horizonException, ok := err.(*horizonclient.Error)

			if ok {

				if horizonException.Problem.Status == http.StatusNotFound {
					return marketOffers, &tErrors.ErrorBlockchainAccountNotActivated{}
				}
				log.Printf("[GetMarketOffer] error is known. Type: %v, Status: %v, Detail: %v, Title: %v, Extras: %v", horizonException.Problem.Type, horizonException.Problem.Status, horizonException.Problem.Detail, horizonException.Problem.Title, horizonException.Problem.Extras)
			}

		}
		return marketOffers, &tErrors.ErrorTemporaryServerError{}
	}

	return marketOffers, nil
}
