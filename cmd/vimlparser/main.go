package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/haya14busa/go-vimlparser"
	"github.com/haya14busa/go-vimlparser/compiler"
)

var neovim = flag.Bool("neovim", false, "use neovim parser")

func main() {
	flag.Parse()

	r := os.Stdin

	if p := flag.Arg(0); p != "" {
		f, err := os.Open(p)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			flag.Usage()
			os.Exit(1)
		}
		r = f
	}

	node, err := vimlparser.ParseFile(r, &vimlparser.ParseOption{Neovim: *neovim})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	c := &compiler.Compiler{Config: compiler.Config{Indent: "  "}}
	if err := c.Compile(os.Stdout, node); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
