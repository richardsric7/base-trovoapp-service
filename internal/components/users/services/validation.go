package users

import (
	"fmt"
	"log"
	"strings"
	usermodels "trovo-wallet-api/internal/components/users/models"
	tErrors "trovo-wallet-api/internal/errors"

	"github.com/nyaruka/phonenumbers"
)

//NormalizeUserRegistrationInfo normalizes user registration info. lowercases usernames and emails., and trims space.
func NormalizeUserRegistrationInfo(user *usermodels.UserRegistrationInfo) {

	user.Username = strings.TrimSpace(strings.ToLower(user.Username))
	user.FirstName = strings.TrimSpace(strings.ToUpper(user.FirstName))
	user.LastName = strings.TrimSpace(strings.ToUpper(user.LastName))
	user.Email = strings.ReplaceAll(strings.TrimSpace(strings.ToLower(user.Email)), " ", "")
	user.Referrer = strings.TrimSpace(strings.ToLower(user.Referrer))
	user.Mobile = strings.ReplaceAll(strings.TrimSpace(user.Mobile), " ", "")
	if len(user.Mobile) > 0 {
		// geoData, _ := usermodels.GetGeoInfo(user.PublicIP)
		num, err := phonenumbers.Parse(user.Mobile, user.MobileCountryCode)
		if err == nil {
			mobile := fmt.Sprintf("+%v-%v", *num.CountryCode, *num.NationalNumber)
			user.Mobile = mobile
		}
	}

}

//ValidateUserRegistrationInfo validates user registration info
func ValidateUserRegistrationInfo(user usermodels.UserRegistrationInfo) error {
	//required parameters
	{
		if len(user.Username) == 0 {
			return &tErrors.ErrorMissingParameter{Parameter: "username"}

		}
		if len(user.Username) > 16 {
			return &tErrors.ErrorInvalidUsernameFormat{Username: "username", Detail: "username cannot be more than 16 characters"}

		}
		onlyNumbers := true
		charCount := 0
		acceptedChars := "abcdefghijklmnopqrstuvwxyz1234567890"
		acceptedPhoneChars := "+-1234567890"
		for _, c := range []byte(strings.ToLower(user.Username)) {

			if !strings.Contains(acceptedChars, string(c)) {
				r := string(c)
				if string(c) == " " {
					r = "whitespace"
				}
				return &tErrors.ErrorInvalidUsernameFormat{Username: "username", Detail: fmt.Sprintf("'%s' is not allowed in usernames", r)}

			}

			if strings.Contains("abcdefghijklmnopqrstuvwxyz", string(c)) {
				if onlyNumbers {
					onlyNumbers = false
				}
				charCount++
			}

		}
		if onlyNumbers {
			return &tErrors.ErrorInvalidUsernameFormat{Username: user.Username, Detail: "username with only numbers are not allowed"}

		}
		// if charCount < 3 {
		// 	return &tErrors.ErrorInvalidUsernameFormat{Username: user.Username, Detail: "username must contain atleast 3 English alphabets"}

		// }

		if len(user.PublicKey) == 0 {
			var x tErrors.ErrorMissingParameter
			x.Parameter = "publicKey"
			return &x
		}

		if len(user.Email) == 0 {
			var x tErrors.ErrorMissingParameter
			x.Parameter = "email"
			return &x
		}
		if len(user.Mobile) == 0 {
			var x tErrors.ErrorMissingParameter
			x.Parameter = "mobile"
			return &x
		}
		if len(user.MobileCountryCode) == 0 {
			var x tErrors.ErrorMissingParameter
			x.Parameter = "mobileCountryCode"
			return &x
		}
		if len(user.LastName) == 0 && user.Corporate == 0 {
			var x tErrors.ErrorMissingParameter
			x.Parameter = "lastname"
			return &x
		}
		if len(user.FirstName) == 0 {
			var x tErrors.ErrorMissingParameter
			x.Parameter = "firstname"
			return &x
		}
		if len(user.LastName) > 50 {

			return &tErrors.ErrorInvalidName{Field: "lastName"}
		}
		if len(user.LastName) > 50 {

			return &tErrors.ErrorInvalidName{Field: "lastName"}
		}
		if len(user.FirstName) > 50 {

			return &tErrors.ErrorInvalidName{Field: "firstName"}
		}
		if len(user.MobileCountryCode) > 2 {

			return &tErrors.CustomError{Param: "mobileCountryCode", Err: "error invalid mobileCountryCode", ErrMessage: "mobileCountryCode must be 2 characters, eg NG, US, CA.", Code: 400}
		}

		//check if first name contains numbers
		for _, c := range []byte(strings.ToLower(user.FirstName)) {
			if strings.Contains("1234567890_", string(c)) && user.Corporate == 0 {
				log.Println("[ValidateUserRegistrationInfo] first name validation failed for ", user)
				return &tErrors.ErrorNameFailedValidation{Detail: fmt.Sprintf("%v not allowed in firstname", string(c))}
			}
		}
		//check if last name contains numbers
		if len(user.LastName) > 0 {
			for _, c := range []byte(strings.ToLower(user.LastName)) {
				if strings.Contains("1234567890_", string(c)) && user.Corporate == 0 {
					log.Println("[ValidateUserRegistrationInfo] last name validation failed for ", user)

					return &tErrors.ErrorNameFailedValidation{Detail: fmt.Sprintf("%v not allowed in lastname", string(c))}
				}
			}
		}
		if !strings.Contains(user.Email, "@") {
			return &tErrors.ErrorEmailFailedValidation{Email: user.Email, Detail: fmt.Sprintf("%v is not an email", user.Email)}

		}
		onlyNumbers = true
		charCount = 0
		emailUser := strings.Split(user.Email, "@")[0]
		acceptedChars = "abcdefghijklmnopqrstuvwxyz_1234567890.-"

		for _, c := range []byte(strings.ToLower(emailUser)) {

			if !strings.Contains(acceptedChars, string(c)) {
				// r := string(c)
				// if string(c) == " " {
				// 	r = "whitespace"
				// }
				return &tErrors.ErrorInvalidEmailFormat{Email: user.Email, Detail: fmt.Sprintf("'%v' is not allowed in an email", string(c))}

			}

			if strings.Contains("abcdefghijklmnopqrstuvwxyz", string(c)) {
				if onlyNumbers {
					onlyNumbers = false
				}
				charCount++
			}

		}
		if onlyNumbers {
			return &tErrors.ErrorInvalidEmailFormat{Email: user.Email, Detail: "email cannot contain only numbers"}

		}
		if charCount < 3 {
			return &tErrors.ErrorInvalidEmailFormat{Email: user.Email, Detail: "email must have at least 3 english characters"}

		}
		//check if mobile contains unaccepted character
		for _, c := range []byte(strings.ToLower(user.Mobile)) {
			if !strings.Contains(acceptedPhoneChars, string(c)) {
				log.Println("[ValidateUserRegistrationInfo] mobile phone validation failed for user ", user.Username)

				return &tErrors.ErrorPhoneNumberValidationFailed{Detail: fmt.Sprintf("mobile number [%v] contains invalid character [%v]", user.Mobile, string(c))}
			}
		}
		num, err := phonenumbers.Parse(user.Mobile, "")
		if err != nil {
			return &tErrors.ErrorInvalidPhoneNumber{}
		}
		if !phonenumbers.IsValidNumber(num) || phonenumbers.GetNumberType(num) == 0 {
			//invalid or not mobile
			return &tErrors.ErrorInvalidPhoneNumber{}
		}

		// //thumbail
		// if user.ImageThumbnail != "" {
		// 	dataURL, thumbnailerr := dataurl.DecodeString(user.ImageThumbnail)
		// 	if thumbnailerr != nil {

		// 		return &tErrors.ErrorInvalidImageThumbnail{}
		// 	}
		// 	if len(dataURL.Data) > 25600 {
		// 		return &tErrors.ErrorInvalidImageThumbnailSize{}
		// 	}

		// }

		// //SM Handles
		// if len(user.Twitter) > 0 && len(user.Twitter) < 2 {

		// 	return &tErrors.ErrorInvalidTwitterHandleFormat{Handle: "twitter", Detail: "length too short"}
		// }
		// if len(user.Telegram) > 0 && len(user.Telegram) < 2 {

		// 	return &tErrors.ErrorInvalidTelegramHandleFormat{Handle: "telegram", Detail: "length too short"}
		// }
		// if len(user.Instagram) > 0 && len(user.Instagram) < 2 {

		// 	return &tErrors.ErrorInvalidInstagramHandleFormat{Handle: "instagram", Detail: "length too short"}
		// }

	}

	return nil

}
