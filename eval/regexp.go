package eval

import "regexp"

func init() {
	Functions["regexp.match"] = regexpMatch
	Functions["regexp.fullmatch"] = regexpFullmatch
	Functions["regexp.count"] = regexpCount
}

// regexpMatch returns true if s matches the regular expression pattern.
func regexpMatch(args ...interface{}) (interface{}, error) {
	if len(args) != 2 {
		return nil, ErrInvalidArgument{}
	}
	s, ok := toString(args[0])
	if !ok {
		return nil, ErrInvalidArgument{}
	}
	pattern, ok := toString(args[1])
	if !ok {
		return nil, ErrInvalidArgument{}
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	return re.MatchString(s), nil
}

// regexpFullmatch returns true if the entire string s matches the pattern.
func regexpFullmatch(args ...interface{}) (interface{}, error) {
	if len(args) != 2 {
		return nil, ErrInvalidArgument{}
	}
	s, ok := toString(args[0])
	if !ok {
		return nil, ErrInvalidArgument{}
	}
	pattern, ok := toString(args[1])
	if !ok {
		return nil, ErrInvalidArgument{}
	}
	re, err := regexp.Compile("^(?:" + pattern + ")$")
	if err != nil {
		return nil, err
	}
	return re.MatchString(s), nil
}

// regexpCount returns the number of non-overlapping matches of pattern in s.
func regexpCount(args ...interface{}) (interface{}, error) {
	if len(args) != 2 {
		return nil, ErrInvalidArgument{}
	}
	s, ok := toString(args[0])
	if !ok {
		return nil, ErrInvalidArgument{}
	}
	pattern, ok := toString(args[1])
	if !ok {
		return nil, ErrInvalidArgument{}
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	return len(re.FindAllStringIndex(s, -1)), nil
}
