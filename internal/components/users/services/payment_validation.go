package users

import (
	"os"
	tPayErrors "trovo-wallet-api/internal/components/payments/errors"
	paymentModels "trovo-wallet-api/internal/components/payments/models"
	"trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/validators"

	"github.com/shopspring/decimal"
)

// ValidatePaymentInfo validates payment information
func ValidatePaymentInfo(paymentInfo *paymentModels.PaymentInfo) (*paymentModels.PaymentInfo, error) {

	//required parameters
	{
		if len(paymentInfo.Amount) == 0 {
			return paymentInfo, &errors.ErrorMissingParameter{Parameter: "amount"}
		}
		minAmountSendable := "0.0000100"
		if os.Getenv("MIN_SENDABLE_AMOUNT") != "" && os.Getenv("MIN_SENDABLE_AMOUNT") != "0" {
			minAmountSendable = os.Getenv("MIN_SENDABLE_AMOUNT")
		}

		if decimal.RequireFromString(paymentInfo.Amount).LessThan(decimal.RequireFromString(minAmountSendable)) {
			return paymentInfo, &tPayErrors.ErrorPaymentAmountBelowMinAllowed{}
		}

		if len(paymentInfo.Destination) == 0 {
			return paymentInfo, &errors.ErrorMissingParameter{Parameter: "destination"}
		}

	}

	//vaidate asset issuer public key

	{
		if len(paymentInfo.AssetIssuer) > 0 {
			err := validators.ValidatePublicKeyFormat(paymentInfo.AssetIssuer)
			if err != nil {
				return paymentInfo, err
			}
		}
	}

	//vaidate channelAccount public key

	{
		if len(paymentInfo.ChannelAccount) > 0 {
			err := validators.ValidatePublicKeyFormat(paymentInfo.ChannelAccount)
			if err != nil {
				return paymentInfo, err
			}
		}
	}

	//Validate asset code

	{
		err := validators.ValidateAssetCodeFormat(paymentInfo.AssetCode)

		if err != nil {
			return paymentInfo, err
		}
	}

	//signature must exist if transaction exists

	{
		if len(paymentInfo.Transaction) > 0 {
			if len(paymentInfo.TransactionSignature) == 0 {
				return paymentInfo, &errors.ErrorMissingParameter{Parameter: "transactionSignature"}
			}
		}
	}

	//check memo

	{
		if len(paymentInfo.Memo) > 28 {
			return paymentInfo, &tPayErrors.ErrorInvalidPaymentMemo{}
		}
	}

	return paymentInfo, nil

}
