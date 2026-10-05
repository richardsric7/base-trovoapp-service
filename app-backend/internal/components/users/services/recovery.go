package users

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"trovo-wallet-api/internal/aa"
	"trovo-wallet-api/internal/basetxn"
	usersDB "trovo-wallet-api/internal/components/users/db"
	userModels "trovo-wallet-api/internal/components/users/models"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/evmkeypair"
	"trovo-wallet-api/internal/gnosissafe"
	"trovo-wallet-api/internal/network"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/mailgun/mailgun-go/v4"
	"gorm.io/gorm/clause"
)

// Account recovery uses Candide's Social Recovery Module
// (RECOVERY_MODULE_ADDRESS, see recovery/): a user who opts in makes the
// platform's recovery guardian the only guardian of each covered wallet -
// their primary wallet and their own sub-wallets without approvers. To
// recover, after the security answers and email OTP, the guardian starts
// replacing the lost key with the new one on each covered wallet. That
// takes effect after the module's recovery period, during which the user
// is notified and can cancel with their current key; then the guardian
// finalizes it. The guardian can only replace owners, never move funds.
//
// Wallets the user shares with approvers (theirs or others') are not
// covered: once the recovery completes, each gets a REPLACE SIGNER
// approval request swapping the old key for the new one, which the other
// co-signers approve.

// Wallet operation kinds of account recovery.
const (
	OperationRecoveryEnable  = "ACCOUNT RECOVERY"
	OperationRecoveryDisable = "ACCOUNT RECOVERY DISABLE"
	OperationRecoveryCancel  = "ACCOUNT RECOVERY CANCEL"
	OperationReplaceSigner   = "REPLACE SIGNER"
)

// recoveryOpContext is what an enable operation keeps (its fee).
type recoveryOpContext struct {
	Fees []sharedconfig.FeeCollection `json:"fees"`
}

func recoveryNotConfigured() error {
	return &tErrors.CustomError{Param: "username", Err: "error-account-recovery-not-configured", ErrMessage: "Account recovery is not available at the moment.", Code: http.StatusServiceUnavailable}
}

// recoveryModule is RECOVERY_MODULE_ADDRESS.
func recoveryModule() (common.Address, error) {
	raw := strings.TrimSpace(os.Getenv("RECOVERY_MODULE_ADDRESS"))
	if !common.IsHexAddress(raw) {
		return common.Address{}, recoveryNotConfigured()
	}
	return common.HexToAddress(raw), nil
}

// recoveryGuardian is the platform's recovery guardian: the Safe at
// ACCOUNT_RECOVERY_GUARDIAN_SAFE, owned by the
// ACCOUNT_RECOVERY_GUARDIAN_SIGNERS keys (recommended), or - when no Safe
// is set - the first of those keys itself.
type recoveryGuardian struct {
	Address common.Address
	Safe    bool
	Signers []*evmkeypair.Full
}

func loadRecoveryGuardian(ctx context.Context, gc *sharedconfig.GlobalConfig) (*recoveryGuardian, error) {
	raw := strings.TrimSpace(os.Getenv("ACCOUNT_RECOVERY_GUARDIAN_SIGNERS"))
	var keys []*evmkeypair.Full
	for _, k := range strings.Split(raw, ";") {
		if k = strings.TrimSpace(k); k == "" {
			continue
		}
		kp, err := evmkeypair.ParseFull(k)
		if err != nil {
			log.Printf("[recoveryGuardian] ACCOUNT_RECOVERY_GUARDIAN_SIGNERS: %v", err)
			return nil, recoveryNotConfigured()
		}
		keys = append(keys, kp)
	}
	if len(keys) == 0 {
		return nil, recoveryNotConfigured()
	}
	if s := strings.TrimSpace(os.Getenv("ACCOUNT_RECOVERY_GUARDIAN_SAFE")); s != "" {
		if !common.IsHexAddress(s) {
			return nil, recoveryNotConfigured()
		}
		safe := common.HexToAddress(s)
		signers, err := gnosissafe.SignersForSafe(ctx, gc.BantuExpansionClient, safe, keys)
		if err != nil {
			log.Printf("[recoveryGuardian] signers for %v: %v", safe.Hex(), err)
			return nil, recoveryNotConfigured()
		}
		return &recoveryGuardian{Address: safe, Safe: true, Signers: signers}, nil
	}
	return &recoveryGuardian{Address: common.HexToAddress(keys[0].Address()), Signers: keys[:1]}, nil
}

// recoveryGuardianAddress is the guardian's address, without checking the
// keys a Safe guardian's owners hold (for building wallets' enable calls).
func recoveryGuardianAddress() (common.Address, error) {
	if s := strings.TrimSpace(os.Getenv("ACCOUNT_RECOVERY_GUARDIAN_SAFE")); s != "" {
		if !common.IsHexAddress(s) {
			return common.Address{}, recoveryNotConfigured()
		}
		return common.HexToAddress(s), nil
	}
	keys := strings.FieldsFunc(os.Getenv("ACCOUNT_RECOVERY_GUARDIAN_SIGNERS"), func(r rune) bool { return r == ';' })
	if len(keys) == 0 {
		return common.Address{}, recoveryNotConfigured()
	}
	kp, err := evmkeypair.ParseFull(strings.TrimSpace(keys[0]))
	if err != nil {
		return common.Address{}, recoveryNotConfigured()
	}
	return common.HexToAddress(kp.Address()), nil
}

// exec has the guardian make calls (it pays their gas) and waits for them.
func (g *recoveryGuardian) exec(ctx context.Context, calls []aa.Call, gc *sharedconfig.GlobalConfig) ([]string, error) {
	client := gc.BantuExpansionClient
	chainID := network.GetBlockchainChainID()
	if g.Safe {
		var sc []gnosissafe.Call
		for _, c := range calls {
			sc = append(sc, gnosissafe.Call{To: c.To, Value: big.NewInt(0), Data: c.Data})
		}
		h, err := gnosissafe.ExecCalls(ctx, client, chainID, g.Address, g.Signers, sc, common.HexToAddress(envOrDefault("SAFE_MULTISEND_CALL_ONLY_ADDRESS", gnosissafe.DefaultMultiSendCallOnlyAddress)))
		if err != nil {
			return nil, err
		}
		return []string{h}, gnosissafe.WaitSuccess(ctx, client, common.HexToHash(h))
	}
	var hashes []string
	for _, c := range calls {
		h, err := gnosissafe.SendTransaction(ctx, client, chainID, g.Signers[0], c.To, c.Data)
		if err != nil {
			return hashes, err
		}
		hashes = append(hashes, h)
		if err := gnosissafe.WaitSuccess(ctx, client, common.HexToHash(h)); err != nil {
			return hashes, err
		}
	}
	return hashes, nil
}

// recoveryWallet is a wallet's recovery state.
type recoveryWallet struct {
	Wallet        userModels.UserWallet
	Owners        []common.Address
	Threshold     int64
	ModuleEnabled bool
	Guarded       bool
}

// ownedRecoveryWallets lists the user's own deployed wallets whose owners
// include signer, without approvers (the wallets recovery covers), with
// their recovery state.
func ownedRecoveryWallets(ctx context.Context, user *userModels.User, signer common.Address, module, guardian common.Address, gc *sharedconfig.GlobalConfig) ([]recoveryWallet, error) {
	var out []recoveryWallet
	for _, w := range user.GetAllWallets(gc) {
		if w.UserID != user.ID || (w.SharedAccessEnabled == 1 && w.NumberOfApprovalsNeeded > 0) {
			continue
		}
		addr := common.HexToAddress(w.ID)
		deployed, err := aa.Deployed(ctx, gc.BantuExpansionClient, addr)
		if err != nil {
			return nil, &tErrors.ErrorTemporaryServerError{}
		}
		if !deployed {
			continue
		}
		owners, threshold, err := aa.OnchainOwners(ctx, gc.BantuExpansionClient, addr)
		if err != nil {
			return nil, &tErrors.ErrorTemporaryServerError{}
		}
		if !containsAddress(owners, signer) {
			continue
		}
		rw := recoveryWallet{Wallet: w, Owners: owners, Threshold: threshold}
		if rw.ModuleEnabled, err = aa.ModuleEnabled(ctx, gc.BantuExpansionClient, addr, module); err != nil {
			return nil, &tErrors.ErrorTemporaryServerError{}
		}
		if rw.ModuleEnabled {
			if rw.Guarded, err = aa.IsGuardian(ctx, gc.BantuExpansionClient, module, addr, guardian); err != nil {
				return nil, &tErrors.ErrorTemporaryServerError{}
			}
		}
		out = append(out, rw)
	}
	return out, nil
}

// guardianRevokeCalls are wallet's calls removing the platform's recovery
// guardian when it gets co-signers (target has more than one owner), so the
// platform can never replace keys on a wallet approvers control. None when
// recovery is not configured or the wallet is not covered.
func guardianRevokeCalls(ctx context.Context, wallet common.Address, target []common.Address, gc *sharedconfig.GlobalConfig) ([]aa.Call, error) {
	if len(target) < 2 {
		return nil, nil
	}
	module, err := recoveryModule()
	if err != nil {
		return nil, nil
	}
	guardian, err := recoveryGuardianAddress()
	if err != nil {
		return nil, nil
	}
	if deployed, err := aa.Deployed(ctx, gc.BantuExpansionClient, wallet); err != nil {
		return nil, &tErrors.ErrorTemporaryServerError{}
	} else if !deployed {
		return nil, nil
	}
	guarded, err := aa.IsGuardian(ctx, gc.BantuExpansionClient, module, wallet, guardian)
	if err != nil {
		return nil, &tErrors.ErrorTemporaryServerError{}
	}
	if !guarded {
		return nil, nil
	}
	return []aa.Call{aa.RevokeGuardianCall(module, guardian)}, nil
}

func containsAddress(list []common.Address, a common.Address) bool {
	for _, x := range list {
		if x == a {
			return true
		}
	}
	return false
}

// recoveryOps prepares one operation per wallet (the user signs each), or
// submits them when the payload carries their signatures.
func recoveryOps(ctx context.Context, kind string, user *userModels.User, wallets []recoveryWallet, callsFor func(rw recoveryWallet, first bool) ([]aa.Call, interface{}, []string, error), payload *userModels.UserAccountRecoveryPayload, gc *sharedconfig.GlobalConfig) (submitted bool, err error) {
	sigs := payload.TransactionSignatures
	txs := payload.Transactions
	if len(sigs) > 0 {
		if len(sigs) != len(txs) {
			return false, &tErrors.CustomError{Param: "transactionSignatures", Err: "error-signatures-missing", ErrMessage: "Sign every transaction (one per wallet).", Code: http.StatusBadRequest}
		}
		// every operation must be one prepared for this user's wallets
		type loaded struct {
			rec *userModels.WalletOperation
			p   *aa.Prepared
		}
		var ops []loaded
		for _, tx := range txs {
			rec, p, e := loadOwnRecoveryOperation(tx, kind, user, gc)
			if e != nil {
				return false, e
			}
			ops = append(ops, loaded{rec, p})
		}
		var hashes []string
		for i, o := range ops {
			h, e := SignSingleOwnerOperation(ctx, o.rec, o.p, user.PrimarySigner, sigs[i], gc)
			if e != nil {
				if len(hashes) > 0 {
					gc.LogDiscordFailedRequest(fmt.Sprintf("[%v] %v: %d of %d operations submitted before %v failed: %v", kind, user.Username, len(hashes), len(ops), o.rec.WalletAddress, e))
				}
				return len(hashes) > 0, e
			}
			recordOperationFees(o.rec, h, gc)
			hashes = append(hashes, h)
		}
		payload.TransactionID = strings.Join(hashes, ",")
		return true, nil
	}

	payload.Transactions, payload.Wallets, payload.Messages = nil, nil, nil
	for i, rw := range wallets {
		calls, opContext, notes, e := callsFor(rw, i == 0)
		if e != nil {
			return false, e
		}
		owner := *user
		w := rw.Wallet
		op, e := PrepareWalletOperation(ctx, kind, user, &owner, &w, calls, 0, opContext, gc)
		if e != nil {
			return false, e
		}
		payload.Transactions = append(payload.Transactions, op.Transaction)
		payload.Wallets = append(payload.Wallets, w.Alias)
		payload.Messages = append(payload.Messages, notes...)
		for _, m := range op.Messages() {
			payload.Messages = append(payload.Messages, fmt.Sprintf("%v: %v", w.Alias, m))
		}
	}
	payload.NetworkPassPhrase = network.GetBlockchainNetworkPassPhrase()
	return false, nil
}

// loadOwnRecoveryOperation loads a recovery operation and checks it is
// for one of user's own wallets.
func loadOwnRecoveryOperation(tx, kind string, user *userModels.User, gc *sharedconfig.GlobalConfig) (*userModels.WalletOperation, *aa.Prepared, error) {
	var rec userModels.WalletOperation
	raw, err := base64.StdEncoding.DecodeString(strings.TrimSpace(tx))
	if err != nil || len(raw) != 32 || gc.DB.Where("id = ?", common.BytesToHash(raw).Hex()).First(&rec).Error != nil {
		return nil, nil, &tErrors.CustomError{Param: "transactions", Err: "error-invalid-transaction", ErrMessage: "Unknown transaction. Please start again.", Code: http.StatusBadRequest}
	}
	w, err := userModels.UserWalletID(rec.WalletAddress).GetWallet(gc.DB, gc)
	if err != nil || w.UserID != user.ID {
		return nil, nil, &tErrors.CustomError{Param: "transactions", Err: "error-transaction-mismatch", ErrMessage: "This transaction does not belong to your wallets.", Code: http.StatusBadRequest}
	}
	return LoadWalletOperation(tx, rec.WalletAddress, kind, gc)
}

// recoveryFeeCalls are the primary wallet's calls paying the account
// recovery fee (ACCOUNT_RECOVERY_FEE, set in USD, paid in its stablecoin).
func recoveryFeeCalls(ctx context.Context, user *userModels.User, primary *userModels.UserWallet, gc *sharedconfig.GlobalConfig) ([]aa.Call, []sharedconfig.FeeCollection, []string, error) {
	fee := primary.GetAccountRecoveryFee(gc)
	if fee.Inactive != 0 || fee.FeeFixed <= 0 || feeExemptProfile(user.Username, gc) {
		return nil, nil, nil, nil
	}
	feeAddr, err := feeWalletAddress(fee.FeeWalletSecretKey, "account recovery fee wallet", gc)
	if err != nil {
		return nil, nil, nil, err
	}
	if !common.IsHexAddress(fee.FeeContractAddress) {
		gc.LogDiscordFailedRequest("[EnableAccountRecovery] ACCOUNT_RECOVERY_FEE has no fee_contract_address (the fee must be a stablecoin)")
		return nil, nil, nil, recoveryNotConfigured()
	}
	asset := basetxn.CreditAsset{Code: fee.FeeAssetCode, Issuer: fee.FeeContractAddress}
	amount, err := usdAmountIn(fee.FeeFixed, fee.FeeAssetCode, fee.FeeContractAddress, gc)
	if err != nil {
		return nil, nil, nil, err
	}
	decimals, err := network.AssetDecimals(ctx, gc.BantuExpansionClient, asset)
	if err != nil {
		return nil, nil, nil, &tErrors.ErrorTemporaryServerError{}
	}
	amount = amount.Truncate(int32(decimals))
	units, err := baseUnits(amount.String(), decimals)
	if err != nil {
		return nil, nil, nil, err
	}
	bal, err := network.B20BalanceOf(gc.BantuExpansionClient, asset.Issuer, primary.ID, decimals)
	if err != nil {
		return nil, nil, nil, &tErrors.ErrorTemporaryServerError{}
	}
	if bal.LessThan(amount) {
		return nil, nil, nil, &tErrors.CustomError{Param: "username", Err: "error-primary-wallet-underfunded", ErrMessage: fmt.Sprintf("%v %v is required on wallet %v to pay the account recovery fee. Please first fund the wallet with at least %v %v.", amount, asset.Code, user.Username, amount.Sub(bal), asset.Code), Code: http.StatusBadRequest}
	}
	note := fmt.Sprintf("%v %v ($%v USD) will be deducted from wallet %v as the account recovery fee.", amount, asset.Code, fee.FeeFixed, user.Username)
	return []aa.Call{transferCall(asset, feeAddr, units)}, []sharedconfig.FeeCollection{feeRecord("ACCOUNT_RECOVERY", primary, user, asset, amount.String(), feeAddr, 0)}, []string{note}, nil
}

// recoveryPeriodNote describes the recovery period for messages.
func recoveryPeriodNote() string {
	if s, err := strconv.ParseInt(strings.TrimSpace(os.Getenv("RECOVERY_PERIOD_SECONDS")), 10, 64); err == nil && s > 0 {
		d := time.Duration(s) * time.Second
		if d >= 24*time.Hour {
			return fmt.Sprintf("%v days", int(d.Hours()/24))
		}
		return d.String()
	}
	return "the recovery period"
}

// EnableAccountRecovery opts the user into account recovery: one operation
// per covered wallet (primary wallet and own sub-wallets without
// approvers, not yet covered) enables the recovery module and makes the
// platform's guardian its only guardian; the primary wallet's also pays the
// fee, the first time. Without signatures it returns the operations
// (transactions, one per wallet in wallets) for the user to sign; with
// them (transactionSignatures, same order) it submits them. Calling it
// again later covers wallets created since, without a fee.
func EnableAccountRecovery(user *userModels.User, payload *userModels.UserAccountRecoveryPayload, gc *sharedconfig.GlobalConfig) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	if user.HasSecurityQuestions == 0 {
		return &tErrors.CustomError{Param: "username", Err: "error security answers not set", ErrMessage: "security answers has not been set for this account."}
	}
	module, err := recoveryModule()
	if err != nil {
		return err
	}
	guardian, err := recoveryGuardianAddress()
	if err != nil {
		return err
	}
	primary, err := userModels.UserWalletID(user.Address).GetWallet(gc.DB, gc)
	if err != nil {
		return &tErrors.ErrorTemporaryServerError{}
	}
	if deployed, err := aa.Deployed(ctx, gc.BantuExpansionClient, common.HexToAddress(primary.ID)); err != nil {
		return &tErrors.ErrorTemporaryServerError{}
	} else if !deployed {
		return &tErrors.CustomError{Param: "username", Err: "error primary account not yet activated", ErrMessage: "Your primary wallet is not activated yet. It activates with its first transaction; then you can turn on account recovery.", Code: http.StatusBadRequest}
	}
	all, err := ownedRecoveryWallets(ctx, user, common.HexToAddress(user.PrimarySigner), module, guardian, gc)
	if err != nil {
		return err
	}
	// the primary wallet first: it pays the fee
	var todo []recoveryWallet
	for _, rw := range all {
		if !rw.Guarded && rw.Wallet.ID == primary.ID {
			todo = append([]recoveryWallet{rw}, todo...)
		} else if !rw.Guarded {
			todo = append(todo, rw)
		}
	}
	if len(todo) == 0 && len(payload.TransactionSignatures) == 0 {
		return &tErrors.CustomError{Param: "username", Err: "error account recovery already enabled.", ErrMessage: "Account recovery already covers all your wallets.", Code: http.StatusConflict}
	}
	chargeFee := user.AccountRecoveryEnabled == 0
	submitted, err := recoveryOps(ctx, OperationRecoveryEnable, user, todo, func(rw recoveryWallet, first bool) ([]aa.Call, interface{}, []string, error) {
		calls := aa.EnableRecoveryCalls(common.HexToAddress(rw.Wallet.ID), module, guardian, rw.ModuleEnabled)
		notes := []string{fmt.Sprintf("Wallet %v: the platform's recovery key can replace your key on it after %v if you lose it. It can never move your funds, and you can cancel any recovery during that time.", rw.Wallet.Alias, recoveryPeriodNote())}
		if !chargeFee || rw.Wallet.ID != primary.ID {
			return calls, nil, notes, nil
		}
		feeCalls, fees, feeNotes, e := recoveryFeeCalls(ctx, user, &rw.Wallet, gc)
		if e != nil {
			return nil, nil, nil, e
		}
		return append(calls, feeCalls...), recoveryOpContext{Fees: fees}, append(notes, feeNotes...), nil
	}, payload, gc)
	if err != nil || !submitted {
		return err
	}
	exp := time.Now().AddDate(1, 0, 0)
	if e := gc.DB.Model(&userModels.User{}).Where("id = ?", user.ID).Updates(map[string]interface{}{"account_recovery_enabled": 1, "account_recovery_expires_on": exp}).Error; e != nil {
		log.Printf("[EnableAccountRecovery] %v: operations submitted but not recorded: %v", user.Username, e)
	}
	user.InvalidateUserCache(gc)
	return nil
}

// DisableAccountRecovery removes the platform's guardian from each covered
// wallet, one operation per wallet the user signs (same two steps as
// EnableAccountRecovery).
func DisableAccountRecovery(user *userModels.User, payload *userModels.UserAccountRecoveryPayload, gc *sharedconfig.GlobalConfig) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	module, err := recoveryModule()
	if err != nil {
		return err
	}
	guardian, err := recoveryGuardianAddress()
	if err != nil {
		return err
	}
	var pending int64
	gc.DB.Model(&userModels.UserAccountRecoveryLog{}).Where("username = ? AND status = ?", user.Username, userModels.AccountRecoveryPending).Count(&pending)
	if pending > 0 {
		return &tErrors.CustomError{Param: "username", Err: "error-recovery-pending", ErrMessage: "A recovery of this account is in progress. Cancel it first.", Code: http.StatusConflict}
	}
	all, err := ownedRecoveryWallets(ctx, user, common.HexToAddress(user.PrimarySigner), module, guardian, gc)
	if err != nil {
		return err
	}
	var todo []recoveryWallet
	for _, rw := range all {
		if rw.Guarded {
			todo = append(todo, rw)
		}
	}
	if len(todo) == 0 && len(payload.TransactionSignatures) == 0 {
		if user.AccountRecoveryEnabled == 1 {
			gc.DB.Model(&userModels.User{}).Where("id = ?", user.ID).Update("account_recovery_enabled", 0)
			user.InvalidateUserCache(gc)
			return nil
		}
		return &tErrors.CustomError{Param: "username", Err: "error account recovery not enabled.", ErrMessage: "Account recovery not enabled."}
	}
	submitted, err := recoveryOps(ctx, OperationRecoveryDisable, user, todo, func(rw recoveryWallet, _ bool) ([]aa.Call, interface{}, []string, error) {
		return []aa.Call{aa.RevokeGuardianCall(module, guardian)}, nil, []string{fmt.Sprintf("Wallet %v will no longer be recoverable by the platform.", rw.Wallet.Alias)}, nil
	}, payload, gc)
	if err != nil || !submitted {
		return err
	}
	gc.DB.Model(&userModels.User{}).Where("id = ?", user.ID).Update("account_recovery_enabled", 0)
	user.InvalidateUserCache(gc)
	return nil
}

// DoAccountRecovery starts recovering user's account onto
// payload.NewSignerAddress, after the security answers and email OTP. With
// commit 0 it only checks and describes; with commit 1 the guardian starts
// replacing the old key on each covered wallet, which completes after the
// recovery period (ProcessAccountRecoveries) unless the user cancels with
// their current key. Nothing changes in the database until then.
func DoAccountRecovery(user *userModels.User, payload *userModels.AccountRecoveryRequest, gc *sharedconfig.GlobalConfig) (multiAccessWallets []userModels.UserWallet, sharedApproverWallets []userModels.WalletPermission, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	payload.Messages = nil
	if !common.IsHexAddress(payload.NewSignerAddress) {
		return nil, nil, &tErrors.CustomError{Param: "newSignerAddress", Err: "error invalid new signer public key.", ErrMessage: "Invalid new signer public key."}
	}
	newKey := common.HexToAddress(payload.NewSignerAddress)
	oldKey := common.HexToAddress(user.PrimarySigner)
	if newKey == oldKey {
		return nil, nil, &tErrors.CustomError{Param: "newSignerAddress", Err: "error invalid new signer public key.", ErrMessage: "The new key is your current key."}
	}
	if _, e := usersDB.GetUser(payload.NewSignerAddress, gc.DB, gc); e == nil {
		return nil, nil, &tErrors.CustomError{Param: "newSignerAddress", Err: "error new signer public key already in use.", ErrMessage: "The new signer public key is already in use on another account."}
	}
	if _, e := usersDB.GetUserFromPrimarySigner(payload.NewSignerAddress, gc.DB, gc); e == nil {
		return nil, nil, &tErrors.CustomError{Param: "newSignerAddress", Err: "error new signer public key already in use.", ErrMessage: "The new signer public key is already in use on another account."}
	}
	if user.AccountRecoveryEnabled == 0 {
		return nil, nil, &tErrors.CustomError{Param: "username", Err: "error account recovery not enabled.", ErrMessage: "Account recovery not enabled."}
	}
	var pending userModels.UserAccountRecoveryLog
	if gc.DB.Where("username = ? AND status = ?", user.Username, userModels.AccountRecoveryPending).First(&pending).Error == nil {
		return nil, nil, &tErrors.CustomError{Param: "username", Err: "error-recovery-pending", ErrMessage: fmt.Sprintf("A recovery of this account is already in progress (completes %v).", pending.ExecuteAfter), Code: http.StatusConflict}
	}
	if !ValidateSecurityAnswers(user, payload.SecurityAnswers, gc) {
		return nil, nil, &tErrors.CustomError{Param: "username", Err: "error invalid security answers", ErrMessage: "Answers to the security questions are invalid."}
	}
	if CheckAccountRecoveryEmailOTP(user, payload.EmailOTP, gc.DB) != nil {
		return nil, nil, &tErrors.CustomError{Param: "username", Err: "error invalid email otp", ErrMessage: "Email OTP is invalid."}
	}
	module, err := recoveryModule()
	if err != nil {
		return nil, nil, err
	}
	guardian, err := loadRecoveryGuardian(ctx, gc)
	if err != nil {
		return nil, nil, err
	}
	all, err := ownedRecoveryWallets(ctx, user, oldKey, module, guardian.Address, gc)
	if err != nil {
		return nil, nil, err
	}
	var covered []recoveryWallet
	for _, rw := range all {
		if rw.Guarded {
			covered = append(covered, rw)
		}
	}
	if len(covered) == 0 {
		return nil, nil, &tErrors.CustomError{Param: "username", Err: "error-no-recoverable-wallet", ErrMessage: "None of your wallets is covered by account recovery.", Code: http.StatusBadRequest}
	}
	payload.Wallets = nil
	for _, rw := range covered {
		payload.Wallets = append(payload.Wallets, rw.Wallet.Alias)
	}
	for _, w := range user.GetAllWallets(gc) {
		if w.UserID == user.ID && w.SharedAccessEnabled == 1 && w.NumberOfApprovalsNeeded > 0 {
			multiAccessWallets = append(multiAccessWallets, w)
		}
	}
	for _, p := range user.WalletsSharedWithUser {
		if p.Permission == "APPROVER" {
			sharedApproverWallets = append(sharedApproverWallets, p)
		}
	}
	payload.Messages = append(payload.Messages, fmt.Sprintf("Your key is replaced on %v after %v. Until then you can cancel from a device that still has your current key.", strings.Join(payload.Wallets, ", "), recoveryPeriodNote()))
	if len(multiAccessWallets)+len(sharedApproverWallets) > 0 {
		payload.Messages = append(payload.Messages, "On wallets you share with approvers, your new key replaces the old one once the other co-signers approve the request you will be sent.")
	}
	if payload.Commit == 0 {
		return multiAccessWallets, sharedApproverWallets, nil
	}

	arl, err := startAccountRecovery(ctx, user, oldKey, newKey, module, guardian, covered, gc)
	if err != nil {
		return nil, nil, err
	}
	RemoveAccountRecoveryEmailOTP(user, payload.EmailOTP, gc.DB)
	payload.ExecuteAfter = arl.ExecuteAfter
	payload.TransactionID = arl.StartTxHashes
	return multiAccessWallets, sharedApproverWallets, nil
}

// startAccountRecovery has the guardian start replacing oldKey with newKey
// on each covered wallet, records the pending recovery and notifies the
// user.
func startAccountRecovery(ctx context.Context, user *userModels.User, oldKey, newKey, module common.Address, guardian *recoveryGuardian, covered []recoveryWallet, gc *sharedconfig.GlobalConfig) (*userModels.UserAccountRecoveryLog, error) {
	var calls []aa.Call
	var walletIDs []string
	for _, rw := range covered {
		addr := common.HexToAddress(rw.Wallet.ID)
		nonce, e := aa.RecoveryNonce(ctx, gc.BantuExpansionClient, module, addr)
		if e != nil {
			return nil, &tErrors.ErrorTemporaryServerError{}
		}
		calls = append(calls, aa.ConfirmRecoveryCall(module, addr, aa.ReplaceOwner(rw.Owners, oldKey, newKey), rw.Threshold, nonce))
		walletIDs = append(walletIDs, rw.Wallet.ID)
	}
	masterRecover := 0
	if user.PrimarySigner == user.Address {
		masterRecover = 1
	}
	// recorded first, so the watcher never sees this recovery on-chain
	// without its record
	arl := userModels.UserAccountRecoveryLog{
		ID: uuid.NewString(), Username: user.Username, OldSignerAddress: strings.ToUpper(oldKey.Hex()), NewSignerAddress: strings.ToUpper(newKey.Hex()), MasterWallet: masterRecover,
		Status: userModels.AccountRecoveryPending, Wallets: strings.Join(walletIDs, ","),
	}
	if e := gc.DB.Omit(clause.Associations).Create(&arl).Error; e != nil {
		log.Printf("[startAccountRecovery] %v: recording: %v", user.Username, e)
		return nil, &tErrors.ErrorTemporaryServerError{}
	}
	hashes, err := guardian.exec(ctx, calls, gc)
	if err != nil {
		log.Printf("[startAccountRecovery] %v: starting recovery (%v): %v", user.Username, hashes, err)
		gc.LogDiscordFailedRequest(fmt.Sprintf("[startAccountRecovery] starting recovery %v of %v failed after %v: %v", arl.ID, user.Username, hashes, err))
		gc.DB.Model(&arl).Updates(map[string]interface{}{"status": userModels.AccountRecoveryFailed, "start_tx_hashes": strings.Join(hashes, ",")})
		return nil, &tErrors.CustomError{Param: "username", Err: "error-recovery-not-started", ErrMessage: "The recovery could not be started. Please try again later.", Code: http.StatusInternalServerError}
	}
	arl.StartTxHashes = strings.Join(hashes, ",")
	req, err := aa.ReadRecoveryRequest(ctx, gc.BantuExpansionClient, module, common.HexToAddress(walletIDs[0]))
	if err != nil || req.ExecutableAt == 0 {
		gc.LogDiscordFailedRequest(fmt.Sprintf("[startAccountRecovery] recovery %v of %v sent (%v) but no request is pending on %v", arl.ID, user.Username, hashes, walletIDs[0]))
		gc.DB.Model(&arl).Update("start_tx_hashes", arl.StartTxHashes)
		return nil, &tErrors.ErrorTemporaryServerError{}
	}
	after := time.Unix(int64(req.ExecutableAt), 0).UTC()
	arl.ExecuteAfter = &after
	if e := gc.DB.Model(&arl).Updates(map[string]interface{}{"start_tx_hashes": arl.StartTxHashes, "execute_after": after}).Error; e != nil {
		log.Printf("[startAccountRecovery] %v: recovery started but not recorded: %v", user.Username, e)
		gc.LogDiscordFailedRequest(fmt.Sprintf("[startAccountRecovery] recovery %v of %v started (%v) but its details were not recorded: %v", arl.ID, user.Username, hashes, e))
	}
	notifyRecovery(user, "Account recovery started",
		fmt.Sprintf("A recovery of your Trovo account %v was requested. Your wallets move to a new key on %v UTC. If this was not you, open the app on a device with your current key and cancel the recovery now.", user.Username, after.Format("2 Jan 2006 15:04")), gc)
	return &arl, nil
}

// CancelAccountRecovery cancels the user's pending recovery: one operation
// per wallet being recovered, calling the module's cancelRecovery, signed
// with the current key (same two steps as EnableAccountRecovery).
func CancelAccountRecovery(user *userModels.User, payload *userModels.UserAccountRecoveryPayload, gc *sharedconfig.GlobalConfig) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	var arl userModels.UserAccountRecoveryLog
	if gc.DB.Where("username = ? AND status = ?", user.Username, userModels.AccountRecoveryPending).First(&arl).Error != nil {
		return &tErrors.CustomError{Param: "username", Err: "error-no-recovery-pending", ErrMessage: "No recovery of this account is in progress.", Code: http.StatusNotFound}
	}
	module, err := recoveryModule()
	if err != nil {
		return err
	}
	var todo []recoveryWallet
	for _, id := range arl.WalletList() {
		w, e := userModels.UserWalletID(id).GetWallet(gc.DB, gc)
		if e != nil {
			continue
		}
		req, e := aa.ReadRecoveryRequest(ctx, gc.BantuExpansionClient, module, common.HexToAddress(id))
		if e != nil {
			return &tErrors.ErrorTemporaryServerError{}
		}
		if req.ExecutableAt > 0 {
			todo = append(todo, recoveryWallet{Wallet: w})
		}
	}
	if arl.ExecuteAfter == nil && len(todo) == 0 {
		return &tErrors.CustomError{Param: "username", Err: "error-recovery-starting", ErrMessage: "The recovery is being started. Please try again in a minute.", Code: http.StatusConflict}
	}
	submitted, err := recoveryOps(ctx, OperationRecoveryCancel, user, todo, func(rw recoveryWallet, _ bool) ([]aa.Call, interface{}, []string, error) {
		return []aa.Call{aa.CancelRecoveryCall(module)}, nil, []string{fmt.Sprintf("Cancels the recovery of wallet %v.", rw.Wallet.Alias)}, nil
	}, payload, gc)
	if err != nil {
		return err
	}
	if submitted || len(todo) == 0 {
		gc.DB.Model(&arl).Update("status", userModels.AccountRecoveryCanceled)
		notifyRecovery(user, "Account recovery canceled", fmt.Sprintf("The recovery of your Trovo account %v was canceled. Your current key stays in control.", user.Username), gc)
	}
	return nil
}

// AccountRecoveryStatus is the public state of a user's latest recovery.
type AccountRecoveryStatus struct {
	Status       string     `json:"status"`
	ExecuteAfter *time.Time `json:"executeAfter"`
	CompletedAt  *time.Time `json:"completedAt"`
}

// GetAccountRecoveryStatus reports the latest recovery of username onto
// newSigner (none: empty status).
func GetAccountRecoveryStatus(username, newSigner string, gc *sharedconfig.GlobalConfig) AccountRecoveryStatus {
	var arl userModels.UserAccountRecoveryLog
	if gc.DB.Where("username = ? AND LOWER(new_signer_address) = ?", username, strings.ToLower(newSigner)).Order("created_at DESC").First(&arl).Error != nil {
		return AccountRecoveryStatus{}
	}
	return AccountRecoveryStatus{Status: arl.Status, ExecuteAfter: arl.ExecuteAfter, CompletedAt: arl.CompletedAt}
}

// ProcessAccountRecoveries finalizes recoveries whose period is over,
// completes them once every wallet has the new key, notices recoveries
// canceled on-chain, and alerts on recoveries started outside this
// backend. It runs in the background (see main.go).
func ProcessAccountRecoveries(ctx context.Context, gc *sharedconfig.GlobalConfig) {
	module, err := recoveryModule()
	if err != nil || gc.BantuExpansionClient == nil {
		return
	}
	// the module compares the period with block time
	head, err := gc.BantuExpansionClient.HeaderByNumber(ctx, nil)
	if err != nil {
		return
	}
	var logs []userModels.UserAccountRecoveryLog
	gc.DB.Where("status = ?", userModels.AccountRecoveryPending).Limit(50).Find(&logs)
	var guardian *recoveryGuardian
	for i := range logs {
		arl := &logs[i]
		if arl.ExecuteAfter == nil {
			// still being started; one that never got its details failed
			if time.Since(arl.CreatedAt) > 15*time.Minute {
				gc.DB.Model(arl).Update("status", userModels.AccountRecoveryFailed)
				gc.LogDiscordFailedRequest(fmt.Sprintf("[ProcessAccountRecoveries] recovery %v of %v was never confirmed started (tx %v) - check its wallets %v on-chain", arl.ID, arl.Username, arl.StartTxHashes, arl.Wallets))
			}
			continue
		}
		newKey, oldKey := common.HexToAddress(arl.NewSignerAddress), common.HexToAddress(arl.OldSignerAddress)
		var finalize []aa.Call
		done, canceled, failed := 0, false, false
		for _, id := range arl.WalletList() {
			addr := common.HexToAddress(id)
			owners, _, e := aa.OnchainOwners(ctx, gc.BantuExpansionClient, addr)
			if e != nil {
				failed = true
				break
			}
			if containsAddress(owners, newKey) && !containsAddress(owners, oldKey) {
				done++
				continue
			}
			req, e := aa.ReadRecoveryRequest(ctx, gc.BantuExpansionClient, module, addr)
			if e != nil {
				failed = true
				break
			}
			if req.ExecutableAt == 0 {
				canceled = true
				break
			}
			if head.Time >= req.ExecutableAt {
				finalize = append(finalize, aa.FinalizeRecoveryCall(module, addr))
			}
		}
		switch {
		case failed:
			continue
		case canceled:
			gc.DB.Model(arl).Update("status", userModels.AccountRecoveryCanceled)
			if done > 0 {
				gc.LogDiscordFailedRequest(fmt.Sprintf("[ProcessAccountRecoveries] recovery %v of %v was canceled after %d wallet(s) moved to the new key - reconcile manually", arl.ID, arl.Username, done))
			}
			if u, e := userModels.Username(arl.Username).GetSimpleUser(gc.DB, gc); e == nil {
				notifyRecovery(&u, "Account recovery canceled", fmt.Sprintf("The recovery of your Trovo account %v was canceled.", arl.Username), gc)
			}
		case len(finalize) > 0:
			if guardian == nil {
				if guardian, err = loadRecoveryGuardian(ctx, gc); err != nil {
					return
				}
			}
			if hashes, e := guardian.exec(ctx, finalize, gc); e != nil {
				log.Printf("[ProcessAccountRecoveries] finalizing %v: %v (%v)", arl.ID, e, hashes)
				gc.LogDiscordFailedRequest(fmt.Sprintf("[ProcessAccountRecoveries] finalizing recovery %v of %v failed: %v", arl.ID, arl.Username, e))
			}
		case done == len(arl.WalletList()):
			completeAccountRecovery(ctx, arl, gc)
		}
	}
	watchRecoveries(ctx, module, gc)
}

// completeAccountRecovery records a recovery whose wallets all have the
// new key: the user's signer becomes the new key, and each wallet they share
// with approvers gets a REPLACE SIGNER request.
func completeAccountRecovery(ctx context.Context, arl *userModels.UserAccountRecoveryLog, gc *sharedconfig.GlobalConfig) {
	user, err := userModels.Username(arl.Username).GetFullUser(gc.DB, gc)
	if err != nil {
		return
	}
	newKey := arl.NewSignerAddress
	now := time.Now().UTC()
	tx := gc.DB.Begin()
	defer tx.Rollback()
	if e := tx.Model(&userModels.User{}).Where("id = ?", user.ID).Updates(map[string]interface{}{"primary_signer": newKey, "last_recovered_account_on": now}).Error; e != nil {
		log.Printf("[completeAccountRecovery] %v: %v", arl.ID, e)
		return
	}
	if e := tx.Model(&userModels.UserWallet{}).Where("user_id = ? AND LOWER(signer) = ?", user.ID, strings.ToLower(arl.OldSignerAddress)).Update("signer", newKey).Error; e != nil {
		log.Printf("[completeAccountRecovery] %v: %v", arl.ID, e)
		return
	}
	if e := tx.Model(arl).Updates(map[string]interface{}{"status": userModels.AccountRecoveryCompleted, "completed_at": now}).Error; e != nil {
		return
	}
	if tx.Commit().Error != nil {
		return
	}
	user.InvalidateUserCache(gc)
	user.InvalidateUserWalletCache(gc)
	user.PrimarySigner = newKey
	requests := createReplaceSignerRequests(ctx, &user, common.HexToAddress(arl.OldSignerAddress), common.HexToAddress(newKey), gc)
	msg := fmt.Sprintf("Your Trovo account %v has been recovered. Import it on your new device with your new key.", user.Username)
	if requests > 0 {
		msg += fmt.Sprintf(" The co-signers of %d shared wallet(s) have been asked to approve your new key.", requests)
	}
	notifyRecovery(&user, "Account recovery complete", msg, gc)
}

// createReplaceSignerRequests asks the co-signers of each wallet the user
// shares with approvers - their own, and those they approve on - to swap
// oldKey for newKey. It returns how many requests it created.
func createReplaceSignerRequests(ctx context.Context, user *userModels.User, oldKey, newKey common.Address, gc *sharedconfig.GlobalConfig) int {
	seen := map[string]bool{}
	var wallets []userModels.UserWallet
	for _, w := range user.GetAllWallets(gc) {
		if w.UserID == user.ID && w.SharedAccessEnabled == 1 && w.NumberOfApprovalsNeeded > 0 {
			wallets = append(wallets, w)
		}
	}
	for _, p := range user.WalletsSharedWithUser {
		if p.Permission != "APPROVER" {
			continue
		}
		if w, e := userModels.UserWalletID(p.WalletAddress).GetWallet(gc.DB, gc); e == nil {
			wallets = append(wallets, w)
		}
	}
	created := 0
	for i := range wallets {
		w := &wallets[i]
		if seen[w.ID] {
			continue
		}
		seen[w.ID] = true
		// a distribution wallet follows its issuing wallet
		var issuing userModels.UserWallet
		if gc.DB.Where("linked_wallet_address = ? AND wallet_type = ?", w.ID, 1).First(&issuing).Error == nil {
			continue
		}
		aw, e := walletOwners(ctx, w, gc)
		if e != nil || !containsAddress(aw.Owners, oldKey) {
			continue
		}
		target := aa.ReplaceOwner(aw.Owners, oldKey, newKey)
		var linked *userModels.UserWallet
		if w.WalletType == 1 && w.LinkedWalletAddress != nil {
			if lw, e := userModels.UserWalletID(*w.LinkedWalletAddress).GetWallet(gc.DB, gc); e == nil {
				linked = &lw
			}
		}
		calls, e := sharedAccessCalls(ctx, w, linked, target, aw.Threshold, gc)
		if e != nil || len(calls) == 0 {
			log.Printf("[createReplaceSignerRequests] %v on %v: %v", user.Username, w.ID, e)
			continue
		}
		owner, e := w.GetWalletOwner(gc.DB, gc)
		if e != nil {
			continue
		}
		op, e := PrepareWalletOperation(ctx, OperationReplaceSigner, user, &owner, w, calls, sharedWalletOperationValidity(), nil, gc)
		if e != nil {
			gc.LogDiscordFailedRequest(fmt.Sprintf("[createReplaceSignerRequests] %v's new key on %v: %v", user.Username, w.Alias, e))
			continue
		}
		others := len(aw.Owners) - 1
		description := fmt.Sprintf("%v recovered their account. Approve to replace their old key %v with their new key %v on wallet %v.", user.Username, oldKey.Hex(), newKey.Hex(), w.Alias)
		if int64(others) < aw.Threshold {
			description += fmt.Sprintf(" Warning: this needs %d approvals but only %d other co-signer(s) remain.", aw.Threshold, others)
			gc.LogDiscordFailedRequest(fmt.Sprintf("[createReplaceSignerRequests] %v: wallet %v needs %d approvals but only %d other co-signers remain", user.Username, w.Alias, aw.Threshold, others))
		}
		pa := userModels.PendingAuth{
			ID: uuid.NewString(), Initiator: user.Username, InitiatorSignerAddress: newKey.Hex(), WalletAddress: w.ID,
			TransactionType: OperationReplaceSigner, Description: description, ApprovalsNeeded: int(aw.Threshold), TransactionXdr: op.Transaction,
		}
		if e := gc.DB.Omit(clause.Associations).Create(&pa).Error; e != nil {
			continue
		}
		created++
		for _, p := range w.GetPermissionList(gc.DB) {
			if p.Permission != "APPROVER" || p.TargetUsername == user.Username {
				continue
			}
			if u, e := userModels.Username(p.TargetUsername).GetSimpleUser(gc.DB, gc); e == nil {
				u.SendPushMessage(fmt.Sprintf("%v needs your approval", user.Username), description, "", map[string]string{"route": "pendingApproval"}, gc)
			}
		}
	}
	return created
}

// unexpectedRecoveryAlert alerts the team (replaced in tests).
var unexpectedRecoveryAlert = func(gc *sharedconfig.GlobalConfig, msg string) { gc.LogDiscordFailedRequest(msg) }

// watchRecoveries alerts on recoveries started on-chain that this backend
// did not start (e.g. a misused guardian key): the wallet's owner must
// cancel them within the recovery period.
func watchRecoveries(ctx context.Context, module common.Address, gc *sharedconfig.GlobalConfig) {
	latest, err := gc.BantuExpansionClient.BlockNumber(ctx)
	if err != nil {
		return
	}
	cursor := userModels.RecoveryWatchCursor{ID: "recovery-module"}
	if gc.DB.First(&cursor, "id = ?", cursor.ID).Error != nil {
		// start from now
		cursor.Block = latest
		gc.DB.Create(&cursor)
		return
	}
	from := cursor.Block + 1
	if from > latest {
		return
	}
	to := latest
	if to-from > 5000 {
		to = from + 5000
	}
	logs, err := gc.BantuExpansionClient.FilterLogs(ctx, ethereum.FilterQuery{FromBlock: new(big.Int).SetUint64(from), ToBlock: new(big.Int).SetUint64(to), Addresses: []common.Address{module}, Topics: [][]common.Hash{{aa.RecoveryStartedTopic()}}})
	if err != nil {
		return
	}
	for _, ev := range aa.ParseRecoveryStarted(logs, module) {
		var n int64
		// started by this backend: by its transaction, or (while it is being
		// started) a pending recovery of the wallet
		gc.DB.Model(&userModels.UserAccountRecoveryLog{}).Where("LOWER(start_tx_hashes) LIKE ? OR (status = ? AND LOWER(wallets) LIKE ?)",
			"%"+strings.ToLower(ev.TxHash.Hex())+"%", userModels.AccountRecoveryPending, "%"+strings.ToLower(ev.Wallet.Hex())+"%").Count(&n)
		if n > 0 {
			continue
		}
		after := time.Unix(int64(ev.ExecutableAt), 0).UTC()
		unexpectedRecoveryAlert(gc, fmt.Sprintf("[watchRecoveries] UNEXPECTED recovery of wallet %v started in %v (new owners %v, executable %v) - not started by this backend", ev.Wallet.Hex(), ev.TxHash.Hex(), ev.NewOwners, after))
		if w, e := userModels.UserWalletID(strings.ToUpper(ev.Wallet.Hex())).GetWallet(gc.DB, gc); e == nil {
			if owner, e := w.GetWalletOwner(gc.DB, gc); e == nil {
				notifyRecovery(&owner, "Urgent: unexpected recovery of your wallet",
					fmt.Sprintf("A recovery of your wallet %v was started that you did not request through Trovo. It takes effect on %v UTC. Open the app now and cancel it.", w.Alias, after.Format("2 Jan 2006 15:04")), gc)
			}
		}
	}
	gc.DB.Model(&cursor).Update("block", to)
}

// notifyRecovery tells the user by push notification and email.
func notifyRecovery(user *userModels.User, title, body string, gc *sharedconfig.GlobalConfig) {
	user.SendPushMessage(title, body, "", map[string]string{"route": "accountRecovery"}, gc)
	if os.Getenv("ENABLE_EMAIL_NOTIFICATIONS") != "1" || os.Getenv("MAILGUN_PRIVATE_API_KEY") == "" || user.Email == "" {
		return
	}
	sender := os.Getenv("DEFAULT_MAIL_SENDER")
	if sender == "" {
		sender = os.Getenv("SUPPORT_EMAIL")
	}
	if sender == "" {
		sender = "support@trovotech.io"
	}
	mg := mailgun.NewMailgun(os.Getenv("MAILGUN_DOMAIN"), os.Getenv("MAILGUN_PRIVATE_API_KEY"))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if _, _, err := mg.Send(ctx, mg.NewMessage(sender, "Trovo: "+title, body, user.Email)); err != nil {
		log.Printf("[notifyRecovery] emailing %v: %v", user.Username, err)
	}
}
