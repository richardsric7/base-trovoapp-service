package sharedconfig

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"gorm.io/gorm"
)

// Background work and shutdown: every lock-holding piece of work registers
// in backgroundWork, so a stopping instance (BeginShutdown) takes no new
// locks and can wait for what it is running to finish
// (WaitForBackgroundWork) before it exits.
var (
	shuttingDown   atomic.Bool
	backgroundWork sync.WaitGroup
)

// BeginShutdown stops this instance from starting new locked work.
func BeginShutdown() { shuttingDown.Store(true) }

// ShuttingDown reports whether BeginShutdown was called.
func ShuttingDown() bool { return shuttingDown.Load() }

// WaitForBackgroundWork waits up to timeout for running locked work to
// finish; it reports whether everything finished.
func WaitForBackgroundWork(timeout time.Duration) bool {
	done := make(chan struct{})
	go func() {
		backgroundWork.Wait()
		close(done)
	}()
	select {
	case <-done:
		return true
	case <-time.After(timeout):
		return false
	}
}

// lockHolder identifies one acquisition: the instance plus a random
// suffix, so two goroutines of the same instance never share a lock.
func lockHolder() string {
	var b [6]byte
	_, _ = rand.Read(b[:])
	return InstanceIdentity() + ":" + hex.EncodeToString(b[:])
}

// renewWhileRunning refreshes a held lock every staleAfter/3 until stop is
// closed, so work that runs longer than staleAfter keeps its lock; a crashed
// holder stops renewing and its lock goes stale.
func renewWhileRunning(db *gorm.DB, name, holder string, staleAfter time.Duration, stop <-chan struct{}) {
	every := staleAfter / 3
	if every < time.Second {
		every = time.Second
	}
	t := time.NewTicker(every)
	defer t.Stop()
	for {
		select {
		case <-stop:
			return
		case <-t.C:
			if ok, err := RenewLock(db, name, holder); err != nil || !ok {
				log.Printf("[lock:%v] renewal failed (held=%v): %v\n", name, ok, err)
			}
		}
	}
}

// runHolding runs fn while holding an acquired lock, renewing it, then
// releases it.
func runHolding(db *gorm.DB, name, holder string, staleAfter time.Duration, fn func()) {
	stop := make(chan struct{})
	go renewWhileRunning(db, name, holder, staleAfter, stop)
	defer func() {
		close(stop)
		if err := ReleaseLock(db, name, holder); err != nil {
			log.Printf("[lock:%v] failed to release: %v\n", name, err)
		}
	}()
	fn()
}

// WithSingletonLock runs fn only if this instance wins the named lock for
// this tick; every other replica running the same periodic loop just skips
// the tick and returns immediately. Non-blocking by design - these loops
// are idempotent periodic sweeps, not one-shot operations, so a losing
// replica should not wait, it should try again on its own next tick.
//
// Reuses the same distributed_locks table as MigrateDB (see
// distributed_lock.go) rather than a Postgres advisory lock, since this
// project's dev/test DB is SQLite.
//
// The lock is renewed while fn runs, so fn may take longer than
// staleAfter; staleAfter only bounds how long a crashed holder blocks the
// others. Waiting between ticks belongs outside fn (in the caller's loop),
// not inside it.
func WithSingletonLock(gc *GlobalConfig, lockName string, staleAfter time.Duration, fn func()) {
	if ShuttingDown() {
		return
	}
	if err := EnsureDistributedLock(gc.DB, lockName); err != nil {
		log.Printf("[WithSingletonLock:%v] failed to ensure lock row: %v\n", lockName, err)
		return
	}
	holder := lockHolder()
	acquired, err := TryAcquireLock(gc.DB, lockName, holder, staleAfter)
	if err != nil {
		log.Printf("[WithSingletonLock:%v] failed to acquire lock: %v\n", lockName, err)
		return
	}
	if !acquired {
		return
	}
	backgroundWork.Add(1)
	defer backgroundWork.Done()
	runHolding(gc.DB, lockName, holder, staleAfter, fn)
}

// ErrLockBusy is returned by WithKeyLock when the lock stayed held by
// someone else for the whole wait.
var ErrLockBusy = errors.New("lock busy")

// WithKeyLock runs fn holding the named lock across all instances, waiting
// up to wait for it (contrast WithSingletonLock, which skips). It is for
// work that must not overlap anywhere - e.g. everything between reading a
// signing key's nonce and its transaction being mined. The lock is renewed
// while fn runs.
func WithKeyLock(db *gorm.DB, lockName string, wait time.Duration, fn func() error) error {
	const staleAfter = 30 * time.Second
	if err := EnsureDistributedLock(db, lockName); err != nil {
		return err
	}
	holder := lockHolder()
	acquired, err := WaitAcquireLock(db, lockName, holder, staleAfter, wait)
	if err != nil {
		return err
	}
	if !acquired {
		return ErrLockBusy
	}
	backgroundWork.Add(1)
	defer backgroundWork.Done()
	var fnErr error
	runHolding(db, lockName, holder, staleAfter, func() { fnErr = fn() })
	return fnErr
}
