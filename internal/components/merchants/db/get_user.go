package merchants

import (
	"errors"
	"log"
	"strings"
	userModels "trovo-wallet-api/internal/components/users/models"
	tErrors "trovo-wallet-api/internal/errors"

	"gorm.io/gorm"
)

//GetUserInfo gets user data
func GetUserInfo(userInfo string, db *gorm.DB) (user userModels.User, err error) {

	// conDB.PrintDBStats("GetUserInfo", db)
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
		e = db.Where(userModels.User{PublicKey: userInfo}).Or(userModels.User{TempPublicKey: &userInfo}).First(&user).Error
	} else if strings.Contains(userInfo, "@") {
		//email is supplied
		e = db.First(&user, userModels.User{Email: strings.ToLower(userInfo)}).Error
	} else {
		//username is supplied
		e = db.First(&user, userModels.User{Username: strings.ToLower(userInfo)}).Error
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
