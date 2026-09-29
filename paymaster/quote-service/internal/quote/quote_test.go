package quote

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"math/big"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"trovo-paymaster-quote-service/internal/chain"
	"trovo-paymaster-quote-service/internal/config"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/shopspring/decimal"
)

type fixture struct {
	ChainID               string `json:"chainId"`
	Paymaster             string `json:"paymaster"`
	QuoteSignerPrivateKey string `json:"quoteSignerPrivateKey"`
	Cases                 []struct {
		Name                          string `json:"name"`
		Sender                        string `json:"sender"`
		Nonce                         string `json:"nonce"`
		InitCode                      string `json:"initCode"`
		CallData                      string `json:"callData"`
		VerificationGasLimit          string `json:"verificationGasLimit"`
		CallGasLimit                  string `json:"callGasLimit"`
		PreVerificationGas            string `json:"preVerificationGas"`
		MaxFeePerGas                  string `json:"maxFeePerGas"`
		MaxPriorityFeePerGas          string `json:"maxPriorityFeePerGas"`
		PaymasterVerificationGasLimit string `json:"paymasterVerificationGasLimit"`
		PaymasterPostOpGasLimit       string `json:"paymasterPostOpGasLimit"`
		Token                         string `json:"token"`
		ValidUntil                    uint64 `json:"validUntil"`
		ValidAfter                    uint64 `json:"validAfter"`
		ExchangeRate                  string `json:"exchangeRate"`
		Hash                          string `json:"hash"`
		Signature                     string `json:"signature"`
	} `json:"cases"`
}

func bi(s string) *big.Int {
	v, ok := new(big.Int).SetString(s, 10)
	if !ok {
		panic(s)
	}
	return v
}

// TestHashMatchesContract checks Hash and the eth_sign signature against
// values produced by TrovoTokenPaymaster.getHash and ethers' signMessage
// (paymaster/contracts/scripts/fixtures.js).
func TestHashMatchesContract(t *testing.T) {
	raw, err := os.ReadFile("testdata/hash_fixture.json")
	if err != nil {
		t.Fatal(err)
	}
	var f fixture
	if err := json.Unmarshal(raw, &f); err != nil {
		t.Fatal(err)
	}
	key, err := crypto.HexToECDSA(strings.TrimPrefix(f.QuoteSignerPrivateKey, "0x"))
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range f.Cases {
		t.Run(c.Name, func(t *testing.T) {
			op := UserOp{
				Sender: common.HexToAddress(c.Sender), Nonce: bi(c.Nonce),
				InitCode: hexutil.MustDecode(c.InitCode), CallData: hexutil.MustDecode(c.CallData),
				VerificationGasLimit: bi(c.VerificationGasLimit), CallGasLimit: bi(c.CallGasLimit),
				PreVerificationGas: bi(c.PreVerificationGas), MaxFeePerGas: bi(c.MaxFeePerGas),
				MaxPriorityFeePerGas:          bi(c.MaxPriorityFeePerGas),
				PaymasterVerificationGasLimit: bi(c.PaymasterVerificationGasLimit),
				PaymasterPostOpGasLimit:       bi(c.PaymasterPostOpGasLimit),
			}
			h := Hash(op, bi(f.ChainID), common.HexToAddress(f.Paymaster), common.HexToAddress(c.Token), c.ValidUntil, c.ValidAfter, bi(c.ExchangeRate))
			if h.Hex() != c.Hash {
				t.Fatalf("hash %s, contract says %s", h.Hex(), c.Hash)
			}
			sig, err := crypto.Sign(accounts.TextHash(h.Bytes()), key)
			if err != nil {
				t.Fatal(err)
			}
			sig[64] += 27
			if hexutil.Encode(sig) != c.Signature {
				t.Fatalf("signature %s, ethers says %s", hexutil.Encode(sig), c.Signature)
			}
		})
	}
}

func TestPaymasterDataLayout(t *testing.T) {
	token := common.HexToAddress("0x1000000000000000000000000000000000000001")
	pm := common.HexToAddress("0x2000000000000000000000000000000000000002")
	sig := bytes.Repeat([]byte{0xab}, 65)
	data := PaymasterData(token, 0x010203040506, 0x0a0b0c0d0e0f, big.NewInt(3000000000), sig)
	full := PaymasterAndData(pm, big.NewInt(150000), big.NewInt(80000), data)
	if len(full) != 181 {
		t.Fatalf("paymasterAndData is %d bytes, contract expects 181", len(full))
	}
	check := func(name string, from, to int, want []byte) {
		if !bytes.Equal(full[from:to], want) {
			t.Errorf("%s [%d:%d] = %x, want %x", name, from, to, full[from:to], want)
		}
	}
	check("paymaster", 0, 20, pm.Bytes())
	check("verification gas", 20, 36, common.LeftPadBytes(big.NewInt(150000).Bytes(), 16))
	check("postOp gas", 36, 52, common.LeftPadBytes(big.NewInt(80000).Bytes(), 16))
	check("token", 52, 72, token.Bytes())
	check("validUntil", 72, 78, []byte{1, 2, 3, 4, 5, 6})
	check("validAfter", 78, 84, []byte{0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f})
	check("rate", 84, 116, common.LeftPadBytes(big.NewInt(3000000000).Bytes(), 32))
	check("signature", 116, 181, sig)
}

func TestTokenCostRoundsUp(t *testing.T) {
	rate := big.NewInt(3000_000000)
	for _, c := range []struct{ wei, want int64 }{{0, 0}, {1, 1}, {1e18, 3000_000000}, {333333333333, 1000}, {333333333334, 1001}} {
		if got := TokenCost(big.NewInt(c.wei), rate); got.Int64() != c.want {
			t.Errorf("TokenCost(%d) = %s, want %d", c.wei, got, c.want)
		}
	}
}

type fakeRates map[string]decimal.Decimal

func (f fakeRates) Convert(route []string) (decimal.Decimal, time.Time, error) {
	v := decimal.NewFromInt(1)
	for _, leg := range route {
		r, ok := f[leg]
		if !ok {
			return decimal.Zero, time.Time{}, errors.New("rate " + leg + " unavailable")
		}
		v = v.Mul(r)
	}
	return v, time.Unix(1700000000, 0), nil
}

type fakeChain struct {
	bounds map[common.Address]chain.TokenBounds
	debt   bool
}

func (f fakeChain) TokenBounds(_ context.Context, t common.Address) (chain.TokenBounds, error) {
	return f.bounds[t], nil
}
func (f fakeChain) HasDebt(context.Context, common.Address) (bool, error) { return f.debt, nil }

var (
	usdc = common.HexToAddress("0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913")
	cngn = common.HexToAddress("0x1000000000000000000000000000000000000001")
)

func spread(v int64) *int64 { return &v }

func newQuoter(t *testing.T) (*Quoter, *fakeChain) {
	key, _ := crypto.GenerateKey()
	fc := &fakeChain{bounds: map[common.Address]chain.TokenBounds{
		usdc: {Enabled: true, MinRate: big.NewInt(1000e6), MaxRate: big.NewInt(10000e6)},
		cngn: {Enabled: true, MinRate: bi("1000000000000"), MaxRate: bi("20000000000000")},
	}}
	dep := big.NewInt(1e18)
	return &Quoter{
		ChainID: big.NewInt(8453), Paymaster: common.HexToAddress("0x2000000000000000000000000000000000000002"), Key: key,
		Tokens: []config.GasToken{
			{Symbol: "USDC", Address: usdc, Decimals: 6, Route: []string{"ETH/USD", "USD/USDC"}},
			{Symbol: "cNGN", Address: cngn, Decimals: 6, Route: []string{"ETH/USD", "USD/NGN", "NGN/CNGN"}, SpreadBps: spread(150)},
		},
		Spread: func(g config.GasToken) int64 {
			if g.SpreadBps != nil {
				return *g.SpreadBps
			}
			return 100
		},
		Rates: fakeRates{
			"ETH/USD": decimal.RequireFromString("3000"), "USD/USDC": decimal.NewFromInt(1),
			"USD/NGN": decimal.RequireFromString("1550.25"), "NGN/CNGN": decimal.NewFromInt(1),
		},
		Chain: fc, Deposit: func() (*big.Int, bool) { return dep, true },
		Validity: 10 * time.Minute, MaxValidity: 24 * time.Hour, ClockSkew: time.Minute,
		PaymasterVerificationGasLimit: 150000, PaymasterPostOpGasLimit: 80000,
		Now: func() time.Time { return time.Unix(1790000000, 0) },
	}, fc
}

func testOp() UserOp {
	return UserOp{
		Sender: common.HexToAddress("0x5a6b0c1c6f1c0f4b1d7aa2d1c3e0bee8a9f2d4c1"), Nonce: big.NewInt(3),
		CallData: []byte{1, 2, 3}, VerificationGasLimit: big.NewInt(300000), CallGasLimit: big.NewInt(200000),
		PreVerificationGas: big.NewInt(60000), MaxFeePerGas: big.NewInt(2e9), MaxPriorityFeePerGas: big.NewInt(1e6),
	}
}

func TestRateAppliesSpread(t *testing.T) {
	q, _ := newQuoter(t)
	usdcTok, _ := q.Token("USDC")
	r, err := q.Rate(usdcTok)
	if err != nil {
		t.Fatal(err)
	}
	// 3000 USD/ETH + 1% default spread = 3030 USDC per ETH
	if r.ExchangeRate.String() != "3030000000" || r.SpreadBps != 100 || !r.MarketRate.Equal(decimal.NewFromInt(3000)) {
		t.Fatalf("USDC rate %+v", r)
	}
	cngnTok, _ := q.Token(cngn.Hex())
	r, err = q.Rate(cngnTok)
	if err != nil {
		t.Fatal(err)
	}
	// 3000 * 1550.25 = 4,650,750 cNGN per ETH, + 1.5% = 4,720,511.25
	if r.ExchangeRate.String() != "4720511250000" || r.SpreadBps != 150 {
		t.Fatalf("cNGN rate %s spread %d", r.ExchangeRate, r.SpreadBps)
	}
}

func TestQuoteIsSignedAndBound(t *testing.T) {
	q, _ := newQuoter(t)
	op := testOp()
	res, err := q.Quote(context.Background(), Request{Op: op, Token: "cNGN", Validity: 24 * time.Hour})
	if err != nil {
		t.Fatal(err)
	}
	if res.ValidUntil != 1790000000+86400 || res.ValidAfter != 1790000000-60 {
		t.Fatalf("validity window %d..%d", res.ValidAfter, res.ValidUntil)
	}
	if len(res.PaymasterAndData) != 181 || !bytes.Equal(res.PaymasterAndData[52:], res.PaymasterData) {
		t.Fatal("paymasterAndData must be the static fields followed by paymasterData")
	}
	op.PaymasterVerificationGasLimit, op.PaymasterPostOpGasLimit = big.NewInt(150000), big.NewInt(80000)
	want := Hash(op, q.ChainID, q.Paymaster, cngn, res.ValidUntil, res.ValidAfter, res.ExchangeRate.ToInt())
	if res.QuoteHash != want {
		t.Fatal("quote hash does not cover the operation")
	}
	sig := append([]byte(nil), res.PaymasterData[64:]...)
	sig[64] -= 27
	pub, err := crypto.SigToPub(accounts.TextHash(want.Bytes()), sig)
	if err != nil || crypto.PubkeyToAddress(*pub) != q.Signer() {
		t.Fatal("signature does not recover to the quote signer")
	}
	// max cost = (300000+200000+150000+80000+60000) gas * 2 gwei
	if res.MaxCostWei.ToInt().String() != "1580000000000000" {
		t.Fatalf("max cost %s", res.MaxCostWei.ToInt())
	}
	if res.MaxTokenCost.ToInt().Cmp(TokenCost(res.MaxCostWei.ToInt(), res.ExchangeRate.ToInt())) != 0 {
		t.Fatal("max token cost must use the contract's rounding")
	}
}

func TestQuoteRefusals(t *testing.T) {
	cases := map[string]struct {
		mutate func(q *Quoter, fc *fakeChain, r *Request)
		status int
		code   string
	}{
		"unknown token":     {func(_ *Quoter, _ *fakeChain, r *Request) { r.Token = "DAI" }, http.StatusBadRequest, "unknown-token"},
		"validity too long": {func(_ *Quoter, _ *fakeChain, r *Request) { r.Validity = 48 * time.Hour }, http.StatusBadRequest, "invalid-validity"},
		"bad op":            {func(_ *Quoter, _ *fakeChain, r *Request) { r.Op.MaxPriorityFeePerGas = big.NewInt(3e9) }, http.StatusBadRequest, "invalid-user-operation"},
		"rate missing": {func(q *Quoter, _ *fakeChain, _ *Request) {
			delete(q.Rates.(fakeRates), "USD/NGN")
		}, http.StatusServiceUnavailable, "rate-unavailable"},
		"disabled on-chain": {func(_ *Quoter, fc *fakeChain, _ *Request) {
			fc.bounds[cngn] = chain.TokenBounds{MinRate: big.NewInt(1), MaxRate: big.NewInt(2)}
		}, http.StatusServiceUnavailable, "token-disabled"},
		"above on-chain max": {func(_ *Quoter, fc *fakeChain, _ *Request) {
			fc.bounds[cngn] = chain.TokenBounds{Enabled: true, MinRate: big.NewInt(1), MaxRate: big.NewInt(2)}
		}, http.StatusServiceUnavailable, "rate-out-of-bounds"},
		"wallet in debt": {func(_ *Quoter, fc *fakeChain, _ *Request) { fc.debt = true }, http.StatusConflict, "outstanding-debt"},
		"deposit too low": {func(q *Quoter, _ *fakeChain, _ *Request) {
			q.Deposit = func() (*big.Int, bool) { return big.NewInt(1), true }
		}, http.StatusServiceUnavailable, "paymaster-deposit-low"},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			q, fc := newQuoter(t)
			req := Request{Op: testOp(), Token: "cNGN"}
			c.mutate(q, fc, &req)
			_, err := q.Quote(context.Background(), req)
			var qe *Error
			if !errors.As(err, &qe) || qe.Status != c.status || qe.Code != c.code {
				t.Fatalf("got %v, want %d %s", err, c.status, c.code)
			}
		})
	}
}
