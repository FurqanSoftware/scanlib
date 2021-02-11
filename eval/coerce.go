package eval

import (
	"strconv"

	"git.furqansoftware.net/toph/scanlib/ast"
)

func toInt(v interface{}) (int, bool) {
	switch v := v.(type) {
	case int:
		return v, true
	case int64:
		return 0, false
	case ast.Number:
		n, err := strconv.ParseInt(string(v), 10, 32)
		if err != nil {
			return 0, false
		}
		return int(n), true
	}
	return 0, false
}

func toInt64(v interface{}) (int64, bool) {
	switch v := v.(type) {
	case int:
		return int64(v), true
	case int64:
		return v, true
	case ast.Number:
		n, err := strconv.ParseInt(string(v), 10, 64)
		if err != nil {
			return 0, false
		}
		return n, true
	}
	return 0, false
}
