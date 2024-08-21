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

func (s *Set[V]) Union(t *Set[V]) {
	for v := range t.m {
		s.Add(v)
	}
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

func (s *Set[V]) Items() []V {
	// Man, it would be nice to use maps.Keys() or write an iterator here.
	items := make([]V, 0, s.Size())
	for v := range s.m {
		items = append(items, v)
	}
	return items
}

func NewSet[V comparable]() *Set[V] {
	s := Set[V]{}
	s.m = make(map[V]struct{})
	return &s
}