package models

import "time"

// The structs below are read-only mirrors of app-backend's P2P module
// (internal/components/p2p/models in the app-backend repo), which now owns
// this data - see internal/db/main.go's AdminDB doc comment: the wallet,
// P2P, and admin schemas all live in the same physical database, so tm-api
// reads these tables directly rather than calling app-backend over HTTP.
// tm-api never writes to any of these tables - all P2P writes go through
// app-backend's own API (including the admin dispute-resolution action,
// which is a server-to-server call to app-backend, not a direct write
// here). Field sets are trimmed to what the admin dashboard actually
// displays, not a full mirror of every column app-backend defines.

type P2POffer struct {
	ID                 string `json:"id"`
	MerchantUsername   string `json:"merchantUsername"`
	MerchantUserID     string `json:"merchantUserId"`
	OfferType          string `json:"offerType"`
	Asset              string `json:"asset"`
	Price              string `json:"price"`
	Currency           string `json:"currency"`
	Status             string `json:"status"`
	AvailabilityStatus string `json:"availabilityStatus"`
	CreatedAt          time.Time `json:"createdAt"`
}

func (P2POffer) TableName() string { return "offers" }

type P2POrder struct {
	ID                   string    `json:"id"`
	OfferID              string    `json:"offerId"`
	CustomerUserID       string    `json:"customerUserId"`
	CustomerUsername     string    `json:"customerUsername"`
	MerchantUserID       string    `json:"merchantUserId"`
	MerchantUsername     string    `json:"merchantUsername"`
	OfferType            string    `json:"offerType"`
	Asset                string    `json:"asset"`
	Currency             string    `json:"currency"`
	Price                string    `json:"price"`
	SpecifiedAssetAmount string    `json:"specifiedAssetAmount"`
	PaymentAmount        string    `json:"paymentAmount"`
	OrderStatus          string    `json:"orderStatus"`
	IsDisputed           bool      `json:"isDisputed"`
	CreatedAt            time.Time `json:"createdAt"`
}

func (P2POrder) TableName() string { return "orders" }

type P2PDispute struct {
	ID       string    `json:"id"`
	OrderID  string    `json:"orderId"`
	Status   string    `json:"status"`
	OpenedAt time.Time `json:"openedAt"`
}

func (P2PDispute) TableName() string { return "disputes" }

// MerchantPerformance/CustomerPerformance back the admin "P2P users" list
// and the generic user-profile trade-stats fields (UserInfo.TradeStats/
// TakerReputation in user_info.go) - both now read real trading
// performance from here instead of the dead legacy P2P schema's
// maker_stats/taker_reputations tables.
type MerchantPerformance struct {
	MerchantID       string `json:"merchantId"`
	MerchantUsername string `json:"merchantUsername"`
	CompletedTrades  int64  `json:"completedTrades"`
	CompletionRate   string `json:"completionRate"`
	DisputesOpened   int64  `json:"disputesOpened"`
	DisputesResolvedAgainstMerchant int64 `json:"disputesResolvedAgainstMerchant"`
}

func (MerchantPerformance) TableName() string { return "merchant_performances" }

type CustomerPerformance struct {
	CustomerID       string `json:"customerId"`
	CustomerUsername string `json:"customerUsername"`
	CompletedTrades  int64  `json:"completedTrades"`
	CompletionRate   string `json:"completionRate"`
	DisputesOpened   int64  `json:"disputesOpened"`
}

func (CustomerPerformance) TableName() string { return "customer_performances" }
