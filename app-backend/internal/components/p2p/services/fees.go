package p2p

import (
	"time"
	p2pModels "trovo-wallet-api/internal/components/p2p/models"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// fiatDecimalPlaces is the truncation precision for PaymentAmount, which is
// denominated in the offer's fiat currency (Currency), not the traded
// asset - 2 matches standard fiat currency precision (cents/kobo/etc.) and
// is independent of the traded asset's own on-chain decimals (see
// CalculateOrderFees' assetDecimals parameter, which every other,
// asset-denominated field is truncated to instead).
const fiatDecimalPlaces = 2

// GetActiveFeeConfiguration returns the current ACTIVE TradeFeeConfiguration
// for a country, falling back to the "" (default/global) country code row
// if no country-specific one is configured.
func GetActiveFeeConfiguration(db *gorm.DB, countryCode string) (p2pModels.TradeFeeConfiguration, error) {
	var cfg p2pModels.TradeFeeConfiguration
	err := db.Where("country_code = ? AND status = ?", countryCode, p2pModels.FeeConfigStatusActive).
		Order("version desc").First(&cfg).Error
	if err == gorm.ErrRecordNotFound {
		err = db.Where("country_code = ? AND status = ?", "", p2pModels.FeeConfigStatusActive).
			Order("version desc").First(&cfg).Error
	}
	return cfg, err
}

// GetActiveFeeWalletConfiguration returns the current ACTIVE
// FeeCollectionWalletConfiguration for a country, with the same
// default-country fallback as GetActiveFeeConfiguration.
func GetActiveFeeWalletConfiguration(db *gorm.DB, countryCode string) (p2pModels.FeeCollectionWalletConfiguration, error) {
	var cfg p2pModels.FeeCollectionWalletConfiguration
	err := db.Where("country_code = ? AND status = ?", countryCode, p2pModels.FeeConfigStatusActive).
		Order("version desc").First(&cfg).Error
	if err == gorm.ErrRecordNotFound {
		err = db.Where("country_code = ? AND status = ?", "", p2pModels.FeeConfigStatusActive).
			Order("version desc").First(&cfg).Error
	}
	return cfg, err
}

// OrderFeeBreakdown is the fully-calculated fee/settlement math for one
// order, computed at order-creation time and snapshotted onto the Order
// (Plan Sections 46-51).
type OrderFeeBreakdown struct {
	BuyerPlatformFee      decimal.Decimal
	BuyerRegulatoryFee    decimal.Decimal
	BuyerPlatformFeeVat   decimal.Decimal
	BuyerRegulatoryFeeVat decimal.Decimal
	BuyerTotalFees        decimal.Decimal
	BuyerTotalVat         decimal.Decimal
	BuyerTotalCharges     decimal.Decimal
	BuyerNetAssetAmount   decimal.Decimal

	SellerPlatformFee       decimal.Decimal
	SellerRegulatoryFee     decimal.Decimal
	SellerPlatformFeeVat    decimal.Decimal
	SellerRegulatoryFeeVat  decimal.Decimal
	SellerTotalFees         decimal.Decimal
	SellerTotalVat          decimal.Decimal
	SellerTotalCharges      decimal.Decimal
	SellerEscrowAssetAmount decimal.Decimal

	CombinedPlatformFee   decimal.Decimal
	CombinedRegulatoryFee decimal.Decimal
	CombinedVat           decimal.Decimal

	PaymentAmount decimal.Decimal
}

// CalculateOrderFees implements the formulas in Plan Section 50, given the
// specified asset amount, the fiat price, the active fee configuration, and
// assetDecimals - the traded asset's own on-chain decimal precision (see
// network.AssetDecimals), which every asset-denominated field below is
// truncated to. This replaces a flat 7-decimal-place assumption (the old
// Stellar stroop precision) that didn't match every B20 token's real
// decimals (USDC/USDT use 6, WBTC uses 8, WETH uses 18) - truncating a real
// 18-decimal asset's fees to only 7 places was quietly rounding away real
// fee revenue on every order, and a 6-decimal asset's fees carried a
// meaningless 8th decimal digit until the wei conversion at settlement
// discarded it anyway.
func CalculateOrderFees(specifiedAssetAmount decimal.Decimal, price decimal.Decimal, cfg p2pModels.TradeFeeConfiguration, assetDecimals uint8) OrderFeeBreakdown {
	buyerPlatformPct := decimal.RequireFromString(orDefault(cfg.BuyerPlatformFeePercent))
	buyerRegulatoryPct := decimal.RequireFromString(orDefault(cfg.BuyerRegulatoryFeePercent))
	sellerPlatformPct := decimal.RequireFromString(orDefault(cfg.SellerPlatformFeePercent))
	sellerRegulatoryPct := decimal.RequireFromString(orDefault(cfg.SellerRegulatoryFeePercent))
	vatPct := decimal.RequireFromString(orDefault(cfg.VatPercent))

	hundred := decimal.NewFromInt(100)
	decimalPlaces := int32(assetDecimals)

	buyerPlatformFee := specifiedAssetAmount.Mul(buyerPlatformPct).Div(hundred).Truncate(decimalPlaces)
	buyerRegulatoryFee := specifiedAssetAmount.Mul(buyerRegulatoryPct).Div(hundred).Truncate(decimalPlaces)
	buyerPlatformFeeVat := buyerPlatformFee.Mul(vatPct).Div(hundred).Truncate(decimalPlaces)
	buyerRegulatoryFeeVat := buyerRegulatoryFee.Mul(vatPct).Div(hundred).Truncate(decimalPlaces)
	buyerTotalFees := buyerPlatformFee.Add(buyerRegulatoryFee)
	buyerTotalVat := buyerPlatformFeeVat.Add(buyerRegulatoryFeeVat)
	buyerTotalCharges := buyerTotalFees.Add(buyerTotalVat)
	buyerNetAssetAmount := specifiedAssetAmount.Sub(buyerTotalCharges)

	sellerPlatformFee := specifiedAssetAmount.Mul(sellerPlatformPct).Div(hundred).Truncate(decimalPlaces)
	sellerRegulatoryFee := specifiedAssetAmount.Mul(sellerRegulatoryPct).Div(hundred).Truncate(decimalPlaces)
	sellerPlatformFeeVat := sellerPlatformFee.Mul(vatPct).Div(hundred).Truncate(decimalPlaces)
	sellerRegulatoryFeeVat := sellerRegulatoryFee.Mul(vatPct).Div(hundred).Truncate(decimalPlaces)
	sellerTotalFees := sellerPlatformFee.Add(sellerRegulatoryFee)
	sellerTotalVat := sellerPlatformFeeVat.Add(sellerRegulatoryFeeVat)
	sellerTotalCharges := sellerTotalFees.Add(sellerTotalVat)
	sellerEscrowAssetAmount := specifiedAssetAmount.Add(sellerTotalCharges)

	// PaymentAmount is fiat (Currency), not the traded asset - it gets its
	// own, currency-appropriate precision rather than the asset's.
	paymentAmount := specifiedAssetAmount.Mul(price).Truncate(fiatDecimalPlaces)

	return OrderFeeBreakdown{
		BuyerPlatformFee:      buyerPlatformFee,
		BuyerRegulatoryFee:    buyerRegulatoryFee,
		BuyerPlatformFeeVat:   buyerPlatformFeeVat,
		BuyerRegulatoryFeeVat: buyerRegulatoryFeeVat,
		BuyerTotalFees:        buyerTotalFees,
		BuyerTotalVat:         buyerTotalVat,
		BuyerTotalCharges:     buyerTotalCharges,
		BuyerNetAssetAmount:   buyerNetAssetAmount,

		SellerPlatformFee:       sellerPlatformFee,
		SellerRegulatoryFee:     sellerRegulatoryFee,
		SellerPlatformFeeVat:    sellerPlatformFeeVat,
		SellerRegulatoryFeeVat:  sellerRegulatoryFeeVat,
		SellerTotalFees:         sellerTotalFees,
		SellerTotalVat:          sellerTotalVat,
		SellerTotalCharges:      sellerTotalCharges,
		SellerEscrowAssetAmount: sellerEscrowAssetAmount,

		CombinedPlatformFee:   buyerPlatformFee.Add(sellerPlatformFee),
		CombinedRegulatoryFee: buyerRegulatoryFee.Add(sellerRegulatoryFee),
		CombinedVat:           buyerTotalVat.Add(sellerTotalVat),

		PaymentAmount: paymentAmount,
	}
}

func orDefault(s string) string {
	if s == "" {
		return "0"
	}
	return s
}

// nowUTC is a small seam kept for testability.
func nowUTC() time.Time {
	return time.Now().UTC()
}
