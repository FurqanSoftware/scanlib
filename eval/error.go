package eval

import (
	"fmt"

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
	Clause int
	Cursor Cursor
}

func (e ErrCheckError) Error() string {
	return fmt.Sprintf("%d:%d: check error (clause %d, cursor %d:%d)", e.Pos.Line, e.Pos.Column, e.Clause, e.Cursor.Ln, e.Cursor.Col)
}

type ErrExpectedEOL struct {
	Pos lexer.Position
}

func (e ErrExpectedEOL) Error() string {
	return fmt.Sprintf("%d:%d: want EOL", e.Pos.Line, e.Pos.Column)

}

type ErrUnexpectedEOL struct {
	Pos lexer.Position
}

func (e ErrUnexpectedEOL) Error() string {
	return fmt.Sprintf("%d:%d: unwanted EOL", e.Pos.Line, e.Pos.Column)
}

type ErrExpectedEOF struct {
	Pos   lexer.Position
	Token []byte
}

func (e ErrExpectedEOF) Error() string {
	return fmt.Sprintf("%d:%d: want EOF, got trailing %q", e.Pos.Line, e.Pos.Column, e.Token)
}

type ErrUnexpectedEOF struct {
	Pos lexer.Position
}

func (e ErrUnexpectedEOF) Error() string {
	return fmt.Sprintf("%d:%d: unwanted EOF", e.Pos.Line, e.Pos.Column)
}

type ErrBadParse struct {
	Want   string
	Got    []byte
	Cursor Cursor
}

func (e ErrBadParse) Error() string {
	return fmt.Sprintf("parse error (cursor %d:%d): want %s, got %q", e.Cursor.Ln, e.Cursor.Col, e.Want, e.Got)
}

type Cursor struct {
	Ln, Col int
}

func catch(err error) {
	if err != nil {
		panic(err)
	}
}
