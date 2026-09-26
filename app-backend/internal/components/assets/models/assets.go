package assets

import (
	"time"
)

// CuratedAsset model struct for CuratedAsset.
type CuratedAsset struct {
	ID                          uint64     `gorm:"primaryKey" json:"-"`
	CreatedAt                   time.Time  `json:"-"`
	UpdatedAt                   time.Time  `json:"-"`
	AssetCode                   string     `gorm:"size:12;unique;not null; default:''" json:"assetCode"`
	AssetName                   string     `gorm:"size:50;null; default:''" json:"assetName"`
	ContractAddress             string     `gorm:"size:56;not null; default:''" json:"contractAddress"`
	Description                 string     `gorm:"not null" json:"description"`
	ImageURL                    *string    `gorm:"null" json:"imageUrl"`
	Website                     string     `gorm:"null;size:100" json:"website"`
	AssetConditions             string     `gorm:"null;size:100" json:"assetConditions"`
	AssetLimit                  float64    `gorm:"type:integer;not null;default:0" json:"assetLimit"` //0 = unlimited
	AssetRedemptionInstructions string     `gorm:"null;" json:"assetRedemptionInstructions"`
	ContactEmail                string     `gorm:"null;size:100" json:"contactEmail"`
	Priority                    uint64     `gorm:"null;" json:"-"`
	AssetClassID                uint64     `gorm:"not null; default:1" json:"assetClassId"`
	AssetClass                  AssetClass `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"assetClass"`
	Organization                string     `gorm:"null;size:100" json:"organization"`
	Withdrawable                uint64     `gorm:"type:integer;not null;default:0" json:"withdrawable"`
	GenerateDepositAddress      uint64     `gorm:"type:integer;not null;default:0" json:"generateDepositAddress"`
	DecimalPlaces               uint64     `gorm:"type:integer;not null;default:7" json:"decimalPlaces"`
	RealAssetImageURL           *string    `gorm:"null;" json:"realAssetImageUrl"`
	Inactive                    uint64     `gorm:"type:integer;not null;default:0" json:"-"`
	ClosedGroup                 *string    `gorm:"null;" json:"closedGroup"`
	// P2PEnabled gates this asset's availability on the P2P marketplace,
	// independent of Withdrawable/Inactive (which gate unrelated wallet
	// features) - curated from the admin panel (tm-api). An asset with
	// this unset cannot be selected when creating a P2P offer, and any
	// existing offer for it disappears from marketplace search.
	// column:p2p_enabled is explicit because GORM's default namer splits
	// "P2PEnabled" on the digit (p2_p_enabled), not into the p2p_enabled
	// column name every other P2P table/query in this codebase expects.
	P2PEnabled bool `gorm:"column:p2p_enabled;not null;default:false" json:"p2pEnabled"`
}

// CuratedAsset model struct for CuratedAsset.
type CuratedSwapAsset struct {
	ID                          uint64     `gorm:"primaryKey" json:"-"`
	CreatedAt                   time.Time  `json:"-"`
	UpdatedAt                   time.Time  `json:"-"`
	AssetCode                   string     `gorm:"size:12;unique;not null" json:"assetCode"`
	AssetName                   string     `gorm:"size:50;null" json:"assetName"`
	ContractAddress             string     `gorm:"size:56;not null;" json:"contractAddress"`
	Description                 string     `gorm:"not null" json:"description"`
	ImageURL                    string     `gorm:"null" json:"imageUrl"`
	Website                     string     `gorm:"null;size:100" json:"website"`
	AssetConditions             string     `gorm:"null;size:100" json:"assetConditions"`
	AssetLimit                  float64    `gorm:"type:integer;not null;default:0" json:"assetLimit"` //0 = unlimited
	AssetRedemptionInstructions string     `gorm:"null;size:100" json:"assetRedemptionInstructions"`
	ContactEmail                string     `gorm:"null;size:100" json:"contactEmail"`
	Priority                    uint64     `gorm:"null;unique" json:"-"`
	AssetClassID                uint64     `gorm:"not null; default:1" json:"assetClassId"`
	AssetClass                  AssetClass `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"assetClass"`
	Organization                string     `gorm:"null;size:100" json:"organization"`
	Withdrawable                uint64     `gorm:"type:integer;not null;default:0" json:"withdrawable"`
	GenerateDepositAddress      uint64     `gorm:"type:integer;not null;default:0" json:"generateDepositAddress"`
	DecimalPlaces               uint64     `gorm:"type:integer;not null;default:7" json:"decimalPlaces"`
	RealAssetImageURL           string     `gorm:"null;" json:"realAssetImageUrl"`
	Inactive                    uint64     `gorm:"type:integer;not null;default:1" json:"-"`
	ClosedGroup                 string     `gorm:"null;" json:"closedGroup"`
}

// PaginatedCuratedAssets returns records sent for search
type PaginatedCuratedAssets struct {
	Pages        int            `json:"pages"`
	CurrentPage  int            `json:"currentPage"`
	TotalRecords int            `json:"totalRecords"`
	Limit        int            `json:"limit"`
	Records      []CuratedAsset `json:"records"`
}

// CuratedAssetOutput holds output for curated assets
type CuratedAssetOutput struct {
	CuratedAsset CuratedAsset `json:"curatedAssets"`
	Trusted      bool         `json:"trusted"`
}

// PaginatedBlockchainAssets returns records sent for search
type PaginatedBlockchainAssets struct {
	PageCursor string            `json:"pageCursor"`
	Assets     []BlockchainAsset `json:"assets"`
}

// BlockchainAsset holds blochcain assets
type BlockchainAsset struct {
	AssetCode       string
	ContractAddress string
	AmountOfTokens  string
	NumOfAccounts   int64
	AuthRequired    bool
	AuthRevocable   bool
	AuthImmutable   bool
	Toml            string
}

// AssetClass model struct for CuratedAsset.
// token, stablecoin, sto (security token) and nft (non fungible token)
type AssetClass struct {
	ID         uint64 `gorm:"primaryKey" json:"-"`
	AssetClass string `gorm:"size:45;unique;not null" json:"assetClass"`
}

// AssetClassOutput model struct for CuratedAsset.
// token, stablecoin, sto (security token) and nft (non fungible token)
type AssetClassOutput struct {
	ID         uint64 `json:"id"`
	AssetClass string `json:"assetClass"`
}

type GasDollarPrice struct {
	ID            string `gorm:"primaryKey"`
	Asset         string `gorm:"index:idx_dollar_price_unique_asset,unique"`
	Source        string `gorm:"size:100;index:idx_dollar_price_unique_asset,unique"`
	AskRate       string `gorm:"size:100"`
	BidRate       string `gorm:"size:100"`
	LastTradeRate string `gorm:"size:100"`
	LastUpdated   time.Time
}
type GasMarketChart struct {
	ID          string `gorm:"primaryKey"`
	Asset       string `gorm:"index:idx_gas_chart_unique_asset,unique"`
	Source      string `gorm:"size:100;index:idx_gas_chart_unique_asset,unique"`
	ChartString string
	LastUpdated time.Time
}
