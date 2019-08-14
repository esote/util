package uuid

import (
	"bytes"
	"testing"
)

// It is possible, but very unlikely, that this could fail due to a collision.
func TestUniqueUUID(t *testing.T) {
	x, err := NewUUID()

	if err != nil {
		t.Fatal(err)
	}

	y, err := NewUUID()

	if err != nil {
		t.Fatal(err)
	}

	if bytes.Equal(x, y) {
		t.Fatal("x == y")
	}
}

func TestUniqueMegaUUID(t *testing.T) {
	x, err := NewMegaUUID()

	if err != nil {
		t.Fatal(err)
	}

	y, err := NewMegaUUID()

	if err != nil {
		t.Fatal(err)
	}

	if bytes.Equal(x, y) {
		t.Fatal("x == y")
	}
}

func BenchmarkUUID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = NewUUID()
	}
}

func BenchmarkMegaUUID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = NewMegaUUID()
	}
}
