package users

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	usermodels "trovo-wallet-api/internal/components/users/models"
	tErrors "trovo-wallet-api/internal/errors"

	"github.com/ecnepsnai/discord"
	"github.com/google/uuid"
	"github.com/nyaruka/phonenumbers"
	"gorm.io/gorm"
)

//UserRegistrationInfoToUser populates user information with registration information
func UserRegistrationInfoToUser(userInfo usermodels.UserRegistrationInfo, user *usermodels.User) {
	user.PublicKey = strings.TrimSpace(strings.ToUpper(userInfo.PublicKey))
	user.Username = strings.TrimSpace(strings.ToLower(userInfo.Username))
	user.Email = strings.TrimSpace(strings.ToLower(userInfo.Email))
	user.ID = uuid.NewString()
	user.FirstName = strings.TrimSpace(strings.ToUpper(userInfo.FirstName))
	user.LastName = strings.TrimSpace(strings.ToUpper(userInfo.LastName))
	if len(userInfo.PushNotificationToken) > 0 {
		user.PushNotificationToken = &userInfo.PushNotificationToken
	}

	if len(userInfo.Mobile) > 0 {
		userInfo.Mobile = strings.TrimSpace(userInfo.Mobile)
		// geoData, _ := user.GetGeoInfo()
		num, err := phonenumbers.Parse(userInfo.Mobile, userInfo.MobileCountryCode)
		if err == nil {
			mobile := fmt.Sprintf("+%v-%v", *num.CountryCode, *num.NationalNumber)
			user.Mobile = &mobile
		} else {

			user.Mobile = &userInfo.Mobile
		}
	}
	if len(userInfo.Referrer) > 2 {
		//Check if referrer exists
		user.Referrer = &userInfo.Referrer
	}

	user.PublicIP = userInfo.PublicIP
}

//UserRegistrationDbChecks checks validaty of user data
func UserRegistrationDbChecks(userInfo usermodels.UserRegistrationInfo, db *gorm.DB) (*usermodels.User, error) {
	discord.WebhookURL = "https://discord.com/api/webhooks/824381163367170058/OXSX51RHd9DyLFbFipjdW3yXmyYC8SWwqd6HiXl6UtDzu75RxS1LzWA800hWereJJumw"
	if len(os.Getenv("REGISTRATION_ERROR_WEBHOOK")) > 50 {
		discord.WebhookURL = os.Getenv("REGISTRATION_ERROR_WEBHOOK")
	}
	// db, err := conDB.OpenDb()
	// if err != nil {
	// 	log.Println("-------------- ------DB error in UserUpdateDbChecks:", err)
	// 	return nil, &tErrors.ErrorTemporaryServerError{}
	// }
	var user usermodels.User
	//check if username is reserved
	_, err := UsernameIsReserved(strings.ToLower(userInfo.Username), db)
	if err != nil {
		return nil, err
	}
	//
	_, err = PublicKeyAlreadyExists(userInfo.PublicKey, db)
	if err != nil {
		return nil, err
	}

	//Check if mobile already exists.
	if len(userInfo.Mobile) > 0 {
		log.Printf("-------------- ------Checking if mobile [%s] already exists\n", userInfo.Mobile)
		err = db.Where("mobile = ?", userInfo.Mobile).First(&user).Error
		if err == nil {
			return nil, &tErrors.ErrorMobileNumberAlreadyExists{Detail: fmt.Sprintf("mobile number [%v] already exists with another account", userInfo.Mobile)}

		}
	}
	//Check if username already exists.
	err = db.Where("username = ?", strings.ToLower(userInfo.Username)).Or("email = ?", strings.ToLower(userInfo.Email)).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {

			if len(userInfo.Referrer) > 2 {
				//Check if referrer exists
				var referrer usermodels.User
				refError := db.Where("username = ?", userInfo.Referrer).First(&referrer).Error
				if refError != nil {
					//referrer record was not found
					discord.Say(fmt.Sprintf("[UserRegistrationDbChecks] Referrer Validation failed for user:%v, referrer: %v", userInfo.Username, userInfo.Referrer))
					return nil, &tErrors.ErrorInvalidReferrer{}
				}

			}
			UserRegistrationInfoToUser(userInfo, &user)
			//append the Rules and score

			return &user, nil
		}
		log.Println("[UserRegistrationDbChecks]", err)
		discord.Say(fmt.Sprintf("[UserRegistrationDbChecks] DB check failed with unknown error for user:%v, error: %v", userInfo.Username, err))

		return nil, &tErrors.ErrorTemporaryServerError{}

	}

	if user.Email == userInfo.Email {
		return nil, &tErrors.ErrorEmailAlreadyExists{Detail: fmt.Sprintf("email [%v] already exists with another account", userInfo.Email)}
	}

	// if len(userInfo.Instagram) > 0 && user.Instagram != nil {
	// 	if *user.Instagram == userInfo.Instagram {
	// 		return nil, &tErrors.ErrorInstagramHandleAlreadyExists{}
	// 	}
	// }

	// if len(userInfo.Twitter) > 0 && user.Twitter != nil {
	// 	if *user.Twitter == userInfo.Twitter {
	// 		return nil, &tErrors.ErrorTwitterHandleAlreadyExists{}
	// 	}
	// }

	// if len(userInfo.Telegram) > 0 && user.Telegram != nil {
	// 	if *user.Telegram == userInfo.Telegram {
	// 		return nil, &tErrors.ErrorTelegramHandleAlreadyExists{}
	// 	}

	// }

	return nil, &tErrors.ErrorUsernameAlreadyExists{Detail: userInfo.Username + " has already been taken by another user"}

}

//UserUpdateDbChecks checks validaty of user data

//UsernameIsReserved check is name is reserved. Status = 0 means not available (reserved). Status = 1 means available
func UsernameIsReserved(username string, db *gorm.DB) (reserved bool, err error) {
	// db, err := conDB.OpenDb()
	// if err != nil {
	// 	log.Println("-------------- ------DB error in RESERVED NAME:", err)
	// 	return
	// }
	username = strings.TrimSpace(username)
	var reservedName usermodels.ReservedName
	if err := db.Where("reserved_name = ? AND status = 0", strings.ToLower(strings.ReplaceAll(username, " ", ""))).First(&reservedName).Error; err != nil {

		return false, nil
	}

	return true, &tErrors.ErrorUsernameIsReserved{}
}

//PublicKeyIAlreadyExists check if public key already exists
func PublicKeyAlreadyExists(publicKey string, db *gorm.DB) (exists bool, err error) {
	// discord.WebhookURL = "https://discord.com/api/webhooks/824381163367170058/OXSX51RHd9DyLFbFipjdW3yXmyYC8SWwqd6HiXl6UtDzu75RxS1LzWA800hWereJJumw"
	// if len(os.Getenv("IMPORT_ERROR_WEBHOOK")) > 50 {
	// 	discord.WebhookURL = os.Getenv("IMPORT_ERROR_WEBHOOK")
	// }
	publicKey = strings.TrimSpace(publicKey)
	var userWallet usermodels.UserWallet
	if err := db.Where("id = ?", strings.ToUpper(strings.ReplaceAll(publicKey, " ", ""))).First(&userWallet).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return true, &tErrors.ErrorTemporaryServerError{}
		}
		return false, nil
	}
	// discord.Say(fmt.Sprintf("[PublicKeyIsBanned] publicKey: %v is banned\n", publicKey))

	return true, &tErrors.CustomError{Param: "publicKey", Err: "error-public-key-already-exists", ErrMessage: fmt.Sprintf("Bantu Address [%v] already exists with another active account", userWallet.ID)}

}
