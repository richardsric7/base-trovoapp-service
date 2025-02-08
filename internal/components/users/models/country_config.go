package users

// Country holds country struct
type Country struct {
	CountryCode   string  `gorm:"size:3;primaryKey" json:"countryCode"`
	SECFeePercent float64 `gorm:"default:0" json:"SECFeePercent"`
	RegionName    string  `json:"regionName"`
}
