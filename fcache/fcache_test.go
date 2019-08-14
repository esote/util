package fcache

import (
	"bytes"
	"crypto/rand"
	"testing"
	"time"
)

func TestSimple(t *testing.T) {
	dur := 3 * time.Second

	fill := func() []byte {
		b := make([]byte, 50)
		rand.Read(b)
		return b
	}

	f, err := NewFCache(dur, fill)

	if err != nil {
		t.Fatal(err)
	}

	var b1, b2 []byte

	if b1, err = f.Next(); err != nil {
		t.Fatal(err)
	}

	time.Sleep(3 * time.Second)

	if b2, err = f.Next(); err != nil {
		t.Fatal(err)
	}

	if bytes.Equal(b1, b2) {
		t.Fatal("separate Next's equal")
	}

	if err = f.Clean(); err != nil {
		t.Fatal(err)
	}

	// calling Next after Clean should still work
	if b1, err = f.Next(); err != nil {
		t.Fatal(err)
	}

	if bytes.Equal(b1, b2) {
		t.Fatal("equal after clean")
	}

	f.Clean()
}
