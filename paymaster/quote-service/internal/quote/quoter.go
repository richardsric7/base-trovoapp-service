package quote

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"net/http"
	"time"

	"trovo-paymaster-quote-service/internal/chain"
	"trovo-paymaster-quote-service/internal/config"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/shopspring/decimal"
)

// Error is a quote failure with the HTTP status it maps to.
type Error struct {
	Status  int    `json:"-"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *Error) Error() string { return e.Code + ": " + e.Message }

func errf(status int, code, format string, args ...interface{}) *Error {
	return &Error{Status: status, Code: code, Message: fmt.Sprintf(format, args...)}
}

// Rates converts 1 ETH along a route (see rates.Book).
type Rates interface {
	Convert(route []string) (decimal.Decimal, time.Time, error)
}

// Chain is the on-chain state a quote depends on (see chain.Paymaster).
type Chain interface {
	TokenBounds(ctx context.Context, token common.Address) (chain.TokenBounds, error)
	HasDebt(ctx context.Context, sender common.Address) (bool, error)
}

// Deposit reports the last known paymaster deposit (wei), if known.
type Deposit func() (*big.Int, bool)

type Quoter struct {
	ChainID     *big.Int
	Paymaster   common.Address
	Key         *ecdsa.PrivateKey
	Tokens      []config.GasToken
	Spread      func(config.GasToken) int64
	Rates       Rates
	Chain       Chain
	Deposit     Deposit
	Validity    time.Duration
	MaxValidity time.Duration
	ClockSkew   time.Duration

	PaymasterVerificationGasLimit uint64
	PaymasterPostOpGasLimit       uint64

	Now func() time.Time
}

// Signer is the address quotes are signed with.
func (q *Quoter) Signer() common.Address { return crypto.PubkeyToAddress(q.Key.PublicKey) }

// Token finds a configured gas token by address or symbol.
func (q *Quoter) Token(ref string) (config.GasToken, bool) {
	for _, t := range q.Tokens {
		if (common.IsHexAddress(ref) && common.HexToAddress(ref) == t.Address) || ref == t.Symbol {
			return t, true
		}
	}
	return config.GasToken{}, false
}

// TokenRate is a gas token's current price of 1 ETH.
type TokenRate struct {
	Token        common.Address  `json:"token" swaggertype:"string"`
	Symbol       string          `json:"symbol"`
	Decimals     uint8           `json:"decimals"`
	MarketRate   decimal.Decimal `json:"marketRate" swaggertype:"string"` // tokens per ETH at market
	QuotedRate   decimal.Decimal `json:"quotedRate" swaggertype:"string"` // tokens per ETH with the spread
	ExchangeRate *big.Int        `json:"-"`
	// QuotedRate in token base units per 1e18 wei, hex: what the paymaster uses
	ExchangeRateHex *hexutil.Big `json:"exchangeRate" swaggertype:"string"`
	SpreadBps       int64        `json:"spreadBps"`
	RatesAsOf       time.Time    `json:"ratesAsOf"`
}

// Rate prices 1 ETH in t, with t's spread applied.
func (q *Quoter) Rate(t config.GasToken) (TokenRate, error) {
	market, asOf, err := q.Rates.Convert(t.Route)
	if err != nil {
		return TokenRate{}, err
	}
	spread := q.Spread(t)
	quoted := market.Mul(decimal.NewFromInt(10000 + spread)).Div(decimal.NewFromInt(10000))
	base := quoted.Shift(int32(t.Decimals)).Truncate(0)
	if !base.IsPositive() {
		return TokenRate{}, fmt.Errorf("rate for %s rounds to zero base units", t.Symbol)
	}
	return TokenRate{
		Token: t.Address, Symbol: t.Symbol, Decimals: t.Decimals,
		MarketRate: market.Round(int32(t.Decimals)), QuotedRate: base.Shift(-int32(t.Decimals)),
		ExchangeRate: base.BigInt(), ExchangeRateHex: (*hexutil.Big)(base.BigInt()), SpreadBps: spread, RatesAsOf: asOf,
	}, nil
}

// TokenCost is TrovoTokenPaymaster.tokenCost: ceil(wei * rate / 1e18).
func TokenCost(wei, rate *big.Int) *big.Int {
	n := new(big.Int).Mul(wei, rate)
	d := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
	q, r := new(big.Int).QuoRem(n, d, new(big.Int))
	if r.Sign() != 0 {
		q.Add(q, big.NewInt(1))
	}
	return q
}

// Request asks for a quote for Op, paid in Token.
type Request struct {
	Op       UserOp
	Token    string
	Validity time.Duration // 0 = the default validity
}

// Result is a signed quote, ready to be put in the UserOperation.
type Result struct {
	Paymaster                     common.Address  `json:"paymaster" swaggertype:"string"`
	PaymasterVerificationGasLimit *hexutil.Big    `json:"paymasterVerificationGasLimit" swaggertype:"string"`
	PaymasterPostOpGasLimit       *hexutil.Big    `json:"paymasterPostOpGasLimit" swaggertype:"string"`
	PaymasterData                 hexutil.Bytes   `json:"paymasterData" swaggertype:"string"`
	PaymasterAndData              hexutil.Bytes   `json:"paymasterAndData" swaggertype:"string"`
	Token                         common.Address  `json:"token" swaggertype:"string"`
	Symbol                        string          `json:"symbol"`
	ExchangeRate                  *hexutil.Big    `json:"exchangeRate" swaggertype:"string"`
	MarketRate                    decimal.Decimal `json:"marketRate" swaggertype:"string"`
	QuotedRate                    decimal.Decimal `json:"quotedRate" swaggertype:"string"`
	SpreadBps                     int64           `json:"spreadBps"`
	ValidAfter                    uint64          `json:"validAfter"`
	ValidUntil                    uint64          `json:"validUntil"`
	MaxCostWei                    *hexutil.Big    `json:"maxCostWei" swaggertype:"string"`
	MaxTokenCost                  *hexutil.Big    `json:"maxTokenCost" swaggertype:"string"`
	MaxTokenCostFormatted         decimal.Decimal `json:"maxTokenCostFormatted" swaggertype:"string"`
	QuoteHash                     common.Hash     `json:"quoteHash" swaggertype:"string"`
	Signer                        common.Address  `json:"signer" swaggertype:"string"`
}

// Quote prices and signs a quote for req.
func (q *Quoter) Quote(ctx context.Context, req Request) (*Result, error) {
	t, ok := q.Token(req.Token)
	if !ok {
		return nil, errf(http.StatusBadRequest, "unknown-token", "%s is not a gas token", req.Token)
	}
	validity := req.Validity
	if validity == 0 {
		validity = q.Validity
	}
	if validity < 0 || validity > q.MaxValidity {
		return nil, errf(http.StatusBadRequest, "invalid-validity", "validity must be between 1s and %s", q.MaxValidity)
	}

	op := req.Op
	if op.PaymasterVerificationGasLimit == nil {
		op.PaymasterVerificationGasLimit = new(big.Int).SetUint64(q.PaymasterVerificationGasLimit)
	}
	if op.PaymasterPostOpGasLimit == nil {
		op.PaymasterPostOpGasLimit = new(big.Int).SetUint64(q.PaymasterPostOpGasLimit)
	}
	if err := op.Validate(); err != nil {
		return nil, errf(http.StatusBadRequest, "invalid-user-operation", "%v", err)
	}

	rate, err := q.Rate(t)
	if err != nil {
		return nil, errf(http.StatusServiceUnavailable, "rate-unavailable", "%v", err)
	}
	bounds, err := q.Chain.TokenBounds(ctx, t.Address)
	if err != nil {
		return nil, errf(http.StatusServiceUnavailable, "chain-unavailable", "reading the paymaster: %v", err)
	}
	if !bounds.Enabled {
		return nil, errf(http.StatusServiceUnavailable, "token-disabled", "%s is not enabled on the paymaster", t.Symbol)
	}
	if rate.ExchangeRate.Cmp(bounds.MinRate) < 0 || rate.ExchangeRate.Cmp(bounds.MaxRate) > 0 {
		return nil, errf(http.StatusServiceUnavailable, "rate-out-of-bounds",
			"%s rate %s is outside the paymaster's bounds [%s, %s]; the paymaster owner must update setToken", t.Symbol, rate.ExchangeRate, bounds.MinRate, bounds.MaxRate)
	}
	if debt, err := q.Chain.HasDebt(ctx, op.Sender); err != nil {
		return nil, errf(http.StatusServiceUnavailable, "chain-unavailable", "reading the paymaster: %v", err)
	} else if debt {
		return nil, errf(http.StatusConflict, "outstanding-debt", "%s owes the paymaster for an earlier operation; settle it (settleDebt) or pay this operation's gas in ETH", op.Sender.Hex())
	}

	maxCost := op.MaxCost()
	if dep, known := q.Deposit(); known && dep.Cmp(maxCost) < 0 {
		return nil, errf(http.StatusServiceUnavailable, "paymaster-deposit-low", "the paymaster deposit cannot cover this operation right now")
	}

	now := q.Now()
	validAfter := uint64(now.Add(-q.ClockSkew).Unix())
	validUntil := uint64(now.Add(validity).Unix())
	hash := Hash(op, q.ChainID, q.Paymaster, t.Address, validUntil, validAfter, rate.ExchangeRate)
	sig, err := crypto.Sign(accounts.TextHash(hash.Bytes()), q.Key)
	if err != nil {
		return nil, errf(http.StatusInternalServerError, "signing-failed", "%v", err)
	}
	sig[64] += 27

	data := PaymasterData(t.Address, validUntil, validAfter, rate.ExchangeRate, sig)
	maxTokenCost := TokenCost(maxCost, rate.ExchangeRate)
	return &Result{
		Paymaster:                     q.Paymaster,
		PaymasterVerificationGasLimit: (*hexutil.Big)(op.PaymasterVerificationGasLimit),
		PaymasterPostOpGasLimit:       (*hexutil.Big)(op.PaymasterPostOpGasLimit),
		PaymasterData:                 data,
		PaymasterAndData:              PaymasterAndData(q.Paymaster, op.PaymasterVerificationGasLimit, op.PaymasterPostOpGasLimit, data),
		Token:                         t.Address,
		Symbol:                        t.Symbol,
		ExchangeRate:                  (*hexutil.Big)(rate.ExchangeRate),
		MarketRate:                    rate.MarketRate,
		QuotedRate:                    rate.QuotedRate,
		SpreadBps:                     rate.SpreadBps,
		ValidAfter:                    validAfter,
		ValidUntil:                    validUntil,
		MaxCostWei:                    (*hexutil.Big)(maxCost),
		MaxTokenCost:                  (*hexutil.Big)(maxTokenCost),
		MaxTokenCostFormatted:         decimal.NewFromBigInt(maxTokenCost, -int32(t.Decimals)),
		QuoteHash:                     hash,
		Signer:                        q.Signer(),
	}, nil
}
