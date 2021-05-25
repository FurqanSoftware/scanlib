package eval

import (
	"errors"
	"fmt"
	"io"
	"reflect"

	"git.furqansoftware.net/toph/scanlib/ast"
)

type Evaluator struct {
	ctx *Context
}

func Evaluate(ctx *Context, n *ast.Source) (e *Evaluator, err error) {
	e = &Evaluator{
		ctx: ctx,
	}
	defer func() {
		v := recover()
		if v == nil {
			return
		}
		switch v := v.(type) {
		case error:
			err = v
		default:
			panic(err)
		}
	}()
	ast.Walk(e, n)
	return e, nil
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
		err := e.varDecl(n)
		catch(err)
		return nil

	case *ast.ScanStmt:
		err := e.scanStmt(n)
		catch(err)
		return nil

	case *ast.ScanlnStmt:
		err := e.scanlnStmt(n)
		catch(err)
		return nil

	case *ast.CheckStmt:
		err := e.checkStmt(n)
		catch(err)
		return nil

	case *ast.IfStmt:
		err := e.ifStmt(n)
		catch(err)
		return nil

	case *ast.ForStmt:
		err := e.forStmt(n)
		catch(err)
		return nil

	case *ast.EOLStmt:
		err := e.eolStmt(n)
		catch(err)
		return nil

	case *ast.EOFStmt:
		err := e.eofStmt(n)
		catch(err)
		return nil
	}

	panic(fmt.Errorf("unreachable, with %T", n))
}

func (e *Evaluator) varDecl(n *ast.VarDecl) error {
	for _, x := range n.VarSpec.IdentList {
		switch {
		case n.VarSpec.Type.TypeName != nil:
			e.ctx.Values[x] = reflect.New(Types[*n.VarSpec.Type.TypeName])

		case n.VarSpec.Type.TypeLit != nil:
			l, err := evalExpr(e.ctx, &n.VarSpec.Type.TypeLit.ArrayType.ArrayLength)
			if err != nil {
				return err
			}
			li, ok := l.(int)
			if !ok {
				return errors.New("invalid array bound")
			}
			// XXX(hjr265): This work's for one dimensional arrays only.
			t := reflect.SliceOf(Types[*n.VarSpec.Type.TypeLit.ArrayType.ElementType.TypeName])
			v := reflect.MakeSlice(t, li, li)
			e.ctx.Values[x] = v
		}
	}
	return nil
}

func (e *Evaluator) scanStmt(n *ast.ScanStmt) error {
	for _, f := range n.RefList {
		v, ok := e.ctx.Values[f.Ident]
		if !ok {
			return ErrUndefined{f.Ident}
		}
		for _, i := range f.Indices {
			r, err := evalExpr(e.ctx, &i)
			if err != nil {
				return err
			}
			ri, ok := r.(int)
			if !ok {
				return ErrNonIntegerIndex{Pos: i.Pos}
			}
			l := v.Len()
			if ri < 0 || ri >= l {
				return ErrInvalidIndex{Pos: i.Pos, Index: ri, Length: l}
			}
			v = v.Index(ri)
		}
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		var err error
		switch v.Type() {
		case Types["bool"]:
			var d bool
			d, err = e.ctx.Input.Bool()
			v.SetBool(d)
		case Types["int"]:
			var d int
			d, err = e.ctx.Input.Int()
			v.SetInt(int64(d))
		case Types["int64"]:
			var d int64
			d, err = e.ctx.Input.Int64()
			v.SetInt(d)
		case Types["float32"]:
			var d float32
			d, err = e.ctx.Input.Float32()
			v.SetFloat(float64(d))
		case Types["float64"]:
			var d float64
			d, err = e.ctx.Input.Float64()
			v.SetFloat(d)
		case Types["string"]:
			var d string
			d, err = e.ctx.Input.String()
			v.SetString(d)
		default:
			return ErrCantScanType{}
		}
		if err != nil {
			if err == io.EOF {
				return ErrUnexpectedEOF{Pos: n.Pos}
			}
			return err
		}
	}
	return nil
}

func (e *Evaluator) scanlnStmt(n *ast.ScanlnStmt) error {
	for _, f := range n.RefList {
		v, ok := e.ctx.Values[f.Ident]
		if !ok {
			return ErrUndefined{f.Ident}
		}
		for _, i := range f.Indices {
			r, err := evalExpr(e.ctx, &i)
			if err != nil {
				return err
			}
			ri, ok := r.(int)
			if !ok {
				return ErrNonIntegerIndex{Pos: i.Pos}
			}
			l := v.Len()
			if ri < 0 || ri >= l {
				return ErrInvalidIndex{Pos: i.Pos, Index: ri, Length: l}
			}
			v = v.Index(ri)
		}
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		var err error
		switch v.Type() {
		case Types["string"]:
			var d string
			d, err = e.ctx.Input.StringLn()
			v.SetString(d)
		default:
			return ErrCantScanType{}
		}
		if err != nil {
			if err == io.EOF {
				return ErrUnexpectedEOF{Pos: n.Pos}
			}
			return err
		}
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
			return ErrCheckError{Pos: n.Pos, Clause: i + 1, Cursor: e.ctx.Input.Cursor}
		}
	}
	return nil
}

func (e *Evaluator) ifStmt(n *ast.IfStmt) error {
	for _, n := range n.Branches {
		cond := false
		if n.Condition == nil {
			cond = true
		} else {
			v, err := evalExpr(e.ctx, n.Condition)
			if err != nil {
				return err
			}
			vb, ok := toBool(v)
			if !ok {
				return errors.New("non-bool used as if condition")
			}
			cond = vb
		}
		if cond {
			ast.Walk(e, &n.Block)
			break
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
		e.ctx.Values[n.Range.Index] = reflect.ValueOf(&i)
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
		return ErrExpectedEOL{Pos: n.Pos}
	}
	return nil
}

func (e *Evaluator) eofStmt(n *ast.EOFStmt) error {
	eof, err := e.ctx.Input.EOF()
	if err != nil {
		return err
	}
	if !eof {
		return ErrExpectedEOF{Pos: n.Pos, Token: e.ctx.Input.Scanner.Bytes()}
	}
	return nil
}

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
			return nil, ErrInvalidOperation{Pos: n.Pos}
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
			return nil, ErrInvalidOperation{Pos: n.Pos}
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
			return nil, ErrInvalidOperation{Pos: n.Pos}
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
			return nil, ErrInvalidOperation{Pos: n.Pos}
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

	case string:
		ri, ok := toString(r)
		if !ok {
			return nil, ErrInvalidOperation{Pos: n.Pos}
		}
		switch n.Operator {
		case "==":
			return l == ri, nil
		case "!=":
			return l != ri, nil
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
			return nil, ErrInvalidOperation{Pos: n.Pos}
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
		default:
			return nil, ErrInvalidOperation{}
		}
	}
	panic("unreachable")
}

func evalPrimary(ctx *Context, n *ast.Primary) (interface{}, error) {
	switch {
	case n.CallExpr != nil:
		args := []interface{}{}
		for _, a := range n.CallExpr.Args {
			v, err := evalExpr(ctx, &a)
			if err != nil {
				return nil, err
			}
			args = append(args, v)
		}
		return Functions[n.CallExpr.Ident](args...)

	case n.Variable != nil:
		v := ctx.Values[n.Variable.Ident]
		for _, i := range n.Variable.Indices {
			r, err := evalExpr(ctx, &i)
			if err != nil {
				return nil, err
			}
			ri, ok := r.(int)
			if !ok {
				return nil, ErrNonIntegerIndex{Pos: i.Pos}
			}
			v = v.Index(ri)
		}
		if v.Kind() == reflect.Ptr {
			return v.Elem().Interface(), nil
		}
		return v.Interface(), nil

	case n.BasicLit != nil:
		return evalBasicLit(ctx, n.BasicLit)

	case n.SubExpr != nil:
		return evalExpr(ctx, n.SubExpr)
	}

	panic("unreachable")
}

func evalBasicLit(ctx *Context, n *ast.BasicLit) (interface{}, error) {
	switch {
	case n.FloatLit != nil:
		return *n.FloatLit, nil

	case n.IntLit != nil:
		return *n.IntLit, nil

	case n.StringLit != nil:
		return *n.StringLit, nil
	}

	panic("unreachable")
}
