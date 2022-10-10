package blockchainalgofuncs

import (
	"crypto/rand"
	"io"
	"os"
)

// Encrypt encryps data with passphrase
func Encrypt(data []byte, passphrase string) (hash []byte, err error) {
	gcm, err := getGCM(passphrase)
	if err != nil {
		return
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

// Decrypt decryps data with passphrase
func Decrypt(data []byte, passphrase string) (hash []byte, err error) {
	gcm, err := getGCM(passphrase)
	if err != nil {
		return
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return
	}
	return plaintext, nil
}

// EncryptFile encrypts data with passphrase and saves as filename
func EncryptFile(filename string, data []byte, passphrase string) {
	f, _ := os.Create(filename)
	defer f.Close()
	hash, err := Encrypt(data, passphrase)
	if err == nil {
		f.Write(hash)
	}

}

// DecryptFile encrypts filename with passphrase
func DecryptFile(filename string, passphrase string) ([]byte, error) {
	data, _ := os.ReadFile(filename)
	return Decrypt(data, passphrase)
}
