package collectors

type stack[T any] struct {
	ts []T
}

func NewStack[T any](n int) *stack[T] {
	return &stack[T]{make([]T, 0, n)}
}

func (s *stack[T]) Push(vs ...T) {
	s.ts = append(s.ts, vs...)
}

func (s *stack[T]) Pop() (t T) {
	t = s.ts[s.Len()-1]
	s.ts = s.ts[:s.Len()-1]
	return
}

func (s *stack[T]) Peek(n int) T {
	return s.ts[s.Len()-n]
}

func (s *stack[T]) Len() int {
	return len(s.ts)
}
