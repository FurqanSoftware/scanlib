package cpp14

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"git.furqansoftware.net/toph/scanlib/ast"
	"git.furqansoftware.net/toph/scanlib/gen/cpp14"
)

func TestGenerate(t *testing.T) {
	fis, err := ioutil.ReadDir("./testdata")
	if err != nil {
		t.Fatal(err)
	}
	for _, fi := range fis {
		codesrc, err := ioutil.ReadFile(filepath.Join("./testdata", fi.Name(), "cpp14.cpp"))
		if os.IsNotExist(err) {
			continue
		}

		t.Run(fi.Name(), func(t *testing.T) {
			specsrc, err := ioutil.ReadFile(filepath.Join("./testdata", fi.Name(), "scanspec"))
			if err != nil {
				t.Fatal(err)
			}
			n, err := ast.ParseString("inputspec", string(specsrc))
			if err != nil {
				t.Fatal(err)
			}

			code, err := cpp14.Generate(n)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(code, codesrc) {
				t.Errorf("want:\n\n%s\n\ngot:\n\n%s", codesrc, code)
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
