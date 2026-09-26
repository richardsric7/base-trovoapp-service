package p2p

import (
	"time"
	p2pModels "trovo-wallet-api/internal/components/p2p/models"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// completionRateString computes completed / (completed + cancelled +
// expired) as a percentage string. Intentionally NOT populated here:
// MedianOrderCompletionTimeSeconds/MedianDisputeResolutionTimeSeconds (a
// true median needs the raw historical samples, not just a running
// average - approximated as equal to the average field instead, see
// runningAverageSeconds) and FalsePaymentClaims/IncompletePaymentClaims/
// UnderpaymentEvents (no unambiguous mapping from a dispute's
// subject+resolution to these three distinct counters was specified
// anywhere in the plan - left at zero rather than guessed). Lifetime/
// Last30Day/Last90DayMetrics are - per this model's own existing doc
// comment - "computed on read/refresh", so they're deliberately not
// written here either.
func completionRateString(completed, cancelled, expired int64) string {
	total := completed + cancelled + expired
	if total == 0 {
		return "0"
	}
	rate := decimal.NewFromInt(completed).Mul(decimal.NewFromInt(100)).Div(decimal.NewFromInt(total))
	return rate.Truncate(2).String()
}

// runningAverageSeconds folds a new sample into a running mean without
// storing the historical sample list.
func runningAverageSeconds(currentAverage int64, sampleCount int64, newSample time.Duration) int64 {
	if sampleCount <= 0 {
		return int64(newSample.Seconds())
	}
	total := currentAverage*sampleCount + int64(newSample.Seconds())
	return total / (sampleCount + 1)
}

// UpdatePerformanceOnCompletion updates MerchantPerformance,
// MerchantOfferPerformance, and CustomerPerformance rolling aggregates
// (Plan Section 15) after an order completes. A merchant rejection is not
// automatically misconduct (Section 15's own note) - only real outcomes
// (completions, cancellations, expiries, disputes) move these numbers.
func UpdatePerformanceOnCompletion(gc *sharedconfig.GlobalConfig, order p2pModels.Order) {
	now := time.Now().UTC()
	volume := decimal.RequireFromString(order.SpecifiedAssetAmount)
	completionTime := now.Sub(order.CreatedAt)

	upsertMerchantPerformance(gc.DB, order.MerchantUserID, order.MerchantUsername, func(p *p2pModels.MerchantPerformance) {
		p.AverageOrderCompletionTimeSeconds = runningAverageSeconds(p.AverageOrderCompletionTimeSeconds, p.CompletedTrades, completionTime)
		p.MedianOrderCompletionTimeSeconds = p.AverageOrderCompletionTimeSeconds
		p.CompletedTrades++
		p.CompletedTradeVolume = decimal.RequireFromString(orDefaultStr(p.CompletedTradeVolume, "0")).Add(volume).String()
		p.CompletionRate = completionRateString(p.CompletedTrades, p.CancelledOrders, p.ExpiredOrders)
		p.LastActivityAt = &now
	})

	upsertMerchantOfferPerformance(gc.DB, order.MerchantUserID, order.OfferID, func(p *p2pModels.MerchantOfferPerformance) {
		p.MedianCompletionTimeSeconds = runningAverageSeconds(p.MedianCompletionTimeSeconds, p.CompletedTrades, completionTime)
		p.CompletedTrades++
		p.CompletionRate = completionRateString(p.CompletedTrades, 0, 0)
	})

	upsertCustomerPerformance(gc.DB, order.CustomerUserID, order.CustomerUsername, func(p *p2pModels.CustomerPerformance) {
		p.AverageOrderCompletionTimeSeconds = runningAverageSeconds(p.AverageOrderCompletionTimeSeconds, p.CompletedTrades, completionTime)
		p.MedianOrderCompletionTimeSeconds = p.AverageOrderCompletionTimeSeconds
		if order.PaymentSentAt != nil {
			confirmationTime := now.Sub(*order.PaymentSentAt)
			p.PaymentConfirmationTimeSeconds = runningAverageSeconds(p.PaymentConfirmationTimeSeconds, p.CompletedTrades, confirmationTime)
		}
		p.CompletedTrades++
		p.CompletedTradeVolume = decimal.RequireFromString(orDefaultStr(p.CompletedTradeVolume, "0")).Add(volume).String()
		p.CompletionRate = completionRateString(p.CompletedTrades, p.CancelledOrders, p.ExpiredOrders)
		p.LastActivityAt = &now
	})
}

// RecordOrderCancelled updates CancelledOrders (and, for the customer,
// CancelledAfterAcceptance when the merchant had already accepted) for
// whichever party initiated the cancellation. Plan Section 15's own note
// that a merchant rejection is not misconduct means RejectOrder
// (AWAITING_APPROVAL) intentionally does not call this - only an actual
// cancellation after the order progressed does.
func RecordOrderCancelled(gc *sharedconfig.GlobalConfig, order p2pModels.Order, cancelledByMerchant bool) {
	now := time.Now().UTC()
	if cancelledByMerchant {
		upsertMerchantPerformance(gc.DB, order.MerchantUserID, order.MerchantUsername, func(p *p2pModels.MerchantPerformance) {
			p.CancelledOrders++
			p.CompletionRate = completionRateString(p.CompletedTrades, p.CancelledOrders, p.ExpiredOrders)
			p.LastActivityAt = &now
		})
		return
	}
	upsertCustomerPerformance(gc.DB, order.CustomerUserID, order.CustomerUsername, func(p *p2pModels.CustomerPerformance) {
		p.CancelledOrders++
		p.CancelledAfterAcceptance++ // CancelOrder/merchant-cancel are only reachable post-AWAITING_APPROVAL
		p.CompletionRate = completionRateString(p.CompletedTrades, p.CancelledOrders, p.ExpiredOrders)
		p.LastActivityAt = &now
	})
}

// RecordOrderExpired updates ExpiredOrders for both parties - used for all
// three of ExpireStaleOrders' timeout paths, including the AWAITING_PAYMENT
// one that ends in OrderStatusCancelled (the underlying event was still a
// timeout, not a deliberate cancellation by either party).
func RecordOrderExpired(gc *sharedconfig.GlobalConfig, order p2pModels.Order) {
	now := time.Now().UTC()
	upsertMerchantPerformance(gc.DB, order.MerchantUserID, order.MerchantUsername, func(p *p2pModels.MerchantPerformance) {
		p.ExpiredOrders++
		p.CompletionRate = completionRateString(p.CompletedTrades, p.CancelledOrders, p.ExpiredOrders)
		p.LastActivityAt = &now
	})
	upsertCustomerPerformance(gc.DB, order.CustomerUserID, order.CustomerUsername, func(p *p2pModels.CustomerPerformance) {
		p.ExpiredOrders++
		p.CompletionRate = completionRateString(p.CompletedTrades, p.CancelledOrders, p.ExpiredOrders)
		p.LastActivityAt = &now
	})
}

// RecordDisputeOpened updates DisputesOpened for both parties to the order,
// regardless of which of them opened it - being party to a dispute is the
// signal being tracked, not fault.
func RecordDisputeOpened(gc *sharedconfig.GlobalConfig, order p2pModels.Order) {
	now := time.Now().UTC()
	upsertMerchantPerformance(gc.DB, order.MerchantUserID, order.MerchantUsername, func(p *p2pModels.MerchantPerformance) {
		p.DisputesOpened++
		p.LastActivityAt = &now
	})
	upsertCustomerPerformance(gc.DB, order.CustomerUserID, order.CustomerUsername, func(p *p2pModels.CustomerPerformance) {
		p.DisputesOpened++
		p.LastActivityAt = &now
	})
}

// RecordDisputeResolved updates DisputesResolved and the
// against/in-favor-of split for both parties, plus the running average
// dispute-resolution time, based on the dispute's outcome.
// DisputeResolutionSplit updates DisputesResolved only - Section 62 itself
// notes the plan doesn't define an automated split-settlement outcome, so
// there is no well-defined "favor" to attribute here either.
func RecordDisputeResolved(gc *sharedconfig.GlobalConfig, order p2pModels.Order, resolution string, resolutionTime time.Duration) {
	now := time.Now().UTC()
	upsertMerchantPerformance(gc.DB, order.MerchantUserID, order.MerchantUsername, func(p *p2pModels.MerchantPerformance) {
		p.AverageDisputeResolutionTimeSeconds = runningAverageSeconds(p.AverageDisputeResolutionTimeSeconds, p.DisputesResolved, resolutionTime)
		p.MedianDisputeResolutionTimeSeconds = p.AverageDisputeResolutionTimeSeconds
		p.DisputesResolved++
		switch resolution {
		case p2pModels.DisputeResolutionInFavorOfSeller:
			p.DisputesResolvedInFavorOfMerchant++
		case p2pModels.DisputeResolutionInFavorOfBuyer:
			p.DisputesResolvedAgainstMerchant++
		}
		p.LastActivityAt = &now
	})
	upsertCustomerPerformance(gc.DB, order.CustomerUserID, order.CustomerUsername, func(p *p2pModels.CustomerPerformance) {
		p.DisputesResolved++
		switch resolution {
		case p2pModels.DisputeResolutionInFavorOfBuyer:
			p.DisputesResolvedInFavorOfCustomer++
		case p2pModels.DisputeResolutionInFavorOfSeller:
			p.DisputesResolvedAgainstCustomer++
		}
		p.LastActivityAt = &now
	})
	upsertMerchantOfferPerformance(gc.DB, order.MerchantUserID, order.OfferID, func(p *p2pModels.MerchantOfferPerformance) {
		p.DisputesResolved++
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

// GetMerchantPerformance and GetCustomerPerformance back the marketplace
// display of trust signals (Plan Section 8/26 - "merchant sees customer
// performance"; marketplace offer cards show merchant performance).
func GetMerchantPerformance(db *gorm.DB, merchantID string) (p2pModels.MerchantPerformance, error) {
	var p p2pModels.MerchantPerformance
	err := db.Where("merchant_id = ?", merchantID).First(&p).Error
	if err == gorm.ErrRecordNotFound {
		return p2pModels.MerchantPerformance{MerchantID: merchantID, CompletionRate: "0"}, nil
	}
	return p, err
}

func GetCustomerPerformance(db *gorm.DB, customerID string) (p2pModels.CustomerPerformance, error) {
	var p p2pModels.CustomerPerformance
	err := db.Where("customer_id = ?", customerID).First(&p).Error
	if err == gorm.ErrRecordNotFound {
		return p2pModels.CustomerPerformance{CustomerID: customerID, CompletionRate: "0"}, nil
	}
	return p, err
}
