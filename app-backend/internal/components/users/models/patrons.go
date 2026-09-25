package users

import "time"

type PatronPackage struct {
	ID               string `json:"id"`
	Description      string `json:"description"`
	PackageListTitle string `json:"packageListTitle"`
	PackageList      string `json:"packageList"`
	Inactive         int    `gorm:"default:0" json:"inactive"`
	PriorityOrder    int    `gorm:"default:1" json:"-"`
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
	ValidTill       time.Time     `gorm:"not null" json:"validTill"` //lifetime is represented by year '9999'

}

type PatronMembershipGrade struct {
	ID            uint64  `gorm:"primaryKey" json:"id"`
	PatronPackage string  `json:"patronPackage"`
	PatronTierID  string  `json:"patronTier"`
	Price         float64 `json:"price"`
}

type UserPatronSubscriptionLog struct {
	ID                    string    `json:"id"`
	CreatedAt             time.Time `json:"createdAt"`
	UpdatedAt             time.Time `json:"updatedAt"`
	Username              string    `json:"username"`
	PatronPackageID       string    `json:"patronPackageId"`
	PatronTierID          string    `json:"patronTierId"`
	ActivePatronPackageID *string   `gorm:"null" json:"activePatronPackageId"`    //valid and used only when effectiveDate is future
	ActivePatronTierID    *string   `gorm:"null" json:"activePatronTierId"`       //valid and used only when effectiveDate is future
	EffectiveDate         time.Time `gorm:"not null" json:"effectiveDate"`        //Holds when this subscription becomes effective.
	ValidTill             time.Time `gorm:"not null" json:"validTill"`            //lifetime is represented by year '9999'
	VatPaid               float64   `gorm:"not null;default:0.00" json:"vatPaid"` //vat paid

}

type PatronSubscriptionInput struct {
	PatronMembershipGradeID uint64   `json:"patronMembershipGradeId"`
	PaymentAssetCode        string   `json:"paymentAssetCode"`
	PaymentContractAddress  string   `json:"paymentContractAddress"`
	Transaction             string   `json:"transaction"`
	TransactionSignature    string   `json:"transactionSignature"`
	TransactionID           string   `json:"transactionId"`
	NetworkPassPhrase       string   `json:"networkPassPhrase"`
	Vat                     string   `json:"vat"`
	VatAmount               string   `json:"vatAmount"`
	AmountToPay             string   `json:"amountToPay"`
	Messages                []string `json:"messages"`
}

type PatronSubscriptionPaymentAsset struct {
	AssetCode       string `gorm:"primaryKey" json:"assetCode"`
	ContractAddress string `json:"contractAddress"`
	Inactive        uint   `gorm:"default:0" json:"-"`
}
