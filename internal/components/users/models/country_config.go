package users

import "trovo-wallet-api/internal/sharedconfig"

// Country holds country struct
type Country struct {
	CountryCode                              string  `gorm:"size:2;primaryKey" json:"countryCode"`
	SECTokenizationFeePercent                float64 `gorm:"default:0" json:"SECTokenizationFeePercent"`
	SECTokenizationFeeFixed                  float64 `gorm:"default:0" json:"SECTokenizationFeeFixed"`
	SECTradeFeePercent                       float64 `gorm:"default:0" json:"SECTradeFeePercent"`
	SECTradeFeeFixed                         float64 `gorm:"default:0" json:"SECTradeFeeFixed"`
	RegulatorName                            string  `gorm:"default:'SECURITY AND EXCHANGE COMMISSION'" json:"regulatorName"`
	RegionName                               string  `json:"regionName"`
	CountryName                              string  `json:"countryName"`
	QuoteCurrencyCode                        string  `json:"quoteCurrencyCode"`
	FiatLabel                                string  `json:"fiatLabel"`
	FiatGlyph                                string  `json:"fiatGlyph"`
	MinTokenizationFee                       float64 `gorm:"default:0" json:"minTokenizationFee"`
	MinTROVBalanceForTokenizationApplication float64 `gorm:"default:600" json:"minTROVBalanceForTokenizationApplication"`
	TokenizationApplicationFee               float64 `gorm:"default:0" json:"tokenizationApplicationFee"`
	TokenizationApplicationFeeAsset          string  `gorm:"default:'TROV:GAXMBPVA2GNG6A3NV6Q664VZASMROS5ZACKSMTPVCRIKPOJIV43A2CTJ'" json:"tokenizationApplicationFeeAsset"`
	VATPercent                               float64 `gorm:"default:0" json:"vatPercent"`
}

type CountryCode string

func (c CountryCode) GetConfig(gc *sharedconfig.GlobalConfig) (cConfig Country) {
	gc.DB.Where("country_code = ?", string(c)).First(&cConfig)
	return
}
func (c CountryCode) GetCustodyFee(custodianID uint64, gc *sharedconfig.GlobalConfig) (cConfig ApprovedAssetCustodian) {
	gc.DB.Where("id = ? AND Asset_Custodian_Country = ?", custodianID, string(c)).First(&cConfig)
	return
}
func (c CountryCode) GetIssuingHouseFee(issuingHouseID uint64, gc *sharedconfig.GlobalConfig) (issuingHouse AssetIssuingHouse) {
	gc.DB.Where("id = ? AND Asset_Issuing_House_Country = ?", issuingHouseID, string(c)).First(&issuingHouse)
	return
}
func (c CountryCode) GetAssetMgtFee(assetManagerID uint64, gc *sharedconfig.GlobalConfig) (cConfig AssetManager) {
	gc.DB.Where("id = ? AND Asset_Manager_Country = ?", assetManagerID, string(c)).First(&cConfig)
	return
}
