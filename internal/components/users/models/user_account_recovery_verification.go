package users

import "time"

// UserAccountRecoveryEmailVerification model for user account recovery verificationInfo
type UserAccountRecoveryEmailVerification struct {
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
	UserID           string    `gorm:"primaryKey"`
	User             User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Email            string    `gorm:"not null"`
	RequestDate      time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
	VerificationCode string
}
