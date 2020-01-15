package pool

import (
	"sync/atomic"
	"testing"
	"time"
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

	p.Close(false)

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

	p.Close(false)

	if i+dropped != n {
		t.Fatalf("%d+%d != %d", i, dropped, n)
	}
}

func TestCloseFlush(t *testing.T) {
	var i uint64

	f := func(args ...interface{}) {
		time.Sleep(50 * time.Millisecond)
		atomic.AddUint64(&i, 1)
	}

	p := New(1, 1)
	p.Enlist(true, f, 0)
	p.Enlist(true, f, 1)
	p.Close(false)

	if i != 2 {
		t.Fatalf("expected %d got %d", 2, i)
	}

	i = 0
	p = New(1, 1)
	p.Enlist(true, f, 0)
	p.Enlist(true, f, 1)
	p.Close(true)

	// The second, buffered item will now be flushed.
	if i != 1 {
		t.Fatalf("expected %d got %d", 1, i)
	}

}
