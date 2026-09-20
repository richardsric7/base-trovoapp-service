package validators

import (
	"trovo-wallet-payment-history-engine/internal/errors"
	"trovo-wallet-payment-history-engine/internal/evmkeypair"
)

// ValidateAddressFormat validates a blockchain public key
func ValidateAddressFormat(publicKey string) error {
	_, err := evmkeypair.ParseAddress(publicKey)

	if err != nil {
		var x errors.ErrorInvalidAddress
		x.Address = publicKey
		return &x
	}
	return nil

}
