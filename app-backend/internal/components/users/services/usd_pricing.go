package users

import (
	"net/http"
	"os"
	"strings"

	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/shopspring/decimal"
)

// usdAmountIn converts a USD-denominated fee or price into units of the
// asset it is paid in. Without a DEX only stablecoins can be priced: the
// dollar asset (DOLLAR_ASSET) and USD-named codes (USDC, USDT, USDB, ...)
// at 1:1, and the naira asset (NAIRA_ASSET) at the USD/cNGN rate.
func usdAmountIn(usd float64, assetCode, contract string, gc *sharedconfig.GlobalConfig) (decimal.Decimal, error) {
	code := strings.ToUpper(strings.TrimSpace(assetCode))
	id := code + ":" + strings.TrimSpace(contract)
	if d := strings.TrimSpace(os.Getenv("DOLLAR_ASSET")); d != "" && strings.EqualFold(id, d) {
		return decimal.NewFromFloat(usd), nil
	}
	if code != "" && (strings.HasPrefix(code, "USD") || strings.HasSuffix(code, "USD")) {
		return decimal.NewFromFloat(usd), nil
	}
	if n := strings.TrimSpace(os.Getenv("NAIRA_ASSET")); n != "" && strings.EqualFold(id, n) {
		a := gc.ConvertUsdToCngn(usd)
		if !a.IsPositive() {
			return decimal.Zero, &tErrors.ErrorTemporaryServerError{}
		}
		return a, nil
	}
	return decimal.Zero, &tErrors.CustomError{Param: "assetCode", Err: "error-asset-not-priceable", ErrMessage: "This fee can currently be paid only in a dollar or naira stablecoin.", Code: http.StatusBadRequest}
}

// feeExemptProfile reports whether username pays no platform service fees
// (see GlobalConfig.FeeExemptUsername).
func feeExemptProfile(username string, gc *sharedconfig.GlobalConfig) bool {
	return gc.FeeExemptUsername(username)
}
