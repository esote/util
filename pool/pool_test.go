package pool

import (
	"sync/atomic"
	"testing"
)

const n = 50000

func TestPool(t *testing.T) {
	var i uint64
	f := func(args ...interface{}) {
		atomic.AddUint64(&i, 1)
	}
	p := New(10, 100)

	for i := 0; i < n; i++ {
		p.Enlist(true, f, i)
	}

	p.Close()

	if i != n {
		t.Fatalf("%d != %d", i, n)
	}
}

func TestPoolNonblocking(t *testing.T) {
	var i uint64
	f := func(args ...interface{}) {
		atomic.AddUint64(&i, 1)
	}

	p := New(10, 100)

	var dropped uint64
	for i := 0; i < n; i++ {
		if !p.Enlist(false, f, i) {
			dropped++
		}
	}

	p.Close()

	if i+dropped != n {
		t.Fatalf("%d+%d != %d", i, dropped, n)
	}
}
