package users

import (
	"encoding/base64"
	tErrors "trovo-wallet-api/internal/errors"

	"context"
	"errors"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
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
	account, err := u.GetBlockchainAccountDetail(temp)
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

//GetBalance gets user blockchain balance and return it as a map of assets  [code:issuer]Balance. Naitve key is [:]
func (u *UserWallet) GetBalance(publicKey string, db *gorm.DB, temp bool) (balances map[string]Balance, err error) {
	balances = make(map[string]Balance)
	account, err := u.GetBlockchainAccountDetail(temp)
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
			amount, _ := strconv.ParseFloat(v.Balance, 64)
			if (temp && (amount == 0)) || (v.Code == "" && temp) {
				// continue

				return
			}

			// askPrice, _ := GetDollarAskPrice(v.Code, v.Issuer)
			// usdPrice, _ := strconv.ParseFloat(askPrice, 64)

			// buyingLiabilities, _ := strconv.ParseFloat(v.BuyingLiabilities, 64)
			sellingLiabilities, _ := strconv.ParseFloat(v.SellingLiabilities, 64)
			availableBalFloat := (amount - sellingLiabilities)
			availableBalance := decimal.NewFromFloat(availableBalFloat).Truncate(7).String()
			// usdValue := decimal.NewFromFloat(usdPrice * availableBalFloat).Truncate(2).String()
			qrCode := ""
			// if !temp {

			// 	p, e := merchantServices.GeneratePaymentData(u.Username, v.Code, v.Issuer, "", "")
			// 	if e == nil {
			// 		qrCode = p.QRCode
			// 	}
			// }
			balance := Balance{AssetIssuer: v.Issuer, AssetCode: v.Code,
				Amount: availableBalance, QRCode: qrCode}
			m.Lock()
			balances[v.Code+":"+v.Issuer] = balance
			m.Unlock()

		}(v)
		wg.Wait()
	}

	return balances, nil
}

// GetAccountThresholds returns user signers
func (u *UserWallet) GetAccountThresholds(temp bool) (thresholds horizon.AccountThresholds) {
	account, err := u.GetBlockchainAccountDetail(temp)
	if err != nil {
		return thresholds
	}

	return account.Thresholds
}

// // GetAccountThresholds returns user signers
// func (u *User) PublicKeyBanned(db *gorm.DB, temp bool) (publicKeyBanned bool) {

// 	e := db.Where("public_key = ?", u.PublicKey).First(&BannedPublicKey{}).Error
// 	return e == nil
// }

//GetBlockchainAccountDetail fetches the bantu account information using public key
func (u *UserWallet) GetBlockchainAccountDetail(temp bool) (clientAccount horizon.Account, err error) {
	client := network.GetBlockchainClient()
	var accountRequest horizonclient.AccountRequest
	if temp {
		if u.TempPublicKey != nil {
			accountRequest = horizonclient.AccountRequest{AccountID: *u.TempPublicKey}

		} else {
			err = &tErrors.ErrorBlockchainAccountNotActivated{}
			return
		}

	} else {
		accountRequest = horizonclient.AccountRequest{AccountID: u.ID}
	}

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
	account, _ := u.GetBlockchainAccountDetail(temp)
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
