package list

import (
	"reflect"
	"testing"
)

func TestContains(t *testing.T) {
	tests := []struct {
		name string
		arr  []int
		v    int
		want bool
	}{
		{
			name: "empty array",
			v:    1,
			want: false,
		},
		{
			name: "array contains the value",
			arr:  []int{1, 2},
			v:    1,
			want: true,
		},
		{
			name: "array doesn't contain the value",
			arr:  []int{1, 3, 4},
			v:    2,
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Contains(tt.arr, tt.v); got != tt.want {
				t.Errorf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTakeN(t *testing.T) {
	tests := []struct {
		name   string
		arr    []int
		n      int
		offset int
		want   []int
	}{
		{
			name: "takes from first",
			arr:  []int{1, 2, 3},
			n:    2,
			want: []int{1, 2},
		},
		{
			name:   "takes from offset",
			arr:    []int{1, 2, 3},
			n:      3,
			offset: 1,
			want:   []int{2, 3},
		},
		{
			name:   "offset is greater than array",
			arr:    []int{1, 2, 3},
			n:      1,
			offset: 3,
			want:   nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TakeN(tt.arr, tt.n, tt.offset); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TakeN() = %v, want %v", got, tt.want)
			}
		})
	}
}
