package cpp14

import (
	"testing"

	"git.furqansoftware.net/toph/scanlib/ast"
)

var sources = map[string]string{
	"add": `A = int(-9999999, 20000000)
B = int(-9999999, 20000000)
eol()
eof()
`,

	"grid": `R = int(1, 25)
C = int(1, 25)
eol()
G = make(int[R])
for i 0 R {
	G[i] = string(C, "*.")
	eol()
}
eof()
`,
}

func TestWalk(t *testing.T) {
	for _, s := range sources {
		n, err := ast.ParseString("inputspec", s)
		if err != nil {
			t.Fatal(err)
		}
		Generate(n)
	}
}
