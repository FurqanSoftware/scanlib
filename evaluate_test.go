package scanlib

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"git.furqansoftware.net/toph/scanlib/ast"
	"git.furqansoftware.net/toph/scanlib/eval"
)

func TestEvaluate(t *testing.T) {
	fis, err := os.ReadDir("./testdata")
	if err != nil {
		t.Fatal(err)
	}
	for _, fi := range fis {
		t.Run(fi.Name(), func(t *testing.T) {
			_, err := os.Stat(filepath.Join("./testdata", fi.Name(), "_skip"))
			if err == nil {
				t.Skip()
			}

			specsrc, err := os.ReadFile(filepath.Join("./testdata", fi.Name(), "scanspec"))
			if err != nil {
				t.Fatal(err)
			}
			n, err := ast.ParseString("inputspec", string(specsrc))
			if err != nil {
				t.Fatal(err)
			}

			pis, err := os.ReadDir(filepath.Join("./testdata", fi.Name(), "inputs"))
			for _, pi := range pis {
				if !strings.HasSuffix(pi.Name(), ".in") {
					continue
				}
				t.Run(pi.Name(), func(t *testing.T) {
					instr, err := os.ReadFile(filepath.Join("./testdata", fi.Name(), "inputs", pi.Name()))
					if err != nil {
						t.Fatal(err)
					}

					errstr, _ := os.ReadFile(filepath.Join("./testdata", fi.Name(), "inputs", strings.TrimSuffix(pi.Name(), ".in")+".err"))

					_, err = eval.Evaluate(n, bytes.NewReader(instr))
					if err != nil {
						if err.Error() != string(errstr) {
							t.Fatalf("want err == %q, got %q", string(errstr), err.Error())
						}
					} else {
						if string(errstr) != "" {
							t.Fatalf("want err == %q, got nil", string(errstr))
						}
					}
				})
			}
		})
	}
}
