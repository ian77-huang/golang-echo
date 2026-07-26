package store

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"time"

	storeCache "github.com/ian77-huang/golang-echo/pkg/store/cache"
	storeRedis "github.com/ian77-huang/golang-echo/pkg/store/redis"

	goRedis "github.com/redis/go-redis/v9"
)

const REDIS_CANNEL_NAME = "Rbbr_Store_Redis"

var ErrKeyNotFound = errors.New("key not found")

type Keys struct {
	Key    string
	Target interface{}
	Source interface{}
}

func New(so *StoreOption) *StoreServer {
	ss := &StoreServer{}

	if so != nil && so.Cache != nil {
		ss.cache = storeCache.New(so.Cache)
	} else {
		ss.cache = storeCache.New(&CacheOption{})
	}
	if so != nil && so.Redis != nil && so.Redis.IsUse {
		if so.Redis.ConnectURL != "" {
			redisServer, err := storeRedis.NewURL(so.Redis.ConnectURL)
			if err == nil {
				ss.redis = redisServer
			}
		}
		if ss.redis == nil && so.Redis.Option != nil {
			ss.redis = storeRedis.New(so.Redis.Option)
		}
	}

	return ss
}

func normalizeExpiration(expiration time.Duration) (cacheExp time.Duration, redisExp time.Duration) {
	if expiration <= 0 {
		return storeCache.NoExpiration, 0
	}
	return expiration, expiration
}

func isNil(i interface{}) bool {
	if i == nil {
		return true
	}

	v := reflect.ValueOf(i)
	switch v.Kind() {
	case reflect.Ptr, reflect.Map, reflect.Slice, reflect.Chan, reflect.Func, reflect.Interface:
		return v.IsNil()
	}

	return false
}

func (s *StoreServer) GetByteKeys(keys *[]Keys) error {
	if keys == nil || len(*keys) == 0 {
		return nil
	}

	needGetRedis := []string{}
	for i := range *keys {
		item := &(*keys)[i]
		source, found := s.cache.Get(item.Key)
		if !found {
			needGetRedis = append(needGetRedis, item.Key)
		} else {
			item.Source = source
		}
	}

	redisResult, err := s.redis.GetByteKeys(needGetRedis)
	if err != nil {
		return err
	}

	for i := range *keys {
		key := &(*keys)[i]
		if key.Source == nil {
			if item, ok := (*redisResult)[key.Key]; ok {
				key.Source = item.Val
				cacheTTL := time.Duration(0)
				if item.Ttl != 0 {
					cacheTTL = item.Ttl
				}
				s.SetCache(key.Key, key.Source, cacheTTL)
			}
		}
		if isNil(key.Target) {
			return errors.New("store target is not nil")
		} else {
			if !isNil(key.Source) {
				// storeVal := StoreValue{
				// 	Val: key.Target,
				// }
				// err := msgpackDecode(key.Source.([]byte), &storeVal)
				err := msgpackDecode(key.Source.([]byte), key.Target)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (s *StoreServer) GetByte(key string, target interface{}) error {
	source, found := s.cache.Get(key)
	if !found {
		if s.redis != nil {
			valSource, ttl, err := s.redis.GetByte(key)
			if err != nil {
				return err
			}
			source = valSource
			cacheTTL := time.Duration(0)
			if ttl != nil {
				cacheTTL = *ttl
			}
			s.SetCache(key, source, cacheTTL)
		}
	}

	if source == nil {
		return nil
	}
	if isNil(target) {
		return fmt.Errorf("target is not nil")
	}

	// storeVal := StoreValue{
	// 	Val: target,
	// }

	// err := msgpackDecode(source.([]byte), &storeVal)
	err := msgpackDecode(source.([]byte), target)
	if err != nil {
		return err
	}

	return nil
}

func (s *StoreServer) Set(key string, val interface{}, expiration time.Duration) error {
	bVal, err := msgpackEncode(val)
	if err != nil {
		return err
	}

	cacheExp, redisExp := normalizeExpiration(expiration)

	if s.redis != nil {
		if err := s.redis.Set(key, bVal, redisExp); err != nil {
			return err
		}
	}
	s.SetCache(key, bVal, cacheExp)
	return nil
}
func (s *StoreServer) MSet(keys []RedisMSET) error {
	if len(keys) == 0 {
		return nil
	}

	for index := range keys {
		bVal, err := msgpackEncode(keys[index].Value)
		if err != nil {
			return err
		}

		cacheExp, redisExp := normalizeExpiration(keys[index].Expiration)

		keys[index].Expiration = redisExp
		keys[index].Value = bVal

		s.SetCache(keys[index].Key, keys[index].Value, cacheExp)
	}

	if s.redis != nil {
		if err := s.redis.MSet(keys); err != nil {
			return err
		}
	}

	return nil
}

func (s *StoreServer) SetCache(key string, val interface{}, expiration time.Duration) {
	s.cache.Set(key, val, expiration)
}

func (s *StoreServer) Close() error {
	if s.redis != nil {
		return s.redis.Close()
	}
	return nil
}

func (s *StoreServer) Delete(key string) error {
	if s.redis != nil {
		return s.redis.Del(key)
	}
	s.DeleteCache(key)
	return nil
}

func (s *StoreServer) DeleteCache(key string) {
	s.cache.Delete(key)
}

func (s *StoreServer) NewRedisPubSub(res func(msg *goRedis.Message)) context.CancelFunc {
	if s.redis == nil {
		return func() {}
	}
	return s.redis.NewPubSub(REDIS_CANNEL_NAME, res)
}

func (s *StoreServer) PublishRedis(message interface{}) (*int64, error) {
	return s.redis.Publish(REDIS_CANNEL_NAME, message)
}
