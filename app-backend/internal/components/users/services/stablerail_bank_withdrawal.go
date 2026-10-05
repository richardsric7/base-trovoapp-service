package users

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
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
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// OperationBankWithdrawal sends cNGN from a user's wallet to their Strails
// smart wallet, from where Strails pays it to their bank account.
const OperationBankWithdrawal = "BANK WITHDRAWAL"

// maxOfframpRequestAttempts is how often Strails is asked to pay out a
// deposited withdrawal before it is left to support.
const maxOfframpRequestAttempts = 10

var accountNumberPattern = regexp.MustCompile(`^[0-9]{10}$`)

// StablerailEnabled reports whether Stablerail is configured and switched on.
func StablerailEnabled(gc *sharedconfig.GlobalConfig) bool {
	var config userModels.StablerailConfig
	gc.DB.First(&config)
	return len(config.ApiKey) > 0 && config.EnableStablerail == 1
}

// stablerailPost POSTs payload to a Stablerail endpoint and decodes the
// reply into out. A non-2xx reply is an error carrying Stablerail's message.
func stablerailPost(path string, payload, out interface{}, gc *sharedconfig.GlobalConfig) error {
	var config userModels.StablerailConfig
	gc.DB.First(&config)
	if len(config.ApiKey) == 0 {
		return fmt.Errorf("no stablerail config found: %v", 404)
	}
	if config.EnableStablerail == 0 {
		return fmt.Errorf("stablerail not enabled: %v", 202)
	}
	baseUrl := config.BaseUrl
	if len(baseUrl) == 0 {
		baseUrl = "https://beta.stablesrail.io/v1"
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, strings.TrimRight(baseUrl, "/")+path, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("x-api-key", config.ApiKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err := (&http.Client{Timeout: 15 * time.Second}).Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		var e struct {
			Message string `json:"message"`
		}
		_ = json.Unmarshal(data, &e)
		return fmt.Errorf("stablerail %v: %d %v", path, resp.StatusCode, e.Message)
	}
	return json.Unmarshal(data, out)
}

// minimumBankWithdrawal is the smallest cNGN amount that can be withdrawn
// to a bank (STABLERAIL_MIN_WITHDRAWAL, default 1000).
func minimumBankWithdrawal() decimal.Decimal {
	if d, err := decimal.NewFromString(strings.TrimSpace(os.Getenv("STABLERAIL_MIN_WITHDRAWAL"))); err == nil && d.IsPositive() {
		return d
	}
	return decimal.NewFromInt(1000)
}

// GetStablerailProfile tells the apps whether the user can deposit from and
// withdraw to a bank.
func GetStablerailProfile(user *userModels.User, gc *sharedconfig.GlobalConfig) userModels.StablerailProfile {
	p := userModels.StablerailProfile{Enabled: StablerailEnabled(gc), MinimumWithdrawal: minimumBankWithdrawal().String()}
	if su := GetStablerailUser(user.Username, gc); len(su.ID) > 0 {
		p.Onboarded, p.OnboardingStatus = true, "completed"
		return p
	}
	var r userModels.StablerailRequest
	if gc.DB.Where("trovo_username = ? AND request_type = ?", user.Username, "Onboarding").Order("created_at DESC").First(&r).Error == nil {
		p.OnboardingStatus = r.Status
	}
	return p
}

// stablerailDepositWallet is the user's Strails smart wallet on Base.
func stablerailDepositWallet(su *userModels.StablerailUser, gc *sharedconfig.GlobalConfig) (common.Address, error) {
	if common.IsHexAddress(su.EvmWallet) {
		return common.HexToAddress(su.EvmWallet), nil
	}
	var res userModels.StablerailUserDetailsResponse
	if err := stablerailPost("/getuserdetails", map[string]string{"userId": su.ID}, &res, gc); err != nil {
		log.Printf("[stablerailDepositWallet] %v: %v", su.TrovoUsername, err)
		return common.Address{}, &tErrors.ErrorTemporaryServerError{}
	}
	w := strings.TrimSpace(res.Data.WalletDetails.EvmWallet)
	if !common.IsHexAddress(w) {
		return common.Address{}, &tErrors.CustomError{Param: "account", Err: "error-fiat-wallet-unavailable", ErrMessage: "Your bank withdrawal wallet is not ready yet. Please try again later.", Code: http.StatusServiceUnavailable}
	}
	gc.DB.Model(&userModels.StablerailUser{}).Where("id = ?", su.ID).Update("evm_wallet", w)
	return common.HexToAddress(w), nil
}

// bankWithdrawalContext is what a bank withdrawal operation needs once mined.
type bankWithdrawalContext struct {
	Withdrawal userModels.StablerailOfframp `json:"withdrawal"`
}

// WithdrawToBank withdraws cNGN from wallet to a Nigerian bank account. Like
// Pay it runs in two steps: without a signature it builds the transfer of
// the cNGN to the user's Strails wallet as a wallet operation and returns it
// (req.Transaction); with the signature it submits it. Once the transfer is
// mined Strails is asked to pay the bank account (see UpdateOfframpStatuses).
// Shared wallets with approvers get an approval request on commit.
func WithdrawToBank(signerUser, walletOwner *userModels.User, wallet *userModels.UserWallet, req *userModels.BankWithdrawalRequest, gc *sharedconfig.GlobalConfig) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	req.Messages = make([]string, 0)
	if wallet.SharedAccessEnabled == 1 && wallet.NumberOfApprovalsNeeded > 0 {
		req.Multiparty = 1
	}
	if !StablerailEnabled(gc) {
		return &tErrors.CustomError{Param: "account", Err: "error-fiat-unavailable", ErrMessage: "Bank withdrawals are not available at this time.", Code: http.StatusServiceUnavailable}
	}

	// step 2: the user signed the operation built in step 1
	if req.Multiparty == 0 && len(req.TransactionSignature) > 0 && len(req.Transaction) > 0 {
		rec, p, err := LoadWalletOperation(req.Transaction, wallet.ID, OperationBankWithdrawal, gc)
		if err != nil {
			return err
		}
		var c bankWithdrawalContext
		if rec.Context == nil || json.Unmarshal([]byte(*rec.Context), &c) != nil {
			return &tErrors.ErrorTemporaryServerError{}
		}
		hash, err := SignSingleOwnerOperation(ctx, rec, p, signerUser.PrimarySigner, req.TransactionSignature, gc)
		if err != nil {
			return err
		}
		req.TransactionID, req.WithdrawalID = hash, c.Withdrawal.ID
		w := c.Withdrawal
		w.UserOpHash, w.Status = &hash, userModels.OfframpDepositing
		// the mined hook may have recorded it already
		if e := gc.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&w).Error; e != nil {
			log.Printf("[WithdrawToBank] saving withdrawal %+v: %v", w, e)
			gc.LogDiscordFailedRequest(fmt.Sprintf("[WithdrawToBank] bank withdrawal of %v submitted (%v) but not recorded: %v", wallet.ID, hash, e))
		}
		walletOwner.InvalidateUserCache(gc)
		return nil
	}
	// a shared wallet's initiator commits the operation built in step 1
	if req.Multiparty == 1 && req.Commit == 1 && len(req.Transaction) > 0 {
		rec, _, err := LoadWalletOperation(req.Transaction, wallet.ID, OperationBankWithdrawal, gc)
		if err != nil {
			return err
		}
		var c bankWithdrawalContext
		if rec.Context == nil || json.Unmarshal([]byte(*rec.Context), &c) != nil {
			return &tErrors.ErrorTemporaryServerError{}
		}
		return createBankWithdrawalApproval(signerUser, wallet, req, &c.Withdrawal, gc)
	}

	// step 1: build
	su := GetStablerailUser(walletOwner.Username, gc)
	if len(su.ID) == 0 {
		return &tErrors.CustomError{Param: "account", Err: "error-fiat-not-onboarded", ErrMessage: "Verify your BVN for fiat transactions before withdrawing to a bank.", Code: http.StatusForbidden}
	}
	amount, e := decimal.NewFromString(strings.TrimSpace(req.Amount))
	if e != nil || !amount.IsPositive() || amount.Exponent() < -2 {
		return &tErrors.CustomError{Param: "amount", Err: "error-invalid-amount", ErrMessage: "Enter an amount in Naira with at most 2 decimal places."}
	}
	if min := minimumBankWithdrawal(); amount.LessThan(min) {
		return &tErrors.CustomError{Param: "amount", Err: "error-invalid-amount", ErrMessage: fmt.Sprintf("The smallest bank withdrawal is NGN %v.", min)}
	}
	req.AccountNumber = strings.TrimSpace(req.AccountNumber)
	if !accountNumberPattern.MatchString(req.AccountNumber) {
		return &tErrors.CustomError{Param: "accountNumber", Err: "error-invalid-account-number", ErrMessage: "Enter the 10-digit bank account number."}
	}
	var bank userModels.StablerailBank
	if gc.DB.Where("bank_code = ?", strings.TrimSpace(req.BankCode)).First(&bank).Error != nil {
		return &tErrors.CustomError{Param: "bankCode", Err: "error-invalid-bank", ErrMessage: "Choose a bank from the list."}
	}
	req.BankName = bank.BankName
	currency := GetTokenizationCurrencyByCode("CNGN", gc.DB)
	if !common.IsHexAddress(currency.ContractAddress) {
		return &tErrors.CustomError{Param: "account", Err: "error-fiat-unavailable", ErrMessage: "Bank withdrawals are not available at this time.", Code: http.StatusServiceUnavailable}
	}
	token := common.HexToAddress(currency.ContractAddress)
	client := gc.BantuExpansionClient
	decimals, err := offerbook.Decimals(ctx, gc.DB, client, token)
	if err != nil {
		return &tErrors.ErrorTemporaryServerError{}
	}
	_, _, _, balance, _, err := network.BlockchainAccountProperties(client, wallet.ID, basetxn.CreditAsset{Code: "CNGN", Issuer: token.Hex()})
	if err != nil {
		return err
	}
	if balance.LessThan(amount) {
		return &tErrors.ErrorUnderfundedAccount{Detail: fmt.Sprintf("Not enough funds. This wallet holds %v CNGN.", balance)}
	}
	deposit, err := stablerailDepositWallet(&su, gc)
	if err != nil {
		return err
	}
	draft := userModels.StablerailOfframp{
		ID: uuid.NewString(), BankCode: bank.BankCode, BankName: bank.BankName, AccountNumber: req.AccountNumber,
		BaseAmount: amount.InexactFloat64(), Ticker: "CNGN", TrovoUsername: walletOwner.Username,
		WalletAddress: wallet.ID, DepositAddress: deposit.Hex(),
	}
	calls := []aa.Call{aa.ERC20Transfer(token, deposit, amount.Shift(int32(decimals)).Truncate(0).BigInt())}
	validity := time.Duration(0)
	if req.Multiparty == 1 {
		validity = sharedWalletOperationValidity()
	}
	op, err := PrepareWalletOperation(ctx, OperationBankWithdrawal, signerUser, walletOwner, wallet, calls, validity, bankWithdrawalContext{Withdrawal: draft}, gc)
	if err != nil {
		return err
	}
	req.Transaction, req.WithdrawalID = op.Transaction, draft.ID
	req.Messages = append(req.Messages, fmt.Sprintf("NGN %v will be paid to %v %v, less the provider's fee. Payment usually arrives within minutes.", amount, bank.BankName, req.AccountNumber))
	req.Messages = append(req.Messages, op.Messages()...)
	req.SignatureRequired = 1
	if req.Multiparty == 1 && req.Commit == 1 {
		return createBankWithdrawalApproval(signerUser, wallet, req, &draft, gc)
	}
	return nil
}

func createBankWithdrawalApproval(signerUser *userModels.User, wallet *userModels.UserWallet, req *userModels.BankWithdrawalRequest, w *userModels.StablerailOfframp, gc *sharedconfig.GlobalConfig) error {
	req.TransactionID = "PENDING_AUTH"
	description := fmt.Sprintf("Bank withdrawal of NGN %v to %v %v", decimal.NewFromFloat(w.BaseAmount), w.BankName, w.AccountNumber)
	if len(req.Messages) > 0 {
		description = fmt.Sprintf("%s\nMessages: %v", description, strings.Join(req.Messages, "\n"))
	}
	infoBytes, _ := json.Marshal(*req)
	info := string(infoBytes)
	pendingAuth := userModels.PendingAuth{
		ID:                     uuid.NewString(),
		Initiator:              signerUser.Username,
		InitiatorSignerAddress: signerUser.PrimarySigner,
		WalletAddress:          wallet.ID,
		TransactionType:        OperationBankWithdrawal,
		Description:            description,
		TransactionSource:      wallet.ID,
		ApprovalsNeeded:        wallet.NumberOfApprovalsNeeded,
		TransactionXdr:         req.Transaction,
		TransactionInfoStr:     &info,
	}
	if err := gc.DB.Omit(clause.Associations).Create(&pendingAuth).Error; err != nil {
		log.Printf("[WithdrawToBank] saving approval request for %v: %v", wallet.ID, err)
		return &tErrors.ErrorTemporaryServerError{}
	}
	return nil
}

// recordBankWithdrawalDeposit marks a withdrawal deposited once its transfer
// to the Strails wallet is mined, and asks Strails to pay it out.
func recordBankWithdrawalDeposit(op userModels.WalletOperation, txHash common.Hash, gc *sharedconfig.GlobalConfig) {
	var c bankWithdrawalContext
	if op.Context == nil || json.Unmarshal([]byte(*op.Context), &c) != nil || c.Withdrawal.ID == "" {
		return
	}
	w := c.Withdrawal
	w.UserOpHash, w.TxHash, w.Status = op.UserOpHash, txHash.Hex(), userModels.OfframpDeposited
	gc.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&w)
	gc.DB.Model(&userModels.StablerailOfframp{}).Where("id = ? AND status = ?", w.ID, userModels.OfframpDepositing).
		Updates(map[string]interface{}{"status": userModels.OfframpDeposited, "tx_hash": txHash.Hex(), "user_op_hash": op.UserOpHash})
	go requestOfframpPayout(w.ID, gc)
}

func init() {
	minedHooks[OperationBankWithdrawal] = recordBankWithdrawalDeposit
}

// requestOfframpPayout asks Strails to pay out a deposited withdrawal. The
// row is claimed first (DEPOSITED or REQUEST_FAILED -> REQUESTING) so only
// one instance asks; Strails itself refuses a second payout of the same
// deposit (the wallet no longer holds it).
func requestOfframpPayout(id string, gc *sharedconfig.GlobalConfig) {
	res := gc.DB.Model(&userModels.StablerailOfframp{}).
		Where("id = ? AND status IN ? AND attempts < ?", id, []string{userModels.OfframpDeposited, userModels.OfframpRequestFailed}, maxOfframpRequestAttempts).
		Updates(map[string]interface{}{"status": userModels.OfframpRequesting, "attempts": gorm.Expr("attempts + 1")})
	if res.Error != nil || res.RowsAffected == 0 {
		return
	}
	var w userModels.StablerailOfframp
	if gc.DB.Where("id = ?", id).First(&w).Error != nil {
		return
	}
	su := GetStablerailUser(w.TrovoUsername, gc)
	r, err := StableRailInitiateOfframp(userModels.StablerailOfframpRequest{
		UserID: su.ID, Amount: w.BaseAmount, AccountNumber: w.AccountNumber, BankCode: w.BankCode, Ticker: nonEmpty(w.Ticker, "CNGN"),
	}, gc)
	if err == nil && (r == nil || r.Data.RequestID == "") {
		err = fmt.Errorf("no request id: %v", func() string {
			if r == nil {
				return "empty reply"
			}
			return r.Message
		}())
	}
	if err != nil {
		log.Printf("[requestOfframpPayout] %v (attempt %d): %v", id, w.Attempts, err)
		gc.DB.Model(&userModels.StablerailOfframp{}).Where("id = ?", id).Updates(map[string]interface{}{"status": userModels.OfframpRequestFailed, "error": err.Error()})
		if w.Attempts >= maxOfframpRequestAttempts {
			gc.LogDiscordFailedRequest(fmt.Sprintf("[requestOfframpPayout] bank withdrawal %v of %v (NGN %v, tx %v) could not be requested from Stablerail after %d attempts: %v. The cNGN is in the user's Stablerail wallet %v.", id, w.TrovoUsername, w.BaseAmount, w.TxHash, w.Attempts, err, w.DepositAddress))
			if u, e := userModels.Username(w.TrovoUsername).GetSimpleUser(gc.DB, gc); e == nil {
				u.SendPushMessage("Bank withdrawal delayed", fmt.Sprintf("Your withdrawal of NGN %v could not be sent to your bank yet. Our support team has been notified and will complete it.", w.BaseAmount), "", map[string]string{"route": "basicTransactionHistory"}, gc)
			}
		}
		return
	}
	gc.DB.Model(&userModels.StablerailOfframp{}).Where("id = ?", id).
		Updates(map[string]interface{}{"request_id": r.Data.RequestID, "status": nonEmpty(r.Data.Status, "pending"), "error": ""})
}

// ListBankWithdrawals lists a user's bank withdrawals, newest first.
func ListBankWithdrawals(username string, gc *sharedconfig.GlobalConfig) ([]userModels.StablerailOfframp, error) {
	out := make([]userModels.StablerailOfframp, 0)
	if err := gc.DB.Where("trovo_username = ?", username).Order("created_at DESC").Limit(100).Find(&out).Error; err != nil {
		return nil, &tErrors.ErrorTemporaryServerError{}
	}
	return out, nil
}
