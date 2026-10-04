package users

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"strings"
	"time"

	"trovo-wallet-api/internal/aa"
	"trovo-wallet-api/internal/basetxn"
	paymentModels "trovo-wallet-api/internal/components/payments/models"
	userModels "trovo-wallet-api/internal/components/users/models"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/evmkeypair"
	"trovo-wallet-api/internal/gnosissafe"
	"trovo-wallet-api/internal/network"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm/clause"
)

// OperationAssetSubscription is the wallet operation kind of a tokenized
// asset purchase paid in a stablecoin.
const OperationAssetSubscription = "ASSET SUBSCRIPTION"

// assetSubscriptionContext is what a purchase operation keeps for when it is
// submitted: the request as it was priced.
type assetSubscriptionContext struct {
	Input userModels.TokenizedAssetSubscriptionInput `json:"input"`
}

// checkCanSubscribe applies the platform's purchase rules: an active,
// KYC-verified owner, an asset on sale, and the per-owner purchase cap.
func checkCanSubscribe(subscriber *userModels.User, wallet *userModels.UserWallet, ta *userModels.TokenizedAsset, amount float64, gc *sharedconfig.GlobalConfig) error {
	if err := subscriber.EnsureNotSuspended(); err != nil {
		return err
	}
	owner, err := wallet.GetWalletOwner(gc.DB, gc)
	if err != nil {
		return &tErrors.CustomError{Param: "walletAddress", Err: "error-invalid-kyc", ErrMessage: "Unable to verify wallet owner."}
	}
	if owner.KYCVerified == 0 {
		return &tErrors.CustomError{Param: "walletAddress", Err: "error-invalid-kyc", ErrMessage: fmt.Sprintf("%v has not passed KYC to purchase this tokenized asset %v.", owner.Username, *ta.AssetCode)}
	}
	if ta.AssetTokenizationStatus != 5 && ta.AssetTokenizationStatus != 6 {
		return &tErrors.CustomError{Param: "tokenizedAssetId", Err: "error-invalid-request", ErrMessage: "Only projects that are in sales can accept purchase."}
	}
	if amount <= 0 {
		return &tErrors.CustomError{Param: "amount", Err: "error-invalid-amount", ErrMessage: "The purchase amount must be greater than zero."}
	}
	capEnd := ta.SalesStart.AddDate(0, 0, ta.CapDurationInDays)
	if capEnd.After(time.Now()) && ta.CapOnPurchase > 0 && ta.CapAmountInFiat > 0 {
		bought := ta.SumAmountBoughtByWalletOwner(wallet.Alias, gc)
		if decimal.NewFromFloat(amount + bought).Truncate(7).GreaterThan(decimal.NewFromFloat(ta.CapAmountInFiat)) {
			most := decimal.NewFromFloat(ta.CapAmountInFiat - bought)
			if !most.IsPositive() {
				return &tErrors.CustomError{Param: "amount", Err: "error-cap-amount-exceeded", ErrMessage: fmt.Sprintf("You have exhausted the allowed purchase cap of %v%v worth of %v at this time.", *ta.AssetQuoteCurrency, ta.CapAmountInFiat, *ta.AssetCode)}
			}
			return &tErrors.CustomError{Param: "amount", Err: "error-cap-amount-exceeded", ErrMessage: fmt.Sprintf("You can only purchase upto %v%v worth of %v at this time.", *ta.AssetQuoteCurrency, most, *ta.AssetCode)}
		}
	}
	return nil
}

// fillValidity is how long a fill authorization stays usable: as long as
// the operation carrying it.
func fillValidity(opValidity time.Duration, gc *sharedconfig.GlobalConfig) time.Duration {
	if opValidity > 0 {
		return opValidity
	}
	return operationBuilder(gc).DefaultValidity
}

// SubscribeToTokenizedAsset buys a tokenized asset with a stablecoin from
// subscriberWallet: one operation of the wallet approves the offer book for
// the payment and fills the asset's sale offer, so payment and delivery
// happen together or not at all. It runs in two steps:
//
//  1. without transactionSignature: prices the purchase (the most of the
//     asset input.Amount of the payment stablecoin buys), has the platform
//     authorize it and returns the operation as transaction;
//  2. with transaction + transactionSignature: submits it and records the
//     subscription.
//
// Shared wallets with approvers get an approval request on commit.
func SubscribeToTokenizedAsset(subscriber *userModels.User, subscriberWallet *userModels.UserWallet, ta *userModels.TokenizedAsset, input *userModels.TokenizedAssetSubscriptionInput, gc *sharedconfig.GlobalConfig) (taSubscription userModels.TokenizedAssetSubscription, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	input.TokenizedAssetID = ta.ID
	input.WalletAddress = subscriberWallet.ID
	input.SubscriberUsername = subscriber.Username
	input.Amount = decimal.NewFromFloat(input.Amount).Truncate(7).InexactFloat64()
	input.NetworkPassPhrase = network.GetBlockchainNetworkPassPhrase()
	if subscriberWallet.SharedAccessEnabled == 1 && subscriberWallet.NumberOfApprovalsNeeded > 0 {
		input.Multiparty = 1
	}
	if subscriberWallet.HasViewOnlyAccess(gc) {
		input.SignatureRequired = 1
	}

	if len(input.TransactionSignature) > 0 && len(input.Transaction) > 0 && input.Multiparty == 0 {
		return submitAssetSubscription(ctx, subscriber, subscriberWallet, ta, input, gc)
	}
	if input.Multiparty == 1 && input.Commit == 1 && len(input.Transaction) > 0 {
		rec, _, e := LoadWalletOperation(input.Transaction, subscriberWallet.ID, OperationAssetSubscription, gc)
		if e != nil {
			return taSubscription, e
		}
		var c assetSubscriptionContext
		if rec.Context == nil || json.Unmarshal([]byte(*rec.Context), &c) != nil {
			return taSubscription, &tErrors.ErrorTemporaryServerError{}
		}
		c.Input.Transaction, c.Input.Multiparty, c.Input.Commit = input.Transaction, 1, 1
		*input = c.Input
		return taSubscription, createAssetSubscriptionApprovalRequest(subscriber, subscriberWallet, ta, input, gc)
	}

	if err = checkCanSubscribe(subscriber, subscriberWallet, ta, input.Amount, gc); err != nil {
		return
	}
	// the payment stablecoin: any tokenization currency (CNGN by default),
	// resolved server-side
	code := strings.ToUpper(strings.TrimSpace(input.PaymentAssetCode))
	if code == "" {
		code = "CNGN"
	}
	currency := GetTokenizationCurrencyByCode(code, gc.DB)
	if currency.AssetCode == "" || !common.IsHexAddress(currency.ContractAddress) {
		err = &tErrors.CustomError{Param: "paymentAssetCode", Err: "error-invalid-payment-asset", ErrMessage: fmt.Sprintf("Payment asset code [%v] you supplied is invalid.", code)}
		return
	}
	if input.PaymentContractAddress != "" && !strings.EqualFold(input.PaymentContractAddress, currency.ContractAddress) {
		err = &tErrors.CustomError{Param: "paymentContractAddress", Err: "error-invalid-payment-asset", ErrMessage: "Payment asset issuer does not match the registered issuer for this currency."}
		return
	}
	if userModels.IsInternalBalanceAsset(currency.AssetCode, currency.ContractAddress, gc) || strings.EqualFold(currency.AssetCode, *ta.AssetCode) {
		err = &tErrors.CustomError{Param: "paymentAssetCode", Err: "error-invalid-payment-asset", ErrMessage: "Payment asset cannot be the internal balance token or the tokenized asset itself."}
		return
	}
	input.PaymentAssetCode, input.PaymentContractAddress = currency.AssetCode, currency.ContractAddress

	o, err := loadAssetOffer(ctx, ta, gc)
	if err != nil {
		return
	}
	validity := time.Duration(0)
	if input.Multiparty == 1 {
		validity = sharedWalletOperationValidity()
	}
	wallet := common.HexToAddress(subscriberWallet.ID)
	purchase, err := priceAssetPurchase(ctx, o, currency.AssetCode, common.HexToAddress(currency.ContractAddress), decimal.NewFromFloat(input.Amount), wallet, wallet, fillValidity(validity, gc), gc)
	if err != nil {
		return
	}
	calls, err := purchase.purchaseCalls(o.Book)
	if err != nil {
		return taSubscription, &tErrors.ErrorTemporaryServerError{}
	}
	owner, err := subscriberWallet.GetWalletOwner(gc.DB, gc)
	if err != nil {
		return taSubscription, &tErrors.ErrorTemporaryServerError{}
	}
	got := purchase.AssetAmountHuman(o.AssetDecimals)
	input.SwappedEstimate = got.String()
	input.Memo = fmt.Sprintf("%v>%v", currency.AssetCode, *ta.AssetCode)
	op, err := PrepareWalletOperation(ctx, OperationAssetSubscription, subscriber, &owner, subscriberWallet, calls, validity, assetSubscriptionContext{Input: *input}, gc)
	if err != nil {
		return
	}
	pay := paymentModels.PaymentInfo{AssetCode: currency.AssetCode, ContractAddress: currency.ContractAddress, AmountToPay: purchase.PaymentHuman().String()}
	if err = checkPaymentLeavesGas(ctx, subscriberWallet, &pay, op.Prepared, gc); err != nil {
		return
	}
	// the platform's (database) record that this wallet may hold the asset
	if e := network.SetWalletAssetAuthorization(subscriberWallet.ID, basetxn.CreditAsset{Code: *ta.AssetCode, Issuer: o.Token.Hex()}, true, *ta.IssuingWalletAddress, "tokenized-asset purchase"); e != nil {
		log.Printf("[SubscribeToTokenizedAsset] authorizing %v for %v: %v", subscriberWallet.ID, *ta.AssetCode, e)
		return taSubscription, &tErrors.ErrorTemporaryServerError{}
	}

	input.Transaction = op.Transaction
	input.SignatureRequired = 1
	input.Messages = append([]string{fmt.Sprintf("You pay %v %v and receive %v %v in this wallet.", purchase.PaymentHuman(), currency.AssetCode, got, *ta.AssetCode)}, op.Messages()...)
	taSubscription.UpdateTokenizedAssetSubscriptionFromInput(subscriber.Username, subscriberWallet, input, ta, gc)
	if input.Multiparty == 1 && input.Commit == 1 {
		err = createAssetSubscriptionApprovalRequest(subscriber, subscriberWallet, ta, input, gc)
	}
	return
}

// submitAssetSubscription submits a signed purchase and records it.
func submitAssetSubscription(ctx context.Context, subscriber *userModels.User, wallet *userModels.UserWallet, ta *userModels.TokenizedAsset, input *userModels.TokenizedAssetSubscriptionInput, gc *sharedconfig.GlobalConfig) (taSubscription userModels.TokenizedAssetSubscription, err error) {
	rec, p, err := LoadWalletOperation(input.Transaction, wallet.ID, OperationAssetSubscription, gc)
	if err != nil {
		return
	}
	var c assetSubscriptionContext
	if rec.Context == nil || json.Unmarshal([]byte(*rec.Context), &c) != nil {
		return taSubscription, &tErrors.ErrorTemporaryServerError{}
	}
	hash, err := SignSingleOwnerOperation(ctx, rec, p, subscriber.PrimarySigner, input.TransactionSignature, gc)
	if err != nil {
		return
	}
	stored := c.Input
	stored.Transaction, stored.TransactionID = input.Transaction, hash
	*input = stored
	taSubscription.UpdateTokenizedAssetSubscriptionFromInput(subscriber.Username, wallet, input, ta, gc)
	taSubscription.TransactionID = hash
	if e := gc.DB.Omit(clause.Associations).Save(&taSubscription).Error; e != nil {
		log.Printf("[SubscribeToTokenizedAsset] operation %v submitted but subscription not saved: %v", hash, e)
		gc.LogDiscordFailedRequest(fmt.Sprintf("[SubscribeToTokenizedAsset] purchase %v of %v by %v submitted but its subscription was not saved: %v", hash, *ta.AssetCode, wallet.ID, e))
	}
	wallet.InvalidateUserCache(gc)
	subscriber.InvalidateUserWalletCache(gc)
	return taSubscription, nil
}

// createAssetSubscriptionApprovalRequest records a shared wallet's purchase
// for its approvers.
func createAssetSubscriptionApprovalRequest(subscriber *userModels.User, wallet *userModels.UserWallet, ta *userModels.TokenizedAsset, input *userModels.TokenizedAssetSubscriptionInput, gc *sharedconfig.GlobalConfig) error {
	input.TransactionID = "PENDING_AUTH"
	description := fmt.Sprintf("Buying Tokenized asset [%v]\n Amount:%v %v,\n Getting:%v %v", *ta.AssetName, input.Amount, input.PaymentAssetCode, input.SwappedEstimate, *ta.AssetCode)
	if len(input.Messages) > 0 {
		description = fmt.Sprintf("%v\nMessages: %v", description, strings.Join(input.Messages, "\n"))
	}
	input.ReturnedDescription = description
	info, _ := json.Marshal(*input)
	infoStr := string(info)
	pendingAuth := userModels.PendingAuth{
		ID:                     uuid.NewString(),
		Initiator:              subscriber.Username,
		InitiatorSignerAddress: subscriber.PrimarySigner,
		WalletAddress:          wallet.ID,
		TransactionType:        "ASSET SUBSCRIPTION",
		Description:            description,
		ApprovalsNeeded:        wallet.NumberOfApprovalsNeeded,
		TransactionXdr:         input.Transaction,
		TransactionInfoStr:     &infoStr,
	}
	if err := gc.DB.Omit(clause.Associations).Create(&pendingAuth).Error; err != nil {
		log.Printf("[SubscribeToTokenizedAsset] saving approval request for %v: %v", wallet.ID, err)
		return &tErrors.ErrorTemporaryServerError{}
	}
	wallet.InvalidateUserCache(gc)
	subscriber.InvalidateUserWalletCache(gc)
	return nil
}

// fiatPurchaseStatement is what a fiat buyer signs: there is nothing for
// their wallet to send (the platform pays on their behalf once the fiat
// payment is confirmed), so their signature records their consent.
func fiatPurchaseStatement(invoiceID string, amount float64, currency, assetCode, wallet string) string {
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("Trovo fiat purchase %v: pay %v %v for %v delivered to wallet %v", invoiceID, decimal.NewFromFloat(amount), currency, assetCode, wallet)))
}

// internalBalanceMinter is the country's internal balance token, the Safe
// that mints it, and the keys of its owners
// (INTERNAL_BALANCE_ISSUING_SIGNERS).
func internalBalanceMinter(ctx context.Context, ta *userModels.TokenizedAsset, gc *sharedconfig.GlobalConfig) (code string, token, safe common.Address, signers []*evmkeypair.Full, err error) {
	fail := &tErrors.CustomError{Param: "tokenizedAssetId", Err: "error-fiat-purchase-not-configured", ErrMessage: "Fiat purchases are not configured for this asset's country.", Code: http.StatusServiceUnavailable}
	if ta.AssetCountryLocation == nil {
		return "", token, safe, nil, fail
	}
	cc := userModels.CountryCode(*ta.AssetCountryLocation).GetConfig(gc)
	if cc.InternalBalanceTokenCode == nil || cc.InternalTokenIssuer == nil || cc.InternalTokenMinterSafe == nil || !common.IsHexAddress(*cc.InternalTokenIssuer) || !common.IsHexAddress(*cc.InternalTokenMinterSafe) {
		return "", token, safe, nil, fail
	}
	candidates, e := getInternalBalanceIssuingSigners()
	if e != nil {
		log.Printf("[internalBalanceMinter] %v", e)
		return "", token, safe, nil, fail
	}
	safe = common.HexToAddress(*cc.InternalTokenMinterSafe)
	signers, e = gnosissafe.SignersForSafe(ctx, gc.BantuExpansionClient, safe, candidates)
	if e != nil {
		log.Printf("[internalBalanceMinter] signers for %v: %v", safe.Hex(), e)
		return "", token, safe, nil, fail
	}
	return *cc.InternalBalanceTokenCode, common.HexToAddress(*cc.InternalTokenIssuer), safe, signers, nil
}

// SubscribeToTokenizedAssetByFiat is the fiat counterpart of
// SubscribeToTokenizedAsset. The buyer pays the payment provider
// (Flutterwave); once the webhook confirms it, the internal balance
// token's minting Safe mints the purchase amount to itself and fills the
// asset's sale offer with the buyer's wallet as recipient
// (DeliverFiatAssetPurchase).
//
// input.ID is required on both calls: it keys the FiatPaymentInvoice and
// the TokenizedAssetSubscription and doubles as the provider's reference.
// Call 1 (no transactionSignature) creates the invoice and returns the
// statement the buyer signs as transaction; call 2 stores the signature.
func SubscribeToTokenizedAssetByFiat(subscriber *userModels.User, subscriberWallet *userModels.UserWallet, ta *userModels.TokenizedAsset, input *userModels.FiatTokenizedAssetSubscriptionInput, gc *sharedconfig.GlobalConfig) (invoice userModels.FiatPaymentInvoice, err error) {
	if strings.TrimSpace(input.ID) == "" {
		return invoice, &tErrors.CustomError{Param: "id", Err: "error-invalid-request", ErrMessage: "id is required."}
	}
	if len(input.TransactionSignature) > 0 {
		if e := gc.DB.Where("id = ? AND username = ? AND payment_type = ?", input.ID, subscriber.Username, "ASSET PURCHASE").First(&invoice).Error; e != nil {
			return invoice, &tErrors.CustomError{Param: "id", Err: "error-invalid-request", ErrMessage: "No fiat asset purchase invoice found for this id."}
		}
		if invoice.TransactionSignature != nil || invoice.Status != "PENDING" || invoice.Transaction == nil {
			return invoice, &tErrors.CustomError{Param: "id", Err: "error-invalid-request", ErrMessage: "This fiat asset purchase invoice has already been signed or processed."}
		}
		if err = verifyStatementSignature(subscriber.PrimarySigner, *invoice.Transaction, input.TransactionSignature); err != nil {
			return
		}
		sig := input.TransactionSignature
		if e := gc.DB.Model(&userModels.FiatPaymentInvoice{}).Where("id = ?", invoice.ID).Update("transaction_signature", sig).Error; e != nil {
			return invoice, &tErrors.ErrorTemporaryServerError{}
		}
		invoice.TransactionSignature = &sig
		return invoice, nil
	}
	var existing userModels.FiatPaymentInvoice
	if gc.DB.Where("id = ?", input.ID).First(&existing).Error == nil {
		if existing.TransactionSignature != nil || !strings.EqualFold(existing.Username, subscriber.Username) {
			return invoice, &tErrors.CustomError{Param: "id", Err: "error-invalid-request", ErrMessage: "This fiat asset purchase invoice has already been signed."}
		}
		return existing, nil
	}

	if err = checkCanSubscribe(subscriber, subscriberWallet, ta, input.Amount, gc); err != nil {
		return
	}
	if subscriberWallet.SharedAccessEnabled == 1 && subscriberWallet.NumberOfApprovalsNeeded > 0 {
		return invoice, &tErrors.CustomError{Param: "walletAddress", Err: "error-invalid-request", ErrMessage: "Fiat purchase is not supported on shared-access wallets that require multiple approvals."}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	code, token, minter, _, err := internalBalanceMinter(ctx, ta, gc)
	if err != nil {
		return
	}
	o, err := loadAssetOffer(ctx, ta, gc)
	if err != nil {
		return
	}
	// price it now so the buyer is not charged for a purchase the offer
	// cannot fill (it is priced again, and authorized, on delivery)
	wallet := common.HexToAddress(subscriberWallet.ID)
	if _, err = priceAssetPurchase(ctx, o, code, token, decimal.NewFromFloat(input.Amount), minter, wallet, time.Minute, gc); err != nil {
		return
	}

	statement := fiatPurchaseStatement(input.ID, input.Amount, *ta.AssetQuoteCurrency, *ta.AssetCode, subscriberWallet.ID)
	sub := userModels.TokenizedAssetSubscription{
		ID: input.ID, TokenizedAssetID: ta.ID, AssetCode: *ta.AssetCode, ContractAddress: o.Token.Hex(),
		WalletAlias: subscriberWallet.Alias, WalletAddress: subscriberWallet.ID, Amount: decimal.NewFromFloat(input.Amount).Truncate(7).InexactFloat64(),
		Price: ta.PricePerToken, SubscriberUsername: subscriber.Username, PaymentAssetCode: code, PaymentContractAddress: "FIAT",
	}
	alias, addr, taID := subscriberWallet.Alias, subscriberWallet.ID, ta.ID
	invoice = userModels.FiatPaymentInvoice{
		ID: input.ID, ServiceProvider: "flutterwave", Username: subscriber.Username, Amount: input.Amount, PaymentType: "ASSET PURCHASE",
		Status: "PENDING", TokenizedAssetID: &taID, WalletAlias: &alias, WalletAddress: &addr, Transaction: &statement,
	}
	dbTX := gc.DB.Begin()
	defer dbTX.Rollback()
	if e := dbTX.Omit(clause.Associations).Create(&sub).Error; e != nil {
		log.Printf("[SubscribeToTokenizedAssetByFiat] saving subscription %v: %v", input.ID, e)
		return invoice, &tErrors.ErrorTemporaryServerError{}
	}
	if e := dbTX.Create(&invoice).Error; e != nil {
		log.Printf("[SubscribeToTokenizedAssetByFiat] saving invoice %v: %v", input.ID, e)
		return invoice, &tErrors.ErrorTemporaryServerError{}
	}
	if e := dbTX.Commit().Error; e != nil {
		return invoice, &tErrors.ErrorTemporaryServerError{}
	}
	return invoice, nil
}

// DeliverFiatAssetPurchase carries out a fiat purchase whose payment was
// confirmed: the internal balance token's minting Safe mints the amount to
// itself, approves the offer book and fills the asset's sale offer with
// the buyer's wallet as recipient, in one Safe transaction its owners'
// keys sign (the first pays the gas). It returns the mined transaction.
func DeliverFiatAssetPurchase(invoice *userModels.FiatPaymentInvoice, gc *sharedconfig.GlobalConfig) (string, error) {
	if invoice.TokenizedAssetID == nil || invoice.WalletAddress == nil {
		return "", fmt.Errorf("invoice %v has no asset or wallet", invoice.ID)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	ta, _, err := GetTokenizedAssetByID(*invoice.TokenizedAssetID, gc.DB)
	if err != nil {
		return "", err
	}
	code, token, minter, signers, err := internalBalanceMinter(ctx, &ta, gc)
	if err != nil {
		return "", err
	}
	o, err := loadAssetOffer(ctx, &ta, gc)
	if err != nil {
		return "", err
	}
	purchase, err := priceAssetPurchase(ctx, o, code, token, decimal.NewFromFloat(invoice.Amount), minter, common.HexToAddress(*invoice.WalletAddress), 10*time.Minute, gc)
	if err != nil {
		return "", err
	}
	fill, err := aa.FillCall(o.Book, purchase.Fill, purchase.Authorization)
	if err != nil {
		return "", err
	}
	calls := []gnosissafe.Call{}
	for _, c := range []aa.Call{aa.ERC20Mint(token, minter, purchase.Payment), aa.ERC20Approve(token, o.Book, purchase.Payment), fill} {
		calls = append(calls, gnosissafe.Call{To: c.To, Value: big.NewInt(0), Data: c.Data})
	}
	hash, err := gnosissafe.ExecCalls(ctx, gc.BantuExpansionClient, network.GetBlockchainChainID(), minter, signers, calls, common.HexToAddress(envOrDefault("SAFE_MULTISEND_CALL_ONLY_ADDRESS", gnosissafe.DefaultMultiSendCallOnlyAddress)))
	if err != nil {
		return "", err
	}
	if err := gnosissafe.WaitSuccess(ctx, gc.BantuExpansionClient, common.HexToHash(hash)); err != nil {
		return hash, err
	}
	if e := network.SetWalletAssetAuthorization(*invoice.WalletAddress, basetxn.CreditAsset{Code: *ta.AssetCode, Issuer: o.Token.Hex()}, true, *ta.IssuingWalletAddress, "tokenized-asset fiat purchase"); e != nil {
		log.Printf("[DeliverFiatAssetPurchase] authorizing %v for %v: %v", *invoice.WalletAddress, *ta.AssetCode, e)
	}
	return hash, nil
}

func envOrDefault(key, def string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return def
}
