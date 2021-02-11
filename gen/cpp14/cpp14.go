// Copyright 2020 Furqan Software Ltd. All rights reserved.

package cpp14

import (
	"bytes"
	"errors"
	"fmt"

	"git.furqansoftware.net/toph/inputlib/ast"
)

type Context struct {
	Source bytes.Buffer
}

func Generate(n *ast.Specification) []byte {
	ctx := Context{
		Source: bytes.Buffer{},
	}
	GenerateSpecification(n, &ctx)
	return nil
}

func GenerateSpecification(n *ast.Specification, ctx *Context) {
	GenerateBlock(&n.Block, ctx)
}

func GenerateBlock(n *ast.Block, ctx *Context) {
	for _, s := range n.Statement {
		GenerateStatement(s, ctx)
	}
}

func GenerateStatement(n *ast.Statement, ctx *Context) {
	switch {
	case n.Assignment != nil:
		GenerateAssignment(n.Assignment, ctx)
	case n.Call != nil:
		GenerateCall(n.Call, ctx)
	case n.For != nil:
		GenerateFor(n.For, ctx)
	}
}

func GenerateAssignment(n *ast.Assignment, ctx *Context) {
	switch {
	case n.Variable != nil:
		ctx.Source.WriteString(fmt.Sprintf("int %s;\n", n.Variable))
		ctx.Source.WriteString(fmt.Sprintf("cin >> %s;\n"))

	case n.Array != nil:
		ctx.Source.WriteString(fmt.Sprintf("int %s;\n", n.Variable))
		//
	}
}

func GenerateValue(n *ast.Value, ctx *Context) {
	switch {
	case n.Make != nil:
		GenerateMake(n.Make, ctx)
	case n.Call != nil:
		GenerateCall(n.Call, ctx)
	case n.Variable != nil:
		//
	case n.Integer != nil:
		//
	case n.String != nil:
		//
	}
}

func GenerateMake(n *ast.Make, ctx *Context) {
	ctx.Source.WriteString(fmt.Sprintf("vector<%s> %s(%s);", n.Array.Variable))
	s, err := GenerateIndex(&n.Array.Index, ctx)
	if err != nil {
		return nil, err
	}
	si, ok := s.(int)
	if !ok {
		return nil, errors.New("bad type")
	}
	return make([]int, si), nil
}

func GenerateCall(n *ast.Call, ctx *Context) {
	f, ok := Functions[n.Name]
	if !ok {
		return nil, errors.New("undefined: " + n.Name)
	}
	args := []interface{}{}
	for _, a := range n.Arguments {
		v, err := GenerateValue(a, ctx)
		args = append(args, v)
	}
	return f(ctx, args...)
}

func GenerateIndex(n *ast.Index, ctx *Context) {
	switch {
	case n.Variable != nil:
		return ctx.Variables[*n.Variable], nil

	case n.Integer != nil:
		return *n.Integer, nil
	}
}

func GenerateFor(n *ast.For, ctx *Context) {
	i, err := GenerateIndex(&n.From, ctx)
	if err != nil {
		return nil, err
	}
	ii, ok := i.(int)
	if !ok {
		return nil, errors.New("bad type")
	}

	j, err := GenerateIndex(&n.To, ctx)
	if err != nil {
		return nil, err
	}
	jj, ok := j.(int)
	if !ok {
		return nil, errors.New("bad type")
	}

	for ; ii < jj; ii++ {
		ctx.Variables[n.Variable] = ii
		_, err := GenerateBlock(&n.Block, ctx)
	}

	return nil, nil
}
