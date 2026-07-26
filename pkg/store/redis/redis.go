package storeRedis

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

const defaultTimeout = 5 * time.Second

var errNotInitialized = errors.New("redis client is not initialized")

func New(ro *RedisOption) *RedisServer {
	if ro == nil {
		ro = &RedisOption{}
	}
	if ro.Addr == "" {
		ro.Addr = "localhost:6379"
	}

	return &RedisServer{redis: redis.NewClient(ro)}
}

func NewURL(url string) (*RedisServer, error) {
	opt, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}

	redisClient := redis.NewClient(opt)

	return &RedisServer{redis: redisClient}, nil
}

func (rs *RedisServer) NewPubSub(channel string, res func(msg *redis.Message)) context.CancelFunc {
	if rs == nil || rs.redis == nil || res == nil {
		return func() {}
	}
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		pubsub := rs.redis.Subscribe(ctx, channel)
		defer pubsub.Close()

		if _, err := pubsub.Receive(ctx); err != nil {
			log.Printf("redis pubsub subscribe failed: %v", err)
			cancel()
			return
		}
		ch := pubsub.Channel()

		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-ch:
				if !ok {
					return
				}
				res(msg)
			}
		}
	}()
	return cancel
}

func (rs *RedisServer) Close() error {
	if rs == nil || rs.redis == nil {
		return nil
	}
	return rs.redis.Close()
}

func newDefaultContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), defaultTimeout)
}

func newTimeoutContext(timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), timeout)
}

type ResultBytes struct {
	Ttl time.Duration
	Val []byte
}
type KeyCmds struct {
	GetCmd *redis.StringCmd
	TtlCmd *redis.DurationCmd
}

func (rs *RedisServer) createKeysWithTTL(ctx context.Context, pipe redis.Pipeliner, keys []string) (map[string]KeyCmds, error) {
	if len(keys) == 0 {
		return make(map[string]KeyCmds), nil
	}

	target := make(map[string]KeyCmds, len(keys))

	for _, key := range keys {
		target[key] = KeyCmds{
			GetCmd: pipe.Get(ctx, key),
			TtlCmd: pipe.TTL(ctx, key),
		}
	}

	_, err := pipe.Exec(ctx)
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, fmt.Errorf("redis pipe exec error: %w", err)
	}

	return target, nil
}
func (rs *RedisServer) getPipeByte(ctx context.Context, keys []string) (*map[string]ResultBytes, error) {
	pipe := rs.redis.Pipeline()
	cmds, err := rs.createKeysWithTTL(ctx, pipe, keys)
	if err != nil {
		return nil, err
	}

	val := map[string]ResultBytes{}

	for key, item := range cmds {
		ttlVal, err := item.TtlCmd.Result()
		if err != nil {
			continue
		}

		valBytes, err := item.GetCmd.Bytes()
		if errors.Is(err, redis.Nil) {
			continue
		} else if err != nil {
			continue
		}

		val[key] = ResultBytes{Ttl: ttlVal, Val: valBytes}
	}

	return &val, nil
}

func (rs *RedisServer) Get(key string) (string, error) {
	if rs == nil || rs.redis == nil {
		return "", errNotInitialized
	}
	ctx, cancel := newDefaultContext()
	defer cancel()
	return rs.redis.Get(ctx, key).Result()
}

func (rs *RedisServer) GetByte(key string) (interface{}, *time.Duration, error) {
	if rs == nil || rs.redis == nil {
		return nil, nil, errNotInitialized
	}

	ctx, cancel := newDefaultContext()
	defer cancel()

	val, err := rs.getPipeByte(ctx, []string{key})

	if err != nil {
		return nil, nil, err
	}
	if val != nil {
		if result, ok := (*val)[key]; ok {
			return result.Val, &result.Ttl, nil
		}
	}

	return nil, nil, nil
}

func (rs *RedisServer) GetByteKeys(keys []string) (*map[string]ResultBytes, error) {
	if len(keys) == 0 {
		val := make(map[string]ResultBytes)
		return &val, nil
	}

	if rs == nil || rs.redis == nil {
		return nil, errNotInitialized
	}
	ctx, cancel := newDefaultContext()
	defer cancel()

	val, err := rs.getPipeByte(ctx, keys)
	if err != nil {
		return nil, err
	}

	return val, nil
}

func (rs *RedisServer) MGet(keys []string) ([]interface{}, error) {
	if rs == nil || rs.redis == nil {
		return nil, errNotInitialized
	}

	ctx, cancel := newDefaultContext()
	defer cancel()

	valsFromSlice, err := rs.redis.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, err
	}

	return valsFromSlice, nil
}

func (rs *RedisServer) Set(key string, val interface{}, expiration time.Duration) error {
	if rs == nil || rs.redis == nil {
		return errNotInitialized
	}
	ctx, cancel := newTimeoutContext(expiration)
	defer cancel()
	return rs.redis.Set(ctx, key, val, expiration).Err()
}

type RedisMSET struct {
	Key        string
	Value      interface{}
	Expiration time.Duration
}

func (rs *RedisServer) MSet(keys []RedisMSET) error {
	if rs == nil || rs.redis == nil {
		return errNotInitialized
	}
	if len(keys) == 0 {
		return nil
	}
	pipe := rs.redis.Pipeline()
	ctx, cancel := newDefaultContext()
	defer cancel()

	for _, item := range keys {
		pipe.Set(ctx, item.Key, item.Value, item.Expiration)
	}

	_, err := pipe.Exec(ctx)
	if err != nil && !errors.Is(err, redis.Nil) {
		return fmt.Errorf("redis pipe exec error: %w", err)
	}

	return nil
}

func (rs *RedisServer) Del(keys ...string) error {
	if rs == nil || rs.redis == nil {
		return errNotInitialized
	}
	ctx, cancel := newDefaultContext()
	defer cancel()
	return rs.redis.Del(ctx, keys...).Err()
}

func (rs *RedisServer) Publish(channel string, message interface{}) (*int64, error) {
	if rs == nil || rs.redis == nil {
		return nil, errNotInitialized
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	receivers, err := rs.redis.Publish(ctx, channel, message).Result()
	if err != nil {
		log.Printf("Publish error: %v", err)
		return nil, err
	}

	return &receivers, nil
}
