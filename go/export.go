package vimlparser

import (
	"fmt"

	"github.com/haya14busa/go-vimlparser/ast"
	"github.com/haya14busa/go-vimlparser/token"
)

func (self *VimLParser) Parse(reader *StringReader) ast.Node {
	return newAstNode(self.parse(reader))
}

func (self *ExprParser) Parse() ast.Node {
	return newAstNode(self.parse())
}

// ----

// newAstNode converts internal node type to ast.Node.
// n.type_ must no be zero value.
// n.pos must no be nil except TOPLEVEL node.
func newAstNode(n *VimNode) ast.Node {
	if n == nil {
		return nil
	}

	// TOPLEVEL doens't have position...?
	var pos ast.Pos
	if p := newPos(n.pos); p != nil {
		pos = *p
	} else {
		pos = ast.Pos{Offset: 0, Line: 1, Column: 1}
	}

	switch n.type_ {

	case NODE_TOPLEVEL:
		return &ast.File{Start: pos, Body: newBody(*n)}

	case NODE_COMMENT:
		return &ast.Comment{
			Quote: pos,
			Text:  n.str,
		}

	case NODE_EXCMD:
		return &ast.Excmd{
			Excmd:   pos,
			ExArg:   newExArg(*n.ea),
			Command: n.str,
		}

	case NODE_FUNCTION:
		attr := ast.FuncAttr{}
		if n.attr != nil {
			attr = ast.FuncAttr{
				Range:   n.attr.range_,
				Abort:   n.attr.abort,
				Dict:    n.attr.dict,
				Closure: n.attr.closure,
			}
		}
		return &ast.Function{
			Func:        pos,
			ExArg:       newExArg(*n.ea),
			Body:        newBody(*n),
			Name:        newAstNode(n.left),
			Params:      newIdents(*n),
			Attr:        attr,
			EndFunction: *newAstNode(n.endfunction).(*ast.EndFunction),
		}

	case NODE_ENDFUNCTION:
		return &ast.EndFunction{
			EndFunc: pos,
			ExArg:   newExArg(*n.ea),
		}

	case NODE_DELFUNCTION:
		return &ast.DelFunction{
			DelFunc: pos,
			ExArg:   newExArg(*n.ea),
			Name:    newAstNode(n.left),
		}

	case NODE_RETURN:
		return &ast.Return{
			Return: pos,
			ExArg:  newExArg(*n.ea),
			Result: newAstNode(n.left),
		}

	case NODE_EXCALL:
		return &ast.ExCall{
			ExCall:   pos,
			ExArg:    newExArg(*n.ea),
			FuncCall: *newAstNode(n.left).(*ast.CallExpr),
		}

	case NODE_LET:
		return &ast.Let{
			Let:   pos,
			ExArg: newExArg(*n.ea),
			Op:    n.op,
			Left:  newAstNode(n.left),
			List:  newList(*n),
			Rest:  newAstNode(n.rest),
			Right: newAstNode(n.right),
		}

	case NODE_UNLET:
		return &ast.UnLet{
			UnLet: pos,
			ExArg: newExArg(*n.ea),
			List:  newList(*n),
		}

	case NODE_LOCKVAR:
		return &ast.LockVar{
			LockVar: pos,
			ExArg:   newExArg(*n.ea),
			Depth:   n.depth,
			List:    newList(*n),
		}

	case NODE_UNLOCKVAR:
		return &ast.UnLockVar{
			UnLockVar: pos,
			ExArg:     newExArg(*n.ea),
			Depth:     n.depth,
			List:      newList(*n),
		}

	case NODE_IF:
		var elifs []ast.ElseIf
		if n.elseif != nil {
			elifs = make([]ast.ElseIf, 0, len(n.elseif))
		}
		for _, node := range n.elseif {
			if node != nil { // conservative
				elifs = append(elifs, *newAstNode(node).(*ast.ElseIf))
			}
		}
		var els *ast.Else
		if n.else_ != nil {
			els = newAstNode(n.else_).(*ast.Else)
		}
		return &ast.If{
			If:        pos,
			ExArg:     newExArg(*n.ea),
			Body:      newBody(*n),
			Condition: newAstNode(n.cond),
			ElseIf:    elifs,
			Else:      els,
			EndIf:     *newAstNode(n.endif).(*ast.EndIf),
		}

	case NODE_ELSEIF:
		return &ast.ElseIf{
			ElseIf:    pos,
			ExArg:     newExArg(*n.ea),
			Body:      newBody(*n),
			Condition: newAstNode(n.cond),
		}

	case NODE_ELSE:
		return &ast.Else{
			Else:  pos,
			ExArg: newExArg(*n.ea),
			Body:  newBody(*n),
		}

	case NODE_ENDIF:
		return &ast.EndIf{
			EndIf: pos,
			ExArg: newExArg(*n.ea),
		}

	case NODE_WHILE:
		return &ast.While{
			While:     pos,
			ExArg:     newExArg(*n.ea),
			Body:      newBody(*n),
			Condition: newAstNode(n.cond),
			EndWhile:  *newAstNode(n.endwhile).(*ast.EndWhile),
		}

	case NODE_ENDWHILE:
		return &ast.EndWhile{
			EndWhile: pos,
			ExArg:    newExArg(*n.ea),
		}

	case NODE_FOR:
		return &ast.For{
			For:    pos,
			ExArg:  newExArg(*n.ea),
			Body:   newBody(*n),
			Left:   newAstNode(n.left),
			List:   newList(*n),
			Rest:   newAstNode(n.rest),
			Right:  newAstNode(n.right),
			EndFor: *newAstNode(n.endfor).(*ast.EndFor),
		}

	case NODE_ENDFOR:
		return &ast.EndFor{
			EndFor: pos,
			ExArg:  newExArg(*n.ea),
		}

	case NODE_CONTINUE:
		return &ast.Continue{
			Continue: pos,
			ExArg:    newExArg(*n.ea),
		}

	case NODE_BREAK:
		return &ast.Break{
			Break: pos,
			ExArg: newExArg(*n.ea),
		}

	case NODE_TRY:
		var catches []ast.Catch
		if n.catch != nil {
			catches = make([]ast.Catch, 0, len(n.catch))
		}
		for _, node := range n.catch {
			if node != nil { // conservative
				catches = append(catches, *newAstNode(node).(*ast.Catch))
			}
		}
		var finally *ast.Finally
		if n.finally != nil {
			finally = newAstNode(n.finally).(*ast.Finally)
		}
		return &ast.Try{
			Try:     pos,
			ExArg:   newExArg(*n.ea),
			Body:    newBody(*n),
			Catch:   catches,
			Finally: finally,
			EndTry:  *newAstNode(n.endtry).(*ast.EndTry),
		}

	case NODE_CATCH:
		return &ast.Catch{
			Catch:   pos,
			ExArg:   newExArg(*n.ea),
			Body:    newBody(*n),
			Pattern: n.pattern,
		}

	case NODE_FINALLY:
		return &ast.Finally{
			Finally: pos,
			ExArg:   newExArg(*n.ea),
			Body:    newBody(*n),
		}

	case NODE_ENDTRY:
		return &ast.EndTry{
			EndTry: pos,
			ExArg:  newExArg(*n.ea),
		}

	case NODE_THROW:
		return &ast.Throw{
			Throw: pos,
			ExArg: newExArg(*n.ea),
			Expr:  newAstNode(n.left),
		}

	case NODE_ECHO, NODE_ECHON, NODE_ECHOMSG, NODE_ECHOERR:
		return &ast.EchoCmd{
			Start:   pos,
			CmdName: n.ea.cmd.name,
			ExArg:   newExArg(*n.ea),
			Exprs:   newList(*n),
		}

	case NODE_ECHOHL:
		return &ast.Echohl{
			Echohl: pos,
			ExArg:  newExArg(*n.ea),
			Name:   n.str,
		}

	case NODE_EXECUTE:
		return &ast.Execute{
			Execute: pos,
			ExArg:   newExArg(*n.ea),
			Exprs:   newList(*n),
		}

	case NODE_TERNARY:
		return &ast.TernaryExpr{
			Ternary:   pos,
			Condition: newAstNode(n.cond),
			Left:      newAstNode(n.left),
			Right:     newAstNode(n.right),
		}

	case NODE_OR, NODE_AND, NODE_EQUAL, NODE_EQUALCI, NODE_EQUALCS,
		NODE_NEQUAL, NODE_NEQUALCI, NODE_NEQUALCS, NODE_GREATER,
		NODE_GREATERCI, NODE_GREATERCS, NODE_GEQUAL, NODE_GEQUALCI,
		NODE_GEQUALCS, NODE_SMALLER, NODE_SMALLERCI, NODE_SMALLERCS,
		NODE_SEQUAL, NODE_SEQUALCI, NODE_SEQUALCS, NODE_MATCH,
		NODE_MATCHCI, NODE_MATCHCS, NODE_NOMATCH, NODE_NOMATCHCI,
		NODE_NOMATCHCS, NODE_IS, NODE_ISCI, NODE_ISCS, NODE_ISNOT,
		NODE_ISNOTCI, NODE_ISNOTCS, NODE_ADD, NODE_SUBTRACT, NODE_CONCAT,
		NODE_MULTIPLY, NODE_DIVIDE, NODE_REMAINDER:
		return &ast.BinaryExpr{
			Left:  newAstNode(n.left),
			OpPos: pos,
			Op:    opToken(n.type_),
			Right: newAstNode(n.right),
		}

	case NODE_NOT, NODE_MINUS, NODE_PLUS:
		return &ast.UnaryExpr{
			OpPos: pos,
			Op:    opToken(n.type_),
			X:     newAstNode(n.left),
		}

	case NODE_SUBSCRIPT:
		return &ast.SubscriptExpr{
			Lbrack: pos,
			Left:   newAstNode(n.left),
			Right:  newAstNode(n.right),
		}

	case NODE_SLICE:
		return &ast.SliceExpr{
			Lbrack: pos,
			X:      newAstNode(n.left),
			Low:    newAstNode(n.rlist[0]),
			High:   newAstNode(n.rlist[1]),
		}

	case NODE_CALL:
		return &ast.CallExpr{
			Lparen: pos,
			Fun:    newAstNode(n.left),
			Args:   newRlist(*n),
		}

	case NODE_DOT:
		return &ast.DotExpr{
			Left:  newAstNode(n.left),
			Dot:   pos,
			Right: *newAstNode(n.right).(*ast.Ident),
		}

	case NODE_NUMBER:
		return &ast.BasicLit{
			ValuePos: pos,
			Kind:     token.NUMBER,
			Value:    n.value.(string),
		}
	case NODE_STRING:
		return &ast.BasicLit{
			ValuePos: pos,
			Kind:     token.STRING,
			Value:    n.value.(string),
		}
	case NODE_LIST:
		return &ast.List{
			Lsquare: pos,
			Values:  newValues(*n),
		}

	case NODE_DICT:
		var kvs []ast.KeyValue
		for _, nn := range n.value.([]interface{}) {
			kv := nn.([]interface{})
			k := newAstNode(kv[0].(*VimNode))
			v := newAstNode(kv[1].(*VimNode))
			kvs = append(kvs, ast.KeyValue{Key: k, Value: v})
		}
		return &ast.Dict{
			Lcurlybrace: pos,
			Entries:     kvs,
		}

	case NODE_OPTION:
		return &ast.BasicLit{
			ValuePos: pos,
			Kind:     token.OPTION,
			Value:    n.value.(string),
		}
	case NODE_IDENTIFIER:
		return &ast.Ident{
			NamePos: pos,
			Name:    n.value.(string),
		}

	case NODE_CURLYNAME:
		var parts []ast.CurlyNamePart
		for _, n := range n.value.([]*VimNode) {
			node := newAstNode(n)
			parts = append(parts, node.(ast.CurlyNamePart))
		}
		return &ast.CurlyName{
			CurlyName: pos,
			Parts:     parts,
		}

	case NODE_ENV:
		return &ast.BasicLit{
			ValuePos: pos,
			Kind:     token.ENV,
			Value:    n.value.(string),
		}

	case NODE_REG:
		return &ast.BasicLit{
			ValuePos: pos,
			Kind:     token.REG,
			Value:    n.value.(string),
		}

	case NODE_CURLYNAMEPART:
		return &ast.CurlyNameLit{
			CurlyNameLit: pos,
			Value:        n.value.(string),
		}

	case NODE_CURLYNAMEEXPR:
		n := n.value.(*VimNode)
		return &ast.CurlyNameExpr{
			CurlyNameExpr: pos,
			Value:         newAstNode(n),
		}

	case NODE_LAMBDA:
		return &ast.LambdaExpr{
			Lcurlybrace: pos,
			Params:      newIdents(*n),
			Expr:        newAstNode(n.left),
		}

	}
	panic(fmt.Errorf("Unknown node type: %v, node: %v", n.type_, n))
}

func newPos(p *pos) *ast.Pos {
	if p == nil {
		return nil
	}
	return &ast.Pos{
		Offset: p.i,
		Line:   p.lnum,
		Column: p.col,
	}
}

func newExArg(ea ExArg) ast.ExArg {
	return ast.ExArg{
		Forceit:    ea.forceit,
		AddrCount:  ea.addr_count,
		Line1:      ea.line1,
		Line2:      ea.line2,
		Flags:      ea.flags,
		DoEcmdCmd:  ea.do_ecmd_cmd,
		DoEcmdLnum: ea.do_ecmd_lnum,
		Append:     ea.append,
		Usefilter:  ea.usefilter,
		Amount:     ea.amount,
		Regname:    ea.regname,
		ForceBin:   ea.force_bin,
		ReadEdit:   ea.read_edit,
		ForceFf:    ea.force_ff,
		ForceEnc:   ea.force_enc,
		BadChar:    ea.bad_char,
		Linepos:    newPos(ea.linepos),
		Cmdpos:     newPos(ea.cmdpos),
		Argpos:     newPos(ea.argpos),
		Cmd:        newCmd(ea.cmd),
		Modifiers:  ea.modifiers,
		Range:      ea.range_,
		Argopt:     ea.argopt,
		Argcmd:     ea.argcmd,
	}
}

func newCmd(c *Cmd) *ast.Cmd {
	if c == nil {
		return nil
	}
	return &ast.Cmd{
		Name:   c.name,
		Minlen: c.minlen,
		Flags:  c.flags,
		Parser: c.parser,
	}
}

func newBody(n VimNode) []ast.Statement {
	var body []ast.Statement
	if n.body != nil {
		body = make([]ast.Statement, 0, len(n.body))
	}
	for _, node := range n.body {
		if node != nil { // conservative
			body = append(body, newAstNode(node))
		}
	}
	return body
}

func newIdents(n VimNode) []ast.Ident {
	var idents []ast.Ident
	if n.rlist != nil {
		idents = make([]ast.Ident, 0, len(n.rlist))
	}
	for _, node := range n.rlist {
		if node != nil { // conservative
			idents = append(idents, *newAstNode(node).(*ast.Ident))
		}
	}
	return idents
}

func newRlist(n VimNode) []ast.Expr {
	var exprs []ast.Expr
	if n.rlist != nil {
		exprs = make([]ast.Expr, 0, len(n.rlist))
	}
	for _, node := range n.rlist {
		if node != nil { // conservative
			exprs = append(exprs, newAstNode(node))
		}
	}
	return exprs
}

func newList(n VimNode) []ast.Expr {
	var list []ast.Expr
	if n.list != nil {
		list = make([]ast.Expr, 0, len(n.list))
	}
	for _, node := range n.list {
		if node != nil { // conservative
			list = append(list, newAstNode(node))
		}
	}
	return list
}

func newValues(n VimNode) []ast.Expr {
	var values []ast.Expr
	for _, v := range n.value.([]interface{}) {
		n := v.(*VimNode)
		values = append(values, newAstNode(n))
	}
	return values
}

func opToken(nodeType int) token.Token {
	switch nodeType {
	case NODE_OR:
		return token.OROR
	case NODE_AND:
		return token.ANDAND
	case NODE_EQUAL:
		return token.EQEQ
	case NODE_EQUALCI:
		return token.EQEQCI
	case NODE_EQUALCS:
		return token.EQEQCS
	case NODE_NEQUAL:
		return token.NEQ
	case NODE_NEQUALCI:
		return token.NEQCI
	case NODE_NEQUALCS:
		return token.NEQCS
	case NODE_GREATER:
		return token.GT
	case NODE_GREATERCI:
		return token.GTCI
	case NODE_GREATERCS:
		return token.GTCS
	case NODE_GEQUAL:
		return token.GTEQ
	case NODE_GEQUALCI:
		return token.GTEQCI
	case NODE_GEQUALCS:
		return token.GTEQCS
	case NODE_SMALLER:
		return token.LT
	case NODE_SMALLERCI:
		return token.LTCI
	case NODE_SMALLERCS:
		return token.LTCS
	case NODE_SEQUAL:
		return token.LTEQ
	case NODE_SEQUALCI:
		return token.LTEQCI
	case NODE_SEQUALCS:
		return token.LTEQCS
	case NODE_MATCH:
		return token.MATCH
	case NODE_MATCHCI:
		return token.MATCHCI
	case NODE_MATCHCS:
		return token.MATCHCS
	case NODE_NOMATCH:
		return token.NOMATCH
	case NODE_NOMATCHCI:
		return token.NOMATCHCI
	case NODE_NOMATCHCS:
		return token.NOMATCHCS
	case NODE_IS:
		return token.IS
	case NODE_ISCI:
		return token.ISCI
	case NODE_ISCS:
		return token.ISCS
	case NODE_ISNOT:
		return token.ISNOT
	case NODE_ISNOTCI:
		return token.ISNOTCI
	case NODE_ISNOTCS:
		return token.ISNOTCS
	case NODE_ADD:
		return token.PLUS
	case NODE_SUBTRACT:
		return token.MINUS
	case NODE_CONCAT:
		return token.DOT
	case NODE_MULTIPLY:
		return token.STAR
	case NODE_DIVIDE:
		return token.SLASH
	case NODE_REMAINDER:
		return token.PERCENT
	case NODE_NOT:
		return token.NOT
	case NODE_MINUS:
		return token.MINUS
	case NODE_PLUS:
		return token.PLUS
	}
	return token.ILLEGAL
}
