// Package rates discovers exchange rates from pluggable sources.
//
// A rate pair "BASE/QUOTE" (e.g. "ETH/USD", "USD/NGN") is priced as the
// number of QUOTE units one BASE unit is worth. Each pair has one or more
// sources; the pair's rate is the median of the sources that answered,
// after dropping outliers.
//
// Source types (the "type" field of a source definition):
//
//	fixed          a constant, e.g. a 1:1 peg
//	http-json      a number read from a JSON HTTP API (exchanges, FX feeds, issuers)
//	chainlink      a Chainlink AggregatorV3 price feed
//	uniswap-v3-twap the time-weighted average price of a Uniswap v3 (or
//	               Aerodrome Slipstream / any v3-compatible) pool - e.g. the
//	               cNGN/USDC pool for cNGN market-rate discovery
//
// Every source also accepts "invert": true (use 1/x) and "multiplier" (a
// decimal factor). New types are added with Register.
package rates

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/shopspring/decimal"
)

// Source fetches one rate.
type Source interface {
	Fetch(ctx context.Context) (decimal.Decimal, error)
}

// Env is what sources may use.
type Env struct {
	Chain  bind.ContractCaller             // on-chain reads (chainlink, uniswap-v3-twap)
	HTTP   *http.Client                    // http-json
	Secret func(key string) (string, bool) // resolves "vault:KEY" references
}

// Factory builds a source from its JSON definition.
type Factory func(def json.RawMessage, env Env) (Source, error)

var registry = map[string]Factory{}

// Register adds a source type. It panics on duplicates (programming error).
func Register(typ string, f Factory) {
	if _, dup := registry[typ]; dup {
		panic("rates: duplicate source type " + typ)
	}
	registry[typ] = f
}

// Types lists the registered source types.
func Types() []string {
	out := make([]string, 0, len(registry))
	for t := range registry {
		out = append(out, t)
	}
	sort.Strings(out)
	return out
}

// common fields of every source definition
type header struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	Invert     bool   `json:"invert"`
	Multiplier string `json:"multiplier"`
}

// named is a built source with its generic adjustments applied.
type named struct {
	name       string
	src        Source
	invert     bool
	multiplier decimal.Decimal
}

func (n named) Fetch(ctx context.Context) (decimal.Decimal, error) {
	v, err := n.src.Fetch(ctx)
	if err != nil {
		return decimal.Zero, err
	}
	if !v.IsPositive() {
		return decimal.Zero, fmt.Errorf("non-positive rate %s", v)
	}
	if n.invert {
		v = decimal.NewFromInt(1).DivRound(v, 36)
	}
	return v.Mul(n.multiplier), nil
}

// Build builds a source from its definition.
func Build(def json.RawMessage, env Env) (string, Source, error) {
	var h header
	if err := json.Unmarshal(def, &h); err != nil {
		return "", nil, fmt.Errorf("invalid source definition: %w", err)
	}
	if h.Type == "" {
		return "", nil, fmt.Errorf("source definition needs a type (one of %s)", strings.Join(Types(), ", "))
	}
	f, ok := registry[h.Type]
	if !ok {
		return "", nil, fmt.Errorf("unknown source type %q (known: %s)", h.Type, strings.Join(Types(), ", "))
	}
	name := h.Name
	if name == "" {
		name = h.Type
	}
	src, err := f(def, env)
	if err != nil {
		return "", nil, fmt.Errorf("source %s: %w", name, err)
	}
	mult := decimal.NewFromInt(1)
	if h.Multiplier != "" {
		m, err := decimal.NewFromString(h.Multiplier)
		if err != nil || !m.IsPositive() {
			return "", nil, fmt.Errorf("source %s: multiplier must be a positive decimal", name)
		}
		mult = m
	}
	return name, named{name: name, src: src, invert: h.Invert, multiplier: mult}, nil
}

// resolveSecret replaces a "vault:KEY" reference with the configuration
// secret's KEY.
func resolveSecret(v string, env Env) (string, error) {
	key, ok := strings.CutPrefix(v, "vault:")
	if !ok {
		return v, nil
	}
	if env.Secret == nil {
		return "", fmt.Errorf("secret %s referenced but no secret store configured", key)
	}
	s, ok := env.Secret(key)
	if !ok {
		return "", fmt.Errorf("secret %s referenced but not set", key)
	}
	return s, nil
}
