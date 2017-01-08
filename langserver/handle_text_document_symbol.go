package langserver

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/haya14busa/go-vimlparser/ast"

	"github.com/sourcegraph/jsonrpc2"
)

func (h *LangHandler) handleTextDocumentSymbols(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (result interface{}, err error) {
	if req.Params == nil {
		return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams}
	}

	var params DocumentSymbolParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		return nil, err
	}

	if f, ok := h.files[params.TextDocument.URI]; ok {
		node, err := f.GetAst()
		if err != nil {
			return nil, err
		}
		return getDocumentSymbols(params, node), nil
	}
	return nil, fmt.Errorf("%s not open", params.TextDocument.URI)
}

func getDocumentSymbols(params DocumentSymbolParams, node ast.Node) []SymbolInformation {
	var symbols []SymbolInformation
	ast.Inspect(node, func(n ast.Node) bool {
		var name string
		var kind SymbolKind
		var pos ast.Pos
		switch x := n.(type) {
		case *ast.Function:
			switch y := x.Name.(type) {
			case *ast.Ident:
				kind = SKFunction
				name = y.Name
				pos = y.NamePos
			}

			if name != "" {
				symbols = append(symbols, SymbolInformation{
					Name: name,
					Kind: kind,
					Location: Location{
						URI: params.TextDocument.URI,
						Range: Range{
							Start: Position{
								Line:      pos.Line - 1,
								Character: pos.Column - 1,
							},
							End: Position{
								Line:      pos.Line - 1,
								Character: pos.Column + len(name) - 1,
							},
						},
					},
				})
				return false
			}
		}
		return true
	})
	return symbols
}
