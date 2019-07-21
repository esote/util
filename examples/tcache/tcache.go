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

type node struct {
	rand []byte
	hash []byte
}

var (
	cache *tcache.TCache
	tree  = []node{{
		rand: []byte{0},
		hash: []byte{0},
	}}
)

func srvNodes(w http.ResponseWriter, r *http.Request) {
	next := cache.Next().(node)

	if _, ok := r.URL.Query()["full"]; ok {
		for _, n := range tree {
			_, _ = fmt.Fprintf(w, "%x : %x\n", n.rand, n.hash)
		}
	} else {
		_, _ = fmt.Fprintf(w, "%x : %x\n", next.rand, next.hash)
	}
}

func main() {
	hash := sha256.New()

	fill := func() interface{} {
		r := make([]byte, sha256.Size)

		if _, err := rand.Read(r); err != nil {
			log.Fatal(err)
		}

		treetop := tree[len(tree)-1]

		hash.Reset()
		hash.Write(append(treetop.rand, treetop.hash...))

		n := node{
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
