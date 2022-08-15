package heap

import (
	"reflect"
	"testing"
)

func TestLimitHeap(t *testing.T) {
	t.Parallel()

	h := NewLimitHeap(3, func(a, b int) bool {
		return a < b
	})
	h.Push(3)
	h.Push(1)
	h.Push(4)
	h.Push(2)

	if !reflect.DeepEqual(h.Len(), 3) {
		t.Errorf("Heap_Len() = %v, want %v", h.Len(), 3)
	}
	if !reflect.DeepEqual(*h.Peek(), 2) {
		t.Errorf("Heap_Peek() = %v, want %v", *h.Peek(), 2)
	}
	if got := h.Pop(); !reflect.DeepEqual(got, 2) {
		t.Errorf("Heap_Pop() = %v, want %v", got, 2)
	}
	if !reflect.DeepEqual(h.Len(), 2) {
		t.Errorf("Heap_Len() = %v, want %v", h.Len(), 2)
	}
	if !reflect.DeepEqual(*h.Peek(), 3) {
		t.Errorf("Heap_Peek() = %v, want %v", *h.Peek(), 3)
	}
}
