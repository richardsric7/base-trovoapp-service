package aa

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
)

// ChainReader is what building an operation reads from the chain
// (*ethclient.Client satisfies it).
type ChainReader interface {
	Caller
	SuggestGasTipCap(ctx context.Context) (*big.Int, error)
	HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error)
}

// Submitter sends a signed operation (the bundler).
type Submitter interface {
	EstimateGas(ctx context.Context, op UserOperation, entryPoint common.Address) (GasEstimate, error)
	Send(ctx context.Context, op UserOperation, entryPoint common.Address) (common.Hash, error)
}

// Quoter prices an operation's gas in a stablecoin (the quote service).
type Quoter interface {
	Quote(ctx context.Context, op UserOperation, token common.Address, validity time.Duration, pmVerificationGas, pmPostOpGas *big.Int) (*Quote, error)
}

// Wallet is a Safe as the database knows it.
type Wallet struct {
	Address common.Address
	// Owners and Threshold are the Safe's current owners (from the
	// database, checked against the chain once deployed).
	Owners    []common.Address
	Threshold int64
	// InitialOwners, InitialThreshold and SaltNonce are what the Safe was
	// (or will be) deployed with - they fix its address. Needed only
	// until it is deployed.
	InitialOwners    []common.Address
	InitialThreshold int64
	SaltNonce        *big.Int
	// InitialModules are extra modules enabled at deployment (besides the
	// Safe4337Module), also part of what fixes its address.
	InitialModules []common.Address
}

// RandomNonceKey returns a random 192-bit EntryPoint nonce key, for an
// operation that may wait (for approvers) alongside others on the same
// wallet.
func RandomNonceKey() (*big.Int, error) {
	b := make([]byte, 24)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}
	return new(big.Int).SetBytes(b), nil
}

// Request describes an operation to prepare.
type Request struct {
	Wallet Wallet
	Calls  []Call
	// GasToken pays gas through the paymaster; nil pays in ETH from the
	// wallet itself.
	GasToken *common.Address
	// Validity is how long the owners have to sign (0 = DefaultValidity).
	// Shared wallets waiting for approvers ask for longer.
	Validity time.Duration
	// NonceKey selects the EntryPoint nonce sequence (nil = key 0).
	// Operations that wait for approvers use their own key, so one pending
	// operation does not invalidate another.
	NonceKey *big.Int
}

// Prepared is an operation waiting for its owners' signatures. It is
// stored (JSON) between building it and submitting it.
type Prepared struct {
	Op         UserOperation    `json:"-"`
	OpRPC      RPC              `json:"op"`
	SafeOpHash common.Hash      `json:"safeOpHash"`
	UserOpHash common.Hash      `json:"userOpHash"`
	ValidAfter uint64           `json:"validAfter"`
	ValidUntil uint64           `json:"validUntil"`
	Owners     []common.Address `json:"owners"`
	Threshold  int64            `json:"threshold"`
	Activation bool             `json:"activation"` // this operation deploys the wallet
	GasToken   *common.Address  `json:"gasToken,omitempty"`
	Quote      *Quote           `json:"quote,omitempty"`
	MaxCostWei *hexutil.Big     `json:"maxCostWei"`
}

// Builder prepares and submits Safe UserOperations.
type Builder struct {
	Config  Config
	Chain   ChainReader
	Bundler Submitter
	Quotes  Quoter // nil disables stablecoin gas
	// Paymaster is TrovoTokenPaymaster (PAYMASTER_ADDRESS): the spender
	// wallets approve for their gas token.
	Paymaster common.Address
	Now       func() time.Time
	// DefaultValidity bounds how long a prepared operation stays signable.
	DefaultValidity time.Duration
	// GasBufferPercent is added to the bundler's gas estimates.
	GasBufferPercent int64
}

// Errors callers map to user-facing messages.
var (
	ErrInsufficientGasFunds = errors.New("aa: the wallet cannot pay for this operation's gas")
	ErrPaymasterUnavailable = errors.New("aa: paying gas in this token is unavailable right now")
)

var maxUint256 = new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 256), big.NewInt(1))

func (b *Builder) now() time.Time {
	if b.Now != nil {
		return b.Now()
	}
	return time.Now()
}

func withBuffer(v *big.Int, pct int64) *big.Int {
	if v == nil {
		return nil
	}
	out := new(big.Int).Mul(v, big.NewInt(100+pct))
	return out.Div(out, big.NewInt(100))
}

// Prepare builds the operation, estimates its gas, prices it (in the gas
// token, if any) and computes the hash the owners sign.
func (b *Builder) Prepare(ctx context.Context, req Request) (*Prepared, error) {
	w := req.Wallet
	if len(w.Owners) == 0 || w.Threshold < 1 || int(w.Threshold) > len(w.Owners) {
		return nil, fmt.Errorf("aa: wallet %s has an invalid owner set", w.Address.Hex())
	}
	deployed, err := Deployed(ctx, b.Chain, w.Address)
	if err != nil {
		return nil, err
	}
	var initCode []byte
	if !deployed {
		if w.SaltNonce == nil || len(w.InitialOwners) == 0 {
			return nil, fmt.Errorf("aa: wallet %s is not deployed and has no deployment parameters", w.Address.Hex())
		}
		if got := b.Config.SafeAddress(w.InitialOwners, w.InitialThreshold, w.SaltNonce, w.InitialModules...); got != w.Address {
			return nil, fmt.Errorf("aa: deployment parameters give %s, not wallet %s", got.Hex(), w.Address.Hex())
		}
		initCode = b.Config.InitCode(w.InitialOwners, w.InitialThreshold, w.SaltNonce, w.InitialModules...)
	}

	calls := req.Calls
	if req.GasToken != nil {
		if b.Quotes == nil || b.Paymaster == (common.Address{}) {
			return nil, ErrPaymasterUnavailable
		}
		// the paymaster pulls its fee with an allowance; the activation
		// operation (and the first one in a newly chosen token) grants it
		approve := !deployed
		if deployed {
			allowance, err := callView(ctx, b.Chain, *req.GasToken, "allowance", w.Address, b.Paymaster)
			if err != nil {
				return nil, err
			}
			approve = allowance[0].(*big.Int).Cmp(new(big.Int).Rsh(maxUint256, 1)) < 0
		}
		if approve {
			calls = append([]Call{ERC20Approve(*req.GasToken, b.Paymaster, maxUint256)}, calls...)
		}
	}
	callData, err := b.Config.CallData(calls)
	if err != nil {
		return nil, err
	}
	nonce, err := b.Config.EntryPointNonce(ctx, b.Chain, w.Address, req.NonceKey)
	if err != nil {
		return nil, err
	}
	tip, err := b.Chain.SuggestGasTipCap(ctx)
	if err != nil {
		return nil, err
	}
	head, err := b.Chain.HeaderByNumber(ctx, nil)
	if err != nil {
		return nil, err
	}
	baseFee := head.BaseFee
	if baseFee == nil {
		baseFee = big.NewInt(0)
	}
	maxFee := new(big.Int).Add(new(big.Int).Mul(baseFee, big.NewInt(2)), tip)

	op := UserOperation{
		Sender: w.Address, Nonce: nonce, InitCode: initCode, CallData: callData,
		// provisional limits for the first quote and the estimate
		VerificationGasLimit: big.NewInt(1_000_000), CallGasLimit: big.NewInt(1_000_000), PreVerificationGas: big.NewInt(200_000),
		MaxFeePerGas: maxFee, MaxPriorityFeePerGas: tip,
		Signature: DummySignature(int(w.Threshold)),
	}

	validity := req.Validity
	if validity <= 0 {
		validity = b.DefaultValidity
	}
	if validity <= 0 {
		validity = 10 * time.Minute
	}

	var quote *Quote
	if req.GasToken != nil {
		quote, err = b.Quotes.Quote(ctx, op, *req.GasToken, validity, nil, nil)
		if err != nil {
			return nil, b.quoteErr(err)
		}
		op.PaymasterAndData = quote.PaymasterAndData
	}

	est, err := b.Bundler.EstimateGas(ctx, op, b.Config.EntryPoint)
	if err != nil {
		return nil, fmt.Errorf("aa: estimating gas: %w", err)
	}
	op.VerificationGasLimit = withBuffer(est.VerificationGasLimit.ToInt(), b.GasBufferPercent)
	op.CallGasLimit = withBuffer(est.CallGasLimit.ToInt(), b.GasBufferPercent)
	op.PreVerificationGas = withBuffer(est.PreVerificationGas.ToInt(), b.GasBufferPercent)

	if req.GasToken != nil {
		var pmv, pmp *big.Int
		if est.PaymasterVerificationGasLimit != nil {
			pmv = withBuffer(est.PaymasterVerificationGasLimit.ToInt(), b.GasBufferPercent)
		}
		if est.PaymasterPostOpGasLimit != nil {
			pmp = withBuffer(est.PaymasterPostOpGasLimit.ToInt(), b.GasBufferPercent)
		}
		op.PaymasterAndData = nil
		quote, err = b.Quotes.Quote(ctx, op, *req.GasToken, validity, pmv, pmp)
		if err != nil {
			return nil, b.quoteErr(err)
		}
		op.PaymasterAndData = quote.PaymasterAndData
	} else {
		bal, err := balanceAt(ctx, b.Chain, w.Address)
		if err == nil && bal.Cmp(op.MaxCost()) < 0 {
			return nil, ErrInsufficientGasFunds
		}
	}

	now := uint64(b.now().Unix())
	validAfter := now - 60
	validUntil := now + uint64(validity/time.Second)
	if quote != nil && quote.ValidUntil != 0 && quote.ValidUntil < validUntil {
		validUntil = quote.ValidUntil
	}
	op.Signature = nil

	p := &Prepared{
		Op: op, SafeOpHash: b.Config.SafeOpHash(op, validAfter, validUntil),
		ValidAfter: validAfter, ValidUntil: validUntil,
		Owners: w.Owners, Threshold: w.Threshold, Activation: !deployed,
		GasToken: req.GasToken, Quote: quote, MaxCostWei: hb(op.MaxCost()),
	}
	p.OpRPC = op.ToRPC()
	p.UserOpHash = b.Config.UserOpHash(op)
	return p, nil
}

func (b *Builder) quoteErr(err error) error {
	var qe *QuoteError
	if errors.As(err, &qe) {
		return fmt.Errorf("%w: %s", ErrPaymasterUnavailable, qe.Message)
	}
	return fmt.Errorf("%w: %v", ErrPaymasterUnavailable, err)
}

func balanceAt(ctx context.Context, c ChainReader, addr common.Address) (*big.Int, error) {
	if br, ok := c.(interface {
		BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error)
	}); ok {
		return br.BalanceAt(ctx, addr, nil)
	}
	return nil, errors.New("balance unavailable")
}

// Submit adds the owners' signatures to a prepared operation and sends it
// to the bundler.
func (b *Builder) Submit(ctx context.Context, p *Prepared, sigs []OwnerSignature) (common.Hash, error) {
	if uint64(b.now().Unix()) > p.ValidUntil {
		return common.Hash{}, errors.New("aa: this operation expired before it was fully signed; start it again")
	}
	sig, err := AssembleSignature(p.SafeOpHash, p.ValidAfter, p.ValidUntil, p.Owners, int(p.Threshold), sigs)
	if err != nil {
		return common.Hash{}, err
	}
	op := p.Op
	op.Signature = sig
	return b.Bundler.Send(ctx, op, b.Config.EntryPoint)
}

// Marshal / UnmarshalPrepared store a prepared operation between the two
// halves of a request (build, then sign and submit).
func (p *Prepared) Marshal() ([]byte, error) {
	p.OpRPC = p.Op.ToRPC()
	return json.Marshal(p)
}

func UnmarshalPrepared(raw []byte) (*Prepared, error) {
	var p Prepared
	if err := json.Unmarshal(raw, &p); err != nil {
		return nil, err
	}
	p.Op = p.OpRPC.ToUserOperation()
	return &p, nil
}

// ToUserOperation converts the RPC form back.
func (r RPC) ToUserOperation() UserOperation {
	op := UserOperation{
		Sender: r.Sender, Nonce: r.Nonce.ToInt(), CallData: r.CallData,
		CallGasLimit: r.CallGasLimit.ToInt(), VerificationGasLimit: r.VerificationGasLimit.ToInt(),
		PreVerificationGas: r.PreVerificationGas.ToInt(), MaxFeePerGas: r.MaxFeePerGas.ToInt(),
		MaxPriorityFeePerGas: r.MaxPriorityFeePerGas.ToInt(), Signature: r.Signature,
	}
	if r.Factory != nil {
		op.InitCode = append(append([]byte{}, r.Factory.Bytes()...), r.FactoryData...)
	}
	if r.Paymaster != nil {
		gas := pack128(r.PaymasterVerificationGasLimit.ToInt(), r.PaymasterPostOpGasLimit.ToInt())
		op.PaymasterAndData = append(append(append([]byte{}, r.Paymaster.Bytes()...), gas...), r.PaymasterData...)
	}
	return op
}
