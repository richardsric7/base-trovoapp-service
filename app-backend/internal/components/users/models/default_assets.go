package users

type DefaultAsset struct {
	ID              uint64 `gorm:"primaryKey" json:"-"`
	AssetCode       string `gorm:"size:12" json:"assetCode"`
	ContractAddress string `gorm:"size:56" json:"contractAddress"`
	ImageURL        string `json:"imageUrl"`
}

type DefaultAssets []DefaultAsset
