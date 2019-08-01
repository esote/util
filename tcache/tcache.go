package tcache

import (
	"errors"
	"sync"
	"time"
)

// TCache (timed cache) is a cache which refreshes only after a certain
// duration.
//
// This is the "memory" version of FCache.
type TCache struct {
	cache interface{}
	dur   time.Duration
	fill  func() interface{}
	last  time.Time
	mu    sync.Mutex
}

// NewTCache creates a new timed cache.
func NewTCache(dur time.Duration, fill func() interface{}) (*TCache, error) {
	if dur <= 0 {
		return nil, errors.New("tcache: duration <= 0")
	} else if fill == nil {
		return nil, errors.New("tcache: fill is nil")
	}

	return &TCache{
		cache: nil,
		dur:   dur,
		fill:  fill,
		last:  time.Now().UTC().Add(-dur),
	}, nil
}

// Next retrieves the value in the cache.
func (t *TCache) Next() interface{} {
	t.mu.Lock()

	if time := time.Now().UTC(); time.Sub(t.last) > t.dur {
		t.last = time
		t.cache = t.fill()
	}

	ret := t.cache

	t.mu.Unlock()

	return ret
}
