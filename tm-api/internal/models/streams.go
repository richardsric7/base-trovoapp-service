package models

import (
	"admin-panel-dashboard/internal/cache"

	"admin-panel-dashboard/internal/trovosdk"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/ethclient"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

// GlobalConfig used to carry a large legacy P2P Telegram-trading-bot
// surface (offer/deal broadcasting, market-price caching, connection
// reports, message-deletion scheduling) - all of it was confirmed to have
// zero live callers anywhere in this service and depended on the now-
// deleted legacy P2P models (Offer/Currency/Asset/PaymentChannel/
// MarketPrice), so it was removed along with that schema. What remains
// here is the live websocket/push-notification broadcasting used by the
// real users/auth components.
type GlobalConfig struct {
	UsersOnline          map[string]UserMessageChannels // map of users that connect to socket and adds self and creates its first UserMessage channel with channel id. user must first check if it exists on the map first then retrieves the existing channels and adds new channel to it.
	GeneralUsers         map[string]chan map[string]interface{}
	AuthUsers            map[string]chan map[string]interface{}
	LoginUsers           map[string]chan map[string]interface{}
	ServiceLink          *trovosdk.ServiceLink
	Mutex                *sync.Mutex
	Cache                *cache.RedisCache
	TrovoWalletDB        *gorm.DB
	BlockchainClient     *ethclient.Client
	BlockchainPassphrase string
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
