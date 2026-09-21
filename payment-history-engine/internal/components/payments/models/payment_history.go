package payments

import (
	"time"
)

// PaymentHistory holds payment information
type PaymentHistory struct {
	ID                    string
	TransactionType       string    `gorm:"index:idx_payment_history_unique_key,unique"`
	TransactionDate       time.Time `json:"transactionDate" gorm:"index:idx_payment_history_tx_time"`
	From                  *string   `json:"from" gorm:"size:150;index:idx_payment_history_from;null"` //trovoWallet alias and name
	FromAddress           string    `json:"fromAddress" gorm:"size:150;index:idx_payment_history_from_pk;not null;index:idx_payment_history_unique_key,unique;index:idx_payment_history_unique_key,unique"`
	To                    *string   `json:"to" gorm:"size:100;index:idx_payment_history_to;null"` //trovoWallet alias and name
	ToAddress             string    `json:"toAddress" gorm:"size:100;index:idx_payment_history_to_pk;not null;index:idx_payment_history_unique_key,unique"`
	Memo                  *string   `json:"memo" gorm:"size:60;null"`
	AssetIssuer           *string   `json:"assetIssuer" gorm:"size:100;null;"`
	AssetCode             string    `json:"assetCode" gorm:"size:12;not null;index:idx_payment_history_unique_key,unique"`
	Amount                string    `json:"amount" gorm:"index:idx_payment_history_unique_key,unique"`
	TransactionID         string    `json:"transactionId" gorm:"size:100;not null;index:idx_payment_history_txid;index:idx_payment_history_unique_key,unique"`
	PT                    string    `json:"-" gorm:"size:100;not null;index:idx_payment_history_unique_key,unique;"`
	SourceAccountSequence string    `json:"-" gorm:"size:100;not null;index:idx_payment_history_unique_key,unique;"`
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
	AssetIssuer     string    `json:"assetIssuer"`
	AssetCode       string    `json:"assetCode"`
	Amount          string    `json:"amount"`
	TransactionID   string    `json:"transactionId"`
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
