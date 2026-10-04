package aa

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strings"
	"sync/atomic"
	"time"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// Bundler is a JSON-RPC client for our self-hosted ERC-4337 bundler
// (BUNDLER_URL).
type Bundler struct {
	URL    string
	HTTP   *http.Client
	nextID atomic.Int64
}

func NewBundler(url string) *Bundler {
	return &Bundler{URL: url, HTTP: &http.Client{Timeout: 30 * time.Second}}
}

// RPCError is a JSON-RPC error from the bundler (e.g. an AAxx validation
// failure).
type RPCError struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data,omitempty"`
}

func (e *RPCError) Error() string { return fmt.Sprintf("bundler error %d: %s", e.Code, e.Message) }

func (b *Bundler) call(ctx context.Context, method string, params []interface{}, out interface{}) error {
	body, err := json.Marshal(map[string]interface{}{"jsonrpc": "2.0", "id": b.nextID.Add(1), "method": method, "params": params})
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, b.URL, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := b.HTTP.Do(req)
	if err != nil {
		return fmt.Errorf("bundler %s: %w", method, err)
	}
	defer res.Body.Close()
	raw, err := io.ReadAll(io.LimitReader(res.Body, 4<<20))
	if err != nil {
		return err
	}
	var env struct {
		Result json.RawMessage `json:"result"`
		Error  *RPCError       `json:"error"`
	}
	if err := json.Unmarshal(raw, &env); err != nil {
		return fmt.Errorf("bundler %s: HTTP %d, invalid response", method, res.StatusCode)
	}
	if env.Error != nil {
		return env.Error
	}
	if out == nil {
		return nil
	}
	return json.Unmarshal(env.Result, out)
}

// GasEstimate is eth_estimateUserOperationGas's answer.
type GasEstimate struct {
	PreVerificationGas            *hexutil.Big `json:"preVerificationGas"`
	VerificationGasLimit          *hexutil.Big `json:"verificationGasLimit"`
	CallGasLimit                  *hexutil.Big `json:"callGasLimit"`
	PaymasterVerificationGasLimit *hexutil.Big `json:"paymasterVerificationGasLimit,omitempty"`
	PaymasterPostOpGasLimit       *hexutil.Big `json:"paymasterPostOpGasLimit,omitempty"`
}

func (b *Bundler) EstimateGas(ctx context.Context, op UserOperation, entryPoint common.Address) (GasEstimate, error) {
	var est GasEstimate
	err := b.call(ctx, "eth_estimateUserOperationGas", []interface{}{op.ToRPC(), entryPoint}, &est)
	return est, err
}

// Send submits a signed operation and returns its userOpHash.
func (b *Bundler) Send(ctx context.Context, op UserOperation, entryPoint common.Address) (common.Hash, error) {
	var h common.Hash
	err := b.call(ctx, "eth_sendUserOperation", []interface{}{op.ToRPC(), entryPoint}, &h)
	return h, err
}

// Receipt is eth_getUserOperationReceipt's answer (nil until included).
type Receipt struct {
	UserOpHash    common.Hash    `json:"userOpHash"`
	Sender        common.Address `json:"sender"`
	Success       bool           `json:"success"`
	Reason        string         `json:"reason"`
	ActualGasCost *hexutil.Big   `json:"actualGasCost"`
	ActualGasUsed *hexutil.Big   `json:"actualGasUsed"`
	Receipt       struct {
		TransactionHash common.Hash  `json:"transactionHash"`
		BlockNumber     *hexutil.Big `json:"blockNumber"`
	} `json:"receipt"`
}

func (b *Bundler) Receipt(ctx context.Context, userOpHash common.Hash) (*Receipt, error) {
	var r *Receipt
	err := b.call(ctx, "eth_getUserOperationReceipt", []interface{}{userOpHash}, &r)
	return r, err
}

// WaitReceipt polls for the operation's receipt until ctx ends.
func (b *Bundler) WaitReceipt(ctx context.Context, userOpHash common.Hash) (*Receipt, error) {
	t := time.NewTicker(2 * time.Second)
	defer t.Stop()
	for {
		r, err := b.Receipt(ctx, userOpHash)
		if err == nil && r != nil {
			return r, nil
		}
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("user operation %s not included: %w", userOpHash.Hex(), ctx.Err())
		case <-t.C:
		}
	}
}

// Caller reads contract state (*ethclient.Client satisfies it).
type Caller interface {
	CallContract(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error)
	CodeAt(ctx context.Context, account common.Address, blockNumber *big.Int) ([]byte, error)
}

func callRaw(ctx context.Context, c Caller, to common.Address, data []byte) ([]byte, error) {
	return c.CallContract(ctx, ethereum.CallMsg{To: &to, Data: data}, nil)
}

func callView(ctx context.Context, c Caller, to common.Address, method string, args ...interface{}) ([]interface{}, error) {
	data := pack(method, args...)
	out, err := c.CallContract(ctx, ethereum.CallMsg{To: &to, Data: data}, nil)
	if err != nil {
		return nil, fmt.Errorf("%s on %s: %w", method, to.Hex(), err)
	}
	return contractsABI.Unpack(method, out)
}

// EntryPointNonce reads the wallet's next UserOperation nonce for a nonce
// key (nil = key 0). Each key is its own sequence, so operations built on
// different keys do not invalidate each other.
func (c Config) EntryPointNonce(ctx context.Context, r Caller, sender common.Address, key ...*big.Int) (*big.Int, error) {
	k := big.NewInt(0)
	if len(key) > 0 && key[0] != nil {
		k = key[0]
	}
	res, err := callView(ctx, r, c.EntryPoint, "getNonce", sender, k)
	if err != nil {
		return nil, err
	}
	return res[0].(*big.Int), nil
}

// Deployed reports whether the wallet has been deployed (activated).
func Deployed(ctx context.Context, r Caller, addr common.Address) (bool, error) {
	code, err := r.CodeAt(ctx, addr, nil)
	if err != nil {
		return false, err
	}
	return len(code) > 0, nil
}

// OnchainOwners reads a deployed Safe's owners and threshold.
func OnchainOwners(ctx context.Context, r Caller, safe common.Address) ([]common.Address, int64, error) {
	res, err := callView(ctx, r, safe, "getOwners")
	if err != nil {
		return nil, 0, err
	}
	th, err := callView(ctx, r, safe, "getThreshold")
	if err != nil {
		return nil, 0, err
	}
	return res[0].([]common.Address), th[0].(*big.Int).Int64(), nil
}

// IsAAError reports whether err is a bundler validation rejection (the
// operation itself is invalid, as opposed to the bundler being unreachable).
func IsAAError(err error) bool {
	e, ok := err.(*RPCError)
	return ok && (strings.Contains(e.Message, "AA") || e.Code == -32500 || e.Code == -32501 || e.Code == -32502)
}
