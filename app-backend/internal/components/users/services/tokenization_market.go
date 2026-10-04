package users

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"strings"
	"time"

	"trovo-wallet-api/internal/aa"
	"trovo-wallet-api/internal/basetxn"
	userModels "trovo-wallet-api/internal/components/users/models"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/evmkeypair"
	"trovo-wallet-api/internal/network"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
)

// Tokenized assets are sold through TrovoOfferBook (market/contracts): the
// mint operation escrows the supply for sale from the asset's distribution
// Safe in an offer priced in each accepted payment token, and every
// purchase is a fill of that offer, authorized by the platform
// (OFFER_AUTHORIZER_PRIVATE_KEY) once the buyer passes the platform's
// checks (KYC, purchase cap, sale window).

// offerBookAddress is OFFER_BOOK_ADDRESS.
func offerBookAddress() (common.Address, error) {
	raw := strings.TrimSpace(os.Getenv("OFFER_BOOK_ADDRESS"))
	if !common.IsHexAddress(raw) {
		return common.Address{}, &tErrors.CustomError{Param: "tokenizedAssetId", Err: "error-offer-book-not-configured", ErrMessage: "Tokenized asset sales are not configured. Please try again later.", Code: http.StatusServiceUnavailable}
	}
	return common.HexToAddress(raw), nil
}

// offerAuthorizerKey is OFFER_AUTHORIZER_PRIVATE_KEY: the key that signs
// fill authorizations. It can only let someone buy at the seller's price;
// it never controls anyone's tokens.
func offerAuthorizerKey() (*ecdsa.PrivateKey, error) {
	kp, err := evmkeypair.ParseFull(strings.TrimSpace(os.Getenv("OFFER_AUTHORIZER_PRIVATE_KEY")))
	if err != nil {
		log.Printf("[offerAuthorizerKey] OFFER_AUTHORIZER_PRIVATE_KEY: %v", err)
		return nil, &tErrors.CustomError{Param: "tokenizedAssetId", Err: "error-offer-book-not-configured", ErrMessage: "Tokenized asset sales are not configured. Please try again later.", Code: http.StatusServiceUnavailable}
	}
	return kp.PrivateKey(), nil
}

// offerPrice converts a price per whole token of the asset, in the quote
// currency, into the offer book's price for a payment token pegged 1:1 to
// that currency: payment base units per asset base unit, as an exact
// fraction.
func offerPrice(pricePerToken decimal.Decimal, assetDecimals, paymentDecimals uint8) (aa.OfferPrice, error) {
	if !pricePerToken.IsPositive() {
		return aa.OfferPrice{}, fmt.Errorf("price must be positive")
	}
	// scale away the price's own fractional digits
	k := int32(0)
	if e := pricePerToken.Exponent(); e < 0 {
		k = -e
	}
	num := pricePerToken.Shift(int32(paymentDecimals) + k).BigInt()
	den := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(assetDecimals)+int64(k)), nil)
	g := new(big.Int).GCD(nil, nil, num, den)
	p := aa.OfferPrice{Num: new(big.Int).Quo(num, g), Den: new(big.Int).Quo(den, g)}
	if !p.Valid() {
		return aa.OfferPrice{}, fmt.Errorf("price %v cannot be represented", pricePerToken)
	}
	return p, nil
}

// assetOffer is a tokenized asset's primary-sale offer as it stands.
type assetOffer struct {
	Book           common.Address
	ID             *big.Int
	Offer          *aa.Offer
	Token          common.Address
	AssetDecimals  uint8
	AssetCode      string
	TokenizedAsset *userModels.TokenizedAsset
}

func loadAssetOffer(ctx context.Context, ta *userModels.TokenizedAsset, gc *sharedconfig.GlobalConfig) (*assetOffer, error) {
	book, err := offerBookAddress()
	if err != nil {
		return nil, err
	}
	token, err := tokenizedAssetContract(ta)
	if err != nil {
		return nil, err
	}
	id, ok := new(big.Int), false
	if ta.OfferBookOfferID != nil {
		id, ok = id.SetString(*ta.OfferBookOfferID, 10)
	}
	if !ok {
		return nil, &tErrors.CustomError{Param: "tokenizedAssetId", Err: "error-not-on-sale", ErrMessage: "This asset has not been put on sale yet.", Code: http.StatusConflict}
	}
	o, err := aa.ReadOffer(ctx, gc.BantuExpansionClient, book, id)
	if err != nil {
		log.Printf("[loadAssetOffer] reading offer %v of %v: %v", id, *ta.AssetCode, err)
		return nil, &tErrors.ErrorTemporaryServerError{}
	}
	dec, err := network.AssetDecimals(ctx, gc.BantuExpansionClient, basetxn.CreditAsset{Code: *ta.AssetCode, Issuer: token})
	if err != nil {
		return nil, &tErrors.ErrorTemporaryServerError{}
	}
	return &assetOffer{Book: book, ID: id, Offer: o, Token: common.HexToAddress(token), AssetDecimals: dec, AssetCode: *ta.AssetCode, TokenizedAsset: ta}, nil
}

// assetPurchase is a priced purchase from an asset's offer.
type assetPurchase struct {
	Fill          aa.FillRequest
	PaymentCode   string
	PaymentDec    uint8
	Payment       *big.Int // payment token base units
	AssetAmount   *big.Int // asset base units
	Authorization []byte
}

// AssetAmountHuman is the asset amount in whole tokens.
func (p *assetPurchase) AssetAmountHuman(assetDecimals uint8) decimal.Decimal {
	return decimal.NewFromBigInt(p.AssetAmount, -int32(assetDecimals))
}

// PaymentHuman is the payment in whole payment tokens.
func (p *assetPurchase) PaymentHuman() decimal.Decimal {
	return decimal.NewFromBigInt(p.Payment, -int32(p.PaymentDec))
}

// priceAssetPurchase prices spending up to spend (whole units of the payment
// token) on the asset's offer, delivered to recipient, and has the platform
// authorize taker to make it, valid for validity.
func priceAssetPurchase(ctx context.Context, o *assetOffer, paymentCode string, paymentToken common.Address, spend decimal.Decimal, taker, recipient common.Address, validity time.Duration, gc *sharedconfig.GlobalConfig) (*assetPurchase, error) {
	if !o.Offer.Open || o.Offer.Remaining.Sign() == 0 {
		return nil, &tErrors.CustomError{Param: "amount", Err: "error-no-liquidity", ErrMessage: fmt.Sprintf("The allocated %v asset has sold out.", o.AssetCode)}
	}
	price, err := aa.ReadOfferPrice(ctx, gc.BantuExpansionClient, o.Book, o.ID, paymentToken)
	if err != nil {
		return nil, &tErrors.ErrorTemporaryServerError{}
	}
	if price.Num.Sign() == 0 {
		return nil, &tErrors.CustomError{Param: "paymentAssetCode", Err: "error-invalid-payment-asset", ErrMessage: fmt.Sprintf("%v cannot be bought with %v.", o.AssetCode, paymentCode)}
	}
	payDec, err := network.AssetDecimals(ctx, gc.BantuExpansionClient, basetxn.CreditAsset{Code: paymentCode, Issuer: paymentToken.Hex()})
	if err != nil {
		return nil, &tErrors.ErrorTemporaryServerError{}
	}
	budget := spend.Shift(int32(payDec)).Truncate(0).BigInt()
	amount := price.AmountFor(budget)
	if amount.Sign() == 0 {
		return nil, &tErrors.CustomError{Param: "amount", Err: "error-amount-too-small", ErrMessage: fmt.Sprintf("%v %v does not buy any %v.", spend, paymentCode, o.AssetCode)}
	}
	if amount.Cmp(o.Offer.Remaining) > 0 {
		most := decimal.NewFromBigInt(price.Cost(o.Offer.Remaining), -int32(payDec))
		return nil, &tErrors.CustomError{Param: "amount", Err: "error-low-liquidity", ErrMessage: fmt.Sprintf("Remaining %v token can be purchased with a maximum of %v %v. Please adjust your purchase amount.", o.AssetCode, most, paymentCode)}
	}
	payment := price.Cost(amount)
	nonce, err := aa.RandomFillNonce()
	if err != nil {
		return nil, &tErrors.ErrorTemporaryServerError{}
	}
	f := aa.FillRequest{
		OfferID: o.ID, PaymentToken: paymentToken, Amount: amount, MaxPayment: payment, Recipient: recipient,
		Nonce: nonce, Deadline: big.NewInt(time.Now().Add(validity).Unix()),
	}
	key, err := offerAuthorizerKey()
	if err != nil {
		return nil, err
	}
	sig, err := aa.SignFill(key, network.GetBlockchainChainID(), o.Book, f, taker)
	if err != nil {
		return nil, &tErrors.ErrorTemporaryServerError{}
	}
	return &assetPurchase{Fill: f, PaymentCode: paymentCode, PaymentDec: payDec, Payment: payment, AssetAmount: amount, Authorization: sig}, nil
}

// purchaseCalls are the buyer's calls making purchase p: allow the book to
// take the payment, then fill.
func (p *assetPurchase) purchaseCalls(book common.Address) ([]aa.Call, error) {
	fill, err := aa.FillCall(book, p.Fill, p.Authorization)
	if err != nil {
		return nil, err
	}
	return []aa.Call{aa.ERC20Approve(p.Fill.PaymentToken, book, p.Payment), fill}, nil
}
