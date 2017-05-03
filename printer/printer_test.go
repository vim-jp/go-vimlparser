package printer

import (
	"bytes"
	"strings"
	"testing"

	"github.com/haya14busa/go-vimlparser"
)

func TestFprint_file(t *testing.T) {
	src := `let _ = 1`
	r := strings.NewReader(src)
	node, err := vimlparser.ParseFile(r, "", nil)
	if err != nil {
		t.Fatal(err)
	}
	buf := new(bytes.Buffer)
	if err := Fprint(buf, node, nil); err == nil {
		t.Error("want error")
	}
}

func TestFprint_expr(t *testing.T) {
	tests := []struct {
		in      string
		want    string
		wantErr bool
	}{
		{in: `xyz`, want: `xyz`},                       // Ident
		{in: `"double quote"`, want: `"double quote"`}, // BasicLit
		{in: `14`, want: `14`},                         // BasicLit
		{in: `x+1`, want: `x + 1`, wantErr: true},
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
