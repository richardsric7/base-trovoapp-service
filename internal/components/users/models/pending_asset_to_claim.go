package users

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"
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
	Commit               int      `json:"commit"`
	Messages             []string `json:"messages"`
}

type BantuAsset struct {
	AssetCode   string
	AssetIssuer string
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
