package collectors

type Mat[T any] [][]T

func NewMat[T any](n int) Mat[T] {
	m := make([][]T, n)
	for i := range m {
		m[i] = make([]T, n)
	}
	return m
}
