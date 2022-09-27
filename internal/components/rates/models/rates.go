package rates

import (
	"time"
)

// Rates model
type CurrencyRates struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	Rates     []byte    `gorm:"null" json:"rates"`
}

// FetchedRates sotres fetched rates
type FetchedRates struct {
	Disclaimer string                 `json:"disclaimer"`
	License    string                 `json:"license"`
	Timestamp  uint64                 `json:"timestamp"`
	Base       string                 `json:"base"`
	Rates      map[string]interface{} `json:"rates"`
}

// Rate stores slices rate
type Rate map[string]interface{}
