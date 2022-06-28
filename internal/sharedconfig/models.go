package sharedconfig

import (
	"context"
	"trovo-wallet-api/internal/cache"

	"firebase.google.com/go/messaging"
	"github.com/stellar/go/clients/horizonclient"
	"gorm.io/gorm"
)

type GlobalConfig struct {
	DynamicLinkServiceURLChan chan string
	PushNotificationClient    *messaging.Client
	PNSContext                context.Context
	RedisCache                *cache.RedisCache
	DB                        *gorm.DB
	RoachDB                   *gorm.DB
	BantuExpansionClient      *horizonclient.Client
	BantuNetworkPassphrase    string
}
