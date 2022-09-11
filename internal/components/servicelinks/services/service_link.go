package servicelinks

import (
	"errors"
	"log"
	"strings"
	merchantdb "trovo-wallet-api/internal/components/servicelinks/db"
	servicelinkModels "trovo-wallet-api/internal/components/servicelinks/models"
	conDB "trovo-wallet-api/internal/db"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/sharedconfig"

	"gorm.io/gorm"
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

	e := db.Where(servicelinkModels.ServiceLink{ApiKey: apiKey}).First(&serviceLink).Error

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

	e := db.Where("owner_username = ?", mInfo).Where("wallet_Username = ?", mInfo).Where("id = ?", authID).First(&authData).Error

	if e != nil {
		if errors.Is(e, gorm.ErrRecordNotFound) {
			//no user was found
			err = &tErrors.CustomError{
				Param:      authID,
				Err:        "error reward data does not exist",
				ErrMessage: "Reward/Airdrop is invalid or it has expired.",
			}
			return
		}
		log.Println("[GetRewardOnlyAuthorizationData] error: ", e)
		err = &tErrors.ErrorTemporaryServerError{}
		return

	}

	if authData.Authorized == 1 {
		if errors.Is(e, gorm.ErrRecordNotFound) {
			//no user was found
			err = &tErrors.CustomError{
				Param:      authID,
				Err:        "error: reward data has expired",
				ErrMessage: "Reward/Airdrop has expired.",
			}
			return
		}
		log.Println("[GetRewardOnlyAuthorizationData] error: ", e)
		err = &tErrors.ErrorTemporaryServerError{}
		return

	}

	return authData, nil

}

// GetEventAuthorizationData gets user authorization data
func GetEventAuthorizationData(mInfo, authID string, db *gorm.DB) (authData servicelinkModels.ServiceLinkAuthorization, err error) {

	conDB.PrintDBStats("GetEventAuthorizationData", db)

	e := db.Where("owner_username = ?", mInfo).Where("wallet_Username = ?", mInfo).Where("id = ?", authID).First(&authData).Error

	if e != nil {
		if errors.Is(e, gorm.ErrRecordNotFound) {
			//no user was found
			err = &tErrors.CustomError{
				Param:      authID,
				Err:        "error: event data does not exist",
				ErrMessage: "Event registration is either invalid or closed or expired.",
			}
			return
		}
		log.Println("[GetEventAuthorizationData] error: ", e)
		err = &tErrors.ErrorTemporaryServerError{}
		return

	}

	if authData.Authorized == 1 {
		if errors.Is(e, gorm.ErrRecordNotFound) {
			//no user was found
			err = &tErrors.CustomError{
				Param:      authID,
				Err:        "error: event data has expired",
				ErrMessage: "Event registration has closed/expired.",
			}
			return
		}
		log.Println("[GetEventAuthorizationData] error: ", e)
		err = &tErrors.ErrorTemporaryServerError{}
		return

	}

	return authData, nil

}

// GetUserForServiceLink gets user information
func GetUserForServiceLink(ID string, serviceLink servicelinkModels.ServiceLink, gc *sharedconfig.GlobalConfig) (userInfoForServiceLink servicelinkModels.User, err error) {
	conDB.PrintDBStats("GetUserForServiceLink", gc.DB)

	userInfoForServiceLink, err = merchantdb.GetServiceLinkUserInfo(ID, gc.DB)

	if err != nil {
		return userInfoForServiceLink, err
	}
	if userInfoForServiceLink.Suspended == 1 {
		return userInfoForServiceLink, &tErrors.ErrorUsernameIsSuspended{}
	}

	return userInfoForServiceLink, nil

}
