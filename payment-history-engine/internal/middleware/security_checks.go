package middleware

import (
	"encoding/base64"
	"errors"
	"log"
	"strings"

	tErrors "trovo-wallet-payment-history-engine/internal/errors"
	"trovo-wallet-payment-history-engine/internal/evmkeypair"
)

// SignString returns an EIP-191 personal_sign, base64 encoded signature of
// toSign - the Base equivalent of the original's Stellar ed25519 signature.
func SignString(toSign string, secretKey string) (string, error) {
	kp, keyPairError := evmkeypair.ParseFull(secretKey)
	if keyPairError != nil {
		return "", keyPairError
	}
	toSign = strings.TrimSpace(toSign)
	signature, err := kp.SignBase64([]byte(toSign))

	if err != nil {
		return "", err
	}

	return signature, nil
}

// SignHttp returns a signed base64 encoded string of fullPathWithQuery+keyParam. keyParam = publicKey+timestamp
func SignHttp(fullPathWithQuery string, keyParam string, secretKey string) (string, error) {
	keyParam = strings.TrimSpace(keyParam)
	fullPathWithQuery = strings.TrimSpace(fullPathWithQuery)
	signature, err := SignString(fullPathWithQuery+keyParam, secretKey)

	if err != nil {
		return "", err
	}

	return signature, nil
}

// SignBase64Txn signs a transaction digest with the secret key. On Base,
// there is no XDR envelope to parse and re-hash: base64Txn is now a
// base64-encoded transaction digest computed upstream, and this function's
// job is only to sign it. networkPassPhrase is vestigial (kept only for
// call-site compatibility).
func SignBase64Txn(secretKey string, base64Txn string, networkPassPhrase string) (string, error) {
	_ = networkPassPhrase

	kp, keyPairError := evmkeypair.ParseFull(secretKey)
	if keyPairError != nil {
		return "", keyPairError
	}

	digest, err := base64.StdEncoding.DecodeString(base64Txn)
	if err != nil {
		return "", errors.New("could not decode transaction digest")
	}

	signature, err := kp.SignBase64(digest)

	if err != nil {
		return "", err
	}

	return signature, nil

}

// VerifySignatureString verifies if the signatures match with the one to be generated from toSign. toSign = publicKey+timestamp
func VerifySignatureString(toSign string, base64Signature string, signerAddress string) error {
	kp, errParsingAddress := evmkeypair.ParseAddress(signerAddress)
	if errParsingAddress != nil {
		return &tErrors.ErrorInvalidAddress{}
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

// VerifyHttpSignature verifies the httpRequest signation retrieved from X-TW-SIGNATURE header. keyParam = signerAddress+timestamp
func VerifyHttpSignature(fullPathWithQuery string, keyParam string, base64Signature string, signerAddress string) error {
	keyParam = strings.TrimSpace(keyParam)
	fullPathWithQuery = strings.TrimSpace(fullPathWithQuery)
	signatureError := VerifySignatureString(fullPathWithQuery+keyParam, base64Signature, signerAddress)

	if signatureError != nil {
		log.Printf("invalid signature: %s\n", signatureError)
		return &tErrors.ErrorInvalidAuthenticationSignature{}
	}

	return nil

}
