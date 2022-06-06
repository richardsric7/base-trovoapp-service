package users

import (
	"log"
	"sort"
	"strings"
	userDB "trovo-wallet-api/internal/components/users/db"
	userModels "trovo-wallet-api/internal/components/users/models"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/network"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/stellar/go/clients/horizonclient"

	"github.com/stellar/go/protocols/horizon"
)

//GetSortedUserBalance gets user blockchain balance
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

//GetBlockchainAccountDetail fetches the bantu account information using public key
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
}
