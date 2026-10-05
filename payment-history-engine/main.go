package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	_ "trovo-wallet-payment-history-engine/docs"

	"trovo-wallet-payment-history-engine/internal/cache"
	"trovo-wallet-payment-history-engine/internal/components/health"
	paymentModels "trovo-wallet-payment-history-engine/internal/components/payments/models"
	paymentServices "trovo-wallet-payment-history-engine/internal/components/payments/services"

	userModels "trovo-wallet-payment-history-engine/internal/components/users/models"
	db "trovo-wallet-payment-history-engine/internal/db"
	"trovo-wallet-payment-history-engine/internal/network"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// startStreaming and trackAddress are read and written from multiple goroutines,
// so they're stored as int32s and accessed exclusively through atomic ops (0 = false, 1 = true)
// instead of as plain bools, which would be a data race under concurrent access.
var startStreaming int32
var trackAddress int32

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

func setTrackAddress(v bool) {
	if v {
		atomic.StoreInt32(&trackAddress, 1)
		return
	}
	atomic.StoreInt32(&trackAddress, 0)
}

func isTrackAddress() bool {
	return atomic.LoadInt32(&trackAddress) == 1
}

// @title Payment History Engine
// @version 1.0
// @description Background worker that indexes on-chain Base (EVM) payment activity - native transfers and B20/ERC-20 Transfer events - for Trovo Wallet's tracked wallets, writing the results to its own RoachDB (CockroachDB) payment-history database. It has no request/response API of its own: the only HTTP surface is the liveness/readiness probe pair documented below, served for monitoring purposes.
// @host localhost:8080
// @BasePath /
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
		// BLOCKCHAIN_NETWORK_PASSPHRASE and BLOCKCHAIN_BASE_RESERVE are
		// vestigial on Base (see network.GetBlockchainNetworkPassPhrase /
		// GetBlockchainBaseReserve) - both handle being unset gracefully,
		// so they're no longer required to boot.
		requiredEnvironmentVariables := []string{"BASE_RPC_URL",
			"MNEMONIC_TEMP_ACCOUNTS", "NATIVE_ASSET_CODE", "ENABLE_CACHING",
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

		errMigrate = roachDB.Raw("DROP TABLE IF EXISTS tracked_addresses").Error
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
	errMigrate = roachDB.AutoMigrate(&paymentModels.TrackedAddress{})
	if errMigrate != nil {
		if !strings.Contains(errMigrate.Error(), "constraint") {
			log.Fatalf("Error migrating TrackedAddress model, error: %v", errMigrate)
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
				e := database.Where("created_at >= ? AND tracked = 0", startTrackingFrom).First(&userModels.UserWallet{}).Error
				if e != nil {
					log.Println("[TRACKUserWallet] error getting user wallet:", e)
					setStartStreaming(true)
					if errors.Is(e, gorm.ErrRecordNotFound) {
						setTrackAddress(true)
					}

					time.Sleep(2 * time.Minute)
					continue
				}

				result := database.Where("created_at >= ? AND tracked = 0", startTrackingFrom).FindInBatches(&userWallets, batchSize, func(tx *gorm.DB, batch int) error {
					for i, u := range userWallets {
						log.Printf("[TRACKUserWallet] processing user wallet %+v of %v\n", u, batch)
						trackError := paymentServices.TrackUserWallet(u, roachDB, database, isTrackAddress(), &redisCache)
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
						setTrackAddress(true)
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
			log.Println("##[TRACKAddresses] started routine to generate payment history of public keys")

			batchSize := 100
			var userAddresses []paymentModels.TrackedAddress
			for {
				e := roachDB.First(&paymentModels.TrackedAddress{}).Error
				if e != nil {
					log.Println("[TRACKAddresses] error getting public keys to process", e)
					if errors.Is(e, gorm.ErrRecordNotFound) {
						setTrackAddress(true)
					}
					time.Sleep(1 * time.Minute)
					continue
				}
				var wg sync.WaitGroup
				// at most backfillConcurrency() addresses at a time: one
				// goroutine per address at once spikes memory and RPC load
				// with many tracked addresses
				slots := make(chan struct{}, backfillConcurrency())
				result := roachDB.FindInBatches(&userAddresses, batchSize, func(tx *gorm.DB, batch int) error {
					for _, u := range userAddresses {
						wg.Add(1)
						slots <- struct{}{}
						log.Printf("[TRACKAddresses] processing user publicKey %v of %v\n", u.Address, batch)
						//spin off worker to process the payment history of the public key
						go func(address string) {
							defer func() { <-slots }()
							MonitorAddressPaymentStream(address, database, roachDB, &wg)
						}(u.Address)
					}

					return nil
				})
				if result.Error != nil {
					if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
						log.Println("[TRACKAddresses]()()()()()()()()()()()()()error occurred during batch processing:", result.Error.Error())

					}
					if errors.Is(e, gorm.ErrRecordNotFound) {
						setTrackAddress(true)
					}

				}
				//All waiting public keys finished. wait for them to complete generating their history
				wg.Wait()
				time.Sleep(1 * time.Minute)
				log.Printf("[TRACKAddresses] Total: %v, Errors:%v restarting process...\n", result.RowsAffected, result.Error)
				setTrackAddress(true)
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
	return nil

}

// GetAccountCursor returns the last processed block number for a single tracked public key's
// genesis backfill, or "0" (start from genesis) if it has never been recorded.
func GetAccountCursor(publicKey string, roachDB *gorm.DB) string {
	var mAccount paymentModels.MonitoredAccountCursor
	e := roachDB.Where("address = ?", publicKey).First(&mAccount).Error
	if e != nil {
		return "0"
	}
	return mAccount.LastCursor
}

// SaveAccountCursor persists the last processed block number for a single tracked public key,
// so its backfill can resume from where it left off instead of always restarting at genesis.
func SaveAccountCursor(publicKey, lastCursor string, roachDB *gorm.DB) error {
	var mAccount paymentModels.MonitoredAccountCursor
	e := roachDB.Where("address = ?", publicKey).First(&mAccount).Error
	if e != nil {
		if !errors.Is(e, gorm.ErrRecordNotFound) {
			return e
		}
		mAccount = paymentModels.MonitoredAccountCursor{Address: publicKey, LastCursor: lastCursor}
		return roachDB.Create(&mAccount).Error
	}
	if mAccount.LastCursor == lastCursor {
		return nil
	}
	mAccount.LastCursor = lastCursor
	return roachDB.Save(&mAccount).Error
}

// transferEventSig is keccak256("Transfer(address,address,uint256)") - the
// standard ERC-20/B20 Transfer event topic0, the Base equivalent of
// Horizon's "payment"/"path_payment" operation types for token movement.
var transferEventSig = crypto.Keccak256Hash([]byte("Transfer(address,address,uint256)"))

var tokenMetaABI abi.ABI

func init() {
	var err error
	tokenMetaABI, err = abi.JSON(strings.NewReader(`[
		{"constant":true,"inputs":[],"name":"symbol","outputs":[{"name":"","type":"string"}],"type":"function"},
		{"constant":true,"inputs":[],"name":"decimals","outputs":[{"name":"","type":"uint8"}],"type":"function"}
	]`))
	if err != nil {
		log.Panicf("invalid embedded token-metadata ABI: %v", err)
	}
}

type tokenMeta struct {
	symbol   string
	decimals int32
}

var tokenMetaCache sync.Map // string(lowercased contract address) -> tokenMeta

// getTokenMeta fetches (and caches) a B20/ERC20 token contract's symbol and
// decimals - Base tokens are self-describing via these calls the way a
// Stellar CreditAsset's Code was self-describing in the operation itself.
func getTokenMeta(client *ethclient.Client, contract string) tokenMeta {
	key := strings.ToLower(contract)
	if v, ok := tokenMetaCache.Load(key); ok {
		return v.(tokenMeta)
	}
	meta := tokenMeta{symbol: contract, decimals: 18}
	readSymbol, readDecimals := false, false
	addr := common.HexToAddress(contract)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if data, err := tokenMetaABI.Pack("symbol"); err == nil {
		if out, err := client.CallContract(ctx, ethereum.CallMsg{To: &addr, Data: data}, nil); err == nil {
			if unpacked, err := tokenMetaABI.Unpack("symbol", out); err == nil && len(unpacked) > 0 {
				if s, ok := unpacked[0].(string); ok && s != "" {
					meta.symbol = s
				}
				readSymbol = true
			}
		}
	}
	if data, err := tokenMetaABI.Pack("decimals"); err == nil {
		if out, err := client.CallContract(ctx, ethereum.CallMsg{To: &addr, Data: data}, nil); err == nil {
			if unpacked, err := tokenMetaABI.Unpack("decimals", out); err == nil && len(unpacked) > 0 {
				if d, ok := unpacked[0].(uint8); ok {
					meta.decimals = int32(d)
					readDecimals = true
				}
			}
		}
	}
	// only cache what was read: a fallback cached while the RPC was failing
	// would record this token's amounts with the wrong decimals for good
	if readSymbol && readDecimals {
		tokenMetaCache.Store(key, meta)
	}
	return meta
}

// swapTransactionType labels a path payment as a plain SWAP, or as a MINT/BURN TOKEN (SWAP ...)
// when one leg's issuer is also that leg's sender/receiver, mirroring the plain "payment"
// operation's MINT TOKEN / BURN TOKEN classification (an issuer sending its own asset is a mint,
// an issuer receiving its own asset is a burn) so mint/burn activity routed through a path
// payment isn't hidden behind a generic swap label. Currently unused: Base has no native DEX to
// source a path-payment-shaped swap event from (see MonitorTradeStream's doc comment) - kept for
// a future DEX/AMM-router integration that would want the same classification.
func swapTransactionType(from, sourceContractAddress, sourceAssetCode, to, destinationContractAddress, destinationAssetCode string) string {
	swapLabel := fmt.Sprintf("SWAP %s>%s", sourceAssetCode, destinationAssetCode)
	isMint := sourceContractAddress != "" && from == sourceContractAddress
	isBurn := destinationContractAddress != "" && to == destinationContractAddress

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

// lookupTrackedWallets finds the alias/name of from/to among this engine's
// tracked wallets - the same "is this address one of our users" gate the
// original Horizon-operation processing applied before writing history.
func lookupTrackedWallets(roachDB *gorm.DB, from, to string) (fromAlias, fromName, toAlias, toName string, found bool) {
	var trackedWallets []paymentModels.TrackedWallet
	dbFetchError := roachDB.Where("address = ? OR address = ?", from, to).Find(&trackedWallets).Error
	if dbFetchError != nil {
		log.Println("[lookupTrackedWallets] unable to find tracked wallets due to:", dbFetchError)
		return "", "", "", "", false
	}
	if len(trackedWallets) == 0 {
		return "", "", "", "", false
	}
	for _, t := range trackedWallets {
		if strings.EqualFold(t.Address, from) {
			fromAlias, fromName = t.Alias, t.Name
		}
		if strings.EqualFold(t.Address, to) {
			toAlias, toName = t.Alias, t.Name
		}
	}
	return fromAlias, fromName, toAlias, toName, true
}

// recordNativeTransfer records an ETH transfer touching a tracked wallet
// under id (unique per transfer) and seq (distinguishing transfers of one
// transaction).
func recordNativeTransfer(from, to string, weiAmount *big.Int, txHash string, blockNumber uint64, id, seq string, blockTime uint64, db, roachDB *gorm.DB) {
	fromAlias, fromName, toAlias, toName, found := lookupTrackedWallets(roachDB, from, to)
	if !found {
		return
	}
	log.Printf("[recordNativeTransfer] found trovo user wallet(s) for tx %v\n", txHash)

	amount := decimal.NewFromBigInt(weiAmount, -18)
	pt := fmt.Sprintf("%020d", blockNumber)

	e := paymentServices.SavePaymentHistory(from, fromAlias, fromName, to, toAlias, toName, "", "", os.Getenv("NATIVE_ASSET_CODE"), amount.String(), paymentModels.NetworkBase, txHash, "PAYMENT", pt, id, seq, time.Unix(int64(blockTime), 0), db)
	if e != nil {
		log.Println("[recordNativeTransfer] unable to save SavePaymentHistory:", e)
	}
}

// processB20TransferLog records one decoded B20/ERC-20 Transfer event log -
// the Base equivalent of Horizon's "payment"/"path_payment" operations for
// a custom (non-native) asset. Mint/burn is detected the idiomatic
// ERC-20 way (Transfer from/to the zero address), replacing Stellar's
// issuer-address heuristic (an issuer account sending/receiving its own
// asset) since B20 tokens have no separate "issuer account" distinct from
// the token contract itself.
func processB20TransferLog(client *ethclient.Client, lg types.Log, blockTime uint64, db, roachDB *gorm.DB) {
	if len(lg.Topics) < 3 || len(lg.Data) < 32 {
		return // not a standard indexed Transfer(address,address,uint256) log
	}
	from := common.HexToAddress(lg.Topics[1].Hex()).Hex()
	to := common.HexToAddress(lg.Topics[2].Hex()).Hex()
	rawAmount := new(big.Int).SetBytes(lg.Data[len(lg.Data)-32:])
	tokenContract := lg.Address.Hex()

	fromAlias, fromName, toAlias, toName, found := lookupTrackedWallets(roachDB, from, to)
	if !found {
		return
	}
	log.Printf("[processB20TransferLog] found trovo user wallet(s) for tx %v log %v\n", lg.TxHash.Hex(), lg.Index)

	meta := getTokenMeta(client, tokenContract)
	amount := decimal.NewFromBigInt(rawAmount, -meta.decimals)

	paymentType := "PAYMENT"
	if from == common.HexToAddress("0x0").Hex() {
		paymentType = "MINT TOKEN"
	} else if to == common.HexToAddress("0x0").Hex() {
		paymentType = "BURN TOKEN"
	}

	pt := fmt.Sprintf("%020d-%010d", lg.BlockNumber, lg.Index)
	id := fmt.Sprintf("%s-%d", lg.TxHash.Hex(), lg.Index)

	e := paymentServices.SavePaymentHistory(from, fromAlias, fromName, to, toAlias, toName, "", tokenContract, meta.symbol, amount.String(), paymentModels.NetworkBase, lg.TxHash.Hex(), paymentType, pt, id, fmt.Sprintf("%d", lg.TxIndex), time.Unix(int64(blockTime), 0), db)
	if e != nil {
		log.Println("[processB20TransferLog] unable to save SavePaymentHistory:", e)
	}
}

// ProcessBlock scans one Base block for native-currency transfers (plain
// value-carrying transactions) and B20/ERC-20 Transfer event logs - the
// Base equivalent of Horizon's per-operation stream, replayed one block at
// a time since Base has no account-agnostic operation feed to subscribe to
// (see MonitorPaymentStream's doc comment).
//
// It reads the block and its Transfer logs before recording anything, and
// returns an error when either read fails, so the caller retries the block
// instead of moving past it (which would lose its payments).
func ProcessBlock(client *ethclient.Client, blockNumber uint64, db, roachDB *gorm.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	block, err := client.BlockByNumber(ctx, new(big.Int).SetUint64(blockNumber))
	if err != nil {
		return fmt.Errorf("fetching block %d: %w", blockNumber, err)
	}
	logs, err := client.FilterLogs(ctx, ethereum.FilterQuery{
		FromBlock: new(big.Int).SetUint64(blockNumber),
		ToBlock:   new(big.Int).SetUint64(blockNumber),
		Topics:    [][]common.Hash{{transferEventSig, userOperationEventSig, safeReceivedSig}},
	})
	if err != nil {
		return fmt.Errorf("fetching logs for block %d: %w", blockNumber, err)
	}
	byTx := map[common.Hash][]types.Log{}
	for _, lg := range logs {
		byTx[lg.TxHash] = append(byTx[lg.TxHash], lg)
	}

	signer := types.LatestSignerForChainID(network.GetBlockchainChainID())
	entryPoint := entryPointAddress()
	for _, tx := range block.Transactions() {
		for _, t := range nativeTransfersOfTx(tx, byTx[tx.Hash()], signer, entryPoint) {
			recordNativeTransfer(t.From.Hex(), t.To.Hex(), t.Value, tx.Hash().Hex(), blockNumber, t.ID, t.Seq, block.Time(), db, roachDB)
		}
	}
	for _, lg := range logs {
		if len(lg.Topics) > 0 && lg.Topics[0] == transferEventSig {
			processB20TransferLog(client, lg, block.Time(), db, roachDB)
		}
	}
	return nil
}

// MonitorPaymentStream watches every new Base block for native transfers
// and B20 Transfer events touching a tracked wallet, and records payment
// history for them. This replaces Horizon's account-agnostic
// client.StreamOperations/StreamPayments SSE stream: Base has no equivalent
// feed to subscribe to, so this polls sequentially block-by-block instead
// (the cursor is a block number, not a Horizon paging token).
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
	lastCursor := GetLastCursor(roachDB)
	var startBlock uint64
	if n, err := strconv.ParseUint(lastCursor, 10, 64); err == nil && lastCursor != "0" {
		startBlock = n + 1
		log.Printf("[MonitorPaymentStream] Starting monitoring from block[%v]\n", startBlock)
	} else {
		log.Println("[MonitorPaymentStream] Starting monitoring from current state of blockchain")
	}

	for {
		bctx, bcancel := context.WithTimeout(context.Background(), 30*time.Second)
		latest, err := client.BlockNumber(bctx)
		bcancel()
		if err != nil {
			log.Printf("[MonitorPaymentStream] error fetching latest block: %v\n", err)
			time.Sleep(5 * time.Second)
			continue
		}
		if startBlock == 0 {
			// No saved cursor - start from the current chain tip rather than
			// replaying the whole chain's history.
			startBlock = latest
		}
		for b := startBlock; b <= latest; b++ {
			if err := ProcessBlock(client, b, db, roachDB); err != nil {
				// retried on the next pass; the cursor stays before it
				log.Printf("[MonitorPaymentStream] %v - retrying\n", err)
				break
			}
			cursor := fmt.Sprintf("%d", b)
			health.RecordOperation(cursor)
			if e := SaveLastCursor(cursor, roachDB); e != nil {
				log.Printf("[MonitorPaymentStream] unable to save last cursor: %v\n", e)
			}
			startBlock = b + 1
		}
		time.Sleep(4 * time.Second)
	}
}

// backfillConcurrency is how many tracked addresses are backfilled at once
// (TRACK_ADDRESS_CONCURRENCY, default 8).
func backfillConcurrency() int {
	if n, err := strconv.Atoi(strings.TrimSpace(os.Getenv("TRACK_ADDRESS_CONCURRENCY"))); err == nil && n > 0 {
		return n
	}
	return 8
}

// backfillChunkBlocks bounds each eth_getLogs range query - many public RPC
// providers cap how many blocks a single filter query may span.
const backfillChunkBlocks = 5000

// MonitorAddressPaymentStream backfills one tracked public key's B20
// token-transfer history from its last saved cursor (or genesis) up to the
// current chain tip, using indexed Transfer-event topic filtering
// (eth_getLogs with the address in the "from" or "to" topic position) so it
// does not need to scan every block. Native-currency history for a single
// historical account has no equally efficient plain-JSON-RPC equivalent
// (native transfers emit no logs to filter by address) - that needs a
// block-indexing service (e.g. an Etherscan/Blockscout-style API). Trovo
// wallets are Safes, though, and a Safe's ETH does leave logs: its own
// SafeReceived events and the EntryPoint's UserOperationEvents for its
// sends, so their ETH history is backfilled too (backfillNativeOfSafe) -
// all but ETH sent to a wallet before it was deployed. This is the Base
// equivalent of Horizon's per-account client.StreamPayments(ForAccount:...).
func MonitorAddressPaymentStream(publicKey string, db, roachDB *gorm.DB, wg *sync.WaitGroup) {
	defer wg.Done()
	checkExists := roachDB.Where("address = ?", publicKey).First(&paymentModels.TrackedWallet{}).Error
	if checkExists != nil {
		log.Printf("[MonitorAddressPaymentStream] aborting because %v could not be found in tracked wallets table:", checkExists)
		return
	}
	if !common.IsHexAddress(publicKey) {
		log.Printf("[MonitorAddressPaymentStream] aborting, %v is not a valid Base address\n", publicKey)
		return
	}
	client := network.GetBlockchainClient()

	lastAccountCursor := GetAccountCursor(publicKey, roachDB)
	var startBlock uint64
	if n, err := strconv.ParseUint(lastAccountCursor, 10, 64); err == nil && lastAccountCursor != "0" {
		startBlock = n + 1
	}
	log.Printf("[MonitorAddressPaymentStream] Starting monitoring for %v from block[%v]\n", publicKey, startBlock)

	bctx, bcancel := context.WithTimeout(context.Background(), 30*time.Second)
	latest, err := client.BlockNumber(bctx)
	bcancel()
	if err != nil {
		log.Printf("[MonitorAddressPaymentStream] error fetching latest block for %v: %v\n", publicKey, err)
		return
	}

	addrTopic := common.BytesToHash(common.HexToAddress(publicKey).Bytes())

	for from := startBlock; from <= latest; from += backfillChunkBlocks {
		to := from + backfillChunkBlocks - 1
		if to > latest {
			to = latest
		}
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		logsFrom, err1 := client.FilterLogs(ctx, ethereum.FilterQuery{
			FromBlock: new(big.Int).SetUint64(from),
			ToBlock:   new(big.Int).SetUint64(to),
			Topics:    [][]common.Hash{{transferEventSig}, {addrTopic}},
		})
		logsTo, err2 := client.FilterLogs(ctx, ethereum.FilterQuery{
			FromBlock: new(big.Int).SetUint64(from),
			ToBlock:   new(big.Int).SetUint64(to),
			Topics:    [][]common.Hash{{transferEventSig}, nil, {addrTopic}},
		})
		cancel()
		if err1 != nil || err2 != nil {
			log.Printf("[MonitorAddressPaymentStream] error fetching logs for %v in range [%v,%v]: %v / %v\n", publicKey, from, to, err1, err2)
			return
		}

		seen := map[string]bool{}
		blockTimes := map[uint64]uint64{}
		for _, lg := range append(logsFrom, logsTo...) {
			key := fmt.Sprintf("%s-%d", lg.TxHash.Hex(), lg.Index)
			if seen[key] {
				continue
			}
			seen[key] = true

			ts, ok := blockTimes[lg.BlockNumber]
			if !ok {
				hctx, hcancel := context.WithTimeout(context.Background(), 15*time.Second)
				header, errH := client.HeaderByNumber(hctx, new(big.Int).SetUint64(lg.BlockNumber))
				hcancel()
				if errH != nil {
					// stop here: the cursor stays before this chunk, so the
					// next pass retries it
					log.Printf("[MonitorAddressPaymentStream] error fetching block %v header for %v: %v\n", lg.BlockNumber, publicKey, errH)
					return
				}
				ts = header.Time
				blockTimes[lg.BlockNumber] = ts
			}
			processB20TransferLog(client, lg, ts, db, roachDB)
		}
		// the address's ETH, if it is a Safe: what it received (its
		// SafeReceived events) and sent (its user operations, from the
		// EntryPoint's UserOperationEvent indexed by sender)
		if err := backfillNativeOfSafe(client, common.HexToAddress(publicKey), from, to, blockTimes, db, roachDB); err != nil {
			log.Printf("[MonitorAddressPaymentStream] error reading ETH history of %v in range [%v,%v]: %v\n", publicKey, from, to, err)
			return
		}

		//persist progress so the next periodic pass over tracked_addresses resumes
		//here instead of re-scanning from genesis.
		if e := SaveAccountCursor(publicKey, fmt.Sprintf("%d", to), roachDB); e != nil {
			log.Printf("[MonitorAddressPaymentStream] unable to save account cursor for %v: %v\n", publicKey, e)
		}
	}

	log.Println("#####[MonitorAddressPaymentStream]...finished backfill pass for", publicKey)
}

// MonitorTradeStream watched Horizon's global trade stream
// (client.StreamTrades) and, for each trade filling one of this app's
// market-making offers, built and submitted a fee-collection transaction
// (see the former ProcessTrade). On Base, market-making offers live on
// TrovoOfferBook and their fills are indexed by app-backend
// (internal/offerbook); market-making fees are no longer charged, so there
// is no fee to collect when an offer is filled, and this stays idle.
func MonitorTradeStream(db, roachDB *gorm.DB, redisCache *cache.RedisCache) {
	checkExists := db.Where("CAST(remaining_quantity AS REAL) > ? AND canceled = 0", 0).First(&paymentModels.MarketOffer{}).Error

	if checkExists != nil {
		if !errors.Is(checkExists, gorm.ErrRecordNotFound) {
			log.Println("[MonitorTradeStream]unknown error while checking market offers:", checkExists)
		}
		log.Println("[MonitorTradeStream] exited because no market offers exists:", checkExists)
		return
	}
	log.Println("[MonitorTradeStream] Base has no trade stream to watch yet (see doc comment) - idling.")
	time.Sleep(5 * time.Minute)
}
