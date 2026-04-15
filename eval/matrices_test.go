package eval

import "testing"

func TestMatricesDimensions(t *testing.T) {
	tests := []struct {
		name string
		g    []string
		r, c int
		want bool
	}{
		{"valid 2x3", []string{"abc", "def"}, 2, 3, true},
		{"wrong rows", []string{"abc"}, 2, 3, false},
		{"wrong cols", []string{"ab", "cd"}, 2, 3, false},
		{"mixed cols", []string{"abc", "de"}, 2, 3, false},
		{"empty", []string{}, 0, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := matricesDimensions(tt.g, tt.r, tt.c)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
