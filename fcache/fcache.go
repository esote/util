package fcache

import (
	"errors"
	"io/ioutil"
	"os"
	"sync"
	"time"
)

// FCache (file cache) is a cache which refreshes only after a certain duration.
//
// This is the "file" version of TCache.
type FCache struct {
	dur  time.Duration
	fill func() []byte
	last time.Time
	mu   sync.Mutex
	name string
}

// NewFCache creates a new file cache.
func NewFCache(dur time.Duration, fill func() []byte) (*FCache, error) {
	if dur <= 0 {
		return nil, errors.New("fcache: duration <= 0")
	} else if fill == nil {
		return nil, errors.New("fcache: fill is nil")
	}

	fcache := &FCache{
		dur:  dur,
		fill: fill,
		last: time.Now().UTC().Add(-dur),
	}

	f, err := ioutil.TempFile("", "*.fcache")

	if err != nil {
		return nil, err
	}

	fcache.name = f.Name()

	return fcache, f.Close()
}

// Next retrieves the value in the cache.
func (f *FCache) Next() ([]byte, error) {
	f.mu.Lock()

	if time := time.Now().UTC(); time.Sub(f.last) > f.dur {
		f.last = time

		if err := ioutil.WriteFile(f.name, f.fill(), 0600); err != nil {
			f.mu.Unlock()
			return nil, err
		}
	}

	ret, err := ioutil.ReadFile(f.name)

	f.mu.Unlock()

	return ret, err
}

// Clean removes the cache file. Future calls to Next will simply recreate the
// file.
func (f *FCache) Clean() error {
	f.last = time.Now().UTC().Add(-f.dur)
	return os.Remove(f.name)
}
