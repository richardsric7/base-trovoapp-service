package mail

import (
	tErrors "admin-panel-dashboard/internal/errors"
	"admin-panel-dashboard/internal/models"
	serverModels "admin-panel-dashboard/internal/server/models"
	"context"
	"fmt"
	"log"
	"os"
	"time"

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
		sender = fmt.Sprintf("TrovoP2P Team <no-reply@%s>", mailgunDomain)
	}
	subject := os.Getenv("EMAIL_VERIFICATION_SUBJECT")
	if subject == "" {
		subject = "Your Trovo Wallet Email Verification Code"
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
func SendEmail(email, body, subject string) (id, resp string, err error) {
	// var mailgunDomain string = "sandbox33b89314ed134a399cf8df6b684395b7.mailgun.org" // e.g. mg.yourcompany.com
	var mailgunDomain string = os.Getenv("MAILGUN_DOMAIN") // e.g. mg.yourcompany.com
	log.Println("starting mail sending for email", email)
	// encode to base64
	// sEnc := base64.StdEncoding.EncodeToString([]byte(data))
	// Create an instance of the Mailgun Client
	mg := mailgun.NewMailgun(mailgunDomain, os.Getenv("MAILGUN_PRIVATE_API_KEY"))
	sender := os.Getenv("MAIL_SENDER")
	if sender == "" {
		sender = fmt.Sprintf("TrovoP2P Team <no-reply@%s>", mailgunDomain)
	}

	// body := ""
	// body := fmt.Sprintf("[%s]", sEnc)
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

func SendEmailWithMailgunTemplate(admin models.AdminUser, body, subject string) (id, resp string, err error) {
	// Fetch Mailgun domain and API key from environment variables
	var mailgunDomain string = os.Getenv("MAILGUN_DOMAIN")
	log.Println("Starting mail sending for email:", admin.Email)

	// Create an instance of the Mailgun Client
	mg := mailgun.NewMailgun(mailgunDomain, os.Getenv("MAILGUN_PRIVATE_API_KEY"))
	sender := os.Getenv("MAIL_SENDER")
	adminSender := os.Getenv("TROVO_ADMIN_MANAGER")
	if sender == "" {
		sender = fmt.Sprintf("%v <no-reply@%s>", adminSender, mailgunDomain)
	}

	// Set the recipient email
	recipient := admin.Email

	// Create the message object and set the template ID
	message := mg.NewMessage(sender, subject, body, recipient)
	message.SetTemplate("create_admin_email_template")
	// lets set the dashboard url
	err = message.AddVariable("dashboard_url", os.Getenv("ADMIN_DASHBOARD_URL"))
	if err != nil {
		return "", "", err
	}
	err = message.AddVariable("first_name", admin.FirstName)
	if err != nil {
		return "", "", err
	}
	// Set a timeout for sending the email
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Send the email message
	resp, id, err = mg.Send(ctx, message)
	if err != nil {
		log.Printf("Failed to send to: %s due to %v\n", admin.Email, err)
		err = &tErrors.CustomError{
			Param:      "email",
			Err:        "error could not send email",
			ErrMessage: "Could not send the email at this time. Please try again later",
			Code:       400,
		}
		return
	}

	log.Printf("ID: %s Resp: %s\n", id, resp)
	return
}

// SendOrganizationInviteEmail sends an organization invitation email using the organization_invite_email template
func SendOrganizationInviteEmail(admin models.AdminUser, inviteLink, organizationName, organizationType string, invite_expires_at time.Time) (id, resp string, err error) {
	// Fetch Mailgun domain and API key from environment variables
	var mailgunDomain string = os.Getenv("MAILGUN_DOMAIN")
	log.Println("Starting organization invite mail sending for email:", admin.Email)

	// Create an instance of the Mailgun Client
	mg := mailgun.NewMailgun(mailgunDomain, os.Getenv("MAILGUN_PRIVATE_API_KEY"))
	sender := os.Getenv("MAIL_SENDER")
	adminSender := os.Getenv("TROVO_ADMIN_MANAGER")
	if sender == "" {
		sender = fmt.Sprintf("%v <no-reply@%s>", adminSender, mailgunDomain)
	}

	// Set the recipient email
	recipient := admin.Email

	// Create the message object and set the template ID
	message := mg.NewMessage(sender, "Trovo Manager Organization Invitation", "", recipient)
	message.SetTemplate("organization_invite_email")

	// Calculate days until expiration
	daysUntilExpiration := int(invite_expires_at.Sub(time.Now()).Hours() / 24)
	if daysUntilExpiration < 1 {
		daysUntilExpiration = 1
	}

	// Get contact email from environment, default to support@trovotech.io
	contactEmail := os.Getenv("TROVO_CONTACT_EMAIL")
	if contactEmail == "" {
		contactEmail = "support@trovotech.io"
	}

	// Set template variables
	templateVars := map[string]string{
		"Year":          fmt.Sprintf("%d", time.Now().Year()),
		"contact_email": contactEmail,
		"days":          fmt.Sprintf("%d", daysUntilExpiration),
		"first_name":    admin.FirstName,
		"invite_link":   inviteLink,
		"organization":  organizationName,
		"role":          organizationType,
	}

	// Add all template variables
	for key, value := range templateVars {
		err = message.AddVariable(key, value)
		if err != nil {
			log.Printf("Failed to add template variable %s: %v", key, err)
			return "", "", err
		}
	}

	// Set a timeout for sending the email
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Send the email message
	resp, id, err = mg.Send(ctx, message)
	if err != nil {
		log.Printf("Failed to send organization invite to: %s due to %v\n", admin.Email, err)
		err = &tErrors.CustomError{
			Param:      "email",
			Err:        "error could not send organization invite email",
			ErrMessage: "Could not send the organization invitation at this time. Please try again later",
			Code:       400,
		}
		return
	}

	log.Printf("Organization invite email sent successfully. ID: %s Resp: %s\n", id, resp)
	return
}

// SendOrganizationOTPEmail sends an OTP email using the organization_otp_email template
func SendOrganizationOTPEmail(email, otpCode string, otpExpiresAt time.Time, server *serverModels.Server) (id, resp string, err error) {
	// Fetch Mailgun domain and API key from environment variables
	var mailgunDomain string = os.Getenv("MAILGUN_DOMAIN")
	log.Println("Starting organization OTP mail sending for email:", email)

	// Get user data from the database
	var member models.OrganizationMember
	if err := server.AdminDB.Where("email = ?", email).First(&member).Error; err != nil {
		log.Printf("Failed to get user data for email %s: %v", email, err)
		// Use email as fallback for first name
		member.FirstName = email
	}

	// Create an instance of the Mailgun Client
	mg := mailgun.NewMailgun(mailgunDomain, os.Getenv("MAILGUN_PRIVATE_API_KEY"))
	sender := os.Getenv("MAIL_SENDER")
	adminSender := os.Getenv("TROVO_ADMIN_MANAGER")
	if sender == "" {
		sender = fmt.Sprintf("%v <no-reply@%s>", adminSender, mailgunDomain)
	}

	// Set the recipient email
	recipient := email

	// Create the message object and set the template ID
	message := mg.NewMessage(sender, "Trovo Manager OTP Verification", "", recipient)
	message.SetTemplate("organization_otp_email")

	// Set template variables
	templateVars := map[string]string{
		"first_name": member.FirstName,
		"otp_code":   otpCode,
		"time":       otpExpiresAt.Format("2006-01-02 15:04:05"),
		"year":       fmt.Sprintf("%d", time.Now().Year()),
	}

	// Add all template variables
	for key, value := range templateVars {
		err = message.AddVariable(key, value)
		if err != nil {
			log.Printf("Failed to add template variable %s: %v", key, err)
			return "", "", err
		}
	}

	// Set a timeout for sending the email
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Send the email message
	resp, id, err = mg.Send(ctx, message)
	if err != nil {
		log.Printf("Failed to send organization OTP to: %s due to %v\n", email, err)
		err = &tErrors.CustomError{
			Param:      "email",
			Err:        "error could not send organization OTP email",
			ErrMessage: "Could not send the OTP verification at this time. Please try again later",
			Code:       400,
		}
		return
	}

	log.Printf("Organization OTP email sent successfully. ID: %s Resp: %s\n", id, resp)
	return
}
