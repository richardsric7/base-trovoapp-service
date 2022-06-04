package users

type DefaultAsset struct {
	ID          uint64
	AssetCode   string `gorm:"size:12" json:"assetCode"`
	AssetIssuer string `gorm:"size:56" json:"assetIssuer"`
}
