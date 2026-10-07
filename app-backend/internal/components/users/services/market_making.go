package users

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/big"
	"net/http"
	"os"
	"strings"
	"time"

	"trovo-wallet-api/internal/aa"
	"trovo-wallet-api/internal/basetxn"
	userModels "trovo-wallet-api/internal/components/users/models"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/network"
	"trovo-wallet-api/internal/offerbook"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm/clause"
)

// Market-making offers are offers on TrovoOfferBook (market/contracts) a
// wallet places with its own funds: SELL escrows the asset, priced in the
// currency; BUY escrows the currency, priced in the asset. Each is one
// wallet operation the user signs (allow the book to take the escrow, then
// createOffer); cancelling (cancelOffer) returns what is left. Swaps fill
// these offers. Only tokens trade on the book.

// Market-making operation kinds.
const (
	OperationMarketOffer       = "MAKE MARKET OFFER"
	OperationCancelMarketOffer = "DELETE MARKET OFFER"
)

// marketOfferContext is what a market-making operation needs once it is
// submitted and mined.
type marketOfferContext struct {
	Offer userModels.MarketOffer `json:"offer"`
	Book  string                 `json:"book"`
}

// inverseOfferPrice is the offer book price, in asset base units per
// currency base unit, of buying the asset at pricePerToken (currency per
// whole asset): what a BUY offer, which sells the currency, charges.
func inverseOfferPrice(pricePerToken decimal.Decimal, assetDecimals, currencyDecimals uint8) (aa.OfferPrice, error) {
	if !pricePerToken.IsPositive() {
		return aa.OfferPrice{}, fmt.Errorf("price must be positive")
	}
	k := int32(0)
	if e := pricePerToken.Exponent(); e < 0 {
		k = -e
	}
	// price = P / 10^k, so 1/price = 10^k / P
	p := pricePerToken.Shift(k).BigInt()
	num := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(k)+int64(assetDecimals)), nil)
	den := new(big.Int).Mul(p, new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(currencyDecimals)), nil))
	g := new(big.Int).GCD(nil, nil, num, den)
	out := aa.OfferPrice{Num: new(big.Int).Quo(num, g), Den: new(big.Int).Quo(den, g)}
	if !out.Valid() {
		return aa.OfferPrice{}, fmt.Errorf("price %v cannot be represented", pricePerToken)
	}
	return out, nil
}

// marketToken is a market side's token, or an error for the native asset.
func marketToken(code, contract, param string) (common.Address, error) {
	if contract == "" || strings.EqualFold(code, os.Getenv("NATIVE_ASSET_CODE")) || !common.IsHexAddress(contract) {
		return common.Address{}, &tErrors.CustomError{Param: param, Err: "error-asset-not-tradable", ErrMessage: fmt.Sprintf("%v cannot be traded. Only tokens can be traded.", nonEmpty(code, os.Getenv("NATIVE_ASSET_CODE")))}
	}
	return common.HexToAddress(contract), nil
}

func errOffersNotConfigured() error {
	return &tErrors.CustomError{Param: "offerType", Err: "error-market-not-configured", ErrMessage: "Market offers are not available at this time. Please try again later.", Code: http.StatusServiceUnavailable}
}

// MakeOffer places a BUY or SELL offer from sourceWallet. Like Pay it runs
// in two steps: without a signature it builds the offer as a wallet
// operation and returns it (offerRequest.Transaction); with the signature
// it submits it. Shared wallets with approvers get an approval request on
// commit.
func MakeOffer(signerUser, walletOwner *userModels.User, sourceWallet *userModels.UserWallet, offerRequest *userModels.MarketOfferRequest, gc *sharedconfig.GlobalConfig) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	offerRequest.Messages = make([]string, 0)
	offerRequest.NetworkPassPhrase = gc.BantuNetworkPassphrase
	if sourceWallet.SharedAccessEnabled == 1 && sourceWallet.NumberOfApprovalsNeeded > 0 {
		offerRequest.Multiparty = 1
	}
	if sourceWallet.HasViewOnlyAccess(gc) {
		offerRequest.SignatureRequired = 1
	}

	// step 2: the user signed the operation built in step 1
	if offerRequest.Multiparty == 0 && len(offerRequest.TransactionSignature) > 0 && len(offerRequest.Transaction) > 0 {
		rec, p, err := LoadWalletOperation(offerRequest.Transaction, sourceWallet.ID, OperationMarketOffer, gc)
		if err != nil {
			return err
		}
		var c marketOfferContext
		if rec.Context == nil || json.Unmarshal([]byte(*rec.Context), &c) != nil {
			return &tErrors.ErrorTemporaryServerError{}
		}
		hash, err := SignSingleOwnerOperation(ctx, rec, p, signerUser.PrimarySigner, offerRequest.TransactionSignature, gc)
		if err != nil {
			return err
		}
		offerRequest.TransactionID = hash
		c.Offer.TransactionID = &hash
		if e := gc.DB.Omit(clause.Associations).Create(&c.Offer).Error; e != nil {
			log.Printf("[MakeOffer] saving market offer %+v: %v", c.Offer, e)
			gc.LogDiscordFailedRequest(fmt.Sprintf("[MakeOffer] market offer of %v submitted (%v) but not recorded: %v", sourceWallet.ID, hash, e))
		}
		walletOwner.InvalidateUserCache(gc)
		return nil
	}
	// a shared wallet's initiator commits the operation built in step 1
	if offerRequest.Multiparty == 1 && offerRequest.Commit == 1 && len(offerRequest.Transaction) > 0 {
		rec, _, err := LoadWalletOperation(offerRequest.Transaction, sourceWallet.ID, OperationMarketOffer, gc)
		if err != nil {
			return err
		}
		var c marketOfferContext
		if rec.Context == nil || json.Unmarshal([]byte(*rec.Context), &c) != nil {
			return &tErrors.ErrorTemporaryServerError{}
		}
		return createMarketOfferApproval(signerUser, sourceWallet, offerRequest, &c.Offer, gc)
	}

	// step 1: build
	book, err := offerbook.BookAddress()
	if err != nil {
		return errOffersNotConfigured()
	}
	offerType := strings.ToUpper(strings.TrimSpace(offerRequest.OfferType))
	if offerType != "BUY" && offerType != "SELL" {
		return &tErrors.CustomError{Param: "offerType", Err: "error-invalid-paramter", ErrMessage: "Offer type can only be either BUY or SELL."}
	}
	qty, e := decimal.NewFromString(offerRequest.Quantity)
	if e != nil || !qty.IsPositive() {
		return &tErrors.CustomError{Param: "quantity", Err: "error-missing-quantity", ErrMessage: "Quantity must be greater than zero."}
	}
	price, e := decimal.NewFromString(offerRequest.PricePerUnit)
	if e != nil || !price.IsPositive() {
		return &tErrors.CustomError{Param: "pricePerUnit", Err: "error-missing-price", ErrMessage: "Price must be greater than zero."}
	}
	asset, err := marketToken(offerRequest.AssetCode, offerRequest.ContractAddress, "assetCode")
	if err != nil {
		return err
	}
	currency, err := marketToken(offerRequest.CurrencyCode, offerRequest.CurrencyIssuer, "currencyCode")
	if err != nil {
		return err
	}
	if asset == currency {
		return &tErrors.CustomError{Param: "currencyCode", Err: "error-invalid-paramter", ErrMessage: "The asset and the currency must differ."}
	}
	client := gc.BantuExpansionClient
	for _, t := range []common.Address{asset, currency} {
		ok, err := aa.Tradable(ctx, client, book, t)
		if err != nil {
			return &tErrors.ErrorTemporaryServerError{}
		}
		if !ok {
			return &tErrors.CustomError{Param: "assetCode", Err: "error-asset-not-tradable", ErrMessage: "This market is not open for trading."}
		}
	}
	assetDec, err := offerbook.Decimals(ctx, gc.DB, client, asset)
	if err != nil {
		return &tErrors.ErrorTemporaryServerError{}
	}
	currencyDec, err := offerbook.Decimals(ctx, gc.DB, client, currency)
	if err != nil {
		return &tErrors.ErrorTemporaryServerError{}
	}

	var sellToken, payToken common.Address
	var sellCode string
	var sellDec uint8
	var amount *big.Int
	var p aa.OfferPrice
	if offerType == "SELL" {
		sellToken, payToken, sellCode, sellDec = asset, currency, offerRequest.AssetCode, assetDec
		amount = qty.Shift(int32(assetDec)).Truncate(0).BigInt()
		p, err = offerPrice(price, assetDec, currencyDec)
	} else {
		sellToken, payToken, sellCode, sellDec = currency, asset, offerRequest.CurrencyCode, currencyDec
		amount = qty.Mul(price).Shift(int32(currencyDec)).Truncate(0).BigInt()
		p, err = inverseOfferPrice(price, assetDec, currencyDec)
	}
	if err != nil {
		return &tErrors.CustomError{Param: "pricePerUnit", Err: "error-invalid-price", ErrMessage: "This price cannot be used. Please use fewer decimal places."}
	}
	if amount.Sign() <= 0 {
		return &tErrors.CustomError{Param: "quantity", Err: "error-missing-quantity", ErrMessage: "Quantity is too small."}
	}
	sellAsset := basetxn.CreditAsset{Code: sellCode, Issuer: sellToken.Hex()}
	_, authorized, _, balance, _, err := network.BlockchainAccountProperties(client, sourceWallet.ID, sellAsset)
	if err != nil {
		return err
	}
	need := decimal.NewFromBigInt(amount, -int32(sellDec))
	if !authorized || balance.LessThan(need) {
		return &tErrors.ErrorUnderfundedAccount{Detail: fmt.Sprintf("Not enough funds. The offer needs %v %v.", need, sellCode)}
	}
	create, err := aa.CreateOfferCall(book, sellToken, amount, common.HexToAddress(sourceWallet.ID), []common.Address{payToken}, []aa.OfferPrice{p})
	if err != nil {
		return &tErrors.ErrorTemporaryServerError{}
	}
	calls := []aa.Call{aa.ERC20Approve(sellToken, book, amount), create}

	ca, ci := asset.Hex(), currency.Hex()
	draft := userModels.MarketOffer{
		ID:                        uuid.NewString(),
		SourceWalletAlias:         sourceWallet.Alias,
		SourceWalletAddress:       sourceWallet.ID,
		MarketMakingWalletAddress: sourceWallet.ID,
		OfferType:                 offerType,
		AssetCode:                 offerRequest.AssetCode,
		ContractAddress:           &ca,
		CurrencyCode:              offerRequest.CurrencyCode,
		CurrencyIssuer:            &ci,
		PricePerUnit:              price.String(),
		Quantity:                  qty.String(),
		FeeChargedOnAsset:         "0",
		FeeValue:                  "0",
		NetQuantity:               qty.String(),
	}
	offerRequest.NetQuantity, offerRequest.FeeChargedOnAsset, offerRequest.FeeValue = draft.NetQuantity, "0", "0"

	validity := time.Duration(0)
	if offerRequest.Multiparty == 1 {
		validity = sharedWalletOperationValidity()
	}
	op, err := PrepareWalletOperation(ctx, OperationMarketOffer, signerUser, walletOwner, sourceWallet, calls, validity, marketOfferContext{Offer: draft, Book: book.Hex()}, gc)
	if err != nil {
		return err
	}
	offerRequest.Transaction = op.Transaction
	offerRequest.Messages = append(offerRequest.Messages, fmt.Sprintf("%v %v of this wallet is held by the offer until it is filled or cancelled.", need, sellCode))
	offerRequest.Messages = append(offerRequest.Messages, op.Messages()...)
	offerRequest.SignatureRequired = 1
	if offerRequest.Multiparty == 1 && offerRequest.Commit == 1 {
		return createMarketOfferApproval(signerUser, sourceWallet, offerRequest, &draft, gc)
	}
	return nil
}

func createMarketOfferApproval(signerUser *userModels.User, sourceWallet *userModels.UserWallet, offerRequest *userModels.MarketOfferRequest, offer *userModels.MarketOffer, gc *sharedconfig.GlobalConfig) error {
	offerRequest.TransactionID = "PENDING_AUTH"
	description := fmt.Sprintf("%v %v %v @ %v %v", offer.OfferType, offer.Quantity, offer.AssetCode, offer.PricePerUnit, offer.CurrencyCode)
	if len(offerRequest.Messages) > 0 {
		description = fmt.Sprintf("%s\nMessages: %v", description, strings.Join(offerRequest.Messages, "\n"))
	}
	offerRequest.ReturnedDescription = description
	// the approvals record the offer once it is submitted
	raw, _ := json.Marshal(*offer)
	info := string(raw)
	pendingAuth := userModels.PendingAuth{
		ID:                     uuid.NewString(),
		Initiator:              signerUser.Username,
		InitiatorSignerAddress: signerUser.PrimarySigner,
		WalletAddress:          sourceWallet.ID,
		TransactionType:        OperationMarketOffer,
		Description:            description,
		TransactionSource:      sourceWallet.ID,
		ApprovalsNeeded:        sourceWallet.NumberOfApprovalsNeeded,
		TransactionXdr:         offerRequest.Transaction,
		TransactionInfoStr:     &info,
	}
	if e := gc.DB.Omit(clause.Associations).Create(&pendingAuth).Error; e != nil {
		log.Printf("[MakeOffer] saving approval request for %v: %v", sourceWallet.ID, e)
		return &tErrors.ErrorTemporaryServerError{}
	}
	return nil
}

// marketOfferBookID is a market offer's id on the book: recorded when its
// operation was mined, else found from that operation's transaction.
func marketOfferBookID(mo *userModels.MarketOffer, gc *sharedconfig.GlobalConfig) (*big.Int, bool) {
	if mo.BlockchainOfferID != nil {
		if id, ok := new(big.Int).SetString(*mo.BlockchainOfferID, 10); ok {
			return id, true
		}
	}
	if mo.TransactionID == nil {
		return nil, false
	}
	var op userModels.WalletOperation
	if gc.DB.Where("user_op_hash = ?", *mo.TransactionID).First(&op).Error != nil || op.TxHash == nil {
		return nil, false
	}
	var o offerbook.Offer
	if gc.DB.Where("created_tx = ? AND seller = ?", *op.TxHash, offerbook.Key(mo.SourceWalletAddress)).First(&o).Error != nil {
		return nil, false
	}
	gc.DB.Model(&userModels.MarketOffer{}).Where("id = ?", mo.ID).Update("blockchain_offer_id", o.ID)
	mo.BlockchainOfferID = &o.ID
	id, ok := new(big.Int).SetString(o.ID, 10)
	return id, ok
}

// recordMarketOffer stores the book id of the offer a mined market-making
// operation created.
func recordMarketOffer(op userModels.WalletOperation, txHash common.Hash, gc *sharedconfig.GlobalConfig) {
	var c marketOfferContext
	if op.Context == nil || json.Unmarshal([]byte(*op.Context), &c) != nil || op.UserOpHash == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	rcpt, err := gc.BantuExpansionClient.TransactionReceipt(ctx, txHash)
	if err != nil {
		log.Printf("[recordMarketOffer] receipt of %v: %v", txHash.Hex(), err)
		return // found later from the index (marketOfferBookID)
	}
	id, ok := aa.ParseOfferCreated(rcpt.Logs, common.HexToAddress(c.Book), common.HexToAddress(op.WalletAddress))
	if !ok {
		return
	}
	s := id.String()
	gc.DB.Model(&userModels.MarketOffer{}).Where("transaction_id = ?", *op.UserOpHash).Update("blockchain_offer_id", s)
}

// recordMarketOfferCancelled marks a market offer cancelled once its
// cancellation was mined.
func recordMarketOfferCancelled(op userModels.WalletOperation, txHash common.Hash, gc *sharedconfig.GlobalConfig) {
	var c marketOfferContext
	if op.Context == nil || json.Unmarshal([]byte(*op.Context), &c) != nil || c.Offer.ID == "" {
		return
	}
	gc.DB.Model(&userModels.MarketOffer{}).Where("id = ?", c.Offer.ID).Update("canceled", 1)
}

func init() {
	minedHooks[OperationMarketOffer] = recordMarketOffer
	minedHooks[OperationCancelMarketOffer] = recordMarketOfferCancelled
}

// MarketOfferStatus is a market offer with its state on the book.
type MarketOfferStatus struct {
	userModels.MarketOffer
	Open      bool   `json:"open"`
	Remaining string `json:"remaining"` // what the offer still sells (whole tokens)
	Sold      string `json:"sold"`      // what it has sold so far, in the token it sells (whole tokens)
	Received  string `json:"received"`  // what it has been paid so far, in the token it buys (whole tokens)
	Fills     int    `json:"fills"`
}

// ListOffers lists sourceWallet's market offers with their state on the
// book.
func ListOffers(sourceWallet *userModels.UserWallet, gc *sharedconfig.GlobalConfig) ([]MarketOfferStatus, error) {
	var offers []userModels.MarketOffer
	if err := gc.DB.Where("source_wallet_address = ?", sourceWallet.ID).Order("created_at DESC").Limit(200).Find(&offers).Error; err != nil {
		return nil, &tErrors.ErrorTemporaryServerError{}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	out := make([]MarketOfferStatus, 0, len(offers))
	for i := range offers {
		s := MarketOfferStatus{MarketOffer: offers[i], Remaining: "0", Sold: "0", Received: "0"}
		if id, ok := marketOfferBookID(&s.MarketOffer, gc); ok {
			var o offerbook.Offer
			if gc.DB.First(&o, "id = ?", id.String()).Error == nil {
				s.Open = o.Open
				sellDec, err := offerbook.Decimals(ctx, gc.DB, gc.BantuExpansionClient, common.HexToAddress(o.SellToken))
				if rem, _ := new(big.Int).SetString(o.Remaining, 10); err == nil && rem != nil {
					s.Remaining = decimal.NewFromBigInt(rem, -int32(sellDec)).String()
				}
				var fills []offerbook.Fill
				gc.DB.Where("offer_id = ?", o.ID).Find(&fills)
				sold, received := new(big.Int), new(big.Int)
				payToken := ""
				for _, f := range fills {
					if a, ok := new(big.Int).SetString(f.Amount, 10); ok {
						sold.Add(sold, a)
					}
					if p, ok := new(big.Int).SetString(f.Payment, 10); ok {
						received.Add(received, p)
					}
					payToken = f.PaymentToken
				}
				s.Fills = len(fills)
				if err == nil {
					s.Sold = decimal.NewFromBigInt(sold, -int32(sellDec)).String()
				}
				if payToken != "" {
					if dec, err := offerbook.Decimals(ctx, gc.DB, gc.BantuExpansionClient, common.HexToAddress(payToken)); err == nil {
						s.Received = decimal.NewFromBigInt(received, -int32(dec)).String()
					}
				}
			}
		}
		out = append(out, s)
	}
	return out, nil
}

// CancelOffer cancels one of sourceWallet's market offers, returning what
// is left of it to the wallet. It runs in the same two steps as MakeOffer.
func CancelOffer(signerUser, walletOwner *userModels.User, sourceWallet *userModels.UserWallet, deleteOfferRequest *userModels.DeleteOfferRequest, gc *sharedconfig.GlobalConfig) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	deleteOfferRequest.Messages = make([]string, 0)
	deleteOfferRequest.NetworkPassPhrase = gc.BantuNetworkPassphrase
	if sourceWallet.SharedAccessEnabled == 1 && sourceWallet.NumberOfApprovalsNeeded > 0 {
		deleteOfferRequest.Multiparty = 1
	}
	if sourceWallet.HasViewOnlyAccess(gc) {
		deleteOfferRequest.SignatureRequired = 1
	}
	if len(deleteOfferRequest.ID) == 0 {
		return &tErrors.ErrorMissingParameter{Parameter: "Id"}
	}
	marketOffer, err := sourceWallet.GetMarketOfferByID(deleteOfferRequest.ID, gc.DB, gc)
	if err != nil {
		return err
	}
	if marketOffer.Canceled == 1 {
		return &tErrors.CustomError{Param: "id", Err: "error-offer-has been canceled", ErrMessage: "Offer cannot be canceled because it has been already been canceled before."}
	}

	// step 2: the user signed the operation built in step 1
	if deleteOfferRequest.Multiparty == 0 && len(deleteOfferRequest.TransactionSignature) > 0 && len(deleteOfferRequest.Transaction) > 0 {
		rec, p, err := LoadWalletOperation(deleteOfferRequest.Transaction, sourceWallet.ID, OperationCancelMarketOffer, gc)
		if err != nil {
			return err
		}
		hash, err := SignSingleOwnerOperation(ctx, rec, p, signerUser.PrimarySigner, deleteOfferRequest.TransactionSignature, gc)
		if err != nil {
			return err
		}
		deleteOfferRequest.TransactionID = hash
		walletOwner.InvalidateUserCache(gc)
		return nil
	}
	if deleteOfferRequest.Multiparty == 1 && deleteOfferRequest.Commit == 1 && len(deleteOfferRequest.Transaction) > 0 {
		if _, _, err := LoadWalletOperation(deleteOfferRequest.Transaction, sourceWallet.ID, OperationCancelMarketOffer, gc); err != nil {
			return err
		}
		return createCancelOfferApproval(signerUser, sourceWallet, deleteOfferRequest, &marketOffer, gc)
	}

	// step 1: build
	book, err := offerbook.BookAddress()
	if err != nil {
		return errOffersNotConfigured()
	}
	id, ok := marketOfferBookID(&marketOffer, gc)
	if !ok {
		return &tErrors.CustomError{Param: "id", Err: "error-offer-not-yet-placed", ErrMessage: "This offer is not on the market yet. Please try again shortly.", Code: http.StatusConflict}
	}
	o, err := aa.ReadOffer(ctx, gc.BantuExpansionClient, book, id)
	if err != nil {
		return &tErrors.ErrorTemporaryServerError{}
	}
	if !o.Open {
		gc.DB.Model(&userModels.MarketOffer{}).Where("id = ?", marketOffer.ID).Update("canceled", 1)
		return &tErrors.CustomError{Param: "id", Err: "error-offer-has been canceled", ErrMessage: "Offer cannot be canceled because it has been already been canceled before."}
	}
	if o.Remaining.Sign() == 0 {
		return &tErrors.CustomError{Param: "id", Err: "error-offer-has been filled", ErrMessage: "Offer cannot be canceled because it has been filled."}
	}
	validity := time.Duration(0)
	if deleteOfferRequest.Multiparty == 1 {
		validity = sharedWalletOperationValidity()
	}
	op, err := PrepareWalletOperation(ctx, OperationCancelMarketOffer, signerUser, walletOwner, sourceWallet, []aa.Call{aa.CancelOfferCall(book, id)}, validity, marketOfferContext{Offer: marketOffer, Book: book.Hex()}, gc)
	if err != nil {
		return err
	}
	deleteOfferRequest.Transaction = op.Transaction
	deleteOfferRequest.Messages = append(deleteOfferRequest.Messages, op.Messages()...)
	deleteOfferRequest.SignatureRequired = 1
	if deleteOfferRequest.Multiparty == 1 && deleteOfferRequest.Commit == 1 {
		return createCancelOfferApproval(signerUser, sourceWallet, deleteOfferRequest, &marketOffer, gc)
	}
	return nil
}

func createCancelOfferApproval(signerUser *userModels.User, sourceWallet *userModels.UserWallet, req *userModels.DeleteOfferRequest, offer *userModels.MarketOffer, gc *sharedconfig.GlobalConfig) error {
	req.TransactionID = "PENDING_AUTH"
	req.ReturnedDescription = fmt.Sprintf("Cancel %v %v %v @ %v %v", offer.OfferType, offer.Quantity, offer.AssetCode, offer.PricePerUnit, offer.CurrencyCode)
	raw, _ := json.Marshal(*req)
	info := string(raw)
	pendingAuth := userModels.PendingAuth{
		ID:                     uuid.NewString(),
		Initiator:              signerUser.Username,
		InitiatorSignerAddress: signerUser.PrimarySigner,
		WalletAddress:          sourceWallet.ID,
		TransactionType:        OperationCancelMarketOffer,
		Description:            req.ReturnedDescription,
		TransactionSource:      sourceWallet.ID,
		ApprovalsNeeded:        sourceWallet.NumberOfApprovalsNeeded,
		TransactionXdr:         req.Transaction,
		TransactionInfoStr:     &info,
	}
	if e := gc.DB.Omit(clause.Associations).Create(&pendingAuth).Error; e != nil {
		log.Printf("[CancelOffer] saving approval request for %v: %v", sourceWallet.ID, e)
		return &tErrors.ErrorTemporaryServerError{}
	}
	return nil
}

// ToFractionInt32 approximates a float64 as a fraction with int32 numerator and denominator.
// It uses a continued fraction algorithm for a good approximation.
func ToFractionInt32(x float64) (num, den int32) {
	if x < 0 {
		return 0, 0 // Handle negative numbers as per requirement
	}

	const (
		maxDen = math.MaxInt32
		maxNum = math.MaxInt32
	)

	// Initial values for the continued fraction
	n1, d1 := int64(math.Floor(x)), int64(1)
	n2, d2 := int64(1), int64(0)
	x -= math.Floor(x)

	for x != 0 {
		a := math.Floor(1 / x)
		n3 := n1*int64(a) + n2
		d3 := d1*int64(a) + d2

		if d3 > maxDen || n3 > maxNum {
			break // Stop if numerator or denominator exceeds int32 limits
		}

		// Update values for the next iteration
		n2, d2 = n1, d1
		n1, d1 = n3, d3
		x = 1/x - a
	}

	return int32(n1), int32(d1)
}
