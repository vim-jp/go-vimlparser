package exporter

import "github.com/haya14busa/go-vimlparser/ast"
import internal "github.com/haya14busa/go-vimlparser/go"

func newExArg(ea internal.ExportExArg) ast.ExArg {
	return ast.ExArg{
		Forceit:    ea.Forceit,
		AddrCount:  ea.AddrCount,
		Line1:      ea.Line1,
		Line2:      ea.Line2,
		Flags:      ea.Flags,
		DoEcmdCmd:  ea.DoEcmdCmd,
		DoEcmdLnum: ea.DoEcmdLnum,
		Append:     ea.Append,
		Usefilter:  ea.Usefilter,
		Amount:     ea.Amount,
		Regname:    ea.Regname,
		ForceBin:   ea.ForceBin,
		ReadEdit:   ea.ReadEdit,
		ForceFf:    ea.ForceFf,
		ForceEnc:   ea.ForceEnc,
		BadChar:    ea.BadChar,
		Linepos:    newPos(ea.Linepos),
		Cmdpos:     newPos(ea.Cmdpos),
		Argpos:     newPos(ea.Argpos),
		Cmd:        newCmd(ea.Cmd),
		Modifiers:  ea.Modifiers,
		Range:      ea.Range,
		Argopt:     ea.Argopt,
		Argcmd:     ea.Argcmd,
	}
}

func newCmd(c *internal.ExportCmd) *ast.Cmd {
	if c == nil {
		return nil
	}
	return &ast.Cmd{
		Name:   c.Name,
		Minlen: c.Minlen,
		Flags:  c.Flags,
		Parser: c.Parser,
	}
}
