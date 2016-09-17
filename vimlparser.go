package vimlparser

import (
	"bufio"
	"fmt"
	"io"

	"github.com/haya14busa/go-vimlparser/ast"
	internal "github.com/haya14busa/go-vimlparser/go"
	"github.com/haya14busa/go-vimlparser/internal/exporter"
)

// ParseOption is option for Parse().
type ParseOption struct {
	Neovim bool
}

// ParseFile parses Vim script.
// filename can be empty.
func ParseFile(r io.Reader, filename string, opt *ParseOption) (node *ast.File, err error) {
	defer func() {
		if r := recover(); r != nil {
			node = nil
			err = fmt.Errorf("%v", r)
			// log.Printf("%s", debug.Stack())
		}
	}()
	lines := readlines(r)
	reader := internal.NewStringReader(lines)
	neovim := false
	if opt != nil {
		neovim = opt.Neovim
	}
	node = exporter.NewNode(internal.NewVimLParser(neovim).Parse(reader)).(*ast.File)
	return
}

// ParseExpr parses Vim script expression.
func ParseExpr(r io.Reader) (node ast.Expr, err error) {
	defer func() {
		if r := recover(); r != nil {
			node = nil
			err = fmt.Errorf("%v", r)
			// log.Printf("%s", debug.Stack())
		}
	}()
	lines := readlines(r)
	reader := internal.NewStringReader(lines)
	p := internal.NewExprParser(reader)
	node = exporter.NewNode(p.Parse())
	return
}

func readlines(r io.Reader) []string {
	lines := []string{}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}
