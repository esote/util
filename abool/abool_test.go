package abool

import "testing"

func TestBool(t *testing.T) {
	b := New()
	if b.IsSet() {
		t.Fatal("already set")
	}
	if !b.Set() {
		t.Fatal("already set")
	}
	if !b.IsSet() {
		t.Fatal("not set")
	}
	if !b.Unset() {
		t.Fatal("not set")
	}
	if b.IsSet() {
		t.Fatal("already set")
	}
}

func BenchmarkAtomicBool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b := New()
		_ = b.Set()
		_ = b.Unset()
		_ = b.IsSet()
	}
}

func BenchmarkNormalBool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b := false
		b = true
		b = false
		_ = b
	}
}
