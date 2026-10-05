package users

import (
	"encoding/json"
	"os"
	"time"
	announcementServices "trovo-wallet-api/internal/components/announcements/services"
	blockchain "trovo-wallet-api/internal/components/assets/blockchain"
	paymentModels "trovo-wallet-api/internal/components/payments/models"
	usersDB "trovo-wallet-api/internal/components/users/db"
	userModels "trovo-wallet-api/internal/components/users/models"
	"trovo-wallet-api/internal/middleware"
	"trovo-wallet-api/internal/sharedconfig"

	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// StartMessage is models for stream start message
type StartMessage struct {
	Stream string `json:"stream"`
}

type StreamObject struct {
	Address string
	Alias   string
	request interface{}
}

// UserWebSocketAPI handles websocket connections
func UserWebSocketAPI(c *gin.Context, gc *sharedconfig.GlobalConfig) {
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("error get WS connection")
		return
	}
	defer ws.Close()
	messageChan := make(chan map[string]interface{}, 20)

	type actionData struct {
		ID              string `json:"id"`
		TransactionHash string `json:"transactionHash"`
		TxType          string `json:"txType"`
	}
	//Read data in ws

	var data userModels.Streams
	var auth userModels.Handshake
	err = ws.ReadJSON(&data)
	defer ws.Close()
	if err != nil {
		log.Println("error read json")
		ws.Close()
		return
	}
	log.Printf("received subscriptionMessage: %+v\n", data)
	authEnable := os.Getenv("ENABLE_WEBSOCKET_AUTH") == "1"
	if authEnable {
		if data.Signer == "null" {
			log.Printf("signer cannot be %v\n", data.Signer)
			ws.Close()
			return
		}
		errWebsocketSignature := middleware.WebSocketAuthenticationChecks(data.Stream+data.Cursor, data.Signature, data.Signer)
		if errWebsocketSignature != nil {
			log.Println("*************####****** Websocket signature failure:", err)
			auth.Auth = false
			auth.Message = "signature could not be verified"
			message := gin.H{"stream": auth, "streamType": "auth"}
			ws.WriteJSON(message)
			return
		}
		log.Println("<<<<>>>>*************####****** Websocket signature verified")
		auth.Auth = true
		auth.Message = "signature verified"
		message := gin.H{"stream": auth, "streamType": "auth"}
		ws.WriteJSON(message)

	}
	data.Stream = strings.ToLower(data.Stream)
	identifier := strings.TrimSpace(strings.ToLower(c.Param("identifier")))

	user, err := usersDB.GetUser(identifier, gc.DB, gc)
	if err != nil {
		auth.Auth = false
		auth.Message = "user could not be authenticated"
		message := gin.H{"stream": auth, "streamType": "auth"}
		ws.WriteJSON(message)
		ws.Close()
		return
	}
	if user.Suspended == 1 {
		auth.Auth = false
		auth.Message = "user account has been suspended"
		message := gin.H{"stream": auth, "streamType": "auth"}
		ws.WriteJSON(message)
		ws.Close()
		return

	}

	if !user.SignerIsValid(data.Signer, gc) && authEnable {
		auth.Auth = false
		auth.Message = "signer mismatch"
		message := gin.H{"stream": auth, "streamType": "auth"}
		ws.WriteJSON(message)
		return
	}

	auth.Auth = true
	auth.Message = "success"
	message := gin.H{"stream": auth, "streamType": "auth"}
	ws.WriteJSON(message)

	// Live P2P order/offer/escrow/dispute updates for this user, fanned
	// out across instances via Redis (see sharedconfig/realtime.go and
	// p2p/services/notifications.go's NotifyUsername) - fire-and-forget,
	// a live nudge alongside the existing push notification, not a
	// catch-up-able feed.
	p2pStreamChan, unregisterP2PStream := sharedconfig.RegisterUserStreamConnection(gc, identifier)
	defer unregisterP2PStream()

	ip := c.ClientIP()
	if len(c.GetHeader("Cf-Connecting-Ip")) > 4 {
		ip = c.GetHeader("Cf-Connecting-Ip")

	}
	announcements, err := announcementServices.HandleGetAnnouncement(ip, gc.DB)
	if err == nil {
		message := gin.H{"stream": announcements, "streamType": "announcements"}
		ws.WriteJSON(message)

	}

	// Live payments: the user's wallets' new payment history (recorded by
	// payment-history-engine), pushed as it appears.
	done := make(chan struct{})
	defer close(done)
	go streamUserPayments(done, &user, data.Cursor, messageChan, gc)
	go keepAlive(done, messageChan)

	for {
		var v map[string]interface{}
		select {
		case v = <-messageChan:
		case v = <-p2pStreamChan:
		}
		if err = ws.WriteJSON(v); err != nil {
			log.Printf("Error Sending stream: %v\nError %v\n", v, err)
			return
		}
	}
}

// send hands msg to a connection's writer unless the connection is done.
func send(done <-chan struct{}, out chan<- map[string]interface{}, msg map[string]interface{}) bool {
	select {
	case out <- msg:
		return true
	case <-done:
		return false
	}
}

// keepAlive sends a keep-alive every 30s until done.
func keepAlive(done <-chan struct{}, out chan<- map[string]interface{}) {
	t := time.NewTicker(30 * time.Second)
	defer t.Stop()
	for {
		select {
		case <-done:
			return
		case <-t.C:
			if !send(done, out, gin.H{"stream": "keep-alive", "streamType": "keep-alive"}) {
				return
			}
		}
	}
}

// streamUserPayments pushes the user's wallets' payments recorded after
// cursor (RFC 3339 or unix seconds; default: now) every few seconds, as
// "payment" messages ("swap" for swaps), until done.
func streamUserPayments(done <-chan struct{}, user *userModels.User, cursor string, out chan<- map[string]interface{}, gc *sharedconfig.GlobalConfig) {
	since := time.Now().UTC()
	if t, err := time.Parse(time.RFC3339, strings.TrimSpace(cursor)); err == nil {
		since = t
	} else if n, err := strconv.ParseInt(strings.TrimSpace(cursor), 10, 64); err == nil && n > 0 {
		since = time.Unix(n, 0).UTC()
	}
	seen := map[string]bool{}
	t := time.NewTicker(5 * time.Second)
	defer t.Stop()
	var addresses []string
	for tick := 0; ; tick++ {
		select {
		case <-done:
			return
		case <-t.C:
		}
		// the wallet list changes rarely: re-read it once a minute
		if tick%12 == 0 {
			addresses = addresses[:0]
			for _, w := range user.GetAllWallets(gc) {
				addresses = append(addresses, strings.ToUpper(w.ID))
			}
		}
		if len(addresses) == 0 {
			continue
		}
		var rows []paymentModels.PaymentHistory
		err := gc.DB.Where("(UPPER(from_address) IN ? OR UPPER(to_address) IN ?) AND transaction_date >= ?", addresses, addresses, since).
			Order("transaction_date ASC").Limit(100).Find(&rows).Error
		if err != nil {
			continue
		}
		for _, r := range rows {
			key := r.TransactionID + "|" + r.PT + "|" + r.SourceAccountSequence + "|" + r.TransactionType
			if seen[key] {
				continue
			}
			seen[key] = true
			kind := "payment"
			if strings.EqualFold(r.TransactionType, "SWAP") {
				kind = "swap"
			}
			if !send(done, out, gin.H{"stream": r, "streamType": kind}) {
				return
			}
			if r.TransactionDate.After(since) {
				since = r.TransactionDate
			}
		}
		if len(seen) > 5000 {
			seen = map[string]bool{}
		}
	}
}

// OrderBookSocketAPI streams an order book (offers on the offer book) as it
// changes.
func OrderBookSocketAPI(c *gin.Context, gc *sharedconfig.GlobalConfig) {
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("error get WS connection")
		return
	}
	defer ws.Close()

	var data userModels.OrderBookStream
	if err = ws.ReadJSON(&data); err != nil {
		log.Println("error read json")
		return
	}
	log.Printf("received subscriptionMessage: %+v\n", data)
	auth := userModels.Handshake{Auth: true, Message: "success"}
	if err = ws.WriteJSON(gin.H{"stream": auth, "streamType": "auth"}); err != nil {
		return
	}

	messageChan := make(chan map[string]interface{}, 20)
	done := make(chan struct{})
	defer close(done)
	go keepAlive(done, messageChan)
	go func() {
		if !send(done, messageChan, gin.H{"stream": "Started Orderbook Stream", "streamType": "notice"}) {
			return
		}
		var last string
		t := time.NewTicker(5 * time.Second)
		defer t.Stop()
		for {
			book, err := blockchain.GetOrderBook(data.AssetCode, data.ContractAddress, data.CurrencyCode, data.CurrencyIssuer)
			if err == nil {
				if b, _ := json.Marshal(book); string(b) != last {
					last = string(b)
					if !send(done, messageChan, gin.H{"stream": book, "streamType": "orderBook"}) {
						return
					}
				}
			}
			select {
			case <-done:
				return
			case <-t.C:
			}
		}
	}()

	for {
		v := <-messageChan
		if err = ws.WriteJSON(v); err != nil {
			log.Printf("Error Sending stream: %v\nError %v\n", v, err)
			return
		}
	}
}

// TradeChartSocketAPI streams a pair's trade chart (fills on the offer
// book) as it changes.
func TradeChartSocketAPI(c *gin.Context, gc *sharedconfig.GlobalConfig) {
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("error get WS connection")
		return
	}
	defer ws.Close()

	var data userModels.ChartStream
	if err = ws.ReadJSON(&data); err != nil {
		log.Println("error read json")
		return
	}
	log.Printf("[TradeChartSocketAPI]received subscriptionMessage: %+v\n", data)
	auth := userModels.Handshake{Auth: true, Message: "success"}
	if err = ws.WriteJSON(gin.H{"stream": auth, "streamType": "auth"}); err != nil {
		return
	}

	messageChan := make(chan map[string]interface{}, 20)
	done := make(chan struct{})
	defer close(done)
	go keepAlive(done, messageChan)
	go func() {
		if !send(done, messageChan, gin.H{"stream": "Started Trade Chart Stream", "streamType": "notice"}) {
			return
		}
		var last string
		t := time.NewTicker(30 * time.Second)
		defer t.Stop()
		for {
			tradeChart, err := blockchain.GetChartRecords(data.AssetCode, data.ContractAddress, data.CurrencyCode, data.CurrencyIssuer, "", "", data.Order, "", data.ChartPeriod, "")
			if err == nil {
				if b, _ := json.Marshal(tradeChart); string(b) != last {
					last = string(b)
					if !send(done, messageChan, gin.H{"stream": tradeChart, "streamType": "tradeChart"}) {
						return
					}
				}
			}
			select {
			case <-done:
				return
			case <-t.C:
			}
		}
	}()

	for {
		v := <-messageChan
		if err = ws.WriteJSON(v); err != nil {
			log.Printf("Error Sending stream: %v\nError %v\n", v, err)
			return
		}
	}
}
