package users

import (
	"crypto/sha256"
	"fmt"
	"os"

	"github.com/stellar/go/keypair"
)

func RecoveryAccountKeypair(username, publicKey string) (*keypair.Full, error) {
	kAccountSalt := "e45nDk4rk4LAhbX"
	kExtraAccountSalt := os.Getenv("ACCOUNT_RECOVERY_SALT")
	if len(kExtraAccountSalt) == 0 {
		kExtraAccountSalt = "7gKsg63jgGHfdtzma8)653$423"
	}
	mnemonic := os.Getenv("MNEMONIC_ACCOUNT_RECOVERY")

	h := sha256.New()
	h.Write([]byte(kAccountSalt))
	h.Write([]byte(kExtraAccountSalt))
	h.Write([]byte(mnemonic))
	h.Write([]byte(username))
	h.Write([]byte(publicKey))

	hashed := h.Sum(nil)

	var rawSeed [32]byte
	copy(rawSeed[:], hashed[0:32])

	return keypair.FromRawSeed([32]byte(rawSeed))

}

func EncodeSha256(str string) string {
	kAccountSalt := "e45nDk4rk4LAhbX"
	kExtraAccountSalt := os.Getenv("ACCOUNT_RECOVERY_SALT")
	if len(kExtraAccountSalt) == 0 {
		kExtraAccountSalt = "7gKsg63jgGHfdtzma8)653$423"
	}
	h := sha256.New()
	h.Write([]byte(kAccountSalt))
	h.Write([]byte(kExtraAccountSalt))
	h.Write([]byte(str))

	hashed := h.Sum(nil)

	return fmt.Sprintf("%x", hashed)

}

func GetRecoveryAccountAddress(username, publicKey string) (recoveryAddress string) {

	kp, err := RecoveryAccountKeypair(username, publicKey)
	if err != nil {
		return ""
	}
	return kp.Address()

}
