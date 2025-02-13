package users

import (
	"log"
	"time"
	"trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/shopspring/decimal"
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
	ID                                          string                          `json:"id"`
	CreatedAt                                   time.Time                       `json:"createdAt"`
	UpdatedAt                                   time.Time                       `json:"updatedAt"`
	InitiatorUsername                           string                          `gorm:"size:50;not null" json:"initiatorUsername"`
	AssetSector                                 *string                         `json:"assetSector"`
	AssetSubSector                              *string                         `json:"assetSubSector"`
	AssetType                                   *string                         `json:"assetType"`
	AssetName                                   *string                         `json:"assetName"`
	ApprovedAssetCustodianID                    uint64                          `gorm:"not null" json:"approvedAssetCustodianId"`
	ApprovedAssetCustodian                      ApprovedAssetCustodian          `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"approvedAssetCustodianInfo"`
	OfferingType                                *string                         `gorm:"default:'PRIVATE'" json:"offeringType"` //PRIVATE, PUBLIC
	ClosedGroupID                               *string                         `gorm:"null" json:"closedGroupId"`
	ClosedGroup                                 ClosedGroup                     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"closedGroupInfo"`
	SecApproval                                 int                             `gorm:"default:0" json:"secApproval"`
	SecApprovalIdNumber                         *string                         `json:"secApprovalIdNumber"`
	IssuingWalletPublicKey                      *string                         `gorm:"size:60" json:"issuingWalletPublicKey"`
	IssuingWalletAlias                          *string                         `gorm:"size:60" json:"issuingWalletAlias"`
	MarketMakingWallet                          *string                         `json:"marketMakingWallet"`
	AssetDescription                            *string                         `json:"assetDescription"`
	AssetCountryLocation                        *string                         `json:"assetCountryLocation"`
	AssetPhysicalAddress                        *string                         `json:"assetPhysicalAddress"`
	AssetLongitude                              *string                         `json:"assetLongitude"`
	AssetLatitude                               *string                         `json:"assetLatitude"`
	OwnershipType                               *string                         `json:"ownershipType"` //DIRECT,THIRD-PARTY
	OwnershipKind                               *string                         `json:"ownershipKind"` //INDIVIDUAL,CORPORATE,
	InitialOwnerPreferredWalletAddress          *string                         `json:"initialOwnerPreferredWalletAddress"`
	AssetOwnerName                              *string                         `json:"assetOwnerName"`
	AssetOwnerAddress                           *string                         `json:"assetOwnerAddress"`
	AssetManagerID                              uint64                          `json:"assetManagerId"`
	AssetManager                                AssetManager                    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"assetManagerInfo"`
	AssetQuoteCurrency                          *string                         `gorm:"default:'NGN'" json:"assetQuoteCurrency"`
	AssetCurrentValue                           float64                         `gorm:"default:0" json:"assetCurrentValue"`
	AssetOwnerRetainedOrContributedValue        float64                         `gorm:"default:0" json:"assetOwnerRetainedOrContributedValue"`
	AssetMscCostOutisdeOfValuation              float64                         `gorm:"default:0" json:"assetMscCostOutisdeOfValuation"`
	ValueOfTokenizedAsset                       float64                         `gorm:"default:0" json:"valueOfTokenizedAsset"`
	ProtectionMethods                           *string                         `json:"protectionMethods"` //csv format
	InsuranceCompanyName                        *string                         `json:"insuranceCompanyName"`
	InsurancePolicyNumber                       *string                         `json:"insurancePolicyNumber"`
	InsurancePolicyHolder                       *string                         `json:"insurancePolicyHolder"`
	PercentageValueOfInsurance                  float64                         `gorm:"default:0" json:"percentageValueOfInsurance"`
	IsFreeFromLiensAndEncumbrances              int                             `gorm:"default:0" json:"IsFreeFromLiensAndEncumbrances"`
	AssetAlreadyExists                          int                             `gorm:"default:1" json:"assetAlreadyExists"`
	VettingStatus                               int                             `gorm:"default:0" json:"vettingStatus"`
	AssetTokenizationDocuments                  []AssetTokenizationDocument     `json:"AssetTokenizationDocuments"`
	ProofOfPaymentDocuments                     []TokenizationFeeProofOfPayment `json:"ProofOfPaymentDocuments"`
	AssetCode                                   *string                         `json:"assetCode"`
	AssetLogo                                   *string                         `json:"assetLogo"`
	AssetWebsite                                *string                         `gorm:"default:'trovotech.io'" json:"assetWebsite"`
	NumberOfTokenToBeIssued                     float64                         `gorm:"default:0" json:"numberOfTokenToBeIssued"`
	MaxNumberOfTokenAvailableForSale            float64                         `gorm:"default:0" json:"maxNumberOfTokenAvailableForSale"`
	FeeInAsset                                  float64                         `gorm:"default:0" json:"feeInAsset"`
	FeeInFiat                                   float64                         `gorm:"default:0" json:"feeInFiat"`
	NumberOfTokenToBeSold                       float64                         `gorm:"default:0" json:"numberOfTokenToBeSold"`
	TotalTokenHeldByManager                     float64                         `gorm:"default:0" json:"totalTokenHeldByManager"`
	WalletToHoldAssetsNotForSale                *string                         `json:"walletToHoldAssetsNotForSale"`
	PricePerToken                               float64                         `gorm:"default:0" json:"pricePerToken"`
	SalesStart                                  time.Time                       `json:"salesStart"`
	SalesEnd                                    time.Time                       `json:"salesEnd"`
	CapOnPurchase                               float64                         `gorm:"default:0" json:"capOnPurchase"`
	CapQuantity                                 float64                         `gorm:"default:0" json:"capQuantity"`
	CapDurationInDays                           int                             `gorm:"default:0" json:"capDurationInDays"`
	ProceedCycle                                *string                         `gorm:"size:50" json:"proceedCycle"`
	TokenizationFeeID                           *uint64                         `gorm:"default:0" json:"tokenizationFeeId"`
	TokenizationFee                             TokenizationFee                 `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"tokenizationFee"`
	SECTokenizationFeePercent                   float64                         `gorm:"default:0" json:"SECTokenizationFeePercent"`
	SECTokenizationFeeValue                     float64                         `gorm:"default:0" json:"SECTokenizationFeeValue"`
	CustodianFeePercent                         float64                         `gorm:"default:0" json:"custodianFeePercent"`
	CustodianFeeValue                           float64                         `gorm:"default:0" json:"custodianFeeValue"`
	AssetManagerFeeValue                        float64                         `gorm:"default:0" json:"assetManagerFeeValue"`
	AssetManagerFeePercent                      float64                         `gorm:"default:0" json:"assetManagerFeePercent"`
	ProceedPayoutCurrency                       *string                         `json:"proceedPayoutCurrency"`
	ProceedPayoutType                           int                             `gorm:"default:0" json:"proceedPayoutType"` // FIAT=1, CRYPTO=0
	ExemptedCountries                           *string                         `json:"exemptedCountries"`
	HasAdditionalKYCRequirements                int                             `gorm:"default:0" json:"hasAdditionalKYCRequirements"`
	AdditionalKYCRequirements                   *string                         `json:"additionalKYCRequirements"`
	InvestorAccreditationRequired               int                             `gorm:"default:0" json:"investorAccreditationRequired"`
	AssetTokenizationStatus                     int                             `gorm:"default:0" json:"assetTokenizationStatus"`
	LastUpdatedBy                               *string                         `gorm:"null" json:"lastUpdatedBy"`
	TokenizationTransaction                     *string                         `gorm:"null" json:"tokenizationTransaction"`
	AgreeTransferTitleToCustodian               int                             `gorm:"default:0" json:"agreeTransferTitleToCustodian"`
	ContractualProtectionRevGuarantees          int                             `gorm:"default:0" json:"contractualProtectionRevGuarantees"`
	ContractualProtectionPerfBond               int                             `gorm:"default:0" json:"contractualProtectionPerfBond"`
	ContractualProtectionSLA                    int                             `gorm:"default:0" json:"contractualProtectionSLA"`
	RiskSharingMechanismPPPs                    int                             `gorm:"default:0" json:"riskSharingMechanismPPPs"`
	RiskSharingMechanismHedgeInstruments        int                             `gorm:"default:0" json:"riskSharingMechanismHedgeInstruments"`
	RiskSharingMechanismCompletionGuarantees    int                             `gorm:"default:0" json:"riskSharingMechanismCompletionGuarantees"`
	IndependentMonitoringList                   *string                         `gorm:"null" json:"independentMonitoringList"`
	ESGSafeguardsSusCerts                       int                             `gorm:"default:0" json:"eSGSafeguardsSusCerts"`
	ESGSafeguardsCommEngPlans                   int                             `gorm:"default:0" json:"eSGSafeguardsCommEngPlans"`
	SecurityMeasuresAccessControl               int                             `gorm:"default:0" json:"securityMeasuresAccessControl"`
	SecurityMeasuresSurveilanceSystems          int                             `gorm:"default:0" json:"securityMeasuresSurveilanceSystems"`
	SecurityMeasuresOnSiteSecurityPersonnel     int                             `gorm:"default:0" json:"securityMeasuresOnSiteSecurityPersonnel"`
	SecurityMeasuresPerimeterSecurity           int                             `gorm:"default:0" json:"securityMeasuresPerimeterSecurity"`
	SecurityMeasuresCriticalInfraProtections    int                             `gorm:"default:0" json:"securityMeasuresCriticalInfraProtections"`
	OtherAssetProtection                        *string                         `gorm:"null" json:"otherAssetProtection"`
	LegalAdvisor                                *string                         `gorm:"null" json:"legalAdvisor"`
	FinancialAdvisor                            *string                         `gorm:"null" json:"financialAdvisor"`
	UndertakingNoLien                           int                             `gorm:"default:0" json:"undertakingNoLien"`
	UndertakingNotCollateral                    int                             `gorm:"default:0" json:"undertakingNotCollateral"`
	UndertakingNoClaims                         int                             `gorm:"default:0" json:"undertakingNoClaims"`
	UndertakingNoForeclosure                    int                             `gorm:"default:0" json:"undertakingNoForeclosure"`
	ComplianceNoViolation                       int                             `gorm:"default:0" json:"complianceNoViolation"`
	ComplianceAllPermits                        int                             `gorm:"default:0" json:"complianceAllPermits"`
	OutstandingFinancialRespNoDebts             int                             `gorm:"default:0" json:"outstandingFinancialRespNoDebts"`
	OutstandingFinancialRespNoHiddenLiabilities int                             `gorm:"default:0" json:"outstandingFinancialRespNoHiddenLiabilities"`
	RiskManagementFullyInsured                  int                             `gorm:"default:0" json:"riskManagementFullyInsured"`
	RiskManagementDeclaredValue                 int                             `gorm:"default:0" json:"riskManagementDeclaredValue"`
	PhysicalConditionSound                      int                             `gorm:"default:0" json:"physicalConditionSound"`
	PhysicalConditionNolease                    int                             `gorm:"default:0" json:"physicalConditionNolease"`
	PhysicalConditionNoUndisclosedEasements     int                             `gorm:"default:0" json:"physicalConditionNoUndisclosedEasements"`
}

type TokenizedAssetID string

type TokenizedAssetJSONInput struct {
	AssetSector                                 string       `json:"assetSector"`
	AssetSubSector                              string       `json:"assetSubSector"`
	AssetType                                   string       `json:"assetType"`
	AssetName                                   string       `json:"assetName"`
	AssetWebsite                                string       `json:"assetWebsite"`
	ApprovedAssetCustodianID                    uint64       `gorm:"not null" json:"approvedAssetCustodianId"`
	OfferingType                                string       `gorm:"default:'PRIVATE'" json:"offeringType"` //PRIVATE, PUBLIC
	ClosedGroupID                               string       `gorm:"null" json:"closedGroupId"`
	SecApproval                                 int          `gorm:"default:0" json:"secApproval"`
	SecApprovalIdNumber                         string       `json:"secApprovalIdNumber"`
	MarketMakingWallet                          string       `json:"marketMakingWallet"`
	AssetDescription                            string       `json:"assetDescription"`
	AssetCountryLocation                        string       `json:"assetCountryLocation"`
	AssetPhysicalAddress                        string       `json:"assetPhysicalAddress"`
	AssetLongitude                              string       `json:"assetLongitude"`
	AssetLatitude                               string       `json:"assetLatitude"`
	OwnershipType                               string       `json:"ownershipType"` //DIRECT,THIRD-PARTY
	OwnershipKind                               string       `json:"ownershipKind"` //INDIVIDUAL,CORPORATE,
	InitialOwnerPreferredWalletAddress          string       `json:"initialOwnerPreferredWalletAddress"`
	AssetOwnerName                              string       `json:"assetOwnerName"`
	AssetOwnerRetainedOrContributedValue        float64      `json:"assetOwnerRetainedOrContributedValue"`
	AssetOwnerAddress                           string       `json:"assetOwnerAddress"`
	AssetManagerID                              uint64       `json:"assetManagerId"`
	AssetManager                                AssetManager `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"assetManagerInfo"`
	AssetQuoteCurrency                          string       `gorm:"default:'NGN'" json:"assetQuoteCurrency"`
	AssetCurrentValue                           float64      `gorm:"default:0" json:"assetCurrentValue"`
	AssetMscCostOutisdeOfValuation              float64      `gorm:"default:0" json:"assetMscCostOutisdeOfValuation"`
	ProtectionMethods                           string       `json:"protectionMethods"` //csv format
	InsuranceCompanyName                        string       `json:"insuranceCompanyName"`
	InsurancePolicyNumber                       string       `json:"insurancePolicyNumber"`
	InsurancePolicyHolder                       string       `json:"insurancePolicyHolder"`
	PercentageValueOfInsurance                  float64      `gorm:"default:0" json:"percentageValueOfInsurance"`
	IsFreeFromLiensAndEncumbrances              int          `gorm:"default:0" json:"IsFreeFromLiensAndEncumbrances"`
	AssetAlreadyExists                          int          `gorm:"default:1" json:"assetAlreadyExists"`
	AssetCode                                   string       `json:"assetCode"`
	AssetLogo                                   string       `json:"assetLogo"`
	NumberOfTokenToBeIssued                     float64      `json:"numberOfTokenToBeIssued"`
	NumberOfTokenToBeSold                       float64      `json:"numberOfTokenToBeSold"`
	TotalTokenHeldByManager                     float64      `json:"totalTokenHeldByManager"`
	WalletToHoldAssetsNotForSale                string       `json:"walletToHoldAssetsNotForSale"` //wallet that the original owner wants to use to receive their portion of tokenized asset that are not meant for sale.
	PricePerToken                               float64      `json:"pricePerToken"`
	SalesStart                                  time.Time    `json:"salesStart"`
	SalesEnd                                    time.Time    `json:"salesEnd"`
	CapOnPurchase                               float64      `gorm:"default:0" json:"capOnPurchase"`
	CapQuantity                                 float64      `gorm:"default:0" json:"capQuantity"`
	CapDurationInDays                           int          `gorm:"default:0" json:"capDurationInDays"`
	ProceedCycle                                string       `gorm:"size:50" json:"proceedCycle"`
	TokenizationFeeID                           uint64       `json:"tokenizationFeeId"`
	ProceedPayoutCurrency                       string       `json:"proceedPayoutCurrency"`
	ProceedPayoutType                           int          `gorm:"default:0" json:"proceedPayoutType"` // FIAT=1, CRYPTO=0. FIAT requires fiat payment method.
	ExemptedCountries                           string       `json:"exemptedCountries"`
	HasAdditionalKYCRequirements                int          `gorm:"default:0" json:"hasAdditionalKYCRequirements"`
	AdditionalKYCRequirements                   string       `json:"additionalKYCRequirements"`
	InvestorAccreditationRequired               int          `gorm:"default:0" json:"investorAccreditationRequired"`
	Messages                                    []string     `json:"messages"`
	AgreeTransferTitleToCustodian               int          `gorm:"default:0" json:"agreeTransferTitleToCustodian"`
	ContractualProtectionRevGuarantees          int          `gorm:"default:0" json:"contractualProtectionRevGuarantees"`
	ContractualProtectionPerfBond               int          `gorm:"default:0" json:"contractualProtectionPerfBond"`
	ContractualProtectionSLA                    int          `gorm:"default:0" json:"contractualProtectionSLA"`
	RiskSharingMechanismPPPs                    int          `gorm:"default:0" json:"riskSharingMechanismPPPs"`
	RiskSharingMechanismHedgeInstruments        int          `gorm:"default:0" json:"riskSharingMechanismHedgeInstruments"`
	RiskSharingMechanismCompletionGuarantees    int          `gorm:"default:0" json:"riskSharingMechanismCompletionGuarantees"`
	IndependentMonitoringList                   string       `gorm:"null" json:"independentMonitoringList"`
	ESGSafeguardsSusCerts                       int          `gorm:"default:0" json:"eSGSafeguardsSusCerts"`
	ESGSafeguardsCommEngPlans                   int          `gorm:"default:0" json:"eSGSafeguardsCommEngPlans"`
	SecurityMeasuresAccessControl               int          `gorm:"default:0" json:"securityMeasuresAccessControl"`
	SecurityMeasuresSurveilanceSystems          int          `gorm:"default:0" json:"securityMeasuresSurveilanceSystems"`
	SecurityMeasuresOnSiteSecurityPersonnel     int          `gorm:"default:0" json:"securityMeasuresOnSiteSecurityPersonnel"`
	SecurityMeasuresPerimeterSecurity           int          `gorm:"default:0" json:"securityMeasuresPerimeterSecurity"`
	SecurityMeasuresCriticalInfraProtections    int          `gorm:"default:0" json:"securityMeasuresCriticalInfraProtections"`
	OtherAssetProtection                        string       `gorm:"null" json:"otherAssetProtection"`
	LegalAdvisor                                string       `gorm:"null" json:"legalAdvisor"`
	FinancialAdvisor                            string       `gorm:"null" json:"financialAdvisor"`
	UndertakingNoLien                           int          `gorm:"default:0" json:"undertakingNoLien"`
	UndertakingNotCollateral                    int          `gorm:"default:0" json:"undertakingNotCollateral"`
	UndertakingNoClaims                         int          `gorm:"default:0" json:"undertakingNoClaims"`
	UndertakingNoForeclosure                    int          `gorm:"default:0" json:"undertakingNoForeclosure"`
	ComplianceNoViolation                       int          `gorm:"default:0" json:"complianceNoViolation"`
	ComplianceAllPermits                        int          `gorm:"default:0" json:"complianceAllPermits"`
	OutstandingFinancialRespNoDebts             int          `gorm:"default:0" json:"outstandingFinancialRespNoDebts"`
	OutstandingFinancialRespNoHiddenLiabilities int          `gorm:"default:0" json:"outstandingFinancialRespNoHiddenLiabilities"`
	RiskManagementFullyInsured                  int          `gorm:"default:0" json:"riskManagementFullyInsured"`
	RiskManagementDeclaredValue                 int          `gorm:"default:0" json:"riskManagementDeclaredValue"`
	PhysicalConditionSound                      int          `gorm:"default:0" json:"physicalConditionSound"`
	PhysicalConditionNolease                    int          `gorm:"default:0" json:"physicalConditionNolease"`
	PhysicalConditionNoUndisclosedEasements     int          `gorm:"default:0" json:"physicalConditionNoUndisclosedEasements"`
}

type ConfirmTokenizedAssetJSONInput struct {
	Messages             []string `json:"messages"`
	Transaction          string   `json:"transaction"`
	TransactionSignature string   `json:"transactionSignature"`
	TransactionID        string   `json:"transactionId"`
	NetworkPassPhrase    string   `json:"networkPassPhrase"`
}

type VetTokenizedAssetJSONInput struct {
	ApprovedAssetCustodianID uint64   `gorm:"not null" json:"approvedAssetCustodianId"`
	AssetManagerID           uint64   `json:"assetManagerId"`
	CountryCode              string   `json:"CountryCode"`
	Messages                 []string `json:"messages"`
}

type TokenizedAssetJSON struct {
	ID                                          string                          `json:"id"`
	CreatedAt                                   time.Time                       `json:"createdAt"`
	UpdatedAt                                   time.Time                       `json:"updatedAt"`
	InitiatorUsername                           string                          `gorm:"size:50;not null" json:"initiatorUsername"`
	AssetSector                                 string                          `json:"assetSector"`
	AssetSubSector                              string                          `json:"assetSubSector"`
	AssetType                                   string                          `json:"assetType"`
	AssetName                                   string                          `json:"assetName"`
	AssetWebsite                                string                          `json:"assetWebsite"`
	ApprovedAssetCustodianID                    uint64                          `gorm:"not null" json:"approvedAssetCustodianId"`
	ApprovedAssetCustodian                      ApprovedAssetCustodian          `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"approvedAssetCustodianInfo"`
	OfferingType                                string                          `gorm:"default:'PRIVATE'" json:"offeringType"` //PRIVATE, PUBLIC
	ClosedGroupID                               string                          `gorm:"null" json:"closedGroupId"`
	ClosedGroup                                 ClosedGroup                     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"closedGroupInfo"`
	SecApproval                                 int                             `gorm:"default:0" json:"secApproval"`
	SecApprovalIdNumber                         string                          `json:"secApprovalIdNumber"`
	IssuingWalletPublicKey                      string                          `gorm:"size:60" json:"issuingWalletPublicKey"`
	IssuingWalletAlias                          string                          `gorm:"size:60" json:"issuingWalletAlias"`
	MarketMakingWallet                          string                          `json:"marketMakingWallet"`
	AssetDescription                            string                          `json:"assetDescription"`
	AssetCountryLocation                        string                          `json:"assetCountryLocation"`
	AssetPhysicalAddress                        string                          `json:"assetPhysicalAddress"`
	AssetLongitude                              string                          `json:"assetLongitude"`
	AssetLatitude                               string                          `json:"assetLatitude"`
	OwnershipType                               string                          `json:"ownershipType"` //DIRECT,THIRD-PARTY
	OwnershipKind                               string                          `json:"ownershipKind"` //INDIVIDUAL,CORPORATE,
	InitialOwnerPreferredWalletAddress          string                          `json:"initialOwnerPreferredWalletAddress"`
	AssetOwnerName                              string                          `json:"assetOwnerName"`
	AssetOwnerRetainedOrContributedValue        float64                         `json:"assetOwnerRetainedOrContributedValue"`
	AssetOwnerAddress                           string                          `json:"assetOwnerAddress"`
	AssetManagerID                              uint64                          `json:"assetManagerId"`
	AssetManager                                AssetManager                    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"assetManagerInfo"`
	AssetQuoteCurrency                          string                          `gorm:"default:'NGN'" json:"assetQuoteCurrency"`
	AssetCurrentValue                           float64                         `gorm:"default:0" json:"assetCurrentValue"`
	AssetMscCostOutisdeOfValuation              float64                         `gorm:"default:0" json:"assetMscCostOutisdeOfValuation"`
	ValueOfTokenizedAsset                       float64                         `gorm:"default:0" json:"valueOfTokenizedAsset"`
	ProtectionMethods                           string                          `json:"protectionMethods"` //csv format
	InsuranceCompanyName                        string                          `json:"insuranceCompanyName"`
	InsurancePolicyNumber                       string                          `json:"insurancePolicyNumber"`
	InsurancePolicyHolder                       string                          `json:"insurancePolicyHolder"`
	PercentageValueOfInsurance                  float64                         `gorm:"default:0" json:"percentageValueOfInsurance"`
	IsFreeFromLiensAndEncumbrances              int                             `gorm:"default:0" json:"IsFreeFromLiensAndEncumbrances"`
	AssetAlreadyExists                          int                             `gorm:"default:1" json:"assetAlreadyExists"`
	VettingStatus                               int                             `gorm:"default:0" json:"vettingStatus"`
	AssetTokenizationDocuments                  []AssetTokenizationDocument     `json:"AssetTokenizationDocuments"`
	ProofOfPaymentDocuments                     []TokenizationFeeProofOfPayment `json:"ProofOfPaymentDocuments"`
	AssetCode                                   string                          `json:"assetCode"`
	AssetLogo                                   string                          `json:"assetLogo"`
	NumberOfTokenToBeIssued                     float64                         `json:"numberOfTokenToBeIssued"`
	MaxNumberOfTokenAvailableForSale            float64                         `gorm:"default:0" json:"maxNumberOfTokenAvailableForSale"`
	FeeInAsset                                  float64                         `gorm:"default:0" json:"feeInAsset"`
	FeeInFiat                                   float64                         `gorm:"default:0" json:"feeInFiat"`
	NumberOfTokenToBeSold                       float64                         `json:"numberOfTokenToBeSold"`
	TotalTokenHeldByManager                     float64                         `json:"totalTokenHeldByManager"`
	WalletToHoldAssetsNotForSale                string                          `json:"walletToHoldAssetsNotForSale"`
	PricePerToken                               float64                         `json:"pricePerToken"`
	SalesStart                                  time.Time                       `json:"salesStart"`
	SalesEnd                                    time.Time                       `json:"salesEnd"`
	CapOnPurchase                               float64                         `gorm:"default:0" json:"capOnPurchase"`
	CapQuantity                                 float64                         `gorm:"default:0" json:"capQuantity"`
	CapDurationInDays                           int                             `gorm:"default:0" json:"capDurationInDays"`
	ProceedCycle                                string                          `gorm:"size:50" json:"proceedCycle"`
	TokenizationFeeID                           uint64                          `json:"tokenizationFeeId"`
	TokenizationFee                             TokenizationFee                 `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"tokenizationFee"`
	SECTokenizationFeePercent                   float64                         `gorm:"default:0" json:"SECTokenizationFeePercent"`
	SECTokenizationFeeValue                     float64                         `gorm:"default:0" json:"SECTokenizationFeeValue"`
	CustodianFeePercent                         float64                         `gorm:"default:0" json:"custodianFeePercent"`
	CustodianFeeValue                           float64                         `gorm:"default:0" json:"custodianFeeValue"`
	AssetManagerFeeValue                        float64                         `gorm:"default:0" json:"assetManagerFeeValue"`
	AssetManagerFeePercent                      float64                         `gorm:"default:0" json:"assetManagerFeePercent"`
	ProceedPayoutCurrency                       string                          `json:"proceedPayoutCurrency"`
	ProceedPayoutType                           int                             `gorm:"default:0" json:"proceedPayoutType"` // FIAT=1, CRYPTO=0
	ExemptedCountries                           string                          `json:"exemptedCountries"`
	HasAdditionalKYCRequirements                int                             `gorm:"default:0" json:"hasAdditionalKYCRequirements"`
	AdditionalKYCRequirements                   string                          `json:"additionalKYCRequirements"`
	InvestorAccreditationRequired               int                             `gorm:"default:0" json:"investorAccreditationRequired"`
	AssetTokenizationStatus                     int                             `gorm:"default:0" json:"assetTokenizationStatus"`
	LastUpdatedBy                               string                          `gorm:"null" json:"lastUpdatedBy"`
	AgreeTransferTitleToCustodian               int                             `gorm:"default:0" json:"agreeTransferTitleToCustodian"`
	ContractualProtectionRevGuarantees          int                             `gorm:"default:0" json:"contractualProtectionRevGuarantees"`
	ContractualProtectionPerfBond               int                             `gorm:"default:0" json:"contractualProtectionPerfBond"`
	ContractualProtectionSLA                    int                             `gorm:"default:0" json:"contractualProtectionSLA"`
	RiskSharingMechanismPPPs                    int                             `gorm:"default:0" json:"riskSharingMechanismPPPs"`
	RiskSharingMechanismHedgeInstruments        int                             `gorm:"default:0" json:"riskSharingMechanismHedgeInstruments"`
	RiskSharingMechanismCompletionGuarantees    int                             `gorm:"default:0" json:"riskSharingMechanismCompletionGuarantees"`
	IndependentMonitoringList                   string                          `gorm:"null" json:"independentMonitoringList"`
	ESGSafeguardsSusCerts                       int                             `gorm:"default:0" json:"eSGSafeguardsSusCerts"`
	ESGSafeguardsCommEngPlans                   int                             `gorm:"default:0" json:"eSGSafeguardsCommEngPlans"`
	SecurityMeasuresAccessControl               int                             `gorm:"default:0" json:"securityMeasuresAccessControl"`
	SecurityMeasuresSurveilanceSystems          int                             `gorm:"default:0" json:"securityMeasuresSurveilanceSystems"`
	SecurityMeasuresOnSiteSecurityPersonnel     int                             `gorm:"default:0" json:"securityMeasuresOnSiteSecurityPersonnel"`
	SecurityMeasuresPerimeterSecurity           int                             `gorm:"default:0" json:"securityMeasuresPerimeterSecurity"`
	SecurityMeasuresCriticalInfraProtections    int                             `gorm:"default:0" json:"securityMeasuresCriticalInfraProtections"`
	OtherAssetProtection                        string                          `gorm:"null" json:"otherAssetProtection"`
	LegalAdvisor                                string                          `gorm:"null" json:"legalAdvisor"`
	FinancialAdvisor                            string                          `gorm:"null" json:"financialAdvisor"`
	UndertakingNoLien                           int                             `gorm:"default:0" json:"undertakingNoLien"`
	UndertakingNotCollateral                    int                             `gorm:"default:0" json:"undertakingNotCollateral"`
	UndertakingNoClaims                         int                             `gorm:"default:0" json:"undertakingNoClaims"`
	UndertakingNoForeclosure                    int                             `gorm:"default:0" json:"undertakingNoForeclosure"`
	ComplianceNoViolation                       int                             `gorm:"default:0" json:"complianceNoViolation"`
	ComplianceAllPermits                        int                             `gorm:"default:0" json:"complianceAllPermits"`
	OutstandingFinancialRespNoDebts             int                             `gorm:"default:0" json:"outstandingFinancialRespNoDebts"`
	OutstandingFinancialRespNoHiddenLiabilities int                             `gorm:"default:0" json:"outstandingFinancialRespNoHiddenLiabilities"`
	RiskManagementFullyInsured                  int                             `gorm:"default:0" json:"riskManagementFullyInsured"`
	RiskManagementDeclaredValue                 int                             `gorm:"default:0" json:"riskManagementDeclaredValue"`
	PhysicalConditionSound                      int                             `gorm:"default:0" json:"physicalConditionSound"`
	PhysicalConditionNolease                    int                             `gorm:"default:0" json:"physicalConditionNolease"`
	PhysicalConditionNoUndisclosedEasements     int                             `gorm:"default:0" json:"physicalConditionNoUndisclosedEasements"`
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
	FeeFiatCap         float64 `json:"feeFiatCap"`
	FeeAssetPercentage float64 `json:"feeAssetPercentage"`
	// FeeAssetCap        float64 `json:"feeAssetCap"`
	FeeDescription string `json:"feeDescription"`
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
}
type AssetManager struct {
	ID                  uint64  `gorm:"" json:"id"`
	AssetManagerName    string  `gorm:"size:100" json:"assetManagerName"`
	AssetManagerAddress string  `json:"assetManagerAddress"`
	AssetManagerCountry string  `gorm:"size:3" json:"assetManagerCountry"`
	FeePercent          float64 `gorm:"default:0" json:"FeePercent"`
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

type ExpressionOfInterest struct {
	ID                 uint64         `gorm:"" json:"-" form:"-"`
	CreatedAt          time.Time      `json:"createdAt"`
	UpdatedAt          time.Time      `json:"updatedAt"`
	TokenizedAssetID   string         `gorm:"not null;size:100" json:"tokenizedAssetId"`
	TokenizedAsset     TokenizedAsset `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"tokenizedAssetInfo"`
	AssetCode          string         `gorm:"not null;size:12" json:"assetCode"`
	AssetIssuer        string         `gorm:"not null;size:100" json:"assetIssuer"`
	WalletAlias        string         `gorm:"not null;size:100" json:"WalletAlias"`
	WalletPublicKey    string         `gorm:"not null;size:100" json:"WalletPublicKey"`
	Amount             float64        `json:"amount"`
	Price              float64        `json:"price"`
	SubscriberUsername string         `gorm:"not null;size:100" json:"SubscriberUsername"`
}

type ExpressionOfInterestInput struct {
	Amount float64 `json:"amount"`
}

type NonExistingAssetValidationAssetDocument struct {
	ID uint64 `gorm:"" json:"-" form:"-"`
}
type NonExistingAssetValidationAssetTokenInfo struct {
	ID uint64 `gorm:"" json:"-" form:"-"`
}

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
func (i IssuingWalletPublicKey) GetTokenization(gc *sharedconfig.GlobalConfig) (t TokenizedAsset) {
	e := gc.DB.Preload(clause.Associations).Order("asset_tokenization_status ASC").Where("issuing_wallet_public_key = ?", string(i)).First(&t).Error
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

func (i IssuingWalletPublicKey) GetTokenizationFeeByID(feeID uint64, gc *sharedconfig.GlobalConfig) (fee TokenizationFee) {
	gc.DB.Preload(clause.Associations).Where("id = ?", feeID).First(&fee)

	return
}

func (t *TokenizedAsset) GetTokenizationFeeByID(feeID uint64, gc *sharedconfig.GlobalConfig) (fee TokenizationFee) {
	gc.DB.Preload(clause.Associations).Where("id = ?", feeID).First(&fee)

	return
}
func (t *TokenizedAsset) GetExpressedInterestByWalletPublicKey(subscriberWalletPublicKey string, gc *sharedconfig.GlobalConfig) (exp ExpressionOfInterest, err error) {
	if t == nil {
		log.Printf("[TokenizedAsset::GetExpressedInterestByWalletPublicKey] Error tokenized asset is nil %v\n", subscriberWalletPublicKey)
		err = &errors.ErrorTemporaryServerError{}
		return
	}

	err = gc.DB.Preload(clause.Associations).Where("Tokenized_Asset_ID = ? AND Wallet_Public_Key = ?", t.ID, subscriberWalletPublicKey).First(&exp).Error

	// log.Printf("[TokenizedAsset::GetExpressedInterestByWalletPublicKey] Error getting expressedInterest for %v, %v\n", subscriberWalletPublicKey, e)

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

func (t *TokenizedAsset) UpdateTokenizedAssetFromInput(ti *TokenizedAssetJSONInput, gc *sharedconfig.GlobalConfig) TokenizedAsset {
	//TODO: set the SEC fee, Custody fee, Asset manager fee and recover feeInFiat from total asset value
	t.HasAdditionalKYCRequirements = ti.HasAdditionalKYCRequirements

	if len(ti.AdditionalKYCRequirements) > 0 && t.HasAdditionalKYCRequirements > 0 {

		t.AdditionalKYCRequirements = &ti.AdditionalKYCRequirements
	} else {
		t.AdditionalKYCRequirements = nil
		t.HasAdditionalKYCRequirements = 0
	}

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

		t.AssetWebsite = &ti.AssetWebsite
	} else {
		t.AssetWebsite = nil
	}

	if len(ti.AssetName) > 0 {

		t.AssetName = &ti.AssetName
	}

	t.ApprovedAssetCustodianID = ti.ApprovedAssetCustodianID

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

		t.AssetDescription = &ti.AssetDescription
	}

	if len(ti.AssetCountryLocation) > 0 {

		t.AssetCountryLocation = &ti.AssetCountryLocation
	}

	if len(ti.AssetPhysicalAddress) > 0 {

		t.AssetPhysicalAddress = &ti.AssetPhysicalAddress
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

		t.AssetOwnerName = &ti.AssetOwnerName
	}

	if len(ti.AssetOwnerAddress) > 0 {

		t.AssetOwnerAddress = &ti.AssetOwnerAddress
	}
	t.AssetOwnerRetainedOrContributedValue = ti.AssetOwnerRetainedOrContributedValue

	t.AssetManagerID = ti.AssetManagerID

	if len(ti.AssetQuoteCurrency) > 0 {

		t.AssetQuoteCurrency = &ti.AssetQuoteCurrency
	} else {
		t.AssetQuoteCurrency = nil
	}

	t.AssetCurrentValue = ti.AssetCurrentValue
	t.AssetMscCostOutisdeOfValuation = ti.AssetMscCostOutisdeOfValuation

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

	if len(ti.AssetCode) > 0 {

		t.AssetCode = &ti.AssetCode
	}
	t.NumberOfTokenToBeIssued = ti.NumberOfTokenToBeIssued

	var feeCompo TokenizationFee
	var feeInAsset float64

	if ti.TokenizationFeeID > 0 {
		// fee has been selected
		t.TokenizationFeeID = &ti.TokenizationFeeID
		feeCompo = t.UpdateTokenizationFeeByID(ti.TokenizationFeeID, gc)

		// Calculate Fees
		feeInAsset = decimal.NewFromFloat(t.NumberOfTokenToBeIssued * (feeCompo.FeeAssetPercentage / 100)).Truncate(7).InexactFloat64()
		t.FeeInAsset = feeInAsset
		t.FeeInFiat = decimal.NewFromFloat(t.AssetCurrentValue * (feeCompo.FeeFiatPercentage / 100)).Truncate(2).InexactFloat64()

	}
	// get SEC tokenization fee.
	var cConfig Country
	var custodyFee, assetMgtFee float64
	if t.AssetCountryLocation != nil {
		cConfig = CountryCode(*t.AssetCountryLocation).GetConfig(gc)

	}
	var secFee float64
	if cConfig.SECTokenizationFeeType == 1 {
		// fixed
		secFee = cConfig.SECTokenizationFee
	} else {
		//percent==0
		secFee = decimal.NewFromFloat(t.AssetCurrentValue * (cConfig.SECTokenizationFee / 100)).Truncate(2).InexactFloat64()
		t.SECTokenizationFeeValue = secFee
		t.SECTokenizationFeePercent = cConfig.SECTokenizationFee
	}

	if t.AssetCurrentValue > 0 {

		custodyFee = decimal.NewFromFloat(t.AssetCurrentValue * (t.CustodianFeePercent / 100)).Truncate(2).InexactFloat64()
		t.CustodianFeeValue = custodyFee

		assetMgtFee = decimal.NewFromFloat(t.AssetCurrentValue * (t.AssetManagerFeePercent / 100)).Truncate(2).InexactFloat64()
		t.AssetManagerFeeValue = assetMgtFee
	}

	//t.FeeInFiat x2 so as to recover the fee in asset
	t.ValueOfTokenizedAsset = decimal.NewFromFloat(t.AssetCurrentValue + t.AssetMscCostOutisdeOfValuation + secFee + custodyFee + assetMgtFee + (t.FeeInFiat * 2)).InexactFloat64()

	if t.NumberOfTokenToBeIssued > 0 && t.ValueOfTokenizedAsset > 0 {

		t.PricePerToken = decimal.NewFromFloat(t.ValueOfTokenizedAsset / t.NumberOfTokenToBeIssued).Truncate(7).InexactFloat64()
		// auto calculate, token to be held is less the fee. token not to be sold
		t.TotalTokenHeldByManager = decimal.NewFromFloat(t.AssetOwnerRetainedOrContributedValue / t.PricePerToken).Truncate(7).InexactFloat64()
		{
			//ensure correct the number of token to be sold.
			maxTokenToBeSold := decimal.NewFromFloat(t.NumberOfTokenToBeIssued - feeInAsset - t.TotalTokenHeldByManager).Truncate(7).InexactFloat64()
			t.MaxNumberOfTokenAvailableForSale = maxTokenToBeSold
			t.NumberOfTokenToBeSold = t.MaxNumberOfTokenAvailableForSale

		}
	}

	if len(ti.WalletToHoldAssetsNotForSale) > 0 {

		t.WalletToHoldAssetsNotForSale = &ti.WalletToHoldAssetsNotForSale
	} else {
		t.WalletToHoldAssetsNotForSale = nil
	}

	/**



	**/

	t.SalesStart = ti.SalesStart
	t.SalesEnd = ti.SalesEnd
	t.CapOnPurchase = ti.CapOnPurchase
	t.CapQuantity = ti.CapQuantity
	t.CapDurationInDays = ti.CapDurationInDays

	if len(ti.ProceedCycle) > 0 {

		t.ProceedCycle = &ti.ProceedCycle
	} else {
		t.ProceedCycle = nil
	}

	if len(ti.ProceedPayoutCurrency) > 0 {

		t.ProceedPayoutCurrency = &ti.ProceedPayoutCurrency
	} else {
		t.ProceedPayoutCurrency = nil
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

		t.LegalAdvisor = &ti.LegalAdvisor
	} else {
		t.LegalAdvisor = nil
	}

	if len(ti.FinancialAdvisor) > 0 {

		t.FinancialAdvisor = &ti.FinancialAdvisor
	} else {
		t.FinancialAdvisor = nil
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
	return *t

}

func (t *TokenizedAsset) UpdateCalculation(gc *sharedconfig.GlobalConfig) {
	//TODO: set the SEC fee, Custody fee, Asset manager fee and recover feeInFiat from total asset value

	var feeCompo TokenizationFee
	var feeInAsset float64

	if *t.TokenizationFeeID > 0 {
		// fee has been selected
		feeCompo = t.UpdateTokenizationFeeByID(*t.TokenizationFeeID, gc)

		// Calculate Fees
		feeInAsset = decimal.NewFromFloat(t.NumberOfTokenToBeIssued * (feeCompo.FeeAssetPercentage / 100)).Truncate(7).InexactFloat64()
		t.FeeInAsset = feeInAsset
		t.FeeInFiat = decimal.NewFromFloat(t.AssetCurrentValue * (feeCompo.FeeFiatPercentage / 100)).Truncate(2).InexactFloat64()

	}
	// get SEC tokenization fee.
	var cConfig Country
	var custodyFee, assetMgtFee float64
	if t.AssetCountryLocation != nil {
		cConfig = CountryCode(*t.AssetCountryLocation).GetConfig(gc)

	}
	var secFee float64
	if cConfig.SECTokenizationFeeType == 1 {
		// fixed
		secFee = cConfig.SECTokenizationFee
	} else {
		//percent==0
		secFee = decimal.NewFromFloat(t.AssetCurrentValue * (cConfig.SECTokenizationFee / 100)).Truncate(2).InexactFloat64()
		t.SECTokenizationFeeValue = secFee
		t.SECTokenizationFeePercent = cConfig.SECTokenizationFee
	}

	if t.AssetCurrentValue > 0 {

		custodyFee = decimal.NewFromFloat(t.AssetCurrentValue * (t.CustodianFeePercent / 100)).Truncate(2).InexactFloat64()
		t.CustodianFeeValue = custodyFee

		assetMgtFee = decimal.NewFromFloat(t.AssetCurrentValue * (t.AssetManagerFeePercent / 100)).Truncate(2).InexactFloat64()
		t.AssetManagerFeeValue = assetMgtFee
	}

	//t.FeeInFiat x2 so as to recover the fee in asset
	t.ValueOfTokenizedAsset = decimal.NewFromFloat(t.AssetCurrentValue + t.AssetMscCostOutisdeOfValuation + secFee + custodyFee + assetMgtFee + (t.FeeInFiat * 2)).InexactFloat64()

	if t.NumberOfTokenToBeIssued > 0 && t.ValueOfTokenizedAsset > 0 {

		t.PricePerToken = decimal.NewFromFloat(t.ValueOfTokenizedAsset / t.NumberOfTokenToBeIssued).Truncate(7).InexactFloat64()

		// auto calculate, token to be held is less the fee. token not to be sold
		t.TotalTokenHeldByManager = decimal.NewFromFloat(t.AssetOwnerRetainedOrContributedValue / t.PricePerToken).Truncate(7).InexactFloat64()

		{
			//ensure correct the number of token to be sold.
			maxTokenToBeSold := decimal.NewFromFloat(t.NumberOfTokenToBeIssued - feeInAsset - t.TotalTokenHeldByManager).Truncate(7).InexactFloat64()
			t.MaxNumberOfTokenAvailableForSale = maxTokenToBeSold
			t.NumberOfTokenToBeSold = t.MaxNumberOfTokenAvailableForSale

		}
	}

}

func (ti *TokenizedAsset) ToJSON(gc *sharedconfig.GlobalConfig) (t TokenizedAssetJSON) {

	t.ID = ti.ID
	t.CreatedAt = ti.CreatedAt
	t.UpdatedAt = ti.UpdatedAt
	t.InitiatorUsername = ti.InitiatorUsername
	t.FeeInAsset = ti.FeeInAsset
	t.FeeInFiat = ti.FeeInFiat
	t.SECTokenizationFeePercent = ti.SECTokenizationFeePercent
	t.SECTokenizationFeeValue = ti.SECTokenizationFeeValue
	t.AssetManagerFeePercent = ti.AssetManagerFeePercent
	t.AssetManagerFeeValue = ti.AssetManagerFeeValue
	t.CustodianFeePercent = ti.CustodianFeePercent
	t.CustodianFeeValue = ti.CustodianFeeValue
	t.VettingStatus = ti.VettingStatus

	t.AssetOwnerRetainedOrContributedValue = ti.AssetOwnerRetainedOrContributedValue

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
	if ti.AssetType != nil {

		t.AssetType = *ti.AssetType
	}
	if ti.AssetName != nil {
		t.AssetName = *ti.AssetName

	}
	if ti.AssetWebsite != nil {
		t.AssetWebsite = *ti.AssetWebsite

	}

	if ti.AssetCode != nil {
		t.AssetCode = *ti.AssetCode

	}
	if ti.AssetLogo != nil {
		t.AssetLogo = *ti.AssetLogo

	}

	if ti.ApprovedAssetCustodianID > 0 {
		t.ApprovedAssetCustodianID = ti.ApprovedAssetCustodianID
		t.ApprovedAssetCustodian = ti.ApprovedAssetCustodian
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
		t.AssetDescription = *ti.AssetDescription
	}
	if ti.AssetCountryLocation != nil {
		t.AssetCountryLocation = *ti.AssetCountryLocation
	}
	if ti.AssetPhysicalAddress != nil {
		t.AssetPhysicalAddress = *ti.AssetPhysicalAddress
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
		t.AssetOwnerName = *ti.AssetOwnerName
	}
	if ti.AssetOwnerAddress != nil {
		t.AssetOwnerAddress = *ti.AssetOwnerAddress
	}
	if ti.AssetManagerID > 0 {
		t.AssetManagerID = ti.AssetManagerID
		t.AssetManager = ti.AssetManager
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

		t.LegalAdvisor = *ti.LegalAdvisor
	}
	if ti.FinancialAdvisor != nil {

		t.FinancialAdvisor = *ti.FinancialAdvisor
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
	return t

}

type PaginatedTokenizedAssets struct {
	Pages        int                  `json:"pages"`
	CurrentPage  int                  `json:"currentPage"`
	TotalRecords int                  `json:"totalRecords"`
	Limit        int                  `json:"limit"`
	Records      []TokenizedAssetJSON `json:"records"`
}

type PaginatedExpressionOfInterest struct {
	Pages        int                    `json:"pages"`
	CurrentPage  int                    `json:"currentPage"`
	TotalRecords int                    `json:"totalRecords"`
	Limit        int                    `json:"limit"`
	Records      []ExpressionOfInterest `json:"records"`
}

func (e *ExpressionOfInterest) UpdateExpressionOfInterestFromInput(subscriberUsername string, subscriberWallet *UserWallet, input *ExpressionOfInterestInput, ta *TokenizedAsset, gc *sharedconfig.GlobalConfig) (ei ExpressionOfInterest) {
	if e == nil {
		return
	}

	if ta.AssetTokenizationStatus < 4 {
		return
	}
	e.TokenizedAssetID = ta.ID
	e.AssetCode = *ta.AssetCode
	e.AssetIssuer = *ta.IssuingWalletPublicKey
	e.WalletAlias = subscriberWallet.Alias
	e.WalletPublicKey = subscriberWallet.ID
	e.Amount = decimal.NewFromFloat(input.Amount).Truncate(7).InexactFloat64()
	e.Price = ta.PricePerToken
	e.SubscriberUsername = subscriberUsername
	return *e
}
