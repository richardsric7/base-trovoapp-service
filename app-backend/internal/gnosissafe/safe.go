// Package gnosissafe is app-backend's Safe (formerly Gnosis Safe) multisig
// client: EIP-712 SafeTx hashing, owner signatures, execTransaction
// submission (single call, or several calls batched atomically through
// Safe's MultiSendCallOnly), and deploying new Safe proxies through the
// canonical SafeProxyFactory.
//
// It is used by P2P escrow settlement (internal/components/p2p/safesigner)
// and by tokenized-asset issuance, where each asset's issuing wallet is a
// Safe that owns the asset's B20 token contract and mints through it.
package gnosissafe

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"math/big"
	"sort"
	"strings"
	"time"

	"trovo-wallet-api/internal/evmkeypair"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Safe operation types.
const (
	OperationCall         = uint8(0)
	OperationDelegateCall = uint8(1)
)

// DefaultMultiSendCallOnlyAddress is Safe's canonical MultiSendCallOnly
// v1.3.0 deployment, at the same address on Base and Base Sepolia (Safe's
// deterministic deployment). MultiSendCallOnly only permits plain `call`s
// for the batched sub-transactions, never a nested delegatecall.
const DefaultMultiSendCallOnlyAddress = "0x40A2aCCbd92BCA938b02010E17A5b8929b49130D"

// Call is one call the Safe should make.
type Call struct {
	To    common.Address
	Value *big.Int
	Data  []byte
}

var (
	safeABI      abi.ABI
	multiSendABI abi.ABI
	factoryABI   abi.ABI
	structTypes  abi.Arguments
)

func init() {
	var err error
	safeABI, err = abi.JSON(strings.NewReader(`[
		{"constant":true,"inputs":[],"name":"nonce","outputs":[{"name":"","type":"uint256"}],"type":"function"},
		{"constant":true,"inputs":[],"name":"getThreshold","outputs":[{"name":"","type":"uint256"}],"type":"function"},
		{"constant":true,"inputs":[],"name":"getOwners","outputs":[{"name":"","type":"address[]"}],"type":"function"},
		{"constant":false,"inputs":[
			{"name":"_owners","type":"address[]"},
			{"name":"_threshold","type":"uint256"},
			{"name":"to","type":"address"},
			{"name":"data","type":"bytes"},
			{"name":"fallbackHandler","type":"address"},
			{"name":"paymentToken","type":"address"},
			{"name":"payment","type":"uint256"},
			{"name":"paymentReceiver","type":"address"}
		],"name":"setup","outputs":[],"type":"function"},
		{"constant":false,"inputs":[
			{"name":"to","type":"address"},
			{"name":"value","type":"uint256"},
			{"name":"data","type":"bytes"},
			{"name":"operation","type":"uint8"},
			{"name":"safeTxGas","type":"uint256"},
			{"name":"baseGas","type":"uint256"},
			{"name":"gasPrice","type":"uint256"},
			{"name":"gasToken","type":"address"},
			{"name":"refundReceiver","type":"address"},
			{"name":"signatures","type":"bytes"}
		],"name":"execTransaction","outputs":[{"name":"","type":"bool"}],"type":"function"}
	]`))
	if err != nil {
		panic("gnosissafe: invalid safe ABI: " + err.Error())
	}
	multiSendABI, err = abi.JSON(strings.NewReader(`[
		{"inputs":[{"internalType":"bytes","name":"transactions","type":"bytes"}],"name":"multiSend","outputs":[],"stateMutability":"payable","type":"function"}
	]`))
	if err != nil {
		panic("gnosissafe: invalid multiSend ABI: " + err.Error())
	}
	factoryABI, err = abi.JSON(strings.NewReader(`[
		{"inputs":[],"name":"proxyCreationCode","outputs":[{"name":"","type":"bytes"}],"stateMutability":"pure","type":"function"},
		{"inputs":[{"name":"_singleton","type":"address"},{"name":"initializer","type":"bytes"},{"name":"saltNonce","type":"uint256"}],"name":"createProxyWithNonce","outputs":[{"name":"proxy","type":"address"}],"stateMutability":"nonpayable","type":"function"}
	]`))
	if err != nil {
		panic("gnosissafe: invalid proxy factory ABI: " + err.Error())
	}

	addressTy, _ := abi.NewType("address", "", nil)
	uint256Ty, _ := abi.NewType("uint256", "", nil)
	bytes32Ty, _ := abi.NewType("bytes32", "", nil)
	uint8Ty, _ := abi.NewType("uint8", "", nil)
	structTypes = abi.Arguments{
		{Type: bytes32Ty}, // SAFE_TX_TYPEHASH
		{Type: addressTy}, // to
		{Type: uint256Ty}, // value
		{Type: bytes32Ty}, // keccak256(data)
		{Type: uint8Ty},   // operation
		{Type: uint256Ty}, // safeTxGas
		{Type: uint256Ty}, // baseGas
		{Type: uint256Ty}, // gasPrice
		{Type: addressTy}, // gasToken
		{Type: addressTy}, // refundReceiver
		{Type: uint256Ty}, // nonce
	}
}

// SafeTxTypehash = keccak256("SafeTx(address to,uint256 value,bytes data,uint8 operation,uint256 safeTxGas,uint256 baseGas,uint256 gasPrice,address gasToken,address refundReceiver,uint256 nonce)")
var SafeTxTypehash = crypto.Keccak256([]byte("SafeTx(address to,uint256 value,bytes data,uint8 operation,uint256 safeTxGas,uint256 baseGas,uint256 gasPrice,address gasToken,address refundReceiver,uint256 nonce)"))

// DomainSeparatorTypehash = keccak256("EIP712Domain(uint256 chainId,address verifyingContract)")
var DomainSeparatorTypehash = crypto.Keccak256([]byte("EIP712Domain(uint256 chainId,address verifyingContract)"))

// ComputeSafeTxHash implements Safe's EIP-712 SafeTx typed-data hash
// exactly as the Safe contract verifies it (getTransactionHash), for a
// transaction with no gas refund (safeTxGas, baseGas and gasPrice all 0).
func ComputeSafeTxHash(chainID *big.Int, safe common.Address, to common.Address, value *big.Int, data []byte, operation uint8, nonce *big.Int) ([]byte, error) {
	addressTy, _ := abi.NewType("address", "", nil)
	uint256Ty, _ := abi.NewType("uint256", "", nil)
	bytes32Ty, _ := abi.NewType("bytes32", "", nil)
	domainArgs := abi.Arguments{{Type: bytes32Ty}, {Type: uint256Ty}, {Type: addressTy}}
	domainPacked, err := domainArgs.Pack(toBytes32(DomainSeparatorTypehash), chainID, safe)
	if err != nil {
		return nil, err
	}
	domainSeparator := crypto.Keccak256(domainPacked)

	zero := big.NewInt(0)
	structPacked, err := structTypes.Pack(
		toBytes32(SafeTxTypehash),
		to,
		value,
		toBytes32(crypto.Keccak256(data)),
		operation,
		zero, zero, zero,
		common.Address{}, common.Address{},
		nonce,
	)
	if err != nil {
		return nil, err
	}
	structHash := crypto.Keccak256(structPacked)

	// EIP-712: keccak256("\x19\x01" || domainSeparator || structHash)
	preimage := append([]byte{0x19, 0x01}, domainSeparator...)
	preimage = append(preimage, structHash...)
	return crypto.Keccak256(preimage), nil
}

// SignSafeTxHash produces Safe's "type 0" (EOA signature over the raw
// hash) signature for each signer, concatenated in ascending signer-address
// order as Safe's checkSignatures requires.
func SignSafeTxHash(safeTxHash []byte, signers []*evmkeypair.Full) ([]byte, error) {
	type sig struct {
		addr common.Address
		raw  []byte
	}
	sigs := make([]sig, 0, len(signers))
	for _, s := range signers {
		signature, err := crypto.Sign(safeTxHash, s.PrivateKey())
		if err != nil {
			return nil, err
		}
		// crypto.Sign returns V in {0,1}; Safe expects {27,28} for a
		// direct (non eth_sign-prefixed) hash signature.
		signature[64] += 27
		sigs = append(sigs, sig{addr: common.HexToAddress(s.Address()), raw: signature})
	}
	sort.Slice(sigs, func(i, j int) bool {
		return bytes.Compare(sigs[i].addr.Bytes(), sigs[j].addr.Bytes()) < 0
	})
	out := make([]byte, 0, 65*len(sigs))
	for _, s := range sigs {
		out = append(out, s.raw...)
	}
	return out, nil
}

// PackMultiSendTx encodes one sub-transaction in MultiSend's packed format:
// operation (1 byte) ++ to (20 bytes) ++ value (32 bytes) ++ data length
// (32 bytes) ++ data.
func PackMultiSendTx(operation uint8, to common.Address, value *big.Int, data []byte) []byte {
	buf := make([]byte, 0, 1+20+32+32+len(data))
	buf = append(buf, operation)
	buf = append(buf, to.Bytes()...)
	buf = append(buf, common.LeftPadBytes(value.Bytes(), 32)...)
	buf = append(buf, common.LeftPadBytes(big.NewInt(int64(len(data))).Bytes(), 32)...)
	buf = append(buf, data...)
	return buf
}

// Nonce reads the Safe's current transaction nonce.
func Nonce(ctx context.Context, client *ethclient.Client, safe common.Address) (*big.Int, error) {
	return callUint(ctx, client, safe, "nonce")
}

// Threshold reads the Safe's signature threshold.
func Threshold(ctx context.Context, client *ethclient.Client, safe common.Address) (*big.Int, error) {
	return callUint(ctx, client, safe, "getThreshold")
}

// Owners reads the Safe's owner list.
func Owners(ctx context.Context, client *ethclient.Client, safe common.Address) ([]common.Address, error) {
	out, err := call(ctx, client, safe, safeABI, "getOwners")
	if err != nil {
		return nil, err
	}
	owners, ok := out[0].([]common.Address)
	if !ok {
		return nil, errors.New("unexpected getOwners() result type")
	}
	return owners, nil
}

// HasCode reports whether a contract is deployed at addr.
func HasCode(ctx context.Context, client *ethclient.Client, addr common.Address) (bool, error) {
	code, err := client.CodeAt(ctx, addr, nil)
	if err != nil {
		return false, err
	}
	return len(code) > 0, nil
}

// SignersForSafe returns the subset of candidate keys that are owners of
// safe, limited to exactly the Safe's on-chain threshold - so a Safe whose
// owners or threshold were changed on-chain is signed for correctly, or
// fails clearly when the configured keys can no longer meet it.
func SignersForSafe(ctx context.Context, client *ethclient.Client, safe common.Address, candidates []*evmkeypair.Full) ([]*evmkeypair.Full, error) {
	owners, err := Owners(ctx, client, safe)
	if err != nil {
		return nil, fmt.Errorf("reading Safe owners: %w", err)
	}
	threshold, err := Threshold(ctx, client, safe)
	if err != nil {
		return nil, fmt.Errorf("reading Safe threshold: %w", err)
	}
	isOwner := map[common.Address]bool{}
	for _, o := range owners {
		isOwner[o] = true
	}
	need := int(threshold.Int64())
	out := make([]*evmkeypair.Full, 0, need)
	for _, kp := range candidates {
		if isOwner[common.HexToAddress(kp.Address())] {
			out = append(out, kp)
			if len(out) == need {
				return out, nil
			}
		}
	}
	return nil, fmt.Errorf("configured signers include %d owner(s) of Safe %s; its threshold is %d", len(out), safe.Hex(), need)
}

// ExecCalls signs and broadcasts calls as one Safe execTransaction,
// returning its hash without waiting for it to be mined - use WaitSuccess
// for that. Once a hash is returned the Safe nonce is spoken for: callers
// must never retry by calling ExecCalls again for the same intent, since
// that signs the next nonce and executes it twice. A single call is
// executed directly; several are batched through MultiSendCallOnly
// (delegatecall), so they all succeed or all revert together. signers must already satisfy the Safe's threshold
// (see SignersForSafe); signers[0] broadcasts and pays gas.
func ExecCalls(ctx context.Context, client *ethclient.Client, chainID *big.Int, safe common.Address, signers []*evmkeypair.Full, calls []Call, multiSendCallOnly common.Address) (string, error) {
	// The Safe's nonce only moves once its transaction is mined, so the
	// lock covers reading it through mining (callers' WaitSuccess then
	// returns at once).
	var hash string
	err := withKeyLock(ctx, "safe:"+strings.ToLower(safe.Hex()), func() error {
		var e error
		hash, e = execCalls(ctx, client, chainID, safe, signers, calls, multiSendCallOnly)
		if e != nil {
			return e
		}
		if e := waitMined(ctx, client, common.HexToHash(hash)); e != nil {
			log.Printf("[gnosissafe] %s on %s not mined before the lock was released: %v", hash, safe.Hex(), e)
		}
		return nil
	})
	return hash, err
}

func execCalls(ctx context.Context, client *ethclient.Client, chainID *big.Int, safe common.Address, signers []*evmkeypair.Full, calls []Call, multiSendCallOnly common.Address) (string, error) {
	if len(calls) == 0 {
		return "", errors.New("no calls to execute")
	}
	if len(signers) == 0 {
		return "", errors.New("no Safe signers supplied")
	}
	to, value, data, op := calls[0].To, calls[0].Value, calls[0].Data, OperationCall
	if value == nil {
		value = big.NewInt(0)
	}
	if len(calls) > 1 {
		ok, err := HasCode(ctx, client, multiSendCallOnly)
		if err != nil {
			return "", fmt.Errorf("checking MultiSendCallOnly code at %s: %w", multiSendCallOnly.Hex(), err)
		}
		if !ok {
			return "", fmt.Errorf("no contract code at MultiSendCallOnly address %s on this network", multiSendCallOnly.Hex())
		}
		var packed []byte
		for _, c := range calls {
			v := c.Value
			if v == nil {
				v = big.NewInt(0)
			}
			packed = append(packed, PackMultiSendTx(OperationCall, c.To, v, c.Data)...)
		}
		msData, err := multiSendABI.Pack("multiSend", packed)
		if err != nil {
			return "", fmt.Errorf("packing multiSend: %w", err)
		}
		to, value, data, op = multiSendCallOnly, big.NewInt(0), msData, OperationDelegateCall
	}

	nonce, err := Nonce(ctx, client, safe)
	if err != nil {
		return "", fmt.Errorf("fetching Safe nonce: %w", err)
	}
	safeTxHash, err := ComputeSafeTxHash(chainID, safe, to, value, data, op, nonce)
	if err != nil {
		return "", fmt.Errorf("computing safeTxHash: %w", err)
	}
	signatures, err := SignSafeTxHash(safeTxHash, signers)
	if err != nil {
		return "", fmt.Errorf("signing safeTxHash: %w", err)
	}
	zero := big.NewInt(0)
	execData, err := safeABI.Pack("execTransaction", to, value, data, op, zero, zero, zero, common.Address{}, common.Address{}, signatures)
	if err != nil {
		return "", fmt.Errorf("packing execTransaction: %w", err)
	}
	return SendTransaction(ctx, client, chainID, signers[0], safe, execData)
}

// DeployConfig names the canonical Safe contracts a new Safe proxy is
// created from.
type DeployConfig struct {
	ProxyFactory    common.Address
	Singleton       common.Address
	FallbackHandler common.Address
}

// PredictSafeAddress returns the CREATE2 address createProxyWithNonce
// will deploy owners/threshold/saltNonce's Safe to - so a retried
// deployment finds and reuses the Safe a previous attempt already created
// instead of reverting on the occupied address.
func PredictSafeAddress(ctx context.Context, client *ethclient.Client, cfg DeployConfig, owners []common.Address, threshold int64, saltNonce *big.Int) (common.Address, []byte, error) {
	initializer, err := safeABI.Pack("setup", owners, big.NewInt(threshold), common.Address{}, []byte{}, cfg.FallbackHandler, common.Address{}, big.NewInt(0), common.Address{})
	if err != nil {
		return common.Address{}, nil, fmt.Errorf("packing Safe setup: %w", err)
	}
	out, err := call(ctx, client, cfg.ProxyFactory, factoryABI, "proxyCreationCode")
	if err != nil {
		return common.Address{}, nil, fmt.Errorf("reading proxyCreationCode: %w", err)
	}
	creationCode, ok := out[0].([]byte)
	if !ok {
		return common.Address{}, nil, errors.New("unexpected proxyCreationCode() result type")
	}
	return create2SafeAddress(cfg.ProxyFactory, cfg.Singleton, creationCode, initializer, saltNonce), initializer, nil
}

// create2SafeAddress mirrors SafeProxyFactory.deployProxy's CREATE2
// derivation: salt = keccak256(keccak256(initializer) ++ saltNonce),
// init code = proxyCreationCode ++ uint256(singleton).
func create2SafeAddress(factory, singleton common.Address, creationCode, initializer []byte, saltNonce *big.Int) common.Address {
	salt := crypto.Keccak256(crypto.Keccak256(initializer), common.LeftPadBytes(saltNonce.Bytes(), 32))
	initCode := append(append([]byte{}, creationCode...), common.LeftPadBytes(singleton.Bytes(), 32)...)
	return crypto.CreateAddress2(factory, toBytes32(salt), crypto.Keccak256(initCode))
}

// DeploySafe creates a Safe proxy with the given owners and threshold
// through SafeProxyFactory.createProxyWithNonce, paid for by deployer, and
// waits for it to be mined. Idempotent per saltNonce: when the Safe
// already exists at its predicted address, that address is returned
// without sending a transaction.
func DeploySafe(ctx context.Context, client *ethclient.Client, chainID *big.Int, deployer *evmkeypair.Full, cfg DeployConfig, owners []common.Address, threshold int64, saltNonce *big.Int) (common.Address, error) {
	if threshold < 1 || int(threshold) > len(owners) {
		return common.Address{}, fmt.Errorf("invalid Safe threshold %d for %d owner(s)", threshold, len(owners))
	}
	for name, addr := range map[string]common.Address{"SafeProxyFactory": cfg.ProxyFactory, "Safe singleton": cfg.Singleton, "fallback handler": cfg.FallbackHandler} {
		ok, err := HasCode(ctx, client, addr)
		if err != nil {
			return common.Address{}, fmt.Errorf("checking %s code at %s: %w", name, addr.Hex(), err)
		}
		if !ok {
			return common.Address{}, fmt.Errorf("no %s contract deployed at %s on this network", name, addr.Hex())
		}
	}

	predicted, initializer, err := PredictSafeAddress(ctx, client, cfg, owners, threshold, saltNonce)
	if err != nil {
		return common.Address{}, err
	}
	if exists, err := HasCode(ctx, client, predicted); err != nil {
		return common.Address{}, err
	} else if exists {
		return predicted, nil
	}

	data, err := factoryABI.Pack("createProxyWithNonce", cfg.Singleton, initializer, saltNonce)
	if err != nil {
		return common.Address{}, fmt.Errorf("packing createProxyWithNonce: %w", err)
	}
	hash, err := SendTransaction(ctx, client, chainID, deployer, cfg.ProxyFactory, data)
	if err != nil {
		return common.Address{}, fmt.Errorf("sending Safe deployment: %w", err)
	}
	if err := WaitSuccess(ctx, client, common.HexToHash(hash)); err != nil {
		return common.Address{}, fmt.Errorf("Safe deployment %s: %w", hash, err)
	}
	if exists, err := HasCode(ctx, client, predicted); err != nil || !exists {
		return common.Address{}, fmt.Errorf("Safe deployment %s mined but no Safe found at predicted address %s", hash, predicted.Hex())
	}
	return predicted, nil
}

// SendTransaction signs data as a call from `from` to `to` and broadcasts
// it, returning the transaction hash. Gas is estimated with a 20% margin.
func SendTransaction(ctx context.Context, client *ethclient.Client, chainID *big.Int, from *evmkeypair.Full, to common.Address, data []byte) (string, error) {
	// the pending nonce counts a transaction once it is broadcast, so the
	// lock covers reading it through broadcasting
	var hash string
	err := withKeyLock(ctx, "eoa:"+strings.ToLower(from.Address()), func() error {
		var e error
		hash, e = sendTransaction(ctx, client, chainID, from, to, data)
		return e
	})
	return hash, err
}

func sendTransaction(ctx context.Context, client *ethclient.Client, chainID *big.Int, from *evmkeypair.Full, to common.Address, data []byte) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	fromAddr := common.HexToAddress(from.Address())
	nonce, err := client.PendingNonceAt(ctx, fromAddr)
	if err != nil {
		return "", err
	}
	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		return "", err
	}
	gasLimit, err := client.EstimateGas(ctx, ethereum.CallMsg{From: fromAddr, To: &to, Data: data})
	if err != nil {
		// A failing estimate usually means the call itself would revert;
		// surface that instead of burning gas on a doomed transaction.
		return "", fmt.Errorf("estimating gas: %w", err)
	}
	tx := types.NewTx(&types.LegacyTx{Nonce: nonce, To: &to, Value: big.NewInt(0), Gas: gasLimit * 12 / 10, GasPrice: gasPrice, Data: data})
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), from.PrivateKey())
	if err != nil {
		return "", err
	}
	if err := client.SendTransaction(ctx, signedTx); err != nil {
		return "", err
	}
	return signedTx.Hash().Hex(), nil
}

// KeyLock, when set, runs fn holding the named lock across every instance
// of the service (app-backend sets it at startup; see
// sharedconfig.WithKeyLock). ExecCalls takes "safe:<address>" and
// SendTransaction "eoa:<address>", so two instances never use the same
// Safe or signing key nonce at once. Unset (tests, tools), nothing is
// locked.
var KeyLock func(ctx context.Context, name string, fn func() error) error

func withKeyLock(ctx context.Context, name string, fn func() error) error {
	if KeyLock == nil {
		return fn()
	}
	return KeyLock(ctx, name, fn)
}

// waitMined waits (up to 2 minutes) for hash to be mined, whatever its
// outcome.
func waitMined(ctx context.Context, client *ethclient.Client, hash common.Hash) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		if _, err := client.TransactionReceipt(ctx, hash); err == nil {
			return nil
		} else if !errors.Is(err, ethereum.NotFound) {
			return err
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}
	}
}

// WaitSuccess polls for hash's receipt (up to 2 minutes) and returns an
// error if it reverted.
func WaitSuccess(ctx context.Context, client *ethclient.Client, hash common.Hash) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	for {
		receipt, err := client.TransactionReceipt(ctx, hash)
		if err == nil {
			if receipt.Status != types.ReceiptStatusSuccessful {
				return fmt.Errorf("transaction %s reverted", hash.Hex())
			}
			return nil
		}
		if !errors.Is(err, ethereum.NotFound) {
			return err
		}
		select {
		case <-ctx.Done():
			return fmt.Errorf("timed out waiting for transaction %s to be mined", hash.Hex())
		case <-ticker.C:
		}
	}
}

func callUint(ctx context.Context, client *ethclient.Client, contract common.Address, method string) (*big.Int, error) {
	out, err := call(ctx, client, contract, safeABI, method)
	if err != nil {
		return nil, err
	}
	v, ok := out[0].(*big.Int)
	if !ok {
		return nil, fmt.Errorf("unexpected %s() result type", method)
	}
	return v, nil
}

func call(ctx context.Context, client *ethclient.Client, contract common.Address, a abi.ABI, method string, args ...interface{}) ([]interface{}, error) {
	data, err := a.Pack(method, args...)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	result, err := client.CallContract(ctx, ethereum.CallMsg{To: &contract, Data: data}, nil)
	if err != nil {
		return nil, err
	}
	out, err := a.Unpack(method, result)
	if err != nil || len(out) == 0 {
		return nil, fmt.Errorf("could not decode %s() result", method)
	}
	return out, nil
}

func toBytes32(b []byte) [32]byte {
	var out [32]byte
	copy(out[:], b)
	return out
}
