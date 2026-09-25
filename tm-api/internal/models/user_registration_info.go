package models

import (
	bantuErrors "admin-panel-dashboard/internal/errors"
	"context"
	"errors"
	"log"
	"os"
	"strings"
	"time"

	"github.com/mailgun/mailgun-go/v4"
)

// UserRegistrationInfo model for user resitration info
type UserRegistrationInfo struct {
	ID               string
	Username         string
	Address          string
	Email            string
	PublicIP         string
	VerificationCode string
	LastName         string
	FirstName        string
	MiddleName       string
	Mobile           string
	CaptchaCode      string
	CaptchaChallenge string
	Referrer         string
}

// Gender          string
// Telegram         string
// Instagram        string
// Twitter          string
// ImageThumbnail   string

func (u *UserRegistrationInfo) VerifyEmailOnMailgun() (validationResult mailgun.EmailVerification, blockEmail bool, err error) {
	// To use the /v4 version of validations define MG_URL in the environment
	// as `https://api.mailgun.net/v4` or set `v.SetAPIBase("https://api.mailgun.net/v4")`
	if os.Getenv("ENABLE_EMAIL_VALIDATION") == "1" {
		// check if email validation passes
		if len(os.Getenv("MAILGUN_VALIDATOR_API_KEY")) == 0 {
			log.Println("MAILGUN_VALIDATOR_API_KEY not set!")
			err = &bantuErrors.ErrorTemporaryServerError{}
			return
		}
		apiKey := os.Getenv("MAILGUN_VALIDATOR_API_KEY")

		// Create an instance of the Validator
		v := mailgun.NewEmailValidator(apiKey)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		validationResult, err = v.ValidateEmail(ctx, u.Email, true)

		if err != nil {
			log.Println("[VerifyEmailOnMailgun] mail validation request error:", err)
			validationResult, err = v.ValidateEmail(ctx, u.Email, true)
			if err != nil {
				log.Println("[VerifyEmailOnMailgun] mail validation request failed second time with error:", err)

				return validationResult, blockEmail, &bantuErrors.ErrorTemporaryServerError{}
			}
		}
		if os.Getenv("SHOW_MAIL_VALIDATION_RESULT") == "1" {
			log.Printf("[VerifyEmailOnMailgun] validation result: %#v\n", validationResult)
		}
		if validationResult.IsDisposableAddress || strings.Contains(validationResult.Reason, "unknown") || strings.Contains(validationResult.Reason, "risk") || strings.Contains(validationResult.Reason, "mail") || strings.Contains(validationResult.Reason, "domain") || strings.Contains(validationResult.Reason, "catch") || strings.Contains(validationResult.Reason, "long") || strings.Contains(validationResult.Reason, "no_mx") {
			log.Printf("[VerifyEmailOnMailgun] %v, user:%#v %#v\n", errors.New("failed mail validation"), *u, validationResult)
			blockEmail = true
			return validationResult, blockEmail, nil
		}

	}
	return
}
