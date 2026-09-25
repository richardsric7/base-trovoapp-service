package middleware

import (
	cache "admin-panel-dashboard/internal/cache"
	"time"
)

// var client redis.Client
// var ctx = context.Background()

// CreateAuth Creates/saves Authentication into the redis cache
func CreateAuth(userid string, td *TokenDetails, redisCache *cache.RedisCache) error {
	// ctx := context.Background()
	client := redisCache.Client
	at := time.Unix(td.AtExpires, 0) // converting Unix to UTC(to Time object)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()
	// fmt.Println("[CreateAuth] Setting Access Token")
	// fmt.Printf("UserID: %v, AccessTime: %v", strconv.Itoa(int(userid)), at.Sub(now))
	errAccess := client.Set(redisCache.Context, td.AccessUUID, userid, at.Sub(now)).Err()
	// fmt.Println("Saving done.....")
	if errAccess != nil {
		return errAccess
	}
	// fmt.Println("[CreateAuth] Setting Refresh Token")
	errRefresh := client.Set(redisCache.Context, td.RefreshUUID, userid, rt.Sub(now)).Err()
	if errRefresh != nil {
		return errRefresh
	}
	return nil
}
