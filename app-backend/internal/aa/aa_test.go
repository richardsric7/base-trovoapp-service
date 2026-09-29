package aa

import (
	"context"
	"encoding/json"
	"errors"
	"math/big"
	"os"
	"strings"
	"testing"
	"time"

	"trovo-wallet-api/internal/evmkeypair"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

// hardhat's well-known test account #2 - the owner in the fixture
const fixtureOwnerKey = "0x5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a"

type safeFixture struct {
	Config struct {
		ProxyFactory, Singleton, ModuleSetup, Safe4337Module, EntryPoint, ProxyCreationCode string
	} `json:"config"`
	Addresses []struct {
		Owners      []string `json:"owners"`
		Threshold   int64    `json:"threshold"`
		SaltNonce   string   `json:"saltNonce"`
		Initializer string   `json:"initializer"`
		Address     string   `json:"address"`
	} `json:"addresses"`
	SafeOp struct {
		ChainID, Sender, Nonce, InitCode, CallData                 string
		VerificationGasLimit, CallGasLimit, PreVerificationGas     string
		MaxPriorityFeePerGas, MaxFeePerGas, PaymasterAndData, Hash string
		UserOpHash, OwnerSignatureWithV4                           string
		ValidAfter, ValidUntil                                     uint64
	} `json:"safeOp"`
}

func bigStr(s string) *big.Int {
	v, ok := new(big.Int).SetString(s, 10)
	if !ok {
		panic(s)
	}
	return v
}

func loadFixture(t *testing.T) (safeFixture, Config) {
	t.Helper()
	raw, err := os.ReadFile("testdata/safe_fixture.json")
	if err != nil {
		t.Fatal(err)
	}
	var f safeFixture
	if err := json.Unmarshal(raw, &f); err != nil {
		t.Fatal(err)
	}
	cfg := Config{
		ChainID:           bigStr(f.SafeOp.ChainID),
		EntryPoint:        common.HexToAddress(f.Config.EntryPoint),
		ProxyFactory:      common.HexToAddress(f.Config.ProxyFactory),
		Singleton:         common.HexToAddress(f.Config.Singleton),
		ModuleSetup:       common.HexToAddress(f.Config.ModuleSetup),
		Safe4337Module:    common.HexToAddress(f.Config.Safe4337Module),
		MultiSendCallOnly: common.HexToAddress(DefaultMultiSendCallOnly),
		ProxyCreationCode: hexutil.MustDecode(f.Config.ProxyCreationCode),
	}
	return f, cfg
}

func addrs(ss []string) []common.Address {
	out := make([]common.Address, len(ss))
	for i, s := range ss {
		out[i] = common.HexToAddress(s)
	}
	return out
}

// TestAddressesMatchRealDeployments checks the initializer and CREATE2
// address against Safes deployed by the real factory
// (paymaster/contracts/test/SafeWallet.test.js).
func TestAddressesMatchRealDeployments(t *testing.T) {
	f, cfg := loadFixture(t)
	if hexutil.Encode(cfg.ProxyCreationCode)[2:] != SafeProxyCreationCode {
		t.Fatal("the fixture factory's proxyCreationCode differs from SafeProxyCreationCode")
	}
	for _, c := range f.Addresses {
		owners := addrs(c.Owners)
		if got := hexutil.Encode(cfg.Initializer(owners, c.Threshold)); got != c.Initializer {
			t.Fatalf("initializer mismatch:\n got %s\nwant %s", got, c.Initializer)
		}
		if got := cfg.SafeAddress(owners, c.Threshold, bigStr(c.SaltNonce)); got != common.HexToAddress(c.Address) {
			t.Fatalf("address %s, factory deployed %s", got.Hex(), c.Address)
		}
		ic := cfg.InitCode(owners, c.Threshold, bigStr(c.SaltNonce))
		if common.BytesToAddress(ic[:20]) != cfg.ProxyFactory {
			t.Fatal("initCode must start with the factory")
		}
	}
}

func fixtureOp(f safeFixture) UserOperation {
	return UserOperation{
		Sender: common.HexToAddress(f.SafeOp.Sender), Nonce: bigStr(f.SafeOp.Nonce),
		InitCode: hexutil.MustDecode(f.SafeOp.InitCode), CallData: hexutil.MustDecode(f.SafeOp.CallData),
		VerificationGasLimit: bigStr(f.SafeOp.VerificationGasLimit), CallGasLimit: bigStr(f.SafeOp.CallGasLimit),
		PreVerificationGas: bigStr(f.SafeOp.PreVerificationGas), MaxPriorityFeePerGas: bigStr(f.SafeOp.MaxPriorityFeePerGas),
		MaxFeePerGas: bigStr(f.SafeOp.MaxFeePerGas), PaymasterAndData: hexutil.MustDecode(f.SafeOp.PaymasterAndData),
	}
}

// TestSafeOpHashAndSignatureMatchModule checks the SafeOp hash against the
// Safe4337Module's getOperationHash, the assembled signature against the one
// the real Safe accepted, and the userOpHash against the EntryPoint.
func TestSafeOpHashAndSignatureMatchModule(t *testing.T) {
	f, cfg := loadFixture(t)
	op := fixtureOp(f)
	h := cfg.SafeOpHash(op, f.SafeOp.ValidAfter, f.SafeOp.ValidUntil)
	if h.Hex() != f.SafeOp.Hash {
		t.Fatalf("SafeOp hash %s, module says %s", h.Hex(), f.SafeOp.Hash)
	}

	owner, _ := evmkeypair.ParseFull(fixtureOwnerKey)
	sig, err := evmkeypair.SignPersonal(owner.PrivateKey(), h.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	full, err := AssembleSignature(h, f.SafeOp.ValidAfter, f.SafeOp.ValidUntil, []common.Address{common.HexToAddress(owner.Address())}, 1,
		[]OwnerSignature{{Owner: common.HexToAddress(owner.Address()), Signature: sig}})
	if err != nil {
		t.Fatal(err)
	}
	if hexutil.Encode(full) != f.SafeOp.OwnerSignatureWithV4 {
		t.Fatalf("signature\n got %s\nwant %s", hexutil.Encode(full), f.SafeOp.OwnerSignatureWithV4)
	}

	op.Signature = full
	if got := cfg.UserOpHash(op); got.Hex() != f.SafeOp.UserOpHash {
		t.Fatalf("userOpHash %s, EntryPoint says %s", got.Hex(), f.SafeOp.UserOpHash)
	}
}

func TestAssembleSignatureRules(t *testing.T) {
	hash := crypto.Keccak256Hash([]byte("op"))
	keys := make([]*evmkeypair.Full, 3)
	owners := make([]common.Address, 3)
	sigs := make([]OwnerSignature, 3)
	for i := range keys {
		keys[i], _ = evmkeypair.Random()
		owners[i] = common.HexToAddress(keys[i].Address())
		s, _ := evmkeypair.SignPersonal(keys[i].PrivateKey(), hash.Bytes())
		sigs[i] = OwnerSignature{Owner: owners[i], Signature: s}
	}

	out, err := AssembleSignature(hash, 1, 2, owners, 2, []OwnerSignature{sigs[2], sigs[0], sigs[1]})
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != 12+2*65 {
		t.Fatalf("expected exactly threshold signatures, got %d bytes", len(out))
	}
	sorted := SortedOwners(owners)
	for i := 0; i < 2; i++ {
		part := append([]byte{}, out[12+i*65:12+(i+1)*65]...)
		if part[64] != 31 && part[64] != 32 {
			t.Fatalf("v must be 27/28 + 4, got %d", part[64])
		}
		part[64] -= 4
		if err := evmkeypair.VerifyPersonal(sorted[i], hash.Bytes(), part); err != nil {
			t.Fatalf("signature %d is not the sorted owner's: %v", i, err)
		}
	}

	stranger, _ := evmkeypair.Random()
	ss, _ := evmkeypair.SignPersonal(stranger.PrivateKey(), hash.Bytes())
	cases := map[string][]OwnerSignature{
		"non-owner":       {{Owner: common.HexToAddress(stranger.Address()), Signature: ss}},
		"duplicate":       {sigs[0], sigs[0]},
		"wrong hash":      {{Owner: owners[0], Signature: ss}},
		"below threshold": {sigs[0]},
	}
	for name, in := range cases {
		if _, err := AssembleSignature(hash, 0, 0, owners, 2, in); err == nil {
			t.Errorf("%s: expected an error", name)
		}
	}
	if _, err := AssembleSignature(hash, 0, 0, owners, 2, []OwnerSignature{sigs[0]}); !errors.Is(err, ErrNotEnoughSignatures) {
		t.Fatalf("expected ErrNotEnoughSignatures, got %v", err)
	}
}

func TestCallDataBatchesThroughMultiSend(t *testing.T) {
	cfg := ConfigFromEnv(big.NewInt(8453))
	token := common.HexToAddress("0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913")
	to := common.HexToAddress("0x1000000000000000000000000000000000000001")

	one, _ := cfg.CallData([]Call{ERC20Transfer(token, to, big.NewInt(5))})
	m, _ := contractsABI.MethodById(one[:4])
	args, _ := m.Inputs.Unpack(one[4:])
	if m.Name != "executeUserOp" || args[0].(common.Address) != token || args[3].(uint8) != 0 {
		t.Fatalf("single call must be a plain executeUserOp call, got %s %v", m.Name, args)
	}

	two, _ := cfg.CallData([]Call{ERC20Approve(token, to, big.NewInt(1)), ERC20Transfer(token, to, big.NewInt(5))})
	args, _ = m.Inputs.Unpack(two[4:])
	if args[0].(common.Address) != cfg.MultiSendCallOnly || args[3].(uint8) != 1 {
		t.Fatal("several calls must delegatecall MultiSendCallOnly")
	}
	if _, err := cfg.CallData(nil); err == nil {
		t.Fatal("no calls must be an error")
	}
}

func TestOwnerManagementCalls(t *testing.T) {
	safe := common.HexToAddress("0x5a6b0c1c6f1c0f4b1d7aa2d1c3e0bee8a9f2d4c1")
	a, b, c := common.HexToAddress("0xa"), common.HexToAddress("0xb"), common.HexToAddress("0xc")
	call, err := SwapOwner(safe, []common.Address{a, b}, b, c)
	if err != nil {
		t.Fatal(err)
	}
	m, _ := contractsABI.MethodById(call.Data[:4])
	args, _ := m.Inputs.Unpack(call.Data[4:])
	if m.Name != "swapOwner" || args[0].(common.Address) != a || args[1].(common.Address) != b || args[2].(common.Address) != c || call.To != safe {
		t.Fatalf("swapOwner must name the previous owner: %v", args)
	}
	call, _ = SwapOwner(safe, []common.Address{a, b}, a, c)
	args, _ = m.Inputs.Unpack(call.Data[4:])
	if args[0].(common.Address) != sentinelOwners {
		t.Fatal("the first owner's previous owner is the sentinel")
	}
	if _, err := RemoveOwner(safe, []common.Address{a}, b, 1); err == nil {
		t.Fatal("removing a non-owner must fail")
	}
	if OwnersHash([]common.Address{a, b}, 2) != OwnersHash([]common.Address{b, a}, 2) || OwnersHash([]common.Address{a, b}, 2) == OwnersHash([]common.Address{a, b}, 1) {
		t.Fatal("OwnersHash must ignore order and cover the threshold")
	}
}

// --- Builder with fakes ---------------------------------------------------

type fakeChain struct {
	code      map[common.Address]bool
	balance   *big.Int
	allowance *big.Int
	nonce     *big.Int
}

func (f *fakeChain) CodeAt(_ context.Context, a common.Address, _ *big.Int) ([]byte, error) {
	if f.code[a] {
		return []byte{1}, nil
	}
	return nil, nil
}

func (f *fakeChain) CallContract(_ context.Context, msg ethereum.CallMsg, _ *big.Int) ([]byte, error) {
	m, err := contractsABI.MethodById(msg.Data[:4])
	if err != nil {
		return nil, err
	}
	switch m.Name {
	case "getNonce":
		return m.Outputs.Pack(f.nonce)
	case "allowance":
		return m.Outputs.Pack(f.allowance)
	}
	return nil, errors.New("unexpected call " + m.Name)
}

func (f *fakeChain) SuggestGasTipCap(context.Context) (*big.Int, error) { return big.NewInt(1e6), nil }
func (f *fakeChain) HeaderByNumber(context.Context, *big.Int) (*types.Header, error) {
	return &types.Header{BaseFee: big.NewInt(1e7)}, nil
}
func (f *fakeChain) BalanceAt(context.Context, common.Address, *big.Int) (*big.Int, error) {
	return f.balance, nil
}

type fakeBundler struct {
	estimated UserOperation
	sent      *UserOperation
}

func (b *fakeBundler) EstimateGas(_ context.Context, op UserOperation, _ common.Address) (GasEstimate, error) {
	b.estimated = op
	return GasEstimate{
		PreVerificationGas: (*hexutil.Big)(big.NewInt(50000)), VerificationGasLimit: (*hexutil.Big)(big.NewInt(400000)),
		CallGasLimit: (*hexutil.Big)(big.NewInt(100000)), PaymasterVerificationGasLimit: (*hexutil.Big)(big.NewInt(90000)),
		PaymasterPostOpGasLimit: (*hexutil.Big)(big.NewInt(50000)),
	}, nil
}

func (b *fakeBundler) Send(_ context.Context, op UserOperation, _ common.Address) (common.Hash, error) {
	b.sent = &op
	return crypto.Keccak256Hash(op.Signature), nil
}

type fakeQuoter struct {
	calls int
	gas   [2]*big.Int
	err   error
}

func (q *fakeQuoter) Quote(_ context.Context, op UserOperation, token common.Address, validity time.Duration, pmv, pmp *big.Int) (*Quote, error) {
	q.calls++
	q.gas = [2]*big.Int{pmv, pmp}
	if q.err != nil {
		return nil, q.err
	}
	pm := common.HexToAddress("0x2000000000000000000000000000000000000002")
	if pmv == nil {
		pmv, pmp = big.NewInt(150000), big.NewInt(80000)
	}
	data := append(append(pm.Bytes(), pack128(pmv, pmp)...), token.Bytes()...)
	return &Quote{Paymaster: pm, PaymasterAndData: data, Token: token, ValidUntil: uint64(time.Unix(1790000000, 0).Add(validity).Unix())}, nil
}

func newTestBuilder(chain *fakeChain, q Quoter) (*Builder, *fakeBundler) {
	bun := &fakeBundler{}
	return &Builder{
		Config: ConfigFromEnv(big.NewInt(8453)), Chain: chain, Bundler: bun, Quotes: q,
		Paymaster: common.HexToAddress("0x2000000000000000000000000000000000000002"),
		Now:       func() time.Time { return time.Unix(1790000000, 0) }, GasBufferPercent: 20,
	}, bun
}

func TestPrepareActivationPaysInTokenAndApproves(t *testing.T) {
	owner, _ := evmkeypair.Random()
	ownerAddr := common.HexToAddress(owner.Address())
	chain := &fakeChain{code: map[common.Address]bool{}, balance: big.NewInt(0), allowance: big.NewInt(0), nonce: big.NewInt(0)}
	q := &fakeQuoter{}
	b, bun := newTestBuilder(chain, q)
	w := Wallet{Owners: []common.Address{ownerAddr}, Threshold: 1, InitialOwners: []common.Address{ownerAddr}, InitialThreshold: 1, SaltNonce: big.NewInt(0)}
	w.Address = b.Config.SafeAddress(w.InitialOwners, 1, w.SaltNonce)
	usdc := common.HexToAddress("0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913")

	p, err := b.Prepare(context.Background(), Request{Wallet: w, Calls: []Call{ERC20Transfer(usdc, common.HexToAddress("0xb0b"), big.NewInt(5))}, GasToken: &usdc})
	if err != nil {
		t.Fatal(err)
	}
	if !p.Activation || len(p.Op.InitCode) == 0 {
		t.Fatal("an undeployed wallet's operation must deploy it")
	}
	// approve + transfer batched through MultiSendCallOnly
	m, _ := contractsABI.MethodById(p.Op.CallData[:4])
	args, _ := m.Inputs.Unpack(p.Op.CallData[4:])
	if args[0].(common.Address) != b.Config.MultiSendCallOnly {
		t.Fatal("activation in a token must batch the paymaster approval with the call")
	}
	if !strings.Contains(hexutil.Encode(args[2].([]byte)), strings.ToLower(b.Paymaster.Hex()[2:])) {
		t.Fatal("the batch must approve the paymaster")
	}
	if q.calls != 2 || q.gas[0].Int64() != 108000 || q.gas[1].Int64() != 60000 {
		t.Fatalf("expected a provisional and a final quote with buffered paymaster gas, got %d calls %v", q.calls, q.gas)
	}
	if p.Op.VerificationGasLimit.Int64() != 480000 || p.Op.CallGasLimit.Int64() != 120000 || p.Op.PreVerificationGas.Int64() != 60000 {
		t.Fatalf("gas limits must be the buffered estimates: %v %v %v", p.Op.VerificationGasLimit, p.Op.CallGasLimit, p.Op.PreVerificationGas)
	}
	if len(bun.estimated.PaymasterAndData) == 0 || len(bun.estimated.Signature) != 12+65 {
		t.Fatal("estimation must use the provisional quote and a dummy signature")
	}
	if p.SafeOpHash != b.Config.SafeOpHash(p.Op, p.ValidAfter, p.ValidUntil) {
		t.Fatal("prepared hash mismatch")
	}

	// stored, reloaded, signed and submitted
	raw, _ := p.Marshal()
	loaded, err := UnmarshalPrepared(raw)
	if err != nil {
		t.Fatal(err)
	}
	if b.Config.SafeOpHash(loaded.Op, loaded.ValidAfter, loaded.ValidUntil) != p.SafeOpHash {
		t.Fatal("a reloaded operation must hash the same")
	}
	sig, _ := evmkeypair.SignPersonal(owner.PrivateKey(), p.SafeOpHash.Bytes())
	if _, err := b.Submit(context.Background(), loaded, []OwnerSignature{{Owner: ownerAddr, Signature: sig}}); err != nil {
		t.Fatal(err)
	}
	if bun.sent == nil || len(bun.sent.Signature) != 12+65 || bun.sent.Signature[12+64] < 31 {
		t.Fatal("the submitted operation must carry the assembled eth_sign signature")
	}
}

func TestPrepareInETHNeedsFunds(t *testing.T) {
	owner, _ := evmkeypair.Random()
	ownerAddr := common.HexToAddress(owner.Address())
	chain := &fakeChain{code: map[common.Address]bool{}, balance: big.NewInt(1), nonce: big.NewInt(0)}
	b, _ := newTestBuilder(chain, nil)
	w := Wallet{Owners: []common.Address{ownerAddr}, Threshold: 1, InitialOwners: []common.Address{ownerAddr}, InitialThreshold: 1, SaltNonce: big.NewInt(0)}
	w.Address = b.Config.SafeAddress(w.InitialOwners, 1, w.SaltNonce)
	_, err := b.Prepare(context.Background(), Request{Wallet: w, Calls: []Call{NativeTransfer(common.HexToAddress("0xb0b"), big.NewInt(1))}})
	if !errors.Is(err, ErrInsufficientGasFunds) {
		t.Fatalf("expected ErrInsufficientGasFunds, got %v", err)
	}

	chain.balance = big.NewInt(1e18)
	p, err := b.Prepare(context.Background(), Request{Wallet: w, Calls: []Call{NativeTransfer(common.HexToAddress("0xb0b"), big.NewInt(1))}})
	if err != nil || len(p.Op.PaymasterAndData) != 0 {
		t.Fatalf("an ETH-paid operation has no paymaster: %v", err)
	}

	// wrong deployment parameters are refused
	w.Address = common.HexToAddress("0xdead")
	if _, err := b.Prepare(context.Background(), Request{Wallet: w, Calls: []Call{NativeTransfer(common.HexToAddress("0xb0b"), big.NewInt(1))}}); err == nil {
		t.Fatal("expected a deployment-parameter mismatch error")
	}
}

func TestPrepareDeployedWithAllowanceSkipsApproval(t *testing.T) {
	owner, _ := evmkeypair.Random()
	ownerAddr := common.HexToAddress(owner.Address())
	w := Wallet{Address: common.HexToAddress("0x5a6b0c1c6f1c0f4b1d7aa2d1c3e0bee8a9f2d4c1"), Owners: []common.Address{ownerAddr}, Threshold: 1}
	chain := &fakeChain{code: map[common.Address]bool{w.Address: true}, allowance: maxUint256, nonce: big.NewInt(4)}
	b, _ := newTestBuilder(chain, &fakeQuoter{})
	usdc := common.HexToAddress("0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913")
	p, err := b.Prepare(context.Background(), Request{Wallet: w, Calls: []Call{ERC20Transfer(usdc, common.HexToAddress("0xb0b"), big.NewInt(5))}, GasToken: &usdc})
	if err != nil {
		t.Fatal(err)
	}
	m, _ := contractsABI.MethodById(p.Op.CallData[:4])
	args, _ := m.Inputs.Unpack(p.Op.CallData[4:])
	if p.Activation || len(p.Op.InitCode) != 0 || args[0].(common.Address) != usdc || p.Op.Nonce.Int64() != 4 {
		t.Fatal("a deployed wallet with an allowance sends the call alone, with its current nonce")
	}

	// the quote service being unavailable surfaces as ErrPaymasterUnavailable
	b.Quotes = &fakeQuoter{err: &QuoteError{Status: 503, Code: "rate-unavailable", Message: "stale"}}
	if _, err := b.Prepare(context.Background(), Request{Wallet: w, Calls: []Call{ERC20Transfer(usdc, common.HexToAddress("0xb0b"), big.NewInt(5))}, GasToken: &usdc}); !errors.Is(err, ErrPaymasterUnavailable) {
		t.Fatalf("expected ErrPaymasterUnavailable, got %v", err)
	}
}

// TestBaseAddressesMatchWalletCore pins the canonical-Base addresses
// wallet-core's primarySafeAddress computes (src/safe.rs), so the app and
// the backend always agree on a user's address.
func TestBaseAddressesMatchWalletCore(t *testing.T) {
	cfg := ConfigFromEnv(big.NewInt(8453))
	owner := []common.Address{common.HexToAddress("0x70997970C51812dc3A010C7d01b50e0d17dc79C8")}
	for salt, want := range map[int64]string{0: "0xe7a9D4D8a9633bea8f6f891C5D98744356A9259F", 1: "0x3A741746d076eCF518186E8644a70C0982Dd584E"} {
		if got := cfg.SafeAddress(owner, 1, big.NewInt(salt)).Hex(); got != want {
			t.Fatalf("salt %d: %s, wallet-core computes %s", salt, got, want)
		}
	}
}
