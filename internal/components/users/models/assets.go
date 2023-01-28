package users

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	assetsDB "trovo-wallet-api/internal/components/assets/db"
	assetModels "trovo-wallet-api/internal/components/assets/models"
	"trovo-wallet-api/internal/network"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/stellar/go/protocols/horizon"
	"github.com/stellar/go/txnbuild"
)

// PendingAssetToClaim holds pensing assets to be claimed
type PendingAssetToClaim struct {
	AssetCode            string   `json:"assetCode"`
	AssetIssuer          string   `json:"assetIssuer"`
	Transaction          string   `json:"transaction"`
	TransactionSignature string   `json:"transactionSignature"`
	TransactionID        string   `json:"transactionId"`
	NetworkPassPhrase    string   `json:"networkPassPhrase"`
	Multiparty           int      `json:"-"`
	TransactionSource    string   `json:"-"`
	SignatureRequired    int      `json:"signatureRequired"`
	ReturnedDescription  string   `json:"-"`
	Commit               int      `json:"commit"`
	Messages             []string `json:"messages"`
}
type Trustline struct {
	AssetCode            string   `json:"assetCode"`
	AssetIssuer          string   `json:"assetIssuer"`
	Transaction          string   `json:"transaction"`
	TransactionSignature string   `json:"transactionSignature"`
	TransactionID        string   `json:"transactionId"`
	NetworkPassPhrase    string   `json:"networkPassPhrase"`
	Multiparty           int      `json:"-"`
	TransactionSource    string   `json:"-"`
	SignatureRequired    int      `json:"signatureRequired"`
	ReturnedDescription  string   `json:"-"`
	Commit               int      `json:"commit"`
	Messages             []string `json:"messages"`
}

// MintingInfo represents model for minting asset
type MintingInfo struct {
	Destination             string            `json:"destination"`
	Memo                    string            `json:"memo"`
	AssetIssuer             string            `json:"assetIssuer"`
	AssetCode               string            `json:"assetCode"`
	Amount                  string            `json:"amount"`
	Transaction             string            `json:"transaction"`
	TransactionSignature    string            `json:"transactionSignature"`
	TransactionID           string            `json:"transactionId"`
	NetworkPassPhrase       string            `json:"networkPassPhrase"`
	DestinationFirstName    string            `json:"destinationFirstName"`
	DestinationLastName     string            `json:"destinationLastName"`
	DestinationThumbnail    string            `json:"destinationThumbnail"`
	ReturnedDescription     string            `json:"-"`
	DestinationVerified     int               `json:"destinationVerified"`
	Multiparty              int               `json:"-"`
	TransactionSource       string            `json:"-"`
	SignatureRequired       int               `json:"signatureRequired"`
	Commit                  int               `json:"commit"`
	ChannelAccount          string            `json:"channelAccount"`
	ChannelAccountSignature string            `json:"channelAccountSignature"`
	Messages                []string          `json:"messages"`
	CallbackURLS            map[string]string `json:"-"`
}
type BantuAsset struct {
	AssetCode   string `json:"assetCode"`
	AssetIssuer string `json:"assetIssuer"`
}

type OrderBook struct {
	Bids     []string   `json:"bids"`
	Asks     []string   `json:"asks"`
	Asset    BantuAsset `json:"asset"`
	Currency BantuAsset `json:"currency"`
}

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

func (i BantuAsset) GetAssetImage(gc *sharedconfig.GlobalConfig) string {
	cacheKey := fmt.Sprintf("url%v_%v", i.AssetCode, i.AssetIssuer)
	ok, response := gc.RedisCache.GetCachedResult(cacheKey)
	if ok {
		return response.(string)
	}
	defaultAssetImageURL := os.Getenv("DEFAULT_ASSET_IMAGE_URL")
	if len(i.AssetCode) == 0 && len(i.AssetIssuer) == 0 {
		return os.Getenv("NATIVE_ASSET_IMAGE_URL")
	}
	if len(strings.TrimSpace(i.AssetIssuer)) != 56 {
		return defaultAssetImageURL
	}
	cassets := assetsDB.GetCuratedAssets(false, gc)

	if len(cassets) == 0 {
		log.Printf("[GetAssetImage] <<<<<<< unable to get curated assets. returning default asset image")
		return defaultAssetImageURL
	}
	v, ok := cassets[i.AssetCode+":"+i.AssetIssuer]
	if !ok {
		return defaultAssetImageURL

	}
	url := v.ImageURL

	gc.RedisCache.StoreResultToCache(cacheKey, url, 0)
	return url
}

func (i BantuAsset) CanDeposit(gc *sharedconfig.GlobalConfig) bool {

	cassets := assetsDB.GetCuratedAssets(false, gc)

	if len(cassets) == 0 {

		return false
	}
	v, ok := cassets[i.AssetCode+":"+i.AssetIssuer]
	if !ok {
		return false

	}
	if v.GenerateDepositAddress == 1 {
		return true
	}

	return false
}
func (i BantuAsset) GetDepositAddresses(walletID string, gc *sharedconfig.GlobalConfig) (depositAddresses []CryptoWalletDepositAddress) {
	depositAddresses = make([]CryptoWalletDepositAddress, 0)
	cacheKey := fmt.Sprintf("depositAddresses_%s_%s", i.AssetCode, walletID)
	{

		// search cache
		ok, rawdata := gc.RedisCache.GetCachedResultRaw(cacheKey)

		if ok {

			log.Printf("GetDepositAddresses[%v], served from cache\n", cacheKey)
			json.Unmarshal(rawdata, &depositAddresses)
			return
		}

	}
	if !i.CanDeposit(gc) {
		return
	}
	//get the deposit addresses
	e := gc.DB.Where("trovo_wallet_public_key = ? AND LOWER(currency) = ?", walletID, strings.ToLower(i.AssetCode)).Find(&depositAddresses).Error
	if e != nil {
		log.Printf("[GetDepositAddresses]Error getting deposit address, error: %v\n", e)
	}
	if len(depositAddresses) > 0 {
		gc.RedisCache.StoreResultToCacheRaw(cacheKey, depositAddresses, 800)
	}

	return
}

func (i BantuAsset) CanWithdraw(gc *sharedconfig.GlobalConfig) bool {

	cassets := assetsDB.GetCuratedAssets(false, gc)

	if len(cassets) == 0 {

		return false
	}
	v, ok := cassets[i.AssetCode+":"+i.AssetIssuer]
	if !ok {
		return false

	}
	if v.Withdrawable == 1 {
		return true
	}

	return false
}

func (i BantuAsset) GetAssetImageFromIssuer(gc *sharedconfig.GlobalConfig) string {
	client := network.GetBlockchainClient()
	cacheKey := fmt.Sprintf("url%v_%v", i.AssetCode, i.AssetIssuer)
	ok, response := gc.RedisCache.GetCachedResult(cacheKey)
	if ok {
		return response.(string)
	}
	defaultAssetImageURL := os.Getenv("DEFAULT_ASSET_IMAGE_URL")
	if len(i.AssetCode) == 0 && len(i.AssetIssuer) == 0 {
		return os.Getenv("NATIVE_ASSET_IMAGE_URL")
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

type Currency string

func (c Currency) GetCurratedAsset(gc *sharedconfig.GlobalConfig) (ca assetModels.CuratedAsset, err error) {
	currency := strings.ToUpper(string(c))
	cas := assetsDB.GetCuratedAssets(false, gc)

	for _, k := range cas {
		if strings.EqualFold(k.AssetCode, currency) {
			if k.Withdrawable == 0 {
				log.Println("[GetCurratedAsset] currency not marked as withdrawable.", currency)
				err = errors.New(currency + " not withdrawable")
				return
			}
			return k, nil

		}
	}
	log.Println("[GetCurratedAsset] currency not available in list of currencies.", currency)
	err = errors.New(currency + " is an invalid withdrawable currency")
	return

}
