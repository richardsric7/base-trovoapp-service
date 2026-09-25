package middleware

import (
	cache "admin-panel-dashboard/internal/cache"
)

// AccessDetails holds access details
type AccessDetails struct {
	AccessUUID string
	UserID     string
}

// FetchAuth is looking up the token metadata in redis
func FetchAuth(authD *AccessDetails, redisCache *cache.RedisCache) (string, error) {
	client := redisCache.Client
	userid, err := client.Get(redisCache.Context, authD.AccessUUID).Result()
	if err != nil {
		return "", err
	}

	return userid, nil
}
