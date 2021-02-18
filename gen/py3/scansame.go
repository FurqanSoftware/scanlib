// Copyright 2020 Furqan Software Ltd. All rights reserved.

package py3

import "git.furqansoftware.net/toph/scanlib/ast"

type scanSame struct {
	varDecl  *ast.VarDecl
	scanStmt *ast.ScanStmt
	eolStmt  *string
}

func (o scanSame) Generate(ctx *Context) error {
	for i, x := range o.scanStmt.RefList {
		if i > 0 {
			ctx.cw.Print(", ")
		}
		ctx.cw.Printf("%s", x.Identifier)
	}
	t := ASTType[*o.varDecl.VarSpec.Type.TypeName]
	if len(o.scanStmt.RefList) == 1 {
		if t == "string" {
			ctx.cw.Printf(" = input()")
		} else {
			ctx.cw.Printf(" = %s(input())", t)
		}
	} else {
		ctx.cw.Printf(" = map(%s, input().split())", t)
	}
	ctx.cw.Println()
	return nil
}

func analyzeBlockScanSame(ctx *Context, n *ast.Block) {
	const (
		szero int = iota
		svar
		sscan
		seol
	)

	oz := scanSame{}

	var state = szero
	for _, n := range n.Statements {
		switch {
		case n.VarDecl != nil:
			if n.VarDecl.VarSpec.Type.TypeName != nil {
				state = svar
				oz.varDecl = n.VarDecl
			}

		case n.ScanStmt != nil:
			if state == svar {
				intersect := true
				idents := map[string]bool{}
				for _, x := range oz.varDecl.VarSpec.IdentList {
					idents[x] = true
				}
				for _, r := range n.ScanStmt.RefList {
					if !idents[r.Identifier] && len(r.Indices) > 0 {
						intersect = false
					}
				}
				if intersect {
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
