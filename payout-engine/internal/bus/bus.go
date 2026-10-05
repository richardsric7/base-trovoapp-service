// Package bus is payout-engine's Redis link to tm-api: tm-api publishes a
// wake-up on CommandsChannel after changing a payout in the database, and
// the engine publishes its progress on EventsChannel and a heartbeat under
// HeartbeatKey. The database stays the source of truth: without Redis the
// engine still polls it, and tm-api still reads the state from it.
package bus

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	CommandsChannel = "payout-engine:commands"
	EventsChannel   = "payout-engine:events"
	HeartbeatKey    = "payout-engine:heartbeat"
)

// Event is what the engine publishes.
type Event struct {
	Type     string    `json:"type"` // e.g. preparing, locked, batch-mined, completed, error
	PayoutID uint64    `json:"payoutId,omitempty"`
	Message  string    `json:"message"`
	At       time.Time `json:"at"`
}

// Bus is a Redis connection, or a no-op when Redis is not configured.
type Bus struct{ rdb *redis.Client }

// New connects when addr is set.
func New(addr, password string, useTLS bool) *Bus {
	if addr == "" {
		return &Bus{}
	}
	opts := &redis.Options{Addr: addr, Password: password}
	if useTLS {
		opts.TLSConfig = &tls.Config{MinVersion: tls.VersionTLS12}
	}
	return &Bus{rdb: redis.NewClient(opts)}
}

// Wakeups delivers a signal for every command tm-api publishes (the
// message itself is only a hint; the engine re-reads the database).
func (b *Bus) Wakeups(ctx context.Context) <-chan struct{} {
	out := make(chan struct{}, 1)
	if b.rdb == nil {
		return out
	}
	go func() {
		for ctx.Err() == nil {
			sub := b.rdb.Subscribe(ctx, CommandsChannel)
			for msg := range sub.Channel() {
				_ = msg
				select {
				case out <- struct{}{}:
				default:
				}
			}
			sub.Close()
			select {
			case <-ctx.Done():
			case <-time.After(5 * time.Second): // reconnect
			}
		}
	}()
	return out
}

// Publish sends an event (best effort).
func (b *Bus) Publish(ctx context.Context, e Event) {
	if b.rdb == nil {
		return
	}
	e.At = time.Now().UTC()
	raw, _ := json.Marshal(e)
	if err := b.rdb.Publish(ctx, EventsChannel, raw).Err(); err != nil {
		log.Printf("[bus] publish: %v", err)
	}
}

// Heartbeat records that the engine is alive (expires after ttl).
func (b *Bus) Heartbeat(ctx context.Context, status string, ttl time.Duration) {
	if b.rdb == nil {
		return
	}
	_ = b.rdb.Set(ctx, HeartbeatKey, status, ttl).Err()
}

// Close closes the connection.
func (b *Bus) Close() {
	if b.rdb != nil {
		b.rdb.Close()
	}
}
