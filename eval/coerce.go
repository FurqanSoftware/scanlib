package eval

func toBool(v interface{}) (bool, bool) {
	switch v := v.(type) {
	case bool:
		return v, true
	}
	return false, false
}

func toInt(v interface{}) (int, bool) {
	switch v := v.(type) {
	case int:
		return v, true
	case int64:
		return int(v), true
	}
	return 0, false
}

func toInt64(v interface{}) (int64, bool) {
	switch v := v.(type) {
	case int:
		return int64(v), true
	case int64:
		return v, true
	}
	return 0, false
}

func toFloat32(v interface{}) (float32, bool) {
	switch v := v.(type) {
	case int:
		return float32(v), true
	case int64:
		return float32(v), true
	case float32:
		return v, true
	case float64:
		return float32(v), true
	}
	return 0, false
}

func toFloat64(v interface{}) (float64, bool) {
	switch v := v.(type) {
	case int:
		return float64(v), true
	case int64:
		return float64(v), true
	case float32:
		return float64(v), true
	case float64:
		return v, true
	}
	return 0, false
}

func toString(v interface{}) (string, bool) {
	switch v := v.(type) {
	case string:
		return v, true
	}
	return "", false
}
