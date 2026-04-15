package eval

import "testing"

func TestNumberPrime(t *testing.T) {
	tests := []struct {
		n    int
		want bool
	}{
		{0, false},
		{1, false},
		{2, true},
		{3, true},
		{4, false},
		{5, true},
		{6, false},
		{7, true},
		{9, false},
		{11, true},
		{25, false},
		{97, true},
	}
	for _, tt := range tests {
		got, err := numberPrime(tt.n)
		if err != nil {
			t.Errorf("numbers.prime(%d): unexpected error: %v", tt.n, err)
			continue
		}
		if got != tt.want {
			t.Errorf("numbers.prime(%d) = %v, want %v", tt.n, got, tt.want)
		}
	}
}

func TestNumberGCD(t *testing.T) {
	tests := []struct {
		a, b, want int
	}{
		{12, 8, 4},
		{7, 13, 1},
		{0, 5, 5},
		{100, 75, 25},
		{-12, 8, 4},
	}
	for _, tt := range tests {
		got, err := numberGCD(tt.a, tt.b)
		if err != nil {
			t.Errorf("numbers.gcd(%d, %d): unexpected error: %v", tt.a, tt.b, err)
			continue
		}
		if got != tt.want {
			t.Errorf("numbers.gcd(%d, %d) = %v, want %v", tt.a, tt.b, got, tt.want)
		}
	}
}

func TestNumberLCM(t *testing.T) {
	tests := []struct {
		a, b, want int
	}{
		{4, 6, 12},
		{3, 5, 15},
		{7, 7, 7},
	}
	for _, tt := range tests {
		got, err := numberLCM(tt.a, tt.b)
		if err != nil {
			t.Errorf("numbers.lcm(%d, %d): unexpected error: %v", tt.a, tt.b, err)
			continue
		}
		if got != tt.want {
			t.Errorf("numbers.lcm(%d, %d) = %v, want %v", tt.a, tt.b, got, tt.want)
		}
	}
}

func TestNumberCoprime(t *testing.T) {
	tests := []struct {
		a, b int
		want bool
	}{
		{3, 5, true},
		{4, 6, false},
		{1, 7, true},
		{14, 21, false},
	}
	for _, tt := range tests {
		got, err := numberCoprime(tt.a, tt.b)
		if err != nil {
			t.Errorf("numbers.coprime(%d, %d): unexpected error: %v", tt.a, tt.b, err)
			continue
		}
		if got != tt.want {
			t.Errorf("numbers.coprime(%d, %d) = %v, want %v", tt.a, tt.b, got, tt.want)
		}
	}
}
