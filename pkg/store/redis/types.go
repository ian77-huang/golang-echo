package storeRedis

import "github.com/redis/go-redis/v9"

type RedisOption = redis.Options

type RedisServer struct {
	redis *redis.Client
}
