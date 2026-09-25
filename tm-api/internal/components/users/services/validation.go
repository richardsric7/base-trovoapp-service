package users

import (
	errors "admin-panel-dashboard/internal/errors"
	models "admin-panel-dashboard/internal/models"
	"fmt"
	"log"
	"strings"

	"github.com/nyaruka/phonenumbers"
)

// ValidateUserUpdateInfo validates user update registration info
func ValidateUserUpdateInfo(user models.UserUpdateInfo, geoData models.IPAPI) error {

	// required parameters
	{

		// if len(user.Mobile) == 0 {
		// 	var x errors.ErrorMissingParameter
		// 	x.Parameter = "mobile"
		// 	return &x
		// }
		// if len(user.LastName) == 0 {
		// 	var x errors.ErrorMissingParameter
		// 	x.Parameter = "lastname"
		// 	return &x
		// }
		// if len(user.FirstName) == 0 {
		// 	var x errors.ErrorMissingParameter
		// 	x.Parameter = "firstname"
		// 	return &x
		// }
		// if len(user.LastName) > 50 {

		// 	return &errors.ErrorInvalidName{Field: "lastname"}
		// }
		// if len(user.FirstName) > 50 {

		// 	return &errors.ErrorInvalidName{Field: "firstname"}
		// }
		// if len(user.MiddleName) > 50 {

		// 	return &errors.ErrorInvalidName{Field: "middlename"}
		// }

		// //check if first name contains numbers
		// for _, c := range []byte(strings.ToLower(user.FirstName)) {
		// 	if strings.Contains("1234567890_", string(c)) {
		// 		log.Println("[ValidateUserUpdateInfo] first name validation failed for ", user)
		// 		return &errors.ErrorNameFailedValidation{Detail: fmt.Sprintf("%v not allowed in firstname", string(c))}
		// 	}
		// }
		// //check if last name contains numbers
		// for _, c := range []byte(strings.ToLower(user.LastName)) {
		// 	if strings.Contains("1234567890_", string(c)) {
		// 		log.Println("[ValidateUserUpdateInfo] last name validation failed for ", user)

		// 		return &errors.ErrorNameFailedValidation{Detail: fmt.Sprintf("%v not allowed in lastname", string(c))}
		// 	}
		// }

		// if !strings.Contains("FM", user.Gender) || len(user.Gender) != 1 {

		// 	return &errors.ErrorInvalidGender{}
		// }
		// check if mobile contains unaccepted character
		// acceptedPhoneChars := "+-1234567890"
		// for _, c := range []byte(strings.ToLower(user.Mobile)) {
		// 	if !strings.Contains(acceptedPhoneChars, string(c)) {
		// 		log.Println("[ValidateUserRegistrationInfo] mobile phone validation failed for ", user)

		// 		return &errors.ErrorPhoneNumberValidationFailed{Detail: fmt.Sprintf("phone number [%v] contains invalid character [%v]", user.Mobile, string(c))}
		// 	}
		// }
		// num, err := phonenumbers.Parse(user.Mobile, geoData.CountryCode)
		// if err != nil {
		// 	return &errors.ErrorPhoneNumberValidationFailed{Detail: fmt.Sprintf("phone number [%v] does not match your location", user.Mobile)}
		// }
		// if !phonenumbers.IsValidNumber(num) || phonenumbers.GetNumberType(num) == 0 {
		// 	//invalid or not mobile
		// 	return &errors.ErrorInvalidPhoneNumber{}
		// }

		// thumbail
		// if user.ImageThumbnail != "" {
		// 	dataURL, thumbnailerr := dataurl.DecodeString(user.ImageThumbnail)
		// 	if thumbnailerr != nil {

		// 		return &errors.ErrorInvalidImageThumbnail{}
		// 	}
		// 	if len(dataURL.Data) > 60600 {
		// 		return &errors.ErrorInvalidImageThumbnailSize{}
		// 	}

		// }

	}

	return nil

}

// NormalizeUserRegistrationInfo normalizes user registration info. lowercases usernames and emails., and trims space.
func NormalizeUserRegistrationInfo(user *models.UserRegistrationInfo) {

	user.Username = strings.TrimSpace(strings.ToLower(user.Username))
	user.FirstName = strings.TrimSpace(strings.ToUpper(user.FirstName))
	user.LastName = strings.TrimSpace(strings.ToUpper(user.LastName))
	user.MiddleName = strings.TrimSpace(strings.ToUpper(user.MiddleName))
	user.Email = strings.ReplaceAll(strings.TrimSpace(strings.ToLower(user.Email)), " ", "")
	user.Referrer = strings.TrimSpace(strings.ToLower(user.Referrer))
	user.Mobile = strings.ReplaceAll(strings.TrimSpace(user.Mobile), " ", "")
	if len(user.Mobile) > 0 {
		geoData, _ := models.GetGeoInfo(user.PublicIP)
		num, err := phonenumbers.Parse(user.Mobile, geoData.CountryCode)
		if err == nil {
			mobile := fmt.Sprintf("+%v-%v", *num.CountryCode, *num.NationalNumber)
			user.Mobile = mobile
		}
	}

}

// ValidateUserRegistrationInfo validates user registration info
func ValidateUserRegistrationInfo(user models.UserRegistrationInfo) error {
	// required parameters
	{
		if len(user.Username) == 0 {
			return &errors.ErrorMissingParameter{Parameter: "username"}

		}
		if len(user.Username) > 16 {
			return &errors.ErrorInvalidUsernameFormat{Username: "username", Detail: "username cannot be more than 16 characters"}

		}
		onlyNumbers := true
		charCount := 0
		acceptedChars := "abcdefghijklmnopqrstuvwxyz_1234567890"
		acceptedPhoneChars := "+-1234567890"
		for _, c := range []byte(strings.ToLower(user.Username)) {

			if !strings.Contains(acceptedChars, string(c)) {
				r := string(c)
				if string(c) == " " {
					r = "whitespace"
				}
				return &errors.ErrorInvalidUsernameFormat{Username: "username", Detail: fmt.Sprintf("'%s' is not allowed in usernames", r)}

			}

			if strings.Contains("abcdefghijklmnopqrstuvwxyz", string(c)) {
				if onlyNumbers {
					onlyNumbers = false
				}
				charCount++
			}

		}
		if onlyNumbers {
			return &errors.ErrorInvalidUsernameFormat{Username: user.Username, Detail: "username with only numbers are not allowed"}

		}
		if charCount < 3 {
			return &errors.ErrorInvalidUsernameFormat{Username: user.Username, Detail: "username must contain atleast 3 English alphabets"}

		}

		if len(user.PublicKey) == 0 {
			var x errors.ErrorMissingParameter
			x.Parameter = "publicKey"
			return &x
		}

		if len(user.Email) == 0 {
			var x errors.ErrorMissingParameter
			x.Parameter = "email"
			return &x
		}
		if len(user.Mobile) == 0 {
			var x errors.ErrorMissingParameter
			x.Parameter = "mobile"
			return &x
		}
		if len(user.LastName) == 0 {
			var x errors.ErrorMissingParameter
			x.Parameter = "lastname"
			return &x
		}
		if len(user.FirstName) == 0 {
			var x errors.ErrorMissingParameter
			x.Parameter = "firstname"
			return &x
		}
		if len(user.LastName) > 50 {

			return &errors.ErrorInvalidName{Field: "lastname"}
		}
		if len(user.FirstName) > 50 {

			return &errors.ErrorInvalidName{Field: "firstname"}
		}
		if len(user.MiddleName) > 50 {

			return &errors.ErrorInvalidName{Field: "middlename"}
		}

		// check if first name contains numbers
		for _, c := range []byte(strings.ToLower(user.FirstName)) {
			if strings.Contains("1234567890_", string(c)) {
				log.Println("[ValidateUserRegistrationInfo] first name validation failed for ", user)
				return &errors.ErrorNameFailedValidation{Detail: fmt.Sprintf("%v not allowed in firstname", string(c))}
			}
		}
		// check if last name contains numbers
		for _, c := range []byte(strings.ToLower(user.LastName)) {
			if strings.Contains("1234567890_", string(c)) {
				log.Println("[ValidateUserRegistrationInfo] last name validation failed for ", user)

				return &errors.ErrorNameFailedValidation{Detail: fmt.Sprintf("%v not allowed in lastname", string(c))}
			}
		}
		// if len(user.Gender) != 0 {

		// 	if !strings.Contains("FM", user.Gender) || len(user.Gender) != 1 {

		// 		return &errors.ErrorInvalidGender{}
		// 	}
		// }
		if !strings.Contains(user.Email, "@") {
			return &errors.ErrorEmailFailedValidation{Email: user.Email, Detail: fmt.Sprintf("%v is not an email", user.Email)}

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
				return &errors.ErrorInvalidParameter{Param: "email", ErrMessage: fmt.Sprintf("'%v' is not allowed in an email", string(c))}

			}

			if strings.Contains("abcdefghijklmnopqrstuvwxyz", string(c)) {
				if onlyNumbers {
					onlyNumbers = false
				}
				charCount++
			}

		}
		if onlyNumbers {
			return &errors.ErrorInvalidParameter{Param: "email", ErrMessage: "email cannot contain only numbers"}

		}
		if charCount < 3 {
			return &errors.ErrorInvalidParameter{Param: "email", ErrMessage: "email must have at least 3 english characters"}

		}
		// check if mobile contains unaccepted character
		for _, c := range []byte(strings.ToLower(user.Mobile)) {
			if !strings.Contains(acceptedPhoneChars, string(c)) {
				log.Println("[ValidateUserRegistrationInfo] mobile phone validation failed for user ", user.Username)

				return &errors.ErrorPhoneNumberValidationFailed{Detail: fmt.Sprintf("mobile number [%v] contains invalid character [%v]", user.Mobile, string(c))}
			}
		}
		num, err := phonenumbers.Parse(user.Mobile, "")
		if err != nil {
			return &errors.ErrorInvalidPhoneNumber{}
		}
		if !phonenumbers.IsValidNumber(num) || phonenumbers.GetNumberType(num) == 0 {
			// invalid or not mobile
			return &errors.ErrorInvalidPhoneNumber{}
		}

	}

	return nil

}
