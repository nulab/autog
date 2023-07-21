package stack

type stack[T any] []T

func New[T any](items ...T) stack[T] {
	s := stack[T]{}
	s.Push(items...)
	return s
}

func (s *stack[T]) Push(item ...T) {
	*s = append(*s, item...)
}

func (s *stack[T]) Pop() T {
	item := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return item
}
