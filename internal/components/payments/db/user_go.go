package payments

import (
	"context"
	"encoding/base64"
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
	conDB "trovo-wallet-api/internal/db"
	dl "trovo-wallet-api/internal/dynamiclinks"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/network"
	pns "trovo-wallet-api/internal/pns"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/mailgun/mailgun-go/v4"
	"github.com/shopspring/decimal"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/protocols/horizon"
	"github.com/stellar/go/txnbuild"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type User struct {
	CreatedAt                time.Time    `json:"createdAt"`
	UpdatedAt                time.Time    `json:"updatedAt"`
	LastUpdatedMobileOn      time.Time    `json:"lastUpdatedMobileOn"`
	LastRecoveredAccountOn   time.Time    `json:"lastRecoveredAccountOn"`
	ID                       string       `json:"id"`
	Username                 string       `gorm:"size:16; index:idx_user_unique_username, unique" json:"username"`
	Email                    string       `gorm:"size:45; index:idx_user_unique_email, unique" json:"email"`
	ImageThumbnailURL        *string      `json:"imageThumbnailURL"`
	FirstName                string       `gorm:"size:50" json:"firstName"`
	LastName                 *string      `gorm:"size:50" json:"lastName"`
	Mobile                   *string      `gorm:"size:16; index:idx_user_unique_phone, unique" json:"mobile"`
	PublicKey                string       `gorm:"size:56; index:idx_user_unique_public_key, unique" json:"publicKey"`
	PrimarySigner            string       `gorm:"size:56; index:idx_user_unique_primary_signer, unique" json:"primarySigner"`
	Referrer                 *string      `gorm:"size:16; index:idx_user_referrer" json:"referrer"`
	ReferralLink             *string      `json:"referralLink"`
	ReferralQrCode           *string      `json:"referralQrCode"`
	PushNotificationToken    *string      `json:"pushNotificationToken"`
	Corporate                int          `gorm:"type:integer;not null; default:0" json:"corporate"`
	MobileVerified           int          `gorm:"type:integer;not null; default:0" json:"mobileVerified"`
	MembershipType           int          `gorm:"type:integer;not null; default:0" json:"membershipType"`
	MembershipExpiry         *time.Time   `json:"membershipExpiry"`
	KYCVerified              int          `gorm:"type:integer;not null; default:0" json:"kycVerified"`
	AccountRecoveryEnabled   int          `gorm:"type:integer;not null; default:0" json:"accountRecoveryEnabled"`
	AccountRecoveryExpiresOn *time.Time   `gorm:"null" json:"accountRecoveryExpiresOn"`
	UserWallets              []UserWallet `json:"userWallets"`
	PublicIP                 string       `gorm:"size:45" json:"publicIP"`
	CountryCode              *string      `gorm:"size:2;null"`
	Latitude                 *float64     `gorm:"null"`
	Longitude                *float64     `gorm:"null"`
	City                     *string      `gorm:"null;size:100"`
	Region                   *string      `gorm:"null;size:100"`
	RegionName               *string      `gorm:"null;size:100"`
	TimeZone                 *string      `gorm:"null;size:100"`
	ISP                      *string      `gorm:"null;size:150"`
	HasSecurityQuestions     int          `gorm:"type:integer;not null;default:0" json:"hasSecurityQuestions"`
	Verified                 int          `gorm:"type:integer;not null;default:0" json:"verified"`
	Suspended                int          `gorm:"type:integer;not null;default:0" json:"suspended"`
	SuspensionReason         *string      `gorm:"null" json:"suspensionReason"`
}

type UserWallet struct {
	CreatedAt              time.Time              `json:"createdAt"`
	UpdatedAt              time.Time              `json:"updatedAt"`
	ID                     string                 `gorm:"size:56" json:"publicKey"`
	TempPublicKey          *string                `gorm:"size:56;index:idx_user_wallet_temp_key;null"`
	Tag                    *string                `gorm:"null;size:16" json:"tag"`
	Description            *string                `gorm:"null;size:100" json:"description"`
	Alias                  string                 `gorm:"size:27; index:idx_unique_alias, unique" json:"alias"` //primaryUsername_tag for sub wallets
	Signer                 string                 `gorm:"size:56; index:idx_user_wallet_signer" json:"signer"`  //if ID is same as signer, then it is a primary wallet
	UserID                 string                 `gorm:"type:integer;not null; default:0;index:idx_user_wallets_user_id" json:"userId"`
	SharedAccessEnabled    int                    `gorm:"type:integer;not null; default:0" json:"sharedAccessEnabled"`
	UserWalletSharedAccess UserWalletSharedAccess `json:"userWalletSharedAccess"`
	Tracked                int                    `gorm:"type:integer;not null;default:0" json:"-"`
	PrimaryWallet          int                    `gorm:"type:integer;not null;default:0" json:"primaryWallet"`
}

type UserWalletSharedAccess struct {
	CreatedAt         time.Time          `json:"createdAt"`
	UpdatedAt         time.Time          `json:"updatedAt"`
	ID                string             `gorm:"" json:"accessId"`
	UserWalletID      string             `gorm:"size:56; index:idx_manage_access_user_wallet_id,unique" json:"walletPublicKey"`
	NumberOfApprovers int                `gorm:"type:integer; default:0" json:"numberOfApprovers"`
	Permissions       []WalletPermission `json:"permissions"`
}
type WalletPermission struct {
	CreatedAt                time.Time `json:"-"`
	UpdatedAt                time.Time `json:"-"`
	ID                       string
	WalletPublicKey          string `gorm:"size:90;not null; index:access_level_permission,unique;index:idx_public_key_shared" json:"walletPublicKey"`
	TargetUsername           string `gorm:"size:16;not null; index:access_level_permission,unique;" json:"targetUsername"`
	Permission               string `gorm:"size:10;not null; index:access_level_permission,unique" json:"permission"`
	UserWalletSharedAccessID string `gorm:"not null;index:idx_wallet_access_wallet_access_id" json:"userWalletSharedAccessId"`
}

// ReservedName holds model struct for ReservedName table
type ReservedName struct {
	ID           uint64 `gorm:"primaryKey"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	ReservedName *string `gorm:"size:50;not null;index:unique_reserved_name, unique;index:idx_reserved_status"`
	Status       *uint64 `gorm:"default:0;index:idx_reserved_status"`
}

type Balance struct {
	AssetIssuer string          `json:"assetIssuer"`
	AssetCode   string          `json:"assetCode"`
	Amount      decimal.Decimal `json:"amount"`
	QRCode      string          `json:"qrCode"`
	ImageURL    string          `json:"imageUrl"`
	UsdPrice    string          `json:"usdPrice"`
	NativePrice string          `json:"nativePrice"`
}

type NFT struct {
	AssetIssuer    string `json:"assetIssuer"`
	AssetCode      string `json:"assetCode"`
	NFTName        string `json:"nftName"`
	NFTDescription string `json:"nftDescription"`
	NFTImageURI    string `json:"nftImageURI"`
}

// Signer model for user
type Signer struct {
	Weight  int    `json:"weight"`
	Key     string `json:"key"`
	Type    string `json:"type"`
	Sponsor string `json:"sponsor"`
}

// Signer model for user
type Thresholds struct {
	LowThreshold    string `json:"low_threshold"`
	MediumThreshold string `json:"medium_threshold"`
	HighThreshold   string `json:"high_threshold"`
}

// AssetBalances holds user balances
type AssetBalances struct {
	Claimed   []Balance `json:"claimed"`
	Unclaimed []Balance `json:"unclaimed"`
}

// MappedNFTBalance holds user NFTs per mapped public key
type MappedNFTBalance map[string][]NFT

// UserBalanceForMerchant holds user balances
type UserBalanceForMerchant struct {
	Balances []Balance `json:"balances"`
}

type MappedBalance map[string]Balance

// PendingAssetToClaim holds pensing assets to be claimed
type PendingAssetToClaim struct {
	AssetCode            string `json:"assetCode"`
	AssetIssuer          string `json:"assetIssuer"`
	Transaction          string `json:"transaction"`
	TransactionSignature string `json:"transactionSignature"`
	TransactionID        string `json:"transactionId"`
	NetworkPassPhrase    string `json:"networkPassPhrase"`
}
type Trustline struct {
	AssetCode            string `json:"assetCode"`
	AssetIssuer          string `json:"assetIssuer"`
	Transaction          string `json:"transaction"`
	TransactionSignature string `json:"transactionSignature"`
	TransactionID        string `json:"transactionId"`
	NetworkPassPhrase    string `json:"networkPassPhrase"`
}

type BantuAsset struct {
	AssetCode   string
	AssetIssuer string
}
type UserWalletSharedAccessID string
type ThirdPartyWalletAccess struct {
	Owner             string `json:"owner"`
	WalletPublicKey   string `json:"walletPublicKey"`
	Permission        string `json:"permission"`
	WalletAlias       string `json:"walletAlias"`
	WalletDescription string `json:"walletDescription"`
}

// UserWalletID is type for wallet/sub-wallet Public Key
type UserWalletID string

func (i BantuAsset) GetDataKey(key string, gc *sharedconfig.GlobalConfig) string {
	client := network.GetBlockchainClient()

	defaultAssetImageURL := os.Getenv("DEFAULT_ASSET_IMAGE_URL")
	if len(strings.TrimSpace(i.AssetIssuer)) != 56 {
		return defaultAssetImageURL
	}
	issuerExists, _, _, _, account, err := network.BlockchainAccountProperties(client, i.AssetIssuer, txnbuild.NativeAsset{})
	if !issuerExists || err != nil {
		return defaultAssetImageURL
	}
	d, ok := account.Data[key]
	if !ok {
		return defaultAssetImageURL
	}
	decData, err := base64.StdEncoding.DecodeString(d)

	if err != nil {
		return defaultAssetImageURL
	}

	return string(decData)

}
func (i BantuAsset) GetAssetImageFromIssuer(gc *sharedconfig.GlobalConfig) string {
	client := network.GetBlockchainClient()
	cacheKey := fmt.Sprintf("url_%v_%v", i.AssetIssuer, i.AssetCode)
	ok, response := gc.RedisCache.GetCachedResult(cacheKey)
	if ok {
		return response.(string)
	}
	defaultAssetImageURL := os.Getenv("DEFAULT_ASSET_IMAGE_URL")
	if len(i.AssetCode) == 0 && len(i.AssetIssuer) == 0 {
		return os.Getenv("XBN_ASSET_IMAGE_URL")
	}
	if len(strings.TrimSpace(i.AssetIssuer)) != 56 {
		return defaultAssetImageURL
	}
	issuerExists, _, _, _, account, err := network.BlockchainAccountProperties(client, i.AssetIssuer, txnbuild.NativeAsset{})
	if !issuerExists || err != nil {
		return defaultAssetImageURL
	}
	key := fmt.Sprintf("imageurl_%v", strings.ToLower(i.AssetCode))
	d, ok := account.Data[key]
	if !ok {
		return defaultAssetImageURL
	}
	decData, err := base64.StdEncoding.DecodeString(d)

	if err != nil {
		return defaultAssetImageURL
	}
	url := string(decData)
	gc.RedisCache.StoreResultToCache(cacheKey, url, 1000000)
	return url
}

// GetBlockchainAccountDataKey fetches the bantu account information using public key
func (i BantuAsset) GetBlockchainAccountDataKey(account *horizon.Account, keys ...string) (dataValues map[string]string) {
	dataValues = make(map[string]string)

	for _, key := range keys {
		d, ok := account.Data[key]
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
		ok, response := gc.RedisCache.GetCachedResult(cacheKey)

		if ok {
			log.Printf("GetBalance[%v], served from cache\n", cacheKey)
			b := response.(map[string]interface{})
			for k, v := range b {
				mi := v.(map[string]interface{})
				usdPrice := "0"
				nativePrice := "0"
				if len(mi["usdPrice"].(string)) > 0 {
					usdPrice = mi["usdPrice"].(string)
				}
				if len(mi["nativePrice"].(string)) > 0 {
					nativePrice = mi["nativePrice"].(string)
				}
				balances[k] = Balance{
					AssetIssuer: mi["assetIssuer"].(string),
					AssetCode:   mi["assetCode"].(string),
					Amount:      decimal.RequireFromString(mi["amount"].(string)),
					QRCode:      mi["qrCode"].(string),
					ImageURL:    mi["imageUrl"].(string),
					UsdPrice:    usdPrice,
					NativePrice: nativePrice,
				}

			}
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
		}

		if !temp && err.Error() == "error-blockchain-account-not-activated" {

			log.Printf("[GetBalance] get blockchain account detail error: %v\n", err)
			//save to cache
			gc.RedisCache.StoreResultToCache(cacheKey, balances, 60)
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
			}
			// log.Printf("[BALANCE] balance: %+v\n", bal)

			m.Lock()
			balances[bal.Code+":"+bal.Issuer] = balance
			m.Unlock()

		}(v)
		wg.Wait()
	}
	//save to cache
	gc.RedisCache.StoreResultToCache(cacheKey, balances, 0)
	return balances, nil
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
		ok, response := gc.RedisCache.GetCachedResult(cacheKey)

		if ok {
			log.Printf("GetNFTBalance[%v], served from cache\n", cacheKey)

			miSlices := response.([]interface{})
			for _, v1 := range miSlices {
				mi := v1.(NFT)
				nfts = append(nfts, mi)
			}

			return nfts, nil
		}

	}

	account, _, err := u.GetBlockchainAccountDetail(temp)
	if err != nil {

		if !temp && err.Error() == "error-blockchain-account-not-activated" {

			//save to cache
			gc.RedisCache.StoreResultToCache(cacheKey, nfts, 0)
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
	gc.RedisCache.StoreResultToCache(cacheKey, nfts, 0)
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

func (u *User) BuildNewSubWallet(subWalletPublicKey, walletTag, walletDescription string, gc *sharedconfig.GlobalConfig) (userWallet UserWallet, err error) {
	walletTag = strings.ToLower(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(walletTag, "_", ""), ".", ""), " ", ""), "%", ""))
	walletDescription = strings.TrimSpace(walletDescription)

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
				log.Println("[BuildNewSubWallet] subwallet does not exist...", errWallet)

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
		Signer:        u.PublicKey,
		UserID:        u.ID,
	}
	// u.UserWallets = append(u.UserWallets, userSubWallet)
	return userSubWallet, nil
}

func (id UserWalletSharedAccessID) String() string {
	return string(id)
}

func (id UserWalletID) String() string {
	return string(id)
}

func (id UserWalletSharedAccessID) GetPermissionAssignment(db *gorm.DB) (assignment UserWalletSharedAccess, err error) {
	e := db.Where("id = ?", string(id)).First(&assignment).Error
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

func (id UserWalletID) GetPermissionList(db *gorm.DB) (accessList []WalletPermission) {
	accessList = make([]WalletPermission, 0)
	db.Preload(clause.Associations).Where("wallet_public_key = ?", string(id)).Find(&accessList)

	return
}

func (u UserWallet) GetPermissionList(db *gorm.DB) (accessList []WalletPermission) {
	accessList = make([]WalletPermission, 0)
	db.Preload(clause.Associations).Where("wallet_public_key = ?", u.ID).Find(&accessList)

	return
}

func (u *UserWallet) PublicKeyHasViewOnlyAccess(gc *sharedconfig.GlobalConfig) (viewOnly bool) {
	viewOnly = true
	if u.ID == "" {
		return false
	}
	accessList := u.GetPermissionList(gc.DB)
	if len(accessList) == 0 {
		return false
	}

	for _, access := range accessList {
		if access.Permission != "VIEW-ONLY" {
			return false
		}
	}

	return
}

func (u *UserWallet) HasViewOnlyAccess(gc *sharedconfig.GlobalConfig) (viewOnly bool) {
	viewOnly = true
	if u.ID == "" {
		return false
	}
	accessList := u.GetPermissionList(gc.DB)
	if len(accessList) == 0 {
		return false
	}

	for _, access := range accessList {
		if access.Permission != "VIEW-ONLY" {
			return false
		}
	}

	return
}

func (id UserWalletID) GetWallet(db *gorm.DB) (wallet UserWallet, err error) {
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
	return
}

func (id UserWalletID) PublicKeyHasViewOnlyAccess(gc *sharedconfig.GlobalConfig) (viewOnly bool) {
	viewOnly = true
	if id == "" {
		return false
	}
	wallet, err := id.GetWallet(gc.DB)
	if err != nil {
		return false
	}

	if wallet.SharedAccessEnabled == 0 {
		return false
	}
	for _, access := range wallet.UserWalletSharedAccess.Permissions {
		if access.Permission != "VIEW-ONLY" {
			return false
		}
	}

	return
}

func (u *User) GetWalletByPublicKey(publicKey string, db *gorm.DB) (wallet UserWallet, err error) {
	e := db.Where("id = ?", publicKey).First(&wallet).Error
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

func (id UserWalletID) GetWalletOwner(db *gorm.DB) (walletOwner User, err error) {
	e := db.Where("id = (SELECT user_id FROM user_wallets WHERE id = ?)", string(id)).First(&walletOwner).Error
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
	walletPermissions := u.Fetch3rdPartyWalletPermissions(gc)
	if len(walletPermissions) == 0 {
		return false
	}
	for _, walletAccess := range walletPermissions {
		if walletAccess.WalletPublicKey == publicKey {
			return true
		}
	}

	return false
}

func (u *User) GetAllWallets(gc *sharedconfig.GlobalConfig) (wallets []UserWallet) {
	wallets = make([]UserWallet, 0)
	gc.DB.Where("user_id = ?", u.ID).Find(&wallets)
	return
}

// Fetch3rdPartyWallets fetches all 3rd party wallets that the user is assigned to manage
func (u *User) Fetch3rdPartyWalletPermissions(gc *sharedconfig.GlobalConfig) (thirdPartyWallets []ThirdPartyWalletAccess) {
	var walletPermissions []WalletPermission
	thirdPartyWallets = make([]ThirdPartyWalletAccess, 0)
	cacheKey := fmt.Sprintf("Fetch3rdPartyWallets_%s", u.ID)

	{

		// search cache for balance
		ok, response := gc.RedisCache.GetCachedResult(cacheKey)

		if ok {
			log.Printf("Fetch3rdPartyWallets [%v], served from cache\n", cacheKey)
			w3rp := response.([]interface{})
			for _, w3 := range w3rp {
				w3i := w3.(map[string]interface{})
				thirdPartyWallets = append(thirdPartyWallets, ThirdPartyWalletAccess{
					Owner:             w3i["owner"].(string),
					WalletPublicKey:   w3i["walletPublicKey"].(string),
					Permission:        w3i["permission"].(string),
					WalletAlias:       w3i["walletAlias"].(string),
					WalletDescription: w3i["walletDescription"].(string),
				})
			}

			return
		}

	}
	e := gc.DB.Where("target_username = ?", u.Username).Find(&walletPermissions).Error
	if e != nil {
		return
	}
	if len(walletPermissions) == 0 {
		log.Println("[Fetch3rdPartyWallets] no wallet permissions found")
		return
	}
	for _, assignedPermission := range walletPermissions {
		//Get the permission assignment
		thirdPartyWallet := ThirdPartyWalletAccess{
			Permission: assignedPermission.Permission,
		}

		//use it to fetch wallet details
		wallet, err := UserWalletID(assignedPermission.WalletPublicKey).GetWallet(gc.DB)
		if err == nil {
			thirdPartyWallet.WalletPublicKey = wallet.ID
			thirdPartyWallet.WalletAlias = wallet.Alias
			if wallet.Description != nil {
				thirdPartyWallet.WalletDescription = *wallet.Description
			}
		}
		//use it to fetch wallet owner details
		owner, err := UserWalletID(assignedPermission.WalletPublicKey).GetWalletOwner(gc.DB)
		if err == nil {
			thirdPartyWallet.Owner = owner.Username
		}

		thirdPartyWallets = append(thirdPartyWallets, thirdPartyWallet)

	}
	//save to cache
	gc.RedisCache.StoreResultToCache(cacheKey, thirdPartyWallets, 120)

	return
}

func GetUser(userInfo string, db *gorm.DB) (user User, err error) {
	conDB.PrintDBStats("GetUserInfo", db)

	//e returns execution errors
	var e error
	if len(userInfo) == 56 {
		//56 char public key is supplied

		subQuery := db.Table("user_wallets").Where("id = ?", userInfo).Or("temp_public_key = ?", &userInfo).Or("signer = ?", userInfo).Select("user_id")
		e = db.Preload(clause.Associations).Where("id IN (?)", subQuery).First(&user).Error
	} else if strings.Contains(userInfo, "_") {
		//alias format is supplied
		subQuery := db.Table("user_wallets").Where("alias = ?", strings.ToLower(userInfo)).Select("user_id")
		e = db.Preload(clause.Associations).Where("id = (?)", subQuery).First(&user).Error

	} else {
		//search by ID and phone number, username, email

		e = db.Preload(clause.Associations).Where("id = ?", userInfo).Or("username = ?", strings.ToLower(userInfo)).Or("mobile = ?", &userInfo).Or("email = ?", userInfo).First(&user).Error
	}

	if e != nil {
		if errors.Is(e, gorm.ErrRecordNotFound) {
			//no user was found
			err = &tErrors.ErrorUserDoesNotExist{Username: userInfo}
			return
		}
		log.Println("[GetUserInfo] error: ", e)
		err = &tErrors.ErrorTemporaryServerError{}
		return

	}

	// log.Printf("user for %v is %v\n", userInfo, user)
	return user, nil

}

// GetWallet gets user wallet data by alias or public key or temp public key
func GetWallet(identifier string, db *gorm.DB) (userWallet UserWallet, temp bool, err error) {
	conDB.PrintDBStats("GetUserInfo", db)

	//e returns execution errors
	var e error
	if len(identifier) == 56 {
		//56 char public key is supplied
		e = db.Preload(clause.Associations).Where("id = ?", identifier).Or("temp_public_key = ?", &identifier).First(&userWallet).Error
		if e == nil {
			if identifier == *userWallet.TempPublicKey {
				temp = true
			}
			return
		}
	} else {
		//username is supplied
		e = db.Preload(clause.Associations).Where("alias = ?", strings.ToLower(identifier)).First(&userWallet).Error

	}

	if e != nil {
		if errors.Is(e, gorm.ErrRecordNotFound) {
			//no user was found
			err = &tErrors.CustomError{Param: "publicKey", Err: "error-wallet-does-not-exist", ErrMessage: fmt.Sprintf("%v is not assigned to any wallet", identifier)}
			return
		}
		log.Println("[GetUserInfo] error: ", e)
		err = &tErrors.ErrorTemporaryServerError{}
		return

	}

	// log.Printf("user for %v is %v\n", userInfo, user)
	return userWallet, temp, nil

}

// UsernameIsReserved check is name is reserved. Status = 0 means not available (reserved). Status = 1 means available
func UsernameIsReserved(username string, db *gorm.DB) (reserved bool, err error) {

	username = strings.TrimSpace(username)
	var reservedName ReservedName
	if err := db.Where("reserved_name = ? AND status = 0", strings.ToLower(strings.ReplaceAll(username, " ", ""))).First(&reservedName).Error; err != nil {

		return false, nil
	}

	return true, &tErrors.ErrorUsernameIsReserved{}
}

// PublicKeyIAlreadyExists check if public key already exists
func PublicKeyAlreadyExists(publicKey string, db *gorm.DB) (exists bool, err error) {
	// discord.WebhookURL = "https://discord.com/api/webhooks/824381163367170058/OXSX51RHd9DyLFbFipjdW3yXmyYC8SWwqd6HiXl6UtDzu75RxS1LzWA800hWereJJumw"
	// if len(os.Getenv("IMPORT_ERROR_WEBHOOK")) > 50 {
	// 	discord.WebhookURL = os.Getenv("IMPORT_ERROR_WEBHOOK")
	// }
	publicKey = strings.TrimSpace(publicKey)
	var userWallet UserWallet
	if err := db.Where("id = ?", strings.ToUpper(strings.ReplaceAll(publicKey, " ", ""))).First(&userWallet).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return true, &tErrors.ErrorTemporaryServerError{}
		}
		return false, nil
	}
	// discord.Say(fmt.Sprintf("[PublicKeyIsBanned] publicKey: %v is banned\n", publicKey))

	return true, &tErrors.CustomError{Param: "publicKey", Err: "error-public-key-already-exists", ErrMessage: fmt.Sprintf("Bantu Address [%v] already exists with another active account", userWallet.ID)}

}
func (u *User) SendPushMessage(title, body, imageURI string, dataPayload map[string]string, gc *sharedconfig.GlobalConfig) {
	//Send push notification to user
	// log.Println(title, body)

	if u.PushNotificationToken == nil {
		return
	}

	pns.SendFirebaseMessage(*u.PushNotificationToken, title, body, imageURI, dataPayload, gc.PushNotificationClient, gc.PNSContext)

}

// GetUserFromPrimarySigner fetches the user linked to the primary signer
func GetUserFromPrimarySigner(publicKey string, db *gorm.DB) (user User, err error) {

	publicKey = strings.TrimSpace(publicKey)
	// var user usermodels.User
	if err := db.Where("primary_signer = ?", strings.ToUpper(strings.ReplaceAll(publicKey, " ", ""))).First(&user).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return user, &tErrors.ErrorTemporaryServerError{}
		}
		return user, &tErrors.CustomError{Param: "primarySigner",
			Err:        "error primary signer does not exist",
			ErrMessage: "Primary Signer does not exist",
			Code:       http.StatusNotFound,
		}
	}
	// discord.Say(fmt.Sprintf("[PublicKeyIsBanned] publicKey: %v is banned\n", publicKey))

	return user, nil

}

func (w *UserWallet) WalletCountApproverAccess(gc *sharedconfig.GlobalConfig) (accessCount uint) {
	if w.SharedAccessEnabled == 0 {
		return 0
	}
	for _, access := range w.UserWalletSharedAccess.Permissions {
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
	for _, access := range w.UserWalletSharedAccess.Permissions {
		if access.Permission == "INITIATOR" {
			accessCount++
		}
	}

	return
}

func (w *UserWallet) HasInitiatorPermissionToPublicKey(ownerSignerPublicKey string, gc *sharedconfig.GlobalConfig) (hasAccess bool) {
	user, err := GetUserFromPrimarySigner(ownerSignerPublicKey, gc.DB)

	if err != nil {
		return false
	}
	walletPermissions := user.Fetch3rdPartyWalletPermissions(gc)
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

	walletPermissions := signerOwner.Fetch3rdPartyWalletPermissions(gc)
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
