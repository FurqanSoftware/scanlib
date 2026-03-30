package eval

import "math"

func init() {
	Functions["math.abs"] = mathAbs
	Functions["math.min"] = mathMin
	Functions["math.max"] = mathMax
	Functions["math.sqrt"] = mathSqrt
	Functions["math.log2"] = mathLog2
	Functions["math.ceil"] = mathCeil
	Functions["math.floor"] = mathFloor
	Functions["math.pow"] = mathPow
	Functions["math.sum"] = mathSum
}

// mathAbs returns the absolute value of n.
func mathAbs(args ...interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, ErrInvalidArgument{}
	}
	switch n := args[0].(type) {
	case int:
		if n < 0 {
			return -n, nil
		}
		return n, nil
	case int64:
		if n < 0 {
			return -n, nil
		}
		return n, nil
	case float32:
		return float32(math.Abs(float64(n))), nil
	case float64:
		return math.Abs(n), nil
	}
	return nil, ErrInvalidArgument{}
}

// mathMin returns the minimum of a and b.
func mathMin(args ...interface{}) (interface{}, error) {
	if len(args) != 2 {
		return nil, ErrInvalidArgument{}
	}
	switch a := args[0].(type) {
	case int:
		b, ok := toInt(args[1])
		if !ok {
			return nil, ErrInvalidArgument{}
		}
		if a < b {
			return a, nil
		}
		return b, nil
	case int64:
		b, ok := toInt64(args[1])
		if !ok {
			return nil, ErrInvalidArgument{}
		}
		if a < b {
			return a, nil
		}
		return b, nil
	case float64:
		b, ok := toFloat64(args[1])
		if !ok {
			return nil, ErrInvalidArgument{}
		}
		return math.Min(a, b), nil
	}
	return nil, ErrInvalidArgument{}
}

// mathMax returns the maximum of a and b.
func mathMax(args ...interface{}) (interface{}, error) {
	if len(args) != 2 {
		return nil, ErrInvalidArgument{}
	}
	switch a := args[0].(type) {
	case int:
		b, ok := toInt(args[1])
		if !ok {
			return nil, ErrInvalidArgument{}
		}
		if a > b {
			return a, nil
		}
		return b, nil
	case int64:
		b, ok := toInt64(args[1])
		if !ok {
			return nil, ErrInvalidArgument{}
		}
		if a > b {
			return a, nil
		}
		return b, nil
	case float64:
		b, ok := toFloat64(args[1])
		if !ok {
			return nil, ErrInvalidArgument{}
		}
		return math.Max(a, b), nil
	}
	return nil, ErrInvalidArgument{}
}

// mathSqrt returns the square root of n.
func mathSqrt(args ...interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, ErrInvalidArgument{}
	}
	n, ok := toFloat64(args[0])
	if !ok {
		return nil, ErrInvalidArgument{}
	}
	return math.Sqrt(n), nil
}

// mathLog2 returns the base-2 logarithm of n.
func mathLog2(args ...interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, ErrInvalidArgument{}
	}
	n, ok := toFloat64(args[0])
	if !ok {
		return nil, ErrInvalidArgument{}
	}
	return math.Log2(n), nil
}

// mathCeil returns the smallest integer not less than n.
func mathCeil(args ...interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, ErrInvalidArgument{}
	}
	n, ok := toFloat64(args[0])
	if !ok {
		return nil, ErrInvalidArgument{}
	}
	return math.Ceil(n), nil
}

// mathFloor returns the largest integer not greater than n.
func mathFloor(args ...interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, ErrInvalidArgument{}
	}
	n, ok := toFloat64(args[0])
	if !ok {
		return nil, ErrInvalidArgument{}
	}
	return math.Floor(n), nil
}

// mathPow returns n raised to the power of e.
func mathPow(args ...interface{}) (interface{}, error) {
	if len(args) != 2 {
		return nil, ErrInvalidArgument{}
	}
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
}

// mathSum returns the sum of the arguments.
func mathSum(args ...interface{}) (interface{}, error) {
	return sum(args...)
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
