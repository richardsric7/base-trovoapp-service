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

func GetAssetProceedCycle(db *gorm.DB) (apo []userModels.ProceedCycle) {
	apo = make([]userModels.ProceedCycle, 0)
	db.Preload(clause.Associations).Find(&apo)

	return
}
func GetAssetProtectionOptions(db *gorm.DB) (apo []userModels.AssetProtectionOption) {
	apo = make([]userModels.AssetProtectionOption, 0)
	db.Preload(clause.Associations).Find(&apo)

	return
}

func GetApprovedAssetCustodians(db *gorm.DB) (custodians []userModels.ApprovedAssetCustodian) {
	custodians = make([]userModels.ApprovedAssetCustodian, 0)
	db.Preload(clause.Associations).Find(&custodians)

	return
}

func GetTokenizationFees(db *gorm.DB) (fees []userModels.TokenizationFee) {
	fees = make([]userModels.TokenizationFee, 0)
	db.Preload(clause.Associations).Find(&fees)

	return
}
func GetTokenizationCurrencies(db *gorm.DB) (currencies []userModels.TokenizationCurrency) {
	currencies = make([]userModels.TokenizationCurrency, 0)
	db.Preload(clause.Associations).Find(&currencies)

	return
}

func GetTokenizedAssetTypesBySubsectorId(subSectorID string, db *gorm.DB) (assetTypes []userModels.TokenizedAssetType) {
	assetTypes = make([]userModels.TokenizedAssetType, 0)
	db.Preload(clause.Associations).Where("tokenized_asset_sub_sector_id = ?", subSectorID).Find(&assetTypes)

	return
}
func GetTokenizationDocumentById(id string, db *gorm.DB) (documents []userModels.AssetTokenizationDocument) {
	documents = make([]userModels.AssetTokenizationDocument, 0)
	db.Preload(clause.Associations).Where("tokenized_asset_id = ?", id).Find(&documents)

	return
}
