package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/haya14busa/go-vimlparser"
	"github.com/haya14busa/go-vimlparser/compiler"
)

var neovim = flag.Bool("neovim", false, "use neovim parser")

func main() {
	flag.Parse()

	opt := &vimlparser.ParseOption{Neovim: *neovim}

	if len(flag.Args()) == 0 {
		if err := parseFile("", os.Stdin, os.Stdout, opt); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	exitCode := 0

	for i, file := range flag.Args() {
		f, err := os.Open(file)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			if i == 0 {
				flag.Usage()
			}
			exitCode = 1
		}
		if err := parseFile(f.Name(), f, os.Stdout, opt); err != nil {
			fmt.Fprintln(os.Stderr, err)
			exitCode = 1
		}
	}

	os.Exit(exitCode)
}

// filename is empty string if r is os.Stdin
func parseFile(filename string, r io.ReadCloser, w io.Writer, opt *vimlparser.ParseOption) error {
	defer r.Close()
	node, err := vimlparser.ParseFile(r, filename, opt)
	if err != nil {
		return err
	}
	c := &compiler.Compiler{Config: compiler.Config{Indent: "  "}}
	if err := c.Compile(w, node); err != nil {
		return err
	}
	return nil
}
