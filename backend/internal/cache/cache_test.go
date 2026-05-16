package cache_test

import (
	"testing"
	"time"

	"github.com/jossbnd/trainwatch/backend/internal/cache"
)

func TestGet_Miss(t *testing.T) {
	c := cache.New[string](time.Minute, time.Minute)
	_, ok := c.Get("missing")
	if ok {
		t.Fatal("expected miss, got hit")
	}
}

func TestGet_Hit(t *testing.T) {
	c := cache.New[string](time.Minute, time.Minute)
	c.Set("k", "v", time.Minute)
	v, ok := c.Get("k")
	if !ok {
		t.Fatal("expected hit, got miss")
	}
	if v != "v" {
		t.Fatalf("got %q, want %q", v, "v")
	}
}

func TestSet_Overwrite(t *testing.T) {
	c := cache.New[int](time.Minute, time.Minute)
	c.Set("k", 1, time.Minute)
	c.Set("k", 2, time.Minute)
	v, ok := c.Get("k")
	if !ok {
		t.Fatal("expected hit after overwrite")
	}
	if v != 2 {
		t.Fatalf("got %d, want 2", v)
	}
}

func TestGet_Expiration(t *testing.T) {
	c := cache.New[string](time.Millisecond, time.Millisecond)
	c.Set("k", "v", time.Millisecond)
	time.Sleep(10 * time.Millisecond)
	_, ok := c.Get("k")
	if ok {
		t.Fatal("expected miss after expiration, got hit")
	}
}
