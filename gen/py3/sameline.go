// Copyright 2020 Furqan Software Ltd. All rights reserved.

package py3

import (
	"strings"

	"git.furqansoftware.net/toph/scanlib/ast"
)

type sameLine struct {
	scanStmt *ast.ScanStmt
}

func (o sameLine) Generate(ctx *Context) error {
	for _, f := range o.scanStmt.RefList {
		ctx.cw.Printf("%s", f.Ident)
		for _, i := range f.Indices {
			ctx.cw.Print("[")
			err := genExpr(ctx, &i)
			if err != nil {
				return err
			}
			ctx.cw.Print("]")
		}
		ctx.cw.Printf(" = %s(_.pop(0))", ctx.types[f.Ident+strings.Repeat("[]", len(f.Indices))])
		ctx.cw.Println()
	}
	return nil
}

func (a *analyzer) sameLine(n *ast.Block) {
	if a.blockEOLs[n] {
		return
	}

	ast.Inspect(n, func(n ast.Node) bool {
		switch n := n.(type) {
		case *ast.Block, *ast.Statement:
			return true

		case *ast.ScanStmt:
			a.ozs[n] = sameLine{
				scanStmt: n,
			}
			return false
		}

		return false
	})
}
