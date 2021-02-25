// Copyright 2020 Furqan Software Ltd. All rights reserved.

package py3

import (
	"fmt"

	"git.furqansoftware.net/toph/scanlib/ast"
)

type analyzer struct {
	ozs map[ast.Node]Optimization

	blockEOLs map[*ast.Block]bool
}

func analyze(n *ast.Source) *analyzer {
	a := analyzer{
		ozs:       map[ast.Node]Optimization{},
		blockEOLs: map[*ast.Block]bool{},
	}
	findBlockEOLs(&a, n)
	fmt.Println(a.blockEOLs)
	ast.Walk(&a, n)
	return &a
}

func (a *analyzer) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return nil
	}

	switch n := n.(type) {
	case *ast.Source, *ast.Statement, *ast.ForStmt, *ast.IfStmt, *ast.IfBranch:
		return a

	case *ast.Block:
		a.multiVar(n)
		a.onlyToken(n)
		a.arrayLine(n)
		a.sameLine(n)
		return a
	}

	return nil
}

func findBlockEOLs(a *analyzer, n *ast.Source) {
	stack := []ast.Node{}
	ast.Inspect(n, func(n ast.Node) bool {
		if n == nil {
			stack = stack[:len(stack)-1]
		}
		switch n := n.(type) {
		case *ast.Source, *ast.Block, *ast.Statement, *ast.ForStmt, *ast.IfStmt:
			stack = append(stack, n)
			return true
		case *ast.EOLStmt:
			for _, n := range stack {
				b, ok := n.(*ast.Block)
				if ok {
					a.blockEOLs[b] = true
				}
			}
			return false
		}
		return false
	})
}

type State int
