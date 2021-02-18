package walk

import (
	"fmt"
	"strings"

	"git.furqansoftware.net/toph/scanlib/ast"
)

func Walk(n *ast.Source) {
	WalkSource(n, 1)
}

func WalkSource(n *ast.Source, d int) {
	pad(d)
	fmt.Println("Source")
	WalkBlock(&n.Block, d+1)
}

func WalkBlock(n *ast.Block, d int) {
	pad(d)
	fmt.Println("Block")
	for _, s := range n.Statement {
		WalkStatement(s, d+1)
	}
}

func WalkStatement(n *ast.Statement, d int) {
	pad(d)
	fmt.Println("Statement")
	switch {
	case n.VarDecl != nil:
		WalkVarDecl(n.VarDecl, d+1)
	case n.ScanStmt != nil:
		WalkScanStmt(n.ScanStmt, d+1)
	case n.CheckStmt != nil:
		WalkCheckStmt(n.CheckStmt, d+1)
	case n.EOLStmt != nil:
		pad(d + 1)
		fmt.Println("EOLStmt")
	case n.EOFStmt != nil:
		pad(d + 1)
		fmt.Println("EOFStmt")
	}
}

func WalkVarDecl(n *ast.VarDecl, d int) {
	pad(d)
	fmt.Println("VarDecl")
	WalkVarSpec(&n.VarSpec, d+1)
}

func WalkScanStmt(n *ast.ScanStmt, d int) {
	pad(d)
	fmt.Println("ScanStmt")
}

func WalkCheckStmt(n *ast.CheckStmt, d int) {
	pad(d)
	fmt.Println("CheckStmt")
	for _, e := range n.ExprList {
		WalkExpression(&e, d+1)
	}
}

func WalkVarSpec(n *ast.VarSpec, d int) {
	pad(d)
	fmt.Println("VarSpec")
	WalkType(&n.Type, d+1)
	for _, t := range n.IdentList {
		pad(d + 1)
		fmt.Println(t)
	}
}

func WalkType(n *ast.Type, d int) {
	pad(d)
	fmt.Println("Type")
	switch {
	case n.TypeName != nil:
		pad(d + 1)
		fmt.Println(*n.TypeName)
	case n.TypeLit != nil:
		WalkTypeLit(n.TypeLit, d+1)
	}
}

func WalkTypeLit(n *ast.TypeLit, d int) {
	pad(d)
	fmt.Println("TypeLit")
	WalkArrayType(n.ArrayType, d+1)
}

func WalkArrayType(n *ast.ArrayType, d int) {
	pad(d)
	fmt.Println("ArrayType")
	WalkType(&n.ElementType, d+1)
}

func WalkExpression(n *ast.Expression, d int) {
	pad(d)
	fmt.Println("Expression")
	if n.Left != nil {
		WalkLogicalOr(n.Left, d+1)
	}
	for _, c := range n.Right {
		WalkOpLogicalOr(c, d+1)
	}
}

func WalkLogicalOr(n *ast.LogicalOr, d int) {
	pad(d)
	fmt.Println("LogicalOr")
	if n.Left != nil {
		WalkLogicalAnd(n.Left, d+1)
	}
	for _, t := range n.Right {
		WalkOpLogicalAnd(t, d+1)
	}
}

func WalkOpLogicalOr(n *ast.OpLogicalOr, d int) {
	pad(d)
	fmt.Println("OpLogicalOr")
	if n.LogicalOr != nil {
		WalkLogicalOr(n.LogicalOr, d+1)
	}
}

func WalkLogicalAnd(n *ast.LogicalAnd, d int) {
	pad(d)
	fmt.Println("LogicalAnd")
	if n.Left != nil {
		WalkRelative(n.Left, d+1)
	}
	for _, t := range n.Right {
		WalkOpRelative(t, d+1)
	}
}

func WalkOpLogicalAnd(n *ast.OpLogicalAnd, d int) {
	pad(d)
	fmt.Println("OpLogicalAnd")
	if n.LogicalAnd != nil {
		WalkLogicalAnd(n.LogicalAnd, d+1)
	}
}

func WalkRelative(n *ast.Relative, d int) {
	pad(d)
	fmt.Println("Relative")
	if n.Left != nil {
		WalkAddition(n.Left, d+1)
	}
	for _, t := range n.Right {
		WalkOpAddition(t, d+1)
	}
}

func WalkOpRelative(n *ast.OpRelative, d int) {
	pad(d)
	fmt.Println("OpRelative")
	WalkOperator(n.Operator, d+1)
	if n.Cmp != nil {
		WalkRelative(n.Cmp, d+1)
	}
}

func WalkAddition(n *ast.Addition, d int) {
	pad(d)
	fmt.Println("Addition")
	if n.Left != nil {
		WalkMultiplication(n.Left, d+1)
	}
	if n.Right != nil {
		WalkMultiplication(n.Left, d+1)
	}
}

func WalkOpAddition(n *ast.OpAddition, d int) {
	pad(d)
	fmt.Println("OpAddition")
	WalkOperator(n.Operator, d+1)
	if n.Term != nil {
		WalkAddition(n.Term, d+1)
	}
}

func WalkMultiplication(n *ast.Multiplication, d int) {
	pad(d)
	fmt.Println("Multiplication")
	if n.Unary != nil {
		WalkUnary(n.Unary, d+1)
	}
	if n.Exponent != nil {
		WalkPrimary(n.Exponent, d+1)
	}
}

func WalkOpMultiplication(n *ast.OpMultiplication, d int) {
	pad(d)
	fmt.Println("OpMultiplication")
	WalkOperator(n.Operator, d+1)
	if n.Factor != nil {
		WalkMultiplication(n.Factor, d+1)
	}
}

func WalkUnary(n *ast.Unary, d int) {
	pad(d)
	fmt.Println("Unary")
	WalkPrimary(n.Value, d+1)
}

func WalkPrimary(n *ast.Primary, d int) {
	pad(d)
	fmt.Println("Primary")
	switch {
	case n.Number != nil:
		pad(d + 1)
		fmt.Println("Number", *n.Number)
	case n.Variable != nil:
		pad(d + 1)
		fmt.Println("Variable", *n.Variable)
	case n.String != nil:
		pad(d + 1)
		fmt.Println("String", *n.String)
	case n.Call != nil:
		// TODO
	case n.Subexpression != nil:
		WalkExpression(n.Subexpression, d+1)
	}
}

func WalkOperator(n ast.Operator, d int) {
	pad(d)
	fmt.Println("Operator", n)
}

func pad(n int) {
	if n > 0 {
		fmt.Print(strings.Repeat("..", n), " ")
	}
}
