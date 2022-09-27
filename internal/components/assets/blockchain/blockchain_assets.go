package assets

import (
	"log"
	"strings"
	models "trovo-wallet-api/internal/components/assets/models"
	bantupayerrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/network"

	"github.com/stellar/go/clients/horizonclient"
)

// GetBlockchainAsset gets blockchain assets
func GetBlockchainAsset(assetCode, assetIssuer, cursor, order string, limit uint) (paginatedBlockchainAssets models.PaginatedBlockchainAssets, err error) {
	var ordering horizonclient.Order
	var assets []models.BlockchainAsset
	paginationToken := ""
	if limit < 1 {
		limit = 25
	}
	if order == "asc" {
		ordering = horizonclient.OrderAsc
	} else {
		ordering = horizonclient.OrderDesc
	}
	client := network.GetBlockchainClient()
	assetRequest := horizonclient.AssetRequest{
		ForAssetCode:   assetCode,
		ForAssetIssuer: assetIssuer,
		Order:          ordering,
		Cursor:         cursor,
		Limit:          limit,
	}

	blockchainAssets, err := client.Assets(assetRequest)

	if err != nil {
		if strings.Contains(err.Error(), "tls") || strings.Contains(err.Error(), "timeout") || strings.Contains(err.Error(), "handshake") || strings.Contains(err.Error(), "read tcp") || strings.Contains(err.Error(), "connection reset by peer") || strings.Contains(err.Error(), "dial tcp") || strings.Contains(err.Error(), "no such host") {
			log.Print("[GetBlockchainAsset]", err)
			return paginatedBlockchainAssets, &bantupayerrors.ErrorTemporaryServerError{}
		}
		// log.Println("[GetAssets]: ", err)
		hError := err.(*horizonclient.Error)
		//something went wrong, verify stage and check approprate action
		rCode, _ := hError.ResultCodes()
		rS, _ := hError.ResultString()
		log.Println("[GetBlockchainAsset] Problem in Transaction:", hError.Problem)
		log.Println("[GetBlockchainAsset] Result Codes in Transaction:", rCode)
		log.Println("[GetBlockchainAsset] Result String in Transaction:", rS)
		log.Printf("[GetBlockchainAsset] Problem in Transaction - RESPONSE: %+v\n", hError.Response)
		return paginatedBlockchainAssets, &bantupayerrors.ErrorTemporaryServerError{}
	}
	for _, v := range blockchainAssets.Embedded.Records {
		paginationToken = v.PagingToken()
		var asset models.BlockchainAsset
		asset.AssetCode = v.Code
		asset.AssetIssuer = v.Issuer
		asset.AmountOfTokens = v.Amount
		asset.NumOfAccounts = int64(v.NumAccounts)
		asset.AuthImmutable = v.Flags.AuthImmutable
		asset.AuthRequired = v.Flags.AuthRequired
		asset.AuthRevocable = v.Flags.AuthRevocable
		asset.Toml = v.Links.Toml.Href
		assets = append(assets, asset)

	}
	paginatedBlockchainAssets.PageCursor = paginationToken
	paginatedBlockchainAssets.Assets = assets
	return paginatedBlockchainAssets, nil
}
