package users

import (
	"net/http"
	"strings"

	assetModels "trovo-wallet-api/internal/components/assets/models"
	userModels "trovo-wallet-api/internal/components/users/models"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/sharedconfig"
)

// GasFeeAssets lists the curated assets users can pay network fees in.
func GasFeeAssets(gc *sharedconfig.GlobalConfig) []assetModels.CuratedAsset {
	assets := make([]assetModels.CuratedAsset, 0)
	gc.DB.Where("gas_fee_eligible = ? AND inactive = ?", true, 0).Order("asset_code").Find(&assets)
	return assets
}

// SetGasFeeAsset sets (or, with an empty code, clears) the stablecoin the
// user's wallets pay network fees in. Without one, or when a wallet holds
// too little of it, fees are paid in ETH.
func SetGasFeeAsset(user *userModels.User, assetCode string, gc *sharedconfig.GlobalConfig) error {
	code := strings.ToUpper(strings.TrimSpace(assetCode))
	var value *string
	if code != "" {
		var n int64
		gc.DB.Model(&assetModels.CuratedAsset{}).Where("asset_code = ? AND gas_fee_eligible = ? AND inactive = ?", code, true, 0).Count(&n)
		if n == 0 {
			return &tErrors.CustomError{Param: "assetCode", Err: "error-gas-fee-asset-not-eligible", ErrMessage: code + " cannot be used to pay network fees.", Code: http.StatusBadRequest}
		}
		value = &code
	}
	if err := gc.DB.Model(&userModels.User{}).Where("id = ?", user.ID).Update("gas_fee_asset", value).Error; err != nil {
		return &tErrors.ErrorTemporaryServerError{}
	}
	user.GasFeeAsset = value
	user.InvalidateUserCache(gc)
	return nil
}
