package swaps

import (
	swapErrors "trovo-wallet-api/internal/components/swaps/errors"
	swapModels "trovo-wallet-api/internal/components/swaps/models"
	"trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/validators"

	"github.com/shopspring/decimal"
)

// ValidateSwapSendPathInfo validates payment information
func ValidateSwapSendInfo(swapInfo *swapModels.SwapSendInfo) error {

	//required parameters
	{
		if len(swapInfo.SourceAmount) == 0 {
			return &errors.ErrorMissingParameter{Parameter: "sourceAmount"}
		}

		if len(swapInfo.DestinationAssetCode) == 0 && len(swapInfo.DestinationContractAddress) > 0 {
			return &errors.ErrorMissingParameter{Parameter: "destinationAssetCode"}
		}
		if len(swapInfo.DestinationContractAddress) == 0 && len(swapInfo.DestinationAssetCode) > 0 {
			return &errors.ErrorMissingParameter{Parameter: "destinationContractAddress"}
		}

	}

	//vaidate asset issuer public key

	{
		if len(swapInfo.SourceContractAddress) > 0 {
			err := validators.ValidateAddressFormat(swapInfo.SourceContractAddress)
			if err != nil {
				return err
			}
		}
		if len(swapInfo.DestinationContractAddress) > 0 {
			err := validators.ValidateAddressFormat(swapInfo.DestinationContractAddress)
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
	if swapInfo.DestinationAssetCode+swapInfo.DestinationContractAddress == swapInfo.SourceAssetCode+swapInfo.SourceContractAddress {
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
			if len(swapInfo.TransactionSignature) == 0 && swapInfo.Commit == 0 {
				return &errors.ErrorMissingParameter{Parameter: "transactionSignature"}
			}
		}
	}

	return nil

}

// ValidateSwapSendPathInfo validates swap information
func ValidateSwapReceiveInfo(swapInfo *swapModels.SwapReceiveInfo) error {

	//required parameters
	{
		if len(swapInfo.DestinationAmount) == 0 {
			return &errors.ErrorMissingParameter{Parameter: "destinationAmount"}
		}

		if len(swapInfo.DestinationAssetCode) == 0 && len(swapInfo.DestinationContractAddress) > 0 {
			return &errors.ErrorMissingParameter{Parameter: "destinationAssetCode"}
		}
		if len(swapInfo.DestinationContractAddress) == 0 && len(swapInfo.DestinationAssetCode) > 0 {
			return &errors.ErrorMissingParameter{Parameter: "destinationContractAddress"}
		}
		if len(swapInfo.SourceAssetCode) == 0 && len(swapInfo.SourceContractAddress) > 0 {
			return &errors.ErrorMissingParameter{Parameter: "sourceAssetCode"}
		}
		if len(swapInfo.SourceContractAddress) == 0 && len(swapInfo.SourceAssetCode) > 0 {
			return &errors.ErrorMissingParameter{Parameter: "sourceContractAddress"}
		}

	}

	//vaidate asset issuer public key

	{
		if len(swapInfo.SourceContractAddress) > 0 {
			err := validators.ValidateAddressFormat(swapInfo.SourceContractAddress)
			if err != nil {
				return err
			}
		}
		if len(swapInfo.DestinationContractAddress) > 0 {
			err := validators.ValidateAddressFormat(swapInfo.DestinationContractAddress)
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
	if swapInfo.DestinationAssetCode+swapInfo.DestinationContractAddress == swapInfo.SourceAssetCode+swapInfo.SourceContractAddress {
		return &swapErrors.ErrorSourceAndDestinationAssetAreSame{}
	}

	//check if amount is valid.

	{
		var err error
		var amount decimal.Decimal
		if amount, err = decimal.NewFromString(swapInfo.DestinationAmount); err != nil {
			return &swapErrors.ErrorInvalidSwapAmount{}
		}

		if amount.LessThan(decimal.NewFromFloat(0.0000001)) {
			return &swapErrors.ErrorInvalidSwapAmount{}
		}

	}

	//signature must exist if transaction exists

	{
		if len(swapInfo.Transaction) > 0 {
			if len(swapInfo.TransactionSignature) == 0 && swapInfo.Commit == 0 {
				return &errors.ErrorMissingParameter{Parameter: "transactionSignature"}
			}
		}
	}

	return nil

}
