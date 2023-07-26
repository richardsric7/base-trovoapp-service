package sms

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	tErrors "trovo-wallet-api/internal/errors"

	"github.com/infobip/infobip-api-go-client/v2"
	"gorm.io/gorm"
)

type SmsProvider struct {
	PhonePrefix string `gorm:"primaryKey" json:"phone_prefix"`
	Provider    string `json:"provider"`
}

func SendSMS(userMobile, messageBody string, db *gorm.DB) error {

	if len(userMobile) == 0 {
		return &tErrors.CustomError{
			Param:      "mobile",
			Err:        "invalid mobile number",
			ErrMessage: "Invalid Mobile Number",
		}
	}

	if len(messageBody) == 0 {
		return &tErrors.CustomError{
			Param:      "mobile",
			Err:        "invalid message body",
			ErrMessage: "Invalid Message Body",
		}
	}

	smsProvider := getSMSProvider(userMobile, db)
	destNumber := strings.ReplaceAll(strings.ReplaceAll(userMobile, "+", ""), "-", "")

	if strings.EqualFold(smsProvider, "infobip") {
		return SendSMSWithInfobip(destNumber, messageBody)
	} else if strings.EqualFold(smsProvider, "termii") {
		return SendSMSWithTermiiGateway(destNumber, messageBody)

	} else {
		return SendSMSWithTermiiGateway(destNumber, messageBody)
	}

}

func SendSMSWithInfobip(destNumber, messageBody string) error {

	baseHost := os.Getenv("INFOBIP_SMS_HOST")
	if len(baseHost) == 0 {
		return &tErrors.ErrorTemporaryServerError{}
	}
	smsAPIKey := os.Getenv("INFOBIP_SMS_API_KEY")
	if len(smsAPIKey) == 0 {
		return &tErrors.ErrorTemporaryServerError{}
	}

	configuration := infobip.NewConfiguration()
	configuration.Host = baseHost
	infobipClient := infobip.NewAPIClient(configuration)
	auth := context.WithValue(context.Background(), infobip.ContextAPIKey, smsAPIKey)

	request := infobip.NewSmsAdvancedTextualRequest()

	destNumber = strings.ReplaceAll(strings.ReplaceAll(destNumber, "+", ""), "-", "")
	destination := infobip.NewSmsDestination(destNumber)

	from := os.Getenv("TROVOWALLET_SMS_FROM")
	if len(from) == 0 {
		from = "TrovoWallet"
	}
	text := messageBody
	message := infobip.NewSmsTextualMessage()
	message.From = &from
	message.Destinations = &[]infobip.SmsDestination{*destination}
	message.Text = &text
	request.Messages = &[]infobip.SmsTextualMessage{*message}
	apiResponse, httpResponse, err := infobipClient.
		SendSmsApi.
		SendSmsMessage(auth).
		SmsAdvancedTextualRequest(*request).
		Execute()
	if err != nil {
		log.Printf("[SendSMS] send sms has error: %v\n", err)
		return &tErrors.CustomError{
			Param:      "mobile",
			Err:        "unable to send OTP at this time",
			ErrMessage: "Unable to send OTP at this time. Please try again later.",
		}
	}
	log.Printf("[SendSMS] API Response: %+v, httpResponse: %+v, Error: %v\n", apiResponse, httpResponse.Status, err)

	return nil

}

func SendSMSWithTermiiGateway(destNumber, messageBody string) error {

	baseHost := os.Getenv("TERMII_SMS_URL")
	if len(baseHost) == 0 {
		log.Println("[SendSMSWithTermiiGateway] TERMII_SMS_URL is not set")
		return &tErrors.ErrorTemporaryServerError{}
	}

	if len(os.Getenv("TERMII_SMS_API_KEY")) == 0 {

		log.Println("[SendSMSWithTermiiGateway] TERMII_SMS_API_KEY is not set")

		return &tErrors.ErrorTemporaryServerError{}
	}

	if len(os.Getenv("TERMII_SMS_SENDER_ID")) == 0 {

		log.Println("[SendSMSWithTermiiGateway] TERMII_SMS_SENDER_ID is not set")

		return &tErrors.ErrorTemporaryServerError{}
	}

	destNumber = strings.ReplaceAll(strings.ReplaceAll(destNumber, "+", ""), "-", "")

	text := messageBody

	senderID := os.Getenv("TERMII_SMS_SENDER_ID")
	channel := "generic"
	log.Println("[SendSMSWithTermiiGateway] sending sms to: ", destNumber, " with message: ", text)
	if strings.HasPrefix(destNumber, "234") {
		senderID = "N-Alert"
		channel = "dnd"
		log.Println("[SendSMSWithTermiiGateway] sending sms with sender ID: ", senderID, " through channel: ", channel)
	}

	url := fmt.Sprintf("https://%s/api/sms/send?to=%s&from=%s&sms=%s&type=plain&channel=%s&api_key=%s", os.Getenv("TERMII_SMS_URL"), destNumber, senderID, url.QueryEscape(text), channel, os.Getenv("TERMII_SMS_API_KEY"))
	log.Println("URL:", url)

	resp, err := http.Post(url, "application/json", nil)

	if err != nil {
		log.Printf("[SendSMSWithTermiiGateway] send sms has error: %v\n", err)
		return &tErrors.CustomError{
			Param:      "mobile",
			Err:        "unable to send OTP at this time",
			ErrMessage: "Unable to send OTP at this time. Please try again later.",
		}
	}
	defer resp.Body.Close()
	if resp.StatusCode > 204 {
		log.Printf("[SendSMSWithTermiiGateway] other send sms has error: %v\n", resp.Status)
		return &tErrors.CustomError{
			Param:      "mobile",
			Err:        "unable to send OTP at this time",
			ErrMessage: "Unable to send OTP at this time. Please try again later.",
		}
	}
	log.Printf("[SendSMSWithTermiiGateway] API Response: %+v\n", resp.Body)

	return nil
}

func getSMSProvider(destination string, db *gorm.DB) string {
	defaultProvider := os.Getenv("DEFAULT_SMS_PROVIDER")
	if len(defaultProvider) == 0 {
		defaultProvider = "termii"
	}

	phoneProfix := strings.Split(destination, "-")
	if len(phoneProfix) < 2 {
		//use default
		return defaultProvider
	}

	var phoneProvider SmsProvider

	err := db.First(&phoneProvider, "phone_prefix = ?", phoneProfix[0]).Error
	if err != nil {
		return defaultProvider
	}
	return phoneProvider.Provider
}
