package heap

type LimitHeap[T any] struct {
	limit int
	*Heap[T]
}

func NewLimitHeap[T any](limit int, comp func(a, b T) bool) *LimitHeap[T] {
	return &LimitHeap[T]{limit: limit, Heap: NewHeap(comp)}
}

func (l *LimitHeap[T]) Push(v T) {
	if l.Len() < l.limit {
		l.Heap.Push(v)
	} else {
		if last := l.Peek(); l.comp(v, *last) {
			l.Heap.Push(v)
			_ = l.Heap.Pop()
		}
	}
}
