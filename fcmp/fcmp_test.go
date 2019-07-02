package fcmp

import (
	"bytes"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"testing"
)

const (
	filex = "test_fx"
	filey = "test_fy"
)

func TestMain(m *testing.M) {
	if _, err := os.Stat(filex); !os.IsNotExist(err) {
		log.Fatalf("test file x '%s' already exists", filex)
	} else if _, err := os.Stat(filey); !os.IsNotExist(err) {
		log.Fatalf("test file y '%s' already exists", filey)
	}

	ret := m.Run()

	if err := os.Remove(filex); err != nil {
		log.Fatal(err)
	} else if err := os.Remove(filey); err != nil {
		log.Fatal(err)
	}

	os.Exit(ret)
}

func TestFuzz(t *testing.T) {
	const reps = 5000

	maxsizes := []int{1, 10, 100, 1000, 10000}

	for _, maxsize := range maxsizes {
		for i := 0; i < reps; i++ {
			strx := make([]byte, rand.Intn(maxsize))
			stry := make([]byte, rand.Intn(maxsize))

			rand.Read(strx)
			rand.Read(stry)

			fx, fy := write_open(strx, stry)

			same, err := FCmp(fx, fy)

			if err != nil {
				t.Error(err)
			} else if bytes.Equal(strx, stry) != same {
				if same {
					t.Error("expected different")
				} else {
					t.Error("expected same")
				}
			}

			fx.Close()
			fy.Close()
		}
	}
}

func TestSame(t *testing.T) {
	str := []byte("abc")

	fx, fy := write_open(str, str)

	if same, err := FCmp(fx, fy); err != nil {
		t.Error(err)
	} else if !same {
		t.Error("expected same")
	}

	fx.Close()
	fy.Close()
}

func TestDifferent(t *testing.T) {
	strx := []byte("ab")
	stry := []byte("abcd")

	fx, fy := write_open(strx, stry)

	if same, err := FCmp(fx, fy); err != nil {
		t.Error(err)
	} else if same {
		t.Error("expected different")
	}

	fx.Close()
	fy.Close()
}

func TestLarge(t *testing.T) {
	size := rand.Intn(10000) + 10000

	strx := make([]byte, size)
	stry := make([]byte, size)

	rand.Read(strx)

	// same
	copy(stry, strx)

	fx, fy := write_open(strx, stry)

	if same, err := FCmp(fx, fy); err != nil {
		t.Error(err)
	} else if !same {
		t.Error("expected same")
	}

	fx.Close()
	fy.Close()

	// different
	rand.Read(stry)

	fx, fy = write_open(strx, stry)

	if same, err := FCmp(fx, fy); err != nil {
		t.Error(err)
	} else if same {
		t.Error("expected different")
	}

	fx.Close()
	fy.Close()
}

func write_open(x []byte, y []byte) (*os.File, *os.File) {
	const perm = 0600

	_ = ioutil.WriteFile(filex, x, perm)
	_ = ioutil.WriteFile(filey, y, perm)

	fx, _ := os.Open(filex)
	fy, _ := os.Open(filey)

	return fx, fy
}
