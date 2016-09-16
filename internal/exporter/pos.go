package exporter

import "github.com/haya14busa/go-vimlparser/ast"
import internal "github.com/haya14busa/go-vimlparser/go"

func newPos(p *internal.ExportPos) *ast.Pos {
	if p == nil {
		return nil
	}
	return &ast.Pos{
		Offset: p.I,
		Line:   p.Lnum,
		Column: p.Col,
	}
}
