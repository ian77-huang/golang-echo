package cache

import (
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
)

var (
	appCache *cache.Cache
	once     sync.Once
)

func New() *cache.Cache {
	once.Do(func() {
		appCache = cache.New(cache.NoExpiration, cache.NoExpiration)
	})
	return appCache
}

func SetCache(key string, value interface{}, expiration time.Duration) {
	appCache.Set(key, value, expiration)
}

func GetCache(key string) (interface{}, bool) {
	return appCache.Get(key)
}
