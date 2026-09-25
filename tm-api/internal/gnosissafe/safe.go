// Package gnosissafe is tm-api's Base equivalent of the Stellar signer-weight
// management the vault-signer feature used to do via txnbuild.SetOptions.
// Base EOAs have no native multisig/weight concept, so wallets this feature
// manages are expected to be deployed Safe{Wallet} (formerly Gnosis Safe)
// contracts on Base - swapping/removing a signer means calling the Safe's
// own OwnerManager (swapOwner/removeOwner) through execTransaction, with
// enough of the Safe's current owners co-signing to meet its threshold.
//
// This package only wraps the handful of Safe v1.3/1.4 calls vault-signer
// needs (getOwners/getThreshold/nonce/getTransactionHash/execTransaction/
// swapOwner/removeOwner) - a minimal, stable, well-known subset of the Safe
// ABI, not a general Safe SDK.
package gnosissafe

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"sort"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// SentinelOwner is the linked-list head Safe's OwnerManager uses - the
// prevOwner value for swapOwner/removeOwner when the target owner is the
// first one returned by getOwners().
const SentinelOwner = "0x0000000000000000000000000000000000000001"

// operationCall is Safe's Enum.Operation.Call (0) - a plain external call,
// as opposed to delegatecall (1). swapOwner/removeOwner are always called
// this way, directly on the Safe itself.
const operationCall = uint8(0)

var safeABI abi.ABI

func init() {
	var err error
	safeABI, err = abi.JSON(strings.NewReader(`[
		{"inputs":[],"name":"getOwners","outputs":[{"internalType":"address[]","name":"","type":"address[]"}],"stateMutability":"view","type":"function"},
		{"inputs":[],"name":"getThreshold","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},
		{"inputs":[],"name":"nonce","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},
		{"inputs":[{"internalType":"address","name":"to","type":"address"},{"internalType":"uint256","name":"value","type":"uint256"},{"internalType":"bytes","name":"data","type":"bytes"},{"internalType":"uint8","name":"operation","type":"uint8"},{"internalType":"uint256","name":"safeTxGas","type":"uint256"},{"internalType":"uint256","name":"baseGas","type":"uint256"},{"internalType":"uint256","name":"gasPrice","type":"uint256"},{"internalType":"address","name":"gasToken","type":"address"},{"internalType":"address","name":"refundReceiver","type":"address"},{"internalType":"uint256","name":"_nonce","type":"uint256"}],"name":"getTransactionHash","outputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"stateMutability":"view","type":"function"},
		{"inputs":[{"internalType":"address","name":"to","type":"address"},{"internalType":"uint256","name":"value","type":"uint256"},{"internalType":"bytes","name":"data","type":"bytes"},{"internalType":"uint8","name":"operation","type":"uint8"},{"internalType":"uint256","name":"safeTxGas","type":"uint256"},{"internalType":"uint256","name":"baseGas","type":"uint256"},{"internalType":"uint256","name":"gasPrice","type":"uint256"},{"internalType":"address","name":"gasToken","type":"address"},{"internalType":"address","name":"refundReceiver","type":"address"},{"internalType":"bytes","name":"signatures","type":"bytes"}],"name":"execTransaction","outputs":[{"internalType":"bool","name":"success","type":"bool"}],"stateMutability":"payable","type":"function"},
		{"inputs":[{"internalType":"address","name":"prevOwner","type":"address"},{"internalType":"address","name":"oldOwner","type":"address"},{"internalType":"address","name":"newOwner","type":"address"}],"name":"swapOwner","outputs":[],"stateMutability":"nonpayable","type":"function"},
		{"inputs":[{"internalType":"address","name":"prevOwner","type":"address"},{"internalType":"address","name":"owner","type":"address"},{"internalType":"uint256","name":"_threshold","type":"uint256"}],"name":"removeOwner","outputs":[],"stateMutability":"nonpayable","type":"function"},
		{"inputs":[{"internalType":"address","name":"owner","type":"address"},{"internalType":"uint256","name":"_threshold","type":"uint256"}],"name":"addOwnerWithThreshold","outputs":[],"stateMutability":"nonpayable","type":"function"}
	]`))
	if err != nil {
		panic(fmt.Sprintf("[gnosissafe] invalid embedded Safe ABI: %v", err))
	}
}

// Safe wraps one deployed Safe contract's address for read/write calls.
type Safe struct {
	Address common.Address
	client  *ethclient.Client
}

// New wraps a deployed Safe's address, validating it's a well-formed Base
// address (not that a Safe is actually deployed there - the first call
// against it will fail clearly if not).
func New(client *ethclient.Client, safeAddress string) (*Safe, error) {
	if !common.IsHexAddress(safeAddress) {
		return nil, fmt.Errorf("invalid safe address %q", safeAddress)
	}
	return &Safe{Address: common.HexToAddress(safeAddress), client: client}, nil
}

// call ABI-encodes method(args...), performs an eth_call against this
// Safe, and returns its single unpacked return value.
func (s *Safe) call(ctx context.Context, method string, args ...interface{}) (interface{}, error) {
	data, err := safeABI.Pack(method, args...)
	if err != nil {
		return nil, fmt.Errorf("encoding %s call: %w", method, err)
	}
	result, err := s.client.CallContract(ctx, ethereum.CallMsg{To: &s.Address, Data: data}, nil)
	if err != nil {
		return nil, fmt.Errorf("calling %s: %w", method, err)
	}
	unpacked, err := safeABI.Unpack(method, result)
	if err != nil {
		return nil, fmt.Errorf("decoding %s result: %w", method, err)
	}
	return unpacked[0], nil
}

// Owners returns the Safe's current owner set, in the linked-list order
// getOwners() itself returns them.
func (s *Safe) Owners(ctx context.Context) ([]common.Address, error) {
	out, err := s.call(ctx, "getOwners")
	if err != nil {
		return nil, err
	}
	return out.([]common.Address), nil
}

// Threshold returns the number of owner signatures execTransaction
// currently requires.
func (s *Safe) Threshold(ctx context.Context) (*big.Int, error) {
	out, err := s.call(ctx, "getThreshold")
	if err != nil {
		return nil, err
	}
	return out.(*big.Int), nil
}

// Nonce returns the Safe's current transaction nonce (its own, unrelated
// to any EOA's account nonce - every execTransaction call, successful or
// not, consumes exactly one).
func (s *Safe) Nonce(ctx context.Context) (*big.Int, error) {
	out, err := s.call(ctx, "nonce")
	if err != nil {
		return nil, err
	}
	return out.(*big.Int), nil
}

// PrevOwner finds target's predecessor in the Safe's owner linked list -
// the prevOwner argument swapOwner/removeOwner both require. Returns
// SentinelOwner if target is the first owner returned by getOwners().
func PrevOwner(owners []common.Address, target common.Address) (common.Address, error) {
	prev := common.HexToAddress(SentinelOwner)
	for _, owner := range owners {
		if owner == target {
			return prev, nil
		}
		prev = owner
	}
	return common.Address{}, fmt.Errorf("%s is not a current owner of this safe", target.Hex())
}

// SwapOwnerCalldata ABI-encodes a swapOwner(prevOwner, oldOwner, newOwner)
// call - the data argument of the execTransaction that performs it.
func SwapOwnerCalldata(prevOwner, oldOwner, newOwner common.Address) ([]byte, error) {
	return safeABI.Pack("swapOwner", prevOwner, oldOwner, newOwner)
}

// RemoveOwnerCalldata ABI-encodes a removeOwner(prevOwner, owner,
// threshold) call.
func RemoveOwnerCalldata(prevOwner, owner common.Address, newThreshold *big.Int) ([]byte, error) {
	return safeABI.Pack("removeOwner", prevOwner, owner, newThreshold)
}

// AddOwnerCalldata ABI-encodes an addOwnerWithThreshold(owner, threshold)
// call - used to restore a previously-removed owner, since Safe has no
// "undo removeOwner" primitive: re-adding is a distinct call, at whatever
// threshold should hold once the owner is back.
func AddOwnerCalldata(owner common.Address, newThreshold *big.Int) ([]byte, error) {
	return safeABI.Pack("addOwnerWithThreshold", owner, newThreshold)
}

// TransactionHash calls the Safe's own getTransactionHash view function
// for a call to `to` with `data` at the Safe's current nonce - deliberately
// not computed locally (EIP-712 domain separator/typehash), so this always
// matches exactly what the Safe itself will check signatures against,
// whatever Safe version is actually deployed.
func (s *Safe) TransactionHash(ctx context.Context, to common.Address, data []byte, nonce *big.Int) ([32]byte, error) {
	out, err := s.call(ctx, "getTransactionHash",
		to, big.NewInt(0), data, operationCall,
		big.NewInt(0), big.NewInt(0), big.NewInt(0),
		common.Address{}, common.Address{}, nonce,
	)
	if err != nil {
		return [32]byte{}, err
	}
	return out.([32]byte), nil
}

// SignTransactionHash produces a Safe-format 65-byte (r,s,v) signature over
// safeTxHash - an EOA owner signing the RAW hash directly (Safe's
// checkNSignatures treats v in {27,28} as exactly that, no EIP-191/
// personal_sign prefix - unlike evmkeypair.Full.SignBase64, which is
// personal_sign and would produce a signature the Safe rejects).
func SignTransactionHash(priv *ecdsa.PrivateKey, safeTxHash [32]byte) ([]byte, error) {
	sig, err := crypto.Sign(safeTxHash[:], priv)
	if err != nil {
		return nil, err
	}
	sig[64] += 27
	return sig, nil
}

// ConcatSignatures orders per-owner signatures the way Safe's
// checkNSignatures requires: ascending by signer address, concatenated
// with no separator. sigs must already contain exactly one signature per
// distinct owner.
func ConcatSignatures(sigs map[common.Address][]byte) []byte {
	owners := make([]common.Address, 0, len(sigs))
	for owner := range sigs {
		owners = append(owners, owner)
	}
	sort.Slice(owners, func(i, j int) bool {
		return strings.ToLower(owners[i].Hex()) < strings.ToLower(owners[j].Hex())
	})
	out := make([]byte, 0, len(owners)*65)
	for _, owner := range owners {
		out = append(out, sigs[owner]...)
	}
	return out
}

// ExecTransactionCalldata ABI-encodes the execTransaction call that
// executes `data` against the Safe itself, with the given aggregated
// owner signatures.
func ExecTransactionCalldata(safeAddr, to common.Address, data []byte, signatures []byte) ([]byte, error) {
	return safeABI.Pack("execTransaction",
		to, big.NewInt(0), data, operationCall,
		big.NewInt(0), big.NewInt(0), big.NewInt(0),
		common.Address{}, common.Address{}, signatures,
	)
}

// ExecResult is what a submitted execTransaction resolved to once mined -
// TxHash always set, Success only meaningful once the receipt is in (a
// mined-but-reverted execTransaction, e.g. insufficient/invalid
// signatures, still has a real TxHash but Success=false).
type ExecResult struct {
	TxHash  string
	Success bool
}

// SubmitOwnerChange builds, signs (from senderKey), submits, and waits for
// the execTransaction that runs innerCalldata (a swapOwner/removeOwner
// call) against the Safe, with signatures already aggregated by the
// caller. senderKey pays gas and is the tx's msg.sender - execTransaction
// itself doesn't require the caller to be an owner, only that signatures
// satisfy the Safe's threshold. Waiting for the receipt (rather than just
// the broadcast succeeding) matches the old Stellar submission's
// synchronous success/failure semantics: SubmitTransactionXDR only
// returned once Horizon knew whether the ledger had actually applied it.
func (s *Safe) SubmitOwnerChange(ctx context.Context, senderKey *ecdsa.PrivateKey, innerCalldata []byte, signatures []byte, chainID *big.Int) (ExecResult, error) {
	execData, err := ExecTransactionCalldata(s.Address, s.Address, innerCalldata, signatures)
	if err != nil {
		return ExecResult{}, fmt.Errorf("encoding execTransaction: %w", err)
	}

	senderAddr := crypto.PubkeyToAddress(senderKey.PublicKey)
	nonce, err := s.client.PendingNonceAt(ctx, senderAddr)
	if err != nil {
		return ExecResult{}, fmt.Errorf("fetching sender nonce: %w", err)
	}
	gasTipCap, err := s.client.SuggestGasTipCap(ctx)
	if err != nil {
		return ExecResult{}, fmt.Errorf("fetching gas tip cap: %w", err)
	}
	head, err := s.client.HeaderByNumber(ctx, nil)
	if err != nil {
		return ExecResult{}, fmt.Errorf("fetching latest header: %w", err)
	}
	gasFeeCap := new(big.Int).Add(gasTipCap, new(big.Int).Mul(head.BaseFee, big.NewInt(2)))

	gasLimit, err := s.client.EstimateGas(ctx, ethereum.CallMsg{
		From: senderAddr,
		To:   &s.Address,
		Data: execData,
	})
	if err != nil {
		return ExecResult{}, fmt.Errorf("estimating gas: %w", err)
	}

	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		GasTipCap: gasTipCap,
		GasFeeCap: gasFeeCap,
		Gas:       gasLimit + gasLimit/5, // 20% headroom
		To:        &s.Address,
		Data:      execData,
	})

	signer := types.LatestSignerForChainID(chainID)
	signedTx, err := types.SignTx(tx, signer, senderKey)
	if err != nil {
		return ExecResult{}, fmt.Errorf("signing execTransaction: %w", err)
	}
	if err := s.client.SendTransaction(ctx, signedTx); err != nil {
		return ExecResult{}, fmt.Errorf("submitting execTransaction: %w", err)
	}

	receipt, err := waitMined(ctx, s.client, signedTx.Hash())
	if err != nil {
		return ExecResult{TxHash: signedTx.Hash().Hex()}, fmt.Errorf("waiting for execTransaction to be mined: %w", err)
	}
	return ExecResult{
		TxHash:  signedTx.Hash().Hex(),
		Success: receipt.Status == types.ReceiptStatusSuccessful,
	}, nil
}

// waitMined polls for txHash's receipt, the same way go-ethereum's own
// accounts/abi/bind.WaitMined does - kept as a small local helper rather
// than pulling in the bind package for this one call.
func waitMined(ctx context.Context, client *ethclient.Client, txHash common.Hash) (*types.Receipt, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Minute)
	defer cancel()
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	for {
		receipt, err := client.TransactionReceipt(ctx, txHash)
		if err == nil {
			return receipt, nil
		}
		if !errors.Is(err, ethereum.NotFound) {
			return nil, err
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-ticker.C:
		}
	}
}
