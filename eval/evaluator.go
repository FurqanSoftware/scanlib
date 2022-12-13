package eval

import (
	"errors"
	"fmt"
	"io"
	"reflect"

	"git.furqansoftware.net/toph/scanlib/ast"
	"github.com/alecthomas/participle/v2/lexer"
)

type evaluator struct {
	Source *ast.Source
	Input  *Input
	Values Values
}

func Evaluate(source *ast.Source, input io.Reader, options ...Option) (values Values, err error) {
	e := evaluator{
		Source: source,
		Values: Values{},
	}
	e.Input, err = newInput(input)
	if err != nil {
		return nil, err
	}

	for _, o := range options {
		o.apply(&e)
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

	ast.Walk(&e, e.Source)
	return e.Values, nil
}

func (e *evaluator) Visit(n ast.Node) (w ast.Visitor) {
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

func (e *evaluator) varDecl(n *ast.VarDecl) error {
	for _, x := range n.VarSpec.IdentList {
		switch {
		case n.VarSpec.Type.TypeName != nil:
			e.Values[x] = reflect.New(Types[*n.VarSpec.Type.TypeName])

		case n.VarSpec.Type.TypeLit != nil:
			l, err := e.expr(&n.VarSpec.Type.TypeLit.ArrayType.ArrayLength)
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
			e.Values[x] = v
		}
	}
	return nil
}

func (e *evaluator) scanStmt(n *ast.ScanStmt) error {
	for _, f := range n.RefList {
		v, ok := e.Values[f.Ident]
		if !ok {
			return ErrUndefined{Pos: n.Pos, Name: f.Ident}
		}
		for _, i := range f.Indices {
			r, err := e.expr(&i)
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
			d, err = e.Input.readBool()
			v.SetBool(d)
		case Types["int"]:
			var d int
			d, err = e.Input.readInt()
			v.SetInt(int64(d))
		case Types["int64"]:
			var d int64
			d, err = e.Input.readInt64()
			v.SetInt(d)
		case Types["float32"]:
			var d float32
			d, err = e.Input.readFloat32()
			v.SetFloat(float64(d))
		case Types["float64"]:
			var d float64
			d, err = e.Input.readFloat64()
			v.SetFloat(d)
		case Types["string"]:
			var d string
			d, err = e.Input.readString()
			v.SetString(d)
		default:
			return ErrCantScanType{}
		}
		if err != nil {
			if err == io.EOF {
				return ErrUnexpectedEOF{Pos: n.Pos}
			}
			return e.enrichError(err, n.Pos)
		}
	}
	return nil
}

func (e *evaluator) scanlnStmt(n *ast.ScanlnStmt) error {
	for _, f := range n.RefList {
		v, ok := e.Values[f.Ident]
		if !ok {
			return ErrUndefined{Pos: n.Pos, Name: f.Ident}
		}
		for _, i := range f.Indices {
			r, err := e.expr(&i)
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
			d, err = e.Input.readStringLn()
			v.SetString(d)
		default:
			return ErrCantScanType{}
		}
		if err != nil {
			if err == io.EOF {
				return ErrUnexpectedEOF{Pos: n.Pos}
			}
			return e.enrichError(err, n.Pos)
		}
	}
	return nil
}

func (e *evaluator) checkStmt(n *ast.CheckStmt) error {
	for i, x := range n.ExprList {
		v, err := e.expr(&x)
		if err != nil {
			return err
		}
		vb, _ := v.(bool)
		if !vb {
			return ErrCheckError{Pos: n.Pos, Clause: i + 1, Cursor: e.Input.Cursor}
		}
	}
	return nil
}

func (e *evaluator) ifStmt(n *ast.IfStmt) error {
	for _, n := range n.Branches {
		cond := false
		if n.Condition == nil {
			cond = true
		} else {
			v, err := e.expr(n.Condition)
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

func (e *evaluator) forStmt(n *ast.ForStmt) error {
	switch {
	case n.Range != nil:
		return e.forStmtRange(n)
	case n.Scan != nil:
		return e.forStmtScan(n)
	case n.Scanln != nil:
		return e.forStmtlnScan(n)
	}
	panic("unreachable")
}

func (e *evaluator) forStmtRange(n *ast.ForStmt) error {
	l, err := e.expr(&n.Range.Low)
	if err != nil {
		return err
	}
	li, ok := toInt(l)
	if !ok {
		return errors.New("invalid loop bound")
	}
	h, err := e.expr(&n.Range.High)
	if err != nil {
		return err
	}
	hi, ok := toInt(h)
	if !ok {
		return errors.New("invalid loop bound")
	}
	for i := li; i < hi; i++ {
		e.Values[n.Range.Index] = reflect.ValueOf(&i)
		ast.Walk(e, &n.Block)
	}
	return nil
}

func (e *evaluator) forStmtScan(n *ast.ForStmt) error {
	for {
		err := e.scanStmt(n.Scan)
		if err != nil {
			if errors.As(err, &ErrUnexpectedEOF{}) {
				break
			}
			return err
		}
		ast.Walk(e, &n.Block)
	}
	return nil
}

func (e *evaluator) forStmtlnScan(n *ast.ForStmt) error {
	for {
		err := e.scanlnStmt(n.Scanln)
		if err != nil {
			if errors.As(err, &ErrUnexpectedEOF{}) {
				break
			}
			return err
		}
		ast.Walk(e, &n.Block)
	}
	return nil
}

func (e *evaluator) eolStmt(n *ast.EOLStmt) error {
	eol, err := e.Input.isAtEOL()
	if err != nil {
		return err
	}
	if !eol {
		return ErrExpectedEOL{Pos: n.Pos, Got: e.Input.Scanner.Bytes(), Cursor: e.Input.Cursor}
	}
	return nil
}

func (e *evaluator) eofStmt(n *ast.EOFStmt) error {
	eof, err := e.Input.isAtEOF()
	if err != nil {
		return err
	}
	if !eof {
		return ErrExpectedEOF{Pos: n.Pos, Got: e.Input.Scanner.Bytes()}
	}
	return nil
}

func (e *evaluator) expr(n *ast.Expr) (interface{}, error) {
	l, err := e.logicalOr(n.Left)
	if err != nil {
		return nil, err
	}
	for _, c := range n.Right {
		r, err := e.opLogicalOr(c, l)
		if err != nil {
			return nil, err
		}
		l = r
	}
	return l, nil
}

func (e *evaluator) logicalOr(n *ast.LogicalOr) (interface{}, error) {
	l, err := e.logicalAnd(n.Left)
	if err != nil {
		return nil, err
	}
	for _, c := range n.Right {
		r, err := e.opLogicalAnd(c, l)
		if err != nil {
			return nil, err
		}
		l = r
	}
	return l, nil
}

func (e *evaluator) opLogicalOr(n *ast.OpLogicalOr, l interface{}) (interface{}, error) {
	r, err := e.logicalOr(n.LogicalOr)
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

func (e *evaluator) logicalAnd(n *ast.LogicalAnd) (interface{}, error) {
	l, err := e.relative(n.Left)
	if err != nil {
		return nil, err
	}
	for _, c := range n.Right {
		r, err := e.opRelative(c, l)
		if err != nil {
			return nil, err
		}
		l = r
	}
	return l, nil
}

func (e *evaluator) opLogicalAnd(n *ast.OpLogicalAnd, l interface{}) (interface{}, error) {
	r, err := e.logicalAnd(n.LogicalAnd)
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

func (e *evaluator) relative(n *ast.Relative) (interface{}, error) {
	l, err := e.addition(n.Left)
	if err != nil {
		return nil, err
	}
	for _, c := range n.Right {
		r, err := e.opAddition(c, l)
		if err != nil {
			return nil, err
		}
		l = r
	}
	return l, nil
}

func (e *evaluator) opRelative(n *ast.OpRelative, l interface{}) (interface{}, error) {
	r, err := e.relative(n.Relative)
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

func (e *evaluator) addition(n *ast.Addition) (interface{}, error) {
	l, err := e.multiplication(n.Left)
	if err != nil {
		return nil, err
	}
	for _, c := range n.Right {
		r, err := e.opMultiplication(c, l)
		if err != nil {
			return nil, err
		}
		l = r
	}
	return l, nil
}

func (e *evaluator) opAddition(n *ast.OpAddition, l interface{}) (interface{}, error) {
	r, err := e.addition(n.Addition)
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

func (e *evaluator) multiplication(n *ast.Multiplication) (interface{}, error) {
	return e.unary(n.Unary)
}

func (e *evaluator) opMultiplication(n *ast.OpMultiplication, l interface{}) (interface{}, error) {
	return l, nil
}

func (e *evaluator) unary(n *ast.Unary) (interface{}, error) {
	switch {
	case n.Value != nil:
		return e.primary(n.Value)

	case n.Negated != nil:
		v, err := e.primary(n.Negated)
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

func (e *evaluator) primary(n *ast.Primary) (interface{}, error) {
	switch {
	case n.CallExpr != nil:
		args := []interface{}{}
		for _, a := range n.CallExpr.Args {
			v, err := e.expr(&a)
			if err != nil {
				return nil, err
			}
			args = append(args, v)
		}
		return Functions[n.CallExpr.Ident](args...)

	case n.Variable != nil:
		v, ok := e.Values[n.Variable.Ident]
		if !ok {
			return nil, ErrUndefined{Pos: n.Pos, Name: n.Variable.Ident}
		}
		for _, i := range n.Variable.Indices {
			r, err := e.expr(&i)
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
		return e.basicLit(n.BasicLit)

	case n.SubExpr != nil:
		return e.expr(n.SubExpr)
	}

	panic("unreachable")
}

func (e *evaluator) basicLit(n *ast.BasicLit) (interface{}, error) {
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

func (e *evaluator) enrichError(err error, pos lexer.Position) error {
	switch err := err.(type) {
	case ErrBadParse:
		err.Pos = pos
	}
	return err
}
