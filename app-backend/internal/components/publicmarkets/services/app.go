package publicmarkets

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"trovo-wallet-api/internal/aa"
	pm "trovo-wallet-api/internal/components/publicmarkets/models"
	userModels "trovo-wallet-api/internal/components/users/models"
	usersvc "trovo-wallet-api/internal/components/users/services"

	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Wallet operation kinds of the app's Public Markets trades.
const (
	OperationBuy  = "PUBLIC MARKETS BUY"
	OperationSell = "PUBLIC MARKETS SELL"
)

// TradeRequest is an app buy or sell. Like every wallet operation it takes
// two calls: without transactionSignature it returns the operation to
// sign; with transaction and transactionSignature it submits it and places
// the order.
type TradeRequest struct {
	WalletAddress        string `json:"walletAddress"`
	Amount               string `json:"amount"`   // buy: funding currency to spend (fee included)
	Quantity             string `json:"quantity"` // sell: tokens
	Transaction          string `json:"transaction"`
	TransactionSignature string `json:"transactionSignature"`
}

// TradeResponse is the quote and the operation to sign, or the order.
type TradeResponse struct {
	Quote       Quote      `json:"quote"`
	Transaction string     `json:"transaction,omitempty"`
	Messages    []string   `json:"messages,omitempty"`
	Order       *OrderView `json:"order,omitempty"`
}

type tradeContext struct {
	AssetID string `json:"assetId"`
	Quote   Quote  `json:"quote"`
}

func (e *Engine) ownWallet(user *userModels.User, address string) (*userModels.UserWallet, error) {
	var w userModels.UserWallet
	if err := e.DB.Where("LOWER(id) = ? AND user_id = ?", strings.ToLower(strings.TrimSpace(address)), user.ID).First(&w).Error; err != nil {
		return nil, refuse(http.StatusBadRequest, "error-invalid-wallet", "walletAddress", "Choose one of your own wallets.")
	}
	if w.SharedAccessEnabled == 1 && w.NumberOfApprovalsNeeded > 0 {
		return nil, refuse(http.StatusBadRequest, "error-shared-wallet", "walletAddress", "Public Markets trades are not available on wallets that need approvers.")
	}
	return &w, nil
}

func (e *Engine) checkTrader(user *userModels.User) error {
	if err := user.EnsureNotSuspended(); err != nil {
		return err
	}
	if user.KYCVerified == 0 {
		return refuse(http.StatusForbidden, "error-invalid-kyc", "", "Complete your KYC to buy or sell Public Markets assets.")
	}
	return nil
}

// AppBuy is the app's creation: the wallet pays the treasury in one
// operation; the order fills once that payment is mined.
func (e *Engine) AppBuy(ctx context.Context, user *userModels.User, a *pm.Asset, r TradeRequest) (*TradeResponse, error) {
	return e.appTrade(ctx, user, a, r, OperationBuy)
}

// AppSell is the app's redemption: the wallet sends its tokens to the
// issuing Safe in one operation; proceeds are paid once the order settles.
func (e *Engine) AppSell(ctx context.Context, user *userModels.User, a *pm.Asset, r TradeRequest) (*TradeResponse, error) {
	return e.appTrade(ctx, user, a, r, OperationSell)
}

func (e *Engine) appTrade(ctx context.Context, user *userModels.User, a *pm.Asset, r TradeRequest, kind string) (*TradeResponse, error) {
	if err := e.checkTrader(user); err != nil {
		return nil, err
	}
	w, err := e.ownWallet(user, r.WalletAddress)
	if err != nil {
		return nil, err
	}
	if r.Transaction != "" && r.TransactionSignature != "" {
		return e.submitTrade(ctx, user, w, a, r, kind)
	}
	if e.Chain == nil || e.Chain.TreasurySafe() == (common.Address{}) || !common.IsHexAddress(a.ContractAddress) || !common.IsHexAddress(a.IssuingSafeAddress) {
		return nil, refuse(http.StatusServiceUnavailable, "error-not-configured", "", "Public Markets trading is not available right now.")
	}
	var q Quote
	var call aa.Call
	var message string
	if kind == OperationBuy {
		if q, err = e.QuoteCreation(ctx, a, d(r.Amount)); err != nil {
			return nil, err
		}
		ft, err := e.fundingToken(ctx)
		if err != nil {
			return nil, err
		}
		call = aa.ERC20Transfer(ft.Address, e.Chain.TreasurySafe(), units(q.Amount, ft.Decimals))
		message = fmt.Sprintf("You pay %s %s (including a %s %s Trovo fee) for about %s %s.", d(q.Amount).StringFixed(2), ft.Code, d(q.Fee).StringFixed(2), ft.Code, q.Quantity, a.AssetCode)
	} else {
		if q, err = e.QuoteRedemption(ctx, a, d(r.Quantity)); err != nil {
			return nil, err
		}
		bal, err := e.Chain.BalanceOf(ctx, common.HexToAddress(a.ContractAddress), common.HexToAddress(w.ID))
		if err != nil {
			return nil, refuse(http.StatusServiceUnavailable, "error-temporary", "", "Your balance could not be read. Please try again.")
		}
		held := decimal.NewFromBigInt(bal, -int32(a.TokenDecimals)).Sub(e.pendingRedemptions(w.ID, a.ID))
		if held.LessThan(d(q.Quantity)) {
			return nil, refuse(http.StatusBadRequest, "error-insufficient-holding", "quantity", "You only hold %s %s.", held, a.AssetCode)
		}
		call = aa.ERC20Transfer(common.HexToAddress(a.ContractAddress), common.HexToAddress(a.IssuingSafeAddress), units(q.Quantity, a.TokenDecimals))
		message = fmt.Sprintf("You sell %s %s and receive about %s %s after the %s%% Trovo fee.", q.Quantity, a.AssetCode, d(q.NetAmount).StringFixed(2), q.FundingAsset, q.FeePercent)
	}
	op, err := usersvc.PrepareWalletOperation(ctx, kind, user, user, w, []aa.Call{call}, 0, tradeContext{AssetID: a.ID, Quote: q}, e.GC)
	if err != nil {
		return nil, err
	}
	return &TradeResponse{Quote: q, Transaction: op.Transaction, Messages: append([]string{message, q.Note}, op.Messages()...)}, nil
}

func (e *Engine) submitTrade(ctx context.Context, user *userModels.User, w *userModels.UserWallet, a *pm.Asset, r TradeRequest, kind string) (*TradeResponse, error) {
	rec, p, err := usersvc.LoadWalletOperation(r.Transaction, w.ID, kind, e.GC)
	if err != nil {
		return nil, err
	}
	var c tradeContext
	if rec.Context == nil || json.Unmarshal([]byte(*rec.Context), &c) != nil || c.AssetID != a.ID {
		return nil, refuse(http.StatusBadRequest, "error-invalid-transaction", "transaction", "This transaction is not a %s trade.", a.AssetCode)
	}
	if _, err := usersvc.SignSingleOwnerOperation(ctx, rec, p, user.PrimarySigner, r.TransactionSignature, e.GC); err != nil {
		return nil, err
	}
	o, err := e.Place(OrderRequest{Asset: a, Channel: pm.ChannelApp, Username: user.Username, WalletAddress: w.ID, WalletAlias: w.Alias,
		Quote: c.Quote, State: pm.StateAwaitingPayment, PaymentOperationID: rec.ID})
	if err != nil {
		e.logf("operation %s submitted but its order was not saved: %v", rec.ID, err)
		if e.GC != nil {
			e.GC.LogDiscordFailedRequest(fmt.Sprintf("[publicmarkets] %s operation %s by %s submitted but the order was not saved: %v", kind, rec.ID, user.Username, err))
		}
		return nil, err
	}
	e.Kick()
	v := e.OrderViewOf(o)
	return &TradeResponse{Quote: c.Quote, Order: &v}, nil
}

// ---------------------------------------------------------------- views

// OrderView is an order with its timeline.
type OrderView struct {
	pm.Order
	Events []pm.OrderEvent `json:"events"`
}

// OrderViewOf loads an order's timeline.
func (e *Engine) OrderViewOf(o *pm.Order) OrderView {
	var evs []pm.OrderEvent
	e.DB.Where("order_id = ?", o.ID).Order("id").Find(&evs)
	return OrderView{Order: *o, Events: evs}
}

// AssetView is an asset as the apps and exchanges see it.
type AssetView struct {
	ID                  string                 `json:"id"`
	AssetCode           string                 `json:"assetCode"`
	Ticker              string                 `json:"ticker"`
	Name                string                 `json:"name"`
	ShortName           string                 `json:"shortName"`
	Market              string                 `json:"market"`
	AssetType           string                 `json:"assetType"`
	Sector              string                 `json:"sector"`
	ISIN                string                 `json:"isin"`
	Description         string                 `json:"description"`
	UnitDescription     string                 `json:"unitDescription"`
	Logo                map[string]string      `json:"logo"`
	Price               string                 `json:"price"`
	PriceSource         string                 `json:"priceSource"`
	PriceAt             *time.Time             `json:"priceAt"`
	PriceLive           bool                   `json:"priceLive"`
	DayChange           string                 `json:"dayChangePercent"`
	PreviousClose       string                 `json:"previousClose"`
	DayHigh             string                 `json:"dayHigh"`
	DayLow              string                 `json:"dayLow"`
	DayVolume           string                 `json:"dayVolume"`
	MarketCap           string                 `json:"marketCap"`
	PERatio             string                 `json:"peRatio"`
	DividendYield       string                 `json:"dividendYield"`
	Coupon              string                 `json:"coupon"`
	MaturityDate        *time.Time             `json:"maturityDate"`
	Status              string                 `json:"status"` // open | halted | coming-soon
	ContractAddress     string                 `json:"contractAddress"`
	TokenDecimals       int                    `json:"tokenDecimals"`
	MinimumBuy          string                 `json:"minimumBuy"`
	FeePercent          string                 `json:"feePercent"`
	FundingAsset        string                 `json:"fundingAsset"`
	Currency            string                 `json:"currency"`
	Session             map[string]interface{} `json:"session"`
	Custody             map[string]interface{} `json:"custody,omitempty"`
	CorporateActions    []pm.CorporateAction   `json:"corporateActions,omitempty"`
	TokensInCirculation string                 `json:"tokensInCirculation,omitempty"`
}

func viewStatus(s string) string {
	switch s {
	case pm.AssetLive:
		return "open"
	case pm.AssetHalted:
		return "halted"
	}
	return "coming-soon"
}

// ViewAsset renders an asset; detail adds custody, supply and corporate
// actions.
func (e *Engine) ViewAsset(ctx context.Context, a *pm.Asset, detail bool) AssetView {
	s := LoadSettings(e.DB)
	p := e.CurrentPrice(a)
	sess := map[string]interface{}{"open": p.Session.Open, "market": a.Market, "opensAt": p.Session.OpenAt, "closesAt": p.Session.Close}
	if p.Session.Open {
		sess["minutesRemaining"] = int(p.Session.Close.Sub(e.now()).Minutes())
		sess["minutesTotal"] = int(p.Session.Close.Sub(p.Session.OpenAt).Minutes())
	}
	v := AssetView{ID: a.ID, AssetCode: a.AssetCode, Ticker: a.Ticker, Name: a.InstrumentName, ShortName: a.ShortName, Market: a.Market, AssetType: a.AssetType,
		Sector: a.Sector, ISIN: a.ISIN, Description: a.Description, UnitDescription: a.UnitDescription,
		Logo:  map[string]string{"background": a.LogoBackground, "foreground": a.LogoForeground, "initials": a.LogoInitials, "url": a.LogoURL},
		Price: a.LastPrice, PriceSource: a.PriceSource, PriceAt: a.PriceAt, PriceLive: p.Fresh, DayChange: DayChange(a).String(), PreviousClose: a.PreviousClose,
		DayHigh: a.DayHigh, DayLow: a.DayLow, DayVolume: a.DayVolume, MarketCap: a.MarketCap, PERatio: a.PERatio, DividendYield: a.DividendYield, Coupon: a.Coupon,
		MaturityDate: a.MaturityDate, Status: viewStatus(a.Status), ContractAddress: a.ContractAddress, TokenDecimals: a.TokenDecimals, MinimumBuy: a.MinimumBuy,
		FeePercent: FeePercent(a, s).String(), FundingAsset: s.FundingAssetCode, Currency: a.Currency, Session: sess}
	if !detail {
		v.TokensInCirculation = e.LedgerTotal(a).String()
		return v
	}
	supply, err := e.Supply(ctx, a)
	if err == nil {
		v.TokensInCirculation = supply.String()
	}
	var c pm.Custodian
	e.DB.First(&c, "custodian_id = ?", a.CustodianID)
	custody := map[string]interface{}{"custodian": e.custodianName(a), "nominee": c.NomineeName, "dealingMember": e.dealingMemberName(a),
		"depository": "CSCS — pooled custody account", "settlement": "Instant from inventory, otherwise T+2"}
	if pos, ok := e.Position(a); ok {
		custody["unitsHeld"] = pos.RealUnitsHeld
		custody["positionAsOf"] = pos.AsOf
	}
	var run pm.ReconciliationRun
	if e.DB.Where("asset_id = ?", a.ID).Order("id DESC").First(&run).Error == nil {
		custody["lastReconciliation"] = map[string]interface{}{"at": run.CreatedAt, "result": run.Result}
	}
	v.Custody = custody
	e.DB.Where("asset_id = ? AND status <> ?", a.ID, pm.ActionCancelledCA).Order("record_date DESC").Limit(5).Find(&v.CorporateActions)
	return v
}

// ListAssets lists the assets open (or halted / coming) to the public,
// filtered by market, asset type and a search.
func (e *Engine) ListAssets(ctx context.Context, market, assetType, search string) []AssetView {
	q := e.DB.Where("status IN ?", []string{pm.AssetLive, pm.AssetHalted, pm.AssetSetup})
	if m := strings.ToUpper(strings.TrimSpace(market)); m != "" {
		q = q.Where("market = ?", m)
	}
	if t := strings.ToUpper(strings.TrimSpace(assetType)); t != "" {
		q = q.Where("asset_type = ?", t)
	}
	if s := strings.ToLower(strings.TrimSpace(search)); s != "" {
		like := "%" + s + "%"
		q = q.Where("LOWER(asset_code) LIKE ? OR LOWER(instrument_name) LIKE ? OR LOWER(ticker) LIKE ?", like, like, like)
	}
	var assets []pm.Asset
	q.Order("asset_code").Find(&assets)
	out := make([]AssetView, 0, len(assets))
	for i := range assets {
		out = append(out, e.ViewAsset(ctx, &assets[i], false))
	}
	return out
}

// ---------------------------------------------------------------- portfolio

// Holding is one asset a user holds across their wallets.
type Holding struct {
	Asset          AssetView         `json:"asset"`
	Quantity       string            `json:"quantity"`
	MarketValue    string            `json:"marketValue"`
	AverageCost    string            `json:"averageCost"`
	CostBasis      string            `json:"costBasis"`
	TotalReturn    string            `json:"totalReturn"`
	ReturnPercent  string            `json:"returnPercent"`
	TodayChange    string            `json:"todayChange"`
	IncomeReceived string            `json:"incomeReceived"`
	Wallets        map[string]string `json:"wallets"` // wallet address -> quantity
}

// Portfolio is a user's Public Markets position (My Stocks).
type Portfolio struct {
	Value         string     `json:"value"`
	CostBasis     string     `json:"costBasis"`
	TotalReturn   string     `json:"totalReturn"`
	ReturnPercent string     `json:"returnPercent"`
	TodayChange   string     `json:"todayChange"`
	Income        string     `json:"income"`
	Holdings      []Holding  `json:"holdings"`
	OpenOrders    []pm.Order `json:"openOrders"`
	Activity      []Activity `json:"activity"`
}

// Activity is one line of recent activity (trades and dividends).
type Activity struct {
	Kind      string    `json:"kind"` // BUY | SELL | DIVIDEND | COUPON
	Title     string    `json:"title"`
	Detail    string    `json:"detail"`
	Received  string    `json:"received"`
	Spent     string    `json:"spent"`
	Reference string    `json:"reference"`
	At        time.Time `json:"at"`
}

func (e *Engine) userWallets(user *userModels.User) []string {
	var ws []userModels.UserWallet
	e.DB.Select("id").Where("user_id = ?", user.ID).Find(&ws)
	out := make([]string, 0, len(ws))
	for _, w := range ws {
		out = append(out, strings.ToLower(w.ID))
	}
	return out
}

// PortfolioOf is a user's holdings (from the ownership ledger), cost and
// returns (from their completed orders), income, open orders and activity.
func (e *Engine) PortfolioOf(ctx context.Context, user *userModels.User) Portfolio {
	wallets := e.userWallets(user)
	out := Portfolio{Holdings: []Holding{}, OpenOrders: []pm.Order{}, Activity: []Activity{}}
	if len(wallets) == 0 {
		return out
	}
	var rows []pm.LedgerEntry
	e.DB.Where("LOWER(wallet_address) IN ? AND balance <> '0'", wallets).Find(&rows)
	byAsset := map[string][]pm.LedgerEntry{}
	for _, r := range rows {
		byAsset[r.AssetID] = append(byAsset[r.AssetID], r)
	}
	var orders []pm.Order
	e.DB.Where("username = ?", user.Username).Order("created_at DESC").Find(&orders)
	var ents []pm.Entitlement
	e.DB.Where("LOWER(wallet_address) IN ? AND status = ?", wallets, pm.EntitlementPaid).Order("paid_at DESC").Find(&ents)
	value, cost, today, income := decimal.Zero, decimal.Zero, decimal.Zero, decimal.Zero
	for assetID, entries := range byAsset {
		var a pm.Asset
		if e.DB.First(&a, "id = ?", assetID).Error != nil {
			continue
		}
		qty := decimal.Zero
		per := map[string]string{}
		for _, r := range entries {
			q := d(r.Balance).Shift(-int32(a.TokenDecimals))
			qty = qty.Add(q)
			per[r.WalletAddress] = q.String()
		}
		if !qty.IsPositive() {
			continue
		}
		// average cost: what completed buys cost per token (sales take
		// tokens out at that average)
		boughtQty, boughtCost := decimal.Zero, decimal.Zero
		for i := len(orders) - 1; i >= 0; i-- {
			o := orders[i]
			if o.AssetID != assetID || o.State != pm.StateComplete {
				continue
			}
			if o.Type == pm.OrderCreation {
				boughtQty, boughtCost = boughtQty.Add(d(o.Quantity)), boughtCost.Add(d(o.Amount))
			} else if boughtQty.IsPositive() {
				avg := boughtCost.Div(boughtQty)
				sold := decimal.Min(d(o.Quantity), boughtQty)
				boughtQty, boughtCost = boughtQty.Sub(sold), boughtCost.Sub(sold.Mul(avg))
			}
		}
		price := d(a.LastPrice)
		avg := price
		if boughtQty.IsPositive() {
			avg = boughtCost.Div(boughtQty)
		}
		mv := qty.Mul(price)
		basis := qty.Mul(avg)
		ret := mv.Sub(basis)
		pct := decimal.Zero
		if basis.IsPositive() {
			pct = ret.Div(basis).Mul(decimal.NewFromInt(100))
		}
		dayChg := mv.Mul(DayChange(&a)).Div(decimal.NewFromInt(100))
		inc := decimal.Zero
		for _, en := range ents {
			if en.AssetID == assetID {
				inc = inc.Add(d(en.NetAmount))
			}
		}
		out.Holdings = append(out.Holdings, Holding{Asset: e.ViewAsset(ctx, &a, false), Quantity: qty.String(), MarketValue: mv.StringFixed(2),
			AverageCost: avg.StringFixed(2), CostBasis: basis.StringFixed(2), TotalReturn: ret.StringFixed(2), ReturnPercent: pct.StringFixed(2),
			TodayChange: dayChg.StringFixed(2), IncomeReceived: inc.StringFixed(2), Wallets: per})
		value, cost, today, income = value.Add(mv), cost.Add(basis), today.Add(dayChg), income.Add(inc)
	}
	sort.Slice(out.Holdings, func(i, j int) bool { return d(out.Holdings[i].MarketValue).GreaterThan(d(out.Holdings[j].MarketValue)) })
	pct := decimal.Zero
	if cost.IsPositive() {
		pct = value.Sub(cost).Div(cost).Mul(decimal.NewFromInt(100))
	}
	out.Value, out.CostBasis, out.TotalReturn, out.ReturnPercent = value.StringFixed(2), cost.StringFixed(2), value.Sub(cost).StringFixed(2), pct.StringFixed(2)
	out.TodayChange, out.Income = today.StringFixed(2), income.StringFixed(2)
	for _, o := range orders {
		switch o.State {
		case pm.StateComplete, pm.StateRejected, pm.StateFailed, pm.StateCancelled:
		default:
			out.OpenOrders = append(out.OpenOrders, o)
		}
		if o.State == pm.StateComplete && len(out.Activity) < 50 {
			if o.Type == pm.OrderCreation {
				out.Activity = append(out.Activity, Activity{Kind: "BUY", Title: "Bought " + o.AssetCode, Detail: o.ID, Received: "+" + o.Quantity + " " + o.AssetCode,
					Spent: "-" + d(o.Amount).StringFixed(2) + " " + o.FundingAssetCode, Reference: o.ID, At: o.CreatedAt})
			} else {
				out.Activity = append(out.Activity, Activity{Kind: "SELL", Title: "Sold " + o.AssetCode, Detail: o.ID, Received: "+" + d(o.NetAmount).StringFixed(2) + " " + o.FundingAssetCode,
					Spent: "-" + o.Quantity + " " + o.AssetCode, Reference: o.ID, At: o.CreatedAt})
			}
		}
	}
	for _, en := range ents {
		var ca pm.CorporateAction
		e.DB.First(&ca, "id = ?", en.CorporateActionID)
		at := en.UpdatedAt
		if en.PaidAt != nil {
			at = *en.PaidAt
		}
		out.Activity = append(out.Activity, Activity{Kind: ca.EventType, Title: cases.Title(language.English).String(strings.ToLower(ca.EventType)) + " · " + strings.TrimSuffix(en.AssetCode, "-T"),
			Detail: fmt.Sprintf("after %s%% WHT", en.WHTPercent), Received: "+" + d(en.NetAmount).StringFixed(2) + " " + LoadSettings(e.DB).FundingAssetCode,
			Spent: ca.Description, Reference: en.ID, At: at})
	}
	sort.Slice(out.Activity, func(i, j int) bool { return out.Activity[i].At.After(out.Activity[j].At) })
	return out
}

// DividendView is a holder's dividend with its breakdown.
type DividendView struct {
	pm.Entitlement
	EventType     string `json:"eventType"`
	Description   string `json:"description"`
	RecordDate    string `json:"recordDate"`
	PayDate       string `json:"payDate"`
	AmountPerUnit string `json:"amountPerUnit"`
	CustodianName string `json:"custodianName"`
}

// DividendsOf lists a user's dividends and coupons (paid and coming).
func (e *Engine) DividendsOf(user *userModels.User, assetCode string) []DividendView {
	wallets := e.userWallets(user)
	out := []DividendView{}
	if len(wallets) == 0 {
		return out
	}
	q := e.DB.Where("LOWER(wallet_address) IN ? AND status <> ?", wallets, pm.EntitlementRetained)
	if assetCode != "" {
		if a, err := e.AssetByCode(assetCode); err == nil {
			q = q.Where("asset_id = ?", a.ID)
		}
	}
	var ents []pm.Entitlement
	q.Order("created_at DESC").Limit(200).Find(&ents)
	for _, en := range ents {
		var ca pm.CorporateAction
		e.DB.First(&ca, "id = ?", en.CorporateActionID)
		var a pm.Asset
		e.DB.First(&a, "id = ?", en.AssetID)
		out = append(out, DividendView{Entitlement: en, EventType: ca.EventType, Description: ca.Description, RecordDate: ca.RecordDate, PayDate: ca.PayDate,
			AmountPerUnit: ca.AmountPerUnit, CustodianName: e.custodianName(&a)})
	}
	return out
}
