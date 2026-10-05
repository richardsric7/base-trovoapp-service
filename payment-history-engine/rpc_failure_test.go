package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/ethereum/go-ethereum/ethclient"
)

// rpcServer answers eth_call with symbol "USDC" / decimals 6, or with HTTP
// 502 while down is set.
func rpcServer(down *atomic.Bool) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if down.Load() {
			w.WriteHeader(http.StatusBadGateway)
			return
		}
		var req struct {
			ID     json.RawMessage `json:"id"`
			Params []struct {
				Data  string `json:"data"`
				Input string `json:"input"`
			} `json:"params"`
		}
		json.NewDecoder(r.Body).Decode(&req)
		data := req.Params[0].Data + req.Params[0].Input
		result := "0x0000000000000000000000000000000000000000000000000000000000000006" // decimals()
		if strings.HasPrefix(data, "0x95d89b41") {                                     // symbol()
			result = "0x0000000000000000000000000000000000000000000000000000000000000020" +
				"0000000000000000000000000000000000000000000000000000000000000004" +
				"5553444300000000000000000000000000000000000000000000000000000000"
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"jsonrpc":"2.0","id":` + string(req.ID) + `,"result":"` + result + `"}`))
	}))
}

func TestTokenMetaIsNotCachedWhileTheRPCFails(t *testing.T) {
	var down atomic.Bool
	srv := rpcServer(&down)
	defer srv.Close()
	client, err := ethclient.Dial(srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	const token = "0x00000000000000000000000000000000000000c1"

	down.Store(true)
	if m := getTokenMeta(client, token); m.decimals != 18 {
		t.Fatalf("fallback while down: %+v", m)
	}
	// once it answers, the real metadata is read (the fallback was not kept)
	down.Store(false)
	if m := getTokenMeta(client, token); m.decimals != 6 || m.symbol != "USDC" {
		t.Fatalf("after recovery: %+v", m)
	}
}

func TestProcessBlockReportsRPCFailure(t *testing.T) {
	var down atomic.Bool
	down.Store(true)
	srv := rpcServer(&down)
	defer srv.Close()
	client, _ := ethclient.Dial(srv.URL)
	// nothing is recorded (no database is touched) and the caller is told,
	// so it retries the block instead of moving past it
	if err := ProcessBlock(client, 123, nil, nil); err == nil {
		t.Fatal("expected an error while the RPC is down")
	}
}

func TestBackfillConcurrency(t *testing.T) {
	if backfillConcurrency() != 8 {
		t.Fatal("default")
	}
	t.Setenv("TRACK_ADDRESS_CONCURRENCY", "3")
	if backfillConcurrency() != 3 {
		t.Fatal("override")
	}
}
