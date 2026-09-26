package p2p

import (
	"time"
	p2pModels "trovo-wallet-api/internal/components/p2p/models"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// UpdatePerformanceOnCompletion updates MerchantPerformance,
// MerchantOfferPerformance, and CustomerPerformance rolling aggregates
// (Plan Section 15) after an order completes. A merchant rejection is not
// automatically misconduct (Section 15's own note) - only real outcomes
// (completions, cancellations, disputes) move these numbers.
func UpdatePerformanceOnCompletion(gc *sharedconfig.GlobalConfig, order p2pModels.Order) {
	now := time.Now().UTC()
	volume := decimal.RequireFromString(order.SpecifiedAssetAmount)

	upsertMerchantPerformance(gc.DB, order.MerchantUserID, order.MerchantUsername, func(p *p2pModels.MerchantPerformance) {
		p.CompletedTrades++
		p.CompletedTradeVolume = decimal.RequireFromString(orDefaultStr(p.CompletedTradeVolume, "0")).Add(volume).String()
		p.LastActivityAt = &now
	})

	upsertMerchantOfferPerformance(gc.DB, order.MerchantUserID, order.OfferID, func(p *p2pModels.MerchantOfferPerformance) {
		p.CompletedTrades++
	})

	upsertCustomerPerformance(gc.DB, order.CustomerUserID, order.CustomerUsername, func(p *p2pModels.CustomerPerformance) {
		p.CompletedTrades++
		p.CompletedTradeVolume = decimal.RequireFromString(orDefaultStr(p.CompletedTradeVolume, "0")).Add(volume).String()
		p.LastActivityAt = &now
	})
}

func upsertMerchantPerformance(db *gorm.DB, merchantID, username string, mutate func(*p2pModels.MerchantPerformance)) {
	var p p2pModels.MerchantPerformance
	err := db.Where("merchant_id = ?", merchantID).First(&p).Error
	if err == gorm.ErrRecordNotFound {
		p = p2pModels.MerchantPerformance{ID: uuid.NewString(), MerchantID: merchantID, MerchantUsername: username}
	}
	mutate(&p)
	db.Save(&p)
}

func upsertMerchantOfferPerformance(db *gorm.DB, merchantID, offerID string, mutate func(*p2pModels.MerchantOfferPerformance)) {
	var p p2pModels.MerchantOfferPerformance
	err := db.Where("offer_id = ?", offerID).First(&p).Error
	if err == gorm.ErrRecordNotFound {
		p = p2pModels.MerchantOfferPerformance{ID: uuid.NewString(), MerchantID: merchantID, OfferID: offerID}
	}
	mutate(&p)
	db.Save(&p)
}

func upsertCustomerPerformance(db *gorm.DB, customerID, username string, mutate func(*p2pModels.CustomerPerformance)) {
	var p p2pModels.CustomerPerformance
	err := db.Where("customer_id = ?", customerID).First(&p).Error
	if err == gorm.ErrRecordNotFound {
		p = p2pModels.CustomerPerformance{ID: uuid.NewString(), CustomerID: customerID, CustomerUsername: username}
	}
	mutate(&p)
	db.Save(&p)
}
