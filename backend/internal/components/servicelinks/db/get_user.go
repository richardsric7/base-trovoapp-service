package servicelinks

import (
	"errors"
	"net/http"
	"strings"
	usersDB "trovo-wallet-api/internal/components/users/db"
	userModels "trovo-wallet-api/internal/components/users/models"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/sharedconfig"

	"gorm.io/gorm"
)

// GetServiceLinkUserInfo gets user data
func GetServiceLinkUserInfo(userInfo string, db *gorm.DB, gc *sharedconfig.GlobalConfig) (user userModels.ServiceLinksUser, err error) {
	// var wallet usermodels.UserWallet
	u, err := usersDB.GetUser(userInfo, gc.DB, gc)

	if err != nil {
		return user, err
	}

	return u.ToServiceLinkUser(gc), nil

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

	return user, nil

}
