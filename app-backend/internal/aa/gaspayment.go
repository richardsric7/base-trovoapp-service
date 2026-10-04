package aa

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// The TrovoTokenPaymaster's debt functions and charge events.
var paymasterABI = mustABI(`[
 {"name":"debt","type":"function","stateMutability":"view","inputs":[{"name":"sender","type":"address"},{"name":"token","type":"address"}],"outputs":[{"type":"uint256"}]},
 {"name":"settleDebt","type":"function","inputs":[{"name":"sender","type":"address"},{"name":"token","type":"address"}],"outputs":[]},
 {"name":"GasPaidInToken","type":"event","anonymous":false,"inputs":[
   {"name":"sender","type":"address","indexed":true},{"name":"token","type":"address","indexed":true},
   {"name":"userOpHash","type":"bytes32","indexed":true},{"name":"exchangeRate","type":"uint256","indexed":false},
   {"name":"tokenCost","type":"uint256","indexed":false},{"name":"actualGasCost","type":"uint256","indexed":false}]},
 {"name":"DebtSettled","type":"event","anonymous":false,"inputs":[
   {"name":"sender","type":"address","indexed":true},{"name":"token","type":"address","indexed":true},
   {"name":"amount","type":"uint256","indexed":false},{"name":"forgiven","type":"bool","indexed":false}]},
 {"name":"GasChargeFailed","type":"event","anonymous":false,"inputs":[
   {"name":"sender","type":"address","indexed":true},{"name":"token","type":"address","indexed":true},
   {"name":"userOpHash","type":"bytes32","indexed":true},{"name":"tokenCost","type":"uint256","indexed":false}]}
]`)

// PaymasterDebt is what sender owes the paymaster in token: gas it could not
// collect after an operation. While a wallet owes any, the paymaster
// refuses to pay its gas.
func PaymasterDebt(ctx context.Context, r Caller, paymaster, sender, token common.Address) (*big.Int, error) {
	data, err := paymasterABI.Pack("debt", sender, token)
	if err != nil {
		return nil, err
	}
	out, err := callRaw(ctx, r, paymaster, data)
	if err != nil {
		return nil, err
	}
	res, err := paymasterABI.Unpack("debt", out)
	if err != nil {
		return nil, err
	}
	return res[0].(*big.Int), nil
}

// SettleDebtCalls are the calls a wallet makes to pay its debt in token:
// allow the paymaster to take it, then settleDebt.
func SettleDebtCalls(paymaster, sender, token common.Address) []Call {
	data, _ := paymasterABI.Pack("settleDebt", sender, token)
	return []Call{ERC20Approve(token, paymaster, maxUint256), {To: paymaster, Value: big.NewInt(0), Data: data}}
}

// GasPayment is how the paymaster charged an operation, from its events.
type GasPayment struct {
	Token     common.Address
	TokenCost *big.Int
	// Failed means the charge failed and was recorded as debt.
	Failed bool
}

// ParseGasPayment finds the paymaster's charge for userOpHash in a
// transaction's logs (nil when it paid no gas for it).
func ParseGasPayment(logs []*types.Log, paymaster common.Address, userOpHash common.Hash) *GasPayment {
	paid := paymasterABI.Events["GasPaidInToken"]
	failed := paymasterABI.Events["GasChargeFailed"]
	for _, l := range logs {
		if l.Address != paymaster || len(l.Topics) != 4 || l.Topics[3] != userOpHash {
			continue
		}
		var ev abi.Event
		switch l.Topics[0] {
		case paid.ID:
			ev = paid
		case failed.ID:
			ev = failed
		default:
			continue
		}
		vals, err := ev.Inputs.NonIndexed().Unpack(l.Data)
		if err != nil {
			continue
		}
		g := &GasPayment{Token: common.BytesToAddress(l.Topics[2].Bytes()), Failed: ev.Name == failed.Name}
		if g.Failed {
			g.TokenCost = vals[0].(*big.Int)
		} else {
			g.TokenCost = vals[1].(*big.Int)
		}
		return g
	}
	return nil
}

// DebtSettlement is a gas debt the wallet paid (DebtSettled, not forgiven).
type DebtSettlement struct {
	Token  common.Address
	Amount *big.Int
}

// ParseDebtSettlements finds the gas debts sender paid in a transaction.
func ParseDebtSettlements(logs []*types.Log, paymaster, sender common.Address) []DebtSettlement {
	ev := paymasterABI.Events["DebtSettled"]
	var out []DebtSettlement
	for _, l := range logs {
		if l.Address != paymaster || len(l.Topics) != 3 || l.Topics[0] != ev.ID || common.BytesToAddress(l.Topics[1].Bytes()) != sender {
			continue
		}
		vals, err := ev.Inputs.NonIndexed().Unpack(l.Data)
		if err != nil || vals[1].(bool) {
			continue
		}
		out = append(out, DebtSettlement{Token: common.BytesToAddress(l.Topics[2].Bytes()), Amount: vals[0].(*big.Int)})
	}
	return out
}
