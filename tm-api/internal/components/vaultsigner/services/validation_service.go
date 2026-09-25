package services

import (
	"strings"

	tErrors "admin-panel-dashboard/internal/errors"
	"admin-panel-dashboard/internal/network"

	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
)

// ValidateSignerValue runs the mandatory pipeline every submitted signer
// value goes through, for both signer slots and personal secrets (Section
// 5c): trim, structural format validation, then ledger activation. Reuses
// this repo's existing typed errors (internal/errors) rather than inventing
// a new error envelope for this feature.
func ValidateSignerValue(raw string) (publicKey string, err error) {
	trimmed := strings.TrimSpace(raw)

	kp, parseErr := keypair.ParseFull(trimmed)
	if parseErr != nil {
		return "", &tErrors.ErrorInvalidPublicKey{PublicKey: trimmed}
	}
	publicKey = kp.Address()

	client := network.GetBlockchainClient()
	if _, horizonErr := client.AccountDetail(horizonclient.AccountRequest{AccountID: publicKey}); horizonErr != nil {
		// Fail closed on any Horizon error, not just a 404 "not found" — an
		// unreachable Horizon must never be treated as "account is fine."
		return "", &tErrors.ErrorBlockchainAccountNotActivated{PublicKey: publicKey}
	}
	return publicKey, nil
}

// IsValidKeypair is the lighter, structural-only check reused in a few
// places this feature needs it without a Horizon round trip: Section 5d's
// "is the old value still a valid keypair" gate, and Section 5e's
// CountValidKeypairs trigger for assignment deletion.
func IsValidKeypair(raw string) bool {
	_, err := keypair.ParseFull(strings.TrimSpace(raw))
	return err == nil
}
