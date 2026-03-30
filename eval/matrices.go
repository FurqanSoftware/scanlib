package eval

func init() {
	Functions["matrices.dimensions"] = matricesDimensions
}

// matricesDimensions returns true if the string array G has exactly R elements,
// each of length C.
func matricesDimensions(args ...interface{}) (interface{}, error) {
	if len(args) != 3 {
		return nil, ErrInvalidArgument{}
	}
	g, ok := args[0].([]string)
	if !ok {
		return nil, ErrInvalidArgument{}
	}
	r, ok := toInt(args[1])
	if !ok {
		return nil, ErrInvalidArgument{}
	}
	c, ok := toInt(args[2])
	if !ok {
		return nil, ErrInvalidArgument{}
	}
	if len(g) != r {
		return false, nil
	}
	for _, row := range g {
		if len(row) != c {
			return false, nil
		}
	}
	return true, nil
}
