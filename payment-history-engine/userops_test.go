package main

import (
	"context"
	"encoding/json"
	"math/big"
	"os"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func pack(t *testing.T, method string, args ...interface{}) []byte {
	t.Helper()
	b, err := userOpABI.Pack(method, args...)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func multiSendTx(to common.Address, value int64, data []byte) []byte {
	b := []byte{0}
	b = append(b, to.Bytes()...)
	b = append(b, common.LeftPadBytes(big.NewInt(value).Bytes(), 32)...)
	b = append(b, common.LeftPadBytes(big.NewInt(int64(len(data))).Bytes(), 32)...)
	return append(b, data...)
}

type packedOp struct {
	Sender             common.Address `abi:"sender"`
	Nonce              *big.Int       `abi:"nonce"`
	InitCode           []byte         `abi:"initCode"`
	CallData           []byte         `abi:"callData"`
	AccountGasLimits   [32]byte       `abi:"accountGasLimits"`
	PreVerificationGas *big.Int       `abi:"preVerificationGas"`
	GasFees            [32]byte       `abi:"gasFees"`
	PaymasterAndData   []byte         `abi:"paymasterAndData"`
	Signature          []byte         `abi:"signature"`
}

func opEvent(t *testing.T, ep, sender common.Address, nonce int64, success bool) types.Log {
	data, err := userOpABI.Events["UserOperationEvent"].Inputs.NonIndexed().Pack(big.NewInt(nonce), success, big.NewInt(1), big.NewInt(1))
	if err != nil {
		t.Fatal(err)
	}
	return types.Log{Address: ep, Topics: []common.Hash{userOperationEventSig, {}, common.BytesToHash(sender.Bytes()), {}}, Data: data}
}

// TestUserOpETHSendsAreFound decodes a bundle the way Trovo wallets send
// ETH: a single executeUserOp transfer, a MultiSend batch (ETH to two
// recipients plus a token call), and a failed operation that must not
// count.
func TestUserOpETHSendsAreFound(t *testing.T) {
	ep := entryPointAddress()
	walletA := common.HexToAddress("0x00000000000000000000000000000000000000a1")
	walletB := common.HexToAddress("0x00000000000000000000000000000000000000b2")
	walletC := common.HexToAddress("0x00000000000000000000000000000000000000c3")
	bob := common.HexToAddress("0x0000000000000000000000000000000000000b0b")
	carol := common.HexToAddress("0x0000000000000000000000000000000000000ca1")
	multiSend := common.HexToAddress("0x9641d764fc13c8B624c04430C7356C1C7C8102e2")
	token := common.HexToAddress("0x00000000000000000000000000000000000000ee")

	single := pack(t, "executeUserOp", bob, big.NewInt(5e15), []byte{}, uint8(0))
	batch := append(multiSendTx(bob, 1e15, nil), multiSendTx(carol, 2e15, nil)...)
	batch = append(batch, multiSendTx(token, 0, []byte{0xa9, 0x05, 0x9c, 0xbb})...)
	batched := pack(t, "executeUserOp", multiSend, big.NewInt(0), pack(t, "multiSend", batch), uint8(1))
	failed := pack(t, "executeUserOp", bob, big.NewInt(9e15), []byte{}, uint8(0))

	ops := []packedOp{
		{Sender: walletA, Nonce: big.NewInt(1), CallData: single, PreVerificationGas: big.NewInt(0)},
		{Sender: walletB, Nonce: big.NewInt(7), CallData: batched, PreVerificationGas: big.NewInt(0)},
		{Sender: walletC, Nonce: big.NewInt(2), CallData: failed, PreVerificationGas: big.NewInt(0)},
	}
	input := pack(t, "handleOps", ops, common.HexToAddress("0x01"))
	logs := []types.Log{opEvent(t, ep, walletA, 1, true), opEvent(t, ep, walletB, 7, true), opEvent(t, ep, walletC, 2, false)}

	got := nativeOpTransfers(input, succeededOps(logs, ep))
	want := []struct {
		from, to common.Address
		value    int64
	}{{walletA, bob, 5e15}, {walletB, bob, 1e15}, {walletB, carol, 2e15}}
	if len(got) != len(want) {
		t.Fatalf("got %d transfers: %+v", len(got), got)
	}
	for i, w := range want {
		if got[i].From != w.from || got[i].To != w.to || got[i].Value.Int64() != w.value {
			t.Fatalf("transfer %d: got %+v, want %+v", i, got[i], w)
		}
	}
	// distinct ids within one transaction
	if got[1].Op != 1 || got[1].Call != 0 || got[2].Call != 1 {
		t.Fatalf("op/call indexes: %+v", got)
	}
}

func TestSafeReceivedIsRead(t *testing.T) {
	safe := common.HexToAddress("0x00000000000000000000000000000000000000a1")
	sender := common.HexToAddress("0x0000000000000000000000000000000000000b0b")
	l := types.Log{Address: safe, Topics: []common.Hash{safeReceivedSig, common.BytesToHash(sender.Bytes())}, Data: common.LeftPadBytes(big.NewInt(42).Bytes(), 32)}
	r := safeReceipts([]types.Log{l, {Address: safe, Topics: []common.Hash{transferEventSig}}})
	if len(r) != 1 || r[0].Safe != safe || r[0].Sender != sender || r[0].Value.Int64() != 42 {
		t.Fatalf("got %+v", r)
	}
}

// TestSafeETHSendsOnLocalChain decodes a real wallet's ETH batch send on
// the local chain: app-backend's TestSwapsOnLocalChain leaves one and
// writes "<tx hash> <wallet>" to AA_SAFE_ETH_TX_FILE. Skipped unless that
// file and AA_LOCAL_STACK are set.
func TestSafeETHSendsOnLocalChain(t *testing.T) {
	txFile, stackFile := os.Getenv("AA_SAFE_ETH_TX_FILE"), os.Getenv("AA_LOCAL_STACK")
	if txFile == "" || stackFile == "" {
		t.Skip("AA_SAFE_ETH_TX_FILE / AA_LOCAL_STACK not set")
	}
	raw, err := os.ReadFile(txFile)
	if err != nil {
		t.Fatal(err)
	}
	parts := strings.Fields(string(raw))
	stackRaw, _ := os.ReadFile(stackFile)
	stack := map[string]string{}
	json.Unmarshal(stackRaw, &stack)
	t.Setenv("ENTRYPOINT_ADDRESS", stack["ENTRYPOINT_ADDRESS"])

	ctx := context.Background()
	client, err := ethclient.Dial("http://127.0.0.1:8545")
	if err != nil {
		t.Fatal(err)
	}
	chainID, _ := client.ChainID(ctx)
	hash := common.HexToHash(parts[0])
	tx, _, err := client.TransactionByHash(ctx, hash)
	if err != nil {
		t.Fatal(err)
	}
	rcpt, err := client.TransactionReceipt(ctx, hash)
	if err != nil {
		t.Fatal(err)
	}
	var logs []types.Log
	for _, l := range rcpt.Logs {
		logs = append(logs, *l)
	}
	got := nativeTransfersOfTx(tx, logs, types.LatestSignerForChainID(chainID), entryPointAddress())
	wallet := common.HexToAddress(parts[1])
	if len(got) != 2 || got[0].From != wallet || got[0].To != common.HexToAddress("0x000000000000000000000000000000000000bEEF") || got[0].Value.Int64() != 1e15 ||
		got[1].To != common.HexToAddress("0x000000000000000000000000000000000000bEE2") || got[1].Value.Int64() != 2e15 || got[0].ID == got[1].ID {
		t.Fatalf("transfers: %+v", got)
	}
}
