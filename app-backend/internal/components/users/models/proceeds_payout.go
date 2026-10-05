package users

import "time"

// ProceedPayout statuses. tm-api moves a payout between them on admin
// actions; payout-engine on its own progress.
//
//	REGISTERED -> (admin: prepare) PREPARE_REQUESTED -> (engine) PREPARING -> LOCKED
//	LOCKED -> (admins: approve, distinct, ApprovalsRequired of them) APPROVED
//	LOCKED -> (admin: re-prepare) PREPARE_REQUESTED          (approvals are dropped)
//	APPROVED -> (admin: confirm funding) FUNDING_CHECK_REQUESTED
//	FUNDING_CHECK_REQUESTED -> (engine: Safe holds the total, executor has gas) PAYING
//	                        -> (engine: not funded) APPROVED, with Note
//	PAYING -> (engine) COMPLETED | COMPLETED_WITH_FAILURES
//	PAYING <-> PAUSED (admin)            COMPLETED_WITH_FAILURES -> (admin: retry failed) PAYING
//	any state before PAYING, or PAUSED -> (admin) CANCELLED
const (
	ProceedPayoutStatusRegistered            = "REGISTERED"
	ProceedPayoutStatusPrepareRequested      = "PREPARE_REQUESTED"
	ProceedPayoutStatusPreparing             = "PREPARING"
	ProceedPayoutStatusLocked                = "LOCKED"
	ProceedPayoutStatusApproved              = "APPROVED"
	ProceedPayoutStatusFundingCheckRequested = "FUNDING_CHECK_REQUESTED"
	ProceedPayoutStatusPaying                = "PAYING"
	ProceedPayoutStatusPaused                = "PAUSED"
	ProceedPayoutStatusCompleted             = "COMPLETED"
	ProceedPayoutStatusCompletedWithFailures = "COMPLETED_WITH_FAILURES"
	ProceedPayoutStatusCancelled             = "CANCELLED"
)

// TokenizedAssetPayoutSchedule (payout item) statuses.
const (
	PayoutItemPending  = "PENDING"  // to be paid
	PayoutItemQueued   = "QUEUED"   // in a batch being executed
	PayoutItemPaid     = "PAID"     // paid on-chain (TxHash), or marked paid by an admin
	PayoutItemFailed   = "FAILED"   // the transfer to it reverts (e.g. a blocked address)
	PayoutItemExcluded = "EXCLUDED" // left out by an admin, or a platform wallet
	PayoutItemSkipped  = "SKIPPED"  // its share rounds to nothing (or is below the minimum)
)

// Payout item kinds, and payout fee types.
const (
	PayoutItemKindHolder = "HOLDER"
	PayoutItemKindFee    = "FEE"
	PayoutItemKindVat    = "VAT"
	PayoutFeeFixed       = "FIXED"
	PayoutFeePercent     = "PERCENT"
	// ProceedPayoutFeeServiceFeeID is the service_fees row holding the
	// payout fee wallet (FeeWalletSecretKey: its address) and the default
	// fee for new payouts (FeeFixed, or FeePercent with Remarks as the cap).
	ProceedPayoutFeeServiceFeeID = "PROCEED_PAYOUT_FEE"
	// fee_collections fee types for payout fees and their VAT
	FeeTypeProceedPayout    = "PROCEED_PAYOUT_FEE"
	FeeTypeProceedPayoutVat = "PROCEED_PAYOUT_FEE_VAT"
)

// ProceedPayoutBatch is one Safe transaction paying a batch of holders
// through MultiSendCallOnly. It is recorded, with its signed transaction,
// before it is broadcast, so an engine that stops mid-way can tell whether
// it was mined.
type ProceedPayoutBatch struct {
	ID              uint64    `json:"id"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
	ProceedPayoutID uint64    `gorm:"index" json:"proceedPayoutId"`
	Status          string    `gorm:"size:16;index" json:"status"` // see PayoutBatch*
	SafeNonce       uint64    `json:"safeNonce"`
	SafeTxHash      string    `gorm:"size:70" json:"safeTxHash"`
	TxHash          string    `gorm:"size:70;index" json:"txHash"`
	RawTx           string    `gorm:"type:text" json:"-"`
	Executor        string    `gorm:"size:50" json:"executor"` // the signer that broadcast it and paid its gas
	ItemCount       int       `json:"itemCount"`
	AmountUnits     string    `gorm:"size:80" json:"amountUnits"`
	GasUsed         uint64    `json:"gasUsed"`
	Error           string    `gorm:"size:500" json:"error"`
}

func (ProceedPayoutBatch) TableName() string { return "proceed_payout_batches" }

const (
	PayoutBatchSubmitted = "SUBMITTED" // signed, recorded, being broadcast / waiting to be mined
	PayoutBatchMined     = "MINED"
	PayoutBatchReverted  = "REVERTED" // its items go back to PENDING (or FAILED once isolated)
	PayoutBatchDropped   = "DROPPED"  // never mined (its Safe nonce was used otherwise); items go back to PENDING
)

// ProceedPayoutApproval is one admin's approval of a locked schedule
// (ScheduleChecksum); a re-prepared schedule needs new approvals.
type ProceedPayoutApproval struct {
	ID               uint64    `json:"id"`
	CreatedAt        time.Time `json:"createdAt"`
	ProceedPayoutID  uint64    `gorm:"index:idx_payout_approval_admin,unique,priority:1" json:"proceedPayoutId"`
	AdminEmail       string    `gorm:"size:150;index:idx_payout_approval_admin,unique,priority:2" json:"adminEmail"`
	ScheduleChecksum string    `gorm:"size:80" json:"scheduleChecksum"`
}

func (ProceedPayoutApproval) TableName() string { return "proceed_payout_approvals" }

// PayoutTokenIndex is payout-engine's replay of a token's Transfer events:
// balances (PayoutTokenBalance) are exact as of ScannedBlock, so each new
// payout of the token only scans the blocks since the last one.
type PayoutTokenIndex struct {
	Token        string    `gorm:"primaryKey;size:50" json:"token"`
	StartBlock   uint64    `json:"startBlock"`
	ScannedBlock uint64    `json:"scannedBlock"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

func (PayoutTokenIndex) TableName() string { return "payout_token_indexes" }

type PayoutTokenBalance struct {
	Token   string `gorm:"primaryKey;size:50"`
	Holder  string `gorm:"primaryKey;size:50"`
	Balance string `gorm:"size:80"` // base units
}

func (PayoutTokenBalance) TableName() string { return "payout_token_balances" }

// PayoutEngineState is the single row (ID 1) through which payout-engine and
// tm-api share the engine's state: its heartbeat and the kill switch that
// stops all payouts at once.
type PayoutEngineState struct {
	ID          uint64     `json:"id"`
	Halted      bool       `json:"halted"`
	HaltReason  string     `gorm:"size:500" json:"haltReason"`
	HaltedBy    string     `gorm:"size:150" json:"haltedBy"`
	HaltedAt    *time.Time `json:"haltedAt"`
	HeartbeatAt *time.Time `json:"heartbeatAt"`
	Instance    string     `gorm:"size:150" json:"instance"`
	Version     string     `gorm:"size:80" json:"version"`
	Activity    string     `gorm:"size:500" json:"activity"` // what it is doing now
	LastError   string     `gorm:"size:1000" json:"lastError"`
	// an admin's request to move the payout Safe's whole balance of
	// SweepToken to the engine's PROCEED_PAYOUT_SWEEP_ADDRESS; refused while a
	// payout of that token is being funded or paid
	SweepToken       string     `gorm:"size:50" json:"sweepToken"`
	SweepRequestedBy string     `gorm:"size:150" json:"sweepRequestedBy"`
	SweepRequestedAt *time.Time `json:"sweepRequestedAt"`
	SweepResult      string     `gorm:"size:500" json:"sweepResult"`
}

func (PayoutEngineState) TableName() string { return "payout_engine_states" }
