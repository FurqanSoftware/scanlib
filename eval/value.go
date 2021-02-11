package eval

type Value struct {
	Type Type
	Data interface{}
}

type Type int

const (
	Invalid Type = iota
	Bool
	Int
	Int64
	String
	Array
)

var ASTType = map[string]Type{
	"bool":   Bool,
	"int":    Int,
	"int64":  Int64,
	"string": String,
}
