// Package evmkeypair is a Base (EVM/secp256k1) keypair implementation that
// mirrors the subset of github.com/stellar/go/keypair's API this codebase
// used (Full, FromRawSeed, ParseFull, MustParseFull, ParseAddress, Address,
// Seed, Verify, SignBase64) so callers ported from the Stellar keypair need
// minimal changes - only the import and the underlying key material differ.
package evmkeypair

import (
	"crypto/ecdsa"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"strings"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// Full holds a secp256k1 private key - the Base equivalent of Stellar's
// ed25519 keypair.Full.
type Full struct {
	priv *ecdsa.PrivateKey
	addr common.Address
}

// AddressOnly holds just a Base address, with no private key - the Base
// equivalent of Stellar's keypair.FromAddress, returned by ParseAddress.
type AddressOnly struct {
	addr common.Address
}

// FromRawSeed derives a keypair from a 32-byte seed, the Base equivalent of
// Stellar's keypair.FromRawSeed. Used across this codebase's deterministic
// server-side signer derivation (recovery accounts, market-making signers,
// bulk-payment signers, temp accounts) - unchanged call sites, new curve.
func FromRawSeed(seed [32]byte) (*Full, error) {
	priv, err := crypto.ToECDSA(seed[:])
	if err != nil {
		return nil, err
	}
	return &Full{priv: priv, addr: crypto.PubkeyToAddress(priv.PublicKey)}, nil
}

// ParseFull parses a hex-encoded secp256k1 private key (0x-prefixed or not)
// - the Base equivalent of Stellar's base32 "S..." secret seed that
// keypair.ParseFull/MustParseFull accepted.
func ParseFull(secretHex string) (*Full, error) {
	secretHex = strings.TrimPrefix(strings.TrimSpace(secretHex), "0x")
	priv, err := crypto.HexToECDSA(secretHex)
	if err != nil {
		return nil, err
	}
	return &Full{priv: priv, addr: crypto.PubkeyToAddress(priv.PublicKey)}, nil
}

// Random generates a fresh keypair from CSPRNG entropy - the Base
// equivalent of Stellar's keypair.Random(), used to mint new issuer/
// distributor-style accounts.
func Random() (*Full, error) {
	priv, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}
	return &Full{priv: priv, addr: crypto.PubkeyToAddress(priv.PublicKey)}, nil
}

// MustParseFull is ParseFull, panicking on error - matches
// keypair.MustParseFull's contract (used for boot-time env/config secrets).
func MustParseFull(secretHex string) *Full {
	kp, err := ParseFull(secretHex)
	if err != nil {
		panic(err)
	}
	return kp
}

// ParseAddress validates a Base address string and returns an address-only
// KP - the Base equivalent of Stellar's keypair.ParseAddress.
func ParseAddress(address string) (*AddressOnly, error) {
	address = strings.TrimSpace(address)
	if !common.IsHexAddress(address) {
		return nil, errors.New("invalid address")
	}
	return &AddressOnly{addr: common.HexToAddress(address)}, nil
}

// Address returns the 0x-prefixed, EIP-55 checksummed address.
func (f *Full) Address() string { return f.addr.Hex() }

// Address returns the 0x-prefixed, EIP-55 checksummed address.
func (a *AddressOnly) Address() string { return a.addr.Hex() }

// Seed returns the hex-encoded private key - the Base equivalent of
// Stellar's base32 "S..." secret seed returned by keypair.Full.Seed().
func (f *Full) Seed() string {
	return hex.EncodeToString(crypto.FromECDSA(f.priv))
}

// PrivateKey exposes the raw ecdsa key for callers that need it directly
// (e.g. raw transaction signing) rather than through the compatibility
// methods below.
func (f *Full) PrivateKey() *ecdsa.PrivateKey { return f.priv }

// SignBase64 signs message with EIP-191 personal_sign (the Base equivalent
// of Stellar's ed25519 detached signature) and returns the base64-encoded
// 65-byte (r,s,v) signature - matching keypair.Full.SignBase64's shape so
// existing call sites don't need to change beyond the import.
func (f *Full) SignBase64(message []byte) (string, error) {
	sig, err := SignPersonal(f.priv, message)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(sig), nil
}

// Verify checks a base64-encoded EIP-191 personal_sign signature against
// this address - the Base equivalent of Stellar's ed25519 Verify.
func (a *AddressOnly) Verify(message []byte, sig []byte) error {
	return VerifyPersonal(a.addr, message, sig)
}

// Verify checks a base64-encoded EIP-191 personal_sign signature against
// this keypair's own address.
func (f *Full) Verify(message []byte, sig []byte) error {
	return VerifyPersonal(f.addr, message, sig)
}

// PersonalSignHash returns the EIP-191 digest ("\x19Ethereum Signed
// Message:\n" + len + message) that SignPersonal/VerifyPersonal sign and
// check - Base's equivalent of a Stellar transaction's network-passphrase
// hash, for offline/off-chain message authentication (HTTP request
// signing, login challenges) rather than on-chain transactions.
func PersonalSignHash(message []byte) common.Hash {
	return common.BytesToHash(accounts.TextHash(message))
}

// SignPersonal produces a 65-byte (r,s,v) recoverable EIP-191 signature,
// with v normalized to {27,28} as personal_sign/ecrecover tooling expects.
func SignPersonal(priv *ecdsa.PrivateKey, message []byte) ([]byte, error) {
	hash := PersonalSignHash(message)
	sig, err := crypto.Sign(hash.Bytes(), priv)
	if err != nil {
		return nil, err
	}
	sig[64] += 27
	return sig, nil
}

// VerifyPersonal recovers the signer address from a 65-byte (r,s,v)
// EIP-191 signature and checks it matches want.
func VerifyPersonal(want common.Address, message []byte, sig []byte) error {
	if len(sig) != 65 {
		return errors.New("invalid signature length")
	}
	hash := PersonalSignHash(message)

	normalized := make([]byte, 65)
	copy(normalized, sig)
	if normalized[64] >= 27 {
		normalized[64] -= 27
	}

	pubKey, err := crypto.SigToPub(hash.Bytes(), normalized)
	if err != nil {
		return err
	}
	recovered := crypto.PubkeyToAddress(*pubKey)
	if recovered != want {
		return errors.New("signature does not match address")
	}
	return nil
}
