package users

import (
	trovowalletErrors "admin-panel-dashboard/internal/errors"
	models "admin-panel-dashboard/internal/models"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ecnepsnai/discord"
	"github.com/google/uuid"
	"github.com/nyaruka/phonenumbers"
	"gorm.io/gorm"
)

// UserRegistrationInfoToUser populates user information with registration information
func UserRegistrationInfoToUser(userInfo models.UserRegistrationInfo, user *models.User) {
	user.Username = strings.TrimSpace(strings.ToLower(userInfo.Username))
	// user.Email = strings.TrimSpace(strings.ToLower(userInfo.Email))
	user.ID = uuid.NewString()
	user.FirstName = strings.TrimSpace(strings.ToUpper(userInfo.FirstName))
	user.LastName = strings.TrimSpace(strings.ToUpper(userInfo.LastName))
	// if len(userInfo.MiddleName) > 0 {
	//	v := strings.TrimSpace(strings.ToUpper(userInfo.MiddleName))
	//	user.MiddleName = &v
	//}
	if len(userInfo.Mobile) > 0 {
		userInfo.Mobile = strings.TrimSpace(userInfo.Mobile)
		geoData, _ := user.GetGeoInfo()
		num, err := phonenumbers.Parse(userInfo.Mobile, geoData.CountryCode)
		if err == nil {
			mobile := fmt.Sprintf("+%v-%v", *num.CountryCode, *num.NationalNumber)
			user.Mobile = mobile
		} else {

			user.Mobile = userInfo.Mobile
		}
	}

	user.PublicIP = userInfo.PublicIP
}

// UserRegistrationDbChecks checks validaty of user data
func UserRegistrationDbChecks(userInfo models.UserRegistrationInfo, db *gorm.DB) (*models.User, error) {
	discord.WebhookURL = "https://discord.com/api/webhooks/824381163367170058/OXSX51RHd9DyLFbFipjdW3yXmyYC8SWwqd6HiXl6UtDzu75RxS1LzWA800hWereJJumw"
	if len(os.Getenv("REGISTRATION_ERROR_WEBHOOK")) > 50 {
		discord.WebhookURL = os.Getenv("REGISTRATION_ERROR_WEBHOOK")
	}
	// db, err := conDB.OpenDb()
	// if err != nil {
	// 	log.Println("-------------- ------DB error in UserUpdateDbChecks:", err)
	// 	return nil, &trovowalletErrors.ErrorTemporaryServerError{}
	// }
	var user models.User
	// check if username is reserved
	_, err := UsernameIsReserved(strings.ToLower(userInfo.Username), db)
	if err != nil {
		return nil, err
	}

	// apply rules engine 1

	// Check if mobile already exists.
	err = db.Where("mobile = ?", userInfo.Mobile).First(&user).Error
	if err == nil {
		return nil, &trovowalletErrors.ErrorMobileNumberAlreadyExists{Detail: fmt.Sprintf("mobile number [%v] already exists with another account", userInfo.Mobile)}

	}
	// Check if username already exists.
	err = db.Where("username = ?", strings.ToLower(userInfo.Username)).Or("email = ?", strings.ToLower(userInfo.Email)).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {

			if len(userInfo.Referrer) > 2 {
				// Check if referrer exists
				var referrer models.User
				refError := db.Where("username = ?", userInfo.Referrer).First(&referrer).Error
				if refError != nil {
					// referrer record was not found
					err = discord.Say(fmt.Sprintf("[UserRegistrationDbChecks] Referrer Validation failed for user:%v, referrer: %v", userInfo.Username, userInfo.Referrer))
					if err != nil {
						return nil, err
					}
					return nil, &trovowalletErrors.ErrorInvalidReferrer{}
				}

			}
			UserRegistrationInfoToUser(userInfo, &user)
			// append the Rules and score

			return &user, nil
		}
		log.Println("[UserRegistrationDbChecks]", err)
		err = discord.Say(fmt.Sprintf("[UserRegistrationDbChecks] DB check failed with unknown error for user:%v, error: %v", userInfo.Username, err))
		if err != nil {
			return nil, err
		}

		return nil, &trovowalletErrors.ErrorTemporaryServerError{}

	}

	// if user.Email == userInfo.Email {
	// 	return nil, &trovowalletErrors.ErrorEmailAlreadyExists{Detail: fmt.Sprintf("email [%v] already exists with another account", userInfo.Email)}
	// }

	return nil, &trovowalletErrors.ErrorUsernameAlreadyExists{Detail: userInfo.Username + " has already been taken by another user"}

}

// UserUpdateDbChecks checks validaty of user data
func UserUpdateDbChecks(owner string, userInfo models.UserUpdateInfo, db *gorm.DB) error {
	// db, err := conDB.OpenDb()
	// if err != nil {
	// 	log.Println("-------------- ------DB error in UserUpdateDbChecks:", err)
	// 	return &trovowalletErrors.ErrorTemporaryServerError{}
	// }
	discord.WebhookURL = "https://discord.com/api/webhooks/824381163367170058/OXSX51RHd9DyLFbFipjdW3yXmyYC8SWwqd6HiXl6UtDzu75RxS1LzWA800hWereJJumw"
	if len(os.Getenv("REGISTRATION_ERROR_WEBHOOK")) > 50 {
		discord.WebhookURL = os.Getenv("REGISTRATION_ERROR_WEBHOOK")
	}
	// var user models.User
	// Check if username already exists.
	// if os.Getenv("ALLOW_UPDATE_OF_BUDS_MOBILE_NUMBER") == "1" {
	// 	err := db.Where("mobile = ?", strings.ToLower(userInfo.Mobile)).Where("username != ?", owner).First(&user).Error
	// 	if err != nil {
	// 		if !errors.Is(err, gorm.ErrRecordNotFound) {
	// 			log.Println("[UserUpdateDbChecks]", err)
	// 			return &trovowalletErrors.ErrorTemporaryServerError{}
	// 		}
	// 	} else {
	// 		if len(userInfo.Mobile) > 0 && user.Mobile != nil {
	// 			if *user.Mobile == userInfo.Mobile {
	// 				return &trovowalletErrors.ErrorMobileNumberAlreadyExists{Detail: "mobile number you are trying to update to, already exists with another account"}
	// 			}
	// 		}
	// 	}
	// }

	// if user.Email == userInfo.Email {
	// 	return nil, &trovowalletErrors.ErrorEmailAlreadyExists{}
	// }

	return nil

}

// UsernameIsReserved check is name is reserved. Status = 0 means not available (reserved). Status = 1 means available
func UsernameIsReserved(username string, db *gorm.DB) (reserved bool, err error) {
	// db, err := conDB.OpenDb()
	// if err != nil {
	// 	log.Println("-------------- ------DB error in RESERVED NAME:", err)
	// 	return
	// }
	username = strings.TrimSpace(username)
	var reservedName models.ReservedName
	if err := db.Where("reserved_name = ? AND status = 0", strings.ToLower(strings.ReplaceAll(username, " ", ""))).First(&reservedName).Error; err != nil {

		return false, nil
	}

	return true, &trovowalletErrors.ErrorUsernameIsReserved{}
}

// PublicKeyIsBanned check if public key is banned.
func PublicKeyIsBanned(publicKey string, db *gorm.DB) (banned bool, err error) {
	// discord.WebhookURL = "https://discord.com/api/webhooks/824381163367170058/OXSX51RHd9DyLFbFipjdW3yXmyYC8SWwqd6HiXl6UtDzu75RxS1LzWA800hWereJJumw"
	// if len(os.Getenv("IMPORT_ERROR_WEBHOOK")) > 50 {
	// 	discord.WebhookURL = os.Getenv("IMPORT_ERROR_WEBHOOK")
	// }
	publicKey = strings.TrimSpace(publicKey)
	var bannedPublicKey models.BannedPublicKey
	if err := db.Where("public_key = ?", strings.ToLower(strings.ReplaceAll(publicKey, " ", ""))).First(&bannedPublicKey).Error; err != nil {

		return false, nil
	}
	// discord.Say(fmt.Sprintf("[PublicKeyIsBanned] publicKey: %v is banned\n", publicKey))

	return true, &trovowalletErrors.ErrorAccountIsBanned{}
}
