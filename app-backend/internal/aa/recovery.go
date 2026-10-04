package aa

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Candide's Social Recovery Module v0.2.0 (recovery/contracts). A wallet
// enables it and adds guardians; guardians can start replacing the wallet's
// owners, which takes effect after the module's recovery period unless an
// owner cancels. Guardians can never execute anything else on the wallet.
var recoveryABI = mustABI(`[
 {"name":"enableModule","type":"function","inputs":[{"name":"module","type":"address"}],"outputs":[]},
 {"name":"addGuardianWithThreshold","type":"function","inputs":[{"name":"_guardian","type":"address"},{"name":"_threshold","type":"uint256"}],"outputs":[]},
 {"name":"revokeGuardianWithThreshold","type":"function","inputs":[{"name":"_prevGuardian","type":"address"},{"name":"_guardian","type":"address"},{"name":"_threshold","type":"uint256"}],"outputs":[]},
 {"name":"cancelRecovery","type":"function","inputs":[],"outputs":[]},
 {"name":"confirmRecovery","type":"function","inputs":[{"name":"_wallet","type":"address"},{"name":"_newOwners","type":"address[]"},{"name":"_newThreshold","type":"uint256"},{"name":"_nonce","type":"uint256"},{"name":"_execute","type":"bool"}],"outputs":[]},
 {"name":"finalizeRecovery","type":"function","inputs":[{"name":"_wallet","type":"address"}],"outputs":[]},
 {"name":"getRecoveryRequest","type":"function","stateMutability":"view","inputs":[{"name":"_wallet","type":"address"}],
  "outputs":[{"name":"request","type":"tuple","components":[
    {"name":"guardiansApprovalCount","type":"uint256"},{"name":"newThreshold","type":"uint256"},{"name":"nonce","type":"uint256"},
    {"name":"executableAt","type":"uint64"},{"name":"newOwners","type":"address[]"}]}]},
 {"name":"nonce","type":"function","stateMutability":"view","inputs":[{"name":"_wallet","type":"address"}],"outputs":[{"type":"uint256"}]},
 {"name":"isGuardian","type":"function","stateMutability":"view","inputs":[{"name":"_wallet","type":"address"},{"name":"_guardian","type":"address"}],"outputs":[{"type":"bool"}]},
 {"name":"guardiansCount","type":"function","stateMutability":"view","inputs":[{"name":"_wallet","type":"address"}],"outputs":[{"type":"uint256"}]},
 {"name":"VERSION","type":"function","stateMutability":"view","inputs":[],"outputs":[{"type":"string"}]},
 {"name":"RecoveryExecuted","type":"event","anonymous":false,"inputs":[
   {"name":"wallet","type":"address","indexed":true},{"name":"recoveryHash","type":"bytes32","indexed":true},
   {"name":"newOwners","type":"address[]","indexed":false},{"name":"newThreshold","type":"uint256","indexed":false},
   {"name":"nonce","type":"uint256","indexed":false},{"name":"executableAt","type":"uint64","indexed":false},
   {"name":"guardiansApprovalCount","type":"uint256","indexed":false}]}
]`)

// sentinelGuardian heads the module's guardian list.
var sentinelGuardian = common.HexToAddress("0x0000000000000000000000000000000000000001")

func recoveryCall(to common.Address, method string, args ...interface{}) Call {
	data, err := recoveryABI.Pack(method, args...)
	if err != nil {
		panic(err)
	}
	return Call{To: to, Value: big.NewInt(0), Data: data}
}

// EnableRecoveryCalls are a wallet's calls enabling the recovery module
// (unless enabled already) and making guardian its only guardian.
func EnableRecoveryCalls(wallet, module, guardian common.Address, moduleEnabled bool) []Call {
	var calls []Call
	if !moduleEnabled {
		calls = append(calls, recoveryCall(wallet, "enableModule", module))
	}
	return append(calls, recoveryCall(module, "addGuardianWithThreshold", guardian, big.NewInt(1)))
}

// RevokeGuardianCall is a wallet's call removing guardian, its only
// guardian, which ends its recovery coverage.
func RevokeGuardianCall(module, guardian common.Address) Call {
	return recoveryCall(module, "revokeGuardianWithThreshold", sentinelGuardian, guardian, big.NewInt(0))
}

// CancelRecoveryCall is a wallet's call cancelling its ongoing recovery.
func CancelRecoveryCall(module common.Address) Call {
	return recoveryCall(module, "cancelRecovery")
}

// ConfirmRecoveryCall is the guardian's call starting the replacement of
// wallet's owners by newOwners (threshold newThreshold) at the module's
// current nonce for the wallet.
func ConfirmRecoveryCall(module, wallet common.Address, newOwners []common.Address, newThreshold int64, nonce *big.Int) Call {
	return recoveryCall(module, "confirmRecovery", wallet, newOwners, big.NewInt(newThreshold), nonce, true)
}

// FinalizeRecoveryCall completes wallet's recovery once its period is over
// (anyone may call it).
func FinalizeRecoveryCall(module, wallet common.Address) Call {
	return recoveryCall(module, "finalizeRecovery", wallet)
}

// RecoveryRequest is a wallet's ongoing recovery (ExecutableAt 0: none).
type RecoveryRequest struct {
	NewOwners    []common.Address
	NewThreshold *big.Int
	Nonce        *big.Int
	ExecutableAt uint64
}

func recoveryView(ctx context.Context, r Caller, module common.Address, method string, args ...interface{}) ([]interface{}, error) {
	data, err := recoveryABI.Pack(method, args...)
	if err != nil {
		return nil, err
	}
	out, err := callRaw(ctx, r, module, data)
	if err != nil {
		return nil, err
	}
	return recoveryABI.Unpack(method, out)
}

// ReadRecoveryRequest reads wallet's ongoing recovery.
func ReadRecoveryRequest(ctx context.Context, r Caller, module, wallet common.Address) (*RecoveryRequest, error) {
	res, err := recoveryView(ctx, r, module, "getRecoveryRequest", wallet)
	if err != nil {
		return nil, err
	}
	q := res[0].(struct {
		GuardiansApprovalCount *big.Int         `json:"guardiansApprovalCount"`
		NewThreshold           *big.Int         `json:"newThreshold"`
		Nonce                  *big.Int         `json:"nonce"`
		ExecutableAt           uint64           `json:"executableAt"`
		NewOwners              []common.Address `json:"newOwners"`
	})
	return &RecoveryRequest{NewOwners: q.NewOwners, NewThreshold: q.NewThreshold, Nonce: q.Nonce, ExecutableAt: q.ExecutableAt}, nil
}

// RecoveryNonce is the module's current nonce for wallet.
func RecoveryNonce(ctx context.Context, r Caller, module, wallet common.Address) (*big.Int, error) {
	res, err := recoveryView(ctx, r, module, "nonce", wallet)
	if err != nil {
		return nil, err
	}
	return res[0].(*big.Int), nil
}

// IsGuardian reports whether guardian guards wallet.
func IsGuardian(ctx context.Context, r Caller, module, wallet, guardian common.Address) (bool, error) {
	res, err := recoveryView(ctx, r, module, "isGuardian", wallet, guardian)
	if err != nil {
		return false, err
	}
	return res[0].(bool), nil
}

// ModuleEnabled reports whether wallet (a deployed Safe) has module enabled.
func ModuleEnabled(ctx context.Context, r Caller, wallet, module common.Address) (bool, error) {
	out, err := callRaw(ctx, r, wallet, pack("isModuleEnabled", module))
	if err != nil {
		return false, err
	}
	res, err := contractsABI.Unpack("isModuleEnabled", out)
	if err != nil {
		return false, err
	}
	return res[0].(bool), nil
}

// RecoveryStarted is a RecoveryExecuted event: a recovery of Wallet began.
type RecoveryStarted struct {
	Wallet       common.Address
	NewOwners    []common.Address
	ExecutableAt uint64
	TxHash       common.Hash
}

// ParseRecoveryStarted finds the recoveries started in logs from module.
func ParseRecoveryStarted(logs []types.Log, module common.Address) []RecoveryStarted {
	ev := recoveryABI.Events["RecoveryExecuted"]
	var out []RecoveryStarted
	for _, l := range logs {
		if l.Address != module || len(l.Topics) != 3 || l.Topics[0] != ev.ID {
			continue
		}
		vals, err := ev.Inputs.NonIndexed().Unpack(l.Data)
		if err != nil {
			continue
		}
		out = append(out, RecoveryStarted{Wallet: common.BytesToAddress(l.Topics[1].Bytes()), NewOwners: vals[0].([]common.Address), ExecutableAt: vals[3].(uint64), TxHash: l.TxHash})
	}
	return out
}

// RecoveryStartedTopic is the RecoveryExecuted event's topic, for log
// queries.
func RecoveryStartedTopic() common.Hash { return recoveryABI.Events["RecoveryExecuted"].ID }

// ReplaceOwner returns owners with old replaced by replacement.
func ReplaceOwner(owners []common.Address, old, replacement common.Address) []common.Address {
	out := make([]common.Address, len(owners))
	for i, o := range owners {
		if o == old {
			o = replacement
		}
		out[i] = o
	}
	return out
}
