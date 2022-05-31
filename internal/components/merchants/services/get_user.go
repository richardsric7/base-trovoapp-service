package merchants

import (
	"trovo-wallet-api/internal/cache"
	merchantUserModels "trovo-wallet-api/internal/components/merchants/models"

	"gorm.io/gorm"
)

//GetUser gets user information
func GetUser(ID string, db *gorm.DB, publicKey string, dynamicLinkServiceUrlChan chan string, redisCache *cache.RedisCache) (userInfo merchantUserModels.User, err error) {

	return

}
