package models

import "time"

// CuratedAsset is a writable mirror of app-backend's curated asset catalog
// (internal/components/assets/models.CuratedAsset in the app-backend repo) -
// unlike the read-only P2P mirrors in internal/models/p2p.go, tm-api DOES
// write to this table directly: curating which assets exist on the
// platform (and, via P2PEnabled, which of those are available on the P2P
// marketplace) is an admin/config concern with no app-backend business
// logic of its own to preserve, the same as the other admin-managed
// catalog tables (KYC/faucet/fee configs) tm-api already owns.
type CuratedAsset struct {
	ID                          uint64    `gorm:"column:id;primaryKey" json:"id"`
	CreatedAt                   time.Time `json:"createdAt"`
	UpdatedAt                   time.Time `json:"updatedAt"`
	AssetCode                   string    `gorm:"column:asset_code" json:"assetCode"`
	AssetName                   string    `gorm:"column:asset_name" json:"assetName"`
	ContractAddress              string    `gorm:"column:contract_address" json:"contractAddress"`
	Description                  string    `gorm:"column:description" json:"description"`
	ImageURL                      *string   `gorm:"column:image_url" json:"imageUrl"`
	Website                        string    `gorm:"column:website" json:"website"`
	AssetConditions                 string    `gorm:"column:asset_conditions" json:"assetConditions"`
	AssetLimit                       float64   `gorm:"column:asset_limit" json:"assetLimit"`
	AssetRedemptionInstructions       string    `gorm:"column:asset_redemption_instructions" json:"assetRedemptionInstructions"`
	ContactEmail                       string    `gorm:"column:contact_email" json:"contactEmail"`
	Priority                            uint64    `gorm:"column:priority" json:"priority"`
	AssetClassID                        uint64    `gorm:"column:asset_class_id" json:"assetClassId"`
	Organization                        string    `gorm:"column:organization" json:"organization"`
	Withdrawable                        uint64    `gorm:"column:withdrawable" json:"withdrawable"`
	GenerateDepositAddress               uint64    `gorm:"column:generate_deposit_address" json:"generateDepositAddress"`
	DecimalPlaces                        uint64    `gorm:"column:decimal_places" json:"decimalPlaces"`
	RealAssetImageURL                     *string   `gorm:"column:real_asset_image_url" json:"realAssetImageUrl"`
	Inactive                              uint64    `gorm:"column:inactive" json:"inactive"`
	ClosedGroup                            *string   `gorm:"column:closed_group" json:"closedGroup"`
	// P2PEnabled gates this asset's availability on the P2P marketplace -
	// an offer cannot be created for an asset with this unset, and any
	// existing offer for it drops out of marketplace search. See
	// app-backend's CuratedAsset.P2PEnabled doc for the full contract.
	P2PEnabled bool `gorm:"column:p2p_enabled" json:"p2pEnabled"`
}

func (CuratedAsset) TableName() string { return "curated_assets" }

// AssetClass is the read-only reference list a curated asset is
// classified under (token/stablecoin/sto/nft) - mirrors app-backend's
// asset_classes table, used here only to populate the admin form's
// category dropdown.
type AssetClass struct {
	ID         uint64 `gorm:"column:id;primaryKey" json:"id"`
	AssetClass string `gorm:"column:asset_class" json:"assetClass"`
}

func (AssetClass) TableName() string { return "asset_classes" }

// CuratedAssetRequest is the create/update payload for the curated-asset
// management endpoint - booleans in, uint64 (0/1) columns out, since
// that's how app-backend's own model stores Withdrawable/
// GenerateDepositAddress/Inactive (a legacy convention predating this
// admin surface, not worth changing here).
type CuratedAssetRequest struct {
	Action                       string  `json:"action" binding:"required,oneof=create update" enums:"create,update"`
	ID                            uint64  `json:"id,omitempty"`
	AssetCode                     string  `json:"assetCode" binding:"required"`
	AssetName                     string  `json:"assetName"`
	ContractAddress                string  `json:"contractAddress"`
	Description                    string  `json:"description"`
	ImageURL                        string  `json:"imageUrl"`
	Website                          string  `json:"website"`
	AssetConditions                   string  `json:"assetConditions"`
	AssetLimit                        float64 `json:"assetLimit"`
	AssetRedemptionInstructions        string  `json:"assetRedemptionInstructions"`
	ContactEmail                       string  `json:"contactEmail"`
	Priority                            uint64  `json:"priority"`
	AssetClassID                        uint64  `json:"assetClassId"`
	Organization                        string  `json:"organization"`
	Withdrawable                        bool    `json:"withdrawable"`
	GenerateDepositAddress               bool    `json:"generateDepositAddress"`
	DecimalPlaces                        uint64  `json:"decimalPlaces"`
	RealAssetImageURL                     string  `json:"realAssetImageUrl"`
	Inactive                              bool    `json:"inactive"`
	ClosedGroup                           string  `json:"closedGroup"`
	P2PEnabled                            bool    `json:"p2pEnabled"`
}

// CuratedAssetListRequest filters the paginated curated-assets admin table.
type CuratedAssetListRequest struct {
	Page         int
	PageSize     int
	AssetCode    string
	AssetClassID uint64
	P2PEnabled   *bool
	Inactive     *bool
}
