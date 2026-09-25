package utils

import (
	"context"
	"fmt"
	"os"

	"github.com/mailgun/mailgun-go/v4"
)

// SendEmailWithMailgunTemplate sends an email using Mailgun's template system
func SendEmailWithMailgunTemplate(to, subject, template string, data map[string]interface{}) error {
	mg := mailgun.NewMailgun(os.Getenv("MAILGUN_DOMAIN"), os.Getenv("MAILGUN_API_KEY"))

	message := mg.NewMessage(
		fmt.Sprintf("Admin Dashboard <noreply@%s>", os.Getenv("MAILGUN_DOMAIN")),
		subject,
		"",
		to,
	)

	// Set template variables
	for k, v := range data {
		message.AddVariable(k, v) //nolint:errcheck //
	}

	// Set template name
	message.SetTemplate(template)

	ctx := context.Background()
	_, _, err := mg.Send(ctx, message)
	return err
}
