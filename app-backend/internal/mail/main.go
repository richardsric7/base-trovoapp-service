package mail

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"time"
	tErrors "trovo-wallet-api/internal/errors"

	"github.com/mailgun/mailgun-go/v4"
)

func SendEmailVerificationCode(email, verificationCode string) (id, resp string, err error) {
	// var mailgunDomain string = "sandbox33b89314ed134a399cf8df6b684395b7.mailgun.org" // e.g. mg.yourcompany.com
	var mailgunDomain string = os.Getenv("MAILGUN_DOMAIN") // e.g. mg.yourcompany.com
	log.Println("starting mail sending for email", email)

	// Create an instance of the Mailgun Client
	mg := mailgun.NewMailgun(mailgunDomain, os.Getenv("MAILGUN_PRIVATE_API_KEY"))
	sender := os.Getenv("MAIL_SENDER")
	if sender == "" {
		sender = fmt.Sprintf("Trovotech <no-reply@%s>", mailgunDomain)
	}
	subject := os.Getenv("EMAIL_VERIFICATION_SUBJECT")
	if subject == "" {
		subject = "Your TrovoApp Email Verification Code"
	}
	// body := ""
	body := fmt.Sprintf("Your verification code is: %s", verificationCode)
	recipient := email

	// The message object allows you to add attachments and Bcc recipients
	message := mg.NewMessage(sender, subject, body, recipient)
	// message.SetTemplate("bantupay_v2")
	// message.SetTemplate(os.Getenv("EMAIL_VERIFICATION_TEMPLATE"))
	// message.AddTemplateVariable("verification_code", verificationCode)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Send the message with a 10 second timeout
	resp, id, err = mg.Send(ctx, message)
	if err != nil {
		log.Printf("failed to send to:%s due to %v\n", email, err)
		err = &tErrors.CustomError{Param: "email",
			Err:        "error could not send verification message",
			ErrMessage: "Could not send verification code at this time. Please try again later",
			Code:       400}
		return
	}
	log.Printf("ID: %s Resp: %s\n", id, resp)

	return

}
func SendEmail(email, data string) (id, resp string, err error) {
	// var mailgunDomain string = "sandbox33b89314ed134a399cf8df6b684395b7.mailgun.org" // e.g. mg.yourcompany.com
	var mailgunDomain string = os.Getenv("MAILGUN_DOMAIN") // e.g. mg.yourcompany.com
	log.Println("starting mail sending for email", email)
	// encode to base64
	sEnc := base64.StdEncoding.EncodeToString([]byte(data))
	// Create an instance of the Mailgun Client
	mg := mailgun.NewMailgun(mailgunDomain, os.Getenv("MAILGUN_PRIVATE_API_KEY"))
	sender := os.Getenv("MAIL_SENDER")
	if sender == "" {
		sender = fmt.Sprintf("Trovotech <no-reply@%s>", mailgunDomain)
	}
	subject := ""
	if subject == "" {
		subject = "new encoded CHANNELS"
	}
	// body := ""
	body := fmt.Sprintf("[%s]", sEnc)
	recipient := email

	// The message object allows you to add attachments and Bcc recipients
	message := mg.NewMessage(sender, subject, body, recipient)
	// message.SetTemplate("bantupay_v2")
	// message.SetTemplate(os.Getenv("EMAIL_VERIFICATION_TEMPLATE"))
	// message.AddTemplateVariable("verification_code", verificationCode)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Send the message with a 10 second timeout
	resp, id, err = mg.Send(ctx, message)
	if err != nil {
		log.Printf("failed to send to:%s due to %v\n", email, err)
		err = &tErrors.CustomError{Param: "email",
			Err:        "error could not send verification message",
			ErrMessage: "Could not send verification code at this time. Please try again later",
			Code:       400}
		return
	}
	log.Printf("ID: %s Resp: %s\n", id, resp)

	return

}
