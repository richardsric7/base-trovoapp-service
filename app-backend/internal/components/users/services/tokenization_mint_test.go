package users

import (
	"encoding/base64"
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
)

var (
	mintSafe  = common.HexToAddress("0x2000000000000000000000000000000000000002")
	mintToken = common.HexToAddress("0x1000000000000000000000000000000000000001")
	mintFee   = common.HexToAddress("0x4000000000000000000000000000000000000004")
)

func TestTokenizationMintPlanSplitsSupplyBetweenTreasuryAndFeeWallet(t *testing.T) {
	plan, err := newTokenizationMintPlan(big.NewInt(8453), mintSafe, mintToken, mintFee, 6, decimal.RequireFromString("1000000"), decimal.RequireFromString("2500.5"))
	if err != nil {
		t.Fatal(err)
	}
	if plan.Treasury != mintSafe {
		t.Fatalf("unsold supply must be minted to the issuing Safe, got %s", plan.Treasury.Hex())
	}
	if plan.TreasuryAmount.String() != "997499500000" || plan.FeeAmount.String() != "2500500000" {
		t.Fatalf("unexpected base-unit amounts: treasury %s fee %s", plan.TreasuryAmount, plan.FeeAmount)
	}

	calls, err := plan.Calls()
	if err != nil {
		t.Fatal(err)
	}
	if len(calls) != 2 {
		t.Fatalf("expected treasury and fee mint calls, got %d", len(calls))
	}
	for i, want := range []struct {
		to     common.Address
		amount *big.Int
	}{{mintSafe, plan.TreasuryAmount}, {mintFee, plan.FeeAmount}} {
		c := calls[i]
		if c.To != mintToken {
			t.Fatalf("call %d must target the token contract, got %s", i, c.To.Hex())
		}
		method, err := tokenizedAssetTokenABI.MethodById(c.Data[:4])
		if err != nil || method.Name != "mint" {
			t.Fatalf("call %d is not mint(): %v", i, err)
		}
		args, err := method.Inputs.Unpack(c.Data[4:])
		if err != nil {
			t.Fatal(err)
		}
		if args[0].(common.Address) != want.to || args[1].(*big.Int).Cmp(want.amount) != 0 {
			t.Fatalf("call %d mints %v to %v, want %v to %v", i, args[1], args[0], want.amount, want.to)
		}
	}

	decoded, err := base64.StdEncoding.DecodeString(plan.ApprovalPayload())
	if err != nil || string(decoded) != plan.Description() {
		t.Fatal("approval payload must be the base64 of the description approvers sign")
	}
	for _, s := range []string{mintSafe.Hex(), mintToken.Hex(), mintFee.Hex(), "8453"} {
		if !strings.Contains(plan.Description(), s) {
			t.Fatalf("description %q does not name %s", plan.Description(), s)
		}
	}
}

func TestTokenizationMintPlanWithoutFeeMintsOnce(t *testing.T) {
	plan, err := newTokenizationMintPlan(big.NewInt(8453), mintSafe, mintToken, mintFee, 18, decimal.RequireFromString("100"), decimal.Zero)
	if err != nil {
		t.Fatal(err)
	}
	calls, _ := plan.Calls()
	if len(calls) != 1 {
		t.Fatalf("expected a single treasury mint, got %d calls", len(calls))
	}
	if strings.Contains(plan.Description(), mintFee.Hex()) {
		t.Fatal("description must not mention a fee mint that will not happen")
	}
}

func TestTokenizationMintPlanRejectsInvalidAmounts(t *testing.T) {
	cases := map[string][2]string{
		"zero supply":           {"0", "0"},
		"fee above supply":      {"10", "11"},
		"negative fee":          {"10", "-1"},
		"more precision than 2": {"10.123", "0"},
	}
	for name, c := range cases {
		if _, err := newTokenizationMintPlan(big.NewInt(1), mintSafe, mintToken, mintFee, 2, decimal.RequireFromString(c[0]), decimal.RequireFromString(c[1])); err == nil {
			t.Errorf("%s: expected an error", name)
		}
	}
}
