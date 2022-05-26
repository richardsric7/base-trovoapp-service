package merchants

import (
	"errors"
	"log"
	"strings"
	"trovo-wallet-api/internal/cache"
	merchantmodels "trovo-wallet-api/internal/components/merchants/models"
	conDB "trovo-wallet-api/internal/db"
	tErrors "trovo-wallet-api/internal/errors"

	"gorm.io/gorm"
)

//GetMerchantInfo gets merchant data
func GetMerchantInfo(mInfo string, db *gorm.DB) (user merchantmodels.Merchant, err error) {

	conDB.PrintDBStats("GetMerchantInfo", db)

	//e returns execution errors
	var e error
	if len(mInfo) == 56 {
		//56 char publick key is supplied
		e = db.Where(merchantmodels.Merchant{PublicKey: mInfo}).First(&user).Error
	} else if strings.Contains(mInfo, "@") {
		//email is supplied
		e = db.First(&user, merchantmodels.Merchant{Email: strings.ToLower(mInfo)}).Error
	} else {
		//username is supplied
		e = db.First(&user, merchantmodels.Merchant{BantupayUsername: strings.ToLower(mInfo)}).Error
	}

	if e != nil {
		if errors.Is(e, gorm.ErrRecordNotFound) {
			//no user was found
			err = &tErrors.ErrorMerchantDoesNotExist{Username: mInfo}
			return
		}
		log.Println("[GetMerchantInfo] error: ", e)
		err = &tErrors.ErrorTemporaryServerError{}
		return

	}

	// log.Printf("user for %v is %v\n", userInfo, user)
	return user, nil

}

//GetLoginSession gets merchant data
func GetLoginSession(mInfo, walletInfo, loginID string, db *gorm.DB) (loginSession merchantmodels.MerchantLoginSession, err error) {

	conDB.PrintDBStats("GetLoginSession", db)
	// db, err := conDB.OpenDb()
	// if err != nil {
	// 	log.Println("-------------- ------DB error in getUserInfo:", err)
	// 	return
	// }
	// var user usermodels.User
	//Check if username already exists.

	//e returns execution errors
	e := db.Where("id = ?", loginID).Where("merchant_username = ?", mInfo).Where("wallet_Username = ?", walletInfo).First(&loginSession).Error

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

//GetUserAuthorization gets user authorization data
func GetUserAuthorizationData(mInfo, walletInfo, authID string, db *gorm.DB) (authData merchantmodels.MerchantAuthorization, err error) {

	conDB.PrintDBStats("GetUserAuthorizationData", db)

	e := db.Where("merchant_username = ?", mInfo).Where("wallet_Username = ?", walletInfo).Where("id = ?", authID).First(&authData).Error

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

//GetRewardOnlyAuthorizationData gets user authorization data
func GetRewardOnlyAuthorizationData(mInfo, authID string, db *gorm.DB) (authData merchantmodels.MerchantAuthorization, err error) {

	conDB.PrintDBStats("GetRewardOnlyAuthorizationData", db)

	e := db.Where("merchant_username = ?", mInfo).Where("wallet_Username = ?", mInfo).Where("id = ?", authID).First(&authData).Error

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

//GetEventAuthorizationData gets user authorization data
func GetEventAuthorizationData(mInfo, authID string, db *gorm.DB) (authData merchantmodels.MerchantAuthorization, err error) {

	conDB.PrintDBStats("GetEventAuthorizationData", db)

	e := db.Where("merchant_username = ?", mInfo).Where("wallet_Username = ?", mInfo).Where("id = ?", authID).First(&authData).Error

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

//GetUserForMerchants gets user information
func GetUserForMerchants(ID string, merchant merchantmodels.Merchant, db *gorm.DB, dynamicLinkServiceUrlChan chan string, redisCache *cache.RedisCache) (userForMerchantInfo merchantmodels.MerchantBudsInfo, err error) {
	conDB.PrintDBStats("GetUserForMerchants", db)

	// user, err := users.GetUserInfo(ID, db)

	// if err != nil {
	// 	return userForMerchantInfo, err
	// }
	// if user.Suspended == 1 {
	// 	return userForMerchantInfo, &tErrors.ErrorUsernameIsSuspended{}
	// }

	return userForMerchantInfo, nil

}
