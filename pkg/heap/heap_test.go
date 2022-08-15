package heap

import (
	"reflect"
	"testing"
)

func TestHeap(t *testing.T) {
	t.Parallel()

	h := NewHeap(func(a, b int) bool {
		return a < b
	})
	h.Push(1)
	h.Push(3)
	h.Push(2)
	h.Push(4)

	if got := h.Len(); !reflect.DeepEqual(got, 4) {
		t.Errorf("Heap_Len() = %v, want %v", got, 4)
	}
	if got := h.Pop(); !reflect.DeepEqual(got, 1) {
		t.Errorf("Heap_Pop() = %v, want %v", got, 1)
	}
	if got := h.Len(); !reflect.DeepEqual(got, 3) {
		t.Errorf("Heap_Len() = %v, want %v", got, 3)
	}
	if got := *h.Peek(); !reflect.DeepEqual(got, 2) {
		t.Errorf("Heap_Peek() = %v, want %v", got, 2)
	}
	if got := h.ToArray(); !reflect.DeepEqual(got, []int{4, 3, 2}) {
		t.Errorf("Heap_ToArray() = %v, want %v", got, []int{4, 3, 2})
	}
}
