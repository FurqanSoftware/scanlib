package py3

var ASTType = map[string]string{
	"bool":   "bool",
	"int":    "int",
	"int64":  "int",
	"string": "string",
}

var ASTZero = map[string]string{
	"bool":   "false",
	"int":    "0",
	"int64":  "0",
	"string": `""`,
}
