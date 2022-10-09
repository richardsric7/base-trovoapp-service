package users

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	userModels "trovo-wallet-api/internal/components/users/models"
	conDB "trovo-wallet-api/internal/db"
	tErrors "trovo-wallet-api/internal/errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// GetUser gets user data by either wallet id or signer or temporary public key
func GetUser(userInfo string, db *gorm.DB) (user userModels.User, err error) {
	conDB.PrintDBStats("GetUserInfo", db)

	//e returns execution errors
	var e error
	if len(userInfo) == 56 {
		//56 char public key is supplied

		subQuery := db.Table("user_wallets").Where("id = ?", userInfo).Or("temp_public_key = ?", &userInfo).Or("signer = ?", userInfo).Select("user_id")
		e = db.Preload(clause.Associations).Where("id IN (?)", subQuery).First(&user).Error
	} else if strings.Contains(userInfo, "_") {
		//alias format is supplied
		subQuery := db.Table("user_wallets").Where("alias = ?", strings.ToLower(userInfo)).Select("user_id")
		e = db.Preload(clause.Associations).Where("id = (?)", subQuery).First(&user).Error

	} else {
		

		e = db.Preload(clause.Associations).Where("id = ?", userInfo).Or("username = ?", strings.ToLower(userInfo)).Or("mobile = ?", &userInfo).Or("email = ?", userInfo).First(&user).Error
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
	return user, nil

}

// GetWallet gets user wallet data by alias or public key or temp public key
func GetWallet(identifier string, db *gorm.DB) (userWallet userModels.UserWallet, temp bool, err error) {
	conDB.PrintDBStats("GetUserInfo", db)

	//e returns execution errors
	var e error
	if len(identifier) == 56 {
		//56 char public key is supplied
		e = db.Preload(clause.Associations).Where("id = ?", identifier).Or("temp_public_key = ?", &identifier).First(&userWallet).Error
		if e == nil {
			if identifier == *userWallet.TempPublicKey {
				temp = true
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
	db.Preload(clause.Associations).Where("public_key = ?", publicKey).Find(&accessList)

	return
}

func UpdatePushNotificationToken(identifier string, pnt *string, db *gorm.DB) {
	user, _ := GetUser(identifier, db)

	if user.PushNotificationToken != pnt {
		user.PushNotificationToken = pnt
		err := db.Save(&user).Error
		if err != nil {
			log.Printf("[UpdatePushNotificationToken] unable to update push notification token for user [%v], due to:[%v]", user.Username, err)
		}
	}
}

// GetUserFromPrimarySigner fetches the user linked to the primary signer
func GetUserFromPrimarySigner(publicKey string, db *gorm.DB) (user userModels.User, err error) {

	publicKey = strings.TrimSpace(publicKey)
	// var user usermodels.User
	if err := db.Where("primary_signer = ?", strings.ToUpper(strings.ReplaceAll(publicKey, " ", ""))).First(&user).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return user, &tErrors.ErrorTemporaryServerError{}
		}
		return user, &tErrors.CustomError{Param: "primarySigner",
			Err:        "error primary signer does not exist",
			ErrMessage: "PrimarySigner does not exist",
			Code:       http.StatusNotFound,
		}
	}
	// discord.Say(fmt.Sprintf("[PublicKeyIsBanned] publicKey: %v is banned\n", publicKey))

	return user, nil

}
