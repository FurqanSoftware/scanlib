package ast

import (
    "strings"

    "github.com/alecthomas/participle/v2/lexer"
)

type Source struct {
    Block Block `@@`
}

type Block struct {
    Statement []*Statement `EOL* ( @@ ( EOL+ @@ )* ) EOL*`
}

type Statement struct {
    Pos lexer.Position

    VarDecl   *VarDecl   `  @@`
    ScanStmt  *ScanStmt  `| @@`
    CheckStmt *CheckStmt `| @@`
    ForStmt   *ForStmt   `| @@`
    EOLStmt   *string    `| @"eol"`
    EOFStmt   *string    `| @"eof"`
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

    ExprList []Expression `"check" @@ ( "," @@ )*`
}

type ForStmt struct {
    Range RangeClause `"for" @@ EOL`
    Block Block       `@@ "end"`
}

type VarSpec struct {
    IdentList []string `@Identifier ( "," @Identifier )*`
    Type      Type     `@@`
}

type Type struct {
    TypeName *string  `  @Identifier`
    TypeLit  *TypeLit `| @@`
}

type TypeLit struct {
    ArrayType *ArrayType `@@`
}

type ArrayType struct {
    ArrayLength Expression `"[" @@ "]"`
    ElementType Type       `@@`
}

type Reference struct {
    Identifier string       `@Identifier`
    Indices    []Expression `( "[" @@ "]" )?`
}

type Expression struct {
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

    LogicalOr *LogicalOr `"||" @@`
}

type LogicalAnd struct {
    Pos lexer.Position

    Left  *Relative     `@@`
    Right []*OpRelative `@@*`
}

type OpLogicalAnd struct {
    Pos lexer.Position

    LogicalAnd *LogicalAnd `"&&" @@`
}

type Relative struct {
    Left  *Addition     `@@`
    Right []*OpAddition `@@*`
}

type OpRelative struct {
    Pos lexer.Position

    Operator Operator  `@("==" | "!=" | "<=" | ">=" | "<" | ">")`
    Cmp      *Relative `@@`
}

type Addition struct {
    Left  *Multiplication     `@@`
    Right []*OpMultiplication `@@*`
}

type OpAddition struct {
    Pos lexer.Position

    Operator Operator  `@("+" | "-")`
    Term     *Addition `@@`
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
    Number        *Number     `  @Number`
    Call          *Call       `| @@`
    Variable      *Variable   `| @@`
    String        *string     `| @String`
    Subexpression *Expression `| "(" @@ ")"`
}

type Number string

type Operator string

func (o *Operator) Capture(s []string) error {
    *o = Operator(strings.Join(s, ""))
    return nil
}

type RangeClause struct {
    Index string     `@Identifier ":="`
    Low   Expression `@@ "..."`
    High  Expression `@@`
}

type Variable struct {
    Identifier string       `@Identifier`
    Indices    []Expression `( "[" @@ "]" )?`
}

type Call struct {
    Name      string       `@Identifier`
    Arguments []Expression `"(" ( @@ ( "," @@ )* )? ")"`
}
