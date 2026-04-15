package eval

import "testing"

func TestArraysSorted(t *testing.T) {
	tests := []struct {
		name string
		arg  interface{}
		want bool
	}{
		{"sorted ints", []int{1, 2, 3}, true},
		{"unsorted ints", []int{3, 1, 2}, false},
		{"equal ints", []int{1, 1, 1}, true},
		{"sorted int64s", []int64{1, 2, 3}, true},
		{"sorted float64s", []float64{1.1, 2.2, 3.3}, true},
		{"sorted strings", []string{"a", "b", "c"}, true},
		{"unsorted strings", []string{"c", "a", "b"}, false},
		{"single element", []int{42}, true},
		{"empty", []int{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := arraysSorted(tt.arg)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArraysDistinct(t *testing.T) {
	tests := []struct {
		name string
		arg  interface{}
		want bool
	}{
		{"distinct ints", []int{1, 2, 3}, true},
		{"duplicate ints", []int{1, 2, 1}, false},
		{"distinct strings", []string{"a", "b", "c"}, true},
		{"duplicate strings", []string{"a", "b", "a"}, false},
		{"single", []int{1}, true},
		{"empty", []int{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := arraysDistinct(tt.arg)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArraysPermutation(t *testing.T) {
	tests := []struct {
		name string
		n    int
		a    []int
		want bool
	}{
		{"valid", 3, []int{2, 3, 1}, true},
		{"identity", 3, []int{1, 2, 3}, true},
		{"wrong length", 3, []int{1, 2}, false},
		{"out of range", 3, []int{1, 2, 4}, false},
		{"duplicate", 3, []int{1, 2, 2}, false},
		{"zero", 3, []int{0, 1, 2}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := arraysPermutation(tt.n, tt.a)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArraysRange(t *testing.T) {
	tests := []struct {
		name string
		a    interface{}
		lo   interface{}
		hi   interface{}
		want bool
	}{
		{"in range ints", []int{1, 2, 3}, 1, 3, true},
		{"out of range ints", []int{1, 2, 4}, 1, 3, false},
		{"in range int64s", []int64{10, 20}, int64(5), int64(25), true},
		{"in range float64s", []float64{1.5, 2.5}, 1.0, 3.0, true},
		{"empty", []int{}, 1, 10, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := arraysRange(tt.a, tt.lo, tt.hi)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArrays_InvalidArg(t *testing.T) {
	if _, err := arraysSorted("not an array"); err == nil {
		t.Error("expected error for non-array arg")
	}
	if _, err := arraysDistinct(42); err == nil {
		t.Error("expected error for non-array arg")
	}
}
