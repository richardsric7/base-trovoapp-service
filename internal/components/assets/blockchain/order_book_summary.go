package assets

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	bantupayerrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/network"
	"trovo-wallet-api/internal/sharedconfig"

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
type Asset struct {
	AssetCode   string `json:"assetCode"`
	AssetIssuer string `json:"assetIssuer"`
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
	StartTime          time.Time
	EndTime            time.Time
	Resolution         time.Duration
	Offset             time.Duration
	BaseAssetCode      string
	BaseAssetIssuer    string
	BaseAssetType      string
	CounterAssetCode   string
	CounterAssetIssuer string
	CounterAssetType   string
	Order              string
	Limit              string
}
type PriceCache struct {
	Price     string `json:"price"`
	PriceType string `json:"priceType,omitempty"`
}

// getTradeAggregate gets trade chart data
func getTradeAggregate(input TradeAggregateInput) (tds horizon.TradeAggregationsPage, err error) {
	client := network.GetBlockchainClient()
	var limit uint
	var order horizonclient.Order
	var baseAssetType, counterAssetType horizonclient.AssetType
	if (len(input.BaseAssetCode) == 0 && len(input.BaseAssetIssuer) == 0) || (input.BaseAssetCode == "native") {
		baseAssetType = horizonclient.AssetTypeNative
		input.BaseAssetIssuer = ""
		input.BaseAssetCode = ""
	} else if input.BaseAssetType == "credit_alphanum4" || len(input.BaseAssetCode) < 5 {
		baseAssetType = horizonclient.AssetType4
	} else if input.BaseAssetType == "credit_alphanum12" || len(input.BaseAssetCode) > 4 {
		baseAssetType = horizonclient.AssetType12
	}

	if (len(input.CounterAssetCode) == 0 && len(input.CounterAssetIssuer) == 0) || (input.CounterAssetCode == "native") {
		counterAssetType = horizonclient.AssetTypeNative
		input.CounterAssetIssuer = ""
		input.CounterAssetCode = ""
	} else if input.CounterAssetType == "credit_alphanum4" || len(input.CounterAssetCode) < 5 {
		counterAssetType = horizonclient.AssetType4
	} else if input.CounterAssetType == "credit_alphanum12" || len(input.CounterAssetCode) > 4 {
		counterAssetType = horizonclient.AssetType12
	}

	if input.Limit == "" {
		limit = 200
	} else {
		plimit, _ := strconv.ParseInt(input.Limit, 10, 64)
		limit = uint(plimit)
	}
	order = horizonclient.OrderDesc
	if input.Order == "asc" {
		order = horizonclient.OrderAsc
	}
	oRequest := horizonclient.TradeAggregationRequest{
		StartTime:          input.StartTime,
		EndTime:            input.EndTime,
		Resolution:         input.Resolution,
		Offset:             input.Offset,
		BaseAssetType:      baseAssetType,
		BaseAssetCode:      input.BaseAssetCode,
		BaseAssetIssuer:    input.BaseAssetIssuer,
		CounterAssetType:   counterAssetType,
		CounterAssetCode:   input.CounterAssetCode,
		CounterAssetIssuer: input.CounterAssetIssuer,
		Order:              order,
		Limit:              limit,
	}
	// fmt.Printf("Offer Request: %+v\n", oRequest)
	oSummary, err := client.TradeAggregations(oRequest)
	if err != nil {
		if strings.Contains(err.Error(), "tls") || strings.Contains(err.Error(), "timeout") || strings.Contains(err.Error(), "handshake") || strings.Contains(err.Error(), "read tcp") || strings.Contains(err.Error(), "connection reset by peer") || strings.Contains(err.Error(), "dial tcp") || strings.Contains(err.Error(), "no such host") {
			log.Println("[getTradeAggregate]", err)
			return tds, &bantupayerrors.ErrorTemporaryServerError{}
		}
		hError := err.(*horizonclient.Error)
		//something went wrong, verify stage and check approprate action
		rCode, _ := hError.ResultCodes()
		rS, _ := hError.ResultString()
		log.Println("\n[getTradeAggregate] Problem in Request:", hError.Problem)
		log.Println("\n[getTradeAggregate] Result Codes in Request:", rCode)
		log.Println("\n[getTradeAggregate] Result String in Request:", rS)
		log.Printf("\n[getTradeAggregate] Problem in Request - RESPONSE: %+v\n", hError.Response)
		log.Println("[getTradeAggregate] Error submitting:", err)
		return tds, &bantupayerrors.ErrorTemporaryServerError{}

	}

	return oSummary, nil

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
	} else if input.BuyingAssetType == "credit_alphanum4" || len(input.BuyingAssetCode) < 5 {
		buyingAssetType = horizonclient.AssetType4
	} else if input.BuyingAssetType == "credit_alphanum12" || len(input.BuyingAssetCode) > 4 {
		buyingAssetType = horizonclient.AssetType12
	}

	if input.Limit == "" {
		limit = 200
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
		log.Println("[client.OrderBookRequest]error:", err)
		if strings.Contains(err.Error(), "tls") || strings.Contains(err.Error(), "timeout") || strings.Contains(err.Error(), "handshake") || strings.Contains(err.Error(), "read tcp") || strings.Contains(err.Error(), "connection reset by peer") || strings.Contains(err.Error(), "dial tcp") || strings.Contains(err.Error(), "no such host") {
			log.Println("[client.OrderBookRequest]", err)
			return orderBookSummary, &bantupayerrors.ErrorTemporaryServerError{}
		}
		// hError := err.(*horizonclient.Error)
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
			return orderBookSummary, &bantupayerrors.ErrorTemporaryServerError{}
		} else {
			// Assertion failed
			return orderBookSummary, &bantupayerrors.ErrorTemporaryServerError{}
		}

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

// GetDollarPrice dollar ask price using USDB
func GetDollarPrice(sellingAssetCode, sellingAssetIssuer string, gc *sharedconfig.GlobalConfig, checkCacheFirst bool) (usdPrice, priceType string, err error) {
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
	cacheKey := fmt.Sprintf("%v.%v_dollar", sellingAssetCode, sellingAssetIssuer)
	if checkCacheFirst {
		ok, concatPriceByte := gc.RedisCache.GetCachedResultRaw(cacheKey)
		if ok {
			json.Unmarshal(concatPriceByte, &priceCache)
			// log.Printf("[GetDollarPrice] cache result: %+v\n", priceCache)

			return priceCache.Price, priceCache.PriceType, nil
		}
	}

	input.SellingAssetCode = sellingAssetCode
	input.SellingAssetIssuer = sellingAssetIssuer
	dollarAsset := strings.Split(os.Getenv("DOLLAR_ASSET"), ":")
	if len(dollarAsset) != 2 {
		return "0", priceType, &bantupayerrors.ErrorTemporaryServerError{}
	}

	if strings.EqualFold(sellingAssetCode, dollarAsset[0]) && strings.EqualFold(sellingAssetIssuer, dollarAsset[1]) {
		//it is dollar asset
		return "1", priceType, nil
	}

	if strings.HasPrefix(sellingAssetCode, "USD") || strings.HasSuffix(sellingAssetCode, "USD") {
		return "1", priceType, nil
	}

	input.BuyingAssetCode = dollarAsset[0]
	input.BuyingAssetIssuer = dollarAsset[1]

	orderBook, err := getBantuOrderBookSummary(input)
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
			return "0", priceType, &bantupayerrors.ErrorTemporaryServerError{}
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

// GetNativeAskPrice native (XBN) ask price
func GetNativeAskPrice(sellingAssetCode, sellingAssetIssuer string, gc *sharedconfig.GlobalConfig, checkCacheFirst bool) (nativePrice string, err error) {
	var priceCache PriceCache
	var nativeCode, nativeIssuer string
	nv := strings.Split(os.Getenv("USE_ASSET_FOR_NATIVE_PRICE"), ":")
	if len(nv) == 2 {
		nativeCode = nv[0]
		nativeIssuer = nv[1]
	}
	cacheKey := fmt.Sprintf("%v.%v_nativePrice", sellingAssetCode, sellingAssetIssuer)
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
	input.SellingAssetIssuer = sellingAssetIssuer
	input.BuyingAssetCode = nativeCode
	input.BuyingAssetIssuer = nativeIssuer

	if sellingAssetCode == nativeCode && sellingAssetIssuer == nativeIssuer {
		return "1", nil
	}
	orderBook, err := getBantuOrderBookSummary(input)

	if len(orderBook.Asks) == 0 || err != nil {
		return "0", &bantupayerrors.ErrorTemporaryServerError{}
	}
	nativePrice = orderBook.Asks[0].Price

	priceCache.Price = nativePrice

	gc.RedisCache.StoreResultToCacheRaw(cacheKey, priceCache, 120)

	return nativePrice, nil
}

// GetOrderBook
func GetOrderBook(assetCode, assetIssuer, currencyCode, currencyIssuer string) (trovoOrderBook OrderBook, err error) {
	trovoOrderBook.Asks = make([]OfferVolume, 0)
	trovoOrderBook.Bids = make([]OfferVolume, 0)
	trovoOrderBook.Currency = Asset{
		AssetCode:   currencyCode,
		AssetIssuer: currencyIssuer,
	}
	trovoOrderBook.Asset = Asset{
		AssetCode:   assetCode,
		AssetIssuer: assetIssuer,
	}
	if strings.EqualFold(assetCode, os.Getenv("NATIVE_ASSET_CODE")) {
		assetCode = ""
		assetIssuer = ""
	}
	if strings.EqualFold(currencyCode, os.Getenv("NATIVE_ASSET_CODE")) {
		currencyCode = ""
		currencyIssuer = ""
	}

	var input OrderBookRequestInput

	input.SellingAssetCode = assetCode
	input.SellingAssetIssuer = assetIssuer
	input.BuyingAssetCode = currencyCode
	input.BuyingAssetIssuer = currencyIssuer

	orderBook, err := getBantuOrderBookSummary(input)
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
func GetChartRecords(assetCode, assetIssuer, currencyCode, currencyIssuer, startTime, endTime, order, limit, chartPeriod, offset string) (tradeChart TradeChart, err error) {
	tradeChart.ChartRecords = make([]ChartRecord, 0)
	// tradeChart.Currency = Asset{
	// 	AssetCode:   currencyCode,
	// 	AssetIssuer: currencyIssuer,
	// }
	// tradeChart.Asset = Asset{
	// 	AssetCode:   assetCode,
	// 	AssetIssuer: assetIssuer,
	// }
	if strings.EqualFold(assetCode, os.Getenv("NATIVE_ASSET_CODE")) {
		assetCode = ""
		assetIssuer = ""
	}
	if strings.EqualFold(currencyCode, os.Getenv("NATIVE_ASSET_CODE")) {
		currencyCode = ""
		currencyIssuer = ""
	}

	var input TradeAggregateInput

	input.BaseAssetCode = assetCode
	input.BaseAssetIssuer = assetIssuer
	input.CounterAssetCode = currencyCode
	input.CounterAssetIssuer = currencyIssuer
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

	tradeChart, err = TransformTradeAggregationInstance(assetCode, assetIssuer, currencyCode, currencyIssuer, &tds)

	return tradeChart, err

}

// TransformTradeAggregationInstance
func TransformTradeAggregationInstance(assetCode, assetIssuer, currencyCode, currencyIssuer string, tds *horizon.TradeAggregationsPage) (tradeChart TradeChart, err error) {
	tradeChart.ChartRecords = make([]ChartRecord, 0)
	if assetCode == "" {
		assetCode = os.Getenv("NATIVE_ASSET_CODE")
		assetIssuer = ""
	}
	if currencyCode == "" {
		currencyCode = os.Getenv("NATIVE_ASSET_CODE")
		currencyIssuer = ""
	}
	tradeChart.Currency = Asset{
		AssetCode:   currencyCode,
		AssetIssuer: currencyIssuer,
	}
	tradeChart.Asset = Asset{
		AssetCode:   assetCode,
		AssetIssuer: assetIssuer,
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
func ProcessOrderBookEvent(orderBook horizon.OrderBookSummary) (trovoOrderBook OrderBook, err error) {
	trovoOrderBook.Asks = make([]OfferVolume, 0)
	trovoOrderBook.Bids = make([]OfferVolume, 0)
	if len(orderBook.Selling.Code) == 0 {
		trovoOrderBook.Asset = Asset{
			AssetCode: os.Getenv("NATIVE_ASSET_CODE"),
		}
	} else {
		trovoOrderBook.Asset = Asset{
			AssetCode:   orderBook.Selling.Code,
			AssetIssuer: orderBook.Selling.Issuer,
		}
	}
	if len(orderBook.Buying.Code) == 0 {
		trovoOrderBook.Currency = Asset{
			AssetCode: os.Getenv("NATIVE_ASSET_CODE"),
		}
	} else {
		trovoOrderBook.Currency = Asset{
			AssetCode:   orderBook.Buying.Code,
			AssetIssuer: orderBook.Buying.Issuer,
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
