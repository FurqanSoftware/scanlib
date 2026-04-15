package eval

import (
	"math"
	"testing"
)

func TestMathAbs(t *testing.T) {
	tests := []struct {
		arg  interface{}
		want interface{}
	}{
		{5, 5},
		{-5, 5},
		{0, 0},
		{int64(-7), int64(7)},
		{int64(7), int64(7)},
		{float32(-2.5), float32(2.5)},
		{float64(-3.14), 3.14},
	}
	for _, tt := range tests {
		got, err := mathAbs(tt.arg)
		if err != nil {
			t.Errorf("math.abs(%v): unexpected error: %v", tt.arg, err)
			continue
		}
		if got != tt.want {
			t.Errorf("math.abs(%v) = %v, want %v", tt.arg, got, tt.want)
		}
	}
}

func TestMathAbs_Invalid(t *testing.T) {
	if _, err := mathAbs("hello"); err == nil {
		t.Error("expected error for string arg")
	}
	if _, err := mathAbs(1, 2); err == nil {
		t.Error("expected error for two args")
	}
}

func TestMathMin(t *testing.T) {
	tests := []struct {
		a, b interface{}
		want interface{}
	}{
		{3, 5, 3},
		{5, 3, 3},
		{int64(10), int64(20), int64(10)},
		{3.5, 2.5, 2.5},
	}
	for _, tt := range tests {
		got, err := mathMin(tt.a, tt.b)
		if err != nil {
			t.Errorf("math.min(%v, %v): unexpected error: %v", tt.a, tt.b, err)
			continue
		}
		if got != tt.want {
			t.Errorf("math.min(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
		}
	}
}

func TestMathMax(t *testing.T) {
	tests := []struct {
		a, b interface{}
		want interface{}
	}{
		{3, 5, 5},
		{5, 3, 5},
		{int64(10), int64(20), int64(20)},
		{3.5, 2.5, 3.5},
	}
	for _, tt := range tests {
		got, err := mathMax(tt.a, tt.b)
		if err != nil {
			t.Errorf("math.max(%v, %v): unexpected error: %v", tt.a, tt.b, err)
			continue
		}
		if got != tt.want {
			t.Errorf("math.max(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
		}
	}
}

func TestMathSqrt(t *testing.T) {
	got, err := mathSqrt(16.0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != 4.0 {
		t.Errorf("math.sqrt(16) = %v, want 4", got)
	}

	got, err = mathSqrt(2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != math.Sqrt(2) {
		t.Errorf("math.sqrt(2) = %v, want %v", got, math.Sqrt(2))
	}
}

func TestMathLog2(t *testing.T) {
	got, err := mathLog2(8.0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != 3.0 {
		t.Errorf("math.log2(8) = %v, want 3", got)
	}
}

func TestMathCeil(t *testing.T) {
	got, err := mathCeil(2.3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != 3.0 {
		t.Errorf("math.ceil(2.3) = %v, want 3", got)
	}
}

func TestMathFloor(t *testing.T) {
	got, err := mathFloor(2.7)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != 2.0 {
		t.Errorf("math.floor(2.7) = %v, want 2", got)
	}
}

func TestMathPow(t *testing.T) {
	tests := []struct {
		name string
		base interface{}
		exp  interface{}
		want interface{}
	}{
		{"int", 2, 10, 1024},
		{"int zero exp", 5, 0, 1},
		{"int64", int64(3), int64(4), int64(81)},
		{"float64", 2.0, 3.0, 8.0},
		{"int negative exp falls to float", 2, -1, 0.5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mathPow(tt.base, tt.exp)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("math.pow(%v, %v) = %v, want %v", tt.base, tt.exp, got, tt.want)
			}
		})
	}
}

func TestMathSum(t *testing.T) {
	tests := []struct {
		name string
		args []interface{}
		want interface{}
	}{
		{"ints", []interface{}{1, 2, 3}, 6},
		{"int64s", []interface{}{int64(10), int64(20)}, int64(30)},
		{"int array", []interface{}{[]int{1, 2, 3, 4}}, 10},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mathSum(tt.args...)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
