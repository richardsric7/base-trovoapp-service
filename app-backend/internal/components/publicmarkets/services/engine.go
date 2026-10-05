package publicmarkets

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"strings"
	"sync"
	"time"

	"trovo-wallet-api/internal/aa"
	pm "trovo-wallet-api/internal/components/publicmarkets/models"
	"trovo-wallet-api/internal/components/publicmarkets/partners"
	userModels "trovo-wallet-api/internal/components/users/models"
	"trovo-wallet-api/internal/gnosissafe"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// Engine is the Public Markets engine. Every channel's orders go through
// the same Engine methods.
type Engine struct {
	DB    *gorm.DB
	Chain Chain // nil: no chain configured (database-only)
	GC    *sharedconfig.GlobalConfig
	Now   func() time.Time
	// Partners builds the clients of Custodians and Dealing Members
	// (mock, REST or manual); Prices is the price feed.
	Partners PartnerFactory
	Prices   partners.PriceFeed
	// Notify sends an app user a push notification (nil: none).
	Notify func(username, title, body string, data map[string]string)

	mu       sync.Mutex
	decimals map[common.Address]uint8
	kick     chan struct{}
}

// NewEngine builds the engine on the global configuration.
func NewEngine(gc *sharedconfig.GlobalConfig, chain Chain) *Engine {
	e := &Engine{DB: gc.DB, Chain: chain, GC: gc, Now: time.Now, decimals: map[common.Address]uint8{}}
	e.Partners = DefaultPartners{E: e}
	e.Prices = PriceFeedFromEnv(e)
	return e
}

func (e *Engine) now() time.Time {
	if e.Now != nil {
		return e.Now().UTC()
	}
	return time.Now().UTC()
}

// ---------------------------------------------------------------- assets

// AssetByCode returns an asset by its code (MTNN-T) or ticker (MTNN).
func (e *Engine) AssetByCode(code string) (*pm.Asset, error) {
	code = strings.ToUpper(strings.TrimSpace(code))
	var a pm.Asset
	if err := e.DB.Where("asset_code = ? OR ticker = ? OR id = ?", code, code, code).First(&a).Error; err != nil {
		return nil, refuse(http.StatusNotFound, "error-asset-not-found", "assetCode", "No Public Markets asset %q.", code)
	}
	return &a, nil
}

func (e *Engine) tokenDecimals(ctx context.Context, token common.Address) (uint8, error) {
	e.mu.Lock()
	if v, ok := e.decimals[token]; ok {
		e.mu.Unlock()
		return v, nil
	}
	e.mu.Unlock()
	if e.Chain == nil {
		return 0, errors.New("no chain configured")
	}
	v, err := e.Chain.Decimals(ctx, token)
	if err != nil {
		return 0, err
	}
	e.mu.Lock()
	e.decimals[token] = v
	e.mu.Unlock()
	return v, nil
}

// FundingToken is the stablecoin orders are paid in (CNGN by default).
type FundingToken struct {
	Code     string
	Address  common.Address
	Decimals int
}

func (e *Engine) fundingToken(ctx context.Context) (FundingToken, error) {
	s := LoadSettings(e.DB)
	code := strings.ToUpper(strings.TrimSpace(s.FundingAssetCode))
	var cur userModels.TokenizationCurrency
	if err := e.DB.Where("asset_code = ?", code).First(&cur).Error; err != nil || !common.IsHexAddress(cur.ContractAddress) {
		return FundingToken{}, refuse(http.StatusServiceUnavailable, "error-funding-not-configured", "", "The funding currency %s has no token contract.", code)
	}
	ft := FundingToken{Code: code, Address: common.HexToAddress(cur.ContractAddress), Decimals: 6}
	if e.Chain != nil {
		dec, err := e.tokenDecimals(ctx, ft.Address)
		if err != nil {
			return ft, err
		}
		ft.Decimals = int(dec)
	}
	return ft, nil
}

// Price is an asset's current reference price and whether it can be
// traded on: fresh, and from a live feed.
type Price struct {
	Value      decimal.Decimal
	Source     string
	At         *time.Time
	Fresh      bool
	MarketOpen bool
	Session    Session
}

// CurrentPrice is the asset's reference price as of now.
func (e *Engine) CurrentPrice(a *pm.Asset) Price {
	s := LoadSettings(e.DB)
	sess := MarketSession(a.Market, e.now(), s)
	p := Price{Value: d(a.LastPrice), Source: a.PriceSource, At: a.PriceAt, MarketOpen: sess.Open, Session: sess}
	if a.PriceAt != nil && p.Value.IsPositive() {
		age := e.now().Sub(*a.PriceAt)
		p.Fresh = sess.Open && a.PriceSource != pm.PriceSynthetic && age <= time.Duration(s.PriceStaleMinutes)*time.Minute
	}
	return p
}

// Supply is the asset's token supply (human units): on-chain, or the
// ledger's total when there is no chain.
func (e *Engine) Supply(ctx context.Context, a *pm.Asset) (decimal.Decimal, error) {
	if e.Chain != nil && common.IsHexAddress(a.ContractAddress) {
		v, err := e.Chain.TotalSupply(ctx, common.HexToAddress(a.ContractAddress))
		if err != nil {
			return decimal.Zero, err
		}
		return decimal.NewFromBigInt(v, -int32(a.TokenDecimals)), nil
	}
	return e.LedgerTotal(a), nil
}

// LedgerTotal is Σ BeneficialOwnershipLedger for an asset (human units).
func (e *Engine) LedgerTotal(a *pm.Asset) decimal.Decimal {
	var rows []pm.LedgerEntry
	e.DB.Where("asset_id = ?", a.ID).Find(&rows)
	total := decimal.Zero
	for _, r := range rows {
		total = total.Add(d(r.Balance))
	}
	return total.Shift(-int32(a.TokenDecimals))
}

// Position is the latest Custodian position of an asset.
func (e *Engine) Position(a *pm.Asset) (pm.CustodianPosition, bool) {
	var p pm.CustodianPosition
	if e.DB.Where("asset_id = ?", a.ID).Order("id DESC").First(&p).Error != nil {
		return p, false
	}
	return p, true
}

// reservedUnits are tokens promised to creations that are backed at the
// Custodian but not yet minted.
func (e *Engine) reservedUnits(a *pm.Asset) decimal.Decimal {
	var orders []pm.Order
	e.DB.Where("asset_id = ? AND type = ? AND state IN ?", a.ID, pm.OrderCreation,
		[]string{pm.StateFilledFromInventory, pm.StateSettlementFinal, pm.StateSubmitted}).Find(&orders)
	total := decimal.Zero
	for _, o := range orders {
		if o.State == pm.StateSettlementFinal && o.TokenTxHash != "" {
			continue // minted already (fast path reaching settlement after chain_final)
		}
		total = total.Add(d(o.Quantity))
	}
	return total
}

// Inventory is what the Custodian holds beyond the token supply and the
// creations already promised: what a creation can be filled from at once.
func (e *Engine) Inventory(ctx context.Context, a *pm.Asset) (decimal.Decimal, error) {
	pos, ok := e.Position(a)
	if !ok {
		return decimal.Zero, nil
	}
	supply, err := e.Supply(ctx, a)
	if err != nil {
		return decimal.Zero, err
	}
	return d(pos.RealUnitsHeld).Sub(supply).Sub(e.reservedUnits(a)), nil
}

// TreasuryFreeCash is the funding stablecoin in the treasury that is not
// owed to anyone: exchange balances, unswept fees, pending redemption
// payouts and dividends, and slow creations' cash committed to the market
// all come off it.
func (e *Engine) TreasuryFreeCash(ctx context.Context) (decimal.Decimal, error) {
	if e.Chain == nil || e.Chain.TreasurySafe() == (common.Address{}) {
		return decimal.Zero, nil
	}
	ft, err := e.fundingToken(ctx)
	if err != nil {
		return decimal.Zero, err
	}
	bal, err := e.Chain.BalanceOf(ctx, ft.Address, e.Chain.TreasurySafe())
	if err != nil {
		return decimal.Zero, err
	}
	free := decimal.NewFromBigInt(bal, -int32(ft.Decimals))
	return free.Sub(e.obligations()), nil
}

func (e *Engine) obligations() decimal.Decimal {
	total := decimal.Zero
	var partnersRows []pm.ExchangePartner
	e.DB.Find(&partnersRows)
	for _, p := range partnersRows {
		total = total.Add(d(p.Balance))
	}
	var orders []pm.Order
	e.DB.Where("(state = ? AND fee_sweep_tx = '' AND fee <> '0') OR (type = ? AND state IN ? AND payout_tx_hash = '' AND channel = ?) OR (type = ? AND path = ? AND state IN ?)",
		pm.StateComplete,
		pm.OrderRedemption, []string{pm.StateSettlementFinal, pm.StateSubmitted, pm.StateChainFinal}, pm.ChannelApp,
		pm.OrderCreation, pm.PathSlow, []string{pm.StatePendingExecution, pm.StateExecuted, pm.StateSettlementFinal, pm.StateSubmitted, pm.StateChainFinal}).Find(&orders)
	for _, o := range orders {
		switch {
		case o.State == pm.StateComplete:
			total = total.Add(d(o.Fee))
		case o.Type == pm.OrderRedemption:
			total = total.Add(d(o.NetAmount))
		default:
			total = total.Add(d(o.NetAmount))
		}
	}
	var ents []pm.Entitlement
	e.DB.Joins("JOIN public_market_corporate_actions ca ON ca.id = public_market_dividend_entitlements.corporate_action_id").
		Where("ca.status IN ? AND public_market_dividend_entitlements.status IN ?", []string{pm.ActionApproved, pm.ActionPaying},
			[]string{pm.EntitlementPending, pm.EntitlementQueued, pm.EntitlementFailed}).Find(&ents)
	for _, en := range ents {
		total = total.Add(d(en.NetAmount)).Add(d(en.WHTAmount))
	}
	return total
}

// ---------------------------------------------------------------- quotes

// Quote prices a creation (by amount paid) or a redemption (by quantity)
// and says whether it fills now (fast or netted) or waits for the next
// session (slow).
type Quote struct {
	AssetCode      string     `json:"assetCode"`
	Side           string     `json:"side"` // BUY | SELL
	Amount         string     `json:"amount"`
	Fee            string     `json:"fee"`
	FeePercent     string     `json:"feePercent"`
	NetAmount      string     `json:"netAmount"`
	Quantity       string     `json:"quantity"`
	Price          string     `json:"price"`
	PriceSource    string     `json:"priceSource"`
	Currency       string     `json:"currency"`
	FundingAsset   string     `json:"fundingAsset"`
	MarketOpen     bool       `json:"marketOpen"`
	Path           string     `json:"path"`
	Note           string     `json:"note"`
	NextSessionAt  *time.Time `json:"nextSessionAt,omitempty"`
	MinimumBuy     string     `json:"minimumBuy"`
	CustodianName  string     `json:"custodianName"`
	SettlementNote string     `json:"settlementNote"`
}

func (e *Engine) custodianName(a *pm.Asset) string {
	var c userModels.ApprovedAssetCustodian
	if e.DB.First(&c, a.CustodianID).Error == nil {
		return c.AssetCustodianName
	}
	return "the Custodian"
}

func (e *Engine) checkTradable(a *pm.Asset) error {
	switch a.Status {
	case pm.AssetLive:
		return nil
	case pm.AssetHalted:
		return refuse(http.StatusConflict, "error-asset-halted", "assetCode", "%s is halted: new orders are not accepted until it is reconciled.", a.AssetCode)
	default:
		return refuse(http.StatusConflict, "error-asset-not-live", "assetCode", "%s is not open for orders yet.", a.AssetCode)
	}
}

// QuoteCreation prices buying with amount (in the funding currency, fee
// included).
func (e *Engine) QuoteCreation(ctx context.Context, a *pm.Asset, amount decimal.Decimal) (Quote, error) {
	s := LoadSettings(e.DB)
	if err := e.checkTradable(a); err != nil {
		return Quote{}, err
	}
	if !amount.IsPositive() {
		return Quote{}, refuse(http.StatusBadRequest, "error-invalid-amount", "amount", "The amount must be greater than zero.")
	}
	if min := d(a.MinimumBuy); min.IsPositive() && amount.LessThan(min) {
		return Quote{}, refuse(http.StatusBadRequest, "error-below-minimum", "amount", "The minimum buy is %s %s.", min.StringFixed(2), s.FundingAssetCode)
	}
	p := e.CurrentPrice(a)
	if !p.Value.IsPositive() {
		return Quote{}, refuse(http.StatusServiceUnavailable, "error-no-price", "assetCode", "%s has no reference price yet.", a.AssetCode)
	}
	feePct := FeePercent(a, s)
	fee := amount.Mul(feePct).Div(decimal.NewFromInt(100)).Round(2)
	net := amount.Sub(fee)
	qty := net.Div(p.Value).Truncate(int32(a.TokenDecimals))
	q := Quote{AssetCode: a.AssetCode, Side: "BUY", Amount: amount.String(), Fee: fee.String(), FeePercent: feePct.String(), NetAmount: net.String(),
		Quantity: qty.String(), Price: p.Value.String(), PriceSource: p.Source, Currency: a.Currency, FundingAsset: s.FundingAssetCode,
		MarketOpen: p.MarketOpen, MinimumBuy: a.MinimumBuy, CustodianName: e.custodianName(a)}
	q.Path, q.Note = pm.PathSlow, ""
	if p.Fresh {
		if inv, err := e.Inventory(ctx, a); err == nil && inv.GreaterThanOrEqual(qty) {
			q.Path = pm.PathFast
			q.Note = fmt.Sprintf("Fills instantly. %s already holds enough %s, so your tokens arrive as soon as you authorize.", q.CustodianName, a.Ticker)
		}
	}
	if q.Path == pm.PathSlow {
		at := p.Session.OpenAt
		q.NextSessionAt = &at
		if p.MarketOpen {
			q.Note = fmt.Sprintf("Fills later today. Your order joins the next %s net order to the broker; tokens arrive once the Custodian confirms settlement. %s is held until then.", a.Market, amount.StringFixed(2)+" "+s.FundingAssetCode)
		} else {
			q.Note = fmt.Sprintf("%s is closed. Your order is queued and fills at the next session (from %s WAT) at that session's price. %s is held until then.", a.Market, at.In(Lagos).Format("Mon 3:04 PM"), amount.StringFixed(2)+" "+s.FundingAssetCode)
		}
	}
	q.SettlementNote = "Instant from inventory, otherwise T+2"
	return q, nil
}

// QuoteRedemption prices selling quantity tokens.
func (e *Engine) QuoteRedemption(ctx context.Context, a *pm.Asset, qty decimal.Decimal) (Quote, error) {
	s := LoadSettings(e.DB)
	if err := e.checkTradable(a); err != nil {
		return Quote{}, err
	}
	qty = qty.Truncate(int32(a.TokenDecimals))
	if !qty.IsPositive() {
		return Quote{}, refuse(http.StatusBadRequest, "error-invalid-quantity", "quantity", "The quantity must be greater than zero.")
	}
	p := e.CurrentPrice(a)
	if !p.Value.IsPositive() {
		return Quote{}, refuse(http.StatusServiceUnavailable, "error-no-price", "assetCode", "%s has no reference price yet.", a.AssetCode)
	}
	feePct := FeePercent(a, s)
	gross := qty.Mul(p.Value).Round(2)
	fee := gross.Mul(feePct).Div(decimal.NewFromInt(100)).Round(2)
	net := gross.Sub(fee)
	q := Quote{AssetCode: a.AssetCode, Side: "SELL", Amount: gross.String(), Fee: fee.String(), FeePercent: feePct.String(), NetAmount: net.String(),
		Quantity: qty.String(), Price: p.Value.String(), PriceSource: p.Source, Currency: a.Currency, FundingAsset: s.FundingAssetCode,
		MarketOpen: p.MarketOpen, MinimumBuy: a.MinimumBuy, CustodianName: e.custodianName(a), Path: pm.PathSlow,
		SettlementNote: "Instant when netted, else T+2"}
	if p.Fresh {
		if free, err := e.TreasuryFreeCash(ctx); err == nil && free.GreaterThanOrEqual(net) {
			q.Path = pm.PathNetted
			q.Note = fmt.Sprintf("Settles instantly. Today's buy demand covers this sale, so %s lands in your wallet right away.", s.FundingAssetCode)
		}
	}
	if q.Path == pm.PathSlow {
		at := p.Session.OpenAt
		q.NextSessionAt = &at
		if p.MarketOpen {
			q.Note = "Settles after the broker sells. Your tokens are held until the Custodian confirms the sale (T+2)."
		} else {
			q.Note = fmt.Sprintf("%s is closed. Your sale is queued for the next session (from %s WAT) at that session's price.", a.Market, at.In(Lagos).Format("Mon 3:04 PM"))
		}
	}
	return q, nil
}

// ---------------------------------------------------------------- orders

func newOrderID(prefix string) string {
	n, _ := rand.Int(rand.Reader, big.NewInt(90000000))
	return fmt.Sprintf("%s%08d", prefix, n.Int64()+10000000)
}

// OrderRequest is an order from any channel.
type OrderRequest struct {
	Asset              *pm.Asset
	Channel            string
	Username           string
	WalletAddress      string
	WalletAlias        string
	ServiceLinkID      string
	PartnerWalletID    string
	ExternalOrderRef   string
	Quote              Quote
	State              string // StateAwaitingPayment (app: payment in flight) or StateQueued
	PaymentOperationID string
	IdempotencyKey     string
}

// Place records an order from a quote. Every channel calls it.
func (e *Engine) Place(r OrderRequest) (*pm.Order, error) {
	prefix, typ := "PM-CR-", pm.OrderCreation
	if r.Quote.Side == "SELL" {
		prefix, typ = "PM-RD-", pm.OrderRedemption
	}
	if r.State == "" {
		r.State = pm.StateQueued
	}
	ft := LoadSettings(e.DB).FundingAssetCode
	o := &pm.Order{ID: newOrderID(prefix), Type: typ, AssetID: r.Asset.ID, AssetCode: r.Asset.AssetCode, Channel: r.Channel,
		ServiceLinkID: r.ServiceLinkID, ExternalOrderRef: r.ExternalOrderRef, Username: r.Username, WalletAddress: r.WalletAddress,
		WalletAlias: r.WalletAlias, PartnerWalletID: r.PartnerWalletID, FundingAssetCode: ft, Amount: r.Quote.Amount, Fee: r.Quote.Fee,
		NetAmount: r.Quote.NetAmount, FeePercent: r.Quote.FeePercent, Quantity: r.Quote.Quantity, ReferencePrice: r.Quote.Price,
		PriceSource: r.Quote.PriceSource, State: r.State, PaymentOperationID: r.PaymentOperationID, IdempotencyKey: r.IdempotencyKey}
	err := e.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(o).Error; err != nil {
			return err
		}
		note := fmt.Sprintf("%s %s %s via %s", strings.ToLower(o.Type), o.Quantity, o.AssetCode, channelLabel(o.Channel))
		if o.State == pm.StateAwaitingPayment {
			note = "Waiting for the wallet's payment to be confirmed on-chain"
		}
		return tx.Create(&pm.OrderEvent{OrderID: o.ID, State: o.State, Note: note}).Error
	})
	if err != nil {
		log.Printf("[publicmarkets] saving order: %v", err)
		return nil, refuse(http.StatusInternalServerError, "error-temporary", "", "The order could not be saved. Please try again.")
	}
	return o, nil
}

func channelLabel(c string) string {
	switch c {
	case pm.ChannelApp:
		return "Trovo App"
	case pm.ChannelExchange:
		return "exchange partner"
	}
	return c
}

// transition moves an order from one of `from` to `to` (atomically, so a
// concurrent move loses), records the event and returns whether it moved.
func (e *Engine) transition(o *pm.Order, from []string, to, note string, updates map[string]interface{}) bool {
	if updates == nil {
		updates = map[string]interface{}{}
	}
	updates["state"] = to
	updates["updated_at"] = e.now()
	moved := false
	e.DB.Transaction(func(tx *gorm.DB) error {
		res := tx.Model(&pm.Order{}).Where("id = ? AND state IN ?", o.ID, from).Updates(updates)
		if res.Error != nil || res.RowsAffected == 0 {
			return errors.New("not moved")
		}
		moved = true
		return tx.Create(&pm.OrderEvent{OrderID: o.ID, State: to, Note: note, CreatedAt: e.now()}).Error
	})
	if moved {
		e.DB.First(o, "id = ?", o.ID)
		e.onTransition(o)
	}
	return moved
}

// note records a step that does not change the state.
func (e *Engine) note(o *pm.Order, note string) {
	e.DB.Create(&pm.OrderEvent{OrderID: o.ID, State: o.State, Note: note, CreatedAt: e.now()})
}

// onTransition tells the order's owner: an app push, or an exchange
// webhook.
func (e *Engine) onTransition(o *pm.Order) {
	switch o.State {
	case pm.StateComplete, pm.StateRejected, pm.StateFailed:
	default:
		return
	}
	if o.Channel == pm.ChannelExchange {
		e.queueOrderWebhook(o)
		return
	}
	if e.Notify == nil || o.Username == "" {
		return
	}
	data := map[string]string{"route": "publicMarketsOrder", "orderId": o.ID, "assetCode": o.AssetCode}
	switch {
	case o.State == pm.StateComplete && o.Type == pm.OrderCreation:
		e.Notify(o.Username, "Purchase complete", fmt.Sprintf("You bought %s %s. It's in your wallet and in My Stocks.", d(o.Quantity).StringFixed(4), o.AssetCode), data)
	case o.State == pm.StateComplete:
		e.Notify(o.Username, "Sale complete", fmt.Sprintf("You sold %s %s. %s %s is in your wallet.", d(o.Quantity).String(), o.AssetCode, d(o.NetAmount).StringFixed(2), o.FundingAssetCode), data)
	default:
		e.Notify(o.Username, "Order not completed", fmt.Sprintf("Your %s order %s was not completed: %s", o.AssetCode, o.ID, o.Note), data)
	}
}

// ---------------------------------------------------------------- the order loop

// ProcessOrders moves every open order one step: payment confirmation,
// routing, minting/burning, payouts.
func (e *Engine) ProcessOrders(ctx context.Context) {
	var orders []pm.Order
	e.DB.Where("state IN ?", []string{pm.StateAwaitingPayment, pm.StateQueued, pm.StateFilledFromInventory, pm.StateSettlementFinal,
		pm.StateSubmitted, pm.StateChainFinal}).Order("created_at").Limit(500).Find(&orders)
	for i := range orders {
		if ctx.Err() != nil {
			return
		}
		o := &orders[i]
		if err := e.step(ctx, o); err != nil {
			e.fail(o, err)
		}
	}
}

// fail counts a failed step: after five, the order fails and Operations
// are alerted (a payment or token move that went through is never undone
// silently).
func (e *Engine) fail(o *pm.Order, err error) {
	log.Printf("[publicmarkets] order %s (%s): %v", o.ID, o.State, err)
	attempts := o.Attempts + 1
	e.DB.Model(&pm.Order{}).Where("id = ?", o.ID).Updates(map[string]interface{}{"attempts": attempts, "last_error": trimTo(err.Error(), 500)})
	if attempts >= 5 {
		if e.transition(o, []string{o.State}, pm.StateFailed, "Stopped after repeated errors: "+trimTo(err.Error(), 300)+". Operations have been alerted.", nil) && e.GC != nil {
			e.GC.LogDiscordFailedRequest(fmt.Sprintf("[publicmarkets] order %s failed: %v", o.ID, err))
		}
	}
}

func trimTo(s string, n int) string {
	if len(s) > n {
		return s[:n]
	}
	return s
}

func (e *Engine) step(ctx context.Context, o *pm.Order) error {
	a := &pm.Asset{}
	if err := e.DB.First(a, "id = ?", o.AssetID).Error; err != nil {
		return err
	}
	switch o.State {
	case pm.StateAwaitingPayment:
		return e.confirmPayment(ctx, o, a)
	case pm.StateQueued:
		return e.route(ctx, o, a)
	case pm.StateFilledFromInventory:
		return e.sendTokenMove(ctx, o, a)
	case pm.StateSettlementFinal:
		if o.TokenTxHash == "" {
			return e.sendTokenMove(ctx, o, a) // slow path: settled, now mint / burn
		}
		return e.finish(ctx, o, a)
	case pm.StateSubmitted:
		return e.confirmTokenMove(ctx, o, a)
	case pm.StateChainFinal:
		return e.afterChainFinal(ctx, o, a)
	}
	return nil
}

// confirmPayment waits for the payment leg: an app wallet's operation
// (stablecoin for a creation, tokens for a redemption) or, for an
// exchange redemption, the move of the tokens out of the customer's
// Trovotech-held wallet.
func (e *Engine) confirmPayment(ctx context.Context, o *pm.Order, a *pm.Asset) error {
	if o.Channel == pm.ChannelExchange {
		return e.exchangeRedemptionTransfer(ctx, o, a)
	}
	var op userModels.WalletOperation
	if err := e.DB.First(&op, "id = ?", o.PaymentOperationID).Error; err != nil {
		return fmt.Errorf("payment operation %s: %w", o.PaymentOperationID, err)
	}
	switch {
	case op.Success != nil && *op.Success:
		tx := ""
		if op.TxHash != nil {
			tx = *op.TxHash
		}
		now := e.now()
		what := "Payment"
		if o.Type == pm.OrderRedemption {
			what = "Tokens received"
		}
		e.transition(o, []string{pm.StateAwaitingPayment}, pm.StateQueued, what+" confirmed on-chain", map[string]interface{}{"payment_tx_hash": tx, "payment_confirmed_at": &now})
	case op.Success != nil || op.Status == userModels.WalletOperationFailed:
		reason := "the wallet's transaction did not go through"
		if op.Error != nil {
			reason = trimTo(*op.Error, 200)
		}
		e.transition(o, []string{pm.StateAwaitingPayment}, pm.StateFailed, "Payment failed: "+reason+". Nothing was charged.", map[string]interface{}{"note": "Payment failed"})
	case op.Status == userModels.WalletOperationPending && e.now().After(op.ExpiresAt):
		e.transition(o, []string{pm.StateAwaitingPayment}, pm.StateFailed, "The payment expired before it was submitted.", map[string]interface{}{"note": "Payment expired"})
	}
	return nil
}

// route decides how a queued order fills: from inventory (creation),
// netted (redemption), or in the next session's net batch.
func (e *Engine) route(ctx context.Context, o *pm.Order, a *pm.Asset) error {
	if a.Status == pm.AssetHalted || a.Status != pm.AssetLive {
		return e.reject(ctx, o, a, "Asset halted by reconciliation; the order was refunded.")
	}
	p := e.CurrentPrice(a)
	if o.Type == pm.OrderCreation {
		qty := d(o.NetAmount).Div(p.Value).Truncate(int32(a.TokenDecimals))
		if p.Fresh && p.Value.IsPositive() {
			inv, err := e.Inventory(ctx, a)
			if err != nil {
				return err
			}
			if inv.GreaterThanOrEqual(qty) && qty.IsPositive() {
				e.transition(o, []string{pm.StateQueued}, pm.StateFilledFromInventory, fmt.Sprintf("%s already holds the units; filled at %s", e.custodianName(a), p.Value),
					map[string]interface{}{"path": pm.PathFast, "quantity": qty.String(), "reference_price": p.Value.String(), "price_source": p.Source})
				return nil
			}
		}
		note := fmt.Sprintf("Rolled into the next %s session's net order to %s", a.Market, e.dealingMemberName(a))
		e.transition(o, []string{pm.StateQueued}, pm.StatePendingExecution, note, map[string]interface{}{"path": pm.PathSlow})
		return nil
	}
	// redemption
	if p.Fresh && p.Value.IsPositive() {
		gross := d(o.Quantity).Mul(p.Value).Round(2)
		fee := gross.Mul(d(o.FeePercent)).Div(decimal.NewFromInt(100)).Round(2)
		net := gross.Sub(fee)
		free, err := e.TreasuryFreeCash(ctx)
		if err != nil {
			return err
		}
		if o.Channel == pm.ChannelExchange || free.GreaterThanOrEqual(net) {
			e.transition(o, []string{pm.StateQueued}, pm.StateFilledFromInventory, "Netted against today's demand; no market trade needed",
				map[string]interface{}{"path": pm.PathNetted, "amount": gross.String(), "fee": fee.String(), "net_amount": net.String(), "reference_price": p.Value.String(), "price_source": p.Source})
			return nil
		}
	}
	e.transition(o, []string{pm.StateQueued}, pm.StatePendingExecution, fmt.Sprintf("Rolled into the next %s session's net order to %s", a.Market, e.dealingMemberName(a)),
		map[string]interface{}{"path": pm.PathSlow})
	return nil
}

func (e *Engine) dealingMemberName(a *pm.Asset) string {
	var dm pm.DealingMember
	if e.DB.First(&dm, a.DealingMemberID).Error == nil {
		return dm.DealingMemberName
	}
	return "the Dealing Member"
}

// reject refuses a queued order and gives back what was paid: stablecoin
// from the treasury (app creation), tokens from the issuing Safe (app
// redemption), or the exchange's balance.
func (e *Engine) reject(ctx context.Context, o *pm.Order, a *pm.Asset, reason string) error {
	switch {
	case o.Channel == pm.ChannelExchange && o.Type == pm.OrderCreation:
		if err := e.creditExchange(o.ServiceLinkID, d(o.Amount), "REFUND", "refund:"+o.ID, "Refund of rejected order "+o.ID, "system"); err != nil {
			return err
		}
	case o.Channel == pm.ChannelExchange:
		// tokens go back to the customer's wallet
		if err := e.returnTokens(ctx, o, a); err != nil {
			return err
		}
	case o.Type == pm.OrderCreation && o.PaymentTxHash != "":
		if err := e.refundStable(ctx, o, d(o.Amount)); err != nil {
			return err
		}
	case o.Type == pm.OrderRedemption && o.PaymentTxHash != "":
		if err := e.returnTokens(ctx, o, a); err != nil {
			return err
		}
	}
	e.transition(o, []string{o.State}, pm.StateRejected, reason, map[string]interface{}{"note": reason})
	return nil
}

func (e *Engine) refundStable(ctx context.Context, o *pm.Order, amount decimal.Decimal) error {
	if e.Chain == nil {
		return errors.New("no chain configured")
	}
	ft, err := e.fundingToken(ctx)
	if err != nil {
		return err
	}
	hash, err := e.Chain.Send(ctx, e.Chain.TreasurySafe(), []gnosissafe.Call{callOf(aa.ERC20Transfer(ft.Address, common.HexToAddress(o.WalletAddress), units(amount.String(), ft.Decimals)))})
	if err != nil {
		return err
	}
	e.DB.Model(&pm.Order{}).Where("id = ?", o.ID).Update("payout_tx_hash", hash)
	e.note(o, "Refunded "+amount.StringFixed(2)+" "+ft.Code+" (tx "+hash+")")
	return e.Chain.Wait(ctx, hash)
}

func (e *Engine) returnTokens(ctx context.Context, o *pm.Order, a *pm.Asset) error {
	if e.Chain == nil {
		return errors.New("no chain configured")
	}
	hash, err := e.Chain.Send(ctx, common.HexToAddress(a.IssuingSafeAddress), []gnosissafe.Call{callOf(aa.ERC20Transfer(common.HexToAddress(a.ContractAddress), common.HexToAddress(o.WalletAddress), units(o.Quantity, a.TokenDecimals)))})
	if err != nil {
		return err
	}
	e.DB.Model(&pm.Order{}).Where("id = ?", o.ID).Update("token_tx_hash", hash)
	e.note(o, "Returned "+o.Quantity+" "+a.AssetCode+" to the wallet (tx "+hash+")")
	return e.Chain.Wait(ctx, hash)
}

// sendTokenMove mints a creation's tokens to the buyer, or burns a
// redemption's tokens (held by the issuing Safe since the sale).
func (e *Engine) sendTokenMove(ctx context.Context, o *pm.Order, a *pm.Asset) error {
	if e.Chain == nil {
		return errors.New("no chain configured")
	}
	if !common.IsHexAddress(a.ContractAddress) || !common.IsHexAddress(a.IssuingSafeAddress) {
		return errors.New("the asset has no token contract or issuing Safe")
	}
	token, issuing := common.HexToAddress(a.ContractAddress), common.HexToAddress(a.IssuingSafeAddress)
	amount := units(o.Quantity, a.TokenDecimals)
	if amount.Sign() <= 0 {
		return errors.New("nothing to move")
	}
	var call aa.Call
	what := "Mint"
	if o.Type == pm.OrderCreation {
		call = aa.ERC20Mint(token, common.HexToAddress(o.WalletAddress), amount)
	} else {
		call, what = aa.ERC20Burn(token, amount), "Burn"
	}
	// claim the step before sending, so it can never be sent twice
	from := o.State
	if !e.transition(o, []string{from}, pm.StateSubmitted, what+" submitted via the Public Markets multisig", map[string]interface{}{"token_tx_hash": "pending"}) {
		return nil
	}
	hash, err := e.Chain.Send(ctx, issuing, []gnosissafe.Call{callOf(call)})
	if err != nil {
		// not broadcast: put it back so it is tried again
		e.DB.Model(&pm.Order{}).Where("id = ?", o.ID).Updates(map[string]interface{}{"state": from, "token_tx_hash": ""})
		return err
	}
	e.DB.Model(&pm.Order{}).Where("id = ?", o.ID).Update("token_tx_hash", hash)
	return nil
}

func (e *Engine) confirmTokenMove(ctx context.Context, o *pm.Order, a *pm.Asset) error {
	if o.TokenTxHash == "" || o.TokenTxHash == "pending" {
		if e.now().Sub(o.UpdatedAt) > 10*time.Minute {
			return errors.New("the token transaction was not recorded; check the issuing Safe before retrying")
		}
		return nil
	}
	if err := e.Chain.Wait(ctx, o.TokenTxHash); err != nil {
		return err
	}
	e.transition(o, []string{pm.StateSubmitted}, pm.StateChainFinal, "Confirmed on-chain (tx "+o.TokenTxHash+")", nil)
	return nil
}

// afterChainFinal settles what is left: a fast creation is final at once
// (the units were already settled at the Custodian); a redemption pays
// out.
func (e *Engine) afterChainFinal(ctx context.Context, o *pm.Order, a *pm.Asset) error {
	if o.Type == pm.OrderCreation {
		if o.SettledAt == nil {
			now := e.now()
			e.transition(o, []string{pm.StateChainFinal}, pm.StateSettlementFinal, e.custodianName(a)+" confirmed the units are held for you", map[string]interface{}{"settled_at": &now})
			return nil
		}
		return e.complete(o, "Ledger updated")
	}
	return e.payRedemption(ctx, o, a)
}

func (e *Engine) finish(ctx context.Context, o *pm.Order, a *pm.Asset) error {
	if o.Type == pm.OrderCreation {
		return e.complete(o, "Ledger updated")
	}
	return e.payRedemption(ctx, o, a)
}

// payRedemption pays a redemption's proceeds: to the app wallet from the
// treasury, or onto the exchange's balance.
func (e *Engine) payRedemption(ctx context.Context, o *pm.Order, a *pm.Asset) error {
	net := d(o.NetAmount)
	if o.Channel == pm.ChannelExchange {
		if err := e.creditExchange(o.ServiceLinkID, net, "REDEMPTION", "redemption:"+o.ID, "Proceeds of "+o.ID, "system"); err != nil {
			return err
		}
		credited := fmt.Sprintf("%s credited to the exchange's balance", net.StringFixed(2))
		if o.SettledAt == nil {
			now := e.now()
			e.transition(o, []string{pm.StateChainFinal}, pm.StateSettlementFinal, credited, map[string]interface{}{"settled_at": &now})
			return nil
		}
		return e.complete(o, credited)
	}
	if o.PayoutTxHash == "" {
		if e.Chain == nil {
			return errors.New("no chain configured")
		}
		free, err := e.TreasuryFreeCash(ctx)
		if err != nil {
			return err
		}
		// this order's own payout is counted in the obligations
		if free.Add(net).LessThan(net) {
			e.DB.Model(&pm.Order{}).Where("id = ?", o.ID).Update("note", "Waiting for the treasury to be funded")
			return nil
		}
		ft, err := e.fundingToken(ctx)
		if err != nil {
			return err
		}
		if !e.claimPayout(o) {
			return nil
		}
		hash, err := e.Chain.Send(ctx, e.Chain.TreasurySafe(), []gnosissafe.Call{callOf(aa.ERC20Transfer(ft.Address, common.HexToAddress(o.WalletAddress), units(net.String(), ft.Decimals)))})
		if err != nil {
			e.DB.Model(&pm.Order{}).Where("id = ?", o.ID).Update("payout_tx_hash", "")
			return err
		}
		e.DB.Model(&pm.Order{}).Where("id = ?", o.ID).Update("payout_tx_hash", hash)
		o.PayoutTxHash = hash
	}
	if o.PayoutTxHash == "pending" {
		if e.now().Sub(o.UpdatedAt) > 10*time.Minute {
			return errors.New("the payout transaction was not recorded; check the treasury before retrying")
		}
		return nil
	}
	if err := e.Chain.Wait(ctx, o.PayoutTxHash); err != nil {
		return err
	}
	paid := fmt.Sprintf("%s %s paid to %s", net.StringFixed(2), o.FundingAssetCode, walletLabel(o))
	if o.SettledAt == nil {
		now := e.now()
		e.transition(o, []string{pm.StateChainFinal}, pm.StateSettlementFinal, paid, map[string]interface{}{"settled_at": &now})
		return nil
	}
	return e.complete(o, paid)
}

func (e *Engine) claimPayout(o *pm.Order) bool {
	res := e.DB.Model(&pm.Order{}).Where("id = ? AND payout_tx_hash = ''", o.ID).Updates(map[string]interface{}{"payout_tx_hash": "pending", "updated_at": e.now()})
	return res.Error == nil && res.RowsAffected == 1
}

func walletLabel(o *pm.Order) string {
	if o.WalletAlias != "" {
		return o.WalletAlias
	}
	return o.WalletAddress
}

func (e *Engine) complete(o *pm.Order, note string) error {
	now := e.now()
	updates := map[string]interface{}{"completed_at": &now}
	if o.SettledAt == nil {
		updates["settled_at"] = &now
	}
	e.transition(o, []string{pm.StateChainFinal, pm.StateSettlementFinal}, pm.StateComplete, note, updates)
	return nil
}
