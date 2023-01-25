package users

import "time"

type PatronPackage struct {
	ID            string `json:"id"`
	Decscription  string `json:"decscription"`
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
	PatronPackageID string        `json:"patronPackageId"`
	PatronPackage   PatronPackage `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	PatronTierID    string        `json:"patronTierId"`
	PatronTier      PatronTier    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	ValidTill       time.Time     `gorm:"not null" json:"validTill"` //lifetime is represented by 'infinity'
}

type PatronMembershipPrice struct {
	ID            uint64 `gorm:"primaryKey" json:"id"`
	PatronPackage string `json:"patronPackage"`
	PatronTierID  string `json:"patronTier"`
}

type UserPatronSubscriptionLog struct {
	ID                    string    `json:"id"`
	CreatedAt             time.Time `json:"createdAt"`
	UpdatedAt             time.Time `json:"updatedAt"`
	Username              string    `json:"username"`
	PatronPackageID       string    `json:"patronPackageId"`
	PatronTierID          string    `json:"patronTierId"`
	ActivePatronPackageID *string   `gorm:"null" json:"activePatronPackageId"` //valid and used only when effectiveDate is future
	ActivePatronTierID    *string   `gorm:"null" json:"activePatronTierId"`    //valid and used only when effectiveDate is future
	EffectiveDate         time.Time `gorm:"not null" json:"effectiveDate"`     //Holds when this subscription becomes effective. lifetime is represented by 'infinity'
	ValidTill             time.Time `gorm:"not null" json:"validTill"`         //lifetime is represented by 'infinity'
}

type PatronSubscriptionInput struct {
	PatronMembershipPriceID uint64   `json:"patronMembershipPriceId"`
	Transaction             string   `json:"transaction"`
	TransactionSignature    string   `json:"transactionSignature"`
	TransactionID           string   `json:"transactionId"`
	NetworkPassPhrase       string   `json:"networkPassPhrase"`
	Messages                []string `json:"messages"`
}
