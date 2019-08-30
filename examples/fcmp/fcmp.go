package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/esote/util/fcmp"
)

func usage() {
	fmt.Fprintln(os.Stderr, "usage: fcmp [-q] FILE1 FILE2")
	os.Exit(1)
}

func main() {
	quiet := flag.Bool("q", false, "quiet")

	flag.Usage = usage

	flag.Parse()

	if len(flag.Args()) < 2 {
		usage()
	}

	equal, err := fcmp.Paths(flag.Args()[0], flag.Args()[1])

	if err != nil {
		log.Fatal(err)
	}

	if !*quiet {
		if equal {
			fmt.Println("equal")
		} else {
			fmt.Println("unequal")
		}
	}

	if equal {
		os.Exit(0)
	}

	os.Exit(1)
}
