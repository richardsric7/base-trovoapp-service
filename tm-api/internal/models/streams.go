package models

import (
	"admin-panel-dashboard/internal/cache"

	"admin-panel-dashboard/internal/trovosdk"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/gin-gonic/gin"
	tb "gopkg.in/tucnak/telebot.v2"
)

type GlobalConfig struct {
	UsersOnline                 map[string]UserMessageChannels // map of users that connect to socket and adds self and creates its first UserMessage channel with channel id. user must first check if it exists on the map first then retrieves the existing channels and adds new channel to it.
	GeneralUsers                map[string]chan map[string]interface{}
	AuthUsers                   map[string]chan map[string]interface{}
	TelegramNotificationChannel chan TGNotification
	TelegramBot                 *tb.Bot
	LoginUsers                  map[string]chan map[string]interface{}
	PaymentChannels             []PaymentChannel
	Currencies                  []Currency
	Assets                      []AssetJSON
	Markets                     []MarketPrice
	ServiceLink                 *trovosdk.ServiceLink
	TelegramDeleteList          TelegramDeleteList
	TelegramNextDelete          map[int64]map[string]StoredMessage
	Mutex                       *sync.Mutex
	Cache                       *cache.RedisCache
	TrovoWalletDB               *gorm.DB
	P2PDB                       *gorm.DB
	CallbackACL                 map[string]string
	BlockchainClient            *ethclient.Client
	BlockchainPassphrase        string
}

type UserMessageChannels map[string]chan map[string]interface{} // map of channel ID and the message channel that is used to share information to specific user connections

func (n *GlobalConfig) BroadcastToGeneralStream(streamType string, payload interface{}) {

	message := gin.H{"stream": payload, "streamType": streamType}

	n.Mutex.Lock()
	for _, conChan := range n.GeneralUsers {
		conChan <- message
	}
	n.Mutex.Unlock()

}

func (n *GlobalConfig) BroadcastToUserConnectionStreams(username, streamType string, payload interface{}) {

	message := gin.H{"stream": payload, "streamType": streamType}

	n.Mutex.Lock()
	for _, conChan := range n.UsersOnline[username] {
		conChan <- message
	}
	n.Mutex.Unlock()

}
func (n *GlobalConfig) BroadcastToLoginID(loginID, streamType string, payload interface{}) {

	message := gin.H{"stream": payload, "streamType": streamType}

	n.Mutex.Lock()
	if conChan, ok := n.LoginUsers[loginID]; ok {
		conChan <- message

	}

	n.Mutex.Unlock()

}

func (n *GlobalConfig) BroadcastToAuthID(authID, streamType string, payload interface{}) {

	message := gin.H{"stream": payload, "streamType": streamType}

	n.Mutex.Lock()
	if conChan, ok := n.AuthUsers[authID]; ok {
		conChan <- message

	}
	n.Mutex.Unlock()

}

func (n *GlobalConfig) SendPNAndSocketUserKYCChanged(userInfo UserInfo, userKYCBefore int) {

	pl := struct {
		User UserInfo `json:"user"`
	}{
		User: userInfo,
	}
	n.BroadcastToUserConnectionStreams(userInfo.Username, "userChanged", pl)

	if ((strings.Split(userInfo.Username, "@"))[1]) == "trovo" && (userKYCBefore == 0 && userInfo.KYCLevel == 1) {
		title := "TrovoP2P: p2p.trovotech.io Merchant access!"
		_, err := n.ServiceLink.SendPushNotification((strings.Split(userInfo.Username, "@"))[0], title, "You have been granted a merchant access on p2p.trovotech.io! Happy Sales!!!", "", "")
		if err != nil {
			log.Printf("Error sending push notification: %v\n", err)
			return
		}

	}

	if ((strings.Split(userInfo.Username, "@"))[1]) == "trovo" && (userKYCBefore == 1 && userInfo.KYCLevel == 0) {
		title := "TrovoP2P: p2p.trovotech.io Merchant access revoked!"
		_, err := n.ServiceLink.SendPushNotification(((strings.Split(userInfo.Username, "@"))[0]), title, "Your market merchant access on p2p.trovotech.io has been revoked!!!", "", "")
		if err != nil {
			log.Printf("Error sending push notification: %v\n", err)
			return
		}

	}
	{
		cacheKeyMaker := fmt.Sprintf("[GET] /v1/users/trades %v", userInfo.Username)
		cacheKeyTaker := fmt.Sprintf("[GET] /v1/users/orders %v", userInfo.Username)
		cacheKeyOfferMaker := fmt.Sprintf("[GET] /v1/users/offers %v", userInfo.Username)
		cacheKeyOfferTaker := fmt.Sprintf("[GET] /v1/users/offers %v", userInfo.Username)
		cacheKeyAllOffers := fmt.Sprintf("[GET] /v1/offers %v", "ALL")
		n.Cache.InvalidateCachedHttpResponse(cacheKeyMaker, cacheKeyTaker, cacheKeyAllOffers, cacheKeyOfferMaker, cacheKeyOfferTaker)

	}
}

func (n *GlobalConfig) SendPNAndSocketUserChanged(userInfo UserInfo) {

	pl := struct {
		User UserInfo `json:"user"`
	}{
		User: userInfo,
	}
	n.BroadcastToUserConnectionStreams(userInfo.Username, "userChanged", pl)

	{
		cacheKeyMaker := fmt.Sprintf("[GET] /v1/users/trades %v", userInfo.Username)
		cacheKeyTaker := fmt.Sprintf("[GET] /v1/users/orders %v", userInfo.Username)
		cacheKeyOfferMaker := fmt.Sprintf("[GET] /v1/users/offers %v", userInfo.Username)
		cacheKeyOfferTaker := fmt.Sprintf("[GET] /v1/users/offers %v", userInfo.Username)
		cacheKeyAllOffers := fmt.Sprintf("[GET] /v1/offers %v", "ALL")
		n.Cache.InvalidateCachedHttpResponse(cacheKeyMaker, cacheKeyTaker, cacheKeyAllOffers, cacheKeyOfferMaker, cacheKeyOfferTaker)

	}
}

func (n *GlobalConfig) SendSocketUserChanged(userInfo UserInfo) {

	pl := struct {
		User UserInfo `json:"user"`
	}{
		User: userInfo,
	}
	n.BroadcastToUserConnectionStreams(userInfo.Username, "userChanged", pl)
	{
		cacheKeyMaker := fmt.Sprintf("[GET] /v1/users/trades %v", userInfo.Username)
		cacheKeyTaker := fmt.Sprintf("[GET] /v1/users/orders %v", userInfo.Username)
		cacheKeyOfferMaker := fmt.Sprintf("[GET] /v1/users/offers %v", userInfo.Username)
		cacheKeyOfferTaker := fmt.Sprintf("[GET] /v1/users/offers %v", userInfo.Username)
		cacheKeyAllOffers := fmt.Sprintf("[GET] /v1/offers %v", "ALL")
		n.Cache.InvalidateCachedHttpResponse(cacheKeyMaker, cacheKeyTaker, cacheKeyAllOffers, cacheKeyOfferMaker, cacheKeyOfferTaker)

	}
}

func (n *GlobalConfig) ResetCache(owner string) {

	{
		cacheKeyMaker := fmt.Sprintf("[GET] /v1/users/trades %v", owner)
		cacheKeyTaker := fmt.Sprintf("[GET] /v1/users/orders %v", owner)
		cacheKeyOfferMaker := fmt.Sprintf("[GET] /v1/users/offers %v", owner)
		cacheKeyAllOffers := fmt.Sprintf("[GET] /v1/offers %v", "ALL")
		n.Cache.InvalidateCachedHttpResponse(cacheKeyMaker, cacheKeyTaker, cacheKeyAllOffers, cacheKeyOfferMaker)

	}
}

func (n *GlobalConfig) SendSocketOfferDeleted(offer Offer) {

	pl := struct {
		Offer interface{} `json:"offer"`
	}{
		Offer: offer.ToOfferJSON(n.TrovoWalletDB),
	}

	{
		// socket notification
		//offer maker
		n.BroadcastToUserConnectionStreams(offer.Maker, "offerChanged", pl)

		{
			// invalidate cache
			// cacheKey := fmt.Sprintf("[GET] /v1/users/verify/%v/%v", targetUser, loginID)
			cacheKey := fmt.Sprintf("[GET] /v1/offers %v", "ALL")
			cacheKeyUser := fmt.Sprintf("[GET] /v1/users/offers %v", offer.Maker)

			n.Cache.InvalidateCachedHttpResponse(cacheKey, cacheKeyUser)
		}

	}
}

func (n *GlobalConfig) SendSocketOfferChanged(offer Offer) {

	pl := struct {
		Offer interface{} `json:"offer"`
	}{
		Offer: offer.ToOfferJSON(n.TrovoWalletDB),
	}

	{
		// socket notification
		//offer maker
		n.BroadcastToUserConnectionStreams(offer.Maker, "offerChanged", pl)

		{
			// invalidate cache
			// cacheKey := fmt.Sprintf("[GET] /v1/users/verify/%v/%v", targetUser, loginID)
			cacheKey := fmt.Sprintf("[GET] /v1/offers %v", "ALL")
			cacheKeyUser := fmt.Sprintf("[GET] /v1/users/offers %v", offer.Maker)

			n.Cache.InvalidateCachedHttpResponse(cacheKey, cacheKeyUser)
		}

	}
}

func (n *GlobalConfig) SendSocketOfferCreated(offer Offer) {

	pl := struct {
		Offer interface{} `json:"offer"`
	}{
		Offer: offer.ToOfferJSON(n.TrovoWalletDB),
	}

	{
		// socket notification
		//offer maker
		n.BroadcastToUserConnectionStreams(offer.Maker, "offerCreated", pl)

		{
			// invalidate cache
			// cacheKey := fmt.Sprintf("[GET] /v1/users/verify/%v/%v", targetUser, loginID)
			cacheKey := fmt.Sprintf("[GET] /v1/offers %v", "ALL")
			cacheKeyUser := fmt.Sprintf("[GET] /v1/users/offers %v", offer.Maker)

			n.Cache.InvalidateCachedHttpResponse(cacheKey, cacheKeyUser)
		}

	}
}

func (n *GlobalConfig) SendPNAndSocketLowLiquidityNotice(offer Offer) {

	{
		// socket notification
		//offer maker
		n.SendSocketLowLiquidityNotice(offer)

		// send Push notification to taker that order has been accepted
		if ((strings.Split(offer.Maker, "@"))[1]) == "trovo" {

			title := "TrovoP2P: Low Liquidity alert!"
			message := fmt.Sprintf("Your %v %v offer for %v on %v has low liquidity. It will be removed from market list. Please increase your liquidity to restore your offer to the list", offer.CurrencyID, offer.OfferType, offer.AssetID, offer.CurrencyPaymentChannelID)

			_, err := n.ServiceLink.SendPushNotification(((strings.Split(offer.Maker, "@"))[0]), title, message, "", "")
			if err != nil {
				log.Printf("Error sending push notification: %v\n", err)
				return
			}

		}

	}
}

func (n *GlobalConfig) SendSocketLowLiquidityNotice(offer Offer) {

	pl := struct {
		Offer interface{} `json:"offer"`
	}{
		Offer: offer.ToOfferJSON(n.TrovoWalletDB),
	}

	{
		// socket notification
		//offer maker
		n.BroadcastToUserConnectionStreams(offer.Maker, "lowLiquidity", pl)

	}
	{
		// invalidate cache
		// cacheKey := fmt.Sprintf("[GET] /v1/users/verify/%v/%v", targetUser, loginID)
		cacheKey := fmt.Sprintf("[GET] /v1/offers %v", "ALL")
		cacheKeyUser := fmt.Sprintf("[GET] /v1/users/offers %v", offer.Maker)

		n.Cache.InvalidateCachedHttpResponse(cacheKey, cacheKeyUser)
	}
}

func (n *GlobalConfig) SendTGLowLiquidityNotice(makerRecord *User, offer Offer, menu *tb.ReplyMarkup) {
	if makerRecord == nil {
		return
	}

	// if makerRecord.Telegram != nil && makerRecord.TelegramNotifications == 1 {
	//
	//	message := fmt.Sprintf("Your %v %v offer for %v on %v has low liquidity. It will be removed from market list. Please increase your liquidity to restore your offer to the list", offer.CurrencyID, offer.OfferType, offer.AssetID, offer.CurrencyPaymentChannelID)
	//
	//	n.TelegramBot.Send(&tb.User{ID: *makerRecord.Telegram}, message, menu)
	//
	//}

}

func (n *GlobalConfig) GetUserPaymentMethods(userRecord *User) (paymentMethods []PaymentMethod) {
	paymentMethods = make([]PaymentMethod, 0)
	if userRecord == nil {
		return
	}
	n.TrovoWalletDB.Preload(clause.Associations).Order("payment_channel_id ASC").Where("username = ?", userRecord.Username).Find(&paymentMethods)
	// log.Printf("GetUserPaymentMethods for %v: %+v\n", userRecord.Username, paymentMethods)
	return paymentMethods
}

func (n *GlobalConfig) GetOfferCurrencyPaymentMethods(userRecord *User, offer OfferJSON) (paymentMethods []PaymentMethod) {
	paymentMethods = make([]PaymentMethod, 0)
	if userRecord == nil {
		return
	}
	n.TrovoWalletDB.Preload(clause.Associations).Order("payment_channel_id ASC").Where("username = ?", userRecord.Username).Where("payment_channel_id = ?", offer.CurrencyPaymentChannelID).Find(&paymentMethods)
	// log.Printf("GetUserPaymentMethods for %v: %+v\n", userRecord.Username, paymentMethods)
	return paymentMethods
}

func (n *GlobalConfig) GetUserPaymentMethod(username, destinationAccount string) (paymentMethod PaymentMethod, err error) {
	if destinationAccount == "" {
		return
	}
	err = n.TrovoWalletDB.Preload(clause.Associations).Where("username = ?", username).Where("destination_account LIKE ?", destinationAccount+"%").First(&paymentMethod).Error

	return paymentMethod, err
}

//
// func (n *GlobalConfig) GetFormattedUserDataForManageUserDisplay(userRecord *User) (formattedText string) {
//	if userRecord == nil {
//		return
//	}
//	formattedText = fmt.Sprintf("<b>Name:</b> %v %v\n", userRecord.FirstName, userRecord.LastName)
//
//	if userRecord.Mobile != nil {
//		formattedText = fmt.Sprintf("%v<b>Mobile:</b> %v\n", formattedText, *userRecord.Mobile)
//	}
//
//	if userRecord.ContactPhone != nil {
//		formattedText = fmt.Sprintf("%v<b>Contact Phone:</b> %v\n", formattedText, *userRecord.ContactPhone)
//	}
//
//	formattedText = fmt.Sprintf("%v<b>Email:</b> %v\n", formattedText, userRecord.Email)
//
//	if userRecord.CountryCode != nil {
//		formattedText = fmt.Sprintf("%v<b>Country:</b> %v\n", formattedText, *userRecord.CountryCode)
//	}
//
//	if userRecord.City != nil {
//		formattedText = fmt.Sprintf("%v<b>City:</b> %v\n", formattedText, *userRecord.City)
//	}
//
//	if userRecord.Suspended == 1 {
//
//		formattedText = fmt.Sprintf("%v<b>Suspended: Yes.</b>\n", formattedText)
//
//		if userRecord.SuspensionReason != nil {
//			formattedText = fmt.Sprintf("%v<b> Reason:</b> %v\n", formattedText, *userRecord.SuspensionReason)
//		}
//	}
//
//	if userRecord.MaxAssetPerOffer == 0 {
//
//		formattedText = fmt.Sprintf("%v<b>Max Asset Per Offer:</b> %v\n", formattedText, os.Getenv("MAX_XBN_PER_OFFER"))
//	} else {
//		formattedText = fmt.Sprintf("%v<b>Max Asset Per Offer:</b> %v\n", formattedText, userRecord.MaxAssetPerOffer)
//	}
//	formattedText = fmt.Sprintf("%v<b>Public Key:</b> %v\n\n", formattedText, userRecord.PublicKey)
//
//	//get blockchain account balance
//	balances, err := bc.GetAccountBalance(userRecord.PublicKey)
//	if err == nil {
//		formattedText = fmt.Sprintf("%v<b>Blockchain Balances:</b> <a href='%v/account/%v'>%v</a>\n", formattedText, os.Getenv("EXPLORER_URL"), userRecord.PublicKey, userRecord.PublicKey)
//		for _, b := range balances {
//			if b.AssetCode == "" {
//				formattedText = fmt.Sprintf("%v<b><i>%v:</i></b> %v\n", formattedText, os.Getenv("NATIVE_ASSET_CODE"), b.Amount)
//			} else {
//				formattedText = fmt.Sprintf("%v<b><i>%v:</i></b> %v\n", formattedText, b.AssetCode, b.Amount)
//			}
//		}
//	}
//	countActive, countInactive, makerOffers := n.GetMakerOffers(userRecord)
//	if len(makerOffers) > 0 {
//		formattedText = fmt.Sprintf("%v<b>Offers: Active: %v, Disabled: %v </b>\n", formattedText, countActive, countInactive)
//
//		if len(makerOffers) > 10 {
//			//trim the offer list
//			makerOffers = makerOffers[0:9]
//		}
//		for i, o := range makerOffers {
//
//			formattedText = fmt.Sprintf("%v<b> ↠%v.: </b> %v\n", formattedText, i+1, o.Formatted)
//
//		}
//	}
//
//	return
//}

func (n *GlobalConfig) GetUserBantuWallets(userRecord *User) (paymentMethods []PaymentMethod) {
	paymentMethods = make([]PaymentMethod, 0)
	if userRecord == nil {
		return
	}
	n.TrovoWalletDB.Raw("SELECT * FROM payment_methods WHERE payment_channel_id = 'Bantu Public Key' AND username = ?", userRecord.Username).Find(&paymentMethods)

	return paymentMethods
}

type TGAdminOfferDisplay struct {
	OfferType      string
	Asset          string
	Liquidity      string
	Currency       string
	Price          string
	PaymentChannel string // Currency payment channel
	Formatted      string
}

func (n *GlobalConfig) GetMakerOffers(userRecord *User) (countActive, countInactive int, offersDiplay []TGAdminOfferDisplay) {
	offersDiplay = make([]TGAdminOfferDisplay, 0)
	if userRecord == nil {
		return
	}
	// if userRecord.KYCLevel == 0 {
	//	return
	//}
	var offers []Offer
	e := n.TrovoWalletDB.Where("maker = ?", userRecord.Username).Find(&offers).Error
	if e == nil {
		for _, offer := range offers {
			if offer.Inactive == 0 {
				countActive++
			} else {
				countInactive++
			}
			o := TGAdminOfferDisplay{
				OfferType:      offer.OfferType,
				Asset:          offer.AssetID,
				Liquidity:      fmt.Sprintf("%v", decimal.NewFromFloat(offer.AssetAmount).String()),
				Currency:       offer.CurrencyID,
				Price:          decimal.NewFromFloat(offer.AssetPrice).String(),
				PaymentChannel: offer.CurrencyPaymentChannelID,
				Formatted:      fmt.Sprintf("%ving %v %v @ %v %v through %v", offer.OfferType, decimal.NewFromFloat(offer.AssetAmount).String(), offer.AssetID, decimal.NewFromFloat(offer.AssetPrice).String(), offer.CurrencyID, offer.CurrencyPaymentChannelID),
			}
			offersDiplay = append(offersDiplay, o)
		}
	}
	return
}

func (n *GlobalConfig) GetPaymentChannels() (paymentChannels []PaymentChannel) {
	paymentChannels = make([]PaymentChannel, 0)

	n.TrovoWalletDB.Order("id ASC").Where("inactive=?", 0).Find(&paymentChannels)
	return paymentChannels
}

func (n *GlobalConfig) GetCurrencies() (currencies []Currency) {
	currencies = make([]Currency, 0)
	n.TrovoWalletDB.Order("id ASC").Where("inactive=?", 0).Find(&currencies)

	// log.Printf("CurrencyList: %+v", currencies)
	return currencies
}

func (n *GlobalConfig) GetAssets() (assets []Asset) {
	assets = make([]Asset, 0)
	n.TrovoWalletDB.Preload(clause.Associations).Order("id ASC").Where("inactive = ?", 0).Find(&assets)

	// log.Printf("AssetList: %+v\n", assets)
	return assets
}

func (n *GlobalConfig) RemoveNextDeleteMsgs(tgUser *tb.User) {
	if tgUser == nil {
		log.Println("[RemoveNextDeleteMsgs] tguser is nil")
		return
	}
	n.Mutex.Lock()
	list, ok := n.TelegramNextDelete[tgUser.ID]
	n.Mutex.Unlock()
	if !ok {
		n.Mutex.Lock()
		n.TelegramNextDelete[tgUser.ID] = make(map[string]StoredMessage)
		n.Mutex.Unlock()
		return
	}

	for _, m := range list {
		err := n.TelegramBot.Delete(&m)
		if err != nil {
			log.Printf("[RemoveNextDeleteMsgs] error deleting message: %v\n", err)
			continue
		}
	}

	n.Mutex.Lock()
	delete(n.TelegramNextDelete, tgUser.ID)
	n.Mutex.Unlock()
}

func (n *GlobalConfig) AddNextDeleteMsg(tgUser *tb.User, msg *tb.Message) {
	if tgUser == nil {
		log.Println("[AddNextDeleteMsg] tguser is nil")

		return
	}
	if msg == nil {
		log.Println("[AddNextDeleteMsg] message is nil")
		return
	}
	n.Mutex.Lock()
	_, ok := n.TelegramNextDelete[tgUser.ID]
	n.Mutex.Unlock()
	if !ok {
		n.Mutex.Lock()
		n.TelegramNextDelete[tgUser.ID] = make(map[string]StoredMessage)
		n.Mutex.Unlock()
	}

	n.Mutex.Lock()
	n.TelegramNextDelete[tgUser.ID][uuid.NewString()] = StoredMessage{
		MessageID: fmt.Sprintf("%v", msg.ID),
		ChatID:    msg.Chat.ID,
	}
	n.Mutex.Unlock()
}

// OrderFormInterface
type OrderFormObject interface {
	GetFormattedOrder() string
}
type OrderForm struct {
	OrderType             string
	Asset                 string
	CurrencyID            string
	PaymentMethodID       string
	PaymentChannelID      string
	ModeOfEntry           string
	Amount                string
	TrustedMerchantOption string
	TrustedMerchant       string
}

func (orderForm OrderForm) GetFormattedOrder() (formattedOrder string) {
	if strings.EqualFold(orderForm.ModeOfEntry, "currency") {
		if !strings.EqualFold(orderForm.TrustedMerchant, "none") {
			// has a trusted merchant
			formattedOrder = fmt.Sprintf("Your order is a request to %v %v %v worth of %v through %v including offers from the merchant %v", orderForm.OrderType, orderForm.Amount, orderForm.CurrencyID, orderForm.Asset, orderForm.PaymentChannelID, orderForm.TrustedMerchant)

		} else {
			// has no trusted merchant
			formattedOrder = fmt.Sprintf("Your order is a request to %v %v %v worth of %v through %v.", orderForm.OrderType, orderForm.Amount, orderForm.CurrencyID, orderForm.Asset, orderForm.PaymentChannelID)

		}

	}
	if strings.EqualFold(orderForm.ModeOfEntry, "asset") {
		if !strings.EqualFold(orderForm.TrustedMerchant, "none") {
			// has a trusted merchant
			formattedOrder = fmt.Sprintf("Your order is a request to %v %v %v for %v through %v including offers from the merchant %v", orderForm.OrderType, orderForm.Amount, orderForm.Asset, orderForm.CurrencyID, orderForm.PaymentChannelID, orderForm.TrustedMerchant)

		} else {
			// has no trusted merchant
			formattedOrder = fmt.Sprintf("Your order is a request to %v %v %v for %v through %v.", orderForm.OrderType, orderForm.Amount, orderForm.Asset, orderForm.CurrencyID, orderForm.PaymentChannelID)

		}

	}
	return
}

type Deal struct {
	OfferID                  string
	Maker                    string
	OfferType                string
	AssetID                  string
	CurrencyID               string
	CurrencyPaymentChannelID string
	AssetPrice               string
	Completed                string
	Trades                   string
	Rank                     string
	ReportsAgainst           string
	FormattedText            string
}

func (n *GlobalConfig) OfferToDeal(o Offer, orderForm OrderForm, i int) (deal Deal, oj OfferJSON) {
	oj = o.ToOfferJSON(n.TrovoWalletDB)
	log.Printf("@@@@@@ @@@@ @@@ offer [%v], %+v\n", i+1, oj)
	dealText := ""
	if orderForm.ModeOfEntry == "currency" {
		if strings.EqualFold(oj.OrderType, "Sell") {
			dealText = fmt.Sprintf("↠ %v. %v %ving @ %v%v (%v%v) | Trades:%v | Rank:%v﹪", i+1, oj.Maker, oj.OfferType, oj.AssetPrice, oj.CurrencyID, (decimal.RequireFromString(orderForm.Amount).Div(decimal.RequireFromString(oj.AssetPrice)).Truncate(3)).String(), oj.AssetID, oj.MakerStats.Trades, oj.MakerStats.Rank)
		} else {
			dealText = fmt.Sprintf("↠ %v. %v %ving @ %v%v. (%v%v) | Trades:%v | Rank:%v﹪", i+1, oj.Maker, oj.OfferType, oj.AssetPrice, oj.CurrencyID, (decimal.RequireFromString(orderForm.Amount).Div(decimal.RequireFromString(oj.AssetPrice)).Truncate(3)).String(), oj.AssetID, oj.MakerStats.Trades, oj.MakerStats.Rank)

		}

	} else {
		// asset amount
		if strings.EqualFold(oj.OrderType, "Sell") {
			dealText = fmt.Sprintf("↠ %v. %v %ving @ %v%v. (%v%v) | Trades:%v | Rank:%v﹪", i+1, oj.Maker, oj.OfferType, oj.AssetPrice, oj.CurrencyID, (decimal.RequireFromString(orderForm.Amount).Mul(decimal.RequireFromString(oj.AssetPrice)).RoundBank(2)).String(), oj.CurrencyID, oj.MakerStats.Trades, oj.MakerStats.Rank)
		} else {
			dealText = fmt.Sprintf("↠ %v. %v %ving @ %v%v. (%v%v) | Trades:%v | Rank:%v﹪", i+1, oj.Maker, oj.OfferType, oj.AssetPrice, oj.CurrencyID, (decimal.RequireFromString(orderForm.Amount).Mul(decimal.RequireFromString(oj.AssetPrice)).RoundBank(2)).String(), oj.CurrencyID, oj.MakerStats.Trades, oj.MakerStats.Rank)

		}
		// dealText = fmt.Sprintf("↠ %v. %v %ving @ %v%v. You get %v%v | Trades:%v | Rank:%v﹪", i+1, oj.Maker, oj.OfferType, oj.AssetPrice, oj.CurrencyID, (decimal.RequireFromString(orderForm.Amount).Mul(decimal.RequireFromString(oj.AssetPrice)).RoundBank(2)).String(), oj.CurrencyID, oj.MakerStats.Trades, oj.MakerStats.Rank)

	}
	deal = Deal{
		OfferID:                  oj.ID,
		Maker:                    oj.Maker,
		OfferType:                oj.OfferType,
		AssetID:                  oj.AssetID,
		CurrencyID:               oj.CurrencyID,
		CurrencyPaymentChannelID: oj.CurrencyPaymentChannelID,
		AssetPrice:               oj.AssetPrice,
		Completed:                oj.MakerStats.Completed,
		Trades:                   oj.MakerStats.Trades,
		Rank:                     oj.MakerStats.Rank,
		ReportsAgainst:           oj.MakerStats.ReportsAgainst,
		FormattedText:            dealText,
	}

	return
}

// func (n *GlobalConfig) GetPublicKeyBalance(publicKey string) (balanceText string) {
//	balances, err := bc.GetAccountBalance(publicKey)
//	if err == nil {
//		balanceText = fmt.Sprintf("%v<b>Blockchain Balances:</b> <a href='%v/account/%v'>%v</a>\n", balanceText, os.Getenv("EXPLORER_URL"), publicKey, publicKey)
//		for _, b := range balances {
//			if b.AssetCode == "" {
//				balanceText = fmt.Sprintf("%v<b><i>%v:</i></b> %v\n", balanceText, "XBN", b.Amount)
//			} else {
//				balanceText = fmt.Sprintf("%v<b><i>%v:</i></b> %v\n", balanceText, b.AssetCode, b.Amount)
//			}
//		}
//	}
//
//	return
//}
//func (n *GlobalConfig) ValidateAssetTxID(txID string, order Order) (ops operations.Operation, err error) {
//	assetCode := ""
//	if order.OfferAssetID != "XBN" {
//		assetCode = order.OfferAssetID
//	}
//	if len(txID) < 50 {
//		err = fmt.Errorf("%s is not a valid bantu blockchain transaction ID", txID)
//		return ops, err
//	}
//	ops, err = bc.GetPaymentForTransactionID(order.OrderEscrowAddress, txID, assetCode)
//	if err != nil {
//		return
//	}
//	if ops.GetType() == "payment" {
//		pmt := interface{}(ops).(operations.Payment)
//		if decimal.RequireFromString(pmt.Amount).LessThan(decimal.NewFromFloat(order.OrderAmount)) {
//			//amount is invalid
//			return ops, &p2pErrors.ErrorInvalidTransaction{}
//		}
//
//	}
//
//	if ops.GetType() == "create_account" {
//		pmt := interface{}(ops).(operations.CreateAccount)
//		if decimal.RequireFromString(pmt.StartingBalance).LessThan(decimal.NewFromFloat(order.OrderAmount)) {
//			//amount is invalid
//			return ops, &p2pErrors.ErrorInvalidTransaction{}
//		}
//	}
//
//	return ops, nil
//
//}

type OfferSearch struct {
	ID, Maker, Page, CurrencyID, AssetID, OfferType, OrderType, PaymentChannelID, CurrencyPaymentChannelID, Limit string
}

type OrderSearch struct {
	ID, Maker, Taker, Page, CurrencyID, AssetID, OfferType, OrderType, ResultGroup, Limit, ForWho string
}

func (n *GlobalConfig) GetMarketDeals(orderForm OrderForm) (dealsText string, offerList []OfferJSON, deals []Deal) {
	offerList = make([]OfferJSON, 0)
	deals = make([]Deal, 0)
	var offers = make([]Offer, 0)

	var err error
	// var proposedAssetAmount, proposedCurrencyAmount float64
	var queryStr string
	if orderForm.ModeOfEntry == "currency" {
		// search by poposedCurrencyAmount on min/max trade

		if strings.EqualFold(orderForm.OrderType, "sell") {
			// order it by highest price at the top
			orderAmount, _ := decimal.RequireFromString(orderForm.Amount).Float64()
			queryStr = `select * from offers where (order_type = ? AND currency_id = ? AND asset_id = ?
				AND (? between min_trade_amount AND max_trade_amount) AND
				offline = 0 AND currency_payment_channel_id = ?
				) OR (order_type = ? AND currency_id = ? AND asset_id = ?
				AND (? between min_trade_amount AND max_trade_amount) AND
				offline = 0 AND currency_payment_channel_id = ? AND maker = ?) ORDER BY asset_price DESC LIMIT 10`
			err = n.TrovoWalletDB.Raw(queryStr, strings.ToTitle(orderForm.OrderType), orderForm.CurrencyID, orderForm.Asset, orderAmount, orderForm.PaymentChannelID,
				strings.ToTitle(orderForm.OrderType), orderForm.CurrencyID, orderForm.Asset, orderAmount, orderForm.PaymentChannelID, orderForm.TrustedMerchant).Find(&offers).Error
		} else {
			// order is buy. order it by highest price at the top
			orderAmount, _ := decimal.RequireFromString(orderForm.Amount).Float64()
			queryStr = `select * from offers where (order_type = ? AND currency_id = ? AND asset_id = ?
				AND (? between min_trade_amount AND max_trade_amount) AND
				offline = 0 AND currency_payment_channel_id = ?) OR (order_type = ? AND currency_id = ? AND asset_id = ?
				AND (? between min_trade_amount AND max_trade_amount) AND
				offline = 0 AND currency_payment_channel_id = ? AND maker = ?) ORDER BY asset_price ASC LIMIT 10`
			err = n.TrovoWalletDB.Raw(queryStr, strings.ToTitle(orderForm.OrderType), orderForm.CurrencyID, orderForm.Asset, orderAmount, orderForm.PaymentChannelID,
				strings.ToTitle(orderForm.OrderType), orderForm.CurrencyID, orderForm.Asset, orderAmount, orderForm.PaymentChannelID, orderForm.TrustedMerchant).Find(&offers).Error

		}

	} else {
		// search by poposedAssetAmount
		if strings.EqualFold(orderForm.OrderType, "sell") {
			// order it by highest price at the top
			orderAmount, _ := decimal.RequireFromString(orderForm.Amount).Float64()
			queryStr = `select * from offers where (order_type = ? AND currency_id = ? AND asset_id = ?
				AND (round((asset_price*?)::numeric,2) between min_trade_amount AND max_trade_amount) AND
				offline = 0 AND currency_payment_channel_id = ?) OR (order_type = ? AND currency_id = ? AND asset_id = ?
				AND (round((asset_price*?)::numeric,2) between min_trade_amount AND max_trade_amount) AND
				offline = 0 AND currency_payment_channel_id = ? AND maker = ?) ORDER BY asset_price DESC LIMIT 10`
			err = n.TrovoWalletDB.Raw(queryStr, strings.ToTitle(orderForm.OrderType), orderForm.CurrencyID, orderForm.Asset, orderAmount, orderForm.PaymentChannelID,
				strings.ToTitle(orderForm.OrderType), orderForm.CurrencyID, orderForm.Asset, orderAmount, orderForm.PaymentChannelID, orderForm.TrustedMerchant).Find(&offers).Error
		} else {
			// order  buyit by lowest price at the top
			orderAmount, _ := decimal.RequireFromString(orderForm.Amount).Float64()
			queryStr = `select * from offers where (order_type = ? AND currency_id = ? AND asset_id = ?
				AND (round((asset_price*?)::numeric,2) between min_trade_amount AND max_trade_amount) AND
				offline = 0 AND currency_payment_channel_id = ?) OR (order_type = ? AND currency_id = ? AND asset_id = ?
				AND (round((asset_price*?)::numeric,2) between min_trade_amount AND max_trade_amount) AND
				offline = 0 AND currency_payment_channel_id = ? AND maker = ?) ORDER BY asset_price ASC LIMIT 10`
			err = n.TrovoWalletDB.Raw(queryStr, strings.ToTitle(orderForm.OrderType), orderForm.CurrencyID, orderForm.Asset, orderAmount, orderForm.PaymentChannelID,
				strings.ToTitle(orderForm.OrderType), orderForm.CurrencyID, orderForm.Asset, orderAmount, orderForm.PaymentChannelID, orderForm.TrustedMerchant).Find(&offers).Error

		}
	}
	if err != nil {
		log.Println("[GetMarketOffers]@@@@@@@@@@@ @@@@@@@@@@@@ error fetching records from db", err)
	}
	if len(offers) > 0 {
		log.Printf("@@@@@@@@@@@@ @@@@@@@@@ @@@@@ [%v] offers retrieved.\n", len(offers))
		for i, o := range offers {
			deal, oj := n.OfferToDeal(o, orderForm, i)
			deals = append(deals, deal)
			offerList = append(offerList, oj)
			dealText := deal.FormattedText
			dealsText = fmt.Sprintf("%v%v\n", dealsText, dealText)
		}
	} else {
		log.Println("[GetMarketOffers]@@@@@@@@@@@ @@@@@@@@@@@@ NO DEALS FETCHED", err)

	}

	log.Printf("@@@@@@ DEALS TEXT: [%v]\n", dealsText)
	return dealsText, offerList, deals
}

func (n *GlobalConfig) GetOfferList(offerSearch OfferSearch) (records PaginatedOffers) {
	var query *gorm.DB
	var err error
	var countQuery *gorm.DB
	var ol []Offer
	offers := make([]OfferJSON, 0)
	// oD := ""

	paymentChannelID := offerSearch.PaymentChannelID
	currencyPaymentChannelID := offerSearch.CurrencyPaymentChannelID
	limit := 5
	if len(offerSearch.Limit) > 0 {
		limitU, _ := strconv.ParseUint(offerSearch.Limit, 10, 64)
		limit = int(limitU)
	}

	pageVal := "1"

	if len(offerSearch.Page) > 0 {
		pageVal = offerSearch.Page
	}

	pageU, _ := strconv.ParseUint(pageVal, 10, 64)
	page := int(pageU)
	currencyID := offerSearch.CurrencyID
	assetID := offerSearch.AssetID
	offerType := offerSearch.OfferType
	orderType := offerSearch.OrderType
	// orderBy := "asset_price"
	// orderDirection := c.DefaultQuery("order", "DESC")
	query = n.TrovoWalletDB.Preload(clause.Associations)
	countQuery = n.TrovoWalletDB.Group("id")

	// if len(orderDirection) > 0 && strings.ToLower(orderDirection) == "desc" {
	// 	oD = "DESC"
	// }
	if strings.EqualFold(orderType, "Buy") || strings.EqualFold(offerType, "Sell") {
		query = query.Order("asset_price ASC")
		countQuery = countQuery.Order("id ASC")

	}
	if strings.EqualFold(orderType, "Sell") || strings.EqualFold(offerType, "Buy") {
		query = query.Order("asset_price DESC")
		countQuery = countQuery.Order("id DESC")

	}

	if len(offerSearch.Maker) > 0 {
		// request is for all offers since owner is supplied
		query = query.Where("maker = ?", offerSearch.Maker)
		countQuery = countQuery.Where("maker = ?", offerSearch.Maker)
	}

	if len(currencyID) > 0 {
		query = query.Where("currency_id = ?", currencyID)
		countQuery = countQuery.Where("currency_id = ?", currencyID)
	}
	if len(offerSearch.ID) > 0 {
		query = query.Where("id = ?", offerSearch.ID)
		countQuery = countQuery.Where("id = ?", offerSearch.ID)
	}
	if len(paymentChannelID) > 0 {
		query = query.Where("payment_channel_id = ?", paymentChannelID)
		countQuery = countQuery.Where("payment_channel_id = ?", paymentChannelID)
	}
	if len(currencyPaymentChannelID) > 0 {
		query = query.Where("currency_payment_channel_id = ?", currencyPaymentChannelID)
		countQuery = countQuery.Where("currency_payment_channel_id = ?", currencyPaymentChannelID)
	}
	if len(assetID) > 0 {
		query = query.Where("asset_id = ?", assetID)
		countQuery = countQuery.Where("asset_id = ?", assetID)
	}
	if len(offerType) > 0 {
		query = query.Where("offer_type = ?", offerType)
		countQuery = countQuery.Where("offer_type = ?", offerType)
	}
	if len(orderType) > 0 {
		query = query.Where("order_type = ?", orderType)
		countQuery = countQuery.Where("order_type = ?", orderType)
	}
	if len(offerSearch.Maker) == 0 {
		// if owner is not supplied, then request is for all active offers
		query = query.Where("inactive = ?", 0)
		// offline
		query = query.Where("offline = ?", 0)
		query = query.Where("asset_amount > ?", 0)
		query = query.Where("max_trade_amount > ? AND max_trade_amount > min_trade_amount", 0)

		countQuery = countQuery.Where("inactive = ?", 0)
		// offline
		countQuery = countQuery.Where("offline = ?", 0)
		countQuery = countQuery.Where("asset_amount > ?", 0)
		countQuery = countQuery.Where("max_trade_amount > ? AND max_trade_amount > min_trade_amount", 0)

	}

	var countR int64

	errCount := countQuery.Find(&[]Offer{}).Count(&countR).Error
	if errCount != nil {
		log.Println("[GetOfferList]Count Error:", errCount)
		return records
	}

	if limit > 0 {
		query.Limit(limit)
	}
	count := int(countR)
	pages := 1
	if count > limit {
		log.Println("count / limit = ", count/limit, "count%limit = ", count%limit)
		pages = count / limit
		if count%limit > 0 {
			pages = pages + 1
		}
	}
	if page > pages {
		page = pages
	}
	if page > 1 {
		log.Println("Offset = ", (page-1)*limit)
		query.Offset(((page - 1) * limit))
	}
	// query = handy.Tif(page > 1, query.Offset(int(page*limit)), query).(*gorm.DB)
	if err = query.Find(&ol).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Println("[GetOfferList]Query Error:", err)
		return
	}
	for i, v := range ol {
		serialNo := i
		if page > 1 {
			serialNo = ((page - 1) * limit) + i
		}
		o := v.ToOfferJSON(n.TrovoWalletDB)
		o.FormattedText = o.ToFormattedText(serialNo)
		offers = append(offers, o)
	}

	// pages = (handy.Tif(count > int64(limit), int((count / int64(limit))), 1)).(int) + handy.Tif(count%int64(limit) > 0, 1, 0).(int)
	log.Println("Number of Pages:", pages)
	log.Println("Current Page:", page)
	log.Println("Total Records Found:", count)
	log.Println("Limit Used:", limit)
	records = PaginatedOffers{CurrentPage: page, Pages: pages, TotalRecords: count, Limit: limit, Records: offers}

	return records
}

type XbnDollarPrice struct {
	ID            string `gorm:"primaryKey"`
	Asset         string `gorm:"index:isx_dollar_price_unique_asset,unique"`
	Source        string `gorm:"size:100;index:isx_dollar_price_unique_asset,unique"`
	AskRate       string `gorm:"size:100"`
	BidRate       string `gorm:"size:100"`
	LastTradeRate string `gorm:"size:100"`
	LastUpdated   time.Time
}

// GetXBNDollarAskPrice dollar ask price using USDB
func (n *GlobalConfig) GetBittrexXBNDollarAskPrice() (xbnPrice XbnDollarPrice, err error) {

	err = n.TrovoWalletDB.Where("source = ?", "BITTREX").Where("asset = ?", "XBN").First(&xbnPrice).Error
	if err != nil {

		return XbnDollarPrice{
			ID:            "",
			Asset:         "XBN",
			Source:        "Bittrex",
			AskRate:       "0",
			BidRate:       "0",
			LastTradeRate: "0",
			LastUpdated:   time.Now(),
		}, err
	}
	return xbnPrice, nil

}

type OrderReport struct {
	OrderType   string
	Asset       string
	TotalTrades uint64
	TotalAmount float64
}
type TradeReport struct {
	OfferType   string
	Asset       string
	TotalTrades uint64
	TotalAmount float64
}

// DailyOrdersReport
func (n *GlobalConfig) DailyOrdersReport(taker string) (report string, err error) {

	// var orderReports []OrderReport
	orderReports := make([]OrderReport, 0)
	err = n.TrovoWalletDB.Raw(`select order_type, offer_asset_id as asset, count(0) as total_trades, sum(order_amount) as total_amount from orders
	Where order_status_id = 11 and offer_taker = ? and updated_at::date = now()::date
   group by order_type, offer_asset_id`, taker).Scan(&orderReports).Error
	if err != nil {
		return "", err
	}
	if len(orderReports) == 0 {
		return "", errors.New("no record")
	}
	xbnDolarPrice, _ := n.GetBittrexXBNDollarAskPrice()
	for i, r := range orderReports {
		rformat := ""
		if strings.EqualFold(r.Asset, "XBN") {
			// give dollar equivalent for XBN
			bidWorth := (decimal.RequireFromString(xbnDolarPrice.BidRate).Mul(decimal.NewFromFloat(r.TotalAmount))).Round(2).String()
			askWorth := (decimal.RequireFromString(xbnDolarPrice.AskRate).Mul(decimal.NewFromFloat(r.TotalAmount))).Round(2).String()
			rformat = fmt.Sprintf("%v. %v %v | Total Trades: %v | Total Amount:  %v(B:$%v[%v]/A:$%v[%v])", i+1, r.OrderType, r.Asset, r.TotalTrades, r.TotalAmount, bidWorth, xbnDolarPrice.BidRate, askWorth, xbnDolarPrice.AskRate)

		} else {
			// do not show for none XBN
			rformat = fmt.Sprintf("%v. %v %v | Total Trades: %v | Total Amount: %v", i+1, r.OrderType, r.Asset, r.TotalTrades, r.TotalAmount)

		}
		report = fmt.Sprintf("%v%v\n", report, rformat)

	}
	return report, nil

}

// DailyOrdersReport
func (n *GlobalConfig) DailyTradesReport(maker string) (report string, err error) {

	var tradeReports []TradeReport
	err = n.TrovoWalletDB.Raw(`select offer_type, offer_asset_id as asset, count(0) as total_trades, sum(order_amount) as total_amount from orders
	Where order_status_id = 11 and offer_maker = ? and updated_at::date = now()::date
   group by offer_type, offer_asset_id`, maker).Scan(&tradeReports).Error
	if err != nil {
		return "", err
	}
	if len(tradeReports) == 0 {
		return "", errors.New("no record")
	}
	xbnDolarPrice, _ := n.GetBittrexXBNDollarAskPrice()
	for i, r := range tradeReports {
		rformat := ""
		if strings.EqualFold(r.Asset, "XBN") {
			// give dollar equivalent for XBN
			bidWorth := (decimal.RequireFromString(xbnDolarPrice.BidRate).Mul(decimal.NewFromFloat(r.TotalAmount))).Round(2).String()
			askWorth := (decimal.RequireFromString(xbnDolarPrice.AskRate).Mul(decimal.NewFromFloat(r.TotalAmount))).Round(2).String()
			rformat = fmt.Sprintf("%v. %v %v | Total Trades: %v | Total Amount:  %v(B:$%v[%v]/A:$%v[%v])", i+1, r.OfferType, r.Asset, r.TotalTrades, r.TotalAmount, bidWorth, xbnDolarPrice.BidRate, askWorth, xbnDolarPrice.AskRate)

		} else {
			// do not show for none XBN
			rformat = fmt.Sprintf("%v. %v %v | Total Trades: %v | Total Amount: %v", i+1, r.OfferType, r.Asset, r.TotalTrades, r.TotalAmount)

		}
		report = fmt.Sprintf("%v%v\n", report, rformat)

	}
	return report, nil

}

// YesterDayOrdersReport
func (n *GlobalConfig) YesterdayOrdersReport(taker string) (report string, err error) {

	var orderReports []OrderReport
	err = n.TrovoWalletDB.Raw(`select order_type, offer_asset_id as asset, count(0) as total_trades, sum(order_amount) as total_amount from orders
	Where order_status_id = 11 and offer_taker = ? and updated_at::date = current_date - interval '1 days'
   group by order_type, offer_asset_id`, taker).Scan(&orderReports).Error
	if err != nil {
		return "", err
	}
	if len(orderReports) == 0 {
		return "", errors.New("no record")
	}
	xbnDolarPrice, _ := n.GetBittrexXBNDollarAskPrice()
	for i, r := range orderReports {
		rformat := ""
		if strings.EqualFold(r.Asset, "XBN") {
			// give dollar equivalent for XBN
			bidWorth := (decimal.RequireFromString(xbnDolarPrice.BidRate).Mul(decimal.NewFromFloat(r.TotalAmount))).Round(2).String()
			askWorth := (decimal.RequireFromString(xbnDolarPrice.AskRate).Mul(decimal.NewFromFloat(r.TotalAmount))).Round(2).String()
			rformat = fmt.Sprintf("%v. %v %v | Total Trades: %v | Total Amount:  %v(B:$%v[%v]/A:$%v[%v])", i+1, r.OrderType, r.Asset, r.TotalTrades, r.TotalAmount, bidWorth, xbnDolarPrice.BidRate, askWorth, xbnDolarPrice.AskRate)

		} else {
			// do not show for none XBN
			rformat = fmt.Sprintf("%v. %v %v | Total Trades: %v | Total Amount: %v", i+1, r.OrderType, r.Asset, r.TotalTrades, r.TotalAmount)

		}
		report = fmt.Sprintf("%v%v\n", report, rformat)

	}
	return report, nil

}

// SevenDayOrdersReport
func (n *GlobalConfig) SevenDayOrdersReport(taker string) (report string, err error) {

	var orderReports []OrderReport
	err = n.TrovoWalletDB.Raw(`select order_type, offer_asset_id as asset, count(0) as total_trades, sum(order_amount) as total_amount from orders
	Where order_status_id = 11 and offer_taker = ? and updated_at::date > current_date - interval '7 days'
   group by order_type, offer_asset_id`, taker).Scan(&orderReports).Error
	if err != nil {
		return "", err
	}
	if len(orderReports) == 0 {
		return "", errors.New("no record")
	}
	xbnDolarPrice, _ := n.GetBittrexXBNDollarAskPrice()
	for i, r := range orderReports {
		rformat := ""
		if strings.EqualFold(r.Asset, "XBN") {
			// give dollar equivalent for XBN
			bidWorth := (decimal.RequireFromString(xbnDolarPrice.BidRate).Mul(decimal.NewFromFloat(r.TotalAmount))).Round(2).String()
			askWorth := (decimal.RequireFromString(xbnDolarPrice.AskRate).Mul(decimal.NewFromFloat(r.TotalAmount))).Round(2).String()
			rformat = fmt.Sprintf("%v. %v %v | Total Trades: %v | Total Amount:  %v(B:$%v[%v]/A:$%v[%v])", i+1, r.OrderType, r.Asset, r.TotalTrades, r.TotalAmount, bidWorth, xbnDolarPrice.BidRate, askWorth, xbnDolarPrice.AskRate)

		} else {
			// do not show for none XBN
			rformat = fmt.Sprintf("%v. %v %v | Total Trades: %v | Total Amount: %v", i+1, r.OrderType, r.Asset, r.TotalTrades, r.TotalAmount)

		}
		report = fmt.Sprintf("%v%v\n", report, rformat)

	}
	return report, nil

}

// YesterdayTradesReport
func (n *GlobalConfig) YesterdayTradesReport(maker string) (report string, err error) {

	var tradeReports []TradeReport
	err = n.TrovoWalletDB.Raw(`select offer_type, offer_asset_id as asset, count(0) as total_trades, sum(order_amount) as total_amount from orders
	Where order_status_id = 11 and offer_maker = ? and updated_at::date = current_date - interval '1 days'
   group by offer_type, offer_asset_id`, maker).Scan(&tradeReports).Error
	if err != nil {
		return "", err
	}
	if len(tradeReports) == 0 {
		return "", errors.New("no record")
	}
	xbnDolarPrice, _ := n.GetBittrexXBNDollarAskPrice()
	for i, r := range tradeReports {
		rformat := ""
		if strings.EqualFold(r.Asset, "XBN") {
			// give dollar equivalent for XBN
			bidWorth := (decimal.RequireFromString(xbnDolarPrice.BidRate).Mul(decimal.NewFromFloat(r.TotalAmount))).Round(2).String()
			askWorth := (decimal.RequireFromString(xbnDolarPrice.AskRate).Mul(decimal.NewFromFloat(r.TotalAmount))).Round(2).String()
			rformat = fmt.Sprintf("%v. %v %v | Total Trades: %v | Total Amount:  %v(B:$%v[%v]/A:$%v[%v])", i+1, r.OfferType, r.Asset, r.TotalTrades, r.TotalAmount, bidWorth, xbnDolarPrice.BidRate, askWorth, xbnDolarPrice.AskRate)

		} else {
			// do not show for none XBN
			rformat = fmt.Sprintf("%v. %v %v | Total Trades: %v | Total Amount: %v", i+1, r.OfferType, r.Asset, r.TotalTrades, r.TotalAmount)

		}
		report = fmt.Sprintf("%v%v\n", report, rformat)

	}
	return report, nil

}

// SevenDayTradesReport
func (n *GlobalConfig) SevenDayTradesReport(maker string) (report string, err error) {

	var tradeReports []TradeReport
	err = n.TrovoWalletDB.Raw(`select offer_type, offer_asset_id as asset, count(0) as total_trades, sum(order_amount) as total_amount from orders
	Where order_status_id = 11 and offer_maker = ? and updated_at::date > current_date - interval '7 days'
   group by offer_type, offer_asset_id`, maker).Scan(&tradeReports).Error
	if err != nil {
		return "", err
	}
	if len(tradeReports) == 0 {
		return "", errors.New("no record")
	}
	xbnDolarPrice, _ := n.GetBittrexXBNDollarAskPrice()
	for i, r := range tradeReports {
		rformat := ""
		if strings.EqualFold(r.Asset, "XBN") {
			// give dollar equivalent for XBN
			bidWorth := (decimal.RequireFromString(xbnDolarPrice.BidRate).Mul(decimal.NewFromFloat(r.TotalAmount))).Round(2).String()
			askWorth := (decimal.RequireFromString(xbnDolarPrice.AskRate).Mul(decimal.NewFromFloat(r.TotalAmount))).Round(2).String()
			rformat = fmt.Sprintf("%v. %v %v | Total Trades: %v | Total Amount:  %v(B:$%v[%v]/A:$%v[%v])", i+1, r.OfferType, r.Asset, r.TotalTrades, r.TotalAmount, bidWorth, xbnDolarPrice.BidRate, askWorth, xbnDolarPrice.AskRate)

		} else {
			// do not show for none XBN
			rformat = fmt.Sprintf("%v. %v %v | Total Trades: %v | Total Amount: %v", i+1, r.OfferType, r.Asset, r.TotalTrades, r.TotalAmount)

		}
		report = fmt.Sprintf("%v%v\n", report, rformat)

	}
	return report, nil

}

// ThirtyDayOrdersReport
func (n *GlobalConfig) ThirtyDayOrdersReport(taker string) (report string, err error) {

	var orderReports []OrderReport
	err = n.TrovoWalletDB.Raw(`select order_type, offer_asset_id as asset, count(0) as total_trades, sum(order_amount) as total_amount from orders
	Where order_status_id = 11 and offer_taker = ? and updated_at::date > current_date - interval '30 days'
   group by order_type, offer_asset_id`, taker).Scan(&orderReports).Error
	if err != nil {
		return "", err
	}
	if len(orderReports) == 0 {
		return "", errors.New("no record")
	}
	xbnDolarPrice, _ := n.GetBittrexXBNDollarAskPrice()
	for i, r := range orderReports {
		rformat := ""
		if strings.EqualFold(r.Asset, "XBN") {
			// give dollar equivalent for XBN
			bidWorth := (decimal.RequireFromString(xbnDolarPrice.BidRate).Mul(decimal.NewFromFloat(r.TotalAmount))).Round(2).String()
			askWorth := (decimal.RequireFromString(xbnDolarPrice.AskRate).Mul(decimal.NewFromFloat(r.TotalAmount))).Round(2).String()
			rformat = fmt.Sprintf("%v. %v %v | Total Trades: %v | Total Amount:  %v(B:$%v[%v]/A:$%v[%v])", i+1, r.OrderType, r.Asset, r.TotalTrades, r.TotalAmount, bidWorth, xbnDolarPrice.BidRate, askWorth, xbnDolarPrice.AskRate)

		} else {
			// do not show for none XBN
			rformat = fmt.Sprintf("%v. %v %v | Total Trades: %v | Total Amount: %v", i+1, r.OrderType, r.Asset, r.TotalTrades, r.TotalAmount)

		}
		report = fmt.Sprintf("%v%v\n", report, rformat)

	}
	return report, nil

}

// ThirtyDayOrdersReport
func (n *GlobalConfig) ThirtyDayTradesReport(maker string) (report string, err error) {

	var tradeReports []TradeReport
	err = n.TrovoWalletDB.Raw(`select offer_type, offer_asset_id as asset, count(0) as total_trades, sum(order_amount) as total_amount from orders
	Where order_status_id = 11 and offer_maker = ? and updated_at::date > current_date - interval '30 days'
   group by offer_type, offer_asset_id`, maker).Scan(&tradeReports).Error
	if err != nil {
		return "", err
	}
	if len(tradeReports) == 0 {
		return "", errors.New("no record")
	}
	xbnDolarPrice, _ := n.GetBittrexXBNDollarAskPrice()
	for i, r := range tradeReports {
		rformat := ""
		if strings.EqualFold(r.Asset, "XBN") {
			// give dollar equivalent for XBN
			bidWorth := (decimal.RequireFromString(xbnDolarPrice.BidRate).Mul(decimal.NewFromFloat(r.TotalAmount))).Round(2).String()
			askWorth := (decimal.RequireFromString(xbnDolarPrice.AskRate).Mul(decimal.NewFromFloat(r.TotalAmount))).Round(2).String()
			rformat = fmt.Sprintf("%v. %v %v | Total Trades: %v | Total Amount:  %v(B:$%v[%v]/A:$%v[%v])", i+1, r.OfferType, r.Asset, r.TotalTrades, r.TotalAmount, bidWorth, xbnDolarPrice.BidRate, askWorth, xbnDolarPrice.AskRate)

		} else {
			// do not show for none XBN
			rformat = fmt.Sprintf("%v. %v %v | Total Trades: %v | Total Amount: %v", i+1, r.OfferType, r.Asset, r.TotalTrades, r.TotalAmount)

		}
		report = fmt.Sprintf("%v%v\n", report, rformat)

	}
	return report, nil

}

// AllTimeOrdersReport
func (n *GlobalConfig) AllTimeOrdersReport(taker string) (report string, err error) {

	var orderReports []OrderReport
	err = n.TrovoWalletDB.Raw(`select order_type, offer_asset_id as asset, count(0) as total_trades, sum(order_amount) as total_amount from orders
	Where order_status_id = 11 and offer_taker = ?
   group by order_type, offer_asset_id`, taker).Scan(&orderReports).Error
	if err != nil {
		return "", err
	}
	if len(orderReports) == 0 {
		return "", errors.New("no record")
	}
	xbnDolarPrice, _ := n.GetBittrexXBNDollarAskPrice()
	for i, r := range orderReports {
		rformat := ""
		if strings.EqualFold(r.Asset, "XBN") {
			// give dollar equivalent for XBN
			bidWorth := (decimal.RequireFromString(xbnDolarPrice.BidRate).Mul(decimal.NewFromFloat(r.TotalAmount))).Round(2).String()
			askWorth := (decimal.RequireFromString(xbnDolarPrice.AskRate).Mul(decimal.NewFromFloat(r.TotalAmount))).Round(2).String()
			rformat = fmt.Sprintf("%v. %v %v | Total Trades: %v | Total Amount:  %v(B:$%v[%v]/A:$%v[%v])", i+1, r.OrderType, r.Asset, r.TotalTrades, r.TotalAmount, bidWorth, xbnDolarPrice.BidRate, askWorth, xbnDolarPrice.AskRate)

		} else {
			// do not show for none XBN
			rformat = fmt.Sprintf("%v. %v %v | Total Trades: %v | Total Amount: %v", i+1, r.OrderType, r.Asset, r.TotalTrades, r.TotalAmount)

		}
		report = fmt.Sprintf("%v%v\n", report, rformat)

	}
	return report, nil

}

// AllTimeTradesReport
func (n *GlobalConfig) AllTimeTradesReport(maker string) (report string, err error) {

	var tradeReports []TradeReport
	err = n.TrovoWalletDB.Raw(`select offer_type, offer_asset_id as asset, count(0) as total_trades, sum(order_amount) as total_amount from orders
	Where order_status_id = 11 and offer_maker = ?
   group by offer_type, offer_asset_id`, maker).Scan(&tradeReports).Error
	if err != nil {
		return "", err
	}
	if len(tradeReports) == 0 {
		return "", errors.New("no record")
	}
	xbnDolarPrice, _ := n.GetBittrexXBNDollarAskPrice()
	for i, r := range tradeReports {
		rformat := ""
		if strings.EqualFold(r.Asset, "XBN") {
			// give dollar equivalent for XBN
			bidWorth := (decimal.RequireFromString(xbnDolarPrice.BidRate).Mul(decimal.NewFromFloat(r.TotalAmount))).Round(2).String()
			askWorth := (decimal.RequireFromString(xbnDolarPrice.AskRate).Mul(decimal.NewFromFloat(r.TotalAmount))).Round(2).String()
			rformat = fmt.Sprintf("%v. %v %v | Total Trades: %v | Total Amount:  %v(B:$%v[%v]/A:$%v[%v])", i+1, r.OfferType, r.Asset, r.TotalTrades, r.TotalAmount, bidWorth, xbnDolarPrice.BidRate, askWorth, xbnDolarPrice.AskRate)

		} else {
			// do not show for none XBN
			rformat = fmt.Sprintf("%v. %v %v | Total Trades: %v | Total Amount: %v", i+1, r.OfferType, r.Asset, r.TotalTrades, r.TotalAmount)

		}
		report = fmt.Sprintf("%v%v\n", report, rformat)

	}
	return report, nil

}

// AllTimeTradesReport
func (n *GlobalConfig) TelegramAllTimeConnectionReport() (total int) {

	err := n.TrovoWalletDB.Raw(`select count(0) as total from users where telegram is not null`).Scan(&total).Error
	if err != nil {
		return 0
	}

	return total

}

// AllTimeTradesReport
func (n *GlobalConfig) TelegramDailyConnectionReport() (total int) {

	err := n.TrovoWalletDB.Raw(`select count(0) as total from users where telegram is not null AND telegram_connected_at::date = current_date`).Scan(&total).Error
	if err != nil {
		return 0
	}

	return total

}

type TGConnections struct {
	DateConnected string
	Total         int
}

type TGMonthConnections struct {
	Week  int
	Total int
}
type TGYearConnections struct {
	Month int
	Total int
}

// Telegram7DConnnectionReport
func (n *GlobalConfig) Telegram7DConnectionReport() (total int) {
	err := n.TrovoWalletDB.Raw(`select count(0) as total from users where telegram is not null AND telegram_connected_at::date > current_date - interval '7 days'`).Scan(&total).Error
	if err != nil {
		return total
	}

	return total

}

// Telegram7DConnnectionReport
func (n *GlobalConfig) Telegram7DConnectionReportBD() (report []TGConnections) {
	report = make([]TGConnections, 0)
	err := n.TrovoWalletDB.Raw(`select telegram_connected_at::date as date_connected, count(0) as total from users where telegram is not null AND telegram_connected_at::date > current_date - interval '7 days'
	group by telegram_connected_at::date`).Scan(&report).Error
	if err != nil {
		return report
	}

	return report

}

// Telegram7DConnnectionReport
func (n *GlobalConfig) TelegramCurrentMonthConnectionReport() (total int) {
	err := n.TrovoWalletDB.Raw(`select count(0) as total from users where telegram is not null AND date_part('year', telegram_connected_at) = date_part('year', now())
	AND date_part('month', telegram_connected_at) = date_part('month', now())`).Scan(&total).Error
	if err != nil {
		return total
	}

	return total

}

// Telegram7DConnnectionReport
func (n *GlobalConfig) TelegramCurrentYearConnectionReport() (total int) {
	err := n.TrovoWalletDB.Raw(`select count(0) as total from users where telegram is not null
	AND date_part('year', telegram_connected_at) = date_part('year', now())`).Scan(&total).Error
	if err != nil {
		return total
	}

	return total

}

// Telegram7DConnnectionReport
func (n *GlobalConfig) TelegramCurrentMonthConnectionReportBD() (report []TGMonthConnections) {
	report = make([]TGMonthConnections, 0)
	err := n.TrovoWalletDB.Raw(`select date_part('week', telegram_connected_at) as week, count(0) as total from users where telegram is not null AND date_part('year', telegram_connected_at) = date_part('year', now())
	AND date_part('month', telegram_connected_at) = date_part('month', now())
	group by date_part('week', telegram_connected_at)`).Scan(&report).Error
	if err != nil {
		return report
	}

	return report

}

// Telegram7DConnnectionReport
func (n *GlobalConfig) TelegramCurrentYearConnectionReportBD() (report []TGYearConnections) {
	report = make([]TGYearConnections, 0)
	err := n.TrovoWalletDB.Raw(`select date_part('month', telegram_connected_at) as month, count(0) as total from users where telegram is not null AND date_part('year', telegram_connected_at) = date_part('year', now())
	group by date_part('month', telegram_connected_at)`).Scan(&report).Error
	if err != nil {
		return report
	}

	return report

}
