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

	msc "trovo-wallet-api/internal/components/announcements/controllers"
	callbacks "trovo-wallet-api/internal/components/callbacks/controllers"
	payments "trovo-wallet-api/internal/components/payments/controllers"
	paymentModels "trovo-wallet-api/internal/components/payments/models"
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
			"SHARED_ACCESS_PAYMENT_FEE_AMOUNT", "CHANNEL_ACCOUNTS", "WALLET_SIGNER_ACTIVATION_AMOUNT", "WALLET_DOMAIN",
			"MNEMONIC_BULK_PAYMENT", "BULK_PAYMENT_SALT", "ENCODER_SALT", "MARKET_MAKING_SALT",
			"MNEMONIC_MARKET_MAKING", "MAX_ISSUED_ASSETS_PER_WALLET", "CHECK_CHANNEL_ACCOUNT_BALANCE",
			"JWT_ACCESS_SECRET", "JWT_TOKEN_EXPIRY", "JWT_REFRESH_TOKEN_EXPIRY",
			"SUBWALLET_FEE_AMOUNT_USD", "SUBWALLET_FEE_ASSET_ISSUER", "SUBWALLET_FEE_ASSET_CODE",
			"SUBWALLET_FEE_WALLET", "DOLLAR_ASSET", "MARKET_MAKING_FEE_ENABLED", "SWAP_FEE_ENABLED",
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
		if os.Getenv("ENABLE_CRYPTO_DEPOSIT_MINTING") == "1" && len(os.Getenv("CRYPTO_DEPOSIT_MINTING_INITIATOR_PUBLIC_KEY")) != 56 {
			log.Println("CRYPTO_DEPOSIT_MINTING_INITIATOR_PUBLIC_KEY environment variable is required when ENABLE_CRYPTO_DEPOSIT_MINTING is set to 1")

			exit = true
		}
		if os.Getenv("MARKET_MAKING_FEE_ENABLED") == "1" && len(os.Getenv("MARKET_MAKING_FEE_WALLET")) != 56 {
			log.Println("MARKET_MAKING_FEE_WALLET environment variable is required when MARKET_MAKING_FEE_ENABLED is set to 1")

			exit = true
		}

		if os.Getenv("ENABLE_NAIRA_ASSET_BY_DEFAULT") == "1" && len(os.Getenv("NAIRA_ASSET")) < 60 {
			log.Println("NAIRA_ASSET environment variable is required when ENABLE_NAIRA_ASSET_BY_DEFAULT is set to 1")

			exit = true
		}
		if os.Getenv("MARKET_MAKING_FEE_ENABLED") == "1" {
			_, err := keypair.ParseFull(os.Getenv("MARKET_MAKING_FEE_WALLET"))
			if err != nil {
				log.Println("MARKET_MAKING_FEE_WALLET  is invalid wallet secret key")

				exit = true
			}
		}
		if os.Getenv("SWAP_FEE_ENABLED") == "1" && len(os.Getenv("SWAP_FEE_WALLET")) != 56 {
			log.Println("SWAP_FEE_WALLET environment variable is required when SWAP_FEE_ENABLED is set to 1")

			exit = true
		}
		if os.Getenv("SWAP_FEE_ENABLED") == "1" {
			_, err := keypair.ParseFull(os.Getenv("SWAP_FEE_WALLET"))
			if err != nil {
				log.Println("SWAP_FEE_WALLET is invalid wallet secret key")

				exit = true
			}
		}

		if os.Getenv("SHARED_ACCESS_FEE_ENABLED") == "1" && len(os.Getenv("SHARED_ACCESS_FEE_WALLET")) != 56 {
			log.Println("SHARED_ACCESS_FEE_WALLET environment variable is required when SHARED_ACCESS_FEE_ENABLED is set to 1")

			exit = true
		}
		if os.Getenv("SHARED_ACCESS_FEE_ENABLED") == "1" {
			_, err := keypair.ParseFull(os.Getenv("SHARED_ACCESS_FEE_WALLET"))
			if err != nil {
				log.Println("SHARED_ACCESS_FEE_WALLET is invalid wallet secret key")

				exit = true
			}
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
	{
		// clear cache
		cacheKeyInfo := "curatedAssets_"
		redisCache.DeleteFromCache(cacheKeyInfo)
		// clear cache
		cacheKey := "[GET] /v1/rates"
		redisCache.InvalidateCachedHttpResponse(cacheKey)
		redisCache.DeleteFromCache(cacheKey)
	}

	cas := strings.Split(os.Getenv("FBDL_SERVICE_URLS"), ",")
	dynamicLinkServiceUrlChan := make(chan string, len(cas))

	if len(cas) >= 1 {
		for _, u := range cas {

			log.Printf("Firebase Dynamic Links service URL to be used: %v", u)
			dynamicLinkServiceUrlChan <- u
		}
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
						rld, errLink := dl.GenerateReferralLinkWithStaticURL(u.Username, dynamicLinkServiceUrl, &globalConfig)
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

	globalConfig.InUseChannelAccounts = make(map[string]*keypair.Full)
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

	{
		if os.Getenv("ENABLE_CRYPTO_DEPOSIT_MINTING") == "1" {
			go func() {
				for {
					var di userModels.CallbackDepositItem

					e := database.Order("created_at ASC").Where("minted = 0").First(&di).Error
					if e != nil {
						log.Println("[MINTING INITIATOR] Unable to locate waiting callback deposits.")
						time.Sleep(60 * time.Second)
						continue
					}
					//preparing minting
					log.Printf("[MINTING INITIATOR] Preparing to mint %v for address %v\n", di.Currency, di.ToAddress)
					da, err := userModels.CryptoDepositAddress(di.ToAddress).GetDetail(di.Currency, &globalConfig)
					if err != nil {
						log.Printf("[MINTING INITIATOR] error getting address owner to mint %v %v, error: %v\n", di.Currency, di.ToAddress, err)

						continue
					}
					//initiate minting
					amountLessFees := ((decimal.RequireFromString((di.Amount).(string)).Sub(decimal.RequireFromString(di.Fees))).Div(decimal.NewFromInt(10).Pow(decimal.NewFromInt(int64(di.Decimal))))).Truncate(7)
					log.Printf("preparing to mint %v %v to %v\n", amountLessFees.String(), di.Currency, da.TrovoWalletPublicKey)

					signerPK := os.Getenv("CRYPTO_DEPOSIT_MINTING_INITIATOR_PUBLIC_KEY")
					signerUser, err := userModels.UserWalletID(signerPK).GetWalletOwner(globalConfig.DB, &globalConfig)
					if err != nil {

						log.Printf("[MINTING INITIATOR] error getting initiator user to mint %v %v, error: %v\n", di.Currency, di.ToAddress, err)

						continue
					}
					ca, err := userModels.Currency(da.Currency).GetCurratedAsset(&globalConfig)
					if err != nil {

						log.Printf("[MINTING INITIATOR] error getting curated asset to mint %v %v, error: %v\n", di.Currency, di.ToAddress, err)

						continue
					}
					sourceWallet, err := userModels.UserWalletID(ca.AssetIssuer).GetWallet(globalConfig.DB, &globalConfig)
					if err != nil {

						log.Printf("[MINTING INITIATOR] error getting initiator user to mint %v %v, error: %v\n", di.Currency, di.ToAddress, err)

						continue
					}
					// fetch deposit
					pdi, err := userServices.GetADepositByID(di.DepositID, &globalConfig)
					if err != nil {

						log.Printf("[MINTING INITIATOR] error getting deposit from service %v %v, error: %v\n", di.Currency, di.DepositID, err)

						continue
					}
					//build deposit
					layout := "2006-01-02T15:04:05.000Z"
					createdAt, _ := time.Parse(layout, pdi.CreatedAt)
					updatedAt, _ := time.Parse(layout, pdi.UpdatedAt)
					depositItem := userModels.CryptoDeposit{
						CreatedAt:            createdAt,
						UpdatedAt:            updatedAt,
						TrovoWalletPublicKey: da.TrovoWalletPublicKey,
						DepositID:            pdi.DepositID,
						TxID:                 pdi.TxID,
						Amount:               amountLessFees.String(),
						Currency:             ca.AssetCode,
						Decimal:              pdi.Decimal,
						Fees:                 pdi.Fees,
						FromAddress:          pdi.FromAddress,
						ToAddress:            pdi.ToAddress,
						IsCompleted:          pdi.IsCompleted,
						IsValid:              pdi.IsValid,
						IsVerified:           pdi.IsVerified,
					}

					mintingInfo := userModels.MintingInfo{
						Destination: da.TrovoWalletPublicKey,
						Memo:        fmt.Sprintf("%v %v", amountLessFees.String(), da.Currency),
						AssetIssuer: ca.AssetIssuer,
						AssetCode:   ca.AssetCode,
						Amount:      amountLessFees.String(),
						Commit:      1,
					}

					dbtx := database.Begin()
					e = dbtx.Create(&depositItem).Error
					if e != nil {
						dbtx.Rollback()
						log.Printf("[MINTING INITIATOR] error creating deposit item. error: %v\nDepositItem: %+v\n", e, depositItem)

						continue
					}
					di.Minted = 1
					e = dbtx.Save(&di).Error
					if e != nil {
						dbtx.Rollback()
						log.Printf("[MINTING INITIATOR] error saving callback item. error: %v\nCallbackDepositItem: %+v\n", e, di)

						continue
					}

					_, _, err = userServices.MintAsset(&signerUser, &sourceWallet, &mintingInfo, &globalConfig)

					if err != nil {
						dbtx.Rollback()
						log.Printf("[MINTING INITIATOR] error MINTING deposit item. error: %v\nDepositItem: %+v\n", err, depositItem)
						continue
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
							u, e := userModels.Username(v.TargetUsername).GetSimpleUser(globalConfig.DB, &globalConfig)
							if e != nil {
								continue
							}

							dataPayload := make(map[string]string)
							dataPayload["route"] = "pendingApproval"

							u.SendPushMessage(fmt.Sprintf("%v %v minting request submitted on %v!", mintingInfo.Amount, da.Currency, sourceWallet.Alias), fmt.Sprintf("Request:\n %v", mintingInfo.ReturnedDescription), "", dataPayload, &globalConfig)

						}
					}

				}
			}()

		}
	}

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
	msc.Init(router, &globalConfig)
	log.Println("##announcements/version services initialized##")
	callbacks.Init(router, callBackRetryChan, &globalConfig)
	log.Println("##callbacks services initialized##")
	//run app
	log.Println("##service started##")
	if len(os.Getenv("PORT")) > 0 {
		log.Println(router.Run(":" + os.Getenv("PORT")))

	} else {
		log.Println(router.Run(":8080"))
	}

}
