package payments

import (
	"time"
)

// NetworkBase is the only network this engine (and the rest of this
// monorepo) currently watches/supports. Persisted per-row (rather than
// left implicit) so a future second network - a Base-native bridge, or
// support for another chain entirely - doesn't need another schema
// change, just a different value written here.
const NetworkBase = "base"

// PaymentHistory holds payment information. Source is the asset/network
// leaving From/FromAddress; Destination is the asset/network arriving at
// To/ToAddress. For a plain transfer (the only kind this engine currently
// detects - see SavePaymentHistory) the two sides are identical, which is
// correct: a plain payment genuinely has the same asset/network on both
// ends. A row where they differ is, by definition, a swap.
type PaymentHistory struct {
	ID              string
	TransactionType string    `gorm:"index:idx_payment_history_unique_key,unique"`
	TransactionDate time.Time `json:"transactionDate" gorm:"index:idx_payment_history_tx_time"`
	From            *string   `json:"from" gorm:"size:150;index:idx_payment_history_from;null"` //trovoWallet alias and name
	FromAddress     string    `json:"fromAddress" gorm:"size:150;index:idx_payment_history_from_pk;not null;index:idx_payment_history_unique_key,unique;index:idx_payment_history_unique_key,unique"`
	To              *string   `json:"to" gorm:"size:100;index:idx_payment_history_to;null"` //trovoWallet alias and name
	ToAddress       string    `json:"toAddress" gorm:"size:100;index:idx_payment_history_to_pk;not null;index:idx_payment_history_unique_key,unique"`
	Memo            *string   `json:"memo" gorm:"size:60;null"`

	SourceNetwork         string  `json:"sourceNetwork" gorm:"size:20;not null;default:'base'"`
	SourceContractAddress *string `json:"sourceContractAddress" gorm:"size:100;null;"`
	SourceAssetCode       string  `json:"sourceAssetCode" gorm:"size:12;not null;default:''"`
	SourceAmount          string  `json:"sourceAmount"`

	DestinationNetwork         string  `json:"destinationNetwork" gorm:"size:20;not null;default:'base'"`
	DestinationContractAddress *string `json:"destinationContractAddress" gorm:"size:100;null;"`
	DestinationAssetCode       string  `json:"destinationAssetCode" gorm:"size:12;not null;index:idx_payment_history_unique_key,unique"`
	DestinationAmount          string  `json:"destinationAmount" gorm:"index:idx_payment_history_unique_key,unique"`

	TransactionID         string `json:"transactionId" gorm:"size:100;not null;index:idx_payment_history_txid;index:idx_payment_history_unique_key,unique"`
	PT                    string `json:"-" gorm:"size:100;not null;index:idx_payment_history_unique_key,unique;"`
	SourceAccountSequence string `json:"-" gorm:"size:100;not null;index:idx_payment_history_unique_key,unique;"`
}

// PaymentHistoryJSON holds payment information in json format
type PaymentHistoryJSON struct {
	TransactionDate time.Time `json:"transactionDate"`
	TransactionType string    `json:"transactionType"`
	From            string    `json:"from"` //trovoWallet alias and name
	FromAddress     string    `json:"fromAddress"`
	To              string    `json:"to"` //trovoWallet alias and name
	ToAddress       string    `json:"toAddress"`
	Memo            string    `json:"memo"`

	SourceNetwork         string `json:"sourceNetwork"`
	SourceContractAddress string `json:"sourceContractAddress"`
	SourceAssetCode       string `json:"sourceAssetCode"`
	SourceAmount          string `json:"sourceAmount"`

	DestinationNetwork         string `json:"destinationNetwork"`
	DestinationContractAddress string `json:"destinationContractAddress"`
	DestinationAssetCode       string `json:"destinationAssetCode"`
	DestinationAmount          string `json:"destinationAmount"`

	TransactionID string `json:"transactionId"`
}
type TrackedWallet struct {
	ID          string  `gorm:""`
	Address     string  `gorm:"index:idx_tracked_wallet_address,unique"`
	TempAddress *string `gorm:"index:idx_tracked_wallet_temp_address,unique"`
	Alias       string  `gorm:"index:idx_tracked_wallet_alias"`
	Name        string  `gorm:"index:idx_tracked_wallet_name"`
}

type TrackedAddress struct {
	Address string `gorm:"primaryKey;"`
}
type MonitoredCursor struct {
	ID         uint64
	LastCursor string `gorm:"size:100"`
}

type MonitoredAccountCursor struct {
	ID         uint64
	Address    string `gorm:"size:100;uniqueIndex"`
	LastCursor string `gorm:"size:100"`
}
