package users

import (
	"errors"
	"log"
	"strings"
	usermodels "trovo-wallet-api/internal/components/users/models"
	conDB "trovo-wallet-api/internal/db"
	bantupayerrors "trovo-wallet-api/internal/errors"

	"gorm.io/gorm"
)

//GetUserInfo gets user data
func GetUserInfo(userInfo string, db *gorm.DB) (user usermodels.User, err error) {

	conDB.PrintDBStats("GetUserInfo", db)
	// db, err := conDB.OpenDb()
	// if err != nil {
	// 	log.Println("-------------- ------DB error in getUserInfo:", err)
	// 	return
	// }
	// var user usermodels.User
	//Check if username already exists.

	//e returns execution errors
	var e error
	if len(userInfo) == 56 {
		//56 char publick key is supplied
		e = db.Where(usermodels.User{PublicKey: userInfo}).Or(usermodels.User{TempPublicKey: &userInfo}).First(&user).Error
	} else if strings.Contains(userInfo, "@") {
		//email is supplied
		e = db.First(&user, usermodels.User{Email: strings.ToLower(userInfo)}).Error
	} else {
		//username is supplied
		e = db.First(&user, usermodels.User{Username: strings.ToLower(userInfo)}).Error
	}

	if e != nil {
		if errors.Is(e, gorm.ErrRecordNotFound) {
			//no user was found
			err = &bantupayerrors.ErrorUserDoesNotExist{Username: userInfo}
			return
		}
		log.Println("[GetUserInfo] error: ", e)
		err = &bantupayerrors.ErrorTemporaryServerError{}
		return

	}

	// log.Printf("user for %v is %v\n", userInfo, user)
	return user, nil

}
