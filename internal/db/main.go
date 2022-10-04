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
	SMS "trovo-wallet-api/internal/sms"

	"github.com/ecnepsnai/discord"
	sqliteEncrypt "github.com/jackfr0st13/gorm-sqlite-cipher"
	"gorm.io/driver/postgres"

	// "gorm.io/driver/sqlite"
	"gorm.io/gorm"
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

		errMigrate = gormDB.AutoMigrate(&users.UserWallet{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error Migrating UserWallet: ", errMigrate)
		}
		errMigrate = gormDB.AutoMigrate(&users.UserWalletSharedAccess{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error Migrating UserWalletSharedAccess: ", errMigrate)
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

		errMigrate = gormDB.AutoMigrate(&users.Permissions{})
		if errMigrate != nil {
			log.Fatalln("[OpenDb]Error Migrating Permissions: ", errMigrate)
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
