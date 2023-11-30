package users

import "time"

type TokenizedAsset struct {
	ID                             string    `json:"id"`
	CreatedAt                      time.Time `json:"createdAt"`
	UpdatedAt                      time.Time `json:"updatedAt"`
	UserID                         string    `json:"userID"`
	MintingWallet                  string    `json:"mintingWallet"`
	MarketMakingWallet             string    `json:"marketMakingWallet"`
	AssetCategory                  string    `json:"assetCategory"`
	AssetName                      string    `json:"assetName"`
	AssetDescription               string    `json:"assetDescription"`
	AssetCountryLocation           string    `json:"assetCountryLocation"`
	AssetPhysicalAddress           string    `json:"assetPhysicalAddress"`
	AssetLongitude                 string    `json:"assetLongitude"`
	AssetLatitude                  string    `json:"assetLatitude"`
	OwnershipType                  string    `json:"ownershipType"`
	AssetCustodianName             string    `json:"assetCustodianName"`
	AssetCustodianAddress          string    `json:"assetCustodianAddress"`
	AssetManagerName               string    `json:"assetManagerName"`
	AssetManagerAddress            string    `json:"assetManagerAddress"`
	AssetOwnerName                 string    `json:"assetOwnerName"`
	AssetOwnerAddress              string    `json:"assetOwnerAddress"`
	OrganizationName               string    `json:"OrganizationName"`
	OrganizationAddress            string    `json:"OrganizationAddress"`
	AssetValue                     string    `json:"assetValue"`
	AssetPercentageForTokenization string    `json:"assetPercentageForTokenization"`
	ValueOfTokenizedAsset          string    `json:"valueOfTokenizedAsset"`
	ProtectionMethods              string    `json:"protectionMethods"`
	InsuranceCompanyName           string    `json:"insuranceCompanyName"`
	InsurancePolicyNumber          string    `json:"insurance_policy_number"`
	InsurancePolicyHolder          string    `json:"insurancePolicyHolder"`
	PercentageValueOfInsurance     string    `json:"percentageValueOfInsurance"`
	IsFreeFromLiensAndEncumbrances int       `json:"IsFreeFromLiensAndEncumbrances"`
	ProofOfAssetsExistence         string    `json:"proofOfAssetsExistence"`
	ProofOfOwnership               string    `json:"proofOfOwnership"`
	AssetStatusVerification        string    `json:"assetStatusVerification"`
	AssetCustodianAgreement        string    `json:"assetCustodianAgreement"`
	ProofOfAssetManager            string    `json:"proofOfAssetManager"`
	AssetProtectionDocument        string    `json:"assetProtectionDocument"`
	AssetValuationCertificate      string    `json:"assetValuationCertificate"`
	AssetCode                      string    `json:"assetCode"`
	AssetLogo                      string    `json:"assetLogo"`
	NumberOfTokenToBeIssued        int       `json:"numberOfTokenToBeIssued"`
	NumberOfTokenToBeSold          int       `json:"numberOfTokenToBeSold"`
	TotalTokenHeldByManager        int       `json:"totalTokenHeldByManager"`
	WalletToHoldAssetsNotForSale   string    `json:"walletToHoldAssetsNotForSale"`
	PricePerToken                  string    `json:"pricePerToken"`
	FundingMethod                  string    `json:"fundingMethod"`
	SalesStart                     time.Time `json:"salesStart"`
	SalesEnd                       time.Time `json:"salesEnd"`
	CapOnPurchase                  bool      `json:"capOnPurchase"`
	CapQuantity                    int       `json:"capQuantity"`
	CapDurationInDays              int       `json:"capDurationInDays"`
	ProceedCycle                   string    `json:"proceedCycle"`
	ProceedPayoutCurrency          string    `json:"proceedPayoutCurrency"`
	ExemptedCountries              string    `json:"exemptedCountries"`
	HasAdditionalKYCRequirements   int       `json:"hasAdditionalKYCRequirements"`
	AdditionalKYCRequirements      string    `json:"additionalKYCRequirements"`
	InvestorAccreditationRequired  int       `json:"investorAccreditationRequired"`
	AssetStatus                    string    `json:"assetStatus"`
}

type TokenizedAssetSector struct {
	ID string `gorm:"primaryKey;size:100" json:"sector"`
}

type TokenizedAssetSubSector struct {
	ID                     string `gorm:"primaryKey;size:100" json:"subSector"`
	TokenizedAssetSectorID string `gorm:"size:100" json:"assetSectorId"`
}

type TokenizedAssetType struct {
	ID                     string `gorm:"primaryKey" json:"id"`
	TokenizedAssetSectorID string `gorm:"size:100;index:idx_asset_type_unique,unique" json:"assetSectorId"`
	AssetType              string `gorm:"size:100;index:idx_asset_type_unique,unique" json:"assetType"`
}
