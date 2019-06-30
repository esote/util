package dcache

import (
	"errors"
	"sync"
)

// DCache (delayed cache) is a cache as a self-populating ring buffer.
// The cache is only repopulated when all values have been withdrawn.
type DCache struct {
	cache []interface{}
	fill  func() interface{}
	index int
	mu    sync.Mutex
	size  int
}

// NewDCache creates a new delayed cache.
func NewDCache(size int, fill func() interface{}) (*DCache, error) {
	if size <= 0 {
		return nil, errors.New("dcache: size <= 0")
	} else if fill == nil {
		return nil, errors.New("dcache: fill is nil")
	}

	return &DCache{
		cache: make([]interface{}, size),
		fill:  fill,
		index: size - 1,
		size:  size,
	}, nil
}

// Next retrieves the next value in the cache. Refilling is done consecutively.
func (d *DCache) Next() interface{} {
	d.mu.Lock()

	if d.index == d.size-1 {
		for i := 0; i < d.size; i++ {
			d.cache[i] = d.fill()
		}
	}

	ret := d.cache[d.index]

	if d.index == 0 {
		d.index = d.size - 1
	} else {
		d.index--
	}

	d.mu.Unlock()

	return ret
}

// NextWg retrieves the next value in the cache. Refilling is done concurrently.
func (d *DCache) NextWg() interface{} {
	d.mu.Lock()

	if d.index == d.size-1 {
		var wg sync.WaitGroup
		wg.Add(d.size)

		for i := 0; i < d.size; i++ {
			go func(i int) {
				d.cache[i] = d.fill()
				wg.Done()
			}(i)
		}

		wg.Wait()
	}

	ret := d.cache[d.index]

	if d.index == 0 {
		d.index = d.size - 1
	} else {
		d.index--
	}

	d.mu.Unlock()

	return ret
}
