package eval

import (
	"errors"
	"fmt"
	"io"

	"git.furqansoftware.net/toph/scanlib/ast"
)

type Evaluator struct {
	ctx *Context
}

func Evaluate(ctx *Context, n *ast.Source) (*Evaluator, error) {
	e := Evaluator{
		ctx: ctx,
	}
	ast.Walk(&e, n)
	return &e, nil

}

func (e *Evaluator) Visit(n ast.Node) (w ast.Visitor) {
	if n == nil {
		return nil
	}

	switch n := n.(type) {
	case *ast.Source:
		return e

	case *ast.Block:
		return e

	case *ast.Statement:
		return e

	case *ast.VarDecl:
		e.varDecl(n)
		return nil

	case *ast.ScanStmt:
		e.scanStmt(n)
		return nil

	case *ast.CheckStmt:
		e.checkStmt(n)
		return nil

	case *ast.ForStmt:
		e.forStmt(n)
		return nil

	case *ast.EOLStmt:
		e.eolStmt(n)
		return nil

	case *ast.EOFStmt:
		e.eofStmt(n)
		return nil
	}

	panic(fmt.Errorf("unreachable, with %T", n))
}

func (e *Evaluator) varDecl(n *ast.VarDecl) error {
	for _, x := range n.VarSpec.IdentList {
		switch {
		case n.VarSpec.Type.TypeName != nil:
			e.ctx.Values[x] = Zero(ASTType[*n.VarSpec.Type.TypeName])

		case n.VarSpec.Type.TypeLit != nil:
			l, err := evalExpr(e.ctx, &n.VarSpec.Type.TypeLit.ArrayType.ArrayLength)
			if err != nil {
				return err
			}
			li, ok := l.(int)
			if !ok {
				return errors.New("invalid array bound")
			}
			e.ctx.Values[x] = MakeArray(ASTType[*n.VarSpec.Type.TypeLit.ArrayType.ElementType.TypeName], []int{li})
		}
	}
	return nil
}

func (e *Evaluator) scanStmt(n *ast.ScanStmt) error {
	for _, f := range n.RefList {
		indices := []int{}
		for _, i := range f.Indices {
			v, err := evalExpr(e.ctx, &i)
			if err != nil {
				return err
			}
			vi, ok := v.(int)
			if !ok {
				return ErrNonIntegerIndex{Pos: Cursor{i.Pos.Line, i.Pos.Column}}
			}
			indices = append(indices, vi)
		}
		v := e.ctx.GetValue(f.Identifier, indices)
		var err error
		switch v.Type {
		case Bool:
			v.Data, err = e.ctx.Input.Bool()
		case Int:
			v.Data, err = e.ctx.Input.Int()
		case Int64:
			v.Data, err = e.ctx.Input.Int64()
		case Float32:
			v.Data, err = e.ctx.Input.Float32()
		case Float64:
			v.Data, err = e.ctx.Input.Float64()
		case String:
			v.Data, err = e.ctx.Input.String()
		default:
			return ErrCantScanType{}
		}
		if err != nil {
			if err == io.EOF {
				return ErrUnexpectedEOF{Pos: Cursor{n.Pos.Line, n.Pos.Column}}
			}
			return err
		}
		e.ctx.SetValue(f.Identifier, indices, v)
	}
	return nil
}

func (e *Evaluator) checkStmt(n *ast.CheckStmt) error {
	for i, x := range n.ExprList {
		v, err := evalExpr(e.ctx, &x)
		if err != nil {
			return err
		}
		vb, _ := v.(bool)
		if !vb {
			return ErrCheckError{Pos: Cursor{n.Pos.Line, n.Pos.Column}, Clause: i + 1}
		}
	}
	return nil
}

func (e *Evaluator) forStmt(n *ast.ForStmt) error {
	l, err := evalExpr(e.ctx, &n.Range.Low)
	if err != nil {
		return err
	}
	li, ok := toInt(l)
	if !ok {
		return errors.New("invalid loop bound")
	}
	h, err := evalExpr(e.ctx, &n.Range.High)
	if err != nil {
		return err
	}
	hi, ok := toInt(h)
	if !ok {
		return errors.New("invalid loop bound")
	}
	for i := li; i < hi; i++ {
		e.ctx.Values[n.Range.Index] = Value{Int, i}
		ast.Walk(e, &n.Block)
	}
	return nil
}

func (e *Evaluator) eolStmt(n *ast.EOLStmt) error {
	eol, err := e.ctx.Input.EOL()
	if err != nil {
		return err
	}
	if !eol {
		return ErrExpectedEOL{Pos: Cursor{n.Pos.Line, n.Pos.Column}}
	}
	return nil
}

func (e *Evaluator) eofStmt(n *ast.EOFStmt) error {
	eof, err := e.ctx.Input.EOF()
	if err != nil {
		return err
	}
	if !eof {
		return ErrExpectedEOF{Pos: Cursor{n.Pos.Line, n.Pos.Column}, Token: e.ctx.Input.Scanner.Bytes()}
	}
	return nil
}

//

func evalExpr(ctx *Context, n *ast.Expr) (interface{}, error) {
	l, err := evalLogicalOr(ctx, n.Left)
	if err != nil {
		return nil, err
	}
	for _, c := range n.Right {
		r, err := evalOpLogicalOr(ctx, c, l)
		if err != nil {
			return nil, err
		}
		l = r
	}
	return l, nil
}

func evalLogicalOr(ctx *Context, n *ast.LogicalOr) (interface{}, error) {
	l, err := evalLogicalAnd(ctx, n.Left)
	if err != nil {
		return nil, err
	}
	for _, c := range n.Right {
		r, err := evalOpLogicalAnd(ctx, c, l)
		if err != nil {
			return nil, err
		}
		l = r
	}
	return l, nil
}

func evalOpLogicalOr(ctx *Context, n *ast.OpLogicalOr, l interface{}) (interface{}, error) {
	r, err := evalLogicalOr(ctx, n.LogicalOr)
	if err != nil {
		return nil, err
	}
	switch l := l.(type) {
	case bool:
		ri, ok := toBool(r)
		if !ok {
			return nil, ErrInvalidOperation{Pos: Cursor{n.Pos.Line, n.Pos.Column}}
		}
		return l || ri, nil
	}
	return nil, ErrInvalidOperation{}
}

func evalLogicalAnd(ctx *Context, n *ast.LogicalAnd) (interface{}, error) {
	l, err := evalRelative(ctx, n.Left)
	if err != nil {
		return nil, err
	}
	for _, c := range n.Right {
		r, err := evalOpRelative(ctx, c, l)
		if err != nil {
			return nil, err
		}
		l = r
	}
	return l, nil
}

func evalOpLogicalAnd(ctx *Context, n *ast.OpLogicalAnd, l interface{}) (interface{}, error) {
	r, err := evalLogicalAnd(ctx, n.LogicalAnd)
	if err != nil {
		return nil, err
	}
	switch l := l.(type) {
	case bool:
		ri, ok := toBool(r)
		if !ok {
			return nil, ErrInvalidOperation{Pos: Cursor{n.Pos.Line, n.Pos.Column}}
		}
		return l && ri, nil
	}
	return nil, ErrInvalidOperation{}
}

func evalRelative(ctx *Context, n *ast.Relative) (interface{}, error) {
	l, err := evalAddition(ctx, n.Left)
	if err != nil {
		return nil, err
	}
	for _, c := range n.Right {
		r, err := evalOpAddition(ctx, c, l)
		if err != nil {
			return nil, err
		}
		l = r
	}
	return l, nil
}

func evalOpRelative(ctx *Context, n *ast.OpRelative, l interface{}) (interface{}, error) {
	r, err := evalRelative(ctx, n.Relative)
	if err != nil {
		return nil, err
	}
	switch l := l.(type) {
	case bool:
		ri, ok := toBool(r)
		if !ok {
			return nil, ErrInvalidOperation{Pos: Cursor{n.Pos.Line, n.Pos.Column}}
		}
		switch n.Operator {
		case "==":
			return l == ri, nil
		case "!=":
			return l != ri, nil
		}

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

func evalAddition(ctx *Context, n *ast.Addition) (interface{}, error) {
	l, err := evalMultiplication(ctx, n.Left)
	if err != nil {
		return nil, err
	}
	for _, c := range n.Right {
		r, err := evalOpMultiplication(ctx, c, l)
		if err != nil {
			return nil, err
		}
		l = r
	}
	return l, nil
}

func evalOpAddition(ctx *Context, n *ast.OpAddition, l interface{}) (interface{}, error) {
	r, err := evalAddition(ctx, n.Addition)
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

func evalMultiplication(ctx *Context, n *ast.Multiplication) (interface{}, error) {
	return evalUnary(ctx, n.Unary)
}

func evalOpMultiplication(ctx *Context, n *ast.OpMultiplication, l interface{}) (interface{}, error) {
	return l, nil
}

func evalUnary(ctx *Context, n *ast.Unary) (interface{}, error) {
	switch {
	case n.Value != nil:
		return evalPrimary(ctx, n.Value)

	case n.Negated != nil:
		v, err := evalPrimary(ctx, n.Negated)
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

func evalPrimary(ctx *Context, n *ast.Primary) (interface{}, error) {
	switch {
	case n.Call != nil:
		args := []interface{}{}
		for _, a := range n.Call.Arguments {
			v, err := evalExpr(ctx, &a)
			if err != nil {
				return nil, err
			}
			args = append(args, v)
		}
		return Functions[n.Call.Name](args...)

	case n.Variable != nil:
		indices := []int{}
		for _, i := range n.Variable.Indices {
			v, err := evalExpr(ctx, &i)
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
	case n.SubExpr != nil:
		return evalExpr(ctx, n.SubExpr)
	}

	panic("unreachable")
}
