package store

type Set[T comparable] struct {
	values map[T]struct{}
}

func NewSet[T comparable](values ...T) *Set[T] { // HL
	m := make(map[T]struct{}, len(values))
	for _, v := range values {
		m[v] = struct{}{}
	}
	return &Set[T]{
		values: m,
	}
}

func (s *Set[T]) Add(values ...T) { // HL
	for _, v := range values {
		s.values[v] = struct{}{}
	}
}
