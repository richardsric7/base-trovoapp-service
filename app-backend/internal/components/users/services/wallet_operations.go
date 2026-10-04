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
	"trovo-wallet-api/internal/basetxn"
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
	for _, m := range w.InitialModuleList() {
		if common.IsHexAddress(m) {
			aw.InitialModules = append(aw.InitialModules, common.HexToAddress(m))
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
	// Notes explain choices made while building it (e.g. unpaid gas).
	Notes []string
}

// Messages are what the app shows before the user signs: the notes, then
// the network fee (and activation).
func (op *PreparedWalletOperation) Messages() []string {
	return append(append([]string{}, op.Notes...), operationMessages(op.Prepared)...)
}

// gasDebt is what a wallet owes the paymaster in one token.
type gasDebt struct {
	amount  decimal.Decimal
	code    string
	payable bool // the wallet holds enough to settle it
	calls   []aa.Call
}

func (d gasDebt) note(settled bool) string {
	if settled {
		return fmt.Sprintf("This also pays %v %v of network fees this wallet still owed from an earlier transaction; this transaction's fee is paid in ETH.", d.amount, d.code)
	}
	return fmt.Sprintf("This wallet still owes %v %v of network fees from an earlier transaction, so its network fees are paid in ETH until that is paid. Hold at least that much %v in this wallet to settle it with your next transaction.", d.amount, d.code, d.code)
}

// gasDebts reads the wallet's gas debts with the paymaster, in each
// gas-eligible stablecoin, with the calls that settle them. While a wallet
// owes any, the paymaster will not pay its gas, so it pays in ETH.
func gasDebts(ctx context.Context, wallet common.Address, gc *sharedconfig.GlobalConfig) []gasDebt {
	b := operationBuilder(gc)
	if b.Paymaster == (common.Address{}) {
		return nil
	}
	var assets []assetModels.CuratedAsset
	gc.DB.Where("gas_fee_eligible = ? AND inactive = ?", true, 0).Find(&assets)
	var out []gasDebt
	for _, a := range assets {
		if !common.IsHexAddress(a.ContractAddress) {
			continue
		}
		token := common.HexToAddress(a.ContractAddress)
		owed, err := aa.PaymasterDebt(ctx, gc.BantuExpansionClient, b.Paymaster, wallet, token)
		if err != nil || owed.Sign() == 0 {
			continue
		}
		decimals, err := network.AssetDecimals(ctx, gc.BantuExpansionClient, basetxn.CreditAsset{Code: a.AssetCode, Issuer: a.ContractAddress})
		if err != nil {
			continue
		}
		d := gasDebt{amount: decimal.NewFromBigInt(owed, -int32(decimals)), code: a.AssetCode}
		if bal, err := network.B20BalanceOf(gc.BantuExpansionClient, a.ContractAddress, wallet.Hex(), decimals); err == nil && !bal.LessThan(d.amount) {
			d.payable, d.calls = true, aa.SettleDebtCalls(b.Paymaster, wallet, token)
		}
		out = append(out, d)
	}
	return out
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
	// an operation that waits for approvers (validity > 0) gets its own
	// nonce sequence, so it neither blocks nor is invalidated by the
	// wallet's other operations
	var nonceKey *big.Int
	if validity > 0 {
		if nonceKey, err = aa.RandomNonceKey(); err != nil {
			return nil, &tErrors.ErrorTemporaryServerError{}
		}
	}
	var notes []string
	p, err := b.Prepare(ctx, aa.Request{Wallet: aw, Calls: calls, GasToken: gasToken(payer, wallet, gc), Validity: validity, NonceKey: nonceKey})
	if err != nil && errors.Is(err, aa.ErrPaymasterUnavailable) {
		// the stablecoin route is unavailable: pay in ETH
		log.Printf("[PrepareWalletOperation] %v; falling back to ETH gas for %v", err, wallet.ID)
		p = nil
		if errors.Is(err, aa.ErrOutstandingDebt) {
			// the wallet owes the paymaster for an earlier operation's gas:
			// settle what it can in this operation
			debts := gasDebts(ctx, aw.Address, gc)
			var settle []aa.Call
			for _, d := range debts {
				if d.payable {
					settle = append(settle, d.calls...)
				}
			}
			if len(settle) > 0 {
				if sp, e := b.Prepare(ctx, aa.Request{Wallet: aw, Calls: append(settle, calls...), Validity: validity, NonceKey: nonceKey}); e == nil {
					p = sp
				}
			}
			for _, d := range debts {
				notes = append(notes, d.note(p != nil && d.payable))
			}
		}
		if p == nil {
			p, err = b.Prepare(ctx, aa.Request{Wallet: aw, Calls: calls, Validity: validity, NonceKey: nonceKey})
		} else {
			err = nil
		}
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
	return &PreparedWalletOperation{Transaction: base64.StdEncoding.EncodeToString(p.SafeOpHash.Bytes()), Prepared: p, Record: rec, Notes: notes}, nil
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
	// only the first to record the outcome applies its effects (the tracker
	// and WalletOperationOutcome can both see the receipt)
	res := gc.DB.Model(&userModels.WalletOperation{}).Where("id = ? AND status = ?", op.ID, userModels.WalletOperationSubmitted).Updates(updates)
	if res.Error != nil || res.RowsAffected == 0 {
		return
	}
	recordGasCharges(op, r.Receipt.TransactionHash, gc)
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
		if hook := minedHooks[op.Kind]; hook != nil {
			hook(op, r.Receipt.TransactionHash, gc)
		}
	}
}

// minedHooks finish flows whose effects depend on a mined operation's
// receipt, by operation kind (registered by the flows in init).
var minedHooks = map[string]func(op userModels.WalletOperation, txHash common.Hash, gc *sharedconfig.GlobalConfig){}

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

// recordGasCharges records the stablecoin gas the paymaster collected for a
// mined operation (and any earlier gas debt it settled) as GAS fee
// collections; a charge the paymaster could not collect becomes the
// wallet's debt and is alerted.
func recordGasCharges(op userModels.WalletOperation, txHash common.Hash, gc *sharedconfig.GlobalConfig) {
	paymaster := operationBuilder(gc).Paymaster
	if paymaster == (common.Address{}) || gc.BantuExpansionClient == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	rcpt, err := gc.BantuExpansionClient.TransactionReceipt(ctx, txHash)
	if err != nil {
		log.Printf("[recordGasCharges] receipt of %v: %v", txHash.Hex(), err)
		return
	}
	wallet := common.HexToAddress(op.WalletAddress)
	record := func(token common.Address, amount *big.Int) {
		var asset assetModels.CuratedAsset
		gc.DB.Where("LOWER(contract_address) = ?", strings.ToLower(token.Hex())).First(&asset)
		decimals, err := network.AssetDecimals(ctx, gc.BantuExpansionClient, basetxn.CreditAsset{Code: asset.AssetCode, Issuer: token.Hex()})
		if err != nil {
			log.Printf("[recordGasCharges] decimals of %v: %v", token.Hex(), err)
			return
		}
		contract := token.Hex()
		h := txHash.Hex()
		f := sharedconfig.FeeCollection{
			ID: gc.GenerateUUIDString(), FromWalletAddress: op.WalletAddress, FeeType: "GAS",
			Amount: decimal.NewFromBigInt(amount, -int32(decimals)).InexactFloat64(), AssetCode: asset.AssetCode,
			ContractAddress: &contract, DestinationWallet: paymaster.Hex(), TransactionHash: &h,
		}
		if w, err := userModels.UserWalletID(op.WalletAddress).GetWallet(gc.DB, gc); err == nil {
			f.FromWalletAlias = w.Alias
			if w.SharedAccessEnabled == 1 {
				f.SharedAccessOperation = 1
			}
			if owner, err := w.GetWalletOwner(gc.DB, gc); err == nil {
				f.FromUsername, f.BelongsToEnterpriseProfile = owner.Username, owner.CreatedByServiceLinkID
			}
		}
		if err := gc.DB.Omit(clause.Associations).Create(&f).Error; err != nil {
			log.Printf("[recordGasCharges] saving GAS fee for %v: %v", h, err)
		}
	}
	if op.UserOpHash != nil {
		if g := aa.ParseGasPayment(rcpt.Logs, paymaster, common.HexToHash(*op.UserOpHash)); g != nil {
			if g.Failed {
				gc.LogDiscordFailedRequest(fmt.Sprintf("[recordGasCharges] the paymaster could not collect %v (base units) of %v gas from %v for %v; it is recorded as the wallet's debt", g.TokenCost, g.Token.Hex(), op.WalletAddress, txHash.Hex()))
			} else {
				record(g.Token, g.TokenCost)
			}
		}
	}
	for _, d := range aa.ParseDebtSettlements(rcpt.Logs, paymaster, wallet) {
		record(d.Token, d.Amount)
	}
}
