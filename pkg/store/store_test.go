package store

import (
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	storeCache "github.com/ian77-huang/golang-echo/pkg/store/cache"
	storeRedis "github.com/ian77-huang/golang-echo/pkg/store/redis"
	goRedis "github.com/redis/go-redis/v9"
)

func newTestStore(t *testing.T) (*miniredis.Miniredis, *StoreServer) {
	t.Helper()
	mr := miniredis.RunT(t)
	ss := &StoreServer{
		cache: storeCache.New(nil),
		redis: storeRedis.New(&goRedis.Options{Addr: mr.Addr()}),
	}
	t.Cleanup(func() { ss.Close() })
	return mr, ss
}

func TestNewWithNilOptions(t *testing.T) {
	ss := New(nil)
	if ss == nil || ss.cache == nil {
		t.Fatal("expected initialized store server with default cache")
	}
	if ss.redis != nil {
		t.Fatal("expected no redis by default")
	}
	ss.Close()

	ss = New(&StoreOption{})
	if ss.cache == nil {
		t.Fatal("expected initialized cache")
	}

	ss = New(&StoreOption{Redis: &StoreOptionRedis{IsUse: false}})
	if ss.redis != nil {
		t.Fatal("expected no redis when IsUse is false")
	}

	ss = New(&StoreOption{Cache: &CacheOption{}})
	if ss.cache == nil {
		t.Fatal("expected initialized cache with explicit option")
	}
}

func TestSetAndGetCacheOnly(t *testing.T) {
	ss := New(nil)

	if err := ss.Set("key", "value", 0); err != nil {
		t.Fatal(err)
	}

	var target string
	if err := ss.GetByte("key", &target); err != nil {
		t.Fatal(err)
	}
	if target != "value" {
		t.Fatalf("target = %q", target)
	}

	if err := ss.GetByte("missing", &target); err != nil {
		t.Fatal(err)
	}

	ss.DeleteCache("key")
	if err := ss.GetByte("key", &target); err != nil {
		t.Fatal(err)
	}
}

func TestSetAndGetWithRedis(t *testing.T) {
	_, ss := newTestStore(t)

	type payload struct {
		Name string
		Age  int
	}
	if err := ss.Set("user:1", payload{Name: "Yien", Age: 30}, time.Minute); err != nil {
		t.Fatal(err)
	}

	var got payload
	if err := ss.GetByte("user:1", &got); err != nil {
		t.Fatal(err)
	}
	if got.Name != "Yien" || got.Age != 30 {
		t.Fatalf("got %#v", got)
	}

	// Miss in cache (fresh server) but hit in redis via a new cache.
	ss2 := &StoreServer{
		cache: storeCache.New(nil),
		redis: ss.redis,
	}
	var got2 payload
	if err := ss2.GetByte("user:1", &got2); err != nil {
		t.Fatal(err)
	}
	if got2.Name != "Yien" {
		t.Fatalf("got2 %#v", got2)
	}
}

func TestDeleteRemovesFromRedis(t *testing.T) {
	mr, ss := newTestStore(t)
	if err := ss.Set("key", "value", time.Minute); err != nil {
		t.Fatal(err)
	}
	if err := ss.Delete("key"); err != nil {
		t.Fatal(err)
	}
	if mr.Exists("key") {
		t.Fatal("expected key to be deleted from redis")
	}
}

func TestGetByteKeys(t *testing.T) {
	mr, ss := newTestStore(t)

	type row struct {
		Num int
	}
	for i, key := range []string{"a", "b"} {
		if err := ss.Set(key, row{Num: i + 1}, time.Minute); err != nil {
			t.Fatal(err)
		}
	}

	// Both keys are in redis; cache is empty so both come from redis.
	keys := []Keys{
		{Key: "a", Target: &row{}},
		{Key: "b", Target: &row{}},
	}
	if err := ss.GetByteKeys(&keys); err != nil {
		t.Fatal(err)
	}
	for _, k := range keys {
		r := k.Target.(*row)
		if r.Num == 0 {
			t.Fatalf("key %q not decoded: %#v", k.Key, r)
		}
	}

	// Second read should come entirely from cache (redis untouched).
	mr.Set("should-not-be-read", "x")
	if err := ss.GetByteKeys(&keys); err != nil {
		t.Fatal(err)
	}

	if err := ss.GetByteKeys(nil); err != nil {
		t.Fatal(err)
	}
	empty := []Keys{}
	if err := ss.GetByteKeys(&empty); err != nil {
		t.Fatal(err)
	}
}

func TestGetByteKeysErrors(t *testing.T) {
	_, ss := newTestStore(t)

	keys := []Keys{{Key: "a", Target: nil}}
	if err := ss.GetByteKeys(&keys); err == nil {
		t.Fatal("expected error for nil target")
	}

	// Missing key from redis leaves target untouched but returns no error.
	keys = []Keys{{Key: "missing", Target: new(int)}}
	if err := ss.GetByteKeys(&keys); err != nil {
		t.Fatalf("missing key should not error: %v", err)
	}
}

func TestMSet(t *testing.T) {
	mr, ss := newTestStore(t)
	if err := ss.MSet([]RedisMSET{
		{Key: "m1", Value: "one", Expiration: time.Minute},
		{Key: "m2", Value: "two", Expiration: 0},
	}); err != nil {
		t.Fatal(err)
	}
	if !mr.Exists("m1") || !mr.Exists("m2") {
		t.Fatal("expected keys in redis")
	}

	ss2 := &StoreServer{cache: storeCache.New(nil), redis: ss.redis}
	var got string
	if err := ss2.GetByte("m1", &got); err != nil || got != "one" {
		t.Fatalf("m1 = %q err=%v", got, err)
	}

	if err := ss.MSet(nil); err != nil {
		t.Fatal(err)
	}
}

func TestStoreValueRoundTrip(t *testing.T) {
	b, err := msgpackEncode(map[string]string{"a": "b"})
	if err != nil {
		t.Fatal(err)
	}
	var target map[string]string
	if err := msgpackDecode(b, &target); err != nil {
		t.Fatal(err)
	}
	if target["a"] != "b" {
		t.Fatalf("target = %#v", target)
	}
}

func TestNormalizeExpiration(t *testing.T) {
	cacheExp, redisExp := normalizeExpiration(0)
	if cacheExp != storeCache.NoExpiration || redisExp != 0 {
		t.Fatalf("expiration = %v, %v", cacheExp, redisExp)
	}
	cacheExp, redisExp = normalizeExpiration(-time.Second)
	if cacheExp != storeCache.NoExpiration || redisExp != 0 {
		t.Fatalf("negative expiration = %v, %v", cacheExp, redisExp)
	}
	cacheExp, redisExp = normalizeExpiration(time.Minute)
	if cacheExp != time.Minute || redisExp != time.Minute {
		t.Fatalf("positive expiration = %v, %v", cacheExp, redisExp)
	}
}

func TestIsNil(t *testing.T) {
	for _, tt := range []struct {
		name string
		val  interface{}
		want bool
	}{
		{"nil", nil, true},
		{"nil pointer", (*int)(nil), true},
		{"nil map", (map[string]int)(nil), true},
		{"nil slice", ([]int)(nil), true},
		{"nil channel", (chan int)(nil), true},
		{"nil func", (func())(nil), true},
		{"empty string", "", false},
		{"int", 1, false},
		{"empty map", map[string]int{}, false},
		{"empty slice", []int{}, false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got := isNil(tt.val); got != tt.want {
				t.Fatalf("isNil(%#v) = %v, want %v", tt.val, got, tt.want)
			}
		})
	}
}

func TestPublishAndPubSub(t *testing.T) {
	_, ss := newTestStore(t)
	if n, err := ss.PublishRedis("hello"); err != nil || n == nil {
		t.Fatalf("publish = %v, %v", n, err)
	}
	cancel := ss.NewRedisPubSub(func(msg *goRedis.Message) {})
	cancel()

	cacheOnly := New(nil)
	if n, err := cacheOnly.PublishRedis("x"); err == nil {
		t.Fatalf("expected error publishing without redis, got %v", n)
	}
	cacheOnly.NewRedisPubSub(func(msg *goRedis.Message) {})()
}
