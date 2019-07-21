package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/esote/util/dcache"
	"github.com/esote/util/uuid"
)

var cacheUUID *dcache.DCache

func srvUUIDs(w http.ResponseWriter, r *http.Request) {
	count := 10

	if str := r.URL.Query().Get("n"); str != "" {
		if n, err := strconv.Atoi(str); err == nil && n > 0 && n <= 10 {
			count = n
		}
	}

	for i := 0; i < count; i++ {
		_, _ = fmt.Fprintf(w, "%x\n", cacheUUID.NextWg().([]byte))
	}
}

func main() {
	fillUUID := func() interface{} {
		ret, _ := uuid.NewUUID()
		return ret
	}

	var err error

	cacheUUID, err = dcache.NewDCache(1000, fillUUID)

	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", srvUUIDs)

	if err = http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
