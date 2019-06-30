package tcache

import (
	"math/rand"
	"sync"
	"testing"
	"time"
)

// Test that the cache refreshes when it should.
func TestRefresh(t *testing.T) {
	dur := 250 * time.Millisecond

	fill := func() interface{} {
		return rand.Intn(3)
	}

	c, err := NewTCache(dur, fill)

	if err != nil {
		t.Error(err)
	}

	var prev interface{}

	for i := 0; i < 11; i++ {
		n := c.Next()

		if n == nil {
			t.Error("nil returned")
		}

		// It is possible for this to fail due to a timing skew, but
		// that should be rare.
		switch {
		case i == 0 || i == 5:
			prev = n
		case i < 5 || (i > 5 && i < 10):
			if n != prev {
				t.Error("refreshed out of sync")
			}
		}

		time.Sleep(50 * time.Millisecond)
	}
}

// Test for race conditions when calling Next concurrently, run with the -race
// flag.
func TestRace(t *testing.T) {
	dur := 3 * time.Millisecond

	fill := func() interface{} {
		return rand.Intn(3)
	}

	c, _ := NewTCache(dur, fill)

	reps := 1000

	var wg sync.WaitGroup
	wg.Add(reps)
	for i := 0; i < 1000; i++ {
		go func() {
			_ = c.Next()
			time.Sleep(time.Millisecond)
			wg.Done()
		}()
	}
	wg.Wait()
}

func BenchmarkNext(b *testing.B) {
	dur := 5 * time.Millisecond

	fill := func() interface{} {
		return expOp()
	}

	c, _ := NewTCache(dur, fill)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = c.Next()
	}
}

// Expected to take more ns/op.
func BenchmarkCacheless(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = expOp()
	}
}

// Simulate an expensive operation.
func expOp() interface{} {
	time.Sleep(time.Millisecond)
	return rand.Intn(3)
}
