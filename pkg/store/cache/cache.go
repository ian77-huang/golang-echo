package cache

import (
	"time"

	"github.com/patrickmn/go-cache"
)

const NoExpiration time.Duration = -1

func New(co *CacheOption) *CacheServer {
	if co == nil {
		co = &CacheOption{}
	}
	if co.DefaultExpiration == 0 {
		co.DefaultExpiration = cache.NoExpiration
	}
	if co.CleanupInterval == 0 {
		co.CleanupInterval = cache.NoExpiration
	}

	return &CacheServer{cache: cache.New(co.DefaultExpiration, co.CleanupInterval)}
}

func (ca *CacheServer) Set(key string, val interface{}, expiration time.Duration) {
	ca.cache.Set(key, val, expiration)
}

func (ca *CacheServer) Get(key string) (interface{}, bool) {
	return ca.cache.Get(key)
}

func (ca *CacheServer) Delete(key string) {
	ca.cache.Delete(key)
}
