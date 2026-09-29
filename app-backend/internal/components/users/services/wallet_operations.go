package users

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"trovo-wallet-api/internal/aa"
	assetModels "trovo-wallet-api/internal/components/assets/models"
	userModels "trovo-wallet-api/internal/components/users/models"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/network"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Every send is a Safe UserOperation (internal/aa): the backend builds it,
// the wallet's owners sign its SafeOp hash with the apps' existing
// signBase64Txn (personal_sign of the base64-decoded "transaction"), and
// the backend submits it to the bundler. The wallet pays its own gas - in
// the stablecoin the user chose (User.GasFeeAsset), through the paymaster,
// or in ETH.

var (
	builderOnce sync.Once
	builder     *aa.Builder
)

// operationBuilder is configured from BUNDLER_URL, PAYMASTER_ADDRESS,
// PAYMASTER_QUOTE_SERVICE_URL, PAYMASTER_QUOTE_SERVICE_API_KEY,
// WALLET_OPERATION_VALIDITY and WALLET_OPERATION_GAS_BUFFER_PERCENT.
func operationBuilder(gc *sharedconfig.GlobalConfig) *aa.Builder {
	builderOnce.Do(func() {
		validity := 10 * time.Minute
		if d, err := time.ParseDuration(strings.TrimSpace(os.Getenv("WALLET_OPERATION_VALIDITY"))); err == nil && d > 0 {
			validity = d
		}
		buffer := int64(20)
		if v, err := strconv.ParseInt(strings.TrimSpace(os.Getenv("WALLET_OPERATION_GAS_BUFFER_PERCENT")), 10, 64); err == nil && v >= 0 {
			buffer = v
		}
		builder = &aa.Builder{
			Config:           network.AAConfig(),
			Chain:            gc.BantuExpansionClient,
			Bundler:          aa.NewBundler(strings.TrimSpace(os.Getenv("BUNDLER_URL"))),
			DefaultValidity:  validity,
			GasBufferPercent: buffer,
		}
		if url := strings.TrimSpace(os.Getenv("PAYMASTER_QUOTE_SERVICE_URL")); url != "" && common.IsHexAddress(os.Getenv("PAYMASTER_ADDRESS")) {
			builder.Quotes = aa.NewQuoteClient(strings.TrimRight(url, "/"), os.Getenv("PAYMASTER_QUOTE_SERVICE_API_KEY"))
			builder.Paymaster = common.HexToAddress(os.Getenv("PAYMASTER_ADDRESS"))
		}
	})
	return builder
}

// walletOwners returns who controls the wallet: its on-chain owners and
// threshold once deployed, otherwise its initial owner.
func walletOwners(ctx context.Context, w *userModels.UserWallet, gc *sharedconfig.GlobalConfig) (aa.Wallet, error) {
	addr := common.HexToAddress(w.ID)
	salt, ok := new(big.Int).SetString(w.SafeSaltNonce, 10)
	if !ok {
		salt = big.NewInt(0)
	}
	aw := aa.Wallet{Address: addr, SaltNonce: salt, InitialThreshold: int64(w.InitialThreshold)}
	if aw.InitialThreshold < 1 {
		aw.InitialThreshold = 1
	}
	for _, o := range w.InitialOwnerList() {
		if common.IsHexAddress(o) {
			aw.InitialOwners = append(aw.InitialOwners, common.HexToAddress(o))
		}
	}
	deployed, err := aa.Deployed(ctx, gc.BantuExpansionClient, addr)
	if err != nil {
		return aw, err
	}
	if deployed {
		owners, threshold, err := aa.OnchainOwners(ctx, gc.BantuExpansionClient, addr)
		if err != nil {
			return aw, err
		}
		aw.Owners, aw.Threshold = owners, threshold
		return aw, nil
	}
	aw.Owners, aw.Threshold = aw.InitialOwners, aw.InitialThreshold
	return aw, nil
}

// gasToken picks the stablecoin to pay gas in: the user's GasFeeAsset if it
// is still gas-fee eligible and the wallet holds some; otherwise nil (ETH).
func gasToken(user *userModels.User, w *userModels.UserWallet, gc *sharedconfig.GlobalConfig) *common.Address {
	if user == nil || user.GasFeeAsset == nil || operationBuilder(gc).Quotes == nil {
		return nil
	}
	var asset assetModels.CuratedAsset
	if err := gc.DB.Where("asset_code = ? AND gas_fee_eligible = ? AND inactive = ?", *user.GasFeeAsset, true, 0).First(&asset).Error; err != nil {
		return nil
	}
	if !common.IsHexAddress(asset.ContractAddress) {
		return nil
	}
	bal, err := network.B20BalanceOf(gc.BantuExpansionClient, asset.ContractAddress, w.ID, uint8(asset.DecimalPlaces))
	if err != nil || !bal.IsPositive() {
		return nil
	}
	token := common.HexToAddress(asset.ContractAddress)
	return &token
}

// PreparedWalletOperation is a built operation waiting for signatures.
type PreparedWalletOperation struct {
	// Transaction is what the app signs: base64 of the SafeOp hash (the
	// apps' signBase64Txn decodes it and personal_signs the bytes).
	Transaction string
	Prepared    *aa.Prepared
	Record      userModels.WalletOperation
}

// operationMessages describes the network fee (and activation) for the
// app to show before the user signs.
func operationMessages(p *aa.Prepared) []string {
	var msgs []string
	if p.Activation {
		msgs = append(msgs, "This is your wallet's first transaction: it also activates the wallet on the network.")
	}
	if p.Quote != nil {
		msgs = append(msgs, fmt.Sprintf("Network fee: up to %v %v, paid from this wallet (any unused part is refunded).", p.Quote.MaxTokenCostFormatted.String(), p.Quote.Symbol))
	} else if p.MaxCostWei != nil {
		msgs = append(msgs, fmt.Sprintf("Network fee: up to %v ETH, paid from this wallet.", decimal.NewFromBigInt(p.MaxCostWei.ToInt(), -18).Round(8).String()))
	}
	return msgs
}

// PrepareWalletOperation builds the Safe operation executing calls from
// wallet, prices its gas and stores it for signing. payer is the user
// whose gas-fee preference applies (the wallet's owner). validity 0 uses
// the default (shared wallets waiting for approvers pass a longer one).
func PrepareWalletOperation(ctx context.Context, kind string, initiator, payer *userModels.User, wallet *userModels.UserWallet, calls []aa.Call, validity time.Duration, opContext interface{}, gc *sharedconfig.GlobalConfig) (*PreparedWalletOperation, error) {
	b := operationBuilder(gc)
	if b.Bundler == nil || os.Getenv("BUNDLER_URL") == "" {
		log.Printf("[PrepareWalletOperation] BUNDLER_URL is not configured")
		return nil, &tErrors.ErrorTemporaryServerError{}
	}
	aw, err := walletOwners(ctx, wallet, gc)
	if err != nil {
		log.Printf("[PrepareWalletOperation] reading wallet %v: %v", wallet.ID, err)
		return nil, &tErrors.ErrorTemporaryServerError{}
	}
	p, err := b.Prepare(ctx, aa.Request{Wallet: aw, Calls: calls, GasToken: gasToken(payer, wallet, gc), Validity: validity})
	if err != nil && errors.Is(err, aa.ErrPaymasterUnavailable) {
		// the stablecoin route is down: fall back to ETH
		log.Printf("[PrepareWalletOperation] %v; falling back to ETH gas for %v", err, wallet.ID)
		p, err = b.Prepare(ctx, aa.Request{Wallet: aw, Calls: calls, Validity: validity})
	}
	if err != nil {
		if errors.Is(err, aa.ErrInsufficientGasFunds) {
			return nil, &tErrors.CustomError{Param: "amount", Err: "error-insufficient-network-fee", ErrMessage: "This wallet does not hold enough to pay the network fee. Add ETH, or a stablecoin you pay network fees in, and try again.", Code: http.StatusBadRequest}
		}
		log.Printf("[PrepareWalletOperation] %v %v: %v", kind, wallet.ID, err)
		return nil, &tErrors.CustomError{Param: "transaction", Err: "error-operation-could-not-be-prepared", ErrMessage: "The transaction could not be prepared. Please check the details and try again.", Code: http.StatusBadRequest}
	}

	raw, err := p.Marshal()
	if err != nil {
		return nil, &tErrors.ErrorTemporaryServerError{}
	}
	rec := userModels.WalletOperation{
		ID:            p.SafeOpHash.Hex(),
		WalletAddress: wallet.ID,
		Kind:          kind,
		Initiator:     initiator.Username,
		Status:        userModels.WalletOperationPending,
		Prepared:      string(raw),
		ExpiresAt:     time.Unix(int64(p.ValidUntil), 0),
		Activation:    p.Activation,
	}
	userOpHash := p.UserOpHash.Hex()
	rec.UserOpHash = &userOpHash
	if p.GasToken != nil {
		rec.GasToken = p.GasToken.Hex()
	}
	if opContext != nil {
		if c, err := json.Marshal(opContext); err == nil {
			s := string(c)
			rec.Context = &s
		}
	}
	if err := gc.DB.Omit(clause.Associations).Create(&rec).Error; err != nil {
		log.Printf("[PrepareWalletOperation] saving %v: %v", rec.ID, err)
		return nil, &tErrors.ErrorTemporaryServerError{}
	}
	return &PreparedWalletOperation{Transaction: base64.StdEncoding.EncodeToString(p.SafeOpHash.Bytes()), Prepared: p, Record: rec}, nil
}

// LoadWalletOperation finds a pending operation by the "transaction" the
// app was given (base64 SafeOp hash), checking it belongs to wallet and is
// of the expected kind.
func LoadWalletOperation(transaction, walletAddress, kind string, gc *sharedconfig.GlobalConfig) (*userModels.WalletOperation, *aa.Prepared, error) {
	hash, err := base64.StdEncoding.DecodeString(strings.TrimSpace(transaction))
	if err != nil || len(hash) != 32 {
		return nil, nil, &tErrors.CustomError{Param: "transaction", Err: "error-invalid-transaction", ErrMessage: "Unknown transaction. Please start again.", Code: http.StatusBadRequest}
	}
	var rec userModels.WalletOperation
	if err := gc.DB.Where("id = ?", common.BytesToHash(hash).Hex()).First(&rec).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, &tErrors.CustomError{Param: "transaction", Err: "error-invalid-transaction", ErrMessage: "Unknown transaction. Please start again.", Code: http.StatusBadRequest}
		}
		return nil, nil, &tErrors.ErrorTemporaryServerError{}
	}
	if !strings.EqualFold(rec.WalletAddress, walletAddress) || rec.Kind != kind {
		return nil, nil, &tErrors.CustomError{Param: "transaction", Err: "error-transaction-mismatch", ErrMessage: "This transaction does not belong to this wallet or request.", Code: http.StatusBadRequest}
	}
	p, err := aa.UnmarshalPrepared([]byte(rec.Prepared))
	if err != nil {
		return nil, nil, &tErrors.ErrorTemporaryServerError{}
	}
	return &rec, p, nil
}

// DecodeAppSignature decodes the app's base64 personal_sign signature.
func DecodeAppSignature(sigB64 string) ([]byte, error) {
	sig, err := base64.StdEncoding.DecodeString(strings.TrimSpace(sigB64))
	if err != nil || len(sig) != 65 {
		return nil, &tErrors.CustomError{Param: "transactionSignature", Err: "error-invalid-signature", ErrMessage: "The transaction signature is invalid.", Code: http.StatusBadRequest}
	}
	return sig, nil
}

// SubmitWalletOperation assembles the owners' signatures into the stored
// operation and sends it to the bundler, returning the userOpHash.
func SubmitWalletOperation(ctx context.Context, rec *userModels.WalletOperation, p *aa.Prepared, sigs []aa.OwnerSignature, gc *sharedconfig.GlobalConfig) (string, error) {
	if rec.Status != userModels.WalletOperationPending {
		return "", &tErrors.CustomError{Param: "transaction", Err: "error-transaction-already-submitted", ErrMessage: "This transaction has already been submitted.", Code: http.StatusConflict}
	}
	// claim it first so a double submit cannot send twice
	res := gc.DB.Model(&userModels.WalletOperation{}).Where("id = ? AND status = ?", rec.ID, userModels.WalletOperationPending).Update("status", userModels.WalletOperationSubmitted)
	if res.Error != nil {
		return "", &tErrors.ErrorTemporaryServerError{}
	}
	if res.RowsAffected == 0 {
		return "", &tErrors.CustomError{Param: "transaction", Err: "error-transaction-already-submitted", ErrMessage: "This transaction has already been submitted.", Code: http.StatusConflict}
	}
	hash, err := operationBuilder(gc).Submit(ctx, p, sigs)
	if err != nil {
		msg := err.Error()
		status := userModels.WalletOperationFailed
		if errors.Is(err, aa.ErrNotEnoughSignatures) || strings.Contains(msg, "does not sign this operation") || strings.Contains(msg, "not an owner") {
			// a bad signature: let the signer retry
			status = userModels.WalletOperationPending
		}
		gc.DB.Model(&userModels.WalletOperation{}).Where("id = ?", rec.ID).Updates(map[string]interface{}{"status": status, "error": msg})
		log.Printf("[SubmitWalletOperation] %v %v: %v", rec.Kind, rec.ID, err)
		if status == userModels.WalletOperationPending {
			return "", &tErrors.CustomError{Param: "transactionSignature", Err: "error-invalid-signature", ErrMessage: "The signature does not authorize this transaction.", Code: http.StatusBadRequest}
		}
		if strings.Contains(msg, "expired") {
			return "", &tErrors.CustomError{Param: "transaction", Err: "error-transaction-expired", ErrMessage: "This transaction expired before it was signed. Please start again.", Code: http.StatusBadRequest}
		}
		return "", &tErrors.CustomError{Param: "transaction", Err: "error-transaction-failed", ErrMessage: "The network did not accept this transaction. Please try again.", Code: http.StatusBadRequest}
	}
	h := hash.Hex()
	gc.DB.Model(&userModels.WalletOperation{}).Where("id = ?", rec.ID).Update("user_op_hash", h)
	rec.Status, rec.UserOpHash = userModels.WalletOperationSubmitted, &h
	return h, nil
}

// SignSingleOwnerOperation is the common case: the wallet's one owner (the
// signing user) has signed.
func SignSingleOwnerOperation(ctx context.Context, rec *userModels.WalletOperation, p *aa.Prepared, signer string, sigB64 string, gc *sharedconfig.GlobalConfig) (string, error) {
	sig, err := DecodeAppSignature(sigB64)
	if err != nil {
		return "", err
	}
	return SubmitWalletOperation(ctx, rec, p, []aa.OwnerSignature{{Owner: common.HexToAddress(signer), Signature: sig}}, gc)
}

// TrackWalletOperations follows submitted operations to inclusion: it
// records the transaction hash and outcome and marks wallets activated by
// their first operation. It runs in the background (see main.go).
func TrackWalletOperations(ctx context.Context, gc *sharedconfig.GlobalConfig) {
	b := operationBuilder(gc)
	bun, ok := b.Bundler.(*aa.Bundler)
	if !ok || bun.URL == "" {
		return
	}
	var ops []userModels.WalletOperation
	if err := gc.DB.Where("status = ? AND user_op_hash IS NOT NULL", userModels.WalletOperationSubmitted).Limit(200).Find(&ops).Error; err != nil {
		return
	}
	for _, op := range ops {
		r, err := bun.Receipt(ctx, common.HexToHash(*op.UserOpHash))
		if err != nil || r == nil {
			if time.Since(op.UpdatedAt) > 30*time.Minute {
				msg := "not included within 30 minutes"
				gc.DB.Model(&userModels.WalletOperation{}).Where("id = ?", op.ID).Updates(map[string]interface{}{"status": userModels.WalletOperationFailed, "error": msg})
			}
			continue
		}
		applyReceipt(op, r, gc)
	}
}

// applyReceipt records a mined operation's outcome and marks the wallets it
// deployed activated.
func applyReceipt(op userModels.WalletOperation, r *aa.Receipt, gc *sharedconfig.GlobalConfig) {
	tx := r.Receipt.TransactionHash.Hex()
	success := r.Success
	updates := map[string]interface{}{"status": userModels.WalletOperationIncluded, "tx_hash": tx, "success": success}
	if !success {
		updates["error"] = "the operation's calls reverted: " + r.Reason
	}
	gc.DB.Model(&userModels.WalletOperation{}).Where("id = ?", op.ID).Updates(updates)
	// the wallet was deployed by this operation, even if its calls reverted
	now := time.Now()
	if op.Activation {
		markActivated(op.WalletAddress, op.WalletAddress, tx, now, gc)
	}
	// wallets it deploys exist only if its calls succeeded
	if success {
		for _, w := range strings.Split(op.Deploys, ",") {
			if w = strings.TrimSpace(w); w != "" {
				markActivated(w, op.WalletAddress, tx, now, gc)
			}
		}
	}
}

// WalletOperationOutcome reports what became of a submitted operation,
// identified by its userOpHash, waiting for it until ctx ends. done is
// false while it is not yet included; once done, success says whether its
// calls succeeded and txHash is the transaction that included it.
func WalletOperationOutcome(ctx context.Context, userOpHash string, gc *sharedconfig.GlobalConfig) (txHash string, success, done bool, err error) {
	var op userModels.WalletOperation
	if err := gc.DB.Where("user_op_hash = ?", userOpHash).First(&op).Error; err != nil {
		return "", false, false, err
	}
	switch op.Status {
	case userModels.WalletOperationIncluded:
		if op.TxHash != nil {
			txHash = *op.TxHash
		}
		return txHash, op.Success != nil && *op.Success, true, nil
	case userModels.WalletOperationFailed, userModels.WalletOperationExpired:
		return "", false, true, nil
	case userModels.WalletOperationPending:
		return "", false, false, nil
	}
	bun, ok := operationBuilder(gc).Bundler.(*aa.Bundler)
	if !ok || bun.URL == "" {
		return "", false, false, nil
	}
	r, err := bun.WaitReceipt(ctx, common.HexToHash(userOpHash))
	if err != nil || r == nil {
		return "", false, false, nil // not included yet
	}
	applyReceipt(op, r, gc)
	return r.Receipt.TransactionHash.Hex(), r.Success, true, nil
}

func markActivated(wallet, by, tx string, at time.Time, gc *sharedconfig.GlobalConfig) {
	gc.DB.Model(&userModels.UserWallet{}).Where("UPPER(id) = ? AND activated = ?", strings.ToUpper(wallet), false).Updates(map[string]interface{}{
		"activated": true, "activated_at": &at, "activation_tx_hash": tx, "activated_by": by,
	})
}

// ExpireWalletOperations marks pending operations past their validity.
func ExpireWalletOperations(gc *sharedconfig.GlobalConfig) {
	gc.DB.Model(&userModels.WalletOperation{}).Where("status = ? AND expires_at < ?", userModels.WalletOperationPending, time.Now()).Update("status", userModels.WalletOperationExpired)
}
