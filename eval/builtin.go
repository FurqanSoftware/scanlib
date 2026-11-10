package eval

import (
	"errors"
	"reflect"
	"strconv"
)

// Functions maps built-in function names to their implementations.
var Functions = map[string]func(args ...interface{}) (interface{}, error){
	"len": func(args ...interface{}) (interface{}, error) {
		if len(args) != 1 {
			return nil, ErrInvalidArgument{}
		}
		if s, ok := args[0].(string); ok {
			return len(s), nil
		}
		if v := reflect.ValueOf(args[0]); v.Kind() == reflect.Slice {
			return v.Len(), nil
		}
		return nil, ErrInvalidArgument{}
	},

	"toInt64": func(args ...interface{}) (interface{}, error) {
		switch n := args[0].(type) {
		case int:
			return int64(n), nil

		case int64:
			return n, nil

		case string:
			base := 10
			if len(args) == 2 {
				var ok bool
				base, ok = toInt(args[1])
				if !ok {
					return 0, errors.New("toInt64: base is not int")
				}
			}
			return strconv.ParseInt(n, base, 64)

		default:
			return 0, errors.New("toInt64: want string")
		}
	},
}
