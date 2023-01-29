package eval

import (
	"errors"
	"math"
	"regexp"
	"strconv"
)

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

	"pow": func(args ...interface{}) (interface{}, error) {
		switch n := args[0].(type) {
		case int:
			exp, ok := toInt(args[1])
			if !ok {
				return nil, ErrInvalidArgument{}
			}
			if exp >= 0 {
				return powInt(n, exp), nil
			}

		case int64:
			exp, ok := toInt64(args[1])
			if !ok {
				return nil, ErrInvalidArgument{}
			}
			if exp >= 0 {
				return powInt64(n, exp), nil
			}
		}

		n, ok := toFloat64(args[0])
		if !ok {
			return nil, ErrInvalidArgument{}
		}
		exp, ok := toFloat64(args[1])
		if !ok {
			return nil, ErrInvalidArgument{}
		}
		return math.Pow(n, exp), nil
	},

	"sum": sum,

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

func powInt(n int, exp int) int {
	r := 1
	for {
		if exp&1 > 0 {
			r *= n
		}
		exp >>= 1
		if exp == 0 {
			break
		}
		n *= n
	}
	return r
}

func powInt64(n int64, exp int64) int64 {
	var r int64 = 1
	for {
		if exp&1 > 0 {
			r *= n
		}
		exp >>= 1
		if exp == 0 {
			break
		}
		n *= n
	}
	return r
}

func sum(args ...interface{}) (interface{}, error) {
	var r interface{} = 0
	for _, a := range args {
		switch a := a.(type) {
		case int:
			ri, _ := toInt(r)
			r = ri + a

		case int64:
			ri, _ := toInt64(r)
			r = ri + a

		case []int:
			args := []interface{}{}
			for _, v := range a {
				args = append(args, v)
			}
			s, err := sum(args...)
			if err != nil {
				return 0, err
			}
			ri, _ := toInt(r)
			r = ri + s.(int)
		}
	}
	return r, nil
}
