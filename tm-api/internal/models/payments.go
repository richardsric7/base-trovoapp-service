package models

import (
	"time"

	"gorm.io/gorm"
)

type PaymentMethod struct {
	ID                   string         `gorm:"size:100;primaryKey;check:,length(id) > 5" json:"id"`
	Username             string         `gorm:"not null;index:idx_user_payment_method,unique;check:,length(username) < 50" json:"username"`
	CreatedAt            time.Time      `gorm:"default:now()" json:"createdAt"`
	UpdatedAt            time.Time      `gorm:"default:now()" json:"updatedAt"`
	PaymentChannel       PaymentChannel `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	PaymentChannelID     string         `gorm:"not null;index:idx_user_payment_method,unique;check:,length(payment_channel_id) > 2" json:"paymentChannelId"`
	Name                 *string        `gorm:"size:100;null;" json:"name"`
	DestinationAccount   string         `gorm:"size:150;not null;index:idx_user_payment_method,unique" json:"destinationAccount"`
	Memo                 *string        `gorm:"size:100;null;" json:"memo"`
	BankName             *string        `gorm:"size:100;null;" json:"bankName"`
	AccountOpeningBranch *string        `gorm:"size:100;null;" json:"accountOpeningBranch"`
	CountryCode          *string        `gorm:"size:3;null;" json:"countryCode"` // country codes are 2 characters and max of 3 characters depending ons tandard.
	CurrencyID           *string        `gorm:"size:10;null;" json:"currencyId"` // CurrencyID for both MoMo and Bank Transfer
}
type PaymentCallback struct {
	TransactionID   string    `gorm:"size:64;primaryKey;check:,length(transaction_id) > 5"`
	CreatedAt       time.Time `gorm:"not null;default:now()" json:"createdAt"`
	Memo            string    `gorm:"not null;size:28;null;index:idx_callback_memo;check:,length(memo) > 20" json:"memo"`
	Sender          string    `gorm:"size 60" json:"sender"`
	Amount          string    `gorm:"size:100" json:"amount"`
	AssetCode       string    `gorm:"size:12" json:"assetCode"`
	ContractAddress string    `gorm:"size:60" json:"contractAddress"`
}
type PaymentMethodJSON struct {
	CreatedAt            time.Time `json:"createdAt"`
	UpdatedAt            time.Time `json:"updatedAt"`
	ID                   string    `gorm:"size:100;primaryKey" json:"id"`
	Username             string    `gorm:"not null;" json:"username"`
	PaymentChannelID     string    `gorm:"not null;" json:"paymentChannelId"`
	PaymentType          string    `gorm:"not null;" json:"paymentType"`
	Name                 string    `gorm:"size:100;null;" json:"name,omitempty"`
	DestinationAccount   string    `gorm:"size:150;not null;" json:"destinationAccount"`
	Memo                 string    `gorm:"size:100;null;" json:"memo,omitempty"`
	BankName             string    `gorm:"size:100;null;" json:"bankName,omitempty"`
	AccountOpeningBranch string    `gorm:"size:100;null;" json:"accountOpeningBranch,omitempty"`
	CountryCode          string    `gorm:"size:100;null;" json:"countryCode,omitempty"`
	CurrencyID           string    `gorm:"size:10;null;" json:"currencyId,omitempty"` // CurrencyID for both MoMo and Bank Transfer
}

type MarketPrice struct {
	Pair           string `json:"paire"`
	Price          string `json:"price"`
	CurrencySymbol string `json:"currencySymbol"`
}
type PaymentCallbackInfo struct {
	Destination     string    `json:"destination"`
	Sender          string    `json:"sender"`
	Amount          string    `json:"amount"`
	AssetCode       string    `json:"assetCode"`
	ContractAddress string    `json:"contractAddress"`
	TransactionID   string    `json:"transactionId"`
	TransactionMemo string    `json:"transactionMemo"`
	TransactionTime time.Time `json:"transactionTime"`
}
type CurrencyPaymentMethod string

func (pmts *PaymentMethod) ToPaymentMethodJSON() (paymentMethod PaymentMethodJSON) {

	paymentMethod = PaymentMethodJSON{
		ID:                 pmts.ID,
		CreatedAt:          pmts.CreatedAt,
		UpdatedAt:          pmts.UpdatedAt,
		Username:           pmts.Username,
		PaymentChannelID:   pmts.PaymentChannelID,
		DestinationAccount: pmts.DestinationAccount,
		PaymentType:        pmts.PaymentChannel.PaymentTypeID,
	}
	if pmts.AccountOpeningBranch != nil {
		paymentMethod.AccountOpeningBranch = *pmts.AccountOpeningBranch
	}
	if pmts.Name != nil {
		paymentMethod.Name = *pmts.Name
	}
	if pmts.BankName != nil {
		paymentMethod.BankName = *pmts.BankName
	}
	if pmts.Memo != nil {
		paymentMethod.Memo = *pmts.Memo
	}
	if pmts.CountryCode != nil {
		paymentMethod.CountryCode = *pmts.CountryCode
	}
	if pmts.CurrencyID != nil {
		paymentMethod.CurrencyID = *pmts.CurrencyID
	}

	return paymentMethod
}
func (pID CurrencyPaymentMethod) ToPaymentMethodJSON(db *gorm.DB) (paymentMethod PaymentMethodJSON) {
	var pmts PaymentMethod
	e := db.Where("id = ?", string(pID)).First(&pmts).Error
	if e != nil {
		return PaymentMethodJSON{}
	}
	paymentMethod = PaymentMethodJSON{
		ID:                 pmts.ID,
		CreatedAt:          pmts.CreatedAt,
		UpdatedAt:          pmts.UpdatedAt,
		Username:           pmts.Username,
		PaymentChannelID:   pmts.PaymentChannelID,
		DestinationAccount: pmts.DestinationAccount,
		PaymentType:        pmts.PaymentChannel.PaymentTypeID,
	}
	if pmts.AccountOpeningBranch != nil {
		paymentMethod.AccountOpeningBranch = *pmts.AccountOpeningBranch
	}
	if pmts.Name != nil {
		paymentMethod.Name = *pmts.Name
	}
	if pmts.BankName != nil {
		paymentMethod.BankName = *pmts.BankName
	}
	if pmts.Memo != nil {
		paymentMethod.Memo = *pmts.Memo
	}
	if pmts.CountryCode != nil {
		paymentMethod.CountryCode = *pmts.CountryCode
	}

	if pmts.CurrencyID != nil {
		paymentMethod.CurrencyID = *pmts.CurrencyID
	}
	return paymentMethod
}
