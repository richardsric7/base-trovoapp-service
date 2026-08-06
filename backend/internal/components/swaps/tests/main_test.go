package swaps

import (
	"testing"

	"github.com/stellar/go/txnbuild"
)

func TestSignature(t *testing.T) {

	gTxn, err := txnbuild.TransactionFromXDR("AAAAAgAAAADirXYFhkNNhc2nJgA9VoQISOQ6e5D0aIzhozUI9xQVrwAAAGQAHJgYAAAAAQAAAAEAAAAAAAAAAAAAAABgHIWaAAAAAAAAAAEAAAAAAAAAAAAAAAAmsDl7eJ11cYWt32hyZHrg64sd4lJqL6lpTCjRDUGcsgAAAAAL/piHAAAAAAAAAAA=")

	if err != nil {
		t.Fatalf("error decoding xdr %v", err)
	}

	txn, ok := gTxn.Transaction()

	if !ok {
		t.Fatalf("error extracting transaction")

	}

	txn, err = txn.AddSignatureBase64("Bantu Testnet", "GDRK25QFQZBU3BONU4TAAPKWQQEERZB2POIPI2EM4GRTKCHXCQK27PY2", "4F5i38Sagc72DBYMqjZvrFNG3IrrASBYnt744xxyfC0kEKImuJqHMmxdeyQq1DtGOJGkFGw2GsXYsqAsJxsICQ==")

	if err != nil {
		t.Fatalf("Failed to verify signature [%v]", err)

	}
}
