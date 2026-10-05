package users

import (
	"log"
	"os"
	"strings"
	bantupayerrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/offerbook"

	"github.com/shopspring/decimal"
)

// OrderBookRequestInput holds orderbook request input bindings
type OrderBookRequestInput struct {
	SellingAssetType       string `json:"selling_asset_type" form:"selling_asset_type"`
	SellingAssetCode       string `json:"selling_asset_code" form:"selling_asset_code"`
	SellingContractAddress string `json:"selling_contract_address" form:"selling_contract_address"`
	BuyingAssetType        string `json:"buying_asset_type" form:"buying_asset_type"`
	BuyingAssetCode        string `json:"buying_asset_code" form:"buying_asset_code"`
	BuyingContractAddress  string `json:"buying_contract_address" form:"buying_contract_address"`
	Limit                  string `json:"limit" form:"limit"`
}

// PriceLevel is the Base equivalent of Stellar's horizon.PriceLevel.
type PriceLevel struct {
	Price  string `json:"price"`
	Amount string `json:"amount"`
}

// OrderBookSummary is the Base equivalent of Stellar's
// horizon.OrderBookSummary - see getBantuOrderBookSummary's doc.
type OrderBookSummary struct {
	Bids []PriceLevel `json:"bids"`
	Asks []PriceLevel `json:"asks"`
}

// getBantuOrderBookSummary reads the order book of the selling asset
// against the buying asset from the offer book index (internal/offerbook).
// Only tokens trade there: a native-asset side has no book (an error).
func getBantuOrderBookSummary(input OrderBookRequestInput) (orderBookSummary OrderBookSummary, err error) {
	asks, bids, err := offerbook.CurrentTextBook(input.SellingContractAddress, input.BuyingContractAddress, input.Limit)
	if err != nil {
		return orderBookSummary, &bantupayerrors.ErrorTemporaryServerError{}
	}
	for _, l := range asks {
		orderBookSummary.Asks = append(orderBookSummary.Asks, PriceLevel{Price: l.Price, Amount: l.Amount})
	}
	for _, l := range bids {
		orderBookSummary.Bids = append(orderBookSummary.Bids, PriceLevel{Price: l.Price, Amount: l.Amount})
	}
	return orderBookSummary, nil
}

// GetDollarAskPrice dollar ask price using USDB
func GetDollarAskPrice(sellingAssetCode, sellingContractAddress string) (usdPrice string, err error) {
	var input OrderBookRequestInput
	var errAssetCode string
	if sellingAssetCode == "" {
		errAssetCode = "native"
	} else {
		errAssetCode = sellingAssetCode
	}
	input.SellingAssetCode = sellingAssetCode
	input.SellingContractAddress = sellingContractAddress
	if os.Getenv("DOLLAR_ASSET") != "" {
		asset := strings.Split(os.Getenv("DOLLAR_ASSET"), ":")
		input.BuyingAssetCode = asset[0]
		input.BuyingContractAddress = asset[1]
	} else {
		input.BuyingAssetCode = "USDB"
		input.BuyingContractAddress = os.Getenv("USDB_B20_TOKEN_ADDRESS")
	}
	orderBook, err := getBantuOrderBookSummary(input)
	if err != nil {
		log.Println("[GetDollarAskPrice] Error fetching dollar price:", err)
		return "0", &bantupayerrors.ErrorTemporaryServerError{}
	}
	if len(orderBook.Asks) == 0 {
		log.Printf("[Error GetDollarAskPrice]: error fetching dollar ASK price for asset %v, err: %v\n", errAssetCode, err)
		return "0", &bantupayerrors.ErrorTemporaryServerError{}
	}
	usdPrice = orderBook.Asks[0].Price

	return usdPrice, nil
}

// GetAvalableMarketQuantity
func GetAvalableMarketQuantity(sellingAssetCode, sellingContractAddress, buyingAssetCode, buyingContractAddress string) (sellingQuantity, buyingQuantity string, err error) {
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
	input.SellingContractAddress = sellingContractAddress

	input.BuyingAssetCode = buyingAssetCode
	input.BuyingContractAddress = buyingContractAddress

	orderBook, err := getBantuOrderBookSummary(input)
	if err != nil {
		log.Printf("[GetAvalableMarketQuantity] Error fetching %v/%v market: %v\n", errBuyingAssetCode, errAssetCode, err)
		return "0", "0", &bantupayerrors.ErrorTemporaryServerError{}
	}
	if len(orderBook.Asks) == 0 {
		log.Printf("[GetAvalableMarketQuantity] Error fetching %v asks for %v: %v\n", errBuyingAssetCode, errAssetCode, err)
		sellingQuantity = "0"
	} else {
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
		totalBids := decimal.Zero

		for _, v := range orderBook.Bids {
			totalBids = totalBids.Add(decimal.RequireFromString(v.Amount))
		}
		buyingQuantity = totalBids.String()
	}

	return sellingQuantity, buyingQuantity, nil
}
