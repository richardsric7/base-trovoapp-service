// Package monitor watches the paymaster's EntryPoint deposit and alerts
// the treasury when it runs low (the treasury then converts collected
// stablecoins to ETH and tops it up).
package monitor

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"sync"
	"time"

	"github.com/shopspring/decimal"
)

// DepositReader reads the current deposit (wei).
type DepositReader func(ctx context.Context) (*big.Int, error)

type Deposit struct {
	Read       DepositReader
	Low        *big.Int // alert below this (wei)
	WebhookURL string   // Discord/Slack-compatible incoming webhook; empty logs only
	Cooldown   time.Duration
	Label      string // e.g. "paymaster 0x… on chain 8453"
	HTTP       *http.Client

	mu        sync.RWMutex
	last      *big.Int
	checkedAt time.Time
	lastErr   error
	alertedAt time.Time
}

// Last returns the last successfully read deposit.
func (d *Deposit) Last() (*big.Int, bool) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	if d.last == nil {
		return nil, false
	}
	return new(big.Int).Set(d.last), true
}

// Status is the monitor's state for the status endpoint.
type Status struct {
	DepositWei   string    `json:"depositWei"`
	DepositETH   string    `json:"depositEth"`
	LowWatermark string    `json:"lowWatermarkEth"`
	Low          bool      `json:"low"`
	CheckedAt    time.Time `json:"checkedAt"`
	Error        string    `json:"error,omitempty"`
}

func (d *Deposit) Status() Status {
	d.mu.RLock()
	defer d.mu.RUnlock()
	s := Status{LowWatermark: eth(d.Low), CheckedAt: d.checkedAt}
	if d.last != nil {
		s.DepositWei = d.last.String()
		s.DepositETH = eth(d.last)
		s.Low = d.last.Cmp(d.Low) < 0
	}
	if d.lastErr != nil {
		s.Error = d.lastErr.Error()
	}
	return s
}

func eth(wei *big.Int) string {
	if wei == nil {
		return ""
	}
	return decimal.NewFromBigInt(wei, -18).String()
}

// Check reads the deposit once and alerts if it is low.
func (d *Deposit) Check(ctx context.Context) {
	v, err := d.Read(ctx)
	d.mu.Lock()
	d.checkedAt = time.Now()
	d.lastErr = err
	if err == nil {
		d.last = v
	}
	shouldAlert := err == nil && v.Cmp(d.Low) < 0 && time.Since(d.alertedAt) >= d.Cooldown
	if shouldAlert {
		d.alertedAt = time.Now()
	}
	d.mu.Unlock()

	if err != nil {
		log.Printf("[deposit] reading deposit: %v", err)
		return
	}
	if shouldAlert {
		msg := fmt.Sprintf("Paymaster deposit low: %s ETH (below %s ETH) - %s. Convert collected stablecoins to ETH and call deposit().", eth(v), eth(d.Low), d.Label)
		log.Printf("[deposit] %s", msg)
		d.alert(ctx, msg)
	}
}

func (d *Deposit) alert(ctx context.Context, msg string) {
	if d.WebhookURL == "" {
		return
	}
	body, _ := json.Marshal(map[string]string{"content": msg, "text": msg})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, d.WebhookURL, bytes.NewReader(body))
	if err != nil {
		log.Printf("[deposit] alert: %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	client := d.HTTP
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	}
	res, err := client.Do(req)
	if err != nil {
		log.Printf("[deposit] alert: %v", err)
		return
	}
	res.Body.Close()
	if res.StatusCode >= 300 {
		log.Printf("[deposit] alert webhook answered HTTP %d", res.StatusCode)
	}
}
