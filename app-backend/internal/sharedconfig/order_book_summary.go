package sharedconfig

import (
	"log"
	"os"
	"strings"
	tErrors "trovo-wallet-api/internal/errors"

	"github.com/ecnepsnai/discord"
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
// horizon.OrderBookSummary. Stellar's DEX order book is a native ledger
// feature with no Base/EVM equivalent - Base price discovery happens
// on-chain via AMM pools (e.g. a Uniswap-v3-style quoter), not a native
// order book. GetBantuOrderBookSummary below is a stub pending that
// integration (tracked, out of scope for this alteration pass); it
// always returns an empty book so callers' existing "no
// asks/bids -> zero/temporary-error" handling degrades gracefully rather
// than crashing.
type OrderBookSummary struct {
	Bids []PriceLevel `json:"bids"`
	Asks []PriceLevel `json:"asks"`
}

// GetBantuOrderBookSummary is a stub: Base has no native on-chain order
// book to query (see OrderBookSummary doc). Wire this to a real Base DEX
// price source (e.g. a Uniswap-v3 Quoter contract call for the
// asset/native pair) when swap pricing needs it to be more than
// unavailable.
func GetBantuOrderBookSummary(input OrderBookRequestInput) (orderBookSummary OrderBookSummary, err error) {
	return orderBookSummary, &tErrors.ErrorTemporaryServerError{}
}

// GetDollarAskPrice dollar ask price using USDB
func (gc *GlobalConfig) GetDollarAskPrice(sellingAssetCode, sellingContractAddress string) (usdPrice string, err error) {
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

	return usdPrice, nil
}

// GetAvalableMarketQuantity
func (gc *GlobalConfig) GetAvalableMarketQuantity(sellingAssetCode, sellingContractAddress, buyingAssetCode, buyingContractAddress string) (sellingQuantity, buyingQuantity string, err error) {
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

	orderBook, err := GetBantuOrderBookSummary(input)
	if err != nil {
		log.Printf("[GetAvalableMarketQuantity] Error fetching %v/%v market: %v\n", errBuyingAssetCode, errAssetCode, err)
		return "0", "0", &tErrors.ErrorTemporaryServerError{}
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

func logDiscordFailedRequest(msg string) {
	discord.WebhookURL = "https://discord.com/api/webhooks/827986576415129663/wqMKp9wxB_fxs9Q3zlMKCNPGENXmD_ueUnL8hVCu1wmRfD2wkXAjfP85k1Ro_2_wGfiY"
	if len(os.Getenv("EXPANSION_NETWORK_ERROR_WEBHOOK")) > 50 {
		discord.WebhookURL = os.Getenv("EXPANSION_NETWORK_ERROR_WEBHOOK")
	}
	discord.Say(msg)
}
