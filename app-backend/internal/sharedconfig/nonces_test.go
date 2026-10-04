package sharedconfig

import (
	"context"
	"path/filepath"
	"sort"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

// fakeChain's pending nonce only counts transactions "broadcast" so far.
type fakeChain struct{ pending atomic.Uint64 }

func (f *fakeChain) PendingNonceAt(context.Context, common.Address) (uint64, error) {
	return f.pending.Load(), nil
}

func TestNonceReservations(t *testing.T) {
	db := lockTestDB(t, filepath.Join(t.TempDir(), "nonces.db"))
	db.AutoMigrate(&NonceReservation{})
	chain := &fakeChain{}
	chain.pending.Store(7)
	key := common.HexToAddress("0x00000000000000000000000000000000000000aa")
	ctx := context.Background()

	// many signatures at once (none broadcast yet) get distinct, consecutive nonces
	var mu sync.Mutex
	var got []int
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			n, err := ReserveNonce(ctx, db, chain, key)
			if err != nil {
				t.Error(err)
				return
			}
			mu.Lock()
			got = append(got, int(n))
			mu.Unlock()
		}()
	}
	wg.Wait()
	sort.Ints(got)
	for i, n := range got {
		if n != 7+i {
			t.Fatalf("nonces %v, want 7..14", got)
		}
	}

	// the immediate-send path continues after them, and gives a nonce back
	// when its transaction was not sent
	n, release, err := NextNonceLocked(ctx, db, chain, key)
	if err != nil || n != 15 {
		t.Fatalf("got %d, %v", n, err)
	}
	release()
	if n, _ := ReserveNonce(ctx, db, chain, key); n != 15 {
		t.Fatalf("released nonce not reused: %d", n)
	}

	// once broadcast, the chain's pending nonce leads
	chain.pending.Store(30)
	if n, _ := ReserveNonce(ctx, db, chain, key); n != 30 {
		t.Fatalf("got %d, want the chain's 30", n)
	}

	// a reservation never sent stops counting after the TTL: the key goes
	// on from the chain's pending nonce instead of waiting behind the gap
	nonceReservationTTL = 200 * time.Millisecond
	defer func() { nonceReservationTTL = 10 * time.Minute }()
	ReserveNonce(ctx, db, chain, key) // 31, never sent
	time.Sleep(300 * time.Millisecond)
	if n, _ := ReserveNonce(ctx, db, chain, key); n != 30 {
		t.Fatalf("got %d, want 30 after the stale reservation", n)
	}
}
