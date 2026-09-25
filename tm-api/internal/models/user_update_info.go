package models

// UserUpdateInfo model for user update info
type UserUpdateInfo struct {
	ContactPhone string `json:"contactPhone"`
}
type KYCUpdateInfo struct {
	KYCLevel uint `json:"kycLevel"`
}
