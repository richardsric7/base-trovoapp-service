package users

import (
	"log"
	"sort"
	"strings"
	"sync"
	"trovo-wallet-api/internal/cache"
	userDB "trovo-wallet-api/internal/components/users/db"
	userModels "trovo-wallet-api/internal/components/users/models"
	dl "trovo-wallet-api/internal/dynamiclinks"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/network"

	"github.com/shopspring/decimal"
	"github.com/stellar/go/clients/horizonclient"
	"gorm.io/gorm"

	"github.com/stellar/go/protocols/horizon"
)

//GetUserBalance gets user blockchain balance
func GetUserBalance(publicKey string, db *gorm.DB, temp bool, dynamicLinkServiceUrlChan chan string, redisCache *cache.RedisCache) (balances []userModels.Balance, err error) {
	unsortedBalances := make(map[string]userModels.Balance)
	account, err := GetBlockchainAccountDetail(publicKey)
	if err != nil {
		return balances, err
	}

	//get username for qrcode is exists
	var user userModels.User
	var errUser error
	if !temp {
		user, errUser = userDB.GetUserInfo(publicKey, db)
	}

	var wg sync.WaitGroup
	var m sync.Mutex
	var keys []string
	for _, v := range account.Balances {
		wg.Add(1)
		go func(v horizon.Balance) {
			defer wg.Done()

			amount, _ := decimal.NewFromString(v.Balance)
			if (temp && (amount.IsZero())) || (v.Code == "" && temp) {

				return
			}

			// buyingLiabilities, _ := decimal.NewFromString(v.BuyingLiabilities)
			sellingLiabilities, _ := decimal.NewFromString(v.SellingLiabilities)
			availableBalance := amount.Sub(sellingLiabilities)
			// availableBalance := decimal.NewFromFloat(availableBalFloat).Truncate(7).String()
			// usdValue := decimal.NewFromFloat(usdPrice * availableBalFloat).Truncate(2).String()
			qrCode := ""
			if !temp {
				key := publicKey
				if errUser == nil {
					key = user.Username
				}
				p, e := dl.GeneratePaymentData(key, v.Code, v.Issuer, "", "", dynamicLinkServiceUrlChan, redisCache)
				if e == nil {
					qrCode = p.QRCode
				}
			}

			balance := userModels.Balance{AssetIssuer: v.Issuer, AssetCode: v.Code,
				Amount: availableBalance, QRCode: qrCode}
			m.Lock()
			keys = append(keys, v.Code+":"+v.Issuer)
			unsortedBalances[v.Code+":"+v.Issuer] = balance
			m.Unlock()

		}(v)
		wg.Wait()
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
