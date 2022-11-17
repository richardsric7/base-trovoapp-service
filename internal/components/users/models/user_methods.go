package users

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"
	"time"
	blockchain "trovo-wallet-api/internal/components/assets/blockchain"
	assets "trovo-wallet-api/internal/components/assets/models"
	dl "trovo-wallet-api/internal/dynamiclinks"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/network"
	pns "trovo-wallet-api/internal/pns"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/mailgun/mailgun-go/v4"
	"github.com/shopspring/decimal"
	"github.com/stellar/go/clients/horizonclient"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

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

// GetSignersWA returns user signers
func (u *UserWallet) GetSignersWA(account *horizon.Account) (signers map[string]Signer) {

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

// SignerIsValidWA checks if the signerKey is valid for this user public key
func (u *User) SignerIsValidWA(signerKey string, account *horizon.Account) bool {
	signer, ok := u.GetSignersWA(account)[signerKey]
	if !ok || signer.Weight < 1 {
		return false
	}

	return true
}

// SignerIsValidWA checks if the signerKey is valid for this user public key
func (u *UserWallet) SignerIsValidWA(signerKey string, account *horizon.Account) bool {
	signer, ok := u.GetSignersWA(account)[signerKey]
	if !ok || signer.Weight < 1 {
		return false
	}

	return true
}

// SignerIsValid checks if the signerKey is valid for this user public key
func (u *UserWallet) SignerIsValid(signerKey string, temp bool) bool {
	signer, ok := u.GetSigners(temp)[signerKey]
	if !ok || signer.Weight < 1 {
		return false
	}

	return true
}

// SignerIsValid checks if the signerKey is valid for this user public key
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

// GetBalance gets user wallet blockchain balance and return it as a map of assets  [code:issuer]Balance. Native key is [:]
func (u *UserWallet) GetBalance(temp bool, gc *sharedconfig.GlobalConfig) (balances map[string]Balance, err error) {
	balances = make(map[string]Balance)
	xbnUsdPrice, _ := blockchain.GetXBNDollarAskPrice(gc.DB)
	xbnNativePrice := "1"
	// log.Println("xbnUsdPrice", xbnUsdPrice)
	cacheKey := fmt.Sprintf("GetBalance_%s", u.ID)
	if temp {
		cacheKey = fmt.Sprintf("GetBalance_%s", *u.TempPublicKey)

	}
	{

		// search cache for balance
		ok, rawdata := gc.RedisCache.GetCachedResultRaw(cacheKey)

		if ok {

			log.Printf("GetBalance[%v], served from cache\n", cacheKey)
			json.Unmarshal(rawdata, &balances)
			return
		}

	}

	account, _, err := u.GetBlockchainAccountDetail(temp)
	if err != nil {
		qrCode := ""
		if !temp {

			p, e := dl.GeneratePaymentData(u.ID, "", "", "", "", gc)
			if e == nil {
				qrCode = p.QRCode
			}

		}
		balances[":"] = Balance{
			AssetIssuer: "",
			AssetCode:   "",
			Amount:      decimal.Zero,
			QRCode:      qrCode,
			ImageURL:    os.Getenv("XBN_ASSET_IMAGE_URL"),
			UsdPrice:    xbnUsdPrice,
			NativePrice: xbnNativePrice,
			InTrade: TradeLiabilties{
				SellingLiabilities: "0",
				BuyingLiabilities:  "0",
			},
		}

		if !temp && err.Error() == "error-blockchain-account-not-activated" {

			log.Printf("[GetBalance] get blockchain account detail error: %v\n", err)
			//save to cache
			gc.RedisCache.StoreResultToCacheRaw(cacheKey, balances, 60)
			return balances, nil
		}
		if !temp {
			return balances, nil
		}
		return balances, err
	}

	var wg sync.WaitGroup
	var m sync.Mutex
	// var keys []string
	for _, v := range account.Balances {
		wg.Add(1)
		go func(bal horizon.Balance) {
			// log.Printf("[BALANCE]%+v\n", v)
			defer wg.Done()
			amount, _ := decimal.NewFromString(bal.Balance)

			assetNativePrice := "0"
			assetUsdPrice := "0"
			if (temp && (amount.IsZero())) || (bal.Code == "" && temp) {
				//if nft or if it has NFT we skip
				return
			}
			if !temp && (strings.HasSuffix(strings.ToLower(bal.Code), "nft")) {
				//if nft skip in balance
				return
			}

			//liability when u hv placed a BUY
			buyingLiabilities, _ := decimal.NewFromString(bal.BuyingLiabilities)
			sellingLiabilities, _ := decimal.NewFromString(bal.SellingLiabilities)
			// availableBalance := amount.Sub(sellingLiabilities)
			availableBalance := amount.Sub(sellingLiabilities.Add(buyingLiabilities))
			// availableBalance := availableBal.Truncate(7).String()
			if bal.Issuer != "" && bal.Code != "" {
				assetNativePrice, _ := blockchain.GetNativeAskPrice(bal.Code, bal.Issuer)

				xbnUsdPriceDec := decimal.RequireFromString(xbnUsdPrice)
				nativePriceDec := decimal.RequireFromString(assetNativePrice)
				assetUsdPrice = nativePriceDec.Mul(xbnUsdPriceDec).Truncate(7).String()
			}
			if bal.Issuer == "" && bal.Code == "" {
				assetUsdPrice = xbnUsdPrice
				assetNativePrice = xbnNativePrice
			}

			qrCode := ""
			if !temp {

				p, e := dl.GeneratePaymentData(u.ID, bal.Code, bal.Issuer, "", "", gc)
				if e == nil {
					qrCode = p.QRCode
				}
			}
			imageUrl := BantuAsset{AssetCode: bal.Code, AssetIssuer: bal.Issuer}.GetAssetImageFromIssuer(gc)
			balance := Balance{AssetIssuer: bal.Issuer,
				AssetCode:   bal.Code,
				Amount:      availableBalance,
				QRCode:      qrCode,
				ImageURL:    imageUrl,
				UsdPrice:    assetUsdPrice,
				NativePrice: assetNativePrice,
				InTrade: TradeLiabilties{
					SellingLiabilities: bal.SellingLiabilities,
					BuyingLiabilities:  bal.BuyingLiabilities,
				},
			}
			// log.Printf("[BALANCE] balance: %+v\n", bal)

			m.Lock()
			balances[bal.Code+":"+bal.Issuer] = balance
			m.Unlock()

		}(v)
		wg.Wait()
	}
	//save to cache
	gc.RedisCache.StoreResultToCacheRaw(cacheKey, balances, 240)
	return balances, nil
}
func (u *User) GetUserNFTs(gc *sharedconfig.GlobalConfig) (userNFTs map[string][]NFT, err error) {

	userNFTs = make(map[string][]NFT)
	for _, wallet := range u.UserWallets {

		var wg sync.WaitGroup
		var m sync.Mutex
		//use go routine to fetch

		wg.Add(1)
		go func(vg2 UserWallet, w *sync.WaitGroup, ml *sync.Mutex) {
			defer w.Done()
			//get only the NFTs in the main wallet
			nfts, errR1 := vg2.GetNFTs(false, gc)

			if errR1 != nil {
				//log server error
				log.Printf("[GetUserNFTs] error getting NFT asset for user:[%s] wallet:[%s] error:[%+v]\n", u.Username, vg2.ID, errR1)

			}
			//Claimed Assets
			ml.Lock()
			userNFTs[vg2.ID] = nfts
			ml.Unlock()

		}(wallet, &wg, &m)
		wg.Wait()

		// log.Println("exited inner wait")

	}
	// log.Println("done...")

	return
}

// GetNFTs gets user wallet blockchain NFT balance and return it as a map of assets  [code:issuer]Balance. Native key is [:]
func (u *UserWallet) GetNFTs(temp bool, gc *sharedconfig.GlobalConfig) (nfts []NFT, err error) {
	nfts = make([]NFT, 0)
	cacheKey := fmt.Sprintf("GetNFTs_%s", u.ID)
	if temp {
		cacheKey = fmt.Sprintf("GetNFTs_%s", *u.TempPublicKey)

	}
	{

		// search cache for balance
		ok, rawdata := gc.RedisCache.GetCachedResultRaw(cacheKey)

		if ok {

			log.Printf("GetNFTs[%v], served from cache\n", cacheKey)
			json.Unmarshal(rawdata, &nfts)
			return
		}

	}

	account, _, err := u.GetBlockchainAccountDetail(temp)
	if err != nil {

		if !temp && err.Error() == "error-blockchain-account-not-activated" {

			//save to cache
			gc.RedisCache.StoreResultToCacheRaw(cacheKey, nfts, 0)
			return nil, nil
		}
		log.Printf("[GetNFTs] get blockchain account detail error: %v\n", err)

		return nil, err
	}

	var wg sync.WaitGroup
	// var keys []string
	for _, v := range account.Balances {
		wg.Add(1)
		go func(v horizon.Balance) {
			defer wg.Done()
			amount, _ := decimal.NewFromString(v.Balance)
			if amount.IsZero() {
				//if amount is zero, skip
				return
			}
			if !(strings.HasSuffix(strings.ToLower(v.Code), "nft")) {
				//if not nft skip
				return
			}
			//get name, description and image URI
			var nftName, nftDescription, nftImageURI string
			{ //TODO fetch NFT names and description and image URL

			}

			nft := NFT{AssetIssuer: v.Issuer, AssetCode: v.Code,
				NFTName: nftName, NFTDescription: nftDescription, NFTImageURI: nftImageURI}
			nfts = append(nfts, nft)

		}(v)
		wg.Wait()
	}
	//save to cache
	gc.RedisCache.StoreResultToCacheRaw(cacheKey, nfts, 0)
	return nfts, nil
}

// GetSortedUserBalance gets user blockchain balance
func (u *UserWallet) GetSortedUserBalance(temp bool, gc *sharedconfig.GlobalConfig) (balances []Balance, err error) {

	//GetBalance
	unsortedBalances, err := u.GetBalance(temp, gc)
	if err != nil {
		return
	}
	// log.Printf("unsorted balance for [%v]:[%+v]", u.ID, unsortedBalances)

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

func (u *UserWallet) GetWalletAssetBalances(gc *sharedconfig.GlobalConfig) (assetBalances AssetBalances, err error) {

	var wg sync.WaitGroup
	var m sync.Mutex
	//use go routine to fetch

	assetBalances.Unclaimed = make([]Balance, 0)
	assetBalances.Claimed = make([]Balance, 0)

	wg.Add(1)
	go func(vg1 *UserWallet, w *sync.WaitGroup, ml *sync.Mutex) {
		defer w.Done()
		unclaimedBalance, errR1 := vg1.GetSortedUserBalance(true, gc)

		if errR1 == nil {
			//Unclaimed Assets
			ml.Lock()
			assetBalances.Unclaimed = unclaimedBalance
			ml.Unlock()

		}
	}(u, &wg, &m)
	wg.Add(1)
	go func(vg2 *UserWallet, w *sync.WaitGroup, ml *sync.Mutex) {
		defer w.Done()
		claimedWalletBalance, errR1 := vg2.GetSortedUserBalance(false, gc)

		if errR1 != nil {
			//log server error
			log.Printf("[GetWalletAssetBalances] error getting claimed wallet balance for wallet:[%s] error:[%+v]\n", vg2.ID, errR1)

		}
		//Claimed Assets
		ml.Lock()
		assetBalances.Claimed = claimedWalletBalance
		ml.Unlock()

	}(u, &wg, &m)
	wg.Wait()

	return
}

func (u *User) GetUserWalletAssetBalances(gc *sharedconfig.GlobalConfig) (userWalletBalances map[string]AssetBalances, err error) {

	// userWalletBalances = make(map[string]userModels.AssetBalances)
	userWalletBalances = make(map[string]AssetBalances)
	for _, wallet := range u.UserWallets {

		var wg sync.WaitGroup
		var m sync.Mutex
		//use go routine to fetch

		var assetBalances AssetBalances
		// assetBalances.Unclaimed = make(map[string]userModels.Balance)
		assetBalances.Unclaimed = make([]Balance, 0)
		// assetBalances.Claimed = make(map[string]userModels.Balance)
		assetBalances.Claimed = make([]Balance, 0)

		wg.Add(1)
		go func(vg1 UserWallet, w *sync.WaitGroup, ml *sync.Mutex) {
			defer w.Done()
			unclaimedBalance, errR1 := vg1.GetSortedUserBalance(true, gc)

			if errR1 == nil {
				//Unclaimed Assets
				ml.Lock()
				assetBalances.Unclaimed = unclaimedBalance
				ml.Unlock()

			}
		}(wallet, &wg, &m)
		wg.Add(1)
		go func(vg2 UserWallet, w *sync.WaitGroup, ml *sync.Mutex) {
			defer w.Done()
			claimedWalletBalance, errR1 := vg2.GetSortedUserBalance(false, gc)

			if errR1 != nil {
				//log server error
				log.Printf("[GetUserWalletAssetBalances] error getting claimed wallet balance for user:[%s] wallet:[%s] error:[%+v]\n", u.Username, vg2.ID, errR1)

			}
			//Claimed Assets
			ml.Lock()
			assetBalances.Claimed = claimedWalletBalance
			ml.Unlock()

		}(wallet, &wg, &m)
		wg.Wait()
		// log.Printf("[GetUserWalletAssetBalances] finished user wallet balance:[%+v]\n", assetBalances)
		m.Lock()
		userWalletBalances[wallet.ID] = assetBalances
		m.Unlock()

		// log.Println("exited inner wait")

	}
	// log.Println("done...")

	return
}

func (id UserWalletID) GetWalletAssetBalances(gc *sharedconfig.GlobalConfig) (assetBalances AssetBalances, err error) {

	var wg sync.WaitGroup
	var m sync.Mutex
	//use go routine to fetch

	assetBalances.Unclaimed = make([]Balance, 0)
	assetBalances.Claimed = make([]Balance, 0)

	u, e := id.GetWallet(gc.DB, gc)
	if e != nil {
		return
	}
	wg.Add(1)
	go func(vg1 *UserWallet, w *sync.WaitGroup, ml *sync.Mutex) {
		defer w.Done()
		unclaimedBalance, errR1 := vg1.GetSortedUserBalance(true, gc)

		if errR1 == nil {
			//Unclaimed Assets
			ml.Lock()
			assetBalances.Unclaimed = unclaimedBalance
			ml.Unlock()

		}
	}(&u, &wg, &m)
	wg.Add(1)
	go func(vg2 *UserWallet, w *sync.WaitGroup, ml *sync.Mutex) {
		defer w.Done()
		claimedWalletBalance, errR1 := vg2.GetSortedUserBalance(false, gc)

		if errR1 != nil {
			//log server error
			log.Printf("[GetWalletAssetBalances] error getting claimed wallet balance for wallet:[%s] error:[%+v]\n", vg2.ID, errR1)

		}
		//Claimed Assets
		ml.Lock()
		assetBalances.Claimed = claimedWalletBalance
		ml.Unlock()

	}(&u, &wg, &m)
	wg.Wait()

	return
}

// GetAccountThresholds returns user signers
func (u *UserWallet) GetAccountThresholds(temp bool) (thresholds horizon.AccountThresholds) {
	account, _, err := u.GetBlockchainAccountDetail(temp)
	if err != nil {
		return thresholds
	}

	return account.Thresholds
}

// GetBlockchainAccountDetail fetches the bantu account information using public key
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
		// log.Printf("[GetBlockchainAccountDetail]: %v, error: [%v]", accountRequest.AccountID, err)
		if strings.Contains(err.Error(), "timeout") || strings.Contains(err.Error(), "handshake") || strings.Contains(err.Error(), "no such host") || strings.Contains(err.Error(), "timeout") || strings.Contains(err.Error(), "dial") {
			log.Printf("[GetBlockchainAccountDetail Network Failure]: %s\n", "Error Connecting to Expansion Service")
			return clientAccount, destinationAccountExists, &tErrors.ErrorTemporaryServerError{}
		} else {
			horizonException, ok := err.(*horizonclient.Error)

			if ok {

				if horizonException.Problem.Status == http.StatusNotFound {
					return clientAccount, false, &tErrors.ErrorBlockchainAccountNotActivated{}
				}
				log.Println("[BlockchainAccountProperties] error is known", horizonException.Problem.Status)
			}

		}
		return clientAccount, destinationAccountExists, &tErrors.ErrorTemporaryServerError{}
	}
	return clientAccount, true, nil
}

// GetBlockchainAssets fetches the blockchain asset information using public key
func (u *UserWallet) GetBlockchainAssets() (assetsPage horizon.AssetsPage, err error) {
	client := network.GetBlockchainClient()
	assetRequest := horizonclient.AssetRequest{ForAssetIssuer: u.ID, Limit: 200}
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

// OwnerOfBlockchainAssetIssued fetches the blockchain asset information using public key
func (u *UserWallet) OwnerOfBlockchainAsset(assetCode string) bool {
	client := network.GetBlockchainClient()
	assetRequest := horizonclient.AssetRequest{ForAssetIssuer: u.ID, ForAssetCode: assetCode}
	assetsPage, err := client.Assets(assetRequest)
	if err != nil {
		return false
	}
	return len(assetsPage.Embedded.Records) > 0

}

// CanIssueMoreAssets check if walet can issue more assets or has reached max limit
func (u *UserWallet) CanIssueMoreAssets() bool {
	assetPage, err := u.GetBlockchainAssets()
	if err != nil {
		return false
	}
	maxCountAssets := 200

	d, err := decimal.NewFromString(os.Getenv("MAX_ISSUED_ASSETS_PER_WALLET"))
	if err != nil {
		return len(assetPage.Embedded.Records) < maxCountAssets
	}

	return decimal.NewFromInt(int64(len(assetPage.Embedded.Records))).LessThan(d)

}

// GetBlockchainAssetsIssuedByIssuer returns blockchain assets issued by the issuer
func (u *UserWallet) GetIssuedBlockchainAssets() (issuedAssets map[string]horizon.AssetStat) {
	issuedAssets = make(map[string]horizon.AssetStat, 0)
	var err error
	assetPage, err := u.GetBlockchainAssets()
	if err != nil {
		return
	}
	//iterate through assetPage
	for _, a := range assetPage.Embedded.Records {
		issuedAssets[a.Code] = a
	}
	return
}

// GetBlockchainAssets fetches the blockchain asset information using public key
func (u Issuer) GetBlockchainAssets() (assetsPage horizon.AssetsPage, err error) {
	client := network.GetBlockchainClient()
	assetRequest := horizonclient.AssetRequest{ForAssetIssuer: string(u), Limit: 200}
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

// OwnerOfBlockchainAssetIssued fetches the blockchain asset information using public key
func (u Issuer) OwnerOfBlockchainAsset(assetCode string) bool {
	client := network.GetBlockchainClient()
	assetRequest := horizonclient.AssetRequest{ForAssetIssuer: string(u), ForAssetCode: assetCode}
	assetsPage, err := client.Assets(assetRequest)
	if err != nil {
		return false
	}
	return len(assetsPage.Embedded.Records) > 0

}

// CanIssueMoreAssets check if walet can issue more assets or has reached max limit
func (u Issuer) CanIssueMoreAssets() bool {
	assetPage, err := u.GetBlockchainAssets()
	if err != nil {
		return false
	}
	maxCountAssets := 200

	d, err := decimal.NewFromString(os.Getenv("MAX_ISSUED_ASSETS_PER_WALLET"))
	if err != nil {
		return len(assetPage.Embedded.Records) < maxCountAssets
	}

	return decimal.NewFromInt(int64(len(assetPage.Embedded.Records))).LessThan(d)

}

// GetBlockchainAssetsIssuedByIssuer returns blockchain assets issued by the issuer
func (u Issuer) GetIssuedBlockchainAssets() (issuedAssets map[string]horizon.AssetStat) {
	issuedAssets = make(map[string]horizon.AssetStat, 0)
	var err error
	assetPage, err := u.GetBlockchainAssets()
	if err != nil {
		return
	}
	//iterate through assetPage
	for _, a := range assetPage.Embedded.Records {
		issuedAssets[a.Code] = a
	}
	return
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

// GetBlockchainAccountDataKey fetches the bantu account information using public key
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
		PrimaryWallet: 1,
	}
	u.UserWallets = append(u.UserWallets, userWallet)
}

func (u *User) BuildNewSubWallet(subWalletPublicKey, walletTag, walletDescription string, walletType int, gc *sharedconfig.GlobalConfig) (userWallet UserWallet, err error) {
	walletTag = strings.ToLower(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(walletTag, "_", ""), ".", ""), " ", ""), "%", ""))
	walletDescription = strings.TrimSpace(walletDescription)
	hasMMSubwallet := false
	hasBPSubWallet := false

	if len(subWalletPublicKey) != 56 || len(walletTag) == 0 {
		log.Println("[BuildNewSubWallet] invalid parameters")
		return userWallet, &tErrors.CustomError{
			Param:      "id",
			Err:        "error-sub-wallet-parameters-invalid",
			ErrMessage: "Sub-wallet parameters are invalid. Ensure public key is 56 characters long and tag is not empty",
			Code:       http.StatusBadRequest,
		}
	}
	if len(walletDescription) == 0 {
		walletDescription = walletTag
	}
	userWallets := u.GetAllWallets(gc)
	{
		//check to ensure sub-wallet does not already exist
		for _, wallet := range userWallets {
			if wallet.ID == subWalletPublicKey {
				log.Printf("[BuildNewSubWallet] wallet [%v] already exists in your account\n", subWalletPublicKey)
				return userWallet, &tErrors.CustomError{
					Param:      "id",
					Err:        "error-sub-wallet-already-exists-in-your-account",
					ErrMessage: "Sub-wallet already exists in your account",
					Code:       http.StatusConflict,
				}
			}
			if wallet.WalletType == 2 {
				hasMMSubwallet = true
			}
			if wallet.WalletType == 3 {
				hasBPSubWallet = true
			}
			if walletType == 2 && hasMMSubwallet {
				return userWallet, &tErrors.CustomError{
					Param:      "id",
					Err:        "error-market-marking-wallet-already-exists-in-your-account",
					ErrMessage: "Market making wallet already exists in your account. You cannot have more than one Market making wallet.",
					Code:       http.StatusConflict,
				}
			}
			if walletType == 3 && hasBPSubWallet {
				return userWallet, &tErrors.CustomError{
					Param:      "id",
					Err:        "error-bulk-payment-wallet-already-exists-in-your-account",
					ErrMessage: "Bulk Payment wallet already exists in your account. You cannot have more than one Bulk Payment wallet.",
					Code:       http.StatusConflict,
				}
			}
			if wallet.Tag != nil {
				if strings.EqualFold(*wallet.Tag, walletTag) {
					log.Printf("[BuildNewSubWallet] wallet tag [%v] already exists in your account\n", walletTag)
					return userWallet, &tErrors.CustomError{
						Param:      "id",
						Err:        "error-wallet-tag-already-exists-in-your-account",
						ErrMessage: fmt.Sprintf("Sub-wallet tag [%v] already exists in your account", walletTag),
						Code:       http.StatusConflict,
					}
				}
			}

		}
	}
	{
		//check if wallet already exists in wallets
		_, errWallet := u.GetWalletByPublicKey(subWalletPublicKey, gc.DB)
		if errWallet != nil {
			if errWallet.Error() != "error-wallet-not-found" {
				log.Println("[BuildNewSubWallet] other service error ...", errWallet)

				return userWallet, errWallet
			}

		} else {
			//wallet already exists.
			log.Printf("[BuildNewSubWallet] wallet [%v] found in another account\n", subWalletPublicKey)

			return userWallet, &tErrors.CustomError{
				Param:      "id",
				Err:        "error-sub-wallet-already-exists-with-another-account",
				ErrMessage: "Sub-wallet already exists with another account",
				Code:       http.StatusConflict,
			}
		}
	}
	tempKP, pErr := network.TempAccountKeypair(subWalletPublicKey)
	var tempPK string
	if pErr != nil {
		return userWallet, &tErrors.CustomError{
			Param:      "id",
			Err:        "error-sub-wallet-public-key-invalid",
			ErrMessage: "Sub-wallet public key is invalid",
			Code:       http.StatusBadRequest,
		}
	}
	if tempKP != nil {
		tempPK = tempKP.Address()
	}

	alias := fmt.Sprintf("%s_%s", u.Username, walletTag)
	userSubWallet := UserWallet{
		ID:            subWalletPublicKey,
		TempPublicKey: &tempPK,
		Tag:           &walletTag,
		Description:   &walletDescription,
		Alias:         alias,
		Signer:        u.PrimarySigner,
		UserID:        u.ID,
		WalletType:    walletType,
	}
	// u.UserWallets = append(u.UserWallets, userSubWallet)
	return userSubWallet, nil
}

func (id UserWalletID) String() string {
	return string(id)
}

func (publicKey UserSigner) String() string {
	return string(publicKey)
}
func (u Username) String() string {
	return string(u)
}

func (publicKey UserSigner) GetOwner(db *gorm.DB, gc *sharedconfig.GlobalConfig) (signerOwner User, err error) {
	cacheKeyInfo := fmt.Sprintf("userObj %v", string(publicKey))

	{

		// search cache for balance
		ok, rawdata := gc.RedisCache.GetCachedResultRaw(cacheKeyInfo)

		if ok {

			log.Printf("GetOwner[%v], served from cache\n", cacheKeyInfo)
			json.Unmarshal(rawdata, &signerOwner)
			return
		}

	}

	e := db.Preload("UserWallets.Permissions").Preload(clause.Associations).Where("primary_signer = ?", string(publicKey)).First(&signerOwner).Error
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

	// log.Printf("user for %v is %v\n", userInfo, user)
	gc.RedisCache.StoreResultToCacheRaw(cacheKeyInfo, signerOwner, 2000)
	cacheKeyUsername := fmt.Sprintf("userObj %v", signerOwner.Username)
	gc.RedisCache.StoreResultToCacheRaw(cacheKeyUsername, signerOwner, 2000)
	cacheKeyEmail := fmt.Sprintf("userObj %v", signerOwner.Email)
	gc.RedisCache.StoreResultToCacheRaw(cacheKeyEmail, signerOwner, 2000)
	cacheKeySigner := fmt.Sprintf("userObj %v", signerOwner.PrimarySigner)
	gc.RedisCache.StoreResultToCacheRaw(cacheKeySigner, signerOwner, 2000)
	cacheKeyUserID := fmt.Sprintf("userObj %v", signerOwner.ID)
	gc.RedisCache.StoreResultToCacheRaw(cacheKeyUserID, signerOwner, 2000)
	return
}

func (id UserWalletID) GetPermissionList(db *gorm.DB) (accessList []WalletPermission) {
	accessList = make([]WalletPermission, 0)
	db.Preload(clause.Associations).Where("wallet_public_key = ?", string(id)).Find(&accessList)

	return
}

func (u *UserWallet) GetPermissionList(db *gorm.DB) (accessList []WalletPermission) {

	if u.Permissions != nil {

		if len(u.Permissions) > 0 {
			return u.Permissions
		}
	}
	accessList = make([]WalletPermission, 0)
	db.Preload(clause.Associations).Where("wallet_public_key = ?", u.ID).Find(&accessList)

	return
}

func (a WalletAlias) GetAccessList(db *gorm.DB, gc *sharedconfig.GlobalConfig) (accessList []WalletPermission) {

	accessList = make([]WalletPermission, 0)
	wallet, e := a.GetWallet(db, gc)
	if e != nil {
		return accessList
	}
	if wallet.SharedAccessEnabled == 0 {
		return accessList
	}
	if len(wallet.Permissions) > 0 {
		accessList = wallet.Permissions
	}

	return
}

func (u *UserWallet) PublicKeyHasViewOnlyAccess(gc *sharedconfig.GlobalConfig) (viewOnly bool) {
	if u == nil {
		return true
	}
	if u.NumberOfApprovalsNeeded > 0 && u.SharedAccessEnabled == 1 {
		return false
	}

	return true
}
func (u *UserWallet) HasViewOnlyAccess(gc *sharedconfig.GlobalConfig) (viewOnly bool) {
	if u == nil {
		return true
	}
	if u.NumberOfApprovalsNeeded > 0 && u.SharedAccessEnabled == 1 {
		return false
	}

	return true
}

func (u *UserWallet) SignerHasAccess(signer *User, gc *sharedconfig.GlobalConfig) (hasAccess bool) {
	hasAccess = false
	if u == nil {
		return false
	}
	if signer == nil {
		return false
	}
	if u.SharedAccessEnabled == 0 {
		return false
	}
	permList := u.GetPermissionList(gc.DB)
	if len(permList) == 0 {
		return false
	}

	for _, access := range permList {
		if access.TargetUsername == signer.Username {
			return true
		}
	}

	return
}

func (u *UserWallet) SignerHasInitatorAccess(signer *User, gc *sharedconfig.GlobalConfig) (hasAccess bool) {
	hasAccess = false
	if u == nil {
		return false
	}
	if signer == nil {
		return false
	}
	if u.SharedAccessEnabled == 0 {
		return false
	}
	permList := u.GetPermissionList(gc.DB)
	if len(permList) == 0 {
		return false
	}

	for _, access := range permList {
		if (access.TargetUsername == signer.Username) && (access.Permission == "INITIATOR") {
			return true
		}
	}

	return
}

func (u *UserWallet) SignerHasApproverAccess(signer *User, gc *sharedconfig.GlobalConfig) (hasAccess bool) {
	hasAccess = false
	if u == nil {
		return false
	}
	if signer == nil {
		return false
	}
	if u.SharedAccessEnabled == 0 {
		return false
	}
	permList := u.GetPermissionList(gc.DB)
	if len(permList) == 0 {
		return false
	}

	for _, access := range permList {
		if (access.TargetUsername == signer.Username) && (access.Permission == "APPROVER") {
			return true
		}
	}

	return
}

func (u *UserWallet) SignerKeyHasAccess(signerKey string, gc *sharedconfig.GlobalConfig) (hasAccess bool) {
	hasAccess = false
	if u == nil {
		return false
	}
	if signerKey == "" {
		return false
	}
	if u.SharedAccessEnabled == 0 {
		return false
	}

	signer, err := UserSigner(signerKey).GetOwner(gc.DB, gc)
	if err != nil {
		return false
	}
	permList := u.GetPermissionList(gc.DB)
	if len(permList) == 0 {
		return false
	}

	for _, access := range permList {
		if access.TargetUsername == signer.Username {
			return true
		}
	}

	return
}

func (u *UserWallet) SignerKeyHasInitatorAccess(signerKey string, gc *sharedconfig.GlobalConfig) (hasAccess bool) {
	hasAccess = false
	if u == nil {
		return false
	}
	if signerKey == "" {
		return false
	}
	if u.SharedAccessEnabled == 0 {
		return false
	}
	signer, err := UserSigner(signerKey).GetOwner(gc.DB, gc)
	if err != nil {
		return false
	}
	permList := u.GetPermissionList(gc.DB)
	if len(permList) == 0 {
		return false
	}

	for _, access := range permList {
		if (access.TargetUsername == signer.Username) && (access.Permission == "INITIATOR") {
			return true
		}
	}

	return
}

func (u *UserWallet) SignerKeyHasApproverAccess(signerKey string, gc *sharedconfig.GlobalConfig) (hasAccess bool) {
	hasAccess = false
	if u == nil {
		return false
	}
	if signerKey == "" {
		return false
	}
	if u.SharedAccessEnabled == 0 {
		return false
	}
	signer, err := UserSigner(signerKey).GetOwner(gc.DB, gc)
	if err != nil {
		return false
	}
	permList := u.GetPermissionList(gc.DB)
	if len(permList) == 0 {
		return false
	}

	for _, access := range permList {
		if (access.TargetUsername == signer.Username) && (access.Permission == "APPROVER") {
			return true
		}
	}

	return
}

func (u *UserWallet) GetWalletOwner(db *gorm.DB, gc *sharedconfig.GlobalConfig) (walletOwner User, err error) {
	cacheKeyInfo := fmt.Sprintf("userObj %v", u.ID)

	{

		// search cache for balance
		ok, rawdata := gc.RedisCache.GetCachedResultRaw(cacheKeyInfo)

		if ok {

			log.Printf("GetWalletOwner[%v], served from cache\n", cacheKeyInfo)
			json.Unmarshal(rawdata, &walletOwner)
			return
		}

	}

	e := db.Preload("UserWallets.Permissions").Preload(clause.Associations).Where("id = ?", u.UserID).First(&walletOwner).Error
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
	gc.RedisCache.StoreResultToCacheRaw(cacheKeyInfo, walletOwner, 2000)
	cacheKeyUsername := fmt.Sprintf("userObj %v", walletOwner.Username)
	gc.RedisCache.StoreResultToCacheRaw(cacheKeyUsername, walletOwner, 2000)
	cacheKeyEmail := fmt.Sprintf("userObj %v", walletOwner.Email)
	gc.RedisCache.StoreResultToCacheRaw(cacheKeyEmail, walletOwner, 2000)
	cacheKeySigner := fmt.Sprintf("userObj %v", walletOwner.PrimarySigner)
	gc.RedisCache.StoreResultToCacheRaw(cacheKeySigner, walletOwner, 2000)
	cacheKeyUserID := fmt.Sprintf("userObj %v", walletOwner.ID)
	gc.RedisCache.StoreResultToCacheRaw(cacheKeyUserID, walletOwner, 2000)
	return
}

func (u *User) GetMartketMakingWallet(db *gorm.DB) (wallet UserWallet, err error) {
	e := db.Preload(clause.Associations).Where("user_id = ? AND wallet_type = 2", u.ID).First(&wallet).Error
	if e != nil {
		if errors.Is(e, gorm.ErrRecordNotFound) {
			//no wallet was found
			err = &tErrors.CustomError{
				Param:      "id",
				Err:        "error-market-making-wallet-not-found",
				ErrMessage: "Market Making Wallet not found",
				Code:       404,
			}
			return
		}
		err = &tErrors.ErrorTemporaryServerError{}
	}
	return
}

func (u *User) GetBulkPaymentWallet(db *gorm.DB) (wallet UserWallet, err error) {
	e := db.Preload(clause.Associations).Where("user_id = ? AND wallet_type = 3", u.ID).First(&wallet).Error
	if e != nil {
		if errors.Is(e, gorm.ErrRecordNotFound) {
			//no wallet was found
			err = &tErrors.CustomError{
				Param:      "id",
				Err:        "error-bulk-payment-wallet-not-found",
				ErrMessage: "Bulk Payment Wallet not found",
				Code:       404,
			}
			return
		}
		err = &tErrors.ErrorTemporaryServerError{}
	}
	return
}

func (id UserWalletID) GetWallet(db *gorm.DB, gc *sharedconfig.GlobalConfig) (wallet UserWallet, err error) {
	cacheKeyInfo := fmt.Sprintf("walletObj_%v", string(id))

	{

		// search cache for balance
		ok, rawdata := gc.RedisCache.GetCachedResultRaw(cacheKeyInfo)

		if ok {

			log.Printf("GetWallet[%v], served from cache\n", cacheKeyInfo)
			json.Unmarshal(rawdata, &wallet)
			return
		}

	}
	// cacheKeyAlias := fmt.Sprintf("walletObj_%v", wallet.Alias)
	// gc.RedisCache.StoreResultToCacheRaw(cacheKeyInfo, wallet, 0)
	// gc.RedisCache.StoreResultToCacheRaw(cacheKeyAlias, wallet, 0)

	e := db.Preload(clause.Associations).Where("id = ?", string(id)).First(&wallet).Error
	if e != nil {
		if errors.Is(e, gorm.ErrRecordNotFound) {
			//no wallet was found
			err = &tErrors.ErrorInvalidWallet{
				PublicKey: string(id),
			}
			return
		}
		err = &tErrors.ErrorTemporaryServerError{}
	}

	cacheKeyAlias := fmt.Sprintf("walletObj_%v", wallet.Alias)

	gc.RedisCache.StoreResultToCacheRaw(cacheKeyInfo, wallet, 0)
	gc.RedisCache.StoreResultToCacheRaw(cacheKeyAlias, wallet, 0)
	return
}

func (a WalletAlias) GetWallet(db *gorm.DB, gc *sharedconfig.GlobalConfig) (wallet UserWallet, err error) {
	cacheKeyInfo := fmt.Sprintf("walletObj_%v", string(a))

	{

		// search cache for balance
		ok, rawdata := gc.RedisCache.GetCachedResultRaw(cacheKeyInfo)

		if ok {

			log.Printf("GetWallet[%v], served from cache\n", cacheKeyInfo)
			json.Unmarshal(rawdata, &wallet)
			return
		}

	}
	// cacheKeyAlias := fmt.Sprintf("walletObj_%v", wallet.Alias)
	// gc.RedisCache.StoreResultToCacheRaw(cacheKeyInfo, wallet, 0)
	// gc.RedisCache.StoreResultToCacheRaw(cacheKeyAlias, wallet, 0)

	e := db.Preload(clause.Associations).Where("alias = ?", string(a)).First(&wallet).Error
	if e != nil {
		if errors.Is(e, gorm.ErrRecordNotFound) {
			//no wallet was found
			err = &tErrors.ErrorInvalidWallet{
				PublicKey: string(a),
			}
			return
		}
		err = &tErrors.ErrorTemporaryServerError{}
	}
	cacheKeyID := fmt.Sprintf("walletObj_%v", wallet.ID)
	gc.RedisCache.StoreResultToCacheRaw(cacheKeyInfo, wallet, 0)
	gc.RedisCache.StoreResultToCacheRaw(cacheKeyID, wallet, 0)
	return
}

func (a ApprovalID) GetSubmittedTransaction(db *gorm.DB) (pendingAuth PendingAuth, err error) {
	e := db.Preload(clause.Associations).Where("id = ?", string(a)).First(&pendingAuth).Error
	if e != nil {
		if errors.Is(e, gorm.ErrRecordNotFound) {
			//no approval request was found
			err = &tErrors.ErrorInvalidRequest{
				ID: string(a),
			}
			return
		}
		err = &tErrors.ErrorTemporaryServerError{}
	}
	return
}

func (id UserWalletID) PublicKeyHasViewOnlyAccess(gc *sharedconfig.GlobalConfig) (viewOnly bool) {
	if id == "" {
		return true
	}
	wallet, err := id.GetWallet(gc.DB, gc)
	if err != nil {
		return true
	}

	if wallet.NumberOfApprovalsNeeded > 0 && wallet.SharedAccessEnabled == 1 {
		return false
	}

	return true
}

func (u *User) GetCuratedSwapList(db *gorm.DB) (list []assets.CuratedSwapAsset) {

	list = make([]assets.CuratedSwapAsset, 0)
	// var swapAsset assets.CuratedSwapAsset
	e := db.Find(&list).Error
	if e != nil {
		return make([]assets.CuratedSwapAsset, 0)
	}
	return list
}

func (u *User) GetWalletByPublicKey(publicKey string, db *gorm.DB) (wallet UserWallet, err error) {
	e := db.Preload(clause.Associations).Where("id = ?", publicKey).First(&wallet).Error
	if e != nil {
		if errors.Is(e, gorm.ErrRecordNotFound) {
			//no wallet was found
			err = &tErrors.CustomError{
				Param:      "publicKey",
				Err:        "error-wallet-not-found",
				ErrMessage: "Wallet not found",
				Code:       404,
			}
			return
		}
		log.Printf("[GetWalletByPublicKey] error: %s\n", e)
		err = &tErrors.ErrorTemporaryServerError{}
	}
	return wallet, nil
}

func (id UserWalletID) GetWalletOwner(db *gorm.DB, gc *sharedconfig.GlobalConfig) (walletOwner User, err error) {
	cacheKeyInfo := fmt.Sprintf("userObj %v", string(id))

	{

		// search cache for balance
		ok, rawdata := gc.RedisCache.GetCachedResultRaw(cacheKeyInfo)

		if ok {

			log.Printf("GetWalletOwner[%v], served from cache\n", cacheKeyInfo)
			json.Unmarshal(rawdata, &walletOwner)
			return
		}

	}
	e := db.Preload("UserWallets.Permissions").Preload(clause.Associations).Where("id = (SELECT user_id FROM user_wallets WHERE id = ?)", string(id)).First(&walletOwner).Error
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
	gc.RedisCache.StoreResultToCacheRaw(cacheKeyInfo, walletOwner, 2000)
	cacheKeyUsername := fmt.Sprintf("userObj %v", walletOwner.Username)
	gc.RedisCache.StoreResultToCacheRaw(cacheKeyUsername, walletOwner, 2000)
	cacheKeyEmail := fmt.Sprintf("userObj %v", walletOwner.Email)
	gc.RedisCache.StoreResultToCacheRaw(cacheKeyEmail, walletOwner, 2000)
	cacheKeySigner := fmt.Sprintf("userObj %v", walletOwner.PrimarySigner)
	gc.RedisCache.StoreResultToCacheRaw(cacheKeySigner, walletOwner, 2000)
	cacheKeyUserID := fmt.Sprintf("userObj %v", walletOwner.ID)
	gc.RedisCache.StoreResultToCacheRaw(cacheKeyUserID, walletOwner, 2000)
	return
}

func (id UserWalletID) GetUserPermissionOnWallet(username, permission string, db *gorm.DB) (walletPermission WalletPermission, err error) {
	e := db.Where("target_username = ? AND wallet_public_key = ? AND permission = ?", username, string(id), strings.ToUpper(permission)).First(&walletPermission).Error
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

func (u Username) GetFullUser(db *gorm.DB, gc *sharedconfig.GlobalConfig) (owner User, err error) {
	cacheKeyInfo := fmt.Sprintf("userObj %v", string(u))

	{

		// search cache for balance
		ok, rawdata := gc.RedisCache.GetCachedResultRaw(cacheKeyInfo)

		if ok {

			log.Printf("GetFullUser[%v], served from cache\n", cacheKeyInfo)
			json.Unmarshal(rawdata, &owner)
			return
		}

	}

	e := db.Preload("UserWallets.Permissions").Preload(clause.Associations).Where("username = ?", string(u)).First(&owner).Error
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
	gc.RedisCache.StoreResultToCacheRaw(cacheKeyInfo, owner, 2000)
	cacheKeyUsername := fmt.Sprintf("userObj %v", owner.Username)
	gc.RedisCache.StoreResultToCacheRaw(cacheKeyUsername, owner, 2000)
	cacheKeyEmail := fmt.Sprintf("userObj %v", owner.Email)
	gc.RedisCache.StoreResultToCacheRaw(cacheKeyEmail, owner, 2000)
	cacheKeySigner := fmt.Sprintf("userObj %v", owner.PrimarySigner)
	gc.RedisCache.StoreResultToCacheRaw(cacheKeySigner, owner, 2000)
	cacheKeyUserID := fmt.Sprintf("userObj %v", owner.ID)
	gc.RedisCache.StoreResultToCacheRaw(cacheKeyUserID, owner, 2000)
	return
}

func (u Username) GetSimpleUser(db *gorm.DB, gc *sharedconfig.GlobalConfig) (owner User, err error) {

	cacheKeyInfo := fmt.Sprintf("userObj %v", string(u))

	{

		// search cache for balance
		ok, rawdata := gc.RedisCache.GetCachedResultRaw(cacheKeyInfo)

		if ok {

			log.Printf("GetSimpleUser[%v], served from cache\n", cacheKeyInfo)
			json.Unmarshal(rawdata, &owner)
			return
		}

	}

	// e := db.Where("username = ?", string(u)).First(&owner).Error
	// if e != nil {
	// 	if errors.Is(e, gorm.ErrRecordNotFound) {
	// 		//no wallet was found
	// 		err = &tErrors.CustomError{
	// 			Param:      "id",
	// 			Err:        "error-account-not-found",
	// 			ErrMessage: "Account not found",
	// 			Code:       404,
	// 		}
	// 		return
	// 	}
	// 	err = &tErrors.ErrorTemporaryServerError{}
	// }
	return u.GetFullUser(db, gc)
}

func (u Username) GetUserPermissionOnWallet(walletPublicKey string, db *gorm.DB) (walletPermission WalletPermission, err error) {
	e := db.Where("target_username = ? AND wallet_public_key = ?", string(u), walletPublicKey).First(&walletPermission).Error
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

func (u *User) HasAccessToPublicKey(publicKey string, gc *sharedconfig.GlobalConfig) (hasAccess bool) {
	// walletPermissions := u.Fetch3rdPartyWalletPermissions(gc)
	// if len(walletPermissions) == 0 {
	// 	return false
	// }
	for _, walletAccess := range u.WalletsSharedWithUser {
		if walletAccess.WalletPublicKey == publicKey {
			return true
		}
	}

	return false
}

func (u *User) GetAllWallets(gc *sharedconfig.GlobalConfig) (wallets []UserWallet) {
	wallets = make([]UserWallet, 0)
	if len(u.UserWallets) == 0 {
		gc.DB.Preload(clause.Associations).Where("user_id = ?", u.ID).Find(&wallets)
		return
	}

	return u.UserWallets
}

func (wp *WalletPermission) ToWalletPermissionInfo(user *User, w *UserWallet, gc *sharedconfig.GlobalConfig) (wpInfo WalletPermissionInfo) {
	var name string
	var pushNotificationToken *string
	if user == nil {
		var u User
		e := gc.DB.Where("username = ?", wp.TargetUsername).First(&u).Error
		if e == nil {
			pushNotificationToken = u.PushNotificationToken
			name = u.FirstName
			if u.LastName != nil {
				name = name + " " + *u.LastName
			}
		}
	} else {
		pushNotificationToken = user.PushNotificationToken
		name = user.FirstName
		if user.LastName != nil {
			name = name + " " + *user.LastName
		}
	}

	wpInfo = WalletPermissionInfo{
		ID:                    wp.ID,
		WalletPublicKey:       wp.WalletPublicKey,
		WalletAlias:           w.Alias,
		TargetUsername:        wp.TargetUsername,
		Name:                  name,
		PushNotificationToken: pushNotificationToken,
		Permission:            wp.Permission,
	}
	return
}

// FetchWalletsPermissionsSharedWithUser fetches all 3rd party wallets that the user is assigned to manage
func (u *User) FetchWalletsPermissionsSharedWithUser(gc *sharedconfig.GlobalConfig) (thirdPartyWallets []WalletsSharedWithUser) {
	// var walletPermissions []WalletPermission
	thirdPartyWallets = make([]WalletsSharedWithUser, 0)
	if u.WalletsSharedWithUser == nil {
		return
	}
	if len(u.WalletsSharedWithUser) == 0 {
		return
	}
	// cacheKey := fmt.Sprintf("FetchWalletsPermissionsSharedWithUser_%s", u.ID)

	// {

	// 	// search cache for balance
	// 	ok, rawdata := gc.RedisCache.GetCachedResultRaw(cacheKey)

	// 	if ok {

	// 		log.Printf("FetchWalletsPermissionsSharedWithUser[%v], served from cache\n", cacheKey)
	// 		json.Unmarshal(rawdata, &thirdPartyWallets)
	// 		return
	// 	}

	// }

	for _, assignedPermission := range u.WalletsSharedWithUser {
		//Get the permission assignment
		thirdPartyWallet := WalletsSharedWithUser{
			Permission: assignedPermission.Permission,
		}

		//use it to fetch wallet details
		wallet, err := UserWalletID(assignedPermission.WalletPublicKey).GetWallet(gc.DB, gc)
		if err != nil {
			return
		}

		thirdPartyWallet.WalletPublicKey = wallet.ID
		thirdPartyWallet.WalletAlias = wallet.Alias
		if wallet.Description != nil {
			thirdPartyWallet.WalletDescription = *wallet.Description
		}
		thirdPartyWallet.AssetBalances, _ = wallet.GetWalletAssetBalances(gc)
		//use it to fetch wallet owner details
		owner, err := UserWalletID(assignedPermission.WalletPublicKey).GetWalletOwner(gc.DB, gc)
		if err != nil {
			return
		}

		thirdPartyWallet.Owner = owner.Username

		if assignedPermission.Permission == "INITIATOR" {
			//set walletSettings SINCE INITIATORS CAN MODIFY WALLET
			walletSettings := WalletSettings{
				NumberOfApprovalsNeeded: wallet.NumberOfApprovalsNeeded,
				WalletType:              wallet.WalletType,
			}
			for _, p := range wallet.Permissions {
				walletSettings.Permissions = append(walletSettings.Permissions, p.ToJSON(gc))
			}
			thirdPartyWallet.WalletSettings = &walletSettings
		}

		thirdPartyWallets = append(thirdPartyWallets, thirdPartyWallet)

	}
	//save to cache
	// gc.RedisCache.StoreResultToCacheRaw(cacheKey, thirdPartyWallets, 0)

	return
}
func (w *UserWallet) WalletCountApproverAccess(gc *sharedconfig.GlobalConfig) (accessCount uint) {
	if w.SharedAccessEnabled == 0 {
		return 0
	}
	if w.Permissions == nil {
		return 0
	}
	for _, access := range w.Permissions {
		if access.Permission == "APPROVER" {
			accessCount++
		}
	}

	return
}

func (w *UserWallet) WalletCountInitiatorAccess(gc *sharedconfig.GlobalConfig) (accessCount uint) {
	if w.SharedAccessEnabled == 0 {
		return 0
	}
	if w.Permissions == nil {
		return 0
	}
	for _, access := range w.Permissions {
		if access.Permission == "INITIATOR" {
			accessCount++
		}
	}

	return
}

func (w *UserWallet) HasInitiatorPermissionToPublicKey(ownerSignerPublicKey string, gc *sharedconfig.GlobalConfig) (hasAccess bool) {
	user, err := UserSigner(ownerSignerPublicKey).GetOwner(gc.DB, gc)

	if err != nil {
		return false
	}
	walletPermissions := user.FetchWalletsPermissionsSharedWithUser(gc)
	if len(walletPermissions) == 0 {
		return false
	}
	for _, walletAccess := range walletPermissions {
		if walletAccess.WalletPublicKey == w.ID && walletAccess.Permission == "INITIATOR" {
			return true
		}
	}

	return false
}

func (w *UserWallet) SignerHasInitiatorPermissionToPublicKey(signerOwner User, gc *sharedconfig.GlobalConfig) (hasAccess bool) {

	walletPermissions := signerOwner.FetchWalletsPermissionsSharedWithUser(gc)
	if len(walletPermissions) == 0 {
		return false
	}
	for _, walletAccess := range walletPermissions {
		if walletAccess.WalletPublicKey == w.ID && walletAccess.Permission == "INITIATOR" {
			return true
		}
	}

	return false
}

func (u *User) GetDefaultAssets(gc *sharedconfig.GlobalConfig) (defaultAssets []DefaultAsset) {
	cacheKey := "GetDefaultAssets_"

	{

		// search cache for balance
		ok, response := gc.RedisCache.GetCachedResult(cacheKey)

		if ok {
			log.Printf("GetDefaultAssets [%v], served from cache\n", cacheKey)
			da := response.([]interface{})
			for _, v := range da {
				vals := v.(map[string]interface{})
				// log.Printf("VALS: [%+v]\n", vals)
				defaultAssets = append(defaultAssets, DefaultAsset{
					AssetCode:   vals["assetCode"].(string),
					AssetIssuer: vals["assetIssuer"].(string),
					ImageURL:    vals["imageUrl"].(string),
				})
			}
			return
		}

	}
	e := gc.DB.Find(&defaultAssets).Error
	if e != nil {
		log.Printf("[User.GetDefaultAssets] Error pulling default Assets, Error: %v", e)
	}
	//save to cache
	gc.RedisCache.StoreResultToCache(cacheKey, defaultAssets, 200)
	return defaultAssets

}

func (u *User) SendPushMessage(title, body, imageURI string, dataPayload map[string]string, gc *sharedconfig.GlobalConfig) {
	//Send push notification to user
	// log.Println(title, body)

	if u.PushNotificationToken == nil {
		return
	}

	pns.SendFirebaseMessage(*u.PushNotificationToken, title, body, imageURI, dataPayload, gc.PushNotificationClient, gc.PNSContext)

}

func (w *UserWallet) InvalidateUserCache(gc *sharedconfig.GlobalConfig) {
	userAccount, err := w.GetWalletOwner(gc.DB, gc)
	if err != nil {
		return
	}
	cacheKey1 := fmt.Sprintf("GetBalance_%s", userAccount.PublicKey)
	cacheKeyUsername := fmt.Sprintf("userObj %v", userAccount.Username)
	cacheKeyEmail := fmt.Sprintf("userObj %v", userAccount.Email)
	cacheKeySigner := fmt.Sprintf("userObj %v", userAccount.PrimarySigner)
	cacheKeyUserID := fmt.Sprintf("userObj %v", userAccount.ID)
	gc.RedisCache.DeleteFromCache(cacheKeyUsername, cacheKeyEmail, cacheKeySigner, cacheKeyUserID)

	gc.RedisCache.DeleteFromCache(cacheKey1)
	userAccount.InvalidateUserWalletCache(gc)
}
func (u Username) InvalidateUserCache(id string, gc *sharedconfig.GlobalConfig) {
	userAccount, err := u.GetFullUser(gc.DB, gc)
	if err != nil {
		return
	}
	cacheKey1 := fmt.Sprintf("GetBalance_%s", userAccount.PublicKey)
	cacheKeyUsername := fmt.Sprintf("userObj %v", userAccount.Username)
	cacheKeyEmail := fmt.Sprintf("userObj %v", userAccount.Email)
	cacheKeySigner := fmt.Sprintf("userObj %v", userAccount.PrimarySigner)
	cacheKeyUserID := fmt.Sprintf("userObj %v", userAccount.ID)
	gc.RedisCache.DeleteFromCache(cacheKeyUsername, cacheKeyEmail, cacheKeySigner, cacheKeyUserID)

	gc.RedisCache.DeleteFromCache(cacheKey1)
	userAccount.InvalidateUserWalletCache(gc)
}
func (u *User) InvalidateUserCache(gc *sharedconfig.GlobalConfig) {
	if u == nil {
		return
	}
	cacheKey1 := fmt.Sprintf("GetBalance_%s", u.PublicKey)
	cacheKeyUsername := fmt.Sprintf("userObj %v", u.Username)
	cacheKeyEmail := fmt.Sprintf("userObj %v", u.Email)
	cacheKeySigner := fmt.Sprintf("userObj %v", u.PrimarySigner)
	cacheKeyUserID := fmt.Sprintf("userObj %v", u.ID)
	cacheKeyPShared := fmt.Sprintf("FetchWalletsPermissionsSharedWithUser_%s", u.ID)
	gc.RedisCache.DeleteFromCache(cacheKeyPShared, cacheKeyUsername, cacheKeyEmail, cacheKeySigner, cacheKeyUserID)

	gc.RedisCache.DeleteFromCache(cacheKey1)
	u.InvalidateUserWalletCache(gc)
}

func (u *User) InvalidateUserWalletCache(gc *sharedconfig.GlobalConfig) {

	if u == nil {
		return
	}
	if u.UserWallets == nil {
		return
	}
	if len(u.UserWallets) == 0 {
		return
	}
	for _, w := range u.UserWallets {
		cacheKey1 := fmt.Sprintf("GetBalance_%s", w.ID)
		cacheKey2 := fmt.Sprintf("GetBalance_%s", *w.TempPublicKey)

		cacheKey3 := fmt.Sprintf("userObj %v", w.Alias)
		cacheKey4 := fmt.Sprintf("userObj %v", w.ID)
		cacheKeySigner := fmt.Sprintf("userObj %v", w.Signer)
		cacheKeyUserID := fmt.Sprintf("userObj %v", w.UserID)
		cacheKeyWalletAlias := fmt.Sprintf("walletObj_%v", w.Alias)
		cacheKeyWalletID := fmt.Sprintf("walletObj_%v", w.ID)
		cacheKeyPShared := fmt.Sprintf("FetchWalletsPermissionsSharedWithUser_%s", w.UserID)
		gc.RedisCache.DeleteFromCache(cacheKeyPShared, cacheKeyWalletAlias, cacheKeyWalletID, cacheKey1, cacheKey2, cacheKey3, cacheKey4, cacheKeySigner, cacheKeyUserID)

	}

}

func (w UserWallet) GetMarketOfferByID(offerID string, db *gorm.DB, gc *sharedconfig.GlobalConfig) (marketOfer MarketOffer, err error) {
	e := db.Where("id = ?", offerID).Where("source_wallet_public_key = ?", w.ID).First(&marketOfer).Error
	if e != nil {
		if errors.Is(e, gorm.ErrRecordNotFound) {
			//no wallet was found
			err = &tErrors.CustomError{
				Param:      "id",
				Err:        "error-offer-not-found",
				ErrMessage: "Offer not found",
				Code:       404,
			}
			return
		}
		err = &tErrors.ErrorTemporaryServerError{}
	}
	return
}
func (s MarketOfferID) GetMarketOffer(db *gorm.DB, gc *sharedconfig.GlobalConfig) (marketOfer MarketOffer, err error) {
	e := db.Where("id = ?", string(s)).First(&marketOfer).Error
	if e != nil {
		if errors.Is(e, gorm.ErrRecordNotFound) {
			//no wallet was found
			err = &tErrors.CustomError{
				Param:      "id",
				Err:        "error-offer-not-found",
				ErrMessage: "Offer not found",
				Code:       404,
			}
			return
		}
		err = &tErrors.ErrorTemporaryServerError{}
	}
	return
}

func (mo *MarketOffer) GetBlockchainOfferDetail(gc *sharedconfig.GlobalConfig) (offer horizon.Offer, err error) {
	if mo.BlockchainOfferID == nil {
		err = &tErrors.ErrorMissingParameter{Parameter: "blockchainOfferId"}
		return
	}

	offer, e := gc.BantuExpansionClient.OfferDetails(*mo.BlockchainOfferID)

	if e != nil {
		log.Println("[GetBlockchainOffer]: ", e)
		if strings.Contains(e.Error(), "timeout") || strings.Contains(e.Error(), "handshake") || strings.Contains(e.Error(), "no such host") || strings.Contains(e.Error(), "timeout") || strings.Contains(e.Error(), "dial") {
			log.Printf("[GetBlockchainOffer Network Failure]: %s\n", "Error Connecting to Blockchain API Service")
			return offer, &tErrors.ErrorTemporaryServerError{}
		}

		return offer, &tErrors.ErrorTemporaryServerError{}

	}
	return offer, nil
}

func (mo *MarketOffer) CancelBlockchainOffer(gc *sharedconfig.GlobalConfig) (offer horizon.Offer, err error) {
	if mo.BlockchainOfferID == nil {
		err = &tErrors.ErrorMissingParameter{Parameter: "blockchainOfferId"}
		return
	}

	offer, e := gc.BantuExpansionClient.OfferDetails(*mo.BlockchainOfferID)

	if e != nil {
		log.Println("[GetBlockchainOffer]: ", e)
		if strings.Contains(e.Error(), "timeout") || strings.Contains(e.Error(), "handshake") || strings.Contains(e.Error(), "no such host") || strings.Contains(e.Error(), "timeout") || strings.Contains(e.Error(), "dial") {
			log.Printf("[GetBlockchainOffer Network Failure]: %s\n", "Error Connecting to Blockchain API Service")
			return offer, &tErrors.ErrorTemporaryServerError{}
		}

		return offer, &tErrors.ErrorTemporaryServerError{}

	}
	return offer, nil
}
