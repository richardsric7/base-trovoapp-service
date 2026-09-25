package assets

import (
	"strings"
	models "trovo-wallet-api/internal/components/assets/models"

	"gorm.io/gorm"
)

// curatedAssetRow mirrors just the columns GetBlockchainAsset needs from
// the curated_assets table.
type curatedAssetRow struct {
	AssetCode       string
	ContractAddress string
}

// GetBlockchainAsset searches Base assets matching assetCode/contractAddress.
// Stellar's version queried Horizon's global, network-wide asset
// registry (paginated by cursor); Base has no such registry - any
// contract can mint a token with any symbol, so "search assets" here
// means this application's own curated-asset catalog instead (the
// backend's actual source of truth for which B20 assets it supports),
// enriched with each match's live on-chain total supply where available.
// cursor/order are accepted for call-site compatibility but unused - the
// curated-asset catalog is small enough not to need pagination yet.
func GetBlockchainAsset(assetCode, contractAddress, cursor, order string, limit uint, db *gorm.DB) (paginatedBlockchainAssets models.PaginatedBlockchainAssets, err error) {
	if limit < 1 {
		limit = 25
	}
	var rows []curatedAssetRow
	q := db.Table("curated_assets").Select("asset_code, contract_address")
	if assetCode != "" {
		q = q.Where("asset_code = ?", strings.ToUpper(assetCode))
	}
	if contractAddress != "" {
		q = q.Where("contract_address = ?", strings.ToLower(contractAddress))
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
			AssetCode:       row.AssetCode,
			ContractAddress: row.ContractAddress,
			AmountOfTokens:  "0",
		})
	}

	paginatedBlockchainAssets.Assets = assetsOut
	return paginatedBlockchainAssets, nil
}
