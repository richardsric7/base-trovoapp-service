// Package market serves the token market (TrovoOfferBook) to the apps over
// plain REST: the pairs with their last 24 hours, each pair's order book,
// recent trades and price candles. Everything is read from the offer book
// index (internal/offerbook); placing and cancelling orders stays with
// /v1/users/trades.
package market

import (
	"context"
	"errors"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"trovo-wallet-api/internal/aa"
	"trovo-wallet-api/internal/offerbook"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
)

// Token is one side of a pair.
type Token struct {
	Code            string `json:"code" example:"TROV"`
	Name            string `json:"name" example:"Trovo Token"`
	ContractAddress string `json:"contractAddress" example:"0xC0FFEE0000000000000000000000000000000003"`
	ImageURL        string `json:"imageUrl"`
}

// Pair is a market with its last 24 hours. Prices are in counter per base.
type Pair struct {
	Base             Token  `json:"base"`
	Counter          Token  `json:"counter"`
	LastPrice        string `json:"lastPrice" example:"152.5"`
	ChangePercent24h string `json:"changePercent24h" example:"-1.25"`
	High24h          string `json:"high24h"`
	Low24h           string `json:"low24h"`
	BaseVolume24h    string `json:"baseVolume24h"`
	CounterVolume24h string `json:"counterVolume24h"`
	Trades24h        int64  `json:"trades24h"`
	BestAsk          string `json:"bestAsk"`
	BestBid          string `json:"bestBid"`
	LastTradeAt      *int64 `json:"lastTradeAt,omitempty"`
}

// Level is one order book price level.
type Level struct {
	Price  string `json:"price"`
	Amount string `json:"amount"`
}

// Book is a pair's order book: asks sell base (amount in base, cheapest
// first), bids buy base (amount in counter, best first).
type Book struct {
	Asks []Level `json:"asks"`
	Bids []Level `json:"bids"`
}

// Trade is one fill, newest first.
type Trade struct {
	Timestamp     int64  `json:"timestamp"`
	Price         string `json:"price"`
	BaseAmount    string `json:"baseAmount"`
	CounterAmount string `json:"counterAmount"`
}

// Candle is one price bucket.
type Candle struct {
	Timestamp     int64  `json:"timestamp"`
	Open          string `json:"open"`
	High          string `json:"high"`
	Low           string `json:"low"`
	Close         string `json:"close"`
	BaseVolume    string `json:"baseVolume"`
	CounterVolume string `json:"counterVolume"`
	Trades        int64  `json:"trades"`
}

// ErrNotConfigured: no offer book (OFFER_BOOK_ADDRESS).
var ErrNotConfigured = errors.New("the market is not available")

// ErrBadPair: base or counter is not a token address, or they are equal.
var ErrBadPair = errors.New("base and counter must be two different token contract addresses")

// tradable reports whether the offer book lists token for trading. A
// variable so tests need no chain.
var tradable = func(ctx context.Context, gc *sharedconfig.GlobalConfig, token common.Address) (bool, error) {
	book, err := offerbook.BookAddress()
	if err != nil {
		return false, ErrNotConfigured
	}
	return aa.Tradable(ctx, gc.BantuExpansionClient, book, token)
}

var tradableCache = struct {
	sync.Mutex
	at  map[string]time.Time
	val map[string]bool
}{at: map[string]time.Time{}, val: map[string]bool{}}

const tradableTTL = 10 * time.Minute

func isTradable(ctx context.Context, gc *sharedconfig.GlobalConfig, token common.Address) (bool, error) {
	k := strings.ToLower(token.Hex())
	tradableCache.Lock()
	if at, ok := tradableCache.at[k]; ok && time.Since(at) < tradableTTL {
		v := tradableCache.val[k]
		tradableCache.Unlock()
		return v, nil
	}
	tradableCache.Unlock()
	v, err := tradable(ctx, gc, token)
	if err != nil {
		return false, err
	}
	tradableCache.Lock()
	tradableCache.at[k], tradableCache.val[k] = time.Now(), v
	tradableCache.Unlock()
	return v, nil
}

// quoteTokens are the currencies pairs are priced in: NAIRA_ASSET and
// DOLLAR_ASSET ("CODE:0xcontract").
func quoteTokens() []Token {
	var out []Token
	for _, env := range []string{"NAIRA_ASSET", "DOLLAR_ASSET"} {
		parts := strings.SplitN(strings.TrimSpace(os.Getenv(env)), ":", 2)
		if len(parts) == 2 && common.IsHexAddress(parts[1]) {
			out = append(out, Token{Code: parts[0], ContractAddress: common.HexToAddress(parts[1]).Hex()})
		}
	}
	return out
}

// listedTokens are the active curated assets that are tokens.
func listedTokens(gc *sharedconfig.GlobalConfig) ([]Token, error) {
	var rows []struct {
		AssetCode       string
		AssetName       string
		ContractAddress string
		ImageURL        *string
	}
	if err := gc.DB.Table("curated_assets").
		Select("asset_code, asset_name, contract_address, image_url").
		Where("inactive = 0 AND contract_address <> ''").
		Order("priority, asset_code").Scan(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]Token, 0, len(rows))
	for _, r := range rows {
		if !common.IsHexAddress(r.ContractAddress) {
			continue
		}
		t := Token{Code: r.AssetCode, Name: r.AssetName, ContractAddress: common.HexToAddress(r.ContractAddress).Hex()}
		if r.ImageURL != nil {
			t.ImageURL = *r.ImageURL
		}
		out = append(out, t)
	}
	return out, nil
}

// Pairs lists every tradable curated token against each quote currency,
// with its last 24 hours, busiest first.
func Pairs(ctx context.Context, gc *sharedconfig.GlobalConfig, now time.Time) ([]Pair, error) {
	if _, err := offerbook.BookAddress(); err != nil {
		return nil, ErrNotConfigured
	}
	tokens, err := listedTokens(gc)
	if err != nil {
		return nil, err
	}
	byAddr := map[string]Token{}
	for _, t := range tokens {
		byAddr[strings.ToLower(t.ContractAddress)] = t
	}
	pairs := []Pair{}
	for _, q := range quoteTokens() {
		if known, ok := byAddr[strings.ToLower(q.ContractAddress)]; ok {
			q = known
		}
		if ok, err := isTradable(ctx, gc, common.HexToAddress(q.ContractAddress)); err != nil || !ok {
			continue
		}
		for _, b := range tokens {
			if strings.EqualFold(b.ContractAddress, q.ContractAddress) {
				continue
			}
			if ok, err := isTradable(ctx, gc, common.HexToAddress(b.ContractAddress)); err != nil || !ok {
				continue
			}
			p, err := pairStats(ctx, gc, b, q, now)
			if err != nil {
				return nil, err
			}
			pairs = append(pairs, p)
		}
	}
	sort.SliceStable(pairs, func(i, j int) bool { return pairs[i].Trades24h > pairs[j].Trades24h })
	return pairs, nil
}

func pairStats(ctx context.Context, gc *sharedconfig.GlobalConfig, base, counter Token, now time.Time) (Pair, error) {
	p := Pair{Base: base, Counter: counter, LastPrice: "0", ChangePercent24h: "0", High24h: "0", Low24h: "0", BaseVolume24h: "0", CounterVolume24h: "0"}
	b, c := common.HexToAddress(base.ContractAddress), common.HexToAddress(counter.ContractAddress)
	// a year of trades: enough for a last price on a quiet market
	trades, err := offerbook.Trades(ctx, gc.DB, gc.BantuExpansionClient, b, c, now.AddDate(-1, 0, 0), time.Time{})
	if err != nil {
		return p, err
	}
	since := now.Add(-24 * time.Hour)
	var open, high, low decimal.Decimal
	var baseVol, counterVol decimal.Decimal
	var count int64
	for _, t := range trades {
		if t.At.Before(since) {
			open = t.Price // the last price before the window
			continue
		}
		if count == 0 {
			if open.IsZero() {
				open = t.Price
			}
			high, low = t.Price, t.Price
		}
		count++
		if t.Price.GreaterThan(high) {
			high = t.Price
		}
		if t.Price.LessThan(low) {
			low = t.Price
		}
		baseVol, counterVol = baseVol.Add(t.BaseAmount), counterVol.Add(t.CounterAmount)
	}
	if n := len(trades); n > 0 {
		last := trades[n-1]
		p.LastPrice = last.Price.String()
		ts := last.At.Unix()
		p.LastTradeAt = &ts
		if count > 0 && open.IsPositive() {
			p.ChangePercent24h = last.Price.Sub(open).Div(open).Mul(decimal.NewFromInt(100)).StringFixed(2)
		}
	}
	if count > 0 {
		p.High24h, p.Low24h = high.String(), low.String()
		p.BaseVolume24h, p.CounterVolume24h = baseVol.String(), counterVol.String()
	}
	p.Trades24h = count
	if book, err := offerbook.OrderBook(ctx, gc.DB, gc.BantuExpansionClient, b, c, 1); err == nil {
		if len(book.Asks) > 0 {
			p.BestAsk = book.Asks[0].Price.String()
		}
		if len(book.Bids) > 0 {
			p.BestBid = book.Bids[0].Price.String()
		}
	}
	return p, nil
}

func pairAddresses(base, counter string) (common.Address, common.Address, error) {
	if !common.IsHexAddress(base) || !common.IsHexAddress(counter) || strings.EqualFold(base, counter) {
		return common.Address{}, common.Address{}, ErrBadPair
	}
	return common.HexToAddress(base), common.HexToAddress(counter), nil
}

// OrderBook is a pair's order book, limit levels per side.
func OrderBook(ctx context.Context, gc *sharedconfig.GlobalConfig, base, counter string, limit int) (Book, error) {
	b, c, err := pairAddresses(base, counter)
	if err != nil {
		return Book{}, err
	}
	book, err := offerbook.OrderBook(ctx, gc.DB, gc.BantuExpansionClient, b, c, limit)
	if err != nil {
		return Book{}, err
	}
	out := Book{Asks: []Level{}, Bids: []Level{}}
	for _, l := range book.Asks {
		out.Asks = append(out.Asks, Level{Price: l.Price.String(), Amount: l.Amount.String()})
	}
	for _, l := range book.Bids {
		out.Bids = append(out.Bids, Level{Price: l.Price.String(), Amount: l.Amount.String()})
	}
	return out, nil
}

// RecentTrades are a pair's last limit trades from the last 30 days,
// newest first.
func RecentTrades(ctx context.Context, gc *sharedconfig.GlobalConfig, base, counter string, limit int, now time.Time) ([]Trade, error) {
	b, c, err := pairAddresses(base, counter)
	if err != nil {
		return nil, err
	}
	trades, err := offerbook.Trades(ctx, gc.DB, gc.BantuExpansionClient, b, c, now.AddDate(0, 0, -30), time.Time{})
	if err != nil {
		return nil, err
	}
	out := []Trade{}
	for i := len(trades) - 1; i >= 0 && len(out) < limit; i-- {
		t := trades[i]
		out = append(out, Trade{Timestamp: t.At.Unix(), Price: t.Price.String(), BaseAmount: t.BaseAmount.String(), CounterAmount: t.CounterAmount.String()})
	}
	return out, nil
}

// Resolutions are the candle sizes the apps can ask for.
var Resolutions = map[string]time.Duration{
	"15m": 15 * time.Minute,
	"1h":  time.Hour,
	"4h":  4 * time.Hour,
	"1d":  24 * time.Hour,
	"1w":  7 * 24 * time.Hour,
}

// Candles are a pair's last limit candles of resolution, oldest first.
// Periods without trades are left out.
func Candles(ctx context.Context, gc *sharedconfig.GlobalConfig, base, counter string, resolution time.Duration, limit int, now time.Time) ([]Candle, error) {
	b, c, err := pairAddresses(base, counter)
	if err != nil {
		return nil, err
	}
	start := now.Add(-time.Duration(limit) * resolution)
	trades, err := offerbook.Trades(ctx, gc.DB, gc.BantuExpansionClient, b, c, start, time.Time{})
	if err != nil {
		return nil, err
	}
	out := []Candle{}
	for _, k := range offerbook.Aggregate(trades, resolution, 0) {
		out = append(out, Candle{
			Timestamp: k.Start.Unix(), Open: k.Open.String(), High: k.High.String(), Low: k.Low.String(), Close: k.Close.String(),
			BaseVolume: k.BaseVolume.String(), CounterVolume: k.CounterVolume.String(), Trades: k.Count,
		})
	}
	return out, nil
}
