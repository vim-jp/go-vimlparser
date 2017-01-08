package langserver

import (
	"bytes"
	"context"
	"fmt"

	"github.com/haya14busa/go-vimlparser"
	"github.com/haya14busa/go-vimlparser/ast"

	"github.com/sourcegraph/jsonrpc2"
)

func NewHandler() jsonrpc2.Handler {
	var langHandler = &LangHandler{
		files: make(map[string]*vimfile),
	}
	return jsonrpc2.HandlerWithError(langHandler.handle)
}

type LangHandler struct {
	files map[string]*vimfile
}

type vimfile struct {
	TextDocumentItem
	Ast      *ast.File
	AstError error
}

func NewVimFile(textDocumentItem TextDocumentItem) (result *vimfile, error error) {
	return &vimfile{
		TextDocumentItem: textDocumentItem,
		Ast:              nil,
		AstError:         nil,
	}, nil
}

func (f *vimfile) GetAst() (result *ast.File, error error) {
	if f.AstError != nil {
		return nil, f.AstError
	} else if f.Ast != nil {
		return f.Ast, nil
	} else {
		opt := &vimlparser.ParseOption{Neovim: false}
		r := bytes.NewBufferString(f.Text)
		ast, err := vimlparser.ParseFile(r, f.URI, opt)
		f.Ast = ast
		f.AstError = err
		return ast, err
	}
}

func (h *LangHandler) handle(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (result interface{}, err error) {
	switch req.Method {
	case "initialize":
		return h.handleInitialize(ctx, conn, req)
	case "textDocument/didOpen":
		return h.handleTextDocumentDidOpen(ctx, conn, req)
	case "textDocument/documentSymbol":
		return h.handleTextDocumentSymbols(ctx, conn, req)
	}

	return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeMethodNotFound, Message: fmt.Sprintf("method not supported: %s", req.Method)}
}
