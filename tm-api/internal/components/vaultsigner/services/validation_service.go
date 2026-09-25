package services

import (
	"strings"

	tErrors "admin-panel-dashboard/internal/errors"
	"admin-panel-dashboard/internal/evmkeypair"
)

// ValidateSignerValue runs the mandatory pipeline every submitted signer
// value goes through, for both signer slots and personal secrets (Section
// 5c): trim, then structural format validation. The pipeline's third step
// used to be a Horizon "is this account activated on-chain" check -
// vestigial on Base, where any well-formed address is valid without a
// separate on-chain account-creation step (same reasoning as app-backend's
// network.BlockchainAccountProperties), so it's dropped rather than ported.
func ValidateSignerValue(raw string) (address string, err error) {
	trimmed := strings.TrimSpace(raw)

	kp, parseErr := evmkeypair.ParseFull(trimmed)
	if parseErr != nil {
		return "", &tErrors.ErrorInvalidAddress{Address: trimmed}
	}
	return kp.Address(), nil
}

// IsValidKeypair is the lighter, structural-only check reused in a few
// places this feature needs it: Section 5d's "is the old value still a
// valid keypair" gate, and Section 5e's CountValidKeypairs trigger for
// assignment deletion.
func IsValidKeypair(raw string) bool {
	_, err := evmkeypair.ParseFull(strings.TrimSpace(raw))
	return err == nil
}
