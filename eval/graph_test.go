package eval

import "testing"

func TestGraphSimple(t *testing.T) {
	tests := []struct {
		name string
		n    int
		u, v []int
		want bool
	}{
		{"simple graph", 3, []int{1, 2}, []int{2, 3}, true},
		{"self loop", 3, []int{1, 2}, []int{1, 3}, false},
		{"duplicate edge", 3, []int{1, 1}, []int{2, 2}, false},
		{"duplicate reversed", 3, []int{1, 2}, []int{2, 1}, false},
		{"empty", 3, []int{}, []int{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := graphSimple(tt.n, tt.u, tt.v)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGraphConnected(t *testing.T) {
	tests := []struct {
		name string
		n    int
		u, v []int
		want bool
	}{
		{"connected", 3, []int{1, 2}, []int{2, 3}, true},
		{"disconnected", 4, []int{1, 3}, []int{2, 4}, false},
		{"single node", 1, []int{}, []int{}, true},
		{"empty graph", 0, []int{}, []int{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := graphConnected(tt.n, tt.u, tt.v)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGraphAcyclic(t *testing.T) {
	tests := []struct {
		name string
		n    int
		u, v []int
		want bool
	}{
		{"tree", 3, []int{1, 2}, []int{2, 3}, true},
		{"cycle", 3, []int{1, 2, 3}, []int{2, 3, 1}, false},
		{"no edges", 3, []int{}, []int{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := graphAcyclic(tt.n, tt.u, tt.v)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGraphTree(t *testing.T) {
	tests := []struct {
		name string
		n    int
		u, v []int
		want bool
	}{
		{"valid tree", 4, []int{1, 2, 3}, []int{2, 3, 4}, true},
		{"too few edges", 4, []int{1, 2}, []int{2, 3}, false},
		{"has cycle", 3, []int{1, 2, 3}, []int{2, 3, 1}, false},
		{"single node", 1, []int{}, []int{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := graphTree(tt.n, tt.u, tt.v)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGraphArgs_Invalid(t *testing.T) {
	if _, err := graphSimple(3); err == nil {
		t.Error("expected error for wrong arg count")
	}
	if _, err := graphSimple("3", []int{1}, []int{2}); err == nil {
		t.Error("expected error for non-int n")
	}
	if _, err := graphSimple(3, []int{1}, []int{2, 3}); err == nil {
		t.Error("expected error for mismatched u/v lengths")
	}
}
