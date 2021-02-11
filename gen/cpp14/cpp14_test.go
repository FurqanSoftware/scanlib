package cpp14

import (
	"testing"

	"git.furqansoftware.net/toph/scanlib/ast"
)

func TestGenerate(t *testing.T) {
	for _, c := range []struct {
		label  string
		source string
		code   string
	}{
		{
			label:  "add",
			source: sourceAdd,
			code:   "#include <iostream>\n\nusing namespace std;\n\nint main() {\n\tint\t A\t,\t B\t;\n\tcin >> A;\n\tcin >> B;\n}\n",
		},
		{
			label:  "grid",
			source: sourceGrid,
			code:   "#include <iostream>\n\nusing namespace std;\n\nint main() {\n\tint\t R\t,\t C\t;\n\tcin >> R;\n\tcin >> C;\n\tstring\t G[\tR\t]\t;\n\tfor (int i = \t0\t; i < \tR\t; ++i) {\n\t\tcin >> G;\n\t}\n}\n",
		},
	} {
		n, err := ast.ParseString("inputspec", c.source)
		if err != nil {
			t.Fatal(err)
		}
		code, err := Generate(n)
		if err != nil {
			t.Fatal(err)
		}
		if string(code) != c.code {
			t.Errorf("want code == %q, got %q", c.code, code)
		}
	}
}

const sourceAdd = `var A, B int
scan A, B
check A >= 0, A < 10, B >= 0, B < 20
eol
eof
`

const sourceGrid = `var R, C int
scan R, C
check R >= 1, R < 25, C >= 1, C < 25
eol
var G [R]string
for i 0 R
	scan G[i]
	check len(G[i]) == C
	check re(G[i], "^[*.]+$")
	eol
end
eof
`
