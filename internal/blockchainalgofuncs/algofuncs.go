package blockchainalgofuncs

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

func GetRecoveryAccountAddress(username, publicKey string) (recoveryAddress string) {

	kp, err := RecoveryAccountKeypair(username, publicKey)
	if err != nil {
		return ""
	}
	return kp.Address()

}

func MarketMakingSignerKeypair(username, walletPublicKey string) (*keypair.Full, error) {
	kAccountSalt := "e45nDk4rk4LAhbX"
	kExtraAccountSalt := os.Getenv("MARKET_MAKING_SALT")
	if len(kExtraAccountSalt) == 0 {
		kExtraAccountSalt = "7gKsLLMJBNGfdvGGSHDUEUzma8fgg653$423"
	}
	mnemonic := os.Getenv("MNEMONIC_MARKET_MAKING")

	h := sha256.New()
	h.Write([]byte(kAccountSalt))
	h.Write([]byte(kExtraAccountSalt))
	h.Write([]byte(mnemonic))
	h.Write([]byte(username))
	h.Write([]byte(walletPublicKey))

	hashed := h.Sum(nil)

	var rawSeed [32]byte
	copy(rawSeed[:], hashed[0:32])

	return keypair.FromRawSeed([32]byte(rawSeed))

}

func GetMarketMakingSignerAddress(username, walletPublicKey string) (mmAddress string) {

	kp, err := MarketMakingSignerKeypair(username, walletPublicKey)
	if err != nil {
		return ""
	}
	return kp.Address()

}

func BulkPaymentSignerKeypair(username, walletPublicKey string) (*keypair.Full, error) {
	kAccountSalt := "e45nDk4rk4LAhbX"
	kExtraAccountSalt := os.Getenv("BULK_PAYMENT_SALT")
	if len(kExtraAccountSalt) == 0 {
		kExtraAccountSalt = "7gKsLLMJBNGfdvGGSHDUEUzma8fgg653$423"
	}
	mnemonic := os.Getenv("MNEMONIC_BULK_PAYMENT")

	h := sha256.New()
	h.Write([]byte(kAccountSalt))
	h.Write([]byte(kExtraAccountSalt))
	h.Write([]byte(mnemonic))
	h.Write([]byte(username))
	h.Write([]byte(walletPublicKey))

	hashed := h.Sum(nil)

	var rawSeed [32]byte
	copy(rawSeed[:], hashed[0:32])

	return keypair.FromRawSeed([32]byte(rawSeed))

}

func GetBulkPaymentSignerAddress(username, walletPublicKey string) (bulkPaymentAddress string) {

	kp, err := MarketMakingSignerKeypair(username, walletPublicKey)
	if err != nil {
		return ""
	}
	return kp.Address()

}

func EncodeSha256(str string) string {
	kAccountSalt := "e45nDk4rk4LAhbX"
	kExtraAccountSalt := os.Getenv("ENCODER_SALT")
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
