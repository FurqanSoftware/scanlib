// Copyright 2020 Furqan Software Ltd. All rights reserved.

package cpp14

import (
	"bytes"
	"sort"

	"git.furqansoftware.net/toph/scanlib/ast"
	"git.furqansoftware.net/toph/scanlib/gen/code"
)

func Generate(n *ast.Source) ([]byte, error) {
	ctx := Context{
		types:    map[string]string{},
		includes: map[string]bool{},
		cw:       code.NewWriter("\t"),
	}
	ctx.includes["iostream"] = true

	ctx.cw.Indent(1)
	err := GenerateSource(&ctx, n)
	if err != nil {
		return nil, err
	}
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
		if t == "string" {
			ctx.includes["string"] = true
		}

		ctx.cw.Printf("%s", t)
		for i, x := range n.VarSpec.IdentList {
			ctx.types[x] = t

			if i > 0 {
				ctx.cw.Printf(",")
			}
			ctx.cw.Printf(" %s", x)
		}
		ctx.cw.Println(";")

	case n.VarSpec.Type.TypeLit != nil:
		t := ASTType[*n.VarSpec.Type.TypeLit.ArrayType.ElementType.TypeName]
		if t == "string" {
			ctx.includes["string"] = true
		}

		ctx.cw.Printf("%s", t)
		for i, x := range n.VarSpec.IdentList {
			ctx.types[x] = "array"
			ctx.types[x+"[]"] = t

			if i > 0 {
				ctx.cw.Printf(",")
			}
			ctx.cw.Printf(" %s[", x)
			err := GenerateExpression(ctx, &n.VarSpec.Type.TypeLit.ArrayType.ArrayLength)
			if err != nil {
				return err
			}
			ctx.cw.Print("]")
		}
		ctx.cw.Println(";")
	}
	return nil
}

func GenerateScanStmt(ctx *Context, n *ast.ScanStmt) error {
	for _, f := range n.RefList {
		ctx.cw.Printf("cin >> %s", f.Identifier)
		for _, i := range f.Indices {
			ctx.cw.Print("[")
			err := GenerateExpression(ctx, &i)
			if err != nil {
				return err
			}
			ctx.cw.Print("]")
		}
		ctx.cw.Print(";")
		ctx.cw.Println()
	}
	return nil
}

func GenerateCheckStmt(ctx *Context, n *ast.CheckStmt) error {
	return nil
}

func GenerateForStmt(ctx *Context, n *ast.ForStmt) error {
	ctx.cw.Printf("for (int %s = ", n.Range.Index)
	err := GenerateExpression(ctx, &n.Range.Low)
	if err != nil {
		return err
	}
	ctx.cw.Printf("; %s < ", n.Range.Index)
	err = GenerateExpression(ctx, &n.Range.High)
	if err != nil {
		return err
	}
	ctx.cw.Printf("; ++i) {")
	ctx.cw.Println()
	ctx.cw.Indent(1)
	GenerateBlock(ctx, &n.Block)
	ctx.cw.Indent(-1)
	ctx.cw.Printf("}")
	ctx.cw.Println()
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
	return GenerateRelative(ctx, n.Cmp)
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
	return GenerateAddition(ctx, n.Term)
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
