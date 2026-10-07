package market

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"trovo-wallet-api/internal/offerbook"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	trov = "0xc0ffee0000000000000000000000000000000003"
	cngn = "0xc0ffee0000000000000000000000000000000001"
	usdc = "0xc0ffee0000000000000000000000000000000002"
)

func setup(t *testing.T) (*sharedconfig.GlobalConfig, time.Time) {
	t.Helper()
	t.Setenv("OFFER_BOOK_ADDRESS", "0x5FbDB2315678afecb367f032d93F642f64180aa3")
	t.Setenv("NAIRA_ASSET", "CNGN:"+cngn)
	t.Setenv("DOLLAR_ASSET", "USDC:"+usdc)
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(offerbook.Models()...); err != nil {
		t.Fatal(err)
	}
	if err := db.Exec(`CREATE TABLE curated_assets (asset_code text, asset_name text, contract_address text, image_url text, inactive integer default 0, priority integer)`).Error; err != nil {
		t.Fatal(err)
	}
	db.Exec(`INSERT INTO curated_assets VALUES ('TROV','Trovo Token',?,NULL,0,1),('CNGN','cNGN',?,NULL,0,2),('USDC','USD Coin',?,NULL,0,3),('OLD','Retired',?,NULL,1,4)`,
		trov, cngn, usdc, "0xc0ffee0000000000000000000000000000000009")
	db.Create(&[]offerbook.Token{{Address: trov, Decimals: 18}, {Address: cngn, Decimals: 6}, {Address: usdc, Decimals: 6}})
	now := time.Date(2026, 10, 7, 12, 0, 0, 0, time.UTC)
	// TROV/CNGN: 150 thirty hours ago, then 160 two hours ago
	db.Create(&[]offerbook.Fill{
		{ID: "a:0", OfferID: "1", SellToken: trov, PaymentToken: cngn, Amount: "1000000000000000000", Payment: "150000000", FilledAt: now.Add(-30 * time.Hour)},
		{ID: "b:0", OfferID: "1", SellToken: trov, PaymentToken: cngn, Amount: "1000000000000000000", Payment: "160000000", FilledAt: now.Add(-2 * time.Hour)},
	})
	// an ask: 2 TROV at 170 CNGN; a bid: 340 CNGN buying TROV at 165
	db.Create(&offerbook.Offer{ID: "1", SellToken: trov, Remaining: "2000000000000000000", Open: true})
	db.Create(&offerbook.Price{OfferID: "1", PaymentToken: cngn, Num: "170", Den: "1000000000000"})
	db.Create(&offerbook.Offer{ID: "2", SellToken: cngn, Remaining: "340000000", Open: true})
	// bid price: CNGN units per TROV unit inverse: 1 CNGN buys 1/165 TROV
	db.Create(&offerbook.Price{OfferID: "2", PaymentToken: trov, Num: "1000000000000", Den: "165"})

	tradable = func(_ context.Context, _ *sharedconfig.GlobalConfig, token common.Address) (bool, error) {
		return token != common.HexToAddress(usdc), nil // USDC is not listed for trading
	}
	t.Cleanup(func() {
		tradableCache.Lock()
		tradableCache.at, tradableCache.val = map[string]time.Time{}, map[string]bool{}
		tradableCache.Unlock()
	})
	return &sharedconfig.GlobalConfig{DB: db}, now
}

func TestPairs(t *testing.T) {
	gc, now := setup(t)
	pairs, err := Pairs(context.Background(), gc, now)
	if err != nil {
		t.Fatal(err)
	}
	// TROV/CNGN only: USDC is not tradable, OLD is inactive, CNGN/CNGN is skipped
	if len(pairs) != 1 {
		t.Fatalf("want 1 pair, got %+v", pairs)
	}
	p := pairs[0]
	if p.Base.Code != "TROV" || p.Counter.Code != "CNGN" || p.Counter.Name != "cNGN" {
		t.Fatalf("pair %+v", p)
	}
	if p.LastPrice != "160" || p.ChangePercent24h != "6.67" || p.Trades24h != 1 || p.High24h != "160" || p.BaseVolume24h != "1" || p.CounterVolume24h != "160" {
		t.Fatalf("stats %+v", p)
	}
	if p.BestAsk != "170" || p.BestBid != "165" {
		t.Fatalf("best ask/bid %q/%q", p.BestAsk, p.BestBid)
	}
}

func TestOrderBookTradesCandles(t *testing.T) {
	gc, now := setup(t)
	ctx := context.Background()
	book, err := OrderBook(ctx, gc, trov, cngn, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(book.Asks) != 1 || book.Asks[0].Price != "170" || book.Asks[0].Amount != "2" || len(book.Bids) != 1 || book.Bids[0].Price != "165" || book.Bids[0].Amount != "340" {
		t.Fatalf("book %+v", book)
	}
	trades, err := RecentTrades(ctx, gc, trov, cngn, 10, now)
	if err != nil {
		t.Fatal(err)
	}
	if len(trades) != 2 || trades[0].Price != "160" || trades[1].Price != "150" {
		t.Fatalf("trades (newest first) %+v", trades)
	}
	candles, err := Candles(ctx, gc, trov, cngn, Resolutions["1d"], 7, now)
	if err != nil {
		t.Fatal(err)
	}
	if len(candles) != 2 || candles[0].Close != "150" || candles[1].Open != "160" || candles[0].Timestamp >= candles[1].Timestamp {
		t.Fatalf("candles %+v", candles)
	}
	if _, err := OrderBook(ctx, gc, trov, trov, 10); err != ErrBadPair {
		t.Fatalf("same token twice: %v", err)
	}
}

func TestRoutes(t *testing.T) {
	gc, fixed := setup(t)
	now = func() time.Time { return fixed }
	t.Cleanup(func() { now = time.Now })
	gin.SetMode(gin.TestMode)
	r := gin.New()
	Init(r, gc)
	get := func(path string) (int, map[string]interface{}) {
		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, path, nil))
		var body map[string]interface{}
		_ = json.Unmarshal(w.Body.Bytes(), &body)
		return w.Code, body
	}
	if code, body := get("/v1/market/pairs"); code != 200 || len(body["pairs"].([]interface{})) != 1 {
		t.Fatalf("pairs %d %v", code, body)
	}
	if code, _ := get("/v1/market/candles?base=" + trov + "&counter=" + cngn + "&resolution=2h"); code != 400 {
		t.Fatalf("bad resolution: %d", code)
	}
	if code, _ := get("/v1/market/orderbook?base=nope&counter=" + cngn); code != 400 {
		t.Fatalf("bad pair: %d", code)
	}
	if code, body := get("/v1/market/trades?base=" + trov + "&counter=" + cngn); code != 200 || len(body["trades"].([]interface{})) != 2 {
		t.Fatalf("trades %d %v", code, body)
	}
	t.Setenv("OFFER_BOOK_ADDRESS", "")
	if code, _ := get("/v1/market/pairs"); code != 503 {
		t.Fatalf("not configured: %d", code)
	}
}
