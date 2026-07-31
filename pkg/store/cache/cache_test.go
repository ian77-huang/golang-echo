package cache

import (
	"testing"
	"time"
)

func TestNewAppliesDefaults(t *testing.T) {
	cs := New(nil)
	if cs == nil || cs.cache == nil {
		t.Fatal("expected initialized cache server")
	}

	cs = New(&CacheOption{})
	if cs.cache == nil {
		t.Fatal("expected initialized cache server with empty options")
	}

	cs = New(&CacheOption{DefaultExpiration: time.Minute, CleanupInterval: time.Minute})
	if cs.cache == nil {
		t.Fatal("expected initialized cache server with explicit options")
	}
}

func TestSetGetDelete(t *testing.T) {
	cs := New(&CacheOption{DefaultExpiration: time.Minute})

	cs.Set("key", "value", NoExpiration)
	if got, found := cs.Get("key"); !found || got != "value" {
		t.Fatalf("Get() = %v, %v", got, found)
	}

	if _, found := cs.Get("missing"); found {
		t.Fatal("expected missing key")
	}

	cs.Delete("key")
	if _, found := cs.Get("key"); found {
		t.Fatal("expected key to be deleted")
	}
}

func TestExpirationEvictsValue(t *testing.T) {
	cs := New(&CacheOption{})
	cs.Set("temp", "value", time.Millisecond*50)

	if got, found := cs.Get("temp"); !found || got != "value" {
		t.Fatalf("Get() = %v, %v", got, found)
	}

	time.Sleep(time.Millisecond * 100)
	if _, found := cs.Get("temp"); found {
		t.Fatal("expected expired key to be evicted")
	}
}
