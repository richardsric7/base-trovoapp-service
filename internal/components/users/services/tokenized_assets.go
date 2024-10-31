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
	"github.com/shopspring/decimal"
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
	db.Order("asset_custodian_name").Find(&custodians)

	return
}

func GetBanks(countryCode string, db *gorm.DB) (banks []userModels.Bank) {
	banks = make([]userModels.Bank, 0)
	db.Order("country_code,bank_name").Where("country_code = ?", countryCode).Find(&banks)

	return
}

func GetAssetManagers(db *gorm.DB) (assetManagers []userModels.AssetManager) {
	assetManagers = make([]userModels.AssetManager, 0)
	db.Order("asset_manager_country, asset_manager_name").Find(&assetManagers)

	return
}

func GetTokenizationFees(db *gorm.DB) (fees []userModels.TokenizationFee) {
	fees = make([]userModels.TokenizationFee, 0)
	// db.Preload(clause.Associations).Where("inactive != ?", 1).Find(&fees)
	db.Where("inactive != ?", 1).Find(&fees)

	return
}

func GetAssetTokenizationDocumentTypes(db *gorm.DB) (docTypes []userModels.AssetTokenizationDocumentType) {
	docTypes = make([]userModels.AssetTokenizationDocumentType, 0)
	// db.Preload(clause.Associations).Where("inactive != ?", 1).Find(&fees)
	db.Order("document_category,document_type_description").Find(&docTypes)

	return
}

func GetTokenizationFeeByID(feeID uint64, db *gorm.DB) (fee userModels.TokenizationFee) {
	db.Preload(clause.Associations).Where("id = ?", feeID).First(&fee)

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

		} else {
			//record not found
			err = &tErrors.CustomError{Param: "tokenizationID", Err: "error-invalid-tokenizationId", ErrMessage: fmt.Sprintf("Tokenization ID %v is invalid.", id)}
			return
		}
	}

	return
}

func GetTokenizedAssetByIssuingWallet(issuingWalletPublicKey string, db *gorm.DB) (tokenizedAsset userModels.TokenizedAsset, NotFound bool, err error) {
	// var ta userModels.TokenizedAsset
	err = db.Preload(clause.Associations).Where("issuing_wallet_public_key = ?", issuingWalletPublicKey).First(&tokenizedAsset).Error

	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			//critical database error occured
			log.Printf("[GetTokenizedAssetByIssuingWallet]error fetching existing tokenization with issuing wallet %v from database  [%v]", issuingWalletPublicKey, err)
			return

		}
		NotFound = true
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

	// initialize message array
	input.Messages = make([]string, 0)

	//check if existing
	ato, NotFound, e := GetTokenizedAssetByIssuingWallet(issuingWallet.ID, gc.DB)

	if e == nil {
		//tokenization existing
		if ato.AssetTokenizationStatus > 0 {
			// error tokenization is already in progress
			log.Printf("[SubmitTokenizationAssetInfo] Error tokenization procesing is in progress and cannot be modified: %v\n", issuingWallet.ID)
			err = &tErrors.CustomError{Param: "issuingWalletPublicKey", Err: "error-tokenization-cannot-be-modified-by-this-method", ErrMessage: "Tokenization cannot be modified by this method."}
			return

		}
		ato = UpdateFromInput(&ato, input, gc)

		ato.LastUpdatedBy = &initiator.Username

	} else {
		if !NotFound {
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
		ato = UpdateFromInput(&ato, input, gc)

	}

	e = gc.DB.Omit(clause.Associations).Save(&ato).Error
	if e != nil {
		log.Printf("[SubmitTokenizationAssetInfo] error saving tokenization to database  [%v] for %v: %v\n", input, initiator.Username, e)

		err = &tErrors.ErrorTemporaryServerError{}

	}
	ato, _ = GetTokenizedAssetByID(ato.ID, gc.DB)
	return ato, err
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
		j := v.ToJSON(gc)
		tListJSON = append(tListJSON, j)

	}

	records = userModels.PaginatedTokenizedAssets{CurrentPage: page, Pages: pages, TotalRecords: count, Limit: limit, Records: tListJSON}

	return records
}

func UpdateFromInput(t *userModels.TokenizedAsset, ti *userModels.TokenizedAssetJSONInput, gc *sharedconfig.GlobalConfig) userModels.TokenizedAsset {
	if ti.HasAdditionalKYCRequirements > 0 && len(ti.AdditionalKYCRequirements) > 0 {

		t.AdditionalKYCRequirements = &ti.AdditionalKYCRequirements

	}

	if len(ti.AssetSector) > 0 {

		t.AssetSector = &ti.AssetSector
	}

	if len(ti.AssetSubSector) > 0 {

		t.AssetSubSector = &ti.AssetSubSector
	}

	if len(ti.AssetType) > 0 {

		t.AssetType = &ti.AssetType
	}

	if len(ti.AssetName) > 0 {

		t.AssetName = &ti.AssetName
	}

	t.ApprovedAssetCustodianID = ti.ApprovedAssetCustodianID

	if len(ti.OfferingType) > 0 {

		t.OfferingType = &ti.OfferingType
	}

	if len(ti.ClosedGroupID) > 0 {

		t.ClosedGroupID = &ti.ClosedGroupID
	}

	if len(ti.SecApprovalIdNumber) > 0 && ti.SecApproval > 0 {

		t.SecApproval = ti.SecApproval

		t.SecApprovalIdNumber = &ti.SecApprovalIdNumber
	}

	if len(ti.MarketMakingWallet) > 0 {

		t.MarketMakingWallet = &ti.MarketMakingWallet
	}

	if len(ti.AssetDescription) > 0 {

		t.AssetDescription = &ti.AssetDescription
	}

	if len(ti.AssetCountryLocation) > 0 {

		t.AssetCountryLocation = &ti.AssetCountryLocation
	}

	if len(ti.AssetPhysicalAddress) > 0 {

		t.AssetPhysicalAddress = &ti.AssetPhysicalAddress
	}

	if len(ti.AssetLongitude) > 0 {

		t.AssetLongitude = &ti.AssetLongitude
	}

	if len(ti.AssetLatitude) > 0 {

		t.AssetLatitude = &ti.AssetLatitude
	}

	if len(ti.OwnershipType) > 0 {

		t.OwnershipType = &ti.OwnershipType
	}

	if len(ti.OwnershipKind) > 0 {

		t.OwnershipKind = &ti.OwnershipKind
	}

	if len(ti.AssetOwnerName) > 0 {

		t.AssetOwnerName = &ti.AssetOwnerName
	}

	if len(ti.AssetOwnerAddress) > 0 {

		t.AssetOwnerAddress = &ti.AssetOwnerAddress
	}

	// if ti.AssetManagerID > 0 {

	t.AssetManagerID = ti.AssetManagerID
	// }

	// if len(ti.AssetManagerAddress) > 0 {

	// 	t.AssetManagerAddress = &ti.AssetManagerAddress
	// }
	if len(ti.AssetQuoteCurrency) > 0 {

		t.AssetQuoteCurrency = &ti.AssetQuoteCurrency
	}

	t.AssetCurrentValue = ti.AssetCurrentValue
	t.AssetMscCostOutisdeOfValuation = ti.AssetMscCostOutisdeOfValuation
	t.ValueOfTokenizedAsset = ti.ValueOfTokenizedAsset

	if len(ti.ProtectionMethods) > 0 {

		t.ProtectionMethods = &ti.ProtectionMethods
	}

	if len(ti.InsuranceCompanyName) > 0 {

		t.InsuranceCompanyName = &ti.InsuranceCompanyName
	}

	if len(ti.InsurancePolicyHolder) > 0 {

		t.InsurancePolicyHolder = &ti.InsurancePolicyHolder
	}

	if len(ti.InsurancePolicyNumber) > 0 {

		t.InsurancePolicyNumber = &ti.InsurancePolicyNumber
	}

	t.PercentageValueOfInsurance = ti.PercentageValueOfInsurance
	t.IsFreeFromLiensAndEncumbrances = ti.IsFreeFromLiensAndEncumbrances
	t.AssetAlreadyExists = ti.AssetAlreadyExists

	if len(ti.AssetCode) > 0 {

		t.AssetCode = &ti.AssetCode
	}

	if len(ti.AssetLogo) > 0 {

		t.AssetLogo = &ti.AssetLogo
	}
	if t.NumberOfTokenToBeIssued > 0 && t.ValueOfTokenizedAsset > 0 {
		totalValuation := (ti.ValueOfTokenizedAsset + ti.AssetMscCostOutisdeOfValuation)
		t.PricePerToken = decimal.NewFromFloat(totalValuation / t.NumberOfTokenToBeIssued).Truncate(7).InexactFloat64()
	}
	var feeCompo userModels.TokenizationFee
	var assetFee float64

	if ti.TokenizationFeeID > 0 {
		// fee has been selected
		t.TokenizationFeeID = &ti.TokenizationFeeID
		feeCompo = t.UpdateTokenizationFeeByID(ti.TokenizationFeeID, gc)

		{
			// Calculate Fees
			assetFee = decimal.NewFromFloat(t.NumberOfTokenToBeIssued * (feeCompo.FeeAssetPercentage / 100)).Truncate(7).InexactFloat64()

		}

	}

	t.NumberOfTokenToBeIssued = ti.NumberOfTokenToBeIssued
	t.NumberOfTokenToBeSold = ti.NumberOfTokenToBeSold
	// auto calculate, token to be held is less the fee

	{
		//ensure correct the number of token to be sold.
		maxTokenToBeSold := t.NumberOfTokenToBeIssued - assetFee

		if t.NumberOfTokenToBeSold > maxTokenToBeSold {
			t.NumberOfTokenToBeSold = maxTokenToBeSold
			ti.NumberOfTokenToBeSold = maxTokenToBeSold
			ti.Messages = append(ti.Messages, fmt.Sprintf("Submitted Number of tokens to be sold has been adjusted to %v to account for deduction of the asset fee of %v. You may wish to adjust your fee option and then adjust the amount to be sold manually again.", maxTokenToBeSold, assetFee))

		}
	}
	t.TotalTokenHeldByManager = t.NumberOfTokenToBeIssued - t.NumberOfTokenToBeSold - assetFee
	if len(ti.WalletToHoldAssetsNotForSale) > 0 {

		t.WalletToHoldAssetsNotForSale = &ti.WalletToHoldAssetsNotForSale
	}

	/**



	**/

	t.SalesStart = ti.SalesStart
	t.SalesEnd = ti.SalesEnd
	t.CapOnPurchase = ti.CapOnPurchase
	t.CapQuantity = ti.CapQuantity
	t.CapDurationInDays = ti.CapDurationInDays

	if len(ti.ProceedCycle) > 0 {

		t.ProceedCycle = &ti.ProceedCycle
	}

	if len(ti.ProceedPayoutCurrency) > 0 {

		t.ProceedPayoutCurrency = &ti.ProceedPayoutCurrency
	}

	if len(ti.ExemptedCountries) > 0 {

		t.ExemptedCountries = &ti.ExemptedCountries
	}

	t.HasAdditionalKYCRequirements = ti.HasAdditionalKYCRequirements

	if len(ti.AdditionalKYCRequirements) > 0 && t.HasAdditionalKYCRequirements > 0 {

		t.AdditionalKYCRequirements = &ti.AdditionalKYCRequirements
	}

	t.InvestorAccreditationRequired = ti.InvestorAccreditationRequired
	return *t
}
