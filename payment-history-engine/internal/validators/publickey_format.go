package validators

import (
	"trovo-wallet-payment-history-engine/internal/errors"
	"trovo-wallet-payment-history-engine/internal/evmkeypair"
)

//ValidatePublicKeyFormat validates a blockchain public key
func ValidatePublicKeyFormat(publicKey string) error {
	_, err := evmkeypair.ParseAddress(publicKey)

	if err != nil {
		var x errors.ErrorInvalidPublicKey
		x.PublicKey = publicKey
		return &x
	}
	return nil

}
