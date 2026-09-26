package models

// CurrencyConfig represents the model for currency configurations in the database
type CurrencyConfig struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Code        string `gorm:"uniqueIndex;not null" json:"code"`
	Description string `json:"description"`
}
