package users

import (
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"strconv"
	"strings"
	userModels "trovo-wallet-api/internal/components/users/models"
	db "trovo-wallet-api/internal/db"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
func GetTokenizedAssetByID(id string, db *gorm.DB) (tokenizedAsset userModels.TokenizedAsset, err error) {
	// var ta userModels.TokenizedAsset
	err = db.Preload(clause.Associations).Where("id = ?", id).First(&tokenizedAsset).Error

	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			//critical database error occured
			log.Printf("[GetTokenizedAssetByID]error fetching existing tokenization with id %v from database  [%v]", id, err)
			return

		}
	}

	return
}

func UploadTokenizationDocument(user *userModels.User, file multipart.File, fileNameWithExt string, input *userModels.AssetTokenizationInputDocument, gc *sharedconfig.GlobalConfig) (string, error) {

	newThumbnail, err := gc.FirebaseStorageUploader.UploadFile(file, fileNameWithExt, "")
	if err != nil {
		return "", err
	}

	//update the thumbnail url
	url := fmt.Sprintf("https://storage.googleapis.com/%v/%v", gc.FirebaseStorageUploader.BucketName, newThumbnail)
	//check if document already saved and then retireve it:
	var documentUpload userModels.AssetTokenizationDocument
	e := gc.DB.Where("tokenized_asset_id = ? AND document_type = ? AND document_title = ?", input.TokenizedAssetID, input.DocumentType, input.DocumentTitle).First(&documentUpload).Error
	if e == nil {
		//existing record match, update
		documentUpload.DocumentUrl = url
		es := gc.DB.Save(&documentUpload).Error
		if es != nil {

			log.Printf("[UploadTokenizationDocument]error saving existing document in database  [%v] for %v: %v\n", input, user.Username, e)
			return "", fmt.Errorf("error saving document %v", input.DocumentTitle)

		}
	} else {
		if !errors.Is(e, gorm.ErrRecordNotFound) {
			//critical database error occured
			log.Printf("[UploadTokenizationDocument]error fetching existing document from database  [%v] for %v: %v\n", input, user.Username, e)
			return "", fmt.Errorf("error saving document %v", input.DocumentTitle)

		}
		documentUpload = userModels.AssetTokenizationDocument{
			TokenizedAssetID: input.TokenizedAssetID,
			DocumentType:     input.DocumentType,
			DocumentTitle:    input.DocumentTitle,
			DocumentUrl:      url,
		}
		e := gc.DB.Create(&documentUpload).Error
		if e != nil {
			log.Printf("[UploadTokenizationDocument]error creating document in database  [%v] for %v: %v\n", input, user.Username, e)
			return "", fmt.Errorf("error saving document %v", input.DocumentTitle)
		}
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

func DeleteTokenizationDocument(user *userModels.User, documentID uint64, gc *sharedconfig.GlobalConfig) (document userModels.AssetTokenizationDocument, err error) {
	// var document userModels.AssetTokenizationDocument
	gc.DB.Where("id = ?", documentID).First(&document)

	if document.ID != documentID || documentID == 0 {
		return document, fmt.Errorf("error invalid document id %v", documentID)
	}
	err = gc.FirebaseStorageUploader.DeleteFile(document.DocumentUrl)
	if err != nil {
		log.Printf("[UploadTokenizationDocument]error deleting existing document in database  [%v] for %v: %v\n", document, user.Username, err)

		return document, err
	}

	e := gc.DB.Delete(&document).Error
	if e != nil {
		log.Printf("[UploadTokenizationDocument]error deleting existing document in database  [%v] for %v: %v\n", document, user.Username, e)
		return document, fmt.Errorf("error deleting document %v", document.DocumentTitle)

	}

	user.InvalidateUserCache(gc)
	owner, err := userModels.Username(user.Username).GetFullUser(gc.DB, gc)
	if err == nil {
		if owner.Username == user.Username {
			user = &owner
		}

	}

	return document, nil
}

func SubmitTokenizationAssetInfo(initiator *userModels.User, issuingWallet *userModels.UserWallet, input *userModels.TokenizedAssetJSONInput, gc *sharedconfig.GlobalConfig) (ato userModels.TokenizedAsset, err error) {

	//check if existing
	e := gc.DB.Where("asset_tokenization_status = ?", 0).First(&ato).Error
	if e == nil {
		//update existing
		ato.UpdateFromInput(input)

		ato.LastUpdatedBy = &initiator.Username

	} else {
		if !errors.Is(e, gorm.ErrRecordNotFound) {
			//critical database error occured
			log.Printf("[SubmitTokenizationAssetInfo]error fetching existing document from database  [%v] for %v: %v\n", input, initiator.Username, e)
			err = &tErrors.ErrorTemporaryServerError{}
			return

		}

		//create new tokenization
		ato = userModels.TokenizedAsset{
			ID:                     uuid.NewString(),
			InitiatorUsername:      initiator.Username,
			IssuingWalletPublicKey: issuingWallet.ID,
			IssuingWalletAlias:     issuingWallet.Alias,
		}
		ato.UpdateFromInput(input)

	}

	e = gc.DB.Save(&ato).Error
	if e != nil {
		log.Printf("[SubmitTokenizationAssetInfo]error saving  tokenization to database  [%v] for %v: %v\n", input, initiator.Username, e)

		err = &tErrors.ErrorTemporaryServerError{}
	}
	ato, _ = GetTokenizedAssetByID(ato.ID, gc.DB)
	return
}

func GetTokenizationList(user *userModels.User, gc *sharedconfig.GlobalConfig, c *gin.Context) (records userModels.PaginatedTokenizedAssets) {
	var err error
	var tokenizedAssetList []userModels.TokenizedAsset
	records.Records = make([]userModels.TokenizedAssetJSON, 0)
	DB, _ := db.OpenDb()
	DBC, _ := db.OpenDb()

	//get all wallets where user has access
	sharedWallets := make([]string, 0)
	for _, w := range user.WalletsSharedWithUser {
		sharedWallets = append(sharedWallets, w.WalletPublicKey)
	}
	var query *gorm.DB
	var countQuery *gorm.DB
	oD := "ASC"
	tokenizationStatus := strings.ToUpper(strings.TrimSpace(c.Query("assetTokenizationStatus")))
	assetDescription := strings.TrimSpace(strings.ToLower(c.Query("assetDescription")))
	assetType := strings.TrimSpace(strings.ToLower(c.Query("assetType")))
	assetName := strings.TrimSpace(strings.ToUpper(c.Query("assetName")))
	assetCode := strings.TrimSpace(strings.ToUpper(c.Query("assetCode")))
	assetSector := strings.TrimSpace(strings.ToLower(c.Query("assetSector")))
	assetSubSector := strings.TrimSpace(strings.ToLower(c.Query("assetSubSector")))
	initiatorUsername := strings.TrimSpace(strings.ToLower(c.Query("initiatorUsername")))
	offeringType := strings.TrimSpace(strings.ToUpper(c.Query("offeringType")))

	hasSecApproval := strings.TrimSpace(c.Query("hasSecApproval"))
	limitU, _ := strconv.ParseUint(strings.TrimSpace(c.DefaultQuery("limit", "25")), 10, 64)
	limit := int(limitU)
	pageU, _ := strconv.ParseUint(strings.TrimSpace(c.DefaultQuery("page", "1")), 10, 64)
	page := int(pageU)

	createdBetween := strings.TrimSpace(c.Query("createdBetween"))

	orderBy := strings.TrimSpace(c.DefaultQuery("orderby", "created_at"))
	orderDirection := c.DefaultQuery("order", "DESC")

	query = DB.Preload(clause.Associations)
	countQuery = DBC.Group("id")

	if len(orderDirection) > 0 && strings.ToLower(orderDirection) == "desc" {
		oD = "DESC"
	}
	if len(orderBy) > 0 {
		query = query.Order(orderBy + " " + oD)
		countQuery = countQuery.Order(orderBy + " " + oD)

	} else {
		query = query.Order("updated_at desc, ownershipType asc, asset_country_location asc")
		// query = query.Order("created_at DESC")
		// countQuery = countQuery.Order("created_at DESC")
		countQuery = countQuery.Order("updated_at desc, ownershipType asc, asset_country_location asc")
	}

	{
		query = query.Where("(issuing_wallet_public_key IN (?))", sharedWallets)
		countQuery = countQuery.Where("(issuing_wallet_public_key IN (?))", sharedWallets)

	}

	if len(initiatorUsername) > 2 {
		query = query.Where("lower(initiator_username) = ?", initiatorUsername)
		countQuery = countQuery.Where("lower(initiator_username) = ?", initiatorUsername)

	}
	if len(tokenizationStatus) > 0 {
		ts, _ := strconv.ParseUint(strings.TrimSpace(tokenizationStatus), 10, 64)
		tsInt := int(ts)
		query = query.Where("Asset_Tokenization_Status = ?", tsInt)
		countQuery = countQuery.Where("Asset_Tokenization_Status = ?", tokenizationStatus)

	}
	if len(hasSecApproval) > 0 {
		ts, _ := strconv.ParseUint(strings.TrimSpace(hasSecApproval), 10, 64)
		tsInt := int(ts)
		query = query.Where("sec_approval = ?", tsInt)
		countQuery = countQuery.Where("sec_approval = ?", tsInt)

	}

	if len(assetDescription) > 2 {

		query = query.Where("lower(asset_Description) LIKE ?", "%"+strings.ToLower(assetDescription)+"%")
		countQuery = countQuery.Where("lower(asset_Description) LIKE ?", "%"+strings.ToLower(assetDescription)+"%")

	}
	if len(offeringType) > 3 {

		query = query.Where("offering_type = ?", offeringType)
		countQuery = countQuery.Where("offering_type = ?", offeringType)

	}
	if len(assetType) > 4 {

		query = query.Where("lower(asset_type) = ?", strings.TrimSpace(assetType))
		countQuery = countQuery.Where("lower(asset_type) = ?", strings.TrimSpace(assetType))

	}
	if len(assetSector) > 0 {

		query = query.Where("lower(asset_sector) = ?", strings.TrimSpace(assetSector))
		countQuery = countQuery.Where("lower(asset_sector) = ?", strings.TrimSpace(assetSector))

	}
	if len(assetSubSector) > 0 {

		query = query.Where("lower(asset_sub_sector) = ?", strings.TrimSpace(assetSubSector))
		countQuery = countQuery.Where("lower(asset_sub_sector) = ?", strings.TrimSpace(assetSubSector))

	}
	if len(assetName) > 2 {

		query = query.Where("upper(asset_name) = ?", strings.TrimSpace(assetName))
		countQuery = countQuery.Where("upper(asset_name) = ?", strings.TrimSpace(assetName))

	}

	if len(assetCode) > 0 {

		query = query.Where("upper(asset_code) = ?", strings.TrimSpace(assetCode))
		countQuery = countQuery.Where("upper(asset_code) = ?", strings.TrimSpace(assetCode))

	}

	if len(createdBetween) == 21 && strings.Contains(createdBetween, "|") {
		// 2020-01-01|2020-02-31 full range date
		dateRange := strings.Split(createdBetween, "|")
		query = query.Where("created_at::date BETWEEN ?::date AND ?::date", dateRange[0], dateRange[1])
		countQuery = countQuery.Where("created_at::date BETWEEN ?::date AND ?::date", dateRange[0], dateRange[1])

	}

	var countR int64

	errCount := countQuery.Find(&[]userModels.TokenizedAsset{}).Count(&countR).Error
	if errCount != nil {
		log.Println("[GetTokenizationList]Count Error:", errCount)
		return records
	}

	if limit > 0 {
		query.Limit(limit)
	}
	count := int(countR)
	pages := 1
	if count > limit {
		// fmt.Println("count / limit = ", count/limit, "count%limit = ", count%limit)
		pages = count / limit
		if count%limit > 0 {
			pages = pages + 1
		}
	}
	if page > pages {
		page = pages
	}
	if page > 1 {
		// fmt.Println("Offset = ", (page-1)*limit)
		query.Offset(((page - 1) * limit))
	}
	if err = query.Find(&tokenizedAssetList).Error; err != nil {
		log.Println("[GetTokenizationList] Query Error:", err)
		return
	}

	tListJSON := make([]userModels.TokenizedAssetJSON, 0)

	for _, v := range tokenizedAssetList {
		j := v.ToJSON()
		tListJSON = append(tListJSON, j)

	}

	records = userModels.PaginatedTokenizedAssets{CurrentPage: page, Pages: pages, TotalRecords: count, Limit: limit, Records: tListJSON}

	return records
}
