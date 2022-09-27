package assets

import (
	"encoding/json"
	"sync"
	blockchain "trovo-wallet-api/internal/components/assets/blockchain"
	models "trovo-wallet-api/internal/components/assets/models"
	bantuerrors "trovo-wallet-api/internal/errors"

	"github.com/shopspring/decimal"
	"github.com/toorop/go-bittrex"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// GetCuratedAssets returns list of Curated Assets
func GetCuratedAssets(includeInactive bool, db *gorm.DB) (assets map[string]models.CuratedAsset, err error) {

	var fetchedAssets []models.CuratedAsset
	tempAssets := make(map[string]models.CuratedAsset)
	var m sync.Mutex
	var dberr error
	if !includeInactive {
		dberr = db.Preload(clause.Associations).Order("priority").Order("asset_code").Where("inactive = ?", 0).Find(&fetchedAssets).Error

	} else {
		dberr = db.Preload(clause.Associations).Order("priority").Order("asset_code").Find(&fetchedAssets).Error
	}

	if dberr != nil {
		return assets, &bantuerrors.ErrorTemporaryServerError{}
	}
	usdPrice, _ := blockchain.GetXBNDollarAskPrice(db)
	var wg sync.WaitGroup
	for _, v := range fetchedAssets {

		wg.Add(1)
		go func(v models.CuratedAsset) {
			defer wg.Done()
			//get native price
			// log.Printf(">>>>>>>>>>>>>>>Fetched Asset: Code: %v, Issuer: %v\n", v.AssetCode, v.AssetIssuer)
			nativePrice, _ := blockchain.GetNativeAskPrice(v.AssetCode, v.AssetIssuer)
			v.NativePrice = nativePrice
			usdPriceFloat := decimal.RequireFromString(usdPrice)
			nativePriceFloat := decimal.RequireFromString(nativePrice)
			v.UsdPrice = nativePriceFloat.Mul(usdPriceFloat).Truncate(7).String()
			m.Lock()
			tempAssets[v.AssetCode+":"+v.AssetIssuer] = v
			m.Unlock()
		}(v)

	}
	wg.Wait()
	tempAssets[":"] = models.CuratedAsset{UsdPrice: usdPrice,
		AssetLogo:    nativeLogo(),
		AssetName:    "Bantu Network Token",
		Description:  "Description: XBN is the native asset and network utility token isued by the Bantu Blockchain Foundation",
		Website:      "www.bantufoundation.org",
		ContactEmail: "ops@bantufoundation.org",
		Priority:     1,
	}
	assets = tempAssets

	return assets, nil
}

// GetAssetClasses returns list of Curated Assets
func GetAssetClasses(db *gorm.DB) (assetClassesOutput []models.AssetClassOutput, err error) {

	var assetClasses []models.AssetClass
	var assetClassO models.AssetClassOutput
	dberr := db.Find(&assetClasses).Error
	if dberr != nil {
		return assetClassesOutput, &bantuerrors.ErrorTemporaryServerError{}
	}
	// fmt.Printf("Assets: %+v\n", assets)
	for _, v := range assetClasses {
		assetClassO.ID = v.ID
		assetClassO.AssetClass = v.AssetClass
		assetClassesOutput = append(assetClassesOutput, assetClassO)
	}

	return assetClassesOutput, nil
}

// GetBittrexChart returns bittrex chart
func GetBittrexChart(db *gorm.DB) (charts []bittrex.Candle, err error) {
	charts = make([]bittrex.Candle, 0)

	var chart models.XbnMarketChart
	dberr := db.Where("source = ?", "bittrex").First(&chart).Error
	if dberr != nil {
		return charts, &bantuerrors.ErrorTemporaryServerError{}
	}

	e := json.Unmarshal([]byte(chart.ChartString), &charts)
	if e != nil {
		return charts, &bantuerrors.ErrorTemporaryServerError{}
	}
	return charts, nil
}
