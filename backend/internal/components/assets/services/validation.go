package assets

import (
	errors "trovo-wallet-api/internal/errors"
)

// ValidateAssetCode validates user registration info
func ValidateAssetCode(assetCode string) error {

	return &errors.ErrorInvalidAssetCode{}

}
