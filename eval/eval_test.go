package eval

import (
	"strings"
	"testing"

	"git.furqansoftware.net/toph/inputlib/ast"
)

func TestEvaluate(t *testing.T) {
	for _, c := range []struct {
		label  string
		source string
		input  string
		err    error
	}{
		{
			label:  "add/ok#01",
			source: sourceAdd,
			input:  "5 6\n",
		},
		{
			label:  "add/ok#02",
			source: sourceAdd,
			input:  "3 2\n",
		},
		{
			label:  "add/check#01",
			source: sourceAdd,
			input:  "-3 2\n",
			err:    ErrCheckError{},
		},
		{
			label:  "add/check#02",
			source: sourceAdd,
			input:  "3 200\n",
			err:    ErrCheckError{},
		},
		{
			label:  "grid/ok#01",
			source: sourceGrid,
			input:  "3 5\n**...\n..*..\n....*\n",
		},
		{
			label:  "grid/ok#01",
			source: sourceGrid,
			input:  "3 5\n**...\n..*..\n.x..*\n",
			err:    ErrCheckError{},
		},
	} {
		t.Run(c.label, func(t *testing.T) {
			n, err := ast.ParseString("inputspec", c.source)
			if err != nil {
				t.Fatal(err)
			}
			ctx := NewContext(strings.NewReader(c.input))
			_, err = Evaluate(ctx, n)
			if c.err != err {
				t.Errorf("want err == %v, got %v", c.err, err)
			}
		})
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
