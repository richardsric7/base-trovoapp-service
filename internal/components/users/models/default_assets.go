package users

type DefaultAsset struct {
	ID          uint64 `gorm:"primaryKey" json:"-"`
	AssetCode   string `gorm:"size:12" json:"assetCode"`
	AssetIssuer string `gorm:"size:56" json:"assetIssuer"`
	ImageURL    string `json:"imageUrl"`
}

type DefaultAssets []DefaultAsset
