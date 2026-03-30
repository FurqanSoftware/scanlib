package eval

func init() {
	Functions["arrays.sorted"] = arraysSorted
	Functions["arrays.distinct"] = arraysDistinct
	Functions["arrays.permutation"] = arraysPermutation
	Functions["arrays.range"] = arraysRange
}

// arraysSorted returns true if the array is sorted in non-decreasing order.
func arraysSorted(args ...interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, ErrInvalidArgument{}
	}
	switch a := args[0].(type) {
	case []int:
		for i := 1; i < len(a); i++ {
			if a[i] < a[i-1] {
				return false, nil
			}
		}
		return true, nil
	case []int64:
		for i := 1; i < len(a); i++ {
			if a[i] < a[i-1] {
				return false, nil
			}
		}
		return true, nil
	case []float32:
		for i := 1; i < len(a); i++ {
			if a[i] < a[i-1] {
				return false, nil
			}
		}
		return true, nil
	case []float64:
		for i := 1; i < len(a); i++ {
			if a[i] < a[i-1] {
				return false, nil
			}
		}
		return true, nil
	case []string:
		for i := 1; i < len(a); i++ {
			if a[i] < a[i-1] {
				return false, nil
			}
		}
		return true, nil
	}
	return nil, ErrInvalidArgument{}
}

// arraysDistinct returns true if all elements in the array are unique.
func arraysDistinct(args ...interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, ErrInvalidArgument{}
	}
	switch a := args[0].(type) {
	case []int:
		seen := map[int]bool{}
		for _, v := range a {
			if seen[v] {
				return false, nil
			}
			seen[v] = true
		}
		return true, nil
	case []int64:
		seen := map[int64]bool{}
		for _, v := range a {
			if seen[v] {
				return false, nil
			}
			seen[v] = true
		}
		return true, nil
	case []float32:
		seen := map[float32]bool{}
		for _, v := range a {
			if seen[v] {
				return false, nil
			}
			seen[v] = true
		}
		return true, nil
	case []float64:
		seen := map[float64]bool{}
		for _, v := range a {
			if seen[v] {
				return false, nil
			}
			seen[v] = true
		}
		return true, nil
	case []string:
		seen := map[string]bool{}
		for _, v := range a {
			if seen[v] {
				return false, nil
			}
			seen[v] = true
		}
		return true, nil
	}
	return nil, ErrInvalidArgument{}
}

// arraysPermutation returns true if the array is a permutation of 1..N.
func arraysPermutation(args ...interface{}) (interface{}, error) {
	if len(args) != 2 {
		return nil, ErrInvalidArgument{}
	}
	n, ok := toInt(args[0])
	if !ok {
		return nil, ErrInvalidArgument{}
	}
	a, ok := args[1].([]int)
	if !ok {
		return nil, ErrInvalidArgument{}
	}
	if len(a) != n {
		return false, nil
	}
	seen := make([]bool, n+1)
	for _, v := range a {
		if v < 1 || v > n || seen[v] {
			return false, nil
		}
		seen[v] = true
	}
	return true, nil
}

// arraysRange returns true if all elements are in [lo, hi].
func arraysRange(args ...interface{}) (interface{}, error) {
	if len(args) != 3 {
		return nil, ErrInvalidArgument{}
	}
	switch a := args[0].(type) {
	case []int:
		lo, ok := toInt(args[1])
		if !ok {
			return nil, ErrInvalidArgument{}
		}
		hi, ok := toInt(args[2])
		if !ok {
			return nil, ErrInvalidArgument{}
		}
		for _, v := range a {
			if v < lo || v > hi {
				return false, nil
			}
		}
		return true, nil
	case []int64:
		lo, ok := toInt64(args[1])
		if !ok {
			return nil, ErrInvalidArgument{}
		}
		hi, ok := toInt64(args[2])
		if !ok {
			return nil, ErrInvalidArgument{}
		}
		for _, v := range a {
			if v < lo || v > hi {
				return false, nil
			}
		}
		return true, nil
	case []float32:
		lo, ok := toFloat32(args[1])
		if !ok {
			return nil, ErrInvalidArgument{}
		}
		hi, ok := toFloat32(args[2])
		if !ok {
			return nil, ErrInvalidArgument{}
		}
		for _, v := range a {
			if v < lo || v > hi {
				return false, nil
			}
		}
		return true, nil
	case []float64:
		lo, ok := toFloat64(args[1])
		if !ok {
			return nil, ErrInvalidArgument{}
		}
		hi, ok := toFloat64(args[2])
		if !ok {
			return nil, ErrInvalidArgument{}
		}
		for _, v := range a {
			if v < lo || v > hi {
				return false, nil
			}
		}
		return true, nil
	}
	return nil, ErrInvalidArgument{}
}
