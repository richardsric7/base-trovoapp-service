package models

import (
	"fmt"
	"strings"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Price string

func (p Price) GetMarketPrices(db *gorm.DB) (bestPrices []MarketPrice) {
	type Asset struct {
		ID              string `gorm:"size:12;primaryKey;" json:"id"`
		Inactive        uint   `gorm:"not null;default:0" json:"-"`
		ShowMarketPrice uint   `gorm:"not null;default:0" json:"-"`
	}
	type Currency struct {
		ID              string `gorm:"size:10;primaryKey" json:"id"`
		CurrencySymbol  string `gorm:"size:100;not null" json:"currencySymbol"`
		Inactive        uint   `gorm:"not null;default:0" json:"-"`
		ShowMarketPrice uint   `gorm:"not null;default:0" json:"-"`
	}
	// asset/currency, XBN/NGN
	bestPrices = make([]MarketPrice, 0)

	var currencies []Currency
	e := db.Where("inactive = ?", 0).Where("show_market_price = ?", 1).Find(&currencies).Error
	if e == nil {
		var assets []Asset
		e = db.Where("inactive = ?", 0).Where("show_market_price = ?", 1).Find(&assets).Error
		if e == nil {
			// begin processing
			for _, asset := range assets {
				for _, currency := range currencies {
					bestPrice := p.GetBestPrice(asset.ID, currency.ID, db)
					if bestPrice[fmt.Sprintf("%v/%v", asset.ID, currency.ID)] != "0" {
						bestPrices = append(bestPrices, MarketPrice{
							Pair:           fmt.Sprintf("%v/%v", asset.ID, currency.ID),
							Price:          bestPrice[fmt.Sprintf("%v/%v", asset.ID, currency.ID)],
							CurrencySymbol: currency.CurrencySymbol,
						})

					}

				}
			}
			// add XBN dolalr price
			// {
			// 	type XbnDollarPrice struct {
			// 	ID            string `gorm:"primaryKey"`
			// 	Asset         string `gorm:"index:idx_dollar_price_unique_asset,unique"`
			// 	Source        string `gorm:"size:100;index:idx_dollar_price_unique_asset,unique"`
			// 	AskRate       string `gorm:"size:100"`
			// 	BidRate       string `gorm:"size:100"`
			// 	LastTradeRate string `gorm:"size:100"`
			// 	LastUpdated   time.Time
			// }
			// var xbnPrice XbnDollarPrice
			// price := "0"
			// err := db.Where("source = ?", "BITTREX").Where("asset = ?", "XBN").First(&xbnPrice).Error
			// if err != nil {

			// } else {
			// 	price = xbnPrice.AskRate
			// }

			// bestPrices = append(bestPrices, MarketPrice{
			// 	Pair:           fmt.Sprintf("%v/%v", "XBN", "USD"),
			// 	Price:          price,
			// 	CurrencySymbol: "$",
			// })
			// }

		}
	}

	return

}

func (p Price) GetBestPrice(asset, currency string, db *gorm.DB) (bestPrice map[string]string) {
	// asset/currency, XBN/NGN
	bestPrice = make(map[string]string)
	type Result struct {
		AssetID    string
		AssetPrice float64
	}
	var result Result
	var e error
	if strings.EqualFold(currency, "USDT") {
		e = db.Raw(`SELECT asset_id, asset_price FROM offers WHERE offer_type = 'Sell' AND currency_id = ?
		AND asset_id = ? AND inactive = 0 AND offline = 0 AND max_trade_amount > min_trade_amount AND min_trade_amount >= 15 ORDER BY asset_price ASC LIMIT 1`, currency, asset).Scan(&result).Error

	} else {
		e = db.Raw(`SELECT asset_id, asset_price FROM offers WHERE offer_type = 'Sell' AND currency_id = ?
	 AND asset_id = ? AND inactive = 0 AND offline = 0 AND max_trade_amount > min_trade_amount AND min_trade_amount > round((15*asset_price)::numeric,2) ORDER BY asset_price ASC LIMIT 1`, currency, asset).Scan(&result).Error

	}
	if e != nil {
		result.AssetID = asset
		result.AssetPrice = 0
	}
	bestPrice[fmt.Sprintf("%v/%v", asset, currency)] = decimal.NewFromFloat(result.AssetPrice).String()
	return

}
