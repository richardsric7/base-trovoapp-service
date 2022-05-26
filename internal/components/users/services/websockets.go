package users

import (
	"trovo-wallet-api/internal/cache"
	announcementServices "trovo-wallet-api/internal/components/announcements/services"
	bc "trovo-wallet-api/internal/components/users/blockchain"
	usersDB "trovo-wallet-api/internal/components/users/db"
	userModels "trovo-wallet-api/internal/components/users/models"
	"trovo-wallet-api/internal/middleware"
	"os"
	"time"

	"trovo-wallet-api/internal/network"

	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/protocols/horizon/effects"
	"github.com/stellar/go/protocols/horizon/operations"
	"gorm.io/gorm"
)

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

//StartMessage is models for stream start message
type StartMessage struct {
	Stream string `json:"stream"`
}

//UserWebSocketAPI handles websocket connections
func UserWebSocketAPI(c *gin.Context, db *gorm.DB, redisCache *cache.RedisCache) {
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	var order = horizonclient.OrderAsc
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
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

	user, err := usersDB.GetUserInfo(identifier, db)
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

	if !user.SignerIsValid(data.Signer, false) && authEnable {
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
	announcements, err := announcementServices.HandleGetAnnouncement(c.ClientIP(), db)
	if err == nil {
		message := gin.H{"stream": announcements, "streamType": "announcements"}
		ws.WriteJSON(message)

	}

	var cursor string
	data.Cursor = strings.ToLower(data.Cursor)
	if data.Cursor == "all" {
		cursor = ""
	} else if len(data.Cursor) > 0 {
		cursor = data.Cursor
		order = horizonclient.OrderAsc
	} else {
		cursor = "now"
	}

	// fmt.Println("Stream Params: Cursor: ", cursor, ", Order :", order, ", Limit: ", limit)
	client := network.GetBlockchainClient()
	var opsRequest, tempopsRequest horizonclient.OperationRequest
	var effectsRequest, tempEffectsRequest horizonclient.EffectRequest
	opsRequest = horizonclient.OperationRequest{
		ForAccount: user.PublicKey,
		Cursor:     cursor,
		Order:      order,
		Join:       "transactions",
	}
	effectsRequest = horizonclient.EffectRequest{
		ForAccount: user.PublicKey,
		Cursor:     cursor,
		Order:      order,
	}
	if user.TempPublicKey != nil {
		tempopsRequest = horizonclient.OperationRequest{
			ForAccount: *user.TempPublicKey,
			Cursor:     cursor,
			Order:      order,
			Join:       "transactions",
		}
		tempEffectsRequest = horizonclient.EffectRequest{
			ForAccount: *user.TempPublicKey,
			Cursor:     cursor,
			Order:      order,
		}

	}

	mainOpsStreamHandler := func(o operations.Operation) {

		// fmt.Println("/////........................Operations Stream received...........////////")
		// fmt.Println(o.GetType())
		if o.GetType() == "change_trust" {
			//invalidate cache
			cacheKey := fmt.Sprintf("[GET] /v2/users/%v", user.Username)
			redisCache.InvalidateCachedHttpResponse(cacheKey)
			cacheKey = fmt.Sprintf("[GET] /v2/users/%v", user.PublicKey)
			redisCache.InvalidateCachedHttpResponse(cacheKey)
		}
		if o.GetType() == "payment" || o.GetType() == "create_account" || strings.Contains(o.GetType(), "path_payment") || o.GetType() == "account_merge" {

			if o.GetType() != "account_merge" {
				// mesg := gin.H{"stream": o, "streamType": "payment"}
				// messageChan <- mesg
				paymentPg := bc.ProcessStreamPaymentOperation(user.PublicKey, o, user, db)
				if strings.Contains(o.GetType(), "path_payment") {
					message := gin.H{"stream": paymentPg, "streamType": "swap"}

					messageChan <- message
				} else {
					message := gin.H{"stream": paymentPg, "streamType": "payment"}

					messageChan <- message
				}
				log.Println("<<<<<<<<<<<<<<<<<<<<<<< payment object stream received<<<<<<<<<<<<<<<<<<")
				log.Println(paymentPg)
				log.Println("<<<<<<<<<<<<<<<<<<<<<<")
			} else {
				opsData := actionData{
					ID:              o.GetID(),
					TransactionHash: o.GetTransactionHash(),
					TxType:          o.GetType(),
				}
				message := gin.H{"stream": opsData, "streamType": "operation"}

				messageChan <- message
			}
		} else if o.GetType() == "manage_data" {
			txType := o.GetType()

			userData := actionData{
				ID:              o.GetID(),
				TransactionHash: o.GetTransactionHash(),
				TxType:          txType,
			}
			message := gin.H{"stream": userData, "streamType": "manage_data"}

			messageChan <- message
		} else {
			opsData := actionData{
				ID:              o.GetID(),
				TransactionHash: o.GetTransactionHash(),
				TxType:          o.GetType(),
			}
			message := gin.H{"stream": opsData, "streamType": "operation"}

			messageChan <- message
		}
	}

	tempOpsStreamHandler := func(o operations.Operation) {

		// fmt.Println("/////........................Operations Stream received...........////////")
		// fmt.Println(o.GetType())
		if o.GetType() == "change_trust" {
			//invalidate cache
			cacheKey := fmt.Sprintf("[GET] /v2/users/%v", user.Username)
			redisCache.InvalidateCachedHttpResponse(cacheKey)
			cacheKey = fmt.Sprintf("[GET] /v2/users/%v", user.TempPublicKey)
			redisCache.InvalidateCachedHttpResponse(cacheKey)
		}
		if o.GetType() == "payment" || o.GetType() == "create_account" || strings.Contains(o.GetType(), "path_payment") || o.GetType() == "account_merge" {

			if o.GetType() != "account_merge" {
				paymentPg := bc.ProcessStreamPaymentOperation(*user.TempPublicKey, o, user, db)
				if strings.Contains(o.GetType(), "path_payment") {
					message := gin.H{"stream": paymentPg, "streamType": "swap"}

					messageChan <- message
				} else {
					message := gin.H{"stream": paymentPg, "streamType": "payment"}

					messageChan <- message
				}

			} else {
				// TODO: handle account merge
				userData := actionData{
					ID:              o.GetID(),
					TransactionHash: o.GetTransactionHash(),
					TxType:          "account_merge",
				}
				message := gin.H{"stream": userData, "streamType": "payment"}

				messageChan <- message
			}
		} else if o.GetType() == "manage_data" {

			userData := actionData{
				ID:              o.GetID(),
				TransactionHash: o.GetTransactionHash(),
				TxType:          o.GetType(),
			}
			message := gin.H{"stream": userData, "streamType": "manage_data"}

			messageChan <- message
		} else {
			opsData := actionData{
				ID:              o.GetID(),
				TransactionHash: o.GetTransactionHash(),
				TxType:          o.GetType(),
			}
			message := gin.H{"stream": opsData, "streamType": "operation"}

			messageChan <- message
		}
	}
	effectsStreamHandler := func(o effects.Effect) {
		if o.GetType() == "change_trust" {
			//invalidate cache
			cacheKey := fmt.Sprintf("[GET] /v2/users/%v", user.Username)
			redisCache.InvalidateCachedHttpResponse(cacheKey)
			cacheKey = fmt.Sprintf("[GET] /v2/users/%v", user.PublicKey)
			redisCache.InvalidateCachedHttpResponse(cacheKey)
		}
		txType := o.GetType()
		if o.GetType() == "payment" || o.GetType() == "create_account" || strings.Contains(o.GetType(), "path_payment") {
			txType = "payment"
		}

		opsData := actionData{
			ID:              o.GetID(),
			TransactionHash: "",
			TxType:          txType,
		}
		message := gin.H{"stream": opsData, "streamType": "effect"}

		messageChan <- message

	}

	streamOperations := func() {
		wg.Add(1)
		defer wg.Done()
		fmt.Printf("Started Streaming OPERATIONS EVENTS for Public Key: %s, Username: %s\n", user.PublicKey, user.Username)
		message := gin.H{"stream": "Started Operations Stream For Main Account", "streamType": "notice"}
		messageChan <- message
		err = client.StreamOperations(ctx, opsRequest, mainOpsStreamHandler)
		if err != nil {
			fmt.Println("stream operation error:", err)
			// return
		}

	}
	tempStreamOperations := func() {
		wg.Add(1)
		defer wg.Done()
		fmt.Printf("Started Streaming OPERATIONS EVENTS for Public Key: %s, Username: %s\n", *user.TempPublicKey, user.Username)
		message := gin.H{"stream": "Started Operations Stream For Temp Account", "streamType": "notice"}
		messageChan <- message
		err = client.StreamOperations(ctx, tempopsRequest, tempOpsStreamHandler)
		if err != nil {
			fmt.Println("stream temp operation error:", err)
			// return
		}

	}

	streamEffects := func() {
		wg.Add(1)
		defer wg.Done()
		fmt.Printf("Started Streaming EFFECTS EVENTS for Public Key: %s, Username: %s\n", user.PublicKey, user.Username)
		message := gin.H{"stream": "Started Effects Stream For Main Account", "streamType": "notice"}
		messageChan <- message
		err = client.StreamEffects(ctx, effectsRequest, effectsStreamHandler)
		if err != nil {
			fmt.Println("stream effects error:", err)
			// return
		}

	}
	tempStreamEffects := func() {
		wg.Add(1)
		defer wg.Done()
		fmt.Printf("Started Streaming EFFECTS EVENTS for Public Key: %s, Username: %s\n", *user.TempPublicKey, user.Username)
		message := gin.H{"stream": "Started Effects Stream For Temp Account", "streamType": "notice"}
		messageChan <- message
		err = client.StreamEffects(ctx, tempEffectsRequest, effectsStreamHandler)
		if err != nil {
			fmt.Println("stream temp effects error:", err)
			// return
		}

	}

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

	if strings.Contains(data.Stream, "all") {
		go streamOperations()
		go streamEffects()
		if user.TempPublicKey != nil {
			go tempStreamOperations()
			go tempStreamEffects()
		}

	} else {

		if data.Stream == "" || strings.Contains(data.Stream, "operations") {
			go streamOperations()
			go streamEffects()
			if user.TempPublicKey != nil {
				go tempStreamOperations()
				go streamEffects()
			}
		}
		if strings.Contains(data.Stream, "effects") {
			go streamEffects()
			if user.TempPublicKey != nil {
				go tempStreamEffects()
			}
		}

	}
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
