package servicelinks

import (
	"encoding/base64"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
	servicelinkModels "trovo-wallet-api/internal/components/servicelinks/models"
	usersDB "trovo-wallet-api/internal/components/users/db"
	userModels "trovo-wallet-api/internal/components/users/models"
	"trovo-wallet-api/internal/evmkeypair"

	conDB "trovo-wallet-api/internal/db"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/sharedconfig"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// GetMerchantInfo gets merchant data
func GetServiceInfo(mInfo, apikey string, db *gorm.DB) (user servicelinkModels.ServiceLink, err error) {

	conDB.PrintDBStats("GetMerchantInfo", db)

	//e returns execution errors
	var e error
	if len(mInfo) == 56 {
		//56 char publick key is supplied
		e = db.Where(servicelinkModels.ServiceLink{PublicKey: mInfo, ApiKey: apikey}).First(&user).Error
	} else {
		//username is supplied
		e = db.First(&user, servicelinkModels.ServiceLink{OwnerUsername: strings.ToLower(mInfo), ApiKey: apikey}).Error
	}

	if e != nil {
		if errors.Is(e, gorm.ErrRecordNotFound) {
			//no user was found
			err = &tErrors.ErrorServiceDoesNotExist{Username: mInfo}
			return
		}
		log.Println("[GetMerchantInfo] error: ", e)
		err = &tErrors.ErrorTemporaryServerError{}
		return

	}

	// log.Printf("user for %v is %v\n", userInfo, user)
	return user, nil

}

// GetServiceLinkByAPIKey gets merchant data by API key
func GetServiceLinkByAPIKey(apiKey string, db *gorm.DB) (serviceLink servicelinkModels.ServiceLink, err error) {

	conDB.PrintDBStats("GetServiceLinkByAPIKey", db)

	//e returns execution errors

	e := db.Where("api_key = ?", apiKey).First(&serviceLink).Error

	if e != nil {
		if errors.Is(e, gorm.ErrRecordNotFound) {
			//no user was found
			err = &tErrors.CustomError{Param: "apiKey", Err: "error-api-key-invalid", ErrMessage: "apiKey is invalid."}
			return
		}
		log.Println("[GetServiceLinkByAPIKey] error: ", e)
		err = &tErrors.ErrorTemporaryServerError{}
		return

	}

	return serviceLink, nil

}

// GetLoginSession gets login data
func GetLoginSession(mInfo, walletInfo, loginID string, db *gorm.DB) (loginSession servicelinkModels.ServiceLinkLoginSession, err error) {

	conDB.PrintDBStats("GetLoginSession", db)

	//e returns execution errors
	e := db.Where("id = ?", loginID).Where("owner_username = ?", mInfo).Where("wallet_Username = ?", walletInfo).First(&loginSession).Error

	if e != nil {
		if errors.Is(e, gorm.ErrRecordNotFound) {
			//no user was found
			err = &tErrors.ErrorLoginSessionDoesNotExist{Username: walletInfo}
			return
		}
		log.Println("[GetLoginSession] error: ", e)
		err = &tErrors.ErrorTemporaryServerError{}
		return

	}

	return loginSession, nil

}

// GetUserAuthorization gets user authorization data
func GetUserAuthorizationData(mInfo, walletInfo, authID string, db *gorm.DB) (authData servicelinkModels.ServiceLinkAuthorization, err error) {

	conDB.PrintDBStats("GetUserAuthorizationData", db)

	e := db.Where("owner_username = ?", mInfo).Where("wallet_Username = ?", walletInfo).Where("id = ?", authID).First(&authData).Error

	if e != nil {

		if errors.Is(e, gorm.ErrRecordNotFound) {
			//no user was found
			err = &tErrors.ErrorAuthorizationDoesNotExist{Username: walletInfo}
			return
		}
		log.Println("[GetUserAuthorizationData] error: ", e)
		err = &tErrors.ErrorTemporaryServerError{}
		return

	}

	return authData, nil

}

// GetRewardOnlyAuthorizationData gets user authorization data
func GetRewardOnlyAuthorizationData(mInfo, authID string, db *gorm.DB) (authData servicelinkModels.ServiceLinkAuthorization, err error) {

	conDB.PrintDBStats("GetRewardOnlyAuthorizationData", db)

	e := db.Where("owner_username = ?", mInfo).Where("id = ?", authID).First(&authData).Error

	if e != nil {

		if errors.Is(e, gorm.ErrRecordNotFound) {
			//no user was found
			err = &tErrors.CustomError{
				Param:      authID,
				Err:        "error reward data does not exist",
				ErrMessage: "Reward/Airdrop is invalid or it has expired.",
				Code:       http.StatusNotFound,
			}
			return
		}
		log.Println("[GetRewardOnlyAuthorizationData] error: ", e)
		err = &tErrors.ErrorTemporaryServerError{}
		return

	}

	if authData.Authorized == 1 {

		//no user was found
		err = &tErrors.CustomError{
			Param:      authID,
			Err:        "error: reward data has expired",
			ErrMessage: "Reward/Airdrop has expired.",
			Code:       http.StatusNotFound,
		}
		return

	}

	return authData, nil

}

// GetEventAuthorizationData gets user authorization data
func GetEventAuthorizationData(mInfo, eventID string, db *gorm.DB) (eventData servicelinkModels.ServiceLinkEvent, err error) {

	conDB.PrintDBStats("GetEventAuthorizationData", db)

	e := db.Where("owner_username = ?", mInfo).Where("id = ?", eventID).First(&eventData).Error

	if e != nil {
		if errors.Is(e, gorm.ErrRecordNotFound) {
			//no user was found
			err = &tErrors.CustomError{
				Param:      eventID,
				Err:        "error: link data does not exist",
				ErrMessage: "Event is either invalid or closed or expired.",
				Code:       http.StatusNotFound,
			}
			return
		}
		log.Println("[GetEventAuthorizationData] error: ", e)

		err = &tErrors.ErrorTemporaryServerError{}
		return

	}

	if eventData.ExpiresAt.Before(time.Now()) {
		log.Println("[GetEventAuthorizationData] error: auth already expired")
		//no user was found
		err = &tErrors.CustomError{
			Param:      eventID,
			Err:        "error link event has expired",
			ErrMessage: "Event is no longer valid.",
			Code:       http.StatusNotFound,
		}
		return

	}

	return eventData, nil

}

// GetUserForServiceLink gets user information
func GetUserForServiceLink(ID string, serviceLink servicelinkModels.ServiceLink, db *gorm.DB, gc *sharedconfig.GlobalConfig) (userInfoForServiceLink userModels.ServiceLinksUser, err error) {
	conDB.PrintDBStats("GetUserForServiceLink", gc.DB)
	user, err := usersDB.GetUser(ID, db, gc)

	if err != nil {
		return userInfoForServiceLink, err
	}

	userInfoForServiceLink = user.ToServiceLinkUser(gc)

	if userInfoForServiceLink.Suspended == 1 {
		return userInfoForServiceLink, &tErrors.ErrorUsernameIsSuspended{}
	}

	return userInfoForServiceLink, nil

}

// GetTransactionSignature gets signed transaction from unsigned transaction information
func GetTransactionSignature(input *servicelinkModels.ServiceLinkTokenizedAssetAuthRequestInput, gc *sharedconfig.GlobalConfig) (output servicelinkModels.ServiceLinkTokenizedAssetAuthRequestInput, err error) {
	if gc.IsValidTokenizedAsset(input.AssetCode) {
		t := gc.GetTokenizedAssetByCode(input.AssetCode)
		if t.AssetTokenizationStatus < 6 {
			return output, &tErrors.CustomError{
				Param:      "assetIssuer",
				Err:        "error-asset-not-yet-available-for-sale",
				ErrMessage: "This tokenized Asset is not yet available for secondary market. Authorization is not allowed at this time.",
				Code:       http.StatusForbidden,
			}
		}

		//start signing transaction
		log.Println("[generateTrustAssetXdr] <<<<<<<<<<<<<<<<<<<<<<<<<<<< signing transaction with issuer key>>>>>>>>>>>>>>>>>>>>>>>>")
		//get atprofile
		var tokenizationIssuerProfileWallet string

		if len(os.Getenv("TOKENIZATION_ISSUING_PROFILE_WALLET")) > 1 {
			tokenizationIssuerProfileWallet = strings.TrimSpace(os.Getenv("TOKENIZATION_ISSUING_PROFILE_WALLET"))
		}

		tokenizationIssuerProfileWalletKP := evmkeypair.MustParseFull(tokenizationIssuerProfileWallet)

		// input.UnsignedTransaction is a base64-encoded digest to sign
		// (see internal/middleware/security_checks.go's SignBase64Txn
		// doc for why - there is no Base equivalent of parsing/re-hashing
		// a Stellar XDR envelope; the digest is computed once upstream
		// when the transaction is built).
		digest, e := base64.StdEncoding.DecodeString(input.UnsignedTransaction)
		if e != nil {
			return output, &tErrors.ErrorInvalidTransaction{}
		}

		signature, e := tokenizationIssuerProfileWalletKP.SignBase64(digest)
		if e != nil {
			log.Println("[GetTransactionSignature] error signing transaction with issuer key to authorize trustline", e)
			return output, &tErrors.ErrorTemporaryServerError{}
		}

		//assign signature
		input.SignedTransaction = signature
		return *input, nil

	}
	return output, &tErrors.CustomError{
		Param:      "assetIssuer",
		Err:        "error-asset-not-valid",
		ErrMessage: "This asset is not a tokenized Asset.",
		Code:       http.StatusForbidden,
	}

}
func GetUserFromPrimarySigner(signerKey string, db *gorm.DB, gc *sharedconfig.GlobalConfig) (user userModels.User, err error) {
	return usersDB.GetUserFromPrimarySigner(signerKey, gc.DB, gc)
}

func UpdateUserKYCStatus(targetUser *userModels.User, kycStatus int, jsonString string, gc *sharedconfig.GlobalConfig) (err error) {
	dbTx := gc.DB.Begin()
	defer dbTx.Rollback()

	targetUser.KYCVerified = kycStatus
	e := dbTx.Omit(clause.Associations).Save(targetUser).Error
	if e != nil {
		log.Printf("[UpdateUserKYCStatus] error saving user KYC %v\n", e)
		return &tErrors.ErrorTemporaryServerError{}

	}
	jstruc := userModels.KycWebhookRequest{
		ServiceProvider: *targetUser.CreatedByServiceLinkID,
		Data:            jsonString,
	}
	e = dbTx.Omit(clause.Associations).Save(&jstruc).Error
	if e != nil {
		log.Printf("[UpdateUserKYCStatus] error saving user KYC Json string%v\n", e)
		return &tErrors.ErrorTemporaryServerError{}

	}
	dbTx.Commit()
	return nil
}
