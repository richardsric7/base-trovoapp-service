package gnosissafe_test

import (
	"context"
	"encoding/json"
	"math/big"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"trovo-wallet-api/internal/evmkeypair"
	"trovo-wallet-api/internal/gnosissafe"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// TestConcurrentSafeExecutionsOnLocalChain sends several transactions from
// one platform Safe at once, as several instances serving requests would
// (e.g. P2P settlements from the escrow Safe). Each would otherwise read the
// same Safe nonce and all but one would fail; with KeyLock (a database lock
// shared by the "instances") they run one after another and all succeed.
// Skipped unless AA_LOCAL_STACK points at paymaster/contracts' local stack.
func TestConcurrentSafeExecutionsOnLocalChain(t *testing.T) {
	stackFile := os.Getenv("AA_LOCAL_STACK")
	if stackFile == "" {
		t.Skip("AA_LOCAL_STACK not set")
	}
	raw, err := os.ReadFile(stackFile)
	if err != nil {
		t.Fatal(err)
	}
	stack := map[string]string{}
	json.Unmarshal(raw, &stack)
	ctx := context.Background()
	client, err := ethclient.Dial("http://127.0.0.1:8545")
	if err != nil {
		t.Fatal(err)
	}
	chainID, _ := client.ChainID(ctx)

	db, err := gorm.Open(sqlite.Open("file:"+filepath.Join(t.TempDir(), "locks.db")+"?_busy_timeout=10000&_journal_mode=WAL"), &gorm.Config{Logger: logger.Discard})
	if err != nil {
		t.Fatal(err)
	}
	db.AutoMigrate(&sharedconfig.DistributedLock{})
	gnosissafe.KeyLock = func(ctx context.Context, name string, fn func() error) error {
		return sharedconfig.WithKeyLock(db, "key:"+name, 3*time.Minute, fn)
	}
	defer func() { gnosissafe.KeyLock = nil }()

	admin, _ := evmkeypair.ParseFull("0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80") // hardhat #0
	owner, _ := evmkeypair.Random()
	pay := func(to common.Address) {
		t.Helper()
		from := common.HexToAddress(admin.Address())
		nonce, _ := client.PendingNonceAt(ctx, from)
		price, _ := client.SuggestGasPrice(ctx)
		tx, _ := types.SignTx(types.NewTx(&types.LegacyTx{Nonce: nonce, To: &to, Value: big.NewInt(1e18), Gas: 100000, GasPrice: price}), types.NewEIP155Signer(chainID), admin.PrivateKey())
		if err := client.SendTransaction(ctx, tx); err != nil {
			t.Fatal(err)
		}
		if err := gnosissafe.WaitSuccess(ctx, client, tx.Hash()); err != nil {
			t.Fatal(err)
		}
	}
	pay(common.HexToAddress(owner.Address()))
	safe, err := gnosissafe.DeploySafe(ctx, client, chainID, admin, gnosissafe.DeployConfig{
		ProxyFactory:    common.HexToAddress(stack["SAFE_PROXY_FACTORY_ADDRESS"]),
		Singleton:       common.HexToAddress(stack["SAFE_SINGLETON_ADDRESS"]),
		FallbackHandler: common.HexToAddress(stack["SAFE_4337_MODULE_ADDRESS"]), // any contract
	}, []common.Address{common.HexToAddress(owner.Address())}, 1, big.NewInt(time.Now().UnixNano()))
	if err != nil {
		t.Fatal(err)
	}
	pay(safe)

	recipient, _ := evmkeypair.Random()
	to := common.HexToAddress(recipient.Address())
	const n = 5
	var wg sync.WaitGroup
	errs := make(chan error, n)
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			h, err := gnosissafe.ExecCalls(ctx, client, chainID, safe, []*evmkeypair.Full{owner}, []gnosissafe.Call{{To: to, Value: big.NewInt(1000)}}, common.HexToAddress(stack["SAFE_MULTISEND_CALL_ONLY_ADDRESS"]))
			if err == nil {
				err = gnosissafe.WaitSuccess(ctx, client, common.HexToHash(h))
			}
			errs <- err
		}()
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		if err != nil {
			t.Fatalf("a concurrent execution failed: %v", err)
		}
	}
	bal, _ := client.BalanceAt(ctx, to, nil)
	if bal.Cmp(big.NewInt(n*1000)) != 0 {
		t.Fatalf("recipient got %v, want %d", bal, n*1000)
	}
	nonce, _ := gnosissafe.Nonce(ctx, client, safe)
	if nonce.Int64() != n {
		t.Fatalf("Safe nonce %v, want %d", nonce, n)
	}
}
