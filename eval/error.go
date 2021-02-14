package eval

import "fmt"

type ErrUndefined struct {
	Name string
}

func (e ErrUndefined) Error() string {
	return "undefined: " + e.Name
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
	Pos Cursor
}

func (e ErrInvalidOperation) Error() string {
	return fmt.Sprintf("%d:%d: invalid operation", e.Pos.Ln, e.Pos.Col)
}

type ErrNonIntegerIndex struct {
	Pos Cursor
}

func (e ErrNonIntegerIndex) Error() string {
	return fmt.Sprintf("%d:%d: non-integer index", e.Pos.Ln, e.Pos.Col)
}

type ErrCheckError struct {
	Pos    Cursor
	Clause int
}

func (e ErrCheckError) Error() string {
	return fmt.Sprintf("%d:%d: check error (clause %d)", e.Pos.Ln, e.Pos.Col, e.Clause)
}

type ErrExpectedEOL struct {
	Pos Cursor
}

func (e ErrExpectedEOL) Error() string {
	return fmt.Sprintf("%d:%d: want EOL", e.Pos.Ln, e.Pos.Col)

}

type ErrUnexpectedEOL struct {
	Pos Cursor
}

func (e ErrUnexpectedEOL) Error() string {
	return fmt.Sprintf("%d:%d: unwanted EOL", e.Pos.Ln, e.Pos.Col)
}

type ErrExpectedEOF struct {
	Pos   Cursor
	Token []byte
}

func (e ErrExpectedEOF) Error() string {
	return fmt.Sprintf("%d:%d: want EOF, got trailing %q", e.Pos.Ln, e.Pos.Col, e.Token)
}

type ErrUnexpectedEOF struct {
	Pos Cursor
}

func (e ErrUnexpectedEOF) Error() string {
	return fmt.Sprintf("%d:%d: unwanted EOF", e.Pos.Ln, e.Pos.Col)
}

type ErrBadParse struct {
	Pos  Cursor
	Want string
	Got  []byte
}

func (e ErrBadParse) Error() string {
	return fmt.Sprintf("Ln %d, Col %d: parse error: want %s, got %q", e.Pos.Ln, e.Pos.Col, e.Want, e.Got)
}

type Cursor struct {
	Ln, Col int
}
