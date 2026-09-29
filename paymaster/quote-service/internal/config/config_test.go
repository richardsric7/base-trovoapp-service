package config

import (
	"strings"
	"testing"
	"time"
)

func validValues() Values {
	return Values{
		"RPC_URL":                  "https://mainnet.base.org",
		"CHAIN_ID":                 "8453",
		"PAYMASTER_ADDRESS":        "0x2000000000000000000000000000000000000002",
		"QUOTE_SIGNER_PRIVATE_KEY": "0x59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d",
		"API_KEYS":                 "backend-key-0123456789, ops-key-0123456789abc",
		"DEFAULT_SPREAD_BPS":       "100",
		"GAS_TOKENS": `[
			{"symbol":"USDC","address":"0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913","decimals":6,"route":["ETH/USD","USD/USDC"]},
			{"symbol":"cNGN","address":"0x1000000000000000000000000000000000000001","decimals":6,"route":["ETH/USD","USD/NGN","NGN/CNGN"],"spreadBps":150}
		]`,
		"RATE_PAIRS": `{
			"ETH/USD": {"sources":[{"type":"fixed","value":"3000"}]},
			"USDC/USD": {"sources":[{"type":"fixed","value":"1"}]},
			"USD/NGN": {"sources":[{"type":"fixed","value":"1550"}],"minSources":1},
			"NGN/CNGN": {"sources":[{"type":"fixed","value":"1"}]}
		}`,
		"QUOTE_MAX_VALIDITY": "24h",
	}
}

func TestParseValid(t *testing.T) {
	c, err := Parse(validValues())
	if err != nil {
		t.Fatal(err)
	}
	if len(c.APIKeys) != 2 || c.ChainID.Int64() != 8453 || c.QuoteMaxValidity != 24*time.Hour || c.HTTPAddr != ":8090" {
		t.Fatalf("unexpected config %+v", c)
	}
	if c.Spread(c.GasTokens[0]) != 100 || c.Spread(c.GasTokens[1]) != 150 {
		t.Fatal("token spread must override the default")
	}
	if c.DepositLowWatermark.String() != "50000000000000000" {
		t.Fatalf("default watermark 0.05 ETH, got %s wei", c.DepositLowWatermark)
	}
}

func TestParseRejects(t *testing.T) {
	cases := map[string]func(v Values){
		"missing signer": func(v Values) { delete(v, "QUOTE_SIGNER_PRIVATE_KEY") },
		"bad signer":     func(v Values) { v["QUOTE_SIGNER_PRIVATE_KEY"] = "0x1234" },
		"short api key":  func(v Values) { v["API_KEYS"] = "short" },
		"bad address":    func(v Values) { v["PAYMASTER_ADDRESS"] = "0x12" },
		"route not at ETH": func(v Values) {
			v["GAS_TOKENS"] = `[{"symbol":"X","address":"0x1000000000000000000000000000000000000001","decimals":6,"route":["USD/NGN"]}]`
		},
		"broken route": func(v Values) {
			v["GAS_TOKENS"] = `[{"symbol":"X","address":"0x1000000000000000000000000000000000000001","decimals":6,"route":["ETH/USD","NGN/CNGN"]}]`
		},
		"unknown pair": func(v Values) {
			v["GAS_TOKENS"] = `[{"symbol":"X","address":"0x1000000000000000000000000000000000000001","decimals":6,"route":["ETH/EUR"]}]`
		},
		"spread too high":      func(v Values) { v["DEFAULT_SPREAD_BPS"] = "9000" },
		"validity too long":    func(v Values) { v["QUOTE_MAX_VALIDITY"] = "200h" },
		"pair without sources": func(v Values) { v["RATE_PAIRS"] = `{"ETH/USD":{"sources":[]}}` },
	}
	for name, mutate := range cases {
		t.Run(name, func(t *testing.T) {
			v := validValues()
			mutate(v)
			if _, err := Parse(v); err == nil {
				t.Fatal("expected an error")
			}
		})
	}
}

func TestInverseLegResolves(t *testing.T) {
	v := validValues()
	// USD/USDC is configured as USDC/USD; the route may use either direction
	if _, err := Parse(v); err != nil {
		t.Fatal(err)
	}
	name, inverted, ok := LookupPair(map[string]int{"USDC/USD": 1}, "USD", "USDC")
	if !ok || !inverted || name != "USDC/USD" {
		t.Fatalf("got %s %v %v", name, inverted, ok)
	}
	if _, _, err := SplitPair("ETH"); err == nil || !strings.Contains(err.Error(), "BASE/QUOTE") {
		t.Fatal("expected a format error")
	}
}
