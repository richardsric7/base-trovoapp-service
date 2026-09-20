package users

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	userModels "trovo-wallet-api/internal/components/users/models"
	conDB "trovo-wallet-api/internal/db"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/sharedconfig"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// GetUser gets user data by either wallet id or signer or temporary public key
func GetUser(userInfo string, db *gorm.DB, gc *sharedconfig.GlobalConfig) (user userModels.User, err error) {
	conDB.PrintDBStats("GetUserInfo", db)
	cacheKeyInfo := fmt.Sprintf("userObj %v", userInfo)
	userInfo = strings.TrimSpace(userInfo)

	{

		// search cache for balance
		ok, rawdata := gc.RedisCache.GetCachedResultRaw(cacheKeyInfo)

		if ok {

			// log.Printf("GetUserFromPrimarySigner[%v], served from cache\n", cacheKeyInfo)
			json.Unmarshal(rawdata, &user)
			if len(user.UserWallets) > 0 {
				return
			}

		}

	}
	//e returns execution errors
	var e error
	if len(userInfo) == 42 {
		//56 char public key is supplied

		subQuery := db.Table("user_wallets").Where("id = ?", userInfo).Or("temp_address = ?", &userInfo).Or("signer = ?", userInfo).Select("user_id")
		e = db.Preload("UserWallets.Permissions").Preload(clause.Associations).Where("id IN (?)", subQuery).First(&user).Error
	} else if strings.Contains(userInfo, "_") {
		//alias format is supplied
		subQuery := db.Table("user_wallets").Where("alias = ?", strings.ToLower(userInfo)).Select("user_id")
		e = db.Preload("UserWallets.Permissions").Preload(clause.Associations).Where("id = (?)", subQuery).First(&user).Error

	} else {

		// e = db.Preload("UserWallets.Permissions").Preload(clause.Associations).Where("id = ?", userInfo).Or("username = ?", strings.ToLower(userInfo)).Or("mobile = ?", &userInfo).Or("email = ?", userInfo).First(&user).Error
		e = db.Preload("UserWallets.Permissions").Preload(clause.Associations).Where("(id = ? OR lower(username) = lower(?) OR mobile = ? OR lower(email) = lower(?))", userInfo, userInfo, userInfo, userInfo).First(&user).Error
	}

	if e != nil {
		if errors.Is(e, gorm.ErrRecordNotFound) {
			//no user was found
			err = &tErrors.ErrorUserDoesNotExist{Username: userInfo}
			return
		}
		log.Println("[GetUserInfo] error: ", e)
		err = &tErrors.ErrorTemporaryServerError{}
		return

	}

	// log.Printf("user for %v is %v\n", userInfo, user)
	gc.RedisCache.StoreResultToCacheRaw(cacheKeyInfo, user, 2000)
	cacheKeyUsername := fmt.Sprintf("userObj %v", user.Username)
	gc.RedisCache.StoreResultToCacheRaw(cacheKeyUsername, user, 2000)
	cacheKeyEmail := fmt.Sprintf("userObj %v", user.Email)
	gc.RedisCache.StoreResultToCacheRaw(cacheKeyEmail, user, 2000)
	cacheKeySigner := fmt.Sprintf("userObj %v", user.PrimarySigner)
	gc.RedisCache.StoreResultToCacheRaw(cacheKeySigner, user, 2000)
	cacheKeyUserID := fmt.Sprintf("userObj %v", user.ID)
	gc.RedisCache.StoreResultToCacheRaw(cacheKeyUserID, user, 2000)

	return user, nil

}

// GetWallet gets user wallet data by alias or public key or temp public key
func GetWallet(identifier string, db *gorm.DB) (userWallet userModels.UserWallet, temp bool, err error) {
	conDB.PrintDBStats("GetUserInfo", db)

	//e returns execution errors
	var e error
	if len(identifier) == 42 {
		//56 char public key is supplied
		e = db.Preload(clause.Associations).Where("id = ?", identifier).Or("temp_address = ?", &identifier).First(&userWallet).Error
		if e == nil {
			if userWallet.TempAddress != nil {
				if identifier == *userWallet.TempAddress {
					temp = true
				}
			}

			return
		}
	} else {
		//username is supplied
		e = db.Preload(clause.Associations).Where("alias = ?", strings.ToLower(identifier)).First(&userWallet).Error

	}

	if e != nil {
		if errors.Is(e, gorm.ErrRecordNotFound) {
			//no user was found
			err = &tErrors.CustomError{Param: "publicKey", Err: "error-wallet-does-not-exist", ErrMessage: fmt.Sprintf("%v is not assigned to any wallet", identifier)}
			return
		}
		log.Println("[GetUserInfo] error: ", e)
		err = &tErrors.ErrorTemporaryServerError{}
		return

	}

	// log.Printf("user for %v is %v\n", userInfo, user)
	return userWallet, temp, nil

}
func GetPermissionList(publicKey string, db *gorm.DB) (accessList []userModels.WalletPermission) {
	accessList = make([]userModels.WalletPermission, 0)
	db.Preload(clause.Associations).Where("wallet_address = ?", publicKey).Find(&accessList)

	return
}

func UpdatePushNotificationToken(identifier string, pnt *string, db *gorm.DB, gc *sharedconfig.GlobalConfig) {
	user, _ := GetUser(identifier, db, gc)

	if user.PushNotificationToken != pnt {
		user.PushNotificationToken = pnt
		err := db.Omit(clause.Associations).Save(&user).Error
		if err != nil {
			log.Printf("[UpdatePushNotificationToken] unable to update push notification token for user [%v], due to:[%v]", user.Username, err)
		}
		user.InvalidateUserCache(gc)
	}
}

// GetUserFromPrimarySigner fetches the user linked to the primary signer
func GetUserFromPrimarySigner(publicKey string, db *gorm.DB, gc *sharedconfig.GlobalConfig) (user userModels.User, err error) {
	cacheKeySigner := fmt.Sprintf("userObj %v", publicKey)

	{

		// search cache for balance
		ok, rawdata := gc.RedisCache.GetCachedResultRaw(cacheKeySigner)

		if ok {

			// log.Printf("GetUserFromPrimarySigner[%v], served from cache\n", cacheKeySigner)
			json.Unmarshal(rawdata, &user)
			return
		}

	}
	publicKey = strings.TrimSpace(publicKey)
	// var user usermodels.User
	e := db.Preload("UserWallets.Permissions").Preload(clause.Associations).Where("primary_signer = ?", strings.ToUpper(strings.ReplaceAll(publicKey, " ", ""))).First(&user).Error
	if e != nil {
		if !errors.Is(e, gorm.ErrRecordNotFound) {
			return user, &tErrors.ErrorTemporaryServerError{}
		}

		return user, &tErrors.CustomError{Param: "primarySigner",
			Err:        "error primary signer does not exist",
			ErrMessage: "PrimarySigner does not exist",
			Code:       http.StatusNotFound,
		}
	}

	// discord.Say(fmt.Sprintf("[AddressIsBanned] publicKey: %v is banned\n", publicKey))
	gc.RedisCache.StoreResultToCacheRaw(cacheKeySigner, user, 2000)
	cacheKeyUsername := fmt.Sprintf("userObj %v", user.Username)
	gc.RedisCache.StoreResultToCacheRaw(cacheKeyUsername, user, 2000)

	cacheKeyEmail := fmt.Sprintf("userObj %v", user.Email)
	gc.RedisCache.StoreResultToCacheRaw(cacheKeyEmail, user, 2000)

	cacheKeyUserID := fmt.Sprintf("userObj %v", user.ID)
	gc.RedisCache.StoreResultToCacheRaw(cacheKeyUserID, user, 2000)
	return user, nil

}

func InvalidateUserCache(id string, gc *sharedconfig.GlobalConfig) {
	userAccount, err := GetUser(id, gc.DB, gc)
	if err != nil {
		return
	}
	cacheKey1 := fmt.Sprintf("GetBalance_%s", userAccount.Address)
	cacheKeyUsername := fmt.Sprintf("userObj %v", userAccount.Username)
	cacheKeyEmail := fmt.Sprintf("userObj %v", userAccount.Email)
	cacheKeySigner := fmt.Sprintf("userObj %v", userAccount.PrimarySigner)
	cacheKeyUserID := fmt.Sprintf("userObj %v", userAccount.ID)
	gc.RedisCache.DeleteFromCache(cacheKeyUsername, cacheKeyEmail, cacheKeySigner, cacheKeyUserID)

	gc.RedisCache.DeleteFromCache(cacheKey1)
	InvalidateUserWalletCache(&userAccount, gc)
}

func InvalidateUserWalletCache(userAccount *userModels.User, gc *sharedconfig.GlobalConfig) {

	if userAccount == nil {
		return
	}
	if userAccount.UserWallets == nil {
		return
	}
	if len(userAccount.UserWallets) == 0 {
		return
	}
	for _, w := range userAccount.UserWallets {
		cacheKey1 := fmt.Sprintf("GetBalance_%s", w.ID)
		cacheKey2 := fmt.Sprintf("GetBalance_%s", func() string {
			if w.TempAddress != nil {
				return *w.TempAddress
			} else {
				return "nil"
			}
		}())

		cacheKey3 := fmt.Sprintf("userObj %v", w.Alias)
		cacheKey4 := fmt.Sprintf("userObj %v", w.ID)
		cacheKeySigner := fmt.Sprintf("userObj %v", w.Signer)
		cacheKeyUserID := fmt.Sprintf("userObj %v", w.UserID)
		cacheKeyWalletAlias := fmt.Sprintf("walletObj_%v", w.Alias)
		cacheKeyWalletID := fmt.Sprintf("walletObj_%v", w.ID)
		gc.RedisCache.DeleteFromCache(cacheKeyWalletAlias, cacheKeyWalletID, cacheKey1, cacheKey2, cacheKey3, cacheKey4, cacheKeySigner, cacheKeyUserID)

	}

}
