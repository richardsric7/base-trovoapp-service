package users

import (
	"sort"
	"time"
	"trovo-wallet-api/internal/basetxn"
	userDB "trovo-wallet-api/internal/components/users/db"
	userModels "trovo-wallet-api/internal/components/users/models"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/network"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/shopspring/decimal"
)

// GetSortedUserBalance gets user blockchain balance
func GetSortedUserBalance(publicKey string, gc *sharedconfig.GlobalConfig) (balances []userModels.Balance, err error) {

	var userWallet userModels.UserWallet
	userWallet, temp, err := userDB.GetWallet(publicKey, gc.DB)
	if err != nil {
		return
	}

	//GetBalance
	unsortedBalances, err := userWallet.GetBalance(temp, gc)
	if err != nil {
		return
	}

	var keys []string
	for key := range unsortedBalances {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	// To perform the ops of restoration from sorted keys
	for _, k := range keys {
		balance := unsortedBalances[k]
		balances = append(balances, balance)

	}

	return balances, nil
}

// GetUserSigners returns user signers
func GetUserSigners(publicKey string) (signers map[string]userModels.Signer) {
	account, err := GetBlockchainAccountDetail(publicKey)
	if err != nil {
		return signers
	}
	signers = make(map[string]userModels.Signer)
	for _, v := range account.Signers {
		signers[v.Key] = v
	}
	return
}

// SignerIsValid checks if signerKey is the account publicKey's own key.
// Base wallets are plain EOAs with exactly one key (see
// userModels.AccountDetail's doc) - there is no on-chain weighted-signer
// concept to check against as there was on Stellar, so "is this key
// valid to act for this account" is just equality (or, for a subwallet/
// recovery-key check, whatever the caller passed as signerKey - callers
// in subwallets.go/account_recovery.go/shared_access.go pass a
// deterministically-derived recovery/market-making/bulk-payment address
// as signerKey precisely to check "does this wallet's own key match one
// of those derived addresses").
func SignerIsValid(publicKey, signerKey string) bool {
	if network.IsAccountSigner(publicKey, signerKey) {
		return true
	}
	account, err := GetBlockchainAccountDetail(publicKey)
	if err != nil {
		return false
	}
	_, ok := account.Signers[signerKey]
	return ok
}

// GetUserAccountThresholds returns user signers
func GetUserAccountThresholds(publicKey string) (thresholds userModels.Thresholds) {
	account, err := GetBlockchainAccountDetail(publicKey)
	if err != nil {
		return thresholds
	}

	return account.Thresholds
}

type AccountDetailResult struct {
	Signers    map[string]userModels.Signer
	Thresholds userModels.Thresholds
}

// GetBlockchainAccountDetail fetches the bantu account information using public key
func GetBlockchainAccountDetail(publicKey string) (result AccountDetailResult, err error) {
	client := network.GetBlockchainClient()
	exists, _, _, _, _, err := network.BlockchainAccountProperties(client, publicKey, basetxn.NativeAsset{})
	if err != nil {
		return result, &tErrors.ErrorTemporaryServerError{}
	}
	if !exists {
		return result, &tErrors.ErrorBlockchainAccountNotActivated{}
	}
	result.Signers = map[string]userModels.Signer{
		publicKey: {Key: publicKey, Weight: 1, Type: "secp256k1_public_key"},
	}
	return result, nil
}

// GetNativeBalance fetches publicKey's native Base balance - a
// convenience wrapper for the common "check funding before continuing"
// pattern, replacing the original's iteration over
// horizon.Account.Balances looking for the code=="" entry.
func GetNativeBalance(publicKey string) (decimal.Decimal, error) {
	client := network.GetBlockchainClient()
	exists, _, nativeBalance, _, _, err := network.BlockchainAccountProperties(client, publicKey, basetxn.NativeAsset{})
	if err != nil {
		return decimal.Zero, &tErrors.ErrorTemporaryServerError{}
	}
	if !exists {
		return decimal.Zero, &tErrors.ErrorBlockchainAccountNotActivated{}
	}
	return nativeBalance, nil
}

// BlockchainAssetIssuedByIssuer fetches the blockchain asset information using public key
func BlockchainAssetIssuedByIssuer(issuerPublicKey, assetCode string) bool {
	db := network.DB()
	if db == nil {
		return false
	}
	var count int64
	db.Table("curated_assets").Where("asset_issuer = ? AND asset_code = ?", issuerPublicKey, assetCode).Count(&count)
	if count > 0 {
		return true
	}
	db.Table("tokenized_assets").Where("issuing_wallet_public_key = ? AND asset_code = ?", issuerPublicKey, assetCode).Count(&count)
	return count > 0
}

// BlockchainAssetIssuedByIssuer fetches the blockchain asset information using public key
func BlockchainAssetLastPaymentSource(toPublicKey, assetCode, issuerPublicKey string, gc *sharedconfig.GlobalConfig) string {
	type PaymentHistory struct {
		ID              string
		TransactionType string    `gorm:"index:idx_payment_history_unique_key,unique"`
		TransactionDate time.Time `json:"transactionDate" gorm:"index:idx_payment_history_tx_time"`
		From            *string   `json:"from" gorm:"size:150;index:idx_payment_history_from;null"` //trovoWallet alias and name
		FromPublicKey   string    `json:"fromPublicKey" gorm:"size:150;index:idx_payment_history_from_pk;not null;index:idx_payment_history_unique_key,unique;index:idx_payment_history_unique_key,unique"`
		To              *string   `json:"to" gorm:"size:56;index:idx_payment_history_to;null"` //trovoWallet alias and name
		ToPublicKey     string    `json:"toPublicKey" gorm:"size:56;index:idx_payment_history_to_pk;not null;index:idx_payment_history_unique_key,unique"`
		Memo            *string   `json:"memo" gorm:"size:28;null"`
		AssetIssuer     *string   `json:"assetIssuer" gorm:"size:56;null;"`
		AssetCode       string    `json:"assetCode" gorm:"size:12;not null;index:idx_payment_history_unique_key,unique"`
		Amount          string    `json:"amount" gorm:"index:idx_payment_history_unique_key,unique"`
		TransactionID   string    `json:"transactionId" gorm:"size:70;not null;index:idx_payment_history_txid;index:idx_payment_history_unique_key,unique"`
		PT              string    `json:"-" gorm:"size:70;not null;index:idx_payment_history_unique_key,unique;"`
	}

	var ph PaymentHistory
	e := gc.DB.Order("transaction_date DESC").Where("to_public_key = ? AND asset_code = ? AND (CASE WHEN asset_issuer IS NULL THEN '' ELSE asset_issuer END) = ?", toPublicKey, assetCode, issuerPublicKey).First(&ph).Error
	if e == nil {
		return ph.FromPublicKey
	}
	return ""
}
