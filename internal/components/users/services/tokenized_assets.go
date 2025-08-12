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
	"time"
	swapErrors "trovo-wallet-api/internal/components/swaps/errors"
	swapModels "trovo-wallet-api/internal/components/swaps/models"
	swapServices "trovo-wallet-api/internal/components/swaps/services"
	db "trovo-wallet-api/internal/db"

	userModels "trovo-wallet-api/internal/components/users/models"
	// db "trovo-wallet-api/internal/db"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/middleware"
	"trovo-wallet-api/internal/network"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/ecnepsnai/discord"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/protocols/horizon"
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

func GetCountryConfigs(db *gorm.DB) (countries []userModels.Country) {
	countries = make([]userModels.Country, 0)
	db.Preload(clause.Associations).Find(&countries)

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

func GetMintingInitiators(db *gorm.DB) (us []userModels.TokenizationMintingInitiator) {
	us = make([]userModels.TokenizationMintingInitiator, 0)

	db.Find(&us)
	return
}

func GetMintingApprovers(db *gorm.DB) (us []userModels.TokenizationMintingApprover) {
	us = make([]userModels.TokenizationMintingApprover, 0)

	db.Find(&us)
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

func GetAssetIssuingHouses(db *gorm.DB) (assetIssuingHouses []userModels.AssetIssuingHouse) {
	assetIssuingHouses = make([]userModels.AssetIssuingHouse, 0)
	db.Order("asset_issuing_House_country, Asset_Issuing_House_name").Find(&assetIssuingHouses)

	return
}

func GetLegalAndProfesionalPartners(db *gorm.DB) (lpp []userModels.LegalAndProfesionalPartner) {
	lpp = make([]userModels.LegalAndProfesionalPartner, 0)
	db.Order("Partner_Country_country, Partner_name").Find(&lpp)

	return
}

func GetRatingAgencies(db *gorm.DB) (ras []userModels.RatingAgency) {
	ras = make([]userModels.RatingAgency, 0)
	db.Order("Agency_Country, Agency_name").Find(&ras)

	return
}

func GetTrustees(db *gorm.DB) (trs []userModels.Trustee) {
	trs = make([]userModels.Trustee, 0)
	db.Order("Trustee_Country, Trustee_name").Find(&trs)

	return
}

func GetAssetManagerByID(id uint64, db *gorm.DB) (assetManager userModels.AssetManager) {
	db.Where("id = ?", id).First(&assetManager)

	return
}

func GetApprovedCustodianByID(id uint64, db *gorm.DB) (custodian userModels.ApprovedAssetCustodian) {
	db.Where("id = ?", id).First(&custodian)

	return
}

func GetAssetIssuingHouseByID(id uint64, db *gorm.DB) (assetIssuingHouse userModels.AssetIssuingHouse) {
	db.Where("id = ?", id).First(&assetIssuingHouse)

	return
}

func GetLegalAndProfesionalPartnerByID(id uint64, db *gorm.DB) (lpp userModels.LegalAndProfesionalPartner) {
	db.Where("id = ?", id).First(&lpp)

	return
}

func GetRatingAgencyByID(id uint64, db *gorm.DB) (ra userModels.RatingAgency) {
	db.Where("id = ?", id).First(&ra)

	return
}

func GetTrusteeByID(id uint64, db *gorm.DB) (tr userModels.Trustee) {
	db.Where("id = ?", id).First(&tr)

	return
}

func GetTokenizationFees(db *gorm.DB) (fees []userModels.TokenizationFee) {
	fees = make([]userModels.TokenizationFee, 0)
	// db.Preload(clause.Associations).Where("inactive != ?", 1).Find(&fees)
	db.Order("id ASC, country_code ASC").Where("inactive != ?", 1).Find(&fees)

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
	// err = db.Preload(clause.Associations).Order("asset_tokenization_status ASC").Where("asset_tokenization_status < 1 AND initiator_username = ?", initiatorUsername).First(&tokenizedAsset).Error

	// if err != nil {
	// 	if !errors.Is(err, gorm.ErrRecordNotFound) {
	// 		//critical database error occured
	// 		log.Printf("[GetTokenizedAssetByID]error fetching existing tokenization with initiatorUsername %v from database  [%v]", initiatorUsername, err)
	// 		return

	// 	} else {
	// 		//record not found
	// 		NotFound = true
	// 		err = &tErrors.CustomError{Param: "tokenizationID", Err: "error-invalid-tokenizationId", ErrMessage: fmt.Sprintf("%v has no tokenized asset inititated", initiatorUsername)}
	// 		return
	// 	}
	// }
	return userModels.Username(initiatorUsername).GetOpenTokenizedAssetByInitiatorUsername(db)

}

// GetFeeReadyTokenizedAssetApplicationByInitiatorUsername get the tokenization that has status 1 and initiated by the initiator username
func GetFeeReadyTokenizedAssetApplicationByInitiatorUsername(initiatorUsername string, db *gorm.DB) (tokenizedAsset userModels.TokenizedAsset, NotFound bool, err error) {

	return userModels.Username(initiatorUsername).GetFeeReadyTokenizedAssetApplicationByInitiatorUsername(db)
}

// GetOpenTokenizedAssetByInitiatorUsername get the tokenization that has status 0 or 1 initiated by the initiator username
func GetOpenTokenizedAssetByInitiatorUsernameAndID(initiatorUsername, tokenizedAssetID string, db *gorm.DB) (tokenizedAsset userModels.TokenizedAsset, NotFound bool, err error) {
	// var ta userModels.TokenizedAsset
	err = db.Preload(clause.Associations).Where("asset_tokenization_status < 2 AND initiator_username = ? AND id = ?", initiatorUsername, tokenizedAssetID).First(&tokenizedAsset).Error

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
		es := gc.DB.Omit(clause.Associations).Save(&documentUpload).Error
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
		e := gc.DB.Omit(clause.Associations).Create(&documentUpload).Error
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
		log.Printf("[UploadTokenizationFeeProofOfPaymentDocument]error uploading proof of payment document for %v: %v\n", user.Username, err)

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
	e := gc.DB.Omit(clause.Associations).Create(&documentUpload).Error
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
	ato, NotFound, err := GetTokenizedAssetByID(tokenizationID, gc.DB)
	if NotFound {
		log.Printf("[DeleteTokenization] error locating existing tokenization in database  [%v] for %v: %v\n", tokenizationID, user.Username, err)
		return tokenizedAsset, fmt.Errorf("error locating tokenization request with ID %v", tokenizationID)
	}

	if err != nil && !NotFound {
		log.Printf("[DeleteTokenization] error locating existing tokenization in database  [%v] for %v: %v\n", tokenizationID, user.Username, err)
		return tokenizedAsset, fmt.Errorf("system error occured while getching tokenization request with ID %v", tokenizationID)

	}

	if ato.AssetTokenizationStatus > 0 {
		err = &tErrors.CustomError{Param: "ID", Err: "error-invalid-cannot-delete", ErrMessage: "You cannot delete this tokenization request because it has passed the editing stage."}

		return ato.ToJSON(gc), err

	}
	//get access permission and see if the user taking action has an INITIATOR access.
	if ato.InitiatorUsername != user.Username {

		err = &tErrors.CustomError{Param: "ID", Err: "error-invalid-access", ErrMessage: "You cannot delete this tokenization request because you did not initiate it."}

		return tokenizedAsset, err

	}

	tx := gc.DB.Begin()
	defer tx.Rollback()

	if len(ato.ProofOfPaymentDocuments) > 0 {
		tx.Delete(&ato.ProofOfPaymentDocuments)
		// if e != nil {
		// 	log.Printf("[DeleteTokenization]error deleting existing tokenization documents in database  [%v] for %v: %v\n", tokenizationID, user.Username, e)
		// 	return ato.ToJSON(gc), fmt.Errorf("error deleting tokenization request with ID %v", tokenizationID)

		// }
	}

	if len(ato.AssetTokenizationDocuments) > 0 {
		tx.Delete(&ato.AssetTokenizationDocuments)
		// if e != nil {
		// 	log.Printf("[DeleteTokenization]error deleting existing tokenization documents in database  [%v] for %v: %v\n", tokenizationID, user.Username, e)
		// 	return ato.ToJSON(gc), fmt.Errorf("error deleting tokenization request with ID %v", tokenizationID)

		// }
	}

	//reset the document since it is purged
	ato.AssetTokenizationDocuments = nil
	ato.ProofOfPaymentDocuments = nil

	// remove wallets
	if ato.MarketMakingWallet != nil {
		distributionWallet, errDistributionWallet := userModels.UserWalletID(*ato.MarketMakingWallet).GetWallet(gc.DB, gc)
		if errDistributionWallet == nil {
			tx.Delete(&distributionWallet)
		}
	}
	// remove issuing wallet
	if ato.IssuingWalletPublicKey != nil {
		issuingWallet, errIssuingWallet := userModels.UserWalletID(*ato.IssuingWalletPublicKey).GetWallet(gc.DB, gc)
		if errIssuingWallet == nil {
			tx.Delete(&issuingWallet)
		}
	}

	e := tx.Delete(&ato).Error
	if e != nil {
		log.Printf("[DeleteTokenization]error deleting existing tokenization in database  [%v] for %v: %v\n", tokenizationID, user.Username, e)
		return ato.ToJSON(gc), fmt.Errorf("error deleting tokenization request with ID %v", tokenizationID)

	}

	//

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

func TrovoManagerDeleteTokenization(user *userModels.User, tokenizationID string, gc *sharedconfig.GlobalConfig) (tokenizedAsset userModels.TokenizedAssetJSON, err error) {
	// var document userModels.AssetTokenizationDocument
	ato, NotFound, err := GetTokenizedAssetByID(tokenizationID, gc.DB)
	if NotFound {
		log.Printf("[DeleteTokenization] error locating existing tokenization in database  [%v] for %v: %v\n", tokenizationID, user.Username, err)
		return tokenizedAsset, fmt.Errorf("error locating tokenization request with ID %v", tokenizationID)
	}

	if err != nil && !NotFound {
		log.Printf("[DeleteTokenization] error locating existing tokenization in database  [%v] for %v: %v\n", tokenizationID, user.Username, err)
		return tokenizedAsset, fmt.Errorf("system error occured while getching tokenization request with ID %v", tokenizationID)

	}

	if ato.AssetTokenizationStatus > 0 {
		err = &tErrors.CustomError{Param: "ID", Err: "error-invalid-cannot-delete", ErrMessage: "You cannot delete this tokenization request because it has passed the editing stage."}

		return ato.ToJSON(gc), err

	}

	tx := gc.DB.Begin()
	defer tx.Rollback()

	if len(ato.ProofOfPaymentDocuments) > 0 {
		tx.Delete(&ato.ProofOfPaymentDocuments)
		// if e != nil {
		// 	log.Printf("[DeleteTokenization]error deleting existing tokenization documents in database  [%v] for %v: %v\n", tokenizationID, user.Username, e)
		// 	return ato.ToJSON(gc), fmt.Errorf("error deleting tokenization request with ID %v", tokenizationID)

		// }
	}

	if len(ato.AssetTokenizationDocuments) > 0 {
		tx.Delete(&ato.AssetTokenizationDocuments)
		// if e != nil {
		// 	log.Printf("[DeleteTokenization]error deleting existing tokenization documents in database  [%v] for %v: %v\n", tokenizationID, user.Username, e)
		// 	return ato.ToJSON(gc), fmt.Errorf("error deleting tokenization request with ID %v", tokenizationID)

		// }
	}

	//reset the document since it is purged
	ato.AssetTokenizationDocuments = nil
	ato.ProofOfPaymentDocuments = nil

	// remove wallets
	if ato.MarketMakingWallet != nil {
		distributionWallet, errDistributionWallet := userModels.UserWalletID(*ato.MarketMakingWallet).GetWallet(gc.DB, gc)
		if errDistributionWallet == nil {
			tx.Delete(&distributionWallet)
		}
	}
	// remove issuing wallet
	if ato.IssuingWalletPublicKey != nil {
		issuingWallet, errIssuingWallet := userModels.UserWalletID(*ato.IssuingWalletPublicKey).GetWallet(gc.DB, gc)
		if errIssuingWallet == nil {
			tx.Delete(&issuingWallet)
		}
	}

	e := tx.Delete(&ato).Error
	if e != nil {
		log.Printf("[DeleteTokenization]error deleting existing tokenization in database  [%v] for %v: %v\n", tokenizationID, user.Username, e)
		return ato.ToJSON(gc), fmt.Errorf("error deleting tokenization request with ID %v", tokenizationID)

	}

	//

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

	e := gc.DB.Omit(clause.Associations).Delete(&document).Error
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

	e := gc.DB.Omit(clause.Associations).Delete(&document).Error
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

// SubmitTokenizationAssetInfoByInitiator for tokenization application request by initiator
func SubmitTokenizationAssetInfoByInitiator(initiator *userModels.User, input *userModels.TokenizedAssetJSONInput, gc *sharedconfig.GlobalConfig) (ato userModels.TokenizedAsset, err error) {

	// initialize message array
	input.Messages = make([]string, 0)

	if initiator.KYCVerified == 0 {
		err = &tErrors.CustomError{Param: "initiatorUsername", Err: "error-invalid-kyc-status", ErrMessage: "You have not met the minimum KYC requirement."}
		return
	}

	if len(GetTokenizationCurrencyByCode(input.AssetQuoteCurrency, gc.DB).AssetCode) == 0 {
		err = &tErrors.CustomError{Param: "assetQuoteCurrency", Err: "error-invalid-asset-quote-currency", ErrMessage: "Asset quote currency you supplied is invalid."}
		return
	}

	if decimal.NewFromFloat(input.NumberOfTokenToBeIssued).GreaterThan(gc.TokenLimitAsDecimal()) {
		err = &tErrors.CustomError{Param: "numberOfTokenToBeIssued", Err: "error-token-limit-exceeded", ErrMessage: fmt.Sprintf("Issued Number of tokens cannot exceed %v", gc.TokenLimitAsString())}
		return
	}

	input.ProceedPayoutCurrency = strings.ToUpper(input.ProceedPayoutCurrency)
	if len(input.ProceedPayoutCurrency) == 0 {
		input.ProceedPayoutCurrency = strings.ToUpper(input.AssetQuoteCurrency)
	}
	if len(GetTokenizationCurrencyByCode(input.ProceedPayoutCurrency, gc.DB).AssetCode) == 0 {
		err = &tErrors.CustomError{Param: "proceedPayoutCurrency", Err: "error-invalid-proceed-payout-currency", ErrMessage: "Proceed Payout currency code you supplied is invalid."}
		return
	}

	//check if existing
	ato, NotFound, e := GetOpenTokenizedAssetByInitiatorUsername(initiator.Username, gc.DB)

	if e == nil {
		//tokenization existing
		if ato.AssetTokenizationStatus > 0 {
			// error tokenization is already in progress
			log.Printf("[SubmitTokenizationAssetInfoByInitiator] Error tokenization procesing is in progress and cannot be modified: %v\n", ato.ID)
			err = &tErrors.CustomError{Param: "issuingWalletPublicKey", Err: "error-tokenization-cannot-be-modified-by-this-method", ErrMessage: "Tokenization cannot be modified by this method."}
			return

		}
		ato = UpdateTokenizedAssetFromInput(&ato, input, gc)

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
		}
		ato = UpdateTokenizedAssetFromInput(&ato, input, gc)

	}

	if ato.AssetCountryLocation == nil {
		log.Println("[SubmitTokenizationAssetInfoByInitiator]error Localtion/country not provided")
		err = &tErrors.CustomError{Err: "error country not provided", ErrMessage: "country of location not provided."}
		return
	}

	if len(*ato.AssetCountryLocation) != 2 {
		log.Printf("[SubmitTokenizationAssetInfoByInitiator]error Location/country not in acceptable format.")
		err = &tErrors.CustomError{Err: "error country not provided", ErrMessage: "country of location not provided in corect format. Expects 2-character formart, eg. NG"}
		return
	}
	//check Trov balance
	countryConfig := userModels.CountryCode(*ato.AssetCountryLocation).GetConfig(gc)
	_, _, _, sourceAccountCustomBalance, _, errCheckBalance := network.BlockchainAccountProperties(gc.BantuExpansionClient, initiator.PublicKey, txnbuild.CreditAsset{Code: "TROV", Issuer: "GAXMBPVA2GNG6A3NV6Q664VZASMROS5ZACKSMTPVCRIKPOJIV43A2CTJ"})
	if errCheckBalance != nil {
		log.Printf("[SubmitTokenizationAssetInfoByInitiator]error checking wallet balance for initiator. error: %v", errCheckBalance)

		err = &tErrors.CustomError{Err: "error-not-enough-trov", ErrMessage: "Error checking TROV balance of primary wallet."}
		return
	}
	if sourceAccountCustomBalance.InexactFloat64() < countryConfig.MinTROVBalanceForTokenizationApplication {

		log.Println("[SubmitTokenizationAssetInfoByInitiator]error Not enough TROV balance to initiate operation.")

		err = &tErrors.CustomError{Err: "error-not-enough-trov", ErrMessage: fmt.Sprintf("%v TROV is required in your primary wallet with alias [%v] to initiate tokenization application. You have %v TROV.", countryConfig.MinTROVBalanceForTokenizationApplication, initiator.Username, sourceAccountCustomBalance.String())}
		return
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
	replacer := strings.NewReplacer("\r", "", "\n", "", " ", "")

	if decimal.NewFromFloat(input.NumberOfTokenToBeIssued).GreaterThan(gc.TokenLimitAsDecimal()) {
		err = &tErrors.CustomError{Param: "numberOfTokenToBeIssued", Err: "error-token-limit-exceeded", ErrMessage: fmt.Sprintf("Issued Number of tokens cannot exceed %v", gc.TokenLimitAsString())}
		return
	}

	if len(strings.TrimSpace(os.Getenv("TOKENIZATION_ISSUING_PROFILE"))) == 0 {
		err = &tErrors.CustomError{Param: "issuingWalletPublicKey", Err: "error-default-issuing-profile-not-set", ErrMessage: "Issuing profile not set."}
		return
	}
	input.AssetQuoteCurrency = strings.ToUpper(input.AssetQuoteCurrency)
	ac := GetTokenizationCurrencyByCode(input.AssetQuoteCurrency, gc.DB)
	if len(ac.AssetCode) == 0 {
		log.Printf("\n\n[SubmitTokenizationAssetInfo] error: quote currency is invalid: %v, received object: %+v\n\n", ac.AssetCode, input)

		err = &tErrors.CustomError{Param: "assetQuoteCurrency", Err: "error-invalid-asset-quote-currency", ErrMessage: fmt.Sprintf("Asset quote currency [%v] you supplied is invalid.", input.AssetQuoteCurrency)}
		return
	}

	input.ProceedPayoutCurrency = strings.ToUpper(input.ProceedPayoutCurrency)
	if len(input.ProceedPayoutCurrency) == 0 {
		input.ProceedPayoutCurrency = strings.ToUpper(input.AssetQuoteCurrency)
	}
	tc := GetTokenizationCurrencyByCode(input.ProceedPayoutCurrency, gc.DB)
	if len(tc.AssetCode) == 0 {
		log.Printf("\n\n[SubmitTokenizationAssetInfo] error: payout currency is invalid: %v, received object: %+v\n\n", tc.AssetCode, input)
		err = &tErrors.CustomError{Param: "proceedPayoutCurrency", Err: "error-invalid-proceed-payout-currency", ErrMessage: fmt.Sprintf("Proceed Payout currency code [%v] you supplied is invalid.", input.ProceedPayoutCurrency)}
		return
	}
	//referesh issuing wallet profile
	userModels.Username(os.Getenv("TOKENIZATION_ISSUING_PROFILE")).InvalidateUserCache(gc)

	input.AssetCode = strings.ToUpper(replacer.Replace(strings.TrimSpace(input.AssetCode)))
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

	//validate submitted initiators/approvers
	{
		if len(input.MintingInitators) > 0 {

			input.MintingInitators = replacer.Replace(input.MintingInitators)
			us := strings.Split(input.MintingInitators, ",")
			for _, v := range us {
				_, e := userModels.Username(v).GetSimpleUser(gc.DB, gc)
				if e != nil {
					err = &tErrors.CustomError{Param: "mintingInitiators", Err: "error-invalid-minting-initiator", ErrMessage: fmt.Sprintf("%v is an invalid username for minting initiator", v)}
					return
				}
			}
		}
		if len(input.MintingApprovers) > 0 {
			input.MintingApprovers = replacer.Replace(input.MintingApprovers)
			us := strings.Split(input.MintingApprovers, ",")
			for _, v := range us {
				_, e := userModels.Username(v).GetSimpleUser(gc.DB, gc)
				if e != nil {
					err = &tErrors.CustomError{Param: "mintingApprovers", Err: "error-invalid-minting-approver", ErrMessage: fmt.Sprintf("%v is an invalid username for minting Approver", v)}
					return
				}
			}
		}
	}

	if ato.AssetCountryLocation == nil {
		err = &tErrors.CustomError{Param: "Id", Err: "error-asset-location-country-not-found", ErrMessage: "You must specify the asset country of location."}
		return
	}
	if len(*ato.AssetCountryLocation) != 2 {
		err = &tErrors.CustomError{Param: "Id", Err: "error-asset-location-country-not-found", ErrMessage: "You must specify the asset country of location in the formart: NG, SA, UK"}
		return
	}

	if ato.ExemptedCountries != nil {
		exc := strings.Split(replacer.Replace(*ato.ExemptedCountries), ",")
		for _, countryCode := range exc {
			if len(countryCode) != 2 {
				err = &tErrors.CustomError{Param: "Id", Err: "error-asset-exempted-country-invalid", ErrMessage: "You must specify the exempted country in the formart: NG, SA, UK"}
				return
			}
		}

	}

	ato = UpdateTokenizedAssetFromInput(&ato, input, gc)

	ato.LastUpdatedBy = &initiator.Username
	var NotIssuedByIssuer bool
	if ato.IssuingWalletAlias != nil && len(strings.TrimSpace(os.Getenv("TOKENIZATION_ISSUING_PROFILE"))) > 0 {
		NotIssuedByIssuer = !strings.HasPrefix(*ato.IssuingWalletAlias, strings.TrimSpace(os.Getenv("TOKENIZATION_ISSUING_PROFILE")))
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
			log.Printf("[SubmitTokenizationAssetInfo.CreateNewSubWallet: stage 1] Error creating issuing wallet [%v], err: %v\n", p.PublicKey, err)
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
			log.Printf("[SubmitTokenizationAssetInfo.CreateNewSubWallet: stage 2] Error creating issuing wallet [%v], err: %v\n", p.PublicKey, err)
			// err = &tErrors.CustomError{Param: "issuingPublicKey", Err: "error-invalid-issuer", ErrMessage: err}
			return
		}

		log.Printf("[SubmitTokenizationAssetInfo.CreateNewSubWallet] Succesfully Created issuing wallet [%v], txID: %v\n", p.PublicKey, p.TransactionID)

		w, e := userModels.UserWalletID(p.PublicKey).GetWallet(gc.DB, gc)
		if e != nil {
			log.Printf("[SubmitTokenizationAssetInfo] Error fetching issuing wallet [%v], err: %v\n", p.PublicKey, e)
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
		w, e := userModels.UserWalletID(*ato.IssuingWalletPublicKey).GetWallet(gc.DB, gc)
		if e != nil {
			log.Printf("[SubmitTokenizationAssetInfo] Error fetching issuing wallet [%v], err: %v\n", ato.IssuingWalletPublicKey, e)
			err = e
			return
		}

		//set issuing wallet
		issuingWallet = w
		ato.IssuingWalletAlias = &w.Alias
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
func VetTokenizationAssetInfo(tokenizationID string, initiator *userModels.User, input *userModels.VetTokenizedAssetJSONInput, gc *sharedconfig.GlobalConfig) (ato userModels.TokenizedAsset, err error) {

	ato, _, errorGetTokenizationByID := GetTokenizedAssetByID(tokenizationID, gc.DB)
	if errorGetTokenizationByID != nil {
		// error tokenization is already in progress
		log.Printf("[VetTokenizationAssetInfo] Error fetching  tokenization with ID: %v\n", tokenizationID)
		err = errorGetTokenizationByID
		return

	}

	if decimal.NewFromFloat(ato.NumberOfTokenToBeIssued).GreaterThan(gc.TokenLimitAsDecimal()) {
		err = &tErrors.CustomError{Param: "numberOfTokenToBeIssued", Err: "error-token-limit-exceeded", ErrMessage: fmt.Sprintf("Issued Number of tokens cannot exceed %v", gc.TokenLimitAsString())}
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

	// check asset manager ID
	// if input.AssetManagerID == 0 {
	// 	log.Printf("[VetTokenizationAssetInfo] Error Invalid Asset Manager ID: %v\n%v\n", input.AssetManagerID, tokenizationID)
	// 	err = &tErrors.CustomError{Param: "assetManagerID", Err: "error-invalid-asset-manager", ErrMessage: "Invalid Asset Manager. None specified."}
	// 	return
	// }

	// check asset manager ID
	am := GetAssetManagerByID(input.AssetManagerID, gc.DB)
	if len(am.AssetManagerName) == 0 {
		log.Printf("[VetTokenizationAssetInfo] Error Invalid Asset Manager ID: %v\n%v\n", input.AssetManagerID, tokenizationID)
		err = &tErrors.CustomError{Param: "assetManagerID", Err: "error-invalid-asset-manager", ErrMessage: "Invalid Asset Manager."}
		return
	}
	// if input.ApprovedAssetCustodianID == 0 {
	// 	log.Printf("[VetTokenizationAssetInfo] Error Invalid Custodian ID: %v\n%v\n", input.ApprovedAssetCustodianID, tokenizationID)
	// 	err = &tErrors.CustomError{Param: "assetManagerID", Err: "error-invalid-custodian", ErrMessage: "Invalid Asset Custodian. None specified."}
	// 	return
	// }

	// check asset manager ID
	ac := GetApprovedCustodianByID(input.ApprovedAssetCustodianID, gc.DB)
	if len(ac.AssetCustodianName) == 0 {
		log.Printf("[VetTokenizationAssetInfo] Error Invalid custodian ID: %v\n%v\n", input.ApprovedAssetCustodianID, tokenizationID)
		err = &tErrors.CustomError{Param: "approvedCustodianId", Err: "error-invalid-asset-custodian", ErrMessage: "Invalid Asset Custodian."}
		return
	}

	// if input.AssetIssuingHouseID == 0 {
	// 	log.Printf("[VetTokenizationAssetInfo] Error Invalid Issuing House ID: %v\n%v\n", input.AssetIssuingHouseID, tokenizationID)
	// 	err = &tErrors.CustomError{Param: "assetIssuingHouseID", Err: "error-invalid-issuing house", ErrMessage: "Invalid Asset Issuing House. None specified."}
	// 	return
	// }

	// check asset manager ID
	ai := GetAssetIssuingHouseByID(input.AssetIssuingHouseID, gc.DB)
	if len(ai.AssetIssuingHouseName) == 0 {
		log.Printf("[VetTokenizationAssetInfo] Error Invalid assetIssuingHouseID ID: %v\n%v\n", input.AssetIssuingHouseID, tokenizationID)
		err = &tErrors.CustomError{Param: "assetIssuingHouseId", Err: "error-invalid-asset-issuing-house", ErrMessage: "Invalid Asset Issuing House."}
		return
	}

	// if input.LegalAndProfesionalPartnerID == 0 {
	// 	log.Printf("[VetTokenizationAssetInfo] Error Invalid Legal and  Professional Partner: %v\n%v\n", input.LegalAndProfesionalPartnerID, tokenizationID)
	// 	err = &tErrors.CustomError{Param: "legalAndProfesionalPartnerID", Err: "error-invalid-legalAndProfesionalPartnerId", ErrMessage: "Invalid Legal And ProfesionalPartner. None specified."}
	// 	return
	// }

	// check asset manager ID
	lpp := GetLegalAndProfesionalPartnerByID(input.LegalAndProfesionalPartnerID, gc.DB)
	if len(lpp.PartnerName) == 0 {
		log.Printf("[VetTokenizationAssetInfo] Error Invalid LegalAndProfesionalPartnerID: %v\n%v\n", input.LegalAndProfesionalPartnerID, tokenizationID)
		err = &tErrors.CustomError{Param: "assetIssuingHouseId", Err: "error-invalid-legalAndProfesionalPartnerId", ErrMessage: "Invalid Legal And ProfesionalPartner."}
		return
	}

	// if input.RatingAgencyID == 0 {
	// 	log.Printf("[VetTokenizationAssetInfo] Error Invalid RatingAgencyID: %v\n%v\n", input.RatingAgencyID, tokenizationID)
	// 	err = &tErrors.CustomError{Param: "ratingAgencyID", Err: "error-invalid-rating-agency", ErrMessage: "Invalid Rating Agency. None specified."}
	// 	return
	// }

	// check rating agency ID
	ra := GetRatingAgencyByID(input.RatingAgencyID, gc.DB)
	if len(ra.AgencyName) == 0 {
		log.Printf("[VetTokenizationAssetInfo] Error Invalid rating agency: %v\n%v\n", input.RatingAgencyID, tokenizationID)
		err = &tErrors.CustomError{Param: "ratingAgency", Err: "error-invalid-rating-agency", ErrMessage: "Invalid Rating Agency."}
		return
	}

	// check trustee ID
	ts := GetTrusteeByID(input.TrusteeID, gc.DB)
	if len(ra.AgencyName) == 0 {
		log.Printf("[VetTokenizationAssetInfo] Error Invalid trustee: %v\n%v\n", input.TrusteeID, tokenizationID)
		err = &tErrors.CustomError{Param: "trustee", Err: "error-invalid-trustee", ErrMessage: "Invalid Trustee."}
		return
	}

	ato.ExcludeSecFee = input.ExcludeSecFee
	ato.ApprovedAssetCustodianID = input.ApprovedAssetCustodianID

	ato.CustodianFeePercent = func() float64 {
		if input.ApprovedAssetCustodianFeePercent > 0 {
			return input.ApprovedAssetCustodianFeePercent
		}
		return ac.FeePercent
	}()
	ato.CustodianFeeFixed = func() float64 {
		if input.ApprovedAssetCustodianFeeFixed > 0 {
			return input.ApprovedAssetCustodianFeeFixed
		}
		return ac.FeeFixed
	}()
	ato.AssetManagerID = input.AssetManagerID
	ato.AssetManagerFeePercent = func() float64 {
		if input.AssetManagerFeePercent > 0 {
			return input.AssetManagerFeePercent
		}
		return am.FeePercent
	}()
	ato.AssetManagerFeeFixed = func() float64 {
		if input.AssetManagerFeeFixed > 0 {
			return input.AssetManagerFeeFixed
		}
		return am.FeeFixed
	}()
	ato.AssetIssuingHouseID = input.AssetIssuingHouseID
	ato.IssuingHouseFeePercent = func() float64 {
		if input.AssetIssuingHouseFeePercent > 0 {
			return input.AssetIssuingHouseFeePercent
		}
		return ai.FeePercent
	}()
	ato.IssuingHouseFeeFixed = func() float64 {
		if input.AssetIssuingHouseFeeFixed > 0 {
			return input.AssetIssuingHouseFeeFixed
		}
		return ai.FeeFixed
	}()

	ato.LegalAndProfesionalPartnerID = input.LegalAndProfesionalPartnerID
	ato.LegalAndProfessionalFeePercent = func() float64 {
		if input.LegalAndProfesionalPartnerFeePercent > 0 {
			return input.LegalAndProfesionalPartnerFeePercent
		}
		return lpp.FeePercent
	}()
	ato.LegalAndProfessionalFeeFixed = func() float64 {
		if input.LegalAndProfesionalPartnerFeeFixed > 0 {
			return input.LegalAndProfesionalPartnerFeeFixed
		}
		return lpp.FeeFixed
	}()
	ato.RatingAgencyID = input.RatingAgencyID

	ato.RatingAgencyFeePercent = func() float64 {
		if input.RatingAgencyFeePercent > 0 {
			return input.RatingAgencyFeePercent
		}
		return ra.FeePercent
	}()
	ato.RatingAgencyFeeFixed = func() float64 {
		if input.RatingAgencyFeeFixed > 0 {
			return input.RatingAgencyFeeFixed
		}
		return ra.FeeFixed
	}()

	ato.TrusteeID = input.TrusteeID
	ato.TrusteeFeePercent = func() float64 {
		if input.TrusteeFeePercent > 0 {
			return input.TrusteeFeePercent
		}
		return ts.FeePercent
	}()
	ato.TrusteeFeeFixed = func() float64 {
		if input.TrusteeFeeFixed > 0 {
			return input.TrusteeFeeFixed
		}
		return ts.FeeFixed
	}()

	if len(input.AssetQuoteCurrency) > 0 {
		ato.AssetQuoteCurrency = &input.AssetQuoteCurrency
	}
	if len(input.ProceedPayoutCurrency) == 0 {
		input.ProceedPayoutCurrency = input.AssetQuoteCurrency
	}
	if len(input.ProceedPayoutCurrency) > 0 {
		input.ProceedPayoutCurrency = strings.ToUpper(input.ProceedPayoutCurrency)
		ato.ProceedPayoutCurrency = &input.ProceedPayoutCurrency
	}
	ato.VettingStatus = 1

	ato.LastUpdatedBy = &initiator.Username

	ato.UpdateCalculation(gc)

	e := gc.DB.Omit(clause.Associations).Save(&ato).Error
	if e != nil {
		log.Printf("[VetTokenizationAssetInfo] error saving tokenization to database  [%v] for %v: %v\n", input, initiator.Username, e)

		err = &tErrors.ErrorTemporaryServerError{}

	}
	ato, _, _ = GetTokenizedAssetByID(ato.ID, gc.DB)
	return ato, nil
}

// AcknowledgeTokenizationFeePayment used by trovoManager
func AcknowledgeTokenizationFeePayment(tokenizationID string, initiator *userModels.User, gc *sharedconfig.GlobalConfig) (ato userModels.TokenizedAsset, err error) {

	ato, _, errorGetTokenizationByID := GetTokenizedAssetByID(tokenizationID, gc.DB)
	if errorGetTokenizationByID != nil {
		// error tokenization is already in progress
		log.Printf("[AcknowledgeTokenizationFeePayment] Error fetching  tokenization with ID: %v\n", tokenizationID)
		err = errorGetTokenizationByID
		return

	}

	//tokenization existing
	if ato.AssetTokenizationStatus != 2 {
		// error tokenization is already in progress
		log.Printf("[AcknowledgeTokenizationFeePayment] Error tokenization application is not awaiting payment confirmation and cannot be modified: %v\n", tokenizationID)
		err = &tErrors.CustomError{Param: "issuingWalletPublicKey", Err: "error-tokenization-cannot-be-modified-by-this-method", ErrMessage: "Tokenization request is not awaiting payment confirmation."}
		return

	}

	ato.AssetTokenizationStatus = 3

	ato.LastUpdatedBy = &initiator.Username

	e := gc.DB.Omit(clause.Associations).Save(&ato).Error
	if e != nil {
		log.Printf("[AcknowledgeTokenizationFeePayment] error saving tokenization to database [%v] %v\n", initiator.Username, e)

		err = &tErrors.ErrorTemporaryServerError{}

	}
	ato, _, _ = GetTokenizedAssetByID(ato.ID, gc.DB)
	return ato, nil
}

// FailTokenizationDueDiligence used by trovoManager
func FailTokenizationDueDiligence(tokenizationID string, trovoManagerUser *userModels.User, reasonForFailure string, gc *sharedconfig.GlobalConfig) (ato userModels.TokenizedAsset, err error) {

	ato, _, errorGetTokenizationByID := GetTokenizedAssetByID(tokenizationID, gc.DB)
	if errorGetTokenizationByID != nil {
		// error tokenization is already in progress
		log.Printf("[FailTokenizationDueDiligence] Error fetching  tokenization with ID: %v\n", tokenizationID)
		err = errorGetTokenizationByID
		return

	}

	//tokenization existing
	if len(reasonForFailure) < 5 {
		// error no reason given
		log.Printf("[FailTokenizationDueDiligence] Error No reason given for due diligence failure %v\n", tokenizationID)
		err = &tErrors.CustomError{Param: "issuingWalletPublicKey", Err: "error-no-reason-given", ErrMessage: "Valid reason for failure of due diligence must be provided."}
		return

	}

	// reset status to 0
	ato.AssetTokenizationStatus = 0
	ato.VettingStatus = 0

	ato.DueDiligenceFail = 1
	ato.DueDiligenceFailureReason = &reasonForFailure

	ato.LastUpdatedBy = &trovoManagerUser.Username

	e := gc.DB.Omit(clause.Associations).Save(&ato).Error
	if e != nil {
		log.Printf("[FailTokenizationDueDiligence] error saving tokenization to database [%v]  %v\n", trovoManagerUser.Username, e)

		err = &tErrors.ErrorTemporaryServerError{}

	}
	ato, _, _ = GetTokenizedAssetByID(ato.ID, gc.DB)
	return ato, nil
}

// UpdateTokenizationSalesDates used by trovoManager to update tokenized asset sales dates
func UpdateTokenizedAssetSalesDates(tokenizationID string, trovoManagerUser *userModels.User, dateInput *userModels.TokenizedAssetSalesDatesInput, gc *sharedconfig.GlobalConfig) (ato userModels.TokenizedAsset, err error) {

	ato, _, errorGetTokenizationByID := GetTokenizedAssetByID(tokenizationID, gc.DB)
	if errorGetTokenizationByID != nil {
		log.Printf("[UpdateTokenizedAssetSalesDates] Error fetching  tokenization with ID: %v\n", tokenizationID)
		err = errorGetTokenizationByID
		return

	}
	if len(dateInput.SalesStart) == 0 && len(dateInput.SalesEnd) == 0 {
		log.Printf("[UpdateTokenizedAssetSalesDates] No sales/end dates to change tokenization with ID: %v\n", tokenizationID)

		return
	}

	if len(dateInput.SalesStart) > 0 {
		// parse date
		eventDate, e := time.Parse("2006-01-02", dateInput.SalesStart)
		if e != nil {
			log.Printf("[UpdateTokenizedAssetSalesDates] Error invalid sales start date format %v. [%v]\n", dateInput.SalesStart, e)
			err = &tErrors.CustomError{Param: "salesStart", Err: "error-invalid-date-format", ErrMessage: "Invalid date format for sales start. Use YYYY-MM-DD."}
			return
		}
		ato.SalesStart = eventDate

	}

	if len(dateInput.SalesEnd) > 0 {
		// parse date
		eventDate, e := time.Parse("2006-01-02", dateInput.SalesEnd)
		if e != nil {
			log.Printf("[UpdateTokenizedAssetSalesDates] Error invalid sales end date format %v. [%v]\n", dateInput.SalesStart, e)
			err = &tErrors.CustomError{Param: "salesEnd", Err: "error-invalid-date-format", ErrMessage: "Invalid date format for sales end. Use YYYY-MM-DD."}
			return
		}
		ato.SalesEnd = eventDate

	}

	if ato.SalesStart.Compare(ato.SalesEnd) == 0 || ato.SalesStart.Compare(ato.SalesEnd) == +1 {
		log.Printf("[UpdateTokenizedAssetSalesDates] Error invalid sales start/end date %v. [%v]\n", dateInput.SalesStart, dateInput.SalesEnd)
		err = &tErrors.CustomError{Param: "salesEnd", Err: "error-invalid-date-order", ErrMessage: "Dates for sales start/end are not logical. End date must be set to be ahead of Start date."}
		return
	}

	ato.LastUpdatedBy = &trovoManagerUser.Username

	e := gc.DB.Omit(clause.Associations).Save(&ato).Error
	if e != nil {
		log.Printf("[UpdateTokenizedAssetSalesDates] error saving tokenization sales dates to database [%v]  %v\n", trovoManagerUser.Username, e)

		err = &tErrors.ErrorTemporaryServerError{}

	}
	ato, _, _ = GetTokenizedAssetByID(ato.ID, gc.DB)
	return ato, nil
}

// ActivatePrimarySalesRoutine used by automation routine to update tokenized assets to begin sales
func ActivatePrimarySalesRoutine(gc *sharedconfig.GlobalConfig) {
	log.Println("##[ActivatePrimarySalesRoutine][CHECK SALES DATES] started routine to update asset sales status")

	batchSize := 1
	var assets []userModels.TokenizedAsset

	{
		//start primnary sales

		e := gc.DB.Where("Sales_Start::date <= now()::date AND Asset_Tokenization_Status = ?", 4).First(&userModels.TokenizedAsset{}).Error
		if e == nil {
			result := gc.DB.Where("Sales_Start::date <= now()::date AND Asset_Tokenization_Status = ?", 4).FindInBatches(&assets, batchSize, func(tx *gorm.DB, batch int) error {
				for i, asset := range assets {

					//check if it has been minted.
					_, err := userModels.BantuAsset{AssetCode: *asset.AssetCode, AssetIssuer: *asset.IssuingWalletPublicKey}.GetBlockchainAssetProperty(gc)
					if err != nil {
						log.Printf("[ActivatePrimarySalesRoutine][CHECK PRIMARY SALES DATES]()()()@@@()()()()FAILED TO CONFIRM MINTING of %v on blockchain due to: %v\n", *asset.AssetCode, err)

						continue
					}

					// assets[i].AssetTokenizationStatus = 5
					assets[i].AssetTokenizationStatus = 5

					log.Printf("[ActivatePrimarySalesRoutine][CHECK PRIMARY SALES DATES]UPDATE ASSET : %v, %v\n", *asset.AssetCode, asset.AssetTokenizationStatus)

					gc.ChannelOfTokenizedAssetIDs <- asset.ID
				}
				e := tx.Omit(clause.Associations).Save(&assets).Error
				if e != nil {
					//saving model failed
					log.Printf("[ActivatePrimarySalesRoutine][CHECK PRIMARY SALES DATES]()()()@@@()()()()FAILED TO UPDATE ASSET  with status due to: %v\n", e)
				}
				time.Sleep(200 * time.Millisecond)

				return nil
			})
			if result.Error != nil {
				if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
					log.Println("[ActivatePrimarySalesRoutine][CHECK PRIMARY SALES DATES]()()()()()()()()()()()()()error occurred during batch processing:", result.Error.Error())

				}

			}
		}

	}

	time.Sleep(15 * time.Minute)

}

// ActivateSecondarySalesRoutine used by automation routine to update tokenized assets to begin sales
func ActivateSecondarySalesRoutine(gc *sharedconfig.GlobalConfig) {
	log.Println("##[ActivateSecondarySalesRoutine][CHECK SALES DATES] started routine to update asset sales status")

	batchSize := 1
	var assets []userModels.TokenizedAsset

	{
		//start secondary sales

		e := gc.DB.Where("Sales_End::date <= now()::date AND Asset_Tokenization_Status = ?", 5).First(&userModels.TokenizedAsset{}).Error
		if e == nil {
			result := gc.DB.Where("Sales_End::date <= now()::date AND Asset_Tokenization_Status = ?", 5).FindInBatches(&assets, batchSize, func(tx *gorm.DB, batch int) error {
				for _, asset := range assets {

					asset.AssetTokenizationStatus = 6

					e := tx.Omit(clause.Associations).Save(&asset).Error
					if e != nil {
						//saving model failed
						log.Printf("[ActivateSecondarySalesRoutine][CHECK SECONDARY SALES DATES]()()()@@@()()()()FAILED TO UPDATE ASSET LIST with status due to: %v\n", e)
					}
					log.Printf("[ActivateSecondarySalesRoutine][CHECK PRIMARY SALES DATES]UPDATE ASSET : %v, %v\n", *asset.AssetCode, asset.AssetTokenizationStatus)

				}

				time.Sleep(200 * time.Millisecond)

				return nil
			})
			if result.Error != nil {
				if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
					log.Println("[ActivateSecondarySalesRoutine][CHECK SECONDARY SALES DATES]()()()()()()()()()()()()()error occurred during batch processing:", result.Error.Error())

				}

			}
		}

	}

	time.Sleep(15 * time.Minute)

}

// SendPNToSuscribersForPrimarySales used by automation routine to send PN for tokenized assets to begin sales
func SendPNToSuscribersForPrimarySales(gc *sharedconfig.GlobalConfig) {
	chant := <-gc.ChannelOfTokenizedAssetIDs
	t := userModels.TokenizedAssetID(chant).GetTokenization(gc)
	if t.ID != chant {
		return
	}
	batchSize := 100
	var ei []userModels.ExpressionOfInterest
	log.Println("[SendPNToSuscribersForPrimarySales]()()()()()()()()()()()()()Starting push notification for asset:", *t.AssetCode)

	e := gc.DB.Where("Tokenized_Asset_ID = ?", t.ID).First(&userModels.ExpressionOfInterest{}).Error
	if e != nil {
		return
	}

	result := gc.DB.Where("Tokenized_Asset_ID = ?", t.ID).FindInBatches(&ei, batchSize, func(tx *gorm.DB, batch int) error {
		for _, ex := range ei {

			u, err := userModels.Username(ex.SubscriberUsername).GetSimpleUser(gc.DB, gc)
			if err != nil {
				continue
			}
			if u.PushNotificationToken != nil {
				dataPayload := make(map[string]string)
				dataPayload["route"] = "assetSubscription"
				msgBody := fmt.Sprintf("You can now go to your Trovo App and purchase %v (%v). Only %v units @ %v %v are available for the primary sale. So, hurry now!", *t.AssetCode, *t.AssetName, t.MaxNumberOfTokenAvailableForSale, t.PricePerToken, *t.AssetQuoteCurrency)
				if t.CapAmountInFiat > 0 {
					msgBody = fmt.Sprintf("You can now go to your Trovo App and purchase %v (%v). Only %v units @ %v%v are  available for the primary sale, capped at %v %v per person.", *t.AssetCode, *t.AssetName, t.MaxNumberOfTokenAvailableForSale, t.PricePerToken, *t.AssetQuoteCurrency, t.CapAmountInFiat, *t.AssetQuoteCurrency)
				}
				u.SendPushMessage(fmt.Sprintf("The asset %v is now live on sale @ %v%v per unit!", *t.AssetCode, t.PricePerToken, *t.AssetQuoteCurrency), msgBody, *t.AssetLogo, dataPayload, gc)
			}

		}

		return nil
	})
	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Println("[SendPNToSuscribersForPrimarySales]()()()()()()()()()()()()()error occurred during batch processing:", result.Error.Error())

		}

	}

	time.Sleep(5 * time.Minute)

}

// ConfirmTokenizationApplicationInfoByInitiator used by original owner/initiator to advance status to 1 and allow for vetting.
func ConfirmTokenizationApplicationInfoByInitiator(initiator *userModels.User, tokenizationID string, taInput *userModels.ConfirmTokenizedAssetJSONInput, gc *sharedconfig.GlobalConfig) (ato userModels.TokenizedAsset, err error) {
	taInput.NetworkPassPhrase = gc.BantuNetworkPassphrase
	taInput.Messages = make([]string, 0)
	//check if existing
	ato, NotFound, e := GetOpenTokenizedAssetByInitiatorUsername(initiator.Username, gc.DB)

	if e == nil {
		//tokenization existing
		if ato.AssetTokenizationStatus > 0 {
			// error tokenization is already in progress
			log.Printf("[ConfirmTokenizationApplicationInfoByInitiator] Error tokenization procesing is in progress and cannot be modified: %v\n", tokenizationID)
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
			log.Printf("[ConfirmTokenizationApplicationInfoByInitiator]error fetching existing tokenization from database ID [%v] for %v: %v\n", tokenizationID, initiator.Username, e)
			err = &tErrors.ErrorTemporaryServerError{}
			return

		}

		log.Printf("[ConfirmTokenizationApplicationInfoByInitiator] Tokenization does not exist: %v\n", tokenizationID)
		err = &tErrors.CustomError{Param: "Id", Err: "error-tokenization-not-found", ErrMessage: "Only existing valid tokenization requests can be confirmed."}
		return

	}
	if ato.AssetCountryLocation == nil {
		err = &tErrors.CustomError{Param: "Id", Err: "error-asset-location-country-not-found", ErrMessage: "You must specify the asset country of location."}
		return
	}
	if len(*ato.AssetCountryLocation) != 2 {
		err = &tErrors.CustomError{Param: "Id", Err: "error-asset-location-country-not-found", ErrMessage: "You must specify the asset country of location in the formart: NG, SA, UK"}
		return
	}

	//begin a database transaction here
	dbTX := gc.DB.Begin()
	defer dbTX.Rollback()

	e = dbTX.Omit(clause.Associations).Save(&ato).Error
	if e != nil {
		log.Printf("[ConfirmTokenizationApplicationInfoByInitiator] error saving tokenization to database  [%+v] for %v: %v\n", ato, initiator.Username, e)

		err = &tErrors.ErrorTemporaryServerError{}

	}

	//check for blockchain action
	wallet, e := userModels.WalletAlias(initiator.Username).GetWallet(dbTX, gc)
	if e != nil {
		log.Printf("[ConfirmTokenizationApplicationInfoByInitiator] error getting primary wallet from database  [%+v] for %v: %v\n", ato, initiator.Username, e)

		err = &tErrors.ErrorTemporaryServerError{}

		return

	}

	xdrBase64, e := generateTokenizationFeeXdr(&wallet, &ato, taInput, gc)
	if e != nil {
		log.Printf("[ConfirmTokenizationApplicationInfoByInitiator] error getting application fee transaction  [%+v] for %v: %v\n", ato, initiator.Username, e)

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

// ConfirmTokenizationFeePaymentByInitiator is used to confirm payment by initiator. After successful call, it moves the tokenized asset status to trovo manager awaiting payment acknowledgement.
func ConfirmTokenizationFeePaymentByInitiator(initiator *userModels.User, tokenizationID string, gc *sharedconfig.GlobalConfig) (ato userModels.TokenizedAsset, err error) {

	//check if existing
	ato, NotFound, e := GetFeeReadyTokenizedAssetApplicationByInitiatorUsername(initiator.Username, gc.DB)

	if e == nil {
		//tokenization existing
		if ato.AssetTokenizationStatus != 1 {
			// error tokenization is already in progress
			log.Printf("[ConfirmTokenizationFeePaymentByInitiator] Error tokenization process not awaiting payment and cannot be modified: %v\n", tokenizationID)
			err = &tErrors.CustomError{Param: "issuingWalletPublicKey", Err: "error-tokenization-cannot-be-modified-by-this-method", ErrMessage: "Tokenization is not awaiting payment, this action cannot be performed."}
			return

		}

		if ato.ID != tokenizationID {
			log.Printf("[ConfirmTokenizationFeePaymentByInitiator] Error tokenization invalid tokenization ID: %v\n", tokenizationID)

			err = &tErrors.CustomError{Param: "issuingWalletPublicKey", Err: "error-tokenization-id-not-valid", ErrMessage: "Invalid tokenization specified."}
			return
		}

		if ato.VettingStatus == 0 {
			log.Printf("[ConfirmTokenizationFeePaymentByInitiator] Error Tokenization request is still being vetted by the team. Please wait until vetting has completed. tokenization ID: %v\n", tokenizationID)

			err = &tErrors.CustomError{Param: "vettingStatus", Err: "error-tokenization-not-vetted", ErrMessage: "Tokenization request is still being vetted by the team. Please wait until vetting has completed."}
			return
		}

		ato.LastUpdatedBy = &initiator.Username
		ato.AssetTokenizationStatus = 2

	} else {
		if !NotFound {
			//critical database error occured
			log.Printf("[ConfirmTokenizationFeePaymentByInitiator]error fetching existing tokenization from database  ID [%v] for %v: %v\n", tokenizationID, initiator.Username, e)
			err = &tErrors.ErrorTemporaryServerError{}
			return

		}

		log.Printf("[ConfirmTokenizationFeePaymentByInitiator] Tokenization does not exist: %v\n", tokenizationID)
		err = &tErrors.CustomError{Param: "Id", Err: "error-tokenization-not-found", ErrMessage: "Only existing valid tokenization requests can be confirmed."}
		return

	}

	e = gc.DB.Omit(clause.Associations).Save(&ato).Error
	if e != nil {
		log.Printf("[ConfirmTokenizationFeePaymentByInitiator] error saving tokenization to database  [%+v] for %v: %v\n", ato, initiator.Username, e)

		err = &tErrors.ErrorTemporaryServerError{}

	}
	ato, _, _ = GetTokenizedAssetByID(ato.ID, gc.DB)
	return ato, err
}

func GetTokenizationList(user *userModels.User, adminList bool, gc *sharedconfig.GlobalConfig, c *gin.Context) (records userModels.PaginatedTokenizedAssets) {
	var err error
	var tokenizedAssetList []userModels.TokenizedAsset
	records.Records = make([]userModels.TokenizedAssetJSON, 0)
	// DB, _ := db.OpenDb()
	// DBC, _ := db.OpenDb()
	DB := gc.DB
	DBC := gc.DB

	onlyWithUserPermission := strings.TrimSpace(strings.ToUpper(c.DefaultQuery("onlyWithUserPermission", "1")))
	if adminList {
		onlyWithUserPermission = "0"
	}

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
	tokenizationStatus := strings.TrimSpace(c.Query("assetTokenizationStatus"))
	salesList := strings.TrimSpace(c.Query("salesList")) //0=primary+awaiting minting,1=secondary sales
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

	orderBy := strings.TrimSpace(c.Query("orderby"))
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
		query = query.Order("updated_at desc, asset_tokenization_status asc, asset_country_location asc")
		countQuery = countQuery.Order("updated_at desc, asset_tokenization_status asc, asset_country_location asc")
	}

	if onlyWithUserPermission == "1" {
		// query = query.Where("(issuing_wallet_public_key IN (?))", sharedWallets)
		// countQuery = countQuery.Where("(issuing_wallet_public_key IN (?))", sharedWallets)
		query = query.Where("lower(initiator_username) = lower(?)", user.Username)
		countQuery = countQuery.Where("lower(initiator_username) = lower(?)", user.Username)
	}

	if len(initiatorUsername) > 2 {
		query = query.Where("lower(initiator_username) = ?", initiatorUsername)
		countQuery = countQuery.Where("lower(initiator_username) = ?", initiatorUsername)

	}

	if len(tokenizationStatus) > 0 && onlyWithUserPermission == "1" {
		//status is specified for permission viewing
		ts, e := decimal.NewFromString(strings.TrimSpace(tokenizationStatus))
		if e != nil {
			ts = decimal.Zero
		}
		tsInt := int(ts.IntPart())
		query = query.Where("asset_tokenization_status = ?", tsInt)
		countQuery = countQuery.Where("asset_tokenization_status = ?", tsInt)

	}

	if adminList {

		query = query.Where("asset_tokenization_status >= ?", int(1))
		countQuery = countQuery.Where("asset_tokenization_status >= ?", int(1))

	}

	if len(tokenizationStatus) > 0 && onlyWithUserPermission == "0" {
		//status is specified for open viewing
		ts, e := decimal.NewFromString(strings.TrimSpace(tokenizationStatus))
		if e != nil {
			ts = decimal.Zero
		}
		tsInt := int(ts.IntPart())
		if tsInt > 3 {
			//market ready status for open list
			query = query.Where("asset_tokenization_status = ?", tsInt)
			countQuery = countQuery.Where("asset_tokenization_status = ?", tsInt)

		}

	}

	//get only market ready list
	if onlyWithUserPermission == "0" && len(tokenizationStatus) == 0 && len(salesList) == 0 && !adminList {
		//status is not specified for open viewing
		query = query.Where("asset_tokenization_status > ?", 3)
		countQuery = countQuery.Where("asset_tokenization_status > ?", 3)

	}

	//get only market ready list
	if onlyWithUserPermission == "0" && len(tokenizationStatus) == 0 && salesList == "0" && !adminList {
		//status is not specified for open viewing
		query = query.Where("asset_tokenization_status > ? AND asset_tokenization_status < 6", 3)
		countQuery = countQuery.Where("asset_tokenization_status > ? AND asset_tokenization_status < 6", 3)

	}

	//get only market ready list
	if onlyWithUserPermission == "0" && len(tokenizationStatus) == 0 && salesList == "1" && !adminList {
		//status is not specified for open viewing
		query = query.Where("asset_tokenization_status = ?", 6)
		countQuery = countQuery.Where("asset_tokenization_status = ?", 6)

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

func GetTokenizedAssetSubscriptionList(user *userModels.User, gc *sharedconfig.GlobalConfig, c *gin.Context) (records userModels.PaginatedTokenizedAssetSubscription) {
	var err error
	var eiList []userModels.TokenizedAssetSubscription
	records.Records = make([]userModels.TokenizedAssetSubscription, 0)
	// DB, _ := db.OpenDb()
	// DBC, _ := db.OpenDb()
	DB := gc.DB
	DBC := gc.DB

	onlySelf := strings.TrimSpace(strings.ToUpper(c.DefaultQuery("onlySelf", "1")))
	//get all wallets where user has access
	sharedWallets := make([]string, 0)
	if onlySelf == "1" {
		for _, w := range user.WalletsSharedWithUser {
			sharedWallets = append(sharedWallets, w.WalletPublicKey)
		}
	}
	var query *gorm.DB
	var countQuery *gorm.DB
	oD := "ASC"

	assetCode := strings.TrimSpace(strings.ToUpper(c.Query("assetCode")))
	subscriberUsername := strings.TrimSpace(c.Query("subscriberUsername"))
	walletPublicKey := strings.TrimSpace(c.Query("walletPublicKey"))

	limitU, _ := strconv.ParseUint(strings.TrimSpace(c.DefaultQuery("limit", "25")), 10, 64)
	limit := int(limitU)
	pageU, _ := strconv.ParseUint(strings.TrimSpace(c.DefaultQuery("page", "1")), 10, 64)
	page := int(pageU)

	createdBetween := strings.TrimSpace(c.Query("createdBetween"))
	amountBetween := strings.TrimSpace(c.Query("amountBetween"))

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
		query = query.Order("created_at DESC")
		countQuery = countQuery.Order("created_at DESC")
	}

	if onlySelf == "1" {
		// (wallet_public_key IN (?))
		query = query.Where("(Subscriber_Username = lower(?) OR wallet_public_key IN (?))", user.Username, sharedWallets)
		countQuery = countQuery.Where("(Subscriber_Username = lower(?) OR wallet_public_key IN (?))", user.Username, sharedWallets)
	}

	if len(walletPublicKey) > 50 {

		query = query.Where("wallet_Public_Key = ?", walletPublicKey)
		countQuery = countQuery.Where("wallet_Public_Key = ?", walletPublicKey)
	}

	if onlySelf == "0" && len(subscriberUsername) > 0 {

		query = query.Where("Subscriber_Username = lower(?)", subscriberUsername)
		countQuery = countQuery.Where("Subscriber_Username = lower(?)", subscriberUsername)
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
	if len(amountBetween) > 0 && strings.Contains(amountBetween, "|") {
		amountRange := strings.Split(amountBetween, "|")
		lAmount, e := decimal.NewFromString(amountRange[0])
		if e != nil {
			lAmount = decimal.Zero
		}
		hAmount, e := decimal.NewFromString(amountRange[1])
		if e != nil {
			hAmount = decimal.Zero
		}
		if lAmount.InexactFloat64() > 0 || hAmount.InexactFloat64() > 0 {
			query = query.Where("amount::numeric BETWEEN ?::numeric AND ?::numeric", lAmount, hAmount)
			countQuery = countQuery.Where("amount::numeric BETWEEN ?::numeric AND ?::numeric", lAmount, hAmount)

		}

	}

	var countR int64

	errCount := countQuery.Find(&[]userModels.TokenizedAssetSubscription{}).Count(&countR).Error
	if errCount != nil {
		log.Println("[GetTokenizedAssetSubscriptionList]Count Error:", errCount)
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
	if err = query.Find(&eiList).Error; err != nil {
		log.Println("[GetTokenizedAssetSubscriptionList] Query Error:", err)
		return
	}

	tListJSON := make([]userModels.TokenizedAssetSubscription, 0)

	for _, v := range eiList {
		j := v
		tListJSON = append(tListJSON, j)

	}

	records = userModels.PaginatedTokenizedAssetSubscription{CurrentPage: page, Pages: pages, TotalRecords: count, Limit: limit, Records: tListJSON}

	return records
}

func GetExpressionOfInterestList(user *userModels.User, gc *sharedconfig.GlobalConfig, c *gin.Context) (records userModels.PaginatedExpressionOfInterest) {
	var err error
	var eiList []userModels.ExpressionOfInterest
	records.Records = make([]userModels.ExpressionOfInterest, 0)
	DB, _ := db.OpenDb()
	DBC, _ := db.OpenDb()

	onlySelf := strings.TrimSpace(strings.ToUpper(c.DefaultQuery("onlySelf", "1")))

	var query *gorm.DB
	var countQuery *gorm.DB
	oD := "ASC"

	assetCode := strings.TrimSpace(strings.ToUpper(c.Query("assetCode")))
	subscriberUsername := strings.TrimSpace(c.Query("subscriberUsername"))

	limitU, _ := strconv.ParseUint(strings.TrimSpace(c.DefaultQuery("limit", "25")), 10, 64)
	limit := int(limitU)
	pageU, _ := strconv.ParseUint(strings.TrimSpace(c.DefaultQuery("page", "1")), 10, 64)
	page := int(pageU)

	createdBetween := strings.TrimSpace(c.Query("createdBetween"))
	amountBetween := strings.TrimSpace(c.Query("amountBetween"))

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
		query = query.Order("created_at DESC")
		countQuery = countQuery.Order("created_at DESC")
	}

	if onlySelf == "1" {

		query = query.Where("Subscriber_Username = lower(?)", user.Username)
		countQuery = countQuery.Where("Subscriber_Username = lower(?)", user.Username)
	}

	if onlySelf == "0" && len(subscriberUsername) > 0 {

		query = query.Where("Subscriber_Username = lower(?)", subscriberUsername)
		countQuery = countQuery.Where("Subscriber_Username = lower(?)", subscriberUsername)
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
	if len(amountBetween) > 0 && strings.Contains(amountBetween, "|") {
		amountRange := strings.Split(amountBetween, "|")
		lAmount, e := decimal.NewFromString(amountRange[0])
		if e != nil {
			lAmount = decimal.Zero
		}
		hAmount, e := decimal.NewFromString(amountRange[1])
		if e != nil {
			hAmount = decimal.Zero
		}
		if lAmount.InexactFloat64() > 0 || hAmount.InexactFloat64() > 0 {
			query = query.Where("amount::numeric BETWEEN ?::numeric AND ?::numeric", lAmount, hAmount)
			countQuery = countQuery.Where("amount::numeric BETWEEN ?::numeric AND ?::numeric", lAmount, hAmount)

		}

	}

	var countR int64

	errCount := countQuery.Find(&[]userModels.ExpressionOfInterest{}).Count(&countR).Error
	if errCount != nil {
		log.Println("[GetExpressionOfInterestList]Count Error:", errCount)
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
	if err = query.Find(&eiList).Error; err != nil {
		log.Println("[GetExpressionOfInterestList] Query Error:", err)
		return
	}

	tListJSON := make([]userModels.ExpressionOfInterest, 0)

	for _, v := range eiList {
		j := v
		tListJSON = append(tListJSON, j)

	}

	records = userModels.PaginatedExpressionOfInterest{CurrentPage: page, Pages: pages, TotalRecords: count, Limit: limit, Records: tListJSON}

	return records
}

func SubscribeToTokenizedAsset(subscriber *userModels.User, subscriberWallet *userModels.UserWallet, ta *userModels.TokenizedAsset, input *userModels.TokenizedAssetSubscriptionInput, gc *sharedconfig.GlobalConfig) (taSubscription userModels.TokenizedAssetSubscription, err error) {
	input.TokenizedAssetID = ta.ID
	input.WalletPublicKey = subscriberWallet.ID
	input.SubscriberUsername = subscriber.Username
	input.Amount = decimal.NewFromFloat(input.Amount).Truncate(7).InexactFloat64()
	var swapInfo swapModels.SwapSendInfo
	swapInfo.Messages = make([]string, 0)
	swapInfo.SourceAmount = decimal.NewFromFloat(input.Amount).Truncate(7).String()
	swapInfo.SwapAmount = swapInfo.SourceAmount
	// get wallet owner
	walletOwner, e := subscriberWallet.GetWalletOwner(gc.DB, gc)
	if e != nil {

		log.Printf("[SubscribeToTokenizedAsset] Error Unable to verify wallet owner for subscribing wallet %v\n", subscriberWallet.Alias)
		err = &tErrors.CustomError{Param: "issuingWalletPublicKey", Err: "error-invalid-kyc", ErrMessage: "Unable to verify wallet owner."}
		return

	}

	if walletOwner.KYCVerified == 0 {

		log.Printf("[SubscribeToTokenizedAsset] Error Wallet owner %v has not met KYC status for asset %v\n", walletOwner.Username, *ta.AssetCode)
		err = &tErrors.CustomError{Param: "issuingWalletPublicKey", Err: "error-invalid-kyc", ErrMessage: fmt.Sprintf("%v has not passed KYC to purchase this tokenized asset %v.", walletOwner.Username, *ta.AssetCode)}
		return

	}
	if ta.AssetTokenizationStatus != 5 && ta.AssetTokenizationStatus != 6 {

		log.Printf("[SubscribeToTokenizedAsset] Error Tokenized asset not in sales yet: %v\n", ta.ID)
		err = &tErrors.CustomError{Param: "issuingWalletPublicKey", Err: "error-invalid-request", ErrMessage: "Only projects that are in sales can accept purchase."}
		return

	}
	capEndDate := ta.SalesStart.AddDate(0, 0, ta.CapDurationInDays)
	if capEndDate.After(time.Now()) && ta.CapOnPurchase > 0 && ta.CapAmountInFiat > 0 && decimal.NewFromFloat(input.Amount+ta.SumAmountBoughtByWalletOwner(subscriberWallet.Alias, gc)).Truncate(7).GreaterThan(decimal.NewFromFloat(ta.CapAmountInFiat)) {

		log.Printf("[SubscribeToTokenizedAsset] Error Tokenized asset Cap exceeded: %v, Amount In Cap: %v\n", ta.ID, decimal.NewFromFloat(ta.CapAmountInFiat).String())
		err = &tErrors.CustomError{Param: "amount", Err: "error-cap-amount-exceeded", ErrMessage: fmt.Sprintf("You can only purchase not more than %v%v worth of %v at this time.", *ta.AssetQuoteCurrency, decimal.NewFromFloat(ta.CapAmountInFiat-ta.SumAmountBoughtByWalletOwner(subscriberWallet.Alias, gc)), *ta.AssetCode)}
		return

	}

	if subscriberWallet.SharedAccessEnabled == 1 && subscriberWallet.NumberOfApprovalsNeeded > 0 {
		input.Multiparty = 1
		swapInfo.Multiparty = 1
	}

	if subscriberWallet.HasViewOnlyAccess(gc) {
		input.SignatureRequired = 1
		swapInfo.SignatureRequired = 1
	}
	dbTX := gc.DB.Begin()
	defer dbTX.Rollback()

	//always save subscriptions afresh
	taSubscription.UpdateTokenizedAssetSubscriptionFromInput(subscriber.Username, subscriberWallet, input, ta, gc)
	if input.Multiparty == 0 {
		e := dbTX.Omit(clause.Associations).Save(&taSubscription).Error
		if e != nil {
			log.Printf("[SubscribeToTokenizedAsset] error saving tokenized asset subscription  to database  [%+v] for %v: %v\n", taSubscription, subscriber.Username, e)

			err = &tErrors.ErrorTemporaryServerError{}
			return
		}
	}

	//begin transaction xdr

	// swapInfo.SourceAmount = decimal.NewFromFloat(input.Amount).Truncate(7).String()
	// swapAmount := decimal.NewFromFloat(input.Amount).Truncate(7)

	client := gc.BantuExpansionClient
	//transform codes and issuer
	swapInfo.DestinationAssetCode = strings.ToUpper(*ta.AssetCode)
	swapInfo.DestinationAssetIssuer = strings.ToUpper(*ta.IssuingWalletPublicKey)
	swapInfo.SourceAssetCode = strings.ToUpper(*ta.AssetQuoteCurrency)
	// get currency
	quoteCurrency := GetTokenizationCurrencyByCode(*ta.AssetQuoteCurrency, dbTX)
	swapInfo.SourceAssetIssuer = strings.ToUpper(quoteCurrency.AssetIssuer)
	if e := swapServices.ValidateSwapSendInfo(&swapInfo); e != nil {
		log.Printf("[SubscribeToTokenizedAsset] error validating tokenized asset subscription [%+v] for %v: %v\n", taSubscription, subscriber.Username, e)

		err = e

		return
	}

	if len(input.TransactionSignature) == 0 {
		xdrBase64, e := generateAssetSubscriptionXdr(subscriberWallet, &swapInfo, gc)
		if e != nil {
			log.Printf("[SubscribeToTokenizedAsset] error generating tokenized asset subscription xdr [%+v] for %v: %v\n", taSubscription, subscriber.Username, e)
			if strings.Contains(e.Error(), "liquid") || strings.Contains(e.Error(), "market") {
				destAsset := os.Getenv("NATIVE_ASSET_CODE")
				sourceAsset := os.Getenv("NATIVE_ASSET_CODE")
				if len(swapInfo.SourceAssetCode) > 0 {
					sourceAsset = swapInfo.SourceAssetCode
				}
				if len(swapInfo.DestinationAssetCode) > 0 {
					destAsset = swapInfo.DestinationAssetCode
				}
				_, b, _ := gc.GetAvalableMarketQuantity(swapInfo.SourceAssetCode, swapInfo.SourceAssetIssuer, swapInfo.DestinationAssetCode, swapInfo.DestinationAssetIssuer)

				emsg := fmt.Sprintf("There is no %v market to exchange for your %v at this time. Please try again later or reduce the quantity of %v to try again.", destAsset, sourceAsset, sourceAsset)
				if b != "0" {
					emsg = fmt.Sprintf("There is only %v %v to exchange for your %v at this time. Please reduce the quantity of %v to try again.", b, destAsset, sourceAsset, sourceAsset)

				}
				err = &tErrors.CustomError{
					Param:      "destinationAssetCode",
					Err:        "error-low-liquidity",
					ErrMessage: emsg,
				}
				return
			}
			err = e
			return
		}

		swapInfo.Transaction = xdrBase64
		input.Transaction = xdrBase64
		input.Messages = swapInfo.Messages
		input.SwappedEstimate = swapInfo.SwappedEstimate
		input.TransactionSource = swapInfo.TransactionSource
		input.Memo = swapInfo.Memo

	}

	input.NetworkPassPhrase = network.GetBlockchainNetworkPassPhrase()

	if len(input.TransactionSignature) == 0 && input.Commit == 0 {
		err = nil
		return taSubscription, nil
	}
	//no need to check this since offer can change, therefore changing the transaction

	if len(input.TransactionSignature) > 0 && (input.Commit == 0 || subscriberWallet.HasViewOnlyAccess(gc)) {
		txnHash, e := network.SubmitXdrWithSignature(client, subscriber.PrimarySigner, input.Transaction, input.TransactionSignature)
		if e != nil {
			log.Printf("[SubscribeToTokenizedAsset] error submitting tokenized asset subscription to blockchain [%+v] for %v: %v\n", taSubscription, subscriber.Username, e)

			err = e
			logDiscordFailedTokenizedAssetSubscription(fmt.Sprintf("Error submitting asset subscription [%+v] transaction: %s", input, err.Error()))
			if strings.Contains(err.Error(), "liquid") {
				destAsset := os.Getenv("NATIVE_ASSET_CODE")
				sourceAsset := os.Getenv("NATIVE_ASSET_CODE")
				if len(swapInfo.SourceAssetCode) > 0 {
					sourceAsset = swapInfo.SourceAssetCode
				}
				if len(swapInfo.DestinationAssetCode) > 0 {
					destAsset = swapInfo.DestinationAssetCode
				}

				_, b, _ := gc.GetAvalableMarketQuantity(swapInfo.SourceAssetCode, swapInfo.SourceAssetIssuer, swapInfo.DestinationAssetCode, swapInfo.DestinationAssetIssuer)

				emsg := fmt.Sprintf("There is no %v market to exchange for your %v at this time. Please try again later or reduce the quantity of %v to try again.", destAsset, sourceAsset, sourceAsset)
				if b != "0" {
					emsg = fmt.Sprintf("There is only %v %v to exchange for your %v at this time. Please reduce the quantity of %v to try again.", b, destAsset, sourceAsset, sourceAsset)
				}

				err = &tErrors.CustomError{
					Param:      "destinationAssetCode",
					Err:        "error-low-liquidity",
					ErrMessage: emsg,
				}

				return
			}
		}
		input.TransactionID = txnHash
		taSubscription.TransactionID = txnHash
		dbTX.Omit(clause.Associations).Save(&taSubscription)
		dbTX.Commit()
		subscriberWallet.InvalidateUserCache(gc)
		return taSubscription, nil

	}
	//multi Party
	if input.Multiparty == 1 {
		input.TransactionID = "PENDING_AUTH"
		taSubscription.TransactionID = "PENDING_AUTH"

		id := uuid.NewString()

		description := fmt.Sprintf("Buying Tokenized asset [%v]\n Amount:%v,\n Getting Apprx:%v %v", *ta.AssetName, fmt.Sprintf("%v %v", input.Amount, *ta.AssetQuoteCurrency), input.SwappedEstimate, *ta.AssetCode)
		if len(swapInfo.Memo) > 0 {
			description = fmt.Sprintf("%v\nMemo: %v", description, input.Memo)

		}
		if len(input.Messages) > 0 {
			var msgs string
			for i, m := range input.Messages {
				msgs = m
				if i < len(input.Messages)-1 {
					msgs = fmt.Sprintf("%s\n", msgs)
				}
			}
			description = fmt.Sprintf("%v\nMessages: %v", description, msgs)

		}
		input.ReturnedDescription = description
		transactionByte, _ := json.Marshal(*input)
		transactionStr := string(transactionByte)
		pendingAuth := userModels.PendingAuth{
			ID:                       id,
			Initiator:                subscriber.Username,
			InitiatorSignerPublicKey: subscriber.PrimarySigner,
			WalletPublicKey:          subscriberWallet.ID,
			TransactionType:          "ASSET SUBSCRIPTION",
			Description:              description,
			TransactionSource:        input.TransactionSource,
			ApprovalsNeeded:          subscriberWallet.NumberOfApprovalsNeeded,
			TransactionXdr:           input.Transaction,
			TransactionInfoStr:       &transactionStr,
		}
		//save and commit this to database
		e := dbTX.Omit(clause.Associations).Create(&pendingAuth).Error
		if e != nil {
			log.Printf("[SubscribeToTokenizedAsset] Error saving asset subscription txn [%+v] transaction on pending auth table: %s\n", pendingAuth, e.Error())
			err = &tErrors.ErrorTemporaryServerError{}

		}

		dbTX.Commit()

		return taSubscription, nil
	}
	log.Println("[SubscribeToTokenizedAsset]UNKNOWN OPTION FOR ACTION")
	err = &tErrors.ErrorTemporaryServerError{}
	return

}

func generateAssetSubscriptionXdr(wallet *userModels.UserWallet, swapInfo *swapModels.SwapSendInfo, gc *sharedconfig.GlobalConfig) (string, error) {
	// baseReserve := network.GetBlockchainBaseReserve()
	swapDestMin := network.GetBlockchainSwapDestinationMin()
	client := gc.BantuExpansionClient
	messages := make([]string, 0)
	nativeAssetCode := os.Getenv("NATIVE_ASSET_CODE")
	var err error
	var amountToSwap decimal.Decimal

	if amountToSwap, err = decimal.NewFromString(swapInfo.SourceAmount); err != nil {
		return "", &swapErrors.ErrorInvalidSwapAmount{}
	}
	// newAmountToSwap := amountToSwap.Truncate(7).String()

	var sourceAsset txnbuild.Asset = txnbuild.NativeAsset{}
	var destinationAsset txnbuild.Asset = txnbuild.NativeAsset{}

	if len(swapInfo.DestinationAssetCode) != 0 && !strings.EqualFold(swapInfo.DestinationAssetCode, nativeAssetCode) {

		destinationAsset = txnbuild.CreditAsset{Code: swapInfo.DestinationAssetCode, Issuer: swapInfo.DestinationAssetIssuer}
	}
	if len(swapInfo.SourceAssetCode) != 0 && !strings.EqualFold(swapInfo.SourceAssetCode, nativeAssetCode) {

		sourceAsset = txnbuild.CreditAsset{Code: swapInfo.SourceAssetCode, Issuer: swapInfo.SourceAssetIssuer}
	}

	swapInfo.Messages = messages
	var ops []txnbuild.Operation = make([]txnbuild.Operation, 0)
	chanAccount := <-gc.ChannelAccounts
	defer func(c *keypair.Full) {
		gc.ChannelAccounts <- c
	}(chanAccount)

	_, _, _, _, chanSourceAccount, _ := network.BlockchainAccountProperties(client, chanAccount.Address(), txnbuild.NativeAsset{})

	sourceAccountExists, _, sourceAccountNativeBalance, sourceAccountCustomBalance, sourceAccount, sourceAccountErr := network.BlockchainAccountProperties(client, wallet.ID, sourceAsset)
	var sourceAccountTrustsDestinationAsset bool
	if !destinationAsset.IsNative() {
		_, sourceAccountTrustsDestinationAsset, _, _, _, _ = network.BlockchainAccountProperties(client, wallet.ID, destinationAsset)

	}

	if sourceAccountErr != nil {
		return "", sourceAccountErr
	}

	if !sourceAccountExists {
		return "", &tErrors.ErrorUnderfundedAccount{}
	}

	if !destinationAsset.IsNative() {

		if !sourceAccountTrustsDestinationAsset {

			//establish trustline
			ops = append(ops, &txnbuild.ChangeTrust{
				Line:          txnbuild.ChangeTrustAssetWrapper{Asset: destinationAsset},
				Limit:         gc.TokenLimitAsString(),
				SourceAccount: wallet.ID,
			})
			// allow trust from issuer to destination wallet
			ops = append(ops, &txnbuild.SetTrustLineFlags{
				Trustor:       wallet.ID,
				Asset:         txnbuild.CreditAsset{Code: swapInfo.DestinationAssetCode, Issuer: swapInfo.DestinationAssetIssuer},
				SetFlags:      []txnbuild.TrustLineFlag{txnbuild.TrustLineAuthorized},
				SourceAccount: swapInfo.DestinationAssetIssuer,
			})

		}
	}

	log.Printf("[generateAssetSubscriptionXdr]obtained source account balance:\n%v balance is %v\n%v balance is %v\n", nativeAssetCode, sourceAccountNativeBalance, sourceAsset.GetCode(), sourceAccountCustomBalance)

	if sourceAccountCustomBalance.LessThan(amountToSwap) {
		return "", &tErrors.ErrorUnderfundedAccount{Detail: fmt.Sprintf("Not enough funds. Needs Extra %v %v or you reduce same from the amount you want to swap.", (amountToSwap).Sub(sourceAccountCustomBalance), sourceAsset.GetCode())}
	}

	//get sendPath
	//using DestinationAccount will get paths to all assets in the destination account.
	//using destinationAssets gets path to only the asset
	destAsset := ""
	if !destinationAsset.IsNative() {
		destAsset = fmt.Sprintf("%s:%s", swapInfo.DestinationAssetCode, swapInfo.DestinationAssetIssuer)
	}
	pathInput := swapModels.SwapSendPathInput{
		DestinationAssets: destAsset,
		SourceAssetCode:   swapInfo.SourceAssetCode,
		SourceAssetIssuer: swapInfo.SourceAssetIssuer,
		SourceAmount:      swapInfo.SwapAmount,
	}
	path, swappedEstimate, err := GetStrictSendPaths(pathInput, client)
	if err != nil {
		log.Println("[generateAssetSubscriptionXdr]error fetching valid swap Path ", err)
		return "", err
	}

	//native asset
	ops = append(ops, &txnbuild.PathPaymentStrictSend{
		SendAsset:     sourceAsset,
		SendAmount:    swapInfo.SwapAmount,
		Destination:   wallet.ID,
		DestAsset:     destinationAsset,
		DestMin:       swapDestMin.String(),
		Path:          path,
		SourceAccount: wallet.ID,
	})

	// Construct the transaction that holds the operations to execute on the network
	var memoSAC, memoDAC string
	memoSAC = swapInfo.SourceAssetCode
	memoDAC = swapInfo.DestinationAssetCode

	memo := fmt.Sprintf("%v>%v", memoSAC, memoDAC)
	log.Println("[generateAssetSubscriptionXdr] Memo:", memo)
	swapInfo.Memo = memo

	var tx *txnbuild.Transaction
	// Construct the transaction that holds the operations to execute on the network
	if swapInfo.Multiparty == 1 {
		swapInfo.TransactionSource = chanSourceAccount.AccountID
		tx, err = txnbuild.NewTransaction(
			txnbuild.TransactionParams{
				SourceAccount:        chanSourceAccount,
				IncrementSequenceNum: true,
				Operations:           ops,
				BaseFee:              2000,
				Preconditions: txnbuild.Preconditions{
					TimeBounds: txnbuild.NewInfiniteTimeout(),
				},
				Memo: txnbuild.MemoText(memo),
			},
		)
	} else {
		tx, err = txnbuild.NewTransaction(
			txnbuild.TransactionParams{
				SourceAccount:        sourceAccount,
				IncrementSequenceNum: true,
				Operations:           ops,
				BaseFee:              2000,
				Preconditions: txnbuild.Preconditions{
					TimeBounds: txnbuild.NewInfiniteTimeout(),
				},
				Memo: txnbuild.MemoText(memo),
			},
		)
	}

	if err != nil {
		log.Println("[generateAssetSubscriptionXdr] error constructing transaction ", err)
		return "", err
	}

	if swapInfo.Multiparty == 1 {

		tx, err = tx.Sign(network.GetBlockchainNetworkPassPhrase(), chanAccount)

		if err != nil {
			log.Println("[generateAssetSubscriptionXdr] error signing transaction with channelAccount key ", err)
			return "", &tErrors.ErrorTemporaryServerError{}
		}
	}

	if !sourceAccountTrustsDestinationAsset {
		log.Printf("[generateAssetSubscriptionXdr] <<<<<<<<<<<<<<<<<<<<<<<<<<<< signing transaction with issuer key>>>>>>>>>>>>>>>>>>>>>>>>:[%v]\n\n", destinationAsset)
		//get atprofile
		var tokenizationIssuerProfileWallet string

		if len(os.Getenv("TOKENIZATION_ISSUING_PROFILE_WALLET")) > 1 {
			tokenizationIssuerProfileWallet = strings.TrimSpace(os.Getenv("TOKENIZATION_ISSUING_PROFILE_WALLET"))
		}

		tokenizationIssuerProfileWalletKP := keypair.MustParseFull(tokenizationIssuerProfileWallet)

		tx, err = tx.Sign(network.GetBlockchainNetworkPassPhrase(), tokenizationIssuerProfileWalletKP)
		if err != nil {
			log.Println("[generateAssetSubscriptionXdr] error signing transaction with issuer key to authorize trustline", err)
			return "", &tErrors.ErrorTemporaryServerError{}
		}
	}

	var xdrBase64 string

	xdrBase64, err = tx.Base64()
	if err != nil {
		log.Println("[generateAssetSubscriptionXdr] error getting txn base64", err)
		return "", err
	}
	swapInfo.Memo = memo
	swapInfo.Messages = messages
	swapInfo.SwappedEstimate = swappedEstimate
	return xdrBase64, nil
}

// GetStrictSendPaths gets Strict Send Paths for Strict Send Path Payment request
func GetStrictSendPaths(pathInput swapModels.SwapSendPathInput, client *horizonclient.Client) (paths []txnbuild.Asset, swappedEstimate string, err error) {
	discord.WebhookURL = "https://discord.com/api/webhooks/824381163367170058/75RxS1LzWA800hWereJJumw"
	if len(os.Getenv("EXPANSION_NETWORK_ERROR_WEBHOOK")) > 50 {
		discord.WebhookURL = os.Getenv("EXPANSION_NETWORK_ERROR_WEBHOOK")
	}
	var swapPaths horizon.PathsPage
	paths = make([]txnbuild.Asset, 0)
	var sourceAssetType horizonclient.AssetType
	if len(pathInput.SourceAssetIssuer) == 0 {
		sourceAssetType = horizonclient.AssetTypeNative
		pathInput.SourceAssetCode = ""
		pathInput.SourceAssetIssuer = ""
	} else if len(pathInput.SourceAssetCode) < 5 && len(pathInput.SourceAssetIssuer) == 56 {
		sourceAssetType = horizonclient.AssetType4
	} else if len(pathInput.SourceAssetCode) > 4 && len(pathInput.SourceAssetCode) <= 12 && len(pathInput.SourceAssetIssuer) == 56 {
		sourceAssetType = horizonclient.AssetType12
	}
	if pathInput.DestinationAccount != "" {
		pathInput.DestinationAssets = ""
	}

	if pathInput.DestinationAssets == "" && pathInput.DestinationAccount == "" {
		pathInput.DestinationAssets = "native"
	}
	if pathInput.DestinationAssets == "native" {
		pathInput.DestinationAccount = ""
	}

	if pathInput.DestinationAssets != "" {
		pathInput.DestinationAccount = ""
	}
	sspr := horizonclient.StrictSendPathsRequest{
		DestinationAccount: pathInput.DestinationAccount,
		DestinationAssets:  pathInput.DestinationAssets,
		SourceAssetType:    sourceAssetType,
		SourceAssetCode:    pathInput.SourceAssetCode,
		SourceAssetIssuer:  pathInput.SourceAssetIssuer,
		SourceAmount:       pathInput.SourceAmount,
	}

	swapPaths, err = client.StrictSendPaths(sspr)

	if err != nil {
		if strings.Contains(err.Error(), "timeout") || strings.Contains(err.Error(), "handshake") || strings.Contains(err.Error(), "read tcp") || strings.Contains(err.Error(), "connection reset by peer") || strings.Contains(err.Error(), "dial tcp") || strings.Contains(err.Error(), "no such host") {
			log.Println("#######################@@@@@@@@@@@@@@@@[client.StrictSendPathsErr] expansion connection problem:", err)
			discord.Say(fmt.Sprintf("[getStrictSendPaths] error connecting to expansion service: %v\nSwapPathRequest: %+v", err, sspr))

			return paths, "", &tErrors.ErrorTemporaryServerError{}
		}
		if strings.Contains(err.Error(), "liquid") {
			destAsset := os.Getenv("NATIVE_ASSET_CODE")
			sourceAsset := os.Getenv("NATIVE_ASSET_CODE")
			if len(pathInput.SourceAssetCode) > 0 {
				sourceAsset = pathInput.SourceAssetCode
			}
			if pathInput.DestinationAssets != "native" {
				destAsset = strings.Split(pathInput.DestinationAssets, ":")[0]
			}
			return paths, "", &tErrors.CustomError{
				Param:      "destinationAssetCode",
				Err:        "error-low-liquidity",
				ErrMessage: fmt.Sprintf("There is not enough %v market to exchange for your %v at this time. Please try again later or reduce the quantity of %v to try again.", destAsset, sourceAsset, sourceAsset),
			}
		}
		horizonException, ok := err.(*horizonclient.Error)

		if ok {

			extraErrors := horizonException.Problem.Extras

			for key, val := range extraErrors {
				log.Printf("Extras: %v is %v\n", key, val)
			}

			resultCodes, e := horizonException.ResultCodes()
			if e != nil {
				log.Println("[client.StrictSendPathsErr] Error getting result codes:", e)

				destAsset := os.Getenv("NATIVE_ASSET_CODE")
				sourceAsset := os.Getenv("NATIVE_ASSET_CODE")
				if len(pathInput.SourceAssetCode) > 0 {
					sourceAsset = pathInput.SourceAssetCode
				}
				if pathInput.DestinationAssets != "native" {
					destAsset = strings.Split(pathInput.DestinationAssets, ":")[0]
				}
				return paths, "", &tErrors.CustomError{
					Param:      "destinationAssetCode",
					Err:        "error-low-liquidity",
					ErrMessage: fmt.Sprintf("There is not enough %v market to exchange for your %v at this time. Please try again later or reduce the quantity of %v to try again.", destAsset, sourceAsset, sourceAsset),
				}

			}

			for key, val := range resultCodes.OperationCodes {
				log.Printf("Result code: %v is %v\n", key, val)
			}
			discord.Say(fmt.Sprintf("[getStrictSendPaths] error submitting: %v\nSwapPathRequest: %+v\nResultCodes: %+v", err, sspr, resultCodes))

		}
		log.Println("[client.StrictSendPathsErr] Error submitting:", err)
		if strings.Contains(err.Error(), "liquid") {
			destAsset := os.Getenv("NATIVE_ASSET_CODE")
			sourceAsset := os.Getenv("NATIVE_ASSET_CODE")
			if len(pathInput.SourceAssetCode) > 0 {
				sourceAsset = pathInput.SourceAssetCode
			}
			if pathInput.DestinationAssets != "native" {
				destAsset = strings.Split(pathInput.DestinationAssets, ":")[0]
			}
			return paths, "", &tErrors.CustomError{
				Param:      "destinationAssetCode",
				Err:        "error-low-liquidity",
				ErrMessage: fmt.Sprintf("There is not enough %v market to exchange for your %v at this time. Please try again later or reduce the quantity of %v to try again.", destAsset, sourceAsset, sourceAsset),
			}
		}
		return paths, "", &tErrors.ErrorTemporaryServerError{}

	}
	// discord.Say(fmt.Sprintf("[getStrictSendPaths] swapPaths: %+v\nRequestParams: %+v", swapPaths, sspr))
	// log.Printf("[getStrictSendPaths] swapPaths: %+v\n", swapPaths)
	if len(swapPaths.Embedded.Records) == 0 {
		destAsset := os.Getenv("NATIVE_ASSET_CODE")
		sourceAsset := os.Getenv("NATIVE_ASSET_CODE")
		if len(pathInput.SourceAssetCode) > 0 {
			sourceAsset = pathInput.SourceAssetCode
		}
		if pathInput.DestinationAssets != "native" {
			destAsset = strings.Split(pathInput.DestinationAssets, ":")[0]
		}
		return paths, "", &tErrors.CustomError{
			Param:      "destinationAssetCode",
			Err:        "error-low-liquidity",
			ErrMessage: fmt.Sprintf("There is not enough %v market to exchange for your %v at this time. Please try again later or reduce the quantity of %v to try again.", destAsset, sourceAsset, sourceAsset),
		}
	}
	destAmountDec, _ := decimal.NewFromString(swapPaths.Embedded.Records[0].DestinationAmount)
	if destAmountDec.LessThan(network.GetBlockchainSwapDestinationMin()) {
		return paths, "", &swapErrors.ErrorSwapAmountTooSmall{}
	}

	bestPath := swapPaths.Embedded.Records[0]
	//build assets
	swappedEstimate = bestPath.DestinationAmount

	for _, v := range bestPath.Path {
		if len(v.Issuer) == 0 {
			paths = append(paths, txnbuild.NativeAsset{})
		} else {
			paths = append(paths, txnbuild.CreditAsset{Code: v.Code, Issuer: v.Issuer})
		}

	}

	return paths, swappedEstimate, nil
}

func ExpressInterest(subscriber *userModels.User, ta *userModels.TokenizedAsset, input *userModels.ExpressionOfInterestInput, gc *sharedconfig.GlobalConfig) (expressedInterest userModels.ExpressionOfInterest, err error) {

	if ta.AssetTokenizationStatus != 4 {

		log.Printf("[ExpressInterest] Error interests cannot be expressed on tokenizations with this status: %v\n", ta.ID)
		err = &tErrors.CustomError{Param: "issuingWalletPublicKey", Err: "error-invalid-request", ErrMessage: "Only projects that are market ready can accept expression of interests."}
		return

	}
	expressedInterest, _ = ta.GetExpressedInterestByUsername(subscriber.Username, gc)
	expressedInterest.UpdateExpressionOfInterestFromInput(subscriber.Username, input, ta, gc)
	e := gc.DB.Omit(clause.Associations).Save(&expressedInterest).Error
	if e != nil {
		log.Printf("[ExpressInterest] error saving expression of interest to database  [%+v] for %v: %v\n", expressedInterest, subscriber.Username, e)

		err = &tErrors.ErrorTemporaryServerError{}

	}
	return
}

func UpdateTokenizedAssetFromInput(t *userModels.TokenizedAsset, ti *userModels.TokenizedAssetJSONInput, gc *sharedconfig.GlobalConfig) userModels.TokenizedAsset {
	return t.UpdateTokenizedAssetFromInput(ti, gc)
}

func generateMintRegulatedTokenizedAssetXdr(t *userModels.TokenizedAsset, gc *sharedconfig.GlobalConfig) (xdrbase64, transactionSource string, messages []string, issuingWallet userModels.UserWallet, err error) {
	replacer := strings.NewReplacer("\r", "", "\n", "", " ", "")

	if decimal.NewFromFloat(t.NumberOfTokenToBeIssued).GreaterThan(gc.TokenLimitAsDecimal()) {
		err = &tErrors.CustomError{Param: "numberOfTokenToBeIssued", Err: "error-token-limit-exceeded", ErrMessage: fmt.Sprintf("Issued Number of tokens cannot exceed %v", gc.TokenLimitAsString())}
		return
	}
	client := gc.BantuExpansionClient
	feeWallet := keypair.MustParseFull(os.Getenv("TOKENIZATION_FEE_WALLET"))
	var ops []txnbuild.Operation = make([]txnbuild.Operation, 0)
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

	if t.MintingApprovers == nil {
		// set default
		u := t.GetMintingApproversInCSV(gc)
		t.MintingApprovers = &u
	}

	if t.MintingInitators == nil {
		//set default
		u := t.GetMintingInitiatorsInCSV(gc)
		t.MintingInitators = &u
	}
	trimmedApprovers := replacer.Replace(*t.MintingApprovers)
	t.MintingApprovers = &trimmedApprovers

	trimmedInitiators := replacer.Replace(*t.MintingInitators)
	t.MintingInitators = &trimmedInitiators
	var aps []string
	var inits []string
	aps = strings.Split(*t.MintingApprovers, ",")
	inits = strings.Split(*t.MintingInitators, ",")

	if len(aps) == 0 || len(inits) == 0 {
		err = &tErrors.CustomError{
			Param:      "numberOfApprovers",
			Err:        "error-no-approver-or-initiator-specified",
			ErrMessage: "No approvers /initiators specified",
			Code:       404,
		}
		return "", "", messages, issuingWallet, err
	}

	if len(aps) < 4 {
		err = &tErrors.CustomError{
			Param:      "numberOfApprovers",
			Err:        "error-approvers-less-than-4",
			ErrMessage: fmt.Sprintf("%v approvers specified. Requires minimum of 4", len(aps)),
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
			TargetUsername: v,
			Permission:     "APPROVER",
		})
	}

	for _, v := range inits {
		permInfo = append(permInfo, userModels.WalletPermissionInfo{
			TargetUsername: v,
			Permission:     "INITIATOR",
		})
	}
	var p userModels.UserWalletSharedAccessInfo

	p = userModels.UserWalletSharedAccessInfo{
		WalletPublicKey:         *t.IssuingWalletPublicKey,
		NumberOfApprovalsNeeded: len(aps) - 2,
		Permissions:             permInfo,
	}
	if err = checkDistributionWalletHasQuoteCurrencyAuthorization(*t.AssetQuoteCurrency, &distributionWallet, gc); err != nil {
		return "", "", messages, issuingWallet, err
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
		//get fresh records.
		tokenizationIssuerUser.InvalidateUserCache(gc)
		tokenizationIssuerUser, _ = userModels.Username(tokenizationIssuerProfile).GetFullUser(gc.DB, gc)
		issuingWallet, e = userModels.UserWalletID(*t.IssuingWalletPublicKey).GetWallet(gc.DB, gc)
		if e != nil {
			err = &tErrors.CustomError{Param: "issuingPublicKey", Err: "error-invalid-issuer", ErrMessage: "error validating issuing wallet."}
			return
		}
		//second submission to blockchain
		_, err = CreateSharedWalletAccess(&tokenizationIssuerUser, &tokenizationIssuerUser, &issuingWallet, &p, gc)
		if err != nil {
			log.Printf("[generateMintRegulatedTokenizedAssetXdr.CreateSharedWalletAccess: stage 2] Error creating shared access on wallet [%v], err: %v\n", issuingWallet.ID, err)
			return
		}

		log.Printf("[generateMintRegulatedTokenizedAssetXdr.CreateSharedWalletAccess] Succesfully Created shared access on issuing wallet [%v], txID: %v\n", issuingWallet.ID, p.TransactionID)

		//get fresh records.
		tokenizationIssuerUser.InvalidateUserCache(gc)
		tokenizationIssuerUser, _ = userModels.Username(tokenizationIssuerProfile).GetFullUser(gc.DB, gc)
		issuingWallet, _ = userModels.UserWalletID(*t.IssuingWalletPublicKey).GetWallet(gc.DB, gc)

	}
	// create distributionWallet trustline to tokenized asset
	ops = append(ops, &txnbuild.ChangeTrust{
		Line:          txnbuild.ChangeTrustAssetWrapper{Asset: txnbuild.CreditAsset{Code: *t.AssetCode, Issuer: *t.IssuingWalletPublicKey}},
		Limit:         gc.TokenLimitAsString(),
		SourceAccount: distributionWallet.ID,
	})

	//create trusline on fee wallet

	//get fee wallet
	ops = append(ops, &txnbuild.ChangeTrust{
		Line:          txnbuild.ChangeTrustAssetWrapper{Asset: txnbuild.CreditAsset{Code: *t.AssetCode, Issuer: *t.IssuingWalletPublicKey}},
		Limit:         gc.TokenLimitAsString(),
		SourceAccount: feeWallet.Address(),
	})

	// allow trust from issuer to fee wallet
	ops = append(ops, &txnbuild.SetTrustLineFlags{
		Trustor:       feeWallet.Address(),
		Asset:         txnbuild.CreditAsset{Code: *t.AssetCode, Issuer: *t.IssuingWalletPublicKey},
		SetFlags:      []txnbuild.TrustLineFlag{txnbuild.TrustLineAuthorized},
		SourceAccount: *t.IssuingWalletPublicKey,
	})

	// allow trust from issuer to distribution wallet
	ops = append(ops, &txnbuild.SetTrustLineFlags{
		Trustor:       distributionWallet.ID,
		Asset:         txnbuild.CreditAsset{Code: *t.AssetCode, Issuer: *t.IssuingWalletPublicKey},
		SetFlags:      []txnbuild.TrustLineFlag{txnbuild.TrustLineAuthorized},
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
	feeInAssetPayment := &txnbuild.Payment{
		Destination:   feeWallet.Address(),
		Amount:        decimal.NewFromFloat(t.FeeInAsset).StringFixed(7),
		Asset:         txnbuild.CreditAsset{Code: *t.AssetCode, Issuer: *t.IssuingWalletPublicKey},
		SourceAccount: distributionWallet.ID,
	}
	if e := feeInAssetPayment.Validate(); e != nil {
		msg := fmt.Sprintf("[generateMintRegulatedTokenizedAssetXdr] FeeInAssetPayment Operation Failed validation: %v. Fee In Asset figure: %v. Fields: %+v", e, decimal.NewFromFloat(t.FeeInAsset).StringFixed(7), *feeInAssetPayment)
		gc.LogDiscordFailedRequest(msg)
		err = &tErrors.CustomError{
			Param:      "IssuingWalletPublicKey",
			Err:        "error-could-not-approve-tokenization",
			ErrMessage: "Could not approve tokenization. A fee Payment operation could not pass validation.",
			Code:       404,
		}
		return "", "", messages, issuingWallet, err
	}

	ops = append(ops, feeInAssetPayment)
	//make market
	fraction := decimal.NewFromFloat(t.PricePerToken).Rat()
	d := int32(fraction.Denom().Int64())
	n := int32(fraction.Num().Int64())
	xdrPrice := xdr.Price{N: xdr.Int32(n), D: xdr.Int32(d)}
	marketOffer := &txnbuild.ManageSellOffer{
		Buying:        txnbuild.CreditAsset{Code: quoteCurrency.AssetCode, Issuer: quoteCurrency.AssetIssuer},
		Amount:        decimal.NewFromFloat(t.MaxNumberOfTokenAvailableForSale).StringFixed(7),
		Selling:       txnbuild.CreditAsset{Code: *t.AssetCode, Issuer: *t.IssuingWalletPublicKey},
		Price:         xdrPrice,
		SourceAccount: distributionWallet.ID,
	}
	if e := marketOffer.Validate(); e != nil {
		msg := fmt.Sprintf("[generateMintRegulatedTokenizedAssetXdr] marketOffer Operation Failed validation: %v. Offer figure: %v. Fields: %+v", e, t.MaxNumberOfTokenAvailableForSale, *marketOffer)
		gc.LogDiscordFailedRequest(msg)
		err = &tErrors.CustomError{
			Param:      "IssuingWalletPublicKey",
			Err:        "error-could-not-approve-tokenization",
			ErrMessage: "Could not approve tokenization. Market offer operation could not pass validation.",
			Code:       404,
		}
		return "", "", messages, issuingWallet, err
	}

	ops = append(ops, marketOffer)
	msg := fmt.Sprintf("[generateMintRegulatedTokenizedAssetXdr] Transaction to mint %v units (selling %v units) of %v @ %v %v generated. Fraction: %v. N: %v, D: %v. XDR Price: %+v", decimal.NewFromFloat(t.NumberOfTokenToBeIssued).StringFixed(7), decimal.NewFromFloat(t.MaxNumberOfTokenAvailableForSale).StringFixed(7), *t.AssetCode, decimal.NewFromFloat(t.PricePerToken).StringFixed(7), *t.AssetQuoteCurrency, fraction.String(), decimal.NewFromInt32(n).String(), decimal.NewFromInt32(d).String(), xdrPrice)
	gc.LogDiscordFailedRequest(msg)
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

	tx, err = tx.Sign(network.GetBlockchainNetworkPassPhrase(), chanAccount, feeWallet)

	if err != nil {
		log.Println("[generateMintRegulatedTokenizedAssetXdr] error signing transaction with channelAccount & fee wallet key ", err)
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
	replacer := strings.NewReplacer("\r", "", "\n", "", " ", "")
	//referesh issuing wallet profile
	userModels.Username(os.Getenv("TOKENIZATION_ISSUING_PROFILE")).InvalidateUserCache(gc)
	ato, _, err = GetTokenizedAssetByID(tokenizationID, gc.DB)
	if err != nil {
		// error tokenization is already in progress
		log.Printf("[MintRegulatedTokenizedAsset] Error fetching tokenization with ID: %v\n", tokenizationID)
		return

	}

	if decimal.NewFromFloat(ato.NumberOfTokenToBeIssued).GreaterThan(gc.TokenLimitAsDecimal()) {
		err = &tErrors.CustomError{Param: "numberOfTokenToBeIssued", Err: "error-token-limit-exceeded", ErrMessage: fmt.Sprintf("Issued Number of tokens cannot exceed %v", gc.TokenLimitAsString())}
		return
	}
	if ato.AssetCountryLocation == nil {
		err = &tErrors.CustomError{Param: "Id", Err: "error-asset-location-country-not-found", ErrMessage: "You must specify the asset country of location."}
		return
	}
	if len(*ato.AssetCountryLocation) != 2 {
		err = &tErrors.CustomError{Param: "Id", Err: "error-asset-location-country-not-found", ErrMessage: "You must specify the asset country of location in the formart: NG, SA, UK"}
		return
	}
	if ato.ExemptedCountries != nil {
		exc := strings.Split(replacer.Replace(*ato.ExemptedCountries), ",")
		for _, countryCode := range exc {
			if len(countryCode) != 2 {
				err = &tErrors.CustomError{Param: "Id", Err: "error-asset-exempted-country-invalid", ErrMessage: "You must specify the exempted country in the formart: NG, SA, UK"}
				return
			}
		}

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

		if (ato.SecApproval == 0 || ato.SecApprovalIdNumber == nil) && *ato.OfferingType == "PUBLIC" {
			log.Printf("[MintRegulatedTokenizedAsset] Error No complete DD information yet. Misssing SEC approval for public offering %v\n", ato.ID)

			err = &tErrors.CustomError{
				Param:      "publicKey",
				Err:        "error-incomplete-data",
				ErrMessage: "Asset Due Duligence information (SEC Aproval Number) must be provided for public offerings.",
				Code:       404,
			}
			return
		}
		if ato.ClosedGroupID == nil && *ato.OfferingType == "PRIVATE" {
			log.Printf("[MintRegulatedTokenizedAsset] Error No complete information yet. Misssing Closed Group for private offering %v\n", ato.ID)

			err = &tErrors.CustomError{
				Param:      "publicKey",
				Err:        "error-incomplete-data",
				ErrMessage: "Asset Closed Group information must be provided for private offerings.",
				Code:       404,
			}
			return
		}

		if ato.AssetCode == nil || ato.AssetName == nil || ato.AssetDescription == nil || ato.AssetLogo == nil {
			log.Printf("[MintRegulatedTokenizedAsset] Error No complete asset token information yet. %v\n", ato.ID)

			err = &tErrors.CustomError{
				Param:      "publicKey",
				Err:        "error-incomplete-data",
				ErrMessage: "Asset token information (Name,Code,Description,Logo) must be provided before Minting can be done.",
				Code:       404,
			}
			return
		}
		if ato.AssetTokenizationStatus < 3 {
			log.Printf("[MintRegulatedTokenizedAsset] Error Process has not reached stage to approve yet. %v\n", ato.ID)

			err = &tErrors.CustomError{
				Param:      "publicKey",
				Err:        "error-status-too-low",
				ErrMessage: "Process Status still too low for approval. Fee Payment confirmation not yet done.",
				Code:       404,
			}
			return
		}
		if ato.AssetTokenizationStatus > 3 {
			log.Printf("[MintRegulatedTokenizedAsset] Error Process has passed stage to mint. %v\n", ato.ID)

			err = &tErrors.CustomError{
				Param:      "publicKey",
				Err:        "error-status-too-high",
				ErrMessage: "Project status has exceeded the stage of minting.",
				Code:       404,
			}
			return
		}

		if ato.IssuingWalletPublicKey == nil {
			log.Printf("[MintRegulatedTokenizedAsset] Error no issuing wallet assigned. Returning this to earlier status. %v\n", ato.ID)
			ato.AssetTokenizationStatus = 2
			gc.DB.Omit(clause.Associations).Save(&ato)
			err = &tErrors.CustomError{
				Param:      "publicKey",
				Err:        "error-no-issuing-wallet",
				ErrMessage: "Issuing Wallet not assigned. Tokenization has been Returned to approproate status.",
				Code:       404,
			}
			return
		}
	}

	xdrBase64, transactionSource, messages, issuingWallet, err := generateMintRegulatedTokenizedAssetXdr(&ato, gc)

	if err != nil {

		return
	}

	dbTX := gc.DB.Begin()
	defer dbTX.Rollback()

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
	e := dbTX.Omit(clause.Associations).Create(&pendingAuth).Error
	if e != nil {
		log.Printf("[MintRegulatedTokenizedAsset] Error saving txn [%+v] on pending auth table: %s\n", pendingAuth, e.Error())
		err = &tErrors.ErrorTemporaryServerError{}
		return
	}
	ato.AssetTokenizationStatus = 4
	// ato.TokenizationTransaction = &xdrBase64
	e = dbTX.Omit(clause.Associations).Save(&ato).Error
	if e != nil {
		log.Printf("[MintRegulatedTokenizedAsset] Error saving txn on tokenizedAsset table: %s\n", e.Error())
		err = &tErrors.ErrorTemporaryServerError{}
		return
	}
	dbTX.Commit()
	//log message
	fraction := decimal.NewFromFloat(ato.PricePerToken).Rat()
	d := int32(fraction.Denom().Int64())
	n := int32(fraction.Num().Int64())
	msg := fmt.Sprintf("Minting of %v units (selling %v units) of %v @ %v %v submitted. Fraction: %v. N: %v, D: %v", decimal.NewFromFloat(ato.NumberOfTokenToBeIssued).StringFixed(7), decimal.NewFromFloat(ato.MaxNumberOfTokenAvailableForSale).StringFixed(7), *ato.AssetCode, decimal.NewFromFloat(ato.PricePerToken).StringFixed(7), *ato.AssetQuoteCurrency, fraction.String(), decimal.NewFromInt32(n).String(), decimal.NewFromInt32(d).String())
	gc.LogDiscordFailedRequest(msg)
	return ato, nil

}
func GetPostTokenizationTrustlineCandidates(gc *sharedconfig.GlobalConfig) (tas []userModels.PostTokenizationTrustlineCandidate) {
	tas = make([]userModels.PostTokenizationTrustlineCandidate, 0)
	gc.DB.Find(&tas)
	return tas
}

func ProcessPostTokenizationTrustline(gc *sharedconfig.GlobalConfig) {
	//get all candidates.
	candidates := GetPostTokenizationTrustlineCandidates(gc)
	for _, candidate := range candidates {
		log.Printf("[ProcessPostTokenizationTrustline] starting to process candidate [%+v]\n", candidate)

		//get untrusted assets
		uts := candidate.GetUntrustedTokenizedAssets(gc)
		if len(uts) == 0 {
			continue

		}
		log.Printf("[ProcessPostTokenizationTrustline] starting to process candidate's UNTRUSTED ASSETS [%v]\n", uts)

		// get wallet
		wallet, err := userModels.UserWalletID(candidate.PublicKey).GetWallet(gc.DB, gc)
		if err != nil {
			log.Printf("[ProcessPostTokenizationTrustline] error getting candidate wallet info %v, err: %v\n", candidate.PublicKey, err)
			gc.LogDiscordFailedRequest(fmt.Sprintf("[ProcessPostTokenizationTrustline] error getting candidate wallet info %v, err: %v", candidate.PublicKey, err))
			continue
		}
		signerUser, err := userModels.Username("tinitiator").GetSimpleUser(gc.DB, gc)
		if err != nil {
			log.Printf("[ProcessPostTokenizationTrustline] error getting tinitator user info %v, err: %v\n", candidate.PublicKey, err)
			gc.LogDiscordFailedRequest(fmt.Sprintf("[ProcessPostTokenizationTrustline] error getting  tinitator user info %v, err: %v", candidate.PublicKey, err))
			continue
		}

		//use the list to now send
		for _, ut := range uts {
			data := strings.Split(ut, ":")
			assetCode, assetIssuer := data[0], data[1]
			// prepare trustline for it.

			trustLineInfo := userModels.Trustline{
				AssetCode:   assetCode,
				AssetIssuer: assetIssuer,
			}

			//call the func to process the trustline

			returnedTrustLineInfo, err := TrustAsset(&signerUser, &wallet, &trustLineInfo, gc)
			if err != nil {
				if err.Error() != "error-duplicate-operation-exists" {
					log.Printf("[ProcessPostTokenizationTrustline] error executing first trustline call command %v, err: %v\nTrustLineInfo: [%+v]", candidate.PublicKey, err, trustLineInfo)
					gc.LogDiscordFailedRequest(fmt.Sprintf("[ProcessPostTokenizationTrustline] error executing first trustline call command %v, err: %v\nTrustLineInfo: [%+v]", candidate.PublicKey, err, trustLineInfo))
				}
				continue
			}

			//go ahead to initiate the second call with commit
			returnedTrustLineInfo.Commit = 1

			returnedTrustLineInfo, err = TrustAsset(&signerUser, &wallet, returnedTrustLineInfo, gc)
			if err != nil {
				if err.Error() != "error-duplicate-operation-exists" {
					log.Printf("[ProcessPostTokenizationTrustline] error executing 2nd trustline call command %v,err: %v\nTrustLineInfo: [%+v]", candidate.PublicKey, err, trustLineInfo)
					gc.LogDiscordFailedRequest(fmt.Sprintf("[ProcessPostTokenizationTrustline] error executing 2nd trustline call command %v, err: %v\nTrustLineInfo: [%+v]", candidate.PublicKey, err, trustLineInfo))

				}
				continue
			}

			if returnedTrustLineInfo.TransactionID == "PENDING_AUTH" {
				//start push notificationMessage
				notificationList := make(map[string]string)
				permissionList := wallet.Permissions
				for _, v := range permissionList {
					u, e := userModels.Username(v.TargetUsername).GetSimpleUser(gc.DB, gc)
					if e != nil {
						continue
					}
					if u.PushNotificationToken == nil {
						continue
					}

					if _, ok := notificationList[*u.PushNotificationToken]; ok {
						continue
					}

					dataPayload := make(map[string]string)
					dataPayload["route"] = "pendingApproval"

					u.SendPushMessage(fmt.Sprintf("%v opt-in request from %v!", trustLineInfo.AssetCode, wallet.Alias), fmt.Sprintf("Request: %v", returnedTrustLineInfo.ReturnedDescription), "", dataPayload, gc)
					notificationList[*u.PushNotificationToken] = v.TargetUsername

				}
			}

		}
	}
}

func generateTokenizationFeeXdr(wallet *userModels.UserWallet, ato *userModels.TokenizedAsset, taInput *userModels.ConfirmTokenizedAssetJSONInput, gc *sharedconfig.GlobalConfig) (txnBase64 string, err error) {

	countryConfig := userModels.CountryCode(*ato.AssetCountryLocation).GetConfig(gc)

	tfa := strings.Split(countryConfig.TokenizationApplicationFeeAsset, ":") //CODE:ISSUER
	feeAmount := decimal.NewFromFloat(countryConfig.TokenizationApplicationFee)
	feeWallet := os.Getenv("TOKENIZATION_APPLICATION_FEE_WALLET")
	feeWalletPK := keypair.MustParseFull(feeWallet)
	assetIssuer := tfa[1]
	assetCode := tfa[0]

	asset := txnbuild.CreditAsset{Code: assetCode, Issuer: assetIssuer}

	sourceAccountExists, _, _, assetAccountFeeBalance, sourceAccount, _ := network.BlockchainAccountProperties(gc.BantuExpansionClient, wallet.ID, asset)

	if !sourceAccountExists {
		return "", &tErrors.CustomError{Param: "publicKey", Err: "error-account-not-activated-on-blockchain", ErrMessage: "The Wallet public key is currently underfunded. Please send about 3XBN to it to activate it before you can perform this task", Code: http.StatusBadRequest}

	}
	if assetAccountFeeBalance.LessThan(feeAmount) {
		return "", &tErrors.CustomError{Param: "publicKey", Err: "error-wallet-underfunded", ErrMessage: fmt.Sprintf("The Wallet is currently underfunded. Please maintain min %v %v balance before you can perform this task", feeAmount.String(), tfa[0]), Code: http.StatusBadRequest}

	}

	var ops []txnbuild.Operation = make([]txnbuild.Operation, 0)
	ops = append(ops, &txnbuild.Payment{
		Destination:   feeWalletPK.Address(),
		Amount:        feeAmount.String(),
		Asset:         asset,
		SourceAccount: wallet.ID,
	})

	taInput.Messages = append(taInput.Messages, fmt.Sprintf("Application fee of %v %v will be charged to your wallet with alias [%v].", feeAmount.String(), assetCode, wallet.Alias))
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

func checkDistributionWalletHasQuoteCurrencyAuthorization(assetQuoteCurrency string, distributionWallet *userModels.UserWallet, gc *sharedconfig.GlobalConfig) (err error) {
	//check if it has naira asset balance and check if the naira asset has authorization

	ndab := strings.Split(os.Getenv("NAIRA_ASSET"), ":")
	nairaAsset := txnbuild.CreditAsset{Code: ndab[0], Issuer: ndab[1]}
	_, ntrusted, _, _, dwalletSource, _ := network.BlockchainAccountProperties(gc.BantuExpansionClient, distributionWallet.ID, nairaAsset)

	if strings.EqualFold(assetQuoteCurrency, ndab[0]) {
		//it has naira asset

		if !ntrusted {

			//adistributor wallet does not trust naura asset yet..
			log.Printf("[checkDistributionWalletHasQuoteCurrencyAuthorization] Distribution wallet [%v] does not yet accept %v\n", distributionWallet.Alias, ndab[0])

			err = &tErrors.CustomError{
				Param:      "IssuingWalletPublicKey",
				Err:        "error-asset-trustline-not-authorized",
				ErrMessage: fmt.Sprintf("Distribution wallet %v not yet authorized to hold %v", distributionWallet.Alias, ndab[0]),
				Code:       404,
			}
			return err
		}

		//get the naira asset balance
		for _, bal := range dwalletSource.Balances {
			if strings.EqualFold(ndab[0], bal.Code) {
				//naira bal....perform the critical check

				bantuAsset := userModels.BantuAsset{
					AssetCode:   ndab[0],
					AssetIssuer: ndab[1],
				}
				bcAsset, e := bantuAsset.GetBlockchainAssetProperty(gc)
				if e != nil {
					err = &tErrors.ErrorTemporaryServerError{}
					return err

				}
				if len(bcAsset.Code) == 0 {
					err = &tErrors.ErrorTemporaryServerError{}
					return err

				}

				if bcAsset.Flags.AuthRequired {
					if bal.IsAuthorized != nil {
						if !*bal.IsAuthorized {
							//authorization has not been given. abort process.
							log.Printf("[checkDistributionWalletHasQuoteCurrencyAuthorization] Error: %v not authorized on distribution wallet [%v]\n", distributionWallet.Alias)

							err = &tErrors.CustomError{
								Param:      "IssuingWalletPublicKey",
								Err:        "error-asset-trustline-not-authorized",
								ErrMessage: fmt.Sprintf("Distribution wallet %v not yet authorized to hold %v", distributionWallet.Alias, ndab[0]),
								Code:       404,
							}
							return err
						}
					} else {
						//authorization has not been given. abort process.
						log.Printf("[checkDistributionWalletHasQuoteCurrencyAuthorization] Error: %v not authorized on distribution wallet [%v]\n", distributionWallet.Alias)

						err = &tErrors.CustomError{
							Param:      "IssuingWalletPublicKey",
							Err:        "error-asset-trustline-not-authorized",
							ErrMessage: fmt.Sprintf("Distribution wallet %v not yet authorized to hold %v", distributionWallet.Alias, ndab[0]),
							Code:       404,
						}
						return err
					}

				}

			}

		}
	}
	return nil
}

func logDiscordFailedTokenizedAssetSubscription(msg string) {
	discord.WebhookURL = "https://discord.com/api/webhooks/827986576415129663/wqMKp9wxB_fxs9Q3zlMKCNPGENXmD_ueUnL8hVCu1wmRfD2wkXAjfP85k1Ro_2_wGfiY"
	if len(os.Getenv("FAILED_PAYMENT_ERROR_WEBHOOK")) > 50 {
		discord.WebhookURL = os.Getenv("FAILED_PAYMENT_ERROR_WEBHOOK")
	}
	discord.Say(msg)
}

func GetSwapEstimate(sourceAssetCode, sourceAssetIssuer, amount, destinationAssetCode, destinationAssetIssuer string, gc *sharedconfig.GlobalConfig) (swappedEstimate string) {
	destAsset := ""

	if len(destinationAssetIssuer) > 10 {
		destAsset = fmt.Sprintf("%s:%s", destinationAssetCode, destinationAssetIssuer)
	}
	pathInput := swapModels.SwapSendPathInput{
		DestinationAssets: destAsset,
		SourceAssetCode:   sourceAssetCode,
		SourceAssetIssuer: sourceAssetIssuer,
		SourceAmount:      amount,
	}
	_, swappedEstimate, err := GetStrictSendPaths(pathInput, gc.BantuExpansionClient)
	if err != nil {
		log.Println("[GetSwapEstimate]error fetching valid swap Path ", err)
		return "0"
	}
	return swappedEstimate
}
