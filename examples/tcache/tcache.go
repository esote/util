package main

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/esote/util/tcache"
)

type Node struct {
	rand []byte
	hash []byte
}

var (
	cache *tcache.TCache
	tree  = []Node{{
		rand: []byte{0},
		hash: []byte{0},
	}}
)

func srvNodes(w http.ResponseWriter, r *http.Request) {
	next := cache.Next().(Node)

	if _, ok := r.URL.Query()["full"]; ok {
		for _, n := range tree {
			fmt.Fprintf(w, "%x : %x\n", n.rand, n.hash)
		}
	} else {
		fmt.Fprintf(w, "%x : %x\n", next.rand, next.hash)
	}
}

func main() {
	hash := sha256.New()

	fill := func() interface{} {
		r := make([]byte, sha256.Size)
		rand.Read(r)

		treetop := tree[len(tree)-1]

		hash.Reset()
		hash.Write(append(treetop.rand, treetop.hash...))

		n := Node{
			rand: r,
			hash: hash.Sum(nil),
		}

		tree = append(tree, n)

		return n
	}

	dur := 5 * time.Second

	var err error

	cache, err = tcache.NewTCache(dur, fill)

	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", srvNodes)

	if err = http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
