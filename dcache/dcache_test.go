package dcache

import (
	"math/rand"
	"testing"
)

// TestNonNil checks that the cache returns filled values, not nil (default
// value of interface{}).
func TestNonNil(t *testing.T) {
	const (
		size = 10
		reps = size*2 + 1
	)

	fill := func() interface{} {
		return rand.Intn(3)
	}

	d, err := NewDCache(size, fill)

	if err != nil {
		t.Error(err)
	}

	for i := 0; i < reps; i++ {
		if d.Next() == nil {
			t.Errorf("%s at index %d", err, i)
		}
	}
}

func BenchmarkNext(b *testing.B) {
	const size = 10

	fill := func() interface{} {
		return 3
	}

	d, _ := NewDCache(size, fill)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = d.Next()
	}
}

// Expected to take more ns/op
func BenchmarkNextWg(b *testing.B) {
	const size = 10

	fill := func() interface{} {
		return 3
	}

	d, _ := NewDCache(size, fill)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = d.NextWg()
	}
}
