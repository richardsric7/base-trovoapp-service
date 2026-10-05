// Package chain is payout-engine's Base access: an RPC client with request
// deadlines, ERC-20 reads and the Transfer-log replay that snapshots token
// holders.
package chain

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"net"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

// Dial connects to the RPC; over HTTP(S) every request gets a 30s deadline
// so a stalled RPC cannot hang the engine.
func Dial(endpoint string) (*ethclient.Client, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return ethclient.Dial(endpoint)
	}
	transport := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           (&net.Dialer{Timeout: 10 * time.Second, KeepAlive: 30 * time.Second}).DialContext,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 30 * time.Second,
		MaxIdleConnsPerHost:   10,
		IdleConnTimeout:       90 * time.Second,
	}
	c, err := rpc.DialOptions(context.Background(), endpoint, rpc.WithHTTPClient(&http.Client{Transport: transport, Timeout: 30 * time.Second}))
	if err != nil {
		return nil, err
	}
	return ethclient.NewClient(c), nil
}

var erc20 abi.ABI

// TransferTopic is keccak256("Transfer(address,address,uint256)").
var TransferTopic = crypto.Keccak256Hash([]byte("Transfer(address,address,uint256)"))

func init() {
	var err error
	erc20, err = abi.JSON(strings.NewReader(`[
		{"name":"balanceOf","type":"function","stateMutability":"view","inputs":[{"type":"address"}],"outputs":[{"type":"uint256"}]},
		{"name":"decimals","type":"function","stateMutability":"view","inputs":[],"outputs":[{"type":"uint8"}]},
		{"name":"transfer","type":"function","inputs":[{"type":"address"},{"type":"uint256"}],"outputs":[{"type":"bool"}]}
	]`))
	if err != nil {
		panic(err)
	}
}

// TransferData is the calldata of token.transfer(to, amount).
func TransferData(to common.Address, amount *big.Int) []byte {
	data, _ := erc20.Pack("transfer", to, amount)
	return data
}

// BalanceOf reads token.balanceOf(holder) (at block, or latest when nil).
func BalanceOf(ctx context.Context, c *ethclient.Client, token, holder common.Address, block *big.Int) (*big.Int, error) {
	data, _ := erc20.Pack("balanceOf", holder)
	out, err := c.CallContract(ctx, ethereum.CallMsg{To: &token, Data: data}, block)
	if err != nil {
		return nil, err
	}
	if len(out) < 32 {
		return nil, fmt.Errorf("balanceOf on %s returned %d bytes", token.Hex(), len(out))
	}
	return new(big.Int).SetBytes(out[:32]), nil
}

// Decimals reads token.decimals().
func Decimals(ctx context.Context, c *ethclient.Client, token common.Address) (int, error) {
	data, _ := erc20.Pack("decimals")
	out, err := c.CallContract(ctx, ethereum.CallMsg{To: &token, Data: data}, nil)
	if err != nil {
		return 0, err
	}
	if len(out) < 32 {
		return 0, fmt.Errorf("decimals on %s returned %d bytes", token.Hex(), len(out))
	}
	return int(new(big.Int).SetBytes(out[:32]).Int64()), nil
}

// Transfer is one ERC-20 Transfer event.
type Transfer struct {
	From, To common.Address
	Amount   *big.Int
}

// Transfers returns token's Transfer events in [from, to], in chain order.
func Transfers(ctx context.Context, c *ethclient.Client, token common.Address, from, to uint64) ([]Transfer, error) {
	logs, err := c.FilterLogs(ctx, ethereum.FilterQuery{
		FromBlock: new(big.Int).SetUint64(from), ToBlock: new(big.Int).SetUint64(to),
		Addresses: []common.Address{token}, Topics: [][]common.Hash{{TransferTopic}},
	})
	if err != nil {
		return nil, err
	}
	sort.SliceStable(logs, func(i, j int) bool {
		if logs[i].BlockNumber != logs[j].BlockNumber {
			return logs[i].BlockNumber < logs[j].BlockNumber
		}
		return logs[i].Index < logs[j].Index
	})
	out := make([]Transfer, 0, len(logs))
	for _, l := range logs {
		if l.Removed || len(l.Topics) != 3 || len(l.Data) < 32 {
			continue
		}
		out = append(out, Transfer{
			From:   common.BytesToAddress(l.Topics[1].Bytes()),
			To:     common.BytesToAddress(l.Topics[2].Bytes()),
			Amount: new(big.Int).SetBytes(l.Data[:32]),
		})
	}
	return out, nil
}

// BlockAtTime is the first block at or after t (block headers are kept by
// every node, unlike historical state), found by binary search.
func BlockAtTime(ctx context.Context, c *ethclient.Client, t time.Time) (uint64, error) {
	head, err := c.HeaderByNumber(ctx, nil)
	if err != nil {
		return 0, err
	}
	target := uint64(t.Unix())
	if head.Time <= target {
		return head.Number.Uint64(), nil
	}
	lo, hi := uint64(0), head.Number.Uint64()
	for lo < hi {
		mid := (lo + hi) / 2
		h, err := c.HeaderByNumber(ctx, new(big.Int).SetUint64(mid))
		if err != nil {
			return 0, err
		}
		if h.Time < target {
			lo = mid + 1
		} else {
			hi = mid
		}
	}
	return lo, nil
}

// ErrReceiptTimeout is returned when a transaction is not mined in time.
var ErrReceiptTimeout = errors.New("transaction not mined in time")

// WaitReceipt waits for hash to be mined.
func WaitReceipt(ctx context.Context, c *ethclient.Client, hash common.Hash, timeout time.Duration) (*types.Receipt, error) {
	deadline := time.Now().Add(timeout)
	for {
		r, err := c.TransactionReceipt(ctx, hash)
		if err == nil && r != nil {
			return r, nil
		}
		if time.Now().After(deadline) {
			return nil, ErrReceiptTimeout
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(2 * time.Second):
		}
	}
}
