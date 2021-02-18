// Copyright 2020 Furqan Software Ltd. All rights reserved.

package py3

import (
	"strings"

	"git.furqansoftware.net/toph/scanlib/ast"
)

type onlyToken struct {
	scanStmt *ast.ScanStmt
	eolStmt  *ast.EOLStmt
}

func (o onlyToken) Generate(ctx *Context) error {
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

func (a *analyzer) onlyToken(n *ast.Block) {
	const (
		zero State = iota
		scan
	)

	oz := onlyToken{}

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

			case *ast.ScanStmt:
				if len(n.RefList) == 1 {
					oz.scanStmt = n
					state++
				}
				return false

			default:
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
				a.ozs[oz.scanStmt] = oz
				a.ozs[oz.eolStmt] = Noop{}
				state = zero
				return false

			default:
				state = zero
				return false
			}
		}

		return false
	})
}
