package scanlib

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"git.furqansoftware.net/toph/scanlib/ast"
	"git.furqansoftware.net/toph/scanlib/gen/cpp14"
	"git.furqansoftware.net/toph/scanlib/gen/go1"
	"git.furqansoftware.net/toph/scanlib/gen/py3"
)

type language struct {
	key   string
	ext   string
	genFn func(*ast.Source) ([]byte, error)
}

var (
	langs = []language{
		{
			key:   "cpp14",
			ext:   ".cpp",
			genFn: cpp14.Generate,
		},
		{
			key:   "go1",
			ext:   ".go",
			genFn: go1.Generate,
		},
		{
			key:   "py3",
			ext:   ".py",
			genFn: py3.Generate,
		},
	}
)

func TestGenerate(t *testing.T) {
	fis, err := ioutil.ReadDir("./testdata")
	if err != nil {
		t.Fatal(err)
	}
	for _, l := range langs {
		t.Run(l.key, func(t *testing.T) {
			for _, fi := range fis {
				codesrc, err := ioutil.ReadFile(filepath.Join("./testdata", fi.Name(), l.key+l.ext))
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

					code, err := l.genFn(n)
					if err != nil {
						t.Fatal(err)
					}
					if !bytes.Equal(code, codesrc) {
						t.Errorf("want:\n\n%s\n\ngot:\n\n%s", codesrc, code)
					}
				})
			}
		})
	}
}
