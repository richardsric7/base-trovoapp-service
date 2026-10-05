package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"math/big"
	"os"
	"strings"
	"time"

	"trovo-wallet-payment-history-engine/internal/network"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"gorm.io/gorm"
)

// Trovo wallets are Safe smart accounts that send through ERC-4337: a
// bundler's transaction calls the EntryPoint's handleOps, and each wallet's
// ETH leaves inside it - an internal call, invisible to a scan of top-level
// transactions. nativeOpTransfers finds those sends by decoding handleOps:
// each operation's Safe4337Module call (executeUserOp, or a delegatecall to
// MultiSendCallOnly for a batch), kept only when the operation succeeded
// (its UserOperationEvent). ETH a Safe receives from other contracts shows
// as the Safe's SafeReceived event (safeReceipts).

// defaultEntryPoint is the ERC-4337 v0.7 EntryPoint (ENTRYPOINT_ADDRESS
// overrides it).
const defaultEntryPoint = "0x0000000071727De22E5E9d8BAf0edAc6f37da032"

func entryPointAddress() common.Address {
	if a := strings.TrimSpace(os.Getenv("ENTRYPOINT_ADDRESS")); common.IsHexAddress(a) {
		return common.HexToAddress(a)
	}
	return common.HexToAddress(defaultEntryPoint)
}

var (
	userOperationEventSig = crypto.Keccak256Hash([]byte("UserOperationEvent(bytes32,address,address,uint256,bool,uint256,uint256)"))
	safeReceivedSig       = crypto.Keccak256Hash([]byte("SafeReceived(address,uint256)"))
)

var userOpABI = func() abi.ABI {
	a, err := abi.JSON(strings.NewReader(`[
 {"name":"handleOps","type":"function","inputs":[
   {"name":"ops","type":"tuple[]","components":[
     {"name":"sender","type":"address"},{"name":"nonce","type":"uint256"},{"name":"initCode","type":"bytes"},
     {"name":"callData","type":"bytes"},{"name":"accountGasLimits","type":"bytes32"},{"name":"preVerificationGas","type":"uint256"},
     {"name":"gasFees","type":"bytes32"},{"name":"paymasterAndData","type":"bytes"},{"name":"signature","type":"bytes"}]},
   {"name":"beneficiary","type":"address"}],"outputs":[]},
 {"name":"executeUserOp","type":"function","inputs":[
   {"name":"to","type":"address"},{"name":"value","type":"uint256"},{"name":"data","type":"bytes"},{"name":"operation","type":"uint8"}],"outputs":[]},
 {"name":"executeUserOpWithErrorString","type":"function","inputs":[
   {"name":"to","type":"address"},{"name":"value","type":"uint256"},{"name":"data","type":"bytes"},{"name":"operation","type":"uint8"}],"outputs":[]},
 {"name":"multiSend","type":"function","inputs":[{"name":"transactions","type":"bytes"}],"outputs":[]},
 {"name":"UserOperationEvent","type":"event","inputs":[
   {"name":"userOpHash","type":"bytes32","indexed":true},{"name":"sender","type":"address","indexed":true},
   {"name":"paymaster","type":"address","indexed":true},{"name":"nonce","type":"uint256","indexed":false},
   {"name":"success","type":"bool","indexed":false},{"name":"actualGasCost","type":"uint256","indexed":false},
   {"name":"actualGasUsed","type":"uint256","indexed":false}]}
]`))
	if err != nil {
		panic(err)
	}
	return a
}()

// opTransfer is ETH a wallet sent inside a user operation.
type opTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Nonce *big.Int
	Op    int // the operation's index in handleOps
	Call  int // the call's index in the operation
}

// opKey identifies a user operation: its wallet and nonce.
func opKey(sender common.Address, nonce *big.Int) string {
	return strings.ToLower(sender.Hex()) + ":" + nonce.String()
}

// succeededOps are the user operations that succeeded, from the
// EntryPoint's UserOperationEvent logs of one transaction.
func succeededOps(logs []types.Log, entryPoint common.Address) map[string]bool {
	ev := userOpABI.Events["UserOperationEvent"]
	out := map[string]bool{}
	for _, l := range logs {
		if l.Address != entryPoint || len(l.Topics) != 4 || l.Topics[0] != userOperationEventSig {
			continue
		}
		vals, err := ev.Inputs.NonIndexed().Unpack(l.Data)
		if err != nil || len(vals) < 2 {
			continue
		}
		if ok, _ := vals[1].(bool); ok {
			out[opKey(common.BytesToAddress(l.Topics[2].Bytes()), vals[0].(*big.Int))] = true
		}
	}
	return out
}

// nativeOpTransfers are the ETH sends of the succeeded operations in a
// handleOps call (input: the transaction's calldata).
func nativeOpTransfers(input []byte, succeeded map[string]bool) []opTransfer {
	if len(input) < 4 {
		return nil
	}
	m, err := userOpABI.MethodById(input[:4])
	if err != nil || m.Name != "handleOps" {
		return nil
	}
	args, err := m.Inputs.Unpack(input[4:])
	if err != nil || len(args) == 0 {
		return nil
	}
	ops, ok := args[0].([]struct {
		Sender             common.Address `json:"sender"`
		Nonce              *big.Int       `json:"nonce"`
		InitCode           []byte         `json:"initCode"`
		CallData           []byte         `json:"callData"`
		AccountGasLimits   [32]byte       `json:"accountGasLimits"`
		PreVerificationGas *big.Int       `json:"preVerificationGas"`
		GasFees            [32]byte       `json:"gasFees"`
		PaymasterAndData   []byte         `json:"paymasterAndData"`
		Signature          []byte         `json:"signature"`
	})
	if !ok {
		return nil
	}
	var out []opTransfer
	for i, op := range ops {
		if !succeeded[opKey(op.Sender, op.Nonce)] {
			continue
		}
		for j, c := range moduleCalls(op.CallData) {
			if c.value.Sign() > 0 {
				out = append(out, opTransfer{From: op.Sender, To: c.to, Value: c.value, Nonce: op.Nonce, Op: i, Call: j})
			}
		}
	}
	return out
}

type innerCall struct {
	to    common.Address
	value *big.Int
}

// moduleCalls are the calls a Safe4337Module callData makes: one call, or
// the calls of a MultiSend batch run by delegatecall.
func moduleCalls(callData []byte) []innerCall {
	if len(callData) < 4 {
		return nil
	}
	m, err := userOpABI.MethodById(callData[:4])
	if err != nil || (m.Name != "executeUserOp" && m.Name != "executeUserOpWithErrorString") {
		return nil
	}
	args, err := m.Inputs.Unpack(callData[4:])
	if err != nil || len(args) != 4 {
		return nil
	}
	to, _ := args[0].(common.Address)
	value, _ := args[1].(*big.Int)
	data, _ := args[2].([]byte)
	operation, _ := args[3].(uint8)
	if operation == 0 {
		if value == nil {
			value = new(big.Int)
		}
		return []innerCall{{to: to, value: value}}
	}
	// a delegatecall: only MultiSend batches move funds this way
	if len(data) < 4 {
		return nil
	}
	ms, err := userOpABI.MethodById(data[:4])
	if err != nil || ms.Name != "multiSend" {
		return nil
	}
	margs, err := ms.Inputs.Unpack(data[4:])
	if err != nil || len(margs) != 1 {
		return nil
	}
	packed, _ := margs[0].([]byte)
	return unpackMultiSend(packed)
}

// unpackMultiSend reads MultiSend's packed transactions: operation (1
// byte), to (20), value (32), data length (32), data.
func unpackMultiSend(b []byte) []innerCall {
	var out []innerCall
	for len(b) >= 85 {
		operation := b[0]
		to := common.BytesToAddress(b[1:21])
		value := new(big.Int).SetBytes(b[21:53])
		lenWord := b[53:85]
		if new(big.Int).SetBytes(lenWord[:24]).Sign() != 0 {
			return out // absurd length
		}
		n := binary.BigEndian.Uint64(lenWord[24:])
		if uint64(len(b)-85) < n {
			return out
		}
		if operation == 0 {
			out = append(out, innerCall{to: to, value: value})
		}
		b = b[85+n:]
	}
	return out
}

// safeReceipt is ETH a Safe received, from its SafeReceived event.
type safeReceipt struct {
	Safe   common.Address
	Sender common.Address
	Value  *big.Int
	Log    types.Log
}

// safeReceipts are a block's SafeReceived events. A Safe emits one for any
// ETH it receives through its fallback - including from a top-level
// transfer or another wallet's user operation, which the caller already
// records and skips.
func safeReceipts(logs []types.Log) []safeReceipt {
	var out []safeReceipt
	for _, l := range logs {
		if len(l.Topics) != 2 || l.Topics[0] != safeReceivedSig || len(l.Data) < 32 {
			continue
		}
		out = append(out, safeReceipt{Safe: l.Address, Sender: common.BytesToAddress(l.Topics[1].Bytes()), Value: new(big.Int).SetBytes(l.Data[len(l.Data)-32:]), Log: l})
	}
	return out
}

// nativeTransfer is one ETH transfer found in a transaction, with the id
// and sequence it is recorded under (the same whether found live or by a
// backfill, so it is never recorded twice).
type nativeTransfer struct {
	From, To common.Address
	Value    *big.Int
	ID, Seq  string
}

// nativeTransfersOfTx are the ETH transfers of one transaction (logs: its
// logs): its own value, the sends of its succeeded user operations, and
// ETH Safes received from contracts (SafeReceived) that neither covers.
func nativeTransfersOfTx(tx *types.Transaction, logs []types.Log, signer types.Signer, entryPoint common.Address) []nativeTransfer {
	var out []nativeTransfer
	recorded := map[string]bool{}
	add := func(t nativeTransfer) {
		recorded[strings.ToLower(t.From.Hex()+">"+t.To.Hex())] = true
		out = append(out, t)
	}
	hash := tx.Hash().Hex()
	if tx.To() != nil && *tx.To() == entryPoint {
		for _, t := range nativeOpTransfers(tx.Data(), succeededOps(logs, entryPoint)) {
			add(nativeTransfer{From: t.From, To: t.To, Value: t.Value, ID: fmt.Sprintf("%s-op%d-%d", hash, t.Op, t.Call), Seq: fmt.Sprintf("%v:%d:%d", t.Nonce, t.Op, t.Call)})
		}
	}
	if tx.To() != nil && tx.Value().Sign() > 0 {
		if from, err := types.Sender(signer, tx); err == nil {
			add(nativeTransfer{From: from, To: *tx.To(), Value: tx.Value(), ID: hash, Seq: fmt.Sprintf("%d", tx.Nonce())})
		}
	}
	for _, r := range safeReceipts(logs) {
		if r.Value.Sign() == 0 || recorded[strings.ToLower(r.Sender.Hex()+">"+r.Safe.Hex())] {
			continue
		}
		out = append(out, nativeTransfer{From: r.Sender, To: r.Safe, Value: r.Value, ID: fmt.Sprintf("%s-r%d", hash, r.Log.Index), Seq: fmt.Sprintf("r%d", r.Log.Index)})
	}
	return out
}

// backfillNativeOfSafe records the ETH a Safe received and sent in blocks
// [from, to] (blockTimes caches block times).
func backfillNativeOfSafe(client *ethclient.Client, safe common.Address, from, to uint64, blockTimes map[uint64]uint64, db, roachDB *gorm.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	entryPoint := entryPointAddress()
	safeTopic := common.BytesToHash(safe.Bytes())
	received, err := client.FilterLogs(ctx, ethereum.FilterQuery{
		FromBlock: new(big.Int).SetUint64(from), ToBlock: new(big.Int).SetUint64(to),
		Addresses: []common.Address{safe}, Topics: [][]common.Hash{{safeReceivedSig}},
	})
	if err != nil {
		return err
	}
	sent, err := client.FilterLogs(ctx, ethereum.FilterQuery{
		FromBlock: new(big.Int).SetUint64(from), ToBlock: new(big.Int).SetUint64(to),
		Addresses: []common.Address{entryPoint}, Topics: [][]common.Hash{{userOperationEventSig}, nil, {safeTopic}},
	})
	if err != nil {
		return err
	}
	signer := types.LatestSignerForChainID(network.GetBlockchainChainID())
	done := map[common.Hash]bool{}
	for _, l := range append(received, sent...) {
		if done[l.TxHash] {
			continue
		}
		done[l.TxHash] = true
		tx, _, err := client.TransactionByHash(ctx, l.TxHash)
		if err != nil {
			return err
		}
		rcpt, err := client.TransactionReceipt(ctx, l.TxHash)
		if err != nil {
			return err
		}
		ts, ok := blockTimes[l.BlockNumber]
		if !ok {
			h, err := client.HeaderByNumber(ctx, new(big.Int).SetUint64(l.BlockNumber))
			if err != nil {
				return err
			}
			ts = h.Time
			blockTimes[l.BlockNumber] = ts
		}
		logs := make([]types.Log, 0, len(rcpt.Logs))
		for _, rl := range rcpt.Logs {
			logs = append(logs, *rl)
		}
		for _, t := range nativeTransfersOfTx(tx, logs, signer, entryPoint) {
			if t.From != safe && t.To != safe {
				continue
			}
			recordNativeTransfer(t.From.Hex(), t.To.Hex(), t.Value, l.TxHash.Hex(), l.BlockNumber, t.ID, t.Seq, ts, db, roachDB)
		}
	}
	return nil
}
