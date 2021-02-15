package ast

import (
	"strings"

	"github.com/alecthomas/participle/v2/lexer"
)

type Source struct {
	Block Block `@@`
}

type Block struct {
	Statement []*Statement `( @@ ( EOL @@ )* ) EOL?`
}

type Statement struct {
	Pos lexer.Position

	VarDecl   *VarDecl   `  @@`
	ScanStmt  *ScanStmt  `| @@`
	CheckStmt *CheckStmt `| @@`
	ForStmt   *ForStmt   `| @@`
	EOLStmt   *string    `| @"eol"`
	EOFStmt   *string    `| @"eof"`
	// Comment *string  `( @Comment`
	// Call    *Call    `| @@`
	// For     *For     `| @@ ) EOL`
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

	Left  *Cmp     `@@`
	Right []*OpCmp `@@*`
}

type Cmp struct {
	Left  *Term     `@@`
	Right []*OpTerm `@@*`
}

type OpCmp struct {
	Pos lexer.Position

	Operator Operator `@("==" | "!=" | "<=" | ">=" | "<" | ">")`
	Cmp      *Cmp     `@@`
}

type Term struct {
	Left  *Factor     `@@`
	Right []*OpFactor `@@*`
}

type OpTerm struct {
	Pos lexer.Position

	Operator Operator `@("+" | "-")`
	Term     *Term    `@@`
}

type Factor struct {
	Base     *Value `@@`
	Exponent *Value `( "^" @@ )?`
}

type OpFactor struct {
	Operator Operator `@("*" | "/")`
	Factor   *Factor  `@@`
}

type Value struct {
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
