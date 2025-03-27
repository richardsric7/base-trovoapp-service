package users

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"
	users "trovo-wallet-api/internal/components/users/db"
	userModels "trovo-wallet-api/internal/components/users/models"
	dl "trovo-wallet-api/internal/dynamiclinks"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/ecnepsnai/discord"
	"github.com/nyaruka/phonenumbers"
	"github.com/shopspring/decimal"
	"gorm.io/gorm/clause"
)

// RegisterUser registers user information
func RegisterUser(userInfo userModels.UserRegistrationInfo, gc *sharedconfig.GlobalConfig) (userModels.UserRegistrationInfo, bool, error) {
	discord.WebhookURL = "https://discord.com/api/webhooks/824381163367170058/OXSX51RHd9DyLFbFipjdW3yXmyYC8SWwqd6HiXl6UtDzu75RxS1LzWA800hWereJJumw"
	if len(os.Getenv("REGISTRATION_ERROR_WEBHOOK")) > 50 {
		discord.WebhookURL = os.Getenv("REGISTRATION_ERROR_WEBHOOK")
	}
	// if banned, errBanned := users.PublicKeyIsBanned(userInfo.PublicKey, gc.DB); banned {
	// 	return userInfo, false, errBanned
	// }
	if _, errExists := users.PublicKeyAlreadyExists(userInfo.PublicKey, gc.DB); errExists != nil {
		return userInfo, false, errExists
	}
	if _, errExists := users.PrimarySignerAlreadyExists(userInfo.PublicKey, gc.DB); errExists != nil {
		return userInfo, false, errExists
	}
	if len(userInfo.Mobile) > 0 {
		// geoData, _ := userModels.GetGeoInfo(userInfo.PublicIP)
		num, err := phonenumbers.Parse(userInfo.Mobile, userInfo.MobileCountryCode)
		if err == nil {
			mobile := fmt.Sprintf("+%v-%v", *num.CountryCode, *num.NationalNumber)
			userInfo.Mobile = mobile
		}
	}
	NormalizeUserRegistrationInfo(&userInfo)

	errValidation := ValidateUserRegistrationInfo(userInfo)

	if errValidation != nil {
		log.Printf("[RegisterUser]  Validation failed for user:%v, Error:%v[%v]", userInfo.Username, errValidation, errValidation.Error())
		discord.Say(fmt.Sprintf("[RegisterUser]  Validation failed for user:%v, Error:%v[%v]", userInfo.Username, errValidation, errValidation.Error()))
		return userInfo, false, errValidation
	}

	user, dbErrors := users.UserRegistrationDbChecks(userInfo, gc.DB)

	if dbErrors != nil {
		discord.Say(fmt.Sprintf("[RegisterUser] Registration DB check failed for user:%v, Error:%v[%v]", userInfo.Username, dbErrors, dbErrors.Error()))
		return userInfo, false, dbErrors
	}

	{

		//do mail validation
		validationResult, blockEmail, _ := user.VerifyEmailOnMailgun()
		if blockEmail {
			discord.Say(fmt.Sprintf("[RegisterUser] Mailgun Mail Validation failed for user:%v, Result:%+v", userInfo.Username, validationResult))
			return userInfo, false, &tErrors.ErrorEmailFailedValidation{Email: user.Email, Detail: fmt.Sprintf("Email validation failed for [%v]. Please put valid email and try again later.", userInfo.Email)}
		}

		emailSent, expectedCode, checkAndSendErr := CheckAndSendVerificationCode(userInfo)

		if checkAndSendErr != nil {
			log.Printf("error processing verification %s\n", checkAndSendErr)
			discord.Say(fmt.Sprintf("[RegisterUser] verification code process failed for user:%v, Error:%v, EnteredCode: %v, ExpectedCode:%v", userInfo.Username, checkAndSendErr, userInfo.VerificationCode, expectedCode))
			return userInfo, emailSent, checkAndSendErr
		}

		if emailSent {
			//no error, but email still sent
			return userInfo, emailSent, checkAndSendErr
		}

	}

	if os.Getenv("REGISTRATION_THROTTLE_PER_IP") != "" && os.Getenv("REGISTRATION_THROTTLE_PER_IP") != "0" {
		//check if same IP dat registered a user is up to 1hr
		var createdAt time.Time
		qry := "select created_at from users where public_ip = ? order by created_at DESC limit 1"
		e := gc.DB.Raw(qry, user.PublicIP).Scan(&createdAt).Error
		if e == nil {
			//record was found...check the time
			t := decimal.RequireFromString(os.Getenv("REGISTRATION_THROTTLE_PER_IP")).IntPart()
			if time.Since(createdAt) < (time.Duration(t) * time.Second) {
				//it is less than 1hr since registration from same IP, reject registration
				discord.Say(fmt.Sprintf("Too many registration from the IP %+v\n", user))
				return userInfo, false, &tErrors.CustomError{
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
	if data, e := dl.GenerateReferralLink(user.Username, gc); e == nil {
		user.ReferralLink = &data.DynamicLink
		user.ReferralQrCode = &data.QRCode
	}
	// build user wallet
	user.BuildPrimaryWallet()
	//save the user

	errCreate := gc.DB.Omit(clause.Associations).Create(user).Error
	if errCreate != nil {

		discord.Say(fmt.Sprintf("[RegisterUser] user creation failed for user:%v, with DB Error:%v\n\n\nFailedData:%+v", userInfo.Username, errCreate, userInfo))
		return userInfo, false, errors.New("unable to create user due to error in information")
	}

	{
		//send to monitoring service
		trackPublicKey := userModels.TrackedPublicKey{
			PublicKey: user.PublicKey,
		}
		errTrack := gc.RoachDB.Create(&trackPublicKey).Error
		if errTrack != nil {
			//if tracking of public key fails, then payment history generation service will pick it up and do justice to it
			discord.Say(fmt.Sprintf("[RegisterUser] tracking public key for payment history failed for user:%v, with DB Error:%v\n\n\nFailedData:%+v", userInfo.Username, errTrack, userInfo))

		}
	}
	return userInfo, false, nil

}
