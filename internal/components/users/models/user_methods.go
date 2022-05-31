package users

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"sort"
	dl "trovo-wallet-api/internal/dynamiclinks"
	tErrors "trovo-wallet-api/internal/errors"

	"context"
	"errors"
	"log"
	"os"
	"strings"
	"sync"
	"time"
	"trovo-wallet-api/internal/cache"
	"trovo-wallet-api/internal/network"

	"github.com/mailgun/mailgun-go/v4"
	"github.com/shopspring/decimal"
	"github.com/stellar/go/clients/horizonclient"
	"gorm.io/gorm"

	"github.com/stellar/go/protocols/horizon"
)

/////User Model convenience Methods

// GetSigners returns user signers
func (u *UserWallet) GetSigners(temp bool) (signers map[string]Signer) {
	account, _, err := u.GetBlockchainAccountDetail(temp)
	if err != nil {
		return signers
	}
	signers = make(map[string]Signer)
	for _, v := range account.Signers {
		signers[v.Key] = Signer{
			Weight:  int(v.Weight),
			Key:     v.Key,
			Type:    v.Type,
			Sponsor: v.Sponsor,
		}
	}
	return
}

// GetSignersWA returns user signers
func (u *User) GetSignersWA(account *horizon.Account) (signers map[string]Signer) {

	signers = make(map[string]Signer)
	for _, v := range account.Signers {
		signers[v.Key] = Signer{
			Weight:  int(v.Weight),
			Key:     v.Key,
			Type:    v.Type,
			Sponsor: v.Sponsor,
		}
	}
	return
}

//SignerIsValidWA checks if the signerKey is valid for this user public key
func (u *User) SignerIsValidWA(signerKey string, account *horizon.Account) bool {
	signer, ok := u.GetSignersWA(account)[signerKey]
	if !ok || signer.Weight < 1 {
		return false
	}

	return true
}

//SignerIsValid checks if the signerKey is valid for this user public key
func (u *UserWallet) SignerIsValid(signerKey string, temp bool) bool {
	signer, ok := u.GetSigners(temp)[signerKey]
	if !ok || signer.Weight < 1 {
		return false
	}

	return true
}

//SignerIsValid checks if the signerKey is valid for this user public key
func (u *User) SignerIsValid(signerKey string, temp bool) bool {
	for _, w := range u.UserWallets {
		if w.ID == w.Signer {
			signer, ok := w.GetSigners(temp)[signerKey]
			if !ok || signer.Weight < 1 {
				return false
			}

			return true
		}
	}
	return false
}

//GetBalance gets user wallet blockchain balance and return it as a map of assets  [code:issuer]Balance. Native key is [:]
func (u *UserWallet) GetBalance(db *gorm.DB, temp bool, dynamicLinkServiceUrlChan chan string, redisCache *cache.RedisCache) (balances map[string]Balance, err error) {
	balances = make(map[string]Balance)
	cacheKey := fmt.Sprintf("GetBalance_%s", u.ID)
	if temp {
		cacheKey = fmt.Sprintf("GetBalance_%s", *u.TempPublicKey)

	}
	{

		// search cache for balance
		ok, response := redisCache.GetCachedResult(cacheKey)

		if ok {
			log.Printf("GetBalance[%v], served from cache\n", cacheKey)
			balances = response.(map[string]Balance)
			return
		}

	}

	account, _, err := u.GetBlockchainAccountDetail(temp)
	if err != nil {
		return balances, err
	}

	var wg sync.WaitGroup
	var m sync.Mutex
	// var keys []string
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
			// availableBalance := availableBal.Truncate(7).String()
			qrCode := ""
			if !temp {

				p, e := dl.GeneratePaymentData(u.ID, v.Code, v.Issuer, "", "", dynamicLinkServiceUrlChan, redisCache)
				if e == nil {
					qrCode = p.QRCode
				}
			}
			balance := Balance{AssetIssuer: v.Issuer, AssetCode: v.Code,
				Amount: availableBalance, QRCode: qrCode}
			m.Lock()
			balances[v.Code+":"+v.Issuer] = balance
			m.Unlock()

		}(v)
		wg.Wait()
	}
	//save to cache
	redisCache.StoreResultToCache(cacheKey, balances, 0)
	return balances, nil
}

//GetSortedUserBalance gets user blockchain balance
func (u *UserWallet) GetSortedUserBalance(db *gorm.DB, temp bool, dynamicLinkServiceUrlChan chan string, redisCache *cache.RedisCache) (balances []Balance, err error) {

	//GetBalance
	unsortedBalances, err := u.GetBalance(db, temp, dynamicLinkServiceUrlChan, redisCache)
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

// GetAccountThresholds returns user signers
func (u *UserWallet) GetAccountThresholds(temp bool) (thresholds horizon.AccountThresholds) {
	account, _, err := u.GetBlockchainAccountDetail(temp)
	if err != nil {
		return thresholds
	}

	return account.Thresholds
}

//GetBlockchainAccountDetail fetches the bantu account information using public key
func (u *UserWallet) GetBlockchainAccountDetail(temp bool) (clientAccount horizon.Account, destinationAccountExists bool, err error) {
	client := network.GetBlockchainClient()
	var accountRequest horizonclient.AccountRequest
	if temp {
		//temp account
		if u.TempPublicKey != nil {
			//temp account has been generated
			accountRequest = horizonclient.AccountRequest{AccountID: *u.TempPublicKey}

		} else {
			//temp account not yet generated
			err = &tErrors.ErrorBlockchainAccountNotActivated{}
			return
		}

	} else {
		//real account
		accountRequest = horizonclient.AccountRequest{AccountID: u.ID}
	}

	clientAccount, err = client.AccountDetail(accountRequest)
	if err != nil {
		log.Println("[GetBlockchainAccountDetail]: ", err)
		if strings.Contains(err.Error(), "timeout") || strings.Contains(err.Error(), "handshake") || strings.Contains(err.Error(), "no such host") || strings.Contains(err.Error(), "timeout") || strings.Contains(err.Error(), "dial") {
			log.Printf("[GetBlockchainAccountDetail Network Failure]: %s\n", "Error Connecting to Expansion Service")
			return clientAccount, destinationAccountExists, &tErrors.ErrorTemporaryServerError{}
		} else {
			horizonException, ok := err.(*horizonclient.Error)

			if ok {

				log.Println("[BlockchainAccountProperties] error is known", horizonException.Problem.Status)

				if horizonException.Problem.Status == http.StatusNotFound {
					return clientAccount, false, &tErrors.ErrorBlockchainAccountNotActivated{}
				}

			}

		}
		return clientAccount, destinationAccountExists, &tErrors.ErrorTemporaryServerError{}
	}
	return clientAccount, true, nil
}

func (u *User) VerifyEmailOnMailgun() (validationResult mailgun.EmailVerification, blockEmail bool, err error) {
	// To use the /v4 version of validations define MG_URL in the environment
	// as `https://api.mailgun.net/v4` or set `v.SetAPIBase("https://api.mailgun.net/v4")`
	if os.Getenv("ENABLE_EMAIL_VALIDATION") == "1" {
		//check if email validation passes
		if len(os.Getenv("MAILGUN_VALIDATOR_API_KEY")) == 0 {
			log.Println("MAILGUN_VALIDATOR_API_KEY not set!")
			err = &tErrors.ErrorTemporaryServerError{}
			return
		}
		apiKey := os.Getenv("MAILGUN_VALIDATOR_API_KEY")

		// Create an instance of the Validator
		v := mailgun.NewEmailValidator(apiKey)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		validationResult, err = v.ValidateEmail(ctx, u.Email, true)

		if err != nil {
			log.Println("[VerifyEmailOnMailgun] mail validation request error:", err)
			validationResult, err = v.ValidateEmail(ctx, u.Email, true)
			if err != nil {
				log.Println("[VerifyEmailOnMailgun] mail validation request failed second time with error:", err)

				return validationResult, blockEmail, &tErrors.ErrorTemporaryServerError{}
			}
		}

		if validationResult.IsDisposableAddress || strings.Contains(validationResult.Reason, "unknown") || strings.Contains(validationResult.Reason, "risk") || strings.Contains(validationResult.Reason, "mail") || strings.Contains(validationResult.Reason, "domain") || strings.Contains(validationResult.Reason, "catch") || strings.Contains(validationResult.Reason, "long") || strings.Contains(validationResult.Reason, "no_mx") {
			log.Printf("[VerifyEmailOnMailgun] %v, user:%#v %#v\n", errors.New("failed mail validation"), *u, validationResult)
			blockEmail = true
			if os.Getenv("SHOW_MAIL_VALIDATION_RESULT") == "1" {
				log.Printf("[VerifyEmailOnMailgun] validation result: %#v\n", validationResult)
			}
			return validationResult, blockEmail, nil
		}

	}
	return
}

//GetBlockchainAccountDataKey fetches the bantu account information using public key
func (u *UserWallet) GetBlockchainAccountDataKey(temp bool, keys ...string) (dataValues map[string]string) {
	dataValues = make(map[string]string)
	account, _, err := u.GetBlockchainAccountDetail(temp)
	if err != nil {
		return
	}
	data, err := u.GetBlockchainAccountData(account)
	if err != nil {
		return
	}
	for _, key := range keys {
		d, ok := data[key]
		if !ok {
			continue
		}
		decData, err := base64.StdEncoding.DecodeString(d)

		if err != nil {
			continue
		} else {
			dataValues[key] = string(decData)
		}

	}

	return dataValues
}

func (u *UserWallet) GetBlockchainAccountData(clientAccount horizon.Account) (accountData map[string]string, err error) {

	if err != nil {
		return
	}

	return clientAccount.Data, nil
}

func (u *User) BuildPrimaryWallet() {
	tempKP, _ := network.TempAccountKeypair(u.PublicKey)
	var tempPK string
	if tempKP != nil {
		tempPK = tempKP.Address()
	}
	description := "Primary/Default wallet"
	userWallet := UserWallet{
		ID:            u.PublicKey,
		TempPublicKey: &tempPK,
		Description:   &description,
		Alias:         u.Username,
		Signer:        u.PublicKey,
		UserID:        u.ID,
	}
	u.UserWallets = append(u.UserWallets, userWallet)
}
func (u *User) BuildNewSubWallet(subWalletPublicKey, walletTag, walletDescription string) {
	{
		//check to ensure sub-wallet does not already exist
		for _, wallet := range u.UserWallets {
			if wallet.ID == subWalletPublicKey {
				return
			}
		}
	}
	tempKP, _ := network.TempAccountKeypair(subWalletPublicKey)
	var tempPK string
	if tempKP != nil {
		tempPK = tempKP.Address()
	}
	walletTag = strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(walletTag, "_", ""), ".", ""), " ", "")
	alias := fmt.Sprintf("%s_%s", u.Username, walletTag)
	userSubWallet := UserWallet{
		ID:            subWalletPublicKey,
		TempPublicKey: &tempPK,
		Tag:           &walletTag,
		Description:   &walletDescription,
		Alias:         alias,
		Signer:        u.PublicKey,
		UserID:        u.ID,
	}
	u.UserWallets = append(u.UserWallets, userSubWallet)
}
func (id UserWalletManagedAccessID) String() string {
	return string(id)
}
func (id UserWalletID) String() string {
	return string(id)
}

func (id UserWalletManagedAccessID) GetAccessAssignment(db *gorm.DB) (assignment *UserWalletManagedAccess, err error) {
	e := db.Where("id = ?", string(id)).First(assignment).Error
	if e != nil {
		if errors.Is(e, gorm.ErrRecordNotFound) {
			//no managed access was found
			err = &tErrors.CustomError{
				Param:      "id",
				Err:        "error-access-assignment-not-found",
				ErrMessage: "Access ID not found",
				Code:       404,
			}
			return
		}
		err = &tErrors.ErrorTemporaryServerError{}
	}
	return
}

func (id UserWalletID) GetWallet(db *gorm.DB) (wallet *UserWallet, err error) {
	e := db.Where("id = ?", string(id)).First(wallet).Error
	if e != nil {
		if errors.Is(e, gorm.ErrRecordNotFound) {
			//no wallet was found
			err = &tErrors.CustomError{
				Param:      "id",
				Err:        "error-wallet-not-found",
				ErrMessage: "Wallet not found",
				Code:       404,
			}
			return
		}
		err = &tErrors.ErrorTemporaryServerError{}
	}
	return
}

func (id UserWalletID) GetWalletOwner(db *gorm.DB) (walletOwner *User, err error) {
	e := db.Where("public_key = ?", string(id)).First(walletOwner).Error
	if e != nil {
		if errors.Is(e, gorm.ErrRecordNotFound) {
			//no wallet was found
			err = &tErrors.CustomError{
				Param:      "id",
				Err:        "error-account-not-found",
				ErrMessage: "Account not found",
				Code:       404,
			}
			return
		}
		err = &tErrors.ErrorTemporaryServerError{}
	}
	return
}

//Fetch3rdPartyWallets fetches all 3rd party wallets that the user is assigned to manage
func (u *User) Fetch3rdPartyWallets(db *gorm.DB) (thirdPartyWallets []ThirdPartyWalletAccess) {
	var walletPermissions []WalletAccess
	thirdPartyWallets = make([]ThirdPartyWalletAccess, 0)
	e := db.Where("username = ?", u.Username).Find(&walletPermissions).Error
	if e != nil {
		return
	}
	for _, assignedPermission := range walletPermissions {
		//Get the permission assignment
		managedAccess, err := UserWalletManagedAccessID(assignedPermission.UserWalletManagedAccessID).GetAccessAssignment(db)
		thirdPartyWallet := ThirdPartyWalletAccess{
			AccessLevel: assignedPermission.AccessLevel,
		}
		if err == nil {
			//use it to fetch wallet details
			wallet, err := UserWalletID(managedAccess.UserWalletID).GetWallet(db)
			if err == nil {
				thirdPartyWallet.PublicKey = wallet.ID
				thirdPartyWallet.WalletAlias = wallet.Alias
				if wallet.Description != nil {
					thirdPartyWallet.WalletDescription = *wallet.Description
				}
			}
			//use it to fetch wallet owner details
			owner, err := UserWalletID(managedAccess.UserWalletID).GetWalletOwner(db)
			if err == nil {
				thirdPartyWallet.Owner = owner.Username
			}
		}
		thirdPartyWallets = append(thirdPartyWallets, thirdPartyWallet)

	}

	return
}
