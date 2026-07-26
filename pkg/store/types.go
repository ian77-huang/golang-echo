package store

import (
	storeCache "github.com/ian77-huang/golang-echo/pkg/store/cache"
	storeRedis "github.com/ian77-huang/golang-echo/pkg/store/redis"
)

type RedisOption = storeRedis.RedisOption
type RedisMSET = storeRedis.RedisMSET

type CacheOption = storeCache.CacheOption

type StoreOptionRedis struct {
	IsUse      bool
	Option     *RedisOption
	ConnectURL string
}

type StoreOption struct {
	Redis *StoreOptionRedis
	Cache *CacheOption
}

type StoreServer struct {
	cache *storeCache.CacheServer
	redis *storeRedis.RedisServer
}

type StoreValue struct {
	Val interface{}
}
