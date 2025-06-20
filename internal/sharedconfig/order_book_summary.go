package sharedconfig

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/network"

	"github.com/ecnepsnai/discord"
	"github.com/shopspring/decimal"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/protocols/horizon"
)

// OrderBookRequestInput holds orderbook request input bindings
type OrderBookRequestInput struct {
	SellingAssetType   string `json:"selling_asset_type" form:"selling_asset_type"`
	SellingAssetCode   string `json:"selling_asset_code" form:"selling_asset_code"`
	SellingAssetIssuer string `json:"selling_asset_issuer" form:"selling_asset_issuer"`
	BuyingAssetType    string `json:"buying_asset_type" form:"buying_asset_type"`
	BuyingAssetCode    string `json:"buying_asset_code" form:"buying_asset_code"`
	BuyingAssetIssuer  string `json:"buying_asset_issuer" form:"buying_asset_issuer"`
	Limit              string `json:"limit" form:"limit"`
}

// getBantuOrderBookSummary gets orderbook on bantu network
func GetBantuOrderBookSummary(input OrderBookRequestInput) (orderBookSummary horizon.OrderBookSummary, err error) {
	discord.WebhookURL = "https://discord.com/api/webhooks/824381163367170058/75RxS1LzWA800hWereJJumw"
	if len(os.Getenv("EXPANSION_NETWORK_ERROR_WEBHOOK")) > 50 {
		discord.WebhookURL = os.Getenv("EXPANSION_NETWORK_ERROR_WEBHOOK")
	}
	client := network.GetBlockchainClient()
	var limit uint
	var sellingAssetType, buyingAssetType horizonclient.AssetType
	if (len(input.SellingAssetCode) == 0 && len(input.SellingAssetIssuer) == 0) || (input.SellingAssetCode == "native") {
		sellingAssetType = horizonclient.AssetTypeNative
		input.SellingAssetIssuer = ""
		input.SellingAssetCode = ""
	} else if input.SellingAssetType == "credit_alphanum4" || len(input.SellingAssetCode) < 5 {
		sellingAssetType = horizonclient.AssetType4
	} else if input.SellingAssetType == "credit_alphanum12" || len(input.SellingAssetCode) > 4 {
		sellingAssetType = horizonclient.AssetType12
	}

	if (len(input.BuyingAssetCode) == 0 && len(input.BuyingAssetIssuer) == 0) || (input.BuyingAssetCode == "native") {
		buyingAssetType = horizonclient.AssetTypeNative
		input.BuyingAssetIssuer = ""
		input.BuyingAssetCode = ""
	} else if input.BuyingAssetType == "credit_alphanum4" || len(input.BuyingAssetType) < 5 {
		buyingAssetType = horizonclient.AssetType4
	} else if input.BuyingAssetType == "credit_alphanum12" || len(input.BuyingAssetType) > 4 {
		buyingAssetType = horizonclient.AssetType12
	}

	if input.Limit == "" {
		limit = 50
	} else {
		plimit, _ := strconv.ParseInt(input.Limit, 10, 64)
		limit = uint(plimit)
	}
	oRequest := horizonclient.OrderBookRequest{
		SellingAssetCode:   input.SellingAssetCode,
		SellingAssetIssuer: input.SellingAssetIssuer,
		SellingAssetType:   sellingAssetType,
		BuyingAssetCode:    input.BuyingAssetCode,
		BuyingAssetIssuer:  input.BuyingAssetIssuer,
		BuyingAssetType:    buyingAssetType,
		Limit:              limit,
	}
	// fmt.Printf("Offer Request: %+v\n", oRequest)
	oSummary, err := client.OrderBook(oRequest)
	if err != nil {
		if strings.Contains(err.Error(), "decoding horizon.Problem") || strings.Contains(err.Error(), "tls") || strings.Contains(err.Error(), "timeout") || strings.Contains(err.Error(), "handshake") || strings.Contains(err.Error(), "read tcp") || strings.Contains(err.Error(), "connection reset by peer") || strings.Contains(err.Error(), "dial tcp") || strings.Contains(err.Error(), "no such host") {
			log.Println("[client.OrderBookRequest]", err)
			logDiscordFailedRequest(fmt.Sprintf("[client.OrderBookRequest]%v", err))
			return orderBookSummary, &tErrors.ErrorTemporaryServerError{}
		}
		// if hError, ok := err.(*horizonclient.Error); ok {} else {}
		if hError, ok := err.(*horizonclient.Error); ok {
			// Assertion succeeded
			//something went wrong, verify stage and check approprate action
			rCode, _ := hError.ResultCodes()
			rS, _ := hError.ResultString()
			log.Println("\n[client.OrderBookRequest] Problem in Request:", hError.Problem)
			log.Println("\n[client.OrderBookRequest] Result Codes in Request:", rCode)
			log.Println("\n[client.OrderBookRequest] Result String in Request:", rS)
			log.Printf("\n[client.OrderBookRequest] Problem in Request - RESPONSE: %+v\n", hError.Response)
			log.Println("[client.OrderBookRequest] Error submitting:", err)
			return orderBookSummary, &tErrors.ErrorTemporaryServerError{}
		} else {
			log.Println("[client.OrderBookRequest]error:", err)
			// Assertion failed
			return orderBookSummary, &tErrors.ErrorTemporaryServerError{}
		}

	}

	return oSummary, nil

}

// GetDollarAskPrice dollar ask price using USDB
func (gc *GlobalConfig) GetDollarAskPrice(sellingAssetCode, sellingAssetIssuer string) (usdPrice string, err error) {
	var input OrderBookRequestInput
	var errAssetCode string
	if sellingAssetCode == "" {
		errAssetCode = "native"
	} else {
		errAssetCode = sellingAssetCode
	}
	input.SellingAssetCode = sellingAssetCode
	input.SellingAssetIssuer = sellingAssetIssuer
	if os.Getenv("DOLLAR_ASSET") != "" {
		//USDB:GBTNUZDIUMWZEGTNQCL5F73PIABCBJ4YQA2VJS7HXTBRDDSTWCE6UNXE
		asset := strings.Split(os.Getenv("DOLLAR_ASSET"), ":")
		input.BuyingAssetCode = asset[0]
		input.BuyingAssetIssuer = asset[1]
	} else {
		input.BuyingAssetCode = "USDB"
		input.BuyingAssetIssuer = "GBTNUZDIUMWZEGTNQCL5F73PIABCBJ4YQA2VJS7HXTBRDDSTWCE6UNXE"
	}
	orderBook, err := GetBantuOrderBookSummary(input)
	if err != nil {
		log.Println("[GetDollarAskPrice] Error fetching dollar price:", err)
		return "0", &tErrors.ErrorTemporaryServerError{}
	}
	if len(orderBook.Asks) == 0 {
		log.Printf("[Error GetDollarAskPrice]: error fetching dollar ASK price for asset %v, err: %v\n", errAssetCode, err)
		return "0", &tErrors.ErrorTemporaryServerError{}
	}
	usdPrice = orderBook.Asks[0].Price
	// else if len(orderBook.Bids) > 0 {
	// 	price = orderBook.Bids[0].Price
	// }

	// fmt.Printf("OrderBookSummary: %+v\n", orderBook)
	return usdPrice, nil
}

// GetAvalableMarketQuantity
func (gc *GlobalConfig) GetAvalableMarketQuantity(sellingAssetCode, sellingAssetIssuer, buyingAssetCode, buyingAssetIssuer string) (sellingQuantity, buyingQuantity string, err error) {
	var input OrderBookRequestInput
	sellingQuantity = "0"
	buyingQuantity = "0"
	var errAssetCode, errBuyingAssetCode string
	if sellingAssetCode == "" {
		errAssetCode = "native"
	} else {
		errAssetCode = sellingAssetCode
	}
	if buyingAssetCode == "" {
		errBuyingAssetCode = "native"
	} else {
		errBuyingAssetCode = sellingAssetCode
	}
	input.SellingAssetCode = sellingAssetCode
	input.SellingAssetIssuer = sellingAssetIssuer

	input.BuyingAssetCode = buyingAssetCode
	input.BuyingAssetIssuer = buyingAssetIssuer

	orderBook, err := GetBantuOrderBookSummary(input)
	if err != nil {
		log.Printf("[GetAvalableMarketQuantity] Error fetching %v/%v market: %v\n", errBuyingAssetCode, errAssetCode, err)
		return "0", "0", &tErrors.ErrorTemporaryServerError{}
	}
	if len(orderBook.Asks) == 0 {
		log.Printf("[GetAvalableMarketQuantity] Error fetching %v asks for %v: %v\n", errBuyingAssetCode, errAssetCode, err)
		sellingQuantity = "0"
	} else {
		//asks exists
		totalAsks := decimal.Zero

		for _, v := range orderBook.Asks {
			totalAsks = totalAsks.Add(decimal.RequireFromString(v.Amount))
		}
		sellingQuantity = totalAsks.String()
	}
	if len(orderBook.Bids) == 0 {
		log.Printf("[GetAvalableMarketQuantity] Error fetching %v bids for %v: %v\n", errBuyingAssetCode, errAssetCode, err)
		buyingQuantity = "0"
	} else {
		//asks exists
		totalBids := decimal.Zero

		for _, v := range orderBook.Bids {
			totalBids = totalBids.Add(decimal.RequireFromString(v.Amount))
		}
		buyingQuantity = totalBids.String()
	}

	return sellingQuantity, buyingQuantity, nil
}

func logDiscordFailedRequest(msg string) {
	discord.WebhookURL = "https://discord.com/api/webhooks/827986576415129663/wqMKp9wxB_fxs9Q3zlMKCNPGENXmD_ueUnL8hVCu1wmRfD2wkXAjfP85k1Ro_2_wGfiY"
	if len(os.Getenv("EXPANSION_NETWORK_ERROR_WEBHOOK")) > 50 {
		discord.WebhookURL = os.Getenv("EXPANSION_NETWORK_ERROR_WEBHOOK")
	}
	discord.Say(msg)
}
