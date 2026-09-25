package users

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
	"trovo-wallet-api/internal/basetxn"
	blockchain "trovo-wallet-api/internal/components/assets/blockchain"
	assetsDB "trovo-wallet-api/internal/components/assets/db"
	assets "trovo-wallet-api/internal/components/assets/models"
	dl "trovo-wallet-api/internal/dynamiclinks"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/network"
	pns "trovo-wallet-api/internal/pns"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/mailgun/mailgun-go/v4"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

/////User Model convenience Methods

// AccountBalance is the Base equivalent of one horizon.Balance entry.
// Code == "" means the native asset, matching the original's convention.
type AccountBalance struct {
	Code    string
	Issuer  string
	Balance string
	// SellingLiabilities/BuyingLiabilities tracked a Stellar account's
	// open DEX orders against this balance; always "0" on Base, which
	// has no native order book (see internal/sharedconfig/
	// order_book_summary.go's doc) to hold such an order.
	SellingLiabilities string
	BuyingLiabilities  string
}

// AccountDetail is the Base equivalent of Stellar's *horizon.Account.
// Signers/Thresholds are vestigial: this app's Base accounts are plain
// EOAs with exactly one signer (themselves) rather than Stellar's
// on-chain weighted multisig, so Signers always reports just the
// account's own address at weight 1 - real multi-party approval
// (subwallets/shared access) is tracked in this app's own DB instead
// (see internal/components/users/services), not via account-level
// signer weights.
type AccountDetail struct {
	ID         string
	Signers    []Signer
	Thresholds Thresholds
	Balances   []AccountBalance
}

func selfSignerAccountDetail(address string) AccountDetail {
	return AccountDetail{
		ID: address,
		Signers: []Signer{
			{Key: address, Weight: 1, Type: "secp256k1_address"},
		},
	}
}

// GetSigners returns user signers
func (u *UserWallet) GetSigners(temp bool, gc *sharedconfig.GlobalConfig) (signers map[string]Signer) {
	account, _, err := u.GetBlockchainAccountDetail(temp, gc)
	if err != nil {
		return signers
	}
	signers = make(map[string]Signer)
	for _, v := range account.Signers {
		signers[v.Key] = v
	}
	return
}

// GetSignersWA returns user signers
func (u *User) GetSignersWA(account *AccountDetail) (signers map[string]Signer) {
	signers = make(map[string]Signer)
	for _, v := range account.Signers {
		signers[v.Key] = v
	}
	return
}

// GetSignersWA returns user signers
func (u *UserWallet) GetSignersWA(account *AccountDetail) (signers map[string]Signer) {
	signers = make(map[string]Signer)
	for _, v := range account.Signers {
		signers[v.Key] = v
	}
	return
}

// SignerIsValidWA checks if the signerKey is valid for this user public key
func (u *User) SignerIsValidWA(signerKey string, account *AccountDetail) bool {
	signer, ok := u.GetSignersWA(account)[signerKey]
	if !ok || signer.Weight < 1 {
		return false
	}

	return true
}

// SignerIsValidWA checks if the signerKey is valid for this user public key
func (u *UserWallet) SignerIsValidWA(signerKey string, account *AccountDetail) bool {
	signer, ok := u.GetSignersWA(account)[signerKey]
	if !ok || signer.Weight < 1 {
		return false
	}

	return true
}

// SignerIsValid checks if the signerKey is valid for this user public key
func (u *UserWallet) SignerIsValid(signerKey string, temp bool, gc *sharedconfig.GlobalConfig) bool {
	signer, ok := u.GetSigners(temp, gc)[signerKey]
	if !ok || signer.Weight < 1 {
		return false
	}

	return true
}

// SignerIsValid checks if the signerKey is valid for this user public key
func (u *User) SignerIsValid(signerKey string, temp bool, gc *sharedconfig.GlobalConfig) bool {
	for _, w := range u.UserWallets {
		if w.ID == w.Signer {
			signer, ok := w.GetSigners(temp, gc)[signerKey]
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
	balances = make(map[string]Balance, 0)
	depositAddresses := make([]CryptoWalletDepositAddress, 0)
	var nativeCode, nativeIssuer, nativeUsdPrice string
	nv := strings.Split(os.Getenv("USE_ASSET_FOR_NATIVE_PRICE"), ":")
	if len(nv) == 2 {
		nativeCode = nv[0]
		nativeIssuer = nv[1]
	}

	// get the ngn to usd rate. then get GAS naira price.
	cngnPrice := gc.GetCngnUsdRate()
	gasUsdPrice, _ := blockchain.GetGASDollarAskPrice(gc.DB)
	gasNativePrice := "1"
	// log.Println("gasUsdPrice", gasUsdPrice)
	cacheKey := fmt.Sprintf("GetBalance_%s", u.ID)
	if temp {
		if u.TempAddress != nil {

			cacheKey = fmt.Sprintf("GetBalance_%s", *u.TempAddress)
		}

	}
	{

		// search cache for balance
		ok, rawdata := gc.RedisCache.GetCachedResultRaw(cacheKey)

		if ok {

			// log.Printf("GetBalance[%v], served from cache\n", cacheKey)
			json.Unmarshal(rawdata, &balances)
			return
		}

	}

	account, _, err := u.GetBlockchainAccountDetail(temp, gc)
	if err != nil {
		qrCode := ""
		if !temp {

			p, e := dl.GeneratePaymentData(u.ID, "", "", "", "", gc)
			if e == nil {
				qrCode = p.QRCode
			}

		}
		balances[":"] = Balance{
			ContractAddress: "",
			AssetCode:       "",
			Amount:          decimal.Zero,
			QRCode:          qrCode,
			ImageURL:        os.Getenv("NATIVE_ASSET_IMAGE_URL"),
			UsdPrice:        gasUsdPrice,
			NativePrice:     gasNativePrice,
			InTrade: TradeLiabilties{
				SellingLiabilities: "0",
				BuyingLiabilities:  "0",
			},
			CryptoWalletDepositAddresses: depositAddresses,
		}

		if !temp && err.Error() == "error-blockchain-account-not-activated" {

			log.Printf("[GetBalance] get blockchain account detail error: %v\n", err)
			//save to cache
			gc.RedisCache.StoreResultToCacheRaw(cacheKey, balances, 60)
			//check if the cached result is nil then build default native placeholder

			return balances, nil
		}
		if !temp {
			return balances, nil
		}
		return balances, err
	}
	if len(nv) == 2 {
		checkCacheFirst := false
		if nativeCode != "CNGN" {
			checkCacheFirst = true
		}
		nairaAssetPrice, _, _ := blockchain.GetNairaPrice(nativeCode, nativeIssuer, gc, checkCacheFirst, true)
		if len(nairaAssetPrice) > 0 && nairaAssetPrice != "0" {
			// convert it to USD using current rate of naira
			nativeUsdPrice = decimal.RequireFromString(nairaAssetPrice).Mul(decimal.NewFromFloat(cngnPrice.Data.NgnToUsd)).Truncate(7).String()
		} else {
			nativeUsdPrice, _, _ = blockchain.GetDollarPrice(nativeCode, nativeIssuer, gc, checkCacheFirst)

		}
	} else {
		nativeUsdPrice = gasUsdPrice
	}

	var wg sync.WaitGroup
	var m sync.Mutex
	// var keys []string
	for _, v := range account.Balances {
		wg.Add(1)
		go func(bal AccountBalance) {
			bantuAsset := BantuAsset{AssetCode: bal.Code, ContractAddress: bal.Issuer}
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
			// buyingLiabilities, _ := decimal.NewFromString(bal.BuyingLiabilities)
			sellingLiabilities, _ := decimal.NewFromString(bal.SellingLiabilities)
			availableBalance := amount.Sub(sellingLiabilities)
			// availableBalance := amount.Sub(sellingLiabilities.Add(buyingLiabilities))
			// availableBalance := amount
			// availableBalance := availableBal.Truncate(7).String()

			if bal.Issuer != nativeIssuer && bal.Code != nativeCode {
				if len(bal.Code) == 0 {
					assetUsdPrice = gasUsdPrice
					assetNativePrice = gasNativePrice
				} else {
					checkCacheFirst := false
					if nativeCode != "CNGN" {
						checkCacheFirst = true
					}
					///////////////////////////////////////////////////////////////////////////////////////////////////
					nairaAssetPrice, _, _ := blockchain.GetNairaPrice(bal.Code, bal.Issuer, gc, checkCacheFirst, bantuAsset.IsEnabled(gc))
					if len(nairaAssetPrice) > 0 && nairaAssetPrice != "0" {
						// convert it to USD using current rate of naira
						nativeUsdPrice = decimal.RequireFromString(nairaAssetPrice).Mul(decimal.NewFromFloat(cngnPrice.Data.NgnToUsd)).Truncate(7).String()
					}
					assetNativePrice, _ = blockchain.GetNativeAskPrice(bal.Code, bal.Issuer, gc, checkCacheFirst, bantuAsset.IsEnabled(gc))
					dollarAsset := strings.Split(os.Getenv("DOLLAR_ASSET"), ":")
					if len(dollarAsset) == 2 {
						if strings.EqualFold(bal.Code, dollarAsset[0]) && strings.EqualFold(bal.Issuer, dollarAsset[1]) {
							//it is dollar asset
							assetUsdPrice = "1"

						} else if strings.HasPrefix(bal.Code, "USD") || strings.HasSuffix(bal.Code, "USD") {
							assetUsdPrice = "1"
						} else if bal.Code == "CNGN" {
							assetUsdPrice = decimal.NewFromFloat(cngnPrice.Data.UsdToNgn).Truncate(7).String()
							assetNativePrice = "1"
						} else {
							nativeUsdPriceDec := decimal.RequireFromString(nativeUsdPrice)
							nativePriceDec := decimal.RequireFromString(assetNativePrice)
							assetUsdPrice = nativePriceDec.Mul(nativeUsdPriceDec).Truncate(7).String()
						}
					} else {
						nativeUsdPriceDec := decimal.RequireFromString(nativeUsdPrice)
						nativePriceDec := decimal.RequireFromString(assetNativePrice)
						assetUsdPrice = nativePriceDec.Mul(nativeUsdPriceDec).Truncate(7).String()
					}

				}

			}
			if bal.Issuer == nativeIssuer && bal.Code == nativeCode {
				if len(nativeCode) == 0 {
					assetUsdPrice = gasUsdPrice
					assetNativePrice = gasNativePrice
				} else {
					assetUsdPrice = nativeUsdPrice
					assetNativePrice, _ = blockchain.GetNativeAskPrice(bal.Code, bal.Issuer, gc, true, bantuAsset.IsEnabled(gc))
				}

			}

			qrCode := ""
			if !temp {

				p, e := dl.GeneratePaymentData(u.ID, bal.Code, bal.Issuer, "", "", gc)
				if e == nil {
					qrCode = p.QRCode
				}
			}
			imageUrl := BantuAsset{AssetCode: bal.Code, ContractAddress: bal.Issuer}.GetAssetImage(gc)
			var quoteCurrency string
			var tokenizedAsset, fundingStructure, exitWithFiat int
			if gc.IsValidTokenizedAsset(bal.Code) {
				t := gc.GetTokenizedAssetByCode(bal.Code)
				assetNativePrice = decimal.NewFromFloat(t.PricePerToken).String()
				assetUsdPrice = decimal.NewFromFloat(t.PricePerToken).String()
				quoteCurrency = *t.AssetQuoteCurrency
				tokenizedAsset = 1
				fundingStructure = t.FundingStructure
				exitWithFiat = t.ExitWithFiat

			}
			balance := Balance{ContractAddress: bal.Issuer,
				AssetCode:        bal.Code,
				Amount:           availableBalance,
				QRCode:           qrCode,
				ImageURL:         imageUrl,
				UsdPrice:         assetUsdPrice,
				NativePrice:      assetNativePrice,
				QuoteCurrency:    quoteCurrency,
				TokenizedAsset:   tokenizedAsset,
				FundingStructure: fundingStructure,
				ExitWithFiat:     exitWithFiat,
				InTrade: TradeLiabilties{
					SellingLiabilities: bal.SellingLiabilities,
					BuyingLiabilities:  bal.BuyingLiabilities,
				},
			}
			if !temp {
				balance.CryptoWalletDepositAddresses = BantuAsset{AssetCode: bal.Code, ContractAddress: bal.Issuer}.GetDepositAddresses(u.ID, gc)
				//can deposit asset, now get the deposit addresses.
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

func (u *User) HasSharedAccessInAnyWallet(gc *sharedconfig.GlobalConfig) (enabled bool) {

	for _, wallet := range u.UserWallets {

		if wallet.SharedAccessEnabled == 1 {
			enabled = true
			break
		}

	}

	return
}

// GetNFTs gets user wallet blockchain NFT balance and return it as a map of assets  [code:issuer]Balance. Native key is [:]
func (u *UserWallet) GetNFTs(temp bool, gc *sharedconfig.GlobalConfig) (nfts []NFT, err error) {
	nfts = make([]NFT, 0)
	cacheKey := fmt.Sprintf("GetNFTs_%s", u.ID)
	if temp {
		if u.TempAddress != nil {
			cacheKey = fmt.Sprintf("GetNFTs_%s", *u.TempAddress)

		}

	}
	{

		// search cache for balance
		ok, rawdata := gc.RedisCache.GetCachedResultRaw(cacheKey)

		if ok {

			// log.Printf("GetNFTs[%v], served from cache\n", cacheKey)
			json.Unmarshal(rawdata, &nfts)
			return
		}

	}

	account, _, err := u.GetBlockchainAccountDetail(temp, gc)
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
		go func(v AccountBalance) {
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

			nft := NFT{ContractAddress: v.Issuer, AssetCode: v.Code,
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
	balances = make([]Balance, 0)
	//GetBalance
	unsortedBalances, err := u.GetBalance(temp, gc)
	if err != nil {
		// qrCode := ""
		// if !temp {

		// 	p, e := dl.GeneratePaymentData(u.ID, "", "", "", "", gc)
		// 	if e == nil {
		// 		qrCode = p.QRCode
		// 	}

		// }
		// gasUsdPrice, _ := blockchain.GetGASDollarAskPrice(gc.DB)
		// gasNativePrice := "1"
		// unsortedBalances[":"] = Balance{
		// 	ContractAddress: "",
		// 	AssetCode:   "",
		// 	Amount:      decimal.Zero,
		// 	QRCode:      qrCode,
		// 	ImageURL:    os.Getenv("NATIVE_ASSET_IMAGE_URL"),
		// 	UsdPrice:    gasUsdPrice,
		// 	NativePrice: gasNativePrice,
		// 	InTrade: TradeLiabilties{
		// 		SellingLiabilities: "0",
		// 		BuyingLiabilities:  "0",
		// 	},
		// 	CryptoWalletDepositAddresses: nil,
		// }
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
	// if len(balances) < 1 {
	// 	qrCode := ""
	// 	if !temp {

	// 		p, e := dl.GeneratePaymentData(u.ID, "", "", "", "", gc)
	// 		if e == nil {
	// 			qrCode = p.QRCode
	// 		}

	// 	}
	// 	gasUsdPrice, _ := blockchain.GetGASDollarAskPrice(gc.DB)
	// 	gasNativePrice := "1"
	// 	balances = append(balances, Balance{
	// 		ContractAddress: "",
	// 		AssetCode:   "",
	// 		Amount:      decimal.Zero,
	// 		QRCode:      qrCode,
	// 		ImageURL:    os.Getenv("NATIVE_ASSET_IMAGE_URL"),
	// 		UsdPrice:    gasUsdPrice,
	// 		NativePrice: gasNativePrice,
	// 		InTrade: TradeLiabilties{
	// 			SellingLiabilities: "0",
	// 			BuyingLiabilities:  "0",
	// 		},
	// 		CryptoWalletDepositAddresses: nil,
	// 	})
	// }
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
			if len(unclaimedBalance) > 0 {
				assetBalances.Unclaimed = unclaimedBalance
			} else {
				assetBalances.Unclaimed = make([]Balance, 0)
			}

			ml.Unlock()

		} else {
			assetBalances.Unclaimed = make([]Balance, 0)
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
			// else {
			// 	//default asset is returned on error
			// 	//Unclaimed Assets
			// 	ml.Lock()
			// 	assetBalances.Unclaimed = unclaimedBalance
			// 	ml.Unlock()
			// }

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
func (u *UserWallet) GetAccountThresholds(temp bool, gc *sharedconfig.GlobalConfig) (thresholds Thresholds) {
	account, _, err := u.GetBlockchainAccountDetail(temp, gc)
	if err != nil {
		return thresholds
	}

	return account.Thresholds
}

// fetchAccountDetail builds a Base AccountDetail for address: native
// balance plus every curated B20 asset's balance (there is no "list all
// balances an account holds" call on Base the way a Stellar account's own
// ledger entry provided - see AccountDetail's doc), cached the same way
// the original cached its Horizon account fetch.
func fetchAccountDetail(address string, gc *sharedconfig.GlobalConfig, cacheKey string) (clientAccount AccountDetail, destinationAccountExists bool, err error) {
	{
		ok, rawdata := gc.RedisCache.GetCachedResultRaw(cacheKey)
		if ok {
			json.Unmarshal(rawdata, &clientAccount)
			return clientAccount, true, nil
		}
	}

	client := network.GetBlockchainClient()
	exists, _, nativeBalance, _, _, err := network.BlockchainAccountProperties(client, address, basetxn.NativeAsset{})
	if err != nil {
		log.Printf("[fetchAccountDetail Network Failure]: %v\n", err)
		return clientAccount, false, &tErrors.ErrorTemporaryServerError{}
	}

	clientAccount = selfSignerAccountDetail(address)
	clientAccount.Balances = append(clientAccount.Balances, AccountBalance{Code: "", Issuer: "", Balance: nativeBalance.String(), SellingLiabilities: "0", BuyingLiabilities: "0"})

	curatedAssets := assetsDB.GetCuratedAssets(false, gc)
	for _, asset := range curatedAssets {
		bal, e := network.B20BalanceOf(client, asset.ContractAddress, address)
		if e != nil {
			continue
		}
		clientAccount.Balances = append(clientAccount.Balances, AccountBalance{Code: asset.AssetCode, Issuer: asset.ContractAddress, Balance: bal.String(), SellingLiabilities: "0", BuyingLiabilities: "0"})
	}

	cacheTimeStr := strings.TrimSpace(os.Getenv("BLOCKCHAIN_DATA_CACHE_LIFETIME"))
	if cacheTimeStr == "" {
		cacheTimeStr = "60" // balances change often, unlike Stellar's account-existence cache
	}
	cacheTime, _ := strconv.Atoi(cacheTimeStr)
	gc.RedisCache.StoreResultToCacheRaw(cacheKey, clientAccount, cacheTime)
	return clientAccount, exists, nil
}

// GetBlockchainAccountDetail fetches the bantu account information using public key
func (u *UserWallet) GetBlockchainAccountDetail(temp bool, gc *sharedconfig.GlobalConfig) (clientAccount AccountDetail, destinationAccountExists bool, err error) {
	address := u.ID
	cacheKey := fmt.Sprintf("bca_%v", u.ID)
	if temp {
		if u.TempAddress == nil {
			err = &tErrors.ErrorBlockchainAccountNotActivated{}
			return
		}
		address = *u.TempAddress
		cacheKey = fmt.Sprintf("bca_%v", *u.TempAddress)
	}
	return fetchAccountDetail(address, gc, cacheKey)
}

type MarketOfferWallet string

// OffersPage is the Base equivalent of Stellar's horizon.OffersPage - see
// GetMarketOffers' doc.
type OffersPage struct {
	Embedded struct {
		Records []MarketOfferRecord
	}
}

// GetMarketOffers fetches the market maker's available liquidity. See
// TokenizedAsset.GetMarketOffers' doc for the Base simplification this
// mirrors (Stellar's DEX order book has no Base equivalent).
func (u MarketOfferWallet) GetMarketOffers(gc *sharedconfig.GlobalConfig) (marketOffers OffersPage, err error) {
	return marketOffers, &tErrors.ErrorTemporaryServerError{}
}

// GetBlockchainAccountDetail fetches the bantu account information using public key
func (id UserWalletID) GetBlockchainAccountDetail(gc *sharedconfig.GlobalConfig) (clientAccount AccountDetail, destinationAccountExists bool, err error) {
	return fetchAccountDetail(string(id), gc, fmt.Sprintf("bca_%v", string(id)))
}

// isIssuerOfAssetCode reports whether issuer has an asset with assetCode
// registered in this app's own catalog. Stellar's version queried
// Horizon's global asset registry (client.Assets(ForContractAddress,
// ForAssetCode)); Base has no such registry, so "is this wallet the
// issuer of this asset" is answered from this backend's own curated/
// tokenized asset tables instead - the same tables that are already this
// app's actual source of truth for which B20 assets it recognizes.
func isIssuerOfAssetCode(issuer, assetCode string) bool {
	db := network.DB()
	if db == nil {
		return false
	}
	var count int64
	db.Table("curated_assets").Where("contract_address = ? AND asset_code = ?", strings.ToLower(issuer), strings.ToUpper(assetCode)).Count(&count)
	if count > 0 {
		return true
	}
	db.Table("tokenized_assets").Where("issuing_wallet_address = ? AND asset_code = ?", strings.ToLower(issuer), strings.ToUpper(assetCode)).Count(&count)
	return count > 0
}

// OwnerOfBlockchainAssetIssued fetches the blockchain asset information using public key
func (u *UserWallet) OwnerOfBlockchainAsset(assetCode string) bool {
	return isIssuerOfAssetCode(u.ID, assetCode)
}

// // CanIssueMoreAssets check if walet can issue more assets or has reached max limit
// func (u *UserWallet) CanIssueMoreAssets() bool {
// 	assetPage, err := u.GetBlockchainAssets()
// 	if err != nil {
// 		return false
// 	}
// 	maxCountAssets := 200

// OwnerOfBlockchainAssetIssued fetches the blockchain asset information using public key
func (u Issuer) OwnerOfBlockchainAsset(assetCode string) bool {
	return isIssuerOfAssetCode(string(u), assetCode)
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

// GetBlockchainAccountDataKey looked up values from a Stellar account's
// on-chain manage_data store. Base/EVM accounts have no equivalent store
// - see GetDataKey's doc in assets.go for the same simplification applied
// there.
func (u *UserWallet) GetBlockchainAccountDataKey(temp bool, gc *sharedconfig.GlobalConfig, keys ...string) (dataValues map[string]string) {
	return make(map[string]string)
}

func (u *User) BuildPrimaryWallet() {
	tempKP, _ := network.TempAccountKeypair(u.Address)
	var tempPK string
	if tempKP != nil {
		tempPK = tempKP.Address()
	}
	description := "Primary/Default wallet"
	userWallet := UserWallet{
		ID:            u.Address,
		TempAddress:   &tempPK,
		Description:   &description,
		Alias:         u.Username,
		Signer:        u.Address,
		UserID:        u.ID,
		PrimaryWallet: 1,
	}
	u.UserWallets = append(u.UserWallets, userWallet)
}

func (u *User) BuildNewSubWallet(subWalletAddress, walletTag, walletDescription string, walletType int, linkedWalletAddress string, gc *sharedconfig.GlobalConfig) (userWallet UserWallet, err error) {
	walletTag = strings.ToLower(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(walletTag, "_", ""), ".", ""), " ", ""), "%", ""))
	walletDescription = strings.TrimSpace(walletDescription)
	hasMMSubwallet := false
	hasBPSubWallet := false

	if len(linkedWalletAddress) > 0 && len(linkedWalletAddress) != 42 {
		log.Println("[BuildNewSubWallet] invalid parameters")
		return userWallet, &tErrors.CustomError{
			Param:      "linkedWalletAddress",
			Err:        "error-sub-wallet-parameters-invalid",
			ErrMessage: "Sub-wallet parameters are invalid. Ensure linkedWallet public key is 56 characters long.",
			Code:       http.StatusBadRequest,
		}
	}
	if len(subWalletAddress) != 42 || len(walletTag) == 0 {
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
	if strings.EqualFold(walletTag, "distribution") {
		log.Printf("[BuildNewSubWallet] wallet tag [%v] is an internal reserved tag. Not allowed.\n", walletTag)
		return userWallet, &tErrors.CustomError{
			Param:      "tag",
			Err:        "error-wallet-tag-not-allowed",
			ErrMessage: fmt.Sprintf("Sub-wallet tag [%v] is a reserved tag for internal use and not allowed to be used as a tag in sub wallet creation.", walletTag),
			Code:       http.StatusConflict,
		}
	}
	{
		//check to ensure sub-wallet does not already exist
		for _, wallet := range userWallets {
			if wallet.ID == subWalletAddress {
				log.Printf("[BuildNewSubWallet] wallet [%v] already exists in your account\n", subWalletAddress)
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
		_, errWallet := u.GetWalletByAddress(subWalletAddress, gc.DB)
		if errWallet != nil {
			if errWallet.Error() != "error-wallet-not-found" {
				log.Println("[BuildNewSubWallet] other service error ...", errWallet)

				return userWallet, errWallet
			}

		} else {
			//wallet already exists.
			log.Printf("[BuildNewSubWallet] wallet [%v] found in another account\n", subWalletAddress)

			return userWallet, &tErrors.CustomError{
				Param:      "id",
				Err:        "error-sub-wallet-already-exists-with-another-account",
				ErrMessage: "Sub-wallet already exists with another account",
				Code:       http.StatusConflict,
			}
		}
	}
	tempKP, pErr := network.TempAccountKeypair(subWalletAddress)
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
		ID:          subWalletAddress,
		TempAddress: &tempPK,
		Tag:         &walletTag,
		Description: &walletDescription,
		Alias:       alias,
		Signer:      u.PrimarySigner,
		UserID:      u.ID,
		WalletType:  walletType,
	}
	if len(linkedWalletAddress) > 0 {
		userSubWallet.LinkedWalletAddress = &linkedWalletAddress
	}

	return userSubWallet, nil
}

func (uw *UserWallet) BuildNewLinkedSubWallet(owner *User, gc *sharedconfig.GlobalConfig) (userWallet UserWallet, err error) {
	if uw.LinkedWalletAddress == nil {
		return userWallet, &tErrors.CustomError{
			Param:      "id",
			Err:        "error-linked-wallet-public-key-invalid",
			ErrMessage: "Linked-wallet public key is invalid",
			Code:       http.StatusBadRequest,
		}
	}
	if *uw.LinkedWalletAddress == "" {
		return userWallet, &tErrors.CustomError{
			Param:      "id",
			Err:        "error-linked-wallet-public-key-invalid",
			ErrMessage: "Linked-wallet public key is invalid",
			Code:       http.StatusBadRequest,
		}
	}

	var walletTag, walletDescription string

	if uw.WalletType == 1 {
		walletTag = *uw.Tag + "-distribution"
		walletDescription = "distribution wallet for issuer wallet" + uw.Alias

	}

	if len(walletDescription) == 0 {
		walletDescription = walletTag
	}
	userWallets := owner.GetAllWallets(gc)
	{
		//check to ensure sub-wallet does not already exist
		for _, wallet := range userWallets {
			if wallet.ID == *uw.LinkedWalletAddress {
				log.Printf("[BuildNewLinkedSubWallet] wallet [%v] already exists in your account\n", *uw.LinkedWalletAddress)
				return userWallet, &tErrors.CustomError{
					Param:      "id",
					Err:        "error-sub-wallet-already-exists-in-your-account",
					ErrMessage: "Sub-wallet already exists in your account",
					Code:       http.StatusConflict,
				}
			}
			// if wallet.WalletType == 2 {
			// 	hasMMSubwallet = true
			// }
			// if wallet.WalletType == 3 {
			// 	hasBPSubWallet = true
			// }

			if wallet.Tag != nil {
				if strings.EqualFold(*wallet.Tag, walletTag) {
					log.Printf("[BuildNewLinkedSubWallet] wallet tag [%v] already exists in your account\n", walletTag)
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
		_, errWallet := owner.GetWalletByAddress(*uw.LinkedWalletAddress, gc.DB)
		if errWallet != nil {
			if errWallet.Error() != "error-wallet-not-found" {
				log.Println("[BuildNewLinkedSubWallet] other service error ...", errWallet)

				return userWallet, errWallet
			}

		} else {
			//wallet already exists.
			log.Printf("[BuildNewLinkedSubWallet] wallet [%v] found in another account\n", uw.LinkedWalletAddress)

			return userWallet, &tErrors.CustomError{
				Param:      "id",
				Err:        "error-sub-wallet-already-exists-with-another-account",
				ErrMessage: "Sub-wallet already exists with another account",
				Code:       http.StatusConflict,
			}
		}
	}
	tempKP, pErr := network.TempAccountKeypair(*uw.LinkedWalletAddress)
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

	alias := fmt.Sprintf("%s_%s", owner.Username, walletTag)
	userSubWallet := UserWallet{
		ID:          *uw.LinkedWalletAddress,
		TempAddress: &tempPK,
		Tag:         &walletTag,
		Description: &walletDescription,
		Alias:       alias,
		Signer:      owner.PrimarySigner,
		UserID:      owner.ID,
		WalletType:  0,
	}
	return userSubWallet, nil
}

func (lw LinkedWalletAddress) String() string {
	return string(lw)
}
func (lw LinkedWalletAddress) IsValid(gc *sharedconfig.GlobalConfig) (w UserWallet, valid bool) {
	e := gc.DB.Where("linked_wallet_address = ?", string(lw)).First(&w).Error
	if e == nil {
		return w, true
	}
	return w, false
}
func (u UserWallet) IsValidLinkedWallet(gc *sharedconfig.GlobalConfig) (w UserWallet, valid bool) {
	e := gc.DB.Where("linked_wallet_address = ?", u.ID).First(&w).Error
	if e == nil {
		return w, true
	}
	return w, false
}

func (lw *LinkedWalletAddress) BuildNewLinkedSubWallet(owner *User, uw *UserWallet, gc *sharedconfig.GlobalConfig) (userWallet UserWallet, err error) {

	return uw.BuildNewLinkedSubWallet(owner, gc)

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

			// log.Printf("GetOwner[%v], served from cache\n", cacheKeyInfo)
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
	db.Preload(clause.Associations).Where("wallet_address = ?", string(id)).Find(&accessList)

	return
}

func (id UserWalletID) UserHasAccess(username, permissionToCheck string, db *gorm.DB) (hasAccess bool) {
	permissionToCheck = strings.ToUpper(strings.TrimSpace(permissionToCheck))
	if permissionToCheck == "" {
		return
	}
	ap := id.GetPermissionList(db)
	for _, p := range ap {
		if p.TargetUsername == username && p.Permission == permissionToCheck {
			hasAccess = true
			break
		}

	}
	return
}

func (id UserWalletID) UserWithInitiatorAccess(username string, db *gorm.DB) (hasAccess bool) {
	permissionToCheck := "INITIATOR"

	ap := id.GetPermissionList(db)
	for _, p := range ap {
		if p.TargetUsername == username && p.Permission == permissionToCheck {
			hasAccess = true
			break
		}

	}
	return
}

func (id UserWalletID) UserWithApproverAccess(username string, db *gorm.DB) (hasAccess bool) {
	permissionToCheck := "APPROVER"

	ap := id.GetPermissionList(db)
	for _, p := range ap {
		if p.TargetUsername == username && p.Permission == permissionToCheck {
			hasAccess = true
			break
		}

	}
	return
}
func (id UserWalletID) UserWithViewOnlyAccess(username string, db *gorm.DB) (hasAccess bool) {
	permissionToCheck := "VIEW-ONLY"

	ap := id.GetPermissionList(db)
	for _, p := range ap {
		if p.TargetUsername == username && p.Permission == permissionToCheck {
			hasAccess = true
			break
		}

	}
	return
}

func (u *UserWallet) GetPermissionList(db *gorm.DB) (accessList []WalletPermission) {

	if u.Permissions != nil {

		if len(u.Permissions) > 0 {
			return u.Permissions
		}
	}
	accessList = make([]WalletPermission, 0)
	db.Preload(clause.Associations).Where("wallet_address = ?", u.ID).Find(&accessList)

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

func (u *UserWallet) AddressHasViewOnlyAccess(gc *sharedconfig.GlobalConfig) (viewOnly bool) {
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

			// log.Printf("GetWalletOwner[%v], served from cache\n", cacheKeyInfo)
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

func (u *UserWallet) GetSwapFee(gc *sharedconfig.GlobalConfig) (serviceFee ServiceFee) {
	// var serviceFee ServiceFee
	gc.DB.Where("id = ? AND inactive = 0", "SWAP_FEE").First(&serviceFee)
	if serviceFee.Inactive == 0 {

		//TODO: check if user has zero swap fees and modify the swap fee

	}
	serviceFee.FeeWalletSecretKey = os.Getenv("SWAP_FEE_WALLET")

	return
}

func (u *UserWallet) GetPaymentFeeWallet(gc *sharedconfig.GlobalConfig) (wallet string) {
	var serviceFee ServiceFee
	gc.DB.Where("id = ? AND inactive = 0", "PAYMENT_FEE").First(&serviceFee)
	if serviceFee.Inactive == 0 {

		//TODO: check if user has zero swap fees and modify the swap fee

	}
	serviceFee.FeeWalletSecretKey = os.Getenv("PAYMENT_FEE_WALLET")

	return serviceFee.FeeWalletSecretKey
}

func (u *UserWallet) GetVATWallet(gc *sharedconfig.GlobalConfig) string {

	return gc.GetVATWallet()
}

func (u *UserWallet) GetVATValue(serviceFee decimal.Decimal, gc *sharedconfig.GlobalConfig) (vat float64) {
	// var serviceFee ServiceFee
	vat = gc.GetVATValue(serviceFee)
	return
}

func (u *UserWallet) GetSharedAccessPaymentFee(gc *sharedconfig.GlobalConfig) (serviceFee ServiceFee) {
	// var serviceFee ServiceFee
	if u.SharedAccessEnabled == 0 {
		return
	}
	gc.DB.Where("id = ? AND inactive = 0", "SHARED_ACCESS_PAYMENT_FEE").First(&serviceFee)
	if serviceFee.Inactive == 0 {

		//TODO: check if user has zero fees and modify the fee

	}
	serviceFee.FeeWalletSecretKey = os.Getenv("SHARED_ACCESS_PAYMENT_FEE_WALLET")

	return
}

func (u *UserWallet) GetPatronFee(gc *sharedconfig.GlobalConfig) (serviceFee ServiceFee) {
	// var serviceFee ServiceFee

	gc.DB.Where("id = ? AND inactive = 0", "PATRON_FEE").First(&serviceFee)
	if serviceFee.Inactive == 0 {

		//TODO: check if user has zero fees and modify the fee

	}
	serviceFee.FeeWalletSecretKey = os.Getenv("PATRON_FEE_WALLET")

	return
}

func (u *UserWallet) GetAccountRecoveryFee(gc *sharedconfig.GlobalConfig) (serviceFee ServiceFee) {
	// var serviceFee ServiceFee

	gc.DB.Where("id = ? AND inactive = 0", "ACCOUNT_RECOVERY_FEE").First(&serviceFee)
	if serviceFee.Inactive == 0 {

		//TODO: check if user has zero fees and modify the fee

	}
	serviceFee.FeeWalletSecretKey = os.Getenv("ACCOUNT_RECOVERY_FEE_WALLET")

	return
}

func (u *UserWallet) GetSubwalletCreationFee(gc *sharedconfig.GlobalConfig) (serviceFee ServiceFee) {
	// var serviceFee ServiceFee

	gc.DB.Where("id = ? AND inactive = 0", "SUBWALLET_CREATION_FEE").First(&serviceFee)
	if serviceFee.Inactive == 0 {

		//TODO: check if user has zero fees and modify the fee

	}
	// SUBWALLET_CREATION_FEE_WALLET
	serviceFee.FeeWalletSecretKey = os.Getenv("SUBWALLET_CREATION_FEE_WALLET")
	return
}

func (u *UserWallet) GetTokenizationApplicationFee(gc *sharedconfig.GlobalConfig) (serviceFee ServiceFee) {
	// var serviceFee ServiceFee

	gc.DB.Where("id = ? AND inactive = 0", "TOKENIZATION_APPLICATION_FEE").First(&serviceFee)
	if serviceFee.Inactive == 0 {

		//TODO: check if user has zero fees and modify the fee

	}
	serviceFee.FeeWalletSecretKey = os.Getenv("TOKENIZATION_APPLICATION_FEE_WALLET")

	return
}

func (t *TokenizedAsset) GetTokenizationFeeWallet(gc *sharedconfig.GlobalConfig) (serviceFee ServiceFee) {
	// var serviceFee ServiceFee

	gc.DB.Where("id = ? AND inactive = 0", "TOKENIZATION_FEE").First(&serviceFee)
	if serviceFee.Inactive == 0 {

		//TODO: check if user has zero fees and modify the fee

	}
	serviceFee.FeeWalletSecretKey = os.Getenv("TOKENIZATION_FEE_WALLET")

	return
}

func (u *UserWallet) GetClosedGroupFee(gc *sharedconfig.GlobalConfig) (serviceFee ServiceFee) {
	// var serviceFee ServiceFee

	gc.DB.Where("id = ? AND inactive = 0", "CLOSED_GROUP_FEE").First(&serviceFee)
	if serviceFee.Inactive == 0 {

		//TODO: check if user has zero fees and modify the fee

	}
	serviceFee.FeeWalletSecretKey = os.Getenv("CLOSED_GROUP_FEE_WALLET")

	return
}

func (u *UserWallet) GetServiceFee(serviceFeeID string, gc *sharedconfig.GlobalConfig) (serviceFee ServiceFee) {
	// var serviceFee ServiceFee
	gc.DB.Where("id = ? AND inactive = 0", serviceFeeID).First(&serviceFee)
	if serviceFee.Inactive == 0 {

		//TODO: check if user has zero swap fees and modify the swap fee

	}

	return
}

func (u *UserWallet) GetActivationFee(activationFeeID string, gc *sharedconfig.GlobalConfig) (activationAmount ActivationAmount) {
	// var serviceFee ServiceFee
	gc.DB.Where("id = ? AND inactive = 0", activationFeeID).First(&activationAmount)
	if activationAmount.Inactive == 0 {

		//TODO: check if user has zero swap fees and modify the swap fee

	}

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

			// log.Printf("GetWallet[%v], served from cache\n", cacheKeyInfo)
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
				Address: string(id),
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

			// log.Printf("GetWallet[%v], served from cache\n", cacheKeyInfo)
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
				Address: string(a),
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

func (id UserWalletID) AddressHasViewOnlyAccess(gc *sharedconfig.GlobalConfig) (viewOnly bool) {
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

func (u *User) GetCuratedSwapList(gc *sharedconfig.GlobalConfig) (list []assets.CuratedSwapAsset) {
	defaultAssetImageURL := os.Getenv("DEFAULT_ASSET_IMAGE_URL")
	list = make([]assets.CuratedSwapAsset, 0)
	// var swapAsset assets.CuratedSwapAsset
	am := assetsDB.GetCuratedAssets(false, gc)

	for _, a := range am {
		item := assets.CuratedSwapAsset{
			AssetCode:                   a.AssetCode,
			ContractAddress:             a.ContractAddress,
			AssetName:                   a.AssetName,
			Description:                 a.Description,
			Website:                     a.Website,
			AssetConditions:             a.AssetConditions,
			AssetLimit:                  a.AssetLimit,
			AssetRedemptionInstructions: a.AssetRedemptionInstructions,
			ContactEmail:                a.ContactEmail,
			AssetClassID:                a.AssetClassID,
			AssetClass:                  a.AssetClass,
			Organization:                a.Organization,
			Withdrawable:                a.Withdrawable,
			DecimalPlaces:               a.DecimalPlaces,
			GenerateDepositAddress:      a.GenerateDepositAddress,
		}
		if a.RealAssetImageURL != nil {
			item.RealAssetImageURL = *a.RealAssetImageURL
		}
		if a.ClosedGroup != nil {
			item.ClosedGroup = *a.ClosedGroup
		}
		if a.ImageURL == nil {
			item.ImageURL = defaultAssetImageURL
		} else {
			item.ImageURL = *a.ImageURL
		}
		list = append(list, item)
	}

	//add the tokenized assets to the list
	return list
}

func (u *User) GetPatronMembership(gc *sharedconfig.GlobalConfig) (membership UserPatronMembership, err error) {
	// var membership UserPatronMembership
	e := gc.DB.Preload(clause.Associations).Where("username = ?", u.Username).First(&membership).Error
	if e != nil {
		log.Printf("[GetPatronMembership] error : %v\n", e)
		err = &tErrors.ErrorTemporaryServerError{}
		return
	}
	return membership, nil
}
func (pid PatronPackageID) GetPatronPackage(gc *sharedconfig.GlobalConfig) (patronPackage PatronPackage, err error) {
	err = gc.DB.Where("id = ?", strings.ToUpper(string(pid))).First(&patronPackage).Error
	return
}

func (pid PatronTierID) GetPatronTier(gc *sharedconfig.GlobalConfig) (patronTier PatronTier, err error) {
	err = gc.DB.Where("id = ?", strings.ToUpper(string(pid))).First(&patronTier).Error
	return
}

func (pid PatronMembershipGradeID) GetPatronMemberShipConfig(gc *sharedconfig.GlobalConfig) (membershipConfig PatronMembershipGrade, err error) {
	e := gc.DB.Where("id = ?", uint64(pid)).First(&membershipConfig).Error
	if e != nil {
		if errors.Is(e, gorm.ErrRecordNotFound) {
			//no approval request was found
			err = &tErrors.ErrorInvalidRequest{
				ID: fmt.Sprintf("%v", uint64(pid)),
			}
			return
		}
		err = &tErrors.ErrorTemporaryServerError{}
	}
	return membershipConfig, nil
}

func (u *User) GetPatronSubscriptionLogs(gc *sharedconfig.GlobalConfig) (patronSubLogs []UserPatronSubscriptionLog) {
	patronSubLogs = make([]UserPatronSubscriptionLog, 0)
	gc.DB.Order("created_at DESC").Where("username = ?", u.Username).Find(&patronSubLogs)
	return
}

func (u *User) GetReferralCount(gc *sharedconfig.GlobalConfig) (referralCount int64) {
	gc.DB.Where("referrer = ?", u.Username).Count(&referralCount)
	return
}

func (d CryptoDepositAddress) GetDetail(currency string, gc *sharedconfig.GlobalConfig) (depositAddress CryptoWalletDepositAddress, err error) {
	address := string(d)
	e := gc.DB.Where("upper(currency) = upper(?) AND deposit_address = ?", currency, address).First(&depositAddress).Error
	if e != nil {
		log.Printf("Error fetching deposit address, error: %v\n", e)
		err = &tErrors.ErrorTemporaryServerError{}
		return
	}
	return depositAddress, nil

}
func (callback *CallbackDepositItem) Save(gc *sharedconfig.GlobalConfig) (err error) {

	e := gc.DB.Omit(clause.Associations).Create(callback).Error
	if e != nil {
		log.Println("[SAVE CALLBACK]error creating callback: ", e)
		//notify failure
		{
			// procedure to notify
		}
		err = &tErrors.ErrorTemporaryServerError{}
		return
	}
	return nil

}

func (callbackObj *CallbackDeposit) SaveDepositCallback(gc *sharedconfig.GlobalConfig) (err error) {
	if !strings.EqualFold(callbackObj.Event, "WALLET_DEPOSIT_COMPLETED") {
		log.Println("[SAVE CALLBACK]error NOT a completed deposit.", callbackObj.Event)
		err = &tErrors.ErrorTemporaryServerError{}
		return
	}
	_, err = CryptoDepositAddress(callbackObj.Data.ToAddress).GetDetail(callbackObj.Data.Currency, gc)
	if err != nil {
		log.Printf("[DEPOSIT CALLBACK FAILED RETRIEVAL] ERROR GETTING DEPOSIT ADDRESS OBJECT FROM DB for [%v, %v, %v], error: [%v]\n", callbackObj.Data.Currency, callbackObj.Data.Network, callbackObj.Data.ToAddress, err)
		//notify failure
		{
			// procedure to notify
		}

		return err
	}
	e := gc.DB.Omit(clause.Associations).Create(&callbackObj.Data).Error
	if e != nil {
		log.Println("[SAVE CALLBACK]error creating callback: ", e)
		//notify failure
		{
			// procedure to notify
		}
		err = &tErrors.ErrorTemporaryServerError{}
		return
	}
	return nil

}

func (u *User) GetDownlines(gc *sharedconfig.GlobalConfig) (lv1, lv2, lv3 []Downline) {
	lv1 = make([]Downline, 0)
	lv1Referrals := make([]string, 0)

	lv2 = make([]Downline, 0)
	lv2Referrals := make([]string, 0)

	lv3 = make([]Downline, 0)
	// lv3Referrals := make([]string, 0)

	gc.DB.Model(User{}).Where("referrer = ?", u.Username).Select("username", "referrer").Find(&lv1)

	//get all the user referrals
	gc.DB.Model(User{}).Where("referrer = ?", u.Username).Pluck("username", &lv1Referrals)

	if len(lv1) > 0 && len(lv1Referrals) > 0 {
		//check for lv2
		gc.DB.Model(User{}).Where("referrer IN ?", lv1Referrals).Select("username", "referrer").Find(&lv2)
		gc.DB.Model(User{}).Where("referrer IN ?", lv1Referrals).Pluck("username", &lv2Referrals)

		if len(lv2) > 0 && len(lv2Referrals) > 0 {
			//check for lv2
			gc.DB.Model(User{}).Where("referrer IN ?", lv2Referrals).Select("username", "referrer").Find(&lv3)
			// gc.DB.Model(User{}).Where("referrer IN ?", lv2Referrals).Pluck("username",&lv3Referrals)

		}
	}

	return
}

func (u *User) GetUplines(gc *sharedconfig.GlobalConfig) (lv1, lv2, lv3 string) {

	if u.Referrer != nil {
		lv1 = *u.Referrer
	} else {
		return
	}

	//get all the user referrals
	lv1User, err := Username(lv1).GetFullUser(gc.DB, gc)

	if err != nil {
		return
	}
	if lv1User.Referrer != nil {
		lv2 = *lv1User.Referrer
	} else {
		return
	}

	//get the user referrals
	lv2User, err := Username(lv2).GetFullUser(gc.DB, gc)

	if err != nil {
		return
	}
	if lv2User.Referrer != nil {
		lv3 = *lv2User.Referrer
	}

	return
}

func (u *User) GetWalletByAddress(publicKey string, db *gorm.DB) (wallet UserWallet, err error) {
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
		log.Printf("[GetWalletByAddress] error: %s\n", e)
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

			// log.Printf("GetWalletOwner[%v], served from cache\n", cacheKeyInfo)
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
	e := db.Where("target_username = ? AND wallet_address = ? AND permission = ?", username, string(id), strings.ToUpper(permission)).First(&walletPermission).Error
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

// GetFullUser get user by username or id
func (u Username) GetFullUser(db *gorm.DB, gc *sharedconfig.GlobalConfig) (owner User, err error) {
	cacheKeyInfo := fmt.Sprintf("userObj %v", string(u))

	{

		// search cache for balance
		ok, rawdata := gc.RedisCache.GetCachedResultRaw(cacheKeyInfo)

		if ok {

			// log.Printf("GetFullUser[%v], served from cache\n", cacheKeyInfo)
			json.Unmarshal(rawdata, &owner)
			return
		}

	}

	e := db.Preload("UserWallets.Permissions").Preload(clause.Associations).Where("(username = ? OR id= ?)", string(u), string(u)).First(&owner).Error
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

			// log.Printf("GetSimpleUser[%v], served from cache\n", cacheKeyInfo)
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

func (u Username) GetUserPermissionOnWallet(walletAddress string, db *gorm.DB) (walletPermission WalletPermission, err error) {
	e := db.Where("target_username = ? AND wallet_address = ?", string(u), walletAddress).First(&walletPermission).Error
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

// GetOpenTokenizedAssetByInitiatorUsername get the tokenization that has status 0 or 1 initiated by the initiator username
func (u Username) GetOpenTokenizedAssetByInitiatorUsername(db *gorm.DB) (tokenizedAsset TokenizedAsset, NotFound bool, err error) {
	// var ta userModels.TokenizedAsset
	initiatorUsername := u.String()
	err = db.Preload(clause.Associations).Order("asset_tokenization_status ASC").Where("asset_tokenization_status < 1 AND initiator_username = ?", initiatorUsername).First(&tokenizedAsset).Error

	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			//critical database error occured
			log.Printf("[GetTokenizedAssetByID]error fetching existing tokenization with initiatorUsername %v from database  [%v]", initiatorUsername, err)
			return

		} else {
			//record not found
			NotFound = true
			err = &tErrors.CustomError{Param: "tokenizationID", Err: "error-invalid-tokenizationId", ErrMessage: fmt.Sprintf("%v has no tokenized asset inititated", initiatorUsername)}
			return
		}
	}

	return
}

// GetFeeReadyTokenizedAssetApplicationByInitiatorUsername get the tokenization that has status 1 and initiated by the initiator username
func (u Username) GetFeeReadyTokenizedAssetApplicationByInitiatorUsername(db *gorm.DB) (tokenizedAsset TokenizedAsset, NotFound bool, err error) {
	initiatorUsername := u.String()
	err = db.Preload(clause.Associations).Order("Created_At DESC").Where("asset_tokenization_status = 1 AND initiator_username = ?", initiatorUsername).First(&tokenizedAsset).Error

	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			//critical database error occured
			log.Printf("[GetTokenizedAssetByID]error fetching existing tokenization with initiatorUsername %v from database  [%v]", initiatorUsername, err)
			return

		} else {
			//record not found
			NotFound = true
			err = &tErrors.CustomError{Param: "tokenizationID", Err: "error-invalid-tokenizationId", ErrMessage: fmt.Sprintf("%v has no tokenized asset inititated", initiatorUsername)}
			return
		}
	}

	return
}
func (u *User) HasAccessToAddress(publicKey string, gc *sharedconfig.GlobalConfig) (hasAccess bool) {
	// walletPermissions := u.Fetch3rdPartyWalletPermissions(gc)
	// if len(walletPermissions) == 0 {
	// 	return false
	// }
	for _, walletAccess := range u.WalletsSharedWithUser {
		if walletAccess.WalletAddress == publicKey {
			return true
		}
	}

	return false
}

// func (u *User) IsTokenizedAsset(code, issuer string, gc *sharedconfig.GlobalConfig) bool {

// 	e := gc.DB.Where("Asset_Tokenization_Status > 4 AND Asset_Code = upper(?) AND Issuing_Wallet_Public_Key = upper(?)", code, issuer).First(&TokenizedAsset{}).Error

// 	return e == nil
// }

func (u *User) HasPassedKYC() bool {

	return u.KYCVerified > 0
}

func (u *User) GetKYCLevel() int {

	return u.KYCVerified
}

func (u *User) GetKycData(gc *sharedconfig.GlobalConfig) (kycData UserKyc, err error) {

	e := gc.DB.Preload(clause.Associations).Where("user_id = ?", u.ID).First(&kycData).Error
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
		WalletAddress:         wp.WalletAddress,
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
		wallet, err := UserWalletID(assignedPermission.WalletAddress).GetWallet(gc.DB, gc)
		if err != nil {
			return
		}

		thirdPartyWallet.WalletAddress = wallet.ID
		thirdPartyWallet.WalletAlias = wallet.Alias
		if wallet.Description != nil {
			thirdPartyWallet.WalletDescription = *wallet.Description
		}
		thirdPartyWallet.AssetBalances, _ = wallet.GetWalletAssetBalances(gc)
		//use it to fetch wallet owner details
		owner, err := UserWalletID(assignedPermission.WalletAddress).GetWalletOwner(gc.DB, gc)
		if err != nil {
			return
		}

		thirdPartyWallet.Owner = owner.Username

		if assignedPermission.Permission == "INITIATOR" {
			//set walletSettings SINCE INITIATORS CAN MODIFY WALLET
			viewOnlyAccess := true
			hasApprover := false
			if wallet.SharedAccessEnabled == 0 {
				viewOnlyAccess = false
			}
			walletSettings := WalletSettings{
				NumberOfApprovalsNeeded: wallet.NumberOfApprovalsNeeded,
				WalletType:              wallet.WalletType,
			}
			for _, p := range wallet.Permissions {
				if p.Permission != "VIEW-ONLY" {
					viewOnlyAccess = false
				}
				if p.Permission == "APPROVER" {
					hasApprover = true
				}
				walletSettings.Permissions = append(walletSettings.Permissions, p.ToJSON(gc))
			}
			if viewOnlyAccess {
				walletSettings.WalletThreshold = 1
			} else if hasApprover {
				walletSettings.WalletThreshold = 2
			} else {
				walletSettings.WalletThreshold = 0
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

func (w *UserWallet) HasInitiatorPermissionToAddress(ownerSignerAddress string, gc *sharedconfig.GlobalConfig) (hasAccess bool) {
	user, err := UserSigner(ownerSignerAddress).GetOwner(gc.DB, gc)

	if err != nil {
		return false
	}
	walletPermissions := user.FetchWalletsPermissionsSharedWithUser(gc)
	if len(walletPermissions) == 0 {
		return false
	}
	for _, walletAccess := range walletPermissions {
		if walletAccess.WalletAddress == w.ID && walletAccess.Permission == "INITIATOR" {
			return true
		}
	}

	return false
}

func (w *UserWallet) SignerHasInitiatorPermissionToAddress(signerOwner User, gc *sharedconfig.GlobalConfig) (hasAccess bool) {

	walletPermissions := signerOwner.FetchWalletsPermissionsSharedWithUser(gc)
	if len(walletPermissions) == 0 {
		return false
	}
	for _, walletAccess := range walletPermissions {
		if walletAccess.WalletAddress == w.ID && walletAccess.Permission == "INITIATOR" {
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
			// log.Printf("GetDefaultAssets [%v], served from cache\n", cacheKey)
			da := response.([]interface{})
			for _, v := range da {
				vals := v.(map[string]interface{})
				// log.Printf("VALS: [%+v]\n", vals)
				defaultAssets = append(defaultAssets, DefaultAsset{
					AssetCode:       vals["assetCode"].(string),
					ContractAddress: vals["contractAddress"].(string),
					ImageURL:        vals["imageUrl"].(string),
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

func (u *User) IsEnterpriseProfile(gc *sharedconfig.GlobalConfig) bool {
	var result string
	e := gc.DB.Table("service_links").Select("id").Where("username = ?", u.Username).Scan(&result).Error
	if e != nil {
		gc.LogDiscordFailedRequest(fmt.Sprintf("[IsEnterpriseProfile] Error verifying if %v is enterprise profile: %v", u.Username, e))
	}
	return len(result) > 0

}
func (u *User) BelongsToAnEnterpriseProfile() bool {
	if u.CreatedByServiceLinkID != nil {
		if len(*u.CreatedByServiceLinkID) > 5 {
			return true
		}
	}
	return false

}

func (u *User) GetFiatActiationAmount(gc *sharedconfig.GlobalConfig) (activationAmount, trovPercent float64) {
	cc := CountryCode(*u.CountryCode).GetConfig(gc)

	// On Base every address "exists" the moment it's derived, unlike
	// Stellar where an account only existed once funded - so "already
	// activated" is judged by native GAS balance against the same
	// STANDARD_WALLET_MINIMUM_BALANCE threshold used elsewhere (e.g.
	// market_making.go, asset_trusts_and_claims.go), not account
	// existence.
	client := network.GetBlockchainClient()
	_, _, nativeBalance, _, _, err := network.BlockchainAccountProperties(client, u.Address, basetxn.NativeAsset{})
	if minBalance, mErr := decimal.NewFromString(os.Getenv("STANDARD_WALLET_MINIMUM_BALANCE")); err == nil && mErr == nil {
		if nativeBalance.GreaterThanOrEqual(minBalance) {
			return 0, cc.TrovTokenActivationPercent
		}
	}
	// get the country fiat

	return cc.FiatActivationAmount, cc.TrovTokenActivationPercent
}
func (u *ServiceLinksUser) SendPushMessage(title, body, imageURI string, dataPayload map[string]string, gc *sharedconfig.GlobalConfig) {
	//Send push notification to user
	// log.Println(title, body)

	if u.PushNotificationToken == nil {
		return
	}

	pns.SendFirebaseMessage(*u.PushNotificationToken, title, body, imageURI, dataPayload, gc.PushNotificationClient, gc.PNSContext)

}

func (u *User) ToServiceLinkUser(gc *sharedconfig.GlobalConfig) (slUser ServiceLinksUser) {

	if u == nil {
		return
	}
	slUser = ServiceLinksUser{
		CreatedAt:             u.CreatedAt,
		UpdatedAt:             u.UpdatedAt,
		LastUpdatedMobileOn:   u.LastUpdatedMobileOn,
		ID:                    u.ID,
		Username:              u.Username,
		Email:                 u.Email,
		ImageThumbnailURL:     u.ImageThumbnailURL,
		FirstName:             u.FirstName,
		LastName:              u.LastName,
		Mobile:                u.Mobile,
		Address:               u.Address,
		PrimarySigner:         u.PrimarySigner,
		PushNotificationToken: u.PushNotificationToken,
		Corporate:             u.Corporate,
		MobileVerified:        u.MobileVerified,
		KYCVerified:           u.KYCVerified,
		Suspended:             u.Suspended,
	}
	return
}

func (u *ServiceLinksUser) ToServiceLinkUserInfo(gc *sharedconfig.GlobalConfig) (slUser ServiceLinkUserInfo) {

	if u == nil {
		return
	}
	slUser = ServiceLinkUserInfo{
		UserData: *u,
	}
	return
}

func (w *UserWallet) InvalidateUserCache(gc *sharedconfig.GlobalConfig) {
	userAccount, err := w.GetWalletOwner(gc.DB, gc)
	if err != nil {
		return
	}

	cacheKey1 := fmt.Sprintf("GetBalance_%s", userAccount.Address)
	cacheKeyUsername := fmt.Sprintf("userObj %v", userAccount.Username)
	cacheKeyEmail := fmt.Sprintf("userObj %v", userAccount.Email)
	cacheKeySigner := fmt.Sprintf("userObj %v", userAccount.PrimarySigner)
	cacheKeyUserID := fmt.Sprintf("userObj %v", userAccount.ID)
	gc.RedisCache.DeleteFromCache(cacheKeyUsername, cacheKeyEmail, cacheKeySigner, cacheKeyUserID)

	gc.RedisCache.DeleteFromCache(cacheKey1)
	userAccount.InvalidateUserWalletCache(gc)
}
func (u Username) InvalidateUserCache(gc *sharedconfig.GlobalConfig) {
	userAccount, err := u.GetFullUser(gc.DB, gc)
	if err != nil {
		return
	}
	userAccount.InvalidateUserCache(gc)
}
func (u *User) InvalidateUserCache(gc *sharedconfig.GlobalConfig) {
	if u == nil {
		return
	}
	cacheKey1 := fmt.Sprintf("GetBalance_%s", u.Address)
	// cacheKeyBCA := fmt.Sprintf("bca_%v", u.Address)

	cacheKeyUsername := fmt.Sprintf("userObj %v", u.Username)
	cacheKeyEmail := fmt.Sprintf("userObj %v", u.Email)
	cacheKeySigner := fmt.Sprintf("userObj %v", u.PrimarySigner)
	cacheKeyUserID := fmt.Sprintf("userObj %v", u.ID)
	cacheKeyPShared := fmt.Sprintf("FetchWalletsPermissionsSharedWithUser_%s", u.ID)
	cacheKeyCuratedAssets := fmt.Sprintf("curatedAssets %v", u.Username)
	gc.RedisCache.DeleteFromCache(cacheKeyPShared, cacheKeyUsername, cacheKeyEmail, cacheKeySigner, cacheKeyUserID, cacheKey1, cacheKeyCuratedAssets)
	u.InvalidateUserWalletCache(gc)
	gc.RedisCache.InvalidateCachedHttpResponse(cacheKeyCuratedAssets)
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
		cacheKey3 := fmt.Sprintf("userObj %v", w.Alias)
		cacheKey4 := fmt.Sprintf("userObj %v", w.ID)
		cacheKeySigner := fmt.Sprintf("userObj %v", w.Signer)
		cacheKeyUserID := fmt.Sprintf("userObj %v", w.UserID)
		cacheKeyWalletAlias := fmt.Sprintf("walletObj_%v", w.Alias)
		cacheKeyWalletID := fmt.Sprintf("walletObj_%v", w.ID)
		cacheKeyPShared := fmt.Sprintf("FetchWalletsPermissionsSharedWithUser_%s", w.UserID)
		cacheKeybca1 := fmt.Sprintf("bca_%v", w.ID)
		if w.TempAddress != nil {
			cacheKeytempW := fmt.Sprintf("GetBalance_%s", *w.TempAddress)
			cacheKeybca2 := fmt.Sprintf("bca_%v", *w.TempAddress)
			gc.RedisCache.DeleteFromCache(cacheKeybca2, cacheKeytempW)

		}

		gc.RedisCache.DeleteFromCache(cacheKeyPShared, cacheKeyWalletAlias, cacheKeyWalletID, cacheKey1, cacheKey3, cacheKey4, cacheKeySigner, cacheKeyUserID, cacheKeybca1)

	}

}

func (w UserWallet) GetMarketOfferByID(offerID string, db *gorm.DB, gc *sharedconfig.GlobalConfig) (marketOfer MarketOffer, err error) {
	e := db.Where("id = ?", offerID).Where("source_wallet_address = ?", w.ID).First(&marketOfer).Error
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

// MarketOfferDetail is the Base equivalent of Stellar's horizon.Offer.
type MarketOfferDetail struct {
	ID     string
	Amount string
	Price  string
}

// GetBlockchainOfferDetail/CancelBlockchainOffer covered Stellar's native
// DEX offers (create/cancel a standing sell order on-ledger); Base has no
// native order book to hold such an offer (see
// internal/sharedconfig/order_book_summary.go's doc) - market-making on
// Base needs a real design (an on-chain limit-order contract, or an
// off-chain maker service) tracked as a follow-up, not attempted here.
// Both are long-tail (market_making.go only) and stubbed accordingly.
func (mo *MarketOffer) GetBlockchainOfferDetail(gc *sharedconfig.GlobalConfig) (offer MarketOfferDetail, err error) {
	if mo.BlockchainOfferID == nil {
		err = &tErrors.ErrorMissingParameter{Parameter: "blockchainOfferId"}
		return
	}
	return offer, &tErrors.ErrorTemporaryServerError{}
}

func (mo *MarketOffer) CancelBlockchainOffer(gc *sharedconfig.GlobalConfig) (offer MarketOfferDetail, err error) {
	if mo.BlockchainOfferID == nil {
		err = &tErrors.ErrorMissingParameter{Parameter: "blockchainOfferId"}
		return
	}
	return offer, &tErrors.ErrorTemporaryServerError{}
}

func (w UserWallet) GetCryptoDepositAddresses(currency string, gc *sharedconfig.GlobalConfig) (cryptoAddresses []CryptoWalletDepositAddress) {
	cryptoAddresses = make([]CryptoWalletDepositAddress, 0)
	e := gc.DB.Where("trovo_wallet_address = ? AND LOWER(currency) = ?", w.ID, strings.ToLower(currency)).Find(&cryptoAddresses).Error
	if e != nil {
		log.Printf("[GetCryptoDepositAddresses] error fetching cryptoAddresses from db %v", e)
	}

	return
}

func (w UserWallet) GetCryptoSubwallet(currency string, gc *sharedconfig.GlobalConfig) (subwallet CryptoSubWallet, err error) {

	var wdlResp CryptoSubwalletResponse

	client := http.DefaultClient
	url := fmt.Sprintf("%s/%s?currency=%s&uid=%s", os.Getenv("ONELIQUIDITY_BASE_URL"), "wallets/v1/sub", currency, w.Alias+"@"+os.Getenv("WALLET_DOMAIN"))

	request, err := http.NewRequest(http.MethodGet, url, nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("ONELIQUIDITY_TOKEN")))
	resp, err := client.Do(request)
	if err != nil {
		log.Println("[CreateCryptoSubwalletRequest] error sending request:", err)
		return
	}
	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		log.Println("[CreateCryptoSubwalletRequest] error response with code: ", resp.StatusCode, resp.Status)
		err = &tErrors.ErrorTemporaryServerError{}
		return
	}

	defer resp.Body.Close()
	//Decode the data
	if err = json.NewDecoder(resp.Body).Decode(&wdlResp); err != nil {
		log.Println("[CreateCryptoSubwalletRequest] error decoding response:", err)
		return
	}

	return wdlResp.Data, nil

}

func (w UserWallet) CreateCryptoSubwalletRequest(currency string, gc *sharedconfig.GlobalConfig) (subwallet CryptoSubWallet, err error) {

	var wdlResp CryptoSubwalletResponse

	client := http.DefaultClient
	url := fmt.Sprintf("%s/%s", os.Getenv("ONELIQUIDITY_BASE_URL"), "wallets/v1/sub")
	jbody, err := json.Marshal(OnliquiditySubWalletInput{
		Currency: currency,
		UID:      w.Alias + "@" + os.Getenv("WALLET_DOMAIN"),
	})
	if err != nil {
		log.Println("[CreateCryptoSubwalletRequest] error sending request:", err)

		return
	}
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(jbody))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("ONELIQUIDITY_TOKEN")))
	resp, err := client.Do(request)
	if err != nil {
		log.Println("[CreateCryptoSubwalletRequest] error sending request:", err)
		return
	}
	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		log.Println("[CreateCryptoSubwalletRequest] error response with code: ", resp.StatusCode, resp.Status)
		err = &tErrors.ErrorTemporaryServerError{}
		return
	}

	defer resp.Body.Close()
	//Decode the data
	if err = json.NewDecoder(resp.Body).Decode(&wdlResp); err != nil {
		log.Println("[CreateCryptoSubwalletRequest] error decoding response:", err)
		return
	}

	return wdlResp.Data, nil

}
