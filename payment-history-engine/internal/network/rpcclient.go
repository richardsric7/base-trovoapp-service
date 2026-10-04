package network

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

// The Base JSON-RPC client. One client per URL is shared by the whole
// process: a new client per call is harmless over HTTP but leaks a
// connection and its goroutines per call over a websocket URL.
//
// Over HTTP(S) every request has a deadline (RPC_TIMEOUT, default 30s), so
// a call made without one cannot hang forever on an RPC that stops
// answering, and the client fails fast while the RPC is down: after
// rpcBreakerFailures consecutive failed requests it refuses requests
// (ErrRPCUnavailable) for rpcBreakerCooldown, then lets requests through
// again - a single failure then reopens it. Requests during an outage thus
// fail at once instead of each holding a goroutine and its memory for the
// full timeout.

// ErrRPCUnavailable is returned (wrapped) by RPC calls while the client is
// failing fast.
var ErrRPCUnavailable = errors.New("blockchain RPC unavailable (recent requests failed); try again shortly")

const rpcBreakerFailures = 3

var rpcBreakerCooldown = 10 * time.Second

var (
	rpcClientsMu sync.Mutex
	rpcClients   = map[string]*ethclient.Client{}
)

// GetBlockchainClient returns the shared Base JSON-RPC client for
// BASE_RPC_URL, the Base equivalent of Stellar's Horizon client.
func GetBlockchainClient() *ethclient.Client {
	endpoint := os.Getenv("BASE_RPC_URL")
	rpcClientsMu.Lock()
	defer rpcClientsMu.Unlock()
	if c, ok := rpcClients[endpoint]; ok {
		return c
	}
	c, err := dialRPC(endpoint)
	if err != nil {
		log.Panicf("[GetBlockchainClient] invalid BASE_RPC_URL %q: %v", endpoint, err)
	}
	rpcClients[endpoint] = c
	return c
}

func rpcTimeout() time.Duration {
	if d, err := time.ParseDuration(strings.TrimSpace(os.Getenv("RPC_TIMEOUT"))); err == nil && d > 0 {
		return d
	}
	return 30 * time.Second
}

func dialRPC(endpoint string) (*ethclient.Client, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		// websocket / IPC: one long-lived connection, reused
		return ethclient.Dial(endpoint)
	}
	transport := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           (&net.Dialer{Timeout: 10 * time.Second, KeepAlive: 30 * time.Second}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   32,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: time.Second,
	}
	hc := &http.Client{Timeout: rpcTimeout(), Transport: newRPCBreaker(transport, endpoint)}
	rc, err := rpc.DialOptions(context.Background(), endpoint, rpc.WithHTTPClient(hc))
	if err != nil {
		return nil, err
	}
	return ethclient.NewClient(rc), nil
}

// rpcBreaker is the failing-fast http.RoundTripper described above.
type rpcBreaker struct {
	next     http.RoundTripper
	endpoint string

	mu        sync.Mutex
	failures  int
	openUntil time.Time
	tripped   bool // opened and not yet recovered
}

func newRPCBreaker(next http.RoundTripper, endpoint string) *rpcBreaker {
	return &rpcBreaker{next: next, endpoint: redactURL(endpoint)}
}

func (b *rpcBreaker) RoundTrip(req *http.Request) (*http.Response, error) {
	b.mu.Lock()
	if time.Now().Before(b.openUntil) {
		b.mu.Unlock()
		return nil, ErrRPCUnavailable
	}
	b.mu.Unlock()

	resp, err := b.next.RoundTrip(req)
	failed := (err != nil && !errors.Is(req.Context().Err(), context.Canceled)) || (err == nil && resp.StatusCode >= 500)
	b.record(failed)
	return resp, err
}

func (b *rpcBreaker) record(failed bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if !failed {
		if b.tripped {
			log.Printf("[rpc] %s is answering again", b.endpoint)
		}
		b.failures, b.tripped = 0, false
		return
	}
	b.failures++
	if b.tripped || b.failures >= rpcBreakerFailures {
		if !b.tripped {
			log.Printf("[rpc] %s failed %d times in a row: failing fast for %v at a time until it answers", b.endpoint, b.failures, rpcBreakerCooldown)
		}
		b.tripped = true
		b.failures = 0
		b.openUntil = time.Now().Add(rpcBreakerCooldown)
	}
}

func redactURL(endpoint string) string {
	u, err := url.Parse(endpoint)
	if err != nil {
		return "the RPC"
	}
	// provider URLs often carry an API key in the path or query
	return u.Scheme + "://" + u.Host
}
