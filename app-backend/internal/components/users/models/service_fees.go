package users

import "time"

type ServiceFee struct {
	CreatedAt          time.Time `json:"createdAt"`
	UpdatedAt          time.Time `json:"updatedAt"`
	LastUpdatedBy      string    `json:"lastUpdatedBy"`
	ID                 string    `json:"id"`
	FeeWalletSecretKey string    `json:"feeWalletSecretKey"` //Secret key
	FeePercent         float64   `gorm:"default:0" json:"feePercent"`
	FeeFixed           float64   `gorm:"default:0" json:"feeFixed"`
	FeeAssetCode       string    `gorm:"default:''" json:"feeAssetCode"`
	FeeContractAddress string    `gorm:"default:''" json:"feeContractAddress"`
	Inactive           int       `gorm:"default:0" json:"inactive"`
	Remarks            string    `gorm:"default:''" json:"remarks"`
}
