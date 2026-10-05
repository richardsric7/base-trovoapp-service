package store

import (
	"errors"
	"os"
	"strconv"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Open connects to app-backend's database (DB_TYPE postgres, the default,
// or sqlite for local runs; DB_CONNECTION_STRING).
func Open() (*gorm.DB, error) {
	dsn := os.Getenv("DB_CONNECTION_STRING")
	if dsn == "" {
		return nil, errors.New("DB_CONNECTION_STRING is not set")
	}
	cfg := &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)}
	if os.Getenv("DB_TYPE") == "sqlite" {
		return gorm.Open(sqlite.Open(dsn), cfg)
	}
	db, err := gorm.Open(postgres.Open(dsn), cfg)
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	maxOpen, _ := strconv.Atoi(os.Getenv("DB_MAX_OPEN_CONNECTIONS"))
	if maxOpen <= 0 {
		maxOpen = 10
	}
	sqlDB.SetMaxOpenConns(maxOpen)
	sqlDB.SetMaxIdleConns(maxOpen)
	sqlDB.SetConnMaxLifetime(time.Hour)
	return db, nil
}
