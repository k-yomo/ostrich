package heap

type Heap[T any] struct {
	data []T
	less func(a, b T) bool
}

// NewHeap initializes new min-heap with given comparing function
func NewHeap[T any](less func(a, b T) bool) *Heap[T] {
	return &Heap[T]{less: less}
}

func (h *Heap[T]) Len() int { return len(h.data) }

func (h *Heap[T]) Push(v T) {
	h.data = append(h.data, v)
	h.up(h.Len() - 1)
}

func (h *Heap[T]) Pop() T {
	n := h.Len() - 1
	if n > 0 {
		h.swap(0, n)
		h.down()
	}
	v := h.data[n]
	h.data = h.data[0:n]
	return v
}

func (h *Heap[T]) Peek() *T {
	return &(h.data[0])
}

// ToArray map heap to array in ascending order
// this will pop all the values from the heap
func (h *Heap[T]) ToArray() []T {
	values := make([]T, 0, len(h.data))
	for h.Len() > 0 {
		values = append([]T{h.Pop()}, values...)
	}
	return values
}

func (h *Heap[T]) swap(i, j int) {
	h.data[i], h.data[j] = h.data[j], h.data[i]
}

func (h *Heap[T]) up(j int) {
	for {
		i := parent(j)
		if i == j || !h.less(h.data[j], h.data[i]) {
			break
		}
		h.swap(i, j)
		j = i
	}
}

func (h *Heap[T]) down() {
	n := h.Len() - 1
	i1 := 0
	for {
		j1 := left(i1)
		if j1 >= n || j1 < 0 {
			break
		}
		j := j1
		j2 := right(i1)
		if j2 < n && h.less(h.data[j2], h.data[j1]) {
			j = j2
		}
		if !h.less(h.data[j], h.data[i1]) {
			break
		}
		h.swap(i1, j)
		i1 = j
	}
}

func parent(i int) int { return (i - 1) / 2 }
func left(i int) int   { return (i * 2) + 1 }
func right(i int) int  { return left(i) + 1 }
