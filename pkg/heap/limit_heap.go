package heap

type LimitHeap[T any] struct {
	limit int
	*Heap[T]
}

// NewLimitHeap initializes min-heap with the given limit.
func NewLimitHeap[T any](limit int, less func(a, b T) bool) *LimitHeap[T] {
	return &LimitHeap[T]{limit: limit, Heap: NewHeap(less)}
}

func (l *LimitHeap[T]) Push(v T) {
	if l.Len() < l.limit {
		l.Heap.Push(v)
	} else {
		if head := l.Peek(); l.less(*head, v) {
			*head = v
			l.Heap.down()
		}
	}
}
