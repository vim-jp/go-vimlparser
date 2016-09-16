package main

import (
	"flag"
	"log"
	"os"

	"github.com/haya14busa/go-vimlparser"
	"github.com/k0kubun/pp"
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

	node, err := vimlparser.ParseFile(r, &vimlparser.ParseOption{Neovim: *neovim})
	if err != nil {
		log.Fatal(err)
	}
	pp.Println(node)
}
