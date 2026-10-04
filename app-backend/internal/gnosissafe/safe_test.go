package gnosissafe

import (
	"bytes"
	"math/big"
	"testing"

	"trovo-wallet-api/internal/evmkeypair"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

// The two typehashes are fixed constants in every Safe v1.3.0+ contract
// (SAFE_TX_TYPEHASH, DOMAIN_SEPARATOR_TYPEHASH); a typo in the type
// strings would silently produce signatures the Safe rejects.
func TestTypehashesMatchSafeContract(t *testing.T) {
	if got := hexutil.Encode(SafeTxTypehash); got != "0xbb8310d486368db6bd6f849402fdd73ad53d316b5a4b2644ad6efe0f941286d8" {
		t.Errorf("SAFE_TX_TYPEHASH = %s", got)
	}
	if got := hexutil.Encode(DomainSeparatorTypehash); got != "0x47e79534a245952e8b16893a336b85a3d9ea9fa8c573f3d803afb92a79469218" {
		t.Errorf("DOMAIN_SEPARATOR_TYPEHASH = %s", got)
	}
}

func TestDefaultMultiSendCallOnlyAddressIsWellFormed(t *testing.T) {
	if !common.IsHexAddress(DefaultMultiSendCallOnlyAddress) || len(DefaultMultiSendCallOnlyAddress) != 42 {
		t.Fatalf("DefaultMultiSendCallOnlyAddress %q is not a 20-byte hex address", DefaultMultiSendCallOnlyAddress)
	}
}

func TestSignSafeTxHashOrdersByOwnerAndRecovers(t *testing.T) {
	hash, err := ComputeSafeTxHash(big.NewInt(84532), common.HexToAddress("0x1111111111111111111111111111111111111111"),
		common.HexToAddress("0x2222222222222222222222222222222222222222"), big.NewInt(0), []byte{0xde, 0xad}, OperationCall, big.NewInt(7))
	if err != nil {
		t.Fatal(err)
	}
	signers := make([]*evmkeypair.Full, 0, 3)
	for i := 0; i < 3; i++ {
		kp, err := evmkeypair.Random()
		if err != nil {
			t.Fatal(err)
		}
		signers = append(signers, kp)
	}
	sigs, err := SignSafeTxHash(hash, signers)
	if err != nil {
		t.Fatal(err)
	}
	if len(sigs) != 65*3 {
		t.Fatalf("expected 3 concatenated 65-byte signatures, got %d bytes", len(sigs))
	}
	var prev common.Address
	for i := 0; i < 3; i++ {
		sig := append([]byte{}, sigs[i*65:(i+1)*65]...)
		if sig[64] != 27 && sig[64] != 28 {
			t.Fatalf("signature %d has v=%d, Safe expects 27/28", i, sig[64])
		}
		sig[64] -= 27
		pub, err := crypto.SigToPub(hash, sig)
		if err != nil {
			t.Fatal(err)
		}
		addr := crypto.PubkeyToAddress(*pub)
		if i > 0 && bytes.Compare(prev.Bytes(), addr.Bytes()) >= 0 {
			t.Fatalf("signatures not in ascending owner order: %s then %s", prev.Hex(), addr.Hex())
		}
		prev = addr
	}
}

func TestPackMultiSendTxLayout(t *testing.T) {
	to := common.HexToAddress("0x3333333333333333333333333333333333333333")
	data := []byte{1, 2, 3}
	packed := PackMultiSendTx(OperationCall, to, big.NewInt(5), data)
	if len(packed) != 1+20+32+32+3 {
		t.Fatalf("unexpected packed length %d", len(packed))
	}
	if packed[0] != OperationCall || !bytes.Equal(packed[1:21], to.Bytes()) {
		t.Fatal("operation/to not packed at the expected offsets")
	}
	if new(big.Int).SetBytes(packed[21:53]).Int64() != 5 || new(big.Int).SetBytes(packed[53:85]).Int64() != 3 {
		t.Fatal("value/data length not packed as 32-byte big-endian words")
	}
	if !bytes.Equal(packed[85:], data) {
		t.Fatal("data not appended")
	}
}

// create2SafeAddress must follow SafeProxyFactory's own derivation, so a
// retried deployment finds the Safe a previous attempt created.
func TestCreate2SafeAddressDerivation(t *testing.T) {
	factory := common.HexToAddress("0x4e1DCf7AD4e460CfD30791CCC4F9c8a4f820ec67")
	singleton := common.HexToAddress("0x29fcB43b46531BcA003ddC8FCB67FFE91900C762")
	creationCode := []byte{0x60, 0x80, 0x60, 0x40}
	initializer := []byte("setup-calldata")
	salt := big.NewInt(42)

	got := create2SafeAddress(factory, singleton, creationCode, initializer, salt)

	wantSalt := crypto.Keccak256(crypto.Keccak256(initializer), common.LeftPadBytes(salt.Bytes(), 32))
	initCode := append(append([]byte{}, creationCode...), common.LeftPadBytes(singleton.Bytes(), 32)...)
	var s [32]byte
	copy(s[:], wantSalt)
	want := crypto.CreateAddress2(factory, s, crypto.Keccak256(initCode))
	if got != want {
		t.Fatalf("got %s want %s", got.Hex(), want.Hex())
	}
	if other := create2SafeAddress(factory, singleton, creationCode, initializer, big.NewInt(43)); other == got {
		t.Fatal("different salt nonces must give different Safe addresses")
	}
}
