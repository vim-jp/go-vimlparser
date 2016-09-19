// Package exporter provides the way to export internal type to public one.
package exporter

import (
	"fmt"

	"github.com/haya14busa/go-vimlparser/ast"
	internal "github.com/haya14busa/go-vimlparser/go"
	"github.com/haya14busa/go-vimlparser/token"
)

// NewNode converts internal node type to ast.Node.
// n.Type must no be zero value.
// n.Pos must no be nil except TOPLEVEL node.
func NewNode(n *internal.ExportNode) ast.Node {
	if n == nil {
		return nil
	}

	// TOPLEVEL doens't have position...?
	var pos ast.Pos
	if p := newPos(n.Pos); p != nil {
		pos = *p
	} else {
		pos = ast.Pos{Offset: 0, Line: 1, Column: 1}
	}

	switch n.Type {

	case internal.NODE_TOPLEVEL:
		return &ast.File{Start: pos, Body: newBody(*n)}

	case internal.NODE_COMMENT:
		return &ast.Comment{
			Quote: pos,
			Text:  n.Str,
		}

	case internal.NODE_EXCMD:
		return &ast.Excmd{
			Excmd:   pos,
			ExArg:   newExArg(*n.Ea),
			Command: n.Str,
		}

	case internal.NODE_FUNCTION:
		attr := ast.FuncAttr{}
		if n.Attr != nil {
			attr = ast.FuncAttr{
				Range:   n.Attr.Range,
				Abort:   n.Attr.Abort,
				Dict:    n.Attr.Dict,
				Closure: n.Attr.Closure,
			}
		}
		return &ast.Function{
			Func:        pos,
			ExArg:       newExArg(*n.Ea),
			Body:        newBody(*n),
			Name:        NewNode(n.Left),
			Params:      newIdents(*n),
			Attr:        attr,
			EndFunction: *NewNode(n.Endfunction).(*ast.EndFunction),
		}

	case internal.NODE_ENDFUNCTION:
		return &ast.EndFunction{
			EndFunc: pos,
			ExArg:   newExArg(*n.Ea),
		}

	case internal.NODE_DELFUNCTION:
		return &ast.DelFunction{
			DelFunc: pos,
			ExArg:   newExArg(*n.Ea),
			Name:    NewNode(n.Left),
		}

	case internal.NODE_RETURN:
		return &ast.Return{
			Return: pos,
			ExArg:  newExArg(*n.Ea),
			Result: NewNode(n.Left),
		}

	case internal.NODE_EXCALL:
		return &ast.ExCall{
			ExCall:   pos,
			ExArg:    newExArg(*n.Ea),
			FuncCall: *NewNode(n.Left).(*ast.CallExpr),
		}

	case internal.NODE_LET:
		return &ast.Let{
			Let:   pos,
			ExArg: newExArg(*n.Ea),
			Op:    n.Op,
			Left:  NewNode(n.Left),
			List:  newList(*n),
			Rest:  NewNode(n.Rest),
			Right: NewNode(n.Right),
		}

	case internal.NODE_UNLET:
		return &ast.UnLet{
			UnLet: pos,
			ExArg: newExArg(*n.Ea),
			List:  newList(*n),
		}

	case internal.NODE_LOCKVAR:
		return &ast.LockVar{
			LockVar: pos,
			ExArg:   newExArg(*n.Ea),
			Depth:   n.Depth,
			List:    newList(*n),
		}

	case internal.NODE_UNLOCKVAR:
		return &ast.UnLockVar{
			UnLockVar: pos,
			ExArg:     newExArg(*n.Ea),
			Depth:     n.Depth,
			List:      newList(*n),
		}

	case internal.NODE_IF:
		var elifs []ast.ElseIf
		if n.Elseif != nil {
			elifs = make([]ast.ElseIf, 0, len(n.Elseif))
		}
		for _, node := range n.Elseif {
			if node != nil { // conservative
				elifs = append(elifs, *NewNode(node).(*ast.ElseIf))
			}
		}
		var els *ast.Else
		if n.Else != nil {
			els = NewNode(n.Else).(*ast.Else)
		}
		return &ast.If{
			If:        pos,
			ExArg:     newExArg(*n.Ea),
			Body:      newBody(*n),
			Condition: NewNode(n.Cond),
			ElseIf:    elifs,
			Else:      els,
			EndIf:     *NewNode(n.Endif).(*ast.EndIf),
		}

	case internal.NODE_ELSEIF:
		return &ast.ElseIf{
			ElseIf:    pos,
			ExArg:     newExArg(*n.Ea),
			Body:      newBody(*n),
			Condition: NewNode(n.Cond),
		}

	case internal.NODE_ELSE:
		return &ast.Else{
			Else:  pos,
			ExArg: newExArg(*n.Ea),
			Body:  newBody(*n),
		}

	case internal.NODE_ENDIF:
		return &ast.EndIf{
			EndIf: pos,
			ExArg: newExArg(*n.Ea),
		}

	case internal.NODE_WHILE:
		return &ast.While{
			While:     pos,
			ExArg:     newExArg(*n.Ea),
			Body:      newBody(*n),
			Condition: NewNode(n.Cond),
			EndWhile:  *NewNode(n.Endwhile).(*ast.EndWhile),
		}

	case internal.NODE_ENDWHILE:
		return &ast.EndWhile{
			EndWhile: pos,
			ExArg:    newExArg(*n.Ea),
		}

	case internal.NODE_FOR:
		return &ast.For{
			For:    pos,
			ExArg:  newExArg(*n.Ea),
			Body:   newBody(*n),
			Left:   NewNode(n.Left),
			List:   newList(*n),
			Rest:   NewNode(n.Rest),
			Right:  NewNode(n.Right),
			EndFor: *NewNode(n.Endfor).(*ast.EndFor),
		}

	case internal.NODE_ENDFOR:
		return &ast.EndFor{
			EndFor: pos,
			ExArg:  newExArg(*n.Ea),
		}

	case internal.NODE_CONTINUE:
		return &ast.Continue{
			Continue: pos,
			ExArg:    newExArg(*n.Ea),
		}

	case internal.NODE_BREAK:
		return &ast.Break{
			Break: pos,
			ExArg: newExArg(*n.Ea),
		}

	case internal.NODE_TRY:
		var catches []ast.Catch
		if n.Catch != nil {
			catches = make([]ast.Catch, 0, len(n.Catch))
		}
		for _, node := range n.Catch {
			if node != nil { // conservative
				catches = append(catches, *NewNode(node).(*ast.Catch))
			}
		}
		var finally *ast.Finally
		if n.Finally != nil {
			finally = NewNode(n.Finally).(*ast.Finally)
		}
		return &ast.Try{
			Try:     pos,
			ExArg:   newExArg(*n.Ea),
			Body:    newBody(*n),
			Catch:   catches,
			Finally: finally,
			EndTry:  *NewNode(n.Endtry).(*ast.EndTry),
		}

	case internal.NODE_CATCH:
		return &ast.Catch{
			Catch:   pos,
			ExArg:   newExArg(*n.Ea),
			Body:    newBody(*n),
			Pattern: n.Pattern,
		}

	case internal.NODE_FINALLY:
		return &ast.Finally{
			Finally: pos,
			ExArg:   newExArg(*n.Ea),
			Body:    newBody(*n),
		}

	case internal.NODE_ENDTRY:
		return &ast.EndTry{
			EndTry: pos,
			ExArg:  newExArg(*n.Ea),
		}

	case internal.NODE_THROW:
		return &ast.Throw{
			Throw: pos,
			ExArg: newExArg(*n.Ea),
			Expr:  NewNode(n.Left),
		}

	case internal.NODE_ECHO, internal.NODE_ECHON, internal.NODE_ECHOMSG, internal.NODE_ECHOERR:
		return &ast.EchoCmd{
			Start:   pos,
			CmdName: n.Ea.Cmd.Name,
			ExArg:   newExArg(*n.Ea),
			Exprs:   newList(*n),
		}

	case internal.NODE_ECHOHL:
		return &ast.Echohl{
			Echohl: pos,
			ExArg:  newExArg(*n.Ea),
			Name:   n.Str,
		}

	case internal.NODE_EXECUTE:
		return &ast.Execute{
			Execute: pos,
			ExArg:   newExArg(*n.Ea),
			Exprs:   newList(*n),
		}

	case internal.NODE_TERNARY:
		return &ast.TernaryExpr{
			Ternary:   pos,
			Condition: NewNode(n.Cond),
			Left:      NewNode(n.Left),
			Right:     NewNode(n.Right),
		}

	case internal.NODE_OR, internal.NODE_AND, internal.NODE_EQUAL, internal.NODE_EQUALCI, internal.NODE_EQUALCS,
		internal.NODE_NEQUAL, internal.NODE_NEQUALCI, internal.NODE_NEQUALCS, internal.NODE_GREATER,
		internal.NODE_GREATERCI, internal.NODE_GREATERCS, internal.NODE_GEQUAL, internal.NODE_GEQUALCI,
		internal.NODE_GEQUALCS, internal.NODE_SMALLER, internal.NODE_SMALLERCI, internal.NODE_SMALLERCS,
		internal.NODE_SEQUAL, internal.NODE_SEQUALCI, internal.NODE_SEQUALCS, internal.NODE_MATCH,
		internal.NODE_MATCHCI, internal.NODE_MATCHCS, internal.NODE_NOMATCH, internal.NODE_NOMATCHCI,
		internal.NODE_NOMATCHCS, internal.NODE_IS, internal.NODE_ISCI, internal.NODE_ISCS, internal.NODE_ISNOT,
		internal.NODE_ISNOTCI, internal.NODE_ISNOTCS, internal.NODE_ADD, internal.NODE_SUBTRACT, internal.NODE_CONCAT,
		internal.NODE_MULTIPLY, internal.NODE_DIVIDE, internal.NODE_REMAINDER:
		return &ast.BinaryExpr{
			Left:  NewNode(n.Left),
			OpPos: pos,
			Op:    opToken(n.Type),
			Right: NewNode(n.Right),
		}

	case internal.NODE_NOT, internal.NODE_MINUS, internal.NODE_PLUS:
		return &ast.UnaryExpr{
			OpPos: pos,
			Op:    opToken(n.Type),
			X:     NewNode(n.Left),
		}

	case internal.NODE_SUBSCRIPT:
		return &ast.SubscriptExpr{
			Lbrack: pos,
			Left:   NewNode(n.Left),
			Right:  NewNode(n.Right),
		}

	case internal.NODE_SLICE:
		return &ast.SliceExpr{
			Lbrack: pos,
			X:      NewNode(n.Left),
			Low:    NewNode(n.Rlist[0]),
			High:   NewNode(n.Rlist[1]),
		}

	case internal.NODE_CALL:
		return &ast.CallExpr{
			Lparen: pos,
			Fun:    NewNode(n.Left),
			Args:   newRlist(*n),
		}

	case internal.NODE_DOT:
		return &ast.DotExpr{
			Left:  NewNode(n.Left),
			Dot:   pos,
			Right: *NewNode(n.Right).(*ast.Ident),
		}

	case internal.NODE_NUMBER:
		return &ast.BasicLit{
			ValuePos: pos,
			Kind:     token.NUMBER,
			Value:    n.Value.(string),
		}
	case internal.NODE_STRING:
		return &ast.BasicLit{
			ValuePos: pos,
			Kind:     token.STRING,
			Value:    n.Value.(string),
		}
	case internal.NODE_LIST:
		return &ast.List{
			Lsquare: pos,
			Values:  newValues(*n),
		}

	case internal.NODE_DICT:
		var kvs []ast.KeyValue
		for _, nn := range n.Value.([]interface{}) {
			kv := nn.([]interface{})
			k := NewNode(internal.NewExportNode(kv[0].(*internal.VimNode)))
			v := NewNode(internal.NewExportNode(kv[1].(*internal.VimNode)))
			kvs = append(kvs, ast.KeyValue{Key: k, Value: v})
		}
		return &ast.Dict{
			Lcurlybrace: pos,
			Entries:     kvs,
		}

	case internal.NODE_OPTION:
		return &ast.BasicLit{
			ValuePos: pos,
			Kind:     token.OPTION,
			Value:    n.Value.(string),
		}
	case internal.NODE_IDENTIFIER:
		return &ast.Ident{
			NamePos: pos,
			Name:    n.Value.(string),
		}

	case internal.NODE_CURLYNAME:
		var parts []ast.CurlyNamePart
		for _, n := range n.Value.([]*internal.VimNode) {
			node := NewNode(internal.NewExportNode(n))
			parts = append(parts, node.(ast.CurlyNamePart))
		}
		return &ast.CurlyName{
			CurlyName: pos,
			Parts:     parts,
		}

	case internal.NODE_ENV:
		return &ast.BasicLit{
			ValuePos: pos,
			Kind:     token.ENV,
			Value:    n.Value.(string),
		}

	case internal.NODE_REG:
		return &ast.BasicLit{
			ValuePos: pos,
			Kind:     token.REG,
			Value:    n.Value.(string),
		}

	case internal.NODE_CURLYNAMEPART:
		return &ast.CurlyNameLit{
			CurlyNameLit: pos,
			Value:        n.Value.(string),
		}

	case internal.NODE_CURLYNAMEEXPR:
		n := n.Value.(*internal.VimNode)
		return &ast.CurlyNameExpr{
			CurlyNameExpr: pos,
			Value:         NewNode(internal.NewExportNode(n)),
		}

	case internal.NODE_LAMBDA:
		return &ast.LambdaExpr{
			Lcurlybrace: pos,
			Params:      newIdents(*n),
			Expr:        NewNode(n.Left),
		}

	}
	panic(fmt.Errorf("Unknown node type: %v, node: %v", n.Type, n))
}

func newBody(n internal.ExportNode) []ast.Statement {
	var body []ast.Statement
	if n.Body != nil {
		body = make([]ast.Statement, 0, len(n.Body))
	}
	for _, node := range n.Body {
		if node != nil { // conservative
			body = append(body, NewNode(node))
		}
	}
	return body
}

func newIdents(n internal.ExportNode) []ast.Ident {
	var idents []ast.Ident
	if n.Rlist != nil {
		idents = make([]ast.Ident, 0, len(n.Rlist))
	}
	for _, node := range n.Rlist {
		if node != nil { // conservative
			idents = append(idents, *NewNode(node).(*ast.Ident))
		}
	}
	return idents
}

func newRlist(n internal.ExportNode) []ast.Expr {
	var exprs []ast.Expr
	if n.Rlist != nil {
		exprs = make([]ast.Expr, 0, len(n.Rlist))
	}
	for _, node := range n.Rlist {
		if node != nil { // conservative
			exprs = append(exprs, NewNode(node))
		}
	}
	return exprs
}

func newList(n internal.ExportNode) []ast.Expr {
	var list []ast.Expr
	if n.List != nil {
		list = make([]ast.Expr, 0, len(n.List))
	}
	for _, node := range n.List {
		if node != nil { // conservative
			list = append(list, NewNode(node))
		}
	}
	return list
}

func newValues(n internal.ExportNode) []ast.Expr {
	var values []ast.Expr
	for _, v := range n.Value.([]interface{}) {
		n := v.(*internal.VimNode)
		values = append(values, NewNode(internal.NewExportNode(n)))
	}
	return values
}
