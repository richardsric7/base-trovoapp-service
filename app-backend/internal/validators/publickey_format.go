package validators

import (
	"trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/evmkeypair"
)

// ValidateAddressFormat validates a Base (EVM) wallet address. Kept its
// original Stellar-era name ("public key") since callers across the
// codebase still use that field name for what is now an address.
func ValidateAddressFormat(publicKey string) error {
	_, err := evmkeypair.ParseAddress(publicKey)

	if err != nil {
		var x errors.ErrorInvalidAddress
		x.Address = publicKey
		return &x
	}
	return nil

}
