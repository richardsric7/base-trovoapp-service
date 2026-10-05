package users

import (
	"context"
	"encoding/json"
	"math/big"
	"os"
	"strings"
	"testing"
	"time"

	"trovo-wallet-api/internal/aa"
	"trovo-wallet-api/internal/evmkeypair"
	"trovo-wallet-api/internal/gnosissafe"
	"trovo-wallet-api/internal/offerbook"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/shopspring/decimal"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// TestSwapsOnLocalChain runs market-making offers and swaps against the
// real contracts (see TestTokenizationOnLocalChain for the local stack and
// the AA_LOCAL_STACK / AA_LOCAL_MARKET settings): a seller wallet places
// two SELL offers, the offer book index picks them up, a swap is routed
// cheapest-first across both and filled by the buyer's wallet operation,
// the index records the fills (order book, trade candles), a BUY offer
// shows as a bid, a swap larger than the book is refused, and a cancelled
// offer leaves the book. It also leaves an ETH batch send in the chain and
// writes its transaction to $AA_SAFE_ETH_TX_FILE (when set) for
// payment-history-engine's TestSafeETHSendsOnLocalChain.
func TestSwapsOnLocalChain(t *testing.T) {
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
	for _, k := range []string{"ENTRYPOINT_ADDRESS", "SAFE_PROXY_FACTORY_ADDRESS", "SAFE_SINGLETON_ADDRESS", "SAFE_MODULE_SETUP_ADDRESS", "SAFE_4337_MODULE_ADDRESS", "SAFE_MULTISEND_CALL_ONLY_ADDRESS"} {
		t.Setenv(k, stack[k])
	}
	t.Setenv("BASE_CHAIN_ID", stack["BASE_CHAIN_ID"])
	t.Setenv("OFFER_BOOK_ADDRESS", market["OFFER_BOOK_ADDRESS"])
	t.Setenv("OFFER_AUTHORIZER_PRIVATE_KEY", "0x5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a") // hardhat #2
	t.Setenv("OFFER_BOOK_CONFIRMATIONS", "0")

	ctx := context.Background()
	client, err := ethclient.Dial("http://127.0.0.1:8545")
	if err != nil {
		t.Fatal(err)
	}
	chainID, _ := client.ChainID(ctx)
	cfg := aa.ConfigFromEnv(chainID)
	b := &aa.Builder{Config: cfg, Chain: client, Bundler: aa.NewBundler(stack["BUNDLER_URL"]), GasBufferPercent: 300, DefaultValidity: 10 * time.Minute}
	book := common.HexToAddress(market["OFFER_BOOK_ADDRESS"])
	cngn := common.HexToAddress(market["MOCK_CNGN_ADDRESS"])             // 6 decimals
	ngni := common.HexToAddress(market["MOCK_INTERNAL_BALANCE_ADDRESS"]) // 18 decimals

	db, err := gorm.Open(sqlite.Open("file:"+t.TempDir()+"/book.db"), &gorm.Config{Logger: logger.Discard})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(offerbook.Models()...); err != nil {
		t.Fatal(err)
	}
	head, _ := client.BlockNumber(ctx)
	t.Setenv("OFFER_BOOK_START_BLOCK", new(big.Int).SetUint64(head+1).String())
	index := func() {
		t.Helper()
		if done, err := offerbook.Index(ctx, db, client, book, 50); err != nil || !done {
			t.Fatalf("indexing: done=%v err=%v", done, err)
		}
	}

	admin, _ := evmkeypair.ParseFull("0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
	erc20, _ := abi.JSON(strings.NewReader(`[{"name":"mint","type":"function","inputs":[{"type":"address"},{"type":"uint256"}],"outputs":[]},
		{"name":"balanceOf","type":"function","stateMutability":"view","inputs":[{"type":"address"}],"outputs":[{"type":"uint256"}]}]`))
	send := func(to common.Address, value *big.Int, data []byte) {
		t.Helper()
		if value.Sign() > 0 {
			h, err := sendValue(ctx, client, chainID, admin, to, value)
			if err != nil {
				t.Fatal(err)
			}
			if err := gnosissafe.WaitSuccess(ctx, client, h); err != nil {
				t.Fatal(err)
			}
			return
		}
		h, err := gnosissafe.SendTransaction(ctx, client, chainID, admin, to, data)
		if err != nil {
			t.Fatal(err)
		}
		if err := gnosissafe.WaitSuccess(ctx, client, common.HexToHash(h)); err != nil {
			t.Fatal(err)
		}
	}
	balance := func(token, holder common.Address) *big.Int {
		t.Helper()
		data, _ := erc20.Pack("balanceOf", holder)
		out, err := client.CallContract(ctx, ethereum.CallMsg{To: &token, Data: data}, nil)
		if err != nil {
			t.Fatal(err)
		}
		return new(big.Int).SetBytes(out)
	}
	mint := func(token, to common.Address, amount *big.Int) {
		data, _ := erc20.Pack("mint", to, amount)
		send(token, big.NewInt(0), data)
	}
	newWallet := func() (aa.Wallet, *evmkeypair.Full) {
		k, _ := evmkeypair.Random()
		owner := common.HexToAddress(k.Address())
		addr := cfg.SafeAddress([]common.Address{owner}, 1, big.NewInt(0))
		send(addr, big.NewInt(1e17), nil) // its gas
		return aa.Wallet{Address: addr, Owners: []common.Address{owner}, Threshold: 1, InitialOwners: []common.Address{owner}, InitialThreshold: 1, SaltNonce: big.NewInt(0)}, k
	}
	run := func(label string, w *aa.Wallet, calls []aa.Call, k *evmkeypair.Full) *aa.Receipt {
		t.Helper()
		p, err := b.Prepare(ctx, aa.Request{Wallet: *w, Calls: calls})
		if err != nil {
			t.Fatalf("%s: prepare: %v", label, err)
		}
		sig, _ := evmkeypair.SignPersonal(k.PrivateKey(), p.SafeOpHash.Bytes())
		h, err := b.Submit(ctx, p, []aa.OwnerSignature{{Owner: common.HexToAddress(k.Address()), Signature: sig}})
		if err != nil {
			t.Fatalf("%s: submit: %v", label, err)
		}
		wctx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()
		r, err := b.Bundler.(*aa.Bundler).WaitReceipt(wctx, h)
		if err != nil {
			t.Fatalf("%s: %v", label, err)
		}
		if !r.Success {
			t.Fatalf("%s reverted: %s", label, r.Reason)
		}
		w.InitialOwners, w.InitialThreshold = nil, 0 // deployed now
		return r
	}
	ngniUnits := func(n string) *big.Int { return decimal.RequireFromString(n).Shift(18).BigInt() }
	cngnUnits := func(n string) *big.Int { return decimal.RequireFromString(n).Shift(6).BigInt() }

	// --- a seller places two SELL offers of NGNI for cNGN
	seller, sellerKey := newWallet()
	mint(ngni, seller.Address, ngniUnits("1200"))
	cheap, _ := offerPrice(decimal.RequireFromString("1.4"), 18, 6)
	dear, _ := offerPrice(decimal.RequireFromString("1.5"), 18, 6)
	o1, _ := aa.CreateOfferCall(book, ngni, ngniUnits("200"), seller.Address, []common.Address{cngn}, []aa.OfferPrice{cheap})
	o2, _ := aa.CreateOfferCall(book, ngni, ngniUnits("1000"), seller.Address, []common.Address{cngn}, []aa.OfferPrice{dear})
	r := run("place offers", &seller, []aa.Call{aa.ERC20Approve(ngni, book, ngniUnits("1200")), o2, o1}, sellerKey)
	index()
	var offers []offerbook.Offer
	db.Order("id").Find(&offers)
	if len(offers) != 2 || !offers[0].Open || offers[0].Seller != offerbook.Key(seller.Address.Hex()) || offers[0].CreatedTx != r.Receipt.TransactionHash.Hex() {
		t.Fatalf("indexed offers: %+v", offers)
	}

	// --- a buyer swaps 500 cNGN: 200 NGNI at 1.4 (280 cNGN), then 220 cNGN
	// at 1.5 buys 146.666... NGNI
	buyer, buyerKey := newWallet()
	mint(cngn, buyer.Address, cngnUnits("1000"))
	route, err := routeSwap(ctx, db, client, book, cngn, ngni, cngnUnits("500"), "CNGN", "NGNI")
	if err != nil {
		t.Fatal(err)
	}
	if len(route.Fills) != 2 || route.Fills[0].Request.Amount.Cmp(ngniUnits("200")) != 0 || route.Fills[0].Payment.Cmp(cngnUnits("280")) != 0 {
		t.Fatalf("route: %+v", route.Fills)
	}
	if route.Pay.Cmp(cngnUnits("500")) > 0 || route.ReceiveHuman().Truncate(4).String() != "346.6666" {
		t.Fatalf("route pays %v, receives %v", route.PayHuman(), route.ReceiveHuman())
	}
	if err := route.authorize(buyer.Address, buyer.Address, time.Now().Add(10*time.Minute)); err != nil {
		t.Fatal(err)
	}
	calls, _ := route.calls()
	run("swap", &buyer, calls, buyerKey)
	if got := balance(ngni, buyer.Address); got.Cmp(route.Receive) != 0 {
		t.Fatalf("buyer received %v, want %v", got, route.Receive)
	}
	if got := balance(cngn, seller.Address); got.Cmp(route.Pay) != 0 {
		t.Fatalf("seller was paid %v, want %v", got, route.Pay)
	}

	// --- the index sees the fills
	index()
	var fills []offerbook.Fill
	db.Find(&fills)
	if len(fills) != 2 {
		t.Fatalf("fills: %+v", fills)
	}
	ob, err := offerbook.OrderBook(ctx, db, client, ngni, cngn, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(ob.Asks) != 1 || ob.Asks[0].Price.String() != "1.5" || ob.Asks[0].Amount.Truncate(4).String() != "853.3333" {
		t.Fatalf("asks: %+v", ob.Asks)
	}
	trades, err := offerbook.Trades(ctx, db, client, ngni, cngn, time.Time{}, time.Time{})
	if err != nil || len(trades) != 2 {
		t.Fatalf("trades %+v %v", trades, err)
	}
	candles := offerbook.Aggregate(trades, time.Hour, 0)
	if len(candles) != 1 || candles[0].Count != 2 || candles[0].Low.String() != "1.4" || candles[0].High.String() != "1.5" {
		t.Fatalf("candles: %+v", candles)
	}

	// --- a BUY offer (the buyer sells 100 cNGN for NGNI at 1.6) is a bid
	bidPrice, err := inverseOfferPrice(decimal.RequireFromString("1.6"), 18, 6)
	if err != nil {
		t.Fatal(err)
	}
	bid, _ := aa.CreateOfferCall(book, cngn, cngnUnits("100"), buyer.Address, []common.Address{ngni}, []aa.OfferPrice{bidPrice})
	run("buy offer", &buyer, []aa.Call{aa.ERC20Approve(cngn, book, cngnUnits("100")), bid}, buyerKey)
	index()
	ob, _ = offerbook.OrderBook(ctx, db, client, ngni, cngn, 10)
	if len(ob.Bids) != 1 || ob.Bids[0].Price.String() != "1.6" || ob.Bids[0].Amount.String() != "100" {
		t.Fatalf("bids: %+v", ob.Bids)
	}

	// --- more than the book holds is refused, saying what is on offer
	if _, err := routeSwap(ctx, db, client, book, cngn, ngni, cngnUnits("5000"), "CNGN", "NGNI"); err == nil || !strings.Contains(err.Error(), "low-liquidity") {
		t.Fatalf("an oversized swap: %v", err)
	}

	// --- cancelling the remaining offer takes it off the book
	dearID, _ := new(big.Int).SetString(offers[0].ID, 10)
	if offers[0].Remaining == "0" {
		dearID, _ = new(big.Int).SetString(offers[1].ID, 10)
	}
	run("cancel", &seller, []aa.Call{aa.CancelOfferCall(book, dearID)}, sellerKey)
	index()
	ob, _ = offerbook.OrderBook(ctx, db, client, ngni, cngn, 10)
	if len(ob.Asks) != 0 {
		t.Fatalf("asks after cancel: %+v", ob.Asks)
	}
	if got := balance(ngni, seller.Address); got.Sign() == 0 {
		t.Fatal("the cancelled offer's tokens did not return to the seller")
	}

	// --- an ETH batch send from a wallet, for payment-history-engine
	sent := run("eth batch", &buyer, []aa.Call{
		aa.NativeTransfer(common.HexToAddress("0x000000000000000000000000000000000000bEEF"), big.NewInt(1e15)),
		aa.NativeTransfer(common.HexToAddress("0x000000000000000000000000000000000000bEE2"), big.NewInt(2e15)),
	}, buyerKey)
	if f := os.Getenv("AA_SAFE_ETH_TX_FILE"); f != "" {
		os.WriteFile(f, []byte(sent.Receipt.TransactionHash.Hex()+" "+buyer.Address.Hex()), 0o600)
	}
}

func TestInverseOfferPrice(t *testing.T) {
	// buying the asset at 1.6 currency per token: a BUY offer sells the
	// currency at 1/1.6 = 0.625 asset per currency
	p, err := inverseOfferPrice(decimal.RequireFromString("1.6"), 18, 6)
	if err != nil {
		t.Fatal(err)
	}
	// 1 cNGN (1e6 base units) costs 0.625 NGNI (0.625e18 base units)
	if got := p.Cost(big.NewInt(1_000_000)); got.String() != "625000000000000000" {
		t.Fatalf("1 cNGN costs %v NGNI base units", got)
	}
	if _, err := inverseOfferPrice(decimal.Zero, 18, 6); err == nil {
		t.Fatal("zero price accepted")
	}
}
