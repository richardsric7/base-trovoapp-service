package assets

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/ecnepsnai/discord"
	"gorm.io/gorm"
)

// PriceLevel is the Base equivalent of Stellar's horizon.PriceLevel.
type PriceLevel struct {
	Price  string `json:"price"`
	Amount string `json:"amount"`
}

// OrderBookSummary is the Base equivalent of Stellar's
// horizon.OrderBookSummary - see GetBantuOrderBookSummary's doc.
type OrderBookSummary struct {
	Bids    []PriceLevel `json:"bids"`
	Asks    []PriceLevel `json:"asks"`
	Selling Asset        `json:"selling"`
	Buying  Asset        `json:"buying"`
}

// TradeAggregationRecord is the Base equivalent of one horizon
// trade-aggregation bucket.
type TradeAggregationRecord struct {
	Timestamp     int64
	TradeCount    int64
	BaseVolume    string
	CounterVolume string
	Average       string
	High          string
	Low           string
	Open          string
	Close         string
}

// TradeAggregationsPage is the Base equivalent of Stellar's
// horizon.TradeAggregationsPage.
type TradeAggregationsPage struct {
	Embedded struct {
		Records []TradeAggregationRecord
	}
}

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
type Asset struct {
	AssetCode       string `json:"assetCode"`
	ContractAddress string `json:"contractAddress"`
}

type OfferVolume struct {
	Price    string `json:"price"`
	Quantity string `json:"quantity"`
}
type OrderBook struct {
	Bids     []OfferVolume `json:"bids"`
	Asks     []OfferVolume `json:"asks"`
	Asset    Asset         `json:"asset"`
	Currency Asset         `json:"currency"`
}
type ChartRecord struct {
	Timestamp      int64  `json:"timestamp"`
	TradeCount     int64  `json:"tradeCount"`
	AssetVolume    string `json:"assetVolume"`
	CurrencyVolume string `json:"currencyVolume"`
	Average        string `json:"average"`
	High           string `json:"high"`
	Low            string `json:"low"`
	Open           string `json:"open"`
	Close          string `json:"close"`
}
type TradeChart struct {
	Asset        Asset         `json:"asset"`
	Currency     Asset         `json:"currency"`
	ChartRecords []ChartRecord `json:"chartRecords"`
}

type TradeAggregateInput struct {
	StartTime              time.Time
	EndTime                time.Time
	Resolution             time.Duration
	Offset                 time.Duration
	BaseAssetCode          string
	BaseContractAddress    string
	BaseAssetType          string
	CounterAssetCode       string
	CounterContractAddress string
	CounterAssetType       string
	Order                  string
	Limit                  string
}
type PriceCache struct {
	Price     string `json:"price"`
	PriceType string `json:"priceType,omitempty"`
}

// getTradeAggregate is a stub: like GetBantuOrderBookSummary below,
// Stellar's trade-aggregation (candlestick chart) data is a native ledger
// feature with no Base/EVM equivalent - Base price/volume history comes
// from indexing on-chain swap events (a follow-up, out of scope for this
// alteration pass), not a network-native aggregation endpoint. Always
// returns an empty page so GetChartRecords' existing handling renders an
// empty chart rather than crashing.
func getTradeAggregate(input TradeAggregateInput) (tds TradeAggregationsPage, err error) {
	return tds, &tErrors.ErrorTemporaryServerError{}
}

// GetBantuOrderBookSummary is a stub: Base has no native on-chain order
// book to query (see OrderBookSummary's doc). Wire this to a real Base
// DEX price source (e.g. a Uniswap-v3 Quoter contract call for the
// asset/currency pair) when swap/price-display features need it to be
// more than unavailable - every caller in this file already falls back
// to a cached price or "0"/temporary-error on a non-nil err, so this
// degrades gracefully rather than crashing.
func GetBantuOrderBookSummary(input OrderBookRequestInput) (orderBookSummary OrderBookSummary, err error) {
	return orderBookSummary, &tErrors.ErrorTemporaryServerError{}
}

// GetGASDollarAskPrice dollar ask price using USDB
func GetGASDollarAskPrice(db *gorm.DB) (usdPrice string, err error) {
	type GasDollarPrice struct {
		ID            string `gorm:"primaryKey"`
		AskRate       string `gorm:"size:100"`
		BidRate       string `gorm:"size:100"`
		LastTradeRate string `gorm:"size:100"`
		Source        string `gorm:"size:100"`
	}
	var gasPrice GasDollarPrice
	err = db.Where("source = ?", "METRICS").First(&gasPrice).Error
	if err != nil {
		return "0", err
	}

	return gasPrice.LastTradeRate, nil
	// return t.Ticker.LastTradeRate, nil

}

// GetDollarPrice dollar ask price using USDB
func GetDollarPrice(sellingAssetCode, sellingContractAddress string, gc *sharedconfig.GlobalConfig, checkCacheFirst bool) (usdPrice, priceType string, err error) {
	var priceCache PriceCache
	var input OrderBookRequestInput
	priceType = "ask"
	usdPrice = "0"
	sellingAssetCode = strings.ToUpper(sellingAssetCode)
	//sell main asset, buying currency (dollar)
	var errAssetCode string
	if sellingAssetCode == "" {
		errAssetCode = "native"
	} else {
		errAssetCode = sellingAssetCode
	}
	cacheKey := fmt.Sprintf("%v.%v_dollar", sellingAssetCode, sellingContractAddress)
	if checkCacheFirst {
		ok, concatPriceByte := gc.RedisCache.GetCachedResultRaw(cacheKey)
		if ok {
			json.Unmarshal(concatPriceByte, &priceCache)
			// log.Printf("[GetDollarPrice] cache result: %+v\n", priceCache)

			return priceCache.Price, priceCache.PriceType, nil
		}
	}

	input.SellingAssetCode = sellingAssetCode
	input.SellingContractAddress = sellingContractAddress
	dollarAsset := strings.Split(os.Getenv("DOLLAR_ASSET"), ":")
	if len(dollarAsset) != 2 {
		return "0", priceType, &tErrors.ErrorTemporaryServerError{}
	}

	if strings.EqualFold(sellingAssetCode, dollarAsset[0]) && strings.EqualFold(sellingContractAddress, dollarAsset[1]) {
		//it is dollar asset
		return "1", priceType, nil
	}

	if strings.HasPrefix(sellingAssetCode, "USD") || strings.HasSuffix(sellingAssetCode, "USD") {
		return "1", priceType, nil
	}

	input.BuyingAssetCode = dollarAsset[0]
	input.BuyingContractAddress = dollarAsset[1]

	orderBook, err := GetBantuOrderBookSummary(input)
	if err != nil {
		log.Printf("[GetDollarPrice] error getting order book summary: %v\n", err)
		//fetch from last stored in cache
		ok, concatPriceByte := gc.RedisCache.GetCachedResultRaw(cacheKey)
		if ok {
			json.Unmarshal(concatPriceByte, &priceCache)
			// log.Printf("[GetDollarPrice] cache result: %+v\n", priceCache)

			return priceCache.Price, priceCache.PriceType, nil
		}

		return

	}

	if len(orderBook.Asks) == 0 {
		priceType = "bid"
		if len(orderBook.Bids) == 0 {

			log.Printf("[Error GetDollarAskPrice]: error fetching dollar ASK price for asset %v, err: %v\n", errAssetCode, err)
			return "0", priceType, &tErrors.ErrorTemporaryServerError{}
		}
		usdPrice = orderBook.Bids[0].Price
		priceCache.Price = usdPrice
		priceCache.PriceType = priceType
		gc.RedisCache.StoreResultToCacheRaw(cacheKey, priceCache, 120)
		return usdPrice, priceType, nil
	}
	usdPrice = orderBook.Asks[0].Price
	priceCache.Price = usdPrice
	priceCache.PriceType = priceType
	gc.RedisCache.StoreResultToCacheRaw(cacheKey, priceCache, 120)
	return usdPrice, priceType, nil
}

// GetNairaPrice dollar ask price using USDB
func GetNairaPrice(sellingAssetCode, sellingContractAddress string, gc *sharedconfig.GlobalConfig, checkCacheFirst, enabledAsset bool) (nairaPrice, priceType string, err error) {
	var priceCache PriceCache

	var input OrderBookRequestInput
	priceType = "ask"
	nairaPrice = "0"
	if !enabledAsset {
		return
	}
	sellingAssetCode = strings.ToUpper(sellingAssetCode)
	//sell main asset, buying currency (dollar)
	var errAssetCode string
	if sellingAssetCode == "" {
		errAssetCode = "native"
	} else {
		errAssetCode = sellingAssetCode
	}
	cacheKey := fmt.Sprintf("%v.%v_naira", sellingAssetCode, sellingContractAddress)
	if checkCacheFirst {
		ok, concatPriceByte := gc.RedisCache.GetCachedResultRaw(cacheKey)
		if ok {
			json.Unmarshal(concatPriceByte, &priceCache)
			// log.Printf("[GetDollarPrice] cache result: %+v\n", priceCache)

			return priceCache.Price, priceCache.PriceType, nil
		}
	}

	input.SellingAssetCode = sellingAssetCode
	input.SellingContractAddress = sellingContractAddress
	nairaAsset := strings.Split(os.Getenv("NAIRA_ASSET"), ":")
	if len(nairaAsset) != 2 {
		return "0", priceType, &tErrors.ErrorTemporaryServerError{}
	}

	if strings.EqualFold(sellingAssetCode, nairaAsset[0]) && strings.EqualFold(sellingContractAddress, nairaAsset[1]) {
		//it is dollar asset
		return "1", priceType, nil
	}

	if strings.HasPrefix(sellingAssetCode, "NGN") || strings.HasSuffix(sellingAssetCode, "NGN") {
		return "1", priceType, nil
	}

	input.BuyingAssetCode = nairaAsset[0]
	input.BuyingContractAddress = nairaAsset[1]

	orderBook, err := GetBantuOrderBookSummary(input)
	if err != nil {
		log.Printf("[GetNairaPrice] error getting order book summary: %v\n", err)
		//fetch from last stored in cache
		ok, concatPriceByte := gc.RedisCache.GetCachedResultRaw(cacheKey)
		if ok {
			json.Unmarshal(concatPriceByte, &priceCache)
			// log.Printf("[GetNairaPrice] cache result: %+v\n", priceCache)

			return priceCache.Price, priceCache.PriceType, nil
		}

		return

	}

	if len(orderBook.Asks) == 0 {
		priceType = "bid"
		if len(orderBook.Bids) == 0 {

			log.Printf("[Error GetNairaAskPrice]: error fetching naira ASK price for asset %v/%v, err: NO Markets\n", errAssetCode, input.BuyingAssetCode)
			return "0", priceType, &tErrors.ErrorTemporaryServerError{}
		}
		nairaPrice = orderBook.Bids[0].Price
		priceCache.Price = nairaPrice
		priceCache.PriceType = priceType
		gc.RedisCache.StoreResultToCacheRaw(cacheKey, priceCache, 120)
		return nairaPrice, priceType, nil
	}
	nairaPrice = orderBook.Asks[0].Price
	priceCache.Price = nairaPrice
	priceCache.PriceType = priceType
	gc.RedisCache.StoreResultToCacheRaw(cacheKey, priceCache, 120)
	return nairaPrice, priceType, nil
}

// GetNativeAskPrice native (GAS) ask price
func GetNativeAskPrice(sellingAssetCode, sellingContractAddress string, gc *sharedconfig.GlobalConfig, checkCacheFirst, isEnabled bool) (nativePrice string, err error) {
	var priceCache PriceCache
	var nativeCode, nativeIssuer string
	if !isEnabled {
		return "0", nil
	}
	nv := strings.Split(os.Getenv("USE_ASSET_FOR_NATIVE_PRICE"), ":")
	if len(nv) == 2 {
		nativeCode = nv[0]
		nativeIssuer = nv[1]
	}
	cacheKey := fmt.Sprintf("%v.%v_nativePrice", sellingAssetCode, sellingContractAddress)
	if checkCacheFirst {
		ok, concatPriceByte := gc.RedisCache.GetCachedResultRaw(cacheKey)
		if ok {
			json.Unmarshal(concatPriceByte, &priceCache)
			// log.Printf("[GetNativeAskPrice] cache result: %+v\n", priceCache)

			return priceCache.Price, nil
		}
	}
	var input OrderBookRequestInput

	input.SellingAssetCode = sellingAssetCode
	input.SellingContractAddress = sellingContractAddress
	input.BuyingAssetCode = nativeCode
	input.BuyingContractAddress = nativeIssuer

	if sellingAssetCode == nativeCode && sellingContractAddress == nativeIssuer {
		return "1", nil
	}
	orderBook, err := GetBantuOrderBookSummary(input)

	if len(orderBook.Asks) == 0 || err != nil {
		return "0", &tErrors.ErrorTemporaryServerError{}
	}
	nativePrice = orderBook.Asks[0].Price

	priceCache.Price = nativePrice

	gc.RedisCache.StoreResultToCacheRaw(cacheKey, priceCache, 120)

	return nativePrice, nil
}

// GetOrderBook
func GetOrderBook(assetCode, contractAddress, currencyCode, currencyIssuer string) (trovoOrderBook OrderBook, err error) {
	trovoOrderBook.Asks = make([]OfferVolume, 0)
	trovoOrderBook.Bids = make([]OfferVolume, 0)
	trovoOrderBook.Currency = Asset{
		AssetCode:       currencyCode,
		ContractAddress: currencyIssuer,
	}
	trovoOrderBook.Asset = Asset{
		AssetCode:       assetCode,
		ContractAddress: contractAddress,
	}
	if strings.EqualFold(assetCode, os.Getenv("NATIVE_ASSET_CODE")) {
		assetCode = ""
		contractAddress = ""
	}
	if strings.EqualFold(currencyCode, os.Getenv("NATIVE_ASSET_CODE")) {
		currencyCode = ""
		currencyIssuer = ""
	}

	var input OrderBookRequestInput

	input.SellingAssetCode = assetCode
	input.SellingContractAddress = contractAddress
	input.BuyingAssetCode = currencyCode
	input.BuyingContractAddress = currencyIssuer

	orderBook, err := GetBantuOrderBookSummary(input)
	if err != nil {
		log.Printf("[GetOrderBook]Error getting order book summary: %v\n", err)
		return
	}

	//process bids

	for _, bid := range orderBook.Bids {
		trovoOrderBook.Bids = append(trovoOrderBook.Bids, OfferVolume{
			Price:    bid.Price,
			Quantity: bid.Amount,
		})
	}

	//process asks

	for _, ask := range orderBook.Asks {
		trovoOrderBook.Asks = append(trovoOrderBook.Asks, OfferVolume{
			Price:    ask.Price,
			Quantity: ask.Amount,
		})
	}

	return trovoOrderBook, nil

}

// GetChartRecords
func GetChartRecords(assetCode, contractAddress, currencyCode, currencyIssuer, startTime, endTime, order, limit, chartPeriod, offset string) (tradeChart TradeChart, err error) {
	tradeChart.ChartRecords = make([]ChartRecord, 0)
	// tradeChart.Currency = Asset{
	// 	AssetCode:   currencyCode,
	// 	ContractAddress: currencyIssuer,
	// }
	// tradeChart.Asset = Asset{
	// 	AssetCode:   assetCode,
	// 	ContractAddress: contractAddress,
	// }
	if strings.EqualFold(assetCode, os.Getenv("NATIVE_ASSET_CODE")) {
		assetCode = ""
		contractAddress = ""
	}
	if strings.EqualFold(currencyCode, os.Getenv("NATIVE_ASSET_CODE")) {
		currencyCode = ""
		currencyIssuer = ""
	}

	var input TradeAggregateInput

	input.BaseAssetCode = assetCode
	input.BaseContractAddress = contractAddress
	input.CounterAssetCode = currencyCode
	input.CounterContractAddress = currencyIssuer
	input.Order = order
	input.Limit = limit

	if len(startTime) > 0 {
		layout := "2006-01-02 15:04:05 -0700 UTC"
		t, err := time.Parse(layout, startTime)
		if err != nil {
			fmt.Printf("[GetChartRecords] error parsing start time %v, error: %v", startTime, err)
		} else {
			input.StartTime = t
		}

	}
	if len(endTime) > 0 {
		layout := "2006-01-02 15:04:05 -0700 UTC"
		t, err := time.Parse(layout, endTime)
		if err != nil {
			fmt.Printf("[GetChartRecords] error parsing end time %v, error: %v", endTime, err)
		} else {
			input.EndTime = t
		}

	}

	if len(chartPeriod) > 0 {
		i, err := strconv.ParseInt(chartPeriod, 10, 64)
		if err != nil {
			input.Resolution = time.Duration(5)
		} else {
			input.Resolution = time.Duration(i)
		}

	} else {
		input.Resolution = time.Duration(5)
	}

	if len(offset) > 0 {
		i, err := strconv.ParseInt(offset, 10, 64)
		if err != nil {
			input.Offset = time.Duration(0)
		} else {
			input.Offset = time.Duration(i)
		}

	}

	tds, err := getTradeAggregate(input)
	if err != nil {
		log.Printf("[GetChartRecords]Error getting order book summary: %v\n", err)
		return
	}

	tradeChart, err = TransformTradeAggregationInstance(assetCode, contractAddress, currencyCode, currencyIssuer, &tds)

	return tradeChart, err

}

// TransformTradeAggregationInstance
func TransformTradeAggregationInstance(assetCode, contractAddress, currencyCode, currencyIssuer string, tds *TradeAggregationsPage) (tradeChart TradeChart, err error) {
	tradeChart.ChartRecords = make([]ChartRecord, 0)
	if assetCode == "" {
		assetCode = os.Getenv("NATIVE_ASSET_CODE")
		contractAddress = ""
	}
	if currencyCode == "" {
		currencyCode = os.Getenv("NATIVE_ASSET_CODE")
		currencyIssuer = ""
	}
	tradeChart.Currency = Asset{
		AssetCode:       currencyCode,
		ContractAddress: currencyIssuer,
	}
	tradeChart.Asset = Asset{
		AssetCode:       assetCode,
		ContractAddress: contractAddress,
	}

	for _, r := range tds.Embedded.Records {
		tradeChart.ChartRecords = append(tradeChart.ChartRecords, ChartRecord{
			Timestamp:      r.Timestamp,
			TradeCount:     r.TradeCount,
			AssetVolume:    r.BaseVolume,
			CurrencyVolume: r.CounterVolume,
			Average:        r.Average,
			High:           r.High,
			Low:            r.Low,
			Open:           r.Open,
			Close:          r.Close,
		})
	}

	return tradeChart, nil

}

// ProcessOrderBookEvent
func ProcessOrderBookEvent(orderBook OrderBookSummary) (trovoOrderBook OrderBook, err error) {
	trovoOrderBook.Asks = make([]OfferVolume, 0)
	trovoOrderBook.Bids = make([]OfferVolume, 0)
	if len(orderBook.Selling.AssetCode) == 0 {
		trovoOrderBook.Asset = Asset{
			AssetCode: os.Getenv("NATIVE_ASSET_CODE"),
		}
	} else {
		trovoOrderBook.Asset = Asset{
			AssetCode:       orderBook.Selling.AssetCode,
			ContractAddress: orderBook.Selling.ContractAddress,
		}
	}
	if len(orderBook.Buying.AssetCode) == 0 {
		trovoOrderBook.Currency = Asset{
			AssetCode: os.Getenv("NATIVE_ASSET_CODE"),
		}
	} else {
		trovoOrderBook.Currency = Asset{
			AssetCode:       orderBook.Buying.AssetCode,
			ContractAddress: orderBook.Buying.ContractAddress,
		}
	}

	//process bids

	for _, bid := range orderBook.Bids {
		trovoOrderBook.Bids = append(trovoOrderBook.Bids, OfferVolume{
			Price:    bid.Price,
			Quantity: bid.Amount,
		})
	}

	//process asks

	for _, ask := range orderBook.Asks {
		trovoOrderBook.Asks = append(trovoOrderBook.Asks, OfferVolume{
			Price:    ask.Price,
			Quantity: ask.Amount,
		})
	}

	return trovoOrderBook, nil

}

func logDiscordFailedRequest(msg string) {
	discord.WebhookURL = "https://discord.com/api/webhooks/827986576415129663/wqMKp9wxB_fxs9Q3zlMKCNPGENXmD_ueUnL8hVCu1wmRfD2wkXAjfP85k1Ro_2_wGfiY"
	if len(os.Getenv("EXPANSION_NETWORK_ERROR_WEBHOOK")) > 50 {
		discord.WebhookURL = os.Getenv("EXPANSION_NETWORK_ERROR_WEBHOOK")
	}
	discord.Say(msg)
}
