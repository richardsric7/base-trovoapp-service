package publicmarkets

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	pm "trovo-wallet-api/internal/components/publicmarkets/models"
	"trovo-wallet-api/internal/components/publicmarkets/partners"
	userModels "trovo-wallet-api/internal/components/users/models"
	appdb "trovo-wallet-api/internal/db"
	"trovo-wallet-api/internal/gnosissafe"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ---------------------------------------------------------------- a fake chain

// fakeChain applies the engine's Safe calls (ERC-20 transfer, mint, burn)
// to in-memory balances and keeps their Transfer logs.
type fakeChain struct {
	mu        sync.Mutex
	block     uint64
	balances  map[common.Address]map[common.Address]*big.Int
	supply    map[common.Address]*big.Int
	decimals  map[common.Address]uint8
	transfers map[common.Address][]Transfer
	sent      []string
	treasury  common.Address
	failSend  bool
}

func newFakeChain(treasury common.Address) *fakeChain {
	return &fakeChain{block: 100, balances: map[common.Address]map[common.Address]*big.Int{}, supply: map[common.Address]*big.Int{},
		decimals: map[common.Address]uint8{}, transfers: map[common.Address][]Transfer{}, treasury: treasury}
}

func (f *fakeChain) bal(token, holder common.Address) *big.Int {
	if f.balances[token] == nil {
		f.balances[token] = map[common.Address]*big.Int{}
	}
	if f.balances[token][holder] == nil {
		f.balances[token][holder] = new(big.Int)
	}
	return f.balances[token][holder]
}

func (f *fakeChain) move(token, from, to common.Address, v *big.Int, hash string, idx uint) error {
	if from != (common.Address{}) {
		if f.bal(token, from).Cmp(v) < 0 {
			return fmt.Errorf("insufficient balance of %s", from.Hex())
		}
		f.bal(token, from).Sub(f.bal(token, from), v)
	} else {
		if f.supply[token] == nil {
			f.supply[token] = new(big.Int)
		}
		f.supply[token].Add(f.supply[token], v)
	}
	if to != (common.Address{}) {
		f.bal(token, to).Add(f.bal(token, to), v)
	} else {
		f.supply[token].Sub(f.supply[token], v)
	}
	f.transfers[token] = append(f.transfers[token], Transfer{From: from, To: to, Value: new(big.Int).Set(v), Block: f.block, TxHash: hash, LogIndex: idx})
	return nil
}

func (f *fakeChain) Decimals(ctx context.Context, token common.Address) (uint8, error) {
	return f.decimals[token], nil
}
func (f *fakeChain) TotalSupply(ctx context.Context, token common.Address) (*big.Int, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.supply[token] == nil {
		return new(big.Int), nil
	}
	return new(big.Int).Set(f.supply[token]), nil
}
func (f *fakeChain) BalanceOf(ctx context.Context, token, holder common.Address) (*big.Int, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return new(big.Int).Set(f.bal(token, holder)), nil
}
func (f *fakeChain) Owner(ctx context.Context, c common.Address) (common.Address, error) {
	return common.Address{}, nil
}
func (f *fakeChain) Send(ctx context.Context, safe common.Address, calls []gnosissafe.Call) (string, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.failSend {
		return "", fmt.Errorf("rpc down")
	}
	f.block++
	hash := fmt.Sprintf("0x%064x", f.block)
	for i, c := range calls {
		sel := fmt.Sprintf("%x", c.Data[:4])
		args := c.Data[4:]
		word := func(n int) []byte { return args[32*n : 32*(n+1)] }
		var err error
		switch sel {
		case "a9059cbb": // transfer(to, amount)
			err = f.move(c.To, safe, common.BytesToAddress(word(0)), new(big.Int).SetBytes(word(1)), hash, uint(i))
		case "40c10f19": // mint(to, amount)
			err = f.move(c.To, common.Address{}, common.BytesToAddress(word(0)), new(big.Int).SetBytes(word(1)), hash, uint(i))
		case "42966c68": // burn(amount)
			err = f.move(c.To, safe, common.Address{}, new(big.Int).SetBytes(word(0)), hash, uint(i))
		default:
			err = fmt.Errorf("unexpected call %s", sel)
		}
		if err != nil {
			return "", err
		}
	}
	f.sent = append(f.sent, hash)
	return hash, nil
}
func (f *fakeChain) Wait(ctx context.Context, hash string) error { return nil }
func (f *fakeChain) PartnerSafe(ctx context.Context, salt *big.Int, deploy bool) (common.Address, error) {
	return common.BytesToAddress(salt.Bytes()[:20]), nil
}
func (f *fakeChain) Head(ctx context.Context) (uint64, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.block + ledgerConfirmations, nil
}
func (f *fakeChain) Transfers(ctx context.Context, token common.Address, from, to uint64) ([]Transfer, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	var out []Transfer
	for _, t := range f.transfers[token] {
		if t.Block >= from && t.Block <= to {
			out = append(out, t)
		}
	}
	return out, nil
}
func (f *fakeChain) BlockAtTime(ctx context.Context, t time.Time) (uint64, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.block, nil
}
func (f *fakeChain) TxTransfers(ctx context.Context, hash string, token common.Address) ([]Transfer, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	var out []Transfer
	for _, t := range f.transfers[token] {
		if t.TxHash == hash {
			out = append(out, t)
		}
	}
	return out, nil
}
func (f *fakeChain) TreasurySafe() common.Address { return f.treasury }

// ---------------------------------------------------------------- fixtures

var (
	treasury = common.HexToAddress("0x00000000000000000000000000000000000000aa")
	issuing  = common.HexToAddress("0x00000000000000000000000000000000000000bb")
	token    = common.HexToAddress("0x00000000000000000000000000000000000000cc")
	cngn     = common.HexToAddress("0x00000000000000000000000000000000000000dd")
	funder   = common.HexToAddress("0x00000000000000000000000000000000000000ee")
	feeW     = common.HexToAddress("0x00000000000000000000000000000000000000f1")
	whtW     = common.HexToAddress("0x00000000000000000000000000000000000000f2")
)

type fixture struct {
	t     *testing.T
	e     *Engine
	chain *fakeChain
	now   time.Time
	asset *pm.Asset
	link  string
	notes []string
}

// a Wednesday, 11:00 WAT: NGX is open
var sessionTime = time.Date(2026, 9, 30, 11, 0, 0, 0, Lagos)

func newFixture(t *testing.T) *fixture {
	t.Helper()
	db := testDB(t)
	appdb.MigrateDB(db)
	f := &fixture{t: t, chain: newFakeChain(treasury), now: sessionTime.UTC(), link: "sl-palmline"}
	f.chain.decimals[token], f.chain.decimals[cngn] = 4, 6
	f.e = &Engine{DB: db, Chain: f.chain, GC: &sharedconfig.GlobalConfig{DB: db}, decimals: map[common.Address]uint8{}}
	f.e.Now = func() time.Time { return f.now }
	f.e.Partners = DefaultPartners{E: f.e}
	f.e.Prices = MockPriceFeed{E: f.e}
	f.e.Notify = func(username, title, body string, data map[string]string) {
		f.notes = append(f.notes, username+": "+title)
	}
	f.e.GC = nil
	s := LoadSettings(db)
	db.Model(&pm.Settings{}).Where("id = 1").Updates(map[string]interface{}{"fee_wallet": feeW.Hex(), "wht_wallet": whtW.Hex()})
	_ = s
	db.Create(&userModels.TokenizationCurrency{AssetCode: "CNGN", ContractAddress: cngn.Hex()})
	db.Create(&userModels.ApprovedAssetCustodian{ID: 1, AssetCustodianName: "Custodian A Ltd", AssetCustodianCountry: "NGA"})
	db.Create(&pm.Custodian{CustodianID: 1, Code: "CUSTA", NomineeName: "Custodian A Nominees Ltd", Mode: pm.ModeMock})
	db.Create(&pm.DealingMember{ID: 1, DealingMemberName: "Broker Partner 2 Ltd", DealingMemberCountry: "NGA", Code: "BP2", Mode: pm.ModeMock})
	priceAt := f.now.Add(-time.Minute)
	a := pm.Asset{ID: "asset-mtnn", AssetCode: "MTNN-T", Ticker: "MTNN", Market: pm.MarketNGX, AssetType: pm.TypeEquity, ISIN: "NGMTNN000002",
		InstrumentName: "MTN Nigeria Communications Plc", ShortName: "MTN Nigeria", Country: "NG", ContractAddress: token.Hex(), IssuingSafeAddress: issuing.Hex(),
		TokenDecimals: 4, CustodianID: 1, DealingMemberID: 1, OmnibusReference: "POOL-MTNN-01", MinimumBuy: "1000", Status: pm.AssetLive,
		LastPrice: "221", PreviousClose: "216.45", PriceSource: pm.PriceMockFeed, PriceAt: &priceAt, CreatedAt: f.now.Add(-24 * time.Hour)}
	db.Create(&a)
	f.asset = &a
	db.Create(&pm.CustodianPosition{AssetID: a.ID, RealUnitsHeld: "1000", AsOf: f.now.Add(-time.Hour), Source: pm.PositionFeed})
	// an exchange with a prefunded balance (credited from real treasury cash)
	db.Create(&pm.ExchangePartner{ServiceLinkID: f.link, Status: "active", SigningSecret: "s3cret", FundingAddress: funder.Hex(), Balance: "0"})
	return f
}

func (f *fixture) fund(holder common.Address, tok common.Address, human string) {
	dec := f.chain.decimals[tok]
	v := decimal.RequireFromString(human).Shift(int32(dec)).BigInt()
	f.chain.mu.Lock()
	f.chain.bal(tok, holder).Add(f.chain.bal(tok, holder), v)
	f.chain.mu.Unlock()
}

func (f *fixture) deposit(amount string) {
	f.fund(treasury, cngn, amount)
	if err := f.e.creditExchange(f.link, decimal.RequireFromString(amount), "DEPOSIT", "dep-"+amount, "test", "test"); err != nil {
		f.t.Fatal(err)
	}
}

func (f *fixture) wallet(ref string) *pm.PartnerWallet {
	tax := "TIN-" + ref
	w, _, err := f.e.ProvisionWallet(context.Background(), f.link, ProvisionRequest{ExternalUserRef: ref, LegalName: "Adaeze Okafor", TaxIdentifier: tax,
		ResidencyCountry: "NGA", Nationality: "NGA", NDPAConsent: boolp(true)})
	if err != nil {
		f.t.Fatal(err)
	}
	var pw pm.PartnerWallet
	f.e.DB.First(&pw, "id = ?", w.WalletID)
	return &pw
}

func boolp(b bool) *bool { return &b }

// run steps the engine like its loops do, n times.
func (f *fixture) run(n int) {
	ctx := context.Background()
	for i := 0; i < n; i++ {
		f.e.DeliverMockEvents(ctx)
		f.e.ProcessPendingEvents()
		f.e.ProcessOrders(ctx)
		f.e.ReleaseApprovedBatches()
		f.e.BuildBatches(ctx)
		f.e.DispatchInstructions()
		f.e.IndexLedgers(ctx)
	}
}

func (f *fixture) order(id string) pm.Order {
	var o pm.Order
	f.e.DB.First(&o, "id = ?", id)
	return o
}

func (f *fixture) tokens(holder common.Address) string {
	b, _ := f.chain.BalanceOf(context.Background(), token, holder)
	return decimal.NewFromBigInt(b, -4).String()
}

// ---------------------------------------------------------------- tests

func TestQuoteFastAndSlow(t *testing.T) {
	f := newFixture(t)
	ctx := context.Background()
	q, err := f.e.QuoteCreation(ctx, f.asset, decimal.NewFromInt(221000))
	if err != nil {
		t.Fatal(err)
	}
	// 0.25% fee: 552.50; 220447.50 / 221 = 997.5
	if q.Fee != "552.5" || q.Quantity != "997.5" || q.Path != pm.PathFast {
		t.Fatalf("quote %+v", q)
	}
	if _, err := f.e.QuoteCreation(ctx, f.asset, decimal.NewFromInt(500)); err == nil {
		t.Fatal("below the minimum buy accepted")
	}
	// more than the 1000 units of inventory: slow
	q, _ = f.e.QuoteCreation(ctx, f.asset, decimal.NewFromInt(300000))
	if q.Path != pm.PathSlow {
		t.Fatalf("expected slow, got %s", q.Path)
	}
	// after hours: slow even with inventory, priced at the last close
	f.now = time.Date(2026, 9, 30, 20, 0, 0, 0, Lagos).UTC()
	q, _ = f.e.QuoteCreation(ctx, f.asset, decimal.NewFromInt(2210))
	if q.Path != pm.PathSlow || q.MarketOpen || q.NextSessionAt == nil {
		t.Fatalf("after hours %+v", q)
	}
}

func TestExchangeFastCreation(t *testing.T) {
	f := newFixture(t)
	f.deposit("100000")
	w := f.wallet("plm_usr_1")
	o, err := f.e.PlaceExchangeCreation(context.Background(), f.link, ExchangeOrderRequest{WalletID: w.ID, AssetCode: "MTNN", Amount: "50000", AmountCurrency: "NGN", ExternalOrderRef: "plm-1"})
	if err != nil {
		t.Fatal(err)
	}
	var p pm.ExchangePartner
	f.e.DB.First(&p, "service_link_id = ?", f.link)
	if p.Balance != "50000" {
		t.Fatalf("balance %s after the order", p.Balance)
	}
	f.run(6)
	got := f.order(o.ID)
	if got.State != pm.StateComplete || got.Path != pm.PathFast {
		t.Fatalf("order %s %s: %s", got.State, got.Path, got.LastError)
	}
	// 50000 - 125 fee = 49875 / 221 = 225.6787 (4 decimals)
	if got.Quantity != "225.6787" || f.tokens(common.HexToAddress(w.WalletAddress)) != "225.6787" {
		t.Fatalf("minted %s, wallet holds %s", got.Quantity, f.tokens(common.HexToAddress(w.WalletAddress)))
	}
	var hooks []pm.WebhookDelivery
	f.e.DB.Where("service_link_id = ? AND event = ?", f.link, "creation.settled").Find(&hooks)
	if len(hooks) != 1 || !strings.Contains(hooks[0].Payload, o.ID) {
		t.Fatalf("webhooks %+v", hooks)
	}
	// the ledger follows the mint, attributed to the exchange
	var row pm.LedgerEntry
	f.e.DB.First(&row, "asset_id = ? AND LOWER(wallet_address) = ?", f.asset.ID, strings.ToLower(w.WalletAddress))
	if row.Balance != "2256787" || row.Channel != pm.ChannelExchange || row.ServiceLinkID != f.link {
		t.Fatalf("ledger %+v", row)
	}
	// reconciliation: 225.6787 tokens backed by 1000 units
	run := f.e.Reconcile(context.Background(), f.asset, "test")
	if run.Result != pm.ReconMatched || run.Inventory != "774.3213" {
		t.Fatalf("reconciliation %+v", run)
	}
}

func TestSlowCreationThroughMockPartners(t *testing.T) {
	f := newFixture(t)
	f.deposit("400000")
	w := f.wallet("plm_usr_2")
	// 300000 buys ~1353 tokens: more than the 1000 units held
	o, err := f.e.PlaceExchangeCreation(context.Background(), f.link, ExchangeOrderRequest{WalletID: w.ID, AssetCode: "MTNN-T", Amount: "300000", ExternalOrderRef: "plm-2"})
	if err != nil {
		t.Fatal(err)
	}
	f.run(3)
	got := f.order(o.ID)
	if got.State != pm.StatePendingExecution || got.BatchID == "" {
		t.Fatalf("order %s batch %q", got.State, got.BatchID)
	}
	var b pm.NetBatch
	f.e.DB.First(&b, "id = ?", got.BatchID)
	// 299250 / 221 = 1354.07 units needed, 1000 in inventory: buy 355
	if b.Side != "BUY" || b.Quantity != "355" || b.Status != pm.BatchReleased {
		t.Fatalf("batch %+v", b)
	}
	var ins []pm.Instruction
	f.e.DB.Where("batch_id = ?", b.ID).Find(&ins)
	if len(ins) != 2 || ins[0].Status != pm.InstrAccepted || ins[1].Status != pm.InstrAccepted {
		t.Fatalf("instructions %+v", ins)
	}
	// the mock broker fills after 20s, the mock Custodian settles after 60s
	f.now = f.now.Add(30 * time.Second)
	f.run(2)
	if got = f.order(o.ID); got.State != pm.StateExecuted {
		t.Fatalf("after the fill: %s", got.State)
	}
	f.now = f.now.Add(60 * time.Second)
	f.run(6)
	got = f.order(o.ID)
	if got.State != pm.StateComplete || got.ExecutedPrice == "" {
		t.Fatalf("order %s (%s)", got.State, got.LastError)
	}
	exec := decimal.RequireFromString(got.ExecutedPrice)
	want := decimal.RequireFromString("299250").Div(exec).Truncate(4)
	if got.Quantity != want.String() || f.tokens(common.HexToAddress(w.WalletAddress)) != want.String() {
		t.Fatalf("quantity %s, want %s at %s", got.Quantity, want, exec)
	}
	// the Custodian's position grew by what settled
	pos, _ := f.e.Position(f.asset)
	if pos.RealUnitsHeld != "1355" || pos.Source != pm.PositionSettlement {
		t.Fatalf("position %+v", pos)
	}
	var events []pm.OrderEvent
	f.e.DB.Where("order_id = ?", o.ID).Order("id").Find(&events)
	var states []string
	for _, ev := range events {
		states = append(states, ev.State)
	}
	if strings.Join(states, ",") != "queued,pending-execution,pending-execution,executed,settlement_final,submitted,chain_final,complete" {
		t.Fatalf("timeline %v", states)
	}
}

func TestNetCreationAboveThresholdNeedsTwoApprovals(t *testing.T) {
	f := newFixture(t)
	f.e.DB.Model(&pm.Settings{}).Where("id = 1").Updates(map[string]interface{}{"net_creation_threshold": "50000", "net_creation_approvers": "a@trovotech.io,b@trovotech.io"})
	f.deposit("400000")
	w := f.wallet("plm_usr_3")
	o, _ := f.e.PlaceExchangeCreation(context.Background(), f.link, ExchangeOrderRequest{WalletID: w.ID, AssetCode: "MTNN-T", Amount: "300000"})
	f.run(2)
	got := f.order(o.ID)
	var b pm.NetBatch
	f.e.DB.First(&b, "id = ?", got.BatchID)
	if b.Status != pm.BatchAwaitingApproval {
		t.Fatalf("batch %s", b.Status)
	}
	f.e.DB.Create(&pm.BatchApproval{BatchID: b.ID, Approver: "a@trovotech.io"})
	f.e.DB.Create(&pm.BatchApproval{BatchID: b.ID, Approver: "outsider@trovotech.io"})
	f.run(1)
	f.e.DB.First(&b, "id = ?", b.ID)
	if b.Status != pm.BatchAwaitingApproval {
		t.Fatal("released with one listed approver")
	}
	f.e.DB.Create(&pm.BatchApproval{BatchID: b.ID, Approver: "b@trovotech.io"})
	f.run(1)
	f.e.DB.First(&b, "id = ?", b.ID)
	if b.Status != pm.BatchReleased {
		t.Fatalf("batch %s after two approvals", b.Status)
	}
}

func TestRejectedBatchRefundsExchange(t *testing.T) {
	f := newFixture(t)
	f.e.DB.Model(&pm.Settings{}).Where("id = 1").Update("net_creation_threshold", "1000")
	f.deposit("400000")
	w := f.wallet("plm_usr_4")
	o, _ := f.e.PlaceExchangeCreation(context.Background(), f.link, ExchangeOrderRequest{WalletID: w.ID, AssetCode: "MTNN-T", Amount: "300000"})
	f.run(2)
	f.e.DB.Model(&pm.NetBatch{}).Where("id = ?", f.order(o.ID).BatchID).Updates(map[string]interface{}{"status": pm.BatchRejected, "updated_at": f.now})
	f.run(1)
	if got := f.order(o.ID); got.State != pm.StateRejected {
		t.Fatalf("order %s", got.State)
	}
	var p pm.ExchangePartner
	f.e.DB.First(&p, "service_link_id = ?", f.link)
	if p.Balance != "400000" {
		t.Fatalf("balance %s after the refund", p.Balance)
	}
}

func TestExchangeNettedRedemption(t *testing.T) {
	f := newFixture(t)
	f.deposit("100000")
	w := f.wallet("plm_usr_5")
	o, _ := f.e.PlaceExchangeCreation(context.Background(), f.link, ExchangeOrderRequest{WalletID: w.ID, AssetCode: "MTNN-T", Amount: "22100"})
	f.run(6)
	held := f.tokens(common.HexToAddress(w.WalletAddress))
	if f.order(o.ID).State != pm.StateComplete {
		t.Fatal("creation did not complete")
	}
	if _, err := f.e.PlaceExchangeRedemption(context.Background(), f.link, ExchangeOrderRequest{WalletID: w.ID, AssetCode: "MTNN-T", Quantity: "1000"}); err == nil {
		t.Fatal("redeemed more than held")
	}
	r, err := f.e.PlaceExchangeRedemption(context.Background(), f.link, ExchangeOrderRequest{WalletID: w.ID, AssetCode: "MTNN-T", Quantity: "50", ExternalOrderRef: "plm-r1"})
	if err != nil {
		t.Fatal(err)
	}
	f.run(8)
	got := f.order(r.ID)
	if got.State != pm.StateComplete || got.Path != pm.PathNetted {
		t.Fatalf("redemption %s %s: %s", got.State, got.Path, got.LastError)
	}
	// 50 x 221 = 11050, less 27.63 fee
	if got.NetAmount != "11022.37" {
		t.Fatalf("net %s", got.NetAmount)
	}
	if left := decimal.RequireFromString(held).Sub(decimal.NewFromInt(50)).String(); f.tokens(common.HexToAddress(w.WalletAddress)) != left {
		t.Fatalf("wallet holds %s, want %s", f.tokens(common.HexToAddress(w.WalletAddress)), left)
	}
	if f.tokens(issuing) != "0" {
		t.Fatalf("issuing Safe kept %s", f.tokens(issuing))
	}
	var p pm.ExchangePartner
	f.e.DB.First(&p, "service_link_id = ?", f.link)
	if want := decimal.RequireFromString("77900").Add(decimal.RequireFromString("11022.37")).String(); p.Balance != want {
		t.Fatalf("balance %s, want %s", p.Balance, want)
	}
	var hooks int64
	f.e.DB.Model(&pm.WebhookDelivery{}).Where("event = ? AND reference = ?", "redemption.settled", r.ID).Count(&hooks)
	if hooks != 1 {
		t.Fatal("no redemption.settled webhook")
	}
}

func TestReconciliationHaltsAndResumes(t *testing.T) {
	f := newFixture(t)
	f.deposit("400000")
	w := f.wallet("plm_usr_6")
	f.e.PlaceExchangeCreation(context.Background(), f.link, ExchangeOrderRequest{WalletID: w.ID, AssetCode: "MTNN-T", Amount: "200000"})
	f.run(6)
	// the Custodian now reports fewer units than tokens issued
	f.e.DB.Create(&pm.CustodianPosition{AssetID: f.asset.ID, RealUnitsHeld: "500", AsOf: f.now, Source: pm.PositionFeed})
	f.e.DB.Model(&pm.Custodian{}).Where("code = ?", "CUSTA").Update("mode", pm.ModeManual) // no position query
	run := f.e.Reconcile(context.Background(), f.asset, "test")
	if run.Result != pm.ReconDrift || !run.Halted {
		t.Fatalf("run %+v", run)
	}
	var a pm.Asset
	f.e.DB.First(&a, "id = ?", f.asset.ID)
	if a.Status != pm.AssetHalted {
		t.Fatal("not halted")
	}
	if _, err := f.e.PlaceExchangeCreation(context.Background(), f.link, ExchangeOrderRequest{WalletID: w.ID, AssetCode: "MTNN-T", Amount: "5000"}); err == nil {
		t.Fatal("order accepted on a halted asset")
	}
	if _, ok := f.e.ResumeAsset(context.Background(), &a, "ops"); ok {
		t.Fatal("resumed while drifting")
	}
	f.e.DB.Create(&pm.CustodianPosition{AssetID: f.asset.ID, RealUnitsHeld: "1000", AsOf: f.now, Source: pm.PositionFeed})
	f.e.DB.First(&a, "id = ?", f.asset.ID)
	if _, ok := f.e.ResumeAsset(context.Background(), &a, "ops"); !ok {
		t.Fatal("not resumed after matching")
	}
}

func TestPartnerEventsAreIdempotent(t *testing.T) {
	f := newFixture(t)
	body := []byte(`{"asOf":"2026-09-29","positions":[{"assetCode":"MTNN-T","unitsHeld":"1200"}]}`)
	if dup, err := f.e.ReceivePartnerEvent("CUSTODIAN:CUSTA", EventPositionFeed, "feed-1", body); err != nil || dup {
		t.Fatalf("first: dup=%v err=%v", dup, err)
	}
	if dup, _ := f.e.ReceivePartnerEvent("CUSTODIAN:CUSTA", EventPositionFeed, "feed-1", body); !dup {
		t.Fatal("redelivery processed twice")
	}
	var n int64
	f.e.DB.Model(&pm.CustodianPosition{}).Where("source = ?", pm.PositionFeed).Count(&n)
	if n != 2 { // the fixture's and this feed's
		t.Fatalf("%d feed positions", n)
	}
	// a Custodian cannot report another Custodian's asset
	f.e.DB.Create(&pm.Custodian{CustodianID: 2, Code: "CUSTB", Mode: pm.ModeMock})
	f.e.ReceivePartnerEvent("CUSTODIAN:CUSTB", EventPositionFeed, "feed-b", []byte(`{"asOf":"2026-09-29","positions":[{"assetCode":"MTNN-T","unitsHeld":"1"}]}`))
	pos, _ := f.e.Position(f.asset)
	if pos.RealUnitsHeld != "1200" {
		t.Fatalf("position %s", pos.RealUnitsHeld)
	}
}

func TestProvisioningFieldErrors(t *testing.T) {
	f := newFixture(t)
	_, _, err := f.e.ProvisionWallet(context.Background(), f.link, ProvisionRequest{ExternalUserRef: "plm_x", LegalName: "F A", ResidencyCountry: "NGA", Nationality: "NGA", NDPAConsent: boolp(true)})
	pe, ok := AsError(err)
	if !ok || pe.Field != "taxIdentifier" || pe.Status != 400 {
		t.Fatalf("err %v", err)
	}
	var rejected pm.PartnerWallet
	if f.e.DB.First(&rejected, "external_user_ref = ? AND status = ?", "plm_x", "rejected").Error != nil {
		t.Fatal("rejected request not kept")
	}
	_, _, err = f.e.ProvisionWallet(context.Background(), f.link, ProvisionRequest{ExternalUserRef: "plm_y", LegalName: "F A", TaxIdentifier: "1", ResidencyCountry: "NGA", Nationality: "NGA"})
	if pe, _ := AsError(err); pe == nil || pe.Field != "ndpaConsent" {
		t.Fatalf("missing consent accepted: %v", err)
	}
	w1 := f.wallet("plm_z")
	w2 := f.wallet("plm_z")
	if w1.ID != w2.ID {
		t.Fatal("same customer got two wallets")
	}
}

func TestExchangeSignature(t *testing.T) {
	f := newFixture(t)
	var p pm.ExchangePartner
	f.e.DB.First(&p, "service_link_id = ?", f.link)
	body := []byte(`{"walletId":"w"}`)
	ts := fmt.Sprint(f.now.Unix())
	if err := f.e.VerifyExchangeSignature(&p, ts, body, signFor("s3cret", ts, body)); err != nil {
		t.Fatal(err)
	}
	if err := f.e.VerifyExchangeSignature(&p, ts, body, signFor("wrong", ts, body)); err == nil {
		t.Fatal("bad signature accepted")
	}
	old := fmt.Sprint(f.now.Add(-10 * time.Minute).Unix())
	if err := f.e.VerifyExchangeSignature(&p, old, body, signFor("s3cret", old, body)); err == nil {
		t.Fatal("stale timestamp accepted")
	}
}

func TestDividendDistribution(t *testing.T) {
	f := newFixture(t)
	f.deposit("200000")
	w := f.wallet("plm_usr_7") // resident, with a tax ID
	f.e.PlaceExchangeCreation(context.Background(), f.link, ExchangeOrderRequest{WalletID: w.ID, AssetCode: "MTNN-T", Amount: "22100"})
	// an app holder (no user record: resident rate) gets tokens through P2P
	app := common.HexToAddress("0x0000000000000000000000000000000000000a11")
	f.run(6)
	f.chain.mu.Lock()
	f.chain.block++
	f.chain.move(token, common.HexToAddress(w.WalletAddress), app, big.NewInt(200000), "0xp2p", 0) // 20 tokens
	f.chain.mu.Unlock()
	f.run(1)
	ca, err := f.e.DeclareCorporateAction(f.asset, CorporateActionNotice{AssetCode: "MTNN-T", EventType: "DIVIDEND", RecordDate: "2026-09-29", PayDate: "2026-09-30", AmountPerUnit: "40", Currency: "NGN"}, pm.SourceManualCA, "", "ops")
	if err != nil {
		t.Fatal(err)
	}
	f.e.SnapshotDueActions(context.Background())
	f.e.DB.First(ca, "id = ?", ca.ID)
	if ca.Status != pm.ActionSnapshotted || ca.HolderCount != 2 {
		t.Fatalf("snapshot %+v", ca)
	}
	var ents []pm.Entitlement
	f.e.DB.Where("corporate_action_id = ?", ca.ID).Order("units").Find(&ents)
	// app: 20 x 40 = 800, 10% WHT 80; exchange: (99.4382 - 20) x 40
	if ents[0].Units != "20" || ents[0].GrossAmount != "800" || ents[0].WHTAmount != "80" || ents[0].NetAmount != "720" || ents[0].Channel != pm.ChannelApp {
		t.Fatalf("app entitlement %+v", ents[0])
	}
	// approvals for the snapshot's checksum
	f.e.DB.Create(&pm.DividendApproval{CorporateActionID: ca.ID, Approver: "a@trovotech.io", Checksum: ca.SnapshotChecksum})
	f.e.DB.Create(&pm.DividendApproval{CorporateActionID: ca.ID, Approver: "b@trovotech.io", Checksum: "stale"})
	f.e.AdvanceDividends(context.Background())
	f.e.DB.First(ca, "id = ?", ca.ID)
	if ca.Status != pm.ActionSnapshotted {
		t.Fatal("approved with a stale approval")
	}
	f.e.DB.Create(&pm.DividendApproval{CorporateActionID: ca.ID, Approver: "c@trovotech.io", Checksum: ca.SnapshotChecksum})
	f.fund(treasury, cngn, "10000") // the dividend cash from the Custodian
	for i := 0; i < 6; i++ {
		f.e.AdvanceDividends(context.Background())
	}
	f.e.DB.First(ca, "id = ?", ca.ID)
	if ca.Status != pm.ActionDistributed || ca.WHTTxHash == "" {
		t.Fatalf("distribution %s (%s)", ca.Status, ca.Note)
	}
	if b, _ := f.chain.BalanceOf(context.Background(), cngn, app); decimal.NewFromBigInt(b, -6).String() != "720" {
		t.Fatalf("app holder paid %s", decimal.NewFromBigInt(b, -6))
	}
	if b, _ := f.chain.BalanceOf(context.Background(), cngn, whtW); !decimal.NewFromBigInt(b, -6).Equal(decimal.RequireFromString(ca.WHTAmount)) {
		t.Fatalf("WHT wallet got %s, want %s", decimal.NewFromBigInt(b, -6), ca.WHTAmount)
	}
	var hook pm.WebhookDelivery
	if f.e.DB.First(&hook, "event = ? AND needs_confirmation = ?", "dividend.paid", true).Error != nil {
		t.Fatal("no dividend.paid webhook")
	}
	var payload map[string]interface{}
	json.Unmarshal([]byte(hook.Payload), &payload)
	if payload["walletId"] != w.ID || payload["withholdingTaxApplied"] == "" {
		t.Fatalf("payload %v", payload)
	}
	if _, err := f.e.Confirm(f.link, ConfirmationRequest{Event: "dividend.paid", WalletID: w.ID}); err != nil {
		t.Fatal(err)
	}
	if _, err := f.e.Confirm(f.link, ConfirmationRequest{Event: "dividend.paid", WalletID: w.ID}); err == nil {
		t.Fatal("confirmed twice")
	}
}

func signFor(secret, ts string, body []byte) string {
	var buf bytes.Buffer
	buf.Write(body)
	return partners.Sign(secret, ts, buf.Bytes())
}

// testDB is a fresh SQLite database, or a fresh schema of the Postgres
// database in PUBLIC_MARKETS_TEST_POSTGRES (a DSN) when it is set.
func testDB(t *testing.T) *gorm.DB {
	t.Helper()
	cfg := &gorm.Config{Logger: logger.Discard}
	dsn := os.Getenv("PUBLIC_MARKETS_TEST_POSTGRES")
	if dsn == "" {
		db, err := gorm.Open(sqlite.Open("file:"+t.TempDir()+"/pm.db"), cfg)
		if err != nil {
			t.Fatal(err)
		}
		return db
	}
	schema := fmt.Sprintf("pmtest_%d", time.Now().UnixNano())
	admin, err := gorm.Open(postgres.Open(dsn), cfg)
	if err != nil {
		t.Fatal(err)
	}
	if err := admin.Exec("CREATE SCHEMA " + schema).Error; err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		admin.Exec("DROP SCHEMA " + schema + " CASCADE")
		if sqlDB, err := admin.DB(); err == nil {
			sqlDB.Close()
		}
	})
	db, err := gorm.Open(postgres.Open(dsn+" search_path="+schema), cfg)
	if err != nil {
		t.Fatal(err)
	}
	return db
}

// Trovo Manager records MANUAL partners' outcomes and declared corporate
// actions as unprocessed partner events ("MANUAL:<admin>"); the engine
// applies them like webhooks.
func TestManualPartnersThroughTrovoManager(t *testing.T) {
	f := newFixture(t)
	f.e.DB.Model(&pm.Custodian{}).Where("code = ?", "CUSTA").Update("mode", pm.ModeManual)
	f.e.DB.Model(&pm.DealingMember{}).Where("code = ?", "BP2").Update("mode", pm.ModeManual)
	f.deposit("400000")
	w := f.wallet("plm_usr_m")
	o, _ := f.e.PlaceExchangeCreation(context.Background(), f.link, ExchangeOrderRequest{WalletID: w.ID, AssetCode: "MTNN-T", Amount: "300000"})
	f.run(3)
	var ins []pm.Instruction
	f.e.DB.Where("batch_id = ?", f.order(o.ID).BatchID).Order("kind").Find(&ins)
	if len(ins) != 2 || ins[0].Status != pm.InstrAccepted || ins[1].Status != pm.InstrAccepted {
		t.Fatalf("instructions %+v", ins)
	}
	var dm, ci pm.Instruction
	for _, in := range ins {
		if in.Kind == pm.InstrDealingOrder {
			dm = in
		} else {
			ci = in
		}
	}
	record := func(kind, ref string, body map[string]string) {
		payload, _ := json.Marshal(body)
		f.e.DB.Create(&pm.PartnerEvent{Source: "MANUAL:ops@trovotech.io", Kind: kind, Reference: ref, Payload: string(payload), CreatedAt: f.now})
	}
	record(EventExecution, "tm:"+dm.ID+":FILLED", map[string]string{"orderId": dm.ID, "status": "FILLED", "executedQuantity": dm.Quantity, "executedPrice": "221.10"})
	f.run(2)
	if got := f.order(o.ID); got.State != pm.StateExecuted {
		t.Fatalf("after the recorded fill: %s", got.State)
	}
	record(EventSettlement, "tm:"+ci.ID+":settlement_final", map[string]string{"instructionId": ci.ID, "status": "settlement_final",
		"settledQuantity": ci.Quantity, "settlementDate": "2026-10-02", "custodianReference": "CSD-77"})
	f.run(6)
	if got := f.order(o.ID); got.State != pm.StateComplete || got.ExecutedPrice != "221.1" {
		t.Fatalf("after the recorded settlement: %s at %s (%s)", got.State, got.ExecutedPrice, got.LastError)
	}
	record(EventCorporateAction, "tm:MTNN-T:DIVIDEND:2026-10-10", map[string]string{"assetCode": "MTNN-T", "eventType": "DIVIDEND",
		"recordDate": "2026-10-10", "payDate": "2026-10-24", "amountPerUnit": "1.5", "currency": "NGN"})
	f.e.ProcessPendingEvents()
	var ca pm.CorporateAction
	if f.e.DB.First(&ca, "asset_id = ?", f.asset.ID).Error != nil || ca.Source != pm.SourceManualCA || ca.DeclaredBy != "ops@trovotech.io" || ca.Status != pm.ActionAnnounced {
		t.Fatalf("declared %+v", ca)
	}
	var pending int64
	f.e.DB.Model(&pm.PartnerEvent{}).Where("processed_at IS NULL").Count(&pending)
	if pending != 0 {
		t.Fatalf("%d events left unprocessed", pending)
	}
}
