package swaps

import (
	swapErrors "trovo-wallet-api/internal/components/swaps/errors"
	swapModels "trovo-wallet-api/internal/components/swaps/models"
	"trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/validators"

	"github.com/shopspring/decimal"
)

//ValidateSwapSendPathInfo validates payment information
func ValidateSwapSendInfo(swapInfo *swapModels.SwapSendInfo) error {

	//required parameters
	{
		if len(swapInfo.SourceAmount) == 0 {
			return &errors.ErrorMissingParameter{Parameter: "sourceAmount"}
		}

		if len(swapInfo.DestinationAssetCode) == 0 && len(swapInfo.DestinationAssetIssuer) > 0 {
			return &errors.ErrorMissingParameter{Parameter: "destinationAssetCode"}
		}
		if len(swapInfo.DestinationAssetIssuer) == 0 && len(swapInfo.DestinationAssetCode) > 0 {
			return &errors.ErrorMissingParameter{Parameter: "destinationAssetIssuer"}
		}

	}

	//vaidate asset issuer public key

	{
		if len(swapInfo.SourceAssetIssuer) > 0 {
			err := validators.ValidatePublicKeyFormat(swapInfo.SourceAssetIssuer)
			if err != nil {
				return err
			}
		}
		if len(swapInfo.DestinationAssetIssuer) > 0 {
			err := validators.ValidatePublicKeyFormat(swapInfo.DestinationAssetIssuer)
			if err != nil {
				return err
			}
		}
	}

	//Validate asset code

	{
		if len(swapInfo.DestinationAssetCode) > 0 {
			err := validators.ValidateAssetCodeFormat(swapInfo.DestinationAssetCode)

			if err != nil {
				return err
			}
		}
		if len(swapInfo.SourceAssetCode) > 0 {
			err := validators.ValidateAssetCodeFormat(swapInfo.SourceAssetCode)

			if err != nil {
				return err
			}
		}

	}

	//check if asset and destination are same
	if swapInfo.DestinationAssetCode+swapInfo.DestinationAssetIssuer == swapInfo.SourceAssetCode+swapInfo.SourceAssetIssuer {
		return &swapErrors.ErrorSourceAndDestinationAssetAreSame{}
	}

	//check if amount is valid.

	{
		var err error
		var amount decimal.Decimal
		if amount, err = decimal.NewFromString(swapInfo.SourceAmount); err != nil {
			return &swapErrors.ErrorInvalidSwapAmount{}
		}

		if amount.LessThan(decimal.NewFromFloat(0.0000001)) {
			return &swapErrors.ErrorInvalidSwapAmount{}
		}

	}

	//signature must exist if transaction exists

	{
		if len(swapInfo.Transaction) > 0 {
			if len(swapInfo.TransactionSignature) == 0 {
				return &errors.ErrorMissingParameter{Parameter: "transactionSignature"}
			}
		}
	}

	return nil

}
