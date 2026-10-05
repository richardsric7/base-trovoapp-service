package models

import "time"

// FeeExemptUser mirrors app-backend's fee_exempt_users table
// (sharedconfig.FeeExemptUser): accounts that pay no platform service fees
// - the platform's own trading or operations accounts. tm-api writes it
// directly, like the other admin-managed fee and catalog tables;
// app-backend reads it whenever it charges a fee.
type FeeExemptUser struct {
	Username  string    `gorm:"column:username;primaryKey;size:100" json:"username"`
	Reason    string    `gorm:"column:reason;size:255" json:"reason"`
	AddedBy   string    `gorm:"column:added_by;size:255" json:"addedBy"`
	CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
}

func (FeeExemptUser) TableName() string { return "fee_exempt_users" }
