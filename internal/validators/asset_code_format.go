package validators

import (
	"regexp"
	"trovo-wallet-api/internal/errors"
)

//ValidateAssetCodeFormat validates a blockchain asset code

func ValidateAssetCodeFormat(assetCode string) error {

	_len := len(assetCode)

	if _len == 0 {
		return nil
	}

	if _len > 12 {
		return &errors.ErrorInvalidAssetCode{}
	}

	matched, _ := regexp.MatchString("(?i)^[a-z0-9][a-z0-9]+$", assetCode)

	if !matched {
		return &errors.ErrorInvalidAssetCode{}
	}

	return nil
}
