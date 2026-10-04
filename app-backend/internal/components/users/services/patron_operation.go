package users

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"trovo-wallet-api/internal/aa"
	"trovo-wallet-api/internal/basetxn"
	paymentModels "trovo-wallet-api/internal/components/payments/models"
	userModels "trovo-wallet-api/internal/components/users/models"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/network"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// OperationPatronSubscription is the wallet operation kind for patron
// subscription payments.
const OperationPatronSubscription = "PATRON SUBSCRIPTION"

// patronContext is what a patron subscription operation saves; Fees is
// read by recordOperationFees.
type patronContext struct {
	MembershipGradeID uint64                       `json:"membershipGradeId"`
	Fees              []sharedconfig.FeeCollection `json:"fees"`
}

// patronPriceIn converts a USD price into the payment asset (see
// usdAmountIn: only dollar and naira stablecoins can be priced without a DEX).
func patronPriceIn(usd float64, assetCode, contract string, gc *sharedconfig.GlobalConfig) (decimal.Decimal, error) {
	return usdAmountIn(usd, assetCode, contract, gc)
}

// submitPatronSubscription is the payment half of a patron subscription,
// paid from the primary wallet:
//
//  1. without transactionSignature: builds the operation paying the price
//     plus VAT to the patron fee wallet and returns it as transaction;
//  2. with transaction + transactionSignature: submits it and, once
//     submitted, commits the subscription (tx).
func submitPatronSubscription(owner *userModels.User, patronSubInput *userModels.PatronSubscriptionInput, priceConfig *userModels.PatronMembershipGrade, tx *gorm.DB, gc *sharedconfig.GlobalConfig) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	primary, err := primaryWalletOf(owner)
	if err != nil {
		return err
	}

	if len(patronSubInput.TransactionSignature) > 0 {
		rec, p, err := LoadWalletOperation(patronSubInput.Transaction, primary.ID, OperationPatronSubscription, gc)
		if err != nil {
			return err
		}
		var c patronContext
		if rec.Context == nil || json.Unmarshal([]byte(*rec.Context), &c) != nil || c.MembershipGradeID != patronSubInput.PatronMembershipGradeID {
			return &tErrors.CustomError{Param: "transaction", Err: "transaction mismatch", ErrMessage: "transaction mismatch, please try again", Code: http.StatusBadRequest}
		}
		hash, err := SignSingleOwnerOperation(ctx, rec, p, owner.PrimarySigner, patronSubInput.TransactionSignature, gc)
		if err != nil {
			logDiscordFailedSubscription(fmt.Sprintf("Error submitting patron subscription [%+v] operation: %s", patronSubInput, err.Error()))
			return err
		}
		if err := tx.Commit().Error; err != nil {
			logDiscordFailedSubscription(fmt.Sprintf("[SubscribeToPatronPackage] operation %v for %v submitted but the subscription was not saved: %v", hash, owner.Username, err))
			return &tErrors.ErrorTemporaryServerError{}
		}
		recordOperationFees(rec, hash, gc)
		patronSubInput.TransactionID = hash
		owner.InvalidateUserCache(gc)
		return nil
	}

	price, err := patronPriceIn(priceConfig.Price, patronSubInput.PaymentAssetCode, patronSubInput.PaymentContractAddress, gc)
	if err != nil {
		return err
	}
	asset := basetxn.CreditAsset{Code: patronSubInput.PaymentAssetCode, Issuer: patronSubInput.PaymentContractAddress}
	decimals, err := network.AssetDecimals(ctx, gc.BantuExpansionClient, asset)
	if err != nil {
		return &tErrors.ErrorTemporaryServerError{}
	}
	price = price.Truncate(int32(decimals))
	vat := decimal.NewFromFloat(gc.GetVATValue(price))
	total := price.Add(vat)
	patronSubInput.Vat = decimal.NewFromFloat(gc.GetVATRate()).String()
	patronSubInput.VatAmount = vat.String()
	patronSubInput.AmountToPay = total.String()

	if bal, err := network.B20BalanceOf(gc.BantuExpansionClient, asset.Issuer, primary.ID, decimals); err != nil {
		return &tErrors.ErrorTemporaryServerError{}
	} else if bal.LessThan(total) {
		return &tErrors.ErrorUnderfundedAccount{Detail: fmt.Sprintf("You need to add at least %v %v to make up for the subscription fee.", total.Sub(bal), asset.Code)}
	}
	feeAddr, err := feeWalletAddress(primary.GetPatronFee(gc).FeeWalletSecretKey, "patron fee wallet", gc)
	if err != nil {
		return err
	}
	units, err := baseUnits(total.String(), decimals)
	if err != nil {
		return err
	}
	c := patronContext{MembershipGradeID: patronSubInput.PatronMembershipGradeID}
	if vat.IsPositive() {
		c.Fees = append(c.Fees, feeRecord("VAT", primary, owner, asset, vat.String(), feeAddr, 0))
	}
	op, err := PrepareWalletOperation(ctx, OperationPatronSubscription, owner, owner, primary, []aa.Call{transferCall(asset, feeAddr, units)}, 0, c, gc)
	if err != nil {
		return err
	}
	pay := paymentModels.PaymentInfo{AssetCode: asset.Code, ContractAddress: asset.Issuer, AmountToPay: total.String()}
	if err := checkPaymentLeavesGas(ctx, primary, &pay, op.Prepared, gc); err != nil {
		return err
	}
	patronSubInput.Transaction = op.Transaction
	patronSubInput.Messages = append(patronSubInput.Messages, fmt.Sprintf("%v %v will be debited from wallet %v to complete the subscription. This is inclusive of VAT (%v %v)", total, asset.Code, owner.Username, vat, asset.Code))
	patronSubInput.Messages = append(patronSubInput.Messages, op.Messages()...)
	return nil
}
