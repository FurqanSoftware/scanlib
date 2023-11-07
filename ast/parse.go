package ast

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

var parser = participle.MustBuild[Source](participle.Lexer(lexer.MustSimple([]lexer.SimpleRule{
	{"comment", `#[^\n]*`},
	{"whitespace", `[ \t]+`},
	{"Float", `\d+\.\d*`},
	{"Int", `\d+`},
	{"String", `"(\\"|[^"])*"`},
	{"Keyword", `end|eof|eol|for|scanln|scan|var`},
	{"Type", `\b(bool|float32|float64|int|int64|string)\b`},
	{"Ident", `[a-zA-Z_][a-zA-Z0-9_]*`},
	{"Punct", `[-[!@#$%^&*()+_={}\|:;"'<,>.?/]|]`},
	{"EOL", `[\n\r]+`},
})),
	participle.Elide("comment"),
	participle.Unquote("String"),
	participle.UseLookahead(2),
)

func ParseString(filename string, s string) (*Source, error) {
	n, err := parser.ParseString(filename, s)
	if err != nil {
		return nil, err
	}
	return n, nil
}
