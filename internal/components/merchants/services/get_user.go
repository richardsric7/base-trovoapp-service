package merchants

import (
	"trovo-wallet-api/internal/cache"
	userModels "trovo-wallet-api/internal/components/users/models"
	conDB "trovo-wallet-api/internal/db"

	"gorm.io/gorm"
)

//GetUser gets user information
func GetUser(ID string, db *gorm.DB, publicKey string, dynamicLinkServiceUrlChan chan string, redisCache *cache.RedisCache) (userInfo userModels.User, err error) {
	conDB.PrintDBStats("GetUser", db)

	return

}
