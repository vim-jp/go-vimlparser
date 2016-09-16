// Package compiler provides compiler from Vim script AST into S-expression
// like format which is the same format as Compiler of vim-vimlparser.
// ref: "go/printer"
package compiler

import "github.com/haya14busa/go-vimlparser/ast"
import "io"

type Config struct {
	Indent string
}

type Compiler struct {
	Config

	// Current state
	output []byte // raw compiler result
	indent int    // current indentation
}

func (c *Compiler) Compile(w io.Writer, node ast.Node) error {
	return c.compile(w, node)
}

func (c *Compiler) compile(w io.Writer, node ast.Node) error {
	return nil
}
