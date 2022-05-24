package cache

import (
	"context"
	"encoding/json"
	"log"
	"os"
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

// CacheHttpResponse only caches if Enabled.
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
		return false
	}
	param := "default"
	if len(os.Getenv("CACHING_PARAMETER")) > 0 {
		param = os.Getenv("CACHING_PARAMETER")
	}
	key = key + param
	r.Client.HSet(r.Context, key, parameters, bytes).Result()
	_, err = r.Client.Expire(r.Context, key, time.Duration(expiryInSeconds)*time.Second).Result()

	if err != nil {
		log.Printf("[CacheHttpResponse] [%v] : %v", key, err)
		return false
	}

	return true

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
