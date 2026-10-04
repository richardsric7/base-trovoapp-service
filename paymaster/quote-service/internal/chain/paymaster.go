// Package chain reads TrovoTokenPaymaster and EntryPoint state.
package chain

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

var contractsABI = func() abi.ABI {
	a, err := abi.JSON(strings.NewReader(`[
 {"name":"tokens","type":"function","stateMutability":"view","inputs":[{"type":"address"}],"outputs":[
   {"name":"enabled","type":"bool"},{"name":"minRate","type":"uint256"},{"name":"maxRate","type":"uint256"}]},
 {"name":"quoteSigners","type":"function","stateMutability":"view","inputs":[{"type":"address"}],"outputs":[{"type":"bool"}]},
 {"name":"paused","type":"function","stateMutability":"view","inputs":[],"outputs":[{"type":"bool"}]},
 {"name":"entryPoint","type":"function","stateMutability":"view","inputs":[],"outputs":[{"type":"address"}]},
 {"name":"debtTokenCount","type":"function","stateMutability":"view","inputs":[{"type":"address"}],"outputs":[{"type":"uint256"}]},
 {"name":"balanceOf","type":"function","stateMutability":"view","inputs":[{"type":"address"}],"outputs":[{"type":"uint256"}]}
]`))
	if err != nil {
		panic(err)
	}
	return a
}()

func call(ctx context.Context, c bind.ContractCaller, to common.Address, method string, args ...interface{}) ([]interface{}, error) {
	data, err := contractsABI.Pack(method, args...)
	if err != nil {
		return nil, err
	}
	out, err := c.CallContract(ctx, ethereum.CallMsg{To: &to, Data: data}, nil)
	if err != nil {
		return nil, fmt.Errorf("%s() on %s: %w", method, to.Hex(), err)
	}
	res, err := contractsABI.Unpack(method, out)
	if err != nil {
		return nil, fmt.Errorf("%s() on %s: %w", method, to.Hex(), err)
	}
	return res, nil
}

// TokenBounds is the paymaster's on-chain configuration of a gas token.
type TokenBounds struct {
	Enabled bool
	MinRate *big.Int
	MaxRate *big.Int
}

// Paymaster reads the deployed paymaster, caching token bounds briefly so
// every quote does not cost an RPC call.
type Paymaster struct {
	Caller     bind.ContractCaller
	Address    common.Address
	EntryPoint common.Address
	CacheTTL   time.Duration

	mu     sync.Mutex
	bounds map[common.Address]cachedBounds
}

type cachedBounds struct {
	b  TokenBounds
	at time.Time
}

func (p *Paymaster) TokenBounds(ctx context.Context, token common.Address) (TokenBounds, error) {
	p.mu.Lock()
	if c, ok := p.bounds[token]; ok && time.Since(c.at) < p.CacheTTL {
		p.mu.Unlock()
		return c.b, nil
	}
	p.mu.Unlock()

	res, err := call(ctx, p.Caller, p.Address, "tokens", token)
	if err != nil {
		return TokenBounds{}, err
	}
	b := TokenBounds{Enabled: res[0].(bool), MinRate: res[1].(*big.Int), MaxRate: res[2].(*big.Int)}
	p.mu.Lock()
	if p.bounds == nil {
		p.bounds = map[common.Address]cachedBounds{}
	}
	p.bounds[token] = cachedBounds{b, time.Now()}
	p.mu.Unlock()
	return b, nil
}

// IsQuoteSigner reports whether the paymaster accepts signer's quotes.
func (p *Paymaster) IsQuoteSigner(ctx context.Context, signer common.Address) (bool, error) {
	res, err := call(ctx, p.Caller, p.Address, "quoteSigners", signer)
	if err != nil {
		return false, err
	}
	return res[0].(bool), nil
}

func (p *Paymaster) Paused(ctx context.Context) (bool, error) {
	res, err := call(ctx, p.Caller, p.Address, "paused")
	if err != nil {
		return false, err
	}
	return res[0].(bool), nil
}

// ConfiguredEntryPoint is the EntryPoint the paymaster was deployed against.
func (p *Paymaster) ConfiguredEntryPoint(ctx context.Context) (common.Address, error) {
	res, err := call(ctx, p.Caller, p.Address, "entryPoint")
	if err != nil {
		return common.Address{}, err
	}
	return res[0].(common.Address), nil
}

// HasDebt reports whether sender owes the paymaster from a failed charge
// (the paymaster rejects its operations until settled).
func (p *Paymaster) HasDebt(ctx context.Context, sender common.Address) (bool, error) {
	res, err := call(ctx, p.Caller, p.Address, "debtTokenCount", sender)
	if err != nil {
		return false, err
	}
	return res[0].(*big.Int).Sign() != 0, nil
}

// Deposit is the paymaster's ETH deposit at the EntryPoint (wei).
func (p *Paymaster) Deposit(ctx context.Context) (*big.Int, error) {
	res, err := call(ctx, p.Caller, p.EntryPoint, "balanceOf", p.Address)
	if err != nil {
		return nil, err
	}
	return res[0].(*big.Int), nil
}

// TokenBalance is an ERC-20 balance.
func TokenBalance(ctx context.Context, c bind.ContractCaller, token, holder common.Address) (*big.Int, error) {
	res, err := call(ctx, c, token, "balanceOf", holder)
	if err != nil {
		return nil, err
	}
	return res[0].(*big.Int), nil
}
