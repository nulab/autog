package collectors

type Deque[T any] struct {
	data []T
	f    int // front index
	b    int // back index
}

func NewDeque[T any](size int) *Deque[T] {
	return &Deque[T]{
		data: make([]T, size*2),
		f:    size,
		b:    size - 1,
	}
}

func (d *Deque[_]) Len() int {
	return d.b - d.f + 1
}

// PushFront pushes a new item to the front of the queue.
// This grows the queue to the left toward the 0 index.
func (d *Deque[T]) PushFront(x T) {
	d.f--
	d.data[d.f] = x
}

// PushBack pushes a new item to the back of the queue.
// This grows the queue to the right toward Len()-1 index.
func (d *Deque[T]) PushBack(x T) {
	d.b++
	d.data[d.b] = x
}

func (d *Deque[T]) PeekFront(i int) T {
	return d.data[d.f+i-1]
}

func (d *Deque[T]) PeekBack(i int) T {
	return d.data[d.b-i+1]
}

func (d *Deque[T]) PopFront() T {
	i := d.f
	d.f++
	return d.data[i]
}

func (d *Deque[T]) PopBack() T {
	i := d.b
	d.b--
	return d.data[i]
}

func (d *Deque[T]) Front() int {
	return d.f
}

func (d *Deque[T]) Back() int {
	return d.b
}
