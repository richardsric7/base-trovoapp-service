package p2p

import "time"

// Refund reasons (Plan Section 45)
const (
	RefundReasonUnderpayment     = "UNDERPAYMENT"
	RefundReasonOverpayment      = "OVERPAYMENT"
	RefundReasonWrongAsset       = "WRONG_ASSET"
	RefundReasonUnknownOrder     = "UNKNOWN_ORDER"
	RefundReasonDuplicatePayment = "DUPLICATE_PAYMENT"
	RefundReasonOrderExpired     = "ORDER_EXPIRED"
	RefundReasonOther            = "OTHER"
)

// Refund tracks a refundable deposit and its claim state. Claimed is set
// before the transfer is submitted (checks-effects-interactions) to prevent
// double-claims/reentrancy.
type Refund struct {
	ID              string     `json:"id" gorm:"primaryKey;size:36"`
	DepositID       string     `json:"depositId" gorm:"size:36;index:idx_p2p_refund_deposit_id"`
	OrderID         string     `json:"orderId" gorm:"size:36;not null;index:idx_p2p_refund_order_id"`
	Sender          string     `json:"sender" gorm:"size:56;not null;index:idx_p2p_refund_sender"`
	Token           string     `json:"token" gorm:"size:12;not null"`
	ContractAddress string     `json:"contractAddress" gorm:"size:56;default:''"`
	Amount          string     `json:"amount" gorm:"size:60;not null"`
	Reason          string     `json:"reason" gorm:"size:30;not null"`
	Claimed         bool       `json:"claimed" gorm:"not null;default:false"`
	ClaimedAt       *time.Time `json:"claimedAt"`
	TransactionHash string     `json:"transactionHash" gorm:"size:100;default:''"`
	CreatedAt       time.Time  `json:"createdAt"`
	UpdatedAt       time.Time  `json:"updatedAt"`
}
