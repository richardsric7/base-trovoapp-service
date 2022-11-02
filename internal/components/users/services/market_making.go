package users

import (
	"log"
	"strings"
	userModels "trovo-wallet-api/internal/components/users/models"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/shopspring/decimal"
	"github.com/stellar/go/txnbuild"
	"github.com/stellar/go/xdr"
)

func MakeOffer(signerUser *userModels.User, sourceWallet *userModels.UserWallet, offerRequest *userModels.MakeOfferRequest, gc *sharedconfig.GlobalConfig) (err error) {
	var currencyAsset, mainAsset txnbuild.Asset
	var offerBCRequest txnbuild.ManageSellOffer
	if offerRequest.AssetCode+offerRequest.AssetIssuer == offerRequest.CurrencyCode+offerRequest.CurrencyIssuer {
		return &tErrors.CustomError{
			Param:      "AssetCode",
			Err:        "error-asset-and-currency-are-the-same",
			ErrMessage: "You cannot make a market against same asset",
		}
	}
	mmWallet, err := userModels.UserWalletID(offerRequest.MarketMakingWalletPK).GetWallet(gc.DB)
	if err != nil {
		log.Printf("[MakeOffer]Error validating Market making wallet: %v", err)
		return err
	}

	if mmWallet.Signer != sourceWallet.Signer {
		return &tErrors.CustomError{
			Param:      "marketMakingWalletPK",
			Err:        "error-unauthorized-access",
			ErrMessage: "Market Making wallet does ",
		}
	}

	fraction := decimal.RequireFromString(offerRequest.PricePerAsset).Rat()
	d := int32(fraction.Denom().Int64())
	n := int32(fraction.Num().Int64())
	if offerRequest.AssetIssuer == "" {
		mainAsset = txnbuild.NativeAsset{}
	} else {
		mainAsset = txnbuild.CreditAsset{Code: offerRequest.AssetCode, Issuer: offerRequest.AssetIssuer}
	}

	if offerRequest.CurrencyIssuer == "" {
		currencyAsset = txnbuild.NativeAsset{}
	} else {
		currencyAsset = txnbuild.CreditAsset{Code: offerRequest.CurrencyCode, Issuer: offerRequest.CurrencyIssuer}
	}

	if strings.EqualFold(offerRequest.OfferType, "Buy") {
		//invert the price fraction
		offerBCRequest = txnbuild.ManageSellOffer{
			Selling:       mainAsset,
			Buying:        currencyAsset,
			Amount:        offerRequest.Quantity,
			Price:         xdr.Price{D: xdr.Int32(n), N: xdr.Int32(d)},
			SourceAccount: sourceWallet.ID,
		}
	}
	if strings.EqualFold(offerRequest.OfferType, "Sell") {
		offerBCRequest = txnbuild.ManageSellOffer{
			Selling:       mainAsset,
			Buying:        currencyAsset,
			Amount:        offerRequest.Quantity,
			Price:         xdr.Price{N: xdr.Int32(n), D: xdr.Int32(d)},
			SourceAccount: sourceWallet.ID,
		}

	}

}
