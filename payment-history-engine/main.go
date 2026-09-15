package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	// bc "trovo-wallet-payment-history-engine/internal/blockchainalgofuncs"
	bc "trovo-wallet-payment-history-engine/internal/blockchainalgofuncs"
	// merchants "trovo-wallet-payment-history-engine/internal/components/merchants/controllers"
	"trovo-wallet-payment-history-engine/internal/cache"
	"trovo-wallet-payment-history-engine/internal/components/health"
	paymentModels "trovo-wallet-payment-history-engine/internal/components/payments/models"
	paymentServices "trovo-wallet-payment-history-engine/internal/components/payments/services"

	// users "trovo-wallet-payment-history-engine/internal/components/users/controllers"
	userModels "trovo-wallet-payment-history-engine/internal/components/users/models"
	db "trovo-wallet-payment-history-engine/internal/db"
	"trovo-wallet-payment-history-engine/internal/network"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"github.com/shopspring/decimal"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/protocols/horizon"
	"github.com/stellar/go/protocols/horizon/operations"
	"github.com/stellar/go/txnbuild"
	"gorm.io/gorm"
)

// startStreaming and trackPublicKey are read and written from multiple goroutines,
// so they're stored as int32s and accessed exclusively through atomic ops (0 = false, 1 = true)
// instead of as plain bools, which would be a data race under concurrent access.
var startStreaming int32
var trackPublicKey int32

func setStartStreaming(v bool) {
	if v {
		atomic.StoreInt32(&startStreaming, 1)
		return
	}
	atomic.StoreInt32(&startStreaming, 0)
}

func isStartStreaming() bool {
	return atomic.LoadInt32(&startStreaming) == 1
}

func setTrackPublicKey(v bool) {
	if v {
		atomic.StoreInt32(&trackPublicKey, 1)
		return
	}
	atomic.StoreInt32(&trackPublicKey, 0)
}

func isTrackPublicKey() bool {
	return atomic.LoadInt32(&trackPublicKey) == 1
}

func main() {

	//setup environment variables
	errEnv := godotenv.Load()
	if errEnv != nil {
		path, _ := os.Getwd()
		log.Printf("could not find or load any .env file from %v...skipping...\n", path)
	}

	//setup DB

	var database, roachDB *gorm.DB
	{

		var err error
		database, err = db.OpenDb()

		if err != nil {
			log.Fatalf("[main]Error opening DB %s", err)
			return
		}
	}
	{

		var err error
		roachDB, err = db.OpenRoachDB()

		if err != nil {
			log.Fatalf("[main]Error opening RoachDB %s", err)
			return
		}
	}

	//check other required environment variables.

	{
		exit := false
		requiredEnvironmentVariables := []string{"EXPANSION_URL", "BLOCKCHAIN_NETWORK_PASSPHRASE",
			"MNEMONIC_TEMP_ACCOUNTS", "BLOCKCHAIN_BASE_RESERVE", "NATIVE_ASSET_CODE", "ENABLE_CACHING",
			"MM_FEE_COLLECTION_CHANNEL_ACCOUNT", "MARKET_MAKING_SALT", "MNEMONIC_MARKET_MAKING",
			"MARKET_MAKING_FEE_WALLET",
		}

		for _, requiredEnvironmentVariable := range requiredEnvironmentVariables {
			if len(os.Getenv(requiredEnvironmentVariable)) == 0 {
				log.Printf("Required environment variable is missing %v", requiredEnvironmentVariable)
				exit = true
			}
		}
		if os.Getenv("ENABLE_EMAIL_VALIDATION") == "1" && len(os.Getenv("MAILGUN_VALIDATOR_API_KEY")) == 0 {
			log.Println("MAILGUN_VALIDATOR_API_KEY environment variable is required when ENABLE_EMAIL_VALIDATION is set to 1")

			exit = true
		}

		if exit {
			return
		}

	}

	if os.Getenv("PURGE_TABLES") == "1" {

		log.Println("DROPPING TRACKING TABLES")
		errMigrate := roachDB.Raw("DROP TABLE IF EXISTS tracked_wallets").Error
		if errMigrate != nil {
			if !strings.Contains(errMigrate.Error(), "constraint") {
				log.Fatalf("Error migrating TrackedWallet model, error: %v", errMigrate)
			}
		}

		errMigrate = roachDB.Raw("DROP TABLE IF EXISTS tracked_public_keys").Error
		if errMigrate != nil {
			if !strings.Contains(errMigrate.Error(), "constraint") {
				log.Fatalf("Error migrating TrackedWallet model, error: %v", errMigrate)
			}
		}
		var tbd paymentModels.MonitoredCursor
		e := roachDB.First(&tbd).Error
		if e == nil {
			roachDB.Delete(&tbd)
		}
		errMigrate = roachDB.Raw("DROP TABLE IF EXISTS monitored_cursors").Error
		if errMigrate != nil {
			if !strings.Contains(errMigrate.Error(), "constraint") {
				log.Fatalf("Error migrating TrackedWallet model, error: %v", errMigrate)
			}
		}

		log.Println("FINISHED DROPPING TRACKING TABLES................")
		return
	}

	//migrate roach DB models if any

	errMigrate := roachDB.AutoMigrate(&paymentModels.TrackedWallet{})
	if errMigrate != nil {
		if !strings.Contains(errMigrate.Error(), "constraint") {
			log.Fatalf("Error migrating TrackedWallet model, error: %v", errMigrate)
		}
	}
	errMigrate = roachDB.AutoMigrate(&paymentModels.TrackedPublicKey{})
	if errMigrate != nil {
		if !strings.Contains(errMigrate.Error(), "constraint") {
			log.Fatalf("Error migrating TrackedPublicKey model, error: %v", errMigrate)
		}
	}
	errMigrate = roachDB.AutoMigrate(&paymentModels.MonitoredCursor{})
	if errMigrate != nil {
		if !strings.Contains(errMigrate.Error(), "constraint") {
			log.Fatalf("Error migrating MonitoredCursor model, error: %v", errMigrate)
		}
	}
	errMigrate = roachDB.AutoMigrate(&paymentModels.MonitoredAccountCursor{})
	if errMigrate != nil {
		if !strings.Contains(errMigrate.Error(), "constraint") {
			log.Fatalf("Error migrating MonitoredAccountCursor model, error: %v", errMigrate)
		}
	}

	//setup redis
	enableCaching := false

	if os.Getenv("ENABLE_CACHING") == "1" {
		if len(os.Getenv("REDIS_HOST")) == 0 || len(os.Getenv("REDIS_PORT")) == 0 {
			log.Fatalln("redis host of port not configured properly")
		}
		if len(os.Getenv("CACHING_PARAMETER")) == 0 {
			log.Fatalln("env CACHING_PARAMETER not configured")
		}
		enableCaching = true
	}

	var redisCli *redis.Client = nil

	if enableCaching {
		redisCli = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%v:%v", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")), // host:port of the redis server
			Password: os.Getenv("REDIS_PASSWORD"),                                            // no password set
			DB:       0,                                                                      // use default DB
		})
	}

	var redisCache cache.RedisCache = cache.RedisCache{
		Enabled: enableCaching,
		Client:  redisCli,
		Context: context.Background(),
	}

	// Start the health probes once the dependencies exist. This engine has no
	// HTTP server of its own, so the probes bring their own - it runs in a
	// goroutine and never blocks payment processing.
	health.Serve(health.New(database, roachDB, &redisCache, network.GetBlockchainClient()))

	if redisCache.Enabled {
		//test redis connection
		log.Println("Testing redis connection...")
		_, err := redisCli.Ping(redisCache.Context).Result()
		if err != nil {
			log.Printf("[main]Error connecting to redis %s", err)
			time.Sleep(time.Second * 5)
			return
		}
		if !redisCache.StoreResultToCache("test", "test", 5) {
			log.Println("[main] unable to store to redis")
			time.Sleep(time.Second * 5)
			return
		}
	}

	{

		//track untracked user wallets
		go func() {
			log.Println("##[TRACKUserWallet] started routine to update user wallets that are not tracked")
			lastTracked := ""
			var startTrackingFrom time.Time
			if len(lastTracked) > 0 {
				startTrackingFrom, _ = time.Parse(time.RFC3339, lastTracked)
			}
			batchSize := 500
			var userWallets []userModels.UserWallet
			log.Println("[TRACKUserWallet] startTrackingFrom:", startTrackingFrom)
			for {
				e := database.Where("created_at::timestamp >= ?::timestamp AND tracked = 0", startTrackingFrom).First(&userModels.UserWallet{}).Error
				if e != nil {
					log.Println("[TRACKUserWallet] error getting user wallet:", e)
					setStartStreaming(true)
					if errors.Is(e, gorm.ErrRecordNotFound) {
						setTrackPublicKey(true)
					}

					time.Sleep(2 * time.Minute)
					continue
				}

				result := database.Where("created_at::timestamp >= ?::timestamp AND tracked = 0", startTrackingFrom).FindInBatches(&userWallets, batchSize, func(tx *gorm.DB, batch int) error {
					for i, u := range userWallets {
						log.Printf("[TRACKUserWallet] processing user wallet %+v of %v\n", u, batch)
						trackError := paymentServices.TrackUserWallet(u, roachDB, database, isTrackPublicKey(), &redisCache)
						if trackError == nil {
							userWallets[i].Tracked = 1
						} else {
							//leave Tracked at 0 so this wallet is retried on the next pass instead of aborting the whole batch/process
							log.Printf("[TRACKUserWallet] error adding user wallet to tracking table, due to: %v\n", trackError)
						}
					}
					e := tx.Save(&userWallets).Error
					if e != nil {
						log.Printf("[TRACKUserWallet] unable to update user wallets for batch %v due to error: %v\n", batch, e)
					} else {
						log.Printf("[TRACKUserWallet] successfully updated %v user wallets for batch %v\n", tx.RowsAffected, batch)
					}

					return nil
				})
				if result.Error != nil {
					if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
						log.Println("[TRACKUserWallet]()()()()()()()()()()()()()error occurred during batch processing:", result.Error.Error())

					}
					if errors.Is(e, gorm.ErrRecordNotFound) {
						setTrackPublicKey(true)
					}

				}

				time.Sleep(5 * time.Second)
				log.Printf("[TRACKUserWallet] Total: %v, Errors:%v restarting process...\n", result.RowsAffected, result.Error)
				setStartStreaming(true)
			}
		}()

	}

	{

		//monitor public keys
		go func() {
			log.Println("##[TRACKPublicKeys] started routine to generate payment history of public keys")

			batchSize := 100
			var userPublicKeys []paymentModels.TrackedPublicKey
			for {
				e := roachDB.First(&paymentModels.TrackedPublicKey{}).Error
				if e != nil {
					log.Println("[TRACKPublicKeys] error getting public keys to process", e)
					if errors.Is(e, gorm.ErrRecordNotFound) {
						setTrackPublicKey(true)
					}
					time.Sleep(1 * time.Minute)
					continue
				}
				var wg sync.WaitGroup
				result := roachDB.FindInBatches(&userPublicKeys, batchSize, func(tx *gorm.DB, batch int) error {
					for _, u := range userPublicKeys {
						wg.Add(1)
						log.Printf("[TRACKPublicKeys] processing user publicKey %v of %v\n", u.PublicKey, batch)
						//spin off worker to process the payment history of the public key
						go MonitorPublicKeyPaymentStream(u.PublicKey, database, roachDB, &wg)
					}

					return nil
				})
				if result.Error != nil {
					if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
						log.Println("[TRACKPublicKeys]()()()()()()()()()()()()()error occurred during batch processing:", result.Error.Error())

					}
					if errors.Is(e, gorm.ErrRecordNotFound) {
						setTrackPublicKey(true)
					}

				}
				//All waiting public keys finished. wait for them to complete generating their history
				wg.Wait()
				time.Sleep(1 * time.Minute)
				log.Printf("[TRACKPublicKeys] Total: %v, Errors:%v restarting process...\n", result.RowsAffected, result.Error)
				setTrackPublicKey(true)
			}
		}()

	}

	//Start processing payment streams
	go func() {
		for {
			if !isStartStreaming() {
				time.Sleep(5 * time.Second)
				continue
			}
			MonitorPaymentStream(database, roachDB)
			time.Sleep(5 * time.Second)
		}
	}()

	for {
		if !isStartStreaming() {
			time.Sleep(5 * time.Second)
			continue
		}
		MonitorTradeStream(database, roachDB, &redisCache)
		time.Sleep(5 * time.Second)
	}

}

func GetLastCursor(roachDB *gorm.DB) (lastCursor string) {
	var mAccount paymentModels.MonitoredCursor
	var envCusor string
	if len(os.Getenv("LAST_CURSOR")) > 0 {
		envCusor = os.Getenv("LAST_CURSOR")
	} else {
		envCusor = "0"
	}

	e := roachDB.First(&mAccount).Error
	if e == nil {
		if decimal.RequireFromString(envCusor).GreaterThan(decimal.RequireFromString(mAccount.LastCursor)) {
			log.Println("[GetLastCursor] returning ENV cursor since it appears more recent.....")
			return envCusor
		}
		log.Println("[GetLastCursor] returning DB cursor since it appears more recent.....")

		return mAccount.LastCursor
	}
	log.Println("[GetLastCursor] returning DB cursor since could not fetch from db.....")
	return envCusor

}

func GetTradeResumeCursor() string {
	var envCusor string
	if len(os.Getenv("TRADE_RESUME_CURSOR")) > 0 {
		envCusor = os.Getenv("TRADE_RESUME_CURSOR")
	} else {
		envCusor = "0"
	}

	return envCusor

}

func SaveLastCursor(lastCursor string, roachDB *gorm.DB) (e error) {
	var mAccount paymentModels.MonitoredCursor
	e = roachDB.First(&mAccount).Error
	if e != nil {
		mAccount = paymentModels.MonitoredCursor{ID: 1, LastCursor: lastCursor}
		e := roachDB.Create(&mAccount).Error
		if e != nil {
			log.Println("[SaveLastCursor]error saving last cursor", e)
			return e
		}
		return e
	}
	if mAccount.LastCursor != lastCursor {
		mAccount.LastCursor = lastCursor
		e = roachDB.Save(&mAccount).Error
		if e != nil {
			log.Println("[SaveLastCursor] unable to save last cursor", e)
			return e
		}
	}
	log.Println("[SaveLastCursor]saved cursor:", lastCursor)
	return nil

}

// GetAccountCursor returns the last processed paging token for a single tracked public key's
// genesis backfill stream, or "0" (start from genesis) if it has never been recorded.
func GetAccountCursor(publicKey string, roachDB *gorm.DB) string {
	var mAccount paymentModels.MonitoredAccountCursor
	e := roachDB.Where("public_key = ?", publicKey).First(&mAccount).Error
	if e != nil {
		return "0"
	}
	return mAccount.LastCursor
}

// SaveAccountCursor persists the last processed paging token for a single tracked public key,
// so its backfill stream can resume from where it left off instead of always restarting at genesis.
func SaveAccountCursor(publicKey, lastCursor string, roachDB *gorm.DB) error {
	var mAccount paymentModels.MonitoredAccountCursor
	e := roachDB.Where("public_key = ?", publicKey).First(&mAccount).Error
	if e != nil {
		if !errors.Is(e, gorm.ErrRecordNotFound) {
			return e
		}
		mAccount = paymentModels.MonitoredAccountCursor{PublicKey: publicKey, LastCursor: lastCursor}
		return roachDB.Create(&mAccount).Error
	}
	if mAccount.LastCursor == lastCursor {
		return nil
	}
	mAccount.LastCursor = lastCursor
	return roachDB.Save(&mAccount).Error
}

func MonitorPaymentStream(db, roachDB *gorm.DB) {
	checkExists := roachDB.First(&paymentModels.TrackedWallet{}).Error
	if checkExists != nil {
		if !errors.Is(checkExists, gorm.ErrRecordNotFound) {
			log.Fatalln("[MonitorPaymentStream]unknown error while checking tracked wallets:", checkExists)
		}
		// No wallets to watch means no work, which the health probes report as
		// idle rather than stalled - processing nothing is correct when there
		// is nothing to process.
		health.SetHasWork(false)
		log.Println("[MonitorPaymentStream] exited because no tracked wallets exists:", checkExists)
		return
	}
	health.SetHasWork(true)
	go func() {
		for {
			log.Println("[MonitorPaymentStream] <<<<<<<<<<<<active and processing transactions>>>>>>>>>>>>>")
			time.Sleep(120 * time.Second)

		}
	}()
	client := network.GetBlockchainClient()
	workerChan := make(chan operations.Operation, 200000)
	lastCursor := GetLastCursor(roachDB)
	var opsRequest horizonclient.OperationRequest
	if len(lastCursor) > 0 && lastCursor != "0" {
		log.Printf("[MonitorPaymentStream] Starting monitoring from cursor[%v]\n", lastCursor)
		opsRequest = horizonclient.OperationRequest{
			Cursor: lastCursor,
			Order:  horizonclient.OrderAsc,
			Join:   "transactions",
		}
	} else {
		log.Println("[MonitorPaymentStream] Starting monitoring from current state of blockchain")

		// opsRequest = horizonclient.OperationRequest{
		// 	Cursor: "0",
		// 	Order:  horizonclient.OrderAsc,
		// 	Join:   "transactions",
		// }
		opsRequest = horizonclient.OperationRequest{
			Join: "transactions",
		}
	}
	worker := func() {
		for {
			o := <-workerChan
			ProcessOperation(o, db, roachDB)
			// Tell the health probes the stream is advancing. Without this the
			// readiness endpoint could only report that the process is alive,
			// which says nothing about whether payment history is still being
			// written - the failure that actually matters here.
			health.RecordOperation(o.PagingToken())
			//persist progress after each operation so a reconnect/restart resumes here instead of
			//replaying the whole stream or skipping whatever happened during the gap. Only one
			//worker consumes workerChan (see below), so these writes stay in stream order.
			if e := SaveLastCursor(o.PagingToken(), roachDB); e != nil {
				log.Printf("[MonitorPaymentStream.worker] unable to save last cursor: %v\n", e)
			}
		}

	}
	{
		//start a single worker: cursor persistence above assumes in-order processing,
		//so do not enable a second concurrent worker on this channel.
		go worker()

	}

	operationsStreamHandler := func(o operations.Operation) {
		//send to worker channel
		workerChan <- o
	}

	ctx, cancel := context.WithCancel(context.Background())

	streamOperations := func() {

		err := client.StreamPayments(ctx, opsRequest, operationsStreamHandler)
		if err != nil {
			log.Printf("[MonitorPaymentStream.StreamPayments]stream error:[%v]", err)
			cancel()
		}

	}

	//Start stream
	streamOperations()
	log.Println("#####[MonitorPaymentStream]...Ending streaming operation")
	//close channels
	// close(workerChan)

}

func MonitorTradeStream(db, roachDB *gorm.DB, redisCache *cache.RedisCache) {
	checkExists := db.Where("remaining_quantity::numeric > ? AND canceled = 0", 0).First(&paymentModels.MarketOffer{}).Error

	if checkExists != nil {
		if !errors.Is(checkExists, gorm.ErrRecordNotFound) {
			log.Println("[MonitorTradeStream]unknown error while checking market offers:", checkExists)
		}
		log.Println("[MonitorTradeStream] exited because no market offers exists:", checkExists)
		return
	}
	go func() {
		for {
			log.Println("[MonitorTradeStream] <<<<<<<<<<<<active and processing trade executions>>>>>>>>>>>>>")
			time.Sleep(120 * time.Second)

		}
	}()
	ctx, cancel := context.WithCancel(context.Background())
	client := network.GetBlockchainClient()
	tradeWorkerChan := make(chan horizon.Trade, 200000)
	resumeCursor := GetTradeResumeCursor()
	var tradeRequest horizonclient.TradeRequest
	if len(resumeCursor) > 0 && resumeCursor != "0" {
		log.Printf("[MonitorTradeStream] Starting monitoring from cursor[%v]\n", resumeCursor)
		tradeRequest = horizonclient.TradeRequest{
			Cursor: resumeCursor,
			Order:  horizonclient.OrderAsc,
		}
	} else {
		log.Println("[MonitorTradeStream] Starting monitoring from current state of blockchain")

		tradeRequest = horizonclient.TradeRequest{}
	}
	worker := func() {
		for {
			tr := <-tradeWorkerChan
			ProcessTrade(tr, db, roachDB, redisCache, cancel)
		}

	}
	{
		//start two workers
		go worker()
		// go worker()

	}

	tradeStreamHandler := func(t horizon.Trade) {
		//send to worker channel
		tradeWorkerChan <- t
	}

	streamTrades := func() {

		err := client.StreamTrades(ctx, tradeRequest, tradeStreamHandler)
		if err != nil {
			log.Printf("[MonitorTradeStream.StreamTrades]stream error:[%v]", err)
			cancel()
		}

	}

	//Start stream
	streamTrades()
	log.Println("#####[MonitorTradeStream]...Ending streaming trades")
	//close channels
	// close(workerChan)

}

func MonitorPublicKeyPaymentStream(publicKey string, db, roachDB *gorm.DB, wg *sync.WaitGroup) {
	defer wg.Done()
	checkExists := roachDB.Where("public_key = ? OR temp_public_key = ?", publicKey, publicKey).First(&paymentModels.TrackedWallet{}).Error
	if checkExists != nil {
		log.Printf("[MonitorPublicKeyPaymentStream] aborting because %v could not be found in tracked wallets table:", checkExists)
		return
	}
	client := network.GetBlockchainClient()
	workerChan := make(chan operations.Operation, 200000)

	//resume from wherever this account's backfill last left off (defaults to genesis "0"
	//the first time this key is ever monitored) instead of always restarting from genesis.
	lastAccountCursor := GetAccountCursor(publicKey, roachDB)
	log.Printf("[MonitorPublicKeyPaymentStream] Starting monitoring for %v from cursor[%v]\n", publicKey, lastAccountCursor)

	opsRequest := horizonclient.OperationRequest{
		ForAccount: publicKey,
		Cursor:     lastAccountCursor,
		Order:      horizonclient.OrderAsc,
		Join:       "transactions",
	}
	ctx, cancel := context.WithCancel(context.Background())
	worker := func() {
		for {
			//select from worker channel or timeout after 3mins of waiting.
			select {
			case o := <-workerChan:
				ProcessOperation(o, db, roachDB)
				//persist progress so the next periodic pass over tracked_public_keys resumes
				//here instead of re-scanning from genesis.
				if e := SaveAccountCursor(publicKey, o.PagingToken(), roachDB); e != nil {
					log.Printf("[MonitorPublicKeyPaymentStream] unable to save account cursor for %v: %v\n", publicKey, e)
				}
			case <-time.After(30 * time.Second):
				log.Printf("[MonitorPublicKeyPaymentStream] timeout waiting for stream on %v, will resume from saved cursor next pass\n", publicKey)
				//NOTE: intentionally does not delete the tracked_public_keys row here anymore.
				//A brand new/unfunded wallet would otherwise hit this timeout on its very first
				//pass and be dropped from backfill forever. Leaving the row in place lets the
				//periodic scan in main() retry it later, resuming cheaply from the saved cursor.
				cancel()
				return

			}

		}

	}
	{
		//start a single worker: cursor persistence above assumes in-order processing,
		//so do not enable a second concurrent worker on this channel.
		go worker()

	}

	operationsStreamHandler := func(o operations.Operation) {
		//send to worker channel
		workerChan <- o
	}

	streamOperations := func() {

		err := client.StreamPayments(ctx, opsRequest, operationsStreamHandler)
		if err != nil {
			log.Printf("[MonitorPublicKeyPaymentStream.StreamPayments]stream error:[%v]", err)
			cancel()
		}

	}

	//Start stream
	streamOperations()
	log.Println("#####[MonitorPublicKeyPaymentStream]...Ending streaming operation")
	//close channels
	// close(workerChan)

}

// swapTransactionType labels a path payment as a plain SWAP, or as a MINT/BURN TOKEN (SWAP ...)
// when one leg's issuer is also that leg's sender/receiver, mirroring the plain "payment"
// operation's MINT TOKEN / BURN TOKEN classification (an issuer sending its own asset is a mint,
// an issuer receiving its own asset is a burn) so mint/burn activity routed through a path
// payment isn't hidden behind a generic swap label.
func swapTransactionType(from, sourceAssetIssuer, sourceAssetCode, to, destinationAssetIssuer, destinationAssetCode string) string {
	swapLabel := fmt.Sprintf("SWAP %s>%s", sourceAssetCode, destinationAssetCode)
	isMint := sourceAssetIssuer != "" && from == sourceAssetIssuer
	isBurn := destinationAssetIssuer != "" && to == destinationAssetIssuer

	switch {
	case isMint && isBurn:
		return fmt.Sprintf("MINT/BURN TOKEN (%s)", swapLabel)
	case isMint:
		return fmt.Sprintf("MINT TOKEN (%s)", swapLabel)
	case isBurn:
		return fmt.Sprintf("BURN TOKEN (%s)", swapLabel)
	default:
		return swapLabel
	}
}

func ProcessOperation(o operations.Operation, db, roachDB *gorm.DB) {
	//cursor persistence lives in each caller's worker (MonitorPaymentStream saves the shared
	//global cursor, MonitorPublicKeyPaymentStream saves a per-account cursor) rather than here,
	//since this same function is shared by both streams and they must not share one cursor.

	if o.GetType() == "payment" {
		log.Println("found payment operation....beginning processing")
		pmt := interface{}(o).(operations.Payment)
		//send out to be saved to db
		assetCode := os.Getenv("NATIVE_ASSET_CODE")
		if len(pmt.Code) > 0 {
			assetCode = pmt.Code
		}
		//check if public key exists within trovo user ecosystem
		trackedWallets := make([]paymentModels.TrackedWallet, 0)
		// if pmt.From == "GCWNKFHXYJ7XW6ZL3UFTKXBSRFK7EKBLXXRZQR6PIK3N2KBKQ74I3RIC" || pmt.To == "GCWNKFHXYJ7XW6ZL3UFTKXBSRFK7EKBLXXRZQR6PIK3N2KBKQ74I3RIC" {
		// 	log.Println("[ProcessOperation]^^^^^^^found payment from/to GCWNKFHXYJ7XW6ZL3UFTKXBSRFK7EKBLXXRZQR6PIK3N2KBKQ74I3RIC")
		// 	time.Sleep(10 * time.Second)
		// }
		// dbFetchError := roachDB.Where(paymentModels.TrackedWallet{PublicKey: pmt.From}).Or(paymentModels.TrackedWallet{PublicKey: pmt.To}).Find(&trackedWallets).Error
		dbFetchError := roachDB.Where("(public_key = ? OR temp_public_key = ?) OR (public_key = ? OR temp_public_key = ?)", pmt.From, pmt.From, pmt.To, pmt.To).Find(&trackedWallets).Error
		var fromAlias, fromName, toAlias, toName string
		if dbFetchError == nil {
			if len(trackedWallets) == 0 {
				log.Println("[ProcessOperation] public key not found in trovo user ecosystem")
				return
			}
			//Trovo user exists
			log.Printf("[ProcessOperation] found trovo user wallet %+v\n", trackedWallets)
			for _, t := range trackedWallets {

				if t.PublicKey == pmt.From {
					fromAlias = t.Alias
					fromName = t.Name
				}
				if t.PublicKey == pmt.To {
					toAlias = t.Alias
					toName = t.Name
				}
			}

			paymentType := "PAYMENT"
			if pmt.From == pmt.Issuer {
				paymentType = "MINT TOKEN"
			} else if pmt.To == pmt.Issuer {
				paymentType = "BURN TOKEN"
			}

			e := paymentServices.SavePaymentHistory(pmt.From, fromAlias, fromName, pmt.To, toAlias, toName, pmt.Transaction.Memo, pmt.Issuer, assetCode, pmt.Amount, pmt.TransactionHash, paymentType, pmt.PT, pmt.ID, fmt.Sprintf("%v", pmt.Transaction.AccountSequence), pmt.LedgerCloseTime, db)
			if e != nil {
				log.Println("[ProcessOperation] unable to save SavePaymentHistory with payment Operation:", e)
				return
			}
		} else {
			log.Println("[ProcessOperation] unable to find tracked wallets for payment due to:", dbFetchError)
		}

	} else if o.GetType() == "create_account" {
		log.Println("found create account operation....beginning processing")
		pmt := interface{}(o).(operations.CreateAccount)

		//send out to be saved to db
		var trackedWallets []paymentModels.TrackedWallet
		// dbFetchError := roachDB.Where(paymentModels.TrackedWallet{PublicKey: pmt.Funder}).Or(paymentModels.TrackedWallet{PublicKey: pmt.Account}).Find(&trackedWallets).Error
		dbFetchError := roachDB.Where("(public_key = ? OR temp_public_key = ?) OR (public_key = ? OR temp_public_key = ?)", pmt.Funder, pmt.Funder, pmt.Account, pmt.Account).Find(&trackedWallets).Error
		var fromAlias, fromName, toAlias, toName string
		if dbFetchError == nil {
			if len(trackedWallets) == 0 {
				log.Println("[ProcessOperation]public key not found in trovo user ecosystem")
				return
			}
			//Trovo user exists
			log.Printf("[ProcessOperation] found trovo user wallet %+v\n", trackedWallets)

			for _, t := range trackedWallets {

				if t.PublicKey == pmt.Funder {
					fromAlias = t.Alias
					fromName = t.Name
				}
				if t.PublicKey == pmt.Account {
					toAlias = t.Alias
					toName = t.Name
				}
			}

			e := paymentServices.SavePaymentHistory(pmt.Funder, fromAlias, fromName, pmt.Account, toAlias, toName, pmt.Transaction.Memo, "", os.Getenv("NATIVE_ASSET_CODE"), pmt.StartingBalance, pmt.TransactionHash, "PAYMENT", pmt.PT, pmt.ID, fmt.Sprintf("%v", pmt.Transaction.AccountSequence), pmt.LedgerCloseTime, db)
			if e != nil {
				log.Println("[ProcessOperation] unable to save SavePaymentHistory with create account operation:", e)
				return
			}
		} else {
			log.Println("[ProcessOperation] unable to find tracked wallets for create account due to:", dbFetchError)
		}

	} else if o.GetType() == "path_payment_strict_send" {
		//payment transaction.
		log.Println("found path payment strict send operation....beginning processing")
		pmt := interface{}(o).(operations.PathPaymentStrictSend)

		//send out to be saved to db
		var trackedWallets []paymentModels.TrackedWallet
		// dbFetchError := roachDB.Where(paymentModels.TrackedWallet{PublicKey: pmt.From}).Or(paymentModels.TrackedWallet{PublicKey: pmt.To}).Find(&trackedWallets).Error
		dbFetchError := roachDB.Where("(public_key = ? OR temp_public_key = ?) OR (public_key = ? OR temp_public_key = ?)", pmt.From, pmt.From, pmt.To, pmt.To).Find(&trackedWallets).Error
		var fromAlias, fromName, toAlias, toName string
		if dbFetchError == nil {
			if len(trackedWallets) == 0 {
				log.Println("[ProcessOperation]public key not found in trovo user ecosystem")
				return
			}
			//Trovo user exists
			log.Printf("[ProcessOperation] found trovo user wallet %+v\n", trackedWallets)

			for _, t := range trackedWallets {

				if t.PublicKey == pmt.From {
					fromAlias = t.Alias
					fromName = t.Name
				}
				if t.PublicKey == pmt.To {
					toAlias = t.Alias
					toName = t.Name
				}
			}
			var sourceAssetCode, destinationAssetCode string
			if pmt.SourceAssetIssuer == "" {
				sourceAssetCode = os.Getenv("NATIVE_ASSET_CODE")
			} else {
				sourceAssetCode = pmt.SourceAssetCode
			}
			if pmt.Issuer == "" {
				destinationAssetCode = os.Getenv("NATIVE_ASSET_CODE")
			} else {
				destinationAssetCode = pmt.Code
			}
			transactionType := swapTransactionType(pmt.From, pmt.SourceAssetIssuer, sourceAssetCode, pmt.To, pmt.Issuer, destinationAssetCode)
			e := paymentServices.SavePaymentHistory(pmt.From, fromAlias, fromName, pmt.To, toAlias, toName, pmt.Transaction.Memo, pmt.Issuer, destinationAssetCode, pmt.Amount, pmt.TransactionHash, transactionType, pmt.PT, pmt.ID, fmt.Sprintf("%v", pmt.Transaction.AccountSequence), pmt.LedgerCloseTime, db)
			if e != nil {
				log.Println("[ProcessOperation] unable to save SavePaymentHistory with strict payment send:", e)
				return
			}
		} else {
			log.Println("[ProcessOperation] unable to find tracked wallets for path payment strict send due to:", dbFetchError)
		}

	} else if o.GetType() == "path_payment" {
		//payment transaction.
		log.Println("found path payment operation....beginning processing")
		pmt := interface{}(o).(operations.PathPayment)

		//send out to be saved to db
		trackedWallets := make([]paymentModels.TrackedWallet, 0)
		// dbFetchError := roachDB.Where(paymentModels.TrackedWallet{PublicKey: pmt.From}).Or(paymentModels.TrackedWallet{PublicKey: pmt.To}).Find(&trackedWallets).Error
		dbFetchError := roachDB.Where("(public_key = ? OR temp_public_key = ?) OR (public_key = ? OR temp_public_key = ?)", pmt.From, pmt.From, pmt.To, pmt.To).Find(&trackedWallets).Error
		var fromAlias, fromName, toAlias, toName string
		if dbFetchError == nil {
			if len(trackedWallets) == 0 {
				log.Println("[ProcessOperation]public key not found in trovo user ecosystem")
				return
			}
			//Trovo user exists
			log.Printf("[ProcessOperation] found trovo user wallet %+v\n", trackedWallets)

			for _, t := range trackedWallets {

				if t.PublicKey == pmt.From {
					fromAlias = t.Alias
					fromName = t.Name
				}
				if t.PublicKey == pmt.To {
					toAlias = t.Alias
					toName = t.Name
				}
			}
			var sourceAssetCode, destinationAssetCode string
			if pmt.SourceAssetIssuer == "" {
				sourceAssetCode = os.Getenv("NATIVE_ASSET_CODE")
			} else {
				sourceAssetCode = pmt.SourceAssetCode
			}
			if pmt.Issuer == "" {
				destinationAssetCode = os.Getenv("NATIVE_ASSET_CODE")
			} else {
				destinationAssetCode = pmt.Code
			}
			transactionType := swapTransactionType(pmt.From, pmt.SourceAssetIssuer, sourceAssetCode, pmt.To, pmt.Issuer, destinationAssetCode)
			e := paymentServices.SavePaymentHistory(pmt.From, fromAlias, fromName, pmt.To, toAlias, toName, pmt.Transaction.Memo, pmt.Issuer, destinationAssetCode, pmt.Amount, pmt.TransactionHash, transactionType, pmt.PT, pmt.ID, fmt.Sprintf("%v", pmt.Transaction.AccountSequence), pmt.LedgerCloseTime, db)
			if e != nil {
				log.Println("[ProcessOperation] unable to save SavePaymentHistory with path payment:", e)
				return
			}
		} else {
			log.Println("[ProcessOperation] unable to find tracked wallets for path payment due to:", dbFetchError)
		}

	}
	{
		//ignore account merges

	}
	//save lastCursor
	// log.Println("saving last cursor", o.PagingToken())

}

func ProcessTrade(tr horizon.Trade, db, roachDB *gorm.DB, redisCache *cache.RedisCache, cancelTradeStream context.CancelFunc) {
	//when market maker sell counter = asset code, base is currency code, mm quantity is asset code (counter) quantity
	//if baseIsSeller ==true, then baseOfferId is market maker ID
	//if baseIsSeller ==false, then counterOfferId is market maker ID
	//when a market maker buy the counter = currency code, base is asset code, mm quantity is currency code (counter) quantity
	blockchainClient := network.GetBlockchainClient()
	var err error
	var mmBaseSignerKeyPair, mmCounterSignerKeyPair *keypair.Full
	var baseAsset, counterAsset txnbuild.Asset
	var baseMO, counterMO paymentModels.MarketOffer
	var ops []txnbuild.Operation = make([]txnbuild.Operation, 0)
	mmFeeWallet := keypair.MustParseFull(os.Getenv("MARKET_MAKING_FEE_WALLET"))
	var processBase, processCounter, signWithFeeWallet bool
	//if baseIsSeller == true, then baseOfferId is market maker ID
	if len(tr.BaseOfferID) > 0 {
		e := db.Where("blockchain_offer_id = ?", tr.BaseOfferID).First(&baseMO).Error

		if e != nil {
			if !errors.Is(e, gorm.ErrRecordNotFound) {
				log.Printf("[ProcessTrade] fatal error fetching offer with BC id %v from db: %v\nCancelling tradeStream now...\n", tr.BaseOfferID, e)
				cancelTradeStream()
				return
			}
			//could not locate the mm linked to trade stop process

		} else {
			//Enable processing the base offer
			processBase = true
			baseAsset = txnbuild.NativeAsset{}
			if len(tr.BaseAssetCode) > 0 {
				//custom asset
				baseAsset = txnbuild.CreditAsset{Code: tr.BaseAssetCode, Issuer: tr.BaseAssetIssuer}
			}
		}

	}
	if len(tr.CounterOfferID) > 0 {
		e := db.Where("blockchain_offer_id = ?", tr.CounterOfferID).First(&counterMO).Error

		if e != nil {
			if !errors.Is(e, gorm.ErrRecordNotFound) {
				log.Printf("[ProcessTrade] fatal error fetching counter offer with BC id %v from db: %v\nCancelling tradeStream now...\n", tr.CounterOfferID, e)
				cancelTradeStream()
				return
			}

		} else {
			//Enable processing the counter offer
			processCounter = true
			counterAsset = txnbuild.NativeAsset{}
			if len(tr.CounterAssetCode) > 0 {
				//custom asset
				counterAsset = txnbuild.CreditAsset{Code: tr.CounterAssetCode, Issuer: tr.CounterAssetIssuer}
			}
		}

	}
	// set transaction of DB
	dbTx := db.Begin()
	defer dbTx.Rollback()
	lastProcessedBaseCursor, lastProcessedCounterCursor := "0", "0"

	if counterMO.LastProcessedCursor != nil {
		lastProcessedCounterCursor = *counterMO.LastProcessedCursor
	}
	if baseMO.LastProcessedCursor != nil {
		lastProcessedBaseCursor = *baseMO.LastProcessedCursor
	}
	if processBase && tr.PT > lastProcessedBaseCursor {

		log.Println("[ProcessTrade]found base asset trade....beginning processing")

		/**
				when baseIsSeller==true {
		        //when market maker sell
				// base = asset code,
				//counter == currency code,
				//mm quantity is asset code (base amount) quantity
				}
				**/
		tradedAmount := decimal.RequireFromString(tr.BaseAmount)
		feeAmount := decimal.RequireFromString(baseMO.FeeValue)
		netQuantity := decimal.RequireFromString(baseMO.NetQuantity)
		//feeAmount/netQuantity gives hw much fee each netquantity gets
		feeSharefactor := feeAmount.Div(netQuantity)
		remainingQuantity := decimal.RequireFromString(baseMO.RemainingQuantity)
		remainingQuantityAfterTrade := remainingQuantity.Sub(tradedAmount)
		feeToTake := feeSharefactor.Mul(tradedAmount).Truncate(7)
		if remainingQuantityAfterTrade.IsZero() || remainingQuantityAfterTrade.IsNegative() {
			//order was filled.
			feeToTake = decimal.RequireFromString(baseMO.RemainingFeeValue)
			remainingQuantityAfterTrade = decimal.RequireFromString(baseMO.RemainingQuantity)
		}

		baseMO.RemainingQuantity = remainingQuantityAfterTrade.Truncate(7).String()
		baseMO.RemainingFeeValue = (decimal.RequireFromString(baseMO.RemainingFeeValue).Sub(feeToTake)).Truncate(7).String()
		baseMO.LastProcessedCursor = &tr.PT

		e := dbTx.Save(&baseMO).Error
		if e != nil {

			log.Printf("[ProcessTrade] fatal error saving mm offer with BC id %v on db: %v\nCancelling tradeStream now...\n", tr.CounterOfferID, e)
			cancelTradeStream()
			return

		}
		trovoAccountUsername := baseMO.SourceWalletAlias
		if strings.Contains(trovoAccountUsername, "_") {
			//subwallet.
			trovoAccountUsername = strings.Split(trovoAccountUsername, "_")[0]
		}
		mmBaseSignerKeyPair, err = bc.MarketMakingSignerKeypair(trovoAccountUsername, baseMO.MarketMakingWalletPublicKey)
		if err != nil {
			log.Printf("[ProcessTrade] fatal error generating signer key for mm offer with BC id %v on db: %v\nCancelling tradeStream now...\n", tr.CounterOfferID, err)
			cancelTradeStream()
			return
		}
		//prepare fee taking operation.

		if !baseAsset.IsNative() {
			//check if it has trustline to it and then create it.
			// _, mmAccountTrustsAsset, mmNativeAccountBalance, baseAssetCustomBalance, mmSourceAcountObject, _ := network.BlockchainAccountProperties(blockchainClient, baseMO.MarketMakingWalletPublicKey, baseAsset)
			_, mmAccountTrustsAsset, _, _, _, _ := network.BlockchainAccountProperties(blockchainClient, mmFeeWallet.Address(), baseAsset)

			if !mmAccountTrustsAsset {
				//establish trustline automatically
				ops = append(ops, &txnbuild.ChangeTrust{
					Line:          txnbuild.ChangeTrustAssetWrapper{Asset: baseAsset},
					Limit:         "900000000000",
					SourceAccount: mmFeeWallet.Address(),
				})
				signWithFeeWallet = true
			}

			ops = append(ops, &txnbuild.Payment{
				Asset:         baseAsset,
				Destination:   mmFeeWallet.Address(),
				Amount:        feeToTake.String(),
				SourceAccount: baseMO.MarketMakingWalletPublicKey,
			})

		} else {
			ops = append(ops, &txnbuild.Payment{
				Asset:         baseAsset,
				Destination:   mmFeeWallet.Address(),
				Amount:        feeToTake.String(),
				SourceAccount: baseMO.MarketMakingWalletPublicKey,
			})
		}

	}
	//processCounter trade
	if processCounter && tr.PT > lastProcessedCounterCursor {
		log.Println("[ProcessTrade]found counter asset trade....beginning processing")

		/**
				when baseIsSeller==true {
		        //when market maker sell
				// base = asset code,
				//counter == currency code,
				//mm quantity is asset code (base amount) quantity
				}
				**/
		tradedAmount := decimal.RequireFromString(tr.CounterAmount)
		feeAmount := decimal.RequireFromString(counterMO.FeeValue)
		netQuantity := decimal.RequireFromString(counterMO.NetQuantity)
		//feeAmount/netQuantity gives hw much fee each netquantity gets
		feeSharefactor := feeAmount.Div(netQuantity)
		remainingQuantity := decimal.RequireFromString(counterMO.RemainingQuantity)
		remainingQuantityAfterTrade := remainingQuantity.Sub(tradedAmount)
		feeToTake := feeSharefactor.Mul(tradedAmount).Truncate(7)
		if remainingQuantityAfterTrade.IsZero() || remainingQuantityAfterTrade.IsNegative() {
			//order was filled.
			feeToTake = decimal.RequireFromString(counterMO.RemainingFeeValue)
			remainingQuantityAfterTrade = decimal.RequireFromString(counterMO.RemainingQuantity)
		}

		counterMO.RemainingQuantity = remainingQuantityAfterTrade.Truncate(7).String()
		counterMO.RemainingFeeValue = (decimal.RequireFromString(counterMO.RemainingFeeValue).Sub(feeToTake)).Truncate(7).String()
		counterMO.LastProcessedCursor = &tr.PT
		e := dbTx.Save(&counterMO).Error
		if e != nil {

			log.Printf("[ProcessTrade] fatal error saving mm counter offer with BC id %v on db: %v\nCancelling tradeStream now...\n", tr.CounterOfferID, e)
			cancelTradeStream()
			return

		}
		trovoAccountUsername := counterMO.SourceWalletAlias
		if strings.Contains(trovoAccountUsername, "_") {
			//subwallet.
			trovoAccountUsername = strings.Split(trovoAccountUsername, "_")[0]
		}
		mmCounterSignerKeyPair, err = bc.MarketMakingSignerKeypair(trovoAccountUsername, counterMO.MarketMakingWalletPublicKey)
		if err != nil {
			log.Printf("[ProcessTrade] fatal error generating signer key for mm counter offer with BC id %v on db: %v\nCancelling tradeStream now...\n", tr.CounterOfferID, err)
			cancelTradeStream()
			return
		}

		//prepare fee taking operation.

		if !counterAsset.IsNative() {
			//check if it has trustline to it and then create it.
			// _, mmAccountTrustsAsset, mmNativeAccountBalance, baseAssetCustomBalance, mmSourceAcountObject, _ := network.BlockchainAccountProperties(blockchainClient, baseMO.MarketMakingWalletPublicKey, baseAsset)
			_, mmAccountTrustsAsset, _, _, _, _ := network.BlockchainAccountProperties(blockchainClient, mmFeeWallet.Address(), counterAsset)

			if !mmAccountTrustsAsset {
				//establish trustline automatically
				ops = append(ops, &txnbuild.ChangeTrust{
					Line:          txnbuild.ChangeTrustAssetWrapper{Asset: counterAsset},
					Limit:         "900000000000",
					SourceAccount: mmFeeWallet.Address(),
				})
				signWithFeeWallet = true
			}

			ops = append(ops, &txnbuild.Payment{
				Asset:         counterAsset,
				Destination:   mmFeeWallet.Address(),
				Amount:        feeToTake.String(),
				SourceAccount: counterMO.MarketMakingWalletPublicKey,
			})

		} else {
			ops = append(ops, &txnbuild.Payment{
				Asset:         counterAsset,
				Destination:   mmFeeWallet.Address(),
				Amount:        feeToTake.String(),
				SourceAccount: counterMO.MarketMakingWalletPublicKey,
			})
		}

	}

	if len(ops) > 0 {
		//operations exist
		chanKP := keypair.MustParseFull(os.Getenv("MM_FEE_COLLECTION_CHANNEL_ACCOUNT"))
		chanAccountExists, _, _, _, chanSourceAcountObject, _ := network.BlockchainAccountProperties(blockchainClient, chanKP.Address(), txnbuild.NativeAsset{})
		if !chanAccountExists {
			log.Println("[ProcessTrade] fatal error, fee collection channel account is not activated")
			cancelTradeStream()
			return
		}

		tx, err := txnbuild.NewTransaction(
			txnbuild.TransactionParams{
				SourceAccount:        chanSourceAcountObject,
				IncrementSequenceNum: true,
				Operations:           ops,
				BaseFee:              4000,
				Preconditions: txnbuild.Preconditions{
					TimeBounds: txnbuild.NewInfiniteTimeout(),
				},
				Memo: txnbuild.MemoText("MM Service Fee"),
			},
		)

		if err != nil {
			log.Println("[ProcessTrade]error constructing transaction", err)
			cancelTradeStream()
			return
		}

		if processBase {
			tx, err = tx.Sign(network.GetBlockchainNetworkPassPhrase(), mmBaseSignerKeyPair)

			if err != nil {
				log.Println("[ProcessTrade] error signing transaction with base account signer key ", err)
				cancelTradeStream()
				return
			}
		}

		if processCounter {
			tx, err = tx.Sign(network.GetBlockchainNetworkPassPhrase(), mmCounterSignerKeyPair)

			if err != nil {
				log.Println("[ProcessTrade] error signing transaction with counter account signer key ", err)
				cancelTradeStream()
				return
			}
		}

		tx, err = tx.Sign(network.GetBlockchainNetworkPassPhrase(), chanKP)

		if err != nil {
			log.Println("[ProcessTrade] error signing transaction with channel account signer key ", err)
			cancelTradeStream()
			return
		}

		if signWithFeeWallet {
			tx, err = tx.Sign(network.GetBlockchainNetworkPassPhrase(), mmFeeWallet)

			if err != nil {
				log.Println("[ProcessTrade] error signing transaction with mm fee wallet signer key ", err)
				cancelTradeStream()
				return
			}
		}

		ht, err := blockchainClient.SubmitTransaction(tx)

		if err != nil {
			log.Println("[ProcessTrade] error signing transaction with mm fee wallet signer key ", err)
			cancelTradeStream()
			return
		}

		dbTx.Commit()
		log.Println("[ProcessTrade] >>>>>>>>> >>>>>>>>>>> successfully  processed the MM fee.", ht.Hash)

	}
	enable := false
	if processBase && enable {
		//save to history
		//send out to be saved to db
		baseCode := os.Getenv("NATIVE_ASSET_CODE")
		// counterCode := os.Getenv("NATIVE_ASSET_CODE")
		if len(tr.BaseAssetCode) > 0 {
			baseCode = tr.BaseAssetCode
		}
		// if len(tr.CounterAssetCode) > 0 {
		// 	counterCode = tr.CounterAssetCode
		// }
		//check if public key exists within trovo user ecosystem
		trackedWallets := make([]paymentModels.TrackedWallet, 0)

		dbFetchError := roachDB.Where("(public_key = ? OR temp_public_key = ?) OR (public_key = ? OR temp_public_key = ?)", tr.BaseAccount, tr.BaseAccount, tr.CounterAccount, tr.CounterAccount).Find(&trackedWallets).Error
		var fromAlias, fromName, toAlias, toName string
		if dbFetchError == nil {
			if len(trackedWallets) == 0 {
				log.Println("[ProcessTrade] public key not found in trovo user ecosystem")
				return
			}
			//Trovo user exists
			log.Printf("[ProcessTrade] found trovo user wallet %+v\n", trackedWallets)
			for _, t := range trackedWallets {

				if t.PublicKey == tr.BaseAccount {
					fromAlias = t.Alias
					fromName = t.Name
				}
				if t.PublicKey == tr.CounterAccount {
					toAlias = t.Alias
					toName = t.Name
				}
			}

			paymentType := "TRADE"

			e := paymentServices.SavePaymentHistory(tr.BaseAccount, fromAlias, fromName, tr.CounterAccount, toAlias, toName, "", tr.BaseAssetIssuer, baseCode, tr.BaseAmount, tr.ID+tr.BaseOfferID, paymentType, tr.PT, tr.ID, fmt.Sprintf("%v", tr.PT), tr.LedgerCloseTime, db)
			if e != nil {
				log.Println("[ProcessTrade] unable to save trade history:", e)
				return
			}

			processCounter = false

		} else {
			log.Println("[ProcessTrade] unable to find tracked wallets for payment due to:", dbFetchError)
		}

	}

	if processCounter && enable {
		//save to history
		//send out to be saved to db
		// baseCode := os.Getenv("NATIVE_ASSET_CODE")
		counterCode := os.Getenv("NATIVE_ASSET_CODE")
		// if len(tr.BaseAssetCode) > 0 {
		// 	baseCode = tr.BaseAssetCode
		// }
		if len(tr.CounterAssetCode) > 0 {
			counterCode = tr.CounterAssetCode
		}
		//check if public key exists within trovo user ecosystem
		trackedWallets := make([]paymentModels.TrackedWallet, 0)

		dbFetchError := roachDB.Where("(public_key = ? OR temp_public_key = ?) OR (public_key = ? OR temp_public_key = ?)", tr.CounterAccount, tr.CounterAccount, tr.BaseAccount, tr.BaseAccount).Find(&trackedWallets).Error
		var fromAlias, fromName, toAlias, toName string
		if dbFetchError == nil {
			if len(trackedWallets) == 0 {
				log.Println("[ProcessTrade] public key not found in trovo user ecosystem")
				return
			}
			//Trovo user exists
			log.Printf("[ProcessTrade] found trovo user wallet %+v\n", trackedWallets)
			for _, t := range trackedWallets {

				if t.PublicKey == tr.CounterAccount {
					fromAlias = t.Alias
					fromName = t.Name
				}
				if t.PublicKey == tr.BaseAccount {
					toAlias = t.Alias
					toName = t.Name
				}
			}

			paymentType := "TRADE"

			e := paymentServices.SavePaymentHistory(tr.CounterAccount, fromAlias, fromName, tr.BaseAccount, toAlias, toName, "", tr.CounterAssetIssuer, counterCode, tr.CounterAmount, tr.ID+tr.CounterOfferID, paymentType, tr.PT, tr.ID+tr.CounterOfferID, fmt.Sprintf("%v%v", tr.PT, tr.CounterOfferID), tr.LedgerCloseTime, db)
			if e != nil {
				log.Println("[ProcessTrade] unable to save trade history:", e)
				return
			}
		} else {
			log.Println("[ProcessTrade] unable to find tracked wallets for payment due to:", dbFetchError)
		}
	}

}
