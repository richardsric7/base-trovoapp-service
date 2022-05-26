package users

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"
	"trovo-wallet-api/internal/cache"
	merchants "trovo-wallet-api/internal/components/merchants/services"
	users "trovo-wallet-api/internal/components/users/db"
	usermodels "trovo-wallet-api/internal/components/users/models"
	bantupayErrors "trovo-wallet-api/internal/errors"

	"github.com/ecnepsnai/discord"
	"github.com/nyaruka/phonenumbers"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

//RegisterUser registers user information
func RegisterUser(userInfo usermodels.UserRegistrationInfo, db *gorm.DB, dynamicLinkServiceUrlChan chan string, redisCache *cache.RedisCache) (usermodels.UserRegistrationInfo, bool, error) {
	discord.WebhookURL = "https://discord.com/api/webhooks/824381163367170058/OXSX51RHd9DyLFbFipjdW3yXmyYC8SWwqd6HiXl6UtDzu75RxS1LzWA800hWereJJumw"
	if len(os.Getenv("REGISTRATION_ERROR_WEBHOOK")) > 50 {
		discord.WebhookURL = os.Getenv("REGISTRATION_ERROR_WEBHOOK")
	}
	// if banned, errBanned := users.PublicKeyIsBanned(userInfo.PublicKey, db); banned {
	// 	return userInfo, false, errBanned
	// }
	if _, errExists := users.PublicKeyAlreadyExists(userInfo.PublicKey, db); errExists != nil {
		return userInfo, false, errExists
	}
	if len(userInfo.Mobile) > 0 {
		geoData, _ := usermodels.GetGeoInfo(userInfo.PublicIP)
		num, err := phonenumbers.Parse(userInfo.Mobile, geoData.CountryCode)
		if err == nil {
			mobile := fmt.Sprintf("+%v-%v", *num.CountryCode, *num.NationalNumber)
			userInfo.Mobile = mobile
		}
	}
	NormalizeUserRegistrationInfo(&userInfo)

	errValidation := ValidateUserRegistrationInfo(userInfo)

	if errValidation != nil {
		discord.Say(fmt.Sprintf("[RegisterUser]  Validation failed for user:%v, Error:%v", userInfo.Username, errValidation))
		return userInfo, false, errValidation
	}

	user, dbErrors := users.UserRegistrationDbChecks(userInfo, db)

	if dbErrors != nil {
		discord.Say(fmt.Sprintf("[RegisterUser] Registration DB check failed for user:%v, Error:%v", userInfo.Username, dbErrors))
		return userInfo, false, dbErrors
	}

	{

		//do mail validation
		validationResult, blockEmail, _ := user.VerifyEmailOnMailgun()
		if blockEmail {
			discord.Say(fmt.Sprintf("[RegisterUser] Mailgun Mail Validation failed for user:%v, Result:%+v", userInfo.Username, validationResult))
			return userInfo, false, &bantupayErrors.ErrorEmailFailedValidation{Email: user.Email, Detail: fmt.Sprintf("Email validation failed for [%v]. Please put valid email and try again later.", userInfo.Email)}
		}

		emailSent, expectedCode, checkAndSendErr := CheckAndSendVerificationCode(userInfo)

		if checkAndSendErr != nil {
			log.Printf("error processing verification %s\n", checkAndSendErr)
			discord.Say(fmt.Sprintf("[RegisterUser] verification code process failed for user:%v, Error:%v, EnteredCode: %v, ExpectedCode:%v", userInfo.Username, checkAndSendErr, userInfo.VerificationCode, expectedCode))
			return userInfo, emailSent, checkAndSendErr
		}

		if checkAndSendErr == nil && emailSent {
			//no error, but email still sent
			return userInfo, emailSent, checkAndSendErr
		}

	}

	if os.Getenv("REGISTRATION_THROTTLE_PER_IP") != "" && os.Getenv("REGISTRATION_THROTTLE_PER_IP") != "0" {
		//check if same IP dat registered a user is up to 1hr
		var createdAt time.Time
		qry := "select created_at from users where public_ip = ? order by created_at DESC limit 1"
		e := db.Raw(qry, user.PublicIP).Scan(&createdAt).Error
		if e == nil {
			//record was found...check the time
			t := decimal.RequireFromString(os.Getenv("REGISTRATION_THROTTLE_PER_IP")).IntPart()
			if time.Since(createdAt) < (time.Duration(t) * time.Second) {
				//it is less than 1hr since registration from same IP, reject registration
				discord.Say(fmt.Sprintf("Too many registration from the IP %+v\n", user))
				return userInfo, false, &bantupayErrors.CustomError{
					Param:      "username",
					Err:        "error-user-registration-failed-validation",
					ErrMessage: "Your registration validation failed at this time.",
				}
			}
		}
	}

	//add geo information
	user.AppendGeoInfo()
	//all checks have passed

	//add referralLink
	if data, e := merchants.GenerateReferralLink(user.Username, dynamicLinkServiceUrlChan, redisCache); e == nil {
		user.ReferralLink = &data.DynamicLink
		user.ReferralQrCode = &data.QRCode
	}

	//save the user

	errCreate := db.Create(user).Error
	if errCreate != nil {

		discord.Say(fmt.Sprintf("[RegisterUser] user creation failed for user:%v, with DB Error:%v\n\n\nFailedData:%+v", userInfo.Username, errCreate, userInfo))
		return userInfo, false, errors.New("unable to create user due to error in information")
	}
	return userInfo, false, nil

}
