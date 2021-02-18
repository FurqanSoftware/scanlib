// Copyright 2020 Furqan Software Ltd. All rights reserved.

package py3

import (
	"strings"

	"git.furqansoftware.net/toph/scanlib/ast"
)

type scanOne struct {
	scanStmt *ast.ScanStmt
	eolStmt  *ast.EOLStmt
}

func (o scanOne) Generate(ctx *Context) error {
	ctx.cw.Print(o.scanStmt.RefList[0].Ident)
	for _, i := range o.scanStmt.RefList[0].Indices {
		ctx.cw.Print("[")
		err := genExpr(ctx, &i)
		if err != nil {
			return err
		}
		ctx.cw.Print("]")
	}
	t := ctx.types[o.scanStmt.RefList[0].Ident+strings.Repeat("[]", len(o.scanStmt.RefList[0].Indices))]
	if t == "string" {
		ctx.cw.Printf(" = input()")
	} else {
		ctx.cw.Printf(" = %s(input())", t)
	}
	ctx.cw.Println()
	return nil
}

func analyzeBlockScanOne(ctx *Context, n *ast.Block) {
	const (
		szero int = iota
		sscan
		seol
	)

	oz := scanOne{}

	var state = szero
	for _, n := range n.Statements {
		switch {
		case n.ScanStmt != nil:
			if state == szero {
				if len(n.ScanStmt.RefList) == 1 {
					state = sscan
					oz.scanStmt = n.ScanStmt
				}
			} else {
				state = szero
			}

		case n.CheckStmt != nil:

		case n.ForStmt != nil:
			state = szero

			analyzeBlock(ctx, &n.ForStmt.Block)

		case n.EOLStmt != nil:
			if state == sscan {
				state = seol
				oz.eolStmt = n.EOLStmt
				ctx.ozs[oz.scanStmt] = oz
				ctx.ozs[oz.eolStmt] = Noop{}

			} else {
				state = szero
			}

		default:
			state = szero
		}
	}
}
