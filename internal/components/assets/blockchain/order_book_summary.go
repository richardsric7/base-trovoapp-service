package assets

import (
	"log"
	"os"
	"strconv"
	"strings"
	bantupayerrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/network"

	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/protocols/horizon"
	"gorm.io/gorm"
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
func getBantuOrderBookSummary(input OrderBookRequestInput) (orderBookSummary horizon.OrderBookSummary, err error) {
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
		if strings.Contains(err.Error(), "tls") || strings.Contains(err.Error(), "timeout") || strings.Contains(err.Error(), "handshake") || strings.Contains(err.Error(), "read tcp") || strings.Contains(err.Error(), "connection reset by peer") || strings.Contains(err.Error(), "dial tcp") || strings.Contains(err.Error(), "no such host") {
			log.Println("[client.OrderBookRequest]", err)
			return orderBookSummary, &bantupayerrors.ErrorTemporaryServerError{}
		}
		hError := err.(*horizonclient.Error)
		//something went wrong, verify stage and check approprate action
		rCode, _ := hError.ResultCodes()
		rS, _ := hError.ResultString()
		log.Println("\n[client.OrderBookRequest] Problem in Request:", hError.Problem)
		log.Println("\n[client.OrderBookRequest] Result Codes in Request:", rCode)
		log.Println("\n[client.OrderBookRequest] Result String in Request:", rS)
		log.Printf("\n[client.OrderBookRequest] Problem in Request - RESPONSE: %+v\n", hError.Response)
		log.Println("[client.OrderBookRequest] Error submitting:", err)
		return orderBookSummary, &bantupayerrors.ErrorTemporaryServerError{}

	}

	return oSummary, nil

}

// GetXBNDollarAskPrice dollar ask price using USDB
func GetXBNDollarAskPrice(db *gorm.DB) (usdPrice string, err error) {
	type XbnDollarPrice struct {
		ID            string `gorm:"primaryKey"`
		AskRate       string `gorm:"size:100"`
		BidRate       string `gorm:"size:100"`
		LastTradeRate string `gorm:"size:100"`
		Source        string `gorm:"size:100"`
	}
	var xbnPrice XbnDollarPrice
	err = db.Where("source = ?", "METRICS").First(&xbnPrice).Error
	if err != nil {
		return "0", err
	}

	return xbnPrice.LastTradeRate, nil
	// return t.Ticker.LastTradeRate, nil

}

// GetDollarAskPrice dollar ask price using USDB
func GetDollarAskPrice(sellingAssetCode, sellingAssetIssuer string) (usdPrice string, err error) {
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
	orderBook, err := getBantuOrderBookSummary(input)

	if len(orderBook.Asks) == 0 {
		log.Printf("[Error GetDollarAskPrice]: error fetching dollar ASK price for asset %v, err: %v\n", errAssetCode, err)
		return "0", &bantupayerrors.ErrorTemporaryServerError{}
	}
	usdPrice = orderBook.Asks[0].Price
	// else if len(orderBook.Bids) > 0 {
	// 	price = orderBook.Bids[0].Price
	// }

	// fmt.Printf("OrderBookSummary: %+v\n", orderBook)
	return usdPrice, nil
}

// GetNativeAskPrice native (XBN) ask price
func GetNativeAskPrice(sellingAssetCode, sellingAssetIssuer string) (nativePrice string, err error) {

	var input OrderBookRequestInput
	if sellingAssetCode == "" {

		return "0", &bantupayerrors.ErrorTemporaryServerError{}
	}
	input.SellingAssetCode = sellingAssetCode
	input.SellingAssetIssuer = sellingAssetIssuer

	orderBook, err := getBantuOrderBookSummary(input)

	if len(orderBook.Asks) == 0 || err != nil {
		return "0", &bantupayerrors.ErrorTemporaryServerError{}
	}
	nativePrice = orderBook.Asks[0].Price
	// else if len(orderBook.Bids) > 0 {
	// 	price = orderBook.Bids[0].Price
	// }

	// fmt.Printf("OrderBookSummary: %+v\n", orderBook)
	return nativePrice, nil
}
