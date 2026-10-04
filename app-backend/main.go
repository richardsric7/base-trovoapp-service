package main

import (
	_ "trovo-wallet-api/docs"

	"trovo-wallet-api/internal/basetxn"
	cache "trovo-wallet-api/internal/cache"
	paymentModels "trovo-wallet-api/internal/components/payments/models"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/evmkeypair"
	"trovo-wallet-api/internal/gnosissafe"
	"trovo-wallet-api/internal/network"
	pns "trovo-wallet-api/internal/pns"
	"trovo-wallet-api/internal/sharedconfig"

	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	msc "trovo-wallet-api/internal/components/announcements/controllers"
	callbacks "trovo-wallet-api/internal/components/callbacks/controllers"
	p2p "trovo-wallet-api/internal/components/p2p/controllers"
	payments "trovo-wallet-api/internal/components/payments/controllers"
	rates "trovo-wallet-api/internal/components/rates/controllers"
	root "trovo-wallet-api/internal/components/root/controllers"
	serviceLinks "trovo-wallet-api/internal/components/servicelinks/controllers"
	swaps "trovo-wallet-api/internal/components/swaps/controllers"
	users "trovo-wallet-api/internal/components/users/controllers"
	userModels "trovo-wallet-api/internal/components/users/models"
	userServices "trovo-wallet-api/internal/components/users/services"
	db "trovo-wallet-api/internal/db"
	dl "trovo-wallet-api/internal/dynamiclinks"
	m "trovo-wallet-api/internal/mail"
	"trovo-wallet-api/internal/middleware"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/shopspring/decimal"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// @title Trovo Wallet API
// @version 1.0
// @description Core wallet API for the Trovo platform — user accounts, crypto assets, deposits/withdrawals, swaps.
// @host localhost:8080
// @BasePath /
//
// @securityDefinitions.apikey SignatureAuth
// @in header
// @name X-TW-SIGNATURE
// @description Request-signing scheme used by app-web and app-mobile for end-user requests. The caller sends four headers together: X-TW-PUBLIC-KEY (the wallet address making the call), X-TW-SIGNER (the address whose private key signed the request), X-TW-TIMESTAMP (a fresh timestamp), and X-TW-SIGNATURE (a signature over signer+timestamp, and for POST/PUT/PATCH over the request body/URI). There is no bearer token: every request is independently signed. See internal/middleware/authentication_middleware.go.
//
// @securityDefinitions.apikey ServiceLinkApiKey
// @in header
// @name X-TW-SERVICE-LINK-API-KEY
// @description Static per-partner API key used by white-label "service link" integrations (and by tm-api for a few admin/server-to-server calls). Issued when a service link is provisioned; see internal/components/servicelinks. See internal/middleware/api_key_middleware.go.
//
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description JWT bearer token used by the internal "trovo-manager" (tm-api admin) routes. Send as "Bearer <token>". See internal/middleware/authentication_middleware.go's JwtTokenAuthMiddleware.
func main() {

	//setup environment variables
	errEnv := godotenv.Load()
	if errEnv != nil {
		path, _ := os.Getwd()
		log.Printf("could not find or load any .env file from %v...skipping...\n", path)
	}

	// MIGRATE_ONLY: this process is a one-shot migrator run by CI BEFORE the app
	// is deployed, while the previous version keeps serving users. It opens only
	// the databases, runs the migrations, and exits — deliberately skipping the
	// ~50 required-env-var check below, which the app needs to serve but a
	// migrator does not. Exits non-zero on any failure so CI blocks the deploy.
	migrateOnly := os.Getenv("MIGRATE_ONLY") == "1"

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
		// The CI migrator only owns the wallet's own Postgres schema. RoachDB is
		// the payment-history service's database, so the migrator neither opens
		// nor migrates it — that stays a boot-time concern on the server, exactly
		// as before. Skipping the open matters: OpenRoachDB log.Fatal-s internally,
		// and with no CDB_CONNECTION_STRING it fails as `host=/tmp user=root`.
		if !migrateOnly {
			var err error
			roachDB, err = db.OpenRoachDB()

			if err != nil {
				log.Fatalf("[main]Error opening RoachDB %s", err)
				return
			}
		}
	}

	//check other required environment variables.

	// A one-shot migrator only needs DB connectivity, so skip the serving app's
	// required-env-var wall. Without this the migrator would hit `exit = true`
	// and `return` (exit code 0) BEFORE reaching the migrations below — a green
	// CI step that migrated nothing.
	if !migrateOnly {
		exit := false
		// BLOCKCHAIN_NETWORK_PASSPHRASE and BLOCKCHAIN_BASE_RESERVE are
		// vestigial on Base (see network.GetBlockchainNetworkPassPhrase /
		// GetBlockchainBaseReserve) - both handle being unset gracefully,
		// so they're no longer required to boot.
		requiredEnvironmentVariables := []string{"BASE_RPC_URL",
			"MAILGUN_PRIVATE_API_KEY", "CDB_CONNECTION_STRING",
			"IPAPI_KEY", "IPAPI_HOST", "VERIFICATION_CODE_SALT", "ENABLE_EMAIL_VALIDATION", "ENABLE_CACHING", "DEFAULT_ASSET_IMAGE_URL",
			"REDIS_HOST", "REDIS_PORT", "DYNAMIC_LINKS_DOMAIN_PREFIX", "DYNAMIC_LINKS_ANDROID_PACKAGE_NAME",
			"DYNAMIC_LINKS_IOS_BUNDLE_ID", "DYNAMIC_LINKS_FALLBACK_BASE_URL", "MAILGUN_DOMAIN", "NATIVE_ASSET_IMAGE_URL",
			"GC", "GOOGLE_PROJECT_ID",
			"NATIVE_ASSET_CODE",
			"CHANNEL_ACCOUNTS", "WALLET_DOMAIN",
			"MNEMONIC_BULK_PAYMENT", "BULK_PAYMENT_SALT", "ENCODER_SALT", "MARKET_MAKING_SALT",
			"MNEMONIC_MARKET_MAKING", "CHECK_CHANNEL_ACCOUNT_BALANCE",
			"JWT_ACCESS_SECRET", "JWT_TOKEN_EXPIRY", "JWT_REFRESH_TOKEN_EXPIRY",
			"DOLLAR_ASSET", "MARKET_MAKING_FEE_ENABLED", "SWAP_FEE_ENABLED",
			"BLOCKCHAIN_DATA_CACHE_LIFETIME", "VAT_WALLET", "SUBWALLET_CREATION_FEE_WALLET", "TOKENIZATION_APPLICATION_FEE_WALLET",
			"TOKENIZATION_FEE_WALLET", "CLOSED_GROUP_FEE_WALLET", "ACCOUNT_RECOVERY_FEE_WALLET", "PATRON_FEE_WALLET",
			"SHARED_ACCESS_PAYMENT_FEE_WALLET", "PAYMENT_FEE_WALLET", "SWAP_FEE_WALLET", "CNGN_PRICE_API_URL",
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
		if os.Getenv("ENABLE_CRYPTO_DEPOSIT_MINTING") == "1" && len(os.Getenv("CRYPTO_DEPOSIT_MINTING_INITIATOR_PUBLIC_KEY")) != 42 {
			log.Println("CRYPTO_DEPOSIT_MINTING_INITIATOR_PUBLIC_KEY environment variable is required when ENABLE_CRYPTO_DEPOSIT_MINTING is set to 1")

			exit = true
		}
		// if os.Getenv("MARKET_MAKING_FEE_ENABLED") == "1" && len(os.Getenv("MARKET_MAKING_FEE_WALLET")) != 42 {
		// 	log.Println("MARKET_MAKING_FEE_WALLET environment variable is required when MARKET_MAKING_FEE_ENABLED is set to 1")

		// 	exit = true
		// }

		if os.Getenv("ENABLE_NAIRA_ASSET_BY_DEFAULT") == "1" && len(os.Getenv("NAIRA_ASSET")) < 60 {
			log.Println("NAIRA_ASSET environment variable is required when ENABLE_NAIRA_ASSET_BY_DEFAULT is set to 1")

			exit = true
		}
		// if os.Getenv("MARKET_MAKING_FEE_ENABLED") == "1" {
		// 	_, err := keypair.ParseFull(os.Getenv("MARKET_MAKING_FEE_WALLET"))
		// 	if err != nil {
		// 		log.Println("MARKET_MAKING_FEE_WALLET  is invalid wallet secret key")

		// 		exit = true
		// 	}
		// }
		// if os.Getenv("SWAP_FEE_ENABLED") == "1" && len(os.Getenv("SWAP_FEE_WALLET")) != 42 {
		// 	log.Println("SWAP_FEE_WALLET environment variable is required when SWAP_FEE_ENABLED is set to 1")

		// 	exit = true
		// }
		// if os.Getenv("SWAP_FEE_ENABLED") == "1" {
		// 	_, err := keypair.ParseFull(os.Getenv("SWAP_FEE_WALLET"))
		// 	if err != nil {
		// 		log.Println("SWAP_FEE_WALLET is invalid wallet secret key")

		// 		exit = true
		// 	}
		// }

		if exit {
			// Non-zero: a missing required env var is a startup failure, and an
			// exit code of 0 would read as success to orchestrators and CI.
			os.Exit(1)
		}

		if os.Getenv("ENABLE_EMAIL_NOTIFICATIONS") == "" {
			log.Println("ENV variable ENABLE_EMAIL_NOTIFICATIONS is not set.")
		}

		if os.Getenv("BLOCKCHAIN_DATA_CACHE_LIFETIME") == "" {
			log.Println("ENV variable BLOCKCHAIN_DATA_CACHE_LIFETIME is not set. Using 3yrs by default")
			os.Setenv("BLOCKCHAIN_DATA_CACHE_LIFETIME", "94608000")
		}

		if os.Getenv("ACCOUNT_DELETION_REQUEST_TEMPLATE") == "" {
			log.Println("ENV variable ACCOUNT_DELETION_REQUEST_TEMPLATE is not set. defaulting to account-deletion-request-template")
		}

		if os.Getenv("ACCOUNT_DELETION_EMAIL_SUBJECT") == "" {
			log.Println("ENV variable ACCOUNT_DELETION_EMAIL_SUBJECT is not set. defaulting to 'TrovoApp Account Deletion Request'")
		}
		if os.Getenv("ACCOUNT_DELETION_DAYS") == "" {
			log.Println("ENV variable ACCOUNT_DELETION_DAYS is not set. defaulting to '30' days")
		}
		if os.Getenv("SUPPORT_EMAIL") == "" {
			log.Println("ENV variable SUPPORT_EMAIL is not set.")
		}

		if os.Getenv("SHORT_LINKS_BASE_URL") == "" {
			log.Println("ENV variable SHORT_LINKS_BASE_URL is not set. Using default https://trovo.app")
		}

		// Wires internal/network's WalletAssetAuthorization DB handle and
		// basetxn's default Builder (see network.SetDB's doc comment: "Called
		// once at startup"). Without this, IsWalletAuthorizedForAsset silently
		// no-ops (authDB stays nil) and every basetxn.Transaction.Sign call
		// across the codebase fails with "no Builder configured". Kept inside
		// this !migrateOnly block since it dials BASE_RPC_URL, which a
		// migrate-only run isn't guaranteed to have set.
		network.SetDB(database)
	}
	log.Println("starting migration")
	//migrate DB models if any
	db.MigrateDB(database)

	// RoachDB migrations are left exactly as they were originally: ungated by
	// DB_AUTOMIGRATE, run on every server boot. They are two small tables, so
	// they cost ~nothing at startup. Only the CI migrator skips them, because it
	// never opened this database.
	if !migrateOnly {
		errMigrate := roachDB.AutoMigrate(&paymentModels.TrackedWallet{})
		if errMigrate != nil {
			if !strings.Contains(errMigrate.Error(), "constraint") {
				log.Fatalf("Error migrating TrackedWallet model, error: %v", errMigrate)
			}
		}
		log.Println("migrating tracked wallet done...")
		errMigrate = roachDB.AutoMigrate(&paymentModels.TrackedAddress{})
		if errMigrate != nil {
			if !strings.Contains(errMigrate.Error(), "constraint") {
				log.Fatalf("Error migrating TrackedAddress model, error: %v", errMigrate)
			}
		}
		log.Println("migrating tracked public key done...")
	}

	//setup redis
	log.Println("migration done...")

	// Migrations are done. Exit 0 so the CI migrate step passes and the deploy
	// proceeds. Any migration failure above already log.Fatal-ed (non-zero exit)
	// and blocks the deploy, leaving the previous version serving. The serving
	// app runs WITHOUT MIGRATE_ONLY and with DB_AUTOMIGRATE=0, so it never
	// migrates and boots straight into serving traffic.
	if migrateOnly {
		log.Println("MIGRATE_ONLY=1 set — migrations complete, exiting without starting the server")
		os.Exit(0)
	}
	enableCaching := false

	if os.Getenv("ENABLE_CACHING") == "1" {
		enableCaching = true
	}

	var redisCli *redis.Client = nil

	// if enableCaching {
	// 	redisCli = redis.NewClient(&redis.Options{
	// 		Addr:     fmt.Sprintf("%v:%v", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")), // host:port of the redis server
	// 		Password: os.Getenv("REDIS_PASSWORD"),                                            // no password set
	// 		DB:       0,                                                                      // use default DB
	// 		TLSConfig: &tls.Config{
	// 			InsecureSkipVerify: false,
	// 		},
	// 	})
	// }
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
		// clear cache
		cacheKeyInfo := "curatedAssets_"
		redisCache.DeleteFromCache(cacheKeyInfo)
		// clear cache
		cacheKey := "[GET] /v1/rates"
		redisCache.InvalidateCachedHttpResponse(cacheKey)
		redisCache.DeleteFromCache(cacheKey)
	}

	//global config
	pnsContext := context.Background()
	pnsClient, _, err := pns.GetFirebaseMessagingClient(pnsContext)
	if err != nil {
		log.Fatalln("Unable to initialize Firebase messaging client:", err)
	}
	storageContext := context.Background()
	storageClient, _, err := pns.GetFirebaseStorageClient(storageContext)
	if err != nil {
		log.Fatalln("Unable to initialize Firebase storage client:", err)
	}

	var globalConfig = sharedconfig.GlobalConfig{
		PNSContext:             pnsContext,
		RedisCache:             &redisCache,
		DB:                     database,
		PushNotificationClient: pnsClient,
		RoachDB:                roachDB,
		BantuExpansionClient:   network.GetBlockchainClient(),
		BantuNetworkPassphrase: network.GetBlockchainNetworkPassPhrase(),
		FirebaseStorageUploader: &sharedconfig.ClientUploader{
			Client:     storageClient,
			ProjectID:  os.Getenv("GOOGLE_PROJECT_ID"),
			BucketName: os.Getenv("STORAGE_BUCKET_NAME"),
			UploadPath: os.Getenv("STORAGE_BUCKET_NAME"),
		},
	}
	if bucketName := strings.TrimSpace(os.Getenv("STAKEHOLDER_DOCUMENTS_BUCKET_NAME")); bucketName != "" {
		globalConfig.StakeholderDocumentStorage = sharedconfig.NewGCSPrivateDocumentStorage(storageClient, bucketName, "stakeholder-documents")
	} else {
		log.Println("STAKEHOLDER_DOCUMENTS_BUCKET_NAME is not configured; stakeholder document storage endpoints will return 503")
	}
	// Starts the shared Redis subscription P2P (and any future) live user
	// event streams relay through - see sharedconfig/realtime.go. A no-op
	// if Redis is disabled.
	globalConfig.StartUserStreamRelay()
	// platform Safes and signing keys are used by one instance at a time
	gnosissafe.KeyLock = func(ctx context.Context, name string, fn func() error) error {
		return sharedconfig.WithKeyLock(globalConfig.DB, "key:"+name, 3*time.Minute, fn)
	}
	gnosissafe.NonceSource = func(ctx context.Context, client *ethclient.Client, from common.Address) (uint64, func(), error) {
		return sharedconfig.NextNonceLocked(ctx, globalConfig.DB, client, from)
	}
	// legacy transactions signed now and sent later reserve their nonce
	network.ReserveNonce = func(ctx context.Context, from common.Address) (uint64, error) {
		return sharedconfig.ReserveNonce(ctx, globalConfig.DB, globalConfig.BantuExpansionClient, from)
	}
	{

		//update referral links for people with no referral link
		go func() {
			log.Println("##[REFLINKROUTINE] started routine to update referral links for people with no referral link")

			batchSize := 1
			var usersWithNoRefLinks []userModels.User

			for {
				// one instance at a time
				sharedconfig.WithSingletonLock(&globalConfig, "referral-links", 2*time.Minute, func() {
					e := database.Where("referral_qr_code is null AND suspended = ?", 0).First(&userModels.User{}).Error
					if e != nil {
						return
					}

					result := database.Where("referral_qr_code is null AND suspended = ?", 0).FindInBatches(&usersWithNoRefLinks, batchSize, func(tx *gorm.DB, batch int) error {
						for i, u := range usersWithNoRefLinks {
							rld, errLink := dl.GenerateReferralLinkWithStaticURL(u.Username, &globalConfig)
							if errLink != nil {
								continue
							}
							//link retrieved
							usersWithNoRefLinks[i].ReferralLink = &rld.DynamicLink
							usersWithNoRefLinks[i].ReferralQrCode = &rld.QRCode
						}

						e := database.Omit(clause.Associations).Save(&usersWithNoRefLinks).Error
						if e != nil {
							//saving model failed
							log.Printf("[REFLINKROUTINE]()()()@@@()()()()FAILED TO UPDATE USER LIST with referral links due to: %v\n", e)
						}
						time.Sleep(200 * time.Millisecond)

						return nil
					})
					if result.Error != nil {
						if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
							log.Println("[REFLINKROUTINE]()()()()()()()()()()()()()error occurred during batch processing:", result.Error.Error())

						}

					}
				})

				time.Sleep(15 * time.Minute)
			}
		}()

	}
	globalConfig.ChannelAccountKeysByAddress = make(map[string]*evmkeypair.Full)
	scas := strings.Split(os.Getenv("CHANNEL_ACCOUNTS"), ",")
	count := decimal.RequireFromString(os.Getenv("CHANNEL_ACCOUNT_MIN_COUNT")).IntPart()

	go func() {
		funder := evmkeypair.MustParseFull(os.Getenv("CHANNEL_ACCOUNT_FUNDER"))
		// instances booting together fund channel accounts one at a time (the
		// second then finds them funded); each funding transaction's nonce is
		// reserved under the funder key's own lock (network.ReserveNonce)
		lockErr := sharedconfig.WithKeyLock(globalConfig.DB, "channel-account-funding", 10*time.Minute, func() (lockedErr error) {
			var channelAccountsCSV string
			if len(scas) >= 1 {

				for _, v := range scas {
					k, e := evmkeypair.ParseFull(strings.ReplaceAll(v, " ", ""))
					if e != nil {
						log.Printf("[PARSE CHANNEL ACCOUNT]error parsing account %v:%v\n", v, e)
						continue
					}
					{
						//check if channel account is currently in use in pending shared access transaction
						var pendingTransaction userModels.PendingAuth
						errFetch := database.Where("transaction_status = 'PENDING' AND transaction_source = ?", k.Address()).First(&pendingTransaction).Error

						if errFetch == nil {
							//record was retrieved. save this in the map
							log.Printf("[ADDING KEY TO IN-USE CHANNEL ACCOUNT LIST] %v\n", k.Address())

							globalConfig.StoreInUseChannelAccount(k)
							//skip adding it to available channel accounts
							continue

						}

					}
					{
						//check if channel account is currently reserved by a fiat asset purchase invoice
						//(the async fiat purchase flow never creates a PendingAuth row, so it isn't caught
						//by the check above)
						var fiatInvoice userModels.FiatPaymentInvoice
						errFetchFiat := database.Where("transaction_source = ?", k.Address()).First(&fiatInvoice).Error

						if errFetchFiat == nil {
							if fiatInvoice.Status == "PENDING" {
								//still mid-flight, awaiting the payment provider's webhook - genuinely in use
								log.Printf("[ADDING KEY TO IN-USE CHANNEL ACCOUNT LIST] %v\n", k.Address())

								globalConfig.StoreInUseChannelAccount(k)
								//skip adding it to available channel accounts
								continue

							}
							if fiatInvoice.Status == "COMPLETED" {
								//the webhook's success path should already have released this account - its
								//presence here means that release didn't run for some reason (e.g. a crash
								//between the status flip and the release call). self-heal: don't mark it
								//in-use, fall through to the available pool below, and clear the stale
								//reference so this isn't re-detected on every future restart.
								log.Printf("[CHANNEL ACCOUNT RELEASE NOT PERSISTED - SELF-HEALING] %v (invoice %v was COMPLETED but still held a channel account)\n", k.Address(), fiatInvoice.ID)
								database.Model(&userModels.FiatPaymentInvoice{}).Where("id = ?", fiatInvoice.ID).Update("transaction_source", nil)
							}
							//any other status (e.g. EXPIRED) already had its channel account released by
							//the ExpireStalePaymentInvoices sweep - nothing to do, falls through to the pool
						}
					}
					// log.Printf("Channel Account to be used:%v\n", k.Address())
					//check minimum balance
					if len(channelAccountsCSV) == 0 {
						channelAccountsCSV = fmt.Sprintf("%s,", k.Seed())
					} else {
						if !strings.HasSuffix(channelAccountsCSV, ",") {
							channelAccountsCSV = fmt.Sprintf("%s,%s,", channelAccountsCSV, k.Seed())
						} else {
							channelAccountsCSV = fmt.Sprintf("%s%s,", channelAccountsCSV, k.Seed())
						}

					}
					if os.Getenv("CHECK_CHANNEL_ACCOUNT_BALANCE") == "0" || os.Getenv("CHECK_CHANNEL_ACCOUNT_BALANCE") == "" {
						if errSeed := sharedconfig.SeedChannelAccount(&globalConfig, k); errSeed != nil {
							log.Printf("[SEED CHANNEL ACCOUNT] error seeding %v: %v\n", k.Address(), errSeed)
						}
						continue
					}

					exists, _, nativeBal, _, _, _ := network.BlockchainAccountProperties(globalConfig.BantuExpansionClient, k.Address(), basetxn.NativeAsset{})
					_, _, _, _, sact, _ := network.BlockchainAccountProperties(globalConfig.BantuExpansionClient, funder.Address(), basetxn.NativeAsset{})

					var ops []basetxn.Operation
					if !exists {
						//fund from the funder

						ops = append(ops, &basetxn.CreateAccount{
							Destination: k.Address(),
							Amount:      os.Getenv("CHANNEL_ACCOUNT_FUNDING_AMOUNT"),
						})
					} else {
						if nativeBal.LessThan(decimal.RequireFromString(os.Getenv("CHANNEL_ACCOUNT_MIN_BALANCE"))) {
							ops = append(ops, &basetxn.Payment{
								Destination: k.Address(),
								Amount:      os.Getenv("CHANNEL_ACCOUNT_FUNDING_AMOUNT"),
								Asset:       basetxn.NativeAsset{},
							})
						}
					}
					if errSeed := sharedconfig.SeedChannelAccount(&globalConfig, k); errSeed != nil {
						log.Printf("[SEED CHANNEL ACCOUNT] error seeding %v: %v\n", k.Address(), errSeed)
					}
					if len(ops) == 0 {
						log.Println("NO OPERATIONS for this wallet", k.Address())
						continue
					}
					tx, err := basetxn.NewTransaction(
						basetxn.TransactionParams{
							SourceAccount:        sact.Address,
							IncrementSequenceNum: true,
							Operations:           ops,
							BaseFee:              2000,
							Memo:                 "Fund channel account",
						},
					)
					if err != nil {
						log.Println("[FUND CHANNEL ACCOUNT] error constructing transaction ", err)
						continue
					}

					tx, err = tx.Sign(globalConfig.BantuNetworkPassphrase, funder)
					if err != nil {
						log.Println("[FUND CHANNEL ACCOUNT] error signing transaction ", err)
						continue
					}

					fundTxRaw, err := tx.Base64()
					if err != nil {
						log.Println("[FUND CHANNEL ACCOUNT] error serializing transaction ", err)
						continue
					}
					hash, err := network.SubmitXdrWithSignature(globalConfig.BantuExpansionClient, funder.Address(), fundTxRaw, "")
					if err != nil {
						log.Println("[FUND CHANNEL ACCOUNT] error constructing transaction ", err)
						continue
					}
					log.Println("[FUND CHANNEL ACCOUNT] success ", hash)

				}
			}
			//

			//check of number of channel accounts is upto specified amount

			// new accounts are generated by one instance per hour, not by every
			// instance that boots
			if len(scas) < int(count) && claimChannelAccountGeneration(globalConfig.DB) {
				_, _, _, _, sact, _ := network.BlockchainAccountProperties(globalConfig.BantuExpansionClient, funder.Address(), basetxn.NativeAsset{})

				b := 0
				var ops []basetxn.Operation
				log.Printf("NUMBER OF SUPPLIED chan account %v is less than the required number %v\n", len(scas), count)
				for i := len(scas); i < int(count); i++ {
					b++
					// /
					k, errRandom := evmkeypair.Random()
					if errRandom != nil {
						log.Printf("[GENERATE CHANNEL ACCOUNT] error generating account: %v\n", errRandom)
						continue
					}

					log.Printf("Channel Account to be used:%v\n", k.Address())
					if errSeed := sharedconfig.SeedChannelAccount(&globalConfig, k); errSeed != nil {
						log.Printf("[SEED CHANNEL ACCOUNT] error seeding %v: %v\n", k.Address(), errSeed)
					}

					//check minimum balance

					if len(channelAccountsCSV) == 0 {
						channelAccountsCSV = fmt.Sprintf("%s,", k.Seed())
					} else {
						if !strings.HasSuffix(channelAccountsCSV, ",") {
							channelAccountsCSV = fmt.Sprintf("%s,%s,", channelAccountsCSV, k.Seed())
						} else {
							channelAccountsCSV = fmt.Sprintf("%s%s,", channelAccountsCSV, k.Seed())
						}

					}

					//fund from the funder

					ops = append(ops, &basetxn.CreateAccount{
						Destination: k.Address(),
						Amount:      os.Getenv("CHANNEL_ACCOUNT_FUNDING_AMOUNT"),
					})

					if b == 97 {
						tx, err := basetxn.NewTransaction(
							basetxn.TransactionParams{
								SourceAccount:        sact.Address,
								IncrementSequenceNum: true,
								Operations:           ops,
								BaseFee:              2000,
								Memo:                 "Fund channel account",
							},
						)
						if err != nil {
							log.Println("[FUND CHANNEL ACCOUNT] error constructing transaction ", err)
							continue
						}

						tx, err = tx.Sign(globalConfig.BantuNetworkPassphrase, funder)
						if err != nil {
							log.Println("[FUND CHANNEL ACCOUNT] error signing transaction ", err)
							continue
						}

						fundTxRaw, err := tx.Base64()
						if err != nil {
							log.Println("[FUND CHANNEL ACCOUNT] error serializing transaction ", err)
							continue
						}
						hash, err := network.SubmitXdrWithSignature(globalConfig.BantuExpansionClient, funder.Address(), fundTxRaw, "")
						if err != nil {
							log.Println("[FUND CHANNEL ACCOUNT] error constructing transaction ", err)
							continue
						}
						log.Println("[FUND CHANNEL ACCOUNT] success ", hash)

						//reset trx
						ops = make([]basetxn.Operation, 0)
						///
						_, _, _, _, sact, _ = network.BlockchainAccountProperties(globalConfig.BantuExpansionClient, funder.Address(), basetxn.NativeAsset{})
						b = 0
					}

				}

				if len(ops) > 0 {
					tx, err := basetxn.NewTransaction(
						basetxn.TransactionParams{
							SourceAccount:        sact.Address,
							IncrementSequenceNum: true,
							Operations:           ops,
							BaseFee:              2000,
							Memo:                 "Fund channel account",
						},
					)
					if err != nil {
						log.Println("[FUND CHANNEL ACCOUNT] error constructing transaction ", err)
						m.SendEmail(os.Getenv("CHANNEL_ACCOUNT_RECEIPIENT"), channelAccountsCSV)
						return
					}

					tx, err = tx.Sign(globalConfig.BantuNetworkPassphrase, funder)
					if err != nil {
						log.Println("[FUND CHANNEL ACCOUNT] error signing transaction ", err)
						m.SendEmail(os.Getenv("CHANNEL_ACCOUNT_RECEIPIENT"), channelAccountsCSV)
						return
					}

					fundTxRaw, err := tx.Base64()
					if err != nil {
						log.Println("[FUND CHANNEL ACCOUNT] error serializing transaction ", err)
						m.SendEmail(os.Getenv("CHANNEL_ACCOUNT_RECEIPIENT"), channelAccountsCSV)
						return
					}
					hash, err := network.SubmitXdrWithSignature(globalConfig.BantuExpansionClient, funder.Address(), fundTxRaw, "")
					if err != nil {
						log.Println("[FUND CHANNEL ACCOUNT] error constructing transaction ", err)
						m.SendEmail(os.Getenv("CHANNEL_ACCOUNT_RECEIPIENT"), channelAccountsCSV)
						return
					}
					log.Println("[FUND CHANNEL ACCOUNT] success ", hash)
				}

				//send this securely to remote service.
				m.SendEmail(os.Getenv("CHANNEL_ACCOUNT_RECEIPIENT"), channelAccountsCSV)
				log.Println("DONE FUNDING CHANNEL ACCOUNTS")
			}
			return
		})
		if lockErr != nil {
			log.Printf("[CHANNEL ACCOUNTS] funding skipped: %v\n", lockErr)
		}
	}()
	if os.Getenv("ENABLE_CRYPTO_WITHDRAWAL_SERVICE") == "1" {
		go func() {
			//LOAD WITHDRAWAL NETWORKS FROM 1L
			cl := strings.Split(os.Getenv("ONELIQUIDITY_WITHDRAWAL_CURRENCY_LIST"), ",")
			if len(cl) == 0 {
				//exit routine
				return
			}
			for {

				for _, c := range cl {
					//fetching currency withdrawal network list
					log.Println("<<<<<FETCHING/UPDATING WITHDRAWAL NETWORK PARAM FOR:", c)
					cacheKey := fmt.Sprintf("%s/%s?currency=%s", os.Getenv("ONELIQUIDITY_BASE_URL"), "wallets/v1/withdrawal/networks", c)
					redisCache.DeleteFromCache(cacheKey)
					userServices.GetWithdrawalNetworks(c, &globalConfig)
				}

				time.Sleep(800 * time.Second)
			}

		}()

	}

	//ACTIVATES PENDING PATRON SUBSCRIPTION
	if os.Getenv("ENABLE_PATRON") == "1" {
		go func() {
			err := UpdateUserPatronMemberships(database)
			if err != nil {
				log.Printf("[MAIN] error updating memberships: %v\n", err)
				//return
			}
			//clear the user cache
			time.Sleep(40 * time.Second)

		}()
	}

	{
		if os.Getenv("ENABLE_CRYPTO_DEPOSIT_MINTING") == "1" {
			go func() {
				for {
					// one instance at a time (the unique deposit ID also
					// refuses a second mint)
					found := false
					sharedconfig.WithSingletonLock(&globalConfig, "crypto-deposit-minting", 2*time.Minute, func() { found = mintNextCryptoDeposit(database, &globalConfig) })
					if !found {
						time.Sleep(60 * time.Second)
					}
				}
			}()
		}
	}

	{
		// Follow submitted wallet operations (Safe UserOperations) to
		// inclusion - recording outcomes and wallet activations - and
		// expire unsigned ones.
		go func() {
			for {
				sharedconfig.WithSingletonLock(&globalConfig, "track-wallet-operations", 30*time.Second, func() {
					ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
					defer cancel()
					userServices.TrackWalletOperations(ctx, &globalConfig)
					userServices.ExpireWalletOperations(&globalConfig)
				})
				time.Sleep(5 * time.Second)
			}
		}()
	}

	{
		// Finalize account recoveries whose recovery period is over, notice
		// cancellations and alert on recoveries started outside the backend.
		go func() {
			for {
				sharedconfig.WithSingletonLock(&globalConfig, "process-account-recoveries", 2*time.Minute, func() {
					ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
					defer cancel()
					userServices.ProcessAccountRecoveries(ctx, &globalConfig)
				})
				time.Sleep(30 * time.Second)
			}
		}()
	}

	{
		//Start processing payment streams
		go func() {
			for {

				sharedconfig.WithSingletonLock(&globalConfig, "monitor-stream", 20*time.Second, func() {
					MonitorStream(&globalConfig)
				})
				time.Sleep(5 * time.Minute)
			}
		}()
	}

	{
		//Start processing Sales
		go func() {

			for {

				sharedconfig.WithSingletonLock(&globalConfig, "activate-primary-sales", 20*time.Second, func() {
					userServices.ActivatePrimarySalesRoutine(&globalConfig)
				})
				time.Sleep(15 * time.Minute)
			}
		}()
		//Start processing Sales
		go func() {

			for {
				sharedconfig.WithSingletonLock(&globalConfig, "activate-secondary-sales", 20*time.Second, func() {
					userServices.ActivateSecondarySalesRoutine(&globalConfig)
				})
				time.Sleep(15 * time.Minute)
			}
		}()

		//Start processing Stablerails onboarding and onramping
		go func() {

			for {
				sharedconfig.WithSingletonLock(&globalConfig, "stablerail-onboarding-onramp", time.Minute, func() {
					userServices.ProcessUpdateStablerailOnboardingStatus(&globalConfig)
					time.Sleep(10 * time.Second)
					userServices.ProcessUpdateStablerailCNGNOnrampStatus(&globalConfig)
				})
				time.Sleep(10 * time.Second)

			}
		}()
	}

	{
		//Start processing Sale Notification for Interests
		go func() {
			for {

				sharedconfig.WithSingletonLock(&globalConfig, "send-pn-primary-sales-subscribers", 20*time.Second, func() {
					userServices.SendPNToSuscribersForPrimarySales(&globalConfig)
				})
				time.Sleep(30 * time.Second)
			}
		}()
	}
	{
		//Start fetching Stablertail supported bank codes
		go func() {

			for {
				sharedconfig.WithSingletonLock(&globalConfig, "stablerail-save-supported-banks", 30*time.Minute, func() {
					_, e := userServices.StablerailSaveSupportedBanks(&globalConfig)
					if e != nil {
						log.Printf("[Error Fetching stablerail banks] %v\n", e)
					}
				})
				time.Sleep(10 * time.Minute)
			}
		}()
	}

	{
		//Expire fiat payment invoices (e.g. fiat asset purchases) that have been stuck PENDING for
		//more than 2 days without a webhook confirmation
		go func() {
			for {
				sharedconfig.WithSingletonLock(&globalConfig, "expire-stale-payment-invoices", 90*time.Minute, func() {
					if e := userServices.ExpireStalePaymentInvoices(&globalConfig); e != nil {
						log.Printf("[MAIN] error expiring stale payment invoices: %v\n", e)
					}
				})
				time.Sleep(30 * time.Minute)
			}
		}()
	}

	//setup router

	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	// still passed to the handlers, but unused: partner callbacks go through
	// sharedconfig.SendCallback, which records and retries them
	callBackRetryChan := make(chan userModels.RetryCallbacks)
	var router *gin.Engine = gin.Default()
	// router.SetTrustedProxies(nil)
	router.Use(middleware.CORSMiddleware())

	root.Init(router)
	log.Println("##root services initialized##")
	users.Init(router, callBackRetryChan, &globalConfig)
	log.Println("##users services initialized##")
	payments.Init(router, callBackRetryChan, &globalConfig)
	log.Println("##payments services initialized##")
	swaps.Init(router, &globalConfig)
	log.Println("##swap services initialized##")

	serviceLinks.Init(router, &globalConfig)
	log.Println("##serviceLinks services initialized##")

	rates.Init(router, &globalConfig)
	log.Println("##rates services initialized##")
	msc.Init(router, &globalConfig)
	log.Println("##announcements/version services initialized##")
	callbacks.Init(router, callBackRetryChan, &globalConfig)

	p2p.Init(router, &globalConfig)
	log.Println("##callbacks services initialized##")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	log.Println("##swagger UI initialized##")
	// partner callbacks that failed are retried by one instance at a time
	go func() {
		for {
			sharedconfig.WithSingletonLock(&globalConfig, "deliver-callbacks", time.Minute, func() {
				ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
				defer cancel()
				sharedconfig.DeliverDueCallbacks(ctx, &globalConfig)
			})
			time.Sleep(10 * time.Second)
		}
	}()

	//run app
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	srv := &http.Server{Addr: ":" + port, Handler: router}
	go func() {
		log.Println("##service started##")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("[main] server: %v", err)
		}
	}()
	waitForShutdown(srv)
}

// waitForShutdown blocks until SIGINT/SIGTERM (e.g. the autoscaler removing
// this instance), then stops taking new background work, lets in-flight
// requests finish and waits for running locked work (background jobs,
// platform-key transactions) before returning. Locks held by an instance
// killed before that go stale and are taken over by the others.
func waitForShutdown(srv *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	grace := 60 * time.Second
	if d, err := time.ParseDuration(strings.TrimSpace(os.Getenv("SHUTDOWN_GRACE_PERIOD"))); err == nil && d > 0 {
		grace = d
	}
	log.Printf("[main] %v: shutting down (grace %v)", sig, grace)
	sharedconfig.BeginShutdown()
	ctx, cancel := context.WithTimeout(context.Background(), grace)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("[main] HTTP shutdown: %v", err)
	}
	deadline, _ := ctx.Deadline()
	if !sharedconfig.WaitForBackgroundWork(time.Until(deadline)) {
		log.Printf("[main] background work still running at the end of the grace period")
	}
	log.Println("[main] stopped")
}

func FindAllPendingSubscriptions(gc *gorm.DB) ([]userModels.UserPatronSubscriptionLog, error) {
	var pendingSubscriptions []userModels.UserPatronSubscriptionLog
	err := gc.Model(&userModels.UserPatronSubscriptionLog{}).
		Where("CAST(effective_date AS DATE) >= CAST(? AS DATE)", time.Now()).
		Find(&pendingSubscriptions).Error
	if err != nil {
		log.Printf("[FindPendingSubscriptions] error: %v\n", err)
		err = &tErrors.ErrorTemporaryServerError{}
		return nil, err
	}

	return pendingSubscriptions, nil
}

func UpdateUserPatronMemberships(db *gorm.DB) error {
	// Step 1: Fetch pending subscriptions
	pendingSubscriptions, err := FindAllPendingSubscriptions(db)
	if err != nil {
		log.Printf("[UpdateUserPatronMemberships] NO PENDING memberships: %v\n", err)
		return err
	}
	if len(pendingSubscriptions) == 0 {
		log.Printf("[UpdateUserPatronMemberships] no pending memberships: %v\n", err)
		return nil
	}

	// Step 2: Extract unique usernames
	uniqueUsernames := make(map[string]struct{})
	for _, subscription := range pendingSubscriptions {
		uniqueUsernames[subscription.Username] = struct{}{}
	}

	// Step 3: Fetch UserPatronMemberships based on usernames
	var memberships []userModels.UserPatronMembership
	err = db.Where("username IN (?)", getUniqueUsernamesSlice(uniqueUsernames)).Find(&memberships).Error
	if err != nil {
		log.Printf("[UpdateUserPatronMemberships] error fetching memberships: %v\n", err)
		err = &tErrors.ErrorTemporaryServerError{}
		return err
	}

	// Step 4: Update UserPatronMemberships with data from subscriptions
	for _, subscription := range pendingSubscriptions {
		for i, membership := range memberships {
			if membership.Username == subscription.Username {
				// Update membership with data from subscription
				memberships[i].PatronPackageID = subscription.PatronPackageID
				memberships[i].PatronTierID = subscription.PatronTierID
				memberships[i].ValidTill = subscription.ValidTill
			}
		}
	}

	// Save the updated memberships back to the database
	for _, membership := range memberships {
		err := db.Omit(clause.Associations).Save(&membership).Error
		if err != nil {
			log.Printf("[UpdateUserPatronMemberships] error updating membership: %v\n", err)
			err = &tErrors.ErrorTemporaryServerError{}
			return err
		}
	}

	return nil
}

func getUniqueUsernamesSlice(uniqueUsernames map[string]struct{}) []string {
	usernames := make([]string, 0, len(uniqueUsernames))
	for username := range uniqueUsernames {
		usernames = append(usernames, username)
	}
	return usernames
}

// MonitorStream watched Horizon's global operation stream
// (client.StreamOperations, every operation network-wide joined with its
// transaction) and, for each one, invalidated the cached balance/wallet/
// user entries of every account it touched (see the former
// ProcessOperation) so cached reads reflected on-chain state promptly.
// Base has no equivalent account-agnostic operation feed - the
// Base-native way to get this is subscribing to eth_subscribe("logs")
// for the native transfers/B20 token contracts this app cares about and
// decoding events (go-ethereum's ethclient.SubscribeFilterLogs), which
// needs a WS-capable RPC endpoint and a real event-decoding design,
// tracked as a follow-up out of scope for this alteration pass. Until
// then, each payment/swap call path invalidates its own affected
// accounts' cache entries directly (see internal/components/*/services),
// so this is a documented no-op rather than a broken poll loop.
func MonitorStream(gc *sharedconfig.GlobalConfig) {
	log.Println("[MonitorStream] Base has no global operation stream to watch (see doc comment) - idling.")
}

// mintNextCryptoDeposit mints the oldest crypto deposit not minted yet; it
// reports whether there was one.
func mintNextCryptoDeposit(database *gorm.DB, gc *sharedconfig.GlobalConfig) bool {
	var di userModels.CallbackDepositItem
	e := database.Order("created_at ASC").Where("minted = 0").First(&di).Error
	if e != nil {
		log.Println("[MINTING INITIATOR] Unable to locate waiting callback deposits.")
		return false
	}

	//preparing minting
	log.Printf("[MINTING INITIATOR] Preparing to mint %v for address %v\n", di.Currency, di.ToAddress)
	da, err := userModels.CryptoDepositAddress(di.ToAddress).GetDetail(di.Currency, gc)
	if err != nil {
		log.Printf("[MINTING INITIATOR] error getting address owner to mint %v %v, error: %v\n", di.Currency, di.ToAddress, err)

		return true
	}
	//initiate minting
	amountLessFees := ((decimal.RequireFromString((di.Amount).(string)).Sub(decimal.RequireFromString(di.Fees))).Div(decimal.NewFromInt(10).Pow(decimal.NewFromInt(int64(di.Decimal))))).Truncate(7)
	log.Printf("preparing to mint %v %v to %v\n", amountLessFees.String(), di.Currency, da.TrovoWalletAddress)

	signerPK := os.Getenv("CRYPTO_DEPOSIT_MINTING_INITIATOR_PUBLIC_KEY")
	signerUser, err := userModels.UserWalletID(signerPK).GetWalletOwner(gc.DB, gc)
	if err != nil {

		log.Printf("[MINTING INITIATOR] error getting initiator user to mint %v %v, error: %v\n", di.Currency, di.ToAddress, err)

		return true
	}
	ca, err := userModels.Currency(da.Currency).GetCurratedAsset(gc)
	if err != nil {

		log.Printf("[MINTING INITIATOR] error getting curated asset to mint %v %v, error: %v\n", di.Currency, di.ToAddress, err)

		return true
	}
	sourceWallet, err := userModels.UserWalletID(ca.ContractAddress).GetWallet(gc.DB, gc)
	if err != nil {

		log.Printf("[MINTING INITIATOR] error getting initiator user to mint %v %v, error: %v\n", di.Currency, di.ToAddress, err)

		return true
	}
	// fetch deposit
	pdi, err := userServices.GetADepositByID(di.DepositID, gc)
	if err != nil {

		log.Printf("[MINTING INITIATOR] error getting deposit from service %v %v, error: %v\n", di.Currency, di.DepositID, err)

		return true
	}
	//build deposit
	layout := "2006-01-02T15:04:05.000Z"
	createdAt, _ := time.Parse(layout, pdi.CreatedAt)
	updatedAt, _ := time.Parse(layout, pdi.UpdatedAt)
	depositItem := userModels.CryptoDeposit{
		CreatedAt:          createdAt,
		UpdatedAt:          updatedAt,
		TrovoWalletAddress: da.TrovoWalletAddress,
		DepositID:          pdi.DepositID,
		TxID:               pdi.TxID,
		Amount:             amountLessFees.String(),
		Currency:           ca.AssetCode,
		Decimal:            pdi.Decimal,
		Fees:               pdi.Fees,
		FromAddress:        pdi.FromAddress,
		ToAddress:          pdi.ToAddress,
		IsCompleted:        pdi.IsCompleted,
		IsValid:            pdi.IsValid,
		IsVerified:         pdi.IsVerified,
	}

	mintingInfo := userModels.MintingInfo{
		Destination:     da.TrovoWalletAddress,
		Memo:            fmt.Sprintf("%v %v", amountLessFees.String(), da.Currency),
		ContractAddress: ca.ContractAddress,
		AssetCode:       ca.AssetCode,
		Amount:          amountLessFees.String(),
		Commit:          1,
	}

	dbtx := database.Begin()
	e = dbtx.Omit(clause.Associations).Create(&depositItem).Error
	if e != nil {
		dbtx.Rollback()
		log.Printf("[MINTING INITIATOR] error creating deposit item. error: %v\nDepositItem: %+v\n", e, depositItem)

		return true
	}
	di.Minted = 1
	e = dbtx.Omit(clause.Associations).Save(&di).Error
	if e != nil {
		dbtx.Rollback()
		log.Printf("[MINTING INITIATOR] error saving callback item. error: %v\nCallbackDepositItem: %+v\n", e, di)

		return true
	}

	_, _, err = userServices.MintAsset(&signerUser, &sourceWallet, &mintingInfo, gc)

	if err != nil {
		dbtx.Rollback()
		log.Printf("[MINTING INITIATOR] error MINTING deposit item. error: %v\nDepositItem: %+v\n", err, depositItem)
		return true
	}
	dbtx.Commit()
	log.Printf("[MINTING INITIATOR]  Minted %v %v to %v\n", mintingInfo.Amount, mintingInfo.AssetCode, mintingInfo.Destination)
	{
		//start push notificationMessage

		permissionList := sourceWallet.Permissions
		for _, v := range permissionList {
			if v.Permission != "APPROVER" {
				continue
			}
			u, e := userModels.Username(v.TargetUsername).GetSimpleUser(gc.DB, gc)
			if e != nil {
				continue
			}

			dataPayload := make(map[string]string)
			dataPayload["route"] = "pendingApproval"

			u.SendPushMessage(fmt.Sprintf("%v %v minting request submitted on %v!", mintingInfo.Amount, da.Currency, sourceWallet.Alias), fmt.Sprintf("Request:\n %v", mintingInfo.ReturnedDescription), "", dataPayload, gc)

		}
	}
	return true
}

// claimChannelAccountGeneration reports whether this instance may generate
// new channel accounts: at most one instance per hour does (the claim is
// left to expire, not released).
func claimChannelAccountGeneration(db *gorm.DB) bool {
	const name = "channel-account-generation"
	if err := sharedconfig.EnsureDistributedLock(db, name); err != nil {
		return false
	}
	ok, err := sharedconfig.TryAcquireLock(db, name, sharedconfig.InstanceIdentity(), time.Hour)
	return err == nil && ok
}
