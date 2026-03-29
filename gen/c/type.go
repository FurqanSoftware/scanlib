package c

var ASTType = map[string]string{
	"bool":    "int",
	"int":     "int",
	"int64":   "long long int",
	"float32": "float",
	"float64": "double",
}

var ScanFormat = map[string]string{
	"int":           "%d",
	"long long int": "%lld",
	"float":         "%f",
	"double":        "%lf",
}
