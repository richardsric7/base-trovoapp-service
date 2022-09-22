package users

import (
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	users "trovo-wallet-api/internal/components/users/models"
	tErrors "trovo-wallet-api/internal/errors"
	tMail "trovo-wallet-api/internal/mail"
	"trovo-wallet-api/internal/sms"

	"time"
	cache "trovo-wallet-api/internal/cache"

	"github.com/gin-gonic/gin"
	"github.com/nyaruka/phonenumbers"
	"gorm.io/gorm"
)

const saltInCode = "7HrtrGCRUPp8j5Gze"

// NumberOfCharactersInVerificationCode num of chars for verification code
const NumberOfCharactersInVerificationCode = 6

// GenerateAccountRecoveryEmailOTP generates email verification code
func GenerateAccountRecoveryEmailOTP(userInfo *users.User) string {
	salt := os.Getenv("VERIFICATION_CODE_SALT")
	dateReference := time.Now().String()

	tohash := sha256.New()
	tohash.Write([]byte(userInfo.Email))
	tohash.Write([]byte(userInfo.PublicKey))
	tohash.Write([]byte(userInfo.Username))
	tohash.Write([]byte(dateReference))
	tohash.Write([]byte(saltInCode))
	tohash.Write([]byte(salt))
	hashed := tohash.Sum(nil)
	data := binary.BigEndian.Uint32(hashed)

	s := strconv.FormatUint(uint64(data), 10)

	return s[0:NumberOfCharactersInVerificationCode]
}

// GenerateEmailVerificationCode generates email verification code
func GenerateEmailVerificationCode(userInfo users.UserRegistrationInfo, salt string) string {

	dateReference := time.Now().Format("2006-01-02")

	tohash := sha256.New()
	tohash.Write([]byte(userInfo.Email))
	tohash.Write([]byte(userInfo.PublicKey))
	tohash.Write([]byte(userInfo.Username))
	tohash.Write([]byte(dateReference))
	tohash.Write([]byte(saltInCode))
	tohash.Write([]byte(salt))
	hashed := tohash.Sum(nil)
	data := binary.BigEndian.Uint32(hashed)

	s := strconv.FormatUint(uint64(data), 10)

	return s[0:NumberOfCharactersInVerificationCode]
}

// GeneratePhoneVerificationCode generates email verification code
func GeneratePhoneVerificationCode(userInfo *users.User, salt string) string {

	dateReference := time.Now().Format("2006-01-02")

	tohash := sha256.New()
	tohash.Write([]byte(*userInfo.Mobile))
	tohash.Write([]byte(userInfo.PublicKey))
	tohash.Write([]byte(userInfo.Username))
	tohash.Write([]byte(dateReference))
	tohash.Write([]byte(saltInCode))
	tohash.Write([]byte(salt))
	hashed := tohash.Sum(nil)
	data := binary.BigEndian.Uint32(hashed)

	s := strconv.FormatUint(uint64(data), 10)

	return s[0:NumberOfCharactersInVerificationCode]
}

// CheckAndSendVerificationCode checks and sends verification code
func CheckAndSendVerificationCode(userInfo users.UserRegistrationInfo) (bool, string, error) {
	//Verification code checks
	salt := os.Getenv("VERIFICATION_CODE_SALT")

	if len(salt) == 0 {
		return false, "", &tErrors.ErrorMissingEnvironmentVariable{EnvironmentVariable: "VERIFICATION_CODE_SALT"}
	}

	expectedVerificationCode := GenerateEmailVerificationCode(userInfo, salt)
	//check verification code
	if len(userInfo.VerificationCode) == 0 {
		//verification code not yet included

		fmt.Printf("Verification code for %s is %s\n", userInfo.Email, expectedVerificationCode)

		_, _, errSendVerificationCode := tMail.SendEmailVerificationCode(userInfo.Email, expectedVerificationCode)

		if errSendVerificationCode != nil {
			return false, expectedVerificationCode, errSendVerificationCode
		}

		//email sent, so return true as second argument
		return true, expectedVerificationCode, nil

	}

	//verification code included.
	if expectedVerificationCode != userInfo.VerificationCode {
		return false, expectedVerificationCode, &tErrors.ErrorInvalidVerificationCode{}
	}

	//email not sent, and no error
	return false, expectedVerificationCode, nil

}

// UpdatePhoneNumber updates phone number
func UpdatePhoneNumber(userInfo *users.User, mobileCountryCode, mobile string, db *gorm.DB, redisCache *cache.RedisCache) error {

	if len(mobile) == 0 {

		return &tErrors.CustomError{
			Param:      "mobile",
			Err:        "no mobile number provided",
			ErrMessage: "no mobile number provided",
		}

	}
	var parsedMobile string
	num, err := phonenumbers.Parse(mobile, mobileCountryCode)

	if err != nil {
		log.Printf("unable to parse mobile number for user %v due to: %v\n", userInfo.Username, err)
		return &tErrors.ErrorTemporaryServerError{}
	}

	if !phonenumbers.IsValidNumber(num) {
		log.Printf("unable to determine mobile number for validity for user %v\n", userInfo.Username)
		return &tErrors.ErrorTemporaryServerError{}
	}

	parsedMobile = fmt.Sprintf("+%v-%v", *num.CountryCode, *num.NationalNumber)
	if userInfo.Mobile != nil {
		if *userInfo.Mobile == parsedMobile {

			return &tErrors.CustomError{
				Param:      "mobile",
				Err:        "mobile number provided is same with one on record",
				ErrMessage: "mobile number provided is same with one on record",
			}

		}
	}

	dbErr := db.First(&users.User{}, "mobile = ?", parsedMobile).Error
	if dbErr != nil {
		if dbErr != gorm.ErrRecordNotFound {
			log.Println("unable to query database due to:", dbErr)
			return &tErrors.ErrorTemporaryServerError{}
		}

	}

	if dbErr == nil {
		return &tErrors.CustomError{
			Param:      "mobile",
			Err:        "mobile number provided is already in use",
			ErrMessage: "mobile number provided is already in use",
		}
	}
	if userInfo.MobileVerified == 1 && userInfo.LastUpdatedMobileOn.Add(time.Hour*24*90).After(time.Now()) {
		//phone modified less than
		return &tErrors.CustomError{
			Param:      "mobile",
			Err:        "mobile phone cannot be modified",
			ErrMessage: "mobile phone cannot be modified until after 90 days of last modification",
		}

	}
	userInfo.Mobile = &parsedMobile
	userInfo.LastUpdatedMobileOn = time.Now()
	{
		geoData, _ := users.GetGeoInfo(userInfo.PublicIP)

		userInfo.CountryCode = &geoData.CountryCode
		userInfo.City = &geoData.City
		userInfo.PublicIP = geoData.IP
		userInfo.Region = &geoData.Region
		userInfo.RegionName = &geoData.RegionName
	}

	saveError := db.Save(userInfo).Error
	if saveError != nil {
		log.Printf("unable to save new mobile number for user %v due to: %v\n", userInfo.Username, saveError)
		return &tErrors.ErrorTemporaryServerError{}
	}

	errCode := SendPhoneVerificationCode(userInfo, db, redisCache)
	if errCode != nil {
		return &tErrors.CustomError{
			Param:      "mobile",
			Err:        "mobile number saved but unable to send verification code",
			ErrMessage: "mobile number has been saved but verification code was not able to send at this time. Please try again later through the verify phone section.",
		}
	}

	return nil

}

// CheckhoneVerificationCode checks and sends verification code
func CheckPhoneVerificationCode(userInfo *users.User, verificationCode string, db *gorm.DB, c *gin.Context) error {
	if userInfo.MobileVerified == 1 {
		//phone already verified. exit with error
		return &tErrors.CustomError{
			Param:      "mobile",
			Err:        "mobile phone already verified",
			ErrMessage: "mobile phone already verified",
		}

	}
	if len(verificationCode) == 0 {
		//phone already verified. exit with error
		return &tErrors.CustomError{
			Param:      "mobile",
			Err:        "no verification code provided",
			ErrMessage: "no verification code provided",
		}

	}
	//Verification code checks
	salt := os.Getenv("VERIFICATION_CODE_SALT")

	if len(salt) == 0 {
		return &tErrors.ErrorMissingEnvironmentVariable{EnvironmentVariable: "VERIFICATION_CODE_SALT"}
	}
	var userPhoneVerification users.UserMobilePhoneVerification
	//get the verificationRecord
	errDB := db.First(&userPhoneVerification, "user_id = ?", userInfo.ID).Error

	if errDB != nil {
		//check error if server error
		if !errors.Is(errDB, gorm.ErrRecordNotFound) {
			log.Printf("[CheckPhoneVerificationCode] Failed to fetch verification for user %v at this time due to Error: %s\n", userInfo.Username, errDB.Error())

			return &tErrors.ErrorTemporaryServerError{}
		}
		//no record found, check conditions
		return &tErrors.CustomError{
			Param:      "mobile",
			Err:        "no verification code requested",
			ErrMessage: "Your have not requested for verification code before. Please press 'send code' to request for verification code",
		}

	}
	//record is found, check if it matches the expected verification code
	if userPhoneVerification.VerificationCode != verificationCode {
		return &tErrors.ErrorInvalidVerificationCode{}
	}

	userInfo.MobileVerified = 1
	errSaveUser := db.Save(userInfo).Error
	if errSaveUser != nil {
		log.Printf("[CheckPhoneVerificationCode]Unable to save mobile verification for user %v at this time due to Error: %s\n", userInfo.Username, errSaveUser.Error())

		return &tErrors.CustomError{
			Param:      "mobileVerified",
			Err:        "error failed to save mobile verified",
			ErrMessage: "Unable to verify mobile phone at this time. Please try again later.",
		}

	}
	//verification code matches

	return nil

}

// CheckAccountRecoveryEmailOTP checks and sends verification code
func CheckAccountRecoveryEmailOTP(userInfo *users.User, verificationCode string, db *gorm.DB, c *gin.Context) error {

	if len(verificationCode) == 0 {
		//phone already verified. exit with error
		return &tErrors.CustomError{
			Param:      "email",
			Err:        "no verification code provided",
			ErrMessage: "no verification code provided",
		}

	}

	var userVerification users.UserAccountRecoveryEmailVerification
	//get the verificationRecord
	errDB := db.First(&userVerification, "user_id = ?", userInfo.ID).Error

	if errDB != nil {
		//check error if server error
		if !errors.Is(errDB, gorm.ErrRecordNotFound) {
			log.Printf("[CheckAccountRecoveryEmailOTP] Failed to fetch verification for user %v at this time due to Error: %s\n", userInfo.Username, errDB.Error())

			return &tErrors.ErrorTemporaryServerError{}
		}
		//no record found, check conditions
		return &tErrors.CustomError{
			Param:      "mobile",
			Err:        "no verification code requested",
			ErrMessage: "Your have not requested for verification code before. Please press 'send code' to request for verification code",
		}

	}
	//record is found, check if it matches the expected verification code
	if userVerification.VerificationCode != verificationCode {
		return &tErrors.ErrorInvalidVerificationCode{}
	}

	//verification code matches

	return nil

}

// SendPhoneVerificationCode checks and sends verification code
func SendPhoneVerificationCode(userInfo *users.User, db *gorm.DB, redisCache *cache.RedisCache) error {
	if userInfo.MobileVerified == 1 {
		//phone already verified. exit with error
		return &tErrors.CustomError{
			Param:      "mobile",
			Err:        "mobile phone already verified",
			ErrMessage: "mobile phone already verified",
		}

	}
	if userInfo.Mobile == nil {

		return &tErrors.CustomError{
			Param:      "mobile",
			Err:        "no mobile number attached to account",
			ErrMessage: "No mobile number attached to account. Please contact support.",
		}

	}

	if len(*userInfo.Mobile) == 0 {

		return &tErrors.CustomError{
			Param:      "mobile",
			Err:        "invalid mobile number attached to account",
			ErrMessage: "Invalid mobile number attached to account. Please contact support.",
		}

	}
	//Verification code checks
	salt := os.Getenv("VERIFICATION_CODE_SALT")

	if len(salt) == 0 {
		return &tErrors.ErrorMissingEnvironmentVariable{EnvironmentVariable: "VERIFICATION_CODE_SALT"}
	}
	var userPhoneVerification users.UserMobilePhoneVerification
	//get the verification Record
	errDB := db.First(&userPhoneVerification, "user_id = ?", userInfo.ID).Error

	tx := db.Begin()
	defer tx.Rollback()

	if errDB != nil {
		//check error if server error
		if !errors.Is(errDB, gorm.ErrRecordNotFound) {
			return &tErrors.ErrorTemporaryServerError{}
		}
		verificationCode := GeneratePhoneVerificationCode(userInfo, salt)
		//record not found...create new record
		userPhoneVerification = users.UserMobilePhoneVerification{
			UserID:           userInfo.ID,
			Mobile:           *userInfo.Mobile,
			VerificationCode: verificationCode,
		}

		errDB = tx.Create(&userPhoneVerification).Error
		if errDB != nil {
			//could not create verification code
			log.Printf("[SendPhoneVerificationCode] Error creating verification code for user %s. Error: %s\n", userInfo.Username, errDB.Error())
			return &tErrors.ErrorTemporaryServerError{}
		}
		//send SMS
		message := fmt.Sprintf("Your BantuPay mobile phone confirmation code is %s. One time use only.", verificationCode)
		errSMS := sms.SendSMS(*userInfo.Mobile, message, db)
		if errSMS != nil {
			log.Printf("[SendPhoneVerificationCode] Error sending verification code for user %s. Error: %s\n", userInfo.Username, errSMS.Error())

			return &tErrors.CustomError{
				Param:      "mobile",
				Err:        "error-cannot-send-sms-to-mobile",
				ErrMessage: "Unable to send OTP to your registered mobile phone at this time. Please ensure your mobile number is correct and then try again later.",
			}
		}
		tx.Commit()
		return nil

	}
	//record is found
	//check when last request was made
	log.Printf("[SendPhoneVerificationCode] username: %s, requested date: %s, reference Date: %s\n", userInfo.Username, userPhoneVerification.RequestDate.Format("2006-01-02"), time.Now().Format("2006-01-02"))

	if userPhoneVerification.RequestDate.Format("2006-01-02") == time.Now().Format("2006-01-02") {
		return &tErrors.CustomError{
			Param:      "mobile",
			Err:        "daily request quota exceeded",
			ErrMessage: "You have already exhausted your request quota for the day. Wait till you recieve the code or you wait for another day",
		}
	}

	verificationCode := GeneratePhoneVerificationCode(userInfo, salt)

	//update the record
	userPhoneVerification.RequestDate = time.Now()
	userPhoneVerification.VerificationCode = verificationCode
	errDB = tx.Save(&userPhoneVerification).Error
	if errDB != nil {
		//could not update verification code
		log.Printf("[SendPhoneVerificationCode] Error updating verification code for user %s. Error: %s\n", userInfo.Username, errDB.Error())
		return &tErrors.ErrorTemporaryServerError{}
	}
	//updated successfully
	//send SMS
	message := fmt.Sprintf("Your BantuPay mobile phone confirmation code is %s. One time use only.", verificationCode)
	errSMS := sms.SendSMS(*userInfo.Mobile, message, db)
	if errSMS != nil {
		log.Printf("[SendPhoneVerificationCode] Error sending verification code for user %s. Error: %s\n", userInfo.Username, errSMS.Error())

		return &tErrors.CustomError{
			Param:      "mobile",
			Err:        "error-cannot-send-sms-to-mobile",
			ErrMessage: "Unable to send OTP to your registered mobile phone at this time. Please ensure your mobile number is correct and then try again later.",
		}
	}
	tx.Commit()
	return nil
}

// SendAccountRecoveryEmailOTP checks and sends verification code
func SendAccountRecoveryEmailOTP(userInfo *users.User, db *gorm.DB) error {

	if userInfo.Email == "" {

		return &tErrors.CustomError{
			Param:      "email",
			Err:        "no email attached to account",
			ErrMessage: "No email attached to account. Please contact support.",
		}

	}

	var userVerification users.UserAccountRecoveryEmailVerification
	//get the verification Record
	errDB := db.First(&userVerification, "user_id = ?", userInfo.ID).Error

	tx := db.Begin()
	defer tx.Rollback()

	if errDB != nil {
		//check error if server error
		if !errors.Is(errDB, gorm.ErrRecordNotFound) {
			return &tErrors.ErrorTemporaryServerError{}
		}
		verificationCode := GenerateAccountRecoveryEmailOTP(userInfo)
		//record not found...create new record
		userVerification = users.UserAccountRecoveryEmailVerification{
			UserID:           userInfo.ID,
			Email:            userInfo.Email,
			VerificationCode: verificationCode,
		}

		errDB = tx.Create(&userVerification).Error
		if errDB != nil {
			//could not create verification code
			log.Printf("[SendAccountRecoveryEmailOTP] Error creating verification code for user %s. Error: %s\n", userInfo.Username, errDB.Error())
			return &tErrors.ErrorTemporaryServerError{}
		}
		//send TOP

		_, _, errSendVerificationCode := tMail.SendEmailVerificationCode(userInfo.Email, verificationCode)

		if errSendVerificationCode != nil {
			return errSendVerificationCode
		}

		tx.Commit()
		return nil

	}
	//record is found
	//check when last request was made
	log.Printf("[SendAccountRecoveryEmailOTP] username: %s, requested date: %s, reference Date: %s\n", userInfo.Username, userVerification.RequestDate.Format("2006-01-02"), time.Now().Format("2006-01-02"))

	if userVerification.RequestDate.Format("2006-01-02") == time.Now().Format("2006-01-02") {
		return &tErrors.CustomError{
			Param:      "mobile",
			Err:        "daily request quota exceeded",
			ErrMessage: "You have already exhausted your request quota for the day. Wait till you recieve the code or you wait for another day",
		}
	}

	verificationCode := GenerateAccountRecoveryEmailOTP(userInfo)

	//update the record
	userVerification.RequestDate = time.Now()
	userVerification.VerificationCode = verificationCode
	errDB = tx.Save(&userVerification).Error
	if errDB != nil {
		//could not update verification code
		log.Printf("[SendAccountRecoveryEmailOTP] Error updating verification code for user %s. Error: %s\n", userInfo.Username, errDB.Error())
		return &tErrors.ErrorTemporaryServerError{}
	}
	//updated successfully
	//send OTP

	_, _, errSendVerificationCode := tMail.SendEmailVerificationCode(userInfo.Email, verificationCode)

	if errSendVerificationCode != nil {
		return errSendVerificationCode
	}
	tx.Commit()
	return nil
}
