package users

import (
	"mime/multipart"
	"time"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/shopspring/decimal"
)

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
	ProceedPayoutCurrency          *string                     `json:"proceedPayoutCurrency"`
	ExemptedCountries              *string                     `json:"exemptedCountries"`
	HasAdditionalKYCRequirements   int                         `gorm:"default:0" json:"hasAdditionalKYCRequirements"`
	AdditionalKYCRequirements      *string                     `json:"additionalKYCRequirements"`
	InvestorAccreditationRequired  int                         `gorm:"default:0" json:"investorAccreditationRequired"`
	AssetTokenizationStatus        int                         `gorm:"default:0" json:"assetTokenizationStatus"`
	LastUpdatedBy                  *string                     `gorm:"null" json:"lastUpdatedBy"`
}

type TokenizedAssetJSONInput struct {
	AssetSector                    string    `json:"assetSector"`
	AssetSubSector                 string    `json:"assetSubSector"`
	AssetType                      string    `json:"assetType"`
	AssetName                      string    `json:"assetName"`
	ApprovedAssetCustodianID       uint64    `gorm:"not null" json:"approvedAssetCustodianId"`
	OfferingType                   string    `gorm:"default:'PRIVATE'" json:"offeringType"` //PRIVATE, PUBLIC
	ClosedGroupID                  string    `gorm:"null" json:"closedGroupId"`
	SecApproval                    int       `gorm:"default:0" json:"secApproval"`
	SecApprovalIdNumber            string    `json:"secApprovalIdNumber"`
	MarketMakingWallet             string    `json:"marketMakingWallet"`
	AssetDescription               string    `json:"assetDescription"`
	AssetCountryLocation           string    `json:"assetCountryLocation"`
	AssetPhysicalAddress           string    `json:"assetPhysicalAddress"`
	AssetLongitude                 string    `json:"assetLongitude"`
	AssetLatitude                  string    `json:"assetLatitude"`
	OwnershipType                  string    `json:"ownershipType"` //DIRECT,THIRD-PARTY
	OwnershipKind                  string    `json:"ownershipKind"` //INDIVIDUAL,CORPORATE,
	AssetOwnerName                 string    `json:"assetOwnerName"`
	AssetOwnerAddress              string    `json:"assetOwnerAddress"`
	AssetManagerName               string    `json:"assetManagerName"`
	AssetManagerAddress            string    `json:"assetManagerAddress"`
	AssetQuoteCurrency             string    `gorm:"default:'NGN'" json:"assetQuoteCurrency"`
	AssetCurrentValue              float64   `gorm:"default:0" json:"assetCurrentValue"`
	AssetPercentageForTokenization float64   `gorm:"default:0" json:"assetPercentageForTokenization"`
	ValueOfTokenizedAsset          float64   `gorm:"default:0" json:"valueOfTokenizedAsset"`
	ProtectionMethods              string    `json:"protectionMethods"` //csv format
	InsuranceCompanyName           string    `json:"insuranceCompanyName"`
	InsurancePolicyNumber          string    `json:"insurancePolicyNumber"`
	InsurancePolicyHolder          string    `json:"insurancePolicyHolder"`
	PercentageValueOfInsurance     float64   `gorm:"default:0" json:"percentageValueOfInsurance"`
	IsFreeFromLiensAndEncumbrances int       `gorm:"default:0" json:"IsFreeFromLiensAndEncumbrances"`
	AssetAlreadyExists             int       `gorm:"default:1" json:"assetAlreadyExists"`
	AssetCode                      string    `json:"assetCode"`
	AssetLogo                      string    `json:"assetLogo"`
	NumberOfTokenToBeIssued        float64   `json:"numberOfTokenToBeIssued"`
	NumberOfTokenToBeSold          float64   `json:"numberOfTokenToBeSold"`
	TotalTokenHeldByManager        float64   `json:"totalTokenHeldByManager"`
	WalletToHoldAssetsNotForSale   string    `json:"walletToHoldAssetsNotForSale"`
	PricePerToken                  float64   `json:"pricePerToken"`
	SalesStart                     time.Time `json:"salesStart"`
	SalesEnd                       time.Time `json:"salesEnd"`
	CapOnPurchase                  float64   `gorm:"default:0" json:"capOnPurchase"`
	CapQuantity                    float64   `gorm:"default:0" json:"capQuantity"`
	CapDurationInDays              int       `gorm:"default:0" json:"capDurationInDays"`
	ProceedCycle                   string    `gorm:"size:50" json:"proceedCycle"`
	TokenizationFeeID              uint64    `json:"tokenizationFeeId"`
	ProceedPayoutCurrency          string    `json:"proceedPayoutCurrency"`
	ExemptedCountries              string    `json:"exemptedCountries"`
	HasAdditionalKYCRequirements   int       `gorm:"default:0" json:"hasAdditionalKYCRequirements"`
	AdditionalKYCRequirements      string    `json:"additionalKYCRequirements"`
	InvestorAccreditationRequired  int       `gorm:"default:0" json:"investorAccreditationRequired"`
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
	MarketMakingWalletPublicKey    string                      `json:"marketMakingWalletPublicKey"`
	AssetDescription               string                      `json:"assetDescription"`
	AssetCountryLocation           string                      `json:"assetCountryLocation"`
	AssetPhysicalAddress           string                      `json:"assetPhysicalAddress"`
	AssetLongitude                 string                      `json:"assetLongitude"`
	AssetLatitude                  string                      `json:"assetLatitude"`
	OwnershipType                  string                      `json:"ownershipType"` //DIRECT,THIRD-PARTY
	OwnershipKind                  string                      `json:"ownershipKind"` //INDIVIDUAL,CORPORATE,
	AssetOwnerName                 string                      `json:"assetOwnerName"`
	AssetOwnerAddress              string                      `json:"assetOwnerAddress"`
	AssetManagerName               string                      `json:"assetManagerName"`
	AssetManagerAddress            string                      `json:"assetManagerAddress"`
	AssetQuoteCurrency             string                      `gorm:"default:'NGN'" json:"assetQuoteCurrency"`
	AssetCurrentValue              float64                     `gorm:"default:0" json:"assetCurrentValue"`
	AssetPercentageForTokenization float64                     `gorm:"default:0" json:"assetPercentageForTokenization"`
	ValueOfTokenizedAsset          float64                     `gorm:"default:0" json:"valueOfTokenizedAsset"`
	ProtectionMethods              string                      `json:"protectionMethods"` //csv format
	InsuranceCompanyName           string                      `json:"insuranceCompanyName"`
	InsurancePolicyNumber          string                      `json:"insurance_policy_number"`
	InsurancePolicyHolder          string                      `json:"insurancePolicyHolder"`
	PercentageValueOfInsurance     string                      `gorm:"default:0" json:"percentageValueOfInsurance"`
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
	ID               uint64 `json:"id"`
	CreatedAt        time.Time
	TokenizedAssetID string `json:"tokenizedAssetId"`
	DocumentType     uint64 `json:"documentType"`
	DocumentTitle    string `json:"documentTitle"`
	DocumentUrl      string `json:"documentUrl"`
}
type AssetTokenizationInputDocument struct {
	ID               uint64
	CreatedAt        time.Time
	TokenizedAssetID string          `form:"tokenizedAssetId"`
	DocumentType     uint64          `form:"documentType"`
	DocumentTitle    string          `form:"documentTitle"`
	DocumentFile     *multipart.File `form:"documentFile"`
	// DocumentFile     string `json:"-"`// this is not included in struct for input. already extracted by c.FormFile
}

type IssuingWalletPublicKey string

func (i IssuingWalletPublicKey) GetTokenization(gc *sharedconfig.GlobalConfig) (t TokenizedAsset) {
	gc.DB.Where("issuing_wallet = ?", string(i)).First(&t)
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

func (t *TokenizedAsset) UpdateFromInput(ti *TokenizedAssetJSONInput) {
	if len(ti.AdditionalKYCRequirements) > 0 {

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

	t.SecApproval = ti.SecApproval

	if len(ti.SecApprovalIdNumber) > 0 {

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

	if len(ti.AssetManagerName) > 0 {

		t.AssetManagerName = &ti.AssetManagerName
	}

	if len(ti.AssetManagerAddress) > 0 {

		t.AssetManagerAddress = &ti.AssetManagerAddress
	}
	if len(ti.AssetQuoteCurrency) > 0 {

		t.AssetQuoteCurrency = &ti.AssetQuoteCurrency
	}

	t.AssetCurrentValue = ti.AssetCurrentValue
	t.AssetPercentageForTokenization = ti.AssetPercentageForTokenization
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
	var fee, feeFactor float64
	{
		// Calculate Fees
		fee = decimal.NewFromFloat(t.NumberOfTokenToBeIssued * feeFactor).Truncate(7).InexactFloat64()

	}
	t.NumberOfTokenToBeIssued = ti.NumberOfTokenToBeIssued
	t.NumberOfTokenToBeSold = ti.NumberOfTokenToBeSold
	// auto calculate, token to be held is less the fee
	t.TotalTokenHeldByManager = t.NumberOfTokenToBeIssued - t.NumberOfTokenToBeSold - fee

	if len(ti.WalletToHoldAssetsNotForSale) > 0 {

		t.WalletToHoldAssetsNotForSale = &ti.WalletToHoldAssetsNotForSale
	}

	/**



	**/

	if t.NumberOfTokenToBeIssued > 0 && t.ValueOfTokenizedAsset > 0 {
		t.PricePerToken = decimal.NewFromFloat(ti.ValueOfTokenizedAsset / t.NumberOfTokenToBeIssued).Truncate(7).InexactFloat64()
	}
	t.SalesStart = ti.SalesStart
	t.SalesEnd = ti.SalesEnd
	t.CapOnPurchase = ti.CapOnPurchase
	t.CapQuantity = ti.CapQuantity
	t.CapDurationInDays = ti.CapDurationInDays

	if len(ti.ProceedCycle) > 0 {

		t.ProceedCycle = &ti.ProceedCycle
	}

	if ti.TokenizationFeeID > 0 {

		t.TokenizationFeeID = &ti.TokenizationFeeID
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

}
