package cache

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// ErrLockBusy is returned by WithLock when the lock stayed held by someone
// else for the whole wait.
var ErrLockBusy = errors.New("lock busy")

var (
	renewScript   = redis.NewScript(`if redis.call("get", KEYS[1]) == ARGV[1] then return redis.call("pexpire", KEYS[1], ARGV[2]) else return 0 end`)
	releaseScript = redis.NewScript(`if redis.call("get", KEYS[1]) == ARGV[1] then return redis.call("del", KEYS[1]) else return 0 end`)
)

// WithLock runs fn holding the named lock across every instance of the
// service, waiting up to wait for it. The lock expires after ttl unless
// renewed, and is renewed every ttl/3 while fn runs, so a crashed holder
// only blocks the others for ttl. Without Redis (single-instance
// development) fn just runs.
func (r *RedisCache) WithLock(ctx context.Context, name string, wait, ttl time.Duration, fn func() error) error {
	if r == nil || !r.Enabled || r.Client == nil {
		return fn()
	}
	var b [12]byte
	_, _ = rand.Read(b[:])
	token := hex.EncodeToString(b[:])
	deadline := time.Now().Add(wait)
	for {
		ok, err := r.Client.SetNX(ctx, name, token, ttl).Result()
		if err != nil {
			return err
		}
		if ok {
			break
		}
		if time.Now().After(deadline) {
			return ErrLockBusy
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(200 * time.Millisecond):
		}
	}
	stop := make(chan struct{})
	go func() {
		t := time.NewTicker(ttl / 3)
		defer t.Stop()
		for {
			select {
			case <-stop:
				return
			case <-t.C:
				if err := renewScript.Run(context.Background(), r.Client, []string{name}, token, ttl.Milliseconds()).Err(); err != nil {
					log.Printf("[lock:%v] renewal failed: %v", name, err)
				}
			}
		}
	}()
	defer func() {
		close(stop)
		if err := releaseScript.Run(context.Background(), r.Client, []string{name}, token).Err(); err != nil {
			log.Printf("[lock:%v] release failed: %v", name, err)
		}
	}()
	return fn()
}
