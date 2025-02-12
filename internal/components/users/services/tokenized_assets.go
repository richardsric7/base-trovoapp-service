package users

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
	userModels "trovo-wallet-api/internal/components/users/models"
	db "trovo-wallet-api/internal/db"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/middleware"
	"trovo-wallet-api/internal/network"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/txnbuild"
	"github.com/stellar/go/xdr"
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
func IsTokenizationMintingApprover(username string, db *gorm.DB) (access bool) {
	if username == "" {
		return
	}
	return db.Where("approver = ?", username).First(&userModels.TokenizationMintingApprover{}).Error == nil

}
func IsTokenizationMintingInitiator(username string, db *gorm.DB) (access bool) {
	if username == "" {
		return
	}
	return db.Where("initiator = ?", username).First(&userModels.TokenizationMintingInitiator{}).Error == nil

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
func GetAssetManagerByID(id uint64, db *gorm.DB) (assetManager userModels.AssetManager) {
	db.Where("id = ?", id).First(&assetManager)

	return
}

func GetTokenizationFees(db *gorm.DB) (fees []userModels.TokenizationFee) {
	fees = make([]userModels.TokenizationFee, 0)
	// db.Preload(clause.Associations).Where("inactive != ?", 1).Find(&fees)
	db.Where("inactive != ?", 1).Find(&fees)

	return
}

func GetTokenizationStatuses(db *gorm.DB) (statuses []userModels.TokenizationStatus) {
	statuses = make([]userModels.TokenizationStatus, 0)
	db.Find(&statuses)

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

func GetTokenizationFeePaymentMethods(db *gorm.DB) (feePMs []userModels.TokenizationFeePaymentMethod) {
	feePMs = make([]userModels.TokenizationFeePaymentMethod, 0)
	db.Preload(clause.Associations).Find(&feePMs)

	return
}

func GetTokenizationCurrencies(db *gorm.DB) (currencies []userModels.TokenizationCurrency) {
	currencies = make([]userModels.TokenizationCurrency, 0)
	db.Preload(clause.Associations).Order("asset_code").Find(&currencies)

	return
}

func GetTokenizationByDocumentID(did uint64, gc *sharedconfig.GlobalConfig) (t userModels.TokenizedAsset) {
	return userModels.AssetTokenizationDocumentID(did).GetTokenization(gc)
}

func GetTokenizationCurrencyByCode(code string, db *gorm.DB) (currency userModels.TokenizationCurrency) {
	db.Where("asset_code = ?", strings.ToUpper(code)).First(&currency)

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

func GetTokenizedAssetByID(id string, db *gorm.DB) (tokenizedAsset userModels.TokenizedAsset, NotFound bool, err error) {
	// var ta userModels.TokenizedAsset
	err = db.Preload(clause.Associations).Where("id = ?", id).First(&tokenizedAsset).Error

	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			//critical database error occured
			log.Printf("[GetTokenizedAssetByID]error fetching existing tokenization with id %v from database  [%v]", id, err)
			return

		} else {
			//record not found
			NotFound = true
			err = &tErrors.CustomError{Param: "tokenizationID", Err: "error-invalid-tokenizationId", ErrMessage: fmt.Sprintf("Tokenization ID %v is invalid.", id)}
			return
		}
	}

	return
}

// GetOpenTokenizedAssetByInitiatorUsername get the tokenization that has status 0 or 1 initiated by the initiator username
func GetOpenTokenizedAssetByInitiatorUsername(initiatorUsername string, db *gorm.DB) (tokenizedAsset userModels.TokenizedAsset, NotFound bool, err error) {
	// var ta userModels.TokenizedAsset
	err = db.Preload(clause.Associations).Where("asset_tokenization_status < 2 AND initiator_username = ?", initiatorUsername).First(&tokenizedAsset).Error

	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			//critical database error occured
			log.Printf("[GetTokenizedAssetByID]error fetching existing tokenization with initiatorUsername %v from database  [%v]", initiatorUsername, err)
			return

		} else {
			//record not found
			NotFound = true
			err = &tErrors.CustomError{Param: "tokenizationID", Err: "error-invalid-tokenizationId", ErrMessage: fmt.Sprintf("%v has no tokenized asset inititated", initiatorUsername)}
			return
		}
	}

	return
}

// GetTokenizedAssetByIssuingWallet gets the tokenized asset by issuing wallet and returns the first one ordered by the asset tokenization status from 0.
func GetTokenizedAssetByIssuingWallet(issuingWalletPublicKey string, db *gorm.DB) (tokenizedAsset userModels.TokenizedAsset, NotFound bool, err error) {
	// var ta userModels.TokenizedAsset
	err = db.Preload(clause.Associations).Order("asset_tokenization_status ASC").Where("issuing_wallet_public_key = ?", issuingWalletPublicKey).First(&tokenizedAsset).Error

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

func UploadTokenizationFeeProofOfPaymentDocument(user *userModels.User, tokenizedAssetID string, file multipart.File, fileNameWithExt string, feepaymentproofinput *userModels.TokenizationFeeProofOfPaymentInput, gc *sharedconfig.GlobalConfig) (string, error) {

	newThumbnail, err := gc.FirebaseStorageUploader.UploadFile(file, fileNameWithExt, "")
	if err != nil {
		return "", err
	}

	//update the thumbnail url
	url := fmt.Sprintf("https://storage.googleapis.com/%v/%v", gc.FirebaseStorageUploader.BucketName, newThumbnail)
	//check if document already saved and then retireve it:

	documentUpload := userModels.TokenizationFeeProofOfPayment{
		TokenizationFeePaymentMethodID: feepaymentproofinput.TokenizationFeePaymentMethodID,
		TokenizedAssetID:               tokenizedAssetID,
		TransactionReference:           &feepaymentproofinput.TransactionReference,
		DocumentUrl:                    url,
	}
	e := gc.DB.Create(&documentUpload).Error
	if e != nil {
		log.Printf("[UploadTokenizationFeeProofOfPaymentDocument]error creating proof of payment document in database for %v: %v\n", user.Username, e)
		return "", fmt.Errorf("error saving proof of payment for tokenization ID %v", tokenizedAssetID)
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

func UploadTokenizationAssetLogo(user *userModels.User, ato *userModels.TokenizedAsset, file multipart.File, fileNameWithExt string, gc *sharedconfig.GlobalConfig) (string, error) {

	newThumbnail, err := gc.FirebaseStorageUploader.UploadFile(file, fileNameWithExt, "")
	if err != nil {
		return "", err
	}

	//update the thumbnail url
	url := fmt.Sprintf("https://storage.googleapis.com/%v/%v", gc.FirebaseStorageUploader.BucketName, newThumbnail)
	//check if document already saved and then retireve it:
	ato.AssetLogo = &url
	e := gc.DB.Omit(clause.Associations).Save(&ato).Error
	if e != nil {
		log.Printf("[UploadTokenizationAssetLogo] error saving asset logo to database [%v] for %v: %v\n", ato.ID, user.Username, e)

		err = &tErrors.ErrorTemporaryServerError{}

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

func DeleteTokenization(user *userModels.User, tokenizationID string, gc *sharedconfig.GlobalConfig) (tokenizedAsset userModels.TokenizedAssetJSON, err error) {
	// var document userModels.AssetTokenizationDocument
	ato, _, err := GetTokenizedAssetByID(tokenizationID, gc.DB)
	if err != nil {
		log.Printf("[DeleteTokenization] error locating existing tokenization in database  [%v] for %v: %v\n", tokenizationID, user.Username, err)
		return tokenizedAsset, fmt.Errorf("error locating tokenization request with ID %v", tokenizationID)

	}
	if ato.IssuingWalletPublicKey != nil {
		if len(*ato.IssuingWalletPublicKey) != 56 {
			return tokenizedAsset, fmt.Errorf("error locating tokenization request with ID %v", tokenizationID)

		}
	}

	if ato.AssetTokenizationStatus > 0 {
		err = &tErrors.CustomError{Param: "ID", Err: "error-invalid-cannot-delete", ErrMessage: "You cannot delete this tokenization request because it has passed the editing stage."}

		return ato.ToJSON(gc), err

	}
	//get access permission and see if the user taking action has an INITIATOR access.
	if ato.IssuingWalletPublicKey != nil {
		if !userModels.UserWalletID(*ato.IssuingWalletPublicKey).UserHasAccess(user.Username, "INITIATOR", gc.DB) {
			err = &tErrors.CustomError{Param: "ID", Err: "error-invalid-access", ErrMessage: "You cannot delete this tokenization request because you do not possess an initiator permission on the tokenization wallet."}

			return tokenizedAsset, fmt.Errorf("error locating tokenization request with ID %v", tokenizationID)

		}
	}

	tx := gc.DB.Begin()
	defer tx.Rollback()
	if len(ato.AssetTokenizationDocuments) > 0 {
		e := tx.Delete(&ato.AssetTokenizationDocuments).Error
		if e != nil {
			log.Printf("[DeleteTokenization]error deleting existing tokenization documents in database  [%v] for %v: %v\n", tokenizationID, user.Username, e)
			return ato.ToJSON(gc), fmt.Errorf("error deleting tokenization request with ID %v", tokenizationID)

		}
	}
	//reset the document since it is purged
	ato.AssetTokenizationDocuments = make([]userModels.AssetTokenizationDocument, 0)
	e := tx.Delete(&ato).Error
	if e != nil {
		log.Printf("[DeleteTokenization]error deleting existing tokenization in database  [%v] for %v: %v\n", tokenizationID, user.Username, e)
		return ato.ToJSON(gc), fmt.Errorf("error deleting tokenization request with ID %v", tokenizationID)

	}
	tx.Commit()
	user.InvalidateUserCache(gc)
	owner, err := userModels.Username(user.Username).GetFullUser(gc.DB, gc)
	if err == nil {
		if owner.Username == user.Username {
			user = &owner
		}

	}

	return ato.ToJSON(gc), nil
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

func DeleteTokenizationFeePaymentDocument(user *userModels.User, documentID uint64, gc *sharedconfig.GlobalConfig) (document userModels.TokenizationFeeProofOfPayment, err error) {
	// var document userModels.AssetTokenizationDocument
	gc.DB.Where("id = ?", documentID).First(&document)

	if document.ID != documentID || documentID == 0 {
		return document, fmt.Errorf("error invalid document id %v", documentID)
	}
	err = gc.FirebaseStorageUploader.DeleteFile(document.DocumentUrl)
	if err != nil {
		log.Printf("[DeleteTokenizationFeePaymentDocument]error deleting existing document in database  [%v] for %v: %v\n", document, user.Username, err)

		return document, err
	}

	e := gc.DB.Delete(&document).Error
	if e != nil {
		log.Printf("[DeleteTokenizationFeePaymentDocument]error deleting existing document in database  [%v] for %v: %v\n", document, user.Username, e)
		return document, fmt.Errorf("error deleting document with ID %v", documentID)

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

func SubmitTokenizationAssetInfoByInitiator(initiator *userModels.User, input *userModels.TokenizedAssetJSONInput, gc *sharedconfig.GlobalConfig) (ato userModels.TokenizedAsset, err error) {

	// initialize message array
	input.Messages = make([]string, 0)
	// check asset manager ID
	// if input.AssetManagerID == 0 {
	// 	log.Printf("[SubmitTokenizationAssetInfoByInitiator] Error Invalid Asset Manager ID: %v\n", issuingWallet.ID)
	// 	err = &tErrors.CustomError{Param: "assetManagerID", Err: "error-invalid-asset-manager", ErrMessage: "Invalid Asset Manager. None specified."}
	// 	return
	// }

	// // check asset manager ID
	// am := GetAssetManagerByID(input.AssetManagerID, gc.DB)
	// if am.ID == 0 {
	// 	log.Printf("[SubmitTokenizationAssetInfoByInitiator] Error Invalid Asset Manager ID: %v\n", issuingWallet.ID)
	// 	err = &tErrors.CustomError{Param: "assetManagerID", Err: "error-invalid-asset-manager", ErrMessage: "Invalid Asset Manager."}
	// 	return
	// }
	//check if existing
	ato, NotFound, e := GetOpenTokenizedAssetByInitiatorUsername(initiator.Username, gc.DB)
	// ato, NotFound, e := GetTokenizedAssetByIssuingWallet(issuingWallet.ID, gc.DB)

	if e == nil {
		//tokenization existing
		if ato.AssetTokenizationStatus > 0 {
			// error tokenization is already in progress
			log.Printf("[SubmitTokenizationAssetInfoByInitiator] Error tokenization procesing is in progress and cannot be modified: %v\n", ato.ID)
			err = &tErrors.CustomError{Param: "issuingWalletPublicKey", Err: "error-tokenization-cannot-be-modified-by-this-method", ErrMessage: "Tokenization cannot be modified by this method."}
			return

		}
		ato = UpdateFromInput(&ato, input, gc)

		ato.LastUpdatedBy = &initiator.Username

	} else {
		if !NotFound {
			//critical database error occured
			log.Printf("[SubmitTokenizationAssetInfoByInitiator]error fetching existing tokenization from database  [%v] for %v: %v\n", input, initiator.Username, e)
			err = &tErrors.ErrorTemporaryServerError{}
			return

		}

		//create new tokenization
		ato = userModels.TokenizedAsset{
			ID:                uuid.NewString(),
			InitiatorUsername: initiator.Username,
			// IssuingWalletPublicKey: issuingWallet.ID,
			// IssuingWalletAlias:     issuingWallet.Alias,
		}
		ato = UpdateFromInput(&ato, input, gc)

	}

	e = gc.DB.Omit(clause.Associations).Save(&ato).Error
	if e != nil {
		log.Printf("[SubmitTokenizationAssetInfoByInitiator] error saving tokenization to database  [%v] for %v: %v\n", input, initiator.Username, e)

		err = &tErrors.ErrorTemporaryServerError{}

	}
	ato, _, _ = GetTokenizedAssetByID(ato.ID, gc.DB)
	return ato, err
}

// SubmitTokenizationAssetInfo used by trovoManager
func SubmitTokenizationAssetInfo(tokenizationID string, initiator *userModels.User, input *userModels.TokenizedAssetJSONInput, gc *sharedconfig.GlobalConfig) (ato userModels.TokenizedAsset, issuingWallet userModels.UserWallet, err error) {
	if len(strings.TrimSpace(os.Getenv("TOKENIZATION_ISSUING_PROFILE"))) == 0 {
		err = &tErrors.CustomError{Param: "issuingWalletPublicKey", Err: "error-default-issuing-profile-not-set", ErrMessage: "Issuing profile not set."}
		return
	}
	input.AssetCode = strings.ToUpper(strings.ReplaceAll(strings.TrimSpace(input.AssetCode), " ", ""))
	if len(input.AssetCode) == 0 {
		err = &tErrors.CustomError{Param: "assetCode", Err: "error-asset-code-not-set", ErrMessage: "Asset code not set."}
		return
	}
	ato, _, errorGetTokenizationByID := GetTokenizedAssetByID(tokenizationID, gc.DB)
	if errorGetTokenizationByID != nil {
		// error tokenization is already in progress
		log.Printf("[SubmitTokenizationAssetInfo] Error fetching  tokenization with ID: %v\n", tokenizationID)
		err = errorGetTokenizationByID
		return

	}

	//tokenization existing
	if ato.AssetTokenizationStatus < 2 {
		// error tokenization is already in progress
		log.Printf("[SubmitTokenizationAssetInfo] Error tokenization information submission is in progress and cannot be modified: %v\n", tokenizationID)
		err = &tErrors.CustomError{Param: "issuingWalletPublicKey", Err: "error-tokenization-cannot-be-modified-by-this-method", ErrMessage: "Tokenization cannot be modified by this method."}
		return

	}
	ato = UpdateFromInput(&ato, input, gc)

	ato.LastUpdatedBy = &initiator.Username
	var NotIssuedByIssuer bool
	if ato.IssuingWalletAlias != nil && len(strings.TrimSpace(os.Getenv("TOKENIZATION_ISSUING_PROFILE"))) > 0 {
		NotIssuedByIssuer = *ato.IssuingWalletAlias != strings.TrimSpace(os.Getenv("TOKENIZATION_ISSUING_PROFILE"))
	}

	if (ato.IssuingWalletPublicKey == nil || NotIssuedByIssuer) && len(input.AssetCode) > 0 && len(os.Getenv("TOKENIZATION_ISSUING_PROFILE")) > 1 && len(os.Getenv("TOKENIZATION_ISSUING_PROFILE_WALLET")) == 56 {

		//create issuing wallet

		//get atprofile
		var tokenizationIssuerProfile, tokenizationIssuerProfileWallet string
		if len(os.Getenv("TOKENIZATION_ISSUING_PROFILE")) > 1 {
			tokenizationIssuerProfile = strings.TrimSpace(os.Getenv("TOKENIZATION_ISSUING_PROFILE"))
		}
		if len(os.Getenv("TOKENIZATION_ISSUING_PROFILE_WALLET")) > 1 {
			tokenizationIssuerProfileWallet = strings.TrimSpace(os.Getenv("TOKENIZATION_ISSUING_PROFILE_WALLET"))
		}
		tokenizationIssuer, e := userModels.Username(tokenizationIssuerProfile).GetFullUser(gc.DB, gc)
		if e != nil {
			err = &tErrors.CustomError{Param: "issuingWalletPublicKey", Err: "error-invallid-issuing-profile", ErrMessage: "Issuing profile not valid."}
			return
		}
		issuer, _ := keypair.Random()
		distributor, _ := keypair.Random()
		tokenizationIssuerProfileWalletKP := keypair.MustParseFull(tokenizationIssuerProfileWallet)

		walletTag := fmt.Sprintf("%v_issuer", input.AssetCode)
		p := userModels.SubWalletInfo{
			PublicKey:             issuer.Address(),
			WalletTag:             walletTag,
			WalletType:            1,
			LinkedWalletPublicKey: distributor.Address(),
		}
		_, err = CreateNewSubWallet(&tokenizationIssuer, &p, gc)

		if err != nil {
			log.Printf("[SubmitTokenizationAssetInfo.CreateNewSubWallet: stage 1] Error creating issuing wallet [%v], err: %v\n", issuer.Address(), err)
			// err = &tErrors.CustomError{Param: "issuingPublicKey", Err: "error-invalid-issuer", ErrMessage: err}
			return
		}

		//sign transactions
		{

			if p.LinkedWalletMustSign == 1 {
				dsigned, e := middleware.SignBase64Txn(distributor.Seed(), p.Transaction, p.NetworkPassPhrase)
				if e != nil {
					log.Printf("[SubmitTokenizationAssetInfo.SignBase64Txn] Error signing issuing wallet with linked wallet [%v], err: %v\n", distributor.Address(), e)
					err = &tErrors.CustomError{Param: "issuingPublicKey", Err: "error-invalid-issuer", ErrMessage: e.Error()}
					return
				}
				p.LinkedWalletSignature = dsigned
			}

			primarySignature, subwalletSignature, e := middleware.SignSubwalletBase64Txn(tokenizationIssuerProfileWalletKP.Seed(), issuer.Seed(), p.Transaction, p.NetworkPassPhrase)
			if e != nil {
				log.Printf("[SubmitTokenizationAssetInfo.SignSubwalletBase64Txn] Error signing issuing wallet with primary and sub wallets [%v] [%v], err: %v\n", tokenizationIssuerProfileWalletKP.Address(), issuer.Address(), e)
				err = &tErrors.CustomError{Param: "issuingPublicKey", Err: "error-invalid-issuer", ErrMessage: e.Error()}
				return

			}

			p.PrimarySignature = primarySignature
			p.SubWalletSignature = subwalletSignature

		}

		//second submission to blockchain
		_, err = CreateNewSubWallet(&tokenizationIssuer, &p, gc)
		if err != nil {
			log.Printf("[SubmitTokenizationAssetInfo.CreateNewSubWallet: stage 2] Error creating issuing wallet [%v], err: %v\n", issuer.Address(), err)
			// err = &tErrors.CustomError{Param: "issuingPublicKey", Err: "error-invalid-issuer", ErrMessage: err}
			return
		}

		log.Printf("[SubmitTokenizationAssetInfo.CreateNewSubWallet] Succesfully Created issuing wallet [%v], txID: %v\n", issuer.Address(), p.TransactionID)

		w, e := userModels.WalletAlias(p.Alias).GetWallet(gc.DB, gc)
		if e != nil {
			log.Printf("[SubmitTokenizationAssetInfo] Error fetching issuing wallet [%v], err: %v\n", issuer.Address(), e)
			err = e
			return
		}

		//set issuing wallet
		issuingWallet = w
		ato.IssuingWalletAlias = &w.Alias

		ato.IssuingWalletPublicKey = &w.ID
		ato.MarketMakingWallet = w.LinkedWalletPublicKey
		ato.WalletToHoldAssetsNotForSale = w.LinkedWalletPublicKey
	}
	if ato.IssuingWalletAlias == nil {
		err = &tErrors.CustomError{Param: "issuingPublicKey", Err: "error-invalid-issuer", ErrMessage: "issuer is empty."}
		return
	}
	if *ato.IssuingWalletAlias == "" {
		err = &tErrors.CustomError{Param: "issuingPublicKey", Err: "error-invalid-issuer", ErrMessage: "issuer is empty."}
		return
	}
	if len(issuingWallet.ID) == 0 {
		//new wallet not generated. populate with existing wallet info.
		w, e := userModels.WalletAlias(*ato.IssuingWalletAlias).GetWallet(gc.DB, gc)
		if e != nil {
			log.Printf("[SubmitTokenizationAssetInfo] Error fetching issuing wallet [%v], err: %v\n", ato.IssuingWalletAlias, e)
			err = e
			return
		}

		//set issuing wallet
		issuingWallet = w
		// ato.IssuingWalletAlias = &w.Alias
		ato.IssuingWalletPublicKey = &w.ID
		ato.MarketMakingWallet = w.LinkedWalletPublicKey
		ato.WalletToHoldAssetsNotForSale = w.LinkedWalletPublicKey
		input.WalletToHoldAssetsNotForSale = *w.LinkedWalletPublicKey
	}
	// initialize message array
	input.Messages = make([]string, 0)
	// check asset manager ID
	if input.AssetManagerID == 0 {
		log.Printf("[SubmitTokenizationAssetInfo] Error Invalid Asset Manager ID: %v\n%v\n", input.AssetManagerID, tokenizationID)
		err = &tErrors.CustomError{Param: "assetManagerID", Err: "error-invalid-asset-manager", ErrMessage: "Invalid Asset Manager. None specified."}
		return
	}

	// check asset manager ID
	am := GetAssetManagerByID(input.AssetManagerID, gc.DB)
	if am.ID == 0 {
		log.Printf("[SubmitTokenizationAssetInfo] Error Invalid Asset Manager ID: %v\n%v\n", input.AssetManagerID, tokenizationID)
		err = &tErrors.CustomError{Param: "assetManagerID", Err: "error-invalid-asset-manager", ErrMessage: "Invalid Asset Manager."}
		return
	}

	e := gc.DB.Omit(clause.Associations).Save(&ato).Error
	if e != nil {
		log.Printf("[SubmitTokenizationAssetInfo] error saving tokenization to database  [%v] for %v: %v\n", input, initiator.Username, e)

		err = &tErrors.ErrorTemporaryServerError{}

	}
	ato, _, _ = GetTokenizedAssetByID(ato.ID, gc.DB)
	return ato, issuingWallet, nil
}

// VetTokenizationAssetInfo used by trovoManager
func VetTokenizationAssetInfo(tokenizationID string, initiator *userModels.User, input *userModels.VetTokenizedAssetJSONInput, gc *sharedconfig.GlobalConfig) (ato userModels.TokenizedAsset, issuingWallet userModels.UserWallet, err error) {

	ato, _, errorGetTokenizationByID := GetTokenizedAssetByID(tokenizationID, gc.DB)
	if errorGetTokenizationByID != nil {
		// error tokenization is already in progress
		log.Printf("[VetTokenizationAssetInfo] Error fetching  tokenization with ID: %v\n", tokenizationID)
		err = errorGetTokenizationByID
		return

	}

	//tokenization existing
	if ato.AssetTokenizationStatus > 1 {
		// error tokenization is already in progress
		log.Printf("[VetTokenizationAssetInfo] Error tokenization information submission has passed vetting stage and cannot be modified: %v\n", tokenizationID)
		err = &tErrors.CustomError{Param: "issuingWalletPublicKey", Err: "error-tokenization-cannot-be-modified-by-this-method", ErrMessage: "Tokenization request has passed the vetting stage."}
		return

	}
	if ato.AssetCountryLocation == nil {

		if len(input.CountryCode) == 0 {
			// error tokenization
			log.Printf("[VetTokenizationAssetInfo] Error tokenization country information not set: %v\n", tokenizationID)
			err = &tErrors.CustomError{Param: "issuingWalletPublicKey", Err: "error-tokenization-country-not-set", ErrMessage: "Tokenization country not set."}
			return
		}
		ato.AssetCountryLocation = &input.CountryCode

	}

	// initialize message array
	input.Messages = make([]string, 0)
	// check asset manager ID
	if input.AssetManagerID == 0 {
		log.Printf("[VetTokenizationAssetInfo] Error Invalid Asset Manager ID: %v\n%v\n", input.AssetManagerID, tokenizationID)
		err = &tErrors.CustomError{Param: "assetManagerID", Err: "error-invalid-asset-manager", ErrMessage: "Invalid Asset Manager. None specified."}
		return
	}

	// check asset manager ID
	am := GetAssetManagerByID(input.AssetManagerID, gc.DB)
	if am.ID == 0 {
		log.Printf("[VetTokenizationAssetInfo] Error Invalid Asset Manager ID: %v\n%v\n", input.AssetManagerID, tokenizationID)
		err = &tErrors.CustomError{Param: "assetManagerID", Err: "error-invalid-asset-manager", ErrMessage: "Invalid Asset Manager."}
		return
	}
	assetMgtConfig := userModels.CountryCode(*ato.AssetCountryLocation).GetAssetMgtFee(input.AssetManagerID, gc)
	assetCustodianConfig := userModels.CountryCode(*ato.AssetCountryLocation).GetCustodyFee(input.ApprovedAssetCustodianID, gc)

	ato.ApprovedAssetCustodianID = input.ApprovedAssetCustodianID
	ato.CustodianFeePercent = assetCustodianConfig.FeePercent
	ato.AssetManagerID = input.AssetManagerID
	ato.AssetManagerFeePercent = assetMgtConfig.FeePercent
	ato.VettingStatus = 1

	ato.LastUpdatedBy = &initiator.Username

	ato.UpdateCalculation(gc)

	e := gc.DB.Omit(clause.Associations).Save(&ato).Error
	if e != nil {
		log.Printf("[VetTokenizationAssetInfo] error saving tokenization to database  [%v] for %v: %v\n", input, initiator.Username, e)

		err = &tErrors.ErrorTemporaryServerError{}

	}
	ato, _, _ = GetTokenizedAssetByID(ato.ID, gc.DB)
	return ato, issuingWallet, nil
}

// ConfirmTokenizationAssetInfo advance status to 1 and allow for vetting.
func ConfirmTokenizationAssetInfo(initiator *userModels.User, tokenizationID string, taInput *userModels.ConfirmTokenizedAssetJSONInput, gc *sharedconfig.GlobalConfig) (ato userModels.TokenizedAsset, err error) {
	taInput.NetworkPassPhrase = gc.BantuNetworkPassphrase
	taInput.Messages = make([]string, 0)
	//check if existing
	ato, NotFound, e := GetOpenTokenizedAssetByInitiatorUsername(initiator.Username, gc.DB)

	if e == nil {
		//tokenization existing
		if ato.AssetTokenizationStatus > 0 {
			// error tokenization is already in progress
			log.Printf("[SubmitTokenizationAssetInfo] Error tokenization procesing is in progress and cannot be modified: %v\n", tokenizationID)
			err = &tErrors.CustomError{Param: "issuingWalletPublicKey", Err: "error-tokenization-cannot-be-modified-by-this-method", ErrMessage: "Tokenization is already in progress, this action cannot be performed."}
			return

		}

		if ato.ID != tokenizationID {
			err = &tErrors.CustomError{Param: "issuingWalletPublicKey", Err: "error-tokenization-id-not-valid", ErrMessage: "Invalid tokenization specified."}
			return
		}

		ato.LastUpdatedBy = &initiator.Username
		ato.AssetTokenizationStatus = 1

	} else {
		if !NotFound {
			//critical database error occured
			log.Printf("[SubmitTokenizationAssetInfo]error fetching existing tokenization from database ID [%v] for %v: %v\n", tokenizationID, initiator.Username, e)
			err = &tErrors.ErrorTemporaryServerError{}
			return

		}

		log.Printf("[SubmitTokenizationAssetInfo] Tokenization does not exist: %v\n", tokenizationID)
		err = &tErrors.CustomError{Param: "Id", Err: "error-tokenization-not-found", ErrMessage: "Only existing valid tokenization requests can be confirmed."}
		return

	}
	//begin a database transaction here
	dbTX := gc.DB.Begin()
	defer dbTX.Rollback()

	e = dbTX.Omit(clause.Associations).Save(&ato).Error
	if e != nil {
		log.Printf("[SubmitTokenizationAssetInfo] error saving tokenization to database  [%+v] for %v: %v\n", ato, initiator.Username, e)

		err = &tErrors.ErrorTemporaryServerError{}

	}

	//check for blockchain action
	wallet, e := userModels.WalletAlias(initiator.Username).GetWallet(dbTX, gc)
	if e != nil {
		log.Printf("[SubmitTokenizationAssetInfo] error getting primary wallet from database  [%+v] for %v: %v\n", ato, initiator.Username, e)

		err = &tErrors.ErrorTemporaryServerError{}

		return

	}

	xdrBase64, e := generateTokenizationFeeXdr(&wallet, taInput, gc)
	if e != nil {
		log.Printf("[SubmitTokenizationAssetInfo] error getting appliction fee transaction  [%+v] for %v: %v\n", ato, initiator.Username, e)

		err = &tErrors.ErrorTemporaryServerError{}

		return

	}

	taInput.Transaction = xdrBase64

	if len(taInput.TransactionSignature) > 0 {
		//submit to network
		txnHash, e := network.SubmitXdrWithSignature(gc.BantuExpansionClient, initiator.PrimarySigner, xdrBase64, taInput.TransactionSignature)
		if e != nil {
			err = e
			return
		}

		dbTX.Commit()
		taInput.TransactionID = txnHash
		err = nil

	}

	ato, _, _ = GetTokenizedAssetByID(ato.ID, gc.DB)
	return ato, err
}

func ConfirmTokenizationAssetPaymentInfo(initiator *userModels.User, tokenizationID string, gc *sharedconfig.GlobalConfig) (ato userModels.TokenizedAsset, err error) {

	//check if existing
	ato, NotFound, e := GetOpenTokenizedAssetByInitiatorUsername(initiator.Username, gc.DB)

	if e == nil {
		//tokenization existing
		if ato.AssetTokenizationStatus != 1 {
			// error tokenization is already in progress
			log.Printf("[SubmitTokenizationAssetInfo] Error tokenization process not awaiting payment and cannot be modified: %v\n", tokenizationID)
			err = &tErrors.CustomError{Param: "issuingWalletPublicKey", Err: "error-tokenization-cannot-be-modified-by-this-method", ErrMessage: "Tokenization is not awaiting payment, this action cannot be performed."}
			return

		}

		if ato.ID != tokenizationID {
			err = &tErrors.CustomError{Param: "issuingWalletPublicKey", Err: "error-tokenization-id-not-valid", ErrMessage: "Invalid tokenization specified."}
			return
		}

		if ato.VettingStatus == 0 {
			err = &tErrors.CustomError{Param: "vettingStatus", Err: "error-tokenization-not-vetted", ErrMessage: "Tokenization request is still being vetted by the team. Please wait until vetting has completed."}
			return
		}

		ato.LastUpdatedBy = &initiator.Username
		ato.AssetTokenizationStatus = 2

	} else {
		if !NotFound {
			//critical database error occured
			log.Printf("[ConfirmTokenizationAssetPaymentInfo]error fetching existing tokenization from database  ID [%v] for %v: %v\n", tokenizationID, initiator.Username, e)
			err = &tErrors.ErrorTemporaryServerError{}
			return

		}

		log.Printf("[ConfirmTokenizationAssetPaymentInfo] Tokenization does not exist: %v\n", tokenizationID)
		err = &tErrors.CustomError{Param: "Id", Err: "error-tokenization-not-found", ErrMessage: "Only existing valid tokenization requests can be confirmed."}
		return

	}
	// if ato.TokenizationTransaction == nil {
	// 	log.Printf("[ConfirmTokenizationAssetPaymentInfo] Tokenization transaction does not exist: %v\n", issuingWallet.ID)
	// 	err = &tErrors.CustomError{Param: "Id", Err: "error-tokenization-transaction-found", ErrMessage: "Transaction could not be generated. Please ensure all mandatory fields are supplied and try again."}
	// 	return
	// }
	e = gc.DB.Omit(clause.Associations).Save(&ato).Error
	if e != nil {
		log.Printf("[ConfirmTokenizationAssetPaymentInfo] error saving tokenization to database  [%+v] for %v: %v\n", ato, initiator.Username, e)

		err = &tErrors.ErrorTemporaryServerError{}

	}
	ato, _, _ = GetTokenizedAssetByID(ato.ID, gc.DB)
	return ato, err
}

func GetTokenizationList(user *userModels.User, gc *sharedconfig.GlobalConfig, c *gin.Context) (records userModels.PaginatedTokenizedAssets) {
	var err error
	var tokenizedAssetList []userModels.TokenizedAsset
	records.Records = make([]userModels.TokenizedAssetJSON, 0)
	DB, _ := db.OpenDb()
	DBC, _ := db.OpenDb()

	onlyWithUserPermission := strings.TrimSpace(strings.ToUpper(c.DefaultQuery("onlyWithUserPermission", "1")))

	//get all wallets where user has access
	// sharedWallets := make([]string, 0)
	// if onlyWithUserPermission == "1" {
	// 	for _, w := range user.WalletsSharedWithUser {
	// 		sharedWallets = append(sharedWallets, w.WalletPublicKey)
	// 	}
	// }

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

	if onlyWithUserPermission == "1" {
		// query = query.Where("(issuing_wallet_public_key IN (?))", sharedWallets)
		// countQuery = countQuery.Where("(issuing_wallet_public_key IN (?))", sharedWallets)
		query = query.Where("lower(initiator_username) = ?", user.Username)
		countQuery = countQuery.Where("lower(initiator_username) = ?", user.Username)
	}

	if len(initiatorUsername) > 2 {
		query = query.Where("lower(initiator_username) = ?", initiatorUsername)
		countQuery = countQuery.Where("lower(initiator_username) = ?", initiatorUsername)

	}
	if len(tokenizationStatus) > 0 && onlyWithUserPermission == "1" {
		ts, _ := strconv.ParseUint(strings.TrimSpace(tokenizationStatus), 10, 64)
		tsInt := int(ts)
		query = query.Where("asset_tokenization_status = ?", tsInt)
		countQuery = countQuery.Where("asset_tokenization_status = ?", tsInt)

	}

	//get only market ready list
	if onlyWithUserPermission == "0" {

		query = query.Where("asset_tokenization_status > ?", 3)
		countQuery = countQuery.Where("asset_tokenization_status > ?", 3)

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
	return t.UpdateFromInput(ti, gc)
}

func generateMintRegulatedTokenizedAssetXdr(t *userModels.TokenizedAsset, gc *sharedconfig.GlobalConfig) (xdrbase64, transactionSource string, messages []string, issuingWallet userModels.UserWallet, err error) {
	client := gc.BantuExpansionClient
	feeWallet := keypair.MustParse(os.Getenv("TOKENIZATION_FEE_WALLET"))
	ops := make([]txnbuild.Operation, 0)
	messages = make([]string, 0)
	var permInfo []userModels.WalletPermissionInfo

	var minBalance = decimal.NewFromFloat(3.0)

	if len(os.Getenv("WALLET_MINIMUM_BALANCE")) > 0 {
		minBalance = decimal.RequireFromString(os.Getenv("WALLET_MINIMUM_BALANCE"))
	}
	if t.AssetQuoteCurrency == nil {
		log.Println("[generateMintRegulatedTokenizedAssetXdr] Error locating quote currency")

		err = &tErrors.CustomError{
			Param:      "publicKey",
			Err:        "error-invalid-quote-currency",
			ErrMessage: "Tokenization does not have valid tokenization currency.",
			Code:       404,
		}
		return "", "", messages, issuingWallet, err
	}
	// get currency
	quoteCurrency := GetTokenizationCurrencyByCode(*t.AssetQuoteCurrency, gc.DB)

	if len(quoteCurrency.AssetIssuer) == 0 {
		log.Printf("[generateMintRegulatedTokenizedAssetXdr] Error locating quote currency %v\n", *t.AssetQuoteCurrency)

		err = &tErrors.CustomError{
			Param:      "publicKey",
			Err:        "error-invalid-quote-currency",
			ErrMessage: "Tokenization does not have valid tokenization currency.",
			Code:       404,
		}
		return "", "", messages, issuingWallet, err
	}

	var aps []userModels.TokenizationMintingApprover
	var inits []userModels.TokenizationMintingInitiator
	gc.DB.Find(&aps)
	gc.DB.Find(&inits)

	if len(aps) == 0 || len(inits) == 0 {
		err = &tErrors.CustomError{
			Param:      "numberOfApprovers",
			Err:        "error-no-approver-or-initiator-specified",
			ErrMessage: "No approvers /initiators specified",
			Code:       404,
		}
		return "", "", messages, issuingWallet, err
	}
	//get atprofile
	var tokenizationIssuerProfile, tokenizationIssuerProfileWallet string
	if len(os.Getenv("TOKENIZATION_ISSUING_PROFILE")) > 1 {
		tokenizationIssuerProfile = strings.TrimSpace(os.Getenv("TOKENIZATION_ISSUING_PROFILE"))
	}
	if len(os.Getenv("TOKENIZATION_ISSUING_PROFILE_WALLET")) > 1 {
		tokenizationIssuerProfileWallet = strings.TrimSpace(os.Getenv("TOKENIZATION_ISSUING_PROFILE_WALLET"))
	}
	var e error
	tokenizationIssuerUser, _ := userModels.Username(tokenizationIssuerProfile).GetFullUser(gc.DB, gc)
	issuingWallet, e = userModels.UserWalletID(*t.IssuingWalletPublicKey).GetWallet(gc.DB, gc)
	if e != nil {
		err = &tErrors.CustomError{Param: "issuingPublicKey", Err: "error-invalid-issuer", ErrMessage: "error validating issuing wallet."}
		return
	}
	distributionWallet, _ := userModels.UserWalletID(*issuingWallet.LinkedWalletPublicKey).GetWallet(gc.DB, gc)

	tokenizationIssuerProfileWalletKP := keypair.MustParseFull(tokenizationIssuerProfileWallet)

	for _, v := range aps {
		permInfo = append(permInfo, userModels.WalletPermissionInfo{
			TargetUsername: v.Approver,
			Permission:     "APPROVER",
		})
	}

	for _, v := range inits {
		permInfo = append(permInfo, userModels.WalletPermissionInfo{
			TargetUsername: v.Initiator,
			Permission:     "INITIATOR",
		})
	}
	var p userModels.UserWalletSharedAccessInfo

	p = userModels.UserWalletSharedAccessInfo{
		WalletPublicKey:         *t.IssuingWalletPublicKey,
		NumberOfApprovalsNeeded: 2,
		Permissions:             permInfo,
	}

	if issuingWallet.SharedAccessEnabled == 0 {

		//create sharedAccess on issuing wallet
		_, errSharedAccess := CreateSharedWalletAccess(&tokenizationIssuerUser, &tokenizationIssuerUser, &issuingWallet, &p, gc)

		if errSharedAccess != nil {
			log.Printf("[generateMintRegulatedTokenizedAssetXdr.SignBase64Txn] Error creating shared access on issuing wallets [%v] [%v], err: %v\n", tokenizationIssuerProfileWalletKP.Address(), errSharedAccess)

			err = &tErrors.CustomError{
				Param:      "IssuingWalletPublicKey",
				Err:        "error-could-not-create-shared-access-on-wallets",
				ErrMessage: "Could not create shared access on issuing wallet",
				Code:       404,
			}
			return "", "", messages, issuingWallet, err
		}

		//sign and submit
		//sign transactions
		if p.SignatureRequired == 1 {
			signedBase64, e := middleware.SignBase64Txn(tokenizationIssuerProfileWalletKP.Seed(), p.Transaction, p.NetworkPassPhrase)
			if e != nil {
				log.Printf("[generateMintRegulatedTokenizedAssetXdr.SignBase64Txn] Error signing issuing wallet with primary wallets [%v] [%v], err: %v\n", tokenizationIssuerProfileWalletKP.Address(), e)
				err = &tErrors.CustomError{Param: "issuingPublicKey", Err: "error-invalid-issuer", ErrMessage: e.Error()}
				return

			}

			p.TransactionSignature = signedBase64
		}

		//second submission to blockchain
		_, err = CreateSharedWalletAccess(&tokenizationIssuerUser, &tokenizationIssuerUser, &issuingWallet, &p, gc)
		if err != nil {
			log.Printf("[generateMintRegulatedTokenizedAssetXdr.CreateSharedWalletAccess: stage 2] Error creating issuing wallet [%v], err: %v\n", tokenizationIssuerProfileWalletKP.Address(), err)
			return
		}

		log.Printf("[generateMintRegulatedTokenizedAssetXdr.CreateSharedWalletAccess] Succesfully Created shared access on issuing wallet [%v], txID: %v\n", tokenizationIssuerProfileWalletKP.Address(), p.TransactionID)
	}
	// create distributionWallet trustline to asset
	ops = append(ops, &txnbuild.ChangeTrust{
		Line:          txnbuild.ChangeTrustAssetWrapper{Asset: txnbuild.CreditAsset{Code: *t.AssetCode, Issuer: *t.IssuingWalletPublicKey}},
		Limit:         "900000000000",
		SourceAccount: distributionWallet.ID,
	})

	//create trusline on fee wallet

	//get fee wallet
	ops = append(ops, &txnbuild.ChangeTrust{
		Line:          txnbuild.ChangeTrustAssetWrapper{Asset: txnbuild.CreditAsset{Code: *t.AssetCode, Issuer: *t.IssuingWalletPublicKey}},
		Limit:         "900000000000",
		SourceAccount: feeWallet.Address(),
	})

	// allow trust from issuer to distribution wallet
	ops = append(ops, &txnbuild.SetTrustLineFlags{
		Trustor:       distributionWallet.ID,
		Asset:         txnbuild.CreditAsset{Code: *t.AssetCode, Issuer: *t.IssuingWalletPublicKey},
		SetFlags:      []txnbuild.TrustLineFlag{txnbuild.TrustLineAuthorized, txnbuild.TrustLineClawbackEnabled},
		SourceAccount: *t.IssuingWalletPublicKey,
	})
	// allow trust from issuer to distribution wallet
	ops = append(ops, &txnbuild.SetTrustLineFlags{
		Trustor:       feeWallet.Address(),
		Asset:         txnbuild.CreditAsset{Code: *t.AssetCode, Issuer: *t.IssuingWalletPublicKey},
		SetFlags:      []txnbuild.TrustLineFlag{txnbuild.TrustLineAuthorized, txnbuild.TrustLineClawbackEnabled},
		SourceAccount: *t.IssuingWalletPublicKey,
	})

	//mint the token to distribution wallet
	ops = append(ops, &txnbuild.Payment{
		Destination:   distributionWallet.ID,
		Amount:        decimal.NewFromFloat(t.NumberOfTokenToBeIssued).StringFixed(7),
		Asset:         txnbuild.CreditAsset{Code: *t.AssetCode, Issuer: *t.IssuingWalletPublicKey},
		SourceAccount: *t.IssuingWalletPublicKey,
	})

	//deduct fee to fee wallet, from distribution wallet
	ops = append(ops, &txnbuild.Payment{
		Destination:   feeWallet.Address(),
		Amount:        decimal.NewFromFloat(t.FeeInAsset).StringFixed(7),
		Asset:         txnbuild.CreditAsset{Code: *t.AssetCode, Issuer: *t.IssuingWalletPublicKey},
		SourceAccount: distributionWallet.ID,
	})

	//make market
	fraction := decimal.NewFromFloat(t.PricePerToken).Rat()
	d := int32(fraction.Denom().Int64())
	n := int32(fraction.Num().Int64())
	ops = append(ops, &txnbuild.ManageSellOffer{
		Buying:        txnbuild.CreditAsset{Code: quoteCurrency.AssetCode, Issuer: quoteCurrency.AssetIssuer},
		Amount:        decimal.NewFromFloat(t.NumberOfTokenToBeSold).StringFixed(7),
		Selling:       txnbuild.CreditAsset{Code: *t.AssetCode, Issuer: *t.IssuingWalletPublicKey},
		Price:         xdr.Price{N: xdr.Int32(n), D: xdr.Int32(d)},
		SourceAccount: distributionWallet.ID,
	})

	//

	//check if issuing account has native enough native balance
	var nativeAsset txnbuild.Asset = txnbuild.NativeAsset{}
	_, _, walletAccountNativeBalance, _, _, errWalletAct := network.BlockchainAccountProperties(client, issuingWallet.ID, nativeAsset)
	if errWalletAct != nil {
		log.Printf("[generateMintRegulatedTokenizedAssetXdr] by [%v] for shared Account Properties error:[%v] \n", issuingWallet.Alias, errWalletAct)

		return "", "", messages, issuingWallet, errWalletAct
	}
	if (walletAccountNativeBalance).LessThan(minBalance) {
		log.Printf("[generateMintRegulatedTokenizedAssetXdr] by [%v] shared WalletAccount underfunded \n", issuingWallet.Alias)

		err = &tErrors.CustomError{
			Param:      "publicKey",
			Err:        "error-wallet-underfunded",
			ErrMessage: fmt.Sprintf("Wallet %v does not have enough %v balance to perform this operation", issuingWallet.Alias, os.Getenv("NATIVE_ASSET_CODE")),
			Code:       404,
		}
		return "", "", messages, issuingWallet, err
	}

	// minting is free. No fee.
	chanAccount := <-gc.ChannelAccounts
	defer func(c *keypair.Full) {
		gc.ChannelAccounts <- c
	}(chanAccount)
	_, _, _, _, chanSourceAccount, _ := network.BlockchainAccountProperties(client, chanAccount.Address(), txnbuild.NativeAsset{})

	transactionSource = chanSourceAccount.AccountID

	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        chanSourceAccount,
			IncrementSequenceNum: true,
			Operations:           ops,
			BaseFee:              txnbuild.MinBaseFee,
			Preconditions: txnbuild.Preconditions{
				TimeBounds: txnbuild.NewInfiniteTimeout(),
			},
			Memo: txnbuild.MemoText("Mint " + *t.AssetCode),
		},
	)
	if err != nil {
		log.Println("[generateMintRegulatedTokenizedAssetXdr] error constructing transaction ", err)
		return "", transactionSource, messages, issuingWallet, err
	}

	tx, err = tx.Sign(network.GetBlockchainNetworkPassPhrase(), chanAccount)

	if err != nil {
		log.Println("[generateMintRegulatedTokenizedAssetXdr] error signing transaction with channelAccount key ", err)
		return "", transactionSource, messages, issuingWallet, &tErrors.ErrorTemporaryServerError{}
	}
	var xdrBase64 string

	xdrBase64, err = tx.Base64()
	if err != nil {
		log.Println("[generateMintRegulatedTokenizedAssetXdr] error getting txn base64", err)
		return "", transactionSource, messages, issuingWallet, err
	}

	t.TokenizationTransaction = &xdrBase64

	return xdrBase64, transactionSource, messages, issuingWallet, nil

}

// MintRegulatedTokenizedAsset mint tokenized assets
func MintRegulatedTokenizedAsset(tokenizationID string, initiator *userModels.User, gc *sharedconfig.GlobalConfig) (ato userModels.TokenizedAsset, err error) {

	ato, _, err = GetTokenizedAssetByID(tokenizationID, gc.DB)
	if err != nil {
		// error tokenization is already in progress
		log.Printf("[MintRegulatedTokenizedAsset] Error fetching tokenization with ID: %v\n", tokenizationID)
		return

	}
	{
		//check staff access
		if !IsTokenizationMintingApprover(initiator.Username, gc.DB) && !IsTokenizationMintingInitiator(initiator.Username, gc.DB) {
			log.Printf("[MintRegulatedTokenizedAsset] Error %v not minting inititor/approver.\n", initiator.Username)
			err = &tErrors.CustomError{
				Param:      "publicKey",
				Err:        "error-access-denied",
				ErrMessage: fmt.Sprintf("%v does not have minting permission to perform this action.", initiator.Username),
				Code:       401,
			}
			return
		}
	}
	//validators
	{

	}

	xdrBase64, transactionSource, messages, issuingWallet, err := generateMintRegulatedTokenizedAssetXdr(&ato, gc)

	if err != nil {

		return
	}
	description := fmt.Sprintf("Minting tokenized asset %v:%v...%v", *ato.AssetCode, issuingWallet.ID[0:4], issuingWallet.ID[51:55])
	mintObj := userModels.TokenMinting{
		TokenizedAssetID:     ato.ID,
		Destination:          *issuingWallet.LinkedWalletPublicKey,
		Amount:               decimal.NewFromFloat(ato.NumberOfTokenToBeIssued).StringFixed(7),
		AssetCode:            *ato.AssetCode,
		AssetIssuer:          issuingWallet.ID,
		Transaction:          xdrBase64,
		NetworkPassPhrase:    gc.BantuNetworkPassphrase,
		TransactionSignature: transactionSource,
		ReturnedDescription:  description,
		Commit:               1,
		Messages:             messages,
	}
	id := uuid.NewString()
	transactionByte, _ := json.Marshal(mintObj)
	transactionStr := string(transactionByte)
	pendingAuth := userModels.PendingAuth{
		ID:                       id,
		Initiator:                initiator.Username,
		InitiatorSignerPublicKey: initiator.PrimarySigner,
		WalletPublicKey:          issuingWallet.ID,
		TransactionType:          "TOKENIZE ASSET",
		Description:              description,
		TransactionSource:        transactionSource,
		ApprovalsNeeded:          issuingWallet.NumberOfApprovalsNeeded,
		TransactionXdr:           xdrBase64,
		TransactionInfoStr:       &transactionStr,
	}
	//save and commit this to database
	e := gc.DB.Create(&pendingAuth).Error
	if e != nil {
		log.Printf("[ClaimPendingAsset] Error saving txn [%+v] on pending auth table: %s\n", pendingAuth, e.Error())
		err = &tErrors.ErrorTemporaryServerError{}
		return
	}

	return ato, nil

}

func generateTokenizationFeeXdr(wallet *userModels.UserWallet, taInput *userModels.ConfirmTokenizedAssetJSONInput, gc *sharedconfig.GlobalConfig) (txnBase64 string, err error) {
	tfa := strings.Split(os.Getenv("TOKENIZATION_APPLICATION_FEE_ASSET"), ":") //CODE:ISSUER
	feeAmount := decimal.RequireFromString(os.Getenv("TOKENIZATION_APPLICATION_FEE_AMOUNT"))
	feeWallet := os.Getenv("TOKENIZATION_APPLICATION_FEE_WALLET")
	feeWalletPK := keypair.MustParseFull(feeWallet)
	assetIssuer := tfa[1]
	assetCode := tfa[0]

	asset := txnbuild.CreditAsset{Code: assetCode, Issuer: assetIssuer}

	sourceAccountExists, _, _, assetAccountBalance, sourceAccount, _ := network.BlockchainAccountProperties(gc.BantuExpansionClient, wallet.ID, asset)

	if !sourceAccountExists {
		return "", &tErrors.CustomError{Param: "publicKey", Err: "error-account-not-activated-on-blockchain", ErrMessage: "The Wallet public key is currently underfunded. Please send about 3XBN to it to activate it before you can perform this task", Code: http.StatusBadRequest}

	}
	if assetAccountBalance.LessThan(feeAmount) {
		return "", &tErrors.CustomError{Param: "publicKey", Err: "error-wallet-underfunded", ErrMessage: fmt.Sprintf("The Wallet is currently underfunded. Please maintain min %v %v balance before you can perform this task", feeAmount.String(), tfa[0]), Code: http.StatusBadRequest}

	}

	var ops []txnbuild.Operation = make([]txnbuild.Operation, 0)
	ops = append(ops, &txnbuild.Payment{
		Destination:   feeWalletPK.Address(),
		Amount:        feeAmount.String(),
		Asset:         asset,
		SourceAccount: wallet.ID,
	})

	taInput.Messages = append(taInput.Messages, fmt.Sprintf("Application fee of %v %x will be charged to your wallet with alias [%v].", feeAmount.String(), assetCode, wallet.Alias))
	var tx *txnbuild.Transaction
	// Construct the transaction that holds the operations to execute on the network

	tx, err = txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        sourceAccount,
			IncrementSequenceNum: true,
			Operations:           ops,
			BaseFee:              txnbuild.MinBaseFee,
			Preconditions: txnbuild.Preconditions{
				TimeBounds: txnbuild.NewInfiniteTimeout(),
			},
			Memo: txnbuild.MemoText("TApplicationFee"),
		},
	)

	if err != nil {
		log.Println("[generateTokenizationFeeXdr]error constructing transaction", err)
		return "", &tErrors.ErrorTemporaryServerError{}
	}

	xdrBase64, err := tx.Base64()

	if err != nil {
		return "", err
	}

	return xdrBase64, nil

}
