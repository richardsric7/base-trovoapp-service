package users

import (
	users "admin-panel-dashboard/internal/components/users/db"
	models "admin-panel-dashboard/internal/models"
	"errors"
	"fmt"
	"os"

	"github.com/ecnepsnai/discord"
	"github.com/knadh/smtppool"
	"github.com/nyaruka/phonenumbers"
	"gorm.io/gorm"
)

// RegisterUser registers user information
func RegisterUser(userInfo models.UserRegistrationInfo, db *gorm.DB, pool *smtppool.Pool) (models.UserRegistrationInfo, bool, error) {
	discord.WebhookURL = "https://discord.com/api/webhooks/824381163367170058/OXSX51RHd9DyLFbFipjdW3yXmyYC8SWwqd6HiXl6UtDzu75RxS1LzWA800hWereJJumw"
	if len(os.Getenv("REGISTRATION_ERROR_WEBHOOK")) > 50 {
		discord.WebhookURL = os.Getenv("REGISTRATION_ERROR_WEBHOOK")
	}
	if banned, errBanned := users.PublicKeyIsBanned(userInfo.PublicKey, db); banned {
		return userInfo, false, errBanned
	}
	if len(userInfo.Mobile) > 0 {
		geoData, _ := models.GetGeoInfo(userInfo.PublicIP)
		num, err := phonenumbers.Parse(userInfo.Mobile, geoData.CountryCode)
		if err == nil {
			mobile := fmt.Sprintf("+%v-%v", *num.CountryCode, *num.NationalNumber)
			userInfo.Mobile = mobile
		}
	}
	NormalizeUserRegistrationInfo(&userInfo)

	errValidation := ValidateUserRegistrationInfo(userInfo)

	if errValidation != nil {
		err := discord.Say(fmt.Sprintf("[RegisterUser]  Validation failed for user:%v, Error:%v", userInfo.Username, errValidation))
		if err != nil {
			return models.UserRegistrationInfo{}, false, err
		}
		return userInfo, false, errValidation
	}

	user, dbErrors := users.UserRegistrationDbChecks(userInfo, db)

	if dbErrors != nil {
		err := discord.Say(fmt.Sprintf("[RegisterUser] Registration DB check failed for user:%v, Error:%v", userInfo.Username, dbErrors))
		if err != nil {
			return models.UserRegistrationInfo{}, false, err
		}
		return userInfo, false, dbErrors
	}

	{

		// //do mail validation
		// validationResult, blockEmail, _ := user.VerifyEmailOnMailgun()
		// if blockEmail {
		// 	discord.Say(fmt.Sprintf("[RegisterUser] Mailgun Mail Validation failed for user:%v, Result:%+v", userInfo.Username, validationResult))
		// 	return userInfo, false, &bantupayErrors.ErrorEmailFailedValidation{Email: user.Email, Detail: fmt.Sprintf("Email validation failed for [%v]. Please put valid email and try again later.", userInfo.Email)}
		// }

	}

	// add geo information
	user.AppendGeoInfo()
	// all checks have passed
	//save the user
	result := db.Create(user)
	if result.Error != nil {

		err := discord.Say(fmt.Sprintf("[RegisterUser] user creation failed for user:%v, with DB Error:%v\n\n\nFailedData:%+v", userInfo.Username, result.Error, userInfo))
		if err != nil {
			return models.UserRegistrationInfo{}, false, err
		}
		return userInfo, false, errors.New("unable to create user due to error in information")
	}
	return userInfo, false, nil

}
