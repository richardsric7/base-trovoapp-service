package firebase

import (
	"context"
	"encoding/base64"
	"log"
	"os"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"firebase.google.com/go/storage"
	"google.golang.org/api/option"
)

func GetFirebaseMessagingClient(ctx context.Context) (fcmClient *messaging.Client, usedContext context.Context, err error) {

	opts := []option.ClientOption{option.WithCredentialsJSON(getDecodedFireBaseKey())}
	if ctx == nil {
		ctx = context.Background()

	}
	app, err := firebase.NewApp(ctx, nil, opts...)
	if err != nil {
		log.Printf("[GetFirebaseMessagingClient] new firebase app: %s", err)
		return nil, ctx, err
	}

	fcmClient, err = app.Messaging(ctx)
	if err != nil {
		log.Printf("[GetFirebaseMessagingClient] error getting messaging client: %s", err)
		return nil, ctx, err
	}
	return fcmClient, ctx, nil

}

func GetFirebaseStorageClient(ctx context.Context) (storageClient *storage.Client, usedContext context.Context, err error) {

	opts := []option.ClientOption{option.WithCredentialsJSON(getDecodedFireBaseKey())}
	if ctx == nil {
		ctx = context.Background()

	}
	app, err := firebase.NewApp(ctx, nil, opts...)
	if err != nil {
		log.Printf("[GetFirebaseStorageClient] new firebase app: %s", err)
		return nil, ctx, err
	}

	storageClient, err = app.Storage(ctx)
	if err != nil {
		log.Printf("[GetFirebaseStorageClient] error getting storage client: %s", err)
		return nil, ctx, err
	}
	return storageClient, ctx, nil

}

func SendFirebaseMessage(recipient, title, body, imageURI string, dataPayload map[string]string, fcmClient *messaging.Client, ctx context.Context) (response string, err error) {
	if fcmClient == nil {
		fcmClient, ctx, err = GetFirebaseMessagingClient(ctx)
		if err != nil {
			log.Printf("[SendFirebaseMessage]error fetching firebaseMessaging client: %s\n", err)
			return
		}
	}
	// defer ctx.Done()
	response, err = fcmClient.Send(ctx, &messaging.Message{
		Notification: &messaging.Notification{
			Title:    title,
			Body:     body,
			ImageURL: imageURI,
		},
		Token: recipient,
		Data:  dataPayload,
	})
	log.Printf("[SendFirebaseMessage] response: %v, error: %v\n", response, err)
	return response, err

}

func SendFirebaseBroadcast(recipients []string, title, body, imageURI string, dataPayload map[string]string, fcmClient *messaging.Client, ctx context.Context) (*messaging.BatchResponse, error) {

	response, err := fcmClient.SendMulticast(ctx, &messaging.MulticastMessage{
		Notification: &messaging.Notification{
			Title:    title,
			Body:     body,
			ImageURL: imageURI,
		},
		Tokens: recipients,
		Data:   dataPayload,
	})
	log.Printf("[SendFirebaseBroadcast] response: %+v, error: %v\n", response, err)
	return response, err
}

func SendPushNotificationBroadcast(recipients []string, title, body, imageURI string, dataPayload map[string]string, fcmClient *messaging.Client, ctx context.Context) (*messaging.BatchResponse, error) {
	var err error
	response := new(messaging.BatchResponse)
	if len(recipients) <= 10000 {
		return SendFirebaseBroadcast(recipients, title, body, imageURI, dataPayload, fcmClient, ctx)
	}
	i := 0
	er := 10000
	for i <= len(recipients)-1 {
		// var rec []string

		response, err = SendFirebaseBroadcast(recipients[i:er-1], title, body, imageURI, dataPayload, fcmClient, ctx)
		if err != nil {
			log.Printf("[SendPushNotificationBroadcast] error from batch range [%v:%v]: %v\n", i, er-1, err)
		}
		if len(recipients)-er <= 10000 {
			i = er - 1
			er = len(recipients)

		} else {
			i = er - 1
			er += 10000
		}

	}
	return response, err
}

func getDecodedFireBaseKey() []byte {

	fireBaseAuthKey := os.Getenv("GC")

	decodedKey, err := base64.StdEncoding.DecodeString(fireBaseAuthKey)
	if err != nil {
		return nil
	}
	return decodedKey
}
