package users

import (
	userModels "trovo-wallet-api/internal/components/users/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func GetTokenizedAssetSectorList(db *gorm.DB) (sectors []userModels.TokenizedAssetSector) {
	sectors = make([]userModels.TokenizedAssetSector, 0)
	db.Preload(clause.Associations).Find(&sectors)

	return
}

func GetTokenizedAssetSubSectorList(db *gorm.DB) (subSectors []userModels.TokenizedAssetSubSector) {
	subSectors = make([]userModels.TokenizedAssetSubSector, 0)
	db.Preload(clause.Associations).Find(&subSectors)

	return
}

func GetTokenizedAssetTypes(db *gorm.DB) (assetTypes []userModels.TokenizedAssetType) {
	assetTypes = make([]userModels.TokenizedAssetType, 0)
	db.Preload(clause.Associations).Find(&assetTypes)

	return
}
func GetTokenizedAssetTypesBySubsectorId(subSectorID string, db *gorm.DB) (assetTypes []userModels.TokenizedAssetType) {
	assetTypes = make([]userModels.TokenizedAssetType, 0)
	db.Preload(clause.Associations).Where("tokenized_asset_sector_id = ?", subSectorID).Find(&assetTypes)

	return
}
