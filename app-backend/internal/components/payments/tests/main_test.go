package payments

import (
	"encoding/base64"
	"testing"

	"trovo-wallet-api/internal/evmkeypair"
)

// TestSignature exercises the EIP-191 personal_sign round trip
// (evmkeypair.Full.SignBase64 / Verify) that replaced Stellar's XDR
// transaction signature verification (txnbuild.TransactionFromXDR +
// AddSignatureBase64) this test previously covered.
func TestSignature(t *testing.T) {
	kp, err := evmkeypair.Random()
	if err != nil {
		t.Fatalf("error generating keypair %v", err)
	}

	message := []byte("Bantu Testnet")

	sigBase64, err := kp.SignBase64(message)
	if err != nil {
		t.Fatalf("error signing message %v", err)
	}

	addr, err := evmkeypair.ParseAddress(kp.Address())
	if err != nil {
		t.Fatalf("error parsing address %v", err)
	}

	sig, err := base64.StdEncoding.DecodeString(sigBase64)
	if err != nil {
		t.Fatalf("error decoding signature %v", err)
	}

	if err := addr.Verify(message, sig); err != nil {
		t.Fatalf("Failed to verify signature [%v]", err)
	}
}
