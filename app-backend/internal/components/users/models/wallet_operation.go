package users

import "time"

// Wallet operation states.
const (
	WalletOperationPending   = "PENDING"   // built, waiting for the owners' signatures
	WalletOperationSubmitted = "SUBMITTED" // sent to the bundler
	WalletOperationIncluded  = "INCLUDED"  // mined; Success says whether its calls succeeded
	WalletOperationFailed    = "FAILED"    // rejected by the bundler or reverted
	WalletOperationExpired   = "EXPIRED"   // not signed in time
)

// WalletOperation is a Safe UserOperation a wallet's owners are asked to
// sign (see internal/aa). Its ID is the SafeOp hash the owners sign; the
// apps receive it base64-encoded as a request's "transaction" and sign it
// with personal_sign.
type WalletOperation struct {
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	// ID is the 0x-hex SafeOp hash.
	ID            string `gorm:"size:70;primaryKey" json:"id"`
	WalletAddress string `gorm:"size:56;not null;index:idx_wallet_operation_wallet" json:"walletAddress"`
	// Kind is what the operation does (PAYMENT, SWAP, SUB WALLET, ...).
	Kind string `gorm:"size:40;not null;index:idx_wallet_operation_kind" json:"kind"`
	// Initiator is the username that started it.
	Initiator string `gorm:"size:40;not null" json:"initiator"`
	Status    string `gorm:"size:20;not null;default:'PENDING';index:idx_wallet_operation_status" json:"status"`
	// Prepared is the aa.Prepared JSON: the operation, its hash, validity,
	// owners and gas quote.
	Prepared string `gorm:"type:text;not null" json:"-"`
	// Context is flow-specific JSON used when the operation is submitted or
	// mined (e.g. the fee records to write).
	Context   *string   `gorm:"type:text;null" json:"-"`
	ExpiresAt time.Time `gorm:"not null" json:"expiresAt"`
	// Activation marks the wallet's first operation, which deploys it.
	Activation bool `gorm:"not null;default:false" json:"activation"`
	// GasToken is the stablecoin paying for gas (empty = ETH).
	GasToken string `gorm:"size:56;not null;default:''" json:"gasToken"`
	// UserOpHash is the EntryPoint's hash, known once built; TxHash once mined.
	UserOpHash *string `gorm:"size:70;null;index:idx_wallet_operation_userop" json:"userOpHash"`
	TxHash     *string `gorm:"size:70;null" json:"txHash"`
	Success    *bool   `gorm:"null" json:"success"`
	Error      *string `gorm:"null" json:"error"`
}
