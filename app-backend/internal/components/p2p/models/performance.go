package p2p

import "time"

// MerchantPerformance is a rolling aggregate of a merchant's trading history,
// used by the marketplace ranking engine and shown on offer/order screens.
type MerchantPerformance struct {
	ID                                 string    `json:"id" gorm:"primaryKey;size:36"`
	MerchantID                         string    `json:"merchantId" gorm:"size:100;not null;uniqueIndex:idx_p2p_merchant_perf_merchant_id"`
	MerchantUsername                   string    `json:"merchantUsername" gorm:"size:70;not null"`
	CompletedTrades                    int64     `json:"completedTrades" gorm:"not null;default:0"`
	CompletedTradeVolume               string    `json:"completedTradeVolume" gorm:"size:60;not null;default:'0'"`
	CompletionRate                     string    `json:"completionRate" gorm:"size:10;not null;default:'0'"`
	AverageOrderCompletionTimeSeconds  int64     `json:"averageOrderCompletionTime" gorm:"not null;default:0"`
	MedianOrderCompletionTimeSeconds   int64     `json:"medianOrderCompletionTime" gorm:"not null;default:0"`
	PaymentConfirmationTimeSeconds     int64     `json:"paymentConfirmationTime" gorm:"not null;default:0"`
	CancelledOrders                    int64     `json:"cancelledOrders" gorm:"not null;default:0"`
	ExpiredOrders                      int64     `json:"expiredOrders" gorm:"not null;default:0"`
	DisputesOpened                     int64     `json:"disputesOpened" gorm:"not null;default:0"`
	DisputesResolved                   int64     `json:"disputesResolved" gorm:"not null;default:0"`
	DisputesResolvedAgainstMerchant    int64     `json:"disputesResolvedAgainstMerchant" gorm:"not null;default:0"`
	DisputesResolvedInFavorOfMerchant  int64     `json:"disputesResolvedInFavorOfMerchant" gorm:"not null;default:0"`
	AverageDisputeResolutionTimeSeconds int64    `json:"averageDisputeResolutionTime" gorm:"not null;default:0"`
	MedianDisputeResolutionTimeSeconds int64     `json:"medianDisputeResolutionTime" gorm:"not null;default:0"`
	LastActivityAt                     *time.Time `json:"lastActivityAt"`
	LifetimeMetrics                    string    `json:"lifetimeMetrics" gorm:"type:text"`  // JSON blob, computed on read/refresh
	Last30DayMetrics                   string    `json:"last30DayMetrics" gorm:"type:text"` // JSON blob
	Last90DayMetrics                   string    `json:"last90DayMetrics" gorm:"type:text"` // JSON blob
	CreatedAt                          time.Time `json:"createdAt"`
	UpdatedAt                          time.Time `json:"updatedAt"`
}

// MerchantOfferPerformance is per-offer performance, referenced by Offer.ID.
type MerchantOfferPerformance struct {
	ID                          string    `json:"id" gorm:"primaryKey;size:36"`
	MerchantID                 string    `json:"merchantId" gorm:"size:100;not null;index:idx_p2p_offer_perf_merchant_id"`
	OfferID                    string    `json:"offerId" gorm:"size:36;not null;uniqueIndex:idx_p2p_offer_perf_offer_id"`
	CompletedTrades            int64     `json:"completedTrades" gorm:"not null;default:0"`
	CompletionRate             string    `json:"completionRate" gorm:"size:10;not null;default:'0'"`
	MedianCompletionTimeSeconds int64    `json:"medianCompletionTime" gorm:"not null;default:0"`
	DisputesResolved           int64     `json:"disputesResolved" gorm:"not null;default:0"`
	CreatedAt                  time.Time `json:"createdAt"`
	UpdatedAt                  time.Time `json:"updatedAt"`
}

// CustomerPerformance is a rolling aggregate of a customer's trading history.
type CustomerPerformance struct {
	ID                                 string     `json:"id" gorm:"primaryKey;size:36"`
	CustomerID                        string     `json:"customerId" gorm:"size:100;not null;uniqueIndex:idx_p2p_customer_perf_customer_id"`
	CustomerUsername                  string     `json:"customerUsername" gorm:"size:70;not null"`
	CompletedTrades                   int64      `json:"completedTrades" gorm:"not null;default:0"`
	CompletedTradeVolume              string     `json:"completedTradeVolume" gorm:"size:60;not null;default:'0'"`
	CompletionRate                    string     `json:"completionRate" gorm:"size:10;not null;default:'0'"`
	AverageOrderCompletionTimeSeconds int64      `json:"averageOrderCompletionTime" gorm:"not null;default:0"`
	MedianOrderCompletionTimeSeconds  int64      `json:"medianOrderCompletionTime" gorm:"not null;default:0"`
	CancelledOrders                   int64      `json:"cancelledOrders" gorm:"not null;default:0"`
	CancelledAfterAcceptance          int64      `json:"cancelledAfterAcceptance" gorm:"not null;default:0"`
	ExpiredOrders                     int64      `json:"expiredOrders" gorm:"not null;default:0"`
	DisputesOpened                    int64      `json:"disputesOpened" gorm:"not null;default:0"`
	DisputesResolved                  int64      `json:"disputesResolved" gorm:"not null;default:0"`
	DisputesResolvedAgainstCustomer   int64      `json:"disputesResolvedAgainstCustomer" gorm:"not null;default:0"`
	DisputesResolvedInFavorOfCustomer int64      `json:"disputesResolvedInFavorOfCustomer" gorm:"not null;default:0"`
	PaymentConfirmationTimeSeconds    int64      `json:"paymentConfirmationTime" gorm:"not null;default:0"`
	FalsePaymentClaims                int64      `json:"falsePaymentClaims" gorm:"not null;default:0"`
	IncompletePaymentClaims           int64      `json:"incompletePaymentClaims" gorm:"not null;default:0"`
	UnderpaymentEvents                int64      `json:"underpaymentEvents" gorm:"not null;default:0"`
	LastActivityAt                    *time.Time `json:"lastActivityAt"`
	LifetimeMetrics                   string     `json:"lifetimeMetrics" gorm:"type:text"`
	Last30DayMetrics                  string     `json:"last30DayMetrics" gorm:"type:text"`
	Last90DayMetrics                  string     `json:"last90DayMetrics" gorm:"type:text"`
	CreatedAt                         time.Time  `json:"createdAt"`
	UpdatedAt                         time.Time  `json:"updatedAt"`
}

// OfferRankingSnapshot is a versioned marketplace ranking calculation result.
type OfferRankingSnapshot struct {
	ID              string    `json:"id" gorm:"primaryKey;size:36"`
	OfferID         string    `json:"offerId" gorm:"size:36;not null;index:idx_p2p_ranking_offer_id"`
	Rank            int       `json:"rank" gorm:"not null;default:0"`
	Score           string    `json:"score" gorm:"size:30;not null;default:'0'"`
	RankingVersion  int       `json:"rankingVersion" gorm:"not null;default:1"`
	CalculatedAt    time.Time `json:"calculatedAt"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}
