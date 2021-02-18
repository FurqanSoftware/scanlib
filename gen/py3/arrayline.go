// Copyright 2020 Furqan Software Ltd. All rights reserved.

package py3

import (
	"git.furqansoftware.net/toph/scanlib/ast"
)

type arrayLine struct {
	varDecl  *ast.VarDecl
	forStmt  *ast.ForStmt
	scanStmt *ast.ScanStmt
	eolStmt  *ast.EOLStmt
}

func (o arrayLine) Generate(ctx *Context) error {
	for i, x := range o.scanStmt.RefList {
		if i > 0 {
			ctx.cw.Print(", ")
		}
		ctx.cw.Printf("%s", x.Ident)
	}
	t := ASTType[*o.varDecl.VarSpec.Type.TypeLit.ArrayType.ElementType.TypeName]
	ctx.cw.Printf(" = map(%s, input().split())", t)
	ctx.cw.Println()
	return nil
}

func (a *analyzer) arrayLine(n *ast.Block) {
	const (
		zero State = iota
		decl
		loop
		scan
	)

	oz := arrayLine{}

	var state State
	ast.Inspect(n, func(n ast.Node) bool {
		if n == nil {
			return false
		}

		switch state {
		case zero:
			switch n := n.(type) {
			case *ast.Block, *ast.Statement:
				return true

			case *ast.VarDecl:
				if len(n.VarSpec.IdentList) == 1 &&
					n.VarSpec.Type.TypeLit != nil &&
					n.VarSpec.Type.TypeLit.ArrayType.ElementType.TypeName != nil {
					oz.varDecl = n
					state++
				}
				return false

			default:
				return false
			}

		case decl:
			switch n := n.(type) {
			case *ast.Statement:
				return true

			case *ast.CheckStmt:
				return false

			case *ast.ForStmt:
				if exprEqInt64(&n.Range.Low, 0) &&
					exprEq(&oz.varDecl.VarSpec.Type.TypeLit.ArrayType.ArrayLength, &n.Range.High) {
					oz.forStmt = n
					state++
					return true
				}

			default:
				state = zero
				return false
			}

		case loop:
			switch n := n.(type) {
			case *ast.Block, *ast.Statement:
				return true

			case *ast.CheckStmt, *ast.RangeClause:
				return false

			case *ast.ScanStmt:
				if len(n.RefList) == 1 &&
					n.RefList[0].Ident == oz.varDecl.VarSpec.IdentList[0] &&
					len(n.RefList[0].Indices) == 1 &&
					exprEqVar(&n.RefList[0].Indices[0], oz.forStmt.Range.Index) {
					oz.scanStmt = n
					state++
				}
				return false

			default:
				state = zero
				return false
			}

		case scan:
			switch n := n.(type) {
			case *ast.Statement:
				return true

			case *ast.CheckStmt:
				return false

			case *ast.EOLStmt:
				oz.eolStmt = n
				a.ozs[oz.forStmt] = oz
				a.ozs[oz.varDecl] = Noop{}
				a.ozs[oz.eolStmt] = Noop{}
				state = zero
				return false

			default:
				state = zero
			}
		}

		return false
	})
}

func exprEq(a, b *ast.Expr) bool {
	av, ok := exprVar(a)
	if !ok {
		return false
	}
	bv, ok := exprVar(b)
	if !ok {
		return false
	}
	return av == bv
}

func exprEqVar(a *ast.Expr, b string) bool {
	av, ok := exprVar(a)
	if !ok {
		return false
	}
	return av == b
}

func exprEqInt64(a *ast.Expr, b int64) bool {
	av, ok := exprInt64(a)
	if !ok {
		return false
	}
	return av == b
}

func exprVar(e *ast.Expr) (string, bool) {
	if len(e.Right) == 0 &&
		len(e.Left.Right) == 0 &&
		len(e.Left.Left.Right) == 0 &&
		e.Left.Left.Left.Left != nil &&
		e.Left.Left.Left.Left.Left != nil &&
		e.Left.Left.Left.Left.Left.Exponent == nil &&
		e.Left.Left.Left.Left.Left.Unary.Value != nil &&
		e.Left.Left.Left.Left.Left.Unary.Value.Variable != nil &&
		len(e.Left.Left.Left.Left.Left.Unary.Value.Variable.Indices) == 0 {
		return e.Left.Left.Left.Left.Left.Unary.Value.Variable.Ident, true
	}
	return "", false
}

func exprInt64(e *ast.Expr) (int64, bool) {
	if len(e.Right) == 0 &&
		len(e.Left.Right) == 0 &&
		len(e.Left.Left.Right) == 0 &&
		e.Left.Left.Left.Left != nil &&
		e.Left.Left.Left.Left.Left != nil &&
		e.Left.Left.Left.Left.Left.Exponent == nil &&
		e.Left.Left.Left.Left.Left.Unary != nil &&
		e.Left.Left.Left.Left.Left.Unary.Value != nil &&
		e.Left.Left.Left.Left.Left.Unary.Value.BasicLit != nil &&
		e.Left.Left.Left.Left.Left.Unary.Value.BasicLit.IntLit != nil {
		return *e.Left.Left.Left.Left.Left.Unary.Value.BasicLit.IntLit, true
	}
	return 0, false
}
