package eval

import "reflect"

var Types = map[string]reflect.Type{
	"bool":    reflect.TypeOf(bool(false)),
	"int":     reflect.TypeOf(int(0)),
	"int64":   reflect.TypeOf(int64(0)),
	"float32": reflect.TypeOf(float32(0)),
	"float64": reflect.TypeOf(float64(0)),
	"string":  reflect.TypeOf(string("")),
}

type Values map[string]reflect.Value
