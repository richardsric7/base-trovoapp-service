package users

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"
	"time"

	"trovo-wallet-api/internal/aa"
	"trovo-wallet-api/internal/basetxn"
	tPayErrors "trovo-wallet-api/internal/components/payments/errors"
	paymentModels "trovo-wallet-api/internal/components/payments/models"
	usersDB "trovo-wallet-api/internal/components/users/db"
	userModels "trovo-wallet-api/internal/components/users/models"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/evmkeypair"
	"trovo-wallet-api/internal/network"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm/clause"
)

// Wallet operation kinds.
const (
	OperationPayment = "PAYMENT"
)

// sharedWalletOperationValidity is how long a shared wallet's operation
// stays signable while it waits for its approvers.
func sharedWalletOperationValidity() time.Duration {
	if d, err := time.ParseDuration(strings.TrimSpace(os.Getenv("SHARED_WALLET_OPERATION_VALIDITY"))); err == nil && d > 0 {
		return d
	}
	return 24 * time.Hour
}

// paymentContext is what a payment operation needs when it is submitted.
type paymentContext struct {
	Fees []sharedconfig.FeeCollection `json:"fees"`
}

// baseUnits converts a human amount into token base units.
func baseUnits(amount string, decimals uint8) (*big.Int, error) {
	d, err := decimal.NewFromString(amount)
	if err != nil || d.IsNegative() {
		return nil, &tPayErrors.ErrorInvalidPaymentAmount{}
	}
	return d.Shift(int32(decimals)).Truncate(0).BigInt(), nil
}

func transferCall(asset basetxn.Asset, to string, amount *big.Int) aa.Call {
	if asset.IsNative() {
		return aa.NativeTransfer(common.HexToAddress(to), amount)
	}
	return aa.ERC20Transfer(common.HexToAddress(asset.GetIssuer()), common.HexToAddress(to), amount)
}

// feeWalletAddress returns the address of a fee wallet setting. The setting
// should hold the address; a private key (the older form) is still
// accepted, and only its address is used.
func feeWalletAddress(secret, label string, gc *sharedconfig.GlobalConfig) (string, error) {
	if s := strings.TrimSpace(secret); common.IsHexAddress(s) {
		return common.HexToAddress(s).Hex(), nil
	}
	kp, err := evmkeypair.ParseFull(secret)
	if err != nil {
		log.Printf("[Pay] %v is not configured: %v", label, err)
		gc.LogDiscordFailedRequest(fmt.Sprintf("[Pay] %v is not configured", label))
		return "", &tErrors.CustomError{Err: "error-fee-wallet-invalid", Param: "feeAmount", ErrMessage: fmt.Sprintf("The %v is not configured. Please try again later.", label)}
	}
	return kp.Address(), nil
}

// authorizeFeeWallet grants a fee wallet the (database) authorization to
// hold a regulated asset.
func authorizeFeeWallet(addr string, asset basetxn.Asset, label string) error {
	if asset.IsNative() || !network.RequiresWalletAuthorization(asset.GetCode()) {
		return nil
	}
	if network.IsWalletAuthorizedForAsset(addr, asset) {
		return nil
	}
	if err := network.SetWalletAssetAuthorization(addr, asset, true, asset.GetIssuer(), "auto-authorized "+label); err != nil {
		log.Printf("[Pay] authorizing %v for %v: %v", label, asset.GetCode(), err)
		return &tErrors.CustomError{Err: "error-could-not-authorize-fee-wallet", Param: "feeAmount", ErrMessage: "Could not authorize the fee wallet to hold this asset."}
	}
	return nil
}

// Pay sends a payment from sourceWallet. It runs in two steps:
//
//  1. without transactionSignature: validates the payment, builds it as a
//     Safe operation (payment, service fee and VAT transfers, plus the
//     wallet's activation on its first send) and returns it as
//     paymentInfo.Transaction for the user to sign;
//  2. with transaction + transactionSignature: submits that operation.
//
// Shared wallets with approvers instead get an approval request on commit.
func Pay(signerUser *userModels.User, sourceWallet *userModels.UserWallet, paymentInfo *paymentModels.PaymentInfo, gc *sharedconfig.GlobalConfig) (*paymentModels.PaymentInfo, *userModels.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	if userModels.IsInternalBalanceAsset(paymentInfo.AssetCode, paymentInfo.ContractAddress, gc) {
		return paymentInfo, nil, &tErrors.CustomError{Param: "assetCode", Err: "error-asset-not-sendable", ErrMessage: "This asset cannot be sent directly."}
	}
	if err := checkRequiredMemo(paymentInfo); err != nil {
		return paymentInfo, nil, err
	}

	walletHasViewOnlyAccess := sourceWallet.HasViewOnlyAccess(gc)
	sourceWalletOwner, _ := sourceWallet.GetWalletOwner(gc.DB, gc)
	if !walletHasViewOnlyAccess {
		paymentInfo.Multiparty = 1
	}
	if err := computePaymentFees(sourceWallet, &sourceWalletOwner, paymentInfo, gc); err != nil {
		return paymentInfo, nil, err
	}
	paymentInfo.NetworkPassPhrase = network.GetBlockchainNetworkPassPhrase()

	// step 2: the user signed the operation built in step 1
	if paymentInfo.Multiparty == 0 && len(paymentInfo.TransactionSignature) > 0 && len(paymentInfo.Transaction) > 0 {
		return submitPayment(ctx, signerUser, sourceWallet, paymentInfo, gc)
	}
	// a shared wallet's initiator commits the operation built in step 1
	if paymentInfo.Multiparty == 1 && paymentInfo.Commit == 1 && len(paymentInfo.Transaction) > 0 {
		if _, _, err := LoadWalletOperation(paymentInfo.Transaction, sourceWallet.ID, OperationPayment, gc); err != nil {
			return paymentInfo, nil, err
		}
		return paymentInfo, nil, createPaymentApprovalRequest(signerUser, sourceWallet, paymentInfo, gc)
	}

	// step 1: build
	calls, destinationUser, fees, err := buildPaymentCalls(ctx, signerUser, sourceWallet, &sourceWalletOwner, paymentInfo, gc)
	if err != nil {
		log.Printf("[Pay] from [%v] to [%v]: %v", sourceWallet.ID, paymentInfo.Destination, err)
		return paymentInfo, nil, err
	}
	validity := time.Duration(0)
	if paymentInfo.Multiparty == 1 {
		validity = sharedWalletOperationValidity()
	}
	op, err := PrepareWalletOperation(ctx, OperationPayment, signerUser, &sourceWalletOwner, sourceWallet, calls, validity, paymentContext{Fees: fees}, gc)
	if err != nil {
		return paymentInfo, nil, err
	}
	if err := checkPaymentLeavesGas(ctx, sourceWallet, paymentInfo, op.Prepared, gc); err != nil {
		return paymentInfo, nil, err
	}
	paymentInfo.Transaction = op.Transaction
	paymentInfo.Messages = append(paymentInfo.Messages, op.Messages()...)
	paymentInfo.SignatureRequired = 1
	if paymentInfo.Multiparty == 1 && paymentInfo.Commit == 1 {
		return paymentInfo, destinationUser, createPaymentApprovalRequest(signerUser, sourceWallet, paymentInfo, gc)
	}
	return paymentInfo, destinationUser, nil
}

// submitPayment submits a signed payment and records its fees.
func submitPayment(ctx context.Context, signerUser *userModels.User, sourceWallet *userModels.UserWallet, paymentInfo *paymentModels.PaymentInfo, gc *sharedconfig.GlobalConfig) (*paymentModels.PaymentInfo, *userModels.User, error) {
	rec, p, err := LoadWalletOperation(paymentInfo.Transaction, sourceWallet.ID, OperationPayment, gc)
	if err != nil {
		return paymentInfo, nil, err
	}
	hash, err := SignSingleOwnerOperation(ctx, rec, p, signerUser.PrimarySigner, paymentInfo.TransactionSignature, gc)
	if err != nil {
		return paymentInfo, nil, err
	}
	paymentInfo.TransactionID = hash
	recordOperationFees(rec, hash, gc)

	sourceWallet.InvalidateUserCache(gc)
	signerUser.InvalidateUserWalletCache(gc)
	signerUser.InvalidateUserCache(gc)
	var destinationUser *userModels.User
	if len(paymentInfo.Destination) != 42 {
		if du, e := usersDB.GetUser(paymentInfo.Destination, gc.DB, gc); e == nil && len(du.Username) > 0 {
			destinationUser = &du
			du.InvalidateUserCache(gc)
			du.InvalidateUserWalletCache(gc)
		}
	}
	return paymentInfo, destinationUser, nil
}

// recordOperationFees writes the fee records stored with an operation once
// it has been submitted (keyed by its userOpHash).
func recordOperationFees(rec *userModels.WalletOperation, userOpHash string, gc *sharedconfig.GlobalConfig) {
	if rec.Context == nil {
		return
	}
	var c paymentContext
	if err := json.Unmarshal([]byte(*rec.Context), &c); err != nil {
		return
	}
	for _, f := range c.Fees {
		f.ID = gc.GenerateUUIDString()
		h := userOpHash
		f.TransactionHash = &h
		if err := gc.DB.Omit(clause.Associations).Create(&f).Error; err != nil {
			log.Printf("[recordOperationFees] saving %v fee for %v: %v", f.FeeType, userOpHash, err)
			gc.LogDiscordFailedRequest(fmt.Sprintf("[recordOperationFees] saving %v fee for %v failed: %v", f.FeeType, userOpHash, err))
		}
	}
}

// createPaymentApprovalRequest records a shared wallet's payment for its
// approvers; they sign the same operation (see approvals).
func createPaymentApprovalRequest(signerUser *userModels.User, sourceWallet *userModels.UserWallet, paymentInfo *paymentModels.PaymentInfo, gc *sharedconfig.GlobalConfig) error {
	paymentInfo.TransactionID = "PENDING_AUTH"
	assetOfPayment := os.Getenv("NATIVE_ASSET_CODE")
	if len(paymentInfo.ContractAddress) == 42 {
		assetOfPayment = fmt.Sprintf("%v (%v)", paymentInfo.AssetCode, paymentInfo.ContractAddress)
	}
	description := fmt.Sprintf("Payment \nFrom: %v, \nTo: %v, \nAmount: %v %v", sourceWallet.Alias, paymentInfo.Destination, paymentInfo.Amount, assetOfPayment)
	if len(paymentInfo.Memo) > 0 {
		description = fmt.Sprintf("%v \nFor: %v", description, paymentInfo.Memo)
	}
	if len(paymentInfo.Messages) > 0 {
		description = fmt.Sprintf("%v \nMessages: %v", description, strings.Join(paymentInfo.Messages, "\n"))
	}
	infoBytes, _ := json.Marshal(*paymentInfo)
	info := string(infoBytes)
	pendingAuth := userModels.PendingAuth{
		ID:                     uuid.NewString(),
		Initiator:              signerUser.Username,
		InitiatorSignerAddress: signerUser.PrimarySigner,
		WalletAddress:          sourceWallet.ID,
		TransactionType:        "PAYMENT",
		Description:            description,
		TransactionSource:      sourceWallet.ID,
		ApprovalsNeeded:        sourceWallet.NumberOfApprovalsNeeded,
		TransactionXdr:         paymentInfo.Transaction,
		TransactionInfoStr:     &info,
	}
	if err := gc.DB.Omit(clause.Associations).Create(&pendingAuth).Error; err != nil {
		log.Printf("[Pay] saving approval request for %v: %v", sourceWallet.ID, err)
		return &tErrors.ErrorTemporaryServerError{}
	}
	return nil
}

// checkRequiredMemo enforces the memo formats exchanges require for
// deposits to their shared addresses.
func checkRequiredMemo(paymentInfo *paymentModels.PaymentInfo) error {
	if len(paymentInfo.Destination) != 42 {
		paymentInfo.Destination = strings.ToLower(paymentInfo.Destination)
		return nil
	}
	paymentInfo.Destination = strings.ToUpper(paymentInfo.Destination)
	memo := strings.ReplaceAll(paymentInfo.Memo, " ", "")
	rules := []struct {
		env string
		ok  func(int) bool
		msg string
	}{
		{"WALLETS_REQUIRE_28_BYTE_MEMO", func(n int) bool { return n == 28 }, "This exchange requires a 28 character memo. If you do not put the exact memo, your funds will be lost."},
		{"WALLETS_REQUIRE_16_BYTE_MEMO", func(n int) bool { return n == 16 }, "This exchange requires a 16 character memo. If you do not put the exact memo, your funds will be lost."},
		{"WALLETS_REQUIRE_VARIABLE_BYTE_MEMO", func(n int) bool { return n >= 9 }, "You are sending to an exchange that requires a memo for all deposits. If you do not put the exact memo, your funds will be lost."},
	}
	for _, r := range rules {
		for _, w := range strings.Split(os.Getenv(r.env), ",") {
			if w != "" && strings.EqualFold(paymentInfo.Destination, w) && !r.ok(len(memo)) {
				return &tErrors.CustomError{Param: "memo", Err: "error invalid memo", ErrMessage: r.msg}
			}
		}
	}
	return nil
}

// computePaymentFees sets the service fee and VAT on paymentInfo (enterprise
// and shared-wallet payments pay a percentage fee).
func computePaymentFees(sourceWallet *userModels.UserWallet, owner *userModels.User, paymentInfo *paymentModels.PaymentInfo, gc *sharedconfig.GlobalConfig) error {
	paymentInfo.AmountToPay = paymentInfo.Amount
	paymentInfo.FeeAmount = "0"
	var feePercent float64
	if owner.BelongsToAnEnterpriseProfile() {
		slf, exists, e := gc.GetServiceLinkFees(*owner.CreatedByServiceLinkID)
		if e != nil {
			return e
		}
		if exists {
			feePercent = float64(slf.PaymentFee)
		} else {
			feePercent = sourceWallet.GetSharedAccessPaymentFee(gc).FeePercent
		}
	} else if paymentInfo.Multiparty == 1 {
		feePercent = sourceWallet.GetSharedAccessPaymentFee(gc).FeePercent
	}
	amount, err := decimal.NewFromString(paymentInfo.Amount)
	if err != nil || !amount.IsPositive() {
		return &tPayErrors.ErrorInvalidPaymentAmount{}
	}
	fee := decimal.NewFromFloat(feePercent)
	feeAmount := amount.Mul(fee).Div(decimal.NewFromInt(100)).Truncate(7)
	vat := decimal.NewFromFloat(gc.GetVATValue(feeAmount))
	paymentInfo.Fee = fee.String()
	paymentInfo.FeeAmount = feeAmount.String()
	paymentInfo.Vat = decimal.NewFromFloat(gc.GetVATRate()).String()
	paymentInfo.VatAmount = vat.String()
	paymentInfo.AmountToPay = amount.Add(feeAmount).Add(vat).String()
	return nil
}

// buildPaymentCalls validates the payment and returns the calls the source
// wallet makes: the payment, then the service fee and VAT transfers.
func buildPaymentCalls(ctx context.Context, owner *userModels.User, sourceWallet *userModels.UserWallet, sourceWalletOwner *userModels.User, paymentInfo *paymentModels.PaymentInfo, gc *sharedconfig.GlobalConfig) ([]aa.Call, *userModels.User, []sharedconfig.FeeCollection, error) {
	client := gc.BantuExpansionClient
	publicKeyPayment := len(paymentInfo.Destination) == 42

	var err error
	if paymentInfo, err = ValidatePaymentInfo(paymentInfo); err != nil {
		return nil, nil, nil, err
	}

	var asset basetxn.Asset = basetxn.NativeAsset{}
	if len(paymentInfo.ContractAddress) > 0 {
		asset = basetxn.CreditAsset{Code: paymentInfo.AssetCode, Issuer: paymentInfo.ContractAddress}
	}

	destinationInfo, getDestinationError := usersDB.GetUser(paymentInfo.Destination, gc.DB, gc)
	if !publicKeyPayment && (getDestinationError != nil || len(destinationInfo.Username) == 0) {
		return nil, nil, nil, &tPayErrors.ErrorPaymentDestinationDoesNotExist{}
	}
	if !publicKeyPayment && !strings.Contains(paymentInfo.Destination, destinationInfo.Username) {
		paymentInfo.Messages = append(paymentInfo.Messages, fmt.Sprintf("Notice: [%v] belongs to the wallet alias [%v] and will be used as the destination.", paymentInfo.Destination, destinationInfo.Username))
		paymentInfo.Destination = destinationInfo.Username
	}
	destinationWallet, _ := usersDB.GetWallet(paymentInfo.Destination, gc.DB)
	if len(destinationWallet.ID) == 0 && !publicKeyPayment {
		return nil, nil, nil, &tPayErrors.ErrorPaymentDestinationDoesNotExist{}
	}
	if destinationInfo.Suspended == 1 {
		return nil, nil, nil, &tErrors.ErrorUsernameIsSuspended{}
	}

	destinationAddress := destinationWallet.ID
	if publicKeyPayment {
		paymentInfo.Destination = strings.ToUpper(paymentInfo.Destination)
		if _, err := evmkeypair.ParseAddress(paymentInfo.Destination); err != nil {
			return nil, nil, nil, &tPayErrors.ErrorInvalidPaymentDestinationAddress{}
		}
		destinationAddress = paymentInfo.Destination
		paymentInfo.Messages = append(paymentInfo.Messages, "You are about to make payment to an address directly. Please be sure of the address as the payment cannot be retrieved after confirmation.")
	} else {
		if gc.IsValidTokenizedAsset(asset.GetCode()) && destinationInfo.KYCVerified == 0 && asset.GetIssuer() != destinationInfo.Address {
			return nil, nil, nil, &tErrors.CustomError{Param: "destination", Err: "error-invalid-kyc", ErrMessage: fmt.Sprintf("%v has not passed/met the KYC requirement to receive this tokenized asset %v.", destinationInfo.Username, asset.GetCode())}
		}
		paymentInfo.DestinationFirstName = destinationInfo.FirstName
		if destinationInfo.LastName != nil {
			paymentInfo.DestinationLastName = *destinationInfo.LastName
		}
		paymentInfo.DestinationVerified = destinationInfo.Verified
		if destinationInfo.ImageThumbnailURL != nil {
			paymentInfo.DestinationThumbnail = *destinationInfo.ImageThumbnailURL
		}
		// non-standard wallets are only funded by sibling wallets
		if destinationWallet.WalletType != 0 && !strings.EqualFold(destinationWallet.Signer, sourceWallet.Signer) &&
			(asset.IsNative() || !destinationWallet.OwnerOfBlockchainAsset(asset.GetCode())) {
			return nil, nil, nil, &tErrors.CustomError{Param: "destination", Err: "error-action-forbidden", ErrMessage: fmt.Sprintf("%v is only allowed to be funded by a sibling wallet. Only wallets belonging to the same account can fund a non-standard wallet.", destinationWallet.Alias)}
		}
	}

	// regulated assets: the destination must be authorized to hold them
	_, destinationAuthorized, _, _, _, destinationErr := network.BlockchainAccountProperties(client, destinationAddress, asset)
	if destinationErr != nil {
		return nil, nil, nil, destinationErr
	}
	if !asset.IsNative() && !destinationAuthorized {
		return nil, nil, nil, &tErrors.CustomError{Param: "destination", Err: "error-destination-cannot-accept-asset", ErrMessage: fmt.Sprintf("%v does not accept the asset %v at this time.", paymentInfo.Destination, asset.GetCode())}
	}

	if !asset.IsNative() && strings.EqualFold(sourceWallet.ID, asset.GetIssuer()) {
		return nil, nil, nil, &tErrors.CustomError{Param: "destination", Err: "error-source-forbidden-to-sending-asset", ErrMessage: fmt.Sprintf("%v, a token minting wallet, is forbidden from sending %v.", sourceWallet.Alias, asset.GetCode())}
	}
	_, sourceAuthorized, nativeBalance, tokenBalance, _, sourceErr := network.BlockchainAccountProperties(client, sourceWallet.ID, asset)
	if sourceErr != nil {
		return nil, nil, nil, sourceErr
	}
	if !sourceAuthorized {
		return nil, nil, nil, &tErrors.ErrorUnderfundedAccount{}
	}
	total, _ := decimal.NewFromString(paymentInfo.AmountToPay)
	balance := tokenBalance
	if asset.IsNative() {
		balance = nativeBalance
	}
	if balance.LessThan(total) {
		return nil, nil, nil, &tErrors.ErrorUnderfundedAccount{}
	}

	decimals, err := network.AssetDecimals(ctx, client, asset)
	if err != nil {
		log.Printf("[buildPaymentCalls] decimals of %v: %v", asset.GetIssuer(), err)
		return nil, nil, nil, &tErrors.ErrorTemporaryServerError{}
	}
	amount, err := baseUnits(paymentInfo.Amount, decimals)
	if err != nil {
		return nil, nil, nil, err
	}
	calls := []aa.Call{transferCall(asset, destinationAddress, amount)}

	var fees []sharedconfig.FeeCollection
	feeAmount, _ := decimal.NewFromString(paymentInfo.FeeAmount)
	serviceFee := sourceWallet.GetSharedAccessPaymentFee(gc)
	if feeAmount.IsPositive() && serviceFee.Inactive == 0 {
		var feeAddr string
		if sourceWalletOwner.BelongsToAnEnterpriseProfile() {
			feeAddr, err = feeWalletAddress(sourceWallet.GetPaymentFeeWallet(gc), "payment fee wallet", gc)
		} else if paymentInfo.Multiparty == 1 {
			feeAddr, err = feeWalletAddress(serviceFee.FeeWalletSecretKey, "shared access fee wallet", gc)
		}
		if err != nil {
			return nil, nil, nil, err
		}
		if feeAddr != "" {
			if err := authorizeFeeWallet(feeAddr, asset, "payment fee wallet"); err != nil {
				return nil, nil, nil, err
			}
			units, err := baseUnits(paymentInfo.FeeAmount, decimals)
			if err != nil {
				return nil, nil, nil, err
			}
			calls = append(calls, transferCall(asset, feeAddr, units))
			paymentInfo.Messages = append(paymentInfo.Messages, fmt.Sprintf("%v%% will be added from wallet %v as service fee.", paymentInfo.Fee, sourceWallet.Alias))
			fees = append(fees, feeRecord("PAYMENT", sourceWallet, sourceWalletOwner, asset, paymentInfo.FeeAmount, feeAddr, paymentInfo.Multiparty))
		}

		vatAmount, _ := decimal.NewFromString(paymentInfo.VatAmount)
		if vatAmount.IsPositive() {
			vatAddr, err := feeWalletAddress(gc.GetVATWallet(), "VAT wallet", gc)
			if err != nil {
				return nil, nil, nil, err
			}
			if err := authorizeFeeWallet(vatAddr, asset, "VAT wallet"); err != nil {
				return nil, nil, nil, err
			}
			units, err := baseUnits(paymentInfo.VatAmount, decimals)
			if err != nil {
				return nil, nil, nil, err
			}
			calls = append(calls, transferCall(asset, vatAddr, units))
			paymentInfo.Messages = append(paymentInfo.Messages, fmt.Sprintf("%v%% VAT on the fee (%v) will be added from wallet %v.", paymentInfo.Vat, paymentInfo.VatAmount, sourceWallet.Alias))
			fees = append(fees, feeRecord("VAT", sourceWallet, sourceWalletOwner, asset, paymentInfo.VatAmount, vatAddr, paymentInfo.Multiparty))
		}
	}

	if publicKeyPayment {
		return calls, nil, fees, nil
	}
	return calls, &destinationInfo, fees, nil
}

func feeRecord(feeType string, w *userModels.UserWallet, owner *userModels.User, asset basetxn.Asset, amount, destination string, multiparty int) sharedconfig.FeeCollection {
	code := asset.GetCode()
	var contract *string
	if asset.IsNative() {
		code = os.Getenv("NATIVE_ASSET_CODE")
	} else {
		c := asset.GetIssuer()
		contract = &c
	}
	a, _ := decimal.NewFromString(amount)
	return sharedconfig.FeeCollection{
		FromUsername: owner.Username, FromWalletAddress: w.ID, FromWalletAlias: w.Alias,
		BelongsToEnterpriseProfile: owner.CreatedByServiceLinkID, FeeType: feeType,
		Amount: a.InexactFloat64(), AssetCode: code, ContractAddress: contract,
		DestinationWallet: destination, SharedAccessOperation: multiparty,
	}
}

// checkPaymentLeavesGas refuses a payment that would leave the wallet
// unable to pay its own network fee, when the fee is paid in the asset
// being sent.
func checkPaymentLeavesGas(ctx context.Context, w *userModels.UserWallet, paymentInfo *paymentModels.PaymentInfo, p *aa.Prepared, gc *sharedconfig.GlobalConfig) error {
	total, _ := decimal.NewFromString(paymentInfo.AmountToPay)
	fail := func(fee string) error {
		return &tErrors.CustomError{Param: "amount", Err: "error-insufficient-for-network-fee", ErrMessage: fmt.Sprintf("This wallet cannot cover the amount plus the network fee (up to %v). Reduce the amount and try again.", fee), Code: 400}
	}
	if p.Quote != nil {
		if len(paymentInfo.ContractAddress) == 0 || !strings.EqualFold(paymentInfo.ContractAddress, p.Quote.Token.Hex()) {
			return nil
		}
		dec, err := network.AssetDecimals(ctx, gc.BantuExpansionClient, basetxn.CreditAsset{Code: paymentInfo.AssetCode, Issuer: paymentInfo.ContractAddress})
		if err != nil {
			return nil
		}
		bal, err := network.B20BalanceOf(gc.BantuExpansionClient, paymentInfo.ContractAddress, w.ID, dec)
		if err == nil && bal.LessThan(total.Add(p.Quote.MaxTokenCostFormatted)) {
			return fail(fmt.Sprintf("%v %v", p.Quote.MaxTokenCostFormatted, p.Quote.Symbol))
		}
		return nil
	}
	if len(paymentInfo.ContractAddress) > 0 || p.MaxCostWei == nil {
		return nil
	}
	bal, err := gc.BantuExpansionClient.BalanceAt(ctx, common.HexToAddress(w.ID), nil)
	if err != nil {
		return nil
	}
	maxFee := decimal.NewFromBigInt(p.MaxCostWei.ToInt(), -18)
	if decimal.NewFromBigInt(bal, -18).LessThan(total.Add(maxFee)) {
		return fail(fmt.Sprintf("%v ETH", maxFee.Round(8)))
	}
	return nil
}
