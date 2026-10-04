package users

import (
	"math/big"
	"strings"
	"testing"

	"trovo-wallet-api/internal/aa"

	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
)

var (
	mintIssuing = common.HexToAddress("0x2000000000000000000000000000000000000002")
	mintDist    = common.HexToAddress("0x3000000000000000000000000000000000000003")
	mintToken   = common.HexToAddress("0x1000000000000000000000000000000000000001")
	mintFee     = common.HexToAddress("0x4000000000000000000000000000000000000004")
	mintBook    = common.HexToAddress("0x5000000000000000000000000000000000000005")
	mintFunds   = common.HexToAddress("0x6000000000000000000000000000000000000006")
	mintCNGN    = common.HexToAddress("0x7000000000000000000000000000000000000007")
)

func testPlan(t *testing.T, supply, fee, sale string, decimals uint8) tokenizationMintPlan {
	t.Helper()
	plan, err := newTokenizationMintPlan(big.NewInt(8453), mintIssuing, mintDist, mintToken, mintFee, decimals, decimal.RequireFromString(supply), decimal.RequireFromString(fee), decimal.RequireFromString(sale))
	if err != nil {
		t.Fatal(err)
	}
	price, err := offerPrice(decimal.RequireFromString("1500"), decimals, 6)
	if err != nil {
		t.Fatal(err)
	}
	plan.Book, plan.Proceeds = mintBook, mintFunds
	plan.PaymentCodes, plan.PaymentTokens, plan.Prices = []string{"CNGN"}, []common.Address{mintCNGN}, []aa.OfferPrice{price}
	return plan
}

func decodeMint(t *testing.T, c aa.Call) (common.Address, *big.Int) {
	t.Helper()
	if c.To != mintToken {
		t.Fatalf("mint must call the token contract, got %s", c.To.Hex())
	}
	args, err := tokenizedAssetTokenABI.Methods["mint"].Inputs.Unpack(c.Data[4:])
	if err != nil {
		t.Fatal(err)
	}
	return args[0].(common.Address), args[1].(*big.Int)
}

func TestTokenizationMintPlanMintsAndOffersTheSale(t *testing.T) {
	plan := testPlan(t, "1000000", "2500.5", "600000", 6)
	if plan.ForSale.String() != "600000000000" || plan.DistributionAmount.String() != "397499500000" || plan.FeeAmount.String() != "2500500000" {
		t.Fatalf("base-unit amounts: for sale %s distribution %s fee %s", plan.ForSale, plan.DistributionAmount, plan.FeeAmount)
	}
	deploy := aa.Call{To: common.HexToAddress("0x9999999999999999999999999999999999999999"), Value: big.NewInt(0), Data: []byte{1}}
	plan.DeployDistribution = &deploy
	calls, err := plan.Calls()
	if err != nil {
		t.Fatal(err)
	}
	if len(calls) != 6 || calls[0].To != deploy.To {
		t.Fatalf("want deploy, three mints, approve and offer; got %d calls", len(calls))
	}
	for i, want := range []struct {
		to     common.Address
		amount *big.Int
	}{{mintIssuing, plan.ForSale}, {mintDist, plan.DistributionAmount}, {mintFee, plan.FeeAmount}} {
		if to, amt := decodeMint(t, calls[1+i]); to != want.to || amt.Cmp(want.amount) != 0 {
			t.Fatalf("mint %d: %v to %v, want %v to %v", i, amt, to, want.amount, want.to)
		}
	}
	// the issuing Safe approves the book and offers, in its own calls (no
	// module call, so a failure reverts the whole mint)
	if calls[4].To != mintToken || calls[5].To != mintBook {
		t.Fatalf("approve %s, offer %s", calls[4].To.Hex(), calls[5].To.Hex())
	}
	d := plan.Description()
	for _, s := range []string{mintIssuing.Hex(), mintDist.Hex(), mintToken.Hex(), mintFee.Hex(), mintBook.Hex(), mintFunds.Hex(), "CNGN", "600000", "8453"} {
		if !strings.Contains(d, s) {
			t.Fatalf("description %q does not name %s", d, s)
		}
	}
}

func TestTokenizationMintPlanWithoutFeeOrDeploy(t *testing.T) {
	plan := testPlan(t, "100", "0", "100", 18)
	calls, _ := plan.Calls()
	if len(calls) != 3 {
		t.Fatalf("expected mint, approve and offer, got %d calls", len(calls))
	}
	if strings.Contains(plan.Description(), mintFee.Hex()) || strings.Contains(plan.Description(), mintDist.Hex()) {
		t.Fatal("description must not mention mints that will not happen")
	}
}

func TestTokenizationMintPlanRejectsInvalidAmounts(t *testing.T) {
	cases := map[string][3]string{
		"zero supply":            {"0", "0", "0"},
		"fee above supply":       {"10", "11", "1"},
		"negative fee":           {"10", "-1", "1"},
		"more precision than 2":  {"10.123", "0", "1"},
		"nothing for sale":       {"10", "0", "0"},
		"sale above supply-fees": {"10", "2", "9"},
	}
	for name, c := range cases {
		if _, err := newTokenizationMintPlan(big.NewInt(1), mintIssuing, mintDist, mintToken, mintFee, 2, decimal.RequireFromString(c[0]), decimal.RequireFromString(c[1]), decimal.RequireFromString(c[2])); err == nil {
			t.Errorf("%s: expected an error", name)
		}
	}
}

func TestOfferPrice(t *testing.T) {
	// 1500 per token, asset 7 decimals, cNGN 6: 1 token (1e7 base units) costs 1500e6
	p, err := offerPrice(decimal.RequireFromString("1500"), 7, 6)
	if err != nil {
		t.Fatal(err)
	}
	if p.Cost(big.NewInt(10_000_000)).Cmp(big.NewInt(1_500_000_000)) != 0 {
		t.Fatalf("1 token costs %v", p.Cost(big.NewInt(10_000_000)))
	}
	// fractional price, 18-decimal payment token: 2.75 per token
	p, err = offerPrice(decimal.RequireFromString("2.75"), 2, 18)
	if err != nil {
		t.Fatal(err)
	}
	want, _ := new(big.Int).SetString("2750000000000000000", 10)
	if p.Cost(big.NewInt(100)).Cmp(want) != 0 {
		t.Fatalf("1 token costs %v", p.Cost(big.NewInt(100)))
	}
	// what a budget buys never costs more than the budget
	budget := big.NewInt(1_000_000_007)
	if amt := p.AmountFor(budget); p.Cost(amt).Cmp(budget) > 0 {
		t.Fatalf("%v buys %v costing %v", budget, amt, p.Cost(amt))
	}
	if _, err := offerPrice(decimal.Zero, 7, 6); err == nil {
		t.Fatal("zero price")
	}
}
