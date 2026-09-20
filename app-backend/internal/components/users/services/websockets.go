package users

import (
	"encoding/json"
	"os"
	"time"
	announcementServices "trovo-wallet-api/internal/components/announcements/services"
	blockchain "trovo-wallet-api/internal/components/assets/blockchain"
	usersDB "trovo-wallet-api/internal/components/users/db"
	userModels "trovo-wallet-api/internal/components/users/models"
	"trovo-wallet-api/internal/middleware"
	"trovo-wallet-api/internal/sharedconfig"

	"context"
	"fmt"
	"log"
	"net/http"
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
	_, cancel := context.WithCancel(context.Background())
	defer cancel()
	messageChan := make(chan map[string]interface{}, 200)
	if err != nil {
		log.Println("error get WS connection")
	}
	defer ws.Close()

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

	if !user.SignerIsValid(data.Signer, false, gc) && authEnable {
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
	ip := c.ClientIP()
	if len(c.GetHeader("Cf-Connecting-Ip")) > 4 {
		ip = c.GetHeader("Cf-Connecting-Ip")

	}
	announcements, err := announcementServices.HandleGetAnnouncement(ip, gc.DB)
	if err == nil {
		message := gin.H{"stream": announcements, "streamType": "announcements"}
		ws.WriteJSON(message)

	}

	// Live operation/effect streaming (payment/swap/trustline notifications
	// pushed over this socket) used Horizon's SSE streams
	// (client.StreamOperations/StreamEffects), keyed by a Stellar
	// account's ledger cursor. Base has no equivalent account-keyed event
	// stream - the Base-native way to get this is subscribing to
	// eth_subscribe("logs") filtered by the user's address/curated B20
	// token contracts (go-ethereum's ethclient.SubscribeFilterLogs),
	// which needs a WS-capable RPC endpoint and a real event-decoding
	// design, tracked as a follow-up out of scope for this alteration
	// pass. This connection still authenticates and stays open (clients
	// get the auth/notice/keep-alive messages below), it just does not
	// push live blockchain events yet.
	message = gin.H{"stream": "Live event streaming is not yet available on Base for this connection.", "streamType": "notice"}
	messageChan <- message

	//try to read from ws and exit if cannot read.
	go func() {
		log.Println("@@@@@@started ws connection keepalive......")
		for {
			time.Sleep(30 * time.Second)
			auth.Auth = true
			auth.Message = "keep-alive"
			message := gin.H{"stream": "keep-alive", "streamType": "keep-alive"}
			messageChan <- message

		}

	}()

	for {
		v := <-messageChan
		err = ws.WriteJSON(v)
		if err != nil {
			log.Printf("Error Sending stream: %v\nError %v\n", v, err)
			ws.Close()
			cancel()
			break
		}
	}
	// wg.Wait()

}

// OrderBookSocketAPI handles websocket connections
func OrderBookSocketAPI(c *gin.Context, gc *sharedconfig.GlobalConfig) {
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)

	_, cancel := context.WithCancel(context.Background())
	defer cancel()
	messageChan := make(chan map[string]interface{}, 200)
	if err != nil {
		log.Println("error get WS connection")
	}
	defer ws.Close()

	var data userModels.OrderBookStream
	var auth userModels.Handshake
	err = ws.ReadJSON(&data)
	defer ws.Close()
	if err != nil {
		log.Println("error read json")
		ws.Close()
		return
	}
	log.Printf("received subscriptionMessage: %+v\n", data)

	auth.Auth = true
	auth.Message = "success"
	message := gin.H{"stream": auth, "streamType": "auth"}
	ws.WriteJSON(message)

	// Order-book streaming used Horizon's client.StreamOrderBooks against
	// Stellar's native DEX order book. Base has no native on-chain order
	// book to stream from - a real implementation needs a DEX/AMM router
	// integration (e.g. indexing swap events from a specific router
	// contract), tracked as a follow-up out of scope for this alteration
	// pass. This connection still authenticates and stays open.
	message = gin.H{"stream": "Order book streaming is not yet available on Base for this connection.", "streamType": "notice"}
	messageChan <- message

	//try to read from ws and exit if cannot read.
	go func() {
		log.Println("@@@@@@started ws connection keepalive......")
		for {
			time.Sleep(30 * time.Second)
			auth.Auth = true
			auth.Message = "keep-alive"
			message := gin.H{"stream": "keep-alive", "streamType": "keep-alive"}
			messageChan <- message

		}

	}()

	for {
		v := <-messageChan
		err = ws.WriteJSON(v)
		if err != nil {
			log.Printf("Error Sending stream: %v\nError %v\n", v, err)
			ws.Close()
			cancel()
			break
		}
	}

}

// TradeChartSocketAPI handles websocket connections
func TradeChartSocketAPI(c *gin.Context, gc *sharedconfig.GlobalConfig) {
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)

	// _, cancel := context.WithCancel(context.Background())
	// defer cancel()
	messageChan := make(chan map[string]interface{}, 200)
	if err != nil {
		log.Println("error get WS connection")
	}
	defer ws.Close()

	var data userModels.ChartStream
	var auth userModels.Handshake
	err = ws.ReadJSON(&data)
	defer ws.Close()
	if err != nil {
		log.Println("error read json")
		ws.Close()
		return
	}
	log.Printf("[TradeChartSocketAPI]received subscriptionMessage: %+v\n", data)

	auth.Auth = true
	auth.Message = "success"
	message := gin.H{"stream": auth, "streamType": "auth"}
	ws.WriteJSON(message)
	var oldTradeChartStr string
	exit := false
	go func() {
		fmt.Println("Started Streaming trade chart")
		message := gin.H{"stream": "Started Trade Chart Stream", "streamType": "notice"}
		messageChan <- message
		for {

			tradeChart, err := blockchain.GetChartRecords(data.AssetCode, data.AssetIssuer, data.CurrencyCode, data.CurrencyIssuer, "", "", data.Order, "", data.ChartPeriod, "")
			if err != nil {
				log.Printf("Error fetching chart. Error %v\n", err)
				exit = true
				message = gin.H{"stream": "Closing Trade Chart Stream", "streamType": "notice"}
				messageChan <- message

				break
			}
			bChart, _ := json.Marshal(tradeChart)
			if string(bChart) != oldTradeChartStr {
				message := gin.H{"stream": tradeChart, "streamType": "tradeChart"}
				messageChan <- message
				oldTradeChartStr = string(bChart)
			}

			time.Sleep(30 * time.Second)

		}

	}()

	//try to read from ws and exit if cannot read.
	go func() {
		log.Println("@@@@@@started ws connection keepalive......")
		for {
			if exit {
				break
			}
			time.Sleep(30 * time.Second)
			auth.Auth = true
			auth.Message = "keep-alive"
			message := gin.H{"stream": "keep-alive", "streamType": "keep-alive"}
			messageChan <- message

		}

	}()

	for {
		if exit {
			ws.Close()
			break
		}
		v := <-messageChan

		err = ws.WriteJSON(v)
		if err != nil {
			log.Printf("Error Sending stream: %v\nError %v\n", v, err)
			ws.Close()
			break
		}
	}

}
