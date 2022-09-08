package merchants

import (
	"errors"
	"log"
	"strings"
	userModels "trovo-wallet-api/internal/components/merchants/models"
	conDB "trovo-wallet-api/internal/db"
	tErrors "trovo-wallet-api/internal/errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// GetMerchantUserInfo gets user data
func GetMerchantUserInfo(userInfo string, db *gorm.DB) (user userModels.User, err error) {
	// var wallet usermodels.UserWallet
	conDB.PrintDBStats("GetMerchantUserInfo", db)

	//e returns execution errors
	var e error
	if len(userInfo) == 56 {
		//56 char public key is supplied

		e = db.Preload(clause.Associations).Where("public_key = ?", userInfo).First(&user).Error
	} else if strings.Contains(userInfo, "@") {
		//email is supplied
		e = db.Preload(clause.Associations).First(&user, userModels.User{Email: strings.ToLower(userInfo)}).Error
	} else {
		//username is supplied
		e = db.Preload(clause.Associations).Where("username = ?", userInfo).First(&user).Error
	}

	if e != nil {
		if errors.Is(e, gorm.ErrRecordNotFound) {
			//no user was found
			err = &tErrors.ErrorUserDoesNotExist{Username: userInfo}
			return
		}
		log.Println("[GetMerchantUserInfo] error: ", e)
		err = &tErrors.ErrorTemporaryServerError{}
		return

	}
	// log.Printf("[GetMerchantUserInfo] %+v\n", user)
	return user, nil

}
