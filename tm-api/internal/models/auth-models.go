package models

type LoginInput struct {
	Username  string `json:"username"`
	LoginType string `json:"loginType"`
}
type LoginCallbackInput struct {
	LoginID    string `json:"loginId"`
	TargetUser string `json:"targetUser"`
}
type AuthorizationCallbackInput struct {
	AuthID     string `json:"authId"`
	TargetUser string `json:"targetUser"`
}
type PendingAuthorization struct {
	AuthID         string `gorm:"size:100;primaryKey" json:"authId"`
	TargetUser     string `gorm:"size:100" json:"targetUser"`
	TargetID       string `gorm:"size:100" json:"targetId"`
	TargetCategory string `gorm:"size:100" json:"targetCategory"`
	Status         uint   `gorm:"default:0" json:"status"`
}
