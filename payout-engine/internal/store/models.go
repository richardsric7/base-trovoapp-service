// Package store holds payout-engine's view of the app-backend tables it
// reads and writes. app-backend owns their schema (its AutoMigrate); these
// structs mirror the columns the engine uses, with explicit table names.
package store

import "time"

// ProceedPayout statuses (app-backend models.ProceedPayoutStatus*).
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
)

// payout item statuses
const (
	ItemPending  = "PENDING"
	ItemQueued   = "QUEUED"
	ItemPaid     = "PAID"
	ItemFailed   = "FAILED"
	ItemExcluded = "EXCLUDED"
	ItemSkipped  = "SKIPPED"
)

// batch statuses
const (
	BatchSubmitted = "SUBMITTED"
	BatchMined     = "MINED"
	BatchReverted  = "REVERTED"
	BatchDropped   = "DROPPED"
)

type ProceedPayout struct {
	ID                          uint64
	CreatedAt                   time.Time
	UpdatedAt                   time.Time
	TokenizedAssetID            string
	Batch                       string
	ProceedPayoutAmount         float64
	AmountPerTokenizedAssetHeld float64
	PaymentScheduleReady        int
	PayoutCompleted             int
	DistributionID              string
	Status                      string
	PayoutAssetCode             string
	PayoutContractAddress       string
	PayoutDecimals              int
	TokenContractAddress        string
	TokenDecimals               int
	PayoutSafeAddress           string
	TotalAmount                 string
	TotalUnits                  string
	AmountPerToken              string
	SnapshotStartBlock          uint64
	SnapshotBlock               uint64
	ScannedBlock                uint64
	SupplyUnits                 string
	EligibleUnits               string
	PayableUnits                string
	RetainedUnits               string
	PaidUnits                   string
	HolderCount                 int
	PaidCount                   int
	FailedCount                 int
	ExcludedCount               int
	ScheduleChecksum            string
	ApprovalsRequired           int
	PreparedBy                  string
	PreparationRequestedAt      *time.Time
	LockedAt                    *time.Time
	ApprovedAt                  *time.Time
	FundingRequestedBy          string
	FundingCheckedAt            *time.Time
	StartedAt                   *time.Time
	CompletedAt                 *time.Time
	StatusBeforePause           string
	Note                        string
	FeeType                     string
	FeeValue                    string
	FeeCap                      string
	FeeSetBy                    string
	FeeUnits                    string
	VatPercent                  string
	VatUnits                    string
	FeeWallet                   string
	VatWallet                   string
	HolderPayable               string `gorm:"column:holder_payable"`
}

func (ProceedPayout) TableName() string { return "proceed_payouts" }

// PayoutItem is a TokenizedAssetPayoutSchedule row: one holder.
type PayoutItem struct {
	ID                             string
	CreatedAt                      time.Time
	UpdatedAt                      time.Time
	TokenizedAssetID               string
	Batch                          string
	PayoutAssetCode                string
	PayoutContractAddress          string
	BeneficiaryAddress             string
	ConfirmedTokenizedAssetBalance float64
	AmountToReceive                float64
	CannotReceiveAsset             int
	ProceedPayoutID                uint64
	BalanceUnits                   string
	AmountUnits                    string
	Status                         string
	Reason                         string
	Username                       string
	BatchID                        uint64
	TxHash                         string
	PaidAt                         *time.Time
	Notified                       int
	ActionBy                       string
	Kind                           string `gorm:"default:'HOLDER'"`
}

// payout item kinds and fee types (app-backend models.PayoutItemKind* / PayoutFee*)
const (
	KindHolder   = "HOLDER"
	KindFee      = "FEE"
	KindVat      = "VAT"
	FeeFixed     = "FIXED"
	FeePercent   = "PERCENT"
	FeeServiceID = "PROCEED_PAYOUT_FEE"
	// fee_collections fee types
	FeeTypePayout    = "PROCEED_PAYOUT_FEE"
	FeeTypePayoutVat = "PROCEED_PAYOUT_FEE_VAT"
)

// ServiceFee is app-backend's service_fees row: the payout fee's wallet
// (PROCEED_PAYOUT_FEE) and the VAT wallet (VAT).
type ServiceFee struct {
	ID                 string
	FeeWalletSecretKey string
	Inactive           int
}

func (ServiceFee) TableName() string { return "service_fees" }

// CountryConfig holds a country's VAT rate.
type CountryConfig struct {
	CountryCode string  `gorm:"primaryKey"`
	VATPercent  float64 `gorm:"column:vat_percent"`
}

func (CountryConfig) TableName() string { return "country_configs" }

// FeeCollection is app-backend's record of a collected fee (TM's fee report).
type FeeCollection struct {
	CreatedAt         time.Time
	UpdatedAt         time.Time
	ID                string `gorm:"primaryKey"`
	FromUsername      string
	FromWalletAddress string
	FromWalletAlias   string
	FeeType           string
	Amount            float64
	AssetCode         string
	ContractAddress   *string
	DestinationWallet string
	TransactionHash   *string
	Processed         int
}

func (FeeCollection) TableName() string { return "fee_collections" }

func (PayoutItem) TableName() string { return "tokenized_asset_payout_schedules" }

type PayoutBatch struct {
	ID              uint64
	CreatedAt       time.Time
	UpdatedAt       time.Time
	ProceedPayoutID uint64
	Status          string
	SafeNonce       uint64
	SafeTxHash      string
	TxHash          string
	RawTx           string
	Executor        string
	ItemCount       int
	AmountUnits     string
	GasUsed         uint64
	Error           string
}

func (PayoutBatch) TableName() string { return "proceed_payout_batches" }

type PayoutApproval struct {
	ID               uint64
	CreatedAt        time.Time
	ProceedPayoutID  uint64
	AdminEmail       string
	ScheduleChecksum string
}

func (PayoutApproval) TableName() string { return "proceed_payout_approvals" }

type TokenIndex struct {
	Token        string `gorm:"primaryKey"`
	StartBlock   uint64
	ScannedBlock uint64
	UpdatedAt    time.Time
}

func (TokenIndex) TableName() string { return "payout_token_indexes" }

type TokenBalance struct {
	Token   string `gorm:"primaryKey"`
	Holder  string `gorm:"primaryKey"`
	Balance string
}

func (TokenBalance) TableName() string { return "payout_token_balances" }

type EngineState struct {
	ID               uint64
	Halted           bool
	HaltReason       string
	HaltedBy         string
	HaltedAt         *time.Time
	HeartbeatAt      *time.Time
	Instance         string
	Version          string
	Activity         string
	LastError        string
	SweepToken       string
	SweepRequestedBy string
	SweepRequestedAt *time.Time
	SweepResult      string
}

func (EngineState) TableName() string { return "payout_engine_states" }

// TokenizedAsset is the part of app-backend's tokenized_assets the engine reads.
type TokenizedAsset struct {
	ID                   string
	CreatedAt            time.Time
	AssetCode            *string
	AssetName            *string
	ContractAddress      *string
	MarketMakingWallet   *string
	AssetCountryLocation *string
}

func (TokenizedAsset) TableName() string { return "tokenized_assets" }

type User struct {
	ID                    string
	Username              string
	PushNotificationToken *string
}

func (User) TableName() string { return "users" }

type UserWallet struct {
	ID     string
	UserID string
	Alias  string
}

func (UserWallet) TableName() string { return "user_wallets" }

// OfferBookOffer is app-backend's index of TrovoOfferBook offers: tokens on
// sale are held by the book, so they are credited back to their sellers.
type OfferBookOffer struct {
	ID        string
	Seller    string
	SellToken string
	Remaining string
	Open      bool
}

func (OfferBookOffer) TableName() string { return "offer_book_offers" }

// AllModels lists them, for tests that build the schema themselves.
func AllModels() []interface{} {
	return []interface{}{&ProceedPayout{}, &PayoutItem{}, &PayoutBatch{}, &PayoutApproval{}, &TokenIndex{}, &TokenBalance{}, &EngineState{}, &TokenizedAsset{}, &User{}, &UserWallet{}, &OfferBookOffer{}, &ServiceFee{}, &CountryConfig{}, &FeeCollection{}}
}
