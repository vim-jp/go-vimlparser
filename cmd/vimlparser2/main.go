package main

import (
	"flag"
	"log"
	"os"

	"github.com/haya14busa/go-vimlparser"
	"github.com/k0kubun/pp"
)

func main() {
	r := os.Stdin

	if p := flag.Arg(0); p != "" {
		f, err := os.Open(p)
		if err != nil {
			log.Fatal(err)
			flag.Usage()
		}
		r = f
	}

	node, err := vimlparser.ParseFile(r, &vimlparser.ParseOption{Neovim: false})
	if err != nil {
		log.Fatal(err)
	}
	pp.Println(node)
}
