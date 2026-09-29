package aa

import (
	"context"
	"encoding/json"
	"math/big"
	"os"
	"testing"
	"time"

	"trovo-wallet-api/internal/evmkeypair"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// TestLocalStackEndToEnd drives the Builder against real contracts: a
// local node running paymaster/contracts/scripts/local-stack.js (EntryPoint,
// Safe + Safe4337Module, mock USDC, TrovoTokenPaymaster and a dev bundler)
// and a running paymaster quote service. It is skipped unless
//
//	AA_LOCAL_STACK=../paymaster/contracts/deployments/local-stack.json
//	AA_RPC_URL=http://127.0.0.1:8545          (default)
//	AA_QUOTE_SERVICE_URL=http://127.0.0.1:8090 AA_QUOTE_SERVICE_API_KEY=...
//
// are set. An undeployed, funded Safe activates paying gas in USDC, sends
// again (pre-charged), then pays gas in ETH.
func TestLocalStackEndToEnd(t *testing.T) {
	stack := os.Getenv("AA_LOCAL_STACK")
	if stack == "" {
		t.Skip("AA_LOCAL_STACK not set")
	}
	raw, err := os.ReadFile(stack)
	if err != nil {
		t.Fatal(err)
	}
	var s map[string]string
	if err := json.Unmarshal(raw, &s); err != nil {
		t.Fatal(err)
	}
	for _, k := range []string{"ENTRYPOINT_ADDRESS", "SAFE_PROXY_FACTORY_ADDRESS", "SAFE_SINGLETON_ADDRESS", "SAFE_MODULE_SETUP_ADDRESS", "SAFE_4337_MODULE_ADDRESS", "SAFE_MULTISEND_CALL_ONLY_ADDRESS"} {
		t.Setenv(k, s[k])
	}
	rpc := os.Getenv("AA_RPC_URL")
	if rpc == "" {
		rpc = "http://127.0.0.1:8545"
	}
	ctx := context.Background()
	client, err := ethclient.Dial(rpc)
	if err != nil {
		t.Fatal(err)
	}
	chainID, err := client.ChainID(ctx)
	if err != nil {
		t.Fatal(err)
	}
	cfg := ConfigFromEnv(chainID)
	usdc := common.HexToAddress(s["MOCK_USDC_ADDRESS"])

	owner, _ := evmkeypair.Random()
	ownerAddr := common.HexToAddress(owner.Address())
	safe := cfg.SafeAddress([]common.Address{ownerAddr}, 1, big.NewInt(0))

	// hardhat's test account #0 funds the Safe
	funder, _ := crypto.HexToECDSA("ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
	send := func(to common.Address, value *big.Int, data []byte) {
		t.Helper()
		from := crypto.PubkeyToAddress(funder.PublicKey)
		nonce, _ := client.PendingNonceAt(ctx, from)
		tip, _ := client.SuggestGasTipCap(ctx)
		head, _ := client.HeaderByNumber(ctx, nil)
		tx := types.NewTx(&types.DynamicFeeTx{ChainID: chainID, Nonce: nonce, To: &to, Value: value, Data: data, Gas: 200000,
			GasTipCap: tip, GasFeeCap: new(big.Int).Add(new(big.Int).Mul(head.BaseFee, big.NewInt(2)), tip)})
		signed, _ := types.SignTx(tx, types.LatestSignerForChainID(chainID), funder)
		if err := client.SendTransaction(ctx, signed); err != nil {
			t.Fatal(err)
		}
		for i := 0; i < 100; i++ {
			if r, _ := client.TransactionReceipt(ctx, signed.Hash()); r != nil {
				return
			}
			time.Sleep(100 * time.Millisecond)
		}
		t.Fatal("funding transaction not mined")
	}
	mint := mustABI(`[{"name":"mint","type":"function","inputs":[{"type":"address"},{"type":"uint256"}],"outputs":[]}]`)
	data, _ := mint.Pack("mint", safe, big.NewInt(50_000000))
	send(usdc, big.NewInt(0), data)

	b := &Builder{Config: cfg, Chain: client, Bundler: NewBundler(s["BUNDLER_URL"]), GasBufferPercent: 20}
	if url := os.Getenv("AA_QUOTE_SERVICE_URL"); url != "" {
		b.Quotes = NewQuoteClient(url, os.Getenv("AA_QUOTE_SERVICE_API_KEY"))
		b.Paymaster = common.HexToAddress(s["PAYMASTER_ADDRESS"])
	}
	w := Wallet{Address: safe, Owners: []common.Address{ownerAddr}, Threshold: 1,
		InitialOwners: []common.Address{ownerAddr}, InitialThreshold: 1, SaltNonce: big.NewInt(0)}
	recipient := common.HexToAddress("0x000000000000000000000000000000000000bEEF")
	before, _ := callView(ctx, client, usdc, "balanceOf", recipient)

	run := func(label string, req Request) {
		t.Helper()
		p, err := b.Prepare(ctx, req)
		if err != nil {
			t.Fatalf("%s: prepare: %v", label, err)
		}
		stored, _ := p.Marshal() // what the backend keeps between the two calls
		p, _ = UnmarshalPrepared(stored)
		sig, _ := evmkeypair.SignPersonal(owner.PrivateKey(), p.SafeOpHash.Bytes()) // the app's signBase64Txn
		hash, err := b.Submit(ctx, p, []OwnerSignature{{Owner: ownerAddr, Signature: sig}})
		if err != nil {
			t.Fatalf("%s: submit: %v", label, err)
		}
		if hash != p.UserOpHash {
			t.Fatalf("%s: bundler hash %s, prepared %s", label, hash.Hex(), p.UserOpHash.Hex())
		}
		wctx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()
		r, err := b.Bundler.(*Bundler).WaitReceipt(wctx, hash)
		if err != nil {
			t.Fatalf("%s: %v", label, err)
		}
		if !r.Success {
			t.Fatalf("%s: operation reverted", label)
		}
		t.Logf("%s: activation=%v gas in token=%v gas=%s wei tx=%s", label, p.Activation, p.GasToken != nil, r.ActualGasCost.ToInt(), r.Receipt.TransactionHash.Hex())
	}

	var gasToken *common.Address
	if b.Quotes != nil {
		gasToken = &usdc
	} else {
		send(safe, big.NewInt(1e17), nil) // no paymaster: activate paying ETH
	}
	run("activation", Request{Wallet: w, Calls: []Call{ERC20Transfer(usdc, recipient, big.NewInt(5_000000))}, GasToken: gasToken})
	if ok, _ := Deployed(ctx, client, safe); !ok {
		t.Fatal("the Safe was not deployed")
	}
	owners, threshold, err := OnchainOwners(ctx, client, safe)
	if err != nil || len(owners) != 1 || owners[0] != ownerAddr || threshold != 1 {
		t.Fatalf("on-chain owners %v/%d, %v", owners, threshold, err)
	}
	run("second transfer", Request{Wallet: w, Calls: []Call{ERC20Transfer(usdc, recipient, big.NewInt(1_000000))}, GasToken: gasToken})

	send(safe, big.NewInt(1e17), nil)
	run("ETH gas", Request{Wallet: w, Calls: []Call{NativeTransfer(recipient, big.NewInt(1e15))}})

	after, _ := callView(ctx, client, usdc, "balanceOf", recipient)
	if got := new(big.Int).Sub(after[0].(*big.Int), before[0].(*big.Int)); got.Int64() != 6_000000 {
		t.Fatalf("recipient received %v USDC base units, want 6000000", got)
	}
}
