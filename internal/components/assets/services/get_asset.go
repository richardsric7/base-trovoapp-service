package assets

import (
	assetblockchain "trovo-wallet-api/internal/components/assets/blockchain"
	assets "trovo-wallet-api/internal/components/assets/db"
	assetmodels "trovo-wallet-api/internal/components/assets/models"

	"gorm.io/gorm"
)

// GetCuratedAssets gets curated information
func GetCuratedAssets(db *gorm.DB) (assetsInfo map[string]assetmodels.CuratedAsset, err error) {

	assetsInfo, _ = assets.GetCuratedAssets(false, db)

	return

}

// GetBlockchainAssets gets curated information
func GetBlockchainAssets(assetCode, assetIssuer, cursor, order string, limit uint, db *gorm.DB) (assetsInfo assetmodels.PaginatedBlockchainAssets, err error) {

	assetsInfo, _ = assetblockchain.GetBlockchainAsset(assetCode, assetIssuer, cursor, order, limit)

	return

}
