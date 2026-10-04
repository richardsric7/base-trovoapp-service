package config

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/shopspring/decimal"
)

// DefaultEntryPoint is the canonical ERC-4337 v0.7 EntryPoint.
const DefaultEntryPoint = "0x0000000071727De22E5E9d8BAf0edAc6f37da032"

// GasToken is one stablecoin wallets may pay gas in.
type GasToken struct {
	Symbol   string         `json:"symbol"`
	Address  common.Address `json:"address"`
	Decimals uint8          `json:"decimals"`
	// Route converts 1 ETH into the token's unit through configured rate
	// pairs, e.g. ["ETH/USD", "USD/NGN", "NGN/CNGN"]. A leg may also be the
	// inverse of a configured pair ("USD/ETH" uses "ETH/USD").
	Route []string `json:"route"`
	// SpreadBps is Trovo's spread over the market rate, in basis points
	// (100 = 1%). Nil falls back to DEFAULT_SPREAD_BPS.
	SpreadBps *int64 `json:"spreadBps,omitempty"`
}

// PairConfig configures how one rate pair ("BASE/QUOTE": 1 BASE = x QUOTE)
// is discovered.
type PairConfig struct {
	// Sources are rate source definitions, each {"name", "type", ...type
	// specific fields}; see internal/rates for the types.
	Sources []json.RawMessage `json:"sources"`
	// MinSources is how many sources must agree for the pair to be usable
	// (default 1).
	MinSources int `json:"minSources"`
	// MaxDeviationBps drops sources further than this from the median
	// (default 500 = 5%; 0 keeps the default, negative disables).
	MaxDeviationBps int64 `json:"maxDeviationBps"`
}

type Config struct {
	HTTPAddr string
	APIKeys  []string

	RPCURL     string
	ChainID    *big.Int
	EntryPoint common.Address
	Paymaster  common.Address
	SignerKey  *ecdsa.PrivateKey

	GasTokens        []GasToken
	RatePairs        map[string]PairConfig
	DefaultSpreadBps int64

	RateRefreshInterval time.Duration
	RateMaxAge          time.Duration
	SourceTimeout       time.Duration

	QuoteValidity    time.Duration
	QuoteMaxValidity time.Duration
	ClockSkew        time.Duration

	PaymasterVerificationGasLimit uint64
	PaymasterPostOpGasLimit       uint64

	DepositCheckInterval time.Duration
	DepositLowWatermark  *big.Int // wei
	AlertWebhookURL      string
	AlertCooldown        time.Duration

	ConfigReloadInterval time.Duration

	// Secret resolves "vault:KEY" references in source definitions (e.g.
	// API keys in HTTP headers) against the same configuration secret.
	Secret func(key string) (string, bool)
}

// Parse validates the configuration values. Required keys: RPC_URL,
// CHAIN_ID, PAYMASTER_ADDRESS, QUOTE_SIGNER_PRIVATE_KEY, API_KEYS,
// GAS_TOKENS, RATE_PAIRS.
func Parse(v Values) (*Config, error) {
	p := parser{v: v}
	c := &Config{
		HTTPAddr:             p.str("HTTP_ADDR", ":8090"),
		RPCURL:               p.required("RPC_URL"),
		ChainID:              p.bigInt("CHAIN_ID", ""),
		EntryPoint:           p.address("ENTRYPOINT_ADDRESS", DefaultEntryPoint),
		Paymaster:            p.address("PAYMASTER_ADDRESS", ""),
		DefaultSpreadBps:     p.int("DEFAULT_SPREAD_BPS", 0),
		RateRefreshInterval:  p.duration("RATE_REFRESH_INTERVAL", "30s"),
		RateMaxAge:           p.duration("RATE_MAX_AGE", "5m"),
		SourceTimeout:        p.duration("RATE_SOURCE_TIMEOUT", "10s"),
		QuoteValidity:        p.duration("QUOTE_VALIDITY", "10m"),
		QuoteMaxValidity:     p.duration("QUOTE_MAX_VALIDITY", "24h"),
		ClockSkew:            p.duration("QUOTE_CLOCK_SKEW", "60s"),
		DepositCheckInterval: p.duration("DEPOSIT_CHECK_INTERVAL", "1m"),
		DepositLowWatermark:  p.eth("DEPOSIT_LOW_WATERMARK_ETH", "0.05"),
		AlertWebhookURL:      p.str("ALERT_WEBHOOK_URL", ""),
		AlertCooldown:        p.duration("ALERT_COOLDOWN", "1h"),
		ConfigReloadInterval: p.duration("CONFIG_RELOAD_INTERVAL", "0s"),
		Secret: func(key string) (string, bool) {
			s, ok := v[key]
			return s, ok && s != ""
		},
	}
	c.PaymasterVerificationGasLimit = uint64(p.int("PAYMASTER_VERIFICATION_GAS_LIMIT", 150000))
	c.PaymasterPostOpGasLimit = uint64(p.int("PAYMASTER_POST_OP_GAS_LIMIT", 80000))

	for _, k := range strings.Split(p.required("API_KEYS"), ",") {
		if k = strings.TrimSpace(k); k != "" {
			if len(k) < 16 {
				p.fail("API_KEYS entries must be at least 16 characters")
			}
			c.APIKeys = append(c.APIKeys, k)
		}
	}
	if len(c.APIKeys) == 0 && p.err == nil {
		p.fail("API_KEYS has no keys")
	}

	if raw := p.required("QUOTE_SIGNER_PRIVATE_KEY"); raw != "" {
		key, err := crypto.HexToECDSA(strings.TrimPrefix(strings.TrimSpace(raw), "0x"))
		if err != nil {
			p.fail("QUOTE_SIGNER_PRIVATE_KEY is not a valid secp256k1 private key")
		}
		c.SignerKey = key
	}

	if raw := p.required("GAS_TOKENS"); raw != "" {
		if err := json.Unmarshal([]byte(raw), &c.GasTokens); err != nil {
			p.fail("GAS_TOKENS is not valid JSON: %v", err)
		}
	}
	if raw := p.required("RATE_PAIRS"); raw != "" {
		if err := json.Unmarshal([]byte(raw), &c.RatePairs); err != nil {
			p.fail("RATE_PAIRS is not valid JSON: %v", err)
		}
	}
	if p.err != nil {
		return nil, p.err
	}
	if err := c.validate(); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Config) validate() error {
	if c.ChainID.Sign() <= 0 {
		return fmt.Errorf("CHAIN_ID must be positive")
	}
	if c.QuoteValidity <= 0 || c.QuoteMaxValidity < c.QuoteValidity {
		return fmt.Errorf("need 0 < QUOTE_VALIDITY <= QUOTE_MAX_VALIDITY")
	}
	if c.QuoteMaxValidity > 7*24*time.Hour {
		return fmt.Errorf("QUOTE_MAX_VALIDITY must not exceed 168h")
	}
	if c.RateRefreshInterval <= 0 || c.RateMaxAge < c.RateRefreshInterval {
		return fmt.Errorf("need 0 < RATE_REFRESH_INTERVAL <= RATE_MAX_AGE")
	}
	if c.DefaultSpreadBps < 0 || c.DefaultSpreadBps > 5000 {
		return fmt.Errorf("DEFAULT_SPREAD_BPS must be between 0 and 5000")
	}
	if len(c.RatePairs) == 0 {
		return fmt.Errorf("RATE_PAIRS has no pairs")
	}
	for name, pc := range c.RatePairs {
		if _, _, err := SplitPair(name); err != nil {
			return fmt.Errorf("RATE_PAIRS: %w", err)
		}
		if len(pc.Sources) == 0 {
			return fmt.Errorf("RATE_PAIRS[%s] has no sources", name)
		}
		if pc.MinSources > len(pc.Sources) {
			return fmt.Errorf("RATE_PAIRS[%s].minSources exceeds its number of sources", name)
		}
	}
	if len(c.GasTokens) == 0 {
		return fmt.Errorf("GAS_TOKENS has no tokens")
	}
	seen := map[common.Address]bool{}
	for i, t := range c.GasTokens {
		if t.Symbol == "" || t.Address == (common.Address{}) {
			return fmt.Errorf("GAS_TOKENS[%d] needs a symbol and an address", i)
		}
		if seen[t.Address] {
			return fmt.Errorf("GAS_TOKENS lists %s twice", t.Address.Hex())
		}
		seen[t.Address] = true
		if t.Decimals > 36 {
			return fmt.Errorf("GAS_TOKENS[%s].decimals is out of range", t.Symbol)
		}
		if t.SpreadBps != nil && (*t.SpreadBps < 0 || *t.SpreadBps > 5000) {
			return fmt.Errorf("GAS_TOKENS[%s].spreadBps must be between 0 and 5000", t.Symbol)
		}
		if err := ValidateRoute(t.Route, c.RatePairs); err != nil {
			return fmt.Errorf("GAS_TOKENS[%s].route: %w", t.Symbol, err)
		}
	}
	return nil
}

// Spread returns the token's spread in basis points.
func (c *Config) Spread(t GasToken) int64 {
	if t.SpreadBps != nil {
		return *t.SpreadBps
	}
	return c.DefaultSpreadBps
}

// SplitPair splits "BASE/QUOTE".
func SplitPair(name string) (string, string, error) {
	base, quote, ok := strings.Cut(name, "/")
	if !ok || base == "" || quote == "" || strings.Contains(quote, "/") || base == quote {
		return "", "", fmt.Errorf("pair %q must look like BASE/QUOTE", name)
	}
	return strings.ToUpper(base), strings.ToUpper(quote), nil
}

// ValidateRoute checks that route starts at ETH, that each leg continues
// from the previous one, and that every leg (or its inverse) is configured.
func ValidateRoute(route []string, pairs map[string]PairConfig) error {
	if len(route) == 0 {
		return fmt.Errorf("is empty")
	}
	at := "ETH"
	for _, leg := range route {
		base, quote, err := SplitPair(leg)
		if err != nil {
			return err
		}
		if base != at {
			return fmt.Errorf("leg %s does not continue from %s", leg, at)
		}
		if _, _, ok := LookupPair(pairs, base, quote); !ok {
			return fmt.Errorf("neither %s/%s nor %s/%s is configured in RATE_PAIRS", base, quote, quote, base)
		}
		at = quote
	}
	return nil
}

// LookupPair finds the configured pair for base/quote: its name, and
// whether it is configured the other way round (inverted).
func LookupPair[T any](pairs map[string]T, base, quote string) (name string, inverted bool, ok bool) {
	for n := range pairs {
		b, q, err := SplitPair(n)
		if err != nil {
			continue
		}
		if b == base && q == quote {
			return n, false, true
		}
		if b == quote && q == base {
			return n, true, true
		}
	}
	return "", false, false
}

type parser struct {
	v   Values
	err error
}

func (p *parser) fail(format string, args ...interface{}) {
	if p.err == nil {
		p.err = fmt.Errorf(format, args...)
	}
}

func (p *parser) str(k, def string) string {
	if s := strings.TrimSpace(p.v[k]); s != "" {
		return s
	}
	return def
}

func (p *parser) required(k string) string {
	s := p.str(k, "")
	if s == "" {
		p.fail("%s is required", k)
	}
	return s
}

func (p *parser) address(k, def string) common.Address {
	s := p.str(k, def)
	if s == "" {
		p.fail("%s is required", k)
		return common.Address{}
	}
	if !common.IsHexAddress(s) {
		p.fail("%s is not an address: %q", k, s)
	}
	return common.HexToAddress(s)
}

func (p *parser) int(k string, def int64) int64 {
	s := p.str(k, "")
	if s == "" {
		return def
	}
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		p.fail("%s must be an integer: %q", k, s)
	}
	return n
}

func (p *parser) bigInt(k, def string) *big.Int {
	s := p.str(k, def)
	if s == "" {
		p.fail("%s is required", k)
		return new(big.Int)
	}
	n, ok := new(big.Int).SetString(s, 10)
	if !ok {
		p.fail("%s must be an integer: %q", k, s)
		return new(big.Int)
	}
	return n
}

func (p *parser) duration(k, def string) time.Duration {
	s := p.str(k, def)
	d, err := time.ParseDuration(s)
	if err != nil || d < 0 {
		p.fail("%s must be a duration like 30s or 10m: %q", k, s)
	}
	return d
}

func (p *parser) eth(k, def string) *big.Int {
	s := p.str(k, def)
	d, err := decimal.NewFromString(s)
	if err != nil || d.IsNegative() {
		p.fail("%s must be an ETH amount like 0.05: %q", k, s)
		return new(big.Int)
	}
	return d.Shift(18).Truncate(0).BigInt()
}
