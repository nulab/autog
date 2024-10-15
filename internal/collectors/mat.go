package collectors

type Mat[T any] [][]T

// NewMat creates a new n-by-n square matrix. All entries are T's zero value. A negative n is treated as 0.
func NewMat[T any](n int) Mat[T] {
	if n < 0 {
		n = 0
	}
	m := make([][]T, n)
	for i := range m {
		m[i] = make([]T, n)
	}
	return m
}
