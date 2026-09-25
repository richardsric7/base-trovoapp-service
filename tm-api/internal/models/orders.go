package models

import (
	"errors"
	"fmt"
	"log"
	"os"

	p2pErrors "admin-panel-dashboard/internal/errors"
	"time"

	"github.com/ecnepsnai/discord"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Order struct {
	ID                                             string `gorm:"size:100;primaryKey;check:,length(id) > 2" json:"id"`
	OrderType                                      string `gorm:"not null" json:"orderType"`
	CreatedAt                                      time.Time
	UpdatedAt                                      time.Time
	ExpiresAt                                      time.Time
	AcceptedAt                                     time.Time
	CancelAfter                                    time.Time   `json:"cancelAfter"`
	OfferID                                        string      `gorm:"not null" json:"offerId"`
	OfferType                                      string      `gorm:"not null" json:"offerType"`
	OfferMaker                                     string      `gorm:"not null;index:idx_order_offer_maker" json:"offerMaker"`
	OfferMakerPhone                                string      `gorm:"not null" json:"offerMakerPhone"`
	OfferMakerCountryCode                          *string     `gorm:"null" json:"offerMakerCountryCode"`
	OfferMaxTimePerTransaction                     uint        `gorm:"not null;default:45" json:"offerMaxTimePerTransaction"`
	OfferTaker                                     string      `gorm:"not null;index:idx_order_offer_taker" json:"offerTaker"`
	OfferTakerPhone                                string      `gorm:"not null" json:"offerTakerPhone"`
	OfferPaymentMethodID                           *string     `gorm:"null;" json:"offerPaymentMethodId"`
	OfferPaymentChannelID                          *string     `gorm:"null" json:"offerPaymentChannelId"`
	OfferPaymentMethodName                         *string     `gorm:"null;" json:"offerPaymentMethodName"` // beneficiary name
	OfferPaymentMethodDestinationAccount           *string     `gorm:"null" json:"offerPaymentMethodDestinationAccount"`
	OfferPaymentMethodMemo                         *string     `gorm:"null;" json:"offerPaymentMethodMemo"`
	OfferPaymentMethodBankName                     *string     `gorm:"null;" json:"offerPaymentMethodBankName"`
	OfferPaymentMethodAccountOpeningBranch         *string     `gorm:"null;" json:"offerPaymentMethodAccountOpeningBranch"`
	OfferCurrencyPaymentMethodID                   string      `gorm:"not null" json:"offerCurrencyPaymentMethodId"`
	OfferCurrencyPaymentChannelID                  string      `gorm:"not null" json:"offerCurrencyPaymentChannelId"`
	OfferCurrencyPaymentMethodName                 *string     `gorm:"null;" json:"offerCurrencyPaymentMethodName"` // beneficiary name
	OfferCurrencyPaymentMethodDestinationAccount   string      `gorm:"not null" json:"offerCurrencyPaymentMethodDestinationAccount"`
	OfferCurrencyPaymentMethodMemo                 *string     `gorm:"null;" json:"offerCurrencyPaymentMethodMemo"`
	OfferCurrencyPaymentMethodBankName             *string     `gorm:"null;" json:"offerPaymentCurrencyMethodBankName"`
	OfferCurrencyPaymentMethodAccountOpeningBranch *string     `gorm:"null;" json:"offerCurrencyPaymentMethodAccountOpeningBranch"`
	OfferCurrencyPaymentMethodCountryCode          *string     `gorm:"null;" json:"offerCurrencyPaymentMethodCountryCode"`
	OfferCurrencyPaymentMethodCurrencyID           *string     `gorm:"null;" json:"offerCurrencyPaymentMethodCurrencyId"`
	OfferCurrencyID                                string      `gorm:"not null" json:"offerCurrencyID"`
	OfferAssetAmount                               float64     `gorm:"not null;check:,offer_asset_amount > 0" json:"offerAssetAmount"`
	OfferAssetID                                   string      `gorm:"size:12;not null" json:"offerAssetId"`
	OfferAssetPrice                                float64     `gorm:"not null;check:,offer_asset_price > 0" json:"offerAssetPrice"`
	OfferMinTradeAmount                            float64     `gorm:"not null;check:,offer_min_trade_amount >= 0" json:"offerMinTradeAmount"`
	OfferMaxTradeAmount                            float64     `gorm:"not null;check:,offer_max_trade_amount >= 0" json:"offerMaxTradeAmount"`
	OfferRemark                                    *string     `gorm:"null" json:"offerRemark"`
	TakerPaymentMethodID                           string      `gorm:"not null;" json:"takerPaymentMethodId"`
	TakerPaymentChannelID                          string      `gorm:"not null" json:"takerPaymentChannelId"`
	TakerPaymentMethodName                         *string     `gorm:"null;" json:"takerPaymentMethodName"`
	TakerPaymentMethodDestinationAccount           string      `gorm:"not null" json:"takerPaymentMethodDestinationAccount"`
	TakerPaymentMethodMemo                         *string     `gorm:"null;" json:"takerPaymentMethodMemo"`
	TakerPaymentMethodBankName                     *string     `gorm:"null;" json:"takerPaymentMethodBankName"`
	TakerPaymentMethodAccountOpeningBranch         *string     `gorm:"null;" json:"takerPaymentMethodAccountOpeningBranch"`
	TakerPaymentMethodCountryCode                  *string     `gorm:"size:10;null;" json:"takerPaymentMethodCountryCode"`
	TakerPaymentMethodCurrencyID                   *string     `gorm:"size:10;null;" json:"takerPaymentMethodCurrencyId"`
	OrderEscrowAddress                             string      `gorm:"not null;index:idx_order_payment_memo,unique" json:"orderEscrowAddress"`
	OrderPaymentMemo                               string      `gorm:"size:28;not null;index:idx_order_payment_memo,unique" json:"orderPaymentMemo"`
	OrderAmount                                    float64     `gorm:"not null;check:,order_amount > 0" json:"orderAmount"`
	OrderMakerFee                                  float64     `gorm:"not null;check:,order_maker_fee >= 0" json:"orderMakerFee"`
	OrderTakerFee                                  float64     `gorm:"not null;check:,order_taker_fee >= 0" json:"orderTakerFee"`
	OrderEscrowTransactionID                       *string     `gorm:"null;" json:"orderEscrowTransactionId"`
	OrderAssetReleaseTransactionID                 *string     `gorm:"null;" json:"orderAssetReleaseTransactionId"`
	OrderStatusID                                  uint        `gorm:"not null;default:1" json:"orderStatusId"`
	OrderStatus                                    OrderStatus `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"orderStatus"`
	DynamicLink                                    *string     `gorm:"null" json:"dynamicLink"`
	QRCode                                         *string     `gorm:"null" json:"qrCode"`
	FiatDepositTransactionID                       *string     `gorm:"null"`
}

type OrderStatus struct {
	ID     uint   `json:"-"`
	Status string `json:"status"`
}

type OrderInput struct {
	OfferID              string `json:"offerId"`
	OrderAmount          string `json:"orderAmount"`
	TakerPaymentMethodID string `json:"takerPaymentMethodId"`
}

type FiatConfirmInput struct {
	FiatDepositTransactionID string `json:"fiatDepositTransactionId"`
}

type PaginatedOrders struct {
	Pages        int         `json:"pages"`
	CurrentPage  int         `json:"currentPage"`
	TotalRecords int         `json:"totalRecords"`
	Limit        int         `json:"limit"`
	Records      []OrderJSON `json:"records"`
}
type TakerReputation struct {
	ID             string `gorm:"primaryKey;" json:"-"`
	Trades         string `json:"trades"`
	Cancelations   string `json:"cancelations"`
	Completed      string `json:"completed"`
	Pending        string `json:"pending"`
	ReportsAgainst string `json:"reportsAgainst"`
	Rank           string `json:"rank"`
}
type OrderChatMessage struct {
	ID        uint64    `gorm:"primaryKey" json:"messageId"`
	CreatedAt time.Time `gorm:"not null;index:idx_order_message_time;" json:"createdAt"`
	OrderID   string    `gorm:"not null" json:"orderId"`
	Order     Order     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Sender    string    `gorm:"not null" json:"sender"`
	Message   string    `gorm:"not null" json:"message"`
}
type OrderAppeal struct {
	ID           uint64    `gorm:"primaryKey" json:"messageId"`
	CreatedAt    time.Time `gorm:"not null" json:"createdAt"`
	OrderID      string    `gorm:"not null" json:"orderId"`
	Order        Order     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	AppealedBy   string    `gorm:"not null" json:"appealedBy"`
	AppealReason string    `gorm:"not null" json:"appealReason"`
	AppealDetail string    `gorm:"not null" json:"appealDetail"`
}

type OrderChatMessageInput struct {
	Message string `json:"message"`
}

type OrderJSON struct {
	ID                                             string             `json:"id"`
	CreatedAt                                      time.Time          `json:"createdAt"`
	ExpiresAt                                      time.Time          `json:"expiresAt"`
	AcceptedAt                                     time.Time          `json:"acceptedAt"`
	CancelAfter                                    time.Time          `json:"cancelAfter"`
	OfferID                                        string             `json:"offerId"`
	OfferType                                      string             `json:"offerType"`
	OfferMaker                                     string             `json:"offerMaker"`
	OfferMakerFullName                             string             `json:"offerMakerFullName"`
	OfferMaxTimePerTransaction                     uint               `json:"offerMaxTimePerTransaction"`
	OfferMakerPhone                                string             `json:"offerMakerPhone"`
	OfferMakerCountryCode                          string             `json:"offerMakerCountryCode"`
	OfferTaker                                     string             `json:"offerTaker"`
	OfferTakerFullName                             string             `json:"offerTakerFullName"`
	OfferTakerPhone                                string             `json:"offerTakerPhone"`
	OfferPaymentMethodID                           string             `json:"offerPaymentMethodId,omitempty"`
	OfferPaymentChannelID                          string             `json:"offerPaymentChannelId,omitempty"`
	OfferCurrencyPaymentChannelID                  string             `json:"offerCurrencyPaymentChannelId,omitempty"`
	OfferCurrencyPaymentMethodID                   string             `json:"offerCurrencyPaymentMethodId,omitempty"`
	OfferCurrencyPaymentMethod                     PaymentMethodJSON  `json:"offerCurrencyPaymentMethod,omitempty"`
	OfferCurrencyPaymentMethodName                 string             `json:"offerCurrencyPaymentMethodName,omitempty"` // beneficiary name
	OfferCurrencyPaymentMethodDestinationAccount   string             `json:"offerCurrencyPaymentMethodDestinationAccount,omitempty"`
	OfferCurrencyPaymentMethodMemo                 string             `json:"offerCurrencyPaymentMethodMemo,omitempty"`
	OfferCurrencyPaymentMethodBankName             string             `json:"offerPaymentCurrencyMethodBankName,omitempty"`
	OfferCurrencyPaymentMethodAccountOpeningBranch string             `json:"offerCurrencyPaymentMethodAccountOpeningBranch,omitempty"`
	OfferCurrencyPaymentMethodCountryCode          string             `json:"offerCurrencyPaymentMethodCountryCode"`
	OfferCurrencyPaymentMethodCurrencyID           string             `json:"offerCurrencyPaymentMethodCurrencyID"`
	OfferPaymentMethodName                         string             `json:"offerPaymentMethodName,omitempty"`
	OfferPaymentMethodDestinationAccount           string             `json:"offerPaymentMethodDestinationAccount,omitempty"`
	OfferPaymentMethodMemo                         string             `json:"offerPaymentMethodMemo,omitempty"`
	OfferPaymentMethodBankName                     string             `json:"offerPaymentMethodBankName,omitempty"`
	OfferPaymentMethodAccountOpeningBranch         string             `json:"offerPaymentMethodAccountOpeningBranch,omitempty"`
	OfferCurrencyID                                string             `json:"offerCurrencyId"`
	OfferAssetAmount                               string             `json:"offerAssetAmount"`
	OfferAssetID                                   string             `json:"offerAssetId"`
	OfferAssetPrice                                string             `json:"offerAssetPrice"`
	OfferMinTradeAmount                            string             `json:"offerMinTradeAmount"`
	OfferMaxTradeAmount                            string             `json:"offerMaxTradeAmount"`
	OfferRemark                                    string             `json:"offerRemark,omitempty"`
	TakerPaymentMethodID                           string             `json:"takerPaymentMethodId"`
	TakerPaymentChannelID                          string             `json:"takerPaymentChannelId"`
	TakerPaymentMethodName                         string             `json:"takerPaymentMethodName"`
	TakerPaymentMethodDestinationAccount           string             `json:"takerPaymentMethodDestinationAccount"`
	TakerPaymentMethodMemo                         string             `json:"takerPaymentMethodMemo,omitempty"`
	TakerPaymentMethodBankName                     string             `json:"takerPaymentMethodBankName,omitempty"`
	TakerPaymentMethodAccountOpeningBranch         string             `json:"takerPaymentMethodAccountOpeningBranch,omitempty"`
	TakerPaymentMethodCountryCode                  string             `json:"takerPaymentMethodCountryCode,omitempty"`
	TakerPaymentMethodCurrencyID                   string             `son:"takerPaymentMethodCurrencyId"`
	OrderType                                      string             `json:"orderType,omitempty"`
	OrderEscrowAddress                             string             `json:"orderEscrowAddress"`
	OrderPaymentMemo                               string             `json:"orderPaymentMemo"`
	OrderAmount                                    string             `json:"orderAmount"`
	OrderMakerFee                                  string             `json:"orderMakerFee"`
	OrderTakerFee                                  string             `json:"orderTakerFee"`
	OrderCurrencyAmount                            string             `json:"orderCurrencyAmount"`
	OrderEscrowTransactionID                       string             `json:"orderEscrowTransactionId,omitempty"`
	OrderStatusID                                  uint               `gorm:"not null;default:1" json:"orderStatusId"`
	OrderStatus                                    OrderStatus        `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"orderStatus"`
	DynamicLink                                    string             `gorm:"null" json:"dynamicLink,omitempty"` // only shows when order is accepted
	QRCode                                         string             `gorm:"null" json:"qrCode,omitempty"`
	FiatDepositTransactionID                       string             `gorm:"null" json:"fiatDepositTransactionId,omitempty"`
	OrderAssetReleaseTransactionID                 string             `json:"orderAssetReleaseTransactionId,omitempty"`
	OrderAssetReleaseExplorerLink                  string             `json:"orderAssetReleaseExplorerLink,omitempty"`
	OrderChatMessages                              []OrderChatMessage `json:"orderChatMessages,omitempty"`
	TakerReputation                                TakerReputation    `json:"takerReputation,omitempty"`
	FormattedText                                  string             `json:"formattedText,omitempty"`
	InAppeal                                       uint               `json:"inAppeal"`
}

type MerchantTradeReward struct {
	ID             uint64    `gorm:"primaryKey" json:"id"`
	CreatedAt      time.Time `gorm:"not null" json:"createdAt"`
	OrderID        string    `gorm:"not null" json:"orderId"`
	Order          Order     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Merchant       string    `gorm:"not null" json:"merchant"`
	Taker          string    `gorm:"not null" json:"taker"`
	TC             string    `gorm:"not null" json:"tc"`
	CR             string    `gorm:"not null" json:"cr"`
	XB             string    `gorm:"not null" json:"xb"`
	MerchantReward string    `gorm:"not null" json:"merchantReward"`
	TakerReward    string    `gorm:"not null" json:"takerReward"`
}

type AssetCast string
type OfferCast string

func (v *Order) GetTakerReputation(db *gorm.DB) (rep TakerReputation) {

	db.First(&rep, TakerReputation{ID: v.OfferTaker})

	return
}

type OrderID string

func (v OrderID) GetOrder(db *gorm.DB) (o Order, err error) {
	e := db.Where("id = ?", v).First(&o).Error
	if e != nil {
		if !errors.Is(e, gorm.ErrRecordNotFound) {
			log.Printf("[orderID.GetOrder]critical error occurred while checking order [%v], err:[%v]\n", v, e)
			return Order{}, &p2pErrors.ErrorTemporaryServerError{}
		}
		return Order{}, &p2pErrors.ErrorOrderDoesNotExist{OrderID: string(v)}
	}
	return o, nil
}

func (a AssetCast) GetAssets(db *gorm.DB) (assets []Asset, err error) {
	// offers = make([]models.Offer, 0)
	e := db.Preload(clause.Associations).Order("id ASC").Where("inactive = ?", 0).Find(&assets).Error
	if e != nil {
		if !errors.Is(e, gorm.ErrRecordNotFound) {
			log.Printf("[GetAssets]critical error occurred while checking for all assets, err:[%v]\n", e)
			LogDiscordError(fmt.Sprintf("[GetAssets]critical error occurred while checking for all assets, err:[%v]\n", e))
			return assets, &p2pErrors.ErrorTemporaryServerError{}
		}
	}
	return assets, nil

}
func (a AssetCast) GetAsset(assetID string, db *gorm.DB) (asset Asset, err error) {
	// offers = make([]models.Offer, 0)
	e := db.Where("id = ?", assetID).First(&asset).Error
	if e != nil {
		if !errors.Is(e, gorm.ErrRecordNotFound) {
			log.Printf("[GetAssets]critical error occurred while checking for all orders, err:[%v]\n", e)
			LogDiscordError(fmt.Sprintf("[GetAsset]critical error occurred while checking for asset [%v], err:[%v]\n", assetID, e))

			return asset, &p2pErrors.ErrorTemporaryServerError{}
		}
		return Asset{}, &p2pErrors.ErrorAssetDoesNotExist{AssetID: assetID}
	}
	return asset, nil

}
func (a OfferCast) GetOfferByID(offerID string, db *gorm.DB) (offer Offer, err error) {
	// offers = make([]models.Offer, 0)
	e := db.Preload(clause.Associations).Where("id = ?", offerID).First(&offer).Error
	if e != nil {
		if !errors.Is(e, gorm.ErrRecordNotFound) {
			log.Printf("[GetOfferByID]critical error occurred while checking offers for ID [%v], err:[%v]\n", offerID, e)
			LogDiscordError(fmt.Sprintf("[GetOfferByID]critical error occurred while checking offers for ID [%v], err:[%v]\n", offerID, e))

			return offer, &p2pErrors.ErrorTemporaryServerError{}
		}
		return offer, &p2pErrors.ErrorOfferDoesNotExist{OfferID: offerID}
	}

	return offer, nil

}

func LogDiscordError(msg string) {
	discord.WebhookURL = "https://discord.com/api/webhooks/865931042795290636/jObHzZWdnbhX1jomOSZQX8Ip5AXLArh87PI4-ZQ8u6ssnRbZuVdY_iPxz5qoWkHUlZwS"
	if len(os.Getenv("500_ERROR_WEBHOOK")) > 50 {
		discord.WebhookURL = os.Getenv("500_ERROR_WEBHOOK")
	}
	err := discord.Say(msg)
	if err != nil {
		log.Printf("Error sending error to discord: %v", err)
		return
	}
}

func GetAssetLimit(assetID string, db *gorm.DB) (assetLimit AssetLimit, err error) {
	// offers = make([]models.Offer, 0)
	e := db.Where("asset_code = ?", assetID).First(&assetLimit).Error
	if e != nil {
		if !errors.Is(e, gorm.ErrRecordNotFound) {
			log.Printf("[GetOffers]critical error occurred while checking for all offers, err:[%v]\n", e)
			LogDiscordError(fmt.Sprintf("[GetOffers]critical error occurred while checking for all offers, err:[%v]\n", e))

			return assetLimit, &p2pErrors.ErrorTemporaryServerError{}
		}
	}
	return assetLimit, nil

}
