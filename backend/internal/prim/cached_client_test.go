package prim

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jossbnd/trainwatch/backend/internal/cache"
)

type mockClient struct {
	calls  int
	visits []StopVisit
	credits int
	err    error
}

func (m *mockClient) FetchStopVisits(_ context.Context, _, _ string) ([]StopVisit, int, error) {
	m.calls++
	return m.visits, m.credits, m.err
}

func newTestCached(inner Client, ttl time.Duration) Client {
	return newCachedWithCache(inner, cache.New[cachedEntry](ttl, ttl), ttl)
}

func TestCachedClient_CacheHit(t *testing.T) {
	mock := &mockClient{visits: []StopVisit{{DirectionRef: TextValue{Value: "A"}}}, credits: 42}
	c := newTestCached(mock, time.Minute)

	v1, cr1, err := c.FetchStopVisits(context.Background(), "stop1", "line1")
	if err != nil {
		t.Fatal(err)
	}
	v2, cr2, err := c.FetchStopVisits(context.Background(), "stop1", "line1")
	if err != nil {
		t.Fatal(err)
	}

	if mock.calls != 1 {
		t.Fatalf("expected 1 upstream call, got %d", mock.calls)
	}
	if len(v1) != 1 || len(v2) != 1 {
		t.Fatal("expected 1 visit in both responses")
	}
	if cr1 != 42 || cr2 != 42 {
		t.Fatalf("expected credits=42, got %d and %d", cr1, cr2)
	}
}

func TestCachedClient_DifferentKeys(t *testing.T) {
	mock := &mockClient{credits: 10}
	c := newTestCached(mock, time.Minute)

	c.FetchStopVisits(context.Background(), "stop1", "line1") //nolint
	c.FetchStopVisits(context.Background(), "stop2", "line1") //nolint

	if mock.calls != 2 {
		t.Fatalf("expected 2 upstream calls for different keys, got %d", mock.calls)
	}
}

func TestCachedClient_ErrorNotCached(t *testing.T) {
	mock := &mockClient{err: errors.New("upstream error")}
	c := newTestCached(mock, time.Minute)

	c.FetchStopVisits(context.Background(), "stop1", "line1") //nolint
	c.FetchStopVisits(context.Background(), "stop1", "line1") //nolint

	if mock.calls != 2 {
		t.Fatalf("expected 2 upstream calls after error, got %d", mock.calls)
	}
}

func TestCachedClient_Expiration(t *testing.T) {
	mock := &mockClient{credits: 5}
	c := newTestCached(mock, 10*time.Millisecond)

	c.FetchStopVisits(context.Background(), "stop1", "line1") //nolint
	time.Sleep(50 * time.Millisecond)
	c.FetchStopVisits(context.Background(), "stop1", "line1") //nolint

	if mock.calls != 2 {
		t.Fatalf("expected 2 upstream calls after expiration, got %d", mock.calls)
	}
}
