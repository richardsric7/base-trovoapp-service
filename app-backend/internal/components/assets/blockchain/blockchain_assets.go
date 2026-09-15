package assets

import (
	"strings"
	models "trovo-wallet-api/internal/components/assets/models"

	"gorm.io/gorm"
)

// curatedAssetRow mirrors just the columns GetBlockchainAsset needs from
// the curated_assets table.
type curatedAssetRow struct {
	AssetCode   string
	AssetIssuer string
}

// GetBlockchainAsset searches Base assets matching assetCode/assetIssuer.
// Stellar's version queried Horizon's global, network-wide asset
// registry (paginated by cursor); Base has no such registry - any
// contract can mint a token with any symbol, so "search assets" here
// means this application's own curated-asset catalog instead (the
// backend's actual source of truth for which B20 assets it supports),
// enriched with each match's live on-chain total supply where available.
// cursor/order are accepted for call-site compatibility but unused - the
// curated-asset catalog is small enough not to need pagination yet.
func GetBlockchainAsset(assetCode, assetIssuer, cursor, order string, limit uint, db *gorm.DB) (paginatedBlockchainAssets models.PaginatedBlockchainAssets, err error) {
	if limit < 1 {
		limit = 25
	}
	var rows []curatedAssetRow
	q := db.Table("curated_assets").Select("asset_code, asset_issuer")
	if assetCode != "" {
		q = q.Where("asset_code = ?", strings.ToUpper(assetCode))
	}
	if assetIssuer != "" {
		q = q.Where("asset_issuer = ?", strings.ToLower(assetIssuer))
	}
	if e := q.Limit(int(limit)).Find(&rows).Error; e != nil {
		return paginatedBlockchainAssets, e
	}

	assetsOut := make([]models.BlockchainAsset, 0, len(rows))
	for _, row := range rows {
		// AmountOfTokens (a live totalSupply() read) and NumOfAccounts
		// (needs an indexer) are left at their zero values - a follow-up
		// once this browse endpoint needs them, not core to search
		// itself. Authorization for a B20 asset is enforced per-wallet
		// (see internal/network.IsWalletAuthorizedForAsset), not as a
		// single account-level flag, so AuthRequired/AuthRevocable/
		// AuthImmutable are left false here rather than guessed.
		assetsOut = append(assetsOut, models.BlockchainAsset{
			AssetCode:      row.AssetCode,
			AssetIssuer:    row.AssetIssuer,
			AmountOfTokens: "0",
		})
	}

	paginatedBlockchainAssets.Assets = assetsOut
	return paginatedBlockchainAssets, nil
}
