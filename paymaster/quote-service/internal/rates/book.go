package rates

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"trovo-paymaster-quote-service/internal/config"

	"github.com/shopspring/decimal"
)

// SourceResult is one source's latest answer.
type SourceResult struct {
	Name    string          `json:"name"`
	Value   decimal.Decimal `json:"value" swaggertype:"string"`
	Error   string          `json:"error,omitempty"`
	Outlier bool            `json:"outlier,omitempty"`
}

// Price is a pair's aggregated rate.
type Price struct {
	Pair      string          `json:"pair"`
	Value     decimal.Decimal `json:"value" swaggertype:"string"`
	UpdatedAt time.Time       `json:"updatedAt"`
	Sources   []SourceResult  `json:"sources"`
	Error     string          `json:"error,omitempty"` // why the last refresh failed, if it did
}

type pair struct {
	name            string
	sources         []named
	minSources      int
	maxDeviationBps int64

	mu   sync.RWMutex
	last Price // last successful aggregation
	err  string
	seen []SourceResult
}

// Book holds every configured pair and its latest rate.
type Book struct {
	pairs         map[string]*pair
	maxAge        time.Duration
	now           func() time.Time
	sourceTimeout time.Duration
}

// NewBook builds every pair's sources.
func NewBook(pairs map[string]config.PairConfig, env Env, maxAge, sourceTimeout time.Duration) (*Book, error) {
	b := &Book{pairs: map[string]*pair{}, maxAge: maxAge, now: time.Now, sourceTimeout: sourceTimeout}
	for name, pc := range pairs {
		base, quote, err := config.SplitPair(name)
		if err != nil {
			return nil, err
		}
		key := base + "/" + quote
		p := &pair{name: key, minSources: pc.MinSources, maxDeviationBps: pc.MaxDeviationBps}
		if p.minSources <= 0 {
			p.minSources = 1
		}
		if p.maxDeviationBps == 0 {
			p.maxDeviationBps = 500
		}
		seen := map[string]bool{}
		for _, def := range pc.Sources {
			srcName, src, err := Build(def, env)
			if err != nil {
				return nil, fmt.Errorf("pair %s: %w", key, err)
			}
			if seen[srcName] {
				return nil, fmt.Errorf("pair %s: two sources named %q", key, srcName)
			}
			seen[srcName] = true
			p.sources = append(p.sources, src.(named))
		}
		b.pairs[key] = p
	}
	return b, nil
}

// Refresh fetches every pair's sources concurrently and aggregates them. A
// pair that fails keeps its last good price (until it is too old to use).
func (b *Book) Refresh(ctx context.Context) {
	var wg sync.WaitGroup
	for _, p := range b.pairs {
		wg.Add(1)
		go func(p *pair) {
			defer wg.Done()
			b.refreshPair(ctx, p)
		}(p)
	}
	wg.Wait()
}

func (b *Book) refreshPair(ctx context.Context, p *pair) {
	results := make([]SourceResult, len(p.sources))
	var wg sync.WaitGroup
	for i, s := range p.sources {
		wg.Add(1)
		go func(i int, s named) {
			defer wg.Done()
			sctx, cancel := context.WithTimeout(ctx, b.sourceTimeout)
			defer cancel()
			v, err := s.Fetch(sctx)
			results[i] = SourceResult{Name: s.name, Value: v}
			if err != nil {
				results[i].Error = err.Error()
			}
		}(i, s)
	}
	wg.Wait()

	value, err := aggregate(results, p.minSources, p.maxDeviationBps)
	p.mu.Lock()
	defer p.mu.Unlock()
	p.seen = results
	if err != nil {
		p.err = err.Error()
		return
	}
	p.err = ""
	p.last = Price{Pair: p.name, Value: value, UpdatedAt: b.now(), Sources: results}
}

// aggregate returns the median of the answering sources, after dropping
// those further than maxDeviationBps from it (negative disables), and
// requires at least minSources to remain. It marks outliers in results.
func aggregate(results []SourceResult, minSources int, maxDeviationBps int64) (decimal.Decimal, error) {
	ok := make([]int, 0, len(results))
	for i, r := range results {
		if r.Error == "" && r.Value.IsPositive() {
			ok = append(ok, i)
		}
	}
	if len(ok) == 0 {
		return decimal.Zero, fmt.Errorf("no source answered")
	}
	vals := func(idx []int) []decimal.Decimal {
		out := make([]decimal.Decimal, len(idx))
		for j, i := range idx {
			out[j] = results[i].Value
		}
		return out
	}
	med := median(vals(ok))
	kept := ok
	if maxDeviationBps > 0 {
		kept = kept[:0:0]
		limit := med.Mul(decimal.NewFromInt(maxDeviationBps)).Div(decimal.NewFromInt(10000))
		for _, i := range ok {
			if results[i].Value.Sub(med).Abs().GreaterThan(limit) {
				results[i].Outlier = true
				continue
			}
			kept = append(kept, i)
		}
		med = median(vals(kept))
	}
	if len(kept) < minSources {
		return decimal.Zero, fmt.Errorf("only %d of the required %d sources answered within the deviation limit", len(kept), minSources)
	}
	return med, nil
}

func median(v []decimal.Decimal) decimal.Decimal {
	if len(v) == 0 {
		return decimal.Zero
	}
	s := append([]decimal.Decimal(nil), v...)
	sort.Slice(s, func(i, j int) bool { return s[i].LessThan(s[j]) })
	if len(s)%2 == 1 {
		return s[len(s)/2]
	}
	return s[len(s)/2-1].Add(s[len(s)/2]).Div(decimal.NewFromInt(2))
}

// Get returns a pair's rate ("BASE/QUOTE"; the inverse of a configured
// pair is derived), refusing one older than the maximum age.
func (b *Book) Get(name string) (decimal.Decimal, time.Time, error) {
	base, quote, err := config.SplitPair(name)
	if err != nil {
		return decimal.Zero, time.Time{}, err
	}
	key, inverted, ok := config.LookupPair(b.pairs, base, quote)
	if !ok {
		return decimal.Zero, time.Time{}, fmt.Errorf("pair %s/%s is not configured", base, quote)
	}
	p := b.pairs[key]
	p.mu.RLock()
	last, lastErr := p.last, p.err
	p.mu.RUnlock()
	if last.UpdatedAt.IsZero() {
		if lastErr == "" {
			lastErr = "not fetched yet"
		}
		return decimal.Zero, time.Time{}, fmt.Errorf("rate %s unavailable: %s", key, lastErr)
	}
	if age := b.now().Sub(last.UpdatedAt); age > b.maxAge {
		return decimal.Zero, time.Time{}, fmt.Errorf("rate %s is stale (%s old; last error: %s)", key, age.Round(time.Second), lastErr)
	}
	if inverted {
		return decimal.NewFromInt(1).DivRound(last.Value, 36), last.UpdatedAt, nil
	}
	return last.Value, last.UpdatedAt, nil
}

// Convert multiplies the rates along a route (e.g. ETH/USD, USD/NGN) and
// returns the result with the oldest rate's time.
func (b *Book) Convert(route []string) (decimal.Decimal, time.Time, error) {
	v := decimal.NewFromInt(1)
	var oldest time.Time
	for _, leg := range route {
		r, at, err := b.Get(leg)
		if err != nil {
			return decimal.Zero, time.Time{}, err
		}
		v = v.Mul(r)
		if oldest.IsZero() || at.Before(oldest) {
			oldest = at
		}
	}
	return v, oldest, nil
}

// Snapshot returns every pair's state, for the rates endpoint.
func (b *Book) Snapshot() []Price {
	out := make([]Price, 0, len(b.pairs))
	for _, p := range b.pairs {
		p.mu.RLock()
		pr := p.last
		pr.Pair = p.name
		if p.seen != nil {
			pr.Sources = p.seen
		}
		pr.Error = p.err
		p.mu.RUnlock()
		out = append(out, pr)
	}
	sort.Slice(out, func(i, j int) bool { return strings.Compare(out[i].Pair, out[j].Pair) < 0 })
	return out
}
