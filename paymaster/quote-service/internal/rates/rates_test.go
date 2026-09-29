package rates

import (
	"context"
	"encoding/json"
	"math/big"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"trovo-paymaster-quote-service/internal/config"

	"github.com/shopspring/decimal"
)

func d(s string) decimal.Decimal { return decimal.RequireFromString(s) }

func TestAggregateMedianAndOutliers(t *testing.T) {
	res := []SourceResult{
		{Name: "a", Value: d("1500")},
		{Name: "b", Value: d("1520")},
		{Name: "c", Value: d("2400")}, // outlier
		{Name: "d", Error: "timeout"},
	}
	v, err := aggregate(res, 2, 500)
	if err != nil {
		t.Fatal(err)
	}
	if !v.Equal(d("1510")) {
		t.Fatalf("median of the kept sources = %s, want 1510", v)
	}
	if !res[2].Outlier || res[0].Outlier {
		t.Fatal("only the far source should be marked an outlier")
	}
	if _, err := aggregate(res, 3, 500); err == nil {
		t.Fatal("expected an error when too few sources agree")
	}
	if v, _ := aggregate(res, 1, -1); !v.Equal(d("1520")) {
		t.Fatalf("with the deviation check disabled the median is 1520, got %s", v)
	}
	if _, err := aggregate([]SourceResult{{Name: "x", Error: "down"}}, 1, 500); err == nil {
		t.Fatal("expected an error when nothing answered")
	}
}

func TestHTTPJSONSource(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Api-Key") != "s3cret" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.Write([]byte(`{"data":{"tickers":[{"last":"1549.75"},{"last":1551.5}]}}`))
	}))
	defer srv.Close()
	env := Env{HTTP: srv.Client(), Secret: func(k string) (string, bool) { return map[string]string{"FX_KEY": "s3cret"}[k], k == "FX_KEY" }}

	for path, want := range map[string]string{"data.tickers.0.last": "1549.75", "data.tickers.1.last": "1551.5"} {
		def := `{"type":"http-json","url":"` + srv.URL + `","path":"` + path + `","headers":{"X-Api-Key":"vault:FX_KEY"}}`
		_, src, err := Build(json.RawMessage(def), env)
		if err != nil {
			t.Fatal(err)
		}
		v, err := src.Fetch(context.Background())
		if err != nil || !v.Equal(d(want)) {
			t.Fatalf("%s: got %s, %v", path, v, err)
		}
	}

	_, src, _ := Build(json.RawMessage(`{"type":"http-json","url":"`+srv.URL+`","path":"data.missing"}`), env)
	if _, err := src.Fetch(context.Background()); err == nil {
		t.Fatal("expected an error without the API key / for a missing path")
	}
	if _, _, err := Build(json.RawMessage(`{"type":"http-json","url":"`+srv.URL+`","path":"x","headers":{"K":"vault:NOPE"}}`), env); err == nil {
		t.Fatal("expected an error for an unset secret reference")
	}
}

func TestInvertAndMultiplier(t *testing.T) {
	_, src, err := Build(json.RawMessage(`{"type":"fixed","value":"0.0005","invert":true,"multiplier":"1.01"}`), Env{})
	if err != nil {
		t.Fatal(err)
	}
	v, _ := src.Fetch(context.Background())
	if !v.Round(6).Equal(d("2020")) {
		t.Fatalf("1/0.0005 * 1.01 = %s, want 2020", v)
	}
	if _, _, err := Build(json.RawMessage(`{"type":"nope"}`), Env{}); err == nil || !strings.Contains(err.Error(), "unknown source type") {
		t.Fatalf("expected unknown type error, got %v", err)
	}
}

func TestTwapPrice(t *testing.T) {
	// cNGN (token0, 6 decimals) / USDC (token1, 6 decimals) pool at ~1550
	// cNGN per USDC: price of token0 in token1 = 1/1550 -> tick = log(1/1550)/log(1.0001)
	tick := int64(-73460) // 1.0001^-73460 ≈ 0.000645 ≈ 1/1550.3
	window := uint32(1800)
	old := big.NewInt(1_000_000)
	now := new(big.Int).Add(old, big.NewInt(tick*int64(window)))

	// base = USDC (token1): price of USDC in cNGN
	p := twapPrice(old, now, window, false, 6, 6)
	if p.LessThan(d("1549")) || p.GreaterThan(d("1552")) {
		t.Fatalf("USDC in cNGN = %s, want ~1550", p)
	}
	// base = cNGN (token0)
	p0 := twapPrice(old, now, window, true, 6, 6)
	if p0.Mul(p).Sub(d("1")).Abs().GreaterThan(d("0.000001")) {
		t.Fatal("the two directions must be reciprocal")
	}
	// decimals: token0 18 decimals, token1 6 decimals
	p18 := twapPrice(big.NewInt(0), big.NewInt(0), window, true, 18, 6)
	if !p18.Equal(d("1000000000000")) {
		t.Fatalf("tick 0 with 18/6 decimals = %s", p18)
	}
	// negative cumulative deltas round toward negative infinity
	if a, b := twapPrice(big.NewInt(0), big.NewInt(-1), window, true, 0, 0), twapPrice(big.NewInt(0), big.NewInt(0), window, true, 0, 0); !a.LessThan(b) {
		t.Fatal("a -1 delta must round to tick -1")
	}
}

func TestBookConvertStaleAndInverse(t *testing.T) {
	pairs := map[string]config.PairConfig{
		"ETH/USD": {Sources: []json.RawMessage{json.RawMessage(`{"type":"fixed","value":"3000"}`)}},
		"NGN/USD": {Sources: []json.RawMessage{json.RawMessage(`{"type":"fixed","value":"0.0005"}`)}},
	}
	b, err := NewBook(pairs, Env{}, 5*time.Minute, time.Second)
	if err != nil {
		t.Fatal(err)
	}
	if _, _, err := b.Convert([]string{"ETH/USD"}); err == nil {
		t.Fatal("rates must be unavailable before the first refresh")
	}
	b.Refresh(context.Background())
	v, _, err := b.Convert([]string{"ETH/USD", "USD/NGN"}) // USD/NGN is the inverse of NGN/USD
	if err != nil {
		t.Fatal(err)
	}
	if !v.Round(6).Equal(d("6000000")) {
		t.Fatalf("ETH in NGN = %s, want 6000000", v)
	}
	b.now = func() time.Time { return time.Now().Add(6 * time.Minute) }
	if _, _, err := b.Convert([]string{"ETH/USD"}); err == nil || !strings.Contains(err.Error(), "stale") {
		t.Fatalf("expected a stale-rate error, got %v", err)
	}
}
