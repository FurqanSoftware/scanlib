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
		WalkCmp(n.Left, d+1)
	}
	for _, c := range n.Right {
		WalkOpCmp(c, d+1)
	}
}

func WalkCmp(n *ast.Cmp, d int) {
	pad(d)
	fmt.Println("Cmp")
	if n.Left != nil {
		WalkTerm(n.Left, d+1)
	}
	for _, t := range n.Right {
		WalkOpTerm(t, d+1)
	}
}

func WalkOpCmp(n *ast.OpCmp, d int) {
	pad(d)
	fmt.Println("OpCmp")
	WalkOperator(n.Operator, d+1)
	if n.Cmp != nil {
		WalkCmp(n.Cmp, d+1)
	}
}

func WalkTerm(n *ast.Term, d int) {
	pad(d)
	fmt.Println("Term")
	if n.Left != nil {
		WalkFactor(n.Left, d+1)
	}
	if n.Right != nil {
		WalkFactor(n.Left, d+1)
	}
}

func WalkOpTerm(n *ast.OpTerm, d int) {
	pad(d)
	fmt.Println("OpTerm")
	WalkOperator(n.Operator, d+1)
	if n.Term != nil {
		WalkTerm(n.Term, d+1)
	}
}

func WalkFactor(n *ast.Factor, d int) {
	pad(d)
	fmt.Println("Factor")
	if n.Base != nil {
		WalkValue(n.Base, d+1)
	}
	if n.Exponent != nil {
		WalkValue(n.Exponent, d+1)
	}
}

func WalkOpFactor(n *ast.OpFactor, d int) {
	pad(d)
	fmt.Println("OpFactor")
	WalkOperator(n.Operator, d+1)
	if n.Factor != nil {
		WalkFactor(n.Factor, d+1)
	}
}

func WalkValue(n *ast.Value, d int) {
	pad(d)
	fmt.Println("Value")
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

// func WalkValue(n *ast.Value, d int) {
// 	pad(d)
// 	fmt.Println("Value")
// 	switch {
// 	case n.Make != nil:
// 		WalkMake(n.Make, d+1)
// 	case n.Call != nil:
// 		WalkCall(n.Call, d+1)
// 	case n.Variable != nil:
// 		pad(d + 1)
// 		fmt.Println("Variable", *n.Variable)
// 	case n.Integer != nil:
// 		pad(d + 1)
// 		fmt.Println("Integer", *n.Integer)
// 	case n.String != nil:
// 		pad(d + 1)
// 		fmt.Println("String", *n.String)
// 	}
// }

// func WalkMake(n *ast.Make, d int) {
// 	pad(d)
// 	fmt.Println("Make")
// 	WalkArray(&n.Array, d+1)
// }

// func WalkArray(n *ast.Array, d int) {
// 	pad(d)
// 	fmt.Println("Array", n.Variable)
// 	WalkIndex(&n.Index, d+1)
// }

// func WalkCall(n *ast.Call, d int) {
// 	pad(d)
// 	fmt.Println("Call")
// 	pad(d + 1)
// 	fmt.Println(n.Name)
// 	for _, a := range n.Arguments {
// 		WalkValue(a, d+1)
// 	}
// }

// func WalkIndex(n *ast.Index, d int) {
// 	pad(d)
// 	fmt.Println("Index")
// 	switch {
// 	case n.Variable != nil:
// 		pad(d + 1)
// 		fmt.Println("Variable", *n.Variable)
// 	case n.Integer != nil:
// 		pad(d + 1)
// 		fmt.Println("Integer", *n.Integer)
// 	}
// }

// func WalkFor(n *ast.For, d int) {
// 	WalkIndex(&n.From, d+1)
// 	WalkIndex(&n.To, d+1)
// 	WalkBlock(&n.Block, d+1)
// }

func pad(n int) {
	if n > 0 {
		fmt.Print(strings.Repeat("..", n), " ")
	}
}
