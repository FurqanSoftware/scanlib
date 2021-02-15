package eval

import (
	"errors"
	"io"

	"git.furqansoftware.net/toph/scanlib/ast"
)

func Evaluate(ctx *Context, n *ast.Source) (interface{}, error) {
	return EvaluateSource(ctx, n)
}

func EvaluateSource(ctx *Context, n *ast.Source) (interface{}, error) {
	return EvaluateBlock(ctx, &n.Block)
}

func EvaluateBlock(ctx *Context, n *ast.Block) (interface{}, error) {
	for _, s := range n.Statement {
		_, err := EvaluateStatement(ctx, s)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func EvaluateStatement(ctx *Context, n *ast.Statement) (interface{}, error) {
	switch {
	case n.VarDecl != nil:
		return EvaluateVarDecl(ctx, n.VarDecl)

	case n.ScanStmt != nil:
		return EvaluateScanStmt(ctx, n.ScanStmt)

	case n.CheckStmt != nil:
		return EvaluateCheckStmt(ctx, n.CheckStmt)

	case n.ForStmt != nil:
		return EvaluateForStmt(ctx, n.ForStmt)

	case n.EOLStmt != nil:
		eol, err := ctx.Input.EOL()
		if err != nil {
			return nil, err
		}
		if !eol {
			return nil, ErrExpectedEOL{Pos: Cursor{n.Pos.Line, n.Pos.Column}}
		}
		return nil, nil

	case n.EOFStmt != nil:
		eof, err := ctx.Input.EOF()
		if err != nil {
			return nil, err
		}
		if !eof {
			return nil, ErrExpectedEOF{Pos: Cursor{n.Pos.Line, n.Pos.Column}, Token: ctx.Input.Scanner.Bytes()}
		}
		return nil, nil
	}
	panic("unreachable")
}

func EvaluateVarDecl(ctx *Context, n *ast.VarDecl) (interface{}, error) {
	for _, x := range n.VarSpec.IdentList {
		switch {
		case n.VarSpec.Type.TypeName != nil:
			ctx.Values[x] = Zero(ASTType[*n.VarSpec.Type.TypeName])

		case n.VarSpec.Type.TypeLit != nil:
			l, err := EvaluateExpression(ctx, &n.VarSpec.Type.TypeLit.ArrayType.ArrayLength)
			if err != nil {
				return nil, err
			}
			li, ok := l.(int)
			if !ok {
				return nil, errors.New("invalid array bound")
			}
			ctx.Values[x] = MakeArray(ASTType[*n.VarSpec.Type.TypeLit.ArrayType.ElementType.TypeName], []int{li})
		}
	}
	return nil, nil
}

func EvaluateScanStmt(ctx *Context, n *ast.ScanStmt) (interface{}, error) {
	for _, f := range n.RefList {
		indices := []int{}
		for _, i := range f.Indices {
			v, err := EvaluateExpression(ctx, &i)
			if err != nil {
				return nil, err
			}
			vi, ok := v.(int)
			if !ok {
				return nil, ErrNonIntegerIndex{Pos: Cursor{i.Pos.Line, i.Pos.Column}}
			}
			indices = append(indices, vi)
		}
		v := ctx.GetValue(f.Identifier, indices)
		var err error
		switch v.Type {
		case Bool:
			v.Data, err = ctx.Input.Bool()
		case Int:
			v.Data, err = ctx.Input.Int()
		case Int64:
			v.Data, err = ctx.Input.Int64()
		case Float32:
			v.Data, err = ctx.Input.Float32()
		case Float64:
			v.Data, err = ctx.Input.Float64()
		case String:
			v.Data, err = ctx.Input.String()
		default:
			return nil, ErrCantScanType{}
		}
		if err != nil {
			if err == io.EOF {
				return nil, ErrUnexpectedEOF{Pos: Cursor{n.Pos.Line, n.Pos.Column}}
			}
			return nil, err
		}
		ctx.SetValue(f.Identifier, indices, v)
	}
	return nil, nil
}

func EvaluateCheckStmt(ctx *Context, n *ast.CheckStmt) (interface{}, error) {
	for i, e := range n.ExprList {
		v, err := EvaluateExpression(ctx, &e)
		if err != nil {
			return nil, err
		}
		vb, _ := v.(bool)
		if !vb {
			return nil, ErrCheckError{Pos: Cursor{n.Pos.Line, n.Pos.Column}, Clause: i + 1}
		}
	}
	return nil, nil
}

func EvaluateForStmt(ctx *Context, n *ast.ForStmt) (interface{}, error) {
	l, err := EvaluateExpression(ctx, &n.Range.Low)
	if err != nil {
		return nil, err
	}
	li, ok := toInt(l)
	if !ok {
		return nil, errors.New("invalid loop bound")
	}
	h, err := EvaluateExpression(ctx, &n.Range.High)
	if err != nil {
		return nil, err
	}
	hi, ok := toInt(h)
	if !ok {
		return nil, errors.New("invalid loop bound")
	}
	for i := li; i < hi; i++ {
		ctx.Values[n.Range.Index] = Value{Int, i}

		_, err := EvaluateBlock(ctx, &n.Block)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func EvaluateExpression(ctx *Context, n *ast.Expression) (interface{}, error) {
	l, err := EvaluateCmp(ctx, n.Left)
	if err != nil {
		return nil, err
	}
	for _, c := range n.Right {
		r, err := EvaluateOpCmp(ctx, c, l)
		if err != nil {
			return nil, err
		}
		l = r
	}
	return l, nil
}

func EvaluateCmp(ctx *Context, n *ast.Cmp) (interface{}, error) {
	l, err := EvaluateTerm(ctx, n.Left)
	if err != nil {
		return nil, err
	}
	for _, c := range n.Right {
		r, err := EvaluateOpTerm(ctx, c, l)
		if err != nil {
			return nil, err
		}
		l = r
	}
	return l, nil
}

func EvaluateOpCmp(ctx *Context, n *ast.OpCmp, l interface{}) (interface{}, error) {
	r, err := EvaluateCmp(ctx, n.Cmp)
	if err != nil {
		return nil, err
	}
	switch l := l.(type) {
	case int:
		ri, ok := toInt(r)
		if !ok {
			return nil, ErrInvalidOperation{Pos: Cursor{n.Pos.Line, n.Pos.Column}}
		}
		switch n.Operator {
		case "==":
			return l == ri, nil
		case "!=":
			return l != ri, nil
		case "<=":
			return l <= ri, nil
		case ">=":
			return l >= ri, nil
		case "<":
			return l < ri, nil
		case ">":
			return l > ri, nil
		}

	case int64:
		ri, ok := toInt64(r)
		if !ok {
			return nil, ErrInvalidOperation{}
		}
		switch n.Operator {
		case "==":
			return l == ri, nil
		case "!=":
			return l != ri, nil
		case "<=":
			return l <= ri, nil
		case ">=":
			return l >= ri, nil
		case "<":
			return l < ri, nil
		case ">":
			return l > ri, nil
		}

	case float32:
		ri, ok := toFloat32(r)
		if !ok {
			return nil, ErrInvalidOperation{}
		}
		switch n.Operator {
		case "==":
			return l == ri, nil
		case "!=":
			return l != ri, nil
		case "<=":
			return l <= ri, nil
		case ">=":
			return l >= ri, nil
		case "<":
			return l < ri, nil
		case ">":
			return l > ri, nil
		}

	case float64:
		ri, ok := toFloat64(r)
		if !ok {
			return nil, ErrInvalidOperation{}
		}
		switch n.Operator {
		case "==":
			return l == ri, nil
		case "!=":
			return l != ri, nil
		case "<=":
			return l <= ri, nil
		case ">=":
			return l >= ri, nil
		case "<":
			return l < ri, nil
		case ">":
			return l > ri, nil
		}
	}
	return nil, ErrInvalidOperation{}
}

func EvaluateTerm(ctx *Context, n *ast.Term) (interface{}, error) {
	l, err := EvaluateFactor(ctx, n.Left)
	if err != nil {
		return nil, err
	}
	for _, c := range n.Right {
		r, err := EvaluateOpFactor(ctx, c, l)
		if err != nil {
			return nil, err
		}
		l = r
	}
	return l, nil
}

func EvaluateOpTerm(ctx *Context, n *ast.OpTerm, l interface{}) (interface{}, error) {
	r, err := EvaluateTerm(ctx, n.Term)
	if err != nil {
		return nil, err
	}
	switch l := l.(type) {
	case int:
		ri, ok := toInt(r)
		if !ok {
			return nil, ErrInvalidOperation{Pos: Cursor{n.Pos.Line, n.Pos.Column}}
		}
		switch n.Operator {
		case "+":
			return l + ri, nil
		case "-":
			return l - ri, nil
		}

	case int64:
		ri, ok := toInt64(r)
		if !ok {
			return nil, ErrInvalidOperation{}
		}
		switch n.Operator {
		case "+":
			return l + ri, nil
		case "-":
			return l - ri, nil
		}

	case float32:
		ri, ok := toFloat32(r)
		if !ok {
			return nil, ErrInvalidOperation{}
		}
		switch n.Operator {
		case "+":
			return l + ri, nil
		case "-":
			return l - ri, nil
		}

	case float64:
		ri, ok := toFloat64(r)
		if !ok {
			return nil, ErrInvalidOperation{}
		}
		switch n.Operator {
		case "+":
			return l + ri, nil
		case "-":
			return l - ri, nil
		}
	}
	return nil, ErrInvalidOperation{}
}

func EvaluateFactor(ctx *Context, n *ast.Factor) (interface{}, error) {
	return EvaluateUnary(ctx, n.Unary)
}

func EvaluateOpFactor(ctx *Context, n *ast.OpFactor, l interface{}) (interface{}, error) {
	return l, nil
}

func EvaluateUnary(ctx *Context, n *ast.Unary) (interface{}, error) {
	switch {
	case n.Value != nil:
		return EvaluateValue(ctx, n.Value)

	case n.Negated != nil:
		v, err := EvaluateValue(ctx, n.Negated)
		if err != nil {
			return nil, err
		}
		switch v := v.(type) {
		case int:
			return -v, nil
		case int64:
			return -v, nil
		case float32:
			return -v, nil
		case float64:
			return -v, nil
		case ast.Number:
			vs := string(v)
			if vs[0] != '-' {
				vs = "-" + vs
			} else {
				vs = vs[1:]
			}
			return ast.Number(vs), nil
		default:
			return nil, ErrInvalidOperation{}
		}
	}
	panic("unreachable")
}

func EvaluateValue(ctx *Context, n *ast.Value) (interface{}, error) {
	switch {
	case n.Call != nil:
		args := []interface{}{}
		for _, a := range n.Call.Arguments {
			v, err := EvaluateExpression(ctx, &a)
			if err != nil {
				return nil, err
			}
			args = append(args, v)
		}
		return Functions[n.Call.Name](args...)

	case n.Variable != nil:
		indices := []int{}
		for _, i := range n.Variable.Indices {
			v, err := EvaluateExpression(ctx, &i)
			if err != nil {
				return nil, err
			}
			vi, ok := v.(int)
			if !ok {
				return nil, ErrNonIntegerIndex{Pos: Cursor{i.Pos.Line, i.Pos.Column}}
			}
			indices = append(indices, vi)
		}
		return ctx.GetValue(n.Variable.Identifier, indices).Data, nil
	case n.Number != nil:
		return *n.Number, nil
	case n.String != nil:
		return *n.String, nil
	case n.Subexpression != nil:
		return EvaluateExpression(ctx, n.Subexpression)
	}
	panic("unreachable")
}
