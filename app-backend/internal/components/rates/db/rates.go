package rates

import (
	"encoding/json"
	ratesModel "trovo-wallet-api/internal/components/rates/models"

	"gorm.io/gorm"
)

// GetRates is service to get rate
func GetRates(db *gorm.DB) (rate ratesModel.Rate) {

	rates := &ratesModel.CurrencyRates{}
	db.First(&rates)
	json.Unmarshal(rates.Rates, &rate)
	return
}
