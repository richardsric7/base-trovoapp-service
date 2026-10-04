package users

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"trovo-wallet-api/internal/aa"
	"trovo-wallet-api/internal/basetxn"
	userModels "trovo-wallet-api/internal/components/users/models"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/network"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm/clause"
)

// OperationTokenizedAssetEarlyExit is the wallet operation kind of an early
// exit: the holder's tokens go back to the asset's distribution wallet.
const OperationTokenizedAssetEarlyExit = "TOKENIZED ASSET EARLY EXIT"

type earlyExitContext struct {
	Input userModels.TokenizedAssetEarlyExitInput `json:"input"`
}

// EarlyExit processes a holder's early exit (pre-maturity redemption) from a
// tokenized asset: one operation of the holder's wallet transfers the
// exiting tokens to the asset's distribution wallet (there is no on-chain
// buy-back), and the payout to the given bank account is settled off-chain
// from the recorded TokenizedAssetEarlyExit. Two steps, as with payments:
// without transactionSignature it returns the operation as transaction;
// with it, it submits and records the exit. Shared wallets with approvers
// get an approval request on commit.
func EarlyExit(initiator *userModels.User, wallet *userModels.UserWallet, ta *userModels.TokenizedAsset, input *userModels.TokenizedAssetEarlyExitInput, gc *sharedconfig.GlobalConfig) (ee userModels.TokenizedAssetEarlyExit, err error) {
	if err = initiator.EnsureNotSuspended(); err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	input.TokenizedAssetID = ta.ID
	input.WalletAddress = wallet.ID
	input.TokenQuantityToExit = decimal.NewFromFloat(input.TokenQuantityToExit).Truncate(7).InexactFloat64()
	input.NetworkPassPhrase = network.GetBlockchainNetworkPassPhrase()
	if wallet.SharedAccessEnabled == 1 && wallet.NumberOfApprovalsNeeded > 0 {
		input.Multiparty = 1
	}
	if wallet.HasViewOnlyAccess(gc) {
		input.SignatureRequired = 1
	}
	walletOwner, e := wallet.GetWalletOwner(gc.DB, gc)
	if e != nil {
		return ee, &tErrors.CustomError{Param: "walletAddress", Err: "error-invalid-kyc", ErrMessage: "Unable to verify wallet owner."}
	}

	// a prepared exit: submit it, or ask the approvers
	if len(input.Transaction) > 0 && (len(input.TransactionSignature) > 0 && input.Multiparty == 0 || input.Multiparty == 1 && input.Commit == 1) {
		rec, p, e := LoadWalletOperation(input.Transaction, wallet.ID, OperationTokenizedAssetEarlyExit, gc)
		if e != nil {
			return ee, e
		}
		var c earlyExitContext
		if rec.Context == nil || json.Unmarshal([]byte(*rec.Context), &c) != nil {
			return ee, &tErrors.ErrorTemporaryServerError{}
		}
		stored := c.Input
		stored.Transaction, stored.TransactionSignature, stored.Commit, stored.Multiparty = input.Transaction, input.TransactionSignature, input.Commit, input.Multiparty
		*input = stored
		var bank userModels.Bank
		gc.DB.Where("id = ?", input.BankID).First(&bank)
		ee = buildTokenizedAssetEarlyExit(walletOwner.Username, wallet.ID, input, ta, &bank, gc)
		if input.Multiparty == 1 {
			return ee, createEarlyExitApprovalRequest(initiator, wallet, ta, input, &ee, gc)
		}
		hash, e := SignSingleOwnerOperation(ctx, rec, p, initiator.PrimarySigner, input.TransactionSignature, gc)
		if e != nil {
			return ee, e
		}
		input.TransactionID, ee.TransactionID = hash, hash
		if e := gc.DB.Omit(clause.Associations).Create(&ee).Error; e != nil {
			log.Printf("[EarlyExit] operation %v submitted but early exit not saved: %v", hash, e)
			gc.LogDiscordFailedRequest(fmt.Sprintf("[EarlyExit] early exit %v of %v %v by %v submitted but not recorded: %v", hash, input.TokenQuantityToExit, *ta.AssetCode, wallet.ID, e))
		}
		wallet.InvalidateUserCache(gc)
		initiator.InvalidateUserWalletCache(gc)
		return ee, nil
	}

	if walletOwner.KYCVerified == 0 {
		return ee, &tErrors.CustomError{Param: "walletAddress", Err: "error-invalid-kyc", ErrMessage: fmt.Sprintf("%v has not passed KYC to exit this tokenized asset %v.", walletOwner.Username, *ta.AssetCode)}
	}
	if ta.AssetTokenizationStatus != 5 && ta.AssetTokenizationStatus != 6 {
		return ee, &tErrors.CustomError{Param: "tokenizedAssetId", Err: "error-invalid-request", ErrMessage: "Only projects that are on sale or trading can accept an early exit."}
	}
	if !ta.MaturityDate.IsZero() && !time.Now().Before(ta.MaturityDate) {
		return ee, &tErrors.CustomError{Param: "tokenizedAssetId", Err: "error-asset-matured", ErrMessage: "This asset has reached maturity. Please use the standard redemption instead of an early exit."}
	}
	if input.TokenQuantityToExit <= 0 {
		return ee, &tErrors.CustomError{Param: "tokenQuantityToExit", Err: "error-invalid-amount", ErrMessage: "You must specify a token quantity greater than zero to exit."}
	}
	var bank userModels.Bank
	if e := gc.DB.Where("id = ?", input.BankID).First(&bank).Error; e != nil {
		return ee, &tErrors.CustomError{Param: "bankId", Err: "error-invalid-bank", ErrMessage: "Please select a valid bank for payout."}
	}
	if ta.MarketMakingWallet == nil || !common.IsHexAddress(*ta.MarketMakingWallet) {
		return ee, &tErrors.CustomError{Param: "tokenizedAssetId", Err: "error-no-distribution-wallet", ErrMessage: "This asset has no distribution wallet."}
	}
	token, err := tokenizedAssetContract(ta)
	if err != nil {
		return
	}
	decimals, e := network.AssetDecimals(ctx, gc.BantuExpansionClient, basetxn.CreditAsset{Code: *ta.AssetCode, Issuer: token})
	if e != nil {
		return ee, &tErrors.ErrorTemporaryServerError{}
	}
	amount := decimal.NewFromFloat(input.TokenQuantityToExit)
	if bal, e := network.B20BalanceOf(gc.BantuExpansionClient, token, wallet.ID, decimals); e != nil {
		return ee, &tErrors.ErrorTemporaryServerError{}
	} else if bal.LessThan(amount) {
		return ee, &tErrors.ErrorUnderfundedAccount{Detail: fmt.Sprintf("You only have %v %v available; cannot exit %v %v.", bal, *ta.AssetCode, amount, *ta.AssetCode)}
	}
	units, err := baseUnits(amount.String(), decimals)
	if err != nil {
		return
	}
	validity := time.Duration(0)
	if input.Multiparty == 1 {
		validity = sharedWalletOperationValidity()
	}
	ee = buildTokenizedAssetEarlyExit(walletOwner.Username, wallet.ID, input, ta, &bank, gc)
	input.Memo = "EARLYEXIT"
	call := aa.ERC20Transfer(common.HexToAddress(token), common.HexToAddress(*ta.MarketMakingWallet), units)
	op, err := PrepareWalletOperation(ctx, OperationTokenizedAssetEarlyExit, initiator, &walletOwner, wallet, []aa.Call{call}, validity, earlyExitContext{Input: *input}, gc)
	if err != nil {
		return
	}
	input.Transaction = op.Transaction
	input.SignatureRequired = 1
	input.Messages = append([]string{fmt.Sprintf("%v %v return to the asset's distribution wallet; about %v %v is paid out to your bank account.", amount, *ta.AssetCode, decimal.NewFromFloat(ee.EstimatedPayoutAmount), ee.PayoutCurrency)}, op.Messages()...)
	if input.Multiparty == 1 && input.Commit == 1 {
		err = createEarlyExitApprovalRequest(initiator, wallet, ta, input, &ee, gc)
	}
	return
}

// createEarlyExitApprovalRequest records a shared wallet's early exit for
// its approvers.
func createEarlyExitApprovalRequest(initiator *userModels.User, wallet *userModels.UserWallet, ta *userModels.TokenizedAsset, input *userModels.TokenizedAssetEarlyExitInput, ee *userModels.TokenizedAssetEarlyExit, gc *sharedconfig.GlobalConfig) error {
	input.TransactionID, ee.TransactionID = "PENDING_AUTH", "PENDING_AUTH"
	description := fmt.Sprintf("Early exit of Tokenized asset [%v]\nQuantity: %v %v,\nEstimated payout: %v %v", *ta.AssetName, decimal.NewFromFloat(input.TokenQuantityToExit), *ta.AssetCode, decimal.NewFromFloat(ee.EstimatedPayoutAmount), *ta.AssetQuoteCurrency)
	if len(input.Messages) > 0 {
		description = fmt.Sprintf("%v\nMessages: %v", description, strings.Join(input.Messages, "\n"))
	}
	input.ReturnedDescription = description
	info, _ := json.Marshal(*input)
	infoStr := string(info)
	pendingAuth := userModels.PendingAuth{
		ID:                     uuid.NewString(),
		Initiator:              initiator.Username,
		InitiatorSignerAddress: initiator.PrimarySigner,
		WalletAddress:          wallet.ID,
		TransactionType:        OperationTokenizedAssetEarlyExit,
		Description:            description,
		ApprovalsNeeded:        wallet.NumberOfApprovalsNeeded,
		TransactionXdr:         input.Transaction,
		TransactionInfoStr:     &infoStr,
	}
	if err := gc.DB.Omit(clause.Associations).Create(&pendingAuth).Error; err != nil {
		log.Printf("[EarlyExit] saving approval request for %v: %v", wallet.ID, err)
		return &tErrors.ErrorTemporaryServerError{}
	}
	wallet.InvalidateUserCache(gc)
	initiator.InvalidateUserWalletCache(gc)
	return nil
}
