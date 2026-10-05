package publicmarkets

import (
	"context"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
)

// ledgerTotal is Σ the beneficial ownership ledger of an asset (tokens).
func (s *Service) ledgerTotal(a *Asset) (decimal.Decimal, int64) {
	var balances []string
	s.DB.Model(&LedgerEntry{}).Where("asset_id = ? AND balance <> '0'", a.ID).Pluck("balance", &balances)
	total := decimal.Zero
	for _, b := range balances {
		total = total.Add(d(b))
	}
	return total.Shift(-int32(a.TokenDecimals)), int64(len(balances))
}

func (s *Service) position(assetID string) (*CustodianPosition, bool) {
	var p CustodianPosition
	if s.DB.Where("asset_id = ?", assetID).Order("id DESC").First(&p).Error != nil {
		return nil, false
	}
	return &p, true
}

func (s *Service) latestRun(assetID string) (*ReconciliationRun, bool) {
	var r ReconciliationRun
	if s.DB.Where("asset_id = ?", assetID).Order("id DESC").First(&r).Error != nil {
		return nil, false
	}
	return &r, true
}

func (s *Service) custodianName(id uint64) string {
	var name string
	s.DB.Table("approved_asset_custodians").Where("id = ?", id).Limit(1).Pluck("asset_custodian_name", &name)
	return name
}

func (s *Service) dealingMemberName(id uint64) string {
	var dm DealingMember
	if s.DB.First(&dm, "id = ?", id).Error != nil {
		return ""
	}
	return dm.DealingMemberName
}

// AssetRow is an asset as the asset list shows it.
type AssetRow struct {
	Asset
	Supply            string     `json:"supply"` // Σ ledger
	Owners            int64      `json:"owners"`
	Position          string     `json:"position"`
	PositionAsOf      *time.Time `json:"positionAsOf"`
	PositionSource    string     `json:"positionSource"`
	CustodianName     string     `json:"custodianName"`
	DealingMemberName string     `json:"dealingMemberName"`
	MarketValue       string     `json:"marketValue"`
	LastReconResult   string     `json:"lastReconResult"`
	LastReconAt       *time.Time `json:"lastReconAt"`
	PriceFresh        bool       `json:"priceFresh"`
	MarketOpen        bool       `json:"marketOpen"`
}

func (s *Service) row(a Asset, st Settings) AssetRow {
	supply, owners := s.ledgerTotal(&a)
	r := AssetRow{Asset: a, Supply: supply.String(), Owners: owners, CustodianName: s.custodianName(a.CustodianID),
		DealingMemberName: s.dealingMemberName(a.DealingMemberID), MarketValue: supply.Mul(d(a.LastPrice)).StringFixed(2),
		MarketOpen: marketOpen(st, a.Market, s.now())}
	if p, ok := s.position(a.ID); ok {
		r.Position, r.PositionAsOf, r.PositionSource = p.RealUnitsHeld, &p.AsOf, p.Source
	}
	if run, ok := s.latestRun(a.ID); ok {
		r.LastReconResult, r.LastReconAt = run.Result, &run.CreatedAt
	}
	r.PriceFresh = a.PriceAt != nil && s.now().Sub(*a.PriceAt) <= time.Duration(st.PriceStaleMinutes)*time.Minute
	return r
}

// AssetFilters narrow the asset list.
type AssetFilters struct {
	Market, Type, Status, Search string
	CustodianID                  uint64
}

// Assets lists the Public Market assets.
func (s *Service) Assets(f AssetFilters) []AssetRow {
	q := s.DB.Model(&Asset{})
	if f.Market != "" {
		q = q.Where("market = ?", strings.ToUpper(f.Market))
	}
	if f.Type != "" {
		q = q.Where("asset_type = ?", strings.ToUpper(f.Type))
	}
	if f.Status != "" {
		q = q.Where("status = ?", strings.ToUpper(f.Status))
	}
	if f.CustodianID != 0 {
		q = q.Where("custodian_id = ?", f.CustodianID)
	}
	if v := strings.TrimSpace(f.Search); v != "" {
		like := "%" + strings.ToLower(v) + "%"
		q = q.Where("LOWER(asset_code) LIKE ? OR LOWER(ticker) LIKE ? OR LOWER(isin) LIKE ? OR LOWER(instrument_name) LIKE ?", like, like, like, like)
	}
	var assets []Asset
	q.Order("asset_code").Find(&assets)
	st := s.LoadSettings()
	out := make([]AssetRow, 0, len(assets))
	for _, a := range assets {
		out = append(out, s.row(a, st))
	}
	return out
}

func (s *Service) asset(id string) (*Asset, error) {
	var a Asset
	if err := s.DB.First(&a, "id = ? OR asset_code = ? OR ticker = ?", id, strings.ToUpper(id), strings.ToUpper(id)).Error; err != nil {
		return nil, notFound("asset", err)
	}
	return &a, nil
}

// Holder is a beneficial owner of an asset.
type Holder struct {
	WalletAddress   string `json:"walletAddress"`
	Channel         string `json:"channel"`
	ServiceLinkID   string `json:"serviceLinkId"`
	ExchangeName    string `json:"exchangeName"`
	Balance         string `json:"balance"`
	PercentOfSupply string `json:"percentOfSupply"`
	Substantial     bool   `json:"substantial"` // above the disclosure threshold (CAMA s.120)
}

// SetupStep is one step of bringing an asset live.
type SetupStep struct {
	Label  string `json:"label"`
	Detail string `json:"detail"`
	Done   bool   `json:"done"`
}

// AssetDetail is one asset's page.
type AssetDetail struct {
	AssetRow
	Custodian        *Custodian          `json:"custodian"`
	DealingMember    *DealingMember      `json:"dealingMember"`
	Holders          []Holder            `json:"holders"`
	Runs             []ReconciliationRun `json:"reconciliationRuns"`
	Prices           []PriceSnapshot     `json:"prices"`
	CorporateActions []CorporateAction   `json:"corporateActions"`
	Steps            []SetupStep         `json:"setupSteps"`
	OpenOrders       int64               `json:"openOrders"`
}

func (s *Service) holders(a *Asset, total decimal.Decimal, threshold decimal.Decimal, limit int) []Holder {
	var rows []LedgerEntry
	s.DB.Where("asset_id = ? AND balance <> '0'", a.ID).Find(&rows)
	names := s.exchangeNames()
	out := make([]Holder, 0, len(rows))
	for _, r := range rows {
		bal := d(r.Balance).Shift(-int32(a.TokenDecimals))
		if !bal.IsPositive() {
			continue
		}
		pct := decimal.Zero
		if total.IsPositive() {
			pct = bal.Div(total).Mul(decimal.NewFromInt(100))
		}
		out = append(out, Holder{WalletAddress: r.WalletAddress, Channel: r.Channel, ServiceLinkID: r.ServiceLinkID, ExchangeName: names[r.ServiceLinkID],
			Balance: bal.String(), PercentOfSupply: pct.StringFixed(2), Substantial: r.Channel != ChannelPlatform && threshold.IsPositive() && pct.GreaterThanOrEqual(threshold)})
	}
	sort.Slice(out, func(i, j int) bool { return d(out[i].Balance).GreaterThan(d(out[j].Balance)) })
	if limit > 0 && len(out) > limit {
		out = out[:limit]
	}
	return out
}

// AssetDetail returns an asset's page.
func (s *Service) AssetDetail(id string) (*AssetDetail, error) {
	a, err := s.asset(id)
	if err != nil {
		return nil, err
	}
	st := s.LoadSettings()
	out := &AssetDetail{AssetRow: s.row(*a, st)}
	var c Custodian
	if s.DB.First(&c, "custodian_id = ?", a.CustodianID).Error == nil {
		out.Custodian = &c
	}
	var dm DealingMember
	if s.DB.First(&dm, "id = ?", a.DealingMemberID).Error == nil {
		out.DealingMember = &dm
	}
	out.Holders = s.holders(a, d(out.Supply), d(st.SubstantialHoldingPercent), 50)
	s.DB.Where("asset_id = ?", a.ID).Order("id DESC").Limit(10).Find(&out.Runs)
	s.DB.Where("asset_id = ?", a.ID).Order("id DESC").Limit(30).Find(&out.Prices)
	s.DB.Where("asset_id = ?", a.ID).Order("created_at DESC").Limit(20).Find(&out.CorporateActions)
	s.DB.Model(&Order{}).Where("asset_id = ? AND state NOT IN ?", a.ID, finalStates).Count(&out.OpenOrders)
	_, hasPosition := s.position(a.ID)
	out.Steps = []SetupStep{
		{"Asset created", a.AssetCode + " · " + a.InstrumentName, true},
		{"Custodian assigned", out.CustodianName, out.Custodian != nil},
		{"Dealing Member assigned", out.DealingMemberName, out.DealingMember != nil},
		{"Token contract registered", a.ContractAddress, a.ContractAddress != ""},
		{"Price set", a.PriceSource, d(a.LastPrice).IsPositive()},
		{"Custodian position recorded", out.Position, hasPosition},
		{"Live", "Creation and redemption open", a.Status != AssetSetup},
	}
	return out, nil
}

var finalStates = []string{StateComplete, StateRejected, StateFailed, StateCancelled}

// AssetRequest creates or edits an asset.
type AssetRequest struct {
	AssetCode            string  `json:"assetCode"`
	Ticker               string  `json:"ticker"`
	Market               string  `json:"market"`
	AssetType            string  `json:"assetType"`
	ISIN                 string  `json:"isin"`
	InstrumentName       string  `json:"instrumentName"`
	ShortName            string  `json:"shortName"`
	Sector               string  `json:"sector"`
	Description          string  `json:"description"`
	UnitDescription      string  `json:"unitDescription"`
	CustodianID          uint64  `json:"custodianId"`
	DealingMemberID      uint64  `json:"dealingMemberId"`
	OmnibusReference     string  `json:"omnibusReference"`
	MinimumBuy           string  `json:"minimumBuy"`
	FeePercent           string  `json:"feePercent"`
	InventoryTargetUnits string  `json:"inventoryTargetUnits"`
	LogoBackground       string  `json:"logoBackground"`
	LogoForeground       string  `json:"logoForeground"`
	LogoInitials         string  `json:"logoInitials"`
	LogoURL              string  `json:"logoUrl"`
	MarketCap            string  `json:"marketCap"`
	PERatio              string  `json:"peRatio"`
	DividendYield        string  `json:"dividendYield"`
	Coupon               string  `json:"coupon"`
	MaturityDate         *string `json:"maturityDate"`
}

func nonNegative(field, v string) error {
	if strings.TrimSpace(v) == "" {
		return nil
	}
	x, err := decimal.NewFromString(strings.TrimSpace(v))
	if err != nil || x.IsNegative() {
		return invalid(field, "%s must be a non-negative number", field)
	}
	return nil
}

func (s *Service) checkPartners(custodianID, dmID uint64) error {
	var c Custodian
	if err := s.DB.First(&c, "custodian_id = ?", custodianID).Error; err != nil || !c.Active {
		return invalid("custodianId", "choose a Custodian configured for Public Markets (Settings › Custodians)")
	}
	var dm DealingMember
	if err := s.DB.First(&dm, "id = ?", dmID).Error; err != nil || !dm.Active {
		return invalid("dealingMemberId", "choose an active Dealing Member")
	}
	return nil
}

// apply copies the editable fields of a request onto an asset.
func (r AssetRequest) apply(a *Asset) error {
	set := func(dst *string, v string) {
		if v = strings.TrimSpace(v); v != "" {
			*dst = v
		}
	}
	set(&a.InstrumentName, r.InstrumentName)
	set(&a.ShortName, r.ShortName)
	set(&a.Sector, r.Sector)
	set(&a.Description, r.Description)
	set(&a.UnitDescription, r.UnitDescription)
	set(&a.OmnibusReference, r.OmnibusReference)
	set(&a.LogoBackground, r.LogoBackground)
	set(&a.LogoForeground, r.LogoForeground)
	set(&a.LogoInitials, r.LogoInitials)
	set(&a.LogoURL, r.LogoURL)
	set(&a.MarketCap, r.MarketCap)
	set(&a.PERatio, r.PERatio)
	set(&a.DividendYield, r.DividendYield)
	set(&a.Coupon, r.Coupon)
	for field, v := range map[string]string{"minimumBuy": r.MinimumBuy, "feePercent": r.FeePercent, "inventoryTargetUnits": r.InventoryTargetUnits} {
		if err := nonNegative(field, v); err != nil {
			return err
		}
	}
	set(&a.MinimumBuy, r.MinimumBuy)
	set(&a.InventoryTargetUnits, r.InventoryTargetUnits)
	if r.FeePercent != "" {
		if d(r.FeePercent).GreaterThan(decimal.NewFromInt(10)) {
			return invalid("feePercent", "feePercent above 10%% looks like a mistake")
		}
		a.FeePercent = strings.TrimSpace(r.FeePercent)
	}
	if r.ISIN != "" {
		isin := strings.ToUpper(strings.TrimSpace(r.ISIN))
		if !isinPattern.MatchString(isin) {
			return invalid("isin", "an ISIN is 12 characters: 2-letter country, 9 alphanumerics and a check digit")
		}
		a.ISIN = isin
	}
	if r.MaturityDate != nil {
		if strings.TrimSpace(*r.MaturityDate) == "" {
			a.MaturityDate = nil
		} else if t, err := time.ParseInLocation("2006-01-02", *r.MaturityDate, Lagos); err == nil {
			a.MaturityDate = &t
		} else {
			return invalid("maturityDate", "maturityDate must be YYYY-MM-DD")
		}
	}
	return nil
}

// CreateAsset adds an asset in SETUP: it opens for orders once its token
// contract is registered, it has a price and a Custodian position, and an
// admin takes it live.
func (s *Service) CreateAsset(r AssetRequest, by string) (*Asset, error) {
	code := strings.ToUpper(strings.TrimSpace(r.AssetCode))
	if !codePattern.MatchString(code) {
		return nil, invalid("assetCode", "assetCode is 1-16 capital letters, digits, dots or dashes (e.g. DANGCEM-T)")
	}
	ticker := strings.ToUpper(strings.TrimSpace(r.Ticker))
	if ticker == "" {
		ticker = strings.TrimSuffix(code, "-T")
	}
	market, typ := strings.ToUpper(strings.TrimSpace(r.Market)), strings.ToUpper(strings.TrimSpace(r.AssetType))
	if market != MarketNGX && market != MarketFMDQ {
		return nil, invalid("market", "market must be NGX or FMDQ")
	}
	if typ != TypeEquity && typ != TypeBond {
		return nil, invalid("assetType", "assetType must be EQUITY or BOND")
	}
	if strings.TrimSpace(r.InstrumentName) == "" {
		return nil, invalid("instrumentName", "instrumentName is required")
	}
	if strings.TrimSpace(r.ISIN) == "" {
		return nil, invalid("isin", "isin is required")
	}
	if err := s.checkPartners(r.CustodianID, r.DealingMemberID); err != nil {
		return nil, err
	}
	var n int64
	s.DB.Model(&Asset{}).Where("asset_code = ? OR isin = ?", code, strings.ToUpper(strings.TrimSpace(r.ISIN))).Count(&n)
	if n > 0 {
		return nil, refuse(http.StatusConflict, "an asset with this code or ISIN already exists")
	}
	now := s.now()
	a := Asset{ID: "pma-" + strings.ToLower(code), AssetCode: code, Ticker: ticker, Market: market, AssetType: typ, CustodianID: r.CustodianID,
		DealingMemberID: r.DealingMemberID, Country: "NG", Currency: "NGN", MinimumBuy: "1000", InventoryTargetUnits: "0", LastPrice: "0",
		PreviousClose: "0", DayHigh: "0", DayLow: "0", Status: AssetSetup, CreatedBy: by, CreatedAt: now, UpdatedAt: now}
	if typ == TypeEquity {
		a.UnitDescription = "1 share"
	} else {
		a.UnitDescription = "100 NGN face value"
	}
	if err := r.apply(&a); err != nil {
		return nil, err
	}
	if a.OmnibusReference == "" {
		a.OmnibusReference = "POOL-" + ticker + "-01"
	}
	if a.LogoInitials == "" {
		a.LogoInitials = trimTo(ticker, 2)
	}
	if err := s.DB.Create(&a).Error; err != nil {
		return nil, err
	}
	return &a, nil
}

// UpdateAsset edits an asset. Its Custodian and Dealing Member change only
// while it is in SETUP (live orders are in flight with them).
func (s *Service) UpdateAsset(id string, r AssetRequest, by string) (*Asset, error) {
	a, err := s.asset(id)
	if err != nil {
		return nil, err
	}
	if (r.CustodianID != 0 && r.CustodianID != a.CustodianID) || (r.DealingMemberID != 0 && r.DealingMemberID != a.DealingMemberID) {
		if a.Status != AssetSetup {
			return nil, refuse(http.StatusConflict, "the Custodian and Dealing Member of a live asset cannot change; halt it and contact engineering")
		}
		cust, dm := a.CustodianID, a.DealingMemberID
		if r.CustodianID != 0 {
			cust = r.CustodianID
		}
		if r.DealingMemberID != 0 {
			dm = r.DealingMemberID
		}
		if err := s.checkPartners(cust, dm); err != nil {
			return nil, err
		}
		a.CustodianID, a.DealingMemberID = cust, dm
	}
	if r.ISIN != "" && a.Status != AssetSetup && !strings.EqualFold(r.ISIN, a.ISIN) {
		return nil, invalid("isin", "the ISIN of a live asset cannot change")
	}
	if err := r.apply(a); err != nil {
		return nil, err
	}
	a.UpdatedAt = s.now()
	if err := s.DB.Save(a).Error; err != nil {
		return nil, err
	}
	return a, nil
}

// RegisterContract records the asset's token after checking it on-chain:
// owned by the asset's issuing Safe, with nothing minted yet.
func (s *Service) RegisterContract(ctx context.Context, id, contract, issuingSafe string) (*Asset, error) {
	a, err := s.asset(id)
	if err != nil {
		return nil, err
	}
	if a.Status != AssetSetup {
		return nil, refuse(http.StatusConflict, "the token of a live asset cannot change")
	}
	if !validAddress(contract) {
		return nil, invalid("contractAddress", "contractAddress must be a 0x address")
	}
	if !validAddress(issuingSafe) {
		return nil, invalid("issuingSafeAddress", "issuingSafeAddress must be a 0x address")
	}
	if s.Chain == nil {
		return nil, refuse(http.StatusServiceUnavailable, "the chain is not configured (BASE_RPC_URL), so the contract cannot be verified")
	}
	var n int64
	s.DB.Model(&Asset{}).Where("LOWER(contract_address) = ? AND id <> ?", strings.ToLower(contract), a.ID).Count(&n)
	if n > 0 {
		return nil, refuse(http.StatusConflict, "this contract is already registered to another asset")
	}
	info, err := s.Chain.Token(ctx, common.HexToAddress(contract))
	if err != nil {
		return nil, refuse(http.StatusBadRequest, "could not read the token contract: %v", err)
	}
	if info.Owner != common.HexToAddress(issuingSafe) {
		return nil, invalid("issuingSafeAddress", "the token's owner is %s, not the issuing Safe", info.Owner.Hex())
	}
	if info.TotalSupply != nil && info.TotalSupply.Sign() != 0 {
		return nil, invalid("contractAddress", "the token already has a supply; register a freshly deployed token")
	}
	if info.Decimals < 0 || info.Decimals > 18 {
		return nil, invalid("contractAddress", "the token's decimals (%d) are not supported", info.Decimals)
	}
	res := s.DB.Model(&Asset{}).Where("id = ? AND status = ?", a.ID, AssetSetup).Updates(map[string]interface{}{
		"contract_address": common.HexToAddress(contract).Hex(), "issuing_safe_address": common.HexToAddress(issuingSafe).Hex(),
		"token_decimals": info.Decimals, "updated_at": s.now()})
	if res.Error != nil {
		return nil, res.Error
	}
	return s.asset(a.ID)
}

// GoLive opens a SETUP asset for orders.
func (s *Service) GoLive(id, by string) (*Asset, error) {
	detail, err := s.AssetDetail(id)
	if err != nil {
		return nil, err
	}
	if detail.Status != AssetSetup {
		return nil, refuse(http.StatusConflict, "the asset is already %s", detail.Status)
	}
	for _, st := range detail.Steps[:len(detail.Steps)-1] {
		if !st.Done {
			return nil, refuse(http.StatusConflict, "not ready: %s", strings.ToLower(st.Label))
		}
	}
	res := s.DB.Model(&Asset{}).Where("id = ? AND status = ?", detail.ID, AssetSetup).Updates(map[string]interface{}{"status": AssetLive, "updated_at": s.now()})
	if res.RowsAffected == 0 {
		return nil, refuse(http.StatusConflict, "the asset changed; reload it")
	}
	return s.asset(detail.ID)
}

// Halt stops creation and redemption of an asset (holders' balances are
// not frozen).
func (s *Service) Halt(id, reason, by string) (*Asset, error) {
	if strings.TrimSpace(reason) == "" {
		return nil, invalid("reason", "a reason is required")
	}
	a, err := s.asset(id)
	if err != nil {
		return nil, err
	}
	now := s.now()
	res := s.DB.Model(&Asset{}).Where("id = ? AND status = ?", a.ID, AssetLive).Updates(map[string]interface{}{
		"status": AssetHalted, "halt_reason": trimTo("Halted by "+by+": "+strings.TrimSpace(reason), 500), "halted_at": &now, "halted_by": by, "updated_at": now})
	if res.RowsAffected == 0 {
		return nil, refuse(http.StatusConflict, "only a live asset can be halted (it is %s)", a.Status)
	}
	return s.asset(a.ID)
}

// Resume asks the engine to reconcile a halted asset and reopen it if the
// fresh run matches (it never reopens on drift).
func (s *Service) Resume(id, by string) (*JobRequest, error) {
	a, err := s.asset(id)
	if err != nil {
		return nil, err
	}
	if a.Status != AssetHalted {
		return nil, refuse(http.StatusConflict, "the asset is not halted")
	}
	return s.requestJob("RESUME", a.AssetCode, by)
}

// SetPrice records a manual reference price (no vendor, or a correction).
func (s *Service) SetPrice(id, price, by string) (*Asset, error) {
	a, err := s.asset(id)
	if err != nil {
		return nil, err
	}
	p, err := decimal.NewFromString(strings.TrimSpace(price))
	if err != nil || !p.IsPositive() {
		return nil, invalid("price", "price must be a positive number")
	}
	now := s.now()
	open := marketOpen(s.LoadSettings(), a.Market, now)
	if err := s.DB.Create(&PriceSnapshot{AssetID: a.ID, Price: p.String(), Source: PriceManual, MarketHours: open, AsOf: now, CapturedAt: now}).Error; err != nil {
		return nil, err
	}
	updates := map[string]interface{}{"last_price": p.String(), "price_source": PriceManual, "price_at": &now, "updated_at": now}
	if !d(a.PreviousClose).IsPositive() {
		updates["previous_close"] = p.String()
	}
	if !d(a.DayHigh).IsPositive() || p.GreaterThan(d(a.DayHigh)) {
		updates["day_high"] = p.String()
	}
	if !d(a.DayLow).IsPositive() || p.LessThan(d(a.DayLow)) {
		updates["day_low"] = p.String()
	}
	s.DB.Model(&Asset{}).Where("id = ?", a.ID).Updates(updates)
	return s.asset(a.ID)
}

// RecordPosition records the Custodian's position for an asset, from its
// statement (MANUAL Custodians, or a correction). Reconciliation uses the
// latest position.
func (s *Service) RecordPosition(id, units, asOf, reference, by string) (*CustodianPosition, error) {
	a, err := s.asset(id)
	if err != nil {
		return nil, err
	}
	u, err := decimal.NewFromString(strings.TrimSpace(units))
	if err != nil || u.IsNegative() || !u.Equal(u.Truncate(0)) {
		return nil, invalid("unitsHeld", "unitsHeld must be a whole, non-negative number of units")
	}
	at := s.now()
	if strings.TrimSpace(asOf) != "" {
		t, err := time.ParseInLocation("2006-01-02", strings.TrimSpace(asOf), Lagos)
		if err != nil {
			return nil, invalid("asOf", "asOf must be YYYY-MM-DD")
		}
		at = t.Add(18 * time.Hour) // the statement's end of day
	}
	if strings.TrimSpace(reference) == "" {
		return nil, invalid("reference", "the statement's reference is required")
	}
	p := CustodianPosition{AssetID: a.ID, RealUnitsHeld: u.String(), AsOf: at.UTC(), Source: PositionManual, Reference: trimTo(reference, 120),
		RecordedBy: by, CreatedAt: s.now()}
	if err := s.DB.Create(&p).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

// ---------------------------------------------------------------- prices

// PriceRow is an asset on the price oracle page.
type PriceRow struct {
	AssetID       string     `json:"assetId"`
	AssetCode     string     `json:"assetCode"`
	Market        string     `json:"market"`
	LastPrice     string     `json:"lastPrice"`
	PreviousClose string     `json:"previousClose"`
	Source        string     `json:"source"`
	CapturedAt    *time.Time `json:"capturedAt"`
	MarketOpen    bool       `json:"marketOpen"`
	Fresh         bool       `json:"fresh"`
	AgeMinutes    int64      `json:"ageMinutes"`
}

// Prices lists each asset's reference price and its freshness.
func (s *Service) Prices() []PriceRow {
	st := s.LoadSettings()
	var assets []Asset
	s.DB.Order("asset_code").Find(&assets)
	out := make([]PriceRow, 0, len(assets))
	now := s.now()
	for _, a := range assets {
		r := PriceRow{AssetID: a.ID, AssetCode: a.AssetCode, Market: a.Market, LastPrice: a.LastPrice, PreviousClose: a.PreviousClose,
			Source: a.PriceSource, CapturedAt: a.PriceAt, MarketOpen: marketOpen(st, a.Market, now), AgeMinutes: -1}
		if a.PriceAt != nil {
			r.AgeMinutes = int64(now.Sub(*a.PriceAt).Minutes())
			// after hours the last close is expected to be old
			r.Fresh = !r.MarketOpen || now.Sub(*a.PriceAt) <= time.Duration(st.PriceStaleMinutes)*time.Minute
		}
		out = append(out, r)
	}
	return out
}

// PriceHistory is an asset's captured prices, newest first.
func (s *Service) PriceHistory(id string, limit int) ([]PriceSnapshot, error) {
	a, err := s.asset(id)
	if err != nil {
		return nil, err
	}
	if limit <= 0 || limit > 1000 {
		limit = 100
	}
	var out []PriceSnapshot
	s.DB.Where("asset_id = ?", a.ID).Order("id DESC").Limit(limit).Find(&out)
	return out, nil
}

// Execution is a Dealing Member fill against the reference price.
type Execution struct {
	InstructionID  string     `json:"instructionId"`
	DealingMember  string     `json:"dealingMember"`
	AssetCode      string     `json:"assetCode"`
	Side           string     `json:"side"`
	Quantity       string     `json:"quantity"`
	ExecutedPrice  string     `json:"executedPrice"`
	ReferencePrice string     `json:"referencePrice"`
	DeviationPct   string     `json:"deviationPercent"`
	ExecutedAt     *time.Time `json:"executedAt"`
}

// Executions are recent fills with their deviation from the reference
// price taken when the batch was built.
func (s *Service) Executions(limit int) []Execution {
	if limit <= 0 || limit > 500 {
		limit = 50
	}
	var batches []NetBatch
	s.DB.Where("executed_price <> '' AND executed_price IS NOT NULL").Order("executed_at DESC").Limit(limit).Find(&batches)
	out := make([]Execution, 0, len(batches))
	for _, b := range batches {
		var in Instruction
		s.DB.Where("batch_id = ? AND kind = ?", b.ID, InstrDealingOrder).First(&in)
		dev := decimal.Zero
		if ref := d(b.ReferencePrice); ref.IsPositive() {
			dev = d(b.ExecutedPrice).Sub(ref).Div(ref).Mul(decimal.NewFromInt(100))
		}
		qty := b.ExecutedQuantity
		if qty == "" {
			qty = b.Quantity
		}
		out = append(out, Execution{InstructionID: in.ID, DealingMember: in.PartnerName, AssetCode: b.AssetCode, Side: b.Side, Quantity: qty,
			ExecutedPrice: b.ExecutedPrice, ReferencePrice: b.ReferencePrice, DeviationPct: dev.StringFixed(2), ExecutedAt: b.ExecutedAt})
	}
	return out
}

func (s *Service) exchangeNames() map[string]string {
	type row struct {
		ID        string
		ShortName string
		LongName  string
	}
	var rows []row
	s.DB.Table("service_links").Select("id, short_name, long_name").Where("id IN (?)", s.DB.Model(&ExchangePartner{}).Select("service_link_id")).Scan(&rows)
	out := map[string]string{}
	for _, r := range rows {
		name := r.LongName
		if name == "" {
			name = r.ShortName
		}
		out[r.ID] = name
	}
	return out
}
