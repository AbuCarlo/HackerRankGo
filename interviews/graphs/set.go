package graphs

type Set[V comparable] struct {
	m map[V]struct{}
}

func (s *Set[V]) Has(v V) bool {
	_, ok := s.m[v]
	return ok
}

func (s *Set[V]) Add(v V) {
	s.m[v] = struct{}{}
}

func (s *Set[V]) Remove(v V) {
	delete(s.m, v)
}

func (s *Set[V]) Clear() {
	s.m = make(map[V]struct{})
}

func (s *Set[V]) Size() int {
	return len(s.m)
}