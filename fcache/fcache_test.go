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
		t.Error(err)
	}

	var b1, b2 []byte

	if b1, err = f.Next(); err != nil {
		t.Error(err)
	}

	time.Sleep(3 * time.Second)

	if b2, err = f.Next(); err != nil {
		t.Error(err)
	}

	if bytes.Equal(b1, b2) {
		t.Error("separate Next's equal")
	}

	if err = f.Clean(); err != nil {
		t.Error(err)
	}

	// calling Next after Clean should still work
	if b1, err = f.Next(); err != nil {
		t.Error(err)
	}

	if bytes.Equal(b1, b2) {
		t.Error("equal after clean")
	}

	f.Clean()
}
