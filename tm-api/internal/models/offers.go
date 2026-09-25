package models

import (
	p2pErrors "admin-panel-dashboard/internal/errors"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Asset struct {
	ID              string  `gorm:"size:12;primaryKey;check:,length(id) > 2" json:"id"`
	AssetIssuer     *string `gorm:"size:56;null" json:"-"`
	ImageThumbnail  *string `gorm:"null" json:"imageThumbnail"`
	Inactive        uint    `gorm:"not null;default:0" json:"-"`
	ShowMarketPrice uint    `gorm:"not null;default:0" json:"-"`
}
type AssetJSON struct {
	ID             string `gorm:"size:12;primaryKey" json:"id"`
	ImageThumbnail string `gorm:"null" json:"imageThumbnail,omitempty"`
}
type Offer struct {
	ID                               string        `gorm:"primaryKey;check:,length(id) > 2" json:"id"`
	CreatedAt                        time.Time     `gorm:"default:now()" json:"createdAt"`
	UpdatedAt                        time.Time     `gorm:"default:now()" json:"updatedAt"`
	OfferType                        string        `gorm:"size:10;not null;index:idx_unique_offer,unique" json:"offerType"`
	OrderType                        string        `gorm:"size:10;not null;index:idx_offer_order_type" json:"orderType"`
	Maker                            string        `gorm:"not null;index:idx_unique_offer,unique;check:,length(maker) < 50" json:"maker"`
	MakerPhone                       string        `gorm:"not null;" json:"makerPhone"`
	PaymentMethodID                  *string       `gorm:"null;" json:"paymentMethodId"`
	PaymentMethod                    PaymentMethod `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	PaymentChannelID                 *string       `gorm:"null;index:idx_offer_payment_channel_id" json:"paymentChannelId"`
	CurrencyPaymentChannelID         string        `gorm:"not null;index:idx_unique_offer,unique;index:idx_offer_currency_payment_channel_id" json:"currencyPaymentChannelId"`
	CurrencyPaymentMethodID          string        `gorm:"not null;index:idx_offer_currency_payment_method_id;" json:"currencyPaymentMethodlId"`
	CurrencyPaymentMethodCountryCode *string       `gorm:"size:10;null;index:idx_offer_currency_payment_method_country_code;" json:"currencyPaymentMethodCountryCode"`
	Currency                         Currency      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	CurrencyID                       string        `gorm:"size:100;not null;index:idx_unique_offer,unique" json:"currencyId"`
	AssetAmount                      float64       `gorm:"not null;check:,asset_amount >= 0" json:"assetAmount"`
	Asset                            Asset         `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	AssetID                          string        `gorm:"size:12;not null;index:idx_unique_offer,unique" json:"assetId"`
	AssetPrice                       float64       `gorm:"not null;check:,asset_price > 0" json:"assetPrice"`
	MinTradeAmount                   float64       `gorm:"not null;check:,min_trade_amount >= 0" json:"minTradeAmount"`
	MaxTradeAmount                   float64       `gorm:"not null;check:,max_trade_amount >= 0" json:"maxTradeAmount"`
	Remark                           *string       `gorm:"null" json:"remark"`
	Inactive                         uint          `gorm:"not null;default:0" json:"-"`
	Offline                          uint          `gorm:"not null;default:0" json:"-"`
	MaxTimePerTransaction            uint          `gorm:"not null;default:45" json:"maxTimePerTransaction"`
	MakerCountryCode                 *string       `gorm:"size:3;null;index:idx_offer_maker_country_code" json:"makerCountryCode"`
}

type OfferJSON struct {
	ID                               string            `json:"id"`
	CreatedAt                        time.Time         `json:"createdAt"`
	UpdatedAt                        time.Time         `json:"updatedAt"`
	OfferType                        string            `json:"offerType"`
	OrderType                        string            `json:"orderType"`
	Maker                            string            `json:"maker,omitempty"`
	MakerPhone                       string            `json:"makerPhone,omitempty"`
	PaymentChannelID                 string            `json:"paymentChannelId,omitempty"`
	PaymentMethodID                  string            `json:"paymentMethodId,omitempty"`
	PaymentMethod                    PaymentMethodJSON `json:"paymentMethod,omitempty"`
	CurrencyPaymentMethodID          string            `json:"currencyPaymentMethodId"`
	CurrencyPaymentChannelID         string            `json:"currencyPaymentChannelId"`
	CurrencyPaymentMethod            PaymentMethodJSON `json:"currencyPaymentMethod"`
	CurrencyPaymentMethodCountryCode string            `json:"currencyPaymentMethodCountryCode"`
	CurrencyID                       string            `json:"currencyId"`
	AssetAmount                      string            `json:"assetAmount"`
	AssetID                          string            `json:"assetId"`
	AssetPrice                       string            `json:"assetPrice"`
	MinTradeAmount                   string            `json:"minTradeAmount"`
	MaxTradeAmount                   string            `json:"maxTradeAmount"`
	MakerFee                         string            `json:"makerFee"`
	// TakerFee                         string      `json:"takerFee"`
	Remark                string     `json:"remark,omitempty"`
	Inactive              uint       `json:"inactive"`
	MakerStats            MakerStats `json:"makerStats"`
	MakerVerified         uint       `json:"makerVerified"`
	MaxTimePerTransaction uint       `gorm:"not null;default:45" json:"maxTimePerTransaction"`
	MakerCountryCode      string     `json:"makerCountryCode"`
	FormattedText         string     `json:"formattedText"`
	Offline               uint       `json:"offline"`
}

type MakerStats struct {
	ID             string `gorm:"primaryKey;" json:"-"`
	Trades         string `json:"trades"`
	Cancelations   string `json:"cancelations"`
	Completed      string `json:"completed"`
	Pending        string `json:"pending"`
	ReportsAgainst string `json:"reportsAgainst"`
	Rank           string `json:"rank"`
}
type TakerStats struct {
	ID     string `gorm:"primaryKey;" json:"-"`
	Trades string `json:"trades"`
	Rank   string `json:"rank"`
}
type PaginatedOffers struct {
	Pages        int         `json:"pages"`
	CurrentPage  int         `json:"currentPage"`
	TotalRecords int         `json:"totalRecords"`
	Limit        int         `json:"limit"`
	Records      []OfferJSON `json:"records"`
}

type PaymentChannel struct {
	ID            string      `gorm:"size:100;primaryKey" json:"id"`
	PaymentType   PaymentType `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	PaymentTypeID string      `gorm:"size:100;not null" json:"paymentType"`
	Inactive      uint        `gorm:"not null;default:0" json:"-"`
}

type PaymentType struct {
	ID       string `gorm:"size:100;primaryKey" json:"id"`
	Inactive uint   `gorm:"not null;default:0" json:"-"`
}

type Currency struct {
	ID              string `gorm:"size:10;primaryKey" json:"id"`
	Name            string `gorm:"size:20;not null" json:"name"`
	CurrencySymbol  string `gorm:"size:100;not null" json:"currencySymbol"`
	Inactive        uint   `gorm:"not null;default:0" json:"-"`
	ShowMarketPrice uint   `gorm:"not null;default:0" json:"-"`
	DecimalPlaces   uint   `gorm:"not null;default:2" json:"decimalPlaces"`
}

// CurrencyConfig represents the model for currency configurations in the database
type CurrencyConfig struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Code        string `gorm:"uniqueIndex;not null" json:"code"`
	Description string `json:"description"`
}

func (pID CurrencyPaymentMethod) GetCurrencyPaymentMethod(db *gorm.DB) (paymentMethod PaymentMethod) {
	e := db.Where("id = ?", string(pID)).First(&paymentMethod).Error
	if e != nil {
		log.Println("[GetCurrencyPaymentMethod]could not get currency payment method, err:", e)

		return PaymentMethod{}
	}

	return paymentMethod
}

func (v *Offer) GetCurrencyPaymentMethod(db *gorm.DB) (paymentMethod PaymentMethod, err error) {
	e := db.Where("id = ?", v.CurrencyPaymentMethodID).First(&paymentMethod).Error
	if e != nil {
		if !errors.Is(e, gorm.ErrRecordNotFound) {
			log.Println("[GetCurrencyPaymentMethod]could not get currency payment method, err:", e)
			// LogDiscordError(fmt.Sprintf("[GetCurrencyPaymentMethod]could not get currency payment method, err:%v", e))

			err = &p2pErrors.ErrorTemporaryServerError{}
			return
		}
		err = &p2pErrors.ErrorPaymentMethodDoesNotExist{ID: v.CurrencyPaymentMethodID}
		return
	}

	return paymentMethod, nil
}
func (v *Offer) GetPaymentMethod(db *gorm.DB) (paymentMethod PaymentMethod, err error) {
	if !strings.EqualFold(v.OfferType, "buy") {
		return
	}
	e := db.Where("id = ?", v.PaymentMethodID).First(&paymentMethod).Error
	if e != nil {
		if !errors.Is(e, gorm.ErrRecordNotFound) {
			log.Println("[GetPaymentMethod]could not get payment method, err:", e)
			// LogDiscordError(fmt.Sprintf("[GetPaymentMethod]could not get payment method, err:%v", e))

			err = &p2pErrors.ErrorTemporaryServerError{}
			return
		}
		err = &p2pErrors.ErrorPaymentMethodDoesNotExist{ID: v.CurrencyPaymentMethodID}
		return
	}

	return paymentMethod, nil
}

func (v *Offer) ToJSON(db *gorm.DB) (offer OfferJSON) {
	return v.ToOfferJSON(db)
}

func (v *Offer) ToOfferJSON(db *gorm.DB) (offer OfferJSON) {

	offer = OfferJSON{
		Maker:                    strings.Split(v.Maker, "@")[0],
		MakerPhone:               v.MakerPhone,
		ID:                       v.ID,
		OfferType:                v.OfferType,
		OrderType:                v.OrderType,
		CurrencyPaymentMethodID:  v.CurrencyPaymentMethodID,
		CurrencyPaymentChannelID: v.CurrencyPaymentChannelID,
		CurrencyPaymentMethod:    CurrencyPaymentMethod(v.CurrencyPaymentMethodID).ToPaymentMethodJSON(db),
		CurrencyID:               v.CurrencyID,
		AssetAmount:              decimal.NewFromFloat(v.AssetAmount).Truncate(7).String(),
		AssetPrice:               decimal.NewFromFloat(v.AssetPrice).Truncate(4).String(),
		AssetID:                  v.AssetID,
		MinTradeAmount:           decimal.NewFromFloat(v.MinTradeAmount).RoundBank(2).String(),
		MaxTradeAmount:           decimal.NewFromFloat(v.MaxTradeAmount).RoundBank(2).String(),
		Inactive:                 v.Inactive,
		MaxTimePerTransaction:    v.MaxTimePerTransaction,
		Offline:                  v.Offline,
	}

	if v.Remark != nil {
		offer.Remark = *v.Remark
	}
	if v.CurrencyPaymentMethodCountryCode != nil {
		offer.CurrencyPaymentMethodCountryCode = *v.CurrencyPaymentMethodCountryCode
	}
	if v.MakerCountryCode != nil {
		offer.MakerCountryCode = *v.MakerCountryCode
	}

	if strings.EqualFold(v.OfferType, "buy") {
		if v.PaymentChannelID != nil {
			offer.PaymentChannelID = *v.PaymentChannelID
		}
		if v.PaymentMethodID != nil {
			offer.PaymentMethodID = *v.PaymentMethodID
		}
		pm, _ := v.GetPaymentMethod(db)
		offer.PaymentMethod = pm.ToPaymentMethodJSON()

	}
	if len(strings.ReplaceAll(os.Getenv("MAKER_FEE"), " ", "")) > 0 {
		offer.MakerFee = strings.ReplaceAll(os.Getenv("MAKER_FEE"), " ", "")
	} else {
		offer.MakerFee = "0"
	}
	// if len(strings.ReplaceAll(os.Getenv("TAKER_FEE"), " ", "")) > 0 {
	// 	offer.TakerFee = strings.ReplaceAll(os.Getenv("TAKER_FEE"), " ", "")
	// } else {
	// 	offer.TakerFee = "0"
	// }

	stats := v.getMakerStat(db)
	offer.MakerStats = stats
	if len(stats.Trades) > 0 && len(stats.Rank) > 0 {
		mTrades := decimal.RequireFromString(stats.Trades)
		mRank := decimal.RequireFromString(stats.Rank)
		if mTrades.GreaterThanOrEqual(decimal.NewFromFloat(200)) && mRank.GreaterThanOrEqual(decimal.NewFromFloat(80)) {
			offer.MakerVerified = 1
		} else if mTrades.GreaterThanOrEqual(decimal.NewFromFloat(100)) && mRank.GreaterThanOrEqual(decimal.NewFromFloat(90)) {
			offer.MakerVerified = 1
		}

	}

	return offer
}

func (offer *OfferJSON) ValidateOffer() error {

	if len(offer.AssetAmount) == 0 {
		return &p2pErrors.ErrorInvalidParameter{Param: "assetAmount", ErrMessage: "Asset amount is empty"}
	}
	if len(offer.AssetID) == 0 {
		return &p2pErrors.ErrorInvalidParameter{Param: "assetId", ErrMessage: "Asset is empty"}
	}
	if len(offer.AssetPrice) == 0 {
		return &p2pErrors.ErrorInvalidParameter{Param: "assetPrice", ErrMessage: "Asset price is empty"}
	}
	if len(offer.PaymentMethodID) == 0 && strings.EqualFold(offer.OfferType, "buy") {
		return &p2pErrors.ErrorInvalidParameter{Param: "paymentMethodId", ErrMessage: "Payment method id is empty"}
	}
	if len(offer.CurrencyID) == 0 {
		return &p2pErrors.ErrorInvalidParameter{Param: "currencyId", ErrMessage: "Currency  is empty"}
	}
	if len(offer.CurrencyPaymentMethodID) == 0 {
		return &p2pErrors.ErrorInvalidParameter{Param: "currencyPaymentMethodId", ErrMessage: offer.CurrencyID + " payment method is empty"}
	}
	if len(offer.MinTradeAmount) == 0 || offer.MinTradeAmount == "null" {
		return &p2pErrors.ErrorInvalidParameter{Param: "minTradeAmount", ErrMessage: "Min trade amount is empty"}
	}
	if len(offer.MaxTradeAmount) == 0 || offer.MaxTradeAmount == "null" {
		return &p2pErrors.ErrorInvalidParameter{Param: "maxTradeAmount", ErrMessage: "Max trade amount is empty"}
	}
	if len(offer.OfferType) == 0 {
		return &p2pErrors.ErrorInvalidParameter{Param: "offerType", ErrMessage: "Offer type is empty"}
	}
	if offer.Inactive > 1 {
		return &p2pErrors.ErrorInvalidParameter{Param: "inactive", ErrMessage: fmt.Sprintf("Value for inactive [%v] is invalid", offer.Inactive)}
	}

	if !(decimal.RequireFromString(offer.AssetAmount).IsPositive()) {
		return &p2pErrors.ErrorInvalidParameter{Param: "assetAmount", ErrMessage: "Asset amount must be greater than 0"}

	}

	if decimal.RequireFromString(offer.AssetAmount).LessThanOrEqual(decimal.RequireFromString("15")) && strings.EqualFold(offer.AssetID, "XBN") {
		return &p2pErrors.ErrorInvalidParameter{Param: "assetAmount", ErrMessage: "Asset amount for XBN must be greater than 15"}

	}
	if !(decimal.RequireFromString(offer.AssetPrice).IsPositive()) {
		return &p2pErrors.ErrorInvalidParameter{Param: "assetPrice", ErrMessage: "Asset price must be greater than 0"}

	}

	if !(decimal.RequireFromString(offer.MinTradeAmount).IsPositive()) {
		return &p2pErrors.ErrorInvalidParameter{Param: "minTradeAmount", ErrMessage: "Min trade amount must be greater than 0"}

	}

	if !((decimal.RequireFromString(offer.MaxTradeAmount).RoundBank(2)).IsPositive()) {
		return &p2pErrors.ErrorInvalidParameter{Param: "maxTradeAmount", ErrMessage: "Max trade amount must be greater than 0"}

	}
	if decimal.RequireFromString(offer.MaxTradeAmount).LessThanOrEqual(decimal.RequireFromString(offer.MinTradeAmount).RoundBank(2)) {
		return &p2pErrors.ErrorInvalidParameter{Param: "maxTradeAmount", ErrMessage: "Max trade amount must be greater than min trade amount"}

	}

	if decimal.RequireFromString(offer.MinTradeAmount).LessThanOrEqual((decimal.RequireFromString("15").Mul(decimal.RequireFromString(offer.AssetPrice))).RoundBank(2)) && strings.EqualFold(offer.AssetID, "XBN") {
		return &p2pErrors.ErrorInvalidParameter{Param: "minTradeAmount", ErrMessage: fmt.Sprintf("Min trade amount must be greater than %v %v", ((decimal.RequireFromString("15").Mul(decimal.RequireFromString(offer.AssetPrice))).RoundBank(2)), offer.CurrencyID)}

	}
	if (decimal.RequireFromString(offer.MaxTradeAmount).RoundBank(2)).GreaterThan((decimal.RequireFromString(offer.AssetAmount).Mul(decimal.RequireFromString(offer.AssetPrice))).RoundBank(2)) {
		return &p2pErrors.ErrorInvalidParameter{Param: "maxTradeAmount", ErrMessage: fmt.Sprintf("Max trade amount must be less than %v %v", (decimal.RequireFromString(offer.AssetAmount).Mul(decimal.RequireFromString(offer.AssetPrice))).RoundBank(2).String(), offer.CurrencyID)}

	}

	return nil
}

func (offer *OfferJSON) ValidateOfferUpdate() error {

	if len(offer.AssetAmount) == 0 {
		return &p2pErrors.ErrorInvalidParameter{Param: "assetAmount", ErrMessage: "Asset amount is empty"}
	}

	if len(offer.AssetPrice) == 0 {
		return &p2pErrors.ErrorInvalidParameter{Param: "assetPrice", ErrMessage: "Asset price is empty"}
	}
	if len(offer.MinTradeAmount) == 0 || offer.MinTradeAmount == "null" {
		return &p2pErrors.ErrorInvalidParameter{Param: "minTradeAmount", ErrMessage: "Min trade amount is empty"}
	}
	if len(offer.MaxTradeAmount) == 0 || offer.MaxTradeAmount == "null" {
		return &p2pErrors.ErrorInvalidParameter{Param: "maxTradeAmount", ErrMessage: "Max trade amount is empty"}
	}
	if offer.Inactive > 1 {
		return &p2pErrors.ErrorInvalidParameter{Param: "inactive", ErrMessage: fmt.Sprintf("Value for inactive [%v] is invalid", offer.Inactive)}
	}

	if !(decimal.RequireFromString(offer.AssetAmount).IsPositive()) {
		return &p2pErrors.ErrorInvalidParameter{Param: "assetAmount", ErrMessage: "Asset amount must be greater than 0"}

	}
	if decimal.RequireFromString(offer.AssetAmount).LessThanOrEqual(decimal.RequireFromString("15")) && strings.EqualFold(offer.AssetID, "XBN") {
		return &p2pErrors.ErrorInvalidParameter{Param: "assetAmount", ErrMessage: "Asset amount for XBN must be greater than 15"}

	}
	if !(decimal.RequireFromString(offer.AssetPrice).IsPositive()) {
		return &p2pErrors.ErrorInvalidParameter{Param: "assetPrice", ErrMessage: "Asset price must be greater than 0"}

	}

	if !(decimal.RequireFromString(offer.MinTradeAmount).IsPositive()) {
		return &p2pErrors.ErrorInvalidParameter{Param: "minTradeAmount", ErrMessage: "Min trade amount must be greater than 0"}

	}

	if !((decimal.RequireFromString(offer.MaxTradeAmount).RoundBank(2)).IsPositive()) {
		return &p2pErrors.ErrorInvalidParameter{Param: "maxTradeAmount", ErrMessage: "Max trade amount must be greater than 0"}

	}
	if decimal.RequireFromString(offer.MaxTradeAmount).LessThanOrEqual(decimal.RequireFromString(offer.MinTradeAmount).RoundBank(2)) {
		return &p2pErrors.ErrorInvalidParameter{Param: "maxTradeAmount", ErrMessage: "Max trade amount must be greater than min trade amount"}

	}
	if decimal.RequireFromString(offer.MinTradeAmount).LessThanOrEqual((decimal.RequireFromString("15").Mul(decimal.RequireFromString(offer.AssetPrice))).RoundBank(2)) && strings.EqualFold(offer.AssetID, "XBN") {
		return &p2pErrors.ErrorInvalidParameter{Param: "minTradeAmount", ErrMessage: fmt.Sprintf("Min trade amount must be greater than %v %v", ((decimal.RequireFromString("15").Mul(decimal.RequireFromString(offer.AssetPrice))).RoundBank(2)), offer.CurrencyID)}

	}
	if (decimal.RequireFromString(offer.MaxTradeAmount).RoundBank(2)).GreaterThan((decimal.RequireFromString(offer.AssetAmount).Mul(decimal.RequireFromString(offer.AssetPrice))).RoundBank(2)) {
		return &p2pErrors.ErrorInvalidParameter{Param: "maxTradeAmount", ErrMessage: fmt.Sprintf("Max trade amount must be less than %v %v", (decimal.RequireFromString(offer.AssetAmount).Mul(decimal.RequireFromString(offer.AssetPrice))).RoundBank(2).String(), offer.CurrencyID)}

	}

	return nil
}

func (v *OfferJSON) ToFormattedText(i int) (formattedText string) {
	// maker := strings.Split(v.Maker, "@")[0]
	account := ""
	payment := ""
	if strings.EqualFold(v.OfferType, "Buy") {
		pm := v.PaymentMethod
		account = pm.DestinationAccount[0:10] + "..."
		payment = fmt.Sprintf("wallet:%v, p/w: %v", account, v.CurrencyPaymentChannelID)

	}
	if strings.EqualFold(v.OfferType, "Sell") {
		cpm := v.CurrencyPaymentMethod
		account := cpm.DestinationAccount
		if len(cpm.DestinationAccount) > 13 {
			account = account[0:10] + "..."
		}

		payment = fmt.Sprintf("%v:%v", account, cpm.PaymentChannelID)

	}
	offerState := "Active"
	if v.Inactive == 1 {
		offerState = "Inactive"
	}
	formattedText = fmt.Sprintf("%v. %v %v %v @ %v %v 🔀 %v| %v| %vmins| Min/Max:%v/%v", i+1, v.OfferType, v.AssetAmount, v.AssetID, v.AssetPrice, v.CurrencyID, payment, offerState, v.MaxTimePerTransaction, v.MinTradeAmount, v.MaxTradeAmount)
	v.FormattedText = formattedText
	return
}

func (v *Offer) getMakerStat(db *gorm.DB) (rating MakerStats) {

	db.First(&rating, MakerStats{ID: v.Maker})

	return
}
