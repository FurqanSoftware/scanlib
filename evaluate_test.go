package scanlib

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"git.furqansoftware.net/toph/scanlib/ast"
	"git.furqansoftware.net/toph/scanlib/eval"
)

func TestEvaluate(t *testing.T) {
	fis, err := ioutil.ReadDir("./testdata")
	if err != nil {
		t.Fatal(err)
	}
	for _, fi := range fis {
		t.Run(fi.Name(), func(t *testing.T) {
			specsrc, err := ioutil.ReadFile(filepath.Join("./testdata", fi.Name(), "scanspec"))
			if err != nil {
				t.Fatal(err)
			}
			n, err := ast.ParseString("inputspec", string(specsrc))
			if err != nil {
				t.Fatal(err)
			}

			pis, err := ioutil.ReadDir(filepath.Join("./testdata", fi.Name(), "inputs"))
			for _, pi := range pis {
				if !strings.HasSuffix(pi.Name(), ".in") {
					continue
				}
				t.Run(pi.Name(), func(t *testing.T) {
					instr, err := ioutil.ReadFile(filepath.Join("./testdata", fi.Name(), "inputs", pi.Name()))
					if err != nil {
						t.Fatal(err)
					}

					errstr, _ := ioutil.ReadFile(filepath.Join("./testdata", fi.Name(), "inputs", strings.TrimSuffix(pi.Name(), ".in")+".err"))

					ctx := eval.NewContext(bytes.NewReader(instr))
					_, err = eval.Evaluate(ctx, n)
					if err != nil {
						if err.Error() != string(errstr) {
							t.Fatalf("want err == %q, got %q", string(errstr), err.Error())
						}
					}
				})
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
