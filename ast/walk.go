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
        for _, n := range n.RefList {
            Walk(v, &n)
        }

    case *ScanlnStmt:
        for _, n := range n.RefList {
            Walk(v, &n)
        }

    case *CheckStmt:
        for _, n := range n.ExprList {
            Walk(v, &n)
        }

    case *IfStmt:
        for _, n := range n.Branches {
            Walk(v, &n)
        }

    case *IfBranch:
        Walk(v, n.Condition)
        Walk(v, &n.Block)

    case *ForStmt:
        Walk(v, &n.Range)
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
        Walk(v, n.Exponent)

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
