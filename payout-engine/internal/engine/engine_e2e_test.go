package engine

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"trovo-payout-engine/internal/bus"
	"trovo-payout-engine/internal/chain"
	"trovo-payout-engine/internal/config"
	"trovo-payout-engine/internal/store"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type staticSigners []config.Signer

func (s staticSigners) Signers(context.Context) ([]config.Signer, error) { return s, nil }

type fakePush struct {
	mu   sync.Mutex
	sent []string
}

func (f *fakePush) Send(_ context.Context, token, title, body string, _ map[string]string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.sent = append(f.sent, token+": "+body)
	return nil
}

// TestPayoutOnLocalChain runs a payout end to end against the local stack
// (paymaster/contracts scripts/local-stack.js and market/contracts
// scripts/local-market.js on a hardhat node; AA_LOCAL_STACK / AA_LOCAL_MARKET
// point at their deployment files): a fresh asset token with holders, one
// of them a platform wallet, a 3-of-4 payout Safe paying cNGN, preparation
// (replayed from the token's transfers), re-preparation after a transfer
// (only new blocks scanned), the funding check refusing an unfunded Safe and
// a tampered schedule, payment in MultiSend batches with a failing transfer
// isolated, retry, push notifications, recovery of a dropped batch and a
// sweep of what is left.
func TestPayoutOnLocalChain(t *testing.T) {
	stackFile, marketFile := os.Getenv("AA_LOCAL_STACK"), os.Getenv("AA_LOCAL_MARKET")
	if stackFile == "" || marketFile == "" {
		t.Skip("AA_LOCAL_STACK / AA_LOCAL_MARKET not set")
	}
	read := func(f string) map[string]string {
		raw, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		m := map[string]string{}
		if err := json.Unmarshal(raw, &m); err != nil {
			t.Fatal(err)
		}
		return m
	}
	stack, market := read(stackFile), read(marketFile)
	ctx := context.Background()
	client, err := ethclient.Dial("http://127.0.0.1:8545")
	if err != nil {
		t.Fatal(err)
	}
	chainID, _ := client.ChainID(ctx)
	admin, _ := crypto.HexToECDSA("ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
	adminAddr := crypto.PubkeyToAddress(admin.PublicKey)

	send := func(to *common.Address, value *big.Int, data []byte) *types.Receipt {
		t.Helper()
		nonce, _ := client.PendingNonceAt(ctx, adminAddr)
		gp, _ := client.SuggestGasPrice(ctx)
		msg := ethereum.CallMsg{From: adminAddr, To: to, Value: value, Data: data}
		gas, err := client.EstimateGas(ctx, msg)
		if err != nil {
			t.Fatalf("estimate: %v", err)
		}
		tx := types.NewTx(&types.LegacyTx{Nonce: nonce, To: to, Value: value, Gas: gas * 2, GasPrice: gp, Data: data})
		signed, _ := types.SignTx(tx, types.NewEIP155Signer(chainID), admin)
		if err := client.SendTransaction(ctx, signed); err != nil {
			t.Fatal(err)
		}
		r, err := chain.WaitReceipt(ctx, client, signed.Hash(), time.Minute)
		if err != nil || r.Status != 1 {
			t.Fatalf("tx failed: %v", err)
		}
		return r
	}
	erc20, _ := abi.JSON(strings.NewReader(`[{"name":"mint","type":"function","inputs":[{"type":"address"},{"type":"uint256"}],"outputs":[]},
		{"name":"transfer","type":"function","inputs":[{"type":"address"},{"type":"uint256"}],"outputs":[{"type":"bool"}]}]`))
	mint := func(token, to common.Address, amount *big.Int) {
		data, _ := erc20.Pack("mint", to, amount)
		send(&token, big.NewInt(0), data)
	}
	units := func(n int64, dec int) *big.Int {
		return new(big.Int).Mul(big.NewInt(n), new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(dec)), nil))
	}

	// a fresh asset token (18 decimals), owned (minted) by the admin
	var art struct {
		ABI      json.RawMessage
		Bytecode string
	}
	raw, _ := os.ReadFile(market["TOKENIZED_ASSET_ARTIFACT"])
	json.Unmarshal(raw, &art)
	tokenABI, _ := abi.JSON(strings.NewReader(string(art.ABI)))
	ctorArgs, _ := tokenABI.Pack("", "Test Asset", "TST", uint8(18), adminAddr)
	deployed := send(nil, big.NewInt(0), append(hexutil.MustDecode(art.Bytecode), ctorArgs...))
	token := deployed.ContractAddress
	cngn := common.HexToAddress(market["MOCK_CNGN_ADDRESS"]) // 6 decimals

	// the payout Safe: 3-of-4
	keys := make([]config.Signer, 4)
	owners := make([]common.Address, 4)
	for i := range keys {
		k, _ := crypto.GenerateKey()
		keys[i] = config.Signer{Key: k, Address: crypto.PubkeyToAddress(k.PublicKey)}
		owners[i] = keys[i].Address
	}
	safeABI, _ := abi.JSON(strings.NewReader(`[{"name":"setup","type":"function","inputs":[{"name":"_owners","type":"address[]"},{"name":"_threshold","type":"uint256"},{"name":"to","type":"address"},{"name":"data","type":"bytes"},{"name":"fallbackHandler","type":"address"},{"name":"paymentToken","type":"address"},{"name":"payment","type":"uint256"},{"name":"paymentReceiver","type":"address"}],"outputs":[]}]`))
	factoryABI, _ := abi.JSON(strings.NewReader(`[{"name":"createProxyWithNonce","type":"function","inputs":[{"name":"_singleton","type":"address"},{"name":"initializer","type":"bytes"},{"name":"saltNonce","type":"uint256"}],"outputs":[{"name":"proxy","type":"address"}]}]`))
	setup, _ := safeABI.Pack("setup", owners, big.NewInt(3), common.Address{}, []byte{}, common.Address{}, common.Address{}, big.NewInt(0), common.Address{})
	createData, _ := factoryABI.Pack("createProxyWithNonce", common.HexToAddress(stack["SAFE_SINGLETON_ADDRESS"]), setup, big.NewInt(time.Now().UnixNano()))
	factory := common.HexToAddress(stack["SAFE_PROXY_FACTORY_ADDRESS"])
	out, err := client.CallContract(ctx, ethereum.CallMsg{From: adminAddr, To: &factory, Data: createData}, nil)
	if err != nil {
		t.Fatal(err)
	}
	payoutSafe := common.BytesToAddress(out[12:32])
	send(&factory, big.NewInt(0), createData)
	send(&keys[0].Address, units(1, 18), nil) // the executor's gas

	// holders: alice 100, bob 50, carol 30 (no Trovo account), the issuing
	// profile's Safe 20 (a platform wallet)
	randomAddress := func() common.Address {
		k, _ := crypto.GenerateKey()
		return crypto.PubkeyToAddress(k.PublicKey)
	}
	alice, bob, carol, issuing, dave := randomAddress(), randomAddress(), randomAddress(), randomAddress(), randomAddress()
	for addr, n := range map[common.Address]int64{alice: 100, bob: 50, carol: 30, issuing: 20} {
		mint(token, addr, units(n, 18))
	}

	db, err := gorm.Open(sqlite.Open("file:"+t.TempDir()+"/payouts.db"), &gorm.Config{Logger: logger.Discard})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(store.AllModels()...); err != nil {
		t.Fatal(err)
	}
	code, name, contract := "TST", "Test Asset", token.Hex()
	db.Create(&store.TokenizedAsset{ID: "ta1", CreatedAt: time.Now().Add(-time.Minute), AssetCode: &code, AssetName: &name, ContractAddress: &contract})
	pushToken := "alice-device"
	db.Create(&store.User{ID: "1", Username: "alice", PushNotificationToken: &pushToken})
	db.Create(&store.User{ID: "2", Username: "atprofile"})
	db.Create(&store.UserWallet{ID: alice.Hex(), UserID: "1", Alias: "alice"})
	db.Create(&store.UserWallet{ID: issuing.Hex(), UserID: "2", Alias: "issuing"})
	db.Create(&store.ProceedPayout{ID: 1, TokenizedAssetID: "ta1", Batch: "TST-d1", DistributionID: "d1", Status: store.StatusPrepareRequested,
		PayoutAssetCode: "CNGN", PayoutContractAddress: cngn.Hex(), TotalAmount: "2000", ApprovalsRequired: 2})

	sweepTo, feeWallet, vatWallet := randomAddress(), randomAddress(), randomAddress()
	t.Setenv("VAT_WALLET", "")
	db.Create(&store.ServiceFee{ID: store.FeeServiceID, FeeWalletSecretKey: feeWallet.Hex()})
	db.Create(&store.ServiceFee{ID: "VAT", FeeWalletSecretKey: vatWallet.Hex()})
	db.Create(&store.CountryConfig{CountryCode: "NG", VATPercent: 7.5})
	cfg := &config.Config{PayoutSafe: payoutSafe, MultiSendCallOnly: common.HexToAddress(stack["SAFE_MULTISEND_CALL_ONLY_ADDRESS"]), IssuingProfile: "atprofile",
		SweepAddress: sweepTo, Confirmations: 0, LogChunk: 1000, CatchUpWithin: 0, BatchSize: 2, MaxBatchGas: 30_000_000,
		MinPayoutUnits: big.NewInt(1), GasPriceMultiplier: 1.25, PollInterval: time.Second, ReceiptTimeout: time.Minute, Instance: "test"}
	pusher := &fakePush{}
	e := &Engine{DB: db, Chain: client, ChainID: chainID, Cfg: cfg, Signers: staticSigners(keys[:3]), Push: pusher, Bus: bus.New("", "", false)}
	if !e.claim() {
		t.Fatal("claim")
	}
	payout := func() store.ProceedPayout {
		var p store.ProceedPayout
		db.First(&p, 1)
		return p
	}
	tickUntil := func(status string) store.ProceedPayout {
		t.Helper()
		for i := 0; i < 30; i++ {
			if err := e.Tick(ctx); err != nil {
				t.Fatal(err)
			}
			if p := payout(); p.Status == status {
				return p
			}
		}
		p := payout()
		t.Fatalf("payout is %s (%s), want %s", p.Status, p.Note, status)
		return p
	}
	item := func(addr common.Address) store.PayoutItem {
		var it store.PayoutItem
		db.Where("proceed_payout_id = 1 AND LOWER(beneficiary_address) = ?", strings.ToLower(addr.Hex())).First(&it)
		return it
	}
	cngnOf := func(a common.Address) *big.Int {
		b, err := chain.BalanceOf(ctx, client, cngn, a, nil)
		if err != nil {
			t.Fatal(err)
		}
		return b
	}

	// preparation: 200 tokens, 2000 CNGN -> 10 CNGN per token; the
	// issuing Safe's 200 CNGN share is retained
	p := tickUntil(store.StatusLocked)
	if p.SupplyUnits != units(200, 18).String() || p.AmountPerToken != "10" || p.PayableUnits != units(1800, 6).String() || p.RetainedUnits != units(200, 6).String() {
		t.Fatalf("locked payout %+v", p)
	}
	if it := item(alice); it.AmountUnits != units(1000, 6).String() || it.Username != "alice" || it.Status != store.ItemPending {
		t.Fatalf("alice %+v", it)
	}
	if it := item(issuing); it.Status != store.ItemExcluded || !strings.Contains(it.Reason, "issuing") {
		t.Fatalf("issuing Safe %+v", it)
	}
	firstSnapshot := p.SnapshotBlock

	// new tokens for dave after the snapshot, then re-preparation (an old
	// approval is dropped): only the new blocks are scanned
	mint(token, dave, units(50, 18)) // supply 250
	db.Create(&store.PayoutApproval{ProceedPayoutID: 1, AdminEmail: "a@trovo", ScheduleChecksum: p.ScheduleChecksum})
	// an admin sets the processing fee meanwhile: 5%, capped at 50 CNGN,
	// plus 7.5% VAT (NG) on it - 53.75 CNGN out of the 2000
	db.Model(&store.ProceedPayout{}).Where("id = 1").Updates(map[string]interface{}{"status": store.StatusPrepareRequested, "fee_type": store.FeePercent, "fee_value": "5", "fee_cap": "50"})
	p = tickUntil(store.StatusLocked)
	var approvals int64
	db.Model(&store.PayoutApproval{}).Where("proceed_payout_id = 1").Count(&approvals)
	if p.SnapshotBlock <= firstSnapshot || p.SupplyUnits != units(250, 18).String() || item(dave).AmountUnits != "389250000" || approvals != 0 ||
		p.FeeUnits != units(50, 6).String() || p.VatUnits != "3750000" || p.AmountPerToken != "7.785" || item(feeWallet).Kind != store.KindFee || item(vatWallet).AmountUnits != "3750000" {
		t.Fatalf("re-prepared payout %+v (dave %+v, approvals %d)", p, item(dave), approvals)
	}
	var idx store.TokenIndex
	db.First(&idx, "token = ?", strings.ToLower(token.Hex()))
	if idx.ScannedBlock != p.SnapshotBlock {
		t.Fatalf("token index at %d, snapshot %d", idx.ScannedBlock, p.SnapshotBlock)
	}

	// funding check: refused without approvals, an unfunded Safe or a
	// tampered schedule
	requestFunding := func() {
		db.Model(&store.ProceedPayout{}).Where("id = 1").Update("status", store.StatusFundingCheckRequested)
	}
	requestFunding()
	if p = tickUntil(store.StatusApproved); !strings.Contains(p.Note, "approvals") {
		t.Fatalf("note %q", p.Note)
	}
	db.Create(&store.PayoutApproval{ProceedPayoutID: 1, AdminEmail: "a@trovo", ScheduleChecksum: p.ScheduleChecksum})
	db.Create(&store.PayoutApproval{ProceedPayoutID: 1, AdminEmail: "b@trovo", ScheduleChecksum: p.ScheduleChecksum})
	requestFunding()
	if p = tickUntil(store.StatusApproved); !strings.Contains(p.Note, "send") {
		t.Fatalf("note %q", p.Note)
	}
	mint(cngn, payoutSafe, units(2000, 6))
	db.Model(&store.PayoutItem{}).Where("id = ?", item(bob).ID).Update("amount_units", units(999, 6).String())
	requestFunding()
	if p = tickUntil(store.StatusApproved); !strings.Contains(p.Note, "changed") {
		t.Fatalf("note %q", p.Note)
	}
	db.Model(&store.PayoutItem{}).Where("id = ?", item(bob).ID).Update("amount_units", "389250000")
	requestFunding()
	tickUntil(store.StatusPaying)

	// paying 7.785 CNGN per token: alice 778.5, bob 389.25, carol 233.55,
	// dave 389.25, and the fee and VAT; carol's
	// transfer is made to fail (more than the Safe holds) and isolated
	db.Model(&store.PayoutItem{}).Where("id = ?", item(carol).ID).Update("amount_units", units(1_000_000, 6).String())
	p = tickUntil(store.StatusCompletedWithFailures)
	if it := item(carol); it.Status != store.ItemFailed || it.CannotReceiveAsset != 1 {
		t.Fatalf("carol %+v", it)
	}
	var all []store.PayoutItem
	db.Where("proceed_payout_id = 1").Find(&all)
	for _, it := range all {
		t.Logf("%s %s %s %s", it.BeneficiaryAddress, it.Status, it.AmountUnits, it.Reason)
	}
	for addr, want := range map[common.Address]int64{alice: 778_500000, bob: 389_250000, dave: 389_250000, feeWallet: 50_000000, vatWallet: 3_750000} {
		if got := cngnOf(addr); got.Cmp(big.NewInt(want)) != 0 {
			t.Fatalf("%s holds %v, want %d CNGN", addr.Hex(), got, want)
		}
		if it := item(addr); it.Status != store.ItemPaid || it.TxHash == "" {
			t.Fatalf("item %+v", it)
		}
	}
	var batches []store.PayoutBatch
	db.Where("proceed_payout_id = 1 AND status = ?", store.BatchMined).Find(&batches)
	if len(batches) != 3 || p.PaidCount != 3 || p.FailedCount != 1 {
		t.Fatalf("%d batches; payout %+v", len(batches), p)
	}

	// retry the failed one (its amount put right), as an admin would
	db.Model(&store.PayoutItem{}).Where("id = ?", item(carol).ID).Updates(map[string]interface{}{"amount_units": "233550000", "status": store.ItemPending, "cannot_receive_asset": 0})
	db.Model(&store.ProceedPayout{}).Where("id = 1").Update("status", store.StatusPaying)
	p = tickUntil(store.StatusCompleted)
	if cngnOf(carol).Cmp(big.NewInt(233_550000)) != 0 || p.PaidUnits != "1790550000" {
		t.Fatalf("carol %v; payout %+v", cngnOf(carol), p)
	}

	// the fee and its VAT are in the fee report
	var fees []store.FeeCollection
	db.Order("fee_type").Find(&fees)
	if len(fees) != 2 || fees[0].FeeType != store.FeeTypePayout || fees[0].Amount != 50 || fees[1].Amount != 3.75 || fees[0].TransactionHash == nil {
		t.Fatalf("fee collections %+v", fees)
	}

	// alice (the only Trovo user paid) is told once
	e.Tick(ctx)
	if len(pusher.sent) != 1 || !strings.Contains(pusher.sent[0], "778.5 CNGN") || !strings.Contains(pusher.sent[0], "Test Asset") {
		t.Fatalf("notifications %v", pusher.sent)
	}

	// a paid schedule is never prepared again (it would pay holders twice)
	var paidItems int64
	db.Model(&store.PayoutItem{}).Where("proceed_payout_id = 1 AND status = ?", store.ItemPaid).Count(&paidItems)
	db.Model(&store.ProceedPayout{}).Where("id = 1").Update("status", store.StatusPrepareRequested)
	var again store.ProceedPayout
	db.First(&again, 1)
	if err := e.prepare(ctx, &again); err != nil {
		t.Fatal(err)
	}
	var stillPaid int64
	db.Model(&store.PayoutItem{}).Where("proceed_payout_id = 1 AND status = ?", store.ItemPaid).Count(&stillPaid)
	db.First(&again, 1)
	if again.Status != store.StatusPaused || stillPaid != paidItems || paidItems == 0 {
		t.Fatalf("re-preparing a paid payout: %s, %d of %d paid items left", again.Status, stillPaid, paidItems)
	}
	db.Model(&store.ProceedPayout{}).Where("id = 1").Update("status", store.StatusCompleted)

	// a batch recorded but never mined, whose Safe nonce was used since,
	// is dropped and its items paid again later
	db.Create(&store.ProceedPayout{ID: 2, TokenizedAssetID: "ta1", Batch: "TST-d2", Status: store.StatusPaused, PayoutContractAddress: cngn.Hex()})
	stale := store.PayoutBatch{ProceedPayoutID: 2, Status: store.BatchSubmitted, SafeNonce: 0, TxHash: common.HexToHash("0x01").Hex(), RawTx: "0x00", CreatedAt: time.Now()}
	db.Create(&stale)
	db.Create(&store.PayoutItem{ID: "x1", ProceedPayoutID: 2, BeneficiaryAddress: dave.Hex(), Status: store.ItemQueued, BatchID: stale.ID, AmountUnits: "1"})
	if err := e.recoverBatches(ctx); err != nil {
		t.Fatal(err)
	}
	var it store.PayoutItem
	db.First(&it, "id = ?", "x1")
	db.First(&stale, stale.ID)
	if stale.Status != store.BatchDropped || it.Status != store.ItemPending {
		t.Fatalf("stale batch %s, item %s", stale.Status, it.Status)
	}
	db.Model(&store.ProceedPayout{}).Where("id = 2").Update("status", store.StatusCancelled)

	// sweep what is left (the issuing Safe's retained share and dust)
	left := cngnOf(payoutSafe)
	db.Model(&store.EngineState{}).Where("id = 1").Updates(map[string]interface{}{"sweep_token": cngn.Hex(), "sweep_requested_by": "a@trovo"})
	if err := e.sweep(ctx); err != nil {
		t.Fatal(err)
	}
	var st store.EngineState
	db.First(&st, 1)
	if cngnOf(payoutSafe).Sign() != 0 || cngnOf(sweepTo).Cmp(left) != 0 || st.SweepToken != "" || !strings.HasPrefix(st.SweepResult, "swept") {
		t.Fatalf("sweep: left %v, swept %v, state %+v", cngnOf(payoutSafe), cngnOf(sweepTo), st)
	}
	fmt.Println("payout safe", payoutSafe.Hex(), "batches", len(batches))
}
