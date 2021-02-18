// Copyright 2020 Furqan Software Ltd. All rights reserved.

package py3

import (
	"bytes"
	"strings"

	"git.furqansoftware.net/toph/scanlib/ast"
	"git.furqansoftware.net/toph/scanlib/gen/code"
)

func Generate(n *ast.Source) ([]byte, error) {
	ctx := Context{
		types: map[string]string{},
		cw:    code.NewWriter("\t"),
		ozs:   map[interface{}]Optimization{},
	}

	analyzeSource(&ctx, n)

	err := GenerateSource(&ctx, n)
	if err != nil {
		return nil, err
	}

	r := bytes.Buffer{}
	if ctx.linevar {
		r.WriteString("_ = None\n")
	}
	r.Write(ctx.cw.Bytes())

	return r.Bytes(), nil
}

func GenerateSource(ctx *Context, n *ast.Source) error {
	return GenerateBlock(ctx, &n.Block)
}

func GenerateBlock(ctx *Context, n *ast.Block) error {
	for _, s := range n.Statements {
		err := GenerateStatement(ctx, s)
		if err != nil {
			return err
		}
	}
	return nil
}

func GenerateStatement(ctx *Context, n *ast.Statement) error {
	switch {
	case n.VarDecl != nil:
		return GenerateVarDecl(ctx, n.VarDecl)

	case n.ScanStmt != nil:
		return GenerateScanStmt(ctx, n.ScanStmt)

	case n.CheckStmt != nil:
		return GenerateCheckStmt(ctx, n.CheckStmt)

	case n.ForStmt != nil:
		return GenerateForStmt(ctx, n.ForStmt)

	case n.EOLStmt != nil:
		oz, ok := ctx.ozs[n.EOLStmt]
		if ok {
			return oz.Generate(ctx)
		}

		ctx.cw.Println("_ = None")
		return nil

	case n.EOFStmt != nil:
		return nil
	}
	panic("unreachable")
}

func GenerateVarDecl(ctx *Context, n *ast.VarDecl) error {
	switch {
	case n.VarSpec.Type.TypeName != nil:
		t := ASTType[*n.VarSpec.Type.TypeName]
		for _, x := range n.VarSpec.IdentList {
			ctx.types[x] = t
		}

	case n.VarSpec.Type.TypeLit != nil:
		t := ASTType[*n.VarSpec.Type.TypeLit.ArrayType.ElementType.TypeName]

		for _, x := range n.VarSpec.IdentList {
			ctx.types[x] = "array"
			ctx.types[x+"[]"] = t

			oz, ok := ctx.ozs[n]
			if ok {
				return oz.Generate(ctx)
			} else {
				ctx.cw.Printf("%s = [%s] * ", x, ASTZero[t])
				err := GenerateExpression(ctx, &n.VarSpec.Type.TypeLit.ArrayType.ArrayLength)
				if err != nil {
					return err
				}
				ctx.cw.Println()
			}
		}
	}

	return nil
}

func GenerateScanStmt(ctx *Context, n *ast.ScanStmt) error {
	oz, ok := ctx.ozs[n]
	if ok {
		return oz.Generate(ctx)
	}

	ctx.linevar = true
	ctx.cw.Println("if _ == None: _ = input().split()")
	for _, f := range n.RefList {
		ctx.cw.Printf("%s", f.Identifier)
		for _, i := range f.Indices {
			ctx.cw.Print("[")
			err := GenerateExpression(ctx, &i)
			if err != nil {
				return err
			}
			ctx.cw.Print("]")
		}
		ctx.cw.Printf(" = %s(_.pop(0))", ctx.types[f.Identifier+strings.Repeat("[]", len(f.Indices))])
		ctx.cw.Println()
	}
	return nil
}

func GenerateCheckStmt(ctx *Context, n *ast.CheckStmt) error {
	return nil
}

func GenerateForStmt(ctx *Context, n *ast.ForStmt) error {
	oz, ok := ctx.ozs[n]
	if ok {
		return oz.Generate(ctx)
	}

	ctx.cw.Printf("for %s in range(", n.Range.Index)
	err := GenerateExpression(ctx, &n.Range.Low)
	if err != nil {
		return err
	}
	ctx.cw.Print(", ")
	err = GenerateExpression(ctx, &n.Range.High)
	if err != nil {
		return err
	}
	ctx.cw.Printf("):")
	ctx.cw.Println()
	ctx.cw.Indent(1)
	GenerateBlock(ctx, &n.Block)
	ctx.cw.Indent(-1)
	return nil
}

func GenerateExpression(ctx *Context, n *ast.Expression) error {
	err := GenerateLogicalOr(ctx, n.Left)
	if err != nil {
		return err
	}
	for _, c := range n.Right {
		err := GenerateOpLogicalOr(ctx, c)
		if err != nil {
			return err
		}
	}
	return nil
}

func GenerateLogicalOr(ctx *Context, n *ast.LogicalOr) error {
	err := GenerateLogicalAnd(ctx, n.Left)
	if err != nil {
		return err
	}
	for _, c := range n.Right {
		err := GenerateOpLogicalAnd(ctx, c)
		if err != nil {
			return err
		}
	}
	return nil
}

func GenerateOpLogicalOr(ctx *Context, n *ast.OpLogicalOr) error {
	ctx.cw.Print("||")
	return GenerateLogicalOr(ctx, n.LogicalOr)
}

func GenerateLogicalAnd(ctx *Context, n *ast.LogicalAnd) error {
	err := GenerateRelative(ctx, n.Left)
	if err != nil {
		return err
	}
	for _, c := range n.Right {
		err := GenerateOpRelative(ctx, c)
		if err != nil {
			return err
		}
	}
	return nil
}

func GenerateOpLogicalAnd(ctx *Context, n *ast.OpLogicalAnd) error {
	ctx.cw.Print("&&")
	return GenerateLogicalAnd(ctx, n.LogicalAnd)
}

func GenerateRelative(ctx *Context, n *ast.Relative) error {
	err := GenerateAddition(ctx, n.Left)
	if err != nil {
		return err
	}
	for _, c := range n.Right {
		err := GenerateOpAddition(ctx, c)
		if err != nil {
			return err
		}
	}
	return nil
}

func GenerateOpRelative(ctx *Context, n *ast.OpRelative) error {
	ctx.cw.Print(string(n.Operator))
	return GenerateRelative(ctx, n.Relative)
}

func GenerateAddition(ctx *Context, n *ast.Addition) error {
	err := GenerateMultiplication(ctx, n.Left)
	if err != nil {
		return err
	}
	for _, c := range n.Right {
		err := GenerateOpMultiplication(ctx, c)
		if err != nil {
			return err
		}
	}
	return nil
}

func GenerateOpAddition(ctx *Context, n *ast.OpAddition) error {
	ctx.cw.Print(string(n.Operator))
	return GenerateAddition(ctx, n.Addition)
}

func GenerateMultiplication(ctx *Context, n *ast.Multiplication) error {
	return GenerateUnary(ctx, n.Unary)
}

func GenerateOpMultiplication(ctx *Context, n *ast.OpMultiplication) error {
	ctx.cw.Print(string(n.Operator))
	return GenerateMultiplication(ctx, n.Factor)
}

func GenerateUnary(ctx *Context, n *ast.Unary) error {
	return GeneratePrimary(ctx, n.Value)
}

func GeneratePrimary(ctx *Context, n *ast.Primary) error {
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

	case n.Subexpression != nil:
		return GenerateExpression(ctx, n.Subexpression)
	}
	panic("unreachable")
}
