package users

import (
	"fmt"
	"log"
	"mime/multipart"
	userModels "trovo-wallet-api/internal/components/users/models"
	"trovo-wallet-api/internal/sharedconfig"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func GetTokenizedAssetSectorList(db *gorm.DB) (sectors []userModels.TokenizedAssetSector) {
	sectors = make([]userModels.TokenizedAssetSector, 0)
	db.Preload(clause.Associations).Order("id").Find(&sectors)

	return
}

func GetTokenizedAssetSubSectorList(db *gorm.DB) (subSectors []userModels.TokenizedAssetSubSector) {
	subSectors = make([]userModels.TokenizedAssetSubSector, 0)
	db.Preload(clause.Associations).Order("id").Find(&subSectors)

	return
}

func GetTokenizedAssetTypes(db *gorm.DB) (assetTypes []userModels.TokenizedAssetType) {
	assetTypes = make([]userModels.TokenizedAssetType, 0)
	db.Preload(clause.Associations).Order("asset_type").Find(&assetTypes)

	return
}

func GetAssetProceedCycle(db *gorm.DB) (apo []userModels.ProceedCycle) {
	apo = make([]userModels.ProceedCycle, 0)
	db.Preload(clause.Associations).Find(&apo)

	return
}
func GetAssetProtectionOptions(db *gorm.DB) (apo []userModels.AssetProtectionOption) {
	apo = make([]userModels.AssetProtectionOption, 0)
	db.Preload(clause.Associations).Order("id").Find(&apo)

	return
}

func GetApprovedAssetCustodians(db *gorm.DB) (custodians []userModels.ApprovedAssetCustodian) {
	custodians = make([]userModels.ApprovedAssetCustodian, 0)
	db.Preload(clause.Associations).Order("asset_custodian_name").Find(&custodians)

	return
}

func GetTokenizationFees(db *gorm.DB) (fees []userModels.TokenizationFee) {
	fees = make([]userModels.TokenizationFee, 0)
	db.Preload(clause.Associations).Find(&fees)

	return
}
func GetTokenizationCurrencies(db *gorm.DB) (currencies []userModels.TokenizationCurrency) {
	currencies = make([]userModels.TokenizationCurrency, 0)
	db.Preload(clause.Associations).Order("asset_code").Find(&currencies)

	return
}
func GetTokenizationPublicAssetAllowedCountries(db *gorm.DB) (c []userModels.TokenizationPublicAssetAllowedCountryCode) {
	c = make([]userModels.TokenizationPublicAssetAllowedCountryCode, 0)
	db.Preload(clause.Associations).Order("id").Find(&c)

	return
}

func GetTokenizedAssetTypesBySubsectorId(subSectorID string, db *gorm.DB) (assetTypes []userModels.TokenizedAssetType) {
	assetTypes = make([]userModels.TokenizedAssetType, 0)
	db.Preload(clause.Associations).Order("asset_type").Where("tokenized_asset_sub_sector_id = ?", subSectorID).Find(&assetTypes)

	return
}
func GetTokenizationDocumentById(id string, db *gorm.DB) (documents []userModels.AssetTokenizationDocument) {
	documents = make([]userModels.AssetTokenizationDocument, 0)
	db.Preload(clause.Associations).Where("tokenized_asset_id = ?", id).Find(&documents)

	return
}

func UploadTokenizationDocument(user *userModels.User, file multipart.File, fileNameWithExt string, input *userModels.AssetTokenizationInputDocument, gc *sharedconfig.GlobalConfig) (string, error) {

	newThumbnail, err := gc.FirebaseStorageUploader.UploadFile(file, fileNameWithExt, "")
	if err != nil {
		return "", err
	}

	//update the thumbnail url
	url := fmt.Sprintf("https://storage.googleapis.com/%v/%v", gc.FirebaseStorageUploader.BucketName, newThumbnail)
	documentUpload := userModels.AssetTokenizationDocument{
		TokenizedAssetID: input.TokenizedAssetID,
		DocumentType:     input.DocumentType,
		DocumentTitle:    input.DocumentTitle,
		DocumentUrl:      url,
	}
	e := gc.DB.Create(&documentUpload).Error
	if e != nil {
		log.Printf("[UploadTokenizationDocument]error saving document [%v] for %v: %v\n", input, user.Username, e)
		return "", fmt.Errorf("error saving document %v", input.DocumentTitle)
	}
	user.InvalidateUserCache(gc)
	owner, err := userModels.Username(user.Username).GetFullUser(gc.DB, gc)
	if err == nil {
		if owner.Username == user.Username {
			user = &owner
		}

	}

	return url, nil
}
