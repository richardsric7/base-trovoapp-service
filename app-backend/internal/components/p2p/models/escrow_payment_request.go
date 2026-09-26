package p2p

import "time"

// P2PEscrowPaymentRequest statuses
const (
	EscrowPaymentRequestStatusPending = "PENDING"
	EscrowPaymentRequestStatusUsed    = "USED"
	EscrowPaymentRequestStatusExpired = "EXPIRED"
)

// P2PEscrowPaymentRequest is the request-security layer wrapped around the
// existing shortlink/QR mechanism (Plan Section 39). The underlying
// dynamic_links row has no binding/expiry/nonce/status of its own - this
// model supplies all of it. Per Section 39's enforcement-timing gap, expiry/
// nonce/status here can only be checked at shortlink-resolve time and at
// deposit-detection time, never at the moment funds actually move on-chain;
// a late/duplicate deposit is refunded (Section 44), not blocked.
type P2PEscrowPaymentRequest struct {
	ID        string     `json:"id" gorm:"primaryKey;size:36"`
	OrderID   string     `json:"orderId" gorm:"size:36;not null;index:idx_p2p_escrow_req_order_id"`
	Token     string     `json:"token" gorm:"size:12;not null"` // asset code
	Amount    string     `json:"amount" gorm:"size:60;not null"`
	Recipient string     `json:"recipient" gorm:"size:56;not null"` // escrow address
	Network   string     `json:"network" gorm:"size:20;not null;default:'Base'"`
	Expiry    time.Time  `json:"expiry"`
	Nonce     string     `json:"nonce" gorm:"size:36;not null;uniqueIndex:idx_p2p_escrow_req_nonce"`
	Status    string     `json:"status" gorm:"size:20;not null;default:'PENDING'"`
	CreatedBy string     `json:"createdBy" gorm:"size:100;not null"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
}

// TableName overrides GORM's default naming, which mangles "P2P" (a
// capital letter following a digit) into "p2_p_escrow_payment_requests".
func (P2PEscrowPaymentRequest) TableName() string {
	return "p2p_escrow_payment_requests"
}

// BlockchainDeposit records every detected deposit transaction against an
// order's escrow address, canonical or not (Plan Section 21/44). Multiple
// rows can exist per order (partial/overpaid/duplicate deposits); Order's own
// EscrowDepositTransactionHash reflects only the canonical successful one.
type BlockchainDeposit struct {
	ID              string    `json:"id" gorm:"primaryKey;size:36"`
	OrderID         string    `json:"orderId" gorm:"size:36;not null;index:idx_p2p_deposit_order_id"`
	Sender          string    `json:"sender" gorm:"size:56;not null"`
	Token           string    `json:"token" gorm:"size:12;not null"`
	ContractAddress string    `json:"contractAddress" gorm:"size:56;default:''"`
	Amount          string    `json:"amount" gorm:"size:60;not null"`
	TransactionHash string    `json:"transactionHash" gorm:"size:100;not null;uniqueIndex:idx_p2p_deposit_tx_hash"`
	IsCanonical     bool      `json:"isCanonical" gorm:"not null;default:false"`
	DetectionSource string    `json:"detectionSource" gorm:"size:30;not null;default:'P2P_API'"` // P2P_API (own deposit endpoint) or PAYMENT_HISTORY_MATCH (third-party via shared link)
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}
