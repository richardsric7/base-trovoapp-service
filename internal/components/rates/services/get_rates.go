package rates

import (
	ratesDB "trovo-wallet-api/internal/components/rates/db"
	ratesModel "trovo-wallet-api/internal/components/rates/models"

	"gorm.io/gorm"
)

// HandleGetRates is service to get rate
func HandleGetRates(db *gorm.DB) (rate ratesModel.Rate) {

	return ratesDB.GetRates(db)

}
