package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
	announcementModels "trovo-wallet-api/internal/components/announcements/models"
	assetModels "trovo-wallet-api/internal/components/assets/models"
	paymentModels "trovo-wallet-api/internal/components/payments/models"
	servicelinkModels "trovo-wallet-api/internal/components/servicelinks/models"
	users "trovo-wallet-api/internal/components/users/models"
	"trovo-wallet-api/internal/dynamiclinks"
	SMS "trovo-wallet-api/internal/sms"

	"github.com/ecnepsnai/discord"
	sqliteEncrypt "github.com/jackfr0st13/gorm-sqlite-cipher"
	"gorm.io/driver/postgres"

	// "gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

// OpenDb method should only run once to reuse the same connection pool
// https://golang.org/pkg/database/sql/#Open
// The returned DB is safe for concurrent use by multiple goroutines and maintains its own pool of idle connections.
// Thus, the Open function should be called just once. It is rarely necessary to close a DB.
var gormDB, roachDB *gorm.DB
var sqlDB *sql.DB // Set package-wide, but not exported
var once sync.Once

func OpenDb() (*gorm.DB, error) {
	dbType := os.Getenv("DB_TYPE")
	if len(dbType) == 0 {
		dbType = "postgres"
		log.Println("ENV DB_TYPE not set, using default postgres")
	}
	// if len(dbType) == 0 {
	// 	return nil, errors.New("empty db type. Check env variable: DB_TYPE")
	// }
	dbConnectionString := os.Getenv("DB_CONNECTION_STRING")
	if len(dbConnectionString) == 0 {
		return nil, errors.New("empty connection string. Check env variable: DB_CONNECTION_STRING")
	}

	var err error

	// if dbType == "sqlite" {
	// 	gormDB, err = gorm.Open(sqlite.Open(dbConnectionString), &gorm.Config{
	// 		QueryFields: true,
	// 	})
	// }

	if dbType == "postgres" {
		var maxoconn, idleCon string
		if os.Getenv("DB_MAX_OPEN_CONNECTIONS") != "" {
			maxoconn = os.Getenv("DB_MAX_OPEN_CONNECTIONS")
		} else {
			maxoconn = "50"
		}
		if os.Getenv("DB_MAX_IDLE_CONNECTIONS") != "" {
			idleCon = os.Getenv("DB_MAX_IDLE_CONNECTIONS")
		} else {
			idleCon = "50"
		}
		maxOpenConn, _ := strconv.ParseInt(maxoconn, 10, 64)
		maxIdleConns, _ := strconv.ParseInt(idleCon, 10, 64)
		// connsMaxIdleTime, _ := strconv.ParseInt(idletime, 10, 64)
		once.Do(func() {
			sqlDB, err = sql.Open("pgx", dbConnectionString)

			sqlDB.SetMaxIdleConns(int(maxIdleConns))
			sqlDB.SetConnMaxIdleTime(5 * time.Second)
			sqlDB.SetMaxOpenConns(int(maxOpenConn))
			sqlDB.SetConnMaxLifetime(24 * time.Hour)
		})

		gormDB, err = gorm.Open(postgres.New(postgres.Config{
			Conn: sqlDB,
		}), &gorm.Config{
			Logger:      logger.Default.LogMode(logger.Silent),
			QueryFields: true,
		})
		if err != nil {
			log.Printf("[OpenDb]failed to connect database, %s\n", err)
			return nil, err
		}

	}
	return gormDB, nil
}

func OpenRoachDB() (*gorm.DB, error) {
	var err error

	roachDB, err = gorm.Open(postgres.Open(os.Getenv("CDB_CONNECTION_STRING")), &gorm.Config{
		Logger:      logger.Default.LogMode(logger.Silent),
		QueryFields: true,
	})
	if err != nil {
		log.Fatal(err)
	}
	return roachDB, nil
}

// OpenSqliteDB opens ecnrypted SQlite connection
func OpenSqliteDB() (*gorm.DB, error) {

	var errDB error
	key := "746373408hhgdf#^hf*bhe)8"
	if os.Getenv("SQLITE_CYPER_PASSPHRASE") != "" {
		key = os.Getenv("SQLITE_CYPER_PASSPHRASE")
	}
	dbname := "dbs/bantupay.sqlite"
	dbnameWithDSN := dbname + fmt.Sprintf("?_pragma_key=%s&_pragma_cipher_page_size=4096", key)
	sqliteDB, errDB := gorm.Open(sqliteEncrypt.Open(dbnameWithDSN), &gorm.Config{
		Logger:      logger.Default.LogMode(logger.Silent),
		QueryFields: true,
	})

	if errDB != nil {
		log.Printf("[OpenDb]failed to connect database, %s\n", errDB)
		return nil, errDB
	}

	return sqliteDB, nil
}

func MigrateDB(gormDB *gorm.DB) {
	if os.Getenv("DB_AUTOMIGRATE") == "1" {
		errMigrate := gormDB.AutoMigrate(&users.User{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error migrating User:", errMigrate)
		}
		errMigrate = gormDB.AutoMigrate(&users.DeletedUserAccount{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error Migrating DeletedUserAccount: ", errMigrate)
		}
		errMigrate = gormDB.AutoMigrate(&users.UserWallet{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error Migrating UserWallet: ", errMigrate)
		}
		errMigrate = gormDB.AutoMigrate(&users.Bank{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error Migrating Bank: ", errMigrate)
		}
		errMigrate = gormDB.AutoMigrate(&users.Country{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error Migrating Country: ", errMigrate)
		}
		errMigrate = gormDB.AutoMigrate(&users.KYCConfig{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error Migrating KYCConfig: ", errMigrate)
		}
		errMigrate = gormDB.AutoMigrate(&users.KYCLevel{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error Migrating KYCLevel: ", errMigrate)
		}
		errMigrate = gormDB.AutoMigrate(&users.SumSubReviewResult{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error Migrating SumSubReviewResult: ", errMigrate)
		}
		errMigrate = gormDB.AutoMigrate(&users.UserKYCProgress{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error Migrating UserKYCProgress: ", errMigrate)
		}
		errMigrate = gormDB.AutoMigrate(&users.UserFiatPaymentMethod{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error Migrating UserFiatPaymentMethod: ", errMigrate)
		}
		errMigrate = gormDB.AutoMigrate(&users.ClosedGroup{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error Migrating ClosedGroup: ", errMigrate)
		}
		errMigrate = gormDB.AutoMigrate(&users.UserClosedGroup{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error Migrating UserClosedGroup: ", errMigrate)
		}
		errMigrate = gormDB.AutoMigrate(&users.TokenizedAsset{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error Migrating TokenizedAsset: ", errMigrate)
		}
		errMigrate = gormDB.AutoMigrate(&users.TokenizedAssetSubscription{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error Migrating TokenizedAssetSubscription: ", errMigrate)
		}
		errMigrate = gormDB.AutoMigrate(&users.ExpressionOfInterest{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error Migrating ExpressionOfInterest: ", errMigrate)
		}

		errMigrate = gormDB.AutoMigrate(&users.ProceedPayout{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error Migrating ProceedPayout: ", errMigrate)
		}

		errMigrate = gormDB.AutoMigrate(&users.TokenizedAssetPayoutSchedule{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error Migrating TokenizedAssetPayoutSchedule: ", errMigrate)
		}

		errMigrate = gormDB.AutoMigrate(&users.TokenizedAssetPayoutEngineTask{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error Migrating TokenizedAssetPayoutEngineTask: ", errMigrate)
		}

		errMigrate = gormDB.AutoMigrate(&users.UserAccountRecoveryLog{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error Migrating UserAccountRecoveryLog: ", errMigrate)
		}
		errMigrate = gormDB.AutoMigrate(&users.PatronPackage{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error Migrating PatronPackage: ", errMigrate)
		}
		errMigrate = gormDB.AutoMigrate(&users.PatronTier{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error Migrating PatronTier: ", errMigrate)
		}
		errMigrate = gormDB.AutoMigrate(&users.UserPatronMembership{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error Migrating UserPatronMembership: ", errMigrate)
		}
		errMigrate = gormDB.AutoMigrate(&users.PatronMembershipGrade{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error Migrating PatronMembershipGrade: ", errMigrate)
		}
		errMigrate = gormDB.AutoMigrate(&users.UserPatronSubscriptionLog{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error Migrating UserPatronSubscriptionLog: ", errMigrate)
		}
		errMigrate = gormDB.AutoMigrate(&users.PatronSubscriptionPaymentAsset{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error Migrating PatronSubscriptionPaymentAsset: ", errMigrate)
		}
		errMigrate = gormDB.AutoMigrate(&users.PendingAuth{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error migrating PendingAuth:", errMigrate)
		}

		errMigrate = gormDB.AutoMigrate(&users.PendingTransactionSignature{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error migrating PendingTransactionSignature:", errMigrate)
		}

		errMigrate = gormDB.AutoMigrate(&users.CryptoWalletDepositAddress{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error Migrating CryptoWalletDepositAddress: ", errMigrate)
		}
		errMigrate = gormDB.AutoMigrate(&users.CryptoDeposit{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error Migrating CryptoDeposit: ", errMigrate)
		}
		errMigrate = gormDB.AutoMigrate(&users.CallbackDepositItem{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error Migrating CallbackDepositItem: ", errMigrate)
		}
		errMigrate = gormDB.AutoMigrate(&users.WithdrawalNetwork{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error Migrating WithdrawalNetwork: ", errMigrate)
		}
		errMigrate = gormDB.AutoMigrate(&users.CryptoWithdrawal{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error Migrating CryptoWithdrawal: ", errMigrate)
		}
		errMigrate = gormDB.AutoMigrate(&users.WithdrawalRequest{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error Migrating WithdrawalRequest: ", errMigrate)
		}
		errMigrate = gormDB.AutoMigrate(&users.WalletPermission{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error Migrating WalletPermission: ", errMigrate)
		}
		errMigrate = gormDB.AutoMigrate(&users.SecurityQuestion{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error Migrating SecretQuestion: ", errMigrate)
		}
		errMigrate = gormDB.AutoMigrate(&users.UserSecurityAnswer{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error Migrating UserSecretAnswer: ", errMigrate)
		}

		errMigrate = gormDB.AutoMigrate(&users.Permission{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error Migrating Permission: ", errMigrate)
		}

		errMigrate = gormDB.AutoMigrate(&users.ReservedName{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error Migrating ReservedName: ", errMigrate)
		}

		errMigrate = gormDB.AutoMigrate(&users.UserAccountRecoveryEmailVerification{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error Migrating UserAccountRecoveryEmailVerification: ", errMigrate)
		}

		errMigrate = gormDB.AutoMigrate(&users.UserMobilePhoneVerification{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error Migrating UserMobilePhoneVerification: ", errMigrate)
		}

		errMigrate = gormDB.AutoMigrate(&users.MarketOffer{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error Migrating MarketOffer: ", errMigrate)
		}

		errMigrate = gormDB.AutoMigrate(&servicelinkModels.ServiceLink{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error Migrating ServiceLink: ", errMigrate)
		}
		errMigrate = gormDB.AutoMigrate(&servicelinkModels.ServiceLinkApiKeyLog{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error Migrating MerchantApiKeyLog: ", errMigrate)
		}
		errMigrate = gormDB.AutoMigrate(&servicelinkModels.ServiceLinkLoginSession{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error Migrating ServiceLinkLoginSession: ", errMigrate)
		}
		errMigrate = gormDB.AutoMigrate(&servicelinkModels.ServiceLinkAuthorization{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error Migrating ServiceLinkAuthorization: ", errMigrate)
		}
		errMigrate = gormDB.AutoMigrate(&servicelinkModels.ServiceLinkEvent{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error Migrating ServiceLinkEvent: ", errMigrate)
		}

		errMigrate = gormDB.AutoMigrate(&SMS.SmsProvider{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error Migrating SmsProvider: ", errMigrate)
		}

		errMigrate = gormDB.AutoMigrate(&users.DefaultAsset{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error Migrating DefaultAsset: ", errMigrate)
		}

		errMigrate = gormDB.AutoMigrate(&paymentModels.PaymentHistory{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error Migrating PaymentHistory: ", errMigrate)
		}

		errMigrate = gormDB.AutoMigrate(&paymentModels.CurrencyRates{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error Migrating CurrencyRates: ", errMigrate)
		}

		errMigrate = gormDB.AutoMigrate(&assetModels.AssetClass{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error Migrating AssetClass: ", errMigrate)
		}

		errMigrate = gormDB.AutoMigrate(&users.TokenizedAssetSector{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error Migrating TokenizedAssetSector: ", errMigrate)
		}

		errMigrate = gormDB.AutoMigrate(&users.TokenizationStatus{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error Migrating TokenizationStatus: ", errMigrate)
		}

		errMigrate = gormDB.AutoMigrate(&users.TokenizedAssetSubSector{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error Migrating TokenizedAssetSubSector: ", errMigrate)
		}

		errMigrate = gormDB.AutoMigrate(&users.TokenizedAssetType{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error Migrating TokenizedAssetType: ", errMigrate)
		}

		errMigrate = gormDB.AutoMigrate(&users.TokenizationFeePaymentMethod{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error Migrating TokenizationFeePaymentMethod: ", errMigrate)
		}

		errMigrate = gormDB.AutoMigrate(&users.TokenizationFeeProofOfPayment{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error Migrating TokenizationFeeProofOfPayment: ", errMigrate)
		}

		errMigrate = gormDB.AutoMigrate(&users.AssetTokenizationDocumentType{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error Migrating AssetTokenizationDocumentType: ", errMigrate)
		}
		errMigrate = gormDB.AutoMigrate(&users.ApprovedAssetCustodian{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error Migrating ApprovedAssetCustodian: ", errMigrate)
		}
		errMigrate = gormDB.AutoMigrate(&users.AssetManager{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error Migrating AssetManager: ", errMigrate)
		}

		errMigrate = gormDB.AutoMigrate(&users.AssetIssuingHouse{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error Migrating AssetIssuingHouse: ", errMigrate)
		}

		errMigrate = gormDB.AutoMigrate(&users.LegalAndProfesionalPartner{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error Migrating LegalAndProfesionalPartner: ", errMigrate)
		}

		errMigrate = gormDB.AutoMigrate(&users.RatingAgency{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error Migrating RatingAgency: ", errMigrate)
		}

		errMigrate = gormDB.AutoMigrate(&users.Trustee{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error Migrating Trustee: ", errMigrate)
		}
		errMigrate = gormDB.AutoMigrate(&users.TokenizationMintingApprover{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error Migrating TokenizationMintingApprover: ", errMigrate)
		}
		errMigrate = gormDB.AutoMigrate(&users.TokenizationMintingInitiator{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error Migrating TokenizationMintingInitiator: ", errMigrate)
		}
		errMigrate = gormDB.AutoMigrate(&users.ExistingAssetValidationAssetInformation{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error Migrating ExistingAssetValidationAssetInformation: ", errMigrate)
		}
		errMigrate = gormDB.AutoMigrate(&users.ExistingAssetValidationAssetDocument{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error Migrating ExistingAssetValidationAssetDocument: ", errMigrate)
		}
		errMigrate = gormDB.AutoMigrate(&users.ExistingAssetValidationAssetTokenInfo{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error Migrating ExistingAssetValidationAssetTokenInfo: ", errMigrate)
		}
		errMigrate = gormDB.AutoMigrate(&users.NonExistingAssetValidationAssetInformation{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error Migrating NonExistingAssetValidationAssetInformation: ", errMigrate)
		}
		errMigrate = gormDB.AutoMigrate(&users.NonExistingAssetValidationAssetDocument{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error Migrating NonExistingAssetValidationAssetDocument: ", errMigrate)
		}
		errMigrate = gormDB.AutoMigrate(&users.NonExistingAssetValidationAssetTokenInfo{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error Migrating NonExistingAssetValidationAssetTokenInfo: ", errMigrate)
		}
		errMigrate = gormDB.AutoMigrate(&users.TokenizationCurrency{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error Migrating TokenizationCurrency: ", errMigrate)
		}
		errMigrate = gormDB.AutoMigrate(&users.TokenizationFee{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error Migrating TokenizationFee: ", errMigrate)
		}
		errMigrate = gormDB.AutoMigrate(&users.AssetProtectionOption{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error Migrating AssetProtectionOption: ", errMigrate)
		}
		errMigrate = gormDB.AutoMigrate(&users.AssetTokenizationDocument{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error Migrating AssetTokenizationDocument: ", errMigrate)
		}
		errMigrate = gormDB.AutoMigrate(&users.TokenizationPublicAssetAllowedCountryCode{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error Migrating TokenizationPublicAssetAllowedCountryCode: ", errMigrate)
		}
		errMigrate = gormDB.AutoMigrate(&users.ProceedCycle{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error Migrating ProceedCycle: ", errMigrate)
		}

		errMigrate = gormDB.AutoMigrate(&users.PostTokenizationTrustlineCandidate{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error Migrating PostTokenizationTrustlineCandidate: ", errMigrate)
		}

		dberr := gormDB.First(&assetModels.AssetClass{}).Error
		if errors.Is(dberr, gorm.ErrRecordNotFound) {
			assetClasses := []assetModels.AssetClass{{AssetClass: "Token"}, {AssetClass: "Stablecoin"}, {AssetClass: "Tokenized Asset"}, {AssetClass: "Non Fungible Token (NFT)"}, {AssetClass: "Reward"}}
			gormDB.Omit(clause.Associations).Create(&assetClasses)
		}

		errMigrate = gormDB.AutoMigrate(&assetModels.CuratedAsset{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error Migrating CuratedAsset: ", errMigrate)
		}

		errMigrate = gormDB.AutoMigrate(&assetModels.XbnDollarPrice{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error Migrating XbnDollarPrice: ", errMigrate)
		}

		errMigrate = gormDB.AutoMigrate(&announcementModels.Announcement{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error Migrating Announcement: ", errMigrate)
		}
		errMigrate = gormDB.AutoMigrate(&announcementModels.AppVersion{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error Migrating AppVersion: ", errMigrate)
		}

		errMigrate = gormDB.AutoMigrate(&users.KycWebhookRequest{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error Migrating KycWebhookRequest: ", errMigrate)
		}
		errMigrate = gormDB.AutoMigrate(&users.DojaWidget{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error Migrating DojaWidget: ", errMigrate)
		}
		errMigrate = gormDB.AutoMigrate(&users.UserDojaKYCProgress{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error Migrating UserDojaKYCProgress: ", errMigrate)
		}

		errMigrate = gormDB.AutoMigrate(&users.FiatPaymentConfig{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error Migrating FiatPaymentConfig: ", errMigrate)
		}

		errMigrate = gormDB.AutoMigrate(&users.PaymentWebhookRequest{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error Migrating PaymentWebhookRequest: ", errMigrate)
		}

		errMigrate = gormDB.AutoMigrate(&users.FiatPayment{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error Migrating FiatPayment: ", errMigrate)
		}

		errMigrate = gormDB.AutoMigrate(&users.FiatPaymentInvoice{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error Migrating FiatPaymentInvoice: ", errMigrate)
		}

		errMigrate = gormDB.AutoMigrate(&users.FaucetConfig{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error Migrating FaucetConfig: ", errMigrate)
		}

		errMigrate = gormDB.AutoMigrate(&dynamiclinks.DynamicLink{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error Migrating DynamicLink: ", errMigrate)
		}

		errMigrate = gormDB.AutoMigrate(&users.JsonForm{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error Migrating JsonForm: ", errMigrate)
		}

		// errMigrate = UserTriggers(gormDB)
		// if errMigrate != nil {
		// 	log.Fatalln("[OpenDb]Error Migrating User Triggers: ", errMigrate)
		// }
	}

}

func PrintDBStats(tag string, db *gorm.DB) {
	sqlDB, _ := db.DB()
	dbStats := sqlDB.Stats()
	log.Printf("[%v], open connections: %v, Max Open Conns: %v, idle connections: %v, Max Idle Closed: %v", tag, dbStats.OpenConnections, dbStats.MaxOpenConnections, dbStats.Idle, dbStats.MaxIdleClosed)
	if dbStats.OpenConnections > (dbStats.MaxOpenConnections - 25) {
		logDiscordDBWarning(fmt.Sprintf("[DB CONNECTION FILLUP WARNING] open connections have reached %v/%v. increase limit!", dbStats.OpenConnections, dbStats.MaxOpenConnections))
	}
}

func logDiscordDBWarning(msg string) {
	discord.WebhookURL = "https://discord.com/api/webhooks/827986576415129663/wqMKp9wxB_fxs9Q3zlMKCNPGENXmD_ueUnL8hVCu1wmRfD2wkXAjfP85k1Ro_2_wGfiY"
	if len(os.Getenv("CONNECTION_WARNING_WEBHOOK")) > 50 {
		discord.WebhookURL = os.Getenv("CONNECTION_WARNING_WEBHOOK")
	}
	discord.Say(msg)
}

// UserTriggers executes creates users triggers and functions in users table
func UserTriggers(db *gorm.DB) error {

	trigger := `DROP TRIGGER IF EXISTS users_au ON users;

	CREATE OR REPLACE FUNCTION users_aut_func()
	   RETURNS TRIGGER 
	   LANGUAGE PLPGSQL
	AS $$
	BEGIN
	   IF NEW.suspended = 1 AND OLD.suspended = 0 THEN
	   INSERT INTO banned_public_keys (public_key) VALUES(NEW.public_key) ON CONFLICT DO NOTHING;
	   END IF;
	   IF NEW.suspended = 0 AND OLD.suspended = 1 THEN
	   DELETE FROM banned_public_keys WHERE public_key = NEW.public_key;
	   END IF;
	   IF NEW.mobile != OLD.mobile THEN
	   DELETE FROM user_mobile_phone_verifications WHERE user_id = NEW.id;
	   END IF;
	   RETURN NEW;
	END;
	$$
	;
	CREATE TRIGGER users_au
	  AFTER UPDATE
	  ON users
	  FOR EACH ROW
	  EXECUTE PROCEDURE users_aut_func();


	  DROP TRIGGER IF EXISTS users_bi ON users;
	  CREATE OR REPLACE FUNCTION users_bit_func()
	  RETURNS TRIGGER 
	  LANGUAGE PLPGSQL
   AS $$

   DECLARE
   phone_prefix   character varying (10);

   BEGIN
	  IF NEW.isp IN (select isp from high_risk_isps where score > 9) THEN
	  RAISE EXCEPTION 'isp threat score > 9';
	  END IF;
	  IF split_part(NEW.email,'@',2) IN (select domain from medium_risk_domains where score > 9) THEN
	  RAISE EXCEPTION 'email domain threat score > 9';
	  END IF;
	  IF NEW.city IN (select city from high_risk_cities where score > 9) THEN
	  RAISE EXCEPTION 'city threat score > 9';
	  END IF;
	  IF NEW.public_ip IN (select ip_address from high_risk_ips where score > 9 and (ip_address = NEW.public_ip or ip_address = split_part(NEW.public_ip,'.',1) or ip_address = split_part(NEW.public_ip,'.',1) || '.'|| split_part(NEW.public_ip,'.',2) or ip_address = split_part(NEW.public_ip,'.',1) || '.'|| split_part(NEW.public_ip,'.',2) || '.' || split_part(NEW.public_ip,'.',3)) limit 1) THEN
	  RAISE EXCEPTION 'ip threat score > 9';
	  END IF;

	  RETURN NEW;
   END;
   $$
   ;
   CREATE TRIGGER users_bi
	 BEFORE INSERT
	 ON users
	 FOR EACH ROW
	 EXECUTE PROCEDURE users_bit_func();


	DROP TRIGGER IF EXISTS users_bu ON users;
	  CREATE OR REPLACE FUNCTION users_but_func()
	  RETURNS TRIGGER 
	  LANGUAGE PLPGSQL
   AS $$

   DECLARE
   phone_prefix   character varying (10);

   BEGIN
	  IF NEW.isp IN (select isp from high_risk_isps where score > 9) THEN
	  RAISE EXCEPTION 'isp threat score > 9';
	  END IF;
	  IF split_part(NEW.email,'@',2) IN (select domain from medium_risk_domains where score > 9) THEN
	  RAISE EXCEPTION 'email domain threat score > 9';
	  END IF;
	  IF NEW.city IN (select city from high_risk_cities where score > 9) THEN
	  RAISE EXCEPTION 'city threat score > 9';
	  END IF;
	  IF NEW.public_ip IN (select ip_address from high_risk_ips where score > 9 and (ip_address = NEW.public_ip or ip_address = split_part(NEW.public_ip,'.',1) or ip_address = split_part(NEW.public_ip,'.',1) || '.'|| split_part(NEW.public_ip,'.',2) or ip_address = split_part(NEW.public_ip,'.',1) || '.'|| split_part(NEW.public_ip,'.',2) || '.' || split_part(NEW.public_ip,'.',3)) limit 1) THEN
	  RAISE EXCEPTION 'ip threat score > 9';
	  END IF;


	  IF NEW.mobile != OLD.mobile THEN
	   NEW.mobile_verified = 0;
	  END IF;

	  RETURN NEW;
   END;
   $$
   ;
   CREATE TRIGGER users_bu
	 BEFORE UPDATE
	 ON users
	 FOR EACH ROW
	 EXECUTE PROCEDURE users_but_func();
	`
	return db.Exec(trigger).Error
}
