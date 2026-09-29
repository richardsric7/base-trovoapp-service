package users

import (
	"math/big"
	"os"
	"testing"

	"trovo-wallet-api/internal/aa"
	"trovo-wallet-api/internal/basetxn"
	paymentModels "trovo-wallet-api/internal/components/payments/models"

	"github.com/ethereum/go-ethereum/common"
)

func TestBaseUnits(t *testing.T) {
	for _, c := range []struct {
		amount string
		dec    uint8
		want   string
	}{{"1", 6, "1000000"}, {"0.1234567", 6, "123456"}, {"2.5", 18, "2500000000000000000"}, {"0", 6, "0"}} {
		got, err := baseUnits(c.amount, c.dec)
		if err != nil || got.String() != c.want {
			t.Errorf("baseUnits(%s, %d) = %v, %v; want %s", c.amount, c.dec, got, err, c.want)
		}
	}
	if _, err := baseUnits("-1", 6); err == nil {
		t.Error("negative amounts must be rejected")
	}
	if _, err := baseUnits("abc", 6); err == nil {
		t.Error("non-numbers must be rejected")
	}
}

func TestTransferCall(t *testing.T) {
	to := "0x1000000000000000000000000000000000000001"
	native := transferCall(basetxn.NativeAsset{}, to, big.NewInt(7))
	if native.To != common.HexToAddress(to) || native.Value.Int64() != 7 || len(native.Data) != 0 {
		t.Fatal("a native transfer sends value to the destination")
	}
	token := "0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913"
	erc := transferCall(basetxn.CreditAsset{Code: "USDC", Issuer: token}, to, big.NewInt(7))
	want := aa.ERC20Transfer(common.HexToAddress(token), common.HexToAddress(to), big.NewInt(7))
	if erc.To != want.To || string(erc.Data) != string(want.Data) || erc.Value.Sign() != 0 {
		t.Fatal("a token transfer calls transfer() on the token contract")
	}
}

func TestCheckRequiredMemo(t *testing.T) {
	exchange := "0X2000000000000000000000000000000000000002"
	os.Setenv("WALLETS_REQUIRE_16_BYTE_MEMO", exchange)
	defer os.Unsetenv("WALLETS_REQUIRE_16_BYTE_MEMO")
	p := &paymentModels.PaymentInfo{Destination: "0x2000000000000000000000000000000000000002", Memo: "short"}
	if checkRequiredMemo(p) == nil {
		t.Fatal("a 16-character memo is required for this exchange")
	}
	p.Memo = "1234567890abcdef"
	if err := checkRequiredMemo(p); err != nil {
		t.Fatal(err)
	}
	u := &paymentModels.PaymentInfo{Destination: "Alice"}
	if checkRequiredMemo(u) != nil || u.Destination != "alice" {
		t.Fatal("usernames are lower-cased and need no memo")
	}
}

func TestOperationMessages(t *testing.T) {
	p := &aa.Prepared{Activation: true}
	if msgs := operationMessages(p); len(msgs) != 1 {
		t.Fatalf("activation message expected, got %v", msgs)
	}
}
