// Copyright 2020 Furqan Software Ltd. All rights reserved.

package py3

import (
	"git.furqansoftware.net/toph/scanlib/ast"
)

type scanArray struct {
	varDecl  *ast.VarDecl
	forStmt  *ast.ForStmt
	scanStmt *ast.ScanStmt
	eolStmt  *ast.EOLStmt
}

func (o scanArray) Generate(ctx *Context) error {
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

func analyzeBlockScanArray(ctx *Context, n *ast.Block) {
	const (
		szero int = iota
		svar
		sfor
		seol
	)

	oz := scanArray{}

	var state = szero
	for _, n := range n.Statements {
		switch {
		case n.VarDecl != nil:
			if state == szero {
				if len(n.VarDecl.VarSpec.IdentList) == 1 &&
					n.VarDecl.VarSpec.Type.TypeLit != nil &&
					n.VarDecl.VarSpec.Type.TypeLit.ArrayType.ElementType.TypeName != nil {
					state = svar
					oz.varDecl = n.VarDecl
				}
			} else {
				state = szero
			}

		case n.CheckStmt != nil:

		case n.ForStmt != nil:
			if state == svar &&
				exprEqInt64(&n.ForStmt.Range.Low, 0) &&
				exprEq(&oz.varDecl.VarSpec.Type.TypeLit.ArrayType.ArrayLength, &n.ForStmt.Range.High) {
				scanstmt, ok := analyzeBlockScanArraySfor(ctx, n.ForStmt, &oz)
				if ok {
					state = sfor
					oz.forStmt = n.ForStmt
					oz.scanStmt = scanstmt
				} else {
					state = szero
				}

			} else {
				state = szero

				analyzeBlock(ctx, &n.ForStmt.Block)
			}

		case n.EOLStmt != nil:
			if state == sfor {
				state = seol
				oz.eolStmt = n.EOLStmt
				ctx.ozs[oz.forStmt] = oz
				ctx.ozs[oz.varDecl] = Noop{}
				ctx.ozs[oz.eolStmt] = Noop{}

			} else {
				state = szero
			}

		default:
			state = szero
		}
	}
}

func analyzeBlockScanArraySfor(ctx *Context, nfor *ast.ForStmt, oz *scanArray) (*ast.ScanStmt, bool) {
	const (
		szero int = iota
		sscan
	)

	var scanstmt *ast.ScanStmt

	var state = szero
	for _, n := range nfor.Block.Statements {
		switch {
		case n.ScanStmt != nil:
			if state == szero {
				if len(n.ScanStmt.RefList) == 1 &&
					n.ScanStmt.RefList[0].Ident == oz.varDecl.VarSpec.IdentList[0] &&
					len(n.ScanStmt.RefList[0].Indices) == 1 &&
					exprEqVar(&n.ScanStmt.RefList[0].Indices[0], nfor.Range.Index) {
					state = sscan
					scanstmt = n.ScanStmt
				}
			} else {
				return nil, false
			}

		case n.CheckStmt != nil:

		default:
			return nil, false
		}
	}
	return scanstmt, true
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
