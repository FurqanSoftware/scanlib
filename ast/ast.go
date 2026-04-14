// Package ast defines the abstract syntax tree for the Scanspec language.
package ast

import (
    "strings"

    "github.com/alecthomas/participle/v2/lexer"
)

// Node is the interface implemented by all AST nodes.
type Node interface {
    node()
}

// Source is the root node of a Scanspec AST.
type Source struct {
    Block Block `parser:"@@"`
}

type Block struct {
    Statements []*Statement `parser:"EOL* ( @@ ( EOL+ @@ )* ) EOL*"`
}

type Statement struct {
    Pos lexer.Position

    VarDecl    *VarDecl    `parser:"  @@"`
    ScanStmt   *ScanStmt   `parser:"| @@"`
    ScanlnStmt *ScanlnStmt `parser:"| @@"`
    CheckStmt  *CheckStmt  `parser:"| @@"`
    IfStmt     *IfStmt     `parser:"| @@"`
    ForStmt    *ForStmt    `parser:"| @@"`
    EOLStmt    *EOLStmt    `parser:"| @@"`
    EOFStmt    *EOFStmt    `parser:"| @@"`
    AssignStmt *AssignStmt `parser:"| @@"`
}

type VarDecl struct {
    VarSpec VarSpec `parser:"'var' @@"`
}

type ScanStmt struct {
    Pos lexer.Position

    RefList []Reference `parser:"'scan' @@ ( ',' @@ )*"`
}

type ScanlnStmt struct {
    Pos lexer.Position

    RefList []Reference `parser:"'scanln' @@ ( ',' @@ )*"`
}

type CheckStmt struct {
    Pos lexer.Position

    ExprList []Expr `parser:"'check' @@ ( ',' @@ )*"`
}

type IfStmt struct {
    Branches []IfBranch `parser:"( @@ ( EOL* 'else' @@ )* ) EOL* 'end'"`
}

type IfBranch struct {
    Condition *Expr `parser:"( 'if' @@ )? EOL+"`
    Block     Block `parser:"@@"`
}

type ForStmt struct {
    Range  *RangeClause `parser:"'for' ( @@"`
    Scan   *ScanStmt    `parser:"| @@"`
    Scanln *ScanlnStmt  `parser:"| @@ ) EOL+"`
    Block  Block        `parser:"@@ 'end'"`
}

type EOLStmt struct {
    Pos lexer.Position

    EOL bool `parser:"@'eol'"`
}

type EOFStmt struct {
    Pos lexer.Position

    EOF bool `parser:"@'eof'"`
}

type VarSpec struct {
    IdentList []string `parser:"@Ident ( ',' @Ident )*"`
    Type      Type     `parser:"@@"`
}

type AssignStmt struct {
    Pos lexer.Position

    Ref   Reference `parser:"@@"`
    Value Expr      `parser:"'=' @@"`
}

type Type struct {
    TypeName *string  `parser:"  @Type"`
    TypeLit  *TypeLit `parser:"| @@"`
}

type TypeLit struct {
    ArrayType *ArrayType `parser:"@@"`
}

type ArrayType struct {
    ArrayLength Expr `parser:"'[' @@ ']'"`
    ElementType Type `parser:"@@"`
}

type Reference struct {
    Pos lexer.Position

    Ident   string `parser:"@Ident"`
    Indices []Expr `parser:"( '[' @@ ']' )*"`
}

type Expr struct {
    Pos    lexer.Position
    Tokens []lexer.Token

    Left  *LogicalOr     `parser:"@@"`
    Right []*OpLogicalOr `parser:"@@*"`
}

type LogicalOr struct {
    Pos lexer.Position

    Left  *LogicalAnd     `parser:"@@"`
    Right []*OpLogicalAnd `parser:"@@*"`
}

type OpLogicalOr struct {
    Pos lexer.Position

    LogicalOr *LogicalOr `parser:"'|' '|' @@"`
}

type LogicalAnd struct {
    Pos lexer.Position

    Left  *Relative     `parser:"@@"`
    Right []*OpRelative `parser:"@@*"`
}

type OpLogicalAnd struct {
    Pos lexer.Position

    LogicalAnd *LogicalAnd `parser:"'&' '&' @@"`
}

type Relative struct {
    Left  *Addition     `parser:"@@"`
    Right []*OpAddition `parser:"@@*"`
}

type OpRelative struct {
    Pos lexer.Position

    Operator Operator  `parser:"@('=' '=' | '!' '=' | '<' '=' | '>' '=' | '<' | '>')"`
    Relative *Relative `parser:"@@"`
}

type Addition struct {
    Left  *Multiplication     `parser:"@@"`
    Right []*OpMultiplication `parser:"@@*"`
}

type OpAddition struct {
    Pos lexer.Position

    Operator Operator  `parser:"@('+' | '-')"`
    Addition *Addition `parser:"@@"`
}

type Multiplication struct {
    Unary    *Unary   `parser:"@@"`
    Exponent *Primary `parser:"( '*' '*' @@ )?"`
}

type OpMultiplication struct {
    Pos lexer.Position

    Operator Operator        `parser:"@('*' | '/')"`
    Factor   *Multiplication `parser:"@@"`
}

type Unary struct {
    Value   *Primary `parser:"( '+'? @@"`
    Negated *Primary `parser:"| '-' @@ )"`
}

type Primary struct {
    Pos lexer.Position

    BasicLit       *BasicLit       `parser:"  @@"`
    ModuleCallExpr *ModuleCallExpr `parser:"| @@"`
    CallExpr       *CallExpr       `parser:"| @@"`
    Variable       *Variable       `parser:"| @@"`
    SubExpr        *Expr           `parser:"| '(' @@ ')'"`
}

type BasicLit struct {
    FloatLit  *float64 `parser:"  @Float"`
    IntLit    *int64   `parser:"| @Int"`
    StringLit *string  `parser:"| @String"`
}

type Operator string

func (o *Operator) Capture(s []string) error {
    *o = Operator(strings.Join(s, ""))
    return nil
}

type RangeClause struct {
    Index string `parser:"@Ident ':' '='"`
    Low   Expr   `parser:"@@ '.' '.' '.'"`
    High  Expr   `parser:"@@"`
}

type Variable struct {
    Ident   string `parser:"@Ident"`
    Indices []Expr `parser:"( '[' @@ ']' )?"`
}

type CallExpr struct {
    Ident string `parser:"@Ident"`
    Args  []Expr `parser:"'(' ( @@ ( ',' @@ )* )? ')'"`
}

type ModuleCallExpr struct {
    Module string `parser:"@Ident '.'"`
    Ident  string `parser:"@Ident '('"`
    Args   []Expr `parser:"( @@ ( ',' @@ )* )? ')'"`
}

func (Source) node()           {}
func (Block) node()            {}
func (Statement) node()        {}
func (VarDecl) node()          {}
func (ScanStmt) node()         {}
func (ScanlnStmt) node()       {}
func (CheckStmt) node()        {}
func (IfStmt) node()           {}
func (IfBranch) node()         {}
func (ForStmt) node()          {}
func (EOLStmt) node()          {}
func (EOFStmt) node()          {}
func (AssignStmt) node()       {}
func (VarSpec) node()          {}
func (Type) node()             {}
func (TypeLit) node()          {}
func (ArrayType) node()        {}
func (Reference) node()        {}
func (Expr) node()             {}
func (LogicalOr) node()        {}
func (OpLogicalOr) node()      {}
func (LogicalAnd) node()       {}
func (OpLogicalAnd) node()     {}
func (Relative) node()         {}
func (OpRelative) node()       {}
func (Addition) node()         {}
func (OpAddition) node()       {}
func (Multiplication) node()   {}
func (OpMultiplication) node() {}
func (Unary) node()            {}
func (Primary) node()          {}
func (BasicLit) node()         {}
func (RangeClause) node()      {}
func (Variable) node()         {}
func (CallExpr) node()         {}
func (ModuleCallExpr) node()   {}
