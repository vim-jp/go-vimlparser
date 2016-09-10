package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/haya14busa/go-vimlparser"
)

var neovim = flag.Bool("neovim", false, "use neovim parser")

func main() {
	flag.Parse()

	r := os.Stdin

	if p := flag.Arg(0); p != "" {
		f, err := os.Open(p)
		if err != nil {
			log.Fatal(err)
			flag.Usage()
		}
		r = f
	}

	node, err := vimlparser.Parse(r, &vimlparser.ParseOption{Neovim: *neovim})
	if err != nil {
		log.Fatal(err)
	}
	if err := vimlparser.Compile(os.Stdout, node); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\n")
}
