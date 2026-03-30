package eval

import (
	"errors"
	"strconv"
)

// Functions maps built-in function names to their implementations.
var Functions = map[string]func(args ...interface{}) (interface{}, error){
	"len": func(args ...interface{}) (interface{}, error) {
		s, ok := args[0].(string)
		if !ok {
			return nil, ErrInvalidArgument{}
		}
		return len(s), nil
	},

	"toInt64": func(args ...interface{}) (interface{}, error) {
		switch n := args[0].(type) {
		case int:
			return int64(n), nil

		case string:
			base := 10
			if len(args) == 2 {
				var ok bool
				base, ok = args[1].(int)
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
