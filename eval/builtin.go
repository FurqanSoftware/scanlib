package eval

import "regexp"

var Functions = map[string]func(args ...interface{}) (interface{}, error){
	"len": func(args ...interface{}) (interface{}, error) {
		s, ok := args[0].(string)
		if !ok {
			return nil, ErrInvalidArgument{}
		}
		return len(s), nil
	},

	"re": func(args ...interface{}) (interface{}, error) {
		s, ok := args[0].(string)
		if !ok {
			return nil, ErrInvalidArgument{}
		}
		expr, ok := args[1].(string)
		if !ok {
			return nil, ErrInvalidArgument{}
		}
		re, err := regexp.Compile(expr)
		if err != nil {
			return nil, err
		}
		return re.MatchString(s), nil
	},
}
