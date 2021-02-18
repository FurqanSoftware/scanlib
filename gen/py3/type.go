package py3

var ASTType = map[string]string{
	"bool":    "bool",
	"int":     "int",
	"int64":   "int",
	"float32": "float",
	"float64": "float",
	"string":  "string",
}

var ASTZero = map[string]string{
	"bool":    "false",
	"int":     "0",
	"int64":   "0",
	"float32": "0.0",
	"float64": "0.0",
	"string":  `""`,
}
