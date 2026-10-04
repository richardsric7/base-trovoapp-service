package sharedconfig

import (
	"fmt"
	"os"
	"time"

	"gorm.io/gorm"
)

// DistributedLock backs a simple mutual-exclusion primitive that works
// identically on SQLite (this project's dev/test DB) and Postgres
// (production) - deliberately NOT built on Postgres advisory locks or
// `SELECT ... FOR UPDATE`, since SQLite supports neither. One row per
// named lock; claiming is a single atomic UPDATE statement, so callers
// never need to pin a connection or manage a transaction's lifetime the
// way session-scoped advisory locks require.
type DistributedLock struct {
	Name     string    `gorm:"primaryKey;size:100" json:"-"`
	LockedBy string    `gorm:"size:150;not null;default:''" json:"-"`
	LockedAt time.Time `json:"-"`
}

func (DistributedLock) TableName() string { return "distributed_locks" }

// InstanceIdentity is a human-readable label for whichever process holds
// a lock (hostname:pid) - used only for observability/debugging, never
// for correctness.
func InstanceIdentity() string {
	host, _ := os.Hostname()
	return fmt.Sprintf("%s:%d", host, os.Getpid())
}

// EnsureDistributedLock makes sure a lock's row exists. Idempotent and
// safe to call from every instance's boot - `ON CONFLICT DO NOTHING` is
// supported by both SQLite (3.24+) and Postgres.
func EnsureDistributedLock(db *gorm.DB, name string) error {
	return db.Exec(
		`INSERT INTO distributed_locks (name, locked_by, locked_at) VALUES (?, '', ?) ON CONFLICT (name) DO NOTHING`,
		name, time.Time{},
	).Error
}

// TryAcquireLock claims a named lock in one atomic UPDATE: it succeeds
// only if the row's WHERE clause still matches at write time, so two
// concurrent callers can never both succeed for the same name. A lock is
// claimable if it's unheld (LockedBy == "") or its holder has gone silent
// for longer than staleAfter - this covers a crashed holder that never
// released, the one thing dropping session-scoped advisory locks (which
// auto-release when their connection dies) costs us.
func TryAcquireLock(db *gorm.DB, name, holder string, staleAfter time.Duration) (bool, error) {
	now := time.Now()
	staleBefore := now.Add(-staleAfter)
	result := db.Exec(
		`UPDATE distributed_locks SET locked_by = ?, locked_at = ? WHERE name = ? AND (locked_by = '' OR locked_at < ?)`,
		holder, now, name, staleBefore,
	)
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected == 1, nil
}

// ReleaseLock releases a lock this holder currently owns. A no-op if some
// other holder has since taken it over via staleness - never errors on
// that, since losing a stale lock you no longer effectively hold isn't a
// failure.
func ReleaseLock(db *gorm.DB, name, holder string) error {
	return db.Exec(
		`UPDATE distributed_locks SET locked_by = '', locked_at = ? WHERE name = ? AND locked_by = ?`,
		time.Time{}, name, holder,
	).Error
}

// RenewLock refreshes a lock this holder still owns; false when it no
// longer does (another holder took it over as stale).
func RenewLock(db *gorm.DB, name, holder string) (bool, error) {
	result := db.Exec(
		`UPDATE distributed_locks SET locked_at = ? WHERE name = ? AND locked_by = ?`,
		time.Now(), name, holder,
	)
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected == 1, nil
}

// WaitAcquireLock polls TryAcquireLock (200ms backoff) until it succeeds
// or timeout elapses. Use this for a one-shot operation like a schema
// migration, where the caller should block briefly rather than skip the
// tick - contrast with WithSingletonLock (used by the periodic background
// sweeps), which never waits.
func WaitAcquireLock(db *gorm.DB, name, holder string, staleAfter, timeout time.Duration) (bool, error) {
	deadline := time.Now().Add(timeout)
	for {
		ok, err := TryAcquireLock(db, name, holder, staleAfter)
		if err != nil || ok {
			return ok, err
		}
		if time.Now().After(deadline) {
			return false, nil
		}
		time.Sleep(200 * time.Millisecond)
	}
}
