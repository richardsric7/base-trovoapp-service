// Package safe builds the payout Safe's transactions. It ports what
// payout-engine needs from app-backend's internal/gnosissafe (EIP-712
// SafeTx hashing, owner signatures, MultiSend packing), with the steps split
// so the engine can record a signed transaction before broadcasting it.
package safe

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math/big"
	"sort"
	"strings"

	"trovo-payout-engine/internal/config"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	OperationCall         = uint8(0)
	OperationDelegateCall = uint8(1)
)

// Call is one call the Safe makes.
type Call struct {
	To   common.Address
	Data []byte
}

var safeABI, multiSendABI abi.ABI

// ExecutionSuccess / ExecutionFailure are the Safe's events for an
// execTransaction (Safe v1.3+).
var (
	ExecutionSuccessTopic = crypto.Keccak256Hash([]byte("ExecutionSuccess(bytes32,uint256)"))
	ExecutionFailureTopic = crypto.Keccak256Hash([]byte("ExecutionFailure(bytes32,uint256)"))
)

func init() {
	var err error
	safeABI, err = abi.JSON(strings.NewReader(`[
		{"constant":true,"inputs":[],"name":"nonce","outputs":[{"name":"","type":"uint256"}],"type":"function"},
		{"constant":true,"inputs":[],"name":"getThreshold","outputs":[{"name":"","type":"uint256"}],"type":"function"},
		{"constant":true,"inputs":[],"name":"getOwners","outputs":[{"name":"","type":"address[]"}],"type":"function"},
		{"constant":false,"inputs":[
			{"name":"to","type":"address"},{"name":"value","type":"uint256"},{"name":"data","type":"bytes"},{"name":"operation","type":"uint8"},
			{"name":"safeTxGas","type":"uint256"},{"name":"baseGas","type":"uint256"},{"name":"gasPrice","type":"uint256"},
			{"name":"gasToken","type":"address"},{"name":"refundReceiver","type":"address"},{"name":"signatures","type":"bytes"}
		],"name":"execTransaction","outputs":[{"name":"","type":"bool"}],"type":"function"}
	]`))
	if err != nil {
		panic(err)
	}
	multiSendABI, err = abi.JSON(strings.NewReader(`[{"inputs":[{"name":"transactions","type":"bytes"}],"name":"multiSend","outputs":[],"stateMutability":"payable","type":"function"}]`))
	if err != nil {
		panic(err)
	}
}

var (
	safeTxTypehash = crypto.Keccak256([]byte("SafeTx(address to,uint256 value,bytes data,uint8 operation,uint256 safeTxGas,uint256 baseGas,uint256 gasPrice,address gasToken,address refundReceiver,uint256 nonce)"))
	domainTypehash = crypto.Keccak256([]byte("EIP712Domain(uint256 chainId,address verifyingContract)"))
)

func b32(b []byte) [32]byte {
	var out [32]byte
	copy(out[:], b)
	return out
}

// TxHash is Safe's EIP-712 SafeTx hash for a transaction without gas refund.
func TxHash(chainID *big.Int, safe, to common.Address, data []byte, operation uint8, nonce *big.Int) ([]byte, error) {
	addressTy, _ := abi.NewType("address", "", nil)
	uint256Ty, _ := abi.NewType("uint256", "", nil)
	bytes32Ty, _ := abi.NewType("bytes32", "", nil)
	uint8Ty, _ := abi.NewType("uint8", "", nil)
	domain, err := abi.Arguments{{Type: bytes32Ty}, {Type: uint256Ty}, {Type: addressTy}}.Pack(b32(domainTypehash), chainID, safe)
	if err != nil {
		return nil, err
	}
	zero := big.NewInt(0)
	structArgs := abi.Arguments{{Type: bytes32Ty}, {Type: addressTy}, {Type: uint256Ty}, {Type: bytes32Ty}, {Type: uint8Ty},
		{Type: uint256Ty}, {Type: uint256Ty}, {Type: uint256Ty}, {Type: addressTy}, {Type: addressTy}, {Type: uint256Ty}}
	st, err := structArgs.Pack(b32(safeTxTypehash), to, zero, b32(crypto.Keccak256(data)), operation, zero, zero, zero, common.Address{}, common.Address{}, nonce)
	if err != nil {
		return nil, err
	}
	pre := append([]byte{0x19, 0x01}, crypto.Keccak256(domain)...)
	return crypto.Keccak256(append(pre, crypto.Keccak256(st)...)), nil
}

// Sign makes each signer's signature over the SafeTx hash, concatenated in
// ascending owner order as Safe's checkSignatures requires.
func Sign(hash []byte, signers []config.Signer) ([]byte, error) {
	type sig struct {
		addr common.Address
		raw  []byte
	}
	sigs := make([]sig, 0, len(signers))
	for _, s := range signers {
		raw, err := crypto.Sign(hash, s.Key)
		if err != nil {
			return nil, err
		}
		raw[64] += 27
		sigs = append(sigs, sig{s.Address, raw})
	}
	sort.Slice(sigs, func(i, j int) bool { return bytes.Compare(sigs[i].addr.Bytes(), sigs[j].addr.Bytes()) < 0 })
	out := make([]byte, 0, 65*len(sigs))
	for _, s := range sigs {
		out = append(out, s.raw...)
	}
	return out, nil
}

// Target is what the Safe executes for calls: the call itself, or a
// MultiSendCallOnly batch (delegatecall) for several, which all succeed or
// all revert together.
func Target(calls []Call, multiSend common.Address) (to common.Address, data []byte, op uint8, err error) {
	if len(calls) == 0 {
		return common.Address{}, nil, 0, errors.New("no calls")
	}
	if len(calls) == 1 {
		return calls[0].To, calls[0].Data, OperationCall, nil
	}
	var packed []byte
	for _, c := range calls {
		packed = append(packed, OperationCall)
		packed = append(packed, c.To.Bytes()...)
		packed = append(packed, make([]byte, 32)...) // value 0
		packed = append(packed, common.LeftPadBytes(big.NewInt(int64(len(c.Data))).Bytes(), 32)...)
		packed = append(packed, c.Data...)
	}
	data, err = multiSendABI.Pack("multiSend", packed)
	return multiSend, data, OperationDelegateCall, err
}

// ExecData is execTransaction's calldata.
func ExecData(to common.Address, data []byte, op uint8, signatures []byte) ([]byte, error) {
	zero := big.NewInt(0)
	return safeABI.Pack("execTransaction", to, zero, data, op, zero, zero, zero, common.Address{}, common.Address{}, signatures)
}

func callUint(ctx context.Context, c *ethclient.Client, safe common.Address, method string) (*big.Int, error) {
	data, _ := safeABI.Pack(method)
	out, err := c.CallContract(ctx, ethereum.CallMsg{To: &safe, Data: data}, nil)
	if err != nil {
		return nil, err
	}
	if len(out) < 32 {
		return nil, fmt.Errorf("%s on %s returned %d bytes (is it a Safe?)", method, safe.Hex(), len(out))
	}
	return new(big.Int).SetBytes(out[:32]), nil
}

// Nonce reads the Safe's nonce.
func Nonce(ctx context.Context, c *ethclient.Client, safe common.Address) (*big.Int, error) {
	return callUint(ctx, c, safe, "nonce")
}

// CheckSigners verifies the signers are owners of safe and meet its threshold.
func CheckSigners(ctx context.Context, c *ethclient.Client, safe common.Address, signers []config.Signer) error {
	th, err := callUint(ctx, c, safe, "getThreshold")
	if err != nil {
		return fmt.Errorf("reading the payout Safe's threshold: %w", err)
	}
	data, _ := safeABI.Pack("getOwners")
	out, err := c.CallContract(ctx, ethereum.CallMsg{To: &safe, Data: data}, nil)
	if err != nil {
		return fmt.Errorf("reading the payout Safe's owners: %w", err)
	}
	vals, err := safeABI.Unpack("getOwners", out)
	if err != nil {
		return err
	}
	owners, _ := vals[0].([]common.Address)
	isOwner := map[common.Address]bool{}
	for _, o := range owners {
		isOwner[o] = true
	}
	for _, s := range signers {
		if !isOwner[s.Address] {
			return fmt.Errorf("signer %s is not an owner of the payout Safe %s (rotated? check the vault manager)", s.Address.Hex(), safe.Hex())
		}
	}
	if th.Cmp(big.NewInt(int64(len(signers)))) > 0 {
		return fmt.Errorf("the payout Safe %s needs %v signatures; %d signers are configured", safe.Hex(), th, len(signers))
	}
	return nil
}

// Executed reports whether a mined receipt executed the Safe transaction
// safeTxHash. Payouts are sent with safeTxGas 0, so a failing inner call
// reverts the whole transaction: a successful receipt executed it, unless
// the Safe logged ExecutionFailure for it. (Safe v1.3+ indexes the hash in
// ExecutionSuccess/ExecutionFailure's topics; older versions put it in the
// data - both are matched.)
func Executed(r *types.Receipt, safe common.Address, safeTxHash []byte) bool {
	if r == nil || r.Status != types.ReceiptStatusSuccessful {
		return false
	}
	for _, l := range r.Logs {
		if l.Address != safe || len(l.Topics) == 0 || l.Topics[0] != ExecutionFailureTopic {
			continue
		}
		if (len(l.Topics) > 1 && bytes.Equal(l.Topics[1].Bytes(), safeTxHash)) || (len(l.Data) >= 32 && bytes.Equal(l.Data[:32], safeTxHash)) {
			return false
		}
	}
	return true
}

// SignedTx is an EIP-1559 transaction from the executor (signers[0]) to the
// Safe, signed and ready to broadcast.
func SignedTx(ctx context.Context, c *ethclient.Client, chainID *big.Int, executor config.Signer, safe common.Address, execData []byte, gasLimit uint64, tipMultiplier float64) (*types.Transaction, error) {
	nonce, err := c.PendingNonceAt(ctx, executor.Address)
	if err != nil {
		return nil, err
	}
	head, err := c.HeaderByNumber(ctx, nil)
	if err != nil {
		return nil, err
	}
	tip, err := c.SuggestGasTipCap(ctx)
	if err != nil {
		tip = big.NewInt(1_000_000) // 0.001 gwei
	}
	tip = mul(tip, tipMultiplier)
	base := big.NewInt(0)
	if head.BaseFee != nil {
		base = head.BaseFee
	}
	feeCap := new(big.Int).Add(new(big.Int).Mul(base, big.NewInt(2)), tip)
	tx := types.NewTx(&types.DynamicFeeTx{ChainID: chainID, Nonce: nonce, GasTipCap: tip, GasFeeCap: feeCap, Gas: gasLimit, To: &safe, Data: execData})
	return types.SignTx(tx, types.LatestSignerForChainID(chainID), executor.Key)
}

func mul(v *big.Int, f float64) *big.Int {
	r, _ := new(big.Float).Mul(new(big.Float).SetInt(v), big.NewFloat(f)).Int(nil)
	return r
}
