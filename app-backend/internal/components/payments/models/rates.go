package payments

import "time"

// Rates model
type CurrencyRates struct {
	ID        uint64    `json:"id"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	Rates     []byte    `gorm:"null" json:"rates"`
}
