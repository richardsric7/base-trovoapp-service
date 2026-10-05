package swaps

import (
	"fmt"
	"net/http"
	"strings"

	swapErrors "trovo-wallet-api/internal/components/swaps/errors"
	swapModels "trovo-wallet-api/internal/components/swaps/models"
	userModels "trovo-wallet-api/internal/components/users/models"
	userServices "trovo-wallet-api/internal/components/users/services"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/shopspring/decimal"
)

// SwapSend function swaps an asset to another asset
func SwapSend(signerUser, walletOwner *userModels.User, wallet *userModels.UserWallet, swapInfo *swapModels.SwapSendInfo, gc *sharedconfig.GlobalConfig) error {
	swapInfo.Messages = make([]string, 0)
	if _, err := decimal.NewFromString(swapInfo.SourceAmount); err != nil {
		return &swapErrors.ErrorInvalidSwapAmount{}
	}
	if wallet.SharedAccessEnabled == 1 && wallet.NumberOfApprovalsNeeded > 0 {
		swapInfo.Multiparty = 1
	}

	var feePercent float64
	if walletOwner.BelongsToAnEnterpriseProfile() {
		slf, exists, e := gc.GetServiceLinkFees(*walletOwner.CreatedByServiceLinkID)
		if e != nil {
			//error occured
			return e
		}
		if !exists {
			//no service fee is configured, use standard fee
			feePercent = wallet.GetSwapFee(gc).FeePercent
		} else {
			//get the enterprise config fee
			feePercent = float64(slf.SwapFee)
		}
	} else {

		feePercent = wallet.GetSwapFee(gc).FeePercent
	}
	fee := decimal.NewFromFloat(feePercent)

	feeAmount := ((fee.Mul(decimal.RequireFromString(swapInfo.SourceAmount))).Div(decimal.NewFromInt(100))).Truncate(7)
	//calculate VAT on the fee amount.
	vatFee := gc.GetVATValue(feeAmount)
	vatRate := decimal.NewFromFloat(gc.GetVATRate()).String()
	swapInfo.Vat = vatRate
	swapInfo.VatAmount = decimal.NewFromFloat(vatFee).String()
	swapAmount := decimal.RequireFromString(swapInfo.SourceAmount).Sub(feeAmount.Add(decimal.NewFromFloat(vatFee)))
	swapInfo.SwapAmount = swapAmount.String()
	swapInfo.Fee = fee.String()
	swapInfo.FeeAmount = feeAmount.String()

	if wallet.HasViewOnlyAccess(gc) {
		swapInfo.SignatureRequired = 1
	}
	//transform codes and issuer
	swapInfo.DestinationAssetCode = strings.ToUpper(swapInfo.DestinationAssetCode)
	swapInfo.DestinationContractAddress = strings.ToUpper(swapInfo.DestinationContractAddress)
	swapInfo.SourceAssetCode = strings.ToUpper(swapInfo.SourceAssetCode)
	swapInfo.SourceContractAddress = strings.ToUpper(swapInfo.SourceContractAddress)
	if userModels.IsInternalBalanceAsset(swapInfo.SourceAssetCode, swapInfo.SourceContractAddress, gc) || userModels.IsInternalBalanceAsset(swapInfo.DestinationAssetCode, swapInfo.DestinationContractAddress, gc) {
		return &tErrors.CustomError{
			Param:      "assetCode",
			Err:        "error-asset-not-sendable",
			ErrMessage: "This asset cannot be swapped directly.",
		}
	}
	if e := ValidateSwapSendInfo(swapInfo); e != nil {
		return e
	}
	if gc.IsValidTokenizedAsset(swapInfo.DestinationAssetCode) {

		t := gc.GetTokenizedAssetByCode(swapInfo.DestinationAssetCode)
		if t.AssetTokenizationStatus < 5 {
			return &tErrors.CustomError{
				Param:      "contractAddress",
				Err:        "error-asset-not-yet-available-for-sale",
				ErrMessage: "This tokenized Asset is not yet available for sale. Swap is not allowed at this time.",
				Code:       http.StatusForbidden,
			}
		}

		//check if user has done KYC
		if walletOwner.KYCVerified == 0 {
			return &tErrors.CustomError{
				Param:      "destinationAssetCode",
				Err:        "error-no-kyc",
				ErrMessage: fmt.Sprintf("%v does not meet KYC requirement to receive the asset %v", walletOwner.Username, swapInfo.DestinationAssetCode),
			}
		}

		//Check if it is still in primary sales
		if gc.IsTokenizedAssetInPrimarySales(swapInfo.DestinationAssetCode) {
			return &tErrors.CustomError{
				Param:      "destinationAssetCode",
				Err:        "error-primary-sales-active",
				ErrMessage: fmt.Sprintf("%v is still in primary sales. Please go to the tokenized asset market place to purchase from there.", swapInfo.DestinationAssetCode),
			}
		}

	}
	// the swap itself: a wallet operation filling offers on the offer book
	return userServices.SwapSend(signerUser, walletOwner, wallet, swapInfo, gc)
}
