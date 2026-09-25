package models

// BannedAddress model for bantu user directory info
type BannedAddress struct {
	Address string `gorm:"size:56;primaryKey" json:"address"`
	Reason  string `json:"reason"`
}
