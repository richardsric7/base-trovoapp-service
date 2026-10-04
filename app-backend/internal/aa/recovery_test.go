package aa

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func TestReplaceOwner(t *testing.T) {
	a, b, c := common.HexToAddress("0xa"), common.HexToAddress("0xb"), common.HexToAddress("0xc")
	got := ReplaceOwner([]common.Address{a, b}, a, c)
	if len(got) != 2 || got[0] != c || got[1] != b {
		t.Fatalf("got %v", got)
	}
	// the input is not modified, and a missing owner changes nothing
	owners := []common.Address{a, b}
	if got := ReplaceOwner(owners, common.HexToAddress("0xd"), c); got[0] != a || got[1] != b || owners[0] != a {
		t.Fatalf("got %v", got)
	}
}

func TestRecoveryCalls(t *testing.T) {
	wallet, module, guardian := common.HexToAddress("0x1001"), common.HexToAddress("0x2002"), common.HexToAddress("0x3003")
	calls := EnableRecoveryCalls(wallet, module, guardian, false)
	if len(calls) != 2 || calls[0].To != wallet || calls[1].To != module {
		t.Fatalf("enable calls %+v", calls)
	}
	if m, _ := recoveryABI.MethodById(calls[0].Data[:4]); m == nil || m.Name != "enableModule" {
		t.Fatal("first call must enable the module")
	}
	if len(EnableRecoveryCalls(wallet, module, guardian, true)) != 1 {
		t.Fatal("an enabled module is not enabled again")
	}
	args, _ := recoveryABI.Methods["revokeGuardianWithThreshold"].Inputs.Unpack(RevokeGuardianCall(module, guardian).Data[4:])
	if args[0].(common.Address) != sentinelGuardian || args[1].(common.Address) != guardian || args[2].(*big.Int).Sign() != 0 {
		t.Fatalf("revoke args %v", args)
	}
	args, _ = recoveryABI.Methods["confirmRecovery"].Inputs.Unpack(ConfirmRecoveryCall(module, wallet, []common.Address{guardian}, 1, big.NewInt(4)).Data[4:])
	if args[0].(common.Address) != wallet || args[3].(*big.Int).Int64() != 4 || !args[4].(bool) {
		t.Fatalf("confirm args %v", args)
	}
}

func TestParseRecoveryStarted(t *testing.T) {
	module, wallet, owner := common.HexToAddress("0x2002"), common.HexToAddress("0x1001"), common.HexToAddress("0x4004")
	ev := recoveryABI.Events["RecoveryExecuted"]
	data, err := ev.Inputs.NonIndexed().Pack([]common.Address{owner}, big.NewInt(1), big.NewInt(0), uint64(1700000000), big.NewInt(1))
	if err != nil {
		t.Fatal(err)
	}
	l := types.Log{Address: module, Topics: []common.Hash{ev.ID, common.BytesToHash(wallet.Bytes()), {}}, Data: data, TxHash: common.HexToHash("0xabc")}
	other := l
	other.Address = common.HexToAddress("0x9")
	got := ParseRecoveryStarted([]types.Log{l, other}, module)
	if len(got) != 1 || got[0].Wallet != wallet || got[0].NewOwners[0] != owner || got[0].ExecutableAt != 1700000000 || got[0].TxHash != l.TxHash {
		t.Fatalf("got %+v", got)
	}
}
