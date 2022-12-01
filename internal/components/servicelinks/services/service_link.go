package servicelinks

import (
	"errors"
	"log"
	"net/http"
	"strings"
	"time"
	servicelinkModels "trovo-wallet-api/internal/components/servicelinks/models"
	usersDB "trovo-wallet-api/internal/components/users/db"
	userModels "trovo-wallet-api/internal/components/users/models"

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
func GetUserFromPrimarySigner(signerKey string, db *gorm.DB, gc *sharedconfig.GlobalConfig) (user userModels.User, err error) {
	return usersDB.GetUserFromPrimarySigner(signerKey, gc.DB, gc)
}
