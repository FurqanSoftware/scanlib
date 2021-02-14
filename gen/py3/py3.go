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
	for _, s := range n.Statement {
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

			ctx.cw.Printf("%s = [%s] * ", x, ASTZero[t])
			err := GenerateExpression(ctx, &n.VarSpec.Type.TypeLit.ArrayType.ArrayLength)
			if err != nil {
				return err
			}
			ctx.cw.Println()
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
	ctx.cw.Println("if _ == None: _ = list(input().split())")
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
	err := GenerateCmp(ctx, n.Left)
	if err != nil {
		return err
	}
	for _, c := range n.Right {
		err := GenerateOpCmp(ctx, c)
		if err != nil {
			return err
		}
	}
	return nil
}

func GenerateCmp(ctx *Context, n *ast.Cmp) error {
	err := GenerateTerm(ctx, n.Left)
	if err != nil {
		return err
	}
	for _, c := range n.Right {
		err := GenerateOpTerm(ctx, c)
		if err != nil {
			return err
		}
	}
	return nil
}

func GenerateOpCmp(ctx *Context, n *ast.OpCmp) error {
	return nil
	// err := GenerateCmp(ctx, n.Cmp)
	// if err != nil {
	// 	return err
	// }
	// switch l := l.(type) {
	// case int:
	// 	r, ok := r.(int)
	// 	if !ok {
	// 		return nil, ErrInvalidOperation{}
	// 	}
	// 	switch n.Operator {
	// 	case "==":
	// 		return l == r, nil
	// 	case "!=":
	// 		return l != r, nil
	// 	case "<=":
	// 		return l <= r, nil
	// 	case ">=":
	// 		return l >= r, nil
	// 	case "<":
	// 		return l < r, nil
	// 	case ">":
	// 		return l > r, nil
	// 	}
	// }
	// panic("unreachable")
}

func GenerateTerm(ctx *Context, n *ast.Term) error {
	err := GenerateFactor(ctx, n.Left)
	if err != nil {
		return err
	}
	for _, c := range n.Right {
		err := GenerateOpFactor(ctx, c)
		if err != nil {
			return err
		}
	}
	return nil
}

func GenerateOpTerm(ctx *Context, n *ast.OpTerm) error {
	return nil
}

func GenerateFactor(ctx *Context, n *ast.Factor) error {
	return GenerateValue(ctx, n.Base)
}

func GenerateOpFactor(ctx *Context, n *ast.OpFactor) error {
	return nil
}

func GenerateValue(ctx *Context, n *ast.Value) error {
	switch {
	case n.Call != nil:
		// args := []interface{}{}
		// for _, a := range n.Call.Arguments {
		// 	v, err := GenerateExpression(ctx, &a)
		// 	if err != nil {
		// 		return err
		// 	}
		// 	args = append(args, v)
		// }
		// return Functions[n.Call.Name](args...)

	case n.Variable != nil:
		// indices := []int{}
		// for _, i := range n.Variable.Indices {
		// 	v, err := GenerateExpression(ctx, &i)
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
