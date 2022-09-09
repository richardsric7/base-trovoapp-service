package merchants

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/bantublockchain/push-notification-service/pkg/shove"
)

// ServiceLink holds ServiceLink data model
type ServiceLink struct {
	ID                         string    `json:"-" gorm:"size:100"`
	CreatedAt                  time.Time `json:"-"`
	UpdatedAt                  time.Time `json:"-"`
	OwnerUsername              string    `json:"TrovoUsername" gorm:"size:100;index:idx_owner_username;index:idx_service_shortname,unique;not null;check:,length(owner_username) > 2"`
	PublicKey                  string    `json:"publicKey" gorm:"size:56;index:idx_service_user_public_key;not null;"`
	ApiKey                     string    `json:"apiKey" gorm:"size:50;index:idx_service_api_key,unique;not null;"`
	ShortName                  string    `json:"shortName" gorm:"size:50;index:idx_service_shortname,unique;not null;"`
	LongName                   string    `json:"longName" gorm:"size:100"`
	LoginPermission            int       `json:"-" gorm:"type:integer;not null;default:0"`
	PaymentPermission          int       `json:"-" gorm:"type:integer;not null;default:0"`
	AuthorizationPermission    int       `json:"-" gorm:"type:integer;not null;default:0"`
	AllowUserInfo              int       `json:"-" gorm:"type:integer;not null;default:0"`
	PushNotificationPermission int       `json:"-" gorm:"type:integer;not null;default:0"`
	IncludePhoneNumbers        int       `json:"-" gorm:"type:integer;not null;default:0"`
	IncludeUserBalances        int       `json:"-" gorm:"type:integer;not null;default:0"`
	Verified                   int       `json:"-" gorm:"type:integer;not null;default:0"`
	RewardOnly                 int       `json:"-" gorm:"type:integer;not null;default:0"`
	Inactive                   int       `json:"inactive" gorm:"type:integer;not null;default:0"`
	Suspended                  int       `json:"-" gorm:"type:integer;not null;default:0"`
	SuspensionReason           *string   `json:"-" gorm:"null"`
}
type ServiceLinkApiKeyLog struct {
	ID            int64     `json:"-"`
	CreatedAt     time.Time `json:"-"`
	UpdatedAt     time.Time `json:"-"`
	ServiceLinkID string    `json:"-"`
	ApiKey        string    `json:"apiKey" gorm:"size:50;index:idx_old_service_api_key,unique;not null;"`
}

// ServiceLinkLoginSession holds user data model
type ServiceLinkLoginSession struct {
	ID             string `gorm:"size;primaryKey"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	ApiKey         string  `gorm:"size:50"`
	OwnerUsername  string  `gorm:"size:100;index:idx_loginsession;not null;check:,length(owner_username) >= 2"`
	WalletUsername string  `gorm:"size:100;not null;index:idx_loginsession"`
	CallbackURL    *string `gorm:"null"`
	Authorized     int     `gorm:"type:integer;not null;default:0"`
}

// ServiceAuthorization holds authorization data model
type ServiceLinkAuthorization struct {
	ID             string `gorm:"size:100;primaryKey"`
	CreatedAt      time.Time
	ExpiresAt      time.Time `gorm:"default:now()"`
	UpdatedAt      time.Time
	ApiKey         string  `gorm:"size:50"`
	OwnerUsername  string  `gorm:"size:100;index:idx_authdata;not null;check:,length(owner_username) >= 2"`
	WalletUsername string  `gorm:"not null;index:idx_authdata"`
	CallbackURL    *string `gorm:"null"`
	Authorized     int     `gorm:"type:integer;not null;default:0"`
}

type ServiceLinkRequestInput struct {
	AuthDescription   string `json:"authDescription,omitempty"`
	DeviceInfo        string `json:"deviceInfo,omitempty"`
	CallbackURL       string `json:"callbackUrl,omitempty"`
	ValidityInMinutes int    `json:"validityInMinutes,omitempty"`
}

type ServiceLinkPushNotificationInput struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}
type Android struct {
	Priority     string               `json:"priority"`
	Visibility   string               `json:"visibility"`
	Notification *AndroidNotification `json:"notification"`
}
type AndroidNotification struct {
	Priority string `json:"notification_priority"`
}

type APNSHeaders struct {
	Priority string `json:"apns-priority,omitempty"`
}

type APNS struct {
	Headers APNSHeaders `json:"headers,omitempty"`
}
type Message struct {
	Title   string `json:"title"`
	Message string `json:"body"`
}
type payload struct {
	To              string   `json:"to,omitempty"`
	Token           string   `json:"token,omitempty"`
	RegistrationIDs []string `json:"registration_ids,omitempty"`
	Notification    Message  `json:"notification"`
	Android         *Android `json:"android,omitempty"`
	APNS            *APNS    `json:"apns,omitempty"`
}

// ServiceLinkBudsInfo model for bantu user directory info
type ServiceLinkBudsInfo struct {
	CreatedAt             string                 `json:"createdAt"`
	Username              string                 `json:"username"`
	PublicKey             string                 `json:"publicKey"`
	Email                 string                 `json:"email"`
	LastName              string                 `json:"lastName"`
	FirstName             string                 `json:"firstName"`
	MiddleName            string                 `json:"middleName"`
	Mobile                string                 `json:"mobile"`
	BantuTalk             string                 `json:"bantuTalk"`
	ImageThumbnail        string                 `json:"imageThumbnail"`
	Verified              int                    `json:"verified"`
	MobileVerified        int                    `json:"mobileVerified"`
	Suspended             int                    `json:"suspended"`
	PushNotificationToken string                 `json:"-"`
	Referrer              string                 `json:"referrer"`
	Wallet                UserBalanceForMerchant `json:"wallet"`
}

func (m *ServiceLinkPushNotificationInput) PushMessage(token string) {

	if os.Getenv("PUSH_NOTIFICATION_SERVICE_MODE") == "redis" {
		//use redis queue
		log.Println("Using redis for PNS:", os.Getenv("PNS_REDIS_HOST"))
		type FCMNotification struct {
			To              string   `json:"to,omitempty"`
			Token           string   `json:"token,omitempty"`
			RegistrationIDs []string `json:"registration_ids,omitempty"`
			Notification    Message  `json:"notification"`
			Android         *Android `json:"android,omitempty"`
			APNS            *APNS    `json:"apns,omitempty"`
		}
		redisURL := fmt.Sprintf("redis://%v:%v", os.Getenv("PNS_REDIS_HOST"), os.Getenv("PNS_REDIS_PORT"))
		log.Println("Using redis for PNS:", redisURL)

		pwd := os.Getenv("PNS_REDIS_PASSWORD")

		client := shove.NewRedisClient(redisURL, pwd)
		message := Message{
			Title:   m.Title,
			Message: m.Message,
		}
		notification := FCMNotification{
			Token:        token,
			Notification: message,
		}

		raw, err := json.Marshal(notification)
		if err != nil {
			log.Printf("[PushMessage] could not unmarshal notification due to [%v]\n", err)
		}
		err = client.PushRaw("fcm", raw)
		if err != nil {
			log.Printf("[PushMessage] could not save raw redis message due to [%v]\n", err)
		} else {
			log.Printf("[PushMessage] Queued Message : [%v]\n", string(raw))

		}

		return
	}
	//use url
	if len(os.Getenv("PUSH_NOTIFICATION_SERVICE_URL")) > 10 && os.Getenv("PUSH_NOTIFICATION_SERVICE_MODE") != "redis" {
		//url exists, use it to push
		//make callback request
		log.Printf("[PushMessage] send FCM POST message [%v] to %v \n", m.Message, os.Getenv("PUSH_NOTIFICATION_SERVICE_URL"))

		message := Message{
			Title:   m.Title,
			Message: m.Message,
		}
		//https://pns-alpha.dev.bantupay.org/api/push/fcm
		jsonPayload := payload{Token: token, Notification: message}
		body, err := json.Marshal(jsonPayload)
		if err != nil {
			log.Printf("[PushMessage] could not send unmarshal message due to [%v]\n", err)

		}
		log.Printf("JSON STRING: [%v]\n", string(body))

		responseBody := bytes.NewBuffer(body)
		//Leverage Go's HTTP Post function to make request
		resp, err := http.Post(os.Getenv("PUSH_NOTIFICATION_SERVICE_URL"), "application/json", responseBody)
		//Handle Error
		if err != nil {
			log.Printf("[PushMessage] could not send FCM POST message due to [%v]\n", err)
			return
		}
		defer resp.Body.Close()
		//Read the response body
		body, err = io.ReadAll(resp.Body)
		if err == nil {

			log.Printf("[PushMessage] Response Body: [%v]\n", string(body))

		}

	}
}

func (m *ServiceLinkPushNotificationInput) PushMessage999Max(tokens []string) {
	if len(tokens) == 0 || len(tokens) > 999 {
		return
	}
	if os.Getenv("PUSH_NOTIFICATION_SERVICE_MODE") == "redis" {
		//use redis queue

		/////

		redisURL := fmt.Sprintf("redis://%v:%v", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
		pwd := os.Getenv("REDIS_PASSWORD")

		client := shove.NewRedisClient(redisURL, pwd)
		message := Message{
			Title:   m.Title,
			Message: m.Message,
		}
		var notification payload
		if len(tokens) < 1000 {
			notification = payload{
				RegistrationIDs: tokens,
				Notification:    message,
			}

			raw, err := json.Marshal(notification)
			if err != nil {
				log.Printf("[PushMessage] could not unmarshal notification due to [%v]\n", err)
			}
			err = client.PushRaw("fcm", raw)
			if err != nil {
				log.Printf("[PushMessage] could not save raw redis message due to [%v]\n", err)
			}
		}
		return
	}
	//use url
	if len(os.Getenv("PUSH_NOTIFICATION_SERVICE_URL")) > 10 && os.Getenv("PUSH_NOTIFICATION_SERVICE_MODE") != "redis" {
		//url exists, use it to push
		//make callback request
		log.Printf("[PushMessage] send FCM POST message [%v] to %v \n", m.Message, os.Getenv("PUSH_NOTIFICATION_SERVICE_URL"))

		message := Message{
			Title:   m.Title,
			Message: m.Message,
		}
		//https://pns-alpha.dev.bantupay.org/api/push/fcm

		if len(tokens) < 1000 {
			jsonPayload := payload{RegistrationIDs: tokens, Notification: message}
			body, err := json.Marshal(jsonPayload)
			if err != nil {
				log.Printf("[PushMessage] could not send unmarshal message due to [%v]\n", err)

			}
			log.Printf("JSON STRING: [%v]\n", string(body))

			responseBody := bytes.NewBuffer(body)
			//Leverage Go's HTTP Post function to make request
			resp, err := http.Post(os.Getenv("PUSH_NOTIFICATION_SERVICE_URL"), "application/json", responseBody)
			//Handle Error
			if err != nil {
				log.Printf("[PushMessage] could not send FCM POST message due to [%v]\n", err)
				return
			}
			defer resp.Body.Close()
			//Read the response body
			body, err = io.ReadAll(resp.Body)
			if err == nil {

				log.Printf("[PushMessage] Response Body: [%v]\n", string(body))

			}
		}

	}
}

func (m *ServiceLinkPushNotificationInput) PushBulkMessage(tokens []string) {

	if len(tokens) < 1000 {
		m.PushMessage999Max(tokens)
	} else {
		//process max of 999 per batch
		batch := make([]string, 0)

		for _, t := range tokens {
			batch = append(batch, t)
			if len(batch) == 999 {
				//send and reset
				m.PushMessage999Max(batch)
				//reset batch
				batch = make([]string, 0)
			}
		}
		//check if tokens still remianed in batch
		if len(batch) > 0 {
			m.PushMessage999Max(batch)
		}

	}

}
