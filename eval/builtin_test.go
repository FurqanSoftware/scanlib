package eval

import (
	"testing"
)

func TestLen(t *testing.T) {
	tests := []struct {
		arg  interface{}
		want int
	}{
		{"", 0},
		{"hello", 5},
		{"abc def", 7},
	}
	for _, tt := range tests {
		got, err := Functions["len"](tt.arg)
		if err != nil {
			t.Errorf("len(%q): unexpected error: %v", tt.arg, err)
			continue
		}
		if got != tt.want {
			t.Errorf("len(%q) = %v, want %v", tt.arg, got, tt.want)
		}
	}
}

func TestLen_InvalidArg(t *testing.T) {
	_, err := Functions["len"](42)
	if err == nil {
		t.Error("len(42): expected error")
	}
}

func TestToInt64(t *testing.T) {
	tests := []struct {
		name string
		args []interface{}
		want int64
	}{
		{"from int", []interface{}{42}, 42},
		{"from string decimal", []interface{}{"123"}, 123},
		{"from string with base", []interface{}{"ff", 16}, 255},
		{"from string binary", []interface{}{"1010", 2}, 10},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Functions["toInt64"](tt.args...)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToInt64_Invalid(t *testing.T) {
	_, err := Functions["toInt64"](3.14)
	if err == nil {
		t.Error("toInt64(3.14): expected error")
	}

	_, err = Functions["toInt64"]("abc")
	if err == nil {
		t.Error("toInt64(\"abc\"): expected error for non-numeric string")
	}
}
