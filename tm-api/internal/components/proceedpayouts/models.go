// Package proceedpayouts is TM's side of token proceeds payouts: it
// registers the payout when a trustee authorizes a stakeholder distribution,
// and gives Trovo admins the controls over the payout-engine's stages
// (prepare, fee, approvals, funding, pause/resume, cancel, item exclusions,
// retries, sweep, kill switch) and the payout, fee and VAT reports.
//
// payout-engine (its own project, no HTTP) does the work. Both share
// app-backend's database (the tables below are mirrors of app-backend's
// models, users/models/proceeds_payout.go and tokenization.go): tm-api
// changes a payout's state with conditional updates, then wakes the engine
// over Redis (payout-engine:commands); the engine also polls.
package proceedpayouts

import "time"

// Payout statuses (see app-backend users/models/proceeds_payout.go).
const (
	StatusRegistered            = "REGISTERED"
	StatusPrepareRequested      = "PREPARE_REQUESTED"
	StatusPreparing             = "PREPARING"
	StatusLocked                = "LOCKED"
	StatusApproved              = "APPROVED"
	StatusFundingCheckRequested = "FUNDING_CHECK_REQUESTED"
	StatusPaying                = "PAYING"
	StatusPaused                = "PAUSED"
	StatusCompleted             = "COMPLETED"
	StatusCompletedWithFailures = "COMPLETED_WITH_FAILURES"
	StatusCancelled             = "CANCELLED"

	ItemPending  = "PENDING"
	ItemQueued   = "QUEUED"
	ItemPaid     = "PAID"
	ItemFailed   = "FAILED"
	ItemExcluded = "EXCLUDED"
	ItemSkipped  = "SKIPPED"

	KindHolder = "HOLDER"
	KindFee    = "FEE"
	KindVat    = "VAT"

	FeeFixed   = "FIXED"
	FeePercent = "PERCENT"

	// FeeServiceID is the service_fees row with the payout fee wallet and
	// the default fee for new payouts.
	FeeServiceID = "PROCEED_PAYOUT_FEE"
	// fee_collections fee types the engine records
	FeeTypePayout    = "PROCEED_PAYOUT_FEE"
	FeeTypePayoutVat = "PROCEED_PAYOUT_FEE_VAT"

	// CommandsChannel wakes payout-engine after a change.
	CommandsChannel = "payout-engine:commands"

	batchSubmitted = "SUBMITTED"
)

// ProceedPayout is one distribution's payout (proceed_payouts).
type ProceedPayout struct {
	ID                     uint64     `json:"id"`
	CreatedAt              time.Time  `json:"createdAt"`
	UpdatedAt              time.Time  `json:"updatedAt"`
	TokenizedAssetID       string     `json:"tokenizedAssetId"`
	Batch                  string     `json:"batch"`
	DistributionID         string     `json:"distributionId"`
	Status                 string     `json:"status"`
	PayoutAssetCode        string     `json:"payoutAssetCode"`
	PayoutContractAddress  string     `json:"payoutContractAddress"`
	PayoutDecimals         int        `json:"payoutDecimals"`
	TokenContractAddress   string     `json:"tokenContractAddress"`
	TokenDecimals          int        `json:"tokenDecimals"`
	PayoutSafeAddress      string     `json:"payoutSafeAddress"`
	TotalAmount            string     `json:"totalAmount"`
	TotalUnits             string     `json:"totalUnits"`
	AmountPerToken         string     `json:"amountPerToken"`
	SnapshotStartBlock     uint64     `json:"snapshotStartBlock"`
	SnapshotBlock          uint64     `json:"snapshotBlock"`
	ScannedBlock           uint64     `json:"scannedBlock"`
	SupplyUnits            string     `json:"supplyUnits"`
	EligibleUnits          string     `json:"eligibleUnits"`
	PayableUnits           string     `json:"payableUnits"`
	RetainedUnits          string     `json:"retainedUnits"`
	PaidUnits              string     `json:"paidUnits"`
	HolderCount            int        `json:"holderCount"`
	PaidCount              int        `json:"paidCount"`
	FailedCount            int        `json:"failedCount"`
	ExcludedCount          int        `json:"excludedCount"`
	ScheduleChecksum       string     `json:"scheduleChecksum"`
	ApprovalsRequired      int        `json:"approvalsRequired"`
	PaymentScheduleReady   int        `json:"paymentScheduleReady"`
	PayoutCompleted        int        `json:"payoutCompleted"`
	PreparedBy             string     `json:"preparedBy"`
	PreparationRequestedAt *time.Time `json:"preparationRequestedAt"`
	LockedAt               *time.Time `json:"lockedAt"`
	ApprovedAt             *time.Time `json:"approvedAt"`
	FundingRequestedBy     string     `json:"fundingRequestedBy"`
	FundingCheckedAt       *time.Time `json:"fundingCheckedAt"`
	StartedAt              *time.Time `json:"startedAt"`
	CompletedAt            *time.Time `json:"completedAt"`
	StatusBeforePause      string     `json:"statusBeforePause"`
	Note                   string     `json:"note"`
	FeeType                string     `json:"feeType"`
	FeeValue               string     `json:"feeValue"`
	FeeCap                 string     `json:"feeCap"`
	FeeSetBy               string     `json:"feeSetBy"`
	FeeUnits               string     `json:"feeUnits"`
	VatPercent             string     `json:"vatPercent"`
	VatUnits               string     `json:"vatUnits"`
	FeeWallet              string     `json:"feeWallet"`
	VatWallet              string     `json:"vatWallet"`
	HolderPayable          string     `json:"holderPayableUnits"`
}

func (ProceedPayout) TableName() string { return "proceed_payouts" }

// PayoutItem is one line of a payout's schedule: a holder, or the fee / VAT.
type PayoutItem struct {
	ID                             string     `json:"id"`
	CreatedAt                      time.Time  `json:"createdAt"`
	UpdatedAt                      time.Time  `json:"updatedAt"`
	TokenizedAssetID               string     `json:"tokenizedAssetId"`
	Batch                          string     `json:"batch"`
	PayoutAssetCode                string     `json:"payoutAssetCode"`
	PayoutContractAddress          string     `json:"payoutContractAddress"`
	BeneficiaryAddress             string     `json:"beneficiaryAddress"`
	ConfirmedTokenizedAssetBalance float64    `json:"confirmedTokenizedAssetBalance"`
	AmountToReceive                float64    `json:"amountToReceive"`
	CannotReceiveAsset             int        `json:"cannotReceiveAsset"`
	ProceedPayoutID                uint64     `json:"proceedPayoutId"`
	BalanceUnits                   string     `json:"balanceUnits"`
	AmountUnits                    string     `json:"amountUnits"`
	Status                         string     `json:"status"`
	Reason                         string     `json:"reason"`
	Username                       string     `json:"username"`
	BatchID                        uint64     `json:"batchId"`
	TxHash                         string     `json:"txHash"`
	PaidAt                         *time.Time `json:"paidAt"`
	ActionBy                       string     `json:"actionBy"`
	Kind                           string     `json:"kind"`
}

func (PayoutItem) TableName() string { return "tokenized_asset_payout_schedules" }

// PayoutBatch is one Safe transaction of a payout (its raw transaction is
// not exposed).
type PayoutBatch struct {
	ID              uint64    `json:"id"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
	ProceedPayoutID uint64    `json:"proceedPayoutId"`
	Status          string    `json:"status"`
	SafeNonce       uint64    `json:"safeNonce"`
	SafeTxHash      string    `json:"safeTxHash"`
	TxHash          string    `json:"txHash"`
	Executor        string    `json:"executor"`
	ItemCount       int       `json:"itemCount"`
	AmountUnits     string    `json:"amountUnits"`
	GasUsed         uint64    `json:"gasUsed"`
	Error           string    `json:"error"`
}

func (PayoutBatch) TableName() string { return "proceed_payout_batches" }

// PayoutApproval is one admin's approval of a locked schedule.
type PayoutApproval struct {
	ID               uint64    `json:"id"`
	CreatedAt        time.Time `json:"createdAt"`
	ProceedPayoutID  uint64    `json:"proceedPayoutId"`
	AdminEmail       string    `json:"adminEmail"`
	ScheduleChecksum string    `json:"scheduleChecksum"`
}

func (PayoutApproval) TableName() string { return "proceed_payout_approvals" }

// EngineState is payout-engine's shared state row (ID 1).
type EngineState struct {
	ID               uint64     `json:"id"`
	Halted           bool       `json:"halted"`
	HaltReason       string     `json:"haltReason"`
	HaltedBy         string     `json:"haltedBy"`
	HaltedAt         *time.Time `json:"haltedAt"`
	HeartbeatAt      *time.Time `json:"heartbeatAt"`
	Instance         string     `json:"instance"`
	Version          string     `json:"version"`
	Activity         string     `json:"activity"`
	LastError        string     `json:"lastError"`
	SweepToken       string     `json:"sweepToken"`
	SweepRequestedBy string     `json:"sweepRequestedBy"`
	SweepRequestedAt *time.Time `json:"sweepRequestedAt"`
	SweepResult      string     `json:"sweepResult"`
}

func (EngineState) TableName() string { return "payout_engine_states" }

// ServiceFee is app-backend's service_fees row; for the payout fee,
// FeeWalletSecretKey holds the fee wallet's address and the default fee is
// FeeFixed, or FeePercent with Remarks as the cap.
type ServiceFee struct {
	CreatedAt          time.Time
	UpdatedAt          time.Time
	LastUpdatedBy      string
	ID                 string `gorm:"primaryKey"`
	FeeWalletSecretKey string
	FeePercent         float64
	FeeFixed           float64
	Inactive           int
	Remarks            string
}

func (ServiceFee) TableName() string { return "service_fees" }

type tokenizationCurrency struct {
	AssetCode       string `gorm:"primaryKey"`
	ContractAddress string
}

func (tokenizationCurrency) TableName() string { return "tokenization_currencies" }

type countryConfig struct {
	CountryCode              string `gorm:"primaryKey"`
	VATPercent               float64
	InternalBalanceTokenCode *string
}

func (countryConfig) TableName() string { return "country_configs" }

// FeeCollection is a collected fee row (the payout fee and its VAT, as the
// engine records them when paid).
type FeeCollection struct {
	ID                string    `json:"id"`
	CreatedAt         time.Time `json:"createdAt"`
	FromUsername      string    `json:"assetCode"`   // the asset whose payout paid it
	FromWalletAddress string    `json:"payoutSafe"`  // the payout Safe
	FromWalletAlias   string    `json:"payoutBatch"` // "payout <batch>"
	FeeType           string    `json:"feeType"`
	Amount            float64   `json:"amount"`
	AssetCode         string    `json:"payoutAssetCode"`
	DestinationWallet string    `json:"destinationWallet"`
	TransactionHash   *string   `json:"transactionHash"`
}

func (FeeCollection) TableName() string { return "fee_collections" }
