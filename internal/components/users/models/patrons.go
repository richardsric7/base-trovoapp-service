package users

import "time"

type PatronPackage struct {
	ID            string `json:"id"`
	PackageName   string `json:"packageName"`
	Inactive      int    `gorm:"default:0" json:"inactive"`
	PriorityOrder int    `gorm:"default:1" json:"-"`
}

type PatronTier struct {
	ID            string `json:"id"`
	Tier          string `json:"tier"`
	CanExpire     int    `gorm:"default:1" json:"canExpire"`
	Inactive      int    `gorm:"default:0" json:"inactive"`
	PriorityOrder int    `gorm:"default:1" json:"-"`
}

type UserPatronMembership struct {
	Username        string        `gorm:"primaryKey" json:"username"`
	PatronPackageID string        `json:"subscriptionPackageId"`
	PatronPackage   PatronPackage `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	PatronTierID    string        `json:"subscriptionTierId"`
	PatronTier      PatronTier    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	ValidTill       time.Time     `gorm:"not null" json:"validTill"` //lifetime is represented by 'infinity'
}

type UserPatronSubscriptionLog struct {
	ID              string    `json:"id"`
	CreatedAt       time.Time `json:"createdAt"`
	Username        string    `json:"username"`
	PatronPackageID string    `json:"subscriptionPackageId"`
	PatronTierID    string    `json:"subscriptionTierId"`
	ValidTill       time.Time `gorm:"not null" json:"validTill"` //lifetime is represented by 'infinity'
}
