package cache

import (
	"time"

	gocache "github.com/patrickmn/go-cache"
)

// Cache is a generic key-value store with per-entry TTL.
type Cache[V any] interface {
	Get(key string) (V, bool)
	Set(key string, value V, ttl time.Duration)
}

type inMemory[V any] struct {
	c *gocache.Cache
}

// New returns an in-memory Cache backed by go-cache.
// cleanupInterval controls how often expired entries are evicted.
func New[V any](defaultTTL, cleanupInterval time.Duration) Cache[V] {
	return &inMemory[V]{c: gocache.New(defaultTTL, cleanupInterval)}
}

func (m *inMemory[V]) Get(key string) (V, bool) {
	v, ok := m.c.Get(key)
	if !ok {
		var zero V
		return zero, false
	}
	return v.(V), true
}

func (m *inMemory[V]) Set(key string, value V, ttl time.Duration) {
	m.c.Set(key, value, ttl)
}
