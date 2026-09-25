package models

import (
	p2pErrors "admin-panel-dashboard/internal/errors"
	"errors"
	"log"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CurrencyCast string

func (c CurrencyCast) GetCurrencies(db *gorm.DB) (currencies []Currency, err error) {

	e := db.Preload(clause.Associations).Order("id ASC").Find(&currencies).Error
	if e != nil {
		if !errors.Is(e, gorm.ErrRecordNotFound) {
			log.Printf("[GetCurrencies]critical error occurred while checking for all GetCurrencies, err:[%v]\n", e)
			return currencies, &p2pErrors.ErrorTemporaryServerError{}
		}
	}
	return currencies, nil

}
