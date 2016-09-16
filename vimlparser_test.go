package vimlparser

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseFile_can_parse(t *testing.T) {
	match, err := filepath.Glob("test/test_*.vim")
	if err != nil {
		t.Fatal(err)
	}
	okErr := "go-vimlparser:Parse: vimlparser:"
	match = append(match, "autoload/vimlparser.vim")
	match = append(match, "go/gocompiler.vim")
	for _, filename := range match {
		if err := checkParse(t, filename); err != nil && !strings.HasPrefix(err.Error(), okErr) {
			t.Errorf("%s: %v", filename, err)
		}
	}
}

func checkParse(t testing.TB, filename string) error {
	in, err := os.Open(filename)
	if err != nil {
		t.Error(err)
	}
	defer in.Close()
	_, err = ParseFile(in, nil)
	return err
}

func BenchmarkParseFile(b *testing.B) {
	filename := "autoload/vimlparser.vim"
	for i := 0; i < b.N; i++ {
		checkParse(b, filename)
	}
}

func TestParse_Compile(t *testing.T) {
	node, err := Parse(strings.NewReader("let x = 1"), nil)
	if err != nil {
		t.Fatal(err)
	}
	b := new(bytes.Buffer)
	if err := Compile(b, node); err != nil {
		t.Fatal(err)
	}
	if got, want := b.String(), "(let = x 1)"; got != want {
		t.Errorf("Compile(Parse(\"let x = 1\")) = %v, want %v", got, want)
	}
}

func TestParse_Compile_err(t *testing.T) {
	want := "go-vimlparser:Parse: vimlparser: E492: Not an editor command: hoge: line 1 col 1"
	_, err := Parse(strings.NewReader("hoge"), nil)
	if err != nil {
		if got := err.Error(); want != got {
			t.Errorf("Parse(\"hoge\") = %v, want %v", got, want)
		}
	}
}

func TestParseExpr_Compile(t *testing.T) {
	node, err := ParseExpr(strings.NewReader("x + 1"))
	if err != nil {
		t.Fatal(err)
	}
	b := new(bytes.Buffer)
	if err := Compile(b, node); err != nil {
		t.Fatal(err)
	}
	if got, want := b.String(), "(+ x 1)"; got != want {
		t.Errorf("Compile(Parse(\"x + 1\")) = %v, want %v", got, want)
	}
}

func TestParseExpr_Compile_err(t *testing.T) {
	want := "go-vimlparser:Parse: vimlparser: unexpected token: /: line 1 col 4"
	_, err := ParseExpr(strings.NewReader("1 // 2"))
	if err != nil {
		if got := err.Error(); want != got {
			t.Errorf("ParseExpr(\"1 // 2\") = %v, want %v", got, want)
		}
	}
}
