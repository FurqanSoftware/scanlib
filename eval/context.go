package eval

import (
	"io"
)

type Context struct {
	Values map[string]Value
	Input  io.Reader
}

func NewContext(input io.Reader) *Context {
	return &Context{
		Values: map[string]Value{},
		Input:  input,
	}
}

func (c Context) GetValue(key string, indices []int) Value {
	r := c.Values[key]
	for _, i := range indices {
		x, ok := r.Data.([]Value)
		if !ok {
			return Value{}
		}
		r = x[i]
	}
	return r
}

func (c Context) SetValue(key string, indices []int, v Value) {
	if len(indices) == 0 {
		c.Values[key] = v
	} else {
		r := c.Values[key]
		for j, i := range indices {
			x, ok := r.Data.([]Value)
			if !ok {
				return
			}
			if j == len(indices)-1 {
				x[i] = v
			} else {
				r = x[i]
			}
		}
	}
}
