package users

import (
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

// GetDataKey looked up a value in a Stellar account's on-chain
// manage_data key/value store (used for asset metadata like image URLs).
// Base/EVM accounts have no equivalent on-chain key/value store - asset
// metadata lives in this app's own curated-asset DB rows instead (see
// GetAssetImage), so this always returns the default.
func (i BantuAsset) GetDataKey(key string, gc *sharedconfig.GlobalConfig) string {
	return os.Getenv("DEFAULT_ASSET_IMAGE_URL")
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
	if url == nil {
		url = &defaultAssetImageURL
	}

	gc.RedisCache.StoreResultToCache(cacheKey, *url, 0)
	return *url
}

func (i BantuAsset) IsEnabled(gc *sharedconfig.GlobalConfig) bool {
	cacheKey := fmt.Sprintf("isenabled%v_%v", i.AssetCode, i.AssetIssuer)
	ok, response := gc.RedisCache.GetCachedResult(cacheKey)
	if ok {
		return response.(bool)
	}
	if len(i.AssetCode) == 0 && len(i.AssetIssuer) == 0 {
		return true
	}

	casset, err := assetsDB.GetCuratedAssetByCodeAndIssuer(i.AssetCode, i.AssetIssuer, false, gc)

	if err != nil {
		log.Printf("[GetAssetImage] <<<<<<< unable to get curated assets")
		return false
	}
	if casset.Inactive == 0 {
		log.Printf("[GetAssetImage] <<<<<<< curated asset is not enabled")
		return false
	}

	gc.RedisCache.StoreResultToCache(cacheKey, true, 120)
	return true
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

			// log.Printf("GetDepositAddresses[%v], served from cache\n", cacheKey)
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

func (i BantuAsset) IsTokenizedAsset(gc *sharedconfig.GlobalConfig) bool {

	code, issuer := i.AssetCode, i.AssetIssuer

	e := gc.DB.Where("Asset_Tokenization_Status > 4 AND Asset_Code = upper(?) AND Issuing_Wallet_Public_Key = upper(?)", code, issuer).First(&TokenizedAsset{}).Error

	return e == nil
}

// GetAssetImageFromIssuer looked up an asset's image URL from the
// issuer's on-chain manage_data store. See GetDataKey's doc - Base has no
// such store, so this always falls through to the DB-backed
// GetAssetImage/DEFAULT_ASSET_IMAGE_URL instead.
func (i BantuAsset) GetAssetImageFromIssuer(gc *sharedconfig.GlobalConfig) string {
	if len(i.AssetCode) == 0 && len(i.AssetIssuer) == 0 {
		return os.Getenv("NATIVE_ASSET_IMAGE_URL")
	}
	return os.Getenv("DEFAULT_ASSET_IMAGE_URL")
}

// BlockchainAssetFlags is the Base equivalent of Stellar's
// horizon.AccountFlags on an asset's issuer.
type BlockchainAssetFlags struct {
	AuthRequired  bool
	AuthRevocable bool
	AuthImmutable bool
}

// BlockchainAssetStat is the Base equivalent of Stellar's
// horizon.AssetStat.
type BlockchainAssetStat struct {
	Code   string
	Flags  BlockchainAssetFlags
	Amount string
}

// GetBlockchainAssetProperty reports whether this asset is registered in
// this app's own catalog (Code is set iff it is - the Base equivalent of
// Stellar's Horizon returning an empty AssetStat for an asset that
// doesn't exist) and whether it requires per-wallet authorization to
// hold/send (see internal/network's WalletAssetAuthorization doc).
// Unlike Stellar, where AuthRequired was an on-chain flag read from
// Horizon, a B20 asset has no such on-chain flag of its own - this app
// treats every tokenized (regulated) asset as requiring authorization,
// matching the same restriction the original enforced via Stellar
// trustlines for its issuer-authorized assets, and every other curated
// asset (e.g. stablecoins) as freely transferable.
func (i BantuAsset) GetBlockchainAssetProperty(gc *sharedconfig.GlobalConfig) (assetStat BlockchainAssetStat, err error) {
	var count int64
	gc.DB.Table("curated_assets").Where("asset_code = ? AND asset_issuer = ?", strings.ToUpper(i.AssetCode), strings.ToLower(i.AssetIssuer)).Count(&count)
	if count == 0 {
		gc.DB.Table("tokenized_assets").Where("asset_code = ? AND issuing_wallet_public_key = ?", strings.ToUpper(i.AssetCode), strings.ToLower(i.AssetIssuer)).Count(&count)
	}
	if count == 0 {
		return assetStat, nil
	}
	assetStat.Code = i.AssetCode
	assetStat.Flags.AuthRequired = gc.IsValidTokenizedAsset(i.AssetCode)
	assetStat.Amount = "0"
	return assetStat, nil
}

// GetBlockchainAccountDataKey looked up values from a Stellar account's
// on-chain manage_data store (e.g. a payment-callback URL an issuer
// published there). Base/EVM accounts have no equivalent store, so this
// always returns an empty map - see GetDataKey's doc for the same
// simplification applied elsewhere in this file.
func (i BantuAsset) GetBlockchainAccountDataKey(account *network.AccountInfo, keys ...string) (dataValues map[string]string) {
	dataValues = make(map[string]string)
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
