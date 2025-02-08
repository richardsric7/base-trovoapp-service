package users

// Country holds country struct
type Country struct {
	CountryCode               string  `gorm:"size:3;primaryKey" json:"countryCode"`
	SECTokenizationFeePercent float64 `gorm:"default:0" json:"SECTokenizationFeePercent"`
	SECTradeFeePercent        float64 `gorm:"default:0" json:"SECTradeFeePercent"`
	RegionName                string  `json:"regionName"`
}
