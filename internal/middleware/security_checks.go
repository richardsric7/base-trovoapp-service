package middleware

import (
	"encoding/base64"
	"errors"
	"log"
	"strings"

	tErrors "trovo-wallet-api/internal/errors"

	"github.com/stellar/go/keypair"
	"github.com/stellar/go/txnbuild"
)

//SignString returns a signed base64 encoded string of toSign
func SignString(toSign string, secretKey string) (string, error) {
	kp, keyPairError := keypair.ParseFull(secretKey)
	if keyPairError != nil {
		return "", keyPairError
	}
	toSign = strings.TrimSpace(toSign)
	// log.Printf("toSign:[%v]\n", toSign)
	signature, err := kp.SignBase64([]byte(toSign))
	// ps, _ := base64.StdEncoding.DecodeString(signature)
	// log.Printf("provided signature:[%v]\n", ps)

	if err != nil {
		return "", err
	}

	return signature, nil
}

//SignHttp returns a signed base64 encoded string of fullPathWithQuery+keyParam. keyParam = publicKey+timestamp
func SignHttp(fullPathWithQuery string, keyParam string, secretKey string) (string, error) {
	// log.Printf("path + string:[%v]\n", fullPathWithQuery+body)
	keyParam = strings.TrimSpace(keyParam)
	fullPathWithQuery = strings.TrimSpace(fullPathWithQuery)
	signature, err := SignString(fullPathWithQuery+keyParam, secretKey)

	if err != nil {
		return "", err
	}

	return signature, nil
}

//SignBase64Txn signs the transaction hash from base64Txn string using the secret key
func SignBase64Txn(secretKey string, base64Txn string, networkPassPhrase string) (string, error) {

	kp, keyPairError := keypair.ParseFull(secretKey)
	if keyPairError != nil {
		return "", keyPairError
	}

	tx, err := txnbuild.TransactionFromXDR(base64Txn)
	if err != nil {
		return "", err
	}

	txn, b := tx.Transaction()

	if !b {
		return "", errors.New("not a txn")
	}

	bytes, err := txn.Hash(networkPassPhrase)

	if err != nil {
		return "", errors.New("could not hash txn")
	}

	signature, err := kp.SignBase64(bytes[:])

	if err != nil {
		return "", err
	}

	return signature, nil

}

//SignBase64Txn signs the transaction hash from base64Txn string using the secret key
func SignSubwalletBase64Txn(primarySecretKey, subWalletSecretKey string, base64Txn string, networkPassPhrase string) (primarySignature, subWalletSignature string, err error) {

	primaryKP, keyPairError := keypair.ParseFull(primarySecretKey)
	if keyPairError != nil {
		return "", "", keyPairError
	}
	subWalletKP, keyPairError := keypair.ParseFull(subWalletSecretKey)
	if keyPairError != nil {
		return "", "", keyPairError
	}

	tx, err := txnbuild.TransactionFromXDR(base64Txn)
	if err != nil {
		return "", "", err
	}

	txn, b := tx.Transaction()

	if !b {
		return "", "", errors.New("not a txn")
	}

	bytes, err := txn.Hash(networkPassPhrase)

	if err != nil {
		return "", "", errors.New("could not hash txn")
	}

	primarySignature, err = primaryKP.SignBase64(bytes[:])

	if err != nil {
		return "", "", err
	}
	subWalletSignature, err = subWalletKP.SignBase64(bytes[:])

	if err != nil {
		return "", "", err
	}

	return primarySignature, subWalletSignature, nil

}

//VerifySignatureString verifies if the signatures match with the one to be generated from toSign. toSign = publicKey+timestamp
func VerifySignatureString(toSign string, base64Signature string, signerPublicKey string) error {
	kp, errParsingPublicKey := keypair.ParseAddress(signerPublicKey)
	if errParsingPublicKey != nil {
		return &tErrors.ErrorInvalidPublicKey{}
	}
	toSign = strings.TrimSpace(toSign)

	providedSignature, errDecoding := base64.StdEncoding.DecodeString(base64Signature)

	if errDecoding != nil {
		log.Printf("[VerifySignatureString] unable to decode base64 signature: [%s], err:[%v]\n", base64Signature, errDecoding)
		return &tErrors.ErrorInvalidAuthenticationSignature{}
	}

	signatureError := kp.Verify([]byte(toSign), providedSignature)

	if signatureError != nil {
		log.Printf("[VerifySignatureString]invalid signature: %s\n", signatureError)
		return &tErrors.ErrorInvalidAuthenticationSignature{}
	}

	return nil

}

//VerifyHttpSignature verifies the httpRequest signation retrieved from X-TW-SIGNATURE header. keyParam = signerPublicKey+timestamp
func VerifyHttpSignature(fullPathWithQuery string, keyParam string, base64Signature string, signerPublicKey string) error {
	keyParam = strings.TrimSpace(keyParam)
	fullPathWithQuery = strings.TrimSpace(fullPathWithQuery)
	signatureError := VerifySignatureString(fullPathWithQuery+keyParam, base64Signature, signerPublicKey)

	if signatureError != nil {
		log.Printf("invalid signature: %s\n", signatureError)
		return &tErrors.ErrorInvalidAuthenticationSignature{}
	}

	return nil

}
