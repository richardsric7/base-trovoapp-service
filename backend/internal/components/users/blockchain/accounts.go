package users

import (
	"log"
	"sort"
	"strings"
	"time"
	userDB "trovo-wallet-api/internal/components/users/db"
	userModels "trovo-wallet-api/internal/components/users/models"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/network"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/stellar/go/clients/horizonclient"

	"github.com/stellar/go/protocols/horizon"
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
		signers[v.Key] = userModels.Signer{
			Weight:  int(v.Weight),
			Key:     v.Key,
			Type:    v.Type,
			Sponsor: v.Sponsor,
		}
	}
	return
}

func SignerIsValid(publicKey, signerKey string) bool {
	signer, ok := GetUserSigners(publicKey)[signerKey]
	if !ok || signer.Weight < 1 {
		return false
	}

	return true
}

// GetUserAccountThresholds returns user signers
func GetUserAccountThresholds(publicKey string) (thresholds horizon.AccountThresholds) {
	account, err := GetBlockchainAccountDetail(publicKey)
	if err != nil {
		return thresholds
	}

	return account.Thresholds
}

// GetBlockchainAccountDetail fetches the bantu account information using public key
func GetBlockchainAccountDetail(publicKey string) (clientAccount horizon.Account, err error) {
	client := network.GetBlockchainClient()
	accountRequest := horizonclient.AccountRequest{AccountID: publicKey}
	clientAccount, err = client.AccountDetail(accountRequest)
	if err != nil {
		log.Println("[GetBlockchainAccountDetail]: ", err)
		if strings.Contains(err.Error(), "timeout") || strings.Contains(err.Error(), "handshake") || strings.Contains(err.Error(), "no such host") || strings.Contains(err.Error(), "timeout") || strings.Contains(err.Error(), "dial") {
			log.Printf("[GetBlockchainAccountDetail Network Failure]: %s\n", "Error Connecting to Expansion Service")
			return clientAccount, &tErrors.ErrorTemporaryServerError{}
		} else if strings.Contains(strings.ToLower(err.Error()), "missing") {

			err = &tErrors.ErrorBlockchainAccountNotActivated{}
		} else {

			err = &tErrors.ErrorTemporaryServerError{}
		}
		return horizon.Account{}, err
	}
	return clientAccount, nil
} //GetBlockchainAccountDetail fetches the bantu account information using public key

// GetBlockchainAssets fetches the blockchain asset information using public key
func GetBlockchainAssets(issuerPublicKey string) (assetsPage horizon.AssetsPage, err error) {
	client := network.GetBlockchainClient()
	assetRequest := horizonclient.AssetRequest{ForAssetIssuer: issuerPublicKey, Limit: 200}
	assetsPage, err = client.Assets(assetRequest)
	if err != nil {
		log.Println("[GetBlockchainAssets]: ", err)
		if strings.Contains(err.Error(), "timeout") || strings.Contains(err.Error(), "handshake") || strings.Contains(err.Error(), "no such host") || strings.Contains(err.Error(), "timeout") || strings.Contains(err.Error(), "dial") {
			log.Printf("[GetBlockchainAssets Network Failure]: %s\n", "Error Connecting to Blockchain API Service")
			return assetsPage, &tErrors.ErrorTemporaryServerError{}
		} else if strings.Contains(strings.ToLower(err.Error()), "missing") {
			err = &tErrors.ErrorBlockchainAccountNotActivated{}
		} else {

			err = &tErrors.ErrorTemporaryServerError{}
		}

		return
	}
	return assetsPage, nil
}

// BlockchainAssetIssuedByIssuer fetches the blockchain asset information using public key
func BlockchainAssetIssuedByIssuer(issuerPublicKey, assetCode string) bool {
	client := network.GetBlockchainClient()
	assetRequest := horizonclient.AssetRequest{ForAssetIssuer: issuerPublicKey, ForAssetCode: assetCode}
	assetsPage, err := client.Assets(assetRequest)
	if err != nil {
		return false
	}
	return len(assetsPage.Embedded.Records) > 0

}

// GetBlockchainAssetsIssuedByIssuer returns blockchain assets issued by the issuer
func GetBlockchainAssetsIssuedByIssuer(issuerPublicKey string) (issuedAssets map[string]horizon.AssetStat) {
	issuedAssets = make(map[string]horizon.AssetStat, 0)
	var err error
	assetPage, err := GetBlockchainAssets(issuerPublicKey)
	if err != nil {
		return
	}
	//iterate through assetPage
	for _, a := range assetPage.Embedded.Records {
		issuedAssets[a.Code] = a
	}
	return
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
