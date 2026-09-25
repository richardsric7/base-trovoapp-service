package models

import (
	"log"
	"time"

	"gorm.io/gorm"
)

// User holds user data model

// type User struct {
//	ID                  string `gorm:"size:100"`
//	Username            string `gorm:"size:100;index:idx_username,unique;not null;check:,length(username) >= 2"`
//	CreatedAt           time.Time
//	UpdatedAt           time.Time
//	LastUpdatedMobileOn *time.Time
//	LastName            string `gorm:"size:50;index:idx_user_last_name;not null"`
//	FirstName           string `gorm:"size:50;index:idx_user_first_name;not null"`
//	//MiddleName               *string  `gorm:"size:50;index:idx_user_middle_name;null"`
//	Mobile *string `gorm:"size:50;null;index:idx_user_mobile_number,unique"`
//	//ContactPhone             *string  `gorm:"size:50;null;index:idx_user_contact_phone_number"`
//	//Telegram                 *int64   `gorm:"null;index:idx_user_telegram,unique"`
//	Email string `gorm:"not null;index:idx_user_email,unique"`
//	//ImageThumbnail           *string  `gorm:"null"`
//	CountryCode              *string  `gorm:"size:2;null"`
//	Latitude                 *float64 `gorm:"null"`
//	Longitude                *float64 `gorm:"null"`
//	City                     *string  `gorm:"null;size:100"`
//	Region                   *string  `gorm:"null;size:100"`
//	RegionName               *string  `gorm:"null;size:100"`
//	TimeZone                 *string  `gorm:"null;size:100"`
//	ISP                      *string  `gorm:"null;size:150"`
//	PublicIP                 *string  `gorm:"null;size:100"`
//	PublicKey                string   `gorm:"null;size:100"`
//	PrimarySigner            *bool    `gorm:"null"`
//	Referrer                 *string  `gorm:"null"`
//	ReferralLink             *string  `gorm:"null"`
//	ReferralQRCode           *string  `gorm:"null"`
//	PushNotificationToken    *string  `gorm:"null"`
//	Corporate                *bool    `gorm:"null"`
//	MobileVerified           *bool    `gorm:"null"`
//	MembershipType           *string  `gorm:"null;size:100"`
//	MembershipExpiry         *time.Time
//	KYCVerified              *bool      `gorm:"null"`
//	AccountRecoveryEnabled   *bool      `gorm:"null"`
//	Verified                 *bool      `gorm:"null"`
//	Suspended                uint       `gorm:"type:integer;not null;default:0"`
//	SuspensionReason         *string    `gorm:"null"`
//	HasSecurityQuestions     *bool      `gorm:"null"`
//	AccountRecoveryExpiresOn *time.Time `gorm:"null"`
//	LastRecoveredAccountOn   *time.Time `gorm:"null"`
//	//BantuTalk                *string    `gorm:"size:100;index:idx_user_bantu_talk,unique;null"`
//	//KYCLevel              uint      `gorm:"type:integer;not null;default:0"`
//	//MaxAssetPerOffer      float64   `gorm:"type:integer;not null;default:0"` //XBN
//	//MaxAssetPerOrder      float64   `gorm:"type:integer;not null;default:0"` //XBN
//	Offline               uint      `gorm:"not null;default:0" json:"offline"`
//	AdminLevel            uint      `gorm:"not null;default:0" json:"adminLevel"`
//	TelegramNotifications uint      `gorm:"not null;default:0" json:"telegramNotifications"`
//	TelegramConnectedAt   time.Time `json:"telegram_connected_at"`
//}

//// User holds user data model
// type User struct {
//	ID                  string `gorm:"size:100"`
//	Username            string `gorm:"size:100;index:idx_username,unique;not null"`
//	CreatedAt           time.Time
//	UpdatedAt           time.Time
//	LastUpdatedMobileOn *time.Time
//	FirstName           string  `gorm:"size:50;not null"`
//	LastName            string  `gorm:"size:50;not null"`
//	Mobile              *string `gorm:"size:50;null"`
//	Email               string  `gorm:"not null;index:idx_user_email,unique"`
//	//ImageThumbnail           *string  `gorm:"null"`
//	CountryCode *string  `gorm:"size:2;null"`
//	Latitude    *float64 `gorm:"null"`
//	Longitude   *float64 `gorm:"null"`
//	City        *string  `gorm:"null;size:100"`
//	Region      *string  `gorm:"null;size:100"`
//	RegionName  *string  `gorm:"null;size:100"`
//	TimeZone    *string  `gorm:"null;size:100"`
//	ISP         *string  `gorm:"null;size:150"`
//	PublicIP    *string  `gorm:"null;size:100"`
//	PublicKey   string   `gorm:"null;size:100"`
//	//PrimarySigner            *bool    `gorm:"null"`
//	Referrer                 *string `gorm:"null"`
//	ReferralLink             *string `gorm:"null"`
//	ReferralQRCode           *string `gorm:"null"`
//	PushNotificationToken    *string `gorm:"null"`
//	Corporate                *bool   `gorm:"null"`
//	MobileVerified           *bool   `gorm:"null"`
//	MembershipType           *string `gorm:"null;size:100"`
//	MembershipExpiry         *time.Time
//	KYCVerified              *bool      `gorm:"null"`
//	AccountRecoveryEnabled   *bool      `gorm:"null"`
//	Verified                 *bool      `gorm:"null"`
//	Suspended                uint       `gorm:"type:integer;not null;default:0"`
//	SuspensionReason         *string    `gorm:"null"`
//	HasSecurityQuestions     *bool      `gorm:"null"`
//	AccountRecoveryExpiresOn *time.Time `gorm:"null"`
//	LastRecoveredAccountOn   *time.Time `gorm:"null"`
//	// Status                   string     `gorm:"not null;default:ACTIVE"`
//}

// User represents the structure of the users table in PostgreSQL.
type User struct {
	CreatedAt                time.Time `json:"created_at" db:"created_at"`
	UpdatedAt                time.Time `json:"updated_at" db:"updated_at"`
	LastUpdatedMobileOn      time.Time `json:"last_updated_mobile_on" db:"last_updated_mobile_on"`
	ID                       string    `json:"id" db:"id"`
	Username                 string    `json:"username" db:"username"`
	Email                    string    `json:"email,omitempty" db:"email"`
	ImageThumbnailURL        string    `json:"image_thumbnail_url,omitempty" db:"image_thumbnail_url"`
	FirstName                string    `json:"first_name" db:"first_name"`
	LastName                 string    `json:"last_name,omitempty" db:"last_name"`
	Mobile                   string    `json:"mobile,omitempty" db:"mobile"`
	Address                  string    `json:"address,omitempty" db:"address"`
	PrimarySigner            string    `json:"primary_signer,omitempty" db:"primary_signer"`
	Referrer                 string    `json:"referrer,omitempty" db:"referrer"`
	ReferralLink             string    `json:"referral_link,omitempty" db:"referral_link"`
	ReferralQRCode           string    `json:"referral_qr_code,omitempty" db:"referral_qr_code"`
	PushNotificationToken    string    `json:"push_notification_token,omitempty" db:"push_notification_token"`
	Corporate                int       `json:"corporate" db:"corporate"`
	MobileVerified           int       `json:"mobile_verified" db:"mobile_verified"`
	MembershipType           int       `json:"membership_type" db:"membership_type"`
	MembershipExpiry         time.Time `json:"membership_expiry,omitempty" db:"membership_expiry"`
	KYCVerified              int       `json:"kyc_verified" db:"kyc_verified"`
	AccountRecoveryEnabled   int       `json:"account_recovery_enabled" db:"account_recovery_enabled"`
	PublicIP                 string    `json:"public_ip,omitempty" db:"public_ip"`
	CountryCode              string    `json:"country_code,omitempty" db:"country_code"`
	Latitude                 float64   `json:"latitude,omitempty" db:"latitude"`
	Longitude                float64   `json:"longitude,omitempty" db:"longitude"`
	City                     string    `json:"city,omitempty" db:"city"`
	Region                   string    `json:"region,omitempty" db:"region"`
	RegionName               string    `json:"region_name,omitempty" db:"region_name"`
	TimeZone                 string    `json:"time_zone,omitempty" db:"time_zone"`
	ISP                      string    `json:"isp,omitempty" db:"isp"`
	Verified                 int       `json:"verified" db:"verified"`
	Suspended                int       `json:"suspended" db:"suspended"`
	SuspensionReason         string    `json:"suspension_reason,omitempty" db:"suspension_reason"`
	HasSecurityQuestions     int       `json:"has_security_questions" db:"has_security_questions"`
	AccountRecoveryExpiresOn time.Time `json:"account_recovery_expires_on,omitempty" db:"account_recovery_expires_on"`
	LastRecoveredAccountOn   time.Time `json:"last_recovered_account_on,omitempty" db:"last_recovered_account_on"`
}

type UserLastLogin struct {
	ID          string    `gorm:"size:100"`
	Username    string    `gorm:"size:100;index:idx_username_last_login;not null;check:,length(username) >= 2"`
	CreatedAt   time.Time `gorm:"default:now()"`
	CountryCode *string   `gorm:"size:2;null"`
	Latitude    *float64  `gorm:"null"`
	Longitude   *float64  `gorm:"null"`
	City        *string   `gorm:"null;size:100"`
	Region      *string   `gorm:"null;size:100"`
	RegionName  *string   `gorm:"null;size:100"`
	TimeZone    *string   `gorm:"null;size:100"`
	ISP         *string   `gorm:"null;size:150"`
	PublicIP    *string   `gorm:"null;size:100"`
}

type UserWallet struct {
	CreatedAt               time.Time          `json:"createdAt"`
	UpdatedAt               time.Time          `json:"updatedAt"`
	ID                      string             `gorm:"size:56" json:"publicKey"`
	TempAddress             *string            `gorm:"size:56;index:idx_user_wallet_temp_key;null"`
	Tag                     *string            `gorm:"null;size:12" json:"tag"`
	Description             *string            `gorm:"null;size:100" json:"description"`
	Alias                   string             `gorm:"size:30; index:idx_unique_alias, unique" json:"alias"` // primaryUsername_tag for sub wallets
	Signer                  string             `gorm:"size:56; index:idx_user_wallet_signer" json:"signer"`  // if ID is same as signer, then it is a primary wallet
	UserID                  string             `gorm:"type:integer;not null; default:0;index:idx_user_wallets_user_id" json:"userId"`
	SharedAccessEnabled     int                `gorm:"type:integer;not null; default:0" json:"sharedAccessEnabled"`
	Tracked                 int                `gorm:"type:integer;not null;default:0" json:"-"`
	PrimaryWallet           int                `gorm:"type:integer;not null;default:0" json:"primaryWallet"`
	NumberOfApprovalsNeeded int                `gorm:"type:integer; default:0" json:"numberOfApprovalsNeeded"`
	WalletType              int                `gorm:"type:integer; default:0" json:"walletType"` // 0=normal, 1= assetIssuing, 2= marketMaking, 3 = bulkPayment
	Permissions             []WalletPermission `gorm:"foreignKey:WalletAddress;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"permissions"`
	SharedAccessCreatedAt   time.Time          `json:"sharedAccessCreatedAt"`
	SharedAccessUpdatedAt   time.Time          `json:"sharedAccessUpdatedAt"`
	FeeDisabled             int                `gorm:"type:integer; not null; default:0" json:"feeDisabled"`
}

type WalletPermission struct {
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
	ID             string    `json:"id"`
	WalletAddress  string    `gorm:"size:90;not null; index:access_level_permission,unique;index:idx_address_shared" json:"walletAddress"`
	TargetUsername string    `gorm:"size:16;not null; index:access_level_permission,unique;" json:"targetUsername"`
	Permission     string    `gorm:"size:10;not null; index:access_level_permission,unique" json:"permission"`
}

// // GetUser gets user information
// func (u *User) ToUserInfo(db *gorm.DB) (userInfo UserInfo) {
//
//	if u.Suspended == 1 {
//		return userInfo
//	}
//
//	userInfo = UserInfo{
//		ID:        u.ID,
//		Username:  u.Username,
//		PublicKey: u.PublicKey,
//		Email:     u.Email,
//		LastName:  u.LastName,
//		FirstName: u.FirstName,
//		//KYCLevel:              u.KYCLevel,
//		//AdminLevel:            u.AdminLevel,
//		//Offline:               u.Offline,
//		//TelegramNotifications: u.TelegramNotifications,
//		//MaxAssetPerOffer:      u.MaxAssetPerOffer,
//		//MaxAssetPerOrder:      u.MaxAssetPerOrder,
//		Suspended: u.Suspended,
//	}
//
//	if u.Mobile != nil {
//		userInfo.Mobile = *u.Mobile
//	}
//	//if u.Telegram != nil {
//	//	userInfo.Telegram = u.Telegram
//	//}
//	//if u.ContactPhone != nil {
//	//	userInfo.ContactPhone = *u.ContactPhone
//	//}
//
//	if u.CountryCode != nil {
//		userInfo.CountryCode = *u.CountryCode
//	}
//	//if u.ImageThumbnail != nil {
//	//	userInfo.ImageThumbnail = *u.ImageThumbnail
//	//}
//	if len(strings.ReplaceAll(os.Getenv("TAKER_FEE"), " ", "")) > 0 {
//		userInfo.TakerFee = strings.ReplaceAll(os.Getenv("TAKER_FEE"), " ", "")
//	} else {
//		userInfo.TakerFee = "0"
//	}
//	userInfo.TradeStats = userInfo.GetMakerStat(db)
//
//	//userInfo.TakerReputation = userInfo.GetTakerReputation(db)
//
//	return userInfo
//
// }
func (u *User) ToggleOffline(db *gorm.DB) (updatedState uint, err error) {

	type User struct {
		ID       string `gorm:"size:100"`
		Username string `gorm:"size:100"`
		Offline  uint   `gorm:"not null;default:0" json:"offline"`
	}
	var userState User
	// username := u.Username
	//log.Printf("[ToggleOffline] User:[%v], Current Offline State:[%v]\n", u.Username, u.Offline)
	//if u.Offline == 1 {
	//
	//	log.Printf("[ToggleOffline] User:[%v], Switching ONLINE\n", u.Username)
	//
	//	err = db.Raw(`WITH o as (update offers set offline = 0 where maker = ?),
	//u as (update users set offline = 0, updated_at = now() where username = ? returning id, username, offline)
	//select * from u`, u.Username, username).Scan(&userState).Error
	//
	//} else {
	//
	//	log.Printf("[ToggleOffline] User:[%v], Switching OFFLINE\n", u.Username)
	//
	//	err = db.Raw(`WITH o as (update offers set offline = 1 where maker = ?),
	//	u as (update users set offline = 1, updated_at = now() where username = ? returning id, username, offline)
	//	select * from u`, u.Username, username).Scan(&userState).Error
	//
	//}

	// if err != nil {
	//	log.Printf("[ToggleOffline] error fetching user data: [%v]\n", err)
	//	return 0, &p2pErrors.CustomError{
	//		Param:      "offline",
	//		Err:        "error toggling offline state",
	//		ErrMessage: "Unable to switch user offline state",
	//	}
	//}

	// u.Offline = budsState
	updatedState = userState.Offline
	log.Printf("[ToggleOffline] User:[%v], State of Offline after execution:[%v]\n", u.Username, userState.Offline)

	return updatedState, nil
}
