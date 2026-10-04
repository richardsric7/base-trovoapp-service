package users

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"strings"
	"testing"
	"time"

	"trovo-wallet-api/internal/aa"
	userModels "trovo-wallet-api/internal/components/users/models"
	"trovo-wallet-api/internal/evmkeypair"
	"trovo-wallet-api/internal/gnosissafe"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/shopspring/decimal"
)

// TestTokenizationOnLocalChain runs the tokenization flows against real
// contracts: the local stack (paymaster/contracts/scripts/local-stack.js,
// with its dev bundler) and the market contracts on the same node
// (market/contracts/scripts/local-market.js). Skipped unless
//
//	AA_LOCAL_STACK=$PWD/../../../../../paymaster/contracts/deployments/local-stack.json
//	AA_LOCAL_MARKET=$PWD/../../../../../market/contracts/deployments/local-market.json
//
// (absolute paths) are set. It mints an asset the way a mint approval does
// (the issuing Safe's first operation, signed by two of four minting
// approvers: deploy the distribution Safe, mint, open the sale offer), buys from it with cNGN from a new wallet, checks a replayed
// authorization is refused, and delivers a fiat purchase from the internal
// balance minting Safe.
func TestTokenizationOnLocalChain(t *testing.T) {
	stackFile, marketFile := os.Getenv("AA_LOCAL_STACK"), os.Getenv("AA_LOCAL_MARKET")
	if stackFile == "" || marketFile == "" {
		t.Skip("AA_LOCAL_STACK / AA_LOCAL_MARKET not set")
	}
	read := func(f string) map[string]string {
		raw, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		m := map[string]string{}
		if err := json.Unmarshal(raw, &m); err != nil {
			t.Fatal(err)
		}
		return m
	}
	stack, market := read(stackFile), read(marketFile)
	for _, k := range []string{"ENTRYPOINT_ADDRESS", "SAFE_PROXY_FACTORY_ADDRESS", "SAFE_SINGLETON_ADDRESS", "SAFE_MODULE_SETUP_ADDRESS", "SAFE_4337_MODULE_ADDRESS", "SAFE_MULTISEND_CALL_ONLY_ADDRESS"} {
		t.Setenv(k, stack[k])
	}
	t.Setenv("BASE_CHAIN_ID", stack["BASE_CHAIN_ID"])
	t.Setenv("OFFER_BOOK_ADDRESS", market["OFFER_BOOK_ADDRESS"])
	t.Setenv("OFFER_AUTHORIZER_PRIVATE_KEY", "0x5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a") // hardhat #2

	ctx := context.Background()
	client, err := ethclient.Dial("http://127.0.0.1:8545")
	if err != nil {
		t.Fatal(err)
	}
	chainID, _ := client.ChainID(ctx)
	gc := &sharedconfig.GlobalConfig{BantuExpansionClient: client}
	cfg := aa.ConfigFromEnv(chainID)
	b := &aa.Builder{Config: cfg, Chain: client, Bundler: aa.NewBundler(stack["BUNDLER_URL"]), GasBufferPercent: 300, DefaultValidity: 10 * time.Minute} // the dev bundler's fixed estimates are too low for a mint
	book := common.HexToAddress(market["OFFER_BOOK_ADDRESS"])
	cngn := common.HexToAddress(market["MOCK_CNGN_ADDRESS"])
	internal := common.HexToAddress(market["MOCK_INTERNAL_BALANCE_ADDRESS"])

	// hardhat #0: funds wallets and owns the offer book
	admin, _ := evmkeypair.ParseFull("0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
	tx := func(to common.Address, value *big.Int, data []byte) {
		t.Helper()
		if value.Sign() > 0 {
			h, err := sendValue(ctx, client, chainID, admin, to, value)
			if err != nil {
				t.Fatal(err)
			}
			if err := gnosissafe.WaitSuccess(ctx, client, h); err != nil {
				t.Fatal(err)
			}
			return
		}
		h, err := gnosissafe.SendTransaction(ctx, client, chainID, admin, to, data)
		if err != nil {
			t.Fatal(err)
		}
		if err := gnosissafe.WaitSuccess(ctx, client, common.HexToHash(h)); err != nil {
			t.Fatal(err)
		}
	}
	erc20, _ := abi.JSON(strings.NewReader(`[{"name":"mint","type":"function","inputs":[{"type":"address"},{"type":"uint256"}],"outputs":[]},
		{"name":"balanceOf","type":"function","stateMutability":"view","inputs":[{"type":"address"}],"outputs":[{"type":"uint256"}]},
		{"name":"setTradable","type":"function","inputs":[{"type":"address"},{"type":"bool"}],"outputs":[]}]`))
	balance := func(token, holder common.Address) *big.Int {
		t.Helper()
		data, _ := erc20.Pack("balanceOf", holder)
		out, err := client.CallContract(ctx, ethereum.CallMsg{To: &token, Data: data}, nil)
		if err != nil {
			t.Fatal(err)
		}
		return new(big.Int).SetBytes(out)
	}
	mint := func(token, to common.Address, amount *big.Int) {
		data, _ := erc20.Pack("mint", to, amount)
		tx(token, big.NewInt(0), data)
	}
	run := func(label string, w aa.Wallet, calls []aa.Call, validity time.Duration, signers ...*evmkeypair.Full) *aa.Receipt {
		t.Helper()
		key, _ := aa.RandomNonceKey()
		p, err := b.Prepare(ctx, aa.Request{Wallet: w, Calls: calls, Validity: validity, NonceKey: key})
		if err != nil {
			t.Fatalf("%s: prepare: %v", label, err)
		}
		var sigs []aa.OwnerSignature
		for _, k := range signers {
			sig, _ := evmkeypair.SignPersonal(k.PrivateKey(), p.SafeOpHash.Bytes())
			sigs = append(sigs, aa.OwnerSignature{Owner: common.HexToAddress(k.Address()), Signature: sig})
		}
		h, err := b.Submit(ctx, p, sigs)
		if err != nil {
			t.Fatalf("%s: submit: %v", label, err)
		}
		wctx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()
		r, err := b.Bundler.(*aa.Bundler).WaitReceipt(wctx, h)
		if err != nil {
			t.Fatalf("%s: %v", label, err)
		}
		t.Logf("%s: success=%v tx=%s", label, r.Success, r.Receipt.TransactionHash.Hex())
		return r
	}

	// --- the asset: issuing and distribution Safes owned by the issuing
	// profile's key and four minting approvers (2 approvals needed)
	keys := make([]*evmkeypair.Full, 5)
	users := make([]*userModels.User, 5)
	for i := range keys {
		keys[i], _ = evmkeypair.Random()
		users[i] = &userModels.User{PrimarySigner: keys[i].Address()}
	}
	issuingD, distD := issuingSafeDeployments(users[0], users[1:], fmt.Sprintf("e2e-%d", time.Now().UnixNano()))
	issuing, dist := common.HexToAddress(issuingD.Address), common.HexToAddress(distD.Address)
	if issuingD.Threshold != 2 || len(issuingD.Owners) != 5 || distD.Threshold != 2 || !strings.EqualFold(distD.Modules[0], issuing.Hex()) {
		t.Fatalf("issuing %+v distribution %+v", issuingD, distD)
	}

	// its token, owned by the issuing Safe, listed on the book
	artifactRaw, err := os.ReadFile(market["TOKENIZED_ASSET_ARTIFACT"])
	if err != nil {
		t.Fatal(err)
	}
	var art struct {
		ABI      json.RawMessage `json:"abi"`
		Bytecode string          `json:"bytecode"`
	}
	json.Unmarshal(artifactRaw, &art)
	tokenABI, _ := abi.JSON(strings.NewReader(string(art.ABI)))
	ctorArgs, _ := tokenABI.Pack("", "Farm Fund", "FARM", uint8(7), issuing)
	token := deployContract(t, ctx, client, chainID, admin, append(common.FromHex(art.Bytecode), ctorArgs...))
	setTradable, _ := erc20.Pack("setTradable", token, true)
	tx(book, big.NewInt(0), setTradable)

	// --- mint: the issuing Safe's first operation, paying its own gas
	tx(issuing, big.NewInt(1e17), nil)
	fundsKey, _ := evmkeypair.Random()
	feeKey, _ := evmkeypair.Random()
	funds, feeWallet := common.HexToAddress(fundsKey.Address()), common.HexToAddress(feeKey.Address())
	plan, err := newTokenizationMintPlan(chainID, issuing, dist, token, feeWallet, 7, decimal.NewFromInt(1000), decimal.NewFromInt(10), decimal.NewFromInt(600))
	if err != nil {
		t.Fatal(err)
	}
	plan.Book, plan.Proceeds = book, funds
	for _, pt := range []struct {
		code string
		addr common.Address
		dec  uint8
	}{{"CNGN", cngn, 6}, {"NGNI", internal, 18}} {
		p, err := offerPrice(decimal.NewFromInt(1500), 7, pt.dec)
		if err != nil {
			t.Fatal(err)
		}
		plan.PaymentCodes, plan.PaymentTokens, plan.Prices = append(plan.PaymentCodes, pt.code), append(plan.PaymentTokens, pt.addr), append(plan.Prices, p)
	}
	salt, _ := new(big.Int).SetString(distD.SaltNonce, 10)
	deploy := cfg.DeploySafeCall(distD.OwnerAddresses(), int64(distD.Threshold), salt, distD.ModuleAddresses()...)
	plan.DeployDistribution = &deploy
	calls, err := plan.Calls()
	if err != nil {
		t.Fatal(err)
	}
	issuingSalt, _ := new(big.Int).SetString(issuingD.SaltNonce, 10)
	iw := aa.Wallet{Address: issuing, Owners: issuingD.OwnerAddresses(), Threshold: 2, InitialOwners: issuingD.OwnerAddresses(), InitialThreshold: 2, SaltNonce: issuingSalt}
	r := run("mint", iw, calls, 24*time.Hour, keys[1], keys[3])
	if !r.Success {
		t.Fatal("mint reverted")
	}
	rcpt, err := client.TransactionReceipt(ctx, r.Receipt.TransactionHash)
	if err != nil {
		t.Fatal(err)
	}
	offerID, ok := aa.ParseOfferCreated(rcpt.Logs, book, issuing)
	if !ok {
		t.Fatal("no offer created")
	}
	unit := new(big.Int).Exp(big.NewInt(10), big.NewInt(7), nil)
	tokens := func(n int64) *big.Int { return new(big.Int).Mul(big.NewInt(n), unit) }
	for _, c := range []struct {
		who  common.Address
		want *big.Int
	}{{dist, tokens(390)}, {feeWallet, tokens(10)}, {book, tokens(600)}, {issuing, tokens(0)}} {
		if got := balance(token, c.who); got.Cmp(c.want) != 0 {
			t.Fatalf("%s holds %v, want %v", c.who.Hex(), got, c.want)
		}
	}

	// --- crypto purchase: 3000 cNGN buys 2 tokens, in one operation of a
	// new buyer wallet (its activation)
	idStr, code, tokenHex := offerID.String(), "FARM", token.Hex()
	ta := &userModels.TokenizedAsset{AssetCode: &code, ContractAddress: &tokenHex, OfferBookOfferID: &idStr}
	o, err := loadAssetOffer(ctx, ta, gc)
	if err != nil {
		t.Fatal(err)
	}
	buyerKey, _ := evmkeypair.Random()
	buyerOwner := common.HexToAddress(buyerKey.Address())
	buyer := cfg.SafeAddress([]common.Address{buyerOwner}, 1, big.NewInt(0))
	mint(cngn, buyer, big.NewInt(5_000_000_000))
	tx(buyer, big.NewInt(1e17), nil)
	purchase, err := priceAssetPurchase(ctx, o, "CNGN", cngn, decimal.NewFromInt(3000), buyer, buyer, 10*time.Minute, gc)
	if err != nil {
		t.Fatal(err)
	}
	if purchase.AssetAmount.Cmp(tokens(2)) != 0 || purchase.Payment.Cmp(big.NewInt(3_000_000_000)) != 0 {
		t.Fatalf("priced %v for %v", purchase.AssetAmount, purchase.Payment)
	}
	pcalls, _ := purchase.purchaseCalls(book)
	bw := aa.Wallet{Address: buyer, Owners: []common.Address{buyerOwner}, Threshold: 1, InitialOwners: []common.Address{buyerOwner}, InitialThreshold: 1, SaltNonce: big.NewInt(0)}
	if r := run("purchase", bw, pcalls, 0, buyerKey); !r.Success {
		t.Fatal("purchase reverted")
	}
	if got := balance(token, buyer); got.Cmp(tokens(2)) != 0 {
		t.Fatalf("buyer holds %v", got)
	}
	if got := balance(cngn, funds); got.Cmp(big.NewInt(3_000_000_000)) != 0 {
		t.Fatalf("proceeds %v", got)
	}
	// a replayed authorization buys nothing
	bw.InitialOwners, bw.InitialThreshold = nil, 0
	if r := run("replayed purchase", bw, pcalls, 0, buyerKey); r.Success {
		t.Fatal("a used authorization must not fill again")
	}
	if got := balance(token, buyer); got.Cmp(tokens(2)) != 0 {
		t.Fatalf("buyer holds %v after replay", got)
	}

	// --- fiat purchase: the internal balance minting Safe pays 1500 for a
	// buyer, as DeliverFiatAssetPurchase does
	minterKey, _ := evmkeypair.Random()
	tx(common.HexToAddress(minterKey.Address()), big.NewInt(1e17), nil)
	minter, err := gnosissafe.DeploySafe(ctx, client, chainID, admin, gnosissafe.DeployConfig{
		ProxyFactory: cfg.ProxyFactory, Singleton: cfg.Singleton,
		FallbackHandler: cfg.Safe4337Module, // any contract; the local stack has no CompatibilityFallbackHandler
	}, []common.Address{common.HexToAddress(minterKey.Address())}, 1, big.NewInt(time.Now().UnixNano()))
	if err != nil {
		t.Fatal(err)
	}
	fiatBuyerKey, _ := evmkeypair.Random()
	fiatBuyer := common.HexToAddress(fiatBuyerKey.Address())
	fp, err := priceAssetPurchase(ctx, o, "NGNI", internal, decimal.NewFromInt(1500), minter, fiatBuyer, 10*time.Minute, gc)
	if err != nil {
		t.Fatal(err)
	}
	fill, _ := aa.FillCall(book, fp.Fill, fp.Authorization)
	var scalls []gnosissafe.Call
	for _, c := range []aa.Call{aa.ERC20Mint(internal, minter, fp.Payment), aa.ERC20Approve(internal, book, fp.Payment), fill} {
		scalls = append(scalls, gnosissafe.Call{To: c.To, Value: big.NewInt(0), Data: c.Data})
	}
	h, err := gnosissafe.ExecCalls(ctx, client, chainID, minter, []*evmkeypair.Full{minterKey}, scalls, cfg.MultiSendCallOnly)
	if err != nil {
		t.Fatal(err)
	}
	if err := gnosissafe.WaitSuccess(ctx, client, common.HexToHash(h)); err != nil {
		t.Fatal(err)
	}
	if got := balance(token, fiatBuyer); got.Cmp(tokens(1)) != 0 {
		t.Fatalf("fiat buyer holds %v", got)
	}
	if got := balance(internal, funds); got.Cmp(new(big.Int).Mul(big.NewInt(1500), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))) != 0 {
		t.Fatalf("internal balance proceeds %v", got)
	}
	left, _ := aa.ReadOffer(ctx, client, book, offerID)
	if left.Remaining.Cmp(tokens(597)) != 0 {
		t.Fatalf("offer has %v left", left.Remaining)
	}
}

// sendValue sends ETH from from to to.
func sendValue(ctx context.Context, client *ethclient.Client, chainID *big.Int, from *evmkeypair.Full, to common.Address, value *big.Int) (common.Hash, error) {
	return sendRaw(ctx, client, chainID, from, &to, value, nil, 100000)
}

// deployContract deploys creation code from from and returns its address.
func deployContract(t *testing.T, ctx context.Context, client *ethclient.Client, chainID *big.Int, from *evmkeypair.Full, code []byte) common.Address {
	t.Helper()
	h, err := sendRaw(ctx, client, chainID, from, nil, big.NewInt(0), code, 3_000_000)
	if err != nil {
		t.Fatal(err)
	}
	if err := gnosissafe.WaitSuccess(ctx, client, h); err != nil {
		t.Fatal(err)
	}
	r, err := client.TransactionReceipt(ctx, h)
	if err != nil {
		t.Fatal(err)
	}
	return r.ContractAddress
}

func sendRaw(ctx context.Context, client *ethclient.Client, chainID *big.Int, from *evmkeypair.Full, to *common.Address, value *big.Int, data []byte, gas uint64) (common.Hash, error) {
	sender := common.HexToAddress(from.Address())
	nonce, err := client.PendingNonceAt(ctx, sender)
	if err != nil {
		return common.Hash{}, err
	}
	tip, _ := client.SuggestGasTipCap(ctx)
	head, err := client.HeaderByNumber(ctx, nil)
	if err != nil {
		return common.Hash{}, err
	}
	tx := types.NewTx(&types.DynamicFeeTx{ChainID: chainID, Nonce: nonce, To: to, Value: value, Data: data, Gas: gas,
		GasTipCap: tip, GasFeeCap: new(big.Int).Add(new(big.Int).Mul(head.BaseFee, big.NewInt(2)), tip)})
	signed, err := types.SignTx(tx, types.LatestSignerForChainID(chainID), from.PrivateKey())
	if err != nil {
		return common.Hash{}, err
	}
	return signed.Hash(), client.SendTransaction(ctx, signed)
}
