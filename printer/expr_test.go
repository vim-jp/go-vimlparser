package printer

import (
	"bytes"
	"strings"
	"testing"

	vimlparser "github.com/haya14busa/go-vimlparser"
)

func TestFprint_expr(t *testing.T) {
	tests := []struct {
		in      string
		want    string
		wantErr bool
	}{
		{in: `xyz`, want: `xyz`},                        // Ident
		{in: `"double quote"`, want: `"double quote"`},  // BasicLit
		{in: `14`, want: `14`},                          // BasicLit
		{in: `+1`, want: `+1`},                          // UnaryExpr
		{in: `-  1`, want: `-1`},                        // UnaryExpr
		{in: `! + - 1`, want: `!+-1`},                   // UnaryExpr
		{in: `x+1`, want: `x + 1`},                      // BinaryExpr
		{in: `1+2*3`, want: `1 + 2 * 3`},                // BinaryExpr
		{in: `1*2+3`, want: `1 * 2 + 3`},                // BinaryExpr
		{in: `(1+2)*(3-4)`, want: `(1 + 2) * (3 - 4)`},  // ParenExpr
		{in: `1+(2*3)`, want: `1 + (2 * 3)`},            // ParenExpr
		{in: `(((x+(1))))`, want: `(x + (1))`},          // ParenExpr
		{in: `x+1==14 ||-1`, want: `x + 1 == 14 || -1`}, // BinaryExpr
	}

	for _, tt := range tests {
		r := strings.NewReader(tt.in)
		node, err := vimlparser.ParseExpr(r)
		if err != nil {
			t.Fatal(err)
		}
		buf := new(bytes.Buffer)
		if err := Fprint(buf, node, nil); err != nil {
			if !tt.wantErr {
				t.Errorf("got unexpected error: %v", err)
			}
			continue
		}
		if got := buf.String(); got != tt.want {
			t.Errorf("got: %v, want: %v", got, tt.want)
		}
	}
}
