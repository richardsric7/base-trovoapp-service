package models

import (
	"time"

	"github.com/shopspring/decimal"
)

const (
	// FiatRecordTypePayment represents rows sourced from fiat_payment table.
	FiatRecordTypePayment = "payment"
	// FiatRecordTypeInvoice represents rows sourced from fiat_payment_invoice table.
	FiatRecordTypeInvoice = "invoice"
)

// FiatPaymentInvoice represents a fiat payment invoice table record.
type FiatPaymentInvoice struct {
	ID              string          `gorm:"column:id;primaryKey" json:"id"`
	CreatedAt       time.Time       `gorm:"column:created_at" json:"created_at"`
	ServiceProvider string          `gorm:"column:service_provider" json:"service_provider"`
	Username        string          `gorm:"column:username" json:"username"`
	Amount          decimal.Decimal `gorm:"column:amount" json:"amount"`
	PaymentType     string          `gorm:"column:payment_type" json:"payment_type"`
	Status          string          `gorm:"column:status" json:"status"`
	Refunded        int64           `gorm:"column:refunded" json:"refunded"`
}

// TableName overrides the default gorm table name.
func (FiatPaymentInvoice) TableName() string {
	return "fiat_payment_invoices"
}

// FiatPayment represents a fiat payment table record.
type FiatPayment struct {
	ID              int64           `gorm:"column:id;primaryKey" json:"id"`
	ServiceProvider string          `gorm:"column:service_provider" json:"service_provider"`
	Username        string          `gorm:"column:username" json:"username"`
	TransactionID   string          `gorm:"column:transaction_id" json:"transaction_id"`
	Amount          decimal.Decimal `gorm:"column:amount" json:"amount"`
	PaymentType     string          `gorm:"column:payment_type" json:"payment_type"`
	CreatedAt       time.Time       `gorm:"column:created_at" json:"created_at"`
}

// TableName overrides the default gorm table name.
func (FiatPayment) TableName() string {
	return "fiat_payments"
}

// FiatPaymentRecord represents a flattened view of fiat payment data from both tables.
type FiatPaymentRecord struct {
	ID              string          `json:"id" example:"12345"`
	RecordType      string          `json:"record_type" example:"payment" enums:"payment,invoice"`
	ServiceProvider string          `json:"service_provider" example:"flutterwave"`
	Username        string          `json:"username" example:"john.doe"`
	TransactionID   *string         `json:"transaction_id,omitempty" example:"TXN-001231"`
	Amount          decimal.Decimal `json:"amount" example:"2500.50"`
	PaymentType     string          `json:"payment_type" example:"bank_transfer"`
	Status          *string         `json:"status,omitempty" example:"completed"`
	Refunded        *int64          `json:"refunded,omitempty" example:"0"`
	CreatedAt       time.Time       `json:"created_at" example:"2024-03-20T10:00:00Z"`
}

// FiatPaymentListFilter encapsulates filtering and pagination options.
type FiatPaymentListFilter struct {
	RecordType      string
	ServiceProvider string
	Username        string
	PaymentType     string
	Status          string
	Page            int
	PageSize        int
}

// PaginatedFiatPaymentResponse represents the paginated response payload.
type PaginatedFiatPaymentResponse struct {
	Data       []FiatPaymentRecord `json:"data"`
	Total      int64               `json:"total" example:"25"`
	Page       int                 `json:"page" example:"1"`
	PageSize   int                 `json:"page_size" example:"10"`
	TotalPages int                 `json:"total_pages" example:"3"`
}
