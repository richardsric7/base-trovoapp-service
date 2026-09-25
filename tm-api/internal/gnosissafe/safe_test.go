package gnosissafe

import (
	"crypto/ecdsa"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func TestPrevOwner(t *testing.T) {
	a := common.HexToAddress("0x1111111111111111111111111111111111111111")
	b := common.HexToAddress("0x2222222222222222222222222222222222222222")
	c := common.HexToAddress("0x3333333333333333333333333333333333333333")
	owners := []common.Address{a, b, c}

	if prev, err := PrevOwner(owners, a); err != nil || prev.Hex() != common.HexToAddress(SentinelOwner).Hex() {
		t.Fatalf("PrevOwner(a) = %v, %v; want sentinel", prev, err)
	}
	if prev, err := PrevOwner(owners, b); err != nil || prev != a {
		t.Fatalf("PrevOwner(b) = %v, %v; want a", prev, err)
	}
	if prev, err := PrevOwner(owners, c); err != nil || prev != b {
		t.Fatalf("PrevOwner(c) = %v, %v; want b", prev, err)
	}
	other := common.HexToAddress("0x4444444444444444444444444444444444444444")
	if _, err := PrevOwner(owners, other); err == nil {
		t.Fatal("PrevOwner(not-an-owner) should error")
	}
}

func genKey(t *testing.T) *ecdsa.PrivateKey {
	t.Helper()
	priv, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	return priv
}

func TestSignAndConcatSignatures(t *testing.T) {
	k1, k2, k3 := genKey(t), genKey(t), genKey(t)
	addr1 := crypto.PubkeyToAddress(k1.PublicKey)
	addr2 := crypto.PubkeyToAddress(k2.PublicKey)
	addr3 := crypto.PubkeyToAddress(k3.PublicKey)

	var hash [32]byte
	copy(hash[:], crypto.Keccak256([]byte("test safe tx")))

	sigs := map[common.Address][]byte{}
	for addr, k := range map[common.Address]*ecdsa.PrivateKey{addr1: k1, addr2: k2, addr3: k3} {
		sig, err := SignTransactionHash(k, hash)
		if err != nil {
			t.Fatalf("SignTransactionHash: %v", err)
		}
		if len(sig) != 65 {
			t.Fatalf("signature length = %d, want 65", len(sig))
		}
		if v := sig[64]; v != 27 && v != 28 {
			t.Fatalf("signature v = %d, want 27 or 28", v)
		}
		// Recover and confirm it matches the signing address - proves this
		// really is a raw-hash signature Safe's ecrecover(hash, v, r, s)
		// path (not EIP-191 personal_sign) will accept.
		normalized := make([]byte, 65)
		copy(normalized, sig)
		normalized[64] -= 27
		pubKey, err := crypto.SigToPub(hash[:], normalized)
		if err != nil {
			t.Fatalf("SigToPub: %v", err)
		}
		if recovered := crypto.PubkeyToAddress(*pubKey); recovered != addr {
			t.Fatalf("recovered signer = %v, want %v", recovered, addr)
		}
		sigs[addr] = sig
	}

	concatenated := ConcatSignatures(sigs)
	if len(concatenated) != 65*3 {
		t.Fatalf("concatenated length = %d, want %d", len(concatenated), 65*3)
	}

	// Must come back out sorted ascending by address, each still a valid
	// 65-byte sig for its slot.
	addrs := []common.Address{addr1, addr2, addr3}
	for i := 0; i < len(addrs); i++ {
		for j := i + 1; j < len(addrs); j++ {
			if addrs[i].Hex() > addrs[j].Hex() {
				addrs[i], addrs[j] = addrs[j], addrs[i]
			}
		}
	}
	for i, addr := range addrs {
		want := sigs[addr]
		got := concatenated[i*65 : (i+1)*65]
		for b := range want {
			if got[b] != want[b] {
				t.Fatalf("signature slot %d byte %d mismatch: got %x want %x", i, b, got[b], want[b])
			}
		}
	}
}

func TestSwapOwnerAndRemoveOwnerCalldataRoundTrip(t *testing.T) {
	prev := common.HexToAddress(SentinelOwner)
	old := common.HexToAddress("0x2222222222222222222222222222222222222222")
	fresh := common.HexToAddress("0x3333333333333333333333333333333333333333")

	swapData, err := SwapOwnerCalldata(prev, old, fresh)
	if err != nil {
		t.Fatalf("SwapOwnerCalldata: %v", err)
	}
	method, err := safeABI.MethodById(swapData[:4])
	if err != nil {
		t.Fatalf("MethodById(swap): %v", err)
	}
	if method.Name != "swapOwner" {
		t.Fatalf("decoded method = %s, want swapOwner", method.Name)
	}

	removeData, err := RemoveOwnerCalldata(prev, old, big.NewInt(1))
	if err != nil {
		t.Fatalf("RemoveOwnerCalldata: %v", err)
	}
	method, err = safeABI.MethodById(removeData[:4])
	if err != nil {
		t.Fatalf("MethodById(remove): %v", err)
	}
	if method.Name != "removeOwner" {
		t.Fatalf("decoded method = %s, want removeOwner", method.Name)
	}
}
