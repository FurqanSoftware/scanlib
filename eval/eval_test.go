package eval

import (
	"strings"
	"testing"

	"git.furqansoftware.net/toph/scanlib/ast"
)

func TestEvaluate(t *testing.T) {
	for _, c := range []struct {
		label  string
		source string
		input  string
		err    error
		errstr string
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
			err:    ErrCheckError{Pos: Cursor{3, 1}, Clause: 1},
		},
		{
			label:  "add/check#02",
			source: sourceAdd,
			input:  "3 200\n",
			err:    ErrCheckError{Pos: Cursor{3, 1}, Clause: 4},
		},
		{
			label:  "add/check#03",
			source: sourceAdd,
			input:  "3 2 33\n",
			err:    ErrExpectedEOL{Pos: Cursor{4, 1}},
		},
		{
			label:  "grid/ok#01",
			source: sourceGrid,
			input:  "3 5\n**...\n..*..\n....*\n",
		},
		{
			label:  "grid/ok#02",
			source: sourceGrid,
			input:  "3 5\n**...\n..*..\n.x..*\n",
			err:    ErrCheckError{Pos: Cursor{9, 2}, Clause: 1},
		},
		{
			label:  "n1018/ok#01",
			source: sourceN1018,
			input:  "1",
		},
		{
			label:  "n1018/ok#02",
			source: sourceN1018,
			input:  "100000000000000000",
		},
		{
			label:  "n1018/check#01",
			source: sourceN1018,
			input:  "1000000000000000001",
			err:    ErrCheckError{Pos: Cursor{3, 1}, Clause: 2},
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

const sourceN1018 = `var N int64
scan N
check N >= 1, N <= 1000000000000000000
`
