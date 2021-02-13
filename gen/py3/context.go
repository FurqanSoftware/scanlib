// Copyright 2020 Furqan Software Ltd. All rights reserved.

package py3

import (
	"git.furqansoftware.net/toph/scanlib/gen/code"
)

type Context struct {
	types map[string]string
	cw    *code.Writer

	linevar bool
	scanarg int

	ozs map[interface{}]Optimization
}

type Optimization interface {
	Generate(ctx *Context) error
}

type Noop struct{}

func (o Noop) Generate(ctx *Context) error {
	return nil
}
