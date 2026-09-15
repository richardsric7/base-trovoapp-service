// Package basetxn is a Base (EVM) transaction-building package that
// mirrors the shape of github.com/stellar/go/txnbuild - the same Asset/
// Operation/Transaction vocabulary this codebase's service layer already
// uses - so ported call sites need minimal changes (mostly the import and
// a few field/const renames), not a rewrite of their business logic.
//
// The one place this cannot be a faithful drop-in: a Stellar transaction
// atomically bundles multiple operations, each with its own source
// account/signer, into one signed envelope. A plain Base/EVM transaction
// is always exactly one call from exactly one signer. Rather than require
// deploying new multisig/multicall contracts (out of scope for an
// alteration pass), a basetxn.Transaction with N operations becomes N
// ordered, independently-signed-and-submitted native transactions -
// application-atomic (we stop at the first failure and report how far we
// got), not chain-atomic. See Transaction.Submit.
//
// A second, narrower divergence: Stellar's trustline system
// (ChangeTrust/SetTrustLineFlags) is an on-chain, per-account opt-in a
// plain ERC-20-shaped B20 token contract has no equivalent for. Ported
// code that builds these operations keeps working unchanged; this
// package resolves them against this application's own DB-backed
// wallet/asset authorization table (internal/components/assets) instead
// of an on-chain call - i.e. "authorized wallets can hold/send a B20
// asset" is enforced by this backend before it will relay a payment, the
// same way the original enforced it by checking a Stellar trustline's
// authorization flag before submitting.
package basetxn

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"trovo-wallet-api/internal/evmkeypair"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

// Asset is a Base asset: either the chain's native token, or a B20 token
// (an ERC-20-shaped fungible token contract). Mirrors Stellar's
// txnbuild.Asset interface.
type Asset interface {
	GetCode() string
	// GetIssuer returns the B20 token contract's address, or "" for the
	// native asset - the Base equivalent of a Stellar asset issuer.
	GetIssuer() string
	IsNative() bool
}

// NativeAsset is the chain's native token (e.g. ETH on Base) - the Base
// equivalent of Stellar's native XLM asset.
type NativeAsset struct{}

func (NativeAsset) GetCode() string   { return "" }
func (NativeAsset) GetIssuer() string { return "" }
func (NativeAsset) IsNative() bool    { return true }

// CreditAsset is a B20 token: Code is its symbol, Issuer is the B20 token
// contract's address. The Base equivalent of Stellar's txnbuild.CreditAsset
// (asset code + issuer account).
type CreditAsset struct {
	Code   string
	Issuer string
}

func (a CreditAsset) GetCode() string   { return a.Code }
func (a CreditAsset) GetIssuer() string { return a.Issuer }
func (a CreditAsset) IsNative() bool    { return false }

// TrustLineFlag mirrors Stellar's txnbuild.TrustLineFlag constants.
type TrustLineFlag uint32

const (
	TrustLineAuthorized                      TrustLineFlag = 1
	TrustLineAuthorizedToMaintainLiabilities TrustLineFlag = 2
	TrustLineClawbackEnabled                 TrustLineFlag = 4
)

// AccountFlag mirrors Stellar's txnbuild.AccountFlag constants (an
// issuer's own authorization policy for its asset).
type AccountFlag uint32

const (
	AuthRequired        AccountFlag = 1
	AuthRevocable       AccountFlag = 2
	AuthImmutable       AccountFlag = 4
	AuthClawbackEnabled AccountFlag = 8
)

// Operation is one step of a Transaction. SourceAccount, when set,
// overrides the Transaction's own SourceAccount as the signer/payer for
// that one step - the Base equivalent of Stellar's per-operation source
// account.
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

// ChangeTrust authorizes (Limit != "0") or removes authorization (Limit
// == "0") for SourceAccount to hold/send Line. Resolved against this
// app's DB authorization table (see package doc), not sent on-chain.
type ChangeTrust struct {
	Line          Asset
	Limit         string
	SourceAccount string
}

func (c ChangeTrust) opSourceAccount() string { return c.SourceAccount }

// SetTrustLineFlags authorizes/deauthorizes Trustor to hold/send Asset -
// the issuer-side counterpart to ChangeTrust. Resolved against this app's
// DB authorization table (see package doc), not sent on-chain.
type SetTrustLineFlags struct {
	Trustor       string
	Asset         Asset
	SetFlags      []TrustLineFlag
	ClearFlags    []TrustLineFlag
	SourceAccount string
}

func (s SetTrustLineFlags) opSourceAccount() string { return s.SourceAccount }


// PathPaymentStrictSend/PathPaymentStrictReceive mirror Stellar's
// same-named operations: swap SendAsset for DestAsset via the network's
// best available price path, either fixing the send amount (StrictSend,
// DestMin bounds the worst acceptable output) or the receive amount
// (StrictReceive, SendMax bounds the worst acceptable input). Stellar
// resolved the "path" against its native on-chain DEX order book; Base
// has no such native DEX (see internal/sharedconfig/
// order_book_summary.go's doc) - a real implementation needs a specific
// DEX router (e.g. Uniswap v3) wired into internal/network, tracked as a
// follow-up. The operation shape is kept so ported swap-building code
// changes minimally once that integration lands; until then, submitting
// either operation fails clearly rather than silently no-opping.
type PathPaymentStrictSend struct {
	SendAsset     Asset
	SendAmount    string
	Destination   string
	DestAsset     Asset
	DestMin       string
	Path          []Asset
	SourceAccount string
}

func (p PathPaymentStrictSend) opSourceAccount() string { return p.SourceAccount }

type PathPaymentStrictReceive struct {
	SendAsset     Asset
	SendMax       string
	Destination   string
	DestAsset     Asset
	DestAmount    string
	Path          []Asset
	SourceAccount string
}

func (p PathPaymentStrictReceive) opSourceAccount() string { return p.SourceAccount }

// ManageSellOffer placed/updated/canceled a standing sell order on
// Stellar's native on-chain DEX order book. Base has no native order
// book to hold one (see internal/sharedconfig/order_book_summary.go's
// doc) - market-making on Base needs a real design (an on-chain
// limit-order contract, or an off-chain maker service), tracked as a
// follow-up and resolved as a no-op here (see package doc's
// Operation-resolution note), same as PathPaymentStrictSend/Receive.
type ManageSellOffer struct {
	Selling       Asset
	Buying        Asset
	Amount        string
	Price         string
	OfferID       int64
	SourceAccount string
}

func (m ManageSellOffer) opSourceAccount() string { return m.SourceAccount }

// CreateAccount has no Base equivalent (any address is valid without an
// on-chain creation step) - kept, with Stellar's own field name (Amount),
// only so ported code compiles unchanged; resolved as a no-op (see
// package doc's Operation-resolution note).
type CreateAccount struct {
	Destination   string
	Amount        string
	SourceAccount string
}

func (c CreateAccount) opSourceAccount() string { return c.SourceAccount }

// Signer mirrors Stellar's txnbuild.Signer (an address + a weight in
// Stellar's weighted-multisig scheme).
type Signer struct {
	Address string
	Weight  uint32
}

// SetOptions.Signer added/removed a co-signer on a Stellar account -
// Stellar's native weighted multisig, which a plain Base EOA has no
// equivalent for. As the user/PLAN notes: doing this properly on Base
// needs every wallet to be a smart-contract account (e.g. a Safe), not a
// plain EOA - out of scope for this alteration pass. Until then, this is
// resolved as an app-layer signer registry (internal/network's
// AccountSigner - the same DB-backed-authorization pattern used for B20
// wallet authorization) rather than a real on-chain capability: it lets
// account-recovery/subwallet flows that check "is this address a signer
// for that account" keep working meaningfully, without granting any
// actual on-chain signing power.
type SetOptions struct {
	Signer     *Signer
	HomeDomain *string // vestigial (Stellar-only), kept for compatibility
	SetFlags   []AccountFlag
	ClearFlags []AccountFlag
	// LowThreshold/MediumThreshold/HighThreshold/MasterWeight configured
	// a Stellar account's weighted-multisig approval thresholds - the
	// shared-access/subwallet approval-count mechanism the original app
	// built on. Vestigial here for the same reason as Signer above: real
	// N-of-M approval on Base needs every participating wallet to be a
	// smart-contract account (e.g. a Safe), not a plain EOA. Kept, with
	// Stellar's own field names/*uint32 shape, purely for call-site
	// compatibility until that smart-account redesign happens.
	LowThreshold    *uint32
	MediumThreshold *uint32
	HighThreshold   *uint32
	MasterWeight    *uint32
	SourceAccount   string
}

func (s SetOptions) opSourceAccount() string { return s.SourceAccount }

// NewThreshold mirrors Stellar's txnbuild.NewThreshold(txnbuild.Threshold(n))
// - a *uint32 pointer helper for SetOptions' vestigial threshold fields.
func NewThreshold(n uint32) *uint32 { return &n }

type ManageData struct {
	Name          string
	Value         []byte
	SourceAccount string
}

func (m ManageData) opSourceAccount() string { return m.SourceAccount }

// TransactionParams mirrors txnbuild.TransactionParams.
type TransactionParams struct {
	SourceAccount string
	Operations    []Operation
	// BaseFee/Timebounds/Memo are accepted for source compatibility with
	// ported call sites but have no effect - Base's EIP-1559 fee market
	// and mempool expiry replace them.
	BaseFee              int64
	Memo                 string
	IncrementSequenceNum bool
}

// PendingPaymentTx is one still-to-be-signed native transaction produced
// from a Payment operation.
type PendingPaymentTx struct {
	Signer common.Address // the account whose key must sign this
	Op     Payment
	tx     *types.Transaction // built lazily at sign time by the caller-supplied Builder
}

// Transaction is a basetxn equivalent of a Stellar txnbuild.Transaction:
// an ordered list of operations built against SourceAccount, resolved and
// (for Payment ops) signed independently per §package-doc, then submitted
// in order.
type Transaction struct {
	SourceAccount string
	Operations    []Operation

	signed map[int][]byte // operation index -> signed raw tx bytes, once Sign has covered it
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

// Builder builds an unsigned, chain-ready native transaction for a
// Payment operation (nonce, gas, chain ID already resolved) - supplied by
// internal/network, which owns the RPC client basetxn itself does not
// depend on, to keep this package chain-client-free and unit-testable.
type Builder interface {
	BuildPaymentTx(ctx context.Context, from common.Address, op Payment) (*types.Transaction, error)
}

// defaultBuilder is wired up once at startup by internal/network (see
// SetDefaultBuilder) - keeping Sign's call-site shape
// (tx, err = tx.Sign(passphrase, keypairs...)) identical to Stellar's
// txnbuild.Transaction.Sign is what let the ~30 existing call sites
// across this codebase's service layer port with an import change only,
// rather than every one threading a context/builder through by hand.
var defaultBuilder Builder

// SetDefaultBuilder wires the Builder every Transaction.Sign call uses.
// Called once at startup (see internal/network).
func SetDefaultBuilder(b Builder) {
	defaultBuilder = b
}

// Sign builds (via the default Builder) and signs every Payment
// operation whose signer address matches one of the supplied keypairs,
// returning t itself for chaining - the Base equivalent of Stellar's
// tx.Sign(passphrase, keypairs...), which likewise only contributes the
// signatures it holds keys for and leaves the rest for a later Sign call
// (e.g. the client's own signature, attached by
// internal/network.SubmitXdrWithSignature). passphrase is accepted only
// for call-site compatibility - see the package doc.
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
		payment, ok := op.(Payment)
		if !ok {
			continue // non-Payment ops are resolved at Submit time, not signed
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
		unsignedTx, err := defaultBuilder.BuildPaymentTx(ctx, signerAddr, payment)
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

// AddSignedRawTx attaches an already-signed raw transaction (e.g. one
// produced client-side, arriving over the API as the "signature" the
// original Stellar flow expected) for operation index i - the Base
// equivalent of txnbuild.Transaction.AddSignatureBase64.
func (t *Transaction) AddSignedRawTx(i int, raw []byte) {
	t.signed[i] = raw
}

// FullySigned reports whether every Payment operation has a signed raw tx.
func (t *Transaction) FullySigned() bool {
	for i, op := range t.Operations {
		if _, ok := op.(Payment); !ok {
			continue
		}
		if _, signed := t.signed[i]; !signed {
			return false
		}
	}
	return true
}

// SignedRawTxs returns the signed raw transactions in operation order,
// for internal/network to broadcast.
func (t *Transaction) SignedRawTxs() [][]byte {
	out := make([][]byte, 0, len(t.signed))
	for i := range t.Operations {
		if raw, ok := t.signed[i]; ok {
			out = append(out, raw)
		}
	}
	return out
}

// Hash returns the keccak256 hash of operation i's signed raw tx (or, if
// unsigned, of its RLP-encodable request - callers needing a
// pre-signature digest should build via Builder and hash that directly).
// Kept named Hash, taking a vestigial passphrase argument, purely so
// ported call sites that called tx.Hash(networkPassPhrase) keep
// compiling; the passphrase has no Base meaning.
func (t *Transaction) Hash(_ string) ([32]byte, error) {
	raws := t.SignedRawTxs()
	if len(raws) == 0 {
		return [32]byte{}, errors.New("transaction has no signed operations to hash")
	}
	return crypto.Keccak256Hash(raws[len(raws)-1]), nil
}

// Base64 hex-encodes the signed raw transactions, one per line, joined -
// kept named Base64 (and returning a string) purely for source
// compatibility with ported call sites; it is not XDR and not literally
// base64 of a single envelope, since Base has no single-envelope
// equivalent (see package doc).
func (t *Transaction) Base64() (string, error) {
	raws := t.SignedRawTxs()
	parts := make([]string, len(raws))
	for i, raw := range raws {
		parts[i] = "0x" + common.Bytes2Hex(raw)
	}
	return strings.Join(parts, ","), nil
}
