// Package quote prices and signs paymaster quotes for TrovoTokenPaymaster.
package quote

import (
	"encoding/binary"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// UserOp is the part of an ERC-4337 v0.7 UserOperation a quote covers
// (everything but the paymaster data and the signature).
type UserOp struct {
	Sender               common.Address
	Nonce                *big.Int
	InitCode             []byte // factory ++ factoryData, empty for a deployed wallet
	CallData             []byte
	VerificationGasLimit *big.Int
	CallGasLimit         *big.Int
	PreVerificationGas   *big.Int
	MaxFeePerGas         *big.Int
	MaxPriorityFeePerGas *big.Int

	PaymasterVerificationGasLimit *big.Int
	PaymasterPostOpGasLimit       *big.Int
}

var max128 = new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 128), big.NewInt(1))

// Validate checks the numeric fields are present and fit their packed slots.
func (op UserOp) Validate() error {
	for name, v := range map[string]*big.Int{
		"nonce": op.Nonce, "verificationGasLimit": op.VerificationGasLimit, "callGasLimit": op.CallGasLimit,
		"preVerificationGas": op.PreVerificationGas, "maxFeePerGas": op.MaxFeePerGas, "maxPriorityFeePerGas": op.MaxPriorityFeePerGas,
		"paymasterVerificationGasLimit": op.PaymasterVerificationGasLimit, "paymasterPostOpGasLimit": op.PaymasterPostOpGasLimit,
	} {
		if v == nil || v.Sign() < 0 {
			return fmt.Errorf("%s is required and must not be negative", name)
		}
		if name != "nonce" && name != "preVerificationGas" && v.Cmp(max128) > 0 {
			return fmt.Errorf("%s does not fit in 128 bits", name)
		}
	}
	if op.Sender == (common.Address{}) {
		return fmt.Errorf("sender is required")
	}
	if op.MaxPriorityFeePerGas.Cmp(op.MaxFeePerGas) > 0 {
		return fmt.Errorf("maxPriorityFeePerGas exceeds maxFeePerGas")
	}
	if len(op.InitCode) != 0 && len(op.InitCode) < 20 {
		return fmt.Errorf("initCode must start with the factory address")
	}
	return nil
}

func pack128(hi, lo *big.Int) [32]byte {
	var out [32]byte
	hi.FillBytes(out[:16])
	lo.FillBytes(out[16:])
	return out
}

func word(v *big.Int) []byte {
	return common.LeftPadBytes(v.Bytes(), 32)
}

// MaxCost is what the EntryPoint may charge the paymaster for op (wei):
// every gas limit times maxFeePerGas.
func (op UserOp) MaxCost() *big.Int {
	gas := new(big.Int).Add(op.VerificationGasLimit, op.CallGasLimit)
	gas.Add(gas, op.PaymasterVerificationGasLimit)
	gas.Add(gas, op.PaymasterPostOpGasLimit)
	gas.Add(gas, op.PreVerificationGas)
	return gas.Mul(gas, op.MaxFeePerGas)
}

// Hash is TrovoTokenPaymaster.getHash: what a quote signer signs (after
// the EIP-191 "\x19Ethereum Signed Message:\n32" prefix).
func Hash(op UserOp, chainID *big.Int, paymaster, token common.Address, validUntil, validAfter uint64, rate *big.Int) common.Hash {
	agl := pack128(op.VerificationGasLimit, op.CallGasLimit)
	pmg := pack128(op.PaymasterVerificationGasLimit, op.PaymasterPostOpGasLimit)
	fees := pack128(op.MaxPriorityFeePerGas, op.MaxFeePerGas)
	opHash := crypto.Keccak256(
		common.LeftPadBytes(op.Sender.Bytes(), 32),
		word(op.Nonce),
		crypto.Keccak256(op.InitCode),
		crypto.Keccak256(op.CallData),
		agl[:],
		pmg[:],
		word(op.PreVerificationGas),
		fees[:],
	)
	return common.BytesToHash(crypto.Keccak256(
		opHash,
		word(chainID),
		common.LeftPadBytes(paymaster.Bytes(), 32),
		common.LeftPadBytes(token.Bytes(), 32),
		word(new(big.Int).SetUint64(validUntil)),
		word(new(big.Int).SetUint64(validAfter)),
		word(rate),
	))
}

func uint48Bytes(v uint64) []byte {
	var b [8]byte
	binary.BigEndian.PutUint64(b[:], v)
	return b[2:]
}

// PaymasterData is the paymaster-specific part of paymasterAndData (the v0.7
// RPC "paymasterData" field): token | validUntil | validAfter | rate | signature.
func PaymasterData(token common.Address, validUntil, validAfter uint64, rate *big.Int, signature []byte) []byte {
	out := make([]byte, 0, 20+6+6+32+len(signature))
	out = append(out, token.Bytes()...)
	out = append(out, uint48Bytes(validUntil)...)
	out = append(out, uint48Bytes(validAfter)...)
	out = append(out, word(rate)...)
	return append(out, signature...)
}

// PaymasterAndData is the packed v0.7 field: paymaster | verification gas |
// postOp gas | PaymasterData.
func PaymasterAndData(paymaster common.Address, verificationGas, postOpGas *big.Int, data []byte) []byte {
	gl := pack128(verificationGas, postOpGas)
	out := append([]byte{}, paymaster.Bytes()...)
	out = append(out, gl[:]...)
	return append(out, data...)
}
