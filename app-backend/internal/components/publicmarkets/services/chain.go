package publicmarkets

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"
	"time"

	"trovo-wallet-api/internal/aa"
	"trovo-wallet-api/internal/evmkeypair"
	"trovo-wallet-api/internal/gnosissafe"
	"trovo-wallet-api/internal/network"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Transfer is one ERC-20 Transfer log.
type Transfer struct {
	From, To common.Address
	Value    *big.Int
	Block    uint64
	TxHash   string
	LogIndex uint
}

// Chain is everything Public Markets does on Base. The engine only uses
// this interface, so its logic is tested without a chain.
type Chain interface {
	Decimals(ctx context.Context, token common.Address) (uint8, error)
	TotalSupply(ctx context.Context, token common.Address) (*big.Int, error)
	BalanceOf(ctx context.Context, token, holder common.Address) (*big.Int, error)
	Owner(ctx context.Context, contract common.Address) (common.Address, error)
	// Send executes calls as one transaction of safe, signed by the
	// Public Markets signers, and returns its hash once broadcast.
	Send(ctx context.Context, safe common.Address, calls []gnosissafe.Call) (string, error)
	// Wait waits for a sent transaction and reports whether it succeeded.
	Wait(ctx context.Context, hash string) error
	// PartnerSafe returns the address of an exchange customer's wallet
	// (a Safe owned by the Public Markets signers) and deploys it when
	// deploy is set.
	PartnerSafe(ctx context.Context, saltNonce *big.Int, deploy bool) (common.Address, error)
	Head(ctx context.Context) (uint64, error)
	Transfers(ctx context.Context, token common.Address, from, to uint64) ([]Transfer, error)
	BlockAtTime(ctx context.Context, t time.Time) (uint64, error)
	// TxTransfers are a mined transaction's Transfer logs of token.
	TxTransfers(ctx context.Context, hash string, token common.Address) ([]Transfer, error)
	// TreasurySafe holds the funding stablecoin (payments in, redemptions
	// and dividends out).
	TreasurySafe() common.Address
}

var erc20 = func() abi.ABI {
	a, err := abi.JSON(strings.NewReader(`[
 {"name":"decimals","type":"function","stateMutability":"view","inputs":[],"outputs":[{"type":"uint8"}]},
 {"name":"totalSupply","type":"function","stateMutability":"view","inputs":[],"outputs":[{"type":"uint256"}]},
 {"name":"balanceOf","type":"function","stateMutability":"view","inputs":[{"name":"a","type":"address"}],"outputs":[{"type":"uint256"}]},
 {"name":"owner","type":"function","stateMutability":"view","inputs":[],"outputs":[{"type":"address"}]}]`))
	if err != nil {
		panic(err)
	}
	return a
}()

var transferTopic = crypto.Keccak256Hash([]byte("Transfer(address,address,uint256)"))

// DefaultFallbackHandler is Safe v1.4.1's CompatibilityFallbackHandler.
const DefaultFallbackHandler = "0xfd0732Dc9E303f09fCEf3a7388Ad10A83459Ec99"

// BaseChain is the Chain on the configured RPC.
type BaseChain struct {
	Client  *ethclient.Client
	ChainID *big.Int
}

// Signers are the Public Markets signer keys: a managed secret
// (PUBLIC_MARKETS_SIGNERS, ";"-separated) whose owners sign the issuing
// Safes, the treasury and exchange customers' wallets.
func Signers() ([]*evmkeypair.Full, error) {
	raw := strings.TrimSpace(os.Getenv("PUBLIC_MARKETS_SIGNERS"))
	if raw == "" {
		return nil, errors.New("PUBLIC_MARKETS_SIGNERS is not configured")
	}
	var out []*evmkeypair.Full
	for i, p := range strings.Split(raw, ";") {
		if p = strings.TrimSpace(p); p == "" {
			continue
		}
		kp, err := evmkeypair.ParseFull(p)
		if err != nil {
			return nil, fmt.Errorf("PUBLIC_MARKETS_SIGNERS entry %d is invalid: %w", i+1, err)
		}
		out = append(out, kp)
	}
	if len(out) == 0 {
		return nil, errors.New("PUBLIC_MARKETS_SIGNERS has no keys")
	}
	return out, nil
}

func (c *BaseChain) call(ctx context.Context, to common.Address, method string, args ...interface{}) ([]interface{}, error) {
	data, err := erc20.Pack(method, args...)
	if err != nil {
		return nil, err
	}
	out, err := c.Client.CallContract(ctx, ethereum.CallMsg{To: &to, Data: data}, nil)
	if err != nil {
		return nil, err
	}
	return erc20.Unpack(method, out)
}

func (c *BaseChain) Decimals(ctx context.Context, token common.Address) (uint8, error) {
	v, err := c.call(ctx, token, "decimals")
	if err != nil || len(v) == 0 {
		return 0, fmt.Errorf("reading decimals of %s: %v", token.Hex(), err)
	}
	return v[0].(uint8), nil
}

func (c *BaseChain) TotalSupply(ctx context.Context, token common.Address) (*big.Int, error) {
	v, err := c.call(ctx, token, "totalSupply")
	if err != nil || len(v) == 0 {
		return nil, fmt.Errorf("reading totalSupply of %s: %v", token.Hex(), err)
	}
	return v[0].(*big.Int), nil
}

func (c *BaseChain) BalanceOf(ctx context.Context, token, holder common.Address) (*big.Int, error) {
	v, err := c.call(ctx, token, "balanceOf", holder)
	if err != nil || len(v) == 0 {
		return nil, fmt.Errorf("reading balance of %s: %v", holder.Hex(), err)
	}
	return v[0].(*big.Int), nil
}

func (c *BaseChain) Owner(ctx context.Context, contract common.Address) (common.Address, error) {
	v, err := c.call(ctx, contract, "owner")
	if err != nil || len(v) == 0 {
		return common.Address{}, fmt.Errorf("reading owner of %s: %v", contract.Hex(), err)
	}
	return v[0].(common.Address), nil
}

func multiSendAddress() common.Address {
	if v := strings.TrimSpace(os.Getenv("SAFE_MULTISEND_CALL_ONLY_ADDRESS")); common.IsHexAddress(v) {
		return common.HexToAddress(v)
	}
	return common.HexToAddress(gnosissafe.DefaultMultiSendCallOnlyAddress)
}

func (c *BaseChain) Send(ctx context.Context, safe common.Address, calls []gnosissafe.Call) (string, error) {
	candidates, err := Signers()
	if err != nil {
		return "", err
	}
	signers, err := gnosissafe.SignersForSafe(ctx, c.Client, safe, candidates)
	if err != nil {
		return "", err
	}
	return gnosissafe.ExecCalls(ctx, c.Client, c.ChainID, safe, signers, calls, multiSendAddress())
}

func (c *BaseChain) Wait(ctx context.Context, hash string) error {
	return gnosissafe.WaitSuccess(ctx, c.Client, common.HexToHash(hash))
}

// partnerSafeOwners are the owners and threshold of exchange customers'
// wallets: every Public Markets signer, with the threshold the managed
// secret standard signs with (at most 3).
func partnerSafeOwners() ([]common.Address, int64, error) {
	signers, err := Signers()
	if err != nil {
		return nil, 0, err
	}
	owners := make([]common.Address, 0, len(signers))
	for _, s := range signers {
		owners = append(owners, common.HexToAddress(s.Address()))
	}
	threshold := int64(len(owners))
	if threshold > 3 {
		threshold = 3
	}
	if v, err := strconv.ParseInt(os.Getenv("PUBLIC_MARKETS_PARTNER_WALLET_THRESHOLD"), 10, 64); err == nil && v >= 1 && v <= int64(len(owners)) {
		threshold = v
	}
	return owners, threshold, nil
}

func deployConfig(chainID *big.Int) gnosissafe.DeployConfig {
	cfg := aa.ConfigFromEnv(chainID)
	fallback := common.HexToAddress(DefaultFallbackHandler)
	if v := strings.TrimSpace(os.Getenv("SAFE_FALLBACK_HANDLER_ADDRESS")); common.IsHexAddress(v) {
		fallback = common.HexToAddress(v)
	}
	return gnosissafe.DeployConfig{ProxyFactory: cfg.ProxyFactory, Singleton: cfg.Singleton, FallbackHandler: fallback}
}

func (c *BaseChain) PartnerSafe(ctx context.Context, saltNonce *big.Int, deploy bool) (common.Address, error) {
	owners, threshold, err := partnerSafeOwners()
	if err != nil {
		return common.Address{}, err
	}
	cfg := deployConfig(c.ChainID)
	if !deploy {
		addr, _, err := gnosissafe.PredictSafeAddress(ctx, c.Client, cfg, owners, threshold, saltNonce)
		return addr, err
	}
	signers, err := Signers()
	if err != nil {
		return common.Address{}, err
	}
	return gnosissafe.DeploySafe(ctx, c.Client, c.ChainID, signers[0], cfg, owners, threshold, saltNonce)
}

func (c *BaseChain) Head(ctx context.Context) (uint64, error) { return c.Client.BlockNumber(ctx) }

func (c *BaseChain) Transfers(ctx context.Context, token common.Address, from, to uint64) ([]Transfer, error) {
	logs, err := c.Client.FilterLogs(ctx, ethereum.FilterQuery{
		FromBlock: new(big.Int).SetUint64(from), ToBlock: new(big.Int).SetUint64(to),
		Addresses: []common.Address{token}, Topics: [][]common.Hash{{transferTopic}},
	})
	if err != nil {
		return nil, err
	}
	out := make([]Transfer, 0, len(logs))
	for _, l := range logs {
		if t, ok := parseTransfer(l); ok {
			out = append(out, t)
		}
	}
	return out, nil
}

func parseTransfer(l types.Log) (Transfer, bool) {
	if len(l.Topics) != 3 || len(l.Data) < 32 || l.Removed {
		return Transfer{}, false
	}
	return Transfer{From: common.BytesToAddress(l.Topics[1].Bytes()), To: common.BytesToAddress(l.Topics[2].Bytes()),
		Value: new(big.Int).SetBytes(l.Data[:32]), Block: l.BlockNumber, TxHash: l.TxHash.Hex(), LogIndex: l.Index}, true
}

// BlockAtTime is the last block at or before t (binary search on headers).
func (c *BaseChain) BlockAtTime(ctx context.Context, t time.Time) (uint64, error) {
	head, err := c.Client.HeaderByNumber(ctx, nil)
	if err != nil {
		return 0, err
	}
	target := uint64(t.Unix())
	if head.Time <= target {
		return head.Number.Uint64(), nil
	}
	lo, hi := uint64(0), head.Number.Uint64()
	for lo < hi {
		mid := (lo + hi + 1) / 2
		h, err := c.Client.HeaderByNumber(ctx, new(big.Int).SetUint64(mid))
		if err != nil {
			return 0, err
		}
		if h.Time <= target {
			lo = mid
		} else {
			hi = mid - 1
		}
	}
	return lo, nil
}

func (c *BaseChain) TxTransfers(ctx context.Context, hash string, token common.Address) ([]Transfer, error) {
	r, err := c.Client.TransactionReceipt(ctx, common.HexToHash(hash))
	if err != nil {
		return nil, err
	}
	if r.Status != types.ReceiptStatusSuccessful {
		return nil, errors.New("the transaction failed")
	}
	var out []Transfer
	for _, l := range r.Logs {
		if l.Address != token || len(l.Topics) == 0 || l.Topics[0] != transferTopic {
			continue
		}
		if t, ok := parseTransfer(*l); ok {
			out = append(out, t)
		}
	}
	return out, nil
}

func (c *BaseChain) TreasurySafe() common.Address {
	return common.HexToAddress(os.Getenv("PUBLIC_MARKETS_TREASURY_SAFE"))
}

// NewBaseChain is the Chain on the shared RPC client, nil when the chain
// is not configured (the engine then only does database work).
func NewBaseChain() Chain {
	if strings.TrimSpace(os.Getenv("BASE_RPC_URL")) == "" {
		return nil
	}
	client := network.GetBlockchainClient()
	if client == nil {
		return nil
	}
	return &BaseChain{Client: client, ChainID: network.GetBlockchainChainID()}
}

// units converts a human amount to base units (truncating).
func units(human string, decimals int) *big.Int {
	return d(human).Shift(int32(decimals)).Truncate(0).BigInt()
}

func callOf(c aa.Call) gnosissafe.Call {
	return gnosissafe.Call{To: c.To, Value: big.NewInt(0), Data: c.Data}
}
