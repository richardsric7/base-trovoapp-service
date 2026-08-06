package middleware

import (
	"log"
	cache "trovo-wallet-api/internal/cache"
	tErrors "trovo-wallet-api/internal/errors"
)

// DeleteAuth delete the UUID from redis server
func DeleteAuth(givenUUID string, redisCache *cache.RedisCache) (int64, error) {
	client := redisCache.Client
	deleted, err := client.Del(redisCache.Context, givenUUID).Result()
	if err != nil {
		log.Printf("unable to delete token [%v], error: %v", givenUUID, err)

		return 0, &tErrors.ErrorInvalidAuthorization{}
	}
	if deleted == 0 {
		return 0, &tErrors.ErrorInvalidAuthorization{}
	}
	return deleted, nil
}
