package ast

// Pos represents node position.
type Pos struct {
	Offset int // offset, starting at 0
	Line   int // line number, starting at 1
	Column int // column number, starting at 1 (byte count)

	// Should I support Filename?
	// Filename string // filename, if any
}
