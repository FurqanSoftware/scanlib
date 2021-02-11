package eval

func Zero(t Type) Value {
	switch t {
	case Bool:
		return Value{Bool, bool(false)}
	case Int:
		return Value{Int, int(0)}
	case Int64:
		return Value{Int64, int64(0)}
	case String:
		return Value{String, string("")}
	}
	return Value{}
}

func MakeArray(t Type, sizes []int) Value {
	if len(sizes) == 0 {
		return Zero(t)
	}
	r := []Value{}
	for i := 0; i < sizes[0]; i++ {
		r = append(r, MakeArray(t, sizes[1:]))
	}
	return Value{Array, r}
}
