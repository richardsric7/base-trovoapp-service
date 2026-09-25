package main

import (
	_ "admin-panel-dashboard/docs" // very important
	"admin-panel-dashboard/internal/cache"
	auth "admin-panel-dashboard/internal/components/auth/controllers"
	"errors"
	"strings"

	// offers "admin-panel-dashboard/internal/components/p2p_offers/controllers"
	//orders "admin-panel-dashboard/internal/components/p2p_orders/controllers"
	//payments "admin-panel-dashboard/internal/components/payments/controllers"
	//paymentServices "admin-panel-dashboard/internal/components/payments/services"
	root "admin-panel-dashboard/internal/components/root/controllers"
	// streams "admin-panel-dashboard/internal/components/streams/controllers"
	//telegram "admin-panel-dashboard/internal/components/telegram/controllers"
	careers "admin-panel-dashboard/internal/components/careers/controllers"
	generalController "admin-panel-dashboard/internal/components/general/controllers"
	health "admin-panel-dashboard/internal/components/health/controllers"
	orgControllers "admin-panel-dashboard/internal/components/organizations/controllers"
	stakeholder "admin-panel-dashboard/internal/components/stakeholder/controllers"
	swagger "admin-panel-dashboard/internal/components/swagger/controllers"
	userMetrics "admin-panel-dashboard/internal/components/usermetrics/controllers"
	users "admin-panel-dashboard/internal/components/users/controllers"
	vaultsigner "admin-panel-dashboard/internal/components/vaultsigner/controllers"
	db "admin-panel-dashboard/internal/db"
	"admin-panel-dashboard/internal/middleware"
	"admin-panel-dashboard/internal/models"
	"admin-panel-dashboard/internal/network"
	"admin-panel-dashboard/internal/observe"
	serverModels "admin-panel-dashboard/internal/server/models"
	"admin-panel-dashboard/internal/trovosdk"
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"sync"

	_ "github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

// @title Admin dashboard metrics and data retrieval
// @version 1.0.0
// @description This API provides endpoints for handling various metrics and data retrieval.(https://dashboarddev.trovotechnologies.com)
// @host
// @BasePath /api/v1
// @securityDefinitions.apikey JwtTokenAuth
// @in header
// @name Authorization
// @description JWT Token for authentication (Trovo Admin)
// @securityDefinitions.apikey OrganizationAuth
// @in header
// @name Authorization
// @description JWT Token for authentication (Organization Member)
func main() {
	// setup environment variables
	errEnv := godotenv.Load()
	if errEnv != nil {
		path, _ := os.Getwd()
		log.Printf("could not find or load any .env file from %v...skipping...\n", path)
	}

	// The admin, wallet, and P2P schemas all live in the same physical
	// database now, so this is the only connection the service opens -
	// trovoWalletDB/p2pdb are kept as separate variable names (rather
	// than threading admindb through every call site below) purely so
	// seedSuperAdmin's and Server{}'s existing signatures don't need to
	// change; all three point at the identical *gorm.DB.
	admindb, err := db.AdminDB()
	if err != nil {
		log.Fatalf("[main]Error opening AdminDB %s", err)
		return
	}
	trovoWalletDB := admindb
	p2pdb := admindb

	// superadmin
	// TOD0 Eventually get superadmin username (obi only) from .env before going live
	var superAdmins []string
	if len(os.Getenv("DEFAULT_SUPER_ADMINS")) > 0 && len(strings.Split(os.Getenv("DEFAULT_SUPER_ADMINS"), ",")) > 0 {
		superAdmins = strings.Split(strings.ReplaceAll(os.Getenv("DEFAULT_SUPER_ADMINS"), " ", ""), ",")
		log.Printf(">>>>>>>>> DEFAULT_SUPER_ADMINS set to: %v\n", os.Getenv("DEFAULT_SUPER_ADMINS"))
	} else {

		// obi, toluwase, riky
		superAdmins = []string{"obi", "toluwase"}
		log.Printf(">>>>>>>>> DEFAULT_SUPER_ADMINS ENV not SET. Using default value of: %+v\n", superAdmins)

	}
	for _, admin := range superAdmins {
		err := seedSuperAdmin(trovoWalletDB, admindb, admin)
		if err != nil {
			log.Println("[SeedSuperAdmin] error:", err)
			return
		}
	}
	SeedRoleConfig(admindb)
	err = SeedSuspensionReasons(admindb)
	if err != nil {
		log.Println("[SeedSuspensionReasons] error:", err)
		return
	}
	enableCaching := false

	if os.Getenv("ENABLE_CACHING") == "1" {
		enableCaching = true
	}
	var redisCli *redis.Client = nil
	redisSkipInsecureVerify := false
	var tlsConfig *tls.Config

	if os.Getenv("REDIS_SKIP_INSECURE_VERIFY") == "1" {
		log.Println("REDIS_SKIP_INSECURE_VERIFY set to 1. Redis will skip secure verify")
	} else {
		log.Println("REDIS_SKIP_INSECURE_VERIFY ENV not set to 1. Redis will use secure verify")
		tlsConfig = &tls.Config{
			InsecureSkipVerify: redisSkipInsecureVerify,
		}
	}
	if enableCaching {
		redisCli = redis.NewClient(&redis.Options{
			Addr:      fmt.Sprintf("%v:%v", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")), // host:port of the redis server
			Password:  os.Getenv("REDIS_PASSWORD"),                                            // no password set
			DB:        0,                                                                      // use default DB
			TLSConfig: tlsConfig,
		})
	}
	var redisCache cache.RedisCache = cache.RedisCache{
		Enabled: enableCaching,
		Client:  redisCli,
		Context: context.Background(),
	}

	// log.Println("p2p connected")
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	var m sync.Mutex
	var gc models.GlobalConfig
	gc.UsersOnline = make(map[string]models.UserMessageChannels)
	gc.GeneralUsers = make(map[string]chan map[string]interface{})
	gc.LoginUsers = make(map[string]chan map[string]interface{})
	gc.AuthUsers = make(map[string]chan map[string]interface{})

	gc.Mutex = &m
	gc.Cache = &redisCache
	// gc.DB = p2pdb
	gc.BlockchainClient = network.GetBlockchainClient()
	gc.BlockchainPassphrase = network.GetBlockchainNetworkPassPhrase()
	serviceLink, err := trovosdk.NewServiceLink()
	if err != nil {
		log.Fatalln("Could not initialize Service Link", err)
		return
	}
	gc.ServiceLink = serviceLink

	s := &serverModels.Server{
		TrovoWalletDB: trovoWalletDB,
		P2P:           p2pdb,
		AdminDB:       admindb,
		GC:            &gc,
	}
	// Crash reporting is started before the router so a panic during route
	// setup is still reported. Version matches what /health reports.
	flushObserve := observe.Init("")
	defer flushObserve()

	// gin.New() rather than gin.Default(): observe.Logger and observe.Recovery
	// replace gin's own logger and recovery middleware. Order matters -
	// RequestID must run first so the ID is available to everything after it,
	// and Recovery must wrap the handlers it protects.
	var router *gin.Engine = gin.New()
	router.Use(observe.RequestID())
	router.Use(observe.Logger())
	router.Use(observe.Recovery())
	router.Use(observe.Metrics())
	router.GET("/metrics", observe.MetricsHandler())

	router.Use(middleware.CORSMiddleware())
	s.SetupRouterParams(router)

	root.Init(router)
	log.Println("##root services initialized##")
	swagger.Init(router, s)
	log.Println("##swagger services initialized##")
	health.Init(router, s)
	log.Println("##health services initialized##")
	auth.Init(router, s)
	log.Println("##auth services initialized##")
	userMetrics.Init(router, s)
	log.Println("##userMetrics services initialized##")
	generalController.Init(router, s)
	log.Println("##generalController services initialized##")
	users.Init(router, s)
	log.Println("##users services initialized##")
	orgControllers.Init(router, s)
	log.Println("##organization services initialized##")
	careers.Init(router, s)
	log.Println("##careers services initialized##")
	stakeholder.Init(router, s)
	log.Println("##stakeholder services initialized##")
	vaultsigner.Init(router, s)
	log.Println("##vaultsigner services initialized##")

	///start server
	s.Start(router)
}
func seedSuperAdmin(dbb, adminDB *gorm.DB, superAdminUsername string) error {

	// Retrieve the user from the walletDB
	user := &models.User{Username: superAdminUsername}

	// Retrieve the user from the walletDB using the struct\
	if errr := dbb.First(user, "username = ?", superAdminUsername).Error; errr != nil {
		return fmt.Errorf("failed to retrieve user: %v", errr)
	}

	// Check if the super admin already exists
	var existingAdmin models.AdminUser
	if err := adminDB.First(&existingAdmin, "username = ?", superAdminUsername).Error; err == nil {
		// Super admin already exists, so no need to seed
		return nil
	}

	// Create the super admin user
	admin := models.AdminUser{
		Username:     user.Username,
		Email:        user.Email,
		Role:         models.SuperAdmin,
		Status:       string(models.ActiveAdmin),
		IsAdmin:      true,
		WalletUserID: user.ID,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
	}

	// Insert the new super admin into the adminDB
	if err := adminDB.Create(&admin).Error; err != nil {
		return fmt.Errorf("failed to create super admin: %v", err)
	}

	log.Printf("[seedSuperAdmin] Super admin %v created successfully\n", superAdminUsername)
	return nil
}

func SeedSuspensionReasons(db *gorm.DB) error {
	suspensionReasons := []models.UserSuspensionReason{
		{ID: 1, Reason: "Violation of terms of service"},
		{ID: 2, Reason: "Violation of community guidelines"},
		{ID: 3, Reason: "Violation of KYC/AML policy"},
		{ID: 4, Reason: "Violation of security policy"},
		{ID: 5, Reason: "Violation of privacy policy"},
		{ID: 6, Reason: "Violation of trading policy"},
		{ID: 7, Reason: "Violation of payment policy"},
		{ID: 8, Reason: "Violation of dispute resolution policy"},
		{ID: 9, Reason: "Violation of support policy"},
	}

	for _, reason := range suspensionReasons {
		var existingReason models.UserSuspensionReason
		// Check if the reason already exists in the database
		if err := db.Where("id = ?", reason.ID).First(&existingReason).Error; err != nil {
			// If not found, create a new reason
			if errors.Is(err, gorm.ErrRecordNotFound) {
				if err = db.Create(&reason).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		}
	}

	return nil
}

func SeedRoleConfig(db *gorm.DB) {
	// Define roles to be seeded
	roles := []models.RoleConfig{
		{RoleName: "SUPER_ADMIN", Description: "Full access to all system features and configurations"},
		{RoleName: "EDIT_LEVEL_ADMIN", Description: "Can edit configurations but with limited permissions"},
		{RoleName: "VIEW_ONLY_ADMIN", Description: "Can view configurations without editing permissions"},
		{RoleName: "SUSPENDED", Description: "Restricted access, user is suspended"},
	}

	// Iterate through each role and check if it exists before creating
	for _, role := range roles {
		var existingRole models.RoleConfig
		if err := db.Where("role_name = ?", role.RoleName).First(&existingRole).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				// Role does not exist, so create it
				if err := db.Create(&role).Error; err != nil {
					log.Printf("Failed to seed role '%s': %v", role.RoleName, err)
				} else {
					log.Printf("Role '%s' seeded successfully.", role.RoleName)
				}
			} else {
				log.Printf("Error checking for role '%s': %v", role.RoleName, err)
			}
		} else {
			log.Printf("Role '%s' already exists, skipping.", role.RoleName)
		}
	}
}
