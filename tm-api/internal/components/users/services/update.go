package users

import (
	usersdb "admin-panel-dashboard/internal/components/users/db"
	"admin-panel-dashboard/internal/errors"
	models "admin-panel-dashboard/internal/models"

	"gorm.io/gorm"
)

// UpdateUser registers user information
func UpdateUser(identifier string, authenticatedAddress string, userInfo models.UserUpdateInfo, publicIP string, db *gorm.DB) (models.UserUpdateInfo, error) {

	user, err := usersdb.GetUserInfo(identifier, db)

	if err != nil {
		return userInfo, err
	}
	if user.Suspended == 1 {
		return userInfo, &errors.ErrorUsernameIsSuspended{}
	}

	userUpdateInfoToUser(userInfo, &user)

	// user.AppendGeoInfo()
	// all checks have passed
	//save the user
	result := db.Save(user)

	return userInfo, result.Error

}

// UpdateUserKYC registers user information
func UpdateUserKYC(identifier string, authenticatedAddress string, userInfo models.KYCUpdateInfo, publicIP string, db *gorm.DB) (models.KYCUpdateInfo, error) {

	user, err := usersdb.GetUserInfo(identifier, db)

	if err != nil {
		return userInfo, err
	}
	if user.Suspended == 1 {
		return userInfo, &errors.ErrorUsernameIsSuspended{}
	}

	userKYCUpdateInfoToUser(userInfo, &user)

	// save the user
	//if userInfo.KYCLevel == 0 {
	//	user.Offline = 1
	//}
	result := db.Save(user)
	if userInfo.KYCLevel == 0 {
		db.Raw("UPDATE offers SET offline = 1 WHERE maker = ?", user.Username)
	}

	return userInfo, result.Error

}

func userKYCUpdateInfoToUser(userInfo models.KYCUpdateInfo, user *models.User) {

	// user.KYCLevel = userInfo.KYCLevel

}

func userUpdateInfoToUser(userInfo models.UserUpdateInfo, user *models.User) {

	// if len(userInfo.ContactPhone) > 9 {
	//	user.ContactPhone = &userInfo.ContactPhone
	//}

}
