// Copyright 2020 Furqan Software Ltd. All rights reserved.

package py3

import (
	"bytes"
	"fmt"
	"strings"

	"git.furqansoftware.net/toph/scanlib/ast"
	"git.furqansoftware.net/toph/scanlib/gen/code"
)

type Generator struct {
	ctx *Context
}

func Generate(n *ast.Source) ([]byte, error) {
	ctx := Context{
		types: map[string]string{},
		cw:    code.NewWriter("\t"),
		ozs:   map[interface{}]Optimization{},
	}

	analyzeSource(&ctx, n)

	g := Generator{
		ctx: &ctx,
	}

	ast.Walk(&g, n)

	r := bytes.Buffer{}
	if ctx.linevar {
		r.WriteString("_ = None\n")
	}
	r.Write(ctx.cw.Bytes())

	return r.Bytes(), nil
}

func (g *Generator) Visit(n ast.Node) (w ast.Visitor) {
	if n == nil {
		return nil
	}

	switch n := n.(type) {
	case *ast.Source:
		return g

	case *ast.Block:
		return g

	case *ast.Statement:
		return g

	case *ast.VarDecl:
		g.varDecl(n)
		return nil

	case *ast.ScanStmt:
		g.scanStmt(n)
		return nil

	case *ast.CheckStmt:
		return g

	case *ast.ForStmt:
		g.forStmt(n)
		return nil

	case *ast.EOLStmt:
		g.eolStmt(n)
		return nil

	case *ast.EOFStmt:
		return g
	}

	panic(fmt.Errorf("unreachable, with %T", n))
}

func (g *Generator) varDecl(n *ast.VarDecl) error {
	switch {
	case n.VarSpec.Type.TypeName != nil:
		t := ASTType[*n.VarSpec.Type.TypeName]
		for _, x := range n.VarSpec.IdentList {
			g.ctx.types[x] = t
		}

	case n.VarSpec.Type.TypeLit != nil:
		t := ASTType[*n.VarSpec.Type.TypeLit.ArrayType.ElementType.TypeName]

		for _, x := range n.VarSpec.IdentList {
			g.ctx.types[x] = "array"
			g.ctx.types[x+"[]"] = t

			oz, ok := g.ctx.ozs[n]
			if ok {
				return oz.Generate(g.ctx)
			} else {
				g.ctx.cw.Printf("%s = [%s] * ", x, ASTZero[t])
				err := genExpr(g.ctx, &n.VarSpec.Type.TypeLit.ArrayType.ArrayLength)
				if err != nil {
					return err
				}
				g.ctx.cw.Println()
			}
		}
	}
	return nil
}

func (g *Generator) scanStmt(n *ast.ScanStmt) error {
	oz, ok := g.ctx.ozs[n]
	if ok {
		return oz.Generate(g.ctx)
	}

	g.ctx.linevar = true
	g.ctx.cw.Println("if _ == None: _ = input().split()")
	for _, f := range n.RefList {
		g.ctx.cw.Printf("%s", f.Identifier)
		for _, i := range f.Indices {
			g.ctx.cw.Print("[")
			err := genExpr(g.ctx, &i)
			if err != nil {
				return err
			}
			g.ctx.cw.Print("]")
		}
		g.ctx.cw.Printf(" = %s(_.pop(0))", g.ctx.types[f.Identifier+strings.Repeat("[]", len(f.Indices))])
		g.ctx.cw.Println()
	}
	return nil
}

func (g *Generator) forStmt(n *ast.ForStmt) error {
	oz, ok := g.ctx.ozs[n]
	if ok {
		return oz.Generate(g.ctx)
	}

	g.ctx.cw.Printf("for %s in range(", n.Range.Index)
	err := genExpr(g.ctx, &n.Range.Low)
	if err != nil {
		return err
	}
	g.ctx.cw.Print(", ")
	err = genExpr(g.ctx, &n.Range.High)
	if err != nil {
		return err
	}
	g.ctx.cw.Printf("):")
	g.ctx.cw.Println()
	g.ctx.cw.Indent(1)
	ast.Walk(g, &n.Block)
	g.ctx.cw.Indent(-1)
	return nil
}

func (g *Generator) eolStmt(n *ast.EOLStmt) error {
	oz, ok := g.ctx.ozs[n]
	if ok {
		return oz.Generate(g.ctx)
	}
	g.ctx.cw.Println("_ = None")
	return nil
}

func genExpr(ctx *Context, n *ast.Expr) error {
	err := genLogicalOr(ctx, n.Left)
	if err != nil {
		return err
	}
	for _, c := range n.Right {
		err := genOpLogicalOr(ctx, c)
		if err != nil {
			return err
		}
	}
	return nil
}

func genLogicalOr(ctx *Context, n *ast.LogicalOr) error {
	err := genLogicalAnd(ctx, n.Left)
	if err != nil {
		return err
	}
	for _, c := range n.Right {
		err := genOpLogicalAnd(ctx, c)
		if err != nil {
			return err
		}
	}
	return nil
}

func genOpLogicalOr(ctx *Context, n *ast.OpLogicalOr) error {
	ctx.cw.Print("||")
	return genLogicalOr(ctx, n.LogicalOr)
}

func genLogicalAnd(ctx *Context, n *ast.LogicalAnd) error {
	err := genRelative(ctx, n.Left)
	if err != nil {
		return err
	}
	for _, c := range n.Right {
		err := genOpRelative(ctx, c)
		if err != nil {
			return err
		}
	}
	return nil
}

func genOpLogicalAnd(ctx *Context, n *ast.OpLogicalAnd) error {
	ctx.cw.Print("&&")
	return genLogicalAnd(ctx, n.LogicalAnd)
}

func genRelative(ctx *Context, n *ast.Relative) error {
	err := genAddition(ctx, n.Left)
	if err != nil {
		return err
	}
	for _, c := range n.Right {
		err := genOpAddition(ctx, c)
		if err != nil {
			return err
		}
	}
	return nil
}

func genOpRelative(ctx *Context, n *ast.OpRelative) error {
	ctx.cw.Print(string(n.Operator))
	return genRelative(ctx, n.Relative)
}

func genAddition(ctx *Context, n *ast.Addition) error {
	err := genMultiplication(ctx, n.Left)
	if err != nil {
		return err
	}
	for _, c := range n.Right {
		err := genOpMultiplication(ctx, c)
		if err != nil {
			return err
		}
	}
	return nil
}

func genOpAddition(ctx *Context, n *ast.OpAddition) error {
	ctx.cw.Print(string(n.Operator))
	return genAddition(ctx, n.Addition)
}

func genMultiplication(ctx *Context, n *ast.Multiplication) error {
	return genUnary(ctx, n.Unary)
}

func genOpMultiplication(ctx *Context, n *ast.OpMultiplication) error {
	ctx.cw.Print(string(n.Operator))
	return genMultiplication(ctx, n.Factor)
}

func genUnary(ctx *Context, n *ast.Unary) error {
	return genPrimary(ctx, n.Value)
}

func genPrimary(ctx *Context, n *ast.Primary) error {
	switch {
	case n.Call != nil:

	case n.Variable != nil:
		ctx.cw.Printf(n.Variable.Identifier)
		return nil

	case n.Number != nil:
		ctx.cw.Printf("%s", *n.Number)
		return nil

	case n.String != nil:
		ctx.cw.Printf("%q", *n.String)
		return nil

	case n.SubExpr != nil:
		return genExpr(ctx, n.SubExpr)
	}
	panic("unreachable")
}
