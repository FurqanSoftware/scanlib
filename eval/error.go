package eval

import (
	"fmt"

	"git.furqansoftware.net/toph/scanlib/ast"
	"github.com/alecthomas/participle/v2/lexer"
)

type ErrUndefined struct {
	Pos  lexer.Position
	Name string
}

func (e ErrUndefined) Error() string {
	return fmt.Sprintf("%d:%d: undefined: "+e.Name, e.Pos.Line, e.Pos.Column)
}

type ErrCantScanType struct{}

func (e ErrCantScanType) Error() string {
	return "can't scan type"
}

type ErrInvalidArgument struct{}

func (e ErrInvalidArgument) Error() string {
	return "invalid argument"
}

type ErrInvalidOperation struct {
	Pos lexer.Position
}

func (e ErrInvalidOperation) Error() string {
	return fmt.Sprintf("%d:%d: invalid operation", e.Pos.Line, e.Pos.Column)
}

type ErrNonIntegerIndex struct {
	Pos lexer.Position
}

func (e ErrNonIntegerIndex) Error() string {
	return fmt.Sprintf("%d:%d: non-integer index", e.Pos.Line, e.Pos.Column)
}

type ErrInvalidIndex struct {
	Pos    lexer.Position
	Index  int
	Length int
}

func (e ErrInvalidIndex) Error() string {
	if e.Index >= 0 {
		return fmt.Sprintf("%d:%d: invalid array index %d (out of bounds for %d-element array)", e.Pos.Line, e.Pos.Column, e.Index, e.Length)
	} else {
		return fmt.Sprintf("%d:%d: invalid array index %d (index must be non-negative)", e.Pos.Line, e.Pos.Column, e.Index)
	}
}

type ErrCheckError struct {
	Pos    lexer.Position
	Cursor Cursor
	Expr   *ast.Expr
	Values Values
}

func (e ErrCheckError) Error() string {
	vars := []string{}
	varsseen := map[string]bool{}
	ast.Inspect(e.Expr, func(n ast.Node) bool {
		switch n := n.(type) {
		case *ast.Variable:
			if varsseen[n.Ident] || len(n.Indices) > 0 {
				return true
			}
			varsseen[n.Ident] = true
			vars = append(vars, n.Ident)
		}
		return len(vars) < 3
	})
	s := []byte{}
	for _, t := range e.Expr.Tokens {
		s = append(s, []byte(t.Value)...)
	}
	msg := fmt.Sprintf("%d:%d: check error %s", e.Pos.Line, e.Pos.Column, ellipsize(s, 20))
	if len(vars) > 0 {
		msg += " ("
		for i, k := range vars {
			if i > 0 {
				msg += ", "
			}
			msg += fmt.Sprintf("%s=%#v", k, e.Values[k].Elem().Interface())
		}
		msg += ")"
	}
	return msg
}

type ErrExpectedEOL struct {
	Pos    lexer.Position
	Got    []byte
	Cursor Cursor
}

func (e ErrExpectedEOL) Error() string {
	return fmt.Sprintf("%d:%d: (cursor %d:%d): want EOL, got %q", e.Pos.Line, e.Pos.Column, e.Cursor.Ln, e.Cursor.Col, e.Got)

}

type ErrUnexpectedEOL struct {
	Pos lexer.Position
}

func (e ErrUnexpectedEOL) Error() string {
	return fmt.Sprintf("%d:%d: unwanted EOL", e.Pos.Line, e.Pos.Column)
}

type ErrExpectedEOF struct {
	Pos lexer.Position
	Got []byte
}

func (e ErrExpectedEOF) Error() string {
	return fmt.Sprintf("%d:%d: want EOF, got trailing %q", e.Pos.Line, e.Pos.Column, e.Got)
}

type ErrUnexpectedEOF struct {
	Pos lexer.Position
}

func (e ErrUnexpectedEOF) Error() string {
	return fmt.Sprintf("%d:%d: unwanted EOF", e.Pos.Line, e.Pos.Column)
}

type ErrBadParse struct {
	Pos    lexer.Position
	Want   string
	Got    []byte
	Cursor Cursor
}

func (e ErrBadParse) Error() string {
	return fmt.Sprintf("%d:%d: parse error (cursor %d:%d): want %s, got %q", e.Pos.Line, e.Pos.Column, e.Cursor.Ln, e.Cursor.Col, e.Want, e.Got)
}

type Cursor struct {
	Ln, Col int
}

func ellipsize(b []byte, n int) string {
	r := []rune(string(b))
	if len(r) > n {
		return string(r[:n]) + "..."
	}
	return string(r)
}

func catch(err error) {
	if err != nil {
		panic(err)
	}
}
