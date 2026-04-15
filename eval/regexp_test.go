package eval

import "testing"

func TestRegexpMatch(t *testing.T) {
	tests := []struct {
		s, pattern string
		want       bool
	}{
		{"hello world", "world", true},
		{"hello world", "^hello", true},
		{"hello world", "^world", false},
		{"abc123", `\d+`, true},
		{"abc", `\d+`, false},
	}
	for _, tt := range tests {
		got, err := regexpMatch(tt.s, tt.pattern)
		if err != nil {
			t.Errorf("regexp.match(%q, %q): unexpected error: %v", tt.s, tt.pattern, err)
			continue
		}
		if got != tt.want {
			t.Errorf("regexp.match(%q, %q) = %v, want %v", tt.s, tt.pattern, got, tt.want)
		}
	}
}

func TestRegexpFullmatch(t *testing.T) {
	tests := []struct {
		s, pattern string
		want       bool
	}{
		{"hello", "hello", true},
		{"hello world", "hello", false},
		{"abc123", `[a-z]+\d+`, true},
		{"abc123!", `[a-z]+\d+`, false},
	}
	for _, tt := range tests {
		got, err := regexpFullmatch(tt.s, tt.pattern)
		if err != nil {
			t.Errorf("regexp.fullmatch(%q, %q): unexpected error: %v", tt.s, tt.pattern, err)
			continue
		}
		if got != tt.want {
			t.Errorf("regexp.fullmatch(%q, %q) = %v, want %v", tt.s, tt.pattern, got, tt.want)
		}
	}
}

func TestRegexpCount(t *testing.T) {
	tests := []struct {
		s, pattern string
		want       int
	}{
		{"aaa", "a", 3},
		{"abab", "ab", 2},
		{"hello", "xyz", 0},
		{"a1b2c3", `\d`, 3},
	}
	for _, tt := range tests {
		got, err := regexpCount(tt.s, tt.pattern)
		if err != nil {
			t.Errorf("regexp.count(%q, %q): unexpected error: %v", tt.s, tt.pattern, err)
			continue
		}
		if got != tt.want {
			t.Errorf("regexp.count(%q, %q) = %v, want %v", tt.s, tt.pattern, got, tt.want)
		}
	}
}

func TestRegexp_InvalidPattern(t *testing.T) {
	_, err := regexpMatch("test", "[invalid")
	if err == nil {
		t.Error("expected error for invalid regex pattern")
	}
}
