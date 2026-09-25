package usermetrics

import (
	"time"

	"gorm.io/gorm"
)

// ServiceLink represents the service_links table
type ServiceLink struct {
	ID                         string    `json:"id" gorm:"primaryKey"`
	CreatedAt                  time.Time `json:"created_at"`
	UpdatedAt                  time.Time `json:"updated_at"`
	OwnerUsername              string    `json:"owner_username"`
	ShortName                  string    `json:"short_name"`
	LongName                   string    `json:"long_name"`
	LoginPermission            int       `json:"login_permission"`
	PaymentPermission          int       `json:"payment_permission"`
	AuthorizationPermission    int       `json:"authorization_permission"`
	AllowUserInfo              int       `json:"allow_user_info"`
	PushNotificationPermission int       `json:"push_notification_permission"`
	IncludePhoneNumbers        int       `json:"include_phone_numbers"`
	IncludeUserBalances        int       `json:"include_user_balances"`
	Verified                   int       `json:"verified"`
	RewardOnly                 int       `json:"reward_only"`
	Inactive                   int       `json:"inactive"`
	Suspended                  int       `json:"suspended"`
}

// GetAllServiceLinks retrieves all service links
func GetAllServiceLinks(db *gorm.DB) ([]ServiceLink, error) {
	var links []ServiceLink
	if err := db.Order("created_at desc").Find(&links).Error; err != nil {
		return nil, err
	}
	return links, nil
}
