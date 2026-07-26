package cache

import (
	"time"

	"github.com/patrickmn/go-cache"
)

type CacheOption struct {
	DefaultExpiration time.Duration
	CleanupInterval   time.Duration
}

type CacheServer struct {
	cache *cache.Cache
}
