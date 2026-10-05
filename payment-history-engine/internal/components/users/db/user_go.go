package users

import (
	"errors"
	"fmt"
	"log"
	"strings"
	userModels "trovo-wallet-payment-history-engine/internal/components/users/models"
	conDB "trovo-wallet-payment-history-engine/internal/db"
	tErrors "trovo-wallet-payment-history-engine/internal/errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// GetUser gets a user by wallet address or signer, wallet alias, or user ID,
// username, mobile or email
func GetUser(userInfo string, db *gorm.DB) (user userModels.User, err error) {
	conDB.PrintDBStats("GetUserInfo", db)

	//e returns execution errors
	var e error
	if len(userInfo) == 42 {
		//42 char address is supplied

		subQuery := db.Table("user_wallets").Where("id = ?", userInfo).Or("signer = ?", userInfo).Select("user_id")
		e = db.Preload(clause.Associations).Where("id = (?)", subQuery).First(&user).Error
	} else if strings.Contains(userInfo, "_") {
		//alias format is supplied
		subQuery := db.Table("user_wallets").Where("alias = ?", strings.ToLower(userInfo)).Select("user_id")
		e = db.Preload(clause.Associations).Where("id = (?)", subQuery).First(&user).Error

	} else {
		//search by ID and phone number, username, email

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

// GetWallet gets a wallet by its address or alias
func GetWallet(identifier string, db *gorm.DB) (user userModels.UserWallet, err error) {
	conDB.PrintDBStats("GetUserInfo", db)

	//e returns execution errors
	var e error
	if len(identifier) == 42 {
		//42 char address is supplied
		e = db.Preload(clause.Associations).Where("id = ?", identifier).First(&user).Error
		if e == nil {
			return
		}
	} else {
		//username is supplied
		e = db.Preload(clause.Associations).Where("alias = ?", strings.ToLower(identifier)).First(&user).Error

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
	return user, nil

}
