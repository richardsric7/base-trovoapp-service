package users

import (
	"fmt"
	"log"
	"time"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/shopspring/decimal"
	"gorm.io/gorm/clause"
)

/*
*
*****DocumentType and codes****
ProofOfAssetExistence = 1
ProofOfAssetOwnership = 2
ProofOfAssetStatusVerification = 3
AssetCustodianAgreement = 4
ProofOfAssetManager = 5
AssetProtectionDocument = 6
AssetValuationCertificate = 7
AssetOwnerGovernmentID = 8
ProofOfAssetCondtion = 9
ThirdPartyTokenizationAgreement = 10
ThirdPartyAssetOwnerBusinessRegistration = 11
ThirdPartyAssetOwnerProofOfAddress = 12
SEC Registration/Tokenization Approval = 13
Compliance With Local Laws/regulation = 14
Compliance With Environmental Standard = 15
Environmental Impact Assessment Report = 16
Proof Of Legal/Financial Counsel = 17
Legal/Financial Advisor's Contract = 18
Proof of existing mortgages or liens n asset = 19
Proof of outstanding loans on asset = 20
Proof of legal dispute or encumbrances on asset = 21
ProofOfAssetAddress = 22
TitleDeedsOrCertificates = 23
OwnershipAgreements = 24

*
*/
type ExistingAssetValidationAssetDocument struct {
	ID                             uint64 `gorm:"" json:"-" form:"-"`
	ProofOfAssetExistence          int    `gorm:"default:1" json:"proofOfAssetExistence"`          //1
	ProofOfAssetAddress            int    `gorm:"default:1" json:"proofOfAssetAddress"`            //22
	ProofOfAssetOwnership          int    `gorm:"default:1" json:"proofOfAssetOwnership"`          //2
	TitleDeedsOrCertificates       int    `gorm:"default:1" json:"titleDeedsOrCertificates"`       //23
	OwnershipAgreements            int    `gorm:"default:1" json:"ownershipAgreements"`            //24
	ProofOfAssetStatusVerification int    `gorm:"default:1" json:"proofOfAssetStatusVerification"` //3
	ProofOfAssetCondtion           int    `gorm:"default:1" json:"proofOfAssetCondtion"`           //9
	AssetOwnerGovernmentID         int    `gorm:"default:1" json:"assetOwnerGovernmentId"`         //8

}

type TokenizedAsset struct {
	ID                             string                      `json:"id"`
	CreatedAt                      time.Time                   `json:"createdAt"`
	UpdatedAt                      time.Time                   `json:"updatedAt"`
	InitiatorUsername              string                      `gorm:"size:50;not null" json:"initiatorUsername"`
	AssetSector                    *string                     `json:"assetSector"`
	AssetSubSector                 *string                     `json:"assetSubSector"`
	AssetType                      *string                     `json:"assetType"`
	AssetName                      *string                     `json:"assetName"`
	ApprovedAssetCustodianID       uint64                      `gorm:"not null" json:"approvedAssetCustodianId"`
	ApprovedAssetCustodian         ApprovedAssetCustodian      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"approvedAssetCustodianInfo"`
	OfferingType                   *string                     `gorm:"default:'PRIVATE'" json:"offeringType"` //PRIVATE, PUBLIC
	ClosedGroupID                  *string                     `gorm:"null" json:"closedGroupId"`
	ClosedGroup                    ClosedGroup                 `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"closedGroupInfo"`
	SecApproval                    int                         `gorm:"default:0" json:"secApproval"`
	SecApprovalIdNumber            *string                     `json:"secApprovalIdNumber"`
	IssuingWalletPublicKey         string                      `gorm:"size:60" json:"issuingWalletPublicKey"`
	IssuingWalletAlias             string                      `gorm:"size:60" json:"issuingWalletAlias"`
	MarketMakingWallet             *string                     `json:"marketMakingWallet"`
	AssetDescription               *string                     `json:"assetDescription"`
	AssetCountryLocation           *string                     `json:"assetCountryLocation"`
	AssetPhysicalAddress           *string                     `json:"assetPhysicalAddress"`
	AssetLongitude                 *string                     `json:"assetLongitude"`
	AssetLatitude                  *string                     `json:"assetLatitude"`
	OwnershipType                  *string                     `json:"ownershipType"` //DIRECT,THIRD-PARTY
	OwnershipKind                  *string                     `json:"ownershipKind"` //INDIVIDUAL,CORPORATE,
	AssetOwnerName                 *string                     `json:"assetOwnerName"`
	AssetOwnerAddress              *string                     `json:"assetOwnerAddress"`
	AssetManagerID                 uint64                      `json:"assetManagerId"`
	AssetManager                   AssetManager                `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"assetManagerInfo"`
	AssetQuoteCurrency             *string                     `gorm:"default:'NGN'" json:"assetQuoteCurrency"`
	AssetCurrentValue              float64                     `gorm:"default:0" json:"assetCurrentValue"`
	AssetMscCostOutisdeOfValuation float64                     `gorm:"default:0" json:"assetMscCostOutisdeOfValuation"`
	ValueOfTokenizedAsset          float64                     `gorm:"default:0" json:"valueOfTokenizedAsset"`
	ProtectionMethods              *string                     `json:"protectionMethods"` //csv format
	InsuranceCompanyName           *string                     `json:"insuranceCompanyName"`
	InsurancePolicyNumber          *string                     `json:"insurancePolicyNumber"`
	InsurancePolicyHolder          *string                     `json:"insurancePolicyHolder"`
	PercentageValueOfInsurance     float64                     `gorm:"default:0" json:"percentageValueOfInsurance"`
	IsFreeFromLiensAndEncumbrances int                         `gorm:"default:0" json:"IsFreeFromLiensAndEncumbrances"`
	AssetAlreadyExists             int                         `gorm:"default:1" json:"assetAlreadyExists"`
	AssetTokenizationDocuments     []AssetTokenizationDocument `json:"AssetTokenizationDocuments"`
	AssetCode                      *string                     `json:"assetCode"`
	AssetLogo                      *string                     `json:"assetLogo"`
	NumberOfTokenToBeIssued        float64                     `gorm:"default:0" json:"numberOfTokenToBeIssued"`
	NumberOfTokenToBeSold          float64                     `gorm:"default:0" json:"numberOfTokenToBeSold"`
	TotalTokenHeldByManager        float64                     `gorm:"default:0" json:"totalTokenHeldByManager"`
	WalletToHoldAssetsNotForSale   *string                     `json:"walletToHoldAssetsNotForSale"`
	PricePerToken                  float64                     `gorm:"default:0" json:"pricePerToken"`
	SalesStart                     time.Time                   `json:"salesStart"`
	SalesEnd                       time.Time                   `json:"salesEnd"`
	CapOnPurchase                  float64                     `gorm:"default:0" json:"capOnPurchase"`
	CapQuantity                    float64                     `gorm:"default:0" json:"capQuantity"`
	CapDurationInDays              int                         `gorm:"default:0" json:"capDurationInDays"`
	ProceedCycle                   *string                     `gorm:"size:50" json:"proceedCycle"`
	TokenizationFeeID              *uint64                     `gorm:"default:0" json:"tokenizationFeeId"`
	TokenizationFee                TokenizationFee             `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"tokenizationFee"`
	ProceedPayoutCurrency          *string                     `json:"proceedPayoutCurrency"`
	ExemptedCountries              *string                     `json:"exemptedCountries"`
	HasAdditionalKYCRequirements   int                         `gorm:"default:0" json:"hasAdditionalKYCRequirements"`
	AdditionalKYCRequirements      *string                     `json:"additionalKYCRequirements"`
	InvestorAccreditationRequired  int                         `gorm:"default:0" json:"investorAccreditationRequired"`
	AssetTokenizationStatus        int                         `gorm:"default:0" json:"assetTokenizationStatus"`
	LastUpdatedBy                  *string                     `gorm:"null" json:"lastUpdatedBy"`
}

type TokenizedAssetJSONInput struct {
	AssetSector                    string       `json:"assetSector"`
	AssetSubSector                 string       `json:"assetSubSector"`
	AssetType                      string       `json:"assetType"`
	AssetName                      string       `json:"assetName"`
	ApprovedAssetCustodianID       uint64       `gorm:"not null" json:"approvedAssetCustodianId"`
	OfferingType                   string       `gorm:"default:'PRIVATE'" json:"offeringType"` //PRIVATE, PUBLIC
	ClosedGroupID                  string       `gorm:"null" json:"closedGroupId"`
	SecApproval                    int          `gorm:"default:0" json:"secApproval"`
	SecApprovalIdNumber            string       `json:"secApprovalIdNumber"`
	MarketMakingWallet             string       `json:"marketMakingWallet"`
	AssetDescription               string       `json:"assetDescription"`
	AssetCountryLocation           string       `json:"assetCountryLocation"`
	AssetPhysicalAddress           string       `json:"assetPhysicalAddress"`
	AssetLongitude                 string       `json:"assetLongitude"`
	AssetLatitude                  string       `json:"assetLatitude"`
	OwnershipType                  string       `json:"ownershipType"` //DIRECT,THIRD-PARTY
	OwnershipKind                  string       `json:"ownershipKind"` //INDIVIDUAL,CORPORATE,
	AssetOwnerName                 string       `json:"assetOwnerName"`
	AssetOwnerAddress              string       `json:"assetOwnerAddress"`
	AssetManagerID                 uint64       `json:"assetManagerId"`
	AssetManager                   AssetManager `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"assetManagerInfo"`
	AssetQuoteCurrency             string       `gorm:"default:'NGN'" json:"assetQuoteCurrency"`
	AssetCurrentValue              float64      `gorm:"default:0" json:"assetCurrentValue"`
	AssetMscCostOutisdeOfValuation float64      `gorm:"default:0" json:"assetMscCostOutisdeOfValuation"`
	ValueOfTokenizedAsset          float64      `gorm:"default:0" json:"valueOfTokenizedAsset"`
	ProtectionMethods              string       `json:"protectionMethods"` //csv format
	InsuranceCompanyName           string       `json:"insuranceCompanyName"`
	InsurancePolicyNumber          string       `json:"insurancePolicyNumber"`
	InsurancePolicyHolder          string       `json:"insurancePolicyHolder"`
	PercentageValueOfInsurance     float64      `gorm:"default:0" json:"percentageValueOfInsurance"`
	IsFreeFromLiensAndEncumbrances int          `gorm:"default:0" json:"IsFreeFromLiensAndEncumbrances"`
	AssetAlreadyExists             int          `gorm:"default:1" json:"assetAlreadyExists"`
	AssetCode                      string       `json:"assetCode"`
	AssetLogo                      string       `json:"assetLogo"`
	NumberOfTokenToBeIssued        float64      `json:"numberOfTokenToBeIssued"`
	NumberOfTokenToBeSold          float64      `json:"numberOfTokenToBeSold"`
	TotalTokenHeldByManager        float64      `json:"totalTokenHeldByManager"`
	WalletToHoldAssetsNotForSale   string       `json:"walletToHoldAssetsNotForSale"`
	PricePerToken                  float64      `json:"pricePerToken"`
	SalesStart                     time.Time    `json:"salesStart"`
	SalesEnd                       time.Time    `json:"salesEnd"`
	CapOnPurchase                  float64      `gorm:"default:0" json:"capOnPurchase"`
	CapQuantity                    float64      `gorm:"default:0" json:"capQuantity"`
	CapDurationInDays              int          `gorm:"default:0" json:"capDurationInDays"`
	ProceedCycle                   string       `gorm:"size:50" json:"proceedCycle"`
	TokenizationFeeID              uint64       `json:"tokenizationFeeId"`
	ProceedPayoutCurrency          string       `json:"proceedPayoutCurrency"`
	ExemptedCountries              string       `json:"exemptedCountries"`
	HasAdditionalKYCRequirements   int          `gorm:"default:0" json:"hasAdditionalKYCRequirements"`
	AdditionalKYCRequirements      string       `json:"additionalKYCRequirements"`
	InvestorAccreditationRequired  int          `gorm:"default:0" json:"investorAccreditationRequired"`
	Messages                       []string     `json:"messages"`
}
type TokenizedAssetJSON struct {
	ID                             string                      `json:"id"`
	CreatedAt                      time.Time                   `json:"createdAt"`
	UpdatedAt                      time.Time                   `json:"updatedAt"`
	InitiatorUsername              string                      `gorm:"size:50;not null" json:"initiatorUsername"`
	AssetSector                    string                      `json:"assetSector"`
	AssetSubSector                 string                      `json:"assetSubSector"`
	AssetType                      string                      `json:"assetType"`
	AssetName                      string                      `json:"assetName"`
	ApprovedAssetCustodianID       uint64                      `gorm:"not null" json:"approvedAssetCustodianId"`
	ApprovedAssetCustodian         ApprovedAssetCustodian      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"approvedAssetCustodianInfo"`
	OfferingType                   string                      `gorm:"default:'PRIVATE'" json:"offeringType"` //PRIVATE, PUBLIC
	ClosedGroupID                  string                      `gorm:"null" json:"closedGroupId"`
	ClosedGroup                    ClosedGroup                 `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"closedGroupInfo"`
	SecApproval                    int                         `gorm:"default:0" json:"secApproval"`
	SecApprovalIdNumber            string                      `json:"secApprovalIdNumber"`
	IssuingWalletPublicKey         string                      `gorm:"size:60" json:"issuingWalletPublicKey"`
	IssuingWalletAlias             string                      `gorm:"size:60" json:"issuingWalletAlias"`
	MarketMakingWallet             string                      `json:"marketMakingWallet"`
	AssetDescription               string                      `json:"assetDescription"`
	AssetCountryLocation           string                      `json:"assetCountryLocation"`
	AssetPhysicalAddress           string                      `json:"assetPhysicalAddress"`
	AssetLongitude                 string                      `json:"assetLongitude"`
	AssetLatitude                  string                      `json:"assetLatitude"`
	OwnershipType                  string                      `json:"ownershipType"` //DIRECT,THIRD-PARTY
	OwnershipKind                  string                      `json:"ownershipKind"` //INDIVIDUAL,CORPORATE,
	AssetOwnerName                 string                      `json:"assetOwnerName"`
	AssetOwnerAddress              string                      `json:"assetOwnerAddress"`
	AssetManagerID                 uint64                      `json:"assetManagerId"`
	AssetManager                   AssetManager                `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"assetManagerInfo"`
	AssetQuoteCurrency             string                      `gorm:"default:'NGN'" json:"assetQuoteCurrency"`
	AssetCurrentValue              float64                     `gorm:"default:0" json:"assetCurrentValue"`
	AssetMscCostOutisdeOfValuation float64                     `gorm:"default:0" json:"assetMscCostOutisdeOfValuation"`
	ValueOfTokenizedAsset          float64                     `gorm:"default:0" json:"valueOfTokenizedAsset"`
	ProtectionMethods              string                      `json:"protectionMethods"` //csv format
	InsuranceCompanyName           string                      `json:"insuranceCompanyName"`
	InsurancePolicyNumber          string                      `json:"insurancePolicyNumber"`
	InsurancePolicyHolder          string                      `json:"insurancePolicyHolder"`
	PercentageValueOfInsurance     float64                     `gorm:"default:0" json:"percentageValueOfInsurance"`
	IsFreeFromLiensAndEncumbrances int                         `gorm:"default:0" json:"IsFreeFromLiensAndEncumbrances"`
	AssetAlreadyExists             int                         `gorm:"default:1" json:"assetAlreadyExists"`
	AssetTokenizationDocuments     []AssetTokenizationDocument `json:"AssetTokenizationDocuments"`
	AssetCode                      string                      `json:"assetCode"`
	AssetLogo                      string                      `json:"assetLogo"`
	NumberOfTokenToBeIssued        float64                     `json:"numberOfTokenToBeIssued"`
	NumberOfTokenToBeSold          float64                     `json:"numberOfTokenToBeSold"`
	TotalTokenHeldByManager        float64                     `json:"totalTokenHeldByManager"`
	WalletToHoldAssetsNotForSale   string                      `json:"walletToHoldAssetsNotForSale"`
	PricePerToken                  float64                     `json:"pricePerToken"`
	SalesStart                     time.Time                   `json:"salesStart"`
	SalesEnd                       time.Time                   `json:"salesEnd"`
	CapOnPurchase                  float64                     `gorm:"default:0" json:"capOnPurchase"`
	CapQuantity                    float64                     `gorm:"default:0" json:"capQuantity"`
	CapDurationInDays              int                         `gorm:"default:0" json:"capDurationInDays"`
	ProceedCycle                   string                      `gorm:"size:50" json:"proceedCycle"`
	TokenizationFeeID              uint64                      `json:"tokenizationFeeId"`
	TokenizationFee                TokenizationFee             `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"tokenizationFee"`
	ProceedPayoutCurrency          string                      `json:"proceedPayoutCurrency"`
	ExemptedCountries              string                      `json:"exemptedCountries"`
	HasAdditionalKYCRequirements   int                         `gorm:"default:0" json:"hasAdditionalKYCRequirements"`
	AdditionalKYCRequirements      string                      `json:"additionalKYCRequirements"`
	InvestorAccreditationRequired  int                         `gorm:"default:0" json:"investorAccreditationRequired"`
	AssetTokenizationStatus        int                         `gorm:"default:0" json:"assetTokenizationStatus"`
	LastUpdatedBy                  string                      `gorm:"null" json:"lastUpdatedBy"`
}

type TokenizedAssetSector struct {
	ID string `gorm:"primaryKey;size:100" json:"sector"`
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

type TokenizationFee struct {
	ID                 uint64  `gorm:"" json:"id"`
	FeeFiatPercentage  float64 `json:"feeFiatPercentage"`
	FeeFiatCap         float64 `json:"feeFiatCap"`
	FeeAssetPercentage float64 `json:"feeAssetPercentage"`
	// FeeAssetCap        float64 `json:"feeAssetCap"`
	FeeDescription string `json:"feeDescription"`
	Inactive       int    `gorm:"default:0" json:"-"`
}

type TokenizedAssetType struct {
	ID                        uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	TokenizedAssetSubSectorID string `gorm:"size:100;index:idx_asset_type_unique,unique" json:"assetSubSectorId"`
	AssetType                 string `gorm:"size:100;index:idx_asset_type_unique,unique" json:"assetType"`
}

type ApprovedAssetCustodian struct {
	ID                    uint64 `gorm:"" json:"id"`
	AssetCustodianName    string `gorm:"size:100" json:"assetCustodianName"`
	AssetCustodianAddress string `json:"assetCustodianAddress"`
	AssetCustodianCountry string `gorm:"size:3" json:"assetCustodianCountry"`
}
type AssetManager struct {
	ID                  uint64 `gorm:"" json:"id"`
	AssetManagerName    string `gorm:"size:100" json:"assetManagerName"`
	AssetManagerAddress string `json:"assetManagerAddress"`
	AssetManagerCountry string `gorm:"size:3" json:"assetManagerCountry"`
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
type NonExistingAssetValidationAssetDocument struct {
	ID uint64 `gorm:"" json:"-" form:"-"`
}
type NonExistingAssetValidationAssetTokenInfo struct {
	ID uint64 `gorm:"" json:"-" form:"-"`
}

func (i IssuingWalletPublicKey) GetTokenization(gc *sharedconfig.GlobalConfig) (t TokenizedAsset) {
	e := gc.DB.Preload(clause.Associations).Where("issuing_wallet_public_key = ?", string(i)).First(&t).Error
	if e != nil {
		log.Printf("[IssuingWalletPublicKey::GetTokenization] Error getting tokenized asset for %v, %v\n", string(i), e)
	}
	return
}

func (i IssuingWalletPublicKey) GetTokenizationFeeByID(feeID uint64, gc *sharedconfig.GlobalConfig) (fee TokenizationFee) {
	gc.DB.Preload(clause.Associations).Where("id = ?", feeID).First(&fee)

	return
}

/**
*****DocumentType and codes****
ProofOfAssetExistence = 1
ProofOfAssetOwnership = 2
ProofOfAssetStatusVerification = 3
AssetCustodianAgreement = 4
ProofOfAssetManager = 5
AssetProtectionDocument = 6
AssetValuationCertificate = 7
AssetOwnerGovernmentID = 8
ProofOfAssetCondtion = 9
ThirdPartyTokenizationAgreement = 10
ThirdPartyAssetOwnerBusinessRegistration = 11
ThirdPartyAssetOwnerProofOfAddress = 12
SEC Registration/Tokenization Approval = 13
Compliance With Local Laws/regulation = 14
Compliance With Environmental Standard = 15
Environmental Impact Assessment Report = 16
Proof Of Legal/Financial Counsel = 17
Legal/Financial Advisor's Contract = 18
Proof of existing mortgages or liens n asset = 19
Proof of outstanding loans on asset = 20
Proof of legal dispute or encumbrances on asset = 21

**/

func (t *TokenizedAsset) GetTokenizationFeeByID(feeID uint64, gc *sharedconfig.GlobalConfig) (fee TokenizationFee) {
	gc.DB.Preload(clause.Associations).Where("id = ?", feeID).First(&fee)

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

func (t *TokenizedAsset) UpdateFromInput(ti *TokenizedAssetJSONInput, gc *sharedconfig.GlobalConfig) TokenizedAsset {
	if ti.HasAdditionalKYCRequirements > 0 && len(ti.AdditionalKYCRequirements) > 0 {

		t.AdditionalKYCRequirements = &ti.AdditionalKYCRequirements

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

	if len(ti.AssetName) > 0 {

		t.AssetName = &ti.AssetName
	}

	t.ApprovedAssetCustodianID = ti.ApprovedAssetCustodianID

	if len(ti.OfferingType) > 0 {

		t.OfferingType = &ti.OfferingType
	}

	if len(ti.ClosedGroupID) > 0 {

		t.ClosedGroupID = &ti.ClosedGroupID
	}

	if len(ti.SecApprovalIdNumber) > 0 && ti.SecApproval > 0 {

		t.SecApproval = ti.SecApproval

		t.SecApprovalIdNumber = &ti.SecApprovalIdNumber
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

	// if len(ti.AssetManagerName) > 0 {

	t.AssetManagerID = ti.AssetManagerID
	// }

	// if len(ti.AssetManagerAddress) > 0 {

	// 	t.AssetManagerAddress = &ti.AssetManagerAddress
	// }
	if len(ti.AssetQuoteCurrency) > 0 {

		t.AssetQuoteCurrency = &ti.AssetQuoteCurrency
	}

	t.AssetCurrentValue = ti.AssetCurrentValue
	t.AssetMscCostOutisdeOfValuation = ti.AssetMscCostOutisdeOfValuation
	t.ValueOfTokenizedAsset = ti.ValueOfTokenizedAsset

	if len(ti.ProtectionMethods) > 0 {

		t.ProtectionMethods = &ti.ProtectionMethods
	}

	if len(ti.InsuranceCompanyName) > 0 {

		t.InsuranceCompanyName = &ti.InsuranceCompanyName
	}

	if len(ti.InsurancePolicyHolder) > 0 {

		t.InsurancePolicyHolder = &ti.InsurancePolicyHolder
	}

	if len(ti.InsurancePolicyNumber) > 0 {

		t.InsurancePolicyNumber = &ti.InsurancePolicyNumber
	}

	t.PercentageValueOfInsurance = ti.PercentageValueOfInsurance
	t.IsFreeFromLiensAndEncumbrances = ti.IsFreeFromLiensAndEncumbrances
	t.AssetAlreadyExists = ti.AssetAlreadyExists

	if len(ti.AssetCode) > 0 {

		t.AssetCode = &ti.AssetCode
	}

	if len(ti.AssetLogo) > 0 {

		t.AssetLogo = &ti.AssetLogo
	}
	if t.NumberOfTokenToBeIssued > 0 && t.ValueOfTokenizedAsset > 0 {
		totalValuation := (ti.ValueOfTokenizedAsset + ti.AssetMscCostOutisdeOfValuation)
		t.PricePerToken = decimal.NewFromFloat(totalValuation / t.NumberOfTokenToBeIssued).Truncate(7).InexactFloat64()
	}
	var feeCompo TokenizationFee
	var assetFee float64

	if ti.TokenizationFeeID > 0 {
		// fee has been selected
		t.TokenizationFeeID = &ti.TokenizationFeeID
		t.UpdateTokenizationFeeByID(ti.TokenizationFeeID, gc)

		{
			// Calculate Fees
			assetFee = decimal.NewFromFloat(t.NumberOfTokenToBeIssued * (feeCompo.FeeAssetPercentage / 100)).Truncate(7).InexactFloat64()

		}

	}

	t.NumberOfTokenToBeIssued = ti.NumberOfTokenToBeIssued
	t.NumberOfTokenToBeSold = ti.NumberOfTokenToBeSold
	// auto calculate, token to be held is less the fee

	{
		//ensure correct the number of token to be sold.
		maxTokenToBeSold := t.NumberOfTokenToBeIssued - assetFee

		if t.NumberOfTokenToBeSold > maxTokenToBeSold {
			t.NumberOfTokenToBeSold = maxTokenToBeSold
			ti.NumberOfTokenToBeSold = maxTokenToBeSold
			ti.Messages = append(ti.Messages, fmt.Sprintf("Submitted Number of tokens to be sold has been adjusted to %v to account for deduction of the asset fee of %v. You may wish to adjust your fee option and then adjust the amount to be sold manually again.", maxTokenToBeSold, assetFee))

		}
	}
	t.TotalTokenHeldByManager = t.NumberOfTokenToBeIssued - t.NumberOfTokenToBeSold - assetFee
	if len(ti.WalletToHoldAssetsNotForSale) > 0 {

		t.WalletToHoldAssetsNotForSale = &ti.WalletToHoldAssetsNotForSale
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
	}

	if len(ti.ProceedPayoutCurrency) > 0 {

		t.ProceedPayoutCurrency = &ti.ProceedPayoutCurrency
	}

	if len(ti.ExemptedCountries) > 0 {

		t.ExemptedCountries = &ti.ExemptedCountries
	}

	t.HasAdditionalKYCRequirements = ti.HasAdditionalKYCRequirements

	if len(ti.AdditionalKYCRequirements) > 0 && t.HasAdditionalKYCRequirements > 0 {

		t.AdditionalKYCRequirements = &ti.AdditionalKYCRequirements
	}

	t.InvestorAccreditationRequired = ti.InvestorAccreditationRequired
	return *t

}

func (ti *TokenizedAsset) ToJSON(gc *sharedconfig.GlobalConfig) (t TokenizedAssetJSON) {

	t.ID = ti.ID
	t.CreatedAt = ti.CreatedAt
	t.UpdatedAt = ti.UpdatedAt
	t.InitiatorUsername = ti.InitiatorUsername

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
	t.IssuingWalletPublicKey = ti.IssuingWalletPublicKey
	t.IssuingWalletAlias = ti.IssuingWalletAlias

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

	return t

}

type PaginatedTokenizedAssets struct {
	Pages        int                  `json:"pages"`
	CurrentPage  int                  `json:"currentPage"`
	TotalRecords int                  `json:"totalRecords"`
	Limit        int                  `json:"limit"`
	Records      []TokenizedAssetJSON `json:"records"`
}
