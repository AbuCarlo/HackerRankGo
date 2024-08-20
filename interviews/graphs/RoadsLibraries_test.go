package graphs

// See https://en.wikipedia.org/wiki/Disjoint-set_data_structure
// See

import (
	// "testing"
	"math/rand"
	"testing"

	"github.com/abucarlo/hackerrank/interviews/graphs/sets"
	"pgregory.net/rapid"
)

type DisjointSets map[int]*sets.Set[int]

type UndirectedGraph struct {
	// TODO: Rename "parents"
	parents map[int]int
	adjacency map[int]*sets.Set[int]
}

func NewUndirectedGraph() *UndirectedGraph {
	g := UndirectedGraph{ make(map[int]int), make(map[int]*sets.Set[int]) }
	return &g
}

func (g *UndirectedGraph) FindDisjoint() []UndirectedGraph {
	// To get distinct values in a slice:
	// https://stackoverflow.com/a/76471309/476942
	// To get a map's keys as a slice, see
	// https://stackoverflow.com/a/69889828/476942
	// I don't want to exploit any of these options here,
	// in case HackerRank doesn't let me.
	// See mainly https://en.wikipedia.org/wiki/Kruskal%27s_algorithm

	disjoints := make(DisjointSets)
	// Now give each subgraph its adjacency matrix.
	for v, _ := range g.parents {
		root := g.FindRoot(v)
		s, has := disjoints[root]
		if !has {
			s = sets.New[int]()
			disjoints[root] = s
		}
		s.Add(v)
		s.Add(root)
		disjoints[root] = s
	}

	var trees []UndirectedGraph 
	for _, s := range disjoints {
		// Construct a new subgraph
		// in which every adjacency list
		// has a size of 1 or 2.
		mst := NewUndirectedGraph()
		seen := sets.New[int]()
		for _, u := range s.Items() {
			if seen.Has(u) {
				continue
			}
			for _, v := range g.adjacency[u].Items() {
				if seen.Has(v) {
					continue
				}
				mst.Insert(u, v)
			}
		}
		
		trees = append(trees, *mst)
	}

	return trees
}

func (g *UndirectedGraph) Insert(u, v int) {
	if u == v {
		panic("u must not == v")
	}
	// Since it's an undirected graph, let's just
	// decide to make the the lower vertex number
	// the root.
	// TODO: Collapse these tests.
	if _, ok := g.parents[v]; !ok {
		g.parents[v] = g.FindRoot(u)
	}
	if _, ok := g.parents[v]; !ok {
		g.parents[v] = v
	}
	if _, ok := g.adjacency[u]; !ok {
		g.adjacency[u] = sets.New[int]()
	}
	if _, ok := g.adjacency[v]; !ok {
		g.adjacency[v] = sets.New[int]()
	}
	g.adjacency[u].Add(v)
	g.adjacency[v].Add(u)
}

func (g *UndirectedGraph) FindRoot(v int) int {
	u, ok := g.parents[v]
	if !ok {
		g.parents[v] = v
		return v
	}
	// Follow the path to the root, compressing all the while.
	// See https://en.wikipedia.org/wiki/Disjoint-set_data_structure#Finding_set_representatives
	for u != g.parents[u] {
		u, g.parents[u] = g.parents[u], g.parents[g.parents[u]]
	}

	return u
}

func TestPathGraph(t *testing.T) {
	f := func(t *rapid.T) {
		path := rapid.Custom[[]int](func(t *rapid.T) []int {
			size:= rapid.IntRange(2, 9999).Draw(t, "size")
			seed := rapid.Int64().Draw(t, "seed")
			result := make([]int, size)
			for i := range result {
				result[i] = i + 1
			}
			source := rand.NewSource(seed)
			r := rand.New(source)
			r.Shuffle(len(result), func(i, j int) { result[i], result[j] = result[j], result[i] })
			return result
		})

		blah := path.Draw(t, "path")
		
		graph := NewUndirectedGraph()
		for i, u := range blah {
			if i == 0 {
				continue
			}
			graph.Insert(blah[i - 1], u)
		}

		fart := graph.FindDisjoint()
		if len(fart) != 1 {
			t.Errorf("A path graph should have 1 disjoint, not %d", len(blah))
		}
	}

	rapid.Check(t, f)
}

func TestStarGraph(t *testing.T) {

	f := func(t *rapid.T) {
		n := rapid.IntRange(1, 9999).Draw(t, "size")
		v := rapid.IntRange(1, n).Filter(func (m int) bool { return m != n }).Draw(t, "axis")
		
		graph := NewUndirectedGraph()
		for u := 1; u <= n; u++ {
			if u != v {
				graph.Insert(u, v)
			}
		}
	
		d := graph.FindDisjoint()
		if len(d) != 1 {
			t.Errorf("A star graph of %d nodes centered on %d should have 1 connected subgraph, not %d", n, v, len(d))
		}
	}
	rapid.Check(t, f)
}