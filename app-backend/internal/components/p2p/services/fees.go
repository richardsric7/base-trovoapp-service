package p2p

import (
	"time"
	p2pModels "trovo-wallet-api/internal/components/p2p/models"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// decimalPlaces matches CuratedAsset.DecimalPlaces' default (7) - every P2P
// amount/fee field is truncated to this precision, following the same
// convention already used by payments/swaps.
const decimalPlaces = 7

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
// specified asset amount, the fiat price, and the active fee configuration.
func CalculateOrderFees(specifiedAssetAmount decimal.Decimal, price decimal.Decimal, cfg p2pModels.TradeFeeConfiguration) OrderFeeBreakdown {
	buyerPlatformPct := decimal.RequireFromString(orDefault(cfg.BuyerPlatformFeePercent))
	buyerRegulatoryPct := decimal.RequireFromString(orDefault(cfg.BuyerRegulatoryFeePercent))
	sellerPlatformPct := decimal.RequireFromString(orDefault(cfg.SellerPlatformFeePercent))
	sellerRegulatoryPct := decimal.RequireFromString(orDefault(cfg.SellerRegulatoryFeePercent))
	vatPct := decimal.RequireFromString(orDefault(cfg.VatPercent))

	hundred := decimal.NewFromInt(100)

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

	paymentAmount := specifiedAssetAmount.Mul(price).Truncate(decimalPlaces)

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
