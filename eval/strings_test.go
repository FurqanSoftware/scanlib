package eval

import "testing"

func TestStringsDistinct(t *testing.T) {
	tests := []struct {
		s    string
		want bool
	}{
		{"abc", true},
		{"abca", false},
		{"", true},
		{"a", true},
	}
	for _, tt := range tests {
		got, err := stringsDistinct(tt.s)
		if err != nil {
			t.Errorf("strings.distinct(%q): unexpected error: %v", tt.s, err)
			continue
		}
		if got != tt.want {
			t.Errorf("strings.distinct(%q) = %v, want %v", tt.s, got, tt.want)
		}
	}
}

func TestStringsSorted(t *testing.T) {
	tests := []struct {
		s    string
		want bool
	}{
		{"abc", true},
		{"cba", false},
		{"aab", true},
		{"", true},
		{"a", true},
	}
	for _, tt := range tests {
		got, err := stringsSorted(tt.s)
		if err != nil {
			t.Errorf("strings.sorted(%q): unexpected error: %v", tt.s, err)
			continue
		}
		if got != tt.want {
			t.Errorf("strings.sorted(%q) = %v, want %v", tt.s, got, tt.want)
		}
	}
}

func TestStringsPalindrome(t *testing.T) {
	tests := []struct {
		s    string
		want bool
	}{
		{"racecar", true},
		{"hello", false},
		{"aba", true},
		{"ab", false},
		{"a", true},
		{"", true},
	}
	for _, tt := range tests {
		got, err := stringsPalindrome(tt.s)
		if err != nil {
			t.Errorf("strings.palindrome(%q): unexpected error: %v", tt.s, err)
			continue
		}
		if got != tt.want {
			t.Errorf("strings.palindrome(%q) = %v, want %v", tt.s, got, tt.want)
		}
	}
}

func TestStringsLowercase(t *testing.T) {
	tests := []struct {
		s    string
		want bool
	}{
		{"abc", true},
		{"ABC", false},
		{"aBc", false},
		{"abc1", false},
		{"", true},
	}
	for _, tt := range tests {
		got, err := stringsLowercase(tt.s)
		if err != nil {
			t.Errorf("strings.lowercase(%q): unexpected error: %v", tt.s, err)
			continue
		}
		if got != tt.want {
			t.Errorf("strings.lowercase(%q) = %v, want %v", tt.s, got, tt.want)
		}
	}
}

func TestStringsUppercase(t *testing.T) {
	tests := []struct {
		s    string
		want bool
	}{
		{"ABC", true},
		{"abc", false},
		{"ABc", false},
		{"AB1", false},
		{"", true},
	}
	for _, tt := range tests {
		got, err := stringsUppercase(tt.s)
		if err != nil {
			t.Errorf("strings.uppercase(%q): unexpected error: %v", tt.s, err)
			continue
		}
		if got != tt.want {
			t.Errorf("strings.uppercase(%q) = %v, want %v", tt.s, got, tt.want)
		}
	}
}

func TestStringsAlpha(t *testing.T) {
	tests := []struct {
		s    string
		want bool
	}{
		{"abcXYZ", true},
		{"abc123", false},
		{"abc xyz", false},
		{"", true},
	}
	for _, tt := range tests {
		got, err := stringsAlpha(tt.s)
		if err != nil {
			t.Errorf("strings.alpha(%q): unexpected error: %v", tt.s, err)
			continue
		}
		if got != tt.want {
			t.Errorf("strings.alpha(%q) = %v, want %v", tt.s, got, tt.want)
		}
	}
}

func TestStringsDigit(t *testing.T) {
	tests := []struct {
		s    string
		want bool
	}{
		{"12345", true},
		{"123a5", false},
		{"", true},
	}
	for _, tt := range tests {
		got, err := stringsDigit(tt.s)
		if err != nil {
			t.Errorf("strings.digit(%q): unexpected error: %v", tt.s, err)
			continue
		}
		if got != tt.want {
			t.Errorf("strings.digit(%q) = %v, want %v", tt.s, got, tt.want)
		}
	}
}

func TestStringsAlphanumeric(t *testing.T) {
	tests := []struct {
		s    string
		want bool
	}{
		{"abc123", true},
		{"abc 123", false},
		{"abc!123", false},
		{"", true},
	}
	for _, tt := range tests {
		got, err := stringsAlphanumeric(tt.s)
		if err != nil {
			t.Errorf("strings.alphanumeric(%q): unexpected error: %v", tt.s, err)
			continue
		}
		if got != tt.want {
			t.Errorf("strings.alphanumeric(%q) = %v, want %v", tt.s, got, tt.want)
		}
	}
}

func TestStringsBinary(t *testing.T) {
	tests := []struct {
		s    string
		want bool
	}{
		{"0101", true},
		{"0102", false},
		{"", true},
	}
	for _, tt := range tests {
		got, err := stringsBinary(tt.s)
		if err != nil {
			t.Errorf("strings.binary(%q): unexpected error: %v", tt.s, err)
			continue
		}
		if got != tt.want {
			t.Errorf("strings.binary(%q) = %v, want %v", tt.s, got, tt.want)
		}
	}
}

func TestStringsContains(t *testing.T) {
	tests := []struct {
		s, sub string
		want   bool
	}{
		{"hello world", "world", true},
		{"hello world", "xyz", false},
		{"abc", "abc", true},
		{"abc", "abcd", false},
		{"", "", true},
	}
	for _, tt := range tests {
		got, err := stringsContains(tt.s, tt.sub)
		if err != nil {
			t.Errorf("strings.contains(%q, %q): unexpected error: %v", tt.s, tt.sub, err)
			continue
		}
		if got != tt.want {
			t.Errorf("strings.contains(%q, %q) = %v, want %v", tt.s, tt.sub, got, tt.want)
		}
	}
}
