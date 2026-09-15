// Package basetxn is a Base (EVM) transaction-building compatibility layer
// that mirrors the subset of github.com/stellar/go/txnbuild this module
// used (NativeAsset, CreditAsset, Payment, TransactionParams, NewTransaction,
// Transaction.Sign/Base64) so the market-making fee-collection transaction
// in main.go's ProcessTrade ports with minimal call-site changes. Unlike
// Stellar's single signed XDR envelope, a Base "transaction" here is really
// a set of independently signed native transactions (one per Payment
// operation, each with its own account/nonce) - see Transaction.Base64.
package basetxn

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"trovo-wallet-payment-history-engine/internal/evmkeypair"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Asset mirrors txnbuild.Asset - either NativeAsset (Base's own gas/native
// token) or CreditAsset (a B20/ERC20-shaped token; Issuer is the token
// contract address).
type Asset interface {
	IsNative() bool
	GetCode() string
	GetIssuer() string
}

type NativeAsset struct{}

func (NativeAsset) IsNative() bool     { return true }
func (NativeAsset) GetCode() string    { return "" }
func (NativeAsset) GetIssuer() string  { return "" }

// CreditAsset is a B20 (ERC-20-shaped) token - the Base equivalent of a
// Stellar issued asset. Issuer is the token contract address.
type CreditAsset struct {
	Code   string
	Issuer string
}

func (c CreditAsset) IsNative() bool    { return false }
func (c CreditAsset) GetCode() string   { return c.Code }
func (c CreditAsset) GetIssuer() string { return c.Issuer }

// Operation mirrors txnbuild.Operation.
type Operation interface {
	opSourceAccount() string
}

// Payment sends Amount (a decimal string, in the asset's own units) of
// Asset to Destination.
type Payment struct {
	Destination   string
	Amount        string
	Asset         Asset
	SourceAccount string
}

func (p Payment) opSourceAccount() string { return p.SourceAccount }

// TransactionParams mirrors txnbuild.TransactionParams.
type TransactionParams struct {
	SourceAccount string
	Operations    []Operation
	// BaseFee/Memo/IncrementSequenceNum are accepted for call-site
	// compatibility but have no effect - Base's EIP-1559 fee market and
	// mempool nonce replace them.
	BaseFee              int64
	Memo                 string
	IncrementSequenceNum bool
}

// Transaction is a basetxn equivalent of txnbuild.Transaction: an ordered
// list of Payment operations, each signed and submitted independently.
type Transaction struct {
	SourceAccount string
	Operations    []Operation

	signed map[int][]byte
}

// NewTransaction mirrors txnbuild.NewTransaction(params).
func NewTransaction(params TransactionParams) (*Transaction, error) {
	if len(params.Operations) == 0 {
		return nil, errors.New("transaction has no operations")
	}
	return &Transaction{
		SourceAccount: params.SourceAccount,
		Operations:    params.Operations,
		signed:        map[int][]byte{},
	}, nil
}

// Builder builds an unsigned, chain-ready native transaction for a Payment
// operation (nonce, gas, chain ID already resolved) - supplied by
// internal/network, which owns the RPC client this package does not depend
// on directly.
type Builder interface {
	BuildPaymentTx(ctx context.Context, from common.Address, op Payment) (*types.Transaction, error)
}

var defaultBuilder Builder

// SetDefaultBuilder wires the Builder every Transaction.Sign call uses.
// Called once at startup (see internal/network).
func SetDefaultBuilder(b Builder) {
	defaultBuilder = b
}

// Sign builds (via the default Builder) and signs every Payment operation
// whose signer address matches one of the supplied keypairs - the Base
// equivalent of Stellar's tx.Sign(passphrase, keypairs...). passphrase is
// accepted only for call-site compatibility.
func (t *Transaction) Sign(passphrase string, keypairs ...*evmkeypair.Full) (*Transaction, error) {
	if defaultBuilder == nil {
		return t, errors.New("basetxn: no Builder configured - call SetDefaultBuilder at startup")
	}
	ctx := context.Background()
	byAddr := map[common.Address]*evmkeypair.Full{}
	for _, kp := range keypairs {
		byAddr[common.HexToAddress(kp.Address())] = kp
	}

	for i, op := range t.Operations {
		payment, ok := op.(*Payment)
		if !ok {
			continue
		}
		if _, already := t.signed[i]; already {
			continue
		}
		source := payment.SourceAccount
		if source == "" {
			source = t.SourceAccount
		}
		signerAddr := common.HexToAddress(source)
		kp, have := byAddr[signerAddr]
		if !have {
			continue
		}
		unsignedTx, err := defaultBuilder.BuildPaymentTx(ctx, signerAddr, *payment)
		if err != nil {
			return t, fmt.Errorf("building payment tx for operation %d: %w", i, err)
		}
		signer := types.LatestSignerForChainID(unsignedTx.ChainId())
		signedTx, err := types.SignTx(unsignedTx, signer, kp.PrivateKey())
		if err != nil {
			return t, fmt.Errorf("signing operation %d: %w", i, err)
		}
		raw, err := signedTx.MarshalBinary()
		if err != nil {
			return t, err
		}
		t.signed[i] = raw
	}
	return t, nil
}

// SignedRawTxs returns the signed raw transactions in operation order, for
// internal/network to broadcast.
func (t *Transaction) SignedRawTxs() [][]byte {
	out := make([][]byte, 0, len(t.signed))
	for i := range t.Operations {
		if raw, ok := t.signed[i]; ok {
			out = append(out, raw)
		}
	}
	return out
}

// Base64 hex-encodes the signed raw transactions, one per line, joined -
// kept named Base64 for source compatibility; it is not literal base64 of a
// single XDR envelope, since Base has no single-envelope equivalent.
func (t *Transaction) Base64() (string, error) {
	raws := t.SignedRawTxs()
	parts := make([]string, len(raws))
	for i, raw := range raws {
		parts[i] = "0x" + common.Bytes2Hex(raw)
	}
	return strings.Join(parts, ","), nil
}
