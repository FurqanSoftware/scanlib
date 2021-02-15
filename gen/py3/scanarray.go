// Copyright 2020 Furqan Software Ltd. All rights reserved.

package py3

import (
	"git.furqansoftware.net/toph/scanlib/ast"
)

type scanArray struct {
	varDecl  *ast.VarDecl
	forStmt  *ast.ForStmt
	scanStmt *ast.ScanStmt
	eolStmt  *string
}

func (o scanArray) Generate(ctx *Context) error {
	for i, x := range o.scanStmt.RefList {
		if i > 0 {
			ctx.cw.Print(", ")
		}
		ctx.cw.Printf("%s", x.Identifier)
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
	for _, n := range n.Statement {
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
				expressionEqNumber(&n.ForStmt.Range.Low, "0") &&
				expressionEq(&oz.varDecl.VarSpec.Type.TypeLit.ArrayType.ArrayLength, &n.ForStmt.Range.High) {
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
	for _, n := range nfor.Block.Statement {
		switch {
		case n.ScanStmt != nil:
			if state == szero {
				if len(n.ScanStmt.RefList) == 1 &&
					n.ScanStmt.RefList[0].Identifier == oz.varDecl.VarSpec.IdentList[0] &&
					len(n.ScanStmt.RefList[0].Indices) == 1 &&
					expressionEqVariable(&n.ScanStmt.RefList[0].Indices[0], nfor.Range.Index) {
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

func expressionEq(a, b *ast.Expression) bool {
	av, ok := expressionVariable(a)
	if !ok {
		return false
	}
	bv, ok := expressionVariable(b)
	if !ok {
		return false
	}
	return av == bv
}

func expressionEqVariable(a *ast.Expression, b string) bool {
	av, ok := expressionVariable(a)
	if !ok {
		return false
	}
	return av == b
}

func expressionEqNumber(a *ast.Expression, b string) bool {
	av, ok := exprNumber(a)
	if !ok {
		return false
	}
	return av == b
}

func expressionVariable(e *ast.Expression) (string, bool) {
	if len(e.Right) == 0 &&
		len(e.Left.Right) == 0 &&
		len(e.Left.Left.Right) == 0 &&
		e.Left.Left.Left.Exponent == nil &&
		e.Left.Left.Left.Unary.Value != nil &&
		e.Left.Left.Left.Unary.Value.Variable != nil &&
		len(e.Left.Left.Left.Unary.Value.Variable.Indices) == 0 {
		return e.Left.Left.Left.Unary.Value.Variable.Identifier, true
	}
	return "", false
}

func exprNumber(e *ast.Expression) (string, bool) {
	if len(e.Right) == 0 &&
		len(e.Left.Right) == 0 &&
		len(e.Left.Left.Right) == 0 &&
		e.Left.Left.Left.Exponent == nil &&
		e.Left.Left.Left.Unary != nil &&
		e.Left.Left.Left.Unary.Value != nil &&
		e.Left.Left.Left.Unary.Value.Number != nil {
		return string(*e.Left.Left.Left.Unary.Value.Number), true
	}
	return "", false
}
