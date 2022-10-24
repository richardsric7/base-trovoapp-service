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
	"trovo-wallet-api/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
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
			"DYNAMIC_LINKS_IOS_BUNDLE_ID", "DYNAMIC_LINKS_FALLBACK_BASE_URL", "FBDL_SERVICE_URLS", "MAILGUN_DOMAIN", "XBN_ASSET_IMAGE_URL",
			"GC", "GOOGLE_PROJECT_ID", "ACCOUNT_RECOVERY_SALT", "MNEMONIC_ACCOUNT_RECOVERY", "RECOVERY_SIGNER_ACTIVATION_AMOUNT",
			"NATIVE_ASSET_CODE", "ACCOUNT_RECOVERY_MINIMUM_BALANCE", "SHARED_ACCESS_FEE_ADDRESS",
			"SHARED_ACCESS_FEE_AMOUNT",
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
				e := database.Where("referral_link is null AND referral_qr_code is null AND suspended = ?", 0).First(&userModels.User{}).Error
				if e != nil {
					time.Sleep(15 * time.Minute)
					continue
				}

				result := database.Where("referral_link is null AND referral_qr_code is null AND suspended = ?", 0).FindInBatches(&usersWithNoRefLinks, batchSize, func(tx *gorm.DB, batch int) error {
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
	log.Println(router.Run())

}
