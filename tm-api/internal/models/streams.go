package models

import (
	"admin-panel-dashboard/internal/cache"

	"admin-panel-dashboard/internal/trovosdk"
	"context"
	"encoding/json"
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
// loginStreamChannelPrefix/authStreamChannelPrefix namespace these two
// one-shot notification streams' Redis Pub/Sub channels from anything else
// that might use the same Redis instance (e.g. app-backend's cache, which
// shares REDIS_URL in some deployments).
const loginStreamChannelPrefix = "tm-api:login-stream:"
const authStreamChannelPrefix = "tm-api:auth-stream:"

func LoginStreamChannel(loginID string) string { return loginStreamChannelPrefix + loginID }
func AuthStreamChannel(authID string) string   { return authStreamChannelPrefix + authID }

func (n *GlobalConfig) BroadcastToLoginID(loginID, streamType string, payload interface{}) {

	message := gin.H{"stream": payload, "streamType": streamType}

	n.Mutex.Lock()
	if conChan, ok := n.LoginUsers[loginID]; ok {
		conChan <- message

	}

	n.Mutex.Unlock()

	// The local map send above only reaches a subscriber connected to
	// *this* instance - a login callback (POST /callbacks/login/...) can
	// land on a different instance than the one holding the admin's SSE
	// connection (GET /login/stream/...) once this service runs 2+
	// replicas behind a load balancer. Publishing lets every instance's
	// SubscribeStreamRelay (see LoginNotificationStream) pick it up
	// regardless of which one received the callback.
	n.publishStreamEvent(LoginStreamChannel(loginID), message)
}

func (n *GlobalConfig) BroadcastToAuthID(authID, streamType string, payload interface{}) {

	message := gin.H{"stream": payload, "streamType": streamType}

	n.Mutex.Lock()
	if conChan, ok := n.AuthUsers[authID]; ok {
		conChan <- message

	}
	n.Mutex.Unlock()

	n.publishStreamEvent(AuthStreamChannel(authID), message)
}

// publishStreamEvent fans a one-shot stream message out over Redis so a
// subscriber on another instance can receive it too. Fire-and-forget:
// Redis PUBLISH is at-most-once and silently no-ops with zero subscribers,
// which matches this stream's existing behavior when nobody's connected -
// no regression from the pre-Redis local-map-only version, since that one
// already silently dropped the message whenever the connection was on a
// different process (previously always a bug, now only when Redis itself
// is unavailable/disabled).
func (n *GlobalConfig) publishStreamEvent(channel string, message gin.H) {
	if n.Cache == nil || !n.Cache.Enabled || n.Cache.Client == nil {
		return
	}
	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("[publishStreamEvent] failed to marshal message for channel [%v]: %v\n", channel, err)
		return
	}
	if err := n.Cache.Client.Publish(n.Cache.Context, channel, data).Err(); err != nil {
		log.Printf("[publishStreamEvent] failed to publish to channel [%v]: %v\n", channel, err)
	}
}

// SubscribeStreamRelay relays messages published to a Redis channel into a
// local delivery channel (ch), non-blocking so it never leaks a goroutine
// waiting on a receiver that already stopped reading - both of this
// model's streams are one-shot (the caller's select loop reads at most one
// message from ch and returns), so once the local map send or an earlier
// relay delivers, nobody reads ch again and a second delivery must be
// droppable, not block forever. Callers must invoke the returned cancel
// func (e.g. via defer) once they stop reading from ch, or this leaks a
// Redis subscription for up to ctx's lifetime.
func (n *GlobalConfig) SubscribeStreamRelay(ctx context.Context, channel string, ch chan map[string]interface{}) func() {
	if n.Cache == nil || !n.Cache.Enabled || n.Cache.Client == nil {
		return func() {}
	}
	pubsub := n.Cache.Client.Subscribe(n.Cache.Context, channel)
	relayCtx, cancel := context.WithCancel(ctx)
	go func() {
		defer pubsub.Close()
		msgCh := pubsub.Channel()
		for {
			select {
			case <-relayCtx.Done():
				return
			case msg, ok := <-msgCh:
				if !ok {
					return
				}
				var payload map[string]interface{}
				if err := json.Unmarshal([]byte(msg.Payload), &payload); err != nil {
					log.Printf("[SubscribeStreamRelay] failed to unmarshal message from channel [%v]: %v\n", channel, err)
					continue
				}
				select {
				case ch <- payload:
				default:
					// ch already has an undelivered message buffered (or
					// nobody's reading anymore) - drop rather than block,
					// this is a best-effort duplicate delivery path.
				}
			}
		}
	}()
	return cancel
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
