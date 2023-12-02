package users

import "time"

type TokenizedAsset struct {
	ID                             string                      `json:"id"`
	CreatedAt                      time.Time                   `json:"createdAt"`
	UpdatedAt                      time.Time                   `json:"updatedAt"`
	OwnerUsername                  string                      `gorm:"not null" json:"ownerUsername"`
	AssetSector                    *string                     `json:"assetSector"`
	AssetSubSector                 *string                     `json:"assetSubSector"`
	AssetType                      *string                     `json:"assetType"`
	AssetName                      *string                     `json:"assetName"`
	ApprovedAssetCustodianID       uint64                      `gorm:"not null" json:"approvedAssetCustodianID"`
	ApprovedAssetCustodian         ApprovedAssetCustodian      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"approvedAssetCustodianInfo"`
	OfferingType                   *string                     `gorm:"default:'PRIVATE'" json:"offeringType"` //PRIVATE, PUBLIC
	ClosedGroupID                  *string                     `gorm:"" json:"closedGroupId"`
	ClosedGroup                    ClosedGroup                 `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"closedGroupInfo"`
	SecApproval                    int                         `gorm:"default:0" json:"secApproval"`
	SecApprovalIdNumber            *string                     `json:"secApprovalIdNumber"`
	MintingWallet                  *string                     `gorm:"size:60" json:"mintingWallet"`
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
	AssetManagerName               *string                     `json:"assetManagerName"`
	AssetManagerAddress            *string                     `json:"assetManagerAddress"`
	AssetQuoteCurrency             *string                     `gorm:"default:'NGN'" json:"assetQuoteCurrency"`
	AssetCurrentValue              float64                     `gorm:"default:0" json:"assetCurrentValue"`
	AssetPercentageForTokenization float64                     `gorm:"default:0" json:"assetPercentageForTokenization"`
	ValueOfTokenizedAsset          float64                     `gorm:"default:0" json:"valueOfTokenizedAsset"`
	ProtectionMethods              *string                     `json:"protectionMethods"` //csv format
	InsuranceCompanyName           *string                     `json:"insuranceCompanyName"`
	InsurancePolicyNumber          *string                     `json:"insurance_policy_number"`
	InsurancePolicyHolder          *string                     `json:"insurancePolicyHolder"`
	PercentageValueOfInsurance     string                      `gorm:"default:0" json:"percentageValueOfInsurance"`
	IsFreeFromLiensAndEncumbrances int                         `gorm:"default:0" json:"IsFreeFromLiensAndEncumbrances"`
	AssetAlreadyExists             int                         `gorm:"default:1" json:"assetAlreadyExists"`
	AssetTokenizationDocuments     []AssetTokenizationDocument `json:"AssetTokenizationDocuments"`
	AssetCode                      *string                     `json:"assetCode"`
	AssetLogo                      *string                     `json:"assetLogo"`
	NumberOfTokenToBeIssued        float64                     `json:"numberOfTokenToBeIssued"`
	NumberOfTokenToBeSold          float64                     `json:"numberOfTokenToBeSold"`
	TotalTokenHeldByManager        float64                     `json:"totalTokenHeldByManager"`
	WalletToHoldAssetsNotForSale   *string                     `json:"walletToHoldAssetsNotForSale"`
	PricePerToken                  float64                     `json:"pricePerToken"`
	SalesStart                     time.Time                   `json:"salesStart"`
	SalesEnd                       time.Time                   `json:"salesEnd"`
	CapOnPurchase                  float64                     `gorm:"default:0" json:"capOnPurchase"`
	CapQuantity                    float64                     `gorm:"default:0" json:"capQuantity"`
	CapDurationInDays              int                         `gorm:"default:0" json:"capDurationInDays"`
	ProceedCycle                   *string                     `gorm:"size:50" json:"proceedCycle"`
	TokenizationFeeID              uint64                      `json:"tokenizationFeeId"`
	ProceedPayoutCurrency          *string                     `json:"proceedPayoutCurrency"`
	ExemptedCountries              *string                     `json:"exemptedCountries"`
	HasAdditionalKYCRequirements   int                         `gorm:"default:0" json:"hasAdditionalKYCRequirements"`
	AdditionalKYCRequirements      *string                     `json:"additionalKYCRequirements"`
	InvestorAccreditationRequired  int                         `gorm:"default:0" json:"investorAccreditationRequired"`
	AssetTokenizationStatus        int                         `gorm:"default:0" json:"assetTokenizationStatus"`
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
}

type TokenizationFee struct {
	ID                 uint64  `gorm:"" json:"id"`
	FeeFiatPercentage  float64 `json:"feeFiatPercentage"`
	FeeFiatCap         float64 `json:"feeFiatCap"`
	FeeAssetPercentage float64 `json:"feeAssetPercentage"`
	FeeAssetCap        float64 `json:"feeAssetCap"`
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

type AssetProtectionOption struct {
	ID string `gorm:"size:100;primaryKey" json:"id"`
}

type ProceedCycle struct {
	ID string `gorm:"size:50" json:"id"`
}

type AssetTokenizationDocument struct {
	ID               uint64
	CreatedAt        time.Time
	TokenizedAssetID string `json:"tokenizedAssetID"`
	DocumentType     uint64 `json:"documentType"`
	DocumentTitle    string `json:"documentTitle"`
	DocumentUrl      string `json:"documentUrl"`
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
Proof of legal dispute or encumberances on asset = 21

**/
