// Package publicmarkets holds the Public Markets tables: tokenized NGX/FMDQ
// instruments held at CSCS through licensed Custodians, the orders that
// create and redeem their tokens, the instructions sent to Custodians and
// Dealing Members, reconciliation, exchange partners and dividends.
//
// Amounts and quantities are decimal strings (human units), like the
// payout tables, so no precision is lost on any database.
package publicmarkets

import "time"

// Asset statuses.
const (
	AssetSetup  = "SETUP"  // created, not yet open for orders
	AssetLive   = "LIVE"   // creation and redemption open
	AssetHalted = "HALTED" // reconciliation drift or an admin halt: new orders rejected
)

// Markets and asset types.
const (
	MarketNGX  = "NGX"
	MarketFMDQ = "FMDQ"
	TypeEquity = "EQUITY"
	TypeBond   = "BOND"
)

// Asset is one tokenized public-market instrument (doc: PublicMarketAsset).
type Asset struct {
	ID             string    `gorm:"primaryKey;size:64" json:"id"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
	AssetCode      string    `gorm:"size:16;uniqueIndex;not null" json:"assetCode"` // e.g. MTNN-T
	Ticker         string    `gorm:"size:16;not null" json:"ticker"`                // e.g. MTNN
	Market         string    `gorm:"size:10;not null" json:"market"`                // NGX | FMDQ
	AssetType      string    `gorm:"size:10;not null" json:"assetType"`             // EQUITY | BOND
	ISIN           string    `gorm:"size:12" json:"isin"`
	InstrumentName string    `gorm:"size:200;not null" json:"instrumentName"`
	ShortName      string    `gorm:"size:100" json:"shortName"`
	Sector         string    `gorm:"size:100" json:"sector"`
	Description    string    `gorm:"type:text" json:"description"`
	Country        string    `gorm:"size:2;not null;default:'NG'" json:"country"`
	Currency       string    `gorm:"size:5;not null;default:'NGN'" json:"currency"`
	// UnitDescription says what one token stands for (one share, or
	// 100 NGN face value of a bond).
	UnitDescription string `gorm:"size:100" json:"unitDescription"`
	// display
	LogoBackground string `gorm:"size:16" json:"logoBackground"`
	LogoForeground string `gorm:"size:16" json:"logoForeground"`
	LogoInitials   string `gorm:"size:6" json:"logoInitials"`
	LogoURL        string `gorm:"size:300" json:"logoUrl"`
	// key statistics shown in the apps (free text from the data vendor or ops)
	MarketCap     string     `gorm:"size:40" json:"marketCap"`
	PERatio       string     `gorm:"size:20" json:"peRatio"`
	DividendYield string     `gorm:"size:40" json:"dividendYield"`
	Coupon        string     `gorm:"size:40" json:"coupon"`
	MaturityDate  *time.Time `json:"maturityDate"`
	// the token
	ContractAddress    string `gorm:"size:64" json:"contractAddress"`
	IssuingSafeAddress string `gorm:"size:64" json:"issuingSafeAddress"` // owns the token, holds redeemed tokens until burnt
	TokenDecimals      int    `gorm:"not null;default:0" json:"tokenDecimals"`
	// custody and execution
	CustodianID      uint64 `gorm:"not null;index" json:"custodianId"`     // approved_asset_custodians.id
	DealingMemberID  uint64 `gorm:"not null;index" json:"dealingMemberId"` // approved_dealing_members.id
	OmnibusReference string `gorm:"size:120" json:"omnibusReference"`      // trovotechAccountReference, e.g. POOL-DANGCEM-01
	// trading rules
	MinimumBuy           string `gorm:"size:40;not null;default:'1000'" json:"minimumBuy"` // in the funding currency
	FeePercent           string `gorm:"size:20" json:"feePercent"`                         // empty: the settings' trade fee
	InventoryTargetUnits string `gorm:"size:40;not null;default:'0'" json:"inventoryTargetUnits"`
	// status
	Status     string     `gorm:"size:12;not null;default:'SETUP';index" json:"status"`
	HaltReason string     `gorm:"size:500" json:"haltReason"`
	HaltedAt   *time.Time `json:"haltedAt"`
	HaltedBy   string     `gorm:"size:150" json:"haltedBy"`
	// latest reference price (also in PriceSnapshot), kept here for listing
	LastPrice     string     `gorm:"size:40;not null;default:'0'" json:"lastPrice"`
	PreviousClose string     `gorm:"size:40;not null;default:'0'" json:"previousClose"`
	DayHigh       string     `gorm:"size:40;not null;default:'0'" json:"dayHigh"`
	DayLow        string     `gorm:"size:40;not null;default:'0'" json:"dayLow"`
	DayVolume     string     `gorm:"size:40" json:"dayVolume"`
	PriceSource   string     `gorm:"size:30" json:"priceSource"`
	PriceAt       *time.Time `json:"priceAt"`
	// ledger indexer progress (Transfer logs of the token)
	LedgerStartBlock   uint64 `json:"-"`
	LedgerScannedBlock uint64 `json:"-"`
	CreatedBy          string `gorm:"size:150" json:"createdBy"`
}

func (Asset) TableName() string { return "public_market_assets" }

// Integration modes of a Custodian or Dealing Member.
const (
	ModeMock   = "MOCK"   // simulated in-process (sandbox, until the partner's API is final)
	ModeREST   = "REST"   // the integration specification's REST contract
	ModeManual = "MANUAL" // operations act on the partner's portal and record it in Trovo Manager
)

// Custodian is the Public Markets integration of an existing Approved Asset
// Custodian (approved_asset_custodians); CSCS is reached only through it.
type Custodian struct {
	CustodianID uint64    `gorm:"primaryKey;autoIncrement:false" json:"custodianId"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Code        string    `gorm:"size:20;uniqueIndex;not null" json:"code"` // sent as X-Partner-Code on callbacks
	NomineeName string    `gorm:"size:150" json:"nomineeName"`
	Mode        string    `gorm:"size:10;not null;default:'MOCK'" json:"mode"`
	Transport   string    `gorm:"size:60" json:"transport"`  // label, e.g. "REST + webhooks", "SFTP (end of day file)"
	AuthScheme  string    `gorm:"size:20" json:"authScheme"` // HMAC | MTLS | NONE
	BaseURL     string    `gorm:"size:300" json:"baseUrl"`
	// CredentialsRef names where the partner's credentials live: an
	// environment variable (env:NAME) or a Vault path (vault://...).
	CredentialsRef string `gorm:"size:200" json:"credentialsRef"`
	FeePercent     string `gorm:"size:20" json:"feePercent"`
	Active         bool   `gorm:"not null;default:true" json:"active"`
}

func (Custodian) TableName() string { return "public_market_custodians" }

// DealingMember is a licensed broker that executes net buys and sells on
// NGX/FMDQ (doc: ApprovedDealingMember).
type DealingMember struct {
	ID                   uint64    `gorm:"primaryKey" json:"id"`
	CreatedAt            time.Time `json:"createdAt"`
	UpdatedAt            time.Time `json:"updatedAt"`
	DealingMemberName    string    `gorm:"size:100;not null" json:"dealingMemberName"`
	DealingMemberAddress string    `json:"dealingMemberAddress"`
	DealingMemberCountry string    `gorm:"size:3;not null" json:"dealingMemberCountry"`
	CSCSMemberCode       string    `gorm:"size:20" json:"cscsMemberCode"`
	RequirementDocument  string    `json:"requirementDocument"`
	FeePercent           float64   `gorm:"default:0" json:"feePercent"`
	FeeFixed             float64   `gorm:"default:0" json:"feeFixed"`
	Code                 string    `gorm:"size:20;uniqueIndex;not null" json:"code"`
	Mode                 string    `gorm:"size:10;not null;default:'MOCK'" json:"mode"`
	AuthScheme           string    `gorm:"size:20" json:"authScheme"`
	BaseURL              string    `gorm:"size:300" json:"baseUrl"`
	CredentialsRef       string    `gorm:"size:200" json:"credentialsRef"`
	Active               bool      `gorm:"not null;default:true" json:"active"`
}

func (DealingMember) TableName() string { return "approved_dealing_members" }

// Position sources.
const (
	PositionFeed       = "FEED"       // the Custodian's daily position feed (push)
	PositionPull       = "PULL"       // on-demand position query
	PositionSettlement = "SETTLEMENT" // previous position plus a settled instruction
	PositionManual     = "MANUAL"     // recorded by operations
)

// CustodianPosition is the real number of units the Custodian holds for the
// pool at a point in time; the latest row is the current position.
type CustodianPosition struct {
	ID            uint64    `gorm:"primaryKey" json:"id"`
	CreatedAt     time.Time `json:"createdAt"`
	AssetID       string    `gorm:"size:64;not null;index" json:"assetId"`
	RealUnitsHeld string    `gorm:"size:40;not null" json:"realUnitsHeld"`
	AsOf          time.Time `json:"asOf"`
	Source        string    `gorm:"size:12;not null" json:"source"`
	Reference     string    `gorm:"size:120" json:"reference"`
	RecordedBy    string    `gorm:"size:150" json:"recordedBy"`
}

func (CustodianPosition) TableName() string { return "custodian_positions" }

// Channels an order or holding comes from.
const (
	ChannelApp      = "TROVO_APP"
	ChannelExchange = "EXCHANGE"
	ChannelPlatform = "PLATFORM" // issuing Safe, treasury
)

// LedgerEntry is who owns what: one row per wallet per asset (doc:
// BeneficialOwnershipLedger). It follows the token's on-chain Transfers,
// so P2P trades re-attribute ownership too.
type LedgerEntry struct {
	ID            uint64    `gorm:"primaryKey" json:"id"`
	UpdatedAt     time.Time `json:"updatedAt"`
	AssetID       string    `gorm:"size:64;not null;uniqueIndex:idx_pm_ledger_asset_wallet" json:"assetId"`
	WalletAddress string    `gorm:"size:64;not null;uniqueIndex:idx_pm_ledger_asset_wallet" json:"walletAddress"`
	Balance       string    `gorm:"size:80;not null;default:'0'" json:"balance"` // base units of the token
	Channel       string    `gorm:"size:12" json:"channel"`
	ServiceLinkID string    `gorm:"size:100;index" json:"serviceLinkId"`
}

func (LedgerEntry) TableName() string { return "beneficial_ownership_ledgers" }

// LedgerMovement is one Transfer of a token as it affects one wallet; the
// record-date snapshot of a dividend sums them up to the record block.
type LedgerMovement struct {
	ID            uint64 `gorm:"primaryKey"`
	AssetID       string `gorm:"size:64;not null;index:idx_pm_mov_asset_block"`
	WalletAddress string `gorm:"size:64;not null;index"`
	Delta         string `gorm:"size:80;not null"` // signed base units
	BlockNumber   uint64 `gorm:"not null;index:idx_pm_mov_asset_block"`
	TxHash        string `gorm:"size:70;not null;uniqueIndex:idx_pm_mov_unique"`
	LogIndex      uint   `gorm:"not null;uniqueIndex:idx_pm_mov_unique"`
	Side          string `gorm:"size:4;not null;uniqueIndex:idx_pm_mov_unique"` // IN | OUT
}

func (LedgerMovement) TableName() string { return "beneficial_ownership_movements" }

// Price sources.
const (
	PriceNGXFeed   = "NGX_FEED"
	PriceFMDQFeed  = "FMDQ_FEED"
	PriceSynthetic = "LAST_CLOSE_SYNTHETIC" // after hours: the last close, labelled
	PriceManual    = "MANUAL"
	PriceMockFeed  = "MOCK_FEED"
)

// PriceSnapshot is a captured reference price (doc: PriceOracleSnapshot).
type PriceSnapshot struct {
	ID          uint64    `gorm:"primaryKey" json:"id"`
	AssetID     string    `gorm:"size:64;not null;index:idx_pm_price_asset_time" json:"assetId"`
	Price       string    `gorm:"size:40;not null" json:"price"`
	Source      string    `gorm:"size:30;not null" json:"source"`
	MarketHours bool      `json:"marketHours"`
	AsOf        time.Time `json:"asOf"` // the vendor's time
	CapturedAt  time.Time `gorm:"index:idx_pm_price_asset_time" json:"capturedAt"`
}

func (PriceSnapshot) TableName() string { return "price_oracle_snapshots" }

// Order types, paths and states (states follow the architecture's model;
// an order's actual sequence is in its events).
const (
	OrderCreation   = "CREATION"
	OrderRedemption = "REDEMPTION"

	PathFast   = "FAST"   // filled from Custodian inventory
	PathSlow   = "SLOW"   // waits for the next session's net instruction
	PathNetted = "NETTED" // a redemption absorbed by inventory and treasury cash

	StateAwaitingPayment     = "awaiting-payment" // app: the wallet operation is not mined yet
	StateQueued              = "queued"
	StateFilledFromInventory = "filled-from-inventory"
	StatePendingExecution    = "pending-execution"
	StateExecuted            = "executed" // the Dealing Member filled the net instruction
	StateSubmitted           = "submitted"
	StateChainFinal          = "chain_final"
	StateSettlementFinal     = "settlement_final"
	StateComplete            = "complete"
	StateRejected            = "rejected"
	StateFailed              = "failed"
	StateCancelled           = "cancelled"
)

// Order is a creation (buy) or redemption (sell) from any channel; every
// channel runs through the same engine.
type Order struct {
	ID               string    `gorm:"primaryKey;size:40" json:"id"` // PM-CR-..., PM-RD-...
	CreatedAt        time.Time `gorm:"index" json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
	Type             string    `gorm:"size:12;not null;index" json:"type"`
	AssetID          string    `gorm:"size:64;not null;index" json:"assetId"`
	AssetCode        string    `gorm:"size:16;not null" json:"assetCode"`
	Channel          string    `gorm:"size:12;not null;index" json:"channel"`
	ServiceLinkID    string    `gorm:"size:100;index" json:"serviceLinkId"`
	ExternalOrderRef string    `gorm:"size:120;index" json:"externalOrderRef"`
	Username         string    `gorm:"size:100;index" json:"username"`
	WalletAddress    string    `gorm:"size:64;not null;index" json:"walletAddress"`
	WalletAlias      string    `gorm:"size:100" json:"walletAlias"`
	PartnerWalletID  string    `gorm:"size:40;index" json:"partnerWalletId"` // exchange-provisioned wallet
	// payment (creation) / proceeds (redemption), in the funding currency
	FundingAssetCode string `gorm:"size:12" json:"fundingAssetCode"`
	FundingContract  string `gorm:"size:64" json:"fundingContract"`
	Amount           string `gorm:"size:40;not null;default:'0'" json:"amount"`    // creation: what the buyer pays (incl. fee); redemption: gross proceeds
	Fee              string `gorm:"size:40;not null;default:'0'" json:"fee"`       // Trovo fee
	NetAmount        string `gorm:"size:40;not null;default:'0'" json:"netAmount"` // creation: invested; redemption: paid out
	FeePercent       string `gorm:"size:20" json:"feePercent"`
	Quantity         string `gorm:"size:40;not null;default:'0'" json:"quantity"` // tokens (estimated until settled for slow creations)
	ReferencePrice   string `gorm:"size:40" json:"referencePrice"`
	PriceSource      string `gorm:"size:30" json:"priceSource"`
	ExecutedPrice    string `gorm:"size:40" json:"executedPrice"`
	Path             string `gorm:"size:8" json:"path"`
	State            string `gorm:"size:24;not null;index" json:"state"`
	Note             string `gorm:"size:500" json:"note"`
	BatchID          string `gorm:"size:60;index" json:"batchId"`
	// on-chain trail
	PaymentOperationID string     `gorm:"size:70;index" json:"paymentOperationId"` // the app wallet's operation
	PaymentTxHash      string     `gorm:"size:70" json:"paymentTxHash"`
	TokenTxHash        string     `gorm:"size:70" json:"tokenTxHash"` // mint (creation) or burn (redemption)
	PayoutTxHash       string     `gorm:"size:70" json:"payoutTxHash"`
	PaymentConfirmedAt *time.Time `json:"paymentConfirmedAt"`
	SettledAt          *time.Time `json:"settledAt"`
	CompletedAt        *time.Time `json:"completedAt"`
	FeeSweepTx         string     `gorm:"size:70;index" json:"-"` // the treasury transfer that moved the fee to the fee wallet
	IdempotencyKey     string     `gorm:"size:120" json:"-"`
	Attempts           int        `gorm:"not null;default:0" json:"-"`
	LastError          string     `gorm:"size:500" json:"-"`
}

func (Order) TableName() string { return "public_market_orders" }

// OrderEvent is one step of an order's timeline.
type OrderEvent struct {
	ID        uint64    `gorm:"primaryKey" json:"-"`
	CreatedAt time.Time `json:"at"`
	OrderID   string    `gorm:"size:40;not null;index" json:"-"`
	State     string    `gorm:"size:24;not null" json:"state"`
	Note      string    `gorm:"size:500" json:"note"`
}

func (OrderEvent) TableName() string { return "public_market_order_events" }

// Net batch statuses.
const (
	BatchAwaitingApproval = "AWAITING_APPROVAL"
	BatchReleased         = "RELEASED"  // instructions dispatched
	BatchExecuted         = "EXECUTED"  // Dealing Member filled
	BatchSettled          = "SETTLED"   // Custodian settlement_final
	BatchProcessed        = "PROCESSED" // its orders minted / paid
	BatchRejected         = "REJECTED"
	BatchFailed           = "FAILED"
	BatchInternal         = "INTERNAL" // fully netted: nothing to send to the market
)

// NetBatch is one session's net instruction for an asset: the unmatched
// net of its slow-path orders, sent once to the Dealing Member and the
// Custodian rather than once per order.
type NetBatch struct {
	ID                string     `gorm:"primaryKey;size:60" json:"id"` // NB-2026-09-29-1015-MTNN-T
	CreatedAt         time.Time  `json:"createdAt"`
	UpdatedAt         time.Time  `json:"updatedAt"`
	AssetID           string     `gorm:"size:64;not null;index" json:"assetId"`
	AssetCode         string     `gorm:"size:16;not null" json:"assetCode"`
	SessionDate       string     `gorm:"size:10;not null;index" json:"sessionDate"`
	Side              string     `gorm:"size:4;not null" json:"side"` // BUY | SELL | NONE
	Quantity          string     `gorm:"size:40;not null" json:"quantity"`
	ReferencePrice    string     `gorm:"size:40;not null" json:"referencePrice"`
	Value             string     `gorm:"size:40;not null" json:"value"`
	CreationOrders    int        `json:"creationOrders"`
	RedemptionOrders  int        `json:"redemptionOrders"`
	Status            string     `gorm:"size:20;not null;index" json:"status"`
	ApprovalsRequired int        `json:"approvalsRequired"`
	ExecutedQuantity  string     `gorm:"size:40" json:"executedQuantity"`
	ExecutedPrice     string     `gorm:"size:40" json:"executedPrice"`
	ExecutedAt        *time.Time `json:"executedAt"`
	SettledAt         *time.Time `json:"settledAt"`
	Note              string     `gorm:"size:500" json:"note"`
}

func (NetBatch) TableName() string { return "public_market_net_batches" }

// BatchApproval is one admin's approval of a net batch above the threshold.
type BatchApproval struct {
	ID        uint64    `gorm:"primaryKey" json:"-"`
	CreatedAt time.Time `json:"at"`
	BatchID   string    `gorm:"size:60;not null;uniqueIndex:idx_pm_batch_approver" json:"batchId"`
	Approver  string    `gorm:"size:150;not null;uniqueIndex:idx_pm_batch_approver" json:"approver"`
}

func (BatchApproval) TableName() string { return "public_market_batch_approvals" }

// Instruction kinds and statuses (the durable outbound queue).
const (
	InstrCustodianCreation   = "CUSTODIAN_CREATION"
	InstrCustodianRedemption = "CUSTODIAN_REDEMPTION"
	InstrDealingOrder        = "DM_ORDER"

	InstrPending   = "PENDING"   // waiting to be sent (or retried)
	InstrAccepted  = "ACCEPTED"  // the partner acknowledged it (202)
	InstrExecuted  = "EXECUTED"  // DM: FILLED
	InstrSettled   = "SETTLED"   // Custodian: settlement_final
	InstrEscalated = "ESCALATED" // retry cap reached: Operations handle it
	InstrHandled   = "HANDLED"   // closed manually by Operations
	InstrRejected  = "REJECTED"  // the partner refused it
)

// Instruction is an outbound instruction to a regulated counterparty. It is
// never dropped: failures retry with backoff, then escalate.
type Instruction struct {
	ID             string     `gorm:"primaryKey;size:64" json:"id"` // also the Idempotency-Key and instructionId/orderId
	CreatedAt      time.Time  `json:"createdAt"`
	UpdatedAt      time.Time  `json:"updatedAt"`
	Kind           string     `gorm:"size:24;not null;index" json:"kind"`
	PartnerType    string     `gorm:"size:16;not null" json:"partnerType"` // CUSTODIAN | DEALING_MEMBER
	PartnerID      uint64     `gorm:"not null;index" json:"partnerId"`
	PartnerName    string     `gorm:"size:150" json:"partnerName"`
	AssetID        string     `gorm:"size:64;not null;index" json:"assetId"`
	AssetCode      string     `gorm:"size:16" json:"assetCode"`
	BatchID        string     `gorm:"size:60;index" json:"batchId"`
	Side           string     `gorm:"size:4;not null" json:"side"`
	Quantity       string     `gorm:"size:40;not null" json:"quantity"`
	Payload        string     `gorm:"type:text" json:"payload"`
	Status         string     `gorm:"size:12;not null;index" json:"status"`
	Attempts       int        `gorm:"not null;default:0" json:"attempts"`
	MaxAttempts    int        `gorm:"not null;default:5" json:"maxAttempts"`
	NextAttemptAt  *time.Time `gorm:"index" json:"nextAttemptAt"`
	LastError      string     `gorm:"size:500" json:"lastError"`
	LastStatusCode int        `json:"lastStatusCode"`
	PartnerRef     string     `gorm:"size:120" json:"partnerRef"` // custodianReference
	AcceptedAt     *time.Time `json:"acceptedAt"`
	ExecutedAt     *time.Time `json:"executedAt"`
	SettledAt      *time.Time `json:"settledAt"`
	HandledBy      string     `gorm:"size:150" json:"handledBy"`
}

func (Instruction) TableName() string { return "public_market_instructions" }

// PartnerEvent is an inbound webhook, kept for idempotency (one row per
// partner reference) and for the audit trail.
type PartnerEvent struct {
	ID          uint64     `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time  `json:"createdAt"`
	Source      string     `gorm:"size:40;not null;uniqueIndex:idx_pm_partner_event" json:"source"` // CUSTODIAN:<code>, DEALING_MEMBER:<code>, EXCHANGE:<id>
	Kind        string     `gorm:"size:40;not null;uniqueIndex:idx_pm_partner_event" json:"kind"`
	Reference   string     `gorm:"size:160;not null;uniqueIndex:idx_pm_partner_event" json:"reference"`
	Payload     string     `gorm:"type:text" json:"payload"`
	ProcessedAt *time.Time `json:"processedAt"`
	Result      string     `gorm:"size:500" json:"result"`
}

func (PartnerEvent) TableName() string { return "public_market_partner_events" }

// PartnerWallet is the audit trail of what an exchange told Trovotech about
// one of its customers, and the individual wallet opened for them (doc:
// WalletProvisioningRequest). Rejected requests are kept too.
type PartnerWallet struct {
	ID               string     `gorm:"primaryKey;size:40" json:"walletId"` // wlt_...
	CreatedAt        time.Time  `json:"createdAt"`
	UpdatedAt        time.Time  `json:"updatedAt"`
	ServiceLinkID    string     `gorm:"size:100;not null;index:idx_pm_pw_ref" json:"serviceLinkId"`
	ExternalUserRef  string     `gorm:"size:120;not null;index:idx_pm_pw_ref" json:"externalUserRef"`
	LegalName        string     `gorm:"size:200" json:"legalName"`
	TaxIdentifier    string     `gorm:"size:60" json:"taxIdentifier"`
	ResidencyCountry string     `gorm:"size:3" json:"residencyCountry"`
	Nationality      string     `gorm:"size:3" json:"nationality"`
	NDPAConsent      bool       `json:"ndpaConsent"`
	WalletAddress    string     `gorm:"size:64;index" json:"walletAddress"` // a Safe owned by the Public Markets signers
	SafeSaltNonce    string     `gorm:"size:80" json:"-"`
	Deployed         bool       `json:"deployed"`
	Status           string     `gorm:"size:12;not null" json:"status"` // active | rejected | suspended
	RejectionReason  string     `gorm:"size:300" json:"rejectionReason"`
	DeployedAt       *time.Time `json:"deployedAt"`
}

func (PartnerWallet) TableName() string { return "wallet_provisioning_requests" }

// IdempotencyRecord remembers an exchange request by its Idempotency-Key
// so a retry returns the original result.
type IdempotencyRecord struct {
	ID            uint64    `gorm:"primaryKey"`
	CreatedAt     time.Time `gorm:"index"`
	ServiceLinkID string    `gorm:"size:100;not null;uniqueIndex:idx_pm_idem"`
	Key           string    `gorm:"size:160;not null;uniqueIndex:idx_pm_idem"`
	Path          string    `gorm:"size:200"`
	RequestHash   string    `gorm:"size:70"`
	Status        int
	Body          string `gorm:"type:text"`
	Done          bool
}

func (IdempotencyRecord) TableName() string { return "public_market_idempotency_keys" }

// Webhook delivery statuses.
const (
	DeliveryPending    = "PENDING"
	DeliveryDelivered  = "DELIVERED"
	DeliveryDeadLetter = "DEAD_LETTER"
)

// WebhookDelivery is an outbound notification to an exchange: at least
// once, with backoff, then dead-lettered (and replayable from Trovo
// Manager). dividend.paid also needs the exchange's confirmation.
type WebhookDelivery struct {
	ID                string     `gorm:"primaryKey;size:40" json:"id"` // evt_...
	CreatedAt         time.Time  `gorm:"index" json:"createdAt"`
	UpdatedAt         time.Time  `json:"updatedAt"`
	ServiceLinkID     string     `gorm:"size:100;not null;index" json:"serviceLinkId"`
	Event             string     `gorm:"size:40;not null;index" json:"event"`
	Reference         string     `gorm:"size:120;index" json:"reference"` // order id / entitlement id
	PartnerWalletID   string     `gorm:"size:40;index" json:"walletId"`
	AssetCode         string     `gorm:"size:16" json:"assetCode"`
	Payload           string     `gorm:"type:text" json:"payload"`
	Status            string     `gorm:"size:12;not null;index" json:"status"`
	Attempts          int        `json:"attempts"`
	NextAttemptAt     *time.Time `gorm:"index" json:"nextAttemptAt"`
	LastResponseCode  int        `json:"lastResponseCode"`
	LastError         string     `gorm:"size:300" json:"lastError"`
	DeliveredAt       *time.Time `json:"deliveredAt"`
	NeedsConfirmation bool       `gorm:"index" json:"needsConfirmation"`
	ConfirmedAt       *time.Time `json:"confirmedAt"`
	ConfirmationDueAt *time.Time `json:"confirmationDueAt"`
	EscalatedAt       *time.Time `json:"escalatedAt"`
}

func (WebhookDelivery) TableName() string { return "public_market_webhook_deliveries" }

// ExchangePartner is the Public Markets configuration of an exchange's
// service link (the exchange calls /v1/trovo-api/public-markets/... with
// its service link API key and signs every request).
type ExchangePartner struct {
	ServiceLinkID string    `gorm:"primaryKey;size:100" json:"serviceLinkId"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
	Status        string    `gorm:"size:12;not null;default:'active'" json:"status"` // active | suspended
	Environment   string    `gorm:"size:12;not null;default:'sandbox'" json:"environment"`
	CallbackURL   string    `gorm:"size:300" json:"callbackUrl"`
	// SigningSecret signs requests in both directions (HMAC-SHA256). After a
	// rotation the previous one keeps working until PreviousSecretExpiresAt.
	SigningSecret           string     `gorm:"size:100" json:"-"`
	PreviousSigningSecret   string     `gorm:"size:100" json:"-"`
	PreviousSecretExpiresAt *time.Time `json:"-"`
	RateLimitTier           string     `gorm:"size:20" json:"rateLimitTier"`
	RevenueShareTier        string     `gorm:"size:20" json:"revenueShareTier"`
	ConfirmationSLAHours    int        `gorm:"not null;default:0" json:"confirmationSlaHours"` // 0: the settings' SLA
	// the exchange's prefunded CNGN balance with Trovotech: creation orders
	// are paid from it, redemption proceeds and dividends credited to it
	FundingAddress string `gorm:"size:64" json:"fundingAddress"` // where its deposits come from
	Balance        string `gorm:"size:40;not null;default:'0'" json:"balance"`
	Reserved       string `gorm:"size:40;not null;default:'0'" json:"reserved"`
	TechContact    string `gorm:"size:150" json:"techContact"`
}

func (ExchangePartner) TableName() string { return "public_market_exchange_partners" }

// ExchangeLedgerEntry is a movement on an exchange's prefunded balance.
type ExchangeLedgerEntry struct {
	ID            uint64    `gorm:"primaryKey" json:"id"`
	CreatedAt     time.Time `json:"createdAt"`
	ServiceLinkID string    `gorm:"size:100;not null;index" json:"serviceLinkId"`
	Kind          string    `gorm:"size:20;not null" json:"kind"`   // DEPOSIT | ORDER | REFUND | REDEMPTION | DIVIDEND | WITHDRAWAL | ADJUSTMENT
	Amount        string    `gorm:"size:40;not null" json:"amount"` // signed
	BalanceAfter  string    `gorm:"size:40;not null" json:"balanceAfter"`
	Reference     string    `gorm:"size:120;uniqueIndex:idx_pm_exl_ref" json:"reference"`
	Note          string    `gorm:"size:300" json:"note"`
	CreatedBy     string    `gorm:"size:150" json:"createdBy"`
}

func (ExchangeLedgerEntry) TableName() string { return "public_market_exchange_ledger" }

// Corporate action types and statuses.
const (
	ActionDividend = "DIVIDEND"
	ActionCoupon   = "COUPON"
	ActionBonus    = "BONUS"
	ActionRights   = "RIGHTS"
	ActionSplit    = "SPLIT"

	ActionAnnounced     = "ANNOUNCED"    // waiting for the record date
	ActionSnapshotted   = "SNAPSHOTTED"  // entitlements fixed, awaiting approval
	ActionApproved      = "APPROVED"     // approved, waiting for funding
	ActionPaying        = "PAYING"       // payments under way
	ActionDistributed   = "DISTRIBUTED"  // all entitlements settled
	ActionNeedsManual   = "NEEDS_MANUAL" // supply-changing events: not automated
	ActionCancelledCA   = "CANCELLED"
	SourceCustodianCA   = "CUSTODIAN_WEBHOOK"
	SourceManualCA      = "MANUAL"
	CashActionFundToken = "CNGN"
)

// CorporateAction is a dividend, coupon or other event on an asset, from a
// Custodian webhook or declared in Trovo Manager. Public Markets runs its
// own distribution: a record-date snapshot, per-owner withholding tax and
// payment from the treasury.
type CorporateAction struct {
	ID                string     `gorm:"primaryKey;size:40" json:"id"` // CA-...
	CreatedAt         time.Time  `json:"createdAt"`
	UpdatedAt         time.Time  `json:"updatedAt"`
	AssetID           string     `gorm:"size:64;not null;index" json:"assetId"`
	AssetCode         string     `gorm:"size:16;not null" json:"assetCode"`
	EventType         string     `gorm:"size:12;not null" json:"eventType"`
	Description       string     `gorm:"size:300" json:"description"`
	RecordDate        string     `gorm:"size:10;not null" json:"recordDate"` // YYYY-MM-DD (Africa/Lagos)
	PayDate           string     `gorm:"size:10" json:"payDate"`
	AmountPerUnit     string     `gorm:"size:40" json:"amountPerUnit"`
	Currency          string     `gorm:"size:5" json:"currency"`
	Source            string     `gorm:"size:20;not null" json:"source"`
	SourceReference   string     `gorm:"size:160" json:"sourceReference"`
	Status            string     `gorm:"size:16;not null;index" json:"status"`
	DeclaredBy        string     `gorm:"size:150" json:"declaredBy"`
	RecordBlock       uint64     `json:"recordBlock"`
	SnapshotAt        *time.Time `json:"snapshotAt"`
	SnapshotChecksum  string     `gorm:"size:70" json:"snapshotChecksum"`
	HolderCount       int        `json:"holderCount"`
	EligibleUnits     string     `gorm:"size:40" json:"eligibleUnits"`
	RetainedUnits     string     `gorm:"size:40" json:"retainedUnits"` // platform wallets (inventory): not paid out
	GrossAmount       string     `gorm:"size:40" json:"grossAmount"`
	WHTAmount         string     `gorm:"size:40" json:"whtAmount"`
	NetAmount         string     `gorm:"size:40" json:"netAmount"`
	ApprovalsRequired int        `json:"approvalsRequired"`
	ApprovedAt        *time.Time `json:"approvedAt"`
	FundedAt          *time.Time `json:"fundedAt"`
	WHTTxHash         string     `gorm:"size:70" json:"whtTxHash"`
	CompletedAt       *time.Time `json:"completedAt"`
	Note              string     `gorm:"size:500" json:"note"`
}

func (CorporateAction) TableName() string { return "public_market_corporate_actions" }

// Entitlement statuses.
const (
	EntitlementPending  = "PENDING"
	EntitlementQueued   = "QUEUED" // in a sent batch
	EntitlementPaid     = "PAID"
	EntitlementFailed   = "FAILED"
	EntitlementRetained = "RETAINED" // platform wallets
)

// Entitlement is one holder's share of a corporate action, fixed at the
// record date and never changed by later balance moves.
type Entitlement struct {
	ID                string     `gorm:"primaryKey;size:40" json:"id"`
	CreatedAt         time.Time  `json:"createdAt"`
	UpdatedAt         time.Time  `json:"updatedAt"`
	CorporateActionID string     `gorm:"size:40;not null;index" json:"corporateActionId"`
	AssetID           string     `gorm:"size:64;not null;index" json:"assetId"`
	AssetCode         string     `gorm:"size:16" json:"assetCode"`
	WalletAddress     string     `gorm:"size:64;not null;index" json:"walletAddress"`
	Channel           string     `gorm:"size:12" json:"channel"`
	ServiceLinkID     string     `gorm:"size:100;index" json:"serviceLinkId"`
	PartnerWalletID   string     `gorm:"size:40;index" json:"partnerWalletId"`
	Username          string     `gorm:"size:100;index" json:"username"`
	Units             string     `gorm:"size:40;not null" json:"units"` // tokens held at the record date
	GrossAmount       string     `gorm:"size:40;not null" json:"grossAmount"`
	WHTPercent        string     `gorm:"size:10;not null" json:"whtPercent"`
	WHTAmount         string     `gorm:"size:40;not null" json:"whtAmount"`
	NetAmount         string     `gorm:"size:40;not null" json:"netAmount"`
	TaxResidency      string     `gorm:"size:3" json:"taxResidency"`
	Status            string     `gorm:"size:12;not null;index" json:"status"`
	BatchTxHash       string     `gorm:"size:70" json:"txHash"`
	PaidAt            *time.Time `json:"paidAt"`
	Note              string     `gorm:"size:300" json:"note"`
}

func (Entitlement) TableName() string { return "public_market_dividend_entitlements" }

// DividendApproval is one admin's approval of a distribution (for its
// snapshot checksum).
type DividendApproval struct {
	ID                uint64    `gorm:"primaryKey" json:"-"`
	CreatedAt         time.Time `json:"at"`
	CorporateActionID string    `gorm:"size:40;not null;uniqueIndex:idx_pm_div_approver" json:"corporateActionId"`
	Approver          string    `gorm:"size:150;not null;uniqueIndex:idx_pm_div_approver" json:"approver"`
	Checksum          string    `gorm:"size:70" json:"checksum"`
}

func (DividendApproval) TableName() string { return "public_market_dividend_approvals" }

// Reconciliation results.
const (
	ReconMatched     = "MATCHED"
	ReconDrift       = "DRIFT"        // token supply above the Custodian position
	ReconLedgerDrift = "LEDGER_DRIFT" // Σ ledger differs from on-chain supply
	ReconNoPosition  = "NO_POSITION"
	ReconError       = "ERROR"
)

// ReconciliationRun is one asset's check of Σ tokens against the units the
// Custodian holds, and of the ledger against on-chain supply.
type ReconciliationRun struct {
	ID                uint64     `gorm:"primaryKey" json:"id"`
	CreatedAt         time.Time  `gorm:"index" json:"createdAt"`
	AssetID           string     `gorm:"size:64;not null;index" json:"assetId"`
	AssetCode         string     `gorm:"size:16" json:"assetCode"`
	TokenSupply       string     `gorm:"size:40" json:"tokenSupply"`
	CustodianPosition string     `gorm:"size:40" json:"custodianPosition"`
	PositionAsOf      *time.Time `json:"positionAsOf"`
	PositionSource    string     `gorm:"size:12" json:"positionSource"`
	LedgerTotal       string     `gorm:"size:40" json:"ledgerTotal"`
	Inventory         string     `gorm:"size:40" json:"inventory"` // position - supply: units held but not tokenized
	Delta             string     `gorm:"size:40" json:"delta"`     // supply - position when positive
	LedgerDelta       string     `gorm:"size:40" json:"ledgerDelta"`
	PriceStale        bool       `json:"priceStale"`
	Result            string     `gorm:"size:16;not null" json:"result"`
	Halted            bool       `json:"halted"`
	Detail            string     `gorm:"size:500" json:"detail"`
	TriggeredBy       string     `gorm:"size:150" json:"triggeredBy"`
}

func (ReconciliationRun) TableName() string { return "public_market_reconciliation_runs" }

// Settings is the single row of Public Markets thresholds and limits,
// edited in Trovo Manager.
type Settings struct {
	ID                        uint      `gorm:"primaryKey" json:"-"`
	UpdatedAt                 time.Time `json:"updatedAt"`
	UpdatedBy                 string    `gorm:"size:150" json:"updatedBy"`
	TradeFeePercent           string    `gorm:"size:20;not null;default:'0.25'" json:"tradeFeePercent"`
	FeeWallet                 string    `gorm:"size:64" json:"feeWallet"`
	FundingAssetCode          string    `gorm:"size:12;not null;default:'CNGN'" json:"fundingAssetCode"`
	NetCreationThreshold      string    `gorm:"size:40;not null;default:'15000000'" json:"netCreationThreshold"` // per asset per day, NGN
	ApprovalsRequired         int       `gorm:"not null;default:2" json:"approvalsRequired"`
	NetCreationApprovers      string    `gorm:"type:text" json:"netCreationApprovers"` // comma-separated admin emails
	DividendApprovers         string    `gorm:"type:text" json:"dividendApprovers"`
	PriceStaleMinutes         int       `gorm:"not null;default:15" json:"priceStaleMinutes"`
	InstructionMaxAttempts    int       `gorm:"not null;default:5" json:"instructionMaxAttempts"`
	SettlementSLAHours        int       `gorm:"not null;default:72" json:"settlementSlaHours"`
	ConfirmationSLAHours      int       `gorm:"not null;default:0" json:"confirmationSlaHours"` // 0 = not set (OI-CONS-10)
	WebhookMaxAttempts        int       `gorm:"not null;default:9" json:"webhookMaxAttempts"`
	BatchIntervalMinutes      int       `gorm:"not null;default:15" json:"batchIntervalMinutes"`
	ReconciliationHour        int       `gorm:"not null;default:6" json:"reconciliationHour"` // WAT
	WHTResidentPercent        string    `gorm:"size:10;not null;default:'10'" json:"whtResidentPercent"`
	WHTNonResidentPercent     string    `gorm:"size:10;not null;default:'10'" json:"whtNonResidentPercent"`
	WHTMissingTaxIDPercent    string    `gorm:"size:10;not null;default:'10'" json:"whtMissingTaxIdPercent"`
	WHTWallet                 string    `gorm:"size:64" json:"whtWallet"`
	SubstantialHoldingPercent string    `gorm:"size:10;not null;default:'5'" json:"substantialHoldingPercent"`
	RateLimitTiers            string    `gorm:"type:text" json:"rateLimitTiers"` // JSON {"Tier 1": 1200, ...} requests/minute
	AUMFeePercent             string    `gorm:"size:10;not null;default:'0.75'" json:"aumFeePercent"`
	FXSpreadPercent           string    `gorm:"size:10;not null;default:'0.50'" json:"fxSpreadPercent"`
	RevenueShareTiers         string    `gorm:"type:text" json:"revenueShareTiers"`
	NGXOpen                   string    `gorm:"size:5;not null;default:'10:00'" json:"ngxOpen"`
	NGXClose                  string    `gorm:"size:5;not null;default:'14:30'" json:"ngxClose"`
	FMDQOpen                  string    `gorm:"size:5;not null;default:'09:00'" json:"fmdqOpen"`
	FMDQClose                 string    `gorm:"size:5;not null;default:'16:00'" json:"fmdqClose"`
	MarketHolidays            string    `gorm:"type:text" json:"marketHolidays"` // comma-separated YYYY-MM-DD
}

func (Settings) TableName() string { return "public_market_settings" }

// JobRequest asks the engine to run something now (from Trovo Manager):
// RECONCILE (an asset code or ALL), RESUME (an asset code), POSITION_FEED
// (a mock Custodian's code) and EXCHANGE_WITHDRAWAL (serviceLinkId:amount).
type JobRequest struct {
	ID          uint64     `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time  `json:"createdAt"`
	Job         string     `gorm:"size:30;not null" json:"job"`
	Target      string     `gorm:"size:64" json:"target"`
	RequestedBy string     `gorm:"size:150" json:"requestedBy"`
	DoneAt      *time.Time `gorm:"index" json:"doneAt"`
	Result      string     `gorm:"size:500" json:"result"`
}

func (JobRequest) TableName() string { return "public_market_job_requests" }

// JobRun records the last run of each background job (Trovo Manager's
// System Health).
type JobRun struct {
	Job       string    `gorm:"primaryKey;size:40" json:"job"`
	LastRunAt time.Time `json:"lastRunAt"`
	Status    string    `gorm:"size:12" json:"status"` // up | degraded | down
	Detail    string    `gorm:"size:300" json:"detail"`
}

func (JobRun) TableName() string { return "public_market_job_runs" }

// MockEvent is a callback the mock Custodian / Dealing Member will deliver
// at DueAt; it goes through the same processing as a real webhook.
type MockEvent struct {
	ID          uint64 `gorm:"primaryKey"`
	CreatedAt   time.Time
	DueAt       time.Time  `gorm:"index"`
	Kind        string     `gorm:"size:30;not null"` // execution | settlement
	Partner     string     `gorm:"size:40"`
	Payload     string     `gorm:"type:text"`
	DeliveredAt *time.Time `gorm:"index"`
}

func (MockEvent) TableName() string { return "public_market_mock_events" }

// MockCustodianBook is the mock Custodian's own record of what it holds,
// the source of its position feed.
type MockCustodianBook struct {
	AssetCode string `gorm:"primaryKey;size:16"`
	UpdatedAt time.Time
	Units     string `gorm:"size:40;not null;default:'0'"`
}

func (MockCustodianBook) TableName() string { return "public_market_mock_custodian_books" }

// All returns every Public Markets model, for AutoMigrate.
func All() []interface{} {
	return []interface{}{
		&Asset{}, &Custodian{}, &DealingMember{}, &CustodianPosition{}, &LedgerEntry{}, &LedgerMovement{},
		&PriceSnapshot{}, &Order{}, &OrderEvent{}, &NetBatch{}, &BatchApproval{}, &Instruction{}, &PartnerEvent{},
		&PartnerWallet{}, &IdempotencyRecord{}, &WebhookDelivery{}, &ExchangePartner{}, &ExchangeLedgerEntry{},
		&CorporateAction{}, &Entitlement{}, &DividendApproval{}, &ReconciliationRun{}, &Settings{}, &JobRequest{},
		&JobRun{}, &MockEvent{}, &MockCustodianBook{},
	}
}
