package publicmarkets

import (
	"context"
	"sort"
	"time"

	pm "trovo-wallet-api/internal/components/publicmarkets/models"

	"github.com/shopspring/decimal"
)

// PollPrices ingests reference prices: live quotes during each market's
// session, and once the session closes the last close, labelled
// LAST_CLOSE_SYNTHETIC (never passed off as live).
func (e *Engine) PollPrices(ctx context.Context) (stale int) {
	s := LoadSettings(e.DB)
	var assets []pm.Asset
	e.DB.Where("status IN ? AND isin <> ''", []string{pm.AssetLive, pm.AssetHalted, pm.AssetSetup}).Find(&assets)
	for i := range assets {
		a := &assets[i]
		sess := MarketSession(a.Market, e.now(), s)
		if !sess.Open {
			e.closeSession(a)
			continue
		}
		q, err := e.Prices.Quote(ctx, a.ISIN, a.Market)
		if err != nil || !q.MarketOpen || !q.Price.IsPositive() {
			if a.PriceAt == nil || e.now().Sub(*a.PriceAt) > time.Duration(s.PriceStaleMinutes)*time.Minute {
				stale++
			}
			if err != nil {
				e.logf("price %s: %v", a.AssetCode, err)
			}
			continue
		}
		source := pm.PriceNGXFeed
		if a.Market == pm.MarketFMDQ {
			source = pm.PriceFMDQFeed
		}
		if _, mock := e.Prices.(MockPriceFeed); mock {
			source = pm.PriceMockFeed
		}
		e.recordPrice(a, q.Price, source, true, q.AsOf)
	}
	return stale
}

// recordPrice stores a snapshot and the asset's listing fields; the first
// live price of a session starts the day (previous close, high/low).
func (e *Engine) recordPrice(a *pm.Asset, price decimal.Decimal, source string, live bool, asOf time.Time) {
	now := e.now()
	if asOf.IsZero() {
		asOf = now
	}
	e.DB.Create(&pm.PriceSnapshot{AssetID: a.ID, Price: price.String(), Source: source, MarketHours: live, AsOf: asOf, CapturedAt: now})
	updates := map[string]interface{}{"last_price": price.String(), "price_source": source, "price_at": &now}
	newDay := a.PriceAt == nil || a.PriceSource == pm.PriceSynthetic || a.PriceAt.In(Lagos).Format("2006-01-02") != now.In(Lagos).Format("2006-01-02")
	if live && newDay {
		if d(a.LastPrice).IsPositive() {
			updates["previous_close"] = a.LastPrice
		}
		updates["day_high"], updates["day_low"] = price.String(), price.String()
	} else if live {
		if price.GreaterThan(d(a.DayHigh)) {
			updates["day_high"] = price.String()
		}
		if low := d(a.DayLow); !low.IsPositive() || price.LessThan(low) {
			updates["day_low"] = price.String()
		}
	}
	e.DB.Model(&pm.Asset{}).Where("id = ?", a.ID).Updates(updates)
}

// closeSession labels the last price as the synthetic close once the
// session is over.
func (e *Engine) closeSession(a *pm.Asset) {
	if a.PriceSource == pm.PriceSynthetic || !d(a.LastPrice).IsPositive() {
		return
	}
	e.recordPrice(a, d(a.LastPrice), pm.PriceSynthetic, false, e.now())
}

// PricePoint is one point of a price chart.
type PricePoint struct {
	At    time.Time `json:"at"`
	Price string    `json:"price"`
}

// PriceHistory is an asset's price chart over a range (1D, 1W, 1M, 3M,
// 1Y, All): every snapshot of the last session for 1D, else the last
// price of each hour (1W) or day.
func (e *Engine) PriceHistory(a *pm.Asset, rng string) []PricePoint {
	now := e.now()
	from, bucket := now.AddDate(0, 0, -1), time.Duration(0)
	switch rng {
	case "1W":
		from, bucket = now.AddDate(0, 0, -7), time.Hour
	case "1M":
		from, bucket = now.AddDate(0, -1, 0), 24*time.Hour
	case "3M":
		from, bucket = now.AddDate(0, -3, 0), 24*time.Hour
	case "1Y":
		from, bucket = now.AddDate(-1, 0, 0), 24*time.Hour
	case "All":
		from, bucket = time.Time{}, 24*time.Hour
	}
	var snaps []pm.PriceSnapshot
	e.DB.Where("asset_id = ? AND captured_at >= ?", a.ID, from).Order("captured_at").Find(&snaps)
	if bucket == 0 {
		out := make([]PricePoint, 0, len(snaps))
		for _, s := range snaps {
			out = append(out, PricePoint{At: s.CapturedAt, Price: s.Price})
		}
		return out
	}
	last := map[int64]pm.PriceSnapshot{}
	for _, s := range snaps {
		last[s.CapturedAt.In(Lagos).Truncate(bucket).Unix()] = s
	}
	keys := make([]int64, 0, len(last))
	for k := range last {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	out := make([]PricePoint, 0, len(keys))
	for _, k := range keys {
		out = append(out, PricePoint{At: last[k].CapturedAt, Price: last[k].Price})
	}
	return out
}

// DayChange is the change against the previous close, in percent.
func DayChange(a *pm.Asset) decimal.Decimal {
	prev, last := d(a.PreviousClose), d(a.LastPrice)
	if !prev.IsPositive() || !last.IsPositive() {
		return decimal.Zero
	}
	return last.Sub(prev).Div(prev).Mul(decimal.NewFromInt(100)).Round(2)
}
