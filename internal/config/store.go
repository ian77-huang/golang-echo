package config

import (
	"time"

	"github.com/ian77-huang/golang-echo/pkg/store"
	"github.com/redis/go-redis/v9"
)

func Store(redisURL string) *store.StoreServer {
	server := store.New(&store.StoreOption{
		Cache: &store.CacheOption{
			DefaultExpiration: 5 * time.Minute,
			CleanupInterval:   30 * time.Minute,
		},
		Redis: &store.StoreOptionRedis{
			IsUse:      true,
			ConnectURL: redisURL,
		},
	})

	server.NewRedisPubSub(func(msg *redis.Message) {

	})

	return server
}
