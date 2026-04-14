package ast

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

var parser = participle.MustBuild[Source](participle.Lexer(lexer.MustSimple([]lexer.SimpleRule{
	{Name: "comment", Pattern: `#[^\n]*`},
	{Name: "whitespace", Pattern: `[ \t]+`},
	{Name: "Float", Pattern: `\d+\.\d*`},
	{Name: "Int", Pattern: `\d+`},
	{Name: "String", Pattern: `"(\\"|[^"])*"`},
	{Name: "Keyword", Pattern: `end|eof|eol|for|scanln|scan|var`},
	{Name: "Type", Pattern: `\b(bool|float32|float64|int|int64|string)\b`},
	{Name: "Ident", Pattern: `[a-zA-Z_][a-zA-Z0-9_]*`},
	{Name: "Punct", Pattern: `[-[!@#$%^&*()+_={}\|:;"'<,>.?/]|]`},
	{Name: "EOL", Pattern: `[\n\r]+`},
})),
	participle.Elide("comment"),
	participle.Unquote("String"),
	participle.UseLookahead(2),
)

// ParseString parses a Scanspec source string and returns the AST.
func ParseString(filename string, s string) (*Source, error) {
	n, err := parser.ParseString(filename, s)
	if err != nil {
		return nil, err
	}
	return n, nil
}
