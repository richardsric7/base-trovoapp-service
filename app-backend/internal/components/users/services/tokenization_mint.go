package users

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"math/big"
	"time"

	"trovo-wallet-api/internal/basetxn"
	userModels "trovo-wallet-api/internal/components/users/models"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/evmkeypair"
	"trovo-wallet-api/internal/gnosissafe"
	"trovo-wallet-api/internal/network"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
)

// tokenizationMintPlan is what minting a tokenized asset does on Base: the
// issuing Safe calls mint() on the asset's token contract - the unsold
// supply to the Safe itself (the asset's treasury) and the platform's
// fee-in-asset to the tokenization fee wallet - in one atomic Safe
// transaction.
type tokenizationMintPlan struct {
	ChainID        *big.Int
	Safe           common.Address
	Token          common.Address
	Decimals       uint8
	Treasury       common.Address
	TreasuryAmount *big.Int // token base units
	FeeWallet      common.Address
	FeeAmount      *big.Int // token base units; may be zero
}

// newTokenizationMintPlan converts the asset's human-unit supply and fee
// into token base units using the contract's own decimals.
func newTokenizationMintPlan(chainID *big.Int, safe, token, feeWallet common.Address, decimals uint8, totalSupply, feeInAsset decimal.Decimal) (tokenizationMintPlan, error) {
	if !totalSupply.IsPositive() {
		return tokenizationMintPlan{}, fmt.Errorf("number of tokens to be issued must be positive")
	}
	if feeInAsset.IsNegative() || feeInAsset.GreaterThan(totalSupply) {
		return tokenizationMintPlan{}, fmt.Errorf("fee in asset (%v) must be between 0 and the issued supply (%v)", feeInAsset, totalSupply)
	}
	total := totalSupply.Shift(int32(decimals))
	fee := feeInAsset.Shift(int32(decimals))
	if !total.Equal(total.Truncate(0)) || !fee.Equal(fee.Truncate(0)) {
		return tokenizationMintPlan{}, fmt.Errorf("issued supply %v / fee %v have more precision than the token's %d decimals", totalSupply, feeInAsset, decimals)
	}
	return tokenizationMintPlan{
		ChainID:        chainID,
		Safe:           safe,
		Token:          token,
		Decimals:       decimals,
		Treasury:       safe,
		TreasuryAmount: total.Sub(fee).BigInt(),
		FeeWallet:      feeWallet,
		FeeAmount:      fee.BigInt(),
	}, nil
}

// Description is the exact, human-readable statement of the mint that
// minting approvers sign, and that is re-derived and compared at execution
// time so the Safe only ever executes what was approved.
func (p tokenizationMintPlan) Description() string {
	d := fmt.Sprintf("Safe %s on chain %s executes %s.mint(%s, %s)", p.Safe.Hex(), p.ChainID.String(), p.Token.Hex(), p.Treasury.Hex(), p.TreasuryAmount.String())
	if p.FeeAmount.Sign() > 0 {
		d += fmt.Sprintf(" and %s.mint(%s, %s)", p.Token.Hex(), p.FeeWallet.Hex(), p.FeeAmount.String())
	}
	return d + fmt.Sprintf(" [amounts in base units, %d decimals]", p.Decimals)
}

// ApprovalPayload is the base64 form of Description that the apps' approval
// screens decode and sign (see signBase64Txn in app-web/app-mobile).
func (p tokenizationMintPlan) ApprovalPayload() string {
	return base64.StdEncoding.EncodeToString([]byte(p.Description()))
}

func (p tokenizationMintPlan) Calls() ([]gnosissafe.Call, error) {
	amounts := []struct {
		to     common.Address
		amount *big.Int
	}{{p.Treasury, p.TreasuryAmount}, {p.FeeWallet, p.FeeAmount}}
	calls := make([]gnosissafe.Call, 0, 2)
	for _, a := range amounts {
		if a.amount.Sign() == 0 {
			continue
		}
		data, err := tokenizedAssetTokenABI.Pack("mint", a.to, a.amount)
		if err != nil {
			return nil, err
		}
		calls = append(calls, gnosissafe.Call{To: p.Token, Value: big.NewInt(0), Data: data})
	}
	return calls, nil
}

// buildTokenizationMintPlan derives the mint plan for t from its current
// state and on-chain token decimals.
func buildTokenizationMintPlan(ctx context.Context, t *userModels.TokenizedAsset, feeWallet string, gc *sharedconfig.GlobalConfig) (tokenizationMintPlan, error) {
	token, err := tokenizedAssetContract(t)
	if err != nil {
		return tokenizationMintPlan{}, err
	}
	if t.IssuingWalletAddress == nil || !common.IsHexAddress(*t.IssuingWalletAddress) {
		return tokenizationMintPlan{}, &tErrors.CustomError{Param: "issuingWalletAddress", Err: "error-no-issuing-wallet", ErrMessage: "Issuing Safe not assigned."}
	}
	decimals, err := network.AssetDecimals(ctx, gc.BantuExpansionClient, basetxn.CreditAsset{Code: *t.AssetCode, Issuer: token})
	if err != nil {
		log.Printf("[buildTokenizationMintPlan] reading decimals of %v: %v\n", token, err)
		return tokenizationMintPlan{}, &tErrors.ErrorTemporaryServerError{}
	}
	plan, err := newTokenizationMintPlan(network.GetBlockchainChainID(), common.HexToAddress(*t.IssuingWalletAddress), common.HexToAddress(token), common.HexToAddress(feeWallet),
		decimals, decimal.NewFromFloat(t.NumberOfTokenToBeIssued), decimal.NewFromFloat(t.FeeInAsset))
	if err != nil {
		return tokenizationMintPlan{}, &tErrors.CustomError{Param: "numberOfTokenToBeIssued", Err: "error-invalid-mint-amounts", ErrMessage: err.Error(), Code: 400}
	}
	return plan, nil
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

// executeTokenizationMint performs an approved TOKENIZE ASSET request: it
// re-derives the mint plan, refuses unless it is exactly what the approvers
// signed (approvedPayload), refuses if the token already has supply (the
// mint already happened - so a retried approval can never mint twice), then
// has the issuing Safe execute it and waits for it to be mined.
func executeTokenizationMint(tkInput *userModels.TokenMinting, approvedPayload string, gc *sharedconfig.GlobalConfig) (network.SubmittedTransaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	ta, _, err := GetTokenizedAssetByID(tkInput.TokenizedAssetID, gc.DB)
	if err != nil {
		return network.SubmittedTransaction{}, err
	}
	feeWallet := ta.GetTokenizationFeeWallet(gc)
	feeKP, err := evmkeypair.ParseFull(feeWallet.FeeWalletSecretKey)
	if err != nil {
		log.Printf("[executeTokenizationMint] tokenization fee wallet not configured: %v\n", err)
		return network.SubmittedTransaction{}, &tErrors.ErrorTemporaryServerError{}
	}
	plan, err := buildTokenizationMintPlan(ctx, &ta, feeKP.Address(), gc)
	if err != nil {
		return network.SubmittedTransaction{}, err
	}
	if plan.ApprovalPayload() != approvedPayload {
		log.Printf("[executeTokenizationMint] plan for %v changed since approval was requested:\napproved: %q\ncurrent:  %q\n", ta.ID, approvedPayload, plan.ApprovalPayload())
		return network.SubmittedTransaction{}, &tErrors.CustomError{Param: "id", Err: "error-please-reject-transaction", ErrMessage: "The minting details changed after this approval was requested. Please reject this request and initiate minting again."}
	}

	supply, err := tokenSupply(ctx, gc.BantuExpansionClient, plan.Token)
	if err != nil {
		log.Printf("[executeTokenizationMint] reading totalSupply of %v: %v\n", plan.Token.Hex(), err)
		return network.SubmittedTransaction{}, &tErrors.ErrorTemporaryServerError{}
	}
	if supply.Sign() != 0 {
		gc.LogDiscordFailedRequest(fmt.Sprintf("[executeTokenizationMint] refusing to mint %v (%v): token %v already has supply %v - reconcile manually", *ta.AssetCode, ta.ID, plan.Token.Hex(), supply))
		return network.SubmittedTransaction{}, &tErrors.CustomError{Param: "id", Err: "error-already-minted", ErrMessage: "This asset's token already has supply on-chain, so it will not be minted again. Please contact engineering to reconcile."}
	}

	candidates, err := issuingSafeSigners()
	if err != nil {
		log.Printf("[executeTokenizationMint] %v\n", err)
		return network.SubmittedTransaction{}, &tErrors.ErrorTemporaryServerError{}
	}
	signers, err := gnosissafe.SignersForSafe(ctx, gc.BantuExpansionClient, plan.Safe, candidates)
	if err != nil {
		log.Printf("[executeTokenizationMint] %v\n", err)
		gc.LogDiscordFailedRequest(fmt.Sprintf("[executeTokenizationMint] cannot sign for issuing Safe %v: %v", plan.Safe.Hex(), err))
		return network.SubmittedTransaction{}, &tErrors.ErrorTemporaryServerError{}
	}
	calls, err := plan.Calls()
	if err != nil {
		return network.SubmittedTransaction{}, &tErrors.ErrorTemporaryServerError{}
	}

	hash, err := gnosissafe.ExecCalls(ctx, gc.BantuExpansionClient, plan.ChainID, plan.Safe, signers, calls, multiSendCallOnlyAddress())
	if err != nil {
		log.Printf("[executeTokenizationMint] executing mint of %v via Safe %v: %v\n", *ta.AssetCode, plan.Safe.Hex(), err)
		return network.SubmittedTransaction{}, &tErrors.CustomError{Param: "id", Err: "error operation failed", ErrMessage: "Minting transaction could not be submitted.", Code: 500}
	}
	// Once broadcast, a retry is safe: it either finds the supply already
	// minted (above) or re-signs the same Safe nonce, which can only execute once.
	if err := gnosissafe.WaitSuccess(ctx, gc.BantuExpansionClient, common.HexToHash(hash)); err != nil {
		log.Printf("[executeTokenizationMint] mint %v of %v not confirmed: %v\n", hash, *ta.AssetCode, err)
		gc.LogDiscordFailedRequest(fmt.Sprintf("[executeTokenizationMint] mint tx %v for %v not confirmed: %v", hash, *ta.AssetCode, err))
		return network.SubmittedTransaction{Hash: hash}, &tErrors.CustomError{Param: "id", Err: "error operation failed", ErrMessage: fmt.Sprintf("Minting transaction %v was not confirmed: %v", hash, err), Code: 500}
	}
	log.Printf("[executeTokenizationMint] minted %v via Safe %v: %v\n", *ta.AssetCode, plan.Safe.Hex(), hash)
	return network.SubmittedTransaction{Hash: hash}, nil
}
