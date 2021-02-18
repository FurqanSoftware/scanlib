package ast

import (
    "strings"

    "github.com/alecthomas/participle/v2/lexer"
)

type Node interface {
    node()
}

type Source struct {
    Block Block `@@`
}

type Block struct {
    Statements []*Statement `EOL* ( @@ ( EOL+ @@ )* ) EOL*`
}

type Statement struct {
    Pos lexer.Position

    VarDecl   *VarDecl   `  @@`
    ScanStmt  *ScanStmt  `| @@`
    CheckStmt *CheckStmt `| @@`
    ForStmt   *ForStmt   `| @@`
    // IfStmt    *IfStmt    `| @@`
    EOLStmt *EOLStmt `| @@`
    EOFStmt *EOFStmt `| @@`
}

type VarDecl struct {
    VarSpec VarSpec `"var" @@`
}

type ScanStmt struct {
    Pos lexer.Position

    RefList []Reference `"scan" @@ ( "," @@ )*`
}

type CheckStmt struct {
    Pos lexer.Position

    ExprList []Expr `"check" @@ ( "," @@ )*`
}

type ForStmt struct {
    Range RangeClause `"for" @@ EOL+`
    Block Block       `@@ "end"`
}

type EOLStmt struct {
    Pos lexer.Position

    EOL bool `@"eol"`
}

type EOFStmt struct {
    Pos lexer.Position

    EOF bool `@"eof"`
}

type VarSpec struct {
    IdentList []string `@Ident ( "," @Ident )*`
    Type      Type     `@@`
}

type Type struct {
    TypeName *string  `  @Type`
    TypeLit  *TypeLit `| @@`
}

type TypeLit struct {
    ArrayType *ArrayType `@@`
}

type ArrayType struct {
    ArrayLength Expr `"[" @@ "]"`
    ElementType Type `@@`
}

type Reference struct {
    Ident   string `@Ident`
    Indices []Expr `( "[" @@ "]" )?`
}

type Expr struct {
    Pos lexer.Position

    Left  *LogicalOr     `@@`
    Right []*OpLogicalOr `@@*`
}

type LogicalOr struct {
    Pos lexer.Position

    Left  *LogicalAnd     `@@`
    Right []*OpLogicalAnd `@@*`
}

type OpLogicalOr struct {
    Pos lexer.Position

    LogicalOr *LogicalOr `"|" "|" @@`
}

type LogicalAnd struct {
    Pos lexer.Position

    Left  *Relative     `@@`
    Right []*OpRelative `@@*`
}

type OpLogicalAnd struct {
    Pos lexer.Position

    LogicalAnd *LogicalAnd `"&" "&" @@`
}

type Relative struct {
    Left  *Addition     `@@`
    Right []*OpAddition `@@*`
}

type OpRelative struct {
    Pos lexer.Position

    Operator Operator  `@("=" "=" | "!" "=" | "<" "=" | ">" "=" | "<" | ">")`
    Relative *Relative `@@`
}

type Addition struct {
    Left  *Multiplication     `@@`
    Right []*OpMultiplication `@@*`
}

type OpAddition struct {
    Pos lexer.Position

    Operator Operator  `@("+" | "-")`
    Addition *Addition `@@`
}

type Multiplication struct {
    Unary    *Unary   `@@`
    Exponent *Primary `( "^" @@ )?`
}

type OpMultiplication struct {
    Operator Operator        `@("*" | "/")`
    Factor   *Multiplication `@@`
}

type Unary struct {
    Value   *Primary `( "+"? @@`
    Negated *Primary `| "-" @@ )`
}

type Primary struct {
    BasicLit *BasicLit `  @@`
    CallExpr *CallExpr `| @@`
    Variable *Variable `| @@`
    SubExpr  *Expr     `| "(" @@ ")"`
}

type BasicLit struct {
    FloatLit  *float64 `  @Float`
    IntLit    *int64   `| @Int`
    StringLit *string  `| @String`
}

type Operator string

func (o *Operator) Capture(s []string) error {
    *o = Operator(strings.Join(s, ""))
    return nil
}

type RangeClause struct {
    Index string `@Ident ":" "="`
    Low   Expr   `@@ "." "." "."`
    High  Expr   `@@`
}

type Variable struct {
    Ident   string `@Ident`
    Indices []Expr `( "[" @@ "]" )?`
}

type CallExpr struct {
    Ident string `@Ident`
    Args  []Expr `"(" ( @@ ( "," @@ )* )? ")"`
}

func (Source) node()           {}
func (Block) node()            {}
func (Statement) node()        {}
func (VarDecl) node()          {}
func (ScanStmt) node()         {}
func (CheckStmt) node()        {}
func (ForStmt) node()          {}
func (EOLStmt) node()          {}
func (EOFStmt) node()          {}
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
func (RangeClause) node()      {}
