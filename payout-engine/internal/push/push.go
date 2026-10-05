// Package push sends Firebase Cloud Messaging notifications, with the same
// credentials app-backend uses (GC: base64 of the service account JSON).
package push

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"google.golang.org/api/option"
)

// Sender sends one notification to a device token.
type Sender interface {
	Send(ctx context.Context, token, title, body string, data map[string]string) error
}

// FCM sends through Firebase.
type FCM struct{ client *messaging.Client }

// NewFromEnv returns an FCM sender, or nil when GC is not set (notifications
// are then skipped).
func NewFromEnv(ctx context.Context) (Sender, error) {
	raw := strings.TrimSpace(os.Getenv("GC"))
	if raw == "" {
		return nil, nil
	}
	key, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		return nil, fmt.Errorf("GC is not base64: %w", err)
	}
	app, err := firebase.NewApp(ctx, nil, option.WithCredentialsJSON(key))
	if err != nil {
		return nil, err
	}
	client, err := app.Messaging(ctx)
	if err != nil {
		return nil, err
	}
	return &FCM{client: client}, nil
}

func (f *FCM) Send(ctx context.Context, token, title, body string, data map[string]string) error {
	_, err := f.client.Send(ctx, &messaging.Message{Notification: &messaging.Notification{Title: title, Body: body}, Token: token, Data: data})
	return err
}
