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
			if err != nil {
				return
			}

			sqlDB.SetMaxIdleConns(int(maxIdleConns))
			sqlDB.SetConnMaxIdleTime(5 * time.Second)
			sqlDB.SetMaxOpenConns(int(maxOpenConn))
			sqlDB.SetConnMaxLifetime(24 * time.Hour)
		})
		if err != nil {
			log.Printf("[OpenDb]failed to open sql connection, %s\n", err)
			return nil, err
		}

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

func PrintDBStats(tag string, db *gorm.DB) {
	sqlDB, _ := db.DB()
	dbStats := sqlDB.Stats()
	log.Printf("[%v], open connections: %v, Max Open Conns: %v, idle connections: %v, Max Idle Closed: %v", tag, dbStats.OpenConnections, dbStats.MaxOpenConnections, dbStats.Idle, dbStats.MaxIdleClosed)
	if dbStats.OpenConnections > (dbStats.MaxOpenConnections - 25) {
		logDiscordDBWarning(fmt.Sprintf("[DB CONNECTION FILLUP WARNING] open connections have reached %v/%v. increase limit!", dbStats.OpenConnections, dbStats.MaxOpenConnections))
	}
}

func logDiscordDBWarning(msg string) {
	// fallback default, used only when CONNECTION_WARNING_WEBHOOK isn't configured for this environment.
	// NOTE: this webhook has been committed to source history - rotate it in Discord's channel settings.
	discord.WebhookURL = "https://discord.com/api/webhooks/827986576415129663/wqMKp9wxB_fxs9Q3zlMKCNPGENXmD_ueUnL8hVCu1wmRfD2wkXAjfP85k1Ro_2_wGfiY"
	if v := os.Getenv("CONNECTION_WARNING_WEBHOOK"); len(v) > 50 {
		discord.WebhookURL = v
	}
	discord.Say(msg)
}
