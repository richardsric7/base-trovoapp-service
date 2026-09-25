package models

// BannedPublicKey model for bantu user directory info
type BannedPublicKey struct {
	PublicKey string `gorm:"size:56;primaryKey" json:"publicKey"`
	Reason    string `json:"reason"`
}
