package users

import (
	"os"
	"time"
	announcementServices "trovo-wallet-api/internal/components/announcements/services"
	bc "trovo-wallet-api/internal/components/users/blockchain"
	usersDB "trovo-wallet-api/internal/components/users/db"
	userModels "trovo-wallet-api/internal/components/users/models"
	"trovo-wallet-api/internal/middleware"
	"trovo-wallet-api/internal/sharedconfig"

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
	PublicKey string
	Alias     string
	request   interface{}
}

// UserWebSocketAPI handles websocket connections
func UserWebSocketAPI(c *gin.Context, gc *sharedconfig.GlobalConfig) {
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
	announcements, err := announcementServices.HandleGetAnnouncement(c.ClientIP(), gc.DB)
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

	//range through the user wallets
	//make a map of the public key and temporary key with the alias
	walletMap := make(map[string]string, 0)
	var opsRequestListMain, opsRequestListTemp, effectRequestListMain, effectRequestListTemp []StreamObject

	for _, userWallet := range user.UserWallets {
		walletMap[userWallet.ID] = userWallet.Alias
		if userWallet.TempPublicKey != nil {
			walletMap[*userWallet.TempPublicKey] = userWallet.Alias
		}

		opsRequest = horizonclient.OperationRequest{
			ForAccount: userWallet.ID,
			Cursor:     cursor,
			Order:      order,
			Join:       "transactions",
		}
		opsRequestListMain = append(opsRequestListMain, StreamObject{userWallet.ID, userWallet.Alias, opsRequest})

		effectsRequest = horizonclient.EffectRequest{
			ForAccount: userWallet.ID,
			Cursor:     cursor,
			Order:      order,
		}
		effectRequestListMain = append(effectRequestListMain, StreamObject{userWallet.ID, userWallet.Alias, effectsRequest})
		if userWallet.TempPublicKey != nil {
			tempopsRequest = horizonclient.OperationRequest{
				ForAccount: *userWallet.TempPublicKey,
				Cursor:     cursor,
				Order:      order,
				Join:       "transactions",
			}
			opsRequestListTemp = append(opsRequestListTemp, StreamObject{*userWallet.TempPublicKey, userWallet.Alias, tempopsRequest})

			tempEffectsRequest = horizonclient.EffectRequest{
				ForAccount: *userWallet.TempPublicKey,
				Cursor:     cursor,
				Order:      order,
			}
			effectRequestListTemp = append(effectRequestListTemp, StreamObject{*userWallet.TempPublicKey, userWallet.Alias, tempEffectsRequest})

		}
	}

	mainOpsStreamHandler := func(o operations.Operation) {

		// fmt.Println("/////........................Operations Stream received...........////////")
		// fmt.Println(o.GetType())
		if o.GetType() == "change_trust" {
			obj := interface{}(o).(operations.AccountMerge)
			{
				//invalidate cache of Account
				if v, ok := walletMap[obj.Account]; ok {
					getCacheKey := fmt.Sprintf("[GET] /v1/users/%v", v)
					historyCacheKey := fmt.Sprintf("[history] %v", v)
					balancesCacheKey := fmt.Sprintf("[balances] %v", v)

					gc.RedisCache.DeleteFromCache(getCacheKey, historyCacheKey, balancesCacheKey)

				}

				getCacheKey := fmt.Sprintf("[GET] /v1/users/%v", obj.Account)
				historyCacheKey := fmt.Sprintf("[history] %v", obj.Account)
				balancesCacheKey := fmt.Sprintf("[balances] %v", obj.Account)
				gc.RedisCache.DeleteFromCache(getCacheKey, historyCacheKey, balancesCacheKey)
			}
			{
				//invalidate cache of Account
				if v, ok := walletMap[obj.Into]; ok {
					getCacheKey := fmt.Sprintf("[GET] /v1/users/%v", v)
					historyCacheKey := fmt.Sprintf("[history] %v", v)
					balancesCacheKey := fmt.Sprintf("[balances] %v", v)

					gc.RedisCache.DeleteFromCache(getCacheKey, historyCacheKey, balancesCacheKey)

				}

				getCacheKey := fmt.Sprintf("[GET] /v1/users/%v", obj.Into)
				historyCacheKey := fmt.Sprintf("[history] %v", obj.Into)
				balancesCacheKey := fmt.Sprintf("[balances] %v", obj.Into)
				gc.RedisCache.DeleteFromCache(getCacheKey, historyCacheKey, balancesCacheKey)
			}

		}
		if o.GetType() == "payment" || o.GetType() == "create_account" || strings.Contains(o.GetType(), "path_payment") || o.GetType() == "account_merge" {

			if o.GetType() != "account_merge" {
				// mesg := gin.H{"stream": o, "streamType": "payment"}
				// messageChan <- mesg
				paymentPg := bc.ProcessStreamPaymentOperation(user.PublicKey, o, user, gc.DB)
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
			obj := interface{}(o).(operations.AccountMerge)

			//invalidate cache
			{
				//invalidate cache of Account
				if v, ok := walletMap[obj.Account]; ok {
					getCacheKey := fmt.Sprintf("[GET] /v1/users/%v", v)
					historyCacheKey := fmt.Sprintf("[history] %v", v)
					balancesCacheKey := fmt.Sprintf("[balances] %v", v)

					gc.RedisCache.DeleteFromCache(getCacheKey, historyCacheKey, balancesCacheKey)

				}

				getCacheKey := fmt.Sprintf("[GET] /v1/users/%v", obj.Account)
				historyCacheKey := fmt.Sprintf("[history] %v", obj.Account)
				balancesCacheKey := fmt.Sprintf("[balances] %v", obj.Account)
				gc.RedisCache.DeleteFromCache(getCacheKey, historyCacheKey, balancesCacheKey)
			}
			{
				//invalidate cache of Account
				if v, ok := walletMap[obj.Into]; ok {
					getCacheKey := fmt.Sprintf("[GET] /v1/users/%v", v)
					historyCacheKey := fmt.Sprintf("[history] %v", v)
					balancesCacheKey := fmt.Sprintf("[balances] %v", v)

					gc.RedisCache.DeleteFromCache(getCacheKey, historyCacheKey, balancesCacheKey)

				}

				getCacheKey := fmt.Sprintf("[GET] /v1/users/%v", obj.Into)
				historyCacheKey := fmt.Sprintf("[history] %v", obj.Into)
				balancesCacheKey := fmt.Sprintf("[balances] %v", obj.Into)
				gc.RedisCache.DeleteFromCache(getCacheKey, historyCacheKey, balancesCacheKey)
			}

		}
		if o.GetType() == "payment" || o.GetType() == "create_account" || strings.Contains(o.GetType(), "path_payment") || o.GetType() == "account_merge" {

			if o.GetType() != "account_merge" {
				obj := interface{}(o).(operations.AccountMerge)
				paymentPg := bc.ProcessStreamPaymentOperation(obj.Into, o, user, gc.DB)
				if strings.Contains(o.GetType(), "path_payment") {
					message := gin.H{"stream": paymentPg, "streamType": "swap"}

					messageChan <- message
				} else {
					message := gin.H{"stream": paymentPg, "streamType": "payment"}

					messageChan <- message
				}

			} else {

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
			cacheKey := fmt.Sprintf("[GET] /v1/users/%v", user.Username)
			gc.RedisCache.InvalidateCachedHttpResponse(cacheKey)
			cacheKey = fmt.Sprintf("[GET] /v1/users/%v", user.PublicKey)
			gc.RedisCache.InvalidateCachedHttpResponse(cacheKey)
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
		for _, opsRequest := range opsRequestListMain {
			wg.Add(1)
			defer wg.Done()
			fmt.Printf("Started Streaming OPERATIONS EVENTS for Public Key: %s, Username: %s\n", opsRequest.PublicKey, opsRequest.Alias)
			message := gin.H{"stream": "Started Operations Stream For Main Account", "streamType": "notice"}
			messageChan <- message
			obj := (opsRequest.request).(horizonclient.OperationRequest)
			err = client.StreamOperations(ctx, obj, mainOpsStreamHandler)
			if err != nil {
				fmt.Println("stream operation error:", err)
				// return
			}
		}

	}
	tempStreamOperations := func() {
		for _, opsRequest := range opsRequestListTemp {
			wg.Add(1)
			defer wg.Done()
			fmt.Printf("Started Streaming OPERATIONS EVENTS for Public Key: %s, Username: %s\n", opsRequest.PublicKey, opsRequest.Alias)
			message := gin.H{"stream": "Started Operations Stream For Temp Account", "streamType": "notice"}
			messageChan <- message
			obj := (opsRequest.request).(horizonclient.OperationRequest)

			err = client.StreamOperations(ctx, obj, tempOpsStreamHandler)
			if err != nil {
				fmt.Println("stream temp operation error:", err)
				// return
			}
		}
	}

	streamEffects := func() {
		for _, effectRequest := range effectRequestListMain {
			wg.Add(1)
			defer wg.Done()
			fmt.Printf("Started Streaming EFFECTS EVENTS for Public Key: %s, Username: %s\n", effectRequest.PublicKey, effectRequest.Alias)
			message := gin.H{"stream": "Started Effects Stream For Main Account", "streamType": "notice"}
			messageChan <- message
			obj := (effectRequest.request).(horizonclient.EffectRequest)

			err = client.StreamEffects(ctx, obj, effectsStreamHandler)
			if err != nil {
				fmt.Println("stream effects error:", err)
				// return
			}
		}
	}
	tempStreamEffects := func() {
		for _, effectRequest := range effectRequestListTemp {
			wg.Add(1)
			defer wg.Done()
			fmt.Printf("Started Streaming EFFECTS EVENTS for Public Key: %s, Username: %s\n", effectRequest.PublicKey, effectRequest.Alias)
			message := gin.H{"stream": "Started Effects Stream For Temp Account", "streamType": "notice"}
			messageChan <- message
			obj := (effectRequest.request).(horizonclient.EffectRequest)

			err = client.StreamEffects(ctx, obj, effectsStreamHandler)
			if err != nil {
				fmt.Println("stream temp effects error:", err)
				// return
			}
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
		if len(opsRequestListTemp) > 0 {
			go tempStreamOperations()
		}
		if len(effectRequestListTemp) > 0 {
			go tempStreamEffects()
		}

	} else {

		if data.Stream == "" || strings.Contains(data.Stream, "operations") {
			go streamOperations()
			go streamEffects()
			if len(opsRequestListTemp) > 0 {
				go tempStreamOperations()
			}
			if len(effectRequestListTemp) > 0 {
				go tempStreamEffects()
			}
		}
		if strings.Contains(data.Stream, "effects") {
			go streamEffects()
			if len(effectRequestListTemp) > 0 {
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
