// Copyright 2020 Furqan Software Ltd. All rights reserved.

package cpp14

import (
	"bytes"
	"fmt"
	"sort"

	"git.furqansoftware.net/toph/scanlib/ast"
	"git.furqansoftware.net/toph/scanlib/gen/code"
)

type Generator struct {
	ctx *Context
}

func Generate(n *ast.Source) ([]byte, error) {
	ctx := Context{
		types:    map[string]string{},
		includes: map[string]bool{},
		cw:       code.NewWriter("\t"),
	}
	ctx.includes["iostream"] = true

	g := Generator{
		ctx: &ctx,
	}

	ctx.cw.Indent(1)
	ast.Walk(&g, n)
	ctx.cw.Indent(-1)

	r := bytes.Buffer{}
	includes := []string{}
	for k := range ctx.includes {
		includes = append(includes, k)
	}
	sort.Strings(includes)
	for _, inc := range includes {
		r.WriteString("#include <" + inc + ">\n")
	}
	r.WriteString("\n")
	r.WriteString("using namespace std;\n")
	r.WriteString("\n")
	r.WriteString("int main() {\n")
	r.Write(ctx.cw.Bytes())
	r.WriteString("\t\n")
	r.WriteString("\treturn 0;\n")
	r.WriteString("}\n")

	return r.Bytes(), nil
}

func (g *Generator) Visit(n ast.Node) (w ast.Visitor) {
	if n == nil {
		return nil
	}

	switch n := n.(type) {
	case *ast.Source, *ast.Block, *ast.Statement:
		return g

	case *ast.CheckStmt, *ast.EOLStmt, *ast.EOFStmt:
		return nil

	case *ast.VarDecl:
		g.varDecl(n)
		return nil

	case *ast.ScanStmt:
		g.scanStmt(n)
		return nil

	case *ast.IfStmt:
		g.ifStmt(n)
		return nil

	case *ast.ForStmt:
		g.forStmt(n)
		return nil
	}

	panic(fmt.Errorf("unreachable, with %T", n))
}

func (g *Generator) varDecl(n *ast.VarDecl) error {
	switch {
	case n.VarSpec.Type.TypeName != nil:
		t := ASTType[*n.VarSpec.Type.TypeName]
		if t == "string" {
			g.ctx.includes["string"] = true
		}

		g.ctx.cw.Printf("%s", t)
		for i, x := range n.VarSpec.IdentList {
			g.ctx.types[x] = t

			if i > 0 {
				g.ctx.cw.Printf(",")
			}
			g.ctx.cw.Printf(" %s", x)
		}
		g.ctx.cw.Println(";")

	case n.VarSpec.Type.TypeLit != nil:
		t := ASTType[*n.VarSpec.Type.TypeLit.ArrayType.ElementType.TypeName]
		if t == "string" {
			g.ctx.includes["string"] = true
		}

		g.ctx.cw.Printf("%s", t)
		for i, x := range n.VarSpec.IdentList {
			g.ctx.types[x] = "array"
			g.ctx.types[x+"[]"] = t

			if i > 0 {
				g.ctx.cw.Printf(",")
			}
			g.ctx.cw.Printf(" %s[", x)
			err := genExpr(g.ctx, &n.VarSpec.Type.TypeLit.ArrayType.ArrayLength)
			if err != nil {
				return err
			}
			g.ctx.cw.Print("]")
		}
		g.ctx.cw.Println(";")
	}
	return nil
}

func (g *Generator) scanStmt(n *ast.ScanStmt) error {
	g.ctx.cw.Printf("cin")
	for _, f := range n.RefList {
		g.ctx.cw.Printf(" >> %s", f.Ident)
		for _, i := range f.Indices {
			g.ctx.cw.Print("[")
			err := genExpr(g.ctx, &i)
			if err != nil {
				return err
			}
			g.ctx.cw.Print("]")
		}
	}
	g.ctx.cw.Print(";")
	g.ctx.cw.Println()
	return nil
}

func (g *Generator) ifStmt(n *ast.IfStmt) error {
	for i, n := range n.Branches {
		if i > 0 {
			g.ctx.cw.Print(" else ")
		}
		if n.Condition != nil {
			g.ctx.cw.Print("if (")
			err := genExpr(g.ctx, n.Condition)
			if err != nil {
				return err
			}
			g.ctx.cw.Printf(") {")
		} else {
			g.ctx.cw.Printf("{")
		}
		g.ctx.cw.Println()
		g.ctx.cw.Indent(1)
		ast.Walk(g, &n.Block)
		g.ctx.cw.Indent(-1)
		g.ctx.cw.Printf("}")
	}
	g.ctx.cw.Println()
	return nil
}

func (g *Generator) forStmt(n *ast.ForStmt) error {
	g.ctx.cw.Printf("for (int %s = ", n.Range.Index)
	err := genExpr(g.ctx, &n.Range.Low)
	if err != nil {
		return err
	}
	g.ctx.cw.Printf("; %s < ", n.Range.Index)
	err = genExpr(g.ctx, &n.Range.High)
	if err != nil {
		return err
	}
	g.ctx.cw.Printf("; ++i) {")
	g.ctx.cw.Println()
	g.ctx.cw.Indent(1)
	ast.Walk(g, &n.Block)
	g.ctx.cw.Indent(-1)
	g.ctx.cw.Printf("}")
	g.ctx.cw.Println()
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
	switch {
	case n.Value != nil:
		return genPrimary(ctx, n.Value)

	case n.Negated != nil:
		ctx.cw.Print("-")
		return genPrimary(ctx, n.Negated)
	}
	panic("unreachable")
}

func genPrimary(ctx *Context, n *ast.Primary) error {
	switch {
	case n.CallExpr != nil:
		// args := []interface{}{}
		// for _, a := range n.Call.Arguments {
		// 	v, err := genExpr(ctx, &a)
		// 	if err != nil {
		// 		return err
		// 	}
		// 	args = append(args, v)
		// }
		// return Functions[n.Call.Name](args...)

	case n.Variable != nil:
		// indices := []int{}
		// for _, i := range n.Variable.Indices {
		// 	v, err := genExpr(ctx, &i)
		// 	if err != nil {
		// 		return err
		// 	}
		// 	vi, ok := v.(int)
		// 	if !ok {
		// 		return nil, ErrNonIntegerIndex{}
		// 	}
		// 	indices = append(indices, vi)
		// }
		// return ctx.GetValue(n.Variable.Identifier, indices).Data, nil
		ctx.cw.Printf(n.Variable.Ident)
		return nil

	case n.BasicLit != nil:
		return genBasicLit(ctx, n.BasicLit)

	case n.SubExpr != nil:
		return genExpr(ctx, n.SubExpr)
	}
	panic("unreachable")
}

func genBasicLit(ctx *Context, n *ast.BasicLit) error {
	switch {
	case n.FloatLit != nil:
		ctx.cw.Printf("%f", *n.FloatLit)
		return nil

	case n.IntLit != nil:
		ctx.cw.Printf("%d", *n.IntLit)
		return nil

	case n.StringLit != nil:
		ctx.cw.Printf("%q", *n.StringLit)
		return nil
	}

	panic("unreachable")
}
