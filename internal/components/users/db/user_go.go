package users

import (
	"errors"
	"log"
	"strings"
	userModels "trovo-wallet-api/internal/components/users/models"
	conDB "trovo-wallet-api/internal/db"
	tErrors "trovo-wallet-api/internal/errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

//GetUserInfo gets user data
func GetUserInfo(userInfo string, db *gorm.DB) (user userModels.User, err error) {
	// var wallet usermodels.UserWallet
	conDB.PrintDBStats("GetUserInfo", db)

	//e returns execution errors
	var e error
	if len(userInfo) == 56 {
		//56 char public key is supplied

		subQuery := db.Table("user_wallets").Where("id = ?", userInfo).Or("temp_public_key = ?", &userInfo).Select("user_id")
		e = db.Preload(clause.Associations).Where("id = (?)", subQuery).First(&user).Error
	} else if strings.Contains(userInfo, "@") {
		//email is supplied
		e = db.Preload(clause.Associations).First(&user, userModels.User{Email: strings.ToLower(userInfo)}).Error
	} else {
		//username is supplied
		// e = db.First(&user, usermodels.User{Username: strings.ToLower(userInfo)}).Error
		subQuery := db.Table("user_wallets").Where("alias = ?", strings.ToLower(userInfo)).Select("user_id")
		e = db.Preload(clause.Associations).Where("id = (?)", subQuery).First(&user).Error
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
