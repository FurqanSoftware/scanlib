package eval

func init() {
	Functions["strings.distinct"] = stringsDistinct
	Functions["strings.sorted"] = stringsSorted
	Functions["strings.palindrome"] = stringsPalindrome
	Functions["strings.lowercase"] = stringsLowercase
	Functions["strings.uppercase"] = stringsUppercase
	Functions["strings.alpha"] = stringsAlpha
	Functions["strings.digit"] = stringsDigit
	Functions["strings.alphanumeric"] = stringsAlphanumeric
	Functions["strings.binary"] = stringsBinary
	Functions["strings.contains"] = stringsContains
}

// stringsDistinct returns true if all characters in s are unique.
func stringsDistinct(args ...interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, ErrInvalidArgument{}
	}
	s, ok := toString(args[0])
	if !ok {
		return nil, ErrInvalidArgument{}
	}
	seen := map[rune]bool{}
	for _, r := range s {
		if seen[r] {
			return false, nil
		}
		seen[r] = true
	}
	return true, nil
}

// stringsSorted returns true if s is sorted in non-decreasing order.
func stringsSorted(args ...interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, ErrInvalidArgument{}
	}
	s, ok := toString(args[0])
	if !ok {
		return nil, ErrInvalidArgument{}
	}
	for i := 1; i < len(s); i++ {
		if s[i] < s[i-1] {
			return false, nil
		}
	}
	return true, nil
}

// stringsPalindrome returns true if s is a palindrome.
func stringsPalindrome(args ...interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, ErrInvalidArgument{}
	}
	s, ok := toString(args[0])
	if !ok {
		return nil, ErrInvalidArgument{}
	}
	r := []rune(s)
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		if r[i] != r[j] {
			return false, nil
		}
	}
	return true, nil
}

// stringsLowercase returns true if s contains only lowercase letters.
func stringsLowercase(args ...interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, ErrInvalidArgument{}
	}
	s, ok := toString(args[0])
	if !ok {
		return nil, ErrInvalidArgument{}
	}
	for _, r := range s {
		if r < 'a' || r > 'z' {
			return false, nil
		}
	}
	return true, nil
}

// stringsUppercase returns true if s contains only uppercase letters.
func stringsUppercase(args ...interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, ErrInvalidArgument{}
	}
	s, ok := toString(args[0])
	if !ok {
		return nil, ErrInvalidArgument{}
	}
	for _, r := range s {
		if r < 'A' || r > 'Z' {
			return false, nil
		}
	}
	return true, nil
}

// stringsAlpha returns true if s contains only letters.
func stringsAlpha(args ...interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, ErrInvalidArgument{}
	}
	s, ok := toString(args[0])
	if !ok {
		return nil, ErrInvalidArgument{}
	}
	for _, r := range s {
		if !('a' <= r && r <= 'z') && !('A' <= r && r <= 'Z') {
			return false, nil
		}
	}
	return true, nil
}

// stringsDigit returns true if s contains only digits.
func stringsDigit(args ...interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, ErrInvalidArgument{}
	}
	s, ok := toString(args[0])
	if !ok {
		return nil, ErrInvalidArgument{}
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false, nil
		}
	}
	return true, nil
}

// stringsAlphanumeric returns true if s contains only letters and digits.
func stringsAlphanumeric(args ...interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, ErrInvalidArgument{}
	}
	s, ok := toString(args[0])
	if !ok {
		return nil, ErrInvalidArgument{}
	}
	for _, r := range s {
		if !('a' <= r && r <= 'z') && !('A' <= r && r <= 'Z') && !('0' <= r && r <= '9') {
			return false, nil
		}
	}
	return true, nil
}

// stringsBinary returns true if s contains only '0' and '1'.
func stringsBinary(args ...interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, ErrInvalidArgument{}
	}
	s, ok := toString(args[0])
	if !ok {
		return nil, ErrInvalidArgument{}
	}
	for _, r := range s {
		if r != '0' && r != '1' {
			return false, nil
		}
	}
	return true, nil
}

// stringsContains returns true if s contains the substring sub.
func stringsContains(args ...interface{}) (interface{}, error) {
	if len(args) != 2 {
		return nil, ErrInvalidArgument{}
	}
	s, ok := toString(args[0])
	if !ok {
		return nil, ErrInvalidArgument{}
	}
	sub, ok := toString(args[1])
	if !ok {
		return nil, ErrInvalidArgument{}
	}
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true, nil
		}
	}
	return false, nil
}
