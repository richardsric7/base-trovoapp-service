package p2p

import "time"

// MerchantPaymentMethod is a merchant's own saved fiat settlement channel,
// selectable from SELL offers by reference (Offer.PaymentMethodID) rather
// than a bespoke value typed into each offer. It is never hard-deleted -
// only deactivated (IsActive) - and cannot be deactivated while any live
// offer still references it (see services.SetPaymentMethodActive): the
// merchant must close the offer or edit it onto a different payment
// method first.
//
// Editing this record (services.UpdatePaymentMethod) propagates to every
// offer currently referencing it, so a live offer always shows its
// merchant's current settlement details - but never to an
// already-created Order.PaymentMethodSnapshot, which is a copy frozen at
// order-creation time and immune to later edits here.
type MerchantPaymentMethod struct {
	ID               string    `json:"id" gorm:"primaryKey;size:36"`
	MerchantUserID   string    `json:"merchantUserId" gorm:"size:100;not null;index:idx_p2p_payment_method_merchant"`
	MerchantUsername string    `json:"merchantUsername" gorm:"size:70;not null"`
	PaymentChannel   string    `json:"paymentChannel" gorm:"size:30;not null;default:''"`
	Provider         string    `json:"provider" gorm:"size:60;not null;default:''"`
	Account          string    `json:"account" gorm:"size:100;not null;default:''"`
	IsActive         bool      `json:"isActive" gorm:"not null;default:true"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

// ToPaymentMethod converts to the denormalized display value embedded on
// an Offer/Order (see offer.go's PaymentMethod).
func (m MerchantPaymentMethod) ToPaymentMethod() PaymentMethod {
	return PaymentMethod{PaymentChannel: m.PaymentChannel, Provider: m.Provider, Account: m.Account}
}
