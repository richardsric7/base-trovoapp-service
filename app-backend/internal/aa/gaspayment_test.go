package aa

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func paymasterLog(t *testing.T, addr common.Address, event string, topics []common.Hash, args ...interface{}) *types.Log {
	t.Helper()
	ev := paymasterABI.Events[event]
	data, err := ev.Inputs.NonIndexed().Pack(args...)
	if err != nil {
		t.Fatal(err)
	}
	return &types.Log{Address: addr, Topics: append([]common.Hash{ev.ID}, topics...), Data: data}
}

// The paymaster's charge for an operation, and the debts a wallet settled,
// are read from the mined transaction's logs.
func TestParseGasPaymentAndDebtSettlements(t *testing.T) {
	pm := common.HexToAddress("0x00000000000000000000000000000000000000aa")
	other := common.HexToAddress("0x00000000000000000000000000000000000000bb")
	wallet := common.HexToAddress("0x0000000000000000000000000000000000000001")
	stranger := common.HexToAddress("0x0000000000000000000000000000000000000002")
	usdc := common.HexToAddress("0x0000000000000000000000000000000000000c0c")
	opA, opB := common.HexToHash("0xa1"), common.HexToHash("0xb2")
	addr := func(a common.Address) common.Hash { return common.BytesToHash(a.Bytes()) }

	logs := []*types.Log{
		// same event from another contract: ignored
		paymasterLog(t, other, "GasPaidInToken", []common.Hash{addr(wallet), addr(usdc), opA}, big.NewInt(1), big.NewInt(999), big.NewInt(1)),
		paymasterLog(t, pm, "GasPaidInToken", []common.Hash{addr(wallet), addr(usdc), opA}, big.NewInt(3000), big.NewInt(12345), big.NewInt(4)),
		paymasterLog(t, pm, "GasChargeFailed", []common.Hash{addr(wallet), addr(usdc), opB}, big.NewInt(777)),
		paymasterLog(t, pm, "DebtSettled", []common.Hash{addr(wallet), addr(usdc)}, big.NewInt(500), false),
		// forgiven debts were not paid by the wallet
		paymasterLog(t, pm, "DebtSettled", []common.Hash{addr(wallet), addr(usdc)}, big.NewInt(600), true),
		// another wallet's settlement
		paymasterLog(t, pm, "DebtSettled", []common.Hash{addr(stranger), addr(usdc)}, big.NewInt(700), false),
	}

	g := ParseGasPayment(logs, pm, opA)
	if g == nil || g.Failed || g.Token != usdc || g.TokenCost.Int64() != 12345 {
		t.Fatalf("charge for opA: %+v", g)
	}
	g = ParseGasPayment(logs, pm, opB)
	if g == nil || !g.Failed || g.TokenCost.Int64() != 777 {
		t.Fatalf("failed charge for opB: %+v", g)
	}
	if g := ParseGasPayment(logs, pm, common.HexToHash("0xc3")); g != nil {
		t.Fatalf("an operation the paymaster did not pay for: %+v", g)
	}

	s := ParseDebtSettlements(logs, pm, wallet)
	if len(s) != 1 || s[0].Token != usdc || s[0].Amount.Int64() != 500 {
		t.Fatalf("settlements: %+v", s)
	}
}

func TestSettleDebtCalls(t *testing.T) {
	pm := common.HexToAddress("0x00000000000000000000000000000000000000aa")
	wallet := common.HexToAddress("0x0000000000000000000000000000000000000001")
	usdc := common.HexToAddress("0x0000000000000000000000000000000000000c0c")
	calls := SettleDebtCalls(pm, wallet, usdc)
	if len(calls) != 2 || calls[0].To != usdc || calls[1].To != pm {
		t.Fatalf("calls: %+v", calls)
	}
	args, err := paymasterABI.Methods["settleDebt"].Inputs.Unpack(calls[1].Data[4:])
	if err != nil || args[0].(common.Address) != wallet || args[1].(common.Address) != usdc {
		t.Fatalf("settleDebt args: %v %v", args, err)
	}
	approve, err := contractsABI.Methods["approve"].Inputs.Unpack(calls[0].Data[4:])
	if err != nil || approve[0].(common.Address) != pm || approve[1].(*big.Int).Cmp(maxUint256) != 0 {
		t.Fatalf("approve args: %v %v", approve, err)
	}
}
