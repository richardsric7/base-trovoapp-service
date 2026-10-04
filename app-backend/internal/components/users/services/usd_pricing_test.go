package users

import "testing"

func TestUSDAmountIn(t *testing.T) {
	t.Setenv("DOLLAR_ASSET", "USDB:0x0000000000000000000000000000000000000001")
	t.Setenv("NAIRA_ASSET", "")
	for _, c := range []struct {
		code, contract string
		want           string
	}{
		{"USDB", "0x0000000000000000000000000000000000000001", "2.5"},
		{"usdc", "0x0000000000000000000000000000000000000002", "2.5"},
		{"USDT", "0x0000000000000000000000000000000000000003", "2.5"},
	} {
		got, err := usdAmountIn(2.5, c.code, c.contract, nil)
		if err != nil || got.String() != c.want {
			t.Fatalf("%s: got %v, %v; want %s", c.code, got, err, c.want)
		}
	}
	// a $2.50 fee must never become 1/2.5 of a token (the pre-port inversion)
	if got, _ := usdAmountIn(2.5, "USDC", "0x02", nil); got.String() == "0.4" {
		t.Fatal("inverted conversion")
	}
	if _, err := usdAmountIn(2.5, "TROV", "0x0000000000000000000000000000000000000004", nil); err == nil {
		t.Fatal("a non-stablecoin cannot be priced without a DEX")
	}
}

func TestFeeExemptProfile(t *testing.T) {
	t.Setenv("TOKENIZATION_ISSUING_PROFILE", "")
	if !feeExemptProfile("atprofile") || feeExemptProfile("alice") {
		t.Fatal("default exempt profile is atprofile")
	}
	t.Setenv("TOKENIZATION_ISSUING_PROFILE", "issuer")
	if !feeExemptProfile("ISSUER") || feeExemptProfile("atprofile") {
		t.Fatal("the configured issuing profile is exempt")
	}
}
