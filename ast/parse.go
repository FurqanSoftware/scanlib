package ast

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer/stateful"
)

var (
	inputLexer = stateful.MustSimple([]stateful.Rule{
		{"Comment", `(?i)#[^\n]*`, nil},
		{"Keyword", "end|eof|eol|for|scan|var", nil},
		{"Operator", "\\|\\||&&|==|!=|<=|>=|<|>", nil},
		{"String", `"(\\"|[^"])*"`, nil},
		{"Number", `[-+]?(\d*\.)?\d+`, nil},
		{"Identifier", `[a-zA-Z_]\w*`, nil},
		{"Punct", `[-[!@#$%^&*()+_={}\|:;"'<,>.?/]|]`, nil},
		{"EOL", `[\n\r]+`, nil},
		{"whitespace", `[ \t]+`, nil},
	})

	inputParser = participle.MustBuild(&Source{},
		participle.Lexer(inputLexer),
		participle.Unquote("String"),
		// participle.UseLookahead(2),
	)
)

func ParseString(filename string, s string) (*Source, error) {
	n := Source{}
	err := inputParser.ParseString(filename, s, &n)
	if err != nil {
		return nil, err
	}
	return &n, nil
}
