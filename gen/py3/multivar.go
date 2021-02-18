// Copyright 2020 Furqan Software Ltd. All rights reserved.

package py3

import "git.furqansoftware.net/toph/scanlib/ast"

type multiVar struct {
	varDecl  *ast.VarDecl
	scanStmt *ast.ScanStmt
	eolStmt  *ast.EOLStmt
}

func (o multiVar) Generate(ctx *Context) error {
	for i, x := range o.scanStmt.RefList {
		if i > 0 {
			ctx.cw.Print(", ")
		}
		ctx.cw.Printf("%s", x.Ident)
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

func (a *analyzer) multiVar(n *ast.Block) {
	const (
		zero State = iota
		decl
		scan
	)

	oz := multiVar{}

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
				if n.VarSpec.Type.TypeName != nil {
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

			case *ast.ScanStmt:
				intersect := true
				idents := map[string]bool{}
				for _, x := range oz.varDecl.VarSpec.IdentList {
					idents[x] = true
				}
				for _, r := range n.RefList {
					if !idents[r.Ident] && len(r.Indices) > 0 {
						intersect = false
					}
				}
				if intersect {
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
				a.ozs[oz.scanStmt] = oz
				a.ozs[oz.varDecl] = Noop{}
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
