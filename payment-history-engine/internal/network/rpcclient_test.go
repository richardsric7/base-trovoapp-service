package network

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

// an RPC answering eth_blockNumber, or stalling while stalled is set
func fakeRPC(stalled *atomic.Bool, hits *atomic.Int32) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits.Add(1)
		if stalled.Load() {
			select {
			case <-r.Context().Done():
			case <-time.After(5 * time.Second):
			}
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"jsonrpc":"2.0","id":1,"result":"0x10"}`))
	}))
}

func TestSharedRPCClientFailsFastAndRecovers(t *testing.T) {
	var stalled atomic.Bool
	var hits atomic.Int32
	srv := fakeRPC(&stalled, &hits)
	defer srv.Close()
	t.Setenv("BASE_RPC_URL", srv.URL)
	t.Setenv("RPC_TIMEOUT", "300ms")
	rpcBreakerCooldown = time.Second
	defer func() { rpcBreakerCooldown = 10 * time.Second }()

	c := GetBlockchainClient()
	if GetBlockchainClient() != c {
		t.Fatal("GetBlockchainClient must return one shared client per URL")
	}
	ctx := context.Background() // no deadline of its own: the client's applies
	if n, err := c.BlockNumber(ctx); err != nil || n != 16 {
		t.Fatalf("healthy RPC: %d, %v", n, err)
	}

	// the RPC stops answering: each call gives up after RPC_TIMEOUT...
	stalled.Store(true)
	for i := 0; i < rpcBreakerFailures; i++ {
		start := time.Now()
		if _, err := c.BlockNumber(ctx); err == nil || time.Since(start) > 2*time.Second {
			t.Fatalf("call %d against a stalled RPC: %v after %v", i, err, time.Since(start))
		}
	}
	// ...then calls fail at once, without reaching the RPC
	before := hits.Load()
	start := time.Now()
	for i := 0; i < 50; i++ {
		if _, err := c.BlockNumber(ctx); !errors.Is(err, ErrRPCUnavailable) && (err == nil || !strings.Contains(err.Error(), ErrRPCUnavailable.Error())) {
			t.Fatalf("expected ErrRPCUnavailable, got %v", err)
		}
	}
	if time.Since(start) > 200*time.Millisecond || hits.Load() != before {
		t.Fatalf("failing fast took %v and reached the RPC %d times", time.Since(start), hits.Load()-before)
	}

	// it answers again: after the cool-down calls go through
	stalled.Store(false)
	time.Sleep(1100 * time.Millisecond)
	if n, err := c.BlockNumber(ctx); err != nil || n != 16 {
		t.Fatalf("after recovery: %d, %v", n, err)
	}
}
