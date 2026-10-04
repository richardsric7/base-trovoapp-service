package sharedconfig

import (
	"context"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"gorm.io/gorm"
)

// Nonces of platform signing keys. A transaction that is signed now and
// sent later (the legacy basetxn path: signed by the platform, then sent
// after the user's or approvers' signatures) cannot hold a lock in
// between, so each signature reserves its nonce instead: the next nonce of
// a key is the larger of the chain's pending nonce and the last
// reservation's next one, taken under the key's lock (the same lock
// gnosissafe.SendTransaction holds), so no two transactions anywhere get
// the same nonce. A reservation that is never sent stops counting after
// nonceReservationTTL: the key then continues from the chain's pending
// nonce rather than waiting forever behind the gap.

var nonceReservationTTL = 10 * time.Minute

// NonceReservation is the next nonce a platform key hands out while its
// last reservation is fresh.
type NonceReservation struct {
	Address    string `gorm:"primaryKey;size:42"`
	Next       uint64
	ReservedAt time.Time
}

// PendingNonceReader reads a key's pending nonce (an *ethclient.Client).
type PendingNonceReader interface {
	PendingNonceAt(ctx context.Context, account common.Address) (uint64, error)
}

func nonceKey(from common.Address) string { return strings.ToLower(from.Hex()) }

// NextNonceLocked reserves from's next nonce; the caller must hold the
// key's lock ("key:eoa:<address>"). release gives the nonce back when the
// transaction was not sent after all.
func NextNonceLocked(ctx context.Context, db *gorm.DB, reader PendingNonceReader, from common.Address) (n uint64, release func(), err error) {
	pending, err := reader.PendingNonceAt(ctx, from)
	if err != nil {
		return 0, nil, err
	}
	n = pending
	key := nonceKey(from)
	var r NonceReservation
	if db.First(&r, "address = ?", key).Error == nil && time.Since(r.ReservedAt) < nonceReservationTTL && r.Next > n {
		n = r.Next
	}
	if err := db.Save(&NonceReservation{Address: key, Next: n + 1, ReservedAt: time.Now()}).Error; err != nil {
		return 0, nil, err
	}
	release = func() {
		// only if nothing was reserved after it
		db.Model(&NonceReservation{}).Where("address = ? AND next = ?", key, n+1).Update("next", n)
	}
	return n, release, nil
}

// ReserveNonce reserves from's next nonce, taking the key's lock.
func ReserveNonce(ctx context.Context, db *gorm.DB, reader PendingNonceReader, from common.Address) (uint64, error) {
	var n uint64
	err := WithKeyLock(db, "key:eoa:"+nonceKey(from), 3*time.Minute, func() error {
		var e error
		n, _, e = NextNonceLocked(ctx, db, reader, from)
		return e
	})
	return n, err
}
