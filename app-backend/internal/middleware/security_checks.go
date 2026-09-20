package middleware

import (
	"encoding/base64"
	"errors"
	"log"
	"strings"

	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/evmkeypair"
)

// SignString returns an EIP-191 personal_sign, base64 encoded signature of
// toSign - the Base equivalent of the original's Stellar ed25519 signature.
func SignString(toSign string, secretKey string) (string, error) {
	kp, keyPairError := evmkeypair.ParseFull(secretKey)
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

// SignHttp returns a signed base64 encoded string of fullPathWithQuery+keyParam. keyParam = signerAddress+timestamp
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

// SignBase64Txn signs a transaction digest with the secret key. On Base,
// there is no XDR envelope to parse and re-hash: the unsigned digest
// (base64Txn - kept named for call-site compatibility, but now a
// base64-encoded 32-byte Base tx/SafeTxHash digest, not an XDR blob) is
// computed once upstream when the transaction is built (internal/network),
// and this function's job is only to sign it. networkPassPhrase is
// vestigial (Stellar network-passphrase domain separation has no Base
// equivalent - chain ID is already baked into the digest upstream) and is
// kept only so existing call sites don't need to change their argument
// count; it is ignored here.
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

// SignSubwalletBase64Txn signs the same transaction digest with up to three
// signers (primary, sub-wallet, and an optional linked wallet) - the Base
// equivalent of co-signing one Stellar multisig-account transaction: each
// signer's personal_sign signature over the same digest is returned for
// the caller to relay together. See SignBase64Txn for why base64Txn is now
// a raw digest and networkPassPhrase is unused.
func SignSubwalletBase64Txn(primarySecretKey, subWalletSecretKey, linkedWalletSecret string, base64Txn string, networkPassPhrase string) (primarySignature, subWalletSignature, linkedWalletSignature string, err error) {
	_ = networkPassPhrase

	primaryKP, keyPairError := evmkeypair.ParseFull(primarySecretKey)
	if keyPairError != nil {
		return "", "", "", keyPairError
	}
	subWalletKP, keyPairError := evmkeypair.ParseFull(subWalletSecretKey)
	if keyPairError != nil {
		return "", "", "", keyPairError
	}

	digest, err := base64.StdEncoding.DecodeString(base64Txn)
	if err != nil {
		return "", "", "", errors.New("could not decode transaction digest")
	}

	primarySignature, err = primaryKP.SignBase64(digest)

	if err != nil {
		return "", "", "", err
	}
	subWalletSignature, err = subWalletKP.SignBase64(digest)

	if err != nil {
		return "", "", "", err
	}

	trimmedLinkedWalletSecret := strings.TrimPrefix(strings.TrimSpace(linkedWalletSecret), "0x")
	if len(trimmedLinkedWalletSecret) == 64 {
		linkedWalletKP, keyPairError := evmkeypair.ParseFull(linkedWalletSecret)
		if keyPairError != nil {
			return "", "", "", keyPairError
		}
		linkedWalletSignature, err = linkedWalletKP.SignBase64(digest)

		if err != nil {
			return "", "", "", err
		}
	}

	return primarySignature, subWalletSignature, linkedWalletSignature, nil

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
