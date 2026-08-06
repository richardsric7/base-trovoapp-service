package validators

import (
	"trovo-wallet-api/internal/errors"

	"github.com/stellar/go/keypair"
)

//ValidatePublicKeyFormat validates a blockchain public key
func ValidatePublicKeyFormat(publicKey string) error {
	_, err := keypair.ParseAddress(publicKey)

	if err != nil {
		var x errors.ErrorInvalidPublicKey
		x.PublicKey = publicKey
		return &x
	}
	return nil

}
