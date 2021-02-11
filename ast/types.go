package ast

import (
	"strings"
)

type Source struct {
	Block Block `@@`
}

type Block struct {
	Statement []*Statement `@@*`
}

type Statement struct {
	VarDecl   *VarDecl   `( @@`
	ScanStmt  *ScanStmt  `| @@`
	CheckStmt *CheckStmt `| @@`
	ForStmt   *ForStmt   `| @@`
	EOLStmt   *string    `| @"eol"`
	EOFStmt   *string    `| @"eof" ) EOL`
	// Comment *string  `( @Comment`
	// Call    *Call    `| @@`
	// For     *For     `| @@ ) EOL`
}

type VarDecl struct {
	VarSpec VarSpec `"var" @@`
}

type ScanStmt struct {
	RefList []Reference `"scan" @@ ( "," @@ )*`
}

type CheckStmt struct {
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
	Left  *Cmp     `@@`
	Right []*OpCmp `@@*`
}

type Cmp struct {
	Left  *Term     `@@`
	Right []*OpTerm `@@*`
}

type OpCmp struct {
	Operator Operator `@("==" | "!=" | "<=" | ">=" | "<" | ">")`
	Cmp      *Cmp     `@@`
}

type Term struct {
	Left  *Factor     `@@`
	Right []*OpFactor `@@*`
}

type OpTerm struct {
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
	Index string     `@Identifier`
	Low   Expression `@@`
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
