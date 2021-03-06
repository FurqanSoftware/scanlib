package eval

import (
	"io"
	"reflect"
)

type Context struct {
	Values map[string]reflect.Value
	Input  *Input
}

func NewContext(input io.Reader) (*Context, error) {
	p, err := NewInput(input)
	if err != nil {
		return nil, err
	}
	return &Context{
		Values: map[string]reflect.Value{},
		Input:  p,
	}, nil
}
