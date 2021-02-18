package ast

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer/stateful"
)

var parser = participle.MustBuild(&Source{},
	participle.Lexer(stateful.MustSimple([]stateful.Rule{
		{"comment", `#[^\n]*`, nil},
		{"whitespace", `[ \t]+`, nil},
		{"Float", `\d+\.\d*`, nil},
		{"Int", `\d+`, nil},
		{"String", `"(\\"|[^"])*"`, nil},
		{"Keyword", `end|eof|eol|for|scan|var`, nil},
		{"Type", `\b(bool|float32|float64|int|int64|string)\b`, nil},
		{"Ident", `[a-zA-Z_][a-zA-Z0-9_]*`, nil},
		{"Punct", `[-[!@#$%^&*()+_={}\|:;"'<,>.?/]|]`, nil},
		{"EOL", `[\n\r]+`, nil},
	})),
	participle.Elide("comment"),
	participle.Unquote("String"),
	participle.UseLookahead(2),
)

func ParseString(filename string, s string) (*Source, error) {
	n := Source{}
	err := parser.ParseString(filename, s, &n)
	if err != nil {
		return nil, err
	}
	return &n, nil
}
