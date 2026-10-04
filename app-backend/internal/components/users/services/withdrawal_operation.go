package users

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"trovo-wallet-api/internal/aa"
	"trovo-wallet-api/internal/basetxn"
	assetModels "trovo-wallet-api/internal/components/assets/models"
	paymentModels "trovo-wallet-api/internal/components/payments/models"
	userModels "trovo-wallet-api/internal/components/users/models"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/network"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm/clause"
)

// OperationCryptoWithdrawal is the wallet operation kind for crypto
// withdrawals through 1Liquidity.
const OperationCryptoWithdrawal = "CRYPTO WITHDRAWAL"

// withdrawalContext is what a withdrawal operation saves when submitted.
type withdrawalContext struct {
	Request userModels.WithdrawalRequest `json:"request"`
}

// queueWithdrawalOperation is the on-chain half of a crypto withdrawal.
// Crypto deposits are minted as Trovo tokens, so a withdrawal burns them
// from the user's wallet (the token contract must be ERC20Burnable); the
// withdrawal request is then paid out off-chain.
//
//  1. without transactionSignature: builds the wallet's operation burning
//     amountSubmitted and returns it as transaction for the user to sign;
//  2. with transaction + transactionSignature: records the withdrawal
//     request and submits the operation.
//
// Shared wallets with approvers get an approval request on commit.
func queueWithdrawalOperation(signerUser *userModels.User, wallet *userModels.UserWallet, wdlInput *userModels.WithdrawalRequestInput, ca assetModels.CuratedAsset, serviceFee decimal.Decimal, gc *sharedconfig.GlobalConfig) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	wdlInput.NetworkPassPhrase = network.GetBlockchainNetworkPassPhrase()

	if wdlInput.Multiparty == 0 && len(wdlInput.TransactionSignature) > 0 && len(wdlInput.Transaction) > 0 {
		return submitWithdrawal(ctx, signerUser, wallet, wdlInput, gc)
	}
	if wdlInput.Multiparty == 1 && wdlInput.Commit == 1 && len(wdlInput.Transaction) > 0 {
		if _, _, err := LoadWalletOperation(wdlInput.Transaction, wallet.ID, OperationCryptoWithdrawal, gc); err != nil {
			return err
		}
		return createWithdrawalApprovalRequest(signerUser, wallet, wdlInput, serviceFee, gc)
	}

	if !common.IsHexAddress(ca.ContractAddress) {
		return &tErrors.CustomError{Param: "currency", Err: "error-currency-not-withdrawable", ErrMessage: wdlInput.Currency + " cannot be withdrawn.", Code: http.StatusBadRequest}
	}
	asset := basetxn.CreditAsset{Code: ca.AssetCode, Issuer: ca.ContractAddress}
	decimals, err := network.AssetDecimals(ctx, gc.BantuExpansionClient, asset)
	if err != nil {
		return &tErrors.ErrorTemporaryServerError{}
	}
	amount := decimal.NewFromFloat(wdlInput.AmountSubmitted)
	if bal, err := network.B20BalanceOf(gc.BantuExpansionClient, ca.ContractAddress, wallet.ID, decimals); err != nil {
		return &tErrors.ErrorTemporaryServerError{}
	} else if bal.LessThan(amount) {
		return &tErrors.CustomError{Param: "submittedAmount", Err: "error insufficient balance for " + wdlInput.Currency, ErrMessage: "Wallet does not have enough balance to withdraw " + wdlInput.Currency}
	}
	units, err := baseUnits(amount.String(), decimals)
	if err != nil {
		return err
	}

	owner, err := wallet.GetWalletOwner(gc.DB, gc)
	if err != nil {
		return &tErrors.ErrorTemporaryServerError{}
	}
	request := userModels.WithdrawalRequest{
		WalletAddress:        wallet.ID,
		WalletAlias:          wallet.Alias,
		UserID:               wallet.UserID,
		Currency:             wdlInput.Currency,
		AmountSubmitted:      wdlInput.AmountSubmitted,
		AmountToWithdraw:     wdlInput.AmountToWithdraw,
		WithdrawalAddress:    wdlInput.WithdrawalAddress,
		WithdrawalMemo:       wdlInput.WithdrawalMemo,
		WithdrawalNetwork:    wdlInput.WithdrawalNetwork,
		WithdrawalServiceFee: wdlInput.WithdrawalServiceFee,
		WithdrawalNetworkFee: wdlInput.WithdrawalNetworkFee,
	}
	validity := time.Duration(0)
	if wdlInput.Multiparty == 1 {
		validity = sharedWalletOperationValidity()
	}
	op, err := PrepareWalletOperation(ctx, OperationCryptoWithdrawal, signerUser, &owner, wallet, []aa.Call{aa.ERC20Burn(common.HexToAddress(ca.ContractAddress), units)}, validity, withdrawalContext{Request: request}, gc)
	if err != nil {
		return err
	}
	// the burn and a network fee paid in the same token must both fit
	burn := paymentModels.PaymentInfo{AssetCode: ca.AssetCode, ContractAddress: ca.ContractAddress, AmountToPay: amount.String()}
	if err := checkPaymentLeavesGas(ctx, wallet, &burn, op.Prepared, gc); err != nil {
		return err
	}
	wdlInput.Transaction = op.Transaction
	wdlInput.SignatureRequired = 1
	wdlInput.Messages = append([]string{fmt.Sprintf("%v %v leaves this wallet; %v %v is sent to %v on %v.", amount, wdlInput.Currency, wdlInput.AmountToWithdraw, wdlInput.Currency, wdlInput.WithdrawalAddress, wdlInput.WithdrawalNetwork)}, op.Messages()...)
	if wdlInput.Multiparty == 1 && wdlInput.Commit == 1 {
		return createWithdrawalApprovalRequest(signerUser, wallet, wdlInput, serviceFee, gc)
	}
	return nil
}

// submitWithdrawal records the withdrawal request of a signed operation and
// submits it.
func submitWithdrawal(ctx context.Context, signerUser *userModels.User, wallet *userModels.UserWallet, wdlInput *userModels.WithdrawalRequestInput, gc *sharedconfig.GlobalConfig) error {
	rec, p, err := LoadWalletOperation(wdlInput.Transaction, wallet.ID, OperationCryptoWithdrawal, gc)
	if err != nil {
		return err
	}
	var c withdrawalContext
	if rec.Context == nil || json.Unmarshal([]byte(*rec.Context), &c) != nil {
		return &tErrors.ErrorTemporaryServerError{}
	}
	request := c.Request
	request.ID = uuid.NewString()

	dbTX := gc.DB.Begin()
	defer dbTX.Rollback()
	if err := dbTX.Omit(clause.Associations).Create(&request).Error; err != nil {
		log.Printf("[QueueWithdrawalRequest] saving withdrawal request for %v: %v", wallet.ID, err)
		return &tErrors.ErrorTemporaryServerError{}
	}
	hash, err := SignSingleOwnerOperation(ctx, rec, p, signerUser.PrimarySigner, wdlInput.TransactionSignature, gc)
	if err != nil {
		return err
	}
	if err := dbTX.Model(&request).Update("transaction_id", hash).Error; err != nil {
		log.Printf("[QueueWithdrawalRequest] saving transaction %v of withdrawal %v: %v", hash, request.ID, err)
	}
	if err := dbTX.Commit().Error; err != nil {
		log.Printf("[QueueWithdrawalRequest] operation %v submitted but its withdrawal request was not saved: %v", hash, err)
		gc.LogDiscordFailedRequest(fmt.Sprintf("[QueueWithdrawalRequest] operation %v for %v (%v %v to %v) submitted but its withdrawal request was not saved: %v", hash, wallet.ID, request.AmountSubmitted, request.Currency, request.WithdrawalAddress, err))
		return &tErrors.ErrorTemporaryServerError{}
	}
	wdlInput.TransactionID = hash
	wallet.InvalidateUserCache(gc)
	signerUser.InvalidateUserWalletCache(gc)
	return nil
}

// createWithdrawalApprovalRequest records a shared wallet's withdrawal for
// its approvers.
func createWithdrawalApprovalRequest(signerUser *userModels.User, wallet *userModels.UserWallet, wdlInput *userModels.WithdrawalRequestInput, serviceFee decimal.Decimal, gc *sharedconfig.GlobalConfig) error {
	wdlInput.TransactionID = "PENDING_AUTH"
	description := fmt.Sprintf("Withdraw %v (%v),\n Amount: %v,\n Withdrawal Address: %v,\n Service Fee: %v,\n Network Fee: %v %v", wdlInput.Currency, wdlInput.WithdrawalNetwork, wdlInput.AmountSubmitted, wdlInput.WithdrawalAddress, serviceFee.String()+"%", wdlInput.WithdrawalNetworkFee, wdlInput.Currency)
	infoBytes, _ := json.Marshal(*wdlInput)
	info := string(infoBytes)
	pendingAuth := userModels.PendingAuth{
		ID:                     uuid.NewString(),
		Initiator:              signerUser.Username,
		InitiatorSignerAddress: signerUser.PrimarySigner,
		WalletAddress:          wallet.ID,
		TransactionType:        "CRYPTO WITHDRAWAL",
		Description:            description,
		TransactionSource:      wallet.ID,
		ApprovalsNeeded:        wallet.NumberOfApprovalsNeeded,
		TransactionXdr:         wdlInput.Transaction,
		TransactionInfoStr:     &info,
	}
	if err := gc.DB.Omit(clause.Associations).Create(&pendingAuth).Error; err != nil {
		log.Printf("[QueueWithdrawalRequest] saving approval request for %v: %v", wallet.ID, err)
		return &tErrors.ErrorTemporaryServerError{}
	}
	wdlInput.ReturnedDescription = description
	return nil
}
