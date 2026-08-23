package users

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"math/big"
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

func GetCountries(db *gorm.DB) (countries []userModels.Country) {
	countries = make([]userModels.Country, 0)
	db.Preload(clause.Associations).Find(&countries)

	return
}

func GetCountryConfigs(db *gorm.DB) (cc []userModels.CountryConfig) {
	cc = make([]userModels.CountryConfig, 0)
	db.Preload(clause.Associations).Find(&cc)

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

func GetTokenizationFormByID(formId uint64, gc *sharedconfig.GlobalConfig) (t userModels.JsonForm) {
	gc.DB.Where("id = ?", formId).First(&t)
	return
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

	if decimal.NewFromFloat(input.NumberOfTokenToBeIssued).GreaterThan(gc.TokenLimitAsDecimal()) {
		err = &tErrors.CustomError{Param: "numberOfTokenToBeIssued", Err: "error-token-limit-exceeded", ErrMessage: fmt.Sprintf("Issued Number of tokens cannot exceed %v", gc.TokenLimitAsString())}
		return
	}

	input.ProceedPayoutCurrency = strings.ToUpper(input.ProceedPayoutCurrency)
	if len(input.ProceedPayoutCurrency) == 0 {
		input.ProceedPayoutCurrency = "CNGN"
	}
	proceedPayoutCurrency := GetTokenizationCurrencyByCode(input.ProceedPayoutCurrency, gc.DB)
	if len(proceedPayoutCurrency.AssetCode) == 0 {
		err = &tErrors.CustomError{Param: "proceedPayoutCurrency", Err: "error-invalid-proceed-payout-currency", ErrMessage: "Proceed Payout currency code you supplied is invalid."}
		return
	}
	if userModels.IsInternalBalanceAsset(proceedPayoutCurrency.AssetCode, proceedPayoutCurrency.AssetIssuer, gc) {
		err = &tErrors.CustomError{Param: "proceedPayoutCurrency", Err: "error-invalid-proceed-payout-currency", ErrMessage: "Proceed Payout currency cannot be the internal balance token."}
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
	countryConfig := userModels.CountryCode(*ato.AssetCountryLocation).GetConfig(gc)
	if countryConfig.InternalBalanceTokenCode == nil {
		log.Printf("[SubmitTokenizationAssetInfoByInitiator]error internal balance token not configured for country: %v\n", *ato.AssetCountryLocation)
		err = &tErrors.CustomError{Param: "assetQuoteCurrency", Err: "error-invalid-asset-quote-currency", ErrMessage: "Internal balance token is not configured for this country."}
		return
	}
	ato.AssetQuoteCurrency = countryConfig.InternalBalanceTokenCode

	//check Trov balance
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
	input.ProceedPayoutCurrency = strings.ToUpper(input.ProceedPayoutCurrency)
	if len(input.ProceedPayoutCurrency) == 0 {
		input.ProceedPayoutCurrency = "CNGN"
	}
	proceedPayoutCurrency := GetTokenizationCurrencyByCode(input.ProceedPayoutCurrency, gc.DB)
	if len(proceedPayoutCurrency.AssetCode) == 0 {
		log.Printf("\n\n[SubmitTokenizationAssetInfo] error: payout currency is invalid: %v, received object: %+v\n\n", proceedPayoutCurrency.AssetCode, input)
		err = &tErrors.CustomError{Param: "proceedPayoutCurrency", Err: "error-invalid-proceed-payout-currency", ErrMessage: fmt.Sprintf("Proceed Payout currency code [%v] you supplied is invalid.", input.ProceedPayoutCurrency)}
		return
	}
	if userModels.IsInternalBalanceAsset(proceedPayoutCurrency.AssetCode, proceedPayoutCurrency.AssetIssuer, gc) {
		err = &tErrors.CustomError{Param: "proceedPayoutCurrency", Err: "error-invalid-proceed-payout-currency", ErrMessage: "Proceed Payout currency cannot be the internal balance token."}
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

	countryConfig := userModels.CountryCode(*ato.AssetCountryLocation).GetConfig(gc)
	if countryConfig.InternalBalanceTokenCode == nil {
		log.Printf("[SubmitTokenizationAssetInfo]error internal balance token not configured for country: %v\n", *ato.AssetCountryLocation)
		err = &tErrors.CustomError{Param: "assetQuoteCurrency", Err: "error-invalid-asset-quote-currency", ErrMessage: "Internal balance token is not configured for this country."}
		return
	}
	ato.AssetQuoteCurrency = countryConfig.InternalBalanceTokenCode

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
	e := gc.DB.Omit(clause.Associations).Save(&ato).Error
	if e != nil {
		log.Printf("[SubmitTokenizationAssetInfo] error saving tokenization to database  [%v] for %v: %v\n", input, initiator.Username, e)

		err = &tErrors.ErrorTemporaryServerError{}

	}
	if (ato.IssuingWalletPublicKey == nil || NotIssuedByIssuer) && len(input.AssetCode) > 0 && len(os.Getenv("TOKENIZATION_ISSUING_PROFILE")) > 1 && len(os.Getenv("TOKENIZATION_ISSUING_PROFILE_WALLET")) == 56 {

		//create issuing wallet
		ato, issuingWallet, err = AssignIssuingWallet(tokenizationID, gc)

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

	e = gc.DB.Omit(clause.Associations).Save(&ato).Error
	if e != nil {
		log.Printf("[SubmitTokenizationAssetInfo] error saving tokenization to database  [%v] for %v: %v\n", input, initiator.Username, e)

		err = &tErrors.ErrorTemporaryServerError{}
		return
	}
	ato, _, _ = GetTokenizedAssetByID(ato.ID, gc.DB)
	return ato, issuingWallet, nil
}

// AssignIssuingWallet assigns issuing wallet to tokenized asset instance
func AssignIssuingWallet(tokenizationID string, gc *sharedconfig.GlobalConfig) (ato userModels.TokenizedAsset, issuingWallet userModels.UserWallet, err error) {
	replacer := strings.NewReplacer("\r", "", "\n", "", " ", "")

	if len(strings.TrimSpace(os.Getenv("TOKENIZATION_ISSUING_PROFILE"))) == 0 {
		err = &tErrors.CustomError{Param: "issuingWalletPublicKey", Err: "error-default-issuing-profile-not-set", ErrMessage: "Issuing profile not set."}
		return
	}

	//referesh issuing wallet profile
	userModels.Username(os.Getenv("TOKENIZATION_ISSUING_PROFILE")).InvalidateUserCache(gc)

	ato, _, errorGetTokenizationByID := GetTokenizedAssetByID(tokenizationID, gc.DB)
	if errorGetTokenizationByID != nil {
		// error tokenization is already in progress
		log.Printf("[AssignIssuingWallet] Error fetching  tokenization with ID: %v\n", tokenizationID)
		err = errorGetTokenizationByID
		return

	}
	if ato.AssetCode == nil {
		err = &tErrors.CustomError{Param: "assetCode", Err: "error-asset-code-not-set", ErrMessage: "Asset code not set."}
		return
	}
	assetCode := (strings.ToUpper(replacer.Replace(strings.TrimSpace(*ato.AssetCode))))
	ato.AssetCode = &assetCode

	var NotIssuedByIssuer bool
	if ato.IssuingWalletAlias != nil && len(strings.TrimSpace(os.Getenv("TOKENIZATION_ISSUING_PROFILE"))) > 0 {
		NotIssuedByIssuer = !strings.HasPrefix(*ato.IssuingWalletAlias, strings.TrimSpace(os.Getenv("TOKENIZATION_ISSUING_PROFILE")))
	}

	if (ato.IssuingWalletPublicKey == nil || NotIssuedByIssuer) && len(*ato.AssetCode) > 0 && len(os.Getenv("TOKENIZATION_ISSUING_PROFILE")) > 1 && len(os.Getenv("TOKENIZATION_ISSUING_PROFILE_WALLET")) == 56 {

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

		walletTag := fmt.Sprintf("%v_issuer", *ato.AssetCode)
		p := userModels.SubWalletInfo{
			PublicKey:             issuer.Address(),
			WalletTag:             walletTag,
			WalletType:            1,
			LinkedWalletPublicKey: distributor.Address(),
		}
		_, err = CreateNewSubWallet(&tokenizationIssuer, &p, gc)

		if err != nil {
			log.Printf("[AssignIssuingWallet.CreateNewSubWallet: stage 1] Error creating issuing wallet [%v], err: %v\n", p.PublicKey, err)
			// err = &tErrors.CustomError{Param: "issuingPublicKey", Err: "error-invalid-issuer", ErrMessage: err}
			return
		}

		//sign transactions
		{

			if p.LinkedWalletMustSign == 1 {
				dsigned, e := middleware.SignBase64Txn(distributor.Seed(), p.Transaction, p.NetworkPassPhrase)
				if e != nil {
					log.Printf("[AssignIssuingWallet.SignBase64Txn] Error signing issuing wallet with linked wallet [%v], err: %v\n", distributor.Address(), e)
					err = &tErrors.CustomError{Param: "issuingPublicKey", Err: "error-invalid-issuer", ErrMessage: e.Error()}
					return
				}
				p.LinkedWalletSignature = dsigned
			}

			primarySignature, subwalletSignature, _, e := middleware.SignSubwalletBase64Txn(tokenizationIssuerProfileWalletKP.Seed(), issuer.Seed(), "", p.Transaction, p.NetworkPassPhrase)
			if e != nil {
				log.Printf("[AssignIssuingWallet.SignSubwalletBase64Txn] Error signing issuing wallet with primary, sub and linked wallets [%v] [%v] [%v], err: %v\n", tokenizationIssuerProfileWalletKP.Address(), issuer.Address(), "", e)
				err = &tErrors.CustomError{Param: "issuingPublicKey", Err: "error-invalid-issuer", ErrMessage: "Invalid primary, subwallet or linked wallet Signer"}
				return

			}

			p.PrimarySignature = primarySignature
			p.SubWalletSignature = subwalletSignature
			// p.LinkedWalletSignature = linkedWalletSignature

		}

		//second submission to blockchain
		_, err = CreateNewSubWallet(&tokenizationIssuer, &p, gc)
		if err != nil {
			log.Printf("[AssignIssuingWallet.CreateNewSubWallet: stage 2] Error creating issuing wallet [%v], err: %v\n", p.PublicKey, err)
			// err = &tErrors.CustomError{Param: "issuingPublicKey", Err: "error-invalid-issuer", ErrMessage: err}
			return
		}

		log.Printf("[AssignIssuingWallet.CreateNewSubWallet] Succesfully Created issuing wallet [%v], txID: %v\n", p.PublicKey, p.TransactionID)

		w, e := userModels.UserWalletID(p.PublicKey).GetWallet(gc.DB, gc)
		if e != nil {
			log.Printf("[AssignIssuingWallet] Error fetching issuing wallet [%v], err: %v\n", p.PublicKey, e)
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
			log.Printf("[AssignIssuingWallet] Error fetching issuing wallet [%v], err: %v\n", ato.IssuingWalletPublicKey, e)
			err = e
			return
		}

		//set issuing wallet
		issuingWallet = w
		ato.IssuingWalletAlias = &w.Alias
		ato.IssuingWalletPublicKey = &w.ID
		ato.MarketMakingWallet = w.LinkedWalletPublicKey
		ato.WalletToHoldAssetsNotForSale = w.LinkedWalletPublicKey
		ato.WalletToHoldAssetsNotForSale = w.LinkedWalletPublicKey
	}

	e := gc.DB.Omit(clause.Associations).Save(&ato).Error
	if e != nil {
		log.Printf("[AssignIssuingWallet] error saving tokenization to database  for %v: %v\n", tokenizationID, e)

		err = &tErrors.ErrorTemporaryServerError{}
		return
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
	if ato.TokenizationFeeID != nil {
		if *ato.TokenizationFeeID == 0 {
			// return it to fee state.
			ato.AssetTokenizationStatus = 0
			gc.DB.Omit(clause.Associations).Save(&ato)
			log.Printf("[VetTokenizationAssetInfo] Error tokenization fee was not selected: %v\n", tokenizationID)
			err = &tErrors.CustomError{Param: "issuingWalletPublicKey", Err: "error-tokenization-fee-not-selected", ErrMessage: "Tokenization fee was not selected. Status has been returned to allow applicant to choose fee."}
			return

		}
	} else {
		//
		ato.AssetTokenizationStatus = 0
		gc.DB.Omit(clause.Associations).Save(&ato)
		log.Printf("[VetTokenizationAssetInfo] Error tokenization fee was not selected: %v\n", tokenizationID)
		err = &tErrors.CustomError{Param: "issuingWalletPublicKey", Err: "error-tokenization-fee-not-selected", ErrMessage: "Tokenization fee was not selected. Status has been returned to allow applicant to choose fee."}
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

	capEndDate := ato.SalesStart.AddDate(0, 0, ato.CapDurationInDays)
	if capEndDate.After(ato.SalesEnd) {
		// cap exceeds
		log.Printf("[VetTokenizationAssetInfo] Error TID: %v cap end date %v exceeds sales end date: %v\n", ato.ID, capEndDate, ato.SalesEnd)
		err = &tErrors.CustomError{Param: "salesEnd", Err: "error-invalid-cap-duration", ErrMessage: fmt.Sprintf("Invalid Cap duration. Cap end date %v exceeds sales end date %v", capEndDate, ato.SalesEnd)}
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

	// Legal Adviser and Financial Adviser are optional A5 stakeholder assignments.
	ato.LegalAdviserID = input.LegalAdviserID
	ato.FinancialAdviserID = input.FinancialAdviserID

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
		return
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

	var NotIssuedByIssuer bool
	if ato.IssuingWalletAlias != nil && len(strings.TrimSpace(os.Getenv("TOKENIZATION_ISSUING_PROFILE"))) > 0 {
		NotIssuedByIssuer = !strings.HasPrefix(*ato.IssuingWalletAlias, strings.TrimSpace(os.Getenv("TOKENIZATION_ISSUING_PROFILE")))
	}
	e := gc.DB.Omit(clause.Associations).Save(&ato).Error
	if e != nil {
		log.Printf("[AcknowledgeTokenizationFeePayment] error saving tokenization to database  for %v (%v): %v\n", tokenizationID, initiator.Username, e)

		err = &tErrors.ErrorTemporaryServerError{}

	}
	assetCodeExists := ato.AssetCode != nil
	if (ato.IssuingWalletPublicKey == nil || NotIssuedByIssuer) && assetCodeExists && len(os.Getenv("TOKENIZATION_ISSUING_PROFILE")) > 1 && len(os.Getenv("TOKENIZATION_ISSUING_PROFILE_WALLET")) == 56 {

		//create issuing wallet
		ato, _, err = AssignIssuingWallet(tokenizationID, gc)

	}
	e = gc.DB.Omit(clause.Associations).Save(&ato).Error
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
	if ato.TokenizationFeeID != nil {
		if *ato.TokenizationFeeID == 0 {
			//
			ato.AssetTokenizationStatus = 0
			gc.DB.Omit(clause.Associations).Save(&ato)
			log.Printf("[ConfirmTokenizationApplicationInfoByInitiator] Error tokenization fee was not selected: %v\n", tokenizationID)
			err = &tErrors.CustomError{Param: "issuingWalletPublicKey", Err: "error-tokenization-fee-not-selected", ErrMessage: "Tokenization fee was not selected. Please go and choose a fee structure of your choice before you can continue."}
			return

		}
	} else {
		//
		ato.AssetTokenizationStatus = 0
		gc.DB.Omit(clause.Associations).Save(&ato)
		log.Printf("[ConfirmTokenizationApplicationInfoByInitiator] Error tokenization fee was not selected: %v\n", tokenizationID)
		err = &tErrors.CustomError{Param: "issuingWalletPublicKey", Err: "error-tokenization-fee-not-selected", ErrMessage: "Tokenization fee was not selected. Please go and choose a fee structure of your choice before you can continue."}
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
	ato.DateSubmitted = time.Now()
	//save again because fee value has been added to object
	e = dbTX.Omit(clause.Associations).Save(&ato).Error
	if e != nil {
		log.Printf("[ConfirmTokenizationApplicationInfoByInitiator] error saving tokenization to database  [%+v] for %v: %v\n", ato, initiator.Username, e)

		err = &tErrors.ErrorTemporaryServerError{}

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

func GetTMTokenizationStat(gc *sharedconfig.GlobalConfig) (data userModels.TokenizationStatForTM) {
	var s userModels.SummaryCat
	data.Message = "tokenized asset statistics fetched successfully"
	totalSql := `WITH total as (SELECT count(0) as count, 
sum(asset_current_value) as total_current_value,
sum(value_of_tokenized_asset) as total_tokenized_value,
sum(number_of_token_to_be_issued) as total_tokens_to_be_issued,
sum(number_of_token_to_be_sold) as total_tokens_to_be_sold,
sum(price_per_token) as total_price_per_token,
avg(price_per_token)::numeric(13,7) as average_price_per_token
FROM tokenized_assets where asset_tokenization_status > ?)

select total.count, total.total_current_value, 
total.total_tokenized_value, 
total.total_tokens_to_be_issued,
total.total_tokens_to_be_sold,
total.total_price_per_token,
total.average_price_per_token
from total`

	unsubmittedSql := `WITH unsubmitted as (SELECT count(0) as count, 
sum(asset_current_value) as total_current_value,
sum(value_of_tokenized_asset) as total_tokenized_value,
sum(number_of_token_to_be_issued) as total_tokens_to_be_issued,
sum(number_of_token_to_be_sold) as total_tokens_to_be_sold,
sum(price_per_token) as total_price_per_token,
avg(price_per_token)::numeric(13,7) as average_price_per_token
FROM tokenized_assets where asset_tokenization_status < ? AND Due_Diligence_Fail=0)

select unsubmitted.count, 
unsubmitted.total_current_value, 
unsubmitted.total_tokenized_value, 
unsubmitted.total_tokens_to_be_issued,
unsubmitted.total_tokens_to_be_sold,
unsubmitted.total_price_per_token,
unsubmitted.average_price_per_token
from unsubmitted`

	submittedSql := `WITH submitted as (SELECT count(0) as count, 
sum(asset_current_value) as total_current_value,
sum(value_of_tokenized_asset) as total_tokenized_value,
sum(number_of_token_to_be_issued) as total_tokens_to_be_issued,
sum(number_of_token_to_be_sold) as total_tokens_to_be_sold,
sum(price_per_token) as total_price_per_token,
avg(price_per_token)::numeric(13,7) as average_price_per_token
FROM tokenized_assets where asset_tokenization_status >0 AND asset_tokenization_status < ?)

select submitted.count, 
submitted.total_current_value, 
submitted.total_tokenized_value, 
submitted.total_tokens_to_be_issued,
submitted.total_tokens_to_be_sold,
submitted.total_price_per_token,
submitted.average_price_per_token
from submitted`

	pendingSql := `WITH pending as (SELECT count(0) as count, 
sum(asset_current_value) as total_current_value,
sum(value_of_tokenized_asset) as total_tokenized_value,
sum(number_of_token_to_be_issued) as total_tokens_to_be_issued,
sum(number_of_token_to_be_sold) as total_tokens_to_be_sold,
sum(price_per_token) as total_price_per_token,
avg(price_per_token)::numeric(13,7) as average_price_per_token
FROM tokenized_assets where asset_tokenization_status = ?)

select pending.count, 
pending.total_current_value, 
pending.total_tokenized_value, 
pending.total_tokens_to_be_issued,
pending.total_tokens_to_be_sold,
pending.total_price_per_token,
pending.average_price_per_token
from pending`

	rejectedSql := `WITH rejected as (SELECT count(0) as count, 
sum(asset_current_value) as total_current_value,
sum(value_of_tokenized_asset) as total_tokenized_value,
sum(number_of_token_to_be_issued) as total_tokens_to_be_issued,
sum(number_of_token_to_be_sold) as total_tokens_to_be_sold,
sum(price_per_token) as total_price_per_token,
avg(price_per_token)::numeric(13,7) as average_price_per_token
FROM tokenized_assets where asset_tokenization_status = 0 AND Due_Diligence_Fail=?)

select rejected.count, 
rejected.total_current_value, 
rejected.total_tokenized_value, 
rejected.total_tokens_to_be_issued,
rejected.total_tokens_to_be_sold,
rejected.total_price_per_token,
rejected.average_price_per_token
from rejected`

	approvedSql := `WITH approved as (SELECT count(0) as count, 
sum(asset_current_value) as total_current_value,
sum(value_of_tokenized_asset) as total_tokenized_value,
sum(number_of_token_to_be_issued) as total_tokens_to_be_issued,
sum(number_of_token_to_be_sold) as total_tokens_to_be_sold,
sum(price_per_token) as total_price_per_token,
avg(price_per_token)::numeric(13,7) as average_price_per_token
FROM tokenized_assets where asset_tokenization_status > ? AND asset_tokenization_status < 8)

select approved.count, 
approved.total_current_value, 
approved.total_tokenized_value, 
approved.total_tokens_to_be_issued,
approved.total_tokens_to_be_sold,
approved.total_price_per_token,
approved.average_price_per_token
from approved`

	liquidatedSql := `WITH liquidated as (SELECT count(0) as count, 
sum(asset_current_value) as total_current_value,
sum(value_of_tokenized_asset) as total_tokenized_value,
sum(number_of_token_to_be_issued) as total_tokens_to_be_issued,
sum(number_of_token_to_be_sold) as total_tokens_to_be_sold,
sum(price_per_token) as total_price_per_token,
avg(price_per_token)::numeric(13,7) as average_price_per_token
FROM tokenized_assets where asset_tokenization_status = ?)

select liquidated.count, 
liquidated.total_current_value, 
liquidated.total_tokenized_value, 
liquidated.total_tokens_to_be_issued,
liquidated.total_tokens_to_be_sold,
liquidated.total_price_per_token,
liquidated.average_price_per_token
from liquidated`

	refundedSql := `WITH refunded as (SELECT count(0) as count, 
sum(asset_current_value) as total_current_value,
sum(value_of_tokenized_asset) as total_tokenized_value,
sum(number_of_token_to_be_issued) as total_tokens_to_be_issued,
sum(number_of_token_to_be_sold) as total_tokens_to_be_sold,
sum(price_per_token) as total_price_per_token,
avg(price_per_token)::numeric(13,7) as average_price_per_token
FROM tokenized_assets where asset_tokenization_status = ?)

select refunded.count, 
refunded.total_current_value, 
refunded.total_tokenized_value, 
refunded.total_tokens_to_be_issued,
refunded.total_tokens_to_be_sold,
refunded.total_price_per_token,
refunded.average_price_per_token
from refunded`

	e := gc.DB.Raw(totalSql, 2).Scan(&s).Error
	if e != nil {
		log.Println("[GetTMTokenizationStat]error:", e)
		data.Message = "error retrieving total category info"
	}
	data.Data.Total = s

	e = gc.DB.Raw(approvedSql, 3).Scan(&s).Error
	if e != nil {
		log.Println("[GetTMTokenizationStat]error:", e)
		data.Message = "error retrieving approved category info"
	}
	data.Data.Approved = s

	e = gc.DB.Raw(submittedSql, 3).Scan(&s).Error
	if e != nil {
		log.Println("[GetTMTokenizationStat]error:", e)
		data.Message = "error retrieving submitted category info"
	}
	data.Data.Submitted = s

	e = gc.DB.Raw(unsubmittedSql, 1).Scan(&s).Error
	if e != nil {
		log.Println("[GetTMTokenizationStat]error:", e)
		data.Message = "error retrieving unsubmitted category info"
	}
	data.Data.Unsubmitted = s

	e = gc.DB.Raw(pendingSql, 3).Scan(&s).Error
	if e != nil {
		log.Println("[GetTMTokenizationStat]error:", e)
		data.Message = "error retrieving pending category info"
	}
	data.Data.Pending = s

	e = gc.DB.Raw(rejectedSql, 1).Scan(&s).Error
	if e != nil {
		log.Println("[GetTMTokenizationStat]error:", e)
		data.Message = "error retrieving rejected category info"

	}
	data.Data.Rejected = s

	e = gc.DB.Raw(refundedSql, 8).Scan(&s).Error
	if e != nil {
		log.Println("[GetTMTokenizationStat]error:", e)
		data.Message = "error retrieving refunded category info"

	}
	data.Data.Refunded = s

	e = gc.DB.Raw(liquidatedSql, 7).Scan(&s).Error
	if e != nil {
		log.Println("[GetTMTokenizationStat]error:", e)
		data.Message = "error retrieving liquidated category info"

	}
	data.Data.Liquidated = s
	data.Status = "OK"
	data.Timestamp = time.Now().String()
	//return the result
	return
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

		log.Printf("[SubscribeToTokenizedAsset] Error Unable to verify wallet owner of subscribing wallet %v\n", subscriberWallet.Alias)
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
		maxPurchase := decimal.NewFromFloat(ta.CapAmountInFiat - ta.SumAmountBoughtByWalletOwner(subscriberWallet.Alias, gc))
		if maxPurchase.IsZero() {
			err = &tErrors.CustomError{Param: "amount", Err: "error-cap-amount-exceeded", ErrMessage: fmt.Sprintf("You have exhausted the allowed purchase cap of %v%v worth of %v at this time.", *ta.AssetQuoteCurrency, ta.CapAmountInFiat, *ta.AssetCode)}

		} else {

			err = &tErrors.CustomError{Param: "amount", Err: "error-cap-amount-exceeded", ErrMessage: fmt.Sprintf("You can only purchase upto %v%v worth of %v at this time.", *ta.AssetQuoteCurrency, maxPurchase.String(), *ta.AssetCode)}
		}
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

	taSubscription.UpdateTokenizedAssetSubscriptionFromInput(subscriber.Username, subscriberWallet, input, ta, gc)

	//begin transaction xdr

	// swapInfo.SourceAmount = decimal.NewFromFloat(input.Amount).Truncate(7).String()
	// swapAmount := decimal.NewFromFloat(input.Amount).Truncate(7)

	client := gc.BantuExpansionClient
	//transform codes and issuer
	swapInfo.DestinationAssetCode = strings.ToUpper(*ta.AssetCode)
	swapInfo.DestinationAssetIssuer = strings.ToUpper(*ta.IssuingWalletPublicKey)

	// resolve payment asset: any approved stablecoin, defaulting to CNGN for unchanged clients.
	// a client-supplied issuer is never trusted verbatim - the issuer is always re-resolved server-side.
	paymentAssetCode := strings.ToUpper(strings.TrimSpace(input.PaymentAssetCode))
	if len(paymentAssetCode) == 0 {
		paymentAssetCode = "CNGN"
	}
	paymentCurrency := GetTokenizationCurrencyByCode(paymentAssetCode, dbTX)
	if len(paymentCurrency.AssetCode) == 0 {
		err = &tErrors.CustomError{Param: "paymentAssetCode", Err: "error-invalid-payment-asset", ErrMessage: fmt.Sprintf("Payment asset code [%v] you supplied is invalid.", paymentAssetCode)}
		return
	}
	if len(input.PaymentAssetIssuer) > 0 && !strings.EqualFold(input.PaymentAssetIssuer, paymentCurrency.AssetIssuer) {
		err = &tErrors.CustomError{Param: "paymentAssetIssuer", Err: "error-invalid-payment-asset", ErrMessage: "Payment asset issuer does not match the registered issuer for this currency."}
		return
	}
	if userModels.IsInternalBalanceAsset(paymentCurrency.AssetCode, paymentCurrency.AssetIssuer, gc) || strings.EqualFold(paymentCurrency.AssetCode, strings.ToUpper(*ta.AssetCode)) {
		err = &tErrors.CustomError{Param: "paymentAssetCode", Err: "error-invalid-payment-asset", ErrMessage: "Payment asset cannot be the internal balance token or the tokenized asset itself."}
		return
	}
	swapInfo.SourceAssetCode = strings.ToUpper(paymentCurrency.AssetCode)
	swapInfo.SourceAssetIssuer = strings.ToUpper(paymentCurrency.AssetIssuer)
	// write the resolved (normalized/defaulted) payment asset back onto both the in-memory
	// subscription row (for the direct, non-multiparty save below) and input itself - input is what
	// gets serialized into PendingAuth.TransactionInfoStr for the multiparty path, and later
	// re-hydrated in approvals.go's ApproveTransaction to build the TokenizedAssetSubscription row
	// there via this same UpdateTokenizedAssetSubscriptionFromInput helper, so it must already carry
	// the resolved values rather than whatever (possibly empty) code the client originally sent
	taSubscription.PaymentAssetCode = paymentCurrency.AssetCode
	taSubscription.PaymentAssetIssuer = paymentCurrency.AssetIssuer
	input.PaymentAssetCode = paymentCurrency.AssetCode
	input.PaymentAssetIssuer = paymentCurrency.AssetIssuer

	//always save subscriptions afresh - done here, after the payment asset is resolved, so the
	//saved row already carries the correct PaymentAssetCode/PaymentAssetIssuer
	if input.Multiparty == 0 {
		e := dbTX.Omit(clause.Associations).Save(&taSubscription).Error
		if e != nil {
			log.Printf("[SubscribeToTokenizedAsset] error saving tokenized asset subscription  to database  [%+v] for %v: %v\n", taSubscription, subscriber.Username, e)

			err = &tErrors.ErrorTemporaryServerError{}
			return
		}
	}

	if e := swapServices.ValidateSwapSendInfo(&swapInfo); e != nil {
		log.Printf("[SubscribeToTokenizedAsset] error validating tokenized asset subscription [%+v] for %v: %v\n", taSubscription, subscriber.Username, e)

		err = e

		return
	}
	// get market offer of the asset
	{
		offers, e := ta.GetMarketOffers(gc)
		if e != nil {
			log.Printf("[SubscribeToTokenizedAsset] error fetching matket offer of tokenized asset [%+v] for %v: %v\n", taSubscription, subscriber.Username, e)

			err = &tErrors.ErrorTemporaryServerError{}
			return
		}
		if len(offers.Embedded.Records) > 0 {
			//get the first
			offer := offers.Embedded.Records[0]
			//get price
			remainingBuyingLiability := decimal.RequireFromString(offer.Amount).Mul(decimal.RequireFromString(offer.Price)).Truncate(7)
			//if the subscription amount is greater than the buying liability, then reject action.
			decPurchaseAmountFiat := decimal.NewFromFloat(input.Amount)

			if decPurchaseAmountFiat.GreaterThan(remainingBuyingLiability) {
				//not enough liquidity
				err = &tErrors.CustomError{Param: "amount", Err: "error-low-liquidity", ErrMessage: fmt.Sprintf("Remaining %v token can be purchased with a maximum of %v %v. Please adjust your purchase amount.", swapInfo.DestinationAssetCode, remainingBuyingLiability.String(), swapInfo.SourceAssetCode)}
				return
			}

			if remainingBuyingLiability.IsZero() {
				//not enough liquidity
				err = &tErrors.CustomError{Param: "amount", Err: "error-no-liquidity", ErrMessage: fmt.Sprintf("The allocated %v asset has sold out.", swapInfo.DestinationAssetCode)}
				return
			}

		}
	}
	if len(input.TransactionSignature) == 0 {
		xdrBase64, e := generateAssetSubscriptionXdr(subscriberWallet, ta, &swapInfo, gc)
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
		subscriber.InvalidateUserWalletCache(gc)
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
		subscriberWallet.InvalidateUserCache(gc)
		subscriber.InvalidateUserWalletCache(gc)
		return taSubscription, nil
	}
	log.Println("[SubscribeToTokenizedAsset]UNKNOWN OPTION FOR ACTION")
	err = &tErrors.ErrorTemporaryServerError{}
	return

}

// SubscribeToTokenizedAssetByFiat is the fiat counterpart of SubscribeToTokenizedAsset. The buyer never
// touches Stellar directly here: the server builds and server-signs (channel account + internal balance
// issuing multisig + tokenization issuing profile wallet if needed) a transaction that mints the
// country's internal balance token to the buyer and swaps it for the tokenized asset, the buyer only
// adds their own signature, and submission to the blockchain is deferred until the fiat payment provider
// (Flutterwave) confirms payment via webhook - see postCallbacksFlutterwaveWebhookHandler.
//
// The client-supplied input.ID is required on both calls and is used as the primary key of both the
// FiatPaymentInvoice row and the TokenizedAssetSubscription row it creates, and is expected to double as
// the payment provider's transaction reference so the webhook can find its way back to this invoice.
//
// Call 1 (input.TransactionSignature empty): generates the transaction xdr using a reserved channel
// account and persists it on a new FiatPaymentInvoice + TokenizedAssetSubscription pair. A retried call 1
// with the same id returns the already-generated invoice as-is rather than regenerating the xdr and
// reserving a second channel account.
// Call 2 (input.TransactionSignature present): loads the existing invoice and stores the buyer's
// signature on it. Nothing is submitted to the blockchain in this call - that only happens later, from
// the webhook.
func SubscribeToTokenizedAssetByFiat(subscriber *userModels.User, subscriberWallet *userModels.UserWallet, ta *userModels.TokenizedAsset, input *userModels.FiatTokenizedAssetSubscriptionInput, gc *sharedconfig.GlobalConfig) (invoice userModels.FiatPaymentInvoice, err error) {
	if len(input.ID) == 0 {
		err = &tErrors.CustomError{Param: "id", Err: "error-missing-parameter", ErrMessage: "id is required."}
		return
	}
	input.TokenizedAssetID = ta.ID
	input.SubscriberUsername = subscriber.Username

	// call 2: buyer's signature is being supplied to finalize a previously-generated invoice
	if len(input.TransactionSignature) > 0 {
		e := gc.DB.Where("id = ? AND username = ?", input.ID, subscriber.Username).First(&invoice).Error
		if e != nil {
			err = &tErrors.CustomError{Param: "id", Err: "error-invalid-request", ErrMessage: "No pending fiat asset purchase invoice found for this id."}
			return
		}
		if invoice.TransactionSignature != nil || invoice.Status != "PENDING" {
			err = &tErrors.CustomError{Param: "id", Err: "error-invalid-request", ErrMessage: "This fiat asset purchase invoice has already been signed or processed."}
			return
		}
		signature := input.TransactionSignature
		e = gc.DB.Model(&userModels.FiatPaymentInvoice{}).Where("id = ?", invoice.ID).Update("transaction_signature", signature).Error
		if e != nil {
			log.Printf("[SubscribeToTokenizedAssetByFiat] error saving transaction signature for invoice %v: %v\n", invoice.ID, e)
			err = &tErrors.ErrorTemporaryServerError{}
			return
		}
		invoice.TransactionSignature = &signature
		return invoice, nil
	}

	// call 1: an already-generated, still-unsigned invoice for this id is reused as-is rather than
	// regenerating the xdr and reserving a second channel account
	{
		var existing userModels.FiatPaymentInvoice
		e := gc.DB.Where("id = ?", input.ID).First(&existing).Error
		if e == nil {
			if existing.TransactionSignature != nil {
				err = &tErrors.CustomError{Param: "id", Err: "error-invalid-request", ErrMessage: "This fiat asset purchase invoice has already been signed."}
				return
			}
			return existing, nil
		}
	}

	walletOwner, e := subscriberWallet.GetWalletOwner(gc.DB, gc)
	if e != nil {
		log.Printf("[SubscribeToTokenizedAssetByFiat] Error Unable to verify wallet owner of subscribing wallet %v\n", subscriberWallet.Alias)
		err = &tErrors.CustomError{Param: "issuingWalletPublicKey", Err: "error-invalid-kyc", ErrMessage: "Unable to verify wallet owner."}
		return
	}
	if walletOwner.KYCVerified == 0 {
		log.Printf("[SubscribeToTokenizedAssetByFiat] Error Wallet owner %v has not met KYC status for asset %v\n", walletOwner.Username, *ta.AssetCode)
		err = &tErrors.CustomError{Param: "issuingWalletPublicKey", Err: "error-invalid-kyc", ErrMessage: fmt.Sprintf("%v has not passed KYC to purchase this tokenized asset %v.", walletOwner.Username, *ta.AssetCode)}
		return
	}
	if ta.AssetTokenizationStatus != 5 && ta.AssetTokenizationStatus != 6 {
		log.Printf("[SubscribeToTokenizedAssetByFiat] Error Tokenized asset not in sales yet: %v\n", ta.ID)
		err = &tErrors.CustomError{Param: "issuingWalletPublicKey", Err: "error-invalid-request", ErrMessage: "Only projects that are in sales can accept purchase."}
		return
	}
	capEndDate := ta.SalesStart.AddDate(0, 0, ta.CapDurationInDays)
	if capEndDate.After(time.Now()) && ta.CapOnPurchase > 0 && ta.CapAmountInFiat > 0 && decimal.NewFromFloat(input.Amount+ta.SumAmountBoughtByWalletOwner(subscriberWallet.Alias, gc)).Truncate(7).GreaterThan(decimal.NewFromFloat(ta.CapAmountInFiat)) {
		log.Printf("[SubscribeToTokenizedAssetByFiat] Error Tokenized asset Cap exceeded: %v, Amount In Cap: %v\n", ta.ID, decimal.NewFromFloat(ta.CapAmountInFiat).String())
		maxPurchase := decimal.NewFromFloat(ta.CapAmountInFiat - ta.SumAmountBoughtByWalletOwner(subscriberWallet.Alias, gc))
		if maxPurchase.IsZero() {
			err = &tErrors.CustomError{Param: "amount", Err: "error-cap-amount-exceeded", ErrMessage: fmt.Sprintf("You have exhausted the allowed purchase cap of %v%v worth of %v at this time.", *ta.AssetQuoteCurrency, ta.CapAmountInFiat, *ta.AssetCode)}
		} else {
			err = &tErrors.CustomError{Param: "amount", Err: "error-cap-amount-exceeded", ErrMessage: fmt.Sprintf("You can only purchase upto %v%v worth of %v at this time.", *ta.AssetQuoteCurrency, maxPurchase.String(), *ta.AssetCode)}
		}
		return
	}
	if subscriberWallet.SharedAccessEnabled == 1 && subscriberWallet.NumberOfApprovalsNeeded > 0 {
		err = &tErrors.CustomError{Param: "walletPublicKey", Err: "error-invalid-request", ErrMessage: "Fiat purchase is not supported on shared-access wallets that require multiple approvals."}
		return
	}

	// get market offer of the asset - liquidity check
	{
		offers, e := ta.GetMarketOffers(gc)
		if e != nil {
			log.Printf("[SubscribeToTokenizedAssetByFiat] error fetching market offer of tokenized asset [%v] for %v: %v\n", ta.ID, subscriber.Username, e)
			err = &tErrors.ErrorTemporaryServerError{}
			return
		}
		if len(offers.Embedded.Records) > 0 {
			offer := offers.Embedded.Records[0]
			remainingBuyingLiability := decimal.RequireFromString(offer.Amount).Mul(decimal.RequireFromString(offer.Price)).Truncate(7)
			decPurchaseAmountFiat := decimal.NewFromFloat(input.Amount)
			if decPurchaseAmountFiat.GreaterThan(remainingBuyingLiability) {
				err = &tErrors.CustomError{Param: "amount", Err: "error-low-liquidity", ErrMessage: fmt.Sprintf("Remaining %v token can be purchased with a maximum of %v. Please adjust your purchase amount.", strings.ToUpper(*ta.AssetCode), remainingBuyingLiability.String())}
				return
			}
			if remainingBuyingLiability.IsZero() {
				err = &tErrors.CustomError{Param: "amount", Err: "error-no-liquidity", ErrMessage: fmt.Sprintf("The allocated %v asset has sold out.", strings.ToUpper(*ta.AssetCode))}
				return
			}
		}
	}

	var swapInfo swapModels.SwapSendInfo
	swapInfo.Messages = make([]string, 0)
	swapInfo.SourceAmount = decimal.NewFromFloat(input.Amount).Truncate(7).String()
	swapInfo.SwapAmount = swapInfo.SourceAmount
	swapInfo.DestinationAssetCode = strings.ToUpper(*ta.AssetCode)
	swapInfo.DestinationAssetIssuer = strings.ToUpper(*ta.IssuingWalletPublicKey)

	xdrBase64, channelAccountPublicKey, e := generateAssetSubscriptionFiatXdr(subscriberWallet, ta, &swapInfo, gc)
	if e != nil {
		log.Printf("[SubscribeToTokenizedAssetByFiat] error generating fiat asset subscription xdr for %v: %v\n", subscriber.Username, e)
		err = e
		return
	}

	countryConfig := userModels.CountryCode(*ta.AssetCountryLocation).GetConfig(gc)

	dbTX := gc.DB.Begin()
	defer dbTX.Rollback()

	var taSubscription userModels.TokenizedAssetSubscription
	taSubscription.ID = input.ID
	taSubscription.TokenizedAssetID = ta.ID
	taSubscription.AssetCode = *ta.AssetCode
	taSubscription.AssetIssuer = *ta.IssuingWalletPublicKey
	taSubscription.WalletAlias = subscriberWallet.Alias
	taSubscription.WalletPublicKey = subscriberWallet.ID
	taSubscription.Amount = decimal.NewFromFloat(input.Amount).Truncate(7).InexactFloat64()
	taSubscription.Price = ta.PricePerToken
	taSubscription.SubscriberUsername = subscriber.Username
	taSubscription.PaymentAssetCode = *countryConfig.InternalBalanceTokenCode
	taSubscription.PaymentAssetIssuer = "FIAT"

	if e := dbTX.Omit(clause.Associations).Create(&taSubscription).Error; e != nil {
		log.Printf("[SubscribeToTokenizedAssetByFiat] error saving tokenized asset subscription [%+v] for %v: %v\n", taSubscription, subscriber.Username, e)
		gc.ReleaseInUseChannelAccount(channelAccountPublicKey)
		err = &tErrors.ErrorTemporaryServerError{}
		return
	}

	walletAlias := subscriberWallet.Alias
	walletPublicKey := subscriberWallet.ID
	tokenizedAssetID := ta.ID

	invoice = userModels.FiatPaymentInvoice{
		ID:                input.ID,
		ServiceProvider:   "flutterwave",
		Username:          subscriber.Username,
		Amount:            input.Amount,
		PaymentType:       "ASSET PURCHASE",
		Status:            "PENDING",
		TokenizedAssetID:  &tokenizedAssetID,
		WalletAlias:       &walletAlias,
		WalletPublicKey:   &walletPublicKey,
		Transaction:       &xdrBase64,
		TransactionSource: &channelAccountPublicKey,
	}
	if e := dbTX.Create(&invoice).Error; e != nil {
		log.Printf("[SubscribeToTokenizedAssetByFiat] error saving fiat payment invoice [%+v] for %v: %v\n", invoice, subscriber.Username, e)
		gc.ReleaseInUseChannelAccount(channelAccountPublicKey)
		err = &tErrors.ErrorTemporaryServerError{}
		return
	}

	dbTX.Commit()
	return invoice, nil
}

func generateAssetSubscriptionXdr(wallet *userModels.UserWallet, ta *userModels.TokenizedAsset, swapInfo *swapModels.SwapSendInfo, gc *sharedconfig.GlobalConfig) (string, error) {
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

	if ta.FundsHoldingWalletPublicKey == nil {
		return "", &tErrors.CustomError{Param: "fundsHoldingWalletPublicKey", Err: "error-funds-holding-wallet-not-set", ErrMessage: "This tokenized asset is not yet configured to accept purchases."}
	}
	if ta.AssetCountryLocation == nil {
		return "", &tErrors.CustomError{Param: "publicKey", Err: "error-invalid-quote-currency", ErrMessage: "Tokenization does not have a valid country of location."}
	}
	// resolve the internal balance token (issued 1:1 on-demand to the buyer) from the asset's country config
	countryConfig := userModels.CountryCode(*ta.AssetCountryLocation).GetConfig(gc)
	if countryConfig.InternalBalanceTokenCode == nil || countryConfig.InternalTokenIssuer == nil {
		return "", &tErrors.CustomError{Param: "publicKey", Err: "error-invalid-quote-currency", ErrMessage: "Tokenization does not have valid tokenization currency."}
	}
	internalBalanceAsset := txnbuild.CreditAsset{Code: *countryConfig.InternalBalanceTokenCode, Issuer: *countryConfig.InternalTokenIssuer}

	internalBalanceIssuingSigners, err := getInternalBalanceIssuingSigners()
	if err != nil {
		log.Println("[generateAssetSubscriptionXdr] error parsing internal balance issuing signers", err)
		return "", &tErrors.ErrorTemporaryServerError{}
	}

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
	_, sourceAccountTrustsInternalBalanceAsset, _, _, _, _ := network.BlockchainAccountProperties(client, wallet.ID, internalBalanceAsset)

	if sourceAccountErr != nil {
		return "", sourceAccountErr
	}

	if !sourceAccountExists {
		return "", &tErrors.ErrorUnderfundedAccount{}
	}

	log.Printf("[generateAssetSubscriptionXdr]obtained source account balance:\n%v balance is %v\n%v balance is %v\n", nativeAssetCode, sourceAccountNativeBalance, sourceAsset.GetCode(), sourceAccountCustomBalance)

	if sourceAccountCustomBalance.LessThan(amountToSwap) {
		return "", &tErrors.ErrorUnderfundedAccount{Detail: fmt.Sprintf("Not enough funds. Needs Extra %v %v or you reduce same from the amount you want to swap.", (amountToSwap).Sub(sourceAccountCustomBalance), sourceAsset.GetCode())}
	}

	// 1. send the stablecoin to the tokenized asset's funds holding wallet
	ops = append(ops, &txnbuild.Payment{
		Destination:   *ta.FundsHoldingWalletPublicKey,
		Amount:        swapInfo.SourceAmount,
		Asset:         sourceAsset,
		SourceAccount: wallet.ID,
	})

	if !sourceAccountTrustsInternalBalanceAsset {
		// 2. establish the buyer's trustline to the internal balance token
		ops = append(ops, &txnbuild.ChangeTrust{
			Line:          txnbuild.ChangeTrustAssetWrapper{Asset: internalBalanceAsset},
			Limit:         gc.TokenLimitAsString(),
			SourceAccount: wallet.ID,
		})
		// 3. authorize that trustline
		ops = append(ops, &txnbuild.SetTrustLineFlags{
			Trustor:       wallet.ID,
			Asset:         internalBalanceAsset,
			SetFlags:      []txnbuild.TrustLineFlag{txnbuild.TrustLineAuthorized},
			SourceAccount: *countryConfig.InternalTokenIssuer,
		})
	}

	// 4. issue the internal balance token 1:1 to the buyer
	ops = append(ops, &txnbuild.Payment{
		Destination:   wallet.ID,
		Amount:        swapInfo.SwapAmount,
		Asset:         internalBalanceAsset,
		SourceAccount: *countryConfig.InternalTokenIssuer,
	})

	if !destinationAsset.IsNative() {

		if !sourceAccountTrustsDestinationAsset {

			// 5. establish trustline to the tokenized asset
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

	//get sendPath from the internal balance token to the tokenized asset
	//using DestinationAccount will get paths to all assets in the destination account.
	//using destinationAssets gets path to only the asset
	destAsset := ""
	if !destinationAsset.IsNative() {
		destAsset = fmt.Sprintf("%s:%s", swapInfo.DestinationAssetCode, swapInfo.DestinationAssetIssuer)
	}
	pathInput := swapModels.SwapSendPathInput{
		DestinationAssets: destAsset,
		SourceAssetCode:   internalBalanceAsset.Code,
		SourceAssetIssuer: internalBalanceAsset.Issuer,
		SourceAmount:      swapInfo.SwapAmount,
	}
	path, swappedEstimate, err := GetStrictSendPaths(pathInput, client)
	if err != nil {
		log.Println("[generateAssetSubscriptionXdr]error fetching valid swap Path ", err)
		return "", err
	}

	// 6. swap the internal balance token for the tokenized asset
	ops = append(ops, &txnbuild.PathPaymentStrictSend{
		SendAsset:     internalBalanceAsset,
		SendAmount:    swapInfo.SwapAmount,
		Destination:   wallet.ID,
		DestAsset:     destinationAsset,
		DestMin:       swapDestMin.String(),
		Path:          path,
		SourceAccount: wallet.ID,
	})

	// 7. remove the internal balance trustline so it never persists outside of this transaction
	ops = append(ops, &txnbuild.ChangeTrust{
		Line:          txnbuild.ChangeTrustAssetWrapper{Asset: internalBalanceAsset},
		Limit:         "0",
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

	// sign for the internal balance issuer (multisig) to authorize the buyer's trustline and issue the internal balance token
	tx, err = tx.Sign(network.GetBlockchainNetworkPassPhrase(), internalBalanceIssuingSigners...)
	if err != nil {
		log.Println("[generateAssetSubscriptionXdr] error signing transaction with internal balance issuing signers ", err)
		return "", &tErrors.ErrorTemporaryServerError{}
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

// generateAssetSubscriptionFiatXdr builds the same purchase transaction as generateAssetSubscriptionXdr
// but drops the buyer's on-chain stablecoin payment leg (fiat settlement, confirmed later by the
// Flutterwave webhook, is what authorizes the mint) and always sources the transaction from a channel
// account rather than the buyer's own sequence account, since submission is deferred until the webhook
// fires - potentially long after this function returns. The channel account is only reserved
// (StoreInUseChannelAccount) once every failable step has already succeeded; on any earlier error it is
// returned to the pool untouched.
func generateAssetSubscriptionFiatXdr(wallet *userModels.UserWallet, ta *userModels.TokenizedAsset, swapInfo *swapModels.SwapSendInfo, gc *sharedconfig.GlobalConfig) (xdrBase64 string, channelAccountPublicKey string, err error) {
	swapDestMin := network.GetBlockchainSwapDestinationMin()
	client := gc.BantuExpansionClient
	messages := make([]string, 0)
	nativeAssetCode := os.Getenv("NATIVE_ASSET_CODE")

	if _, e := decimal.NewFromString(swapInfo.SourceAmount); e != nil {
		return "", "", &swapErrors.ErrorInvalidSwapAmount{}
	}
	if ta.AssetCountryLocation == nil {
		return "", "", &tErrors.CustomError{Param: "publicKey", Err: "error-invalid-quote-currency", ErrMessage: "Tokenization does not have a valid country of location."}
	}
	// resolve the internal balance token (issued 1:1 on-demand to the buyer) from the asset's country config
	countryConfig := userModels.CountryCode(*ta.AssetCountryLocation).GetConfig(gc)
	if countryConfig.InternalBalanceTokenCode == nil || countryConfig.InternalTokenIssuer == nil {
		return "", "", &tErrors.CustomError{Param: "publicKey", Err: "error-invalid-quote-currency", ErrMessage: "Tokenization does not have valid tokenization currency."}
	}
	internalBalanceAsset := txnbuild.CreditAsset{Code: *countryConfig.InternalBalanceTokenCode, Issuer: *countryConfig.InternalTokenIssuer}

	internalBalanceIssuingSigners, e := getInternalBalanceIssuingSigners()
	if e != nil {
		log.Println("[generateAssetSubscriptionFiatXdr] error parsing internal balance issuing signers", e)
		return "", "", &tErrors.ErrorTemporaryServerError{}
	}

	var destinationAsset txnbuild.Asset = txnbuild.NativeAsset{}
	if len(swapInfo.DestinationAssetCode) != 0 && !strings.EqualFold(swapInfo.DestinationAssetCode, nativeAssetCode) {
		destinationAsset = txnbuild.CreditAsset{Code: swapInfo.DestinationAssetCode, Issuer: swapInfo.DestinationAssetIssuer}
	}

	swapInfo.Messages = messages

	sourceAccountExists, _, _, _, _, sourceAccountErr := network.BlockchainAccountProperties(client, wallet.ID, txnbuild.NativeAsset{})
	if sourceAccountErr != nil {
		return "", "", sourceAccountErr
	}
	if !sourceAccountExists {
		return "", "", &tErrors.ErrorUnderfundedAccount{}
	}

	var sourceAccountTrustsDestinationAsset bool
	if !destinationAsset.IsNative() {
		_, sourceAccountTrustsDestinationAsset, _, _, _, _ = network.BlockchainAccountProperties(client, wallet.ID, destinationAsset)
	}
	_, sourceAccountTrustsInternalBalanceAsset, _, _, _, _ := network.BlockchainAccountProperties(client, wallet.ID, internalBalanceAsset)

	var ops []txnbuild.Operation = make([]txnbuild.Operation, 0)

	if !sourceAccountTrustsInternalBalanceAsset {
		// 1. establish the buyer's trustline to the internal balance token
		ops = append(ops, &txnbuild.ChangeTrust{
			Line:          txnbuild.ChangeTrustAssetWrapper{Asset: internalBalanceAsset},
			Limit:         gc.TokenLimitAsString(),
			SourceAccount: wallet.ID,
		})
		// 2. authorize that trustline
		ops = append(ops, &txnbuild.SetTrustLineFlags{
			Trustor:       wallet.ID,
			Asset:         internalBalanceAsset,
			SetFlags:      []txnbuild.TrustLineFlag{txnbuild.TrustLineAuthorized},
			SourceAccount: *countryConfig.InternalTokenIssuer,
		})
	}

	// 3. issue the internal balance token 1:1 to the buyer, in place of the on-chain stablecoin
	// payment leg - the fiat payment confirmed later by the webhook is what authorizes this mint
	ops = append(ops, &txnbuild.Payment{
		Destination:   wallet.ID,
		Amount:        swapInfo.SwapAmount,
		Asset:         internalBalanceAsset,
		SourceAccount: *countryConfig.InternalTokenIssuer,
	})

	if !destinationAsset.IsNative() && !sourceAccountTrustsDestinationAsset {
		// 4. establish trustline to the tokenized asset
		ops = append(ops, &txnbuild.ChangeTrust{
			Line:          txnbuild.ChangeTrustAssetWrapper{Asset: destinationAsset},
			Limit:         gc.TokenLimitAsString(),
			SourceAccount: wallet.ID,
		})
		ops = append(ops, &txnbuild.SetTrustLineFlags{
			Trustor:       wallet.ID,
			Asset:         txnbuild.CreditAsset{Code: swapInfo.DestinationAssetCode, Issuer: swapInfo.DestinationAssetIssuer},
			SetFlags:      []txnbuild.TrustLineFlag{txnbuild.TrustLineAuthorized},
			SourceAccount: swapInfo.DestinationAssetIssuer,
		})
	}

	//get sendPath from the internal balance token to the tokenized asset
	destAsset := ""
	if !destinationAsset.IsNative() {
		destAsset = fmt.Sprintf("%s:%s", swapInfo.DestinationAssetCode, swapInfo.DestinationAssetIssuer)
	}
	pathInput := swapModels.SwapSendPathInput{
		DestinationAssets: destAsset,
		SourceAssetCode:   internalBalanceAsset.Code,
		SourceAssetIssuer: internalBalanceAsset.Issuer,
		SourceAmount:      swapInfo.SwapAmount,
	}
	path, swappedEstimate, pathErr := GetStrictSendPaths(pathInput, client)
	if pathErr != nil {
		log.Println("[generateAssetSubscriptionFiatXdr] error fetching valid swap path ", pathErr)
		return "", "", pathErr
	}

	// 5. swap the internal balance token for the tokenized asset
	ops = append(ops, &txnbuild.PathPaymentStrictSend{
		SendAsset:     internalBalanceAsset,
		SendAmount:    swapInfo.SwapAmount,
		Destination:   wallet.ID,
		DestAsset:     destinationAsset,
		DestMin:       swapDestMin.String(),
		Path:          path,
		SourceAccount: wallet.ID,
	})

	// 6. remove the internal balance trustline so it never persists outside of this transaction
	ops = append(ops, &txnbuild.ChangeTrust{
		Line:          txnbuild.ChangeTrustAssetWrapper{Asset: internalBalanceAsset},
		Limit:         "0",
		SourceAccount: wallet.ID,
	})

	memo := fmt.Sprintf("FIAT>%v", swapInfo.DestinationAssetCode)
	swapInfo.Memo = memo

	// only pull a channel account off the pool once every step above that can fail has already
	// succeeded - on any earlier error nothing was reserved, so nothing needs to be returned
	chanAccount := <-gc.ChannelAccounts
	reserved := false
	defer func() {
		if !reserved {
			gc.ChannelAccounts <- chanAccount
		}
	}()

	_, _, _, _, chanSourceAccount, chanErr := network.BlockchainAccountProperties(client, chanAccount.Address(), txnbuild.NativeAsset{})
	if chanErr != nil {
		return "", "", chanErr
	}
	swapInfo.TransactionSource = chanSourceAccount.AccountID

	tx, txErr := txnbuild.NewTransaction(
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
	if txErr != nil {
		log.Println("[generateAssetSubscriptionFiatXdr] error constructing transaction ", txErr)
		return "", "", txErr
	}

	tx, err = tx.Sign(network.GetBlockchainNetworkPassPhrase(), chanAccount)
	if err != nil {
		log.Println("[generateAssetSubscriptionFiatXdr] error signing transaction with channel account key ", err)
		return "", "", &tErrors.ErrorTemporaryServerError{}
	}

	tx, err = tx.Sign(network.GetBlockchainNetworkPassPhrase(), internalBalanceIssuingSigners...)
	if err != nil {
		log.Println("[generateAssetSubscriptionFiatXdr] error signing transaction with internal balance issuing signers ", err)
		return "", "", &tErrors.ErrorTemporaryServerError{}
	}

	if !sourceAccountTrustsDestinationAsset {
		var tokenizationIssuerProfileWallet string
		if len(os.Getenv("TOKENIZATION_ISSUING_PROFILE_WALLET")) > 1 {
			tokenizationIssuerProfileWallet = strings.TrimSpace(os.Getenv("TOKENIZATION_ISSUING_PROFILE_WALLET"))
		}
		tokenizationIssuerProfileWalletKP := keypair.MustParseFull(tokenizationIssuerProfileWallet)
		tx, err = tx.Sign(network.GetBlockchainNetworkPassPhrase(), tokenizationIssuerProfileWalletKP)
		if err != nil {
			log.Println("[generateAssetSubscriptionFiatXdr] error signing transaction with issuer key to authorize trustline", err)
			return "", "", &tErrors.ErrorTemporaryServerError{}
		}
	}

	xdrBase64, err = tx.Base64()
	if err != nil {
		log.Println("[generateAssetSubscriptionFiatXdr] error getting txn base64", err)
		return "", "", err
	}

	gc.StoreInUseChannelAccount(chanAccount)
	reserved = true
	channelAccountPublicKey = chanAccount.Address()

	swapInfo.Messages = messages
	swapInfo.SwappedEstimate = swappedEstimate
	return xdrBase64, channelAccountPublicKey, nil
}

// getInternalBalanceIssuingSigners parses the multisig internal balance issuer's signers from the
// INTERNAL_BALANCE_ISSUING_SIGNERS env var, a CSV of secret seeds. All parsed signers are returned so
// the caller can sign with each of them - extra valid signatures beyond the account's multisig threshold
// are harmless on Stellar.
func getInternalBalanceIssuingSigners() ([]*keypair.Full, error) {
	var signers []*keypair.Full
	for _, v := range strings.Split(os.Getenv("INTERNAL_BALANCE_ISSUING_SIGNERS"), ",") {
		v = strings.TrimSpace(v)
		if len(v) == 0 {
			continue
		}
		kp, e := keypair.ParseFull(v)
		if e != nil {
			return nil, e
		}
		signers = append(signers, kp)
	}
	if len(signers) == 0 {
		return nil, errors.New("INTERNAL_BALANCE_ISSUING_SIGNERS not configured")
	}
	return signers, nil
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

	TOKENIZATION_FEE := t.GetTokenizationFeeWallet(gc)
	feeWallet := keypair.MustParseFull(TOKENIZATION_FEE.FeeWalletSecretKey)
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
	if t.AssetCountryLocation == nil {
		log.Println("[generateMintRegulatedTokenizedAssetXdr] Error locating asset country location")

		err = &tErrors.CustomError{
			Param:      "publicKey",
			Err:        "error-invalid-quote-currency",
			ErrMessage: "Tokenization does not have a valid country of location.",
			Code:       404,
		}
		return "", "", messages, issuingWallet, err
	}
	// resolve quote currency + issuer from the country's internal balance token config
	countryConfig := userModels.CountryCode(*t.AssetCountryLocation).GetConfig(gc)
	if countryConfig.InternalBalanceTokenCode == nil || countryConfig.InternalTokenIssuer == nil {
		log.Printf("[generateMintRegulatedTokenizedAssetXdr] Error locating quote currency %v\n", *t.AssetQuoteCurrency)

		err = &tErrors.CustomError{
			Param:      "publicKey",
			Err:        "error-invalid-quote-currency",
			ErrMessage: "Tokenization does not have valid tokenization currency.",
			Code:       404,
		}
		return "", "", messages, issuingWallet, err
	}
	quoteCurrency := userModels.TokenizationCurrency{AssetCode: *countryConfig.InternalBalanceTokenCode, AssetIssuer: *countryConfig.InternalTokenIssuer}

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

	if len(strings.TrimSpace(os.Getenv("INTERNAL_BALANCE_AUTHORIZER_WALLET"))) == 0 {
		err = &tErrors.CustomError{Param: "publicKey", Err: "error-internal-balance-authorizer-not-set", ErrMessage: "Internal balance token authorizer wallet not set."}
		return
	}
	authorizerKP := keypair.MustParseFull(strings.TrimSpace(os.Getenv("INTERNAL_BALANCE_AUTHORIZER_WALLET")))

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
	if err = checkDistributionWalletHasQuoteCurrencyAuthorization(quoteCurrency.AssetCode, quoteCurrency.AssetIssuer, &distributionWallet, gc); err != nil {
		return "", "", messages, issuingWallet, err
	}

	if issuingWallet.SharedAccessEnabled == 0 {

		//create sharedAccess on issuing wallet
		_, errSharedAccess := CreateSharedWalletAccess(&tokenizationIssuerUser, &tokenizationIssuerUser, &issuingWallet, &p, gc)

		if errSharedAccess != nil {
			log.Printf("[generateMintRegulatedTokenizedAssetXdr.SignBase64Txn] Error creating shared access on issuing wallets [%v], err: %v\n", tokenizationIssuerProfileWalletKP.Address(), errSharedAccess)

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
				log.Printf("[generateMintRegulatedTokenizedAssetXdr.SignBase64Txn] Error signing issuing wallet with primary wallets [%v], err: %v\n", tokenizationIssuerProfileWalletKP.Address(), e)
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

	// create + authorize distributionWallet trustline to the quote currency (internal balance token)
	ops = append(ops, &txnbuild.ChangeTrust{
		Line:          txnbuild.ChangeTrustAssetWrapper{Asset: txnbuild.CreditAsset{Code: quoteCurrency.AssetCode, Issuer: quoteCurrency.AssetIssuer}},
		Limit:         gc.TokenLimitAsString(),
		SourceAccount: distributionWallet.ID,
	})
	ops = append(ops, &txnbuild.SetTrustLineFlags{
		Trustor:       distributionWallet.ID,
		Asset:         txnbuild.CreditAsset{Code: quoteCurrency.AssetCode, Issuer: quoteCurrency.AssetIssuer},
		SetFlags:      []txnbuild.TrustLineFlag{txnbuild.TrustLineAuthorized},
		SourceAccount: quoteCurrency.AssetIssuer,
	})

	//mint the token to distribution wallet
	ops = append(ops, &txnbuild.Payment{
		Destination:   distributionWallet.ID,
		Amount:        decimal.NewFromFloat(t.NumberOfTokenToBeIssued).StringFixed(7),
		Asset:         txnbuild.CreditAsset{Code: *t.AssetCode, Issuer: *t.IssuingWalletPublicKey},
		SourceAccount: *t.IssuingWalletPublicKey,
	})
	if t.FeeInAsset > 0 {
		//deduct fee to fee wallet, from distribution wallet
		feeInAssetPayment := &txnbuild.Payment{
			Destination:   feeWallet.Address(),
			Amount:        decimal.NewFromFloat(t.FeeInAsset).StringFixed(7),
			Asset:         txnbuild.CreditAsset{Code: *t.AssetCode, Issuer: *t.IssuingWalletPublicKey},
			SourceAccount: distributionWallet.ID,
		}
		// if e := feeInAssetPayment.Validate(); e != nil {
		// 	msg := fmt.Sprintf("[generateMintRegulatedTokenizedAssetXdr] FeeInAssetPayment Operation Failed validation: %v. Fee In Asset figure: %v. Fields: %+v", e, decimal.NewFromFloat(t.FeeInAsset).StringFixed(7), *feeInAssetPayment)
		// 	gc.LogDiscordFailedRequest(msg)
		// 	err = &tErrors.CustomError{
		// 		Param:      "IssuingWalletPublicKey",
		// 		Err:        "error-could-not-approve-tokenization",
		// 		ErrMessage: "Could not approve tokenization. A fee Payment operation could not pass validation.",
		// 		Code:       404,
		// 	}
		// 	return "", "", messages, issuingWallet, err
		// }

		ops = append(ops, feeInAssetPayment)
	}

	//make market
	n, d := ToFractionInt32(t.PricePerToken)

	num := decimal.NewFromInt32(n).BigInt()
	den := decimal.NewFromInt32(d).BigInt()

	if !fitsInInt32(den) || !fitsInInt32(num) {
		//does not fit int32. return error.
		msg := fmt.Sprintf("[generateMintRegulatedTokenizedAssetXdr] Denominator Or Numerator does not fit into Int32. D: %v, N: %v", den.String(), num.String())
		gc.LogDiscordFailedRequest(msg)
		err = &tErrors.CustomError{
			Param:      "pricePerToken",
			Err:        "error-price-per-token-out-of-range",
			ErrMessage: "Current price per token cannot be correctly represented on the blockchain. Please consider raising the total number of tokens to be issued.",
			Code:       404,
		}
		return "", "", messages, issuingWallet, err
	}

	xdrInt32N := xdr.Int32(n)
	xdrInt32D := xdr.Int32(d)
	xdrPrice := xdr.Price{N: xdrInt32N, D: xdrInt32D}
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
	msg := fmt.Sprintf("[generateMintRegulatedTokenizedAssetXdr] Transaction to mint %v units (selling %v units) of %v @ %v %v generated. N: %v, D: %v. xdrN: %v, xdrD: %v. XDR Price: %+v", decimal.NewFromFloat(t.NumberOfTokenToBeIssued).StringFixed(7), decimal.NewFromFloat(t.MaxNumberOfTokenAvailableForSale).StringFixed(7), *t.AssetCode, decimal.NewFromFloat(t.PricePerToken).StringFixed(7), *t.AssetQuoteCurrency, decimal.NewFromInt32(n).String(), decimal.NewFromInt32(d).String(), xdrInt32N, xdrInt32D, xdrPrice)
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

	tx, err = tx.Sign(network.GetBlockchainNetworkPassPhrase(), chanAccount, feeWallet, authorizerKP)

	if err != nil {
		log.Println("[generateMintRegulatedTokenizedAssetXdr] error signing transaction with channelAccount, fee wallet & internal balance authorizer key ", err)
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
	if ato.TokenizationFeeID != nil {
		if *ato.TokenizationFeeID == 0 {
			//
			ato.AssetTokenizationStatus = 0
			ato.VettingStatus = 0
			gc.DB.Omit(clause.Associations).Save(&ato)
			log.Printf("[MintRegulatedTokenizedAsset] Error tokenization fee was not selected: %v\n", tokenizationID)
			err = &tErrors.CustomError{Param: "issuingWalletPublicKey", Err: "error-tokenization-fee-not-selected", ErrMessage: "Tokenization fee was not selected. Status has been returnd to allow applicant select fee."}
			return

		}
	} else {
		//
		ato.AssetTokenizationStatus = 0
		ato.VettingStatus = 0
		gc.DB.Omit(clause.Associations).Save(&ato)
		log.Printf("[MintRegulatedTokenizedAsset] Error tokenization fee was not selected: %v\n", tokenizationID)
		err = &tErrors.CustomError{Param: "issuingWalletPublicKey", Err: "error-tokenization-fee-not-selected", ErrMessage: "Tokenization fee was not selected. Status has been returned to allow applicant select fee."}
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

	//validators
	/////

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

	var NotIssuedByIssuer bool
	if ato.IssuingWalletAlias != nil && len(strings.TrimSpace(os.Getenv("TOKENIZATION_ISSUING_PROFILE"))) > 0 {
		NotIssuedByIssuer = !strings.HasPrefix(*ato.IssuingWalletAlias, strings.TrimSpace(os.Getenv("TOKENIZATION_ISSUING_PROFILE")))
	}
	e := gc.DB.Omit(clause.Associations).Save(&ato).Error
	if e != nil {
		log.Printf("[MintRegulatedTokenizedAsset] error saving tokenization to database  for %v (%v): %v\n", tokenizationID, initiator.Username, e)

		err = &tErrors.ErrorTemporaryServerError{}
		return
	}
	assetCodeExists := ato.AssetCode != nil
	if (ato.IssuingWalletPublicKey == nil || NotIssuedByIssuer) && assetCodeExists && len(os.Getenv("TOKENIZATION_ISSUING_PROFILE")) > 1 && len(os.Getenv("TOKENIZATION_ISSUING_PROFILE_WALLET")) == 56 {

		//create issuing wallet
		ato, _, err = AssignIssuingWallet(tokenizationID, gc)
		if err != nil {
			return
		}
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
	/////

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
	e = dbTX.Omit(clause.Associations).Create(&pendingAuth).Error
	if e != nil {
		log.Printf("[MintRegulatedTokenizedAsset] Error saving txn [%+v] on pending auth table: %s\n", pendingAuth, e.Error())
		err = &tErrors.ErrorTemporaryServerError{}
		return
	}
	ato.AssetTokenizationStatus = 4
	ato.DateOfApproval = time.Now()
	// ato.TokenizationTransaction = &xdrBase64
	e = dbTX.Omit(clause.Associations).Save(&ato).Error
	if e != nil {
		log.Printf("[MintRegulatedTokenizedAsset] Error saving txn on tokenizedAsset table: %s\n", e.Error())
		err = &tErrors.ErrorTemporaryServerError{}
		return
	}
	dbTX.Commit()
	//log message
	//make market
	n, d := ToFractionInt32(ato.PricePerToken)

	msg := fmt.Sprintf("Minting of %v units (selling %v units) of %v @ %v %v submitted. N: %v, D: %v", decimal.NewFromFloat(ato.NumberOfTokenToBeIssued).StringFixed(7), decimal.NewFromFloat(ato.MaxNumberOfTokenAvailableForSale).StringFixed(7), *ato.AssetCode, decimal.NewFromFloat(ato.PricePerToken).StringFixed(7), *ato.AssetQuoteCurrency, decimal.NewFromInt32(n).String(), decimal.NewFromInt32(d).String())
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
					if !strings.Contains(err.Error(), "sale") {
						//skip all sale related error
						gc.LogDiscordFailedRequest(fmt.Sprintf("[ProcessPostTokenizationTrustline] error executing first trustline call command %v, err: %v\nTrustLineInfo: [%+v]", candidate.PublicKey, err, trustLineInfo))

					}
				}
				continue
			}

			//go ahead to initiate the second call with commit
			returnedTrustLineInfo.Commit = 1

			returnedTrustLineInfo, err = TrustAsset(&signerUser, &wallet, returnedTrustLineInfo, gc)
			if err != nil {
				if err.Error() != "error-duplicate-operation-exists" {
					log.Printf("[ProcessPostTokenizationTrustline] error executing 2nd trustline call command %v,err: %v\nTrustLineInfo: [%+v]", candidate.PublicKey, err, trustLineInfo)
					if !strings.Contains(err.Error(), "sale") {
						//skip all sale related errors
						gc.LogDiscordFailedRequest(fmt.Sprintf("[ProcessPostTokenizationTrustline] error executing 2nd trustline call command %v, err: %v\nTrustLineInfo: [%+v]", candidate.PublicKey, err, trustLineInfo))
					}

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

	// countryConfig := userModels.CountryCode(*ato.AssetCountryLocation).GetConfig(gc)

	// tfa := strings.Split(countryConfig.TokenizationApplicationFeeAsset, ":") //CODE:ISSUER
	TOKENIZATION_APPLICATION_FEE := wallet.GetTokenizationApplicationFee(gc)
	if len(TOKENIZATION_APPLICATION_FEE.FeeWalletSecretKey) == 0 {
		gc.LogDiscordFailedRequest("TOKENIZATION_APPLICATION_FEE secret key not configuired in service fee table")
		return "", &tErrors.ErrorTemporaryServerError{}
	}
	feeAmount := decimal.NewFromFloat(TOKENIZATION_APPLICATION_FEE.FeeFixed)

	feeWalletPK := keypair.MustParseFull(TOKENIZATION_APPLICATION_FEE.FeeWalletSecretKey)
	assetIssuer := TOKENIZATION_APPLICATION_FEE.FeeAssetIssuer
	assetCode := TOKENIZATION_APPLICATION_FEE.FeeAssetCode

	asset := txnbuild.CreditAsset{Code: assetCode, Issuer: assetIssuer}

	sourceAccountExists, _, _, assetAccountFeeBalance, sourceAccount, _ := network.BlockchainAccountProperties(gc.BantuExpansionClient, wallet.ID, asset)

	if !sourceAccountExists {
		return "", &tErrors.CustomError{Param: "publicKey", Err: "error-account-not-activated-on-blockchain", ErrMessage: "The Wallet public key is currently underfunded. Please send about 3XBN to it to activate it before you can perform this task", Code: http.StatusBadRequest}

	}
	if assetAccountFeeBalance.LessThan(feeAmount) {
		return "", &tErrors.CustomError{Param: "publicKey", Err: "error-wallet-underfunded", ErrMessage: fmt.Sprintf("The Wallet is currently underfunded. Please maintain min %v %v balance before you can perform this task", feeAmount.String(), assetCode), Code: http.StatusBadRequest}

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

	//save application fee to the object.
	ato.TokenizationApplicationFee = TOKENIZATION_APPLICATION_FEE.FeeFixed
	ato.TokenizationApplicationFeeAsset = fmt.Sprintf("%v:%v", TOKENIZATION_APPLICATION_FEE.FeeAssetCode, TOKENIZATION_APPLICATION_FEE.FeeAssetIssuer)
	// e := gc.DB.Omit(clause.Associations).Save(ato).Error
	// if e != nil {
	// 	log.Println("[generateTokenizationFeeXdr]error saving the application fee to tokenization object", err)
	// 	logDiscordFailedRecovery("[generateTokenizationFeeXdr] error saving the application fee to tokenization object")
	// 	return "", &tErrors.ErrorTemporaryServerError{}
	// }
	return xdrBase64, nil

}

func checkDistributionWalletHasQuoteCurrencyAuthorization(assetQuoteCurrencyCode, assetQuoteCurrencyIssuer string, distributionWallet *userModels.UserWallet, gc *sharedconfig.GlobalConfig) (err error) {
	//check if the distribution wallet has the quote currency trustline and, if the quote currency requires authorization, that it is authorized

	quoteAsset := txnbuild.CreditAsset{Code: assetQuoteCurrencyCode, Issuer: assetQuoteCurrencyIssuer}
	_, trusted, _, _, dwalletSource, _ := network.BlockchainAccountProperties(gc.BantuExpansionClient, distributionWallet.ID, quoteAsset)

	if !trusted {
		//distributor wallet does not trust the quote currency yet.
		log.Printf("[checkDistributionWalletHasQuoteCurrencyAuthorization] Distribution wallet [%v] does not yet accept %v\n", distributionWallet.Alias, assetQuoteCurrencyCode)

		err = &tErrors.CustomError{
			Param:      "IssuingWalletPublicKey",
			Err:        "error-asset-trustline-not-authorized",
			ErrMessage: fmt.Sprintf("Distribution wallet %v not yet authorized to hold %v", distributionWallet.Alias, assetQuoteCurrencyCode),
			Code:       404,
		}
		return err
	}

	//get the quote currency balance entry
	for _, bal := range dwalletSource.Balances {
		if strings.EqualFold(assetQuoteCurrencyCode, bal.Code) {
			//perform the critical check

			bantuAsset := userModels.BantuAsset{
				AssetCode:   assetQuoteCurrencyCode,
				AssetIssuer: assetQuoteCurrencyIssuer,
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
						log.Printf("[checkDistributionWalletHasQuoteCurrencyAuthorization] Error: %v not authorized on distribution wallet [%v]\n", assetQuoteCurrencyCode, distributionWallet.Alias)

						err = &tErrors.CustomError{
							Param:      "IssuingWalletPublicKey",
							Err:        "error-asset-trustline-not-authorized",
							ErrMessage: fmt.Sprintf("Distribution wallet %v not yet authorized to hold %v", distributionWallet.Alias, assetQuoteCurrencyCode),
							Code:       404,
						}
						return err
					}
				} else {
					//authorization has not been given. abort process.
					log.Printf("[checkDistributionWalletHasQuoteCurrencyAuthorization] Error: %v not authorized on distribution wallet [%v]\n", assetQuoteCurrencyCode, distributionWallet.Alias)

					err = &tErrors.CustomError{
						Param:      "IssuingWalletPublicKey",
						Err:        "error-asset-trustline-not-authorized",
						ErrMessage: fmt.Sprintf("Distribution wallet %v not yet authorized to hold %v", distributionWallet.Alias, assetQuoteCurrencyCode),
						Code:       404,
					}
					return err
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

func fitsInInt32(x *big.Int) bool {
	return x.Cmp(big.NewInt(math.MaxInt32)) <= 0 && x.Cmp(big.NewInt(math.MinInt32)) >= 0
}

// parseEarlyExitPercentage parses a tokenization's free-text EarlyExitPenalty/EarlyExitFee field
// (e.g. "5" or "5%") into a float64 percentage, returning 0 when nil or not a plain number.
func parseEarlyExitPercentage(s *string) float64 {
	if s == nil {
		return 0
	}
	trimmed := strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(*s), "%"))
	if len(trimmed) == 0 {
		return 0
	}
	v, err := strconv.ParseFloat(trimmed, 64)
	if err != nil {
		return 0
	}
	return v
}

// buildTokenizedAssetEarlyExit computes the NAV/penalty/payout figures for an early exit and assembles the
// TokenizedAssetEarlyExit record. Shared between EarlyExit (building the record on the initial request) and
// ApproveTransaction (rebuilding the same record once a shared-access wallet's approvals are complete).
func buildTokenizedAssetEarlyExit(walletOwnerUsername, walletPublicKey string, input *userModels.TokenizedAssetEarlyExitInput, ta *userModels.TokenizedAsset, bank *userModels.Bank, gc *sharedconfig.GlobalConfig) (ee userModels.TokenizedAssetEarlyExit) {
	navPerToken := ta.CurrentNAVPerToken
	if navPerToken <= 0 {
		navPerToken = ta.PricePerToken
	}
	penaltyPct := parseEarlyExitPercentage(ta.EarlyExitPenalty) + parseEarlyExitPercentage(ta.EarlyExitFee)
	payoutPricePerToken := decimal.NewFromFloat(navPerToken).Mul(decimal.NewFromFloat(1).Sub(decimal.NewFromFloat(penaltyPct).Div(decimal.NewFromInt(100))))
	estimatedPayout := payoutPricePerToken.Mul(decimal.NewFromFloat(input.TokenQuantityToExit)).Truncate(7)

	ee.ID = gc.GenerateUUIDString()
	ee.TokenizedAssetID = ta.ID
	ee.WalletUsername = walletOwnerUsername
	ee.WalletPublicKey = walletPublicKey
	ee.AccountNumber = input.AccountNumber
	ee.AccountName = input.AccountName
	ee.PayoutCurrency = *ta.AssetQuoteCurrency
	ee.BankID = input.BankID
	ee.Bank = *bank
	ee.TokenQuantityToExit = input.TokenQuantityToExit
	ee.CurrentNAVPerToken = navPerToken
	ee.EarlyExitPenaltyAndFees = penaltyPct
	ee.PayoutPricePerToken, _ = payoutPricePerToken.Float64()
	ee.EstimatedPayoutAmount, _ = estimatedPayout.Float64()
	return
}

// generateEarlyExitPaymentXdr builds a plain Payment operation moving the exiting token quantity from the
// holder's wallet into the tokenization's distribution wallet (the issuing wallet's LinkedWalletPublicKey).
// There is no on-chain buy-back/liquidity for an early exit — the tokens simply return to issuer custody,
// and the holder is settled to their bank account off-chain from the persisted TokenizedAssetEarlyExit record.
func generateEarlyExitPaymentXdr(wallet *userModels.UserWallet, distributionWallet *userModels.UserWallet, ta *userModels.TokenizedAsset, amount string, multiparty int, gc *sharedconfig.GlobalConfig) (xdrBase64, transactionSource string, err error) {
	client := gc.BantuExpansionClient
	asset := txnbuild.CreditAsset{Code: strings.ToUpper(*ta.AssetCode), Issuer: strings.ToUpper(*ta.IssuingWalletPublicKey)}

	var chanAccount *keypair.Full
	var chanSourceAccount *horizon.Account
	if multiparty == 1 {
		chanAccount = <-gc.ChannelAccounts
		defer func(c *keypair.Full) {
			gc.ChannelAccounts <- c
		}(chanAccount)
		_, _, _, _, chanSourceAccount, _ = network.BlockchainAccountProperties(client, chanAccount.Address(), txnbuild.NativeAsset{})
	}

	sourceAccountExists, sourceAccountTrustsAsset, _, sourceAccountBalance, sourceAccount, sourceAccountErr := network.BlockchainAccountProperties(client, wallet.ID, asset)
	if sourceAccountErr != nil {
		return "", "", sourceAccountErr
	}
	if !sourceAccountExists {
		return "", "", &tErrors.ErrorUnderfundedAccount{}
	}
	if !sourceAccountTrustsAsset {
		return "", "", &tErrors.CustomError{Param: "walletPublicKey", Err: "error-no-trustline", ErrMessage: fmt.Sprintf("Your wallet does not hold %v.", *ta.AssetCode)}
	}

	amountDec, e := decimal.NewFromString(amount)
	if e != nil {
		return "", "", &tErrors.CustomError{Param: "tokenQuantityToExit", Err: "error-invalid-amount", ErrMessage: "Invalid token quantity."}
	}
	if sourceAccountBalance.LessThan(amountDec) {
		return "", "", &tErrors.ErrorUnderfundedAccount{Detail: fmt.Sprintf("You only have %v %v available; cannot exit %v %v.", sourceAccountBalance.String(), *ta.AssetCode, amountDec.String(), *ta.AssetCode)}
	}

	ops := []txnbuild.Operation{
		&txnbuild.Payment{
			Destination:   distributionWallet.ID,
			Amount:        amountDec.String(),
			Asset:         asset,
			SourceAccount: wallet.ID,
		},
	}

	var tx *txnbuild.Transaction
	if multiparty == 1 {
		transactionSource = chanSourceAccount.AccountID
		tx, err = txnbuild.NewTransaction(
			txnbuild.TransactionParams{
				SourceAccount:        chanSourceAccount,
				IncrementSequenceNum: true,
				Operations:           ops,
				BaseFee:              2000,
				Preconditions: txnbuild.Preconditions{
					TimeBounds: txnbuild.NewInfiniteTimeout(),
				},
				Memo: txnbuild.MemoText("EARLYEXIT"),
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
				Memo: txnbuild.MemoText("EARLYEXIT"),
			},
		)
	}
	if err != nil {
		log.Println("[generateEarlyExitPaymentXdr] error constructing transaction", err)
		return "", "", &tErrors.ErrorTemporaryServerError{}
	}

	if multiparty == 1 {
		tx, err = tx.Sign(network.GetBlockchainNetworkPassPhrase(), chanAccount)
		if err != nil {
			log.Println("[generateEarlyExitPaymentXdr] error signing transaction with channelAccount key", err)
			return "", "", &tErrors.ErrorTemporaryServerError{}
		}
	}

	xdrBase64, err = tx.Base64()
	if err != nil {
		return "", "", err
	}
	return xdrBase64, transactionSource, nil
}

// EarlyExit processes a holder's early exit (pre-maturity redemption) from a tokenized asset market fund.
// It debits the exiting token quantity from the holder's own wallet into the tokenization's distribution
// wallet via a plain Payment operation (there is no on-chain buy-back/liquidity for an early exit) and
// records the payout/settlement details for the requested bank account so it can be settled manually.
func EarlyExit(initiator *userModels.User, wallet *userModels.UserWallet, ta *userModels.TokenizedAsset, input *userModels.TokenizedAssetEarlyExitInput, gc *sharedconfig.GlobalConfig) (ee userModels.TokenizedAssetEarlyExit, err error) {
	input.TokenizedAssetID = ta.ID
	input.WalletPublicKey = wallet.ID
	input.TokenQuantityToExit = decimal.NewFromFloat(input.TokenQuantityToExit).Truncate(7).InexactFloat64()

	walletOwner, e := wallet.GetWalletOwner(gc.DB, gc)
	if e != nil {
		log.Printf("[EarlyExit] Error Unable to verify wallet owner of exiting wallet %v\n", wallet.Alias)
		err = &tErrors.CustomError{Param: "walletPublicKey", Err: "error-invalid-kyc", ErrMessage: "Unable to verify wallet owner."}
		return
	}

	if walletOwner.KYCVerified == 0 {
		log.Printf("[EarlyExit] Error Wallet owner %v has not met KYC status for asset %v\n", walletOwner.Username, *ta.AssetCode)
		err = &tErrors.CustomError{Param: "walletPublicKey", Err: "error-invalid-kyc", ErrMessage: fmt.Sprintf("%v has not passed KYC to exit this tokenized asset %v.", walletOwner.Username, *ta.AssetCode)}
		return
	}

	if ta.AssetTokenizationStatus != 5 && ta.AssetTokenizationStatus != 6 {
		log.Printf("[EarlyExit] Error Tokenized asset not open for trading yet: %v\n", ta.ID)
		err = &tErrors.CustomError{Param: "tokenizedAssetId", Err: "error-invalid-request", ErrMessage: "Only projects that are on sale or trading can accept an early exit."}
		return
	}

	if !ta.MaturityDate.IsZero() && !time.Now().Before(ta.MaturityDate) {
		log.Printf("[EarlyExit] Error Tokenized asset has already matured: %v\n", ta.ID)
		err = &tErrors.CustomError{Param: "tokenizedAssetId", Err: "error-asset-matured", ErrMessage: "This asset has reached maturity. Please use the standard redemption instead of an early exit."}
		return
	}

	if input.TokenQuantityToExit <= 0 {
		err = &tErrors.CustomError{Param: "tokenQuantityToExit", Err: "error-invalid-amount", ErrMessage: "You must specify a token quantity greater than zero to exit."}
		return
	}

	var bank userModels.Bank
	if e := gc.DB.Where("id = ?", input.BankID).First(&bank).Error; e != nil {
		log.Printf("[EarlyExit] Error locating bank %v: %v\n", input.BankID, e)
		err = &tErrors.CustomError{Param: "bankId", Err: "error-invalid-bank", ErrMessage: "Please select a valid bank for payout."}
		return
	}

	ee = buildTokenizedAssetEarlyExit(walletOwner.Username, wallet.ID, input, ta, &bank, gc)
	estimatedPayout := decimal.NewFromFloat(ee.EstimatedPayoutAmount)

	if wallet.SharedAccessEnabled == 1 && wallet.NumberOfApprovalsNeeded > 0 {
		input.Multiparty = 1
	}
	if wallet.HasViewOnlyAccess(gc) {
		input.SignatureRequired = 1
	}

	issuingWallet, e := userModels.UserWalletID(*ta.IssuingWalletPublicKey).GetWallet(gc.DB, gc)
	if e != nil {
		log.Printf("[EarlyExit] Error locating issuing wallet for tokenized asset %v: %v\n", ta.ID, e)
		err = &tErrors.CustomError{Param: "tokenizedAssetId", Err: "error-invalid-issuer", ErrMessage: "Unable to validate issuing wallet."}
		return
	}
	if issuingWallet.LinkedWalletPublicKey == nil {
		log.Printf("[EarlyExit] Error issuing wallet has no distribution wallet linked for tokenized asset %v\n", ta.ID)
		err = &tErrors.CustomError{Param: "tokenizedAssetId", Err: "error-invalid-issuer", ErrMessage: "This asset has no distribution wallet configured."}
		return
	}
	distributionWallet, e := userModels.UserWalletID(*issuingWallet.LinkedWalletPublicKey).GetWallet(gc.DB, gc)
	if e != nil {
		log.Printf("[EarlyExit] Error locating distribution wallet for tokenized asset %v: %v\n", ta.ID, e)
		err = &tErrors.CustomError{Param: "tokenizedAssetId", Err: "error-invalid-issuer", ErrMessage: "Unable to validate distribution wallet."}
		return
	}

	dbTX := gc.DB.Begin()
	defer dbTX.Rollback()

	if input.Multiparty == 0 {
		e := dbTX.Omit(clause.Associations).Create(&ee).Error
		if e != nil {
			log.Printf("[EarlyExit] error saving tokenized asset early exit to database [%+v] for %v: %v\n", ee, walletOwner.Username, e)
			err = &tErrors.ErrorTemporaryServerError{}
			return
		}
	}

	if len(input.TransactionSignature) == 0 {
		xdrBase64, transactionSource, e := generateEarlyExitPaymentXdr(wallet, &distributionWallet, ta, decimal.NewFromFloat(input.TokenQuantityToExit).Truncate(7).String(), input.Multiparty, gc)
		if e != nil {
			log.Printf("[EarlyExit] error generating early exit xdr [%+v] for %v: %v\n", ee, walletOwner.Username, e)
			err = e
			return
		}
		input.Transaction = xdrBase64
		input.TransactionSource = transactionSource
		input.Memo = "EARLYEXIT"
	}

	input.NetworkPassPhrase = network.GetBlockchainNetworkPassPhrase()

	if len(input.TransactionSignature) == 0 && input.Commit == 0 {
		err = nil
		return ee, nil
	}

	client := gc.BantuExpansionClient

	if len(input.TransactionSignature) > 0 && (input.Commit == 0 || wallet.HasViewOnlyAccess(gc)) {
		txnHash, e := network.SubmitXdrWithSignature(client, initiator.PrimarySigner, input.Transaction, input.TransactionSignature)
		if e != nil {
			log.Printf("[EarlyExit] error submitting early exit transaction to blockchain [%+v] for %v: %v\n", ee, walletOwner.Username, e)
			err = e
			return
		}
		input.TransactionID = txnHash
		ee.TransactionID = txnHash
		dbTX.Omit(clause.Associations).Save(&ee)
		dbTX.Commit()
		wallet.InvalidateUserCache(gc)
		initiator.InvalidateUserWalletCache(gc)
		return ee, nil
	}

	//multi party
	if input.Multiparty == 1 {
		input.TransactionID = "PENDING_AUTH"
		ee.TransactionID = "PENDING_AUTH"

		id := uuid.NewString()
		description := fmt.Sprintf("Early exit of Tokenized asset [%v]\nQuantity: %v %v,\nEstimated payout: %v %v", *ta.AssetName, decimal.NewFromFloat(input.TokenQuantityToExit).String(), *ta.AssetCode, estimatedPayout.String(), *ta.AssetQuoteCurrency)
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
			Initiator:                initiator.Username,
			InitiatorSignerPublicKey: initiator.PrimarySigner,
			WalletPublicKey:          wallet.ID,
			TransactionType:          "TOKENIZED ASSET EARLY EXIT",
			Description:              description,
			TransactionSource:        input.TransactionSource,
			ApprovalsNeeded:          wallet.NumberOfApprovalsNeeded,
			TransactionXdr:           input.Transaction,
			TransactionInfoStr:       &transactionStr,
		}
		e := dbTX.Omit(clause.Associations).Create(&pendingAuth).Error
		if e != nil {
			log.Printf("[EarlyExit] Error saving early exit txn [%+v] transaction on pending auth table: %s\n", pendingAuth, e.Error())
			err = &tErrors.ErrorTemporaryServerError{}
			return
		}
		dbTX.Omit(clause.Associations).Save(&ee)
		dbTX.Commit()
		wallet.InvalidateUserCache(gc)
		initiator.InvalidateUserWalletCache(gc)
		return ee, nil
	}

	log.Println("[EarlyExit]UNKNOWN OPTION FOR ACTION")
	err = &tErrors.ErrorTemporaryServerError{}
	return
}
