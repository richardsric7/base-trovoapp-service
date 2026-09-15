package blockchainalgofuncs

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
)

//createHash returns 32 byte hex hash of they key string
func createHash(key string) []byte {
	hash := sha256.Sum256([]byte(key))
	return hash[:]
}

//getGCM returns the Galois Counter Mode of the passphrase
func getGCM(passphrase string) (gcm cipher.AEAD, err error) {
	block, _ := aes.NewCipher(createHash(passphrase))
	gcm, err = cipher.NewGCM(block)
	return
}
