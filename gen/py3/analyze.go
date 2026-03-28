// Copyright 2020 Furqan Software Ltd. All rights reserved.

package py3

import (
	"git.furqansoftware.net/toph/scanlib/ast"
)

type analyzer struct {
	ozs map[ast.Node]Optimization

	blockEOLs    map[*ast.Block]bool
	blockAssigns map[*ast.Block]bool
	scanned      map[string]bool
}

func analyze(n *ast.Source) *analyzer {
	a := analyzer{
		ozs:          map[ast.Node]Optimization{},
		blockEOLs:    map[*ast.Block]bool{},
		blockAssigns: map[*ast.Block]bool{},
		scanned:      map[string]bool{},
	}
	findBlockEOLs(&a, n)
	findBlockAssigns(&a, n)
	findScanned(&a, n)
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
		case *ast.Source, *ast.Block, *ast.Statement, *ast.ForStmt, *ast.IfStmt, *ast.IfBranch:
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

func findBlockAssigns(a *analyzer, n *ast.Source) {
	stack := []ast.Node{}
	ast.Inspect(n, func(n ast.Node) bool {
		if n == nil {
			stack = stack[:len(stack)-1]
		}
		switch n := n.(type) {
		case *ast.Source, *ast.Block, *ast.Statement, *ast.ForStmt, *ast.IfStmt, *ast.IfBranch:
			stack = append(stack, n)
			return true
		case *ast.AssignStmt:
			for _, n := range stack {
				b, ok := n.(*ast.Block)
				if ok {
					a.blockAssigns[b] = true
				}
			}
			return false
		}
		return false
	})
}

func findScanned(a *analyzer, n *ast.Source) {
	ast.Inspect(n, func(n ast.Node) bool {
		switch n := n.(type) {
		case *ast.Source, *ast.Block, *ast.Statement, *ast.ForStmt, *ast.IfStmt, *ast.IfBranch:
			return true
		case *ast.ScanStmt:
			for _, ref := range n.RefList {
				a.scanned[ref.Ident] = true
			}
			return false
		case *ast.ScanlnStmt:
			for _, ref := range n.RefList {
				a.scanned[ref.Ident] = true
			}
			return false
		}
		return false
	})
}

type State int
