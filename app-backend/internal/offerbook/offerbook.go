// Package offerbook keeps a database index of TrovoOfferBook
// (market/contracts): every offer, its prices and every fill, built from
// the contract's events. Swaps route over it, and the order book, prices
// and trade charts are read from it. The chain stays the authority: swaps
// re-read the offers they use on-chain before pricing a trade.
package offerbook

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"trovo-wallet-api/internal/aa"
	"trovo-wallet-api/internal/basetxn"
	"trovo-wallet-api/internal/network"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Offer is an offer on the book. Addresses are lowercase hex; amounts are
// base units in decimal.
type Offer struct {
	ID                string `gorm:"primaryKey;size:80"`
	Seller            string `gorm:"size:42;index"`
	SellToken         string `gorm:"size:42;index"`
	ProceedsRecipient string `gorm:"size:42"`
	Remaining         string `gorm:"size:80"`
	Open              bool   `gorm:"index"`
	CreatedBlock      uint64
	CreatedTx         string `gorm:"size:66;index"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// TableName keeps the index's tables together.
func (Offer) TableName() string { return "offer_book_offers" }

// Price is an offer's price in one payment token: buying n base units
// costs ceil(n*Num/Den) base units of the payment token.
type Price struct {
	OfferID      string `gorm:"primaryKey;size:80"`
	PaymentToken string `gorm:"primaryKey;size:42;index"`
	Num          string `gorm:"size:80"`
	Den          string `gorm:"size:80"`
}

func (Price) TableName() string { return "offer_book_prices" }

// Fill is one purchase from an offer.
type Fill struct {
	ID           string `gorm:"primaryKey;size:80"` // tx hash:log index
	OfferID      string `gorm:"size:80;index"`
	Taker        string `gorm:"size:42;index"`
	Recipient    string `gorm:"size:42"`
	SellToken    string `gorm:"size:42;index:idx_offer_book_fill_pair"`
	PaymentToken string `gorm:"size:42;index:idx_offer_book_fill_pair"`
	Amount       string `gorm:"size:80"`
	Payment      string `gorm:"size:80"`
	Block        uint64
	TxHash       string    `gorm:"size:66;index"`
	FilledAt     time.Time `gorm:"index"`
}

func (Fill) TableName() string { return "offer_book_fills" }

// Cursor is where indexing a book continues: the next block to index.
type Cursor struct {
	ID    string `gorm:"primaryKey;size:80"`
	Block uint64
}

func (Cursor) TableName() string { return "offer_book_cursors" }

// Token is a token's decimals, read once.
type Token struct {
	Address  string `gorm:"primaryKey;size:42"`
	Decimals uint8
}

func (Token) TableName() string { return "offer_book_tokens" }

// Models are the index's tables, for migration.
func Models() []interface{} {
	return []interface{}{&Offer{}, &Price{}, &Fill{}, &Cursor{}, &Token{}}
}

// ErrNotConfigured: OFFER_BOOK_ADDRESS is not set.
var ErrNotConfigured = errors.New("offer book not configured (OFFER_BOOK_ADDRESS)")

// BookAddress is OFFER_BOOK_ADDRESS.
func BookAddress() (common.Address, error) {
	raw := strings.TrimSpace(os.Getenv("OFFER_BOOK_ADDRESS"))
	if !common.IsHexAddress(raw) {
		return common.Address{}, ErrNotConfigured
	}
	return common.HexToAddress(raw), nil
}

func key(a common.Address) string { return strings.ToLower(a.Hex()) }

// Key is the index's form of an address (lowercase hex), or "" when s is
// not an address.
func Key(s string) string {
	s = strings.TrimSpace(s)
	if !common.IsHexAddress(s) {
		return ""
	}
	return key(common.HexToAddress(s))
}

// the index the order book helpers read (set by Use)
var (
	useMu  sync.RWMutex
	useDB  *gorm.DB
	useEth *ethclient.Client
)

// Use sets the database and client the package-level readers use.
func Use(db *gorm.DB, client *ethclient.Client) {
	useMu.Lock()
	defer useMu.Unlock()
	useDB, useEth = db, client
}

func current() (*gorm.DB, *ethclient.Client) {
	useMu.RLock()
	defer useMu.RUnlock()
	return useDB, useEth
}

// ---------------------------------------------------------------------------
// Indexing

func confirmations() uint64 {
	if n, err := strconv.ParseUint(strings.TrimSpace(os.Getenv("OFFER_BOOK_CONFIRMATIONS")), 10, 64); err == nil {
		return n
	}
	return 2
}

func startBlock() uint64 {
	n, _ := strconv.ParseUint(strings.TrimSpace(os.Getenv("OFFER_BOOK_START_BLOCK")), 10, 64)
	return n
}

// chunkBlocks bounds one eth_getLogs range (providers cap it).
const chunkBlocks = 5000

// Index brings the index up to the chain (less OFFER_BOOK_CONFIRMATIONS
// blocks), at most maxChunks log queries per call; it reports whether it
// caught up. Each chunk's events and the cursor are saved together.
func Index(ctx context.Context, db *gorm.DB, client *ethclient.Client, book common.Address, maxChunks int) (bool, error) {
	head, err := client.BlockNumber(ctx)
	if err != nil {
		return false, err
	}
	conf := confirmations()
	if head < conf {
		return true, nil
	}
	safe := head - conf
	cur := Cursor{ID: key(book)} // Block: the next block to index
	if db.First(&cur, "id = ?", cur.ID).Error != nil {
		cur.Block = startBlock()
		if err := db.Create(&cur).Error; err != nil {
			return false, err
		}
	}
	for i := 0; i < maxChunks; i++ {
		from := cur.Block
		if from > safe {
			return true, nil
		}
		to := safe
		if to-from >= chunkBlocks {
			to = from + chunkBlocks - 1
		}
		logs, err := client.FilterLogs(ctx, ethereum.FilterQuery{
			FromBlock: new(big.Int).SetUint64(from), ToBlock: new(big.Int).SetUint64(to),
			Addresses: []common.Address{book}, Topics: [][]common.Hash{aa.OfferBookTopics()},
		})
		if err != nil {
			return false, err
		}
		times := map[uint64]time.Time{}
		for _, l := range logs {
			if _, ok := times[l.BlockNumber]; ok {
				continue
			}
			h, err := client.HeaderByNumber(ctx, new(big.Int).SetUint64(l.BlockNumber))
			if err != nil {
				return false, err
			}
			times[l.BlockNumber] = time.Unix(int64(h.Time), 0).UTC()
		}
		for _, l := range logs {
			if ev, ok := aa.DecodeBookEvent(l); ok {
				if err := ensureTokens(ctx, db, client, ev); err != nil {
					return false, err
				}
			}
		}
		err = db.Transaction(func(tx *gorm.DB) error {
			for _, l := range logs {
				ev, ok := aa.DecodeBookEvent(l)
				if !ok || l.Removed {
					continue
				}
				if err := apply(tx, ev, times[l.BlockNumber]); err != nil {
					return err
				}
			}
			return tx.Model(&Cursor{}).Where("id = ?", cur.ID).Update("block", to+1).Error
		})
		if err != nil {
			return false, err
		}
		cur.Block = to + 1
	}
	return cur.Block > safe, nil
}

func ensureTokens(ctx context.Context, db *gorm.DB, client *ethclient.Client, ev aa.BookEvent) error {
	for _, t := range []common.Address{ev.SellToken, ev.PaymentToken} {
		if t == (common.Address{}) {
			continue
		}
		var n int64
		db.Model(&Token{}).Where("address = ?", key(t)).Count(&n)
		if n > 0 {
			continue
		}
		d, err := network.AssetDecimals(ctx, client, basetxn.CreditAsset{Code: "TOKEN", Issuer: t.Hex()})
		if err != nil {
			return err
		}
		if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&Token{Address: key(t), Decimals: d}).Error; err != nil {
			return err
		}
	}
	return nil
}

func apply(tx *gorm.DB, ev aa.BookEvent, at time.Time) error {
	id := ev.OfferID.String()
	switch ev.Kind {
	case "created":
		o := Offer{ID: id, Seller: key(ev.Seller), SellToken: key(ev.SellToken), ProceedsRecipient: key(ev.ProceedsRecipient),
			Remaining: ev.Amount.String(), Open: true, CreatedBlock: ev.Log.BlockNumber, CreatedTx: ev.Log.TxHash.Hex()}
		return tx.Clauses(clause.OnConflict{UpdateAll: true}).Create(&o).Error
	case "price":
		p := Price{OfferID: id, PaymentToken: key(ev.PaymentToken), Num: ev.Price.Num.String(), Den: ev.Price.Den.String()}
		if ev.Price.Num.Sign() == 0 {
			return tx.Where("offer_id = ? AND payment_token = ?", p.OfferID, p.PaymentToken).Delete(&Price{}).Error
		}
		return tx.Clauses(clause.OnConflict{UpdateAll: true}).Create(&p).Error
	case "cancelled":
		return tx.Model(&Offer{}).Where("id = ?", id).Updates(map[string]interface{}{"open": false, "remaining": "0"}).Error
	case "filled":
		var o Offer
		if err := tx.First(&o, "id = ?", id).Error; err != nil {
			return fmt.Errorf("fill of unknown offer %v: %w", id, err)
		}
		rem, _ := new(big.Int).SetString(o.Remaining, 10)
		if rem == nil {
			rem = new(big.Int)
		}
		rem.Sub(rem, ev.Amount)
		if rem.Sign() < 0 {
			rem.SetInt64(0)
		}
		if err := tx.Model(&Offer{}).Where("id = ?", id).Update("remaining", rem.String()).Error; err != nil {
			return err
		}
		f := Fill{
			ID: fmt.Sprintf("%s:%d", ev.Log.TxHash.Hex(), ev.Log.Index), OfferID: id, Taker: key(ev.Taker), Recipient: key(ev.Recipient),
			SellToken: o.SellToken, PaymentToken: key(ev.PaymentToken), Amount: ev.Amount.String(), Payment: ev.Payment.String(),
			Block: ev.Log.BlockNumber, TxHash: ev.Log.TxHash.Hex(), FilledAt: at,
		}
		return tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&f).Error
	}
	return nil
}

// ---------------------------------------------------------------------------
// Reading

// Candidate is an open offer selling one token for another, at its price.
type Candidate struct {
	OfferID   *big.Int
	Seller    common.Address
	Remaining *big.Int
	Price     aa.OfferPrice
}

// cheaper reports whether price a costs less per unit than b.
func cheaper(a, b aa.OfferPrice) int {
	return new(big.Int).Mul(a.Num, b.Den).Cmp(new(big.Int).Mul(b.Num, a.Den))
}

// Candidates are the open offers selling sellToken that accept
// paymentToken, cheapest first (then oldest).
func Candidates(db *gorm.DB, sellToken, paymentToken common.Address) ([]Candidate, error) {
	type row struct {
		ID        string
		Seller    string
		Remaining string
		Num       string
		Den       string
	}
	var rows []row
	err := db.Table("offer_book_offers o").
		Select("o.id, o.seller, o.remaining, p.num, p.den").
		Joins("JOIN offer_book_prices p ON p.offer_id = o.id").
		Where("o.open = ? AND o.sell_token = ? AND p.payment_token = ?", true, key(sellToken), key(paymentToken)).
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	var out []Candidate
	for _, r := range rows {
		id, _ := new(big.Int).SetString(r.ID, 10)
		rem, _ := new(big.Int).SetString(r.Remaining, 10)
		num, _ := new(big.Int).SetString(r.Num, 10)
		den, _ := new(big.Int).SetString(r.Den, 10)
		if id == nil || rem == nil || rem.Sign() <= 0 || num == nil || den == nil || den.Sign() == 0 || num.Sign() == 0 {
			continue
		}
		out = append(out, Candidate{OfferID: id, Seller: common.HexToAddress(r.Seller), Remaining: rem, Price: aa.OfferPrice{Num: num, Den: den}})
	}
	sort.SliceStable(out, func(i, j int) bool {
		if c := cheaper(out[i].Price, out[j].Price); c != 0 {
			return c < 0
		}
		return out[i].OfferID.Cmp(out[j].OfferID) < 0
	})
	return out, nil
}

// Decimals reads a token's decimals from the index, else the chain.
func Decimals(ctx context.Context, db *gorm.DB, client *ethclient.Client, token common.Address) (uint8, error) {
	var t Token
	if db.First(&t, "address = ?", key(token)).Error == nil {
		return t.Decimals, nil
	}
	if client == nil {
		return 0, fmt.Errorf("decimals of %v unknown", token.Hex())
	}
	d, err := network.AssetDecimals(ctx, client, basetxn.CreditAsset{Code: "TOKEN", Issuer: token.Hex()})
	if err != nil {
		return 0, err
	}
	db.Clauses(clause.OnConflict{DoNothing: true}).Create(&Token{Address: key(token), Decimals: d})
	return d, nil
}

// HumanPrice is p as an amount of the payment token per whole sell token.
func HumanPrice(p aa.OfferPrice, sellDec, payDec uint8) decimal.Decimal {
	num := decimal.NewFromBigInt(p.Num, 0)
	den := decimal.NewFromBigInt(p.Den, 0)
	return num.Div(den).Shift(int32(sellDec) - int32(payDec))
}

// Level is one price level of an order book side.
type Level struct {
	Price  decimal.Decimal
	Amount decimal.Decimal
}

// Book is the order book of base against counter, in the form Stellar's
// order book had: asks sell base (amount in base); bids buy base with
// counter (amount in counter); both priced in counter per base.
type Book struct {
	Asks []Level
	Bids []Level
}

// OrderBook reads the order book of base against counter.
func OrderBook(ctx context.Context, db *gorm.DB, client *ethclient.Client, base, counter common.Address, limit int) (Book, error) {
	var b Book
	baseDec, err := Decimals(ctx, db, client, base)
	if err != nil {
		return b, err
	}
	counterDec, err := Decimals(ctx, db, client, counter)
	if err != nil {
		return b, err
	}
	asks, err := Candidates(db, base, counter)
	if err != nil {
		return b, err
	}
	for _, c := range asks {
		b.Asks = addLevel(b.Asks, HumanPrice(c.Price, baseDec, counterDec), decimal.NewFromBigInt(c.Remaining, -int32(baseDec)))
	}
	bids, err := Candidates(db, counter, base)
	if err != nil {
		return b, err
	}
	for _, c := range bids {
		perCounter := HumanPrice(c.Price, counterDec, baseDec) // base per counter
		if !perCounter.IsPositive() {
			continue
		}
		price := decimal.NewFromInt(1).DivRound(perCounter, 18)
		b.Bids = addLevel(b.Bids, price, decimal.NewFromBigInt(c.Remaining, -int32(counterDec)))
	}
	// asks cheapest first; bids best (highest) first
	sort.SliceStable(b.Bids, func(i, j int) bool { return b.Bids[i].Price.GreaterThan(b.Bids[j].Price) })
	if limit > 0 {
		if len(b.Asks) > limit {
			b.Asks = b.Asks[:limit]
		}
		if len(b.Bids) > limit {
			b.Bids = b.Bids[:limit]
		}
	}
	return b, nil
}

func addLevel(levels []Level, price, amount decimal.Decimal) []Level {
	if n := len(levels); n > 0 && levels[n-1].Price.Equal(price) {
		levels[n-1].Amount = levels[n-1].Amount.Add(amount)
		return levels
	}
	return append(levels, Level{Price: price, Amount: amount})
}

// Trade is a fill seen from the base/counter pair.
type Trade struct {
	At            time.Time
	BaseAmount    decimal.Decimal
	CounterAmount decimal.Decimal
	Price         decimal.Decimal // counter per base
}

// Trades are the fills between base and counter (either direction) from
// start (zero: any) to end (zero: now), oldest first.
func Trades(ctx context.Context, db *gorm.DB, client *ethclient.Client, base, counter common.Address, start, end time.Time) ([]Trade, error) {
	baseDec, err := Decimals(ctx, db, client, base)
	if err != nil {
		return nil, err
	}
	counterDec, err := Decimals(ctx, db, client, counter)
	if err != nil {
		return nil, err
	}
	q := db.Where("(sell_token = ? AND payment_token = ?) OR (sell_token = ? AND payment_token = ?)", key(base), key(counter), key(counter), key(base))
	if !start.IsZero() {
		q = q.Where("filled_at >= ?", start)
	}
	if !end.IsZero() {
		q = q.Where("filled_at < ?", end)
	}
	var fills []Fill
	if err := q.Order("filled_at, id").Find(&fills).Error; err != nil {
		return nil, err
	}
	out := make([]Trade, 0, len(fills))
	for _, f := range fills {
		amount, err1 := decimal.NewFromString(f.Amount)
		payment, err2 := decimal.NewFromString(f.Payment)
		if err1 != nil || err2 != nil {
			continue
		}
		t := Trade{At: f.FilledAt}
		if f.SellToken == key(base) {
			t.BaseAmount, t.CounterAmount = amount.Shift(-int32(baseDec)), payment.Shift(-int32(counterDec))
		} else {
			t.BaseAmount, t.CounterAmount = payment.Shift(-int32(baseDec)), amount.Shift(-int32(counterDec))
		}
		if !t.BaseAmount.IsPositive() {
			continue
		}
		t.Price = t.CounterAmount.DivRound(t.BaseAmount, 18)
		out = append(out, t)
	}
	return out, nil
}

// Candle is one trade-aggregation bucket.
type Candle struct {
	Start         time.Time
	Count         int64
	BaseVolume    decimal.Decimal
	CounterVolume decimal.Decimal
	Open, High    decimal.Decimal
	Low, Close    decimal.Decimal
	Average       decimal.Decimal // volume-weighted
}

// Aggregate buckets trades into candles of resolution, aligned to offset
// past the epoch; empty buckets are left out.
func Aggregate(trades []Trade, resolution, offset time.Duration) []Candle {
	if resolution <= 0 {
		return nil
	}
	var out []Candle
	for _, t := range trades {
		start := time.Unix(0, (t.At.UnixNano()-int64(offset))/int64(resolution)*int64(resolution)+int64(offset)).UTC()
		n := len(out)
		if n == 0 || !out[n-1].Start.Equal(start) {
			out = append(out, Candle{Start: start, Open: t.Price, High: t.Price, Low: t.Price})
			n++
		}
		c := &out[n-1]
		c.Count++
		c.BaseVolume = c.BaseVolume.Add(t.BaseAmount)
		c.CounterVolume = c.CounterVolume.Add(t.CounterAmount)
		if t.Price.GreaterThan(c.High) {
			c.High = t.Price
		}
		if t.Price.LessThan(c.Low) {
			c.Low = t.Price
		}
		c.Close = t.Price
		c.Average = c.CounterVolume.DivRound(c.BaseVolume, 18)
	}
	return out
}

// ---------------------------------------------------------------------------
// Package-level readers for the order book helpers (see Use)

// ErrUnavailable: the index is not set up, or a side is not a token.
var ErrUnavailable = errors.New("order book unavailable")

// CurrentOrderBook is OrderBook on the index set by Use, for two token
// contract addresses (anything else, e.g. the native asset, has no book).
func CurrentOrderBook(base, counter string, limit int) (Book, error) {
	db, client := current()
	b, c := Key(base), Key(counter)
	if db == nil || b == "" || c == "" || b == c {
		return Book{}, ErrUnavailable
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	return OrderBook(ctx, db, client, common.HexToAddress(b), common.HexToAddress(c), limit)
}

// CurrentCandles is Aggregate(Trades(...)) on the index set by Use.
func CurrentCandles(base, counter string, start, end time.Time, resolution, offset time.Duration) ([]Candle, error) {
	db, client := current()
	b, c := Key(base), Key(counter)
	if db == nil || b == "" || c == "" || b == c {
		return nil, ErrUnavailable
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	trades, err := Trades(ctx, db, client, common.HexToAddress(b), common.HexToAddress(c), start, end)
	if err != nil {
		return nil, err
	}
	return Aggregate(trades, resolution, offset), nil
}

// TextLevel is a price level as text, as the order book responses carry it.
type TextLevel struct {
	Price  string
	Amount string
}

// CurrentTextBook is CurrentOrderBook with its levels as text; limit is
// the levels per side ("" or invalid: 20, at most 200).
func CurrentTextBook(base, counter, limit string) (asks, bids []TextLevel, err error) {
	n, e := strconv.Atoi(strings.TrimSpace(limit))
	if e != nil || n <= 0 {
		n = 20
	}
	if n > 200 {
		n = 200
	}
	b, err := CurrentOrderBook(base, counter, n)
	if err != nil {
		return nil, nil, err
	}
	for _, l := range b.Asks {
		asks = append(asks, TextLevel{Price: l.Price.String(), Amount: l.Amount.String()})
	}
	for _, l := range b.Bids {
		bids = append(bids, TextLevel{Price: l.Price.String(), Amount: l.Amount.String()})
	}
	return asks, bids, nil
}
