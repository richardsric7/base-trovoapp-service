package assets

import (
	"encoding/json"
	"log"
	"sync"
	models "trovo-wallet-api/internal/components/assets/models"
	bantuerrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/sharedconfig"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// GetCuratedAssets returns list of Curated Assets
func GetCuratedAssets(includeInactive bool, gc *sharedconfig.GlobalConfig) (assets map[string]models.CuratedAsset) {
	var fetchedAssets []models.CuratedAsset
	tempAssets := make(map[string]models.CuratedAsset)
	assets = make(map[string]models.CuratedAsset)

	cacheKeyInfo := "curatedAssets_"
	{

		// search cache for balance
		ok, rawdata := gc.RedisCache.GetCachedResultRaw(cacheKeyInfo)

		if ok {

			// log.Printf("[GetCuratedAssets] %v, served from cache\n", cacheKeyInfo)
			json.Unmarshal(rawdata, &assets)
			return
		}

	}

	var m sync.Mutex
	var dberr error
	if !includeInactive {
		dberr = gc.DB.Preload(clause.Associations).Order("priority").Order("asset_code").Where("inactive = ?", 0).Find(&fetchedAssets).Error

	} else {
		dberr = gc.DB.Preload(clause.Associations).Order("priority").Order("asset_code").Find(&fetchedAssets).Error
	}

	if dberr != nil {
		log.Printf("[GetCuratedAssets]error getting assets: %v\n", dberr)
		return assets
	}
	// usdPrice, _ := blockchain.GetXBNDollarAskPrice(db)
	var wg sync.WaitGroup
	for _, v := range fetchedAssets {

		wg.Add(1)
		go func(v models.CuratedAsset) {
			defer wg.Done()
			//get native price
			// log.Printf(">>>>>>>>>>>>>>>Fetched Asset: Code: %v, Issuer: %v\n", v.AssetCode, v.AssetIssuer)
			// nativePrice, _ := blockchain.GetNativeAskPrice(v.AssetCode, v.AssetIssuer)
			// v.NativePrice = nativePrice
			// usdPriceFloat := decimal.RequireFromString(usdPrice)
			// nativePriceFloat := decimal.RequireFromString(nativePrice)
			// v.UsdPrice = nativePriceFloat.Mul(usdPriceFloat).Truncate(7).String()
			m.Lock()
			tempAssets[v.AssetCode+":"+v.AssetIssuer] = v
			m.Unlock()
		}(v)

	}
	wg.Wait()
	// tempAssets[":"] = models.CuratedAsset{
	// 	ImageURL:     nativeLogo(),
	// 	AssetName:    "Bantu Network Token",
	// 	Description:  "Description: XBN is the native asset and network utility token issued by the Bantu Blockchain Foundation",
	// 	Website:      "www.bantufoundation.org",
	// 	ContactEmail: "ops@bantufoundation.org",
	// 	Priority:     1,
	// }
	assets = tempAssets
	gc.RedisCache.StoreResultToCacheRaw(cacheKeyInfo, assets, 0)

	return assets
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
