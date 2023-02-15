package ast

import "fmt"

type Visitor interface {
    Visit(Node) Visitor
}

func Walk(v Visitor, n Node) {
    v = v.Visit(n)
    if v == nil {
        return
    }

    switch n := n.(type) {
    case *Source:
        Walk(v, &n.Block)

    case *Block:
        for _, n := range n.Statements {
            Walk(v, n)
        }

    case *Statement:
        switch {
        case n.VarDecl != nil:
            Walk(v, n.VarDecl)
        case n.ScanStmt != nil:
            Walk(v, n.ScanStmt)
        case n.ScanlnStmt != nil:
            Walk(v, n.ScanlnStmt)
        case n.CheckStmt != nil:
            Walk(v, n.CheckStmt)
        case n.IfStmt != nil:
            Walk(v, n.IfStmt)
        case n.ForStmt != nil:
            Walk(v, n.ForStmt)
        case n.EOLStmt != nil:
            Walk(v, n.EOLStmt)
        case n.EOFStmt != nil:
            Walk(v, n.EOFStmt)
        }

    case *VarDecl:
        Walk(v, &n.VarSpec)

    case *ScanStmt:
        for i := range n.RefList {
            Walk(v, &n.RefList[i])
        }

    case *ScanlnStmt:
        for i := range n.RefList {
            Walk(v, &n.RefList[i])
        }

    case *CheckStmt:
        for i := range n.ExprList {
            Walk(v, &n.ExprList[i])
        }

    case *IfStmt:
        for i := range n.Branches {
            Walk(v, &n.Branches[i])
        }

    case *IfBranch:
        Walk(v, n.Condition)
        Walk(v, &n.Block)

    case *ForStmt:
        switch {
        case n.Range != nil:
            Walk(v, n.Range)
        case n.Scan != nil:
            Walk(v, n.Scan)
        case n.Scanln != nil:
            Walk(v, n.Scanln)
        }
        Walk(v, &n.Block)

    case *EOLStmt:

    case *EOFStmt:

    case *VarSpec:
        Walk(v, &n.Type)

    case *Type:
        Walk(v, n.TypeLit)

    case *TypeLit:
        Walk(v, n.ArrayType)

    case *ArrayType:
        Walk(v, &n.ArrayLength)

    case *Expr:
        Walk(v, n.Left)
        for _, n := range n.Right {
            Walk(v, n)
        }

    case *LogicalOr:
        Walk(v, n.Left)
        for _, n := range n.Right {
            Walk(v, n)
        }

    case *OpLogicalOr:
        Walk(v, n.LogicalOr)

    case *LogicalAnd:
        Walk(v, n.Left)
        for _, n := range n.Right {
            Walk(v, n)
        }

    case *OpLogicalAnd:
        Walk(v, n.LogicalAnd)

    case *Relative:
        Walk(v, n.Left)
        for _, n := range n.Right {
            Walk(v, n)
        }

    case *OpRelative:
        Walk(v, n.Relative)

    case *Addition:
        Walk(v, n.Left)
        for _, n := range n.Right {
            Walk(v, n)
        }

    case *OpAddition:
        Walk(v, n.Addition)

    case *Multiplication:
        Walk(v, n.Unary)
        if n.Exponent != nil {
            Walk(v, n.Exponent)
        }

    case *OpMultiplication:
        Walk(v, n.Factor)

    case *Unary:
        switch {
        case n.Value != nil:
            Walk(v, n.Value)
        case n.Negated != nil:
            Walk(v, n.Negated)
        }

    case *Primary:
        switch {
        // case n.BasicLit != nil:
        //     Walk(v, n.BasicLit)
        case n.CallExpr != nil:
            Walk(v, n.CallExpr)
        case n.Variable != nil:
            Walk(v, n.Variable)
        case n.SubExpr != nil:
            Walk(v, n.SubExpr)
        }

    case *Variable:
        for i := range n.Indices {
            Walk(v, &n.Indices[i])
        }

    case *CallExpr:
        for i := range n.Args {
            Walk(v, &n.Args[i])
        }

    default:
        panic(fmt.Errorf("unreachable, with %T", n))
    }

    v.Visit(nil)
}

type inspector func(Node) bool

func (f inspector) Visit(node Node) Visitor {
    if f(node) {
        return f
    }
    return nil
}

func Inspect(node Node, f func(Node) bool) {
    Walk(inspector(f), node)
}
