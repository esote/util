package uuid

import "testing"

// It is possible, but very unlikely, that this could fail due to a collision.
func TestUniqueUUID(t *testing.T) {
	x, err := NewUUID()

	if err != nil {
		t.Error(err)
	}

	y, err := NewUUID()

	if err != nil {
		t.Error(err)
	}

	if equal(x, y) {
		t.Errorf("x == y")
	}
}

func TestUniqueMegaUUID(t *testing.T) {
	x, err := NewMegaUUID()

	if err != nil {
		t.Error(err)
	}

	y, err := NewMegaUUID()

	if err != nil {
		t.Error(err)
	}

	if equal(x, y) {
		t.Errorf("x == y")
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

func equal(x []byte, y []byte) bool {
	if len(x) != len(y) {
		return false
	}

	for i := 0; i < len(x); i++ {
		if x[i] != y[i] {
			return false
		}
	}

	return true
}
