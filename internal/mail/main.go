package mail

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/mailgun/mailgun-go/v4"
)

func SendEmailVerificationCode(email, verificationCode string) (id, resp string, err error) {
	var mailgunDomain string = "sandbox33b89314ed134a399cf8df6b684395b7.mailgun.org" // e.g. mg.yourcompany.com
	log.Println("starting mail sending for email", email)

	// Create an instance of the Mailgun Client
	mg := mailgun.NewMailgun(mailgunDomain, os.Getenv("MAILGUN_PRIVATE_API_KEY"))
	sender := os.Getenv("MAIL_SENDER")
	if sender == "" {
		sender = "Trovotech <noreply@email.bantupay.org>"
	}
	subject := os.Getenv("EMAIL_VERIFICATION_SUBJECT")
	if subject == "" {
		subject = "Your Trovo Wallet Email Verification Code"
	}
	body := ""
	recipient := email

	// The message object allows you to add attachments and Bcc recipients
	message := mg.NewMessage(sender, subject, body, recipient)
	// message.SetTemplate("bantupay_v2")
	message.SetTemplate(os.Getenv("EMAIL_VERIFICATION_TEMPLATE"))
	message.AddTemplateVariable("verification_code", verificationCode)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Send the message with a 10 second timeout
	resp, id, err = mg.Send(ctx, message)
	if err != nil {
		log.Println("failed to send to:\n", email)
		return
	}
	log.Printf("ID: %s Resp: %s\n", id, resp)

	return

}
