package payments

import (
	"time"
)

// NetworkBase is the only network app-backend currently watches/supports.
// Persisted per-row (rather than left implicit) so a future second
// network - a Base-native bridge, or support for another chain entirely -
// doesn't need another schema change, just a different value written here.
// Kept in sync with the identical constant in payment-history-engine,
// this table's writer.
const NetworkBase = "base"

// PaymentHistory holds payment information. Source is the asset/network
// leaving From/FromAddress; Destination is the asset/network arriving at
// To/ToAddress. For a plain transfer the two sides are identical; a row
// where they differ is, by definition, a swap.
type PaymentHistory struct {
	ID              string
	TransactionType string    `gorm:"index:idx_payment_history_unique_key,unique"`
	TransactionDate time.Time `json:"transactionDate" gorm:"index:idx_payment_history_tx_time"`
	From            *string   `json:"from" gorm:"size:150;index:idx_payment_history_from;null"` //trovoWallet alias and name
	FromAddress     string    `json:"fromAddress" gorm:"size:150;index:idx_payment_history_from_pk;not null"`
	To              *string   `json:"to" gorm:"size:100;index:idx_payment_history_to;null"` //trovoWallet alias and name
	ToAddress       string    `json:"toAddress" gorm:"size:100;index:idx_payment_history_to_pk;not null;"`
	Memo            *string   `json:"memo" gorm:"size:60;null"`

	SourceNetwork         string  `json:"sourceNetwork" gorm:"size:20;not null;default:'base'"`
	SourceContractAddress *string `json:"sourceContractAddress" gorm:"size:100;null;"`
	SourceAssetCode       string  `json:"sourceAssetCode" gorm:"size:12;not null;default:''"`
	SourceAmount          string  `json:"sourceAmount"`

	DestinationNetwork         string  `json:"destinationNetwork" gorm:"size:20;not null;default:'base'"`
	DestinationContractAddress *string `json:"destinationContractAddress" gorm:"size:100;null;"`
	DestinationAssetCode       string  `json:"destinationAssetCode" gorm:"size:12;not null;"`
	DestinationAmount          string  `json:"destinationAmount" gorm:"index:idx_amount_ph"`

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

type PaginatedPaymentHistory struct {
	Pages        int                  `json:"pages"`
	CurrentPage  int                  `json:"currentPage"`
	TotalRecords int                  `json:"totalRecords"`
	Limit        int                  `json:"limit"`
	Records      []PaymentHistoryJSON `json:"records"`
}

func (ph *PaymentHistory) ToJSON() (json PaymentHistoryJSON) {
	if ph == nil {
		return
	}
	json = PaymentHistoryJSON{
		TransactionDate:      ph.TransactionDate,
		TransactionType:      ph.TransactionType,
		FromAddress:          ph.FromAddress,
		ToAddress:            ph.ToAddress,
		SourceNetwork:        ph.SourceNetwork,
		SourceAssetCode:      ph.SourceAssetCode,
		SourceAmount:         ph.SourceAmount,
		DestinationNetwork:   ph.DestinationNetwork,
		DestinationAssetCode: ph.DestinationAssetCode,
		DestinationAmount:    ph.DestinationAmount,
		TransactionID:        ph.TransactionID,
	}
	if ph.From != nil {
		json.From = *ph.From
	}
	if ph.To != nil {
		json.To = *ph.To
	}
	if ph.SourceContractAddress != nil {
		json.SourceContractAddress = *ph.SourceContractAddress
	}
	if ph.DestinationContractAddress != nil {
		json.DestinationContractAddress = *ph.DestinationContractAddress
	}
	if ph.Memo != nil {
		json.Memo = *ph.Memo
	}
	return
}
