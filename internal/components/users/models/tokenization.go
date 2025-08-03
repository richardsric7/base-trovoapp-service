package users

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"log"
	"strings"
	"time"
	"trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
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
	NumberOfTokenToBeSold                        float64                         `json:"numberOfTokenToBeSold"`
	TotalTokenHeldByManager                      float64                         `json:"totalTokenHeldByManager"`
	WalletToHoldAssetsNotForSale                 string                          `json:"walletToHoldAssetsNotForSale"`
	PricePerToken                                float64                         `json:"pricePerToken"`
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

	for key, _ := range knownAssets {
		untrusted = append(untrusted, key)
	}

	return untrusted
}

func (p *ProceedPayout) CreateBatch() error {
	if len(p.TokenizedAssetID) == 0 {
		return &errors.CustomError{Err: "error invalid tokenizedAssetId", ErrMessage: "tokenized asset identification is invalid"}
	}
	if p.TokenizedAsset.ProceedCycle == nil {
		return &errors.CustomError{Err: "error proceed-cycle-not-set", ErrMessage: "Proceed Cycle was not set for this project"}

	}
	if *p.TokenizedAsset.ProceedCycle != "None" {
		return &errors.CustomError{Err: "error proceed-cycle-not-set", ErrMessage: "Proceed Cycle was not set for this project"}

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
		err = &errors.ErrorTemporaryServerError{}
		return
	}

	err = gc.DB.Preload(clause.Associations).Where("Tokenized_Asset_ID = ? AND Subscriber_Username = ?", t.ID, subscriber).First(&exp).Error

	return
}

func (t *TokenizedAsset) CountExpressedInterests(gc *sharedconfig.GlobalConfig) (count int64) {
	if t == nil {
		log.Println("[TokenizedAsset::CountExpressedInterests] Error tokenized asset is nil")
		return
	}

	gc.DB.Model(ExpressionOfInterest{}).Where("Tokenized_Asset_ID = ?", t.ID).Count(&count)

	return
}

func (t *TokenizedAsset) CountNumberOfSubscribers(gc *sharedconfig.GlobalConfig) (count int64) {
	if t == nil {
		log.Println("[TokenizedAsset::CountNumberOfSubscribers] Error tokenized asset is nil")
		return
	}

	gc.DB.Model(TokenizedAssetSubscription{}).Where("Tokenized_Asset_ID = ?", t.ID).Count(&count)

	return
}

// SumQuantitySold suma the total amount in fiat sold so far
func (t *TokenizedAsset) SumQuantitySold(gc *sharedconfig.GlobalConfig) (sum float64) {
	if t == nil {
		log.Println("[TokenizedAsset::SumQuantitySold] Error tokenized asset is nil")
		return
	}
	gc.DB.Model(TokenizedAssetSubscription{}).Where("Tokenized_Asset_ID = ?", t.ID).Select("case when sum(Amount) is not null then sum(Amount) else 0 end").Row().Scan(&sum)

	return
}

// SumAmountBoughtByWalletOwner sums the amout that is purchased by the wallet owner
func (t *TokenizedAsset) SumAmountBoughtByWalletOwner(walletAlias string, gc *sharedconfig.GlobalConfig) (sum float64) {
	if t == nil {
		log.Println("[TokenizedAsset::SumAmountBoughtByWalletOwner] Error tokenized asset is nil")
		return
	}
	ownerUsername := strings.Split(walletAlias, "_")[0]

	gc.DB.Model(TokenizedAssetSubscription{}).Where("Tokenized_Asset_ID = ? AND Wallet_Alias LIKE ?", t.ID, strings.ToLower(ownerUsername)+"%").Select("case when sum(Amount) is not null then sum(Amount) else 0 end").Row().Scan(&sum)

	return
}

func (t *TokenizedAsset) GetTokenizedAssetSubscriptionByWalletPublicKey(subscriberWalletPublicKey string, gc *sharedconfig.GlobalConfig) (sub TokenizedAssetSubscription, err error) {
	if t == nil {
		log.Printf("[TokenizedAsset::GetTokenizedAssetSubscriptionByWalletPublicKey] Error tokenized asset is nil %v\n", subscriberWalletPublicKey)
		err = &errors.ErrorTemporaryServerError{}
		return
	}

	err = gc.DB.Preload(clause.Associations).Where("Tokenized_Asset_ID = ? AND Wallet_Public_Key = ?", t.ID, subscriberWalletPublicKey).First(&sub).Error

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
	} else {

		t.TokenizationApplicationFee = t.CountryConfig.TokenizationApplicationFee
		t.TokenizationApplicationFeeAsset = t.CountryConfig.TokenizationApplicationFeeAsset

	}

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
	t.QuantityOfTokensSold = ti.SumQuantitySold(gc)

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
func nn(s *string) bool {
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
