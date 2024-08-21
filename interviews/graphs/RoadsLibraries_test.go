package graphs

// See https://en.wikipedia.org/wiki/Disjoint-set_data_structure
// See

import (
	"math/rand"
	"testing"

	"pgregory.net/rapid"
)

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

type DisjointSets map[int32]*Set[int32]

type UndirectedGraph struct {
	adjacency map[int32]*Set[int32]
}

func NewUndirectedGraph() *UndirectedGraph {
	g := UndirectedGraph{make(map[int32]*Set[int32])}
	return &g
}

// "size" / "order" cf. https://en.wikipedia.org/wiki/Graph_(discrete_mathematics)#Graph

func (g *UndirectedGraph) Order() int32 {
	return int32(len(g.adjacency))
}

func (g *UndirectedGraph) Size() int32 {
	result := 0
	for _, s := range g.adjacency {
		result += s.Size()
	}
	result /= 2
	return int32(result)
}

func FindRoot(parents map[int32]int32, v int32) int32 {
	parent := parents[v]
	// Follow the path to the root, compressing all the while.
	// See https://en.wikipedia.org/wiki/Disjoint-set_data_structure#Finding_set_representatives
	for parent != parents[parent] {
		parent, parents[parent] = parents[parent], parents[parents[parent]]
	}

	return parent
}

func (g *UndirectedGraph) FindDisconnected() []UndirectedGraph {
	// To get distinct values in a slice:
	// https://stackoverflow.com/a/76471309/476942
	// To get a map's keys as a slice, see
	// https://stackoverflow.com/a/69889828/476942
	// I don't want to exploit any of these options here,
	// in case HackerRank doesn't let me.
	// See mainly https://en.wikipedia.org/wiki/Kruskal%27s_algorithm
	parents := make(map[int32]int32)
	disjoints := make(DisjointSets)

	for v := range g.adjacency {
		parents[v] = v
		disjoints[v] = NewSet[int32]()
		disjoints[v].Add(v)
	}

	// Now give each subgraph its adjacency matrix.
	for u, s := range g.adjacency {
		for _, v := range s.Items() {
			x := FindRoot(parents, u)
			y := FindRoot(parents, v)
			if x == y {
				continue
			}
			if disjoints[x].Size() < disjoints[y].Size() {
				x, y = y, x
			}
			parents[y] = x
			disjoints[x].Union(disjoints[y])
			delete(disjoints, y)
		}
	}

	var trees []UndirectedGraph
	for _, s := range disjoints {
		// Construct a new subgraph
		// in which every adjacency list
		// has a size of 1 or 2.
		mst := NewUndirectedGraph()
		sparents := make(map[int32]int32)
		sdisjoints := make(DisjointSets)
		for _, v := range s.Items() {
			sparents[v] = v
			sdisjoints[v] = NewSet[int32]()
			sdisjoints[v].Add(v)
		}
		for _, u := range s.Items() {
			for _, v := range g.adjacency[u].Items() {
				x, y := FindRoot(sparents, u), FindRoot(sparents, v)
				if x == y {
					continue
				}
				if sdisjoints[x].Size() < sdisjoints[y].Size() {
					x, y = y, x
				}
				sparents[y] = x
				sdisjoints[x].Union(sdisjoints[y])
				delete(sdisjoints, y)
				mst.Insert(u, v)
			}
		}

		trees = append(trees, *mst)
	}

	return trees
}

func (g *UndirectedGraph) Insert(u, v int32) {
	if u == v {
		panic("u must not == v")
	}
	// Since it's an undirected graph, let's just
	// decide to make the the lower vertex number
	// the parent.
	if u > v {
		u, v = v, u
	}

	if _, ok := g.adjacency[v]; !ok {
		g.adjacency[v] = NewSet[int32]()
	}

	if _, ok := g.adjacency[u]; !ok {
		g.adjacency[u] = NewSet[int32]()
	}

	g.adjacency[u].Add(v)
	g.adjacency[v].Add(u)
}


func roadsAndLibraries(order int32, library int32, road int32, edges [][]int32) int64 {
	graph := NewUndirectedGraph()
	for _, edge := range edges {
		u, v := edge[0], edge[1]
		graph.Insert(u, v)
	}

	trees := graph.FindDisconnected()
	if library > road {
		result := int64(0)
		for _, t := range trees {
			result += int64(library) + (int64(len(t.adjacency)-1))*int64(road)
		}
		// Correction: there might be vertices with no edges.
		disconnected := order - int32(len(graph.adjacency))
		result += int64(disconnected * library)
		return result
	}

	return int64(library) * int64(order)
}

// https://en.wikipedia.org/wiki/Path_graph
func TestPathGraph(t *testing.T) {
	f := func(t *rapid.T) {
		path := rapid.Custom[[]int32](func(t *rapid.T) []int32 {
			order := rapid.Int32Range(2, 9999).Draw(t, "order")
			seed := rapid.Int64().Draw(t, "seed")
			result := make([]int32, order)
			for i := range result {
				result[i] = int32(i + 1)
			}
			source := rand.NewSource(seed)
			r := rand.New(source)
			r.Shuffle(len(result), func(i, j int) { result[i], result[j] = result[j], result[i] })
			return result
		})

		vertices := path.Draw(t, "path")

		graph := NewUndirectedGraph()
		for i, u := range vertices {
			if i == 0 {
				continue
			}
			graph.Insert(vertices[i-1], u)
		}

		trees := graph.FindDisconnected()
		if len(trees) != 1 {
			t.Errorf("A path graph should have 1 connected component, not %d", len(vertices))
		}
		if trees[0].Order() - 1 != trees[0].Size() {
			t.Errorf("In an MST, |V| - 1 == |E| (got %d, %d)", trees[0].Order(), trees[0].Size())
		}
	}

	rapid.Check(t, f)
}

// https://en.wikipedia.org/wiki/Star_(graph_theory)
func TestStarGraph(t *testing.T) {

	f := func(t *rapid.T) {
		order := rapid.Int32Range(1, 9999).Draw(t, "order")
		v := rapid.Int32Range(1, order).Filter(func(m int32) bool { return m != order }).Draw(t, "axis")

		graph := NewUndirectedGraph()
		for u := int32(1); u <= order; u++ {
			if u != v {
				graph.Insert(u, v)
			}
		}

		trees := graph.FindDisconnected()
		if len(trees) != 1 {
			t.Errorf("A star graph of %d nodes centered on %d should have 1 connected subgraph, not %d", order, v, len(trees))
		}
		if trees[0].Order() - 1 != trees[0].Size() {
			t.Errorf("In an MST, |V| - 1 == |E| (got %d, %d)", trees[0].Order(), trees[0].Size())
		}
	}
	rapid.Check(t, f)
}
func TestStronglyConnected(t *testing.T) {
	f := func(t *rapid.T) {
		// There will not be more than 100000 edges.
		order := rapid.Int32Range(2, 316).Draw(t, "order")
		graph := NewUndirectedGraph()
		for u := int32(1); u <= order; u++ {
			for v := int32(1); v <= order; v++ {
				if u != v {
					graph.Insert(u, v)
				}
			}
		}

		d := graph.FindDisconnected()
		if len(d) != 1 {
			t.Errorf("A strongly connected graph of %d nodes should have 1 connected subgraph, not %d", order, len(d))
		}
		for _, tree := range d {
			if tree.Size() != tree.Order()-1 {
				t.Errorf("MST %v has order %d, size %d", tree, tree.Order(), tree.Size())
			}
		}
	}
	rapid.Check(t, f)
}

func TestSamples(t *testing.T) {
	type Test struct {
		n        int32
		library  int32
		road     int32
		vertices [][]int32
		expected int64
	}

	tests := []Test{
		{3, 2, 1, [][]int32{{1, 2}, {3, 1}, {2, 3}}, 4},
		{6, 2, 5, [][]int32{{1, 3}, {3, 4}, {2, 4}, {1, 2}, {2, 3}, {5, 6}}, 12},
		{5, 6, 1, [][]int32{{1, 2}, {1, 3}, {1, 4}}, 15},
		// // Test Case 1
		{9, 91, 84, [][]int32{{8, 2}, {2, 9}}, 805},
		{5, 92, 23, [][]int32{{2, 1},
			{5, 3},
			{5, 1},
			{3, 4},
			{3, 1},
			{5, 4},
			{4, 1},
			{5, 2},
			{4, 2}}, 184},
		{8, 10, 55, [][]int32{{6, 4}, {3, 2}, {7, 1}}, 80},
		{1, 5, 3, [][]int32{}, 5},
		{2, 102, 1, [][]int32{}, 204},
	}

	for _, test := range tests {
		actual := roadsAndLibraries(test.n, test.library, test.road, test.vertices)
		if actual != test.expected {
			t.Errorf("Expected %d; got %d", test.expected, actual)
		}
	}
}
