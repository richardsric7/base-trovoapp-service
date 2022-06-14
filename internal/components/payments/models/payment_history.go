package payments

import (
	"time"
)

//PaymentHistory holds payment information
type PaymentHistory struct {
	ID              string
	TransactionType string
	TransactionDate time.Time `json:"transactionDate" gorm:"index:idx_payment_history_tx_time"`
	From            *string   `json:"from" gorm:"size:150;index:idx_payment_history_from;null"` //trovoWallet alias and name
	FromPublicKey   string    `json:"fromPublicKey" gorm:"size:150;index:idx_payment_history_from_pk;not null"`
	To              *string   `json:"to" gorm:"size:56;index:idx_payment_history_to;null"` //trovoWallet alias and name
	ToPublicKey     string    `json:"toPublicKey" gorm:"size:56;index:idx_payment_history_to_pk;not null"`
	Memo            *string   `json:"memo" gorm:"size:28;null"`
	AssetIssuer     *string   `json:"assetIssuer" gorm:"size:56;null"`
	AssetCode       string    `json:"assetCode" gorm:"size:12;not null"`
	Amount          string    `json:"amount"`
	TransactionID   string    `json:"transactionId" gorm:"size:70;not null;index:idx_payment_history_txid"`
}

//PaymentHistoryJSON holds payment information in json format
type PaymentHistoryJSON struct {
	TransactionDate time.Time `json:"transactionDate"`
	TransactionType string    `json:"transactionType"`
	From            string    `json:"from"` //trovoWallet alias and name
	FromPublicKey   string    `json:"fromPublicKey"`
	To              string    `json:"to"` //trovoWallet alias and name
	ToPublicKey     string    `json:"toPublicKey"`
	Memo            string    `json:"memo"`
	AssetIssuer     string    `json:"assetIssuer"`
	AssetCode       string    `json:"assetCode"`
	Amount          string    `json:"amount"`
	TransactionID   string    `json:"transactionId"`
}
