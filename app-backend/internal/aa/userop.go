package aa

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"math/big"
	"sort"

	"trovo-wallet-api/internal/evmkeypair"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

// UserOperation is an EntryPoint v0.7 UserOperation in its unpacked (RPC)
// form. InitCode is factory ++ factoryData (empty for a deployed wallet);
// PaymasterAndData is paymaster ++ verification gas ++ postOp gas ++ data
// (empty when the wallet pays its own gas in ETH).
type UserOperation struct {
	Sender               common.Address
	Nonce                *big.Int
	InitCode             []byte
	CallData             []byte
	CallGasLimit         *big.Int
	VerificationGasLimit *big.Int
	PreVerificationGas   *big.Int
	MaxFeePerGas         *big.Int
	MaxPriorityFeePerGas *big.Int
	PaymasterAndData     []byte
	Signature            []byte
}

func zeroIfNil(v *big.Int) *big.Int {
	if v == nil {
		return new(big.Int)
	}
	return v
}

func pack128(hi, lo *big.Int) []byte {
	out := make([]byte, 32)
	zeroIfNil(hi).FillBytes(out[:16])
	zeroIfNil(lo).FillBytes(out[16:])
	return out
}

// SafeOp typehashes (Safe4337Module v0.3.0).
var (
	safeOpTypehash = crypto.Keccak256Hash([]byte("SafeOp(address safe,uint256 nonce,bytes initCode,bytes callData,uint128 verificationGasLimit,uint128 callGasLimit,uint256 preVerificationGas,uint128 maxPriorityFeePerGas,uint128 maxFeePerGas,bytes paymasterAndData,uint48 validAfter,uint48 validUntil,address entryPoint)"))
	domainTypehash = crypto.Keccak256Hash([]byte("EIP712Domain(uint256 chainId,address verifyingContract)"))
)

func word(v *big.Int) []byte { return common.LeftPadBytes(zeroIfNil(v).Bytes(), 32) }

func wordU64(v uint64) []byte { return word(new(big.Int).SetUint64(v)) }

// SafeOpHash is what the wallet's owners sign: the Safe4337Module's EIP-712
// hash of the operation (Safe4337Module.getOperationHash). validAfter and
// validUntil (unix seconds, 0 = unbounded) travel in the signature.
func (c Config) SafeOpHash(op UserOperation, validAfter, validUntil uint64) common.Hash {
	domain := crypto.Keccak256(domainTypehash.Bytes(), word(c.ChainID), common.LeftPadBytes(c.Safe4337Module.Bytes(), 32))
	structHash := crypto.Keccak256(
		safeOpTypehash.Bytes(),
		common.LeftPadBytes(op.Sender.Bytes(), 32),
		word(op.Nonce),
		crypto.Keccak256(op.InitCode),
		crypto.Keccak256(op.CallData),
		word(op.VerificationGasLimit),
		word(op.CallGasLimit),
		word(op.PreVerificationGas),
		word(op.MaxPriorityFeePerGas),
		word(op.MaxFeePerGas),
		crypto.Keccak256(op.PaymasterAndData),
		wordU64(validAfter),
		wordU64(validUntil),
		common.LeftPadBytes(c.EntryPoint.Bytes(), 32),
	)
	return crypto.Keccak256Hash([]byte{0x19, 0x01}, domain, structHash)
}

// UserOpHash is the EntryPoint's getUserOpHash (what bundlers and receipts
// identify an operation by).
func (c Config) UserOpHash(op UserOperation) common.Hash {
	inner := crypto.Keccak256(
		common.LeftPadBytes(op.Sender.Bytes(), 32),
		word(op.Nonce),
		crypto.Keccak256(op.InitCode),
		crypto.Keccak256(op.CallData),
		pack128(op.VerificationGasLimit, op.CallGasLimit),
		word(op.PreVerificationGas),
		pack128(op.MaxPriorityFeePerGas, op.MaxFeePerGas),
		crypto.Keccak256(op.PaymasterAndData),
	)
	return crypto.Keccak256Hash(inner, common.LeftPadBytes(c.EntryPoint.Bytes(), 32), word(c.ChainID))
}

func uint48Bytes(v uint64) []byte {
	var b [8]byte
	binary.BigEndian.PutUint64(b[:], v)
	return b[2:]
}

// OwnerSignature is one owner's personal_sign signature (65 bytes, v 27/28)
// over the SafeOp hash, as the apps produce it.
type OwnerSignature struct {
	Owner     common.Address
	Signature []byte
}

// VerifyOwnerSignature checks sig is owner's personal_sign of safeOpHash.
func VerifyOwnerSignature(owner common.Address, safeOpHash common.Hash, sig []byte) error {
	return evmkeypair.VerifyPersonal(owner, safeOpHash.Bytes(), sig)
}

// ErrNotEnoughSignatures is returned when fewer valid signatures than the
// threshold were supplied.
var ErrNotEnoughSignatures = errors.New("aa: not enough owner signatures")

// AssembleSignature builds the UserOperation signature the Safe4337Module
// expects: validAfter (6 bytes) ++ validUntil (6 bytes) ++ the owners'
// signatures sorted by owner address, each converted to Safe's eth_sign
// form (v + 4). Every signature is verified first; duplicates and
// non-owners are rejected, and exactly threshold signatures are used (Safe
// rejects trailing bytes).
func AssembleSignature(safeOpHash common.Hash, validAfter, validUntil uint64, owners []common.Address, threshold int, sigs []OwnerSignature) ([]byte, error) {
	isOwner := map[common.Address]bool{}
	for _, o := range owners {
		isOwner[o] = true
	}
	seen := map[common.Address]bool{}
	valid := make([]OwnerSignature, 0, len(sigs))
	for _, s := range sigs {
		if !isOwner[s.Owner] {
			return nil, fmt.Errorf("aa: %s is not an owner of this wallet", s.Owner.Hex())
		}
		if seen[s.Owner] {
			return nil, fmt.Errorf("aa: duplicate signature from %s", s.Owner.Hex())
		}
		if err := VerifyOwnerSignature(s.Owner, safeOpHash, s.Signature); err != nil {
			return nil, fmt.Errorf("aa: signature from %s does not sign this operation: %w", s.Owner.Hex(), err)
		}
		seen[s.Owner] = true
		valid = append(valid, s)
	}
	if len(valid) < threshold {
		return nil, fmt.Errorf("%w: have %d, need %d", ErrNotEnoughSignatures, len(valid), threshold)
	}
	sort.Slice(valid, func(i, j int) bool { return bytes.Compare(valid[i].Owner.Bytes(), valid[j].Owner.Bytes()) < 0 })
	valid = valid[:threshold]

	out := append(uint48Bytes(validAfter), uint48Bytes(validUntil)...)
	for _, s := range valid {
		sig := append([]byte{}, s.Signature...)
		if sig[64] < 27 {
			sig[64] += 27
		}
		sig[64] += 4 // eth_sign signature type
		out = append(out, sig...)
	}
	return out, nil
}

// DummySignature is a well-formed placeholder signature for gas estimation
// before the owners have signed (threshold signatures that recover to
// some address; validation fails but gas use is representative).
func DummySignature(threshold int) []byte {
	out := append(uint48Bytes(0), uint48Bytes(0)...)
	one := make([]byte, 65)
	for i := range 32 {
		one[i] = 0x01
		one[32+i] = 0x01
	}
	one[64] = 31
	for range threshold {
		out = append(out, one...)
	}
	return out
}

// RPC is the JSON form bundlers take (eth_sendUserOperation, v0.7).
type RPC struct {
	Sender                        common.Address  `json:"sender"`
	Nonce                         *hexutil.Big    `json:"nonce"`
	Factory                       *common.Address `json:"factory,omitempty"`
	FactoryData                   hexutil.Bytes   `json:"factoryData,omitempty"`
	CallData                      hexutil.Bytes   `json:"callData"`
	CallGasLimit                  *hexutil.Big    `json:"callGasLimit"`
	VerificationGasLimit          *hexutil.Big    `json:"verificationGasLimit"`
	PreVerificationGas            *hexutil.Big    `json:"preVerificationGas"`
	MaxFeePerGas                  *hexutil.Big    `json:"maxFeePerGas"`
	MaxPriorityFeePerGas          *hexutil.Big    `json:"maxPriorityFeePerGas"`
	Paymaster                     *common.Address `json:"paymaster,omitempty"`
	PaymasterVerificationGasLimit *hexutil.Big    `json:"paymasterVerificationGasLimit,omitempty"`
	PaymasterPostOpGasLimit       *hexutil.Big    `json:"paymasterPostOpGasLimit,omitempty"`
	PaymasterData                 hexutil.Bytes   `json:"paymasterData,omitempty"`
	Signature                     hexutil.Bytes   `json:"signature"`
}

func hb(v *big.Int) *hexutil.Big { return (*hexutil.Big)(zeroIfNil(v)) }

// ToRPC converts op to the bundler JSON form.
func (op UserOperation) ToRPC() RPC {
	r := RPC{
		Sender: op.Sender, Nonce: hb(op.Nonce), CallData: op.CallData,
		CallGasLimit: hb(op.CallGasLimit), VerificationGasLimit: hb(op.VerificationGasLimit),
		PreVerificationGas: hb(op.PreVerificationGas), MaxFeePerGas: hb(op.MaxFeePerGas),
		MaxPriorityFeePerGas: hb(op.MaxPriorityFeePerGas), Signature: op.Signature,
	}
	if len(op.InitCode) >= 20 {
		f := common.BytesToAddress(op.InitCode[:20])
		r.Factory, r.FactoryData = &f, op.InitCode[20:]
	}
	if len(op.PaymasterAndData) >= 52 {
		p := common.BytesToAddress(op.PaymasterAndData[:20])
		r.Paymaster = &p
		r.PaymasterVerificationGasLimit = hb(new(big.Int).SetBytes(op.PaymasterAndData[20:36]))
		r.PaymasterPostOpGasLimit = hb(new(big.Int).SetBytes(op.PaymasterAndData[36:52]))
		r.PaymasterData = op.PaymasterAndData[52:]
	}
	return r
}

// MaxCost is the most the operation can cost in wei (all gas limits at
// maxFeePerGas) - what a wallet paying in ETH must hold.
func (op UserOperation) MaxCost() *big.Int {
	gas := new(big.Int).Add(zeroIfNil(op.VerificationGasLimit), zeroIfNil(op.CallGasLimit))
	gas.Add(gas, zeroIfNil(op.PreVerificationGas))
	if len(op.PaymasterAndData) >= 52 {
		gas.Add(gas, new(big.Int).SetBytes(op.PaymasterAndData[20:36]))
		gas.Add(gas, new(big.Int).SetBytes(op.PaymasterAndData[36:52]))
	}
	return gas.Mul(gas, zeroIfNil(op.MaxFeePerGas))
}
