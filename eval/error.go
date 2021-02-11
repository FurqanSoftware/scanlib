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

type ErrInvalidOperation struct{}

func (e ErrInvalidOperation) Error() string {
	return "invalid operation"
}

type ErrNonIntegerIndex struct{}

func (e ErrNonIntegerIndex) Error() string {
	return "non-integer index"
}

type ErrCheckError struct {
	Pos    Cursor
	Clause int
}

func (e ErrCheckError) Error() string {
	return fmt.Sprintf("%d:%d: check error (clause %d)", e.Pos.Ln, e.Pos.Col, e.Clause)
}

type ErrExpectedEOL struct{}

func (e ErrExpectedEOL) Error() string {
	return "expected EOL"
}

type ErrUnexpectedEOL struct{}

func (e ErrUnexpectedEOL) Error() string {
	return "unexpected EOL"
}

type ErrExpectedEOF struct{}

func (e ErrExpectedEOF) Error() string {
	return "expected EOF"
}

type ErrUnexpectedEOF struct{}

func (e ErrUnexpectedEOF) Error() string {
	return "unexpected EOF"
}

type Cursor struct {
	Ln, Col int
}
