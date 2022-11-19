package servicelinks

import (
	"time"
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
	EventPermission            int       `json:"-" gorm:"type:integer;not null;default:0"`
	AllowUserInfo              int       `json:"-" gorm:"type:integer;not null;default:0"`
	PushNotificationPermission int       `json:"-" gorm:"type:integer;not null;default:0"`
	IncludePhoneNumbers        int       `json:"-" gorm:"type:integer;not null;default:0"`
	IncludeUserBalances        int       `json:"-" gorm:"type:integer;not null;default:0"`
	Verified                   int       `json:"-" gorm:"type:integer;not null;default:0"`
	// RewardOnly                 int       `json:"-" gorm:"type:integer;not null;default:0"`
	Inactive         int     `json:"inactive" gorm:"type:integer;not null;default:0"`
	Suspended        int     `json:"-" gorm:"type:integer;not null;default:0"`
	SuspensionReason *string `json:"-" gorm:"null"`
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

// ServiceLinkEvent holds event data model
type ServiceLinkEvent struct {
	ID            string `gorm:"size:100;primaryKey"`
	CreatedAt     time.Time
	ExpiresAt     time.Time `gorm:"default:now()"`
	UpdatedAt     time.Time
	ApiKey        string  `gorm:"size:50"`
	OwnerUsername string  `gorm:"size:100;index:idx_eventdata;not null;check:,length(owner_username) >= 2"`
	CallbackURL   *string `gorm:"null"`
}

type ServiceLinkRequestInput struct {
	AuthDescription   string `json:"authDescription,omitempty"`
	DeviceInfo        string `json:"deviceInfo,omitempty"`
	CallbackURL       string `json:"callbackUrl,omitempty"`
	ValidityInMinutes int    `json:"validityInMinutes,omitempty"`
}
type ServiceLinkEventRequestInput struct {
	EventDescription  string `json:"eventDescription,omitempty"`
	DeviceInfo        string `json:"deviceInfo,omitempty"`
	CallbackURL       string `json:"callbackUrl,omitempty"`
	ValidityInMinutes int    `json:"validityInMinutes,omitempty"`
}

type ServiceLinkPushNotificationInput struct {
	Title    string `json:"title"`
	Message  string `json:"message"`
	ImageURI string `json:"imageUri"`
	Action   string `json:"action"` //login, payment, 2fa, event
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
