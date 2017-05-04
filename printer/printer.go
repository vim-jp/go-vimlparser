// Package printer implements printing of AST nodes.
//
// This is WIP package. DO NOT USE.
package printer

import (
	"errors"
	"fmt"
	"io"

	"github.com/haya14busa/go-vimlparser/ast"
)

// A Config node controls the output of Fprint.
type Config struct{}

// Fprint "pretty-prints" an AST node to output for a given configuration cfg.
func Fprint(output io.Writer, node ast.Node, cfg *Config) error {
	var p printer
	p.init(cfg)
	if err := p.printNode(node); err != nil {
		return err
	}
	if _, err := output.Write(p.output); err != nil {
		return err
	}
	return nil
}

type printer struct {
	Config

	// Current state
	output []byte // raw printer result
}

func (p *printer) init(cfg *Config) {
	if cfg == nil {
		cfg = &Config{}
	}
	p.Config = *cfg
}

func (p *printer) writeString(s string) {
	p.output = append(p.output, s...)
}

func (p *printer) printNode(node ast.Node) error {
	switch n := node.(type) {
	case *ast.File:
		return p.file(n)
	case ast.Expr:
		return p.expr(n)
	case ast.Statement:
		return p.stmt(n)
	default:
		return fmt.Errorf("go-vimlparser/printer: unsupported node type %T", node)
	}
}

func (p *printer) file(f *ast.File) error {
	return errors.New("Not implemented: printer.file")
}

func (p *printer) expr(expr ast.Expr) error {
	switch n := expr.(type) {
	// case *ast.TernaryExpr:
	// case *ast.BinaryExpr:
	// case *ast.UnaryExpr:
	// case *ast.SubscriptExpr:
	// case *ast.SliceExpr:
	// case *ast.CallExpr:
	// case *ast.DotExpr:
	// case *ast.List:
	// case *ast.Dict:
	// case *ast.CurlyName:
	// case *ast.CurlyNameLit:
	// case *ast.CurlyNameExpr:
	case *ast.BasicLit:
		p.writeString(n.Value)
	case *ast.Ident:
		p.writeString(n.Name)
	// case *ast.LambdaExpr:
	// case *ast.ParenExpr:
	default:
		return fmt.Errorf("unsupported expr type %T", n)
	}
	return nil
}

func (p *printer) stmt(node ast.Statement) error {
	return errors.New("Not implemented: printer.stmt")
}
