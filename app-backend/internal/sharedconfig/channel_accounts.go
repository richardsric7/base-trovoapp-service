package sharedconfig

import (
	"errors"
	"fmt"
	"log"
	"time"
	"trovo-wallet-api/internal/evmkeypair"

	"gorm.io/gorm/clause"
)

// ChannelAccount tracks which blockchain signing key (channel account) is
// free to hand out. Deliberately address + status only, never key
// material - the actual signing key stays in the in-process
// GlobalConfig.ChannelAccountKeysByAddress map, loaded from CHANNEL_ACCOUNTS
// at boot on every instance, exactly like PendingAuth/FiatPaymentInvoice
// never store key material either. This table is what coordinates *which*
// address is free across instances; the map is purely a local cache of
// key material this instance already has in memory.
type ChannelAccount struct {
	Address   string    `gorm:"primaryKey;size:64" json:"-"`
	Status    string    `gorm:"size:20;not null;default:'available';index" json:"-"` // available | checked_out | in_use
	ClaimedBy string    `gorm:"size:150;not null;default:''" json:"-"`
	ClaimedAt time.Time `json:"-"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (ChannelAccount) TableName() string { return "channel_accounts" }

// SeedChannelAccount caches kp's signing key locally and ensures its
// channel_accounts row exists as available. Idempotent and safe to call
// from every instance's boot - clause.OnConflict{DoNothing:true} is GORM's
// portable upsert-guard, translated correctly for both SQLite (3.24+) and
// Postgres, so this never resets an already-claimed row back to available.
func SeedChannelAccount(gc *GlobalConfig, kp *evmkeypair.Full) error {
	if kp == nil {
		return nil
	}
	address := kp.Address()
	gc.ChannelAccountKeysMutex.Lock()
	if gc.ChannelAccountKeysByAddress == nil {
		gc.ChannelAccountKeysByAddress = make(map[string]*evmkeypair.Full)
	}
	gc.ChannelAccountKeysByAddress[address] = kp
	gc.ChannelAccountKeysMutex.Unlock()

	return gc.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&ChannelAccount{
		Address: address,
		Status:  "available",
	}).Error
}

func (gc *GlobalConfig) lookupChannelAccountKey(address string) (*evmkeypair.Full, bool) {
	gc.ChannelAccountKeysMutex.RLock()
	defer gc.ChannelAccountKeysMutex.RUnlock()
	k, ok := gc.ChannelAccountKeysByAddress[address]
	return k, ok
}

// CheckoutChannelAccount claims one available channel account across every
// instance sharing this database and returns its signing key plus a
// release closure - a drop-in replacement for the old
// `chanAccount := <-gc.ChannelAccounts; defer func(c){gc.ChannelAccounts<-c}(chanAccount)`
// pattern.
//
// Claiming is SELECT-candidate then UPDATE ... WHERE address = ? AND
// status = 'available', checked via RowsAffected, retried on collision -
// deliberately not `SELECT ... FOR UPDATE SKIP LOCKED` (Postgres/MySQL8+
// only, no SQLite equivalent). Two racing callers' SELECTs can return the
// same candidate address, but only one UPDATE wins that row; the loser's
// RowsAffected is 0 and it retries with a fresh SELECT. Under real
// contention this costs an occasional extra round trip, not a correctness
// gap - the UPDATE's WHERE-then-write is atomic per row on both SQLite and
// Postgres regardless of concurrent callers.
func CheckoutChannelAccount(gc *GlobalConfig) (*evmkeypair.Full, func(), error) {
	holder := InstanceIdentity()
	const maxAttempts = 8
	for attempt := 0; attempt < maxAttempts; attempt++ {
		var candidate string
		if err := gc.DB.Raw(`SELECT address FROM channel_accounts WHERE status = 'available' ORDER BY address LIMIT 1`).Scan(&candidate).Error; err != nil {
			return nil, nil, fmt.Errorf("checkout channel account: %w", err)
		}
		if candidate == "" {
			return nil, nil, errors.New("no channel account is currently available")
		}

		now := time.Now()
		result := gc.DB.Exec(
			`UPDATE channel_accounts SET status = 'checked_out', claimed_by = ?, claimed_at = ? WHERE address = ? AND status = 'available'`,
			holder, now, candidate,
		)
		if result.Error != nil {
			return nil, nil, fmt.Errorf("checkout channel account: %w", result.Error)
		}
		if result.RowsAffected != 1 {
			// Another caller claimed this exact row between our SELECT and
			// UPDATE - retry with a fresh candidate rather than fail.
			continue
		}

		key, ok := gc.lookupChannelAccountKey(candidate)
		if !ok {
			// Claimed a row this instance has no signing key cached for
			// (shouldn't happen if boot seeding ran correctly) - release
			// it back rather than leak a permanently checked-out row.
			gc.DB.Exec(`UPDATE channel_accounts SET status = 'available', claimed_by = '', claimed_at = ? WHERE address = ? AND claimed_by = ?`,
				time.Time{}, candidate, holder)
			return nil, nil, fmt.Errorf("checkout channel account: no signing key cached locally for address %v", candidate)
		}

		release := func() {
			if err := gc.DB.Exec(`UPDATE channel_accounts SET status = 'available', claimed_by = '', claimed_at = ? WHERE address = ? AND claimed_by = ?`,
				time.Time{}, candidate, holder).Error; err != nil {
				log.Printf("[CheckoutChannelAccount] failed to release %v: %v\n", candidate, err)
			}
		}
		return key, release, nil
	}
	return nil, nil, errors.New("checkout channel account: exhausted retries under contention")
}

// StoreInUseChannelAccount reserves kp for a longer-lived, out-of-band
// release (a pending on-chain transaction or fiat invoice tracked by its
// own row elsewhere) rather than CheckoutChannelAccount's short signing-op
// borrow. Keeps its exact pre-existing signature - callers
// (approvals.go, fiatPayment.go, boot-time reconciliation) are unchanged.
func (gc *GlobalConfig) StoreInUseChannelAccount(kp *evmkeypair.Full) {
	if kp == nil {
		return
	}
	address := kp.Address()
	gc.ChannelAccountKeysMutex.Lock()
	if gc.ChannelAccountKeysByAddress == nil {
		gc.ChannelAccountKeysByAddress = make(map[string]*evmkeypair.Full)
	}
	gc.ChannelAccountKeysByAddress[address] = kp
	gc.ChannelAccountKeysMutex.Unlock()

	now := time.Now()
	holder := InstanceIdentity()
	// This may be the first time this address's row exists - boot-time
	// in-use detection runs before SeedChannelAccount would otherwise
	// create it, so upsert rather than assume the row is already there.
	result := gc.DB.Exec(`UPDATE channel_accounts SET status = 'in_use', claimed_by = ?, claimed_at = ? WHERE address = ?`, holder, now, address)
	if result.Error != nil {
		log.Printf("[StoreInUseChannelAccount] failed to claim %v: %v\n", address, result.Error)
		return
	}
	if result.RowsAffected == 0 {
		if err := gc.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&ChannelAccount{
			Address: address, Status: "in_use", ClaimedBy: holder, ClaimedAt: now,
		}).Error; err != nil {
			log.Printf("[StoreInUseChannelAccount] failed to seed+claim %v: %v\n", address, err)
		}
	}
}

// ReleaseInUseChannelAccount frees an address reserved by
// StoreInUseChannelAccount. Keeps its exact pre-existing signature.
// Unconditional (no "only if this instance's map had it" guard the old
// in-memory version needed) - this also fixes a latent bug where a
// release webhook landing on a different instance than the one that
// originally checked the account in used to silently no-op.
func (gc *GlobalConfig) ReleaseInUseChannelAccount(pk string) {
	if len(pk) == 0 {
		return
	}
	if err := gc.DB.Exec(`UPDATE channel_accounts SET status = 'available', claimed_by = '', claimed_at = ? WHERE address = ?`, time.Time{}, pk).Error; err != nil {
		log.Printf("[ReleaseInUseChannelAccount] failed to release %v: %v\n", pk, err)
	}
}
