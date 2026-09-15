package users

import "time"

// UserMobilePhoneVerification model for user phone verificationInfo
type UserMobilePhoneVerification struct {
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
	UserID           string    `gorm:"primaryKey"`
	User             User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Mobile           string    `gorm:"not null"`
	RequestDate      time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
	VerificationCode string
	// MobileVerified    int `gorm:"type:integer;not null;default:0"`
}
