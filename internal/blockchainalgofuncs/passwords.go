package blockchainalgofuncs

import (
	"golang.org/x/crypto/bcrypt"
)

//HashPassword hashes passwords using bCrypt Algo
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

//CheckPasswordHash checks password hash to see if it matches
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
