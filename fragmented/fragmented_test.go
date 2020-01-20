package fragmented

import (
	"bytes"
	"io"
	"testing"
)

type bufnext struct {
	bytes.Buffer
}

func (b *bufnext) Next() (io.ReadWriteCloser, error) {
	return b, nil
}

func (b *bufnext) Close() error {
	return nil
}

func TestFragment(t *testing.T) {
	const s = "A test string of moderate length"

	var b bufnext
	f := New(&b, 3)
	defer func() {
		if err := f.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	if n, err := f.Write([]byte(s)); err != nil {
		t.Fatal(err)
	} else if n != len(s) {
		t.Fatal("n != len(s)")
	}

	p := make([]byte, len(s))

	if n, err := f.Read(p); err != nil {
		t.Fatal(err)
	} else if n != len(s) {
		t.Fatal("n != len(s)")
	} else if string(p) != s {
		t.Fatal("p != s")
	}
}
