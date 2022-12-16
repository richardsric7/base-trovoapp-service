package main

import (
	cache "trovo-wallet-api/internal/cache"
	"trovo-wallet-api/internal/network"
	pns "trovo-wallet-api/internal/pns"
	"trovo-wallet-api/internal/sharedconfig"

	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	payments "trovo-wallet-api/internal/components/payments/controllers"
	paymentModels "trovo-wallet-api/internal/components/payments/models"
	rates "trovo-wallet-api/internal/components/rates/controllers"
	root "trovo-wallet-api/internal/components/root/controllers"
	serviceLinks "trovo-wallet-api/internal/components/servicelinks/controllers"
	swaps "trovo-wallet-api/internal/components/swaps/controllers"
	users "trovo-wallet-api/internal/components/users/controllers"
	userModels "trovo-wallet-api/internal/components/users/models"
	db "trovo-wallet-api/internal/db"
	dl "trovo-wallet-api/internal/dynamiclinks"
	m "trovo-wallet-api/internal/mail"
	"trovo-wallet-api/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"github.com/shopspring/decimal"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/txnbuild"
	"gorm.io/gorm"
)

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
			"MNEMONIC_TEMP_ACCOUNTS", "BLOCKCHAIN_BASE_RESERVE", "MAILGUN_PRIVATE_API_KEY", "CDB_CONNECTION_STRING",
			"IPAPI_KEY", "IPAPI_HOST", "VERIFICATION_CODE_SALT", "ENABLE_EMAIL_VALIDATION", "ENABLE_CACHING", "DEFAULT_ASSET_IMAGE_URL",
			"REDIS_HOST", "REDIS_PORT", "DYNAMIC_LINKS_API_KEY", "DYNAMIC_LINKS_DOMAIN_PREFIX", "DYNAMIC_LINKS_ANDROID_PACKAGE_NAME",
			"DYNAMIC_LINKS_IOS_BUNDLE_ID", "DYNAMIC_LINKS_FALLBACK_BASE_URL", "FBDL_SERVICE_URLS", "MAILGUN_DOMAIN", "NATIVE_ASSET_IMAGE_URL",
			"GC", "GOOGLE_PROJECT_ID", "ACCOUNT_RECOVERY_SALT", "MNEMONIC_ACCOUNT_RECOVERY", "RECOVERY_SIGNER_ACTIVATION_AMOUNT",
			"NATIVE_ASSET_CODE", "ACCOUNT_RECOVERY_MINIMUM_BALANCE", "SHARED_ACCESS_FEE_ADDRESS",
			"SHARED_ACCESS_FEE_AMOUNT", "CHANNEL_ACCOUNTS", "WALLET_SIGNER_ACTIVATION_AMOUNT", "WALLET_DOMAIN",
			"MNEMONIC_BULK_PAYMENT", "BULK_PAYMENT_SALT", "ENCODER_SALT", "MARKET_MAKING_SALT",
			"MNEMONIC_MARKET_MAKING", "MAX_ISSUED_ASSETS_PER_WALLET", "CHECK_CHANNEL_ACCOUNT_BALANCE",
			"JWT_ACCESS_SECRET", "JWT_TOKEN_EXPIRY", "JWT_REFRESH_TOKEN_EXPIRY",
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
	log.Println("starting migration")
	//migrate DB models if any
	db.MigrateDB(database)

	errMigrate := roachDB.AutoMigrate(&paymentModels.TrackedWallet{})
	if errMigrate != nil {
		if !strings.Contains(errMigrate.Error(), "constraint") {
			log.Fatalf("Error migrating TrackedWallet model, error: %v", errMigrate)
		}
	}
	log.Println("migrating tracked wallet done...")
	errMigrate = roachDB.AutoMigrate(&paymentModels.TrackedPublicKey{})
	if errMigrate != nil {
		if !strings.Contains(errMigrate.Error(), "constraint") {
			log.Fatalf("Error migrating TrackedPublicKey model, error: %v", errMigrate)
		}
	}
	log.Println("migrating tracked public key done...")

	//setup redis
	log.Println("migration done...")
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

	cas := strings.Split(os.Getenv("FBDL_SERVICE_URLS"), ",")
	dynamicLinkServiceUrlChan := make(chan string, len(cas))

	if len(cas) >= 1 {
		for _, u := range cas {

			log.Printf("Firebase Dynamic Links service URL to be used: %v", u)
			dynamicLinkServiceUrlChan <- u
		}
	}

	{

		//update referral links for people with no referral link
		go func() {
			log.Println("##[REFLINKROUTINE] started routine to update referral links for people with no referral link")

			batchSize := 1
			var usersWithNoRefLinks []userModels.User
			dynamicLinkServiceUrl := <-dynamicLinkServiceUrlChan
			defer func() {
				dynamicLinkServiceUrlChan <- dynamicLinkServiceUrl
			}()
			for {
				e := database.Where("referral_qr_code is null AND suspended = ?", 0).First(&userModels.User{}).Error
				if e != nil {
					time.Sleep(15 * time.Minute)
					continue
				}

				result := database.Where("referral_qr_code is null AND suspended = ?", 0).FindInBatches(&usersWithNoRefLinks, batchSize, func(tx *gorm.DB, batch int) error {
					for i, u := range usersWithNoRefLinks {
						rld, errLink := dl.GenerateReferralLinkWithStaticURL(u.Username, dynamicLinkServiceUrl, &redisCache)
						if errLink != nil {
							continue
						}
						//link retrieved
						usersWithNoRefLinks[i].ReferralLink = &rld.DynamicLink
						usersWithNoRefLinks[i].ReferralQrCode = &rld.QRCode
					}

					e := database.Save(&usersWithNoRefLinks).Error
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

				time.Sleep(15 * time.Minute)
			}
		}()

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
		DynamicLinkServiceURLChan: dynamicLinkServiceUrlChan,
		PNSContext:                pnsContext,
		RedisCache:                &redisCache,
		DB:                        database,
		PushNotificationClient:    pnsClient,
		RoachDB:                   roachDB,
		BantuExpansionClient:      network.GetBlockchainClient(),
		BantuNetworkPassphrase:    network.GetBlockchainNetworkPassPhrase(),
		FirebaseStorageUploader: &sharedconfig.ClientUploader{
			Client:     storageClient,
			ProjectID:  os.Getenv("GOOGLE_PROJECT_ID"),
			BucketName: os.Getenv("STORAGE_BUCKET_NAME"),
			UploadPath: os.Getenv("STORAGE_BUCKET_NAME"),
		},
	}
	// clear cache
	cacheKey := "[GET] /v1/rates"
	globalConfig.RedisCache.InvalidateCachedHttpResponse(cacheKey)
	globalConfig.RedisCache.DeleteFromCache(cacheKey)
	scas := strings.Split(os.Getenv("CHANNEL_ACCOUNTS"), ",")
	count := decimal.RequireFromString(os.Getenv("CHANNEL_ACCOUNT_MIN_COUNT")).IntPart()
	if len(scas) > int(count) {
		globalConfig.ChannelAccounts = make(chan *keypair.Full, len(scas))
	} else {
		globalConfig.ChannelAccounts = make(chan *keypair.Full, count)
	}

	go func() {
		var channelAccountsCSV string
		funder := keypair.MustParseFull(os.Getenv("CHANNEL_ACCOUNT_FUNDER"))
		if len(scas) >= 1 {

			for _, v := range scas {
				k, e := keypair.ParseFull(strings.ReplaceAll(v, " ", ""))
				if e != nil {
					log.Printf("[PARSE CHANNEL ACCOUNT]error parsing account %v:%v\n", v, e)
					continue
				}
				log.Printf("Channel Account to be used:%v\n", k.Address())
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
					globalConfig.ChannelAccounts <- k
					continue
				}

				exists, _, nativeBal, _, _, _ := network.BlockchainAccountProperties(globalConfig.BantuExpansionClient, k.Address(), txnbuild.NativeAsset{})
				_, _, _, _, sact, _ := network.BlockchainAccountProperties(globalConfig.BantuExpansionClient, funder.Address(), txnbuild.NativeAsset{})

				var ops []txnbuild.Operation
				if !exists {
					//fund from the funder

					ops = append(ops, &txnbuild.CreateAccount{
						Destination: k.Address(),
						Amount:      os.Getenv("CHANNEL_ACCOUNT_FUNDING_AMOUNT"),
					})
				} else {
					if nativeBal.LessThan(decimal.RequireFromString(os.Getenv("CHANNEL_ACCOUNT_MIN_BALANCE"))) {
						ops = append(ops, &txnbuild.Payment{
							Destination: k.Address(),
							Amount:      os.Getenv("CHANNEL_ACCOUNT_FUNDING_AMOUNT"),
							Asset:       txnbuild.NativeAsset{},
						})
					}
				}
				globalConfig.ChannelAccounts <- k
				if len(ops) == 0 {
					log.Println("NO OPeRATIONS for this wallet", k.Address())
					continue
				}
				tx, err := txnbuild.NewTransaction(
					txnbuild.TransactionParams{
						SourceAccount:        sact,
						IncrementSequenceNum: true,
						Operations:           ops,
						BaseFee:              txnbuild.MinBaseFee,
						Preconditions: txnbuild.Preconditions{
							TimeBounds: txnbuild.NewInfiniteTimeout(),
						},
						Memo: txnbuild.MemoText("Fund channel account"),
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

				hTranx, err := globalConfig.BantuExpansionClient.SubmitTransaction(tx)
				if err != nil {
					log.Println("[FUND CHANNEL ACCOUNT] error constructing transaction ", err)
					continue
				}
				log.Println("[FUND CHANNEL ACCOUNT] success ", hTranx.Hash)

			}
		}
		//

		//check of number of channel accounts is upto specified amount

		if len(scas) < int(count) {
			_, _, _, _, sact, _ := network.BlockchainAccountProperties(globalConfig.BantuExpansionClient, funder.Address(), txnbuild.NativeAsset{})

			b := 0
			var ops []txnbuild.Operation
			log.Printf("NUMBER OF SUPPLIED chan account %v is less than the required number %v\n", len(scas), count)
			for i := len(scas); i < int(count); i++ {
				b++
				// /
				k := keypair.MustRandom()

				log.Printf("Channel Account to be used:%v\n", k.Address())
				globalConfig.ChannelAccounts <- k

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

				ops = append(ops, &txnbuild.CreateAccount{
					Destination: k.Address(),
					Amount:      os.Getenv("CHANNEL_ACCOUNT_FUNDING_AMOUNT"),
				})

				if b == 97 {
					tx, err := txnbuild.NewTransaction(
						txnbuild.TransactionParams{
							SourceAccount:        sact,
							IncrementSequenceNum: true,
							Operations:           ops,
							BaseFee:              txnbuild.MinBaseFee,
							Preconditions: txnbuild.Preconditions{
								TimeBounds: txnbuild.NewInfiniteTimeout(),
							},
							Memo: txnbuild.MemoText("Fund channel account"),
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

					hTranx, err := globalConfig.BantuExpansionClient.SubmitTransaction(tx)
					if err != nil {
						log.Println("[FUND CHANNEL ACCOUNT] error constructing transaction ", err)
						continue
					}
					log.Println("[FUND CHANNEL ACCOUNT] success ", hTranx.Hash)

					//reset trx
					ops = make([]txnbuild.Operation, 0)
					///
					_, _, _, _, sact, _ = network.BlockchainAccountProperties(globalConfig.BantuExpansionClient, funder.Address(), txnbuild.NativeAsset{})
					b = 0
				}

			}

			tx, err := txnbuild.NewTransaction(
				txnbuild.TransactionParams{
					SourceAccount:        sact,
					IncrementSequenceNum: true,
					Operations:           ops,
					BaseFee:              txnbuild.MinBaseFee,
					Preconditions: txnbuild.Preconditions{
						TimeBounds: txnbuild.NewInfiniteTimeout(),
					},
					Memo: txnbuild.MemoText("Fund channel account"),
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

			hTranx, err := globalConfig.BantuExpansionClient.SubmitTransaction(tx)
			if err != nil {
				log.Println("[FUND CHANNEL ACCOUNT] error constructing transaction ", err)
				m.SendEmail(os.Getenv("CHANNEL_ACCOUNT_RECEIPIENT"), channelAccountsCSV)
				return
			}
			log.Println("[FUND CHANNEL ACCOUNT] success ", hTranx.Hash)

			//send this securely to remote service.
			m.SendEmail(os.Getenv("CHANNEL_ACCOUNT_RECEIPIENT"), channelAccountsCSV)
			log.Println("DONE FUNDING CHANNEL ACCOUNTS")
		}

	}()
	//setup router

	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	callBackRetryChan := make(chan userModels.RetryCallbacks, 200000)
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

	//run app
	log.Println("##service started##")
	if len(os.Getenv("PORT")) > 0 {
		log.Println(router.Run(":" + os.Getenv("PORT")))

	} else {
		log.Println(router.Run(":8080"))
	}

}
