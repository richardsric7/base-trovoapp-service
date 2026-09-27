package sharedconfig

import (
	"log"
	"time"
)

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
// staleAfter should be a few multiples of the loop's own tick interval:
// long enough that a legitimately still-running fn isn't stolen from
// under itself, short enough to recover promptly if a holder crashed
// mid-run without releasing.
func WithSingletonLock(gc *GlobalConfig, lockName string, staleAfter time.Duration, fn func()) {
	if err := EnsureDistributedLock(gc.DB, lockName); err != nil {
		log.Printf("[WithSingletonLock:%v] failed to ensure lock row: %v\n", lockName, err)
		return
	}
	holder := InstanceIdentity()
	acquired, err := TryAcquireLock(gc.DB, lockName, holder, staleAfter)
	if err != nil {
		log.Printf("[WithSingletonLock:%v] failed to acquire lock: %v\n", lockName, err)
		return
	}
	if !acquired {
		return
	}
	defer func() {
		if err := ReleaseLock(gc.DB, lockName, holder); err != nil {
			log.Printf("[WithSingletonLock:%v] failed to release lock: %v\n", lockName, err)
		}
	}()
	fn()
}
