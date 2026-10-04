package users

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"strings"
	"time"

	"trovo-wallet-api/internal/aa"
	"trovo-wallet-api/internal/basetxn"
	userModels "trovo-wallet-api/internal/components/users/models"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/network"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"gorm.io/gorm/clause"
)

// OperationTokenizationMint is the wallet operation kind of a tokenized
// asset's mint, made by its issuing Safe once the minting approvers sign.
const OperationTokenizationMint = "TOKENIZE ASSET"

// tokenizationMintContext is what the mint operation keeps for when it is
// mined.
type tokenizationMintContext struct {
	TokenizedAssetID string `json:"tokenizedAssetId"`
	Seller           string `json:"seller"`
	Book             string `json:"book"`
}

// tokenizationMintPlan is what minting a tokenized asset does on Base, in
// one operation of its issuing Safe:
//
//  1. deploy the distribution Safe (with the issuing Safe as its module),
//     unless it already exists;
//  2. mint ForSale to itself, the supply not for sale (less the
//     platform's fee-in-asset) to the distribution Safe, and the fee to
//     the tokenization fee wallet;
//  3. put ForSale on sale in TrovoOfferBook at the asset's price in each
//     payment token, with sale proceeds paid to the funds holding wallet.
//
// The issuing Safe is the offer's seller (it has the same owners as the
// distribution Safe), so every call is the Safe's own and a failing one
// reverts the whole operation - a call made through the distribution
// Safe's module would only report failure, not revert.
type tokenizationMintPlan struct {
	ChainID            *big.Int
	Issuing            common.Address
	Distribution       common.Address
	Token              common.Address
	Decimals           uint8
	DistributionAmount *big.Int // token base units not for sale; may be zero
	FeeWallet          common.Address
	FeeAmount          *big.Int // token base units; may be zero
	ForSale            *big.Int // token base units
	Book               common.Address
	Proceeds           common.Address
	PaymentCodes       []string
	PaymentTokens      []common.Address
	Prices             []aa.OfferPrice
	// DeployDistribution deploys the distribution Safe (nil when it exists).
	DeployDistribution *aa.Call
}

// newTokenizationMintPlan converts the asset's human-unit supply, fee and
// amount for sale into token base units using the contract's decimals.
func newTokenizationMintPlan(chainID *big.Int, issuing, distribution, token, feeWallet common.Address, decimals uint8, totalSupply, feeInAsset, forSale decimal.Decimal) (tokenizationMintPlan, error) {
	if !totalSupply.IsPositive() {
		return tokenizationMintPlan{}, fmt.Errorf("number of tokens to be issued must be positive")
	}
	if feeInAsset.IsNegative() || feeInAsset.GreaterThan(totalSupply) {
		return tokenizationMintPlan{}, fmt.Errorf("fee in asset (%v) must be between 0 and the issued supply (%v)", feeInAsset, totalSupply)
	}
	if !forSale.IsPositive() || forSale.GreaterThan(totalSupply.Sub(feeInAsset)) {
		return tokenizationMintPlan{}, fmt.Errorf("tokens for sale (%v) must be positive and at most the issued supply less the fee (%v)", forSale, totalSupply.Sub(feeInAsset))
	}
	total, fee, sale := totalSupply.Shift(int32(decimals)), feeInAsset.Shift(int32(decimals)), forSale.Shift(int32(decimals))
	for _, d := range []decimal.Decimal{total, fee, sale} {
		if !d.Equal(d.Truncate(0)) {
			return tokenizationMintPlan{}, fmt.Errorf("issued supply %v / fee %v / for sale %v have more precision than the token's %d decimals", totalSupply, feeInAsset, forSale, decimals)
		}
	}
	return tokenizationMintPlan{
		ChainID:            chainID,
		Issuing:            issuing,
		Distribution:       distribution,
		Token:              token,
		Decimals:           decimals,
		DistributionAmount: total.Sub(fee).Sub(sale).BigInt(),
		FeeWallet:          feeWallet,
		FeeAmount:          fee.BigInt(),
		ForSale:            sale.BigInt(),
	}, nil
}

func (p tokenizationMintPlan) human(v *big.Int) string {
	return decimal.NewFromBigInt(v, -int32(p.Decimals)).String()
}

// Description states the mint for the minting approvers.
func (p tokenizationMintPlan) Description() string {
	d := fmt.Sprintf("Issuing wallet %s mints %s of token %s and offers it for sale on %s (proceeds to %s)", p.Issuing.Hex(), p.human(p.ForSale), p.Token.Hex(), p.Book.Hex(), p.Proceeds.Hex())
	if p.DistributionAmount.Sign() > 0 {
		d += fmt.Sprintf(", mints %s not for sale to the distribution wallet %s", p.human(p.DistributionAmount), p.Distribution.Hex())
	}
	if p.FeeAmount.Sign() > 0 {
		d += fmt.Sprintf(", mints %s to the tokenization fee wallet %s", p.human(p.FeeAmount), p.FeeWallet.Hex())
	}
	d += "; price"
	for i, code := range p.PaymentCodes {
		if i > 0 {
			d += ","
		}
		d += fmt.Sprintf(" %s/%s %s base units per token base unit", p.Prices[i].Num, p.Prices[i].Den, code)
	}
	return d + fmt.Sprintf(" [chain %s]", p.ChainID)
}

// Calls are the issuing Safe's calls carrying out the plan.
func (p tokenizationMintPlan) Calls() ([]aa.Call, error) {
	var calls []aa.Call
	if p.DeployDistribution != nil {
		calls = append(calls, *p.DeployDistribution)
	}
	calls = append(calls, aa.ERC20Mint(p.Token, p.Issuing, p.ForSale))
	if p.DistributionAmount.Sign() > 0 {
		calls = append(calls, aa.ERC20Mint(p.Token, p.Distribution, p.DistributionAmount))
	}
	if p.FeeAmount.Sign() > 0 {
		calls = append(calls, aa.ERC20Mint(p.Token, p.FeeWallet, p.FeeAmount))
	}
	offer, err := aa.CreateOfferCall(p.Book, p.Token, p.ForSale, p.Proceeds, p.PaymentTokens, p.Prices)
	if err != nil {
		return nil, err
	}
	return append(calls, aa.ERC20Approve(p.Token, p.Book, p.ForSale), offer), nil
}

// salePaymentTokens are the tokens the asset's sale accepts, each pegged
// 1:1 to its quote currency: the tokenization payment stablecoins and the
// country's internal balance token (which fiat purchases pay in).
func salePaymentTokens(t *userModels.TokenizedAsset, gc *sharedconfig.GlobalConfig) (codes []string, tokens []common.Address) {
	var currencies []userModels.TokenizationCurrency
	gc.DB.Find(&currencies)
	seen := map[common.Address]bool{}
	add := func(code, addr string) {
		if !common.IsHexAddress(addr) {
			return
		}
		a := common.HexToAddress(addr)
		if !seen[a] {
			seen[a] = true
			codes = append(codes, strings.ToUpper(code))
			tokens = append(tokens, a)
		}
	}
	for _, c := range currencies {
		add(c.AssetCode, c.ContractAddress)
	}
	if t.AssetCountryLocation != nil {
		cc := userModels.CountryCode(*t.AssetCountryLocation).GetConfig(gc)
		if cc.InternalBalanceTokenCode != nil && cc.InternalTokenIssuer != nil {
			add(*cc.InternalBalanceTokenCode, *cc.InternalTokenIssuer)
		}
	}
	return codes, tokens
}

// buildTokenizationMintPlan derives the mint plan for t from its current
// state and the chain.
func buildTokenizationMintPlan(ctx context.Context, t *userModels.TokenizedAsset, gc *sharedconfig.GlobalConfig) (tokenizationMintPlan, *userModels.UserWallet, error) {
	fail := func(param, code, msg string) (tokenizationMintPlan, *userModels.UserWallet, error) {
		return tokenizationMintPlan{}, nil, &tErrors.CustomError{Param: param, Err: code, ErrMessage: msg, Code: http.StatusBadRequest}
	}
	token, err := tokenizedAssetContract(t)
	if err != nil {
		return tokenizationMintPlan{}, nil, err
	}
	if t.IssuingWalletAddress == nil {
		return fail("issuingWalletAddress", "error-no-issuing-wallet", "Issuing wallet not assigned.")
	}
	issuing, err := userModels.UserWalletID(*t.IssuingWalletAddress).GetWallet(gc.DB, gc)
	if err != nil || issuing.LinkedWalletAddress == nil {
		return fail("issuingWalletAddress", "error-no-issuing-wallet", "Issuing wallet not assigned.")
	}
	dist, err := userModels.UserWalletID(*issuing.LinkedWalletAddress).GetWallet(gc.DB, gc)
	if err != nil {
		return fail("issuingWalletAddress", "error-no-distribution-wallet", "The issuing wallet has no distribution wallet.")
	}
	if t.FundsHoldingWalletAddress == nil || !common.IsHexAddress(*t.FundsHoldingWalletAddress) {
		return fail("fundsHoldingWalletAddress", "error-funds-holding-wallet-not-set", "Set the funds holding wallet, which receives the sale proceeds, before minting.")
	}
	feeWallet, err := feeWalletAddress(os.Getenv("TOKENIZATION_FEE_WALLET"), "tokenization fee wallet", gc)
	if err != nil {
		return tokenizationMintPlan{}, nil, err
	}
	book, err := offerBookAddress()
	if err != nil {
		return tokenizationMintPlan{}, nil, err
	}
	client := gc.BantuExpansionClient
	decimals, err := network.AssetDecimals(ctx, client, basetxn.CreditAsset{Code: *t.AssetCode, Issuer: token})
	if err != nil {
		log.Printf("[buildTokenizationMintPlan] reading decimals of %v: %v", token, err)
		return tokenizationMintPlan{}, nil, &tErrors.ErrorTemporaryServerError{}
	}
	plan, err := newTokenizationMintPlan(network.GetBlockchainChainID(), common.HexToAddress(issuing.ID), common.HexToAddress(dist.ID), common.HexToAddress(token), common.HexToAddress(feeWallet),
		decimals, decimal.NewFromFloat(t.NumberOfTokenToBeIssued), decimal.NewFromFloat(t.FeeInAsset), decimal.NewFromFloat(t.MaxNumberOfTokenAvailableForSale))
	if err != nil {
		return fail("numberOfTokenToBeIssued", "error-invalid-mint-amounts", err.Error())
	}
	plan.Book = book
	plan.Proceeds = common.HexToAddress(*t.FundsHoldingWalletAddress)

	// the book must list the asset and its payment tokens (the book's owner
	// lists them; see market/INTEGRATION.md)
	if ok, err := aa.Tradable(ctx, client, book, plan.Token); err != nil {
		return tokenizationMintPlan{}, nil, &tErrors.ErrorTemporaryServerError{}
	} else if !ok {
		return fail("contractAddress", "error-token-not-listed", fmt.Sprintf("The offer book does not list %v (%v) yet. Its owner must list it (setTradable) before minting.", *t.AssetCode, token))
	}
	price := decimal.NewFromFloat(t.PricePerToken)
	codes, tokens := salePaymentTokens(t, gc)
	for i, pt := range tokens {
		if ok, err := aa.Tradable(ctx, client, book, pt); err != nil || !ok {
			log.Printf("[buildTokenizationMintPlan] %v (%v) is not tradable on the offer book; not accepted for %v", codes[i], pt.Hex(), *t.AssetCode)
			continue
		}
		payDec, err := network.AssetDecimals(ctx, client, basetxn.CreditAsset{Code: codes[i], Issuer: pt.Hex()})
		if err != nil {
			return tokenizationMintPlan{}, nil, &tErrors.ErrorTemporaryServerError{}
		}
		p, err := offerPrice(price, decimals, payDec)
		if err != nil {
			return fail("pricePerToken", "error-invalid-price", err.Error())
		}
		plan.PaymentCodes = append(plan.PaymentCodes, codes[i])
		plan.PaymentTokens = append(plan.PaymentTokens, pt)
		plan.Prices = append(plan.Prices, p)
	}
	if len(plan.PaymentTokens) == 0 {
		return fail("assetQuoteCurrency", "error-no-payment-token", "None of the sale's payment tokens is listed on the offer book.")
	}

	if deployed, err := aa.Deployed(ctx, client, plan.Distribution); err != nil {
		return tokenizationMintPlan{}, nil, &tErrors.ErrorTemporaryServerError{}
	} else if !deployed {
		owners := make([]common.Address, 0)
		for _, o := range dist.InitialOwnerList() {
			owners = append(owners, common.HexToAddress(o))
		}
		modules := make([]common.Address, 0)
		for _, m := range dist.InitialModuleList() {
			modules = append(modules, common.HexToAddress(m))
		}
		salt, ok := new(big.Int).SetString(dist.SafeSaltNonce, 10)
		if !ok {
			return tokenizationMintPlan{}, nil, &tErrors.ErrorTemporaryServerError{}
		}
		call := network.AAConfig().DeploySafeCall(owners, int64(dist.InitialThreshold), salt, modules...)
		plan.DeployDistribution = &call
	}
	return plan, &issuing, nil
}

// tokenSupply reads the token contract's totalSupply().
func tokenSupply(ctx context.Context, r contractReader, token common.Address) (*big.Int, error) {
	res, err := readToken(ctx, r, token, "totalSupply")
	if err != nil {
		return nil, err
	}
	supply, ok := res[0].(*big.Int)
	if !ok {
		return nil, fmt.Errorf("unexpected totalSupply() result type")
	}
	return supply, nil
}

// prepareTokenizationMint builds the issuing Safe's mint operation for the
// minting approvers to sign. It is refused if the token already has supply
// (the asset was minted), so an asset is minted once.
func prepareTokenizationMint(ctx context.Context, initiator *userModels.User, t *userModels.TokenizedAsset, gc *sharedconfig.GlobalConfig) (*PreparedWalletOperation, tokenizationMintPlan, *userModels.UserWallet, error) {
	plan, issuing, err := buildTokenizationMintPlan(ctx, t, gc)
	if err != nil {
		return nil, plan, nil, err
	}
	if supply, err := tokenSupply(ctx, gc.BantuExpansionClient, plan.Token); err != nil {
		log.Printf("[prepareTokenizationMint] reading totalSupply of %v: %v", plan.Token.Hex(), err)
		return nil, plan, nil, &tErrors.ErrorTemporaryServerError{}
	} else if supply.Sign() != 0 {
		return nil, plan, nil, &tErrors.CustomError{Param: "contractAddress", Err: "error-already-minted", ErrMessage: "This asset's token contract already has supply on-chain.", Code: http.StatusBadRequest}
	}
	calls, err := plan.Calls()
	if err != nil {
		return nil, plan, nil, &tErrors.ErrorTemporaryServerError{}
	}
	profile, err := issuing.GetWalletOwner(gc.DB, gc)
	if err != nil {
		return nil, plan, nil, &tErrors.ErrorTemporaryServerError{}
	}
	op, err := PrepareWalletOperation(ctx, OperationTokenizationMint, initiator, &profile, issuing, calls, sharedWalletOperationValidity(),
		tokenizationMintContext{TokenizedAssetID: t.ID, Seller: plan.Issuing.Hex(), Book: plan.Book.Hex()}, gc)
	if err != nil {
		if ce, ok := err.(*tErrors.CustomError); ok && ce.Err == "error-insufficient-network-fee" {
			ce.ErrMessage = fmt.Sprintf("The issuing wallet %v needs ETH (or a stablecoin the issuing profile pays network fees in) to pay for minting. Send some to it and try again.", issuing.ID)
		}
		return nil, plan, nil, err
	}
	deploys := []string{}
	if plan.DeployDistribution != nil {
		deploys = append(deploys, strings.ToUpper(plan.Distribution.Hex()))
	}
	gc.DB.Model(&userModels.WalletOperation{}).Where("id = ?", op.Record.ID).Update("deploys", strings.Join(deploys, ","))
	return op, plan, issuing, nil
}

// recordTokenizationMint stores the sale offer a mined mint created; the
// asset counts as minted from then on.
func recordTokenizationMint(op userModels.WalletOperation, txHash common.Hash, gc *sharedconfig.GlobalConfig) {
	var c tokenizationMintContext
	if op.Context == nil || json.Unmarshal([]byte(*op.Context), &c) != nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	rcpt, err := gc.BantuExpansionClient.TransactionReceipt(ctx, txHash)
	if err != nil {
		log.Printf("[recordTokenizationMint] receipt of %v: %v", txHash.Hex(), err)
		gc.LogDiscordFailedRequest(fmt.Sprintf("[recordTokenizationMint] mint of tokenization %v mined in %v but its receipt could not be read: %v", c.TokenizedAssetID, txHash.Hex(), err))
		return
	}
	id, ok := aa.ParseOfferCreated(rcpt.Logs, common.HexToAddress(c.Book), common.HexToAddress(c.Seller))
	if !ok {
		gc.LogDiscordFailedRequest(fmt.Sprintf("[recordTokenizationMint] mint of tokenization %v (%v) created no offer", c.TokenizedAssetID, txHash.Hex()))
		return
	}
	s := id.String()
	if err := gc.DB.Model(&userModels.TokenizedAsset{}).Where("id = ?", c.TokenizedAssetID).Omit(clause.Associations).Update("offer_book_offer_id", s).Error; err != nil {
		log.Printf("[recordTokenizationMint] saving offer %v of %v: %v", s, c.TokenizedAssetID, err)
		return
	}
	log.Printf("[recordTokenizationMint] tokenization %v minted (%v), sale offer %v", c.TokenizedAssetID, txHash.Hex(), s)
}

func init() {
	minedHooks[OperationTokenizationMint] = recordTokenizationMint
}
