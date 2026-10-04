package rates

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
)

const onchainABI = `[
 {"name":"decimals","type":"function","stateMutability":"view","inputs":[],"outputs":[{"type":"uint8"}]},
 {"name":"latestRoundData","type":"function","stateMutability":"view","inputs":[],"outputs":[
   {"name":"roundId","type":"uint80"},{"name":"answer","type":"int256"},{"name":"startedAt","type":"uint256"},
   {"name":"updatedAt","type":"uint256"},{"name":"answeredInRound","type":"uint80"}]},
 {"name":"token0","type":"function","stateMutability":"view","inputs":[],"outputs":[{"type":"address"}]},
 {"name":"token1","type":"function","stateMutability":"view","inputs":[],"outputs":[{"type":"address"}]},
 {"name":"observe","type":"function","stateMutability":"view","inputs":[{"name":"secondsAgos","type":"uint32[]"}],"outputs":[
   {"name":"tickCumulatives","type":"int56[]"},{"name":"secondsPerLiquidityCumulativeX128s","type":"uint160[]"}]}
]`

var parsedOnchainABI = mustABI(onchainABI)

func mustABI(s string) abi.ABI {
	a, err := abi.JSON(strings.NewReader(s))
	if err != nil {
		panic(err)
	}
	return a
}

func call(ctx context.Context, c bind.ContractCaller, to common.Address, method string, args ...interface{}) ([]interface{}, error) {
	data, err := parsedOnchainABI.Pack(method, args...)
	if err != nil {
		return nil, err
	}
	out, err := c.CallContract(ctx, ethereum.CallMsg{To: &to, Data: data}, nil)
	if err != nil {
		return nil, fmt.Errorf("%s() on %s: %w", method, to.Hex(), err)
	}
	res, err := parsedOnchainABI.Unpack(method, out)
	if err != nil {
		return nil, fmt.Errorf("%s() on %s: %w", method, to.Hex(), err)
	}
	return res, nil
}

func parseDuration(s, def string) (time.Duration, error) {
	if s == "" {
		s = def
	}
	d, err := time.ParseDuration(s)
	if err != nil || d <= 0 {
		return 0, fmt.Errorf("invalid duration %q", s)
	}
	return d, nil
}

// chainlink:
//
//	{"type": "chainlink", "feed": "0x…", "maxAge": "1h"}
//
// Reads a Chainlink AggregatorV3 feed; answers older than maxAge (default
// 1h - set it just above the feed's heartbeat) are rejected.
type chainlinkSource struct {
	chain  bind.ContractCaller
	feed   common.Address
	maxAge time.Duration
	now    func() time.Time

	mu       sync.Mutex
	loaded   bool
	decimals int32
}

func init() {
	Register("chainlink", func(def json.RawMessage, env Env) (Source, error) {
		var c struct {
			Feed   string `json:"feed"`
			MaxAge string `json:"maxAge"`
		}
		if err := json.Unmarshal(def, &c); err != nil {
			return nil, err
		}
		if !common.IsHexAddress(c.Feed) {
			return nil, fmt.Errorf(`"feed" must be an address`)
		}
		maxAge, err := parseDuration(c.MaxAge, "1h")
		if err != nil {
			return nil, fmt.Errorf(`"maxAge": %w`, err)
		}
		if env.Chain == nil {
			return nil, fmt.Errorf("needs an RPC connection")
		}
		return &chainlinkSource{chain: env.Chain, feed: common.HexToAddress(c.Feed), maxAge: maxAge, now: time.Now}, nil
	})
}

func (s *chainlinkSource) Fetch(ctx context.Context) (decimal.Decimal, error) {
	s.mu.Lock()
	if !s.loaded {
		res, err := call(ctx, s.chain, s.feed, "decimals")
		if err != nil {
			s.mu.Unlock()
			return decimal.Zero, err
		}
		s.decimals, s.loaded = int32(res[0].(uint8)), true
	}
	s.mu.Unlock()
	res, err := call(ctx, s.chain, s.feed, "latestRoundData")
	if err != nil {
		return decimal.Zero, err
	}
	answer := res[1].(*big.Int)
	updatedAt := res[3].(*big.Int)
	if answer.Sign() <= 0 {
		return decimal.Zero, fmt.Errorf("feed %s answered %s", s.feed.Hex(), answer)
	}
	age := s.now().Sub(time.Unix(updatedAt.Int64(), 0))
	if age > s.maxAge {
		return decimal.Zero, fmt.Errorf("feed %s last updated %s ago (maxAge %s)", s.feed.Hex(), age.Round(time.Second), s.maxAge)
	}
	return decimal.NewFromBigInt(answer, -s.decimals), nil
}

// uniswap-v3-twap:
//
//	{"type": "uniswap-v3-twap", "pool": "0x…", "baseToken": "0x…", "window": "30m"}
//
// The time-weighted average price of baseToken, in the pool's other token,
// over window (default 30m). Works with any Uniswap v3-compatible pool that
// implements observe() (Uniswap v3, Aerodrome Slipstream, PancakeSwap v3).
// This is how the cNGN market rate is discovered: e.g. baseToken = USDC on
// the cNGN/USDC pool gives cNGN per USDC, i.e. a market USD/NGN rate.
type uniV3Source struct {
	chain  bind.ContractCaller
	pool   common.Address
	base   common.Address
	window uint32

	mu      sync.Mutex
	loaded  bool
	baseIs0 bool
	dec0    int32
	dec1    int32
}

func init() {
	Register("uniswap-v3-twap", func(def json.RawMessage, env Env) (Source, error) {
		var c struct {
			Pool      string `json:"pool"`
			BaseToken string `json:"baseToken"`
			Window    string `json:"window"`
		}
		if err := json.Unmarshal(def, &c); err != nil {
			return nil, err
		}
		if !common.IsHexAddress(c.Pool) || !common.IsHexAddress(c.BaseToken) {
			return nil, fmt.Errorf(`"pool" and "baseToken" must be addresses`)
		}
		w, err := parseDuration(c.Window, "30m")
		if err != nil || w < time.Minute || w > 7*24*time.Hour {
			return nil, fmt.Errorf(`"window" must be between 1m and 168h`)
		}
		if env.Chain == nil {
			return nil, fmt.Errorf("needs an RPC connection")
		}
		return &uniV3Source{chain: env.Chain, pool: common.HexToAddress(c.Pool), base: common.HexToAddress(c.BaseToken), window: uint32(w / time.Second)}, nil
	})
}

func (s *uniV3Source) load(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.loaded {
		return nil
	}
	r0, err := call(ctx, s.chain, s.pool, "token0")
	if err != nil {
		return err
	}
	r1, err := call(ctx, s.chain, s.pool, "token1")
	if err != nil {
		return err
	}
	t0, t1 := r0[0].(common.Address), r1[0].(common.Address)
	if s.base != t0 && s.base != t1 {
		return fmt.Errorf("baseToken %s is not in pool %s (%s/%s)", s.base.Hex(), s.pool.Hex(), t0.Hex(), t1.Hex())
	}
	d0, err := call(ctx, s.chain, t0, "decimals")
	if err != nil {
		return err
	}
	d1, err := call(ctx, s.chain, t1, "decimals")
	if err != nil {
		return err
	}
	s.baseIs0, s.dec0, s.dec1, s.loaded = s.base == t0, int32(d0[0].(uint8)), int32(d1[0].(uint8)), true
	return nil
}

func (s *uniV3Source) Fetch(ctx context.Context) (decimal.Decimal, error) {
	if err := s.load(ctx); err != nil {
		return decimal.Zero, err
	}
	res, err := call(ctx, s.chain, s.pool, "observe", []uint32{s.window, 0})
	if err != nil {
		return decimal.Zero, fmt.Errorf("%w (the pool may need more observation cardinality for a %ds window)", err, s.window)
	}
	cum := res[0].([]*big.Int)
	if len(cum) != 2 {
		return decimal.Zero, fmt.Errorf("observe returned %d tick cumulatives", len(cum))
	}
	return twapPrice(cum[0], cum[1], s.window, s.baseIs0, s.dec0, s.dec1), nil
}

// twapPrice converts two tick cumulatives into the human-unit price of the
// base token in the other token.
func twapPrice(cumOld, cumNew *big.Int, window uint32, baseIs0 bool, dec0, dec1 int32) decimal.Decimal {
	delta := new(big.Int).Sub(cumNew, cumOld)
	w := big.NewInt(int64(window))
	tick := new(big.Int).Quo(delta, w)
	// round towards negative infinity, like Uniswap's OracleLibrary
	if delta.Sign() < 0 && new(big.Int).Rem(delta, w).Sign() != 0 {
		tick.Sub(tick, big.NewInt(1))
	}
	// raw token1 per raw token0
	raw := math.Pow(1.0001, float64(tick.Int64()))
	// human token1 per human token0
	price := decimal.NewFromFloat(raw).Shift(dec0 - dec1)
	if baseIs0 {
		return price
	}
	return decimal.NewFromInt(1).DivRound(price, 36)
}
