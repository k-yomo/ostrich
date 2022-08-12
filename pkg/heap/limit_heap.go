package heap

type LimitHeap[T any] struct {
	limit int
	*Heap[T]
}

func NewLimitHeap[T any](limit int, comp func(a, b T) bool) *LimitHeap[T] {
	return &LimitHeap[T]{limit: limit, Heap: NewHeap(comp)}
}

func (l *LimitHeap[T]) Push(v T) {
	l.Heap.Push(v)
	if l.Len() > l.limit {
		_ = l.Heap.Pop()
	}
}
