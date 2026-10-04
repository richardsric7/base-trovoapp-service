package network_test

import (
	"context"
	"math/big"
	"math/rand"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"trovo-wallet-api/internal/basetxn"
	"trovo-wallet-api/internal/evmkeypair"
	"trovo-wallet-api/internal/gnosissafe"
	"trovo-wallet-api/internal/network"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// TestPlatformKeyNoncesOnLocalChain uses one platform key from both send
// paths at once, as several instances would: legacy transactions signed
// now and broadcast later (basetxn, e.g. swap fee payments) and immediate
// sends (gnosissafe.SendTransaction). Every transaction must get its own
// nonce, and the legacy ones may be broadcast in any order. Skipped unless
// AA_LOCAL_CHAIN=1 with a hardhat node on 127.0.0.1:8545. It switches the
// node to interval mining while it runs, so run local-chain packages one at
// a time (go test -p 1).
func TestPlatformKeyNoncesOnLocalChain(t *testing.T) {
	if os.Getenv("AA_LOCAL_CHAIN") != "1" {
		t.Skip("AA_LOCAL_CHAIN not set")
	}
	t.Setenv("BASE_RPC_URL", "http://127.0.0.1:8545")
	t.Setenv("BASE_CHAIN_ID", "31337")
	ctx := context.Background()
	client, err := ethclient.Dial("http://127.0.0.1:8545")
	if err != nil {
		t.Fatal(err)
	}
	chainID, _ := client.ChainID(ctx)
	db, err := gorm.Open(sqlite.Open("file:"+filepath.Join(t.TempDir(), "n.db")+"?_busy_timeout=10000&_journal_mode=WAL"), &gorm.Config{Logger: logger.Discard})
	if err != nil {
		t.Fatal(err)
	}
	db.AutoMigrate(&sharedconfig.DistributedLock{}, &sharedconfig.NonceReservation{})
	network.SetDB(nil) // wires basetxn's builder
	gnosissafe.KeyLock = func(ctx context.Context, name string, fn func() error) error {
		return sharedconfig.WithKeyLock(db, "key:"+name, 3*time.Minute, fn)
	}
	gnosissafe.NonceSource = func(ctx context.Context, c *ethclient.Client, from common.Address) (uint64, func(), error) {
		return sharedconfig.NextNonceLocked(ctx, db, c, from)
	}
	network.ReserveNonce = func(ctx context.Context, from common.Address) (uint64, error) {
		return sharedconfig.ReserveNonce(ctx, db, client, from)
	}
	defer func() { gnosissafe.KeyLock, gnosissafe.NonceSource, network.ReserveNonce = nil, nil, nil }()

	// a funded platform key
	admin, _ := evmkeypair.ParseFull("0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80") // hardhat #0
	key, _ := evmkeypair.Random()
	keyAddr := common.HexToAddress(key.Address())
	{
		from := common.HexToAddress(admin.Address())
		n, _ := client.PendingNonceAt(ctx, from)
		price, _ := client.SuggestGasPrice(ctx)
		tx, _ := types.SignTx(types.NewTx(&types.LegacyTx{Nonce: n, To: &keyAddr, Value: big.NewInt(1e18), Gas: 21000, GasPrice: price}), types.NewEIP155Signer(chainID), admin.PrivateKey())
		if err := client.SendTransaction(ctx, tx); err != nil {
			t.Fatal(err)
		}
		if err := gnosissafe.WaitSuccess(ctx, client, tx.Hash()); err != nil {
			t.Fatal(err)
		}
	}

	// mine on an interval so gapped transactions can wait in the pool, as on
	// a real chain
	client.Client().CallContext(ctx, nil, "evm_setAutomine", false)
	client.Client().CallContext(ctx, nil, "evm_setIntervalMining", 300)
	defer func() {
		client.Client().CallContext(ctx, nil, "evm_setIntervalMining", 0)
		client.Client().CallContext(ctx, nil, "evm_setAutomine", true)
	}()

	payee, _ := evmkeypair.Random()
	other, _ := evmkeypair.Random()
	const legacy, immediate = 5, 3
	var mu sync.Mutex
	var signed []string
	var wg sync.WaitGroup
	for i := 0; i < legacy; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			tx, err := basetxn.NewTransaction(basetxn.TransactionParams{SourceAccount: key.Address(), Operations: []basetxn.Operation{
				basetxn.Payment{Destination: payee.Address(), Amount: "0.000001", Asset: basetxn.NativeAsset{}, SourceAccount: key.Address()},
			}})
			if err == nil {
				tx, err = tx.Sign("", key)
			}
			if err != nil {
				t.Error(err)
				return
			}
			raw, _ := tx.Base64()
			mu.Lock()
			signed = append(signed, raw)
			mu.Unlock()
		}()
	}
	for i := 0; i < immediate; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err := gnosissafe.SendTransaction(ctx, client, chainID, key, common.HexToAddress(other.Address()), nil); err != nil {
				t.Error(err)
			}
		}()
	}
	wg.Wait()
	if t.Failed() {
		return
	}
	// the legacy ones are broadcast later, in any order
	rand.Shuffle(len(signed), func(i, j int) { signed[i], signed[j] = signed[j], signed[i] })
	for _, raw := range signed {
		if _, err := network.SubmitXdrWithSignature(client, key.Address(), raw, ""); err != nil {
			t.Fatalf("broadcasting a legacy transaction: %v", err)
		}
	}
	deadline := time.Now().Add(30 * time.Second)
	for {
		n, _ := client.NonceAt(ctx, keyAddr, nil)
		if n == legacy+immediate {
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("key nonce %d, want %d mined", n, legacy+immediate)
		}
		time.Sleep(300 * time.Millisecond)
	}
	if bal, _ := client.BalanceAt(ctx, common.HexToAddress(payee.Address()), nil); bal.Cmp(big.NewInt(legacy*1e12)) != 0 {
		t.Fatalf("payee got %v, want %d", bal, int64(legacy)*1e12)
	}
}
