package cache

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisCache struct {
	Enabled bool
	Client  *redis.Client
	Context context.Context
}

type httpResponse struct {
	Status int
	Value  interface{}
}

type cacheResult struct {
	Value interface{}
}

// CacheTelegramDeleteList only caches if Enabled.
func (r *RedisCache) CacheTelegramDeleteList(stringData string) bool {
	if !r.Enabled {
		return false
	}
	key := "delList"
	param := "telegram"
	if len(os.Getenv("CACHING_PARAMETER")) > 0 {
		param = os.Getenv("CACHING_PARAMETER")
	}
	key = key + param
	// r.Client.HSet(r.Context, key, "telegram", data).Result()
	// _, err := r.Client.Expire(r.Context, key, time.Duration(240)*time.Second).Result()
	err := r.Client.Set(r.Context, key, stringData, time.Duration(240)*time.Second).Err()
	if err != nil {
		log.Printf("[CacheTelegramDeleteList] [%v] : %v", key, err)
		return false
	}

	return true

}

// CachedTelegramDeleteList only caches if Enabled.
func (r *RedisCache) CachedTelegramDeleteList() (res bool, stringData string) {
	if !r.Enabled {
		return
	}
	param := "telegram"
	if len(os.Getenv("CACHING_PARAMETER")) > 0 {
		param = os.Getenv("CACHING_PARAMETER")
	}
	key := "delList" + param
	stringData, err := r.Client.Get(r.Context, key).Result()

	if err != nil {
		return
	}
	stringData = strings.TrimSpace(stringData)

	log.Printf("[CachedTelegramDeleteList] [%v], %v\n", key, err)

	return true, stringData

}

// CacheHttpResponseWithParameters only caches if Enabled.
func (r *RedisCache) CacheHttpResponseWithParameters(key string, parameters string, status int, response interface{}, expiryInSeconds int) bool {
	if !r.Enabled {
		return false
	}

	_httpResponse := httpResponse{
		Status: status,
		Value:  response,
	}

	bytes, err := json.Marshal(_httpResponse)

	if err != nil {
		log.Printf("[CacheHttpResponseWithParameters] failed to marshal data to bytes due to : %v\n", err)
		return false
	}
	param := "default"
	if len(os.Getenv("CACHING_PARAMETER")) > 0 {
		param = os.Getenv("CACHING_PARAMETER")
	}
	key = key + param
	_, err = r.Client.HSet(r.Context, key, parameters, bytes).Result()
	if err != nil {
		return false
	}
	_, err = r.Client.Expire(r.Context, key, time.Duration(expiryInSeconds)*time.Second).Result()

	if err != nil {
		log.Printf("[CacheHttpResponseWithParameters] failed to store [%v] in cache due to : %v\n", key, err)
		return false
	}

	return true

}

// StoreResultToCache only caches if Enabled.
func (r *RedisCache) StoreResultToCache(key string, toCache interface{}, expiryInSeconds int) bool {
	if !r.Enabled {
		return false
	}

	_cacheResult := cacheResult{
		Value: toCache,
	}
	if expiryInSeconds == 0 {
		expiryInSeconds = 120
	}
	bytes, err := json.Marshal(_cacheResult)

	if err != nil {
		log.Printf("[StoreResultToCache] failed to marshal data to bytes due to : %v\n", err)

		return false
	}
	param := "default"
	if len(os.Getenv("CACHING_PARAMETER")) > 0 {
		param = os.Getenv("CACHING_PARAMETER")
	}
	key = key + param
	_, err = r.Client.HSet(r.Context, key, param, bytes).Result()
	if err != nil {
		return false
	}
	_, err = r.Client.Expire(r.Context, key, time.Duration(expiryInSeconds)*time.Second).Result()

	if err != nil {
		log.Printf("[StoreResultToCache] failed to store [%v] in cache due to : %v\n", key, err)
		return false
	}

	return true

}

// StoreResultToCacheRaw only caches if Enabled.
func (r *RedisCache) StoreResultToCacheRaw(key string, toCache interface{}, expiryInSeconds int) bool {
	if !r.Enabled {
		return false
	}

	if expiryInSeconds == 0 {
		expiryInSeconds = 180
	}
	bytes, err := json.Marshal(toCache)

	if err != nil {
		log.Printf("[StoreResultToCacheRaw] failed to marshal data to bytes due to : %v\n", err)

		return false
	}
	param := "default"
	if len(os.Getenv("CACHING_PARAMETER")) > 0 {
		param = os.Getenv("CACHING_PARAMETER")
	}
	key = key + param
	_, err = r.Client.HSet(r.Context, key, param, bytes).Result()
	if err != nil {
		return false
	}
	_, err = r.Client.Expire(r.Context, key, time.Duration(expiryInSeconds)*time.Second).Result()

	if err != nil {
		log.Printf("[StoreResultToCacheRaw] failed to store [%v] in cache due to : %v\n", key, err)
		return false
	}

	return true

}

// GetCachedResultRaw only caches if Enabled.
func (r *RedisCache) GetCachedResultRaw(key string) (bool, []byte) {
	if !r.Enabled {
		return false, nil
	}
	param := "default"
	if len(os.Getenv("CACHING_PARAMETER")) > 0 {
		param = os.Getenv("CACHING_PARAMETER")
	}
	key = key + param

	p, err := r.Client.HGet(r.Context, key, param).Result()

	if err != nil {
		return false, nil
	}

	return true, []byte(p)

}

// GetCachedResult only caches if Enabled.
func (r *RedisCache) GetCachedResult(key string) (bool, interface{}) {
	if !r.Enabled {
		return false, ""
	}
	param := "default"
	if len(os.Getenv("CACHING_PARAMETER")) > 0 {
		param = os.Getenv("CACHING_PARAMETER")
	}
	key = key + param
	var _cachedResult cacheResult

	p, err := r.Client.HGet(r.Context, key, param).Result()

	if err != nil {
		return false, ""
	}

	err = json.Unmarshal([]byte(p), &_cachedResult)

	if err != nil {
		log.Printf("[GetCachedResult] unable to unmarshal cache result for [%v], due to %v\n", key, err)
		return false, ""
	}

	return true, _cachedResult.Value

}

// CacheHttpResponse only caches if Enabled.
func (r *RedisCache) CacheHttpResponse(key string, status int, response interface{}, expiryInSeconds int) bool {
	param := "default"
	if len(os.Getenv("CACHING_PARAMETER")) > 0 {
		param = os.Getenv("CACHING_PARAMETER")
	}
	return r.CacheHttpResponseWithParameters(key, param, status, response, expiryInSeconds)
}

// CacheHttpResponse only caches if Enabled.
func (r *RedisCache) CachedHttpResponseWithParameters(key string, parameters string) (bool, int, interface{}) {
	if !r.Enabled {
		return false, 0, ""
	}
	param := "default"
	if len(os.Getenv("CACHING_PARAMETER")) > 0 {
		param = os.Getenv("CACHING_PARAMETER")
	}
	key = key + param
	var _httpResponse httpResponse

	p, err := r.Client.HGet(r.Context, key, parameters).Result()

	if err != nil {
		return false, 0, ""
	}

	err = json.Unmarshal([]byte(p), &_httpResponse)

	if err != nil {
		log.Printf("[CachedHttpResponse] [%v], %v", key, err)
		return false, 0, ""
	}

	return true, _httpResponse.Status, _httpResponse.Value

}

// CacheHttpResponse only caches if Enabled.
func (r *RedisCache) CachedHttpResponse(key string) (bool, int, interface{}) {
	param := "default"
	if len(os.Getenv("CACHING_PARAMETER")) > 0 {
		param = os.Getenv("CACHING_PARAMETER")
	}
	return r.CachedHttpResponseWithParameters(key, param)

}

func (r *RedisCache) InvalidateCachedHttpResponse(keys ...string) bool {
	if !r.Enabled {
		return false
	}
	modKeys := make([]string, 0)
	param := "default"
	if len(os.Getenv("CACHING_PARAMETER")) > 0 {
		param = os.Getenv("CACHING_PARAMETER")
	}
	for _, i := range keys {
		modKeys = append(modKeys, i+param)
	}
	log.Printf("[InvalidateCachedHttpResponse] [%v]\n", modKeys)

	_, err := r.Client.Del(r.Context, modKeys...).Result()

	return err == nil

}

func (r *RedisCache) DeleteFromCache(keys ...string) bool {
	if !r.Enabled {
		return false
	}
	modKeys := make([]string, 0)
	param := "default"
	if len(os.Getenv("CACHING_PARAMETER")) > 0 {
		param = os.Getenv("CACHING_PARAMETER")
	}
	for _, i := range keys {
		modKeys = append(modKeys, i+param)
	}
	log.Printf("[DeleteFromCache] [%v]\n", modKeys)

	_, err := r.Client.Del(r.Context, modKeys...).Result()

	return err == nil

}
