package users

import "time"

type ServiceFee struct {
	CreatedAt          time.Time `json:"createdAt"`
	UpdatedAt          time.Time `json:"updatedAt"`
	LastUpdatedBy      string    `json:"lastUpdatedBy"`
	ID                 string    `json:"id"`
	FeeWalletSecretKey string    `json:"feeWalletSecretKey"` //Secret key
	FeePercent         float64   `gorm:"default:0" json:"feePercent"`
	Inactive           int       `gorm:"default:0" json:"inactive"`
}
