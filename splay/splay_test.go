package splay

import (
	"bytes"
	"testing"
)

func TestSplay(t *testing.T) {
	s, err := NewSplay("testdata", 2)

	if err != nil {
		t.Fatal(err)
	}

	name := "star-trek"
	data := []byte("I doubt there are many oak trees on Tagus.\n")

	if err = s.Write(name, data); err != nil {
		t.Fatal(err)
	}

	if out, err := s.Read(name); err != nil {
		t.Fatal(err)
	} else if !bytes.Equal(out, data) {
		t.Fatal("bytes read mismatch")
	}

	if err = s.Remove(name); err != nil {
		t.Fatal(err)
	}

	_ = s.Write(name, data)

	if err = s.RemoveAll(); err != nil {
		t.Fatal(err)
	}
}

func TestNewSplay(t *testing.T) {
	if _, err := NewSplay("", 2); err == nil {
		t.Fatal("expected empty name error")
	}

	if _, err := NewSplay("testdata", 0); err == nil {
		t.Fatal("expected cutoff zero error")
	}
}

func TestSplayName(t *testing.T) {
	s, _ := NewSplay("testdata", 2)

	if err := s.Write("st", []byte("abc")); err == nil {
		t.Fatal("expected name error")
	}

	_ = s.RemoveAll()
}
