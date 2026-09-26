package users

import (
	tErrors "trovo-wallet-api/internal/errors"
)

// EnsureNotSuspended is the single check every transaction-creating flow
// (payments, swaps, P2P offers/orders, shared-wallet access grants,
// tokenized-asset subscriptions/minting, service links, ...) should call
// against every user whose wallet the flow is about to move funds through
// or grant access on. An admin-suspended account (Suspended == 1, set via
// tm-api) must not be usable on either side of a transaction.
func (u User) EnsureNotSuspended() error {
	if u.Suspended == 1 {
		return &tErrors.ErrorUsernameIsSuspended{Username: u.Username}
	}
	return nil
}
