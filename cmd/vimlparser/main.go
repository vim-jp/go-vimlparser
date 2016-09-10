package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/haya14busa/go-vimlparser"
)

func main() {
	flag.Parse()

	r := os.Stdin

	if p := flag.Arg(0); p != "" {
		f, err := os.Open(p)
		if err != nil {
			log.Fatal(err)
		}
		r = f
	}

	node, err := vimlparser.Parse(r)
	if err != nil {
		log.Fatal(err)
	}
	if err := vimlparser.Compile(os.Stdout, node); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\n")
}
