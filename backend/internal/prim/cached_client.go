package prim

import (
	"context"
	"time"

	"github.com/jossbnd/trainwatch/backend/internal/cache"
)

type cachedEntry struct {
	visits  []StopVisit
	credits int
}

type cachedClient struct {
	inner Client
	c     cache.Cache[cachedEntry]
	ttl   time.Duration
}

// NewCached wraps inner with an in-memory cache keyed by (stopRef, lineRef).
// Errors are never cached. On a cache hit, the last known credit count is returned.
func NewCached(inner Client, ttl time.Duration) Client {
	return newCachedWithCache(inner, cache.New[cachedEntry](ttl, 5*time.Minute), ttl)
}

// newCachedWithCache allows injecting a custom cache and TTL in tests.
func newCachedWithCache(inner Client, c cache.Cache[cachedEntry], ttl time.Duration) Client {
	return &cachedClient{inner: inner, c: c, ttl: ttl}
}

func (cc *cachedClient) FetchStopVisits(ctx context.Context, stopRef, lineRef string) ([]StopVisit, int, error) {
	key := stopRef + "|" + lineRef

	if entry, ok := cc.c.Get(key); ok {
		return entry.visits, entry.credits, nil
	}

	visits, credits, err := cc.inner.FetchStopVisits(ctx, stopRef, lineRef)
	if err != nil {
		return nil, credits, err
	}

	cc.c.Set(key, cachedEntry{visits: visits, credits: credits}, cc.ttl)
	return visits, credits, nil
}
