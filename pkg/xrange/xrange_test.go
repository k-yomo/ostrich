package xrange

import "testing"

func TestRange_Len(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		r    Range
		want int
	}{
		{
			r:    Range{From: 0, To: 4},
			want: 5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.Len(); got != tt.want {
				t.Errorf("Len() = %v, want %v", got, tt.want)
			}
		})
	}
}
