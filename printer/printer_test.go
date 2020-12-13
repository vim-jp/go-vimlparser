package printer

import (
	"bytes"
	"strings"
	"testing"

	"github.com/vim-jp/go-vimlparser"
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
