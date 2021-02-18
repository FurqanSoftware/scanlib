// Copyright 2020 Furqan Software Ltd. All rights reserved.

package py3

import (
	"git.furqansoftware.net/toph/scanlib/ast"
)

type analyzer struct {
	ozs map[ast.Node]Optimization
}

func analyze(n *ast.Source) *analyzer {
	a := analyzer{
		ozs: map[ast.Node]Optimization{},
	}
	ast.Walk(&a, n)
	return &a
}

func (a *analyzer) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return nil
	}

	switch n := n.(type) {
	case *ast.Source, *ast.Statement, *ast.ForStmt:
		return a

	case *ast.Block:
		a.multiVar(n)
		a.onlyToken(n)
		a.arrayLine(n)
		return a
	}

	return nil
}

type State int
