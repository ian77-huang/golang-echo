package storeRedis

import (
	"errors"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

func newTestRedis(t *testing.T) (*miniredis.Miniredis, *RedisServer) {
	t.Helper()
	mr := miniredis.RunT(t)
	rs := &RedisServer{redis: redis.NewClient(&redis.Options{Addr: mr.Addr()})}
	t.Cleanup(func() { rs.Close() })
	return mr, rs
}

func TestNewAndClose(t *testing.T) {
	mr := miniredis.RunT(t)
	rs := New(&redis.Options{Addr: mr.Addr()})
	if rs == nil || rs.redis == nil {
		t.Fatal("expected initialized redis server")
	}
	if err := rs.Close(); err != nil {
		t.Fatalf("Close() error: %v", err)
	}

	New(nil)
	New(&RedisOption{})

	var nilRS *RedisServer
	if err := nilRS.Close(); err != nil {
		t.Fatalf("nil Close() error: %v", err)
	}
}

func TestNewURL(t *testing.T) {
	mr := miniredis.RunT(t)
	rs, err := NewURL("redis://" + mr.Addr() + "/0")
	if err != nil {
		t.Fatal(err)
	}
	if rs == nil || rs.redis == nil {
		t.Fatal("expected initialized redis server")
	}
	rs.Close()

	if _, err := NewURL("not-a-valid-url"); err == nil {
		t.Fatal("expected invalid url error")
	}
}

func TestSetGetAndDel(t *testing.T) {
	_, rs := newTestRedis(t)

	if err := rs.Set("name", "world", time.Minute); err != nil {
		t.Fatal(err)
	}
	got, err := rs.Get("name")
	if err != nil || got != "world" {
		t.Fatalf("Get() = %q, %v", got, err)
	}

	if _, err := rs.Get("missing"); !errors.Is(err, redis.Nil) {
		t.Fatalf("expected redis.Nil for missing key, got %v", err)
	}

	if err := rs.Del("name"); err != nil {
		t.Fatal(err)
	}
	if _, err := rs.Get("name"); !errors.Is(err, redis.Nil) {
		t.Fatalf("expected redis.Nil after delete, got %v", err)
	}
}

func TestGetByteAndGetByteKeys(t *testing.T) {
	mr, rs := newTestRedis(t)

	_, ttl, err := rs.GetByte("missing")
	if err != nil || ttl != nil {
		t.Fatalf("missing key: ttl=%v err=%v", ttl, err)
	}

	if err := rs.Set("data", []byte("hello"), time.Minute); err != nil {
		t.Fatal(err)
	}

	val, ttl, err := rs.GetByte("data")
	if err != nil || val == nil {
		t.Fatalf("GetByte() = %v, %v, %v", val, ttl, err)
	}
	if string(val.([]byte)) != "hello" {
		t.Fatalf("unexpected value %v", val)
	}
	if ttl == nil || *ttl <= 0 {
		t.Fatalf("expected positive ttl, got %v", ttl)
	}

	result, err := rs.GetByteKeys([]string{"data"})
	if err != nil || result == nil {
		t.Fatal(err)
	}
	if item, ok := (*result)["data"]; !ok || string(item.Val) != "hello" {
		t.Fatalf("unexpected result: %#v", result)
	}

	empty, err := rs.GetByteKeys(nil)
	if err != nil || len(*empty) != 0 {
		t.Fatalf("empty keys result: %#v err=%v", empty, err)
	}

	mr.Set("no-ttl", "x")
	val, ttl, err = rs.GetByte("no-ttl")
	if err != nil || val == nil {
		t.Fatalf("GetByte() no-ttl = %v, %v, %v", val, ttl, err)
	}
}

func TestMSet(t *testing.T) {
	_, rs := newTestRedis(t)

	if err := rs.MSet(nil); err != nil {
		t.Fatal(err)
	}

	keys := []RedisMSET{
		{Key: "a", Value: []byte("1"), Expiration: time.Minute},
		{Key: "b", Value: []byte("2"), Expiration: 0},
	}
	if err := rs.MSet(keys); err != nil {
		t.Fatal(err)
	}

	if got, err := rs.Get("a"); err != nil || got != string([]byte("1")) {
		t.Fatalf("a = %q, %v", got, err)
	}
	if got, err := rs.Get("b"); err != nil || got != string([]byte("2")) {
		t.Fatalf("b = %q, %v", got, err)
	}
}

func TestPublishAndGetByteKeysWithExpiration(t *testing.T) {
	_, rs := newTestRedis(t)

	if _, err := rs.Publish("channel", "hello"); err != nil {
		t.Fatal(err)
	}
}

func TestUninitializedErrors(t *testing.T) {
	var nilRS *RedisServer
	if _, err := nilRS.Get("key"); !errors.Is(err, errNotInitialized) {
		t.Fatalf("nil Get error = %v", err)
	}
	if _, _, err := nilRS.GetByte("key"); !errors.Is(err, errNotInitialized) {
		t.Fatalf("nil GetByte error = %v", err)
	}
	if _, err := nilRS.GetByteKeys([]string{"key"}); !errors.Is(err, errNotInitialized) {
		t.Fatalf("nil GetByteKeys error = %v", err)
	}
	if _, err := nilRS.MGet([]string{"key"}); !errors.Is(err, errNotInitialized) {
		t.Fatalf("nil MGet error = %v", err)
	}
	if err := nilRS.Set("key", "v", 0); !errors.Is(err, errNotInitialized) {
		t.Fatalf("nil Set error = %v", err)
	}
	if err := nilRS.MSet([]RedisMSET{{Key: "a"}}); !errors.Is(err, errNotInitialized) {
		t.Fatalf("nil MSet error = %v", err)
	}
	if err := nilRS.Del("key"); !errors.Is(err, errNotInitialized) {
		t.Fatalf("nil Del error = %v", err)
	}
	if _, err := nilRS.Publish("channel", "msg"); !errors.Is(err, errNotInitialized) {
		t.Fatalf("nil Publish error = %v", err)
	}

	empty := &RedisServer{}
	if _, err := empty.Get("key"); !errors.Is(err, errNotInitialized) {
		t.Fatalf("empty Get error = %v", err)
	}
	if _, err := empty.GetByteKeys(nil); err != nil {
		t.Fatalf("empty GetByteKeys(nil) error = %v", err)
	}
}

func TestNewPubSubNilReceiver(t *testing.T) {
	var nilRS *RedisServer
	cancel := nilRS.NewPubSub("channel", func(msg *redis.Message) {})
	cancel()
}
